---
title: Quick start
metaTitle: "Quick start: add and switch Claude Code accounts"
description: Add two Claude Code accounts with ccam add, give them shortcuts such as ccwork, switch between them, and pick what they share.
updated: 2026-10-20
published: 2026-10-19
order: 2
translation: source
---

Add each account with `ccam add`, then start Claude Code with `ccam <name>` or with a shortcut you name yourself, such as `ccwork`. That is all most people need.

## 1. Add your first account

```sh
ccam add
```

The wizard asks for a name such as `personal` or `work`, which items the account should share with your other accounts, and whether to create a shortcut command (see step 4). Then it runs `claude auth login`, which opens the browser so you can sign in. On a remote machine, the browser may end on a `http://localhost:…/callback` page that does not load; paste that address into the terminal and ccam finishes the sign-in. When it finishes you see:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Bring over your existing settings

If you used Claude Code before, your settings, `CLAUDE.md`, extensions and MCP servers are in `~/.claude`. Copy them into the shared settings, so that every account gets them:

```sh
ccam migrate
```

ccam shows what it will copy and asks first. It only reads `~/.claude`, never copies a login, and leaves conversation history out unless you add `--history`. See [copying an existing setup](/docs/sharing#copy-an-existing-setup) for the details.

## 3. Add the second account

```sh
ccam add work
```

Sign in with the other account in the browser.

## 4. Switch accounts

`ccam work` starts Claude Code as the work account. For something shorter, give the account a shortcut command of its own, such as `ccwork`. After the shared items, the wizard asks:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Press `y`, then type a name, or press Enter to use the suggested one:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

After the login, ccam confirms it:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

If you skipped the question, add a shortcut at any time:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Now any of these start Claude Code as the work account:

```sh
ccwork
ccam work
ccam 2
```

Running `ccam` alone opens a picker: use the arrow keys and Enter, or press the account's number.

Everything after the shortcut or the account name goes to Claude Code unchanged, so you keep using the options you know:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Choose what to share

By default, MCP servers, `settings.json`, the global `CLAUDE.md`, and agents, commands, skills and plugins are shared between accounts. Conversation history is not. Turn it on for both accounts if you want to continue conversations across them:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam shows what it will move before it changes anything. See [sharing](/docs/sharing) for the details.

## 6. When an account reaches its limit

```sh
ccam usage
ccam next --continue
```

`ccam usage` shows the 5-hour and weekly usage of every account. `ccam next --continue` opens the same conversation in the next account. In scripts, `ccam --fallback work -p "…"` switches by itself. See [usage limits](/docs/usage-limits).

## Already using ~/.claude?

Register it instead of signing in again:

```sh
ccam import
```

The folder stays where it is, and plain `claude` keeps working as before.
