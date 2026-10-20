package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtShowsWhatEachAccountLoads(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "mcp", "add", "ctx", "--", "npx", "ctx-server")
	// Extensions move to a group of work alone; personal keeps none.
	h.ok("ccam", "group", "set", "default", "--remove", "extensions", "--fresh", "--yes")
	h.ok("ccam", "group", "add", "ext", "--items", "extensions", "--members", "work", "--yes")
	// A skill written through work's shared skills folder lands in the group.
	skill := filepath.Join(h.profileDir("work"), "skills", "review", "SKILL.md")
	os.MkdirAll(filepath.Dir(skill), 0o700)
	os.WriteFile(skill, []byte("---\nname: review\n---\n"), 0o600)
	// personal has a skill of its own, and a server of its own that clashes
	// with the group's.
	own := filepath.Join(h.profileDir("personal"), "skills", "mine", "SKILL.md")
	os.MkdirAll(filepath.Dir(own), 0o700)
	os.WriteFile(own, []byte("x"), 0o600)
	os.WriteFile(filepath.Join(h.profileDir("personal"), ".claude.json"), []byte(`{"mcpServers":{"ctx":{"command":"other"}}}`), 0o600)

	project := filepath.Join(h.home, "proj")
	os.MkdirAll(filepath.Join(project, ".claude", "commands"), 0o700)
	os.WriteFile(filepath.Join(project, ".claude", "commands", "ship.md"), []byte("ship"), 0o600)

	r := h.runIn(project, "", nil, "ccam", "ext")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	for _, want := range []string{"Skills", "review", "mine", "MCP servers", "ctx", "This project", "commands: ship",
		"! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json"} {
		if !strings.Contains(r.stdout, want) {
			t.Errorf("missing %q in:\n%s", want, r.stdout)
		}
	}
	lines := strings.Split(r.stdout, "\n")
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) == 3 && f[0] == "review" && (f[1] != "-" || f[2] != "ext") {
			t.Errorf("review row: %q (personal stopped sharing, work shares it)", l)
		}
		if len(f) == 3 && f[0] == "mine" && (f[1] != "account" || f[2] != "-") {
			t.Errorf("mine row: %q", l)
		}
	}

	out := h.runIn(project, "", nil, "ccam", "ext", "work", "--json").stdout
	var rep struct {
		Accounts []struct {
			Name   string `json:"name"`
			Skills []struct{ Name, Source string }
		} `json:"accounts"`
		Problems []struct{ Account, Problem string } `json:"problems"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil || len(rep.Accounts) != 1 || rep.Accounts[0].Name != "work" ||
		len(rep.Accounts[0].Skills) != 1 || rep.Accounts[0].Skills[0].Source != "ext" || len(rep.Problems) != 0 {
		t.Fatalf("json: %v\n%s", err, out)
	}
}
