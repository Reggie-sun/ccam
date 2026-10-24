package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// record is a Claude Code record (sessions/<pid>.json, or daemon.lock) with
// the given fields.
func record(fields map[string]any) string {
	data, _ := json.Marshal(fields)
	return string(data)
}

// sleeper starts a process that runs until the test ends, and returns its
// pid and the start time Claude Code would record for it.
func sleeper(t *testing.T) (*exec.Cmd, string) {
	t.Helper()
	self, _ := os.Executable()
	c := exec.Command(self, "-test.run=^TestHelperSleeps$")
	c.Env = append(os.Environ(), "SESSION_HELPER_SLEEP=1")
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Process.Kill(); c.Wait() })
	return c, claudeStart(t, c.Process.Pid)
}

// finished is a process that has ended, and the start time Claude Code
// recorded for it while it ran: what a session that crashed leaves behind.
func finished(t *testing.T) (int, string) {
	t.Helper()
	c, start := sleeper(t)
	c.Process.Kill()
	c.Wait()
	return c.Process.Pid, start
}

// TestHelperSleeps is the process sleeper starts.
func TestHelperSleeps(t *testing.T) {
	if os.Getenv("SESSION_HELPER_SLEEP") == "" {
		t.Skip("only run by sleeper")
	}
	time.Sleep(time.Minute)
}

// zombie runs a process under a helper that never reaps it, and returns its
// pid and what ends it: after that it has ended, but stays a zombie while
// the helper runs, as an orphan does in a container whose PID 1 reaps
// nothing.
func zombie(t *testing.T) (int, func()) {
	t.Helper()
	self, _ := os.Executable()
	c := exec.Command(self, "-test.run=^TestHelperLeavesAZombie$")
	c.Env = append(os.Environ(), "SESSION_HELPER_ZOMBIE=1")
	in, _ := c.StdinPipe()
	out, _ := c.StdoutPipe()
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Process.Kill(); c.Wait() })
	r := bufio.NewReader(out)
	line, _ := r.ReadString('\n')
	pid, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		t.Fatalf("the helper said %q", line)
	}
	// Ended while the helper still holds it, so that its pid names nothing else.
	t.Cleanup(func() {
		if p, err := os.FindProcess(pid); err == nil {
			p.Kill()
		}
	})
	return pid, func() {
		in.Close()
		if line, _ := r.ReadString('\n'); line != "ended\n" {
			t.Fatalf("the helper said %q", line)
		}
	}
}

// TestHelperLeavesAZombie is the helper zombie starts: it runs a sleeper and
// says its pid; once its input ends, it ends the sleeper without reaping it,
// says so, and waits.
func TestHelperLeavesAZombie(t *testing.T) {
	if os.Getenv("SESSION_HELPER_ZOMBIE") == "" {
		t.Skip("only run by zombie")
	}
	c, _ := sleeper(t)
	fmt.Println(c.Process.Pid)
	io.Copy(io.Discard, os.Stdin)
	c.Process.Kill()
	waitEnded(t, c)
	fmt.Println("ended")
	time.Sleep(time.Minute)
}

// A process that has ended no longer runs, also while its parent has not
// reaped it (a zombie): neither ccam's marker nor Claude Code's record of
// it counts then.
func TestUnreapedProcessesDoNotRun(t *testing.T) {
	if _, ok := StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	pid, end := zombie(t)
	now, ok := StartTime(pid)
	if !ok {
		t.Fatal("no start time while it runs")
	}
	markers, dir := t.TempDir(), t.TempDir()
	os.WriteFile(filepath.Join(markers, strconv.Itoa(pid)), []byte(now), 0o600)
	os.Mkdir(filepath.Join(dir, "sessions"), 0o700)
	os.WriteFile(filepath.Join(dir, "sessions", strconv.Itoa(pid)+".json"),
		[]byte(record(map[string]any{"pid": pid, "procStart": claudeStart(t, pid), "startedAt": time.Now().UnixMilli()})), 0o644)
	if Count(markers) != 1 || !ClaudeRuns(dir) {
		t.Fatal("it does not count while it runs")
	}
	end()
	_, started := StartTime(pid)
	if n, runs := Count(markers), ClaudeRuns(dir); started || n != 0 || runs {
		t.Fatalf("after it ended: start time %v, %d marker(s), Claude Code's record counts %v", started, n, runs)
	}
}

