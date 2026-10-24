package main

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

// ucrtArgs is os.Args[1:], or with FAKE_ARGV=ucrt the arguments as a
// program built with Microsoft's C runtime reads its command line (node.exe,
// which an npm install's claude.cmd runs): unlike Go, it takes "" inside
// quotes as one quote and stays in quotes.
func ucrtArgs() []string {
	if os.Getenv("FAKE_ARGV") != "ucrt" {
		return os.Args[1:]
	}
	r := []rune(windows.UTF16PtrToString(windows.GetCommandLine()))
	i := 0
	// The program name: up to the next quote, or the next blank.
	if i < len(r) && r[i] == '"' {
		for i++; i < len(r) && r[i] != '"'; i++ {
		}
		i++
	} else {
		for ; i < len(r) && r[i] != ' ' && r[i] != '\t'; i++ {
		}
	}
	var args []string
	for {
		for ; i < len(r) && (r[i] == ' ' || r[i] == '\t'); i++ {
		}
		if i >= len(r) {
			return args
		}
		var arg strings.Builder
		inquote := false
		for i < len(r) {
			c := r[i]
			switch {
			case !inquote && (c == ' ' || c == '\t'):
			case c == '\\':
				n := 0
				for ; i < len(r) && r[i] == '\\'; i++ {
					n++
				}
				if i < len(r) && r[i] == '"' {
					arg.WriteString(strings.Repeat(`\`, n/2))
					if n%2 == 1 {
						arg.WriteByte('"')
						i++
					}
				} else {
					arg.WriteString(strings.Repeat(`\`, n))
				}
				continue
			case c == '"' && inquote && i+1 < len(r) && r[i+1] == '"':
				arg.WriteByte('"')
				i += 2
				continue
			case c == '"':
				inquote = !inquote
				i++
				continue
			default:
				arg.WriteRune(c)
				i++
				continue
			}
			break
		}
		args = append(args, arg.String())
	}
}
