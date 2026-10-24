//go:build !windows

package main

import "os"

// ucrtArgs is os.Args[1:]: only Windows programs parse their command line.
func ucrtArgs() []string { return os.Args[1:] }
