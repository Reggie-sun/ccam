---
title: MCPサーバー、設定、履歴の共有
description: ccamがログインを分けたまま、MCPサーバー、settings.json、CLAUDE.md、エージェント、スキル、会話履歴をアカウント間で共有するしくみです。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

各アカウントはそれぞれ独自のログインを持ちます。それ以外のものは、項目ごとに共有できます。ccamはアカウントのファイルやフォルダーを共有グループのフォルダーにリンクするので、あるアカウントで行った変更が、そのグループのすべてのアカウントに反映されます。

## 共有できる項目

| 項目 | オプション | 共有されるもの |
| --- | --- | --- |
| MCPサーバー | `--mcp` | グループの`mcp.json`。起動時にClaude Codeに渡されます |
| 設定 | `--settings` | `settings.json`: 権限、フック、モデル、ステータスライン |
| グローバルメモリー | `--memory` | `CLAUDE.md` |
| 拡張機能 | `--extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| 会話履歴 | `--history` | `projects/`と`history.jsonl` |

新しいアカウントは、会話履歴以外のすべてを共有します。

ログインファイルである`.credentials.json`と`.claude.json`が共有されることはありません。ccamは決められたリストにあるパスだけをリンクし、この2つはそのリストに含まれていません。

## グループ

`on`は`default`グループを意味します。名前を指定して、一部のアカウントだけで項目を共有することもできます。

```sh
ccam share work --history team
ccam share side --history team
```

これで`work`と`side`は互いに会話を共有しますが、ほかのアカウントとは共有しません。

## 共有を有効にする

アカウントで項目の共有を開始すると、ccamは次の処理を行います。

1. 移動、コピー、リンクするすべてのファイルの計画を表示し、開始する前に確認を求めます。
2. グループがまだ空であれば、アカウントのファイルをグループのフォルダーに移動します。空でなければ、グループにまだないファイルを追加します。
3. 同じファイルの内容が異なる場合は、グループのバージョンを残すか、自分のバージョンを使うか、JSONをマージするかを尋ねます。
4. 元のファイルを、アカウントのフォルダー内の`.ccam-backup/<time>/`に保存します。
5. 元のファイルを、グループへのリンクに置き換えます。

会話ファイルの名前はセッションIDに基づいているため、複数のアカウントの履歴を競合なく結合できます。

## 共有を無効にする

```sh
ccam share work --history off
```

リンクが削除され、アカウントにはグループの現在の内容のコピーが作られます。代わりに空の状態から始めたい場合は、`--fresh`を追加します。グループの内容は、ほかのアカウントのためにそのまま残ります。

## 既存の環境をコピーする

ccamを使う前からClaude Codeを使っていた場合、設定は`~/.claude`にあります。`ccam migrate`はそれをグループにコピーし、その項目を共有するすべてのアカウントで使えるようにします。

```sh
ccam migrate
```

コピーされるもの:

- `settings.json`とグローバルの`CLAUDE.md`
- 拡張機能の`agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/`
- ユーザーのMCPサーバー（`~/.claude.json`の`mcpServers`項目）。グループの`mcp.json`にコピーされます
- 会話履歴（`projects/`と`history.jsonl`）。ただし`--history`を付けた場合だけです

コピー元のフォルダーは読み取るだけです。中のものを移動、リンク、変更することはなく、通常の`claude`もこれまでどおり動作します。ログインファイルの`.credentials.json`と`.claude.json`がコピーされることはありません。

ccamは、何かをコピーする前に計画を表示して確認を求めます。

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

最後に、ccamは各項目をそのグループで共有しているアカウントを一覧表示します。共有していないアカウントがあれば、そのアカウントで共有を有効にする`ccam share`コマンドを表示します。

グループにすでに内容の異なるファイルがある場合、ccamは共有を有効にするときと同じように、グループのものを残すか、自分のものを使うかを尋ねます。`settings.json`などのJSONファイルはマージすることもでき、`history.jsonl`は行単位でマージできます。グループの以前のファイルは、グループのフォルダー内の`.ccam-backup/<time>/`に保存されます。

| オプション | 効果 |
| --- | --- |
| `--settings`、`--memory`、`--extensions`、`--mcp`、`--history` | 指定した項目だけをコピーします。指定しない場合は、会話履歴以外のすべてをコピーします |
| `-g`、`--group <group>` | `default`の代わりに、このグループにコピーします |
| `-y`、`--yes` | 確認なしでコピーします。内容が異なるファイルは、グループのものを残します |

`ccam --dry-run migrate`のようにコマンドの前に`--dry-run`を付けると、何もコピーせずに計画だけを確認できます。別のフォルダーからコピーするには、`ccam migrate ~/.claude-work --group work`のようにフォルダーを指定します。

## MCPサーバー

Claude CodeはMCPサーバーをログインと一緒に`.claude.json`に保存するため、このファイルは共有できません。ccamは共有サーバーをグループの`mcp.json`に保存し、`--mcp-config=<group>/mcp.json`を付けてClaude Codeを起動します。自分で指定した`--mcp-config`オプションも、その上に追加で適用されます。

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントにだけ保存されます。次にそのアカウントを起動したとき、ccamはそのサーバー名と、サーバーをグループに移動するためのコマンドを表示します。

```sh
ccam mcp import work gh
```

OAuthでサインインするMCPサーバーは、トークンをアカウントのログインと一緒に保存します。そのため、各アカウントで一度ずつ認可してください。

サーバーを`claude mcp list`や`/mcp`に通常のユーザーサーバーとして表示したい場合は、`~/.config/ccam/config.toml`で`mcp_strategy = "sync"`を設定します。するとccamは、起動前にClaude Code自身のファイルロックのもとで、バックアップを取ったうえで、各アカウントの`.claude.json`にサーバーを書き込みます。

## 同じフォルダーで2つのアカウントを使う場合

会話履歴を共有するアカウントは、同じファイルに書き込みます。同じグループの別のアカウントがすでに実行中のフォルダーで、2つ目のClaude Codeを起動しようとすると、ccamは警告を表示し、続行するかどうかを尋ねます。

## 問題が起きたとき

```sh
ccam doctor
ccam doctor --fix
```

たとえばバックアップを手動で復元した後などに、共有パスが通常のファイルになってしまった場合、ccamはClaude Codeを起動する前に停止し、`ccam doctor --fix`を案内します。このコマンドは、そのファイルをグループにマージし、バックアップを残します。
