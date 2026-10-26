---
title: 共享的 MCP 服务器
description: 用 ccam mcp add、add-json、list、remove 和 import 为每个组维护一份 MCP 服务器列表，通过 --mcp-config 传入，或同步到每个账号中。
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

在各个组下方，`ccam mcp list` 会列出所有到不了其他账号的服务器：

- 账号为自己添加的服务器显示为 `only in account work: pg`；
- 为某一个文件夹添加的服务器（普通的 `claude mcp add`，见下文）显示为 `only in account side, for /home/me/code/app only: notes`；
- 账号自己的服务器与组服务器同名时显示 `account work has its own …, used in its sessions in place of the group's`；
- 当前文件夹的 `.mcp.json`，连同尚未批准它的账号一起显示（`waiting for approval in`）；
- `in no account, only in ~/.claude.json (plain claude's)` 的服务器，并附上复制它们的 [`ccam migrate --mcp`](/features/migrate) 命令。

添加组中已有的名称会替换它的定义，并打印 `updated` 而不是 `added`。只有属于该组、且该组共享 MCP 时（`ccam group set <group> --add mcp`），账号才会获得这些服务器；如果该组不共享 MCP，`add` 会提示你。

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
| `--dry-run`（全局） | 只有 `mcp list` 支持它；`add`、`add-json`、`remove` 和 `import` 会拒绝它，不做任何更改 | 关闭 |

`add` 需要 `--` 之后的命令或 `--url` 之一；否则会失败并显示 `a command to run (after --) or --url is required`。`--header` 用于 `--url` 服务器，`--env` 用于命令服务器。`add-json` 接受与 `claude mcp add-json` 相同的 JSON，定义必须是带有 `command` 或 `url` 的对象。`import` 总是使用账号共享 MCP 所用的组，`list` 显示所有组，因此 `-g` 对它们不起作用。

## 工作原理

涉及的文件：

- `~/.local/share/ccam/groups/<group>/mcp.json`：共享的服务器，格式为 `{"mcpServers": {…}}`，也就是 `--mcp-config` 读取的格式。它通过 `fsutil.AtomicWrite` 以 `0600` 权限写入，因为服务器条目中经常带有 API 密钥，并且只在锁 `<group>/.mcp.json.lock` 下写入，因此同时运行两个 `ccam mcp add` 命令时，两者的服务器都会保留。
- `~/.local/share/ccam/profiles/<name>/.claude.json`：Claude Code 自己的文件，保存账号的服务器。ccam 只用 `profile.TopLevelValue` 扫描器读取其中的 `mcpServers` 值，在启动路径上从不解析整个文件。
- `~/.local/state/ccam/mcp/<name>/`：`seen` 文件、`sync-stamp`（带有 ccam 写入的每个副本的哈希值，以便区分没人改过的副本和你编辑过的副本）、`sync-checked`（sync 策略上次看到的内容）、`kept`（账号退出 MCP 共享时保留的服务器）以及 `.claude.json` 的备份。

有两种策略，由 `~/.config/ccam/config.toml` 中的 `mcp_strategy` 设置。默认是 flag 策略：

<!-- motion:mcp-flag -->

1. ccam 找到账号的 MCP 组（`launcher.applyMCP`），并用一次 `stat` 确认组的 `mcp.json` 存在，而不解析它（`mcp.FlagArgs`）。
2. ccam 把 `--mcp-config=<path>` 放在 `claude` 之后、你自己的参数之前，因此你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍会在此基础上生效（`Plan.InsertArgs`）。
3. ccam 检查账号的 `.claude.json` 自上次启动以来是否有变化，只有变化时才解析其中的服务器（`mcp.NewLocal`）。
4. ccam 用 Claude Code 替换自身（`syscall.Exec`）。

