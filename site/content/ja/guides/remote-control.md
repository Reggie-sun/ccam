---
title: ccamでRemote Controlを使う
description: ccam <name> rcで任意のアカウントのClaude Code Remote Controlを起動する方法、使えるオプション、共有MCPサーバーのコピーを勧める理由、問題の解決法です。
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc`は、アカウント`work`としてClaude CodeのRemote Controlを起動します。これで、このコンピューター上のセッションをclaude.ai/codeやClaudeモバイルアプリから操作できます。動作は`claude rc`と同じですが、違いが1つあります。アカウントがMCPサーバーを共有している場合、ccamはそれらをアカウントの`.claude.json`にコピーするかどうかを一度だけ確認します。Remote Controlのセッションは、それ以外の場所からサーバーを読まないためです。

## 起動する

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[ショートカットコマンド](/features/shortcuts)でも同じように動きます。`ccwork rc`は`ccam work rc`と同じです。作業したいフォルダーでコマンドを実行してから、`work`と同じClaudeアカウントでサインインしたclaude.ai/codeかClaudeアプリを開いてください。

## オプション

ここでは2種類のオプションが使われます。アカウント名の前に書くかccamの設定に置くccam自身のオプションと、`rc`の後に書いてそのままClaude Codeに渡るRemote Controlのオプションです。

### ccamのオプションと設定

| オプションまたは設定 | 動作 | デフォルト |
| --- | --- | --- |
| `--dry-run`（アカウントの前） | 環境変数とコマンド、MCPの通知があればそれも表示します。確認もコピーもしません | オフ |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | アカウントがMCPサーバーを共有していて、その`.claude.json`にグループの現在のリストがないときに表示されます。`y`でコピーし、`N`またはEnterでコピーせずに起動します | `N` |
| `~/.config/ccam/config.toml`の`mcp_strategy = "sync"` | Remote Controlかどうかにかかわらず、起動のたびにグループのサーバーを`.claude.json`にコピーするので、確認は一切ありません | `"flag"` |
| `ccam share work --mcp off` | このアカウントとのMCPサーバーの共有をやめます。Remote Controlはアカウント自身のサーバーだけで起動し、確認もありません | オン |

### Remote Controlのオプション

`rc`より後のものはすべてClaude Codeが受け取ります。以下はClaude Code 2.1.304の`claude remote-control`のオプションです。お使いのバージョンの一覧は`ccam work rc --help`で表示できます。

| オプション | 動作 |
| --- | --- |
| `--name <name>` | claude.ai/codeに表示されるセッション名 |
| `--remote-control-session-name-prefix <prefix>` | 自動生成されるセッション名の接頭辞（デフォルト: ホスト名） |
| `-c`, `--continue` | Remote Controlがこのフォルダーで最後に記録したセッションに再接続（約4時間以内） |
| `--session-id <id>` | 指定したセッションに再接続 |
| `--permission-mode <mode>` | 新しいセッションの権限モード: `acceptEdits`、`auto`、`bypassPermissions`、`default`、`dontAsk`、`plan`のいずれか |
| `--spawn <mode>` | `same-dir`（デフォルト）、`worktree`（セッションごとにgit worktreeを1つ）、`session`（セッション1つだけで、終わると終了） |
| `--capacity <N>` | `same-dir`と`worktree`モードで同時に開けるセッションの最大数（デフォルト32） |
| `--[no-]create-session-in-dir` | 起動時に現在のフォルダーで最初のセッションを作る（デフォルトはオン） |
| `--[no-]chrome` | 新しいセッションでClaude in Chromeを使う（デフォルトはオフ） |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | デバッグ出力と詳細出力 |

`rc`はアカウント名のすぐ後に書いてください。ccamは`-`で始まらない最初の引数をClaude Codeのコマンドとみなすので、`ccam work rc --name laptop`では`rc`を見ますが、`ccam work --model opus rc`では`opus`を見ます。

## ccamが行うこと

起動はいつもの流れ（[Claude Codeの起動](/features/launch)）をたどり、MCPの手順だけが変わります。

1. ccamは、通常の起動と同じようにアカウントを見つけ、共有リンクを確認し、`claude`の場所を解決します。
2. コマンドが`rc`か`remote-control`であることを確認し、`--mcp-config`を付けません。
3. アカウントがMCPサーバーを共有している場合、グループの`mcp.json`を`~/.local/state/ccam/mcp/work/`の`sync-stamp`と比較します。`.claude.json`に現在のサーバーがすでに入っている場合や、グループにサーバーがない場合は、そのまま手順6に進みます。
4. それ以外の場合は、セッションで使えなくなるサーバーの名前を示します。

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y`と答えると、`.claude.json`のバックアップを保存し、Claude Code自身のロックを取得して、グループのサーバーを`mcpServers`に追加します。アカウントのほかのサーバーは残します。`N`と答えた場合、ターミナルがない場合、アカウントが実行中の場合は、何も変更せずにその旨を表示します。

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. セッションを記録し、`CLAUDE_CONFIG_DIR`をアカウントのフォルダーに設定してから、`exec`で自分自身を`claude rc`に置き換えます。

