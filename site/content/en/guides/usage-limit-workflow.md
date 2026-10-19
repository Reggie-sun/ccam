---
title: What to do when Claude Code reaches its usage limit
description: Check how much room each Claude Code account has left and continue the same conversation on another account when one reaches its limit.
updated: 2026-10-19
order: 3
translation: source
---

When one account reaches its limit, continue the same conversation in another account you own: share conversation history once, then run `ccam next --continue`. If you have only one account, wait for the reset time Claude Code shows.

## Before it happens

Set up the accounts so that the switch takes one command:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

The status line now shows the account and its 5-hour usage while you work, so a limit does not come as a surprise.

## When the limit is reached

Leave Claude Code (`/exit`), then:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Open the conversation in the account with the most room:

```sh
ccam next --auto --continue
```

Claude Code starts as `work` with the most recent conversation of the current folder, including everything said so far.

## Things to keep in mind

- Run `ccam next --continue` in the same folder you worked in; conversations are stored per project folder.
- Avoid running two accounts on the same conversation at the same time. ccam warns you if you try.
- The other account must also share conversation history, in the same group.
- Whether you may switch between accounts like this depends on each account's terms.

See [usage limits](/docs/usage-limits) for every option of `ccam usage` and `ccam next`.
