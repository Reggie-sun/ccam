---
title: 通过组共享
description: ccam 的组通过一个文件夹在其账号之间共享选定的项；ccam group 按计划移动数据，并提供备份和冲突选择。
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "e4f05b0fd572"
---

每个账号都有自己的文件夹，但其中的大部分内容并不需要分开。**组**规定共享什么、由谁共享：组中的每个账号都把这些项链接到同一个文件夹 `~/.local/share/ccam/groups/<group>/`，因此它们读写的都是同一份副本。登录令牌永远不会被共享。

一个账号可以属于多个组，但每一项只能来自一个组。典型的设置是：一个组用于设置和扩展，另一个组用于在两个账号之间共享对话历史：

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

你添加的第一个账号会创建 `default`，之后添加的账号都会加入它，除非你另行指定。

## 用法

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

每个会移动数据的命令都会打印计划，对每个与组中副本不同的文件进行询问，并在做出任何更改之前询问 `Proceed? [y/N]`：

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

同样的更改也可以在[仪表盘](/features/dashboard)的 **Groups** 选项卡中完成，在某个账号上按 `s` 可以选择它所属的组。

## 项

| 项 | 涵盖的内容 |
| --- | --- |
| `mcp` | MCP 服务器，通过组的 `mcp.json` 实现（[共享的 MCP 服务器](/features/mcp)） |
| `settings` | `settings.json` |
| `memory` | 全局 `CLAUDE.md` |
| `extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| `history` | `projects/` 和 `history.jsonl`，这样另一个账号就可以 `--continue` |

账号的任何组都不共享的项，会留在账号自己的文件夹中。

## 命令和选项

| 命令 | 作用 |
| --- | --- |
| `ccam group` 或 `ccam group list` | 列出组 |
| `ccam group show <group>` | 组共享的内容、它的账号、文件夹和 MCP 服务器 |
| `ccam group add <group> --items <items> [--members <accounts>]` | 创建组；其账号的这些项会移入组中 |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | 为组中的每个账号更改组共享的内容 |
| `ccam group join <group> <account>...` | 把账号加入组 |
| `ccam group leave <group> <account>...` | 把账号移出组；每个账号保留一份副本，包括 MCP 服务器 |
| `ccam group rename <group> <new-name>` | 只在配置中重命名 |
| `ccam group remove <group> [--leave]` | 删除组 |

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--from <account>` | `add`、`set`：其副本最先进入组的账号 | 默认账号，然后按编号 |
| `--replace` | `join`：先离开共享相同项的组 | 关闭 |
| `--fresh` | 停止共享某一项的账号从空白开始，而不是带着一份副本（导入的文件夹仍会保留一份） | 关闭 |
| `--leave` | `remove`：先把组中的账号移出 | 组中有账号时必需 |
| `--prefer <group\|mine\|merge>` | 如何决定冲突 | 询问；使用 `--yes` 时取每个文件的默认值 |
| `-y`、`--yes` | 不询问直接执行；冲突按 `--prefer` 处理 | 关闭 |
| `--dry-run`（全局） | 打印计划，不做任何更改 | 关闭 |

账号可以用名称或编号指定；项用逗号分隔的列表指定。`ccam add` 接受 `--group <group>`（可重复）或 `--no-group`。

## 工作原理

涉及的文件：

- `~/.config/ccam/config.toml`：`[[group]]` 条目，包括 `name`、`items`、`members`，以及在文件夹名与组名不同时的 `dir`。
- `~/.local/share/ccam/profiles/<name>/`：账号的 `CLAUDE_CONFIG_DIR`。
- `~/.local/share/ccam/groups/<folder>/`：共享的副本。

读取 `config.toml` 时，每个账号的组会被解析为每项一个文件夹。五个账号大约需要 25 µs，因此启动时仍然只读取 `config.toml`。如果某个配置会让一个账号从两个组获得同一项，就属于错误，在修正之前什么都不会启动。

更改是在配置的副本上进行的：

