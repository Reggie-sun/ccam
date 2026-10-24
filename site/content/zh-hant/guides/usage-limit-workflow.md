---
title: Claude Code 達到用量限制時該怎麼辦
description: 等待額度重置，或立即在另一個 Claude Code 帳號接續同一個對話。逐步說明 5 小時限制和每週限制如何運作。
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

當一個帳號達到用量限制（常被稱為 rate limit）時，你有兩個選擇：等到限制重置，或現在就在你擁有的另一個帳號接續同一個對話。使用 ccam 時，第二種做法只需要一個指令：`ccam next --auto --continue`。

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
ccam group set default --add history
ccam statusline install
```

現在，工作時狀態列會顯示目前的帳號和它的 5 小時用量，達到限制時就不會措手不及。

達到限制後，先離開 Claude Code（`/exit`），然後執行：

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

在剩餘額度最多的帳號中開啟這個對話：

```sh
ccam next --auto --continue
```

Claude Code 會以 `work` 帳號啟動，並開啟目前資料夾中最近的對話，包含之前的所有內容。共用對話記錄（上面的 `ccam group set default --add history`）能讓它維持同一個對話；沒有共用時，Claude Code 會把對話複製到 `work`，ccam 也會告訴你第一次回覆會重新讀取多少 token。

在 Claude Code 中安裝 [ccam 外掛](/docs/install#claude-code-plugin)後，因上限而停止的工作階段會顯示通知並記錄已滿的區間，讓 `ccam next --auto` 在重設前不會再次選擇該帳號，而 `/ccam:switch`（或 `/ccam:switch work`）則會為 `ccam next --continue` 記下這個對話，什麼都不必猜。

**Team 和 Enterprise 席位。** 自動選擇（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換。Anthropic 說明中心關於 [Team 方案](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)的文章寫道：「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Team 方案的每個席位都指派給一位具名的個人，而且只有這個人可以使用），以及「Users can't route their work through another user's account to use that user's remaining usage」（使用者不能經由其他使用者的帳號處理自己的工作，藉此使用對方剩餘的用量）；關於 [Enterprise 方案](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)的文章也是如此。因此，在席位之間使用自動選擇可能會違反 Anthropic 的條款。在依賴這項功能之前，請先向你所屬的組織和 Anthropic 確認。個人方案（Pro、Max）的情況請參閱[常見問題中關於條款的說明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)。

## ccam 會自動切換帳號嗎？

不會。你工作時 ccam 並沒有在執行：它啟動 Claude Code 後就用 Claude Code 取代自己。何時切換由你決定，`ccam next --auto` 會選出 5 小時用量最低的帳號，並略過用量已達 100%、未登入或屬於其他組織的帳號。在 `config.toml` 中設定 `relay = true` 時，ccam 會改為在工作階段旁等待，讓 `/ccam:switch` 能自行接續下去，並沿用工作階段的選項、權限模式和推理強度。另一些工具則是替換執行中工作階段的登入資訊；[ccam 與其他帳號切換工具](/compare/account-switchers)說明了兩者的差異。

## 注意事項

- 在你工作時所在的同一個資料夾中執行 `ccam next --continue`；對話是依專案資料夾儲存的。
- 避免讓兩個帳號同時處理同一個對話。你這樣做時，ccam 會發出警告。
- 沒有共用對話記錄時，對話會被複製到另一個帳號，之後兩份副本各自發展；`/rewind` 用來還原程式碼的檢查點會留在第一個帳號。
- 新帳號上的第一則訊息可能會比平常消耗更多額度，因為該帳號還沒有快取這個對話。
- 是否可以這樣在帳號之間切換，取決於各帳號的條款。

`ccam usage` 和 `ccam next` 的所有選項請參閱[用量限制](/docs/usage-limits)。
