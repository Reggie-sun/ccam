---
title: Облікові записи з ключем API
description: ccam add --api додає обліковий запис із ключем API для API Anthropic або шлюзу; ключ читається під час запуску зі змінної чи команди й не зберігається.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Обліковий запис із ключем API запускає Claude Code з ключем API замість входу за підпискою Claude: для API Anthropic або для кінцевої точки, яка розуміє той самий формат Messages, як-от шлюз LiteLLM або OpenRouter. ccam не зберігає ключ. Він читає його під час кожного запуску облікового запису зі змінної середовища або з команди, наприклад менеджера паролів, і передає його лише цьому процесу Claude Code.

## Використання

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

Обліковий запис із ключем API спільно використовує налаштування, `CLAUDE.md`, розширення та сервери MCP, як будь-який інший обліковий запис, і працює з короткими командами, `ccam folder` і `ccam ext`.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `--api` | Додати обліковий запис із ключем API замість входу | вимкнено |
| `--key-from env:NAME` | Читати ключ зі змінної середовища `NAME` | |
| `--key-from command:COMMAND` | Виконати `COMMAND` (через `sh -c` або `cmd /C` у Windows) і використати перший виведений рядок | |
| `--base-url <url>` | Кінцева точка (`ANTHROPIC_BASE_URL`) | API Anthropic |
| `--auth api-key\|bearer` | Передавати ключ як `ANTHROPIC_API_KEY` (`x-api-key`) або як `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, що зазвичай потрібно шлюзам) | `api-key` |
| `--model <model>` | Модель (`ANTHROPIC_MODEL`) | модель Claude Code за замовчуванням |

Налаштування зберігаються в `config.toml`:

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

## Як це працює

Під час кожного запуску облікового запису, після звичайних перевірок і до старту Claude Code:

1. ccam читає ключ: `env:` із середовища, `command:` виконанням команди з обмеженням у 30 секунд. Команда може запитати пароль у терміналі; її потік помилок теж іде в термінал. Порожній ключ зупиняє запуск із повідомленням про те, де очікувався ключ.
2. Він видаляє всі успадковані змінні, які визначають, як Claude Code виконує вхід або якого провайдера використовує: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` і `CLAUDE_CODE_USE_FOUNDRY`. Ключ, заданий у вашій оболонці для чогось іншого, не потрапить до цього облікового запису, а вхід за підпискою не візьме гору над його ключем.
3. Він задає `ANTHROPIC_BASE_URL` і `ANTHROPIC_MODEL`, якщо їх указано, а ключ — як `ANTHROPIC_API_KEY` або `ANTHROPIC_AUTH_TOKEN`.

Облікові записи зі входом за підпискою, як і раніше, отримують успадковане середовище. В інтерактивному сеансі Claude Code може один раз запитати, чи використовувати ключ, знайдений у середовищі.

## Безпека та граничні випадки

- Ключ ніколи не записується в `config.toml`, кеш або журнал, а `ccam --dry-run gw` його не читає: джерело `command:` під час пробного запуску не виконується.
- `ccam add --api` одразу перевіряє джерело `env:` і попереджає, якщо змінна порожня; джерело `command:` уперше виконується під час запуску.
- `ccam login` і `ccam logout` відмовляються працювати з обліковим записом із ключем API, бо він не має входу. `ccam doctor` показує його кінцеву точку та джерело ключа й попереджає, коли змінна `env:` порожня.
- `ccam usage` показує для нього `API`. `ccam next --auto`, `--fallback` і `ccam pick` вибирають лише серед облікових записів із підпискою, бо ключ API не має 5-годинного чи тижневого ліміту для порівняння.
- Спільний `settings.json` із блоком `env` або `apiKeyHelper` застосовується до кожного облікового запису групи, зокрема до облікових записів із ключем API. Не тримайте налаштування провайдера в спільних налаштуваннях.

## Пов'язані сторінки

- [Керування обліковими записами](/features/accounts)
- [Запуск Claude Code](/features/launch)
- [Довідник команд](/docs/commands)
