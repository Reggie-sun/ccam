---
title: 查看扩展
description: ccam ext 显示每个账号加载的技能、命令、插件、MCP 服务器和钩子，每一项的来源，以及哪些地方出了问题。
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

`ccam ext` 回答这样的问题："为什么这个技能（或 MCP 服务器）在一个账号中能用，在另一个账号中却不能用？"它为每个账号列出 Claude Code 加载的扩展以及每一项的来源：某个共享组、该账号自己的文件夹，或者你所在的项目。表格下方列出会让扩展无法按预期加载的问题。它不做任何更改。

## 用法

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

每个单元格表示该项在该账号中的来源：共享组（`default`）；位于账号自己的文件夹中时为 `account`；账号只为当前文件夹拥有的 MCP 服务器为 `local`；账号没有该项时为 `-`。已安装但未启用的插件会在来源后显示 `(off)`。`!` 标记有问题的项，问题在表格下方说明；不属于任何账号的问题（例如只有直接运行的 `claude` 才有的服务器）放在最后，不带名称。`ccam ext --json` 输出同样的报告（`accounts`、当前项目、`problems`），Claude Code 通过网络访问的 MCP 服务器带有 `remote`；没有任何账号时，它输出 `accounts: []` 的报告。

## 读取的内容

| 类型 | 位置 | 显示的来源 |
| --- | --- | --- |
| 技能 | `skills/` 中每个含有 `SKILL.md` 的文件夹 | `skills/` 链接到的组，或 `account` |
| 命令、代理 | `commands/` 和 `agents/` 下的 `.md` 文件，子文件夹显示为 `folder/name` | 组，或 `account` |
| 插件 | `plugins/installed_plugins.json`，由账号和项目设置中的 `enabledPlugins` 启用或关闭；只为某个项目安装的插件只出现在该项目中 | 组，或 `account`；未启用时为 `(off)` |
| MCP 服务器 | 组的 `mcp.json`（共享 MCP 时）、账号的 `.claude.json`（它自己的服务器，以及它只为当前文件夹拥有的服务器），以及其插件的服务器，显示为 `plugin:<plugin>:<server>` | 组、`account` 或 `local` |
| 钩子 | `settings.json` 中的 `hooks`，按事件计数 | 组，或 `account` |
| 当前项目 | 当前文件夹中的 `.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json` 和 `.mcp.json` | 只列出一次；尚未信任该文件夹的账号会因其 `.mcp.json` 服务器而得到一个问题 |

## 报告的问题

- **损坏的链接**：`skills/`、`commands/` 或 `agents/` 中目标已不存在的链接，或者链接到不存在位置的共享文件夹（`skills`、`commands`、`agents`、`plugins`）（用 `ccam doctor --fix` 修复，它会重新创建缺少的组文件夹）。
- **没有 `SKILL.md` 的技能文件夹**：Claude Code 会跳过它。
- **在两个地方定义不同的 MCP 服务器**：组的 `mcp.json` 和账号的 `.claude.json`，或者其中之一和项目的 `.mcp.json`。两处定义完全相同（`mcp_strategy = "sync"` 就是这样写入的）不算问题。用 `ccam mcp import <account> <server>` 把账号自己的服务器移到组中。
- **不是有效 JSON 的 `installed_plugins.json`**、`enabledPlugins` 中没有列出的插件，以及文件夹不存在的插件。
- **只有一个账号加载的服务器**：只在该账号的 `.claude.json` 中、其组内其他账号不加载的服务器，以及位于该账号针对当前文件夹的本地作用域中的服务器（不加 `-s user` 时 `claude mcp add` 就是这样做的）。每一项都附有共享它的 `ccam mcp import` 命令。直接运行的 `claude` 所用 `~/.claude.json` 中没有任何账号加载的服务器，附有 `ccam migrate --mcp` 命令。
- **等待批准的项目服务器**：该账号尚未信任此文件夹，所以 Claude Code 会询问一次，在此之前项目 `.mcp.json` 中的服务器处于等待状态。问题会说明在哪里批准它们：设置组中各账号的（或该账号自己的）`settings.json` 中的 `"enabledMcpjsonServers"`。对 git 仓库主工作树给予的信任也适用于它的其他工作树。
- **登录和连接器**：在此账号中需要登录、而另一个账号已登录的 MCP 服务器（登录只对一个账号有效；用 `/mcp` 登录），以及其他账号用过而此账号未用过的 claude.ai 连接器。

## 安全性与边界情况

- `ccam ext` 不做任何更改，也不发出网络请求。它读取服务器名称和钩子数量，从不读取 MCP 定义中 `env` 或请求头的值。为了检查登录，它会读取每个账号登录中的 MCP 登录名称（在 macOS 上从钥匙串读取），从不读取令牌，并且只针对加载了另一个账号也加载的服务器的账号。
- 它显示的是磁盘上的内容。被托管设置阻止的服务器仍会显示出来；在设置中关闭的插件显示为 `(off)`。
- 在项目文件夹中运行它，即可包含该项目的 `.claude` 文件夹和 `.mcp.json`。

## 相关内容

- [通过组共享](/features/share)
- [共享的 MCP 服务器](/features/mcp)
- [检查与修复](/features/doctor)
