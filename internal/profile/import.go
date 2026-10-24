package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// Candidate is an existing Claude Code config folder that could be imported.
type Candidate struct {
	Path  string
	Name  string // suggested profile name ("" if none could be derived)
	Email string
}

// FindCandidates lists ~/.claude and ~/.claude-* folders not yet registered.
func FindCandidates(p paths.Paths, cfg *config.Config) []Candidate {
	matches, _ := filepath.Glob(filepath.Join(p.Home, ".claude-*"))
	all := append([]string{filepath.Join(p.Home, ".claude")}, matches...)
	var out []Candidate
	for _, path := range all {
		if info, err := os.Stat(path); err != nil || !info.IsDir() || registered(p, cfg, path) {
			continue
		}
		email, _ := ReadEmail(path)
		if email == "" && path == filepath.Join(p.Home, ".claude") {
			email, _ = ReadEmail(p.Home) // default layout keeps .claude.json in $HOME
		}
		name := SuggestName(path)
		if _, taken := cfg.Find(name); taken {
			name = "" // ~/.claude-work next to an account work: the folder needs a name of its own
		}
		out = append(out, Candidate{Path: path, Name: name, Email: email})
	}
	return out
}

// PlainClaudeName names plain claude's folder where an account name would go.
const PlainClaudeName = "plain claude (~/.claude)"

// PlainClaude returns ~/.claude, the folder plain claude uses, when no
// account is registered there, or "".
func PlainClaude(p paths.Paths, cfg *config.Config) string {
	dir := filepath.Join(p.Home, ".claude")
	if registered(p, cfg, dir) {
		return ""
	}
	return dir
}

func registered(p paths.Paths, cfg *config.Config, path string) bool {
	for i := range cfg.Profiles {
		if Dir(p, &cfg.Profiles[i]) == path {
			return true
		}
	}
	return false
}

// SuggestName derives a profile name: ~/.claude → main, ~/.claude-work →
// work, ~/.claude-2 → acct2. Returns "" when nothing valid comes out.
func SuggestName(path string) string {
	base := filepath.Base(path)
	name := "main"
	if base != ".claude" {
		name = strings.ToLower(strings.TrimLeft(strings.TrimPrefix(base, ".claude"), "-_."))
		if _, err := strconv.Atoi(name); err == nil {
			name = "acct" + name
		}
	}
	if config.ValidateName(name) != nil {
		return ""
	}
	return name
}

// Import registers an existing folder in place (it is not moved or copied:
// copying would duplicate the refresh token, and the copies would log each
// other out). Notes describe anything that was changed on disk.
func Import(p paths.Paths, cfg *config.Config, path, name string) (*config.Profile, []string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		return nil, nil, fmt.Errorf("not a folder: %s", path)
	}
	if registered(p, cfg, path) {
		return nil, nil, fmt.Errorf("folder is already registered: %s", path)
	}
	if name == "" {
		name = SuggestName(path)
	}
	if name == "" {
		return nil, nil, fmt.Errorf("cannot derive a name. Give one with --name")
	}
	if err := config.ValidateName(name); err != nil {
		return nil, nil, err
	}
	if _, ok := cfg.Find(name); ok {
		return nil, nil, fmt.Errorf("an account named %q already exists (give another name with --name)", name)
	}

	var notes []string
	// Without CLAUDE_CONFIG_DIR, Claude Code keeps .claude.json in $HOME,
	// next to ~/.claude. Link it in so `ccam <name>` and plain `claude` see
	// the same account and settings (Claude Code writes through the link).
	inner := filepath.Join(path, ".claude.json")
	outer := filepath.Join(p.Home, ".claude.json")
	switch {
	case path == filepath.Join(p.Home, ".claude") && fsutil.Exists(outer) && !fsutil.Exists(inner):
		if err := fsutil.Link(outer, inner, false); err != nil {
			return nil, nil, err
		}
		notes = append(notes, fmt.Sprintf("created link %s → %s", inner, outer))
	case OwnClaudeJSON(p, path):
		// Left from a run with CLAUDE_CONFIG_DIR=~/.claude, say; which of
		// the two is the right one only the user knows.
		notes = append(notes, fmt.Sprintf("%s has a .claude.json of its own, which `ccam %s` uses, while plain `claude` uses %s: "+
			"their login details, MCP servers and project settings are kept apart. To keep one, replace %s with a link to %s",
			path, name, outer, inner, outer))
	}

	email, _ := ReadEmail(path)
	// Recorded as `ccam login` records it, so that a later login into
	// another account inside a session is noticed.
	id, _ := ReadIdentity(path)
	cfg.DropMember(name)
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: cfg.NextSlot(), Email: email, Dir: path, Login: id})
	if len(cfg.Profiles) == 1 {
		cfg.Default = name
	}
	return &cfg.Profiles[len(cfg.Profiles)-1], notes, nil
}

// OwnClaudeJSON reports whether dir, when it is ~/.claude, has a
// .claude.json of its own beside ~/.claude.json, which plain claude uses,
// instead of the link import makes to it.
func OwnClaudeJSON(p paths.Paths, dir string) bool {
	inner, outer := filepath.Join(dir, ".claude.json"), filepath.Join(p.Home, ".claude.json")
	return dir == filepath.Join(p.Home, ".claude") && fsutil.Exists(inner) && fsutil.Exists(outer) && !sameFile(inner, outer)
}

// sameFile reports whether a and b are one file (a link to the other).
func sameFile(a, b string) bool {
	ia, err := os.Stat(a)
	if err != nil {
		return false
	}
	ib, err := os.Stat(b)
	return err == nil && os.SameFile(ia, ib)
}
