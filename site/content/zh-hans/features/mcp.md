---
title: 共享的 MCP 服务器
description: 用 ccam mcp add、add-json、list、remove 和 import 为每个组维护一份 MCP 服务器列表，通过 --mcp-config 传入，或同步到每个账号中。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
---

Claude Code 把用户 MCP 服务器保存在 `.claude.json` 中，而这个文件同时保存着账号的登录，所以它不能在账号之间共享。ccam 把共享的服务器保存在组文件 `mcp.json` 中，并在启动时把它交给该组的每个账号。用 `ccam mcp` 可以一次性为所有账号添加服务器，用 `ccam mcp import` 可以把你在某个会话中添加的服务器转给其他账号。

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

更改从每个账号下次启动时生效；正在运行的会话保留它启动时的服务器：

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

添加组中已有的名称会替换它的定义，并打印 `updated` 而不是 `added`。只有与该组共享 MCP 的账号（`ccam share <name> --mcp on`）才会获得这些服务器；如果还没有任何账号共享 MCP，`add` 会提示你。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `-g`、`--group <group>` | `add`、`add-json` 和 `remove` 的目标组 | 如果只有一个组被某些账号用来共享 MCP，就是该组，否则为 `default` |
| `-- <command> [args…]`（`add`） | 把这个命令作为 stdio 服务器运行 | |
| `--url <url>`（`add`） | 远程服务器 URL（http/sse） | |
| `--transport <stdio\|http\|sse>`（`add`） | `--url` 服务器的传输方式 | `http` |
| `--header <'Name: value'>`（`add`） | HTTP 请求头，可重复使用 | 无 |
| `-e`、`--env <KEY=VALUE>`（`add`） | stdio 服务器的环境变量，可重复使用 | 无 |
| `--force`（`import`） | 替换组中同名但定义不同的服务器；旧的 `mcp.json` 会被备份 | 关闭 |
| `--keep-local`（`import`） | 保留账号 `.claude.json` 中的条目 | 关闭 |
| `-h`、`--help` | 各命令的帮助 | |
| `--dry-run`（全局） | 只有 `mcp list` 接受它；`add`、`add-json`、`remove` 和 `import` 会拒绝它，不做任何更改 | 关闭 |

`add` 需要 `--` 之后的命令或 `--url` 之一；否则会失败并显示 `a command to run (after --) or --url is required`。`--header` 用于 `--url` 服务器，`--env` 用于命令服务器。`add-json` 接受与 `claude mcp add-json` 相同的 JSON，定义必须是带有 `command` 或 `url` 的对象。`import` 总是使用账号共享 MCP 所用的组，`list` 显示所有组，因此 `-g` 对它们不起作用。

## 工作原理

涉及的文件：

- `~/.local/share/ccam/groups/<group>/mcp.json`：共享的服务器，格式为 `{"mcpServers": {…}}`，也就是 `--mcp-config` 读取的格式。它通过 `fsutil.AtomicWrite` 以 `0600` 权限写入，因为服务器条目中经常带有 API 密钥。
- `~/.local/share/ccam/profiles/<name>/.claude.json`：Claude Code 自己的文件，保存账号的服务器。ccam 只用 `profile.TopLevelValue` 扫描器读取其中的 `mcpServers` 值，在启动路径上从不解析整个文件。
- `~/.local/state/ccam/mcp/<name>/`：`seen` 文件、`sync-stamp` 和 `.claude.json` 的备份。

