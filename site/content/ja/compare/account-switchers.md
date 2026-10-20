---
title: ccamとほかのClaude Codeアカウント切り替えツール
description: ccamとclaude-swap、CCS、clauth、CAAM、cuxなど、Claude Codeのアカウントを切り替えるツールを比較し、用途に合う方式を説明します。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Claude Codeのアカウントを切り替えるツールの方式は2つに分かれます。多くは1つの`~/.claude`フォルダーの中でログインを入れ替えるので、使えるアカウントは一度に1つです。ccamを含むほかのツールは、`CLAUDE_CONFIG_DIR`でアカウントごとに専用のフォルダーを割り当てるので、どのアカウントもログインしたままで、複数を同時に実行できます。ccamはさらに設定とMCPサーバーを共有し、別のアカウントのフォルダーで会話を続けられます。

## ログインを入れ替えるか、フォルダーを分けるか

**ログインを入れ替える方式**では、フォルダーが1つしかないため、設定、MCPサーバー、履歴も1組だけです。切り替えは速く、実行中のセッションが新しいログインを引き継ぐこともできます。ただし、使えるアカウントは一度に1つで、すべてのターミナルとエディターに同時に適用されます。

**フォルダーを分ける方式**は、[Claude Codeのドキュメントで説明されている方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)です。どのアカウントもログインしたままで、仕事用と個人用のセッションを並行して実行できます。その代わり、ツールで共有しない限り、各フォルダーは設定、MCPサーバー、履歴をそれぞれ別に持つ状態から始まります。

## ツールの比較

2026-10-19時点の各プロジェクトのドキュメントに基づきます。「記載なし」は、プロジェクトが言及していないという意味で、機能がないという意味ではありません。

| ツール | 方式 | 複数アカウントの同時使用 | 5時間と週間の使用量 | 上限に達したとき | プラットフォーム |
| --- | --- | --- | --- | --- | --- |
| ccam | フォルダーを分ける | 可 | 可、`ccam usage` | `ccam next --auto --continue`を実行し、最も余裕があるアカウントで会話を続ける | Linux、macOS、WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | ログインを入れ替える。`cswap run`ではフォルダーを分ける | 可（`cswap run`） | 可 | `cswap auto`がしきい値でログインを切り替え、実行中のセッションはそのまま続く | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | フォルダーを分ける | 可 | 一部（macOSのメニューバーアプリ） | アカウントプールがローカルプロキシ経由でリクエストを次のアカウントに回す | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | ログインを入れ替える。`clauth start`ではフォルダーを分ける | 可（`clauth start`） | 可 | フォールバックチェーンに沿って切り替える | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | ログインを入れ替える、または分離したプロファイル | 可（分離したプロファイル） | 可 | `caam run`が非対話型のコマンドを別のアカウントで再試行する | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | ログインを入れ替え、`claude`をラップする | 記載なし | 可 | ログインを入れ替えて同じ会話を再開する | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | ログインを入れ替える | 不可 | 可 | 記載なし | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | ログインを入れ替える | 不可 | 可 | 100%で切り替え（オプション） | macOS |

ccswapはclaude-swapのフォークで、Codexのログインも管理します。これらのツールのいくつかは、CodexやGeminiなどほかのコーディングエージェントにも対応しています。

## ccamが重視していること

- どのアカウントもログインしたままです。切り替えでログアウトすることはなく、`ccam work`と`ccam personal`を2つのターミナルで同時に実行できます。
- MCPサーバー、`settings.json`、`CLAUDE.md`、エージェント、スキル、プラグインは、項目ごと、アカウントごとにリンクでフォルダー間で共有されます。ログインは決して共有されません。
- 履歴を共有していれば、`ccam next --auto --continue`という1つのコマンドで、最も余裕があるアカウントのフォルダーで同じ会話を開けます。
- ccamはClaude Codeを起動し、自身をClaude Codeに置き換えるため、バックグラウンドでは何も動きません。使用量を監視したり、自動で切り替えたりはしません。切り替えるタイミングはあなたが決めます。

## ほかのツールのほうが得意なこと

- claude-swap、CCS、clauth、CAAM、cuxは、WSLなしのWindowsで動作します。ccamは動作しません。
- claude-swap、clauth、cuxは、アカウントが使用量のしきい値を超えると自動で切り替えます。ccamは切り替えをあなたに任せます。
- CCSwitcher、Claude Switcher、claude-swapは、macOSのメニューバーに使用量を表示します。
- CCS、CAAM、ccswap、Claude Switcherは、Codexやほかのコーディングエージェントも管理します。

## どれが合うか

一度に1つのアカウントしか使わない場合、切り替えを自動で行いたい場合、WSLなしのWindowsを使っている場合は、ログインを入れ替えるツールを選んでください。複数のアカウントを並行して使い、1つの共通の環境を保ち、次のアカウントで同じ会話を続けたい場合は、ccamを選んでください。何も共有する必要がなければ、[普通のシェルエイリアス](/compare/aliases)で十分です。

複数のアカウントを使ってよいかどうかは、各アカウントの規約によります。使う前に確認してください。
