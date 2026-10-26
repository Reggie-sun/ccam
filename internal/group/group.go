// Package group carries a change of groups over to the account folders:
// one share plan for each account whose shared items change, built so that
// a later plan sees what earlier ones will put into the group folders.
package group

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

// Options tune Build.
type Options struct {
	Seed  string // the account whose copy goes into a group first (--from)
	Fresh bool   // accounts that stop sharing an item start without it (not one registered in place)
}

// Move is one item that changes for an account; From and To are group
// names, "" for not shared.
type Move struct{ Item, From, To string }

// Account is what the change does to one account.
type Account struct {
	Name     string
	Dir      string
	Imported bool // registered in place (ccam import), so plain claude may use it
	Moves    []Move
	Plan     *share.Plan
	// DropMCP lists the servers ccam copied into .claude.json (Sync)
	// unchanged since, removed when the account moves MCP to another group
	// or stops sharing it with --fresh; ChangedMCP those changed since,
	// which are the account's now and stay.
	DropMCP, ChangedMCP []string
	unsync              bool
	// SeedMCP copies the mcp.json of the group the account leaves into the
	// new group it gets MCP from, when that one has none yet: the account
	// keeps its servers, as it keeps the files of the other items.
	SeedMCP *[2]string
	// KeepMCP is the folder of the group whose servers an account that
	// stops sharing MCP keeps as its own (KeepNames, copied into its
	// .claude.json); OwnMCP are its own servers of the same names, which
	// stay as they are. With --fresh, StopMCP lists the servers it stops
	// loading instead.
	KeepMCP           string
	KeepNames, OwnMCP []string
	StopMCP           []string
	// DropKept lists the copies kept when the account stopped sharing MCP
	// that the group it joins has servers of: they go, since the group
	// gives them again and a copy would hide its later changes.
	DropKept []string
	dropKept bool
	// OnlyMCP lists the account's own servers the group it joins lacks:
	// they stay its own unless `ccam mcp import` shares them.
	OnlyMCP []string
	// mcp is the group folder the account gets MCP from after the change
	// ("" for none).
	mcp string
}

// Change is the work to get from one config to another.
type Change struct {
	Accounts []*Account
	paths    paths.Paths
	// Every account's shares in the config the change starts from and in
	// the one it leads to.
	from, to map[string]config.Share
	fresh    bool   // Options.Fresh, which the plans of a resumed change need too
	note     []byte // what Begin noted, or the note Resume finishes
	// held is lockPath, taken by Begin and kept until Apply or Abort: no
	// doctor --fix finishes the change between the save and Apply.
	held *lock.Flock
	// groupPlugins are the plugins folders of the groups an account takes
	// its extensions out of. Claude Code records plugin paths through the
	// link of the account that installed them, so before any plan runs the
	// paths there through an account's plugins folder (accountPlugins)
	// point at the group folder itself: the accounts that stay do not load
	// those plugins from the leaving account's folder, and a copy taken
	// from the group then points at itself.
	groupPlugins, accountPlugins []string
}

