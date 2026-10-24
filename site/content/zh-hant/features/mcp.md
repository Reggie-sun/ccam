---
title: 共用的 MCP 伺服器
description: 用 ccam mcp add、add-json、list、remove 和 import 為每個群組維護一份 MCP 伺服器清單，透過 --mcp-config 傳入或同步到每個帳號。
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code 把使用者 MCP 伺服器存放在 `.claude.json` 中，而這個檔案也存放著帳號的登入資訊，所以無法在帳號之間共用。ccam 把共用的伺服器保存在群組檔案 `mcp.json` 中，並在啟動時交給群組中的每個帳號。用 `ccam mcp` 為所有帳號一次新增伺服器，用 `ccam mcp import` 把你在某個工作階段中新增的伺服器移給其他帳號。

## 用法

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

變更會從每個帳號下次啟動時開始生效；正在執行的工作階段會繼續使用啟動時的伺服器：

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

在群組下方，`ccam mcp list` 會列出每個不會送達其他帳號的伺服器：帳號為自己新增的伺服器顯示為 `only in account work: pg`，為某個資料夾新增的伺服器（一般的 `claude mcp add`，見下文）顯示為 `only in account side, for /home/me/code/app only: notes`，帳號自己的伺服器與群組伺服器同名時顯示為 `account work has its own …, used in its sessions in place of the group's`，這個資料夾的 `.mcp.json` 會附上其核准仍處於 `waiting for approval in` 的帳號，而 `in no account, only in ~/.claude.json (plain claude's)` 的伺服器則附上可複製它們的 [`ccam migrate --mcp`](/features/migrate) 指令。

新增群組中已有的名稱時，會取代它的定義，並印出 `updated` 而不是 `added`。只有在該群組中、而且該群組共用 MCP 時（`ccam group set <group> --add mcp`），帳號才會取得這些伺服器；如果群組沒有共用 MCP，`add` 會告訴你。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `-g`、`--group <group>` | `add`、`add-json` 和 `remove` 的目標群組 | 有帳號用來共用 MCP 的唯一群組，否則為 `default` |
| `-- <command> [args…]`（`add`） | 以 stdio 伺服器的形式執行這個指令 | |
| `--url <url>`（`add`） | 遠端伺服器的 URL（http/sse） | |
| `--transport <stdio\|http\|sse>`（`add`） | `--url` 伺服器的傳輸方式 | `http` |
| `--header <'Name: value'>`（`add`） | HTTP 標頭，可重複指定 | 無 |
| `-e`、`--env <KEY=VALUE>`（`add`） | stdio 伺服器的環境變數，可重複指定 | 無 |
| `--force`（`import`） | 取代群組中同名但定義不同的伺服器；舊的 `mcp.json` 會先備份 | 關閉 |
| `--keep-local`（`import`） | 保留帳號 `.claude.json` 中的項目 | 關閉 |
| `-h`、`--help` | 各指令的說明 | |
| `--dry-run`（全域） | 只有 `mcp list` 接受；`add`、`add-json`、`remove` 和 `import` 會拒絕它，且不變更任何東西 | 關閉 |

`add` 需要 `--` 後面的指令或 `--url` 其中之一，否則會失敗並顯示 `a command to run (after --) or --url is required`。`--header` 適用於 `--url` 伺服器，`--env` 適用於指令伺服器。`add-json` 接受和 `claude mcp add-json` 相同的 JSON，定義必須是含有 `command` 或 `url` 的物件。`import` 一律使用帳號用來共用 MCP 的群組，`list` 則顯示所有群組，所以 `-g` 不會改變它們。

## 運作方式

涉及的檔案：

- `~/.local/share/ccam/groups/<group>/mcp.json`：共用的伺服器，格式為 `{"mcpServers": {…}}`，也就是 `--mcp-config` 讀取的格式。它以 `fsutil.AtomicWrite` 和權限 `0600` 寫入，因為伺服器項目常常帶有 API 金鑰，而且只在持有鎖 `<group>/.mcp.json.lock` 時寫入，所以同時執行兩個 `ccam mcp add` 時，兩者的伺服器都會保留。
- `~/.local/share/ccam/profiles/<name>/.claude.json`：Claude Code 自己的檔案，存放帳號的伺服器。ccam 只用 `profile.TopLevelValue` 掃描器讀取其中的 `mcpServers` 值，在啟動流程中絕不解析整個檔案。
- `~/.local/state/ccam/mcp/<name>/`：`seen` 檔案、`sync-stamp`（附有 ccam 寫入的每份副本的雜湊值，讓它分辨沒有人改過的副本和你編輯過的副本）、`sync-checked`（sync 策略上次看到的內容）、`kept`（帳號退出 MCP 共用時保留的伺服器），以及 `.claude.json` 的備份。

