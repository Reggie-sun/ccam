package share

import (
	"bytes"
	"errors"
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
	OpUnlink                   // remove the link From, which pointed to To when planned
	OpLink                     // create the link (and the group dir for dir entries)
	OpCopyBack                 // replace the link To with a copy of the group content From (sharing off)
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
	if _, ok := o.Pending[path]; ok {
		return true
	}
	// A folder an earlier plan moves in holds path if its source does. A
	// miss goes on looking: another plan may bring the file itself (above),
	// and the map is read in random order.
	for k, src := range o.Pending {
		if rest, ok := strings.CutPrefix(path, k+string(filepath.Separator)); ok && src != "" && fsutil.Exists(filepath.Join(src, rest)) {
			return true
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

// content is where the content of group path, or of a file inside a group
// folder, can be read now: an earlier plan may still have to move it in.
func (o Options) content(path string) string {
	if fsutil.Exists(path) {
		return path
	}
	if src := o.Pending[path]; src != "" {
		return src
	}
	for k, src := range o.Pending {
		if rest, ok := strings.CutPrefix(path, k+string(filepath.Separator)); ok && src != "" {
			if in := filepath.Join(src, rest); fsutil.Exists(in) {
				return in
			}
		}
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
				if stage := stagePath(dir, e); !exists && fsutil.Exists(stage) {
					// A copy back stopped after removing the link: the
					// copy is complete.
					pl.Steps = append(pl.Steps, Step{Op: OpAdopt, Item: item, Entry: e, From: stage, To: mine})
					continue
				}
				if isLink && strings.HasPrefix(linkTarget, groupsRoot+string(filepath.Separator)) {
					if opt.Fresh || !fsutil.Exists(linkTarget) {
						pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine, To: linkTarget})
						continue
					}
					// The copy takes the link's place once it is complete.
					pl.Steps = append(pl.Steps, Step{Op: OpCopyBack, Item: item, Entry: e, From: linkTarget, To: mine})
					if e.Name == "plugins" {
						// The copy's registries still point into the group.
						pl.Steps = append(pl.Steps, Step{Op: OpRewritePlugins, Item: item, Entry: e, From: linkTarget + string(filepath.Separator), To: mine})
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
				pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine, To: linkTarget})
			case exists && groupLinksTo(grp, mine):
				// The group's entry is a link to this very entry (made by
				// hand): it holds nothing of its own, so the account's
				// content takes its place.
				target, _ := fsutil.LinkTarget(grp)
				pl.Steps = append(pl.Steps,
					Step{Op: OpUnlink, Item: item, Entry: e, From: grp, To: target},
					Step{Op: OpAdopt, Item: item, Entry: e, From: mine, To: grp})
				if e.Name == "plugins" {
					pl.Steps = append(pl.Steps, Step{Op: OpRewritePlugins, Item: item, Entry: e, From: mine + string(filepath.Separator), To: grp})
				}
				pl.Notes = append(pl.Notes, fmt.Sprintf("%s: group %s held only a link to this account's own; its content becomes the group's", e.Name, to))
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
		// Inside a folder the group's copy of a file stays; say which of
		// the account's differ, since they then live only in the backup.
		differ := differing(e, mine, grp, opt)
		if n := len(differ); n > 0 {
			files, verb := "1 file", "differs"
			if n > 1 {
				files, verb = fmt.Sprintf("%d files", n), "differ"
			}
			if n > 3 {
				differ = append(differ[:3:3], "…")
			}
			pl.Notes = append(pl.Notes, fmt.Sprintf("%s: %s %s from the group's copy (%s); the group's is kept, the account's stays in %s",
				e.Name, files, verb, strings.Join(differ, ", "), pl.backupPath(e)))
		}
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

// differing lists the files in the folder mine (relative to it) that the
// group folder grp, as it will be when the plan runs, has with other
// content. Plugin registries are merged instead and left out, and so is
// what cannot be read: moving does not need to read it.
func differing(e Entry, mine, grp string, opt Options) []string {
	var out []string
	_ = filepath.WalkDir(mine, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.Type().IsRegular() {
			return nil
		}
		rel, _ := filepath.Rel(mine, path)
		if registry(e, rel) {
			return nil
		}
		theirs := opt.content(filepath.Join(grp, rel))
		if info, err := os.Lstat(theirs); err != nil || !info.Mode().IsRegular() {
			return nil
		}
		if same, err := fsutil.SameContent(path, theirs); err == nil && !same {
			out = append(out, rel)
		}
		return nil
	})
	return out
}

// registry reports whether rel, in the folder of entry e, is a plugin
// registry, which an OpMergeRegistry step merges.
func registry(e Entry, rel string) bool {
	return e.Name == "plugins" && slices.Contains(profile.PluginFiles, rel)
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
		if s.Op == OpAdopt && s.From == stagePath(pl.Dir, s.Entry) {
			lines = append(lines, line{op: s.Op, text: fmt.Sprintf("restore  %s: %s → %s (copied back from the group before a stop)", s.Entry.Name, from, to)})
			continue
		}
		if s.Op == OpCopyBack {
			// The copy takes the place of the link.
			u := key{OpUnlink, filepath.Dir(to), ""}
			if _, ok := names[u]; !ok {
				lines = append(lines, line{op: OpUnlink, key: &u})
			}
			names[u] = append(names[u], s.Entry.Name)
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
		if base != "" {
			n = strings.ReplaceAll(n, base+string(filepath.Separator), "")
		}
		fmt.Fprintln(w, "  Note: "+n)
	}
}

// Apply executes the plan in order. Every step either moves data into the
// group or into the backup folder; nothing is deleted but an account's
// copies of files the group holds byte for byte.
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
		return mergeDir(s.Entry, s.From, s.To)
	case OpWriteGroup:
		return resolve(s.Conflict)
	case OpBackup:
		return fsutil.Move(s.From, s.To)
	case OpUnlink:
		return unlink(s.From, s.To)
	case OpLink:
		return linkEntry(s.Entry, s.From, s.To)
	case OpCopyBack:
		// Staged in the backup folder (same disk, a name doctor knows).
		stage := stagePath(filepath.Dir(s.To), s.Entry)
		err := copyInto(s.Entry, s.From, s.To, stage, s.From)
		os.Remove(filepath.Dir(stage)) // when it was made only for the copy
		return err
	case OpCopyToGroup:
		return copyInto(s.Entry, s.From, s.To, filepath.Join(filepath.Dir(s.To), ".ccam-incomplete-"+s.Entry.Name), "")
	case OpMergeRegistry:
		return MergeRegistry(s.From, s.To, s.Conflict.GroupSave, filepath.Dir(s.From))
	case OpRewritePlugins:
		_, err := profile.RewritePluginPaths(s.To, s.From, s.To+string(filepath.Separator))
		return err
	}
	return nil
}

