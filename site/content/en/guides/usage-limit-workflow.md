---
title: What to do when Claude Code hits its usage limit
description: Wait for the reset or continue the same conversation on another Claude Code account right away. How the 5-hour and weekly limits work, step by step.
updated: 2026-10-25
published: 2026-10-19
order: 3
translation: source
---

When an account reaches its usage limit (often called a rate limit), you have two choices: wait until the limit resets, or continue the same conversation now on another account you own. With ccam the second takes one command, `ccam next --auto --continue`.

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
ccam group set default --add history
ccam statusline install
```

The status line now shows the account and its 5-hour usage while you work, so a limit does not come as a surprise.

When the limit is reached, leave Claude Code (`/exit`), then:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Open the conversation in the account with the most room:

```sh
ccam next --auto --continue
```

Claude Code starts as `work` with the most recent conversation of the current folder, including everything said so far. Sharing history (`ccam group set default --add history` above) keeps it one conversation; without it, Claude Code copies the conversation into `work`, and ccam tells you how many tokens that first reply reads again.

With the [ccam plugin](/docs/install#claude-code-plugin) in Claude Code, a session that stops at its limit shows a notification and records the full window, so `ccam next --auto` does not pick that account again before its reset, and `/ccam:switch` (or `/ccam:switch work`) notes the conversation for `ccam next --continue` without leaving anything to guess.

**Team and Enterprise seats.** Automatic choice (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n`, and [`auto_switch`](/features/next#auto-switch) when you turn it on) also switches between seats of the same Claude Team or Enterprise organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Using automatic choice between seats may therefore violate Anthropic's terms. Whether you use it there is your call: check with your organization and Anthropic before you rely on it. Personal plans (Pro, Max) are covered by the [terms note in the FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Does ccam switch accounts automatically?

Only if you turn it on. By default ccam is not running while you work: it starts Claude Code and replaces itself with it, so you decide when to switch, and `ccam next --auto` picks the account with the lowest 5-hour usage, skips accounts that are at 100% or not logged in, and stays within one Team or Enterprise organization, or among personal plans. With `relay = true` in `config.toml`, ccam waits beside the session instead, so `/ccam:switch` carries on by itself, keeping the session's options, permission mode and effort level. With `auto_switch = true` (Linux and macOS; experimental), ccam also makes the switch itself: when Claude Code stops the session at its limit and another account like it has room, ccam ends the session and continues the conversation there ([conditions and costs](/features/next#auto-switch)). Some other tools swap the login of a running session instead; [ccam and other account switchers](/compare/account-switchers) explains the difference.

## Things to keep in mind

- Run `ccam next --continue` in the same folder you worked in; conversations are stored per project folder.
- Avoid running two accounts on the same conversation at the same time. ccam warns you if you try.
- Without shared history, the conversation is copied to the other account and the two copies go their own ways; the checkpoints `/rewind` restores code from stay with the first account.
- The first message on the new account may use more of its limit than usual, because that account has not cached the conversation yet. `ccam next --continue --brief` starts it on a short summary instead.
- Whether you may switch between accounts like this depends on each account's terms.

See [usage limits](/docs/usage-limits) for every option of `ccam usage` and `ccam next`.
