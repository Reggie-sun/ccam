---
title: 命令参考
description: ccam 的所有命令和选项，从 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a82175c30d3
---

在命令前加上 `--dry-run`，可以查看它会做什么而不实际执行，例如 `ccam --dry-run work`。它适用于启动 Claude Code（`ccam <name>`、`ccam run`、`ccam next`）以及 `share`、`migrate`、`shortcut`、`update` 和 `doctor --fix`。其他会做出更改的命令会拒绝它，而不会执行。

<!-- commands -->

## 启动 Claude Code

`ccam <name|number> [arguments]` 以该账号启动 Claude Code。ccam 不会读取账号名之后的参数，而是按相同顺序把它们传给 `claude`；标准输入、退出码和信号的行为与直接运行 `claude` 完全相同。账号名之后的 `--` 是可选的。

`ccam` 带参数但不带账号名时（例如 `ccam -p "summarise"`），会使用默认账号。用 `ccam use <name>` 设置默认账号。

`ccwork` 这样的快捷命令也以同样的方式启动对应账号：`ccwork --resume` 的效果与 `ccam work --resume` 相同。参见[快捷命令](/docs/commands#shortcuts)。

## 添加和删除账号

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

不指定文件夹的 `ccam import` 会查找 `~/.claude` 和 `~/.claude-*`，并逐个询问。导入的文件夹在原位置注册；ccam 从不移动或复制登录信息。

当 Claude Code 正以某个账号运行时，`ccam rename`、`ccam remove` 和 `ccam share` 会拒绝更改该账号。请先退出该会话。

`ccam list --json` 会输出每个账号及其槽位、快捷命令、邮箱和文件夹。其中的 `share` 对象把每一项（`mcp`、`settings`、`memory`、`extensions`、`history`）映射到一个组名或 `local`。

## 快捷命令

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

快捷命令是你自己命名的命令，用来以某个账号启动 Claude Code。它和 `ccam work` 一样，会把所有参数原样传给 Claude Code，所以 `ccwork --resume` 和 `ccwork -p "hi"` 都会按预期工作。它是 `~/.local/bin` 中指向 ccam 的链接。

快捷命令是可选的。`ccam add` 向导会询问是否创建；`ccam add <name> --shortcut <command>` 则不询问直接创建。`ccam list` 会显示每个账号的快捷命令，没有时显示 `-`。

名称可以使用 a-z、0-9、`.`、`-` 和 `_`。名称不能是 ccam 的命令、其他账号的快捷命令，也不能与 `PATH` 上的其他程序同名。

`ccam rename` 会保留快捷命令的名称，`ccam remove` 会连同账号一起删除快捷命令。`ccam shims install` 会重新创建所有账号的链接，例如在你移动了 ccam 之后。

## 共享

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

每个共享项可以设为 `on`（即 `default` 组）、`off` 或一个组名。参见[共享](/docs/sharing)。

## 复制现有配置

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` 会把 Claude Code 文件夹中的 `settings.json`、`CLAUDE.md`、扩展和用户 MCP 服务器复制到一个组中，让共享这些项的每个账号都能用上。该文件夹只会被读取。选项 `--settings`、`--memory`、`--extensions`、`--mcp` 和 `--history` 只复制你指定的项。不加这些选项时，在终端中会显示一个清单供你选择，除对话历史外的所有项都已勾选；在其他情况下，会复制除对话历史以外的所有内容。`-g` 或 `--group` 用于选择组；`-y` 或 `--yes` 不询问直接复制，文件内容不同时保留组中的副本。它还会把这些服务器已保存的登录复制到每个账号自己的登录文件中。参见[复制现有配置](/docs/sharing#copy-an-existing-setup)。

## MCP 服务器

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` 会把只有该账号拥有的服务器移到它的组中。如果组中已有同名但定义不同的服务器，该服务器会被跳过；加上 `--force` 可替换组中的副本，ccam 会先备份它。

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
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

`ccam update` 在替换 ccam 之前会校验版本的 SHA-256 校验和，`--version v0.2.0` 用于安装特定版本。通过 `.deb` 软件包安装的 ccam 则用 apt 更新。参见[更新](/docs/install#update)。

## 文件

| 路径 | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 账号及其共享内容 |
| `~/.local/share/ccam/profiles/<name>/` | 每个账号的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共享的设置、MCP 服务器、扩展和历史记录 |
| `~/.local/state/ccam/` | 用量缓存和最近使用时间 |
| `~/.local/bin/<shortcut>` | `ccwork` 这样的快捷命令，即指向 ccam 的链接 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 变量可以改变这些文件夹的位置。
