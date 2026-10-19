---
title: 用 ccam 使用 Remote Control
description: 用 ccam <name> rc 以任意账号启动 Claude Code 的 Remote Control：哪些选项适用、ccam 为什么询问是否复制共享的 MCP 服务器，以及常见问题的解决方法。
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

`ccam work rc` 以账号 `work` 启动 Claude Code 的 Remote Control，这样你就可以在 claude.ai/code 或 Claude 移动应用中操作这台电脑上的会话。它的用法与 `claude rc` 相同，只有一处区别：如果该账号共享 MCP 服务器，ccam 会询问一次是否把它们复制到账号的 `.claude.json` 中，因为 Remote Control 的会话不会从其他任何地方读取服务器。

## 启动

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[快捷命令](/features/shortcuts)的用法相同：`ccwork rc` 就是 `ccam work rc`。在你要工作的文件夹中运行命令，然后打开 claude.ai/code 或 Claude 应用，并用与 `work` 相同的 Claude 账号登录。

## 选项

这里会用到两类选项：ccam 自己的选项，写在账号名之前或放在 ccam 的设置中；以及 Remote Control 的选项，写在 `rc` 之后，原样传给 Claude Code。

### ccam 选项和设置

| 选项或设置 | 作用 | 默认值 |
| --- | --- | --- |
| `--dry-run`（写在账号之前） | 打印环境变量和命令，有 MCP 提示时也一并打印。既不询问也不复制 | 关 |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | 当账号共享 MCP 服务器、而其 `.claude.json` 中没有组的当前列表时询问。`y` 复制，`N` 或回车则不复制直接启动 | `N` |
| `~/.config/ccam/config.toml` 中的 `mcp_strategy = "sync"` | 每次启动前都把组的服务器复制到 `.claude.json`，无论是否为 Remote Control，因此从不询问 | `"flag"` |
| `ccam share work --mcp off` | 停止与该账号共享 MCP 服务器；之后 Remote Control 只带账号自己的服务器启动，也不会询问 | 开 |

### Remote Control 选项

`rc` 之后的所有内容都交给 Claude Code。下表是 Claude Code 2.1.304 中 `claude remote-control` 的选项；`ccam work rc --help` 会列出你所用版本的选项。

| 选项 | 作用 |
| --- | --- |
| `--name <name>` | 在 claude.ai/code 中显示的会话名称 |
| `--remote-control-session-name-prefix <prefix>` | 自动生成的会话名称的前缀（默认：主机名） |
| `-c`, `--continue` | 重新连接 Remote Control 为此文件夹最后记录的会话，约 4 小时内有效 |
| `--session-id <id>` | 重新连接指定的会话 |
| `--permission-mode <mode>` | 新会话的权限模式：`acceptEdits`、`auto`、`bypassPermissions`、`default`、`dontAsk` 或 `plan` |
| `--spawn <mode>` | `same-dir`（默认）、`worktree`（每个会话一个 git worktree）或 `session`（只有一个会话，会话结束即退出） |
| `--capacity <N>` | `same-dir` 和 `worktree` 模式下同时运行的最大会话数（默认 32） |
| `--[no-]create-session-in-dir` | 启动时在当前文件夹中创建第一个会话（默认开） |
| `--[no-]chrome` | 在新会话中使用 Claude in Chrome（默认关） |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | 调试输出和详细输出 |

把 `rc` 紧跟在账号名之后。ccam 会把第一个不以 `-` 开头的参数当作 Claude Code 的命令，所以在 `ccam work rc --name laptop` 中它看到的是 `rc`，而在 `ccam work --model opus rc` 中看到的是 `opus`。

## ccam 做了什么

启动沿用平常的流程（[启动 Claude Code](/features/launch)），只有 MCP 这一步不同：

1. 与任何一次启动一样，ccam 找到账号，检查其共享链接，并解析 `claude` 的位置。
2. 它发现命令是 `rc` 或 `remote-control`，于是不加 `--mcp-config`。
3. 如果账号共享 MCP 服务器，它会把组的 `mcp.json` 与 `~/.local/state/ccam/mcp/work/` 中的 `sync-stamp` 进行比较。如果 `.claude.json` 已经包含当前的服务器，或者组中没有服务器，就直接跳到第 6 步。
4. 否则，它会列出会话将缺少的服务器：

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. 回答 `y` 时，它会先备份 `.claude.json`，持有 Claude Code 自己的锁，把组的服务器加入 `mcpServers`，并保留账号的其他服务器。回答 `N`、没有终端，或账号正在运行时，它不做任何更改，并给出提示：

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. 它记录会话，把 `CLAUDE_CONFIG_DIR` 设为账号文件夹，然后通过 `exec` 把自己替换为 `claude rc`。

之后 Remote Control 会作为服务器持续运行。它为你启动的每个会话都是这个进程的子进程，因此都以 `work` 的身份运行：使用 `work` 的登录信息、设置和对话历史，并计入 `work` 的用量限制。服务器运行期间，ccam 会把该账号视为正在运行。

试运行可以看出区别。普通启动会把组的服务器作为参数传入，Remote Control 则不会：

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## 为什么这样设计

**Remote Control 不接受 `--mcp-config`。** 普通启动时，ccam 以 `--mcp-config=<path>` 的形式把组的服务器交给 Claude Code，从不写入 Claude Code 的文件。`claude remote-control` 会拒绝出现在其命令名之前的这个选项，在命令名之后也不接受，所以在 ccam 0.6.0 之前，凡是共享 MCP 服务器的账号，`ccam work rc` 都会失败。现在 ccam 不再加这个选项。

