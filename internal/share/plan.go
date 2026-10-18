package share

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
	OpRewritePlugins           // point plugin registry paths at the group folder
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
	return [...]string{"그룹 것 유지", "내 것으로 덮기", "JSON 병합", "줄 단위 병합"}[r]
}

// Conflict is a file that exists in both places with different content.
type Conflict struct {
	Item       string
	Entry      Entry
	Mine       string
	Group      string
	GroupSave  string // where the group copy is saved before it is overwritten
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
}

// Build compares the current share settings with target and lists the
// steps to get there. It does not touch the disk except to read it.
func Build(dir, groupsRoot string, current, target config.Share, opt Options) (*Plan, error) {
	pl := &Plan{Dir: dir, stamp: fsutil.Stamp(time.Now())}
	for _, item := range config.Items {
		from, to := current.Get(item), target.Get(item)
		for _, e := range Layout[item] {
			if Never[e.Name] {
				return nil, fmt.Errorf("내부 오류: %s 는 공유할 수 없습니다", e.Name)
			}
			mine := filepath.Join(dir, e.Name)
			info, err := os.Lstat(mine)
			exists := err == nil
			isLink := exists && info.Mode()&os.ModeSymlink != 0
			var linkTarget string
			if isLink {
				linkTarget, _ = os.Readlink(mine)
			}

			if to == config.Local {
				if isLink && strings.HasPrefix(linkTarget, groupsRoot+string(filepath.Separator)) {
					pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine})
					if !opt.Fresh && fsutil.Exists(linkTarget) {
						pl.Steps = append(pl.Steps, Step{Op: OpCopyBack, Item: item, Entry: e, From: linkTarget, To: mine})
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
			case foreign && rerr == nil && !fsutil.Exists(grp):
				// A user-managed link (dotfiles): its content seeds the group,
				// the link itself goes to the backup, the target is untouched.
				pl.Steps = append(pl.Steps,
					Step{Op: OpCopyToGroup, Item: item, Entry: e, From: resolved, To: grp},
					Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
			case foreign:
				pl.Steps = append(pl.Steps, Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
				pl.Notes = append(pl.Notes, fmt.Sprintf("%s: 기존 링크(→ %s)는 백업하고 %s 그룹 내용을 씁니다", e.Name, linkTarget, to))
			case isLink:
				// Linked to another group: relink, leaving that group as it is.
				pl.Steps = append(pl.Steps, Step{Op: OpUnlink, Item: item, Entry: e, From: mine})
				if from != config.Local && from != to {
					pl.Notes = append(pl.Notes, fmt.Sprintf("%s: %s 그룹 내용은 그대로 두고 %s 그룹으로 연결합니다", e.Name, from, to))
				}
			case exists:
				if err := planAdopt(pl, item, e, mine, grp); err != nil {
					return nil, err
				}
			}
			pl.Steps = append(pl.Steps, Step{Op: OpLink, Item: item, Entry: e, From: grp, To: mine})
		}
	}
	return pl, nil
}

// planAdopt handles a real file/dir in the profile that should become shared.
func planAdopt(pl *Plan, item string, e Entry, mine, grp string) error {
	// Plugin registries hold absolute paths through this profile folder.
	rewrite := Step{Op: OpRewritePlugins, Item: item, Entry: e, From: mine + string(filepath.Separator), To: grp}
	if !fsutil.Exists(grp) {
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
				if fsutil.Exists(filepath.Join(mine, name)) && fsutil.Exists(filepath.Join(grp, name)) {
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
	same, err := sameFile(mine, grp)
	if err != nil {
		return err
	}
	if !same {
		c := &Conflict{Item: item, Entry: e, Mine: mine, Group: grp,
			GroupSave: filepath.Join(pl.Dir, BackupDir, pl.stamp, "group", e.Name)}
		switch {
		case e.Name == "history.jsonl":
			c.Candidates = []Resolution{MergeLines, KeepGroup, UseMine}
		case strings.HasSuffix(e.Name, ".json"):
			c.Candidates = []Resolution{KeepGroup, UseMine, MergeJSON}
		default:
			c.Candidates = []Resolution{KeepGroup, UseMine}
		}
		c.Choice = c.Candidates[0]
		pl.Conflicts = append(pl.Conflicts, c)
		pl.Steps = append(pl.Steps, Step{Op: OpWriteGroup, Item: item, Entry: e, From: mine, To: grp, Conflict: c})
	}
	pl.Steps = append(pl.Steps, Step{Op: OpBackup, Item: item, Entry: e, From: mine, To: pl.backupPath(e)})
	return nil
}

func (pl *Plan) backupPath(e Entry) string {
	return filepath.Join(pl.Dir, BackupDir, pl.stamp, e.Name)
}

func sameFile(a, b string) (bool, error) {
	x, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	y, err := os.ReadFile(b)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return bytes.Equal(x, y), err
}

// Describe writes the plan in plain Korean, one step per line.
func (pl *Plan) Describe(w io.Writer) {
	if pl.Empty() {
		fmt.Fprintln(w, "  바꿀 것이 없습니다.")
		return
	}
	for _, s := range pl.Steps {
		var line string
		switch s.Op {
		case OpAdopt:
			line = fmt.Sprintf("옮기기   %s → %s (그룹의 첫 내용)", s.From, s.To)
		case OpMergeDir:
			line = fmt.Sprintf("합치기   %s 의 새 파일 → %s", s.From, s.To)
		case OpWriteGroup:
			line = fmt.Sprintf("충돌     %s ↔ %s: %s", s.From, s.To, s.Conflict.Choice)
		case OpBackup:
			line = fmt.Sprintf("백업     %s → %s", s.From, s.To)
		case OpUnlink:
			line = fmt.Sprintf("링크 해제 %s", s.From)
		case OpLink:
			line = fmt.Sprintf("링크     %s → %s", s.To, s.From)
		case OpCopyBack, OpCopyToGroup:
			line = fmt.Sprintf("복사     %s → %s", s.From, s.To)
		case OpMergeRegistry:
			line = fmt.Sprintf("합치기   %s → %s (플러그인 목록)", s.From, s.To)
		case OpRewritePlugins:
			line = fmt.Sprintf("경로     %s 의 플러그인 경로를 그룹 폴더로", s.To)
		}
		fmt.Fprintln(w, "  "+line)
	}
	for _, n := range pl.Notes {
		fmt.Fprintln(w, "  참고: "+n)
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
		if err := fsutil.CopyFile(s.To, s.Conflict.GroupSave); err != nil {
			return err
		}
		mine, err := os.ReadFile(s.From)
		if err != nil {
			return err
		}
		grp, err := os.ReadFile(s.To)
		if err != nil {
			return err
		}
		merged, err := MergeJSONBytes(grp, mine)
		if err != nil {
			return err
		}
		return fsutil.AtomicWrite(s.To, merged, 0o600)
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
			if err := fsutil.Move(path, target); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		if fsutil.Exists(target) {
			return nil
		}
		return fsutil.Move(path, target)
	})
}

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
