package transcript

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSlug(t *testing.T) {
	long := "/home/me/" + strings.Repeat("very-long-folder-name/", 12) + "project"
	cases := map[string]string{
		"/home/jhyunwoo/workspace/ccam": "-home-jhyunwoo-workspace-ccam",
		"/a/b_c.d e":                    "-a-b-c-d-e",
		"/x/한":                          "-x--",
		"/x/😀":                          "-x---",
		`C:\Users\me\code`:              "C--Users-me-code",
		// Claude Code 2.1.304 cuts at 200 and adds the path's hash.
		long: "-home-me-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folder-name-very-long-folde-g15bg1",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func write(t *testing.T, path, data string, mod time.Time) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	os.Chtimes(path, mod, mod)
}

func TestNewestSizeAndTitle(t *testing.T) {
	cfg := t.TempDir()
	cwd := t.TempDir()
	dir := ProjectDir(cfg, cwd)
	now := time.Now()
	write(t, filepath.Join(dir, "old.jsonl"), "{}\n", now.Add(-time.Hour))
	conv := strings.Join([]string{
		`{"type":"user","message":{"role":"user","content":"hi"}}`,
		`{"type":"assistant","message":{"usage":{"input_tokens":5,"cache_read_input_tokens":100,"cache_creation_input_tokens":20,"output_tokens":7}}}`,
		`{"type":"ai-title","aiTitle":"Fix the tests"}`,
		`{"type":"assistant","isSidechain":true,"message":{"usage":{"input_tokens":999999}}}`,
		`{"type":"system","subtype":"turn_duration"}`,
	}, "\n") + "\n"
	write(t, filepath.Join(dir, "abc-123.jsonl"), conv, now)
	os.MkdirAll(filepath.Join(dir, "abc-123", "subagents"), 0o700) // folders are not conversations
	path, _, ok := Newest(dir)
	if !ok || ID(path) != "abc-123" {
		t.Fatalf("%q %v", path, ok)
	}
	if n := Size(path); n != 132 {
		t.Fatalf("size %d (the subagent's reply is not the conversation's)", n)
	}
	if got := Title(path); got != "Fix the tests" {
		t.Fatal(got)
	}
	write(t, path, conv+`{"type":"custom-title","customTitle":"release prep"}`+"\n", now)
	if got := Title(path); got != "release prep" {
		t.Fatal(got)
	}
	if _, _, ok := Newest(filepath.Join(cfg, "projects", "nothing")); ok {
		t.Fatal("no folder")
	}
}
