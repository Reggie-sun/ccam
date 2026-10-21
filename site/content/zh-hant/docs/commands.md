---
title: 指令參考
metaTitle: "ccam 指令參考"
description: 所有 ccam 指令和選項，從 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

在指令前加上 `--dry-run`，就能查看它會做什麼而不實際執行，例如 `ccam --dry-run work`。它適用於啟動 Claude Code（`ccam <name>`、`ccam run`、`ccam next`、`ccam attach`）以及 `group`、`migrate`、`shortcut`、`update` 和 `doctor --fix`。其他會做出變更的指令則會拒絕它，而不會執行。

<!-- commands -->

## 啟動 Claude Code

`ccam <name|number> [arguments]` 會以該帳號啟動 Claude Code。ccam 不會讀取名稱後面的參數，而是依相同順序傳給 `claude`；標準輸入、結束代碼和訊號的行為都和直接執行 `claude` 完全相同。名稱後面的 `--` 可加可不加。

只有參數而沒有帳號名稱的 `ccam`，例如 `ccam -p "summarise"`，會使用為目前資料夾設定的帳號，否則使用預設帳號。分別用 `ccam folder <name> [folder]` 和 `ccam use <name>` 設定。

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

帳號名稱前的 `--fallback` 用於 `-p`/`--print` 執行。帳號達到用量上限時，ccam 會以剩餘額度最多的已登入帳號再次執行同一個指令，並且只輸出完成的那次執行的結果。請參閱[依資料夾使用帳號](/features/folder)和[達到上限時重新執行](/features/next#run-again-at-the-limit)。

`ccwork` 這類快捷指令也會以同樣的方式啟動對應的帳號：`ccwork --resume` 的效果和 `ccam work --resume` 相同。請參閱[快捷指令](/docs/commands#shortcuts)。

## 新增和移除帳號

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

不指定資料夾的 `ccam import` 會尋找 `~/.claude` 和 `~/.claude-*`，並逐一詢問。匯入的資料夾會在原位置登錄；ccam 絕不會移動或複製登入資訊。

`ccam add --api` 會新增一個使用 API 金鑰而不是 Claude 登入的帳號；ccam 會在每次啟動時讀取金鑰，絕不儲存它。請參閱 [API 金鑰帳號](/features/api)。

當 Claude Code 正以某個帳號執行時，`ccam rename`、`ccam remove` 和 `ccam group` 系列指令會拒絕變更該帳號。請先結束該工作階段。

`ccam list --json` 會輸出每個帳號及其槽位、快捷指令、電子郵件和資料夾。其中的 `groups` 陣列會列出該帳號所屬的群組，`share` 物件則會把每個項目（`mcp`、`settings`、`memory`、`extensions`、`history`）對應到它所來自的群組，或 `local`。

## 快捷指令

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

快捷指令是你自己命名的指令，用來以某個帳號啟動 Claude Code。它和 `ccam work` 一樣，會把所有參數原封不動地傳給 Claude Code，所以 `ccwork --resume` 和 `ccwork -p "hi"` 都會如預期運作。它是 `~/.local/bin` 中指向 ccam 的連結。

快捷指令是選用的。`ccam add` 設定精靈會詢問是否建立；`ccam add <name> --shortcut <command>` 則不詢問直接建立。`ccam list` 會顯示每個帳號的快捷指令，沒有時顯示 `-`。

名稱可以使用 a-z、0-9、`.`、`-` 和 `_`。名稱不能是 ccam 的指令、其他帳號的快捷指令，也不能和 `PATH` 上的其他程式同名。

`ccam rename` 會保留快捷指令的名稱，`ccam remove` 會連同帳號一起刪除快捷指令。`ccam shims install` 會重新建立所有帳號的連結，例如在你移動 ccam 之後。

## 共用

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

群組會在其中的帳號之間共用項目；一個帳號可以屬於多個群組，但每個項目只會來自一個群組。請參閱[用群組共用](/features/share)。

## 複製現有的設定

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` 會把 Claude Code 資料夾中的 `settings.json`、`CLAUDE.md`、擴充功能和使用者 MCP 伺服器複製到一個群組，讓共用這些項目的每個帳號都能使用。該資料夾只會被讀取。選項 `--settings`、`--memory`、`--extensions`、`--mcp` 和 `--history` 只會複製你指定的項目。不加這些選項時，在終端機中會顯示一份核取清單讓你選擇，除了對話記錄以外的項目都已勾選；在其他情況下，會複製對話記錄以外的所有內容。`-g` 或 `--group` 用來選擇群組；`-y` 或 `--yes` 不詢問直接複製，檔案內容不同時保留群組中的副本。它也會把這些伺服器已儲存的登入複製到每個帳號自己的登入檔案。請參閱[複製現有的設定](/docs/sharing#copy-an-existing-setup)。

## MCP 伺服器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` 會把只有該帳號擁有的伺服器移到它的群組。如果群組已有同名但定義不同的伺服器，該伺服器會被略過；加上 `--force` 可取代群組中的副本，ccam 會先備份它。

## 用量和切換

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`、`--fallback` 和 `ccam pick` 的選擇方式相同：優先選低於 `auto_max_usage`（`config.toml`，預設 100）的帳號，依 `auto_strategy`（`headroom`、`expiring` 或 `balance`）選擇，會計入 `--model` 的每週區間，絕不選擇用量 100% 的帳號，而且只在起始帳號所屬的組織內選擇（加上 `--cross-org` 或設定 `cross_org = true` 可擴大範圍）。無法選出帳號時，它們會以代碼 75 結束，加上 `--wait` 時則等待第一個帳號重設。即使帳號之間沒有共用對話記錄，`ccam next --continue` 也能把對話帶過去：此時 Claude Code 會把對話複製到新的帳號。和 `ccam <name>` 不同，`ccam next` 會讀取自己的選項，所以 Claude Code 的選項要放在 `--` 後面；少了它時 ccam 會以 `unknown flag` 停止。在 Claude Code 中，外掛的 `/ccam:switch [account]` 會把對話交接出去，而 `ccam hook` 就是它的 hook 所執行的指令。請參閱[切換到下一個帳號](/features/next)、[為指令碼選擇帳號](/features/pick)和[所有帳號的背景工作階段](/features/agents)。

## 維護

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` 會在取代 ccam 之前檢查發行版本的 SHA-256 總和檢查碼，`--version v0.2.0` 則可安裝特定的發行版本。從 `.deb` 套件安裝的 ccam 則改用 apt 更新。請參閱[更新](/docs/install#update)。

## 檔案

| 路徑 | 內容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 帳號、群組和設定（`theme`、`usage_source`、`relay`、`cross_org` …） |
| `~/.local/share/ccam/profiles/<name>/` | 每個帳號的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共用的設定、MCP 伺服器、擴充功能和對話記錄 |
| `~/.local/state/ccam/` | 用量快取、上次使用時間，以及 `/ccam:switch` 和上限 hook 留下的記錄（`handoff/`） |
| `~/.local/bin/<shortcut>` | `ccwork` 這類快捷指令，也就是指向 ccam 的連結 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 變數可以變更這些資料夾的位置。在 Windows 上，`~` 是 `%USERPROFILE%`，快捷指令則是 `<shortcut>.exe`。
