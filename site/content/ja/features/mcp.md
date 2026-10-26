---
title: 共有MCPサーバー
description: ccam mcp add、add-json、list、remove、importで、グループごとに1つのMCPサーバーの一覧を管理します。一覧は--mcp-configで渡すか、各アカウントに同期します。
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
---

Claude CodeはユーザーのMCPサーバーを`.claude.json`に保存します。これはアカウントのログインを保持するのと同じファイルなので、アカウント間で共有できません。ccamは共有サーバーをグループのファイル`mcp.json`に保存し、起動時にグループのすべてのアカウントに渡します。`ccam mcp`を使うと、サーバーを一度追加するだけですべてのアカウントで使えます。あるセッション内で追加したサーバーをほかのアカウントに移すには、`ccam mcp import`を使います。

## 使い方

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

変更は各アカウントの次の起動から適用されます。実行中のセッションは、起動時のサーバーを使い続けます。

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

グループの下に、`ccam mcp list`はほかのアカウントに届かないすべてのサーバーを示します。

- アカウントが自分用に追加したサーバーは`only in account work: pg`
- 1つのフォルダー用に追加したサーバー（オプションなしの`claude mcp add`。後述）は`only in account side, for /home/me/code/app only: notes`
- アカウント自身のサーバーがグループのサーバーと同じ名前を持つ場合は`account work has its own …, used in its sessions in place of the group's`
- このフォルダーの`.mcp.json`は、承認を待っているアカウントとともに`waiting for approval in`として示されます
- `in no account, only in ~/.claude.json (plain claude's)`のサーバーは、それらをコピーする[`ccam migrate --mcp`](/features/migrate)のコマンドとともに示されます

グループにすでにある名前を追加すると、その定義が置き換えられ、`added`ではなく`updated`と表示されます。サーバーがアカウントに届くのは、そのアカウントがそのグループに入っていて、グループがMCPを共有している場合（`ccam group set <group> --add mcp`）だけです。グループがMCPを共有していない場合は、`add`がその旨を表示します。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `-g`、`--group <group>` | `add`、`add-json`、`remove`の対象グループ | いずれかのアカウントがMCPを共有しているグループが1つだけならそのグループ、それ以外は`default` |
| `-- <command> [args…]`（`add`） | このコマンドをstdioサーバーとして実行します | |
| `--url <url>`（`add`） | リモートサーバーのURL（http/sse） | |
| `--transport <stdio\|http\|sse>`（`add`） | `--url`サーバーのトランスポート | `http` |
| `--header <'Name: value'>`（`add`） | HTTPヘッダー。複数回指定できます | なし |
| `-e`、`--env <KEY=VALUE>`（`add`） | stdioサーバーの環境変数。複数回指定できます | なし |
| `--force`（`import`） | 同じ名前で内容の異なるグループのサーバーを置き換えます。以前の`mcp.json`はバックアップされます | オフ |
| `--keep-local`（`import`） | アカウントの`.claude.json`の項目を残します | オフ |
| `-h`、`--help` | 各コマンドのヘルプ | |
| `--dry-run`（グローバル） | `mcp list`だけが対応します。`add`、`add-json`、`remove`、`import`は拒否し、何も変更しません | オフ |

`add`には、`--`の後のコマンドか`--url`のどちらかが必要です。どちらもない場合は`a command to run (after --) or --url is required`で失敗します。`--header`は`--url`サーバーに、`--env`はコマンドのサーバーに適用されます。`add-json`は`claude mcp add-json`と同じJSONを受け取り、定義は`command`か`url`を持つオブジェクトでなければなりません。`import`は常にアカウントがMCPを共有しているグループを使い、`list`はすべてのグループを表示するので、`-g`で変わることはありません。

## しくみ

関係するファイルは次のとおりです。

- `~/.local/share/ccam/groups/<group>/mcp.json`: 共有サーバーです。形式は`{"mcpServers": {…}}`で、`--mcp-config`が読む形式です。サーバーの項目にはAPIキーが含まれることが多いため、`fsutil.AtomicWrite`を使ってモード`0600`で書き込まれます。書き込みはロック`<group>/.mcp.json.lock`を取ったときだけ行うので、2つの`ccam mcp add`コマンドを同時に実行しても、どちらのサーバーも残ります。
- `~/.local/share/ccam/profiles/<name>/.claude.json`: アカウントのサーバーを保持する、Claude Code自身のファイルです。ccamは`profile.TopLevelValue`スキャナーでその`mcpServers`の値だけを読み、起動経路でファイル全体を解析することはありません。
- `~/.local/state/ccam/mcp/<name>/`: `seen`ファイル、`sync-stamp`（ccamが書いた各コピーのハッシュを含むので、誰も変更していないコピーと、あなたが編集したコピーを区別できます）、`sync-checked`（sync戦略が最後に見た内容）、`kept`（MCPの共有から抜けたときにアカウントが持ち続けたサーバー）、`.claude.json`のバックアップです。

