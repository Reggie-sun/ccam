---
title: Switching to the next account
description: ccam next opens Claude Code as another account; --auto picks the most room, --continue carries the conversation over, and --fallback reruns -p jobs.
updated: 2026-10-25
order: 8
translation: source
---

`ccam next` opens Claude Code as another account without you naming it. Use it when an account hits its usage limit: plain `ccam next` takes the account after the one you used last, `--auto` takes the one with the most room, and `--continue` carries this folder's last conversation over to the new account, whether or not the two share history. Inside Claude Code, the [ccam plugin](#inside-claude-code) tells you when a limit stops a session and `/ccam:switch` hands the conversation on. For scripts that run `claude -p`, `ccam --fallback` makes the switch by itself (see [Run again at the limit](#run-again-at-the-limit)), and with [`auto_switch`](#auto-switch) turned on (Linux and macOS), ccam also moves an interactive session that Claude Code stops at its limit.

**Team and Enterprise seats.** Automatic choice (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n`, and [`auto_switch`](/features/next#auto-switch) when you turn it on) also switches between seats of the same Claude Team or Enterprise organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Using automatic choice between seats may therefore violate Anthropic's terms. Whether you use it there is your call: check with your organization and Anthropic before you rely on it. Personal plans (Pro, Max) are covered by the [terms note in the FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Usage

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next -c --brief            # a new conversation that starts from a summary of the last one
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
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam prints the switch as `→ <from> → <to>` on standard error, with the reason for an `--auto` choice, then starts Claude Code as the new account the same way `ccam <name>` does. `<from>` is the account you used last, or with `--continue` the account the conversation belongs to. Arguments after `--` go to `claude` unchanged.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--auto` | Pick the logged-in account by usage (see below), skipping accounts at 100% | off |
| `--max-usage <1-100>` | With `--auto`: prefer accounts below this usage | `auto_max_usage`, else 100 |
| `--strategy headroom\|expiring\|balance` | With `--auto`: lowest 5-hour usage, the soonest 5-hour reset, or the weekly quota that would go unused first | `auto_strategy`, else `headroom` |
| `--cross-org` | With `--auto`: also choose accounts of another organization (see below) | `cross_org`, else off |
| `--wait[=<duration>]` | With `--auto`: when no account has room, wait for the first one to reset, at most this long. The value goes after `=`: `--wait 2h` is refused with the hint `--wait takes its value after '=': --wait=2h` | off; `5h` when given alone |
| `-c`, `--continue` | Carry this folder's last conversation on as the new account (see below) | off |
| `--brief` | With `--continue`: start a new conversation as the new account from ccam's summary of the last one, instead of resuming all of it ([see below](#brief)) | off |
| `-- <claude args>` | Arguments passed to `claude` after ccam's own | none |
| `--dry-run` | Honoured: prints `→ <from> → <to>` and the command line it would run, without starting Claude Code or recording a launch | off |

`--auto` and `--continue` work together. With `--continue`, the `--resume` arguments are placed before your own.

Claude Code's own options, such as `--dangerously-skip-permissions` or `--model`, go after `--`. Unlike `ccam <name>`, which passes everything after the name on, `ccam next` reads its arguments itself, so without the `--` it stops with `ccam: unknown flag: --dangerously-skip-permissions`.

## How it works

Every launch through ccam writes the time to `~/.local/state/ccam/last-used/<name>`. `ccam next` reads these files, checks which accounts are logged in, and chooses.

1. `cli.candidates` builds one `next.Candidate` per account from `config.toml`: the account, its last launch time (`launcher.LastUsed`) and whether it has a working login (`creds.Has`: a Claude login in `.credentials.json`, or on macOS in the Keychain, that has not expired as a whole or been cleared by Claude Code). A folder that uses an API key instead (a Console login, or an account added with `ccam add --api`) has no plan limits and is never chosen.
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

A window counts as 100% until its reset when Claude Code refused a request there: the [ccam plugin](#inside-claude-code)'s hook records it (`~/.local/state/ccam/usage/<name>.limit.json`), whatever the saved numbers say. So an account that has just hit its limit is not chosen again before its reset, even when its numbers are old.

Settings for all accounts go at the top of `config.toml`, above the first `[[group]]`, `[[profile]]` or `[[folder]]`. A setting below one of them belongs to that table: ccam names it at every launch, and commands that save the config refuse until it is moved.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizations.** Automatic choices stay with accounts like the one you leave: the same team or enterprise organization, or personal plans (Pro, Max) together, so work from a company account does not continue on a private one or the other way round. ccam knows each account's organization from its login (`ccam login` records it; older logins are read from `.claude.json`); an account it cannot tell fits anywhere. Two folders logged in to the same Claude account share its limits, so one never takes over from the other. `--cross-org`, or `cross_org = true` in `config.toml`, lifts the organization rule.

When nothing can be chosen, ccam exits with code 75 (`EX_TEMPFAIL`), so a script can tell "try later" from a failure. With `--wait`, it instead sleeps until the first full account's blocking window resets (30 seconds after the reset time, when the endpoint shows the new window), prints `ccam: no account has room; waiting until 15:10`, and chooses again; when no account says when it resets, it looks again after 10 minutes. If nothing frees up before the wait ends, it exits with 75. The lookup follows the rules of `ccam usage`: numbers younger than 60 seconds are reused, each account has 3 seconds, an account that does not answer, or whose access token has expired, is judged by its saved numbers, and a limit record lies over them. A saved window whose reset time has passed counts as 0%, so an account whose window has reset since is not passed over. `--auto` never renews tokens. Unlike plain rotation, `--auto` can choose the account you used last if it still has the most room.

### --continue

1. ccam looks for this folder's last conversation in every account: of the `<session>.jsonl` files in `projects/<folder>/` (the folder name is the working folder, links resolved, with every character other than a letter or digit turned into `-`, as Claude Code names it), the one with the newest prompt or reply. Sessions without one, such as a session that only ran `/mcp`, `/usage` or `/status`, or an empty file, are skipped. If a session left a note (a usage limit or `/ccam:switch`, see below) and no conversation has a prompt or reply more than 2 minutes after it, that exact session is used instead; the rows Claude Code writes at `/exit` do not count.
2. The account it belongs to is the one you switch from, and the new account is chosen as above, within the same organization. When several accounts share that history, it belongs to the one whose `.claude.json` says it ran the conversation last (`lastSessionId`), else to the one used last.
3. If the new account shares that history (the same group folder), ccam starts `claude --resume <session id>`, that very conversation (`--continue` would open whichever is newest in the folder).
4. Otherwise it starts `claude --resume <the conversation's file> --fork-session`: Claude Code copies the whole conversation, under a new session id, into the new account's own history, and the original stays as it was. From then on the new account's `--continue` finds it.

ccam says how big the conversation is: the new account reads all of it again with its first reply. Prompt caches are never shared between organizations, and every personal plan is an organization of its own, so when the two accounts belong to different ones ccam adds `prompt caches are not shared between organizations, so the first reply reads it all again`. When the old account is signed in to MCP servers the new one is not, ccam names them: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. Sharing `history` in a group avoids the copy, keeps one conversation, and also shares the checkpoints `/rewind` restores code from (`file-history/`).

### --brief

`--continue` resumes the whole conversation, and the new account reads all of it again with its first reply. `ccam next --continue --brief` starts a new conversation instead, from a summary that ccam writes of the old one without asking a model:

```console
$ ccam next --auto --continue --brief
→ personal → work (12% of its 5 hours used, the least)
  starting a new conversation with ccam's summary of "Fix the flaky tests" instead of resuming it (about 182k tokens)
  summary (6k characters): /home/you/.local/state/ccam/brief/4f1c….md
```

The summary holds, in this order and only where there is something to say: the conversation's title; the path of its transcript, for Claude to search when it needs a detail; your first request, when there were more than four; the summary Claude Code wrote when it last compacted the conversation; your last four requests, with `!` commands and slash commands as you typed them; Claude's last reply; the todo list; the 10 files Claude changed last; and the folder's `git status --short --branch`, up to 800 characters. The todo list keeps at most 10 items: finished items go first, then the steps furthest ahead, so the current step and the next ones stay. Long parts are cut in the middle, and the whole summary stays under 16,000 characters.

ccam saves it as `~/.local/state/ccam/brief/<session>.md`, readable only by you, and the next `--brief` removes summaries older than 7 days. Claude Code starts with one prompt, and attaches the file it mentions: `Read @"<the summary>", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here. Then reply with one or two sentences on where the work stands, and wait for my next message.` A prompt you give after `--` takes the place of the last sentence (`Then: <your prompt>`), and your other options after `--` follow it. Claude Code ends a mention at a `"` and cuts it at a `#`, so when the summary's path has either, ccam mentions it by its path from the working folder; when that path has one too, the prompt names the file, and Claude reads it with its Read tool, which asks first because the file is outside the working folders (in a `-p` run in Manual mode, Claude Code denies the read and Claude goes on without the summary).

`ccam --dry-run next --continue --brief` prints the summary and saves nothing. The old conversation stays as it was, and the new one knows only what the summary says, so use `--brief` when the task matters more than every detail of how you got there. Without `--continue`, ccam stops with `--brief works with --continue`; when it cannot read the conversation or finds nothing in it to summarize, it stops and says that `ccam next --continue` without `--brief` resumes it whole. The [relay](#relay) and [`auto_switch`](#auto-switch) always resume the whole conversation. Claude Code's own `--brief` option goes after `--`, like its other options.

## Inside Claude Code

The [ccam plugin](/docs/install#claude-code-plugin) adds two hooks (Claude Code 2.1.78 or later):

- When a turn stops at a usage limit, it notes the conversation and shows a desktop notification in terminals that support one (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` It knows every message Claude Code counts as a usage limit: the 5-hour, weekly, plain and model limits (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` for a model's allowance), usage credits and spend limits, and an organization that stopped the account. When Claude Code has written the refused request into the conversation, the hook also records the full window until its reset ([the limit record](/features/usage#the-limit-record)), except for an organization's stop, which its admins can lift at any time. A short server throttle or an API key's rate limit is not a usage limit and is left alone, and so are a subagent's stop and sessions nobody sits at (`-p`, background sessions). Claude Code 2.1.234 and later also waits in the session and continues by itself after a 5-hour reset; the notification gives you the choice to switch instead, and with [`auto_switch`](#auto-switch) on (Linux and macOS), ccam switches for you when another account has room. When no other account may take over, because no other account belongs to the same Team or Enterprise organization, or none is another personal plan (with `cross_org` on, any other Claude account may), the notification ends with `You have no other personal plans to carry it on as.` or `You have no other accounts in Acme to carry it on as.` instead.
- `/ccam:switch [account]` notes the conversation and ends the turn before any request is made: type `/exit`, then `ccam next --continue` carries it on, as the account you named or as the next one. In a session that has no conversation yet, it says there is nothing to carry on.

The note lives in `~/.local/state/ccam/handoff/<account>/<session>.json`, one file per session, so two sessions of one account never overwrite each other's. It is kept for 12 hours and used once, when the next session is about to start; a launch that fails leaves it for the next try. It belongs to the folder the session started in, even after a `cd` in Claude Code's shell. `ccam rename` updates the names in it; if the account it names has since been removed, ccam says so and chooses as usual.

### Relay

With `relay = true` at the top of `config.toml`, ccam stays beside interactive sessions instead of replacing itself with Claude Code, and after `/ccam:switch` it carries that session's conversation on by itself (a note another session left is not acted on): the hook ends Claude Code, and ccam runs `ccam next --continue` with that session's note. Without an account named, it takes the one with the most room, counting the conversation's model. On Windows the hook cannot end Claude Code: type `/exit`, and ccam carries on. That costs one waiting ccam process per session, which is why it is off by default.

The relayed session keeps what resuming a conversation does not restore: options such as `--add-dir`, `--mcp-config`, the tool lists, `--settings`, `--plugin-dir` and the system prompt options, the conversation's permission mode and the session's effort level. Claude Code itself restores `--model`, `--name` and `--agent` from the conversation. Anything else (`--worktree`, `--debug`, `--remote-control`, options ccam does not know) is named in `ccam: not carried over: …`. The switch never asks the shared-history question, and what you typed while it ran is dropped.

The relay carries on only when Claude Code ended the way `/ccam:switch` asked it to: not when you typed `/exit` first, pressed Ctrl+C while it ended, or ccam got a hangup or another signal, and not when the terminal is gone. The note then stays, and `ccam next --continue` carries the conversation on later. A session started under an older ccam's relay, or one nobody sits at (`-p`, background sessions), is not ended by the hook; it gets the `/exit` message instead. `claude remote-control` runs without a relay.

### `auto_switch`

With `auto_switch = true` at the top of `config.toml` (Linux and macOS), ccam carries a session on by itself when Claude Code stops it at a usage limit: the plugin's hook picks another account with room and ends Claude Code, and ccam, waiting beside the session, continues the conversation there, the way the [relay](#relay) does after `/ccam:switch`. It is off by default and experimental.

```toml
# ~/.config/ccam/config.toml
auto_switch = true   # Linux and macOS; turns the relay on too
```

```console
ccam: personal hit its 5-hour limit; continuing this conversation as side
→ personal → side
  continuing "Fix the flaky tests" (about 182k tokens); prompt caches are not shared between organizations, so the first reply reads it all again; side does not share personal's history, so Claude Code copies it into side
```

ccam moves a session only when all of these hold:

1. ccam started it in a terminal and waits beside it. `auto_switch` turns the [relay](#relay) on, so this holds for the interactive sessions you start after setting it. A session started before, or under an older ccam, gets a notification that ends ``This session started before auto_switch was on or under an older ccam, so ccam cannot move it: type /exit, then run `ccam next --continue`.`` Sessions nobody sits at (`-p`, background sessions), `claude remote-control` and Claude Code started without ccam are never moved.
2. The account has the [ccam plugin](/docs/install#claude-code-plugin) turned on and runs Claude Code 2.1.214 or later, which tells the hook its process.
3. Claude Code wrote the refused request into the conversation, with the time it resets, for that account's 5-hour, weekly, Opus or Sonnet limit. A spend limit, usage credits running out, an organization's stop, the Fable allowance (`You've reached your Fable limit.`), a server throttle, an API key's rate limit and a subagent's stop never move a session.
4. ccam knows the account's plan, from its `.claude.json` or from what `ccam login` recorded. If not, the notification ends `ccam does not know personal's plan, so auto_switch leaves this session.`
5. The session is not connected to Remote Control. As another account it would not reconnect to the app you drive it from, so it stays, and the notification says so.
6. The limit resets more than 15 minutes from now. For a shorter wait ccam leaves the session to Claude Code: `ccam: personal hit its 5-hour limit, which resets in 9m; ccam leaves this session to wait for it.`
7. You have not typed a prompt since the limit.
8. Another account has room. ccam chooses as [`--auto`](#auto) does (`auto_max_usage`, `auto_strategy`, the conversation's model), among accounts like the one you leave: another seat of the same Team or Enterprise organization, or another personal plan (Pro, Max). `cross_org` and `--cross-org` do not apply here. The account must be logged in with a login Claude Code can renew, must not be logged in to the same Claude account, must not use an API key, must have a plan ccam knows, must not have usage credits on as far as ccam's own usage lookup shows, and must have trusted this folder in Claude Code. When those last two leave out every account with room, the notification names them: `side has not trusted this folder yet, so the conversation was not moved there. Open the folder once as side to allow it.` or `side has room but its usage credits are on, so auto_switch did not move the conversation there; /ccam:switch side moves it.` A window with a [limit record](/features/usage#the-limit-record) counts as full, and an account another session moved to in the last 2 minutes comes after the others. The choice has to be made within about 6 seconds; when the usage lookups take longer, the session stays. With `usage_source = "claude"`, `auto_switch` decides from saved numbers alone, and ccam does not read whether usage credits are on from Claude Code's lookup or the status line, so it counts them as off.

Then the hook ends Claude Code, and ccam, beside it, shows a desktop notification in terminals that support one, such as `ccam: personal hit its 5-hour limit (resets 3:29am). Moving this conversation to side (12% of its 5 hours used, the least).` It prints the first line above and runs `ccam next --continue` for that session as the chosen account, which prints the other two and starts Claude Code with the session's options, permission mode and effort level, as the [relay](#relay) does. The new session's first message is this fixed text, the only thing ccam sends in your place (a move after `/ccam:switch` sends nothing):

```text
Your usage limit was reached and ccam moved this conversation to another account. Continue the task you were working on when the limit was reached; do not repeat work that is already complete.
```

If that account reaches its limit too, the same happens there. The limit records keep the chain from going back to an account before its reset, so when every account is full the session stays on the last one.

When no other account has room, nothing moves: Claude Code handles the limit as usual, waiting in the session or offering its options, and the notification says which account frees first, for example `ccam: personal hit its 5-hour limit (resets 3:29am). None of your other personal plans has room now; side frees first, at 12:30am. Claude Code waits for personal here; /ccam:switch side after 12:30am moves it sooner.` For seats it says `None of your other accounts in Acme has room now`. To stop a move that has started, press Ctrl+C while Claude Code closes: ``ccam: switch cancelled. `ccam next --continue` carries the conversation on later.``

A move has costs. Claude Code restarts, which takes a few seconds, and the background shells, dev servers and agents that Claude started stop with the old session. The new account reads the whole conversation again with its first reply, as with [`--continue`](#continue), and ccam says how big it is; a move never uses [`--brief`](#brief).

`auto_switch` never switches before Claude Code reports a limit, never watches your usage in the background, and never copies, renews or passes on a login. On Windows it is ignored, and `ccam doctor` says so; there `/ccam:switch` and `/exit` still carry a conversation on with `relay = true`.

**Anthropic's terms.** On Team and Enterprise seats, `auto_switch` moves your work to another seat of the same organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). For personal plans (Pro, Max), Anthropic's terms do not say whether your own accounts may take over from each other automatically at a limit. The [Consumer Terms](https://www.anthropic.com/legal/consumer-terms) forbid accessing the Services "through automated or non-human means, whether through a bot, script, or otherwise" unless Anthropic explicitly permits it, and a move sends the fixed message above in your place. The [Usage Policy](https://www.anthropic.com/legal/aup) that takes effect on November 12, 2026 lists "Intentionally bypass or interfere with capabilities, restrictions, or guardrails established within our products" among the things it prohibits. ccam does not claim that `auto_switch` complies with these terms. Turning it on is your call, and so is the risk: Anthropic may suspend an account, or a Team or Enterprise organization's access, when it believes the terms are broken.

While `auto_switch` is on, `ccam doctor` warns about an account without the ccam plugin and about a Claude Code older than 2.1.214, and for seats repeats Anthropic's seat rule ([checks](/features/doctor)).

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
2. Piped standard input is read once, up to 10 MiB, and given to every attempt. Like Claude Code, ccam waits at most 3 seconds for it to start (not with `--input-format stream-json`); a pipe that stays open and sends nothing gives the run no input, with a line on standard error.
3. ccam starts `claude` as a child process instead of `exec`, holds its standard output back, and passes standard error through as it comes.
4. When `claude` exits with 0, or fails for another reason, the held output and the exit code are passed on unchanged. So is a run ended by a signal (an exit code above 128): it was asked to stop, and is not run again.
5. A failure counts as the usage limit when Claude Code's message says so, any of the usage-limit messages the [hook](#inside-claude-code) knows, also inside `--output-format json` or `stream-json` output, or when the account's usage shows its 5-hour, weekly or `--model` window at 100% (a limit record included).
6. Then the output of that attempt is dropped and the same command runs as another logged-in account: the one with the lowest 5-hour usage, then accounts whose usage cannot be looked up, in number order. Accounts at 100% are skipped, accounts with an API key login are never tried, and each account is tried at most once.
7. With `-c`, `--continue`, `-r` or `--resume` among the arguments, only accounts in the same history group are tried, so the conversation can be found.
8. Only accounts of the same organization are tried, as with `--auto` (`cross_org = true` lifts this).

The next account is chosen by the same rules as `--auto` (`auto_max_usage`, `auto_strategy`, the `--model` window). When no account is left, ccam prints the last attempt's output, says that no logged-in account has room left, and exits with code 75. With `--wait` before the account name (`ccam --fallback --wait=2h work -p …`), it waits for the first account to reset instead, then may try every account again; it always sleeps at least once first, so an account that has just hit its limit is not run again at once. The duration bounds all the waiting of the run together, not each wait. Written with a space (`--wait 30m`), it is refused with the same hint as `ccam next`. A leading `--` before Claude Code's arguments works as for any launch (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` prints the command of the first attempt.

Interactive sessions are not switched on their own unless you turn on [`auto_switch`](#auto-switch). With `ccam statusline install`, the status line adds `→ ccam next --auto --continue` once the 5-hour or weekly usage reaches 90% ([Usage limits](/features/usage)).

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
- `--fallback` holds standard output until an attempt ends, so streaming output such as `--output-format stream-json` arrives in one piece at the end. A client that streams `stream-json` input and waits for replies before it sends more cannot use `--fallback`. Files the first attempt changed before it hit the limit stay changed; the next attempt starts the prompt again.

## Related

- [Usage limits](/features/usage)
- [Usage limits and switching accounts](/docs/usage-limits)
- [Share MCP servers, settings and history](/docs/sharing)
- [Usage limit workflow](/guides/usage-limit-workflow)
