---
title: アカウント間の共有
description: ccam shareが、計画、バックアップ、競合時の選択を伴って、設定、CLAUDE.md、拡張機能、履歴をグループのフォルダーにリンクするしくみと、offでコピーを戻すしくみです。
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

各アカウントは独自のフォルダーを持ちますが、その中身の多くは分けておく必要がありません。`ccam share`は、1つのアカウントの1つの項目ごとに共有を有効または無効にします。アカウントのファイルやフォルダーはグループのフォルダーに移動し、シンボリックリンクに置き換えられるので、グループのすべてのアカウントが同じコピーを読み書きします。`ccam add`の後にアカウントが共有するものを変更したり、いくつかのアカウントだけで独自のグループを作ったりするときに使います。

## 使い方

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

項目のオプションを付けない場合、`ccam share`は現在の状態を表示します。

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

オプションを付けると、ccamは計画を表示し、グループのコピーと内容が異なるファイルごとに確認し、何かを変更する前に`Proceed? [y/N]`と尋ねます。新しいアカウントは会話履歴以外のすべてを共有します。各項目をユーザーの視点から説明したものは、[MCPサーバー、設定、履歴の共有](/docs/sharing)を参照してください。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCPサーバー | 変更なし |
| `--settings <on\|off\|group>` | `settings.json` | 変更なし |
| `--memory <on\|off\|group>` | グローバルの`CLAUDE.md` | 変更なし |
| `--extensions <on\|off\|group>` | `agents/`、`commands/`、`skills/`、`plugins/`、`output-styles/`、`rules/` | 変更なし |
| `--history <on\|off\|group>` | `projects/`と`history.jsonl` | 変更なし |
| `--fresh` | 共有を無効にするとき、グループの内容をコピーせずに空の状態から始めます | オフ |
| `--prefer <group\|mine\|merge>` | 競合の決め方 | 確認します。`--yes`では各項目のデフォルト |
| `-y`、`--yes` | 確認なしで実行します。競合は`--prefer`に従います | オフ |
| `-h`、`--help` | `share`のヘルプ | |
| `--dry-run`（グローバル） | 対応しています。競合について確認し、計画を表示しますが、何も変更しません | オフ |

値の`on`は`default`グループ、`off`はローカル（共有しない）を意味し、それ以外の値はグループ名です。グループ名は`~/.local/share/ccam/groups/`の下のフォルダー名になるため、`ccam share side --history ../x`は`invalid group name: "../x"`で失敗します。アカウントは名前か番号で指定し、変更されるのは指定した項目だけです。

## しくみ

関係するファイルは次のとおりです。

- `~/.config/ccam/config.toml`: 各アカウントの`[profile.share]`テーブルで、項目ごとに1つの値を持ちます。
- `~/.local/share/ccam/profiles/<name>/`: アカウントの`CLAUDE_CONFIG_DIR`です。
- `~/.local/share/ccam/groups/<group>/`: 共有されるコピーです。

1. `internal/cli/share.go`が設定を読み込み、`config.Share.Set`で各`--<item>`の値をアカウントの設定のコピーに適用します。アカウントが実行中の場合は拒否します。
2. `share.Build`が現在の設定と目標の設定を比較し、手順を列挙します。ディスクは読み取るだけで、`share.Layout`のすべてのパスについて`os.Lstat`と`os.Readlink`を使います。計画は操作の一覧です。グループの空の場所への移動（adopt）、グループにないファイルの移動（merge dir）、解決した競合のグループへの書き込み、バックアップ、リンクの削除、リンクの作成、コピーの戻し、プラグインのパスの書き換え、プラグインのレジストリのマージがあります。
3. 確認か`--prefer`で、競合を決めます。
4. `Plan.Describe`が計画を表示し、ccamが確認を求めます。
5. 何かを移動する前に、目標の設定を`config.toml`に書き込みます。`Plan.Apply`が途中で止まった場合、`ccam doctor --fix`や再実行は、保存された目標に向けて計画を立て、変更を完了します。
6. `Plan.Apply`が手順を順に実行します。どの手順もデータをグループかバックアップフォルダーに移動するだけで、何も削除しません。

`share.Layout`は、各項目が対象とするパスの許可リストです。ここにないものがグループに入ることはありません。`share.Never`は`.credentials.json`と`.claude.json`を指定しており、どちらかがレイアウトに現れた場合、`Build`は内部エラーで停止します。`mcp`項目にはパスがありません。サーバーの一覧は`.claude.json`の中にあるため、[MCPの共有](/features/mcp)は代わりにグループの`mcp.json`を通じて行われ、`--mcp`は設定を変更するだけです。

次の図は、`side`が履歴の共有を有効にし、グループにすでに内容がある場合に、1つのフォルダー`projects/`がどうなるかを示しています。

<!-- motion:share-on -->

