---
title: 依資料夾使用帳號
description: ccam folder 讓你在未指定帳號名稱就啟動 Claude Code 時，某個資料夾及其下所有資料夾都使用同一個帳號，例如公司程式碼使用 work。
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` 會讓在 `~/code/company` 及其下任何資料夾中未指定帳號名稱的每次啟動，都以 `work` 啟動 Claude Code。在該資料夾之外，則和以前一樣啟動預設帳號。指定帳號名稱時，例如 `ccam personal` 或 `ccpersonal` 這類快捷指令，一律會啟動該帳號。

## 用法

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

以下啟動會使用該資料夾的帳號：

- `ccam -p "…"`，以及其他第一個參數是 Claude Code 選項的啟動
- 不帶名稱的 `ccam run`，以及 `ccam --dry-run`
- `ccam --fallback -p "…"` 的第一次嘗試
- 選擇器，也就是單獨的 `ccam`，開啟時游標會停在該資料夾的帳號上

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `<account>` | 要使用的帳號名稱或編號 | |
| `[folder]` | 資料夾，可以是相對路徑或絕對路徑 | 目前的資料夾 |
| `--remove` | 移除完全相同的那個資料夾的項目 | 關閉 |
| `--dry-run` | 會被接受：印出變更但不儲存 | 關閉 |

## 運作方式

每個資料夾都是 `~/.config/ccam/config.toml` 中的一個 `[[folder]]` 項目，而這個檔案本來就是每次啟動都會讀取的：

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` 會把資料夾解析為已解開符號連結的絕對路徑並儲存。再次設定同一個資料夾時，會取代它的帳號。
2. 未指定名稱的啟動會呼叫 `launcher.HereProfile`。沒有任何項目時，它會立即採用預設帳號。
3. 否則，它會以字串比對工作目錄和每個 `path`：只比對資料夾本身或其下的資料夾，絕不會比對到只是開頭字母相同的資料夾（`/code/company` 不涵蓋 `/code/companyx`）。最深的項目優先，所以 `~/code` 可以是 `personal`，而 `~/code/company` 是 `work`。
4. 只有在工作目錄原樣比對不到任何項目時，ccam 才會解開它的符號連結，再比對一次。
5. 在 Windows 上，比對時不分大小寫，和 Windows 本身一樣。

除了 `config.toml` 之外不讀取任何檔案，所以在資料夾中啟動不會多花可測量的時間（不論有沒有項目都是 4.0 ms，`hyperfine`，300 次）。

## 安全性和特殊情況

- `ccam rename` 會把資料夾項目轉移到新名稱，`ccam remove` 會刪除被移除帳號的項目。
- 如果項目的帳號已不存在，在該資料夾中未指定名稱的啟動會停止：`<folder> is set to account "<name>", which does not exist. Change it with ccam folder`。`ccam doctor` 會回報這類項目，以及資料夾已不存在的項目。
- 在從上層資料夾取得帳號的資料夾中執行 `ccam folder --remove` 時，不會移除任何東西，而是告訴你是哪個上層資料夾。
- `ccam next` 依輪替順序或用量選擇帳號，不會參考資料夾。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [快捷指令](/features/shortcuts)
- [管理帳號](/features/accounts)
- [指令參考](/docs/commands)
