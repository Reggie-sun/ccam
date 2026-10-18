// Package creds reads a profile's Claude Code login without changing it.
//
// Linux keeps it in <config dir>/.credentials.json. macOS keeps it in the
// login Keychain under a name derived from the config dir, and falls back
// to the same file when the Keychain refuses writes (docs/claude-layout.md).
package creds

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// ErrNone means no login is stored.
var ErrNone = errors.New("no stored login")

// Source says where Read found the login.
type Source string

const (
	File     Source = "file"
	Keychain Source = "keychain"
)

// FilePath is the plaintext credentials file of a config dir.
func FilePath(dir string) string { return filepath.Join(dir, ".credentials.json") }

// DefaultStore reports whether Claude Code should use the unsuffixed
// "default" secure-storage entry for dir: true for ~/.claude registered in
// place, so `ccs <name>` and plain `claude` share one macOS Keychain login.
// ccam sets CLAUDE_SECURESTORAGE_CONFIG_DIR="" for such profiles.
func DefaultStore(dir, home string) bool {
	return dir == filepath.Join(home, ".claude")
}

// KeychainService is the Keychain item name Claude Code 2.1.302 uses:
// "Claude Code-credentials", plus "-" and the first 8 hex digits of the
// SHA-256 of the config dir when CLAUDE_CONFIG_DIR selects it.
func KeychainService(dir, home string) string {
	const base = "Claude Code-credentials"
	if DefaultStore(dir, home) {
		return base
	}
	sum := sha256.Sum256([]byte(dir))
	return base + "-" + hex.EncodeToString(sum[:])[:8]
}

var accountRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// KeychainAccount mirrors Claude Code: $USER, or a fixed name when unusable.
func KeychainAccount() string {
	u := os.Getenv("USER")
	if !accountRe.MatchString(u) {
		return "claude-code-user"
	}
	return u
}

func readFile(dir string) ([]byte, Source, error) {
	data, err := os.ReadFile(FilePath(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, "", ErrNone
	}
	return data, File, err
}
