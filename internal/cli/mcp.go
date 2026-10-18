package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func mcpCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "공유 MCP 서버 관리 (그룹의 mcp.json)",
	}
	var group string
	cmd.PersistentFlags().StringVarP(&group, "group", "g", "", "대상 그룹 (기본: MCP를 공유하는 그룹이 하나면 그 그룹, 아니면 default)")
	cmd.AddCommand(mcpListCmd(e), mcpAddCmd(e, &group), mcpAddJSONCmd(e, &group), mcpRemoveCmd(e, &group), mcpImportCmd(e))
	return cmd
}

// mcpGroups lists groups that some profile shares MCP through.
func mcpGroups(cfg *config.Config) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range cfg.Profiles {
		if g := p.Share.Get("mcp"); g != config.Local && !seen[g] {
			seen[g] = true
			out = append(out, g)
		}
	}
	sort.Strings(out)
	return out
}

func pickGroup(cfg *config.Config, flag string) string {
	if flag != "" {
		return flag
	}
	if gs := mcpGroups(cfg); len(gs) == 1 {
		return gs[0]
	}
	return "default"
}

func mcpListCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "그룹별 공유 MCP 서버와 계정에만 있는 서버",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			groups := mcpGroups(cfg)
			if len(groups) == 0 {
				e.printf("MCP를 공유하는 계정이 없습니다. `ccam share <이름> --mcp on` 으로 켜세요.\n")
			}
			shared := map[string]map[string]json.RawMessage{}
			for _, g := range groups {
				store, err := mcp.Load(mcp.Path(e.paths.GroupDir(g)))
				if err != nil {
					return err
				}
				shared[g] = store.Servers
				var members []string
				for _, p := range cfg.Profiles {
					if p.Share.Get("mcp") == g {
						members = append(members, p.Name)
					}
				}
				e.printf("%s 그룹 (%s)\n", g, strings.Join(members, ", "))
				if len(store.Servers) == 0 {
					e.printf("  (없음) `ccam mcp add` 로 추가\n")
				}
				for _, n := range store.Names() {
					e.printf("  %-20s %s\n", n, describeServer(store.Servers[n]))
				}
			}
			for i := range cfg.Profiles {
				p := &cfg.Profiles[i]
				local, err := mcp.LocalServers(profile.Dir(e.paths, p))
				if err != nil || len(local) == 0 {
					continue
				}
				var only []string
				for n := range local {
					if _, ok := shared[p.Share.Get("mcp")][n]; !ok {
						only = append(only, n)
					}
				}
				if len(only) > 0 {
					sort.Strings(only)
					e.printf("%s 계정에만 있음: %s\n", p.Name, strings.Join(only, ", "))
				}
			}
			return nil
		},
	}
}

func describeServer(raw json.RawMessage) string {
	var v struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
		URL     string   `json:"url"`
	}
	json.Unmarshal(raw, &v)
	if v.URL != "" {
		t := v.Type
		if t == "" {
			t = "http"
		}
		return t + " " + v.URL
	}
	return "stdio " + strings.TrimSpace(v.Command+" "+strings.Join(v.Args, " "))
}

func mcpAddCmd(e *env, group *string) *cobra.Command {
	var envs, headers []string
	var transport, url string
	cmd := &cobra.Command{
		Use:   "add <이름> [-- 명령 인자…] | add <이름> --url <주소>",
		Short: "그룹에 MCP 서버 추가",
		Example: "  ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code\n" +
			"  ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def := map[string]any{}
			switch {
			case url != "":
				if transport == "" || transport == "stdio" {
					transport = "http"
				}
				def["type"], def["url"] = transport, url
				if len(headers) > 0 {
					h := map[string]string{}
					for _, kv := range headers {
						k, v, ok := strings.Cut(kv, ":")
						if !ok {
							return fmt.Errorf("--header 는 '이름: 값' 형식입니다: %q", kv)
						}
						h[strings.TrimSpace(k)] = strings.TrimSpace(v)
					}
					def["headers"] = h
				}
			case len(args) > 1:
				def["type"], def["command"], def["args"] = "stdio", args[1], args[2:]
				if len(envs) > 0 {
					m := map[string]string{}
					for _, kv := range envs {
						k, v, ok := strings.Cut(kv, "=")
						if !ok {
							return fmt.Errorf("--env 는 KEY=VALUE 형식입니다: %q", kv)
						}
						m[k] = v
					}
					def["env"] = m
				}
			default:
				return fmt.Errorf("실행할 명령(-- 뒤) 이나 --url 이 필요합니다")
			}
			raw, _ := json.Marshal(def)
			return addServer(e, *group, args[0], raw)
		},
	}
	cmd.Flags().StringArrayVarP(&envs, "env", "e", nil, "환경변수 KEY=VALUE (stdio)")
	cmd.Flags().StringVar(&url, "url", "", "원격 서버 주소 (http/sse)")
	cmd.Flags().StringVar(&transport, "transport", "", "stdio|http|sse")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "HTTP 헤더 '이름: 값'")
	return cmd
}

