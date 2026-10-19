---
title: Managing accounts
description: Add, log in, list, rename, remove and import Claude Code accounts with ccam add, ccam login, ccam list, ccam rename, ccam remove and ccam import.
updated: 2026-10-19
order: 2
translation: source
---

Every ccam account is a separate Claude Code config folder with its own login. These commands create, sign in, list, rename and remove accounts, and register folders you already have.

## Usage

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Wherever a command takes `<name|number>`, you can give the account name or its slot number, the `#` column of `ccam list`.

`ccam add` on a terminal with no options opens a wizard. It asks for the name, then shows a checklist of shared items (everything except conversation history is checked), then asks `Create a shortcut command for work? [y/N]` and, on `y`, `Shortcut name [ccwork]:`. Any of `--share`, `--shortcut` or `--no-login`, or a pipe instead of a terminal, skips the wizard; then a name is required and no shortcut is made unless `--shortcut` names one.

`ccam list` shows the slot, name, shortcut, email, shared items and last use, with `*` before the default account:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

With `--json`, each account is an object with `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (left out when the account was never started) and `share`. `share` maps every item (`mcp`, `settings`, `memory`, `extensions`, `history`) to its group name or to `"local"`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `ccam add --no-login` | Create the account without logging in; log in later with `ccam login <name>` | off |
| `ccam add --share <item=group>` | Share an item with a group, or `item=local` to keep it private. Repeatable | `mcp`, `settings`, `memory`, `extensions` in `default`; `history` local |
| `ccam add --shortcut <command>` | Also create this shortcut command in `~/.local/bin` | none |
| `ccam list --json` | Print the accounts as JSON | off |
| `ccam remove -y`, `--yes` | Remove without asking | off |
| `ccam import --name <name>` | Account name for the folder you give | derived from the folder name |
| `ccam import -y`, `--yes` | Import every folder found without asking | off |
| `--dry-run` | Honoured by `ccam list` and `ccam completion`, which change nothing. Refused by `add`, `login`, `logout`, `use`, `rename`, `remove` and `import`, with `` --dry-run is not supported by `ccam add`; nothing was changed `` (naming the command) | off |

`ccam login`, `ccam logout`, `ccam use` and `ccam rename` have no options of their own. `ccam completion` takes `bash` or `zsh`; anything else fails with `only bash and zsh are supported`.

Account names are 1 to 32 characters of `a-z`, `0-9`, `-` and `_`, starting with a letter or digit. A name cannot be all digits, because `ccam 2` means slot 2. A name cannot be a ccam command either (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), so `ccam <name>` is never ambiguous. The list is `config.Reserved` in `internal/config`.

## How it works

`ccam add work --shortcut ccwork` runs these steps in `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` checks the name, and `shim.Check` checks the shortcut before anything is created.
2. `profile.Add` creates `~/.local/share/ccam/profiles/work/` with mode `0700` and gives the account the lowest free slot number (`Config.NextSlot`), so a removed account's number is reused. The first account becomes the default.
3. `applyShare` links the new, empty folder into its share groups with `share.Build` and then `Apply`, for example `settings.json` to `groups/default/settings.json`. If this fails, the new folder is removed, because a leftover folder would block every retry with `folder already exists`.
4. The account is saved to `~/.config/ccam/config.toml` through `fsutil.AtomicWrite`: a temporary file, `fsync`, then `rename`, so a crash never leaves half a config.
5. `runLogin` starts `claude auth login` with `CLAUDE_CONFIG_DIR` set to the account folder. If the login fails, the account stays, and ccam prints the error and tells you to run `ccam login work` later.
6. `profile.ReadEmail` takes `oauthAccount.emailAddress` out of `profiles/work/.claude.json`. That file can be several megabytes, so `TopLevelValue` scans for the one key instead of parsing the whole file.
7. `shim.Install` links `~/.local/bin/ccwork` to the ccam binary, which looks up `argv[0]` in `config.toml` when it starts.

A login can take minutes, and another ccam may save the config meanwhile. So commands that wait save through `env.update`, which reloads `config.toml`, applies the one change (here, the email) and saves. `ccam remove` reloads after its prompt for the same reason.

`ccam login` runs step 5 and 6 for an existing account. `ccam logout` runs `claude auth logout` in the account and clears the saved email. `ccam use` only rewrites `default` in `config.toml`.

### Signing in on a remote machine

Over SSH, the browser half of `claude auth login` can break. Claude Code waits for the sign-in on a callback server at `http://localhost:<port>/callback` on the server, but the browser that finishes the sign-in runs on your computer, where that address leads nowhere and the browser shows an error such as `ERR_CONNECTION_REFUSED`. The address still carries the one-time code, so ccam can deliver it:

