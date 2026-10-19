---
title: Frequently asked questions
description: Answers about ccam, Claude Code accounts, usage limits, sharing, privacy and supported systems.
updated: 2026-10-19
order: 6
translation: source
---

## Is ccam an official Anthropic tool?

No. ccam is an unofficial open-source project and is not affiliated with or endorsed by Anthropic. It uses the `CLAUDE_CONFIG_DIR` setting that Claude Code documents for running several accounts side by side.

## Am I allowed to use several Claude Code accounts?

That depends on the terms of each account and, for work accounts, your organisation's policy. Check them before you use several accounts. ccam only makes switching between accounts you already have easier.

## Does ccam get around usage limits?

No. Each account keeps its own limits. ccam shows them and lets you continue on another account you own when one runs out.

## How do I use two Claude Code accounts on one computer?

Give each account its own configuration folder with `CLAUDE_CONFIG_DIR`. ccam does this for you: run `ccam add` for each account, then start them with `claude1` and `claude2`. The [guide](/guides/multiple-claude-code-accounts) also shows how to do it by hand.

## How do I continue a conversation on another account?

Share conversation history between the two accounts with `ccam share <name> --history on`, then run `ccam next --continue`. Claude Code opens the most recent conversation of the current folder in the next account.

## Does ccam see or send my login?

ccam never copies or shares a login. It reads an account's access token only to ask Anthropic's usage endpoint for that account's usage, when you run `ccam usage`, `ccam next --auto` or `ccam doctor`. ccam has no telemetry and sends nothing anywhere else.

## Where are my accounts stored?

In `~/.local/share/ccam/profiles/<name>/`, one folder per account. On macOS, Claude Code keeps the login itself in the Keychain, under a name derived from that folder.

## What happens to my existing ~/.claude?

Nothing, unless you run `ccam import`. Import registers the folder where it is, so plain `claude` and `ccam <name>` keep using the same account and settings.

## Which systems does ccam support?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine and others) and macOS. On Windows, use it inside WSL.

## Why does an MCP server show up in only one account?

Servers added inside a session with `claude mcp add` are stored for that account only. ccam tells you the next time you start the account; move the server to the shared group with `ccam mcp import <account> <server>`.

## Does ccam slow Claude Code down?

Starting Claude Code through ccam adds about 4 milliseconds, measured with `scripts/perf-check.sh` on a laptop. ccam reads one small file, checks the shared links and then replaces itself with `claude`, so it is not running while you work.

## How do I uninstall ccam?

Delete `ccam` and the `claude1`… links from `~/.local/bin`, then `~/.config/ccam`, `~/.local/share/ccam` and `~/.local/state/ccam` if you no longer need the accounts. See [install](/docs/install#uninstall).
