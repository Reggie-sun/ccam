---
title: 检查与修复
description: 用 ccam doctor 检查共享链接、文件权限、登录和 Claude Code 版本，并用 ccam doctor --fix 修复可以安全修复的问题。
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` 检查你的 ccam 安装：Claude Code 可执行文件、每个账号的文件夹、登录和共享链接、共享的 MCP 服务器、用量接口以及你的快捷命令。`ccam doctor --fix` 修复它能在不丢失数据的前提下修复的问题，并先把它要替换的文件移到备份文件夹中。当 `ccam <name>` 因共享链接问题而停止、你手动恢复了文件之后，或者 Claude Code 更新之后，都可以运行它。

## 用法

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

每条检查结果占一行，并带有一个标记：`✓` 正常，`·` 信息，`!` 警告，`✗` 错误。与账号相关的结果以账号名开头。

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

使用 `--fix` 时，每条已修复的结果都以 `→ fixed` 结尾：

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--fix` | 修复可以安全修复的问题（更改的文件会被备份） | 关闭 |
| `--dry-run`（全局） | 执行所有检查但不做任何修复，即使带有 `--fix` | 关闭 |

`ccam doctor` 不接受参数。它支持 `--dry-run`：`ccam --dry-run doctor --fix` 打印与 `ccam doctor` 相同的结果，且不做任何更改。

只要至少有一条 `✗` 结果未被修复，退出码就是 `1`，否则为 `0`。警告和信息永远不会让运行失败。在不带 `--fix` 的情况下运行失败时，ccam 会以一个空行和 `` `ccam doctor --fix` repairs what it can. `` 结尾。无法加载配置时，也会以 `1` 和一条 `ccam:` 错误结束。

## 工作原理

`internal/cli/doctor.go` 加载 `~/.config/ccam/config.toml`，并调用 `internal/doctor/doctor.go` 中的 `doctor.Run`；后者按固定顺序执行检查并返回结果。`--fix` 作为 `Options.Fix` 传入，在试运行时设为 false。

<!-- motion:doctor-checks -->

中间那一组检查对每个账号各执行一次。各项检查以及 `--fix` 对它们的处理如下：

1. **Claude Code 版本**（`checkClaude`）。像启动器一样解析 `claude` 可执行文件，并运行 `claude --version`。如果版本不以 `2.1.302` 开头（ccam 的文件夹布局就是针对这个版本验证的，即 `doctor.LayoutVersion`），会加上 `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`。找不到可执行文件或 `--version` 失败属于错误。不修复。
2. **环境中的凭据**（`checkAuthEnv`）。当设置了 `CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN` 时发出警告：从该 shell 启动的每个账号都会继承它，而且 Claude Code 会优先使用它而不是账号自己的登录。不修复；请自行取消设置。
3. **文件夹**（`checkProfile`）。账号文件夹不存在属于错误，并结束该账号的检查。权限中设置了任何组或其他用户的位属于警告；`--fix` 会执行 `chmod 0700`。Windows 没有权限位，所以在 Windows 上会跳过这项检查以及下面的 `0600` 检查；用户配置文件中的文件夹由其访问控制列表保护。
4. **登录文件**（`checkCredentials`）。用 `lstat` 检查 `.credentials.json`。符号链接属于错误，因为每个账号必须有自己的令牌，且不会被修复。权限不是 `0600` 属于警告；`--fix` 会执行 `chmod 0600`。然后读取登录（在 macOS 上先从钥匙串读取）。没有登录属于警告，并提示使用 `ccam login <name>`。刷新令牌过期是错误 `login expired`；访问令牌过期只是信息，因为 Claude Code 会在下次启动时为它续期。这些都不会被修复。
5. **共享链接**（`share.Check`）。只使用 `lstat` 和 `readlink`，把 `share.Layout` 中的每个路径与账号的共享设置进行比较。它报告四种情况：缺少链接、链接指向别处、应为链接的位置是普通文件，以及某项已改为本地却仍残留链接。缺少链接属于警告，其他属于错误。
6. **未知条目**（`checkUnknownEntries`）。列出账号文件夹中 ccam 不认识（`share.Known`）的顶层名称，作为 Claude Code 布局可能已变化的提示。仅为信息。
7. **插件路径**（`checkPlugins`）。插件注册表（`installed_plugins.json`、`known_marketplaces.json`）保存的是绝对路径。在每个组的 `plugins/` 中，经过某个账号文件夹的路径会在该账号被删除或重命名时失效。`--fix` 用 `profile.RewritePluginPaths` 把它们改写为组文件夹路径，对注册表文件进行原子写入。
8. **MCP 服务器**（`checkMCP`）。每个组的 `mcp.json` 只加载一次。无法读取的文件属于错误；未通过 `mcp.Validate` 的服务器属于警告。不修复。
9. **用量响应格式**（`checkUsageFormat`）。不使用缓存，并采用通常的 3 秒时限，依次为各个账号查询用量接口，直到有一个响应为止。如果响应已无法解析，它会警告 `ccam usage` 将退回到状态栏缓存。这是唯一使用网络的检查。
10. **快捷命令**（`checkShims`）。仅在有账号设置了快捷命令时执行。`~/.local/bin` 中缺少链接属于警告；`--fix` 调用 `shim.Install`，它只创建缺少的链接，从不触碰其他文件。bin 文件夹不在 `PATH` 上属于警告，不会被修复。
11. **文件夹设置**（`checkFolders`）。仅在存在 [`ccam folder`](/features/folder) 条目时执行。条目对应的账号不存在属于错误，条目对应的文件夹已不存在属于警告。两者都不会被修复，因为删除条目会改变在那里启动的账号。

