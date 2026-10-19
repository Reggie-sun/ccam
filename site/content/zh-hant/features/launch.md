---
title: 啟動 Claude Code
description: 用 ccam <name|number>、ccam run 或 ccam 選擇器以任一帳號啟動 Claude Code；參數原封不動地傳給 claude，--dry-run 會顯示要執行的指令。
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 401eed57e533
---

`ccam work` 會以帳號 `work` 啟動 Claude Code，名稱後面的所有內容都原樣傳給 `claude`。凡是你會輸入 `claude` 的地方都可以使用它：在終端機中、在腳本中，或搭配管線輸入。

## 用法

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

指定帳號的名稱，或它在 `ccam list` 中的編號。ccam 不會讀取名稱後面的參數，所以 `ccam work --help` 顯示的是 Claude Code 的說明。名稱後面的 `--` 可加可不加。標準輸入、結束代碼和訊號的行為都和 `claude` 相同。

當第一個參數以 `-` 開頭，而且不是 ccam 自己的旗標時，ccam 會帶著所有參數啟動預設帳號。預設帳號是用 `ccam use` 設定的帳號；如果你沒有設定，則是編號最小的帳號。

### 選擇器

不帶參數的 `ccam` 在標準輸入和輸出都是終端機時會開啟選擇器，否則會印出 `ccam list`。它會顯示每個帳號的編號、電子郵件、5 小時用量、上次使用時間和共用的項目；用量會在每個帳號回應時陸續填入。

| 按鍵 | 作用 |
| --- | --- |
| `Enter` | 以選取的帳號啟動 Claude Code |
| `1` 到 `9` | 啟動該編號的帳號 |
| `Up`/`Down`、`k`/`j` | 移動選取項目 |
| `a` | 關閉選擇器並執行 `ccam add` 設定精靈 |
| `s` | 編輯選取帳號的共用項目（`Space` 或 `x` 切換，`Enter` 套用，`Esc` 返回） |
| `d` | 把選取的帳號設為預設帳號 |
| `r` | 略過短期快取，重新取得用量 |
| `q`、`Esc`、`Ctrl+C` | 結束 |

