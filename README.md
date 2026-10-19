# ccam — Claude Code multi-account manager

**English** · [한국어](README.ko.md) · Website: https://ccam.dev

- ccam is a free, open-source command-line tool that lets you use several Claude Code accounts on one computer and switch between them with one command.
- Each account keeps its own login in a separate CLAUDE_CONFIG_DIR, while MCP servers, settings, CLAUDE.md, extensions and conversation history can be shared between accounts, so you can continue a conversation on another account when one reaches its usage limit.
- Install it on Linux or macOS with one command, then run ccam add for each account.
- ccam is an unofficial project and is not affiliated with or endorsed by Anthropic.

`ccam` is a single Go binary for adding and switching between several Claude Code accounts in an Ubuntu terminal. Each account gets its own `CLAUDE_CONFIG_DIR`, so logins stay separate, and you choose per item whether MCP servers, settings, the global CLAUDE.md, extensions and conversation history are shared.

> **Check the terms first.** Whether using several accounts is allowed depends on each account's Anthropic terms and your organization's policy. ccam makes switching convenient; it does not work around usage limits.

All program output (messages, prompts, help) is in English.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux (Ubuntu, Debian, Fedora, Arch, Alpine, …) and macOS are supported. This installs `~/.local/bin/ccam` after verifying the release checksum. Alternatives: the `.deb` from [Releases](https://github.com/jhyunwoo/ccam/releases), or `go install github.com/jhyunwoo/ccam/cmd/ccam@latest` (Go 1.26+). Claude Code (`claude`) must be installed.

## Quick start

```bash
ccam add             # name → shared items → browser login → shortcut claude1
ccam add work        # second account → claude2

claude2              # Claude Code as account 2
ccam work --resume   # by name; everything after the name goes to claude unchanged
ccam                 # picker (Enter launch, 1-9 jump, s sharing, a add)
ccam usage           # 5-hour and weekly usage of every account
ccam next --continue # out of quota? continue the same conversation on the next account
```

An existing `~/.claude` can be registered in place with `ccam import` (nothing is moved or copied).

## Commands

| Command | What it does |
| --- | --- |
| `ccam` | Picker (prints the list when not on a terminal) |
| `ccam <name\|number> [claude args]` | Run `claude` as that account; args are passed through verbatim |
| `claude1` … `claudeN` | Launch by number (`ccam shims install --names` adds `claude-<name>`) |
| `ccam add [name]` | Add an account (wizard on a terminal; `--no-login`, `--share item=group`) |
| `ccam list [--json]` | List accounts |
| `ccam use <name>` | Default account for name-less runs such as `ccam -p "…"` |
| `ccam next [--auto] [--continue]` | Next account; `--auto` picks the lowest 5-hour usage |
| `ccam usage [--watch] [--json] [-v]` | Usage (`ccam u`) |
| `ccam share <name> --history on …` | Change sharing; prints a migration plan and asks first |
| `ccam mcp add\|add-json\|list\|remove\|import` | Shared MCP servers |
| `ccam login\|logout <name>` | Log in again / remove the token |
| `ccam rename` / `ccam remove` | Rename / remove (folder goes to a trash folder) |
| `ccam import [dir]` | Register an existing config folder |
| `ccam doctor [--fix]` | Check links, permissions, login, version, usage format; repair |
| `ccam statusline install` | Show `[account] · model · 5h N%` in Claude Code's status line |
| `ccam completion bash\|zsh` | Shell completion |

Put `--dry-run` first to print what would happen (`ccam --dry-run work`).

## Sharing

Every item is `local` or a group name (default `default`); only accounts in the same group share it. MCP servers are passed with `--mcp-config <group>/mcp.json` at launch (`.claude.json` untouched; `mcp_strategy = "sync"` merges into it instead). Settings, CLAUDE.md, extension folders and `projects/` + `history.jsonl` are symlinked to the group folder. `.credentials.json` and `.claude.json` are never shared. Turning sharing on or off backs originals up to `.ccam-backup/<time>/` and lets you resolve differing files (keep group / use mine / merge JSON).

## Usage lookup

Queries every account's usage in parallel with its own token (3 s per account, 60 s cache). The endpoint is undocumented and may change; on failure the last known or status-line values are shown and `ccam doctor` reports format changes. Expired tokens are not refreshed unless you pass `--refresh`, and never while that account is running.

## Performance

`ccam <name>` reads one file, checks links with `lstat`, and `exec`s into `claude`: about 4 ms p95 added on a typical laptop, 8 MB peak memory. `CCAM_TRACE=1` prints per-step timings; `CCAM_PPROF=<file>` writes a CPU profile.

## Development

```bash
go test ./...
golangci-lint run
scripts/perf-check.sh           # launch-path targets with hyperfine
scripts/bench-compare.sh main   # fail on >20% benchmark regressions (benchstat)
```

Design and decision log: [docs/PLAN.md](docs/PLAN.md). Claude Code folder findings: [docs/claude-layout.md](docs/claude-layout.md).

## License

[MIT](LICENSE). ccam is an unofficial project and is not affiliated with or endorsed by Anthropic.
