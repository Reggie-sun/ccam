---
title: Sharing between accounts
description: How ccam share links settings, CLAUDE.md, extensions and history into a group folder, with a plan, backups and conflict choices, and how off copies back.
updated: 2026-10-20
order: 4
translation: source
---

Every account has its own folder, but most of what is in it does not need to be separate. `ccam share` turns sharing on or off for one account and one item at a time: the account's file or folder is moved into a group folder and replaced by a symbolic link, so every account in the group reads and writes the same copy. Use it after `ccam add` to change what an account shares, or to put a few accounts in a group of their own.

## Usage

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Without item options, `ccam share` prints the current state:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

With options, ccam prints a plan, asks about each file that differs from the group's copy, and asks `Proceed? [y/N]` before it changes anything. A new account shares everything except conversation history; [Share MCP servers, settings and history](/docs/sharing) describes each item from the user's side.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP servers | unchanged |
| `--settings <on\|off\|group>` | `settings.json` | unchanged |
| `--memory <on\|off\|group>` | the global `CLAUDE.md` | unchanged |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | unchanged |
| `--history <on\|off\|group>` | `projects/` and `history.jsonl` | unchanged |
| `--fresh` | when turning sharing off, start empty instead of copying the group's content | off |
| `--prefer <group\|mine\|merge>` | how to decide conflicts | ask; with `--yes`, each item's default |
| `-y`, `--yes` | run without asking; conflicts follow `--prefer` | off |
| `-h`, `--help` | help for `share` | |
| `--dry-run` (global) | honoured: asks about conflicts, prints the plan, changes nothing | off |

The value `on` means the group `default`, `off` means local (not shared), and any other value is a group name. Group names become folder names under `~/.local/share/ccam/groups/`, so `ccam share side --history ../x` fails with `invalid group name: "../x"`. The account is given by name or by number, and only the items you name change.

## How it works

The files involved:

- `~/.config/ccam/config.toml`: each account's `[profile.share]` table, one value per item.
- `~/.local/share/ccam/profiles/<name>/`: the account's `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<group>/`: the shared copies.

1. `internal/cli/share.go` reads the config, applies each `--<item>` value to a copy of the account's settings with `config.Share.Set`, and refuses if the account is running.
2. `share.Build` compares the current settings with the target and lists the steps. It only reads the disk: `os.Lstat` and `os.Readlink` for every path in `share.Layout`. The plan is a list of operations: move into an empty group slot (adopt), move the files the group lacks (merge dir), write a resolved conflict into the group, back up, unlink, link, copy back, rewrite plugin paths and merge the plugin registry.
3. Conflicts are decided, by the prompt or by `--prefer`.
4. `Plan.Describe` prints the plan, and ccam asks.
5. The target settings are written to `config.toml` before anything moves. If `Plan.Apply` stops halfway, `ccam doctor --fix` and a rerun build their plan towards the saved target and finish the change.
6. `Plan.Apply` runs the steps in order. Every step moves data into the group or into the backup folder; nothing is deleted.

`share.Layout` is an allow list of the paths each item covers. Anything not on it can never reach a group. `share.Never` names `.credentials.json` and `.claude.json`, and `Build` stops with an internal error if either ever appears in the layout. The `mcp` item has no paths: the server list lives inside `.claude.json`, so [MCP sharing](/features/mcp) works through the group's `mcp.json` instead, and `--mcp` only changes the setting.

This is what happens to one folder, `projects/`, when `side` turns history on and the group already has content:

<!-- motion:share-on -->

If the group has nothing at that path yet, the account is the first to share it and its folder is moved there as it is, without a backup. Otherwise `mergeDir` moves each file the group lacks into the group, and files that exist on both sides stay in the account's folder, so the group's version is kept. The account's folder then goes to `.ccam-backup/<time>/` inside the account folder, and `os.Symlink` puts a link in its place. Conversation files are named after their session ID, so histories from several accounts combine without clashes.

