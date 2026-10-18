package usage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/lock"
)

// TokenURL and clientID are what Claude Code uses to refresh its OAuth
// token. Not verified against the live service by ccam's tests (doing so
// would rotate a real refresh token); see docs/claude-layout.md.
const (
	TokenURL = "https://console.anthropic.com/v1/oauth/token"
	clientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
)

// refresh renews an expired access token. Refresh tokens rotate, so doing
// this while a Claude Code session of the profile runs would log that
// session out: the profile lock must be free, and is held meanwhile.
func (f *Fetcher) refresh(ctx context.Context, a Account) (oauthCred, error) {
	if a.Imported {
		return oauthCred{}, fmt.Errorf("제자리 등록한 계정은 그냥 claude 로도 실행될 수 있어 갱신하지 않습니다")
	}
	l, err := lock.TryFlock(a.LockPath)
	if err == lock.ErrBusy {
		return oauthCred{}, fmt.Errorf("실행 중인 세션이 있어 갱신하지 않았습니다")
	}
	if err != nil {
		return oauthCred{}, err
	}
	defer l.Unlock()

	path := filepath.Join(a.Dir, ".credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return oauthCred{}, err
	}
	doc := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return oauthCred{}, err
	}
	oauth := map[string]any{}
	if err := json.Unmarshal(doc["claudeAiOauth"], &oauth); err != nil {
		return oauthCred{}, err
	}
	rt, _ := oauth["refreshToken"].(string)
	if rt == "" {
		return oauthCred{}, fmt.Errorf("refresh token 이 없습니다")
	}
	body, _ := json.Marshal(map[string]string{"grant_type": "refresh_token", "refresh_token": rt, "client_id": clientID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.TokenURL, bytes.NewReader(body))
	if err != nil {
		return oauthCred{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.Client.Do(req)
	if err != nil {
		return oauthCred{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return oauthCred{}, fmt.Errorf("토큰 갱신 실패 (HTTP %d) → `ccam login %s`", resp.StatusCode, a.Name)
	}
	var tok struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&tok); err != nil || tok.AccessToken == "" {
		return oauthCred{}, fmt.Errorf("토큰 갱신 응답을 읽지 못했습니다")
	}
	expires := f.Now().UnixMilli() + tok.ExpiresIn*1000
	oauth["accessToken"] = tok.AccessToken
	oauth["expiresAt"] = expires
	if tok.RefreshToken != "" {
		oauth["refreshToken"] = tok.RefreshToken
	}
	raw, _ := json.Marshal(oauth)
	doc["claudeAiOauth"] = raw
	out, err := json.Marshal(doc)
	if err != nil {
		return oauthCred{}, err
	}
	if err := fsutil.AtomicWrite(path, out, 0o600); err != nil {
		return oauthCred{}, err
	}
	return oauthCred{AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken, ExpiresAt: expires}, nil
}
