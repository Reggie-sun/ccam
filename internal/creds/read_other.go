//go:build !darwin

package creds

// Read returns the stored login JSON from <dir>/.credentials.json.
func Read(dir, home string) ([]byte, Source, error) { return readFile(dir) }

// Has reports whether a login is stored.
func Has(dir, home string) bool {
	_, _, err := readFile(dir)
	return err == nil
}
