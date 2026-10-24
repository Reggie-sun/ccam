package cli

import "testing"

// A `go install …@vX` build names its release in the build info; any other
// module version is a development build.
func TestReleaseOf(t *testing.T) {
	for in, want := range map[string]string{
		"v0.10.3":                               "0.10.3",
		"v1.2.0-rc.1":                           "1.2.0-rc.1",
		"(devel)":                               "dev",
		"":                                      "dev",
		"v0.10.3+dirty":                         "dev",
		"v0.10.4-0.20261022010203-c90ba8e12345": "dev",
		"v0.10.4-0.20261022010203-c90ba8e12345+dirty": "dev",
		"v0.0.0-20261022010203-c90ba8e12345":          "dev",
		"v1.2.4-rc.1.0.20261022010203-c90ba8e12345":   "dev",
	} {
		if got := releaseOf(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}

func TestVersion(t *testing.T) {
	if v := version(); v != "dev" {
		t.Fatalf("a test binary is a development build: %q", v)
	}
	old := Version
	t.Cleanup(func() { Version = old })
	Version = "0.9.0"
	if v := version(); v != "0.9.0" {
		t.Fatalf("the ldflags value wins: %q", v)
	}
}
