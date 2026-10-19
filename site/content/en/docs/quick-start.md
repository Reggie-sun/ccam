---
title: Quick start
description: Add two Claude Code accounts with ccam add, switch between them with claude1, claude2 or ccam <name>, and pick what they share.
updated: 2026-10-19
order: 2
translation: source
---

Add each account with `ccam add`, then start Claude Code with `claude1`, `claude2` or `ccam <name>`. That is all most people need.

## 1. Add your first account

```sh
ccam add
```

The wizard asks for a name such as `personal` or `work`, and which items the account should share with your other accounts. Then it runs `claude auth login`, which opens the browser so you can sign in. When it finishes you see:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam prints its own messages in English; the commands and options are the same in every language.

## 2. Add the second account

```sh
ccam add work
```

Sign in with the other account in the browser. It becomes `claude2`.

## 3. Switch

Any of these start Claude Code as the work account:

```sh
claude2
ccam work
ccam 2
```

Everything after the account name goes to Claude Code unchanged, so you keep using the options you know:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Running `ccam` alone opens a picker: use the arrow keys and Enter, or press the account's number.

## 4. Choose what to share

By default, MCP servers, `settings.json`, the global `CLAUDE.md`, and agents, commands, skills and plugins are shared between accounts. Conversation history is not. Turn it on for both accounts if you want to continue conversations across them:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam shows what it will move before it changes anything. See [sharing](/docs/sharing) for the details.

## 5. When an account reaches its limit

```sh
ccam usage
ccam next --continue
```

`ccam usage` shows the 5-hour and weekly usage of every account. `ccam next --continue` opens the same conversation in the next account. See [usage limits](/docs/usage-limits).

## Already using ~/.claude?

Register it instead of signing in again:

```sh
ccam import
```

The folder stays where it is, and plain `claude` keeps working as before.
