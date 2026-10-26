---
title: 复制现有配置
description: 用 ccam migrate 把现有 ~/.claude 中的设置、CLAUDE.md、扩展和 MCP 服务器复制到 ccam 的组中。
updated: 2026-10-26
order: 6
translation: machine
sourceHash: "f835fbe84ade"
---

如果你在使用 ccam 之前就用过 Claude Code，你的配置都在 `~/.claude` 和 `~/.claude.json` 中。`ccam migrate` 会把它们复制到一个共享组中，让通过该组共享这些项的每个账号一开始就拥有你现有的配置。源文件夹只会被读取，直接运行 `claude` 也会像以前一样工作。在第一次 `ccam add` 之后运行一次；如果 `~/.claude.json` 中有 MCP 服务器，第一次 `ccam add` 会在提示中列出它们。

## 用法

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

在终端中，如果不带项目选项也不带 `--yes`，ccam 会先显示一个清单：

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space（或 `x`）切换光标所在的项，方向键或 `j`/`k` 移动光标，Enter 显示计划，Esc、`q` 或 Ctrl+C 会打印 `Cancelled.`。源中没有内容的项会标记为 `not found`。如果你取消勾选所有项，ccam 会打印 `Nothing selected.` 并停止。

然后 ccam 打印计划，并在复制之前征求确认：

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

组只影响它自己的账号，而且只涉及它共享的项，因此摘要会列出组中的账号。如果组不共享某个被复制的项，或者有些账号不在组中，ccam 会打印用来更改这些的 `ccam group set <group> --add <items>` 和 `ccam group join <group> <accounts>` 命令。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `[folder]` | 要从中复制的 Claude Code 文件夹 | `~/.claude` |
| `--settings` | 复制 `settings.json` | 关闭 |
| `--memory` | 复制全局 `CLAUDE.md` | 关闭 |
| `--extensions` | 复制 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 关闭 |
| `--mcp` | 复制用户 MCP 服务器（`.claude.json` 中的 `mcpServers`） | 关闭 |
| `--history` | 复制对话历史（`projects/`、`history.jsonl`） | 关闭 |
| `-g`、`--group <group>` | 要复制到的共享组 | `default` |
| `-y`、`--yes` | 不询问直接复制；发生冲突时保留组中的副本 | 关闭 |
| `--dry-run`（全局） | 打印计划，并在任何询问或复制之前停止 | 关闭 |

在终端中，如果你没有给出任何项目选项，就在清单中选择要复制的项。不在终端中或带有 `--yes` 时，ccam 复制除对话历史以外的所有内容。默认不包括对话历史，因为 `projects/` 可能很大。

`--group on` 表示 `default`。`--group off` 和 `--group local` 会被拒绝，并显示 `--group takes a group name, not "off"`。`ccam migrate` 支持 `--dry-run`，试运行时清单依然会出现，方便你预览所选内容。

## 工作原理

ccam 在复制任何内容之前会先构建完整的计划。它选定项目（`internal/cli/migrate.go`），然后读取源和组并构建计划（`migrate.Build`）。在你确认之后、执行 `Plan.Apply` 之前，不会写入任何内容。

1. **找到源。** ccam 会拒绝不存在的路径、不是文件夹的路径，以及不包含任何已知 Claude Code 条目的文件夹（`checkSource`）。如果只复制 MCP 服务器，只要 `~/.claude.json` 中有服务器，`~/.claude` 为空或不存在也没关系。ccam 找到属于该文件夹的 `.claude.json`（`migrate.ClaudeJSON`）：在文件夹内；或者对于内部没有该文件的 `~/.claude`，则是它旁边的 `~/.claude.json`。未设置 `CLAUDE_CONFIG_DIR` 时，Claude Code 就把它放在那里。
2. **把项目映射到路径。** 每个路径都来自 `share.Layout`，也就是组所使用的同一份允许列表。MCP 服务器通过 `mcp.LocalServers` 从 `.claude.json` 中读出，它使用 `profile.TopLevelValue` 扫描器，而不是解析整个文件。
3. **只复制，不移动。** 文件通过 `fsutil.AtomicWrite` 或 `fsutil.CopyFile` 写入 `~/.local/share/ccam/groups/<group>/`，后者会在目标旁边写入一个临时文件再重命名，因此被中断的复制永远不会留下不完整的文件。MCP 服务器在锁下写入组的 `mcp.json`：另一个 ccam 进程在此期间添加的服务器会保留，为 [Remote Control](/guides/remote-control) 制作的副本会在之后更新。

