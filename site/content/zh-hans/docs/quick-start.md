---
title: 快速入门
description: 用 ccam add 添加两个 Claude Code 账号，用 claude1、claude2 或 ccam <name> 在它们之间切换，并选择要共享的内容。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: c6b1c1c3ebaf
---

用 `ccam add` 添加每个账号，然后用 `claude1`、`claude2` 或 `ccam <name>` 启动 Claude Code。大多数人只需要这些。

## 1. 添加第一个账号

```sh
ccam add
```

向导会询问账号名称（例如 `personal` 或 `work`），以及该账号要与其他账号共享哪些项。然后它会运行 `claude auth login`，打开浏览器让你登录。完成后你会看到：

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

ccam 目前用韩语输出自己的消息；命令和选项在所有语言中都相同。

## 2. 添加第二个账号

```sh
ccam add work
```

在浏览器中用另一个账号登录。它会成为 `claude2`。

## 3. 切换

以下任一命令都会以 work 账号启动 Claude Code：

```sh
claude2
ccam work
ccam 2
```

账号名之后的所有内容都会原样传给 Claude Code，因此你可以继续使用熟悉的选项：

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

只运行 `ccam` 会打开一个选择器：用方向键和回车键选择，或者直接按账号的编号。

## 4. 选择要共享的内容

默认情况下，MCP 服务器、`settings.json`、全局 `CLAUDE.md`，以及 agents、commands、skills 和 plugins 会在账号之间共享。对话历史不共享。如果你想跨账号继续对话，请为两个账号都开启它：

```sh
ccam share personal --history on
ccam share work --history on
```

ccam 在做任何更改之前，会先显示它将移动哪些内容。详情参见[共享](/docs/sharing)。

## 5. 当某个账号达到用量限制时

```sh
ccam usage
ccam next --continue
```

`ccam usage` 显示每个账号的 5 小时用量和每周用量。`ccam next --continue` 在下一个账号中打开同一段对话。参见[用量限制](/docs/usage-limits)。

## 已经在使用 ~/.claude？

注册它，而不必重新登录：

```sh
ccam import
```

该文件夹保留在原位置，直接运行 `claude` 也会像以前一样工作。
