---
title: 安裝 ccam
metaTitle: "在 Linux、macOS 和 Windows 上安裝 ccam"
description: 在 Linux、macOS 或 Windows 上用一個指令、.deb 套件或從原始碼安裝 ccam，並用 ccam update 保持最新版本。需要先安裝 Claude Code。
updated: 2026-10-23
order: 1
translation: machine
sourceHash: "af21dcc24169"
---

在 Linux 或 macOS 的終端機中執行：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

這個腳本會下載適合你系統的發行版本，檢查其 SHA-256 總和檢查碼，並把 `ccam` 放入 `~/.local/bin`；之後你選擇的快捷指令（例如 `ccwork`）也會由 ccam 建立在這個資料夾中。如果該資料夾不在你的 `PATH` 中，腳本會印出要加到 `~/.bashrc` 或 `~/.zshrc` 的那一行。

在 Windows 上，請在 PowerShell 中執行：

```powershell
irm https://ccam.dev/install.ps1 | iex
```

它會下載 `ccam_windows_<arch>.zip`，檢查其 SHA-256 總和檢查碼，把 `ccam.exe` 放入 `%USERPROFILE%\.local\bin`（也就是 Claude Code 安裝程式放置 `claude.exe` 的資料夾），並在該資料夾尚未加入時，把它加到你的使用者 `PATH`，並保留 `%USERPROFILE%\bin` 這類項目原樣。完成後請開啟新的終端機。再次執行它會取代 ccam，即使仍有 Claude Code 工作階段從舊的執行檔執行也一樣（舊檔會被重新命名為 `ccam.exe.old`，若更舊的那個仍在使用中，則為 `ccam.exe.old-<number>`），並更新屬於硬連結的快捷指令：`OK 2 shortcuts now run the new ccam`。

ccam 會替你啟動 Claude Code，所以請先安裝 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)，並確認 `claude --version` 可以正常執行。

## 支援的系統

- x86-64 和 ARM64 上的 Linux：Ubuntu、Debian、Fedora、Arch、Alpine 及其他發行版。執行檔為靜態連結，沒有任何相依套件。
- Apple 晶片和 Intel 上的 macOS。
- x86-64 和 ARM64 上的 Windows 10 和 11，可在 PowerShell、Windows 終端機或 Git Bash 中使用。共用資料夾對每位使用者都可用。共用檔案（`settings.json`、`CLAUDE.md`、`history.jsonl`）需要開發人員模式（設定 > 系統 > 開發人員專用），因為 Windows 只允許這個模式或系統管理員建立指向檔案的連結；請參閱[Windows 上的共用](/features/share#on-windows)。沒有開發人員模式時，第一次 `ccam add` 建立的群組 `default` 只包含 MCP 伺服器和擴充功能，並會說明開啟開發人員模式後如何加入其餘項目。
- 在 WSL 中，ccam 的運作方式與在任何 Linux 系統上相同。

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

從發行標籤（`@latest` 或 `@vX.Y.Z`）建置的版本會回報該版本號，所以 `ccam update` 可以用在它上面；從 `@main` 或超過標籤的本機簽出建置的版本則會回報 `dev`。

**指定版本。** 設定 `CCAM_VERSION`：

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell 自動補齊

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

自動補齊涵蓋 `ccam` 的所有指令和選項，也包括帳號名稱。

## Claude Code 外掛

ccam 外掛會教 Claude Code 本身使用 ccam：檢查哪個帳號還有額度、用 `ccam pick` 在不同帳號上啟動平行的 `claude -p` 工作程序，以及用 `ccam ext` 說明為什麼缺少某個技能或 MCP 伺服器。請在 Claude Code 中安裝：

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

它會新增一個技能、`/ccam:usage`、`/ccam:ext` 和 `/ccam:switch` 這幾個指令，以及兩個 hook：工作階段因用量上限而停止時，你會收到終端機通知；`/ccam:switch [account]` 則會把對話交給另一個帳號（[詳細說明](/features/next#inside-claude-code)）。這些 hook 會執行 `ccam`，所以它必須在 Claude Code 看得到的 `PATH` 上。這個技能絕不會執行登入或移動檔案的指令，而是建議你執行。

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

儀表板（在終端機中執行 `ccam`）也會在每次開啟時尋找新的發行版本，有新版本時會提供更新：按下 `U` 並確認即可。其他指令都不會尋找，ccam 也絕不會未經詢問就更新。若要安裝特定的發行版本（包括較舊的版本），請加上 `--version`。`ccam --dry-run update` 會顯示將下載和取代哪些內容，但不會做任何變更。

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

重新執行安裝腳本也會更新 `~/.local/bin` 中的 ccam。如果你是從發行標籤用 `go install` 建置 ccam，也可以使用 `ccam update`，或再執行一次同樣的指令。

**Windows。** `ccam update` 會下載 Windows 版的 `.zip`。在 Windows 上無法取代執行中的程式，所以舊的執行檔會被重新命名為 `ccam.exe.old`，並在下次更新時刪除；若仍有工作階段從先前的 `ccam.exe.old` 執行，名稱則為 `ccam.exe.old-<time>`。屬於硬連結的快捷指令會在之後立即更新。

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

在 Windows 上，於 PowerShell 中：

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

刪除 `~/.local/share/ccam` 也會一併刪除存放在其中的登入資訊和對話記錄。你[匯入](/docs/commands)的 `~/.claude` 資料夾則會留在原處，但只要它還在某個群組中，它所共用的檔案就位於 `~/.local/share/ccam/groups/`。請先對它執行 `ccam remove <name>`（或 `ccam group leave`）：它會帶著所共用的一切內容的副本離開群組。
