---
title: Claude Codeの起動
description: ccam <name|number>、ccam run、ダッシュボードで任意のアカウントのClaude Codeを起動します。引数はそのままclaudeに渡り、--dry-runで実行内容を確認できます。
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "4812b197f143"
---

`ccam work`は、アカウント`work`としてClaude Codeを起動し、名前より後のものはすべてそのまま`claude`に渡します。ターミナル、スクリプト、パイプで入力を渡す場合など、`claude`と入力するあらゆる場所で使えます。

## 使い方

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

アカウントの名前か、`ccam list`に表示される番号を指定します。ccamは名前より後の引数を読まないので、`ccam work --help`はClaude Codeのヘルプを表示します。名前の後の`--`は省略できます。標準入力、終了コード、シグナルは`claude`と同じように動作します。

最初の引数が`-`で始まり、ccam自身のフラグでない場合、ccamはすべての引数を付けて、名前を指定しないときのアカウントを起動します。それは[`ccam folder`](/features/folder)で現在のフォルダーに設定したアカウントで、なければデフォルトのアカウントです。デフォルトは`ccam use`で設定したアカウントで、設定していない場合は番号が最も小さいアカウントです。

### ダッシュボード

引数なしの`ccam`は、標準入力と標準出力がターミナルであれば[ダッシュボード](/features/dashboard)を開き、そうでなければ`ccam list`を出力します。アカウントに移動し（矢印キーか、番号の入力で）、`Enter`を押すと起動します。`n`は最も余裕があるアカウントを起動し、`o`では`--continue`、`--resume`、そのほかの引数を選べます。ダッシュボードはClaude Codeが起動する前に閉じるので、先にターミナルが元の状態に戻ります。

シェルに`CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`が設定されていると、Claude Codeはアカウント自身のログインではなくそのキーやトークンを使うことがあるため、ccamは起動前にその旨を標準エラー出力に表示します（`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`）。ccamは変数を読むだけです。[APIキーのアカウント](/features/api)では、代わりにこれらの変数を消去します。

