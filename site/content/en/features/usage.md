---
title: Usage limits
description: ccam usage shows the 5-hour, weekly and per-model (Fable) limits of every account at once as bars, and ccam statusline puts the account name and usage in Claude Code.
updated: 2026-10-21
order: 7
translation: source
---

`ccam usage` shows how much of the 5-hour limit, the weekly limit and each model's own weekly limit (such as Fable) every account has used, as a bar per limit with the time until it resets. Use it before you start a long task, or when Claude Code tells you a limit is near, to see which account has room left. `ccam statusline` is the status line command for Claude Code: it shows the account name and usage inside Claude Code and saves the numbers for `ccam usage`.

## Usage

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Each account gets a line with its number, its name and how old its numbers are; `(cached)` means they come from a saved file, not from the endpoint. Under it, one line per limit shows a 20-cell bar, the percentage and the time until the window resets: `5-hour`, `Weekly` (all models), then each model's own weekly window, such as `Fable`. A window the account does not have shows `-`. A percentage turns yellow from 80% and red at 100%. On a terminal every account appears at once as `fetching…`, and its lines are drawn as it answers. The last line names the account [`ccam next --auto`](/features/next) would choose by the same rules (`auto_max_usage`, `auto_strategy`, the running sessions, accounts of the last used account's organization only unless `cross_org` is set, every window at 100% skipped, so an account whose week is full is never named), and only appears when you have more than one account. An account that is not logged in shows `not logged in → ccam login <name>`, and an API key account shows `no limits (API)`.

After `ccam statusline install`, Claude Code shows a line such as `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--json` | Print the reports as JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | off |
| `--watch[=N]` | Redraw every `N` seconds until you press Ctrl+C. `--watch` alone, or a value of 0 or less, means 60 | off |
| `-v`, `--verbose` | Add a `source:` line under each account: where its numbers came from (the usage endpoint, Claude Code, the status line or a saved file) | off |
| `--no-cache` | Ignore the 60-second cache and ask the endpoint again | off |
| `--source live\|claude\|cache` | `live` asks the endpoint with the account's token; `claude` has Claude Code look it up (below); `cache` shows saved numbers only | `usage_source`, else `live` |
| `--refresh` | Renew expired access tokens, only for accounts that are not running | off |
| `--dry-run` | Accepted: `ccam usage` changes nothing but its own cache. Refused together with `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

Any other `--source` value fails with `--source must be live, claude or cache`.

### Usage from Claude Code itself

With `usage_source = "claude"` in `config.toml` (or `--source claude`), ccam does not read the login token at all. For each account it starts `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` in the account's folder and sends Claude Code's `get_usage` request on standard input (Claude Code 2.1.270 or later); Claude Code looks the usage up with its own login and answers with the 5-hour, weekly and model windows. `--safe-mode` keeps the account's hooks, MCP servers and plugins from starting, and no session is saved. Each account takes a Claude Code start, about a second, so this source is slower than `live`; the accounts are asked at once, each within 20 seconds, and the 60-second cache applies as usual. Claude Code marks the request as experimental, so if it stops answering, ccam shows the saved numbers and `ccam doctor` reports it. A `--dry-run` does not start Claude Code (it writes to the account's folder) and shows the saved numbers. `usage_source` takes `live` or `claude`; `cache` is a `--source` choice only. `ccam usage` takes no arguments.

`ccam statusline` reads Claude Code's status line JSON from standard input. Its one option is `--wrap <command>`: an existing status line command whose first output line gets the account name in front. `ccam statusline install` takes no options. Both refuse `--dry-run` with `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## How it works

The numbers come from two places: Anthropic's undocumented usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, called with each account's own login, and the values Claude Code hands to its status line command. All of it lives in `internal/usage`.

1. `cli.usageCmd` loads `config.toml` and builds one `usage.Account` per account: its folder `~/.local/share/ccam/profiles/<name>/`, its session folder (to see whether it is running) and whether it was imported in place.
2. `Fetcher.FetchAll` starts one goroutine per account. All of them share one `http.Client`, so connections to the endpoint are reused. Connecting may take 2 s, and each account gets 3 s in total. The client does not follow redirects, since a redirect would resend the token somewhere else, and reads at most 1 MB of the answer.
3. As each account finishes, its row is redrawn. An account that fails or does not answer in time falls back to saved numbers, so one slow account never holds up the others.

Each account runs in its own goroutine against the endpoint; the one that times out shows its saved numbers instead.

<!-- motion:usage-fetch -->

`personal` and `work` answer in time and show live numbers. `side` takes longer than 3 s, so `fallback()` reads its two cache files and uses whichever was written last. `usage.Table()` draws a block per account and `usage.Best()` picks the account named on the last line.

For one account, `Fetcher.Fetch` makes these checks in order:

<!-- motion:usage-cache -->

