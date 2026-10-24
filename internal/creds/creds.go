// Package creds reads a profile's Claude Code login without changing it.
//
// Linux keeps it in <config dir>/.credentials.json. macOS keeps it in the
// login Keychain under a name derived from the config dir, and falls back
// to the same file when the Keychain refuses writes. Windows keeps it in the
// file too, unless Claude Code's Credential Manager backend is on, which
// ccam does not read (docs/claude-layout.md).
package creds

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/jhyunwoo/ccam/internal/profile"
)

var (
	// ErrNone means no login is stored.
	ErrNone = errors.New("no stored login")
	// ErrCredentialManager means Claude Code keeps the folder's login in
	// Windows Credential Manager, which ccam does not read: whether it is
	// logged in cannot be told.
	ErrCredentialManager = errors.New("the login is kept in Windows Credential Manager, which ccam cannot read")
)

// Source says where Read found the login.
type Source string

const (
	File              Source = "file"
	Keychain          Source = "keychain"
	CredentialManager Source = "credential-manager"
)

// Stale says why the login stored for dir no longer works (see Judge), or
// "not logged in" when nothing is stored. "" means it still works or cannot
// be judged; an expired access token is renewed by Claude Code. An API key
// login (APIKey) has no Claude login and works.
func Stale(dir, home string, now time.Time) string {
	data, _, err := Read(dir, home)
	why := "not logged in"
	switch {
	case err == nil:
		why = Judge(data, now)
	case !errors.Is(err, ErrNone):
		return ""
	}
	if why == "not logged in" && APIKey(dir) {
		return ""
	}
	return why
}

// Has reports whether a working Claude login is stored for dir: what next,
// pick and --fallback may move work onto. One kept where ccam cannot read
// it (Windows Credential Manager) counts as working. An API key login does
// not count: it bills per token, and the plan leaves it out of the
// automatic choices (Stale still calls it working).
func Has(dir, home string) bool {
	data, _, err := Read(dir, home)
	switch {
	case errors.Is(err, ErrCredentialManager):
		return true
	case err != nil:
		return false
	}
	return Judge(data, time.Now()) == ""
}

// APIKey reports whether dir is logged in with an API key (a Console
// login): Claude Code keeps the key in .claude.json, and the login file, if
// any, holds MCP sign-ins only.
func APIKey(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		return false
	}
	_, ok, _ := profile.TopLevelValue(data, "primaryApiKey")
	return ok
}

// Judge says why stored login data (what Read returns) no longer works, as
// Claude Code 2.1.308 judges it: "not logged in" when it holds no Claude
// login (what a logout leaves, or a file with only MCP sign-ins), "login
// expired" when its refresh token has expired or Claude Code emptied it
// after the server refused it (refreshToken "", which /status shows as
// "Login: Expired"). "" means it works, or it cannot be parsed.
func Judge(data []byte, now time.Time) string {
	var c struct {
		OAuth *struct {
			RefreshToken          *string `json:"refreshToken"`
			RefreshTokenExpiresAt int64   `json:"refreshTokenExpiresAt"`
		} `json:"claudeAiOauth"`
	}
	if json.Unmarshal(data, &c) != nil {
		return ""
	}
	switch o := c.OAuth; {
	case o == nil:
		return "not logged in"
	case o.RefreshTokenExpiresAt > 0 && o.RefreshTokenExpiresAt < now.UnixMilli():
		return "login expired"
	case o.RefreshToken != nil && *o.RefreshToken == "":
		return "login expired"
	}
	return ""
}

// MCPSignIns returns the MCP server sign-ins in stored login data
// (mcpOAuth), keyed as Claude Code keys them: "<server>|<hash of its
// definition>". An entry without an access or refresh token is left out:
// Claude Code keeps those for a server that still needs a sign-in, or whose
// tokens it dropped, and treats them as no sign-in.
func MCPSignIns(data []byte) map[string]json.RawMessage {
	var doc struct {
		MCP map[string]json.RawMessage `json:"mcpOAuth"`
	}
	if json.Unmarshal(data, &doc) != nil {
		return nil
	}
	for k, e := range doc.MCP {
		var t struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}
		if json.Unmarshal(e, &t) != nil || t.AccessToken == "" && t.RefreshToken == "" {
			delete(doc.MCP, k)
		}
	}
	return doc.MCP
}

