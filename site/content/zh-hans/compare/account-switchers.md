---
title: ccam 与其他 Claude Code 账号切换工具对比
description: ccam 与 claude-swap、CCS、clauth、CAAM、cux 等 Claude Code 账号切换工具的对比，以及哪种方式适合哪种用途。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

切换 Claude Code 账号的工具采用两种方式之一。大多数工具在同一个 `~/.claude` 文件夹中替换登录信息，所以同一时间只有一个账号处于活动状态。另一些工具（包括 ccam）用 `CLAUDE_CONFIG_DIR` 为每个账号提供单独的文件夹，所以每个账号都保持登录状态，并且可以同时运行多个账号。ccam 还提供共享的设置和 MCP 服务器，并能在另一个账号的文件夹中继续对话。

## 替换同一份登录信息，还是使用单独的文件夹？

**替换登录信息**只保留一套设置、MCP 服务器和历史记录，因为只有一个文件夹。切换很快，正在运行的会话也能用上新的登录信息。但同一时间只有一个账号处于活动状态，而且对所有终端和编辑器同时生效。

**单独的文件夹**是 [Claude Code 文档中介绍的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。每个账号都保持登录，你可以让工作会话和个人会话并行运行。代价是每个文件夹一开始都有各自的设置、MCP 服务器和历史记录，除非有工具来共享它们。

## 各工具对比

信息取自各项目自己的文档（2026-10-19）。"未说明"表示项目没有提及，而不是没有该功能。

| 工具 | 方式 | 同时使用多个账号 | 5 小时和每周用量 | 达到限制时 | 平台 |
| --- | --- | --- | --- | --- | --- |
| ccam | 单独的文件夹 | 是 | 是，`ccam usage` | 你运行 `ccam next --auto --continue`，在余量最多的账号中继续对话 | Linux、macOS、WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | 替换登录信息；用 `cswap run` 时使用单独的文件夹 | 是，使用 `cswap run` | 是 | `cswap auto` 在达到阈值时切换登录信息；正在运行的会话继续进行 | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 单独的文件夹 | 是 | 部分支持，在 macOS 菜单栏应用中 | 账号池通过本地代理把请求转给下一个账号 | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 替换登录信息；用 `clauth start` 时使用单独的文件夹 | 是，使用 `clauth start` | 是 | 沿备用链依次切换 | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 替换登录信息，或使用隔离的配置 | 是，使用隔离的配置 | 是 | `caam run` 在另一个账号上重试非交互式命令 | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | 替换登录信息并包装 `claude` | 未说明 | 是 | 替换登录信息并恢复同一个对话 | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 替换登录信息 | 否 | 是 | 未说明 | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 替换登录信息 | 否 | 是 | 可选在 100% 时切换 | macOS |

ccswap 是 claude-swap 的一个分支，还能管理 Codex 的登录。其中一些工具也支持其他编程代理，例如 Codex 或 Gemini。

## ccam 的侧重点

- 每个账号都保持登录。切换从不会让某个账号退出登录，`ccam work` 和 `ccam personal` 可以同时在两个终端中运行。
- MCP 服务器、`settings.json`、`CLAUDE.md`、agents、skills 和插件通过链接在文件夹之间共享，可以逐项、按账号选择。登录信息永远不会共享。
- 共享历史后，只需一条命令 `ccam next --auto --continue`，就能在余量最多的账号的文件夹中打开同一个对话。
- ccam 启动 Claude Code 后就用它替换自身，所以没有任何东西在后台运行。它不会监视你的用量，也不会自行切换；何时切换由你决定。

## 其他工具做得更多的地方

- claude-swap、CCS、clauth、CAAM 和 cux 无需 WSL 即可在 Windows 上运行。ccam 不行。
- claude-swap、clauth 和 cux 会在账号超过用量阈值时自行切换。ccam 把切换留给你。
- CCSwitcher、Claude Switcher 和 claude-swap 会在 macOS 菜单栏中显示用量。
- CCS、CAAM、ccswap 和 Claude Switcher 还能管理 Codex 和其他编程代理。

## 哪种适合你？

如果你一次只用一个账号、希望自动切换，或者在没有 WSL 的 Windows 上，选择替换登录信息的工具。如果你想让多个账号并行运行、共用一套配置，并在下一个账号中继续同一个对话，选择 ccam。如果什么都不需要共享，[普通的 shell 别名](/compare/aliases)就够了。

是否可以使用多个账号取决于每个账号的条款。使用前请先确认。
