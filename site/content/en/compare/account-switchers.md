---
title: ccam and other Claude Code account switchers
description: How ccam compares with claude-swap, CCS, clauth, CAAM, cux and other tools for switching Claude Code accounts, and which approach fits which use.
updated: 2026-10-20
published: 2026-10-19
order: 2
translation: source
---

Tools for switching Claude Code accounts take one of two approaches. Most swap the login inside the one `~/.claude` folder, so one account is active at a time. Others, ccam among them, give each account its own folder with `CLAUDE_CONFIG_DIR`, so every account stays signed in and several can run at once. ccam adds shared settings and MCP servers, and continues a conversation in another account's folder.

## Swapping one login or separate folders?

**Swapping a login** keeps one set of settings, MCP servers and history, because there is only one folder. Switching is quick, and a running session can pick up the new login. But only one account is active at a time, for every terminal and editor at once.

To swap, the tool keeps a copy of every other account's login and writes it into place. Claude Code replaces the refresh token in a login as it renews it, so each saved copy has to be kept in step with the sessions that are running. Accounts that are logged out or need `/login` again after a switch are among the problems these projects most often fix.

**Separate folders** is the method [Claude Code's documentation describes](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Every account stays logged in, and you can run a work and a personal session side by side. Each login stays in its own folder, where Claude Code renews it as usual. The price is that each folder starts with its own settings, MCP servers and history, unless a tool shares them.

## How the tools compare

Taken from each project's own documentation and source code on 2026-10-20. "Not documented" means the project does not say, not that the feature is missing.

| Tool | Approach | Accounts at once | 5-hour and weekly usage | When a limit is reached | Platforms |
| --- | --- | --- | --- | --- | --- |
| ccam | Separate folders | Yes | Yes, `ccam usage` | You run `ccam next --auto --continue` to continue the conversation in the account with the most room; `ccam --fallback` reruns `-p` jobs there by itself | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Swaps the login; separate folders with `cswap run` | Yes, with `cswap run` | Yes | `cswap auto` switches the login at a threshold; the running session continues | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Separate folders | Yes | Partly | Account pools pass requests to the next account through a local proxy | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Swaps the login; separate folders with `clauth start` | Yes, with `clauth start` | Yes | Switches down a fallback chain | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Swaps the login, or isolated profiles | Yes, with isolated profiles | Yes | `caam run` retries non-interactive commands on another account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Swaps the login and wraps `claude` | No, one active account | Yes | Swaps the login and resumes the same conversation | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Swaps the login | No | Yes | Optional switch at a threshold | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Swaps the login | No | Yes | Optional switch at 100% | macOS |

ccswap is a fork of claude-swap that also manages Codex logins. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) and [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) also give each account its own folder. Several of these tools also handle other coding agents, such as Codex or Gemini.

## How they are built

| Tool | Written in | While Claude Code runs | Copies logins between places | What accounts can share |
| --- | --- | --- | --- | --- |
| ccam | Go, one static binary | Nothing: ccam replaces itself with `claude` (on Windows, which cannot, it waits for it) | No | Each item on its own: local or a named group, per account |
| claude-swap | Python 3.12+, installed with uv or pipx | Nothing with `cswap run`; `cswap auto` and the menu bar app keep running | Yes, when it swaps the login | With `cswap run`: settings, `CLAUDE.md` and extensions together, history optional; MCP servers copied at each launch |
| CCS | Node.js 18+ or Bun | CCS stays as the parent process; pools add a local proxy | No for its account folders; pools keep their own logins | Shared or local per account, history in context groups; MCP servers copied |
| clauth | Rust | clauth stays as the parent process | Yes, when it swaps the login | All of `~/.claude`, or nothing |
| CAAM | Go | Depends on the mode; an optional daemon | Yes, in its vault mode | Fixed by the mode; MCP servers copied at each launch |
| cux | Go | A wrapper stays around `claude` | Yes | Everything: one folder |
| CCSwitcher, Claude Switcher | macOS apps (Swift, Python) | The menu bar app | Yes | Everything: one folder |

## What ccam does well

### Logins stay where Claude Code put them

ccam never copies a Claude login from one place to another. Each account's login stays in its own folder, or on macOS in that folder's own Keychain item, and Claude Code is the program that renews it. There is no saved copy that can fall out of step. An existing `~/.claude` registered with `ccam import` keeps the Keychain item that plain `claude` uses, so it stays signed in either way.