// Build compares old and next (an edited copy of old) and plans the folder
// changes. The first content of a group comes from the group folder, then
// opt.Seed, then the default account, then the others by number.
func Build(p paths.Paths, old, next *config.Config, opt Options) (*Change, error) {
	order := make([]*config.Profile, 0, len(next.Profiles))
	for i := range next.Profiles {
		order = append(order, &next.Profiles[i])
	}
	rank := func(pr *config.Profile) int {
		switch pr.Name {
		case opt.Seed:
			return 0
		case next.Default:
			return 1
		}
		return 2
	}
	slices.SortStableFunc(order, func(a, b *config.Profile) int {
		return cmp.Or(cmp.Compare(rank(a), rank(b)), cmp.Compare(a.Slot, b.Slot))
	})
	ch := &Change{paths: p, from: shares(old), to: shares(next), fresh: opt.Fresh, accountPlugins: accountPlugins(p, old, next)}
	pending := map[string]string{}
	// A group without an mcp.json starts with the one of the group the
	// first account (in this order) moving to it comes from. Known before
	// any account is planned, so that one joining it from no group sees
	// those servers too.
	seeds := map[string]string{} // group mcp.json → the one it starts as
	for _, prof := range order {
		fm := config.Local
		if o, ok := old.Find(prof.Name); ok {
			fm = o.Share.Get("mcp")
		}
		tm := prof.Share.Get("mcp")
		if fm == tm || fm == config.Local || tm == config.Local {
			continue
		}
		src, dst := mcp.Path(p.GroupDir(fm)), mcp.Path(p.GroupDir(tm))
		if _, ok := seeds[dst]; !ok && fsutil.Exists(src) && !fsutil.Exists(dst) {
			seeds[dst] = src
		}
	}
	seeded := map[string]bool{}
	for _, prof := range order {
		var from config.Share
		if o, ok := old.Find(prof.Name); ok {
			from = o.Share
		}
		if from.Equal(prof.Share) {
			continue
		}
		if g := from.Get("extensions"); g != config.Local && g != prof.Share.Get("extensions") {
			if dir := filepath.Join(p.GroupDir(g), "plugins"); !slices.Contains(ch.groupPlugins, dir) {
				ch.groupPlugins = append(ch.groupPlugins, dir)
			}
		}
		dir := profile.Dir(p, prof)
		// A folder registered in place (ccam import) held what it shared as
		// its own files before, and plain claude goes on using it: what it
		// stops sharing it keeps a copy of, --fresh or not.
		aopt := opt
		aopt.Fresh = opt.Fresh && prof.Dir == ""
		pl, err := share.Build(dir, p.GroupsRoot(), from, prof.Share, share.Options{Fresh: aopt.Fresh, Pending: pending})
		if err != nil {
			return nil, fmt.Errorf("account %s: %w", prof.Name, err)
		}
		for k, v := range pl.Created() {
			// Content an account brings beats an empty folder another makes.
			if old, ok := pending[k]; !ok || old == "" {
				pending[k] = v
			}
		}
		a := &Account{Name: prof.Name, Dir: dir, Imported: prof.Dir != "", Plan: pl}
		for _, item := range config.Items {
			if from.Get(item) == prof.Share.Get(item) {
				continue
			}
			a.Moves = append(a.Moves, Move{Item: item, From: groupName(old, prof.Name, item), To: groupName(next, prof.Name, item)})
		}
		if err := planMCP(p, a, from.Get("mcp"), prof.Share.Get("mcp"), aopt, seeds, seeded); err != nil {
			return nil, fmt.Errorf("account %s: %w", prof.Name, err)
		}
		ch.Accounts = append(ch.Accounts, a)
	}
	return ch, nil
}

