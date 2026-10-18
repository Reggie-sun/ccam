---
title: CLAUDE_CONFIG_DIR 详解
description: CLAUDE_CONFIG_DIR 在 Claude Code 中的作用、文件夹中存储了什么、macOS 把登录信息保存在哪里，以及哪些文件可以在账号之间共享。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` 告诉 Claude Code 使用哪个文件夹来代替 `~/.claude`。所有与账号相关的内容都存放在那里：登录信息、设置、MCP 服务器、插件和对话历史。把它指向不同的文件夹，就能得到相互独立的账号。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## 文件夹里有什么

已使用 Claude Code 2.1.302 核实：

| 路径 | 内容 | 可共享 |
| --- | --- | --- |
| `.credentials.json` | 登录令牌（Linux）、MCP OAuth 令牌 | 否 |
| `.claude.json` | 账号详情、用户级 MCP 服务器、各项目的状态 | 否（MCP 服务器另行处理） |
| `settings.json` | 权限、hooks、模型、状态栏、插件 | 是 |
| `CLAUDE.md` | 全局指令 | 是 |
| `agents/`、`commands/`、`skills/`、`plugins/` | 扩展 | 是 |
| `projects/` | 以 `.jsonl` 文件保存的对话、各项目的记忆 | 是 |
| `history.jsonl` | 提示历史 | 是 |
| `sessions/`、`shell-snapshots/`、`cache/` 等 | 运行时状态 | 否 |

## 两个容易踩坑的细节

**不设置 `CLAUDE_CONFIG_DIR` 时，`.claude.json` 不在 `~/.claude` 里面。** 这时 Claude Code 把它保存在 `~/.claude.json`，与该文件夹并列。如果你开始使用 `CLAUDE_CONFIG_DIR=~/.claude`，Claude Code 会改为查找 `~/.claude/.claude.json`，于是找不到你的账号。`ccam import` 通过把这两者链接起来解决这个问题。

**在 macOS 上，登录信息保存在钥匙串中，而不是文件夹里。** Claude Code 以 `Claude Code-credentials` 为名称存储它；设置了 `CLAUDE_CONFIG_DIR` 时，名称后面还会加上文件夹路径的短哈希。因此不同的文件夹会得到不同的钥匙串条目。如果钥匙串拒绝写入，Claude Code 会退而使用文件夹中的 `.credentials.json`。

## 在文件夹之间共享

符号链接是可行的：Claude Code 会通过链接的 `settings.json` 或 `projects/` 文件夹进行读写，并保留链接。MCP 服务器是例外，因为它们与登录信息共用 `.claude.json`；请改用 `--mcp-config` 传入它们。

ccam 正是基于这些规则构建的。关于它如何链接文件夹，参见[共享](/docs/sharing)；想了解全貌，参见[多账号指南](/guides/multiple-claude-code-accounts)。
