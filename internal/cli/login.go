package cli

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/unix"

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
	stars, restore := maskTyping(e)
	defer restore()
	// The relay may stay blocked on a read after claude exits; nothing reads
	// the terminal after a login, and the line it may still take is not lost
	// to anyone.
	go relayLogin(e.in, stars, toClaude, e.stderr, deliverCallback)
	return c.Wait()
}

// maskTyping turns the terminal's echo off for the login so that the relay
// can show each character typed as a '*' on the returned writer. It returns
// a nil writer when stdin or stderr is not a terminal, and a function that
// restores the terminal.
func maskTyping(e *env) (io.Writer, func()) {
	if !isTerminal(e.stdin) || !isTerminal(e.stderr) {
		return nil, func() {}
	}
	fd := e.stdin.(*os.File).Fd()
	saved, err := term.GetState(fd)
	if err != nil {
		return nil, func() {}
	}
	masked := *saved
	masked.Lflag &^= unix.ECHO | unix.ICANON
	masked.Cc[unix.VMIN], masked.Cc[unix.VTIME] = 1, 0
	if term.SetState(fd, &masked) != nil {
		return nil, func() {}
	}
	// Ctrl-C still ends ccam, but not before the echo is back on.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	var once sync.Once
	restore := func() {
		once.Do(func() {
			signal.Stop(sig)
			close(sig)
			_ = term.Restore(fd, saved)
		})
	}
	go func() {
		if s, ok := <-sig; ok {
			restore()
			_ = syscall.Kill(os.Getpid(), s.(syscall.Signal))
		}
	}()
	return e.stderr, restore
}

// relayLogin copies lines from in to claude until in ends, delivering pasted
// callback addresses instead of copying them. With stars set, the terminal
// does not echo and each character typed is shown there as a '*'.
func relayLogin(in *bufio.Reader, stars io.Writer, toClaude io.WriteCloser, msg io.Writer, deliver func(*url.URL) error) {
	defer toClaude.Close()
	for {
		line, err := readLine(in, stars)
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

// readLine reads a line from in, newline included. With stars set, the
// terminal is not in line mode, so readLine does the editing itself: Enter
// ends the line, Backspace erases, Ctrl-D on an empty line ends the input.
func readLine(in *bufio.Reader, stars io.Writer) (string, error) {
	if stars == nil {
		return in.ReadString('\n')
	}
	var line []rune
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			return string(line), err
		}
		switch {
		case r == '\r' || r == '\n':
			fmt.Fprint(stars, "\n")
			return string(line) + "\n", nil
		case r == 0x7f || r == '\b':
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Fprint(stars, "\b \b")
			}
		case r == 0x04:
			if len(line) == 0 {
				return "", io.EOF
			}
		case r < ' ':
			// other control keys: ignored
		default:
			line = append(line, r)
			fmt.Fprint(stars, "*")
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
