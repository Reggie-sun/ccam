---
title: Introducing ccam
description: Why ccam exists, what it does for people with several Claude Code accounts, and how it stays out of the way.
updated: 2026-10-26
published: 2026-10-19
order: 1
translation: source
---

ccam is a small command-line tool for people who use more than one Claude Code account on the same computer: a personal and a work account, or several accounts for one long project.

## The problem

Claude Code supports several accounts through [`CLAUDE_CONFIG_DIR`](/guides/claude-config-dir): one folder per account. In practice, each folder is an island. Settings, `CLAUDE.md`, MCP servers, agents and skills have to be set up again in every folder and kept in sync by hand. Conversations stay where they started. And when one account reaches its usage limit, you find out by trying.

## What ccam does

- `ccam add` creates the folder, runs the login and, if you want one, gives the account a shortcut you name yourself, such as `ccwork`.
- `ccam work` and `ccwork` start Claude Code as that account and pass every argument through unchanged.
- MCP servers, settings, `CLAUDE.md` and extensions are shared through links, item by item and account by account. Logins are never shared.
- `ccam usage` shows the 5-hour and weekly usage of every account; `ccam next --continue` picks up the conversation in the next account.

## Staying out of the way

A tool you run before every session should not cost anything you notice. ccam reads one configuration file, checks the shared links with one or two system calls each, and then replaces itself with `claude`. It adds about 4 milliseconds, and by default it is not running while you work.

It is also careful with your files. Every change to sharing prints a plan before it starts and keeps the originals in a backup folder. `ccam doctor` checks links, permissions and logins, and repairs what it safely can.

## Try it

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

The [quick start](/docs/quick-start) walks through adding two accounts, and [how to use multiple Claude Code accounts](/guides/multiple-claude-code-accounts) explains the approach with and without ccam.

ccam is open source and unofficial: it is not affiliated with or endorsed by Anthropic. Whether you may use more than one account depends on the terms that apply to each one. Feedback and translations are welcome on [GitHub](https://github.com/jhyunwoo/ccam).
