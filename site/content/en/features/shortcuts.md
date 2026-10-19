---
title: Shortcut commands
description: Give an account its own command such as ccwork with ccam shortcut, and recreate the links in ~/.local/bin with ccam shims install.
updated: 2026-10-19
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
ccam shims install            # recreate missing links of all accounts
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
| `--dry-run` | `ccam shortcut` honours it: it checks the name and prints `Would set the shortcut of work to ccwork.` or `Would remove ccwork, the shortcut of work.` without saving or linking | off |

`ccam shims install` takes no arguments and refuses `--dry-run` with `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Name rules, checked in this order:

- 1 to 64 characters of `a-z`, `0-9`, `.`, `-` and `_`, starting with a letter or digit;
- not `ccam`, `claude` or a ccam command such as `run` or `usage`;
- not another account's shortcut;
- not a file in `~/.local/bin` that is not already a link to ccam;
- not the name of another executable on your `PATH`.

Giving a command and `--remove` together stops with `give a command or --remove, not both`. Setting the name an account already has prints `work: ccwork (unchanged)` and recreates its link if it is missing.

## How it works

Setting a name runs `checkShortcut`, which calls `shim.Check` with the bin folder and the real path of the ccam binary:

<!-- motion:shortcuts-check -->

The name moves down the checks in order. Beside each check is a name it refuses: `Work` fails the pattern, `run` is a ccam command, `ccside` belongs to another account, `ccfile` is an existing file in `~/.local/bin`, and `git` would hide `/usr/bin/git`. Only a name that passes every check is saved and linked.

1. `shim.Check` matches the name against `^[a-z0-9][a-z0-9._-]{0,63}$`, then `config.Reserved`, then `cfg.FindShortcut`. It calls `os.Lstat` on `~/.local/bin/<name>`, and for every other folder on `PATH` it calls `os.Stat` and refuses an executable file that is not a link to ccam.
2. ccam saves `config.toml` with `fsutil.AtomicWrite` (temporary file, `fsync`, rename).
3. `shim.Remove` deletes the old shortcut's link, but only if it is a link to ccam.
4. `shim.Install` creates the link of every configured shortcut that is missing, with `os.Symlink` to the ccam binary, resolved with `os.Executable` and `filepath.EvalSymlinks`. It never removes anything.

A link counts as ccam's own when `os.Readlink` returns the ccam path, or when the link and the ccam path resolve to the same file. Every remove and every overwrite check goes through this test, so ccam does not touch files it did not create.

When you run a shortcut, the name you typed is the only thing ccam needs:

<!-- motion:shortcuts-link -->

The shell finds `ccwork` on `PATH`, the kernel follows the link and starts the ccam binary, and `argv[0]` is still `ccwork`. In `app.Main`, any name other than `ccam` loads `config.toml` once and calls `shim.Lookup`, which finds the account whose `shortcut` equals `filepath.Base(argv[0])`. `launcher.LaunchWith` then takes the same steps as [`ccam <name>`](/features/launch) with the config already loaded, and passes `argv[1:]` to `claude` untouched.

Why a symlink and not a shell alias or a wrapper script: a symlink works in every shell and in scripts, it adds no extra shell start (about 1 to 2 ms), and the launch still reads one file. `scripts/perf-check.sh` measures a shortcut with `hyperfine` against the same 10 ms (p95) limit as `ccam <name>`.

## Safety and edge cases

- **A link no account claims.** If a link to ccam in `~/.local/bin` has a name that no account uses, for example after you edited `shortcut` in `config.toml` by hand, ccam does not fall back to the default account. It stops with `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Set the name again with `ccam shortcut`, or delete the link.
- **Existing files.** ccam never replaces a file it did not create. `ccam shims install` stops with `<path> already exists and is not a ccam link, so it was left alone`.
- **A moved ccam binary.** Old links then point to a path that no longer exists, and ccam does not count them as its own. Delete them, then run `ccam shims install`.
- **Rename and remove.** `ccam rename` keeps the shortcut name, and the link keeps working, because it points to ccam and ccam looks up the name in `config.toml`. `ccam remove` deletes the account's link. `ccam shortcut <account> <new>` deletes the old link before it creates the new one.
- **PATH.** When ccam creates a link and `~/.local/bin` is not on your `PATH`, it prints `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` followed by the `export PATH=...` line to add.
- **Checking later.** `ccam doctor` reports `missing shortcuts in /home/me/.local/bin: ccwork` and warns when the folder is not on `PATH`. `ccam doctor --fix` creates the missing links.

## Related

- [Launching Claude Code](/features/launch)
- [Managing accounts](/features/accounts)
- [Command reference](/docs/commands#shortcuts)
- [Install ccam](/docs/install)
