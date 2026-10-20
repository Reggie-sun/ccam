---
title: API 金鑰帳號
description: ccam add --api 為 Anthropic 的 API 或閘道新增 API 金鑰帳號；金鑰在啟動時從變數或指令讀取，絕不儲存。
updated: 2026-10-20
order: 13
translation: machine
sourceHash: "a903c92f2c50"
---

API 金鑰帳號會以 API 金鑰而不是 Claude 訂閱登入來啟動 Claude Code：可以是 Anthropic 的 API，也可以是支援相同 Messages 格式的端點，例如 LiteLLM 或 OpenRouter 閘道。ccam 不會儲存金鑰。每次啟動該帳號時，它會從環境變數或密碼管理工具之類的指令讀取金鑰，並只交給那個 Claude Code 行程。

## 用法

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

API 金鑰帳號和其他帳號一樣加入群組（預設為 `default`：設定、`CLAUDE.md`、擴充功能和 MCP 伺服器），也能搭配快捷指令、`ccam folder` 和 `ccam ext` 使用。

## 選項

| 選項 | 作用 | 預設值 |
| --- | --- | --- |
| `--api` | 新增 API 金鑰帳號，而不是登入 | 關閉 |
| `--key-from env:NAME` | 從環境變數 `NAME` 讀取金鑰 | |
| `--key-from command:COMMAND` | 執行 `COMMAND`（使用 `sh -c`，在 Windows 上使用 `cmd /C`），並使用它印出的第一行 | |
| `--base-url <url>` | 端點（`ANTHROPIC_BASE_URL`） | Anthropic 的 API |
| `--auth api-key\|bearer` | 以 `ANTHROPIC_API_KEY`（`x-api-key`）或 `ANTHROPIC_AUTH_TOKEN`（`Authorization: Bearer`，閘道通常需要這種方式）傳送金鑰 | `api-key` |
| `--model <model>` | 模型（`ANTHROPIC_MODEL`） | Claude Code 的預設值 |

這些設定儲存在 `config.toml` 中：

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

## 運作方式

每次啟動該帳號時，在一般的檢查之後、Claude Code 啟動之前：

1. ccam 讀取金鑰：`env:` 從環境讀取，`command:` 則以 30 秒的時間限制執行指令。指令可以在終端機上要求輸入密碼；它的錯誤輸出也會送到終端機。金鑰為空時，啟動會停止，並顯示應該從哪裡取得金鑰。
2. 它會移除所有繼承而來、決定 Claude Code 登入方式或使用哪個供應商的變數：金鑰和權杖（`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`、`CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`）、端點（`ANTHROPIC_BASE_URL`、`ANTHROPIC_UNIX_SOCKET`、`ANTHROPIC_MODEL`），以及供應商切換開關（`CLAUDE_CODE_USE_BEDROCK`、`CLAUDE_CODE_USE_VERTEX`、`CLAUDE_CODE_USE_FOUNDRY`、`CLAUDE_CODE_USE_ANTHROPIC_AWS`、`CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`、`CLAUDE_CODE_USE_MANTLE`、`CLAUDE_CODE_USE_GATEWAY`）。為了其他用途而設定在 Shell 中的金鑰不會傳到這個帳號，訂閱登入也不會蓋過它的金鑰。
3. 有指定時，它會設定 `ANTHROPIC_BASE_URL` 和 `ANTHROPIC_MODEL`，並把金鑰設為 `ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN`。

使用訂閱登入的帳號和以前一樣，保留它們繼承的環境。在互動式工作階段中，Claude Code 可能會詢問一次是否使用它在環境中找到的金鑰。

## 安全性和特殊情況

- 金鑰絕不會寫入 `config.toml`、快取或記錄檔，而且 `ccam --dry-run gw` 不會讀取它：試執行不會執行 `command:` 來源。
- `ccam add --api` 會立即檢查 `env:` 來源，並在變數為空時發出警告；`command:` 來源則在啟動時才第一次執行。
- `ccam login` 和 `ccam logout` 會拒絕 API 金鑰帳號，因為它沒有登入。`ccam doctor` 會顯示它的端點和金鑰來源，並在 `env:` 變數為空時發出警告。
- `ccam usage` 會為它顯示 `API`。`ccam next --auto`、`--fallback` 和 `ccam pick` 只在訂閱帳號中選擇，因為 API 金鑰沒有可以比較的 5 小時或每週限制。
- 含有 `env` 區塊或 `apiKeyHelper` 的共用 `settings.json` 會套用到群組中的每個帳號，包括 API 金鑰帳號。請不要把供應商設定放在共用的設定中，或把 API 金鑰帳號放到不含 `settings` 的群組。

## 相關內容

- [管理帳號](/features/accounts)
- [啟動 Claude Code](/features/launch)
- [指令參考](/docs/commands)
