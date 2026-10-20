---
title: アカウントの管理
description: ccam add、ccam login、ccam list、ccam rename、ccam remove、ccam importで、Claude Codeのアカウントを追加、ログイン、一覧、名前変更、削除、インポートします。
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 9eb992214132
---

ccamの各アカウントは、独自のログインを持つ別々のClaude Codeの設定フォルダーです。ここで紹介するコマンドで、アカウントの作成、サインイン、一覧表示、名前変更、削除を行い、すでにあるフォルダーを登録します。

## 使い方

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

コマンドが`<name|number>`を受け取る場所では、アカウント名か、`ccam list`の`#`列にあるスロット番号を指定できます。

ターミナルでオプションなしに`ccam add`を実行すると、ウィザードが開きます。名前を尋ね、次に共有する項目のチェックリストを表示し（会話履歴以外のすべてにチェックが入っています）、続いて`Create a shortcut command for work? [y/N]`と尋ね、`y`と答えると`Shortcut name [ccwork]:`と尋ねます。`--share`、`--shortcut`、`--no-login`のいずれかを付けた場合や、ターミナルではなくパイプから実行した場合は、ウィザードは開きません。その場合は名前が必須で、`--shortcut`で名前を指定しない限りショートカットは作成されません。

`ccam list`は、スロット、名前、ショートカット、メールアドレス、共有項目、最終使用日時を表示し、デフォルトのアカウントの前に`*`を付けます。

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json`を付けると、各アカウントは`slot`、`shortcut`、`name`、`email`、`default`、`dir`、`last_used`（一度も起動していないアカウントでは省略されます）、`share`を持つオブジェクトになります。`share`は、各項目（`mcp`、`settings`、`memory`、`extensions`、`history`）をグループ名または`"local"`に対応付けます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `ccam add --no-login` | ログインせずにアカウントを作成します。あとで`ccam login <name>`でログインします | オフ |
| `ccam add --share <item=group>` | 項目をグループで共有します。`item=local`にするとその項目は共有しません。複数回指定できます | `mcp`、`settings`、`memory`、`extensions`は`default`、`history`はローカル |
| `ccam add --shortcut <command>` | このショートカットコマンドも`~/.local/bin`に作成します | なし |
| `ccam add --api --key-from <source>` | ログインの代わりにAPIキーを使うアカウント。[APIキーのアカウント](/features/api)を参照 | オフ |
| `ccam list --json` | アカウントをJSONで出力します | オフ |
| `ccam remove -y`、`--yes` | 確認なしで削除します | オフ |
| `ccam import --name <name>` | 指定したフォルダーのアカウント名 | フォルダー名から決まります |
| `ccam import -y`、`--yes` | 見つかったすべてのフォルダーを確認なしでインポートします | オフ |
| `--dry-run` | 何も変更しない`ccam list`と`ccam completion`では使えます。`add`、`login`、`logout`、`use`、`rename`、`remove`、`import`では拒否され、`` --dry-run is not supported by `ccam add`; nothing was changed ``（コマンド名入り）と表示されます | オフ |

`ccam login`、`ccam logout`、`ccam use`、`ccam rename`には独自のオプションはありません。`ccam completion`は`bash`、`zsh`、`powershell`のいずれかを受け取り、それ以外を指定すると`only bash, zsh and powershell are supported`で失敗します。

アカウント名には、`a-z`、`0-9`、`-`、`_`を使って1～32文字で指定し、先頭は英字か数字にします。`ccam 2`はスロット2を意味するため、数字だけの名前は使えません。また、`ccam <name>`が曖昧にならないように、ccamのコマンド名（`add`、`list`、`ls`、`run`、`use`、`next`、`usage`、`u`、`share`、`mcp`、`login`、`logout`、`rename`、`remove`、`rm`、`import`、`migrate`、`doctor`、`shims`、`shortcut`、`completion`、`help`、`version`、`statusline`、`update`、`folder`、`pick`、`ext`）も使えません。この一覧は`internal/config`の`config.Reserved`です。

## しくみ

`ccam add work --shortcut ccwork`は、`addProfile`（`internal/cli/profiles.go`）で次の手順を実行します。

<!-- motion:accounts-add -->

1. 何かを作成する前に、`config.ValidateName`が名前を、`shim.Check`がショートカットを確認します。
2. `profile.Add`が`~/.local/share/ccam/profiles/work/`をモード`0700`で作成し、空いている最小のスロット番号（`Config.NextSlot`）をアカウントに割り当てます。そのため、削除したアカウントの番号は再利用されます。最初のアカウントがデフォルトになります。
3. `applyShare`が、`share.Build`と`Apply`を使って、作成したばかりの空のフォルダーを共有グループにリンクします。たとえば`settings.json`を`groups/default/settings.json`にリンクします。これが失敗した場合は、新しいフォルダーを削除します。フォルダーが残ると、再試行のたびに`folder already exists`で失敗するためです。
4. アカウントは`fsutil.AtomicWrite`を通じて`~/.config/ccam/config.toml`に保存されます。一時ファイルに書き込み、`fsync`してから`rename`するので、クラッシュしても設定が中途半端な状態で残ることはありません。
5. `runLogin`が、`CLAUDE_CONFIG_DIR`をアカウントのフォルダーに設定して`claude auth login`を起動します。ログインに失敗してもアカウントは残り、ccamはエラーを表示して、あとで`ccam login work`を実行するよう案内します。
6. `profile.ReadEmail`が、`profiles/work/.claude.json`から`oauthAccount.emailAddress`を取り出します。このファイルは数メガバイトになることがあるため、`TopLevelValue`はファイル全体を解析せずに、その1つのキーだけを探します。
7. `shim.Install`が`~/.local/bin/ccwork`をccamのバイナリにリンクします。ccamは起動時に`config.toml`で`argv[0]`を調べます。

ログインには数分かかることがあり、その間に別のccamが設定を保存するかもしれません。そのため、待ち時間のあるコマンドは`env.update`を通じて保存します。これは`config.toml`を読み込み直し、1つの変更（ここではメールアドレス）だけを適用して保存します。`ccam remove`も同じ理由で、確認の後に読み込み直します。

`ccam login`は、既存のアカウントについて手順5と6を実行します。`ccam logout`はアカウント内で`claude auth logout`を実行し、保存済みのメールアドレスを消去します。`ccam use`は`config.toml`の`default`を書き換えるだけです。

### リモートマシンでのサインイン

SSH経由では、`claude auth login`のブラウザー側の処理がうまくいかないことがあります。Claude Codeはサーバー上の`http://localhost:<port>/callback`にあるコールバックサーバーでサインインを待ちますが、サインインを完了するブラウザーは手元のコンピューターで動いているため、そのアドレスはどこにもつながらず、ブラウザーには`ERR_CONNECTION_REFUSED`のようなエラーが表示されます。それでもアドレスにはワンタイムコードが含まれているので、ccamがそれを届けられます。

