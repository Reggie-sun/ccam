---
title: Claude Codeが使用量の上限に達したときの対処法
description: 各Claude Codeアカウントの残りの余裕を確認し、あるアカウントが上限に達したら別のアカウントで同じ会話を続けます。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

あるアカウントが上限に達したら、自分が所有する別のアカウントで同じ会話を続けます。一度だけ会話履歴を共有しておき、`ccam next --continue`を実行してください。アカウントが1つしかない場合は、Claude Codeが表示するリセット時刻まで待ちます。

## 上限に達する前に

1つのコマンドで切り替えられるように、アカウントを設定しておきます。

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

これで、作業中はステータスラインにアカウントと5時間の使用量が表示されるので、上限に達して驚くことはありません。

## 上限に達したとき

Claude Codeを終了（`/exit`）してから、次のコマンドを実行します。

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

## 注意点

- `ccam next --continue`は、作業していたのと同じフォルダーで実行してください。会話はプロジェクトのフォルダーごとに保存されます。
- 2つのアカウントで同じ会話を同時に使うのは避けてください。そうしようとすると、ccamが警告します。
- もう一方のアカウントも、同じグループで会話履歴を共有している必要があります。
- このようにアカウントを切り替えてよいかどうかは、各アカウントの規約によります。

`ccam usage`と`ccam next`のすべてのオプションについては、[使用量の上限](/docs/usage-limits)を参照してください。
