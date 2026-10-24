---
title: 檢查擴充功能
description: ccam ext 顯示每個帳號載入的技能、指令、外掛、MCP 伺服器和 hook，各自的來源，以及哪些已損壞。
updated: 2026-10-23
order: 14
translation: machine
sourceHash: "45740e4efbef"
---

`ccam ext` 回答「為什麼這個技能（或 MCP 伺服器）在一個帳號中能用，在另一個帳號中卻不行？」。它會為每個帳號列出 Claude Code 載入的擴充功能，以及每一項的來源：共用群組、帳號自己的資料夾，或你所在的專案。表格下方會列出讓擴充功能無法如預期載入的問題。它只會讀取。

## 用法

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

每個儲存格表示該項目在該帳號中的來源：共用群組（`default`）、位於帳號自己資料夾中時為 `account`，帳號只為這個資料夾設定的 MCP 伺服器為 `local`，帳號沒有該項目時為 `-`。已安裝但未開啟的外掛，會在來源後面顯示 `(off)`。`!` 標示有問題的項目，說明列在表格下方；不屬於任何帳號的問題（例如只有直接執行 `claude` 時才有的伺服器）會不帶名稱列在最後。`ccam ext --json` 會印出同樣的報告（`accounts`、專案、`problems`），Claude Code 透過網路連線的 MCP 伺服器會帶有 `remote`；沒有任何帳號時，會印出 `accounts: []` 的報告。

## 讀取的內容

| 種類 | 位置 | 顯示的來源 |
| --- | --- | --- |
| 技能 | `skills/` 中每個含有 `SKILL.md` 的資料夾 | `skills/` 連結到的群組，或 `account` |
| 指令、代理 | `commands/` 和 `agents/` 底下的 `.md` 檔案，子資料夾顯示為 `folder/name` | 群組，或 `account` |
| 外掛 | `plugins/installed_plugins.json`，由帳號和專案設定中的 `enabledPlugins` 開啟或關閉；只為某個專案安裝的外掛只會出現在該專案中 | 群組，或 `account`；未開啟時為 `(off)` |
| MCP 伺服器 | 群組的 `mcp.json`（共用 MCP 時）、帳號的 `.claude.json`（它自己的伺服器，以及只為這個資料夾設定的伺服器），以及其外掛的伺服器，顯示為 `plugin:<plugin>:<server>` | 群組、`account` 或 `local` |
| Hook | `settings.json` 中的 `hooks`，依事件計數 | 群組，或 `account` |
| 這個專案 | 目前資料夾中的 `.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json` 和 `.mcp.json` | 列出一次；尚未信任該資料夾的帳號，會為其 `.mcp.json` 伺服器得到一個問題 |

## 回報的問題

- **損壞的連結**：`skills/`、`commands/` 或 `agents/` 中目標已不存在的連結，或連結到不存在位置的共用資料夾（`skills`、`commands`、`agents`、`plugins`）（用 `ccam doctor --fix` 修復，它會重新建立缺少的群組資料夾）。
- **沒有 `SKILL.md` 的技能資料夾**：Claude Code 會略過它。
- **在兩個地方定義不同的 MCP 伺服器**：群組的 `mcp.json` 和帳號的 `.claude.json`，或其中之一和專案的 `.mcp.json`。兩個地方的設定相同（`mcp_strategy = "sync"` 就會這樣寫入）並不是問題。用 `ccam mcp import <account> <server>` 把帳號自己的伺服器移到群組中。
- **不是有效 JSON 的 `installed_plugins.json`**、`enabledPlugins` 中沒有列出的外掛，以及資料夾不存在的外掛。
- **只有一個帳號載入的伺服器**：只在這個帳號的 `.claude.json` 中、群組中其他帳號不會載入的伺服器，以及位於該帳號針對這個資料夾的 local 範圍中的伺服器（除非加上 `-s user`，否則 `claude mcp add` 就會這樣做）。每一項都會附上可共用它的 `ccam mcp import` 指令。直接執行 `claude` 所用的 `~/.claude.json` 中、沒有任何帳號載入的伺服器，則會附上 `ccam migrate --mcp` 指令。
- **等待核准的專案伺服器**：該帳號尚未信任這個資料夾，所以 Claude Code 會詢問一次，在那之前專案的 `.mcp.json` 伺服器都會處於等待狀態。問題會說明要在哪裡核准它們：設定群組中各帳號，或這個帳號的 `settings.json` 中的 `"enabledMcpjsonServers"`。對 git 儲存庫主要工作樹的信任，在它的其他工作樹中也有效。
- **登入和連接器**：在這個帳號中需要登入、而另一個帳號已登入的 MCP 伺服器（登入只對一個帳號有效；請用 `/mcp` 登入），以及另一個帳號用過、這個帳號尚未用過的 claude.ai 連接器。

## 安全性和特殊情況

- `ccam ext` 不會變更任何東西，也不會發出網路請求。它讀取伺服器名稱和 hook 數量，絕不讀取 MCP 定義中 `env` 或標頭的值。為了檢查登入，它會讀取每個帳號登入資訊中 MCP 登入的名稱（在 macOS 上從鑰匙圈讀取），絕不讀取權杖，而且只針對載入了另一個帳號也載入之伺服器的帳號。
- 它顯示的是磁碟上的內容。被受管理設定封鎖的伺服器仍然會出現；在設定中關閉的外掛會顯示 `(off)`。
- 在專案資料夾中執行它，就會包含該專案的 `.claude` 資料夾和 `.mcp.json`。

## 相關內容

- [帳號之間的共用](/features/share)
- [共用的 MCP 伺服器](/features/mcp)
- [檢查和修復](/features/doctor)
