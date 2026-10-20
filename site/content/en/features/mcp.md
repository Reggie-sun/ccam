---
title: Shared MCP servers
description: Keep one list of MCP servers per group with ccam mcp add, add-json, list, remove and import, passed with --mcp-config or synced into each account.
updated: 2026-10-20
order: 5
translation: source
---

Claude Code keeps user MCP servers in `.claude.json`, the same file that holds the account's login, so the file cannot be shared between accounts. ccam keeps shared servers in a group file, `mcp.json`, and hands it to every account of the group at launch. Use `ccam mcp` to add a server once for all your accounts, and `ccam mcp import` to move a server you added inside one session to the others.

## Usage

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

A change applies from the next launch of each account; running sessions keep the servers they started with:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Adding a name the group already has replaces its definition and prints `updated` instead of `added`. Servers reach an account only if it is in that group and the group shares MCP (`ccam group set <group> --add mcp`); when the group does not, `add` says so.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `-g`, `--group <group>` | target group for `add`, `add-json` and `remove` | the only group some account shares MCP through, otherwise `default` |
| `-- <command> [args…]` (`add`) | run this command as a stdio server | |
| `--url <url>` (`add`) | remote server URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport of a `--url` server | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, repeatable | none |
| `-e`, `--env <KEY=VALUE>` (`add`) | environment variable for a stdio server, repeatable | none |
| `--force` (`import`) | replace a group server of the same name that differs; the old `mcp.json` is backed up | off |
| `--keep-local` (`import`) | keep the entries in the account's `.claude.json` | off |
| `-h`, `--help` | help for each command | |
| `--dry-run` (global) | honoured by `mcp list` only; `add`, `add-json`, `remove` and `import` refuse it and change nothing | off |

`add` needs either a command after `--` or `--url`; otherwise it fails with `a command to run (after --) or --url is required`. `--header` applies to `--url` servers and `--env` to command servers. `add-json` takes the same JSON as `claude mcp add-json`, and the definition must be an object with a `command` or a `url`. `import` always uses the group the account shares MCP through, and `list` shows every group, so `-g` does not change them.

## How it works

The files involved:

- `~/.local/share/ccam/groups/<group>/mcp.json`: the shared servers, `{"mcpServers": {…}}`, the format `--mcp-config` reads. It is written with `fsutil.AtomicWrite` and mode `0600`, because server entries often carry API keys.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code's own file with the account's servers. ccam reads only its `mcpServers` value with the `profile.TopLevelValue` scanner and never parses the whole file on the launch path.
- `~/.local/state/ccam/mcp/<name>/`: the `seen` file, the `sync-stamp` and `.claude.json` backups.

