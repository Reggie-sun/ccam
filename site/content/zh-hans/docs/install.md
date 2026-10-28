---
title: 安装 ccam
metaTitle: "在 Linux、macOS 和 Windows 上安装 ccam"
description: 用一条命令在 Linux、macOS 或 Windows 上安装 ccam，也可以使用 .deb 软件包或 go install 安装，并用 ccam update 保持最新。需要先安装 Claude Code。
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam 会为你启动 Claude Code，所以请先安装 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)，并确认 `claude --version` 可以正常运行。

在 Linux 或 macOS 的终端中运行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

脚本会下载适用于你系统的版本，校验其 SHA-256 校验和，然后把 `ccam` 放入 `~/.local/bin`；之后你选择的快捷命令（例如 `ccwork`）也会由 ccam 创建在这个文件夹中。如果该文件夹不在你的 `PATH` 中，脚本会打印需要添加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

在 Windows 上，请在 PowerShell 中运行：

```powershell
irm https://ccam.dev/install.ps1 | iex
```

它会下载 `ccam_windows_<arch>.zip`，校验其 SHA-256 校验和，然后把 `ccam.exe` 放入 `%USERPROFILE%\.local\bin`，也就是 Claude Code 安装程序存放 `claude.exe` 的文件夹。如果该文件夹还不在你的用户 `PATH` 中，脚本会把它加进去，并保留 `%USERPROFILE%\bin` 这类条目原样不变。之后请打开一个新的终端。

## 支持的系统

