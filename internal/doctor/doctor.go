// Package doctor checks a ccam installation and repairs what it safely can.
package doctor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/usage"
)

type Level int

const (
	OK Level = iota
	Info
	Warn
	Error
)

type Finding struct {
	Level   Level
	Profile string
	Area    string // what was checked: "Login", "Share links", "Claude Code"…
	Msg     string
	Fixed   bool
}

type Options struct {
	Fix bool
	// FixOnly limits Fix to the findings of one account, or with General to
	// the findings of no account. Empty repairs everything.
	FixOnly  string
	CcamPath string         // the real ccam binary, for shim checks
	Usage    *usage.Fetcher // nil skips the live usage-format check
	// DryRun, for `--dry-run doctor --fix`: say what repairing the share
	// links would set aside, without changing anything.
	DryRun bool
	// Reload loads the config again once Fix holds the lock of the changes
	// of groups: one loaded before may miss a change another ccam finished
	// meanwhile. nil keeps the config Run was given.
	Reload func() (*config.Config, error)
}

// General is FixOnly for the findings of no account; no account can be
// named so (config names start with a letter or digit).
const General = "(general)"

// fixes reports whether Run repairs the findings of profile prof ("" for
// those of no account).
func (o Options) fixes(prof string) bool {
	switch {
	case !o.Fix:
		return false
	case o.FixOnly == General:
		return prof == ""
	}
	return o.FixOnly == "" || o.FixOnly == prof
}

// LayoutVersion is the Claude Code version docs/claude-layout.md was checked with.
const LayoutVersion = "2.1.309"

type report struct {
	out  []Finding
	area string // the Area of the findings added next
}

func (r *report) add(level Level, prof, format string, a ...any) *Finding {
	r.out = append(r.out, Finding{Level: level, Profile: prof, Area: r.area, Msg: fmt.Sprintf(format, a...)})
	return &r.out[len(r.out)-1]
}

// links says how checkProfile may repair an account's share links.
type links int

const (
	linksRepair links = iota // toward the config
	// A change of groups that stopped halfway repairs them in its order:
	// only what is missing is made now, and a dry run shows the change's
	// plan.
	linksWait
	// Another ccam changes the groups right now: nothing is repaired.
	linksLeave
)

// Run checks everything and, with opt.Fix, repairs safe problems.
func Run(p paths.Paths, cfg *config.Config, opt Options) []Finding {
	r := &report{}
	claude := checkClaude(r, p, cfg)
	checkAuthEnv(r)
	// Share links are repaired, and a stopped change of groups finished,
	// under the lock a change of groups holds, toward the config saved by
	// then.
	var held *group.Lock
	busy := false
	if opt.Fix {
		l, err := group.TryLock(p)
		switch {
		case err == nil:
			defer l.Unlock()
			held = l
			if opt.Reload != nil {
				if cfg, err = opt.Reload(); err != nil {
					r.area = "Config"
					r.add(Error, "", "could not load the config again: %v", err)
					return r.out
				}
			}
		case errors.Is(err, group.ErrChanging):
			busy = true
		default:
			busy = true
			r.area = "Groups"
			r.add(Error, "", "share links are not repaired: %v", err)
		}
	}
	// Before the accounts' links, which it makes in the change's order.
	noted, mode, previews := checkPendingChange(r, p, cfg, held, opt.fixes(""), opt.DryRun)
	// For a dry run: what the repairs of the accounts before would have put
	// into the groups, as share.Options.Pending; the stopped change comes
	// first.
	pending := map[string]string{}
	for _, name := range noted {
		if pl := previews[name]; pl != nil {
			addCreated(pending, pl)
		}
	}
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		m := linksRepair
		switch {
		case busy:
			m = linksLeave
		case slices.Contains(noted, prof.Name):
			m = mode
		}
		checkProfile(r, p, prof, opt.fixes(prof.Name), opt.DryRun, m, previews[prof.Name], pending)
	}
	checkPlugins(r, p, cfg, opt.fixes("") && !busy)
	checkMCP(r, p, cfg)
	checkUsageFormat(r, p, cfg, opt.Usage)
	checkShims(r, p, cfg, opt)
	checkFolders(r, cfg)
	checkGroups(r, cfg)
	checkIdentities(r, p, cfg)
	checkAutoSwitch(r, p, cfg, claude, runtime.GOOS)
	return r.out
}