// planMCP plans what happens to the MCP servers of account a, whose MCP
// group changes from the folder fm to tm (config.Local: none).
func planMCP(p paths.Paths, a *Account, fm, tm string, opt Options, seeds map[string]string, seeded map[string]bool) error {
	if tm != config.Local {
		a.mcp = tm
	}
	if fm == tm {
		return nil
	}
	state := launcher.MCPStateDir(p, a.Name)
	switch {
	case fm == config.Local:
		// Copies kept when it stopped sharing MCP before would hide the
		// group's servers of those names.
		dst := mcp.Path(p.GroupDir(tm))
		dst = cmp.Or(seeds[dst], dst)
		var err error
		if a.DropKept, err = mcp.KeptPlan(a.Dir, dst, state); err != nil {
			return err
		}
		a.dropKept = len(a.DropKept) > 0
		if group, err := mcp.Load(dst); err == nil {
			own, _ := mcp.LocalServers(a.Dir)
			for _, n := range slices.Sorted(maps.Keys(own)) {
				if _, shared := group.Servers[n]; !shared {
					a.OnlyMCP = append(a.OnlyMCP, n)
				}
			}
		}
	case tm == config.Local && !opt.Fresh:
		// It keeps the group's servers, as it keeps the files of the other items.
		a.KeepMCP = p.GroupDir(fm)
		var err error
		if a.KeepNames, a.OwnMCP, err = mcp.KeepPlan(a.Dir, a.KeepMCP, state); err != nil {
			return err
		}
	default:
		// The old group's servers copied into .claude.json would load on,
		// and under the flag strategy win over the new group's.
		a.unsync = true
		a.DropMCP, a.ChangedMCP = mcp.Copied(a.Dir, state)
		src, dst := mcp.Path(p.GroupDir(fm)), mcp.Path(p.GroupDir(tm))
		switch {
		case tm == config.Local:
			if group, err := mcp.Load(src); err == nil {
				own, _ := mcp.LocalServers(a.Dir)
				for _, n := range group.Names() {
					if _, ok := own[n]; !ok || slices.Contains(a.DropMCP, n) {
						a.StopMCP = append(a.StopMCP, n)
					}
				}
			}
		case seeds[dst] == src && !seeded[dst]:
			a.SeedMCP = &[2]string{src, dst}
			seeded[dst] = true
		}
	}
	// A .claude.json ccam cannot read would stop the change at this step,
	// after the config is saved; KeepPlan checks it for Keep.
	if a.dropKept || a.unsync && mcp.HasCopy(state) {
		return mcp.Readable(a.Dir)
	}
	return nil
}

// accountPlugins lists the plugins folder of every account in cfgs.
func accountPlugins(p paths.Paths, cfgs ...*config.Config) []string {
	var out []string
	for _, c := range cfgs {
		for i := range c.Profiles {
			if dir := filepath.Join(profile.Dir(p, &c.Profiles[i]), "plugins"); !slices.Contains(out, dir) {
				out = append(out, dir)
			}
		}
	}
	return out
}

// freePlugins points the plugin paths in the group plugins folders dirs
// that run through one of the accounts' plugins folders at the group
// folder itself.
func freePlugins(dirs, accounts []string) error {
	sep := string(filepath.Separator)
	for _, dir := range dirs {
		for _, a := range accounts {
			if profile.PluginPathsUnder(dir, a+sep) == 0 {
				continue
			}
			if _, err := profile.RewritePluginPaths(dir, a+sep, dir+sep); err != nil {
				return fmt.Errorf("could not point the plugin paths in %s at the group folder: %w", dir, err)
			}
		}
	}
	return nil
}

func groupName(c *config.Config, account, item string) string {
	if g := c.GroupFor(account, item); g != nil {
		return g.Name
	}
	return ""
}

// Empty reports whether the change does nothing.
func (ch *Change) Empty() bool { return len(ch.Accounts) == 0 }

// Names lists the accounts the change touches.
func (ch *Change) Names() []string {
	out := make([]string, len(ch.Accounts))
	for i, a := range ch.Accounts {
		out[i] = a.Name
	}
	return out
}

// Check refuses a change that cannot be applied now: an account in it is
// running (it keeps writing into the folders being moved), or Windows
// cannot make the file links it needs.
func (ch *Change) Check(running func(name, dir string) bool) error {
	var busy []string
	for _, a := range ch.Accounts {
		if running(a.Name, a.Dir) {
			busy = append(busy, a.Name)
		}
	}
	if len(busy) > 0 {
		verb := "is"
		if len(busy) > 1 {
			verb = "are"
		}
		return fmt.Errorf("%s %s running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)", strings.Join(busy, ", "), verb)
	}
	for _, a := range ch.Accounts {
		for _, s := range a.Plan.Steps {
			if s.Op != share.OpLink || s.Entry.Dir {
				continue
			}
			root := ch.paths.GroupsRoot()
			if err := os.MkdirAll(root, 0o700); err != nil {
				return err
			}
			if err := fsutil.CanLinkFiles(root); err != nil {
				return fmt.Errorf("%v, and %s needs a link for %s. Turn it on, or leave settings, memory and history out of the group", err, a.Name, s.Entry.Name)
			}
			return nil
		}
	}
	return nil
}

