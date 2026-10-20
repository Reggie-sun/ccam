---
title: 管理帳號
description: 用 ccam add、ccam login、ccam list、ccam rename、ccam remove 和 ccam import 新增、登入、列出、重新命名、移除和匯入 Claude Code 帳號。
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 9eb992214132
---

每個 ccam 帳號都是一個獨立的 Claude Code 設定資料夾，有自己的登入資訊。這些指令用來建立、登入、列出、重新命名和移除帳號，也能登錄你已經有的資料夾。

## 用法

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

凡是接受 `<name|number>` 的指令，都可以指定帳號名稱或它的槽位編號，也就是 `ccam list` 的 `#` 欄。

在終端機中不加任何選項執行 `ccam add`，會開啟設定精靈。它會先詢問名稱，接著顯示共用項目的核取清單（除了對話記錄以外都已勾選），然後詢問 `Create a shortcut command for work? [y/N]`，回答 `y` 後再詢問 `Shortcut name [ccwork]:`。只要加上 `--share`、`--shortcut` 或 `--no-login` 任一個選項，或是輸入來自管線而不是終端機，就會略過設定精靈；這時必須指定名稱，而且除非 `--shortcut` 指定了名稱，否則不會建立快捷指令。

`ccam list` 會顯示槽位、名稱、快捷指令、電子郵件、共用的項目和上次使用時間，預設帳號前面會標示 `*`：

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

加上 `--json` 時，每個帳號都是一個物件，包含 `slot`、`shortcut`、`name`、`email`、`default`、`dir`、`last_used`（帳號從未啟動過時省略）和 `share`。`share` 會把每個項目（`mcp`、`settings`、`memory`、`extensions`、`history`）對應到它的群組名稱或 `"local"`。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `ccam add --no-login` | 建立帳號但不登入；之後再用 `ccam login <name>` 登入 | 關閉 |
| `ccam add --share <item=group>` | 與某個群組共用一個項目，或用 `item=local` 讓它保持私有。可重複指定 | `mcp`、`settings`、`memory`、`extensions` 在 `default`；`history` 為 local |
| `ccam add --shortcut <command>` | 同時在 `~/.local/bin` 中建立這個快捷指令 | 無 |
| `ccam add --api --key-from <source>` | 使用 API 金鑰而不是登入的帳號；請參閱 [API 金鑰帳號](/features/api) | 關閉 |
| `ccam list --json` | 以 JSON 輸出帳號 | 關閉 |
| `ccam remove -y`、`--yes` | 不詢問直接移除 | 關閉 |
| `ccam import --name <name>` | 你指定的資料夾要使用的帳號名稱 | 從資料夾名稱推導 |
| `ccam import -y`、`--yes` | 不詢問直接匯入找到的每個資料夾 | 關閉 |
| `--dry-run` | `ccam list` 和 `ccam completion` 不會變更任何東西，所以接受它。`add`、`login`、`logout`、`use`、`rename`、`remove` 和 `import` 會拒絕它，並顯示 `` --dry-run is not supported by `ccam add`; nothing was changed ``（其中為該指令的名稱） | 關閉 |

`ccam login`、`ccam logout`、`ccam use` 和 `ccam rename` 沒有自己的選項。`ccam completion` 接受 `bash`、`zsh` 或 `powershell`；其他值會失敗並顯示 `only bash, zsh and powershell are supported`。

帳號名稱為 1 到 32 個字元，可使用 `a-z`、`0-9`、`-` 和 `_`，並以字母或數字開頭。名稱不能全是數字，因為 `ccam 2` 代表槽位 2。名稱也不能是 ccam 的指令（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`folder`、`pick`、`ext`），所以 `ccam <name>` 永遠不會有歧義。這份清單是 `internal/config` 中的 `config.Reserved`。

## 運作方式

`ccam add work --shortcut ccwork` 會在 `addProfile`（`internal/cli/profiles.go`）中執行以下步驟：

<!-- motion:accounts-add -->

