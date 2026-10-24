package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMCPFlagStrategy(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	groupMCP := filepath.Join(h.home, "data", "groups", "default", "mcp.json")

	// No group servers yet: no flag.
	h.ok("ccam", "work", "-c")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"-c"}) {
		t.Fatalf("args %q", args)
	}

	h.ok("ccam", "mcp", "add", "fs", "-e", "TOKEN=x", "--", "npx", "-y", "server-fs", "/code")
	h.ok("ccam", "mcp", "add", "ctx", "--url", "https://mcp.example/mcp", "--header", "Authorization: Bearer t")
	if info, _ := os.Stat(groupMCP); perms && info.Mode().Perm() != 0o600 {
		t.Fatal("mcp.json must be 0600")
	}
	h.ok("ccam", "work", "--mcp-config", "mine.json", "-c")
	_, args := h.called()
	want := []string{"--mcp-config=" + groupMCP, "--mcp-config", "mine.json", "-c"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args %q", args)
	}
	// --mcp-config is variadic in claude: the space form would swallow a
	// following prompt or subcommand as another config file.
	h.ok("ccam", "work", "fix the bug")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"--mcp-config=" + groupMCP, "fix the bug"}) {
		t.Fatalf("prompt after flag: %q", args)
	}
	// claude remote-control refuses a --mcp-config before it and takes none
	// after it, so ccam leaves the flag out and says the sessions lack the
	// servers; also after an option with its value.
	for _, verb := range []string{"remote-control", "rc"} {
		r := h.run("", nil, "ccam", "work", "--model", "opus", verb, "--name", "x")
		if r.code != 0 {
			t.Fatalf("%s exited %d: %s", verb, r.code, r.stderr)
		}
		if _, args := h.called(); !reflect.DeepEqual(args, []string{"--model", "opus", verb, "--name", "x"}) {
			t.Fatalf("%s args %q", verb, args)
		}
		if !strings.Contains(r.stderr, "Remote Control") || !strings.Contains(r.stderr, `mcp_strategy = "sync"`) {
			t.Fatalf("%s stderr: %s", verb, r.stderr)
		}
		// -d takes the verb as its filter: a session, with the servers.
		h.ok("ccam", "work", "-d", verb)
		if _, args := h.called(); !reflect.DeepEqual(args, []string{"--mcp-config=" + groupMCP, "-d", verb}) {
			t.Fatalf("-d %s args %q", verb, args)
		}
	}
	// Claude Code finds a background session by its name only first.
	for _, verb := range []string{"logs", "attach", "stop", "kill", "respawn", "rm"} {
		h.ok("ccam", "work", verb, "myjob")
		if _, args := h.called(); !reflect.DeepEqual(args, []string{verb, "myjob"}) {
			t.Fatalf("%s args %q", verb, args)
		}
	}
	// Without a terminal nobody can answer the question, so .claude.json stays as it was.
	if local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json")); strings.Contains(string(local), "server-fs") {
		t.Fatalf(".claude.json changed without an answer: %s", local)
	}
	var store struct {
		Servers map[string]map[string]any `json:"mcpServers"`
	}
	data, _ := os.ReadFile(groupMCP)
	json.Unmarshal(data, &store)
	if store.Servers["fs"]["command"] != "npx" || store.Servers["ctx"]["type"] != "http" {
		t.Fatalf("store %s", data)
	}
	out := h.ok("ccam", "mcp", "list")
	if !strings.Contains(out, "stdio npx -y server-fs /code") || !strings.Contains(out, "http https://mcp.example/mcp") {
		t.Fatalf("list:\n%s", out)
	}
	h.ok("ccam", "mcp", "remove", "fs")
	if out := h.ok("ccam", "mcp", "list"); strings.Contains(out, "server-fs") {
		t.Fatal("remove failed")
	}

	// Local MCP off: no flag even though the group has servers; the
	// account keeps them in its own .claude.json.
	out = h.ok("ccam", "group", "leave", "default", "work", "--yes")
	if !strings.Contains(out, "copy     MCP servers ctx: "+filepath.FromSlash("groups/default/mcp.json")+" → "+filepath.FromSlash("profiles/work/.claude.json")) {
		t.Fatalf("leave plan:\n%s", out)
	}
	h.ok("ccam", "work")
	if _, args := h.called(); len(args) != 0 {
		t.Fatalf("mcp off still passes %q", args)
	}
	if local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json")); !strings.Contains(string(local), "mcp.example") || !strings.Contains(string(local), "dev@company.com") {
		t.Fatalf("work kept no copy: %s", local)
	}
}

