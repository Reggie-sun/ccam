---
title: 使用量の上限とアカウントの切り替え
description: ccam usageですべてのClaude Codeアカウントの5時間と週間の使用量を確認し、ccam next --autoで最も余裕があるアカウントを開きます。
updated: 2026-10-25
order: 5
translation: machine
sourceHash: "7c15419446c3"
---

`ccam usage`は、各アカウントが5時間の上限、週間の上限、モデルごとの週間上限（Fableなど）をどれだけ使ったかを表示します。`ccam next --auto`は最も余裕があるアカウントを開き、`--continue`を付けると現在の会話を引き継ぎます。

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

アカウントごとに、数値がどれだけ前のものかを示す行があり、その下に上限ごとのバーとリセットまでの時間が並びます。5時間の枠、全モデル共通の週間の枠、そしてFableのようなモデルごとの週間の枠です。バーは80%から黄色、100%で赤になります。最後の行は、`ccam next --auto`が開くアカウントを示します。`auto_strategy = "expiring"`または`"balance"`の場合は、代わりに`→ next --auto picks <name>: <reason>`と表示されます。

**TeamとEnterpriseのシート。** 自動的な選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。そのため、シート間で自動的な選択を使うと、Anthropicの規約に違反するおそれがあります。そこで使うかどうかはご自身の判断です。自動的な選択に頼る前に、所属する組織とAnthropicに確認してください。個人プラン（Pro、Max）については、[FAQの規約に関する説明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を参照してください。

## 数値の取得元

ccamは、各アカウント自身のログインを使って、すべてのアカウントについて同時にAnthropicの使用量エンドポイントに問い合わせます。待つのは1アカウントにつき最大3秒で、応答があったアカウントから順に表示します。

このエンドポイントは公開ドキュメントに記載されておらず、変更される可能性があります。リクエストが失敗した場合、ccamは最後に取得した数値を表示し、キャッシュであることを示します。数値の保存後にリセットされた時間枠は0%と表示されます。応答の形式が変わった場合は、`ccam doctor`が知らせます。

Claude Codeが使用量の上限でセッションを止めると、[ccamプラグイン](/docs/install#claude-code-plugin)が、どの時間枠が上限に達したかと、それがいつリセットされるかを記録します。それまでの間、その時間枠はほかの数値にかかわらず100%と表示されるので、`ccam next --auto`と`--fallback`がそのアカウントを早すぎる時点で再び選ぶことはありません。

取得してから60秒未満の数値は再利用されます。オプションは次のとおりです。

| オプション | 効果 |
| --- | --- |
| `--watch[=seconds]` | 継続的に更新します。デフォルトは60秒ごと |
| `--json` | 機械可読な形式で出力します |
| `-v` | 各アカウントの数値の取得元も表示します |
| `--no-cache` | 常に問い合わせ直します |
| `--source cache` | 問い合わせを行わず、保存済みの数値だけを表示します |

## 期限切れのログイン

アカウントのアクセストークンの期限が切れていても、ccamはトークンを更新しません。同じログインを使っているClaude Codeのセッションがあるときに更新すると、そのセッションがサインアウトされる可能性があるためです。そのアカウントには``access token expired; `ccam <name>` renews it when it starts``と表示されます。`ccam <name>`で一度起動すれば、Claude Codeがトークンを更新します。ログイン全体の期限が切れた場合や、Claude Codeがログインを消去した場合は`` login expired → `ccam login <name>` ``と表示され、新しくログインする必要があります。

`ccam usage --refresh`は期限切れのトークンを自分で更新しますが、対象は、Claude Codeがどう起動されたかにかかわらず、実行中でないアカウントだけです。macOSのキーチェーンに保存されたログインや、インポートした`~/.claude`のログインは更新せず、更新には最大30秒かかることがあります。

## ステータスラインでの使用量表示

```sh
ccam statusline install
```

これで、Claude Codeのステータスラインにアカウント名とその使用量が表示されます。たとえば`[work] · Opus · 5h 24% · 7d 41%`のようになります。90%以上になると、ステータスラインは`ccam next --auto --continue`も提案します。すでにステータスラインのコマンドを設定している場合、ccamはそれを残し、先頭にアカウント名を追加します。ステータスラインは取得した数値の保存も行い、エンドポイントが応答しないときは`ccam usage`がその数値を使います。

## 切り替え

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next`はログインしていないアカウントを飛ばし、`--auto`は使用量が100%のアカウント（Claude Codeが上限に達したと報告した時間枠も、そのリセットまで含みます）と、別の組織のアカウントを飛ばします。2つのアカウントが会話履歴を共有していれば、`--continue`は同じ会話を開きます。共有していなければ、Claude Codeが会話を新しいアカウントにコピーします。詳しくは[共有](/docs/sharing)と[次のアカウントへの切り替え](/features/next#continue)を参照してください。

`config.toml`の先頭で`auto_switch = true`を設定すると（LinuxとmacOS。デフォルトではオフ）、Claude Codeが対話型のセッションを上限で止め、同じ種類の別のアカウントに余裕があるときに、ccamがこの切り替えを自動で行います。いつ動作するか、そして決して行わないことは、[`auto_switch`](/features/next#auto-switch)を参照してください。

`claude -p`を実行するスクリプトや定期実行ジョブでは、アカウント名の前に`--fallback`を付けます。そのアカウントが上限に達すると、ccamは同じコマンドを最も余裕があるアカウントで再実行し、最後まで完了した実行の出力だけを表示します。

```sh
ccam --fallback work -p "review the diff" --output-format json
```

ccamがどのように判断するかは、[上限に達したら再実行](/features/next#run-again-at-the-limit)を参照してください。

[使用量の上限に達したときのワークフローガイド](/guides/usage-limit-workflow)で、一連の流れを例に沿って説明しています。
