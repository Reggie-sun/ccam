// Package mcp shares MCP servers through a group's mcp.json.
//
// Two strategies (docs/claude-layout.md):
//   - flag: pass `--mcp-config <group>/mcp.json` at launch; .claude.json is never touched.
//   - sync: merge the group servers into the profile's .claude.json before launch,
//     under Claude Code's own lock, with a backup.
package mcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	if s.Servers == nil {
		s.Servers = map[string]json.RawMessage{}
	}
	return s, nil
}

// Save writes with 0600: server entries often carry API keys.
func (s *Store) Save(path string) error {
	data, err := fsutil.MarshalJSON(s)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(path, data, 0o600)
}

func (s *Store) Names() []string { return sortedKeys(s.Servers) }

func sortedKeys(m map[string]json.RawMessage) []string { return slices.Sorted(maps.Keys(m)) }

// SameDef compares two server definitions by meaning: key order and
// spacing do not matter.
func SameDef(a, b json.RawMessage) bool {
	var x, y any
	if json.Unmarshal(a, &x) != nil || json.Unmarshal(b, &y) != nil {
		return bytes.Equal(a, b)
	}
	return reflect.DeepEqual(x, y)
}

// Validate checks that a server definition is a JSON object with a command or url.
func Validate(def json.RawMessage) error {
	var v struct {
		Command string `json:"command"`
		URL     string `json:"url"`
	}
	if err := json.Unmarshal(def, &v); err != nil {
		return fmt.Errorf("an MCP server definition must be a JSON object: %w", err)
	}
	if v.Command == "" && v.URL == "" {
		return fmt.Errorf("an MCP server definition needs a command or url")
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
	// One argument: claude's --mcp-config is variadic, so "--mcp-config X"
	// would also eat a following prompt or subcommand as a config file.
	return []string{"--mcp-config=" + path}
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

func syncStampPath(stateDir string) string { return filepath.Join(stateDir, "sync-stamp") }

// syncStamp identifies a version of a group's mcp.json; the folder is in
// it so that moving to another group syncs again.
func syncStamp(groupDir string, mcpInfo os.FileInfo) string {
	return strconv.FormatInt(mcpInfo.ModTime().UnixNano(), 10) + " " + groupDir
}

// Synced reports whether the profile's .claude.json holds the group's current
// servers as the last Sync wrote them (true when the group has none).
func Synced(groupDir, stateDir string) bool {
	info, err := os.Stat(Path(groupDir))
	if err != nil {
		return true
	}
	data, err := os.ReadFile(syncStampPath(stateDir))
	return err == nil && strings.SplitN(string(data), "\n", 2)[0] == syncStamp(groupDir, info)
}

// HasCopy reports whether Sync has ever copied group servers into the profile.
func HasCopy(stateDir string) bool {
	_, err := os.Stat(syncStampPath(stateDir))
	return err == nil
}

// Copied lists the servers the last Sync copied into the profile.
func Copied(stateDir string) []string {
	data, err := os.ReadFile(syncStampPath(stateDir))
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")[1:]
}

// Unsync removes the servers Sync copied, for an account that stops sharing
// MCP: left in .claude.json they would still load, and under the flag
// strategy win over the servers of a group it shares MCP with later.
// The caller makes sure the profile is not running.
func Unsync(profileDir, stateDir string) error {
	if !HasCopy(stateDir) {
		return nil
	}
	if names := Copied(stateDir); len(names) > 0 && fsutil.Exists(filepath.Join(profileDir, ".claude.json")) {
		if err := RemoveLocal(profileDir, stateDir, names); err != nil {
			return err
		}
	}
	return os.Remove(syncStampPath(stateDir))
}

// Sync merges the group's servers into <profileDir>/.claude.json.
// stateDir keeps the sync stamp, the names ccam manages, and backups.
// The caller must make sure no Claude Code session runs in this profile.
func Sync(profileDir, groupDir, stateDir string) error {
	mcpInfo, err := os.Stat(Path(groupDir))
	if err != nil {
		return nil // nothing shared yet
	}
	stampPath := syncStampPath(stateDir)
	stamp := syncStamp(groupDir, mcpInfo)
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
			return fmt.Errorf("could not read .claude.json (left untouched): %w", err)
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
			return fmt.Errorf("could not read mcpServers in .claude.json: %w", err)
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
	// Not json.Marshal: it would write && and > in commands as \u0026 and \u003e.
	raw, err := fsutil.MarshalJSON(servers)
	if err != nil {
		return err
	}
	doc["mcpServers"] = raw
	out, err := fsutil.MarshalJSON(doc)
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
	raw, err := fsutil.MarshalJSON(servers)
	if err != nil {
		return err
	}
	doc["mcpServers"] = raw
	out, err := fsutil.MarshalJSON(doc)
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
