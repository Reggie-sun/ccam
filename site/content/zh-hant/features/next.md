---
title: 切換到下一個帳號
description: ccam next 會以另一個帳號開啟 Claude Code；--auto 選擇剩餘額度最多的帳號，--continue 把對話帶過去，--fallback 會重新執行 -p 工作。
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` 不需要你指定名稱，就能以另一個帳號開啟 Claude Code。當某個帳號達到用量限制時使用它：單純的 `ccam next` 會選擇你上次使用的帳號的下一個帳號，`--auto` 會選擇剩餘額度最多的帳號，`--continue` 則會把這個資料夾的上一段對話帶到新帳號，無論兩個帳號是否共用對話記錄。在 Claude Code 中，[ccam 外掛](#inside-claude-code)會在用量上限讓工作階段停止時通知你，`/ccam:switch` 則會把對話交接出去。對於執行 `claude -p` 的指令碼，`ccam --fallback` 會自行完成切換（請參閱[達到上限時重新執行](#run-again-at-the-limit)）。

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
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam 會在標準錯誤上以 `→ <from> → <to>` 印出切換，`--auto` 選擇時還會附上理由，然後以和 `ccam <name>` 相同的方式用新帳號啟動 Claude Code。`<from>` 是你上次使用的帳號；加上 `--continue` 時，則是這段對話所屬的帳號。`--` 後面的參數會原封不動地傳給 `claude`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--auto` | 依用量（見下文）選擇已登入的帳號，略過用量為 100% 的帳號 | 關閉 |
| `--max-usage <1-100>` | 搭配 `--auto`：優先選擇用量低於此值的帳號 | `auto_max_usage`，否則為 100 |
| `--strategy headroom\|expiring\|balance` | 搭配 `--auto`：5 小時用量最低、5 小時區間最早重設，或每週額度最先會用不完的帳號 | `auto_strategy`，否則為 `headroom` |
| `--cross-org` | 搭配 `--auto`：也選擇其他組織的帳號（見下文） | `cross_org`，否則為關閉 |
| `--wait[=<duration>]` | 搭配 `--auto`：沒有帳號還有額度時，等待第一個帳號重設，最多等這麼久 | 關閉；單獨指定時為 `5h` |
| `-c`、`--continue` | 把這個資料夾的上一段對話帶到新帳號接續（見下文） | 關閉 |
| `-- <claude args>` | 接在 ccam 自己的參數之後傳給 `claude` 的參數 | 無 |
| `--dry-run` | 支援：印出 `→ <from> → <to>` 和要執行的命令列，不啟動 Claude Code，也不記錄啟動 | 關閉 |

`--auto` 和 `--continue` 可以一起使用。使用 `--continue` 時，`--resume` 參數會放在你自己的參數之前。

Claude Code 自己的選項，例如 `--dangerously-skip-permissions` 或 `--model`，要放在 `--` 後面。`ccam <name>` 會把名稱後面的所有內容原封不動地傳下去，`ccam next` 則會自己讀取參數，所以少了 `--` 時會以 `ccam: unknown flag: --dangerously-skip-permissions` 停止。

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

`--auto` 依用量而不是槽位順序決定。`cli.pickAuto` 保留已登入的帳號，用和 `ccam usage` 相同的 `Fetcher.FetchAll` 一次查詢它們所有的用量，再交給 `usage.Policy.Choose()` 選擇：

<!-- motion:next-auto -->

帳號依它的 5 小時和每週區間判斷；當 Claude Code 的參數中（`--` 後面）有 `--model` 時，也會依該模型的每週區間判斷（名稱中含有 `opus`、`sonnet`、`haiku` 或 `fable`）。`Policy.Choose()` 把帳號分成三類：

