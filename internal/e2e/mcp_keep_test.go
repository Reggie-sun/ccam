package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/jhyunwoo/ccam/internal/transcript"
)

// setStrategy puts mcp_strategy at the top of config.toml, where TOML reads
// it as ccam's (below a table it would belong to that table).
func setStrategy(t *testing.T, h *harness, strategy string) {
	t.Helper()
	path := filepath.Join(h.home, "config", "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var keep []string
	for _, l := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(l, "mcp_strategy") {
			keep = append(keep, l)
		}
	}
	os.WriteFile(path, []byte("mcp_strategy = '"+strategy+"'\n"+strings.Join(keep, "\n")), 0o600)
}

// claudeJSON reads an account's .claude.json: its user-scope servers
// (name → command), its local-scope servers of folder key, and the rest.
func claudeJSON(t *testing.T, path, key string) (user, local map[string]string, doc map[string]any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	commands := func(v any) map[string]string {
		out := map[string]string{}
		m, _ := v.(map[string]any)
		for n, d := range m {
			def, _ := d.(map[string]any)
			out[n], _ = def["command"].(string)
		}
		return out
	}
	user = commands(doc["mcpServers"])
	projects, _ := doc["projects"].(map[string]any)
	entry, _ := projects[key].(map[string]any)
	return user, commands(entry["mcpServers"]), doc
}

