---
title: 快捷指令
description: 用 ccam shortcut 為帳號建立 ccwork 這類專屬指令，並用 ccam shims install 重新建立 ~/.local/bin 中的連結。
updated: 2026-10-23
order: 3
translation: machine
sourceHash: "38e2d1aba397"
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
ccam shims install            # recreate missing or outdated links of all accounts
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
- 不能是 `~/.local/bin` 中不屬於 ccam 連結的檔案（或這個帳號自己由 ccam 建立的快捷指令）；
- 不能和你 `PATH` 上的其他執行檔同名。

載入 `config.toml` 時也會以同樣的規則檢查：手動編輯的快捷指令若不符合格式、是 `claude` 或 `ccam`，或同時指定給兩個帳號，在修正之前每個指令都會失敗，並顯示指出該帳號的訊息。

同時指定指令和 `--remove` 會以 `give a command or --remove, not both` 停止。設定帳號已有的名稱時，會印出 `work: ccwork (unchanged)`，並在連結不存在時重新建立。

## 運作方式

設定名稱時會執行 `checkShortcut`，它以 bin 資料夾和 ccam 執行檔的實際路徑呼叫 `shim.Check`：

<!-- motion:shortcuts-check -->

名稱會依序通過各項檢查。每項檢查旁邊都有一個會被它拒絕的名稱：`Work` 不符合格式，`run` 是 ccam 指令，`ccside` 屬於另一個帳號，`ccfile` 是 `~/.local/bin` 中已存在的檔案，`git` 則會遮蔽 `/usr/bin/git`。只有通過所有檢查的名稱才會被儲存並建立連結。

1. `shim.Check` 先把名稱和 `^[a-z0-9][a-z0-9._-]{0,63}$` 比對，再比對 `config.Reserved`，然後是 `cfg.FindShortcut`。它對 `~/.local/bin/<name>` 呼叫 `os.Lstat`，並對 `PATH` 上的其他每個資料夾呼叫 `os.Stat`，拒絕不是指向 ccam 連結的執行檔。
2. ccam 以 `fsutil.AtomicWrite`（暫存檔、`fsync`、rename）儲存 `config.toml`。
3. `shim.Remove` 刪除舊快捷指令的連結，但只在它是指向 ccam 的連結時才刪除。
4. `shim.Install` 依名稱順序走訪每個已設定的快捷指令，為缺少連結的快捷指令建立連結，以 `os.Symlink` 指向 ccam 執行檔；執行檔路徑以 `os.Executable` 和 `filepath.EvalSymlinks` 解析。某個快捷指令的位置被檔案佔用時，不會妨礙其他快捷指令；錯誤訊息會逐一列出。

ccam 會把它建立的每個連結，以及該連結所執行的執行檔，記錄在 `~/.local/state/ccam/shortcuts.json` 中。當 `os.Readlink` 傳回 ccam 的路徑，或連結和 ccam 路徑解析到同一個檔案時，這個連結就算是 ccam 自己的。已記錄的連結，若其執行檔之後被移動或取代（ccam 被移到其他資料夾、重新安裝），就算是過時的，`shim.Install` 會把它指向目前這個 ccam。指向另一個仍已安裝之 ccam 的已記錄連結會維持原狀。每次移除和每次覆寫檢查都會經過這些測試，所以 ccam 不會碰它沒有建立的檔案，例如你自己的 ccam 副本或你自己建立的連結。

執行快捷指令時，ccam 唯一需要的就是你輸入的名稱：

<!-- motion:shortcuts-link -->

shell 在 `PATH` 上找到 `ccwork`，核心跟隨連結啟動 ccam 執行檔，而 `argv[0]` 仍然是 `ccwork`。在 `app.Main` 中，`ccam` 以外的任何名稱都會載入一次 `config.toml` 並呼叫 `shim.Lookup`，找出 `shortcut` 等於 `filepath.Base(argv[0])` 的帳號。接著 `launcher.LaunchWith` 會在設定已載入的情況下執行和 [`ccam <name>`](/features/launch) 相同的步驟，並把 `argv[1:]` 原封不動地傳給 `claude`。

