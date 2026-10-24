# Claude Code config layout (Phase 0 findings)

Checked on 2026-10-18 with Claude Code **2.1.302** on Ubuntu (WSL2); the list of top-level entries again with
**2.1.309** on 2026-10-23 (`doctor.LayoutVersion`). Experiments used throw-away `CLAUDE_CONFIG_DIR` folders; the real
`~/.claude` was only read (key names, never values).

## Files inside a profile folder

| Path (relative to `CLAUDE_CONFIG_DIR`) | Contents | ccam handling | Verified |
| --- | --- | --- | --- |
| `.credentials.json` | `claudeAiOauth` (accessToken, refreshToken, expiresAt, scopes, subscriptionType, rateLimitTier) and `mcpOAuth` (per-MCP-server OAuth tokens) | always local, 0600 | yes |
| `.claude.json` | `oauthAccount.emailAddress`, `mcpServers`, `projects`, caches (incl. `cachedUsageUtilization`) | file local, only `mcpServers` shareable | yes |
| `.claude.json.lock`, `backups/` | Claude Code's own lock and `.claude.json` backups | local | yes |
| `settings.json` | permissions, hooks, model, env, statusLine, enabledPlugins | shareable (symlink) | yes |
| `CLAUDE.md` | global memory | shareable (symlink) | by docs |
| `agents/` `commands/` `skills/` `plugins/` `output-styles/` `rules/` | extensions | shareable as one toggle (symlinks) | yes (`skills`, `plugins`, `output-styles`, `rules` exist) |
| `skills/synced/` | skills synced from claude.ai: `<org>_<account>/` per account, each with `manifest.json`, `.last-complete-round` and the skill folders; `.bucket-*` markers beside them | shared with extensions (inside `skills/`); per-account folder names, so accounts do not overwrite each other | yes (2026-10-21, ccam ext skips it) |
| `projects/` | `<slug>/<session-id>.jsonl` transcripts and `<slug>/memory/` (auto-memory) | shareable with history | yes |
| `history.jsonl` | prompt history (interactive only; `-p` does not write it) | shared with history | yes |
| `sessions/` `session-env/` `shell-snapshots/` `state/` `cache/` `downloads/` `security/` `bridge-spawn/` `policy-limits.json*` `remote-settings.json*` `mcp-needs-auth-cache.json` | runtime state and caches | always local | yes |
| `stats-cache.json`, `logs/`, `teams/`, `usage-data/`, `.claude.json.backup`, `active-time.json`, `server-sessions.json`, `mcp-discovery-cache`, `gh-pr-status-cache.json`, `file-transfers`, `mcp-skill-archives`, `chrome`, `remote-control`, `scratch` and the other names below | per machine state Claude Code 2.1.309 lists itself | always local (`share.known`) | yes (binary, 2026-10-23) |

**Gotcha: default location.** Without `CLAUDE_CONFIG_DIR`, `.claude.json` is `~/.claude.json`
(next to `~/.claude/`, not inside it). With `CLAUDE_CONFIG_DIR=X` it is `X/.claude.json`.
So `ccam import ~/.claude` must copy `~/.claude.json` into the profile folder; registering
`~/.claude` in place would make Claude Code look for `~/.claude/.claude.json` and lose the account.

