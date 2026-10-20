---
title: 使用量の上限
description: ccam usageはすべてのアカウントの5時間と週間の上限を一度に表示し、ccam statuslineはClaude Codeにアカウント名と使用量を表示します。
updated: 2026-10-20
order: 7
translation: machine
sourceHash: 525174ceed0d
---

`ccam usage`は、各アカウントが5時間と週間の上限をどれだけ使ったかを、それぞれの時間枠がリセットされるまでの時間とともに表示します。長い作業を始める前や、Claude Codeが上限に近いことを知らせたときに使い、どのアカウントに余裕が残っているかを確認します。`ccam statusline`はClaude Code用のステータスラインのコマンドです。Claude Codeの中にアカウント名と使用量を表示し、その数値を`ccam usage`のために保存します。

## 使い方

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

ターミナルでは、各行に`fetching…`と表示された表がすぐに現れ、各アカウントの応答があるたびにその行が描き直されます。`As of`は数値がどれだけ前のものかを示し、`(cached)`はエンドポイントではなく保存済みのファイルから取得した数値であることを意味します。最後の行は、5時間の使用量が100%未満で最も少ないアカウントを示し、アカウントが2つ以上ある場合だけ表示されます。ログインしていないアカウントには`not logged in → ccam login <name>`、APIキーのアカウントには`no limits (API)`と表示されます。

`ccam statusline install`の後は、Claude Codeに`[work] · Opus · 5h 24% · 7d 41%`のような行が表示されます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--json` | 結果をJSONで出力します。`name`、`slot`、`status`、`limits`、`source`、`at`、`error`を含みます | オフ |
| `--watch[=N]` | Ctrl+Cを押すまで`N`秒ごとに描き直します。`--watch`だけの場合や、0以下の値の場合は60です | オフ |
| `-v`、`--verbose` | モデル別の週間上限と、`source:`の行（`live`、`cache`、`statusline`）を表示します | オフ |
| `--no-cache` | 60秒のキャッシュを無視し、エンドポイントに問い合わせ直します | オフ |
| `--source live\|cache` | `cache`はエンドポイントに問い合わせず、保存済みの数値だけを表示します。`live`はデフォルトの動作です | live |
| `--refresh` | 期限切れのアクセストークンを更新します。対象は実行中でないアカウントだけです | オフ |
| `--dry-run` | 受け付けます。`ccam usage`は自身のキャッシュ以外を変更しません。`--refresh`と一緒に使うと`--dry-run cannot be combined with --refresh, which rewrites login tokens`で拒否されます | オフ |

`--source`にそれ以外の値を指定すると、`--source must be live or cache`で失敗します。`ccam usage`は引数を取りません。

`ccam statusline`は、Claude Codeのステータスライン用のJSONを標準入力から読みます。オプションは`--wrap <command>`の1つだけで、既存のステータスラインのコマンドを指定すると、その出力の1行目の前にアカウント名が付きます。`ccam statusline install`はオプションを取りません。どちらも`--dry-run`を`` --dry-run is not supported by `ccam statusline`; nothing was changed ``で拒否します。

## しくみ

数値の取得元は2つあります。各アカウント自身のログインで呼び出す、Anthropicの公開ドキュメントに記載されていない使用量エンドポイント`GET https://api.anthropic.com/api/oauth/usage`と、Claude Codeがステータスラインのコマンドに渡す値です。これらはすべて`internal/usage`にあります。

1. `cli.usageCmd`が`config.toml`を読み込み、アカウントごとに1つの`usage.Account`を作ります。内容は、フォルダー`~/.local/share/ccam/profiles/<name>/`、セッションのフォルダー（実行中かどうかを確認するため）、その場所のままインポートしたかどうかです。
2. `Fetcher.FetchAll`がアカウントごとに1つのgoroutineを起動します。すべてが1つの`http.Client`を共有するので、エンドポイントへの接続は再利用されます。接続には最大2秒、各アカウントには合計3秒かかってもかまいません。リダイレクトするとトークンが別の場所に再送されてしまうため、クライアントはリダイレクトをたどらず、応答は最大1 MBまでしか読みません。
3. 各アカウントの処理が終わるたびに、その行が描き直されます。失敗したアカウントや時間内に応答しなかったアカウントは保存済みの数値にフォールバックするので、1つの遅いアカウントがほかのアカウントを待たせることはありません。

