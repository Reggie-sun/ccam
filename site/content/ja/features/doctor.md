---
title: 点検と修復
description: ccam doctorで共有リンク、ファイルの権限、ログイン、Claude Codeのバージョンを点検し、ccam doctor --fixで安全に修復できるものを修復します。
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor`はccamの環境を点検します。対象は、Claude Codeのバイナリ、各アカウントのフォルダー、ログイン、共有リンク、共有MCPサーバー、使用量エンドポイント、ショートカットコマンドです。`ccam doctor --fix`は、データを失わずに修復できるものを修復し、置き換えるファイルは先にバックアップフォルダーに移動します。`ccam <name>`が共有リンクの問題で停止したとき、ファイルを手動で復元した後、Claude Codeを更新した後に実行してください。

## 使い方

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

最初の行ですべてのアカウントに一度に記号が付くので、ほかを読む前にどこに問題があるかがわかります。その下にはアカウントごとのボックスがあり、どのアカウントにも属さない結果（Claude Code、環境変数、プラグイン、MCPサーバー、ショートカットコマンド、フォルダーのルール、グループ、設定、`auto_switch`）は`General`ボックスに入ります。ボックスのタイトルには、そのアカウントの最も重い記号と問題の数が表示されます。ボックス内の各行は点検した領域1つで、常に`Folder`、`Login`、`Share links`の順に並び、そのあとにほかの領域が続きます。記号は`✔`が正常、`ℹ`が情報、`⚠`が警告、`✘`がエラーです。問題のないフォルダーは`✔ Folder  OK`と表示されます。

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix`を付けると、修復された結果の末尾に`→ fixed`が付き、その結果は問題として数えられなくなります。

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