<!-- motion:accounts-relay -->

`runLogin`はターミナルをClaude Codeに渡しません。Claude Codeの標準入力をパイプにつなぎ、`relayLogin`が入力された各行をそこにコピーします。そのため、Claude Codeの`Paste code here`プロンプトに貼り付けたコードもそのまま届きます。行がコールバックアドレスの場合、`callbackURL`は、ホストが`localhost`、`127.0.0.1`、`[::1]`のいずれかで、ポートが指定され、パスが`/callback`で、`code`と`state`の両方がある場合にだけ受け付けます。`http://`は省略できます。続いて`deliverCallback`が、ブラウザーが本来行うのと同じように、10秒の制限時間を付け、リダイレクトをたどらずに、サーバー上でそのアドレスにリクエストを送ります。リクエストは、待機中の`claude auth login`のコールバックサーバー（図の破線）に届きます。Claude Codeはコードを受け取ってログインを保存します。Linuxではアカウントのフォルダー内の`.credentials.json`に保存されます。貼り付けた行によって、ccamがほかのホストに接続することはありません。

ccamはログインを始める前に次のメッセージを表示します。

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

成功すると`→ Passed the sign-in to Claude Code on this machine.`と表示します。

ログインがターミナルで入力を待っている間、ccamはターミナルのエコーと行モードをオフにし、入力した文字ごとに`*`を表示します。そのため、貼り付けたコードやコールバックアドレスは画面に表示されません。行の編集は`readLine`が自分で行います。`Enter`で行を送信し、`Backspace`で1文字消し、空の行で`Ctrl+D`を押すと入力が終わります。ターミナルの設定は、ログインが終わったときに元に戻ります。`Ctrl+C`、`SIGTERM`、`SIGHUP`を受けたときも元に戻り、これらは従来どおりccamを終了します。標準入力と標準エラー出力がターミナルでない場合は、行を届いたとおりにコピーします。