為什麼用符號連結，而不是 shell 別名或包裝腳本：符號連結在每種 shell 和腳本中都能運作，不會多啟動一次 shell（大約 1 到 2 ms），而且啟動時仍然只讀取一個檔案。`scripts/perf-check.sh` 會用 `hyperfine` 測量快捷指令，標準和 `ccam <name>` 一樣是 10 ms（p95）。

### 在 Windows 上

快捷指令是 `%USERPROFILE%\.local\bin` 中的 `ccwork.exe`，`shim.Name` 會在查詢前去掉 `.exe`（不分大小寫），所以 `ccwork`、`ccwork.exe` 和 `CCWORK.EXE` 都會啟動 `work`。ccam 會先建立符號連結，Windows 只在開發人員模式下或對系統管理員允許這麼做。沒有這項權限時，它會建立永久連結：同一個 `ccam.exe` 檔案的第二個名稱，不需要任何權限。永久連結和 `ccam.exe` 是同一個檔案時，就會被視為 ccam 自己的連結。永久連結會一直執行它建立時所指的檔案，所以當 `ccam update` 以新檔案取代 `ccam.exe` 時，它會用 `shim.Relink` 更新每個以永久連結建立的快捷指令，並印出 `✓ shortcut ccwork renewed`。再次執行 `install.ps1` 也會透過 `ccam shims install` 做同樣的事（`OK 2 shortcuts now run the new ccam`），而對於任何屬於較舊 ccam 副本的快捷指令，`ccam shims install` 和 `ccam doctor --fix` 也會這樣做。名稱檢查會在其他 `PATH` 資料夾中，以 `PATHEXT` 中的每個副檔名（`.exe`、`.cmd` 等）尋找 `<name>`。

請在 shell（PowerShell、cmd.exe、Git Bash）中輸入快捷指令。從檔案總管、「執行」對話方塊或 `Start-Process` 開啟的符號連結會以單純的 `ccam` 啟動，開啟的是儀表板，而不是該帳號。

## 安全性和特殊情況

- **沒有帳號認領的連結。** 如果 `~/.local/bin` 中某個指向 ccam 的連結名稱沒有任何帳號使用（例如你手動編輯了 `config.toml` 中的 `shortcut` 之後），ccam 不會改用預設帳號，而是以 `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` 停止。請用 `ccam shortcut` 重新設定名稱，或刪除該連結。
- **已存在的檔案。** ccam 絕不取代它沒有建立的檔案。`ccam shims install` 會回報 `<path> already exists and is not a ccam link, so it was left alone`，並仍會建立其他連結。接著 `ccam add` 和 `ccam import` 會印出 `shortcut: ccwork (not created; see above)`。
- **移動過的 ccam 執行檔。** 這時 ccam 建立的連結會指向已不存在的路徑。`ccam shims install`（或 `ccam doctor --fix`）會把它們指向你用來執行它的 ccam。
- **重新命名和移除。** `ccam rename` 會保留快捷指令的名稱，連結也繼續有效，因為它指向 ccam，而 ccam 會在 `config.toml` 中查詢名稱。`ccam remove` 會刪除帳號的連結。`ccam shortcut <account> <new>` 會先刪除舊連結再建立新連結。
- **PATH。** 當 ccam 建立連結而 `~/.local/bin` 不在你的 `PATH` 上時，它會印出 `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`，後面接著要加入的 `export PATH=...` 那一行。在 Windows 上，它會改為印出把該資料夾加到使用者 `PATH` 的 PowerShell 指令。
- **事後檢查。** `ccam doctor` 會回報 `missing shortcuts in /home/me/.local/bin: ccwork` 和執行較舊 ccam 的快捷指令，並在該資料夾不在 `PATH` 上時發出警告。`ccam doctor --fix` 會建立缺少的連結並更新過時的連結，而且只對它實際建立的連結顯示 `→ fixed`。執行另一個仍已安裝之 ccam 的快捷指令可以運作，只屬於資訊；`ccam doctor --fix` 會把它指向你用來執行它的 ccam，`ccam add` 和 `ccam import` 則會把它顯示為 ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [管理帳號](/features/accounts)
- [指令參考](/docs/commands#shortcuts)
- [安裝 ccam](/docs/install)
