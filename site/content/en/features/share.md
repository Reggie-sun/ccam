---
title: Sharing with groups
description: A ccam group shares chosen items between its accounts through one folder; ccam group moves the data with a plan, backups and conflict choices.
updated: 2026-10-26
order: 4
translation: source
---

Every account has its own folder, but most of what is in it does not need to be separate. A **group** says what is shared and by whom: every account in the group links those items to one folder, `~/.local/share/ccam/groups/<group>/`, so they all read and write the same copy. Login tokens are never shared.

An account can be in several groups, but each item comes from one group only. A typical setup is one group for settings and extensions, and a second one for conversation history between two accounts:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

The first account you add creates `default`, and later accounts join it unless you say otherwise.

## Usage

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Every command that moves data prints a plan, asks about each file that differs from the group's copy, and asks `Proceed? [y/N]` before it changes anything:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

You can make the same changes on the **Groups** tab of the [dashboard](/features/dashboard); press `s` on an account to choose its groups.

## Items

| Item | What it covers |
| --- | --- |
| `mcp` | MCP servers, through the group's `mcp.json` ([Shared MCP servers](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | the global `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` and `history.jsonl`, so another account can `--continue` |

If none of an account's groups shares an item, that item stays in the account's own folder.

## Commands and options

| Command | What it does |
| --- | --- |
| `ccam group` or `ccam group list` | list the groups |
| `ccam group show <group>` | what the group shares, its accounts, its folder and MCP servers |
| `ccam group add <group> --items <items> [--members <accounts>]` | create a group; its accounts' items move into it |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | change what the group shares, for every account in it |
| `ccam group join <group> <account>...` | put accounts in the group |
| `ccam group leave <group> <account>...` | take accounts out; each keeps a copy, MCP servers included |
| `ccam group rename <group> <new-name>` | rename in the config only |
| `ccam group remove <group> [--leave]` | remove a group |

| Option | What it does | Default |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: the account whose copy goes into the group first | the default account, then by number |
| `--replace` | `join`: leave the groups that share the same items first | off |
| `--fresh` | accounts that stop sharing an item start without it instead of with a copy (an imported folder still keeps one) | off |
| `--leave` | `remove`: take the group's accounts out first | needed when it has accounts |
| `--prefer <group\|mine\|merge>` | how to decide conflicts | ask; with `--yes`, each file's default |
| `-y`, `--yes` | run without asking; conflicts follow `--prefer` | off |
| `--dry-run` (global) | print the plan, change nothing | off |

Accounts are given by name or number; items as a comma-separated list. `ccam add` takes `--group <group>` (repeatable) or `--no-group`.

## How it works

The files involved:

- `~/.config/ccam/config.toml`: the `[[group]]` entries: `name`, `items`, `members` and, when the folder's name differs, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: the account's `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<folder>/`: the shared copies.

When `config.toml` is read, each account's groups are resolved into one folder per item. This takes about 25 µs for five accounts, so a launch still reads only `config.toml`. A config where an account would get an item from two groups is an error, and nothing starts until it is fixed.

A change works on a copy of the config:

1. The command edits the copy (`AddGroup`, `SetGroupItems`, `JoinGroup` …), which checks the rules.
2. ccam compares old and new and, for each account whose items change, builds a plan (`internal/group`, `share.Build`). The plans are built one after another, and each sees what the earlier ones will put into the group folders, so only the first account moves its copy in. The group's first copy comes from, in this order: the group folder (if it already has one), the `--from` account, the default account, then the other accounts by number.
3. ccam refuses if one of those accounts is running, however Claude Code was started (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), and on Windows without Developer Mode if a file would need a link.
4. ccam resolves conflicts (by asking, or with `--prefer`), prints the plan and asks for confirmation.
5. Right before saving, ccam checks again. If an account started while ccam waited for your answers, the change stops; it also stops if another ccam process changed the same groups in the meantime (`the groups changed while this plan waited; nothing was moved, make the change again`). The new config is saved before anything moves, and the change is noted in `~/.local/state/ccam/group-change.json`. If a step stops halfway, the note stays: `ccam doctor --fix` finishes the change in its order, with the answers you gave, MCP steps included, and until then no other group change starts (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``).
6. The MCP steps run first. When accounts take extensions out of a group, the group's plugin lists are pointed at the group folder wherever their paths run through an account's folder (`group plugin lists:` in the plan). Then the plans run in order. Every step moves data into the group or into the backup folder; only files identical to the group's copy are deleted instead of backed up.

ccam keeps an allow list of the paths each item covers (`share.Layout`). Anything not on it can never reach a group. ccam also names `.credentials.json` and `.claude.json` as never shared (`share.Never`), and stops with an internal error if either ever appears in the layout (`Build`). The `mcp` item has no paths: the server list lives inside `.claude.json`, so [MCP sharing](/features/mcp) works through the group's `mcp.json`.

This is what happens to one folder, `projects/`, when `side` starts sharing history and the group already has content:

<!-- motion:share-on -->

If the group has nothing at that path yet, the account is the first to share it and its folder is moved there as it is. Otherwise each file the group lacks is moved into the group, and files that exist on both sides stay in the account's folder, so the group's version is kept. Files identical to the group's are deleted, and the plan names the ones that differ (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

