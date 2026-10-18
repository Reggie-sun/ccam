// Package doctor checks a ccam installation and repairs what it safely can.
package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/usage"
)

type Level int

const (
	OK Level = iota
	Info
	Warn
	Error
)

type Finding struct {
	Level   Level
	Profile string
	Msg     string
	Fixed   bool
}

func (f Finding) String() string {
	mark := [...]string{"✓", "·", "!", "✗"}[f.Level]
	prefix := ""
	if f.Profile != "" {
		prefix = f.Profile + ": "
	}
	suffix := ""
	if f.Fixed {
		suffix = " → 고침"
	}
	return fmt.Sprintf("%s %s%s%s", mark, prefix, f.Msg, suffix)
}

type Options struct {
	Fix      bool
	CcamPath string         // the real ccam binary, for shim checks
	Usage    *usage.Fetcher // nil skips the live usage-format check
}

// LayoutVersion is the Claude Code version docs/claude-layout.md was checked with.
const LayoutVersion = "2.1.302"

type report struct {
	out []Finding
}

func (r *report) add(level Level, prof, format string, a ...any) *Finding {
	r.out = append(r.out, Finding{Level: level, Profile: prof, Msg: fmt.Sprintf(format, a...)})
	return &r.out[len(r.out)-1]
}

// Run checks everything and, with opt.Fix, repairs safe problems.
func Run(p paths.Paths, cfg *config.Config, opt Options) []Finding {
	r := &report{}
	checkClaude(r, p, cfg)
	for i := range cfg.Profiles {
		checkProfile(r, p, &cfg.Profiles[i], opt.Fix)
	}
	checkPlugins(r, p, cfg, opt.Fix)
	checkMCP(r, p, cfg)
	checkUsageFormat(r, p, cfg, opt.Usage)
	checkShims(r, p, cfg, opt)
	return r.out
}

// Failed reports whether any finding is still an unfixed error.
func Failed(fs []Finding) bool {
	for _, f := range fs {
		if f.Level == Error && !f.Fixed {
			return true
		}
	}
	return false
}

func checkClaude(r *report, p paths.Paths, cfg *config.Config) {
	bin, err := launcher.ResolveClaude(p, cfg.ClaudeBin)
	if err != nil {
		r.add(Error, "", "%v", err)
		return
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		r.add(Error, "", "%s --version 실패: %v", bin, err)
		return
	}
	version := strings.TrimSpace(string(out))
	r.add(OK, "", "Claude Code %s (%s)", version, bin)
	if !strings.HasPrefix(version, LayoutVersion) {
		r.add(Info, "", "ccam 은 Claude Code %s 의 폴더 구성으로 확인됐습니다. 아래 '알려지지 않은 항목'이 있으면 구성이 바뀐 것일 수 있습니다", LayoutVersion)
	}
}

// checkUsageFormat queries one logged-in account so a changed response
// format of the undocumented usage endpoint is noticed.
func checkUsageFormat(r *report, p paths.Paths, cfg *config.Config, f *usage.Fetcher) {
	if f == nil {
		return
	}
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		a := usage.Account{Name: prof.Name, Slot: prof.Slot, Dir: profile.Dir(p, prof), SessionDir: launcher.SessionDir(p, prof.Name)}
		rep := f.Fetch(context.Background(), a, usage.Options{NoCache: true})
		switch {
		case rep.FormatChanged:
			r.add(Warn, prof.Name, "사용량 응답 형식이 바뀌었습니다. `ccam usage` 는 상태줄 캐시로 대신 보여줍니다 (%s)", rep.Err)
			return
		case rep.Status == usage.StatusOK:
			r.add(OK, prof.Name, "사용량 조회 정상")
			return
		}
	}
}

func checkProfile(r *report, p paths.Paths, prof *config.Profile, fix bool) {
	dir := profile.Dir(p, prof)
	info, err := os.Stat(dir)
	if err != nil {
		r.add(Error, prof.Name, "프로필 폴더가 없습니다: %s", dir)
		return
	}
	if info.Mode().Perm()&0o077 != 0 {
		f := r.add(Warn, prof.Name, "폴더 권한이 %o 입니다 (0700 이어야 함)", info.Mode().Perm())
		if fix && os.Chmod(dir, 0o700) == nil {
			f.Fixed = true
		}
	}
	checkCredentials(r, prof, dir, fix)

	problems := share.Check(dir, prof.Share, p.GroupsRoot())
	if len(problems) == 0 {
		r.add(OK, prof.Name, "공유 링크 정상")
	}
	for _, pr := range problems {
		level := Error
		if pr.Kind == share.MissingLink {
			level = Warn
		}
		r.add(level, prof.Name, "%s", pr)
	}
	if fix && len(problems) > 0 {
		// Build with the target as both sides: it re-derives every step from
		// what is actually on disk, keeps group content on conflicts and
		// backs up anything it replaces.
		pl, err := share.Build(dir, p.GroupsRoot(), prof.Share, prof.Share, share.Options{})
		if err == nil {
			err = pl.Apply()
		}
		if err != nil {
			r.add(Error, prof.Name, "공유 링크 수리 실패: %v", err)
			return
		}
		for i := len(r.out) - len(problems); i < len(r.out); i++ {
			r.out[i].Fixed = true
		}
	}
	checkUnknownEntries(r, prof, dir)
}