// checkPendingChange reports a change of groups that stopped halfway, and
// with fix (held is the lock) finishes it: what the plans toward the saved
// config cannot show, the MCP servers it copies, the answers given to its
// conflicts, and the order of its accounts, which decides the first
// content of a group (--from). A dry run previews it instead. It returns
// the accounts of a change it did not finish, how checkProfile may repair
// their links, and for a dry run the plans the change would carry out.
func checkPendingChange(r *report, p paths.Paths, cfg *config.Config, held *group.Lock, fix, dry bool) ([]string, links, map[string]*share.Plan) {
	if !group.Pending(p) {
		return nil, linksRepair, nil
	}
	r.area = "Groups"
	if held == nil && group.Changing(p) {
		// The note stays until that ccam is done.
		r.add(Info, "", "a group change is in progress right now; its accounts' share links are left to it")
		return group.Noted(p), linksLeave, nil
	}
	f := r.add(Error, "", "the last group change stopped halfway (ccam doctor --fix finishes it; noted in %s)", group.NotePath(p))
	if dry {
		ch, err := group.Preview(p, cfg)
		if err != nil {
			r.add(Error, "", "could not plan the rest of the last group change: %v", err)
			return group.Noted(p), linksWait, nil
		}
		previews := map[string]*share.Plan{}
		if ch != nil && len(ch.Accounts) > 0 {
			r.add(Info, "", "--fix finishes it in its order (%s), with the answers given to its conflicts", strings.Join(ch.Names(), ", "))
			for _, a := range ch.Accounts {
				previews[a.Name] = a.Plan
			}
		}
		return group.Noted(p), linksWait, previews
	}
	if !fix || held == nil {
		return group.Noted(p), linksWait, nil
	}
	noted := group.Noted(p)
	ch, err := group.Resume(held, cfg, group.Running(p))
	if err != nil {
		r.add(Error, "", "could not finish the last group change: %v\n"+
			"If the cause cannot be fixed, removing %s gives up the rest of that change; ccam doctor --fix then repairs the links", err, group.NotePath(p))
		return noted, linksWait, nil
	}
	f.Fixed = true
	if ch == nil || len(ch.Accounts) == 0 {
		r.add(OK, "", "finished the last group change (none of its accounts is in the config anymore)")
		return nil, linksRepair, nil
	}
	r.add(OK, "", "finished the last group change (%s)", strings.Join(ch.Names(), ", "))
	for _, a := range ch.Accounts {
		setAside(r, a.Name, a.Plan, true)
	}
	return nil, linksRepair, nil
}

// addCreated adds what plan pl puts into the groups to pending, for the
// dry-run plans of the accounts after it (share.Options.Pending).
func addCreated(pending map[string]string, pl *share.Plan) {
	for k, v := range pl.Created() {
		if old, ok := pending[k]; !ok || old == "" {
			pending[k] = v
		}
	}
}

// checkIdentities reports a folder now logged in to another account than
// ccam recorded (a /login inside a session), folders logged in to one
// Claude account, which share its limits, and folders holding one login
// (copies, also of plain claude's ~/.claude), which log each other out.
func checkIdentities(r *report, p paths.Paths, cfg *config.Config) {
	r.area = "Login"
	seen := map[string]string{}   // account and organization → first profile
	logins := map[string]string{} // login token (creds.LoginKey) → first profile
	if dir := profile.PlainClaude(p, cfg); dir != "" {
		if k := creds.LoginKey(dir); k != "" {
			logins[k] = profile.PlainClaudeName
		}
	}
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if prof.API != nil {
			continue
		}
		dir := profile.Dir(p, prof)
		if k := creds.LoginKey(dir); k != "" {
			if first, ok := logins[k]; ok {
				r.add(Error, prof.Name, "holds the same login token as %s (a copied folder): whichever renews it first logs the other out. Run `ccam login %s` to give it its own (not `ccam logout %s`, which would log %s out too)", first, prof.Name, prof.Name, first)
				continue
			}
			logins[k] = prof.Name
		}
		id, err := profile.ReadIdentity(dir)
		if err != nil || id == nil {
			continue
		}
		recorded := ""
		switch was := prof.Login; {
		case was != nil && was.AccountUUID != "" && was.AccountUUID != id.AccountUUID:
			recorded = orDash(was.Email)
		case was == nil && prof.Email != "" && id.Email != "" && !strings.EqualFold(prof.Email, id.Email):
			// No login record (an older ccam's account, an imported one):
			// the email it recorded tells.
			recorded = prof.Email
		}
		if recorded != "" {
			r.add(Warn, prof.Name, "now logged in as %s, but ccam recorded %s; log in again with `ccam login %s` if that is not what you meant", orDash(id.Email), recorded, prof.Name)
		}
		k := id.AccountUUID + "/" + id.OrgUUID
		if first, ok := seen[k]; ok {
			r.add(Warn, prof.Name, "logged in to the same Claude account as %s (%s); the two share their usage limits", first, orDash(id.Email))
			continue
		}
		seen[k] = prof.Name
	}
}

