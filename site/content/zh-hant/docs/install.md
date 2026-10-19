---
title: 安裝 ccam
description: 在 Linux 或 macOS 上用一個指令、.deb 套件或從原始碼安裝 ccam。需要先安裝 Claude Code。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

在 Linux 或 macOS 的終端機中執行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

這個腳本會下載適合你系統的發行版本，檢查其 SHA-256 總和檢查碼，並把 `ccam` 放入 `~/.local/bin`；之後 `ccam add` 會在同一個資料夾建立 `claude1`、`claude2`… 快捷指令。如果該資料夾不在你的 `PATH` 中，腳本會印出要加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

ccam 會替你啟動 Claude Code，所以請先安裝 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)，並確認 `claude --version` 可以正常執行。

## 支援的系統

- x86-64 和 ARM64 上的 Linux：Ubuntu、Debian、Fedora、Arch、Alpine 及其他發行版。執行檔為靜態連結，沒有任何相依套件。
- Apple 晶片和 Intel 上的 macOS。
- 不直接支援 Windows。在 WSL 中，ccam 的運作方式與在任何 Linux 系統上相同。

## 其他安裝方式

**Debian 和 Ubuntu 套件。** 從[發行版本頁面](https://github.com/jhyunwoo/ccam/releases)下載 `ccam_<version>_linux_amd64.deb` 並安裝：

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

這個套件會安裝 `/usr/bin/ccam`。

**從原始碼安裝。** 需要 Go 1.26 或更新版本：

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**指定版本。** 設定 `CCAM_VERSION`：

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Shell 自動補齊

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

自動補齊涵蓋 `ccam` 的所有指令和選項，也包括帳號名稱。

## 檢查安裝

```sh
ccam version
ccam doctor
```

`ccam doctor` 會檢查是否找得到 Claude Code，以及每個帳號資料夾、連結和登入是否都正常。接著請閱讀[快速入門](/docs/quick-start)。

## 解除安裝

移除 ccam 指令和帳號快捷指令；如果不再需要這些帳號，再移除資料資料夾：

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

刪除 `~/.local/share/ccam` 也會一併刪除存放在其中的登入資訊和對話記錄。你[匯入](/docs/commands)的 `~/.claude` 資料夾則會留在原處。
