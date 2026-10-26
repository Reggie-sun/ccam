---
title: ccam 与 shell 别名对比
description: 逐项对比两种运行多个 Claude Code 账号的方式：使用 CLAUDE_CONFIG_DIR 的 shell 别名，以及 ccam。
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "c30d84bceba0"
---

为每个账号设置一个 shell 别名，是运行多个 Claude Code 账号最简单的方法。ccam 在底层做的是同一件事，并增加了共享、切换和用量跟踪。

## 别名方式

在 `~/.bashrc` 或 `~/.zshrc` 中为每个额外账号添加一行，打开新的终端，然后运行一次该别名来登录：

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` 仍然是你的第一个账号，`claude-work` 启动第二个账号。这正是 [Claude Code 文档所介绍的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。

## ccam 方式

安装 ccam 并添加账号，这里还为它设置了一个专属的快捷命令：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` 和 `ccwork` 都会启动该账号，而且两者在脚本中同样可用。

## 并排对比

| | Shell 别名 | ccam |
| --- | --- | --- |
| 独立的登录 | 是 | 是 |
| 每个账号的设置步骤 | 编辑 `~/.bashrc`，然后登录 | `ccam add` |
| 能否在脚本和其他程序中使用 | 否，别名只在 shell 中有效 | 是，`ccam work` 和 `ccwork` 这样的快捷命令都是真正的命令 |
| 共享设置、`CLAUDE.md`、扩展 | 手动复制或链接 | 是，逐项选择，带备份 |
| 共享 MCP 服务器 | 手动，每个账号分别设置 | 只需运行一次 `ccam mcp add` |
| 在另一个账号中继续对话 | 否 | `ccam next --continue` |
| 查看所有账号的用量 | 否 | `ccam usage` |
| 选择余量最多的账号 | 否 | `ccam next --auto` |
| 达到限制时以另一个账号重新运行 `-p` 任务 | 否 | `ccam --fallback` |
| 每个项目文件夹使用各自的账号 | 否 | `ccam folder` |
| 在状态栏中显示账号名 | 否 | `ccam statusline install` |
| 修复损坏的链接 | 手动 | `ccam doctor --fix` |
| 额外软件 | 无 | 一个静态二进制文件 |

如果你只有两个账号，并且不想在它们之间共享任何内容，就选择别名。如果你希望多个账号用起来像同一套环境，就选择 ccam。你也可以先用别名，以后再迁移到 ccam：`ccam import` 会在原位置注册已有的 `~/.claude-*` 文件夹。

[如何在一台电脑上使用多个账号](/guides/multiple-claude-code-accounts)逐步介绍了这两种方法，[ccam 与其他账号切换工具](/compare/account-switchers)则把 ccam 与那些切换同一份登录信息的工具做了对比。