1. **有餘裕**：每個區間都低於 `auto_max_usage`（`config.toml`，預設 100）。策略會在其中選擇：`headroom`（預設）選 5 小時用量最低的帳號，`expiring` 選 5 小時區間最早重設的帳號，讓它剩下的額度在失效前用掉；`balance` 選在每週區間重設前，平均每天剩餘每週額度最多的帳號，並把每個執行中的 ccam 工作階段計入該帳號，讓每週額度不會閒置，平行的工作階段也會分散開來。
2. **超過**：低於 100%，但達到或超過 `auto_max_usage`。只有在沒有帳號有餘裕時才會使用，先選 5 小時用量最低的帳號。
3. **已滿**：有區間達到 100%，或 5 小時用量未知。絕不會被選中。

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**組織。** 自動選擇只會在和你要離開的帳號同類的帳號之間進行：同一個團隊或企業組織，或者個人方案（Pro、Max）彼此之間，所以公司帳號的工作不會在私人帳號上繼續，反之亦然。ccam 從每個帳號的登入得知它的組織（`ccam login` 會記錄；較舊的登入則從 `.claude.json` 讀取）；無法判斷的帳號可以和任何帳號搭配。登入同一個 Claude 帳號的兩個資料夾共用它的用量上限，所以其中一個絕不會接替另一個。`--cross-org`，或在 `config.toml` 中設定 `cross_org = true`，會解除組織規則。

無法選出帳號時，ccam 會以代碼 75（`EX_TEMPFAIL`）結束，讓指令碼能分辨「稍後再試」和失敗。加上 `--wait` 時，它會改為等到第一個已滿帳號的阻擋區間重設（重設時間後 30 秒，屆時端點會顯示新的區間），印出 `ccam: no account has room; waiting until 15:10`，然後再選擇一次；沒有任何帳號說明何時重設時，它會在 10 分鐘後再查看。如果等待結束前都沒有帳號空出來，它會以 75 結束。查詢遵循 `ccam usage` 的規則：60 秒內取得的數字會重複使用，每個帳號最多等待 3 秒，沒有回應或登入已過期的帳號則以已儲存的數字判斷。已儲存的時間區間如果已經過了重設時間，會視為 0%，所以在那之後已經重設的帳號不會被跳過。`--auto` 絕不更新權杖。和單純的輪替不同，如果你上次使用的帳號仍然剩餘最多額度，`--auto` 也可能選擇它。

### --continue

1. ccam 會在每個帳號中尋找這個資料夾的上一段對話：`projects/<folder>/` 中最新的 `<session>.jsonl`（資料夾名稱是把工作資料夾中字母和數字以外的每個字元都換成 `-`，和 Claude Code 的命名方式相同）。如果工作階段留下了記錄（用量上限或 `/ccam:switch`，見下文），而且之後沒有寫入更新的內容，就改用那個工作階段。
2. 它所屬的帳號就是你要切換離開的帳號，新帳號則依上述方式在同一個組織內選出。
3. 如果新帳號共用那份對話記錄（同一個群組資料夾），ccam 會啟動 `claude --resume <session id>`，也就是那段對話本身（`--continue` 則會開啟資料夾中最新的那段對話，不論是哪一段）。
4. 否則它會啟動 `claude --resume <the conversation's file> --fork-session`：Claude Code 會以新的工作階段 ID，把整段對話複製到新帳號自己的對話記錄中，原本的對話維持不變。之後新帳號的 `--continue` 就找得到它。

ccam 會告訴你這段對話有多大：新帳號在第一次回覆時會重新讀取全部內容。提示快取絕不會在組織之間共用，所以加上 `--cross-org` 時，第一次回覆會以完整大小計算。在群組中共用 `history` 可以避免複製、維持同一個對話，也會共用 `/rewind` 用來還原程式碼的檢查點（`file-history/`）。

## 在 Claude Code 中

