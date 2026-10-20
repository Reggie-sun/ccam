---
title: 用量限制
description: ccam usage 一次顯示每個帳號的 5 小時和每週限制，ccam statusline 則在 Claude Code 中顯示帳號名稱和用量。
updated: 2026-10-20
order: 7
translation: machine
sourceHash: f08a2fdc6f0b
---

`ccam usage` 會顯示每個帳號已使用了多少 5 小時限制和每週限制，以及每個時間區間距離重設的時間。在開始一項長時間的工作之前，或 Claude Code 告訴你快要達到限制時，可以用它查看哪個帳號還有剩餘額度。`ccam statusline` 是 Claude Code 的狀態列指令：它會在 Claude Code 中顯示帳號名稱和用量，並為 `ccam usage` 儲存這些數字。

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

在終端機中，表格會立即出現，每一列顯示 `fetching…`，並在各帳號回應時重新繪製該列。`As of` 表示數字的新舊程度；`(cached)` 代表數字來自已儲存的檔案，而不是端點。最後一行會列出 5 小時用量最低且低於 100% 的帳號，只在你有多個帳號時才會出現。未登入的帳號會顯示 `not logged in → ccam login <name>`，API 金鑰帳號則顯示 `no limits (API)`。

執行 `ccam statusline install` 之後，Claude Code 會顯示 `[work] · Opus · 5h 24% · 7d 41%` 這樣的一行。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--json` | 以 JSON 輸出報告：`name`、`slot`、`status`、`limits`、`source`、`at`、`error` | 關閉 |
| `--watch[=N]` | 每 `N` 秒重新繪製，直到你按下 Ctrl+C。只寫 `--watch`，或值為 0 以下時，代表 60 | 關閉 |
| `-v`、`--verbose` | 顯示各模型的每週限制，以及 `source:` 一行（`live`、`cache` 或 `statusline`） | 關閉 |
| `--no-cache` | 忽略 60 秒快取，重新查詢端點 | 關閉 |
| `--source live\|cache` | `cache` 絕不查詢端點，只顯示已儲存的數字；`live` 是預設行為 | live |
| `--refresh` | 更新過期的存取權杖，只針對沒有在執行的帳號 | 關閉 |
| `--dry-run` | 接受：`ccam usage` 除了自己的快取以外不變更任何東西。和 `--refresh` 一起使用時會被拒絕：`--dry-run cannot be combined with --refresh, which rewrites login tokens` | 關閉 |

其他任何 `--source` 值都會失敗並顯示 `--source must be live or cache`。`ccam usage` 不接受參數。

`ccam statusline` 從標準輸入讀取 Claude Code 的狀態列 JSON。它唯一的選項是 `--wrap <command>`：一個現有的狀態列指令，ccam 會在它輸出的第一行前面加上帳號名稱。`ccam statusline install` 沒有選項。兩者都會以 `` --dry-run is not supported by `ccam statusline`; nothing was changed `` 拒絕 `--dry-run`。

## 運作方式

數字來自兩個地方：Anthropic 沒有公開文件的用量端點 `GET https://api.anthropic.com/api/oauth/usage`（以每個帳號自己的登入資訊呼叫），以及 Claude Code 交給狀態列指令的值。這些全都位於 `internal/usage`。

1. `cli.usageCmd` 載入 `config.toml`，並為每個帳號建立一個 `usage.Account`：它的資料夾 `~/.local/share/ccam/profiles/<name>/`、它的工作階段資料夾（用來判斷是否正在執行），以及它是否為就地匯入。
2. `Fetcher.FetchAll` 為每個帳號啟動一個 goroutine。它們共用同一個 `http.Client`，所以連到端點的連線會被重複使用。建立連線最多 2 秒，每個帳號總共有 3 秒。用戶端不跟隨重新導向，因為重新導向會把權杖重新送到別的地方，而且最多只讀取 1 MB 的回應。
3. 每個帳號完成時，會重新繪製它的那一列。失敗或沒有及時回應的帳號會改用已儲存的數字，所以一個緩慢的帳號絕不會拖住其他帳號。

每個帳號都在自己的 goroutine 中查詢端點；逾時的那個帳號改為顯示已儲存的數字。

<!-- motion:usage-fetch -->

`personal` 和 `work` 及時回應，顯示即時的數字。`side` 花了超過 3 秒，所以 `fallback()` 讀取它的兩個快取檔案，並使用最後寫入的那一個。`usage.Table()` 繪製各列，`usage.Best()` 選出最後一行列出的帳號。

對於單一帳號，`Fetcher.Fetch` 依序進行以下檢查：

<!-- motion:usage-cache -->

1. 加上 `--source cache` 時，直接交給 `fallback()`。
2. 除非設定了 `--no-cache`，否則 `~/.local/state/ccam/usage/<name>.json` 中不到 60 秒（`CacheTTL`）的即時結果會直接傳回。連續執行兩次 `ccam usage`，每個帳號只會花一次請求。
3. `readCredentials()` 讀取 `.credentials.json`（或 macOS 鑰匙圈）而不變更它，只解碼 `claudeAiOauth` 鍵。沒有訂閱登入時代表 `not logged in`；如果 `.claude.json` 中有 API 金鑰，則為 `no limits (API)`。
4. 過期的存取權杖不會被使用。沒有 `--refresh` 時，帳號會顯示已儲存的數字，並標示為已過期。
5. 否則 ccam 會呼叫端點。401 或 403 視為登入過期，其他任何狀態或無法讀取的回應內容視為錯誤，兩者都會改用已儲存的數字。
6. 正常的回應會連同測量時間寫入 `<name>.json`。

