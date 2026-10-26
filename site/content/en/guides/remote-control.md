---
title: Using Remote Control with ccam
description: Start Claude Code Remote Control as any account with ccam <name> rc, see which options apply, why ccam offers to copy shared MCP servers, and fix problems.
updated: 2026-10-26
order: 4
translation: source
---

`ccam work rc` starts Claude Code's Remote Control as the `work` account, so you can control sessions on this computer from claude.ai/code or the Claude mobile app. It works like `claude rc` with one difference: if the account shares MCP servers, ccam asks once whether to copy them into the account's `.claude.json`, because that file is the only place Remote Control sessions read servers from.

## Start it

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

A [shortcut](/features/shortcuts) works the same way: `ccwork rc` is `ccam work rc`. Run the command in the folder you want to work in, then open claude.ai/code or the Claude app, signed in with the same Claude account as `work`.

## Options

Two kinds of options meet here: ccam's own, which go before the account name or live in ccam's settings, and Remote Control's, which go after `rc` and reach Claude Code unchanged.

### ccam options and settings

| Option or setting | What it does | Default |
| --- | --- | --- |
| `--dry-run` (before the account) | Print the environment and command, and the MCP notice if there is one. Neither asks nor copies | off |
| `mcp_strategy = "sync"` at the top of `~/.config/ccam/config.toml` | Copy the group's servers into `.claude.json` before every launch, Remote Control or not, so nothing is ever asked | `"flag"` |
| `ccam group leave <group> work` | Take the account out of the group that shares MCP with it. The group's servers are copied into its `.claude.json` as its own, so Remote Control has them and nothing is asked; with `--fresh` it starts with only the account's own servers | in `default` |

When the account shares MCP servers and its `.claude.json` lacks the group's current list, or a copied server is gone from it, ccam asks `Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]`. `y` copies them; `N` or Enter starts without them.

### Remote Control options

Everything after `rc` belongs to Claude Code. These are the options of `claude remote-control` in Claude Code 2.1.304; `ccam work rc --help` prints the list for your version.

| Option | What it does |
| --- | --- |
| `--name <name>` | Name of the session shown in claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefix for generated session names (default: the host name) |
| `-c`, `--continue` | Reattach to the session Remote Control last recorded for this folder, within about 4 hours |
| `--session-id <id>` | Reattach to a specific session |
| `--permission-mode <mode>` | Permission mode for new sessions: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` or `plan` |
| `--spawn <mode>` | `same-dir` (default), `worktree` (a git worktree per session) or `session` (one session, exits when it ends) |
| `--capacity <N>` | Maximum number of concurrent sessions in `same-dir` and `worktree` mode (default 32) |
| `--[no-]create-session-in-dir` | Create a first session in the current folder at start (default on) |
| `--[no-]chrome` | Claude in Chrome for new sessions (default off) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Debug and verbose output |

ccam reads Claude Code's command the way Claude Code does, skipping options together with their values, so `ccam work rc --name laptop` and `ccam work --model opus rc` both start Remote Control. An option ccam does not know counts as one without a value; to be safe, put `rc` straight after the account name.

## What ccam does

The launch takes the usual path ([Launching Claude Code](/features/launch)) with one change in the MCP step:

1. ccam finds the account, checks its share links and resolves `claude`, as for any launch.
2. It sees that the command is `rc` or `remote-control` and leaves `--mcp-config` out.
3. If the account shares MCP servers, it compares the group's `mcp.json` with the `sync-stamp` in `~/.local/state/ccam/mcp/work/`, and checks that every group server is still in `.claude.json`. When `.claude.json` already holds the current servers, or the group has none, it goes straight to step 6.
4. Otherwise it names the servers the sessions would miss:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

5. On `y` it saves a backup of `.claude.json`, takes Claude Code's own lock, adds the group's servers to `mcpServers` and keeps the account's other servers; if the account has its own server with the same name as a group server, ccam keeps the account's server and says so. On `N`, without a terminal, or while the account is running, it changes nothing and says so:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. It records the session, sets `CLAUDE_CONFIG_DIR` to the account folder and replaces itself with `claude rc` through `exec`.

Remote Control then keeps running as a server. Every session it starts for you is a child of that process, so it runs as `work`: with `work`'s login, settings and conversation history, and counted against `work`'s usage limit. ccam treats the account as running while a Remote Control server it started is running. A server started with plain `claude remote-control` counts only while it has sessions; ccam cannot see one that has none yet.

The dry run shows the difference. A normal launch carries the group's servers as an argument, Remote Control does not:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Why it works this way

**Remote Control takes no `--mcp-config`.** For normal launches ccam hands the group's servers to Claude Code as `--mcp-config=<path>` and never writes Claude Code's files. `claude remote-control` refuses that option before its command name and accepts none after it, so before ccam 0.6.0, `ccam work rc` failed for every account that shared MCP servers. ccam now leaves the option out.

**Its sessions read servers only from `.claude.json`.** The sessions Remote Control starts do not inherit options from the command line. The one place they find user MCP servers is the account's `.claude.json`. That file also holds the account's login details and per-folder state, so ccam cannot link it to a shared copy as it does with `settings.json`. Copying the servers into it is the only way to give the sessions the shared list.

**ccam asks before writing.** With the default flag strategy, ccam leaves `.claude.json` alone at launch. Remote Control is the one launch where it has to write there, so it asks first, keeps a backup, and touches nothing but `mcpServers`. (Outside launches, leaving or joining MCP sharing and `ccam mcp import` write there too, as their plans say.)

**It does not write while the account is running.** Claude Code rewrites `.claude.json` during a session, and a change made from outside would be lost or overwritten. That is why ccam skips the copy when any session of the account is running, including another Remote Control server and Claude Code started without ccam.

**It asks once per change.** After a copy, `sync-stamp` holds the modification time of the group's `mcp.json`. ccam asks again only when that file changes, for example after `ccam mcp add`, or when a copied server is gone from `.claude.json`. If you answer `N`, there is no copy yet, so it asks on every Remote Control launch.

**`ccam mcp` keeps the copies current.** An entry in `.claude.json` wins over `--mcp-config` when the names clash, so an old copy would hide a newer definition on normal launches too. `ccam mcp add`, `add-json`, `remove` and `import`, and `ccam migrate`, therefore update the copies of every account that is not running, and say so when a copy changed:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Without a terminal the answer is no.** Remote Control often runs from a script, a service or `nohup`, where nobody can answer. ccam then starts without the servers instead of waiting.

## Problems and fixes

### The shared MCP servers are missing in Remote Control sessions

The question was answered `N`, there was no terminal, or the account was running. Stop Remote Control, close the account's other sessions, and run `ccam work rc` again in a terminal; answer `y`. To stop being asked, set `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. Servers already in a session do not change while it runs, so start a new session after the copy.

