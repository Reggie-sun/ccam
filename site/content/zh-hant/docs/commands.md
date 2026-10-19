---
title: 指令參考
description: 所有 ccam 指令和選項，從 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a5a33b1a2f4
---

在任何指令前加上 `--dry-run`，就能查看它會做什麼而不實際執行，例如 `ccam --dry-run work`。

<!-- commands -->

## 啟動 Claude Code

`ccam <name|number> [arguments]` 會以該帳號啟動 Claude Code。ccam 不會讀取名稱後面的參數，而是依相同順序傳給 `claude`；標準輸入、結束代碼和訊號的行為都和直接執行 `claude` 完全相同。名稱後面的 `--` 可加可不加。

只有參數而沒有帳號名稱的 `ccam`，例如 `ccam -p "summarise"`，會使用預設帳號。用 `ccam use <name>` 設定預設帳號。

`ccwork` 這類快捷指令也會以同樣的方式啟動對應的帳號：`ccwork --resume` 的效果和 `ccam work --resume` 相同。請參閱[快捷指令](/docs/commands#shortcuts)。

## 新增和移除帳號

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

不指定資料夾的 `ccam import` 會尋找 `~/.claude` 和 `~/.claude-*`，並逐一詢問。匯入的資料夾會在原位置登錄；ccam 絕不會移動或複製登入資訊。

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
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

每個項目可以設為 `on`（`default` 群組）、`off` 或群組名稱。請參閱[共用](/docs/sharing)。

## 複製現有的設定

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` 會把 Claude Code 資料夾中的 `settings.json`、`CLAUDE.md`、擴充功能和使用者 MCP 伺服器複製到一個群組，讓共用這些項目的每個帳號都能使用。該資料夾只會被讀取。選項 `--settings`、`--memory`、`--extensions`、`--mcp` 和 `--history` 只會複製你指定的項目；不加這些選項時，會複製對話記錄以外的所有內容。`-g` 或 `--group` 用來選擇群組；`-y` 或 `--yes` 不詢問直接複製，檔案內容不同時保留群組中的副本。請參閱[複製現有的設定](/docs/sharing#copy-an-existing-setup)。

## MCP 伺服器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## 用量和切換

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## 維護

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
```

`ccam update` 會在取代 ccam 之前檢查發行版本的 SHA-256 總和檢查碼，`--version v0.2.0` 則可安裝特定的發行版本。從 `.deb` 套件安裝的 ccam 則改用 apt 更新。請參閱[更新](/docs/install#update)。

## 檔案

| 路徑 | 內容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 帳號及其共用的項目 |
| `~/.local/share/ccam/profiles/<name>/` | 每個帳號的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共用的設定、MCP 伺服器、擴充功能和對話記錄 |
| `~/.local/state/ccam/` | 用量快取和上次使用時間 |
| `~/.local/bin/<shortcut>` | `ccwork` 這類快捷指令，也就是指向 ccam 的連結 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 變數可以變更這些資料夾的位置。