その後、Remote Controlはサーバーとして動き続けます。あなたのために起動するセッションはすべてこのプロセスの子なので、`work`として動きます。つまり`work`のログイン、設定、会話履歴を使い、`work`の使用量上限に加算されます。サーバーが動いている間、ccamはこのアカウントを実行中として扱います。

違いはドライランで確認できます。通常の起動はグループのサーバーを引数として渡しますが、Remote Controlは渡しません。

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## このように動く理由

**Remote Controlは`--mcp-config`を受け付けません。** 通常の起動では、ccamはグループのサーバーを`--mcp-config=<path>`としてClaude Codeに渡し、Claude Codeのファイルには決して書き込みません。`claude remote-control`はコマンド名の前のこのオプションを拒否し、後ろでも受け付けません。そのため、ccam 0.6.0より前は、MCPサーバーを共有するすべてのアカウントで`ccam work rc`が失敗していました。現在のccamはこのオプションを付けません。

**そのセッションは`.claude.json`からしかサーバーを読みません。** Remote Controlが起動するセッションは、コマンドラインのオプションを引き継ぎません。ユーザーMCPサーバーを見つける場所は、アカウントの`.claude.json`だけです。このファイルにはアカウントのログイン情報やフォルダーごとの状態も入っているため、ccamは`settings.json`のように共有コピーへリンクすることができません。セッションに共有リストを渡すには、サーバーをこのファイルにコピーするしかありません。

**ccamは書き込む前に確認します。** デフォルトのフラグ戦略では、ccamは`.claude.json`に手を触れないと約束しています。Remote Controlはそこに書き込む必要がある唯一のケースなので、まず確認し、バックアップを残し、`mcpServers`以外には何も触れません。

**アカウントの実行中は書き込みません。** Claude Codeはセッション中に`.claude.json`を書き直すため、外部から加えた変更は失われるか上書きされます。そのため、ccamはアカウントのセッションが1つでも実行中ならコピーを見送ります。別のRemote Controlサーバーもこれに含まれます。

**確認は変更ごとに一度だけです。** コピーの後、`sync-stamp`にはグループの`mcp.json`の更新日時が記録されます。ccamが再び確認するのは、たとえば`ccam mcp add`の後など、そのファイルが変わったときだけです。`N`と答えた場合はまだコピーがないので、Remote Controlを起動するたびに確認します。

**`ccam mcp`がコピーを最新に保ちます。** 名前が重なった場合は`.claude.json`の項目が`--mcp-config`より優先されるため、古いコピーは通常の起動でも新しい定義を隠してしまいます。そこで、`ccam mcp add`、`add-json`、`remove`、`import`は、実行中でないすべてのアカウントのコピーを更新します。

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**ターミナルがなければ答えは「いいえ」です。** Remote Controlはスクリプトやサービス、`nohup`など、誰も答えられない場所で動くことがよくあります。その場合、ccamは待たずにサーバーなしで起動します。

