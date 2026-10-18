package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
)

func addCmd(e *env) *cobra.Command {
	var noLogin bool
	var shares []string
	cmd := &cobra.Command{
		Use:   "add [이름]",
		Short: "계정 추가: 프로필 생성 → 로그인 → 단축 명령",
		Long:  "터미널에서 옵션 없이 실행하면 마법사가 이름과 공유 항목을 묻습니다.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive(e) && !cmd.Flags().Changed("share") && !noLogin {
				cfg, err := e.load()
				if err != nil {
					return err
				}
				name := ""
				if len(args) == 1 {
					if err := tui.ValidateNewName(cfg)(args[0]); err != nil {
						return err
					}
					name = args[0]
				}
				return addInteractive(e, cfg, name)
			}
			if len(args) == 0 {
				return fmt.Errorf("계정 이름이 필요합니다: ccam add <이름>")
			}
			share := config.DefaultShare()
			for _, s := range shares {
				item, group, ok := cutEq(s)
				if !ok {
					return fmt.Errorf("--share 는 항목=그룹 형식입니다 (예: history=team): %q", s)
				}
				if err := share.Set(item, group); err != nil {
					return err
				}
			}
			return addProfile(e, args[0], share, !noLogin)
		},
	}
	cmd.Flags().BoolVar(&noLogin, "no-login", false, "로그인은 나중에 (`ccam login <이름>`)")
	cmd.Flags().StringArrayVar(&shares, "share", nil, "공유 설정 항목=그룹|local (반복 가능)")
	return cmd
}

func addProfile(e *env, name string, share config.Share, login bool) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	prof, err := profile.Add(e.paths, cfg, name, share)
	if err != nil {
		return err
	}
	if err := applyShare(e, cfg, prof); err != nil {
		return err
	}
	if err := e.save(cfg); err != nil {
		return err
	}
	if login {
		if err := loginProfile(e, cfg, name); err != nil {
			fmt.Fprintf(e.stderr, "로그인하지 못했습니다: %v\n나중에 `ccam login %s` 로 다시 시도하세요.\n", err, name)
		}
	}
	if _, err := installShims(e, cfg); err != nil {
		fmt.Fprintln(e.stderr, "단축 명령을 만들지 못했습니다:", err)
	}
	prof, _ = cfg.Find(name)
	email := prof.Email
	if email == "" {
		email = "로그인 전"
	}
	e.printf("✓ %s (%s) 추가됨 — 단축 명령: claude%d\n", name, email, prof.Slot)
	return nil
}

// loginProfile runs `claude auth login` inside the profile and caches the email.
func loginProfile(e *env, cfg *config.Config, sel string) error {
	prof, err := cfg.MustFind(sel)
	if err != nil {
		return err
	}
	if err := runClaude(e, cfg, prof, "auth", "login"); err != nil {
		return err
	}
	email, err := profile.ReadEmail(profile.Dir(e.paths, prof))
	if err != nil {
		return fmt.Errorf("로그인 후 이메일을 읽지 못했습니다: %w", err)
	}
	prof.Email = email
	return e.save(cfg)
}

// runClaude runs a claude subcommand (not exec) inside a profile.
func runClaude(e *env, cfg *config.Config, prof *config.Profile, args ...string) error {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return err
	}
	c := exec.Command(bin, args...)
	c.Env = append(os.Environ(), "CLAUDE_CONFIG_DIR="+profile.Dir(e.paths, prof), "CCAM_PROFILE="+prof.Name)
	c.Stdin, c.Stdout, c.Stderr = e.stdin, e.stdout, e.stderr
	return c.Run()
}

func installShims(e *env, cfg *config.Config) ([]string, error) {
	self, err := ccamPath()
	if err != nil {
		return nil, err
	}
	created, err := shim.Install(e.paths.Bin, self, cfg)
	if err == nil && !shim.InPath(e.paths.Bin) {
		fmt.Fprintf(e.stderr, "참고: %s 가 PATH 에 없습니다. ~/.bashrc 나 ~/.zshrc 에 추가하세요:\n  export PATH=\"%s:$PATH\"\n", e.paths.Bin, e.paths.Bin)
	}
	return created, err
}

func cutEq(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}

type listRow struct {
	Slot     int          `json:"slot"`
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Default  bool         `json:"default"`
	Dir      string       `json:"dir"`
	LastUsed *time.Time   `json:"last_used,omitempty"`
	Share    config.Share `json:"share"`
}

func listCmd(e *env) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "계정 목록",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProfiles(e, asJSON)
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "JSON 으로 출력")
	return cmd
}

func listProfiles(e *env, asJSON bool) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	def, _ := cfg.DefaultProfile()
	rows := make([]listRow, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		r := listRow{Slot: p.Slot, Name: p.Name, Email: p.Email, Default: def == p, Dir: profile.Dir(e.paths, p), Share: p.Share}
		if t := launcher.LastUsed(e.paths, p.Name); !t.IsZero() {
			r.LastUsed = &t
		}
		rows = append(rows, r)
	}
	if asJSON {
		enc := json.NewEncoder(e.stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}
	if len(rows) == 0 {
		e.printf("등록된 계정이 없습니다. `ccam add <이름>` 으로 추가하세요.\n")
		return nil
	}
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, " \t번호\t이름\t이메일\t공유\t마지막 사용")
	for _, r := range rows {
		mark := ""
		if r.Default {
			mark = "*"
		}
		email := r.Email
		if email == "" {
			email = "(로그인 전)"
		}
		last := "-"
		if r.LastUsed != nil {
			last = Ago(time.Since(*r.LastUsed))
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", mark, "claude"+strconv.Itoa(r.Slot), r.Name, email, shareSummary(r.Share), last)
	}
	return tw.Flush()
}