// checkGroups reports group entries that do nothing, and account names that
// became ccam commands since the account was added.
func checkGroups(r *report, cfg *config.Config) {
	r.area = "Groups"
	for _, g := range cfg.Groups {
		for _, m := range g.Members {
			if !cfg.IsAccount(m) {
				r.add(Warn, "", "group %s lists %s, which is not an account (ccam group leave %s %s)", g.Name, m, g.Name, m)
			}
		}
	}
	r.area = "Config"
	for _, u := range cfg.Unknown() {
		r.add(Warn, "", "config.toml %s; commands that change the config refuse until it is fixed", u)
	}
	if s := cfg.UsageSource; s != "" && s != "live" && s != "claude" {
		r.add(Warn, "", "usage_source %q is neither live nor claude; ccam asks the usage endpoint (live)", s)
	}
	if t := cfg.Theme; t != "" && t != "auto" && !slices.Contains(config.Themes, t) {
		r.add(Warn, "", "theme %q is not one ccam knows (auto, %s); the screens follow the terminal's background", t, strings.Join(config.Themes, ", "))
	}
	r.area = "Name"
	for _, prof := range cfg.Profiles {
		if config.Reserved[prof.Name] {
			r.add(Warn, prof.Name, "`ccam %s` is now a ccam command; start this account with `ccam run %s` or rename it (ccam rename %s <new-name>)", prof.Name, prof.Name, prof.Name)
		}
	}
}

// checkAPIKey reports where an API key account gets its key. An env:
// source is read to tell whether it is set; a command: source is not run,
// since it may ask for a password.
func checkAPIKey(r *report, name string, api *config.API) {
	if env, ok := strings.CutPrefix(api.KeyFrom, "env:"); ok {
		if _, err := launcher.ReadAPIKey(api.KeyFrom); err != nil {
			r.add(Warn, name, "%s; $%s is empty, so the account cannot start", api.Label(), env)
			return
		}
	}
	r.add(OK, name, "%s, key from %s", api.Label(), api.KeyFrom)
}

// checkFolders reports `ccam folder` entries that no longer work. They are
// left for the user: removing one would change which account starts there.
func checkFolders(r *report, cfg *config.Config) {
	r.area = "Folder rules"
	for _, f := range cfg.Folders {
		if _, ok := cfg.Find(f.Account); !ok {
			r.add(Error, "", "%s is set to account %q, which does not exist (ccam folder <account> %s, or ccam folder --remove %s)", f.Path, f.Account, f.Path, f.Path)
			continue
		}
		if info, err := os.Stat(f.Path); err != nil || !info.IsDir() {
			r.add(Warn, f.Account, "the folder %s set for this account is missing (ccam folder --remove %s)", f.Path, f.Path)
		}
	}
}

// Failed reports whether any finding is still an unfixed error.
func Failed(fs []Finding) bool {
	for _, f := range fs {
		if f.Level == Error && !f.Fixed {
			return true
		}
	}
	return false
}

// checkClaude finds Claude Code and returns what `claude --version` printed
// ("2.1.310 (Claude Code)"), or "" when it could not.
func checkClaude(r *report, p paths.Paths, cfg *config.Config) string {
	r.area = "Claude Code"
	bin, err := launcher.ResolveClaude(p, cfg.ClaudeBin)
	if err != nil {
		r.add(Error, "", "%v", err)
		return ""
	}
	c, err := launcher.ClaudeCommand(context.Background(), bin, "--version")
	if err != nil {
		r.add(Error, "", "%v", err)
		return ""
	}
	out, err := c.Output()
	if err != nil {
		r.add(Error, "", "%s --version failed: %v", bin, err)
		return ""
	}
	version := strings.TrimSpace(string(out))
	r.add(OK, "", "Claude Code %s (%s)", version, bin)
	return version
}

