// Package shim manages account shortcuts: symlinks in ~/.local/bin, named by
// the user, that point at the ccam binary, which picks the account from
// argv[0]. Shortcuts are opt-in (config Profile.Shortcut).
//
// On Windows a shortcut is <name>.exe: a symbolic link where Windows allows
// one (Developer Mode), else a hard link, which `ccam update` renews.
package shim

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
)

// exe is the file name suffix of programs.
var exe = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

// Name is the command argv0 runs: its file name, without .exe on Windows.
func Name(argv0 string) string {
	base := filepath.Base(argv0)
	if runtime.GOOS == "windows" {
		base = strings.ToLower(base)
		base = strings.TrimSuffix(base, ".exe")
	}
	return base
}

// file is the path of shortcut name in binDir.
func file(binDir, name string) string { return filepath.Join(binDir, name+exe) }

var shortcutRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

// Lookup maps the name ccam was run as to the profile whose shortcut it is.
// ok is false for ccam itself and for any name no account uses.
func Lookup(cfg *config.Config, argv0 string) (string, bool) {
	if p, ok := cfg.FindShortcut(Name(argv0)); ok {
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
	link := file(binDir, name)
	if _, err := os.Lstat(link); err == nil && !isOurs(link, ccamPath) {
		return fmt.Errorf("%s already exists and is not a ccam shortcut", link)
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" || filepath.Clean(dir) == filepath.Clean(binDir) {
			continue
		}
		if path, ok := program(dir, name); ok && !isOurs(path, ccamPath) {
			return fmt.Errorf("%q would hide %s", name, path)
		}
	}
	return nil
}

// program finds a program called name in dir: an executable file, or on
// Windows a file with one of the PATHEXT extensions.
func program(dir, name string) (string, bool) {
	if runtime.GOOS != "windows" {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		return path, err == nil && !info.IsDir() && info.Mode()&0o111 != 0
	}
	exts := os.Getenv("PATHEXT")
	if exts == "" {
		exts = ".COM;.EXE;.BAT;.CMD"
	}
	for _, ext := range filepath.SplitList(exts) {
		path := filepath.Join(dir, name+strings.ToLower(ext))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, true
		}
	}
	return "", false
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
		link := file(binDir, name)
		if isOurs(link, ccamPath) {
			continue
		}
		if _, err := os.Lstat(link); err == nil {
			return created, fmt.Errorf("%s already exists and is not a ccam link, so it was left alone", link)
		}
		if err := linkProgram(ccamPath, link); err != nil {
			return created, err
		}
		created = append(created, name)
	}
	return created, nil
}

// symlink is os.Symlink; tests replace it to see the Windows fallback.
var symlink = os.Symlink

// linkProgram links link to ccamPath: a symbolic link, or on Windows
// without the right to make one, a hard link.
func linkProgram(ccamPath, link string) error {
	err := symlink(ccamPath, link)
	if err == nil || runtime.GOOS != "windows" {
		return err
	}
	if herr := os.Link(ccamPath, link); herr != nil {
		return errors.Join(err, herr)
	}
	return nil
}

// Relink renews the shortcuts that are hard links to old, the binary
// `ccam update` replaced, so they run ccamPath. A running shortcut cannot
// be deleted on Windows, so it is renamed to <name>.exe.old first.
func Relink(binDir, ccamPath, old string, cfg *config.Config) []string {
	oldInfo, err := os.Stat(old)
	if err != nil {
		return nil
	}
	var renewed []string
	for _, p := range cfg.Profiles {
		if p.Shortcut == "" {
			continue
		}
		link := file(binDir, p.Shortcut)
		info, err := os.Lstat(link)
		if err != nil || !info.Mode().IsRegular() || !os.SameFile(info, oldInfo) {
			continue
		}
		os.Remove(link + ".old")
		if os.Remove(link) != nil && os.Rename(link, link+".old") != nil {
			continue
		}
		if linkProgram(ccamPath, link) == nil {
			renewed = append(renewed, p.Shortcut)
		}
	}
	return renewed
}

// Remove deletes the shortcut link name if it is ours.
func Remove(binDir, ccamPath, name string) {
	if name == "" {
		return
	}
	if link := file(binDir, name); isOurs(link, ccamPath) {
		os.Remove(link)
	}
}

// Linked reports whether binDir/name is our link.
func Linked(binDir, ccamPath, name string) bool {
	return isOurs(file(binDir, name), ccamPath)
}

// isOurs reports whether link runs ccamPath: a symbolic link to it, or a
// hard link of it (Windows shortcuts without the right to symlink).
func isOurs(link, ccamPath string) bool {
	target, err := os.Readlink(link)
	if err != nil {
		info, err := os.Lstat(link)
		self, serr := os.Stat(ccamPath)
		return err == nil && serr == nil && info.Mode().IsRegular() && os.SameFile(info, self)
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
