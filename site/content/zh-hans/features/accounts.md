---
title: 管理账号
description: 用 ccam add、ccam login、ccam list、ccam rename、ccam remove 和 ccam import 添加、登录、列出、重命名、删除和导入 Claude Code 账号。
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

每个 ccam 账号都是一个独立的 Claude Code 配置文件夹，拥有自己的登录。这些命令用于创建、登录、列出、重命名和删除账号，以及注册你已有的文件夹。

## 用法

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

凡是命令接受 `<name|number>` 的地方，你都可以给出账号名或它的槽位号，即 `ccam list` 中的 `#` 列。

在终端中不带选项运行 `ccam add` 会打开向导。它先询问名称和账号的登录方式（现在用 Claude 登录、使用 API 密钥，或稍后再用 Claude 登录），然后显示可加入的组，其中 `default` 已勾选（还没有任何组时，显示的是组 `default` 将要共享的项），接着询问 `Create a shortcut command for work? [y/N]`，最后显示一份摘要供你确认。按 `Esc` 返回上一步。只要给出 `--group`、`--no-group`、`--shortcut` 或 `--no-login` 中的任何一个，或者输入来自管道而不是终端，就会跳过向导；此时必须指定名称，并且除非 `--shortcut` 指定了名称，否则不会创建快捷命令。

不带 `--group` 时，新账号会加入组 `default`。第一个账号会创建这个组，共享 MCP 服务器、设置、`CLAUDE.md` 和扩展；参见[用组共享](/features/share)。

`ccam list` 显示槽位、名称、快捷命令、邮箱、组和最近使用时间，默认账号前面标有 `*`：

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

使用 `--json` 时，每个账号是一个对象，包含 `slot`、`shortcut`、`name`、`email`、`default`、`dir`、`last_used`（账号从未启动过时省略）、`groups` 和 `share`。`share` 把每一项（`mcp`、`settings`、`memory`、`extensions`、`history`）映射到提供它的组，或映射到 `"local"`。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `ccam add --no-login` | 创建账号但不登录；之后用 `ccam login <name>` 登录 | 关闭 |
| `ccam add --group <group>` | 加入这个组，而不是 `default`。可重复使用；`default` 不存在时会被创建 | `default` |
| `ccam add --no-group` | 不加入任何组：账号不共享任何内容 | 关闭 |
| `ccam add --shortcut <command>` | 同时在 `~/.local/bin` 中创建这个快捷命令 | 无 |
| `ccam add --api --key-from <source>` | 使用 API 密钥而不是登录的账号；参见 [API 密钥账号](/features/api) | 关闭 |
| `ccam list --json` | 以 JSON 输出账号 | 关闭 |
| `ccam remove -y`、`--yes` | 不询问直接删除 | 关闭 |
| `ccam import --name <name>` | 你所指定文件夹的账号名 | 由文件夹名推导 |
| `ccam import -y`、`--yes` | 不询问直接导入找到的所有文件夹 | 关闭 |
| `--dry-run` | `ccam list` 和 `ccam completion` 接受它，因为它们不做任何更改。`add`、`login`、`logout`、`use`、`rename`、`remove` 和 `import` 会拒绝它，并输出 `` --dry-run is not supported by `ccam add`; nothing was changed ``（其中是对应的命令名） | 关闭 |

`ccam login`、`ccam logout`、`ccam use` 和 `ccam rename` 没有自己的选项。`ccam completion` 接受 `bash`、`zsh` 或 `powershell`；其他值会失败并显示 `only bash, zsh and powershell are supported`。

账号名由 1 到 32 个 `a-z`、`0-9`、`-` 和 `_` 字符组成，必须以字母或数字开头。名称不能全是数字，因为 `ccam 2` 表示槽位 2。名称也不能是 ccam 的命令（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`folder`、`pick`、`ext`、`group`、`groups`），这样 `ccam <name>` 永远不会有歧义。这份列表是 `internal/config` 中的 `config.Reserved`。

## 工作原理

`ccam add work --shortcut ccwork` 在 `addProfile`（`internal/cli/profiles.go`）中执行以下步骤：

<!-- motion:accounts-add -->

