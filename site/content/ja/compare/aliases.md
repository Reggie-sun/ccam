---
title: ccamとシェルエイリアスの比較
description: 複数のClaude Codeアカウントを使う方法として、CLAUDE_CONFIG_DIRを使うシェルエイリアスとccamを機能ごとに比較します。
updated: 2026-10-26
order: 1
translation: machine
sourceHash: c30d84bceba0
---

アカウントごとにシェルエイリアスを作るのが、複数のClaude Codeアカウントを使う最も簡単な方法です。ccamも内部では同じことを行い、そのうえで共有、切り替え、使用量の確認の機能を加えています。

## エイリアスでの設定

追加するアカウントごとに`~/.bashrc`または`~/.zshrc`に1行を加え、新しいターミナルを開いて、エイリアスを一度実行してログインします。

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude`は1つ目のアカウントのまま使え、`claude-work`で2つ目のアカウントを起動します。これは[Claude Codeのドキュメントで説明されている方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)です。

## ccamでの設定

ccamをインストールしてアカウントを追加します。ここではアカウント専用のショートカットコマンドも作ります。

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work`と`ccwork`のどちらでもアカウントを起動でき、どちらもスクリプトから使えます。

## 並べて比較

| | シェルエイリアス | ccam |
| --- | --- | --- |
| ログインを分ける | 可 | 可 |
| アカウントごとの設定作業 | `~/.bashrc`を編集し、ログインする | `ccam add` |
| スクリプトやほかのプログラムからの利用 | 不可。エイリアスはシェル内でのみ有効 | 可。`ccam work`も、`ccwork`のようなショートカットも実際のコマンド |
| 設定、`CLAUDE.md`、拡張機能の共有 | 手動でコピーまたはリンク | 可、項目ごとに選択、バックアップ付き |
| MCPサーバーの共有 | アカウントごとに手動で | `ccam mcp add`を一度実行 |
| 別のアカウントで会話を続ける | 不可 | `ccam next --continue` |
| すべてのアカウントの使用量 | 不可 | `ccam usage` |
| 最も余裕があるアカウントを選ぶ | 不可 | `ccam next --auto` |
| 上限に達した`-p`ジョブを別のアカウントで再実行する | 不可 | `ccam --fallback` |
| プロジェクトのフォルダーごとのアカウント | 不可 | `ccam folder` |
| ステータスラインへのアカウント名の表示 | 不可 | `ccam statusline install` |
| 壊れたリンクの修復 | 手動 | `ccam doctor --fix` |
| 追加のソフトウェア | なし | 静的バイナリ1つ |

アカウントが2つで、アカウント間で何も共有したくない場合は、エイリアスを選んでください。複数のアカウントを1つの環境のように扱いたい場合は、ccamを選んでください。最初はエイリアスで始めて、後からccamに移行することもできます。`ccam import`は、既存の`~/.claude-*`フォルダーをその場所のまま登録します。

[1台のコンピューターで複数のアカウントを使う方法](/guides/multiple-claude-code-accounts)で両方の方法を順を追って説明しています。[ccamとほかのアカウント切り替えツール](/compare/account-switchers)では、1つのログインを入れ替えるツールとccamを比較しています。
