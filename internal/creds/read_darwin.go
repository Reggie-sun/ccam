package creds

import (
	"context"
	"os/exec"
	"time"
)

// Read returns the stored login JSON, Keychain first, then the fallback file.
func Read(dir, home string) ([]byte, Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "security", "find-generic-password",
		"-a", KeychainAccount(), "-w", "-s", KeychainService(dir, home)).Output()
	if err == nil && len(out) > 0 {
		return out, Keychain, nil
	}
	return readFile(dir)
}

// Has reports whether a login is stored, without reading the secret.
func Has(dir, home string) bool {
	err := exec.Command("security", "find-generic-password",
		"-a", KeychainAccount(), "-s", KeychainService(dir, home)).Run()
	if err == nil {
		return true
	}
	_, _, err = readFile(dir)
	return err == nil
}