ターミナルでは、ボックスはその幅に合わせられ、長いメッセージは折り返され、ターミナル自身の色で表示されます（`NO_COLOR`で色をオフにできます）。ボックスが29列より狭くなることはないので、幅の狭いターミナルでも単語が途中で切れることはありません。パイプやファイルに送ると、色は付かず、メッセージが折り返されたり切られたりすることもないので、`ccam doctor | grep settings.json`で行全体が見つかります。このとき、すべてのボックスは最も長い行の幅になります。ダッシュボードのHealthタブにも同じボックスが表示されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--fix` | 安全に修復できるものを修復します（変更するファイルはバックアップされます） | オフ |
| `--dry-run`（グローバル） | `--fix`を付けても、すべての点検を実行するだけで何も修復しません | オフ |

`ccam doctor`は引数を取りません。`--dry-run`に対応しており、`ccam --dry-run doctor --fix`は`ccam doctor`の結果を表示し、`--fix`が退避するものを付け加え、何も変更しません。

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

終了コードは、修復されていない`✘`の結果が1つ以上残っている場合は`1`、それ以外は`0`です。警告と情報で失敗することはありません。`--fix`なしで失敗した場合、ccamは空行と`` `ccam doctor --fix` repairs what it can. ``を表示して終わります。設定を読み込めない場合も、`ccam:`のエラーとともに`1`で終了します。

## しくみ

`internal/cli/doctor.go`が`~/.config/ccam/config.toml`を読み込み、`internal/doctor/doctor.go`の`doctor.Run`を呼び出します。これは決まった順序で点検を実行し、それぞれが属するアカウントと点検した領域を付けて結果を返します。`tui.HealthReport`がそれらをボックスにまとめます。`--fix`は`Options.Fix`として渡され、ドライランではfalseになります。

<!-- motion:doctor-checks -->

中央のブロックはアカウントごとに1回実行されます。各点検と、それぞれに対して`--fix`が行うことは次のとおりです。

1. **Claude Codeのバージョン**（`checkClaude`）。起動時と同じ方法で`claude`のバイナリを特定し、`claude --version`を実行します。バイナリがない場合、`claude_bin`が実行可能ファイルでない場合、`--version`が失敗する場合はエラーです。修復はありません。
2. **環境変数の認証情報**（`checkAuthEnv`）。`CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`が設定されている場合に警告します。そのシェルから起動したすべてのアカウントがこれを引き継ぎ、Claude Codeはアカウント自身のログインよりこれを優先します。修復はありません。自分で設定を解除してください。
3. **フォルダー**（`checkProfile`）。アカウントのフォルダーがない場合はエラーで、そのアカウントの点検はそこで終わります。グループやその他のユーザーのビットが立っている権限は警告で、`--fix`は`chmod 0700`を実行します。インポートした`~/.claude`に、`~/.claude.json`とは別に独自の`.claude.json`がある場合は情報です。その場合、`ccam main`と通常の`claude`はログインの詳細とMCPサーバーを別々に持ちます。Windowsには権限ビットがないため、この確認と後述の`0600`の確認はスキップされます。ユーザープロファイル内のフォルダーは、そのアクセス制御リストで保護されます。
4. **ログインファイル**（`checkCredentials`）。`.credentials.json`を`lstat`で確認します。各アカウントは独自のトークンを持つ必要があるため、シンボリックリンクはエラーで、修復されません。`0600`以外のモードは警告で、`--fix`は`chmod 0600`を実行します。続いてログインを読み取ります（macOSでは先にキーチェーンから）。ログインがない場合や、ログインファイルにログインが入っていない場合は警告で、`ccam login <name>`を案内します。リフレッシュトークンの期限切れや、サーバーに拒否された後にClaude Codeが消去したログインは、エラー`` login expired → `ccam login <name>` ``です。アクセストークンの期限切れは、次の起動時にClaude Codeが更新するため情報にとどまります。Windows資格情報マネージャーに保存されたログインは、ccamが読み取れないため情報です。有効なログインでは、フォルダーが現在ログインしているメールアドレスとプランが表示されます（`login valid (dev@company.com, Claude Team · Acme)`）。いずれも修復されませんが、`ccam login --expired`で、ログインが必要なすべてのアカウントにログインできます。
5. **共有リンク**（`share.Check`）。`share.Layout`のすべてのパスをアカウントのグループと比較します。使うのは`lstat`と`readlink`だけです。報告するのは5種類で、リンクがない、リンクが別の場所を指している、リンクがあるべき場所に通常のファイルがある、現在はローカルになった項目のリンクが残っている、そしてリンクは正しいがリンク先のグループのフォルダーがない（`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`。起動時にはこれを確認しません）、です。リンクがないものとグループのフォルダーがないものは警告、それ以外はエラーです。グループからのコピーの戻しが途中で止まったものも報告します。`--fix`は準備済みの完全なコピーを所定の場所に置きますが、その後（起動によって）パスが新しく作られていた場合は、完全なコピーのあるフォルダーを示し、整理をあなたに任せます。新しいバージョンのccamで項目に追加されたパス（履歴の`file-history/`、`tasks/`、`paste-cache/`、設定の`keybindings.json`）は、アカウントがまだ自分のコピーを持っている間は警告にとどまります: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`。これらが原因で起動が止まることはありません。
6. **不明な項目**（`checkUnknownEntries`）。アカウントのフォルダーの最上位にある、ccamが知らない名前（`share.Known`）を一覧表示します。Claude Codeのフォルダー構成が変わった手がかりになります。情報のみです。
7. **プラグインのパス**（`checkPlugins`）。プラグインのレジストリ（`installed_plugins.json`、`known_marketplaces.json`）は絶対パスを保存します。各グループの`plugins/`で、特定のアカウントのフォルダーを経由するパスは、そのアカウントを削除または名前変更すると壊れます。`--fix`は`profile.RewritePluginPaths`でそれらをグループのフォルダーに書き換えます。レジストリファイルはアトミックに書き込まれます。
8. **MCPサーバー**（`checkMCP`）。各グループの`mcp.json`を1回ずつ読み込みます。読めないファイルはエラー、`mcp.Validate`に通らないサーバーは警告です。通常の`claude`の`~/.claude.json`にあり、どのアカウントも読み込まないサーバーは警告で、それらをコピーする`ccam migrate --mcp`のコマンドを示します（複数のグループがMCPを共有している場合は`--group <a|b>`付き）。グループの1つのアカウントだけが持つサーバーは情報で、それを共有する`ccam mcp import`のコマンドを示します。また、ほかのアカウントがサインインしているMCPサーバーにサインインしていないアカウントには、`/mcp`でサインインするよう警告します（ccamがサインインをコピーすることはありません）。修復はありません。
9. **使用量の応答形式**（`checkUsageFormat`）。キャッシュを使わず、応答があるまでアカウントを1つずつ使って使用量エンドポイントに問い合わせます。全体で4秒の制限時間です。`usage_source = "claude"`の場合は、代わりに1つのアカウントについてClaude Codeに問い合わせます。応答が解析できなくなっている場合は、`ccam usage`がステータスラインのキャッシュにフォールバックすることを警告します。どのアカウントも応答しなかった場合は、情報としてその旨を示します（`could not check the usage response format: …`）。ネットワークを使うのはこの点検だけです。
10. **ショートカット**（`checkShims`）。ショートカットを持つアカウントがある場合だけ実行します。`~/.local/bin`にリンクがない場合は警告で、ccamが作ったショートカットが、置き換えられたり移動されたりした古いccamをまだ実行している場合も警告です（`shortcuts in <bin> that run another ccam (one replaced or moved): …`）。`--fix`は`shim.Install`を呼び出します。これはccamが作って`~/.local/state/ccam/shortcuts.json`に記録したリンクだけを作成・更新し、ほかのファイルには触れません。警告の末尾に`→ fixed`が付くのは、それらのリンクがすべて作られた場合だけで、作れなかったものがあれば理由が追加されます。まだインストールされている別のccamを実行するショートカットは動作するので情報で、`--fix`はそれをこのccamに向け直します。binフォルダーが`PATH`にない場合は警告で、修復されません。
11. **フォルダー**（`checkFolders`）。[`ccam folder`](/features/folder)のエントリーがある場合だけ実行します。アカウントが存在しないエントリーはエラーで、フォルダーがなくなったエントリーは警告です。どちらも修復されません。エントリーを削除すると、そこで起動するアカウントが変わってしまうためです。
12. **グループと設定**（`checkGroups`）。アカウントではないグループのメンバー、名前がその後ccamのコマンド名になったアカウント（`ccam run <name>`で起動するか、名前を変更してください）、不明な`theme`や`usage_source`、そしてccamが読まない`config.toml`の行は警告です。最後のものはたいてい、`[[group]]`、`[[profile]]`、`[[folder]]`のテーブルの下に書かれ、そのテーブルに属してしまった設定です（`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`）。設定を変更するコマンドは、それが直るまで拒否します。修復はありません。
13. **ログイン**（`checkIdentities`）。各フォルダーがどのアカウントでログインしているかを`.claude.json`から読み取ります。`ccam login`が記録したものとは別のアカウントでログインしているフォルダー（たとえばセッション内で`/login`を実行した後）は警告です。インポートしたアカウントや古いccamが作ったアカウントのように、ログインの記録がないアカウントは、記録されたメールアドレスで比較します。別のフォルダーと同じClaudeアカウントにログインしているフォルダーも警告です。その2つは使用量の上限を共有するためです。別のフォルダーと同じログイントークンを持つフォルダー、つまりコピーしたフォルダー（通常の`claude`の`~/.claude`も含みます）はエラーです。先に更新した方がもう一方をログアウトさせるためで、`ccam login <name>`を実行して独自のログインを持たせてください。修復はありません。
14. **`auto_switch`**（`checkAutoSwitch`）。`auto_switch = true`のときだけ行います。`General`ボックスの情報が`auto_switch`の働きを説明し、[その説明](/features/next#auto-switch)にリンクします。アカウントの中にシートのあるプラン（Claude Team、Claude Enterprise）があれば、プランごとに1行で、該当するアカウントを挙げ、`auto_switch`が1つの組織のシート間でセッションを移すことを伝え、Anthropicのシートに関する規則をヘルプセンターのリンクとともに引用し、選択をあなたに委ねます。2.1.214より古いClaude Codeは、フックに自身のプロセス（`CLAUDE_PID`）を伝えないため、どのセッションも移せないので、そこで警告になります。ccamプラグインがない、またはプラグインがオフになっているアカウントには、そのアカウント自身のボックスに警告が出て、Claude Codeでそれを直すコマンドが示されます: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam`（または`/plugin enable ccam@ccam`）。APIキーのアカウントは対象外です。Windowsでは、`auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.`という1行だけが出ます。修復はありません。

アカウントの点検の前に、途中で止まった**グループの変更**が`General`ボックスのエラーとして示されます: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``。`--fix`は、計画された順序で、競合に対して与えられた回答のとおりにそれを完了させます。それまでは、ほかのグループの変更、名前の変更、削除は始まりません。その変更に含まれるアカウントは、完了するまで足りないリンクが作られるだけです。完了できない場合、エラーには、そのファイルを削除すると変更の残りをあきらめることになると示されます。`ccam --dry-run doctor --fix`は、その順序（`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`）と、各アカウントの部分で何が退避されるかを示します。ほかのccamがいま実行している変更は情報で、その実行中は`--fix`は共有リンクを修復しません（`its share links are not repaired while another ccam changes the groups; …`）。`--fix`は修復の間`~/.local/state/ccam/group-change.lock`を保持し、変更を完了した後に`config.toml`を読み込み直します。

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

この例では`settings.json`が通常のファイルです。グループにはすでに内容の異なる`settings.json`があるため、競合はデフォルトの選択である`KeepGroup`になり、グループのコピーが残ります。アカウントのファイルはアカウントのフォルダー内の`.ccam-backup/<time>/settings.json`に移動し、`os.Symlink`がそのパスをグループにリンクします。`--fix`はその旨を示します: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`。ほかのケースは次のとおりです。

- グループにまだコピーがない場合は、アカウントのファイルがグループに移動し、グループの最初の内容になります。
- フォルダーの場合は、グループにないファイルがグループに移動します。グループがバイト単位で同じものを持つファイルは捨て、内容の異なるファイルはバックアップに移動し、メモにその名前が示されます（`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`）。
- グループのフォルダーがない場合は、先にそれを作ります。
- `history.jsonl`の場合、デフォルトは`merge by line`です。グループの行の後に、自分の新しい行が続きます。グループの以前のコピーは`.ccam-backup/<time>/group/`に保存されます。
- 別のグループへのリンクは置き換えられます。ccamの外にあるフォルダーへのリンクはバックアップに移動し、グループが空であれば、その内容がグループの最初の内容になります。
- 現在はローカルになった項目のリンクは削除され、アカウントにはグループの内容のコピーが作られます。

リンクの削除と、グループがバイト単位で同じものをすでに持つファイルを捨てることを除き、すべての手順はデータをグループかバックアップフォルダーに移動します。何も失われません。

## 安全性と特殊なケース

- 通常の`ccam doctor`は何も変更しません。スクリプトで確実にしたい場合は、`--fix`と一緒に`--dry-run`を使ってください。
- 実行中のアカウントでは、Claude Codeがどう起動されたかにかかわらず、`doctor --fix`は何も失うことのない、足りないリンクと足りないグループのフォルダーを作るだけです。ほかに何か残っている場合は、`` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``を付け加えます。修復で移動するファイルに、セッションが書き込み続けるためです。
- 共有の修復に失敗した場合、ccamはそのアカウントについて`could not repair share links: <error>`と報告し、不明な項目の点検をスキップします。ほかのアカウントの点検は続きます。
- `ccam group`による変更が途中で失敗した場合は、新しいグループはすでに保存されていて、変更は`~/.local/state/ccam/group-change.json`に記録されています。ccamは`ccam doctor --fix`で完了するよう案内し、これがMCPの手順も含めて変更の残りを実行します。
- シンボリックリンクになった`.credentials.json`、期限切れのログイン、環境変数の認証情報は、報告されるだけで変更されることはありません。
- プラグインのレジストリの書き換えはアトミックに書き込まれますが、`.ccam-backup/`にはコピーされません。

## 関連項目

- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [既存の環境のコピー](/features/migrate)
- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [コマンド](/docs/commands)