有两种策略，由 `~/.config/ccam/config.toml` 中的 `mcp_strategy` 设置。默认是 flag 策略：

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` 找到账号的 MCP 组。`mcp.FlagArgs` 用一次 `stat` 确认组的 `mcp.json` 存在，而不解析它。
2. `Plan.InsertArgs` 把 `--mcp-config=<path>` 放在 `claude` 之后、你自己的参数之前，因此你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍会在此基础上生效。
3. `mcp.NewLocal` 检查账号的 `.claude.json` 自上次启动以来是否有变化，只有变化时才解析其中的服务器。
4. `syscall.Exec` 用 Claude Code 替换 ccam。

路径作为单个参数 `--mcp-config=<path>` 传入，而不是两个参数。Claude Code 的 `--mcp-config` 接受多个值，因此 `claude --mcp-config a.json "prompt"` 会把提示当作第二个配置文件读取，并以 `MCP config file not found` 失败。单参数形式到此为止，后面跟着的提示或 `mcp list` 这样的子命令都能正常工作。Claude Code 会在账号自己的服务器之外再加载 `--mcp-config` 中的服务器；名称冲突时，`.claude.json` 中的条目优先。flag 策略从不写入 `.claude.json`。

在会话中用 `claude mcp add` 添加的服务器只会进入该账号的 `.claude.json`。ccam 当时无法看到它，因为 Claude Code 已经替换了 ccam，所以它会在下次启动时检查。`seen` 文件保存 `.claude.json` 的修改时间和已经提示过的服务器名称。当时间不同时，ccam 把这些名称与组进行比较，并为每个新服务器打印一次提示，不等待输入：

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

使用 `mcp_strategy = "sync"` 时，这些服务器会作为普通的用户服务器出现在 `claude mcp list` 和 `/mcp` 中，ccam 会在启动前把它们写入每个账号的 `.claude.json`：

<!-- motion:mcp-sync -->

1. `mcp.Sync` 把组的 `mcp.json` 的修改时间与 `sync-stamp` 进行比较。如果没有变化，就不打开账号的 `.claude.json`。
2. 它获取 Claude Code 自己的锁，即文件旁边用 `mkdir` 创建的 `.claude.json.lock` 文件夹。它最多等待 2 秒，并把超过 10 秒的锁视为失效，与 Claude Code 的锁库做法相同。
3. 它把当前的 `.claude.json` 保存到状态文件夹中的 `backups/`，并保留最新的 5 份。
4. 它在 `mcpServers` 中设置组中的每个服务器，移除它以前写入但组中已不存在的服务器，保留账号自己的服务器，并以文件原有的权限原子写入。
5. 它记录新的时间戳和它所管理的名称，然后 Claude Code 在不带 `--mcp-config` 的情况下启动。

当该账号有会话正在运行时，会跳过同步并显示 `ccam: account work is already running, skipping MCP sync`，因为 Claude Code 运行时会改写 `.claude.json`。同步失败会打印 `ccam: MCP sync failed (launching anyway): …` 并照常启动。对会话中新增服务器的提示只属于 flag 策略。

`ccam mcp import` 把账号自己的服务器移到它的组中：

<!-- motion:mcp-import -->

ccam 用 `mcp.LocalServers` 读取账号的服务器，并用 `mcp.SameDef` 把每个服务器与组中的同名服务器进行比较，比较时忽略键的顺序和空白。组中没有的服务器，或者组中定义相同的服务器，会被复制。定义不同的服务器会被跳过，因为组中的定义是其他所有账号正在使用的：

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

使用 `--force` 时，组中的服务器会被替换，之前的 `mcp.json` 会先保存到组文件夹内的 `.ccam-backup/<time>/mcp.json`。然后 `mcp.RemoveLocal` 在同一把锁下、并在状态文件夹中备份之后，从账号的 `.claude.json` 中移除已复制的条目。之所以需要这样做，是因为名称冲突时账号的条目会优先于组中的条目。`--keep-local` 会保留它们。

## 安全性与边界情况

- **登录。** 通过 OAuth 登录的服务器把令牌保存在每个账号的 `.credentials.json`（`mcpOAuth`）中，而不是服务器列表中，所以共享服务器并不会共享它的登录。请在每个账号中用 `/mcp` 登录一次，或者用 [`ccam migrate`](/features/migrate) 从 `~/.claude` 复制这些登录。组永远不会得到登录文件。
- **正在运行的账号。** `import` 会复制到组中，但保留正在运行的账号自己的条目，并提示你；会话结束后请再运行一次。同步会跳过正在运行的账号。
- **损坏的文件。** 如果 `.claude.json` 不是有效的 JSON，同步不会改动它，并报告错误。如果 2 秒内无法获取锁，ccam 会报告 `could not lock … (Claude Code is writing the file)`。
- **不共享 MCP 的账号。** `import` 会拒绝，并提示你先运行 `ccam share side --mcp on`。
- **组名。** `-g` 只接受组名；它会成为 `groups/` 下的一个文件夹，因此 `../x` 这样的名称会被拒绝。

## 相关内容

- [在账号之间共享](/features/share)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [复制现有配置](/features/migrate)
- [所有命令](/docs/commands)
