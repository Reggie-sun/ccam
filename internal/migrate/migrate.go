// Package migrate copies an existing Claude Code config folder (usually
// ~/.claude) into a share group, so accounts that share through that group
// start with the setup the user already has.
//
// The source is only read: nothing in it is moved, linked or changed.
// Only paths listed in share.Layout are considered, so .credentials.json and
// .claude.json are never copied; MCP servers are read out of .claude.json
// with profile.TopLevelValue (via mcp.LocalServers) and stored in the
// group's mcp.json.
package migrate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

// DefaultItems is what migrate copies when no item is named. History is
// left out because projects/ can be large.
var DefaultItems = []string{"settings", "memory", "extensions", "mcp"}

type Options struct {
	Source     string    // folder to copy from
	Home       string    // user home: ~/.claude keeps its .claude.json there
	GroupsRoot string    // paths.GroupsRoot()
	Group      string    // target group name
	Items      []string  // items to copy (keys of share.Layout)
	Accounts   []Account // accounts that may receive the source's MCP sign-ins
}

// State says what happens to one path.
type State int

const (
	Missing  State = iota // not in the source
	Linked                // the source already is the group's copy (a share link)
	Same                  // the group already has the same content
	New                   // the group lacks it: copied whole
	Merge                 // folder in both: files the group lacks are copied
	Conflict              // file in both with different content
)

// Entry is one share.Layout path of the source.
type Entry struct {
	Item       string
	Entry      share.Entry
	Source     string // path in the source folder (may be a link)
	From       string // Source with links resolved
	To         string // path in the group
	State      State
	Files      int   // New/Merge: files that will be copied
	Bytes      int64 // New/Merge: their total size
	Differ     int   // Merge: files in both that differ; the group's are kept
	Registries []string
	Conflict   *share.Conflict
}

// MCP is the plan for the user-scope MCP servers of the source.
type MCP struct {
	From      string // the .claude.json read ("" when there is none)
	New       []string
	Same      []string
	Conflicts []*share.Conflict // one per server; Entry.Name is "MCP server <name>"
	servers   map[string]json.RawMessage
	store     string // the group's mcp.json
}

type Plan struct {
	Source    string
	Group     string
	GroupDir  string
	BackupDir string // where replaced group files are saved
	Entries   []*Entry
	MCP       *MCP    // nil when mcp is not selected
	Logins    *Logins // nil when there are no MCP sign-ins to copy
	Conflicts []*share.Conflict
}

// ClaudeJSON returns the .claude.json that belongs to a config folder: inside
// it, or for ~/.claude without one inside, ~/.claude.json next to it
// (docs/claude-layout.md). "" when there is none.
func ClaudeJSON(source, home string) string {
	if p := filepath.Join(source, ".claude.json"); fsutil.Exists(p) {
		return p
	}
	if source == filepath.Join(home, ".claude") {
		if p := filepath.Join(home, ".claude.json"); fsutil.Exists(p) {
			return p
		}
	}
	return ""
}

// HasContent reports whether dir holds settings, memory or extensions worth
// migrating. Only lstat calls.
func HasContent(dir string) bool {
	for _, item := range []string{"settings", "memory", "extensions"} {
		for _, e := range share.Layout[item] {
			if fsutil.Exists(filepath.Join(dir, e.Name)) {
				return true
			}
		}
	}
	return false
}

// Found reports, for each item, whether the source has anything to copy:
// lstat calls for files and folders, the MCP server list for mcp.
func Found(source, home string) map[string]bool {
	found := map[string]bool{}
	for _, item := range config.Items {
		for _, e := range share.Layout[item] {
			if fsutil.Exists(filepath.Join(source, e.Name)) {
				found[item] = true
			}
		}
	}
	if from := ClaudeJSON(source, home); from != "" {
		servers, err := mcp.LocalServers(filepath.Dir(from))
		found["mcp"] = err == nil && len(servers) > 0
	}
	return found
}

