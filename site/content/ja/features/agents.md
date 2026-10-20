---
title: すべてのアカウントのバックグラウンドセッション
description: ccam agentsは、すべてのアカウントのClaude Codeのバックグラウンドエージェントと開いているセッションを1つの表に一覧表示し、ccam attachはIDで1つを開きます。
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Codeはバックグラウンドセッション（`claude --bg`、エージェントビュー）をアカウントごとに保持します。`claude agents`が表示するのは、実行しているアカウントのセッションだけです。`ccam agents`はすべてのアカウントに一度に問い合わせてまとめて一覧表示し、`ccam attach <id>`はバックグラウンドセッションを、それを開始したアカウントで開きます。

## 使い方

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State`は、バックグラウンドセッションの状態（`working`、`blocked`、`done`、`failed`、`stopped`）と、そのプロセスが何をしているか（`busy`、`waiting`、`idle`）です。開いている対話型のセッションには`open`と表示されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--all` | 完了、失敗、停止したバックグラウンドセッションも一覧表示します | オフ |
| `--json` | セッションをJSONで出力します。それぞれに所属する`account`が付きます | オフ |

`ccam attach <id> [claude args]`は、IDより後の引数を`claude attach`に渡します。`--dry-run`を付けると、コマンドを表示します。

## しくみ

1. ccamはすべてのアカウントについて、アカウントの`CLAUDE_CONFIG_DIR`を設定して`claude agents --json`（指定があれば`--all`付き）を同時に実行し、それぞれ15秒以内の応答を待ちます。これはアカウント自身のセッション一覧を読むだけで、Claude Codeのバックグラウンドサービスは起動しません。
2. 各エントリーにアカウント名を付け、一覧を開始時刻の新しい順に並べます。応答しないアカウントは、標準エラー出力に名前を表示して一覧から除きます。その場合、ccamは終了コード1で終わるので、スクリプトが失敗を空の一覧と取り違えることはありません。
3. `ccam attach`は、Claude Codeがバックグラウンドセッションを保存する`jobs/<id>`を各アカウントのフォルダーで探し、そのアカウントとして`claude attach <id>`を実行します。`ccam <name>`がClaude Codeを起動するのと同じ方法です。

## 安全性と特殊なケース

- バックグラウンドセッションのIDは1つのアカウントに属します。別のアカウントとして実行した`claude attach`ではIDが見つかりません。これを解決するのが`ccam attach`です。
- 2つのアカウントが、同じ文字で始まるIDのセッションを持つことがあります。その場合、`ccam attach`は両方を示し、IDをもっと長く指定するよう求めます。
- ccamは共有MCPサーバーを`--mcp-config=<file>`という1つの引数でClaude Codeに渡します。`claude agents`と`claude attach`はこれを受け付け、バックグラウンドセッションは再起動するときもそのファイルを使い続けます。
- APIキーのアカウントのバックグラウンドサービスは、それを起動したシェルからキーを受け取ります。そのようなアカウントの最初のバックグラウンドセッションはccam経由で起動し（`ccam gw --bg …`）、サービスがキーを持てるようにしてください。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [スクリプト用のアカウント選び](/features/pick)
- [ダッシュボード](/features/dashboard)
