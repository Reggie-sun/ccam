package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/share"
)

// Account is an account that may receive the source's MCP sign-ins.
type Account struct {
	Name       string
	Dir        string // its CLAUDE_CONFIG_DIR
	MCP        string // its share group for mcp
	Extensions string // its share group for extensions
	SessionDir string // markers of its running sessions
	Imported   bool   // registered in place: plain claude may write its login file too
}

// Logins copies the sign-ins of MCP servers (mcpOAuth in the source's login)
// into the own login file of each account that gets those servers through
// the group. They never go into the group: login files stay per account.
//
// Claude Code keys a sign-in by the server name and a hash of its type, url
// and headers, so a copy only applies where the definition is the same as
// the source's; servers whose group definition is kept get none.
type Logins struct {
	entries map[string]json.RawMessage // mcpOAuth key → entry, from the source
	targets []*loginTarget
	Skipped []string // "name (reason)"
	written []string // accounts Apply gave sign-ins to
}

type loginTarget struct {
	acct    Account
	servers bool     // gets the user servers' sign-ins (shares mcp through the group)
	plugins bool     // gets plugin servers' sign-ins (shares extensions through the group)
	missing []string // keys its login file lacks, as planned (before conflicts are decided)
}

func buildLogins(pl *Plan, opt Options) (*Logins, error) {
	data, src, err := creds.Read(pl.Source, opt.Home)
	if errors.Is(err, creds.ErrNone) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read the login of %s: %w", pl.Source, err)
	}
	_ = src // read only; the Keychain is never written
	raw, ok, err := profile.TopLevelValue(data, "mcpOAuth")
	if err != nil || !ok {
		return nil, err
	}
	l := &Logins{}
	if err := json.Unmarshal(raw, &l.entries); err != nil || len(l.entries) == 0 {
		return nil, nil // nothing usable to copy
	}
	// Plugin servers' sign-ins go along with the plugins they belong to.
	servers := slices.Contains(opt.Items, "mcp")
	plugins := slices.Contains(opt.Items, "extensions") && fsutil.Exists(filepath.Join(pl.Source, "plugins"))
	for _, a := range opt.Accounts {
		t := &loginTarget{acct: a, servers: servers && a.MCP == pl.Group, plugins: plugins && a.Extensions == pl.Group}
		if !t.servers && !t.plugins || filepath.Clean(a.Dir) == filepath.Clean(pl.Source) {
			continue
		}
		if reason := t.skip(opt.Home); reason != "" {
			l.Skipped = append(l.Skipped, a.Name+" ("+reason+")")
			continue
		}
		have, err := signIns(creds.FilePath(a.Dir))
		if err != nil {
			l.Skipped = append(l.Skipped, a.Name+" (could not read its login)")
			continue
		}
		for _, k := range l.keys(t, pl.MCP.planned()) {
			if _, ok := have[k]; !ok {
				t.missing = append(t.missing, k)
			}
		}
		if len(t.missing) > 0 {
			l.targets = append(l.targets, t)
		}
	}
	return l, nil
}

// planned lists the servers the group may end up with in the source's
// definition: new, already the same, or differing (if "use mine" is chosen).
func (m *MCP) planned() map[string]bool {
	out := map[string]bool{}
	if m == nil {
		return out
	}
	for _, n := range append(append([]string{}, m.New...), m.Same...) {
		out[n] = true
	}
	for _, c := range m.Conflicts {
		out[serverName(c)] = true
	}
	return out
}

func signIns(path string) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	have := map[string]json.RawMessage{}
	raw, ok, err := profile.TopLevelValue(data, "mcpOAuth")
	if err == nil && ok {
		err = json.Unmarshal(raw, &have)
	}
	return have, err
}

// skip says why an account cannot get sign-ins now, or "".
func (t *loginTarget) skip(home string) string {
	switch {
	case t.acct.Imported:
		return "registered in place; plain claude may be using its login file"
	case session.Active(t.acct.SessionDir):
		return "running; Claude Code rewrites its login file"
	}
	_, src, err := creds.Read(t.acct.Dir, home)
	switch {
	case errors.Is(err, creds.ErrNone):
		return "not logged in yet"
	case err != nil:
		return "could not read its login"
	case src == creds.Keychain:
		return "its login is in the macOS Keychain, which ccam does not write"
	}
	return ""
}

// keys lists the sign-ins t gets, given the servers the group ends up with
// in the source's definition (copied, already the same, or "use mine").
func (l *Logins) keys(t *loginTarget, sameDef map[string]bool) []string {
	var out []string
	for k := range l.entries {
		name, _, _ := strings.Cut(k, "|")
		if t.servers && sameDef[name] || t.plugins && strings.HasPrefix(name, "plugin:") {
			out = append(out, k)
		}
	}
	slices.Sort(out)
	return out
}

func (l *Logins) apply(pl *Plan) error {
	same := pl.MCP.planned()
	if m := pl.MCP; m != nil {
		for _, c := range m.Conflicts {
			if c.Choice != share.UseMine {
				delete(same, serverName(c)) // the group keeps its own definition
			}
		}
	}
	for _, t := range l.targets {
		// The account may have started since the plan was made.
		if session.Active(t.acct.SessionDir) {
			l.Skipped = append(l.Skipped, t.acct.Name+" (started meanwhile)")
			continue
		}
		n, err := addSignIns(creds.FilePath(t.acct.Dir), l.entries, l.keys(t, same))
		if err != nil {
			return fmt.Errorf("MCP sign-ins for %s: %w", t.acct.Name, err)
		}
		if n > 0 {
			l.written = append(l.written, t.acct.Name)
		}
	}
	return nil
}

// addSignIns adds the entries the login file lacks; ones it has (perhaps
// newer) are kept. Every other key of the file is written back unchanged.
func addSignIns(path string, entries map[string]json.RawMessage, keys []string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	doc := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return 0, err
	}
	have := map[string]json.RawMessage{}
	if raw, ok := doc["mcpOAuth"]; ok {
		if err := json.Unmarshal(raw, &have); err != nil {
			return 0, err
		}
	}
	added := 0
	for _, k := range keys {
		if _, ok := have[k]; !ok {
			have[k] = entries[k]
			added++
		}
	}
	if added == 0 {
		return 0, nil
	}
	raw, err := json.Marshal(have)
	if err != nil {
		return 0, err
	}
	doc["mcpOAuth"] = raw
	out, err := json.Marshal(doc)
	if err != nil {
		return 0, err
	}
	return added, fsutil.AtomicWrite(path, out, 0o600)
}

func (l *Logins) changes() bool { return len(l.targets) > 0 }

// describe is the plan line: which servers' sign-ins go to which accounts.
func (l *Logins) describe() string {
	names := map[string]bool{}
	var to []string
	for _, t := range l.targets {
		to = append(to, t.acct.Name)
		for _, k := range t.missing {
			name, _, _ := strings.Cut(k, "|")
			names[name] = true
		}
	}
	s := "none needed"
	if len(to) > 0 {
		s = strings.Join(slices.Sorted(maps.Keys(names)), ", ") + " → " + strings.Join(to, ", ")
	}
	if len(l.Skipped) > 0 {
		s += "; skipped: " + strings.Join(l.Skipped, ", ")
	}
	return s
}
