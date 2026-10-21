---
title: 檢查和修復
description: 用 ccam doctor 檢查共用連結、檔案權限、登入和 Claude Code 版本，並用 ccam doctor --fix 修復可以安全修復的問題。
updated: 2026-10-21
order: 9
translation: machine
sourceHash: "9be63350cb69"
---

`ccam doctor` 會檢查你的 ccam 安裝：Claude Code 執行檔、每個帳號的資料夾、登入和共用連結、共用的 MCP 伺服器、用量端點，以及你的快捷指令。`ccam doctor --fix` 會修復不會遺失資料的問題，並先把要取代的檔案移到備份資料夾。當 `ccam <name>` 因共用連結問題而停止、你手動還原檔案之後，或 Claude Code 更新之後，請執行它。

## 用法

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

第一行一次標出所有帳號，讓你在閱讀其他內容之前就知道問題在哪裡。下方每個帳號各有一個方框，不屬於任何帳號的檢查結果（Claude Code、環境變數、外掛、MCP 伺服器、快捷指令、資料夾規則、群組和設定）則放在 `General` 方框中。方框標題會顯示該帳號最嚴重的標記和問題數量。方框內每一列是一個檢查過的項目，一律依序為 `Folder`、`Login` 和 `Share links`，然後是其他項目，前面有一個標記：`✔` 正常、`ℹ` 資訊、`⚠` 警告、`✘` 錯誤。沒有任何問題的資料夾會顯示 `✔ Folder  OK`。

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

加上 `--fix` 時，每個已修復的項目結尾會加上 `→ fixed`，而且不再算作問題：

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

