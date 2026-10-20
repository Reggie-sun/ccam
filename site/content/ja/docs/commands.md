---
title: コマンドリファレンス
metaTitle: "ccamコマンドリファレンス"
description: ccam addやccam <name>から、ccam usage、ccam next、ccam doctorまで、ccamのすべてのコマンドとオプションを説明します。
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

コマンドの前に`--dry-run`を付けると、実際には実行せずに何をするかを確認できます。たとえば`ccam --dry-run work`のように使います。これが使えるのは、Claude Codeを起動するとき（`ccam <name>`、`ccam run`、`ccam next`）と、`share`、`migrate`、`shortcut`、`update`、`doctor --fix`です。何かを変更するそれ以外のコマンドは、実行せずにこのオプションを拒否します。

<!-- commands -->

## Claude Codeの起動

`ccam <name|number> [arguments]`は、そのアカウントでClaude Codeを起動します。ccamは名前より後の引数を読まず、同じ順序で`claude`に渡します。標準入力、終了コード、シグナルも`claude`そのものとまったく同じように動作します。名前の後の`--`は省略できます。

`ccam -p "summarise"`のように、引数はあってもアカウント名がない`ccam`は、現在のフォルダーに設定されたアカウントを使い、なければデフォルトのアカウントを使います。これらは`ccam folder <name> [folder]`と`ccam use <name>`で設定します。

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

アカウント名の前に置く`--fallback`は、`-p`/`--print`での実行用です。アカウントが使用量の上限に達すると、ccamは同じコマンドを、ログイン済みで最も余裕があるアカウントで再実行し、最後まで完了した実行の出力だけを表示します。詳しくは[フォルダーごとのアカウント](/features/folder)と[上限に達したら再実行](/features/next#run-again-at-the-limit)を参照してください。

`ccwork`のようなショートカットも、同じようにそのアカウントを起動します。`ccwork --resume`は`ccam work --resume`と同じ動作です。詳しくは[ショートカット](/docs/commands#shortcuts)を参照してください。

## アカウントの追加と削除

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

フォルダーを指定せずに`ccam import`を実行すると、`~/.claude`と`~/.claude-*`を探し、それぞれについて確認します。インポートしたフォルダーはその場所のまま登録されます。ccamがログインを移動またはコピーすることはありません。

`ccam add --api`は、Claudeのログインの代わりにAPIキーを使うアカウントを追加します。ccamは起動のたびにキーを読み、保存することはありません。[APIキーのアカウント](/features/api)を参照してください。

`ccam rename`、`ccam remove`、`ccam share`は、そのアカウントでClaude Codeが実行中の間は、アカウントを変更しません。先にセッションを終了してください。

`ccam list --json`は、すべてのアカウントをスロット、ショートカット、メールアドレス、フォルダーとともに出力します。その`share`オブジェクトは、各項目（`mcp`、`settings`、`memory`、`extensions`、`history`）をグループ名または`local`に対応付けます。

## ショートカット

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

ショートカットは、1つのアカウントでClaude Codeを起動する、自分で名前を付けたコマンドです。`ccam work`と同じように、すべての引数をそのままClaude Codeに渡すので、`ccwork --resume`や`ccwork -p "hi"`も期待どおりに動作します。実体は`~/.local/bin`にあるccamへのリンクです。

ショートカットは任意です。`ccam add`のウィザードが作成するかどうかを尋ね、`ccam add <name> --shortcut <command>`なら確認なしで作成します。`ccam list`は各アカウントのショートカットを表示し、ない場合は`-`と表示します。

名前に使える文字は、a-z、0-9、`.`、`-`、`_`です。ccamのコマンド、ほかのアカウントのショートカット、`PATH`上にあるほかのプログラムと同じ名前は使えません。

`ccam rename`ではショートカットの名前は変わらず、`ccam remove`ではアカウントと一緒にショートカットも削除されます。`ccam shims install`は、たとえばccamを移動した後などに、すべてのアカウントのリンクを作り直します。

## 共有

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

各項目には、`on`（`default`グループ）、`off`、またはグループ名を指定します。詳しくは[共有](/docs/sharing)を参照してください。

## 既存の環境のコピー

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`は、Claude Codeのフォルダーから`settings.json`、`CLAUDE.md`、拡張機能、ユーザーのMCPサーバーをグループにコピーし、それらの項目を共有するすべてのアカウントで使えるようにします。フォルダーは読み取るだけです。`--settings`、`--memory`、`--extensions`、`--mcp`、`--history`の各オプションを付けると、指定した項目だけをコピーします。何も付けない場合、ターミナルではチェックリストが表示され、会話履歴以外のすべてにチェックが入った状態から選べます。ターミナル以外では、会話履歴以外のすべてをコピーします。`-g`または`--group`でグループを選び、`-y`または`--yes`を付けると確認なしでコピーし、内容が異なるファイルはグループのものを残します。また、サーバーの保存済みのサインイン情報を各アカウント自身のログインファイルにコピーします。詳しくは[既存の環境をコピーする](/docs/sharing#copy-an-existing-setup)を参照してください。

## MCPサーバー

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>`は、そのアカウントだけが持つサーバーを、そのアカウントのグループに移動します。グループがすでに同じ名前を別の定義に使っているサーバーはスキップされます。`--force`を付けるとグループのものを置き換えます。その際、ccamは先にバックアップを取ります。

## 使用量と切り替え

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`、`--fallback`、`ccam pick`は同じ方法で選びます。まず`auto_max_usage`（`config.toml`、デフォルトは100）未満のアカウントを優先し、`auto_strategy`（`headroom`または`expiring`）に従って選びます。`--model`の週間の枠も数え、100%のアカウントは決して選びません。選べるアカウントがない場合は終了コード75で終了するか、`--wait`を付けていれば最初のアカウントがリセットされるまで待ちます。[次のアカウントへの切り替え](/features/next)と[スクリプト用のアカウント選び](/features/pick)を参照してください。

## メンテナンス

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update`は、ccamを置き換える前にリリースのSHA-256チェックサムを検証します。`--version v0.2.0`を指定すると、特定のリリースをインストールします。`.deb`パッケージからインストールしたccamは、代わりにaptで更新します。詳しくは[アップデート](/docs/install#update)を参照してください。

## ファイル

| パス | 内容 |
| --- | --- |
| `~/.config/ccam/config.toml` | アカウントと、それぞれが共有する項目 |
| `~/.local/share/ccam/profiles/<name>/` | 各アカウントの`CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共有される設定、MCPサーバー、拡張機能、履歴 |
| `~/.local/state/ccam/` | 使用量のキャッシュと最終使用日時 |
| `~/.local/bin/<shortcut>` | `ccwork`のようなショートカットコマンド（ccamへのリンク） |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME`、`XDG_STATE_HOME`の各変数を設定すると、これらのフォルダーの場所が変わります。Windowsでは`~`は`%USERPROFILE%`で、ショートカットは`<shortcut>.exe`です。
