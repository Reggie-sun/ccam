---
title: Inspecting extensions
description: ccam ext shows the skills, commands, plugins, MCP servers and hooks each account loads, where each comes from, and what is broken.
updated: 2026-10-20
order: 14
translation: source
---

`ccam ext` answers "why does this skill (or MCP server) work in one account and not the other?". It lists, for every account, the extensions Claude Code loads and where each one comes from: a share group, the account's own folder, or the project you are in. Below the table it lists what stops an extension from loading as you expect. It only reads.

## Usage

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

A cell names the source of that item in that account: the share group (`default`), `account` when it lives in the account's own folder, or `-` when the account does not have it. `!` marks an item with a problem, explained below the table.

## What it reads

| Kind | Where | Source shown |
| --- | --- | --- |
| Skills | each folder in `skills/` with a `SKILL.md` | the group `skills/` links to, or `account` |
| Commands, agents | `.md` files under `commands/` and `agents/`, subfolders as `folder/name` | the group, or `account` |
| Plugins | `plugins/installed_plugins.json` | the group, or `account` |
| MCP servers | the group's `mcp.json` (when MCP is shared) and the account's `.claude.json` | the group, or `account` |
| Hooks | `hooks` in `settings.json`, counted per event | the group, or `account` |
| This project | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` and `.mcp.json` in the current folder | listed once, for every account |

## Problems it reports

- **Broken links**: a link in `skills/`, `commands/` or `agents/` whose target is gone, or a shared folder (`skills`, `commands`, `agents`, `plugins`) that links to a missing place (repair it with `ccam doctor --fix`).
- **A skill folder without `SKILL.md`**: Claude Code skips it.
- **An MCP server defined differently in two places**: the group's `mcp.json` and the account's `.claude.json`, or one of them and the project's `.mcp.json`. The same settings in both places (which `mcp_strategy = "sync"` writes) is not a problem. Move an account's own server into the group with `ccam mcp import <account> <server>`.
- **`installed_plugins.json` that is not valid JSON.**

## Safety and edge cases

- `ccam ext` changes nothing and makes no network requests. It reads server names and hook counts, never the values of `env` or headers in MCP definitions.
- It shows what is on disk. A plugin that is installed but turned off in settings, or a server blocked by managed settings, still appears.
- Run it in the project folder to include that project's `.claude` folder and `.mcp.json`.

## Related

- [Sharing between accounts](/features/share)
- [Shared MCP servers](/features/mcp)
- [Checking and repairing](/features/doctor)
