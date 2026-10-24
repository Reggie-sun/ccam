---
title: グループによる共有
description: ccamのグループは、選んだ項目を1つのフォルダーを通じてアカウント間で共有します。ccam groupは、計画、バックアップ、競合時の選択を伴ってデータを移動します。
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "e4f05b0fd572"
---

各アカウントは独自のフォルダーを持ちますが、その中身の多くは分けておく必要がありません。**グループ**は、何を誰と共有するかを決めます。グループのすべてのアカウントがそれらの項目を1つのフォルダー`~/.local/share/ccam/groups/<group>/`にリンクするので、全員が同じコピーを読み書きします。ログイントークンが共有されることはありません。

1つのアカウントは複数のグループに入れますが、各項目を受け取るのは1つのグループからだけです。典型的な構成は、設定と拡張機能のためのグループを1つと、2つのアカウントの間で会話履歴を共有する2つ目のグループです。

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

最初に追加したアカウントが`default`を作成し、その後のアカウントは、特に指定しない限りそこに加わります。

## 使い方

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

データを移動するコマンドはどれも、計画を表示し、グループのコピーと内容が異なるファイルごとに確認し、何かを変更する前に`Proceed? [y/N]`と尋ねます。

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

同じ変更は[ダッシュボード](/features/dashboard)の**Groups**タブでも行えます。アカウントの上で`s`を押すと、そのアカウントのグループを選べます。

## 項目

| 項目 | 対象 |
| --- | --- |
| `mcp` | MCPサーバー。グループの`mcp.json`を通じて共有します（[共有MCPサーバー](/features/mcp)） |
| `settings` | `settings.json` |
| `memory` | グローバルの`CLAUDE.md` |
| `extensions` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` |
| `history` | `projects/`と`history.jsonl`。これにより、別のアカウントで`--continue`できます |

アカウントのどのグループも共有していない項目は、アカウント自身のフォルダーに残ります。

## コマンドとオプション

| コマンド | 効果 |
| --- | --- |
| `ccam group`または`ccam group list` | グループを一覧表示します |
| `ccam group show <group>` | グループが共有する項目、アカウント、フォルダー、MCPサーバーを表示します |
| `ccam group add <group> --items <items> [--members <accounts>]` | グループを作成します。そのアカウントの項目はグループに移動します |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | グループのすべてのアカウントについて、グループが共有する項目を変更します |
| `ccam group join <group> <account>...` | アカウントをグループに入れます |
| `ccam group leave <group> <account>...` | アカウントをグループから外します。各アカウントは、MCPサーバーも含めてコピーを持ち続けます |
| `ccam group rename <group> <new-name>` | 設定の中だけで名前を変更します |
| `ccam group remove <group> [--leave]` | グループを削除します |

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--from <account>` | `add`、`set`: 最初にグループにコピーを入れるアカウント | デフォルトのアカウント、次に番号順 |
| `--replace` | `join`: 同じ項目を共有しているグループから先に抜けます | オフ |
| `--fresh` | 項目の共有をやめるアカウントは、コピーを持たずに空の状態から始めます（インポートしたフォルダーにはそれでもコピーが残ります） | オフ |
| `--leave` | `remove`: 先にグループのアカウントを外します | アカウントがあるときは必須 |
| `--prefer <group\|mine\|merge>` | 競合の決め方 | 確認します。`--yes`では各ファイルのデフォルト |
| `-y`、`--yes` | 確認なしで実行します。競合は`--prefer`に従います | オフ |
| `--dry-run`（グローバル） | 計画を表示し、何も変更しません | オフ |

アカウントは名前か番号で、項目はカンマ区切りのリストで指定します。`ccam add`は`--group <group>`（繰り返し指定できます）か`--no-group`を受け付けます。

## しくみ

関係するファイルは次のとおりです。

- `~/.config/ccam/config.toml`: `[[group]]`のエントリーです。`name`、`items`、`members`を持ち、フォルダーの名前が異なる場合は`dir`も持ちます。
- `~/.local/share/ccam/profiles/<name>/`: アカウントの`CLAUDE_CONFIG_DIR`です。
- `~/.local/share/ccam/groups/<folder>/`: 共有されるコピーです。

`config.toml`を読み込むと、各アカウントのグループが、項目ごとに1つのフォルダーに解決されます。5つのアカウントで約25 µsしかかからないので、起動時に読むのは引き続き`config.toml`だけです。あるアカウントが1つの項目を2つのグループから受け取ることになる設定はエラーとなり、修正するまで何も起動しません。

変更は設定のコピーに対して行われます。

