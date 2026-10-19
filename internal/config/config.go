// Package config reads and writes ~/.config/ccam/config.toml.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"

	"github.com/pelletier/go-toml/v2"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

const CurrentVersion = 1

const (
	MCPFlag = "flag"
	MCPSync = "sync"
)

// Local means "not shared" for a share item.
const Local = "local"

type Config struct {
	Version     int       `toml:"version"`
	Default     string    `toml:"default,omitempty"`
	ClaudeBin   string    `toml:"claude_bin,omitempty"`
	MCPStrategy string    `toml:"mcp_strategy,omitempty"`
	NameShims   bool      `toml:"name_shims,omitempty"` // also create claude-<name>
	Profiles    []Profile `toml:"profile"`
}

type Profile struct {
	Name  string `toml:"name"`
	Slot  int    `toml:"slot"`
	Email string `toml:"email,omitempty"`
	Dir   string `toml:"dir,omitempty"` // set for imported folders registered in place
	Share Share  `toml:"share"`
}

// Share maps each shareable item to Local or a group name.
type Share struct {
	MCP        string `toml:"mcp,omitempty"`
	Settings   string `toml:"settings,omitempty"`
	Memory     string `toml:"memory,omitempty"`
	Extensions string `toml:"extensions,omitempty"`
	History    string `toml:"history,omitempty"`
}

// Items lists share item keys in display order.
var Items = []string{"mcp", "settings", "memory", "extensions", "history"}

func (s *Share) Get(item string) string {
	var v string
	switch item {
	case "mcp":
		v = s.MCP
	case "settings":
		v = s.Settings
	case "memory":
		v = s.Memory
	case "extensions":
		v = s.Extensions
	case "history":
		v = s.History
	}
	if v == "" {
		return Local
	}
	return v
}

func (s *Share) Set(item, group string) error {
	if group == "" || group == "off" {
		group = Local
	}
	if group == "on" {
		group = "default"
	}
	if group != Local && !nameRe.MatchString(group) {
		return fmt.Errorf("invalid group name: %q", group)
	}
	switch item {
	case "mcp":
		s.MCP = group
	case "settings":
		s.Settings = group
	case "memory":
		s.Memory = group
	case "extensions":
		s.Extensions = group
	case "history":
		s.History = group
	default:
		return fmt.Errorf("unknown share item: %q (mcp, settings, memory, extensions, history)", item)
	}
	return nil
}

// Equal compares by meaning: "" and "local" are the same.
func (s Share) Equal(o Share) bool {
	for _, item := range Items {
		if s.Get(item) != o.Get(item) {
			return false
		}
	}
	return true
}

// DefaultShare is what `ccam add` uses unless told otherwise.
func DefaultShare() Share {
	return Share{MCP: "default", Settings: "default", Memory: "default", Extensions: "default", History: Local}
}

// Reserved holds subcommand names that cannot be profile names,
// so `ccam <name>` is never ambiguous.
var Reserved = map[string]bool{
	"add": true, "list": true, "ls": true, "run": true, "use": true, "next": true,
	"usage": true, "u": true, "share": true, "mcp": true, "login": true, "logout": true,
	"rename": true, "remove": true, "rm": true, "import": true, "migrate": true, "doctor": true,
	"shims": true, "completion": true, "help": true, "version": true, "statusline": true,
}

var nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,31}$`)

func ValidateName(name string) error {
	if !nameRe.MatchString(name) {
		return fmt.Errorf("profile names are 1-32 characters of a-z, 0-9, - and _: %q", name)
	}
	if _, err := strconv.Atoi(name); err == nil {
		return fmt.Errorf("names cannot be all digits (they look like account numbers): %q", name)
	}
	if Reserved[name] {
		return fmt.Errorf("%q is a ccam command, so it cannot be a profile name", name)
	}
	return nil
}

// Load reads the config. A missing file yields an empty config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &Config{Version: CurrentVersion}, nil
	}
	if err != nil {
		return nil, err
	}
	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	if c.Version > CurrentVersion {
		return nil, fmt.Errorf("%s was written by a newer ccam (format version %d)", path, c.Version)
	}
	if c.Version == 0 {
		c.Version = CurrentVersion
	}
	return &c, nil
}

func (c *Config) Save(path string) error {
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(path, data, 0o600)
}

// Find resolves a profile by name or slot number.
func (c *Config) Find(sel string) (*Profile, bool) {
	if n, err := strconv.Atoi(sel); err == nil {
		for i := range c.Profiles {
			if c.Profiles[i].Slot == n {
				return &c.Profiles[i], true
			}
		}
		return nil, false
	}
	for i := range c.Profiles {
		if c.Profiles[i].Name == sel {
			return &c.Profiles[i], true
		}
	}
	return nil, false
}

// MustFind is Find with a user-facing error.
func (c *Config) MustFind(sel string) (*Profile, error) {
	if p, ok := c.Find(sel); ok {
		return p, nil
	}
	if len(c.Profiles) == 0 {
		return nil, fmt.Errorf("no accounts yet. Add one first with `ccam add <name>`")
	}
	return nil, fmt.Errorf("no account %q (see `ccam list`)", sel)
}

// DefaultProfile returns the default profile, or the lowest slot if none is set.
func (c *Config) DefaultProfile() (*Profile, error) {
	if c.Default != "" {
		return c.MustFind(c.Default)
	}
	if len(c.Profiles) == 0 {
		return nil, fmt.Errorf("no accounts yet. Add one first with `ccam add <name>`")
	}
	best := &c.Profiles[0]
	for i := range c.Profiles {
		if c.Profiles[i].Slot < best.Slot {
			best = &c.Profiles[i]
		}
	}
	return best, nil
}

// NextSlot returns the lowest unused slot number, starting at 1.
func (c *Config) NextSlot() int {
	used := map[int]bool{}
	for _, p := range c.Profiles {
		used[p.Slot] = true
	}
	n := 1
	for used[n] {
		n++
	}
	return n
}

// Remove drops a profile and clears the default if it pointed at it.
func (c *Config) Remove(name string) {
	out := c.Profiles[:0]
	for _, p := range c.Profiles {
		if p.Name != name {
			out = append(out, p)
		}
	}
	c.Profiles = out
	if c.Default == name {
		c.Default = ""
	}
}

func (c *Config) Strategy() string {
	if c.MCPStrategy == MCPSync {
		return MCPSync
	}
	return MCPFlag
}