// Build looks at the source and the group and lists what Apply will do.
// It only reads the disk.
func Build(opt Options) (*Plan, error) {
	src, err := filepath.Abs(opt.Source)
	if err != nil {
		return nil, err
	}
	if err := checkSource(src); err != nil {
		return nil, err
	}
	groupDir := filepath.Join(opt.GroupsRoot, opt.Group)
	pl := &Plan{Source: src, Group: opt.Group, GroupDir: groupDir}
	// One folder per run, even for two runs in the same second, so a backup
	// never overwrites an earlier one.
	pl.BackupDir = fsutil.UniqueDir(filepath.Join(groupDir, share.BackupDir, fsutil.Stamp(time.Now())))
	want := map[string]bool{}
	for _, item := range opt.Items {
		if _, ok := share.Layout[item]; !ok {
			return nil, fmt.Errorf("unknown item: %q", item)
		}
		want[item] = true
	}
	for _, item := range config.Items {
		if !want[item] {
			continue
		}
		if item == "mcp" {
			m, err := buildMCP(pl, src, opt.Home)
			if err != nil {
				return nil, err
			}
			pl.MCP = m
			continue
		}
		for _, e := range share.Layout[item] {
			if share.Never[e.Name] {
				return nil, fmt.Errorf("internal error: %s cannot be copied", e.Name)
			}
			en, err := buildEntry(pl, item, e)
			if err != nil {
				return nil, err
			}
			pl.Entries = append(pl.Entries, en)
		}
	}
	if want["mcp"] || want["extensions"] {
		if pl.Logins, err = buildLogins(pl, opt); err != nil {
			return nil, err
		}
	}
	return pl, nil
}

func checkSource(src string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("no Claude Code folder at %s. Give the folder to copy from: ccam migrate <folder>", src)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a folder: %s", src)
	}
	names, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, d := range names {
		if n := d.Name(); share.Known(n) && n != share.BackupDir && !share.Never[n] {
			return nil
		}
	}
	return fmt.Errorf("%s does not look like a Claude Code folder (no settings.json, CLAUDE.md, skills/, projects/ …)", src)
}

func buildEntry(pl *Plan, item string, e share.Entry) (*Entry, error) {
	en := &Entry{Item: item, Entry: e, Source: filepath.Join(pl.Source, e.Name), To: filepath.Join(pl.GroupDir, e.Name)}
	resolved, err := filepath.EvalSymlinks(en.Source)
	if err != nil {
		return en, nil // missing, or a dangling link
	}
	info, err := os.Stat(resolved)
	if err != nil || info.IsDir() != e.Dir {
		return en, nil
	}
	en.From = resolved
	if g, err := filepath.EvalSymlinks(en.To); err == nil && g == resolved {
		en.State = Linked
		return en, nil
	}
	if !fsutil.Exists(en.To) {
		en.State = New
		if !e.Dir {
			en.Files, en.Bytes = 1, info.Size()
			return en, nil
		}
		return en, scan(en)
	}
	if e.Dir {
		if err := scan(en); err != nil {
			return nil, err
		}
		if en.Files == 0 && en.Differ == 0 && len(en.Registries) == 0 {
			en.State = Same
		} else {
			en.State = Merge
		}
		return en, nil
	}
	same, err := fsutil.SameContent(resolved, en.To)
	if err != nil {
		return nil, err
	}
	if same {
		en.State = Same
		return en, nil
	}
	en.State = Conflict
	en.Conflict = share.NewConflict(item, e, en.Source, en.To, filepath.Join(pl.BackupDir, e.Name))
	pl.Conflicts = append(pl.Conflicts, en.Conflict)
	return en, nil
}