`config.toml`で`relay = true`を設定すると、ccamは対話型のセッションを子プロセスとして起動してそのそばに留まり、`/ccam:switch`の後に会話を自動で引き継げるようにします（[詳細](/features/next#inside-claude-code)）。設定しない場合、ccamは自身をClaude Codeに置き換え、そばでは何も動作しません。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--dry-run` | ccamが実行する環境変数とコマンドを表示して終了します | オフ |
| `--fallback` | `-p`/`--print`での実行用。アカウントが使用量の上限に達すると、同じコマンドを、ログイン済みで最も余裕があるアカウントで再実行します（[詳細](/features/next#run-again-at-the-limit)） | オフ |
| `-v`、`--version` | `ccam <version>`を表示します（最初の引数の場合のみ） | |
| `-h`、`--help` | `ccam --help`や`ccam run --help`として、ccamのヘルプを表示します | |

引数の規則:

- ccamが自身のフラグを読むのは、アカウント名より前だけです。名前より後の引数は、`--dry-run`や`--help`も含めてすべてClaude Codeのものです。
- `--dry-run`は、`ccam --dry-run work`、`ccam run --dry-run work`、`ccam --dry-run`（名前を指定しないときのアカウント）の形で使えます。拒否されることはなく、そのとおりに動作します。`--dry-run`と`--fallback`はどちらの順序でも指定できます。
- `ccam run <name|number> [claude args]`は、`ccam <name>`の長い形式です。名前なしの`ccam run`は、フォルダーのアカウントかデフォルトのアカウントを使います。
- アカウント名には`add`、`list`、`run`、`usage`などのccamのコマンド名を使えないので、`ccam <word>`が曖昧になることはありません。

ドライランは、そのままシェル（WindowsではPowerShell）に貼り付けられる1行を出力します。

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## しくみ

起動のたびに、コマンドラインから`exec`まで同じ手順をたどります。

<!-- motion:launch -->

各ボックスはその手順の実行中に強調表示され、その下の行は、Claude Codeが起動するときの環境変数と引数です。

1. `cmd/ccam/main.go`が`app.Main`を呼び出します。ccamが別の名前で起動された場合、それは[ショートカット](/features/shortcuts)なので、直接起動処理に進みます。
2. `app.parse`は、cobraが関わる前に`os.Args`を調べます。`cli.Execute`に進むのは、`ccam`単体（ダッシュボード）、ccamのコマンド名（`config.Reserved`）、`-h`/`--help`、`-v`/`--version`だけです。それ以外はすべて起動です。
3. `launcher.Launch`が`~/.config/ccam/config.toml`を読み込みます。解析するファイルはこれだけです。
4. `launcher.Prepare`が名前か番号でアカウントを見つけ（指定がない場合は、`HereProfile`が作業ディレクトリと`config.toml`の`[[folder]]`エントリーを文字列として比較し、フォルダーのアカウントを選びます）、先頭の`--`を1つ取り除き、`ResolveClaude`で`claude`の実行ファイルを特定します。PATHの検索結果は、ファイルの更新日時とともに`~/.local/state/ccam/claude-path`にキャッシュされるので、ccamが起動のたびに`PATH`を検索することはありません。ファイルが変更されたか、なくなった場合は、検索し直します。`config.toml`の`claude_bin`のパスは、そのまま使われます。ループになるため、ccam自身である`claude`は、inodeで比較して拒否します。
5. `checkLinks`が、アカウントのフォルダー`~/.local/share/ccam/profiles/<name>/`を共有設定と比較します。使う`share.Check`は`lstat`と`readlink`だけを呼び出すので、システムコールは項目ごとに1回か2回です。足りないリンクはその場で作成します。データを失う可能性があるものがあれば、起動を停止します。
6. `applyMCP`が共有MCPサーバーを追加します（後述）。
7. `session.Register`がセッションを記録します。ccamは、自身のPIDを名前とし、プロセスの開始時刻を内容とするマーカーを`$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`（`XDG_RUNTIME_DIR`がない場合は`~/.local/state/ccam/locks/`）に書き込みます。会話履歴を共有するアカウントの場合は、まず`historyCheck`が`history-<group>-<hash>`（ハッシュは現在のフォルダーから求めます）で実行中のセッションを探し、そこにも登録します。`exec`はPIDを変えないので、マーカーはそのままClaude Codeのプロセス自体を指すことになります。保持し続けるロックではなくマーカーを使うのは、Claude Codeが起動したプロセスがロックを引き継ぎ、セッション終了後もロックを持ち続けてしまうためです。
8. `Plan.Exec`がターミナルのタイトルを`claude · <name>`に設定し（標準エラー出力がターミナルの場合）、起動時刻を`fsync`なしで`~/.local/state/ccam/last-used/<name>`に書き込み、`syscall.Exec`を呼び出します。

Claude Codeは、`CLAUDE_CONFIG_DIR`にアカウントのフォルダー、`CCAM_PROFILE`にアカウント名が設定された状態で実行されます。ccamは、これらの変数について引き継いだ値をすべて取り除きます。macOSで`ccam import`で登録した`~/.claude`フォルダーの場合は、`CLAUDE_SECURESTORAGE_CONFIG_DIR`にも空の値を設定します。これにより、そのアカウントと通常の`claude`が同じキーチェーンのログインを使います。

### execを使う理由

`syscall.Exec`はccamのプロセスを`claude`に置き換えます。そのため、ccamがメモリーに残ることはなく、キー入力、シグナル、終了コードはClaude Code自身のものになり、引数はシェルの文字列ではなく配列のまま保たれます。

Windowsには`exec`がありません。そこでは`Plan.Exec`が、同じ環境変数、引数、コンソールで`claude.exe`を子プロセスとして起動し、終了を待って、その終了コードで終了します。Ctrl+Cはコンソールから直接Claude Codeに届くので、ccamはそれを無視します。セッションのマーカーにはccamのPIDが入ります。ccamはClaude Codeとちょうど同じ時間だけ動作し、マーカーはClaude Codeの終了時に削除されます。待っている間、ccamはほかに何もしません。`claude mcp add`で追加したサーバーの検出など、セッション後に必要な処理は次の起動時に行います。ccamは`.claude.json`の更新日時を`~/.local/state/ccam/mcp/<name>/seen`と比較し、変更があった場合だけファイルを解析します。

### コマンドラインの振り分け

3つのコマンドラインが、`internal/app`の中でそれぞれ別の経路をたどります。

<!-- motion:launch-route -->

ショートカットは`argv[0]`で、`ccam <name>`は最初の引数で判定され、どちらも`syscall.Exec`で終わります。cobraに到達するのは、`ccam usage`のようなccamのコマンドだけです。

### 共有MCPサーバー

デフォルトの`flag`戦略では、ccamはグループのMCPサーバーを引数としてClaude Codeに渡し、自分で指定した引数より前に置きます。

<!-- motion:launch-mcp -->

`mcp.FlagArgs`は、グループの`mcp.json`がある場合にだけ引数を返し、`Plan.InsertArgs`がそれを`claude`の直後に置きます。自分で指定した`--mcp-config`や`--strict-mcp-config`も、その上に追加で適用されます。Claude Codeの`--mcp-config`は複数の値を取るので、`--mcp-config=<path>`という1つの引数にしています。2つの引数に分けると、続くプロンプトも設定ファイルとして読まれてしまいます。

`config.toml`で`mcp_strategy = "sync"`を設定すると、ccamは代わりにアカウントの`.claude.json`にサーバーを書き込みます。アカウントがすでに実行中の場合は、これをスキップします。

### Remote Control

`ccam work rc`（または`ccam work remote-control`）は、`work`としてClaude CodeのRemote Controlを起動します。`claude remote-control`は動詞の前の`--mcp-config`を拒否し、動詞の後でも受け付けません。また、そのセッションはアカウントの`.claude.json`からしかMCPサーバーを読みません。そのため、ここではフラグ戦略の動作が異なります。

<!-- motion:launch-rc -->

1. `remoteControl`は、アカウントより後の引数のうち`-`で始まらない最初のものをClaude Codeの動詞とみなします。それが`remote-control`か`rc`の場合、`mcp.FlagArgs`は何も追加しません。
2. `mcp.Synced`は、グループの`mcp.json`の更新日時をアカウントの`sync-stamp`と比較します。`.claude.json`にグループの現在のサーバーがすでに入っている場合や、グループにサーバーがない場合は、確認せずにRemote Controlを起動します。
3. それ以外の場合、ccamはセッションで使えなくなるサーバーの名前を示して確認します。

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y`と答えると、`mcp.Sync`が`sync`戦略と同じ方法でサーバーをコピーします。Claude Codeのロックを取得し、バックアップを取ったうえで、アカウント自身のサーバーは残します。`N`と答えた場合やターミナルがない場合は、それらのサーバーなしでRemote Controlを起動し、起動のたびにサーバーをコピーする`mcp_strategy = "sync"`を案内します。

ccamは、サーバーがコピーされるまでRemote Controlの起動のたびに確認し、コピーした後はグループの`mcp.json`が再び変わったときに一度だけ確認します。その後の通常の起動では`.claude.json`のコピーが`--mcp-config`より優先されるため、`ccam mcp add`、`remove`、`import`はほかのアカウントが持つコピーも更新します（[共有MCPサーバー](/features/mcp)）。

### 10 msの予算

ccamがClaude Codeの起動前に追加する時間は、最大10 ms（p95）です。起動経路では、cobraのコマンドツリーの構築、Bubble Teaの起動、ネットワーク接続を一切行いません。また、パッケージの`init()`関数の合計が5 msを超えると`TestPackageInitIsCheap`が失敗します。`scripts/perf-check.sh`は`hyperfine`で`ccam --dry-run work`、`ccam work`、ショートカットを計測し、10 msを超えると失敗します。`internal/launcher`の`BenchmarkLaunchPath`は、これらの手順をループで実行します。手順ごとの時間を確認するには、`CCAM_TRACE=1`を設定します。

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 安全性と特殊なケース

- 不明な名前や番号を指定すると`no account "<name>"`で停止し、`ccam list`を案内します。アカウントが1つもない場合は、`ccam add <name>`で追加するよう求めます。
- 共有パスが通常のファイルになったり、間違った場所を指していたりする場合、ccamは`stopped because of share link problems:`で停止し、問題を一覧表示して、`ccam doctor --fix`を実行するよう求めます。起動経路でデータを移動することはありません。
- このアカウントと会話履歴を共有する（同じグループの）別のClaude Codeが同じフォルダーで実行中の場合、ccamは履歴が混ざる可能性があることを警告し、`Run anyway? [y/N]`と尋ねます。ターミナルがない場合は、警告を表示して続行します。
- `--dry-run`は、MCPの同期、新しいサーバーの通知、履歴の確認、セッションのマーカー、タイトル、最終使用日時をスキップし、execも行いません。共有リンクの確認は行い、通常の起動と同じように足りないリンクを作成します。
- `work`がすでに実行中のときに`ccam work rc`を実行すると、`.claude.json`には触れずにその旨を表示し、グループのサーバーなしでRemote Controlを起動します。サーバーをコピーするには、セッションが終了した後にもう一度実行してください。`--dry-run`ではメッセージを表示するだけで、確認もコピーも行いません。
- `mcp_strategy = "sync"`の場合、同期に失敗すると`ccam: MCP sync failed (launching anyway): …`と表示され、Claude Codeはそれまでのサーバーで起動します。

## 関連項目

- [ショートカットコマンド](/features/shortcuts)
- [ccamでRemote Controlを使う](/guides/remote-control)
- [アカウントの管理](/features/accounts)
- [コマンドリファレンス](/docs/commands)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [CLAUDE_CONFIG_DIRのしくみ](/guides/claude-config-dir)
