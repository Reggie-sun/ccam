---
title: ccam compared with shell aliases
description: Shell aliases with CLAUDE_CONFIG_DIR compared with ccam for running several Claude Code accounts, feature by feature.
updated: 2026-10-26
published: 2026-10-19
order: 1
translation: source
---

A shell alias per account is the simplest way to run several Claude Code accounts. ccam does the same thing underneath and adds sharing, switching and usage tracking.

## The alias setup

Add one line per extra account to `~/.bashrc` or `~/.zshrc`, open a new terminal and run the alias once to log in:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` stays your first account and `claude-work` starts the second. This is the method [Claude Code's documentation describes](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## The ccam setup

Install ccam and add the account, here with a shortcut command of its own:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` and `ccwork` start the account, and both work in scripts too.

## Side by side

| | Shell alias | ccam |
| --- | --- | --- |
| Separate logins | Yes | Yes |
| Setup per account | Edit `~/.bashrc`, log in | `ccam add` |
| Works in scripts and other programs | No, aliases are shell-only | Yes, `ccam work` and shortcuts such as `ccwork` are real commands |
| Shared settings, `CLAUDE.md`, extensions | Copy or link by hand | Yes, chosen per item, with backups |
| Shared MCP servers | By hand, per account | `ccam mcp add` once |
| Continue a conversation in another account | No | `ccam next --continue` |
| Usage of every account | No | `ccam usage` |
| Pick the account with the most room | No | `ccam next --auto` |
| Run a `-p` job again as another account at the limit | No | `ccam --fallback` |
| An account for each project folder | No | `ccam folder` |
| Account name in the status line | No | `ccam statusline install` |
| Repair broken links | By hand | `ccam doctor --fix` |
| Extra software | None | One static binary |

Choose aliases if you have two accounts and want nothing shared between them. Choose ccam if you want the accounts to feel like one setup. You can start with aliases and move to ccam later: `ccam import` registers existing `~/.claude-*` folders where they are.

[How to use multiple accounts on one computer](/guides/multiple-claude-code-accounts) explains both approaches step by step, and [ccam and other account switchers](/compare/account-switchers) compares ccam with tools that swap one login.
