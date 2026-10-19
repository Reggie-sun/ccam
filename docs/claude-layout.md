# Claude Code config layout (Phase 0 findings)

Checked on 2026-10-18 with Claude Code **2.1.302** on Ubuntu (WSL2). Experiments used
throw-away `CLAUDE_CONFIG_DIR` folders; the real `~/.claude` was only read (key names, never values).

## Files inside a profile folder

| Path (relative to `CLAUDE_CONFIG_DIR`) | Contents | ccam handling | Verified |
| --- | --- | --- | --- |
| `.credentials.json` | `claudeAiOauth` (accessToken, refreshToken, expiresAt, scopes, subscriptionType, rateLimitTier) and `mcpOAuth` (per-MCP-server OAuth tokens) | always local, 0600 | yes |
| `.claude.json` | `oauthAccount.emailAddress`, `mcpServers`, `projects`, caches (incl. `cachedUsageUtilization`) | file local, only `mcpServers` shareable | yes |
| `.claude.json.lock`, `backups/` | Claude Code's own lock and `.claude.json` backups | local | yes |
| `settings.json` | permissions, hooks, model, env, statusLine, enabledPlugins | shareable (symlink) | yes |
| `CLAUDE.md` | global memory | shareable (symlink) | by docs |
| `agents/` `commands/` `skills/` `plugins/` `output-styles/` `rules/` | extensions | shareable as one toggle (symlinks) | yes (`skills`, `plugins`, `output-styles`, `rules` exist) |
| `projects/` | `<slug>/<session-id>.jsonl` transcripts and `<slug>/memory/` (auto-memory) | shareable with history | yes |
| `history.jsonl` | prompt history (interactive only; `-p` does not write it) | shared with history | yes |
| `sessions/` `session-env/` `shell-snapshots/` `state/` `cache/` `downloads/` `security/` `bridge-spawn/` `policy-limits.json*` `remote-settings.json*` `mcp-needs-auth-cache.json` | runtime state and caches | always local | yes |

**Gotcha: default location.** Without `CLAUDE_CONFIG_DIR`, `.claude.json` is `~/.claude.json`
(next to `~/.claude/`, not inside it). With `CLAUDE_CONFIG_DIR=X` it is `X/.claude.json`.
So `ccam import ~/.claude` must copy `~/.claude.json` into the profile folder; registering
`~/.claude` in place would make Claude Code look for `~/.claude/.claude.json` and lose the account.

**Gotcha: MCP OAuth.** OAuth tokens for remote MCP servers live in `.credentials.json` (`mcpOAuth`),
so a shared MCP server that needs OAuth must be authorised once per account.

## Symlinks and plugins

- Claude Code writes through a symlinked `settings.json` (checked with `claude plugin install/disable`); the link stays a link.
- `plugins/installed_plugins.json` stores absolute `installPath`s under the config dir that installed the plugin,
  e.g. `<profile>/plugins/cache/<marketplace>/<plugin>/<version>`. `known_marketplaces.json` stores the marketplace source path.

## macOS credentials

Read from the Claude Code 2.1.302 binary (the same JavaScript runs on every platform):

- Service name: `"Claude Code" + OAUTH_FILE_SUFFIX + "-credentials" + suffix`. `OAUTH_FILE_SUFFIX` is `""` in production.
  `suffix` is `"-" + sha256(configDir)[:8]` (hex, NFC-normalized path) when `CLAUDE_CONFIG_DIR` is set, else `""`.
- `CLAUDE_SECURESTORAGE_CONFIG_DIR` overrides the path used for the suffix; set to `""` it selects the unsuffixed
  default item. ccam sets it to `""` for an imported `~/.claude` on macOS so `ccam main` and plain `claude` share one login.
- Account: `$USER` if it matches `^[a-zA-Z0-9._-]+$`, else `claude-code-user`.
- Read with `security find-generic-password -a <account> -w -s <service>`; when the Keychain rejects writes,
  Claude Code falls back to `<config dir>/.credentials.json`.
- ccam only reads the Keychain (usage lookup, doctor, `next`); `--refresh` is refused for Keychain logins.
  Not yet checked on a real Mac with a real login.