// Running is the usual check for Check.
func Running(p paths.Paths) func(name, dir string) bool {
	return func(name, dir string) bool { return launcher.Running(p, name, dir) }
}

// Recheck makes edit again on cur, the config as it is right before the
// change is saved, after questions that may have waited, and checks the
// change again. Another ccam may have changed the groups meanwhile: then
// cur, or cur after edit, no longer has the shares the plan was made for,
// and the plan would not match what is on disk. An account may have
// started meanwhile.
func (ch *Change) Recheck(cur *config.Config, edit func(*config.Config) error, running func(name, dir string) bool) error {
	changed := errors.New("the groups changed while this plan waited; nothing was moved, make the change again")
	if !maps.EqualFunc(shares(cur), ch.from, config.Share.Equal) {
		return changed
	}
	if err := edit(cur); err != nil {
		return err
	}
	if !maps.EqualFunc(shares(cur), ch.to, config.Share.Equal) {
		return changed
	}
	return ch.Check(running)
}

// shares maps every account of c to what it shares.
func shares(c *config.Config) map[string]config.Share {
	out := make(map[string]config.Share, len(c.Profiles))
	for _, p := range c.Profiles {
		out[p.Name] = p.Share
	}
	return out
}

// Conflict is a file conflict in one account's plan.
type Conflict struct {
	Account string
	*share.Conflict
}

func (ch *Change) Conflicts() []Conflict {
	var out []Conflict
	for _, a := range ch.Accounts {
		for _, c := range a.Plan.Conflicts {
			out = append(out, Conflict{Account: a.Name, Conflict: c})
		}
	}
	return out
}

// Describe writes the change, account by account, with paths relative to
// ccam's data folder.
func (ch *Change) Describe(w io.Writer) {
	if ch.Empty() {
		fmt.Fprintln(w, "  Nothing to change.")
		return
	}
	short := func(path string) string {
		if rel, err := filepath.Rel(ch.paths.Data, path); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
		return path
	}
	sep := string(filepath.Separator)
	header := false
	for _, dir := range ch.groupPlugins {
		for _, a := range ch.accountPlugins {
			if profile.PluginPathsUnder(dir, a+sep) == 0 {
				continue
			}
			if !header {
				fmt.Fprintln(w, "group plugin lists:")
				header = true
			}
			fmt.Fprintf(w, "  paths    plugin paths under %s → %s\n", short(a), short(dir))
		}
	}
	for _, a := range ch.Accounts {
		parts := make([]string, len(a.Moves))
		for i, m := range a.Moves {
			switch {
			case m.From == "":
				parts[i] = fmt.Sprintf("+%s (group %s)", m.Item, m.To)
			case m.To == "":
				parts[i] = fmt.Sprintf("-%s (group %s)", m.Item, m.From)
			default:
				parts[i] = fmt.Sprintf("%s: group %s → %s", m.Item, m.From, m.To)
			}
		}
		fmt.Fprintf(w, "%s: %s\n", a.Name, strings.Join(parts, ", "))
		if !a.Plan.Empty() {
			a.Plan.DescribeUnder(w, ch.paths.Data)
		}
		for _, m := range a.Moves {
			if m.Item == "mcp" && m.To != "" {
				fmt.Fprintf(w, "  MCP servers come from group %s at the next launch\n", m.To)
			}
		}
		if a.SeedMCP != nil {
			fmt.Fprintf(w, "  copy     %s → %s (the new group's first MCP servers)\n", short(a.SeedMCP[0]), short(a.SeedMCP[1]))
		}
		if len(a.DropMCP) > 0 {
			fmt.Fprintf(w, "  remove   MCP servers ccam copied into .claude.json: %s (a backup is kept)\n", strings.Join(a.DropMCP, ", "))
		}
		if len(a.ChangedMCP) > 0 {
			fmt.Fprintf(w, "  keep     MCP servers changed in .claude.json since ccam copied them: %s\n", strings.Join(a.ChangedMCP, ", "))
		}
		if len(a.KeepNames) > 0 {
			claudeJSON := filepath.Join(a.Dir, ".claude.json")
			backup := ""
			if fsutil.Exists(claudeJSON) {
				backup = " (a backup is kept)"
			}
			fmt.Fprintf(w, "  copy     MCP servers %s: %s → %s%s\n", strings.Join(a.KeepNames, ", "), short(mcp.Path(a.KeepMCP)), short(claudeJSON), backup)
		}
		if len(a.OwnMCP) > 0 {
			fmt.Fprintf(w, "  keep     %s's own MCP servers %s (the group's servers with those names differ)\n", a.Name, strings.Join(a.OwnMCP, ", "))
		}
		if len(a.StopMCP) > 0 {
			fmt.Fprintf(w, "  MCP servers %s stop loading for %s (--fresh)\n", strings.Join(a.StopMCP, ", "), a.Name)
		}
		if len(a.DropKept) > 0 {
			fmt.Fprintf(w, "  remove   MCP servers %s kept in .claude.json when %s left a group: the group gives them again (a backup is kept)\n", strings.Join(a.DropKept, ", "), a.Name)
		}
		if len(a.OnlyMCP) > 0 {
			fmt.Fprintf(w, "  Note: %s's own MCP servers %s stay its own; `ccam mcp import %s` shares them with the group\n", a.Name, strings.Join(a.OnlyMCP, ", "), a.Name)
		}
		if a.Imported {
			fmt.Fprintf(w, "  Note: %s is an imported folder, so plain `claude` sees this change too if it uses it\n", a.Dir)
			if ch.fresh && slices.ContainsFunc(a.Moves, func(m Move) bool { return m.To == "" }) {
				fmt.Fprintf(w, "  Note: %s keeps a copy of what it stops sharing: --fresh applies to ccam's own folders only\n", a.Name)
			}
		}
	}
}