// The owner's case. Account A (personal) has a server X of its own and a
// server Y for one folder (what plain `claude mcp add` does); the group
// shares Z. B (work) gets Z; the launch hint, ext and doctor name X and Y;
// A keeps Z when it leaves the group; sync never replaces A's own X.
func TestMCPOwnServersAndLeavingTheGroup(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "Z", `{"type":"stdio","command":"z-mcp"}`)
	groupMCP := filepath.Join(h.home, "data", "groups", "default", "mcp.json")
	project := filepath.Join(h.home, "proj")
	os.MkdirAll(project, 0o700)
	key, _ := filepath.EvalSymlinks(project)
	a := filepath.Join(h.profileDir("personal"), ".claude.json")
	doc, _ := json.Marshal(map[string]any{
		"oauthAccount": map[string]any{"emailAddress": "me@example.com"},
		"mcpServers":   map[string]any{"X": map[string]any{"type": "stdio", "command": "x-own"}},
		"projects":     map[string]any{key: map[string]any{"hasTrustDialogAccepted": true, "mcpServers": map[string]any{"Y": map[string]any{"type": "stdio", "command": "y-local"}}}},
	})
	os.WriteFile(a, doc, 0o600)

	// Switching to B: its sessions get Z.
	h.ok("ccam", "work", "-c")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"--mcp-config=" + groupMCP, "-c"}) {
		t.Fatalf("work: %q", args)
	}
	// The hint names X and Y once.
	r := h.runIn(project, "", nil, "ccam", "personal")
	for _, want := range []string{"MCP servers only account personal has: X", "for one folder only (claude mcp add without -s user): Y (" + key + ")", "ccam mcp import personal Y"} {
		if !strings.Contains(r.stderr, want) {
			t.Errorf("hint lacks %q:\n%s", want, r.stderr)
		}
	}
	h.called()
	if r := h.runIn(project, "", nil, "ccam", "personal"); strings.Contains(r.stderr, "MCP servers") {
		t.Errorf("hint again: %s", r.stderr)
	}
	h.called()
	// ext, in the folder: X and Y in A's column only, and why.
	r = h.runIn(project, "", nil, "ccam", "ext")
	rows := map[string]string{}
	for _, l := range strings.Split(r.stdout, "\n") {
		if f := strings.Fields(l); len(f) == 3 {
			rows[f[0]] = f[1] + " " + f[2]
		}
	}
	if rows["X"] != "account -" || rows["Y"] != "local -" || rows["Z"] != "default default" {
		t.Errorf("ext rows %v:\n%s", rows, r.stdout)
	}
	for _, want := range []string{"! personal: MCP server X is only in this account's .claude.json", "! personal: MCP server Y is in this account's local scope for this folder"} {
		if !strings.Contains(r.stdout, want) {
			t.Errorf("ext lacks %q:\n%s", want, r.stdout)
		}
	}
	// doctor too.
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "MCP servers the other accounts of group default do not load: X; for one folder only: Y ("+key+") → `ccam mcp import personal X Y`") {
		t.Errorf("doctor:\n%s", out)
	}

	// Leaving the group keeps Z in A, beside its own servers.
	out := h.ok("ccam", "group", "leave", "default", "personal", "--yes")
	if !strings.Contains(out, "copy     MCP servers Z: "+filepath.FromSlash("groups/default/mcp.json")+" → "+filepath.FromSlash("profiles/personal/.claude.json")) {
		t.Fatalf("leave plan:\n%s", out)
	}
	h.ok("ccam", "personal", "-c")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"-c"}) {
		t.Fatalf("after leaving: %q", args)
	}
	user, local, rest := claudeJSON(t, a, key)
	if !reflect.DeepEqual(user, map[string]string{"X": "x-own", "Z": "z-mcp"}) || local["Y"] != "y-local" || rest["oauthAccount"] == nil {
		t.Fatalf("after leaving: %v %v %v", user, local, rest)
	}

	// Back in the group the kept copy of Z goes, and the group gives Z again.
	out = h.ok("ccam", "group", "join", "default", "personal", "--yes")
	if !strings.Contains(out, "remove   MCP servers Z kept in .claude.json when personal left a group") {
		t.Fatalf("join plan:\n%s", out)
	}
	if user, local, _ := claudeJSON(t, a, key); !reflect.DeepEqual(user, map[string]string{"X": "x-own"}) || local["Y"] != "y-local" {
		t.Fatalf("after joining again: %v %v", user, local)
	}

	// Under sync, the group's own X does not replace A's.
	setStrategy(t, h, "sync")
	h.ok("ccam", "mcp", "add-json", "X", `{"type":"stdio","command":"x-group"}`)
	r = h.run("", nil, "ccam", "personal", "-c")
	if !strings.Contains(r.stderr, "account personal has its own MCP servers X, so group default's of those names are not copied there") {
		t.Errorf("sync: %s", r.stderr)
	}
	if user, _, _ := claudeJSON(t, a, key); !reflect.DeepEqual(user, map[string]string{"X": "x-own", "Z": "z-mcp"}) {
		t.Fatalf("after sync: %v", user)
	}
	if r := h.run("", nil, "ccam", "personal", "-c"); strings.Contains(r.stderr, "its own") {
		t.Errorf("said again: %s", r.stderr)
	}
	h.ok("ccam", "work", "-c")
	if user, _, _ := claudeJSON(t, filepath.Join(h.profileDir("work"), ".claude.json"), key); user["X"] != "x-group" || user["Z"] != "z-mcp" {
		t.Fatalf("work after sync: %v", user)
	}
	// A copy removed in a session comes back at the next launch.
	_, _, rest = claudeJSON(t, a, key)
	rest["mcpServers"] = map[string]any{"X": map[string]any{"type": "stdio", "command": "x-own"}}
	doc, _ = json.Marshal(rest)
	os.WriteFile(a, doc, 0o600)
	r = h.run("", nil, "ccam", "personal", "-c")
	if !strings.Contains(r.stderr, "put back group default's MCP servers that were gone from account personal's .claude.json: Z") {
		t.Errorf("restore: %s", r.stderr)
	}
	// Leaving under sync keeps both: X its own, Z as the group has it.
	h.ok("ccam", "group", "leave", "default", "personal", "--yes")
	if user, _, _ := claudeJSON(t, a, key); !reflect.DeepEqual(user, map[string]string{"X": "x-own", "Z": "z-mcp"}) {
		t.Fatalf("after leaving under sync: %v", user)
	}
}

