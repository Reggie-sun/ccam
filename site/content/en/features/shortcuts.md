---
title: Shortcut commands
description: Give an account its own command such as ccwork with ccam shortcut, and recreate the links in ~/.local/bin with ccam shims install.
updated: 2026-10-26
order: 3
translation: source
---

A shortcut is a command you name yourself that starts Claude Code as one account: `ccwork --resume` does what `ccam work --resume` does. It is a link to ccam in `~/.local/bin`, so it works in bash, zsh, scripts and anything else that runs programs from `PATH`. Shortcuts are optional; create one for the accounts you start often.

## Usage

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccam confirms each change:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

An account has at most one shortcut, stored as `shortcut` in `~/.config/ccam/config.toml`. The `ccam add` wizard asks `Create a shortcut command for <name>? [y/N]` and then `Shortcut name [cc<name>]:`. Outside the wizard, `ccam add` creates one only with `--shortcut`. `ccam list` shows each account's shortcut, or `-`.

A shortcut passes every argument to Claude Code, the same way `ccam <name>` does after the name. That includes `--dry-run` and `--help`: to preview a launch, use `ccam --dry-run work`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--remove` | Delete the account's shortcut and its link | off |
| `-h`, `--help` | Help for the command | |
| `--dry-run` | `ccam shortcut` supports it: it checks the name and prints `Would set the shortcut of work to ccwork.` or `Would remove ccwork, the shortcut of work.` without saving or linking | off |

`ccam shims install` takes no arguments and refuses `--dry-run` with `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Name rules, checked in this order:

