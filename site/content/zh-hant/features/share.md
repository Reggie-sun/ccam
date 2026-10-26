---
title: 用群組共用
description: ccam 群組透過一個資料夾，在其中的帳號之間共用選定的項目；ccam group 會依計畫移動資料，並提供備份和衝突選擇。
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "3f94dd13dbf1"
---

每個帳號都有自己的資料夾，但裡面大部分的內容並不需要分開。**群組**決定共用什麼、由誰共用：群組中的每個帳號都會把這些項目連結到同一個資料夾 `~/.local/share/ccam/groups/<group>/`，所以它們讀寫的都是同一份副本。登入權杖絕不會共用。

一個帳號可以屬於多個群組，但每個項目只會來自一個群組。常見的設定是一個群組用於設定和擴充功能，另一個群組用於兩個帳號之間的對話記錄：

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

你新增的第一個帳號會建立 `default`，之後的帳號除非你另外指定，否則都會加入它。

## 用法

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

每個會移動資料的指令都會印出計畫，針對每個和群組副本內容不同的檔案詢問你，並在變更任何東西之前詢問 `Proceed? [y/N]`：

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

你可以在[儀表板](/features/dashboard)的 **Groups** 分頁做同樣的變更；在帳號上按 `s` 可以選擇它的群組。

## 項目

| 項目 | 涵蓋的內容 |
| --- | --- |
| `mcp` | MCP 伺服器，透過群組的 `mcp.json`（[共用的 MCP 伺服器](/features/mcp)） |
| `settings` | `settings.json` |
| `memory` | 全域 `CLAUDE.md` |
| `extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| `history` | `projects/` 和 `history.jsonl`，讓另一個帳號可以 `--continue` |

如果帳號所屬的群組都沒有共用某個項目，該項目會留在帳號自己的資料夾中。

## 指令和選項

| 指令 | 作用 |
| --- | --- |
| `ccam group` 或 `ccam group list` | 列出群組 |
| `ccam group show <group>` | 群組共用的項目、它的帳號、資料夾和 MCP 伺服器 |
| `ccam group add <group> --items <items> [--members <accounts>]` | 建立群組；其中帳號的項目會移入群組 |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | 為群組中的每個帳號變更群組共用的項目 |
| `ccam group join <group> <account>...` | 把帳號加入群組 |
| `ccam group leave <group> <account>...` | 把帳號移出群組；每個帳號都會保留一份副本，MCP 伺服器也包括在內 |
| `ccam group rename <group> <new-name>` | 只在設定中重新命名 |
| `ccam group remove <group> [--leave]` | 移除群組 |

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--from <account>` | `add`、`set`：先把哪個帳號的副本放入群組 | 預設帳號，其次依編號 |
| `--replace` | `join`：先離開共用相同項目的群組 | 關閉 |
| `--fresh` | 停止共用某個項目的帳號從空白開始，而不是帶著副本（匯入的資料夾仍會保留一份） | 關閉 |
| `--leave` | `remove`：先把群組的帳號移出 | 群組中有帳號時必須指定 |
| `--prefer <group\|mine\|merge>` | 如何決定衝突 | 詢問；加上 `--yes` 時採用各檔案的預設值 |
| `-y`、`--yes` | 不詢問直接執行；衝突依 `--prefer` 決定 | 關閉 |
| `--dry-run`（全域） | 印出計畫，不變更任何東西 | 關閉 |

帳號可以用名稱或編號指定；項目則以逗號分隔的清單指定。`ccam add` 接受 `--group <group>`（可重複指定）或 `--no-group`。

## 運作方式

涉及的檔案：

- `~/.config/ccam/config.toml`：`[[group]]` 項目，包含 `name`、`items`、`members`，以及資料夾名稱不同時的 `dir`。
- `~/.local/share/ccam/profiles/<name>/`：帳號的 `CLAUDE_CONFIG_DIR`。
- `~/.local/share/ccam/groups/<folder>/`：共用的副本。

讀取 `config.toml` 時，每個帳號的群組會被解析成每個項目對應一個資料夾。五個帳號大約只需要 25 µs，所以啟動時仍然只讀取 `config.toml`。如果設定中有帳號會從兩個群組取得同一個項目，就是錯誤，在修正之前什麼都不會啟動。

變更是在設定的副本上進行的：