// claude mcp list|get|remove read only .claude.json, .mcp.json and the
// plugins: ccam passes them no --mcp-config and says where the group's
// servers are.
func TestMCPListSaysTheGroupsServers(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-mcp"}`)
	for _, verb := range []string{"list", "get", "remove"} {
		args := []string{"ccam", "work", "mcp", verb}
		if verb != "list" {
			args = append(args, "fs")
		}
		r := h.run("", nil, args...)
		if _, got := h.called(); !reflect.DeepEqual(got, args[2:]) {
			t.Fatalf("%s: %q", verb, got)
		}
		if !strings.Contains(r.stderr, "`claude mcp "+verb+"` sees only account work's own MCP servers; its sessions also load group default's: fs") {
			t.Fatalf("%s: %s", verb, r.stderr)
		}
	}
	r := h.run("", nil, "ccam", "work", "mcp", "add", "x", "--", "y")
	if _, got := h.called(); len(got) == 0 || !strings.HasPrefix(got[0], "--mcp-config=") || strings.Contains(r.stderr, "sees only") {
		t.Fatalf("mcp add: %q %s", got, r.stderr)
	}
}

// A server `claude mcp add` gave one folder (local scope) can be shared:
// import copies it to the group and takes it out of that folder's entry,
// leaving the entry's other keys and the user scope as they were.
func TestMCPImportOfAFolderServer(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	a := filepath.Join(h.profileDir("personal"), ".claude.json")
	os.WriteFile(a, []byte(`{"oauthAccount":{"emailAddress":"me@example.com"},"mcpServers":{"own":{"command":"o"}},`+
		`"projects":{"/p":{"hasTrustDialogAccepted":true,"mcpServers":{"loc":{"command":"l"},"stay":{"command":"s"}}},"/q":{"allowedTools":["x"]}}}`), 0o600)
	out := h.ok("ccam", "mcp", "list")
	if !strings.Contains(out, "only in account personal, for /p only: loc, stay") || !strings.Contains(out, "only in account personal: own") {
		t.Fatalf("list:\n%s", out)
	}
	out = h.ok("ccam", "mcp", "import", "personal", "loc")
	if !strings.Contains(out, "Copied to group default: loc") || !strings.Contains(out, "every account in it gets them in every folder") {
		t.Fatalf("import:\n%s", out)
	}
	if data, _ := os.ReadFile(filepath.Join(h.home, "data", "groups", "default", "mcp.json")); !strings.Contains(string(data), `"loc"`) {
		t.Fatalf("group: %s", data)
	}
	user, local, doc := claudeJSON(t, a, "/p")
	q := doc["projects"].(map[string]any)["/q"].(map[string]any)
	if !reflect.DeepEqual(user, map[string]string{"own": "o"}) || !reflect.DeepEqual(local, map[string]string{"stay": "s"}) || q["allowedTools"] == nil ||
		doc["projects"].(map[string]any)["/p"].(map[string]any)["hasTrustDialogAccepted"] != true {
		t.Fatalf("after import: %v %v %v", user, local, doc)
	}
	if r := h.run("", nil, "ccam", "mcp", "import", "personal", "none"); r.code == 0 || !strings.Contains(r.stderr, `account personal has no server "none"`) {
		t.Fatalf("unknown: %+v", r)
	}
}

// ccam commands that change the group's servers at the same moment all
// keep their change.
func TestMCPAddsAtOnceKeepEveryServer(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// h.run, not h.ok: t.Fatalf may not be called from another goroutine.
			if r := h.run("", nil, "ccam", "mcp", "add-json", fmt.Sprint("s", i), `{"type":"stdio","command":"x"}`); r.code != 0 {
				t.Errorf("s%d: %+v", i, r)
			}
		}()
	}
	wg.Wait()
	var store struct {
		Servers map[string]any `json:"mcpServers"`
	}
	data, _ := os.ReadFile(filepath.Join(h.home, "data", "groups", "default", "mcp.json"))
	if json.Unmarshal(data, &store); len(store.Servers) != 8 {
		t.Fatalf("kept %d of 8: %s", len(store.Servers), data)
	}
}

