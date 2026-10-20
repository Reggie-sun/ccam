---
title: Command reference
metaTitle: ccam command reference
description: Every ccam command and option, from ccam add and ccam <name> to ccam usage, ccam next and ccam doctor.
updated: 2026-10-20
published: 2026-10-19
order: 3
translation: source
---

Put `--dry-run` before a command to see what it would do without doing it, for example `ccam --dry-run work`. It works when starting Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) and with `group`, `migrate`, `shortcut`, `update` and `doctor --fix`. Other commands that change something refuse it instead of running.

<!-- commands -->

## Starting Claude Code

`ccam <name|number> [arguments]` starts Claude Code as that account. ccam does not read the arguments after the name; it passes them to `claude` in the same order, and standard input, the exit code and signals behave exactly as with `claude` itself. A `--` after the name is optional.

`ccam` with arguments but no account name, such as `ccam -p "summarise"`, uses the account set for the current folder, else the default account. Set them with `ccam folder <name> [folder]` and `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` before the account name is for `-p`/`--print` runs. When the account reaches its usage limit, ccam runs the same command again as the logged-in account with the most room and prints only the output of the run that finished. See [Accounts by folder](/features/folder) and [Run again at the limit](/features/next#run-again-at-the-limit).

A shortcut such as `ccwork` starts its account the same way: `ccwork --resume` does what `ccam work --resume` does. See [shortcuts](/docs/commands#shortcuts).

## Adding and removing accounts

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` without a folder looks for `~/.claude` and `~/.claude-*` and asks about each one. Imported folders are registered where they are; ccam never moves or copies a login.

`ccam add --api` adds an account that uses an API key instead of a Claude login; ccam reads the key at each launch and never stores it. See [API key accounts](/features/api).

`ccam rename`, `ccam remove` and the `ccam group` commands refuse to change an account while Claude Code is running as that account. Quit the session first.

`ccam list --json` prints every account with its slot, shortcut, email and folder. Its `groups` array lists the account's groups, and its `share` object maps each item (`mcp`, `settings`, `memory`, `extensions`, `history`) to the group it comes from, or `local`.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

A group shares items between its accounts; an account can be in several groups, but each item comes from one group only. See [sharing with groups](/features/share).

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
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` and `ccam pick` choose the same way: below `auto_max_usage` (`config.toml`, default 100) first, by `auto_strategy` (`headroom`, `expiring` or `balance`), counting the `--model` weekly window, never an account at 100%, and within the organization of the account they start from (`--cross-org` or `cross_org = true` to look further). When none can be chosen they exit with code 75, or with `--wait` wait for the first account to reset. `ccam next --continue` carries a conversation over even between accounts that do not share history: Claude Code then copies it into the new account. Inside Claude Code, the plugin's `/ccam:switch [account]` hands a conversation on, and `ccam hook` is what its hooks run. See [Switching to the next account](/features/next), [Picking an account for scripts](/features/pick) and [Background sessions of every account](/features/agents).

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` checks the release's SHA-256 checksum before it replaces ccam, and `--version v0.2.0` installs a particular release. A ccam installed from the `.deb` package is updated with apt instead. See [update](/docs/install#update).

## Files

| Path | Contents |
| --- | --- |
| `~/.config/ccam/config.toml` | Accounts, groups and settings (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | Each account's `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | Shared settings, MCP servers, extensions and history |
| `~/.local/state/ccam/` | Usage cache, last-used times and the notes `/ccam:switch` and the limit hook leave (`handoff/`) |
| `~/.local/bin/<shortcut>` | Shortcut commands such as `ccwork`, links to ccam |

The `XDG_CONFIG_HOME`, `XDG_DATA_HOME` and `XDG_STATE_HOME` variables move these folders. On Windows, `~` is `%USERPROFILE%` and shortcuts are `<shortcut>.exe`.
