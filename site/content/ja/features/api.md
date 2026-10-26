---
title: APIキーのアカウント
description: ccam add --apiは、AnthropicのAPIやゲートウェイ用のAPIキーのアカウントを追加します。キーは起動時に変数かコマンドから読み、保存しません。
updated: 2026-10-26
order: 13
translation: machine
sourceHash: 771d749f1f90
---

APIキーのアカウントは、Claudeのサブスクリプションのログインの代わりに、APIキーでClaude Codeを起動します。キーは、AnthropicのAPI用でも、LiteLLMやOpenRouterのゲートウェイのように同じMessages形式を話す任意のエンドポイント用でもかまいません。ccamはキーを保存しません。アカウントを起動するたびに、環境変数か、パスワードマネージャーなどのコマンドからキーを読み、そのClaude Codeのプロセスだけに渡します。

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
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

APIキーのアカウントも、ほかのアカウントと同じようにグループに加わり（デフォルトは`default`グループで、設定、`CLAUDE.md`、拡張機能、MCPサーバーを共有します）、ショートカット、`ccam folder`、`ccam ext`でも使えます。

## オプション

| オプション | 効果 | デフォルト |
| --- | --- | --- |
| `--api` | ログインする代わりにAPIキーのアカウントを追加します | オフ |
| `--key-from env:NAME` | 環境変数`NAME`からキーを読みます | |
| `--key-from command:COMMAND` | `COMMAND`を実行し（`sh -c`、Windowsでは`cmd /d /s /c "…"`に入力したとおりに、引用符も残して）、出力の最初の行を使います | |
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

1. ccamがキーを読みます。`env:`は環境から、`command:`はコマンドを30秒のタイムアウトで実行して読みます。コマンドはターミナルでパスワードを尋ねることができ、エラー出力もターミナルに出ます。キーが空の場合は、キーをどこで探したかを示すメッセージを出して起動を中止します。
2. Claude Codeのサインイン方法や使うプロバイダーを決める、引き継いだ変数をすべて削除します。
   - キーとトークン: `ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`、`CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`
   - エンドポイント: `ANTHROPIC_BASE_URL`、`ANTHROPIC_UNIX_SOCKET`、`ANTHROPIC_MODEL`
   - プロバイダーの切り替え: `CLAUDE_CODE_USE_BEDROCK`、`CLAUDE_CODE_USE_VERTEX`、`CLAUDE_CODE_USE_FOUNDRY`、`CLAUDE_CODE_USE_ANTHROPIC_AWS`、`CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`、`CLAUDE_CODE_USE_MANTLE`、`CLAUDE_CODE_USE_GATEWAY`

   ほかの用途でシェルに設定したキーがこのアカウントに届くことはなく、サブスクリプションのログインがこのアカウントのキーより優先されることもありません。
3. 指定があれば`ANTHROPIC_BASE_URL`と`ANTHROPIC_MODEL`を設定し、キーを`ANTHROPIC_API_KEY`または`ANTHROPIC_AUTH_TOKEN`として設定します。

サブスクリプションのログインを使うアカウントは、これまでどおり引き継いだ環境をそのまま使います。対話型セッションでは、環境にあるキーを使うかどうかを、Claude Codeが一度尋ねることがあります。

## 安全性と特殊なケース

- キーが`config.toml`、キャッシュ、ログに書き込まれることはなく、`ccam --dry-run gw`もキーを読みません。ドライランでは`command:`のソースは実行されません。代わりにドライランは、エンドポイント、モデル、キーの取得元を表示します。

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  起動時に取り除かれる、引き継いだ変数は、先頭に`env -u <name>`として表示されます。`env:`のソースの変数が空の場合、ドライランは起動と同じように失敗します。
- `ccam add --api`は`env:`のソースをその場で確認し、変数が空なら警告します。`command:`のソースは起動時まで実行されません。
- APIキーのアカウントにはログインがないので、`ccam login`と`ccam logout`は拒否します。`ccam doctor`はエンドポイントとキーのソースを表示し、`env:`の変数が空なら警告します。
- `ccam usage`では`no limits (API)`と表示されます。APIキーには比べるべき5時間や週間の上限がないので、`ccam next --auto`、`--fallback`、`ccam pick`はサブスクリプションのアカウントの中からだけ選びます。
- `env`ブロックや`apiKeyHelper`を含む共有の`settings.json`は、APIキーのアカウントも含め、グループのすべてのアカウントに適用されます。プロバイダーの設定は共有の設定に入れないか、APIキーのアカウントを`settings`を共有しないグループに入れてください。

## 関連項目

- [アカウントの管理](/features/accounts)
- [Claude Codeの起動](/features/launch)
- [コマンドリファレンス](/docs/commands)
