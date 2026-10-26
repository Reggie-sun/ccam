---
title: 所有账号的后台会话
description: ccam agents 在一张表中列出你所有账号的 Claude Code 后台代理和打开的会话，ccam attach 可以按 ID 打开其中一个。
updated: 2026-10-26
order: 16
translation: machine
sourceHash: "06281fb13119"
---

Claude Code 按账号保存后台会话（`claude --bg`、代理视图）：`claude agents` 只显示它当前所用账号的会话。`ccam agents` 会同时询问每个账号并把结果列在一起，`ccam attach <id|name>` 则以启动该后台会话的账号打开它。

## 用法

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` 显示后台会话的状态（`working`、`blocked`、`done`、`failed`、`stopped`）以及其进程正在做什么（`busy`、`waiting`、`idle`）；打开的交互式会话显示为 `open`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--all` | 同时列出已完成、失败或已停止的后台会话 | 关闭 |
| `--json` | 以 JSON 输出会话，每个会话都带有它所属的 `account` | 关闭 |

`ccam attach <id|name> [claude args]` 会把 ID 之后的所有内容原样传给 `claude attach`，包括各种标志：`ccam attach 49ec --dry-run` 会把 `--dry-run` 交给 Claude Code。如果想改为打印命令，请把 `--dry-run` 放在 ID 之前（`ccam --dry-run attach 49ec` 或 `ccam attach --dry-run 49ec`）。

## 工作原理

1. ccam 为每个账号运行 `claude agents --json`（指定了 `--all` 时也会带上），使用该账号的 `CLAUDE_CONFIG_DIR`，所有账号并行进行，每个超时为 15 秒。这只读取账号自己的会话列表，不会启动 Claude Code 的后台服务。
2. 每个条目都会加上账号名，列表按启动时间排序，最新的在前。没有响应的账号会在标准错误中列出名称，并被排除在列表之外，随后 ccam 以退出码 1 退出，这样脚本就不会把失败误当成空列表。
3. `ccam attach` 会在每个账号的文件夹中查找 `jobs/<id>`（Claude Code 在这里保存后台会话），然后以该账号运行 `claude attach <id>`，方式与 `ccam <name>` 启动 Claude Code 相同。不是 ID 前缀的词会与会话名称进行匹配，不区分大小写；看起来像 ID 的词永远不会被当作名称。

## 安全性与边界情况

- 后台会话的 ID 只属于一个账号：以另一个账号运行的 `claude attach` 找不到它们，而 `ccam attach` 正好解决了这个问题。
- 如果两个账号中的会话有相同的 ID 前缀或名称中有相同的词，`ccam attach` 会列出它们，并要求你输入更长的 ID 或名称。
- ccam 以单个参数 `--mcp-config=<file>` 把共享的 MCP 服务器传给 Claude Code，`claude agents` 和 `claude attach` 都接受这种写法；后台会话重启时会继续使用该文件。
- API 密钥账号的后台服务从启动它的 shell 获取密钥。请通过 ccam 启动这类账号的第一个后台会话（`ccam gw --bg …`），让服务拿到密钥。

## 相关内容

- [启动 Claude Code](/features/launch)
- [为脚本选择账号](/features/pick)
- [仪表盘](/features/dashboard)
