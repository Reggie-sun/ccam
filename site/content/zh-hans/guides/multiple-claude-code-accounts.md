---
title: 如何在一台电脑上使用两个 Claude Code 账号
description: 借助 CLAUDE_CONFIG_DIR 实现 Claude Code 多账号并行运行：可以用 shell 别名手动设置，也可以用 ccam 自动完成。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

为每个账号指定单独的配置文件夹。Claude Code 从 `CLAUDE_CONFIG_DIR` 环境变量中读取文件夹位置，每个文件夹都保存各自的登录信息、设置和历史记录。你可以用 shell 别名手动设置，也可以交给 ccam 完成。

## 手动设置：使用别名

这是 Claude Code 文档中介绍的方法。在 `~/.bashrc` 或 `~/.zshrc` 中为每个额外的账号添加一行：

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

打开一个新终端并运行 `claude-work`。因为文件夹是新的，Claude Code 会要求你登录。从此以后，`claude` 是你的第一个账号，`claude-work` 是第二个。

这种方法适合两个账号。需要注意：

- 每个文件夹一开始都是空的：没有设置、没有 `CLAUDE.md`、没有 MCP 服务器，也没有 agents 或 skills。你需要复制或链接所需的内容，并自己保持这些副本同步。
- 别名只存在于交互式 shell 中。脚本、Makefile 和其他程序看不到它们。
- MCP 服务器保存在 `.claude.json` 中，与登录信息在一起，所以不能简单地在文件夹之间链接。
- 对话只留在开始时所用的账号中，所以你无法在另一个账号中继续对话。
- 没有办法一次查看所有账号的用量。

## 使用 ccam

ccam 把同样的思路自动化，并补上了缺少的部分：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

现在，`claude1` 和 `claude2`，或者 `ccs personal` 和 `ccs work`，都可以启动这两个账号。因为它们是真正的程序而不是别名，所以在脚本中也能使用。设置、`CLAUDE.md`、MCP 服务器和扩展从一开始就共享；你还可以开启对话历史共享，用 `ccs next --continue` 在另一个账号上继续对话。

如果你已经在使用 `~/.claude`，请运行 `ccam import`，而不是再次添加该账号。

## 该选哪一种

如果只有两个账号且不需要共享任何内容，用别名就够了。如果你想要共享设置和 MCP 服务器、使用更多账号、一眼看清用量或继续对话，ccam 可以省去这些繁琐的管理工作。[对比](/compare/aliases)页面并排列出了两者的差异。
