package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/session"
)

type env struct {
	t       *testing.T
	f       *Fetcher
	root    string
	hits    atomic.Int32
	refresh atomic.Int32
}

// newEnv starts a fake usage API. The token decides the behaviour:
// "ok-<fixture>" → 200 with that fixture, "slow" → hangs, "bad" → 401, "boom" → 500.
func newEnv(t *testing.T) *env {
	e := &env{t: t, root: t.TempDir()}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			e.refresh.Add(1)
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if body["grant_type"] != "refresh_token" || body["refresh_token"] != "rt-old" {
				w.WriteHeader(400)
				return
			}
			fmt.Fprint(w, `{"access_token":"ok-new-format.json","refresh_token":"rt-new","expires_in":3600}`)
			return
		}
		e.hits.Add(1)
		if r.Header.Get("anthropic-beta") != "oauth-2025-04-20" {
			w.WriteHeader(400)
			return
		}
		tok := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		switch {
		case strings.HasPrefix(tok, "ok-"):
			http.ServeFile(w, r, filepath.Join("testdata", strings.TrimPrefix(tok, "ok-")))
		case tok == "slow":
			select {
			case <-r.Context().Done():
			case <-time.After(10 * time.Second):
			}
		case tok == "weird":
			fmt.Fprint(w, `{"brand_new_shape":true}`)
		case tok == "bad":
			w.WriteHeader(401)
		case tok == "limited":
			w.Header().Set("Retry-After", "120")
			w.WriteHeader(429)
		default:
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(srv.Close)
	e.f = NewFetcher(filepath.Join(e.root, "state"))
	e.f.URL = srv.URL + "/usage"
	e.f.TokenURL = srv.URL + "/token"
	return e
}

func (e *env) account(name, token string, expiresAt int64) Account {
	dir := filepath.Join(e.root, name)
	os.MkdirAll(dir, 0o700)
	if token != "" {
		cred := fmt.Sprintf(`{"claudeAiOauth":{"accessToken":%q,"refreshToken":"rt-old","expiresAt":%d},"mcpOAuth":{"x":"keep"}}`, token, expiresAt)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(cred), 0o600)
	}
	return Account{Name: name, Dir: dir, SessionDir: filepath.Join(e.root, "run", name)}
}

var future = time.Now().Add(time.Hour).UnixMilli()

func TestFetchAllParallelWithOneHungAccount(t *testing.T) {
	e := newEnv(t)
	accts := []Account{
		e.account("a", "ok-new-format.json", future),
		e.account("b", "ok-old-format.json", future),
		e.account("c", "slow", future),
		e.account("d", "ok-scoped-only.json", future),
		e.account("e", "ok-new-format.json", future),
	}
	var arrived []string
	start := time.Now()
	reports := e.f.FetchAll(context.Background(), accts, Options{}, func(i int, r Report) {
		arrived = append(arrived, r.Name)
	})
	took := time.Since(start)
	if took > accountTimeout+500*time.Millisecond {
		t.Fatalf("took %v", took)
	}
	if arrived[len(arrived)-1] != "c" {
		t.Fatalf("hung account should arrive last: %v", arrived)
	}
	for _, r := range reports {
		if r.Name == "c" {
			if r.Status != StatusError || r.Limits != nil {
				t.Fatalf("hung: %+v", r)
			}
			continue
		}
		if r.Status != StatusOK || r.Limits == nil || r.Source != "live" {
			t.Fatalf("%s: %+v", r.Name, r)
		}
	}
	if reports[1].Limits.FiveHour.Percent != 100 {
		t.Fatal("order must follow input")
	}
}

