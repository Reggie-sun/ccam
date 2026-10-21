---
title: Switching to the next account
description: ccam next opens Claude Code as another account; --auto picks the most room, --continue carries the conversation over, and --fallback reruns -p jobs.
updated: 2026-10-21
order: 8
translation: source
---

`ccam next` opens Claude Code as another account without you naming it. Use it when an account hits its usage limit: plain `ccam next` takes the account after the one you used last, `--auto` takes the one with the most room, and `--continue` carries this folder's last conversation over to the new account, whether or not the two share history. Inside Claude Code, the [ccam plugin](#inside-claude-code) tells you when a limit stops a session and `/ccam:switch` hands the conversation on. For scripts that run `claude -p`, `ccam --fallback` makes the switch by itself (see [Run again at the limit](#run-again-at-the-limit)).

## Usage

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam prints the switch as `→ <from> → <to>` on standard error, with the reason for an `--auto` choice, then starts Claude Code as the new account the same way `ccam <name>` does. `<from>` is the account you used last, or with `--continue` the account the conversation belongs to. Arguments after `--` go to `claude` unchanged.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--auto` | Pick the logged-in account by usage (see below), skipping accounts at 100% | off |
| `--max-usage <1-100>` | With `--auto`: prefer accounts below this usage | `auto_max_usage`, else 100 |
| `--strategy headroom\|expiring\|balance` | With `--auto`: lowest 5-hour usage, the soonest 5-hour reset, or the weekly quota that would go unused first | `auto_strategy`, else `headroom` |
| `--cross-org` | With `--auto`: also choose accounts of another organization (see below) | `cross_org`, else off |
| `--wait[=<duration>]` | With `--auto`: when no account has room, wait for the first one to reset, at most this long | off; `5h` when given alone |
| `-c`, `--continue` | Carry this folder's last conversation on as the new account (see below) | off |
| `-- <claude args>` | Arguments passed to `claude` after ccam's own | none |
| `--dry-run` | Honoured: prints `→ <from> → <to>` and the command line it would run, without starting Claude Code or recording a launch | off |

`--auto` and `--continue` work together. With `--continue`, the `--resume` arguments are placed before your own.

Claude Code's own options, such as `--dangerously-skip-permissions` or `--model`, go after `--`. Unlike `ccam <name>`, which passes everything after the name on, `ccam next` reads its arguments itself, so without the `--` it stops with `ccam: unknown flag: --dangerously-skip-permissions`.

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

`--auto` decides by usage instead of slot order. `cli.pickAuto` keeps the logged-in accounts, asks for all their usage at once with the same `Fetcher.FetchAll` that `ccam usage` uses, and lets `usage.Policy.Choose()` choose:

<!-- motion:next-auto -->

An account is judged by its 5-hour and weekly windows, and, when `--model` is among the arguments for Claude Code (after `--`), by that model's weekly window too (`opus`, `sonnet`, `haiku` or `fable` in the name). `Policy.Choose()` sorts the accounts in three:

1. **Room**: every window below `auto_max_usage` (`config.toml`, default 100). The strategy chooses among them: `headroom` (default) takes the lowest 5-hour usage; `expiring` takes the 5-hour window that resets soonest, so its remaining quota is used before it is lost; `balance` takes the account with the most weekly quota left per day until its week resets, counting each running ccam session against it, so weekly quota is not left unused and parallel sessions spread out.
2. **Over**: below 100% but at or over `auto_max_usage`. Used only when no account has room, the lowest 5-hour usage first.
3. **Full**: a window at 100%, or unknown 5-hour usage. Never chosen.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizations.** Automatic choices stay with accounts like the one you leave: the same team or enterprise organization, or personal plans (Pro, Max) together, so work from a company account does not continue on a private one or the other way round. ccam knows each account's organization from its login (`ccam login` records it; older logins are read from `.claude.json`); an account it cannot tell fits anywhere. Two folders logged in to the same Claude account share its limits, so one never takes over from the other. `--cross-org`, or `cross_org = true` in `config.toml`, lifts the organization rule.

When nothing can be chosen, ccam exits with code 75 (`EX_TEMPFAIL`), so a script can tell "try later" from a failure. With `--wait`, it instead sleeps until the first full account's blocking window resets (30 seconds after the reset time, when the endpoint shows the new window), prints `ccam: no account has room; waiting until 15:10`, and chooses again; when no account says when it resets, it looks again after 10 minutes. If nothing frees up before the wait ends, it exits with 75. The lookup follows the rules of `ccam usage`: numbers younger than 60 seconds are reused, each account has 3 seconds, and an account that does not answer, or whose login has expired, is judged by its saved numbers. A saved window whose reset time has passed counts as 0%, so an account whose window has reset since is not passed over. `--auto` never renews tokens. Unlike plain rotation, `--auto` can choose the account you used last if it still has the most room.

### --continue

1. ccam looks for this folder's last conversation in every account: the newest `<session>.jsonl` in `projects/<folder>/` (the folder name is the working folder with every character other than a letter or digit turned into `-`, as Claude Code names it). If the session left a note (a usage limit or `/ccam:switch`, see below) and nothing newer was written since, that exact session is used instead.
2. The account it belongs to is the one you switch from, and the new account is chosen as above, within the same organization.
3. If the new account shares that history (the same group folder), ccam starts `claude --resume <session id>`, that very conversation (`--continue` would open whichever is newest in the folder).
4. Otherwise it starts `claude --resume <the conversation's file> --fork-session`: Claude Code copies the whole conversation, under a new session id, into the new account's own history, and the original stays as it was. From then on the new account's `--continue` finds it.