1. With `--source cache`, it goes straight to `fallback()`.
2. Unless `--no-cache` is set, a live result in `~/.local/state/ccam/usage/<name>.json` younger than 60 seconds (`CacheTTL`) is returned as it is. Running `ccam usage` twice in a row costs one request per account.
3. `readCredentials()` reads `.credentials.json` (or the macOS Keychain) without changing it, and only the `claudeAiOauth` key is decoded. No subscription login means `not logged in`, or `no limits (API)` when `.claude.json` holds an API key.
4. An expired access token is not used. Without `--refresh`, the account shows its saved numbers marked as expired.
5. Otherwise ccam calls the endpoint. A 401 or 403 counts as an expired login, any other status or a body it cannot read counts as an error, and both fall back to saved numbers. A 429 means the endpoint is being asked too often: ccam writes `<name>.pause` with the time to ask again (`Retry-After`, or 5 minutes, at most an hour), shows saved numbers with `rate limited by the usage endpoint; asking again after 15:10`, and does not call the endpoint for that account until then, even with `--no-cache`. Retrying at once would only keep the limit going.
6. A good answer is written to `<name>.json` with the time it was measured.

`fallback()` reads both `<name>.json` and `<name>.statusline.json` and keeps the newer one. Before it shows saved numbers, it checks each window's reset time: a window that has reset since the numbers were saved shows 0%, because the old percentage no longer holds. That is why `personal` above shows `0%` with no reset time.

The endpoint has two response formats. The newer one has a `limits` array with `session` (5-hour), `weekly_all` and `weekly_scoped` (per model) entries. The older one has `five_hour`, `seven_day`, `seven_day_opus` and `seven_day_sonnet` objects. `ParseLive` reads both: the `limits` array wins, and the old objects fill in what it does not cover. A body with no limits at all is reported as a changed format, which `ccam doctor` tells you about.

### Renewing tokens

Refresh tokens rotate: renewing one gives you a new refresh token and invalidates the old one. If ccam renewed the login while Claude Code was running as the same account, that session would be signed out. So `ccam usage` reads tokens and never writes them, unless you pass `--refresh`. Then `Fetcher.refresh`:

1. refuses accounts registered in place with `ccam import`, since plain `claude` may use that login without ccam knowing;
2. refuses logins stored in the macOS Keychain;
3. takes a lock file, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (under `~/.local/state/ccam/locks/` when that variable is unset), so two ccam processes never renew the same token twice;
4. under the lock, refuses if a session of the account is running, and skips the request if another process already renewed the token;
5. posts the refresh token to the token endpoint and writes the new tokens back to `.credentials.json` with `fsutil.AtomicWrite`, keeping every other key in the file.

### The status line collector

Claude Code runs the status line command and pipes it a JSON document. ccam reads the model name and `rate_limits` from it:

<!-- motion:usage-statusline -->

1. `runStatusline` finds the account from `CCAM_PROFILE`, which ccam sets when it starts Claude Code, or else by matching `CLAUDE_CONFIG_DIR` against the account folders.
2. `ParseStatusLine` reads `rate_limits.five_hour` and `rate_limits.seven_day`. The JSON has no rate limits for API key accounts, or before the session's first response; then nothing is saved.
3. `SaveStatusLine` writes `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` uses it when the endpoint does not answer.
4. ccam prints `[name] · model · 5h N% · 7d N%`. With `--wrap`, it runs your command through `sh -c` with the same input and prints `[name]` followed by the command's first line instead. On Windows it runs the command the way Claude Code does: with Git Bash when it is installed, else with PowerShell.
5. When the 5-hour or weekly usage is at 90% or more and there is another account, the line ends with `→ ccam next --auto --continue`, for example `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

The wrapped command gets 2 seconds. After that it is killed and the line shows the account name alone, so a slow command cannot block Claude Code's status line.

`ccam statusline install` writes `statusLine` into each account's `settings.json`. When an account already has a status line command, ccam keeps it as `--wrap '<old command>'`. On Windows the path to ccam is written with forward slashes and `~` for the home folder (`~/.local/bin/ccam.exe statusline`), which both Git Bash and PowerShell understand. A `settings.json` shared through a group is written once, through its link, so every account in the group gets it. Running it again prints `(already set)` for files that already use ccam.

## Safety and edge cases

- `ccam usage` never changes a token file unless you pass `--refresh`. When a login has expired, run the account once with `ccam <name>` and Claude Code renews it.
- `--refresh` never touches running accounts, imported accounts or Keychain logins. A refused refresh shows the saved numbers and the reason.
- The usage endpoint is not publicly documented and may change. When it does, you still see saved numbers, and `ccam doctor` reports the format change.
- With no saved numbers and no answer, the account line shows `lookup failed:` and the reason, for example `no saved usage` with `--source cache`.
- The cache files hold only percentages and reset times, never tokens.

## Related

- [Usage limits and switching accounts](/docs/usage-limits)
- [Switching to the next account](/features/next)
- [Usage limit workflow](/guides/usage-limit-workflow)
- [Commands](/docs/commands)
