---
title: 点検と修復
description: ccam doctorで共有リンク、ファイルの権限、ログイン、Claude Codeのバージョンを点検し、ccam doctor --fixで安全に修復できるものを修復します。
updated: 2026-10-21
order: 9
translation: machine
sourceHash: "9be63350cb69"
---

`ccam doctor`はccamの環境を点検します。対象は、Claude Codeのバイナリ、各アカウントのフォルダー、ログイン、共有リンク、共有MCPサーバー、使用量エンドポイント、ショートカットコマンドです。`ccam doctor --fix`は、データを失わずに修復できるものを修復し、置き換えるファイルは先にバックアップフォルダーに移動します。`ccam <name>`が共有リンクの問題で停止したとき、ファイルを手動で復元した後、Claude Codeを更新した後に実行してください。

## 使い方

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

最初の行ですべてのアカウントに一度に記号が付くので、ほかを読む前にどこに問題があるかがわかります。その下にはアカウントごとのボックスがあり、どのアカウントにも属さない結果（Claude Code、環境変数、プラグイン、MCPサーバー、ショートカットコマンド、フォルダーのルール、グループ、設定）は`General`ボックスに入ります。ボックスのタイトルには、そのアカウントの最も重い記号と問題の数が表示されます。ボックス内の各行は点検した領域1つで、常に`Folder`、`Login`、`Share links`の順に並び、そのあとにほかの領域が続きます。記号は`✔`が正常、`ℹ`が情報、`⚠`が警告、`✘`がエラーです。問題のないフォルダーは`✔ Folder  OK`と表示されます。

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix`を付けると、修復された結果の末尾に`→ fixed`が付き、その結果は問題として数えられなくなります。

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

ターミナルでは、ボックスはその幅に合わせられ、長いメッセージは折り返され、ターミナル自身の色で表示されます（`NO_COLOR`で色をオフにできます）。パイプやファイルに送ると、色は付かず、メッセージが折り返されたり切られたりすることもないので、`ccam doctor | grep settings.json`で行全体が見つかります。このとき、すべてのボックスは最も長い行の幅になります。ダッシュボードのHealthタブにも同じボックスが表示されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--fix` | 安全に修復できるものを修復します（変更するファイルはバックアップされます） | オフ |
| `--dry-run`（グローバル） | `--fix`を付けても、すべての点検を実行するだけで何も修復しません | オフ |

`ccam doctor`は引数を取りません。`--dry-run`に対応しており、`ccam --dry-run doctor --fix`は`ccam doctor`と同じ結果を表示し、何も変更しません。

終了コードは、修復されていない`✘`の結果が1つ以上残っている場合は`1`、それ以外は`0`です。警告と情報で失敗することはありません。`--fix`なしで失敗した場合、ccamは空行と`` `ccam doctor --fix` repairs what it can. ``を表示して終わります。設定を読み込めない場合も、`ccam:`のエラーとともに`1`で終了します。

## しくみ

`internal/cli/doctor.go`が`~/.config/ccam/config.toml`を読み込み、`internal/doctor/doctor.go`の`doctor.Run`を呼び出します。これは決まった順序で点検を実行し、それぞれが属するアカウントと点検した領域を付けて結果を返します。`tui.HealthReport`がそれらをボックスにまとめます。`--fix`は`Options.Fix`として渡され、ドライランではfalseになります。

<!-- motion:doctor-checks -->

中央のブロックはアカウントごとに1回実行されます。各点検と、それぞれに対して`--fix`が行うことは次のとおりです。

