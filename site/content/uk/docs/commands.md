---
title: Довідник команд
metaTitle: Довідник команд ccam
description: Усі команди та параметри ccam, від ccam add і ccam <name> до ccam usage, ccam next і ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

Додайте `--dry-run` перед командою, щоб побачити, що вона зробила б, не виконуючи її, наприклад `ccam --dry-run work`. Це працює під час запуску Claude Code (`ccam <name>`, `ccam run`, `ccam next`) і з `share`, `migrate`, `shortcut`, `update` та `doctor --fix`. Інші команди, які щось змінюють, відхиляють цей параметр і не виконуються.

<!-- commands -->

## Запуск Claude Code

`ccam <name|number> [arguments]` запускає Claude Code від імені цього облікового запису. ccam не читає аргументи після назви; він передає їх у `claude` у тому самому порядку, а стандартне введення, код виходу та сигнали поводяться так само, як із самим `claude`. `--` після назви необов'язковий.

`ccam` з аргументами, але без назви облікового запису, наприклад `ccam -p "summarise"`, використовує обліковий запис, заданий для поточної теки, а якщо його немає — обліковий запис за замовчуванням. Задайте їх командами `ccam folder <name> [folder]` і `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` перед назвою облікового запису призначений для запусків із `-p`/`--print`. Коли обліковий запис досягає ліміту використання, ccam знову виконує ту саму команду від імені облікового запису з входом, де найбільше запасу, і виводить лише результат запуску, який завершився. Див. [Облікові записи за теками](/features/folder) і [Повторний запуск при ліміті](/features/next#run-again-at-the-limit).

Коротка команда, наприклад `ccwork`, запускає свій обліковий запис так само: `ccwork --resume` робить те саме, що й `ccam work --resume`. Див. [короткі команди](/docs/commands#shortcuts).

## Додавання та видалення облікових записів

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` без теки шукає `~/.claude` і `~/.claude-*` та запитує про кожну з них. Імпортовані теки реєструються там, де вони є; ccam ніколи не переміщує й не копіює дані входу.

`ccam add --api` додає обліковий запис, який використовує ключ API замість входу Claude; ccam читає ключ під час кожного запуску й ніколи його не зберігає. Див. [Облікові записи з ключем API](/features/api).

`ccam rename`, `ccam remove` і `ccam share` відмовляються змінювати обліковий запис, поки Claude Code працює від його імені. Спершу завершіть сеанс.

`ccam list --json` виводить кожен обліковий запис із його слотом, короткою командою, електронною поштою й текою. Його об'єкт `share` зіставляє кожен елемент (`mcp`, `settings`, `memory`, `extensions`, `history`) з назвою групи або з `local`.

## Короткі команди

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Коротка команда є командою з назвою, яку ви вибираєте самі, і запускає Claude Code від імені одного облікового запису. Як і `ccam work`, вона передає всі аргументи в Claude Code без змін, тож `ccwork --resume` і `ccwork -p "hi"` працюють так, як ви очікуєте. Це посилання на ccam у `~/.local/bin`.

Короткі команди необов'язкові. Майстер `ccam add` запитує, чи створити її, `ccam add <name> --shortcut <command>` створює її без запитань, а `ccam list` показує коротку команду кожного облікового запису або `-`, якщо її немає.

У назвах можна використовувати символи a-z, 0-9, `.`, `-` і `_`. Назва не може збігатися з командою ccam, з короткою командою іншого облікового запису чи з назвою іншої програми у вашому `PATH`.

`ccam rename` зберігає назву короткої команди, а `ccam remove` видаляє її разом з обліковим записом. `ccam shims install` повторно створює посилання всіх облікових записів, наприклад після того, як ви перемістили ccam.

## Спільне використання

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Кожен елемент приймає `on` (група `default`), `off` або назву групи. Див. [спільне використання](/docs/sharing).

## Копіювання наявного налаштування

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` копіює `settings.json`, `CLAUDE.md`, розширення та користувацькі сервери MCP з теки Claude Code до групи, щоб їх отримав кожен обліковий запис, який спільно використовує ці елементи. Тека лише читається. Параметри `--settings`, `--memory`, `--extensions`, `--mcp` і `--history` копіюють лише вказані елементи. Без них у терміналі ви вибираєте елементи в контрольному списку, де позначено все, крім історії розмов; в інших випадках копіюється все, крім історії розмов. `-g` або `--group` вибирає групу, а `-y` або `--yes` копіює без запитань і залишає копію групи, якщо файл відрізняється. Він також копіює збережені входи серверів до власного файлу входу кожного облікового запису. Див. [копіювання наявного налаштування](/docs/sharing#copy-an-existing-setup).

## Сервери MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` переміщує до групи облікового запису сервери, які є лише в ньому. Сервер, назву якого група вже використовує для іншого визначення, пропускається; додайте `--force`, щоб замінити копію групи, а ccam спершу створить її резервну копію.

## Використання та перемикання

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` і `ccam pick` вибирають однаково: спершу облікові записи нижче `auto_max_usage` (`config.toml`, за замовчуванням 100), за стратегією `auto_strategy` (`headroom` або `expiring`), з урахуванням тижневого вікна `--model` і ніколи не обліковий запис на 100%. Коли вибрати нікого, вони завершуються з кодом 75 або, з `--wait`, чекають на скидання першого облікового запису. Див. [Перемикання на наступний обліковий запис](/features/next) і [Вибір облікового запису для скриптів](/features/pick).

## Обслуговування

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` перевіряє контрольну суму SHA-256 випуску, перш ніж замінити ccam, а `--version v0.2.0` встановлює певний випуск. ccam, встановлений з пакета `.deb`, натомість оновлюється через apt. Див. [оновлення](/docs/install#update).

## Файли

| Шлях | Вміст |
| --- | --- |
| `~/.config/ccam/config.toml` | Облікові записи та те, що вони спільно використовують |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` кожного облікового запису |
| `~/.local/share/ccam/groups/<group>/` | Спільні налаштування, сервери MCP, розширення та історія |
| `~/.local/state/ccam/` | Кеш використання та час останнього використання |
| `~/.local/bin/<shortcut>` | Короткі команди, наприклад `ccwork`; посилання на ccam |

Змінні `XDG_CONFIG_HOME`, `XDG_DATA_HOME` і `XDG_STATE_HOME` змінюють розташування цих тек. У Windows `~` — це `%USERPROFILE%`, а короткі команди — це `<shortcut>.exe`.
