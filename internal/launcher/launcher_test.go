package launcher

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
)

func setup(t testing.TB) (paths.Paths, *config.Config) {
	root := t.TempDir()
	os.Setenv("CCAM_HOME", root)
	t.Cleanup(func() { os.Unsetenv("CCAM_HOME") })
	p := paths.FromEnv()
	fake := filepath.Join(root, "fake-claude")
	os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755)
	os.MkdirAll(p.ProfileDir("personal"), 0o700)
	os.MkdirAll(p.ProfileDir("work"), 0o700)
	cfg := &config.Config{Default: "work", ClaudeBin: fake, Profiles: []config.Profile{
		{Name: "personal", Slot: 1},
		{Name: "work", Slot: 2},
	}}
	return p, cfg
}

func TestPrepareSelectsProfile(t *testing.T) {
	p, cfg := setup(t)
	cases := []struct {
		sel     string
		args    []string
		profile string
		argv    []string
	}{
		{"work", []string{"--resume"}, "work", []string{"claude", "--resume"}},
		{"1", []string{"--model", "opus", "--continue"}, "personal", []string{"claude", "--model", "opus", "--continue"}},
		{"", []string{"-p", "résume ça"}, "work", []string{"claude", "-p", "résume ça"}},
		{"work", []string{"--", "--version"}, "work", []string{"claude", "--version"}},
		{"work", []string{"--help"}, "work", []string{"claude", "--help"}},
		{"work", []string{"-p", "corrige les tests — ü", "--output-format", "json"}, "work", []string{"claude", "-p", "corrige les tests — ü", "--output-format", "json"}},
		{"work", []string{"a", "--", "b"}, "work", []string{"claude", "a", "--", "b"}},
	}
	for _, c := range cases {
		pl, err := Prepare(p, cfg, c.sel, c.args)
		if err != nil {
			t.Fatal(err)
		}
		if pl.Profile.Name != c.profile || !reflect.DeepEqual(pl.Argv, c.argv) {
			t.Errorf("%q %q: got %s %q", c.sel, c.args, pl.Profile.Name, pl.Argv)
		}
		if pl.Dir != p.ProfileDir(c.profile) {
			t.Errorf("dir %s", pl.Dir)
		}
	}
}

func TestPrepareUnknownProfile(t *testing.T) {
	p, cfg := setup(t)
	if _, err := Prepare(p, cfg, "nope", nil); err == nil {
		t.Fatal("want error")
	}
}

func TestEnvReplacesConfigDir(t *testing.T) {
	env := buildEnv([]string{"A=1", "CLAUDE_CONFIG_DIR=/old", "CCAM_PROFILE=x", "CLAUDE_SECURESTORAGE_CONFIG_DIR=/y"}, "/new", "work", false)
	want := []string{"A=1", "CLAUDE_CONFIG_DIR=/new", "CCAM_PROFILE=work"}
	if !reflect.DeepEqual(env, want) {
		t.Fatalf("%q", env)
	}
	env = buildEnv(nil, "/h/.claude", "main", true)
	want = []string{"CLAUDE_CONFIG_DIR=/h/.claude", "CCAM_PROFILE=main", "CLAUDE_SECURESTORAGE_CONFIG_DIR="}
	if !reflect.DeepEqual(env, want) {
		t.Fatalf("default store: %q", env)
	}
}

func TestInsertArgsGoesBeforeUserArgs(t *testing.T) {
	pl := &Plan{Argv: []string{"claude", "--mcp-config", "mine.json"}}
	pl.InsertArgs("--mcp-config", "group.json")
	want := []string{"claude", "--mcp-config", "group.json", "--mcp-config", "mine.json"}
	if !reflect.DeepEqual(pl.Argv, want) {
		t.Fatalf("%q", pl.Argv)
	}
}

func TestDryRunQuotes(t *testing.T) {
	pl := &Plan{Dir: "/d", Bin: "/bin/claude", Argv: []string{"claude", "-p", "it's ünïcode", "--x"}}
	var buf bytes.Buffer
	pl.DryRun(&buf)
	want := `CLAUDE_CONFIG_DIR=/d /bin/claude -p 'it'\''s ünïcode' --x` + "\n"
	if buf.String() != want {
		t.Fatalf("got %s", buf.String())
	}
}

func TestResolveClaudeCachesLookup(t *testing.T) {
	p, _ := setup(t)
	binDir := t.TempDir()
	fake := filepath.Join(binDir, "claude-x")
	os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755)
	t.Setenv("PATH", binDir)
	got, err := ResolveClaude(p, "claude-x")
	if err != nil || got != fake {
		t.Fatalf("%q %v", got, err)
	}
	if data, _ := os.ReadFile(filepath.Join(p.State, "claude-path")); !strings.HasPrefix(string(data), fake+"\t") {
		t.Fatalf("cache %q", data)
	}
	// Cache hit works even with an empty PATH.
	t.Setenv("PATH", "")
	if got, err := ResolveClaude(p, "claude-x"); err != nil || got != fake {
		t.Fatalf("cache miss: %q %v", got, err)
	}
	// A changed file invalidates the cache.
	os.Chtimes(fake, time.Now().Add(time.Hour), time.Now().Add(time.Hour))
	if _, err := ResolveClaude(p, "claude-x"); err == nil {
		t.Fatal("stale cache should trigger a lookup, which fails with empty PATH")
	}
}

func TestLastUsed(t *testing.T) {
	p, _ := setup(t)
	if !LastUsed(p, "work").IsZero() {
		t.Fatal("expected zero")
	}
	MarkUsed(p, "work")
	if time.Since(LastUsed(p, "work")) > time.Minute {
		t.Fatal("not recorded")
	}
}

func BenchmarkLaunchPath(b *testing.B) {
	p, cfg := setup(b)
	for i := range cfg.Profiles {
		cfg.Profiles[i].Share = config.DefaultShare()
	}
	if err := cfg.Save(p.ConfigFile()); err != nil {
		b.Fatal(err)
	}
	// First run creates the share links; the loop measures the steady state.
	if err := Launch(p, "work", []string{"--continue"}, true, io.Discard); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Launch(p, "work", []string{"--continue"}, true, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}
