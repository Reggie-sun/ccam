//go:build !windows

package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/usage"
)

// Ctrl+C while `ccam usage --refresh` waits for the token server ends the
// lookups, not ccam: the server may have rotated the refresh token already,
// and its answer is saved before ccam exits.
func TestUsageRefreshOutlivesCtrlC(t *testing.T) {
	c := newTestCLI(t, nil)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		syscall.Kill(os.Getpid(), syscall.SIGINT)
		time.Sleep(300 * time.Millisecond)
		fmt.Fprint(w, `{"access_token":"pct-5","refresh_token":"FAKE-RT-NEW","expires_in":3600}`)
	}))
	t.Cleanup(srv.Close)
	orig := newFetcher
	newFetcher = func(dir string) *usage.Fetcher {
		f := orig(dir)
		f.TokenURL = srv.URL
		return f
	}
	t.Cleanup(func() { newFetcher = orig })
	cred := filepath.Join(c.p.ProfileDir("work"), ".credentials.json")
	os.WriteFile(cred, []byte(`{"claudeAiOauth":{"accessToken":"stale","refreshToken":"FAKE-RT-OLD","expiresAt":1}}`), 0o600)
	out, _, code := c.run("", "usage", "--refresh")
	if data, _ := os.ReadFile(cred); !strings.Contains(string(data), "FAKE-RT-NEW") {
		t.Fatalf("the renewed login was not saved: %s", data)
	}
	// It was interrupted, and says so: no timeout is blamed.
	if code != 130 || strings.Contains(out, "no response") {
		t.Fatalf("exit %d:\n%s", code, out)
	}
}
