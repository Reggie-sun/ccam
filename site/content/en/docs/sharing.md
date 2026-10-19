---
title: Share MCP servers, settings and history
description: How ccam shares MCP servers, settings.json, CLAUDE.md, agents, skills and conversation history between Claude Code accounts while keeping every login separate.
updated: 2026-10-19
order: 4
translation: source
---

Each account keeps its own login. Everything else can be shared item by item: ccam links the account's file or folder to a shared group folder, so a change made in one account appears in all accounts of that group.

## What can be shared

| Item | Option | What is shared |
| --- | --- | --- |
| MCP servers | `--mcp` | The group's `mcp.json`, passed to Claude Code at start |
| Settings | `--settings` | `settings.json`: permissions, hooks, model, status line |
| Global memory | `--memory` | `CLAUDE.md` |
| Extensions | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Conversation history | `--history` | `projects/` and `history.jsonl` |

New accounts share everything except conversation history.

The login files `.credentials.json` and `.claude.json` are never shared. ccam only links paths from a fixed list, and these two are not on it.

## Groups

`on` means the `default` group. You can also give a name, so that only some accounts share an item:

```sh
ccam share work --history team
ccam share side --history team
```

`work` and `side` now share conversations with each other, but not with your other accounts.

## Turning sharing on

When an account starts sharing an item, ccam:

1. prints a plan of every file it will move, copy or link, and asks before it starts;
2. moves the account's files into the group folder if the group is still empty, or adds the files the group does not have yet;
3. when the same file differs, asks whether to keep the group's version, use yours, or merge the JSON;
4. keeps your original in `.ccam-backup/<time>/` inside the account folder;
5. replaces the original with a link to the group.

Conversation files are named after their session ID, so histories from several accounts combine without conflicts.

## Turning sharing off

```sh
ccam share work --history off
```

The link is removed and the account gets its own copy of the group's current content. Add `--fresh` to start empty instead. The group keeps its content for the other accounts.

## Copy an existing setup

If you used Claude Code before ccam, your settings are in `~/.claude`. `ccam migrate` copies them into a group, so that every account sharing those items gets them:

```sh
ccam migrate
```

What it copies:

- `settings.json` and the global `CLAUDE.md`;
- the extensions `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` and `rules/`;
- your user MCP servers, the `mcpServers` entry of `~/.claude.json`, into the group's `mcp.json`;
- conversation history (`projects/` and `history.jsonl`), but only when you add `--history`.

The source folder is only read: nothing in it is moved, linked or changed, and plain `claude` keeps working as before. The login files `.credentials.json` and `.claude.json` are never copied.

ccam prints the plan and asks before it copies anything:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

At the end, ccam lists the accounts that share each item through the group. If some accounts do not, it prints the `ccam share` commands that turn sharing on for them.

When the group already has a different copy of a file, ccam asks, as it does when you turn sharing on: keep the group's copy or use yours. `settings.json` and other JSON files can also be merged, and `history.jsonl` line by line. The group's previous copies are saved in `.ccam-backup/<time>/` inside the group folder.

| Option | Effect |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copy only the items you name; without them, everything except conversation history |
| `-g`, `--group <group>` | Copy into this group instead of `default` |
| `-y`, `--yes` | Copy without asking; when a file differs, keep the group's copy |

Put `--dry-run` before the command, as in `ccam --dry-run migrate`, to see the plan without copying anything. To copy from another folder, name it: `ccam migrate ~/.claude-work --group work`.

## MCP servers

Claude Code stores MCP servers in `.claude.json`, together with the login, so that file cannot be shared. ccam keeps shared servers in the group's `mcp.json` and starts Claude Code with `--mcp-config=<group>/mcp.json`. Your own `--mcp-config` options still apply on top.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

A server you add inside a session with `claude mcp add` is saved for that account only. The next time you start that account, ccam names the server and the command that moves it to the group:

```sh
ccam mcp import work gh
```

MCP servers that sign in with OAuth store their tokens with the account's login, so authorise them once in each account.

If you prefer the servers to appear in `claude mcp list` and `/mcp` as ordinary user servers, set `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam then writes them into each account's `.claude.json` before starting, under Claude Code's own file lock and with a backup.

## Two accounts in the same folder

Accounts that share conversation history write to the same files. If you start a second Claude Code in a folder where another account of the same group is already running, ccam warns you and asks whether to continue.

## When something breaks

```sh
ccam doctor
ccam doctor --fix
```

If a shared path became an ordinary file, for example after you restored a backup by hand, ccam stops before starting Claude Code and points to `ccam doctor --fix`, which merges the file into the group and keeps a backup.