**Top-level entries (2.1.309, 2026-10-23).** Claude Code keeps its own list of what a config folder holds per
machine (the set it leaves out of a self-hosted runner's host config snapshot, next to `HostConfigOverCapError` in the
2.1.309 binary): `.claude.json`, `.claude.json.backup`, `.credentials.json`, `projects`, `sessions`, `todos`,
`shell-snapshots`, `statsig`, `file-history`, `history.jsonl`, `ide`, `logs`, `backups`, `.session_ingress_token`,
`policy-limits.json` and `remote-settings.json` with their `.signature.json`, `.signature-iat.json` (`${name}.signature-iat.json`),
`.stamp.json` and `.synced-at` companions, `remote-settings-helper-consent`, `remote-settings-consent.json`,
`hfi-auth.json`, `ccr-home-seed.json`, `daemon`, `jobs`, `teams`, `usage-data`, `shares`, `state`, `uploads`, `feedback`,
`feedback-bundles`, `plans`, `telemetry`, `dump-prompts`, `debug`, `traces`, `startup-perf`, `cache`,
`mcp-discovery-cache`, `mcp-needs-auth-cache.json`, `gh-pr-status-cache.json`, `tasks`, `local`, `antproto.json`, `ccr`,
`session-env`, `bridge-spawn`, `active-time.json`, `loop.md`, `server-sessions.json`, `image-cache`, `paste-cache`,
`file-transfers`, `mcp-skill-archives`, `stats-cache.json`, `computer-use.lock`, `server.lock`, `api-dumps`, `chrome`,
`downloads`, `local-settings`, `project-settings`, `remote`, `remote-control`, `scratch`, `seed-admin`, `storage-v2`,
`systemd`, `.cc-writes`. All are in `share.known` (known to ccam; never shared, except `tasks` and `paste-cache`, which
the history item shares with `projects/`), and any `policy-limits.json.*` or
`remote-settings.json.*` name counts as known, since Claude Code names those companions after the file. The five
real profiles on this machine and `~/.claude` (names only) held nothing else that ccam did not know.

**Gotcha: MCP OAuth.** OAuth tokens for remote MCP servers live in `.credentials.json` (`mcpOAuth`),
so a shared MCP server that needs OAuth must be authorised once per account.

## Symlinks and plugins

- Claude Code writes through a symlinked `settings.json` (checked with `claude plugin install/disable`); the link stays a link.
- `plugins/installed_plugins.json` stores absolute `installPath`s under the config dir that installed the plugin,
  e.g. `<profile>/plugins/cache/<marketplace>/<plugin>/<version>`. `known_marketplaces.json` stores the marketplace source path.
- Loading such a path from another config dir (2.1.309, `claude mcp list` in a scratch `HOME`/`CLAUDE_CONFIG_DIR`, a
  plugin with a stdio MCP server that logs `CLAUDE_PLUGIN_ROOT`, 2026-10-23): when the recorded `.../plugins` folder does
  not exist, Claude Code reads the plugin from its own `plugins/` folder instead; when that folder exists but lacks the
  plugin, the plugin (and its MCP servers) is not loaded at all. So a group's list recording paths through an account that
  then has a `plugins/` folder of its own (it left the group) loses those plugins for the other accounts: a change of
  groups points such paths at the group folder before any account takes its extensions out of the group.

## macOS credentials

Read from the Claude Code 2.1.302 binary (the same JavaScript runs on every platform), checked again in 2.1.308
(2026-10-22); the macOS writer was read from the darwin-arm64 2.1.306 build by the review that found the rename bug:

- Service name: `"Claude Code" + OAUTH_FILE_SUFFIX + "-credentials" + suffix`. `OAUTH_FILE_SUFFIX` is `""` in production.
  `suffix` is `"-" + sha256(path)[:8]` (hex) of the config dir as given in `CLAUDE_CONFIG_DIR`, in Unicode NFC (macOS may
  hand out a folder's name decomposed), when that variable is set, else `""`. `creds.KeychainService` normalizes too.
- `CLAUDE_SECURESTORAGE_CONFIG_DIR` overrides the path used for the suffix; set to `""` it selects the unsuffixed
  default item. ccam sets it to `""` for an imported `~/.claude` on macOS (and Windows) so `ccam main` and plain `claude`
  share one login. It is more than a name: the folder it names (`~/.claude` when empty) is where Claude Code keeps the
  plaintext fallback `.credentials.json` and the `.storage-write`/`.oauth_refresh` locks, and Claude Code creates it.
- Account: `$USER` when not empty, else Bun's `os.userInfo().username`, which is `$USER` again (`""` when it is empty)
  or `unknown` when it is unset (checked with Bun 1.4.2: not the name in the user database); used when it matches
  `^[a-zA-Z0-9._-]+$`, else `claude-code-user` (so also for an empty `$USER`).
- Read with `security find-generic-password -a <account> -w -s <service>`; written with `security -i` (`add-generic-password
  -U -a … -s … -X <hex>`) when that line is at most 4032 characters, else with those arguments; exit code 44 means no
  item. When the Keychain rejects writes, Claude Code falls back to the plaintext file. MCP sign-ins (`mcpOAuth`) live in
  the same item as the login.
- Because the item is named after the folder's path, moving the folder leaves the login behind (and a later account
  whose folder gets that path would inherit it). ccam therefore never moves a macOS account's folder: `ccam rename`
  keeps it (`folder` in config.toml), `ccam remove` first runs `claude auth logout` there, which deletes the item, and
  `ccam add` stops when an item for the new folder's path still holds a login or MCP sign-ins (left by an older
  ccam; Claude Code's logout writes back an item with only the device key `coworkRemoteDevice`, which does not count),
  naming the `security delete-generic-password` command. ccam itself still only reads the Keychain; `--refresh` is refused for
  Keychain logins. Not yet checked on a real Mac with a real login.

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
- `claude mcp add-json <name> <json>` without `-s user` writes the local scope (below).
- **Variadic gotcha:** `claude --mcp-config a.json "prompt"` treats the prompt as another config file
  ("MCP config file not found"). ccam therefore passes `--mcp-config=<path>` as one argument, which
  stops there; a following prompt, subcommand (`mcp list`) or the user's own `--mcp-config` all work.

