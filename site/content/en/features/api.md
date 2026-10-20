---
title: API key accounts
description: ccam add --api adds an API key account for Anthropic's API or a gateway; the key is read at launch from a variable or a command, never stored.
updated: 2026-10-20
order: 13
translation: source
---

An API key account starts Claude Code with an API key instead of a Claude subscription login: Anthropic's API, or an endpoint that speaks the same Messages format, such as a LiteLLM or OpenRouter gateway. ccam does not store the key. It reads it each time the account starts, from an environment variable or from a command such as a password manager, and hands it only to that Claude Code process.

## Usage

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

An API key account shares settings, `CLAUDE.md`, extensions and MCP servers like any other account, and works with shortcuts, `ccam folder` and `ccam ext`.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--api` | Add an API key account instead of logging in | off |
| `--key-from env:NAME` | Read the key from the environment variable `NAME` | |
| `--key-from command:COMMAND` | Run `COMMAND` (with `sh -c`, or `cmd /C` on Windows) and use the first line it prints | |
| `--base-url <url>` | The endpoint (`ANTHROPIC_BASE_URL`) | Anthropic's API |
| `--auth api-key\|bearer` | Send the key as `ANTHROPIC_API_KEY` (`x-api-key`) or as `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, which gateways usually want) | `api-key` |
| `--model <model>` | The model (`ANTHROPIC_MODEL`) | Claude Code's default |

The settings are stored in `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## How it works

At each launch of the account, after the usual checks and before Claude Code starts:

1. ccam reads the key: `env:` from the environment, `command:` by running the command with a 30-second limit. A command can ask for a password on the terminal; its error output goes to the terminal too. An empty key stops the launch with a message saying where the key was expected.
2. It removes every inherited variable that chooses how Claude Code signs in or which provider it uses: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` and `CLAUDE_CODE_USE_FOUNDRY`. A key set in your shell for something else cannot reach this account, and a subscription login cannot win over its key.
3. It sets `ANTHROPIC_BASE_URL` and `ANTHROPIC_MODEL` when given, and the key as `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN`.

Accounts with a subscription login keep the environment they inherit, as before. In an interactive session Claude Code may ask once whether to use the key it finds in the environment.

## Safety and edge cases

- The key is never written to `config.toml`, a cache or a log, and `ccam --dry-run gw` does not read it: a `command:` source is not run by a dry run.
- `ccam add --api` checks an `env:` source right away and warns if the variable is empty; a `command:` source is first run at launch.
- `ccam login` and `ccam logout` refuse an API key account, since it has no login. `ccam doctor` shows its endpoint and key source, and warns when an `env:` variable is empty.
- `ccam usage` shows `API` for it. `ccam next --auto`, `--fallback` and `ccam pick` choose among subscription accounts only, because an API key has no 5-hour or weekly limit to compare.
- A shared `settings.json` with an `env` block or `apiKeyHelper` applies to every account in the group, API key accounts included. Keep provider settings out of shared settings.

## Related

- [Managing accounts](/features/accounts)
- [Launching Claude Code](/features/launch)
- [Command reference](/docs/commands)
