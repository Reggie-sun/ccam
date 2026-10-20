package share

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
)

type Op int

const (
	OpAdopt          Op = iota // move mine into an empty group slot (first sharer)
	OpMergeDir                 // move my files the group lacks into the group dir
	OpWriteGroup               // write the resolved conflict content into the group
	OpBackup                   // move mine into .ccam-backup/<stamp>/
	OpUnlink                   // remove the link
	OpLink                     // create the link (and the group dir for dir entries)
	OpCopyBack                 // copy the group content into the profile (sharing off)
	OpCopyToGroup              // copy a foreign link's target into an empty group slot
	OpMergeRegistry            // merge my plugin registry JSON into the group's (group backed up)
	OpRewritePlugins           // point plugin registry paths in To from prefix From to To
)

// Resolution decides a file conflict when turning sharing on.
type Resolution int

const (
	KeepGroup Resolution = iota
	UseMine
	MergeJSON  // settings.json: deep merge, group wins on scalar clashes
	MergeLines // history.jsonl: group lines then my new lines
)

func (r Resolution) String() string {
	return [...]string{"keep group", "use mine", "merge JSON", "merge by line"}[r]
}

// Conflict is a file that exists in both places with different content.
type Conflict struct {
	Item      string
	Entry     Entry
	Mine      string
	Group     string
	GroupSave string // where the group copy is saved before it is overwritten
	// Seed is where the group copy is while it is planned but not yet
	// moved into the group by an earlier plan of the same change.
	Seed       string
	Choice     Resolution
	Candidates []Resolution
}

type Step struct {
	Op       Op
	Item     string
	Entry    Entry
	From, To string
	Conflict *Conflict
}

type Plan struct {
	Dir       string
	Steps     []Step
	Conflicts []*Conflict
	Notes     []string
	stamp     string
}

func (pl *Plan) Empty() bool { return len(pl.Steps) == 0 }

// Options tune Build.
type Options struct {
	Fresh bool // sharing off: start empty instead of copying the group content
	// Pending lists group paths that plans applied before this one will
	// create, each mapped to where its content is now ("" for a new empty
	// folder). Build treats them as already there.
	Pending map[string]string
}

// exists reports whether path will exist when the plan runs.
func (o Options) exists(path string) bool {
	if fsutil.Exists(path) {
		return true
	}
	for k, src := range o.Pending {
		if path == k {
			return true
		}
		if rest, ok := strings.CutPrefix(path, k+string(filepath.Separator)); ok && src != "" {
			return fsutil.Exists(filepath.Join(src, rest))
		}
	}
	return false
}

// empty reports whether group path holds nothing when the plan runs:
// missing, an empty folder, or a folder an earlier plan only creates.
func (o Options) empty(path string) bool {
	if src, ok := o.Pending[path]; ok && !fsutil.Exists(path) {
		return src == ""
	}
	if !fsutil.Exists(path) {
		return !o.exists(path)
	}
	entries, err := os.ReadDir(path)
	return err == nil && len(entries) == 0
}

// content is where the content of group path can be read now.
func (o Options) content(path string) string {
	if src := o.Pending[path]; src != "" && !fsutil.Exists(path) {
		return src
	}
	return path
}

// Created lists the group paths the plan creates, for Options.Pending of
// the plans after it.
func (pl *Plan) Created() map[string]string {
	out := map[string]string{}
	for _, s := range pl.Steps {
		switch s.Op {
		case OpAdopt, OpCopyToGroup:
			out[s.To] = s.From
		case OpLink:
			if s.Entry.Dir {
				if _, ok := out[s.From]; !ok {
					out[s.From] = ""
				}
			}
		case OpMergeDir:
			// Plugin registries the group lacks move in with the merge;
			// later accounts then merge theirs into it.
			if s.Entry.Name == "plugins" {
				for _, name := range profile.PluginFiles {
					src, dst := filepath.Join(s.From, name), filepath.Join(s.To, name)
					if _, ok := out[dst]; !ok && fsutil.Exists(src) && !fsutil.Exists(dst) {
						out[dst] = src
					}
				}
			}
		}
	}
	return out
}

