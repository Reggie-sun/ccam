---
title: Switching to the next account
description: ccam next opens Claude Code as the account after the one you used last; --auto picks the most room, --continue resumes, and --fallback reruns -p jobs.
updated: 2026-10-20
order: 8
translation: source
---

`ccam next` opens Claude Code as another account without you naming it. Use it when an account hits its usage limit: plain `ccam next` takes the account after the one you used last, `--auto` takes the one with the most room in its 5-hour window, and `--continue` carries on the last conversation in the new account. For scripts that run `claude -p`, `ccam --fallback` makes that switch by itself (see [Run again at the limit](#run-again-at-the-limit)).

## Usage

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam prints the switch as `→ <from> → <to>` on standard error, then starts Claude Code as the new account the same way `ccam <name>` does. `<from>` is the account you used last. Arguments after `--` go to `claude` unchanged.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--auto` | Pick the logged-in account with the lowest 5-hour usage, skipping accounts at 100% | off |
| `-c`, `--continue` | Only consider accounts that share conversation history with the account you used last, and start `claude --continue` | off |
| `-- <claude args>` | Arguments passed to `claude` after ccam's own | none |
| `--dry-run` | Honoured: prints `→ <from> → <to>` and the command line it would run, without starting Claude Code or recording a launch | off |

`--auto` and `--continue` work together. With `--continue`, `--continue` is placed before your own arguments.

## How it works

Every launch through ccam writes the time to `~/.local/state/ccam/last-used/<name>`. `ccam next` reads these files, checks which accounts are logged in, and chooses.

1. `cli.candidates` builds one `next.Candidate` per account from `config.toml`: the account, its last launch time (`launcher.LastUsed`) and whether it has a login (`creds.Has`: a readable `.credentials.json`, or on macOS a Keychain item).
2. Without `--auto`, `next.Rotate` sorts the accounts by slot number and finds the one launched most recently. That is the account you are switching away from. If no account has been launched yet, it is the lowest slot.
3. From there it walks forward in slot order, wrapping from the last slot back to the first, and returns the first account that is logged in.
4. `launcher.LaunchWith` prepares and starts the chosen account exactly as `ccam <name>` would: it checks the account's links, applies MCP sharing, registers the session and replaces ccam with `claude` through `syscall.Exec`.

The walk starts at the account used last and stops at the first logged-in account after it:

<!-- motion:next-rotate -->

`work` was launched most recently, so the walk starts after it. `side` has no login and is skipped, and `lab` is chosen. Had `lab` been logged out too, the walk would wrap around to `personal` (`% len` in the diagram). Rotation never returns the account you started from: with no other logged-in account, it stops with an error.

### --auto

`--auto` decides by usage instead of slot order. `cli.pickAuto` keeps the logged-in accounts, asks for all their usage at once with the same `Fetcher.FetchAll` that `ccam usage` uses, and lets `usage.Best()` choose:

<!-- motion:next-auto -->

`usage.Best()` takes the lowest 5-hour percentage and skips accounts at 100% or with unknown 5-hour usage. The lookup follows the rules of `ccam usage`: numbers younger than 60 seconds are reused, each account has 3 seconds, and an account that does not answer, or whose login has expired, is judged by its saved numbers. A saved window whose reset time has passed counts as 0%, so an account whose window has reset since is not passed over. `--auto` never renews tokens. Unlike plain rotation, `--auto` can choose the account you used last if it still has the most room.

### --continue

Claude Code's `--continue` opens the most recent conversation in the current folder, which only works if the new account can see that conversation. So with `--continue`, ccam only considers accounts in the same `history` group as the account you used last, then puts `--continue` in front of your arguments. Accounts whose history is `local` share with nobody, so `--continue` from such an account always fails. Turn sharing on first with `ccam share <name> --history on` (see [sharing](/docs/sharing)).

## Run again at the limit

`ccam next` is for an interactive session that has stopped. A non-interactive run, `claude -p` in a script or a scheduled job, cannot wait for you, so `--fallback` before the account name lets ccam switch for it:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` works only together with `-p` or `--print`; otherwise ccam stops and points to `ccam next --auto --continue`.
2. Standard input is read once and given to every attempt.
3. ccam starts `claude` as a child process instead of `exec`, holds its standard output back, and passes standard error through as it comes.
4. When `claude` exits with 0, or fails for another reason, the held output and the exit code are passed on unchanged.
5. A failure counts as the usage limit when Claude Code's message says so (`You've hit your limit`, `usage limit reached`), or when the account's usage lookup shows its 5-hour or weekly window at 100%.
6. Then the output of that attempt is dropped and the same command runs as another logged-in account: the one with the lowest 5-hour usage, then accounts whose usage cannot be looked up, in number order. Accounts at 100% are skipped, and each account is tried at most once.
7. With `-c`, `--continue`, `-r` or `--resume` among the arguments, only accounts in the same history group are tried, so the conversation can be found.

When no account is left, ccam prints the last attempt's output, says that no logged-in account has room left, and exits with that attempt's code. `ccam --dry-run --fallback work -p …` prints the command of the first attempt.

Interactive sessions are not switched on their own. Instead, with `ccam statusline install`, the status line adds `→ ccam next --auto --continue` once the 5-hour or weekly usage reaches 90% ([Usage limits](/features/usage)).

## Safety and edge cases

- Accounts without a login are never chosen. With no other logged-in account, plain `ccam next` stops with `no logged-in account other than <name>`.
- With `--continue` and no account sharing history: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- With `--auto` and no logged-in account: `no logged-in account to pick`. When every account is full or its usage is unknown: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` makes network requests even with `--dry-run`. It only reads, and never writes token files.
- The chosen account starts through the same path as `ccam <name>`, so the same checks apply, including the warning when another Claude Code sharing the same history is already running in this folder.
- `last-used` is only written when Claude Code actually starts. `--dry-run` leaves it alone, so a dry run does not move the rotation forward.
- `--fallback` holds standard output until an attempt ends, so streaming output such as `--output-format stream-json` arrives in one piece at the end. Files the first attempt changed before it hit the limit stay changed; the next attempt starts the prompt again.

## Related

- [Usage limits](/features/usage)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Share MCP servers, settings and history](/docs/sharing)
- [Usage limit workflow](/guides/usage-limit-workflow)
