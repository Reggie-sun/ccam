---
title: 快速入门
metaTitle: "快速入门：添加和切换 Claude Code 账号"
description: 用 ccam add 添加两个 Claude Code 账号，为它们设置 ccwork 这样的快捷命令，在它们之间切换，并选择要共享的内容。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 745a9904e32a
---

用 `ccam add` 添加每个账号，然后用 `ccam <name>` 或你自己命名的快捷命令（例如 `ccwork`）启动 Claude Code。大多数人只需要这些。

## 1. 添加第一个账号

```sh
ccam add
```

向导会询问账号名称（例如 `personal` 或 `work`）、该账号要与其他账号共享哪些项，以及是否创建快捷命令（见第 4 步）。然后它会运行 `claude auth login`，打开浏览器让你登录。完成后你会看到：

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. 迁移现有设置

如果你以前用过 Claude Code，你的设置、`CLAUDE.md`、扩展和 MCP 服务器都在 `~/.claude` 中。把它们复制到共享设置中，每个账号就都能用上：

```sh
ccam migrate
```

ccam 会先显示要复制的内容并征求确认。它只读取 `~/.claude`，从不复制登录信息；除非你加上 `--history`，否则不会复制对话历史。详情参见[复制现有配置](/docs/sharing#copy-an-existing-setup)。

## 3. 添加第二个账号

```sh
ccam add work
```

在浏览器中用另一个账号登录。

## 4. 切换账号

`ccam work` 会以 work 账号启动 Claude Code。想要更短的命令，可以给账号起一个专属的快捷命令，例如 `ccwork`。选完共享项后，向导会询问：

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

按 `y`，然后输入名称，或者直接按回车键使用建议的名称：

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

登录完成后，ccam 会显示确认信息：

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

如果你跳过了这个问题，随时都可以添加快捷命令：

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

现在，以下任一命令都会以 work 账号启动 Claude Code：

```sh
ccwork
ccam work
ccam 2
```

只运行 `ccam` 会打开一个选择器：用方向键和回车键选择，或者直接按账号的编号。

快捷命令或账号名之后的所有内容都会原样传给 Claude Code，因此你可以继续使用熟悉的选项：

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 选择要共享的内容

默认情况下，MCP 服务器、`settings.json`、全局 `CLAUDE.md`，以及 agents、commands、skills 和 plugins 会在账号之间共享。对话历史不共享。如果你想跨账号继续对话，请为两个账号都开启它：

```sh
ccam share personal --history on
ccam share work --history on
```

ccam 在做任何更改之前，会先显示它将移动哪些内容。详情参见[共享](/docs/sharing)。

## 6. 当某个账号达到用量限制时

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
