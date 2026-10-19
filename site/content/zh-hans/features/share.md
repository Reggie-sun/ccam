---
title: 在账号之间共享
description: ccam share 如何把设置、CLAUDE.md、扩展和历史记录链接到组文件夹（附带计划、备份和冲突选择），以及 off 如何把内容复制回来。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

每个账号都有自己的文件夹，但其中的大部分内容并不需要分开。`ccam share` 每次为一个账号的一个项开启或关闭共享：账号的文件或文件夹被移入组文件夹，并替换为一个符号链接，这样组中的每个账号读写的都是同一份副本。在 `ccam add` 之后可以用它更改账号共享的内容，或者让几个账号单独组成一个组。

## 用法

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

不带共享项选项时，`ccam share` 打印当前状态：

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

带选项时，ccam 会打印计划，对每个与组中副本不同的文件进行询问，并在做出任何更改之前询问 `Proceed? [y/N]`。新账号共享除对话历史以外的所有内容；[共享 MCP 服务器、设置和历史记录](/docs/sharing)从用户的角度介绍了每一项。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP 服务器 | 不变 |
| `--settings <on\|off\|group>` | `settings.json` | 不变 |
| `--memory <on\|off\|group>` | 全局 `CLAUDE.md` | 不变 |
| `--extensions <on\|off\|group>` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 不变 |
| `--history <on\|off\|group>` | `projects/` 和 `history.jsonl` | 不变 |
| `--fresh` | 关闭共享时从空白开始，而不是复制组中的内容 | 关闭 |
| `--prefer <group\|mine\|merge>` | 如何决定冲突 | 询问；使用 `--yes` 时取每项的默认值 |
| `-y`、`--yes` | 不询问直接执行；冲突按 `--prefer` 处理 | 关闭 |
| `-h`、`--help` | `share` 的帮助 | |
| `--dry-run`（全局） | 支持：询问冲突，打印计划，不做任何更改 | 关闭 |

值 `on` 表示组 `default`，`off` 表示本地（不共享），其他任何值都是组名。组名会成为 `~/.local/share/ccam/groups/` 下的文件夹名，因此 `ccam share side --history ../x` 会失败并显示 `invalid group name: "../x"`。账号可以用名称或编号指定，只有你指定的项会改变。

## 工作原理

涉及的文件：

- `~/.config/ccam/config.toml`：每个账号的 `[profile.share]` 表，每项一个值。
- `~/.local/share/ccam/profiles/<name>/`：账号的 `CLAUDE_CONFIG_DIR`。
- `~/.local/share/ccam/groups/<group>/`：共享的副本。

1. `internal/cli/share.go` 读取配置，用 `config.Share.Set` 把每个 `--<item>` 的值应用到账号设置的副本上；如果账号正在运行，则拒绝执行。
2. `share.Build` 比较当前设置和目标设置，并列出步骤。它只读取磁盘：对 `share.Layout` 中的每个路径调用 `os.Lstat` 和 `os.Readlink`。计划是一系列操作：移入组中空着的位置（adopt）、移入组中没有的文件（merge dir）、把已解决的冲突写入组、备份、取消链接、链接、复制回来、改写插件路径以及合并插件注册表。
3. 通过提示或 `--prefer` 决定冲突。
4. `Plan.Describe` 打印计划，ccam 征求确认。
5. 在移动任何内容之前，目标设置会先写入 `config.toml`。如果 `Plan.Apply` 中途停止，`ccam doctor --fix` 或重新运行会朝着已保存的目标构建计划，完成这次更改。
6. `Plan.Apply` 按顺序执行各步骤。每一步都是把数据移入组或备份文件夹；没有任何内容会被删除。

`share.Layout` 是一份允许列表，列出每一项涵盖的路径。不在列表中的内容永远不会进入组。`share.Never` 列出了 `.credentials.json` 和 `.claude.json`，如果其中任何一个出现在布局中，`Build` 会以内部错误停止。`mcp` 这一项没有路径：服务器列表位于 `.claude.json` 内部，所以 [MCP 共享](/features/mcp)改为通过组的 `mcp.json` 实现，`--mcp` 只改变设置。

下面是当 `side` 开启历史记录共享、而组中已有内容时，`projects/` 这个文件夹会发生的变化：

<!-- motion:share-on -->