// record is a change Apply has not finished, in the state folder while it
// runs, so that `ccam doctor --fix` (Resume) can finish one that stopped
// halfway: the accounts in the order of their plans, since the first
// content of a group follows it, and their MCP steps, which a plan toward
// the saved config no longer shows.
type record struct {
	Accounts []recorded `json:"accounts"`
	Fresh    bool       `json:"fresh,omitempty"` // --fresh: items an account stops sharing start empty
	// Plugins are the group plugins folders whose paths through account
	// folders the change points at the group folder (groupPlugins).
	Plugins []string `json:"plugins,omitempty"`
}

type recorded struct {
	Name string `json:"name"`
	// MCP is the group folder the account gets MCP from after the change
	// ("" for none): the MCP steps are done only while its config says so.
	MCP      string     `json:"mcp,omitempty"`
	Seed     *[2]string `json:"seed_mcp,omitempty"`
	Unsync   bool       `json:"unsync_mcp,omitempty"`
	Keep     string     `json:"keep_mcp,omitempty"`
	DropKept bool       `json:"drop_kept_mcp,omitempty"`
	// Choices are the answers to the account's conflicts, every one: the
	// first candidate is not always keeping the group's copy (history.jsonl
	// merges by line). Entry name -> Resolution.String().
	Choices map[string]string `json:"choices,omitempty"`
}

// NotePath is the note of a change of groups Apply has not finished.
func NotePath(p paths.Paths) string { return filepath.Join(p.State, "group-change.json") }

// Pending reports whether a change of groups stopped halfway, which
// `ccam doctor --fix` finishes, or is being carried out (Changing).
func Pending(p paths.Paths) bool { return fsutil.Exists(NotePath(p)) }

// ErrPending refuses a change of groups while an earlier one is unfinished:
// its note would be lost, and a plan toward the saved config does not show
// what is left of it.
var ErrPending = errors.New("the last group change stopped halfway; finish it with `ccam doctor --fix` first")

// lockPath is held from Begin to the end of Apply (or Abort), and by
// doctor --fix (TryLock) while it finishes a change and repairs links: a
// note is never written over another one, and doctor --fix never touches a
// change that is still being carried out.
func lockPath(p paths.Paths) string { return filepath.Join(p.State, "group-change.lock") }

