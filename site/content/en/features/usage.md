---
title: Usage limits
description: ccam usage shows the 5-hour and weekly limits of every account at once, and ccam statusline puts the account name and usage in Claude Code.
updated: 2026-10-20
order: 7
translation: source
---

`ccam usage` shows how much of the 5-hour and weekly limits each account has used, with the time until each window resets. Use it before you start a long task, or when Claude Code tells you a limit is near, to see which account has room left. `ccam statusline` is the status line command for Claude Code: it shows the account name and usage inside Claude Code and saves the numbers for `ccam usage`.

## Usage

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

On a terminal the table appears at once with `fetching…` in each row, and each row is redrawn as its account answers. `As of` says how old the numbers are; `(cached)` means they come from a saved file, not from the endpoint. The last line names the account with the lowest 5-hour usage below 100%, and only appears when you have more than one account. An account that is not logged in shows `not logged in → ccam login <name>`, and an API key account shows `no limits (API)`.

After `ccam statusline install`, Claude Code shows a line such as `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--json` | Print the reports as JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | off |
| `--watch[=N]` | Redraw every `N` seconds until you press Ctrl+C. `--watch` alone, or a value of 0 or less, means 60 | off |
| `-v`, `--verbose` | Show per-model weekly limits and a `source:` line (`live`, `cache` or `statusline`) | off |
| `--no-cache` | Ignore the 60-second cache and ask the endpoint again | off |
| `--source live\|cache` | `cache` never asks the endpoint and shows saved numbers only; `live` is the default behaviour | live |
| `--refresh` | Renew expired access tokens, only for accounts that are not running | off |
| `--dry-run` | Accepted: `ccam usage` changes nothing but its own cache. Refused together with `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

Any other `--source` value fails with `--source must be live or cache`. `ccam usage` takes no arguments.

`ccam statusline` reads Claude Code's status line JSON from standard input. Its one option is `--wrap <command>`: an existing status line command whose first output line gets the account name in front. `ccam statusline install` takes no options. Both refuse `--dry-run` with `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## How it works

The numbers come from two places: Anthropic's undocumented usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, called with each account's own login, and the values Claude Code hands to its status line command. All of it lives in `internal/usage`.

1. `cli.usageCmd` loads `config.toml` and builds one `usage.Account` per account: its folder `~/.local/share/ccam/profiles/<name>/`, its session folder (to see whether it is running) and whether it was imported in place.
2. `Fetcher.FetchAll` starts one goroutine per account. All of them share one `http.Client`, so connections to the endpoint are reused. Connecting may take 2 s, and each account gets 3 s in total. The client does not follow redirects, since a redirect would resend the token somewhere else, and reads at most 1 MB of the answer.
3. As each account finishes, its row is redrawn. An account that fails or does not answer in time falls back to saved numbers, so one slow account never holds up the others.

Each account runs in its own goroutine against the endpoint; the one that times out shows its saved numbers instead.

<!-- motion:usage-fetch -->

`personal` and `work` answer in time and show live numbers. `side` takes longer than 3 s, so `fallback()` reads its two cache files and uses whichever was written last. `usage.Table()` draws the rows and `usage.Best()` picks the account named on the last line.

For one account, `Fetcher.Fetch` makes these checks in order:

<!-- motion:usage-cache -->

1. With `--source cache`, it goes straight to `fallback()`.
2. Unless `--no-cache` is set, a live result in `~/.local/state/ccam/usage/<name>.json` younger than 60 seconds (`CacheTTL`) is returned as it is. Running `ccam usage` twice in a row costs one request per account.
3. `readCredentials()` reads `.credentials.json` (or the macOS Keychain) without changing it, and only the `claudeAiOauth` key is decoded. No subscription login means `not logged in`, or `no limits (API)` when `.claude.json` holds an API key.
4. An expired access token is not used. Without `--refresh`, the account shows its saved numbers marked as expired.
5. Otherwise ccam calls the endpoint. A 401 or 403 counts as an expired login, any other status or a body it cannot read counts as an error, and both fall back to saved numbers.
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
- With no saved numbers and no answer, the row shows `lookup failed:` and the reason, for example `no saved usage` with `--source cache`.
- The cache files hold only percentages and reset times, never tokens.

## Related

- [Usage limits and switching accounts](/docs/usage-limits)
- [Switching to the next account](/features/next)
- [Usage limit workflow](/guides/usage-limit-workflow)
- [Commands](/docs/commands)