- x86-64 和 ARM64 上的 Linux：Ubuntu、Debian、Fedora、Arch、Alpine 及其他发行版。二进制文件是静态的，没有任何依赖。
- Apple 芯片和 Intel 处理器上的 macOS。
- x86-64 和 ARM64 上的 Windows 10 和 11，可在 PowerShell、Windows Terminal 或 Git Bash 中使用：
  - 共享文件夹无需开启开发人员模式即可使用。
  - 共享文件（`settings.json`、`CLAUDE.md`、`history.jsonl`）需要开启开发人员模式（设置 > 系统 > 开发者选项），因为 Windows 只允许在该模式下或由管理员创建指向文件的链接；参见[在 Windows 上共享](/features/share#on-windows)。
  - 未开启开发人员模式时，第一次 `ccam add` 创建的 `default` 组只包含 MCP 服务器和扩展，并说明开启开发人员模式后如何添加其余各项。
- 在 WSL 中，ccam 的使用方式与在其他任何 Linux 系统上相同。

## 其他安装方式

**Debian 和 Ubuntu 软件包。** 从[发布页面](https://github.com/jhyunwoo/ccam/releases)下载 `ccam_<version>_linux_amd64.deb` 并安装：

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

该软件包会安装 `/usr/bin/ccam`。

**从源码安装。** 需要 Go 1.26 或更高版本：

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

如果你安装的是发布标签（`@latest` 或 `@vX.Y.Z`），ccam 会报告该版本号，`ccam update` 和仪表盘提供的更新都能正常使用。从 `@main` 构建的版本，或从在某个标签之后还有提交的本地检出构建的版本，会报告 `dev`，除非你传入 `--version`，否则 `ccam update` 会拒绝更新它。

**安装特定版本。** 设置 `CCAM_VERSION`：

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell 补全

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

补全覆盖 `ccam` 的所有命令和选项，包括账号名。

## Claude Code 插件

ccam 插件让 Claude Code 自己学会使用 ccam：查看哪个账号还有余量，用 `ccam pick` 在不同账号上启动并行的 `claude -p` 任务，以及用 `ccam ext` 解释为什么缺少某个技能或 MCP 服务器。在 Claude Code 中安装它：

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

它会添加：

- 一个技能，它会建议登录或移动文件的命令，但从不自己运行它们；
- `/ccam:usage`、`/ccam:ext` 和 `/ccam:switch` 这几个命令；
- 两个钩子（[详情](/features/next#inside-claude-code)）。一个在会话因用量限制而停止时显示一条终端通知；在 Linux 和 macOS 上，开启 [`auto_switch`](/features/next#auto-switch) 后，ccam 随后会把会话移到另一个还有余量的账号。另一个让 `/ccam:switch [account]` 把对话交给另一个账号。

这些钩子会运行 `ccam`，所以它必须位于 Claude Code 看到的 `PATH` 中。

## 检查安装

```sh
ccam version
ccam doctor
```

`ccam doctor` 会检查能否找到 Claude Code，以及每个账号文件夹、链接和登录是否正常。接下来请阅读[快速入门](/docs/quick-start)。

## 更新

```sh
ccam update --check
ccam update
```

`ccam update --check` 只告诉你是否有更新的版本。`ccam update` 会下载适用于你系统的版本，用该版本的 `checksums.txt` 校验其 SHA-256 校验和，然后一步替换 `ccam` 文件。如果校验和不匹配，不会做任何更改。`ccwork` 这样的快捷命令是指向 ccam 的链接，因此会立即运行新版本。

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

仪表盘（在终端中不带参数运行 `ccam`）也会在每次打开时检查新版本，有新版本时会提供更新：按 `U` 并确认即可。其他命令都不会检查更新，ccam 也从不未经询问就更新。要安装特定版本（包括较旧的版本），请传入 `--version`。`ccam --dry-run update` 会显示将要下载和替换的内容，但不做任何更改。

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

重新运行安装脚本也会更新 `~/.local/bin` 中的 ccam。如果你是从发布标签用 `go install` 构建的 ccam，可以运行 `ccam update`，也可以再运行一次同样的 `go install` 命令。

**Windows。** `ccam update` 会下载适用于 Windows 的 `.zip`。在 Windows 上无法替换正在运行的程序，所以旧的二进制文件会被重命名为 `ccam.exe.old`，下一次更新时会将其删除；如果仍有会话在使用更早的 `ccam.exe.old`，则名称为 `ccam.exe.old-<time>`。属于硬链接的快捷命令会紧接着被更新。再次运行 `install.ps1` 的效果相同，即使仍有 Claude Code 会话在使用旧的二进制文件也可以；它更新属于硬链接的快捷命令时会打印 `OK 2 shortcuts now run the new ccam`。

**Debian 和 Ubuntu 软件包。** `ccam update` 不会替换由 apt 安装的 `/usr/bin/ccam`，而是打印以下步骤。从[发布页面](https://github.com/jhyunwoo/ccam/releases)下载新的 `ccam_<version>_linux_amd64.deb`，覆盖安装旧版本：

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## 卸载

`ccam uninstall` 会删除 ccam 及其创建的一切。它会先列出将要删除的内容，并征求确认：

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- 每个账号都会像 `ccam remove` 那样被删除。位于 ccam 文件夹中的账号会先退出登录，然后连同其文件夹一起删除，其中的对话历史和设置也会一并删除。
- 你[导入](/docs/commands)的 `~/.claude` 文件夹会保留在原位置，供直接运行的 `claude` 使用。如果它在某个组中，它会拿回一份它共享过的全部内容的副本。它的 `settings.json` 中的 ccam 状态栏会让位给 ccam 所包装的原状态栏（如果有）。如果它装有 ccam 插件，ccam 会打印需要在 Claude Code 中运行的 `/plugin uninstall` 命令。
- 快捷命令、`~/.config/ccam`、`~/.local/share/ccam`（包括 ccam 的回收文件夹）和 `~/.local/state/ccam` 会被删除，最后删除 `ccam` 命令本身。

只要有账号正在运行，ccam 就会拒绝执行；请先退出该账号的 Claude Code 会话。由 apt 安装的 `/usr/bin/ccam` 会保留；之后请用 `sudo apt remove ccam` 删除它。在 Windows 上，正在运行的程序无法删除自身，因此 `ccam uninstall` 最后会打印删除 `ccam.exe` 的 `Remove-Item` 命令。