有兩種策略，由 `~/.config/ccam/config.toml` 中的 `mcp_strategy` 設定。預設是 flag 策略：

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` 找出帳號的 MCP 群組。`mcp.FlagArgs` 以一次 `stat` 確認群組的 `mcp.json` 存在，而不解析它。
2. `Plan.InsertArgs` 把 `--mcp-config=<path>` 放在 `claude` 的正後方、你自己的參數之前，所以你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍會疊加套用。
3. `mcp.NewLocal` 檢查帳號的 `.claude.json` 自上次啟動後是否變更，只有變更時才解析其中的伺服器。
4. `syscall.Exec` 以 Claude Code 取代 ccam。

路徑以單一參數 `--mcp-config=<path>` 傳入，而不是兩個參數。Claude Code 的 `--mcp-config` 接受多個值，所以 `claude --mcp-config a.json "prompt"` 會把提示詞當成第二個設定檔讀取，並以 `MCP config file not found` 失敗。單一參數的寫法會在那裡結束，後面接的提示詞或子指令都能正常運作。`ccam work mcp list`、`get` 和 `remove` 不會加上這個旗標，因為 Claude Code 在那裡會拒絕它；ccam 會改為在標準錯誤中列出群組的伺服器（[詳細說明](/features/launch#shared-mcp-servers)）。Claude Code 會在帳號自己的伺服器之外再載入 `--mcp-config` 的伺服器；名稱衝突時，`.claude.json` 中的項目優先。flag 策略只在以下情況寫入 `.claude.json`：Remote Control、帳號退出或重新加入 MCP 共用時（[見下文](#leaving-and-joining)），以及 `ccam mcp import`。

`claude remote-control`（`rc`）不接受 `--mcp-config`，它的工作階段只從 `.claude.json` 讀取伺服器。對於 `ccam work rc`，ccam 不加上這個參數，並在終端機上詢問是否要像 `sync` 策略那樣，先建立備份再把群組的伺服器複製到 `.claude.json`。在副本存在之前它會一直詢問，之後只有在群組的 `mcp.json` 變更，或某個已複製的伺服器消失後，才會再次詢問（[Remote Control](/features/launch#remote-control)）。由於名稱衝突時這種副本優先於 `--mcp-config`，`ccam mcp add`、`add-json`、`remove` 和 `import` 會更新群組中其他帳號持有的副本：

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

只有在副本有變更時才會出現 `✓ Updated` 這一行。較舊版 ccam 建立的副本也會以同樣的方式更新，`ccam migrate` 之後的副本也一樣。

在工作階段中用 `claude mcp add` 新增的伺服器只會存到該帳號的 `.claude.json`。一般的 `claude mcp add` 使用 local 範圍：該伺服器只會在那一個資料夾中為那一個帳號載入（加上 `-s user` 則會讓它在該帳號的每個資料夾中都可用）。當下 ccam 看不到它，因為 Claude Code 已經取代了 ccam，所以它會在下次啟動時檢查。`seen` 檔案記錄 `.claude.json` 的修改時間和已經回報過的伺服器名稱。時間不同時，ccam 會把名稱和群組比對，並為每個新伺服器印出一次提示，不會等待輸入，兩種策略都一樣：

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

這個提示也會列出帳號自己擁有、與群組伺服器同名的伺服器，因為在該帳號的工作階段中，帳號的伺服器優先。

使用 `mcp_strategy = "sync"` 時，這些伺服器會以一般使用者伺服器的形式出現在 `claude mcp list` 和 `/mcp` 中，ccam 會在啟動前把它們寫入每個帳號的 `.claude.json`：

<!-- motion:mcp-sync -->

1. `mcp.Sync` 把群組 `mcp.json` 的修改時間和 `sync-stamp` 比對，並把帳號的 `.claude.json` 和它上次看到的內容比對。如果兩者都沒有變更，就不做其他事。否則它還會確認每個群組伺服器都仍在 `.claude.json` 中，並把消失的副本放回去（`put back group default's MCP servers that were gone from account work's .claude.json: ctx`）。
2. 它取得 Claude Code 自己的鎖，也就是檔案旁邊以 `mkdir` 建立的資料夾 `.claude.json.lock`。它最多等待 2 秒，並和 Claude Code 的鎖函式庫一樣，把超過 10 秒的鎖視為已失效。
3. 它把目前的 `.claude.json` 儲存到狀態資料夾中的 `backups/`，並保留最新的 5 份。
4. 它在 `mcpServers` 中設定每個群組伺服器，移除它之前寫入、但群組已經沒有的伺服器，並只在有變更時以檔案原本的權限原子寫入。它只變更或移除自己寫入、且之後沒有人改動過的副本。帳號自己擁有、與群組伺服器同名的伺服器仍屬於該帳號，ccam 會告知一次：`account work has its own MCP servers fs, so group default's of those names are not copied there`。
5. 它記錄新的時間戳記以及每份副本的雜湊值，然後 Claude Code 會在不帶 `--mcp-config` 的情況下啟動。