// A change of groups that stops halfway is finished by doctor --fix, MCP
// servers of the new group included; another change waits for that.
func TestStoppedGroupChangeIsFinishedByDoctor(t *testing.T) {
	if !perms || os.Geteuid() == 0 {
		t.Skip("needs folder permissions")
	}
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "github", `{"type":"stdio","command":"gh"}`)
	h.ok("ccam", "group", "add", "two", "--items", "mcp,settings", "--yes")
	two := filepath.Join(h.home, "data", "groups", "two")
	os.MkdirAll(two, 0o500)
	r := h.run("", nil, "ccam", "group", "join", "two", "work", "--replace", "--yes")
	os.Chmod(two, 0o700)
	if r.code == 0 || !strings.Contains(r.stderr, "finish with `ccam doctor --fix`") {
		t.Fatalf("join: %+v", r)
	}
	if r := h.run("", nil, "ccam", "group", "join", "two", "personal", "--yes"); r.code == 0 || !strings.Contains(r.stderr, "finish it with `ccam doctor --fix` first") {
		t.Fatalf("a second change: %+v", r)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "the last change of groups stopped halfway") {
		t.Fatalf("doctor:\n%s", out)
	}
	if out := h.run("", nil, "ccam", "doctor", "--fix").stdout; !strings.Contains(out, "finished the last change of groups (work)") {
		t.Fatalf("doctor --fix:\n%s", out)
	}
	h.ok("ccam", "work", "-c")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"--mcp-config=" + filepath.Join(two, "mcp.json"), "-c"}) {
		t.Fatalf("work: %q", args)
	}
	h.ok("ccam", "group", "join", "two", "personal", "--replace", "--yes")
}

// Servers plain Claude Code has in ~/.claude.json that no account loads
// are named at the first add, by mcp list, ext and doctor, until migrate
// copies them.
func TestHomeMCPServersAreNamed(t *testing.T) {
	h := newHarness(t)
	home := filepath.Join(h.home, "home")
	os.MkdirAll(filepath.Join(home, ".claude", "projects"), 0o700)
	os.WriteFile(filepath.Join(home, ".claude.json"), []byte(`{"mcpServers":{"dokploy-mcp":{"command":"d"},"context7":{"type":"http","url":"https://c"}}}`), 0o600)
	r := h.run("", []string{"FAKE_EMAIL=me@example.com"}, "ccam", "add", "personal")
	if !strings.Contains(r.stdout, "context7, dokploy-mcp") || !strings.Contains(r.stdout, "`ccam migrate --mcp`") {
		t.Fatalf("first add: %+v", r)
	}
	h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "add", "work")
	want := "has MCP servers no account loads: context7, dokploy-mcp"
	if out := h.ok("ccam", "mcp", "list"); !strings.Contains(out, "only in "+filepath.Join(home, ".claude.json")+" (plain claude's): context7, dokploy-mcp") {
		t.Errorf("mcp list:\n%s", out)
	}
	if out := h.ok("ccam", "ext"); !strings.Contains(out, want) {
		t.Errorf("ext:\n%s", out)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, want) {
		t.Errorf("doctor:\n%s", out)
	}
	h.ok("ccam", "migrate", "--mcp", "--yes")
	if out := h.run("", nil, "ccam", "doctor").stdout; strings.Contains(out, "no account loads") {
		t.Errorf("doctor after migrate:\n%s", out)
	}
}

// A server one account is signed in to and another is not shows in ext,
// doctor and when the conversation moves on; no sign-in is copied.
func TestMCPSignInGaps(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "remote", `{"type":"http","url":"https://mcp.example/mcp"}`)
	signIn := func(account, token string) {
		path := filepath.Join(h.profileDir(account), ".credentials.json")
		var c map[string]any
		data, _ := os.ReadFile(path)
		json.Unmarshal(data, &c)
		c["mcpOAuth"] = map[string]any{"remote|0123456789abcdef": map[string]any{"accessToken": token, "refreshToken": token}}
		data, _ = json.Marshal(c)
		os.WriteFile(path, data, 0o600)
	}
	signIn("personal", "t")
	signIn("work", "") // what a health check leaves: no sign-in
	before, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".credentials.json"))
	if out := h.ok("ccam", "ext"); !strings.Contains(out, "! work: MCP server remote needs a sign-in in this account: personal signed in") {
		t.Errorf("ext:\n%s", out)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "not signed in to MCP servers other accounts are signed in to: remote (personal)") {
		t.Errorf("doctor:\n%s", out)
	}
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("personal"), "projects", transcript.Slug(real), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte(`{"type":"user","message":{"role":"user","content":"hi"},"timestamp":"2026-10-23T00:00:00Z"}`+"\n"), 0o600)
	r := h.runIn(project, "", nil, "ccam", "next", "--continue", "--", "-p", "go on")
	if r.code != 0 || !strings.Contains(r.stderr, "work is not signed in to MCP servers personal is: remote (sign in with /mcp there") {
		t.Fatalf("next --continue: %+v", r)
	}
	if after, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".credentials.json")); string(after) != string(before) {
		t.Fatal("a sign-in was copied")
	}
}