- 1 to 64 characters of `a-z`, `0-9`, `.`, `-` and `_`, starting with a letter or digit;
- not `ccam`, `claude` or a ccam command such as `run` or `usage`;
- not another account's shortcut;
- not an existing file in `~/.local/bin`, unless it is a link ccam created (such as this account's own shortcut);
- not the name of another executable on your `PATH`.

`config.toml` is checked against the same rules when it loads. If you edit a shortcut by hand so that it breaks the pattern, is `claude` or `ccam`, or is used by two accounts, every command fails with a message naming the account until you fix it.

Giving a command and `--remove` together stops with `give a command or --remove, not both`. Setting the name an account already has prints `work: ccwork (unchanged)` and recreates its link if it is missing.

## How it works

Setting a name runs `checkShortcut`, which calls `shim.Check` with the bin folder and the real path of the ccam binary:

<!-- motion:shortcuts-check -->

The name moves down the checks in order. Beside each check is a name it refuses: `Work` fails the pattern, `run` is a ccam command, `ccside` belongs to another account, `ccfile` is an existing file in `~/.local/bin`, and `git` would hide `/usr/bin/git`. Only a name that passes every check is saved and linked.

1. ccam matches the name against `^[a-z0-9][a-z0-9._-]{0,63}$`, then `config.Reserved`, then `cfg.FindShortcut` (`shim.Check`). It calls `os.Lstat` on `~/.local/bin/<name>`, and for every other folder on `PATH` it calls `os.Stat` and refuses an executable file that is not a link to ccam.
2. ccam saves `config.toml` with `fsutil.AtomicWrite` (temporary file, `fsync`, rename).
3. ccam deletes the old shortcut's link, but only if it is a link to ccam (`shim.Remove`).
4. ccam goes through every configured shortcut in name order and creates the link of each one that is missing (`shim.Install`), with `os.Symlink` to the ccam binary, resolved with `os.Executable` and `filepath.EvalSymlinks`. A file in the way of one shortcut does not stop the others; the error names each.

ccam records every link it creates, and the binary each link runs, in `~/.local/state/ccam/shortcuts.json`. A link counts as ccam's own when `os.Readlink` returns the ccam path, or when the link and the ccam path resolve to the same file. A recorded link whose binary was moved or replaced since (the ccam binary moved to another folder, a reinstall) is outdated, and ccam points it at this ccam binary (`shim.Install`). A recorded link to another installed copy of ccam is left as it is. Every remove and every overwrite check goes through these tests, so ccam does not touch files it did not create, such as your own copy of ccam or a link you made yourself.

When you run a shortcut, the name you typed is the only thing ccam needs:

<!-- motion:shortcuts-link -->

The shell finds `ccwork` on `PATH`, the kernel follows the link and starts the ccam binary, and `argv[0]` is still `ccwork`. In `app.Main`, any name other than `ccam` loads `config.toml` once and calls `shim.Lookup`, which finds the account whose `shortcut` equals `filepath.Base(argv[0])`. ccam then takes the same steps as [`ccam <name>`](/features/launch) with the config already loaded, and passes `argv[1:]` to `claude` untouched (`launcher.LaunchWith`).

Why a symlink and not a shell alias or a wrapper script: a symlink works in every shell and in scripts, it avoids starting an extra shell, which would cost about 1 to 2 ms, and the launch still reads one file. `scripts/perf-check.sh` measures a shortcut with `hyperfine` against the same 10 ms (p95) limit as `ccam <name>`.

### On Windows

A shortcut is `ccwork.exe` in `%USERPROFILE%\.local\bin`, and ccam drops the `.exe` (case-insensitively) before the lookup (`shim.Name`), so `ccwork`, `ccwork.exe` and `CCWORK.EXE` all start `work`. ccam first creates a symbolic link, which Windows allows in Developer Mode or for an administrator. Without that right it creates a hard link: a second name for the same `ccam.exe` file, which needs no privilege. A hard link counts as ccam's own when it is the same file as `ccam.exe`.

A hard link always runs the file it was created from, so when `ccam update` puts a new file in place of `ccam.exe`, it renews every hard-linked shortcut with `shim.Relink` and prints `✓ shortcut ccwork renewed`. Running `install.ps1` again does the same through `ccam shims install` (`OK 2 shortcuts now run the new ccam`), and so do `ccam shims install` and `ccam doctor --fix` for any shortcut that is a copy of an older version of ccam. The name check looks for `<name>` with each extension in `PATHEXT` (`.exe`, `.cmd` and so on) in the other `PATH` folders.

Type a shortcut in a shell (PowerShell, cmd.exe, Git Bash). A symbolic link opened from Explorer, the Run dialog or `Start-Process` starts as plain `ccam`, which opens the dashboard instead of the account.

## Safety and edge cases

- **A link no account claims.** If a link to ccam in `~/.local/bin` has a name that no account uses, for example after you edited `shortcut` in `config.toml` by hand, ccam does not fall back to the default account. It stops with `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Set the name again with `ccam shortcut`, or delete the link.
- **Existing files.** ccam never replaces a file it did not create. `ccam shims install` reports `<path> already exists and is not a ccam link, so it was left alone` and still creates the other links. `ccam add` and `ccam import` then print `shortcut: ccwork (not created; see above)`.
- **A moved ccam binary.** If you move the ccam binary, the links ccam created point to a path that no longer exists. `ccam shims install` (or `ccam doctor --fix`) points them at the ccam binary you run it with.
- **Rename and remove.** `ccam rename` keeps the shortcut name, and the link keeps working, because it points to ccam and ccam looks up the name in `config.toml`. `ccam remove` deletes the account's link. `ccam shortcut <account> <new>` deletes the old link before it creates the new one.
- **PATH.** When ccam creates a link and `~/.local/bin` is not on your `PATH`, it prints `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` followed by the `export PATH=...` line to add. On Windows it prints the PowerShell line that adds the folder to your user `PATH` instead.
- **Checking later.** `ccam doctor` reports `missing shortcuts in /home/me/.local/bin: ccwork`, shortcuts that run an older version of ccam, and warns when the folder is not on `PATH`. `ccam doctor --fix` creates the missing links and renews the outdated ones, and says `→ fixed` only for those it created. A shortcut that runs another ccam that is still installed works, so `ccam doctor` lists it as information only; `ccam doctor --fix` points it at the ccam binary you run it with, and `ccam add` and `ccam import` show it as ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## Related

- [Launching Claude Code](/features/launch)
- [Managing accounts](/features/accounts)
- [Command reference](/docs/commands#shortcuts)
- [Install ccam](/docs/install)
