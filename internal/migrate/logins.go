package migrate

import (
	"encoding/json"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/share"
)

// Account is an account that may get the source's MCP servers through the
// group.
type Account struct {
	Name       string
	Dir        string // its CLAUDE_CONFIG_DIR
	MCP        string // its share group for mcp
	Extensions string // its share group for extensions
}

// Logins lists the MCP servers the source is signed in to that accounts get
// through the group but have to sign in to themselves, with /mcp: the
// sign-ins (mcpOAuth in the source's login) are not copied. A sign-in
// renews itself with a refresh token only one holder can use; once one copy
// renews it, every other copy (the source's too) is refused and Claude Code
// clears it, and a server that notices the reuse revokes them all.
//
// Claude Code keys a sign-in by the server name and a hash of its type, url
// and headers, so an account signed in to the same definition already needs
// nothing; servers whose group definition is kept do not count. Nothing is
// written: the accounts' login files are only read.
type Logins struct {
	entries map[string]json.RawMessage // mcpOAuth key → a sign-in of the source that holds a token
	targets []*loginTarget
}

type loginTarget struct {
	acct    Account
	servers bool     // gets the user servers (shares mcp through the group)
	plugins bool     // gets plugin servers (shares extensions through the group)
	missing []string // keys it has no sign-in for, as planned (before conflicts are decided)
}

func buildLogins(pl *Plan, opt Options) *Logins {
	data, _, err := creds.Read(pl.Source, opt.Home)
	if err != nil {
		return nil // no login, or one ccam cannot read: nothing to tell
	}
	l := &Logins{entries: creds.MCPSignIns(data)}
	if len(l.entries) == 0 {
		return nil
	}
	// Plugin servers come along with the plugins they belong to.
	servers := slices.Contains(opt.Items, "mcp")
	plugins := slices.Contains(opt.Items, "extensions") && fsutil.Exists(filepath.Join(pl.Source, "plugins"))
	for _, a := range opt.Accounts {
		t := &loginTarget{acct: a, servers: servers && a.MCP == pl.Group, plugins: plugins && a.Extensions == pl.Group}
		if !t.servers && !t.plugins || filepath.Clean(a.Dir) == filepath.Clean(pl.Source) {
			continue
		}
		var have map[string]json.RawMessage
		if own, _, err := creds.Read(a.Dir, opt.Home); err == nil {
			have = creds.MCPSignIns(own)
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
	if len(l.targets) == 0 {
		return nil
	}
	return l
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

// keys lists the sign-ins t would need, given the servers the group ends up
// with in the source's definition (copied, already the same, or "use mine").
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

// needs says which accounts must sign in to which servers, "gh in work,
// side; linear in work", given the servers the group ends up with in the
// source's definition; "" when none.
func (l *Logins) needs(sameDef map[string]bool) string {
	who := map[string][]string{} // server → accounts
	for _, t := range l.targets {
		for _, k := range l.keys(t, sameDef) {
			name, _, _ := strings.Cut(k, "|")
			if slices.Contains(t.missing, k) && !slices.Contains(who[name], t.acct.Name) {
				who[name] = append(who[name], t.acct.Name)
			}
		}
	}
	var parts []string
	for _, name := range slices.Sorted(maps.Keys(who)) {
		parts = append(parts, name+" in "+strings.Join(who[name], ", "))
	}
	return strings.Join(parts, "; ")
}

// signInDefs lists the servers the group ends up with in the source's
// definition: those it may, while the conflicts are still to be decided
// ("use mine" may be chosen), else as they are decided now.
func (pl *Plan) signInDefs(asking bool) map[string]bool {
	same := pl.MCP.planned()
	if m := pl.MCP; m != nil && !asking {
		for _, c := range m.Conflicts {
			if c.Choice != share.UseMine {
				delete(same, serverName(c)) // the group keeps its own definition
			}
		}
	}
	return same
}

// SignIns says which accounts must sign in to which copied servers with
// /mcp, "gh in work, side", with the conflict choices as they are now; ""
// when none.
func (pl *Plan) SignIns() string {
	if pl.Logins == nil {
		return ""
	}
	return pl.Logins.needs(pl.signInDefs(false))
}
