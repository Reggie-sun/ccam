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
)

// rotatingServer is a token endpoint that spends the refresh token it is
// sent as soon as the request arrives, as a rotating OAuth server does, and
// then runs during (if set) before it answers with rt-new. calls counts the
// requests.
func rotatingServer(t *testing.T, during func(r *http.Request)) (*httptest.Server, *atomic.Int32) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["refresh_token"] != "rt-old" {
			w.WriteHeader(400)
			return
		}
		if during != nil {
			during(r)
		}
		fmt.Fprint(w, `{"access_token":"ok-new-format.json","refresh_token":"rt-new","expires_in":3600}`)
	}))
	t.Cleanup(srv.Close)
	return srv, &calls
}

// stored reads the login file back.
func stored(t *testing.T, a Account) (oauth map[string]any, mcp map[string]any) {
	t.Helper()
	var doc struct {
		OAuth map[string]any `json:"claudeAiOauth"`
		MCP   map[string]any `json:"mcpOAuth"`
	}
	data, err := os.ReadFile(filepath.Join(a.Dir, ".credentials.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("%v: %s", err, data)
	}
	return doc.OAuth, doc.MCP
}

func expiredAccount(e *env, name string) Account {
	return e.account(name, "stale", time.Now().Add(-time.Hour).UnixMilli())
}

// A token server slower than a usage lookup's 3 s has still renewed the
// login: its answer is the only refresh token that works now, so ccam waits
// for it as Claude Code does.
func TestRefreshWaitsForASlowTokenServer(t *testing.T) {
	e := newEnv(t)
	srv, _ := rotatingServer(t, func(*http.Request) { time.Sleep(accountTimeout + 500*time.Millisecond) })
	e.f.TokenURL = srv.URL
	a := expiredAccount(e, "x")
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if oauth, _ := stored(t, a); r.Status != StatusOK || oauth["refreshToken"] != "rt-new" {
		t.Fatalf("%+v, stored %v", r, oauth)
	}
}

// Ctrl+C in `ccam usage --refresh --watch` cancels the lookups, but not a
// renewal already sent; a lookup cancelled before sends none.
func TestRefreshOutlivesACancel(t *testing.T) {
	e := newEnv(t)
	ctx, cancel := context.WithCancel(context.Background())
	srv, calls := rotatingServer(t, func(*http.Request) {
		cancel()
		time.Sleep(300 * time.Millisecond)
	})
	e.f.TokenURL = srv.URL
	a := expiredAccount(e, "x")
	e.f.Fetch(ctx, a, Options{Refresh: true})
	if oauth, _ := stored(t, a); oauth["refreshToken"] != "rt-new" {
		t.Fatalf("the renewed login was dropped: %v", oauth)
	}
	b := expiredAccount(e, "y")
	if r := e.f.Fetch(ctx, b, Options{Refresh: true}); calls.Load() != 1 || r.Status != StatusExpired {
		t.Fatalf("a cancelled lookup sent a renewal: %d %+v", calls.Load(), r)
	}
}

// Claude Code writes the login file while ccam waits for the token server
// (an MCP sign-in, its own renewal, a logout): ccam changes only the tokens
// of the file as it is then, keeps a login renewed meanwhile, and does not
// undo a logout.
func TestRefreshKeepsWhatWasWrittenMeanwhile(t *testing.T) {
	e := newEnv(t)
	var path string
	write := func(data string) func(*http.Request) {
		return func(*http.Request) { os.WriteFile(path, []byte(data), 0o600) }
	}
	cases := []struct {
		name, meanwhile string
		check           func(r Report, oauth, mcp map[string]any) bool
	}{
		{"MCP sign-in added", `{"claudeAiOauth":{"accessToken":"stale","refreshToken":"rt-old","expiresAt":1},"mcpOAuth":{"x":"keep","linear|abcd":{"accessToken":"lin"}}}`,
			func(r Report, oauth, mcp map[string]any) bool {
				return r.Status == StatusOK && oauth["refreshToken"] == "rt-new" && mcp["x"] == "keep" && mcp["linear|abcd"] != nil
			}},
		{"renewed by Claude Code", fmt.Sprintf(`{"claudeAiOauth":{"accessToken":"ok-old-format.json","refreshToken":"rt-cc","expiresAt":%d}}`, future),
			func(r Report, oauth, mcp map[string]any) bool {
				return r.Status == StatusOK && oauth["refreshToken"] == "rt-cc" && oauth["accessToken"] == "ok-old-format.json"
			}},
		{"cleared by Claude Code", `{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0},"mcpOAuth":{"x":"keep"}}`,
			func(r Report, oauth, mcp map[string]any) bool {
				return r.Status == StatusOK && oauth["refreshToken"] == "rt-new"
			}},
		{"logged out", `{"coworkRemoteDevice":{"k":"v"}}`,
			func(r Report, oauth, mcp map[string]any) bool {
				return r.Status == StatusExpired && oauth == nil && strings.Contains(r.Err, "logged out")
			}},
		{"login null", `{"claudeAiOauth":null}`,
			func(r Report, oauth, mcp map[string]any) bool { return r.Status == StatusExpired && oauth == nil }},
	}
	for _, c := range cases {
		a := expiredAccount(e, strings.ReplaceAll(c.name, " ", "-"))
		path = filepath.Join(a.Dir, ".credentials.json")
		srv, _ := rotatingServer(t, write(c.meanwhile))
		e.f.TokenURL = srv.URL
		r := e.f.Fetch(context.Background(), a, Options{Refresh: true, NoCache: true})
		if oauth, mcp := stored(t, a); !c.check(r, oauth, mcp) {
			t.Errorf("%s: %+v; stored %v %v", c.name, r, oauth, mcp)
		}
	}
	// A file gone meanwhile stays gone.
	a := expiredAccount(e, "gone")
	path = filepath.Join(a.Dir, ".credentials.json")
	srv, _ := rotatingServer(t, func(*http.Request) { os.Remove(path) })
	e.f.TokenURL = srv.URL
	if r := e.f.Fetch(context.Background(), a, Options{Refresh: true}); r.Status != StatusExpired || fileExists(path) {
		t.Fatalf("%+v", r)
	}
}

// Claude Code rewrites the login file under .storage-write.lock, reading it
// first: ccam waits for that lock before it reads the file and writes, or
// Claude Code's write would bring the spent refresh token back.
func TestRefreshWaitsForClaudeCodesWriteLock(t *testing.T) {
	e := newEnv(t)
	a := expiredAccount(e, "x")
	path := filepath.Join(a.Dir, ".credentials.json")
	lockDir := filepath.Join(a.Dir, ".storage-write.lock")
	done := make(chan struct{})
	srv, _ := rotatingServer(t, func(*http.Request) {
		os.Mkdir(lockDir, 0o700)
		data, _ := os.ReadFile(path) // read under the lock, written back later
		go func() {
			defer close(done)
			time.Sleep(300 * time.Millisecond)
			var doc map[string]any
			json.Unmarshal(data, &doc)
			doc["mcpOAuth"] = map[string]any{"x": "keep", "linear|abcd": map[string]any{"accessToken": "lin"}}
			out, _ := json.Marshal(doc)
			os.WriteFile(path, out, 0o600)
			os.Remove(lockDir)
		}()
	})
	e.f.TokenURL = srv.URL
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	<-done
	if oauth, mcp := stored(t, a); r.Status != StatusOK || oauth["refreshToken"] != "rt-new" || mcp["linear|abcd"] == nil {
		t.Fatalf("%+v; stored %v %v", r, oauth, mcp)
	}
	if fileExists(lockDir) {
		t.Fatal("ccam left Claude Code's write lock behind")
	}
}

// While Claude Code renews the login (its .oauth_refresh.lock, or the older
// <folder>.lock, is held), ccam does not send the same refresh token; a lock
// a minute old is stale, as for Claude Code.
func TestRefreshLeavesClaudeCodesRenewalToIt(t *testing.T) {
	e := newEnv(t)
	srv, calls := rotatingServer(t, nil)
	e.f.TokenURL = srv.URL
	a := expiredAccount(e, "x")
	for _, held := range []string{filepath.Join(a.Dir, ".oauth_refresh.lock"), a.Dir + ".lock"} {
		os.Mkdir(held, 0o700)
		r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
		if calls.Load() != 0 || r.Status != StatusExpired || !strings.Contains(r.Err, "being renewed by Claude Code") {
			t.Fatalf("%s held: %d %+v", filepath.Base(held), calls.Load(), r)
		}
		old := time.Now().Add(-2 * time.Minute)
		os.Chtimes(held, old, old)
	}
	if r := e.f.Fetch(context.Background(), a, Options{Refresh: true}); calls.Load() != 1 || r.Status != StatusOK {
		t.Fatalf("stale locks: %d %+v", calls.Load(), r)
	}
	for _, l := range []string{filepath.Join(a.Dir, ".oauth_refresh.lock"), a.Dir + ".lock"} {
		if fileExists(l) {
			t.Fatalf("%s left behind", l)
		}
	}
}

// A renewal sent whose answer never came may have rotated the token: the
// message says how to get the login back. One that never reached the
// server says nothing of the kind.
func TestRefreshSaysWhenTheLoginMayBeLost(t *testing.T) {
	e := newEnv(t)
	hangUp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err == nil {
			conn.Close()
		}
	}))
	t.Cleanup(hangUp.Close)
	e.f.TokenURL = hangUp.URL
	a := expiredAccount(e, "x")
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusExpired || !strings.HasPrefix(r.Err, "login may be lost: ") || !strings.Contains(r.Err, "run `ccam login x`") {
		t.Fatalf("%+v", r)
	}
	if out := strings.Join(Table([]*Report{&r}, RenderOptions{Now: time.Now()}), "\n"); !strings.Contains(out, "ccam login x") {
		t.Fatalf("the table hides it: %s", out)
	}
	// Also above saved numbers, which most accounts have.
	e.f.writeCache("x", "live", cached{At: time.Now().Add(-time.Hour), Limits: &Limits{FiveHour: Window{Known: true, Percent: 40}}})
	r = e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if out := strings.Join(Table([]*Report{&r}, RenderOptions{Now: time.Now()}), "\n"); r.Limits == nil || !strings.Contains(out, "login may be lost: ") || strings.Contains(out, "token not renewed") || !strings.Contains(out, "ccam login x") || !strings.Contains(out, "40%") {
		t.Fatalf("with saved numbers: %s", out)
	}
	closed := httptest.NewServer(http.NotFoundHandler())
	e.f.TokenURL = closed.URL
	closed.Close()
	if r := e.f.Fetch(context.Background(), a, Options{Refresh: true}); r.Status != StatusExpired || strings.Contains(r.Err, "ccam login") || !strings.HasPrefix(r.Err, "token not renewed: ") {
		t.Fatalf("not sent: %+v", r)
	}
}

