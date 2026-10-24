---
title: 检查与修复
description: 用 ccam doctor 检查共享链接、文件权限、登录和 Claude Code 版本，并用 ccam doctor --fix 修复可以安全修复的问题。
updated: 2026-10-23
order: 9
translation: machine
sourceHash: "11ad34630dcf"
---

`ccam doctor` 检查你的 ccam 安装：Claude Code 可执行文件、每个账号的文件夹、登录和共享链接、共享的 MCP 服务器、用量接口以及你的快捷命令。`ccam doctor --fix` 修复它能在不丢失数据的前提下修复的问题，并先把它要替换的文件移到备份文件夹中。当 `ccam <name>` 因共享链接问题而停止、你手动恢复了文件之后，或者 Claude Code 更新之后，都可以运行它。

## 用法

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

第一行一次性标出所有账号，让你在阅读其他内容之前就知道问题在哪里。下面每个账号有一个框，不属于任何账号的结果（Claude Code、环境变量、插件、MCP 服务器、快捷命令、文件夹规则、组和设置）放在 `General` 框中。框的标题显示该账号最严重的标记和问题数量。框内每一行是一个检查过的方面，总是依次为 `Folder`、`Login` 和 `Share links`，然后是其他方面，并带有一个标记：`✔` 正常，`ℹ` 信息，`⚠` 警告，`✘` 错误。没有任何问题的文件夹显示 `✔ Folder  OK`。

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

使用 `--fix` 时，每条已修复的结果都以 `→ fixed` 结尾，并且不再算作问题：

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

