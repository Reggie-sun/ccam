---
title: ccam compared with shell aliases
description: Shell aliases with CLAUDE_CONFIG_DIR against ccam for running several Claude Code accounts, side by side.
updated: 2026-10-19
order: 1
translation: source
---

A shell alias per account is the simplest way to run several Claude Code accounts. ccam does the same thing underneath and adds sharing, switching and usage tracking.

| | Shell alias | ccam |
| --- | --- | --- |
| Separate logins | Yes | Yes |
| Setup per account | Edit `~/.bashrc`, log in | `ccam add` |
| Works in scripts and other programs | No, aliases are shell-only | Yes, `claude2` is a real command |
| Shared settings, `CLAUDE.md`, extensions | Copy or link by hand | Per item, with backups |
| Shared MCP servers | By hand, per account | `ccam mcp add` once |
| Continue a conversation in another account | No | `ccs next --continue` |
| Usage of every account | No | `ccs usage` |
| Pick the account with the most room | No | `ccs next --auto` |
| Account name in the status line | No | `ccam statusline install` |
| Repair broken links | By hand | `ccam doctor --fix` |
| Extra software | None | One static binary |

Choose aliases if you have two accounts and want nothing shared between them. Choose ccam if you want the accounts to feel like one setup. You can start with aliases and move to ccam later: `ccam import` registers existing `~/.claude-*` folders where they are.

[How to use two accounts on one computer](/guides/multiple-claude-code-accounts) explains both approaches step by step.
