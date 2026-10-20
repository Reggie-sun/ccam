package cli

import (
	"runtime"
	"testing"
)

func TestStatusLineCommand(t *testing.T) {
	if runtime.GOOS != "windows" {
		if got := statusLineCommand("/home/me/.local/bin/ccam", "/home/me"); got != "/home/me/.local/bin/ccam statusline" {
			t.Fatal(got)
		}
		if got := statusLineCommand("/opt/my tools/ccam", "/home/me"); got != "'/opt/my tools/ccam' statusline" {
			t.Fatal(got)
		}
		return
	}
	// Git Bash eats backslashes and PowerShell runs a quoted path as a
	// string, so: forward slashes, ~ for the home folder, quotes only when
	// there is no other way.
	for self, want := range map[string]string{
		`C:\Users\Jane Doe\.local\bin\ccam.exe`: "~/.local/bin/ccam.exe statusline",
		`D:\tools\ccam.exe`:                     "D:/tools/ccam.exe statusline",
		`D:\my tools\ccam.exe`:                  `"D:/my tools/ccam.exe" statusline`,
	} {
		if got := statusLineCommand(self, `C:\Users\Jane Doe`); got != want {
			t.Errorf("%s: %s", self, got)
		}
	}
}
