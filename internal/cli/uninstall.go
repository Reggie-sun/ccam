package cli

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/update"
)

func uninstallCmd(e *env) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove ccam: every account, the shortcut commands, ccam's folders and ccam itself",
		Long: "Removes every account as `ccam remove` does, then deletes ccam's folders (config,\n" +
			"data with the trash, state) and the ccam command.\n\n" +
			"An account in ccam's folders is logged out first (not when plain claude holds the\n" +
			"same login) and deleted with its folder: its conversations and settings too.\n" +
			"A folder registered in place (ccam import), such as ~/.claude, stays for plain\n" +
			"claude: it leaves its groups with a copy of what they share, and the ccam status\n" +
			"line in its settings.json gives way to the one it wrapped, if any.\n\n" +
			"A ccam installed from the .deb package stays: remove it with `sudo apt remove ccam`.\n" +
			"On Windows a running program cannot delete itself: ccam says which file to delete.",
		Example: "  ccam --dry-run uninstall   # show what would be removed\n" +
			"  ccam uninstall",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(e, yes)
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "uninstall without asking")
	return cmd
}

func runUninstall(e *env, yes bool) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	if err := group.Unfinished(e.paths); err != nil {
		return err
	}
	for i := range cfg.Profiles {
		if err := refuseRunning(e, &cfg.Profiles[i]); err != nil {
			return err
		}
	}
	self, err := ccamPath()
	if err != nil {
		return err
	}
	owned := update.New(version(), self, runtime.GOOS, runtime.GOARCH).Managed()
	e.printf("This removes ccam:\n")
	var names, shortcuts []string
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		names = append(names, p.Name)
		dir := profile.Dir(e.paths, p)
		switch {
		case p.Dir != "" && !p.Share.Equal(config.Share{}):
			e.printf("  %s: %s stays; it leaves its groups with a copy of what they share\n", p.Name, dir)
		case p.Dir != "":
			e.printf("  %s: %s stays\n", p.Name, dir)
		case p.API == nil && storedLogin(e, p) && sameLogin(e, cfg, p) != profile.PlainClaudeName:
			e.printf("  %s: logged out, then deleted with its folder %s\n", p.Name, dir)
		default:
			e.printf("  %s: deleted with its folder %s\n", p.Name, dir)
		}
		if p.Shortcut != "" {
			shortcuts = append(shortcuts, p.Shortcut)
		}
	}
	if len(shortcuts) > 0 {
		e.printf("  shortcut commands: %s\n", strings.Join(shortcuts, ", "))
	}
	for _, f := range ccamFolders(e) {
		e.printf("  folder %s\n", f)
	}
	switch {
	case owned:
		e.printf("  %s stays: a package installed it (sudo apt remove ccam)\n", self)
	case runtime.GOOS == "windows":
		e.printf("  %s: delete it yourself once this ends\n", self)
	default:
		e.printf("  the command %s\n", self)
	}
	if e.dryRun {
		return nil
	}
	if !yes && !e.confirm("Uninstall ccam? Deleted conversations and logins cannot be brought back.") {
		e.printf("Cancelled.\n")
		return nil
	}

	for i, name := range names {
		p := &cfg.Profiles[i]
		imported, dir := p.Dir != "", profile.Dir(e.paths, p)
		msg, err := removeAccount(e, name)
		if err != nil {
			return fmt.Errorf("%w\nccam is still installed; run `ccam uninstall` again once that is fixed", err)
		}
		e.printf("✓ %s\n", msg)
		if !imported {
			continue
		}
		// Plain claude goes on using the folder, without ccam.
		settings := filepath.Join(dir, "settings.json")
		if changed, err := removeStatusLine(settings); err != nil {
			fmt.Fprintf(e.stderr, "ccam: the ccam status line stays in %s: %v\n", settings, err)
		} else if changed {
			e.printf("✓ ccam status line taken out of %s\n", settings)
		}
		for _, it := range ext.ScanAccount(name, dir, e.paths.GroupsRoot(), config.Share{}, nil, "").Plugins {
			if n, _, _ := strings.Cut(it.Name, "@"); n == "ccam" && !it.Off {
				e.printf("  %s still has the ccam plugin; in plain claude: /plugin uninstall %s\n", dir, it.Name)
			}
		}
	}
	// Again: removing the accounts takes locks, which may make one.
	for _, f := range ccamFolders(e) {
		if err := os.RemoveAll(f); err != nil {
			return err
		}
		e.printf("✓ %s deleted\n", f)
	}
	switch {
	case owned:
		e.printf("ccam's data is gone; remove the package with: sudo apt remove ccam\n")
	case runtime.GOOS == "windows":
		e.printf("ccam's data is gone; delete the command with: Remove-Item %s\n", launcher.PSString(self))
	default:
		if err := os.Remove(self); err != nil {
			return err
		}
		e.printf("✓ %s deleted\n", self)
	}
	return nil
}

// ccamFolders are the folders of ccam's own that exist, none inside another.
func ccamFolders(e *env) []string {
	var out []string
	for _, f := range []string{e.paths.Config, e.paths.Data, e.paths.State, e.paths.Runtime} {
		if !fsutil.Exists(f) {
			continue
		}
		inside := false
		for _, o := range out {
			if rel, err := filepath.Rel(o, f); err == nil && !strings.HasPrefix(rel, "..") {
				inside = true
			}
		}
		if !inside {
			out = append(out, f)
		}
	}
	return out
}

// removeStatusLine takes the ccam status line out of the settings file at
// path, putting back the command it wraps, if any. A wrapped command it
// cannot read is left as it is.
func removeStatusLine(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	doc := map[string]any{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return false, err
	}
	sl, _ := doc["statusLine"].(map[string]any)
	command, _ := sl["command"].(string)
	rest, ours := ccamStatusLine(command)
	if !ours {
		return false, nil
	}
	switch wrapped, ok := wrappedStatusLine(rest); {
	case !ok:
		return false, nil
	case wrapped == "":
		delete(doc, "statusLine")
	default:
		sl["command"] = wrapped
	}
	out, err := fsutil.MarshalJSON(doc)
	if err != nil {
		return false, err
	}
	return true, fsutil.AtomicWrite(path, out, 0o600)
}

// wrappedStatusLine reads the options after `ccam statusline` as
// installStatusLine writes them: the command they wrap, "" for none.
func wrappedStatusLine(opts string) (string, bool) {
	if opts == "" {
		return "", true
	}
	if s, ok := strings.CutPrefix(opts, " --wrap-base64 "); ok {
		b, err := base64.StdEncoding.DecodeString(s)
		return string(b), err == nil
	}
	if s, ok := strings.CutPrefix(opts, " --wrap "); ok {
		w, rest, ok := shWord(s + " ")
		return w, ok && rest == " "
	}
	return "", false
}