// unlink removes the link at path if it still points to target. What took
// its place since the plan was made (a file an editor saved there) stays,
// and the step fails; a link already gone is fine.
func unlink(path, target string) error {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if got, ok := fsutil.LinkTargetOf(path, info); !ok || got != target {
		return fmt.Errorf("%s changed since the plan was made, so it was left as it is", path)
	}
	return os.Remove(path)
}

// rename is os.Rename; tests replace it.
var rename = os.Rename

// copyInto copies src, entry e's content, to dst through stage, a path
// beside dst that never holds anything but this copy: dst changes only
// once the copy is complete, so a copy that stops halfway is never taken
// for a finished one. It replaces the link to src at dst (link, for a copy
// back) or an empty group folder (link ""); anything else there since the
// plan was made stays, and the step fails.
func copyInto(e Entry, src, dst, stage, link string) error {
	if err := os.RemoveAll(stage); err != nil { // left by a run that stopped
		return err
	}
	err := fsutil.CopyTree(src, stage)
	if err == nil {
		if link != "" {
			err = unlink(dst, link)
		} else if info, lerr := os.Lstat(dst); lerr == nil && !info.IsDir() {
			err = fmt.Errorf("%s is there since the plan was made, so it was left as it is", dst)
		} else if lerr == nil {
			err = os.Remove(dst) // fails unless the folder is still empty
		}
	}
	if err != nil {
		os.RemoveAll(stage)
		return err
	}
	if err := rename(stage, dst); err != nil {
		// With the link back the copy is made again by the next run;
		// without, the account would find nothing there.
		if link != "" && fsutil.Link(link, dst, e.Dir) != nil {
			return fmt.Errorf("%w (the copy is in %s)", err, stage)
		}
		os.RemoveAll(stage)
		return err
	}
	return nil
}

// mergeDir moves every file of src that dst lacks into dst and drops the
// ones dst has byte for byte, which a backup would only pile up. Files
// present in both with different content stay in src (which then goes to
// backup), and so do plugin registries, which a later step merges. The
// folders this leaves empty go too; src itself stays for the backup step.
func mergeDir(e Entry, src, dst string) error {
	var emptied []string
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if path != src && sameFile(path, target) {
			// The group holds a link the user made into this account:
			// it is the only copy, which takes the link's place. Reached
			// through a linked folder above, it stays for the backup.
			if _, ok := fsutil.LinkTarget(target); ok {
				if err := os.Remove(target); err != nil {
					return err
				}
				if err := fsutil.MoveWithin(path, target, src); err != nil {
					return err
				}
			}
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			if fsutil.Exists(target) {
				if path != src {
					emptied = append(emptied, path)
				}
				return nil
			}
			if err := fsutil.MoveWithin(path, target, src); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		if !fsutil.Exists(target) {
			return fsutil.MoveWithin(path, target, src)
		}
		if !d.Type().IsRegular() || registry(e, rel) {
			return nil
		}
		if info, err := os.Lstat(target); err != nil || !info.Mode().IsRegular() {
			return nil
		}
		if same, err := fsutil.SameContent(path, target); err != nil || !same {
			return nil
		}
		return os.Remove(path)
	})
	if err != nil {
		return err
	}
	// Deepest first; a folder that still holds something stays.
	for _, dir := range slices.Backward(emptied) {
		os.Remove(dir)
	}
	return nil
}

// stagePath is where a copy back of entry e into the account folder dir
// is made before it takes the link's place.
func stagePath(dir string, e Entry) string {
	return filepath.Join(dir, BackupDir, "incomplete-"+e.Name)
}

// groupLinksTo reports whether grp is a link that resolves to mine.
func groupLinksTo(grp, mine string) bool {
	_, link := fsutil.LinkTarget(grp)
	return link && sameFile(grp, mine)
}

// sameFile reports whether a and b, links followed, are one file.
func sameFile(a, b string) bool {
	ia, err := os.Stat(a)
	if err != nil {
		return false
	}
	ib, err := os.Stat(b)
	return err == nil && os.SameFile(ia, ib)
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
