package creds

import (
	"context"
	"os/exec"
	"time"

	"golang.org/x/text/unicode/norm"
)

// nfc is s in Unicode NFC, as Claude Code hashes a config dir's path.
func nfc(s string) string { return norm.NFC.String(s) }

// keychainTimeout bounds one call of the security command, which can hang
// on a locked or busy Keychain.
const keychainTimeout = 3 * time.Second

// Read returns the stored login JSON, Keychain first, then the fallback file.
func Read(dir, home string) ([]byte, Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), keychainTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, "security", "find-generic-password",
		"-a", KeychainAccount(), "-w", "-s", KeychainService(dir, home)).Output()
	if err == nil && len(out) > 0 {
		return out, Keychain, nil
	}
	return readFile(dir)
}
