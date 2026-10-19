---
title: 快捷指令
description: 用 ccam shortcut 為帳號建立 ccwork 這類專屬指令，並用 ccam shims install 重新建立 ~/.local/bin 中的連結。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

快捷指令是你自己命名、用來以某個帳號啟動 Claude Code 的指令：`ccwork --resume` 的效果和 `ccam work --resume` 相同。它是 `~/.local/bin` 中指向 ccam 的連結，所以在 bash、zsh、腳本以及任何從 `PATH` 執行程式的地方都能使用。快捷指令是選用的；請為你常啟動的帳號建立一個。

## 用法

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam 會確認每項變更：

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

一個帳號最多只有一個快捷指令，以 `shortcut` 儲存在 `~/.config/ccam/config.toml` 中。`ccam add` 設定精靈會詢問 `Create a shortcut command for <name>? [y/N]`，接著詢問 `Shortcut name [cc<name>]:`。在設定精靈之外，`ccam add` 只有加上 `--shortcut` 時才會建立快捷指令。`ccam list` 會顯示每個帳號的快捷指令，沒有時顯示 `-`。

快捷指令會把所有參數傳給 Claude Code，方式和 `ccam <name>` 處理名稱後面的參數相同。這也包括 `--dry-run` 和 `--help`：若要預覽啟動，請使用 `ccam --dry-run work`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--remove` | 刪除帳號的快捷指令及其連結 | 關閉 |
| `-h`、`--help` | 指令的說明 | |
| `--dry-run` | `ccam shortcut` 支援它：會檢查名稱並印出 `Would set the shortcut of work to ccwork.` 或 `Would remove ccwork, the shortcut of work.`，但不儲存也不建立連結 | 關閉 |

`ccam shims install` 不接受參數，並會以 `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` 拒絕 `--dry-run`。

名稱規則，依以下順序檢查：

- 1 到 64 個字元，可使用 `a-z`、`0-9`、`.`、`-` 和 `_`，並以字母或數字開頭；
- 不能是 `ccam`、`claude`，或 `run`、`usage` 這類 ccam 指令；
- 不能是其他帳號的快捷指令；
- 不能是 `~/.local/bin` 中尚未連結到 ccam 的檔案；
- 不能和你 `PATH` 上的其他執行檔同名。

同時指定指令和 `--remove` 會以 `give a command or --remove, not both` 停止。設定帳號已有的名稱時，會印出 `work: ccwork (unchanged)`，並在連結不存在時重新建立。

## 運作方式

設定名稱時會執行 `checkShortcut`，它以 bin 資料夾和 ccam 執行檔的實際路徑呼叫 `shim.Check`：

<!-- motion:shortcuts-check -->

名稱會依序通過各項檢查。每項檢查旁邊都有一個會被它拒絕的名稱：`Work` 不符合格式，`run` 是 ccam 指令，`ccside` 屬於另一個帳號，`ccfile` 是 `~/.local/bin` 中已存在的檔案，`git` 則會遮蔽 `/usr/bin/git`。只有通過所有檢查的名稱才會被儲存並建立連結。

1. `shim.Check` 先把名稱和 `^[a-z0-9][a-z0-9._-]{0,63}$` 比對，再比對 `config.Reserved`，然後是 `cfg.FindShortcut`。它對 `~/.local/bin/<name>` 呼叫 `os.Lstat`，並對 `PATH` 上的其他每個資料夾呼叫 `os.Stat`，拒絕不是指向 ccam 連結的執行檔。
2. ccam 以 `fsutil.AtomicWrite`（暫存檔、`fsync`、rename）儲存 `config.toml`。
3. `shim.Remove` 刪除舊快捷指令的連結，但只在它是指向 ccam 的連結時才刪除。
4. `shim.Install` 為每個已設定但缺少連結的快捷指令建立連結，以 `os.Symlink` 指向 ccam 執行檔；執行檔路徑以 `os.Executable` 和 `filepath.EvalSymlinks` 解析。它絕不移除任何東西。

當 `os.Readlink` 傳回 ccam 的路徑，或連結和 ccam 路徑解析到同一個檔案時，這個連結就算是 ccam 自己的。每次移除和每次覆寫檢查都會經過這項測試，所以 ccam 不會碰它沒有建立的檔案。

執行快捷指令時，ccam 唯一需要的就是你輸入的名稱：

<!-- motion:shortcuts-link -->

shell 在 `PATH` 上找到 `ccwork`，核心跟隨連結啟動 ccam 執行檔，而 `argv[0]` 仍然是 `ccwork`。在 `app.Main` 中，`ccam` 以外的任何名稱都會載入一次 `config.toml` 並呼叫 `shim.Lookup`，找出 `shortcut` 等於 `filepath.Base(argv[0])` 的帳號。接著 `launcher.LaunchWith` 會在設定已載入的情況下執行和 [`ccam <name>`](/features/launch) 相同的步驟，並把 `argv[1:]` 原封不動地傳給 `claude`。

為什麼用符號連結，而不是 shell 別名或包裝腳本：符號連結在每種 shell 和腳本中都能運作，不會多啟動一次 shell（大約 1 到 2 ms），而且啟動時仍然只讀取一個檔案。`scripts/perf-check.sh` 會用 `hyperfine` 測量快捷指令，標準和 `ccam <name>` 一樣是 10 ms（p95）。

## 安全性和特殊情況

- **沒有帳號認領的連結。** 如果 `~/.local/bin` 中某個指向 ccam 的連結名稱沒有任何帳號使用（例如你手動編輯了 `config.toml` 中的 `shortcut` 之後），ccam 不會改用預設帳號，而是以 `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` 停止。請用 `ccam shortcut` 重新設定名稱，或刪除該連結。
- **已存在的檔案。** ccam 絕不取代它沒有建立的檔案。`ccam shims install` 會以 `<path> already exists and is not a ccam link, so it was left alone` 停止。
- **移動過的 ccam 執行檔。** 這時舊的連結會指向已不存在的路徑，ccam 不會把它們視為自己的連結。請刪除它們，然後執行 `ccam shims install`。
- **重新命名和移除。** `ccam rename` 會保留快捷指令的名稱，連結也繼續有效，因為它指向 ccam，而 ccam 會在 `config.toml` 中查詢名稱。`ccam remove` 會刪除帳號的連結。`ccam shortcut <account> <new>` 會先刪除舊連結再建立新連結。
- **PATH。** 當 ccam 建立連結而 `~/.local/bin` 不在你的 `PATH` 上時，它會印出 `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`，後面接著要加入的 `export PATH=...` 那一行。
- **事後檢查。** `ccam doctor` 會回報 `missing shortcuts in /home/me/.local/bin: ccwork`，並在該資料夾不在 `PATH` 上時發出警告。`ccam doctor --fix` 會建立缺少的連結。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [管理帳號](/features/accounts)
- [指令參考](/docs/commands#shortcuts)
- [安裝 ccam](/docs/install)
