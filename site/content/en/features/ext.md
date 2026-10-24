---
title: Inspecting extensions
description: ccam ext shows the skills, commands, plugins, MCP servers and hooks each account loads, where each comes from, and what is broken.
updated: 2026-10-23
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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

A cell names the source of that item in that account: the share group (`default`), `account` when it lives in the account's own folder, `local` for an MCP server the account has for this folder only, or `-` when the account does not have it. A plugin that is installed but not turned on shows `(off)` after its source. `!` marks an item with a problem, explained below the table; problems that belong to no account, such as servers only plain `claude` has, come last without a name. `ccam ext --json` prints the same report (`accounts`, the project, `problems`), with `remote` on MCP servers that Claude Code reaches over the network; with no accounts it prints the report with `accounts: []`.

## What it reads

| Kind | Where | Source shown |
| --- | --- | --- |
| Skills | each folder in `skills/` with a `SKILL.md` | the group `skills/` links to, or `account` |
| Commands, agents | `.md` files under `commands/` and `agents/`, subfolders as `folder/name` | the group, or `account` |
| Plugins | `plugins/installed_plugins.json`, turned on or off by `enabledPlugins` in the account's and the project's settings; a plugin installed for one project only appears in that project | the group, or `account`; `(off)` when not turned on |
| MCP servers | the group's `mcp.json` (when MCP is shared), the account's `.claude.json` (its own servers, and those it has for this folder only), and the servers of its plugins, as `plugin:<plugin>:<server>` | the group, `account` or `local` |
| Hooks | `hooks` in `settings.json`, counted per event | the group, or `account` |
| This project | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` and `.mcp.json` in the current folder | listed once; an account that has not trusted the folder gets a problem for its `.mcp.json` servers |

## Problems it reports

- **Broken links**: a link in `skills/`, `commands/` or `agents/` whose target is gone, or a shared folder (`skills`, `commands`, `agents`, `plugins`) that links to a missing place (repair it with `ccam doctor --fix`, which makes a missing group folder again).
- **A skill folder without `SKILL.md`**: Claude Code skips it.
- **An MCP server defined differently in two places**: the group's `mcp.json` and the account's `.claude.json`, or one of them and the project's `.mcp.json`. The same settings in both places (which `mcp_strategy = "sync"` writes) is not a problem. Move an account's own server into the group with `ccam mcp import <account> <server>`.
- **`installed_plugins.json` that is not valid JSON**, a plugin that `enabledPlugins` does not name, and a plugin whose folder is missing.
- **Servers only one account loads**: a server only in this account's `.claude.json`, which the other accounts of its group do not load, and one in the account's local scope for this folder (what `claude mcp add` does unless given `-s user`). Each comes with the `ccam mcp import` command that shares it. Servers in plain `claude`'s `~/.claude.json` that no account loads come with the `ccam migrate --mcp` command.
- **A project server waiting for approval**: the account has not trusted the folder yet, so Claude Code asks once, and until then the project's `.mcp.json` servers wait. The problem says where to approve them: `"enabledMcpjsonServers"` in the `settings.json` of the accounts of the settings group, or of this account. Trust given to a git repository's main worktree counts in its other worktrees.
- **Sign-ins and connectors**: an MCP server that needs a sign-in in this account while another account is signed in to it (a sign-in holds for one account only; sign in with `/mcp`), and claude.ai connectors another account has used and this one has not.

## Safety and edge cases

- `ccam ext` changes nothing and makes no network requests. It reads server names and hook counts, never the values of `env` or headers in MCP definitions. For the sign-in check it reads the names of the MCP sign-ins in each account's login (on macOS from the Keychain), never the tokens, and only for accounts that load a server another account loads too.
- It shows what is on disk. A server blocked by managed settings still appears; a plugin turned off in settings shows `(off)`.
- Run it in the project folder to include that project's `.claude` folder and `.mcp.json`.

## Related

- [Sharing between accounts](/features/share)
- [Shared MCP servers](/features/mcp)
- [Checking and repairing](/features/doctor)