// lockWait is how long Begin and Apply wait for the lock.
var lockWait = 10 * time.Second

// ErrChanging means another ccam holds the lock of the changes of groups:
// it carries one out, or doctor --fix finishes or repairs one.
var ErrChanging = errors.New("another ccam process is changing the groups right now")

// takeLock takes lockPath; then says what to do when another ccam holds it.
func takeLock(p paths.Paths, then string) (*lock.Flock, error) {
	l, err := lock.WaitFlock(lockPath(p), lockWait)
	if errors.Is(err, lock.ErrBusy) {
		return nil, fmt.Errorf("%w; %s", ErrChanging, then)
	}
	return l, err
}

// Lock is the lock of the changes of groups, held by doctor --fix while it
// finishes a stopped change (Resume) and repairs share links.
type Lock struct {
	p paths.Paths
	f *lock.Flock
}

// TryLock takes the lock without waiting: ErrChanging when another ccam
// holds it.
func TryLock(p paths.Paths) (*Lock, error) {
	f, err := lock.TryFlock(lockPath(p))
	if errors.Is(err, lock.ErrBusy) {
		return nil, ErrChanging
	}
	if err != nil {
		return nil, err
	}
	return &Lock{p: p, f: f}, nil
}

// Unlock lets go of the lock.
func (l *Lock) Unlock() { l.f.Unlock() }

// Unfinished refuses a new change of groups while another ccam carries
// one out (ErrChanging) or the last one stopped halfway (ErrPending).
func Unfinished(p paths.Paths) error {
	switch {
	case !Pending(p):
		return nil
	case Changing(p):
		return fmt.Errorf("%w; try again when it is done", ErrChanging)
	}
	return ErrPending
}

// Changing reports whether another ccam carries a change of groups out
// right now: then its note is not one that stopped.
func Changing(p paths.Paths) bool {
	l, err := lock.TryFlock(lockPath(p))
	if err != nil {
		return errors.Is(err, lock.ErrBusy)
	}
	l.Unlock()
	return false
}