There are two strategies, set by `mcp_strategy` in `~/.config/ccam/config.toml`. The default is the flag strategy:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` finds the account's MCP group. `mcp.FlagArgs` checks with one `stat` that the group's `mcp.json` exists, without parsing it.
2. `Plan.InsertArgs` puts `--mcp-config=<path>` right after `claude`, before your own arguments, so your own `--mcp-config` or `--strict-mcp-config` still apply on top.
3. `mcp.NewLocal` checks whether the account's `.claude.json` changed since the last launch, and only then parses its servers.
4. `syscall.Exec` replaces ccam with Claude Code.

The path is passed as one argument, `--mcp-config=<path>`, not as two. Claude Code's `--mcp-config` takes several values, so `claude --mcp-config a.json "prompt"` would read the prompt as a second config file and fail with `MCP config file not found`. The one-argument form stops there, and a prompt or a subcommand such as `mcp list` that follows works. Claude Code loads `--mcp-config` servers in addition to the account's own; on a name clash the entry in `.claude.json` wins. The flag strategy writes `.claude.json` only for Remote Control.

`claude remote-control` (`rc`) takes no `--mcp-config`, and its sessions read servers only from `.claude.json`. For `ccam work rc`, ccam leaves the flag out and, on a terminal, asks whether to copy the group's servers into `.claude.json`, with a backup, as the `sync` strategy does. It asks until a copy exists, then again only after the group's `mcp.json` changes ([Remote Control](/features/launch#remote-control)). Because such a copy wins over `--mcp-config` on a name clash, `ccam mcp add`, `add-json`, `remove` and `import` update the copies that other accounts of the group hold:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

A server added inside a session with `claude mcp add` lands in that account's `.claude.json` only. ccam cannot see it at the time, because Claude Code has replaced it, so it checks at the next launch. The `seen` file holds the modification time of `.claude.json` and the server names already reported. When the time differs, ccam compares the names with the group and prints a hint for each new server once, without waiting for input:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

With `mcp_strategy = "sync"`, the servers appear in `claude mcp list` and `/mcp` as ordinary user servers, and ccam writes them into each account's `.claude.json` before launch:

<!-- motion:mcp-sync -->

1. `mcp.Sync` compares the modification time of the group's `mcp.json` with `sync-stamp`. If it has not changed, the account's `.claude.json` is not opened.
2. It takes Claude Code's own lock, the folder `.claude.json.lock` next to the file, created with `mkdir`. It waits up to 2 seconds and treats a lock older than 10 seconds as stale, as Claude Code's lock library does.
3. It saves the current `.claude.json` to `backups/` in the state folder and keeps the newest 5.
4. It sets every group server in `mcpServers`, removes servers it wrote before that the group no longer has, keeps the account's own servers, and writes the file atomically with its original mode.
5. It records the new stamp and the names it manages, and Claude Code starts without `--mcp-config`.

Sync is skipped, with `ccam: account work is already running, skipping MCP sync`, when a session of the account is running, because Claude Code rewrites `.claude.json` while it runs. A failed sync prints `ccam: MCP sync failed (launching anyway): …` and launches anyway. The hint for servers added inside a session belongs to the flag strategy only.

`ccam mcp import` moves an account's own servers into its group:

<!-- motion:mcp-import -->

ccam reads the account's servers with `mcp.LocalServers` and compares each one with the group's server of the same name using `mcp.SameDef`, which ignores key order and spacing. A server the group lacks or has with the same definition is copied. A different one is skipped, because the group's definition is what every other account uses:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

With `--force` the group's server is replaced, and the previous `mcp.json` is first saved to `.ccam-backup/<time>/mcp.json` inside the group folder. Then `mcp.RemoveLocal` removes the copied entries from the account's `.claude.json`, under the same lock and after a backup in the state folder. This is needed because the account's entry would win over the group's on a name clash. `--keep-local` leaves them.

## Safety and edge cases

- **Sign-ins.** Servers that sign in with OAuth keep their tokens in each account's `.credentials.json` (`mcpOAuth`), not in the server list, so sharing a server does not share its sign-in. Sign in once per account with `/mcp`, or copy the sign-ins from `~/.claude` with [`ccam migrate`](/features/migrate). The group never receives a login file.
- **Running accounts.** `import` copies to the group but leaves a running account's own entries and says so; run it again after the session ends. Sync skips running accounts.
- **Copies for Remote Control.** A running account's copy is not touched: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` With `mcp_strategy = "sync"` nothing needs updating, because every launch syncs.
- **Broken files.** If `.claude.json` is not valid JSON, sync leaves it untouched and reports the error. If the lock cannot be taken in 2 seconds, ccam reports `could not lock … (Claude Code is writing the file)`.
- **Accounts that do not share MCP.** `import` refuses and tells you to put the account in a group that shares MCP first (`ccam group join <group> side`).
- **Group names.** `-g` takes a group name only; it becomes a folder under `groups/`, so names such as `../x` are refused.

## Related

- [Sharing between accounts](/features/share)
- [Using Remote Control with ccam](/guides/remote-control)
- [Share MCP servers, settings and history](/docs/sharing)
- [Copy an existing setup](/features/migrate)
- [All commands](/docs/commands)
