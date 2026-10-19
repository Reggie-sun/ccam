package cli

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// Signing in on a remote machine (over SSH) breaks the browser half of
// `claude auth login`: Claude Code waits for the sign-in on a callback server
// at http://localhost:<port>/callback on that machine, but the browser that
// finishes the sign-in runs on the user's own computer, where the address
// leads nowhere. The address still carries the one-time code, so ccam lets
// the user paste it and requests it on the machine where Claude Code listens.

// runLogin runs `claude auth login` in the profile. Lines typed meanwhile go
// to claude (for its "Paste code here" prompt), except a pasted callback
// address, which ccam requests itself.
func runLogin(e *env, cfg *config.Config, prof *config.Profile) error {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return err
	}
	c := exec.Command(bin, "auth", "login")
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, prof), prof.Name)
	c.Stdout, c.Stderr = e.stdout, e.stderr
	toClaude, err := c.StdinPipe()
	if err != nil {
		return err
	}
	fmt.Fprintln(e.stderr, "On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.")
	if err := c.Start(); err != nil {
		return err
	}
	// The relay may stay blocked on a read after claude exits; nothing reads
	// the terminal after a login, and the line it may still take is not lost
	// to anyone.
	go relayLogin(e.in, toClaude, e.stderr, deliverCallback)
	return c.Wait()
}

// relayLogin copies lines from in to claude until in ends, delivering pasted
// callback addresses instead of copying them.
func relayLogin(in *bufio.Reader, toClaude io.WriteCloser, msg io.Writer, deliver func(*url.URL) error) {
	defer toClaude.Close()
	for {
		line, err := in.ReadString('\n')
		if u, ok := callbackURL(line); ok {
			if derr := deliver(u); derr != nil {
				fmt.Fprintf(msg, "ccam: could not pass the sign-in to Claude Code: %v\n", derr)
			} else {
				fmt.Fprintln(msg, "→ Passed the sign-in to Claude Code on this machine.")
			}
		} else if line != "" {
			if _, werr := io.WriteString(toClaude, line); werr != nil {
				return // claude has exited
			}
		}
		if err != nil {
			return
		}
	}
}

// callbackURL recognises the address Claude Code's sign-in redirects to:
// http://localhost:<port>/callback?code=…&state=… (also 127.0.0.1 or [::1],
// with or without http://). Nothing else is ever requested, so a pasted line
// cannot make ccam contact another host.
func callbackURL(line string) (*url.URL, bool) {
	s := strings.TrimSpace(line)
	if !strings.Contains(s, "://") {
		s = "http://" + s
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "http" || u.Port() == "" || u.Path != "/callback" || u.User != nil {
		return nil, false
	}
	switch u.Hostname() {
	case "localhost", "127.0.0.1", "::1":
	default:
		return nil, false
	}
	q := u.Query()
	if q.Get("code") == "" || q.Get("state") == "" {
		return nil, false
	}
	return u, true
}

// deliverCallback requests the callback address the way the browser would
// have. Claude Code answers with a redirect to a "signed in" page, which is
// not followed.
func deliverCallback(u *url.URL) error {
	client := &http.Client{
		Timeout:       10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	resp, err := client.Get(u.String())
	if err != nil {
		return fmt.Errorf("%v (is the login still waiting? Start it again if it ended)", err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("the sign-in was answered with HTTP %d; the code may have expired, so start the login again", resp.StatusCode)
	}
	return nil
}
