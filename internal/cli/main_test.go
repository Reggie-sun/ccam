package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// stubClaude is a claude that does nothing: on Unix a script that exits 0,
// on Windows a claude.exe that is no program, so starting it fails. TestMain
// puts its folder first on PATH, so that a test whose claude_bin is empty
// (or names "claude") never finds the Claude Code installed on the machine:
// on Windows, where there is no `true`, tests did run the real one.
var stubClaude string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ccam-cli-test")
	if err != nil {
		panic(err)
	}
	stubClaude = writeStubClaude(dir)
	os.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	// A launch that is no dry run would put the stub claude in place of this
	// test process, which exits 0: go test would report ok and skip every
	// test after it. Reaching exec fails the run instead.
	launcher.BeforeExec = func() {
		panic("a test reached syscall.Exec, which would replace the test process with claude: use --dry-run, or a claude_bin that cannot be found")
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func writeStubClaude(dir string) string {
	path, data := filepath.Join(dir, "claude"), "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		path, data = filepath.Join(dir, "claude.exe"), "not a program\n"
	}
	if err := os.WriteFile(path, []byte(data), 0o755); err != nil {
		panic(err)
	}
	return path
}

// No test can reach the claude installed on the machine through PATH.
func TestClaudeOnPathIsTheStub(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	got, err := launcher.ResolveClaude(paths.FromEnv(), "")
	if err != nil || got != stubClaude {
		t.Fatalf("claude on PATH is %q (%v), not the stub %q", got, err, stubClaude)
	}
}