1. 命令编辑这份副本（`AddGroup`、`SetGroupItems`、`JoinGroup` …），同时检查规则。
2. `internal/group` 比较新旧配置，并为每个项发生变化的账号用 `share.Build` 构建计划。这些计划依次构建，每个计划都能看到前面的计划将放入组文件夹的内容，因此只有第一个账号会把它的副本移进去。第一份副本的来源依次是：组文件夹中已有的副本、`--from`、默认账号，然后是按编号排列的其他账号。
3. 如果其中某个账号正在运行，无论 Claude Code 是怎样启动的，ccam 都会拒绝执行（``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``）；在没有开启开发人员模式的 Windows 上，如果某个文件需要链接，也会拒绝。
4. 通过提示或 `--prefer` 决定冲突，打印计划，然后 ccam 征求确认。
5. 就在保存之前，ccam 会再检查一次：在等待回答期间启动的账号会让更改停止，另一个 ccam 对同一些组的更改也会如此（`the groups changed while this plan waited; nothing was moved, make the change again`）。在移动任何内容之前，新配置会先保存，并且这次更改会记录在 `~/.local/state/ccam/group-change.json` 中。如果某一步中途停止，这条记录会保留：`ccam doctor --fix` 会按原来的顺序、使用你给出的回答完成更改，包括 MCP 相关的步骤；在此之前不会开始其他组更改（``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``）。
6. MCP 相关的步骤最先执行。当账号把扩展移出某个组时，只要组的插件列表中的路径经过某个账号的文件夹，就会被改为指向组文件夹（计划中的 `group plugin lists:`）。然后各计划按顺序执行。每一步都是把数据移入组或备份文件夹；只有组中已有的逐字节相同的文件会被丢弃，而不是备份。

`share.Layout` 是一份允许列表，列出每一项涵盖的路径。不在列表中的内容永远不会进入组。`share.Never` 列出了 `.credentials.json` 和 `.claude.json`，如果其中任何一个出现在布局中，`Build` 会以内部错误停止。`mcp` 这一项没有路径：服务器列表位于 `.claude.json` 内部，所以 [MCP 共享](/features/mcp)通过组的 `mcp.json` 实现。

下面是当 `side` 开始共享历史记录、而组中已有内容时，`projects/` 这个文件夹会发生的变化：

<!-- motion:share-on -->

如果组中该路径下还没有内容，这个账号就是第一个共享它的账号，它的文件夹会被原样移过去。否则，组中没有的每个文件都会被移入组中，两边都有的文件留在账号的文件夹中，从而保留组中的版本。与组中相同的文件会被丢弃，计划会列出内容不同的文件（`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`）。然后账号的文件夹被移到账号文件夹内的 `.ccam-backup/<time>/`，原位置换成一个链接。如果组中的条目本身就是指向加入账号自己条目的链接，账号的内容会进入组中取代它。对话文件以会话 ID 命名，因此多个账号的历史记录可以合并而不会冲突。`plugins/` 中的注册表保存的是经过账号文件夹的绝对路径；ccam 会把它们指向组的文件夹，如果组中已有注册表，则把账号的条目合并进去。

离开会反转链接：

<!-- motion:share-off -->

ccam 先把组中当前的内容复制到账号文件夹内的 `.ccam-backup/incomplete-<name>`，复制完成后才替换链接，因此中途停止的复制会让链接保持原样，`ccam doctor --fix` 会重新复制（这种收尾的计划会显示 `restore <entry>: … (copied back from the group before a stop)`）。只有当链接仍像制定计划时那样指向组文件夹内部时，ccam 才会移除它；此后发生变化的路径，例如原来链接的位置变成了普通文件，会保持原样，更改也会停止。使用 `--fresh` 时会跳过复制。组会为其他账号保留其内容。如果一个账号把某一项移到文件夹仍为空的组，它会带上自己原先看到的内容的副本，因此文件得以保留；这份副本在完整之前暂存在 `<group>/.ccam-incomplete-<name>` 中。停止共享 MCP 的账号会保留组的服务器：它们会写入该账号自己的 `.claude.json`（会保留备份），与它自己的同名服务器放在一起，后者保持不变。使用 `--fresh` 时，计划会改为列出将不再加载的服务器。再次加入时，没人改过的副本会被移除。离开 `extensions` 时，组的插件列表中经过离开账号文件夹的路径也会被改为指向组文件夹；对其他账号来说，它们的含义不变。

两边都存在但内容不同的单个文件属于冲突：

<!-- motion:share-conflict -->

