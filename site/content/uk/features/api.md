---
title: Облікові записи з ключем API
description: ccam add --api додає обліковий запис із ключем API для API Anthropic або шлюзу; ключ читається під час запуску зі змінної чи команди й не зберігається.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: "771d749f1f90"
---

Обліковий запис із ключем API запускає Claude Code з ключем API замість входу за підпискою Claude. Ключ може бути для API Anthropic або для будь-якої кінцевої точки, яка розуміє той самий формат Messages, як-от шлюз LiteLLM або OpenRouter. ccam не зберігає ключ. Він читає його під час кожного запуску облікового запису зі змінної середовища або з команди, наприклад менеджера паролів, і передає його лише цьому процесу Claude Code.

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
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Обліковий запис із ключем API входить до груп, як будь-який інший обліковий запис (за замовчуванням до групи `default`, яка спільно використовує налаштування, `CLAUDE.md`, розширення та сервери MCP), і працює з короткими командами, `ccam folder` і `ccam ext`.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `--api` | Додати обліковий запис із ключем API замість входу | вимкнено |
| `--key-from env:NAME` | Читати ключ зі змінної середовища `NAME` | |
| `--key-from command:COMMAND` | Виконати `COMMAND` (через `sh -c` або, у Windows, так, ніби його введено в `cmd /d /s /c "…"`, зі збереженням лапок) і використати перший виведений рядок | |
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

1. ccam читає ключ: `env:` із середовища, `command:` виконанням команди з тайм-аутом 30 секунд. Команда може запитати пароль у терміналі; її потік помилок теж іде в термінал. Порожній ключ зупиняє запуск із повідомленням про те, де очікувався ключ.
2. Він видаляє всі успадковані змінні, які визначають, як Claude Code виконує вхід або якого провайдера використовує:
   - ключі й токени: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - кінцеву точку: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - перемикачі провайдера: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   Ключ, заданий у вашій оболонці для чогось іншого, не потрапить до цього облікового запису, а вхід за підпискою не матиме переваги над його ключем.
3. Він задає `ANTHROPIC_BASE_URL` і `ANTHROPIC_MODEL`, якщо їх указано, а ключ — як `ANTHROPIC_API_KEY` або `ANTHROPIC_AUTH_TOKEN`.

Облікові записи зі входом за підпискою, як і раніше, отримують успадковане середовище. В інтерактивному сеансі Claude Code може один раз запитати, чи використовувати ключ, знайдений у середовищі.

## Безпека та граничні випадки

- Ключ ніколи не записується в `config.toml`, кеш або журнал, а `ccam --dry-run gw` його не читає: джерело `command:` під час пробного запуску не виконується. Натомість пробний запуск показує кінцеву точку, модель і те, звідки береться ключ:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Успадкована змінна, яку запуск видаляє, показується як `env -u <name>` на початку. Із джерелом `env:`, змінна якого порожня, пробний запуск завершується помилкою так само, як завершився б запуск.
- `ccam add --api` одразу перевіряє джерело `env:` і попереджає, якщо змінна порожня; джерело `command:` не виконується до запуску.
- `ccam login` і `ccam logout` відмовляються працювати з обліковим записом із ключем API, бо він не має входу. `ccam doctor` показує його кінцеву точку та джерело ключа й попереджає, коли змінна `env:` порожня.
- `ccam usage` показує для нього `no limits (API)`. `ccam next --auto`, `--fallback` і `ccam pick` вибирають лише серед облікових записів із підпискою, бо ключ API не має 5-годинного чи тижневого ліміту для порівняння.
- Спільний `settings.json` із блоком `env` або `apiKeyHelper` застосовується до кожного облікового запису групи, зокрема до облікових записів із ключем API. Не тримайте налаштування провайдера в спільних налаштуваннях або додайте обліковий запис із ключем API до групи без `settings`.

## Пов'язані сторінки

- [Керування обліковими записами](/features/accounts)
- [Запуск Claude Code](/features/launch)
- [Довідник команд](/docs/commands)
