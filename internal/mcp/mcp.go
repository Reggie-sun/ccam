// Package mcp shares MCP servers through a group's mcp.json.
//
// Two strategies (docs/claude-layout.md):
//   - flag: pass `--mcp-config <group>/mcp.json` at launch; .claude.json is never touched.
//   - sync: merge the group servers into the profile's .claude.json before launch,
//     under Claude Code's own lock, with a backup.
//
// Whatever ccam writes into an account's .claude.json it notes with a hash
// of each definition (the sync stamp, or the kept copies). It only ever
// replaces or deletes an entry it wrote itself and nobody changed since;
// every other entry is the account's own and stays.
package mcp

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

// lockName is the file ccam locks beside a group's mcp.json while it
// changes it. Claude Code never writes that mcp.json, only ccam does.
const lockName = ".mcp.json.lock"

// Lock takes ccam's lock for the group's mcp.json at path, waiting while
// another ccam command changes it.
func Lock(path string) (unlock func(), err error) {
	l, err := lock.WaitFlock(filepath.Join(filepath.Dir(path), lockName), 10*time.Second)
	if err != nil {
		return nil, err
	}
	return l.Unlock, nil
}

// Update loads the group's mcp.json at path, lets change edit it and saves
// it, under Lock: two ccam commands that change one group's servers at the
// same moment both keep their change. Nothing is saved when change fails.
func Update(path string, change func(*Store) error) error {
	unlock, err := Lock(path)
	if err != nil {
		return err
	}
	defer unlock()
	s, err := Load(path)
	if err != nil {
		return err
	}
	if err := change(s); err != nil {
		return err
	}
	return s.Save(path)
}