グループのそのパスにまだ何もない場合、そのアカウントが最初の共有者となり、フォルダーはバックアップなしでそのままグループに移動します。そうでない場合は、`mergeDir`がグループにない各ファイルをグループに移動し、両方にあるファイルはアカウントのフォルダーに残すので、グループのバージョンが保たれます。その後、アカウントのフォルダーはアカウントのフォルダー内の`.ccam-backup/<time>/`に移動し、`os.Symlink`がその場所にリンクを置きます。会話ファイルの名前はセッションIDに基づいているため、複数のアカウントの履歴を競合なく結合できます。

移動には`rename(2)`を使い、ファイルシステムをまたぐ場合はコピーと削除にフォールバックします。移動したツリーの中にあり、ツリーの外を指していた相対リンクは絶対リンクに書き換えられるので（`fsutil.PortableLink`）、新しい場所からも解決できます。`plugins/`のレジストリ（`installed_plugins.json`、`known_marketplaces.json`）は、アカウントのフォルダーを経由する絶対パスを保存しています。ccamはそれらをグループのフォルダーを指すように変更し、グループにすでにレジストリがある場合は、アカウントの項目をそこにマージします。

項目の共有を無効にすると、リンクが逆の処理で戻されます。

<!-- motion:share-off -->

ccamは、リンクがgroupsフォルダーの中を指している場合だけリンクを削除し、`CopyTree`でグループの現在の内容をアカウントにコピーします。`--fresh`を付けるとコピーをスキップし、アカウントはそのパスに何もない状態から始めます。グループの内容は、ほかのアカウントのためにそのまま残ります。項目をあるグループから別のグループに移すと、リンクを削除して新しいグループにリンクし、古いグループはそのまま残します。計画にはその旨が`Note:`の行で示されます。

両方にあり内容が異なる1つのファイルは、競合になります。

<!-- motion:share-conflict -->

`share.NewConflict`は、そのファイルに合った選択肢を列挙します。`history.jsonl`では、`merge by line`、`keep group`、`use mine`の順で、マージが先頭です。それ以外の`.json`ファイルでは`keep group`、`use mine`、`merge JSON`です。それ以外のファイルでは`keep group`と`use mine`です。最初の選択肢がデフォルトです。

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON`（`MergeJSONBytes`）は、オブジェクトを再帰的にマージし、配列を重複なしで結合します。両方が同じキーに異なる値を設定している場合はグループの値を残すので、ほかのアカウントの動作は変わりません。`merge by line`はグループの行を残し、そこにない自分の行を追加します。`keep group`以外のすべての選択肢では、まずグループのコピーをアカウントのフォルダー内の`.ccam-backup/<time>/group/`に保存し、それから結果を`fsutil.AtomicWrite`で書き込みます。自分のファイルは常にバックアップに移動します。

`--prefer group`、`mine`、`merge`を指定すると、確認なしですべての競合を決めます。`merge`は、そのファイルが提供するほうの`merge JSON`か`merge by line`を意味し、どちらも提供しないファイルはデフォルトのままです。`--yes`だけを付けた場合は、各ファイルのデフォルトを使います。

起動のたびに、`share.Check`がレイアウトの各パスを設定と比較します。使うのは`lstat`と`readlink`だけなので、起動経路でも十分に軽い処理です。足りないリンクは、失うものがないのでその場で作成します。リンクがあるべき場所に通常のファイルがある場合、リンクが間違った場所を指している場合、古いリンクが残っている場合は、起動を停止します。

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix`は、保存された設定を現在と目標の両方に指定して`share.Build`を実行します。これはディスク上の状態からすべての手順を導き、競合ではグループの内容を残し、置き換えるものはすべてバックアップします。

## 安全性と特殊なケース

- **実行中のアカウント。** アカウントとしてClaude Codeが実行中の間は、ccamは共有の変更を拒否します（`account work is running; quit it first, then change its sharing`）。セッションが移動するファイルに書き込み続け、その間に追加された行が失われてしまうためです。`--dry-run`では計画が表示されます。
- **バックアップ。** 計画ごとに、アカウントのフォルダー内に専用の`.ccam-backup/<time>/`フォルダーが作られます。同じ秒に2つの計画が始まった場合は、番号の接尾辞が付きます。ccamがこれらを削除することはありません。
- **途中での失敗。** 手順が失敗した場合、ccamはエラーと元のファイルの場所を表示し、`ccam doctor --fix`で完了するよう案内します。目標はすでに保存されているので、修復はそこに向けて続きます。
- **自分で作ったリンク。** パスがすでに自分で作ったリンク（たとえばdotfilesのリポジトリへのリンク）で、グループが空の場合、ccamはリンク先の内容をグループにコピーし、リンク自体をバックアップに移動します。リンク先には触れません。
- **ログインファイル。** `.credentials.json`と`.claude.json`が、グループにリンク、移動、コピーされることはありません。
- **同じフォルダーで2つのアカウントを使う場合。** 履歴を共有するアカウントは同じファイルに書き込みます。同じ履歴グループのアカウントがすでに実行中のフォルダーで2つ目のアカウントを起動すると、ccamは警告を表示し、続行する前に確認を求めます。

## 関連項目

- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [共有MCPサーバー](/features/mcp)
- [既存の環境のコピー](/features/migrate)
- [すべてのコマンド](/docs/commands)