The account's folder then goes to `.ccam-backup/<time>/` inside the account folder, and a link takes its place. If the group's entry is itself a link to the joining account's own entry, the account's content goes into the group in its place. Conversation files are named after their session ID, so histories from several accounts combine without clashes. The `plugins/` registries store absolute paths through the account folder; ccam points them at the group's folder and merges the account's entries into a registry the group already has.

Leaving reverses the link:

<!-- motion:share-off -->

ccam copies the group's current content into `.ccam-backup/incomplete-<name>` in the account folder first, and replaces the link only once the copy is complete. If the copy stops, the link stays in place and `ccam doctor --fix` copies again (the plan of such a finish shows `restore <entry>: … (copied back from the group before a stop)`). The group keeps its content for the other accounts.

- ccam removes the link only if it still points into the groups folder as it did when the plan was created. A path that changed since, such as a real file where the link was, is left as it is and the change stops.
- With `--fresh`, the copy is skipped.
- An account that moves an item to a group whose folder is still empty brings a copy of what it saw, so it keeps its files. That copy is staged in `<group>/.ccam-incomplete-<name>` until it is complete.
- An account that stops sharing MCP keeps the group's servers: they go into its own `.claude.json` (a backup is kept), except where the account already has a server of the same name, which stays as it is. With `--fresh` the plan lists the servers that stop loading instead. Joining again removes the copies that nobody changed.
- Leaving `extensions` also points the group's plugin lists at the group folder where their paths ran through the leaving account's folder; for the other accounts, the rewritten paths lead to the same plugins.

A single file that exists on both sides with different content is a conflict:

<!-- motion:share-conflict -->

You choose for single files such as `settings.json`; inside folders the group's copy stays, and the plan names the files that differ. The choices depend on the file; the first one is the default:

| File | Choices |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| Other `.json` files | `keep group`, `use mine`, `merge JSON` |
| Anything else | `keep group`, `use mine` |

`merge JSON` merges objects recursively, joins arrays without duplicates, and keeps the group's value when both sides set the same key differently, so the other accounts see no change. Every choice except `keep group` first saves the group's copy in the account's `.ccam-backup/<time>/group/`. Your own file always goes to the backup.

At every launch, ccam compares each path in the layout with the account's groups using `lstat` and `readlink` only (`share.Check`). A missing link is created on the spot. A real file where a link belongs, a link to the wrong place, or a stale link stops the launch and points to `ccam doctor --fix`.

### Renaming, removing and the old format

- `ccam group rename` changes only the config. The folder keeps its name and the group records it as `dir`, so nothing moves and nothing can be lost.
- `ccam group remove` refuses a group with accounts unless you add `--leave`. Each account then keeps a copy, and the folder goes to `~/.local/share/ccam/trash/` unless another group uses it.
- Before groups, `config.toml` stored a `[profile.share]` table per account. ccam reads that format as groups without moving anything. On the first save it writes the new format and keeps the old file as `config.toml.bak`. Items shared by the same accounts become one group named after the folder; the rest get names such as `default-history` and keep using the same folder.

### On Windows

Windows creates symbolic links only in Developer Mode (Settings > System > For developers) or for an administrator. Without that right, a shared folder becomes a directory junction, which works the same for Claude Code. A shared file (`settings.json`, `CLAUDE.md`, `history.jsonl`) can only be a symbolic link, so a group change that needs one is refused before anything moves. For the same reason the first `ccam add` there creates the `default` group with only `mcp` and `extensions` and says so, and `ccam add --group` into a group that shares settings, memory or history is refused.

## Safety and edge cases

- **Running accounts.** A change is refused while Claude Code is running as one of the accounts it moves, because the session keeps writing into those files, whether ccam started it or not (plain `claude`, an IDE, background sessions). The check is repeated right before saving. Leaving and joining affect only the named account; leaving `extensions` also rewrites paths in the group's plugin lists, without changing what they mean for the others. `--dry-run` still prints the plan.
- **Backups.** Each plan gets its own `.ccam-backup/<time>/` folder in the account folder. ccam never deletes them.
- **Your own links.** If a path is already a link you made, for example to a dotfiles repository, and the group is empty, ccam copies the link's target into the group and moves the link itself to the backup. The target is not touched.
- **Imported folders.** An account registered in place with `ccam import` (such as `~/.claude`) changes in place, so plain `claude` sees the change too; the plan says so. When it leaves a group it always keeps a copy, even with `--fresh`, and `ccam remove` takes it out of its groups that way before it forgets the account. Its shared files live in the group folder while it is in a group, so remove the account (or leave the groups) before you delete `~/.local/share/ccam`.
- **Two accounts in one folder.** Accounts that share history write the same files. When a second one starts in a folder where one already runs, ccam warns and asks.
- **Manual edits.** If `members` lists a name that is not an account, ccam ignores it and `ccam doctor` reports it. The names `local`, `on`, `off` and `none` cannot be group names.

## Related

- [The dashboard](/features/dashboard)
- [Share MCP servers, settings and history](/docs/sharing)
- [Shared MCP servers](/features/mcp)
- [Copying an existing setup](/features/migrate)
- [Command reference](/docs/commands)