戦略は2つあり、`~/.config/ccam/config.toml`の`mcp_strategy`で設定します。デフォルトはフラグ戦略です。

<!-- motion:mcp-flag -->

1. ccamはアカウントのMCPグループを見つけ（`launcher.applyMCP`）、グループの`mcp.json`があるかどうかを、解析せずに1回の`stat`で確認します（`mcp.FlagArgs`）。
2. ccamは`--mcp-config=<path>`を`claude`の直後、自分で指定した引数より前に置きます。そのため、自分で指定した`--mcp-config`や`--strict-mcp-config`も、その上に追加で適用されます（`Plan.InsertArgs`）。
3. ccamは、前回の起動以降にアカウントの`.claude.json`が変更されたかを確認し、変更があった場合だけそのサーバーを解析します（`mcp.NewLocal`）。
4. ccamは自身をClaude Codeに置き換えます（`syscall.Exec`）。

パスは2つではなく、`--mcp-config=<path>`という1つの引数で渡されます。Claude Codeの`--mcp-config`は複数の値を取るため、`claude --mcp-config a.json "prompt"`ではプロンプトが2つ目の設定ファイルとして読まれ、`MCP config file not found`で失敗します。1つの引数の形式ならこれを避けられるので、後に続くプロンプトやサブコマンドも動作します。`ccam work mcp list`、`get`、`remove`にはフラグを付けません。Claude Codeがそこではフラグを拒否するためで、代わりにccamがグループのサーバーを標準エラー出力に示します（[詳細](/features/launch#shared-mcp-servers)）。Claude Codeは、アカウント自身のサーバーに加えて`--mcp-config`のサーバーを読み込みます。名前が重なった場合は`.claude.json`の項目が優先されます。フラグ戦略が`.claude.json`に書き込むのは、Remote Controlの場合、アカウントがMCPの共有から抜けるときやもう一度加わるとき（[後述](#leaving-and-joining)）、そして`ccam mcp import`の場合だけです。

`claude remote-control`（`rc`）は`--mcp-config`を受け付けず、そのセッションは`.claude.json`からしかサーバーを読みません。`ccam work rc`では、ccamはこのフラグを付けず、ターミナル上であれば、`sync`戦略と同じようにバックアップを取ってグループのサーバーを`.claude.json`にコピーするかどうかを確認します。コピーができるまで確認し、その後はグループの`mcp.json`が変わったときか、コピーしたサーバーがなくなったときにだけ再び確認します（[Remote Control](/features/launch#remote-control)）。このコピーは名前が重なった場合に`--mcp-config`より優先されるため、`ccam mcp add`、`add-json`、`remove`、`import`は、グループのほかのアカウントが持つコピーも更新します。

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated`の行は、コピーが変わったときにだけ表示されます。古いバージョンのccamが作成したコピーも同じように最新の状態にされ、`ccam migrate`の後のコピーも同様です。

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントの`.claude.json`にだけ保存されます。オプションなしの`claude mcp add`はlocalスコープを使うので、サーバーはその1つのアカウントの、その1つのフォルダーでだけ読み込まれます（`-s user`を付けると、どのフォルダーでもそのアカウントで使えるようになります）。ccamはそれが起きた時点では検出できません。Claude Codeがすでにccamのプロセスを置き換えているためです。そのため、次の起動時に確認します。`seen`ファイルには、`.claude.json`の更新日時と、すでに通知したサーバー名が保存されています。日時が異なる場合、ccamは名前をグループと比較し、新しいサーバーごとにヒントを一度だけ、入力を待たずに表示します。これはどちらの戦略でも行います。

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

アカウント自身のサーバーがグループのサーバーと同じ名前を持つ場合も、ヒントでそれを示します。そのアカウントのセッションでは、アカウントのものが優先されるためです。

`mcp_strategy = "sync"`の場合、サーバーは`claude mcp list`や`/mcp`に通常のユーザーサーバーとして表示されます。ccamは起動前に、各アカウントの`.claude.json`にサーバーを書き込みます。

<!-- motion:mcp-sync -->

1. ccamは、グループの`mcp.json`の更新日時を`sync-stamp`と比較し、アカウントの`.claude.json`を前回見た内容と比較します（`mcp.Sync`）。どちらも変わっていなければ、ほかには何もしません。変わっていれば、グループのすべてのサーバーがまだ`.claude.json`にあるかも確認し、なくなったコピーを戻します（`put back group default's MCP servers that were gone from account work's .claude.json: ctx`）。
2. Claude Code自身のロックを取得します。これはファイルの隣にある`.claude.json.lock`というフォルダーで、`mkdir`で作成します。最大2秒待ち、Claude Codeのロックライブラリと同じく、10秒より古いロックは古くなったものとして扱います。
3. 現在の`.claude.json`を状態フォルダーの`backups/`に保存し、新しいものを5つ残します。
4. `mcpServers`にグループのすべてのサーバーを設定し、以前に自分が書き込んだサーバーのうちグループにもうないものを削除して、何かが変わったときだけ、元のモードのままファイルをアトミックに書き込みます。変更や削除をするのは、自分が書き込み、その後誰も変更していないコピーだけです。グループのサーバーと同じ名前のサーバーをアカウントが独自に持っている場合、ccamはアカウントのサーバーを残し、一度だけその旨を知らせます: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`。
5. 各コピーのハッシュとともに新しいスタンプを記録します。Claude Codeは`--mcp-config`なしで起動します。

アカウントのセッションが実行中の場合は、Claude Codeがどう起動されたかにかかわらず、`ccam: account work is already running, skipping MCP sync`と表示して同期をスキップします。Claude Codeは実行中に`.claude.json`を書き換えるためです。同期に失敗すると`ccam: MCP sync failed (launching anyway): …`と表示され、そのまま起動します。

### 抜けるときと加わるとき

MCPの共有をやめたアカウント（`ccam group leave`、`ccam group set --remove mcp`、`ccam group remove --leave`）は、グループのサーバーを持ち続けます。計画がそれらをバックアップを取ったうえでアカウントの`.claude.json`にコピーし、同じ名前のアカウント自身のサーバーもそのまま残ります。`--fresh`を付けると何も受け取らず、計画に読み込まれなくなるサーバーが示されます。それらのサーバーを持つグループにもう一度加わると、誰も変更していない残されたコピーは削除されるので、その後のグループの編集がアカウントに届きます。

`ccam mcp import`は、アカウント自身のサーバーをそのグループに移動します。

<!-- motion:mcp-import -->

ccamはアカウントのサーバー、つまりアカウント自身のサーバーと1つのフォルダーだけに持つサーバーを読み、`mcp.SameDef`を使って、それぞれをグループの同じ名前のサーバーと比較します。この比較は、キーの順序や空白を無視します。フォルダーのサーバーは、どのフォルダーでも使えるグループのサーバーになります（`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`）。1つの名前に、フォルダーごとに異なる2つの定義がある場合は拒否されます。プラグインのサーバーは、プラグインが持ち込むので対象になりません。グループにないサーバーと、グループに同じ定義で存在するサーバーはコピーされます。定義が異なるサーバーはスキップされます。グループの定義は、ほかのすべてのアカウントが使っているためです。

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force`を付けるとグループのサーバーが置き換えられ、以前の`mcp.json`は先にグループのフォルダー内の`.ccam-backup/<time>/mcp.json`に保存されます。続いてccamは、同じロックのもとで、状態フォルダーにバックアップを取った後、コピーした項目（アカウント自身のものとフォルダーのもの）を1回の書き込みでアカウントの`.claude.json`から削除します（`mcp.Remove`）。ccamは、グループの`mcp.json`を読み込んだときからそれらの項目を削除し終えるまでロックしたままにするので、その間に誰かがグループからサーバーを削除することはありません。名前が重なるとアカウントの項目がグループのものより優先されてしまうため、この処理が必要です。`--keep-local`を付けると、項目はアカウントの`.claude.json`に残ります。

## 安全性と特殊なケース

- **サインイン。** OAuthでサインインするサーバーは、トークンをサーバーの一覧ではなく各アカウントの`.credentials.json`（`mcpOAuth`）に保存します。そのため、サーバーを共有してもサインインは共有されません。各アカウントで一度ずつ`/mcp`でサインインしてください。ccamがサインインをコピーすることはなく、[`ccam migrate`](/features/migrate)もコピーしません。コピーしたリフレッシュトークンは、ほかのすべての持ち主をログアウトさせるためです。[`ccam ext`](/features/ext)と`ccam doctor`は、ほかのアカウントはサインインしているのに、そのアカウントはサインインしていないサーバーを示し、`ccam next --continue`は、新しいアカウントに足りないサーバーを示します。グループがログインファイルを受け取ることはありません。
- **実行中のアカウント。** `import`はグループにコピーしますが、実行中のアカウント自身の項目は残し、その旨を表示します。セッションが終了した後にもう一度実行してください。同期は実行中のアカウントをスキップします。どちらも、ccamを使わずに起動したClaude Codeセッションも検出します。
- **Remote Control用のコピー。** 実行中のアカウントのコピーには触れず、次のように表示します。`Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"`の場合は、起動のたびに同期するため、更新は必要ありません。
- **壊れたファイル。** `.claude.json`が有効なJSONでない場合、同期はファイルに触れずにエラーを報告します。2秒以内にロックを取得できない場合、ccamは`could not lock … (Claude Code is writing the file)`と報告します。
- **MCPを共有していないアカウント。** `import`は拒否し、先にMCPを共有するグループにアカウントを入れる（`ccam group join <group> side`）よう案内します。
- **グループ名。** `-g`はグループ名だけを受け取ります。名前は`groups/`の下のフォルダーになるため、`../x`のような名前は拒否されます。

## 関連項目

- [グループによる共有](/features/share)
- [ccamでRemote Controlを使う](/guides/remote-control)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [既存の環境のコピー](/features/migrate)
- [コマンドリファレンス](/docs/commands)