func mcpAddJSONCmd(e *env, group *string) *cobra.Command {
	return &cobra.Command{
		Use:   "add-json <이름> <json>",
		Short: "JSON 정의로 MCP 서버 추가 (claude mcp add-json 과 같은 형식)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addServer(e, *group, args[0], json.RawMessage(args[1]))
		},
	}
}

func addServer(e *env, groupFlag, name string, def json.RawMessage) error {
	if err := mcp.Validate(def); err != nil {
		return err
	}
	cfg, err := e.load()
	if err != nil {
		return err
	}
	g := pickGroup(cfg, groupFlag)
	path := mcp.Path(e.paths.GroupDir(g))
	store, err := mcp.Load(path)
	if err != nil {
		return err
	}
	_, existed := store.Servers[name]
	store.Servers[name] = def
	if err := store.Save(path); err != nil {
		return err
	}
	verb := "추가"
	if existed {
		verb = "갱신"
	}
	e.printf("✓ %s 그룹에 %s %s — 다음 실행부터 적용\n", g, name, verb)
	if len(mcpGroups(cfg)) == 0 {
		e.printf("참고: 아직 MCP를 공유하는 계정이 없습니다. `ccam share <이름> --mcp %s`\n", g)
	}
	return nil
}

func mcpRemoveCmd(e *env, group *string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <이름>",
		Short: "그룹에서 MCP 서버 삭제",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			g := pickGroup(cfg, *group)
			path := mcp.Path(e.paths.GroupDir(g))
			store, err := mcp.Load(path)
			if err != nil {
				return err
			}
			if _, ok := store.Servers[args[0]]; !ok {
				return fmt.Errorf("%s 그룹에 %q 서버가 없습니다", g, args[0])
			}
			delete(store.Servers, args[0])
			if err := store.Save(path); err != nil {
				return err
			}
			e.printf("✓ %s 그룹에서 %s 삭제\n", g, args[0])
			return nil
		},
	}
}

func mcpImportCmd(e *env) *cobra.Command {
	var keepLocal bool
	cmd := &cobra.Command{
		Use:   "import <계정> [서버…]",
		Short: "계정에만 있는 MCP 서버를 그 계정의 공유 그룹으로 옮기기",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			g := prof.Share.Get("mcp")
			if g == config.Local {
				return fmt.Errorf("%s 계정은 MCP 공유가 꺼져 있습니다. `ccam share %s --mcp on` 먼저", prof.Name, prof.Name)
			}
			dir := profile.Dir(e.paths, prof)
			local, err := mcp.LocalServers(dir)
			if err != nil {
				return err
			}
			names := args[1:]
			if len(names) == 0 {
				for n := range local {
					names = append(names, n)
				}
				sort.Strings(names)
			}
			if len(names) == 0 {
				e.printf("%s 계정에만 있는 MCP 서버가 없습니다.\n", prof.Name)
				return nil
			}
			path := mcp.Path(e.paths.GroupDir(g))
			store, err := mcp.Load(path)
			if err != nil {
				return err
			}
			for _, n := range names {
				def, ok := local[n]
				if !ok {
					return fmt.Errorf("%s 계정에 %q 서버가 없습니다", prof.Name, n)
				}
				store.Servers[n] = def
			}
			if err := store.Save(path); err != nil {
				return err
			}
			e.printf("✓ %s 그룹으로 복사: %s\n", g, strings.Join(names, ", "))
			if keepLocal {
				return nil
			}
			// Local entries win over --mcp-config on name clashes, so drop them.
			if launcher.Running(e.paths, prof.Name) {
				e.printf("참고: %s 계정이 실행 중이라 계정 쪽 항목은 그대로 둡니다. 종료 후 다시 실행하면 지웁니다.\n", prof.Name)
				return nil
			}
			if err := mcp.RemoveLocal(dir, filepath.Join(e.paths.State, "mcp", prof.Name), names); err != nil {
				return fmt.Errorf("계정 쪽 항목 삭제 실패: %w", err)
			}
			e.printf("✓ %s 계정의 개별 항목 삭제 (.claude.json 백업함)\n", prof.Name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&keepLocal, "keep-local", false, "계정의 .claude.json 에서는 지우지 않음")
	return cmd
}
