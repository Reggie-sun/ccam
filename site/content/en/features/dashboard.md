---
title: The dashboard
description: ccam on a terminal opens a full-screen dashboard of accounts, groups, folder rules, extensions and checks; every change shows its plan first.
updated: 2026-10-24
order: 15
translation: source
---

`ccam` without arguments opens a full-screen dashboard when standard input and output are a terminal; otherwise it prints `ccam list`. The dashboard shows every account with its usage and groups, and does what the commands do: launch, add, log in, change groups, set folder rules, check and repair. Launching closes it; adding an account or logging in leaves it for the terminal and comes back.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

The line under the tab names says in one sentence what the current tab is for. On a wider terminal the Accounts tab also shows each account's email.

Each time it opens, the dashboard asks GitHub in the background whether a newer ccam release exists; the first frame never waits for it, and offline nothing is shown. When there is one, a line under the title recommends it:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` asks first, then updates the way [`ccam update`](/features/update) does: the release is checked against its SHA-256 checksums and replaces ccam in one step. The dashboard you are in is still the old version, so quit and start `ccam` again. When the update cannot be done here, for example for a ccam installed from the `.deb` package, the dashboard shows the same advice as `ccam update` and keeps the offer.

The first frame draws from `config.toml` alone, in well under 100 ms; usage fills in as each account answers. Each limit has its own bar: the 5-hour window, the weekly window for all models, and Fable's own weekly window; `-` means the account has no such window. A percentage turns yellow from 80% and red at 100%, and values from a cache after a failed lookup are drawn faint. A window Claude Code reported full shows 100% until its reset. An expired access token shows a faint `⚠ token expired` (Claude Code renews it at the next start), while `⚠ login expired` means the account needs `ccam login`.

## Tabs

| Tab | What it shows | Keys |
| --- | --- | --- |
| Accounts | number, name, email, 5-hour, weekly and Fable bars, last use, groups; `✔` the default, `▸` this folder's account, `•` running (also Claude Code started without ccam: plain `claude`, an IDE, background sessions and their daemon) | below |
| Groups | each group, what it shares and its accounts, and the accounts in no group | `enter` edit, `n` new, `e` rename, `x` remove |
| Folders | the [`ccam folder`](/features/folder) rules and the one that applies here, also when you reached the folder through a symbolic link (new rules are stored with links resolved) | `n` use an account here, `enter` change, `x` remove |
| Extensions | what [`ccam ext`](/features/ext) shows, each extension once with the groups it comes from, `off:` the accounts where a plugin is not turned on, servers an account has for this folder only among its own, `This folder's project`, and problems that belong to no account | `r` reload |
| Health | what [`ccam doctor`](/features/doctor) reports, as a box per account (below) | `f` repair the selected box, `F` repair all (`doctor --fix`), `↑`/`↓` move between boxes |

`←`/`→` or `Tab` switch tabs, `?` lists the keys of the current tab, `r` looks usage up again without the short cache, and `q` or `Esc` quits.

## Keys on the Accounts tab

| Key | What it does |
| --- | --- |
| `Enter` | Start Claude Code as the selected account |
| `1`–`9` | Move to the account with that number; type `1` `2` quickly for 12. Only `Enter` launches |
| `↑`/`↓`, `k`/`j` | Move |
| `n` | Start the account [`ccam next --auto`](/features/next) would choose |
| `o` | Start with options: `--continue`, `--resume`, other arguments, or a dry run |
| `u` | Usage details: 5-hour, weekly and model windows with their resets and where the numbers come from |
| `a` | Add an account (the wizard below) |
| `l` | Log in again |
| `s` | Choose the account's groups |
| `d` | Make it the default account |
| `f` | Use it in the current folder |
| `e` | Rename it |
| `c` | Set or remove its shortcut command |
| `x` | Remove it (asks first, starting on No) |

