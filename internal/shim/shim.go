// Package shim manages account shortcuts: symlinks in ~/.local/bin, named by
// the user, that point at the ccam binary, which picks the account from
// argv[0]. Shortcuts are opt-in (config Profile.Shortcut).
//
// On Windows a shortcut is <name>.exe: a symbolic link where Windows allows
// one (Developer Mode), else a hard link. A hard link keeps running the
// binary an update or a new install replaced; `ccam update`, install.ps1
// (through `ccam shims install`) and `ccam doctor --fix` renew it.
//
// ccam records each link it makes and the binary it made it run (state
// shortcuts.json): only such a link counts as ccam's once that binary moved
// or was replaced (stale). Other links and copies of ccam are the user's.
package shim

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
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
// program on PATH (or a file in the bin folder that is not our link, or
// prof's own shortcut that ccam made).
func Check(name string, p paths.Paths, ccamPath string, cfg *config.Config, prof string) error {
	if !config.IsShortcut(name) {
		return fmt.Errorf("shortcut names are 1-64 characters of a-z, 0-9, '.', '-' and '_', starting with a letter or digit: %q", name)
	}
	if name == "ccam" || name == "claude" || config.Reserved[name] {
		return fmt.Errorf("%q is taken by ccam or Claude Code", name)
	}
	own := false
	if q, ok := cfg.FindShortcut(name); ok && q.Name != prof {
		return fmt.Errorf("%q is already the shortcut of account %s", name, q.Name)
	} else if ok {
		own = true
	}
	link := file(p.Bin, name)
	if _, err := os.Lstat(link); err == nil && !isOurs(link, ccamPath) {
		if rec := readRecord(p); !own || !made(link, rec) && !stale(link, rec) {
			return fmt.Errorf("%s already exists and is not a ccam shortcut", link)
		}
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" || filepath.Clean(dir) == filepath.Clean(p.Bin) {
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

// Install creates the link of every configured shortcut that is missing,
// and points the stale ones (see stale) at ccamPath; a link ccam made for
// another ccam that is still there stays as it is. A file in the way of one
// shortcut does not stop the others; the error names each. It never
// removes anything else: renaming or removing a shortcut deletes its link
// with Remove, and files that are not ccam's links are never touched.
func Install(p paths.Paths, ccamPath string, cfg *config.Config) (created []string, err error) {
	var want []string
	for _, p := range cfg.Profiles {
		if p.Shortcut != "" && !slices.Contains(want, p.Shortcut) {
			want = append(want, p.Shortcut)
		}
	}
	if len(want) == 0 {
		return nil, nil
	}
	if err := os.MkdirAll(p.Bin, 0o755); err != nil {
		return nil, err
	}
	slices.Sort(want)
	rec := readRecord(p)
	var errs []error
	for _, name := range want {
		link := file(p.Bin, name)
		if isOurs(link, ccamPath) {
			rec.add(link, runs(link, ccamPath)) // made before ccam recorded its links
			continue
		}
		if _, err := os.Lstat(link); err == nil {
			isStale := stale(link, rec)
			if !isStale && made(link, rec) {
				continue // it runs the ccam it was made for
			}
			if !isStale {
				errs = append(errs, fmt.Errorf("%s already exists and is not a ccam link, so it was left alone", link))
				continue
			}
			if err := remove(link); err != nil {
				errs = append(errs, fmt.Errorf("%s runs another ccam, and could not be replaced: %w", link, err))
				continue
			}
		}
		if err := linkProgram(ccamPath, link); err != nil {
			errs = append(errs, err)
			continue
		}
		rec.add(link, runs(link, ccamPath))
		created = append(created, name)
	}
	if err := rec.save(p); err != nil {
		errs = append(errs, err)
	}
	return created, errors.Join(errs...)
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

// Relink renews the stale shortcuts that are copies of another ccam
// binary, as hard links (Windows) of the one `ccam update` or install.ps1
// replaced are, so they run ccamPath.
func Relink(p paths.Paths, ccamPath string, cfg *config.Config) []string {
	rec := readRecord(p)
	var renewed []string
	for _, prof := range cfg.Profiles {
		if prof.Shortcut == "" {
			continue
		}
		link := file(p.Bin, prof.Shortcut)
		if isOurs(link, ccamPath) || !staleCopy(link) || !stale(link, rec) || remove(link) != nil {
			continue
		}
		if linkProgram(ccamPath, link) == nil {
			rec.add(link, runs(link, ccamPath))
			renewed = append(renewed, prof.Shortcut)
		}
	}
	_ = rec.save(p)
	return renewed
}

// remove deletes a shortcut, or one that is running (Windows refuses to
// delete it) renames to <name>.exe.old.
func remove(link string) error {
	if os.Remove(link) == nil {
		return nil
	}
	_, err := fsutil.MoveAside(link)
	return err
}

// Remove deletes the shortcut link name if it is ours, or one ccam made.
func Remove(p paths.Paths, ccamPath, name string) {
	if name == "" {
		return
	}
	rec := readRecord(p)
	if link := file(p.Bin, name); isOurs(link, ccamPath) || made(link, rec) || stale(link, rec) {
		os.Remove(link)
		delete(rec.Links, link)
		_ = rec.save(p)
	}
}

// Other returns the ccam binary shortcut name runs when it is a link ccam
// made for another ccam that is still there (Install leaves those alone).
func Other(p paths.Paths, ccamPath, name string) (string, bool) {
	link, rec := file(p.Bin, name), readRecord(p)
	if isOurs(link, ccamPath) || !made(link, rec) || stale(link, rec) {
		return "", false
	}
	return rec.Links[link], true
}

// Repoint points the shortcuts names that run another ccam (Other) at
// ccamPath, which doctor --fix asks for, and returns the ones it changed.
func Repoint(p paths.Paths, ccamPath string, names []string) ([]string, error) {
	rec := readRecord(p)
	var done []string
	var errs []error
	for _, name := range names {
		link := file(p.Bin, name)
		if isOurs(link, ccamPath) || !made(link, rec) || stale(link, rec) {
			continue
		}
		if err := remove(link); err != nil {
			errs = append(errs, fmt.Errorf("%s runs another ccam, and could not be replaced: %w", link, err))
			continue
		}
		if err := linkProgram(ccamPath, link); err != nil {
			errs = append(errs, err)
			continue
		}
		rec.add(link, runs(link, ccamPath))
		done = append(done, name)
	}
	if err := rec.save(p); err != nil {
		errs = append(errs, err)
	}
	return done, errors.Join(errs...)
}

// Stale reports whether the shortcut name runs another ccam (see stale).
func Stale(p paths.Paths, name string) bool {
	return stale(file(p.Bin, name), readRecord(p))
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

// record lists the links ccam made, each with the binary it made it run.
type record struct {
	Links map[string]string `json:"links"`
	dirty bool
}

func recordPath(p paths.Paths) string { return filepath.Join(p.State, "shortcuts.json") }

func readRecord(p paths.Paths) *record {
	r := &record{}
	if data, err := os.ReadFile(recordPath(p)); err == nil {
		_ = json.Unmarshal(data, r)
	}
	if r.Links == nil {
		r.Links = map[string]string{}
	}
	return r
}

func (r *record) add(link, ccamPath string) {
	if r.Links[link] != ccamPath {
		r.Links[link], r.dirty = ccamPath, true
	}
}

func (r *record) save(p paths.Paths) error {
	if !r.dirty {
		return nil
	}
	data, err := fsutil.MarshalJSON(r)
	if err == nil {
		err = os.MkdirAll(p.State, 0o700)
	}
	if err == nil {
		err = fsutil.AtomicWrite(recordPath(p), data, 0o600)
	}
	return err
}

// runs is the binary link runs, as the record keeps it: a symbolic link's
// target (from the link's folder when relative), else ccamPath, of which
// ccam made link a hard link.
func runs(link, ccamPath string) string {
	target, err := os.Readlink(link)
	if err != nil {
		return ccamPath
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(link), target)
	}
	return target
}

// made reports whether link is still what ccam made (recorded): a symbolic
// link to the binary it recorded, or a copy of a ccam binary (a hard link).
func made(link string, r *record) bool {
	target, ok := r.Links[link]
	if !ok {
		return false
	}
	if _, err := os.Readlink(link); err == nil {
		return runs(link, "") == target
	}
	return staleCopy(link)
}

// stale reports whether link, which ccam made, no longer runs the binary
// it made it run: that one moved or was replaced (an update, a reinstall),
// so the link points nowhere or, as a hard link, still runs the older
// binary. On Windows a copy of a ccam binary that ccam did not record is
// stale too: hard links made before ccam recorded its links. Removing a
// link never touches what it points to.
func stale(link string, r *record) bool {
	if _, ok := r.Links[link]; !ok {
		return runtime.GOOS == "windows" && staleCopy(link)
	}
	if !made(link, r) {
		return false // changed since: the user's now
	}
	self, err := os.Stat(r.Links[link])
	if err != nil {
		return true
	}
	info, err := os.Stat(link)
	return err != nil || !os.SameFile(info, self)
}

// ccamBuild is in the build information Go writes into every ccam binary.
const ccamBuild = "path\tgithub.com/jhyunwoo/ccam/cmd/ccam\nmod\tgithub.com/jhyunwoo/ccam"

// staleCopy reports whether link is a file of a ccam binary: on Windows a
// hard link of one that an update or a reinstall has replaced since, which
// still runs the older version.
func staleCopy(link string) bool {
	info, err := os.Lstat(link)
	if err != nil || !info.Mode().IsRegular() {
		return false
	}
	data, err := os.ReadFile(link)
	return err == nil && bytes.Contains(data, []byte(ccamBuild))
}

// InPath reports whether dir is on PATH, so we can tell the user to add it.
func InPath(dir string) bool {
	return slices.ContainsFunc(filepath.SplitList(os.Getenv("PATH")), func(d string) bool {
		return filepath.Clean(d) == filepath.Clean(dir) // "~/.local/bin/" is the same folder
	})
}