[ccam 外掛](/docs/install#claude-code-plugin)會新增兩個 hook（Claude Code 2.1.78 或更新版本）：

- 當某一輪因用量上限而停止時（`You've hit your session limit`、每週上限或模型上限），它會記下這段對話，並在支援的終端機（iTerm2、WezTerm、Windows Terminal、kitty、Ghostty）中顯示桌面通知：`ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` 短暫的伺服器節流或 API 金鑰的速率限制不是用量上限，不會處理。Claude Code 2.1.234 及更新版本也會在工作階段中等待，並在 5 小時區間重設後自行繼續；通知讓你可以選擇改為切換帳號。
- `/ccam:switch [account]` 會記下這段對話，並在發出任何請求之前結束這一輪：輸入 `/exit`，接著 `ccam next --continue` 就會以你指定的帳號或下一個帳號接續它。

記錄存放在 `~/.local/state/ccam/handoff/<account>.json`，保留 12 小時，而且只會在下一個工作階段即將啟動時使用一次；啟動失敗時會保留下來，供下次嘗試。它屬於工作階段啟動時所在的資料夾，即使之後在 Claude Code 的 shell 中執行了 `cd` 也一樣。如果它指名的帳號後來被重新命名或移除，ccam 會告知這件事，並照常選擇。在 `config.toml` 中設定 `relay = true` 時，ccam 會留在互動式工作階段旁邊，而不是以 Claude Code 取代自己，並在 `/ccam:switch` 之後自行接續那個工作階段的對話（其他工作階段留下的記錄不會被處理）。這樣每個工作階段都要多一個等待中的 ccam 行程，所以預設為關閉。接力的工作階段會以 `ccam next --continue` 啟動，不帶第一個工作階段啟動時使用的參數，例如 `--model`。

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
8. 和 `--auto` 一樣，只會嘗試同一個組織的帳號（`cross_org = true` 會解除這個限制）。

下一個帳號依和 `--auto` 相同的規則選擇（`auto_max_usage`、`auto_strategy`、`--model` 的區間）。沒有帳號可用時，ccam 會印出最後一次嘗試的輸出，說明沒有任何已登入的帳號還有剩餘額度，並以代碼 75 結束。在帳號名稱前加上 `--wait`（`ccam --fallback --wait=2h work -p …`）時，它會改為等待第一個帳號重設，然後可能再次嘗試每個帳號；它一定會先等待至少一次，所以剛達到上限的帳號不會立刻再次執行。`ccam --dry-run --fallback work -p …` 會印出第一次嘗試的指令。

互動式工作階段不會自動切換。取而代之的是，執行 `ccam statusline install` 後，當 5 小時或每週用量達到 90% 時，狀態列會加上 `→ ccam next --auto --continue`（[用量限制](/features/usage)）。

## 安全性和特殊情況

- 沒有登入資訊的帳號絕不會被選中。沒有其他已登入的帳號時，單純的 `ccam next` 會以 `no logged-in account other than <name>` 停止。
- 加上 `--continue` 但資料夾中沒有對話時：`no account has a conversation in <folder> to continue`。
- 組織規則讓所有帳號都被排除時：`` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``。
- 加上 `--auto` 但沒有已登入的帳號時：`no logged-in account to pick`。當每個帳號都已用滿或用量未知時：`` no account with known usage and headroom (see `ccam usage`) ``。兩者都以結束代碼 75 結束。
- `--max-usage`、`--strategy`、`--wait` 和 `--cross-org` 只能搭配 `--auto` 使用，超出 1-100 的值或其他策略會被拒絕。
- 分支只會複製對話，不會複製帳號的其他狀態：除非共用對話記錄，否則原帳號 `file-history/` 中的檢查點不會一起帶過去。
- 即使加上 `--dry-run`，`--auto` 仍會發出網路請求。它只會讀取，絕不寫入權杖檔案。
- 所選的帳號會透過和 `ccam <name>` 相同的流程啟動，所以適用相同的檢查，包括另一個共用相同對話記錄的 Claude Code 已在這個資料夾中執行時的警告。
- `last-used` 只在 Claude Code 實際啟動時才寫入。`--dry-run` 不會碰它，所以試執行不會讓輪替往前推進。
- `--fallback` 會暫存標準輸出直到一次嘗試結束，所以 `--output-format stream-json` 這類串流輸出會在最後一次全部送出。第一次嘗試在達到上限前變更的檔案會維持變更後的狀態；下一次嘗試會從頭重新執行提示詞。

## 相關內容

- [用量限制](/features/usage)
- [用量限制與切換帳號](/docs/usage-limits)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [用量限制的工作流程](/guides/usage-limit-workflow)