1. 在创建任何内容之前，`config.ValidateName` 检查名称，`shim.Check` 检查快捷命令。
2. `profile.Add` 以 `0700` 权限创建 `~/.local/share/ccam/profiles/work/`，并为账号分配最小的空闲槽位号（`Config.NextSlot`），因此已删除账号的编号会被重复使用。第一个账号会成为默认账号。
3. 账号加入它的组（除非 `--group` 或 `--no-group` 另有指定，否则为 `default`；第一个账号会创建 `default`），`applyShare` 先用 `share.Build` 再用 `Apply`，把新建的空文件夹链接到这些组，例如把 `settings.json` 链接到 `groups/default/settings.json`。如果这一步失败，新文件夹会被删除，因为残留的文件夹会让每次重试都因 `folder already exists` 而失败。
4. 账号通过 `fsutil.AtomicWrite` 保存到 `~/.config/ccam/config.toml`：先写临时文件，执行 `fsync`，再 `rename`，因此即使崩溃也不会留下写了一半的配置。
5. `runLogin` 在 `CLAUDE_CONFIG_DIR` 设为账号文件夹的情况下启动 `claude auth login`，并用 `--email` 传入账号上一次的邮箱，这样登录页面会预先填好。如果登录失败，账号仍会保留，ccam 会打印错误，并提示你稍后运行 `ccam login work`。
6. `profile.ReadIdentity` 从 `profiles/work/.claude.json` 中取出 `oauthAccount`：邮箱、账号 ID 和组织 ID、组织名称以及套餐（`claude_max`、`claude_team` …）。这个文件可能有好几 MB，所以 `TopLevelValue` 只扫描这一个键，而不解析整个文件。ccam 把它们保存在 `config.toml` 的 `[profile.login]` 中。
7. `shim.Install` 把 `~/.local/bin/ccwork` 链接到 ccam 可执行文件，后者启动时会在 `config.toml` 中查找 `argv[0]`。

登录可能需要几分钟，期间另一个 ccam 可能会保存配置。因此需要等待的命令通过 `env.update` 保存：它重新加载 `config.toml`，只应用这一项更改（这里是邮箱），然后保存。出于同样的原因，`ccam remove` 在提示之后也会重新加载。

`ccam login` 对已有账号执行第 5 步和第 6 步。Claude Code 只把邮箱当作登录页面的提示，所以一次登录可能进入与该 ccam 账号原先不同的 Claude 账号；ccam 会把新的登录与 `[profile.login]` 比较，并指出这一点（`ccam: this login is other@company.com, but work was dev@company.com before.`）。它会保留新的登录，并告诉你如何撤销：先 `ccam logout work`，再 `ccam login work`，然后用正确的邮箱登录。当另一个 ccam 账号登录到同一个 Claude 账号时，ccam 会警告这两个账号共享用量限制。`ccam login --expired` 会依次登录每个从未登录过、或登录（其刷新令牌）已过期的账号。

`ccam logout` 在账号中运行 `claude auth logout` 并清除保存的邮箱；登录记录会保留，因此下一次登录会与它进行比较。`ccam use` 只改写 `config.toml` 中的 `default`。

当账号属于不止一个组织时，`ccam list` 会增加一列 `Organization`（个人套餐显示为 `personal`），`ccam list --json` 则为每个账号提供 `plan` 和 `org`。

### 在远程机器上登录

通过 SSH 连接时，`claude auth login` 中浏览器那一半可能会出问题。Claude Code 在服务器上的 `http://localhost:<port>/callback` 回调服务器上等待登录完成，但完成登录的浏览器运行在你的电脑上，在那里这个地址指向不了任何地方，浏览器会显示 `ERR_CONNECTION_REFUSED` 之类的错误。这个地址中仍带有一次性代码，所以 ccam 可以把它转交过去：

<!-- motion:accounts-relay -->

`runLogin` 不会把你的终端直接交给 Claude Code。它把 Claude Code 的标准输入连接到一个管道，由 `relayLogin` 把你输入的每一行复制进去，因此在 Claude Code 的 `Paste code here` 提示处粘贴的代码依然能送达。当某一行是回调地址时，`callbackURL` 只在以下条件都满足时才接受它：主机是 `localhost`、`127.0.0.1` 或 `[::1]`，指定了端口，路径为 `/callback`，并且同时带有 `code` 和 `state`；`http://` 可以省略。然后 `deliverCallback` 像浏览器那样，在服务器上请求这个地址，时限为 10 秒，且不跟随重定向。请求到达正在等待的 `claude auth login` 的回调服务器（图中的虚线）。Claude Code 收到代码并保存登录，在 Linux 上保存到账号文件夹中的 `.credentials.json`。粘贴的内容永远不会让 ccam 连接其他主机。

