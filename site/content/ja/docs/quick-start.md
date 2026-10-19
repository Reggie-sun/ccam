---
title: クイックスタート
description: ccam addで2つのClaude Codeアカウントを追加し、claude1、claude2、ccam <name>で切り替え、共有する項目を選びます。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

`ccam add`で各アカウントを追加し、`claude1`、`claude2`、または`ccam <name>`でClaude Codeを起動します。ほとんどの人にとって、必要なのはこれだけです。

## 1. 最初のアカウントを追加する

```sh
ccam add
```

ウィザードが、`personal`や`work`のような名前と、そのアカウントがほかのアカウントと共有する項目を尋ねます。次に`claude auth login`を実行し、サインインできるようにブラウザーを開きます。完了すると、次のように表示されます。

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccamは独自のメッセージを英語で表示します。コマンドとオプションはどの言語でも同じです。

## 2. 2つ目のアカウントを追加する

```sh
ccam add work
```

ブラウザーで、もう一方のアカウントでサインインします。このアカウントが`claude2`になります。

## 3. 切り替える

次のどのコマンドでも、仕事用のアカウントでClaude Codeが起動します。

```sh
claude2
ccam work
ccam 2
```

アカウント名より後のものはすべてそのままClaude Codeに渡されるので、使い慣れたオプションをそのまま使えます。

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

`ccam`だけを実行すると、選択画面が開きます。矢印キーとEnterキーを使うか、アカウントの番号を押してください。

## 4. 共有する項目を選ぶ

デフォルトでは、MCPサーバー、`settings.json`、グローバルの`CLAUDE.md`、そしてエージェント、コマンド、スキル、プラグインがアカウント間で共有されます。会話履歴は共有されません。アカウントをまたいで会話を続けたい場合は、両方のアカウントで会話履歴の共有を有効にします。

```sh
ccam share personal --history on
ccam share work --history on
```

ccamは、何かを変更する前に、何を移動するかを表示します。詳しくは[共有](/docs/sharing)を参照してください。

## 5. アカウントが上限に達したとき

```sh
ccam usage
ccam next --continue
```

`ccam usage`は、すべてのアカウントの5時間と週間の使用量を表示します。`ccam next --continue`は、次のアカウントで同じ会話を開きます。詳しくは[使用量の上限](/docs/usage-limits)を参照してください。

## すでに~/.claudeを使っている場合

もう一度サインインする代わりに、既存のフォルダーを登録します。

```sh
ccam import
```

フォルダーはその場所のまま残り、通常の`claude`もこれまでどおり動作します。
