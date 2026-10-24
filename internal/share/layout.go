// Package share links profile folders to shared group folders.
//
// Every shareable path is listed in Layout (an allow list): anything not
// listed, .credentials.json above all, can never end up in a group.
package share

import "strings"

// Entry is one path inside a profile folder.
type Entry struct {
	Name string
	Dir  bool
	// Late marks a path added to an item after accounts already shared it:
	// such an account may still have its own copy, which a launch leaves
	// alone and `ccam doctor --fix` (or the next group change) moves in.
	Late bool
}

// Layout maps share items to the profile paths they cover.
// The "mcp" item has no paths; it is handled by the mcp package.
// See docs/claude-layout.md for where this comes from.
var Layout = map[string][]Entry{
	"settings": {{Name: "settings.json"}, {Name: "keybindings.json", Late: true}},
	"memory":   {{Name: "CLAUDE.md"}},
	"extensions": {
		{Name: "agents", Dir: true}, {Name: "commands", Dir: true}, {Name: "skills", Dir: true},
		{Name: "plugins", Dir: true}, {Name: "output-styles", Dir: true}, {Name: "rules", Dir: true},
	},
	// file-history/ holds the backups /rewind restores code from, tasks/
	// the task lists and paste-cache/ the large pastes history.jsonl points
	// to: without them a conversation carried on elsewhere loses those.
	"history": {
		{Name: "projects", Dir: true}, {Name: "history.jsonl"},
		{Name: "file-history", Dir: true, Late: true}, {Name: "tasks", Dir: true, Late: true}, {Name: "paste-cache", Dir: true, Late: true},
	},
	"mcp": nil,
}

// Never lists per-account files that must never be linked or moved.
var Never = map[string]bool{".credentials.json": true, ".claude.json": true}

// BackupDir is where originals go before a link replaces them.
const BackupDir = ".ccam-backup"

// known lists every top-level name seen in a profile folder (Claude Code
// 2.1.302 to 2.1.309, docs/claude-layout.md). doctor reports anything else.
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
	// Claude Code 2.1.304: also per account, never shared.
	"tasks": true, "paste-cache": true, "keybindings.json": true, "uploads": true, "image-cache": true,
	"jobs": true, "daemon": true, "daemon.json": true, "daemon.lock": true, "daemon.log": true, "daemon.status.json": true,
	"feedback": true, "policy-limits.json.signature.json": true, "remote-settings.json.signature.json": true,
	".oauth_refresh.lock": true,
	// Claude Code 2.1.308: the login file's write lock and the refresh
	// lock's owner record (ccam takes both locks around `usage --refresh`).
	".storage-write.lock": true, ".oauth_refresh.lock.owner": true,
	// Claude Code 2.1.309: the rest of its own list of what a config
	// folder holds per machine (kept out of its host config snapshots);
	// all per account, never shared.
	".claude.json.backup": true, ".session_ingress_token": true, "active-time.json": true, "antproto.json": true,
	"api-dumps": true, "ccr": true, "ccr-home-seed.json": true, "chrome": true, "computer-use.lock": true,
	"dump-prompts": true, "feedback-bundles": true, "file-transfers": true, "gh-pr-status-cache.json": true,
	"hfi-auth.json": true, "local": true, "local-settings": true, "logs": true, "loop.md": true,
	"mcp-discovery-cache": true, "mcp-skill-archives": true, "project-settings": true, "remote": true,
	"remote-control": true, "remote-settings-consent.json": true, "remote-settings-helper-consent": true,
	"scratch": true, "seed-admin": true, "server-sessions.json": true, "server.lock": true, "shares": true,
	"startup-perf": true, "stats-cache.json": true, "storage-v2": true, "systemd": true, "teams": true,
	"traces": true, "usage-data": true, ".cc-writes": true,
	// ccam's own: the flock its processes take turns on before
	// .claude.json.lock (lock.DirLock).
	".claude.json.lock.ccam": true,
}

// Known reports whether name is one ccam knows in a profile folder. Claude
// Code names the companions of its policy and remote settings files after
// them (.signature.json, .signature-iat.json, .stamp.json, .synced-at), so
// any such name counts.
func Known(name string) bool {
	return known[name] || strings.HasPrefix(name, "policy-limits.json.") || strings.HasPrefix(name, "remote-settings.json.")
}