路径作为单个参数 `--mcp-config=<path>` 传入，而不是两个参数。Claude Code 的 `--mcp-config` 接受多个值，因此 `claude --mcp-config a.json "prompt"` 会把提示当作第二个配置文件读取，并以 `MCP config file not found` 失败。单参数形式避免了这个问题，因此后面跟着的提示或子命令都能正常工作。`ccam work mcp list`、`get` 和 `remove` 不加这个参数，因为 Claude Code 在那里会拒绝它；ccam 会改为在标准错误中列出组的服务器（[详情](/features/launch#shared-mcp-servers)）。Claude Code 会在账号自己的服务器之外再加载 `--mcp-config` 中的服务器；名称冲突时，`.claude.json` 中的条目优先。flag 策略只在以下情况写入 `.claude.json`：Remote Control、账号退出或重新加入 MCP 共享时（[见下文](#leaving-and-joining)），以及执行 `ccam mcp import` 时。

`claude remote-control`（`rc`）不接受 `--mcp-config`，它的会话只从 `.claude.json` 读取服务器。对于 `ccam work rc`，ccam 不加这个参数，并在终端上询问是否像 `sync` 策略那样，先做备份再把组的服务器复制到 `.claude.json`。在副本存在之前它会一直询问，之后只有在组的 `mcp.json` 变化或某个已复制的服务器消失后才会再次询问（[Remote Control](/features/launch#remote-control)）。由于名称冲突时这种副本优先于 `--mcp-config`，`ccam mcp add`、`add-json`、`remove` 和 `import` 会更新组中其他账号持有的副本：

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

只有副本发生变化时才会出现 `✓ Updated` 这一行。旧版本 ccam 创建的副本也会以同样的方式更新，`ccam migrate` 之后的副本也是如此。

在会话中用 `claude mcp add` 添加的服务器只会进入该账号的 `.claude.json`。普通的 `claude mcp add` 使用本地作用域：服务器只在那一个文件夹中为那一个账号加载（`-s user` 则让该账号在所有文件夹中都能使用它）。添加发生时 ccam 无法看到，因为 Claude Code 已经替换了 ccam 进程，所以它会在下次启动时检查。`seen` 文件保存 `.claude.json` 的修改时间和已经提示过的服务器名称。当时间不同时，ccam 把这些名称与组进行比较，并为每个新服务器打印一次提示，不等待输入，两种策略都是如此：

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

提示还会列出账号自己拥有的、与组服务器同名的服务器，因为在该账号的会话中，账号自己的服务器优先。

使用 `mcp_strategy = "sync"` 时，这些服务器会作为普通的用户服务器出现在 `claude mcp list` 和 `/mcp` 中，ccam 会在启动前把它们写入每个账号的 `.claude.json`：

<!-- motion:mcp-sync -->

1. ccam 把组的 `mcp.json` 的修改时间与 `sync-stamp` 进行比较，并把账号的 `.claude.json` 与它上次看到的内容进行比较（`mcp.Sync`）。如果两者都没有变化，就不做其他任何事。否则它还会检查组的每个服务器是否仍在 `.claude.json` 中，并把消失的副本放回去（`put back group default's MCP servers that were gone from account work's .claude.json: ctx`）。
2. 它获取 Claude Code 自己的锁，即文件旁边用 `mkdir` 创建的 `.claude.json.lock` 文件夹。它最多等待 2 秒，并把超过 10 秒的锁视为失效，与 Claude Code 的锁库做法相同。
3. 它把当前的 `.claude.json` 保存到状态文件夹中的 `backups/`，并保留最新的 5 份。
4. 它在 `mcpServers` 中设置组中的每个服务器，移除它以前写入但组中已不存在的服务器，并且只在有变化时以文件原有的权限原子写入。它只更改或删除它自己写入、此后没有被任何人更改的副本。如果账号自己有一个与组服务器同名的服务器，ccam 会保留账号的服务器，并说明一次：`account work has its own MCP servers fs, so group default's servers with those names are not copied there`。
5. 它记录新的时间戳以及每个副本的哈希值，然后 Claude Code 在不带 `--mcp-config` 的情况下启动。

当该账号有会话正在运行时，无论 Claude Code 是怎样启动的，都会跳过同步并显示 `ccam: account work is already running, skipping MCP sync`，因为 Claude Code 运行时会改写 `.claude.json`。同步失败会打印 `ccam: MCP sync failed (launching anyway): …` 并照常启动。

### 退出与加入

停止共享 MCP 的账号（`ccam group leave`、`ccam group set --remove mcp`、`ccam group remove --leave`）会保留组的服务器：计划会先做备份，再把它们复制到该账号的 `.claude.json` 中，账号自己的同名服务器保持不变。使用 `--fresh` 时它不会得到任何服务器，计划会列出将不再加载的服务器。当它再次加入拥有这些服务器的组时，没人改过的保留副本会被移除，这样组之后的修改就能到达它。

`ccam mcp import` 把账号自己的服务器移到它的组中：

<!-- motion:mcp-import -->

ccam 读取账号的服务器，包括它自己的服务器以及它只为某一个文件夹拥有的服务器，并用 `mcp.SameDef` 把每个服务器与组中的同名服务器进行比较，比较时忽略键的顺序和空白。文件夹服务器会成为对所有文件夹有效的组服务器（`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`）；同一个名称有两个不同的文件夹定义时会被拒绝。插件的服务器不会被列入，因为它们由插件提供。组中没有的服务器，或者组中定义相同的服务器，会被复制。定义不同的服务器会被跳过，因为组中的定义是其他所有账号正在使用的：

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

使用 `--force` 时，组中的服务器会被替换，之前的 `mcp.json` 会先保存到组文件夹内的 `.ccam-backup/<time>/mcp.json`。然后 ccam 在同一把锁下、并在状态文件夹中备份之后，从账号的 `.claude.json` 中一次性移除已复制的条目，包括它自己的和文件夹的（`mcp.Remove`）。从读取组的 `mcp.json` 起，ccam 一直保持它的锁定，直到这些条目被移除，这样就不会有人在此期间把该服务器从组中删除。之所以需要这样做，是因为名称冲突时账号的条目会优先于组中的条目。使用 `--keep-local` 时，这些条目会保留在账号的 `.claude.json` 中。

## 安全性与边界情况

- **登录。** 通过 OAuth 登录的服务器把令牌保存在每个账号的 `.credentials.json`（`mcpOAuth`）中，而不是服务器列表中，所以共享服务器并不会共享它的登录。请在每个账号中用 `/mcp` 登录一次。ccam 从不复制登录，[`ccam migrate`](/features/migrate) 也不例外：复制过去的刷新令牌会让其他所有持有者退出登录。当某个账号未登录、而另一个账号已登录某些服务器时，[`ccam ext`](/features/ext) 和 `ccam doctor` 会列出这些服务器，`ccam next --continue` 也会说明新账号缺少哪些登录。组永远不会得到登录文件。
- **正在运行的账号。** `import` 会复制到组中，但保留正在运行的账号自己的条目，并提示你；会话结束后请再运行一次。同步会跳过正在运行的账号。两者也都能发现不通过 ccam 启动的 Claude Code 会话。
- **为 Remote Control 制作的副本。** 正在运行的账号的副本不会被改动，ccam 会提示 `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` 使用 `mcp_strategy = "sync"` 时无需更新，因为每次启动都会同步。
- **损坏的文件。** 如果 `.claude.json` 不是有效的 JSON，同步不会改动它，并报告错误。如果 2 秒内无法获取锁，ccam 会报告 `could not lock … (Claude Code is writing the file)`。
- **不共享 MCP 的账号。** `import` 会拒绝，并提示你先把该账号放进一个共享 MCP 的组（`ccam group join <group> side`）。
- **组名。** `-g` 只接受组名；它会成为 `groups/` 下的一个文件夹，因此 `../x` 这样的名称会被拒绝。

## 相关内容

- [通过组共享](/features/share)
- [用 ccam 使用 Remote Control](/guides/remote-control)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [复制现有配置](/features/migrate)
- [命令参考](/docs/commands)
