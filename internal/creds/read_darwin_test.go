package creds

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// A Keychain that does not answer must not hold up the choice of an
// account: Has gives up like Read does.
func TestHasGivesUpOnAHangingKeychain(t *testing.T) {
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "security"), []byte("#!/bin/sh\nexec sleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	start := time.Now()
	if Has(t.TempDir(), t.TempDir()) {
		t.Fatal("no login is stored")
	}
	if took := time.Since(start); took > keychainTimeout+5*time.Second {
		t.Fatalf("took %v", took)
	}
}