// scan counts what copying en.From into en.To adds, and which files exist
// on both sides with different content.
func scan(en *Entry) error {
	return filepath.WalkDir(en.From, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || share.Never[d.Name()] {
			return nil
		}
		rel, _ := filepath.Rel(en.From, path)
		target := filepath.Join(en.To, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !fsutil.Exists(target) {
			en.Files++
			en.Bytes += info.Size()
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		same, err := fsutil.SameContent(path, target)
		if err != nil || same {
			return err
		}
		if en.Entry.Name == "plugins" && rel == filepath.Base(rel) && isRegistry(rel) {
			en.Registries = append(en.Registries, rel)
			return nil
		}
		en.Differ++
		return nil
	})
}

func isRegistry(name string) bool {
	for _, n := range profile.PluginFiles {
		if n == name {
			return true
		}
	}
	return false
}

func buildMCP(pl *Plan, src, home string) (*MCP, error) {
	m := &MCP{From: ClaudeJSON(src, home), store: mcp.Path(pl.GroupDir)}
	if m.From == "" {
		return m, nil
	}
	servers, err := mcp.LocalServers(filepath.Dir(m.From))
	if err != nil {
		return nil, fmt.Errorf("could not read the MCP servers in %s: %w", m.From, err)
	}
	m.servers = servers
	store, err := mcp.Load(m.store)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(servers))
	for n := range servers {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		have, ok := store.Servers[n]
		switch {
		case !ok:
			m.New = append(m.New, n)
		case mcp.SameDef(have, servers[n]):
			m.Same = append(m.Same, n)
		default:
			c := &share.Conflict{
				Item:       "mcp",
				Entry:      share.Entry{Name: "MCP server " + n},
				Mine:       m.From + " (mcpServers." + n + ")",
				Group:      m.store,
				GroupSave:  filepath.Join(pl.BackupDir, "mcp.json"),
				Candidates: []share.Resolution{share.KeepGroup, share.UseMine},
			}
			m.Conflicts = append(m.Conflicts, c)
			pl.Conflicts = append(pl.Conflicts, c)
		}
	}
	return m, nil
}

func serverName(c *share.Conflict) string { return strings.TrimPrefix(c.Entry.Name, "MCP server ") }

// changes reports whether Apply would write anything.
func (en *Entry) changes() bool {
	switch en.State {
	case New:
		return true
	case Merge:
		return en.Files > 0 || len(en.Registries) > 0
	case Conflict:
		return en.Conflict.Choice != share.KeepGroup
	}
	return false
}

func (m *MCP) changes() bool {
	if len(m.New) > 0 {
		return true
	}
	for _, c := range m.Conflicts {
		if c.Choice != share.KeepGroup {
			return true
		}
	}
	return false
}

// Changes reports whether Apply would write anything with the current
// conflict choices.
func (pl *Plan) Changes() bool {
	for _, en := range pl.Entries {
		if en.changes() {
			return true
		}
	}
	return pl.MCP != nil && pl.MCP.changes() || pl.Logins != nil && pl.Logins.changes()
}

// Empty reports whether there is nothing to copy and nothing to decide.
func (pl *Plan) Empty() bool { return len(pl.Conflicts) == 0 && !pl.Changes() }

