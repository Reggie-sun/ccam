---
title: 共用 MCP 伺服器、設定和對話記錄
description: ccam 群組如何在 Claude Code 帳號間共用 MCP 伺服器、settings.json、CLAUDE.md、skills 和對話記錄，登入則各自獨立。
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

每個帳號都保有自己的登入資訊。其他所有內容都可以透過**群組**共用：群組會指定它共用的項目和其中的帳號，ccam 會把這些帳號的檔案連結到群組的資料夾，因此在一個帳號中所做的變更，會出現在所有這些帳號中。

## 可以共用的項目

| 項目 | 共用的內容 |
| --- | --- |
| `mcp` | MCP 伺服器：群組的 `mcp.json`，在啟動時傳給 Claude Code |
| `settings` | `settings.json`：權限、hooks、模型、狀態列 |
| `memory` | 全域 `CLAUDE.md` |
| `extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| `history` | 對話記錄：`projects/` 和 `history.jsonl` |

你的第一個帳號會建立群組 `default`，它共用對話記錄以外的所有項目，之後的帳號都會加入它。

登入檔案 `.credentials.json` 和 `.claude.json` 永遠不會共用。ccam 只連結固定清單中的路徑，而這兩個檔案不在清單上。

本頁說明日常的操作。每個選項都在功能頁面中有詳細說明：[用群組共用](/features/share)、[共用的 MCP 伺服器](/features/mcp)和[複製現有的設定](/features/migrate)。

## 群組

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

一個帳號可以屬於多個群組，但每個項目只會來自一個群組。所以如果要讓 `work` 和 `side` 共用對話記錄，而不和你的其他帳號共用，就不要讓 `default` 共用對話記錄，並像上面那樣為這兩個帳號另建一個群組。`ccam share work` 會顯示帳號共用哪些項目，以及透過哪個群組共用。[儀表板](/features/dashboard)的 Groups 分頁也能做同樣的事。

## 帳號開始共用時

當帳號開始共用某個項目時（它加入群組，或它的群組新增了這個項目），ccam 會：

1. 為每個帳號印出一份計畫，列出它將移動、複製或連結的每個檔案，並在開始前徵求你的同意；
2. 如果群組仍是空的，就把帳號的檔案移到群組資料夾，否則只加入群組還沒有的檔案；
3. 遇到同一個檔案內容不同時，詢問要保留群組的版本、使用你的版本，還是合併 JSON；
4. 將你的原始檔案保留在帳號資料夾內的 `.ccam-backup/<time>/` 中；
5. 以指向群組的連結取代原始檔案。

對話檔案以工作階段 ID 命名，因此多個帳號的對話記錄可以合併而不會衝突。

當 Claude Code 正以某個帳號執行時（不論它是如何啟動的，背景工作階段也算在內），ccam 不會變更該帳號，因為工作階段會持續寫入正在移動的檔案。請先結束工作階段。ccam 在儲存前會再檢查一次，所以在計畫等待確認期間啟動的工作階段也會被發現。

## 帳號停止共用時

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

連結會被移除，帳號會取得一份群組目前內容的副本；群組的 MCP 伺服器會寫入它自己的 `.claude.json`，與它自己的伺服器並列。加上 `--fresh` 則改為從空白開始；此時計畫會列出將不再載入的 MCP 伺服器。`~/.claude` 這類匯入的資料夾一定會保留一份副本。群組會為其他帳號保留其內容。

## 複製現有的設定

如果你在使用 ccam 之前就用過 Claude Code，你的設定都在 `~/.claude` 中。`ccam migrate` 會把它們複製到一個群組，讓共用這些項目的每個帳號都能使用：

```sh
ccam migrate
```

在終端機中，ccam 會先顯示項目的核取清單。除了對話記錄以外的項目都已勾選；按 Space 切換項目，按 Enter 查看計畫。資料夾中沒有內容的項目會標示為 `not found`。

它會複製：

- `settings.json` 和全域 `CLAUDE.md`；
- 擴充功能 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/` 和 `rules/`；
- 你的使用者 MCP 伺服器，也就是 `~/.claude.json` 中的 `mcpServers` 項目，會複製到群組的 `mcp.json`；
- 對話記錄（`projects/` 和 `history.jsonl`），但只有在你加上 `--history` 時才會複製。

