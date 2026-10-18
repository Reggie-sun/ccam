// Package shim manages the claude1…claudeN and ccs symlinks in ~/.local/bin.
// Each one points at the ccam binary, which picks the profile from argv[0].
package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/jhyunwoo/ccam/internal/config"
)

var (
	slotRe = regexp.MustCompile(`^claude([0-9]+)$`)
	nameRe = regexp.MustCompile(`^claude-([a-z0-9][a-z0-9_-]*)$`)
)

// Selector maps an executable name to a profile selector.
// claude2 -> "2", claude-work -> "work". ok is false for ccam/ccs/anything else.
func Selector(argv0 string) (string, bool) {
	base := filepath.Base(argv0)
	if m := slotRe.FindStringSubmatch(base); m != nil {
		return m[1], true
	}
	if m := nameRe.FindStringSubmatch(base); m != nil {
		return m[1], true
	}
	return "", false
}

// Install makes binDir match cfg: one claudeN per profile, claude-<name> when
// enabled, and ccs. Stale links that point at ccam are removed; files that
// are not our links are never touched.
func Install(binDir, ccamPath string, cfg *config.Config) (created []string, err error) {
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return nil, err
	}
	want := map[string]bool{"ccs": true}
	for _, p := range cfg.Profiles {
		want["claude"+strconv.Itoa(p.Slot)] = true
		if cfg.NameShims {
			want["claude-"+p.Name] = true
		}
	}
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if want[name] || !isOurs(filepath.Join(binDir, name), ccamPath) {
			continue
		}
		if _, ok := Selector(name); ok {
			os.Remove(filepath.Join(binDir, name))
		}
	}
	for name := range want {
		link := filepath.Join(binDir, name)
		if isOurs(link, ccamPath) {
			continue
		}
		if _, err := os.Lstat(link); err == nil {
			return created, fmt.Errorf("%s 가 이미 있고 ccam 링크가 아니라서 건드리지 않았습니다", link)
		}
		if err := os.Symlink(ccamPath, link); err != nil {
			return created, err
		}
		created = append(created, name)
	}
	return created, nil
}

// RemoveFor deletes the shims of one profile.
func RemoveFor(binDir, ccamPath string, prof *config.Profile) {
	for _, name := range []string{"claude" + strconv.Itoa(prof.Slot), "claude-" + prof.Name} {
		link := filepath.Join(binDir, name)
		if isOurs(link, ccamPath) {
			os.Remove(link)
		}
	}
}

func isOurs(link, ccamPath string) bool {
	target, err := os.Readlink(link)
	if err != nil {
		return false
	}
	if target == ccamPath {
		return true
	}
	a, err1 := filepath.EvalSymlinks(link)
	b, err2 := filepath.EvalSymlinks(ccamPath)
	return err1 == nil && err2 == nil && a == b
}

// InPath reports whether dir is on PATH, so we can tell the user to add it.
func InPath(dir string) bool {
	for _, d := range filepath.SplitList(os.Getenv("PATH")) {
		if d == dir {
			return true
		}
	}
	return false
}