// Begin notes the change, right before the config it leads to is saved,
// and holds lockPath until Apply or Abort; Apply removes the note when it
// is done, Abort when the config was not saved. It refuses while another
// change is noted (ErrPending).
func (ch *Change) Begin() error {
	if ch.Empty() {
		return nil
	}
	l, err := takeLock(ch.paths, "try again when it is done")
	if err != nil {
		return err
	}
	if Pending(ch.paths) {
		l.Unlock()
		return ErrPending
	}
	rec := record{Fresh: ch.fresh, Plugins: ch.groupPlugins}
	for _, a := range ch.Accounts {
		r := recorded{Name: a.Name, MCP: a.mcp, Seed: a.SeedMCP, Unsync: a.unsync, Keep: a.KeepMCP, DropKept: a.dropKept}
		for _, c := range a.Plan.Conflicts {
			if r.Choices == nil {
				r.Choices = map[string]string{}
			}
			r.Choices[c.Entry.Name] = c.Choice.String()
		}
		rec.Accounts = append(rec.Accounts, r)
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err == nil {
		err = fsutil.AtomicWrite(NotePath(ch.paths), data, 0o600)
	}
	if err != nil {
		l.Unlock()
		return fmt.Errorf("could not record the group change, so nothing was changed: %w", err)
	}
	ch.note, ch.held = data, l
	return nil
}

// Abort removes the note of a change whose config was not saved.
func (ch *Change) Abort() {
	ch.forget()
	ch.release()
}

// release lets go of the lock Begin took.
func (ch *Change) release() {
	if ch.held != nil {
		ch.held.Unlock()
		ch.held = nil
	}
}

// forget removes the change's note, unless another change noted itself
// over it meanwhile.
func (ch *Change) forget() {
	if ch.note == nil {
		return
	}
	if ch.noted() {
		os.Remove(NotePath(ch.paths))
	}
	ch.note = nil
}

// noted reports whether the change's note is still there: doctor --fix
// may have finished the change and removed it.
func (ch *Change) noted() bool {
	data, err := os.ReadFile(NotePath(ch.paths))
	return err == nil && bytes.Equal(data, ch.note)
}

// Apply makes the change on disk: the MCP steps first, which only read the
// groups' mcp.json files that no plan touches, then the plugin paths of
// the groups accounts leave, then the plans in order. Save the new config
// before (after Begin): if Apply stops halfway, `ccam doctor --fix`
// finishes the change.
func (ch *Change) Apply() error {
	if ch.held == nil {
		l, err := takeLock(ch.paths, "when it is done, finish this change with `ccam doctor --fix`")
		if err != nil {
			return err
		}
		ch.held = l
	}
	defer ch.release()
	if ch.note != nil && !ch.noted() {
		return nil // doctor --fix finished it meanwhile
	}
	if err := ch.applyMCP(); err != nil {
		return err
	}
	if err := freePlugins(ch.groupPlugins, ch.accountPlugins); err != nil {
		return fmt.Errorf("%w\n%s", err, finish)
	}
	for _, a := range ch.Accounts {
		if err := a.Plan.Apply(); err != nil {
			msg := finish
			if b := filepath.Join(a.Dir, share.BackupDir); fsutil.Exists(b) {
				msg = "The originals are in " + b + ". " + finish
			}
			return fmt.Errorf("account %s: %w\n%s", a.Name, err, msg)
		}
	}
	ch.forget()
	return nil
}

const finish = "Fix the cause, then finish with `ccam doctor --fix`"

// applyMCP does the MCP steps of every account, the seeds first.
func (ch *Change) applyMCP() error {
	for _, a := range ch.Accounts {
		if a.SeedMCP != nil {
			if err := mcp.Seed(a.SeedMCP[0], a.SeedMCP[1]); err != nil {
				return fmt.Errorf("account %s: could not copy the MCP servers: %w\n%s", a.Name, err, finish)
			}
		}
	}
	for _, a := range ch.Accounts {
		if err := a.applyMCP(ch.paths); err != nil {
			return fmt.Errorf("account %s: %w\n%s", a.Name, err, finish)
		}
	}
	return nil
}

// applyMCP does the account's MCP steps, all but the seed.
func (a *Account) applyMCP(p paths.Paths) error {
	state := launcher.MCPStateDir(p, a.Name)
	switch {
	case a.KeepMCP != "":
		if err := mcp.Keep(a.Dir, a.KeepMCP, state); err != nil {
			return fmt.Errorf("could not copy the group's MCP servers into .claude.json: %w", err)
		}
	case a.unsync:
		if err := mcp.Unsync(a.Dir, state); err != nil {
			return fmt.Errorf("could not remove the copied MCP servers: %w", err)
		}
	}
	if a.dropKept {
		if err := mcp.DropKept(a.Dir, p.GroupDir(a.mcp), state); err != nil {
			return fmt.Errorf("could not remove the kept copies of the group's MCP servers: %w", err)
		}
	}
	return nil
}

// readNote reads the note of a stopped change; nil, nil when there is none.
func readNote(p paths.Paths) (*record, []byte, error) {
	data, err := os.ReadFile(NotePath(p))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var rec record
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, nil, fmt.Errorf("could not read %s: %w; remove that file, then run `ccam doctor --fix` again to repair the links", NotePath(p), err)
	}
	return &rec, data, nil
}