// The token endpoint may give the renewed refresh token a new lifetime;
// Claude Code then stores it, and a stale date would later make the login
// look expired. Without one in the answer the old date stays.
func TestRefreshStoresTheNewRefreshTokenLifetime(t *testing.T) {
	e := newEnv(t)
	answer := `{"access_token":"ok-new-format.json","refresh_token":"rt-new","expires_in":3600,"refresh_token_expires_in":2592000}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, answer) }))
	t.Cleanup(srv.Close)
	e.f.TokenURL = srv.URL
	a := expiredAccount(e, "x")
	path := filepath.Join(a.Dir, ".credentials.json")
	os.WriteFile(path, fmt.Appendf(nil, `{"claudeAiOauth":{"accessToken":"stale","refreshToken":"rt-old","expiresAt":1,"refreshTokenExpiresAt":%d}}`, future), 0o600)
	if r := e.f.Fetch(context.Background(), a, Options{Refresh: true}); r.Status != StatusOK {
		t.Fatalf("%+v", r)
	}
	oauth, _ := stored(t, a)
	if at, _ := oauth["refreshTokenExpiresAt"].(float64); at < float64(time.Now().Add(29*24*time.Hour).UnixMilli()) {
		t.Fatalf("stored %v", oauth)
	}
	answer = `{"access_token":"ok-new-format.json","refresh_token":"rt-newer","expires_in":3600.0}`
	os.WriteFile(path, []byte(`{"claudeAiOauth":{"accessToken":"stale","refreshToken":"rt-old","expiresAt":1,"refreshTokenExpiresAt":4102444800000}}`), 0o600)
	if r := e.f.Fetch(context.Background(), a, Options{Refresh: true, NoCache: true}); r.Status != StatusOK {
		t.Fatalf("%+v", r)
	}
	if oauth, _ := stored(t, a); oauth["refreshTokenExpiresAt"] != float64(4102444800000) || oauth["refreshToken"] != "rt-newer" {
		t.Fatalf("stored %v", oauth)
	}
}

// A login made with another OAuth client is renewed with that client, as
// Claude Code does; others with Claude Code's.
func TestRefreshSendsTheLoginsClient(t *testing.T) {
	e := newEnv(t)
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		got = append(got, body["client_id"])
		fmt.Fprint(w, `{"access_token":"ok-new-format.json","refresh_token":"rt-new","expires_in":3600}`)
	}))
	t.Cleanup(srv.Close)
	e.f.TokenURL = srv.URL
	a := expiredAccount(e, "own")
	b := expiredAccount(e, "other")
	os.WriteFile(filepath.Join(b.Dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"stale","refreshToken":"rt-old","expiresAt":1,"clientId":"other-client"}}`), 0o600)
	e.f.Fetch(context.Background(), a, Options{Refresh: true})
	e.f.Fetch(context.Background(), b, Options{Refresh: true})
	if len(got) != 2 || got[0] != clientID || got[1] != "other-client" {
		t.Fatalf("%v", got)
	}
}

