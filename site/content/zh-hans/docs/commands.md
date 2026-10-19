---
title: 命令参考
description: ccam 的所有命令和选项，从 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

在任何命令前加上 `--dry-run`，可以查看它会做什么而不实际执行，例如 `ccam --dry-run work`。

<!-- commands -->

## 启动 Claude Code

`ccam <name|number> [arguments]` 以该账号启动 Claude Code。ccam 不会读取账号名之后的参数，而是按相同顺序把它们传给 `claude`；标准输入、退出码和信号的行为与直接运行 `claude` 完全相同。账号名之后的 `--` 是可选的。

`ccam` 带参数但不带账号名时（例如 `ccam -p "summarise"`），会使用默认账号。用 `ccam use <name>` 设置默认账号。

快捷命令 `claude1`、`claude2` 等是 `~/.local/bin` 中指向 ccam 的链接。`ccam shims install --names` 还会为每个账号创建 `claude-<name>`。

## 添加和删除账号

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

不指定文件夹的 `ccam import` 会查找 `~/.claude` 和 `~/.claude-*`，并逐个询问。导入的文件夹在原位置注册；ccam 从不移动或复制登录信息。

## 共享

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

每个共享项可以设为 `on`（即 `default` 组）、`off` 或一个组名。参见[共享](/docs/sharing)。

## MCP 服务器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## 用量与切换

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## 维护

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## 文件

| 路径 | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 账号及其共享内容 |
| `~/.local/share/ccam/profiles/<name>/` | 每个账号的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共享的设置、MCP 服务器、扩展和历史记录 |
| `~/.local/state/ccam/` | 用量缓存和最近使用时间 |
| `~/.local/bin/claude1`、`claude2`… | 指向 ccam 的链接 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 变量可以改变这些文件夹的位置。
