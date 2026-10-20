---
title: スクリプト用のアカウント選び
description: ccam pickは、スクリプトに最も余裕があるアカウントを表示し、同時に起動したワーカーには別々のアカウントを渡し、余裕がなければ75で終了します。
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

`ccam pick`は、ログイン済みのアカウントのうち最も余裕があるアカウントの名前を表示します。スクリプトはその名前でClaude Codeを起動できます。余裕のあるアカウントが十分にある間は、同時に起動したワーカーがそれぞれ別のアカウントを受け取るので、`claude -p`のジョブの一群が1つのアカウントを使い切ることなく、複数のアカウントに分散されます。

## 使い方

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

終了コード: 標準出力に名前を出したときは`0`、余裕のあるアカウントがないか、ログイン済みのアカウントがないときは`75`（`EX_TEMPFAIL`）、不正なフラグなどそのほかのエラーでは`1`です。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--json` | アカウント、その5時間と週間の使用量、数値の取得元を表示します | オフ |
| `--model <model>` | `opus`など、そのモデルの週間の枠も数えます | なし |
| `--max-usage <1-100>` | 使用量がこの値未満のアカウントを優先します | `auto_max_usage`、なければ100 |
| `--strategy headroom\|expiring` | 5時間の使用量が最も少ないもの、または5時間の枠のリセットが最も早いもの | `auto_strategy`、なければ`headroom` |
| `--wait[=<duration>]` | 余裕のあるアカウントがないとき、最初のアカウントがリセットされるまで最長でこの時間だけ待ちます | オフ。値なしで指定すると`5h` |
| `--dry-run` | 選ぶだけで、選択を記録しません | オフ |

## しくみ

1. ccamはログイン済みのアカウントを一覧にします（APIキーのアカウントには上限がないので除外します）。
2. ロック（`~/.local/state/ccam/picks.lock`）を取るので、同時に実行された選択は1つずつ順番に行われます。
3. `ccam usage`と同じ方法で、60秒のキャッシュを使い、すべてのアカウントの使用量を一度に問い合わせます。
4. 直近2分以内に選ばれたアカウント（`~/.local/state/ccam/picks/<name>`）をまず除き、残りの中から`ccam next --auto`の[選択の規則](/features/next#auto)で選びます。`auto_max_usage`未満を優先し、戦略に従い、100%のアカウントは決して選びません。
5. 余裕のあるアカウントがすべて最近選ばれていた場合は、規則がそれらすべての中からもう一度選びます。そのため、3つのアカウントに対する4つ目のワーカーは、どれか1つを共有します。
6. 選択を記録してロックを解放し、名前を表示します。

2分の期間が重要なのは、新しいジョブの使用量が数値に表れるまでに時間がかかるからです。これがないと、同時に起動した5つのワーカーがすべて、同じアカウントを最も余裕があるアカウントと見なしてしまいます。

## 安全性と特殊なケース

- `ccam pick`は使用量を読み、選択した時刻を書き込むだけです。`ccam usage`が行う以上に、トークンを更新したり読んだりすることはありません。
- 選んでもアカウントが予約されるわけではありません。2分が過ぎれば、2つのスクリプトが同じアカウントで起動することもあります。そのアカウントが上限に達したときに立て直すには、`--fallback`を使います。
- `--wait`で待っている間はロックを持たないので、そのほかの選択はその間も進みます。
- `ccam pick`は[`ccam folder`](/features/folder)を無視し、使用量だけで選びます。

## 関連項目

- [次のアカウントへの切り替え](/features/next)
- [上限に達したら再実行](/features/next#run-again-at-the-limit)
- [使用量の上限](/features/usage)
- [APIキーのアカウント](/features/api)
