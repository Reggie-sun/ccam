---
title: 快速入門
description: 用 ccam add 新增兩個 Claude Code 帳號，為它們設定 ccwork 這類快捷指令並在帳號之間切換，再選擇要共用的項目。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: d040e34d2a17
---

用 `ccam add` 新增每個帳號，然後用 `ccam <name>` 或你自己命名的快捷指令（例如 `ccwork`）啟動 Claude Code。大多數人只需要這些。

## 1. 新增第一個帳號

```sh
ccam add
```

設定精靈會要求你輸入名稱，例如 `personal` 或 `work`，並詢問這個帳號要和其他帳號共用哪些項目，以及是否建立快捷指令（請見第 4 步）。接著它會執行 `claude auth login`，開啟瀏覽器讓你登入。在遠端機器上，瀏覽器可能會停在無法載入的 `http://localhost:…/callback` 頁面；把這個網址貼到終端機中，ccam 就會完成登入。完成後你會看到：

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. 移轉現有的設定

如果你以前用過 Claude Code，你的設定、`CLAUDE.md`、擴充功能和 MCP 伺服器都在 `~/.claude` 中。把它們複製到共用設定，每個帳號就都能使用：

```sh
ccam migrate
```

ccam 會先顯示要複製的內容並詢問你。它只會讀取 `~/.claude`，絕不複製登入資訊；除非加上 `--history`，否則不會複製對話記錄。詳情請參閱[複製現有的設定](/docs/sharing#copy-an-existing-setup)。

## 3. 新增第二個帳號

```sh
ccam add work
```

在瀏覽器中用另一個帳號登入。

## 4. 切換帳號

`ccam work` 會以 work 帳號啟動 Claude Code。想要更短的指令，可以為帳號取一個專屬的快捷指令，例如 `ccwork`。選完共用項目後，設定精靈會詢問：

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

按下 `y`，然後輸入名稱，或直接按 Enter 使用建議的名稱：

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

登入完成後，ccam 會顯示確認訊息：

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

如果你跳過了這個問題，隨時都可以新增快捷指令：

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

現在，以下任一指令都會以 work 帳號啟動 Claude Code：

```sh
ccwork
ccam work
ccam 2
```

單獨執行 `ccam` 會開啟選擇器：用方向鍵和 Enter 選擇，或直接按下帳號的編號。

快捷指令或帳號名稱後面的所有內容都會原封不動地傳給 Claude Code，所以你可以繼續使用熟悉的選項：

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 選擇要共用的項目

預設情況下，MCP 伺服器、`settings.json`、全域 `CLAUDE.md`，以及 agents、commands、skills 和 plugins 都會在帳號之間共用，對話記錄則不會。如果想在帳號之間接續對話，請為兩個帳號都開啟對話記錄共用：

```sh
ccam share personal --history on
ccam share work --history on
```

ccam 在變更任何內容之前，會先顯示它將移動哪些項目。詳情請參閱[共用](/docs/sharing)。

## 6. 帳號達到限制時

```sh
ccam usage
ccam next --continue
```

`ccam usage` 會顯示每個帳號的 5 小時用量和每週用量。`ccam next --continue` 會在下一個帳號開啟同一段對話。請參閱[用量限制](/docs/usage-limits)。

## 已經在使用 ~/.claude？

直接登錄它，不必重新登入：

```sh
ccam import
```

資料夾會留在原處，直接執行 `claude` 也會和以前一樣運作。
