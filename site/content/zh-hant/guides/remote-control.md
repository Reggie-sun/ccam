---
title: 透過 ccam 使用 Remote Control
description: 用 ccam <name> rc 以任一帳號啟動 Claude Code 的 Remote Control：哪些選項適用、ccam 為什麼提議複製共用的 MCP 伺服器，以及如何解決問題。
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "5f0b90ea53d6"
---

`ccam work rc` 會以帳號 `work` 啟動 Claude Code 的 Remote Control，讓你可以從 claude.ai/code 或 Claude 行動應用程式操作這台電腦上的工作階段。它的用法和 `claude rc` 一樣，只有一個差別：如果該帳號共用 MCP 伺服器，ccam 會詢問一次是否要把它們複製到帳號的 `.claude.json`，因為 Remote Control 的工作階段不會從其他任何地方讀取伺服器。

## 啟動

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[快捷指令](/features/shortcuts)的用法相同：`ccwork rc` 就是 `ccam work rc`。在你要工作的資料夾中執行指令，然後開啟 claude.ai/code 或 Claude 應用程式，並用與 `work` 相同的 Claude 帳號登入。

## 選項

這裡會用到兩類選項：ccam 自己的選項，寫在帳號名稱之前或放在 ccam 的設定中；以及 Remote Control 的選項，寫在 `rc` 之後，原封不動地傳給 Claude Code。

### ccam 選項和設定

| 選項或設定 | 作用 | 預設值 |
| --- | --- | --- |
| `--dry-run`（寫在帳號之前） | 印出環境變數和指令，有 MCP 提示時也一併印出。不會詢問，也不會複製 | 關閉 |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | 當帳號共用 MCP 伺服器，而它的 `.claude.json` 缺少群組目前的清單時詢問。`y` 會複製，`N` 或 Enter 則不複製直接啟動 | `N` |
| `~/.config/ccam/config.toml` 中的 `mcp_strategy = "sync"` | 每次啟動前都把群組的伺服器複製到 `.claude.json`，不論是否為 Remote Control，因此永遠不會詢問 | `"flag"` |
| `ccam group leave <group> work` | 把帳號移出與它共用 MCP 的群組；之後 Remote Control 只會帶著帳號自己的伺服器啟動，也不會詢問 | 在 `default` 中 |

### Remote Control 選項

`rc` 之後的所有內容都交給 Claude Code。下表是 Claude Code 2.1.304 中 `claude remote-control` 的選項；`ccam work rc --help` 會列出你所用版本的選項。

| 選項 | 作用 |
| --- | --- |
| `--name <name>` | 在 claude.ai/code 中顯示的工作階段名稱 |
| `--remote-control-session-name-prefix <prefix>` | 自動產生的工作階段名稱前綴（預設：主機名稱） |
| `-c`, `--continue` | 重新連上 Remote Control 為此資料夾最後記錄的工作階段，約 4 小時內有效 |
| `--session-id <id>` | 重新連上指定的工作階段 |
| `--permission-mode <mode>` | 新工作階段的權限模式：`acceptEdits`、`auto`、`bypassPermissions`、`default`、`dontAsk` 或 `plan` |
| `--spawn <mode>` | `same-dir`（預設）、`worktree`（每個工作階段一個 git worktree）或 `session`（只有一個工作階段，結束時即退出） |
| `--capacity <N>` | `same-dir` 和 `worktree` 模式下可同時執行的工作階段上限（預設 32） |
| `--[no-]create-session-in-dir` | 啟動時在目前資料夾建立第一個工作階段（預設開啟） |
| `--[no-]chrome` | 在新工作階段中使用 Claude in Chrome（預設關閉） |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | 除錯輸出和詳細輸出 |

請把 `rc` 緊接在帳號名稱之後。ccam 會把第一個不以 `-` 開頭的參數當成 Claude Code 的指令，所以在 `ccam work rc --name laptop` 中它看到的是 `rc`，但在 `ccam work --model opus rc` 中看到的是 `opus`。