`fallback()` 會讀取 `<name>.json` 和 `<name>.statusline.json`，並保留較新的那一個。在顯示已儲存的數字之前，它會檢查每個時間區間的重設時間：在數字儲存之後已經重設的時間區間會顯示 0%，因為舊的百分比已不再成立。這就是上面的 `personal` 顯示 `0%` 且沒有重設時間的原因。

端點有兩種回應格式。較新的格式有一個 `limits` 陣列，包含 `session`（5 小時）、`weekly_all` 和 `weekly_scoped`（各模型）項目。較舊的格式有 `five_hour`、`seven_day`、`seven_day_opus` 和 `seven_day_sonnet` 物件。`ParseLive` 兩種都能讀取：以 `limits` 陣列為準，舊物件補上它沒有涵蓋的部分。完全沒有限制資料的回應內容會被回報為格式已變更，`ccam doctor` 會告訴你這件事。

### 更新權杖

更新權杖會輪替：更新一次會得到新的更新權杖，並讓舊的失效。如果 ccam 在 Claude Code 以同一個帳號執行時更新了登入，那個工作階段就會被登出。所以 `ccam usage` 只讀取權杖、絕不寫入，除非你加上 `--refresh`。這時 `Fetcher.refresh` 會：

1. 拒絕以 `ccam import` 就地登錄的帳號，因為直接執行的 `claude` 可能在 ccam 不知情的情況下使用那份登入；
2. 拒絕存放在 macOS 鑰匙圈中的登入資訊；
3. 取得鎖檔 `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock`（該變數未設定時位於 `~/.local/state/ccam/locks/` 下），讓兩個 ccam 行程絕不會重複更新同一個權杖；
4. 在持有鎖的情況下，如果該帳號有工作階段正在執行就拒絕，如果另一個行程已經更新了權杖就略過請求；
5. 把更新權杖送到權杖端點，並以 `fsutil.AtomicWrite` 把新的權杖寫回 `.credentials.json`，保留檔案中的其他所有鍵。

### 狀態列收集器

Claude Code 會執行狀態列指令，並透過管線傳給它一份 JSON 文件。ccam 從中讀取模型名稱和 `rate_limits`：

<!-- motion:usage-statusline -->

1. `runStatusline` 從 `CCAM_PROFILE`（ccam 啟動 Claude Code 時會設定）找出帳號，否則就把 `CLAUDE_CONFIG_DIR` 和各帳號資料夾比對。
2. `ParseStatusLine` 讀取 `rate_limits.five_hour` 和 `rate_limits.seven_day`。API 金鑰帳號的 JSON 中沒有頻率限制，工作階段的第一個回應之前也沒有；這時不會儲存任何東西。
3. `SaveStatusLine` 寫入 `~/.local/state/ccam/usage/<name>.statusline.json`。端點沒有回應時，`ccam usage` 會使用它。
4. ccam 印出 `[name] · model · 5h N% · 7d N%`。加上 `--wrap` 時，它會以相同的輸入透過 `sh -c` 執行你的指令，改為印出 `[name]` 再接上該指令的第一行。在 Windows 上，它會用和 Claude Code 相同的方式執行指令：有安裝 Git Bash 時使用 Git Bash，否則使用 PowerShell。
5. 當 5 小時或每週用量達到 90% 以上，而且還有其他帳號時，這一行的結尾會加上 `→ ccam next --auto --continue`，例如 `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`。

被包裝的指令有 2 秒時間。超過之後會被終止，該行只顯示帳號名稱，所以緩慢的指令無法卡住 Claude Code 的狀態列。

`ccam statusline install` 會把 `statusLine` 寫入每個帳號的 `settings.json`。如果帳號已經有狀態列指令，ccam 會以 `--wrap '<old command>'` 保留它。在 Windows 上，ccam 的路徑會以正斜線和代表主資料夾的 `~` 寫入（`~/.local/bin/ccam.exe statusline`），Git Bash 和 PowerShell 都能理解。透過群組共用的 `settings.json` 只會經由它的連結寫入一次，群組中的每個帳號都會取得這項設定。再次執行時，對已經使用 ccam 的檔案會印出 `(already set)`。

## 安全性和特殊情況

- 除非你加上 `--refresh`，否則 `ccam usage` 絕不變更權杖檔案。登入過期時，用 `ccam <name>` 啟動該帳號一次，Claude Code 就會更新它。
- `--refresh` 絕不碰正在執行的帳號、匯入的帳號或鑰匙圈中的登入資訊。被拒絕的更新會顯示已儲存的數字和原因。
- 用量端點沒有公開的文件，可能會變更。變更時你仍會看到已儲存的數字，`ccam doctor` 也會回報格式的變更。
- 既沒有已儲存的數字也沒有回應時，該列會顯示 `lookup failed:` 和原因，例如加上 `--source cache` 時的 `no saved usage`。
- 快取檔案只存放百分比和重設時間，絕不存放權杖。

## 相關內容

- [用量限制與切換帳號](/docs/usage-limits)
- [切換到下一個帳號](/features/next)
- [用量限制的工作流程](/guides/usage-limit-workflow)
- [指令參考](/docs/commands)
