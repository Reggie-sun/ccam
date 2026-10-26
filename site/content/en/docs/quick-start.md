---
title: Quick start
metaTitle: "Quick start: add and switch Claude Code accounts"
description: Add two Claude Code accounts with ccam add, give them shortcuts such as ccwork, switch between them, and pick what they share.
updated: 2026-10-26
published: 2026-10-19
order: 2
translation: source
---

Add each account with `ccam add`, then start Claude Code with `ccam <name>` or with a shortcut you name yourself, such as `ccwork`. That is all most people need.

Already use Claude Code? Read [Already using ~/.claude?](#already-using-claude) at the end of this page before you start.

## 1. Add your first account

```sh
ccam add
```

The wizard asks for a name such as `personal` or `work`, which items the account should share with your other accounts, and whether to create a shortcut command (see step 3). Then it runs `claude auth login`, which opens the browser so you can sign in. On a remote machine, your browser may end up on an `http://localhost:…/callback` page that does not load; paste that address into the terminal and ccam finishes the sign-in. When sign-in is complete, you see:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Bring over your existing settings

If you have used Claude Code before, your settings, `CLAUDE.md` and extensions are in `~/.claude`, and your MCP servers are in `~/.claude.json`. Copy them into the shared `default` group so that every account gets them:

```sh
ccam migrate
```

ccam shows what it will copy and asks first. It only reads `~/.claude`, never copies a login, and leaves conversation history out unless you add `--history`. See [copying an existing setup](/docs/sharing#copy-an-existing-setup) for details.

## 3. Add the second account

```sh
ccam add work
```

The wizard skips the name, which you gave on the command line, and shows the groups to join with `default` checked. This time, give the account a shortcut command of its own, such as `ccwork`, so you can start it with one word. After the groups, the wizard asks:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Press `y`, then type a name, or press Enter to use the suggested one:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Then sign in with the other account in the browser. ccam confirms the account and its shortcut:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Switch accounts

`ccam work` starts Claude Code as the work account. If you skipped the shortcut question, add a shortcut at any time:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Each of these commands now starts Claude Code as the work account (`2` is its number in `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Running `ccam` alone opens the [dashboard](/features/dashboard): move to an account with the arrow keys or its number, then press Enter.

Everything after the shortcut or the account name goes to Claude Code unchanged, so you keep using the options you know:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Choose what to share

By default, MCP servers, `settings.json`, the global `CLAUDE.md` and extensions (agents, commands, skills and plugins) are shared between accounts through the `default` group. Conversation history is not. Add it to the group if you want to continue conversations across accounts:

```sh
ccam group set default --add history
```

ccam shows what it will move for each account before it changes anything. See [sharing with groups](/features/share) for details.

## 6. When an account reaches its limit

```sh
ccam usage
ccam next --continue
```

`ccam usage` shows the 5-hour and weekly usage of every account. `ccam next --continue` opens the same conversation in the next account. In scripts, `ccam --fallback work -p "…"` switches by itself. See [usage limits](/docs/usage-limits).

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.

## Already using ~/.claude?

Register it instead of signing in again:

```sh
ccam import
```

The folder stays where it is, and plain `claude` keeps working as before.
