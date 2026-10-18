package app

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		argv []string
		want route
	}{
		{[]string{"ccs", "work", "--resume"}, route{kind: routeLaunch, sel: "work", args: []string{"--resume"}}},
		{[]string{"ccs", "2"}, route{kind: routeLaunch, sel: "2", args: []string{}}},
		{[]string{"/x/claude2", "--model", "opus", "--continue"}, route{kind: routeLaunch, sel: "2", args: []string{"--model", "opus", "--continue"}}},
		{[]string{"claude-work", "add"}, route{kind: routeLaunch, sel: "work", args: []string{"add"}}},
		{[]string{"ccs", "work", "--help"}, route{kind: routeLaunch, sel: "work", args: []string{"--help"}}},
		{[]string{"ccs", "-p", "요약해줘"}, route{kind: routeLaunch, args: []string{"-p", "요약해줘"}}},
		{[]string{"ccs", "work", "--", "--version"}, route{kind: routeLaunch, sel: "work", args: []string{"--", "--version"}}},
		{[]string{"ccs", "run", "work", "--", "-c"}, route{kind: routeLaunch, sel: "work", args: []string{"--", "-c"}}},
		{[]string{"ccs", "run", "-c"}, route{kind: routeLaunch, args: []string{"-c"}}},
		{[]string{"ccs", "--dry-run", "work", "-c"}, route{kind: routeLaunch, sel: "work", args: []string{"-c"}, dry: true}},
		{[]string{"ccs", "--dry-run"}, route{kind: routeLaunch, dry: true}},
		{[]string{"ccs"}, route{kind: routeCLI}},
		{[]string{"ccam", "list", "--json"}, route{kind: routeCLI, args: []string{"list", "--json"}}},
		{[]string{"ccam", "--help"}, route{kind: routeCLI, args: []string{"--help"}}},
		{[]string{"ccam", "--version"}, route{kind: routeCLI, args: []string{"--version"}}},
		{[]string{"ccs", "--dry-run", "next"}, route{kind: routeCLI, args: []string{"--dry-run", "next"}, dry: false}},
	}
	for _, c := range cases {
		got := parse(c.argv)
		if got.kind != c.want.kind || got.sel != c.want.sel || got.dry != c.want.dry || !sameArgs(got.args, c.want.args) {
			t.Errorf("%q:\n got %+v\nwant %+v", c.argv, got, c.want)
		}
	}
}

func sameArgs(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}
