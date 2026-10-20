---
title: Share MCP servers, settings and history
description: How ccam groups share MCP servers, settings.json, CLAUDE.md, skills and conversation history between Claude Code accounts, with every login kept separate.
updated: 2026-10-20
published: 2026-10-19
order: 4
translation: source
---

Each account keeps its own login. Everything else can be shared through **groups**: a group names the items it shares and the accounts in it, and ccam links each of those accounts' files to the group's folder, so a change made in one account appears in all of them.

## What can be shared

| Item | What is shared |
| --- | --- |
| `mcp` | MCP servers: the group's `mcp.json`, passed to Claude Code at start |
| `settings` | `settings.json`: permissions, hooks, model, status line |
| `memory` | the global `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | conversation history: `projects/` and `history.jsonl` |

Your first account creates the group `default`, which shares everything except conversation history, and later accounts join it.

The login files `.credentials.json` and `.claude.json` are never shared. ccam only links paths from a fixed list, and these two are not on it.

This page covers the everyday tasks. Every option is described on the feature pages: [sharing with groups](/features/share), [shared MCP servers](/features/mcp) and [copying an existing setup](/features/migrate).

## Groups

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

An account can be in several groups, but each item comes from one group only. So to share history between `work` and `side` but not your other accounts, keep history out of `default` and give the two a group of their own, as above. `ccam share work` shows what an account shares and through which group. The [dashboard](/features/dashboard) does the same on its Groups tab.

## When an account starts sharing

When an account starts sharing an item (it joins a group, or its group adds the item), ccam:

1. prints a plan of every file it will move, copy or link for each account, and asks before it starts;
2. moves the account's files into the group folder if the group is still empty, or adds the files the group does not have yet;
3. when the same file differs, asks whether to keep the group's version, use yours, or merge the JSON;
4. keeps your original in `.ccam-backup/<time>/` inside the account folder;
5. replaces the original with a link to the group.

Conversation files are named after their session ID, so histories from several accounts combine without conflicts.

ccam does not change an account while Claude Code is running as it, because the session keeps writing into the files being moved. Quit it first.

## When an account stops sharing

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

The link is removed and the account gets its own copy of the group's current content. Add `--fresh` to start empty instead. The group keeps its content for the other accounts.

## Copy an existing setup

If you used Claude Code before ccam, your settings are in `~/.claude`. `ccam migrate` copies them into a group, so that every account sharing those items gets them:

```sh
ccam migrate
```

On a terminal, ccam first shows a checklist of the items. Everything except conversation history is checked; press Space to change an item and Enter to see the plan. Items the folder has nothing for are marked `not found`.

What it copies:

- `settings.json` and the global `CLAUDE.md`;
- the extensions `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` and `rules/`;
- your user MCP servers, the `mcpServers` entry of `~/.claude.json`, into the group's `mcp.json`;
- the saved sign-ins (OAuth) of the MCP servers it copies, into each account's own login file;
- conversation history (`projects/` and `history.jsonl`), but only when you add `--history`.

The source folder is only read: nothing in it is moved, linked or changed, and plain `claude` keeps working as before. The login files `.credentials.json` and `.claude.json` are never copied into the group.

MCP servers that sign in with OAuth keep their tokens in the login, not in the server list. So that you do not have to sign in to them again, ccam copies those sign-ins into the login file of every account that gets the servers through the group, unless the account already has its own. It skips accounts that are running, not logged in yet, registered in place with `ccam import`, or keep their login in the macOS Keychain, and names them; sign in to the servers there with `/mcp`. An account that starts sharing MCP servers later gets the sign-ins when you run `ccam migrate` again.

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
Group default: personal
To share them:
  ccam group join default work
```

At the end, ccam lists the group's accounts. If the group does not share one of the items yet, or some accounts are not in it, it prints the `ccam group` commands that change that.

When the group already has a different copy of a file, ccam asks, as it does when you turn sharing on: keep the group's copy or use yours. `settings.json` and other JSON files can also be merged, and `history.jsonl` line by line. The group's previous copies are saved in `.ccam-backup/<time>/` inside the group folder.

| Option | Effect |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copy only the items you name, without the checklist; without them and outside a terminal, everything except conversation history |
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

If the group already has a different server with that name, ccam skips it and says so. `--force` replaces the group's server and keeps the old `mcp.json` in the group's `.ccam-backup/` folder.

MCP servers that sign in with OAuth store their tokens with the account's login, so authorise them once in each account, or copy the sign-ins from `~/.claude` with `ccam migrate`.

If you prefer the servers to appear in `claude mcp list` and `/mcp` as ordinary user servers, set `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam then writes them into each account's `.claude.json` before starting, under Claude Code's own file lock and with a backup.

## Two accounts in the same folder

Accounts that share conversation history write to the same files. If you start a second Claude Code in a folder where another account of the same group is already running, ccam warns you and asks whether to continue.

## When something breaks

```sh
ccam doctor
ccam doctor --fix
```

If a shared path became an ordinary file, for example after you restored a backup by hand, ccam stops before starting Claude Code and points to `ccam doctor --fix`, which merges the file into the group and keeps a backup.