// Build compares the current share settings with target and lists the
// steps to get there. It does not touch the disk except to read it.
func Build(dir, groupsRoot string, current, target config.Share, opt Options) (*Plan, error) {
	// One backup folder per plan, even for two plans in the same second.
	stamp := fsutil.UniqueDir(filepath.Join(dir, BackupDir, fsutil.Stamp(time.Now())))
	pl := &Plan{Dir: dir, stamp: filepath.Base(stamp)}
	for _, item := range config.Items {
		from, to := current.Get(item), target.Get(item)
		for _, e := range Layout[item] {
			if Never[e.Name] {
				return nil, fmt.Errorf("internal error: %s cannot be shared", e.Name)
			}
			mine := filepath.Join(dir, e.Name)
			info, err := os.Lstat(mine)
			exists := err == nil
			var linkTarget string
			isLink := false
			if exists {
				linkTarget, isLink = fsutil.LinkTargetOf(mine, info)
			}

			if to == config.Local {
				if isLink && strings.HasPrefix(linkTarget, groupsRoot+string(filepath.Separator)) {
					pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine})
					if !opt.Fresh && fsutil.Exists(linkTarget) {
						pl.Steps = append(pl.Steps, Step{Op: OpCopyBack, Item: item, Entry: e, From: linkTarget, To: mine})
						if e.Name == "plugins" {
							// The copy's registries still point into the group.
							pl.Steps = append(pl.Steps, Step{Op: OpRewritePlugins, Item: item, Entry: e, From: linkTarget + string(filepath.Separator), To: mine})
						}
					}
				}
				continue
			}

			grp := filepath.Join(groupsRoot, to, e.Name)
			if isLink && linkTarget == grp {
				continue
			}
			foreign := isLink && !strings.HasPrefix(linkTarget, groupsRoot+string(filepath.Separator))
			resolved, rerr := filepath.EvalSymlinks(mine)
			switch {
			case foreign && rerr == nil && !opt.exists(grp):
				// A user-managed link (dotfiles): its content seeds the group,
				// the link itself goes to the backup, the target is untouched.
				pl.Steps = append(pl.Steps,
					Step{Op: OpCopyToGroup, Item: item, Entry: e, From: resolved, To: grp},
					Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
			case foreign:
				pl.Steps = append(pl.Steps, Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
				pl.Notes = append(pl.Notes, fmt.Sprintf("%s: backing up the existing link (→ %s) and using group %s", e.Name, linkTarget, to))
			case isLink:
				// Linked to another group: relink, leaving that group as it
				// is. A new group starts with a copy of what the account saw.
				if opt.empty(grp) && fsutil.Exists(linkTarget) {
					pl.Steps = append(pl.Steps, Step{Op: OpCopyToGroup, Item: item, Entry: e, From: linkTarget, To: grp})
					if e.Name == "plugins" {
						pl.Steps = append(pl.Steps, Step{Op: OpRewritePlugins, Item: item, Entry: e, From: linkTarget + string(filepath.Separator), To: grp})
					}
				} else if from != config.Local && from != to {
					pl.Notes = append(pl.Notes, fmt.Sprintf("%s: leaving group %s as it is and linking to group %s", e.Name, from, to))
				}
				pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine})
			case exists:
				if err := planAdopt(pl, item, e, mine, grp, opt); err != nil {
					return nil, err
				}
			}
			pl.Steps = append(pl.Steps, Step{Op: OpLink, Item: item, Entry: e, From: grp, To: mine})
		}
	}
	return pl, nil
}

