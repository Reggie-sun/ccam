---
title: ccam 与 shell 别名对比
description: 并排对比两种运行多个 Claude Code 账号的方式：使用 CLAUDE_CONFIG_DIR 的 shell 别名，以及 ccam。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6fc4200a8a65
---

为每个账号设置一个 shell 别名，是运行多个 Claude Code 账号最简单的方法。ccam 在底层做的是同一件事，并增加了共享、切换和用量跟踪。

| | Shell 别名 | ccam |
| --- | --- | --- |
| 独立的登录 | 是 | 是 |
| 每个账号的设置步骤 | 编辑 `~/.bashrc`，然后登录 | `ccam add` |
| 能否在脚本和其他程序中使用 | 否，别名只在 shell 中有效 | 是，`ccam work` 和 `ccwork` 这样的快捷命令都是真正的命令 |
| 共享设置、`CLAUDE.md`、扩展 | 手动复制或链接 | 逐项设置，带备份 |
| 共享 MCP 服务器 | 手动，每个账号分别设置 | 只需运行一次 `ccam mcp add` |
| 在另一个账号中继续对话 | 否 | `ccam next --continue` |
| 查看所有账号的用量 | 否 | `ccam usage` |
| 选择余量最多的账号 | 否 | `ccam next --auto` |
| 在状态栏中显示账号名 | 否 | `ccam statusline install` |
| 修复损坏的链接 | 手动 | `ccam doctor --fix` |
| 额外软件 | 无 | 一个静态二进制文件 |

如果你只有两个账号，并且不想在它们之间共享任何内容，就选择别名。如果你希望多个账号用起来像同一套环境，就选择 ccam。你也可以先用别名，以后再迁移到 ccam：`ccam import` 会在原位置注册已有的 `~/.claude-*` 文件夹。

[如何在一台电脑上使用多个账号](/guides/multiple-claude-code-accounts)逐步介绍了这两种方法，[ccam 与其他账号切换工具](/compare/account-switchers)则把 ccam 与那些切换同一份登录信息的工具做了对比。