Moving uses `rename(2)` and falls back to copy and delete across file systems. A relative link inside the moved tree that pointed outside it is rewritten as an absolute link (`fsutil.PortableLink`), so it still resolves from the new place. The `plugins/` registries (`installed_plugins.json`, `known_marketplaces.json`) store absolute paths through the account folder; ccam points them at the group's folder and, when the group already has a registry, merges the account's entries into it.

Turning an item off reverses the link:

<!-- motion:share-off -->

ccam removes the link only if it points into the groups folder, then copies the group's current content into the account with `CopyTree`. With `--fresh` the copy is skipped and the account starts with nothing at that path. The group keeps its content for the other accounts. Moving an item from one group to another unlinks it and links it to the new group, leaving the old group as it is; the plan says so in a `Note:` line.

A single file that exists on both sides with different content is a conflict:

<!-- motion:share-conflict -->

`share.NewConflict` lists the choices that fit the file. `history.jsonl` offers `merge by line`, `keep group` and `use mine`, with merging first. Other `.json` files offer `keep group`, `use mine` and `merge JSON`. Anything else offers `keep group` and `use mine`. The first choice is the default:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) merges objects recursively, joins arrays without duplicates, and keeps the group's value when both sides set the same key differently, so other accounts see no change in behaviour. `merge by line` keeps the group's lines, then adds your lines it lacks. Every choice except `keep group` first saves the group's copy to `.ccam-backup/<time>/group/` in the account folder, then writes the result with `fsutil.AtomicWrite`. Your own file always goes to the backup.

`--prefer group`, `mine` or `merge` decides every conflict without asking; `merge` means `merge JSON` or `merge by line`, whichever the file offers, and a file that offers neither keeps the default. `--yes` alone takes each file's default.

At every launch, `share.Check` compares each path in the layout with the setting using `lstat` and `readlink` only, which keeps it cheap enough for the launch path. A missing link is created on the spot, because there is nothing to lose. A real file where a link belongs, a link to the wrong place, or a stale link stops the launch:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` runs `share.Build` with the saved settings on both sides, which derives every step from what is on disk, keeps the group's content on conflicts and backs up anything it replaces.

### On Windows

Windows creates symbolic links only in Developer Mode (Settings > System > For developers) or for an administrator. `fsutil.Link` tries a symbolic link first. Without that right, a shared folder (`agents`, `skills`, `projects` and the others) becomes a directory junction, which any user may create, and works the same for Claude Code. A shared file (`settings.json`, `CLAUDE.md`, `history.jsonl`) can only be a symbolic link, so ccam stops with `Windows links files only in Developer Mode` and suggests turning it on or keeping that item local. `share.Check` reads junctions as links (`fsutil.LinkTarget`), and removing an account folder removes a junction itself, never the group content it points to.

## Safety and edge cases

- **Running accounts.** ccam refuses to change sharing while Claude Code runs as the account (`account work is running; quit it first, then change its sharing`), because the session keeps writing into the files being moved and lines it appends meanwhile would be lost. `--dry-run` still prints the plan.
- **Backups.** Each plan gets its own `.ccam-backup/<time>/` folder in the account folder, with a numbered suffix when two plans start in the same second. ccam never deletes them.
- **Partial failures.** If a step fails, ccam prints the error and where the originals are, and tells you to finish with `ccam doctor --fix`. The target is already saved, so the repair continues towards it.
- **Your own links.** If a path is already a link you made, for example to a dotfiles repository, and the group is empty, ccam copies the link's target into the group and moves the link itself to the backup. The target is not touched.
- **Login files.** `.credentials.json` and `.claude.json` are never linked, moved or copied into a group.
- **Two accounts in one folder.** Accounts that share history write the same files. When a second account of the same history group starts in a folder where one already runs, ccam warns and asks before it continues.

## Related

- [Share MCP servers, settings and history](/docs/sharing)
- [Shared MCP servers](/features/mcp)
- [Copy an existing setup](/features/migrate)
- [All commands](/docs/commands)
