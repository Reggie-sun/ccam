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
	for _, path := range append(scripts, filepath.Join(root, "install.sh")) {
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
	r = h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "add", "work", "--share", "history=team")
	outputs = append(outputs, r.stdout+r.stderr)
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

	// Seed the status line cache so usage has numbers to show.
	input, err := os.ReadFile(filepath.Join(repoRoot(), "internal", "usage", "testdata", "statusline-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	r = h.run(string(input), []string{"CCAM_PROFILE=work"}, "ccam", "statusline")
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
		{"ccam", "--dry-run", "share", "work", "--history", "off", "--yes"},
		{"ccam", "--dry-run", "share", "personal", "--history", "team", "--prefer", "merge"},
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
		// Error paths.
		{"ccam", "nobody"},
		{"ccam", "9"},
		{"ccam", "add", "list"},
		{"ccam", "add", "Bad"},
		{"ccam", "add", "work"},
		{"ccam", "add", "x", "--share", "history"},
		{"ccam", "share", "work", "--mcp", "bad/group"},
		{"ccam", "share", "work", "--history", "off", "--prefer", "nope"},
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
	} {
		run("", argv...)
	}
	os.WriteFile(filepath.Join(h.home, "home", ".claude", "settings.json"), []byte(`{"theme":"dark"}`), 0o600)
	run("2\nn\n", "ccam", "migrate", "--settings")               // conflict prompt, then no
	os.Symlink(ccamBin, filepath.Join(h.home, "bin", "claude7")) // stale shortcut name
	run("", "claude7")
	run("n\n", "ccam", "remove", "spare")
	run("y\n", "ccam", "remove", "spare")
	run("", "ccam", "work", "-p", "hello")

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
