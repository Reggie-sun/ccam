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
		{[]string{"ccam", "work", "--resume"}, route{kind: routeLaunch, sel: "work", args: []string{"--resume"}}},
		{[]string{"ccam", "2"}, route{kind: routeLaunch, sel: "2", args: []string{}}},
		{[]string{"ccam", "work", "--help"}, route{kind: routeLaunch, sel: "work", args: []string{"--help"}}},
		{[]string{"ccam", "-p", "résume ça"}, route{kind: routeLaunch, args: []string{"-p", "résume ça"}}},
		{[]string{"ccam", "work", "--", "--version"}, route{kind: routeLaunch, sel: "work", args: []string{"--", "--version"}}},
		{[]string{"ccam", "run", "work", "--", "-c"}, route{kind: routeLaunch, sel: "work", args: []string{"--", "-c"}}},
		{[]string{"ccam", "run", "-c"}, route{kind: routeLaunch, args: []string{"-c"}}},
		{[]string{"ccam", "run", "--dry-run", "work", "-c"}, route{kind: routeLaunch, sel: "work", args: []string{"-c"}, dry: true}},
		{[]string{"ccam", "run", "--help"}, route{kind: routeCLI, args: []string{"help", "run"}}},
		{[]string{"ccam", "run", "work", "--help"}, route{kind: routeLaunch, sel: "work", args: []string{"--help"}}},
		{[]string{"ccam", "--dry-run", "work", "-c"}, route{kind: routeLaunch, sel: "work", args: []string{"-c"}, dry: true}},
		{[]string{"ccam", "--dry-run"}, route{kind: routeLaunch, dry: true}},
		{[]string{"ccam"}, route{kind: routeCLI}},
		{[]string{"ccam", "list", "--json"}, route{kind: routeCLI, args: []string{"list", "--json"}}},
		{[]string{"ccam", "--help"}, route{kind: routeCLI, args: []string{"--help"}}},
		{[]string{"ccam", "--version"}, route{kind: routeCLI, args: []string{"version"}}},
		{[]string{"ccam", "-v"}, route{kind: routeCLI, args: []string{"version"}}},
		{[]string{"ccam", "work", "-v"}, route{kind: routeLaunch, sel: "work", args: []string{"-v"}}},
		{[]string{"ccam", "__complete", "w"}, route{kind: routeCLI, args: []string{"__complete", "w"}}},
		{[]string{"ccam", "--dry-run", "next"}, route{kind: routeCLI, args: []string{"--dry-run", "next"}, dry: false}},
		{[]string{"ccam", "--fallback", "work", "-p", "hi"}, route{kind: routeLaunch, sel: "work", args: []string{"-p", "hi"}, fallback: true}},
		{[]string{"ccam", "--fallback", "-p", "hi"}, route{kind: routeLaunch, args: []string{"-p", "hi"}, fallback: true}},
		{[]string{"ccam", "--dry-run", "--fallback", "work", "-p"}, route{kind: routeLaunch, sel: "work", args: []string{"-p"}, dry: true, fallback: true}},
		{[]string{"ccam", "--fallback", "--dry-run", "work", "-p"}, route{kind: routeLaunch, sel: "work", args: []string{"-p"}, dry: true, fallback: true}},
		{[]string{"ccam", "run", "--fallback", "work", "-p"}, route{kind: routeLaunch, sel: "work", args: []string{"-p"}, fallback: true}},
		{[]string{"ccam", "work", "--fallback"}, route{kind: routeLaunch, sel: "work", args: []string{"--fallback"}}}, // claude's argument
		{[]string{"ccam", "--fallback", "usage"}, route{kind: routeCLI, args: []string{"--fallback"}}},
	}
	for _, c := range cases {
		got := parse(c.argv)
		if got.kind != c.want.kind || got.sel != c.want.sel || got.dry != c.want.dry || got.fallback != c.want.fallback || !sameArgs(got.args, c.want.args) {
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
