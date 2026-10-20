---
title: 复制现有配置
description: 用 ccam migrate 把现有 ~/.claude 中的设置、CLAUDE.md、扩展、MCP 服务器和 MCP 登录复制到 ccam 的组中。
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

如果你在使用 ccam 之前就用过 Claude Code，你的配置都在 `~/.claude` 和 `~/.claude.json` 中。`ccam migrate` 会把它们复制到一个共享组中，让通过该组共享这些项的每个账号一开始就拥有你现有的配置。源文件夹只会被读取，直接运行 `claude` 也会像以前一样工作。在第一次 `ccam add` 之后运行一次，之后每当有账号开始共享 MCP 服务器时再运行一次。

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
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space（或 `x`）切换光标所在的项，方向键或 `j`/`k` 移动光标，Enter 显示计划，Esc、`q` 或 Ctrl+C 会打印 `Cancelled.`。源中没有内容的项会标记为 `not found`。如果你取消勾选所有项，ccam 会打印 `Nothing selected.` 并停止。

然后 ccam 打印计划，并在复制之前征求确认：

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

组只会到达它自己的账号，而且只涉及它共享的项，因此摘要会列出组中的账号。如果组不共享某个被复制的项，或者有些账号不在组中，ccam 会打印用来更改这些的 `ccam group set <group> --add <items>` 和 `ccam group join <group> <accounts>` 命令。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `[folder]` | 要从中复制的 Claude Code 文件夹 | `~/.claude` |
| `--settings` | 复制 `settings.json` | 关闭 |
| `--memory` | 复制全局 `CLAUDE.md` | 关闭 |
| `--extensions` | 复制 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 关闭 |
| `--mcp` | 复制用户 MCP 服务器（`.claude.json` 中的 `mcpServers`） | 关闭 |
| `--history` | 复制对话历史（`projects/`、`history.jsonl`） | 关闭 |
| `-g`、`--group string` | 要复制到的共享组 | `default` |
| `-y`、`--yes` | 不询问直接复制；发生冲突时保留组中的副本 | 关闭 |
| `--dry-run`（全局） | 打印计划，并在任何询问或复制之前停止 | 关闭 |

不带项目选项时，在终端中由清单决定。不在终端中或带有 `--yes` 时，ccam 复制除对话历史以外的所有内容。默认不包括对话历史，因为 `projects/` 可能很大。

`--group on` 表示 `default`。`--group off` 和 `--group local` 会被拒绝，并显示 `--group takes a group name, not "off"`。`ccam migrate` 支持 `--dry-run`，试运行时清单依然会出现，方便你预览所选内容。

## 工作原理

`internal/cli/migrate.go` 选定项目，然后 `migrate.Build` 读取源和组并返回一个计划。在你确认之后、执行 `Plan.Apply` 之前，不会写入任何内容。

1. **找到源。** `checkSource` 会拒绝不存在的路径、不是文件夹的路径，以及不包含任何已知 Claude Code 条目的文件夹。`migrate.ClaudeJSON` 找到属于该文件夹的 `.claude.json`：在文件夹内；或者对于内部没有该文件的 `~/.claude`，则是它旁边的 `~/.claude.json`。未设置 `CLAUDE_CONFIG_DIR` 时，Claude Code 就把它放在那里。
2. **把项目映射到路径。** 每个路径都来自 `share.Layout`，也就是组所使用的同一份允许列表。MCP 服务器通过 `mcp.LocalServers` 从 `.claude.json` 中读出，它使用 `profile.TopLevelValue` 扫描器，而不是解析整个文件。
3. **只复制，不移动。** 文件通过 `fsutil.AtomicWrite` 或 `fsutil.CopyFile` 写入 `~/.local/share/ccam/groups/<group>/`；MCP 服务器写入组的 `mcp.json`。

每一项都把一个源路径映射到一个组路径：

<!-- motion:migrate-items -->

两行历史记录是虚线，因为只有加上 `--history` 才会复制它们。`.credentials.json` 和 `.claude.json` 列在 `share.Never` 中，永远不会进入组：如果它们出现在布局中，`Build` 会返回内部错误，遍历文件夹时也会按名称跳过它们。

`buildEntry` 为每个路径确定一个状态，计划中的每一行描述的就是这个状态：

<!-- motion:migrate-states -->

