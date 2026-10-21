---
title: 儀表板
description: 在終端機中執行 ccam 會開啟全螢幕儀表板，顯示帳號、群組、資料夾規則、擴充功能和檢查結果；每項變更都會先顯示計畫。
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

不帶參數的 `ccam` 在標準輸入和輸出都是終端機時會開啟全螢幕儀表板，否則會印出 `ccam list`。儀表板會顯示每個帳號及其用量和群組，並能完成各個指令的工作：啟動、新增、登入、變更群組、設定資料夾規則、檢查和修復。啟動帳號時儀表板會關閉；新增帳號或登入時則會暫時離開儀表板、交給終端機，完成後再回來。

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

分頁名稱下方的一行用一句話說明目前分頁的用途。在較寬的終端機上，Accounts 分頁還會顯示每個帳號的電子郵件。

第一個畫面只根據 `config.toml` 繪製，遠低於 100 ms；用量會在每個帳號回應時陸續填入。每項限制都有自己的長條：5 小時區間、所有模型共用的每週區間，以及 Fable 自己的每週區間；`-` 表示該帳號沒有這個區間。百分比從 80% 起變成黃色，達到 100% 時變成紅色；查詢失敗後改用快取的數值會以淡色顯示。

## 分頁

| 分頁 | 顯示的內容 | 按鍵 |
| --- | --- | --- |
| Accounts | 編號、名稱、電子郵件、5 小時、每週和 Fable 用量長條、上次使用時間、群組；`✔` 為預設帳號，`▸` 為這個資料夾的帳號，`•` 為執行中 | 見下文 |
| Groups | 每個群組、它共用的項目和它的帳號，以及不屬於任何群組的帳號 | `enter` 編輯，`n` 新增，`e` 重新命名，`x` 移除 |
| Folders | [`ccam folder`](/features/folder) 規則，以及適用於這裡的規則 | `n` 在這裡使用某個帳號，`enter` 變更，`x` 移除 |
| Extensions | [`ccam ext`](/features/ext) 顯示的內容，每個擴充功能只列一次，並附上它來自的群組 | `r` 重新載入 |
| Health | [`ccam doctor`](/features/doctor) 回報的內容，每個帳號一個方框（見下文） | `f` 修復選取的方框，`F` 修復全部（`doctor --fix`），`↑`/`↓` 在方框之間移動 |

`←`/`→` 或 `Tab` 切換分頁，`?` 列出目前分頁的按鍵，`r` 會略過短期快取重新查詢用量，`q` 或 `Esc` 則結束。

## Accounts 分頁的按鍵

| 按鍵 | 作用 |
| --- | --- |
| `Enter` | 以選取的帳號啟動 Claude Code |
| `1`–`9` | 移到該編號的帳號；快速輸入 `1` `2` 可移到 12。只有 `Enter` 會啟動 |
| `↑`/`↓`、`k`/`j` | 移動 |
| `n` | 啟動 [`ccam next --auto`](/features/next) 會選擇的帳號 |
| `o` | 帶選項啟動：`--continue`、`--resume`、其他參數，或試執行 |
| `u` | 用量詳細資料：5 小時、每週和各模型的區間，以及它們的重設時間和數字的來源 |
| `a` | 新增帳號（使用下文的設定精靈） |
| `l` | 重新登入 |
| `s` | 選擇帳號的群組 |
| `d` | 設為預設帳號 |
| `f` | 在目前的資料夾中使用它 |
| `e` | 重新命名 |
| `c` | 設定或移除它的快捷指令 |
| `x` | 移除它（會先詢問，游標預設停在 No） |

## Health 分頁

開啟 Health 分頁時，它會執行 [`ccam doctor`](/features/doctor) 的檢查並繪製相同的畫面：第一行為每個帳號顯示一個標記，下方每個帳號各有一個方框。方框內每一列是一個檢查過的項目，`Folder`、`Login` 和 `Share links` 排在最前面，標記為 `✔` 正常、`ℹ` 資訊、`⚠` 警告或 `✘` 錯誤，讓你一眼看出哪個帳號在哪裡有問題。不屬於任何帳號的檢查結果放在 `General` 方框中。

```console
 ✘ personal   ⚠ work

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

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` 在方框之間移動，選取的方框會以強調色邊框顯示。`f` 只修復選取的方框，方式與 `ccam doctor --fix` 相同：在帳號的方框上只修復該帳號，在 `General` 方框上修復不屬於任何帳號的檢查結果（外掛路徑和快捷指令）；其他帳號只會重新檢查。它會先詢問，沒有可修復項目的方框只會顯示 `nothing to repair in <name>`。`F` 會像 `ccam doctor --fix` 一樣一次修復全部；`r` 會重新檢查。

## 變更群組

在 Groups 分頁中按 `Enter` 會開啟編輯器：左邊是項目，右邊是帳號，用 `Tab` 切換，用 `Space` 勾選或取消。其中某個帳號已經從另一個群組取得的項目會顯示為灰色並附上原因，會重複取得某個項目的帳號也一樣；一個項目只能來自一個群組。按 `Enter` 會顯示計畫：

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

和群組副本內容不同的檔案會先詢問你（保留群組的版本、使用我的版本，或合併）。接著變更會在背景以和 [`ccam group`](/features/share) 相同的程式碼執行：先儲存設定，再依序移動每個帳號的資料，原始檔案則放到 `.ccam-backup/`。會移動執行中帳號資料的變更會被拒絕，並列出該帳號的名稱。

## 新增帳號

`a` 會開啟 `ccam add` 在終端機中也使用的設定精靈：

1. 帳號名稱。
2. 登入方式：現在以 Claude 訂閱登入（在瀏覽器中）、[API 金鑰](/features/api)，或稍後再以 Claude 訂閱登入。
3. 使用 API 金鑰時：ccam 每次啟動時從哪裡讀取金鑰（`env:NAME` 或 `command:COMMAND`）、端點、金鑰的傳送方式和模型。
4. 要加入的群組，其中 `default` 已勾選；還沒有任何群組時，則是群組 `default` 共用的項目。
5. 快捷指令，如果你需要的話。
6. 讓你確認的摘要。

`Esc` 會回到上一步，在第一步時則會關閉設定精靈。

## 色彩

色彩依循 Claude Code 終端機介面的角色配色，所以 ccam 放在它旁邊看起來很協調；群組名稱會根據名稱取得一種顏色。ccam 會向終端機查詢背景色，並選擇淺色或深色主題。你也可以在 `config.toml` 中自行設定：

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 主題適用於紅綠色盲，`-ansi` 主題則使用終端機本身的 16 種顏色。`NO_COLOR` 會關閉色彩。

## 安全性和特殊情況

- 儀表板呼叫的函式和指令相同，所以每條規則都成立：移動資料之前先顯示計畫、拒絕變更執行中的帳號、保留備份，而且登入權杖絕不會離開帳號。
- 儀表板沒有試執行：不帶名稱的 `ccam --dry-run` 會印出這個資料夾的帳號會如何啟動，在某個帳號上按 `o` 則可以只對這一次啟動進行試執行。
- 在 Claude Code 啟動或登入開啟之前，終端機會先還原。
- 畫面太小放不下對話框時，會先省略分頁的列；過長的行會以 `…` 截斷，而不是換行。

## 相關內容

- [用群組共用](/features/share)
- [啟動 Claude Code](/features/launch)
- [管理帳號](/features/accounts)
- [用量限制](/features/usage)