來源資料夾只會被讀取：其中的內容都不會被移動、連結或變更，直接執行 `claude` 也會和以前一樣運作。登入檔案 `.credentials.json` 和 `.claude.json` 絕不會被複製到群組。

使用 OAuth 登入的 MCP 伺服器會把權杖存放在登入資訊中，而不是伺服器清單中。ccam 不會複製這些登入：複製的登入持有相同的更新權杖，第一次更新時就會讓其他所有持有者（包括 `~/.claude`）從該伺服器登出。計畫和摘要會列出需要自行用 `/mcp` 登入這些伺服器的帳號。

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
Group default: personal
To share them:
  ccam group join default work
```

最後，ccam 會列出群組的帳號。如果群組還沒有共用其中某個項目，或有些帳號不在群組中，它會印出用來變更的 `ccam group` 指令。

如果群組中已有內容不同的檔案，ccam 會像開啟共用時一樣詢問你：保留群組中的副本，或使用你的副本。`settings.json` 等 JSON 檔案也可以合併，`history.jsonl` 則可以逐行合併。群組原本的副本會保存在群組資料夾內的 `.ccam-backup/<time>/`。

| 選項 | 作用 |
| --- | --- |
| `--settings`、`--memory`、`--extensions`、`--mcp`、`--history` | 只複製你指定的項目，不顯示核取清單；不指定且不在終端機中時，複製對話記錄以外的所有內容 |
| `-g`、`--group <group>` | 複製到這個群組，而不是 `default` |
| `-y`、`--yes` | 不詢問直接複製；檔案內容不同時保留群組中的副本 |

在指令前加上 `--dry-run`，例如 `ccam --dry-run migrate`，就能只查看計畫而不複製任何內容。若要從其他資料夾複製，請指定該資料夾：`ccam migrate ~/.claude-work --group work`。

## MCP 伺服器

Claude Code 將 MCP 伺服器和登入資訊一起存放在 `.claude.json` 中，所以這個檔案無法共用。ccam 將共用的伺服器保存在群組的 `mcp.json` 中，並以 `--mcp-config=<group>/mcp.json` 啟動 Claude Code。你自己的 `--mcp-config` 選項仍會疊加套用。

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

你在工作階段中用 `claude mcp add` 新增的伺服器只會儲存在該帳號，而且只適用於你當時所在的資料夾（Claude Code 的 local 範圍；加上 `-s user` 則會讓它在該帳號的每個資料夾中都可用）。下次啟動該帳號時，ccam 會列出這個伺服器的名稱與其資料夾，以及將它移到群組的指令：

```sh
ccam mcp import work gh
```

如果群組中已有同名的不同伺服器，ccam 會略過它並告訴你。`--force` 會取代群組的伺服器，並把舊的 `mcp.json` 保留在群組的 `.ccam-backup/` 資料夾中。

使用 OAuth 登入的 MCP 伺服器會將權杖和帳號的登入資訊存放在一起，所以請在每個帳號中用 `/mcp` 各授權一次。當某個帳號尚未登入某個伺服器、而另一個帳號已登入時，`ccam ext` 和 `ccam doctor` 會列出這些伺服器。

如果你希望這些伺服器以一般使用者伺服器的形式出現在 `claude mcp list` 和 `/mcp` 中，請在 `~/.config/ccam/config.toml` 的最上方、第一個 `[[group]]`、`[[profile]]` 或 `[[folder]]` 之前設定 `mcp_strategy = "sync"`。ccam 會在啟動前將它們寫入每個帳號的 `.claude.json`，過程中使用 Claude Code 自己的檔案鎖，並保留備份。帳號自己擁有的同名伺服器仍屬於該帳號，ccam 會告訴你。

## 兩個帳號在同一個資料夾

共用對話記錄的帳號會寫入相同的檔案。如果你在某個資料夾中啟動第二個 Claude Code，而同一群組的另一個帳號已經在該資料夾中執行，ccam 會警告你並詢問是否繼續。

## 發生問題時

```sh
ccam doctor
ccam doctor --fix
```

如果某個共用路徑變成了一般檔案（例如你手動還原了備份），ccam 會在啟動 Claude Code 前停止，並提示你使用 `ccam doctor --fix`；這個指令會重新將它連結到群組：群組沒有副本時，該檔案會成為群組的副本；否則群組保留自己的副本，你的檔案則放入備份，`--fix` 會指出備份的位置。中途停止的群組變更也由 `ccam doctor --fix` 完成。