**Decision: `mcp_strategy = "flag"` by default.** It never touches `.claude.json`, and stacking with
user-given `--mcp-config` works. `sync` stays available as an opt-in.

### Scopes, subcommands and sign-ins (2.1.308, 2026-10-23)

Checked with the real binary in a scratch `HOME`/`CLAUDE_CONFIG_DIR` (`mcp add/add-json/get/list/remove` and
`--help` only), and by reading the binary for the argument dispatch.

- `claude mcp add` without `-s` uses the **local** scope: `projects[<key>].mcpServers` of the account's
  `.claude.json`, loaded for that account only. The key is the folder with symbolic links resolved; inside a
  git repository it is the top of the **main** worktree (`claude mcp add` in `repo/sub` or in a worktree of
  `repo` files under `repo`, and every folder of the repository and its worktrees loads it); for a worktree
  of a bare repository it is the bare repository's folder (`proj.git`). Outside git a folder is its own key:
  `plain/sub` does not load `plain`'s servers. The message `[project: <cwd>]` names the current folder, not
  the key. ccam: `mcp.ProjectKey`, `mcp.FolderServers`.
- One name in both scopes: the local one wins (`mcp get` shows `Scope: Local config`).
- `claude mcp list`, `get` and `remove` read `.claude.json`, the project's `.mcp.json` and the plugins only:
  with `--mcp-config=<group>/mcp.json` in front, `list` says "No MCP servers configured" and `get`/`remove`
  "No MCP server named …" (exit 1); after the verb, `--mcp-config` is an unknown option. ccam leaves the flag
  out there and prints the group's servers on stderr.
