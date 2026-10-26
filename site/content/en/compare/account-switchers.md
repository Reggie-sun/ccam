---
title: ccam and other Claude Code account switchers
description: How ccam compares with claude-swap, CCS, clausona, clauth, CAAM, cux and other tools for switching Claude Code accounts, and which approach fits which use.
updated: 2026-10-26
published: 2026-10-19
order: 2
translation: source
---

Tools for switching Claude Code accounts take one of two approaches. Most swap the login inside the one `~/.claude` folder, so one account is active at a time. Others, ccam among them, give each account its own folder with `CLAUDE_CONFIG_DIR`, so every account stays signed in and several can run at once. ccam adds shared settings and MCP servers, and continues a conversation in another account's folder.

## Swapping one login or separate folders?

**Swapping a login** keeps one set of settings, MCP servers and history, because there is only one folder. Switching is quick, and a running session can pick up the new login. But only one account is active at a time, for every terminal and editor at once.

To swap, the tool keeps a copy of every other account's login and writes it into place. Claude Code replaces the refresh token in a login as it renews it, so each saved copy has to be kept in sync with the sessions that are running. Accounts that are logged out or need `/login` again after a switch are among the problems these projects most often fix.

**Separate folders** is the method [Claude Code's documentation describes](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Every account stays logged in, and you can run a work and a personal session side by side. Each login stays in its own folder, where Claude Code renews it as usual. The price is that each folder starts with its own settings, MCP servers and history, unless a tool shares them.

## How the tools compare

Taken from each project's own documentation and source code on 2026-10-20 (clausona: v0.12.0-beta). "Not documented" means the project does not say, not that the feature is missing.

| Tool | Approach | Accounts at once | 5-hour and weekly usage | When a limit is reached | Platforms |
| --- | --- | --- | --- | --- | --- |
| ccam | Separate folders | Yes | Yes, `ccam usage` | A notification in the session; `ccam next --auto --continue` or `/ccam:switch` moves the conversation to the account with the most room, with or without shared history; optional `auto_switch` (Linux, macOS) moves it by itself; `ccam --fallback` reruns `-p` jobs; `ccam pick` spreads workers | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Swaps the login; separate folders with `cswap run` | Yes, with `cswap run` | Yes | `cswap auto` switches the login at a threshold; the running session continues | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Separate folders | Yes | Partly | Account pools pass requests to the next account through a local proxy | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Separate folders | Yes | Yes | Routes choose an account with quota before each run, can wait for a reset, and hand out distinct accounts to workers | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Swaps the login; separate folders with `clauth start` | Yes, with `clauth start` | Yes | Switches down a fallback chain | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Swaps the login, or isolated profiles | Yes, with isolated profiles | Yes | `caam run` retries non-interactive commands on another account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Swaps the login and wraps `claude` | No, one active account | Yes | Swaps the login and resumes the same conversation | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Swaps the login | No | Yes | Optional switch at a threshold | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Swaps the login | No | Yes | Optional switch at 100% | macOS |

ccswap is a fork of claude-swap that also manages Codex logins. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) and [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) also give each account its own folder. Several of these tools also handle other coding agents, such as Codex or Gemini.

## How they are built

