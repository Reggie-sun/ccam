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
	h.ok("ccs", "work", "-c")
	if _, args := h.called(); !reflect.DeepEqual(args, []string{"-c"}) {
		t.Fatalf("args %q", args)
	}

	h.ok("ccam", "mcp", "add", "fs", "-e", "TOKEN=x", "--", "npx", "-y", "server-fs", "/code")
	h.ok("ccam", "mcp", "add", "ctx", "--url", "https://mcp.example/mcp", "--header", "Authorization: Bearer t")
	if info, _ := os.Stat(groupMCP); info.Mode().Perm() != 0o600 {
		t.Fatal("mcp.json must be 0600")
	}
	h.ok("ccs", "work", "--mcp-config", "mine.json", "-c")
	_, args := h.called()
	want := []string{"--mcp-config", groupMCP, "--mcp-config", "mine.json", "-c"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args %q", args)
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

	// Local MCP off: no flag even though the group has servers.
	h.ok("ccam", "share", "work", "--mcp", "off", "--yes")
	h.ok("ccs", "work")
	if _, args := h.called(); len(args) != 0 {
		t.Fatalf("mcp off still passes %q", args)
	}
}

func TestMCPNewLocalServerHintAndImport(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	// Simulate `claude mcp add` inside a work session.
	os.WriteFile(filepath.Join(h.profileDir("work"), ".claude.json"),
		[]byte(`{"oauthAccount":{"emailAddress":"dev@company.com"},"mcpServers":{"gh":{"type":"stdio","command":"gh-mcp"}}}`), 0o600)
	r := h.run("", nil, "ccs", "work")
	if !strings.Contains(r.stderr, "gh") || !strings.Contains(r.stderr, "ccam mcp import work gh") {
		t.Fatalf("no hint: %q", r.stderr)
	}
	if r := h.run("", nil, "ccs", "work"); strings.Contains(r.stderr, "gh") {
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
	h.ok("ccs", "work", "-c")
	_, args := h.called()
	if !reflect.DeepEqual(args, []string{"-c"}) {
		t.Fatalf("sync must not add flags: %q", args)
	}
	local, _ := os.ReadFile(filepath.Join(h.profileDir("work"), ".claude.json"))
	if !strings.Contains(string(local), "fs-mcp") || !strings.Contains(string(local), "dev@company.com") {
		t.Fatalf(".claude.json after sync: %s", local)
	}
}

func TestProfileLockHeldWhileClaudeRuns(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	lockPath := filepath.Join(h.home, "run", "profile-work.lock")
	h.run("", []string{"FAKE_LOCK_PROBE=" + lockPath}, "ccs", "work")
	if b, _ := os.ReadFile(lockPath + ".result"); strings.TrimSpace(string(b)) != "held" {
		t.Fatalf("lock while running: %q", b)
	}
	h.run("", []string{"FAKE_LOCK_PROBE=" + filepath.Join(h.home, "run", "profile-personal.lock")}, "ccs", "work")
	if b, _ := os.ReadFile(filepath.Join(h.home, "run", "profile-personal.lock.result")); strings.TrimSpace(string(b)) != "free" {
		t.Fatalf("other profile lock: %q", b)
	}
}