- A project's `.mcp.json` servers need the folder's **trust**, which is per account
  (`projects[<folder>].hasTrustDialogAccepted`): trust of the folder itself or, inside git, of a folder up to
  the repository's top counts, and so does trust of the local-scope key (the main worktree's folder, or the
  bare repository's) in every worktree; a folder above those does not, nor a parent outside git. Claude Code
  moves an approval in `projects[…].enabledMcpjsonServers` into `<project>/.claude/settings.local.json`,
  which then approves the servers for every account that trusts the folder. `enabledMcpjsonServers` or `enableAllProjectMcpServers`
  in the user `settings.json` (the group's, when shared) approve them without trust; the same keys in the
  project's `.claude/settings.json` do not; `disabledMcpjsonServers` in the user settings hides them.
  `claude mcp list` shows the others as "Pending approval"; an interactive session asks once, `-p` takes them
  (`by: "non_interactive"`, read in 2.1.306).
- Sign-ins of remote servers (`mcpOAuth` in the login) are keyed `<name>|` + 16 hex digits of a SHA-256 of
  `{type, url, headers || {}}` (function `Gl` in 2.1.308). A health check (`mcp list`) writes an entry with
  an empty `accessToken` and no `refreshToken` for a server that waits for a sign-in; only an entry with a
  token is a sign-in. `claudeAiMcpEverConnected` in `.claude.json` lists the claude.ai connectors a session of
  the account connected; connectors come with the claude.ai account.
- Plugins' MCP servers load only for plugins `enabledPlugins` names `true` (user settings, then the project's
  `.claude/settings.json`, which overrides them either way; a plugin not named is off), from the plugin's
  `.mcp.json` and its `plugin.json` `mcpServers` (an object, a path or a list; each file flat or wrapped in
  `mcpServers`), all of them, named `plugin:<plugin>:<server>`. Plugin servers with the same command and
  arguments are listed once. Checked with `claude mcp list` and a directory marketplace in a scratch folder
  (`settings.local.json` overriding the project's was checked with 2.1.306 by the review).
- Argument dispatch: `remote-control`, `rc`, `remote`, `sync` and `bridge` as the **first** argument start
  Remote Control before any option is read; after root options only `remote-control`/`rc` (commander), and
  Remote Control refuses every root option but its allow-list (`--verbose`, `--debug`, `--debug-to-stderr`,
  `--debug-file`, `--ax-screen-reader`, `--workload`, `--session-id`, `--name`,
  `--remote-control-session-name-prefix`, `--plugin-dir`, `--plugin-dir-no-mcp`, `--plugin-url`, `--add-dir`,
  `--ide`, `--chrome`, `--model`, `--effort`, `--fallback-model`, `--betas`, `--thinking`, `--thinking-display`,
  `--max-thinking-tokens`, `--autocompact`, `--allowedTools`, `--dangerously-skip-permissions`,
  `--allow-dangerously-skip-permissions`, `--enable-auto-mode`, `--bare`): `--mcp-config` is refused. With
  `-p` before `--` there are no subcommands (read in 2.1.306). The background-session verbs `attach`, `logs`, `stop`, `kill`,
  `respawn`, `rm` count by name only as the first argument; behind root options only with an 8-hex id (or
  `respawn --all`), else the words become a prompt. ccam (`launcher.claudeVerb`) therefore passes no
  `--mcp-config` before Remote Control, a background-session verb given first, or `mcp list|get|remove`;
  every other subcommand takes it without harm.

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

On Windows, 2.1.309's PowerShell lookups (read in the binary) prefer pwsh.exe
(`C:\Program Files\PowerShell\7`, `…\7-preview`, `%LOCALAPPDATA%\Microsoft\WindowsApps`,
`~\.dotnet\tools`) over `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`; which lookup
the status line goes through was not traced. So a status line command may meet PowerShell 7 or
Windows PowerShell 5.1, which pass quoted arguments to programs differently; ccam's own text in it
needs no quotes on Windows (`--wrap-base64`).

### npm install on Windows

`npm i -g @anthropic-ai/claude-code` puts `claude.cmd` in `%APPDATA%\npm` (a cmd-shim:
`"%_prog%" "%dp0%\node_modules\@anthropic-ai\claude-code\cli.js" %*`, or the older
`@IF EXIST "%~dp0\node.exe"` form). Windows runs a batch file through cmd.exe, which parses the
command line itself; node.exe then reads `%*` by the C runtime's rules (`""` inside quotes is one
quote). Checked on the Windows 11 host with node.exe and both shim forms (2026-10-23).

### Token refresh (unverified)

`ccam usage --refresh` posts `grant_type=refresh_token` to `https://platform.claude.com/v1/oauth/token`
with Claude Code's public client id. Claude Code 2.1.304 uses this address (`TOKEN_URL` in its bundle; earlier
versions used `https://console.anthropic.com/v1/oauth/token`, which still answers). A probe with an invalid refresh
token on 2026-10-20 got a JSON `rate_limit_error` (HTTP 429) from both, so both endpoints are live; ccam reports HTTP 429 as
"try again later" rather than asking for a new login. This was **not** exercised against the live service: refresh tokens
rotate, so a test refresh would log out the running session that owns the token. Only unit-tested
against a fake server; it refuses to run while Claude Code runs on that profile (see "Running processes") and while
Claude Code's refresh lock is held, and follows Claude Code's own rules for the request and the save (see "Logins").

## Logins (2.1.308)