1. 在建立任何東西之前，`config.ValidateName` 會檢查名稱，`shim.Check` 會檢查快捷指令。
2. `profile.Add` 以權限 `0700` 建立 `~/.local/share/ccam/profiles/work/`，並給帳號最小的可用槽位編號（`Config.NextSlot`），因此已移除帳號的編號會被重複使用。第一個帳號會成為預設帳號。
3. `applyShare` 先用 `share.Build` 再用 `Apply`，把新建的空資料夾連結到它的共用群組，例如把 `settings.json` 連結到 `groups/default/settings.json`。如果這一步失敗，新資料夾會被移除，因為殘留的資料夾會讓之後每次重試都因 `folder already exists` 而失敗。
4. 帳號會透過 `fsutil.AtomicWrite` 儲存到 `~/.config/ccam/config.toml`：先寫入暫存檔、`fsync`，再 `rename`，所以當機時絕不會留下只寫了一半的設定。
5. `runLogin` 會在 `CLAUDE_CONFIG_DIR` 設為帳號資料夾的狀態下啟動 `claude auth login`。如果登入失敗，帳號仍會保留，ccam 會印出錯誤，並告訴你之後執行 `ccam login work`。
6. `profile.ReadEmail` 從 `profiles/work/.claude.json` 取出 `oauthAccount.emailAddress`。這個檔案可能有好幾 MB，所以 `TopLevelValue` 只掃描那一個鍵，而不解析整個檔案。
7. `shim.Install` 把 `~/.local/bin/ccwork` 連結到 ccam 執行檔；ccam 啟動時會在 `config.toml` 中查詢 `argv[0]`。

登入可能要花上幾分鐘，這段期間另一個 ccam 可能會儲存設定。因此需要等待的指令會透過 `env.update` 儲存：它會重新載入 `config.toml`，套用那一項變更（這裡是電子郵件），再儲存。基於同樣的理由，`ccam remove` 會在詢問之後重新載入。

`ccam login` 會對現有的帳號執行第 5 和第 6 步。`ccam logout` 會在帳號中執行 `claude auth logout`，並清除儲存的電子郵件。`ccam use` 只會改寫 `config.toml` 中的 `default`。

### 在遠端機器上登入

透過 SSH 連線時，`claude auth login` 在瀏覽器端的步驟可能會失敗。Claude Code 在伺服器上的 `http://localhost:<port>/callback` 開啟回呼伺服器等待登入，但完成登入的瀏覽器在你的電腦上執行，那個位址在你的電腦上連不到任何地方，瀏覽器會顯示 `ERR_CONNECTION_REFUSED` 之類的錯誤。這個位址仍然帶著一次性代碼，所以 ccam 可以替你轉交：

<!-- motion:accounts-relay -->

`runLogin` 不會把你的終端機直接交給 Claude Code。它把 Claude Code 的標準輸入接到一個管線，由 `relayLogin` 把你輸入的每一行複製進去，所以在 Claude Code 的 `Paste code here` 提示貼上的代碼仍然會送達。當某一行是回呼位址時，`callbackURL` 只在以下條件都成立時才接受它：主機是 `localhost`、`127.0.0.1` 或 `[::1]`，有指定連接埠，路徑是 `/callback`，而且 `code` 和 `state` 都存在；`http://` 可以省略。接著 `deliverCallback` 會像瀏覽器原本會做的那樣，在伺服器上請求該位址，時間限制為 10 秒，而且不跟隨重新導向。這個請求會抵達正在等待的 `claude auth login` 的回呼伺服器（圖中的虛線）。Claude Code 收到代碼並儲存登入資訊，在 Linux 上是存到帳號資料夾中的 `.credentials.json`。貼上的內容絕不會讓 ccam 連線到其他主機。

ccam 在開始登入前會印出：

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

成功時會印出 `→ Passed the sign-in to Claude Code on this machine.`

登入在終端機上等待輸入時，ccam 會關閉終端機的回顯和行模式，你每輸入一個字元就顯示一個 `*`，所以貼上的代碼或回呼位址不會出現在畫面上。行編輯由 `readLine` 自己處理：`Enter` 送出這一行，`Backspace` 刪除一個字元，在空行上按 `Ctrl+D` 結束輸入。登入結束時會還原終端機設定；收到 `Ctrl+C`、`SIGTERM` 或 `SIGHUP` 時也會還原，它們仍會結束 ccam。如果標準輸入和標準錯誤不是終端機，各行會照原樣複製。

