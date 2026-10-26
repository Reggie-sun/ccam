---
title: Shared MCP servers
description: Keep one list of MCP servers per group with ccam mcp add, add-json, list, remove and import, passed with --mcp-config or synced into each account.
updated: 2026-10-26
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

Below the groups, `ccam mcp list` names every server that does not reach the other accounts:

- `only in account work: pg` for a server an account added for itself;
- `only in account side, for /home/me/code/app only: notes` for one added for one folder (plain `claude mcp add`, see below);
- `account work has its own …, used in its sessions in place of the group's` when an account's own server has a group server's name;
- this folder's `.mcp.json`, with the accounts whose approval it is `waiting for approval in`;
- servers that are `in no account, only in ~/.claude.json (plain claude's)`, with the [`ccam migrate --mcp`](/features/migrate) command that copies them.

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
| `--dry-run` (global) | supported by `mcp list` only; `add`, `add-json`, `remove` and `import` refuse it and change nothing | off |

`add` needs either a command after `--` or `--url`; otherwise it fails with `a command to run (after --) or --url is required`. `--header` applies to `--url` servers and `--env` to command servers. `add-json` takes the same JSON as `claude mcp add-json`, and the definition must be an object with a `command` or a `url`. `import` always uses the group the account shares MCP through, and `list` shows every group, so `-g` does not change them.

## How it works

The files involved:

- `~/.local/share/ccam/groups/<group>/mcp.json`: the shared servers, `{"mcpServers": {…}}`, the format `--mcp-config` reads. It is written with `fsutil.AtomicWrite` and mode `0600`, because server entries often carry API keys, and only under the lock `<group>/.mcp.json.lock`, so two `ccam mcp add` commands run at the same time both keep their servers.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code's own file with the account's servers. ccam reads only its `mcpServers` value with the `profile.TopLevelValue` scanner and never parses the whole file on the launch path.
- `~/.local/state/ccam/mcp/<name>/`: the `seen` file, the `sync-stamp` (with a hash of each copy ccam wrote, so it can tell a copy nobody changed from one you edited), `sync-checked` (what the sync strategy saw last), `kept` (the servers the account kept when it left MCP sharing) and `.claude.json` backups.

There are two strategies, set by `mcp_strategy` in `~/.config/ccam/config.toml`. The default is the flag strategy:

<!-- motion:mcp-flag -->

1. ccam finds the account's MCP group (`launcher.applyMCP`) and checks with one `stat` that the group's `mcp.json` exists, without parsing it (`mcp.FlagArgs`).
2. ccam puts `--mcp-config=<path>` right after `claude`, before your own arguments, so your own `--mcp-config` or `--strict-mcp-config` still apply on top (`Plan.InsertArgs`).
3. ccam checks whether the account's `.claude.json` changed since the last launch, and only then parses its servers (`mcp.NewLocal`).
4. ccam replaces itself with Claude Code (`syscall.Exec`).

