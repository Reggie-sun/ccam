---
title: 安装 ccam
metaTitle: "在 Linux、macOS 和 Windows 上安装 ccam"
description: 用一条命令在 Linux、macOS 或 Windows 上安装 ccam，也可以使用 .deb 软件包或从源码安装，并用 ccam update 保持最新。需要先安装 Claude Code。
updated: 2026-10-21
order: 1
translation: machine
sourceHash: "1e3a53623de0"
---

在 Linux 或 macOS 的终端中运行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

脚本会下载适用于你系统的版本，校验其 SHA-256 校验和，然后把 `ccam` 放入 `~/.local/bin`；之后你选择的快捷命令（例如 `ccwork`）也会由 ccam 创建在这个文件夹中。如果该文件夹不在你的 `PATH` 中，脚本会打印需要添加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

在 Windows 上，请在 PowerShell 中运行：

```powershell
irm https://ccam.dev/install.ps1 | iex
```

它会下载 `ccam_windows_<arch>.zip`，校验其 SHA-256 校验和，把 `ccam.exe` 放入 `%USERPROFILE%\.local\bin`（Claude Code 安装程序存放 `claude.exe` 的文件夹），如果该文件夹还不在你的用户 `PATH` 中，就把它加进去。之后请打开一个新的终端。

ccam 会为你启动 Claude Code，所以请先安装 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)，并确认 `claude --version` 可以正常运行。

## 支持的系统

- x86-64 和 ARM64 上的 Linux：Ubuntu、Debian、Fedora、Arch、Alpine 及其他发行版。二进制文件是静态的，没有任何依赖。
- Apple 芯片和 Intel 处理器上的 macOS。
- x86-64 和 ARM64 上的 Windows 10 和 11，可在 PowerShell、Windows Terminal 或 Git Bash 中使用。共享文件夹对所有用户都可用。共享文件（`settings.json`、`CLAUDE.md`、`history.jsonl`）需要开启开发人员模式（设置 > 系统 > 开发者选项），因为 Windows 只允许在该模式下或由管理员创建指向文件的链接；参见[在 Windows 上共享](/features/share#on-windows)。
- 在 WSL 中，ccam 的使用方式与在任何 Linux 系统上相同。

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

它会添加一个技能、`/ccam:usage`、`/ccam:ext` 和 `/ccam:switch` 这几个命令，以及两个钩子：当会话因用量限制而停止时，你会收到一条终端通知，`/ccam:switch [account]` 会把对话交给另一个账号（[详情](/features/next#inside-claude-code)）。这些钩子会运行 `ccam`，所以它必须位于 Claude Code 看到的 `PATH` 中。该技能从不运行登录或移动文件的命令，而是建议你去运行。

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

仪表盘（在终端中运行 `ccam`）也会在每次打开时查找新版本，有新版本时会提供更新：按 `U` 并确认即可。其他命令都不会查找，ccam 也从不未经询问就更新。要安装特定版本（包括较旧的版本），请传入 `--version`。`ccam --dry-run update` 会显示将要下载和替换的内容，但不做任何更改。

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

重新运行安装脚本也会更新 `~/.local/bin` 中的 ccam。如果你是用 `go install` 构建的 ccam，再运行一次同样的命令即可。

**Windows。** `ccam update` 会下载适用于 Windows 的 `.zip`。在 Windows 上无法替换正在运行的程序，所以旧的二进制文件会被重命名为 `ccam.exe.old`，下一次更新时会将其删除。

**Debian 和 Ubuntu 软件包。** `ccam update` 不会替换由 apt 安装的 `/usr/bin/ccam`，而是打印以下步骤。从[发布页面](https://github.com/jhyunwoo/ccam/releases)下载新的 `ccam_<version>_linux_amd64.deb`，覆盖安装旧版本：

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## 卸载

趁 ccam 还在，先删除账号的快捷命令。`ccam list` 会在 Shortcut 列中显示它们；用 `ccam shortcut <account> --remove` 逐个删除，或者自己从 `~/.local/bin` 中删除这些链接。然后删除 ccam 命令；如果不再需要这些账号，再删除数据文件夹：

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

在 Windows 上，在 PowerShell 中运行：

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

删除 `~/.local/share/ccam` 也会删除其中存储的登录信息和对话历史。你[导入](/docs/commands)的 `~/.claude` 文件夹会保留在原位置。
