---
title: 启动 Claude Code
description: 用 ccam <name|number>、ccam run 或 ccam 选择器以任意账号启动 Claude Code；参数原样传给 claude，--dry-run 会显示将执行的命令。
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work` 以账号 `work` 启动 Claude Code，账号名之后的所有内容都原样传给 `claude`。凡是你会输入 `claude` 的地方都可以用它：在终端中、在脚本中，或者通过管道输入。

## 用法

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

给出账号名或它在 `ccam list` 中的编号。ccam 不会读取账号名之后的参数，所以 `ccam work --help` 显示的是 Claude Code 的帮助。账号名之后的 `--` 是可选的。标准输入、退出码和信号的行为与 `claude` 相同。

当第一个参数以 `-` 开头且不是 ccam 自己的选项时，ccam 会用全部参数启动一个不指定名称的账号。这个账号是用 [`ccam folder`](/features/folder) 为当前文件夹设置的账号，否则就是默认账号：用 `ccam use` 设置的账号；如果没有设置，则是编号最小的账号。

### 选择器

不带参数的 `ccam` 在标准输入和输出都是终端时打开选择器，否则打印 `ccam list`。它显示每个账号的编号、邮箱、5 小时用量、最近使用时间和共享项；用量会在每个账号返回结果时逐个填入。

| 按键 | 作用 |
| --- | --- |
| `Enter` | 以选中的账号启动 Claude Code |
| `1` 到 `9` | 启动对应编号的账号 |
| `Up`/`Down`、`k`/`j` | 移动选择 |
| `a` | 关闭选择器并运行 `ccam add` 向导 |
| `s` | 编辑选中账号的共享项（`Space` 或 `x` 切换，`Enter` 应用，`Esc` 返回） |
| `d` | 把选中的账号设为默认账号 |
| `r` | 不使用短期缓存，重新获取用量 |
| `q`、`Esc`、`Ctrl+C` | 退出 |

选择器会在执行操作之前关闭，因此终端会在 Claude Code 启动前恢复原状。按 `d` 或更改共享之后，选择器会再次打开。更改共享时会像 `ccam share` 一样先显示计划并征求确认。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--dry-run` | 打印 ccam 将使用的环境变量和命令，然后退出 | 关闭 |
| `--fallback` | 用于 `-p`/`--print` 运行：当账号达到用量限制时，以余量最多的已登录账号重新运行同一个命令（[详情](/features/next#run-again-at-the-limit)） | 关闭 |
| `-v`、`--version` | 打印 `ccam <version>`（仅作为第一个参数时） | |
| `-h`、`--help` | ccam 的帮助，形式为 `ccam --help` 或 `ccam run --help` | |

参数规则：

- ccam 只读取账号名之前它自己的选项。账号名之后的每个参数都属于 Claude Code，包括 `--dry-run` 和 `--help`。
- `--dry-run` 可以写成 `ccam --dry-run work`、`ccam run --dry-run work` 和 `ccam --dry-run`（不指定名称的账号）。它会被执行，而不是被拒绝。`--dry-run` 和 `--fallback` 的顺序可以任意。
- `ccam run <name|number> [claude args]` 是 `ccam <name>` 的完整形式。不带账号名的 `ccam run` 使用当前文件夹的账号或默认账号。
- 账号名不能是 `add`、`list`、`run` 或 `usage` 这样的 ccam 命令名，因此 `ccam <word>` 永远不会有歧义。

试运行会打印一行可以直接粘贴到 shell 中（在 Windows 上则是 PowerShell 中）的命令：

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## 工作原理

每次启动都经过相同的步骤，从命令行一直到 `exec`：

<!-- motion:launch -->

每个方框会在对应步骤执行时亮起，下面的几行是 Claude Code 启动时使用的环境变量和参数。

1. `cmd/ccam/main.go` 调用 `app.Main`。如果 ccam 是以其他名称启动的，那就是一个[快捷命令](/features/shortcuts)，直接进入启动器。
2. `app.parse` 在 cobra 介入之前查看 `os.Args`。只有单独的 `ccam`（选择器）、ccam 命令名（`config.Reserved`）、`-h`/`--help` 和 `-v`/`--version` 会进入 `cli.Execute`。其他一切都是启动。
3. `launcher.Launch` 读取 `~/.config/ccam/config.toml`，这是它解析的唯一文件。
4. `launcher.Prepare` 按名称或编号找到账号（不指定名称时，`HereProfile` 会取当前文件夹的账号，它把工作目录与 `config.toml` 中的 `[[folder]]` 条目按字符串进行比较），去掉开头的一个 `--`，并用 `ResolveClaude` 解析 `claude` 可执行文件。PATH 查找结果连同文件的修改时间一起缓存在 `~/.local/state/ccam/claude-path` 中，因此 ccam 不会在每次启动时搜索 `PATH`。文件变化或不存在时会重新查找。`config.toml` 中的 `claude_bin` 路径会被原样使用。如果 `claude` 就是 ccam 自身（按 inode 比较），ccam 会拒绝它，因为那会造成循环。
5. `checkLinks` 用 `share.Check` 把账号文件夹 `~/.local/share/ccam/profiles/<name>/` 与它的共享设置进行比较；它只调用 `lstat` 和 `readlink`，每项一到两次系统调用。缺少的链接会当场创建。任何可能丢失数据的情况都会中止启动。
6. `applyMCP` 添加共享的 MCP 服务器（见下文）。
7. `session.Register` 记录会话。ccam 把一个以其 PID 命名、内容为进程启动时间的标记写入 `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`（没有 `XDG_RUNTIME_DIR` 时写入 `~/.local/state/ccam/locks/`）。对于共享对话历史的账号，`historyCheck` 会先在 `history-<group>-<hash>` 中查找正在运行的会话（hash 由当前文件夹计算），并同样在那里注册。`exec` 会保留 PID，因此该标记随后指向的就是 Claude Code 进程本身。之所以使用标记而不是持有锁，是因为 Claude Code 启动的进程会继承锁，并在会话结束后继续持有它。
8. `Plan.Exec` 把终端标题设为 `claude · <name>`（当标准错误是终端时），不经 `fsync` 把启动时间写入 `~/.local/state/ccam/last-used/<name>`，然后调用 `syscall.Exec`。

Claude Code 运行时，`CLAUDE_CONFIG_DIR` 设为账号文件夹，`CCAM_PROFILE` 设为账号名。ccam 会移除这两个变量继承来的值。对于在 macOS 上用 `ccam import` 注册的 `~/.claude` 文件夹，它还会把 `CLAUDE_SECURESTORAGE_CONFIG_DIR` 设为空值，让该账号和直接运行的 `claude` 使用钥匙串中的同一个登录。

### 为什么用 exec

`syscall.Exec` 用 `claude` 替换 ccam 进程，因此 ccam 不会留在内存中，按键、信号和退出码都直接属于 Claude Code，参数始终是一个数组，而不是 shell 字符串。

Windows 没有 `exec`。在 Windows 上，`Plan.Exec` 以子进程方式启动 `claude.exe`，使用相同的环境变量、参数和控制台，等待它结束，并以它的退出码退出。Ctrl+C 会直接从控制台到达 Claude Code，所以 ccam 忽略它。会话标记记录的是 ccam 的 PID，它的存活时间与 Claude Code 完全相同，并在 Claude Code 退出时被删除。等待期间，ccam 不做任何其他事情。会话结束后的工作（例如发现用 `claude mcp add` 添加的服务器）在下次启动时进行：ccam 把 `.claude.json` 的修改时间与 `~/.local/state/ccam/mcp/<name>/seen` 进行比较，只有文件变化时才解析它。

### 命令行如何分流

三条命令行，在 `internal/app` 中走三条路径：

<!-- motion:launch-route -->

快捷命令按 `argv[0]` 匹配，`ccam <name>` 按第一个参数匹配，两者最终都进入 `syscall.Exec`。只有 `ccam usage` 这样的 ccam 命令才会到达 cobra。

### 共享的 MCP 服务器

在默认的 `flag` 策略下，ccam 把组中的 MCP 服务器作为一个参数传给 Claude Code，放在你的参数之前：

<!-- motion:launch-mcp -->

`mcp.FlagArgs` 只有在组的 `mcp.json` 存在时才返回参数，`Plan.InsertArgs` 把它放在 `claude` 之后。你自己的 `--mcp-config` 或 `--strict-mcp-config` 仍会在此基础上生效。它是单个参数 `--mcp-config=<path>`，因为 Claude Code 的 `--mcp-config` 接受多个值：如果写成两个参数，它会把后面的提示也当作配置文件读取。

在 `config.toml` 中设置 `mcp_strategy = "sync"` 时，ccam 改为把服务器写入账号的 `.claude.json`；如果账号已在运行，则跳过这一步。

### Remote Control

`ccam work rc`（或 `ccam work remote-control`）以账号 `work` 启动 Claude Code 的 Remote Control。`claude remote-control` 会拒绝出现在其动词之前的 `--mcp-config`，在动词之后也不接受，而且它的会话只从账号的 `.claude.json` 读取 MCP 服务器。因此 flag 策略在这里的工作方式不同：

<!-- motion:launch-rc -->

1. `remoteControl` 把账号之后第一个不以 `-` 开头的参数当作 Claude Code 的动词。当它是 `remote-control` 或 `rc` 时，`mcp.FlagArgs` 不添加任何内容。
2. `mcp.Synced` 把组的 `mcp.json` 的修改时间与账号的 `sync-stamp` 比较。如果 `.claude.json` 已经包含组当前的服务器，或者组中没有服务器，Remote Control 会直接启动，不会询问。
3. 否则，ccam 会列出会话将缺少的服务器并询问：

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. 回答 `y` 时，`mcp.Sync` 按 `sync` 策略的方式复制这些服务器：持有 Claude Code 的锁，先做备份，并保留账号自己的服务器。回答 `N` 或没有终端时，Remote Control 会在没有这些服务器的情况下启动，ccam 会提示 `mcp_strategy = "sync"`，它会在每次启动前复制这些服务器。

在服务器被复制之前，ccam 每次启动 Remote Control 都会询问；复制之后，只有在组的 `mcp.json` 再次变化时才会询问一次。在之后的普通启动中，`.claude.json` 中的副本优先于 `--mcp-config`，因此 `ccam mcp add`、`remove` 和 `import` 会更新其他账号持有的副本（[共享的 MCP 服务器](/features/mcp)）。

### 10 ms 预算

ccam 在 Claude Code 启动前最多增加 10 ms（p95）。启动路径从不构建 cobra 命令树、启动 Bubble Tea 或打开网络连接；如果各包 `init()` 函数总共耗时超过 5 ms，`TestPackageInitIsCheap` 就会失败。`scripts/perf-check.sh` 用 `hyperfine` 测量 `ccam --dry-run work`、`ccam work` 和一个快捷命令，超过 10 ms 即失败；`internal/launcher` 中的 `BenchmarkLaunchPath` 会循环运行这些步骤。要查看每一步的耗时，请设置 `CCAM_TRACE=1`：

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 安全性与边界情况

- 未知的名称或编号会以 `no account "<name>"` 停止，并提示使用 `ccam list`。如果一个账号都没有，ccam 会请你用 `ccam add <name>` 添加一个。
- 如果某个共享路径变成了普通文件或指向错误的位置，ccam 会以 `stopped because of share link problems:` 停止，列出这些问题，并请你运行 `ccam doctor --fix`。它不会在启动路径上移动数据。
- 如果另一个与此账号共享对话历史（同一组）的 Claude Code 正在同一文件夹中运行，ccam 会警告历史记录可能混在一起，并询问 `Run anyway? [y/N]`。没有终端时，它会显示警告并继续。
- `--dry-run` 会跳过 MCP 同步、新服务器提示、历史检查、会话标记、标题和最近使用时间，也不会执行 exec。它仍会检查共享链接，并像正常启动一样创建缺少的链接。
- 如果在 `work` 已在运行时执行 `ccam work rc`，ccam 不会改动 `.claude.json`，会说明这一点，并在没有组中服务器的情况下启动 Remote Control；会话退出后再运行一次即可复制它们。使用 `--dry-run` 时，它只打印这条消息，既不询问也不复制。
- 使用 `mcp_strategy = "sync"` 时，同步失败会打印 `ccam: MCP sync failed (launching anyway): …`，Claude Code 会以它原有的服务器启动。

## 相关内容

- [快捷命令](/features/shortcuts)
- [用 ccam 使用 Remote Control](/guides/remote-control)
- [管理账号](/features/accounts)
- [命令参考](/docs/commands)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [CLAUDE_CONFIG_DIR 的工作原理](/guides/claude-config-dir)
