---
title: 檢查擴充功能
description: ccam ext 顯示每個帳號載入的技能、指令、外掛、MCP 伺服器和 hook，各自的來源，以及哪些已損壞。
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

每個儲存格表示該項目在該帳號中的來源：共用群組（`default`）、位於帳號自己資料夾中時為 `account`，帳號沒有該項目時為 `-`。`!` 標示有問題的項目，說明列在表格下方。

## 讀取的內容

| 種類 | 位置 | 顯示的來源 |
| --- | --- | --- |
| 技能 | `skills/` 中每個含有 `SKILL.md` 的資料夾 | `skills/` 連結到的群組，或 `account` |
| 指令、代理 | `commands/` 和 `agents/` 底下的 `.md` 檔案，子資料夾顯示為 `folder/name` | 群組，或 `account` |
| 外掛 | `plugins/installed_plugins.json` | 群組，或 `account` |
| MCP 伺服器 | 群組的 `mcp.json`（共用 MCP 時）和帳號的 `.claude.json` | 群組，或 `account` |
| Hook | `settings.json` 中的 `hooks`，依事件計數 | 群組，或 `account` |
| 這個專案 | 目前資料夾中的 `.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json` 和 `.mcp.json` | 列出一次，適用於每個帳號 |

## 回報的問題

- **損壞的連結**：`skills/`、`commands/` 或 `agents/` 中目標已不存在的連結，或連結到不存在位置的共用資料夾（`skills`、`commands`、`agents`、`plugins`）（用 `ccam doctor --fix` 修復）。
- **沒有 `SKILL.md` 的技能資料夾**：Claude Code 會略過它。
- **在兩個地方定義不同的 MCP 伺服器**：群組的 `mcp.json` 和帳號的 `.claude.json`，或其中之一和專案的 `.mcp.json`。兩個地方的設定相同（`mcp_strategy = "sync"` 就會這樣寫入）並不是問題。用 `ccam mcp import <account> <server>` 把帳號自己的伺服器移到群組中。
- **不是有效 JSON 的 `installed_plugins.json`。**

## 安全性和特殊情況

- `ccam ext` 不會變更任何東西，也不會發出網路請求。它讀取伺服器名稱和 hook 數量，絕不讀取 MCP 定義中 `env` 或標頭的值。
- 它顯示的是磁碟上的內容。已安裝但在設定中關閉的外掛，或被受管理設定封鎖的伺服器，仍然會出現。
- 在專案資料夾中執行它，就會包含該專案的 `.claude` 資料夾和 `.mcp.json`。

## 相關內容

- [帳號之間的共用](/features/share)
- [共用的 MCP 伺服器](/features/mcp)
- [檢查和修復](/features/doctor)