在终端中，框会适应终端宽度，过长的消息会换行，并使用终端自己的颜色（`NO_COLOR` 可以关闭颜色）。框的宽度永远不小于 29 列，因此在较窄的终端上也不会有单词被截断。输出到管道或文件时，没有颜色，任何消息都不会被换行或截断，因此 `ccam doctor | grep settings.json` 能找到完整的一行；此时每个框都与最长的一行一样宽。仪表盘的 Health 选项卡显示同样的框。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--fix` | 修复可以安全修复的问题（更改的文件会被备份） | 关闭 |
| `--dry-run`（全局） | 执行所有检查但不做任何修复，即使带有 `--fix` | 关闭 |

`ccam doctor` 不接受参数。它支持 `--dry-run`：`ccam --dry-run doctor --fix` 打印 `ccam doctor` 的结果，加上 `--fix` 会放到一边的内容，且不做任何更改：

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

只要至少有一条 `✘` 结果未被修复，退出码就是 `1`，否则为 `0`。警告和信息永远不会让运行失败。在不带 `--fix` 的情况下运行失败时，ccam 会以一个空行和 `` `ccam doctor --fix` repairs what it can. `` 结尾。无法加载配置时，也会以 `1` 和一条 `ccam:` 错误结束。

## 工作原理

`internal/cli/doctor.go` 加载 `~/.config/ccam/config.toml`，并调用 `internal/doctor/doctor.go` 中的 `doctor.Run`；后者按固定顺序执行检查并返回结果，每条结果都带有它所属的账号和检查的方面。`tui.HealthReport` 把它们归入各个框。`--fix` 作为 `Options.Fix` 传入，在试运行时设为 false。

<!-- motion:doctor-checks -->

中间那一组检查对每个账号各执行一次。各项检查以及 `--fix` 对它们的处理如下：

1. **Claude Code 版本**（`checkClaude`）。像启动器一样解析 `claude` 可执行文件，并运行 `claude --version`。找不到可执行文件、`claude_bin` 不是可执行文件，或者 `--version` 失败，都属于错误。不修复。
2. **环境中的凭据**（`checkAuthEnv`）。当设置了 `CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN` 时发出警告：从该 shell 启动的每个账号都会继承它，而且 Claude Code 会优先使用它而不是账号自己的登录。不修复；请自行取消设置。
3. **文件夹**（`checkProfile`）。账号文件夹不存在属于错误，并结束该账号的检查。权限中设置了任何组或其他用户的位属于警告；`--fix` 会执行 `chmod 0700`。导入的 `~/.claude` 在 `~/.claude.json` 之外还有自己的 `.claude.json` 时属于信息：此时 `ccam main` 和直接运行的 `claude` 各自保存登录信息和 MCP 服务器。Windows 没有权限位，所以在 Windows 上会跳过这项检查以及下面的 `0600` 检查；用户配置文件中的文件夹由其访问控制列表保护。
4. **登录文件**（`checkCredentials`）。用 `lstat` 检查 `.credentials.json`。符号链接属于错误，因为每个账号必须有自己的令牌，且不会被修复。权限不是 `0600` 属于警告；`--fix` 会执行 `chmod 0600`。然后读取登录（在 macOS 上先从钥匙串读取）。没有登录，或者登录文件中没有登录，属于警告，并提示使用 `ccam login <name>`。刷新令牌过期，或者服务器拒绝后 Claude Code 清除了登录，属于错误 `` login expired → `ccam login <name>` ``；访问令牌过期只是信息，因为 Claude Code 会在下次启动时为它续期。保存在 Windows 凭据管理器中的登录属于信息，因为 ccam 无法读取它。有效的登录会显示该文件夹当前登录的邮箱和套餐（`login valid (dev@company.com, Claude Team · Acme)`）。这些都不会被修复；`ccam login --expired` 会登录所有需要登录的账号。
5. **共享链接**（`share.Check`）。只使用 `lstat` 和 `readlink`，把 `share.Layout` 中的每个路径与账号的组进行比较。它报告五种情况：缺少链接、链接指向别处、应为链接的位置是普通文件、某项已改为本地却仍残留链接，以及链接正确但其组文件夹不存在（`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`；启动时不检查这一点）。缺少链接和组文件夹不存在属于警告，其他属于错误。从组复制回来时中途停止的情况也会被报告：`--fix` 会把暂存的完整副本放到位，但如果该路径在此之后（因启动而）被重新创建，它会指出存放完整副本的文件夹，由你自行处理。在较新版本的 ccam 中加入某一项的路径（历史记录的 `file-history/`、`tasks/` 和 `paste-cache/`，设置的 `keybindings.json`），在账号仍有自己的副本时只算警告：`is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`。启动绝不会因为它们而中止。
6. **未知条目**（`checkUnknownEntries`）。列出账号文件夹中 ccam 不认识（`share.Known`）的顶层名称，作为 Claude Code 布局可能已变化的提示。仅为信息。
7. **插件路径**（`checkPlugins`）。插件注册表（`installed_plugins.json`、`known_marketplaces.json`）保存的是绝对路径。在每个组的 `plugins/` 中，经过某个账号文件夹的路径会在该账号被删除或重命名时失效。`--fix` 用 `profile.RewritePluginPaths` 把它们改写为组文件夹路径，对注册表文件进行原子写入。
8. **MCP 服务器**（`checkMCP`）。每个组的 `mcp.json` 只加载一次。无法读取的文件属于错误；未通过 `mcp.Validate` 的服务器属于警告。直接运行的 `claude` 所用 `~/.claude.json` 中没有任何账号加载的服务器属于警告，并给出复制它们的 `ccam migrate --mcp` 命令（当有多个组共享 MCP 时带上 `--group <a|b>`）。组中只有一个账号拥有的服务器属于信息，并给出共享它们的 `ccam mcp import` 命令；如果某个账号未登录其他账号已登录的 MCP 服务器，会收到一条警告，提示用 `/mcp` 登录（ccam 从不复制登录）。不修复。
9. **用量响应格式**（`checkUsageFormat`）。不使用缓存，依次为各个账号查询用量接口，直到有一个响应为止，全部在一个 4 秒时限内完成。设置 `usage_source = "claude"` 时，改为向 Claude Code 查询一个账号。如果响应已无法解析，它会警告 `ccam usage` 将退回到状态栏缓存；没有任何账号响应时，会作为信息说明这一点（`could not check the usage response format: …`）。这是唯一使用网络的检查。
10. **快捷命令**（`checkShims`）。仅在有账号设置了快捷命令时执行。`~/.local/bin` 中缺少链接属于警告，ccam 创建的快捷命令仍运行着已被替换或移动的旧版 ccam 时也属于警告（`shortcuts in <bin> that run another ccam (one replaced or moved): …`）；`--fix` 调用 `shim.Install`，它只创建和更新 ccam 创建并记录在 `~/.local/state/ccam/shortcuts.json` 中的链接，从不触碰其他文件。只有在这些链接都已创建时，警告才以 `→ fixed` 结尾；如果某个链接无法创建，会附上原因。运行另一个仍已安装的 ccam 的快捷命令可以正常使用，属于信息；`--fix` 会把它指向当前这个 ccam。bin 文件夹不在 `PATH` 上属于警告，不会被修复。
11. **文件夹设置**（`checkFolders`）。仅在存在 [`ccam folder`](/features/folder) 条目时执行。条目对应的账号不存在属于错误，条目对应的文件夹已不存在属于警告。两者都不会被修复，因为删除条目会改变在那里启动的账号。
12. **组和设置**（`checkGroups`）。不是账号的组成员、名称后来变成了 ccam 命令的账号（用 `ccam run <name>` 启动它或给它改名）、未知的 `theme` 或 `usage_source`，以及 ccam 不读取的 `config.toml` 行，都属于警告。最后这种情况通常是写在 `[[group]]`、`[[profile]]` 或 `[[folder]]` 表下面的设置，它会被当作该表的一部分（`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`）；在修正之前，更改配置的命令会拒绝执行。不做修复。
13. **登录**（`checkIdentities`）。从每个文件夹的 `.claude.json` 读取它以谁的身份登录。如果文件夹现在登录的账号与 `ccam login` 记录的不同（例如在会话中执行了 `/login` 之后），属于警告；没有登录记录的账号（例如导入的账号或旧版 ccam 创建的账号）会按记录的邮箱进行比较。如果文件夹与另一个文件夹登录的是同一个 Claude 账号，也属于警告：两者会共用用量限制。如果文件夹与另一个文件夹持有同一个登录令牌，即复制的文件夹（包括直接运行的 `claude` 所用的 `~/.claude`），属于错误，因为先续期的一方会让另一方退出登录；运行 `ccam login <name>` 让它拥有自己的登录。不做修复。

在检查各个账号之前，中途停止的**组更改**会在 `General` 框中作为错误报告：``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``。`--fix` 会按原计划的顺序、使用对其冲突给出的回答完成它；在此之前，不会开始其他组更改、重命名或删除。在它完成之前，其中的账号只会补上缺少的链接。如果无法完成，错误信息会说明删除那个文件就会放弃更改的剩余部分。`ccam --dry-run doctor --fix` 会显示顺序（`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`）以及每个账号的那部分会放到一边的内容。另一个 ccam 正在执行的更改属于信息，在它运行期间，`--fix` 不修复任何共享链接（`its share links are not repaired while another ccam changes the groups; …`）。`--fix` 在修复期间持有 `~/.local/state/ccam/group-change.lock`，并在完成一次更改后重新读取 `config.toml`。

### 修复共享链接

启动路径会当场创建缺少的链接，但遇到任何可能丢失数据的情况都会停止，并提示你使用这个命令：

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

对于任何共享链接问题，`--fix` 都会以账号的共享设置同时作为当前状态和目标状态调用 `share.Build`，然后调用 `Plan.Apply`。`Build` 根据磁盘上的实际情况重新推导每一步，与 `ccam group` 更改使用的计划相同，因此修复遵循与加入组相同的规则。

<!-- motion:doctor-repair -->

这里 `settings.json` 是一个普通文件。由于组中已有一个不同的 `settings.json`，冲突采用默认选择 `KeepGroup`，保留组中的副本。账号的文件被移到账号文件夹内的 `.ccam-backup/<time>/settings.json`，然后 `os.Symlink` 把该路径链接到组。`--fix` 会说明这一点：`settings.json differed from the group's copy, which it now uses; this account's copy is in …`。其他情况：

- 如果组中还没有副本，账号的文件会被移入组中，成为组的第一份内容。
- 对于文件夹，组中没有的文件会被移入组中。与组中逐字节相同的文件会被丢弃，内容不同的文件进入备份，并在一条说明中列出（`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`）。
- 缺少的组文件夹会先被创建。
- 对于 `history.jsonl`，默认是 `merge by line`：先是组中的行，然后是你新增的行。组中原有的副本保存在 `.ccam-backup/<time>/group/` 中。
- 指向其他组的链接会被替换。指向 ccam 之外某个文件夹的链接会进入备份，如果组为空，它的内容会成为组的初始内容。
- 某项已改为本地却仍残留的链接会被移除，账号会得到组中内容的一份副本。

除了移除链接以及丢弃组中已有的逐字节相同的文件之外，每一步都是把数据移入组或备份文件夹。不会丢失任何内容。

## 安全性与边界情况

- 单纯的 `ccam doctor` 从不更改任何内容。在脚本中，请把 `--dry-run` 与 `--fix` 一起使用以确保安全。
- 对于正在运行的账号，无论 Claude Code 是怎样启动的，`doctor --fix` 只创建缺失的链接和缺失的组文件夹，这不会丢失任何内容。如果还有其他问题未处理，它会加上 `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``，因为会话会继续写入修复时要移动的文件。
- 如果共享修复失败，ccam 会为该账号报告 `could not repair share links: <error>`，并跳过它的未知条目检查。其他账号仍会被检查。
- 如果 `ccam group` 的更改中途失败，新的组已经保存，该更改会记录在 `~/.local/state/ccam/group-change.json` 中；ccam 会提示你用 `ccam doctor --fix` 完成，它会执行剩余的部分，包括 MCP 相关的步骤。
- 作为符号链接的 `.credentials.json`、过期的登录和环境中的凭据只会被报告，从不会被更改。
- 插件注册表的改写以原子方式写入，但不会复制到 `.ccam-backup/`。

## 相关内容

- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [复制现有配置](/features/migrate)
- [用量限制与切换账号](/docs/usage-limits)
- [命令](/docs/commands)
