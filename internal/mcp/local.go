package mcp

import (
	"encoding/json"
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/trace"
)

// FolderServers returns the local-scope servers of a profile's .claude.json
// by folder: projects[<folder>].mcpServers, where `claude mcp add` puts a
// server unless given -s user. Such a server loads for this account only,
// and only in that folder (for a git repository: the top of its main
// worktree, which every folder of the repository and its worktrees uses,
// see ProjectKey).
func FolderServers(profileDir string) (map[string]map[string]json.RawMessage, error) {
	data, err := os.ReadFile(filepath.Join(profileDir, ".claude.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return folderServers(data)
}

func folderServers(data []byte) (map[string]map[string]json.RawMessage, error) {
	raw, ok, err := profile.TopLevelValue(data, "projects")
	if err != nil || !ok {
		return nil, err
	}
	return projectServers(raw)
}

// projectServers is folderServers for the value of projects.
func projectServers(raw []byte) (map[string]map[string]json.RawMessage, error) {
	projects, err := entries(raw)
	if err != nil {
		return nil, err
	}
	out := map[string]map[string]json.RawMessage{}
	for folder, entry := range projects {
		// The entries can be large (older versions kept prompts there): only
		// their mcpServers is decoded.
		raw, ok, err := profile.TopLevelValue(entry, "mcpServers")
		if err != nil || !ok {
			continue
		}
		var servers map[string]json.RawMessage
		if json.Unmarshal(raw, &servers) == nil && len(servers) > 0 {
			out[folder] = servers
		}
	}
	return out, nil
}

// projectEntries splits .claude.json's projects into folder → entry,
// without decoding the entries.
func projectEntries(data []byte) (map[string][]byte, error) {
	raw, ok, err := profile.TopLevelValue(data, "projects")
	if err != nil || !ok {
		return nil, err
	}
	return entries(raw)
}

// entries splits a JSON object into key → raw value.
func entries(raw []byte) (map[string][]byte, error) {
	out := map[string][]byte{}
	return out, profile.EachValue(raw, func(key string, value []byte) error {
		out[key] = value
		return nil
	})
}

// scopes reads both scopes of servers from .claude.json in one pass: the
// user scope and, by folder, the local scope.
func scopes(data []byte) (user map[string]json.RawMessage, folders map[string]map[string]json.RawMessage, err error) {
	var users, projects []byte
	err = profile.EachValue(data, func(key string, value []byte) error {
		switch key {
		case "mcpServers":
			users = value
		case "projects":
			projects = value
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	if users != nil {
		if err := json.Unmarshal(users, &user); err != nil {
			return nil, nil, err
		}
	}
	if projects != nil {
		if folders, err = projectServers(projects); err != nil {
			return nil, nil, err
		}
	}
	return user, folders, nil
}

// ProjectKey is the folder Claude Code files the local-scope servers of
// dir under (2.1.308): dir with symbolic links resolved, or inside a git
// repository the top folder of its main worktree (for a worktree of a bare
// repository, the repository's folder).
func ProjectKey(dir string) string {
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		dir = real
	}
	top, ok := gitFolder(dir)
	if !ok {
		return dir
	}
	if main, ok := mainWorktree(top); ok {
		return main
	}
	return top
}

// gitFolder is the nearest folder from dir up that holds a .git.
func gitFolder(dir string) (string, bool) {
	for d := dir; ; {
		if _, err := os.Lstat(filepath.Join(d, ".git")); err == nil {
			return d, true
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}

// mainWorktree is the top folder of the main worktree when top is another
// worktree of a repository (its .git is a file naming the worktree's git
// folder, whose commondir names the main .git), or the folder of a bare
// repository top is a worktree of.
func mainWorktree(top string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(top, ".git"))
	if err != nil {
		return "", false // a folder, the main worktree itself
	}
	gitdir, ok := strings.CutPrefix(strings.TrimSpace(string(data)), "gitdir:")
	if !ok {
		return "", false
	}
	gitdir = filepath.Clean(abs(top, strings.TrimSpace(gitdir)))
	common, err := os.ReadFile(filepath.Join(gitdir, "commondir"))
	if err != nil {
		return "", false
	}
	c := filepath.Clean(abs(gitdir, strings.TrimSpace(string(common))))
	if filepath.Base(c) != ".git" {
		return c, true // a bare repository: Claude Code takes its folder (2.1.308)
	}
	return filepath.Dir(c), true
}

func abs(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

// sameFolder compares a folder key of .claude.json with a path. Keys are
// Claude Code's spelling; on Windows case and slashes do not matter.
func sameFolder(key, path string) bool {
	k := filepath.Clean(filepath.FromSlash(key))
	if runtime.GOOS == "windows" {
		return strings.EqualFold(k, path)
	}
	return k == path
}

// FolderLookup finds the entry of folders whose key is the folder dir.
func FolderLookup[V any](folders map[string]V, dir string) (V, bool) {
	if v, ok := folders[dir]; ok {
		return v, true
	}
	for k, v := range folders {
		if sameFolder(k, dir) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// Trusts reports whether the account of profileDir accepted Claude Code's
// trust dialog for the folder dir (projects[<folder>].hasTrustDialogAccepted,
// per account): for dir's ProjectKey (in another worktree, the main one's
// folder), for dir itself or, inside a git repository, for a folder between
// dir and the top of the repository. A folder above those does not count
// (2.1.308).
func Trusts(profileDir, dir string) bool {
	data, err := os.ReadFile(filepath.Join(profileDir, ".claude.json"))
	if err != nil {
		return false
	}
	projects, err := projectEntries(data)
	if err != nil || len(projects) == 0 {
		return false
	}
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		dir = real
	}
	trusted := func(folder string) bool {
		entry, ok := FolderLookup(projects, folder)
		if !ok {
			return false
		}
		v, ok, _ := profile.TopLevelValue(entry, "hasTrustDialogAccepted")
		return ok && string(v) == "true"
	}
	if trusted(ProjectKey(dir)) {
		return true
	}
	top, inGit := gitFolder(dir)
	for d := dir; ; {
		if trusted(d) {
			return true
		}
		parent := filepath.Dir(d)
		if !inGit || d == top || parent == d {
			return false
		}
		d = parent
	}
}

// Local is what NewLocal found that the account has alone.
type Local struct {
	User   []string            // user-scope servers the group lacks
	Folder map[string][]string // local-scope servers (claude mcp add's default), by folder
	// Shadow lists user-scope servers with a group server's name and
	// another definition: in a session of the account they win over the
	// group's, which --mcp-config passes.
	Shadow []string
}

// NewLocal reports servers that appeared in the profile's .claude.json
// (e.g. via `claude mcp add` inside a session) that the account has alone,
// each once. It only parses when .claude.json changed since the last call.
func NewLocal(profileDir, groupMCPPath, seenPath string) (Local, error) {
	var l Local
	path := filepath.Join(profileDir, ".claude.json")
	info, err := os.Stat(path)
	if err != nil {
		return l, nil
	}
	stamp := strconv.FormatInt(info.ModTime().UnixNano(), 10)
	seen := map[string]bool{}
	if data, err := os.ReadFile(seenPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if lines[0] == stamp {
			return l, nil
		}
		for _, n := range lines[1:] {
			seen[n] = true
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return l, err
	}
	trace.Step("read .claude.json")
	user, folders, err := scopes(data)
	if err != nil {
		return l, err
	}
	trace.Step("scan .claude.json")
	group, err := Load(groupMCPPath)
	if err != nil {
		return l, err
	}
	// Seen are the names (user scope), "<name>\t<folder>" (local scope) and
	// "<name>\t!" (a user-scope server with a group server's name).
	var now []string
	for _, n := range sortedKeys(user) {
		now = append(now, n)
		def, shared := group.Servers[n]
		switch {
		case !shared:
			if !seen[n] {
				l.User = append(l.User, n)
			}
		case !SameDef(def, user[n]):
			key := n + "\t!"
			now = append(now, key)
			if !seen[key] {
				l.Shadow = append(l.Shadow, n)
			}
		}
	}
	for _, folder := range slices.Sorted(maps.Keys(folders)) {
		for _, n := range sortedKeys(folders[folder]) {
			key := n + "\t" + folder
			now = append(now, key)
			if !seen[key] {
				if l.Folder == nil {
					l.Folder = map[string][]string{}
				}
				l.Folder[folder] = append(l.Folder[folder], n)
			}
		}
	}
	_ = fsutil.WriteFast(seenPath, []byte(stamp+"\n"+strings.Join(now, "\n")))
	return l, nil
}

// SignedIn lists the MCP servers the account of profileDir is signed in to
// (a stored token in its login, mcpOAuth), by the names Claude Code keys
// them with; a plugin's is plugin:<plugin>:<server>. A sign-in works for
// its account only, so ccam never copies one: the other accounts sign in
// themselves, with /mcp.
func SignedIn(profileDir, home string) map[string]bool {
	data, _, err := creds.Read(profileDir, home)
	if err != nil {
		return nil
	}
	out := map[string]bool{}
	for k := range creds.MCPSignIns(data) {
		// "<server>|<16 hex digits of a hash of its type, url and headers>"
		if i := strings.LastIndexByte(k, '|'); i > 0 {
			out[k[:i]] = true
		}
	}
	return out
}
