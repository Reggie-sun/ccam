// Package share links profile folders to shared group folders.
//
// Every shareable path is listed in Layout (an allow list): anything not
// listed, .credentials.json above all, can never end up in a group.
package share

// Entry is one path inside a profile folder.
type Entry struct {
	Name string
	Dir  bool
}

// Layout maps share items to the profile paths they cover.
// The "mcp" item has no paths; it is handled by the mcp package.
// See docs/claude-layout.md for where this comes from.
var Layout = map[string][]Entry{
	"settings": {{"settings.json", false}},
	"memory":   {{"CLAUDE.md", false}},
	"extensions": {
		{"agents", true}, {"commands", true}, {"skills", true},
		{"plugins", true}, {"output-styles", true}, {"rules", true},
	},
	"history": {{"projects", true}, {"history.jsonl", false}},
	"mcp":     nil,
}

// Never lists per-account files that must never be linked or moved.
var Never = map[string]bool{".credentials.json": true, ".claude.json": true}

// BackupDir is where originals go before a link replaces them.
const BackupDir = ".ccam-backup"

// known lists every top-level name seen in a profile folder (Claude Code
// 2.1.302, docs/claude-layout.md). doctor reports anything else.
var known = map[string]bool{
	".credentials.json": true, ".claude.json": true, ".claude.json.lock": true,
	".last-cleanup": true, ".last-update-result.json": true, BackupDir: true,
	"backups": true, "bridge-spawn": true, "cache": true, "downloads": true,
	"history.jsonl": true, "mcp-needs-auth-cache.json": true, "output-styles": true,
	"plugins": true, "policy-limits.json": true, "policy-limits.json.stamp.json": true,
	"projects": true, "remote-settings.json": true, "remote-settings.json.synced-at": true,
	"rules": true, "security": true, "session-env": true, "sessions": true,
	"settings.json": true, "shell-snapshots": true, "skills": true, "state": true,
	"agents": true, "commands": true, "CLAUDE.md": true, "todos": true, "statsig": true,
	"ide": true, "plans": true, "file-history": true, "debug": true, "telemetry": true,
	"settings.local.json": true,
}

func Known(name string) bool { return known[name] }