| Tool | Written in | While Claude Code runs | Copies logins between places | What accounts can share |
| --- | --- | --- | --- | --- |
| ccam | Go, one static binary | Nothing by default: ccam replaces itself with `claude` (on Windows, where that is not possible, it waits for Claude Code to exit); with `relay` or `auto_switch` on, ccam keeps running alongside interactive sessions | No | Groups: each names the items it shares and its accounts; an account can be in several |
| claude-swap | Python 3.12+, installed with uv or pipx | Nothing with `cswap run`; `cswap auto` and the menu bar app keep running | Yes, when it swaps the login | With `cswap run`: settings, `CLAUDE.md` and extensions together, history optional; MCP servers copied at each launch |
| CCS | Node.js 18+ or Bun | CCS stays as the parent process; pools add a local proxy | No for its account folders; pools keep their own logins | Shared or local per account, history in context groups; MCP servers copied |
| clausona | TypeScript, Node.js 20+ | Plain `claude` runs through a shell function; `clausona run` keeps Node as the parent; a background Node process after each exit | No, but it renews idle accounts' tokens itself to read their usage | Everything except per-account files, from one primary; history optional; servers added with `claude mcp add` are not shared |
| clauth | Rust | clauth stays as the parent process | Yes, when it swaps the login | All of `~/.claude`, or nothing |
| CAAM | Go | Depends on the mode; an optional daemon | Yes, in its vault mode | Fixed by the mode; MCP servers copied at each launch |
| cux | Go | A wrapper stays around `claude` | Yes | Everything: one folder |
| CCSwitcher, Claude Switcher | macOS apps (Swift, Python) | The menu bar app | Yes | Everything: one folder |

## What ccam does well

### Logins stay where Claude Code put them

ccam never copies a Claude login from one place to another. Each account's login stays in its own folder, or on macOS in that folder's own Keychain item, and Claude Code is the program that renews it. There is no saved copy that can fall out of sync. An existing `~/.claude` registered with `ccam import` keeps the Keychain item that plain `claude` uses, so it stays signed in either way.

ccam renews an expired token itself only when you run `ccam usage --refresh`. It takes a lock first and refuses while that account is running, however Claude Code was started (including background sessions and their daemon), and while Claude Code itself is renewing that login, so it never rotates a token under a live session.

### By default, nothing runs between you and Claude Code

ccam reads one file, checks its links and replaces itself with `claude` (`exec`). By default, while you work there is no ccam process, wrapper, proxy or daemon: standard input, signals and the exit code belong to Claude Code, and its requests go straight to Anthropic. Windows cannot replace a process, so there ccam starts Claude Code with the same console and waits for its exit code, doing nothing else. With `relay = true` or `auto_switch = true`, ccam keeps running the same way alongside interactive sessions, so it can continue a conversation on another account; Claude Code's requests still go straight to Anthropic.

Starting an account this way adds about 4 ms, under 6 ms at the 95th percentile, with 8 MB of peak memory ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) with hyperfine, 300 runs). On the same laptop, starting the Python interpreter alone takes 13 ms and Node.js 31 ms, before a program does anything.

### Groups decide what is shared

MCP servers, `settings.json`, `CLAUDE.md`, extensions (agents, commands, skills, plugins) and conversation history are five separate items. A group names the items it shares and the accounts in it, so two work accounts can share history in a `work` group while a personal account keeps its own, and all three share the same MCP servers through `default`. `ccam group set work --add memory` updates every account in the group in one plan; the [dashboard](/features/dashboard) does the same with checkboxes and grays out an item an account already gets from another group.

- MCP servers are passed with `--mcp-config` at launch, so by default a launch never touches `.claude.json`, the file Claude Code rewrites while it runs. ccam writes there only in steps that say so, under Claude Code's lock and with a backup: for Remote Control after asking, when an account leaves or rejoins MCP sharing, and for `ccam mcp import`. `ccam mcp add` adds a server once for the whole group.
- The other items are links into the group folder, so a change made in one account shows up in the others at once. Nothing is copied at launch.
- Logins, `.claude.json` and running state are never shared.

### Every change is shown first and backed up

`ccam group` and `ccam migrate` print a plan and ask before they move anything. When a file differs between an account and its group, you choose: keep the group's, use yours, or merge the JSON; inside folders the group's copy stays and the plan names the files that differ. Originals go to `.ccam-backup/<time>/`, a removed account is logged out and its folder goes to a trash folder, and files are replaced in one rename, so an interrupted write leaves the old file. `--dry-run` shows what launching, sharing, migrating, shortcuts, updates and repairs would do. `ccam doctor --fix` repairs broken links and permissions.

### The same conversation in the next account

