---
title: 切换到下一个账号
description: ccam next 以另一个账号打开 Claude Code；--auto 选择余量最多的账号，--continue 把对话带到新账号，--fallback 重新运行 -p 任务。
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` 无需你指定账号名，就能以另一个账号打开 Claude Code。当某个账号达到用量限制时使用它：单纯的 `ccam next` 选择你上次使用的账号之后的那个账号，`--auto` 选择余量最多的账号，`--continue` 则把此文件夹中的上一次对话带到新账号继续，无论两个账号是否共享历史记录。在 Claude Code 中，[ccam 插件](#inside-claude-code)会在用量限制让会话停下时提醒你，`/ccam:switch` 则把对话交接出去。对于运行 `claude -p` 的脚本，`ccam --fallback` 会自动完成这种切换（参见[达到限制时重新运行](#run-again-at-the-limit)）。

## 用法

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam 在标准错误上以 `→ <from> → <to>` 的形式打印这次切换，对 `--auto` 的选择还会附上理由，然后像 `ccam <name>` 一样以新账号启动 Claude Code。`<from>` 是你上次使用的账号；使用 `--continue` 时，则是那次对话所属的账号。`--` 之后的参数原样传给 `claude`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--auto` | 按用量选择已登录的账号（见下文），跳过用量已达 100% 的账号 | 关闭 |
| `--max-usage <1-100>` | 与 `--auto` 一起使用：优先选择用量低于此值的账号 | `auto_max_usage`，否则为 100 |
| `--strategy headroom\|expiring\|balance` | 与 `--auto` 一起使用：5 小时用量最低的账号、5 小时时间窗口最早重置的账号，或每周配额最先会被浪费的账号 | `auto_strategy`，否则为 `headroom` |
| `--cross-org` | 与 `--auto` 一起使用：也选择其他组织的账号（见下文） | `cross_org`，否则为关闭 |
| `--wait[=<duration>]` | 与 `--auto` 一起使用：没有账号还有余量时，等待第一个账号重置，最多等待这么久 | 关闭；单独给出时为 `5h` |
| `-c`、`--continue` | 以新账号继续此文件夹中的上一次对话（见下文） | 关闭 |
| `-- <claude args>` | 放在 ccam 自身参数之后传给 `claude` 的参数 | 无 |
| `--dry-run` | 支持：打印 `→ <from> → <to>` 和将要执行的命令行，不启动 Claude Code，也不记录启动 | 关闭 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 时，`--resume` 参数会放在你自己的参数之前。

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

`--auto` 按用量而不是槽位顺序做决定。`cli.pickAuto` 保留已登录的账号，用 `ccam usage` 所用的同一个 `Fetcher.FetchAll` 一次性查询它们的用量，然后交给 `usage.Policy.Choose()` 选择：

<!-- motion:next-auto -->

账号按其 5 小时和每周时间窗口来判断；当传给 Claude Code 的参数中有 `--model` 时，还会计入该模型的每周时间窗口（名称中含有 `opus`、`sonnet`、`haiku` 或 `fable`）。`Policy.Choose()` 把账号分为三类：

1. **有余量**：每个时间窗口都低于 `auto_max_usage`（`config.toml`，默认 100）。由策略在其中选择：`headroom`（默认）选 5 小时用量最低的账号；`expiring` 选 5 小时时间窗口最早重置的账号，让它剩余的配额在失效前被用掉；`balance` 选在每周时间窗口重置之前平均每天剩余配额最多的账号，并把每个正在运行的 ccam 会话计入其中，这样每周配额不会被浪费，并行的会话也会分散开。
2. **超出**：低于 100%，但达到或超过 `auto_max_usage`。只在没有账号有余量时使用，先选 5 小时用量最低的。
3. **已满**：某个时间窗口已达 100%，或 5 小时用量未知。从不选择。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**组织。** 自动选择只会落在与你离开的账号同类的账号上：同一个团队或企业组织，或者个人套餐（Pro、Max）之间，这样公司账号的工作不会在私人账号上继续，反之亦然。ccam 从每个账号的登录中得知其所属组织（`ccam login` 会记录它；较早的登录则从 `.claude.json` 读取）；无法判断组织的账号可以归入任何一类。登录到同一个 Claude 账号的两个文件夹共享同一份限制，所以它们永远不会互相接替。`--cross-org`，或 `config.toml` 中的 `cross_org = true`，会取消这条组织规则。

