// Package group carries a change of groups over to the account folders:
// one share plan for each account whose shared items change, built so that
// a later plan sees what earlier ones will put into the group folders.
package group

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

// Options tune Build.
type Options struct {
	Seed  string // the account whose copy goes into a group first (--from)
	Fresh bool   // accounts that stop sharing an item start without it
}

// Move is one item that changes for an account; From and To are group
// names, "" for not shared.
type Move struct{ Item, From, To string }

// Account is what the change does to one account.
type Account struct {
	Name     string
	Dir      string
	Imported bool // registered in place (ccam import), so plain claude may use it
	Moves    []Move
	Plan     *share.Plan
	// DropMCP lists the servers ccam copied into .claude.json, removed
	// when the account stops sharing MCP or moves it to another group.
	DropMCP []string
	unsync  bool
	// SeedMCP copies the mcp.json of the group the account leaves into the
	// new group it gets MCP from, when that one has none yet: the account
	// keeps its servers, as it keeps the files of the other items.
	SeedMCP *[2]string
}

// Change is the work to get from one config to another.
type Change struct {
	Accounts []*Account
	paths    paths.Paths
}

// Build compares old and next (an edited copy of old) and plans the folder
// changes. The first content of a group comes from the group folder, then
// opt.Seed, then the default account, then the others by number.
func Build(p paths.Paths, old, next *config.Config, opt Options) (*Change, error) {
	order := make([]*config.Profile, 0, len(next.Profiles))
	for i := range next.Profiles {
		order = append(order, &next.Profiles[i])
	}
	rank := func(pr *config.Profile) int {
		switch pr.Name {
		case opt.Seed:
			return 0
		case next.Default:
			return 1
		}
		return 2
	}
	slices.SortStableFunc(order, func(a, b *config.Profile) int {
		return cmp.Or(cmp.Compare(rank(a), rank(b)), cmp.Compare(a.Slot, b.Slot))
	})
	ch := &Change{paths: p}
	pending := map[string]string{}
	for _, prof := range order {
		var from config.Share
		if o, ok := old.Find(prof.Name); ok {
			from = o.Share
		}
		if from.Equal(prof.Share) {
			continue
		}
		dir := profile.Dir(p, prof)
		pl, err := share.Build(dir, p.GroupsRoot(), from, prof.Share, share.Options{Fresh: opt.Fresh, Pending: pending})
		if err != nil {
			return nil, fmt.Errorf("account %s: %w", prof.Name, err)
		}
		for k, v := range pl.Created() {
			// Content an account brings beats an empty folder another makes.
			if old, ok := pending[k]; !ok || old == "" {
				pending[k] = v
			}
		}
		a := &Account{Name: prof.Name, Dir: dir, Imported: prof.Dir != "", Plan: pl}
		for _, item := range config.Items {
			if from.Get(item) == prof.Share.Get(item) {
				continue
			}
			a.Moves = append(a.Moves, Move{Item: item, From: groupName(old, prof.Name, item), To: groupName(next, prof.Name, item)})
		}
		if fm, tm := from.Get("mcp"), prof.Share.Get("mcp"); fm != config.Local && tm != fm {
			// The old group's servers copied into .claude.json would load
			// on, and under the flag strategy win over the new group's.
			a.unsync = true
			a.DropMCP = mcp.Copied(launcher.MCPStateDir(p, prof.Name))
			src, dst := mcp.Path(p.GroupDir(fm)), mcp.Path(p.GroupDir(tm))
			if _, planned := pending[dst]; tm != config.Local && !planned && fsutil.Exists(src) && !fsutil.Exists(dst) {
				a.SeedMCP = &[2]string{src, dst}
				pending[dst] = src
			}
		}
		ch.Accounts = append(ch.Accounts, a)
	}
	return ch, nil
}

func groupName(c *config.Config, account, item string) string {
	if g := c.GroupFor(account, item); g != nil {
		return g.Name
	}
	return ""
}

// Empty reports whether the change does nothing.
func (ch *Change) Empty() bool { return len(ch.Accounts) == 0 }

// Names lists the accounts the change touches.
func (ch *Change) Names() []string {
	out := make([]string, len(ch.Accounts))
	for i, a := range ch.Accounts {
		out[i] = a.Name
	}
	return out
}

