---
title: 拡張機能の確認
description: ccam extは、各アカウントが読み込むスキル、コマンド、プラグイン、MCPサーバー、フックと、それぞれの出どころ、壊れている箇所を表示します。
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

各セルは、そのアカウントでのその項目の出どころを示します。共有グループ（`default`）、アカウント自身のフォルダーにある場合は`account`、そのアカウントにない場合は`-`です。`!`は問題のある項目を示し、その説明は表の下にあります。

## 読み取る場所

| 種類 | 場所 | 表示される出どころ |
| --- | --- | --- |
| スキル | `skills/`の中で`SKILL.md`があるフォルダー | `skills/`のリンク先のグループ、または`account` |
| コマンド、エージェント | `commands/`と`agents/`の下の`.md`ファイル。サブフォルダーは`folder/name` | グループ、または`account` |
| プラグイン | `plugins/installed_plugins.json` | グループ、または`account` |
| MCPサーバー | グループの`mcp.json`（MCPを共有しているとき）と、アカウントの`.claude.json` | グループ、または`account` |
| フック | `settings.json`の`hooks`。イベントごとに数える | グループ、または`account` |
| このプロジェクト | 現在のフォルダーの`.claude/skills`、`.claude/commands`、`.claude/agents`、`.claude/settings.json`、`.claude/settings.local.json`、`.mcp.json` | すべてのアカウントについて1回だけ表示 |

## 報告する問題

- **壊れたリンク**: `skills/`、`commands/`、`agents/`の中のリンク先が消えたリンク、または存在しない場所にリンクしている共有フォルダー（`skills`、`commands`、`agents`、`plugins`）です（`ccam doctor --fix`で修復できます）。
- **`SKILL.md`がないスキルのフォルダー**: Claude Codeはスキップします。
- **2か所で異なる定義になっているMCPサーバー**: グループの`mcp.json`とアカウントの`.claude.json`、またはそのどちらかとプロジェクトの`.mcp.json`です。両方の設定が同じであれば（`mcp_strategy = "sync"`はそのように書き込みます）問題ではありません。アカウント独自のサーバーをグループに移すには、`ccam mcp import <account> <server>`を使います。
- **有効なJSONではない`installed_plugins.json`。**

## 安全性と特殊なケース

- `ccam ext`は何も変更せず、ネットワークリクエストも行いません。読むのはサーバー名とフックの数だけで、MCPの定義にある`env`やヘッダーの値を読むことはありません。
- ディスク上の状態をそのまま表示します。インストールされていても設定でオフになっているプラグインや、管理設定でブロックされたサーバーも表示されます。
- プロジェクトのフォルダーで実行すると、そのプロジェクトの`.claude`フォルダーと`.mcp.json`も含まれます。

## 関連項目

- [アカウント間の共有](/features/share)
- [共有MCPサーバー](/features/mcp)
- [点検と修復](/features/doctor)