// Seed copies the group's mcp.json src to dst unless dst exists (under the
// lock for dst): the first servers of a group an account moves to.
func Seed(src, dst string) error {
	unlock, err := Lock(dst)
	if err != nil {
		return err
	}
	defer unlock()
	if fsutil.Exists(dst) {
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(dst, data, 0o600)
}

// SameDef compares two server definitions by meaning: key order and
// spacing do not matter.
func SameDef(a, b json.RawMessage) bool {
	var x, y any
	if json.Unmarshal(a, &x) != nil || json.Unmarshal(b, &y) != nil {
		return bytes.Equal(a, b)
	}
	return reflect.DeepEqual(x, y)
}

// defHash names a server definition by its meaning, as SameDef compares
// them, without keeping what it holds (API keys) in ccam's notes.
func defHash(def json.RawMessage) string {
	var v any
	if json.Unmarshal(def, &v) == nil {
		if b, err := json.Marshal(v); err == nil { // object keys come out sorted
			def = b
		}
	}
	sum := sha256.Sum256(def)
	return hex.EncodeToString(sum[:8])
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
	servers, err := userServers(data)
	trace.Step("scan .claude.json")
	return servers, err
}

func userServers(data []byte) (map[string]json.RawMessage, error) {
	raw, ok, err := profile.TopLevelValue(data, "mcpServers")
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

// stamp is what ccam wrote into an account's .claude.json: the sync stamp
// (stateDir/sync-stamp) of the sync strategy and of the Remote Control
// copy, or the copies an account kept when it stopped sharing MCP
// (stateDir/kept). A file of lines:
//
//	<mcp.json mtime> <group folder>   the version of the group's mcp.json
//	<name>\t<hash>                    a server ccam wrote, as it wrote it
//	<name>\t!<hash>                   the account's own server of a group server's
//	                                  name, left alone; <hash> is the group's
//	                                  definition, which was said to be unused there
//	<name>                            a copy an older ccam listed without its hash
type stamp struct {
	head   string
	copies map[string]string // name → hash of what ccam wrote, "" when not known
	own    map[string]string // name → hash of the group's definition not used
}

func syncStampPath(stateDir string) string { return filepath.Join(stateDir, "sync-stamp") }

func keptPath(stateDir string) string { return filepath.Join(stateDir, "kept") }

// checkedPath notes the group's mcp.json (the stamp's head) and the
// profile's .claude.json (time and size) as Sync last left them complete:
// while both are unchanged, SyncFresh tells without reading either.
func checkedPath(stateDir string) string { return filepath.Join(stateDir, "sync-checked") }

func checkedNote(head, profileDir string) string {
	info, err := os.Stat(filepath.Join(profileDir, ".claude.json"))
	if err != nil {
		return ""
	}
	return head + "\n" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + " " + strconv.FormatInt(info.Size(), 10) + "\n"
}

// noteChecked records that Sync found .claude.json complete for head;
// losing it only costs the next launch a read.
func noteChecked(head, profileDir, stateDir string) {
	if note := checkedNote(head, profileDir); note != "" {
		_ = fsutil.WriteFast(checkedPath(stateDir), []byte(note))
	}
}

// SyncFresh reports, from two file times and one small file, that Sync
// has nothing to do: neither the group's mcp.json nor the profile's
// .claude.json changed since Sync last left them complete. A launch then
// skips Sync and the running check before it.
func SyncFresh(profileDir, groupDir, stateDir string) bool {
	info, err := os.Stat(Path(groupDir))
	if err != nil {
		return true // nothing shared yet: Sync does nothing
	}
	data, err := os.ReadFile(checkedPath(stateDir))
	note := checkedNote(syncStamp(groupDir, info), profileDir)
	return err == nil && note != "" && string(data) == note
}

// syncStamp identifies a version of a group's mcp.json; the folder is in
// it so that moving to another group syncs again.
func syncStamp(groupDir string, mcpInfo os.FileInfo) string {
	return strconv.FormatInt(mcpInfo.ModTime().UnixNano(), 10) + " " + groupDir
}

// readStamp reads a stamp; without the file it is empty and ok is false.
func readStamp(path string) (st *stamp, ok bool) {
	st = &stamp{copies: map[string]string{}, own: map[string]string{}}
	data, err := os.ReadFile(path)
	if err != nil {
		return st, false
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	st.head = lines[0]
	for _, l := range lines[1:] {
		i := strings.LastIndexByte(l, '\t')
		switch {
		case l == "":
		case i < 0:
			st.copies[l] = ""
		case strings.HasPrefix(l[i+1:], "!"):
			st.own[l[:i]] = l[i+2:]
		default:
			st.copies[l[:i]] = l[i+1:]
		}
	}
	return st, true
}

func (st *stamp) save(path string) error {
	lines := []string{st.head}
	for _, n := range slices.Sorted(maps.Keys(st.copies)) {
		lines = append(lines, n+"\t"+st.copies[n])
	}
	for _, n := range slices.Sorted(maps.Keys(st.own)) {
		lines = append(lines, n+"\t!"+st.own[n])
	}
	return fsutil.AtomicWrite(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600)
}

// legacy reports whether the stamp lists a copy without its hash.
func (st *stamp) legacy() bool {
	for _, h := range st.copies {
		if h == "" {
			return true
		}
	}
	return false
}

// group loads the group's mcp.json the stamp's copies came from; nil when
// that cannot be read.
func (st *stamp) group() *Store {
	_, dir, ok := strings.Cut(st.head, " ")
	if !ok || dir == "" {
		return nil
	}
	s, err := Load(Path(dir))
	if err != nil {
		return nil
	}
	return s
}

// ours reports whether cur, the account's entry name, is a copy ccam wrote
// that nobody changed since. A copy listed without its hash counts only
// while it is the same as the group's definition.
func (st *stamp) ours(name string, cur json.RawMessage, group *Store) bool {
	h, ok := st.copies[name]
	switch {
	case !ok:
		return false
	case h != "":
		return defHash(cur) == h
	case group == nil:
		return false
	}
	def, ok := group.Servers[name]
	return ok && SameDef(def, cur)
}

// missing reports whether a server of group is not in the profile's
// .claude.json, which Sync then puts back. A .claude.json ccam cannot read
// counts as complete: Sync could not change it either.
func missing(profileDir string, group *Store) bool {
	local, err := LocalServers(profileDir)
	if err != nil {
		return false
	}
	for n := range group.Servers {
		if _, ok := local[n]; !ok {
			return true
		}
	}
	return false
}

// Synced reports whether the profile's .claude.json holds the group's
// current servers as Sync copies them (true when the group has none):
// Sync ran for the group's mcp.json as it is, and none of its servers is
// gone from .claude.json since.
func Synced(profileDir, groupDir, stateDir string) bool {
	info, err := os.Stat(Path(groupDir))
	if err != nil {
		return true
	}
	st, ok := readStamp(syncStampPath(stateDir))
	if !ok || st.head != syncStamp(groupDir, info) {
		return false
	}
	group, err := Load(Path(groupDir))
	return err != nil || !missing(profileDir, group)
}

// NoteCopies gives the copies an older ccam listed in the sync stamp
// without their hashes those hashes, while the group's mcp.json is still
// the version it synced: a copy that equals the group's definition then is
// the one it wrote. Called before that file changes, so that the copies
// follow the change instead of counting as the account's own. It writes
// the stamp only, never .claude.json.
func NoteCopies(profileDir, groupDir, stateDir string) error {
	st, ok := readStamp(syncStampPath(stateDir))
	if !ok || !st.legacy() {
		return nil
	}
	info, err := os.Stat(Path(groupDir))
	if err != nil || st.head != syncStamp(groupDir, info) {
		return nil // changed since: a copy cannot be told from the account's own
	}
	group, err := Load(Path(groupDir))
	if err != nil {
		return err
	}
	local, err := LocalServers(profileDir)
	if err != nil {
		return err
	}
	for n, h := range st.copies {
		if cur, ok := local[n]; ok && h == "" && st.ours(n, cur, group) {
			st.copies[n] = defHash(cur)
		}
	}
	return st.save(syncStampPath(stateDir))
}

// HasCopy reports whether Sync has ever copied group servers into the profile.
func HasCopy(stateDir string) bool {
	_, err := os.Stat(syncStampPath(stateDir))
	return err == nil
}

// Copied splits the servers Sync copied into the profile's .claude.json
// into those still as it wrote them, which a move to another group or a
// --fresh leave removes, and those changed since, which are the account's
// own now and stay.
func Copied(profileDir, stateDir string) (drop, keep []string) {
	st, ok := readStamp(syncStampPath(stateDir))
	if !ok || len(st.copies) == 0 {
		return nil, nil
	}
	local, _ := LocalServers(profileDir)
	group := st.group()
	for _, n := range slices.Sorted(maps.Keys(st.copies)) {
		cur, ok := local[n]
		switch {
		case !ok:
		case st.ours(n, cur, group):
			drop = append(drop, n)
		default:
			keep = append(keep, n)
		}
	}
	return drop, keep
}

// Unsync removes the copies Sync made that are still as it made them, for
// an account that moves to another group or stops sharing MCP with
// --fresh: left in .claude.json they would still load, and under the flag
// strategy win over the new group's servers of the same name. Copies
// changed since stay. The caller makes sure the profile is not running.
func Unsync(profileDir, stateDir string) error {
	st, ok := readStamp(syncStampPath(stateDir))
	if !ok {
		return nil
	}
	if len(st.copies) > 0 {
		group := st.group()
		err := editServers(profileDir, stateDir, false, func(servers map[string]json.RawMessage) bool {
			changed := false
			for n := range st.copies {
				if cur, ok := servers[n]; ok && st.ours(n, cur, group) {
					delete(servers, n)
					changed = true
				}
			}
			return changed
		})
		if err != nil {
			return err
		}
	}
	os.Remove(checkedPath(stateDir))
	return removeIfThere(syncStampPath(stateDir))
}

// Result says what Sync did besides copying, for the user to hear.
type Result struct {
	// Own lists the account's own servers that have a group server's name
	// and another definition: they stay, and the group's is not used in
	// that account. Each is said once for each definition of the group's.
	Own []string
	// Restored lists copies that were gone from .claude.json (removed in a
	// session, or the file was made anew) and are back.
	Restored []string
	// Wrote says whether .claude.json was written.
	Wrote bool
}

// Sync copies the group's servers into <profileDir>/.claude.json, under
// Claude Code's lock and with a backup, and notes what it wrote in
// stateDir (the sync stamp). A server of the account's own under a group
// server's name stays as it is. Copies of servers the group no longer has
// go, unless they changed since. Unless the group's mcp.json changed or one
// of its servers is gone from .claude.json, Sync only reads.
// The caller must make sure no Claude Code session runs in this profile.
func Sync(profileDir, groupDir, stateDir string) (Result, error) {
	var res Result
	info, err := os.Stat(Path(groupDir))
	if err != nil {
		return res, nil // nothing shared yet
	}
	group, err := Load(Path(groupDir))
	if err != nil {
		return res, err
	}
	head := syncStamp(groupDir, info)
	old, _ := readStamp(syncStampPath(stateDir))
	// A stamp of an older ccam has no hashes: while the group's mcp.json is
	// still the one it synced, its copies are the group's definitions, and
	// a full pass notes their hashes.
	if old.head == head && !old.legacy() && !missing(profileDir, group) {
		noteChecked(head, profileDir, stateDir)
		return res, nil
	}
	next := &stamp{head: head, copies: map[string]string{}, own: map[string]string{}}
	err = editServers(profileDir, stateDir, true, func(servers map[string]json.RawMessage) bool {
		changed := false
		for n := range old.copies {
			if _, shared := group.Servers[n]; shared {
				continue
			}
			if cur, ok := servers[n]; ok && old.ours(n, cur, group) {
				delete(servers, n)
				changed = true
			}
		}
		for _, n := range group.Names() {
			def := group.Servers[n]
			cur, ok := servers[n]
			switch {
			case !ok:
				if _, was := old.copies[n]; was {
					res.Restored = append(res.Restored, n)
				}
			case old.ours(n, cur, group):
			case SameDef(cur, def):
				continue // the account's own, as the group has it: nothing to copy
			default:
				h := defHash(def)
				if old.own[n] != h {
					res.Own = append(res.Own, n)
				}
				next.own[n] = h
				continue
			}
			if !ok || !SameDef(cur, def) {
				servers[n] = def
				changed = true
			}
			next.copies[n] = defHash(def)
		}
		res.Wrote = changed
		return changed
	})
	if err != nil {
		return Result{}, err
	}
	if err := next.save(syncStampPath(stateDir)); err != nil {
		return res, err
	}
	noteChecked(head, profileDir, stateDir)
	return res, nil
}

// KeepPlan says what Keep does for an account that stops sharing the
// group in groupDir: copied lists the group's servers it gets in its
// .claude.json, own its servers of the same names with other definitions,
// which stay as they are.
func KeepPlan(profileDir, groupDir, stateDir string) (copied, own []string, err error) {
	group, err := Load(Path(groupDir))
	if err != nil {
		return nil, nil, err
	}
	if err := Readable(profileDir); err != nil {
		return nil, nil, err
	}
	local, err := LocalServers(profileDir)
	if err != nil {
		return nil, nil, err
	}
	st, _ := readStamp(syncStampPath(stateDir))
	prev, _ := readStamp(keptPath(stateDir)) // an earlier Keep's copies, which Keep replaces too
	for _, n := range group.Names() {
		cur, ok := local[n]
		switch {
		case !ok || st.ours(n, cur, group) || prev.ours(n, cur, nil):
			copied = append(copied, n)
		case !SameDef(cur, group.Servers[n]):
			own = append(own, n)
		}
	}
	return copied, own, nil
}

// Readable reports why ccam could not change the MCP servers in the
// profile's .claude.json: the file is there and does not read as a JSON
// object. A change of groups checks it before it saves the config, so that
// it does not stop at the account's MCP step.
func Readable(profileDir string) error {
	path := filepath.Join(profileDir, ".claude.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err == nil {
		var doc map[string]json.RawMessage
		if err = json.Unmarshal(data, &doc); err == nil {
			return nil
		}
	}
	return fmt.Errorf("could not read %s: %w (the change edits the MCP servers in it; repair the file first)", path, err)
}

// Keep gives an account that stops sharing MCP the servers it had from the
// group in groupDir, as entries of its own .claude.json (with a backup,
// under Claude Code's lock): those it lacks, and the copies Sync made as
// the group has them now. Its own servers of those names stay, and Sync's
// copies of servers the group no longer has go. Keep notes what it wrote
// (stateDir/kept), so that a group the account joins later gives those
// servers again in place of these copies (DropKept). The caller ensures
// the account is not running.
func Keep(profileDir, groupDir, stateDir string) error {
	group, err := Load(Path(groupDir))
	if err != nil {
		return err
	}
	st, _ := readStamp(syncStampPath(stateDir))
	// Copies an earlier Keep noted (a change that is finished again) stay
	// noted while unchanged.
	prev, _ := readStamp(keptPath(stateDir))
	kept := &stamp{copies: map[string]string{}, own: map[string]string{}}
	if info, err := os.Stat(Path(groupDir)); err == nil {
		kept.head = syncStamp(groupDir, info)
	}
	if local, err := LocalServers(profileDir); err == nil {
		for n, h := range prev.copies {
			if cur, ok := local[n]; ok && prev.ours(n, cur, nil) {
				kept.copies[n] = h
			}
		}
	}
	if len(group.Servers) > 0 || len(st.copies) > 0 {
		err := editServers(profileDir, stateDir, len(group.Servers) > 0, func(servers map[string]json.RawMessage) bool {
			changed := false
			for n := range st.copies {
				if _, shared := group.Servers[n]; shared {
					continue
				}
				if cur, ok := servers[n]; ok && st.ours(n, cur, group) {
					delete(servers, n)
					changed = true
				}
			}
			for _, n := range group.Names() {
				def := group.Servers[n]
				cur, ok := servers[n]
				if _, noted := kept.copies[n]; ok && !noted && !st.ours(n, cur, group) {
					continue // the account's own
				}
				if !ok || !SameDef(cur, def) {
					servers[n] = def
					changed = true
				}
				kept.copies[n] = defHash(def)
			}
			return changed
		})
		if err != nil {
			return err
		}
	}
	if len(kept.copies) > 0 {
		if err := kept.save(keptPath(stateDir)); err != nil {
			return err
		}
	}
	os.Remove(checkedPath(stateDir))
	return removeIfThere(syncStampPath(stateDir))
}

// KeptPlan lists the copies Keep made for the account that the group it
// joins has servers of: DropKept takes them out, since the group gives
// those servers again, and under the flag strategy a copy would hide the
// group's later changes. Copies changed since are the account's own.
// groupMCP is the mcp.json of that group (or the one it starts from).
func KeptPlan(profileDir, groupMCP, stateDir string) ([]string, error) {
	st, ok := readStamp(keptPath(stateDir))
	if !ok {
		return nil, nil
	}
	group, err := Load(groupMCP)
	if err != nil {
		return nil, err
	}
	local, _ := LocalServers(profileDir)
	var out []string
	for _, n := range slices.Sorted(maps.Keys(st.copies)) {
		if _, shared := group.Servers[n]; !shared {
			continue
		}
		if cur, ok := local[n]; ok && st.ours(n, cur, nil) {
			out = append(out, n)
		}
	}
	return out, nil
}

// DropKept takes the copies KeptPlan lists out of .claude.json (with a
// backup, under Claude Code's lock) and forgets the others, which stay as
// the account's own. The caller ensures the account is not running.
func DropKept(profileDir, groupDir, stateDir string) error {
	st, ok := readStamp(keptPath(stateDir))
	if !ok {
		return nil
	}
	group, err := Load(Path(groupDir))
	if err != nil {
		return err
	}
	err = editServers(profileDir, stateDir, false, func(servers map[string]json.RawMessage) bool {
		changed := false
		for n := range st.copies {
			if _, shared := group.Servers[n]; !shared {
				continue
			}
			if cur, ok := servers[n]; ok && st.ours(n, cur, nil) {
				delete(servers, n)
				changed = true
			}
		}
		return changed
	})
	if err != nil {
		return err
	}
	return removeIfThere(keptPath(stateDir))
}

// Remove deletes servers from a profile's .claude.json: user from its
// mcpServers, and for each folder of folders those names from the folder's
// local scope (projects[<folder>].mcpServers). One write, with a backup,
// under Claude Code's lock; the caller ensures the profile is not running.
func Remove(profileDir, stateDir string, user []string, folders map[string][]string) error {
	return editDoc(profileDir, stateDir, false, func(doc map[string]json.RawMessage) (bool, error) {
		changed := false
		if len(user) > 0 {
			servers, err := objectAt(doc, "mcpServers")
			if err != nil {
				return false, err
			}
			if deleteAll(servers, user) {
				if doc["mcpServers"], err = fsutil.MarshalJSON(servers); err != nil {
					return false, err
				}
				changed = true
			}
		}
		if len(folders) == 0 {
			return changed, nil
		}
		projects, err := objectAt(doc, "projects")
		if err != nil {
			return false, err
		}
		moved := false
		for folder, names := range folders {
			entry, err := objectAt(projects, folder)
			if err != nil {
				return false, err
			}
			servers, err := objectAt(entry, "mcpServers")
			if err != nil {
				return false, err
			}
			if !deleteAll(servers, names) {
				continue
			}
			if entry["mcpServers"], err = fsutil.MarshalJSON(servers); err != nil {
				return false, err
			}
			if projects[folder], err = fsutil.MarshalJSON(entry); err != nil {
				return false, err
			}
			moved = true
		}
		if moved {
			if doc["projects"], err = fsutil.MarshalJSON(projects); err != nil {
				return false, err
			}
		}
		return changed || moved, nil
	})
}

// objectAt decodes the object at key of m (empty when it is not there).
func objectAt(m map[string]json.RawMessage, key string) (map[string]json.RawMessage, error) {
	out := map[string]json.RawMessage{}
	raw, ok := m[key]
	if !ok {
		return out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("could not read %s in .claude.json: %w", key, err)
	}
	if out == nil {
		out = map[string]json.RawMessage{}
	}
	return out, nil
}

// deleteAll deletes names from m and reports whether one was there.
func deleteAll(m map[string]json.RawMessage, names []string) bool {
	found := false
	for _, n := range names {
		if _, ok := m[n]; ok {
			delete(m, n)
			found = true
		}
	}
	return found
}

// editServers lets edit change the mcpServers of a profile's .claude.json;
// see editDoc.
func editServers(profileDir, stateDir string, create bool, edit func(servers map[string]json.RawMessage) bool) error {
	return editDoc(profileDir, stateDir, create, func(doc map[string]json.RawMessage) (bool, error) {
		servers, err := objectAt(doc, "mcpServers")
		if err != nil || !edit(servers) {
			return false, err
		}
		// Not json.Marshal: it would write && and > in commands as & and >.
		if doc["mcpServers"], err = fsutil.MarshalJSON(servers); err != nil {
			return false, err
		}
		return true, nil
	})
}

// editDoc lets edit change a profile's .claude.json and writes the file
// when edit reports a change: under Claude Code's lock, after a backup in
// stateDir. Without a .claude.json, create says whether to make one.
func editDoc(profileDir, stateDir string, create bool, edit func(doc map[string]json.RawMessage) (bool, error)) error {
	claudeJSON := filepath.Join(profileDir, ".claude.json")
	unlock, err := lock.DirLock(claudeJSON+".lock", 2*time.Second)
	if err != nil {
		return err
	}
	defer unlock()
	doc := map[string]json.RawMessage{}
	perm := fs.FileMode(0o600)
	data, err := os.ReadFile(claudeJSON)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("could not read .claude.json (left untouched): %w", err)
		}
		if doc == nil {
			doc = map[string]json.RawMessage{}
		}
		if info, err := os.Stat(claudeJSON); err == nil {
			perm = info.Mode().Perm()
		}
	case !errors.Is(err, fs.ErrNotExist):
		return err
	case !create:
		return nil
	}
	changed, err := edit(doc)
	if err != nil || !changed {
		return err
	}
	out, err := fsutil.MarshalJSON(doc)
	if err != nil {
		return err
	}
	if data != nil {
		if err := backup(claudeJSON, data, filepath.Join(stateDir, "backups")); err != nil {
			return err
		}
	}
	return fsutil.AtomicWrite(claudeJSON, out, perm)
}

func removeIfThere(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
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
