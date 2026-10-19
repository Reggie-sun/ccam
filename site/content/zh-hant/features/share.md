---
title: 帳號之間的共用
description: ccam share 如何把設定、CLAUDE.md、擴充功能和對話記錄連結到群組資料夾，並提供計畫、備份和衝突選擇，以及 off 如何複製回來。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

每個帳號都有自己的資料夾，但裡面大部分的內容並不需要分開。`ccam share` 一次針對一個帳號的一個項目開啟或關閉共用：帳號的檔案或資料夾會被移到群組資料夾，並以符號連結取代，所以群組中的每個帳號讀寫的都是同一份副本。在 `ccam add` 之後用它來變更帳號共用的項目，或讓幾個帳號組成自己的群組。

## 用法

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

不加項目選項時，`ccam share` 會印出目前的狀態：

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

加上選項時，ccam 會印出計畫，針對每個和群組副本內容不同的檔案詢問你，並在變更任何東西之前詢問 `Proceed? [y/N]`。新帳號會共用對話記錄以外的所有項目；[共用 MCP 伺服器、設定和對話記錄](/docs/sharing)從使用者的角度說明每個項目。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP 伺服器 | 不變 |
| `--settings <on\|off\|group>` | `settings.json` | 不變 |
| `--memory <on\|off\|group>` | 全域 `CLAUDE.md` | 不變 |
| `--extensions <on\|off\|group>` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 不變 |
| `--history <on\|off\|group>` | `projects/` 和 `history.jsonl` | 不變 |
| `--fresh` | 關閉共用時從空白開始，而不是複製群組的內容 | 關閉 |
| `--prefer <group\|mine\|merge>` | 如何決定衝突 | 詢問；加上 `--yes` 時採用各項目的預設值 |
| `-y`、`--yes` | 不詢問直接執行；衝突依 `--prefer` 決定 | 關閉 |
| `-h`、`--help` | `share` 的說明 | |
| `--dry-run`（全域） | 支援：會詢問衝突並印出計畫，但不變更任何東西 | 關閉 |

值 `on` 代表群組 `default`，`off` 代表 local（不共用），其他任何值都是群組名稱。群組名稱會成為 `~/.local/share/ccam/groups/` 下的資料夾名稱，所以 `ccam share side --history ../x` 會失敗並顯示 `invalid group name: "../x"`。帳號可以用名稱或編號指定，只有你指定的項目會變更。

## 運作方式

涉及的檔案：

- `~/.config/ccam/config.toml`：每個帳號的 `[profile.share]` 表格，每個項目一個值。
- `~/.local/share/ccam/profiles/<name>/`：帳號的 `CLAUDE_CONFIG_DIR`。
- `~/.local/share/ccam/groups/<group>/`：共用的副本。

1. `internal/cli/share.go` 讀取設定，用 `config.Share.Set` 把每個 `--<item>` 的值套用到帳號設定的副本上，如果帳號正在執行就拒絕。
2. `share.Build` 比較目前的設定和目標設定，並列出步驟。它只讀取磁碟：對 `share.Layout` 中的每個路徑執行 `os.Lstat` 和 `os.Readlink`。計畫是一連串操作：移入群組中的空位（adopt）、移入群組缺少的檔案（merge dir）、把已決定的衝突結果寫入群組、備份、取消連結、建立連結、複製回來、改寫外掛路徑，以及合併外掛登錄檔。
3. 衝突由提示或 `--prefer` 決定。
4. `Plan.Describe` 印出計畫，然後 ccam 詢問你。
5. 在移動任何東西之前，目標設定會先寫入 `config.toml`。如果 `Plan.Apply` 中途停止，`ccam doctor --fix` 和重新執行都會朝已儲存的目標建立計畫，完成這項變更。
6. `Plan.Apply` 依序執行各步驟。每個步驟都是把資料移到群組或備份資料夾；不會刪除任何東西。

`share.Layout` 是每個項目涵蓋路徑的允許清單。不在清單上的東西永遠不會進入群組。`share.Never` 列出 `.credentials.json` 和 `.claude.json`，只要其中之一出現在結構清單中，`Build` 就會以內部錯誤停止。`mcp` 項目沒有路徑：伺服器清單存放在 `.claude.json` 內，所以 [MCP 共用](/features/mcp)改為透過群組的 `mcp.json` 運作，`--mcp` 只會變更設定。

以下是 `side` 開啟對話記錄共用、而群組已經有內容時，`projects/` 這個資料夾會發生的事：

<!-- motion:share-on -->

