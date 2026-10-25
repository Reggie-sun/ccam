---
title: 快速入门
metaTitle: "快速入门：添加和切换 Claude Code 账号"
description: 用 ccam add 添加两个 Claude Code 账号，为它们设置 ccwork 这样的快捷命令，在它们之间切换，并选择要共享的内容。
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

用 `ccam add` 添加每个账号，然后用 `ccam <name>` 或你自己命名的快捷命令（例如 `ccwork`）启动 Claude Code。大多数人只需要这些。

## 1. 添加第一个账号

```sh
ccam add
```

向导会询问账号名称（例如 `personal` 或 `work`）、该账号要与其他账号共享哪些项，以及是否创建快捷命令（见第 4 步）。然后它会运行 `claude auth login`，打开浏览器让你登录。在远程机器上，浏览器可能停在一个无法加载的 `http://localhost:…/callback` 页面；把这个地址粘贴到终端中，ccam 就会完成登录。完成后你会看到：

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

只运行 `ccam` 会打开[仪表盘](/features/dashboard)：用方向键或账号的编号移到某个账号，然后按回车键。

快捷命令或账号名之后的所有内容都会原样传给 Claude Code，因此你可以继续使用熟悉的选项：

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 选择要共享的内容

默认情况下，MCP 服务器、`settings.json`、全局 `CLAUDE.md`，以及 agents、commands、skills 和 plugins 会通过组 `default` 在账号之间共享。对话历史不共享。如果你想跨账号继续对话，请把它加入这个组：

```sh
ccam group set default --add history
```

ccam 在做任何更改之前，会先显示它将为每个账号移动哪些内容。详情参见[通过组共享](/features/share)。

## 6. 当某个账号达到用量限制时

```sh
ccam usage
ccam next --continue
```

`ccam usage` 显示每个账号的 5 小时用量和每周用量。`ccam next --continue` 在下一个账号中打开同一段对话。在脚本中，`ccam --fallback work -p "…"` 会自动切换。参见[用量限制](/docs/usage-limits)。

**Team 和 Enterprise 席位。** 自动选择（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`，以及你开启后的 [`auto_switch`](/features/next#auto-switch)）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换。Anthropic 帮助中心关于 [Team 套餐](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的说明写道："Each seat on a Team plan is assigned to one named individual, and only that person may use it"（Team 套餐的每个席位都分配给一位指定的个人，只有此人可以使用），以及 "Users can't route their work through another user's account to use that user's remaining usage"（用户不能通过另一位用户的账号完成自己的工作，以使用该用户剩余的用量）；关于 [Enterprise 套餐](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的说明也是如此。因此，在席位之间使用自动选择可能违反 Anthropic 的条款。在这些席位之间是否使用，由你自己决定：在依赖这项功能之前，请先向你所在的组织和 Anthropic 确认。个人套餐（Pro、Max）的情况请参见[常见问题中关于条款的说明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)。

## 已经在使用 ~/.claude？

注册它，而不必重新登录：

```sh
ccam import
```

该文件夹保留在原位置，直接运行 `claude` 也会像以前一样工作。
