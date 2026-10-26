package usage

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/session"
)

// TokenURL and clientID are what Claude Code uses to refresh its OAuth
// token (TOKEN_URL and CLIENT_ID in the 2.1.304 bundle). Not verified with a
// real refresh by ccam's tests (doing so would rotate a real refresh token);
// see docs/claude-layout.md.
const (
	TokenURL = "https://platform.claude.com/v1/oauth/token"
	clientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
)

// refreshTimeout bounds a token renewal once sent: Claude Code's own
// (timeout:30000 on its token request, 2.1.308).
const refreshTimeout = 30 * time.Second

// Claude Code's locks in the account folder (2.1.308): it renews the login
// under .oauth_refresh.lock and, as older versions did, <realpath(dir)>.lock,
// both mkdir locks it counts stale after a minute, and rewrites
// .credentials.json under .storage-write.lock (stale after 15 s, which it
// waits for up to about 7 s).
const (
	refreshLockStale = time.Minute
	storageLockStale = 15 * time.Second
	storageLockWait  = 5 * time.Second
)

// RefreshLock is the lock refresh holds, beside an account's session
// folder; rename, remove and logout take it too.
func RefreshLock(sessionDir string) string { return sessionDir + ".refresh.lock" }

// refresh renews an expired access token. Refresh tokens rotate, so doing
// this while Claude Code runs on the profile would log it out: it is
// refused while a session ccam started, or any Claude Code process that
// lists itself in the folder (a background session, its daemon), runs, and
// while one holds the lock it renews the login under.
func (f *Fetcher) refresh(ctx context.Context, a Account) (oauthCred, error) {
	if a.Imported {
		return oauthCred{}, fmt.Errorf("not refreshing an account registered in place, since plain claude may also be using it")
	}
	// Serialize ccam's own refreshes of this account; everything below is
	// checked under the lock so two ccams never rotate the token twice.
	l, err := lock.TryFlock(RefreshLock(a.SessionDir))
	if errors.Is(err, lock.ErrBusy) {
		return oauthCred{}, fmt.Errorf("another ccam process is using this login right now (a refresh, rename, remove or logout)")
	}
	if err != nil {
		return oauthCred{}, err
	}
	defer l.Unlock()
	if session.InUse(a.SessionDir, a.Dir) {
		return oauthCred{}, fmt.Errorf("a session is running")
	}
	release, err := claudeRefreshLock(a.Dir)
	if err != nil {
		return oauthCred{}, err
	}
	defer release()

	path := creds.FilePath(a.Dir)
	data, err := os.ReadFile(path)
	if err != nil {
		return oauthCred{}, err
	}
	current, ok := storedOAuth(data)
	if !ok {
		return oauthCred{}, fmt.Errorf("the login file holds no login")
	}
	if current.AccessToken != "" && !current.expired(f.Now()) {
		return current, nil // refreshed by someone else meanwhile
	}
	rt := current.RefreshToken
	if rt == "" {
		return oauthCred{}, fmt.Errorf("no refresh token")
	}
	// A lookup already given up (Ctrl+C in --watch) starts no renewal.
	if err := ctx.Err(); err != nil {
		return oauthCred{}, err
	}
	// Once sent, the request is never given up early, not even on Ctrl+C:
	// the server may rotate the refresh token as it reads it, and then only
	// its answer still holds the login.
	sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), refreshTimeout)
	defer cancel()
	var sent atomic.Bool
	sendCtx = httptrace.WithClientTrace(sendCtx, &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) { sent.Store(true) }})
	// The client the login was made with, as Claude Code sends it.
	body, _ := json.Marshal(map[string]string{"grant_type": "refresh_token", "refresh_token": rt, "client_id": cmp.Or(current.ClientID, clientID)})
	req, err := http.NewRequestWithContext(sendCtx, http.MethodPost, f.TokenURL, bytes.NewReader(body))
	if err != nil {
		return oauthCred{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cmp.Or(f.TokenClient, f.Client).Do(req)
	if err != nil {
		if sent.Load() {
			return oauthCred{}, lostRenewal(a.Name, err)
		}
		return oauthCred{}, fmt.Errorf("token refresh failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return oauthCred{}, fmt.Errorf("token refresh rate limited (HTTP 429); try again later")
	}
	if resp.StatusCode != http.StatusOK {
		return oauthCred{}, fmt.Errorf("token refresh failed (HTTP %d) → `ccam login %s`", resp.StatusCode, a.Name)
	}
	// Numbers as floats: an answer that cannot be read is a lost login.
	var tok struct {
		AccessToken      string   `json:"access_token"`
		RefreshToken     string   `json:"refresh_token"`
		ExpiresIn        float64  `json:"expires_in"`
		RefreshExpiresIn *float64 `json:"refresh_token_expires_in"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&tok); err != nil || tok.AccessToken == "" {
		return oauthCred{}, lostRenewal(a.Name, errors.New("its answer could not be read"))
	}
	// Without a new refresh token, or its lifetime, in the answer the old
	// one stays, as in Claude Code.
	now := f.Now().UnixMilli()
	got := oauthCred{AccessToken: tok.AccessToken, RefreshToken: cmp.Or(tok.RefreshToken, rt), ExpiresAt: now + int64(tok.ExpiresIn*1000)}
	if tok.RefreshExpiresIn != nil {
		got.refreshExpiresAt = now + int64(*tok.RefreshExpiresIn*1000)
	}
	return saveRenewed(path, rt, got, a.Name)
}

// lostRenewal is a renewal that was sent but whose answer ccam did not get:
// the server may have rotated the refresh token, which then no longer works.
func lostRenewal(name string, err error) error {
	return fmt.Errorf("%w: the token refresh got no usable answer (%v); the server may have renewed the login without ccam saving it: if `ccam %s` then asks you to log in, run `ccam login %s`", errLoginLost, err, name, name)
}

// claudeRefreshLock takes the locks Claude Code renews dir's login under. It
// does not wait: a held lock means a Claude Code (one ccam cannot see, such
// as one on another machine sharing the folder) is renewing it now, and
// spending one refresh token twice logs one of the two out.
func claudeRefreshLock(dir string) (func(), error) {
	busy := func(err error) error {
		if errors.Is(err, lock.ErrBusy) {
			return errors.New("its login is being renewed by Claude Code right now; try again in a minute")
		}
		return err
	}
	unlock, err := lock.DirLockStale(filepath.Join(dir, ".oauth_refresh.lock"), 0, refreshLockStale)
	if err != nil {
		return nil, busy(err)
	}
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		real = dir
	}
	unlockLegacy, err := lock.DirLockStale(real+".lock", 0, refreshLockStale)
	if err != nil {
		unlock()
		return nil, busy(err)
	}
	return func() { unlockLegacy(); unlock() }, nil
}

// storedOAuth reads claudeAiOauth out of login file data; ok is false when
// the file holds no login (none, or null after a logout).
func storedOAuth(data []byte) (oauthCred, bool) {
	var doc struct {
		OAuth *oauthCred `json:"claudeAiOauth"`
	}
	if json.Unmarshal(data, &doc) != nil || doc.OAuth == nil {
		return oauthCred{}, false
	}
	return *doc.OAuth, true
}

// errLoggedOut means the login file no longer holds a login.
var errLoggedOut = errors.New("logged out")

// errLoginLost starts the errors of a renewal the server may have made
// without ccam saving it: the old refresh token is likely spent.
var errLoginLost = errors.New("login may be lost")

// saveRenewed puts a renewed login into the login file as it is now, under
// Claude Code's write lock: only the tokens of claudeAiOauth change, so what
// Claude Code wrote meanwhile (an MCP sign-in) stays. When the file holds
// another refresh token by then (Claude Code renewed it too), that login
// stands, as Claude Code decides; when the login is gone (a logout), nothing
// is written.
func saveRenewed(path, posted string, got oauthCred, name string) (oauthCred, error) {
	// The server has rotated the refresh token already: past the wait ccam
	// writes without the lock rather than drop the only one that works.
	if unlock, err := lock.DirLockStale(filepath.Join(filepath.Dir(path), ".storage-write.lock"), storageLockWait, storageLockStale); err == nil {
		defer unlock()
	}
	var err error
	for try := range 3 {
		if try > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		var stored oauthCred
		var out []byte
		stored, out, err = merged(path, posted, got)
		switch {
		case errors.Is(err, errLoggedOut):
			return oauthCred{}, fmt.Errorf("%s was logged out while its login was renewed; nothing was saved", name)
		case err != nil:
			continue
		case out == nil:
			return stored, nil // renewed again meanwhile; that login is the one in use
		}
		if err = fsutil.AtomicWrite(path, out, 0o600); err == nil {
			return got, nil
		}
	}
	return oauthCred{}, fmt.Errorf("%w: the renewed login could not be saved (%v); the old one may no longer work: run `ccam login %s`", errLoginLost, err, name)
}

// merged is the login file at path with got's tokens in claudeAiOauth; out
// is nil, with the stored login, when the file holds a refresh token other
// than the one posted.
func merged(path, posted string, got oauthCred) (stored oauthCred, out []byte, err error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return oauthCred{}, nil, errLoggedOut
	}
	if err != nil {
		return oauthCred{}, nil, err
	}
	doc := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return oauthCred{}, nil, err
	}
	stored, ok := storedOAuth(data)
	oauth := map[string]any{}
	if !ok || json.Unmarshal(doc["claudeAiOauth"], &oauth) != nil {
		return oauthCred{}, nil, errLoggedOut
	}
	// Claude Code saves only over the refresh token it posted or "" (what it
	// leaves for a dead one); a login without one, or with another, stays.
	if rt, ok := oauth["refreshToken"].(string); !ok || rt != "" && rt != posted {
		return stored, nil, nil
	}
	oauth["accessToken"], oauth["refreshToken"], oauth["expiresAt"] = got.AccessToken, got.RefreshToken, got.ExpiresAt
	if got.refreshExpiresAt > 0 {
		oauth["refreshTokenExpiresAt"] = got.refreshExpiresAt
	}
	doc["claudeAiOauth"], _ = json.Marshal(oauth)
	out, err = json.Marshal(doc)
	return stored, out, err
}
