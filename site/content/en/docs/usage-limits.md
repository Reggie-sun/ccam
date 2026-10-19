---
title: Usage limits and switching accounts
description: See the 5-hour and weekly usage of every Claude Code account with ccam usage, and open the account with the most room with ccam next --auto.
updated: 2026-10-19
published: 2026-10-19
order: 5
translation: source
---

`ccam usage` shows how much of the 5-hour and weekly limits each account has used. `ccam next --auto` opens the account with the most room, and `--continue` carries the current conversation over.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

The columns are: account, 5-hour usage, time until it resets, weekly usage, time until it resets, and how old the numbers are.

## Where the numbers come from

ccam asks Anthropic's usage endpoint for every account at the same time, using each account's own login. It waits at most 3 seconds per account and shows each account as soon as it answers.

This endpoint is not publicly documented and may change. When a request fails, ccam shows the last known numbers and marks them as cached. A window that has reset since those numbers were saved shows 0%. `ccam doctor` tells you if the response format has changed.

Numbers fetched less than 60 seconds ago are reused. Options:

| Option | Effect |
| --- | --- |
| `--watch [seconds]` | Refresh continuously, every 60 seconds by default |
| `--json` | Machine-readable output |
| `-v` | Also show model-specific weekly limits and where each number came from |
| `--no-cache` | Always ask again |
| `--source cache` | Never ask; show saved numbers only |

## Expired logins

When an account's access token has expired, ccam does not renew it, because renewing it while a Claude Code session uses the same login could sign that session out. Start the account once with `ccam <name>` and Claude Code renews it.

`ccam usage --refresh` renews expired tokens itself, but only for accounts that are not running. It never renews logins stored in the macOS Keychain or in an imported `~/.claude`.

## Usage in the status line

```sh
ccam statusline install
```

Claude Code then shows the account name and its usage in the status line, for example `[work] · Opus · 5h 24% · 7d 41%`. If you already have a status line command, ccam keeps it and puts the account name in front. The status line also saves its numbers, which `ccam usage` uses when the endpoint does not answer.

## Switching

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` skips accounts that are not logged in, and `--auto` skips accounts at 100%. To continue a conversation, both accounts must share conversation history; see [sharing](/docs/sharing).

The [usage limit workflow guide](/guides/usage-limit-workflow) walks through a full example.
