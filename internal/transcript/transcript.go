// Package transcript finds Claude Code's saved conversations in a config
// folder: <config>/projects/<slug of the working folder>/<session id>.jsonl.
// The entries inside are internal to Claude Code and change between
// versions, so only the few fields used here are read, and anything else
// is skipped.
package transcript

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

// maxSlug is where Claude Code (2.1.304) cuts long folder names.
const maxSlug = 200

// Slug is the folder name Claude Code uses under projects/ for the working
// folder dir (with symbolic links resolved): every UTF-16 unit that is not
// an ASCII letter or digit becomes '-'; a name over 200 characters keeps
// the first 200 and a hash of the whole path.
func Slug(dir string) string {
	units := utf16.Encode([]rune(dir))
	var b strings.Builder
	for _, u := range units {
		if u < 0x80 && (u >= 'a' && u <= 'z' || u >= 'A' && u <= 'Z' || u >= '0' && u <= '9') {
			b.WriteByte(byte(u))
		} else {
			b.WriteByte('-')
		}
	}
	s := b.String()
	if len(s) <= maxSlug {
		return s
	}
	// Java's String.hashCode over the UTF-16 units, as a 32-bit integer.
	var h int32
	for _, u := range units {
		h = 31*h + int32(u)
	}
	v := int64(h)
	if v < 0 {
		v = -v
	}
	return s[:maxSlug] + "-" + strconv.FormatInt(v, 36)
}

// ProjectDir is where config folder configDir keeps the conversations of
// the working folder cwd.
func ProjectDir(configDir, cwd string) string {
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = real
	}
	return filepath.Join(configDir, "projects", Slug(cwd))
}

// Newest is the most recently written conversation in projectDir.
func Newest(projectDir string) (path string, mod time.Time, ok bool) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return "", time.Time{}, false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(mod) {
			path, mod, ok = filepath.Join(projectDir, e.Name()), info.ModTime(), true
		}
	}
	return path, mod, ok
}

// ID is the session id of a conversation file.
func ID(path string) string { return strings.TrimSuffix(filepath.Base(path), ".jsonl") }

// tailSize is how much of the end of a conversation Size reads.
const tailSize = 1 << 20

// Size is about how many tokens the next reply in the conversation reads:
// the prompt of its last reply (input, cache read and cache creation) plus
// that reply. 0 when unknown.
func Size(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	if info, err := f.Stat(); err == nil && info.Size() > tailSize {
		if _, err := f.Seek(info.Size()-tailSize, io.SeekStart); err != nil {
			return 0
		}
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return 0
	}
	lines := bytes.Split(data, []byte("\n"))
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"usage"`)) || !bytes.Contains(l, []byte(`"assistant"`)) {
			continue
		}
		var e struct {
			Type      string `json:"type"`
			Sidechain bool   `json:"isSidechain"`
			Message   struct {
				Usage *struct {
					Input       int `json:"input_tokens"`
					CacheRead   int `json:"cache_read_input_tokens"`
					CacheCreate int `json:"cache_creation_input_tokens"`
					Output      int `json:"output_tokens"`
				} `json:"usage"`
			} `json:"message"`
		}
		if json.Unmarshal(l, &e) != nil || e.Type != "assistant" || e.Sidechain || e.Message.Usage == nil {
			continue
		}
		u := e.Message.Usage
		return u.Input + u.CacheRead + u.CacheCreate + u.Output
	}
	return 0
}

// Title is the conversation's name: the one given with /rename, else the
// one Claude Code made up; "" when it has neither.
func Title(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	custom, ai := "", ""
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64<<10), 16<<20)
	for sc.Scan() {
		l := sc.Bytes()
		if !bytes.Contains(l, []byte(`-title"`)) {
			continue
		}
		var e struct {
			Type   string `json:"type"`
			Custom string `json:"customTitle"`
			AI     string `json:"aiTitle"`
		}
		if json.Unmarshal(l, &e) != nil {
			continue
		}
		switch e.Type {
		case "custom-title":
			custom = e.Custom
		case "ai-title":
			ai = e.AI
		}
	}
	if custom != "" {
		return custom
	}
	return ai
}
