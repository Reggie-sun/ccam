---
title: Accounts by folder
description: ccam folder makes a folder and everything below it use one account when you start Claude Code without an account name, such as work for your company code.
updated: 2026-10-26
order: 11
translation: source
---

After `ccam folder work ~/code/company`, any launch without an account name in `~/code/company` or its subfolders starts Claude Code as `work`. Outside it, the default account starts as before. Naming an account, with `ccam personal` or a shortcut such as `ccpersonal`, always starts that account.

## Usage

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

These launches use the folder's account:

- `ccam -p "…"` and any other launch whose first argument is an option for Claude Code
- `ccam run` without a name, and `ccam --dry-run`
- `ccam --fallback -p "…"` for its first attempt
- the dashboard, `ccam` alone, which starts with its cursor on the folder's account, even when you reached the folder through a symbolic link

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `<account>` | Name or number of the account to use | |
| `[folder]` | The folder, relative or absolute | the current folder |
| `--remove` | Remove the entry for exactly that folder | off |
| `--dry-run` | Supported: print the change without saving it | off |

## How it works

Each folder is one `[[folder]]` entry in `~/.config/ccam/config.toml`, the file every launch reads anyway:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` resolves the folder to an absolute path with symbolic links resolved and saves it. Setting a folder again replaces its account.
2. A launch without a name calls `launcher.HereProfile`. With no entries it takes the default account at once.
3. Otherwise it compares the working directory with each `path` as strings: the folder itself or one below it, never a folder that only starts with the same letters (`/code/company` does not cover `/code/companyx`). The deepest entry wins, so `~/code` can be `personal` while `~/code/company` is `work`.
4. Only if the working directory matches nothing as given does ccam resolve its symbolic links and compare again.
5. On Windows and macOS the comparison ignores upper and lower case, as their file systems usually do.

No file besides `config.toml` is read, so a launch inside a folder adds no measurable time (4.0 ms with or without entries, `hyperfine`, 300 runs).

## Safety and edge cases

- `ccam rename` carries the folder entries to the new name, and `ccam remove` deletes the entries of the removed account.
- An entry whose account no longer exists stops a launch without a name in that folder: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` reports such entries, and entries whose folder is gone.
- If you run `ccam folder --remove` in a subfolder that inherits its account from a parent folder, ccam names the parent folder and removes nothing.
- `ccam next` chooses by rotation or usage and does not look at folders.

## Related

- [Launching Claude Code](/features/launch)
- [Shortcut commands](/features/shortcuts)
- [Managing accounts](/features/accounts)
- [Command reference](/docs/commands)
