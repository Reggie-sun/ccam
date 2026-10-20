---
title: 更新 ccam
description: ccam update 會在檢查 SHA-256 之後以最新的 GitHub 發行版本取代 ccam；--check 只回報，--version 安裝指定的發行版本標籤。
updated: 2026-10-20
order: 10
translation: machine
sourceHash: d91499b95d3d
---

`ccam update` 會從 GitHub 下載 ccam 的發行版本，以該發行版本的 `checksums.txt` 驗證，然後透過一次 rename 取代正在執行的執行檔。如果你是用安裝腳本或壓縮檔安裝 ccam，請使用它。ccam 絕不會自行檢查更新：只有這個指令會連線到 GitHub，所以其他所有指令在離線時都一樣能運作。

## 用法

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

沒有更新的發行版本時，ccam 會印出 `ccam 0.4.0 is up to date.`，且不變更任何東西。以 `--version` 指定你目前正在執行的版本時，會印出 `ccam 0.4.0 is already installed.`

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--check` | 只回報是否有更新的發行版本 | 關閉 |
| `--version <tag>` | 安裝這個發行版本標籤，例如 `v0.3.0`，可以比目前的版本新或舊 | 最新版本 |
| `--dry-run` | 支援：查詢最新的標籤（除非 `--version` 已指定），對執行檔執行相同的檢查，並印出下載、總和檢查碼和取代的路徑，但不下載 | 關閉 |

`--check` 和 `--version` 不能一起使用。標籤必須是 `vMAJOR.MINOR.PATCH` 的形式，可以帶有預發行版本或組建後綴；開頭的 `v` 可以省略。其他寫法都會失敗，例如 `--version takes a release tag such as v0.3.0, not "1.2"`。

環境變數 `CCAM_RELEASES_URL` 可以取代 `https://github.com/jhyunwoo/ccam/releases`，用於鏡像站或測試。它的結構必須相同：`<url>/latest` 會重新導向到 `<url>/tag/<tag>`，檔案則位於 `<url>/download/<tag>/<asset>`。

## 運作方式

程式碼位於 `internal/update`，只有 `ccam update` 會匯入它。啟動流程絕不會呼叫它。

1. `cli.updateCmd` 以 `os.Executable` 找出正在執行的執行檔，並以 `filepath.EvalSymlinks` 解析符號連結。被取代的就是這個實際路徑。
2. 沒有 `--version` 時，`Updater.Latest` 會請求 `<releases>/latest` 但不跟隨重新導向，並從 `Location` 標頭（`…/tag/v0.5.0`）讀取標籤。這裡使用的是發行版本頁面而不是 GitHub API，所以不受 API 每小時 60 次請求的限制。查詢的時間限制為 10 秒。
3. `parseVersion` 以嚴格的格式檢查標籤。標籤會出現在下載 URL 和印出的 shell 指令中，所以任何含有 `/`、`%`、空白或 shell 字元的內容都無法通過。`compareVersions` 會把它和正在執行的版本比較；預發行版本排在其正式版本之前。
4. `managed()` 和 `writable()` 決定 ccam 是否能取代這個執行檔（見下文）。
5. `Updater.Fetch` 同時下載 `ccam_<os>_<arch>.tar.gz`（在 Windows 上為 `ccam_windows_<arch>.zip`）和 `checksums.txt`，兩者合計的時間限制為 2 分鐘。每個下載的上限為 256 MB。
6. 在做任何其他事之前，壓縮檔的 SHA-256 必須和 `checksums.txt` 中對應的那一行相符。缺少那一行或總和檢查碼不同時，更新會以 `nothing was changed` 停止。
7. `extract()` 讀取 tar，取出第一個名為 `ccam` 的一般檔案（`extractZip` 則從 zip 中取出 `ccam.exe`）。資料夾、連結和其他項目都會被忽略，執行檔的上限同樣是 256 MB。
8. `Updater.Install` 以 `fsutil.AtomicWrite` 寫入新的執行檔：在執行檔旁邊建立暫存檔 `.ccam.tmp*`，套用舊檔案的權限位元，寫入磁碟後再 rename 覆蓋舊檔案。

從查詢標籤到 rename 的完整流程：

<!-- motion:update-pipeline -->

壓縮檔和 `checksums.txt` 會並行下載，在總和檢查碼相符之前不會寫入任何東西。rename 是原子操作，所以當機或磁碟已滿時，留下的不是舊執行檔就是新執行檔，絕不會只有一半。`~/.local/bin/ccwork` 這類快捷指令是指向 ccam 執行檔的符號連結，所以不需任何變更就會執行新版本。

在 Windows 上，執行中程式的檔案無法被取代或刪除，只能重新命名。`Updater.Install` 會把新的執行檔寫入 `ccam.exe.new`，把執行中的 `ccam.exe` 重新命名為 `ccam.exe.old`，再把新檔案重新命名為 `ccam.exe`。如果最後一步失敗，舊的執行檔會被改回原來的名稱。下次更新時會刪除 `ccam.exe.old`。以永久連結而非符號連結建立的快捷指令（[快捷指令](/features/shortcuts)）會在之後立即更新，所以它們也會執行新版本。

### 由套件管理的安裝

從 `.deb` 安裝的 ccam 屬於 apt，如果繞過 apt 取代它，套件資料庫的內容就會出錯。在下載之前，ccam 會檢查執行檔歸誰管理：

<!-- motion:update-managed -->

1. 位於 `/usr/bin`、`/bin`、`/usr/sbin` 或 `/sbin` 的執行檔視為由套件管理。
2. 在其他位置時，如果已安裝 `dpkg` 且 `dpkg -S <path>` 成功，同樣視為由套件管理。
3. 否則 ccam 會在執行檔旁邊建立並移除一個暫存檔 `.ccam-update-*`，確認自己能在那裡寫入。

對於由套件管理的執行檔，ccam 會停止並印出 apt 的步驟：

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

在 Linux 以外的系統上，只會印出安裝腳本那一行：在 Windows 上是 `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`。資料夾無法寫入時，ccam 會建議重新執行安裝腳本，把 ccam 安裝到 `~/.local/bin`。

## 安全性和特殊情況

- 開發版組建的版本會顯示為 `dev`，沒有可比較的發行版本，所以 `ccam update` 和 `ccam update --check` 會拒絕執行。`ccam update --version vX.Y.Z` 可以在它上面安裝發行版本。
- 下載失敗、找不到你平台的檔案，或總和檢查碼不符時，已安裝的執行檔都會維持原樣。
- 不存在的發行版本會以 `there is no release <tag>` 失敗。當 GitHub 以 HTTP 403 或 429 拒絕時，ccam 會說明這可能是頻率限制，並請你稍後再試。
- `--check` 只會讀取。它不會檢查執行檔是否由套件管理或是否可寫入。
- 以 `--version` 指定較舊的標籤就是回復舊版的方法。

## 相關內容

- [安裝](/docs/install)
- [指令參考](/docs/commands)
- [常見問題](/docs/faq)
