package creds

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeychainService(t *testing.T) {
	// Vector: sha256("/Users/me/.claude-work")[:8] == "1e91dd84".
	if got := KeychainService("/Users/me/.claude-work", "/Users/me"); got != "Claude Code-credentials-1e91dd84" {
		t.Fatalf("got %q", got)
	}
	if got := KeychainService("/Users/me/.claude", "/Users/me"); got != "Claude Code-credentials" {
		t.Fatalf("default store: %q", got)
	}
}

func TestKeychainAccount(t *testing.T) {
	t.Setenv("USER", "dev.name")
	if KeychainAccount() != "dev.name" {
		t.Fatal("plain user")
	}
	t.Setenv("USER", "bad user")
	if KeychainAccount() != "claude-code-user" {
		t.Fatal("unsafe user name")
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := readFile(dir); err != ErrNone {
		t.Fatalf("missing: %v", err)
	}
	os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"x":1}`), 0o600)
	if data, src, err := readFile(dir); err != nil || src != File || string(data) != `{"x":1}` {
		t.Fatalf("%s %s %v", data, src, err)
	}
}
