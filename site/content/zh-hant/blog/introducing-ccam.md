---
title: ccam 介紹
description: ccam 為什麼存在、它能為擁有多個 Claude Code 帳號的人做些什麼，以及它如何做到不礙事。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: b116c91bb8d4
---

ccam 是一個小型命令列工具，適合在同一台電腦上使用多個 Claude Code 帳號的人：例如一個個人帳號和一個工作帳號，或是為一個長期專案準備的多個帳號。

## 問題所在

Claude Code 透過 `CLAUDE_CONFIG_DIR` 支援多個帳號：每個帳號一個資料夾。實際使用時，每個資料夾都是一座孤島。設定、`CLAUDE.md`、MCP 伺服器、agents 和 skills 都得在每個資料夾中重新設定，並手動保持一致。對話會留在開始時的地方。而當某個帳號達到用量限制時，你要等到嘗試時才會發現。

## ccam 能做什麼

- `ccam add` 會建立資料夾、執行登入；如果你需要，還會為帳號設定一個你自己命名的快捷指令，例如 `ccwork`。
- `ccam work` 和 `ccwork` 會以該帳號啟動 Claude Code，並將所有參數原封不動地傳遞過去。
- MCP 伺服器、設定、`CLAUDE.md` 和擴充功能透過連結共用，可以逐項、逐帳號設定。登入資訊永遠不會共用。
- `ccam usage` 會顯示每個帳號的 5 小時用量和每週用量；`ccam next --continue` 會在下一個帳號接續對話。

## 不礙事

一個在每次工作階段前都會執行的工具，不應該帶來任何你察覺得到的負擔。ccam 會讀取一個設定檔，對每個共用連結各用一次系統呼叫進行檢查，然後以 `claude` 取代自身。它大約只增加 4 毫秒，而且你工作時它並不在執行。

它也會謹慎處理你的檔案。每次變更共用設定前，都會先印出計畫，並將原始檔案保留在備份資料夾中。`ccam doctor` 會檢查連結、權限和登入狀態，並修復可以安全修復的部分。

## 試試看

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam 是開放原始碼的非官方工具：它與 Anthropic 沒有關聯，也未獲得 Anthropic 認可。能否使用多個帳號取決於各帳號的條款。歡迎在 [GitHub](https://github.com/jhyunwoo/ccam) 上提供意見回饋和翻譯。
