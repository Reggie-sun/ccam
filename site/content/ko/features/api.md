---
title: API 키 계정
description: ccam add --api는 Anthropic API나 게이트웨이용 API 키 계정을 추가합니다. 키는 실행할 때 변수나 명령에서 읽으며 저장하지 않습니다.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API 키 계정은 Claude 구독 로그인 대신 API 키로 Claude Code를 시작합니다. Anthropic API를 쓰거나, LiteLLM이나 OpenRouter 게이트웨이처럼 같은 Messages 형식을 쓰는 엔드포인트를 쓸 수 있습니다. ccam은 키를 저장하지 않습니다. 계정을 시작할 때마다 환경 변수나 비밀번호 관리자 같은 명령에서 키를 읽고, 그 Claude Code 프로세스에만 넘깁니다.

## 사용법

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

API 키 계정도 다른 계정처럼 설정, `CLAUDE.md`, 확장, MCP 서버를 공유하며, 단축 명령, `ccam folder`, `ccam ext`와 함께 쓸 수 있습니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--api` | 로그인하는 대신 API 키 계정을 추가합니다 | 꺼짐 |
| `--key-from env:NAME` | 환경 변수 `NAME`에서 키를 읽습니다 | |
| `--key-from command:COMMAND` | `COMMAND`를 실행하고(`sh -c`, Windows에서는 `cmd /C`) 출력한 첫 줄을 씁니다 | |
| `--base-url <url>` | 엔드포인트(`ANTHROPIC_BASE_URL`) | Anthropic API |
| `--auth api-key\|bearer` | 키를 `ANTHROPIC_API_KEY`(`x-api-key`)로 보내거나, `ANTHROPIC_AUTH_TOKEN`(`Authorization: Bearer`, 게이트웨이는 보통 이쪽을 요구)으로 보냅니다 | `api-key` |
| `--model <model>` | 모델(`ANTHROPIC_MODEL`) | Claude Code의 기본값 |

설정은 `config.toml`에 저장됩니다.

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

## 동작 방식

계정을 시작할 때마다, 평소의 검사를 마친 뒤 Claude Code가 시작되기 전에 다음을 합니다.

1. ccam이 키를 읽습니다. `env:`는 환경에서, `command:`는 30초 제한으로 명령을 실행해서 읽습니다. 명령은 터미널에서 비밀번호를 물을 수 있고, 오류 출력도 터미널로 갑니다. 키가 비어 있으면 키를 어디서 찾으려 했는지 알리는 메시지와 함께 시작을 멈춥니다.
2. Claude Code의 로그인 방식이나 사용할 제공자를 정하는, 물려받은 변수를 모두 지웁니다. `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`입니다. 다른 용도로 셸에 설정한 키가 이 계정에 닿을 수 없고, 구독 로그인이 이 계정의 키보다 우선할 수도 없습니다.
3. 지정했다면 `ANTHROPIC_BASE_URL`과 `ANTHROPIC_MODEL`을 설정하고, 키를 `ANTHROPIC_API_KEY`나 `ANTHROPIC_AUTH_TOKEN`으로 설정합니다.

구독 로그인을 쓰는 계정은 지금까지처럼 물려받은 환경을 그대로 씁니다. 대화형 세션에서는 Claude Code가 환경에서 찾은 키를 쓸지 한 번 물어볼 수 있습니다.

## 안전장치와 예외 상황

- 키는 `config.toml`, 캐시, 로그 어디에도 기록되지 않으며, `ccam --dry-run gw`는 키를 읽지 않습니다. 드라이 런에서는 `command:` 소스를 실행하지 않습니다.
- `ccam add --api`는 `env:` 소스를 바로 확인하고 변수가 비어 있으면 경고합니다. `command:` 소스는 시작할 때 처음 실행됩니다.
- API 키 계정에는 로그인이 없으므로 `ccam login`과 `ccam logout`은 거부합니다. `ccam doctor`는 엔드포인트와 키 소스를 보여 주고, `env:` 변수가 비어 있으면 경고합니다.
- `ccam usage`는 이 계정을 `API`로 표시합니다. API 키에는 비교할 5시간·주간 한도가 없으므로, `ccam next --auto`, `--fallback`, `ccam pick`은 구독 계정 중에서만 고릅니다.
- `env` 블록이나 `apiKeyHelper`가 들어 있는 공유 `settings.json`은 API 키 계정을 포함해 그룹의 모든 계정에 적용됩니다. 제공자 설정은 공유 설정에 넣지 마세요.

## 관련 문서

- [계정 관리](/features/accounts)
- [Claude Code 시작하기](/features/launch)
- [명령어 레퍼런스](/docs/commands)
