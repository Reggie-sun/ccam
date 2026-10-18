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
		out = append(out, Candidate{Path: path, Name: SuggestName(path), Email: email})
	}
	return out
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
		return nil, nil, fmt.Errorf("폴더가 아닙니다: %s", path)
	}
	if registered(p, cfg, path) {
		return nil, nil, fmt.Errorf("이미 등록된 폴더입니다: %s", path)
	}
	if name == "" {
		name = SuggestName(path)
	}
	if name == "" {
		return nil, nil, fmt.Errorf("이름을 정할 수 없습니다. --name 으로 지정하세요")
	}
	if err := config.ValidateName(name); err != nil {
		return nil, nil, err
	}
	if _, ok := cfg.Find(name); ok {
		return nil, nil, fmt.Errorf("이미 있는 계정 이름입니다: %q (--name 으로 다른 이름을 주세요)", name)
	}

	var notes []string
	// Without CLAUDE_CONFIG_DIR, Claude Code keeps .claude.json in $HOME,
	// next to ~/.claude. Link it in so `ccs <name>` and plain `claude` see
	// the same account and settings (Claude Code writes through the link).
	inner := filepath.Join(path, ".claude.json")
	outer := filepath.Join(p.Home, ".claude.json")
	if path == filepath.Join(p.Home, ".claude") && !fsutil.Exists(inner) && fsutil.Exists(outer) {
		if err := os.Symlink(outer, inner); err != nil {
			return nil, nil, err
		}
		notes = append(notes, fmt.Sprintf("%s → %s 링크를 만들었습니다", inner, outer))
	}

	email, _ := ReadEmail(path)
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: cfg.NextSlot(), Email: email, Dir: path})
	if len(cfg.Profiles) == 1 {
		cfg.Default = name
	}
	return &cfg.Profiles[len(cfg.Profiles)-1], notes, nil
}
