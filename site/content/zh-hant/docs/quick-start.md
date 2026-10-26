---
title: 快速入門
metaTitle: "快速入門：新增和切換 Claude Code 帳號"
description: 用 ccam add 新增兩個 Claude Code 帳號，為它們設定 ccwork 這類快捷指令並在帳號之間切換，再選擇要共用的項目。
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

用 `ccam add` 新增每個帳號，然後用 `ccam <name>` 或你自己命名的快捷指令（例如 `ccwork`）啟動 Claude Code。大多數人只需要這些。

已經在使用 Claude Code 了嗎？開始之前，請先閱讀本頁最後的[已經在使用 ~/.claude？](#already-using-claude)。

## 1. 新增第一個帳號

```sh
ccam add
```

設定精靈會要求你輸入名稱，例如 `personal` 或 `work`，並詢問這個帳號要和其他帳號共用哪些項目，以及是否建立快捷指令（請見第 3 步）。接著它會執行 `claude auth login`，開啟瀏覽器讓你登入。在遠端機器上，你的瀏覽器可能會停在無法載入的 `http://localhost:…/callback` 頁面；把這個網址貼到終端機中，ccam 就會完成登入。登入完成後，你會看到：

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. 移轉現有的設定

如果你以前用過 Claude Code，你的設定、`CLAUDE.md` 和擴充功能都在 `~/.claude` 中，MCP 伺服器則在 `~/.claude.json` 中。把它們複製到共用的 `default` 群組，每個帳號就都能使用：

```sh
ccam migrate
```

ccam 會先顯示要複製的內容並詢問你。它只會讀取 `~/.claude`，絕不複製登入資訊；除非加上 `--history`，否則不會複製對話記錄。詳情請參閱[複製現有的設定](/docs/sharing#copy-an-existing-setup)。

## 3. 新增第二個帳號

```sh
ccam add work
```

設定精靈會略過名稱（你已在命令列中指定），並顯示可加入的群組，其中 `default` 已勾選。這次請為帳號取一個專屬的快捷指令，例如 `ccwork`，這樣只要一個詞就能啟動它。選完群組後，設定精靈會詢問：

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

按下 `y`，然後輸入名稱，或直接按 Enter 使用建議的名稱：

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

接著在瀏覽器中用另一個帳號登入。ccam 會確認帳號及其快捷指令：

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. 切換帳號

`ccam work` 會以 work 帳號啟動 Claude Code。如果你跳過了快捷指令的問題，隨時都可以新增快捷指令：

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

現在，以下每個指令都會以 work 帳號啟動 Claude Code（`2` 是它在 `ccam list` 中的編號）：

```sh
ccwork
ccam work
ccam 2
```

單獨執行 `ccam` 會開啟[儀表板](/features/dashboard)：用方向鍵或帳號的編號移到某個帳號，再按 Enter。

快捷指令或帳號名稱後面的所有內容都會原封不動地傳給 Claude Code，所以你可以繼續使用熟悉的選項：

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 選擇要共用的項目

預設情況下，MCP 伺服器、`settings.json`、全域 `CLAUDE.md` 和擴充功能（agents、commands、skills 和 plugins）都會透過 `default` 群組在帳號之間共用，對話記錄則不會。如果想在帳號之間接續對話，請把它加到群組中：

```sh
ccam group set default --add history
```

ccam 在變更任何內容之前，會先顯示它將為每個帳號移動哪些項目。詳情請參閱[用群組共用](/features/share)。

## 6. 帳號達到限制時

```sh
ccam usage
ccam next --continue
```

`ccam usage` 會顯示每個帳號的 5 小時用量和每週用量。`ccam next --continue` 會在下一個帳號開啟同一段對話。在指令碼中，`ccam --fallback work -p "…"` 會自行切換。請參閱[用量限制](/docs/usage-limits)。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。

## 已經在使用 ~/.claude？

直接登錄它，不必重新登入：

```sh
ccam import
```

資料夾會留在原處，直接執行 `claude` 也會和以前一樣運作。
