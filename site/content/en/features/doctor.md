---
title: Checking and repairing
description: Check share links, file permissions, logins and the Claude Code version with ccam doctor, and repair what is safe to repair with ccam doctor --fix.
updated: 2026-10-26
order: 9
translation: source
---

`ccam doctor` checks your ccam installation: the Claude Code binary, every account's folder, login and share links, the shared MCP servers, the usage endpoint and your shortcut commands. `ccam doctor --fix` repairs what it can do without losing data and moves any file it replaces into a backup folder first. Run it when `ccam <name>` stops with a share link problem, after you restore files by hand, or after a Claude Code update.

## Usage

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

The first line marks every account at once, so you see where the problems are before you read anything else. Below it, each account gets a box, and findings that belong to no account (Claude Code, environment variables, plugins, MCP servers, shortcuts, folder rules, groups and settings, `auto_switch`) go in a `General` box. The box title shows the account's worst mark and how many problems it has. Each row inside is one area that was checked, always `Folder`, `Login` and `Share links` in that order, then any others, with a mark: `✔` OK, `ℹ` information, `⚠` a warning, `✘` an error. A folder with nothing wrong shows `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

With `--fix`, every finding that was repaired ends in `→ fixed`, and it no longer counts as a problem:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

On a terminal, the boxes fit the terminal's width, long messages wrap and the terminal's own colors are used (`NO_COLOR` turns them off). A box is never narrower than 29 columns, so no word is cut off on a narrow terminal. Sent to a pipe or a file, the output has no colors and no message is wrapped or cut, so `ccam doctor | grep settings.json` finds the whole line; every box is then as wide as the longest line. The dashboard's Health tab shows the same boxes.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--fix` | Repair what is safe to repair (changed files are backed up) | off |
| `--dry-run` (global) | Run every check but repair nothing, even with `--fix` | off |

`ccam doctor` takes no arguments. It supports `--dry-run`: `ccam --dry-run doctor --fix` prints the findings of `ccam doctor`, adds what `--fix` would set aside, and changes nothing:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

The exit code is `1` when at least one `✘` finding is left unfixed, and `0` otherwise. Warnings and information never fail the run. When the run fails without `--fix`, ccam ends with a blank line and `` `ccam doctor --fix` repairs what it can. ``. A config that cannot be loaded also ends with `1` and a `ccam:` error.

## How it works

The checks run in a fixed order, and the findings appear in one box per account. ccam loads `~/.config/ccam/config.toml` (`internal/cli/doctor.go`) and runs the checks, each finding naming the account it belongs to and the area it checked (`doctor.Run` in `internal/doctor/doctor.go`). It then groups the findings into the boxes (`tui.HealthReport`). `--fix` is passed through as `Options.Fix`, set to false in a dry run.

<!-- motion:doctor-checks -->

The middle block runs once per account. The checks, and what `--fix` does for each:

1. **Claude Code version** (`checkClaude`). Resolves the `claude` binary as the launcher does and runs `claude --version`. A missing binary, a `claude_bin` that is not an executable file, or a failing `--version` is an error. No fix.
2. **Credentials in the environment** (`checkAuthEnv`). Warns when `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN` is set: every account launched from that shell inherits it, and Claude Code prefers it over the account's own login. No fix; unset it yourself.
3. **Folder** (`checkProfile`). A missing account folder is an error and ends that account's checks. Permissions with any group or other bits set are a warning; `--fix` runs `chmod 0700`. An imported `~/.claude` with a `.claude.json` of its own beside `~/.claude.json` is information: `ccam main` and plain `claude` then keep separate login details and MCP servers. Windows has no permission bits, so this check and the `0600` check below are skipped there; the folder in your user profile is protected by its access list.
4. **Login file** (`checkCredentials`). `.credentials.json` is checked with `lstat`. A symlink is an error, because each account must have its own token, and is not fixed. A mode other than `0600` is a warning; `--fix` runs `chmod 0600`.

   Then the login is read (on macOS from the Keychain first). A missing login, or a login file with no login in it, is a warning that points to `ccam login <name>`. An expired refresh token, or a login Claude Code cleared after the server refused it, is the error `` login expired → `ccam login <name>` ``; an expired access token is only information, since Claude Code refreshes it on the next launch. A login kept in Windows Credential Manager is information, since ccam cannot read it. A valid login shows the email the folder is logged in as now and the plan (`login valid (dev@company.com, Claude Team · Acme)`). None of these are fixed; `ccam login --expired` logs in every account that needs it.
