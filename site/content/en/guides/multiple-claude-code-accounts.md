---
title: Use multiple Claude Code accounts on one computer
description: Run several Claude Code accounts side by side with CLAUDE_CONFIG_DIR, by hand with shell aliases or automatically with ccam, and switch between them.
updated: 2026-10-25
published: 2026-10-19
order: 1
translation: source
---

Give each account its own configuration folder. Claude Code reads the folder from the `CLAUDE_CONFIG_DIR` environment variable, and each folder holds its own login, settings and history, so every account stays signed in and you can run them at the same time. You can set this up by hand with a shell alias, or let ccam do it.

## How do I set it up by hand?

This is the approach [Claude Code's documentation describes](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Add a line per extra account to `~/.bashrc` or `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Open a new terminal and run `claude-work`. Claude Code asks you to log in, because the folder is new. From then on, `claude` is your first account and `claude-work` the second. Add one alias per account, each with its own folder.

On Windows without WSL, a PowerShell function in your profile (`notepad $PROFILE`) does the same:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## What does the alias approach leave out?

It works well for two accounts. Things to know:

- Each folder starts empty: no settings, no `CLAUDE.md`, no MCP servers, no agents or skills. You copy or link what you need, and keep the copies in step yourself.
- Aliases exist only in interactive shells. Scripts, Makefiles and other programs do not see them.
- MCP servers live in `.claude.json`, next to the login, so they cannot simply be linked between folders.
- Conversations stay in the account where they started, so you cannot continue one in another account.
- There is nothing that shows the usage of all accounts at once.

## Can I switch accounts with /login instead?

Yes, if you only switch now and then. Typing `/login` inside Claude Code [switches the account](https://code.claude.com/docs/en/quickstart) of the folder you are in. Only one account is signed in at a time, so you log in again at every switch, and you cannot run two accounts side by side. Separate folders avoid both.

## How do I set it up with ccam?

ccam automates the same idea and adds the missing parts:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Now `ccam personal` and `ccam work` start the two accounts, and so does `ccwork`, the shortcut you gave the work account. Because they are real programs rather than aliases, they also work in scripts. Run `ccam` on its own to pick an account from a list, and `ccam use work` to make one the default.

Settings, `CLAUDE.md`, MCP servers and extensions are shared from the start. `ccam next --continue` carries a conversation over to the other account, and shared conversation history, which you can turn on, keeps it one conversation. `ccam usage` shows the 5-hour and weekly usage of every account at once.

If you already use `~/.claude`, run `ccam import` instead of adding that account again, or `ccam migrate` to copy its settings, MCP servers and extensions into the shared settings.

## How do I know which account is running?

Inside Claude Code, `/status` shows the account you are signed in to. With ccam, run `ccam statusline install` once and the status line shows the account name and its usage all the time, for example `[work] · Opus · 5h 24% · 7d 41%`.

## Does this work on Windows and WSL?

Claude Code runs [natively on Windows and inside WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` works in both, so the alias or PowerShell function above is enough on native Windows. ccam runs natively on Windows too: install it in PowerShell with `irm https://ccam.dev/install.ps1 | iex`, or inside WSL with the Linux command. Shortcuts become `ccwork.exe`, and sharing files such as `settings.json` needs Developer Mode ([Install ccam](/docs/install)).

## Troubleshooting

**Claude Code asks me to log in every time.** The first start in a new folder always asks. If it keeps asking, the login has expired or was removed: run `ccam login work` (or start the alias once and log in).

**An MCP server shows up in only one account.** Servers added inside a session with `claude mcp add` belong to that account, and without `-s user` to that folder only. Move one to all accounts with `ccam mcp import work <server>`.

**Settings differ between accounts.** `ccam share work` prints what the account shares and through which group; `ccam group join default work` makes it share what `default` shares, after showing what will change. `ccam doctor --fix` repairs broken links and keeps a backup of every file it changes.

**Two accounts open the same conversation.** Avoid working on one conversation in two accounts at the same time. ccam warns you when you try.

## Which one should I choose?

For two accounts and nothing shared, an alias is enough. When you want shared settings and MCP servers, more accounts, usage at a glance or continuing conversations, ccam saves the bookkeeping. The [comparison with aliases](/compare/aliases) lists the differences side by side, and [ccam and other account switchers](/compare/account-switchers) covers tools that swap one login instead.

Whether you may use several accounts depends on each account's terms. Check them before you do.

**Team and Enterprise seats.** Automatic choice (`ccam next`, `--fallback`, `ccam pick`, the dashboard's `n`, and [`auto_switch`](/features/next#auto-switch) when you turn it on) also switches between seats of the same Claude Team or Enterprise organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it" and "Users can't route their work through another user's account to use that user's remaining usage" ([Team plans](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); likewise for [Enterprise plans](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Using automatic choice between seats may therefore violate Anthropic's terms. Whether you use it there is your call: check with your organization and Anthropic before you rely on it. Personal plans (Pro, Max) are covered by the note above.
