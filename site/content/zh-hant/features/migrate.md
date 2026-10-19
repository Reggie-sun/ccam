---
title: 複製現有的設定
description: 用 ccam migrate 把現有 ~/.claude 中的設定、CLAUDE.md、擴充功能、MCP 伺服器和 MCP 登入複製到 ccam 的共用群組。
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

如果你在使用 ccam 之前就用過 Claude Code，你的設定都在 `~/.claude` 和 `~/.claude.json` 中。`ccam migrate` 會把它們複製到一個共用群組，讓透過該群組共用這些項目的每個帳號，一開始就擁有你現有的設定。來源資料夾只會被讀取，直接執行 `claude` 也會和以前一樣運作。請在第一次 `ccam add` 之後執行一次，並在每次有帳號開始共用 MCP 伺服器時再執行一次。

## 用法

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

在終端機中，沒有指定項目旗標也沒有 `--yes` 時，ccam 會先顯示一份核取清單：

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space（或 `x`）切換游標所在的項目，方向鍵或 `j`/`k` 移動游標，Enter 顯示計畫，Esc、`q` 或 Ctrl+C 會印出 `Cancelled.`。來源中沒有內容的項目會標示為 `not found`。如果你取消勾選所有項目，ccam 會印出 `Nothing selected.` 並停止。

接著 ccam 會印出計畫，並在複製前詢問你：

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

