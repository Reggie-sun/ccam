---
title: Command reference
description: Every ccam command and option, from ccam add and ccam <name> to ccam usage, ccam next and ccam doctor.
updated: 2026-10-19
order: 3
translation: source
---

Put `--dry-run` before a command to see what it would do without doing it, for example `ccam --dry-run work`. It works when starting Claude Code (`ccam <name>`, `ccam run`, `ccam next`) and with `share`, `migrate`, `shortcut`, `update` and `doctor --fix`. Other commands that change something refuse it instead of running.

<!-- commands -->

## Starting Claude Code

`ccam <name|number> [arguments]` starts Claude Code as that account. ccam does not read the arguments after the name; it passes them to `claude` in the same order, and standard input, the exit code and signals behave exactly as with `claude` itself. A `--` after the name is optional.

`ccam` with arguments but no account name, such as `ccam -p "summarise"`, uses the default account. Set it with `ccam use <name>`.

A shortcut such as `ccwork` starts its account the same way: `ccwork --resume` does what `ccam work --resume` does. See [shortcuts](/docs/commands#shortcuts).

## Adding and removing accounts

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` without a folder looks for `~/.claude` and `~/.claude-*` and asks about each one. Imported folders are registered where they are; ccam never moves or copies a login.

`ccam rename`, `ccam remove` and `ccam share` refuse to change an account while Claude Code is running as that account. Quit the session first.

`ccam list --json` prints every account with its slot, shortcut, email and folder. Its `share` object maps each item (`mcp`, `settings`, `memory`, `extensions`, `history`) to a group name or `local`.

## Shortcuts

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

A shortcut is a command you name yourself that starts Claude Code as one account. Like `ccam work`, it passes every argument to Claude Code unchanged, so `ccwork --resume` and `ccwork -p "hi"` work as you would expect. It is a link to ccam in `~/.local/bin`.

Shortcuts are optional. The `ccam add` wizard asks whether to create one, `ccam add <name> --shortcut <command>` creates one without asking, and `ccam list` shows each account's shortcut, or `-` when it has none.

Names use the characters a-z, 0-9, `.`, `-` and `_`. A name cannot be a ccam command, another account's shortcut, or the name of another program on your `PATH`.

`ccam rename` keeps the shortcut name, and `ccam remove` deletes the shortcut together with the account. `ccam shims install` recreates the links of all accounts, for example after you have moved ccam.

## Sharing

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Each item takes `on` (the `default` group), `off`, or a group name. See [sharing](/docs/sharing).

## Copying an existing setup

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copies `settings.json`, `CLAUDE.md`, extensions and user MCP servers from a Claude Code folder into a group, so every account that shares those items gets them. The folder is only read. The options `--settings`, `--memory`, `--extensions`, `--mcp` and `--history` copy only the items you name. Without them, a checklist on a terminal lets you choose, with everything except conversation history checked; elsewhere everything except conversation history is copied. `-g` or `--group` picks the group, and `-y` or `--yes` copies without asking and keeps the group's copy when a file differs. It also copies the servers' saved sign-ins into each account's own login file. See [copying an existing setup](/docs/sharing#copy-an-existing-setup).

## MCP servers

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` moves the servers that only that account has into its group. A server whose name the group already uses for a different definition is skipped; add `--force` to replace the group's copy, which ccam backs up first.

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
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

`ccam update` checks the release's SHA-256 checksum before it replaces ccam, and `--version v0.2.0` installs a particular release. A ccam installed from the `.deb` package is updated with apt instead. See [update](/docs/install#update).

## Files

| Path | Contents |
| --- | --- |
| `~/.config/ccam/config.toml` | Accounts and what they share |
| `~/.local/share/ccam/profiles/<name>/` | Each account's `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | Shared settings, MCP servers, extensions and history |
| `~/.local/state/ccam/` | Usage cache and last-used times |
| `~/.local/bin/<shortcut>` | Shortcut commands such as `ccwork`, links to ccam |

The `XDG_CONFIG_HOME`, `XDG_DATA_HOME` and `XDG_STATE_HOME` variables move these folders.
