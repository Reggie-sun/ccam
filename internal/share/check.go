package share

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
)

type ProblemKind int

const (
	MissingLink   ProblemKind = iota // nothing at the path; safe to link
	WrongTarget                      // a link, but to somewhere else
	NotLinked                        // a real file/dir where a link belongs
	StaleLink                        // item is local but the path still links to a group
	MissingTarget                    // the right link, but the group folder it needs is gone
	Staged                           // item is local, nothing at the path, and a stopped copy back left the copy in Got
	StagedAside                      // the same, but something was made at the path since (a launch): nothing moves the copy
)

type Problem struct {
	Kind  ProblemKind
	Item  string
	Entry Entry
	Path  string // in the profile
	Want  string // expected link target ("" for StaleLink)
	Got   string // current link target, if any
}

func (pr Problem) String() string {
	switch pr.Kind {
	case MissingLink:
		return fmt.Sprintf("%s: share link is missing (→ %s)", pr.Path, pr.Want)
	case WrongTarget:
		return fmt.Sprintf("%s: points elsewhere (%s, expected %s)", pr.Path, pr.Got, pr.Want)
	case NotLinked:
		return fmt.Sprintf("%s: set to shared but is a real file, not a link", pr.Path)
	case MissingTarget:
		return fmt.Sprintf("%s: the group folder it links to is missing (%s), so Claude Code cannot save there", pr.Path, needed(pr.Entry, pr.Want))
	case Staged:
		return fmt.Sprintf("%s: missing, since copying it back from the group stopped; the copy is in %s", pr.Path, pr.Got)
	case StagedAside:
		return fmt.Sprintf("%s: copying it back from the group stopped, and it was made anew since; the complete copy is in %s: move what you need from it, then remove that folder", pr.Path, pr.Got)
	default:
		return fmt.Sprintf("%s: set to local but still linked to group %s", pr.Path, pr.Got)
	}
}

// Safe reports whether QuickFix repairs pr: there is nothing to lose.
func (pr Problem) Safe() bool {
	return pr.Kind == MissingLink || pr.Kind == MissingTarget || pr.Kind == Staged
}

// Waits reports whether pr is an account's own copy, or own link, of a
// path added to an item in a later ccam (Entry.Late): it never stops a
// launch, and doctor --fix shares it (the group's copy wins if it has one).
func (pr Problem) Waits(groupsRoot string) bool {
	return pr.Entry.Late && (pr.Kind == NotLinked ||
		pr.Kind == WrongTarget && !strings.HasPrefix(pr.Got, groupsRoot+string(filepath.Separator)))
}

// Check compares a profile folder with its share settings using only
// lstat/readlink, so it is cheap enough for the launch path.
func Check(dir string, s config.Share, groupsRoot string) []Problem {
	return check(dir, s, groupsRoot, false)
}

// CheckTargets is Check that also finds right links whose group folder is
// gone (MissingTarget). That is an lstat more per link, which a launch does
// not spend (a third more time); doctor does.
func CheckTargets(dir string, s config.Share, groupsRoot string) []Problem {
	return check(dir, s, groupsRoot, true)
}

func check(dir string, s config.Share, groupsRoot string, targets bool) []Problem {
	var out []Problem
	var there map[string]bool // group folders found there or not, by path
	if targets {
		there = map[string]bool{}
	}
	for _, item := range config.Items {
		group := s.Get(item)
		for _, e := range Layout[item] {
			path := filepath.Join(dir, e.Name)
			info, err := os.Lstat(path)
			var target string
			isLink := false
			if err == nil {
				target, isLink = fsutil.LinkTargetOf(path, info)
			}
			if group == config.Local {
				if isLink && strings.HasPrefix(target, groupsRoot+string(filepath.Separator)) {
					out = append(out, Problem{Kind: StaleLink, Item: item, Entry: e, Path: path, Got: target})
				}
				if stage := stagePath(dir, e); targets && fsutil.Exists(stage) {
					kind := Staged
					if err == nil {
						kind = StagedAside
					}
					out = append(out, Problem{Kind: kind, Item: item, Entry: e, Path: path, Got: stage})
				}
				continue
			}
			want := filepath.Join(groupsRoot, group, e.Name)
			switch {
			case err != nil:
				out = append(out, Problem{Kind: MissingLink, Item: item, Entry: e, Path: path, Want: want})
			case !isLink:
				out = append(out, Problem{Kind: NotLinked, Item: item, Entry: e, Path: path, Want: want})
			case target != want:
				out = append(out, Problem{Kind: WrongTarget, Item: item, Entry: e, Path: path, Want: want, Got: target})
			case targets:
				// A file not written yet is fine; the folder it goes in, or a
				// linked folder, must be there for Claude Code to save into.
				folder := needed(e, want)
				ok, seen := there[folder]
				if !seen {
					ok = fsutil.Exists(folder)
					there[folder] = ok
				}
				if !ok {
					out = append(out, Problem{Kind: MissingTarget, Item: item, Entry: e, Path: path, Want: want, Got: target})
				}
			}
		}
	}
	return out
}

// needed is the group folder a link to target needs: the folder itself for
// a directory entry, the one holding the file otherwise.
func needed(e Entry, target string) string {
	if e.Dir {
		return target
	}
	return filepath.Dir(target)
}

// QuickFix repairs MissingLink, MissingTarget and Staged problems (nothing
// to lose) and returns the ones it could not fix safely.
func QuickFix(problems []Problem) ([]Problem, error) {
	var rest []Problem
	for _, pr := range problems {
		var err error
		switch pr.Kind {
		case MissingLink:
			err = linkEntry(pr.Entry, pr.Want, pr.Path)
		case MissingTarget:
			err = os.MkdirAll(needed(pr.Entry, pr.Want), 0o700)
		case Staged:
			if err = fsutil.Move(pr.Got, pr.Path); err == nil {
				os.Remove(filepath.Dir(pr.Got)) // when it was made only for the copy
			}
		default:
			rest = append(rest, pr)
		}
		if err != nil {
			return nil, err
		}
	}
	return rest, nil
}

// linkEntry makes path a link to target (a junction for a folder on
// Windows without Developer Mode), creating the group dir for directory
// entries so the link is never dangling.
func linkEntry(e Entry, target, path string) error {
	if err := os.MkdirAll(needed(e, target), 0o700); err != nil {
		return err
	}
	return fsutil.Link(target, path, e.Dir)
}