群組只會送達透過它共用某個項目的帳號，所以摘要會逐項列出這些帳號。對於尚未共用的帳號，ccam 會印出開啟共用的 `ccam share <name> --<item> on` 指令。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `[folder]` | 要從中複製的 Claude Code 資料夾 | `~/.claude` |
| `--settings` | 複製 `settings.json` | 關閉 |
| `--memory` | 複製全域 `CLAUDE.md` | 關閉 |
| `--extensions` | 複製 `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 關閉 |
| `--mcp` | 複製使用者 MCP 伺服器（`.claude.json` 中的 `mcpServers`） | 關閉 |
| `--history` | 複製對話記錄（`projects/`、`history.jsonl`） | 關閉 |
| `-g`、`--group string` | 要複製到的共用群組 | `default` |
| `-y`、`--yes` | 不詢問直接複製；發生衝突時保留群組的副本 | 關閉 |
| `--dry-run`（全域） | 印出計畫，並在任何詢問或複製之前停止 | 關閉 |

沒有指定項目旗標時，在終端機中由核取清單決定。不在終端機中或加上 `--yes` 時，ccam 會複製對話記錄以外的所有內容。對話記錄預設不包含在內，因為 `projects/` 可能很大。

`--group on` 代表 `default`。`--group off` 和 `--group local` 會被拒絕，並顯示 `--group takes a group name, not "off"`。`ccam migrate` 支援 `--dry-run`，而且試執行時仍會顯示核取清單，讓你預覽所選的項目。

## 運作方式

`internal/cli/migrate.go` 會選出項目，接著 `migrate.Build` 讀取來源和群組並傳回一份計畫。在你確認、執行 `Plan.Apply` 之前，不會寫入任何東西。

1. **找出來源。** `checkSource` 會拒絕不存在、不是資料夾，或沒有任何已知 Claude Code 項目的路徑。`migrate.ClaudeJSON` 會找出屬於該資料夾的 `.claude.json`：在資料夾內，或者對於內部沒有這個檔案的 `~/.claude`，則是它旁邊的 `~/.claude.json`。沒有設定 `CLAUDE_CONFIG_DIR` 時，Claude Code 就把它放在那裡。
2. **把項目對應到路徑。** 每個路徑都來自 `share.Layout`，也就是 `ccam share` 使用的同一份允許清單。MCP 伺服器是用 `mcp.LocalServers` 從 `.claude.json` 讀出的，它使用 `profile.TopLevelValue` 掃描器，而不解析整個檔案。
3. **只複製，絕不移動。** 檔案以 `fsutil.AtomicWrite` 或 `fsutil.CopyFile` 寫到 `~/.local/share/ccam/groups/<group>/`；MCP 伺服器則寫入群組的 `mcp.json`。

每個項目都把一個來源路徑對應到一個群組路徑：

<!-- motion:migrate-items -->

兩列對話記錄是虛線，因為只有加上 `--history` 時才會複製。`.credentials.json` 和 `.claude.json` 列在 `share.Never` 中，絕不會進入群組：如果其中之一出現在結構清單中，`Build` 會傳回內部錯誤，走訪資料夾時也會依名稱略過它們。

`buildEntry` 會給每個路徑一個狀態，計畫中的每一行描述的就是這個狀態：

<!-- motion:migrate-states -->

- **Missing**：來源中沒有，或在應該是資料夾的地方是檔案。每個項目只列出一次，標示為 `not found, skipped`。
- **Linked**：來源路徑已經解析到群組的副本，例如來源是 ccam 的帳號資料夾時。
- **Same**：群組已經有相同的內容。
- **New**：群組沒有，所以整個複製過去。
- **Merge**：兩邊都有的資料夾。`copyMissing` 只加入群組沒有的檔案；內容不同的檔案維持群組中的版本，計畫會顯示 `group copies kept`。對於 `plugins/`，登錄檔 `installed_plugins.json` 和 `known_marketplaces.json` 會被合併，外掛的絕對路徑會改寫為指向群組的資料夾。
- **Conflict**：兩邊都有但內容不同的檔案。你可以逐檔選擇：`keep group`、`use mine`，`.json` 檔案另外有 `merge JSON`（深層合併，值衝突時以群組為準），`history.jsonl` 則有 `merge by line`。定義不同的 MCP 伺服器也算衝突，可選擇 `keep group` 或 `use mine`。列出的第一個選擇是按下 Enter 時的預設值；`--yes` 一律保留群組的副本。

核取清單本身是 `tui.MigratePicker`，建立在 `ccam add` 也使用的 `checklist` 元件上；`migrate.Found` 透過 lstat 呼叫和 MCP 伺服器清單標示缺少的項目。

### MCP 登入

使用 OAuth 登入的 MCP 伺服器會把權杖存放在登入檔案中，也就是 `.credentials.json` 的 `mcpOAuth` 底下，而不是伺服器清單中。Claude Code 以 `<name>|` 加上伺服器 `type`、`url` 和 `headers` 的 `sha256` 雜湊值前 16 個字元作為每個項目的鍵。這個鍵與資料夾無關，所以只要伺服器定義完全相同，複製過去的項目在其他帳號中也能使用。因此 `internal/migrate/logins.go` 中的 `buildLogins` 只會複製群組最終採用來源定義的那些伺服器的項目：新的、原本就相同的，或你在衝突中選擇了 `use mine` 的。外掛伺服器的登入（鍵以 `plugin:` 開頭）會在你複製擴充功能且來源有 `plugins/` 時一併複製。

<!-- motion:migrate-signins -->

這些項目會寫入每個接收帳號自己的 `.credentials.json`，以 `fsutil.AtomicWrite` 和權限 `0600` 寫入。它們絕不會進入群組。帳號只有在透過目標群組共用 `mcp`（外掛伺服器則是 `extensions`）時才會收到，而且只加入它缺少的鍵：它已經有的項目（可能是較新的）會保留。計畫中的 `sign-ins` 那一行會列出伺服器、接收的帳號和被略過的帳號。符合以下情況的帳號會被 ccam 略過，並列出名稱和原因：

- 以 `ccam import` 就地登錄的帳號，因為直接執行的 `claude` 可能正在使用它的登入檔案；
- 正在執行的帳號，因為 Claude Code 會改寫它的登入檔案（寫入前會再檢查一次，回報為 `started meanwhile`）；
- 尚未登入的帳號；
- 把登入資訊存放在 macOS 鑰匙圈中的帳號，ccam 不會寫入鑰匙圈。

請在被略過的帳號中用 `/mcp` 登入這些伺服器。之後才開始共用 MCP 伺服器的帳號，會在你再次執行 `ccam migrate` 時取得這些登入；摘要會以 ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` 提醒你。

## 安全性和特殊情況

- 來源只會被讀取。其中的檔案都不會被移動、連結或變更，來源的登入資訊也絕不會被寫入，包括在 macOS 上可能從鑰匙圈讀取的情況。
- 在群組檔案因 `use mine`、合併或現有 `mcp.json` 中新增 MCP 伺服器而被取代之前，群組的副本會儲存在 `<group>/.ccam-backup/<time>/`。資料夾名稱每次執行都不同，所以同一秒內執行兩次也不會覆蓋彼此的備份。摘要會印出這個路徑。
- 加上 `--yes` 時，每個衝突都會保留群組的副本，ccam 也會告訴你如何改為自行選擇：不加 `--yes` 執行同一個指令。
- 如果沒有任何變更，ccam 會印出 `Nothing to copy.`；如果每個衝突都保留群組的副本，則印出 `Nothing to copy: the group keeps its own copies.`
- 如果複製中途失敗，ccam 會停止並印出 `The source folder was not changed. Group files replaced so far are saved in`，後面接著備份資料夾。
- 來源不存在時，會以 `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` 結束。

## 相關內容

- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [檢查和修復](/features/doctor)
- [多個 Claude Code 帳號](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