// LoginKey identifies the login in dir's login file by its refresh token: a
// SHA-256 of it, never the token. Two folders with one key are copies of one
// login, and whichever renews it first logs the other out. "" when the file
// holds no refresh token; a macOS Keychain login is not copied with a
// folder.
func LoginKey(dir string) string {
	data, _, err := readFile(dir)
	if err != nil {
		return ""
	}
	var c struct {
		OAuth *struct {
			RefreshToken string `json:"refreshToken"`
		} `json:"claudeAiOauth"`
	}
	if json.Unmarshal(data, &c) != nil || c.OAuth == nil || c.OAuth.RefreshToken == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(c.OAuth.RefreshToken))
	return hex.EncodeToString(sum[:])
}

// FilePath is the plaintext credentials file of a config dir.
func FilePath(dir string) string { return filepath.Join(dir, ".credentials.json") }

// DefaultStore reports whether Claude Code should use the unsuffixed
// "default" secure-storage entry for dir: true for ~/.claude registered in
// place, so `ccam <name>` and plain `claude` share one macOS Keychain (or
// Windows Credential Manager) login. ccam sets
// CLAUDE_SECURESTORAGE_CONFIG_DIR="" for such profiles.
func DefaultStore(dir, home string) bool {
	return dir == filepath.Join(home, ".claude")
}

// KeychainService is the Keychain item name Claude Code 2.1.308 uses:
// "Claude Code-credentials", plus "-" and the first 8 hex digits of the
// SHA-256 of the config dir in Unicode NFC when CLAUDE_CONFIG_DIR selects it
// (macOS may hand out a folder's name decomposed; nfc normalizes only in
// the macOS build, the one that reads the Keychain).
func KeychainService(dir, home string) string {
	const base = "Claude Code-credentials"
	if DefaultStore(dir, home) {
		return base
	}
	sum := sha256.Sum256([]byte(nfc(dir)))
	return base + "-" + hex.EncodeToString(sum[:])[:8]
}

// accountRe is compiled on first use: a package-level regexp would be
// compiled at init in every ccam process, launches included.
var accountRe = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`^[a-zA-Z0-9._-]+$`) })

// KeychainAccount mirrors Claude Code: $USER when not empty, else Bun's
// os.userInfo().username (Bun runs Claude Code), which is $USER again: ""
// when it is empty, "unknown" when it is unset. A name that is not plain
// becomes a fixed one.
func KeychainAccount() string {
	u, ok := os.LookupEnv("USER")
	if !ok {
		u = "unknown"
	}
	if !accountRe().MatchString(u) {
		return "claude-code-user"
	}
	return u
}

// NamedByFolder names the store where Claude Code keeps dir's login under
// a name made from dir's path, which a move of the folder leaves behind:
// the macOS Keychain (Claude Code writes there first on macOS), or Windows
// Credential Manager when Claude Code uses it for dir. "" when the login is
// the folder's own file and moves with it.
func NamedByFolder(dir string) string {
	switch {
	case runtime.GOOS == "darwin":
		return "the macOS Keychain"
	case runtime.GOOS == "windows" && credentialManager(dir, os.Getenv):
		return "Windows Credential Manager"
	}
	return ""
}

// credentialManager reports whether Claude Code keeps dir's login in Windows
// Credential Manager: with CLAUDE_CODE_FORCE_WINDOWS_CREDMAN=1, or the
// tengu_windows_credman feature flag it caches in the folder's .config.json,
// else in its .claude.json (2.1.304). It moves the login there at its next
// write and deletes .credentials.json.
func credentialManager(dir string, getenv func(string) string) bool {
	if getenv("CLAUDE_CODE_FORCE_WINDOWS_CREDMAN") == "1" {
		return true
	}
	data, err := os.ReadFile(filepath.Join(dir, ".config.json"))
	if err != nil || len(data) == 0 {
		if data, err = os.ReadFile(filepath.Join(dir, ".claude.json")); err != nil {
			return false
		}
	}
	// .claude.json can be megabytes: only the one key is decoded.
	raw, ok, err := profile.TopLevelValue(data, "cachedGrowthBookFeatures")
	if err != nil || !ok {
		return false
	}
	var f struct {
		CredMan bool `json:"tengu_windows_credman"`
	}
	return json.Unmarshal(raw, &f) == nil && f.CredMan
}

func readFile(dir string) ([]byte, Source, error) {
	data, err := os.ReadFile(FilePath(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, "", ErrNone
	}
	return data, File, err
}
