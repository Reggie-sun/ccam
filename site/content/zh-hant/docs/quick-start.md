---
title: 快速入門
description: 用 ccam add 新增兩個 Claude Code 帳號，用 claude1、claude2 或 ccam <name> 切換，並選擇要共用的項目。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

用 `ccam add` 新增每個帳號，然後用 `claude1`、`claude2` 或 `ccam <name>` 啟動 Claude Code。大多數人只需要這些。

## 1. 新增第一個帳號

```sh
ccam add
```

設定精靈會要求你輸入名稱，例如 `personal` 或 `work`，並詢問這個帳號要和其他帳號共用哪些項目。接著它會執行 `claude auth login`，開啟瀏覽器讓你登入。完成後你會看到：

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam 以英文顯示自己的訊息；指令和選項在每種語言中都相同。

## 2. 新增第二個帳號

```sh
ccam add work
```

在瀏覽器中用另一個帳號登入。這個帳號會成為 `claude2`。

## 3. 切換

以下任一指令都會以 work 帳號啟動 Claude Code：

```sh
claude2
ccam work
ccam 2
```

帳號名稱後面的所有內容都會原封不動地傳給 Claude Code，所以你可以繼續使用熟悉的選項：

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

單獨執行 `ccam` 會開啟選擇器：用方向鍵和 Enter 選擇，或直接按下帳號的編號。

## 4. 選擇要共用的項目

預設情況下，MCP 伺服器、`settings.json`、全域 `CLAUDE.md`，以及 agents、commands、skills 和 plugins 都會在帳號之間共用，對話記錄則不會。如果想在帳號之間接續對話，請為兩個帳號都開啟對話記錄共用：

```sh
ccam share personal --history on
ccam share work --history on
```

ccam 在變更任何內容之前，會先顯示它將移動哪些項目。詳情請參閱[共用](/docs/sharing)。

## 5. 帳號達到限制時

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
