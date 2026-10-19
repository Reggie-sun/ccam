// Package shim manages account shortcuts: symlinks in ~/.local/bin, named by
// the user, that point at the ccam binary, which picks the account from
// argv[0]. Shortcuts are opt-in (config Profile.Shortcut).
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
	slotRe     = regexp.MustCompile(`^claude([0-9]+)$`)
	nameRe     = regexp.MustCompile(`^claude-([a-z0-9][a-z0-9_-]*)$`)
	shortcutRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)
)

// legacyAlias is the short ccs alias older versions linked to ccam.
const legacyAlias = "ccs"

// Lookup maps the name ccam was run as to a profile selector: a configured
// shortcut, or for a config from before shortcuts were stored (version 1)
// the old claude<slot> / claude-<name> links. ok is false for ccam itself.
func Lookup(cfg *config.Config, argv0 string) (string, bool) {
	base := filepath.Base(argv0)
	if p, ok := cfg.FindShortcut(base); ok {
		return p.Name, true
	}
	if cfg.Version < 2 {
		return legacySelector(base)
	}
	return "", false
}

// Stale reports whether argv0 looks like a shortcut of an older ccam
// (claude<N>, claude-<name>) that the config no longer has.
func Stale(argv0 string) bool {
	_, ok := legacySelector(filepath.Base(argv0))
	return ok
}

// legacySelector: claude2 -> "2", claude-work -> "work".
func legacySelector(base string) (string, bool) {
	if m := slotRe.FindStringSubmatch(base); m != nil {
		return m[1], true
	}
	if m := nameRe.FindStringSubmatch(base); m != nil {
		return m[1], true
	}
	return "", false
}

// Suggest is the shortcut offered for an account: "cc" + its name (work ->
// ccwork). Account names already use only allowed characters. Callers run
// Check on it and offer nothing when it fails.
func Suggest(name string) string { return "cc" + name }

// Adopt upgrades a version 1 config, where every profile had a claude<slot>
// link: each profile keeps the old link that is still on disk and points at
// ccam (claude<slot>, else claude-<name>) as its shortcut. Profiles without
// one get none. The caller saves the config.
func Adopt(binDir, ccamPath string, cfg *config.Config) {
	if cfg.Version >= 2 {
		return
	}
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		if p.Shortcut != "" {
			continue
		}
		for _, name := range []string{"claude" + strconv.Itoa(p.Slot), "claude-" + p.Name} {
			if _, taken := cfg.FindShortcut(name); !taken && isOurs(filepath.Join(binDir, name), ccamPath) {
				p.Shortcut = name
				break
			}
		}
	}
	cfg.Version = 2
}

// Check validates a shortcut name for profile prof: a safe file name, not a
// ccam command, not another account's shortcut, and not hiding another
// program on PATH (or a file in binDir that is not our link).
func Check(name, binDir, ccamPath string, cfg *config.Config, prof string) error {
	if !shortcutRe.MatchString(name) {
		return fmt.Errorf("shortcut names are 1-64 characters of a-z, 0-9, '.', '-' and '_', starting with a letter or digit: %q", name)
	}
	if name == "ccam" || name == "claude" || config.Reserved[name] {
		return fmt.Errorf("%q is taken by ccam or Claude Code", name)
	}
	if p, ok := cfg.FindShortcut(name); ok && p.Name != prof {
		return fmt.Errorf("%q is already the shortcut of account %s", name, p.Name)
	}
	link := filepath.Join(binDir, name)
	if _, err := os.Lstat(link); err == nil && !isOurs(link, ccamPath) {
		return fmt.Errorf("%s already exists and is not a ccam shortcut", link)
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" || filepath.Clean(dir) == filepath.Clean(binDir) {
			continue
		}
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() && info.Mode()&0o111 != 0 && !isOurs(path, ccamPath) {
			return fmt.Errorf("%q would hide %s", name, path)
		}
	}
	return nil
}

// Install makes binDir match cfg: one link per configured shortcut. Old
// claude<N>, claude-<name> and ccs links that point at ccam and are no
// longer configured are removed; files that are not our links are never
// touched.
func Install(binDir, ccamPath string, cfg *config.Config) (created []string, err error) {
	want := map[string]bool{}
	for _, p := range cfg.Profiles {
		if p.Shortcut != "" {
			want[p.Shortcut] = true
		}
	}
	entries, err := os.ReadDir(binDir)
	if err != nil && len(want) == 0 {
		return nil, nil // nothing to create, nothing to clean
	}
	for _, e := range entries {
		name := e.Name()
		if want[name] || !isOurs(filepath.Join(binDir, name), ccamPath) {
			continue
		}
		if _, ok := legacySelector(name); ok || name == legacyAlias {
			os.Remove(filepath.Join(binDir, name))
		}
	}
	if len(want) == 0 {
		return nil, nil
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return nil, err
	}
	for name := range want {
		link := filepath.Join(binDir, name)
		if isOurs(link, ccamPath) {
			continue
		}
		if _, err := os.Lstat(link); err == nil {
			return created, fmt.Errorf("%s already exists and is not a ccam link, so it was left alone", link)
		}
		if err := os.Symlink(ccamPath, link); err != nil {
			return created, err
		}
		created = append(created, name)
	}
	return created, nil
}

// Remove deletes the shortcut link name if it is ours.
func Remove(binDir, ccamPath, name string) {
	if name == "" {
		return
	}
	if link := filepath.Join(binDir, name); isOurs(link, ccamPath) {
		os.Remove(link)
	}
}

// Linked reports whether binDir/name is our link.
func Linked(binDir, ccamPath, name string) bool {
	return isOurs(filepath.Join(binDir, name), ccamPath)
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
