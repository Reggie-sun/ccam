---
title: 共有MCPサーバー
description: ccam mcp add、add-json、list、remove、importで、グループごとに1つのMCPサーバーの一覧を管理します。一覧は--mcp-configで渡すか、各アカウントに同期します。
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
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
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
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
```

グループにすでにある名前を追加すると、その定義が置き換えられ、`added`ではなく`updated`と表示されます。サーバーがアカウントに届くのは、そのアカウントがそのグループとMCPを共有している場合（`ccam share <name> --mcp on`）だけです。MCPを共有しているアカウントがまだない場合は、`add`がその旨を表示します。

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: 共有サーバーです。形式は`{"mcpServers": {…}}`で、`--mcp-config`が読む形式です。サーバーの項目にはAPIキーが含まれることが多いため、`fsutil.AtomicWrite`を使ってモード`0600`で書き込まれます。
- `~/.local/share/ccam/profiles/<name>/.claude.json`: アカウントのサーバーを保持する、Claude Code自身のファイルです。ccamは`profile.TopLevelValue`スキャナーでその`mcpServers`の値だけを読み、起動経路でファイル全体を解析することはありません。
- `~/.local/state/ccam/mcp/<name>/`: `seen`ファイル、`sync-stamp`、`.claude.json`のバックアップです。

戦略は2つあり、`~/.config/ccam/config.toml`の`mcp_strategy`で設定します。デフォルトはフラグ戦略です。

<!-- motion:mcp-flag -->

1. `launcher.applyMCP`がアカウントのMCPグループを見つけます。`mcp.FlagArgs`は、グループの`mcp.json`があるかどうかを、解析せずに1回の`stat`で確認します。
2. `Plan.InsertArgs`が、`--mcp-config=<path>`を`claude`の直後、自分で指定した引数より前に置きます。そのため、自分で指定した`--mcp-config`や`--strict-mcp-config`も、その上に追加で適用されます。
3. `mcp.NewLocal`が、前回の起動以降にアカウントの`.claude.json`が変更されたかを確認し、変更があった場合だけそのサーバーを解析します。
4. `syscall.Exec`がccamをClaude Codeに置き換えます。

パスは2つではなく、`--mcp-config=<path>`という1つの引数で渡されます。Claude Codeの`--mcp-config`は複数の値を取るため、`claude --mcp-config a.json "prompt"`ではプロンプトが2つ目の設定ファイルとして読まれ、`MCP config file not found`で失敗します。1つの引数の形式ならそこで区切られるので、続くプロンプトや`mcp list`のようなサブコマンドも動作します。Claude Codeは、アカウント自身のサーバーに加えて`--mcp-config`のサーバーを読み込みます。名前が重なった場合は`.claude.json`の項目が優先されます。フラグ戦略が`.claude.json`に書き込むのは、Remote Controlの場合だけです。

`claude remote-control`（`rc`）は`--mcp-config`を受け付けず、そのセッションは`.claude.json`からしかサーバーを読みません。`ccam work rc`では、ccamはこのフラグを付けず、ターミナル上であれば、`sync`戦略と同じようにバックアップを取ってグループのサーバーを`.claude.json`にコピーするかどうかを確認します。コピーができるまで確認し、その後はグループの`mcp.json`が変わったときにだけ再び確認します（[Remote Control](/features/launch#remote-control)）。このコピーは名前が重なった場合に`--mcp-config`より優先されるため、`ccam mcp add`、`add-json`、`remove`、`import`は、グループのほかのアカウントが持つコピーも更新します。

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

セッション内で`claude mcp add`を使って追加したサーバーは、そのアカウントの`.claude.json`にだけ保存されます。その時点ではccamはClaude Codeに置き換えられていて検出できないため、次の起動時に確認します。`seen`ファイルには、`.claude.json`の更新日時と、すでに通知したサーバー名が保存されています。日時が異なる場合、ccamは名前をグループと比較し、新しいサーバーごとにヒントを一度だけ、入力を待たずに表示します。

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"`の場合、サーバーは`claude mcp list`や`/mcp`に通常のユーザーサーバーとして表示されます。ccamは起動前に、各アカウントの`.claude.json`にサーバーを書き込みます。

