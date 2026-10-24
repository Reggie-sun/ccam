---
title: CLAUDE_CONFIG_DIR 說明
description: CLAUDE_CONFIG_DIR 在 Claude Code 中的作用、資料夾裡存放什麼、macOS 把登入資訊放在哪裡，以及哪些檔案可以在帳號之間共用。
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR` 告訴 Claude Code 要使用哪個資料夾來取代 `~/.claude`。所有和帳號相關的內容都存放在這裡：登入資訊、設定、MCP 伺服器、外掛和對話記錄。將它指向不同的資料夾，就能得到彼此獨立的帳號。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR 是官方支援的嗎？

是的。Claude Code 在[環境變數](https://code.claude.com/docs/en/env-vars)清單中將它列為覆寫設定資料夾的方式，並註明它「適合同時執行多個帳號」；其[驗證指南](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)也用它來登入多個帳號。你可以在 Shell 或使用者設定中設定它，但不能在專案設定中設定。

## 資料夾裡有什麼？

已使用 Claude Code 2.1.309 驗證：

| 路徑 | 內容 | 可共用 |
| --- | --- | --- |
| `.credentials.json` | 登入權杖（Linux、Windows）、MCP OAuth 權杖 | 否 |
| `.claude.json` | 帳號詳細資料、使用者 MCP 伺服器、各專案的狀態 | 否（MCP 伺服器另行處理） |
| `settings.json` | 權限、hooks、模型、狀態列、外掛 | 是 |
| `CLAUDE.md` | 全域指示 | 是 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 擴充功能 | 是 |
| `projects/` | 以 `.jsonl` 檔案儲存的對話、各專案的記憶 | 是 |
| `history.jsonl` | 提示詞記錄 | 是 |
| `sessions/`, `shell-snapshots/`, `cache/` 等 | 執行階段狀態 | 否 |

## 有哪些容易出錯的地方？

**沒有設定 `CLAUDE_CONFIG_DIR` 時，`.claude.json` 不在 `~/.claude` 裡面。** 此時 Claude Code 會把它放在資料夾旁邊的 `~/.claude.json`。如果你開始使用 `CLAUDE_CONFIG_DIR=~/.claude`，Claude Code 會改為尋找 `~/.claude/.claude.json`，因而找不到你的帳號。`ccam import` 會將這兩者連結起來以處理這個問題。如果 `~/.claude` 已經有自己的 `.claude.json`（來自先前以 `CLAUDE_CONFIG_DIR=~/.claude` 執行的結果），ccam 會保留兩個檔案並告訴你，因為只有你知道哪一個才是對的。

**在 macOS 上，登入資訊存放在鑰匙圈中，而不是資料夾裡。** Claude Code 以 `Claude Code-credentials` 這個名稱儲存登入資訊；設定了 `CLAUDE_CONFIG_DIR` 時，名稱後面還會加上 `-` 以及資料夾路徑 SHA-256 的前 8 個十六進位數字，並歸在 `$USER` 所指名的帳戶下。因此不同的資料夾會有各自的鑰匙圈項目，而被移動或重新命名的資料夾會失去它的登入；這就是 `ccam rename` 在 macOS 上會保留資料夾名稱的原因。如果鑰匙圈拒絕寫入，Claude Code 會改用資料夾中的 `.credentials.json`，[其文件](https://code.claude.com/docs/en/authentication#credential-management)中也有說明。

**並非所有登入方式都存放在資料夾中。** 未使用 API 金鑰的 Claude Console 登入存放在資料夾之外，所以獨立的資料夾無法把兩個這類登入分開。claude.ai 登入和 API 金鑰則依資料夾分別保存。

## 在 Windows 上，這個資料夾在哪裡？

預設是 `%USERPROFILE%\.claude`，登入資訊存放在其中的 `.credentials.json`；若 Claude Code 設定為使用 Windows 認證管理員，則存放在那裡，名稱和 macOS 一樣由資料夾路徑產生。為每個額外的帳號把 `CLAUDE_CONFIG_DIR` 設為另一個資料夾，方法與 Linux 和 macOS 相同。

## 資料夾之間可以共用檔案嗎？

符號連結可以正常運作：Claude Code 會透過連結讀寫 `settings.json` 或 `projects/` 資料夾，並保留該連結。MCP 伺服器是例外，因為它們和登入資訊共用 `.claude.json`；請改用 `--mcp-config` 傳入。

ccam 正是建立在這些規則之上。關於 ccam 如何連結資料夾，請參閱[共用](/docs/sharing)；整體概念請參閱[多帳號指南](/guides/multiple-claude-code-accounts)。
