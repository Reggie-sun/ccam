---
title: Claude Code 達到用量限制時該怎麼辦
description: 等待額度重置，或立即在另一個 Claude Code 帳號接續同一個對話。逐步說明 5 小時限制和每週限制如何運作。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

當一個帳號達到用量限制（常被稱為 rate limit）時，你有兩個選擇：等到限制重置，或現在就在你擁有的另一個帳號接續同一個對話。使用 ccam 時，只要帳號之間共用了對話記錄，第二種做法只需要一個指令：`ccam next --auto --continue`。

## Claude Code 有哪些限制？

在 Pro 和 Max 方案中，Claude Code 有[5 小時工作階段限制和每週限制](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)。這兩個限制與網頁版和應用程式中的 Claude 共用，用量會同時計入兩者。其中一個用完時，Claude Code 會[顯示一則訊息](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)，例如「You've hit your session limit · resets 3:45pm」。切換模型沒有幫助，因為這些限制涵蓋所有模型。

## 做法 1：等待重置

如果你只有一個帳號，就只能等待。較新版本的 Claude Code 會[在開啟的工作階段中等待](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)，限制重置後自動繼續任務。在 Claude Code 中執行 `/usage` 可以查看你已用掉多少方案額度。

## 做法 2：現在就在另一個帳號接續

如果你有第二個帳號，就能立刻接著同一個對話繼續工作。只需設定一次：

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

現在，工作時狀態列會顯示目前的帳號和它的 5 小時用量，達到限制時就不會措手不及。

達到限制後，先離開 Claude Code（`/exit`），然後執行：

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

在剩餘額度最多的帳號中開啟這個對話：

```sh
ccam next --auto --continue
```

Claude Code 會以 `work` 帳號啟動，並開啟目前資料夾中最近的對話，包含之前的所有內容。

## ccam 會自動切換帳號嗎？

不會。你工作時 ccam 並沒有在執行：它啟動 Claude Code 後就用 Claude Code 取代自己。何時切換由你決定，`ccam next --auto` 會選出 5 小時用量最低的帳號，並略過用量已達 100% 或未登入的帳號。另一些工具則是替換執行中工作階段的登入資訊；[ccam 與其他帳號切換工具](/compare/account-switchers)說明了兩者的差異。

## 注意事項

- 在你工作時所在的同一個資料夾中執行 `ccam next --continue`；對話是依專案資料夾儲存的。
- 避免讓兩個帳號同時處理同一個對話。你這樣做時，ccam 會發出警告。
- 另一個帳號也必須在同一個群組中共用對話記錄。
- 新帳號上的第一則訊息可能會比平常消耗更多額度，因為該帳號還沒有快取這個對話。
- 是否可以這樣在帳號之間切換，取決於各帳號的條款。

`ccam usage` 和 `ccam next` 的所有選項請參閱[用量限制](/docs/usage-limits)。
