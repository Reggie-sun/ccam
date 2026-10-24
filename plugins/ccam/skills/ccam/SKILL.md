---
name: ccam
description: Use when the user has several Claude Code accounts managed by ccam and asks about usage limits, which account has room, switching accounts after a limit, running parallel `claude -p` workers on different accounts, or why a skill, command or MCP server is missing in one account. Triggers include "which account has room", "I hit my limit", "run these tasks on my other accounts", "spread the work over my accounts", "ccam".
---

# Working with ccam

ccam gives each Claude Code account its own `CLAUDE_CONFIG_DIR` and starts Claude Code as the account you name. It is installed as the `ccam` command. These commands are safe to run from this session; they only read, except where noted.

## See usage

```sh
ccam usage --json     # 5-hour and weekly usage of every account
ccam list --json      # accounts, emails, what they share
```

In the JSON, `limits.five_hour.percent` and `limits.weekly.percent` are 0-100, and `resets_at` says when each window starts over. `status` is `api` for API key accounts, which have no limits.

## When the current account reaches its limit

You cannot move this session to another account yourself; the user can. Tell them to type `/ccam:switch` (or `/ccam:switch <account>`): it notes this conversation and ends the turn. With `relay = true` in ccam's `config.toml`, ccam then ends this session and carries the conversation on as the other account by itself, with the permission mode, effort level and options such as `--add-dir` it had. Otherwise they type `/exit` and run, in this folder:

```sh
ccam next --continue          # as the account /ccam:switch named, else the next one
ccam next --auto --continue   # or the account with the most room
```

When a session stops at a usage limit, the ccam plugin notes it the same way, and `--auto` leaves that account out. If the accounts share conversation history, this very conversation continues; if not, Claude Code copies it into the new account and the original stays as it is.

## Run work on other accounts

For non-interactive work, start a new `claude -p` through ccam. Pick the account with `ccam pick`; calls made close together return different accounts, so parallel workers spread out:

```sh
ccam "$(ccam pick)" -p "fix the failing test in pkg/parse" --output-format json
```

Or let ccam switch when an account hits its limit mid-job:

```sh
ccam --fallback "$(ccam pick)" -p "…" --output-format json
```

- Exit code 75 means no account has room now. Do not retry in a loop; report it, or add `--wait` (`ccam pick --wait`, `ccam --fallback --wait …`), which waits for the first account to reset, up to 5 hours by default.
- `ccam pick --model opus` also counts that model's weekly limit.
- Give each worker its own git worktree when they edit the same repository.
- Output of a `--fallback` run is held until it finishes, so stream-json arrives at the end.

## A skill, command or MCP server is missing in one account

```sh
ccam ext            # what every account loads, and from where
ccam ext work --json
ccam doctor         # broken share links, logins, permissions
```

`ccam ext` shows each skill, command, agent, plugin, MCP server and hook per account, with its source (a share group name, `account`, `local` for one folder only, or `project`), and lists problems such as broken links, skill folders without `SKILL.md`, MCP servers defined differently in two places or only in one account, and servers another account is signed in to. Suggest `ccam doctor --fix` for broken share links and `ccam mcp import <account> <server>` to share a server the account added on its own.

- Plain `claude mcp add` gives one account a server for one folder only; `ccam mcp add` (or `claude mcp add -s user`, then `ccam mcp import`) shares it.
- A server that needs a sign-in (OAuth) is signed in once per account, with `/mcp` in a session of that account; ccam never copies sign-ins. claude.ai connectors belong to each claude.ai account.
- With the default `mcp_strategy` (flag), `claude mcp list` does not show the group's servers (sessions get them with `--mcp-config`): use `ccam ext` or `ccam mcp list` to check. Under `mcp_strategy = "sync"` ccam copies them into the account's `.claude.json`, so `claude mcp list` shows them.
- Servers left in plain claude's `~/.claude.json` load in no account until `ccam migrate --mcp` copies them; suggest it, do not run it.

## Do not

- Do not read or print `.credentials.json` or any token; ccam never needs them in a conversation.
- Do not run `ccam add`, `ccam login`, `ccam group` (other than `ccam group` and `ccam group show`, which only read), `ccam remove` or `ccam migrate` yourself: they log in through a browser or move files. Suggest the command to the user instead.