// checkAutoSwitch says, with auto_switch = true, what it does and, for the
// plans that have seats among the accounts, what Anthropic's help center
// says about seats (D7 b). It warns about what keeps a session where it is:
// an account without the ccam plugin turned on, whose hook ends Claude Code,
// and a Claude Code before 2.1.214, which tells hooks no CLAUDE_PID. Windows
// ignores auto_switch (D6). claude is what `claude --version` printed.
func checkAutoSwitch(r *report, p paths.Paths, cfg *config.Config, claude, goos string) {
	if !cfg.AutoSwitch {
		return
	}
	r.area = "auto_switch"
	if goos == "windows" {
		r.add(Info, "", "auto_switch is ignored on Windows; with relay = true, /ccam:switch followed by /exit still moves a conversation to another account.")
		return
	}
	r.add(Info, "", "auto_switch is on: at a usage limit, ccam moves an interactive session it started to another account like the one it leaves (a seat of the same organization, or another personal plan) with room. Whether that fits Anthropic's terms is your call: https://ccam.dev/en/features/next/#auto-switch")
	var v [3]int
	if n, _ := fmt.Sscanf(claude, "%d.%d.%d", &v[0], &v[1], &v[2]); n == 3 && slices.Compare(v[:], []int{2, 1, 214}) < 0 {
		r.add(Warn, "", "auto_switch needs Claude Code 2.1.214 or later, which tells hooks its process (CLAUDE_PID); this one is %d.%d.%d, so no session is moved. Update it with `claude update`", v[0], v[1], v[2])
	}
	seats := map[string][]string{} // accounts by plan
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if prof.API != nil {
			continue
		}
		// Who the account is as the automatic choice reads it.
		if id, _ := profile.LoggedInAs(p, prof); id != nil {
			seats[id.OrgType] = append(seats[id.OrgType], prof.Name)
		}
		var plugin ext.Item
		for _, it := range ext.ScanAccount(prof.Name, profile.Dir(p, prof), p.GroupsRoot(), prof.Share, groupNames(cfg), "").Plugins {
			if name, _, _ := strings.Cut(it.Name, "@"); name == "ccam" {
				plugin = it
			}
		}
		switch {
		case plugin.Off:
			r.add(Warn, prof.Name, "auto_switch is on, but %s has the ccam plugin turned off, so its sessions are not moved. In Claude Code as %s: /plugin enable %s", prof.Name, prof.Name, plugin.Name)
		case plugin.Name == "" || plugin.Problem != "":
			r.add(Warn, prof.Name, "auto_switch is on, but %s does not have the ccam plugin, so its sessions are not moved. In Claude Code as %s: /plugin install ccam --marketplace jhyunwoo/ccam", prof.Name, prof.Name)
		}
	}
	// Quoted as the site quotes them (checked 2026-10-24).
	for _, s := range []struct{ orgType, plan, says, link string }{
		{"claude_team", "Claude Team", `"Each seat on a Team plan is assigned to one named individual, and only that person may use it." and "Users can't route their work through another user's account to use that user's remaining usage."`,
			"https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129"},
		{"claude_enterprise", "Claude Enterprise", `"Each Enterprise seat is assigned to one named individual, and only that person may use it." and "Users can't use another user's account to get around their own spend limits or usage limits."`,
			"https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0"},
	} {
		if names := seats[s.orgType]; len(names) > 0 {
			r.add(Info, "", "%s seats (%s): auto_switch moves a session between seats of one organization. Anthropic's help center says %s Whether you use it there is your call: %s", s.plan, strings.Join(names, ", "), s.says, s.link)
		}
	}
}

// checkAuthEnv warns about credentials in the environment: every account
// inherits them, and Claude Code prefers them over the account's own login.
func checkAuthEnv(r *report) {
	r.area = "Environment"
	for _, v := range []string{"CLAUDE_CODE_OAUTH_TOKEN", "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN"} {
		if os.Getenv(v) != "" {
			r.add(Warn, "", "%s is set in this shell; every account launched from it may use that credential instead of its own login", v)
		}
	}
}

// checkUsageFormat queries one logged-in account so a changed response
// format of the undocumented usage endpoint is noticed.
func checkUsageFormat(r *report, p paths.Paths, cfg *config.Config, f *usage.Fetcher) {
	if f == nil {
		return
	}
	r.area = "Usage"
	// One time limit for all accounts: on a network that does not answer,
	// each lookup would wait out its own in turn. Claude Code's lookup
	// (usage_source = "claude") is not cut short, as it may be renewing
	// the login; it is asked of one account only.
	ctx := context.Background()
	if f.Source != "claude" {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, usage.LookupTimeout+time.Second)
		defer cancel()
	}
	var failed *usage.Report
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		a := usage.Account{Name: prof.Name, Slot: prof.Slot, Dir: profile.Dir(p, prof), SessionDir: launcher.SessionDir(p, prof.Name)}
		rep := f.Fetch(ctx, a, usage.Options{NoCache: true})
		switch {
		case rep.FormatChanged:
			r.add(Warn, prof.Name, "the usage response format changed. `ccam usage` falls back to the status line cache (%s)", rep.Err)
			return
		case rep.Status == usage.StatusOK:
			r.add(OK, prof.Name, "usage lookup OK")
			return
		}
		if rep.Status == usage.StatusError {
			if failed == nil {
				failed = &rep
			}
			if f.Source == "claude" || ctx.Err() != nil {
				break
			}
		}
	}
	if failed != nil {
		r.add(Info, failed.Name, "could not check the usage response format: %s", failed.Err)
	}
}

