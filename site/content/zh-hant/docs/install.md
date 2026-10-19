---
title: 安裝 ccam
metaTitle: "在 Linux 和 macOS 上安裝 ccam"
description: 在 Linux 或 macOS 上用一個指令、.deb 套件或從原始碼安裝 ccam，並用 ccam update 保持最新版本。需要先安裝 Claude Code。
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
---

在 Linux 或 macOS 的終端機中執行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

這個腳本會下載適合你系統的發行版本，檢查其 SHA-256 總和檢查碼，並把 `ccam` 放入 `~/.local/bin`；之後你選擇的快捷指令（例如 `ccwork`）也會由 ccam 建立在這個資料夾中。如果該資料夾不在你的 `PATH` 中，腳本會印出要加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

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

## 更新

```sh
ccam update --check
ccam update
```

`ccam update --check` 只會告訴你是否有更新的發行版本。`ccam update` 會下載適合你系統的發行版本，以該版本的 `checksums.txt` 檢查其 SHA-256 總和檢查碼，再一次取代 `ccam` 檔案。如果總和檢查碼不符，就不會變更任何東西。`ccwork` 這類快捷指令是指向 ccam 的連結，因此會立即執行新版本。

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam 從不自動檢查更新，只有在你執行這個指令時才會尋找新的發行版本。若要安裝特定的發行版本（包括較舊的版本），請加上 `--version`。`ccam --dry-run update` 會顯示將下載和取代哪些內容，但不會做任何變更。

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

重新執行安裝腳本也會更新 `~/.local/bin` 中的 ccam。如果你是用 `go install` 建置 ccam，再執行一次同樣的指令即可。

**Debian 和 Ubuntu 套件。** `ccam update` 不會取代由 apt 安裝的 `/usr/bin/ccam`，而是改為印出以下步驟。從[發行版本頁面](https://github.com/jhyunwoo/ccam/releases)下載新的 `ccam_<version>_linux_amd64.deb`，覆蓋安裝舊版：

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## 解除安裝

請趁 ccam 還在時，先移除帳號的快捷指令。`ccam list` 會在 Shortcut 欄顯示它們；用 `ccam shortcut <account> --remove` 逐一移除，或自行從 `~/.local/bin` 刪除這些連結。接著移除 ccam 指令；如果不再需要這些帳號，再移除資料資料夾：

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

刪除 `~/.local/share/ccam` 也會一併刪除存放在其中的登入資訊和對話記錄。你[匯入](/docs/commands)的 `~/.claude` 資料夾則會留在原處。