如果群組在該路徑還沒有任何內容，這個帳號就是第一個共用它的帳號，它的資料夾會原樣移過去，不做備份。否則 `mergeDir` 會把群組缺少的每個檔案移到群組中，兩邊都有的檔案則留在帳號的資料夾中，因此保留的是群組的版本。接著帳號的資料夾會移到帳號資料夾內的 `.ccam-backup/<time>/`，再由 `os.Symlink` 在原位置建立連結。對話檔案以工作階段 ID 命名，因此多個帳號的對話記錄可以合併而不會衝突。

移動使用 `rename(2)`，跨檔案系統時則改為複製後刪除。被移動的目錄樹中，指向樹外部的相對連結會改寫為絕對連結（`fsutil.PortableLink`），所以從新位置仍能解析。`plugins/` 的登錄檔（`installed_plugins.json`、`known_marketplaces.json`）儲存的是經過帳號資料夾的絕對路徑；ccam 會把它們改為指向群組的資料夾，而且當群組已經有登錄檔時，會把帳號的項目合併進去。

關閉某個項目會反轉連結：

<!-- motion:share-off -->

ccam 只在連結指向 groups 資料夾內時才移除它，然後用 `CopyTree` 把群組目前的內容複製到帳號中。加上 `--fresh` 時會略過複製，帳號在該路徑從空白開始。群組會為其他帳號保留其內容。把項目從一個群組移到另一個群組時，會取消連結並連結到新群組，舊群組維持原狀；計畫會以 `Note:` 這一行說明。

兩邊都有但內容不同的單一檔案就是衝突：

<!-- motion:share-conflict -->

`share.NewConflict` 會列出適合該檔案的選擇。`history.jsonl` 提供 `merge by line`、`keep group` 和 `use mine`，合併排在第一個。其他 `.json` 檔案提供 `keep group`、`use mine` 和 `merge JSON`。其他所有檔案提供 `keep group` 和 `use mine`。第一個選擇是預設值：

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON`（`MergeJSONBytes`）會遞迴合併物件、合併陣列並去除重複，當兩邊對同一個鍵設定了不同的值時保留群組的值，所以其他帳號的行為不會改變。`merge by line` 會保留群組的各行，再加上群組沒有的你的各行。除了 `keep group` 以外的每個選擇，都會先把群組的副本儲存到帳號資料夾中的 `.ccam-backup/<time>/group/`，再以 `fsutil.AtomicWrite` 寫入結果。你自己的檔案一律會放到備份。

`--prefer group`、`mine` 或 `merge` 會不詢問直接決定每個衝突；`merge` 代表 `merge JSON` 或 `merge by line`，視檔案提供哪一個而定，兩者都不提供的檔案則維持預設值。只加 `--yes` 時會採用每個檔案的預設值。

每次啟動時，`share.Check` 只用 `lstat` 和 `readlink` 把結構清單中的每個路徑和設定比對，所以成本低到足以放在啟動流程中。缺少的連結會當場建立，因為沒有任何東西會遺失。應該是連結的地方是實際檔案、連結指向錯誤的位置，或殘留過時的連結，都會中止啟動：

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` 會以已儲存的設定同時作為兩邊執行 `share.Build`，根據磁碟上的實際內容推導每個步驟，發生衝突時保留群組的內容，並備份所有被取代的東西。

## 安全性和特殊情況

- **正在執行的帳號。** 當 Claude Code 正以該帳號執行時，ccam 會拒絕變更共用（`account work is running; quit it first, then change its sharing`），因為工作階段會持續寫入正在移動的檔案，期間附加的各行會遺失。`--dry-run` 仍會印出計畫。
- **備份。** 每份計畫在帳號資料夾中都有自己的 `.ccam-backup/<time>/` 資料夾，兩份計畫在同一秒開始時會加上編號後綴。ccam 絕不刪除它們。
- **部分失敗。** 如果某個步驟失敗，ccam 會印出錯誤和原始檔案的位置，並告訴你用 `ccam doctor --fix` 完成。目標已經儲存，所以修復會朝它繼續進行。
- **你自己的連結。** 如果某個路徑已經是你建立的連結（例如指向 dotfiles 儲存庫），而群組是空的，ccam 會把連結的目標複製到群組中，並把連結本身移到備份。目標不會被碰到。
- **登入檔案。** `.credentials.json` 和 `.claude.json` 絕不會被連結、移動或複製到群組中。
- **兩個帳號在同一個資料夾。** 共用對話記錄的帳號會寫入相同的檔案。當同一個對話記錄群組的第二個帳號在已有帳號執行的資料夾中啟動時，ccam 會警告並在繼續之前詢問你。

## 相關內容

- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [共用的 MCP 伺服器](/features/mcp)
- [複製現有的設定](/features/migrate)
- [所有指令](/docs/commands)
