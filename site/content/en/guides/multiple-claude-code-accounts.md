---
title: How to use two Claude Code accounts on one computer
description: Run several Claude Code accounts side by side with CLAUDE_CONFIG_DIR, by hand with shell aliases or automatically with ccam.
updated: 2026-10-19
order: 1
translation: source
---

Give each account its own configuration folder. Claude Code reads the folder from the `CLAUDE_CONFIG_DIR` environment variable, and each folder holds its own login, settings and history. You can set this up by hand with a shell alias, or let ccam do it.

## By hand, with an alias

This is the approach Claude Code's documentation describes. Add a line per extra account to `~/.bashrc` or `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Open a new terminal and run `claude-work`. Claude Code asks you to log in, because the folder is new. From then on, `claude` is your first account and `claude-work` the second.

This works well for two accounts. Things to know:

- Each folder starts empty: no settings, no `CLAUDE.md`, no MCP servers, no agents or skills. You copy or link what you need, and keep the copies in step yourself.
- Aliases exist only in interactive shells. Scripts, Makefiles and other programs do not see them.
- MCP servers live in `.claude.json`, next to the login, so they cannot simply be linked between folders.
- Conversations stay in the account where they started, so you cannot continue one in another account.
- There is nothing that shows the usage of all accounts at once.

## With ccam

ccam automates the same idea and adds the missing parts:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Now `claude1` and `claude2`, or `ccs personal` and `ccs work`, start the two accounts. Because they are real programs rather than aliases, they also work in scripts. Settings, `CLAUDE.md`, MCP servers and extensions are shared from the start, and you can turn on shared conversation history to continue a conversation on the other account with `ccs next --continue`.

If you already use `~/.claude`, run `ccam import` instead of adding that account again.

## Which one to choose

For two accounts and nothing shared, an alias is enough. When you want shared settings and MCP servers, more accounts, usage at a glance or continuing conversations, ccam saves the bookkeeping. The [comparison](/compare/aliases) lists the differences side by side.