1. **Claude Codeのバージョン**（`checkClaude`）。起動時と同じ方法で`claude`のバイナリを特定し、`claude --version`を実行します。バージョンが、ccamのフォルダー構成を検証したバージョン（`doctor.LayoutVersion`）である`2.1.302`で始まらない場合は、情報`ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`を追加します。バイナリがない場合や`--version`が失敗する場合はエラーです。修復はありません。
2. **環境変数の認証情報**（`checkAuthEnv`）。`CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`が設定されている場合に警告します。そのシェルから起動したすべてのアカウントがこれを引き継ぎ、Claude Codeはアカウント自身のログインよりこれを優先します。修復はありません。自分で設定を解除してください。
3. **フォルダー**（`checkProfile`）。アカウントのフォルダーがない場合はエラーで、そのアカウントの点検はそこで終わります。グループやその他のユーザーのビットが立っている権限は警告で、`--fix`は`chmod 0700`を実行します。Windowsには権限ビットがないため、この確認と後述の`0600`の確認はスキップされます。ユーザープロファイル内のフォルダーは、そのアクセス制御リストで保護されます。
4. **ログインファイル**（`checkCredentials`）。`.credentials.json`を`lstat`で確認します。各アカウントは独自のトークンを持つ必要があるため、シンボリックリンクはエラーで、修復されません。`0600`以外のモードは警告で、`--fix`は`chmod 0600`を実行します。続いてログインを読み取ります（macOSでは先にキーチェーンから）。ログインがない場合は警告で、`ccam login <name>`を案内します。リフレッシュトークンの期限切れは`login expired`エラーです。アクセストークンの期限切れは、次の起動時にClaude Codeが更新するため情報にとどまります。有効なログインでは、メールアドレスとプランが表示されます（`login valid (dev@company.com, Claude Team · Acme)`）。いずれも修復されませんが、`ccam login --expired`で、ログインが必要なすべてのアカウントにログインできます。
5. **共有リンク**（`share.Check`）。`share.Layout`のすべてのパスをアカウントのグループと比較します。使うのは`lstat`と`readlink`だけです。報告するのは4種類で、リンクがない、リンクが別の場所を指している、リンクがあるべき場所に通常のファイルがある、現在はローカルになった項目のリンクが残っている、です。リンクがないものは警告、それ以外はエラーです。新しいバージョンのccamで項目に追加されたパス（履歴の`file-history/`、`tasks/`、`paste-cache/`、設定の`keybindings.json`）は、アカウントがまだ自分のコピーを持っている間は警告にとどまります: `is not shared with the group yet (ccam doctor --fix moves it in)`。これらが原因で起動が止まることはありません。
6. **不明な項目**（`checkUnknownEntries`）。アカウントのフォルダーの最上位にある、ccamが知らない名前（`share.Known`）を一覧表示します。Claude Codeのフォルダー構成が変わった手がかりになります。情報のみです。
7. **プラグインのパス**（`checkPlugins`）。プラグインのレジストリ（`installed_plugins.json`、`known_marketplaces.json`）は絶対パスを保存します。各グループの`plugins/`で、特定のアカウントのフォルダーを経由するパスは、そのアカウントを削除または名前変更すると壊れます。`--fix`は`profile.RewritePluginPaths`でそれらをグループのフォルダーに書き換えます。レジストリファイルはアトミックに書き込まれます。
8. **MCPサーバー**（`checkMCP`）。各グループの`mcp.json`を1回ずつ読み込みます。読めないファイルはエラー、`mcp.Validate`に通らないサーバーは警告です。修復はありません。
9. **使用量の応答形式**（`checkUsageFormat`）。キャッシュを使わず、通常の3秒の制限時間で、応答があるまでアカウントを1つずつ使って使用量エンドポイントに問い合わせます。応答が解析できなくなっている場合は、`ccam usage`がステータスラインのキャッシュにフォールバックすることを警告します。ネットワークを使うのはこの点検だけです。
10. **ショートカット**（`checkShims`）。ショートカットを持つアカウントがある場合だけ実行します。`~/.local/bin`にリンクがない場合は警告で、`--fix`は`shim.Install`を呼び出します。これは足りないリンクを作成するだけで、ほかのファイルには触れません。binフォルダーが`PATH`にない場合は警告で、修復されません。
11. **フォルダー**（`checkFolders`）。[`ccam folder`](/features/folder)のエントリーがある場合だけ実行します。アカウントが存在しないエントリーはエラーで、フォルダーがなくなったエントリーは警告です。どちらも修復されません。エントリーを削除すると、そこで起動するアカウントが変わってしまうためです。
12. **グループと設定**（`checkGroups`）。アカウントではないグループのメンバー、名前がその後ccamのコマンド名になったアカウント（`ccam run <name>`で起動するか、名前を変更してください）、`config.toml`の不明な`theme`や`usage_source`は警告です。修復はありません。
13. **ログイン**（`checkIdentities`）。各フォルダーがどのアカウントでログインしているかを`.claude.json`から読み取ります。`ccam login`が記録したものとは別のアカウントでログインしているフォルダー（たとえばセッション内で`/login`を実行した後）は警告です。別のフォルダーと同じClaudeアカウントにログインしているフォルダーも警告です。その2つは使用量の上限を共有するためです。修復はありません。

