---
title: Аккаунты с ключом API
description: ccam add --api добавляет аккаунт с ключом API для API Anthropic или шлюза; ключ читается при запуске из переменной или команды и не сохраняется.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Аккаунт с ключом API запускает Claude Code с ключом API вместо входа по подписке Claude: для API Anthropic или для эндпоинта, который понимает тот же формат Messages, например шлюза LiteLLM или OpenRouter. ccam не хранит ключ. Он читает его при каждом запуске аккаунта из переменной окружения или из команды, например менеджера паролей, и передаёт его только этому процессу Claude Code.

## Использование

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

Аккаунт с ключом API делит настройки, `CLAUDE.md`, расширения и MCP-серверы, как любой другой аккаунт, и работает с короткими командами, `ccam folder` и `ccam ext`.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `--api` | Добавить аккаунт с ключом API вместо входа | выкл. |
| `--key-from env:NAME` | Читать ключ из переменной окружения `NAME` | |
| `--key-from command:COMMAND` | Выполнить `COMMAND` (через `sh -c` или `cmd /C` в Windows) и использовать первую выведенную строку | |
| `--base-url <url>` | Эндпоинт (`ANTHROPIC_BASE_URL`) | API Anthropic |
| `--auth api-key\|bearer` | Передавать ключ как `ANTHROPIC_API_KEY` (`x-api-key`) или как `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, что обычно нужно шлюзам) | `api-key` |
| `--model <model>` | Модель (`ANTHROPIC_MODEL`) | модель Claude Code по умолчанию |

Настройки хранятся в `config.toml`:

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

## Как это работает

При каждом запуске аккаунта, после обычных проверок и до старта Claude Code:

1. ccam читает ключ: `env:` из окружения, `command:` выполнением команды с ограничением в 30 секунд. Команда может запросить пароль в терминале; её поток ошибок тоже идёт в терминал. Пустой ключ останавливает запуск с сообщением о том, где ожидался ключ.
2. Он удаляет все унаследованные переменные, которые определяют, как Claude Code выполняет вход или какого провайдера использует: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` и `CLAUDE_CODE_USE_FOUNDRY`. Ключ, заданный в вашей оболочке для чего-то другого, не попадёт в этот аккаунт, а вход по подписке не возьмёт верх над его ключом.
3. Он задаёт `ANTHROPIC_BASE_URL` и `ANTHROPIC_MODEL`, если они указаны, а ключ — как `ANTHROPIC_API_KEY` или `ANTHROPIC_AUTH_TOKEN`.

Аккаунты со входом по подписке, как и раньше, получают унаследованное окружение. В интерактивной сессии Claude Code может один раз спросить, использовать ли ключ, найденный в окружении.

## Безопасность и особые случаи

- Ключ никогда не записывается в `config.toml`, кэш или журнал, а `ccam --dry-run gw` его не читает: источник `command:` при пробном запуске не выполняется.
- `ccam add --api` сразу проверяет источник `env:` и предупреждает, если переменная пуста; источник `command:` впервые выполняется при запуске.
- `ccam login` и `ccam logout` отказываются работать с аккаунтом с ключом API, поскольку у него нет входа. `ccam doctor` показывает его эндпоинт и источник ключа и предупреждает, если переменная `env:` пуста.
- `ccam usage` показывает для него `API`. `ccam next --auto`, `--fallback` и `ccam pick` выбирают только среди аккаунтов с подпиской, потому что у ключа API нет 5-часового или недельного лимита для сравнения.
- Общий `settings.json` с блоком `env` или `apiKeyHelper` применяется к каждому аккаунту группы, включая аккаунты с ключом API. Не храните настройки провайдера в общих настройках.

## См. также

- [Управление аккаунтами](/features/accounts)
- [Запуск Claude Code](/features/launch)
- [Справочник команд](/docs/commands)
