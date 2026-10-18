---
title: CLAUDE_CONFIG_DIR explained
description: What CLAUDE_CONFIG_DIR does in Claude Code, what is stored in the folder, where macOS keeps the login, and which files can be shared between accounts.
updated: 2026-10-19
order: 2
translation: source
---

`CLAUDE_CONFIG_DIR` tells Claude Code which folder to use instead of `~/.claude`. Everything account-specific lives there: the login, settings, MCP servers, plugins and conversation history. Pointing it at different folders gives you separate accounts.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## What is inside the folder

Checked with Claude Code 2.1.302:

| Path | Contents | Shareable |
| --- | --- | --- |
| `.credentials.json` | Login tokens (Linux), MCP OAuth tokens | No |
| `.claude.json` | Account details, user MCP servers, per-project state | No (MCP servers separately) |
| `settings.json` | Permissions, hooks, model, status line, plugins | Yes |
| `CLAUDE.md` | Global instructions | Yes |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensions | Yes |
| `projects/` | Conversations as `.jsonl` files, per-project memory | Yes |
| `history.jsonl` | Prompt history | Yes |
| `sessions/`, `shell-snapshots/`, `cache/` and others | Runtime state | No |

## Two details that catch people out

**Without `CLAUDE_CONFIG_DIR`, `.claude.json` is not inside `~/.claude`.** Claude Code then keeps it at `~/.claude.json`, next to the folder. If you start using `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code looks for `~/.claude/.claude.json` instead and does not find your account. `ccam import` handles this by linking the two.

**On macOS the login is in the Keychain, not in the folder.** Claude Code stores it under the name `Claude Code-credentials`, followed by a short hash of the folder path when `CLAUDE_CONFIG_DIR` is set. Separate folders therefore get separate Keychain entries. If the Keychain refuses the write, Claude Code falls back to `.credentials.json` in the folder.

## Sharing between folders

Symbolic links work: Claude Code reads and writes through a linked `settings.json` or `projects/` folder and keeps the link. MCP servers are the exception, because they share `.claude.json` with the login; pass them with `--mcp-config` instead.

ccam builds on exactly these rules. See [sharing](/docs/sharing) for how it links folders, and [the guide to multiple accounts](/guides/multiple-claude-code-accounts) for the bigger picture.