各アカウントはそれぞれのgoroutineでエンドポイントに問い合わせ、タイムアウトしたアカウントは代わりに保存済みの数値を表示します。

<!-- motion:usage-fetch -->

`personal`と`work`は時間内に応答し、最新の数値を表示します。`side`は3秒より長くかかるので、`fallback()`が2つのキャッシュファイルを読み、最後に書き込まれたほうを使います。`usage.Table()`が行を描画し、`usage.Best()`が最後の行に示すアカウントを選びます。

1つのアカウントについて、`Fetcher.Fetch`は次の順序で確認します。

<!-- motion:usage-cache -->

1. `--source cache`の場合は、直接`fallback()`に進みます。
2. `--no-cache`が指定されていなければ、`~/.local/state/ccam/usage/<name>.json`にある、取得してから60秒（`CacheTTL`）未満の最新の結果をそのまま返します。`ccam usage`を続けて2回実行しても、リクエストはアカウントごとに1回で済みます。
3. `readCredentials()`が`.credentials.json`（またはmacOSのキーチェーン）を変更せずに読み、`claudeAiOauth`キーだけをデコードします。サブスクリプションのログインがない場合は`not logged in`、`.claude.json`にAPIキーがある場合は`no limits (API)`になります。
4. 期限切れのアクセストークンは使いません。`--refresh`がない場合、そのアカウントは期限切れの印を付けて保存済みの数値を表示します。
5. それ以外の場合、ccamはエンドポイントを呼び出します。401や403はログインの期限切れ、それ以外のステータスや読めない本文はエラーとして扱い、どちらの場合も保存済みの数値にフォールバックします。429は、エンドポイントへの問い合わせが多すぎることを意味します。ccamは`<name>.pause`に次に問い合わせる時刻（`Retry-After`、なければ5分、最長1時間）を書き込み、保存済みの数値を`rate limited by the usage endpoint; asking again after 15:10`とともに表示し、その時刻までは`--no-cache`を付けてもそのアカウントのエンドポイントを呼び出しません。すぐに再試行しても、制限が続くだけだからです。
6. 正常な応答は、計測した時刻とともに`<name>.json`に書き込まれます。

`fallback()`は`<name>.json`と`<name>.statusline.json`の両方を読み、新しいほうを使います。保存済みの数値を表示する前に、各時間枠のリセット時刻を確認します。数値の保存後にリセットされた時間枠は、古い割合がもう当てはまらないため0%と表示します。上の例で`personal`がリセット時刻なしで`0%`と表示されているのはこのためです。

エンドポイントの応答形式は2つあります。新しい形式には`limits`配列があり、`session`（5時間）、`weekly_all`、`weekly_scoped`（モデル別）の項目を含みます。古い形式には`five_hour`、`seven_day`、`seven_day_opus`、`seven_day_sonnet`のオブジェクトがあります。`ParseLive`は両方を読みます。`limits`配列が優先され、それに含まれないものを古いオブジェクトで補います。上限がまったく含まれない本文は形式の変更として報告され、`ccam doctor`がそれを知らせます。

### トークンの更新

リフレッシュトークンはローテーションします。更新すると新しいリフレッシュトークンが発行され、古いものは無効になります。同じアカウントでClaude Codeが実行中のときにccamがログインを更新すると、そのセッションはサインアウトされてしまいます。そのため、`ccam usage`はトークンを読むだけで、`--refresh`を付けない限り書き込みません。このオプションを付けると、`Fetcher.refresh`が次の処理を行います。

