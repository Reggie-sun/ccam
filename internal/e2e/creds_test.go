package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// writeLogin puts a login file into an account's folder (fake tokens).
func (h *harness) writeLogin(name, data string) string {
	h.t.Helper()
	path := filepath.Join(h.profileDir(name), ".credentials.json")
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		h.t.Fatal(err)
	}
	return path
}

// Claude Code's logout deletes the whole login file, so the account's MCP
// server sign-ins go with the login; a running session of the account is
// signed out too. ccam says so and asks first; with nothing to lose, or
// --yes, it does not ask.
func TestLogoutSaysWhatGoesWithTheLogin(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cred := h.writeLogin("work", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT"},"mcpOAuth":{`+
		`"plugin:cloudflare:cloudflare|aaaa":{"accessToken":"FAKE-AT"},"lin|bbbb":{"accessToken":"","discoveryState":{}}}}`)
	r := h.run("n\n", nil, "ccam", "logout", "work")
	if r.code != 0 || !strings.Contains(r.stdout, "removes work's MCP server sign-ins (plugin:cloudflare:cloudflare)") ||
		strings.Contains(r.stdout, "lin") || !strings.Contains(r.stdout, "Cancelled.") || !fileExists(cred) {
		t.Fatalf("%+v", r)
	}
	if r := h.run("y\n", nil, "ccam", "logout", "work"); r.code != 0 || fileExists(cred) {
		t.Fatalf("%+v", r)
	}
	if _, args := h.called(); strings.Join(args, " ") != "auth logout" {
		t.Fatalf("claude got %v", args)
	}
	if out := h.ok("ccam", "logout", "--help"); !strings.Contains(out, "MCP server sign-ins") {
		t.Fatalf("help: %s", out)
	}

	// personal is running: logging out signs its session out as well.
	claude := exec.Command(fakeBin)
	claude.Env = append(append([]string{}, h.env...), "CLAUDE_CONFIG_DIR="+h.profileDir("personal"),
		"FAKE_RECORD=interactive", "FAKE_SLEEP=60", "FAKE_CLAUDE_LOG="+filepath.Join(h.home, "plain.log"))
	if err := claude.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { claude.Process.Kill(); claude.Wait() })
	waitFile(t, filepath.Join(h.profileDir("personal"), "sessions", strconv.Itoa(claude.Process.Pid)+".json"))
	personal := filepath.Join(h.profileDir("personal"), ".credentials.json")
	if r := h.run("", nil, "ccam", "logout", "personal"); !strings.Contains(r.stdout, "personal is running") || !fileExists(personal) {
		t.Fatalf("%+v", r)
	}
	if r := h.run("", nil, "ccam", "logout", "personal", "--yes"); r.code != 0 || fileExists(personal) {
		t.Fatalf("%+v", r)
	}
}

