// Package mcp shares MCP servers through a group's mcp.json.
//
// Two strategies (docs/claude-layout.md):
//   - flag: pass `--mcp-config <group>/mcp.json` at launch; .claude.json is never touched.
//   - sync: merge the group servers into the profile's .claude.json before launch,
//     under Claude Code's own lock, with a backup.
package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/trace"
)

// Store is a group's mcp.json, in the format --mcp-config reads.
type Store struct {
	Servers map[string]json.RawMessage `json:"mcpServers"`
}

func Path(groupDir string) string { return filepath.Join(groupDir, "mcp.json") }

func Load(path string) (*Store, error) {
	s := &Store{Servers: map[string]json.RawMessage{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("%s 읽기 실패: %w", path, err)
	}
	if s.Servers == nil {
		s.Servers = map[string]json.RawMessage{}
	}
	return s, nil
}

// Save writes with 0600: server entries often carry API keys.
func (s *Store) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(path, append(data, '\n'), 0o600)
}

func (s *Store) Names() []string { return sortedKeys(s.Servers) }

func sortedKeys(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Validate checks that a server definition is a JSON object with a command or url.
func Validate(def json.RawMessage) error {
	var v struct {
		Command string `json:"command"`
		URL     string `json:"url"`
	}
	if err := json.Unmarshal(def, &v); err != nil {
		return fmt.Errorf("MCP 서버 정의는 JSON 객체여야 합니다: %w", err)
	}
	if v.Command == "" && v.URL == "" {
		return fmt.Errorf("MCP 서버 정의에 command 나 url 이 필요합니다")
	}
	return nil
}

// LocalServers returns the user-scope servers stored in a profile's .claude.json.
func LocalServers(profileDir string) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(filepath.Join(profileDir, ".claude.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	trace.Step("read .claude.json")
	raw, ok, err := profile.TopLevelValue(data, "mcpServers")
	trace.Step("scan .claude.json")
	if err != nil || !ok {
		return nil, err
	}
	var servers map[string]json.RawMessage
	return servers, json.Unmarshal(raw, &servers)
}

// FlagArgs returns the launch args for the flag strategy, or nil when the
// group has no mcp.json yet. One stat, no parsing.
func FlagArgs(groupDir string) []string {
	path := Path(groupDir)
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	return []string{"--mcp-config", path}
}

// NewLocal reports servers that appeared in the profile's .claude.json
// (e.g. via `claude mcp add` inside a session) and are not in the group.
// It only parses when .claude.json changed since the last call.
func NewLocal(profileDir, groupMCPPath, seenPath string) ([]string, error) {
	info, err := os.Stat(filepath.Join(profileDir, ".claude.json"))
	if err != nil {
		return nil, nil
	}
	stamp := strconv.FormatInt(info.ModTime().UnixNano(), 10)
	seen := map[string]bool{}
	if data, err := os.ReadFile(seenPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if lines[0] == stamp {
			return nil, nil
		}
		for _, n := range lines[1:] {
			seen[n] = true
		}
	}
	local, err := LocalServers(profileDir)
	if err != nil {
		return nil, err
	}
	group, err := Load(groupMCPPath)
	if err != nil {
		return nil, err
	}
	names := sortedKeys(local)
	var fresh []string
	for _, n := range names {
		if _, shared := group.Servers[n]; !seen[n] && !shared {
			fresh = append(fresh, n)
		}
	}
	_ = fsutil.WriteFast(seenPath, []byte(stamp+"\n"+strings.Join(names, "\n")))
	return fresh, nil
}

// Sync merges the group's servers into <profileDir>/.claude.json.
// stateDir keeps the sync stamp, the names ccam manages, and backups.
// The caller must make sure no Claude Code session runs in this profile.
func Sync(profileDir, groupDir, stateDir string) error {
	mcpInfo, err := os.Stat(Path(groupDir))
	if err != nil {
		return nil // nothing shared yet
	}
	stampPath := filepath.Join(stateDir, "sync-stamp")
	stamp := strconv.FormatInt(mcpInfo.ModTime().UnixNano(), 10)
	var managed []string
	if data, err := os.ReadFile(stampPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if lines[0] == stamp {
			return nil
		}
		managed = lines[1:]
	}
	group, err := Load(Path(groupDir))
	if err != nil {
		return err
	}

	claudeJSON := filepath.Join(profileDir, ".claude.json")
	unlock, err := lock.DirLock(claudeJSON+".lock", 2*time.Second)
	if err != nil {
		return err
	}
	defer unlock()

	doc := map[string]json.RawMessage{}
	perm := fs.FileMode(0o600)
	if data, err := os.ReadFile(claudeJSON); err == nil {
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf(".claude.json 읽기 실패 (손대지 않음): %w", err)
		}
		if err := backup(claudeJSON, data, filepath.Join(stateDir, "backups")); err != nil {
			return err
		}
		if info, err := os.Stat(claudeJSON); err == nil {
			perm = info.Mode().Perm()
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	servers := map[string]json.RawMessage{}
	if raw, ok := doc["mcpServers"]; ok {
		if err := json.Unmarshal(raw, &servers); err != nil {
			return fmt.Errorf(".claude.json 의 mcpServers 를 읽지 못했습니다: %w", err)
		}
	}
	for _, n := range managed {
		if _, still := group.Servers[n]; !still {
			delete(servers, n)
		}
	}
	for n, def := range group.Servers {
		servers[n] = def
	}
	raw, err := json.Marshal(servers)
	if err != nil {
		return err
	}
	doc["mcpServers"] = raw
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	if err := fsutil.AtomicWrite(claudeJSON, out, perm); err != nil {
		return err
	}
	return fsutil.WriteFast(stampPath, []byte(stamp+"\n"+strings.Join(group.Names(), "\n")))
}

// RemoveLocal deletes servers from a profile's .claude.json, with a backup,
// under Claude Code's lock. The caller ensures the profile is not running.
func RemoveLocal(profileDir, stateDir string, names []string) error {
	claudeJSON := filepath.Join(profileDir, ".claude.json")
	unlock, err := lock.DirLock(claudeJSON+".lock", 2*time.Second)
	if err != nil {
		return err
	}
	defer unlock()
	data, err := os.ReadFile(claudeJSON)
	if err != nil {
		return err
	}
	info, err := os.Stat(claudeJSON)
	if err != nil {
		return err
	}
	doc := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	servers := map[string]json.RawMessage{}
	if raw, ok := doc["mcpServers"]; ok {
		if err := json.Unmarshal(raw, &servers); err != nil {
			return err
		}
	}
	for _, n := range names {
		delete(servers, n)
	}
	raw, _ := json.Marshal(servers)
	doc["mcpServers"] = raw
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	if err := backup(claudeJSON, data, filepath.Join(stateDir, "backups")); err != nil {
		return err
	}
	return fsutil.AtomicWrite(claudeJSON, out, info.Mode().Perm())
}

const keepBackups = 5

func backup(src string, data []byte, dir string) error {
	name := filepath.Base(src) + "." + time.Now().Format("20060102-150405.000000000")
	if err := fsutil.AtomicWrite(filepath.Join(dir, name), data, 0o600); err != nil {
		return err
	}
	old, _ := filepath.Glob(filepath.Join(dir, filepath.Base(src)+".*"))
	sort.Strings(old)
	for len(old) > keepBackups {
		os.Remove(old[0])
		old = old[1:]
	}
	return nil
}