// Describe writes one line per path. asking says conflicts are still to be
// decided (otherwise their current choice is shown).
func (pl *Plan) Describe(w io.Writer, asking bool) {
	fmt.Fprintf(w, "Copy from %s into group %s (%s):\n", pl.Source, pl.Group, pl.GroupDir)
	last := ""
	label := func(item string) string {
		if item == last {
			return ""
		}
		last = item
		return item
	}
	// MCP first: the order of config.Items, which conflicts are asked in.
	if m := pl.MCP; m != nil {
		name := ".claude.json"
		if m.From != "" && filepath.Dir(m.From) != pl.Source {
			name = "~/.claude.json" // ClaudeJSON only looks outside the source for ~/.claude
		}
		var desc string
		switch {
		case m.From == "":
			desc = "not found, skipped"
		case len(m.servers) == 0:
			desc = "no MCP servers, skipped"
		default:
			var parts []string
			if len(m.New) > 0 {
				parts = append(parts, "new: "+strings.Join(m.New, ", "))
			}
			if len(m.Same) > 0 {
				parts = append(parts, "already in the group: "+strings.Join(m.Same, ", "))
			}
			if len(m.Conflicts) > 0 {
				names := make([]string, len(m.Conflicts))
				for i, c := range m.Conflicts {
					names[i] = serverName(c)
				}
				parts = append(parts, "differ: "+strings.Join(names, ", "))
			}
			desc = strings.Join(parts, "; ")
		}
		fmt.Fprintf(w, "  %-11s %-15s %s\n", label("mcp"), name, desc)
		for _, c := range m.Conflicts {
			fmt.Fprintf(w, "  %-11s %-15s %s\n", "", serverName(c), conflictText(c, asking))
		}
	}
	if l := pl.Logins; l != nil {
		fmt.Fprintf(w, "  %-11s %-15s %s\n", label("mcp"), "sign-ins", l.describe())
	}
	// Paths missing from the source go on one line per item, after the rest.
	var missing []string
	for i, en := range pl.Entries {
		name := en.Entry.Name
		if en.Entry.Dir {
			name += "/"
		}
		if en.State == Missing {
			missing = append(missing, name)
		} else {
			fmt.Fprintf(w, "  %-11s %-15s %s\n", label(en.Item), name, en.describe(asking))
		}
		if len(missing) > 0 && (i == len(pl.Entries)-1 || pl.Entries[i+1].Item != en.Item) {
			fmt.Fprintf(w, "  %-11s not found, skipped: %s\n", label(en.Item), strings.Join(missing, ", "))
			missing = nil
		}
	}
}

func (en *Entry) describe(asking bool) string {
	switch en.State {
	case Linked:
		return "already the group's copy (linked), skipped"
	case Same:
		return "same as the group's copy, skipped"
	case New:
		if en.Entry.Dir {
			return fmt.Sprintf("new: %s, %s", files(en.Files), Size(en.Bytes))
		}
		return "new: " + Size(en.Bytes)
	case Merge:
		var parts []string
		if en.Files > 0 {
			parts = append(parts, fmt.Sprintf("add %s the group lacks (%s)", files(en.Files), Size(en.Bytes)))
		}
		if len(en.Registries) > 0 {
			parts = append(parts, "merge the plugin list")
		}
		if en.Differ > 0 {
			verb := "differ"
			if en.Differ == 1 {
				verb = "differs"
			}
			parts = append(parts, fmt.Sprintf("%s %s, group copies kept", files(en.Differ), verb))
		}
		return strings.Join(parts, "; ")
	default:
		return conflictText(en.Conflict, asking)
	}
}

func conflictText(c *share.Conflict, asking bool) string {
	if asking {
		names := make([]string, len(c.Candidates))
		for i, r := range c.Candidates {
			names[i] = r.String()
		}
		return "differs from the group's copy; you choose: " + strings.Join(names, " / ")
	}
	return "differs from the group's copy: " + c.Choice.String()
}

func servers(n int) string {
	if n == 1 {
		return "1 MCP server"
	}
	return fmt.Sprintf("%d MCP servers", n)
}

func files(n int) string {
	if n == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", n)
}

// Size renders a byte count as B, KB, MB or GB (powers of 1024).
func Size(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	case n < 1<<30:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	default:
		return fmt.Sprintf("%.1f GB", float64(n)/(1<<30))
	}
}

// Apply copies into the group. The source is only read. Group files that
// get replaced are saved under BackupDir first.
func (pl *Plan) Apply() error {
	for _, en := range pl.Entries {
		if err := apply(pl, en); err != nil {
			return fmt.Errorf("%s: %w", en.Entry.Name, err)
		}
	}
	if pl.MCP != nil && pl.MCP.changes() {
		if err := pl.MCP.apply(pl); err != nil {
			return fmt.Errorf("MCP servers: %w", err)
		}
	}
	if pl.Logins != nil && pl.Logins.changes() {
		return pl.Logins.apply(pl)
	}
	return nil
}

