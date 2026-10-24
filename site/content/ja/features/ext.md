---
title: 拡張機能の確認
description: ccam extは、各アカウントが読み込むスキル、コマンド、プラグイン、MCPサーバー、フックと、それぞれの出どころ、壊れている箇所を表示します。
updated: 2026-10-23
order: 14
translation: machine
sourceHash: "45740e4efbef"
---

`ccam ext`は、「このスキル（またはMCPサーバー）が、あるアカウントでは動くのに別のアカウントでは動かないのはなぜか」に答えます。アカウントごとに、Claude Codeが読み込む拡張機能と、それぞれの出どころ（共有グループ、アカウント自身のフォルダー、今いるプロジェクト）を一覧にします。表の下には、拡張機能が期待どおりに読み込まれない原因を挙げます。読み取るだけです。

## 使い方

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

各セルは、そのアカウントでのその項目の出どころを示します。共有グループ（`default`）、アカウント自身のフォルダーにある場合は`account`、そのアカウントがこのフォルダーだけに持つMCPサーバーは`local`、そのアカウントにない場合は`-`です。インストールされていてもオンになっていないプラグインには、出どころの後に`(off)`が付きます。`!`は問題のある項目を示し、その説明は表の下にあります。通常の`claude`だけが持つサーバーのように、どのアカウントにも属さない問題は、名前なしで最後に示されます。`ccam ext --json`は同じ報告（`accounts`、プロジェクト、`problems`）を出力し、Claude Codeがネットワーク経由で接続するMCPサーバーには`remote`が付きます。アカウントがない場合は、`accounts: []`の報告を出力します。

## 読み取る場所

| 種類 | 場所 | 表示される出どころ |
| --- | --- | --- |
| スキル | `skills/`の中で`SKILL.md`があるフォルダー | `skills/`のリンク先のグループ、または`account` |
| コマンド、エージェント | `commands/`と`agents/`の下の`.md`ファイル。サブフォルダーは`folder/name` | グループ、または`account` |
| プラグイン | `plugins/installed_plugins.json`。アカウントとプロジェクトの設定の`enabledPlugins`でオンまたはオフになります。1つのプロジェクトだけにインストールしたプラグインは、そのプロジェクトでだけ表示されます | グループ、または`account`。オンになっていなければ`(off)` |
| MCPサーバー | グループの`mcp.json`（MCPを共有しているとき）、アカウントの`.claude.json`（アカウント自身のサーバーと、このフォルダーだけに持つサーバー）、そしてプラグインのサーバー（`plugin:<plugin>:<server>`として） | グループ、`account`、または`local` |
| フック | `settings.json`の`hooks`。イベントごとに数える | グループ、または`account` |
| このプロジェクト | 現在のフォルダーの`.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json`、`.mcp.json` | 1回だけ表示。そのフォルダーを信頼していないアカウントには、その`.mcp.json`のサーバーについて問題が示されます |

## 報告する問題

- **壊れたリンク**: `skills/`、`commands/`、`agents/`の中のリンク先が消えたリンク、または存在しない場所にリンクしている共有フォルダー（`skills`、`commands`、`agents`、`plugins`）です（`ccam doctor --fix`で修復できます。グループのフォルダーがなければ作り直します）。
- **`SKILL.md`がないスキルのフォルダー**: Claude Codeはスキップします。
- **2か所で異なる定義になっているMCPサーバー**: グループの`mcp.json`とアカウントの`.claude.json`、またはそのどちらかとプロジェクトの`.mcp.json`です。両方の設定が同じであれば（`mcp_strategy = "sync"`はそのように書き込みます）問題ではありません。アカウント独自のサーバーをグループに移すには、`ccam mcp import <account> <server>`を使います。
- **有効なJSONではない`installed_plugins.json`**、`enabledPlugins`に名前のないプラグイン、フォルダーがないプラグイン。
- **1つのアカウントだけが読み込むサーバー**: このアカウントの`.claude.json`にだけあり、グループのほかのアカウントが読み込まないサーバーと、このフォルダーについてアカウントのlocalスコープにあるサーバー（`-s user`を付けない`claude mcp add`がそうします）です。それぞれに、それを共有する`ccam mcp import`のコマンドが付きます。通常の`claude`の`~/.claude.json`にあり、どのアカウントも読み込まないサーバーには、`ccam migrate --mcp`のコマンドが付きます。
- **承認待ちのプロジェクトのサーバー**: アカウントがまだそのフォルダーを信頼していないので、Claude Codeが一度確認し、それまでプロジェクトの`.mcp.json`のサーバーは待機します。問題の説明には、承認する場所が示されます。設定グループのアカウント、またはこのアカウントの`settings.json`にある`"enabledMcpjsonServers"`です。gitリポジトリのメインのワークツリーに与えた信頼は、ほかのワークツリーでも有効です。
- **サインインとコネクター**: ほかのアカウントはサインインしているのに、このアカウントではサインインが必要なMCPサーバー（サインインは1つのアカウントにだけ有効です。`/mcp`でサインインしてください）と、ほかのアカウントが使ったことがあり、このアカウントは使ったことのないclaude.aiのコネクターです。

## 安全性と特殊なケース

- `ccam ext`は何も変更せず、ネットワークリクエストも行いません。読むのはサーバー名とフックの数だけで、MCPの定義にある`env`やヘッダーの値を読むことはありません。サインインの確認では、各アカウントのログイン（macOSではキーチェーン）にあるMCPサインインの名前を読みますが、トークンは読みません。また、それを行うのは、ほかのアカウントも読み込むサーバーを読み込むアカウントについてだけです。
- ディスク上の状態をそのまま表示します。管理設定でブロックされたサーバーも表示され、設定でオフになっているプラグインには`(off)`が付きます。
- プロジェクトのフォルダーで実行すると、そのプロジェクトの`.claude`フォルダーと`.mcp.json`も含まれます。

## 関連項目

- [アカウント間の共有](/features/share)
- [共有MCPサーバー](/features/mcp)
- [点検と修復](/features/doctor)