選擇器會在執行動作之前關閉，所以在 Claude Code 啟動前終端機就已經還原。按 `d` 或變更共用之後，它會再次開啟。變更共用時會像 `ccam share` 一樣先顯示計畫並詢問你。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--dry-run` | 印出 ccam 會使用的環境變數和指令，然後結束 | 關閉 |
| `-v`、`--version` | 印出 `ccam <version>`（只能作為第一個參數） | |
| `-h`、`--help` | ccam 的說明，可用 `ccam --help` 或 `ccam run --help` | |

參數規則：

- ccam 只讀取帳號名稱前面的自身旗標。名稱後面的每個參數都屬於 Claude Code，包括 `--dry-run` 和 `--help`。
- `--dry-run` 可以寫成 `ccam --dry-run work`、`ccam run --dry-run work` 和 `ccam --dry-run`（預設帳號）。它會被接受，不會被拒絕。
- `ccam run <name|number> [claude args]` 是 `ccam <name>` 的完整寫法。不帶名稱的 `ccam run` 會使用預設帳號。
- 帳號名稱不能是 `add`、`list`、`run` 或 `usage` 這類 ccam 指令名稱，所以 `ccam <word>` 永遠不會有歧義。

試執行會印出一行可以直接貼到 shell 的指令：

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## 運作方式

每次啟動都經過相同的步驟，從命令列一直到 `exec`：

<!-- motion:launch -->

每個方塊會在其步驟執行時亮起，下方的各行是 Claude Code 啟動時使用的環境變數和參數。

1. `cmd/ccam/main.go` 呼叫 `app.Main`。如果 ccam 是以其他名稱啟動的，那就是[快捷指令](/features/shortcuts)，會直接交給啟動器。
2. `app.parse` 在 cobra 介入之前檢查 `os.Args`。只有單獨的 `ccam`（選擇器）、ccam 指令名稱（`config.Reserved`）、`-h`/`--help` 和 `-v`/`--version` 會交給 `cli.Execute`。其他一切都是啟動。
3. `launcher.Launch` 讀取 `~/.config/ccam/config.toml`，這是它唯一解析的檔案。
4. `launcher.Prepare` 依名稱或編號找到帳號，去掉開頭的一個 `--`，並用 `ResolveClaude` 找出 `claude` 執行檔。PATH 查詢的結果會連同檔案的修改時間快取在 `~/.local/state/ccam/claude-path`，所以 ccam 不會在每次啟動時搜尋 `PATH`。檔案變更或不存在時會重新查詢。`config.toml` 中的 `claude_bin` 路徑會直接使用。如果 `claude` 就是 ccam 本身（以 inode 比對），ccam 會拒絕，因為那樣會無限循環。
5. `checkLinks` 用 `share.Check` 把帳號資料夾 `~/.local/share/ccam/profiles/<name>/` 和它的共用設定比對，過程中只呼叫 `lstat` 和 `readlink`，每個項目一到兩次系統呼叫。缺少的連結會當場建立。任何可能遺失資料的情況都會中止啟動。
6. `applyMCP` 加入共用的 MCP 伺服器（見下文）。
7. `session.Register` 記錄工作階段。ccam 會在 `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`（沒有 `XDG_RUNTIME_DIR` 時為 `~/.local/state/ccam/locks/`）寫入一個以其 PID 命名、記錄行程啟動時間的標記。對於共用對話記錄的帳號，`historyCheck` 會先在 `history-<group>-<hash>`（hash 取自目前的資料夾）中尋找正在執行的工作階段，並同樣在那裡登錄。`exec` 會保留 PID，所以標記之後指的就是 Claude Code 行程本身。使用標記而不是持有鎖，是因為 Claude Code 啟動的行程會繼承鎖，並在工作階段結束後繼續持有它。
8. `Plan.Exec` 把終端機標題設為 `claude · <name>`（當標準錯誤是終端機時），在不 `fsync` 的情況下把啟動時間寫入 `~/.local/state/ccam/last-used/<name>`，然後呼叫 `syscall.Exec`。

Claude Code 執行時，`CLAUDE_CONFIG_DIR` 會設為帳號資料夾，`CCAM_PROFILE` 會設為帳號名稱。ccam 會移除這些變數繼承而來的任何值。在 macOS 上，對於用 `ccam import` 登錄的 `~/.claude` 資料夾，它還會把 `CLAUDE_SECURESTORAGE_CONFIG_DIR` 設為空值，讓該帳號和直接執行的 `claude` 使用同一個鑰匙圈登入。

### 為什麼用 exec

`syscall.Exec` 會以 `claude` 取代 ccam 行程，所以 ccam 不會留在記憶體中，按鍵、訊號和結束代碼都直接屬於 Claude Code，參數也始終是陣列，絕不會變成 shell 字串。工作階段之後的工作，例如發現用 `claude mcp add` 新增的伺服器，會在下次啟動時進行：ccam 會把 `.claude.json` 的修改時間和 `~/.local/state/ccam/mcp/<name>/seen` 比對，只在檔案變更時才解析它。

### 命令列如何分派

三種命令列，在 `internal/app` 中走三條路徑：

<!-- motion:launch-route -->

快捷指令依 `argv[0]` 比對，`ccam <name>` 依第一個參數比對，兩者最後都是 `syscall.Exec`。只有 `ccam usage` 這類 ccam 指令才會進入 cobra。

### 共用的 MCP 伺服器

使用預設的 `flag` 策略時，ccam 會以參數形式把群組的 MCP 伺服器傳給 Claude Code，放在你的參數之前：

<!-- motion:launch-mcp -->

`mcp.FlagArgs` 只在群組的 `mcp.json` 存在時才傳回參數，`Plan.InsertArgs` 會把它放在 `claude` 的正後方。你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍會疊加套用。它是單一參數 `--mcp-config=<path>`，因為 Claude Code 的 `--mcp-config` 接受多個值：如果寫成兩個參數，後面的提示詞也會被當成設定檔讀取。

在 `config.toml` 中設定 `mcp_strategy = "sync"` 時，ccam 會改為把伺服器寫入帳號的 `.claude.json`；如果帳號已經在執行，就略過這一步。

### Remote Control

`ccam work rc`（或 `ccam work remote-control`）會以帳號 `work` 啟動 Claude Code 的 Remote Control。`claude remote-control` 會拒絕出現在動詞之前的 `--mcp-config`，在動詞之後也不接受，而且它的工作階段只從帳號的 `.claude.json` 讀取 MCP 伺服器。因此 flag 策略在這裡的運作方式不同：

<!-- motion:launch-rc -->

1. `remoteControl` 會把帳號之後第一個不以 `-` 開頭的參數當成 Claude Code 的動詞。當它是 `remote-control` 或 `rc` 時，`mcp.FlagArgs` 不會加入任何東西。
2. `mcp.Synced` 會比較群組 `mcp.json` 的修改時間和帳號的 `sync-stamp`。如果 `.claude.json` 已經含有群組目前的伺服器，或群組沒有伺服器，Remote Control 會直接啟動，不會詢問。
3. 否則，ccam 會列出工作階段將缺少的伺服器並詢問：

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. 回答 `y` 時，`mcp.Sync` 會用 `sync` 策略的方式複製它們：取得 Claude Code 的鎖、先建立備份，並保留帳號自己的伺服器。回答 `N` 或沒有終端機時，Remote Control 會在沒有這些伺服器的情況下啟動，ccam 會提示 `mcp_strategy = "sync"`，它會在每次啟動前複製這些伺服器。

在伺服器被複製之前，ccam 每次啟動 Remote Control 都會詢問；複製之後，只有在群組的 `mcp.json` 再次變更時才會詢問一次。在之後的一般啟動中，`.claude.json` 中的副本優先於 `--mcp-config`，所以 `ccam mcp add`、`remove` 和 `import` 會更新其他帳號持有的副本（[共用的 MCP 伺服器](/features/mcp)）。

### 10 ms 的預算

ccam 在 Claude Code 啟動前最多增加 10 ms（p95）。啟動流程絕不建立 cobra 指令樹、不啟動 Bubble Tea，也不開啟網路連線；如果套件的 `init()` 函式總共花費超過 5 ms，`TestPackageInitIsCheap` 就會失敗。`scripts/perf-check.sh` 用 `hyperfine` 測量 `ccam --dry-run work`、`ccam work` 和一個快捷指令，超過 10 ms 就失敗；`internal/launcher` 中的 `BenchmarkLaunchPath` 則會以迴圈執行這些步驟。若要查看每個步驟花費的時間，請設定 `CCAM_TRACE=1`：

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 安全性和特殊情況

- 未知的名稱或編號會以 `no account "<name>"` 停止，並提示你查看 `ccam list`。完全沒有帳號時，ccam 會請你用 `ccam add <name>` 新增一個。
- 如果某個共用路徑變成了一般檔案或指向錯誤的位置，ccam 會以 `stopped because of share link problems:` 停止，列出這些問題，並請你執行 `ccam doctor --fix`。它不會在啟動流程中移動資料。
- 如果另一個和這個帳號共用對話記錄（同一個群組）的 Claude Code 正在同一個資料夾中執行，ccam 會警告對話記錄可能混在一起，並詢問 `Run anyway? [y/N]`。沒有終端機時，它會顯示警告並繼續。
- `--dry-run` 會略過 MCP 同步、新伺服器通知、對話記錄檢查、工作階段標記、標題和上次使用時間，也不會 exec。它仍會檢查共用連結，並像一般啟動一樣建立缺少的連結。
- 如果在 `work` 已經在執行時執行 `ccam work rc`，ccam 不會碰 `.claude.json`，會告訴你這一點，並在沒有群組伺服器的情況下啟動 Remote Control；工作階段結束後再執行一次即可複製它們。使用 `--dry-run` 時，它只會印出訊息，既不詢問也不複製。
- 使用 `mcp_strategy = "sync"` 時，同步失敗會印出 `ccam: MCP sync failed (launching anyway): …`，Claude Code 會以它原本就有的伺服器啟動。

## 相關內容

- [快捷指令](/features/shortcuts)
- [透過 ccam 使用 Remote Control](/guides/remote-control)
- [管理帳號](/features/accounts)
- [指令參考](/docs/commands)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [CLAUDE_CONFIG_DIR 的運作方式](/guides/claude-config-dir)
