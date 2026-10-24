package cli

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCallbackURL(t *testing.T) {
	for _, s := range []string{
		"http://localhost:34841/callback?code=abc&state=xyz\n",
		"  localhost:34841/callback?code=abc&state=xyz  ",
		"http://127.0.0.1:5000/callback?state=xyz&code=abc",
		"http://[::1]:5000/callback?code=abc&state=xyz",
	} {
		if _, ok := callbackURL(s); !ok {
			t.Errorf("not recognised: %q", s)
		}
	}
	for _, s := range []string{
		"",
		"abc123", // a code for the "Paste code here" prompt
		"http://localhost/callback?code=a&state=b", // no port
		"http://example.com:80/callback?code=a&state=b",
		"http://localhost.example.com:80/callback?code=a&state=b",
		"https://localhost:443/callback?code=a&state=b",
		"http://localhost:1/other?code=a&state=b",
		"http://localhost:1/callback?code=a",
		"http://user@localhost:1/callback?code=a&state=b",
		"file://localhost:1/callback?code=a&state=b",
	} {
		if _, ok := callbackURL(s); ok {
			t.Errorf("accepted: %q", s)
		}
	}
}

type pipe struct {
	bytes.Buffer
	closed bool
}

func (p *pipe) Close() error { p.closed = true; return nil }

func TestRelayLoginDeliversCallbackAndForwardsTheRest(t *testing.T) {
	var got url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		http.Redirect(w, r, "https://example.com/signed-in", http.StatusFound) // not followed
	}))
	defer srv.Close()

	in := bufio.NewReader(strings.NewReader("pasted-code\n" + srv.URL + "/callback?code=c1&state=s1\nlast"))
	var toClaude pipe
	var msg bytes.Buffer
	relayLogin(&lineSource{in: in}, nil, &toClaude, &msg, deliverCallback, nil)

	if got.Get("code") != "c1" || got.Get("state") != "s1" {
		t.Fatalf("callback got %v", got)
	}
	if toClaude.String() != "pasted-code\nlast" || !toClaude.closed {
		t.Fatalf("claude got %q (closed %v)", toClaude.String(), toClaude.closed)
	}
	if !strings.Contains(msg.String(), "Passed the sign-in") {
		t.Fatalf("message %q", msg.String())
	}
}

func TestRelayLoginReportsAFailedDelivery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad state", http.StatusBadRequest)
	}))
	defer srv.Close()
	in := bufio.NewReader(strings.NewReader(srv.URL + "/callback?code=c&state=s\n"))
	var toClaude pipe
	var msg bytes.Buffer
	relayLogin(&lineSource{in: in}, nil, &toClaude, &msg, deliverCallback, nil)
	if !strings.Contains(msg.String(), "HTTP 400") || toClaude.Len() != 0 {
		t.Fatalf("message %q, claude got %q", msg.String(), toClaude.String())
	}
}

func TestReadLineShowsStars(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("ab\x7fc\x1bd\r\x7fé\n\x04"))
	var stars bytes.Buffer
	for _, want := range []string{"acd\n", "é\n"} {
		if line, err := readLine(in, &stars); line != want || err != nil {
			t.Fatalf("got %q, %v; want %q", line, err, want)
		}
	}
	if line, err := readLine(in, &stars); line != "" || err != io.EOF {
		t.Fatalf("Ctrl-D: got %q, %v", line, err)
	}
	if got := stars.String(); got != "**\b \b**\n*\n" {
		t.Fatalf("stars %q", got)
	}
}
