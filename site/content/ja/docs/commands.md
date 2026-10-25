---
title: コマンドリファレンス
metaTitle: "ccamコマンドリファレンス"
description: ccam addやccam <name>から、ccam usage、ccam next、ccam doctorまで、ccamのすべてのコマンドとオプションを説明します。
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "abba525cae48"
---

コマンドの前に`--dry-run`を付けると、実際には実行せずに何をするかを確認できます。たとえば`ccam --dry-run work`のように使います。これが使えるのは、Claude Codeを起動するとき（`ccam <name>`、`ccam run`、`ccam next`、`ccam attach`）と、`group`、`migrate`、`shortcut`、`update`、`doctor --fix`です。何かを変更するそれ以外のコマンドは、実行せずにこのオプションを拒否します。

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

**TeamとEnterpriseのシート。** 自動的な選択（`ccam next`、`--fallback`、`ccam pick`、ダッシュボードの`n`、有効にした場合の[`auto_switch`](/features/next#auto-switch)）は、同じClaude TeamまたはEnterpriseの組織に属するシート間でも切り替えます。Anthropicのヘルプセンターには、「Each seat on a Team plan is assigned to one named individual, and only that person may use it」（Teamプランの各シートは1人の指定ユーザーに割り当てられ、そのユーザーのみが使用できます）、「Users can't route their work through another user's account to use that user's remaining usage」（ユーザーは別のユーザーのアカウントを通じて作業をルーティングして、そのユーザーの残りの使用量を使用することはできません）と書かれています（[Teamプラン](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)。[Enterpriseプラン](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)も同様です）。そのため、シート間で自動的な選択を使うと、Anthropicの規約に違反するおそれがあります。そこで使うかどうかはご自身の判断です。自動的な選択に頼る前に、所属する組織とAnthropicに確認してください。個人プラン（Pro、Max）については、[FAQの規約に関する説明](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)を参照してください。

`ccwork`のようなショートカットも、同じようにそのアカウントを起動します。`ccwork --resume`は`ccam work --resume`と同じ動作です。詳しくは[ショートカット](/docs/commands#shortcuts)を参照してください。

## アカウントの追加と削除

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

フォルダーを指定せずに`ccam import`を実行すると、`~/.claude`と`~/.claude-*`を探し、それぞれについて確認します。失敗したフォルダーは`Skipped <path>: <reason>`と報告され、ほかのフォルダーは登録されます。名前がすでにあるアカウントの名前と同じフォルダーには`--name`が必要で、これはフォルダーと一緒に指定します。インポートしたフォルダーはその場所のまま登録されます。ccamがログインを移動またはコピーすることはありません。

`ccam login --expired`は、Claude Codeがログインを消去したアカウントにもログインします。`ccam remove`は、管理しているアカウントのフォルダーをごみ箱に移す前に、そのアカウントをログアウトします（同じログインをほかのアカウントが持っている場合は除きます）。インポートしたアカウントを`ccam remove`すると、共有しているもののコピーを持ってグループから抜けます。macOSでは、キーチェーンがフォルダーのパスでログインを保存しているため、`ccam rename`はフォルダーの名前を変えません。

`ccam add --api`は、Claudeのログインの代わりにAPIキーを使うアカウントを追加します。ccamは起動のたびにキーを読み、保存することはありません。[APIキーのアカウント](/features/api)を参照してください。

`ccam rename`、`ccam remove`、`ccam group`の各コマンドは、そのアカウントでClaude Codeが実行中の間は、どう起動されたかにかかわらず（バックグラウンドセッションも含みます）、アカウントを変更しません。先にセッションを終了してください。グループの変更が途中で終わっている間も拒否します: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``。

`ccam list --json`は、すべてのアカウントをスロット、ショートカット、メールアドレス、フォルダーとともに出力します。その`groups`配列はアカウントのグループを一覧し、`share`オブジェクトは、各項目（`mcp`、`settings`、`memory`、`extensions`、`history`）を、その項目の出どころのグループか`local`に対応付けます。

## ショートカット

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

ショートカットは、1つのアカウントでClaude Codeを起動する、自分で名前を付けたコマンドです。`ccam work`と同じように、すべての引数をそのままClaude Codeに渡すので、`ccwork --resume`や`ccwork -p "hi"`も期待どおりに動作します。実体は`~/.local/bin`にあるccamへのリンクです。

ショートカットは任意です。`ccam add`のウィザードが作成するかどうかを尋ね、`ccam add <name> --shortcut <command>`なら確認なしで作成します。`ccam list`は各アカウントのショートカットを表示し、ない場合は`-`と表示します。

名前に使える文字は、a-z、0-9、`.`、`-`、`_`です。ccamのコマンド、ほかのアカウントのショートカット、`PATH`上にあるほかのプログラムと同じ名前は使えません。また、自分で置いたccamのコピーやリンクのように、ccamが作っていない`~/.local/bin`内のファイルと同じ名前も使えません。`config.toml`を手で編集して付けたショートカットが、これらの規則に反している場合、`claude`か`ccam`である場合、または2つのアカウントに付けられている場合は、`config.toml`の読み込みが失敗し、そのアカウントの名前が示されます。

`ccam rename`ではショートカットの名前は変わらず、`ccam remove`ではアカウントと一緒にショートカットも削除されます。`ccam shims install`は、たとえばccamを移動した後などに、すべてのアカウントのリンクを作り直します。1つのショートカットの場所にファイルがあっても、ほかのショートカットは作り直されます。

## 共有

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

グループは、そのアカウントの間で項目を共有します。1つのアカウントは複数のグループに入れますが、各項目を受け取るのは1つのグループからだけです。詳しくは[グループによる共有](/features/share)を参照してください。

## 既存の環境のコピー

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`は、Claude Codeのフォルダーから`settings.json`、`CLAUDE.md`、拡張機能、ユーザーのMCPサーバーをグループにコピーし、それらの項目を共有するすべてのアカウントで使えるようにします。フォルダーは読み取るだけです。`--settings`、`--memory`、`--extensions`、`--mcp`、`--history`の各オプションを付けると、指定した項目だけをコピーします。何も付けない場合、ターミナルではチェックリストが表示され、会話履歴以外のすべてにチェックが入った状態から選べます。ターミナル以外では、会話履歴以外のすべてをコピーします。`-g`または`--group`でグループを選び、`-y`または`--yes`を付けると確認なしでコピーし、内容が異なるファイルはグループのものを残します。MCPのサインイン情報はコピーしません。コピーしたサインインは、ほかの持ち主をログアウトさせてしまうためです。計画には、`/mcp`でサインインするアカウントが示されます。`--mcp`は、`~/.claude`が空か存在しなくても、`~/.claude.json`にサーバーがあれば使えます。詳しくは[既存の環境をコピーする](/docs/sharing#copy-an-existing-setup)を参照してください。

## MCPサーバー

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>`は、そのアカウントだけが持つサーバーを、そのアカウントのグループに移動します。1つのフォルダーにだけ持っているサーバー（オプションなしの`claude mcp add`）も対象で、移動後はどのフォルダーでも読み込まれます。グループがすでに同じ名前を別の定義に使っているサーバーはスキップされます。`--force`を付けるとグループのものを置き換えます。その際、ccamは先にバックアップを取ります。

## 使用量と切り替え

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`、`--fallback`、`ccam pick`は同じ方法で選びます。まず`auto_max_usage`（`config.toml`、デフォルトは100）未満のアカウントを優先し、`auto_strategy`（`headroom`、`expiring`、`balance`）に従って選びます。`--model`の週間の枠も数え、100%のアカウント（Claude Codeが上限に達したと報告した時間枠は、そのリセットまで100%と数えます）とAPIキーのログインは決して選ばず、起点となるアカウントと同じ組織の中から選びます（範囲を広げるには`--cross-org`か`cross_org = true`）。`--wait`の値は`=`の後に書きます（`--wait=2h`）。スペースで区切って書くと、ヒントとともに拒否されます。選べるアカウントがない場合は終了コード75で終了するか、`--wait`を付けていれば最初のアカウントがリセットされるまで待ちます。`ccam next --continue`は、履歴を共有していないアカウントの間でも会話を持ち越します。その場合はClaude Codeが会話を新しいアカウントにコピーします。`--brief`を付けると、新しいアカウントは代わりに新しい会話を始めます。その元になるのは、ccamがモデルを使わずに書く、前の会話の要約です（[`--brief`](/features/next#brief)）。`ccam <name>`と違い、`ccam next`は自分のオプションを読むので、Claude Codeのオプションは`--`の後に置きます。付けないとccamは`unknown flag`で止まります。Claude Codeの中では、プラグインの`/ccam:switch [account]`が会話を引き渡し、そのフックが実行するのが`ccam hook`です。`config.toml`で`auto_switch = true`を設定すると（LinuxとmacOS。デフォルトではオフ）、上限のフックは、ccamが起動した対話型のセッションを終了させ、余裕のある別のアカウントで続けます。移す先は同じ種類のアカウントだけで、同じ組織の別のシートか、別の個人プランです（[`auto_switch`](/features/next#auto-switch)）。[次のアカウントへの切り替え](/features/next)、[スクリプト用のアカウント選び](/features/pick)、[すべてのアカウントのバックグラウンドセッション](/features/agents)を参照してください。

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
| `~/.config/ccam/config.toml` | アカウント、グループ、設定（`theme`、`usage_source`、`relay`、`auto_switch`、`cross_org` …）。設定はファイルの先頭、最初の`[[group]]`、`[[profile]]`、`[[folder]]`より上に書きます。テーブルの下に書いた設定はそのテーブルに属するので、ccamは起動のたびと`ccam doctor`でそれを示し、設定を保存するコマンドは、それが移されるまで拒否します。ccamがファイルを保存するとき、コメントは残りません |
| `~/.local/share/ccam/profiles/<name>/` | 各アカウントの`CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 共有される設定、MCPサーバー、拡張機能、履歴 |
| `~/.local/state/ccam/` | 使用量のキャッシュと上限の記録（`usage/<name>.json`、`usage/<name>.limit.json`）、最終使用日時、セッションの目印とロック（`locks/`）、ccamが作ったショートカット（`shortcuts.json`）、終わっていないグループの変更（`group-change.json`）、`/ccam:switch`と上限のフックが残すメモ（`handoff/<account>/<session>.json`）、`ccam pick`か`auto_switch`が各アカウントを最後に選んだ時刻（`picks/<name>`。選択は2分間有効です）、そして7日間保存される`--brief`の要約（`brief/<session>.md`） |
| `~/.local/bin/<shortcut>` | `ccwork`のようなショートカットコマンド（ccamへのリンク） |

`XDG_CONFIG_HOME`、`XDG_DATA_HOME`、`XDG_STATE_HOME`の各変数を設定すると、これらのフォルダーの場所が変わります。Windowsでは`~`は`%USERPROFILE%`で、ショートカットは`<shortcut>.exe`です。