func apply(pl *Plan, en *Entry) error {
	switch en.State {
	case New:
		if !en.Entry.Dir {
			data, err := os.ReadFile(en.From)
			if err != nil {
				return err
			}
			return fsutil.AtomicWrite(en.To, data, 0o600)
		}
		if err := copyMissing(en.From, en.To); err != nil {
			return err
		}
	case Merge:
		if err := copyMissing(en.From, en.To); err != nil {
			return err
		}
		for _, name := range en.Registries {
			if err := share.MergeRegistry(filepath.Join(en.From, name), filepath.Join(en.To, name), filepath.Join(pl.BackupDir, en.Entry.Name, name), en.Source, en.From); err != nil {
				return err
			}
		}
	case Conflict:
		return en.Conflict.Resolve()
	default:
		return nil
	}
	if en.Entry.Name != "plugins" {
		return nil
	}
	// Plugin registries hold absolute paths through the source folder; the
	// group's copy must point at its own files.
	sep := string(filepath.Separator)
	for _, prefix := range []string{en.Source + sep, en.From + sep} {
		if _, err := profile.RewritePluginPaths(en.To, prefix, en.To+sep); err != nil {
			return err
		}
	}
	return nil
}

// copyMissing copies every file, link and folder of src that dst lacks.
// Nothing in dst is overwritten and nothing in src changes.
func copyMissing(src, dst string) error {
	// Parents the walk would create with the source's mode stay private.
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if share.Never[d.Name()] {
			return nil
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		switch {
		case d.IsDir():
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(target, info.Mode().Perm()|0o700)
		case fsutil.Exists(target):
			return nil
		case fsutil.IsLink(path, d):
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return fsutil.CopyLink(fsutil.PortableLink(path, link, src), target)
		default:
			return fsutil.CopyFile(path, target)
		}
	})
}

func (m *MCP) apply(pl *Plan) error {
	store, err := mcp.Load(m.store)
	if err != nil {
		return err
	}
	if fsutil.Exists(m.store) {
		if err := fsutil.CopyFile(m.store, filepath.Join(pl.BackupDir, "mcp.json")); err != nil {
			return err
		}
	}
	for _, n := range m.New {
		store.Servers[n] = m.servers[n]
	}
	for _, c := range m.Conflicts {
		if c.Choice == share.UseMine {
			store.Servers[serverName(c)] = m.servers[serverName(c)]
		}
	}
	return store.Save(m.store)
}

// Copied lists what Apply wrote, for the summary ("settings.json",
// "skills/", "2 MCP servers").
func (pl *Plan) Copied() []string {
	var out []string
	if m := pl.MCP; m != nil && m.changes() {
		n := len(m.New)
		for _, c := range m.Conflicts {
			if c.Choice == share.UseMine {
				n++
			}
		}
		out = append(out, servers(n))
	}
	for _, en := range pl.Entries {
		if en.changes() {
			name := en.Entry.Name
			if en.Entry.Dir {
				name += "/"
			}
			out = append(out, name)
		}
	}
	if l := pl.Logins; l != nil && len(l.written) > 0 {
		out = append(out, "MCP sign-ins for "+strings.Join(l.written, ", "))
	}
	return out
}

// Kept lists conflicts where the group's version stayed.
func (pl *Plan) Kept() []string {
	var out []string
	for _, c := range pl.Conflicts {
		if c.Choice == share.KeepGroup {
			out = append(out, c.Entry.Name)
		}
	}
	return out
}

// Items lists the selected items that exist in the source, in display order.
func (pl *Plan) Items() []string {
	seen := map[string]bool{}
	var out []string
	for _, en := range pl.Entries {
		if en.State != Missing && !seen[en.Item] {
			seen[en.Item] = true
			out = append(out, en.Item)
		}
	}
	if pl.MCP != nil && len(pl.MCP.servers) > 0 {
		out = append(out, "mcp")
	}
	order := map[string]int{}
	for i, item := range config.Items {
		order[item] = i
	}
	sort.SliceStable(out, func(i, j int) bool { return order[out[i]] < order[out[j]] })
	return out
}