當帳號有工作階段正在執行時（不論 Claude Code 是如何啟動的），會略過同步並顯示 `ccam: account work is already running, skipping MCP sync`，因為 Claude Code 在執行期間會改寫 `.claude.json`。同步失敗會印出 `ccam: MCP sync failed (launching anyway): …`，並照常啟動。

### 退出與加入

停止共用 MCP 的帳號（`ccam group leave`、`ccam group set --remove mcp`、`ccam group remove --leave`）會保留群組的伺服器：計畫會先建立備份，再把它們複製到該帳號的 `.claude.json`，而它自己的同名伺服器會保留。加上 `--fresh` 時則不會取得任何伺服器，計畫會列出將不再載入的伺服器。當它再次加入擁有這些伺服器的群組時，沒有人改動過的保留副本會被移除，讓群組之後的編輯能送達它。

`ccam mcp import` 會把帳號自己的伺服器移到它的群組：

<!-- motion:mcp-import -->

ccam 讀取帳號的伺服器，包括它自己的伺服器和只為某個資料夾設定的伺服器，並用 `mcp.SameDef` 把每一個和群組中同名的伺服器比對，比對時會忽略鍵的順序和空白。資料夾伺服器會成為適用於每個資料夾的群組伺服器（`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`）；同一個名稱有兩種不同的資料夾定義時會被拒絕。外掛的伺服器不會列入，因為外掛會自己帶來它們。群組沒有的伺服器，或群組中定義相同的伺服器，會被複製。定義不同的會被略過，因為群組的定義是其他所有帳號在使用的：

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

加上 `--force` 時，群組的伺服器會被取代，原本的 `mcp.json` 會先儲存到群組資料夾內的 `.ccam-backup/<time>/mcp.json`。接著 `mcp.Remove` 會在同一個鎖之下、先在狀態資料夾中備份後，以一次寫入從帳號的 `.claude.json` 移除已複製的項目，包括它自己的和該資料夾的項目。群組的 `mcp.json` 從載入起就保持鎖定，直到這些項目被移除為止，所以這段期間沒有人能從群組中移除該伺服器。這一步是必要的，因為名稱衝突時帳號的項目會優先於群組的項目。`--keep-local` 會保留它們。

## 安全性和特殊情況

- **登入。** 使用 OAuth 登入的伺服器會把權杖存放在每個帳號的 `.credentials.json`（`mcpOAuth`）中，而不是伺服器清單中，所以共用伺服器並不會共用它的登入。請在每個帳號中用 `/mcp` 登入一次。ccam 絕不複製登入，[`ccam migrate`](/features/migrate) 也不會：複製的更新權杖會讓其他所有持有者登出。[`ccam ext`](/features/ext) 和 `ccam doctor` 會列出某個帳號尚未登入、而另一個帳號已登入的伺服器，`ccam next --continue` 則會說明新帳號缺少哪些登入。群組絕不會取得登入檔案。
- **正在執行的帳號。** `import` 會複製到群組，但會保留正在執行的帳號自己的項目並告訴你；請在工作階段結束後再執行一次。同步會略過正在執行的帳號。兩者也都會把不透過 ccam 啟動的 Claude Code 算在內。
- **為 Remote Control 建立的副本。** 正在執行的帳號的副本不會被更動，ccam 會顯示 `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` 使用 `mcp_strategy = "sync"` 時不需要更新，因為每次啟動都會同步。
- **損壞的檔案。** 如果 `.claude.json` 不是有效的 JSON，同步不會碰它，並回報錯誤。如果 2 秒內無法取得鎖，ccam 會回報 `could not lock … (Claude Code is writing the file)`。
- **不共用 MCP 的帳號。** `import` 會拒絕，並告訴你先把該帳號加入共用 MCP 的群組（`ccam group join <group> side`）。
- **群組名稱。** `-g` 只接受群組名稱；它會成為 `groups/` 下的資料夾，所以 `../x` 這類名稱會被拒絕。

## 相關內容

- [帳號之間的共用](/features/share)
- [透過 ccam 使用 Remote Control](/guides/remote-control)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [複製現有的設定](/features/migrate)
- [所有指令](/docs/commands)