<!-- motion:mcp-sync -->

1. `mcp.Sync`が、グループの`mcp.json`の更新日時を`sync-stamp`と比較します。変更がなければ、アカウントの`.claude.json`は開きません。
2. Claude Code自身のロックを取得します。これはファイルの隣にある`.claude.json.lock`というフォルダーで、`mkdir`で作成します。最大2秒待ち、Claude Codeのロックライブラリと同じく、10秒より古いロックは古くなったものとして扱います。
3. 現在の`.claude.json`を状態フォルダーの`backups/`に保存し、新しいものを5つ残します。
4. `mcpServers`にグループのすべてのサーバーを設定し、以前に自分が書き込んだサーバーのうちグループにもうないものを削除し、アカウント自身のサーバーは残して、元のモードのままファイルをアトミックに書き込みます。
5. 新しいスタンプと、自分が管理するサーバー名を記録します。Claude Codeは`--mcp-config`なしで起動します。

アカウントのセッションが実行中の場合は、`ccam: account work is already running, skipping MCP sync`と表示して同期をスキップします。Claude Codeは実行中に`.claude.json`を書き換えるためです。同期に失敗すると`ccam: MCP sync failed (launching anyway): …`と表示され、そのまま起動します。セッション内で追加したサーバーのヒントは、フラグ戦略だけの機能です。

`ccam mcp import`は、アカウント自身のサーバーをそのグループに移動します。

<!-- motion:mcp-import -->

ccamは`mcp.LocalServers`でアカウントのサーバーを読み、`mcp.SameDef`を使って、それぞれをグループの同じ名前のサーバーと比較します。この比較は、キーの順序や空白を無視します。グループにないサーバーと、グループに同じ定義で存在するサーバーはコピーされます。定義が異なるサーバーはスキップされます。グループの定義は、ほかのすべてのアカウントが使っているためです。

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force`を付けるとグループのサーバーが置き換えられ、以前の`mcp.json`は先にグループのフォルダー内の`.ccam-backup/<time>/mcp.json`に保存されます。続いて`mcp.RemoveLocal`が、同じロックのもとで、状態フォルダーにバックアップを取った後、コピーした項目をアカウントの`.claude.json`から削除します。名前が重なるとアカウントの項目がグループのものより優先されてしまうため、この処理が必要です。`--keep-local`を付けると項目を残します。

## 安全性と特殊なケース

- **サインイン。** OAuthでサインインするサーバーは、トークンをサーバーの一覧ではなく各アカウントの`.credentials.json`（`mcpOAuth`）に保存します。そのため、サーバーを共有してもサインインは共有されません。各アカウントで一度ずつ`/mcp`でサインインするか、[`ccam migrate`](/features/migrate)で`~/.claude`からサインイン情報をコピーしてください。グループがログインファイルを受け取ることはありません。
- **実行中のアカウント。** `import`はグループにコピーしますが、実行中のアカウント自身の項目は残し、その旨を表示します。セッションが終了した後にもう一度実行してください。同期は実行中のアカウントをスキップします。
- **Remote Control用のコピー。** 実行中のアカウントのコピーには触れず、次のように表示します。`Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"`の場合は、起動のたびに同期するため、更新は必要ありません。
- **壊れたファイル。** `.claude.json`が有効なJSONでない場合、同期はファイルに触れずにエラーを報告します。2秒以内にロックを取得できない場合、ccamは`could not lock … (Claude Code is writing the file)`と報告します。
- **MCPを共有していないアカウント。** `import`は拒否し、先に`ccam share side --mcp on`を実行するよう案内します。
- **グループ名。** `-g`はグループ名だけを受け取ります。名前は`groups/`の下のフォルダーになるため、`../x`のような名前は拒否されます。

## 関連項目

- [アカウント間の共有](/features/share)
- [MCPサーバー、設定、履歴の共有](/docs/sharing)
- [既存の環境のコピー](/features/migrate)
- [すべてのコマンド](/docs/commands)