## Login

- `claude auth login [--claudeai|--console] [--email <e>] [--sso]` logs in without opening a session.
- `claude auth status --json` prints `{loggedIn, authMethod, apiProvider, configDirectory, ...}`; exit code 1 when logged out.
- ccam `add`/`login` therefore run `claude auth login`, then read the email from `.claude.json` (`oauthAccount.emailAddress`).

## MCP

- `claude mcp add-json <name> <json> -s user` writes to `.claude.json` → `mcpServers` in the config dir.
- `--mcp-config` is variadic and repeatable: `--mcp-config a.json --mcp-config b.json` and
  `--mcp-config a.json b.json` both load servers from both files.
- Servers from `--mcp-config` are loaded **in addition to** user servers in `.claude.json`.
  On a name clash the `.claude.json` entry wins (only one `demo` appears).
- `--strict-mcp-config` given by the user still works (it only drops non-flag sources).
- **Variadic gotcha:** `claude --mcp-config a.json "prompt"` treats the prompt as another config file
  ("MCP config file not found"). ccam therefore passes `--mcp-config=<path>` as one argument, which
  stops there; a following prompt, subcommand (`mcp list`) or the user's own `--mcp-config` all work.

**Decision: `mcp_strategy = "flag"` by default.** It never touches `.claude.json`, and stacking with
user-given `--mcp-config` works. `sync` stays available as an opt-in.

## Conversation sharing

Profile B with `projects/` symlinked to profile A's: a conversation made in A was resumed in B
with `claude -p --continue` (B recalled a word only said in A). Sharing `projects/` also shares
per-project auto-memory (`projects/<slug>/memory/`).

## Usage

### Live endpoint

`GET https://api.anthropic.com/api/oauth/usage` with
`Authorization: Bearer <claudeAiOauth.accessToken>` and `anthropic-beta: oauth-2025-04-20` → 200.

Top-level non-null keys seen: `five_hour`, `limits`, `extra_usage`, `spend`, `member_dashboard_available`.
`seven_day` and many model-specific windows may be `null`.

```json
{
  "five_hour": {"utilization": 0.0, "resets_at": "2026-10-18T18:10:00.299994+00:00"},
  "seven_day": null,
  "limits": [
    {"kind": "session", "group": "session", "percent": 0, "severity": "normal",
     "resets_at": "2026-10-18T18:10:00.299994+00:00", "scope": null, "is_active": true},
    {"kind": "weekly_scoped", "group": "weekly", "percent": 0, "severity": "normal",
     "resets_at": "2026-10-24T18:00:00+00:00",
     "scope": {"model": {"id": null, "display_name": "Fable"}, "surface": null}, "is_active": false}
  ]
}
```

- New format: `limits[]` with `kind` in `session` / `weekly_all` / `weekly_scoped`, `percent` 0-100, ISO `resets_at`.
- Old format: `five_hour` / `seven_day` objects with `utilization` (0-100) and ISO `resets_at`.
- Claude Code caches the same payload in `.claude.json` → `cachedUsageUtilization.utilization`.

### Status line

Claude Code pipes JSON to the `statusLine` command. `rate_limits` appears only for Pro/Max
subscribers and only after the first API response:
`rate_limits.five_hour.used_percentage` (0-100), `rate_limits.five_hour.resets_at` (Unix seconds),
same for `seven_day`. Each window may be absent. The command inherits Claude Code's environment,
so `CLAUDE_CONFIG_DIR` (and ccam's `CCAM_PROFILE`) identifies the profile.
Verified end to end: `ccam statusline` inside a real session received `rate_limits.five_hour` after the
first reply (`seven_day` was absent on that account, matching the live endpoint's null `seven_day`).

### Token refresh (unverified)

`ccam usage --refresh` posts `grant_type=refresh_token` to `https://console.anthropic.com/v1/oauth/token`
with Claude Code's public client id. This was **not** exercised against the live service: refresh tokens
rotate, so a test refresh would log out the running session that owns the token. Only unit-tested
against a fake server; it refuses to run while a session of that profile holds the profile lock.
