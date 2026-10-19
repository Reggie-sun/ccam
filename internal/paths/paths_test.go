package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCCAMHomeOverridesEverything(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CCAM_HOME", root)
	t.Setenv("XDG_CONFIG_HOME", "/nope")
	p := FromEnv()
	for _, dir := range []string{p.Config, p.Data, p.State, p.Bin, p.Runtime, p.Home} {
		if !strings.HasPrefix(dir, root+string(filepath.Separator)) {
			t.Errorf("%s escapes CCAM_HOME %s", dir, root)
		}
	}
}

func TestXDG(t *testing.T) {
	t.Setenv("CCAM_HOME", "")
	t.Setenv("HOME", "/h")
	t.Setenv("XDG_CONFIG_HOME", "/cfg")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("XDG_RUNTIME_DIR", "/run/user/1")
	p := FromEnv()
	want := Paths{Config: "/cfg/ccam", Data: "/h/.local/share/ccam", State: "/h/.local/state/ccam", Bin: "/h/.local/bin", Runtime: "/run/user/1/ccam", Home: "/h"}
	if p != want {
		t.Fatalf("got %+v want %+v", p, want)
	}
}

func TestRelativeXDGIsIgnored(t *testing.T) {
	t.Setenv("CCAM_HOME", "")
	t.Setenv("HOME", "/h")
	t.Setenv("XDG_DATA_HOME", "rel/data")
	t.Setenv("XDG_RUNTIME_DIR", "rel/run")
	p := FromEnv()
	if p.Data != "/h/.local/share/ccam" || p.Runtime != p.State+"/locks" {
		t.Fatalf("got %+v", p)
	}
}

func TestUnsetHomeStaysAbsolute(t *testing.T) {
	t.Setenv("CCAM_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	if p := FromEnv(); !filepath.IsAbs(p.Config) || !filepath.IsAbs(p.Bin) {
		t.Fatalf("relative paths without HOME: %+v", p)
	}
}
