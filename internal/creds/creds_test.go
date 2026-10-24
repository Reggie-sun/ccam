package creds

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestKeychainService(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the Keychain is macOS's; the vectors are Unix paths")
	}
	// Vector: sha256("/Users/me/.claude-work")[:8] == "1e91dd84".
	if got := KeychainService("/Users/me/.claude-work", "/Users/me"); got != "Claude Code-credentials-1e91dd84" {
		t.Fatalf("got %q", got)
	}
	if got := KeychainService("/Users/me/.claude", "/Users/me"); got != "Claude Code-credentials" {
		t.Fatalf("default store: %q", got)
	}
	// Claude Code hashes the path in NFC, so a name macOS hands out
	// decomposed names the same item (values from its own function on Bun).
	dirs := []string{"/Users/me/Documents/caf\u00e9/.claude-work"}
	if runtime.GOOS == "darwin" { // the build that reads the Keychain normalizes
		dirs = append(dirs, "/Users/me/Documents/cafe\u0301/.claude-work")
	}
	for _, dir := range dirs {
		if got := KeychainService(dir, "/Users/me"); got != "Claude Code-credentials-7ef176dd" {
			t.Fatalf("%q: got %q", dir, got)
		}
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
	t.Setenv("USER", "")
	if KeychainAccount() != "claude-code-user" {
		t.Fatal("empty USER")
	}
	// Unset, Claude Code takes Bun's os.userInfo().username, which is
	// "unknown" without USER (not the name in the user database).
	os.Unsetenv("USER")
	if got := KeychainAccount(); got != "unknown" {
		t.Fatalf("USER unset: %q", got)
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

// What Claude Code leaves in the login file decides whether the login still
// works: it empties both tokens when the server refuses the refresh token,
// and a logout leaves a file without a login when it keeps other keys.
func TestJudge(t *testing.T) {
	now := time.Now()
	later := now.Add(20 * 24 * time.Hour).UnixMilli()
	for _, c := range []struct{ name, data, want string }{
		{"works", `{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt","expiresAt":1,"refreshTokenExpiresAt":` + strconv.FormatInt(later, 10) + `}}`, ""},
		{"access token expired", `{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt","expiresAt":1}}`, ""},
		{"no refresh token kept", `{"claudeAiOauth":{"accessToken":"at","expiresAt":4102444800000}}`, ""},
		{"cleared by Claude Code", `{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0,"refreshTokenExpiresAt":` + strconv.FormatInt(later, 10) + `},"mcpOAuth":{}}`, "login expired"},
		{"refresh token emptied", `{"claudeAiOauth":{"accessToken":"at","refreshToken":"","expiresAt":4102444800000}}`, "login expired"},
		{"refresh token expired", `{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt","refreshTokenExpiresAt":2}}`, "login expired"},
		{"left after a logout", `{"coworkRemoteDevice":{"k":"v"}}`, "not logged in"},
		{"only MCP sign-ins", `{"mcpOAuth":{"gh|aa":{"accessToken":"","discoveryState":{}}}}`, "not logged in"},
		{"login null", `{"claudeAiOauth":null}`, "not logged in"},
		{"unreadable", `not json`, ""},
	} {
		if got := Judge([]byte(c.data), now); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

// login --expired, doctor, next, pick and --fallback all ask Stale or Has:
// a login Claude Code has cleared is no working login.
func TestStaleAndHasJudgeTheLogin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS asks the Keychain first")
	}
	now := time.Now()
	for _, c := range []struct {
		data, stale string
		has         bool
	}{
		{"", "not logged in", false},
		{`{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0}}`, "login expired", false},
		{`{"coworkRemoteDevice":{}}`, "not logged in", false},
		{`{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt","expiresAt":1}}`, "", true},
	} {
		dir := t.TempDir()
		if c.data != "" {
			os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(c.data), 0o600)
		}
		if got := Stale(dir, t.TempDir(), now); got != c.stale {
			t.Errorf("%s: Stale %q, want %q", c.data, got, c.stale)
		}
		if got := Has(dir, t.TempDir()); got != c.has {
			t.Errorf("%s: Has %v", c.data, got)
		}
	}
}

func TestMCPSignInsLeaveOutPlaceholders(t *testing.T) {
	got := MCPSignIns([]byte(`{"claudeAiOauth":{"accessToken":"own"},"mcpOAuth":{` +
		`"gh|aaaa":{"serverName":"gh","accessToken":"gh-token","refreshToken":"gh-rt"},` +
		`"lin|bbbb":{"serverName":"lin","accessToken":"","refreshToken":"lin-rt"},` +
		`"plugin:p:srv|cccc":{"serverName":"plugin:p:srv","accessToken":"","discoveryState":{"authorizationServerUrl":"https://x"}},` +
		`"dropped|dddd":{"accessToken":"","expiresAt":0}}}`))
	var keys []string
	for k := range got {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	if !slices.Equal(keys, []string{"gh|aaaa", "lin|bbbb"}) {
		t.Fatalf("got %v", keys)
	}
	if MCPSignIns([]byte(`{"claudeAiOauth":{}}`)) != nil || MCPSignIns([]byte(`nope`)) != nil {
		t.Fatal("no sign-ins")
	}
}

func TestLoginKeyTellsCopiesApart(t *testing.T) {
	write := func(data string) string {
		dir := t.TempDir()
		if data != "" {
			os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(data), 0o600)
		}
		return dir
	}
	a := write(`{"claudeAiOauth":{"accessToken":"at","refreshToken":"FAKE-RT-SAME"}}`)
	b := write(`{"claudeAiOauth":{"accessToken":"other","refreshToken":"FAKE-RT-SAME"},"mcpOAuth":{}}`)
	c := write(`{"claudeAiOauth":{"accessToken":"at","refreshToken":"FAKE-RT-OTHER"}}`)
	if k := LoginKey(a); k == "" || k != LoginKey(b) || k == LoginKey(c) {
		t.Fatalf("keys %q %q %q", LoginKey(a), LoginKey(b), LoginKey(c))
	}
	if len(LoginKey(a)) != 64 || LoginKey(a) == "FAKE-RT-SAME" {
		t.Fatal("the key is a digest, never the token")
	}
	for _, data := range []string{"", `{"claudeAiOauth":{"refreshToken":""}}`, `{"mcpOAuth":{}}`} {
		if k := LoginKey(write(data)); k != "" {
			t.Fatalf("%q: %q", data, k)
		}
	}
}

// Claude Code 2.1.304 on Windows keeps logins in Credential Manager when
// forced by its variable or by the feature flag it caches in the folder.
func TestCredentialManagerFlag(t *testing.T) {
	env := map[string]string{}
	getenv := func(k string) string { return env[k] }
	dir := t.TempDir()
	if credentialManager(dir, getenv) {
		t.Fatal("nothing says so")
	}
	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"cachedGrowthBookFeatures":{"tengu_windows_credman":true},"projects":{}}`), 0o600)
	if !credentialManager(dir, getenv) {
		t.Fatal("flag in .claude.json")
	}
	// .config.json, when not empty, is read instead.
	os.WriteFile(filepath.Join(dir, ".config.json"), []byte(`{"cachedGrowthBookFeatures":{"tengu_windows_credman":false}}`), 0o600)
	if credentialManager(dir, getenv) {
		t.Fatal(".config.json comes first")
	}
	os.WriteFile(filepath.Join(dir, ".config.json"), nil, 0o600)
	if !credentialManager(dir, getenv) {
		t.Fatal("an empty .config.json falls back to .claude.json")
	}
	env["CLAUDE_CODE_FORCE_WINDOWS_CREDMAN"] = "1"
	if !credentialManager(t.TempDir(), getenv) {
		t.Fatal("forced")
	}
}

// Where ccam cannot read the login, it neither calls the account logged
// out nor asks to log it in again.
func TestCredentialManagerLoginIsNotMisreported(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows only")
	}
	t.Setenv("CLAUDE_CODE_FORCE_WINDOWS_CREDMAN", "1")
	dir := t.TempDir()
	if _, src, err := Read(dir, t.TempDir()); !errors.Is(err, ErrCredentialManager) || src != CredentialManager {
		t.Fatalf("%s %v", src, err)
	}
	if !Has(dir, t.TempDir()) || Stale(dir, t.TempDir(), time.Now()) != "" {
		t.Fatal("misreported as logged out")
	}
	// A login file still there is read as before.
	os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"","refreshToken":""}}`), 0o600)
	if Has(dir, t.TempDir()) || Stale(dir, t.TempDir(), time.Now()) != "login expired" {
		t.Fatal("the file")
	}
}

