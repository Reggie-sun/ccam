---
title: 在一台電腦上使用多個 Claude Code 帳號
description: 透過 CLAUDE_CONFIG_DIR 同時執行多個 Claude Code 帳號並在帳號之間切換，可以用 Shell 別名手動設定，或交給 ccam 自動處理。
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "381de7a473b2"
---

為每個帳號指定獨立的設定資料夾。Claude Code 會從 `CLAUDE_CONFIG_DIR` 環境變數讀取資料夾位置，每個資料夾各自保存登入資訊、設定和對話記錄，所以每個帳號都會保持登入，你也可以同時執行它們。你可以用 Shell 別名手動設定，或交給 ccam 處理。

## 如何手動設定？

這是 [Claude Code 文件中說明的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。在 `~/.bashrc` 或 `~/.zshrc` 中，為每個額外的帳號加上一行：

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

開啟新的終端機並執行 `claude-work`。因為資料夾是新的，Claude Code 會要求你登入。之後，`claude` 是你的第一個帳號，`claude-work` 則是第二個。

在沒有 WSL 的 Windows 上，可以在 PowerShell 設定檔（`notepad $PROFILE`）中加入一個函式，效果相同：

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## 別名做法有哪些不足？

這種做法在兩個帳號時很好用。需要注意的事項：

- 每個資料夾一開始都是空的：沒有設定、沒有 `CLAUDE.md`、沒有 MCP 伺服器，也沒有 agents 或 skills。你需要自己複製或連結需要的內容，並自行讓這些副本保持同步。
- 別名只存在於互動式 Shell 中。腳本、Makefile 和其他程式都看不到它們。
- MCP 伺服器存放在 `.claude.json` 中，和登入資訊放在一起，所以無法直接在資料夾之間連結。
- 對話會留在開始時的帳號中，因此無法在另一個帳號接續。
- 沒有辦法一次查看所有帳號的用量。

## 可以改用 /login 切換帳號嗎？

可以，前提是你只是偶爾切換。在 Claude Code 中輸入 `/login` 會[切換目前資料夾的帳號](https://code.claude.com/docs/en/quickstart)。同一時間只有一個帳號處於登入狀態，所以每次切換都要重新登入，也無法讓兩個帳號同時執行。使用獨立的資料夾就能避免這兩個問題。

## 如何用 ccam 設定？

ccam 將同樣的做法自動化，並補上缺少的部分：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

現在 `ccam personal` 和 `ccam work` 就能啟動這兩個帳號，你為 work 帳號設定的快捷指令 `ccwork` 也可以。因為它們是真正的程式而不是別名，所以在腳本中也能使用。單獨執行 `ccam` 可以從清單中選擇帳號，執行 `ccam use work` 則可將某個帳號設為預設。

設定、`CLAUDE.md`、MCP 伺服器和擴充功能從一開始就會共用。`ccam next --continue` 會在另一個帳號接續對話。如果你開啟對話記錄共用，兩個帳號會處理同一段對話，而不是一份副本。`ccam usage` 會一次顯示所有帳號的 5 小時用量和每週用量。

如果你已經在使用 `~/.claude`，請執行 `ccam import`，不必再次新增該帳號；或執行 `ccam migrate`，把其中的設定、MCP 伺服器和擴充功能複製到共用的 `default` 群組。

## 如何知道目前執行的是哪個帳號？

在 Claude Code 中，`/status` 會顯示你目前登入的帳號。使用 ccam 時，執行一次 `ccam statusline install`，狀態列就會一直顯示帳號名稱和用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。

## 在 Windows 和 WSL 上可以使用嗎？

Claude Code [可以在 Windows 上原生執行，也可以在 WSL 中執行](https://code.claude.com/docs/en/setup)。`CLAUDE_CONFIG_DIR` 在兩者中都有效，所以在原生 Windows 上，用上面的別名或 PowerShell 函式就夠了。ccam 也能在 Windows 上原生執行：在 PowerShell 中用 `irm https://ccam.dev/install.ps1 | iex` 安裝，或在 WSL 中用 Linux 的指令安裝。快捷指令會變成 `ccwork.exe`，而共用 `settings.json` 這類檔案需要開發人員模式（[安裝 ccam](/docs/install)）。

## 疑難排解

**Claude Code 每次都要求我登入。** 在新資料夾中第一次啟動時一定會要求登入。如果一直要求登入，表示登入已過期或已被移除：執行 `ccam login work`（或透過別名啟動一次並登入）。

**某個 MCP 伺服器只出現在一個帳號中。** 在工作階段中用 `claude mcp add` 新增的伺服器屬於該帳號，而且沒有加上 `-s user` 時只屬於該資料夾。用 `ccam mcp import work <server>` 將它移到所有帳號。

**各帳號的設定不一致。** `ccam share work` 會列出該帳號共用的內容，以及透過哪個群組共用；`ccam group join default work` 會先顯示將要變更的內容，再讓它共用 `default` 共用的項目。`ccam doctor --fix` 會修復損壞的連結，並為它修改的每個檔案保留備份。

**兩個帳號開啟了同一個對話。** 避免在兩個帳號中同時處理同一個對話。你這樣做時，ccam 會發出警告。

## 該選擇哪一種？

如果只有兩個帳號，而且不需要共用任何內容，用別名就夠了。如果你想共用設定和 MCP 伺服器、使用更多帳號、一眼看到用量或接續對話，ccam 能幫你省去這些管理工作。[與別名的比較](/compare/aliases)並列了兩者的差異，[ccam 與其他帳號切換工具](/compare/account-switchers)則介紹了那些切換同一組登入資訊的工具。

能否使用多個帳號，取決於適用於各帳號的條款。開始之前請先閱讀這些條款。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。
