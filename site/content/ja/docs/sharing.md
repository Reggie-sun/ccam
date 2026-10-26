---
title: MCPサーバー、設定、履歴の共有
description: ccamのグループが、ログインを分けたまま、MCPサーバー、settings.json、CLAUDE.md、スキル、会話履歴をClaude Codeのアカウント間で共有するしくみです。
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 2d8029ad6ec8
---

各アカウントはそれぞれ独自のログインを持ちます。それ以外のものは、**グループ**を通じて共有できます。グループは共有する項目と所属するアカウントを決め、ccamはそれらのアカウントのファイルをグループのフォルダーにリンクするので、あるアカウントで行った変更がすべてのアカウントに反映されます。

## 共有できる項目

| 項目 | 共有されるもの |
| --- | --- |
| `mcp` | MCPサーバー: グループの`mcp.json`。起動時にClaude Codeに渡されます |
| `settings` | `settings.json`: 権限、フック、モデル、ステータスライン |
| `memory` | グローバルの`CLAUDE.md` |
| `extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| `history` | 会話履歴: `projects/`と`history.jsonl` |

最初のアカウントが、会話履歴以外のすべてを共有する`default`グループを作成し、その後のアカウントはそこに加わります。

ログインファイルである`.credentials.json`と`.claude.json`が共有されることはありません。ccamは決められたリストにあるパスだけをリンクし、この2つはそのリストに含まれていません。

このページでは日常的な作業を説明します。すべてのオプションは機能ページで説明しています。[グループによる共有](/features/share)、[共有MCPサーバー](/features/mcp)、[既存の環境のコピー](/features/migrate)を参照してください。

## グループ

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

1つのアカウントは複数のグループに入れますが、各項目を受け取るのは1つのグループからだけです。そのため、`work`と`side`の間だけで履歴を共有し、ほかのアカウントとは共有しないようにするには、上の例のように、履歴を`default`に含めず、2つのアカウント専用のグループを作ります。`ccam share work`は、アカウントが何を、どのグループを通じて共有しているかを表示します。[ダッシュボード](/features/dashboard)のGroupsタブでも同じことができます。

## アカウントが共有を始めるとき

アカウントが項目の共有を始めると（グループに加わったときや、グループに項目が追加されたとき）、ccamは次の処理を行います。

1. アカウントごとに、移動、コピー、リンクするすべてのファイルの計画を表示し、開始する前に確認を求めます。
2. グループがまだ空であれば、アカウントのファイルをグループのフォルダーに移動します。空でなければ、グループにまだないファイルを追加します。
3. 同じファイルの内容が異なる場合は、グループのバージョンを残すか、自分のバージョンを使うか、JSONをマージするかを尋ねます。
4. 元のファイルを、アカウントのフォルダー内の`.ccam-backup/<time>/`に保存します。
5. 元のファイルを、グループへのリンクに置き換えます。

会話ファイルの名前はセッションIDに基づいているため、複数のアカウントの履歴を競合なく結合できます。

そのアカウントとしてClaude Codeが実行中の間は、Claude Codeがどう起動されたかにかかわらず（バックグラウンドセッションも含みます）、ccamはそのアカウントを変更しません。移動するファイルに、セッションが書き込み続けるためです。先にセッションを終了してください。ccamは保存する直前にもう一度確認するので、ccamがあなたの回答を待っている間に始まったセッションも見逃しません。

## アカウントが共有をやめるとき

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

リンクが削除され、アカウントにはグループの現在の内容のコピーが作られます。グループのMCPサーバーは、そのアカウント自身のサーバーと並んで、アカウントの`.claude.json`に入ります。代わりに空の状態から始めたい場合は、`--fresh`を追加します。その場合、計画には読み込まれなくなるMCPサーバーが表示されます。`~/.claude`のようなインポートしたフォルダーには、常にコピーが残ります。グループの内容は、ほかのアカウントのためにそのまま残ります。

## 既存の環境をコピーする

ccamを使う前からClaude Codeを使っていた場合、設定は`~/.claude`にあります。`ccam migrate`はそれをグループにコピーし、その項目を共有するすべてのアカウントで使えるようにします。

```sh
ccam migrate
```

ターミナルでは、ccamは最初に項目のチェックリストを表示します。会話履歴以外のすべてにチェックが入っています。Spaceキーで項目を切り替え、Enterキーで計画を表示します。フォルダーに該当するものがない項目には`not found`と表示されます。

コピーされるもの:

- `settings.json`とグローバルの`CLAUDE.md`
- 拡張機能の`agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/`
- ユーザーのMCPサーバー（`~/.claude.json`の`mcpServers`項目）。グループの`mcp.json`にコピーされます
- 会話履歴（`projects/`と`history.jsonl`）。ただし`--history`を付けた場合だけです

コピー元のフォルダーは読み取るだけです。中のものを移動、リンク、変更することはなく、通常の`claude`もこれまでどおり動作します。ログインファイルの`.credentials.json`と`.claude.json`がグループにコピーされることはありません。

OAuthでサインインするMCPサーバーは、トークンをサーバーの一覧ではなくログインに保存します。ccamはそのサインイン情報をコピーしません。コピーしたサインインは同じリフレッシュトークンを持つので、最初の更新で、`~/.claude`を含むほかのすべての持ち主がそのサーバーからログアウトされるためです。計画と結果のまとめには、`/mcp`で自分でそれらのサーバーにサインインする必要があるアカウントが示されます。

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
Group default: personal
To share them:
  ccam group join default work
```

