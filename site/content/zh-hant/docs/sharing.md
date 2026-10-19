---
title: 共用 MCP 伺服器、設定和對話記錄
description: ccam 如何在 Claude Code 帳號間共用 MCP 伺服器、settings.json、CLAUDE.md、skills 和對話記錄，登入則各自獨立。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

每個帳號都保有自己的登入資訊。其他所有內容都可以逐項共用：ccam 會將帳號的檔案或資料夾連結到共用的群組資料夾，因此在一個帳號中所做的變更，會出現在該群組的所有帳號中。

## 可以共用的項目

| 項目 | 選項 | 共用的內容 |
| --- | --- | --- |
| MCP 伺服器 | `--mcp` | 群組的 `mcp.json`，在啟動時傳給 Claude Code |
| 設定 | `--settings` | `settings.json`：權限、hooks、模型、狀態列 |
| 全域記憶 | `--memory` | `CLAUDE.md` |
| 擴充功能 | `--extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| 對話記錄 | `--history` | `projects/` 和 `history.jsonl` |

新帳號會共用對話記錄以外的所有項目。

登入檔案 `.credentials.json` 和 `.claude.json` 永遠不會共用。ccam 只連結固定清單中的路徑，而這兩個檔案不在清單上。

## 群組

`on` 代表 `default` 群組。你也可以指定名稱，讓只有部分帳號共用某個項目：

```sh
ccam share work --history team
ccam share side --history team
```

現在 `work` 和 `side` 彼此共用對話，但不與你的其他帳號共用。

## 開啟共用

當帳號開始共用某個項目時，ccam 會：

1. 印出一份計畫，列出它將移動、複製或連結的每個檔案，並在開始前徵求你的同意；
2. 如果群組仍是空的，就把帳號的檔案移到群組資料夾，否則只加入群組還沒有的檔案；
3. 遇到同一個檔案內容不同時，詢問要保留群組的版本、使用你的版本，還是合併 JSON；
4. 將你的原始檔案保留在帳號資料夾內的 `.ccam-backup/<time>/` 中；
5. 以指向群組的連結取代原始檔案。

對話檔案以工作階段 ID 命名，因此多個帳號的對話記錄可以合併而不會衝突。

## 關閉共用

```sh
ccam share work --history off
```

連結會被移除，帳號會取得一份群組目前內容的副本。加上 `--fresh` 則改為從空白開始。群組會為其他帳號保留其內容。

## 複製現有的設定

如果你在使用 ccam 之前就用過 Claude Code，你的設定都在 `~/.claude` 中。`ccam migrate` 會把它們複製到一個群組，讓共用這些項目的每個帳號都能使用：

```sh
ccam migrate
```

它會複製：

- `settings.json` 和全域 `CLAUDE.md`；
- 擴充功能 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/` 和 `rules/`；
- 你的使用者 MCP 伺服器，也就是 `~/.claude.json` 中的 `mcpServers` 項目，會複製到群組的 `mcp.json`；
- 對話記錄（`projects/` 和 `history.jsonl`），但只有在你加上 `--history` 時才會複製。

來源資料夾只會被讀取：其中的內容都不會被移動、連結或變更，直接執行 `claude` 也會和以前一樣運作。登入檔案 `.credentials.json` 和 `.claude.json` 絕不會被複製。

ccam 在複製任何內容之前，會先顯示計畫並詢問你：

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

最後，ccam 會列出透過該群組共用每個項目的帳號。如果有帳號尚未共用，它會顯示為這些帳號開啟共用的 `ccam share` 指令。

如果群組中已有內容不同的檔案，ccam 會像開啟共用時一樣詢問你：保留群組中的副本，或使用你的副本。`settings.json` 等 JSON 檔案也可以合併，`history.jsonl` 則可以逐行合併。群組原本的副本會保存在群組資料夾內的 `.ccam-backup/<time>/`。

| 選項 | 作用 |
| --- | --- |
| `--settings`、`--memory`、`--extensions`、`--mcp`、`--history` | 只複製你指定的項目；不指定時，複製對話記錄以外的所有內容 |
| `-g`、`--group <group>` | 複製到這個群組，而不是 `default` |
| `-y`、`--yes` | 不詢問直接複製；檔案內容不同時保留群組中的副本 |

在指令前加上 `--dry-run`，例如 `ccam --dry-run migrate`，就能只查看計畫而不複製任何內容。若要從其他資料夾複製，請指定該資料夾：`ccam migrate ~/.claude-work --group work`。

## MCP 伺服器

Claude Code 將 MCP 伺服器和登入資訊一起存放在 `.claude.json` 中，所以這個檔案無法共用。ccam 將共用的伺服器保存在群組的 `mcp.json` 中，並以 `--mcp-config=<group>/mcp.json` 啟動 Claude Code。你自己的 `--mcp-config` 選項仍會疊加套用。

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

你在工作階段中用 `claude mcp add` 新增的伺服器只會儲存在該帳號。下次啟動該帳號時，ccam 會列出這個伺服器的名稱，以及將它移到群組的指令：

```sh
ccam mcp import work gh
```

使用 OAuth 登入的 MCP 伺服器會將權杖和帳號的登入資訊存放在一起，所以請在每個帳號中各授權一次。

如果你希望這些伺服器以一般使用者伺服器的形式出現在 `claude mcp list` 和 `/mcp` 中，請在 `~/.config/ccam/config.toml` 中設定 `mcp_strategy = "sync"`。ccam 會在啟動前將它們寫入每個帳號的 `.claude.json`，過程中使用 Claude Code 自己的檔案鎖，並保留備份。

## 兩個帳號在同一個資料夾

共用對話記錄的帳號會寫入相同的檔案。如果你在某個資料夾中啟動第二個 Claude Code，而同一群組的另一個帳號已經在該資料夾中執行，ccam 會警告你並詢問是否繼續。

## 發生問題時

```sh
ccam doctor
ccam doctor --fix
```

如果某個共用路徑變成了一般檔案（例如你手動還原了備份），ccam 會在啟動 Claude Code 前停止，並提示你使用 `ccam doctor --fix`；這個指令會將該檔案合併到群組中，並保留備份。