1. 指令編輯副本（`AddGroup`、`SetGroupItems`、`JoinGroup` …），並由這些函式檢查規則。
2. ccam 比較新舊設定，並為每個項目有所變更的帳號建立計畫（`internal/group`、`share.Build`）。計畫會逐一建立，每一份都看得到前面的計畫將放入群組資料夾的內容，所以只有第一個帳號會把它的副本移進去。群組第一份副本的來源依序是：群組資料夾（如果它已經有副本）、`--from` 指定的帳號、預設帳號，接著是依編號排列的其他帳號。
3. 如果其中有帳號正在執行（不論 Claude Code 是如何啟動的；``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``），或在沒有開發人員模式的 Windows 上有檔案需要連結，ccam 會拒絕。
4. ccam 解決衝突（詢問你，或依 `--prefer`），印出計畫，並請你確認。
5. 在儲存之前，ccam 會再檢查一次。如果有帳號在 ccam 等待你回答期間啟動，變更會停止；如果另一個 ccam 行程在這段期間變更了相同的群組，變更也會停止（`the groups changed while this plan waited; nothing was moved, make the change again`）。在移動任何東西之前，新的設定會先儲存，這次變更也會記錄在 `~/.local/state/ccam/group-change.json` 中。如果某個步驟中途停止，這份記錄會保留下來：`ccam doctor --fix` 會依原本的順序、以你給過的回答完成這項變更，包括 MCP 的步驟；在那之前，不會開始其他群組變更（``the last group change stopped halfway; finish it with `ccam doctor --fix` first``）。
6. MCP 的步驟最先執行。當帳號把擴充功能移出群組時，群組的外掛清單中經過帳號資料夾的路徑，會改為指向群組資料夾（計畫中的 `group plugin lists:`）。接著計畫依序執行。每個步驟都是把資料移到群組或備份資料夾；只有和群組副本完全相同的檔案會被刪除，而不是放到備份。

ccam 保有一份每個項目涵蓋路徑的允許清單（`share.Layout`）。不在清單上的東西永遠不會進入群組。ccam 也把 `.credentials.json` 和 `.claude.json` 列為絕不共用（`share.Never`），只要其中之一出現在結構清單中，就會以內部錯誤停止（`Build`）。`mcp` 項目沒有路徑：伺服器清單存放在 `.claude.json` 內，所以 [MCP 共用](/features/mcp)是透過群組的 `mcp.json` 運作的。

以下是 `side` 開始共用對話記錄、而群組已經有內容時，`projects/` 這個資料夾會發生的事：

<!-- motion:share-on -->

如果群組在該路徑還沒有任何內容，這個帳號就是第一個共用它的帳號，它的資料夾會原樣移過去。否則群組缺少的每個檔案都會移到群組中，兩邊都有的檔案則留在帳號的資料夾中，因此保留的是群組的版本。和群組內容相同的檔案會被刪除，計畫會列出內容不同的檔案（`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`）。

接著帳號的資料夾會移到帳號資料夾內的 `.ccam-backup/<time>/`，並由連結取代它的位置。如果群組的項目本身就是指向這個加入帳號自己項目的連結，帳號的內容會取代它放入群組中。對話檔案以工作階段 ID 命名，因此多個帳號的對話記錄可以合併而不會衝突。`plugins/` 的登錄檔儲存的是經過帳號資料夾的絕對路徑；ccam 會把它們改為指向群組的資料夾，而且當群組已經有登錄檔時，會把帳號的項目合併進去。

離開會反轉連結：

<!-- motion:share-off -->

ccam 會先把群組目前的內容複製到帳號資料夾中的 `.ccam-backup/incomplete-<name>`，並在複製完成後才取代連結。如果複製中途停止，連結會維持原狀，`ccam doctor --fix` 則會再複製一次（這種完成動作的計畫會顯示 `restore <entry>: … (copied back from the group before a stop)`）。群組會為其他帳號保留其內容。

- 只有當連結仍像建立計畫時那樣指向 groups 資料夾內時，ccam 才會移除連結。之後有所變更的路徑，例如原本是連結的地方變成了實際檔案，會維持原狀，變更也會停止。
- 加上 `--fresh` 時會略過複製。
- 帳號把某個項目移到資料夾仍是空的群組時，會帶著它看到的內容的副本過去，所以它的檔案都會保留。這份副本在完成之前會暫存在 `<group>/.ccam-incomplete-<name>` 中。
- 停止共用 MCP 的帳號會保留群組的伺服器：它們會寫入它自己的 `.claude.json`（會保留備份），但帳號已經有同名伺服器的除外，那些伺服器會維持原狀。加上 `--fresh` 時，計畫則會列出將不再載入的伺服器。再次加入時，沒有人改動過的副本會被移除。
- 離開 `extensions` 時，群組的外掛清單中經過離開帳號資料夾的路徑也會改為指向群組資料夾；對其他帳號而言，改寫後的路徑會指向相同的外掛。

兩邊都有但內容不同的單一檔案就是衝突：

<!-- motion:share-conflict -->

