---
title: 在一台电脑上使用多个 Claude Code 账号
description: 借助 CLAUDE_CONFIG_DIR 让多个 Claude Code 账号并行运行：可以用 shell 别名手动设置，也可以用 ccam 自动完成，并在账号之间切换。
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "381de7a473b2"
---

为每个账号指定单独的配置文件夹。Claude Code 从 `CLAUDE_CONFIG_DIR` 环境变量中读取文件夹位置，每个文件夹都保存各自的登录信息、设置和历史记录，所以每个账号都保持登录状态，你也可以同时运行它们。你可以用 shell 别名手动设置，也可以交给 ccam 完成。

## 如何手动设置？

这是 [Claude Code 文档中介绍的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。在 `~/.bashrc` 或 `~/.zshrc` 中为每个额外的账号添加一行：

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

打开一个新终端并运行 `claude-work`。因为文件夹是新的，Claude Code 会要求你登录。从此以后，`claude` 是你的第一个账号，`claude-work` 是第二个。

在没有 WSL 的 Windows 上，可以在 PowerShell 配置文件（`notepad $PROFILE`）中添加一个函数，效果相同：

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## 别名方式有哪些不足？

它适合两个账号。需要注意：

- 每个文件夹一开始都是空的：没有设置、没有 `CLAUDE.md`、没有 MCP 服务器，也没有 agents 或 skills。你需要复制或链接所需的内容，并自己保持这些副本同步。
- 别名只存在于交互式 shell 中。脚本、Makefile 和其他程序看不到它们。
- MCP 服务器保存在 `.claude.json` 中，与登录信息在一起，所以不能简单地在文件夹之间链接。
- 对话只留在开始时所用的账号中，所以你无法在另一个账号中继续对话。
- 没有办法一次查看所有账号的用量。

## 能否改用 /login 切换账号？

可以，前提是你只是偶尔切换。在 Claude Code 中输入 `/login` 会[切换当前文件夹的账号](https://code.claude.com/docs/en/quickstart)。同一时间只有一个账号处于登录状态，所以每次切换都要重新登录，也无法让两个账号并行运行。使用单独的文件夹可以避免这两个问题。

## 如何用 ccam 设置？

ccam 把同样的思路自动化，并补上了缺少的部分：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

现在，`ccam personal` 和 `ccam work` 可以启动这两个账号，你为 work 账号设置的快捷命令 `ccwork` 也可以。因为它们是真正的程序而不是别名，所以在脚本中也能使用。单独运行 `ccam` 可以从列表中选择账号，运行 `ccam use work` 可以把某个账号设为默认账号。

设置、`CLAUDE.md`、MCP 服务器和扩展从一开始就共享。`ccam next --continue` 会在另一个账号上继续对话。如果你开启对话历史共享，两个账号处理的就是同一个对话，而不是一份副本。`ccam usage` 会一次显示所有账号的 5 小时用量和每周用量。

如果你已经在使用 `~/.claude`，请运行 `ccam import`，而不是再次添加该账号；或者运行 `ccam migrate`，把其中的设置、MCP 服务器和扩展复制到共享的 `default` 组中。

## 如何知道当前运行的是哪个账号？

在 Claude Code 中，`/status` 会显示你当前登录的账号。使用 ccam 时，运行一次 `ccam statusline install`，状态栏就会一直显示账号名称及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。

## 在 Windows 和 WSL 上可以用吗？

Claude Code [可以在 Windows 上原生运行，也可以在 WSL 中运行](https://code.claude.com/docs/en/setup)。`CLAUDE_CONFIG_DIR` 在两者中都有效，所以在原生 Windows 上，用上面的别名或 PowerShell 函数就够了。ccam 也可以在 Windows 上原生运行：在 PowerShell 中用 `irm https://ccam.dev/install.ps1 | iex` 安装，或者在 WSL 中用 Linux 的命令安装。快捷命令会变成 `ccwork.exe`，共享 `settings.json` 这样的文件需要开启开发人员模式（[安装 ccam](/docs/install)）。

## 故障排除

**Claude Code 每次都要求我登录。** 在新文件夹中首次启动时总会要求登录。如果一直要求登录，说明登录已过期或被删除：运行 `ccam login work`（或者通过别名启动一次并登录）。

**某个 MCP 服务器只出现在一个账号中。** 在会话中用 `claude mcp add` 添加的服务器属于该账号，不加 `-s user` 时还只属于当前文件夹。用 `ccam mcp import work <server>` 把它移到所有账号。

**各账号的设置不一致。** `ccam share work` 会列出该账号共享的内容以及通过哪个组共享；`ccam group join default work` 会先显示将要发生的变更，再让它共享 `default` 所共享的内容。`ccam doctor --fix` 会修复损坏的链接，并为它修改的每个文件保留备份。

**两个账号打开了同一个对话。** 避免在两个账号中同时处理同一个对话。你这样做时，ccam 会发出警告。

## 该选哪一种？

如果只有两个账号且不需要共享任何内容，用别名就够了。如果你想要共享设置和 MCP 服务器、使用更多账号、一眼看清用量或继续对话，ccam 可以帮你省去这些繁琐的管理工作。[与别名的对比](/compare/aliases)并排列出了两者的差异，[ccam 与其他账号切换工具](/compare/account-switchers)则介绍了那些切换同一份登录信息的工具。

你能否使用多个账号，取决于适用于每个账号的条款。开始之前请先阅读这些条款。

**Team 和 Enterprise 席位。** 自动选择账号（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换，这可能违反 Anthropic 的条款。在依赖这项功能之前，请先阅读[常见问题](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，并向你所在的组织和 Anthropic 确认。