// Where the login store is named after the folder's path, ccam rename keeps
// the folder: always on macOS, on Windows when Credential Manager is used.
func TestNamedByFolder(t *testing.T) {
	dir := t.TempDir()
	want := map[string]string{"darwin": "the macOS Keychain"}[runtime.GOOS]
	if got := NamedByFolder(dir); got != want {
		t.Fatalf("%s: %q", runtime.GOOS, got)
	}
	t.Setenv("CLAUDE_CODE_FORCE_WINDOWS_CREDMAN", "1")
	want = map[string]string{"darwin": "the macOS Keychain", "windows": "Windows Credential Manager"}[runtime.GOOS]
	if got := NamedByFolder(dir); got != want {
		t.Fatalf("%s with Credential Manager: %q", runtime.GOOS, got)
	}
}

// An API key (Console) login keeps its key in .claude.json; once it signs
// in to an MCP server, its login file holds that sign-in only. It still
// works: login --expired leaves it alone. It is no Claude login, though, so
// Has (what the automatic choices go by) is false.
func TestAPIKeyLoginWithAnMCPSignInWorks(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS asks the Keychain first")
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"mcpOAuth":{"gh|abc":{"accessToken":"a","refreshToken":"r"}}}`), 0o600)
	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"primaryApiKey":"sk-ant-x"}`), 0o600)
	if got := Stale(dir, t.TempDir(), time.Now()); got != "" {
		t.Errorf("Stale %q", got)
	}
	if Has(dir, t.TempDir()) {
		t.Error("Has true")
	}
	os.Remove(filepath.Join(dir, ".credentials.json"))
	if got := Stale(dir, t.TempDir(), time.Now()); got != "" || Has(dir, t.TempDir()) {
		t.Errorf("without a login file: Stale %q, Has %v", got, Has(dir, t.TempDir()))
	}
	if !APIKey(dir) {
		t.Error("APIKey false")
	}
}