最後に、ccamはグループのアカウントを一覧表示します。グループがまだ共有していない項目がある場合や、グループに入っていないアカウントがある場合は、それを変更する`ccam group`コマンドを表示します。

グループにすでに内容の異なるファイルがある場合、ccamは共有を有効にするときと同じように、グループのものを残すか、自分のものを使うかを尋ねます。`settings.json`などのJSONファイルはマージすることもでき、`history.jsonl`は行単位でマージできます。グループの以前のファイルは、グループのフォルダー内の`.ccam-backup/<time>/`に保存されます。

| オプション | 効果 |
| --- | --- |
| `--settings`、`--memory`、`--extensions`、`--mcp`、`--history` | 指定した項目だけを、チェックリストなしでコピーします。項目のオプションを指定せずにターミナル以外で実行した場合、ccamは会話履歴以外のすべてをコピーします。 |
| `-g`、`--group <group>` | `default`の代わりに、このグループにコピーします |
| `-y`、`--yes` | 確認なしでコピーします。内容が異なるファイルは、グループのものを残します |

`ccam --dry-run migrate`のようにコマンドの前に`--dry-run`を付けると、何もコピーせずに計画だけを確認できます。別のフォルダーからコピーするには、`ccam migrate ~/.claude-work --group work`のようにフォルダーを指定します。

## MCPサーバー

Claude CodeはMCPサーバーをログインと一緒に`.claude.json`に保存するため、このファイルは共有できません。ccamは共有サーバーをグループの`mcp.json`に保存し、`--mcp-config=<group>/mcp.json`を付けてClaude Codeを起動します。自分で指定した`--mcp-config`オプションも、その上に追加で適用されます。

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントにだけ、しかもそのときいたフォルダーにだけ保存されます（Claude Codeのlocalスコープです。`-s user`を付けると、どのフォルダーでもそのアカウントで使えるようになります）。次にそのアカウントを起動したとき、ccamはそのサーバー名とそのフォルダー、そしてサーバーをグループに移動するためのコマンドを表示します。

```sh
ccam mcp import work gh
```

グループにすでに同じ名前の別のサーバーがある場合、ccamはそれをスキップし、その旨を表示します。`--force`を付けるとグループのサーバーを置き換え、以前の`mcp.json`をグループの`.ccam-backup/`フォルダーに保存します。

OAuthでサインインするMCPサーバーは、トークンをアカウントのログインと一緒に保存します。そのため、各アカウントで一度ずつ`/mcp`でサインインしてください。`ccam ext`と`ccam doctor`は、ほかのアカウントはサインインしているのに、そのアカウントはサインインしていないサーバーを示します。

サーバーを`claude mcp list`や`/mcp`に通常のユーザーサーバーとして表示したい場合は、`~/.config/ccam/config.toml`の先頭、最初の`[[group]]`、`[[profile]]`、`[[folder]]`より上に`mcp_strategy = "sync"`を設定します。するとccamは、起動前にClaude Code自身のファイルロックのもとで、バックアップを取ったうえで、各アカウントの`.claude.json`にサーバーを書き込みます。グループのサーバーと同じ名前のサーバーをアカウントが独自に持っている場合、ccamはアカウントのサーバーを残し、その旨を知らせます。

## 同じフォルダーで2つのアカウントを使う場合

会話履歴を共有するアカウントは、同じファイルに書き込みます。同じグループの別のアカウントがすでに実行中のフォルダーで、2つ目のClaude Codeセッションを起動しようとすると、ccamは警告を表示し、続行するかどうかを尋ねます。

## 問題が起きたとき

```sh
ccam doctor
ccam doctor --fix
```

たとえばバックアップを手動で復元した後などに、共有パスが通常のファイルになってしまった場合、ccamはClaude Codeを起動する前に停止し、`ccam doctor --fix`を案内します。この修復は、そのパスをもう一度グループにリンクします。グループにコピーがなければ、そのファイルがグループのコピーになります。グループにコピーがあれば、グループは自分のコピーを残し、あなたのファイルはバックアップに移り、`--fix`がその場所を示します。途中で止まったグループの変更も、`ccam doctor --fix`が完了させます。