1. コマンドがコピーを編集し（`AddGroup`、`SetGroupItems`、`JoinGroup` …）、その際に規則を確認します。
2. `internal/group`が変更前と変更後を比較し、項目が変わるアカウントごとに`share.Build`で計画を立てます。計画は1つずつ順に立てられ、それぞれがそれより前の計画でグループのフォルダーに入るものを考慮するので、コピーを移動するのは最初のアカウントだけです。最初のコピーは、グループのフォルダーにすでにあればそれを使い、なければ`--from`、デフォルトのアカウント、ほかのアカウントの番号順で決まります。
3. それらのアカウントのいずれかが、Claude Codeがどう起動されたかにかかわらず実行中の場合（``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``）や、開発者モードでないWindowsでファイルにリンクが必要になる場合、ccamは拒否します。
4. 確認か`--prefer`で競合を決め、計画を表示してから、ccamが確認を求めます。
5. 保存する直前に、ccamはもう一度確認します。質問の間に起動したアカウントがあれば変更は止まり、同じグループに対するほかのccamの変更があった場合も止まります（`the groups changed while this plan waited; nothing was moved, make the change again`）。何かを移動する前に新しい設定を保存し、変更を`~/.local/state/ccam/group-change.json`に記録します。手順が途中で止まった場合、その記録が残ります。`ccam doctor --fix`が、あなたの回答のとおりに、MCPの手順も含めて元の順序で変更を完了させます。それまでは、ほかのグループの変更は始まりません（``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``）。
6. まずMCPの手順を実行します。アカウントがグループから拡張機能を外す場合、グループのプラグインの一覧で、パスがアカウントのフォルダーを経由しているものは、グループのフォルダーを指すように変更します（計画では`group plugin lists:`）。その後、計画を順に実行します。どの手順もデータをグループかバックアップフォルダーに移動します。グループがバイト単位で同じものをすでに持っているファイルだけは、バックアップせずに捨てます。

`share.Layout`は、各項目が対象とするパスの許可リストです。ここにないものがグループに入ることはありません。`share.Never`は`.credentials.json`と`.claude.json`を指定しており、どちらかがレイアウトに現れた場合、`Build`は内部エラーで停止します。`mcp`項目にはパスがありません。サーバーの一覧は`.claude.json`の中にあるため、[MCPの共有](/features/mcp)はグループの`mcp.json`を通じて行われます。

次の図は、`side`が履歴の共有を始め、グループにすでに内容がある場合に、1つのフォルダー`projects/`がどうなるかを示しています。

<!-- motion:share-on -->

グループのそのパスにまだ何もない場合、そのアカウントが最初の共有者となり、フォルダーはそのままグループに移動します。そうでない場合は、グループにない各ファイルをグループに移動し、両方にあるファイルはアカウントのフォルダーに残すので、グループのバージョンが保たれます。グループのものと同じファイルは捨て、内容の異なるファイルは計画に名前が示されます（`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`）。その後、アカウントのフォルダーはアカウントのフォルダー内の`.ccam-backup/<time>/`に移動し、その場所にリンクが置かれます。グループの項目そのものが、加わるアカウント自身の項目へのリンクになっている場合は、アカウントの内容がその代わりにグループに入ります。会話ファイルの名前はセッションIDに基づいているため、複数のアカウントの履歴を競合なく結合できます。`plugins/`のレジストリは、アカウントのフォルダーを経由する絶対パスを保存しています。ccamはそれらをグループのフォルダーを指すように変更し、グループにすでにレジストリがある場合は、アカウントの項目をそこにマージします。

グループから抜けると、リンクが逆の処理で戻されます。

<!-- motion:share-off -->

ccamはまず、グループの現在の内容をアカウントのフォルダー内の`.ccam-backup/incomplete-<name>`にコピーし、コピーが完了してから初めてリンクを置き換えます。そのため、コピーが途中で止まってもリンクはそのまま残り、`ccam doctor --fix`がもう一度コピーします（そのような完了処理の計画には`restore <entry>: … (copied back from the group before a stop)`と表示されます）。リンクを削除するのは、計画を立てたときと同じようにリンクがまだgroupsフォルダーの中を指している場合だけです。リンクがあった場所に実際のファイルがあるなど、その後変わったパスはそのまま残し、変更は止まります。`--fresh`を付けるとコピーをスキップします。グループの内容は、ほかのアカウントのためにそのまま残ります。フォルダーがまだ空のグループに項目を移すアカウントは、それまで見ていた内容のコピーを持ち込むので、ファイルを失いません。そのコピーは、完了するまで`<group>/.ccam-incomplete-<name>`に置かれます。MCPの共有をやめたアカウントは、グループのサーバーを持ち続けます。それらはアカウント自身の`.claude.json`に入り（バックアップは残ります）、同じ名前のアカウント自身のサーバーもそのまま残ります。`--fresh`を付けた場合は、代わりに読み込まれなくなるサーバーが計画に示されます。もう一度加わると、誰も変更していないコピーは削除されます。`extensions`から抜けるときも、グループのプラグインの一覧で、パスが抜けるアカウントのフォルダーを経由していたものはグループのフォルダーを指すように変更します。ほかのアカウントにとって意味は変わりません。

両方にあり内容が異なる1つのファイルは、競合になります。

<!-- motion:share-conflict -->

