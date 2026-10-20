---
title: Picking an account for scripts
description: ccam pick prints the account with the most room for a script, gives workers started together different accounts, and exits 75 when none has room.
updated: 2026-10-20
order: 12
translation: source
---

`ccam pick` prints the name of the logged-in account with the most room, for a script to start Claude Code with. Workers started at the same time each get a different account while there are enough accounts with room, so a batch of `claude -p` jobs spreads over your accounts instead of draining one.

## Usage

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live"
}
```

Exit codes: `0` with a name on standard output, `75` (`EX_TEMPFAIL`) when no account has room or none is logged in, `1` for other errors such as a bad flag.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--json` | Print the account, its 5-hour and weekly usage and where the numbers come from | off |
| `--model <model>` | Count that model's weekly window too, such as `opus` | none |
| `--max-usage <1-100>` | Prefer accounts below this usage | `auto_max_usage`, else 100 |
| `--strategy headroom\|expiring` | Lowest 5-hour usage, or the soonest 5-hour reset | `auto_strategy`, else `headroom` |
| `--wait[=<duration>]` | When no account has room, wait for the first one to reset, at most this long | off; `5h` when given alone |
| `--dry-run` | Choose without recording the pick | off |

## How it works

1. ccam lists the logged-in accounts (API key accounts have no limits and are left out).
2. It takes a lock (`~/.local/state/ccam/picks.lock`), so picks running at the same time happen one after another.
3. It looks up every account's usage at once, the same way `ccam usage` does, with the 60-second cache.
4. Accounts picked within the last 2 minutes (`~/.local/state/ccam/picks/<name>`) are left out first, and the [choice rules](/features/next#auto) of `ccam next --auto` pick among the rest: below `auto_max_usage` first, by the strategy, never an account at 100%.
5. If every account with room was picked recently, the rules choose among all of them again, so a fourth worker on three accounts shares one.
6. The pick is recorded and the lock released, then the name is printed.

The 2-minute window matters because a new job's usage takes a while to show in the numbers: without it, five workers started together would all see the same account as the one with the most room.

## Safety and edge cases

- `ccam pick` only reads usage and writes the pick time; it never renews or reads tokens beyond what `ccam usage` does.
- A pick does not reserve the account. Two scripts may still start on the same account after 2 minutes; use `--fallback` to recover when that account hits its limit.
- Waiting with `--wait` holds no lock, so other picks go on meanwhile.
- `ccam pick` ignores [`ccam folder`](/features/folder): it chooses by usage only.

## Related

- [Switching to the next account](/features/next)
- [Run again at the limit](/features/next#run-again-at-the-limit)
- [Usage limits](/features/usage)
- [API key accounts](/features/api)
