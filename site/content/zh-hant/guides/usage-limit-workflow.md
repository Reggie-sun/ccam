---
title: Claude Code 達到用量限制時該怎麼辦
description: 查看每個 Claude Code 帳號還剩多少額度，並在某個帳號達到限制時，在另一個帳號繼續同一段對話。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

當某個帳號達到用量限制時，可以在你擁有的另一個帳號繼續同一段對話：先設定一次對話記錄共用，然後執行 `ccam next --continue`。如果你只有一個帳號，請等到 Claude Code 顯示的重設時間。

## 事前準備

先設定好帳號，讓切換只需要一個指令：

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

現在你工作時，狀態列會顯示帳號及其 5 小時用量，達到限制時就不會措手不及。

## 達到限制時

離開 Claude Code（`/exit`），然後執行：

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

在剩餘額度最多的帳號開啟這段對話：

```sh
ccam next --auto --continue
```

Claude Code 會以 `work` 啟動，並開啟目前資料夾中最近的一段對話，包含到目前為止的所有內容。

## 注意事項

- 請在原本工作的同一個資料夾中執行 `ccam next --continue`；對話是依專案資料夾分別儲存的。
- 避免同時用兩個帳號處理同一段對話。如果你這麼做，ccam 會發出警告。
- 另一個帳號也必須在同一個群組中共用對話記錄。
- 能否像這樣在帳號之間切換，取決於各帳號的條款。

關於 `ccam usage` 和 `ccam next` 的所有選項，請參閱[用量限制](/docs/usage-limits)。
