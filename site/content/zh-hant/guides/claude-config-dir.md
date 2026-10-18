---
title: CLAUDE_CONFIG_DIR 說明
description: CLAUDE_CONFIG_DIR 在 Claude Code 中的作用、資料夾裡存放什麼、macOS 把登入資訊放在哪裡，以及哪些檔案可以在帳號之間共用。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` 告訴 Claude Code 要使用哪個資料夾來取代 `~/.claude`。所有和帳號相關的內容都存放在這裡：登入資訊、設定、MCP 伺服器、外掛和對話記錄。將它指向不同的資料夾，就能得到彼此獨立的帳號。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## 資料夾裡有什麼

已使用 Claude Code 2.1.302 驗證：

| 路徑 | 內容 | 可共用 |
| --- | --- | --- |
| `.credentials.json` | 登入權杖（Linux）、MCP OAuth 權杖 | 否 |
| `.claude.json` | 帳號詳細資料、使用者 MCP 伺服器、各專案的狀態 | 否（MCP 伺服器另行處理） |
| `settings.json` | 權限、hooks、模型、狀態列、外掛 | 是 |
| `CLAUDE.md` | 全域指示 | 是 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 擴充功能 | 是 |
| `projects/` | 以 `.jsonl` 檔案儲存的對話、各專案的記憶 | 是 |
| `history.jsonl` | 提示詞記錄 | 是 |
| `sessions/`, `shell-snapshots/`, `cache/` 等 | 執行階段狀態 | 否 |

## 兩個容易出錯的細節

**沒有設定 `CLAUDE_CONFIG_DIR` 時，`.claude.json` 不在 `~/.claude` 裡面。** 此時 Claude Code 會把它放在資料夾旁邊的 `~/.claude.json`。如果你開始使用 `CLAUDE_CONFIG_DIR=~/.claude`，Claude Code 會改為尋找 `~/.claude/.claude.json`，因而找不到你的帳號。`ccam import` 會將這兩者連結起來以處理這個問題。

**在 macOS 上，登入資訊存放在鑰匙圈中，而不是資料夾裡。** Claude Code 以 `Claude Code-credentials` 這個名稱儲存登入資訊；設定了 `CLAUDE_CONFIG_DIR` 時，名稱後面還會加上資料夾路徑的簡短雜湊值。因此不同的資料夾會有各自的鑰匙圈項目。如果鑰匙圈拒絕寫入，Claude Code 會改用資料夾中的 `.credentials.json`。

## 在資料夾之間共用

符號連結可以正常運作：Claude Code 會透過連結讀寫 `settings.json` 或 `projects/` 資料夾，並保留該連結。MCP 伺服器是例外，因為它們和登入資訊共用 `.claude.json`；請改用 `--mcp-config` 傳入。

ccam 正是建立在這些規則之上。關於 ccam 如何連結資料夾，請參閱[共用](/docs/sharing)；整體概念請參閱[多帳號指南](/guides/multiple-claude-code-accounts)。
