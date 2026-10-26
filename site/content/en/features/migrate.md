---
title: Copying an existing setup
description: Copy the settings, CLAUDE.md, extensions and MCP servers of an existing ~/.claude into a ccam group with ccam migrate.
updated: 2026-10-26
order: 6
translation: source
---

If you used Claude Code before ccam, your setup lives in `~/.claude` and `~/.claude.json`. `ccam migrate` copies it into a share group, so every account that shares those items through the group starts with the setup you already have. The source folder is only read, and plain `claude` keeps working as before. Run it once after your first `ccam add`; when `~/.claude.json` has MCP servers, the first `ccam add` names them in its tip.

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
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space (or `x`) toggles the item under the cursor, the arrow keys or `j`/`k` move, Enter shows the plan, and Esc, `q` or Ctrl+C prints `Cancelled.`. Items the source has nothing for are marked `not found`. If you uncheck everything, ccam prints `Nothing selected.` and stops.

Then ccam prints the plan and asks before it copies:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

A group affects only its own accounts, and only the items it shares, so the summary lists the group's accounts. When the group does not share one of the copied items, or some accounts are not in it, ccam prints the `ccam group set <group> --add <items>` and `ccam group join <group> <accounts>` commands that change that.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `[folder]` | The Claude Code folder to copy from | `~/.claude` |
| `--settings` | Copy `settings.json` | off |
| `--memory` | Copy the global `CLAUDE.md` | off |
| `--extensions` | Copy `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | off |
| `--mcp` | Copy user MCP servers (`mcpServers` in `.claude.json`) | off |
| `--history` | Copy conversation history (`projects/`, `history.jsonl`) | off |
| `-g`, `--group <group>` | The share group to copy into | `default` |
| `-y`, `--yes` | Copy without asking; on a conflict the group's copy is kept | off |
| `--dry-run` (global) | Print the plan and stop before any question or copy | off |

On a terminal, if you give no item flags, you choose the items in the checklist. Outside a terminal, or with `--yes`, ccam copies everything except conversation history. History is left out by default because `projects/` can be large.

`--group on` means `default`. `--group off` and `--group local` are refused with `--group takes a group name, not "off"`. `ccam migrate` supports `--dry-run`, and the checklist still appears in a dry run so you can preview a selection.

## How it works

ccam builds a complete plan before it copies anything. It picks the items (`internal/cli/migrate.go`), then reads the source and the group and builds the plan (`migrate.Build`). Nothing is written until `Plan.Apply`, after you confirm.

1. **Find the source.** ccam refuses a path that is missing, not a folder, or holds no known Claude Code entry (`checkSource`). For MCP servers alone, an empty or missing `~/.claude` is fine as long as `~/.claude.json` has servers. ccam finds the `.claude.json` that belongs to the folder (`migrate.ClaudeJSON`): inside it, or for `~/.claude` without one inside, `~/.claude.json` next to it. That is where Claude Code keeps it when `CLAUDE_CONFIG_DIR` is not set.
2. **Map items to paths.** Every path comes from `share.Layout`, the same allow list groups use. The MCP servers are read out of `.claude.json` with `mcp.LocalServers`, which uses the `profile.TopLevelValue` scanner rather than parsing the whole file.
3. **Copy, never move.** Files go to `~/.local/share/ccam/groups/<group>/` with `fsutil.AtomicWrite` or `fsutil.CopyFile`, which writes a temporary file beside the target and renames it, so an interrupted copy never leaves a partial file. MCP servers go into the group's `mcp.json` under its lock: a server that another ccam process adds in the meantime is kept, and the copies made for [Remote Control](/guides/remote-control) are updated afterwards.

Each item maps one source path to one group path:

<!-- motion:migrate-items -->

The two history rows are dashed because they are copied only with `--history`. `.credentials.json` and `.claude.json` are listed in `share.Never` and never reach the group: ccam stops with an internal error if one ever shows up in the layout (`Build`), and the folder walk skips them by name.

ccam gives every path a state, which is what each plan line describes (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: not in the source, or a file where a folder is expected. Listed once per item as `not found, skipped`.
- **Linked**: the source path already resolves to the group's copy, for example when the source is a ccam profile folder.
- **Same**: the group already has identical content.
- **New**: the group lacks it, so it is copied whole.
- **Merge**: a folder on both sides. ccam adds only the files the group lacks (`copyMissing`); files that differ stay as the group has them and the plan says `group copies kept`. For `plugins/`, the registries `installed_plugins.json` and `known_marketplaces.json` are merged, and the absolute plugin paths are rewritten to point at the group's folder.
- **Conflict**: a file on both sides with different content. You choose per file: `keep group`, `use mine`, plus `merge JSON` for `.json` files (a deep merge in which the group wins on clashing values) or `merge by line` for `history.jsonl`. An MCP server whose definition differs is a conflict too, with `keep group` or `use mine`. The first choice listed is the default when you press Enter; `--yes` always keeps the group's copy.

The checklist itself is `tui.MigratePicker` built on the `checklist` component that `ccam add` uses; ccam marks missing items with lstat calls and the MCP server list (`migrate.Found`).

### MCP sign-ins

MCP servers that sign in with OAuth keep their tokens in the login file, `.credentials.json` under `mcpOAuth`, not in the server list. ccam does not copy them. A copied sign-in holds the same refresh token as the original, and the server rotates that token at the first renewal: every other holder, the source `~/.claude` included, is then logged out of that server. So for each copied server the source is signed in to, the plan's `sign-ins` line lists the group's accounts that must sign in themselves with `/mcp`, and the summary repeats it. An entry without a token is not a sign-in and is not listed, nor is an account that has its own sign-in already, or a server whose definition the group keeps. No login file is written.

## Safety and edge cases

- The source is only read. No file in it is moved, linked or changed, and the source login is never written, including on macOS where it may be read from the Keychain.
- Before a group file is replaced, by `use mine`, a merge, or new MCP servers in an existing `mcp.json`, the group's copy is saved in `<group>/.ccam-backup/<time>/`. The folder name is unique per run, so two runs in the same second do not overwrite each other's backups. The summary prints the path.
- With `--yes`, every conflict keeps the group's copy, and ccam says how to choose instead: run the same command without `--yes`.
- If nothing would change, ccam prints `Nothing to copy.`; if every conflict keeps the group's copy, `Nothing to copy: the group keeps its own copies.`
- If a copy fails halfway, ccam stops with `The source folder was not changed. Group files replaced so far are saved in` followed by the backup folder. Running it again copies the files that were not copied completely.
- When several groups share MCP, the commands that `ccam doctor`, `ccam ext` and `ccam mcp list` suggest read `ccam migrate --mcp --group <a|b>`.
- A missing source ends with `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Related

- [Share MCP servers, settings and history](/docs/sharing)
- [Checking and repairing](/features/doctor)
- [Use multiple Claude Code accounts on one computer](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR explained](/guides/claude-config-dir)