Read from the 2.1.308 binary (Windows parts from the 2.1.304 `claude.exe` on this machine's Windows host), 2026-10-22;
nothing was run against a real login.

- Renewal: `POST TOKEN_URL {grant_type: "refresh_token", refresh_token, client_id (the stored clientId, else
  CLIENT_ID), scope (the stored scopes)}` with `timeout: 30000`. An answer without `refresh_token` keeps the old one;
  `refresh_token_expires_in` is stored as `refreshTokenExpiresAt` only when present. A profile lookup (10 s) may follow
  before the save. On exit Claude Code waits 2 s for a renewal in flight (`OAuth refresh still in flight after 2000ms`).
- Locks, all proper-lockfile mkdir locks in the secure-storage folder (`CLAUDE_SECURESTORAGE_CONFIG_DIR`, else the
  config dir): a renewal holds `.oauth_refresh.lock` (stale after 60 s, touched every 5 s) and the older
  `<realpath(dir)>.lock`, tried 5 times 1-2 s apart; its holder writes `.oauth_refresh.lock.owner` (pid, start time,
  lock birth time), and another process takes the lock over early only when that record proves the holder gone (a
  lock without a record waits for the 60 s). Every read-modify-write of the login (renewals, MCP sign-ins, the
  dead-token clear) holds `.storage-write.lock` (stale after 15 s, 10 retries 100-1000 ms apart) and re-reads under it.
  A renewal is saved only when the stored `refreshToken` is still the one posted or `""`; otherwise the newer login stays
  (`adopted_newer_write`). ccam's `--refresh` takes the first two locks without waiting and saves under the third
  (waiting up to 5 s, then writing anyway, since the server has rotated the token by then), with the same save rule,
  changing only `accessToken`, `refreshToken`, `expiresAt` (and `refreshTokenExpiresAt` when the answer gives a
  lifetime) of the file as it is then; it sends the login's `clientId` but no `scope`.
- A dead login: when the server refuses the refresh token (`invalid_grant`), Claude Code keeps `.credentials.json` with
  `refreshToken: ""`, `accessToken: ""`, `expiresAt: 0` (other keys, `refreshTokenExpiresAt` included, stay). It treats
  `refreshToken === ""` as dead (`/status`: "Login: Expired — log in again"). `claude auth status --json` says
  `loggedIn: true` whenever `accessToken` is not empty, so it is no test of a working login. `creds.Judge` follows the
  first rule: no `claudeAiOauth` is not logged in; an empty refresh token or a past `refreshTokenExpiresAt` is expired.
- `claude auth logout` revokes the refresh tokens of `claudeAiOauth`, `designOauth` and `cloudSessionsOauth` (`POST
  TOKEN_URL/revoke`, 5 s each, failures ignored: "continuing with local logout"), then deletes the whole store, the MCP
  sign-ins with it, writing back only `coworkRemoteDevice`. `claude auth login` deletes only the Anthropic keys
  (`claudeAiOauth`, `organizationUuid`, trusted-device keys, `designOauth`, `cloudSessionsOauth`) and keeps `mcpOAuth`,
  and does not revoke. So a login as another account needs no logout, and a logout of a copied folder would revoke
  the other copy's login too.
- MCP sign-ins (`mcpOAuth`, keyed `<server>|<hash of type, url, headers>`): an entry with an empty `accessToken` and no
  `refreshToken` is no sign-in (a server that still needs one, or one whose tokens were dropped); Claude Code registers as
  a public client (`token_endpoint_auth_method: "none"`), whose refresh tokens rotate, and clears an entry whose refresh
  is refused. A sign-in therefore works for one holder only, which is why `ccam migrate` no longer copies them.
- Windows Credential Manager (2.1.304 Windows build; the Linux build always uses the file): used when
  `CLAUDE_CODE_FORCE_WINDOWS_CREDMAN=1`, or the cached feature flag `cachedGrowthBookFeatures.tengu_windows_credman` is
  true in the config folder's `.config.json` (when not empty), else its `.claude.json`. Entries go through
  `Bun.secrets` with the Keychain's service name and the account `claude-code-user`; values over 2400 characters are
  split (`<name>#m` holds `{n, l}`, `<name>#0`…, `<name>#p`, base64). The store reads Credential Manager first and the
  file after; its first successful write deletes `.credentials.json`. ccam does not read Credential Manager (the
  target name Bun builds is not verified): with the flag on and no file it treats the login as unknown, not logged out
  (`creds.ErrCredentialManager`), keeps the folder on `ccam rename` as on macOS, and passes
  `CLAUDE_SECURESTORAGE_CONFIG_DIR=""` for an imported `~/.claude` there too.

## Usage limits in a session (2.1.306)

Read from the Claude Code 2.1.306 binary and from the main-thread limit rows of real transcripts on this machine
(key names and rate-limit values only), 2026-10-22.

### The limit row

When the server refuses a request at a usage limit, Claude Code appends an assistant row to the conversation
(2.1.304 and 2.1.306 wrote the same keys; seven rows seen, all `five_hour` on a Team seat):

- `type: "assistant"`, `isSidechain` (true for a subagent's row), `isApiErrorMessage: true`,
  `apiError: "usage_limit_reached"`, `error: "rate_limit"`, `apiErrorStatus: 429`, `timestamp`, `perTurnEffort`
  (`xhigh`, `max`), `entrypoint`, `version`.
- `message.model: "<synthetic>"` with all-zero `usage`, and the text Claude Code shows
  (`You've hit your session limit · resets 3:30am (Asia/Seoul)`).
- `apiErrorParams.rate_limit_info` is an object: `status: "rejected"`, `resetsAt` (Unix seconds, an integer),
  `rateLimitType` (`five_hour`, `seven_day`, `seven_day_opus`, `seven_day_sonnet`, `seven_day_overage_included`,
  `overage`), `overageStatus`, `overageDisabledReason` (`org_level_disabled`, `out_of_credits` seen), `isUsingOverage`.
  `quotaLimits` repeats it with `unifiedRateLimitFallbackAvailable`.
- A model that needs usage credits (the Fable allowance on a Fable model) is another apiError,
  `model_requires_usage_credits`, without `rate_limit_info`.

ccam reads only this row (`transcript.LastLimit`): `hook stop-failure` records the window as full until `resetsAt`
(`state/usage/<name>.limit.json`) for `five_hour`, `seven_day`, `seven_day_opus` and `seven_day_sonnet`.
The row's `<synthetic>` usage is not a reply, so `transcript.Size` and `transcript.Model` skip it.

### The StopFailure hook

- Input: `session_id`, `transcript_path`, `cwd`, `permission_mode` (not set for StopFailure), `agent_id` and
  `agent_type` (set for a subagent's loop, which fires the hook too), `effort: {level}` (only for models with an
  effort level), `error`, `error_details`, `last_assistant_message` (the stopping message's text blocks joined with
  newlines and trimmed: for a limit, the same message as the row). The hook is not waited for, and the row may reach
  the file after it starts (writes are flushed every 100 ms), so ccam reads again for up to 1.5 s, and only takes a
  row with that text which no prompt of the person's came after (in a shared history an older row can be another
  account's). When `last_assistant_message` is empty, the last row that ends the last turn is taken.
- Environment: `CLAUDE_PID`, `CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_CHILD_SESSION=1`, `CLAUDE_EFFORT`, and
  `CLAUDE_CODE_SESSION_ATTENDED`, which is `"0"` for `CLAUDE_CODE_SESSION_KIND` bg, daemon and daemon-worker
  sessions, for teammates, and for a `-p` session (a non-interactive session counts as attended only when it is not a
  child session and its entrypoint is a desktop or remote one). `CLAUDE_CODE_ENTRYPOINT` is `cli` in an interactive
  session and `sdk-cli` in `-p`. ccam's hook does nothing when `agent_id` is set, when the session is unattended,
  or when the entrypoint starts with `sdk-`.

### Limit messages

`You've hit your ` + the window's name (`session limit`, `weekly limit`, `Opus limit`, `Sonnet limit` (`weekly
limit` on Pro and Enterprise), `Fable limit`) + ` · resets <time>`, or just `limit` (overage refused) or
`usage limit`. Spend limits (the account's, or its team's, organization's or channel's) use the same start:
`monthly spend limit`, `individual spend limit`, `org's monthly spend limit`, `team's shared budget`, `channel's
monthly spend limit`; when a window is full as well, the message ends `· your session limit resets <time>`.
Organization stops start otherwise (`Your org is out of usage`, `Your group's usage limit is set to $0`, `Your usage
allocation has been disabled by your admin`, `This service is disabled for your org`). Claude Code counts a message as a
usage limit when it starts with `You've hit your`, `You've reached your` (`You've reached your Fable limit.`),
`You're out of usage credits`, `Your org is out of usage · …`, `Your seat type doesn't include …`, `Your usage
allocation has been disabled by your admin`, `Your group's usage limit is set to $0`, `Fable 5 requires usage
credits` or `You're out of extra usage`, or matches `^(?:Fable|Haiku|Mythos|Opus|Sonnet)(?: [^·\n]{1,40})? requires
usage credits\.`; `handoff.Classify` mirrors this. The server's throttle (`Server is temporarily limiting requests
(not your usage limit)`) and an API key's 429 are no usage limits.

Reset times: `3:45pm` (or `9am`) within 24 hours, else `Oct 24, 9am`, with the year when it differs, in the local
time zone; the message adds ` (<zone>)`. `handoff.FormatReset` writes the same without the zone.

### Usage through Claude Code

The `get_usage` answer has `rate_limits_available: true` and `rate_limits: null` when the login has plan limits but
Claude Code's own lookup did not succeed (its status was not `ok` or `seeded`); `false` means no plan limits apply (API
key, Bedrock, Vertex). The usage endpoint's `extra_usage.is_enabled` says whether usage credits are on
(`Limits.Credits` in ccam).

## Carrying a session on (2.1.308)

Read from `claude --help` of 2.1.308 (run with a scratch `HOME` and `CLAUDE_CONFIG_DIR`) and the 2.1.308 binary,
2026-10-22. ccam's relay hop (`launcher.CarryArgs`, `WithState`) and the history check build on these.

- `--permission-mode` lists `acceptEdits`, `auto`, `bypassPermissions`, `manual`, `dontAsk`, `plan` and also takes
  `default`; any other value stops Claude Code with exit 1 (`option '--permission-mode <mode>' argument 'bogus' is
  invalid`). An unknown `--effort` value only warns (`Unknown --effort value … ignoring it`). So a hop passes a mode
  only when it is one of these.
- How many values a flag takes (commander): `<values...>` flags take the next argument whatever it is, then every
  following one up to the next starting with `-` (`--add-dir`, `--mcp-config`, `--allowedTools`, `--disallowedTools`,
  `--tools`, `--betas`, `--file`), so a prompt right after one is one of its values; `--flag=value` takes no more. `[value]` flags take the
  next argument unless it is a flag (`-r/--resume`, `-d/--debug`, `-w/--worktree`, `--remote-control`, `--from-pr`,
  `--teleport`, `--cloud`); `<value>` flags always take the next one. New since the 2.1.306 list in the plan:
  `--bg`/`--background`, `--cloud`, `--desktop`, `--environment`, `--file`, `--plugin-url`, `--permission-prompts`.
- Subcommands that open no conversation: `agents`, `auth`, `auto-mode`, `doctor`, `gateway`, `import`, `install`,
  `logs`, `mcp`, `plugin`/`plugins`, `purge`, `rm`, `setup-token`, `stop`/`kill`, `ultrareview`, `update`/`upgrade`.
  `attach` and `respawn` run sessions, and `remote-control` (not in the help) serves them. ccam asks nothing about
  a shared history for the first kind, nor for `--version`/`--help`, nor for a `-p` run that resumes nothing.
- Rows that are no prompt or reply: `cost-state` (`sessionId`, `totalCostUSD` …) and `last-prompt` (`lastPrompt`,
  `leafUuid`), written when a session ends, possibly long after its last message; `permission-mode` rows are
  last-wins. A session that only ran local commands (`/mcp`, `/usage`, `/status`) has user rows of `<command-name>`
  and `<local-command-…>` text, or `system`/`local_command` rows, and no reply (seen by the review in real
  conversations). `transcript.Latest` ranks a folder's conversations by their last prompt or reply and skips such
  sessions and empty files.
- When a session's cost is saved into its project's entry of the account's own `.claude.json`, `lastSessionId` is set
  to the current session id (`projects[<folder>].lastSessionId`). With shared history this is how ccam tells which
  account ran a conversation (`ranLast`); several accounts list the same id when the conversation moved on with
  `--resume`.

## Running processes (2.1.308)

Read from the 2.1.308 binary (Windows parts from the 2.1.304 `claude.exe` on this machine's Windows host) and from the
live records on this machine (key names, pids and start times only), 2026-10-22. ccam's `session.ClaudeRuns` reads
these records, `session.InUse` adds ccam's own markers, and every command that moves, rewrites or removes what a
running Claude Code writes (`rename`, `remove`, group changes, `doctor --fix`, `mcp import`, the MCP copy and sync,
`usage --refresh`, `migrate`'s sign-ins) asks them right before it changes anything.

- Every Claude Code session lists itself in `sessions/<pid>.json` of its config folder (`CLAUDE_CONFIG_DIR`, else
  `~/.claude`) when it starts, however it was started, and removes the file when it exits. A killed process leaves it
  behind: `~/.claude/sessions/` here still holds one written by 2.1.269 in September. Registration is skipped only for
  diskless sessions, agent-team teammates and child sessions inside a tmux that marks them (code read), so `-p` runs
  and background workers are listed too.
- Keys seen: `pid`, `sessionId`, `cwd`, `startedAt`, `procStart`, `version`, `kind` (`interactive`, or the
  `CLAUDE_CODE_SESSION_KIND` `bg`, `daemon`, `daemon-worker`), `entrypoint`, `status` (`busy`, `idle`, `shell` seen),
  `pidDomain`, `peerProtocol`, `peerFeatures`, `messagingSocketPath`, `name`, `nameSource`, `nameSince`,
  `updatedAt`, `statusUpdatedAt`, `bridgeSessionId` (Remote Control connected); a background worker adds `jobId`, a
  worker the daemon starts ahead of time `spare`. The `<pid>.<64 hex>.key` files beside them hold peer tokens; ccam
  never opens them.
- `procStart` is how Claude Code tells a live process from a reused pid: on Linux field 22 of `/proc/<pid>/stat`, the
  same value ccam's `session.StartTime` reads (equal for the three live sessions here); on macOS and other Unix
  systems what `ps -o lstart= -p <pid>` prints with `LC_ALL=C TZ=UTC` (`Thu Oct 22 03:28:37 2026`, code read); on
  Windows the creation time as a FILETIME from Claude Code's native helper (`Bun.ant.windowsProcessCreationTime`),
  or, without it, .NET ticks of `Win32_Process.CreationDate` from PowerShell, which Claude Code itself does not compare
  with a FILETIME (code read). Readers take `procStartFt` before `procStart`. The value is left out when Claude Code
  could not read it (its `ps` call has a 1 s timeout); `startedAt` (ms since 1970) is in every record, `daemon.lock`
  too. Claude Code counts no pid 1, and removes only records whose pid is gone (code read).
- ccam counts a record while its pid (above 1) runs and its start time matches (to the second on macOS and Windows);
  a process that has ended but that its parent has not reaped (a zombie, which an orphan stays in a container whose
  PID 1 reaps nothing) does not run. Without a start time it can compare, the process must not have started more than
  a minute after the record's `startedAt`, so a record whose pid was reused does not block. A record of a process ccam
  itself started (the `usage_source = "claude"` probe, `claude -p --no-session-persistence`, which registers too) does
  not count; on Windows the parent comes from a process snapshot, and one that has ended or started after the child
  is none. Records are never removed by ccam.
- The daemon behind background sessions (`--bg`, `claude agents`) starts through a path of its own and is not in
  `sessions/` (code read): it holds `daemon.lock` in the config folder, `{pid, version, jsonPath, logPath, startedAt, origin, spawnedBy, procStart (procStartFt on Windows),
  launchTarget, processWrapper}`, and its workers list themselves as `kind: "bg"`. It renews the account's login
  itself (`auth: proactive refresh succeeded` in `daemon.log`) and lived for hours here (`uptime=87781s`); `claude
  daemon --help`: "the daemon runs on demand and exits when the last client disconnects", `claude daemon stop --any`
  stops a transient one, `claude stop <id>` one background session.
- `claude remote-control` is seen through the sessions it spawns, which register like any session (code read, not
  run); the server itself writes only `projects/<slug>/bridge-pointer.json` (`source`, `pid`, `procStart`, …), which
  ccam does not read: with shared history that folder is the group's, and one account's server would mark them all.
  A Remote Control server ccam started is seen through ccam's marker.
- A storage backend gated off on this machine (storageV5) keeps the session record and the daemon lock elsewhere when
  it is on; ccam would then see only the sessions it started itself.

