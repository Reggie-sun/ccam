---
title: What to do when Claude Code hits its usage limit
description: Wait for the reset or continue the same conversation on another Claude Code account right away. How the 5-hour and weekly limits work, step by step.
updated: 2026-10-19
published: 2026-10-19
order: 3
translation: source
---

When an account reaches its usage limit (often called a rate limit), you have two choices: wait until the limit resets, or continue the same conversation now on another account you own. With ccam the second takes one command, `ccam next --auto --continue`, once conversation history is shared between the accounts.

## Which limits does Claude Code have?

On Pro and Max plans, Claude Code has a [five-hour session limit and a weekly limit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Both are shared with Claude on the web and in the apps, and usage counts against both at the same time. When one runs out, Claude Code [shows a message](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) such as "You've hit your session limit · resets 3:45pm". Switching models does not help, because the limits cover all models.

## Option 1: wait for the reset

If you have one account, wait. Recent versions of Claude Code [wait in the open session](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) and continue the task on their own once the limit resets. `/usage` inside Claude Code shows how much of your plan's limits you have used.

## Option 2: continue on another account now

If you have a second account, you can carry on straight away with the same conversation. Set it up once:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

The status line now shows the account and its 5-hour usage while you work, so a limit does not come as a surprise.

When the limit is reached, leave Claude Code (`/exit`), then:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Open the conversation in the account with the most room:

```sh
ccam next --auto --continue
```

Claude Code starts as `work` with the most recent conversation of the current folder, including everything said so far.

## Does ccam switch accounts automatically?

No. ccam is not running while you work: it starts Claude Code and replaces itself with it. You decide when to switch, and `ccam next --auto` picks the account with the lowest 5-hour usage and skips accounts that are at 100% or not logged in. Some other tools swap the login of a running session instead; [ccam and other account switchers](/compare/account-switchers) explains the difference.

## Things to keep in mind

- Run `ccam next --continue` in the same folder you worked in; conversations are stored per project folder.
- Avoid running two accounts on the same conversation at the same time. ccam warns you if you try.
- The other account must also share conversation history, in the same group.
- The first message on the new account may use more of its limit than usual, because that account has not cached the conversation yet.
- Whether you may switch between accounts like this depends on each account's terms.

See [usage limits](/docs/usage-limits) for every option of `ccam usage` and `ccam next`.
