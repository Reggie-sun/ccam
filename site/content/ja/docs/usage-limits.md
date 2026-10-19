---
title: 使用量の上限とアカウントの切り替え
description: ccam usageですべてのClaude Codeアカウントの5時間と週間の使用量を確認し、ccam next --autoで最も余裕があるアカウントを開きます。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 8ab41501c548
---

`ccam usage`は、各アカウントが5時間と週間の上限をどれだけ使ったかを表示します。`ccam next --auto`は最も余裕があるアカウントを開き、`--continue`を付けると現在の会話を引き継ぎます。

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

各列は、アカウント、5時間の使用量、リセットまでの時間、週間の使用量、リセットまでの時間、数値がどれだけ前のものか、を表します。

## 数値の取得元

ccamは、各アカウント自身のログインを使って、すべてのアカウントについて同時にAnthropicの使用量エンドポイントに問い合わせます。待つのは1アカウントにつき最大3秒で、応答があったアカウントから順に表示します。

このエンドポイントは公開ドキュメントに記載されておらず、変更される可能性があります。リクエストが失敗した場合、ccamは最後に取得した数値を表示し、キャッシュであることを示します。応答の形式が変わった場合は、`ccam doctor`が知らせます。

取得してから60秒未満の数値は再利用されます。オプションは次のとおりです。

| オプション | 効果 |
| --- | --- |
| `--watch [seconds]` | 継続的に更新します。デフォルトは60秒ごと |
| `--json` | 機械可読な形式で出力します |
| `-v` | モデル別の週間上限と、各数値の取得元も表示します |
| `--no-cache` | 常に問い合わせ直します |
| `--source cache` | 問い合わせを行わず、保存済みの数値だけを表示します |

## 期限切れのログイン

アカウントのアクセストークンの期限が切れていても、ccamはトークンを更新しません。同じログインを使っているClaude Codeのセッションがあるときに更新すると、そのセッションがサインアウトされる可能性があるためです。`ccam <name>`でそのアカウントを一度起動すれば、Claude Codeがトークンを更新します。

`ccam usage --refresh`は期限切れのトークンを自分で更新しますが、対象は実行中でないアカウントだけです。macOSのキーチェーンに保存されたログインや、インポートした`~/.claude`のログインは更新しません。

## ステータスラインでの使用量表示

```sh
ccam statusline install
```

これで、Claude Codeのステータスラインにアカウント名とその使用量が表示されます。たとえば`[work] · Opus · 5h 24% · 7d 41%`のようになります。すでにステータスラインのコマンドを設定している場合、ccamはそれを残し、先頭にアカウント名を追加します。ステータスラインは取得した数値の保存も行い、エンドポイントが応答しないときは`ccam usage`がその数値を使います。

## 切り替え

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next`はログインしていないアカウントを飛ばし、`--auto`は使用量が100%のアカウントを飛ばします。会話を続けるには、両方のアカウントが会話履歴を共有している必要があります。詳しくは[共有](/docs/sharing)を参照してください。

[使用量の上限に達したときのワークフローガイド](/guides/usage-limit-workflow)で、一連の流れを例に沿って説明しています。
