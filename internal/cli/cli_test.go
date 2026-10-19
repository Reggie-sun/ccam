package cli

import (
	"os"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

// The launcher hands only config.Reserved names to cobra, so a command
// missing there would be launched as an account name instead.
func TestReservedMatchesCommands(t *testing.T) {
	root := newRoot(&env{})
	names := map[string]bool{"help": true} // cobra adds help itself
	for _, c := range root.Commands() {
		names[c.Name()] = true
		for _, a := range c.Aliases {
			names[a] = true
		}
	}
	for n := range names {
		if !config.Reserved[n] {
			t.Errorf("command %q is not in config.Reserved", n)
		}
	}
	for n := range config.Reserved {
		if !names[n] {
			t.Errorf("config.Reserved has %q, which is no command", n)
		}
	}
}

func TestDevNullIsNotATerminal(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if isTerminal(f) {
		t.Fatal("/dev/null counted as a terminal; `ccam add work >/dev/null` would start the wizard")
	}
}
