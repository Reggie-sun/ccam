package cli

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func (c *testCLI) env() (*env, *config.Config) {
	e := &env{paths: c.p, stdin: strings.NewReader(""), in: bufio.NewReader(strings.NewReader("")), stdout: os.Stdout, stderr: os.Stderr}
	cfg, err := e.load()
	if err != nil {
		c.t.Fatal(err)
	}
	return e, cfg
}

func TestFallbackTargetPrefersRoomThenUnknown(t *testing.T) {
	cases := []struct {
		name  string
		usage map[string]int
		tried []string
		want  string
	}{
		{"lowest 5-hour usage", map[string]int{"personal": 100, "work": 60, "side": 20}, []string{"personal"}, "side"},
		{"full accounts are skipped", map[string]int{"personal": 100, "work": 100, "side": 30}, []string{"personal"}, "side"},
		{"unknown usage comes after room", map[string]int{"personal": 100, "work": -1, "side": 90}, []string{"personal"}, "side"},
		{"unknown usage is still tried", map[string]int{"personal": 100, "work": -1, "side": 100}, []string{"personal"}, "work"},
		{"tried accounts are not", map[string]int{"personal": 100, "work": 10, "side": 100}, []string{"personal", "work"}, ""},
		{"logged-out accounts are not", map[string]int{"personal": 100}, []string{"personal"}, ""},
	}
	for _, tc := range cases {
		c := newTestCLI(t, tc.usage)
		e, cfg := c.env()
		tried := map[string]bool{}
		for _, n := range tc.tried {
			tried[n] = true
		}
		from, _ := cfg.Find(tc.tried[0])
		got := ""
		if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "hi"}, usage.Policy{}); p != nil {
			got = p.Name
		}
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFallbackTargetKeepsHistoryWhenResuming(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 50, "side": 10})
	e, cfg := c.env()
	cfg.Profiles[2].Share.History = config.Local // side does not share history
	from := &cfg.Profiles[0]
	tried := map[string]bool{"personal": true}
	if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "go on", "--continue"}, usage.Policy{}); p == nil || p.Name != "work" {
		t.Fatalf("--continue must stay in the history group: %v", p)
	}
	if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "new"}, usage.Policy{}); p == nil || p.Name != "side" {
		t.Fatalf("a new run may use any account: %v", p)
	}
}

func TestLimitedByMessageOrUsage(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 40})
	e, cfg := c.env()
	personal, work := &cfg.Profiles[0], &cfg.Profiles[1]
	if !limited(e, cfg, work, []byte("You've hit your limit · resets 3pm"), nil) {
		t.Error("message on stdout")
	}
	if !limited(e, cfg, work, nil, []byte("Claude AI usage limit reached|1760000000")) {
		t.Error("older message on stderr")
	}
	if !limited(e, cfg, personal, []byte("API Error: 500"), nil) {
		t.Error("full by usage")
	}
	if limited(e, cfg, work, []byte("API Error: 500"), nil) {
		t.Error("another failure with room left is not a limit")
	}
}

func TestPrintMode(t *testing.T) {
	for args, want := range map[string]bool{
		"-p hi": true, "--print hi": true, "hi -p": true, "--resume": false, "-- -p": false, "": false,
	} {
		if got := printMode(strings.Fields(args)); got != want {
			t.Errorf("%q: %v", args, got)
		}
	}
}
