---
title: 安装 ccam
description: 用一条命令在 Linux 或 macOS 上安装 ccam，也可以使用 .deb 软件包或从源码安装。需要先安装 Claude Code。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

在 Linux 或 macOS 的终端中运行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

脚本会下载适用于你系统的版本，校验其 SHA-256 校验和，然后把 `ccam` 放入 `~/.local/bin`；之后 `ccam add` 会在同一文件夹中创建 `claude1`、`claude2`… 快捷命令。如果该文件夹不在你的 `PATH` 中，脚本会打印需要添加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

ccam 会为你启动 Claude Code，所以请先安装 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)，并确认 `claude --version` 可以正常运行。

## 支持的系统

- x86-64 和 ARM64 上的 Linux：Ubuntu、Debian、Fedora、Arch、Alpine 及其他发行版。二进制文件是静态的，没有任何依赖。
- Apple 芯片和 Intel 处理器上的 macOS。
- 不直接支持 Windows。在 WSL 中，ccam 的使用方式与在任何 Linux 系统上相同。

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

## Shell 补全

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

补全覆盖 `ccam` 的所有命令和选项，包括账号名。

## 检查安装

```sh
ccam version
ccam doctor
```

`ccam doctor` 会检查能否找到 Claude Code，以及每个账号文件夹、链接和登录是否正常。接下来请阅读[快速入门](/docs/quick-start)。

## 卸载

删除 ccam 命令和账号快捷命令；如果不再需要这些账号，再删除数据文件夹：

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

删除 `~/.local/share/ccam` 也会删除其中存储的登录信息和对话历史。你[导入](/docs/commands)的 `~/.claude` 文件夹会保留在原位置。
