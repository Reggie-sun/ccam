package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/mcp"
)

// mcp import must not silently replace the definition every other account
// of the group uses.
func TestMCPImportKeepsDifferentGroupServer(t *testing.T) {
	c := newTestCLI(t, nil)
	cfg, _ := config.Load(c.p.ConfigFile())
	if err := cfg.AddGroup("default", []string{"mcp"}, []string{"personal", "work", "side"}); err != nil {
		t.Fatal(err)
	}
	cfg.Save(c.p.ConfigFile())
	store := mcp.Path(c.p.GroupDir("default"))
	os.MkdirAll(c.p.GroupDir("default"), 0o700)
	os.WriteFile(store, []byte(`{"mcpServers":{"github":{"command":"gh","env":{"TOKEN":"A"}}}}`), 0o600)
	os.WriteFile(filepath.Join(c.p.ProfileDir("work"), ".claude.json"),
		[]byte(`{"mcpServers":{"github":{"command":"gh","env":{"TOKEN":"B"}},"fs":{"command":"fs"}}}`), 0o600)

	out, errOut, code := c.run("", "mcp", "import", "work")
	if code != 0 || !strings.Contains(out, "Skipped github") || !strings.Contains(out, "Copied to group default: fs") {
		t.Fatalf("%d\n%s%s", code, out, errOut)
	}
	if b, _ := os.ReadFile(store); !strings.Contains(string(b), `"A"`) || !strings.Contains(string(b), `"fs"`) {
		t.Fatalf("group store: %s", b)
	}
	if b, _ := os.ReadFile(filepath.Join(c.p.ProfileDir("work"), ".claude.json")); !strings.Contains(string(b), `"B"`) {
		t.Fatalf("skipped server must stay in the account: %s", b)
	}

	out, _, code = c.run("", "mcp", "import", "work", "github", "--force")
	if code != 0 || !strings.Contains(out, "previous mcp.json is saved") {
		t.Fatalf("--force: %d %s", code, out)
	}
	if b, _ := os.ReadFile(store); !strings.Contains(string(b), `"B"`) {
		t.Fatalf("--force did not replace: %s", b)
	}
	saved, _ := filepath.Glob(filepath.Join(c.p.GroupDir("default"), ".ccam-backup", "*", "mcp.json"))
	if len(saved) != 1 {
		t.Fatalf("backup: %v", saved)
	}
	if b, _ := os.ReadFile(saved[0]); !strings.Contains(string(b), `"A"`) {
		t.Fatalf("backup content: %s", b)
	}
}

func TestMCPGroupFlagIsValidated(t *testing.T) {
	c := newTestCLI(t, nil)
	_, errOut, code := c.run("", "mcp", "add", "x", "--group", "../../escape", "--", "echo")
	if code == 0 || !strings.Contains(errOut, "invalid group name") {
		t.Fatalf("%d %q", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(c.p.Data, "..", "escape")); err == nil {
		t.Fatal("wrote outside the groups folder")
	}
}
