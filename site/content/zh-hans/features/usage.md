---
title: 用量限制
description: ccam usage 一次显示每个账号的 5 小时限额和每周限额，ccam statusline 在 Claude Code 中显示账号名和用量。
updated: 2026-10-24
order: 7
translation: machine
sourceHash: "a14b1da71a1a"
---

`ccam usage` 显示每个账号已用掉多少 5 小时限额、每周限额以及各模型自己的每周限额（例如 Fable），每项限额一个用量条，并附上距离重置的时间。在开始一项长任务之前，或者当 Claude Code 提示你快到限额时，用它查看哪个账号还有余量。`ccam statusline` 是供 Claude Code 使用的状态栏命令：它在 Claude Code 中显示账号名和用量，并为 `ccam usage` 保存这些数字。

## 用法

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

每个账号先有一行，显示编号、名称以及这些数字是多久之前的；`(cached)` 表示它们来自保存的文件，而不是来自接口。下面每项限额占一行，显示 20 格的用量条、百分比和距离窗口重置的时间：`5-hour`、`Weekly`（所有模型），然后是各模型自己的每周窗口，例如 `Fable`。账号没有的窗口显示 `-`。百分比从 80% 起变为黄色，达到 100% 时变为红色。当 Claude Code 因用量限制拒绝请求时，它报告已满的时间窗口在重置前显示 100%，无论保存的数字是多少（参见[用量限制记录](#the-limit-record)）。在终端中，所有账号会先一起显示为 `fetching…`，每个账号返回结果时再绘制它的各行。最后一行列出 [`ccam next --auto`](/features/next) 按相同规则会选择的账号（`auto_max_usage`、`auto_strategy`、正在运行的会话；除非设置了 `cross_org`，否则只考虑上次使用的账号所属组织中的账号；跳过任何时间窗口已达 100% 的账号，因此每周额度已用满的账号永远不会被列出），只有在你有多个账号时才会出现。使用默认的 `headroom` 策略时，这一行为 `→ Most headroom: <name>`；使用 `expiring` 或 `balance` 时，它会给出理由，例如 `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`。

**Team 和 Enterprise 席位。** 自动选择（`ccam next`、`--fallback`、`ccam pick`、仪表盘中的 `n`）也会在同一个 Claude Team 或 Enterprise 组织的席位之间切换。Anthropic 帮助中心关于 [Team 套餐](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的说明写道："Each seat on a Team plan is assigned to one named individual, and only that person may use it"（Team 套餐的每个席位都分配给一位指定的个人，只有此人可以使用），以及 "Users can't route their work through another user's account to use that user's remaining usage"（用户不能通过另一位用户的账号完成自己的工作，以使用该用户剩余的用量）；关于 [Enterprise 套餐](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的说明也是如此。因此，在席位之间使用自动选择可能违反 Anthropic 的条款。在依赖这项功能之前，请先向你所在的组织和 Anthropic 确认。个人套餐（Pro、Max）的情况请参见[常见问题中关于条款的说明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)。

未登录的账号显示 `not logged in → ccam login <name>`，使用 API 密钥的账号显示 `no limits (API)`。访问令牌过期时显示 ``access token expired; `ccam <name>` renews it when it starts``，如果有保存的数字，就显示在这些数字上方：Claude Code 会在下次启动时为它续期。整个登录已过期时（因为其刷新令牌已过期，或者服务器拒绝后 Claude Code 清除了它），显示 `` login expired → `ccam login <name>` ``。

运行 `ccam statusline install` 之后，Claude Code 会显示 `[work] · Opus · 5h 24% · 7d 41%` 这样的一行。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--json` | 以 JSON 输出报告：`name`、`slot`、`status`、`limits`（账号开启了用量积分时带有 `credits`）、`source`、`at`、`error`，以及在用量限制记录把时间窗口设为 100% 时的 `limited` 和 `limited_windows`。没有任何账号时输出 `[]` | 关闭 |
| `--watch[=N]` | 每 `N` 秒重绘一次，直到你按下 Ctrl+C。单独的 `--watch`，或者值小于等于 0，表示 60。值写在 `=` 之后：`--watch 10` 会被拒绝，并给出提示 | 关闭 |
| `-v`、`--verbose` | 在每个账号下添加一行 `source:`，说明数字来自哪里（用量接口、Claude Code、状态栏、保存的文件，或 Claude Code 报告的用量限制）。当用量限制记录覆盖在这些数字之上时，后面会跟着 `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` 这样的一行 | 关闭 |
| `--no-cache` | 忽略 60 秒缓存，重新查询接口 | 关闭 |
| `--source live\|claude\|cache` | `live` 用账号的令牌查询接口；`claude` 让 Claude Code 去查询（见下文）；`cache` 只显示已保存的数字 | `usage_source`，否则为 `live` |
| `--refresh` | 为过期的访问令牌续期，仅限未在运行的账号。一次续期最多可能需要 30 秒 | 关闭 |
| `--dry-run` | 接受：`ccam usage` 除了自己的缓存之外不更改任何内容。与 `--refresh` 一起使用时会被拒绝：`--dry-run cannot be combined with --refresh, which rewrites login tokens` | 关闭 |

`--source` 的其他任何值都会失败并显示 `--source must be live, claude or cache`。

### 由 Claude Code 自己查询用量

在 `config.toml` 的顶部（第一个 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之上）设置 `usage_source = "claude"`，或使用 `--source claude` 时，ccam 完全不读取登录令牌。它为每个账号在账号文件夹中启动 `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`，并通过标准输入发送 Claude Code 的 `get_usage` 请求（Claude Code 2.1.270 或更高版本）；Claude Code 用自己的登录查询用量，并返回 5 小时、每周和各模型的时间窗口。`--safe-mode` 让账号的钩子、MCP 服务器和插件不会启动，也不会保存会话。每个账号都要启动一次 Claude Code，大约一秒，所以这个来源比 `live` 慢；所有账号会同时查询，每个账号限时 20 秒，60 秒缓存照常适用。访问令牌已过期或将在 5 分钟内过期的登录有 50 秒时间，因为 Claude Code 会在启动时为它续期，而被打断的续期可能会丢失登录。放弃的查询会用 SIGTERM 请求 Claude Code 结束，这样它会完成已经开始的续期（Windows 没有这种信号，在那里它会被直接停止）。整个登录已过期时，会显示 `` login expired → `ccam login <name>` ``，不会启动 Claude Code。Claude Code 把这个请求标为实验性功能，所以如果它不再响应，ccam 会显示已保存的数字，`ccam doctor` 也会报告这一情况。当 Claude Code 回答它自己的查询失败时，ccam 会显示已保存的数字并附上 `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`，而不是显示为没有限制的套餐。`--dry-run` 不会启动 Claude Code（它会写入账号文件夹），而是显示已保存的数字。`usage_source` 接受 `live` 或 `claude`；`cache` 只能作为 `--source` 的选项。`ccam usage` 不接受参数。

`ccam statusline` 从标准输入读取 Claude Code 的状态栏 JSON。它唯一的选项是 `--wrap <command>`：一条现有的状态栏命令，其输出的第一行前面会加上账号名。`ccam statusline install` 没有选项。两者都会拒绝 `--dry-run`，并显示 `` --dry-run is not supported by `ccam statusline`; nothing was changed ``。

## 工作原理

这些数字来自两个地方：Anthropic 没有公开文档的用量接口 `GET https://api.anthropic.com/api/oauth/usage`（用每个账号自己的登录调用），以及 Claude Code 传给其状态栏命令的值。所有相关代码都在 `internal/usage` 中。

1. `cli.usageCmd` 加载 `config.toml`，并为每个账号构建一个 `usage.Account`：它的文件夹 `~/.local/share/ccam/profiles/<name>/`、它的会话标记（结合 Claude Code 自己在账号文件夹中的记录，用来判断它是否正在运行），以及它是否是原地导入的。
2. `Fetcher.FetchAll` 为每个账号启动一个 goroutine。它们共用一个 `http.Client`，因此到接口的连接会被复用。建立连接最多 2 s，每个账号总共 3 s。客户端不跟随重定向，因为重定向会把令牌再次发送到别处；它最多读取 1 MB 的响应。
3. 每个账号完成时，重绘它所在的行。失败或未能按时响应的账号会退回到保存的数字，因此一个慢的账号永远不会拖住其他账号。

每个账号在自己的 goroutine 中查询接口；超时的那个账号改为显示保存的数字。

<!-- motion:usage-fetch -->

`personal` 和 `work` 按时响应，显示实时数字。`side` 耗时超过 3 s，因此 `fallback()` 改为显示它保存的数字。`usage.Table()` 为每个账号绘制一块，`usage.Best()` 选出最后一行列出的账号。

对于单个账号，`Fetcher.Fetch` 按以下顺序进行检查：

<!-- motion:usage-cache -->

1. 使用 `--source cache` 时，直接进入 `fallback()`。
2. 除非设置了 `--no-cache`，否则 `~/.local/state/ccam/usage/<name>.json` 中测量时间不到 60 秒（`CacheTTL`）的实时结果会被使用，并按当前时间判断：此后已经重置的时间窗口显示 0%，日期晚于当前时钟的文件不算新鲜。连续运行两次 `ccam usage`，每个账号只会发出一次请求。
3. `readCredentials()` 读取 `.credentials.json`（或 macOS 钥匙串）而不修改它，只解码 `claudeAiOauth` 键。没有订阅登录时显示 `not logged in`；如果 `.claude.json` 中保存的是 API 密钥，则显示 `no limits (API)`。刷新令牌已过期或已被 Claude Code 清除的登录显示为 `` login expired → `ccam login <name>` ``，不会查询。
4. 过期的访问令牌不会被使用。不带 `--refresh` 时，该账号显示保存的数字，并附上 ``access token expired; `ccam <name>` renews it when it starts``。
5. 否则 ccam 调用接口。401 或 403 视为登录已过期（`login expired: authentication failed (401)`），其他任何状态码或无法读取的响应体视为错误，两者都会退回到保存的数字。被 Ctrl+C 结束的查询显示 `cancelled`。429 表示对接口的请求过于频繁：ccam 会写入 `<name>.pause`，记下再次请求的时间（`Retry-After`，否则为 5 分钟，最多一小时），显示保存的数字并附上 `rate limited by the usage endpoint; asking again after 15:10`，并且在那之前不会为该账号调用接口，即使使用 `--no-cache` 也是如此。立即重试只会让限制持续下去。
6. 正常的响应会连同测量时间一起写入 `<name>.json`。

`fallback()` 取最新保存的查询结果（`<name>.json`，或来自 Claude Code 查询的 `<name>.claude.json`）；如果 `<name>.statusline.json` 是在它之后写入的，就把状态栏的 5 小时和每周时间窗口覆盖在上面；状态栏没有的各模型时间窗口则保持查询结果保存时的样子。在显示保存的数字之前，它会检查每个时间窗口的重置时间：在这些数字保存之后已经重置的时间窗口显示 0%，因为旧的百分比已不再适用。这就是上面的 `personal` 显示 `0%` 且没有重置时间的原因。

接口有两种响应格式。较新的格式有一个 `limits` 数组，其中包含 `session`（5 小时）、`weekly_all` 和 `weekly_scoped`（按模型）条目。较旧的格式有 `five_hour`、`seven_day`、`seven_day_opus` 和 `seven_day_sonnet` 对象。`ParseLive` 两种都能读取：以 `limits` 数组为准，旧对象补充它没有涵盖的部分。完全没有限额的响应体会被报告为格式已变化，`ccam doctor` 会告诉你这一点。

### 用量限制记录

当 Claude Code 因用量限制停止某一轮时，[ccam 插件](/docs/install#claude-code-plugin)的钩子会从对话中读取被拒绝的请求：哪个时间窗口已满（5 小时、每周、Opus 或 Sonnet）以及它何时重置。它把这些写入 `~/.local/state/ccam/usage/<name>.limit.json`。`Fetch` 会把这条记录覆盖在它从所有来源得知的其他内容之上：记录中的每个时间窗口在重置前都按 100% 计算，无论保存的数字多旧或多低，因此 `ccam next --auto`、`ccam pick` 和 `--fallback` 在那之前不会再次选择该账号。如果没有其他已知信息，账号显示的就只有这条记录（使用 `-v` 时为 `source: the usage limit Claude Code reported`）。无人值守的会话（`-p`、后台会话）不写入记录，组织的停用（用量耗尽、组的额度设为 $0）也不写入，因为组织管理员随时可以解除它。`ccam login` 登录到另一个 Claude 账号时会丢弃这条记录；删除该文件也可以清除它。

### 为令牌续期

刷新令牌是轮换的：续期一次会得到新的刷新令牌，旧的随之失效。如果 ccam 在 Claude Code 正以同一账号运行时为登录续期，该会话就会被登出。因此 `ccam usage` 只读取令牌，从不写入，除非你加上 `--refresh`。此时 `Fetcher.refresh` 会：

1. 拒绝用 `ccam import` 原地注册的账号，因为直接运行的 `claude` 可能会在 ccam 不知情的情况下使用那个登录；
2. 拒绝保存在 macOS 钥匙串中的登录；
3. 获取一个锁文件 `~/.local/state/ccam/locks/profile-<name>.refresh.lock`，这样两个 ccam 进程永远不会对同一个令牌续期两次；持有锁期间，该账号显示 `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`。`ccam rename`、`remove` 和 `logout` 也会获取这个锁，并在续期进行时拒绝执行；
4. 在持有锁的情况下，如果该账号有会话正在运行则拒绝，无论它是怎样启动的（通过 ccam、直接运行的 `claude`、IDE，或者后台会话及其背后会自行续期登录的守护进程），如果另一个进程已经为令牌续期则跳过请求；
5. 当 Claude Code 在账号文件夹中持有它自己的续期锁（`.oauth_refresh.lock`、`<folder>.lock`）时，不发送任何内容：此时 Claude Code 正在为该登录续期；
6. 把刷新令牌发送到令牌接口，并最多等待 30 秒的响应，这也是 Claude Code 给自己的续期留出的时间。请求一旦发出就永远不会被打断，即使按 Ctrl+C 也不会，因为服务器可能已经轮换了刷新令牌，只有它的响应还保有该登录：Ctrl+C 或挂断会结束其他查询，ccam 会在保存响应后退出（按 Ctrl+C 后退出码为 130）；
7. 在 Claude Code 的写入锁（`.storage-write.lock`）下保存，并且只更改当时 `.credentials.json` 中的令牌字段：Claude Code 在此期间写入的内容（MCP 登录）会保留，Claude Code 在此期间续期的登录会被保留，退出登录也永远不会被撤销。当响应给出新刷新令牌的有效期时，会像 Claude Code 一样把它保存为 `refreshTokenExpiresAt`。

没有成功的续期会说明原因，显示在账号那一行或其保存的数字上方：例如 `token not renewed: a session is running`；如果请求已经发出，但其响应无法使用或保存，则显示 `login may be lost: …`，并建议在 Claude Code 随后要求你登录时运行 `ccam login <name>`。

### 状态栏收集器

Claude Code 运行状态栏命令，并通过管道传给它一个 JSON 文档。ccam 从中读取模型名称和 `rate_limits`：

<!-- motion:usage-statusline -->

1. `runStatusline` 根据 `CCAM_PROFILE`（ccam 启动 Claude Code 时设置）找到账号，否则通过把 `CLAUDE_CONFIG_DIR` 与各账号文件夹进行匹配来确定。
2. `ParseStatusLine` 读取 `rate_limits.five_hour` 和 `rate_limits.seven_day`。对于使用 API 密钥的账号，或者在会话的第一次响应之前，JSON 中没有速率限制；这时不会保存任何内容。
3. `SaveStatusLine` 写入 `~/.local/state/ccam/usage/<name>.statusline.json`。当接口没有响应时，`ccam usage` 会使用它。
4. ccam 打印 `[name] · model · 5h N% · 7d N%`。使用 `--wrap` 时，它会通过 `sh -c` 以相同的输入运行你的命令，并按 Claude Code 显示状态栏的方式打印其输出：每一个非空行，第一行前面加上 `[name]`。命令失败时只显示账号名。在 Windows 上，它按 Claude Code 的方式运行命令：安装了 Git Bash 时使用 Git Bash，否则使用 PowerShell。
5. 当 5 小时或每周用量达到 90% 或以上，并且还有其他账号时，这一行末尾会加上 `→ ccam next --auto --continue`，例如 `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`。

被包装的命令没有自己的时间限制：Claude Code 给状态栏命令 10 分钟时间，并在新的更新开始时停止正在进行的运行，因此慢的命令的表现与不使用 ccam 时一样。当 Claude Code 结束 ccam 时，该命令也会结束（在 Linux 和 macOS 上，连同它启动的所有进程）。

`ccam statusline install` 把 `statusLine` 写入每个账号的 `settings.json`，并保留该对象的其他键，例如 `padding`、`refreshInterval` 和 `hideVimModeIndicator`。如果账号已有状态栏命令，ccam 会把它保留为 `--wrap '<old command>'`；在 Windows 上则保留为 `--wrap-base64 <base64>`，在两种 shell 中都不需要加引号。已经在运行某个 ccam 的状态栏（来自其他路径，或者是旧版 ccam 按其方式加了引号的写法）会被替换而不是再包装一层，并保留它所包装的内容。在 Windows 上，ccam 的路径用正斜杠书写，主文件夹写成 `~`（`~/.local/bin/ccam.exe statusline`），Git Bash 和 PowerShell 都能识别。需要加引号的路径会按 Claude Code 所用的 shell 完整写出：Git Bash 用 `'…'`，PowerShell 用 `& '…'`。通过组共享的 `settings.json` 只会经由其链接写入一次，因此组中的每个账号都会得到它。再次运行时，对已经使用 ccam 的文件会打印 `(already set)`。

## 安全性与边界情况

- 除非你加上 `--refresh`，否则 `ccam usage` 从不修改令牌文件。访问令牌过期时，用 `ccam <name>` 启动一次该账号，Claude Code 就会为它续期。整个登录已过期时需要运行 `ccam login <name>`。
- `--refresh` 从不触碰正在运行的账号（包括不通过 ccam 启动的 Claude Code）、导入的账号或钥匙串中的登录。续期被拒绝时，会显示保存的数字和原因（`token not renewed: …`）。
- 用量接口没有公开文档，可能会发生变化。发生变化时，你仍能看到保存的数字，`ccam doctor` 会报告格式变化。
- 既没有保存的数字也没有响应时，该账号的那一行显示 `lookup failed:` 以及原因，例如使用 `--source cache` 时的 `no saved usage`。
- 缓存文件只保存百分比和重置时间，从不保存令牌。

## 相关内容

- [用量限制与切换账号](/docs/usage-limits)
- [切换到下一个账号](/features/next)
- [用量限制工作流程](/guides/usage-limit-workflow)
- [命令](/docs/commands)