**Team and Enterprise seats.** Automatic choice (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n`) also switches between seats of the same Claude Team or Enterprise organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Using automatic choice between seats may therefore violate Anthropic's terms. Check with your organization and Anthropic before you rely on it. Personal plans (Pro, Max) are covered by the [terms note in the FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## The Health tab

The Health tab runs the checks of [`ccam doctor`](/features/doctor) when you open it and draws the same view: a first line with a mark per account, then a box per account. Each row in a box is one area that was checked, `Folder`, `Login` and `Share links` first, marked `✔` OK, `ℹ` information, `⚠` a warning or `✘` an error, so you see at once which account has a problem and where. Findings that belong to no account go in a `General` box.

```console
 ✘ personal   ⚠ work

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

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` move from box to box, and the selected box has an accent-colored border. `f` repairs only the selected box, the way `ccam doctor --fix` would: on an account's box it repairs that account alone, and on the `General` box the findings of no account (plugin paths and shortcuts); every other account is only checked again. It asks first (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), and a box with nothing to repair just says `nothing to repair in <name>`. `F` repairs everything at once, like `ccam doctor --fix`, and `r` checks again.

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

A file that differs from the group's copy is asked about first (keep the group's, use mine, or merge). Then the change runs in the background with the same code as [`ccam group`](/features/share): the config is saved first, each account's data moves in order, and the originals go to `.ccam-backup/`; an account that leaves MCP sharing keeps the group's servers in its `.claude.json`. A change that would move the data of a running account is refused with the account's name, also when the account started while the plan was open, and a plan that another ccam's change of groups made stale is refused with `the groups changed while this plan waited; nothing was moved, make the change again`. While another ccam changes the groups, a new change, rename or remove says `another ccam is changing the groups right now`; after a change that stopped halfway, they wait until `F` on the Health tab (or `ccam doctor --fix`) has finished it.

## Adding an account

`a` opens the wizard that `ccam add` also uses on a terminal:

1. The account name.
2. How it signs in: a Claude subscription now (in the browser), an [API key](/features/api), or a Claude subscription later.
3. For an API key: where ccam reads the key at each launch (`env:NAME` or `command:COMMAND`), the endpoint, how the key is sent and the model.
4. The groups to join, with `default` checked; before there are groups, what group `default` shares.
5. A shortcut command, if you want one.
6. A summary to confirm.

`Esc` goes back one step, and on the first step closes the wizard.

## Removing an account

`x` asks before anything happens, starting on No, and says what the removal does: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` For an API key account there is no login to revoke, and when another account holds the same login it is not logged out. An imported account keeps its folder: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Colors

The colors follow the roles of Claude Code's terminal UI, so ccam looks at home next to it; group names get a color from their name. ccam asks the terminal for its background and picks a light or dark theme. Set one yourself at the top of `config.toml`, above the first `[[group]]`, `[[profile]]` or `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

The `-daltonized` themes are for red-green color blindness and the `-ansi` themes use the terminal's own 16 colors. `NO_COLOR` turns colors off.

## Safety and edge cases

- The dashboard calls the same functions as the commands, so every rule holds: plans before data moves, refusals for running accounts, backups, and login tokens that never leave an account.
- The dashboard has no dry run: `ccam --dry-run` without a name prints how this folder's account would start, and `o` on an account offers a dry run of that one launch.
- The terminal is restored before Claude Code starts or the login opens.
- A screen too small for a dialog drops the tab's rows first; lines are cut with `…` rather than wrapped. On a short screen the blank lines and the tab description go next, so the line of keys stays. Below 40 columns or 8 rows the dashboard shows only `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` while a change, a removal, an update or a repair is being applied waits for it: the status line says `ccam quits when it is done.`
- After `U` has updated ccam, the offer does not come back when the dashboard returns from adding an account or a login.

## Related

- [Sharing with groups](/features/share)
- [Launching Claude Code](/features/launch)
- [Managing accounts](/features/accounts)
- [Usage limits](/features/usage)
