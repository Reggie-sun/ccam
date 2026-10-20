---
title: 用量限制
description: ccam usage 一次显示每个账号的 5 小时限额和每周限额，ccam statusline 在 Claude Code 中显示账号名和用量。
updated: 2026-10-20
order: 7
translation: machine
sourceHash: f08a2fdc6f0b
---

`ccam usage` 显示每个账号已用掉多少 5 小时限额和每周限额，以及每个时间窗口距离重置还有多久。在开始一项长任务之前，或者当 Claude Code 提示你快到限额时，用它查看哪个账号还有余量。`ccam statusline` 是供 Claude Code 使用的状态栏命令：它在 Claude Code 中显示账号名和用量，并为 `ccam usage` 保存这些数字。

## 用法

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

在终端中，表格会立即出现，每一行显示 `fetching…`，每个账号返回结果时重绘对应的行。`As of` 表示这些数字是多久之前的；`(cached)` 表示它们来自保存的文件，而不是来自接口。最后一行列出 5 小时用量低于 100% 且最低的账号，只有在你有多个账号时才会出现。未登录的账号显示 `not logged in → ccam login <name>`，使用 API 密钥的账号显示 `no limits (API)`。

运行 `ccam statusline install` 之后，Claude Code 会显示 `[work] · Opus · 5h 24% · 7d 41%` 这样的一行。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--json` | 以 JSON 输出报告：`name`、`slot`、`status`、`limits`、`source`、`at`、`error` | 关闭 |
| `--watch[=N]` | 每 `N` 秒重绘一次，直到你按下 Ctrl+C。单独的 `--watch`，或者值小于等于 0，表示 60 | 关闭 |
| `-v`、`--verbose` | 显示按模型区分的每周限额，以及一行 `source:`（`live`、`cache` 或 `statusline`） | 关闭 |
| `--no-cache` | 忽略 60 秒缓存，重新查询接口 | 关闭 |
| `--source live\|cache` | `cache` 从不查询接口，只显示已保存的数字；`live` 是默认行为 | live |
| `--refresh` | 为过期的访问令牌续期，仅限未在运行的账号 | 关闭 |
| `--dry-run` | 接受：`ccam usage` 除了自己的缓存之外不更改任何内容。与 `--refresh` 一起使用时会被拒绝：`--dry-run cannot be combined with --refresh, which rewrites login tokens` | 关闭 |

`--source` 的其他任何值都会失败并显示 `--source must be live or cache`。`ccam usage` 不接受参数。

`ccam statusline` 从标准输入读取 Claude Code 的状态栏 JSON。它唯一的选项是 `--wrap <command>`：一条现有的状态栏命令，其输出的第一行前面会加上账号名。`ccam statusline install` 没有选项。两者都会拒绝 `--dry-run`，并显示 `` --dry-run is not supported by `ccam statusline`; nothing was changed ``。

## 工作原理

这些数字来自两个地方：Anthropic 没有公开文档的用量接口 `GET https://api.anthropic.com/api/oauth/usage`（用每个账号自己的登录调用），以及 Claude Code 传给其状态栏命令的值。所有相关代码都在 `internal/usage` 中。

1. `cli.usageCmd` 加载 `config.toml`，并为每个账号构建一个 `usage.Account`：它的文件夹 `~/.local/share/ccam/profiles/<name>/`、它的会话文件夹（用来判断它是否正在运行），以及它是否是原地导入的。
2. `Fetcher.FetchAll` 为每个账号启动一个 goroutine。它们共用一个 `http.Client`，因此到接口的连接会被复用。建立连接最多 2 s，每个账号总共 3 s。客户端不跟随重定向，因为重定向会把令牌再次发送到别处；它最多读取 1 MB 的响应。
3. 每个账号完成时，重绘它所在的行。失败或未能按时响应的账号会退回到保存的数字，因此一个慢的账号永远不会拖住其他账号。

每个账号在自己的 goroutine 中查询接口；超时的那个账号改为显示保存的数字。

<!-- motion:usage-fetch -->

`personal` 和 `work` 按时响应，显示实时数字。`side` 耗时超过 3 s，因此 `fallback()` 读取它的两个缓存文件，使用最后写入的那个。`usage.Table()` 绘制各行，`usage.Best()` 选出最后一行列出的账号。

对于单个账号，`Fetcher.Fetch` 按以下顺序进行检查：

<!-- motion:usage-cache -->

1. 使用 `--source cache` 时，直接进入 `fallback()`。
2. 除非设置了 `--no-cache`，否则 `~/.local/state/ccam/usage/<name>.json` 中不到 60 秒（`CacheTTL`）的实时结果会被原样返回。连续运行两次 `ccam usage`，每个账号只会发出一次请求。
3. `readCredentials()` 读取 `.credentials.json`（或 macOS 钥匙串）而不修改它，只解码 `claudeAiOauth` 键。没有订阅登录时显示 `not logged in`；如果 `.claude.json` 中保存的是 API 密钥，则显示 `no limits (API)`。
4. 过期的访问令牌不会被使用。不带 `--refresh` 时，该账号显示保存的数字，并标记为已过期。
5. 否则 ccam 调用接口。401 或 403 视为登录已过期，其他任何状态码或无法读取的响应体视为错误，两者都会退回到保存的数字。
6. 正常的响应会连同测量时间一起写入 `<name>.json`。

