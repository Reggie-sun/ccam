---
title: 共用的 MCP 伺服器
description: 用 ccam mcp add、add-json、list、remove 和 import 為每個群組維護一份 MCP 伺服器清單，透過 --mcp-config 傳入或同步到每個帳號。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

新增群組中已有的名稱時，會取代它的定義，並印出 `updated` 而不是 `added`。只有和該群組共用 MCP 的帳號（`ccam share <name> --mcp on`）才會取得這些伺服器；如果還沒有任何帳號共用 MCP，`add` 會告訴你。

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

- `~/.local/share/ccam/groups/<group>/mcp.json`：共用的伺服器，格式為 `{"mcpServers": {…}}`，也就是 `--mcp-config` 讀取的格式。它以 `fsutil.AtomicWrite` 和權限 `0600` 寫入，因為伺服器項目常常帶有 API 金鑰。
- `~/.local/share/ccam/profiles/<name>/.claude.json`：Claude Code 自己的檔案，存放帳號的伺服器。ccam 只用 `profile.TopLevelValue` 掃描器讀取其中的 `mcpServers` 值，在啟動流程中絕不解析整個檔案。
- `~/.local/state/ccam/mcp/<name>/`：`seen` 檔案、`sync-stamp` 和 `.claude.json` 的備份。

有兩種策略，由 `~/.config/ccam/config.toml` 中的 `mcp_strategy` 設定。預設是 flag 策略：

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` 找出帳號的 MCP 群組。`mcp.FlagArgs` 以一次 `stat` 確認群組的 `mcp.json` 存在，而不解析它。
2. `Plan.InsertArgs` 把 `--mcp-config=<path>` 放在 `claude` 的正後方、你自己的參數之前，所以你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍會疊加套用。
3. `mcp.NewLocal` 檢查帳號的 `.claude.json` 自上次啟動後是否變更，只有變更時才解析其中的伺服器。
4. `syscall.Exec` 以 Claude Code 取代 ccam。

路徑以單一參數 `--mcp-config=<path>` 傳入，而不是兩個參數。Claude Code 的 `--mcp-config` 接受多個值，所以 `claude --mcp-config a.json "prompt"` 會把提示詞當成第二個設定檔讀取，並以 `MCP config file not found` 失敗。單一參數的寫法會在那裡結束，後面接的提示詞或 `mcp list` 這類子指令都能正常運作。Claude Code 會在帳號自己的伺服器之外再載入 `--mcp-config` 的伺服器；名稱衝突時，`.claude.json` 中的項目優先。flag 策略絕不寫入 `.claude.json`。

在工作階段中用 `claude mcp add` 新增的伺服器只會存到該帳號的 `.claude.json`。當下 ccam 看不到它，因為 Claude Code 已經取代了 ccam，所以它會在下次啟動時檢查。`seen` 檔案記錄 `.claude.json` 的修改時間和已經回報過的伺服器名稱。時間不同時，ccam 會把名稱和群組比對，並為每個新伺服器印出一次提示，不會等待輸入：

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

使用 `mcp_strategy = "sync"` 時，這些伺服器會以一般使用者伺服器的形式出現在 `claude mcp list` 和 `/mcp` 中，ccam 會在啟動前把它們寫入每個帳號的 `.claude.json`：

<!-- motion:mcp-sync -->

1. `mcp.Sync` 把群組 `mcp.json` 的修改時間和 `sync-stamp` 比對。如果沒有變更，就不會開啟帳號的 `.claude.json`。
2. 它取得 Claude Code 自己的鎖，也就是檔案旁邊以 `mkdir` 建立的資料夾 `.claude.json.lock`。它最多等待 2 秒，並和 Claude Code 的鎖函式庫一樣，把超過 10 秒的鎖視為已失效。
3. 它把目前的 `.claude.json` 儲存到狀態資料夾中的 `backups/`，並保留最新的 5 份。
4. 它在 `mcpServers` 中設定每個群組伺服器，移除它之前寫入、但群組已經沒有的伺服器，保留帳號自己的伺服器，並以檔案原本的權限原子寫入。
5. 它記錄新的時間戳記和它所管理的名稱，然後 Claude Code 會在不帶 `--mcp-config` 的情況下啟動。

當帳號有工作階段正在執行時，會略過同步並顯示 `ccam: account work is already running, skipping MCP sync`，因為 Claude Code 在執行期間會改寫 `.claude.json`。同步失敗會印出 `ccam: MCP sync failed (launching anyway): …`，並照常啟動。關於工作階段中新增伺服器的提示只屬於 flag 策略。

`ccam mcp import` 會把帳號自己的伺服器移到它的群組：

<!-- motion:mcp-import -->

ccam 用 `mcp.LocalServers` 讀取帳號的伺服器，並用 `mcp.SameDef` 把每一個和群組中同名的伺服器比對，比對時會忽略鍵的順序和空白。群組沒有的伺服器，或群組中定義相同的伺服器，會被複製。定義不同的會被略過，因為群組的定義是其他所有帳號在使用的：

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

加上 `--force` 時，群組的伺服器會被取代，原本的 `mcp.json` 會先儲存到群組資料夾內的 `.ccam-backup/<time>/mcp.json`。接著 `mcp.RemoveLocal` 會在同一個鎖之下、先在狀態資料夾中備份後，從帳號的 `.claude.json` 移除已複製的項目。這一步是必要的，因為名稱衝突時帳號的項目會優先於群組的項目。`--keep-local` 會保留它們。

## 安全性和特殊情況

- **登入。** 使用 OAuth 登入的伺服器會把權杖存放在每個帳號的 `.credentials.json`（`mcpOAuth`）中，而不是伺服器清單中，所以共用伺服器並不會共用它的登入。請在每個帳號中用 `/mcp` 登入一次，或用 [`ccam migrate`](/features/migrate) 從 `~/.claude` 複製這些登入。群組絕不會取得登入檔案。
- **正在執行的帳號。** `import` 會複製到群組，但會保留正在執行的帳號自己的項目並告訴你；請在工作階段結束後再執行一次。同步會略過正在執行的帳號。
- **損壞的檔案。** 如果 `.claude.json` 不是有效的 JSON，同步不會碰它，並回報錯誤。如果 2 秒內無法取得鎖，ccam 會回報 `could not lock … (Claude Code is writing the file)`。
- **不共用 MCP 的帳號。** `import` 會拒絕，並告訴你先執行 `ccam share side --mcp on`。
- **群組名稱。** `-g` 只接受群組名稱；它會成為 `groups/` 下的資料夾，所以 `../x` 這類名稱會被拒絕。

## 相關內容

- [帳號之間的共用](/features/share)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [複製現有的設定](/features/migrate)
- [所有指令](/docs/commands)
