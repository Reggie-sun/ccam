---
title: Copying an existing setup
description: Copy the settings, CLAUDE.md, extensions, MCP servers and MCP sign-ins of an existing ~/.claude into a ccam group with ccam migrate.
updated: 2026-10-19
order: 6
translation: source
---

If you used Claude Code before ccam, your setup lives in `~/.claude` and `~/.claude.json`. `ccam migrate` copies it into a share group, so every account that shares those items through the group starts with the setup you already have. The source folder is only read, and plain `claude` keeps working as before. Run it once after your first `ccam add`, and again whenever an account starts sharing MCP servers.

## Usage

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

On a terminal, without item flags and without `--yes`, ccam first shows a checklist:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (or `x`) toggles the item under the cursor, the arrow keys or `j`/`k` move, Enter shows the plan, and Esc, `q` or Ctrl+C prints `Cancelled.`. Items the source has nothing for are marked `not found`. If you uncheck everything, ccam prints `Nothing selected.` and stops.

Then ccam prints the plan and asks before it copies:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

A group only reaches its own accounts, and only for the items it shares, so the summary lists the group's accounts. When the group does not share one of the copied items, or some accounts are not in it, ccam prints the `ccam group set <group> --add <items>` and `ccam group join <group> <accounts>` commands that change that.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `[folder]` | The Claude Code folder to copy from | `~/.claude` |
| `--settings` | Copy `settings.json` | off |
| `--memory` | Copy the global `CLAUDE.md` | off |
| `--extensions` | Copy `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | off |
| `--mcp` | Copy user MCP servers (`mcpServers` in `.claude.json`) | off |
| `--history` | Copy conversation history (`projects/`, `history.jsonl`) | off |
| `-g`, `--group string` | The share group to copy into | `default` |
| `-y`, `--yes` | Copy without asking; on a conflict the group's copy is kept | off |
| `--dry-run` (global) | Print the plan and stop before any question or copy | off |

With no item flag, the checklist decides on a terminal. Outside a terminal, or with `--yes`, ccam copies everything except conversation history. History is left out by default because `projects/` can be large.

`--group on` means `default`. `--group off` and `--group local` are refused with `--group takes a group name, not "off"`. `ccam migrate` honours `--dry-run`, and the checklist still appears in a dry run so you can preview a selection.

## How it works

`internal/cli/migrate.go` picks the items, then `migrate.Build` reads the source and the group and returns a plan. Nothing is written until `Plan.Apply`, after you confirm.

1. **Find the source.** `checkSource` refuses a path that is missing, not a folder, or holds no known Claude Code entry. `migrate.ClaudeJSON` finds the `.claude.json` that belongs to the folder: inside it, or for `~/.claude` without one inside, `~/.claude.json` next to it. That is where Claude Code keeps it when `CLAUDE_CONFIG_DIR` is not set.
2. **Map items to paths.** Every path comes from `share.Layout`, the same allow list groups use. The MCP servers are read out of `.claude.json` with `mcp.LocalServers`, which uses the `profile.TopLevelValue` scanner rather than parsing the whole file.
3. **Copy, never move.** Files go to `~/.local/share/ccam/groups/<group>/` with `fsutil.AtomicWrite` or `fsutil.CopyFile`; MCP servers go into the group's `mcp.json`.

Each item maps one source path to one group path:

<!-- motion:migrate-items -->

The two history rows are dashed because they are copied only with `--history`. `.credentials.json` and `.claude.json` are listed in `share.Never` and never reach the group: `Build` returns an internal error if one ever shows up in the layout, and the folder walk skips them by name.

`buildEntry` gives every path a state, which is what each plan line describes:

<!-- motion:migrate-states -->

- **Missing**: not in the source, or a file where a folder is expected. Listed once per item as `not found, skipped`.
- **Linked**: the source path already resolves to the group's copy, for example when the source is a ccam profile folder.
- **Same**: the group already has identical content.
- **New**: the group lacks it, so it is copied whole.
- **Merge**: a folder on both sides. `copyMissing` adds only the files the group lacks; files that differ stay as the group has them and the plan says `group copies kept`. For `plugins/`, the registries `installed_plugins.json` and `known_marketplaces.json` are merged, and the absolute plugin paths are rewritten to point at the group's folder.
- **Conflict**: a file on both sides with different content. You choose per file: `keep group`, `use mine`, plus `merge JSON` for `.json` files (a deep merge in which the group wins on clashing values) or `merge by line` for `history.jsonl`. An MCP server whose definition differs is a conflict too, with `keep group` or `use mine`. The first choice listed is the default when you press Enter; `--yes` always keeps the group's copy.

The checklist itself is `tui.MigratePicker` built on the `checklist` component that `ccam add` uses; `migrate.Found` marks missing items with lstat calls and the MCP server list.

### MCP sign-ins

MCP servers that sign in with OAuth keep their tokens in the login file, `.credentials.json` under `mcpOAuth`, not in the server list. Claude Code keys each entry as `<name>|` plus the first 16 characters of a `sha256` hash of the server's `type`, `url` and `headers`. The key does not depend on the folder, so a copied entry works in another account as long as the server definition is identical. `buildLogins` in `internal/migrate/logins.go` therefore copies only the entries of servers the group ends up with in the source's definition: new, already the same, or a conflict where you chose `use mine`. Sign-ins of plugin servers (keys starting with `plugin:`) go along when you copy extensions and the source has `plugins/`.

<!-- motion:migrate-signins -->

The entries go into each receiving account's own `.credentials.json`, written with `fsutil.AtomicWrite` at mode `0600`. They never go into the group. An account receives them only if it shares `mcp` (or, for plugin servers, `extensions`) through the target group, and only the keys it lacks are added: entries it already has, perhaps newer, are kept. The plan line `sign-ins` lists the servers, the receiving accounts and the skipped ones. ccam skips an account, and names it with the reason, when it is:

- registered in place with `ccam import`, since plain `claude` may be using its login file;
- running, since Claude Code rewrites its login file (checked again right before writing, reported as `started meanwhile`);
- not logged in yet;
- keeping its login in the macOS Keychain, which ccam does not write.

Sign in to the servers in a skipped account with `/mcp`. An account that starts sharing MCP servers later gets the sign-ins when you run `ccam migrate` again; the summary says so with ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Safety and edge cases

- The source is only read. No file in it is moved, linked or changed, and the source login is never written, including on macOS where it may be read from the Keychain.
- Before a group file is replaced, by `use mine`, a merge, or new MCP servers in an existing `mcp.json`, the group's copy is saved in `<group>/.ccam-backup/<time>/`. The folder name is unique per run, so two runs in the same second do not overwrite each other's backups. The summary prints the path.
- With `--yes`, every conflict keeps the group's copy, and ccam says how to choose instead: run the same command without `--yes`.
- If nothing would change, ccam prints `Nothing to copy.`; if every conflict keeps the group's copy, `Nothing to copy: the group keeps its own copies.`
- If a copy fails halfway, ccam stops with `The source folder was not changed. Group files replaced so far are saved in` followed by the backup folder.
- A missing source ends with `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Related

- [Share MCP servers, settings and history](/docs/sharing)
- [Checking and repairing](/features/doctor)
- [Multiple Claude Code accounts](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
