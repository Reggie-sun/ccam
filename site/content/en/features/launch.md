---
title: Launching Claude Code
description: Start Claude Code as any account with ccam <name|number>, ccam run or the ccam picker; arguments go to claude unchanged, and --dry-run shows the command.
updated: 2026-10-19
order: 1
translation: source
---

Launching is what you do with ccam most of the time. `ccam work` starts Claude Code as the account `work`, and everything you type after the name goes to `claude` as it is. Use it wherever you would type `claude`: in a terminal, in a script, or with input piped in.

## Usage

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

You can name an account or give its number from `ccam list`. ccam does not read the arguments after the name, so `ccam work --help` shows Claude Code's help, not ccam's. A `--` after the name is optional: `ccam work -- --version` and `ccam work --version` both run `claude --version`. Standard input, the exit code and signals behave as with `claude`, because ccam is replaced by Claude Code rather than running it as a child.

When the first argument starts with `-` and is not one of ccam's own flags, ccam starts the default account with all arguments. The default is the account set with `ccam use`, or the account with the lowest number if you have not set one.

### The picker

`ccam` without arguments opens the picker when both standard input and standard output are a terminal. It lists every account with its number, email, 5-hour usage, last use and shared items. The usage column fills in as each account answers. Without a terminal, `ccam` prints the same table as `ccam list`.

| Key | What it does |
| --- | --- |
| `Enter` | Start Claude Code as the selected account |
| `1`–`9` | Start the account with that number |
| `Up`/`Down`, `k`/`j` | Move the selection |
| `a` | Close the picker and run the `ccam add` wizard |
| `s` | Edit the shared items of the selected account (`Space` or `x` toggles, `Enter` applies, `Esc` goes back) |
| `d` | Make the selected account the default |
| `r` | Fetch usage again, without the short cache |
| `q`, `Esc`, `Ctrl+C` | Quit |

The picker closes before it acts, so the terminal is restored before Claude Code starts. After `d` or a sharing change it opens again. A sharing change shows the plan and asks before it moves anything, as `ccam share` does.

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--dry-run` | Print the environment and command ccam would run, then exit | off |
| `-v`, `--version` | Print `ccam <version>` (only as the first argument) | |
| `-h`, `--help` | ccam's help, as `ccam --help` or `ccam run --help` | |

Argument rules:

- ccam reads its own flags only before the account name. After the name, every argument belongs to Claude Code, including `--dry-run` and `--help`.
- `--dry-run` works as `ccam --dry-run work`, `ccam run --dry-run work` and `ccam --dry-run` (the default account). It is honoured, not refused.
- `ccam run <name|number> [claude args]` is the long form of `ccam <name>`. `ccam run` without a name uses the default account.
- Account names cannot be ccam command names such as `add`, `list`, `run` or `usage`, so `ccam <word>` is never ambiguous.

The dry run prints one line you could paste into a shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## How it works

Each launch takes the same steps, from the command line to `exec`:

<!-- motion:launch -->

Each box lights up as its step runs, and the lines underneath are the environment and arguments Claude Code starts with.

1. `cmd/ccam/main.go` calls `app.Main`. If ccam was started under another name, it is a [shortcut](/features/shortcuts) and goes to the launcher directly.
2. `app.parse` looks at `os.Args` before cobra is involved. Only a ccam command name (`config.Reserved`), `-h`/`--help` or `-v`/`--version` goes to `cli.Execute`. Everything else is a launch.
3. `launcher.Launch` reads `~/.config/ccam/config.toml`, the only file it parses.
4. `launcher.Prepare` finds the account by name or number, drops one leading `--`, and resolves the `claude` executable with `ResolveClaude`. The PATH lookup is cached in `~/.local/state/ccam/claude-path` together with the file's modification time, so ccam does not search `PATH` on every launch. A changed or missing file triggers a new lookup. ccam refuses a `claude` that is ccam itself, compared by inode, because that would loop.
5. `checkLinks` compares the account folder `~/.local/share/ccam/profiles/<name>/` with its sharing settings using `share.Check`, which only calls `lstat` and `readlink`, one system call or two per item. A missing link is created on the spot. Anything that could lose data stops the launch.
6. `applyMCP` adds the shared MCP servers (see below).
7. `historyCheck` and `session.Register` record the session. ccam writes a marker named after its PID, holding the process start time, to `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`. `exec` keeps the PID, so the marker then names the Claude Code process itself. A marker, not a held lock, because processes that Claude Code starts would inherit a lock and keep it after the session ends.
8. `Plan.Exec` sets the terminal title to `claude · <name>` (when standard error is a terminal), writes the launch time to `~/.local/state/ccam/last-used/<name>` without `fsync`, and calls `syscall.Exec`.

Claude Code runs with `CLAUDE_CONFIG_DIR` set to the account folder and `CCAM_PROFILE` set to the account name. ccam removes any values of these that it inherited. For a `~/.claude` folder registered with `ccam import` on macOS, it also sets `CLAUDE_SECURESTORAGE_CONFIG_DIR` to an empty value, so that the account and plain `claude` use the same Keychain login.

### Why exec

ccam does not stay in memory while you work. `syscall.Exec` replaces the ccam process with `claude`, so keystrokes, signals and the exit code are Claude Code's own, and arguments are passed as an array, never joined into a shell string. Anything ccam wants to do after a session, such as noticing MCP servers you added with `claude mcp add`, happens at the next launch: ccam compares the modification time of the account's `.claude.json` with a note in `~/.local/state/ccam/mcp/<name>/seen` and parses the file only when it changed.

### How a command line is routed

Three command lines, three routes through `internal/app`:

<!-- motion:launch-route -->

A shortcut is matched on `argv[0]` against `config.toml`. `ccam <name>` is matched on the first argument. Both end in the same `syscall.Exec`. Only a ccam command such as `ccam usage` reaches cobra and runs inside ccam.

### Shared MCP servers

With the default `flag` strategy, ccam passes the group's MCP servers to Claude Code as an argument, placed before yours:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` returns an argument only if the group's `mcp.json` exists, and `Plan.InsertArgs` puts it right after `claude`. Your own `--mcp-config` or `--strict-mcp-config` then still applies on top. It is one argument, `--mcp-config=<path>`, because Claude Code's `--mcp-config` takes several values: written as two arguments, it would also read a following prompt as a config file.

