---
title: Довідник команд
description: Усі команди та параметри ccam, від ccam add і ccam <name> до ccam usage, ccam next і ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Додайте `--dry-run` перед будь-якою командою, щоб побачити, що вона зробила б, не виконуючи її, наприклад `ccam --dry-run work`.

<!-- commands -->

## Запуск Claude Code

`ccam <name|number> [arguments]` запускає Claude Code від імені цього облікового запису. ccam не читає аргументи після назви; він передає їх у `claude` у тому самому порядку, а стандартне введення, код виходу та сигнали поводяться так само, як із самим `claude`. `--` після назви необов'язковий.

`ccam` з аргументами, але без назви облікового запису, наприклад `ccam -p "summarise"`, використовує обліковий запис за замовчуванням. Задайте його командою `ccam use <name>`.

Короткі команди `claude1`, `claude2` і так далі є посиланнями на ccam у `~/.local/bin`. `ccam shims install --names` також створює `claude-<name>` для кожного облікового запису.

## Додавання та видалення облікових записів

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` без теки шукає `~/.claude` і `~/.claude-*` та запитує про кожну з них. Імпортовані теки реєструються там, де вони є; ccam ніколи не переміщує й не копіює дані входу.

## Спільне використання

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Кожен елемент приймає `on` (група `default`), `off` або назву групи. Див. [спільне використання](/docs/sharing).

## Сервери MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Використання та перемикання

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Обслуговування

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Файли

| Шлях | Вміст |
| --- | --- |
| `~/.config/ccam/config.toml` | Облікові записи та те, що вони спільно використовують |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` кожного облікового запису |
| `~/.local/share/ccam/groups/<group>/` | Спільні налаштування, сервери MCP, розширення та історія |
| `~/.local/state/ccam/` | Кеш використання та час останнього використання |
| `~/.local/bin/claude1`, `claude2`… | Посилання на ccam |

Змінні `XDG_CONFIG_HOME`, `XDG_DATA_HOME` і `XDG_STATE_HOME` змінюють розташування цих тек.
