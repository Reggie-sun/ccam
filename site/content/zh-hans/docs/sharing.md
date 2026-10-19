---
title: 共享 MCP 服务器、设置和历史记录
description: ccam 如何在 Claude Code 账号之间共享 MCP 服务器、settings.json、CLAUDE.md、agents、skills 和对话历史，登录则各自独立。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

每个账号都保留自己的登录。其他所有内容都可以逐项共享：ccam 把账号中的文件或文件夹链接到一个共享组文件夹，因此在一个账号中所做的更改会出现在该组的所有账号中。

## 可以共享哪些内容

| 共享项 | 选项 | 共享的内容 |
| --- | --- | --- |
| MCP 服务器 | `--mcp` | 组中的 `mcp.json`，在启动时传给 Claude Code |
| 设置 | `--settings` | `settings.json`：权限、hooks、模型、状态栏 |
| 全局记忆 | `--memory` | `CLAUDE.md` |
| 扩展 | `--extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| 对话历史 | `--history` | `projects/` 和 `history.jsonl` |

新账号默认共享除对话历史以外的所有内容。

登录文件 `.credentials.json` 和 `.claude.json` 永远不会被共享。ccam 只链接一份固定列表中的路径，而这两个文件不在列表中。

## 组

`on` 表示 `default` 组。你也可以指定一个组名，让只有部分账号共享某一项：

```sh
ccam share work --history team
ccam share side --history team
```

现在 `work` 和 `side` 彼此共享对话，但不与你的其他账号共享。

## 开启共享

当某个账号开始共享某一项时，ccam 会：

1. 打印一份计划，列出它将移动、复制或链接的每个文件，并在开始前征求你的确认；
2. 如果组文件夹仍为空，就把账号的文件移入组文件夹；否则只添加组中还没有的文件；
3. 当同一个文件的内容不同时，询问你是保留组中的版本、使用你的版本，还是合并 JSON；
4. 把你的原始文件保存在账号文件夹内的 `.ccam-backup/<time>/` 中；
5. 用指向该组的链接替换原始文件。

对话文件以会话 ID 命名，因此多个账号的历史记录可以合并而不会冲突。

## 关闭共享

```sh
ccam share work --history off
```

链接会被移除，该账号会得到组中当前内容的一份副本。如果想从空白开始，请加上 `--fresh`。组会为其他账号保留其内容。

## 复制现有配置

如果你在使用 ccam 之前就用过 Claude Code，你的设置都在 `~/.claude` 中。`ccam migrate` 会把它们复制到一个组中，让共享这些项的每个账号都能用上：

```sh
ccam migrate
```

它会复制：

- `settings.json` 和全局 `CLAUDE.md`；
- 扩展 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/` 和 `rules/`；
- 你的用户 MCP 服务器，即 `~/.claude.json` 中的 `mcpServers` 条目，复制到组的 `mcp.json` 中；
- 对话历史（`projects/` 和 `history.jsonl`），但仅在你加上 `--history` 时复制。

源文件夹只会被读取：其中的任何内容都不会被移动、链接或修改，直接运行 `claude` 也会像以前一样工作。登录文件 `.credentials.json` 和 `.claude.json` 从不复制。

ccam 在复制任何内容之前，会先显示计划并征求确认：

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

最后，ccam 会列出通过该组共享每一项的账号。如果有账号尚未共享，它会显示为这些账号开启共享的 `ccam share` 命令。

如果组中已有内容不同的同名文件，ccam 会像开启共享时一样询问你：保留组中的副本，还是使用你的副本。`settings.json` 等 JSON 文件还可以合并，`history.jsonl` 可以按行合并。组中原有的副本会保存在组文件夹内的 `.ccam-backup/<time>/` 中。

| 选项 | 作用 |
| --- | --- |
| `--settings`、`--memory`、`--extensions`、`--mcp`、`--history` | 只复制你指定的项；不指定时复制除对话历史以外的所有内容 |
| `-g`、`--group <group>` | 复制到这个组，而不是 `default` |
| `-y`、`--yes` | 不询问直接复制；文件内容不同时保留组中的副本 |

在命令前加上 `--dry-run`，例如 `ccam --dry-run migrate`，即可只查看计划而不复制任何内容。要从其他文件夹复制，请指定该文件夹：`ccam migrate ~/.claude-work --group work`。

## MCP 服务器

Claude Code 把 MCP 服务器和登录信息一起存储在 `.claude.json` 中，所以这个文件不能共享。ccam 把共享的服务器保存在组的 `mcp.json` 中，并用 `--mcp-config=<group>/mcp.json` 启动 Claude Code。你自己的 `--mcp-config` 选项仍会在此基础上生效。

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

在会话中用 `claude mcp add` 添加的服务器只为该账号保存。下次启动该账号时，ccam 会列出这个服务器的名称，以及把它移到组中的命令：

```sh
ccam mcp import work gh
```

通过 OAuth 登录的 MCP 服务器会把令牌和账号的登录信息存放在一起，因此需要在每个账号中分别授权一次。

如果你希望这些服务器作为普通的用户级服务器出现在 `claude mcp list` 和 `/mcp` 中，请在 `~/.config/ccam/config.toml` 中设置 `mcp_strategy = "sync"`。这样 ccam 会在启动前把它们写入每个账号的 `.claude.json`，写入时使用 Claude Code 自己的文件锁，并保留备份。

## 同一文件夹中的两个账号

共享对话历史的账号会写入相同的文件。如果你在某个文件夹中启动第二个 Claude Code，而同组的另一个账号已经在该文件夹中运行，ccam 会发出警告并询问是否继续。

## 出现问题时

```sh
ccam doctor
ccam doctor --fix
```

如果某个共享路径变成了普通文件（例如在你手动恢复备份之后），ccam 会在启动 Claude Code 之前停止，并提示你运行 `ccam doctor --fix`；该命令会把这个文件合并到组中，并保留备份。
