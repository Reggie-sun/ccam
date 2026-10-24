package cli

import (
	"bytes"
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// Where files cannot be linked (Windows without Developer Mode), the first
// `ccam add` still works: group default shares only folders and MCP
// servers. A group that shares files is refused before anything is made.
func TestAddWhereFilesCannotBeLinked(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	c := &testCLI{t: t, p: paths.FromEnv()}
	orig := canLinkFiles
	canLinkFiles = func(string) error {
		return errors.New("files can be linked on Windows only in Developer Mode (Settings > System > For developers)")
	}
	t.Cleanup(func() { canLinkFiles = orig })

	out := c.ok("add", "work", "--no-login")
	if !strings.Contains(out, "group default shares only mcp, extensions") || !strings.Contains(out, "`ccam group set default --add settings,memory`") {
		t.Fatalf("no note:\n%s", out)
	}
	g, ok := c.cfg().FindGroup("default")
	if !ok || !g.HasMember("work") || !slices.Equal(g.Items, []string{"mcp", "extensions"}) {
		t.Fatalf("group default: %+v", g)
	}
	for _, name := range []string{"settings.json", "CLAUDE.md"} {
		if _, ok := fsutil.LinkTarget(filepath.Join(c.p.ProfileDir("work"), name)); ok {
			t.Fatalf("%s is linked", name)
		}
	}
	if _, ok := fsutil.LinkTarget(filepath.Join(c.p.ProfileDir("work"), "skills")); !ok {
		t.Fatal("skills is not linked")
	}
	if out := c.ok("add", "home", "--no-login"); strings.Contains(out, "Note:") {
		t.Fatalf("joining that group needs no note:\n%s", out)
	}

	c.ok("group", "add", "docs", "--items", "memory")
	_, errOut, code := c.run("", "add", "side", "--no-login", "--group", "docs")
	if code == 0 || !strings.Contains(errOut, "Developer Mode") || !strings.Contains(errOut, "group docs shares memory") {
		t.Fatalf("%d %s", code, errOut)
	}
	if fsutil.Exists(c.p.ProfileDir("side")) {
		t.Fatal("a refused add left its folder")
	}
	if _, ok := c.cfg().Find("side"); ok {
		t.Fatal("side was registered")
	}
}

// Where files can be linked, group default shares what it always did.
func TestAddWhereFilesCanBeLinked(t *testing.T) {
	if err := fsutil.CanLinkFiles(t.TempDir()); err != nil {
		t.Skip(err)
	}
	t.Setenv("CCAM_HOME", t.TempDir())
	c := &testCLI{t: t, p: paths.FromEnv()}
	if out := c.ok("add", "work", "--no-login"); strings.Contains(out, "Note:") {
		t.Fatal(out)
	}
	if g, _ := c.cfg().FindGroup("default"); !slices.Equal(g.Items, []string{"mcp", "settings", "memory", "extensions"}) {
		t.Fatalf("%+v", g)
	}
}

// The wizard may pick only items that link files: group default then
// shares nothing yet, and the Note says so in plain words.
func TestAddWhereNoPickedItemCanBeLinked(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	c := &testCLI{t: t, p: paths.FromEnv()}
	orig := canLinkFiles
	canLinkFiles = func(string) error { return errors.New("files can be linked on Windows only in Developer Mode") }
	t.Cleanup(func() { canLinkFiles = orig })
	e, _ := c.env()
	var out bytes.Buffer
	e.stdout = &out
	if err := addProfile(e, newAccount{name: "work", groups: []string{"default"}, items: []string{"settings", "memory"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "so group default shares nothing yet. Once it is on, `ccam group set default --add settings,memory` shares them.") {
		t.Fatalf("note:\n%s", out.String())
	}
	if g, ok := c.cfg().FindGroup("default"); !ok || len(g.Items) != 0 || !g.HasMember("work") {
		t.Fatalf("%+v", g)
	}
}