func checkProfile(r *report, p paths.Paths, prof *config.Profile, fix, dry bool, mode links, preview *share.Plan, pending map[string]string) {
	dir := profile.Dir(p, prof)
	r.area = "Folder"
	info, err := os.Stat(dir)
	if err != nil {
		r.add(Error, prof.Name, "profile folder is missing: %s", dir)
		return
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		f := r.add(Warn, prof.Name, "folder permissions are %o (should be 0700)", info.Mode().Perm())
		if fix && os.Chmod(dir, 0o700) == nil {
			f.Fixed = true
		}
	}
	if prof.Dir != "" && profile.OwnClaudeJSON(p, dir) {
		r.add(Info, prof.Name, "%s has a .claude.json of its own, so `ccam %s` and plain `claude` (which uses %s) keep separate login details and MCP servers; to keep one, replace it with a link to %s",
			dir, prof.Name, filepath.Join(p.Home, ".claude.json"), filepath.Join(p.Home, ".claude.json"))
	}
	checkCredentials(r, prof, dir, p.Home, fix)

	r.area = "Share links"
	problems := share.CheckTargets(dir, prof.Share, p.GroupsRoot())
	if len(problems) == 0 {
		r.add(OK, prof.Name, "share links OK")
	}
	for _, pr := range problems {
		switch {
		case pr.Waits(p.GroupsRoot()):
			r.add(Warn, prof.Name, "%s is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)", pr.Path)
		case pr.Safe() && fix:
			r.add(Warn, prof.Name, "%s", pr)
		case pr.Kind == share.MissingTarget:
			r.add(Warn, prof.Name, "%s (ccam doctor --fix makes it)", pr)
		case pr.Kind == share.Staged:
			r.add(Warn, prof.Name, "%s (ccam doctor --fix puts it in place)", pr)
		case pr.Kind == share.StagedAside:
			// Only the user can tell which copy to keep; --fix leaves it.
			r.add(Warn, prof.Name, "%s", pr)
		case pr.Safe():
			// Usual right after an upgrade that shares more paths.
			r.add(Warn, prof.Name, "%s (the next launch or ccam doctor --fix makes it)", pr)
		default:
			r.add(Error, prof.Name, "%s", pr)
		}
	}
	base := len(r.out) - len(problems) // the findings of problems
	switch {
	case len(problems) == 0 || !fix && !dry:
	case mode == linksLeave:
		// Another ccam's change of groups is between its steps: a repair
		// now would make one of them fail.
		r.add(Info, prof.Name, "its share links are not repaired while another ccam changes the groups; run `ccam doctor --fix` again when it is done")
	case mode == linksWait && dry:
		// What finishing the stopped change does for this account.
		if preview != nil {
			setAside(r, prof.Name, preview, false)
		}
	case launcher.Running(p, prof.Name, dir):
		if !fix {
			break
		}
		// A running session keeps writing into what a repair would move;
		// only make the links and group folders that are missing, which
		// loses nothing.
		rest, err := share.QuickFix(problems)
		if err != nil {
			r.add(Error, prof.Name, "could not repair share links: %v", err)
			return
		}
		for i, pr := range problems {
			r.out[base+i].Fixed = pr.Safe()
		}
		if len(rest) > 0 {
			r.add(Warn, prof.Name, "is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest")
		}
	case mode == linksWait && fix:
		// The stopped change of groups repairs these in its own order;
		// only what is missing is made now.
		rest, err := share.QuickFix(problems)
		if err != nil {
			r.add(Error, prof.Name, "could not repair share links: %v", err)
			return
		}
		for i, pr := range problems {
			r.out[base+i].Fixed = pr.Safe()
		}
		if len(rest) > 0 {
			r.add(Warn, prof.Name, "the rest of its share links are repaired when the last group change is finished (ccam doctor --fix)")
		}
	case fix:
		// Build with the target as both sides: it re-derives every step from
		// what is actually on disk, keeps group content on conflicts and
		// backs up anything it replaces. It leaves a link that is right as
		// it is, so a group folder that is gone is made first.
		_, err := share.QuickFix(problems)
		var pl *share.Plan
		if err == nil {
			pl, err = share.Build(dir, p.GroupsRoot(), prof.Share, prof.Share, share.Options{})
		}
		if err == nil {
			err = pl.Apply()
		}
		if err != nil {
			r.add(Error, prof.Name, "could not repair share links: %v", err)
			return
		}
		for i, pr := range problems {
			r.out[base+i].Fixed = pr.Kind != share.StagedAside
		}
		setAside(r, prof.Name, pl, true)
	default:
		// A dry run: the plan --fix would make, after the accounts before.
		pl, err := share.Build(dir, p.GroupsRoot(), prof.Share, prof.Share, share.Options{Pending: pending})
		if err != nil {
			r.add(Error, prof.Name, "could not plan the repair of share links: %v", err)
			break
		}
		addCreated(pending, pl)
		setAside(r, prof.Name, pl, false)
	}
	checkUnknownEntries(r, prof, dir)
}