func checkCredentials(r *report, prof *config.Profile, dir string, fix bool) {
	cred := filepath.Join(dir, ".credentials.json")
	info, err := os.Lstat(cred)
	if err != nil {
		r.add(Warn, prof.Name, "로그인 정보가 없습니다 → `ccam login %s`", prof.Name)
		return
	}
	if info.Mode()&os.ModeSymlink != 0 {
		r.add(Error, prof.Name, ".credentials.json 이 심볼릭 링크입니다. 토큰은 계정마다 따로 있어야 합니다")
		return
	}
	if info.Mode().Perm() != 0o600 {
		f := r.add(Warn, prof.Name, ".credentials.json 권한이 %o 입니다 (0600 이어야 함)", info.Mode().Perm())
		if fix && os.Chmod(cred, 0o600) == nil {
			f.Fixed = true
		}
	}
	data, err := os.ReadFile(cred)
	if err != nil {
		r.add(Error, prof.Name, ".credentials.json 읽기 실패: %v", err)
		return
	}
	var c struct {
		OAuth *struct {
			ExpiresAt             int64 `json:"expiresAt"`
			RefreshTokenExpiresAt int64 `json:"refreshTokenExpiresAt"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(data, &c); err != nil || c.OAuth == nil {
		r.add(Info, prof.Name, "구독 로그인 토큰이 없습니다 (API 키 계정일 수 있음)")
		return
	}
	now := time.Now().UnixMilli()
	switch {
	case c.OAuth.RefreshTokenExpiresAt > 0 && c.OAuth.RefreshTokenExpiresAt < now:
		r.add(Error, prof.Name, "로그인이 만료됐습니다 → `ccam login %s`", prof.Name)
	case c.OAuth.ExpiresAt < now:
		r.add(Info, prof.Name, "접속 토큰이 만료됐습니다. 다음에 실행하면 Claude Code가 갱신합니다")
	default:
		r.add(OK, prof.Name, "로그인 유효 (%s)", orDash(prof.Email))
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// checkUnknownEntries flags files Claude Code started creating that ccam
// does not know about, a hint that the layout changed.
func checkUnknownEntries(r *report, prof *config.Profile, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var unknown []string
	for _, e := range entries {
		if !share.Known(e.Name()) {
			unknown = append(unknown, e.Name())
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		r.add(Info, prof.Name, "알려지지 않은 항목 (Claude Code 구성이 바뀌었을 수 있음): %s", strings.Join(unknown, ", "))
	}
}

// checkPlugins finds plugin paths that point through one profile's folder
// into a shared group. They break when that profile is removed or renamed,
// so --fix rewrites them to the group folder itself.
func checkPlugins(r *report, p paths.Paths, cfg *config.Config, fix bool) {
	groups, _ := os.ReadDir(p.GroupsRoot())
	sep := string(filepath.Separator)
	for _, g := range groups {
		pluginsDir := filepath.Join(p.GroupDir(g.Name()), "plugins")
		n := 0
		for i := range cfg.Profiles {
			pre := filepath.Join(profile.Dir(p, &cfg.Profiles[i]), "plugins") + sep
			n += countPrefix(pluginsDir, pre)
		}
		if n == 0 {
			continue
		}
		f := r.add(Warn, "", "%s 그룹의 플러그인 경로 %d개가 특정 계정 폴더를 거칩니다", g.Name(), n)
		if !fix {
			continue
		}
		f.Fixed = true
		for i := range cfg.Profiles {
			pre := filepath.Join(profile.Dir(p, &cfg.Profiles[i]), "plugins") + sep
			if _, err := profile.RewritePluginPaths(pluginsDir, pre, pluginsDir+sep); err != nil {
				f.Fixed = false
			}
		}
	}
}

func countPrefix(pluginsDir, prefix string) int {
	n := 0
	for _, name := range []string{"installed_plugins.json", "known_marketplaces.json"} {
		data, err := os.ReadFile(filepath.Join(pluginsDir, name))
		if err != nil {
			continue
		}
		var doc any
		if json.Unmarshal(data, &doc) != nil {
			continue
		}
		profile.RewriteStrings(doc, func(s string) string {
			if strings.HasPrefix(s, prefix) {
				n++
			}
			return s
		})
	}
	return n
}

func checkMCP(r *report, p paths.Paths, cfg *config.Config) {
	seen := map[string]bool{}
	for _, prof := range cfg.Profiles {
		g := prof.Share.Get("mcp")
		if g == config.Local || seen[g] {
			continue
		}
		seen[g] = true
		store, err := mcp.Load(mcp.Path(p.GroupDir(g)))
		if err != nil {
			r.add(Error, "", "%v", err)
			continue
		}
		for _, n := range store.Names() {
			if err := mcp.Validate(store.Servers[n]); err != nil {
				r.add(Warn, "", "%s 그룹 MCP 서버 %s: %v", g, n, err)
			}
		}
	}
}

func checkShims(r *report, p paths.Paths, cfg *config.Config, opt Options) {
	if opt.CcamPath == "" || len(cfg.Profiles) == 0 {
		return
	}
	var missing []string
	for _, prof := range cfg.Profiles {
		name := "claude" + strconv.Itoa(prof.Slot)
		if _, err := os.Lstat(filepath.Join(p.Bin, name)); err != nil {
			missing = append(missing, name)
		}
	}
	if _, err := os.Lstat(filepath.Join(p.Bin, "ccs")); err != nil {
		missing = append(missing, "ccs")
	}
	if len(missing) == 0 {
		r.add(OK, "", "단축 명령 정상 (%s)", p.Bin)
	} else {
		f := r.add(Warn, "", "단축 명령 없음: %s", strings.Join(missing, ", "))
		if opt.Fix {
			_, err := shim.Install(p.Bin, opt.CcamPath, cfg)
			f.Fixed = err == nil
		}
	}
	if !shim.InPath(p.Bin) {
		r.add(Warn, "", "%s 가 PATH 에 없습니다", p.Bin)
	}
}
