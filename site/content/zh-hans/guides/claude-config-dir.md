---
title: CLAUDE_CONFIG_DIR 详解
description: CLAUDE_CONFIG_DIR 在 Claude Code 中的作用、文件夹中存储了什么、macOS 把登录信息保存在哪里，以及哪些文件可以在账号之间共享。
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR` 告诉 Claude Code 使用哪个文件夹来代替 `~/.claude`。所有与账号相关的内容都存放在那里：登录信息、设置、MCP 服务器、插件和对话历史。把它指向不同的文件夹，就能得到相互独立的账号。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR 是官方支持的吗？

是的。Claude Code 在其[环境变量](https://code.claude.com/docs/en/env-vars)列表中把它列为覆盖配置文件夹的方式，并注明它"可用于并行运行多个账号"；其[身份验证指南](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)也用它来登录多个账号。你可以在 shell 或用户设置中设置它，但不能在项目设置中设置。

## 文件夹里有什么？

已使用 Claude Code 2.1.309 核实：

| 路径 | 内容 | 可共享 |
| --- | --- | --- |
| `.credentials.json` | 登录令牌（Linux、Windows）、MCP OAuth 令牌 | 否 |
| `.claude.json` | 账号详情、用户级 MCP 服务器、各项目的状态 | 否（MCP 服务器另行处理） |
| `settings.json` | 权限、hooks、模型、状态栏、插件 | 是 |
| `CLAUDE.md` | 全局指令 | 是 |
| `agents/`、`commands/`、`skills/`、`plugins/` | 扩展 | 是 |
| `projects/` | 以 `.jsonl` 文件保存的对话、各项目的记忆 | 是 |
| `history.jsonl` | 提示历史 | 是 |
| `sessions/`、`shell-snapshots/`、`cache/` 等 | 运行时状态 | 否 |

## 有哪些容易踩坑的地方？

**不设置 `CLAUDE_CONFIG_DIR` 时，`.claude.json` 不在 `~/.claude` 里面。** 这时 Claude Code 把它保存在 `~/.claude.json`，与该文件夹并列。如果你开始使用 `CLAUDE_CONFIG_DIR=~/.claude`，Claude Code 会改为查找 `~/.claude/.claude.json`，于是找不到你的账号。`ccam import` 通过把这两者链接起来解决这个问题。如果 `~/.claude` 中已经有自己的 `.claude.json`（来自之前用 `CLAUDE_CONFIG_DIR=~/.claude` 运行的结果），ccam 会保留这两个文件并说明这一点，因为只有你知道哪一个是正确的。

**在 macOS 上，登录信息保存在钥匙串中，而不是文件夹里。** Claude Code 以 `Claude Code-credentials` 为名称存储它；设置了 `CLAUDE_CONFIG_DIR` 时，名称后面还会加上 `-` 和文件夹路径的 SHA-256 的前 8 个十六进制数字，账户为 `$USER` 所指定的账户。因此不同的文件夹会得到不同的钥匙串条目，而被移动或重命名的文件夹会失去它的登录；这就是 `ccam rename` 在 macOS 上保留文件夹名称的原因。如果钥匙串拒绝写入，Claude Code 会退而使用文件夹中的 `.credentials.json`，[其文档](https://code.claude.com/docs/en/authentication#credential-management)中也有说明。

**并非所有登录方式都保存在文件夹中。** 不使用 API 密钥的 Claude Console 登录保存在文件夹之外，所以单独的文件夹无法把两个这样的登录分开。claude.ai 登录和 API 密钥则按文件夹分别保存。

## 在 Windows 上，这个文件夹在哪里？

默认是 `%USERPROFILE%\.claude`，登录信息保存在其中的 `.credentials.json` 里；当 Claude Code 设置为使用 Windows 凭据管理器时，则保存在凭据管理器中，名称与 macOS 上一样由文件夹路径生成。为每个额外的账号把 `CLAUDE_CONFIG_DIR` 设为另一个文件夹，方法与 Linux 和 macOS 相同。

## 文件夹之间可以共享文件吗？

符号链接是可行的：Claude Code 会通过链接的 `settings.json` 或 `projects/` 文件夹进行读写，并保留链接。MCP 服务器是例外，因为它们与登录信息共用 `.claude.json`；请改用 `--mcp-config` 传入它们。

ccam 正是基于这些规则构建的。关于它如何链接文件夹，参见[共享](/docs/sharing)；想了解全貌，参见[多账号指南](/guides/multiple-claude-code-accounts)。