1. `ccam import`でその場所のまま登録したアカウントは拒否します。ccamの知らないところで通常の`claude`がそのログインを使っている可能性があるためです。
2. macOSのキーチェーンに保存されたログインは拒否します。
3. ロックファイル`$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock`（その変数が設定されていない場合は`~/.local/state/ccam/locks/`の下）を取得し、2つのccamのプロセスが同じトークンを二重に更新しないようにします。
4. ロックを保持した状態で、アカウントのセッションが実行中なら拒否し、別のプロセスがすでにトークンを更新していればリクエストをスキップします。
5. リフレッシュトークンをトークンエンドポイントにPOSTし、新しいトークンを`fsutil.AtomicWrite`で`.credentials.json`に書き戻します。ファイル内のほかのキーはすべて残します。

### ステータスラインの収集

Claude Codeはステータスラインのコマンドを実行し、JSONドキュメントをパイプで渡します。ccamはそこからモデル名と`rate_limits`を読み取ります。

<!-- motion:usage-statusline -->

1. `runStatusline`は、ccamがClaude Codeを起動するときに設定する`CCAM_PROFILE`からアカウントを特定します。ない場合は、`CLAUDE_CONFIG_DIR`をアカウントのフォルダーと照合します。
2. `ParseStatusLine`が`rate_limits.five_hour`と`rate_limits.seven_day`を読み取ります。APIキーのアカウントや、セッションの最初の応答より前では、JSONにレート制限が含まれません。その場合は何も保存しません。
3. `SaveStatusLine`が`~/.local/state/ccam/usage/<name>.statusline.json`に書き込みます。エンドポイントが応答しないとき、`ccam usage`はこれを使います。
4. ccamは`[name] · model · 5h N% · 7d N%`を出力します。`--wrap`を付けた場合は、同じ入力でコマンドを`sh -c`で実行し、代わりに`[name]`に続けてコマンドの出力の1行目を表示します。Windowsでは、Claude Codeと同じ方法でコマンドを実行します。Git Bashがインストールされていればそれを使い、なければPowerShellを使います。
5. 5時間または週間の使用量が90%以上で、ほかのアカウントがある場合、行の最後に`→ ccam next --auto --continue`が付きます。たとえば`[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`のようになります。

ラップしたコマンドの制限時間は2秒です。それを過ぎるとコマンドは終了させられ、行にはアカウント名だけが表示されるので、遅いコマンドがClaude Codeのステータスラインを止めることはありません。

`ccam statusline install`は、各アカウントの`settings.json`に`statusLine`を書き込みます。アカウントにすでにステータスラインのコマンドがある場合、ccamはそれを`--wrap '<old command>'`として残します。Windowsでは、ccamへのパスをスラッシュとホームフォルダーを表す`~`で書き込みます（`~/.local/bin/ccam.exe statusline`）。これはGit BashとPowerShellのどちらでも解釈できます。グループで共有されている`settings.json`は、リンクを通じて1回だけ書き込まれるので、グループのすべてのアカウントに反映されます。もう一度実行すると、すでにccamを使っているファイルについては`(already set)`と表示します。

## 安全性と特殊なケース

- `--refresh`を付けない限り、`ccam usage`がトークンファイルを変更することはありません。ログインの期限が切れたときは、`ccam <name>`でそのアカウントを一度起動すれば、Claude Codeが更新します。
- `--refresh`は、実行中のアカウント、インポートしたアカウント、キーチェーンのログインには一切触れません。更新を拒否した場合は、保存済みの数値とその理由を表示します。
- 使用量エンドポイントは公開ドキュメントに記載されておらず、変更される可能性があります。変更された場合でも保存済みの数値は表示され、`ccam doctor`が形式の変更を報告します。
- 保存済みの数値も応答もない場合、その行には`lookup failed:`と理由が表示されます。たとえば`--source cache`では`no saved usage`となります。
- キャッシュファイルに保存されるのは割合とリセット時刻だけで、トークンが保存されることはありません。

## 関連項目

- [使用量の上限とアカウントの切り替え](/docs/usage-limits)
- [次のアカウントへの切り替え](/features/next)
- [使用量の上限に達したときのワークフロー](/guides/usage-limit-workflow)
- [コマンド](/docs/commands)