## ccam 做了什麼

啟動會走平常的流程（[啟動 Claude Code](/features/launch)），只有 MCP 這一步不同：

1. 和任何一次啟動一樣，ccam 會找到帳號、檢查它的共用連結，並解析 `claude` 的位置。
2. 它發現指令是 `rc` 或 `remote-control`，於是不加上 `--mcp-config`。
3. 如果帳號共用 MCP 伺服器，它會比較群組的 `mcp.json` 和 `~/.local/state/ccam/mcp/work/` 中的 `sync-stamp`。如果 `.claude.json` 已經含有目前的伺服器，或群組沒有伺服器，就直接跳到第 6 步。
4. 否則，它會列出工作階段將缺少的伺服器：

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. 回答 `y` 時，它會先備份 `.claude.json`、取得 Claude Code 自己的鎖、把群組的伺服器加入 `mcpServers`，並保留帳號的其他伺服器。回答 `N`、沒有終端機，或帳號正在執行時，它不會做任何變更，並會顯示提示：

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. 它會記錄工作階段，把 `CLAUDE_CONFIG_DIR` 設為帳號資料夾，然後透過 `exec` 把自己替換成 `claude rc`。

之後 Remote Control 會以伺服器的形式持續執行。它為你啟動的每個工作階段都是這個程序的子程序，因此都以 `work` 的身分執行：使用 `work` 的登入資訊、設定和對話記錄，並計入 `work` 的用量限制。伺服器執行期間，ccam 會把該帳號視為正在執行。

試執行可以看出差別。一般啟動會把群組的伺服器當成參數傳入，Remote Control 則不會：

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## 為什麼這樣運作

**Remote Control 不接受 `--mcp-config`。** 一般啟動時，ccam 會以 `--mcp-config=<path>` 的形式把群組的伺服器交給 Claude Code，而且從不寫入 Claude Code 的檔案。`claude remote-control` 會拒絕出現在指令名稱之前的這個選項，在指令名稱之後也不接受，所以在 ccam 0.6.0 之前，只要帳號共用 MCP 伺服器，`ccam work rc` 就會失敗。現在 ccam 不再加上這個選項。

**它的工作階段只從 `.claude.json` 讀取伺服器。** Remote Control 啟動的工作階段不會繼承命令列上的選項。它們找到使用者 MCP 伺服器的唯一位置是帳號的 `.claude.json`。這個檔案也存放帳號的登入資訊和各資料夾的狀態，所以 ccam 無法像處理 `settings.json` 那樣把它連結到共用的副本。要讓工作階段拿到共用清單，唯一的辦法就是把伺服器複製進這個檔案。

**ccam 寫入前會先詢問。** 在預設的 flag 策略下，ccam 承諾不更動 `.claude.json`。Remote Control 是唯一必須寫入該檔案的情況，所以它會先詢問、保留備份，而且只更動 `mcpServers`。

**帳號執行中時不寫入。** Claude Code 會在工作階段期間重寫 `.claude.json`，從外部做的變更會遺失或被覆寫。因此只要該帳號有任何工作階段正在執行，包括另一個 Remote Control 伺服器，ccam 就會略過複製。

**每次變更只問一次。** 複製之後，`sync-stamp` 會記錄群組 `mcp.json` 的修改時間。只有在該檔案變更時（例如執行 `ccam mcp add` 之後），ccam 才會再次詢問。如果你回答 `N`，因為還沒有副本，每次啟動 Remote Control 都會詢問。

**`ccam mcp` 會讓副本保持最新。** 名稱衝突時，`.claude.json` 中的項目優先於 `--mcp-config`，所以舊的副本在一般啟動中也會蓋過較新的定義。因此 `ccam mcp add`、`add-json`、`remove` 和 `import` 會更新所有未在執行的帳號中的副本：

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**沒有終端機時，答案就是「否」。** Remote Control 常在腳本、服務或 `nohup` 中執行，沒有人能回答。這時 ccam 不會等待，而是不帶這些伺服器直接啟動。