`fallback()` 同时读取 `<name>.json` 和 `<name>.statusline.json`，并保留较新的那个。在显示保存的数字之前，它会检查每个时间窗口的重置时间：在这些数字保存之后已经重置的时间窗口显示 0%，因为旧的百分比已不再适用。这就是上面的 `personal` 显示 `0%` 且没有重置时间的原因。

接口有两种响应格式。较新的格式有一个 `limits` 数组，其中包含 `session`（5 小时）、`weekly_all` 和 `weekly_scoped`（按模型）条目。较旧的格式有 `five_hour`、`seven_day`、`seven_day_opus` 和 `seven_day_sonnet` 对象。`ParseLive` 两种都能读取：以 `limits` 数组为准，旧对象补充它没有涵盖的部分。完全没有限额的响应体会被报告为格式已变化，`ccam doctor` 会告诉你这一点。

### 为令牌续期

刷新令牌是轮换的：续期一次会得到新的刷新令牌，旧的随之失效。如果 ccam 在 Claude Code 正以同一账号运行时为登录续期，该会话就会被登出。因此 `ccam usage` 只读取令牌，从不写入，除非你加上 `--refresh`。此时 `Fetcher.refresh` 会：

1. 拒绝用 `ccam import` 原地注册的账号，因为直接运行的 `claude` 可能会在 ccam 不知情的情况下使用那个登录；
2. 拒绝保存在 macOS 钥匙串中的登录；
3. 获取一个锁文件 `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock`（该变量未设置时位于 `~/.local/state/ccam/locks/` 下），这样两个 ccam 进程永远不会对同一个令牌续期两次；
4. 在持有锁的情况下，如果该账号有会话正在运行则拒绝，如果另一个进程已经为令牌续期则跳过请求；
5. 把刷新令牌发送到令牌接口，并用 `fsutil.AtomicWrite` 把新令牌写回 `.credentials.json`，保留文件中的其他所有键。

### 状态栏收集器

Claude Code 运行状态栏命令，并通过管道传给它一个 JSON 文档。ccam 从中读取模型名称和 `rate_limits`：

<!-- motion:usage-statusline -->

1. `runStatusline` 根据 `CCAM_PROFILE`（ccam 启动 Claude Code 时设置）找到账号，否则通过把 `CLAUDE_CONFIG_DIR` 与各账号文件夹进行匹配来确定。
2. `ParseStatusLine` 读取 `rate_limits.five_hour` 和 `rate_limits.seven_day`。对于使用 API 密钥的账号，或者在会话的第一次响应之前，JSON 中没有速率限制；这时不会保存任何内容。
3. `SaveStatusLine` 写入 `~/.local/state/ccam/usage/<name>.statusline.json`。当接口没有响应时，`ccam usage` 会使用它。
4. ccam 打印 `[name] · model · 5h N% · 7d N%`。使用 `--wrap` 时，它会通过 `sh -c` 以相同的输入运行你的命令，改为打印 `[name]` 加上该命令输出的第一行。在 Windows 上，它按 Claude Code 的方式运行命令：安装了 Git Bash 时使用 Git Bash，否则使用 PowerShell。
5. 当 5 小时或每周用量达到 90% 或以上，并且还有其他账号时，这一行末尾会加上 `→ ccam next --auto --continue`，例如 `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`。

被包装的命令有 2 秒时间。超时后它会被终止，这一行只显示账号名，因此慢的命令无法阻塞 Claude Code 的状态栏。

`ccam statusline install` 把 `statusLine` 写入每个账号的 `settings.json`。如果账号已有状态栏命令，ccam 会把它保留为 `--wrap '<old command>'`。在 Windows 上，ccam 的路径用正斜杠书写，主文件夹写成 `~`（`~/.local/bin/ccam.exe statusline`），Git Bash 和 PowerShell 都能识别。通过组共享的 `settings.json` 只会经由其链接写入一次，因此组中的每个账号都会得到它。再次运行时，对已经使用 ccam 的文件会打印 `(already set)`。

## 安全性与边界情况

- 除非你加上 `--refresh`，否则 `ccam usage` 从不修改令牌文件。登录过期时，用 `ccam <name>` 启动一次该账号，Claude Code 就会为它续期。
- `--refresh` 从不触碰正在运行的账号、导入的账号或钥匙串中的登录。续期被拒绝时，会显示保存的数字和原因。
- 用量接口没有公开文档，可能会发生变化。发生变化时，你仍能看到保存的数字，`ccam doctor` 会报告格式变化。
- 既没有保存的数字也没有响应时，该行显示 `lookup failed:` 以及原因，例如使用 `--source cache` 时的 `no saved usage`。
- 缓存文件只保存百分比和重置时间，从不保存令牌。

## 相关内容

- [用量限制与切换账号](/docs/usage-limits)
- [切换到下一个账号](/features/next)
- [用量限制工作流程](/guides/usage-limit-workflow)
- [命令](/docs/commands)