<!-- motion:accounts-relay -->

`runLogin` does not hand your terminal to Claude Code. It connects Claude Code's standard input to a pipe and `relayLogin` copies each line you type into it, so a code pasted at Claude Code's `Paste code here` prompt still arrives. When a line is a callback address, `callbackURL` accepts it only if the host is `localhost`, `127.0.0.1` or `[::1]`, a port is given, the path is `/callback` and both `code` and `state` are present; `http://` may be left out. `deliverCallback` then requests that address on the server with a 10-second limit and without following redirects, the way the browser would have. The request reaches the callback server of the waiting `claude auth login` (the dashed line in the diagram). Claude Code receives the code and saves the login, on Linux in `.credentials.json` in the account folder. A pasted line can never make ccam contact another host.

ccam prints this before the login starts:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

On success it prints `→ Passed the sign-in to Claude Code on this machine.`

While the login waits for input on a terminal, ccam turns off the terminal's echo and line mode and shows a `*` for each character you type, so a pasted code or callback address does not appear on screen. `readLine` does the line editing itself: `Enter` sends the line, `Backspace` erases one character, and `Ctrl+D` on an empty line ends the input. The terminal is restored when the login ends, and also on `Ctrl+C`, `SIGTERM` or `SIGHUP`, which still end ccam. Without a terminal on standard input and standard error, lines are copied as they come.

### Rename, remove and import

`ccam rename` and `ccam remove` act on the account's folder and on what ccam keeps under the account's name elsewhere. `nameState` lists those paths: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` and the usage caches `~/.local/state/ccam/usage/<name>.json` and `<name>.statusline.json`. `ccam remove work` handles them in this order:

<!-- motion:accounts-remove -->

`ccam remove` first refuses if the account is running, then asks. If the account shares extensions, `RewritePluginPaths` rewrites the group's `installed_plugins.json` and `known_marketplaces.json`, because Claude Code stores absolute plugin paths that may run through this account's folder. `profile.Trash` moves the folder to `~/.local/share/ccam/trash/<name>-<time>/`, the account leaves `config.toml`, the `nameState` paths are deleted so a later account with the same name does not inherit them, and the shortcut link is removed. To recover, move the folder back from the trash.

`ccam rename` renames the folder, rewrites the plugin paths in its `plugins/` from the old folder to the new one, carries the default along, and moves the `nameState` paths to the new name. The shortcut command keeps its name. If the plugin rewrite fails, the folder is renamed back, since the config is not saved.

`ccam import` registers a folder where it is, as `dir` in `config.toml`. It never moves or copies it: a copied refresh token would make the copies log each other out. Without a folder argument it looks for unregistered `~/.claude` and `~/.claude-*` folders and asks about each; names come from the folder (`~/.claude` is `main`, `~/.claude-work` is `work`, `~/.claude-2` is `acct2`), and a folder whose name cannot be derived is listed with the `--name` command to use. Without `CLAUDE_CONFIG_DIR`, Claude Code keeps `.claude.json` in your home folder next to `~/.claude`, so when you import `~/.claude` and it has no `.claude.json` of its own, ccam creates the link `~/.claude/.claude.json` to `~/.claude.json`. Plain `claude` and `ccam main` then see the same account. Imported accounts share nothing until you run `ccam share`, and no shortcut is created.

`ccam completion bash|zsh` prints a completion script that also completes account names.

## Safety and edge cases

- `ccam rename` and `ccam remove` refuse while the account runs (`refuseRunning`), with `work is running; quit its Claude Code session first`, because the session keeps writing into the folder.
- `ccam remove` on an imported account only unregisters it: `✓ main unregistered (folder ... left in place)`. Your `~/.claude` is not moved to the trash.
- `ccam add` refuses a name whose folder already exists in `profiles/`, and points to `ccam import`.
- `ccam import` refuses a folder that is already registered, a path that is not a folder, and a name that is taken.
- `ccam logout` and `ccam login` do not check whether the account is running.
- If the relay cannot reach Claude Code, for example after the login ended, it prints `ccam: could not pass the sign-in to Claude Code:` with the reason. An HTTP error usually means the code expired; start the login again.

## Related

- [Launching Claude Code](/features/launch)
- [Shortcut commands](/features/shortcuts)
- [Share MCP servers, settings and history](/docs/sharing)
- [Multiple Claude Code accounts](/guides/multiple-claude-code-accounts)
