---
title: 為指令碼選擇帳號
description: ccam pick 為指令碼印出剩餘額度最多的帳號，讓同時啟動的工作程序使用不同帳號，沒有帳號有額度時以 75 結束。
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

`ccam pick` 會印出剩餘額度最多的已登入帳號名稱，讓指令碼用它啟動 Claude Code。只要有足夠多還有額度的帳號，同時啟動的工作程序就會各自拿到不同的帳號，所以一批 `claude -p` 工作會分散到你的各個帳號，而不是把一個帳號用光。

## 用法

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live"
}
```

結束代碼：`0` 表示在標準輸出上印出了名稱，`75`（`EX_TEMPFAIL`）表示沒有帳號還有額度或沒有帳號登入，`1` 表示其他錯誤，例如錯誤的旗標。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--json` | 印出帳號、它的 5 小時和每週用量，以及這些數字的來源 | 關閉 |
| `--model <model>` | 也計入該模型的每週區間，例如 `opus` | 無 |
| `--max-usage <1-100>` | 優先選擇用量低於此值的帳號 | `auto_max_usage`，否則為 100 |
| `--strategy headroom\|expiring` | 5 小時用量最低，或 5 小時區間最早重設 | `auto_strategy`，否則為 `headroom` |
| `--wait[=<duration>]` | 沒有帳號還有額度時，等待第一個帳號重設，最多等這麼久 | 關閉；單獨指定時為 `5h` |
| `--dry-run` | 只選擇，不記錄這次選擇 | 關閉 |

## 運作方式

1. ccam 列出已登入的帳號（API 金鑰帳號沒有限制，因此不列入）。
2. 它取得一個鎖定（`~/.local/state/ccam/picks.lock`），所以同時執行的選擇會依序進行。
3. 它以和 `ccam usage` 相同的方式一次查詢所有帳號的用量，並使用 60 秒的快取。
4. 先排除過去 2 分鐘內被選過的帳號（`~/.local/state/ccam/picks/<name>`），再由 `ccam next --auto` 的[選擇規則](/features/next#auto)在其餘帳號中選擇：優先選低於 `auto_max_usage` 的帳號，依策略選擇，絕不選擇用量 100% 的帳號。
5. 如果每個還有額度的帳號最近都被選過，規則會再次在所有帳號中選擇，所以三個帳號上的第四個工作程序會和其他工作程序共用一個帳號。
6. 記錄這次選擇並釋放鎖定，然後印出名稱。

2 分鐘的時間窗很重要，因為新工作的用量要過一段時間才會反映在數字上：沒有它的話，同時啟動的五個工作程序都會把同一個帳號看成剩餘額度最多的帳號。

## 安全性和特殊情況

- `ccam pick` 只讀取用量並寫入選擇的時間；除了 `ccam usage` 會做的事之外，它絕不更新或讀取權杖。
- 選擇並不會保留帳號。2 分鐘後，兩個指令碼仍可能在同一個帳號上啟動；請使用 `--fallback`，在該帳號達到上限時恢復。
- 以 `--wait` 等待時不持有鎖定，所以其他選擇可以同時進行。
- `ccam pick` 會忽略 [`ccam folder`](/features/folder)：它只依用量選擇。

## 相關內容

- [切換到下一個帳號](/features/next)
- [達到上限時重新執行](/features/next#run-again-at-the-limit)
- [用量限制](/features/usage)
- [API 金鑰帳號](/features/api)