// setAside says what a share link repair of account name, plan pl, does
// (done: has done) with the account's own copies of paths the group has
// its own copy of: the account sees the group's, and its own goes to the
// backup folder.
func setAside(r *report, name string, pl *share.Plan, done bool) {
	for _, c := range pl.Conflicts {
		saved := ""
		for _, s := range pl.Steps {
			if s.Op == share.OpBackup && s.Entry.Name == c.Entry.Name {
				saved = s.To
			}
		}
		switch {
		case done && c.Choice == share.UseMine:
			r.add(Info, name, "%s differed from the group's copy; the group now uses this account's, as chosen for the group change (the group's copy is in %s, this account's file in %s)", c.Entry.Name, c.GroupSave, saved)
		case done && c.Choice == share.MergeJSON:
			r.add(Info, name, "%s differed from the group's copy; this account's settings were merged into it, as chosen for the group change (the group's copy is in %s, this account's file in %s)", c.Entry.Name, c.GroupSave, saved)
		case done && c.Choice == share.MergeLines:
			r.add(Info, name, "%s differed from the group's copy; this account's lines were added to it, and its file is in %s", c.Entry.Name, saved)
		case done:
			r.add(Info, name, "%s differed from the group's copy, which it now uses; this account's copy is in %s", c.Entry.Name, saved)
		case c.Choice == share.UseMine:
			r.add(Info, name, "%s differs from the group's copy: --fix would give the group this account's, as chosen for the group change (the group's copy goes to %s, this account's file to %s)", c.Entry.Name, c.GroupSave, saved)
		case c.Choice == share.MergeJSON:
			r.add(Info, name, "%s differs from the group's copy: --fix would merge this account's settings into it, as chosen for the group change (the group's copy goes to %s, this account's file to %s)", c.Entry.Name, c.GroupSave, saved)
		case c.Choice == share.MergeLines:
			r.add(Info, name, "%s differs from the group's copy: --fix would add this account's lines to it, and move its file to %s", c.Entry.Name, saved)
		default:
			r.add(Info, name, "%s differs from the group's copy: --fix would use the group's, and move this account's copy to %s", c.Entry.Name, saved)
		}
	}
	for _, n := range pl.Notes {
		r.add(Info, name, "%s", n)
	}
}