当没有可选的账号时，ccam 以退出码 75（`EX_TEMPFAIL`）退出，这样脚本就能把"稍后再试"和失败区分开。使用 `--wait` 时，它会改为等待，直到第一个已满账号中造成阻塞的时间窗口重置（重置时间之后 30 秒，此时接口会显示新的时间窗口），打印 `ccam: no account has room; waiting until 15:10`，然后重新选择；如果没有账号说明何时重置，它会在 10 分钟后再查看。如果在等待结束前没有账号空出来，它以 75 退出。查询遵循 `ccam usage` 的规则：60 秒内获取的数字会被复用，每个账号最多等待 3 秒，没有响应或登录已过期的账号按其保存的数字判断。重置时间已过的已保存时间窗口按 0% 计算，因此之后已经重置的账号不会被漏掉。`--auto` 从不为令牌续期。与单纯的轮换不同，如果你上次使用的账号余量仍然最多，`--auto` 可能会选中它。

### --continue

1. ccam 在每个账号中查找此文件夹的上一次对话：`projects/<folder>/` 中最新的 `<session>.jsonl`（文件夹名就是工作文件夹，其中字母和数字以外的每个字符都换成 `-`，与 Claude Code 的命名方式一致）。如果会话留下了一条记录（用量限制或 `/ccam:switch`，见下文），并且之后没有写入更新的内容，则改用那个会话。
2. 对话所属的账号就是你要切换离开的账号，新账号按上述规则在同一组织内选择。
3. 如果新账号共享那份历史记录（同一个组文件夹），ccam 启动 `claude --resume <session id>`，也就是那次对话本身（`--continue` 会打开文件夹中最新的对话，不管是哪一次）。
4. 否则它启动 `claude --resume <the conversation's file> --fork-session`：Claude Code 以新的会话 ID 把整个对话复制到新账号自己的历史记录中，原对话保持不变。此后新账号的 `--continue` 就能找到它。

ccam 会告诉你对话有多大：新账号在第一次回复时会重新读取全部内容。提示缓存从不在组织之间共享，所以使用 `--cross-org` 时，第一次回复要按完整大小计费。在组中共享 `history` 可以避免复制，只保留一份对话，同时也共享 `/rewind` 用来恢复代码的检查点（`file-history/`）。

## 在 Claude Code 中