## 問題與解決方法

### Remote Control 工作階段中缺少共用的 MCP 伺服器

原因是詢問時回答了 `N`、當時沒有終端機，或帳號正在執行。停止 Remote Control，關閉該帳號的其他工作階段，然後在終端機中重新執行 `ccam work rc`，並回答 `y`。如果希望永遠不被詢問，請在 `~/.config/ccam/config.toml` 中設定 `mcp_strategy = "sync"`。已開啟的工作階段中的伺服器在執行期間不會改變，所以複製之後請啟動新的工作階段。

### "Account work is running, so its .claude.json is left alone"

`work` 仍有工作階段開著：可能是一般的 `ccam work`、快捷指令，或另一個 Remote Control 伺服器。結束它，然後重新執行 `ccam work rc`。其他帳號執行的工作階段不算在內。

### ccam 每次都會詢問

你還沒有複製這些伺服器。回答一次 `y`，或設定 `mcp_strategy = "sync"`。如果你完全不想在 Remote Control 中使用這些伺服器，就停止與該帳號共用 MCP：把它移出那個群組（`ccam group leave default work`），或把 MCP 放在一個該帳號不在其中的專屬群組。

### `ccam work rc` 立刻失敗，並出現與 `--mcp-config` 有關的錯誤

你的 ccam 比 0.6.0 舊，會把 `--mcp-config` 傳給 Remote Control。請執行 `ccam update`。如果 ccam 已是最新版本仍出現同樣的錯誤，表示在 `rc` 之前有一個帶值的選項，例如 `ccam work --model opus rc`，於是 ccam 把 `opus` 當成了指令；請把 `rc` 緊接在帳號名稱之後。

### 用 `ccam mcp` 做的變更沒有反映到 Remote Control

`ccam mcp add` 和 `remove` 只會更新未在執行的帳號中的副本，並會提示：`Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` 停止 Remote Control 後重新啟動；ccam 會發現 `mcp.json` 已變更，並再次詢問。

### 複製取代了我自己的伺服器

如果帳號自己有一個和群組伺服器同名的伺服器，複製會取代它。先前的 `.claude.json` 存放在 `~/.local/state/ccam/mcp/work/backups/`，那裡會保留最新的五份。若想讓所有帳號都使用你的版本，請在下次啟動 Remote Control 之前，用 `ccam mcp import work <server> --force` 把它移入群組。

### 我想移除複製過去的伺服器

`ccam mcp remove <server>` 會把伺服器從群組和每個副本中移除。若要回到複製之前的狀態，請在該帳號沒有任何工作階段執行時，從 `~/.local/state/ccam/mcp/work/backups/` 還原最新的備份。

### Remote Control 詢問是否信任該資料夾，或在沒有終端機時直接退出

資料夾的信任設定是按帳號存放在 `.claude.json` 中的，所以你以 `personal` 身分信任過的資料夾，對 `work` 來說是新的。先在該資料夾中執行一次 `ccam work` 並接受信任，然後再在那裡啟動 Remote Control。

### "You must be logged in with a Claude account that has a subscription"

Remote Control 需要以 Claude 訂閱帳號登入。在 `ccam list` 中查看該帳號的電子郵件，必要時用 `ccam login work` 重新登入。

### 手機上看不到工作階段

Claude 應用程式或 claude.ai/code 登入的是另一個 Claude 帳號。把它和 `ccam list` 為 `work` 顯示的電子郵件比對一下。

### ccam 警告另一個 Claude Code 在此資料夾中共用對話記錄

同一記錄群組中的另一個帳號正在同一個資料夾中執行。如果兩個帳號的工作階段延續同一段對話，記錄可能會混在一起。只有在它們處理不同對話時才回答 `y`，否則請在另一個資料夾中啟動 Remote Control。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [共用的 MCP 伺服器](/features/mcp)
- [帳號之間的共用](/features/share)
- [管理帳號](/features/accounts)
