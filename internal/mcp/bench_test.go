package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// bigClaudeJSON is a 1.6 MB .claude.json: 400 projects with old prompt
// histories, one local-scope server each, user servers at the end.
func bigClaudeJSON() []byte {
	var b bytes.Buffer
	b.WriteString(`{"numStartups":3,"projects":{`)
	for i := 0; i < 400; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `"/p%d":{"allowedTools":[],"mcpServers":{"n%d":{"command":"x"}},"history":[`, i, i)
		for j := 0; j < 20; j++ {
			if j > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`"` + strings.Repeat("x", 200) + `"`)
		}
		b.WriteString(`]}`)
	}
	b.WriteString(`},"oauthAccount":{"emailAddress":"a@b"},"mcpServers":{"a":{"command":"b"}}}`)
	return b.Bytes()
}

// typicalClaudeJSON is the size of the owner's (about 80 KB): mostly
// caches, a few projects.
func typicalClaudeJSON() []byte {
	var b bytes.Buffer
	b.WriteString(`{"numStartups":300,"cachedGrowthBookFeatures":{`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `"tengu_feature_%d":{"enabled":true,"value":"%s"}`, i, strings.Repeat("v", 30))
	}
	b.WriteString(`},"projects":{`)
	for i := 0; i < 6; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `"/home/me/p%d":{"allowedTools":[],"mcpServers":{},"hasTrustDialogAccepted":true,"lastSessionId":"%s"}`, i, strings.Repeat("a", 36))
	}
	b.WriteString(`},"oauthAccount":{"emailAddress":"a@b"},"mcpServers":{}}`)
	return b.Bytes()
}

// The launch's hint after Claude Code rewrote .claude.json: one read, both
// scopes scanned.
func BenchmarkNewLocalChanged(b *testing.B) {
	for name, data := range map[string][]byte{"typical": typicalClaudeJSON(), "big": bigClaudeJSON()} {
		b.Run(name, func(b *testing.B) {
			dir, seen := b.TempDir(), filepath.Join(b.TempDir(), "seen")
			os.WriteFile(filepath.Join(dir, ".claude.json"), data, 0o600)
			group := filepath.Join(b.TempDir(), "mcp.json")
			b.SetBytes(int64(len(data)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				os.Remove(seen) // as if .claude.json changed since the last launch
				if _, err := NewLocal(dir, group, seen); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// A launch under the sync strategy when nothing changed: the stamp, the
// group's mcp.json and the check that no copy is gone.
func BenchmarkSyncUnchanged(b *testing.B) {
	prof, group, state := b.TempDir(), b.TempDir(), b.TempDir()
	os.WriteFile(filepath.Join(prof, ".claude.json"), bigClaudeJSON(), 0o600)
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"g"}`)}}).Save(Path(group))
	later := time.Now().Add(-time.Minute)
	os.Chtimes(Path(group), later, later)
	if _, err := Sync(prof, group, state); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Sync(prof, group, state); err != nil {
			b.Fatal(err)
		}
	}
}