对于 `settings.json` 这样的单个文件，由你来选择；对于文件夹内部，保留组的副本，计划会列出内容不同的文件。`history.jsonl` 提供 `merge by line`、`keep group` 和 `use mine`，合并排在第一位。其他 `.json` 文件提供 `keep group`、`use mine` 和 `merge JSON`。其他文件提供 `keep group` 和 `use mine`。`merge JSON` 递归合并对象，合并数组并去除重复项，当两边对同一个键设置了不同的值时保留组中的值，因此其他账号看不到任何变化。除 `keep group` 之外的每个选项都会先把组中的副本保存到账号的 `.ccam-backup/<time>/group/`。你自己的文件总会进入备份。

每次启动时，`share.Check` 只用 `lstat` 和 `readlink` 把布局中的每个路径与账号的组进行比较。缺少的链接会当场创建。应为链接的位置是普通文件、链接指向错误的位置，或者残留的过时链接，都会中止启动，并提示使用 `ccam doctor --fix`。

### 重命名、删除和旧格式

- `ccam group rename` 只更改配置。文件夹保留原名，组把它记录为 `dir`，因此不会移动任何内容，也不会丢失任何内容。
- 对于有账号的组，`ccam group remove` 会拒绝执行，除非加上 `--leave`。这样每个账号都会保留一份副本，文件夹会移到 `~/.local/share/ccam/trash/`，除非另一个组还在使用它。
- 在引入组之前，`config.toml` 为每个账号保存一个 `[profile.share]` 表。ccam 会把这种格式读作组，不移动任何内容。第一次保存时，它会写入新格式，并把旧文件保留为 `config.toml.bak`。由相同账号共享的项会合并为一个以文件夹命名的组；其余的项会得到 `default-history` 这样的名称，并继续使用同一个文件夹。

### 在 Windows 上

Windows 只在开发人员模式（设置 > 系统 > 开发者选项）下或为管理员创建符号链接。没有这项权限时，共享的文件夹会变成目录联接，对 Claude Code 来说效果相同。共享的文件（`settings.json`、`CLAUDE.md`、`history.jsonl`）只能是符号链接，所以需要符号链接的组更改会在移动任何内容之前被拒绝。出于同样的原因，在那里第一次运行 `ccam add` 时创建的 `default` 组只包含 `mcp` 和 `extensions`，并会说明这一点；用 `ccam add --group` 加入共享设置、memory 或历史记录的组会被拒绝。

## 安全性与边界情况

- **正在运行的账号。** 当 Claude Code 正以某个要移动的账号运行时，更改会被拒绝，因为会话会继续写入这些文件，无论它是不是由 ccam 启动的（直接运行的 `claude`、IDE、后台会话）。就在保存之前会再检查一次。离开和加入只涉及所指定的账号；离开 `extensions` 时还会改写组的插件列表中的路径，但不会改变它们对其他账号的含义。`--dry-run` 仍会打印计划。
- **备份。** 每个计划在账号文件夹中都有自己的 `.ccam-backup/<time>/` 文件夹。ccam 从不删除它们。
- **你自己的链接。** 如果某个路径已经是你自己创建的链接（例如指向 dotfiles 仓库），且组为空，ccam 会把链接的目标复制到组中，并把链接本身移到备份中。目标不会被改动。
- **导入的文件夹。** 用 `ccam import` 在原位置注册的账号（例如 `~/.claude`）会在原位置被更改，因此直接运行 `claude` 也会看到变化；计划中会注明这一点。它离开组时总会保留一份副本，即使使用 `--fresh` 也是如此；`ccam remove` 在忘记该账号之前，也会以这种方式让它离开各个组。它在组中时，共享的文件存放在组文件夹中，所以在删除 `~/.local/share/ccam` 之前，请先删除该账号（或离开这些组）。
- **同一文件夹中的两个账号。** 共享历史记录的账号会写入相同的文件。当第二个账号在已有账号运行的文件夹中启动时，ccam 会发出警告并征求确认。
- **手动编辑。** 不是账号的成员会被忽略，并由 `ccam doctor` 报告。`local`、`on`、`off` 和 `none` 不能用作组名。

## 相关内容

- [仪表盘](/features/dashboard)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [共享的 MCP 服务器](/features/mcp)
- [复制现有配置](/features/migrate)
- [所有命令](/docs/commands)