ccam says how big the conversation is: the new account reads all of it again with its first reply. Prompt caches are never shared between organizations, so with `--cross-org` that first reply costs the full size. Sharing `history` in a group avoids the copy, keeps one conversation, and also shares the checkpoints `/rewind` restores code from (`file-history/`).

## Inside Claude Code

The [ccam plugin](/docs/install#claude-code-plugin) adds two hooks (Claude Code 2.1.78 or later):

- When a turn stops at a usage limit (`You've hit your session limit`, the weekly or a model limit), it notes the conversation and shows a desktop notification in terminals that support one (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` A short server throttle or an API key's rate limit is not a usage limit and is left alone. Claude Code 2.1.234 and later also waits in the session and continues by itself after a 5-hour reset; the notification gives you the choice to switch instead.
- `/ccam:switch [account]` notes the conversation and ends the turn before any request is made: type `/exit`, then `ccam next --continue` carries it on, as the account you named or as the next one.

The note lives in `~/.local/state/ccam/handoff/<account>.json` for 12 hours and is used once, when the next session is about to start; a launch that fails leaves it for the next try. It belongs to the folder the session started in, even after a `cd` in Claude Code's shell. If the account it names has since been renamed or removed, ccam says so and chooses as usual. With `relay = true` in `config.toml`, ccam stays beside interactive sessions instead of replacing itself with Claude Code, and after `/ccam:switch` it carries that session's conversation on by itself (a note another session left is not acted on). That costs one waiting ccam process per session, which is why it is off by default. The relayed session starts with `ccam next --continue`, without the arguments the first one was started with, such as `--model`.

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
8. Only accounts of the same organization are tried, as with `--auto` (`cross_org = true` lifts this).

The next account is chosen by the same rules as `--auto` (`auto_max_usage`, `auto_strategy`, the `--model` window). When no account is left, ccam prints the last attempt's output, says that no logged-in account has room left, and exits with code 75. With `--wait` before the account name (`ccam --fallback --wait=2h work -p …`), it waits for the first account to reset instead, then may try every account again; it always sleeps at least once first, so an account that has just hit its limit is not run again at once. `ccam --dry-run --fallback work -p …` prints the command of the first attempt.

Interactive sessions are not switched on their own. Instead, with `ccam statusline install`, the status line adds `→ ccam next --auto --continue` once the 5-hour or weekly usage reaches 90% ([Usage limits](/features/usage)).

## Safety and edge cases

- Accounts without a login are never chosen. With no other logged-in account, plain `ccam next` stops with `no logged-in account other than <name>`.
- With `--continue` and no conversation in the folder: `no account has a conversation in <folder> to continue`.
- When the organization rule leaves no account: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- With `--auto` and no logged-in account: `no logged-in account to pick`. When every account is full or its usage is unknown: `` no account with known usage and headroom (see `ccam usage`) ``. Both end with exit code 75.
- `--max-usage`, `--strategy`, `--wait` and `--cross-org` work only with `--auto`, and a value outside 1-100 or another strategy is refused.
- A fork copies the conversation, not the account's other state: checkpoints of the original account's `file-history/` do not come along unless history is shared.
- `--auto` makes network requests even with `--dry-run`. It only reads, and never writes token files.
- The chosen account starts through the same path as `ccam <name>`, so the same checks apply, including the warning when another Claude Code sharing the same history is already running in this folder.
- `last-used` is only written when Claude Code actually starts. `--dry-run` leaves it alone, so a dry run does not move the rotation forward.
- `--fallback` holds standard output until an attempt ends, so streaming output such as `--output-format stream-json` arrives in one piece at the end. Files the first attempt changed before it hit the limit stay changed; the next attempt starts the prompt again.

## Related

- [Usage limits](/features/usage)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Share MCP servers, settings and history](/docs/sharing)
- [Usage limit workflow](/guides/usage-limit-workflow)
