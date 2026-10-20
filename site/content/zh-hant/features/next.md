---
title: 切換到下一個帳號
description: ccam next 會以你上次使用的帳號的下一個帳號開啟 Claude Code；--auto 選擇剩餘額度最多的帳號，--continue 接續對話，--fallback 會重新執行 -p 工作。
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 7bc55c2dd520
---

`ccam next` 不需要你指定名稱，就能以另一個帳號開啟 Claude Code。當某個帳號達到用量限制時使用它：單純的 `ccam next` 會選擇你上次使用的帳號的下一個帳號，`--auto` 會選擇 5 小時區間內剩餘額度最多的帳號，`--continue` 會在新帳號中接續上一段對話。對於執行 `claude -p` 的指令碼，`ccam --fallback` 會自行完成這個切換（請參閱[達到上限時重新執行](#run-again-at-the-limit)）。

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

ccam 會在標準錯誤上以 `→ <from> → <to>` 印出切換，然後以和 `ccam <name>` 相同的方式用新帳號啟動 Claude Code。`<from>` 是你上次使用的帳號。`--` 後面的參數會原封不動地傳給 `claude`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--auto` | 選擇已登入且 5 小時用量最低的帳號，略過用量為 100% 的帳號 | 關閉 |
| `-c`、`--continue` | 只考慮和你上次使用的帳號共用對話記錄的帳號，並啟動 `claude --continue` | 關閉 |
| `-- <claude args>` | 接在 ccam 自己的參數之後傳給 `claude` 的參數 | 無 |
| `--dry-run` | 支援：印出 `→ <from> → <to>` 和要執行的命令列，不啟動 Claude Code，也不記錄啟動 | 關閉 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 時，`--continue` 會放在你自己的參數之前。

## 運作方式

每次透過 ccam 啟動，都會把時間寫入 `~/.local/state/ccam/last-used/<name>`。`ccam next` 會讀取這些檔案，檢查哪些帳號已登入，然後做出選擇。

1. `cli.candidates` 根據 `config.toml` 為每個帳號建立一個 `next.Candidate`：帳號、它上次啟動的時間（`launcher.LastUsed`），以及它是否有登入資訊（`creds.Has`：可讀取的 `.credentials.json`，或在 macOS 上的鑰匙圈項目）。
2. 沒有 `--auto` 時，`next.Rotate` 會依槽位編號排序帳號，並找出最近一次啟動的帳號。那就是你要切換離開的帳號。如果還沒有任何帳號啟動過，則是槽位最小的帳號。
3. 從那裡開始依槽位順序往前走，從最後一個槽位繞回第一個，並傳回第一個已登入的帳號。
4. `launcher.LaunchWith` 會完全像 `ccam <name>` 一樣準備並啟動所選的帳號：檢查帳號的連結、套用 MCP 共用、登錄工作階段，並透過 `syscall.Exec` 以 `claude` 取代 ccam。

走訪從上次使用的帳號開始，並在它之後的第一個已登入帳號停下：

<!-- motion:next-rotate -->

`work` 是最近一次啟動的帳號，所以走訪從它之後開始。`side` 沒有登入，因此被略過，最後選擇了 `lab`。如果 `lab` 也已登出，走訪就會繞回 `personal`（圖中的 `% len`）。輪替絕不會傳回你出發時的帳號：沒有其他已登入的帳號時，它會以錯誤停止。

### --auto

`--auto` 依用量而不是槽位順序決定。`cli.pickAuto` 保留已登入的帳號，用和 `ccam usage` 相同的 `Fetcher.FetchAll` 一次查詢它們所有的用量，再交給 `usage.Best()` 選擇：

<!-- motion:next-auto -->

`usage.Best()` 選擇 5 小時百分比最低的帳號，並略過用量為 100% 或 5 小時用量未知的帳號。查詢遵循 `ccam usage` 的規則：60 秒內取得的數字會重複使用，每個帳號最多等待 3 秒，沒有回應或登入已過期的帳號則以已儲存的數字判斷。已儲存的時間區間如果已經過了重設時間，會視為 0%，所以在那之後已經重設的帳號不會被跳過。`--auto` 絕不更新權杖。和單純的輪替不同，如果你上次使用的帳號仍然剩餘最多額度，`--auto` 也可能選擇它。

### --continue

Claude Code 的 `--continue` 會開啟目前資料夾中最近的對話，而這只有在新帳號看得到那段對話時才行得通。所以加上 `--continue` 時，ccam 只會考慮和你上次使用的帳號在同一個 `history` 群組中的帳號，然後把 `--continue` 放在你的參數之前。對話記錄為 `local` 的帳號不和任何帳號共用，所以從這樣的帳號執行 `--continue` 一定會失敗。請先用 `ccam share <name> --history on` 開啟共用（請參閱[共用](/docs/sharing)）。

## 達到上限時重新執行

`ccam next` 用於已經停下來的互動式工作階段。非互動式的執行，也就是指令碼或排程工作中的 `claude -p`，無法等你處理，所以在帳號名稱前加上 `--fallback`，就能讓 ccam 替它切換：

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

1. `--fallback` 只能和 `-p` 或 `--print` 一起使用；否則 ccam 會停止，並指向 `ccam next --auto --continue`。
2. 標準輸入只讀取一次，並提供給每次嘗試。
3. ccam 會把 `claude` 作為子行程啟動，而不是使用 `exec`；它會暫存標準輸出，標準錯誤則即時傳遞。
4. 當 `claude` 以 0 結束，或因其他原因失敗時，暫存的輸出和結束代碼會原樣傳回。
5. 當 Claude Code 的訊息表明是用量上限（`You've hit your limit`、`usage limit reached`），或該帳號的用量查詢顯示 5 小時或每週區間已達 100% 時，失敗會被視為達到用量上限。
6. 這時會丟棄該次嘗試的輸出，並以另一個已登入的帳號執行同一個指令：先是 5 小時用量最低的帳號，接著是查不到用量的帳號，依編號順序。已達 100% 的帳號會被略過，每個帳號最多嘗試一次。
7. 參數中有 `-c`、`--continue`、`-r` 或 `--resume` 時，只會嘗試同一個對話記錄群組中的帳號，這樣才找得到那段對話。

沒有帳號可用時，ccam 會印出最後一次嘗試的輸出，說明沒有任何已登入的帳號還有剩餘額度，並以該次嘗試的結束代碼結束。`ccam --dry-run --fallback work -p …` 會印出第一次嘗試的指令。

互動式工作階段不會自動切換。取而代之的是，執行 `ccam statusline install` 後，當 5 小時或每週用量達到 90% 時，狀態列會加上 `→ ccam next --auto --continue`（[用量限制](/features/usage)）。

## 安全性和特殊情況

- 沒有登入資訊的帳號絕不會被選中。沒有其他已登入的帳號時，單純的 `ccam next` 會以 `no logged-in account other than <name>` 停止。
- 加上 `--continue` 但沒有帳號共用對話記錄時：`` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``。
- 加上 `--auto` 但沒有已登入的帳號時：`no logged-in account to pick`。當每個帳號都已用滿或用量未知時：`` no account with known usage and headroom (see `ccam usage`) ``。
- 即使加上 `--dry-run`，`--auto` 仍會發出網路請求。它只會讀取，絕不寫入權杖檔案。
- 所選的帳號會透過和 `ccam <name>` 相同的流程啟動，所以適用相同的檢查，包括另一個共用相同對話記錄的 Claude Code 已在這個資料夾中執行時的警告。
- `last-used` 只在 Claude Code 實際啟動時才寫入。`--dry-run` 不會碰它，所以試執行不會讓輪替往前推進。
- `--fallback` 會暫存標準輸出直到一次嘗試結束，所以 `--output-format stream-json` 這類串流輸出會在最後一次全部送出。第一次嘗試在達到上限前變更的檔案會維持變更後的狀態；下一次嘗試會從頭重新執行提示詞。

## 相關內容

- [用量限制](/features/usage)
- [用量限制與切換帳號](/docs/usage-limits)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [用量限制的工作流程](/guides/usage-limit-workflow)
