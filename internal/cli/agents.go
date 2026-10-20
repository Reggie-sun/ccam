package cli

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// agentEntry is one session `claude agents --json` lists (Claude Code
// 2.1.304); background sessions have an id and a state.
type agentEntry struct {
	Account    string `json:"account"`
	ID         string `json:"id,omitempty"`
	Kind       string `json:"kind"`
	Name       string `json:"name,omitempty"`
	State      string `json:"state,omitempty"`
	Status     string `json:"status,omitempty"`
	WaitingFor string `json:"waitingFor,omitempty"`
	Cwd        string `json:"cwd"`
	StartedAt  int64  `json:"startedAt"`
	SessionID  string `json:"sessionId,omitempty"`
	PID        int    `json:"pid,omitempty"`
}

// agentsTimeout bounds one account's `claude agents --json`.
const agentsTimeout = 15 * time.Second

func agentsCmd(e *env) *cobra.Command {
	var all, asJSON bool
	cmd := &cobra.Command{
		Use:   "agents [--all] [--json]",
		Short: "List the Claude Code sessions of every account: background agents and open sessions",
		Long: "Runs `claude agents --json` in every account at once and lists the sessions\n" +
			"together, with the account each belongs to. Background sessions belong to the\n" +
			"account that started them; open one with `ccam attach <id>`.",
		Example: "  ccam agents\n  ccam agents --all     # also finished background sessions\n  ccam attach 49ec0028",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			list, failed := listAgents(e, cfg, all)
			if asJSON {
				if list == nil {
					list = []agentEntry{}
				}
				data, _ := json.MarshalIndent(list, "", "  ")
				e.printf("%s\n", data)
			} else if len(list) > 0 || len(failed) == 0 {
				printAgents(e, list)
			}
			for _, f := range failed {
				fmt.Fprintln(e.stderr, "ccam:", f)
			}
			if len(failed) > 0 {
				// A script must not take a failure for an empty list.
				return fmt.Errorf("%d of %d accounts could not be asked", len(failed), max(len(cfg.Profiles), 1))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "also list finished background sessions")
	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	return cmd
}

// listAgents asks every account in parallel; an account that fails is
// named in failed and left out.
func listAgents(e *env, cfg *config.Config, all bool) (list []agentEntry, failed []string) {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return nil, []string{err.Error()}
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := accountAgents(e, bin, p, all)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s: could not list its sessions: %v", p.Name, err))
				return
			}
			list = append(list, got...)
		}()
	}
	wg.Wait()
	slices.SortFunc(list, func(a, b agentEntry) int { return cmp.Compare(b.StartedAt, a.StartedAt) })
	slices.Sort(failed)
	return list, failed
}

func accountAgents(e *env, bin string, p *config.Profile, all bool) ([]agentEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), agentsTimeout)
	defer cancel()
	args := []string{"agents", "--json"}
	if all {
		args = append(args, "--all")
	}
	c := exec.CommandContext(ctx, bin, args...)
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, p), p.Name)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, strings.TrimSpace(stderr.String()))
	}
	var list []agentEntry
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("unexpected output of claude agents --json: %w", err)
	}
	for i := range list {
		list[i].Account = p.Name
	}
	return list, nil
}

func printAgents(e *env, list []agentEntry) {
	if len(list) == 0 {
		e.printf("No sessions. Start one in the background with `ccam <account> --bg \"…\"`.\n")
		return
	}
	home, _ := os.UserHomeDir()
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Account\tID\tName\tState\tFolder\tStarted")
	for _, a := range list {
		id, state := cmp.Or(a.ID, "-"), cmp.Or(a.State, a.Status, "-")
		if a.Kind == "interactive" {
			state = "open (" + cmp.Or(a.Status, "interactive") + ")"
		} else if a.Status != "" && a.Status != a.State {
			state += ", " + a.Status
		}
		if a.WaitingFor != "" {
			state += ": " + a.WaitingFor
		}
		dir := a.Cwd
		if rel, err := filepath.Rel(home, dir); err == nil && home != "" && !strings.HasPrefix(rel, "..") {
			dir = filepath.Join("~", rel)
		}
		started := "-"
		if a.StartedAt > 0 {
			started = Ago(time.Since(time.UnixMilli(a.StartedAt)))
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", a.Account, id, cmp.Or(a.Name, "-"), state, dir, started)
	}
	_ = tw.Flush()
	e.printf("Open a background session with `ccam attach <id>`.\n")
}

func attachCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "attach <id> [claude args]",
		Short: "Open a background session of whichever account started it",
		Long: "Finds the account whose folder has the background session (jobs/<id>) and runs\n" +
			"`claude attach <id>` as that account. A unique start of the id is enough.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			var hits []string // account/id
			for i := range cfg.Profiles {
				p := &cfg.Profiles[i]
				entries, _ := os.ReadDir(filepath.Join(profile.Dir(e.paths, p), "jobs"))
				for _, en := range entries {
					if en.IsDir() && strings.HasPrefix(en.Name(), args[0]) {
						hits = append(hits, p.Name+"/"+en.Name())
					}
				}
			}
			switch len(hits) {
			case 0:
				return fmt.Errorf("no account has a background session %q (see `ccam agents --all`)", args[0])
			case 1:
			default:
				return fmt.Errorf("%q fits more than one session: %s; give more of the id", args[0], strings.Join(hits, ", "))
			}
			account, id, _ := strings.Cut(hits[0], "/")
			fmt.Fprintf(e.stderr, "→ %s (%s)\n", id, account)
			return launcher.LaunchWith(e.paths, cfg, account, append([]string{"attach", id}, args[1:]...), e.dryRun, e.stdout)
		},
	}
}
