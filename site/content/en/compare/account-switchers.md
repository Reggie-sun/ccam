---
title: ccam and other Claude Code account switchers
description: How ccam compares with claude-swap, CCS, clauth, CAAM, cux and other tools for switching Claude Code accounts, and which approach fits which use.
updated: 2026-10-19
published: 2026-10-19
order: 2
translation: source
---

Tools for switching Claude Code accounts take one of two approaches. Most swap the login inside the one `~/.claude` folder, so one account is active at a time. Others, ccam among them, give each account its own folder with `CLAUDE_CONFIG_DIR`, so every account stays signed in and several can run at once. ccam adds shared settings and MCP servers, and continues a conversation in another account's folder.

## Swapping one login or separate folders?

**Swapping a login** keeps one set of settings, MCP servers and history, because there is only one folder. Switching is quick, and a running session can pick up the new login. But only one account is active at a time, for every terminal and editor at once.

**Separate folders** is the method [Claude Code's documentation describes](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Every account stays logged in, and you can run a work and a personal session side by side. The price is that each folder starts with its own settings, MCP servers and history, unless a tool shares them.

## How the tools compare

Taken from each project's own documentation on 2026-10-19. "Not documented" means the project does not say, not that the feature is missing.

| Tool | Approach | Accounts at once | 5-hour and weekly usage | When a limit is reached | Platforms |
| --- | --- | --- | --- | --- | --- |
| ccam | Separate folders | Yes | Yes, `ccam usage` | You run `ccam next --auto --continue` to continue the conversation in the account with the most room | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Swaps the login; separate folders with `cswap run` | Yes, with `cswap run` | Yes | `cswap auto` switches the login at a threshold; the running session continues | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Separate folders | Yes | Partly, in a macOS menu bar app | Account pools pass requests to the next account through a local proxy | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Swaps the login; separate folders with `clauth start` | Yes, with `clauth start` | Yes | Switches down a fallback chain | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Swaps the login, or isolated profiles | Yes, with isolated profiles | Yes | `caam run` retries non-interactive commands on another account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Swaps the login and wraps `claude` | Not documented | Yes | Swaps the login and resumes the same conversation | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Swaps the login | No | Yes | Not documented | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Swaps the login | No | Yes | Optional switch at 100% | macOS |

ccswap is a fork of claude-swap that also manages Codex logins. Several of these tools also handle other coding agents, such as Codex or Gemini.

## What ccam focuses on

- Every account stays logged in. Switching never logs one out, and `ccam work` and `ccam personal` can run in two terminals at once.
- MCP servers, `settings.json`, `CLAUDE.md`, agents, skills and plugins are shared between folders through links, item by item and account by account. Logins are never shared.
- With shared history, one command, `ccam next --auto --continue`, opens the same conversation in the folder of the account with the most room.
- ccam starts Claude Code and replaces itself with it, so nothing runs in the background. It does not watch your usage or switch on its own; you decide when.

## Where other tools do more

- claude-swap, CCS, clauth, CAAM and cux run on Windows without WSL. ccam does not.
- claude-swap, clauth and cux switch on their own when an account crosses a usage threshold. ccam leaves the switch to you.
- CCSwitcher, Claude Switcher and claude-swap show usage in the macOS menu bar.
- CCS, CAAM, ccswap and Claude Switcher also manage Codex and other coding agents.

## Which one fits?

Choose a login-swapping tool if you use one account at a time, want switching to happen automatically, or are on Windows without WSL. Choose ccam if you want accounts side by side, one shared setup and the same conversation in the next account. If you need nothing shared, [plain shell aliases](/compare/aliases) are enough.

Whether you may use several accounts depends on each account's terms. Check them before you do.
