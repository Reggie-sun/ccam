// Package paths resolves every directory ccam touches.
//
// CCAM_HOME overrides everything so tests never touch the real home folder.
// Otherwise the XDG base directory variables are honoured.
package paths

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Config  string // ~/.config/ccam
	Data    string // ~/.local/share/ccam
	State   string // ~/.local/state/ccam
	Bin     string // ~/.local/bin
	Runtime string // $XDG_RUNTIME_DIR/ccam, for lock files
	Home    string // user home, used by import
}

func FromEnv() Paths {
	if root := os.Getenv("CCAM_HOME"); root != "" {
		return Paths{
			Config:  filepath.Join(root, "config"),
			Data:    filepath.Join(root, "data"),
			State:   filepath.Join(root, "state"),
			Bin:     filepath.Join(root, "bin"),
			Runtime: filepath.Join(root, "run"),
			Home:    filepath.Join(root, "home"),
		}
	}
	home, _ := os.UserHomeDir()
	p := Paths{
		Config: filepath.Join(xdg("XDG_CONFIG_HOME", home, ".config"), "ccam"),
		Data:   filepath.Join(xdg("XDG_DATA_HOME", home, ".local/share"), "ccam"),
		State:  filepath.Join(xdg("XDG_STATE_HOME", home, ".local/state"), "ccam"),
		Bin:    filepath.Join(home, ".local", "bin"),
		Home:   home,
	}
	if rt := os.Getenv("XDG_RUNTIME_DIR"); rt != "" {
		p.Runtime = filepath.Join(rt, "ccam")
	} else {
		p.Runtime = filepath.Join(p.State, "locks")
	}
	return p
}

func xdg(env, home, def string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return filepath.Join(home, def)
}

func (p Paths) ConfigFile() string            { return filepath.Join(p.Config, "config.toml") }
func (p Paths) ProfileDir(name string) string { return filepath.Join(p.Data, "profiles", name) }
func (p Paths) GroupsRoot() string            { return filepath.Join(p.Data, "groups") }
func (p Paths) GroupDir(name string) string   { return filepath.Join(p.GroupsRoot(), name) }
func (p Paths) TrashDir() string              { return filepath.Join(p.Data, "trash") }