如果组中该路径下还没有内容，这个账号就是第一个共享它的账号，它的文件夹会被原样移过去，不做备份。否则，`mergeDir` 会把组中没有的每个文件移入组中，两边都有的文件留在账号的文件夹中，从而保留组中的版本。然后账号的文件夹被移到账号文件夹内的 `.ccam-backup/<time>/`，`os.Symlink` 在原位置放上一个链接。对话文件以会话 ID 命名，因此多个账号的历史记录可以合并而不会冲突。

移动使用 `rename(2)`，跨文件系统时退回为复制后删除。被移动的目录树中指向树外的相对链接会被改写为绝对链接（`fsutil.PortableLink`），这样从新位置仍能解析。`plugins/` 中的注册表（`installed_plugins.json`、`known_marketplaces.json`）保存的是经过账号文件夹的绝对路径；ccam 会把它们指向组的文件夹，如果组中已有注册表，则把账号的条目合并进去。

关闭某一项会反转链接：

<!-- motion:share-off -->

只有当链接指向组文件夹内部时，ccam 才会移除它，然后用 `CopyTree` 把组中当前的内容复制到账号中。使用 `--fresh` 时会跳过复制，账号在该路径下从空白开始。组会为其他账号保留其内容。把某一项从一个组移到另一个组时，会先取消链接再链接到新组，旧组保持不变；计划会在一行 `Note:` 中说明这一点。

两边都存在但内容不同的单个文件属于冲突：

<!-- motion:share-conflict -->

`share.NewConflict` 列出适用于该文件的选项。`history.jsonl` 提供 `merge by line`、`keep group` 和 `use mine`，合并排在第一位。其他 `.json` 文件提供 `keep group`、`use mine` 和 `merge JSON`。其他文件提供 `keep group` 和 `use mine`。第一个选项是默认值：

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON`（`MergeJSONBytes`）递归合并对象，合并数组并去除重复项，当两边对同一个键设置了不同的值时保留组中的值，因此其他账号的行为不会改变。`merge by line` 保留组中的行，再加上组中没有的你的行。除 `keep group` 之外的每个选项都会先把组中的副本保存到账号文件夹内的 `.ccam-backup/<time>/group/`，然后用 `fsutil.AtomicWrite` 写入结果。你自己的文件总会进入备份。

`--prefer group`、`mine` 或 `merge` 会不加询问地决定所有冲突；`merge` 表示该文件提供的 `merge JSON` 或 `merge by line`，两者都不提供的文件则采用默认值。单独使用 `--yes` 时，每个文件都采用其默认值。

每次启动时，`share.Check` 只用 `lstat` 和 `readlink` 把布局中的每个路径与设置进行比较，因此开销足够小，可以放在启动路径上。缺少的链接会当场创建，因为不会丢失任何内容。应为链接的位置是普通文件、链接指向错误的位置，或者残留的过时链接，都会中止启动：

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` 以已保存的设置同时作为两边运行 `share.Build`，从而根据磁盘上的实际情况推导每一步，冲突时保留组中的内容，并备份它所替换的任何内容。

## 安全性与边界情况

- **正在运行的账号。** 当 Claude Code 正以该账号运行时，ccam 拒绝更改共享（`account work is running; quit it first, then change its sharing`），因为会话会继续写入正在移动的文件，期间追加的行会丢失。`--dry-run` 仍会打印计划。
- **备份。** 每个计划在账号文件夹中都有自己的 `.ccam-backup/<time>/` 文件夹，两个计划在同一秒开始时会加上编号后缀。ccam 从不删除它们。
- **部分失败。** 如果某一步失败，ccam 会打印错误以及原始文件所在的位置，并提示你用 `ccam doctor --fix` 完成。目标已经保存，因此修复会朝着它继续进行。
- **你自己的链接。** 如果某个路径已经是你自己创建的链接（例如指向 dotfiles 仓库），且组为空，ccam 会把链接的目标复制到组中，并把链接本身移到备份中。目标不会被改动。
- **登录文件。** `.credentials.json` 和 `.claude.json` 永远不会被链接、移动或复制到组中。
- **同一文件夹中的两个账号。** 共享历史记录的账号会写入相同的文件。当同一历史组的第二个账号在已有账号运行的文件夹中启动时，ccam 会发出警告，并在继续之前征求确认。

## 相关内容

- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [共享的 MCP 服务器](/features/mcp)
- [复制现有配置](/features/migrate)
- [所有命令](/docs/commands)
