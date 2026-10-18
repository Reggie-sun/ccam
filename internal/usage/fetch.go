package usage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// Endpoint is the undocumented usage API (docs/claude-layout.md).
const Endpoint = "https://api.anthropic.com/api/oauth/usage"

const (
	CacheTTL       = 60 * time.Second
	connectTimeout = 2 * time.Second
	accountTimeout = 3 * time.Second
)

type Status string

const (
	StatusOK        Status = "ok"
	StatusExpired   Status = "expired"    // token expired; values may come from cache
	StatusLoggedOut Status = "logged_out" // no credentials
	StatusAPI       Status = "api"        // API key account, no subscription limits
	StatusError     Status = "error"      // request failed; values may come from cache
	StatusPending   Status = "pending"    // not answered yet (live view only)
)

// Account is one profile to query.
type Account struct {
	Name     string
	Slot     int
	Dir      string // CLAUDE_CONFIG_DIR
	LockPath string // held while a session started by ccam runs
	Imported bool   // registered in place; plain `claude` may use it without ccam's lock
}

// Report is one account's result.
type Report struct {
	Name   string    `json:"name"`
	Slot   int       `json:"slot"`
	Status Status    `json:"status"`
	Limits *Limits   `json:"limits,omitempty"`
	Source string    `json:"source,omitempty"` // live | cache | statusline
	At     time.Time `json:"at,omitempty"`     // when the values were measured
	Err    string    `json:"error,omitempty"`
	// FormatChanged means the endpoint answered but the parser did not
	// understand it; doctor reports this.
	FormatChanged bool `json:"format_changed,omitempty"`
}

type Options struct {
	NoCache bool
	Source  string // "" or "live": live first; "cache": cached values only
	Refresh bool   // refresh expired tokens (only when the profile is not running)
}

// Fetcher queries accounts in parallel over one shared HTTP client.
type Fetcher struct {
	Client   *http.Client
	URL      string
	TokenURL string
	StateDir string // ~/.local/state/ccam/usage
	Now      func() time.Time
}

func NewFetcher(stateDir string) *Fetcher {
	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         (&net.Dialer{Timeout: connectTimeout}).DialContext,
		TLSHandshakeTimeout: connectTimeout,
		ForceAttemptHTTP2:   true,
		MaxIdleConnsPerHost: 8,
	}
	return &Fetcher{
		Client:   &http.Client{Transport: tr},
		URL:      Endpoint,
		TokenURL: TokenURL,
		StateDir: stateDir,
		Now:      time.Now,
	}
}

// FetchAll queries every account at once and calls onResult as each one
// finishes (from multiple goroutines, serialized). It returns all reports
// in input order.
func (f *Fetcher) FetchAll(ctx context.Context, accts []Account, opt Options, onResult func(int, Report)) []Report {
	out := make([]Report, len(accts))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, a := range accts {
		wg.Add(1)
		go func(i int, a Account) {
			defer wg.Done()
			r := f.Fetch(ctx, a, opt)
			mu.Lock()
			out[i] = r
			if onResult != nil {
				onResult(i, r)
			}
			mu.Unlock()
		}(i, a)
	}
	wg.Wait()
	return out
}

// Fetch queries one account, falling back to cached values on failure.
func (f *Fetcher) Fetch(ctx context.Context, a Account, opt Options) Report {
	r := Report{Name: a.Name, Slot: a.Slot}
	if opt.Source == "cache" {
		return f.fallback(a, r, StatusOK, "")
	}
	if !opt.NoCache {
		if c, ok := f.readCache(a.Name, "live"); ok && f.Now().Sub(c.At) < CacheTTL {
			r.Status, r.Limits, r.Source, r.At = StatusOK, c.Limits, "live", c.At
			return r
		}
	}
	cred, err := readCredentials(a.Dir)
	switch {
	case errors.Is(err, errNoOAuth):
		if hasAPIKey(a.Dir) {
			r.Status = StatusAPI
			return r
		}
		r.Status = StatusLoggedOut
		return r
	case err != nil:
		return f.fallback(a, r, StatusError, err.Error())
	}
	if cred.expired(f.Now()) {
		if !opt.Refresh {
			return f.fallback(a, r, StatusExpired, "")
		}
		if cred, err = f.refresh(ctx, a); err != nil {
			return f.fallback(a, r, StatusExpired, err.Error())
		}
	}
	ctx, cancel := context.WithTimeout(ctx, accountTimeout)
	defer cancel()
	limits, status, err := f.live(ctx, cred.AccessToken)
	if err != nil {
		r.FormatChanged = errors.Is(err, errFormat)
		return f.fallback(a, r, status, err.Error())
	}
	r.Status, r.Limits, r.Source, r.At = StatusOK, &limits, "live", f.Now()
	f.writeCache(a.Name, "live", cached{At: r.At, Limits: &limits})
	return r
}

