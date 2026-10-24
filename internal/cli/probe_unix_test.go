//go:build !windows

package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// A probe given up (its time over, Ctrl+C in --watch) asks Claude Code to
// end: on SIGTERM it finishes saving a token renewal it has started, which a
// SIGKILL would cut off.
func TestProbeLetsClaudeCodeFinish(t *testing.T) {
	c := newTestCLI(t, nil)
	dir := t.TempDir()
	saved := filepath.Join(dir, "saved")
	fake := filepath.Join(dir, "claude")
	script := "#!/bin/sh\ntrap 'sleep 1; echo renewed > \"" + saved + "\"; exit 143' TERM\nwhile :; do sleep 0.1; done\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	start := time.Now()
	claudeProbe(&env{paths: c.p}, &config.Config{ClaudeBin: fake})(ctx, usage.Account{Name: "work", Dir: c.p.ProfileDir("work")})
	if !fileThere(saved) {
		t.Fatalf("Claude Code was stopped before it could finish (after %v)", time.Since(start))
	}
	if took := time.Since(start); took > probeExitWait+2*time.Second {
		t.Fatalf("took %v", took)
	}
}
