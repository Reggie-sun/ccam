---
title: Checking and repairing
description: Check share links, file permissions, logins and the Claude Code version with ccam doctor, and repair what is safe to repair with ccam doctor --fix.
updated: 2026-10-19
order: 9
translation: source
---

`ccam doctor` checks your ccam installation: the Claude Code binary, every account's folder, login and share links, the shared MCP servers, the usage endpoint and your shortcut commands. `ccam doctor --fix` repairs what it can do without losing data and moves any file it replaces into a backup folder first. Run it when `ccam <name>` stops with a share link problem, after you restore files by hand, or after a Claude Code update.

## Usage

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Each finding is one line with a mark: `✓` OK, `·` information, `!` a warning, `✗` an error. Account findings start with the account name.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

With `--fix`, every finding that was repaired ends in `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--fix` | Repair what is safe to repair (changed files are backed up) | off |
| `--dry-run` (global) | Run every check but repair nothing, even with `--fix` | off |

`ccam doctor` takes no arguments. It honours `--dry-run`: `ccam --dry-run doctor --fix` prints the same findings as `ccam doctor` and changes nothing.

The exit code is `1` when at least one `✗` finding is left unfixed, and `0` otherwise. Warnings and information never fail the run. When the run fails without `--fix`, ccam ends with a blank line and `` `ccam doctor --fix` repairs what it can. ``. A config that cannot be loaded also ends with `1` and a `ccam:` error.

## How it works

`internal/cli/doctor.go` loads `~/.config/ccam/config.toml` and calls `doctor.Run` in `internal/doctor/doctor.go`, which runs the checks in a fixed order and returns findings. `--fix` is passed through as `Options.Fix`, set to false in a dry run.

<!-- motion:doctor-checks -->

The middle block runs once per account. The checks, and what `--fix` does for each:

1. **Claude Code version** (`checkClaude`). Resolves the `claude` binary as the launcher does and runs `claude --version`. If the version does not start with `2.1.302`, the version ccam's folder layout was verified against (`doctor.LayoutVersion`), it adds `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. A missing binary or a failing `--version` is an error. No fix.
2. **Credentials in the environment** (`checkAuthEnv`). Warns when `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN` is set: every account launched from that shell inherits it, and Claude Code prefers it over the account's own login. No fix; unset it yourself.
3. **Folder** (`checkProfile`). A missing account folder is an error and ends that account's checks. Permissions with any group or other bits set are a warning; `--fix` runs `chmod 0700`.
4. **Login file** (`checkCredentials`). `.credentials.json` is checked with `lstat`. A symlink is an error, because each account must have its own token, and is not fixed. A mode other than `0600` is a warning; `--fix` runs `chmod 0600`. Then the login is read (on macOS from the Keychain first). No login is a warning that points to `ccam login <name>`. An expired refresh token is the error `login expired`; an expired access token is only information, since Claude Code refreshes it on the next launch. None of these are fixed.
5. **Share links** (`share.Check`). Compares every path in `share.Layout` with the account's share settings, using only `lstat` and `readlink`. It reports four kinds: a missing link, a link that points elsewhere, a real file where a link belongs, and a link left over for an item that is now local. A missing link is a warning, the others are errors.
6. **Unknown entries** (`checkUnknownEntries`). Lists top-level names in the account folder that ccam does not know (`share.Known`), as a hint that the Claude Code layout changed. Information only.
7. **Plugin paths** (`checkPlugins`). Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) store absolute paths. In each group's `plugins/`, paths that go through one account's folder break when that account is removed or renamed. `--fix` rewrites them to the group folder with `profile.RewritePluginPaths`, an atomic write of the registry file.
8. **MCP servers** (`checkMCP`). Loads each group's `mcp.json` once. A file that cannot be read is an error; a server that fails `mcp.Validate` is a warning. No fix.
9. **Usage response format** (`checkUsageFormat`). Asks the usage endpoint, without the cache and with the usual 3-second limit, for one account after another until one answers. If the response no longer parses, it warns that `ccam usage` falls back to the status line cache. This is the only check that uses the network.
10. **Shortcuts** (`checkShims`). Only when an account has a shortcut. Missing links in `~/.local/bin` are a warning; `--fix` calls `shim.Install`, which only creates missing links and never touches other files. A bin folder that is not on `PATH` is a warning and is not fixed.

### Repairing share links

The launch path creates a missing link on the spot, but stops for anything that could lose data and points here:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

For any share link problem, `--fix` calls `share.Build` with the account's share settings as both the current and the target state, then `Plan.Apply`. `Build` re-derives every step from what is on disk, the same plan `ccam share` uses, so the repair follows the same rules as turning sharing on.

<!-- motion:doctor-repair -->

Here `settings.json` is a real file. Because the group already has a different `settings.json`, the conflict takes its default choice, `KeepGroup`, so the group's copy stays. The account's file moves to `.ccam-backup/<time>/settings.json` inside the account folder, and `os.Symlink` links the path to the group. Other cases:

- If the group has no copy yet, the account's file moves into the group and becomes its first content.
- For a folder, files the group lacks move into the group; the rest goes to the backup.
- For `history.jsonl`, the default is `merge by line`: the group's lines, then your new ones. The group's previous copy is saved in `.ccam-backup/<time>/group/`.
- A link to another group is replaced. A link to a folder outside ccam goes to the backup, and its content seeds the group if the group is empty.
- A link left for an item that is now local is removed, and the account gets a copy of the group's content.

Apart from removing links, every step moves data into the group or the backup folder. Nothing is deleted.

## Safety and edge cases

- A plain `ccam doctor` never changes anything. Use `--dry-run` with `--fix` to be sure in scripts.
- `doctor --fix` does not check whether an account is running. Quit Claude Code for an account before repairing its share links, because the session keeps writing to those files.
- If a share repair fails, ccam reports `could not repair share links: <error>` for that account and skips its unknown-entry check. The other accounts are still checked.
- If `ccam share` fails halfway, it saves the target settings first and tells you to finish with `ccam doctor --fix`.
- A symlinked `.credentials.json`, an expired login and environment credentials are reported but never changed.
- Plugin registry rewrites are written atomically but not copied to `.ccam-backup/`.

## Related

- [Share MCP servers, settings and history](/docs/sharing)
- [Copying an existing setup](/features/migrate)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Commands](/docs/commands)