5. **Share links** (`share.Check`). Compares every path in `share.Layout` with the account's groups, using only `lstat` and `readlink`. It reports five kinds: a missing link, a link that points elsewhere, a real file where a link belongs, a link left over for an item that is now local, and a correct link whose group folder is missing (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; a launch does not check this). A missing link and a missing group folder are warnings, the others are errors.

   A copy back from the group that stopped halfway is reported too: `--fix` puts a complete staged copy in place, but if a launch has recreated the path since then, it names the folder that holds the complete copy so you can sort it out yourself. Paths added to an item in a later version of ccam (`file-history/`, `tasks/` and `paste-cache/` for history, `keybindings.json` for settings) are only a warning while the account still has its own copy: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. A launch never stops for them.
6. **Unknown entries** (`checkUnknownEntries`). Lists top-level names in the account folder that ccam does not know (`share.Known`), as a hint that the Claude Code layout changed. Information only.
7. **Plugin paths** (`checkPlugins`). Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) store absolute paths. In each group's `plugins/`, paths that go through one account's folder break when that account is removed or renamed. `--fix` rewrites them to the group folder with `profile.RewritePluginPaths`, an atomic write of the registry file.
8. **MCP servers** (`checkMCP`). Loads each group's `mcp.json` once. A file that cannot be read is an error; a server that fails `mcp.Validate` is a warning. Servers in plain `claude`'s `~/.claude.json` that no account loads are a warning naming the `ccam migrate --mcp` command that copies them (with `--group <a|b>` when several groups share MCP). Servers only one account of a group has are information with the `ccam mcp import` command that shares them, and an account not signed in to MCP servers other accounts are signed in to gets a warning to sign in with `/mcp` (ccam never copies sign-ins). No fix.
9. **Usage response format** (`checkUsageFormat`). Asks the usage endpoint, without the cache, for one account after another until one answers, within a single 4-second timeout. With `usage_source = "claude"` it asks Claude Code for one account instead. If the response no longer parses, it warns that `ccam usage` falls back to the status line cache; when no account answered, it says so as information (`could not check the usage response format: …`). This is the only check that uses the network.
10. **Shortcuts** (`checkShims`). Only when an account has a shortcut. Missing links in `~/.local/bin` are a warning, and so are shortcuts ccam created that still run an older version of ccam, one that was replaced or moved (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` calls `shim.Install`, which creates and renews only the links ccam created and recorded in `~/.local/state/ccam/shortcuts.json`, and never touches other files. The warning ends in `→ fixed` only when each of those links was created; if one could not be, the reason is added. A shortcut that runs another ccam that is still installed works, so `ccam doctor` lists it as information only; `--fix` points it at this ccam binary. A bin folder that is not on `PATH` is a warning and is not fixed.
11. **Folders** (`checkFolders`). Only when [`ccam folder`](/features/folder) entries exist. An entry whose account does not exist is an error, and an entry whose folder is gone is a warning. Neither is fixed, because removing an entry would change which account starts there.
12. **Groups and settings** (`checkGroups`). A group member that is not an account, an account whose name has since become a ccam command (start it with `ccam run <name>` or rename it), an unknown `theme` or `usage_source`, and lines of `config.toml` that ccam does not read are warnings. The last kind is usually a setting placed below a `[[group]]`, `[[profile]]` or `[[folder]]` table, which makes it part of that table (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); commands that change the config refuse until it is fixed. No fix.
13. **Logins** (`checkIdentities`). Reads who each folder is logged in as from its `.claude.json`. A folder now logged in to a different account than the one `ccam login` recorded (for example after `/login` inside a session) is a warning; an account without a login record, such as one imported or created by an older version of ccam, is compared by its recorded email. A folder logged in to the same Claude account as another one is a warning too: the two share their usage limits. A folder that holds the same login token as another, a copied folder (plain `claude`'s `~/.claude` included), is an error, since whichever renews it first logs the other out; run `ccam login <name>` to give it its own. No fix.
14. **`auto_switch`** (`checkAutoSwitch`). Only when `auto_switch = true`. Information in the `General` box says what `auto_switch` does and links to [its description](/features/next#auto-switch). For each plan with seats among the accounts (Claude Team, Claude Enterprise), a line names them, says that `auto_switch` moves a session between seats of one organization, quotes Anthropic's seat rule with the help center's link, and leaves the choice to you. Claude Code versions older than 2.1.214, which do not tell hooks their process (`CLAUDE_PID`), are a warning there, since no session can be moved.

    An account without the ccam plugin, or with the plugin turned off, gets a warning in its own box with the command that fixes it in Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (or `/plugin enable ccam@ccam`). API key accounts are left out. On Windows the only line is `auto_switch is ignored on Windows; with relay = true, /ccam:switch followed by /exit still moves a conversation to another account.` No fix.

Before the accounts, a **group change** that stopped halfway is an error in the `General` box: ``the last group change stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` finishes it in the order it was planned, with the answers given to its conflicts; until then no other group change, rename or remove starts. Until it is finished, ccam only creates missing links for the accounts involved. If it cannot be finished, the error says that removing that file gives up the rest of the change.

`ccam --dry-run doctor --fix` shows the order (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) and what each account's part would set aside. A group change that another ccam process is carrying out right now is information, and while it is in progress, `--fix` repairs no share links (`its share links are not repaired while another ccam changes the groups; …`). `--fix` holds `~/.local/state/ccam/group-change.lock` while it repairs, and reads `config.toml` again after finishing a change.

### Repairing share links

The launch path creates a missing link on the spot, but stops for anything that could lose data and points here:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

For any share link problem, `--fix` calls `share.Build` with the account's share settings as both the current and the target state, then `Plan.Apply`. ccam re-derives every step from what is on disk (`Build`), giving the same plan a `ccam group` change uses, so the repair follows the same rules as joining a group.

<!-- motion:doctor-repair -->

Here `settings.json` is a real file. Because the group already has a different `settings.json`, the conflict takes its default choice, `keep group` (`KeepGroup` in the diagram), so the group's copy stays. The account's file moves to `.ccam-backup/<time>/settings.json` inside the account folder, and ccam links the path to the group (`os.Symlink`). `--fix` says so: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Other cases:

- If the group has no copy yet, the account's file moves into the group and becomes its first content.
- For a folder, files the group lacks move into the group. Files identical to the group's copy are deleted, and the ones that differ go to the backup, named in a note (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- A missing group folder is created first.
- For `history.jsonl`, the default is `merge by line`: the group's lines, then your new ones. The group's previous copy is saved in `.ccam-backup/<time>/group/`.
- A link to another group is replaced. A link to a folder outside ccam goes to the backup, and its content seeds the group if the group is empty.
- A link left for an item that is now local is removed, and the account gets a copy of the group's content.

Apart from removing links and deleting files identical to the group's copy, every step moves data into the group or the backup folder. Nothing is lost.

## Safety and edge cases

- A plain `ccam doctor` never changes anything. Use `--dry-run` with `--fix` to be sure in scripts.
- For an account that is running, however Claude Code was started, `doctor --fix` only creates missing links and missing group folders, which loses nothing. When something else is left, it adds `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, because the session keeps writing to the files a repair would move.
- If a share repair fails, ccam reports `could not repair share links: <error>` for that account and skips its unknown-entry check. The other accounts are still checked.
- If a `ccam group` change fails halfway, the new groups are already saved and the change is noted in `~/.local/state/ccam/group-change.json`; ccam tells you to finish with `ccam doctor --fix`, which carries out the rest of it, MCP steps included.
- A symlinked `.credentials.json`, an expired login and environment credentials are reported but never changed.
- Plugin registry rewrites are written atomically but not copied to `.ccam-backup/`.

## Related

- [Share MCP servers, settings and history](/docs/sharing)
- [Copying an existing setup](/features/migrate)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Command reference](/docs/commands)