`settings.json`のような単一のファイルでは、あなたが選びます。フォルダーの中ではグループのコピーが残り、計画に内容の異なるファイルの名前が示されます。`history.jsonl`では、`merge by line`、`keep group`、`use mine`の順で、マージが先頭です。それ以外の`.json`ファイルでは`keep group`、`use mine`、`merge JSON`です。それ以外のファイルでは`keep group`と`use mine`です。`merge JSON`は、オブジェクトを再帰的にマージし、配列を重複なしで結合します。両方が同じキーに異なる値を設定している場合はグループの値を残すので、ほかのアカウントの動作は変わりません。`keep group`以外のすべての選択肢では、まずグループのコピーをアカウントの`.ccam-backup/<time>/group/`に保存します。自分のファイルは常にバックアップに移動します。

起動のたびに、`share.Check`がレイアウトの各パスをアカウントのグループと比較します。使うのは`lstat`と`readlink`だけです。足りないリンクはその場で作成します。リンクがあるべき場所に通常のファイルがある場合、リンクが間違った場所を指している場合、古いリンクが残っている場合は、起動を停止して`ccam doctor --fix`を案内します。

### 名前の変更、削除、古い形式

- `ccam group rename`は設定だけを変更します。フォルダーの名前はそのままで、グループはそれを`dir`として記録するので、何も移動せず、何も失われません。
- `ccam group remove`は、`--leave`を付けない限り、アカウントがあるグループの削除を拒否します。付けた場合は各アカウントがコピーを持ち続け、フォルダーは、ほかのグループが使っていなければ`~/.local/share/ccam/trash/`に移動します。
- グループが導入される前、`config.toml`はアカウントごとに`[profile.share]`テーブルを保存していました。ccamはその形式を、何も移動せずにグループとして読み込みます。最初に保存するときに新しい形式で書き込み、古いファイルを`config.toml.bak`として残します。同じアカウントが共有している項目は、フォルダーの名前を付けた1つのグループになります。残りの項目には`default-history`のような名前が付き、同じフォルダーを使い続けます。

### Windowsの場合

Windowsがシンボリックリンクを作成するのは、開発者モード（設定 > システム > 開発者向け）のときか、管理者の場合だけです。その権限がない場合、共有フォルダーはディレクトリジャンクションになり、Claude Codeからは同じように動作します。共有ファイル（`settings.json`、`CLAUDE.md`、`history.jsonl`）はシンボリックリンクにしかできないので、それが必要になるグループの変更は、何かを移動する前に拒否されます。同じ理由で、そこでの最初の`ccam add`は`mcp`と`extensions`だけを含むグループ`default`を作ってその旨を知らせ、設定、メモリー、履歴を共有するグループへの`ccam add --group`は拒否されます。

## 安全性と特殊なケース

- **実行中のアカウント。** 移動の対象となるアカウントのいずれかとしてClaude Codeが実行中の間は、変更を拒否します。ccamが起動したかどうかにかかわらず（通常の`claude`、IDE、バックグラウンドセッション）、セッションがそれらのファイルに書き込み続けるためです。この確認は保存する直前にもう一度行います。グループから抜けるときや加わるときは、指定したアカウントだけが対象です。`extensions`から抜けるときはグループのプラグインの一覧のパスも書き換えますが、ほかのアカウントにとっての意味は変わりません。`--dry-run`では計画が表示されます。
- **バックアップ。** 計画ごとに、アカウントのフォルダー内に専用の`.ccam-backup/<time>/`フォルダーが作られます。ccamがこれらを削除することはありません。
- **自分で作ったリンク。** パスがすでに自分で作ったリンク（たとえばdotfilesのリポジトリへのリンク）で、グループが空の場合、ccamはリンク先の内容をグループにコピーし、リンク自体をバックアップに移動します。リンク先には触れません。
- **インポートしたフォルダー。** `ccam import`でその場所のまま登録したアカウント（`~/.claude`など）はその場所で変更されるので、通常の`claude`からも変更が見えます。計画にはその旨が示されます。グループから抜けるときは、`--fresh`を付けても常にコピーが残ります。`ccam remove`も、アカウントを登録から外す前に、この方法でグループから抜けさせます。グループに入っている間、共有ファイルはグループのフォルダーにあるので、`~/.local/share/ccam`を削除する前にアカウントを削除してください（またはグループから抜けてください）。
- **同じフォルダーで2つのアカウントを使う場合。** 履歴を共有するアカウントは同じファイルに書き込みます。すでに1つのアカウントが実行中のフォルダーで2つ目のアカウントを起動すると、ccamは警告を表示し、確認を求めます。
- **手作業での編集。** アカウントではないメンバーは無視され、`ccam doctor`が報告します。`local`、`on`、`off`、`none`はグループ名に使えません。

## 関連項目

- [ダッシュボード](/features/dashboard)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [共有MCPサーバー](/features/mcp)
- [既存の環境のコピー](/features/migrate)
- [すべてのコマンド](/docs/commands)
