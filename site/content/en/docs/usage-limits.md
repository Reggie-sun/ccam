---
title: Usage limits and switching accounts
description: See the 5-hour and weekly usage of every Claude Code account with ccam usage, and open the account with the most room with ccam next --auto.
updated: 2026-10-21
published: 2026-10-19
order: 5
translation: source
---

`ccam usage` shows how much of the 5-hour limit, the weekly limit and each model's own weekly limit (such as Fable) every account has used. `ccam next --auto` opens the account with the most room, and `--continue` carries the current conversation over.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Each account gets a line with how old its numbers are, then one bar per limit with the time until it resets: the 5-hour window, the weekly window for all models, and each model's own weekly window, such as Fable. From 80% a bar turns yellow, and at 100% red. The last line names the account `ccam next --auto` would open.

## Where the numbers come from

ccam asks Anthropic's usage endpoint for every account at the same time, using each account's own login. It waits at most 3 seconds per account and shows each account as soon as it answers.

This endpoint is not publicly documented and may change. When a request fails, ccam shows the last known numbers and marks them as cached. A window that has reset since those numbers were saved shows 0%. `ccam doctor` tells you if the response format has changed.

Numbers fetched less than 60 seconds ago are reused. Options:

| Option | Effect |
| --- | --- |
| `--watch [seconds]` | Refresh continuously, every 60 seconds by default |
| `--json` | Machine-readable output |
| `-v` | Also show where each account's numbers came from |
| `--no-cache` | Always ask again |
| `--source cache` | Never ask; show saved numbers only |

## Expired logins

When an account's access token has expired, ccam does not renew it, because renewing it while a Claude Code session uses the same login could sign that session out. Start the account once with `ccam <name>` and Claude Code renews it.

`ccam usage --refresh` renews expired tokens itself, but only for accounts that are not running. It never renews logins stored in the macOS Keychain or in an imported `~/.claude`.

## Usage in the status line

```sh
ccam statusline install
```

Claude Code then shows the account name and its usage in the status line, for example `[work] · Opus · 5h 24% · 7d 41%`. From 90% on, the line also suggests `ccam next --auto --continue`. If you already have a status line command, ccam keeps it and puts the account name in front. The status line also saves its numbers, which `ccam usage` uses when the endpoint does not answer.

## Switching

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` skips accounts that are not logged in, and `--auto` skips accounts at 100% and accounts of another organization. When the two accounts share conversation history, `--continue` opens the same conversation; otherwise Claude Code copies it into the new account. See [sharing](/docs/sharing) and [switching to the next account](/features/next#continue).

For scripts and scheduled jobs that run `claude -p`, put `--fallback` before the account name. When that account reaches its limit, ccam runs the same command again as the account with the most room, and only the output of the run that finished is printed:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

See [Run again at the limit](/features/next#run-again-at-the-limit) for how ccam decides.

The [usage limit workflow guide](/guides/usage-limit-workflow) walks through a full example.