// planAdopt handles a real file/dir in the profile that should become shared.
func planAdopt(pl *Plan, item string, e Entry, mine, grp string, opt Options) error {
	// Plugin registries hold absolute paths through this profile folder.
	rewrite := Step{Op: OpRewritePlugins, Item: item, Entry: e, From: mine + string(filepath.Separator), To: grp}
	if !opt.exists(grp) {
		pl.Steps = append(pl.Steps, Step{Op: OpAdopt, Item: item, Entry: e, From: mine, To: grp})
		if e.Name == "plugins" {
			pl.Steps = append(pl.Steps, rewrite)
		}
		return nil
	}
	if e.Dir {
		pl.Steps = append(pl.Steps, Step{Op: OpMergeDir, Item: item, Entry: e, From: mine, To: grp})
		if e.Name == "plugins" {
			for _, name := range profile.PluginFiles {
				if fsutil.Exists(filepath.Join(mine, name)) && opt.exists(filepath.Join(grp, name)) {
					pl.Steps = append(pl.Steps, Step{Op: OpMergeRegistry, Item: item, Entry: e,
						From: filepath.Join(mine, name), To: filepath.Join(grp, name),
						Conflict: &Conflict{GroupSave: filepath.Join(pl.Dir, BackupDir, pl.stamp, "group", e.Name, name)}})
				}
			}
		}
		pl.Steps = append(pl.Steps, Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
		if e.Name == "plugins" {
			pl.Steps = append(pl.Steps, rewrite)
		}
		return nil
	}
	same, err := fsutil.SameContent(mine, opt.content(grp))
	if err != nil {
		return err
	}
	if !same {
		c := NewConflict(item, e, mine, grp, filepath.Join(pl.Dir, BackupDir, pl.stamp, "group", e.Name))
		if src := opt.content(grp); src != grp {
			c.Seed = src
		}
		pl.Conflicts = append(pl.Conflicts, c)
		pl.Steps = append(pl.Steps, Step{Op: OpWriteGroup, Item: item, Entry: e, From: mine, To: grp, Conflict: c})
	}
	pl.Steps = append(pl.Steps, Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
	return nil
}

// NewConflict lists the resolutions that fit the file and picks the first
// as the default choice. groupSave is where the group copy is saved before
// any choice other than KeepGroup overwrites it.
func NewConflict(item string, e Entry, mine, group, groupSave string) *Conflict {
	c := &Conflict{Item: item, Entry: e, Mine: mine, Group: group, GroupSave: groupSave}
	switch {
	case e.Name == "history.jsonl":
		c.Candidates = []Resolution{MergeLines, KeepGroup, UseMine}
	case strings.HasSuffix(e.Name, ".json"):
		c.Candidates = []Resolution{KeepGroup, UseMine, MergeJSON}
	default:
		c.Candidates = []Resolution{KeepGroup, UseMine}
	}
	c.Choice = c.Candidates[0]
	return c
}

func (pl *Plan) backupPath(e Entry) string {
	return filepath.Join(pl.Dir, BackupDir, pl.stamp, e.Name)
}

// Describe writes the plan in plain words, one step per line.
func (pl *Plan) Describe(w io.Writer) { pl.DescribeUnder(w, "") }

// describeRank orders the lines of a description: what goes into the
// group, conflicts, backups, then the links.
var describeRank = map[Op]int{
	OpAdopt: 0, OpCopyToGroup: 0, OpMergeDir: 1, OpWriteGroup: 2, OpMergeRegistry: 2,
	OpBackup: 3, OpUnlink: 4, OpCopyBack: 5, OpLink: 6, OpRewritePlugins: 7,
}

// DescribeUnder is Describe with paths below base written relative to it.
// Steps of one kind between the same two folders share a line.
func (pl *Plan) DescribeUnder(w io.Writer, base string) {
	if pl.Empty() {
		fmt.Fprintln(w, "  Nothing to change.")
		return
	}
	short := func(path string) string {
		if rel, ok := strings.CutPrefix(path, base+string(filepath.Separator)); ok && base != "" {
			return rel
		}
		return path
	}
	type key struct {
		op       Op
		from, to string
	}
	type line struct {
		op   Op
		key  *key // a line that lists entries
		text string
	}
	var lines []line
	names := map[key][]string{}
	for _, s := range pl.Steps {
		from, to := short(s.From), short(s.To)
		var k key
		switch s.Op {
		case OpAdopt, OpMergeDir, OpBackup, OpCopyBack, OpCopyToGroup:
			k = key{s.Op, filepath.Dir(from), filepath.Dir(to)}
		case OpLink:
			k = key{s.Op, filepath.Dir(to), filepath.Dir(from)}
		case OpUnlink:
			k = key{s.Op, filepath.Dir(from), ""}
		case OpWriteGroup:
			lines = append(lines, line{op: s.Op, text: fmt.Sprintf("conflict %s ↔ %s: %s", from, to, s.Conflict.Choice)})
			continue
		case OpMergeRegistry:
			lines = append(lines, line{op: s.Op, text: fmt.Sprintf("merge    %s → %s (plugin list)", from, to)})
			continue
		case OpRewritePlugins:
			lines = append(lines, line{op: s.Op, text: fmt.Sprintf("paths    plugin paths under %s → %s", strings.TrimSuffix(from, string(filepath.Separator)), to)})
			continue
		}
		if _, ok := names[k]; !ok {
			lines = append(lines, line{op: s.Op, key: &k})
		}
		names[k] = append(names[k], s.Entry.Name)
	}
	slices.SortStableFunc(lines, func(a, b line) int { return describeRank[a.op] - describeRank[b.op] })
	for _, l := range lines {
		if k := l.key; k != nil {
			list := strings.Join(names[*k], ", ")
			switch k.op {
			case OpAdopt:
				l.text = fmt.Sprintf("move     %s: %s → %s (first content of the group)", list, k.from, k.to)
			case OpMergeDir:
				l.text = fmt.Sprintf("merge    new files in %s: %s → %s", list, k.from, k.to)
			case OpBackup:
				l.text = fmt.Sprintf("back up  %s: %s → %s", list, k.from, k.to)
			case OpUnlink:
				l.text = fmt.Sprintf("unlink   %s in %s", list, k.from)
			case OpLink:
				l.text = fmt.Sprintf("link     %s: %s → %s", list, k.from, k.to)
			case OpCopyBack, OpCopyToGroup:
				l.text = fmt.Sprintf("copy     %s: %s → %s", list, k.from, k.to)
			}
		}
		fmt.Fprintln(w, "  "+l.text)
	}
	for _, n := range pl.Notes {
		fmt.Fprintln(w, "  Note: "+n)
	}
}

// Apply executes the plan in order. Every step either moves data into the
// group or into the backup folder; nothing is deleted.
func (pl *Plan) Apply() error {
	for _, s := range pl.Steps {
		if err := apply(s); err != nil {
			return fmt.Errorf("%s: %w", s.Entry.Name, err)
		}
	}
	return nil
}

func apply(s Step) error {
	switch s.Op {
	case OpAdopt:
		return fsutil.Move(s.From, s.To)
	case OpMergeDir:
		return mergeDir(s.From, s.To)
	case OpWriteGroup:
		return resolve(s.Conflict)
	case OpBackup:
		return fsutil.Move(s.From, s.To)
	case OpUnlink:
		return os.Remove(s.From)
	case OpLink:
		return linkEntry(s.Entry, s.From, s.To)
	case OpCopyBack, OpCopyToGroup:
		return fsutil.CopyTree(s.From, s.To)
	case OpMergeRegistry:
		return MergeRegistry(s.From, s.To, s.Conflict.GroupSave, filepath.Dir(s.From))
	case OpRewritePlugins:
		_, err := profile.RewritePluginPaths(s.To, s.From, s.To+string(filepath.Separator))
		return err
	}
	return nil
}

// mergeDir moves every file of src that dst lacks into dst. Files present
// in both with different content stay in src (which then goes to backup).
func mergeDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if fsutil.Exists(target) {
				return nil
			}
			if err := fsutil.MoveWithin(path, target, src); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		if fsutil.Exists(target) {
			return nil
		}
		return fsutil.MoveWithin(path, target, src)
	})
}

