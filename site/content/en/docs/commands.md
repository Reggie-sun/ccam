---
title: Command reference
description: Every ccam command and option, from ccam add and ccam <name> to ccam usage, ccam next and ccam doctor.
updated: 2026-10-19
order: 3
translation: source
---

Put `--dry-run` before any command to see what it would do without doing it, for example `ccam --dry-run work`.

<!-- commands -->

## Starting Claude Code

`ccam <name|number> [arguments]` starts Claude Code as that account. ccam does not read the arguments after the name; it passes them to `claude` in the same order, and standard input, the exit code and signals behave exactly as with `claude` itself. A `--` after the name is optional.

`ccam` with arguments but no account name, such as `ccam -p "summarise"`, uses the default account. Set it with `ccam use <name>`.

The shortcuts `claude1`, `claude2` and so on are links to ccam in `~/.local/bin`. `ccam shims install --names` also creates `claude-<name>` for every account.

## Adding and removing accounts

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` without a folder looks for `~/.claude` and `~/.claude-*` and asks about each one. Imported folders are registered where they are; ccam never moves or copies a login.

## Sharing

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Each item takes `on` (the `default` group), `off`, or a group name. See [sharing](/docs/sharing).

## MCP servers

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Usage and switching

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Files

| Path | Contents |
| --- | --- |
| `~/.config/ccam/config.toml` | Accounts and what they share |
| `~/.local/share/ccam/profiles/<name>/` | Each account's `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | Shared settings, MCP servers, extensions and history |
| `~/.local/state/ccam/` | Usage cache and last-used times |
| `~/.local/bin/claude1`, `claude2`… | Links to ccam |

The `XDG_CONFIG_HOME`, `XDG_DATA_HOME` and `XDG_STATE_HOME` variables move these folders.
