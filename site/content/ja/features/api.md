---
title: APIキーのアカウント
description: ccam add --apiは、AnthropicのAPIやゲートウェイ用のAPIキーのアカウントを追加します。キーは起動時に変数かコマンドから読み、保存しません。
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

APIキーのアカウントは、Claudeのサブスクリプションのログインの代わりに、APIキーでClaude Codeを起動します。接続先はAnthropicのAPIか、LiteLLMやOpenRouterのゲートウェイのように同じMessages形式を話すエンドポイントです。ccamはキーを保存しません。アカウントを起動するたびに、環境変数か、パスワードマネージャーなどのコマンドからキーを読み、そのClaude Codeのプロセスだけに渡します。

## 使い方

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

APIキーのアカウントも、ほかのアカウントと同じように設定、`CLAUDE.md`、拡張機能、MCPサーバーを共有し、ショートカット、`ccam folder`、`ccam ext`でも使えます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--api` | ログインする代わりにAPIキーのアカウントを追加します | オフ |
| `--key-from env:NAME` | 環境変数`NAME`からキーを読みます | |
| `--key-from command:COMMAND` | `COMMAND`を実行し（`sh -c`、Windowsでは`cmd /C`）、出力の最初の行を使います | |
| `--base-url <url>` | エンドポイント（`ANTHROPIC_BASE_URL`） | AnthropicのAPI |
| `--auth api-key\|bearer` | キーを`ANTHROPIC_API_KEY`（`x-api-key`）として送るか、`ANTHROPIC_AUTH_TOKEN`（`Authorization: Bearer`。ゲートウェイは通常こちらを求めます）として送ります | `api-key` |
| `--model <model>` | モデル（`ANTHROPIC_MODEL`） | Claude Codeのデフォルト |

設定は`config.toml`に保存されます。

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## しくみ

アカウントを起動するたびに、通常の確認のあと、Claude Codeが起動する前に次のことを行います。

1. ccamがキーを読みます。`env:`は環境から、`command:`はコマンドを30秒の制限時間で実行して読みます。コマンドはターミナルでパスワードを尋ねることができ、エラー出力もターミナルに出ます。キーが空の場合は、キーをどこで探したかを示すメッセージを出して起動を中止します。
2. Claude Codeのサインイン方法や使うプロバイダーを決める、引き継いだ変数をすべて削除します。`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_BASE_URL`、`ANTHROPIC_MODEL`、`CLAUDE_CODE_USE_BEDROCK`、`CLAUDE_CODE_USE_VERTEX`、`CLAUDE_CODE_USE_FOUNDRY`です。ほかの用途でシェルに設定したキーがこのアカウントに届くことはなく、サブスクリプションのログインがこのアカウントのキーより優先されることもありません。
3. 指定があれば`ANTHROPIC_BASE_URL`と`ANTHROPIC_MODEL`を設定し、キーを`ANTHROPIC_API_KEY`または`ANTHROPIC_AUTH_TOKEN`として設定します。

サブスクリプションのログインを使うアカウントは、これまでどおり引き継いだ環境をそのまま使います。対話型セッションでは、環境にあるキーを使うかどうかを、Claude Codeが一度尋ねることがあります。

## 安全性と特殊なケース

- キーが`config.toml`、キャッシュ、ログに書き込まれることはなく、`ccam --dry-run gw`もキーを読みません。ドライランでは`command:`のソースは実行されません。
- `ccam add --api`は`env:`のソースをその場で確認し、変数が空なら警告します。`command:`のソースは起動時に初めて実行されます。
- APIキーのアカウントにはログインがないので、`ccam login`と`ccam logout`は拒否します。`ccam doctor`はエンドポイントとキーのソースを表示し、`env:`の変数が空なら警告します。
- `ccam usage`では`API`と表示されます。APIキーには比べるべき5時間や週間の上限がないので、`ccam next --auto`、`--fallback`、`ccam pick`はサブスクリプションのアカウントの中からだけ選びます。
- `env`ブロックや`apiKeyHelper`を含む共有の`settings.json`は、APIキーのアカウントも含め、グループのすべてのアカウントに適用されます。プロバイダーの設定は共有の設定に入れないでください。

## 関連項目

- [アカウントの管理](/features/accounts)
- [Claude Codeの起動](/features/launch)
- [コマンドリファレンス](/docs/commands)