func TestFetchUsesFreshCacheThenFallsBack(t *testing.T) {
	e := newEnv(t)
	// Before the fixture's 5-hour reset (2026-10-18T18:10Z).
	base := time.Date(2026, 10, 18, 17, 0, 0, 0, time.UTC)
	e.f.Now = func() time.Time { return base }
	a := e.account("a", "ok-new-format.json", future)
	ctx := context.Background()
	e.f.Fetch(ctx, a, Options{})
	e.f.Fetch(ctx, a, Options{})
	if e.hits.Load() != 1 {
		t.Fatalf("second call within 60s should hit the cache, hits=%d", e.hits.Load())
	}
	e.f.Fetch(ctx, a, Options{NoCache: true})
	if e.hits.Load() != 2 {
		t.Fatal("--no-cache must query")
	}
	// Server now fails: old values come back marked as cache.
	e.f.Now = func() time.Time { return base.Add(10 * time.Minute) }
	a = e.account("a", "boom", future)
	r := e.f.Fetch(ctx, a, Options{})
	if r.Status != StatusError || r.Source != "cache" || r.Limits == nil || r.Limits.FiveHour.Percent != 72 {
		t.Fatalf("fallback: %+v", r)
	}
	// Once the 5-hour window has reset, the saved 72% no longer applies;
	// the weekly window (resets 10-22) still does.
	e.f.Now = func() time.Time { return base.Add(2 * time.Hour) }
	r = e.f.Fetch(ctx, a, Options{})
	if l := r.Limits; l == nil || !l.FiveHour.Known || l.FiveHour.Percent != 0 || !l.FiveHour.ResetsAt.IsZero() || l.Weekly.Percent != 41 {
		t.Fatalf("lapsed window: %+v", l)
	}
}

func TestFetchStatusLineCacheAndSourceCache(t *testing.T) {
	e := newEnv(t)
	a := e.account("a", "boom", future)
	SaveStatusLine(e.f.StateDir, "a", Limits{FiveHour: Window{Percent: 23, Known: true}}, time.Now())
	r := e.f.Fetch(context.Background(), a, Options{})
	if r.Source != "statusline" || r.Limits.FiveHour.Percent != 23 {
		t.Fatalf("%+v", r)
	}
	hits := e.hits.Load()
	r = e.f.Fetch(context.Background(), a, Options{Source: "cache"})
	if e.hits.Load() != hits || r.Status != StatusOK || r.Source != "statusline" {
		t.Fatalf("--source cache must not query: %+v", r)
	}
	r = e.f.Fetch(context.Background(), e.account("z", "boom", future), Options{Source: "cache"})
	if r.Status != StatusError {
		t.Fatalf("no cache: %+v", r)
	}
}

func TestFetchExpiredLoggedOutAndAPI(t *testing.T) {
	e := newEnv(t)
	past := time.Now().Add(-time.Hour).UnixMilli()
	r := e.f.Fetch(context.Background(), e.account("x", "ok-new-format.json", past), Options{})
	if r.Status != StatusExpired || e.hits.Load() != 0 {
		t.Fatalf("expired must not query: %+v", r)
	}
	if r := e.f.Fetch(context.Background(), e.account("bad", "bad", future), Options{}); r.Status != StatusExpired {
		t.Fatalf("401: %+v", r)
	}
	if r := e.f.Fetch(context.Background(), e.account("none", "", 0), Options{}); r.Status != StatusLoggedOut {
		t.Fatalf("no creds: %+v", r)
	}
	api := e.account("api", "", 0)
	os.WriteFile(filepath.Join(api.Dir, ".claude.json"), []byte(`{"primaryApiKey":"sk-x"}`), 0o600)
	if r := e.f.Fetch(context.Background(), api, Options{}); r.Status != StatusAPI {
		t.Fatalf("api: %+v", r)
	}
}

func TestRefreshOnlyWhenNotRunning(t *testing.T) {
	e := newEnv(t)
	past := time.Now().Add(-time.Hour).UnixMilli()
	a := e.account("x", "stale", past)

	// This test process stands in for a running claude.
	if err := session.Register(a.SessionDir); err != nil {
		t.Fatal(err)
	}
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusExpired || e.refresh.Load() != 0 || !strings.Contains(r.Err, "running") {
		t.Fatalf("refresh while running: %+v", r)
	}
	os.RemoveAll(a.SessionDir)

	r = e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusOK || e.refresh.Load() != 1 {
		t.Fatalf("refresh: %+v", r)
	}
	data, _ := os.ReadFile(filepath.Join(a.Dir, ".credentials.json"))
	var doc struct {
		OAuth struct {
			AccessToken, RefreshToken string
			ExpiresAt                 int64
		} `json:"claudeAiOauth"`
		MCP map[string]string `json:"mcpOAuth"`
	}
	json.Unmarshal(data, &doc)
	if doc.OAuth.AccessToken != "ok-new-format.json" || doc.OAuth.RefreshToken != "rt-new" || doc.OAuth.ExpiresAt < future-60000 || doc.MCP["x"] != "keep" {
		t.Fatalf("credentials after refresh: %s", data)
	}
	if info, _ := os.Stat(filepath.Join(a.Dir, ".credentials.json")); runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatal("perm")
	}
}

