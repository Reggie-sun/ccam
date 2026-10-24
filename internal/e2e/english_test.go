package e2e

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"
)

// All program output is English. unicode.Hangul covers the Jamo
// (U+1100–U+11FF), compatibility Jamo (U+3130–U+318F), Jamo extended
// (U+A960–U+A97F, U+D7B0–U+D7FF) and syllable (U+AC00–U+D7AF) blocks.
func hasHangul(s string) bool {
	return strings.ContainsFunc(s, func(r rune) bool { return unicode.Is(unicode.Hangul, r) })
}

// hangulLines returns "file:line: text" for every line with Hangul.
func hangulLines(path, text string) []string {
	var out []string
	for i, line := range strings.Split(text, "\n") {
		if hasHangul(line) {
			out = append(out, path+":"+strconv.Itoa(i+1)+": "+strings.TrimSpace(line))
		}
	}
	return out
}

// TestSourceHasNoHangul fails on Hangul anywhere in non-test Go source
// under cmd/ and internal/ (string and rune literals, also when written as
// escapes, comments and identifiers) and in the shell scripts users run.
func TestSourceHasNoHangul(t *testing.T) {
	root := repoRoot()
	var bad []string
	for _, dir := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return err
			}
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, path)
			bad = append(bad, hangulLines(rel, string(src))...)
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
			if err != nil {
				return err
			}
			ast.Inspect(f, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || (lit.Kind != token.STRING && lit.Kind != token.CHAR) {
					return true
				}
				if v, err := strconv.Unquote(lit.Value); err == nil && hasHangul(v) && !hasHangul(lit.Value) {
					bad = append(bad, rel+":"+strconv.Itoa(fset.Position(lit.Pos()).Line)+": escaped Hangul in "+lit.Value)
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	scripts, _ := filepath.Glob(filepath.Join(root, "scripts", "*.sh"))
	for _, path := range append(scripts, filepath.Join(root, "install.sh"), filepath.Join(root, "install.ps1")) {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		rel, _ := filepath.Rel(root, path)
		bad = append(bad, hangulLines(rel, string(src))...)
	}
	if len(bad) > 0 {
		t.Fatalf("program output must be English; Hangul found:\n%s", strings.Join(bad, "\n"))
	}
}

var subcommandRe = regexp.MustCompile(`(?m)^  (\S+)\s{2,}\S`)

// subcommands lists the "Available Commands" of a help text.
func subcommands(help string) []string {
	_, after, ok := strings.Cut(help, "Available Commands:\n")
	if !ok {
		return nil
	}
	section, _, _ := strings.Cut(after, "\n\n")
	var out []string
	for _, m := range subcommandRe.FindAllStringSubmatch(section+"\n", -1) {
		out = append(out, m[1])
	}
	return out
}

// TestOutputHasNoHangul runs the built binary through help for every
// command, the common commands, dry runs, prompts and error paths, and
// fails if stdout or stderr contains Hangul.
func TestOutputHasNoHangul(t *testing.T) {
	h := newHarness(t)
	h.seedClaude() // an existing ~/.claude: add's migrate hint, import, migrate
	var outputs []string
	run := func(stdin string, argv ...string) {
		r := h.run(stdin, nil, argv...)
		outputs = append(outputs, "$ "+strings.Join(argv, " ")+"\n"+r.stdout+r.stderr)
	}

	// Empty state, then the accounts every other check uses.
	run("", "ccam")
	run("", "ccam", "list")
	run("", "ccam", "usage")
	run("", "ccam", "work")
	r := h.run("", []string{"FAKE_EMAIL=me@example.com"}, "ccam", "add", "personal")
	outputs = append(outputs, r.stdout+r.stderr)
	r = h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "add", "work")
	outputs = append(outputs, r.stdout+r.stderr)
	run("", "ccam", "group", "add", "team", "--items", "history", "--members", "work", "--yes")
	run("", "ccam", "add", "later", "--no-login")

	// Help for every command, found by walking "Available Commands".
	queue := [][]string{{}}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		r := h.run("", nil, append(append([]string{"ccam"}, path...), "--help")...)
		outputs = append(outputs, r.stdout+r.stderr)
		run("", append([]string{"ccam", "help"}, path...)...)
		for _, sub := range subcommands(r.stdout) {
			queue = append(queue, append(append([]string{}, path...), sub))
		}
	}
	if len(outputs) < 40 {
		t.Fatalf("help walk found too few commands (%d outputs)", len(outputs))
	}

	// Seed the status line cache so usage has numbers to show, with the
	// fixture's fixed reset times moved ahead of the clock (a window past
	// its reset shows as 0%).
	raw, err := os.ReadFile(filepath.Join(repoRoot(), "internal", "usage", "testdata", "statusline-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	input := strings.NewReplacer("1792346400", strconv.FormatInt(time.Now().Add(2*time.Hour).Unix(), 10),
		"1792634400", strconv.FormatInt(time.Now().Add(72*time.Hour).Unix(), 10)).Replace(string(raw))
	if input == string(raw) {
		t.Fatal("the fixture's resets_at changed; update this test")
	}
	r = h.run(input, []string{"CCAM_PROFILE=work"}, "ccam", "statusline")
	outputs = append(outputs, r.stdout+r.stderr)

	for _, argv := range [][]string{
		{"ccam"},
		{"ccam", "list"},
		{"ccam", "list", "--json"},
		{"ccam", "version"},
		{"ccam", "doctor"},
		{"ccam", "usage"},
		{"ccam", "usage", "--source", "cache", "-v"},
		{"ccam", "usage", "--json", "--source", "cache"},
		{"ccam", "--dry-run", "work", "-c"},
		{"ccam", "--dry-run", "next"},
		{"ccam", "--dry-run", "next", "--continue"},
		{"ccam", "--dry-run", "next", "--auto"},
		{"ccam", "--dry-run", "group", "leave", "team", "work", "--yes"},
		{"ccam", "--dry-run", "group", "join", "team", "personal", "--prefer", "merge"},
		{"ccam", "--dry-run", "group", "set", "default", "--remove", "memory"},
		{"ccam", "group"},
		{"ccam", "group", "show", "default"},
		{"ccam", "share", "work"},
		{"ccam", "use", "work"},
		{"ccam", "mcp", "list"},
		{"ccam", "mcp", "add", "fs", "--", "npx", "server"},
		{"ccam", "mcp", "add-json", "web", `{"type":"http","url":"https://example.com/mcp"}`},
		{"ccam", "mcp", "list"},
		{"ccam", "mcp", "remove", "web"},
		{"ccam", "mcp", "import", "work"},
		{"ccam", "shortcut", "work"},
		{"ccam", "--dry-run", "shortcut", "work", "ccwork"},
		{"ccam", "shortcut", "work", "ccwk"},
		{"ccam", "shortcut", "work", "ccjob"},
		{"ccam", "shortcut", "work"},
		{"ccam", "doctor"},
		{"ccam", "shims", "install"},
		{"ccam", "shortcut", "work", "--remove"},
		{"ccam", "shortcut", "work", "--remove"},
		{"ccam", "statusline", "install"},
		{"ccam", "import"},
		{"ccam", "--dry-run", "migrate"},
		{"ccam", "migrate", "--yes", "--history"},
		{"ccam", "migrate", "--yes"},
		{"ccam", "rename", "later", "spare"},
		{"ccam", "update"}, // a development build refuses
		{"ccam", "update", "--check"},
		// Error paths.
		{"ccam", "nobody"},
		{"ccam", "9"},
		{"ccam", "add", "list"},
		{"ccam", "add", "Bad"},
		{"ccam", "add", "work"},
		{"ccam", "add", "x", "--share", "history"},
		{"ccam", "add", "x", "--group", "nope"},
		{"ccam", "share", "work", "--mcp", "bad/group"},
		{"ccam", "group", "join", "team", "personal", "--prefer", "nope"},
		{"ccam", "group", "add", "local", "--items", "mcp"},
		{"ccam", "group", "add", "x", "--items", "tokens"},
		{"ccam", "group", "set", "default", "--add", "history"},
		{"ccam", "group", "remove", "default"},
		{"ccam", "group", "show", "nope"},
		{"ccam", "mcp", "add", "nothing"},
		{"ccam", "mcp", "add", "h", "--url", "https://x", "--header", "nocolon"},
		{"ccam", "mcp", "add-json", "bad", "[]"},
		{"ccam", "mcp", "remove", "nothere"},
		{"ccam", "usage", "--source", "bogus"},
		{"ccam", "completion", "fish"},
		{"ccam", "import", "/nonexistent"},
		{"ccam", "migrate", "/nonexistent"},
		{"ccam", "shortcut", "work", "Bad"},
		{"ccam", "shortcut", "work", "x", "--remove"},
		{"ccam", "add", "x2", "--no-login", "--shortcut", "list"},
		{"ccam", "migrate", "--group", "bad/group"},
		{"ccam", "use"},
		{"ccam", "bogus-command", "--bogus-flag"},
		{"ccam", "list", "--bogus-flag"},
		{"ccam", "update", "--version", "latest"},
		{"ccam", "update", "--check", "--version", "v1.0.0"},
	} {
		run("", argv...)
	}
	os.WriteFile(filepath.Join(h.home, "home", ".claude", "settings.json"), []byte(`{"theme":"dark"}`), 0o600)
	run("2\nn\n", "ccam", "migrate", "--settings")                  // conflict prompt, then no
	os.Symlink(ccamBin, filepath.Join(h.home, "bin", "ccgone"+exe)) // link no account claims
	run("", "ccgone")
	run("n\n", "ccam", "remove", "spare")
	run("y\n", "ccam", "remove", "spare")
	run("", "ccam", "work", "-p", "hello")
	r = h.run("", []string{"CCAM_RELEASES_URL=http://127.0.0.1:1"}, "ccam", "update", "--version", "v9.9.9") // offline
	outputs = append(outputs, r.stdout+r.stderr)

	var bad []string
	for _, o := range outputs {
		if hasHangul(o) {
			bad = append(bad, o)
		}
	}
	if len(bad) > 0 {
		t.Fatalf("program output must be English; Hangul in:\n%s", strings.Join(bad, "\n---\n"))
	}
}
