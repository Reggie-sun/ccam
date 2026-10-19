---
title: ccam 與 Shell 別名的比較
description: 並列比較用 CLAUDE_CONFIG_DIR 搭配 Shell 別名和使用 ccam 來執行多個 Claude Code 帳號的差異。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6fc4200a8a65
---

為每個帳號設定一個 Shell 別名，是執行多個 Claude Code 帳號最簡單的方法。ccam 在底層做的是同一件事，並加上共用、切換和用量追蹤功能。

| | Shell 別名 | ccam |
| --- | --- | --- |
| 獨立的登入 | 是 | 是 |
| 每個帳號的設定步驟 | 編輯 `~/.bashrc`，登入 | `ccam add` |
| 可在腳本和其他程式中使用 | 否，別名只在 Shell 中有效 | 是，`ccam work` 和 `ccwork` 這類快捷指令都是真正的指令 |
| 共用設定、`CLAUDE.md`、擴充功能 | 手動複製或連結 | 逐項設定，並保留備份 |
| 共用 MCP 伺服器 | 手動，每個帳號各設定一次 | 只需執行一次 `ccam mcp add` |
| 在另一個帳號接續對話 | 否 | `ccam next --continue` |
| 所有帳號的用量 | 否 | `ccam usage` |
| 選擇剩餘額度最多的帳號 | 否 | `ccam next --auto` |
| 在狀態列顯示帳號名稱 | 否 | `ccam statusline install` |
| 修復損壞的連結 | 手動 | `ccam doctor --fix` |
| 額外的軟體 | 無 | 一個靜態執行檔 |

如果你有兩個帳號，而且不想在它們之間共用任何內容，就選擇別名。如果你希望這些帳號用起來像同一套環境，就選擇 ccam。你也可以先用別名，之後再改用 ccam：`ccam import` 會在原位置登錄現有的 `~/.claude-*` 資料夾。

[如何在一台電腦上使用多個帳號](/guides/multiple-claude-code-accounts)逐步說明了這兩種做法，[ccam 與其他帳號切換工具](/compare/account-switchers)則比較了 ccam 和那些切換同一組登入資訊的工具。
