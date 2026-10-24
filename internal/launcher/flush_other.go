//go:build !linux && !darwin

package launcher

import "os"

// discardInput leaves input alone here.
func discardInput(*os.File) bool { return false }
