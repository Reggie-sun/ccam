---
title: Claude Codeが使用量の上限に達したときの対処法
description: リセットを待つか、別のClaude Codeアカウントですぐに同じ会話を続けます。5時間の上限と週間の上限のしくみと対処の手順を説明します。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

アカウントが使用量の上限（レート制限とも呼ばれます）に達したときの選択肢は2つです。上限がリセットされるまで待つか、自分が所有する別のアカウントで同じ会話をすぐに続けるかです。ccamを使う場合、アカウント間で会話履歴を共有しておけば、後者は`ccam next --auto --continue`という1つのコマンドで済みます。

## Claude Codeにはどんな上限がありますか？

ProプランとMaxプランでは、Claude Codeに[5時間のセッション上限と週間の上限](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)があります。どちらもWebやアプリのClaudeと共通で、使用量は両方の上限に同時に加算されます。どちらかを使い切ると、Claude Codeは「You've hit your session limit · resets 3:45pm」のような[メッセージを表示します](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)。上限はすべてのモデルに適用されるため、モデルを切り替えても解決しません。

## 方法1：リセットを待つ

アカウントが1つなら、待つことになります。最近のバージョンのClaude Codeは、[開いているセッションのまま待機し](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)、上限がリセットされると自動的にタスクを再開します。Claude Codeの中で`/usage`を実行すると、プランの上限をどれだけ使ったかを確認できます。

## 方法2：今すぐ別のアカウントで続ける

2つ目のアカウントがあれば、同じ会話をすぐに続けられます。一度だけ設定しておきます。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

これで、作業中はステータスラインにアカウントと5時間の使用量が表示されるので、上限に達して驚くことはありません。

上限に達したら、Claude Codeを終了（`/exit`）してから、次のコマンドを実行します。

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

最も余裕があるアカウントで会話を開きます。

```sh
ccam next --auto --continue
```

Claude Codeが`work`として起動し、現在のフォルダーで最も新しい会話を、それまでのやり取りをすべて含めて開きます。

## ccamは自動でアカウントを切り替えますか？

いいえ。作業中にccamは動いていません。ccamはClaude Codeを起動し、自身をClaude Codeに置き換えます。切り替えるタイミングはあなたが決めます。`ccam next --auto`は、5時間の使用量が最も少ないアカウントを選び、100%に達したアカウントやログインしていないアカウントは飛ばします。実行中のセッションのログインを入れ替えるツールもあります。違いは[ccamとほかのアカウント切り替えツール](/compare/account-switchers)で説明しています。

## 注意点

- `ccam next --continue`は、作業していたのと同じフォルダーで実行してください。会話はプロジェクトのフォルダーごとに保存されます。
- 2つのアカウントで同じ会話を同時に使うのは避けてください。そうしようとすると、ccamが警告します。
- もう一方のアカウントも、同じグループで会話履歴を共有している必要があります。
- 新しいアカウントではまだ会話がキャッシュされていないため、最初のメッセージで通常より多く上限を消費することがあります。
- このようにアカウントを切り替えてよいかどうかは、各アカウントの規約によります。

`ccam usage`と`ccam next`のすべてのオプションについては、[使用量の上限](/docs/usage-limits)を参照してください。
