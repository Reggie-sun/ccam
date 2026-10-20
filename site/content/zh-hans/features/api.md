---
title: API 密钥账号
description: ccam add --api 为 Anthropic 的 API 或网关添加 API 密钥账号；密钥在启动时从变量或命令中读取，从不保存。
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API 密钥账号用 API 密钥而不是 Claude 订阅登录来启动 Claude Code：可以是 Anthropic 的 API，也可以是使用相同 Messages 格式的接口，例如 LiteLLM 或 OpenRouter 网关。ccam 不保存密钥。每次启动该账号时，它都会从环境变量或某个命令（例如密码管理器）中读取密钥，并只把它交给那个 Claude Code 进程。

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
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

API 密钥账号和其他账号一样共享设置、`CLAUDE.md`、扩展和 MCP 服务器，也可以配合快捷命令、`ccam folder` 和 `ccam ext` 使用。

## 选项

| 选项 | 作用 | 默认值 |
| --- | --- | --- |
| `--api` | 添加 API 密钥账号，而不是登录 | 关闭 |
| `--key-from env:NAME` | 从环境变量 `NAME` 读取密钥 | |
| `--key-from command:COMMAND` | 运行 `COMMAND`（用 `sh -c`，在 Windows 上用 `cmd /C`），使用它输出的第一行 | |
| `--base-url <url>` | 接口地址（`ANTHROPIC_BASE_URL`） | Anthropic 的 API |
| `--auth api-key\|bearer` | 把密钥作为 `ANTHROPIC_API_KEY`（`x-api-key`）发送，或作为 `ANTHROPIC_AUTH_TOKEN`（`Authorization: Bearer`，网关通常需要这种）发送 | `api-key` |
| `--model <model>` | 模型（`ANTHROPIC_MODEL`） | Claude Code 的默认值 |

这些设置保存在 `config.toml` 中：

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

## 工作原理

每次启动该账号时，在常规检查之后、Claude Code 启动之前：

1. ccam 读取密钥：`env:` 从环境中读取，`command:` 通过运行该命令读取，限时 30 秒。命令可以在终端上询问密码，它的错误输出也会显示在终端上。如果密钥为空，启动会停止，并显示一条说明预期在哪里找到密钥的消息。
2. 它会删除所有继承来的、决定 Claude Code 如何登录或使用哪个提供商的变量：`ANTHROPIC_API_KEY`、`ANTHROPIC_AUTH_TOKEN`、`CLAUDE_CODE_OAUTH_TOKEN`、`ANTHROPIC_BASE_URL`、`ANTHROPIC_MODEL`、`CLAUDE_CODE_USE_BEDROCK`、`CLAUDE_CODE_USE_VERTEX` 和 `CLAUDE_CODE_USE_FOUNDRY`。你在 shell 中为其他用途设置的密钥无法到达这个账号，订阅登录也无法压过它的密钥。
3. 如果给出了 `ANTHROPIC_BASE_URL` 和 `ANTHROPIC_MODEL`，就设置它们，并把密钥设置为 `ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN`。

使用订阅登录的账号和以前一样，保留它们继承的环境。在交互式会话中，Claude Code 可能会询问一次是否使用它在环境中找到的密钥。

## 安全性与边界情况

- 密钥从不写入 `config.toml`、缓存或日志，`ccam --dry-run gw` 也不会读取它：试运行不会运行 `command:` 来源。
- `ccam add --api` 会立即检查 `env:` 来源，如果该变量为空就发出警告；`command:` 来源在启动时才第一次运行。
- API 密钥账号没有登录，因此 `ccam login` 和 `ccam logout` 会拒绝它。`ccam doctor` 会显示它的接口地址和密钥来源，并在 `env:` 变量为空时发出警告。
- `ccam usage` 对它显示 `API`。`ccam next --auto`、`--fallback` 和 `ccam pick` 只在订阅账号中选择，因为 API 密钥没有可供比较的 5 小时或每周限制。
- 带有 `env` 块或 `apiKeyHelper` 的共享 `settings.json` 会作用于组中的每个账号，包括 API 密钥账号。请不要把提供商设置放进共享设置。

## 相关内容

- [管理账号](/features/accounts)
- [启动 Claude Code](/features/launch)
- [命令参考](/docs/commands)