// Preview returns the plans Resume would carry out for the change of
// groups that stopped halfway, in the change's order and with the answers
// given to its conflicts, without changing anything (nil when nothing is
// noted). Each plan counts on what the ones before it put into the groups.
func Preview(p paths.Paths, cfg *config.Config) (*Change, error) {
	rec, _, err := readNote(p)
	if rec == nil {
		return nil, err
	}
	ch := &Change{paths: p}
	pending := map[string]string{}
	for _, r := range rec.Accounts {
		prof, ok := cfg.Find(r.Name)
		if !ok {
			continue
		}
		a := &Account{Name: prof.Name, Dir: profile.Dir(p, prof), Imported: prof.Dir != ""}
		pl, err := share.Build(a.Dir, p.GroupsRoot(), prof.Share, prof.Share, share.Options{Fresh: rec.Fresh && !a.Imported, Pending: pending})
		if err != nil {
			return nil, fmt.Errorf("account %s: %w", a.Name, err)
		}
		choose(pl, r.Choices)
		for k, v := range pl.Created() {
			if old, ok := pending[k]; !ok || old == "" {
				pending[k] = v
			}
		}
		a.Plan = pl
		ch.Accounts = append(ch.Accounts, a)
	}
	return ch, nil
}

// Resume finishes a change of groups that stopped halfway (Pending): its
// MCP steps that still match cfg, its plugin paths, then the links of its
// accounts, in the order of the change and with the answers given to their
// conflicts. The caller holds l, taken before cfg was loaded. It refuses
// while one of the accounts runs. The change it returns has the plans it
// carried out (nil when the note was gone).
func Resume(l *Lock, cfg *config.Config, running func(name, dir string) bool) (*Change, error) {
	p := l.p
	rec, data, err := readNote(p)
	if rec == nil {
		return nil, err
	}
	ch := &Change{paths: p, note: data}
	var shares []config.Share
	var choices []map[string]string
	for _, r := range rec.Accounts {
		prof, ok := cfg.Find(r.Name)
		if !ok {
			continue
		}
		dir := profile.Dir(p, prof)
		if running(prof.Name, dir) {
			return nil, fmt.Errorf("%s is running; quit Claude Code there, then run `ccam doctor --fix` again to finish the group change", prof.Name)
		}
		a := &Account{Name: prof.Name, Dir: dir, Imported: prof.Dir != ""}
		// A config that was not saved as planned leads elsewhere: then the
		// account's MCP steps are not done.
		if cur := prof.Share.Get("mcp"); cur == r.MCP || cur == config.Local && r.MCP == "" {
			a.mcp, a.unsync, a.KeepMCP, a.dropKept = r.MCP, r.Unsync, r.Keep, r.DropKept
			if r.Seed != nil && fsutil.Exists(r.Seed[0]) {
				a.SeedMCP = r.Seed
			}
		}
		ch.Accounts = append(ch.Accounts, a)
		shares = append(shares, prof.Share)
		choices = append(choices, r.Choices)
	}
	if err := ch.applyMCP(); err != nil {
		return nil, err
	}
	if err := freePlugins(rec.Plugins, accountPlugins(p, cfg)); err != nil {
		return nil, err
	}
	// Each plan is made right before it runs, so that it sees what the
	// ones before it put into the groups.
	for i, a := range ch.Accounts {
		pl, err := share.Build(a.Dir, p.GroupsRoot(), shares[i], shares[i], share.Options{Fresh: rec.Fresh && !a.Imported})
		if err == nil {
			choose(pl, choices[i])
			err = pl.Apply()
		}
		if err != nil {
			return nil, fmt.Errorf("account %s: %w", a.Name, err)
		}
		a.Plan = pl
	}
	ch.forget()
	return ch, nil
}

// choose gives the conflicts of pl the answers noted for them, where the
// answer is still one of a conflict's candidates.
func choose(pl *share.Plan, choices map[string]string) {
	for _, c := range pl.Conflicts {
		for _, cand := range c.Candidates {
			if choices[c.Entry.Name] == cand.String() {
				c.Choice = cand
			}
		}
	}
}

// Noted lists the accounts of the change of groups that stopped halfway,
// or nil when there is none or its note cannot be read.
func Noted(p paths.Paths) []string {
	data, err := os.ReadFile(NotePath(p))
	if err != nil {
		return nil
	}
	var rec record
	if json.Unmarshal(data, &rec) != nil {
		return nil
	}
	names := make([]string, len(rec.Accounts))
	for i, r := range rec.Accounts {
		names[i] = r.Name
	}
	return names
}
