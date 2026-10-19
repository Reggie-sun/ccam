// Package shim manages account shortcuts: symlinks in ~/.local/bin, named by
// the user, that point at the ccam binary, which picks the account from
// argv[0]. Shortcuts are opt-in (config Profile.Shortcut).
package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/jhyunwoo/ccam/internal/config"
)

var shortcutRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

// Lookup maps the name ccam was run as to the profile whose shortcut it is.
// ok is false for ccam itself and for any name no account uses.
func Lookup(cfg *config.Config, argv0 string) (string, bool) {
	if p, ok := cfg.FindShortcut(filepath.Base(argv0)); ok {
		return p.Name, true
	}
	return "", false
}

// Suggest is the shortcut offered for an account: "cc" + its name (work ->
// ccwork). Account names already use only allowed characters. Callers run
// Check on it and offer nothing when it fails.
func Suggest(name string) string { return "cc" + name }

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

// Install creates the link of every configured shortcut that is missing.
// It never removes anything: renaming or removing a shortcut deletes its
// link with Remove, and files that are not our links are never touched.
func Install(binDir, ccamPath string, cfg *config.Config) (created []string, err error) {
	want := map[string]bool{}
	for _, p := range cfg.Profiles {
		if p.Shortcut != "" {
			want[p.Shortcut] = true
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
	return slices.ContainsFunc(filepath.SplitList(os.Getenv("PATH")), func(d string) bool {
		return filepath.Clean(d) == filepath.Clean(dir) // "~/.local/bin/" is the same folder
	})
}