`ccam next --auto --continue` asks every account for its usage at once, opens the one with the most room and continues the conversation. With shared history it is the same conversation; without, Claude Code copies it into the new account (`--resume <file> --fork-session`), and ccam says how many tokens the first reply reads again.

When a limit stops a session, the ccam plugin's hook shows a notification, and `/ccam:switch` hands the conversation on. By default you decide when to switch. With `auto_switch` on (Linux, macOS), ccam moves the session itself, but only after Claude Code has stopped it at a limit, and only to an account of the same kind as the one it leaves: another seat of the same organization, or another personal plan ([details](/features/next#auto-switch)).

Nothing watches your usage in the background, and by default automatic choices never move work out of a Team or Enterprise organization or into one, so a company conversation never lands on a private account by accident. The status line suggests the command once an account reaches 90%.

Scripts cannot wait for you, so `ccam --fallback work -p "…"` switches for them: when the account reaches its limit, the same command runs again as the account with the most room, and only the output of the run that finished is printed. A failure that is not a usage limit is passed on as it is. For parallel workers, `ccam pick` hands each one a different account with room, and both exit with code 75 when no account has any left.

### The right account in each folder

`ccam folder work ~/code/company` makes every launch without an account name inside that folder use `work`, while the rest of your projects keep the default account. The setting lives in the same `config.toml` every launch reads, so it costs no time.

### One binary, checked downloads

ccam is a static binary of about 10 MB with no runtime to install, for Linux, macOS and Windows. The install scripts (`install.sh`, `install.ps1`) and `ccam update` check every download against the release's SHA-256 checksums, and ccam never updates without asking. It has no telemetry: the only requests it makes are usage lookups to Anthropic for your accounts, in parallel with a 3-second timeout and a 60-second cache, the dashboard's check on GitHub for a newer release when it opens, and the downloads you ask for. With `usage_source = "claude"`, even those go through Claude Code, and ccam sends no token.

## Where other tools do more

- claude-swap, clauth, cux, CAAM and CCSwitcher can switch a running interactive session on their own when an account crosses a usage threshold. ccam switches `-p` runs with `--fallback`, and moves an interactive session by itself only with `auto_switch` on (Linux, macOS) and only after Claude Code has stopped it at a limit; otherwise you move it, with a notification at the limit and `/ccam:switch`.
- On Windows without Developer Mode, ccam shares folders but cannot share single files such as `settings.json` (the first `ccam add` then shares only MCP servers and extensions); claude-swap and CCS copy them there instead.
- CCSwitcher, Claude Switcher and claude-swap show usage in the macOS menu bar.
- CCS, clausona, CAAM, clauth, ccswap and Claude Switcher also manage Codex and other coding agents.
- clausona also turns extensions on and off per account (ccam's `ccam ext` only shows them), tracks local token cost itself (ccam hands that to [ccusage](https://ccusage.com): `ccam ccusage`), and ships Claude Code plugins that run fleets of workers in tmux and other terminals.
- CAAM, CCSwitcher and Claude Switcher install with Homebrew, and clauth and CAAM sign their release checksums. ccam installs with its script, a `.deb` or `go install`, and checks checksums without a signature.

## Which one fits?

Choose a login-swapping tool if you use one account at a time, or if you want a running session to switch accounts without restarting, even before it hits the limit. Choose ccam if you want accounts side by side, one shared setup, the same conversation in the next account and, by default, nothing running alongside Claude Code, on Linux, macOS or Windows. If you need nothing shared, [plain shell aliases](/compare/aliases) are enough.

Whether you may use more than one account depends on the terms that apply to each one. Read those terms before you start.

**Team and Enterprise seats.** Automatic account selection (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n` and, once you turn it on, [`auto_switch`](/features/next#auto-switch)) also switches between seats of the same Claude Team or Enterprise organization, which may violate Anthropic's terms. Read the [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) and check with your organization and Anthropic before you rely on it.
