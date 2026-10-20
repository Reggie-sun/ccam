---
title: The dashboard
description: ccam on a terminal opens a full-screen dashboard of accounts, groups, folder rules, extensions and checks; every change shows its plan first.
updated: 2026-10-20
order: 15
translation: source
---

`ccam` without arguments opens a full-screen dashboard when standard input and output are a terminal; otherwise it prints `ccam list`. The dashboard shows every account with its usage and groups, and does what the commands do: launch, add, log in, change groups, set folder rules, check and repair. Launching closes it; adding an account or logging in leaves it for the terminal and comes back.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

The first frame draws from `config.toml` alone, in well under 100 ms; usage fills in as each account answers. A percentage turns yellow from 80% and red at 100%, and values from a cache after a failed lookup are drawn faint.

## Tabs

| Tab | What it shows | Keys |
| --- | --- | --- |
| Accounts | number, name, email, 5-hour bar, 7-day usage, last use, groups; `✔` the default, `▸` this folder's account, `•` running | below |
| Groups | each group, what it shares and its accounts, and the accounts in no group | `enter` edit, `n` new, `e` rename, `x` remove |
| Folders | the [`ccam folder`](/features/folder) rules and the one that applies here | `n` use an account here, `enter` change, `x` remove |
| Extensions | what [`ccam ext`](/features/ext) shows, each extension once with the groups it comes from | `r` reload |
| Health | what [`ccam doctor`](/features/doctor) reports | `f` repair (`doctor --fix`) |

`←`/`→` or `Tab` switch tabs, `?` lists the keys of the current tab, `r` looks usage up again without the short cache, and `q` or `Esc` quits.

## Keys on the Accounts tab

| Key | What it does |
| --- | --- |
| `Enter` | Start Claude Code as the selected account |
| `1`–`9` | Move to the account with that number; type `1` `2` quickly for 12. Only `Enter` launches |
| `↑`/`↓`, `k`/`j` | Move |
| `n` | Start the account [`ccam next --auto`](/features/next) would choose |
| `o` | Start with options: `--continue`, `--resume`, other arguments, or a dry run |
| `u` | Usage details: 5-hour, 7-day and model windows with their resets and where the numbers come from |
| `a` | Add an account (the wizard below) |
| `l` | Log in again |
| `s` | Choose the account's groups |
| `d` | Make it the default account |
| `f` | Use it in the current folder |
| `e` | Rename it |
| `c` | Set or remove its shortcut command |
| `x` | Remove it (asks first, starting on No) |

## Changing groups

On the Groups tab, `Enter` opens the editor: the items on the left, the accounts on the right, `Tab` to switch, `Space` to toggle. An item that one of the accounts already gets from another group is greyed out with the reason, as is an account that would get an item twice; an item can come from one group only. `Enter` shows the plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

A file that differs from the group's copy is asked about first (keep the group's, use mine, or merge). Then the change runs in the background with the same code as [`ccam group`](/features/share): the config is saved first, each account's data moves in order, and the originals go to `.ccam-backup/`. A change that would move the data of a running account is refused with the account's name.

## Adding an account

`a` opens the wizard that `ccam add` also uses on a terminal:

1. The account name.
2. How it signs in: a Claude subscription now (in the browser), an [API key](/features/api), or a Claude subscription later.
3. For an API key: where ccam reads the key at each launch (`env:NAME` or `command:COMMAND`), the endpoint, how the key is sent and the model.
4. The groups to join, with `default` checked; before there are groups, what group `default` shares.
5. A shortcut command, if you want one.
6. A summary to confirm.

`Esc` goes back one step, and on the first step closes the wizard.

## Colors

The colors follow the roles of Claude Code's terminal UI, so ccam looks at home next to it; group names get a color from their name. ccam asks the terminal for its background and picks a light or dark theme. Set one yourself in `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

The `-daltonized` themes are for red-green color blindness and the `-ansi` themes use the terminal's own 16 colors. `NO_COLOR` turns colors off.

## Safety and edge cases

- The dashboard calls the same functions as the commands, so every rule holds: plans before data moves, refusals for running accounts, backups, and login tokens that never leave an account.
- The dashboard has no dry run: `ccam --dry-run` without a name prints how this folder's account would start, and `o` on an account offers a dry run of that one launch.
- The terminal is restored before Claude Code starts or the login opens.
- A screen too small for a dialog drops the tab's rows first; lines are cut with `…` rather than wrapped.

## Related

- [Sharing with groups](/features/share)
- [Launching Claude Code](/features/launch)
- [Managing accounts](/features/accounts)
- [Usage limits](/features/usage)
