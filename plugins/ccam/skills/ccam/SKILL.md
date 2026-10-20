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

This session cannot move itself to another account. Tell the user to quit and run:

```sh
ccam next --auto --continue
```

It opens the account with the most room and continues this conversation there (both accounts must share conversation history; `ccam share <name> --history on`).

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

`ccam ext` shows each skill, command, agent, plugin, MCP server and hook per account, with its source (a share group name, `account`, or `project`), and lists problems such as broken links, skill folders without `SKILL.md`, and MCP servers defined differently in two places. Suggest `ccam doctor --fix` for broken share links and `ccam mcp import <account> <server>` to share a server the account added on its own.

## Do not

- Do not read or print `.credentials.json` or any token; ccam never needs them in a conversation.
- Do not run `ccam add`, `ccam login`, `ccam share`, `ccam remove` or `ccam migrate` yourself: they log in through a browser or move files. Suggest the command to the user instead.
