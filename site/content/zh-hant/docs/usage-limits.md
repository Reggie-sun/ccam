---
title: 用量限制與切換帳號
description: 用 ccam usage 以長條查看每個 Claude Code 帳號的 5 小時、每週和 Fable 用量，並用 ccam next --auto 開啟剩餘額度最多的帳號。
updated: 2026-10-21
order: 5
translation: machine
sourceHash: "0e64e3251c69"
---

`ccam usage` 會顯示每個帳號已使用了多少 5 小時限制、每週限制，以及各模型自己的每週限制（例如 Fable）。`ccam next --auto` 會開啟剩餘額度最多的帳號，而 `--continue` 會把目前的對話帶過去。

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

每個帳號先有一行，顯示數字的新舊程度，接著每項限制各有一個長條，並附上距離重設的時間：5 小時區間、所有模型共用的每週區間，以及各模型自己的每週區間，例如 Fable。長條從 80% 起變成黃色，達到 100% 時變成紅色。最後一行會列出 `ccam next --auto` 會開啟的帳號。

## 數字從哪裡來

ccam 會使用每個帳號自己的登入資訊，同時向 Anthropic 的用量端點查詢所有帳號。每個帳號最多等待 3 秒，並在帳號一有回應時就顯示出來。

這個端點沒有公開的文件，可能會變更。請求失敗時，ccam 會顯示最後已知的數字，並標示為快取資料。在這些數字儲存之後已經重設的時間區間會顯示 0%。如果回應格式有所變更，`ccam doctor` 會告訴你。

60 秒內取得的數字會重複使用。選項：

| 選項 | 效果 |
| --- | --- |
| `--watch [seconds]` | 持續重新整理，預設每 60 秒一次 |
| `--json` | 機器可讀的輸出 |
| `-v` | 另外顯示每個帳號的數字來源 |
| `--no-cache` | 一律重新查詢 |
| `--source cache` | 絕不查詢；只顯示已儲存的數字 |

## 過期的登入

當帳號的存取權杖過期時，ccam 不會更新它，因為在某個 Claude Code 工作階段使用同一個登入時更新權杖，可能會讓該工作階段被登出。用 `ccam <name>` 啟動該帳號一次，Claude Code 就會更新權杖。

`ccam usage --refresh` 會自行更新過期的權杖，但只針對沒有在執行的帳號。它絕不會更新存放在 macOS 鑰匙圈或匯入的 `~/.claude` 中的登入資訊。

## 狀態列中的用量

```sh
ccam statusline install
```

之後 Claude Code 會在狀態列中顯示帳號名稱及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。達到 90% 以上時，這一行也會建議執行 `ccam next --auto --continue`。如果你已經有狀態列指令，ccam 會保留它，並在前面加上帳號名稱。狀態列也會儲存它的數字，當端點沒有回應時，`ccam usage` 會使用這些數字。

## 切換

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` 會略過未登入的帳號，`--auto` 則會略過已達 100% 的帳號和屬於其他組織的帳號。兩個帳號共用對話記錄時，`--continue` 會開啟同一段對話；否則 Claude Code 會把對話複製到新的帳號。請參閱[共用](/docs/sharing)和[切換到下一個帳號](/features/next#continue)。

對於執行 `claude -p` 的指令碼和排程工作，請在帳號名稱前加上 `--fallback`。該帳號達到上限時，ccam 會以剩餘額度最多的帳號再次執行同一個指令，並且只輸出完成的那次執行的結果：

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccam 如何決定，請參閱[達到上限時重新執行](/features/next#run-again-at-the-limit)。

[用量限制工作流程指南](/guides/usage-limit-workflow)會帶你走過一個完整的範例。