// Check refuses a change that cannot be applied now: an account in it is
// running (it keeps writing into the folders being moved), or Windows
// cannot make the file links it needs.
func (ch *Change) Check(running func(name string) bool) error {
	var busy []string
	for _, a := range ch.Accounts {
		if running(a.Name) {
			busy = append(busy, a.Name)
		}
	}
	if len(busy) > 0 {
		verb := "is"
		if len(busy) > 1 {
			verb = "are"
		}
		return fmt.Errorf("%s %s running; quit Claude Code there first, then change the groups", strings.Join(busy, ", "), verb)
	}
	for _, a := range ch.Accounts {
		for _, s := range a.Plan.Steps {
			if s.Op != share.OpLink || s.Entry.Dir {
				continue
			}
			root := ch.paths.GroupsRoot()
			if err := os.MkdirAll(root, 0o700); err != nil {
				return err
			}
			if err := fsutil.CanLinkFiles(root); err != nil {
				return fmt.Errorf("%v, and %s needs a link for %s. Turn it on, or leave settings, memory and history out of the group", err, a.Name, s.Entry.Name)
			}
			return nil
		}
	}
	return nil
}

// Running is the usual check for Check.
func Running(p paths.Paths) func(string) bool {
	return func(name string) bool { return launcher.Running(p, name) }
}

// Conflict is a file conflict in one account's plan.
type Conflict struct {
	Account string
	*share.Conflict
}

func (ch *Change) Conflicts() []Conflict {
	var out []Conflict
	for _, a := range ch.Accounts {
		for _, c := range a.Plan.Conflicts {
			out = append(out, Conflict{Account: a.Name, Conflict: c})
		}
	}
	return out
}

// Describe writes the change, account by account, with paths relative to
// ccam's data folder.
func (ch *Change) Describe(w io.Writer) {
	if ch.Empty() {
		fmt.Fprintln(w, "  Nothing to change.")
		return
	}
	short := func(path string) string {
		if rel, err := filepath.Rel(ch.paths.Data, path); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
		return path
	}
	for _, a := range ch.Accounts {
		parts := make([]string, len(a.Moves))
		for i, m := range a.Moves {
			switch {
			case m.From == "":
				parts[i] = fmt.Sprintf("+%s (group %s)", m.Item, m.To)
			case m.To == "":
				parts[i] = fmt.Sprintf("-%s (group %s)", m.Item, m.From)
			default:
				parts[i] = fmt.Sprintf("%s: group %s → %s", m.Item, m.From, m.To)
			}
		}
		fmt.Fprintf(w, "%s: %s\n", a.Name, strings.Join(parts, ", "))
		if !a.Plan.Empty() {
			a.Plan.DescribeUnder(w, ch.paths.Data)
		}
		for _, m := range a.Moves {
			if m.Item == "mcp" && m.To != "" {
				fmt.Fprintf(w, "  MCP servers come from group %s at the next launch\n", m.To)
			}
		}
		if a.SeedMCP != nil {
			fmt.Fprintf(w, "  copy     %s → %s (the new group's first MCP servers)\n", short(a.SeedMCP[0]), short(a.SeedMCP[1]))
		}
		if len(a.DropMCP) > 0 {
			fmt.Fprintf(w, "  remove   MCP servers ccam copied into .claude.json: %s (a backup is kept)\n", strings.Join(a.DropMCP, ", "))
		}
		if a.Imported {
			fmt.Fprintf(w, "  Note: %s is an imported folder, so plain `claude` sees this change too if it uses it\n", a.Dir)
		}
	}
}

// Apply runs the plans in order. Save the new config before: if a plan stops
// halfway, `ccam doctor --fix` then finishes the change.
func (ch *Change) Apply() error {
	for _, a := range ch.Accounts {
		if err := a.Plan.Apply(); err != nil {
			return fmt.Errorf("account %s: %w\nThe originals are in %s. Fix the cause, then finish with `ccam doctor --fix`", a.Name, err, filepath.Join(a.Dir, share.BackupDir))
		}
		if a.SeedMCP != nil {
			if err := fsutil.CopyFile(a.SeedMCP[0], a.SeedMCP[1]); err != nil {
				return fmt.Errorf("account %s: could not copy the MCP servers: %w", a.Name, err)
			}
		}
		if a.unsync {
			if err := mcp.Unsync(a.Dir, launcher.MCPStateDir(ch.paths, a.Name)); err != nil {
				return fmt.Errorf("account %s: could not remove the copied MCP servers: %w", a.Name, err)
			}
		}
	}
	return nil
}
