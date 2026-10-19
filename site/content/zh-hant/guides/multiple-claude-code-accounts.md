---
title: 如何在一台電腦上使用兩個 Claude Code 帳號
description: 透過 CLAUDE_CONFIG_DIR 同時執行多個 Claude Code 帳號，可以用 Shell 別名手動設定，或交給 ccam 自動處理。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

為每個帳號指定獨立的設定資料夾，就能在一台電腦上使用多個 Claude Code 帳號。Claude Code 會從 `CLAUDE_CONFIG_DIR` 環境變數讀取資料夾位置，每個資料夾各自保存登入資訊、設定和對話記錄。你可以用 Shell 別名手動設定，或交給 ccam 處理。

## 手動設定：使用別名

這是 Claude Code 文件中說明的方法。在 `~/.bashrc` 或 `~/.zshrc` 中，為每個額外的帳號加上一行：

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

開啟新的終端機並執行 `claude-work`。因為資料夾是新的，Claude Code 會要求你登入。之後，`claude` 是你的第一個帳號，`claude-work` 則是第二個。

這種方法在兩個帳號時很好用。需要注意的事項：

- 每個資料夾一開始都是空的：沒有設定、沒有 `CLAUDE.md`、沒有 MCP 伺服器，也沒有 agents 或 skills。你需要自己複製或連結需要的內容，並自行讓這些副本保持一致。
- 別名只存在於互動式 Shell 中。腳本、Makefile 和其他程式都看不到它們。
- MCP 伺服器存放在 `.claude.json` 中，和登入資訊放在一起，所以無法直接在資料夾之間連結。
- 對話會留在開始時的帳號中，因此無法在另一個帳號接續。
- 這種做法無法一次查看所有帳號的用量。

## 使用 ccam

ccam 將同樣的做法自動化，並補上缺少的部分：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

現在 `claude1` 和 `claude2`，或 `ccam personal` 和 `ccam work`，就能啟動這兩個帳號。因為它們是真正的程式而不是別名，所以在腳本中也能使用。設定、`CLAUDE.md`、MCP 伺服器和擴充功能從一開始就會共用；你也可以開啟對話記錄共用，再用 `ccam next --continue` 在另一個帳號接續對話。

如果你已經在使用 `~/.claude`，請執行 `ccam import`，不必再次新增該帳號。

## 該選擇哪一種

如果只有兩個帳號，而且不需要共用任何內容，用別名就夠了。如果你想共用設定和 MCP 伺服器、使用更多帳號、一眼看到用量或接續對話，ccam 能幫你省去這些管理工作。[比較](/compare/aliases)頁面並列了兩者的差異。