With `mcp_strategy = "sync"` in `config.toml`, no argument is added. ccam writes the group's servers into the account's `.claude.json` before the launch instead, and skips that step if the account is already running.

### The 10 ms budget

ccam adds at most 10 ms (p95) before Claude Code starts. On the launch path ccam never builds the cobra command tree, starts Bubble Tea or opens a network connection, and package `init()` functions do no file, terminal or network work (`TestPackageInitIsCheap` fails if they take more than 5 ms in total). `scripts/perf-check.sh` measures `ccam --dry-run work`, `ccam work` and a shortcut with `hyperfine` and fails above 10 ms, and `BenchmarkLaunchPath` in `internal/launcher` runs the steps in a loop. To see the time per step, set `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Safety and edge cases

- An unknown name or number stops with ``ccam: no account "<name>" (see `ccam list`)``. With no accounts at all, ccam asks you to add one with `ccam add <name>`.
- If a shared path became an ordinary file or points to the wrong place, ccam stops with `stopped because of share link problems:`, lists them, and asks you to run `ccam doctor --fix`. It does not move data on the launch path.
- If the account folder is missing, ccam stops and points to `ccam doctor`.
- If another Claude Code that shares conversation history with this account (same group) runs in the same folder, ccam warns that the history can mix and asks `Run anyway? [y/N]`. Without a terminal it shows the warning and continues.
- `--dry-run` skips the MCP sync, the new-server notice, the history check, the session marker, the title and the last-use time, and does not exec. It still checks share links, and creates a missing link as a normal launch would.
- With `mcp_strategy = "sync"`, a failed sync prints `ccam: MCP sync failed (launching anyway): …` and Claude Code starts with the servers it already had.

## Related

- [Shortcut commands](/features/shortcuts)
- [Managing accounts](/features/accounts)
- [Command reference](/docs/commands)
- [Share MCP servers, settings and history](/docs/sharing)
- [How CLAUDE_CONFIG_DIR works](/guides/claude-config-dir)