func TestClaudeRunsReadsClaudeCodesRecords(t *testing.T) {
	if _, ok := StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	self := os.Getpid()
	now := claudeStart(t, self)
	mine, _ := StartTime(self)
	begun, _ := wallStart(mine)
	after := time.Now().UnixMilli()             // written after this process started
	before := begun.Add(-time.Hour).UnixMilli() // written before it: the PID was reused since
	dead, deadStart := finished(t)
	session := "sessions/" + strconv.Itoa(self) + ".json"
	for _, c := range []struct {
		name, file string
		fields     map[string]any
		want       bool
	}{
		{"a session", session, map[string]any{"pid": self, "kind": "interactive", "procStart": now, "startedAt": after}, true},
		{"a reused pid", session, map[string]any{"pid": self, "procStart": otherStart, "startedAt": after}, false},
		{"a session that ended", "sessions/" + strconv.Itoa(dead) + ".json", map[string]any{"pid": dead, "procStart": deadStart}, false},
		{"no start time, written after the process started", session, map[string]any{"pid": self, "startedAt": after}, true},
		{"no start time, written before the process started", session, map[string]any{"pid": self, "startedAt": before}, false},
		{"an unknown start time, written after", session, map[string]any{"pid": self, "procStart": "soon", "startedAt": after}, true},
		{"an unknown start time, written before", session, map[string]any{"pid": self, "procStart": "soon", "startedAt": before}, false},
		{"nothing to tell", session, map[string]any{"pid": self}, true},
		{"the daemon", "daemon.lock", map[string]any{"pid": self, "version": "2.1.308", "procStart": now, "startedAt": after}, true},
		{"the daemon on Windows", "daemon.lock", map[string]any{"pid": self, "version": "2.1.308", "procStartFt": now, "startedAt": after}, true},
		{"a daemon that ended", "daemon.lock", map[string]any{"pid": dead, "procStartFt": deadStart}, false},
		{"pid 1", "sessions/1.json", map[string]any{"pid": 1, "startedAt": 1000}, false},
		{"no record: a key file", "sessions/" + strconv.Itoa(self) + ".0123abcd.key", map[string]any{"pid": self}, false},
		{"no record: another name", "sessions/notes.json", map[string]any{"pid": self}, false},
		{"no record: no pid", session, map[string]any{"procStart": "1"}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, c.file)
			os.MkdirAll(filepath.Dir(path), 0o700)
			if err := os.WriteFile(path, []byte(record(c.fields)), 0o644); err != nil {
				t.Fatal(err)
			}
			if got := ClaudeRuns(dir); got != c.want {
				t.Fatalf("ClaudeRuns = %v, want %v for %s", got, c.want, record(c.fields))
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatal("Claude Code's record was removed")
			}
		})
	}
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sessions"), 0o700)
	os.WriteFile(filepath.Join(dir, session), []byte("{"), 0o644)
	if ClaudeRuns(dir) || ClaudeRuns(t.TempDir()) {
		t.Fatal("a record that is no JSON, or none")
	}
}

// A Claude Code this ccam started itself (the usage probe of usage_source
// = "claude") is not someone else's session.
func TestClaudeRunsLeavesOutItsOwnChildren(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		t.Skip("no process tree here")
	}
	child, start := sleeper(t)
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sessions"), 0o700)
	os.WriteFile(filepath.Join(dir, "sessions", strconv.Itoa(child.Process.Pid)+".json"),
		[]byte(record(map[string]any{"pid": child.Process.Pid, "procStart": start, "startedAt": time.Now().UnixMilli()})), 0o644)
	if ClaudeRuns(dir) {
		t.Fatal("this process's own child counted")
	}
}

func TestInUseAsksMarkersAndRecords(t *testing.T) {
	if _, ok := StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	markers, dir := filepath.Join(t.TempDir(), "markers"), t.TempDir()
	if InUse(markers, dir) {
		t.Fatal("nothing runs")
	}
	if err := Register(markers); err != nil {
		t.Fatal(err)
	}
	if !InUse(markers, dir) {
		t.Fatal("a session ccam started")
	}
	Unregister(markers)
	os.Mkdir(filepath.Join(dir, "sessions"), 0o700)
	os.WriteFile(filepath.Join(dir, "sessions", strconv.Itoa(os.Getpid())+".json"),
		[]byte(record(map[string]any{"pid": os.Getpid(), "procStart": claudeStart(t, os.Getpid())})), 0o644)
	if !InUse(markers, dir) {
		t.Fatal("a session Claude Code lists")
	}
}