**它的会话只从 `.claude.json` 读取服务器。** Remote Control 启动的会话不会继承命令行中的选项。它们找到用户 MCP 服务器的唯一位置是账号的 `.claude.json`。这个文件还保存着账号的登录信息和各文件夹的状态，所以 ccam 不能像处理 `settings.json` 那样把它链接到共享副本。要让会话拿到共享列表，唯一的办法就是把服务器复制进这个文件。

**ccam 写入前会先询问。** 在默认的 flag 策略下，ccam 承诺不改动 `.claude.json`。Remote Control 是唯一必须写入该文件的情况，所以它会先询问、保留备份，并且只改动 `mcpServers`。

**账号运行时不写入。** Claude Code 会在会话期间重写 `.claude.json`，从外部做的更改会丢失或被覆盖。因此只要该账号有任何会话在运行，包括另一个 Remote Control 服务器，ccam 就会跳过复制。

**每次变化只问一次。** 复制之后，`sync-stamp` 会记录组的 `mcp.json` 的修改时间。只有当该文件发生变化时（例如在 `ccam mcp add` 之后），ccam 才会再次询问。如果你回答 `N`，由于还没有副本，每次启动 Remote Control 都会询问。

**`ccam mcp` 会让副本保持最新。** 名称冲突时，`.claude.json` 中的条目优先于 `--mcp-config`，所以旧的副本在普通启动中也会遮住较新的定义。因此 `ccam mcp add`、`add-json`、`remove` 和 `import` 会更新所有未在运行的账号中的副本：

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**没有终端时，答案就是“否”。** Remote Control 常常在脚本、服务或 `nohup` 中运行，没有人能回答。此时 ccam 不会等待，而是不带这些服务器直接启动。

## 问题与解决方法

### Remote Control 会话中缺少共享的 MCP 服务器

原因是询问时回答了 `N`、当时没有终端，或者账号正在运行。停止 Remote Control，关闭该账号的其他会话，然后在终端中重新运行 `ccam work rc`，并回答 `y`。如果希望永远不被询问，请在 `~/.config/ccam/config.toml` 中设置 `mcp_strategy = "sync"`。已打开会话中的服务器在运行期间不会改变，所以复制之后请启动新的会话。

### "Account work is running, so its .claude.json is left alone"

`work` 仍有会话在运行：可能是普通的 `ccam work`、快捷命令，或另一个 Remote Control 服务器。退出它，然后重新运行 `ccam work rc`。其他账号运行的会话不算在内。

### ccam 每次都询问

你还没有复制这些服务器。回答一次 `y`，或者设置 `mcp_strategy = "sync"`。如果你根本不想在 Remote Control 中使用这些服务器，就停止与该账号共享 MCP：`ccam share work --mcp off`。

### `ccam work rc` 立即失败，并报出与 `--mcp-config` 有关的错误

你的 ccam 早于 0.6.0，会把 `--mcp-config` 传给 Remote Control。请运行 `ccam update`。如果 ccam 已是最新版本仍出现同样的错误，说明在 `rc` 之前有一个带值的选项，例如 `ccam work --model opus rc`，于是 ccam 把 `opus` 当成了命令；请把 `rc` 紧跟在账号名之后。

### 用 `ccam mcp` 做的更改没有反映到 Remote Control

`ccam mcp add` 和 `remove` 只更新未在运行的账号中的副本，并会提示：`Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` 停止 Remote Control 后重新启动；ccam 会发现 `mcp.json` 已经变化，并再次询问。

### 复制替换了我自己的服务器

如果账号自己有一个与组中服务器同名的服务器，复制会替换它。之前的 `.claude.json` 保存在 `~/.local/state/ccam/mcp/work/backups/` 中，那里保留最新的五份。若想让所有账号都使用你的版本，请在下次启动 Remote Control 之前，用 `ccam mcp import work <server> --force` 把它移入组中。

### 我想删除复制过去的服务器

`ccam mcp remove <server>` 会把服务器从组和每个副本中删除。若要回到复制之前的状态，请在该账号没有任何会话运行时，从 `~/.local/state/ccam/mcp/work/backups/` 恢复最新的备份。

### Remote Control 询问是否信任该文件夹，或在没有终端时直接退出

文件夹信任是按账号保存在 `.claude.json` 中的，所以你以 `personal` 身份信任过的文件夹，对 `work` 来说是新的。先在该文件夹中运行一次 `ccam work` 并接受信任，然后再在那里启动 Remote Control。

### "You must be logged in with a Claude account that has a subscription"

Remote Control 需要使用 Claude 订阅账号登录。在 `ccam list` 中查看该账号的电子邮件，必要时用 `ccam login work` 重新登录。

### 手机上看不到会话

Claude 应用或 claude.ai/code 登录的是另一个 Claude 账号。把它与 `ccam list` 为 `work` 显示的电子邮件对比一下。

### ccam 警告另一个 Claude Code 在此文件夹中共享对话历史

同一历史组中的另一个账号正在同一文件夹中运行。如果两个账号的会话继续同一段对话，历史可能会混在一起。只有在它们处理不同对话时才回答 `y`，否则请在另一个文件夹中启动 Remote Control。

## 相关内容

- [启动 Claude Code](/features/launch)
- [共享的 MCP 服务器](/features/mcp)
- [在账号之间共享](/features/share)
- [管理账号](/features/accounts)
