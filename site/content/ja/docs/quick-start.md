---
title: クイックスタート
metaTitle: "クイックスタート：Claude Codeアカウントの追加と切替"
description: ccam addで2つのClaude Codeアカウントを追加し、ccworkのようなショートカットを付けて切り替え、共有する項目を選びます。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: a10d91d2a64e
---

`ccam add`で各アカウントを追加し、`ccam <name>`、または`ccwork`のように自分で名前を付けたショートカットでClaude Codeを起動します。ほとんどの人にとって、必要なのはこれだけです。

## 1. 最初のアカウントを追加する

```sh
ccam add
```

ウィザードが、`personal`や`work`のような名前、そのアカウントがほかのアカウントと共有する項目、ショートカットコマンドを作成するかどうか（ステップ4を参照）を尋ねます。次に`claude auth login`を実行し、サインインできるようにブラウザーを開きます。リモートのマシンでは、ブラウザーが最後に読み込まれない`http://localhost:…/callback`のページで止まることがあります。そのアドレスをターミナルに貼り付けると、ccamがサインインを完了します。完了すると、次のように表示されます。

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. 既存の設定を引き継ぐ

以前からClaude Codeを使っている場合、設定、`CLAUDE.md`、拡張機能、MCPサーバーは`~/.claude`にあります。これらを共有設定にコピーすると、すべてのアカウントで使えるようになります。

```sh
ccam migrate
```

ccamは、コピーする内容を表示してから確認を求めます。`~/.claude`は読み取るだけで、ログインは一切コピーせず、`--history`を付けない限り会話履歴も含めません。詳しくは[既存の環境をコピーする](/docs/sharing#copy-an-existing-setup)を参照してください。

## 3. 2つ目のアカウントを追加する

```sh
ccam add work
```

ブラウザーで、もう一方のアカウントでサインインします。

## 4. アカウントを切り替える

`ccam work`で、仕事用のアカウントでClaude Codeが起動します。もっと短くしたい場合は、`ccwork`のような専用のショートカットコマンドをアカウントに付けます。共有する項目を選んだ後、ウィザードが次のように尋ねます。

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y`を押してから名前を入力するか、Enterキーを押して提案された名前を使います。

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

ログインが終わると、ccamが結果を表示します。

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

この質問を飛ばした場合でも、いつでもショートカットを追加できます。

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

これで、次のどのコマンドでも、仕事用のアカウントでClaude Codeが起動します。

```sh
ccwork
ccam work
ccam 2
```

`ccam`だけを実行すると、選択画面が開きます。矢印キーとEnterキーを使うか、アカウントの番号を押してください。

ショートカットやアカウント名より後のものはすべてそのままClaude Codeに渡されるので、使い慣れたオプションをそのまま使えます。

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 共有する項目を選ぶ

デフォルトでは、MCPサーバー、`settings.json`、グローバルの`CLAUDE.md`、そしてエージェント、コマンド、スキル、プラグインがアカウント間で共有されます。会話履歴は共有されません。アカウントをまたいで会話を続けたい場合は、両方のアカウントで会話履歴の共有を有効にします。

```sh
ccam share personal --history on
ccam share work --history on
```

ccamは、何かを変更する前に、何を移動するかを表示します。詳しくは[共有](/docs/sharing)を参照してください。

## 6. アカウントが上限に達したとき

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
