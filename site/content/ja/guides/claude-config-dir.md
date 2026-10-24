---
title: CLAUDE_CONFIG_DIRの解説
description: Claude CodeのCLAUDE_CONFIG_DIRの役割、フォルダーの中身、macOSでのログインの保存先、アカウント間で共有できるファイルを説明します。
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR`は、`~/.claude`の代わりにどのフォルダーを使うかをClaude Codeに指定します。ログイン、設定、MCPサーバー、プラグイン、会話履歴など、アカウント固有のものはすべてそこに保存されます。別々のフォルダーを指定すれば、アカウントを分けられます。

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIRは公式にサポートされていますか？

はい。Claude Codeは[環境変数](https://code.claude.com/docs/en/env-vars)の一覧で、設定フォルダーを変更する方法としてこれを挙げており、「複数のアカウントを並行して実行するのに便利」と説明しています。[認証ガイド](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)でも、複数のアカウントでログインするためにこれを使っています。シェルやユーザー設定で設定できますが、プロジェクトの設定では設定できません。

## フォルダーの中身は？

Claude Code 2.1.309で確認した内容です。

| パス | 内容 | 共有 |
| --- | --- | --- |
| `.credentials.json` | ログイントークン（Linux、Windows）、MCPのOAuthトークン | 不可 |
| `.claude.json` | アカウント情報、ユーザーのMCPサーバー、プロジェクトごとの状態 | 不可（MCPサーバーは別の方法で可） |
| `settings.json` | 権限、フック、モデル、ステータスライン、プラグイン | 可 |
| `CLAUDE.md` | グローバルな指示 | 可 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 拡張機能 | 可 |
| `projects/` | `.jsonl`ファイル形式の会話、プロジェクトごとのメモリー | 可 |
| `history.jsonl` | プロンプトの履歴 | 可 |
| `sessions/`, `shell-snapshots/`, `cache/`など | 実行時の状態 | 不可 |

## つまずきやすい点は？

**`CLAUDE_CONFIG_DIR`を設定していない場合、`.claude.json`は`~/.claude`の中にありません。** このときClaude Codeは、フォルダーの隣の`~/.claude.json`にこのファイルを置きます。`CLAUDE_CONFIG_DIR=~/.claude`を使い始めると、Claude Codeは代わりに`~/.claude/.claude.json`を探すため、アカウントが見つかりません。`ccam import`は、この2つをリンクすることでこの問題に対処します。`CLAUDE_CONFIG_DIR=~/.claude`で以前に実行したことがあり、`~/.claude`にすでに独自の`.claude.json`がある場合、どちらが正しいかはあなたにしかわからないので、ccamは両方のファイルをそのまま残し、その旨を知らせます。

**macOSでは、ログインはフォルダーではなくキーチェーンにあります。** Claude Codeはログインを`Claude Code-credentials`という名前で保存し、`CLAUDE_CONFIG_DIR`が設定されている場合は、その後ろに`-`と、フォルダーパスのSHA-256の先頭16進8桁を付けます。アカウントは`$USER`で指定されたものです。そのため、フォルダーが別であればキーチェーンの項目も別になり、フォルダーを移動したり名前を変えたりするとそのログインは失われます。macOSで`ccam rename`がフォルダーの名前を変えないのはこのためです。キーチェーンへの書き込みが拒否された場合、[ドキュメント](https://code.claude.com/docs/en/authentication#credential-management)に書かれているとおり、Claude Codeは代わりにフォルダー内の`.credentials.json`を使います。

**すべてのサインインがフォルダーに保存されるわけではありません。** APIキーを使わないClaude Consoleのサインインはフォルダーの外に保存されるため、フォルダーを分けてもこの種のサインインを2つ使い分けることはできません。claude.aiのログインとAPIキーはフォルダーごとに保存されます。

## Windowsではフォルダーはどこにありますか？

デフォルトでは`%USERPROFILE%\.claude`で、ログインはその中の`.credentials.json`に保存されます。Claude CodeがWindows資格情報マネージャーを使うように設定されている場合は、macOSと同じくフォルダーパスから作った名前でそこに保存されます。追加するアカウントごとに、LinuxやmacOSと同じように`CLAUDE_CONFIG_DIR`を別のフォルダーに設定します。

## フォルダー間でファイルを共有できますか？

シンボリックリンクは機能します。Claude Codeは、リンクされた`settings.json`や`projects/`フォルダーをリンク経由で読み書きし、リンクをそのまま維持します。例外はMCPサーバーです。MCPサーバーはログインと同じ`.claude.json`に保存されるため、代わりに`--mcp-config`で渡してください。

ccamは、まさにこのルールに基づいて作られています。フォルダーをどのようにリンクするかは[共有](/docs/sharing)を、全体像は[複数アカウントのガイド](/guides/multiple-claude-code-accounts)を参照してください。