// Resolve writes c.Choice into the group file, saving the group copy to
// c.GroupSave first unless the group's version is kept.
func (c *Conflict) Resolve() error { return resolve(c) }

func resolve(c *Conflict) error {
	if c.Choice == KeepGroup {
		return nil
	}
	// Every other choice rewrites the shared file; keep what other
	// accounts had.
	if err := fsutil.CopyFile(c.Group, c.GroupSave); err != nil {
		return err
	}
	switch c.Choice {
	case UseMine:
		data, err := os.ReadFile(c.Mine)
		if err != nil {
			return err
		}
		return fsutil.AtomicWrite(c.Group, data, 0o600)
	case MergeJSON:
		mine, err := os.ReadFile(c.Mine)
		if err != nil {
			return err
		}
		grp, err := os.ReadFile(c.Group)
		if err != nil {
			return err
		}
		merged, err := MergeJSONBytes(grp, mine)
		if err != nil {
			return err
		}
		return fsutil.AtomicWrite(c.Group, merged, 0o600)
	case MergeLines:
		mine, err := os.ReadFile(c.Mine)
		if err != nil {
			return err
		}
		grp, err := os.ReadFile(c.Group)
		if err != nil {
			return err
		}
		return fsutil.AtomicWrite(c.Group, mergeLines(grp, mine), 0o600)
	}
	return nil
}

func mergeLines(base, extra []byte) []byte {
	seen := map[string]bool{}
	var out bytes.Buffer
	for _, chunk := range [][]byte{base, extra} {
		for _, line := range strings.Split(string(chunk), "\n") {
			if line == "" || seen[line] {
				continue
			}
			seen[line] = true
			out.WriteString(line)
			out.WriteByte('\n')
		}
	}
	return out.Bytes()
}
