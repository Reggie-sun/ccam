---
title: Background sessions of every account
description: ccam agents lists the Claude Code background agents and open sessions of all your accounts in one table, and ccam attach opens one by its id.
updated: 2026-10-20
order: 16
translation: source
---

Claude Code keeps background sessions (`claude --bg`, agent view) per account: `claude agents` only shows the sessions of the account it runs as. `ccam agents` asks every account at once and lists them together, and `ccam attach <id>` opens a background session as whichever account started it.

## Usage

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` is the background session's state (`working`, `blocked`, `done`, `failed`, `stopped`) and what its process is doing (`busy`, `waiting`, `idle`); an open interactive session shows `open`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--all` | Also list background sessions that are done, failed or stopped | off |
| `--json` | Print the sessions as JSON, each with the `account` it belongs to | off |

`ccam attach <id> [claude args]` passes the arguments after the id to `claude attach`. With `--dry-run` it prints the command.

## How it works

1. For every account, ccam runs `claude agents --json` (with `--all` when given) with the account's `CLAUDE_CONFIG_DIR`, all at once, each within 15 seconds. This reads the account's own session list and does not start Claude Code's background service.
2. Each entry gets the account's name, and the list is sorted by start time, newest first. An account that does not answer is named on standard error and left out, and ccam then exits with code 1, so a script does not take a failure for an empty list.
3. `ccam attach` looks for `jobs/<id>` in each account's folder, where Claude Code keeps a background session, and runs `claude attach <id>` as that account, the same way `ccam <name>` starts Claude Code.

## Safety and edge cases

- Background session ids belong to one account: `claude attach` run as another account cannot find them, which is what `ccam attach` sorts out.
- Two accounts could hold sessions whose ids start the same; `ccam attach` then names both and asks for more of the id.
- ccam passes its shared MCP servers to Claude Code with `--mcp-config=<file>` in one argument, which `claude agents` and `claude attach` accept; background sessions keep using that file when they restart.
- An API key account's background service takes its key from the shell that started it. Start the first background session of such an account through ccam (`ccam gw --bg …`), so the service has the key.

## Related

- [Launching Claude Code](/features/launch)
- [Picking an account for scripts](/features/pick)
- [The dashboard](/features/dashboard)