### 共有リンクの修復

起動時には、足りないリンクはその場で作成されますが、データを失う可能性があるものについては停止し、ここを案内します。

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

共有リンクの問題に対して、`--fix`はアカウントの共有設定を現在の状態と目標の状態の両方に指定して`share.Build`を呼び出し、続いて`Plan.Apply`を呼び出します。`Build`はディスク上の状態からすべての手順を導き直します。これは`ccam group`による変更が使うのと同じ計画なので、修復はグループに加わるときと同じ規則に従います。

<!-- motion:doctor-repair -->

この例では`settings.json`が通常のファイルです。グループにはすでに内容の異なる`settings.json`があるため、競合はデフォルトの選択である`KeepGroup`になり、グループのコピーが残ります。アカウントのファイルはアカウントのフォルダー内の`.ccam-backup/<time>/settings.json`に移動し、`os.Symlink`がそのパスをグループにリンクします。ほかのケースは次のとおりです。

- グループにまだコピーがない場合は、アカウントのファイルがグループに移動し、グループの最初の内容になります。
- フォルダーの場合は、グループにないファイルがグループに移動し、残りはバックアップに移動します。
- `history.jsonl`の場合、デフォルトは`merge by line`です。グループの行の後に、自分の新しい行が続きます。グループの以前のコピーは`.ccam-backup/<time>/group/`に保存されます。
- 別のグループへのリンクは置き換えられます。ccamの外にあるフォルダーへのリンクはバックアップに移動し、グループが空であれば、その内容がグループの最初の内容になります。
- 現在はローカルになった項目のリンクは削除され、アカウントにはグループの内容のコピーが作られます。

リンクの削除を除き、すべての手順はデータをグループかバックアップフォルダーに移動します。何も削除されません。

## 安全性と特殊なケース

- 通常の`ccam doctor`は何も変更しません。スクリプトで確実にしたい場合は、`--fix`と一緒に`--dry-run`を使ってください。
- 実行中のアカウントでは、`doctor --fix`は足りないリンクを作るだけで、`` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``を付け加えます。修復で移動するファイルに、セッションが書き込み続けるためです。
- 共有の修復に失敗した場合、ccamはそのアカウントについて`could not repair share links: <error>`と報告し、不明な項目の点検をスキップします。ほかのアカウントの点検は続きます。
- `ccam group`による変更が途中で失敗した場合は、新しいグループはすでに保存されているので、ccamは`ccam doctor --fix`で完了するよう案内します。
- シンボリックリンクになった`.credentials.json`、期限切れのログイン、環境変数の認証情報は、報告されるだけで変更されることはありません。
- プラグインのレジストリの書き換えはアトミックに書き込まれますが、`.ccam-backup/`にはコピーされません。

## 関連項目

- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [既存の環境のコピー](/features/migrate)
- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [コマンド](/docs/commands)