### "Account work is running, so its .claude.json is left alone"

A session of `work` is still open: a normal `ccam work`, a shortcut, or another Remote Control server. Exit it and run `ccam work rc` again. Sessions that another account runs do not count.

### ccam asks every time

You have not copied the servers yet. Answer `y` once, or set `mcp_strategy = "sync"`. If you never want the servers in Remote Control, stop sharing MCP with the account without keeping them: `ccam group leave default work --fresh`, or keep MCP in a group of its own that the account is not in. A plain `ccam group leave` copies the group's servers into the account's `.claude.json`, so Remote Control has them.

### `ccam work rc` fails immediately with an error about `--mcp-config`

Your ccam is older than 0.6.0 and passes `--mcp-config` to Remote Control. Run `ccam update`. With a current ccam, the same error means that an unknown option with a value came before `rc`, and ccam mistook the value for the command; put `rc` right after the account name.

### A change made with `ccam mcp` does not reach Remote Control

`ccam mcp add` and `remove` update copies only in accounts that are not running, and say so: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Stop Remote Control and start it again; ccam sees the changed `mcp.json` and asks again.

### A server of my own is used instead of the group's

If the account has its own server with the same name as a group server, the copy leaves it in place and says so: `account work has its own MCP servers fs, so group default's of those names are not copied there`. To use the group's, remove the account's (`claude mcp remove -s user <name>`); to give your version to every account, move it into the group with `ccam mcp import work <server> --force`. Backups of `.claude.json` are in `~/.local/state/ccam/mcp/work/backups/`, which keeps the newest five.

### I want to remove the copied servers

`ccam mcp remove <server>` removes a server from the group and from every copy. To return to the state before the copy, restore the newest backup from `~/.local/state/ccam/mcp/work/backups/` while no session of the account is running.

### Remote Control asks whether to trust the folder, or exits without a terminal

Folder trust is stored per account in `.claude.json`, so a folder you trusted as `personal` is new to `work`. Run `ccam work` once in that folder and accept, then start Remote Control there.

### "You must be logged in with a Claude account that has a subscription"

Remote Control needs a Claude subscription login. Check the account's email in `ccam list`, and sign in again with `ccam login work` if needed.

### The session does not appear on the phone

The Claude app or claude.ai/code is signed in with a different Claude account. Compare it with the email `ccam list` shows for `work`.

### ccam warns that another Claude Code session shares conversation history in this folder

Another account in the same history group is running in the same folder. Sessions that continue the same conversation from both accounts can mix up the history. Answer `y` only if they will work on different conversations, or start Remote Control in another folder.

## Related

- [Launching Claude Code](/features/launch)
- [Shared MCP servers](/features/mcp)
- [Sharing with groups](/features/share)
- [Managing accounts](/features/accounts)