ccam renews an expired token itself only when you run `ccam usage --refresh`. It takes a lock first and refuses while that account is running, so it never rotates a token under a live session.

### Nothing runs between you and Claude Code

ccam reads one file, checks its links and replaces itself with `claude` (`exec`). While you work there is no ccam process, wrapper, proxy or daemon: standard input, signals and the exit code belong to Claude Code, and its requests go straight to Anthropic. Windows cannot replace a process, so there ccam starts Claude Code with the same console and waits for its exit code, doing nothing else.

Starting an account this way adds about 4 ms, under 6 ms at the 95th percentile, with 8 MB of peak memory ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) with hyperfine, 300 runs). On the same laptop, starting the Python interpreter alone takes 13 ms and Node.js 31 ms, before a program does anything.

### Share item by item, account by account

MCP servers, `settings.json`, `CLAUDE.md`, extensions (agents, commands, skills, plugins) and conversation history are five separate items. Each account sets each item to `local` or to a group, so two work accounts can share history in a `work` group while a personal account keeps its own, and all three share the same MCP servers.

- MCP servers are passed with `--mcp-config` at launch, so by default ccam never edits `.claude.json`, the file Claude Code rewrites while it runs. `ccam mcp add` adds a server once for the whole group.
- The other items are links into the group folder, so a change made in one account shows up in the others at once. Nothing is copied at launch.
- Logins, `.claude.json` and running state are never shared.

### Every change is shown first and backed up

`ccam share` and `ccam migrate` print a plan and ask before they move anything. When a file differs between an account and its group, you choose: keep the group's, use yours, or merge the JSON. Originals go to `.ccam-backup/<time>/`, a removed account's folder goes to a trash folder, and files are replaced in one rename, so an interrupted write leaves the old file. `--dry-run` shows what launching, sharing, migrating, shortcuts, updates and repairs would do. `ccam doctor --fix` repairs broken links and permissions.

### The same conversation in the next account

With shared history, `ccam next --auto --continue` asks every account for its usage at once, opens the one with the most room and continues the conversation. You decide when to switch; ccam does not watch your usage in the background. The status line suggests the command once an account reaches 90%.

Scripts cannot wait for you, so `ccam --fallback work -p "…"` switches for them: when the account reaches its limit, the same command runs again as the account with the most room, and only the output of the run that finished is printed. A failure that is not a usage limit is passed on as it is.

### The right account in each folder

`ccam folder work ~/code/company` makes every launch without an account name inside that folder use `work`, while the rest of your projects keep the default account. The setting lives in the same `config.toml` every launch reads, so it costs no time.

### One binary, checked downloads

ccam is a static binary of about 10 MB with no runtime to install, for Linux, macOS and Windows. The install scripts (`install.sh`, `install.ps1`) and `ccam update` check every download against the release's SHA-256 checksums, and ccam never checks for updates on its own. It has no telemetry: the only requests it makes are usage lookups to Anthropic for your accounts, in parallel with a 3-second limit and a 60-second cache, and the downloads you ask for.

## Where other tools do more

- claude-swap, clauth, cux, CAAM and CCSwitcher can switch a running interactive session on their own when an account crosses a usage threshold. ccam switches `-p` runs with `--fallback` and leaves interactive sessions to you.
- On Windows without Developer Mode, ccam shares folders but cannot share single files such as `settings.json`; claude-swap and CCS copy them there instead.
- CCSwitcher, Claude Switcher and claude-swap show usage in the macOS menu bar.
- CCS, CAAM, clauth, ccswap and Claude Switcher also manage Codex and other coding agents.
- CAAM, CCSwitcher and Claude Switcher install with Homebrew, and clauth and CAAM sign their release checksums. ccam installs with its script, a `.deb` or `go install`, and checks checksums without a signature.

## Which one fits?

Choose a login-swapping tool if you use one account at a time or want running sessions to switch automatically. Choose ccam if you want accounts side by side, one shared setup, the same conversation in the next account and nothing running beside Claude Code, on Linux, macOS or Windows. If you need nothing shared, [plain shell aliases](/compare/aliases) are enough.

Whether you may use several accounts depends on each account's terms. Check them before you do.
