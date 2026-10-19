---
title: 常見問題
description: 關於 ccam、Claude Code 帳號、用量限制、共用、隱私和支援系統的解答。
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 1355fa6ef3c6
---

## ccam 是 Anthropic 的官方工具嗎？

不是。ccam 是非官方的開放原始碼專案，與 Anthropic 沒有關聯，也未獲得 Anthropic 認可。它使用的是 Claude Code 文件中說明的 `CLAUDE_CONFIG_DIR` 設定，用來同時執行多個帳號。

## 我可以使用多個 Claude Code 帳號嗎？

這取決於各帳號的條款；如果是工作帳號，也取決於你所屬組織的政策。使用多個帳號前請先確認這些規定。ccam 只是讓你在既有的帳號之間切換更方便。

## ccam 能規避用量限制嗎？

不能。每個帳號都保有各自的限制。ccam 會顯示這些限制，並在某個帳號用完額度時，讓你在自己擁有的另一個帳號繼續。

## 如何在一台電腦上使用兩個 Claude Code 帳號？

用 `CLAUDE_CONFIG_DIR` 為每個帳號指定獨立的設定資料夾。ccam 會替你完成這件事：為每個帳號執行 `ccam add`，再用 `ccam personal` 和 `ccam work`，或你自己命名的快捷指令（例如 `ccwork`）啟動。[指南](/guides/multiple-claude-code-accounts)也說明了如何手動設定。

## 如何在另一個帳號繼續對話？

用 `ccam share <name> --history on` 讓兩個帳號共用對話記錄，然後執行 `ccam next --continue`。Claude Code 會在下一個帳號開啟目前資料夾中最近的一段對話。

## ccam 會看到或傳送我的登入資訊嗎？

ccam 絕不會複製或共用登入資訊。它讀取帳號的存取權杖只有一個用途：在你執行 `ccam usage`、`ccam next --auto` 或 `ccam doctor` 時，向 Anthropic 的用量端點查詢該帳號的用量。ccam 沒有遙測功能，也不會將任何資料傳送到其他地方。

## 我的帳號存放在哪裡？

存放在 `~/.local/share/ccam/profiles/<name>/`，每個帳號一個資料夾。在 macOS 上，Claude Code 會把登入資訊本身存放在鑰匙圈中，名稱由該資料夾衍生而來。

## 現有的 ~/.claude 會怎麼樣？

不會有任何變化，除非你執行 `ccam import` 或 `ccam migrate`。匯入會在原位置登錄該資料夾，因此直接執行 `claude` 和 `ccam <name>` 都會繼續使用相同的帳號和設定。`ccam migrate` 會把其中的設定、`CLAUDE.md`、擴充功能和 MCP 伺服器複製到共用設定，而且只會讀取該資料夾。

## ccam 支援哪些系統？

Linux（Ubuntu、Debian、Fedora、Arch、Alpine 等）和 macOS。在 Windows 上，請在 WSL 中使用。

## 為什麼某個 MCP 伺服器只出現在一個帳號中？

在工作階段中用 `claude mcp add` 新增的伺服器只會儲存在該帳號。下次啟動該帳號時，ccam 會提醒你；用 `ccam mcp import <account> <server>` 將伺服器移到共用群組。

## ccam 會拖慢 Claude Code 嗎？

透過 ccam 啟動 Claude Code 大約會增加 4 毫秒，這是在一台筆記型電腦上用 `scripts/perf-check.sh` 測得的。ccam 會讀取一個小檔案、檢查共用連結，然後以 `claude` 取代自身，所以你工作時它並不在執行。

## 如何更新 ccam？

執行 `ccam update`。它會下載最新的發行版本，檢查總和檢查碼，然後取代 `ccam` 檔案；`ccam update --check` 只會告訴你是否有更新的發行版本。ccam 不會自行更新。如果你安裝的是 `.deb` 套件，請改用 apt 安裝新的 `.deb`。請參閱[更新](/docs/install#update)。

## 如何解除安裝 ccam？

用 `ccam shortcut <account> --remove` 移除每個帳號的快捷指令，再從 `~/.local/bin` 刪除 `ccam`；如果不再需要這些帳號，再刪除 `~/.config/ccam`、`~/.local/share/ccam` 和 `~/.local/state/ccam`。請參閱[安裝](/docs/install#uninstall)。
