---
title: 介绍 ccam
description: ccam 为什么存在，它能为拥有多个 Claude Code 账号的人做什么，以及它如何做到不打扰你的工作。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 985e2d369c34
---

ccam 是一个小型命令行工具，适合在同一台电脑上使用多个 Claude Code 账号的人：比如一个个人账号加一个工作账号，或者为一个长期项目准备的多个账号。

## 问题所在

Claude Code 通过 `CLAUDE_CONFIG_DIR` 支持多个账号：每个账号一个文件夹。但实际上，每个文件夹都是一座孤岛。设置、`CLAUDE.md`、MCP 服务器、agents 和 skills 都要在每个文件夹中重新配置，并手动保持同步。对话只留在它开始的地方。而当一个账号达到用量限制时，你只有在尝试使用时才会发现。

## ccam 能做什么

- `ccam add` 创建文件夹、运行登录，并为账号分配一个快捷命令：`claude1`、`claude2`，依此类推。
- `ccam work` 和 `claude2` 以该账号启动 Claude Code，并原样传递所有参数。
- MCP 服务器、设置、`CLAUDE.md` 和扩展通过链接共享，可以逐项、逐个账号地设置。登录信息永远不共享。
- `ccam usage` 显示每个账号的 5 小时用量和每周用量；`ccam next --continue` 在下一个账号中接着进行对话。

## 不打扰你的工作

每次会话前都要运行的工具，不应该带来任何能察觉到的开销。ccam 读取一个配置文件，对每个共享链接只用一次系统调用进行检查，然后用 `claude` 替换自身。它只增加约 4 毫秒，而且在你工作时它并不在运行。

它也会谨慎对待你的文件。每次更改共享设置时，它都会在开始前打印计划，并把原始文件保存在备份文件夹中。`ccam doctor` 检查链接、权限和登录状态，并修复可以安全修复的问题。

## 试一试

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam 是开源的非官方项目：它与 Anthropic 没有关联，也未获得 Anthropic 的认可。你能否使用多个账号，取决于每个账号适用的条款。欢迎在 [GitHub](https://github.com/jhyunwoo/ccam) 上提供反馈和翻译。