每一项都把一个源路径映射到一个组路径：

<!-- motion:migrate-items -->

两行历史记录是虚线，因为只有加上 `--history` 才会复制它们。`.credentials.json` 和 `.claude.json` 列在 `share.Never` 中，永远不会进入组：如果它们出现在布局中，ccam 会以内部错误停止（`Build`），遍历文件夹时也会按名称跳过它们。

ccam 为每个路径确定一个状态，计划中的每一行描述的就是这个状态（`buildEntry`）：

<!-- motion:migrate-states -->

- **Missing**：源中没有，或者应为文件夹的位置是一个文件。每项列出一次，显示为 `not found, skipped`。
- **Linked**：源路径已经解析到组中的副本，例如当源是 ccam 的账号文件夹时。
- **Same**：组中已有完全相同的内容。
- **New**：组中没有，因此整体复制。
- **Merge**：两边都有的文件夹。ccam 只添加组中没有的文件（`copyMissing`）；内容不同的文件保持组中的版本，计划中显示 `group copies kept`。对于 `plugins/`，注册表 `installed_plugins.json` 和 `known_marketplaces.json` 会被合并，插件的绝对路径会被改写为指向组的文件夹。
- **Conflict**：两边都有但内容不同的文件。你可以逐个文件选择：`keep group`、`use mine`，对于 `.json` 文件还有 `merge JSON`（深度合并，值冲突时组中的值优先），对于 `history.jsonl` 还有 `merge by line`。定义不同的 MCP 服务器也属于冲突，可选 `keep group` 或 `use mine`。列出的第一个选项是按 Enter 时的默认值；`--yes` 总是保留组中的副本。

清单本身是 `tui.MigratePicker`，基于 `ccam add` 使用的 `checklist` 组件构建；ccam 通过 lstat 调用和 MCP 服务器列表标记不存在的项（`migrate.Found`）。

### MCP 登录

通过 OAuth 登录的 MCP 服务器把令牌保存在登录文件 `.credentials.json` 的 `mcpOAuth` 下，而不是服务器列表中。ccam 不会复制它们。复制过去的登录与原登录持有同一个刷新令牌，而服务器会在第一次续期时轮换该令牌：此后其他所有持有者（包括源 `~/.claude`）都会退出该服务器的登录。因此，对于源已登录的每个复制过去的服务器，计划中的 `sign-ins` 一行会列出组中需要自己用 `/mcp` 登录的账号，摘要也会再说明一次。没有令牌的条目不算登录，不会列出；已经有自己登录的账号，以及定义由组保留的服务器，也不会列出。不会写入任何登录文件。

## 安全性与边界情况

- 源只会被读取。其中的任何文件都不会被移动、链接或修改，源的登录也从不会被写入，包括在可能从钥匙串读取登录的 macOS 上。
- 在组中的文件被替换之前（无论是因为 `use mine`、合并，还是现有 `mcp.json` 中加入了新的 MCP 服务器），组中的副本会保存在 `<group>/.ccam-backup/<time>/` 中。每次运行的文件夹名都是唯一的，因此同一秒内的两次运行不会覆盖彼此的备份。摘要会打印这个路径。
- 使用 `--yes` 时，每个冲突都保留组中的副本，ccam 会告诉你如何改为自己选择：不带 `--yes` 运行同一命令。
- 如果没有任何需要更改的内容，ccam 会打印 `Nothing to copy.`；如果每个冲突都保留组中的副本，则打印 `Nothing to copy: the group keeps its own copies.`
- 如果复制中途失败，ccam 会停止并显示 `The source folder was not changed. Group files replaced so far are saved in`，后面跟着备份文件夹。再次运行会复制那些没有完整复制的文件。
- 当有多个组共享 MCP 时，`ccam doctor`、`ccam ext` 和 `ccam mcp list` 建议的命令为 `ccam migrate --mcp --group <a|b>`。
- 源不存在时，会以 `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` 结束。

## 相关内容

- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [检查与修复](/features/doctor)
- [在一台电脑上使用多个 Claude Code 账号](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR 详解](/guides/claude-config-dir)