登录开始之前，ccam 会打印：

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

成功时打印 `→ Passed the sign-in to Claude Code on this machine.`

登录在终端上等待输入时，ccam 会关闭终端的回显和行模式，你每输入一个字符就显示一个 `*`，因此粘贴的代码或回调地址不会出现在屏幕上。行编辑由 `readLine` 自己处理：`Enter` 发送这一行，`Backspace` 删除一个字符，在空行上按 `Ctrl+D` 结束输入。登录结束时会恢复终端设置；收到 `Ctrl+C`、`SIGTERM` 或 `SIGHUP` 时也会恢复，它们仍会结束 ccam。如果标准输入和标准错误不是终端，各行会按原样复制。

### 重命名、删除和导入

`ccam rename` 和 `ccam remove` 作用于账号的文件夹，以及 ccam 在其他位置以该账号名保存的内容。`nameState` 列出了这些路径：`~/.local/state/ccam/last-used/<name>`、`~/.local/state/ccam/mcp/<name>`，以及用量缓存 `~/.local/state/ccam/usage/<name>.json` 和 `<name>.statusline.json`。`ccam remove work` 按以下顺序处理它们：

<!-- motion:accounts-remove -->

如果账号正在运行，`ccam remove` 会先拒绝，否则询问确认。如果账号共享扩展，`RewritePluginPaths` 会改写组中的 `installed_plugins.json` 和 `known_marketplaces.json`，因为 Claude Code 保存的是插件的绝对路径，这些路径可能经过该账号的文件夹。`profile.Trash` 把文件夹移到 `~/.local/share/ccam/trash/<name>-<time>/`，账号从 `config.toml` 中移除，`nameState` 中的路径被删除，以免之后同名的账号继承它们，快捷命令的链接也会被删除。要恢复账号，把文件夹从回收站移回即可。

`ccam rename` 重命名文件夹，把其 `plugins/` 中的插件路径从旧文件夹改写为新文件夹，保留默认账号设置，并把 `nameState` 中的路径移到新名称下。快捷命令保持原名。如果插件路径改写失败，文件夹会被改回原名，因为配置不会被保存。

`ccam import` 把文件夹在原位置注册为 `config.toml` 中的 `dir`。它从不移动或复制文件夹：复制的刷新令牌会让各个副本互相把对方登出。不指定文件夹时，它会查找尚未注册的 `~/.claude` 和 `~/.claude-*` 文件夹并逐个询问；名称取自文件夹（`~/.claude` 为 `main`，`~/.claude-work` 为 `work`，`~/.claude-2` 为 `acct2`），无法推导名称的文件夹会被列出，并附上应使用的 `--name` 命令。没有 `CLAUDE_CONFIG_DIR` 时，Claude Code 把 `.claude.json` 放在你的主目录中、`~/.claude` 旁边，因此当你导入 `~/.claude` 而它没有自己的 `.claude.json` 时，ccam 会创建链接 `~/.claude/.claude.json`，指向 `~/.claude.json`。这样直接运行 `claude` 和运行 `ccam main` 看到的是同一个账号。导入的账号不属于任何组，直到你把它加入某个组（`ccam group join <group> <name>`），也不会创建快捷命令。

`ccam completion bash|zsh` 会打印一个补全脚本，它也能补全账号名。

## 安全性与边界情况

- 账号运行时，`ccam rename` 和 `ccam remove` 会拒绝执行（`refuseRunning`），并显示 `work is running; quit its Claude Code session first`，因为会话会继续写入该文件夹。
- 对导入的账号执行 `ccam remove` 只会取消注册：`✓ main unregistered (folder ... left in place)`。你的 `~/.claude` 不会被移到回收站。
- 如果 `profiles/` 中已存在同名文件夹，`ccam add` 会拒绝该名称，并提示使用 `ccam import`。
- `ccam import` 会拒绝已注册的文件夹、不是文件夹的路径，以及已被占用的名称。
- `ccam logout` 和 `ccam login` 不检查账号是否正在运行。
- 如果转交无法到达 Claude Code（例如登录已经结束），它会打印 `ccam: could not pass the sign-in to Claude Code:` 以及原因。HTTP 错误通常表示代码已过期；请重新开始登录。

## 相关内容

- [启动 Claude Code](/features/launch)
- [快捷命令](/features/shortcuts)
- [共享 MCP 服务器、设置和历史记录](/docs/sharing)
- [多个 Claude Code 账号](/guides/multiple-claude-code-accounts)