### 名前の変更、削除、インポート

`ccam rename`と`ccam remove`は、アカウントのフォルダーと、ccamがアカウント名で別の場所に保持しているものを対象にします。それらのパスは`nameState`が列挙します。`~/.local/state/ccam/last-used/<name>`、`~/.local/state/ccam/mcp/<name>`、使用量のキャッシュである`~/.local/state/ccam/usage/<name>.json`と`<name>.statusline.json`です。`ccam remove work`は次の順序でこれらを処理します。

<!-- motion:accounts-remove -->

`ccam remove`は、まずアカウントが実行中なら拒否し、そうでなければ確認を求めます。アカウントが拡張機能を共有している場合は、`RewritePluginPaths`がグループの`installed_plugins.json`と`known_marketplaces.json`を書き換えます。Claude Codeはプラグインの絶対パスを保存しており、それがこのアカウントのフォルダーを経由している可能性があるためです。`profile.Trash`がフォルダーを`~/.local/share/ccam/trash/<name>-<time>/`に移動し、アカウントは`config.toml`から外れます。あとで同じ名前のアカウントを作ったときに引き継がれないように`nameState`のパスは削除され、ショートカットのリンクも削除されます。元に戻すには、ゴミ箱からフォルダーを戻します。

`ccam rename`はフォルダーの名前を変更し、その`plugins/`内のプラグインのパスを古いフォルダーから新しいフォルダーに書き換え、デフォルトの設定を引き継ぎ、`nameState`のパスを新しい名前に移動します。ショートカットコマンドの名前は変わりません。プラグインのパスの書き換えに失敗した場合は、設定が保存されないため、フォルダーの名前を元に戻します。

`ccam import`は、フォルダーをその場所のまま`config.toml`の`dir`として登録します。移動もコピーもしません。リフレッシュトークンをコピーすると、コピー同士が互いをログアウトさせてしまうためです。フォルダーを指定しない場合は、登録されていない`~/.claude`と`~/.claude-*`のフォルダーを探し、それぞれについて確認します。名前はフォルダーから決まり（`~/.claude`は`main`、`~/.claude-work`は`work`、`~/.claude-2`は`acct2`）、名前を決められないフォルダーは、使うべき`--name`コマンドとともに一覧表示されます。`CLAUDE_CONFIG_DIR`がない場合、Claude Codeは`.claude.json`をホームフォルダーの`~/.claude`の隣に置きます。そのため、`~/.claude`をインポートしたときに独自の`.claude.json`がなければ、ccamは`~/.claude.json`を指すリンク`~/.claude/.claude.json`を作成します。これで通常の`claude`と`ccam main`が同じアカウントを参照します。インポートしたアカウントは、`ccam share`を実行するまで何も共有せず、ショートカットも作成されません。

`ccam completion bash|zsh`は、アカウント名も補完する補完スクリプトを出力します。

## 安全性と特殊なケース

- `ccam rename`と`ccam remove`は、アカウントの実行中は`work is running; quit its Claude Code session first`と表示して拒否します（`refuseRunning`）。セッションがフォルダーに書き込み続けるためです。
- インポートしたアカウントに`ccam remove`を実行すると、登録が解除されるだけです（`✓ main unregistered (folder ... left in place)`）。`~/.claude`がゴミ箱に移動されることはありません。
- `ccam add`は、`profiles/`にすでにフォルダーがある名前を拒否し、`ccam import`を案内します。
- `ccam import`は、すでに登録されているフォルダー、フォルダーではないパス、使用済みの名前を拒否します。
- `ccam logout`と`ccam login`は、アカウントが実行中かどうかを確認しません。
- ログインが終了した後などで、中継がClaude Codeに届かない場合は、`ccam: could not pass the sign-in to Claude Code:`と理由が表示されます。HTTPエラーは通常、コードの期限切れを意味します。ログインをやり直してください。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [ショートカットコマンド](/features/shortcuts)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [Claude Codeの複数アカウント](/guides/multiple-claude-code-accounts)
