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
	MissingLink ProblemKind = iota // nothing at the path; safe to link
	WrongTarget                    // a link, but to somewhere else
	NotLinked                      // a real file/dir where a link belongs
	StaleLink                      // item is local but the path still links to a group
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
	default:
		return fmt.Sprintf("%s: set to local but still linked to group %s", pr.Path, pr.Got)
	}
}

// Check compares a profile folder with its share settings using only
// lstat/readlink, so it is cheap enough for the launch path.
func Check(dir string, s config.Share, groupsRoot string) []Problem {
	var out []Problem
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
			}
		}
	}
	return out
}

// QuickFix repairs MissingLink problems (nothing to lose) and returns the
// ones it could not fix safely.
func QuickFix(problems []Problem) ([]Problem, error) {
	var rest []Problem
	for _, pr := range problems {
		if pr.Kind != MissingLink {
			rest = append(rest, pr)
			continue
		}
		if err := linkEntry(pr.Entry, pr.Want, pr.Path); err != nil {
			return nil, err
		}
	}
	return rest, nil
}

// linkEntry makes path a link to target (a junction for a folder on
// Windows without Developer Mode), creating the group dir for directory
// entries so the link is never dangling.
func linkEntry(e Entry, target, path string) error {
	if e.Dir {
		if err := os.MkdirAll(target, 0o700); err != nil {
			return err
		}
	} else if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return err
	}
	return fsutil.Link(target, path, e.Dir)
}
