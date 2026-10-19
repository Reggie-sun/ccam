// Package config reads and writes ~/.config/ccam/config.toml.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

const (
	MCPFlag = "flag"
	MCPSync = "sync"
)

// Local means "not shared" for a share item.
const Local = "local"

type Config struct {
	Default     string    `toml:"default,omitempty"`
	ClaudeBin   string    `toml:"claude_bin,omitempty"`
	MCPStrategy string    `toml:"mcp_strategy,omitempty"`
	Profiles    []Profile `toml:"profile"`
}

type Profile struct {
	Name  string `toml:"name"`
	Slot  int    `toml:"slot"`
	Email string `toml:"email,omitempty"`
	Dir   string `toml:"dir,omitempty"` // set for imported folders registered in place
	// Shortcut is the command in the bin folder that runs this account
	// ("" = none). It links to the ccam binary, which looks argv[0] up here.
	Shortcut string `toml:"shortcut,omitempty"`
	Share    Share  `toml:"share"`
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

// field returns the setting of item, or nil for an unknown item.
func (s *Share) field(item string) *string {
	switch item {
	case "mcp":
		return &s.MCP
	case "settings":
		return &s.Settings
	case "memory":
		return &s.Memory
	case "extensions":
		return &s.Extensions
	case "history":
		return &s.History
	}
	return nil
}

func (s *Share) Get(item string) string {
	if f := s.field(item); f != nil && *f != "" {
		return *f
	}
	return Local
}

func (s *Share) Set(item, group string) error {
	if group == "" || group == "off" {
		group = Local
	}
	if group == "on" {
		group = "default"
	}
	if group != Local {
		if err := ValidateGroup(group); err != nil {
			return err
		}
	}
	f := s.field(item)
	if f == nil {
		return fmt.Errorf("unknown share item: %q (mcp, settings, memory, extensions, history)", item)
	}
	*f = group
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

// Summary renders the shared items, e.g. "mcp settings md ext history@team";
// "" when nothing is shared.
func (s Share) Summary() string {
	short := map[string]string{"mcp": "mcp", "settings": "settings", "memory": "md", "extensions": "ext", "history": "history"}
	var parts []string
	for _, item := range Items {
		g := s.Get(item)
		if g == Local {
			continue
		}
		if g != "default" {
			parts = append(parts, short[item]+"@"+g)
		} else {
			parts = append(parts, short[item])
		}
	}
	return strings.Join(parts, " ")
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
	"shims": true, "shortcut": true, "completion": true, "help": true, "version": true, "statusline": true,
	"update": true,
}

// isName reports whether s matches ^[a-z0-9][a-z0-9_-]{0,31}$, the rule for
// account and group names. Load checks every name on each launch, so this
// is a loop rather than a regexp.
func isName(s string) bool {
	if len(s) == 0 || len(s) > 32 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || i > 0 && (c == '_' || c == '-')) {
			return false
		}
	}
	return true
}

// ValidateGroup checks a share group name; groups are folder names.
func ValidateGroup(group string) error {
	if !isName(group) {
		return fmt.Errorf("invalid group name: %q", group)
	}
	return nil
}

func ValidateName(name string) error {
	if !isName(name) {
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
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	// Names and groups become folder names; a hand edit such as "../x"
	// must not point ccam outside its own folders.
	for i := range c.Profiles {
		p := &c.Profiles[i]
		if !isName(p.Name) {
			return nil, fmt.Errorf("%s: invalid account name %q", path, p.Name)
		}
		for _, item := range Items {
			if g := p.Share.Get(item); g != Local && !isName(g) {
				return nil, fmt.Errorf("%s: account %s has an invalid %s group %q", path, p.Name, item, g)
			}
		}
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

// FindShortcut returns the profile whose shortcut command is name.
func (c *Config) FindShortcut(name string) (*Profile, bool) {
	if name == "" {
		return nil, false
	}
	for i := range c.Profiles {
		if c.Profiles[i].Shortcut == name {
			return &c.Profiles[i], true
		}
	}
	return nil, false
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
	c.Profiles = slices.DeleteFunc(c.Profiles, func(p Profile) bool { return p.Name == name })
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