// shareSummary renders e.g. "mcp settings memory ext" for shared items.
func shareSummary(s config.Share) string {
	short := map[string]string{"mcp": "mcp", "settings": "settings", "memory": "md", "extensions": "ext", "history": "history"}
	out := ""
	for _, item := range config.Items {
		g := s.Get(item)
		if g == config.Local {
			continue
		}
		if out != "" {
			out += " "
		}
		out += short[item]
		if g != "default" {
			out += "@" + g
		}
	}
	if out == "" {
		return "-"
	}
	return out
}

// Ago renders a short Korean relative time.
func Ago(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "방금"
	case d < time.Hour:
		return fmt.Sprintf("%d분 전", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d시간 전", int(d.Hours()))
	default:
		return fmt.Sprintf("%d일 전", int(d.Hours()/24))
	}
}

func useCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "use <이름|번호>",
		Short: "기본 계정 지정",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			cfg.Default = p.Name
			if err := e.save(cfg); err != nil {
				return err
			}
			e.printf("✓ 기본 계정: %s\n", p.Name)
			return nil
		},
	}
}

func removeCmd(e *env) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:     "remove <이름|번호>",
		Aliases: []string{"rm"},
		Short:   "계정 삭제 (폴더는 휴지통으로 이동)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			prof := *p
			if !yes && !e.confirm(fmt.Sprintf("%s (%s) 계정을 삭제할까요?", prof.Name, profile.Dir(e.paths, &prof))) {
				e.printf("취소했습니다.\n")
				return nil
			}
			trashed, err := profile.Trash(e.paths, &prof)
			if err != nil {
				return err
			}
			cfg.Remove(prof.Name)
			if err := e.save(cfg); err != nil {
				return err
			}
			if self, err := ccamPath(); err == nil {
				shim.RemoveFor(e.paths.Bin, self, &prof)
			}
			if trashed != "" {
				e.printf("✓ %s 삭제됨 (폴더는 %s 로 옮김)\n", prof.Name, trashed)
			} else {
				e.printf("✓ %s 등록 해제됨 (폴더 %s 는 그대로 둠)\n", prof.Name, profile.Dir(e.paths, &prof))
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "확인 없이 삭제")
	return cmd
}

func renameCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "rename <이름> <새 이름>",
		Short: "계정 이름 변경 (단축 명령 갱신)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			old, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			oldName := old.Name
			if err := profile.Rename(e.paths, cfg, oldName, args[1]); err != nil {
				return err
			}
			if err := e.save(cfg); err != nil {
				return err
			}
			if _, err := installShims(e, cfg); err != nil {
				fmt.Fprintln(e.stderr, "단축 명령 갱신 실패:", err)
			}
			e.printf("✓ %s → %s\n", oldName, args[1])
			return nil
		},
	}
}

func loginCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "login <이름|번호>",
		Short: "다시 로그인",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if err := loginProfile(e, cfg, args[0]); err != nil {
				return err
			}
			p, _ := cfg.Find(args[0])
			e.printf("✓ %s 로그인: %s\n", p.Name, p.Email)
			return nil
		},
	}
}

func logoutCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "logout <이름|번호>",
		Short: "토큰 삭제",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			if err := runClaude(e, cfg, p, "auth", "logout"); err != nil {
				return err
			}
			p.Email = ""
			return e.save(cfg)
		},
	}
}

func shimsCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{Use: "shims", Short: "단축 명령(claude1…, ccs) 관리"}
	var names bool
	install := &cobra.Command{
		Use:   "install",
		Short: "claude1…claudeN 과 ccs 를 다시 만듭니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("names") && cfg.NameShims != names {
				cfg.NameShims = names
				if err := e.save(cfg); err != nil {
					return err
				}
			}
			created, err := installShims(e, cfg)
			if err != nil {
				return err
			}
			e.printf("✓ 단축 명령 %d개 새로 만듦 (%s)\n", len(created), e.paths.Bin)
			return nil
		},
	}
	install.Flags().BoolVar(&names, "names", false, "claude-<이름> 단축 명령도 만들기")
	cmd.AddCommand(install)
	return cmd
}

func runCmd(e *env) *cobra.Command {
	// Only reached for `ccam help run`; real runs never get to cobra.
	return &cobra.Command{
		Use:                "run <이름|번호> [-- claude 인자]",
		Short:              "해당 계정으로 claude 실행 (`ccs <이름>` 과 같음)",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			sel := ""
			if len(args) > 0 && args[0] != "--" && args[0][0] != '-' {
				sel, args = args[0], args[1:]
			}
			return launcher.Launch(e.paths, sel, args, e.dryRun, e.stdout)
		},
	}
}
