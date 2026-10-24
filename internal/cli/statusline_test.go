package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestStatusLineCommand(t *testing.T) {
	if runtime.GOOS != "windows" {
		if got := statusLineCommand("/home/me/.local/bin/ccam", "/home/me", "sh"); got != "/home/me/.local/bin/ccam statusline" {
			t.Fatal(got)
		}
		if got := statusLineCommand("/opt/my tools/ccam", "/home/me", "sh"); got != "'/opt/my tools/ccam' statusline" {
			t.Fatal(got)
		}
		return
	}
	// Git Bash eats backslashes and PowerShell runs a quoted path as a
	// string, so: forward slashes, ~ for the home folder, quotes only when
	// there is no other way, as the shell Claude Code uses writes them.
	bash, ps := `C:\Program Files\Git\bin\bash.exe`, "powershell"
	for _, c := range []struct{ self, sh, want string }{
		{`C:\Users\Jane Doe\.local\bin\ccam.exe`, ps, "~/.local/bin/ccam.exe statusline"},
		{`C:\Users\Jane Doe\.local\bin\ccam.exe`, bash, "~/.local/bin/ccam.exe statusline"},
		{`D:\tools\ccam.exe`, ps, "D:/tools/ccam.exe statusline"},
		{`D:\my tools\ccam.exe`, bash, `'D:/my tools/ccam.exe' statusline`},
		{`D:\my tools\ccam.exe`, ps, `& 'D:/my tools/ccam.exe' statusline`},
		{`C:\Users\Jane Doe\it's $here\ccam.exe`, bash, `'C:/Users/Jane Doe/it'\''s $here/ccam.exe' statusline`},
		{`C:\Users\Jane Doe\it's $here\ccam.exe`, ps, `& 'C:/Users/Jane Doe/it''s $here/ccam.exe' statusline`},
		{"D:\\don\u2019t\\ccam.exe", ps, "& 'D:/don\u2019\u2019t/ccam.exe' statusline"},
	} {
		if got := statusLineCommand(c.self, `C:\Users\Jane Doe`, c.sh); got != c.want {
			t.Errorf("%s in %s: %s", c.self, c.sh, got)
		}
	}
}

// Install wraps the command of the existing status line and keeps its
// other options.
func TestInstallStatusLineKeepsOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	os.WriteFile(path, []byte(`{"statusLine":{"type":"command","command":"printf hi","padding":2,"refreshInterval":5,"hideVimModeIndicator":true},"theme":"dark"}`), 0o600)
	if changed, err := installStatusLine(path, "/bin/ccam statusline"); !changed || err != nil {
		t.Fatal(changed, err)
	}
	var doc map[string]any
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, &doc)
	sl := doc["statusLine"].(map[string]any)
	want := "/bin/ccam statusline --wrap 'printf hi'"
	if runtime.GOOS == "windows" {
		want = "/bin/ccam statusline --wrap-base64 cHJpbnRmIGhp" // no quotes for Git Bash and PowerShell to differ on
	}
	if sl["command"] != want || sl["type"] != "command" || sl["padding"] != 2.0 ||
		sl["refreshInterval"] != 5.0 || sl["hideVimModeIndicator"] != true || doc["theme"] != "dark" {
		t.Fatalf("%s", data)
	}
	if changed, _ := installStatusLine(path, "/bin/ccam statusline"); changed {
		t.Fatal("a second install changed it again")
	}
}

// A wrapped status line shows as Claude Code would show it: every line
// that is not blank, trimmed, however long the command takes, and nothing
// but the label when it fails.
func TestStatuslineWrapKeepsRowsLikeClaudeCode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sh commands")
	}
	c := newTestCLI(t, nil)
	t.Setenv("CCAM_PROFILE", "work")
	for command, want := range map[string]string{
		"echo mine":                     "[work] mine\n",
		`printf 'line1\n\n  line2  \n'`: "[work] line1\nline2\n",
		`printf '\nreal\n'`:             "[work] real\n",
		"sleep 2.5; echo slow-but-fine": "[work] slow-but-fine\n",
		"echo partial; exit 3":          "[work]\n",
		"true":                          "[work]\n",
		"sleep 3 & echo bg":             "[work] bg\n",
	} {
		if out, _, _ := c.run("{}", "statusline", "--wrap", command); out != want {
			t.Errorf("%s: %q, want %q", command, out, want)
		}
	}
}

// A status line that runs a ccam already, written by an older ccam (other
// quoting) or for a ccam at another path, is replaced, keeping what it
// wraps: wrapping it would run ccam inside ccam, and show the label twice.
func TestInstallStatusLineReplacesAnOlderCcamOne(t *testing.T) {
	ours := "'/opt/my tools/ccam' statusline"
	for old, want := range map[string]string{
		`"/opt/my tools/ccam" statusline`:                                ours,
		`"C:/Program Files/ccam/ccam.exe" statusline --wrap-base64 eA==`: ours + " --wrap-base64 eA==",
		`& 'C:/Program Files/ccam/ccam.exe' statusline`:                  ours,
		`~/.local/bin/ccam statusline --wrap 'printf hi'`:                ours + " --wrap 'printf hi'",
		`'/home/o'\''neil/old/ccam' statusline --wrap 'my-line'`:         ours + " --wrap 'my-line'",
		`& 'C:/Users/O''Neil/tools/ccam.exe' statusline`:                 ours,
		"& 'C:/Users/O\u2019\u2019Neil/tools/ccam.exe' statusline":       ours,
		"/usr/bin/ccamx statusline":                                      "wrapped",
		`'/opt/o'\''ccam' statusline`:                                    "wrapped",
		"ccam statuslines":                                               "wrapped",
	} {
		path := filepath.Join(t.TempDir(), "settings.json")
		data, _ := json.Marshal(map[string]any{"statusLine": map[string]any{"type": "command", "command": old}})
		os.WriteFile(path, data, 0o600)
		if _, err := installStatusLine(path, ours); err != nil {
			t.Fatal(err)
		}
		var doc struct {
			StatusLine struct{ Command string } `json:"statusLine"`
		}
		data, _ = os.ReadFile(path)
		json.Unmarshal(data, &doc)
		got := doc.StatusLine.Command
		if want == "wrapped" {
			if !strings.HasPrefix(got, ours+" --wrap") {
				t.Errorf("%s: %s, want it wrapped", old, got)
			}
		} else if got != want {
			t.Errorf("%s: %s, want %s", old, got, want)
		}
	}
}
