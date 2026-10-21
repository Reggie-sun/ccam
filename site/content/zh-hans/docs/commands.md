---
title: 命令参考
metaTitle: "ccam 命令参考"
description: ccam 的所有命令和选项，从 ccam add、ccam <name> 到 ccam usage、ccam next 和 ccam doctor。
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

在命令前加上 `--dry-run`，可以查看它会做什么而不实际执行，例如 `ccam --dry-run work`。它适用于启动 Claude Code（`ccam <name>`、`ccam run`、`ccam next`、`ccam attach`）以及 `group`、`migrate`、`shortcut`、`update` 和 `doctor --fix`。其他会做出更改的命令会拒绝它，而不会执行。

<!-- commands -->

## 启动 Claude Code

`ccam <name|number> [arguments]` 以该账号启动 Claude Code。ccam 不会读取账号名之后的参数，而是按相同顺序把它们传给 `claude`；标准输入、退出码和信号的行为与直接运行 `claude` 完全相同。账号名之后的 `--` 是可选的。

`ccam` 带参数但不带账号名时（例如 `ccam -p "summarise"`），会使用为当前文件夹设置的账号，否则使用默认账号。分别用 `ccam folder <name> [folder]` 和 `ccam use <name>` 设置。

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

账号名之前的 `--fallback` 用于 `-p`/`--print` 运行。当账号达到用量限制时，ccam 会以余量最多的已登录账号重新运行同一个命令，并只打印最终完成的那次运行的输出。参见[按文件夹使用账号](/features/folder)和[达到限制时重新运行](/features/next#run-again-at-the-limit)。

`ccwork` 这样的快捷命令也以同样的方式启动对应账号：`ccwork --resume` 的效果与 `ccam work --resume` 相同。参见[快捷命令](/docs/commands#shortcuts)。

## 添加和删除账号

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

不指定文件夹的 `ccam import` 会查找 `~/.claude` 和 `~/.claude-*`，并逐个询问。导入的文件夹在原位置注册；ccam 从不移动或复制登录信息。

`ccam add --api` 会添加一个使用 API 密钥而不是 Claude 登录的账号；ccam 在每次启动时读取密钥，从不保存它。参见 [API 密钥账号](/features/api)。

当 Claude Code 正以某个账号运行时，`ccam rename`、`ccam remove` 和各个 `ccam group` 命令会拒绝更改该账号。请先退出该会话。

`ccam list --json` 会输出每个账号及其槽位、快捷命令、邮箱和文件夹。其中的 `groups` 数组列出账号所在的组，`share` 对象则把每一项（`mcp`、`settings`、`memory`、`extensions`、`history`）映射到提供它的组，或映射到 `local`。

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

组在其账号之间共享各项；一个账号可以属于多个组，但每一项只能来自一个组。参见[用组共享](/features/share)。

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

`next --auto`、`--fallback` 和 `ccam pick` 的选择方式相同：优先选择低于 `auto_max_usage`（`config.toml`，默认 100）的账号，按 `auto_strategy`（`headroom`、`expiring` 或 `balance`）选择，并计入 `--model` 的每周时间窗口，从不选择已达 100% 的账号，并且只在出发账号所属的组织内选择（用 `--cross-org` 或 `cross_org = true` 扩大范围）。当没有可选的账号时，它们以退出码 75 退出，或者在使用 `--wait` 时等待第一个账号重置。即使两个账号不共享历史记录，`ccam next --continue` 也能把对话带过去：此时由 Claude Code 把对话复制到新账号中。与 `ccam <name>` 不同，`ccam next` 会读取自己的选项，所以 Claude Code 的选项要放在 `--` 之后；没有它时 ccam 会以 `unknown flag` 停止。在 Claude Code 中，插件的 `/ccam:switch [account]` 用于交接对话，`ccam hook` 则是插件的钩子所运行的命令。参见[切换到下一个账号](/features/next)、[为脚本选择账号](/features/pick)和[所有账号的后台会话](/features/agents)。

## 维护

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

`ccam update` 在替换 ccam 之前会校验版本的 SHA-256 校验和，`--version v0.2.0` 用于安装特定版本。通过 `.deb` 软件包安装的 ccam 则用 apt 更新。参见[更新](/docs/install#update)。

## 文件

| 路径 | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | 账号、组和设置（`theme`、`usage_source`、`relay`、`cross_org` …） |
| `~/.local/share/ccam/profiles/<name>/` | 每个账号的 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共享的设置、MCP 服务器、扩展和历史记录 |
| `~/.local/state/ccam/` | 用量缓存、最近使用时间，以及 `/ccam:switch` 和用量限制钩子留下的记录（`handoff/`） |
| `~/.local/bin/<shortcut>` | `ccwork` 这样的快捷命令，即指向 ccam 的链接 |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME` 和 `XDG_STATE_HOME` 变量可以改变这些文件夹的位置。在 Windows 上，`~` 是 `%USERPROFILE%`，快捷命令是 `<shortcut>.exe`。