## 問題と解決方法

### Remote Controlのセッションに共有MCPサーバーがない

確認に`N`と答えたか、ターミナルがなかったか、アカウントが実行中だったためです。Remote Controlを止め、アカウントのほかのセッションを閉じてから、ターミナルで`ccam work rc`をもう一度実行し、`y`と答えてください。確認を一切なくすには、`~/.config/ccam/config.toml`で`mcp_strategy = "sync"`を設定します。すでに開いているセッションのサーバーは実行中に変わらないので、コピーの後は新しいセッションを始めてください。

### 「Account work is running, so its .claude.json is left alone」と表示される

`work`のセッションがまだ開いています。通常の`ccam work`、ショートカットコマンド、または別のRemote Controlサーバーです。それを終了して`ccam work rc`をもう一度実行してください。ほかのアカウントが実行しているセッションは関係ありません。

### ccamが毎回確認してくる

まだサーバーをコピーしていないためです。一度`y`と答えるか、`mcp_strategy = "sync"`を設定してください。Remote Controlでこれらのサーバーを使うつもりがないなら、このアカウントとのMCP共有をやめます: `ccam share work --mcp off`。

### `ccam work rc`が`--mcp-config`についてのエラーですぐに失敗する

お使いのccamが0.6.0より古く、Remote Controlに`--mcp-config`を渡しています。`ccam update`を実行してください。最新のccamで同じエラーが出る場合は、`ccam work --model opus rc`のように値を取るオプションが`rc`より前にあり、ccamが`opus`をコマンドとみなしています。`rc`をアカウント名のすぐ後に置いてください。

### `ccam mcp`での変更がRemote Controlに反映されない

`ccam mcp add`と`remove`は実行中でないアカウントのコピーだけを更新し、その場合は次のように表示します: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Controlを止めてから起動し直してください。ccamが変更された`mcp.json`に気づき、再び確認します。

### コピーで自分のサーバーが置き換えられた

アカウントにグループのサーバーと同じ名前の独自サーバーがあった場合、コピーがそれを置き換えます。以前の`.claude.json`は`~/.local/state/ccam/mcp/work/backups/`にあり、最新の5つが保管されます。自分の定義をすべてのアカウントで使うには、次にRemote Controlを起動する前に`ccam mcp import work <server> --force`でグループに移してください。

### コピーしたサーバーを削除したい

`ccam mcp remove <server>`は、グループとすべてのコピーからサーバーを削除します。コピー前の状態に戻すには、アカウントのセッションが1つも動いていないときに、`~/.local/state/ccam/mcp/work/backups/`から最新のバックアップを復元してください。

### Remote Controlがフォルダーを信頼するか確認してくる、またはターミナルなしで終了する

フォルダーの信頼はアカウントごとに`.claude.json`に保存されるため、`personal`で信頼したフォルダーも`work`にとっては初めてのフォルダーです。そのフォルダーで一度`ccam work`を実行して信頼を承認してから、そこでRemote Controlを起動してください。

### 「You must be logged in with a Claude account that has a subscription」と表示される

Remote ControlにはClaudeのサブスクリプションでのログインが必要です。`ccam list`でアカウントのメールアドレスを確認し、必要なら`ccam login work`でログインし直してください。

### スマートフォンにセッションが表示されない

Claudeアプリまたはclaude.ai/codeが、別のClaudeアカウントでサインインしています。`ccam list`が`work`について表示するメールアドレスと比べてください。

### このフォルダーで別のClaude Codeが会話履歴を共有しているとccamが警告する

同じ履歴グループの別のアカウントが、同じフォルダーで実行中です。両方のアカウントから同じ会話を続けるセッションがあると、履歴が混ざることがあります。それぞれ別の会話で作業する場合にだけ`y`と答えるか、別のフォルダーでRemote Controlを起動してください。

## 関連項目

- [Claude Codeの起動](/features/launch)
- [共有MCPサーバー](/features/mcp)
- [アカウント間の共有](/features/share)
- [アカウントの管理](/features/accounts)