在終端機中，方框會配合終端機的寬度，過長的訊息會換行，並使用終端機本身的色彩（`NO_COLOR` 可關閉色彩）。輸出到管道或檔案時，沒有色彩，任何訊息都不會換行或截斷，所以 `ccam doctor | grep settings.json` 能找到完整的一行；此時每個方框都與最長的一行同寬。儀表板的 Health 分頁也會顯示同樣的方框。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--fix` | 修復可以安全修復的問題（變更的檔案會先備份） | 關閉 |
| `--dry-run`（全域） | 執行所有檢查，但不修復任何東西，即使加上 `--fix` 也一樣 | 關閉 |

`ccam doctor` 不接受參數。它支援 `--dry-run`：`ccam --dry-run doctor --fix` 會印出和 `ccam doctor` 相同的檢查結果，且不變更任何東西。

只要還有至少一個 `✘` 項目未修復，結束代碼就是 `1`，否則為 `0`。警告和資訊絕不會讓執行失敗。沒有加 `--fix` 而執行失敗時，ccam 最後會印出一行空白和 `` `ccam doctor --fix` repairs what it can. ``。無法載入設定時，也會以 `1` 結束並顯示 `ccam:` 錯誤。

## 運作方式

`internal/cli/doctor.go` 會載入 `~/.config/ccam/config.toml`，並呼叫 `internal/doctor/doctor.go` 中的 `doctor.Run`；它以固定順序執行各項檢查並傳回結果，每個結果都附上所屬的帳號和檢查的項目。`tui.HealthReport` 會把它們分組到各個方框。`--fix` 會以 `Options.Fix` 傳入，在試執行時設為 false。

<!-- motion:doctor-checks -->

中間的區塊會對每個帳號執行一次。各項檢查以及 `--fix` 對它們的處理如下：

1. **Claude Code 版本**（`checkClaude`）。以和啟動器相同的方式找出 `claude` 執行檔，並執行 `claude --version`。如果版本不是以 `2.1.302` 開頭（這是 ccam 驗證資料夾結構時所用的版本，`doctor.LayoutVersion`），會加上 資訊 `ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`。找不到執行檔或 `--version` 失敗則為錯誤。不修復。
2. **環境變數中的登入資訊**（`checkAuthEnv`）。設定了 `CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN` 時發出警告：從該 shell 啟動的每個帳號都會繼承它，而且 Claude Code 會優先使用它，而不是帳號自己的登入。不修復；請自行取消設定。
3. **資料夾**（`checkProfile`）。帳號資料夾不存在時為錯誤，並結束該帳號的檢查。權限中設定了任何群組或其他使用者的位元時為警告；`--fix` 會執行 `chmod 0700`。Windows 沒有權限位元，所以在那裡會略過這項檢查和下方的 `0600` 檢查；位於使用者設定檔中的資料夾由其存取清單保護。
4. **登入檔案**（`checkCredentials`）。以 `lstat` 檢查 `.credentials.json`。符號連結是錯誤，因為每個帳號都必須有自己的權杖，而且不會修復。權限不是 `0600` 時為警告；`--fix` 會執行 `chmod 0600`。接著讀取登入資訊（在 macOS 上會先從鑰匙圈讀取）。沒有登入是警告，並提示使用 `ccam login <name>`。更新權杖過期是錯誤 `login expired`；存取權杖過期只是資訊，因為 Claude Code 會在下次啟動時更新它。有效的登入會顯示電子郵件和方案（`login valid (dev@company.com, Claude Team · Acme)`）。這些都不會修復；`ccam login --expired` 會登入所有需要重新登入的帳號。
5. **共用連結**（`share.Check`）。只用 `lstat` 和 `readlink`，把 `share.Layout` 中的每個路徑和帳號的群組比對。它會回報四種情況：缺少連結、連結指向別處、應該是連結的地方是實際檔案，以及已改為 local 的項目仍殘留連結。缺少連結是警告，其他都是錯誤。在較新版本的 ccam 中新增到某個項目的路徑（對話記錄的 `file-history/`、`tasks/` 和 `paste-cache/`，設定的 `keybindings.json`），只要帳號仍有自己的副本，就只會是警告：`is not shared with the group yet (ccam doctor --fix moves it in)`。啟動絕不會因為它們而停止。
6. **未知項目**（`checkUnknownEntries`）。列出帳號資料夾中 ccam 不認得的頂層名稱（`share.Known`），作為 Claude Code 資料夾結構可能已變更的提示。僅供參考。
7. **外掛路徑**（`checkPlugins`）。外掛登錄檔（`installed_plugins.json`、`known_marketplaces.json`）儲存的是絕對路徑。在每個群組的 `plugins/` 中，經過某個帳號資料夾的路徑會在該帳號被移除或重新命名時失效。`--fix` 會用 `profile.RewritePluginPaths` 把它們改寫為群組資料夾，以原子寫入的方式更新登錄檔。
8. **MCP 伺服器**（`checkMCP`）。每個群組的 `mcp.json` 只載入一次。無法讀取的檔案是錯誤；未通過 `mcp.Validate` 的伺服器是警告。不修復。
9. **用量回應格式**（`checkUsageFormat`）。不使用快取，以一般的 3 秒限制逐一向用量端點查詢帳號，直到有一個回應為止。如果回應已經無法解析，它會警告 `ccam usage` 將改用狀態列的快取。這是唯一使用網路的檢查。
10. **快捷指令**（`checkShims`）。只在有帳號設定了快捷指令時執行。`~/.local/bin` 中缺少連結是警告；`--fix` 會呼叫 `shim.Install`，它只建立缺少的連結，絕不碰其他檔案。bin 資料夾不在 `PATH` 上時為警告，不會修復。
11. **資料夾設定**（`checkFolders`）。只在有 [`ccam folder`](/features/folder) 項目時執行。項目的帳號不存在時為錯誤，項目的資料夾已不存在時為警告。兩者都不會修復，因為移除項目會改變在那裡啟動的帳號。
12. **群組和設定**（`checkGroups`）。不是帳號的群組成員、名稱後來變成 ccam 指令的帳號（用 `ccam run <name>` 啟動它，或重新命名），以及 `config.toml` 中未知的 `theme` 或 `usage_source`，都是警告。不會修復。
13. **登入身分**（`checkIdentities`）。從每個資料夾的 `.claude.json` 讀取它登入的身分。資料夾現在登入的帳號和 `ccam login` 記錄的不同時（例如在工作階段中執行 `/login` 之後）是警告；資料夾和另一個資料夾登入同一個 Claude 帳號時也是警告，因為兩者會共用用量上限。不會修復。

### 修復共用連結

啟動流程會當場建立缺少的連結，但遇到任何可能遺失資料的情況就會停止，並提示你使用這個指令：

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

遇到任何共用連結問題時，`--fix` 會以帳號的共用設定同時作為目前狀態和目標狀態呼叫 `share.Build`，再執行 `Plan.Apply`。`Build` 會根據磁碟上的實際內容重新推導每個步驟，產生的計畫和 `ccam group` 變更使用的相同，所以修復遵循的規則和加入群組時一樣。

<!-- motion:doctor-repair -->

這裡的 `settings.json` 是實際檔案。因為群組已經有內容不同的 `settings.json`，衝突會採用預設選擇 `KeepGroup`，保留群組的副本。帳號的檔案會移到帳號資料夾內的 `.ccam-backup/<time>/settings.json`，並由 `os.Symlink` 把該路徑連結到群組。其他情況：

- 如果群組還沒有副本，帳號的檔案會移到群組中，成為群組的第一份內容。
- 如果是資料夾，群組沒有的檔案會移到群組中，其餘的放到備份。
- `history.jsonl` 的預設方式是 `merge by line`：先是群組的各行，再接上你新增的各行。群組原本的副本會保存在 `.ccam-backup/<time>/group/`。
- 指向其他群組的連結會被取代。指向 ccam 以外資料夾的連結會放到備份，如果群組是空的，就以其內容作為群組的初始內容。
- 已改為 local 的項目所殘留的連結會被移除，帳號會取得一份群組內容的副本。

除了移除連結以外，每個步驟都是把資料移到群組或備份資料夾。不會刪除任何東西。

## 安全性和特殊情況

- 單純執行 `ccam doctor` 絕不會變更任何東西。在腳本中若要確保不變更，請把 `--dry-run` 和 `--fix` 一起使用。
- 對於正在執行的帳號，`doctor --fix` 只會建立缺少的連結，並加上 `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``，因為工作階段會持續寫入修復時要移動的檔案。
- 如果共用修復失敗，ccam 會對該帳號回報 `could not repair share links: <error>`，並略過它的未知項目檢查。其他帳號仍會繼續檢查。
- 如果 `ccam group` 變更中途失敗，新的群組已經儲存，ccam 會告訴你用 `ccam doctor --fix` 完成。
- 符號連結的 `.credentials.json`、過期的登入和環境變數中的登入資訊只會回報，絕不會變更。
- 外掛登錄檔的改寫以原子方式寫入，但不會複製到 `.ccam-backup/`。

## 相關內容

- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [複製現有的設定](/features/migrate)
- [用量限制與切換帳號](/docs/usage-limits)
- [指令參考](/docs/commands)