The path is passed as one argument, `--mcp-config=<path>`, not as two. Claude Code's `--mcp-config` takes several values, so `claude --mcp-config a.json "prompt"` would read the prompt as a second config file and fail with `MCP config file not found`. The one-argument form avoids this, so a prompt or subcommand after it works. `ccam work mcp list`, `get` and `remove` get no flag, since Claude Code refuses it there; ccam names the group's servers on standard error instead ([details](/features/launch#shared-mcp-servers)). Claude Code loads `--mcp-config` servers in addition to the account's own; on a name clash the entry in `.claude.json` wins. The flag strategy writes `.claude.json` only for Remote Control, when an account leaves MCP sharing or joins it again ([below](#leaving-and-joining)), and for `ccam mcp import`.

`claude remote-control` (`rc`) takes no `--mcp-config`, and its sessions read servers only from `.claude.json`. For `ccam work rc`, ccam leaves the flag out and, on a terminal, asks whether to copy the group's servers into `.claude.json`, with a backup, as the `sync` strategy does. It asks until a copy exists, then again only after the group's `mcp.json` changes or a copied server is gone ([Remote Control](/features/launch#remote-control)). Because such a copy wins over `--mcp-config` on a name clash, `ccam mcp add`, `add-json`, `remove` and `import` update the copies that other accounts of the group hold:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

The `✓ Updated` line appears only when a copy changed. A copy that an older version of ccam created is brought up to date the same way, and so is one after `ccam migrate`.

A server added inside a session with `claude mcp add` lands in that account's `.claude.json` only. Plain `claude mcp add` uses the local scope: the server loads for that one account in that one folder (`-s user` makes it available to the account in every folder). ccam cannot see this when it happens, because Claude Code has already replaced the ccam process, so it checks at the next launch. The `seen` file holds the modification time of `.claude.json` and the server names already reported. When the time differs, ccam compares the names with the group and prints a hint for each new server once, without waiting for input, with both strategies:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

The hint also names an account's own server that has a group server's name, since the account's wins in its sessions.

With `mcp_strategy = "sync"`, the servers appear in `claude mcp list` and `/mcp` as ordinary user servers, and ccam writes them into each account's `.claude.json` before launch:

<!-- motion:mcp-sync -->

1. ccam compares the modification time of the group's `mcp.json` with `sync-stamp`, and the account's `.claude.json` with what it saw last (`mcp.Sync`). If neither changed, nothing else happens. Otherwise it also checks that every group server is still in `.claude.json`, and puts back a copy that is gone (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. It takes Claude Code's own lock, the folder `.claude.json.lock` next to the file, created with `mkdir`. It waits up to 2 seconds and treats a lock older than 10 seconds as stale, as Claude Code's lock library does.
3. It saves the current `.claude.json` to `backups/` in the state folder and keeps the newest 5.
4. It sets every group server in `mcpServers`, removes servers it wrote before that the group no longer has, and writes the file atomically with its original mode, only when something changed. It changes or removes only the copies it wrote that nobody changed since. If the account has its own server with the same name as a group server, ccam keeps the account's server and says so once: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. It records the new stamp with a hash of each copy, and Claude Code starts without `--mcp-config`.

Sync is skipped, with `ccam: account work is already running, skipping MCP sync`, when a session of the account is running, however Claude Code was started, because Claude Code rewrites `.claude.json` while it runs. A failed sync prints `ccam: MCP sync failed (launching anyway): …` and launches anyway.

### Leaving and joining

An account that stops sharing MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) keeps the group's servers: the plan copies them into its `.claude.json`, with a backup, and its own servers of the same names stay. With `--fresh` it gets none, and the plan lists the servers that stop loading. When it joins a group that has those servers again, the kept copies that nobody changed are removed, so the group's later edits reach it.

`ccam mcp import` moves an account's own servers into its group:

<!-- motion:mcp-import -->

ccam reads the account's servers, its own and those it has for one folder only, and compares each one with the group's server of the same name using `mcp.SameDef`, which ignores key order and spacing. A folder server becomes a group server for every folder (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); one name with two different folder definitions is refused. Plugin servers are not offered, since the plugin brings them. A server the group lacks or has with the same definition is copied. A different one is skipped, because the group's definition is what every other account uses:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

With `--force` the group's server is replaced, and the previous `mcp.json` is first saved to `.ccam-backup/<time>/mcp.json` inside the group folder. Then ccam removes the copied entries from the account's `.claude.json`, its own and the folder's in one write, under the same lock and after a backup in the state folder (`mcp.Remove`). ccam keeps the group's `mcp.json` locked from when it reads the file until those entries are removed, so nobody removes the server from the group in between. This is needed because the account's entry would win over the group's on a name clash. With `--keep-local`, the entries stay in the account's `.claude.json`.

## Safety and edge cases

- **Sign-ins.** Servers that sign in with OAuth keep their tokens in each account's `.credentials.json` (`mcpOAuth`), not in the server list, so sharing a server does not share its sign-in. Sign in once per account with `/mcp`. ccam never copies sign-ins, and neither does [`ccam migrate`](/features/migrate): a copied refresh token logs out every other holder. [`ccam ext`](/features/ext) and `ccam doctor` name the servers an account is not signed in to while another account is, and `ccam next --continue` says which ones the new account lacks. The group never receives a login file.
- **Running accounts.** `import` copies to the group but leaves a running account's own entries and says so; run it again after the session ends. Sync skips running accounts. Both also detect Claude Code sessions started without ccam.
- **Copies for Remote Control.** A running account's copy is not touched: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` With `mcp_strategy = "sync"` nothing needs updating, because every launch syncs.
- **Broken files.** If `.claude.json` is not valid JSON, sync leaves it untouched and reports the error. If the lock cannot be taken in 2 seconds, ccam reports `could not lock … (Claude Code is writing the file)`.
- **Accounts that do not share MCP.** `import` refuses and tells you to put the account in a group that shares MCP first (`ccam group join <group> side`).
- **Group names.** `-g` takes a group name only; it becomes a folder under `groups/`, so names such as `../x` are refused.

## Related

- [Sharing with groups](/features/share)
- [Using Remote Control with ccam](/guides/remote-control)
- [Share MCP servers, settings and history](/docs/sharing)
- [Copying an existing setup](/features/migrate)
- [Command reference](/docs/commands)
