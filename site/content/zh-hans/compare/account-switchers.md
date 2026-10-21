---
title: ccam 与其他 Claude Code 账号切换工具对比
description: ccam 与 claude-swap、CCS、clausona、clauth、CAAM、cux 等 Claude Code 账号切换工具的对比，以及哪种方式适合哪种用途。
updated: 2026-10-21
order: 2
translation: machine
sourceHash: "c74be64d2978"
---

切换 Claude Code 账号的工具采用两种方式之一。大多数工具在同一个 `~/.claude` 文件夹中替换登录信息，所以同一时间只有一个账号处于活动状态。另一些工具（包括 ccam）用 `CLAUDE_CONFIG_DIR` 为每个账号提供单独的文件夹，所以每个账号都保持登录状态，并且可以同时运行多个账号。ccam 还提供共享的设置和 MCP 服务器，并能在另一个账号的文件夹中继续对话。

## 替换同一份登录信息，还是使用单独的文件夹？

**替换登录信息**只保留一套设置、MCP 服务器和历史记录，因为只有一个文件夹。切换很快，正在运行的会话也能用上新的登录信息。但同一时间只有一个账号处于活动状态，而且对所有终端和编辑器同时生效。

为了替换，工具会为其他每个账号保存一份登录信息的副本，并在切换时把它写回原位。Claude Code 每次续期登录时都会换掉其中的刷新令牌，所以每份保存的副本都必须和正在运行的会话保持同步。切换后账号被登出，或者需要重新 `/login`，是这些项目最常修复的问题之一。