func checkCredentials(r *report, prof *config.Profile, dir, home string, fix bool) {
	r.area = "Login"
	cred := creds.FilePath(dir)
	if info, err := os.Lstat(cred); err == nil {
		if _, link := fsutil.LinkTarget(cred); link {
			r.add(Error, prof.Name, ".credentials.json is a symlink. Each account must have its own token")
			return
		}
		// Windows has no permission bits; the profile folder's ACL protects it.
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
			f := r.add(Warn, prof.Name, ".credentials.json permissions are %o (should be 0600)", info.Mode().Perm())
			if fix && os.Chmod(cred, 0o600) == nil {
				f.Fixed = true
			}
		}
	}
	if api := prof.API; api != nil {
		checkAPIKey(r, prof.Name, api)
		return
	}
	data, _, err := creds.Read(dir, home)
	if errors.Is(err, creds.ErrNone) {
		r.add(Warn, prof.Name, "not logged in → `ccam login %s`", prof.Name)
		return
	}
	if errors.Is(err, creds.ErrCredentialManager) {
		r.add(Info, prof.Name, "the login is kept in Windows Credential Manager, which ccam cannot check")
		return
	}
	if err != nil {
		r.add(Error, prof.Name, "could not read login: %v", err)
		return
	}
	var c struct {
		OAuth *struct {
			ExpiresAt int64 `json:"expiresAt"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(data, &c); err != nil {
		r.add(Info, prof.Name, "no subscription login token (may be an API key account)")
		return
	}
	if c.OAuth == nil {
		// What a logout leaves, or MCP sign-ins only; an API key login keeps
		// its key in .claude.json.
		var key string
		if ok, _ := profile.DecodeFileKey(filepath.Join(dir, ".claude.json"), "primaryApiKey", &key); ok {
			r.add(Info, prof.Name, "no subscription login token (an API key login)")
		} else {
			r.add(Warn, prof.Name, "not logged in → `ccam login %s`", prof.Name)
		}
		return
	}
	now := time.Now()
	switch {
	case creds.Judge(data, now) != "":
		// The refresh token expired, or Claude Code emptied it when the
		// server refused it: Claude Code cannot renew this login.
		r.add(Error, prof.Name, "login expired → `ccam login %s`", prof.Name)
	case c.OAuth.ExpiresAt < now.UnixMilli():
		r.add(Info, prof.Name, "access token expired. Claude Code refreshes it on the next launch")
	default:
		who := orDash(prof.Email)
		if id, err := profile.ReadIdentity(dir); err == nil && id != nil {
			if id.Email != "" {
				who = id.Email // whom the folder is logged in as now
			}
			if plan := id.PlanName(); plan != "" {
				who += ", " + plan
				if id.Pool() != config.Personal && id.OrgName != "" {
					who += " · " + id.OrgName
				}
			}
		}
		r.add(OK, prof.Name, "login valid (%s)", who)
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// checkUnknownEntries flags files Claude Code started creating that ccam
// does not know about, a hint that the layout changed.
func checkUnknownEntries(r *report, prof *config.Profile, dir string) {
	r.area = "Folder"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var unknown []string
	for _, e := range entries {
		if !share.Known(e.Name()) {
			unknown = append(unknown, e.Name())
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		r.add(Info, prof.Name, "unknown entries (ccam knows the folder layout of Claude Code %s; it may have changed): %s", LayoutVersion, strings.Join(unknown, ", "))
	}
}

// checkPlugins finds plugin paths that point through one profile's folder
// into a shared group. They break when that profile is removed or renamed,
// so --fix rewrites them to the group folder itself.
func checkPlugins(r *report, p paths.Paths, cfg *config.Config, fix bool) {
	r.area = "Plugins"
	groups, _ := os.ReadDir(p.GroupsRoot())
	sep := string(filepath.Separator)
	for _, g := range groups {
		pluginsDir := filepath.Join(p.GroupDir(g.Name()), "plugins")
		n := 0
		for i := range cfg.Profiles {
			pre := filepath.Join(profile.Dir(p, &cfg.Profiles[i]), "plugins") + sep
			n += profile.PluginPathsUnder(pluginsDir, pre)
		}
		if n == 0 {
			continue
		}
		f := r.add(Warn, "", "%d plugin paths in group %s go through one account's folder", n, g.Name())
		if !fix {
			continue
		}
		f.Fixed = true
		for i := range cfg.Profiles {
			pre := filepath.Join(profile.Dir(p, &cfg.Profiles[i]), "plugins") + sep
			if _, err := profile.RewritePluginPaths(pluginsDir, pre, pluginsDir+sep); err != nil {
				f.Fixed = false
			}
		}
	}
}

func checkMCP(r *report, p paths.Paths, cfg *config.Config) {
	r.area = "MCP"
	for _, g := range cfg.Groups {
		if !g.Has("mcp") {
			continue
		}
		store, err := mcp.Load(mcp.Path(p.GroupDir(g.Folder())))
		if err != nil {
			r.add(Error, "", "%v", err)
			continue
		}
		for _, n := range store.Names() {
			if err := mcp.Validate(store.Servers[n]); err != nil {
				r.add(Warn, "", "group %s MCP server %s: %v", g.Name, n, err)
			}
		}
	}
	checkHomeMCP(r, p, cfg)
	checkAccountMCP(r, p, cfg)
}

// checkHomeMCP reports the user-scope servers of plain Claude Code's
// ~/.claude.json that no account loads: once Claude Code runs as a ccam
// account it reads that account's .claude.json, so they seem gone.
func checkHomeMCP(r *report, p paths.Paths, cfg *config.Config) {
	var accts, groups []string
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		accts = append(accts, profile.Dir(p, prof))
		if f := prof.Share.Get("mcp"); f != config.Local {
			groups = append(groups, p.GroupDir(f))
		}
	}
	from, names := migrate.HomeOnly(filepath.Join(p.Home, ".claude"), p.Home, accts, groups)
	if len(names) == 0 {
		return
	}
	r.add(Warn, "", "%s (plain claude's) has MCP servers no account loads: %s → `%s` copies them into the shared group", from, strings.Join(names, ", "), migrate.MCPCommand(cfg))
}

// checkAccountMCP reports, per account, the servers no other account of
// its group loads (in its .claude.json, or for one folder only), and the
// servers other accounts are signed in to and it is not: a sign-in holds
// for one account only, and ccam never copies one.
func checkAccountMCP(r *report, p paths.Paths, cfg *config.Config) {
	var sets []ext.Set
	dirs := map[string]string{}
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		dir := profile.Dir(p, prof)
		s := ext.ScanAccount(prof.Name, dir, p.GroupsRoot(), prof.Share, groupNames(cfg), "")
		sets = append(sets, s)
		dirs[prof.Name] = dir
		g := cfg.GroupFor(prof.Name, "mcp")
		if g == nil || len(slices.DeleteFunc(slices.Clone(g.Members), func(m string) bool { return !cfg.IsAccount(m) })) < 2 {
			continue
		}
		var own, names []string
		for _, it := range s.MCP {
			if it.Source == ext.Account && !ext.PluginServer(it.Name) {
				own = append(own, it.Name)
				names = append(names, it.Name)
			}
		}
		var local []string
		folders, _ := mcp.FolderServers(dir)
		for _, folder := range slices.Sorted(maps.Keys(folders)) {
			for _, n := range slices.Sorted(maps.Keys(folders[folder])) {
				local = append(local, n+" ("+folder+")")
				if !slices.Contains(names, n) {
					names = append(names, n)
				}
			}
		}
		if len(names) == 0 {
			continue
		}
		var parts []string
		if len(own) > 0 {
			parts = append(parts, strings.Join(own, ", "))
		}
		if len(local) > 0 {
			parts = append(parts, "for one folder only: "+strings.Join(local, ", "))
		}
		r.add(Info, prof.Name, "MCP servers the other accounts of group %s do not load: %s → `ccam mcp import %s %s` shares them", g.Name, strings.Join(parts, "; "), prof.Name, strings.Join(names, " "))
	}
	signedIn := map[string]map[string]bool{}
	for _, a := range ext.SignInAccounts(sets) {
		signedIn[a] = mcp.SignedIn(dirs[a], p.Home)
	}
	missing := map[string][]string{}
	for _, g := range ext.SignInGaps(sets, signedIn) {
		for _, a := range g.Missing {
			missing[a] = append(missing[a], g.Server+" ("+strings.Join(g.SignedIn, ", ")+")")
		}
	}
	for _, a := range slices.Sorted(maps.Keys(missing)) {
		r.add(Warn, a, "not signed in to MCP servers other accounts are signed in to: %s → sign in with /mcp in a session of %s (a sign-in holds for one account; ccam does not copy them)", strings.Join(missing[a], ", "), a)
	}
}

// groupNames maps group folders to group names where they differ.
func groupNames(cfg *config.Config) map[string]string {
	out := map[string]string{}
	for _, g := range cfg.Groups {
		if g.Dir != "" {
			out[g.Dir] = g.Name
		}
	}
	return out
}

func checkShims(r *report, p paths.Paths, cfg *config.Config, opt Options) {
	if opt.CcamPath == "" {
		return
	}
	r.area = "Shortcuts"
	var names, missing, older, others []string
	var runs []string // what each of others runs
	for _, prof := range cfg.Profiles {
		if prof.Shortcut == "" {
			continue
		}
		names = append(names, prof.Shortcut)
		switch bin, other := shim.Other(p, opt.CcamPath, prof.Shortcut); {
		case shim.Linked(p.Bin, opt.CcamPath, prof.Shortcut):
		case shim.Stale(p, prof.Shortcut):
			older = append(older, prof.Shortcut)
		case other:
			others = append(others, prof.Shortcut)
			runs = append(runs, fmt.Sprintf("%s → %s", prof.Shortcut, bin))
		default:
			missing = append(missing, prof.Shortcut)
		}
	}
	if len(names) == 0 {
		return // shortcuts are opt-in
	}
	if len(missing)+len(older)+len(others) == 0 {
		r.add(OK, "", "shortcuts OK (%s)", strings.Join(names, ", "))
	}
	if len(missing)+len(older) > 0 {
		var parts []string
		if len(missing) > 0 {
			parts = append(parts, fmt.Sprintf("missing shortcuts in %s: %s", p.Bin, strings.Join(missing, ", ")))
		}
		if len(older) > 0 {
			parts = append(parts, fmt.Sprintf("shortcuts in %s that run another ccam (one replaced or moved): %s", p.Bin, strings.Join(older, ", ")))
		}
		f := r.add(Warn, "", "%s", strings.Join(parts, "; "))
		if opt.fixes("") {
			// Also renews hard links (Windows) of the binary an update or a
			// new install replaced.
			created, err := shim.Install(p, opt.CcamPath, cfg)
			f.Fixed = err == nil && !slices.ContainsFunc(slices.Concat(missing, older), func(n string) bool { return !slices.Contains(created, n) })
			if err != nil {
				f.Msg += " (" + strings.ReplaceAll(err.Error(), "\n", "; ") + ")"
			}
		}
	}
	if len(others) > 0 {
		// They work; another ccam made them, and add or import leave them
		// to it. An explicit --fix moves them to this one.
		f := r.add(Info, "", "shortcuts in %s that run another ccam, which is still installed: %s; `ccam doctor --fix` points them at this one (%s)", p.Bin, strings.Join(runs, ", "), opt.CcamPath)
		if opt.fixes("") {
			done, err := shim.Repoint(p, opt.CcamPath, others)
			f.Fixed = err == nil && len(done) == len(others)
			if err != nil {
				f.Msg += " (" + strings.ReplaceAll(err.Error(), "\n", "; ") + ")"
			}
		}
	}
	if !shim.InPath(p.Bin) {
		r.add(Warn, "", "%s is not in PATH", p.Bin)
	}
}
