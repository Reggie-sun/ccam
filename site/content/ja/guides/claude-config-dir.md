---
title: CLAUDE_CONFIG_DIRの解説
description: Claude CodeのCLAUDE_CONFIG_DIRの役割、フォルダーの中身、macOSでのログインの保存先、アカウント間で共有できるファイルを説明します。
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR`は、`~/.claude`の代わりにどのフォルダーを使うかをClaude Codeに指定します。ログイン、設定、MCPサーバー、プラグイン、会話履歴など、アカウント固有のものはすべてそこに保存されます。別々のフォルダーを指定すれば、アカウントを分けられます。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## フォルダーの中身

Claude Code 2.1.302で確認した内容です。

| パス | 内容 | 共有 |
| --- | --- | --- |
| `.credentials.json` | ログイントークン（Linux）、MCPのOAuthトークン | 不可 |
| `.claude.json` | アカウント情報、ユーザーのMCPサーバー、プロジェクトごとの状態 | 不可（MCPサーバーは別の方法で可） |
| `settings.json` | 権限、フック、モデル、ステータスライン、プラグイン | 可 |
| `CLAUDE.md` | グローバルな指示 | 可 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 拡張機能 | 可 |
| `projects/` | `.jsonl`ファイル形式の会話、プロジェクトごとのメモリー | 可 |
| `history.jsonl` | プロンプトの履歴 | 可 |
| `sessions/`, `shell-snapshots/`, `cache/`など | 実行時の状態 | 不可 |

## つまずきやすい2つの点

**`CLAUDE_CONFIG_DIR`を設定していない場合、`.claude.json`は`~/.claude`の中にありません。** このときClaude Codeは、フォルダーの隣の`~/.claude.json`にこのファイルを置きます。`CLAUDE_CONFIG_DIR=~/.claude`を使い始めると、Claude Codeは代わりに`~/.claude/.claude.json`を探すため、アカウントが見つかりません。`ccam import`は、この2つをリンクすることでこの問題に対処します。

**macOSでは、ログインはフォルダーではなくキーチェーンにあります。** Claude Codeはログインを`Claude Code-credentials`という名前で保存し、`CLAUDE_CONFIG_DIR`が設定されている場合は、その後ろにフォルダーパスの短いハッシュを付けます。そのため、フォルダーが別であればキーチェーンの項目も別になります。キーチェーンへの書き込みが拒否された場合、Claude Codeは代わりにフォルダー内の`.credentials.json`を使います。

## フォルダー間での共有

シンボリックリンクは機能します。Claude Codeは、リンクされた`settings.json`や`projects/`フォルダーをリンク経由で読み書きし、リンクをそのまま維持します。例外はMCPサーバーです。MCPサーバーはログインと同じ`.claude.json`に保存されるため、代わりに`--mcp-config`で渡してください。

ccamは、まさにこのルールに基づいて作られています。フォルダーをどのようにリンクするかは[共有](/docs/sharing)を、全体像は[複数アカウントのガイド](/guides/multiple-claude-code-accounts)を参照してください。