**单独的文件夹**是 [Claude Code 文档中介绍的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。每个账号都保持登录，你可以让工作会话和个人会话并行运行。每份登录信息都留在各自的文件夹中，由 Claude Code 照常在那里续期。代价是每个文件夹一开始都有各自的设置、MCP 服务器和历史记录，除非有工具来共享它们。

## 各工具对比

信息取自各项目自己的文档和源代码（2026-10-20；clausona：v0.12.0-beta）。"未说明"表示项目没有提及，而不是没有该功能。

| 工具 | 方式 | 同时使用多个账号 | 5 小时和每周用量 | 达到限制时 | 平台 |
| --- | --- | --- | --- | --- | --- |
| ccam | 单独的文件夹 | 是 | 是，`ccam usage` | 会话中弹出通知；`ccam next --auto --continue` 或 `/ccam:switch` 会把对话带到余量最多的账号，无论是否共享历史；`ccam --fallback` 会自动在那里重新运行 `-p` 任务，`ccam pick` 把并行任务分散到各个账号 | Linux、macOS、Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | 替换登录信息；用 `cswap run` 时使用单独的文件夹 | 是，使用 `cswap run` | 是 | `cswap auto` 在达到阈值时切换登录信息；正在运行的会话继续进行 | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 单独的文件夹 | 是 | 部分支持 | 账号池通过本地代理把请求转给下一个账号 | Linux、macOS、Windows |
| [clausona](https://github.com/larcane97/clausona) | 单独的文件夹 | 是 | 是 | 每次运行前由路由选出还有配额的账号，可以等待重置，并给各个并行任务分配不同的账号 | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 替换登录信息；用 `clauth start` 时使用单独的文件夹 | 是，使用 `clauth start` | 是 | 沿备用链依次切换 | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 替换登录信息，或使用隔离的配置 | 是，使用隔离的配置 | 是 | `caam run` 在另一个账号上重试非交互式命令 | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | 替换登录信息并包装 `claude` | 否，只有一个活动账号 | 是 | 替换登录信息并恢复同一个对话 | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 替换登录信息 | 否 | 是 | 可选在达到阈值时切换 | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 替换登录信息 | 否 | 是 | 可选在 100% 时切换 | macOS |

ccswap 是 claude-swap 的一个分支，还能管理 Codex 的登录。[claude-account](https://github.com/hamzarehmandeveloper/claude-account)、[claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) 和 [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) 也为每个账号提供单独的文件夹。其中一些工具也支持其他编程代理，例如 Codex 或 Gemini。

## 各工具的实现

| 工具 | 编写语言 | Claude Code 运行时 | 是否在不同位置之间复制登录信息 | 账号之间可以共享什么 |
| --- | --- | --- | --- | --- |
| ccam | Go，单个静态二进制文件 | 什么也没有：ccam 用 `claude` 替换自身（Windows 无法这样做，ccam 会等待它结束；可选的中继进程会留在它旁边） | 否 | 组：每个组指定它共享的项和它的账号；一个账号可以属于多个组 |
| claude-swap | Python 3.12+，用 uv 或 pipx 安装 | 用 `cswap run` 时什么也没有；`cswap auto` 和菜单栏应用会持续运行 | 是，在替换登录信息时 | 用 `cswap run` 时：设置、`CLAUDE.md` 和扩展一起共享，历史记录可选；MCP 服务器在每次启动时复制 |
| CCS | Node.js 18+ 或 Bun | CCS 作为父进程留着；账号池会加一个本地代理 | 账号文件夹不复制；账号池保存自己的登录信息 | 每个账号选择共享或本地，历史记录按上下文分组；MCP 服务器会被复制 |
| clausona | TypeScript，Node.js 20+ | 直接运行的 `claude` 经过一个 shell 函数；`clausona run` 让 Node 作为父进程留着；每次退出后有一个后台 Node 进程 | 否，但它会自行为闲置账号的令牌续期，以读取它们的用量 | 除每个账号自己的文件外，全部来自一个主账号；历史记录可选；用 `claude mcp add` 添加的服务器不共享 |
| clauth | Rust | clauth 作为父进程留着 | 是，在替换登录信息时 | 整个 `~/.claude`，或者什么都不共享 |
| CAAM | Go | 取决于模式；可选的守护进程 | 是，在其保管库模式下 | 由模式决定；MCP 服务器在每次启动时复制 |
| cux | Go | 一个包装器始终包在 `claude` 外面 | 是 | 全部：只有一个文件夹 |
| CCSwitcher, Claude Switcher | macOS 应用（Swift、Python） | 菜单栏应用 | 是 | 全部：只有一个文件夹 |

## ccam 的优势

### 登录信息留在 Claude Code 存放的位置

ccam 从不把 Claude 的登录信息从一个地方复制到另一个地方。每个账号的登录信息都留在自己的文件夹中，在 macOS 上则留在该文件夹专属的钥匙串项目中，由 Claude Code 负责续期。不存在会失去同步的保存副本。用 `ccam import` 注册的现有 `~/.claude` 继续使用普通 `claude` 所用的钥匙串项目，所以无论用哪种方式启动，都保持登录。

只有在你运行 `ccam usage --refresh` 时，ccam 才会自己续期过期的令牌。它会先加锁，并在该账号正在运行时拒绝续期，所以绝不会在活动会话底下轮换令牌。

### 你和 Claude Code 之间没有任何东西在运行

ccam 读取一个文件，检查它的链接，然后用 `claude` 替换自身（`exec`）。你工作时没有 ccam 进程、包装器、代理或守护进程：标准输入、信号和退出码都属于 Claude Code，它的请求直接发往 Anthropic。Windows 无法替换进程，所以在 Windows 上 ccam 用同一个控制台启动 Claude Code，并等待它的退出码，除此之外什么也不做。

以这种方式启动账号大约只多花 4 ms，第 95 百分位低于 6 ms，峰值内存 8 MB（用 hyperfine 运行 [`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) 300 次）。在同一台笔记本电脑上，仅启动 Python 解释器就要 13 ms，Node.js 要 31 ms，这还是在程序做任何事情之前。

### 由组决定共享什么

MCP 服务器、`settings.json`、`CLAUDE.md`、扩展（agents、commands、skills、插件）和对话历史是五个独立的项目。组指定它共享的项以及其中的账号，所以两个工作账号可以在 `work` 组中共享历史记录，个人账号保留自己的历史记录，而三个账号都通过 `default` 使用同样的 MCP 服务器。`ccam group set work --add memory` 会在一个计划中移动该组的所有账号；[仪表盘](/features/dashboard)用复选框完成同样的操作，并把账号已经从另一个组获得的项显示为灰色。

- MCP 服务器在启动时通过 `--mcp-config` 传入，所以默认情况下 ccam 从不编辑 `.claude.json`，也就是 Claude Code 运行时会改写的那个文件。`ccam mcp add` 一次就能为整个组添加一个服务器。
- 其他项目是指向组文件夹的链接，所以在一个账号中做的修改会立即出现在其他账号中。启动时不复制任何东西。
- 登录信息、`.claude.json` 和运行状态从不共享。

### 每项更改都先展示并备份

`ccam group` 和 `ccam migrate` 会先打印计划，并在移动任何东西之前征求确认。当某个文件在账号和组之间不同时，由你选择：保留组的版本、使用你的版本，或者合并 JSON。原始文件放到 `.ccam-backup/<time>/`，被移除账号的文件夹放进回收文件夹，文件通过一次重命名完成替换，所以写入中断时旧文件仍然保留。`--dry-run` 会显示启动、共享、迁移、快捷命令、更新和修复将做什么。`ccam doctor --fix` 会修复损坏的链接和权限。

### 在下一个账号中继续同一个对话

`ccam next --auto --continue` 会同时查询所有账号的用量，打开余量最多的账号并继续对话。共享历史时，这就是同一个对话；不共享时，Claude Code 会把它复制到新账号中（`--resume <file> --fork-session`），ccam 会告诉你第一次回复要重新读取多少 token。当用量限制使会话停止时，ccam 插件的钩子会显示一条通知，`/ccam:switch` 会把对话交接过去。何时切换由你决定；没有任何东西在后台监视你的用量，自动选择也只在同一个组织内进行，所以公司的对话绝不会意外落到私人账号上。账号用量达到 90% 时，状态栏会提示这个命令。

脚本无法等你操作，所以 `ccam --fallback work -p "…"` 会替它们切换：当账号达到限制时，同一个命令会以余量最多的账号重新运行，并且只打印最终完成的那次运行的输出。不是用量限制造成的失败会原样传递。对于并行任务，`ccam pick` 会给每个任务分配一个还有余量的不同账号；当没有账号还有余量时，两者都以退出码 75 退出。

### 每个文件夹使用合适的账号

`ccam folder work ~/code/company` 让该文件夹中所有不指定账号名称的启动都使用 `work`，而你的其他项目继续使用默认账号。这个设置保存在每次启动都会读取的同一个 `config.toml` 中，所以不会增加任何时间。

### 单个二进制文件，经过校验的下载

ccam 是一个约 10 MB 的静态二进制文件，无需安装任何运行时，支持 Linux、macOS 和 Windows。安装脚本（`install.sh`、`install.ps1`）和 `ccam update` 会用发布版本的 SHA-256 校验和检查每个下载，ccam 也从不未经询问就更新。它没有遥测：它发出的请求只有向 Anthropic 查询你的账号用量（并行，3 秒超时，60 秒缓存）、仪表盘打开时向 GitHub 询问最新版本，以及你要求的下载。使用 `usage_source = "claude"` 时，连这些请求也通过 Claude Code 发出，ccam 从不接触令牌。

## 其他工具做得更多的地方

- claude-swap、clauth、cux、CAAM 和 CCSwitcher 可以在账号超过用量阈值时自行切换正在运行的交互式会话。ccam 用 `--fallback` 切换 `-p` 运行，交互式会话则留给你，并在达到限制时发出通知、提供 `/ccam:switch`。
- 在未开启开发人员模式的 Windows 上，ccam 可以共享文件夹，但无法共享 `settings.json` 这样的单个文件；claude-swap 和 CCS 在这种情况下会改为复制它们。
- CCSwitcher、Claude Switcher 和 claude-swap 会在 macOS 菜单栏中显示用量。
- CCS、clausona、CAAM、clauth、ccswap 和 Claude Switcher 还能管理 Codex 和其他编程代理。
- clausona 还能按账号开启或关闭扩展（ccam 的 `ccam ext` 只是显示它们），记录本地的令牌费用，并提供在 tmux 等终端中运行一批并行任务的 Claude Code 插件。
- CAAM、CCSwitcher 和 Claude Switcher 可以用 Homebrew 安装，clauth 和 CAAM 会为发布的校验和签名。ccam 通过安装脚本、`.deb` 或 `go install` 安装，校验校验和但没有签名。

## 哪种适合你？

如果你一次只用一个账号，或者希望正在运行的会话自动切换，选择替换登录信息的工具。如果你想让多个账号并行运行、共用一套配置、在下一个账号中继续同一个对话，并且不希望有任何东西在 Claude Code 旁边运行，无论在 Linux、macOS 还是 Windows 上，都选择 ccam。如果什么都不需要共享，[普通的 shell 别名](/compare/aliases)就够了。

是否可以使用多个账号取决于每个账号的条款。使用前请先确认。