func TestRateLimitedLookupWaitsBeforeAskingAgain(t *testing.T) {
	e := newEnv(t)
	base := time.Now() // the token in e.account expires an hour from now
	e.f.Now = func() time.Time { return base }
	a := e.account("a", "limited", future)
	r := e.f.Fetch(context.Background(), a, Options{NoCache: true})
	if r.Status != StatusError || !strings.Contains(r.Err, "rate limited") || e.hits.Load() != 1 {
		t.Fatalf("first: %+v hits=%d", r, e.hits.Load())
	}
	// Within Retry-After, even --no-cache does not ask again.
	e.f.Now = func() time.Time { return base.Add(time.Minute) }
	r = e.f.Fetch(context.Background(), a, Options{NoCache: true})
	if e.hits.Load() != 1 || !strings.Contains(r.Err, "rate limited") {
		t.Fatalf("during the pause: %+v hits=%d", r, e.hits.Load())
	}
	e.f.Now = func() time.Time { return base.Add(3 * time.Minute) }
	e.f.Fetch(context.Background(), a, Options{NoCache: true})
	if e.hits.Load() != 2 {
		t.Fatalf("after the pause: hits=%d", e.hits.Load())
	}
}

func TestRetryAfter(t *testing.T) {
	at := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	for in, want := range map[string]time.Duration{
		"":                              defaultPause,
		"30":                            30 * time.Second,
		"Tue, 20 Oct 2026 12:10:00 GMT": 10 * time.Minute,
		"soon":                          defaultPause,
		"999999":                        maxPause,
	} {
		if got := retryAfter(in, at); got != want {
			t.Errorf("%q: %v", in, got)
		}
	}
}

func TestFormatChangeIsFlagged(t *testing.T) {
	e := newEnv(t)
	r := e.f.Fetch(context.Background(), e.account("w", "weird", future), Options{})
	if !r.FormatChanged || r.Status != StatusError {
		t.Fatalf("%+v", r)
	}
	if r := e.f.Fetch(context.Background(), e.account("b", "boom", future), Options{}); r.FormatChanged {
		t.Fatal("HTTP 500 is not a format change")
	}
}

func TestRefreshRateLimitedSaysWait(t *testing.T) {
	e := newEnv(t)
	limited := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	t.Cleanup(limited.Close)
	e.f.TokenURL = limited.URL
	a := e.account("x", "stale", time.Now().Add(-time.Hour).UnixMilli())
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusExpired || !strings.Contains(r.Err, "try again later") || strings.Contains(r.Err, "ccam login") {
		t.Fatalf("%+v", r)
	}
	data, _ := os.ReadFile(filepath.Join(a.Dir, ".credentials.json"))
	if !strings.Contains(string(data), `"rt-old"`) {
		t.Fatalf("credentials must be untouched: %s", data)
	}
}

func TestTokenURLIsClaudeCodes(t *testing.T) {
	// Claude Code 2.1.304 refreshes at this address (TOKEN_URL in its bundle).
	if TokenURL != "https://platform.claude.com/v1/oauth/token" {
		t.Fatal(TokenURL)
	}
}

func TestRefreshRefusedForImportedProfile(t *testing.T) {
	e := newEnv(t)
	a := e.account("x", "stale", time.Now().Add(-time.Hour).UnixMilli())
	a.Imported = true
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusExpired || e.refresh.Load() != 0 || !strings.Contains(r.Err, "claude") {
		t.Fatalf("%+v", r)
	}
}
