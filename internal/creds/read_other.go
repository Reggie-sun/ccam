//go:build !darwin

package creds

import (
	"errors"
	"os"
	"runtime"
)

// nfc leaves s as it is: no Keychain is read here (see KeychainService).
func nfc(s string) string { return s }

// Read returns the stored login JSON from <dir>/.credentials.json. On
// Windows, a missing file with Claude Code's Credential Manager backend on
// gives ErrCredentialManager: the login may be there.
func Read(dir, home string) ([]byte, Source, error) {
	data, src, err := readFile(dir)
	if errors.Is(err, ErrNone) && runtime.GOOS == "windows" && credentialManager(dir, os.Getenv) {
		return nil, CredentialManager, ErrCredentialManager
	}
	return data, src, err
}
