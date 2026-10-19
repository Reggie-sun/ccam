---
title: 切换到下一个账号
description: ccam next 以你上次使用的账号之后的账号打开 Claude Code；--auto 选择 5 小时用量最低的账号，--continue 继续当前对话。
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` 无需你指定账号名，就能以另一个账号打开 Claude Code。当某个账号达到用量限制时使用它：单纯的 `ccam next` 选择你上次使用的账号之后的那个账号，`--auto` 选择 5 小时时间窗口内余量最多的账号，`--continue` 则在新账号中继续上一次对话。

## 用法

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam 在标准错误上以 `→ <from> → <to>` 的形式打印这次切换，然后像 `ccam <name>` 一样以新账号启动 Claude Code。`<from>` 是你上次使用的账号。`--` 之后的参数原样传给 `claude`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--auto` | 选择已登录账号中 5 小时用量最低的账号，跳过用量已达 100% 的账号 | 关闭 |
| `-c`、`--continue` | 只考虑与你上次使用的账号共享对话历史的账号，并启动 `claude --continue` | 关闭 |
| `-- <claude args>` | 放在 ccam 自身参数之后传给 `claude` 的参数 | 无 |
| `--dry-run` | 支持：打印 `→ <from> → <to>` 和将要执行的命令行，不启动 Claude Code，也不记录启动 | 关闭 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 时，`--continue` 会放在你自己的参数之前。

## 工作原理

每次通过 ccam 启动都会把时间写入 `~/.local/state/ccam/last-used/<name>`。`ccam next` 读取这些文件，检查哪些账号已登录，然后做出选择。

1. `cli.candidates` 根据 `config.toml` 为每个账号构建一个 `next.Candidate`：账号本身、它最近一次启动的时间（`launcher.LastUsed`），以及它是否有登录（`creds.Has`：可读的 `.credentials.json`，或在 macOS 上的钥匙串项）。
2. 不带 `--auto` 时，`next.Rotate` 按槽位号对账号排序，并找出最近一次启动的账号。这就是你要切换离开的账号。如果还没有任何账号启动过，则取槽位号最小的账号。
3. 从那里开始，它按槽位顺序向后查找，到最后一个槽位后回到第一个，返回第一个已登录的账号。
4. `launcher.LaunchWith` 像 `ccam <name>` 一样准备并启动所选账号：检查账号的链接，应用 MCP 共享，注册会话，并通过 `syscall.Exec` 用 `claude` 替换 ccam。

查找从上次使用的账号开始，在它之后的第一个已登录账号处停止：

<!-- motion:next-rotate -->

`work` 是最近一次启动的，因此从它之后开始查找。`side` 没有登录，被跳过，于是选中 `lab`。如果 `lab` 也未登录，查找会绕回到 `personal`（图中的 `% len`）。轮换永远不会返回你出发时的账号：如果没有其他已登录的账号，它会以错误停止。

### --auto

`--auto` 按用量而不是槽位顺序做决定。`cli.pickAuto` 保留已登录的账号，用 `ccam usage` 所用的同一个 `Fetcher.FetchAll` 一次性查询它们的用量，然后交给 `usage.Best()` 选择：

<!-- motion:next-auto -->

`usage.Best()` 选择 5 小时百分比最低的账号，并跳过用量已达 100% 或 5 小时用量未知的账号。查询遵循 `ccam usage` 的规则：60 秒内获取的数字会被复用，每个账号最多等待 3 秒，没有响应或登录已过期的账号按其保存的数字判断。重置时间已过的已保存时间窗口按 0% 计算，因此之后已经重置的账号不会被漏掉。`--auto` 从不为令牌续期。与单纯的轮换不同，如果你上次使用的账号余量仍然最多，`--auto` 可能会选中它。

### --continue

Claude Code 的 `--continue` 会打开当前文件夹中最近的一次对话，这只有在新账号能看到那次对话时才有效。因此使用 `--continue` 时，ccam 只考虑与你上次使用的账号处于同一个 `history` 组的账号，然后把 `--continue` 放在你的参数前面。历史记录为 `local` 的账号不与任何账号共享，所以从这样的账号使用 `--continue` 总会失败。请先用 `ccam share <name> --history on` 开启共享（参见[共享](/docs/sharing)）。

## 安全性与边界情况

- 没有登录的账号永远不会被选中。如果没有其他已登录的账号，单纯的 `ccam next` 会以 `no logged-in account other than <name>` 停止。
- 使用 `--continue` 但没有账号共享历史时：`` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``。
- 使用 `--auto` 但没有已登录账号时：`no logged-in account to pick`。当每个账号都已用满或用量未知时：`` no account with known usage and headroom (see `ccam usage`) ``。
- 即使带有 `--dry-run`，`--auto` 也会发出网络请求。它只读取，从不写入令牌文件。
- 所选账号通过与 `ccam <name>` 相同的路径启动，因此适用相同的检查，包括当另一个共享同一历史记录的 Claude Code 已在此文件夹中运行时发出的警告。
- `last-used` 只在 Claude Code 实际启动时写入。`--dry-run` 不会改动它，因此试运行不会推进轮换。

## 相关内容

- [用量限制](/features/usage)
- [用量限制与切换账号](/docs/usage-limits)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [用量限制工作流程](/guides/usage-limit-workflow)