- **Missing**：源中没有，或者应为文件夹的位置是一个文件。每项列出一次，显示为 `not found, skipped`。
- **Linked**：源路径已经解析到组中的副本，例如当源是 ccam 的账号文件夹时。
- **Same**：组中已有完全相同的内容。
- **New**：组中没有，因此整体复制。
- **Merge**：两边都有的文件夹。`copyMissing` 只添加组中没有的文件；内容不同的文件保持组中的版本，计划中显示 `group copies kept`。对于 `plugins/`，注册表 `installed_plugins.json` 和 `known_marketplaces.json` 会被合并，插件的绝对路径会被改写为指向组的文件夹。
- **Conflict**：两边都有但内容不同的文件。你可以逐个文件选择：`keep group`、`use mine`，对于 `.json` 文件还有 `merge JSON`（深度合并，值冲突时组中的值优先），对于 `history.jsonl` 还有 `merge by line`。定义不同的 MCP 服务器也属于冲突，可选 `keep group` 或 `use mine`。列出的第一个选项是按 Enter 时的默认值；`--yes` 总是保留组中的副本。

清单本身是 `tui.MigratePicker`，基于 `ccam add` 使用的 `checklist` 组件构建；`migrate.Found` 通过 lstat 调用和 MCP 服务器列表标记不存在的项。

### MCP 登录

通过 OAuth 登录的 MCP 服务器把令牌保存在登录文件 `.credentials.json` 的 `mcpOAuth` 下，而不是服务器列表中。Claude Code 用 `<name>|` 加上由服务器的 `type`、`url` 和 `headers` 计算出的 `sha256` 哈希的前 16 个字符作为每个条目的键。这个键与文件夹无关，因此只要服务器定义完全相同，复制过去的条目在其他账号中也能使用。所以 `internal/migrate/logins.go` 中的 `buildLogins` 只复制那些最终在组中采用源定义的服务器的条目：新增的、本来就相同的，或者在冲突中你选择了 `use mine` 的。当你复制扩展且源中有 `plugins/` 时，插件服务器的登录（以 `plugin:` 开头的键）也会一起复制。

<!-- motion:migrate-signins -->

这些条目写入每个接收账号自己的 `.credentials.json`，通过 `fsutil.AtomicWrite` 以 `0600` 权限写入。它们永远不会进入组。只有通过目标组共享 `mcp`（对于插件服务器则是 `extensions`）的账号才会接收它们，并且只添加账号中没有的键：账号已有的条目（可能更新）会被保留。计划中的 `sign-ins` 一行列出服务器、接收账号和被跳过的账号。出现以下情况时，ccam 会跳过该账号，并列出其名称和原因：

- 用 `ccam import` 原地注册的账号，因为直接运行的 `claude` 可能正在使用它的登录文件；
- 正在运行的账号，因为 Claude Code 会改写它的登录文件（写入前会再检查一次，报告为 `started meanwhile`）；
- 尚未登录的账号；
- 把登录保存在 macOS 钥匙串中的账号，ccam 不会写入钥匙串。

请在被跳过的账号中用 `/mcp` 登录这些服务器。之后才开始共享 MCP 服务器的账号，会在你再次运行 `ccam migrate` 时获得这些登录；摘要会用 ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` 提示这一点。

## 安全性与边界情况

- 源只会被读取。其中的任何文件都不会被移动、链接或修改，源的登录也从不会被写入，包括在可能从钥匙串读取登录的 macOS 上。
- 在组中的文件被替换之前（无论是因为 `use mine`、合并，还是现有 `mcp.json` 中加入了新的 MCP 服务器），组中的副本会保存在 `<group>/.ccam-backup/<time>/` 中。每次运行的文件夹名都是唯一的，因此同一秒内的两次运行不会覆盖彼此的备份。摘要会打印这个路径。
- 使用 `--yes` 时，每个冲突都保留组中的副本，ccam 会告诉你如何改为自己选择：不带 `--yes` 运行同一命令。
- 如果没有任何需要更改的内容，ccam 会打印 `Nothing to copy.`；如果每个冲突都保留组中的副本，则打印 `Nothing to copy: the group keeps its own copies.`
- 如果复制中途失败，ccam 会停止并显示 `The source folder was not changed. Group files replaced so far are saved in`，后面跟着备份文件夹。
- 源不存在时，会以 `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` 结束。

## 相关内容

- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [检查与修复](/features/doctor)
- [多个 Claude Code 账号](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