// A login copied into another account's folder is revoked by the logout
// there too: ccam says so and asks.
func TestLogoutSaysWhenAnotherAccountHoldsTheLogin(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.writeLogin("work", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-SAME"}}`)
	cred := h.writeLogin("personal", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-SAME"}}`)
	r := h.run("n\n", nil, "ccam", "logout", "personal")
	if !strings.Contains(r.stdout, "work holds the same login (a copied folder)") || !strings.Contains(r.stdout, "Cancelled.") || !fileExists(cred) {
		t.Fatalf("%+v", r)
	}
}

// A removed account is logged out first, in its folder before the folder
// moves: Claude Code revokes the token and deletes what it stores for that
// folder, so no working login waits in the trash (or, on macOS, in a
// Keychain item a later account of the name would inherit). An imported
// folder stays live, and an account without a login has nothing to revoke.
func TestRemoveLogsOutFirst(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	// A login copied into another account's folder would be revoked there
	// too: then it is not logged out.
	h.ok("ccam", "add", "copy", "--no-login")
	h.writeLogin("copy", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-SAME"}}`)
	h.writeLogin("personal", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-SAME"}}`)
	os.Remove(h.log)
	if r := h.run("", nil, "ccam", "remove", "copy", "--yes"); r.code != 0 || !strings.Contains(r.stdout, "not logged out, since personal holds the same login") || fileExists(h.log) {
		t.Fatalf("%+v", r)
	}
	r := h.run("", nil, "ccam", "remove", "work", "--yes")
	if r.code != 0 || !strings.Contains(r.stdout, "logged out first") {
		t.Fatalf("%+v", r)
	}
	if dir, args := h.called(); dir != h.profileDir("work") || strings.Join(args, " ") != "auth logout" {
		t.Fatalf("claude got %s %v", dir, args)
	}
	trash, _ := filepath.Glob(filepath.Join(h.home, "data", "trash", "work-*"))
	if len(trash) != 1 || fileExists(filepath.Join(trash[0], ".credentials.json")) || !fileExists(filepath.Join(trash[0], ".claude.json")) {
		t.Fatalf("trash %v", trash)
	}

	h.ok("ccam", "add", "later", "--no-login")
	other := filepath.Join(h.home, "home", ".claude-other")
	os.MkdirAll(other, 0o700)
	os.WriteFile(filepath.Join(other, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-OTHER"}}`), 0o600)
	h.ok("ccam", "import", other)
	for _, name := range []string{"later", "other"} {
		os.Remove(h.log)
		if r := h.run("", nil, "ccam", "remove", name, "--yes"); r.code != 0 || fileExists(h.log) || strings.Contains(r.stdout, "logged out") {
			t.Fatalf("%s: %+v", name, r)
		}
	}
	if !fileExists(filepath.Join(other, ".credentials.json")) {
		t.Fatal("an imported folder was logged out")
	}

	// A logout that fails does not stop the removal, but says so.
	if r := h.run("", []string{"FAKE_EXIT=3"}, "ccam", "remove", "personal", "--yes"); r.code != 0 || !strings.Contains(r.stdout, "could not be logged out first") {
		t.Fatalf("%+v", r)
	}
}

// Claude Code empties both tokens when the server refuses the refresh
// token; such a login is not working, whatever its dates say.
func TestClearedLoginCountsAsExpired(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.writeLogin("work", `{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0,"refreshTokenExpiresAt":4102444800000}}`)
	h.ok("ccam", "personal")
	if r := h.run("", nil, "ccam", "--dry-run", "next"); r.code == 0 || !strings.Contains(r.stderr, "no logged-in account other than personal") {
		t.Fatalf("next went to the cleared login: %+v", r)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "login expired → `ccam login work`") {
		t.Fatalf("doctor:\n%s", out)
	}
	r := h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "login", "--expired")
	if r.code != 0 || !strings.Contains(r.stdout, "[1/1] work (login expired)") {
		t.Fatalf("%+v", r)
	}
}

// An API key (Console) login bills per token: it counts as logged in for
// login --expired, but next, pick and --fallback never move work onto it.
func TestConsoleLoginIsNoAutomaticChoice(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	os.Remove(filepath.Join(h.profileDir("work"), ".credentials.json"))
	os.WriteFile(filepath.Join(h.profileDir("work"), ".claude.json"), []byte(`{"primaryApiKey":"sk-ant-fake"}`), 0o600)
	h.ok("ccam", "personal")
	if r := h.run("", nil, "ccam", "--dry-run", "next"); r.code == 0 || !strings.Contains(r.stderr, "no logged-in account other than personal") {
		t.Fatalf("next went to the Console login: %+v", r)
	}
	runs := filepath.Join(h.home, "runs")
	r := h.run("", []string{"FAKE_RUNS=" + runs, "FAKE_LIMITED=personal"}, "ccam", "--fallback", "personal", "-p", "hi")
	if ran, _ := os.ReadFile(runs); r.code != 75 || strings.Join(strings.Fields(string(ran)), ",") != "personal" || !strings.Contains(r.stderr, "no other logged-in account has room") {
		t.Fatalf("fallback ran %q: %+v", ran, r)
	}
	if r := h.run("", nil, "ccam", "pick"); r.code == 0 || strings.Contains(r.stdout, "work") {
		t.Fatalf("pick chose the Console login: %+v", r)
	}
	if r := h.run("", nil, "ccam", "login", "--expired"); r.code != 0 || strings.Contains(r.stdout, "] work") {
		t.Fatalf("login --expired: %+v", r)
	}
}

// Folders copied from one another (cp -r ~/.claude ~/.claude-2) hold one
// refresh token, and whichever renews it first logs the other out: import
// and doctor say so, without the token.
func TestImportSaysWhenFoldersAreCopies(t *testing.T) {
	h := newHarness(t)
	home := filepath.Join(h.home, "home")
	for _, name := range []string{".claude-one", ".claude-two"} {
		dir := filepath.Join(home, name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"emailAddress":"me@example.com","accountUuid":"u1"}}`), 0o600)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-SAME"}}`), 0o600)
	}
	r := h.run("", nil, "ccam", "import", "--yes")
	if r.code != 0 || !strings.Contains(r.stdout, "holds the same login as one (a copied folder)") || !strings.Contains(r.stdout, "`ccam login two`") ||
		strings.Contains(r.stdout+r.stderr, "FAKE-RT") {
		t.Fatalf("%+v", r)
	}
	if out := h.run("", nil, "ccam", "doctor").stdout; !strings.Contains(out, "same login token as one") || strings.Contains(out, "FAKE-RT") {
		t.Fatalf("doctor:\n%s", out)
	}
}

// remove checks that config.toml can be saved before it logs the account
// out: a removal that then fails to save must not have revoked the login.
// When the folder cannot go to the trash after the logout, the error says
// the account is logged out now.
func TestRemoveLogsOutOnlyWhenItCanFinish(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cred := h.writeLogin("work", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT"}}`)
	config := filepath.Join(h.home, "config", "config.toml")
	saved, _ := os.ReadFile(config)
	os.WriteFile(config, append(append([]byte{}, saved...), "\nrelay = true\n"...), 0o600)
	os.Remove(h.log)
	if r := h.run("", nil, "ccam", "remove", "work", "--yes"); r.code == 0 || !strings.Contains(r.stderr, "not saved") || strings.Contains(r.stderr, "logged out") {
		t.Fatalf("%+v", r)
	}
	if fileExists(h.log) || !fileExists(cred) {
		t.Fatal("logged out although config.toml could not be saved")
	}
	os.WriteFile(config, saved, 0o600)

	if !perms || os.Geteuid() == 0 {
		return
	}
	cfgDir := filepath.Join(h.home, "config")
	os.Chmod(cfgDir, 0o500) // the folder cannot take the saved file
	r := h.run("", nil, "ccam", "remove", "work", "--yes")
	os.Chmod(cfgDir, 0o700)
	if r.code == 0 || fileExists(h.log) || !fileExists(cred) {
		t.Fatalf("read-only config folder: %+v, logged out %v", r, fileExists(h.log))
	}

	trash := filepath.Join(h.home, "data", "trash")
	os.MkdirAll(trash, 0o700)
	os.Chmod(trash, 0o500)
	r = h.run("", nil, "ccam", "remove", "work", "--yes")
	os.Chmod(trash, 0o700)
	if r.code == 0 || !strings.Contains(r.stderr, "work is still registered, but logged out now") {
		t.Fatalf("trash failure: %+v", r)
	}
}
