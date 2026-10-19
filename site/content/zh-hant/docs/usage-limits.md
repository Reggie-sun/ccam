---
title: 用量限制與切換帳號
description: 用 ccam usage 查看每個 Claude Code 帳號的 5 小時和每週用量，並用 ccam next --auto 開啟剩餘額度最多的帳號。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 8ab41501c548
---

`ccam usage` 會顯示每個帳號已使用了多少 5 小時限制和每週限制。`ccam next --auto` 會開啟剩餘額度最多的帳號，而 `--continue` 會把目前的對話帶過去。

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

各欄依序為：帳號、5 小時用量、距離重設的時間、每週用量、距離重設的時間，以及數字的新舊程度。

## 數字從哪裡來

ccam 會使用每個帳號自己的登入資訊，同時向 Anthropic 的用量端點查詢所有帳號。每個帳號最多等待 3 秒，並在帳號一有回應時就顯示出來。

這個端點沒有公開的文件，可能會變更。請求失敗時，ccam 會顯示最後已知的數字，並標示為快取資料。如果回應格式有所變更，`ccam doctor` 會告訴你。

60 秒內取得的數字會重複使用。選項：

| 選項 | 效果 |
| --- | --- |
| `--watch [seconds]` | 持續重新整理，預設每 60 秒一次 |
| `--json` | 機器可讀的輸出 |
| `-v` | 另外顯示各模型的每週限制，以及每個數字的來源 |
| `--no-cache` | 一律重新查詢 |
| `--source cache` | 絕不查詢；只顯示已儲存的數字 |

## 過期的登入

當帳號的存取權杖過期時，ccam 不會更新它，因為在某個 Claude Code 工作階段使用同一個登入時更新權杖，可能會讓該工作階段被登出。用 `ccam <name>` 啟動該帳號一次，Claude Code 就會更新權杖。

`ccam usage --refresh` 會自行更新過期的權杖，但只針對沒有在執行的帳號。它絕不會更新存放在 macOS 鑰匙圈或匯入的 `~/.claude` 中的登入資訊。

## 狀態列中的用量

```sh
ccam statusline install
```

之後 Claude Code 會在狀態列中顯示帳號名稱及其用量，例如 `[work] · Opus · 5h 24% · 7d 41%`。如果你已經有狀態列指令，ccam 會保留它，並在前面加上帳號名稱。狀態列也會儲存它的數字，當端點沒有回應時，`ccam usage` 會使用這些數字。

## 切換

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` 會略過未登入的帳號，`--auto` 則會略過已達 100% 的帳號。要接續對話，兩個帳號都必須共用對話記錄；請參閱[共用](/docs/sharing)。

[用量限制工作流程指南](/guides/usage-limit-workflow)會帶你走過一個完整的範例。