// A login without a refresh token written while ccam waits for the answer
// stays, as Claude Code's save rule keeps it.
func TestRefreshKeepsALoginWithoutARefreshToken(t *testing.T) {
	e := newEnv(t)
	a := expiredAccount(e, "x")
	path := filepath.Join(a.Dir, ".credentials.json")
	srv, _ := rotatingServer(t, func(*http.Request) {
		os.WriteFile(path, fmt.Appendf(nil, `{"claudeAiOauth":{"accessToken":"ok-old-format.json","expiresAt":%d}}`, future), 0o600)
	})
	e.f.TokenURL = srv.URL
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if oauth, _ := stored(t, a); r.Status != StatusOK || oauth["accessToken"] != "ok-old-format.json" || oauth["refreshToken"] != nil {
		t.Fatalf("%+v %v", r, oauth)
	}
}

// When the renewed login cannot be written, the message says the old one
// may be gone.
func TestRefreshSaysWhenTheLoginCannotBeSaved(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("needs a folder this user cannot write")
	}
	e := newEnv(t)
	a := expiredAccount(e, "x")
	srv, _ := rotatingServer(t, func(*http.Request) { os.Chmod(a.Dir, 0o500) })
	t.Cleanup(func() { os.Chmod(a.Dir, 0o700) })
	e.f.TokenURL = srv.URL
	r := e.f.Fetch(context.Background(), a, Options{Refresh: true})
	if r.Status != StatusExpired || !strings.HasPrefix(r.Err, "login may be lost: ") || !strings.Contains(r.Err, "could not be saved") || !strings.Contains(r.Err, "ccam login x") {
		t.Fatalf("%+v", r)
	}
}

// A login renewed by someone else between the expiry ccam saw and the
// lock it then took is used as it is: sending the refresh token again
// would rotate it a second time and could log the account out.
func TestRefreshUsesALoginRenewedMeanwhile(t *testing.T) {
	e := newEnv(t)
	srv, calls := rotatingServer(t, nil)
	e.f.TokenURL = srv.URL
	a := e.account("x", "ok-renewed", time.Now().Add(time.Hour).UnixMilli())
	got, err := e.f.refresh(context.Background(), a, 0)
	if err != nil || calls.Load() != 0 || got.AccessToken != "ok-renewed" {
		t.Fatalf("%+v %v, %d token requests", got, err, calls.Load())
	}
}