func (f *Fetcher) live(ctx context.Context, token string) (Limits, Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.URL, nil)
	if err != nil {
		return Limits{}, StatusError, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("User-Agent", "ccam")
	resp, err := f.Client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return Limits{}, StatusError, fmt.Errorf("응답 없음 (%s 초과)", accountTimeout)
		}
		return Limits{}, StatusError, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Limits{}, StatusError, err
	}
	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return Limits{}, StatusExpired, fmt.Errorf("인증 실패 (%d)", resp.StatusCode)
	case resp.StatusCode != http.StatusOK:
		return Limits{}, StatusError, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	l, err := ParseLive(body)
	if err != nil {
		return Limits{}, StatusError, fmt.Errorf("%w: %v", errFormat, err)
	}
	return l, StatusOK, nil
}

// fallback fills r with the newest cached values (live cache or status line).
func (f *Fetcher) fallback(a Account, r Report, status Status, errMsg string) Report {
	r.Status, r.Err = status, errMsg
	var best cached
	found := false
	for _, src := range []string{"live", "statusline"} {
		if c, ok := f.readCache(a.Name, src); ok && (!found || c.At.After(best.At)) {
			best, found = c, true
			best.source = src
		}
	}
	if !found {
		if status == StatusOK {
			r.Status, r.Err = StatusError, "저장된 사용량이 없습니다"
		}
		return r
	}
	r.Limits, r.At = best.Limits, best.At
	r.Source = "cache"
	if best.source == "statusline" {
		r.Source = "statusline"
	}
	return r
}

type cached struct {
	At     time.Time `json:"at"`
	Limits *Limits   `json:"limits"`
	source string
}

func (f *Fetcher) cachePath(name, src string) string {
	if src == "live" {
		return filepath.Join(f.StateDir, name+".json")
	}
	return filepath.Join(f.StateDir, name+"."+src+".json")
}

func (f *Fetcher) readCache(name, src string) (cached, bool) {
	data, err := os.ReadFile(f.cachePath(name, src))
	if err != nil {
		return cached{}, false
	}
	var c cached
	if json.Unmarshal(data, &c) != nil || c.Limits == nil {
		return cached{}, false
	}
	return c, true
}

func (f *Fetcher) writeCache(name, src string, c cached) {
	data, _ := json.Marshal(c)
	_ = fsutil.WriteFast(f.cachePath(name, src), data)
}

// SaveStatusLine stores values captured by the status line collector.
func SaveStatusLine(stateDir, name string, l Limits, at time.Time) error {
	data, err := json.Marshal(cached{At: at, Limits: &l})
	if err != nil {
		return err
	}
	return fsutil.WriteFast(filepath.Join(stateDir, name+".statusline.json"), data)
}

var (
	errNoOAuth = errors.New("no subscription login")
	errFormat  = errors.New("사용량 응답 형식이 바뀌었습니다")
)

type oauthCred struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"` // Unix ms
}

func (c oauthCred) expired(now time.Time) bool {
	return c.ExpiresAt > 0 && now.UnixMilli() >= c.ExpiresAt
}

// readCredentials reads only; ccam never rewrites the token file except
// for an explicit --refresh.
func readCredentials(dir string) (oauthCred, error) {
	var c oauthCred
	ok, err := profile.DecodeFileKey(filepath.Join(dir, ".credentials.json"), "claudeAiOauth", &c)
	if errors.Is(err, fs.ErrNotExist) || (err == nil && (!ok || c.AccessToken == "")) {
		return c, errNoOAuth
	}
	return c, err
}

func hasAPIKey(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		return false
	}
	_, ok, _ := profile.TopLevelValue(data, "primaryApiKey")
	return ok
}
