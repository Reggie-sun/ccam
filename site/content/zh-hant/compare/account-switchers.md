---
title: ccam 與其他 Claude Code 帳號切換工具比較
description: ccam 與 claude-swap、CCS、clauth、CAAM、cux 等 Claude Code 帳號切換工具的比較，以及哪種做法適合哪種用途。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

切換 Claude Code 帳號的工具採用兩種做法之一。大多數工具在同一個 `~/.claude` 資料夾中替換登入資訊，所以同一時間只有一個帳號處於使用中。另一些工具（包括 ccam）則用 `CLAUDE_CONFIG_DIR` 為每個帳號提供獨立的資料夾，因此每個帳號都會保持登入，也能同時執行多個帳號。ccam 還提供共用的設定和 MCP 伺服器，並能在另一個帳號的資料夾中接續對話。

## 替換同一組登入資訊，還是使用獨立的資料夾？

**替換登入資訊**只保留一套設定、MCP 伺服器和對話記錄，因為只有一個資料夾。切換很快，執行中的工作階段也能套用新的登入資訊。但同一時間只有一個帳號處於使用中，而且會同時影響所有終端機和編輯器。

**獨立的資料夾**是 [Claude Code 文件中說明的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。每個帳號都保持登入，你可以讓工作和個人的工作階段同時執行。代價是每個資料夾一開始都有各自的設定、MCP 伺服器和對話記錄，除非有工具幫你共用。

## 各工具比較

資料取自各專案自己的文件（2026-10-19）。「未說明」表示專案沒有提到，而不是缺少該功能。

| 工具 | 做法 | 同時使用多個帳號 | 5 小時和每週用量 | 達到限制時 | 平台 |
| --- | --- | --- | --- | --- | --- |
| ccam | 獨立的資料夾 | 是 | 是，`ccam usage` | 由你執行 `ccam next --auto --continue`，在剩餘額度最多的帳號中接續對話 | Linux、macOS、WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | 替換登入資訊；`cswap run` 時使用獨立的資料夾 | 是，使用 `cswap run` | 是 | `cswap auto` 在達到門檻時切換登入資訊；執行中的工作階段會繼續 | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 獨立的資料夾 | 是 | 部分支援，在 macOS 選單列 App 中 | 帳號池透過本機 Proxy 把請求轉給下一個帳號 | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 替換登入資訊；`clauth start` 時使用獨立的資料夾 | 是，使用 `clauth start` | 是 | 沿著備援鏈依序切換 | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 替換登入資訊，或使用隔離的設定檔 | 是，使用隔離的設定檔 | 是 | `caam run` 在另一個帳號上重試非互動式指令 | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | 替換登入資訊並包裝 `claude` | 未說明 | 是 | 替換登入資訊並恢復同一個對話 | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 替換登入資訊 | 否 | 是 | 未說明 | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 替換登入資訊 | 否 | 是 | 可選擇在 100% 時切換 | macOS |

ccswap 是 claude-swap 的分支，也能管理 Codex 的登入。其中幾個工具也支援其他程式設計代理，例如 Codex 或 Gemini。

## ccam 的重點

- 每個帳號都保持登入。切換時絕不會讓帳號登出，`ccam work` 和 `ccam personal` 可以同時在兩個終端機中執行。
- MCP 伺服器、`settings.json`、`CLAUDE.md`、agents、skills 和外掛透過連結在資料夾之間共用，可以逐項、依帳號選擇。登入資訊絕不共用。
- 共用對話記錄後，只要一個指令 `ccam next --auto --continue`，就能在剩餘額度最多的帳號的資料夾中開啟同一個對話。
- ccam 啟動 Claude Code 後就用它取代自己，所以沒有任何東西在背景執行。它不會監看你的用量，也不會自行切換；何時切換由你決定。

## 其他工具做得更多的地方

- claude-swap、CCS、clauth、CAAM 和 cux 不需要 WSL 就能在 Windows 上執行。ccam 不行。
- claude-swap、clauth 和 cux 會在帳號超過用量門檻時自行切換。ccam 把切換交給你。
- CCSwitcher、Claude Switcher 和 claude-swap 會在 macOS 選單列中顯示用量。
- CCS、CAAM、ccswap 和 Claude Switcher 也能管理 Codex 和其他程式設計代理。

## 哪一種適合你？

如果你一次只用一個帳號、希望自動切換，或是在沒有 WSL 的 Windows 上，請選擇替換登入資訊的工具。如果你想讓多個帳號同時執行、共用一套設定，並在下一個帳號接續同一個對話，請選擇 ccam。如果什麼都不需要共用，[一般的 Shell 別名](/compare/aliases)就夠了。

是否可以使用多個帳號取決於各帳號的條款。使用前請先確認。