[ccam 插件](/docs/install#claude-code-plugin)会添加两个钩子（Claude Code 2.1.78 或更高版本）：

- 当某一轮因用量限制而停止时（`You've hit your session limit`、每周限制或模型限制），它会记下这次对话，并在支持桌面通知的终端（iTerm2、WezTerm、Windows Terminal、kitty、Ghostty）中显示一条通知：`ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` 服务器短暂的限流或 API 密钥的速率限制不属于用量限制，不会处理。Claude Code 2.1.234 及更高版本还会在会话中等待，并在 5 小时时间窗口重置后自动继续；这条通知让你可以选择改为切换账号。
- `/ccam:switch [account]` 会记下这次对话，并在发出任何请求之前结束这一轮：输入 `/exit`，然后 `ccam next --continue` 会以你指定的账号或下一个账号继续这次对话。

这条记录保存在 `~/.local/state/ccam/handoff/<account>.json` 中，有效期 12 小时，只在下一个会话即将启动时使用一次；启动失败时它会保留下来，供下次尝试。它属于会话启动时所在的文件夹，即使之后在 Claude Code 的 shell 中执行了 `cd` 也是如此。如果它指定的账号后来被重命名或删除，ccam 会说明这一点，并照常选择。在 `config.toml` 中设置 `relay = true` 时，ccam 会留在交互式会话旁边，而不是用 Claude Code 替换自身；在 `/ccam:switch` 之后，它会自己继续那个会话的对话（不会处理其他会话留下的记录）。代价是每个会话多一个等待中的 ccam 进程，所以默认关闭。接力的会话以 `ccam next --continue` 启动，不带第一个会话启动时使用的参数，例如 `--model`。

## 达到限制时重新运行

`ccam next` 用于已经停下来的交互式会话。非交互式运行，即脚本或定时任务中的 `claude -p`，无法等你操作，所以在账号名之前加上 `--fallback`，就能让 ccam 替它切换：

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` 只能与 `-p` 或 `--print` 一起使用；否则 ccam 会停止，并提示使用 `ccam next --auto --continue`。
2. 标准输入只读取一次，并提供给每一次尝试。
3. ccam 以子进程方式启动 `claude`，而不是使用 `exec`；它会暂存标准输出，标准错误则随产生随时传出。
4. 当 `claude` 以 0 退出，或因其他原因失败时，暂存的输出和退出码会原样传递。
5. 当 Claude Code 的消息表明达到了用量限制（`You've hit your limit`、`usage limit reached`），或者该账号的用量查询显示其 5 小时或每周时间窗口已达 100% 时，这次失败就被视为达到用量限制。
6. 此时会丢弃这次尝试的输出，并以另一个已登录账号运行同一个命令：先选 5 小时用量最低的账号，然后是无法查询用量的账号，按编号顺序。用量已达 100% 的账号会被跳过，每个账号最多尝试一次。
7. 如果参数中有 `-c`、`--continue`、`-r` 或 `--resume`，只会尝试同一历史组中的账号，以便找到那次对话。
8. 与 `--auto` 一样，只会尝试同一组织的账号（`cross_org = true` 会取消这一限制）。

下一个账号按与 `--auto` 相同的规则选择（`auto_max_usage`、`auto_strategy`、`--model` 的时间窗口）。当没有剩余账号时，ccam 会打印最后一次尝试的输出，说明没有已登录账号还有余量，并以退出码 75 退出。在账号名之前加上 `--wait`（`ccam --fallback --wait=2h work -p …`）时，它会改为等待第一个账号重置，然后可能再次尝试每个账号；它总是先至少等待一次，所以刚刚达到限制的账号不会立刻再次运行。`ccam --dry-run --fallback work -p …` 会打印第一次尝试的命令。

交互式会话不会自动切换。作为替代，使用 `ccam statusline install` 后，当 5 小时或每周用量达到 90% 时，状态栏会加上 `→ ccam next --auto --continue`（[用量限制](/features/usage)）。

## 安全性与边界情况

- 没有登录的账号永远不会被选中。如果没有其他已登录的账号，单纯的 `ccam next` 会以 `no logged-in account other than <name>` 停止。
- 使用 `--continue` 但文件夹中没有对话时：`no account has a conversation in <folder> to continue`。
- 当组织规则排除了所有账号时：`` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``。
- 使用 `--auto` 但没有已登录账号时：`no logged-in account to pick`。当每个账号都已用满或用量未知时：`` no account with known usage and headroom (see `ccam usage`) ``。两者都以退出码 75 结束。
- `--max-usage`、`--strategy`、`--wait` 和 `--cross-org` 只能与 `--auto` 一起使用，超出 1-100 的值或其他策略会被拒绝。
- 复刻（fork）只复制对话，不复制账号的其他状态：除非共享历史记录，原账号 `file-history/` 中的检查点不会一起带过去。
- 即使带有 `--dry-run`，`--auto` 也会发出网络请求。它只读取，从不写入令牌文件。
- 所选账号通过与 `ccam <name>` 相同的路径启动，因此适用相同的检查，包括当另一个共享同一历史记录的 Claude Code 已在此文件夹中运行时发出的警告。
- `last-used` 只在 Claude Code 实际启动时写入。`--dry-run` 不会改动它，因此试运行不会推进轮换。
- `--fallback` 会暂存标准输出直到一次尝试结束，所以 `--output-format stream-json` 这样的流式输出会在最后一次性到达。第一次尝试在达到限制之前修改过的文件会保持修改后的状态；下一次尝试会重新开始执行提示。

## 相关内容

- [用量限制](/features/usage)
- [用量限制与切换账号](/docs/usage-limits)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [用量限制工作流程](/guides/usage-limit-workflow)