### 重新命名、移除和匯入

`ccam rename` 和 `ccam remove` 會處理帳號的資料夾，以及 ccam 在其他地方以該帳號名稱保存的內容。`nameState` 列出這些路徑：`~/.local/state/ccam/last-used/<name>`、`~/.local/state/ccam/mcp/<name>`，以及用量快取 `~/.local/state/ccam/usage/<name>.json` 和 `<name>.statusline.json`。`ccam remove work` 依以下順序處理它們：

<!-- motion:accounts-remove -->

如果帳號正在執行，`ccam remove` 會先拒絕，否則會詢問你。如果帳號共用擴充功能，`RewritePluginPaths` 會改寫群組的 `installed_plugins.json` 和 `known_marketplaces.json`，因為 Claude Code 儲存的是外掛的絕對路徑，而這些路徑可能經過這個帳號的資料夾。`profile.Trash` 會把資料夾移到 `~/.local/share/ccam/trash/<name>-<time>/`，帳號會從 `config.toml` 中移除，`nameState` 的路徑會被刪除，讓之後同名的帳號不會繼承它們，快捷指令的連結也會被移除。若要復原，把資料夾從垃圾桶移回原處即可。

`ccam rename` 會重新命名資料夾，把其 `plugins/` 中的外掛路徑從舊資料夾改寫為新資料夾，保留預設帳號的設定，並把 `nameState` 的路徑移到新名稱下。快捷指令的名稱不變。如果改寫外掛路徑失敗，因為設定尚未儲存，資料夾會被改回原名。

`ccam import` 會在原位置登錄資料夾，記錄為 `config.toml` 中的 `dir`。它絕不移動或複製資料夾：複製後的更新權杖會讓各個副本互相登出。不指定資料夾時，它會尋找尚未登錄的 `~/.claude` 和 `~/.claude-*` 資料夾，並逐一詢問；名稱取自資料夾（`~/.claude` 為 `main`，`~/.claude-work` 為 `work`，`~/.claude-2` 為 `acct2`），無法推導出名稱的資料夾會列出應使用的 `--name` 指令。沒有 `CLAUDE_CONFIG_DIR` 時，Claude Code 會把 `.claude.json` 放在家目錄中 `~/.claude` 的旁邊，所以當你匯入 `~/.claude` 而它沒有自己的 `.claude.json` 時，ccam 會建立指向 `~/.claude.json` 的連結 `~/.claude/.claude.json`。這樣直接執行 `claude` 和 `ccam main` 會看到同一個帳號。匯入的帳號在你執行 `ccam share` 之前不共用任何項目，也不會建立快捷指令。

`ccam completion bash|zsh` 會印出一份自動完成腳本，也能自動完成帳號名稱。

## 安全性和特殊情況

- 帳號執行時，`ccam rename` 和 `ccam remove` 會拒絕（`refuseRunning`），並顯示 `work is running; quit its Claude Code session first`，因為工作階段會持續寫入該資料夾。
- 對匯入的帳號執行 `ccam remove` 只會取消登錄：`✓ main unregistered (folder ... left in place)`。你的 `~/.claude` 不會被移到垃圾桶。
- 如果 `profiles/` 中已有同名資料夾，`ccam add` 會拒絕該名稱，並提示你使用 `ccam import`。
- `ccam import` 會拒絕已登錄的資料夾、不是資料夾的路徑，以及已被使用的名稱。
- `ccam logout` 和 `ccam login` 不會檢查帳號是否正在執行。
- 如果轉交無法送達 Claude Code（例如登入已經結束），會印出 `ccam: could not pass the sign-in to Claude Code:` 和原因。HTTP 錯誤通常表示代碼已過期；請重新開始登入。

## 相關內容

- [啟動 Claude Code](/features/launch)
- [快捷指令](/features/shortcuts)
- [共用 MCP 伺服器、設定和對話記錄](/docs/sharing)
- [多個 Claude Code 帳號](/guides/multiple-claude-code-accounts)
