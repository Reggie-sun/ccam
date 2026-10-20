---
title: 所有帳號的背景工作階段
description: ccam agents 會在一張表格中列出你所有帳號的 Claude Code 背景代理程式和開啟中的工作階段，ccam attach 則依 ID 開啟其中一個。
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code 會依帳號分別保存背景工作階段（`claude --bg`、agent view）：`claude agents` 只會顯示它所執行的那個帳號的工作階段。`ccam agents` 會同時查詢所有帳號並一起列出，`ccam attach <id>` 則會以啟動該背景工作階段的帳號開啟它。

## 用法

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` 是背景工作階段的狀態（`working`、`blocked`、`done`、`failed`、`stopped`）以及它的行程正在做什麼（`busy`、`waiting`、`idle`）；開啟中的互動式工作階段會顯示 `open`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--all` | 也列出已完成、失敗或已停止的背景工作階段 | 關閉 |
| `--json` | 以 JSON 印出工作階段，每個都附上它所屬的 `account` | 關閉 |

`ccam attach <id> [claude args]` 會把 ID 後面的參數傳給 `claude attach`。加上 `--dry-run` 時會印出指令。

## 運作方式

1. ccam 會為每個帳號，以該帳號的 `CLAUDE_CONFIG_DIR` 同時執行 `claude agents --json`（指定 `--all` 時也會加上），每個都限制在 15 秒內。這只會讀取帳號自己的工作階段清單，不會啟動 Claude Code 的背景服務。
2. 每個項目都會附上帳號的名稱，清單依啟動時間排序，最新的在前。沒有回應的帳號會在標準錯誤中列出名稱，並且不列入清單，接著 ccam 會以結束代碼 1 結束，讓指令碼不會把失敗誤認為空清單。
3. `ccam attach` 會在每個帳號的資料夾中尋找 `jobs/<id>`，也就是 Claude Code 存放背景工作階段的位置，然後以該帳號執行 `claude attach <id>`，方式和 `ccam <name>` 啟動 Claude Code 相同。

## 安全性和特殊情況

- 背景工作階段的 ID 屬於單一帳號：以另一個帳號執行的 `claude attach` 找不到它們，而這正是 `ccam attach` 要解決的問題。
- 兩個帳號可能有 ID 開頭相同的工作階段；這時 `ccam attach` 會列出兩者，並要求你輸入更長的 ID。
- ccam 以單一參數 `--mcp-config=<file>` 把共用的 MCP 伺服器傳給 Claude Code，`claude agents` 和 `claude attach` 都接受這種寫法；背景工作階段重新啟動時也會繼續使用那個檔案。
- API 金鑰帳號的背景服務會從啟動它的 Shell 取得金鑰。請透過 ccam（`ccam gw --bg …`）啟動這種帳號的第一個背景工作階段，讓背景服務取得金鑰。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [為指令碼選擇帳號](/features/pick)
- [儀表板](/features/dashboard)
