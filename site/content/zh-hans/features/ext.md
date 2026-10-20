---
title: 查看扩展
description: ccam ext 显示每个账号加载的技能、命令、插件、MCP 服务器和钩子，每一项的来源，以及哪些地方出了问题。
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` 回答这样的问题："为什么这个技能（或 MCP 服务器）在一个账号中能用，在另一个账号中却不能用？"它为每个账号列出 Claude Code 加载的扩展以及每一项的来源：某个共享组、该账号自己的文件夹，或者你所在的项目。表格下方列出会让扩展无法按预期加载的问题。它只读取。

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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

每个单元格表示该项在该账号中的来源：共享组（`default`）；位于账号自己的文件夹中时为 `account`；账号没有该项时为 `-`。`!` 标记有问题的项，问题在表格下方说明。

## 读取的内容

| 类型 | 位置 | 显示的来源 |
| --- | --- | --- |
| 技能 | `skills/` 中每个含有 `SKILL.md` 的文件夹 | `skills/` 链接到的组，或 `account` |
| 命令、代理 | `commands/` 和 `agents/` 下的 `.md` 文件，子文件夹显示为 `folder/name` | 组，或 `account` |
| 插件 | `plugins/installed_plugins.json` | 组，或 `account` |
| MCP 服务器 | 组的 `mcp.json`（共享 MCP 时）和账号的 `.claude.json` | 组，或 `account` |
| 钩子 | `settings.json` 中的 `hooks`，按事件计数 | 组，或 `account` |
| 当前项目 | 当前文件夹中的 `.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json` 和 `.mcp.json` | 只列出一次，适用于每个账号 |

## 报告的问题

- **损坏的链接**：`skills/`、`commands/` 或 `agents/` 中目标已不存在的链接，或者链接到不存在位置的共享文件夹（`skills`、`commands`、`agents`、`plugins`）（用 `ccam doctor --fix` 修复）。
- **没有 `SKILL.md` 的技能文件夹**：Claude Code 会跳过它。
- **在两个地方定义不同的 MCP 服务器**：组的 `mcp.json` 和账号的 `.claude.json`，或者其中之一和项目的 `.mcp.json`。两处设置相同（`mcp_strategy = "sync"` 就是这样写入的）不算问题。用 `ccam mcp import <account> <server>` 把账号自己的服务器移到组中。
- **不是有效 JSON 的 `installed_plugins.json`。**

## 安全性与边界情况

- `ccam ext` 不做任何更改，也不发出网络请求。它读取服务器名称和钩子数量，从不读取 MCP 定义中 `env` 或请求头的值。
- 它显示的是磁盘上的内容。已安装但在设置中关闭的插件，或被托管设置阻止的服务器，仍会显示出来。
- 在项目文件夹中运行它，即可包含该项目的 `.claude` 文件夹和 `.mcp.json`。

## 相关内容

- [在账号之间共享](/features/share)
- [共享的 MCP 服务器](/features/mcp)
- [检查与修复](/features/doctor)
