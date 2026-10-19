---
title: 指令參考
description: 所有 ccam 指令和選項，從 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

在任何指令前加上 `--dry-run`，就能查看它會做什麼而不實際執行，例如 `ccam --dry-run work`。

<!-- commands -->

## 啟動 Claude Code

`ccam <name|number> [arguments]` 會以該帳號啟動 Claude Code。ccam 不會讀取名稱後面的參數，而是依相同順序傳給 `claude`；標準輸入、結束代碼和訊號的行為都和直接執行 `claude` 完全相同。名稱後面的 `--` 可加可不加。

只有參數而沒有帳號名稱的 `ccam`，例如 `ccam -p "summarise"`，會使用預設帳號。用 `ccam use <name>` 設定預設帳號。

快捷指令 `claude1`、`claude2` 等是 `~/.local/bin` 中指向 ccam 的連結。`ccam shims install --names` 還會為每個帳號建立 `claude-<name>`。

## 新增和移除帳號

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

不指定資料夾的 `ccam import` 會尋找 `~/.claude` 和 `~/.claude-*`，並逐一詢問。匯入的資料夾會在原位置登錄；ccam 絕不會移動或複製登入資訊。

## 共用

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

每個項目可以設為 `on`（`default` 群組）、`off` 或群組名稱。請參閱[共用](/docs/sharing)。

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
ccam shims install
```

## 檔案

| 路徑 | 內容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 帳號及其共用的項目 |
| `~/.local/share/ccam/profiles/<name>/` | 每個帳號的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共用的設定、MCP 伺服器、擴充功能和對話記錄 |
| `~/.local/state/ccam/` | 用量快取和上次使用時間 |
| `~/.local/bin/claude1`, `claude2`… | 指向 ccam 的連結 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 變數可以變更這些資料夾的位置。
