// Package doctor checks a ccam installation and repairs what it safely can.
package doctor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
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
const LayoutVersion = "2.1.302"

type report struct {
	out  []Finding
	area string // the Area of the findings added next
}

func (r *report) add(level Level, prof, format string, a ...any) *Finding {
	r.out = append(r.out, Finding{Level: level, Profile: prof, Area: r.area, Msg: fmt.Sprintf(format, a...)})
	return &r.out[len(r.out)-1]
}

// Run checks everything and, with opt.Fix, repairs safe problems.
func Run(p paths.Paths, cfg *config.Config, opt Options) []Finding {
	r := &report{}
	checkClaude(r, p, cfg)
	checkAuthEnv(r)
	for i := range cfg.Profiles {
		checkProfile(r, p, &cfg.Profiles[i], opt.fixes(cfg.Profiles[i].Name))
	}
	checkPlugins(r, p, cfg, opt.fixes(""))
	checkMCP(r, p, cfg)
	checkUsageFormat(r, p, cfg, opt.Usage)
	checkShims(r, p, cfg, opt)
	checkFolders(r, cfg)
	checkGroups(r, cfg)
	checkIdentities(r, p, cfg)
	return r.out
}

// checkIdentities reports a folder now logged in to another account than
// ccam recorded (a /login inside a session), and folders logged in to one
// Claude account, which share its limits.
func checkIdentities(r *report, p paths.Paths, cfg *config.Config) {
	r.area = "Login"
	seen := map[string]string{} // account and organization → first profile
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if prof.API != nil {
			continue
		}
		id, err := profile.ReadIdentity(profile.Dir(p, prof))
		if err != nil || id == nil {
			continue
		}
		if was := prof.Login; was != nil && was.AccountUUID != "" && was.AccountUUID != id.AccountUUID {
			r.add(Warn, prof.Name, "now logged in as %s, but ccam recorded %s; log in again with `ccam login %s` if that is not what you meant", orDash(id.Email), orDash(was.Email), prof.Name)
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

func checkClaude(r *report, p paths.Paths, cfg *config.Config) {
	r.area = "Claude Code"
	bin, err := launcher.ResolveClaude(p, cfg.ClaudeBin)
	if err != nil {
		r.add(Error, "", "%v", err)
		return
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		r.add(Error, "", "%s --version failed: %v", bin, err)
		return
	}
	version := strings.TrimSpace(string(out))
	r.add(OK, "", "Claude Code %s (%s)", version, bin)
	if !strings.HasPrefix(version, LayoutVersion) {
		r.add(Info, "", "ccam was verified against the folder layout of Claude Code %s. Any 'unknown entries' below may mean the layout changed", LayoutVersion)
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
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		a := usage.Account{Name: prof.Name, Slot: prof.Slot, Dir: profile.Dir(p, prof), SessionDir: launcher.SessionDir(p, prof.Name)}
		rep := f.Fetch(context.Background(), a, usage.Options{NoCache: true})
		switch {
		case rep.FormatChanged:
			r.add(Warn, prof.Name, "the usage response format changed. `ccam usage` falls back to the status line cache (%s)", rep.Err)
			return
		case rep.Status == usage.StatusOK:
			r.add(OK, prof.Name, "usage lookup OK")
			return
		}
	}
}

func checkProfile(r *report, p paths.Paths, prof *config.Profile, fix bool) {
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
	checkCredentials(r, prof, dir, p.Home, fix)

	r.area = "Share links"
	problems := share.Check(dir, prof.Share, p.GroupsRoot())
	if len(problems) == 0 {
		r.add(OK, prof.Name, "share links OK")
	}
	for _, pr := range problems {
		switch {
		case pr.Waits(p.GroupsRoot()):
			r.add(Warn, prof.Name, "%s is not shared with the group yet (ccam doctor --fix moves it in)", pr.Path)
		case pr.Kind == share.MissingLink && fix:
			r.add(Warn, prof.Name, "%s", pr)
		case pr.Kind == share.MissingLink:
			// Usual right after an upgrade that shares more paths.
			r.add(Warn, prof.Name, "%s (the next launch or ccam doctor --fix makes it)", pr)
		default:
			r.add(Error, prof.Name, "%s", pr)
		}
	}
	if fix && len(problems) > 0 && launcher.Running(p, prof.Name) {
		// A running session keeps writing into what a repair would move;
		// only make the links that are missing, as a launch does.
		base := len(r.out) - len(problems)
		if _, err := share.QuickFix(problems); err != nil {
			r.add(Error, prof.Name, "could not repair share links: %v", err)
			return
		}
		for i, pr := range problems {
			r.out[base+i].Fixed = pr.Kind == share.MissingLink
		}
		r.add(Warn, prof.Name, "is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest")
	} else if fix && len(problems) > 0 {
		// Build with the target as both sides: it re-derives every step from
		// what is actually on disk, keeps group content on conflicts and
		// backs up anything it replaces.
		pl, err := share.Build(dir, p.GroupsRoot(), prof.Share, prof.Share, share.Options{})
		if err == nil {
			err = pl.Apply()
		}
		if err != nil {
			r.add(Error, prof.Name, "could not repair share links: %v", err)
			return
		}
		for i := len(r.out) - len(problems); i < len(r.out); i++ {
			r.out[i].Fixed = true
		}
	}
	checkUnknownEntries(r, prof, dir)
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
	if err != nil {
		r.add(Error, prof.Name, "could not read login: %v", err)
		return
	}
	var c struct {
		OAuth *struct {
			ExpiresAt             int64 `json:"expiresAt"`
			RefreshTokenExpiresAt int64 `json:"refreshTokenExpiresAt"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(data, &c); err != nil || c.OAuth == nil {
		r.add(Info, prof.Name, "no subscription login token (may be an API key account)")
		return
	}
	now := time.Now().UnixMilli()
	switch {
	case c.OAuth.RefreshTokenExpiresAt > 0 && c.OAuth.RefreshTokenExpiresAt < now:
		r.add(Error, prof.Name, "login expired → `ccam login %s`", prof.Name)
	case c.OAuth.ExpiresAt < now:
		r.add(Info, prof.Name, "access token expired. Claude Code refreshes it on the next launch")
	default:
		who := orDash(prof.Email)
		if id, err := profile.ReadIdentity(dir); err == nil && id != nil {
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
		r.add(Info, prof.Name, "unknown entries (Claude Code layout may have changed): %s", strings.Join(unknown, ", "))
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
			n += countPrefix(pluginsDir, pre)
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

func countPrefix(pluginsDir, prefix string) int {
	n := 0
	for _, name := range profile.PluginFiles {
		data, err := os.ReadFile(filepath.Join(pluginsDir, name))
		if err != nil {
			continue
		}
		var doc any
		if json.Unmarshal(data, &doc) != nil {
			continue
		}
		profile.RewriteStrings(doc, func(s string) string {
			if strings.HasPrefix(s, prefix) {
				n++
			}
			return s
		})
	}
	return n
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
}

func checkShims(r *report, p paths.Paths, cfg *config.Config, opt Options) {
	if opt.CcamPath == "" {
		return
	}
	r.area = "Shortcuts"
	var names, missing []string
	for _, prof := range cfg.Profiles {
		if prof.Shortcut == "" {
			continue
		}
		names = append(names, prof.Shortcut)
		if !shim.Linked(p.Bin, opt.CcamPath, prof.Shortcut) {
			missing = append(missing, prof.Shortcut)
		}
	}
	if len(names) == 0 {
		return // shortcuts are opt-in
	}
	if len(missing) == 0 {
		r.add(OK, "", "shortcuts OK (%s)", strings.Join(names, ", "))
	} else {
		f := r.add(Warn, "", "missing shortcuts in %s: %s", p.Bin, strings.Join(missing, ", "))
		if opt.fixes("") {
			_, err := shim.Install(p.Bin, opt.CcamPath, cfg)
			f.Fixed = err == nil
		}
	}
	if !shim.InPath(p.Bin) {
		r.add(Warn, "", "%s is not in PATH", p.Bin)
	}
}
