---
title: コマンドリファレンス
description: ccam addやccam <name>から、ccam usage、ccam next、ccam doctorまで、ccamのすべてのコマンドとオプションを説明します。
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

どのコマンドでも前に`--dry-run`を付けると、実際には実行せずに何をするかを確認できます。たとえば`ccam --dry-run work`のように使います。

<!-- commands -->

## Claude Codeの起動

`ccam <name|number> [arguments]`は、そのアカウントでClaude Codeを起動します。ccamは名前より後の引数を読まず、同じ順序で`claude`に渡します。標準入力、終了コード、シグナルも`claude`そのものとまったく同じように動作します。名前の後の`--`は省略できます。

`ccam -p "summarise"`のように、引数はあってもアカウント名がない`ccam`は、デフォルトのアカウントを使います。デフォルトのアカウントは`ccam use <name>`で設定します。

ショートカットの`claude1`、`claude2`などは、`~/.local/bin`にあるccamへのリンクです。`ccam shims install --names`を実行すると、アカウントごとに`claude-<name>`も作成されます。

## アカウントの追加と削除

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

フォルダーを指定せずに`ccam import`を実行すると、`~/.claude`と`~/.claude-*`を探し、それぞれについて確認します。インポートしたフォルダーはその場所のまま登録されます。ccamがログインを移動またはコピーすることはありません。

## 共有

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

各項目には、`on`（`default`グループ）、`off`、またはグループ名を指定します。詳しくは[共有](/docs/sharing)を参照してください。

## MCPサーバー

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## 使用量と切り替え

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## メンテナンス

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ファイル

| パス | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | アカウントと、それぞれが共有する項目 |
| `~/.local/share/ccam/profiles/<name>/` | 各アカウントの`CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共有される設定、MCPサーバー、拡張機能、履歴 |
| `~/.local/state/ccam/` | 使用量のキャッシュと最終使用日時 |
| `~/.local/bin/claude1`, `claude2`… | ccamへのリンク |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME`、`XDG_STATE_HOME`の各変数を設定すると、これらのフォルダーの場所が変わります。