// Project servers wait for approval in an account that has not trusted
// the folder: ext says so and how to approve them for every account.
func TestProjectMCPApprovalPerAccount(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	project := filepath.Join(h.home, "proj")
	os.MkdirAll(project, 0o700)
	key, _ := filepath.EvalSymlinks(project)
	os.WriteFile(filepath.Join(project, ".mcp.json"), []byte(`{"mcpServers":{"projsrv":{"command":"/bin/cat"}}}`), 0o600)
	doc, _ := json.Marshal(map[string]any{"projects": map[string]any{key: map[string]any{"hasTrustDialogAccepted": true}}})
	os.WriteFile(filepath.Join(h.profileDir("personal"), ".claude.json"), doc, 0o600)
	out := h.runIn(project, "", nil, "ccam", "ext").stdout
	if !strings.Contains(out, "! work: has not trusted this folder yet") || strings.Contains(out, "! personal: has not trusted") || strings.Contains(out, "), for every account:") {
		t.Fatalf("ext:\n%s", out)
	}
	if !slices.Contains(strings.Fields(out), `"enabledMcpjsonServers":`) {
		t.Errorf("no approval advice:\n%s", out)
	}
	out = h.runIn(project, "", nil, "ccam", "mcp", "list").stdout
	if !strings.Contains(out, "this folder's .mcp.json: projsrv") || !strings.Contains(out, "waiting for approval in account work, which has not trusted the folder: projsrv") ||
		strings.Contains(out, "account personal, which") {
		t.Errorf("mcp list:\n%s", out)
	}
}

// An account with the Remote Control copy an older ccam made (its sync
// stamp lists the names only) gets the group's change of that server
// after the upgrade: under the flag strategy the copy in .claude.json wins
// over the group's, so a copy left behind would hide the change.
func TestMCPCopyOfAnOlderCcamFollowsTheGroup(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-1"}`)
	groupDir := filepath.Join(h.home, "data", "groups", "default")
	info, err := os.Stat(filepath.Join(groupDir, "mcp.json"))
	if err != nil {
		t.Fatal(err)
	}
	a := filepath.Join(h.profileDir("personal"), ".claude.json")
	os.WriteFile(a, []byte(`{"oauthAccount":{"emailAddress":"me@example.com"},"mcpServers":{"fs":{"type":"stdio","command":"fs-1"},"own":{"type":"stdio","command":"own"}}}`), 0o600)
	state := filepath.Join(h.home, "state", "mcp", "personal")
	os.MkdirAll(state, 0o700)
	os.WriteFile(filepath.Join(state, "sync-stamp"), []byte(fmt.Sprintf("%d %s\nfs\n", info.ModTime().UnixNano(), groupDir)), 0o600)

	out := h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-2"}`)
	if user, _, _ := claudeJSON(t, a, ""); user["fs"] != "fs-2" || user["own"] != "own" {
		t.Fatalf("personal's servers %v after:\n%s", user, out)
	}
	if strings.Contains(out, "has its own") || !strings.Contains(out, "✓ Updated the copy in account personal's .claude.json") {
		t.Errorf("the copy was called the account's own:\n%s", out)
	}
	// A group server of the name of the account's own: nothing is copied,
	// and ccam says so instead of an update.
	out = h.ok("ccam", "mcp", "add-json", "own", `{"type":"stdio","command":"group-own"}`)
	if user, _, _ := claudeJSON(t, a, ""); user["own"] != "own" || strings.Contains(out, "Updated the copy") || !strings.Contains(out, "account personal has its own own") {
		t.Errorf("personal's servers %v after:\n%s", user, out)
	}
}