### 修复共享链接

启动路径会当场创建缺少的链接，但遇到任何可能丢失数据的情况都会停止，并提示你使用这个命令：

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

对于任何共享链接问题，`--fix` 都会以账号的共享设置同时作为当前状态和目标状态调用 `share.Build`，然后调用 `Plan.Apply`。`Build` 根据磁盘上的实际情况重新推导每一步，与 `ccam share` 使用的计划相同，因此修复遵循与开启共享相同的规则。

<!-- motion:doctor-repair -->

这里 `settings.json` 是一个普通文件。由于组中已有一个不同的 `settings.json`，冲突采用默认选择 `KeepGroup`，保留组中的副本。账号的文件被移到账号文件夹内的 `.ccam-backup/<time>/settings.json`，然后 `os.Symlink` 把该路径链接到组。其他情况：

- 如果组中还没有副本，账号的文件会被移入组中，成为组的第一份内容。
- 对于文件夹，组中没有的文件会被移入组中；其余的进入备份。
- 对于 `history.jsonl`，默认是 `merge by line`：先是组中的行，然后是你新增的行。组中原有的副本保存在 `.ccam-backup/<time>/group/` 中。
- 指向其他组的链接会被替换。指向 ccam 之外某个文件夹的链接会进入备份，如果组为空，它的内容会成为组的初始内容。
- 某项已改为本地却仍残留的链接会被移除，账号会得到组中内容的一份副本。

除了移除链接之外，每一步都是把数据移入组或备份文件夹。没有任何内容会被删除。

## 安全性与边界情况

- 单纯的 `ccam doctor` 从不更改任何内容。在脚本中，请把 `--dry-run` 与 `--fix` 一起使用以确保安全。
- `doctor --fix` 不检查账号是否正在运行。修复某个账号的共享链接之前，请先退出该账号的 Claude Code，因为会话会继续写入这些文件。
- 如果共享修复失败，ccam 会为该账号报告 `could not repair share links: <error>`，并跳过它的未知条目检查。其他账号仍会被检查。
- 如果 `ccam share` 中途失败，它会先保存目标设置，并提示你用 `ccam doctor --fix` 完成。
- 作为符号链接的 `.credentials.json`、过期的登录和环境中的凭据只会被报告，从不会被更改。
- 插件注册表的改写以原子方式写入，但不会复制到 `.ccam-backup/`。

## 相关内容

- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [复制现有配置](/features/migrate)
- [用量限制与切换账号](/docs/usage-limits)
- [命令](/docs/commands)
