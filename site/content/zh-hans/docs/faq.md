---
title: 常见问题
metaTitle: "ccam 常见问题：Claude Code 多账号"
description: 关于 ccam、Claude Code 多账号、用量限制、共享、隐私和支持的操作系统的常见问题解答。
updated: 2026-10-21
order: 6
translation: machine
sourceHash: "bc944d7dcab9"
---

## ccam 是 Anthropic 的官方工具吗？

不是。ccam 是一个非官方的开源项目，与 Anthropic 没有关联，也未获得 Anthropic 的认可。它使用的是 [Claude Code 文档](https://code.claude.com/docs/en/env-vars)中说明的 `CLAUDE_CONFIG_DIR` 设置，用于并行运行多个账号。

## 使用多个 Claude Code 账号是否被允许？

这取决于每个账号适用的条款；如果是工作账号，还取决于你所在组织的政策。在使用多个账号之前，请先确认这些规定。ccam 只是让你在已有的账号之间切换起来更方便。

## ccam 能绕过用量限制吗？

不能。每个账号都保留各自的用量限制。ccam 会显示这些限制，并在某个账号用完时，让你在自己拥有的另一个账号上继续。

## 如何在一台电脑上使用两个 Claude Code 账号？

用 `CLAUDE_CONFIG_DIR` 为每个账号指定单独的配置文件夹。ccam 会替你完成这件事：为每个账号运行 `ccam add`，然后用 `ccam personal` 和 `ccam work`，或者你自己命名的快捷命令（例如 `ccwork`）启动它们。[指南](/guides/multiple-claude-code-accounts)也介绍了如何手动设置。

## 如何在另一个账号上继续对话？

用 `ccam group set default --add history`（或者为它们单独建一个组：`ccam group add pair --items history --members work,side`）在这些账号之间共享对话历史，然后运行 `ccam next --continue`。Claude Code 会在下一个账号中打开当前文件夹最近的一次对话。

## ccam 会查看或发送我的登录信息吗？

ccam 从不在账号之间共享你的 Claude 登录信息。它读取账号的访问令牌，只是为了在你运行 `ccam usage`、`ccam next --auto` 或 `ccam doctor` 时，向 Anthropic 的用量接口查询该账号的用量。`ccam migrate` 会把 MCP 服务器已保存的登录从 `~/.claude` 复制到同一台电脑上你各个账号自己的登录文件中。ccam 没有遥测。除了这些查询，它唯一会自行发出的请求是仪表盘向 GitHub 询问最新的 ccam 版本，其中不包含任何关于你或你的账号的信息。

## 如何在远程服务器上登录？

像平常一样在 SSH 会话中运行 `ccam add` 或 `ccam login <name>`。Claude Code 在服务器上等待登录，所以你电脑上的浏览器最后会停在一个无法加载的 `http://localhost:<port>/callback?code=…` 页面。从浏览器复制这个地址并粘贴到终端中：ccam 会把它交给服务器上的 Claude Code，由它完成登录。这个地址包含一次性代码，请不要分享给他人。

## 我的账号存储在哪里？

在 `~/.local/share/ccam/profiles/<name>/` 中，每个账号一个文件夹。在 macOS 上，Claude Code 把登录信息本身保存在钥匙串中，条目名称由该文件夹推导而来。

## 我现有的 ~/.claude 会怎样？

不会有任何变化，除非你运行 `ccam import` 或 `ccam migrate`。导入会在原位置注册该文件夹，因此直接运行 `claude` 和 `ccam <name>` 都会继续使用同一个账号和同一套设置。`ccam migrate` 会把其中的设置、`CLAUDE.md`、扩展和 MCP 服务器复制到共享设置中，并且只读取该文件夹。

## ccam 支持哪些系统？

Linux（Ubuntu、Debian、Fedora、Arch、Alpine 等）、macOS 以及 Windows 10 和 11（原生运行或在 WSL 中运行）。在 Windows 上，共享 `settings.json` 这样的文件需要开启开发人员模式；共享文件夹则无需开启（[在 Windows 上共享](/features/share#on-windows)）。

## 为什么某个 MCP 服务器只出现在一个账号中？

在会话中用 `claude mcp add` 添加的服务器只为该账号保存。下次启动该账号时，ccam 会提醒你；用 `ccam mcp import <account> <server>` 把服务器移到共享组。

## ccam 会拖慢 Claude Code 吗？

通过 ccam 启动 Claude Code 大约增加 4 毫秒，这是在一台笔记本电脑上用 `scripts/perf-check.sh` 测得的。ccam 读取一个小文件，检查共享链接，然后用 `claude` 替换自身，所以在你工作时它并不在运行。Windows 无法替换进程，所以在 Windows 上 ccam 启动 Claude Code 后会等待它结束，除此之外什么也不做。

## 如何更新 ccam？

运行 `ccam update`。它会下载最新版本，校验其校验和，然后替换 `ccam` 文件；`ccam update --check` 只告诉你是否有更新的版本。仪表盘也会在有更新版本时告诉你，按 `U` 会先询问再安装；ccam 从不自行更新。如果你安装的是 `.deb` 软件包，请改用 apt 安装新的 `.deb`。参见[更新](/docs/install#update)。

## 如何卸载 ccam？

用 `ccam shortcut <account> --remove` 删除每个账号的快捷命令，再从 `~/.local/bin` 中删除 `ccam`；如果不再需要这些账号，再删除 `~/.config/ccam`、`~/.local/share/ccam` 和 `~/.local/state/ccam`。参见[安装](/docs/install#uninstall)。