func TestMCPNewLocalServerHintAndImport(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	// Simulate `claude mcp add` inside a work session.
	os.WriteFile(filepath.Join(h.profileDir("work"), ".claude.json"),
		[]byte(`{"oauthAccount":{"emailAddress":"dev@company.com"},"mcpServers":{"gh":{"type":"stdio","command":"gh-mcp"}}}`), 0o600)
	r := h.run("", nil, "ccam", "work")
	if !strings.Contains(r.stderr, "gh") || !strings.Contains(r.stderr, "ccam mcp import work gh") {
		t.Fatalf("no hint: %q", r.stderr)
	}
	if r := h.run("", nil, "ccam", "work"); strings.Contains(r.stderr, "gh") {
		t.Fatal("hint repeated without changes")
	}
	h.ok("ccam", "mcp", "import", "work", "gh")
	data, _ := os.ReadFile(filepath.Join(h.home, "data", "groups", "default", "mcp.json"))
	if !strings.Contains(string(data), "gh-mcp") {
		t.Fatalf("not imported: %s", data)
	}
	local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json"))
	if strings.Contains(string(local), "gh-mcp") || !strings.Contains(string(local), "dev@company.com") {
		t.Fatalf("local copy: %s", local)
	}
}

func TestMCPSyncStrategy(t *testing.T) {
	h := newHarness(t)
	cfgPath := filepath.Join(h.home, "config", "config.toml")
	data, _ := os.ReadFile(cfgPath)
	os.WriteFile(cfgPath, append(data, []byte("mcp_strategy = \"sync\"\n")...), 0o600)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-mcp"}`)
	h.ok("ccam", "work", "-c")
	_, args := h.called()
	if !reflect.DeepEqual(args, []string{"-c"}) {
		t.Fatalf("sync must not add flags: %q", args)
	}
	local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json"))
	if !strings.Contains(string(local), "fs-mcp") || !strings.Contains(string(local), "dev@company.com") {
		t.Fatalf(".claude.json after sync: %s", local)
	}
}

// A copy in .claude.json (made for Remote Control) would win over --mcp-config,
// so changing the group updates it.
func TestMCPCopyFollowsGroupChanges(t *testing.T) {
	h := newHarness(t)
	cfgPath := filepath.Join(h.home, "config", "config.toml")
	data, _ := os.ReadFile(cfgPath)
	os.WriteFile(cfgPath, append(data, []byte("mcp_strategy = \"sync\"\n")...), 0o600)
	h.addAccounts()
	h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-mcp"}`)
	h.ok("ccam", "work", "-c") // leaves a copy in work's .claude.json
	h.called()
	data, _ = os.ReadFile(cfgPath)
	os.WriteFile(cfgPath, []byte(strings.Replace(string(data), "mcp_strategy = 'sync'", "mcp_strategy = 'flag'", 1)), 0o600)

	out := h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-mcp-2"}`)
	if !strings.Contains(out, "Updated the copy in account work") {
		t.Fatalf("add output:\n%s", out)
	}
	local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json"))
	if !strings.Contains(string(local), "fs-mcp-2") {
		t.Fatalf(".claude.json after add: %s", local)
	}
	h.ok("ccam", "mcp", "remove", "fs")
	if local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json")); strings.Contains(string(local), "fs-mcp") {
		t.Fatalf(".claude.json after remove: %s", local)
	}
}

func TestSessionMarkerNamesTheClaudeProcess(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	dir := filepath.Join(h.home, "run", "profile-work")
	h.run("", []string{"FAKE_SESSION_PROBE=" + dir}, "ccam", "work")
	if b, _ := os.ReadFile(dir + ".result"); strings.TrimSpace(string(b)) != "held" {
		t.Fatalf("marker while running: %q", b)
	}
	other := filepath.Join(h.home, "run", "profile-personal")
	h.run("", []string{"FAKE_SESSION_PROBE=" + other}, "ccam", "work")
	if b, _ := os.ReadFile(other + ".result"); strings.TrimSpace(string(b)) != "free" {
		t.Fatalf("other profile: %q", b)
	}
}