單一檔案（例如 `settings.json`）由你選擇；資料夾內則保留群組的副本，並由計畫列出內容不同的檔案。可用的選擇依檔案而定；第一個是預設值：

| 檔案 | 選擇 |
| --- | --- |
| `history.jsonl` | `merge by line`、`keep group`、`use mine` |
| 其他 `.json` 檔案 | `keep group`、`use mine`、`merge JSON` |
| 其他所有檔案 | `keep group`、`use mine` |

`merge JSON` 會遞迴合併物件、合併陣列並去除重複，當兩邊對同一個鍵設定了不同的值時保留群組的值，所以其他帳號看不到任何變化。除了 `keep group` 以外的每個選擇，都會先把群組的副本儲存到帳號的 `.ccam-backup/<time>/group/`。你自己的檔案一律會放到備份。

每次啟動時，ccam 只用 `lstat` 和 `readlink` 把結構清單中的每個路徑和帳號的群組比對（`share.Check`）。缺少的連結會當場建立。應該是連結的地方是實際檔案、連結指向錯誤的位置，或殘留過時的連結，都會中止啟動，並提示使用 `ccam doctor --fix`。

### 重新命名、移除和舊格式

- `ccam group rename` 只會變更設定。資料夾維持原名，群組會把它記錄為 `dir`，所以什麼都不會移動，也不會遺失任何東西。
- 除非加上 `--leave`，否則 `ccam group remove` 會拒絕移除仍有帳號的群組。加上之後，每個帳號都會保留一份副本，而資料夾會移到 `~/.local/share/ccam/trash/`，除非還有其他群組使用它。
- 在有群組之前，`config.toml` 會為每個帳號儲存一個 `[profile.share]` 表格。ccam 會把這種格式當作群組讀取，不移動任何東西。第一次儲存時，它會寫入新格式，並把舊檔案保留為 `config.toml.bak`。由相同帳號共用的項目會成為一個以資料夾命名的群組；其餘的則取得 `default-history` 之類的名稱，並繼續使用同一個資料夾。

### 在 Windows 上

Windows 只有在開發人員模式（設定 > 系統 > 開發人員專用）下或對系統管理員才會建立符號連結。沒有這項權限時，共用的資料夾會變成目錄連接點（directory junction），對 Claude Code 來說運作方式相同。共用的檔案（`settings.json`、`CLAUDE.md`、`history.jsonl`）只能是符號連結，所以需要這種連結的群組變更會在移動任何東西之前被拒絕。基於同樣的理由，在那裡第一次 `ccam add` 建立的 `default` 群組只包含 `mcp` 和 `extensions`，並會告訴你；以 `ccam add --group` 加入共用設定、記憶或對話記錄的群組則會被拒絕。

## 安全性和特殊情況

- **正在執行的帳號。** 當 Claude Code 正以變更所移動的其中一個帳號執行時，變更會被拒絕，因為工作階段會持續寫入這些檔案，不論它是否由 ccam 啟動（直接執行 `claude`、IDE、背景工作階段）。儲存前會再檢查一次。離開和加入只影響所指定的帳號；離開 `extensions` 時也會改寫群組外掛清單中的路徑，但不會改變它們對其他帳號的意義。`--dry-run` 仍會印出計畫。
- **備份。** 每份計畫在帳號資料夾中都有自己的 `.ccam-backup/<time>/` 資料夾。ccam 絕不刪除它們。
- **你自己的連結。** 如果某個路徑已經是你建立的連結（例如指向 dotfiles 儲存庫），而群組是空的，ccam 會把連結的目標複製到群組中，並把連結本身移到備份。目標不會被碰到。
- **匯入的資料夾。** 用 `ccam import` 在原位置登錄的帳號（例如 `~/.claude`）會在原位置變更，所以直接執行的 `claude` 也會看到變更；計畫會說明這一點。它離開群組時一定會保留一份副本，即使加上 `--fresh` 也一樣，而 `ccam remove` 在忘記這個帳號之前，也會以這種方式讓它離開群組。它在群組中時，共用的檔案位於群組資料夾中，所以在刪除 `~/.local/share/ccam` 之前，請先移除這個帳號（或離開群組）。
- **兩個帳號在同一個資料夾。** 共用對話記錄的帳號會寫入相同的檔案。當第二個帳號在已有帳號執行的資料夾中啟動時，ccam 會警告並詢問你。
- **手動編輯。** 如果 `members` 列出了不是帳號的名稱，ccam 會忽略它，並由 `ccam doctor` 回報。`local`、`on`、`off` 和 `none` 不能作為群組名稱。

## 相關內容

- [儀表板](/features/dashboard)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [共用的 MCP 伺服器](/features/mcp)
- [複製現有的設定](/features/migrate)
- [指令參考](/docs/commands)
