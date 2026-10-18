---
title: Справочник команд
description: Все команды и параметры ccam, от ccam add и ccs <name> до ccam usage, ccam next и ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` и `ccs` являются одной и той же программой; используйте ту, что вам удобнее. Поставьте `--dry-run` перед любой командой, чтобы увидеть, что она сделает, ничего не выполняя, например `ccs --dry-run work`.

<!-- commands -->

## Запуск Claude Code

`ccs <name|number> [arguments]` запускает Claude Code от имени этого аккаунта. ccam не читает аргументы после имени; он передаёт их в `claude` в том же порядке, а стандартный ввод, код возврата и сигналы работают точно так же, как с самим `claude`. Указывать `--` после имени необязательно.

`ccs` с аргументами, но без имени аккаунта, например `ccs -p "summarise"`, использует аккаунт по умолчанию. Задать его можно командой `ccam use <name>`.

Короткие команды `claude1`, `claude2` и так далее являются ссылками на ccam в `~/.local/bin`. `ccam shims install --names` дополнительно создаёт `claude-<name>` для каждого аккаунта.

## Добавление и удаление аккаунтов

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` без указания папки ищет `~/.claude` и `~/.claude-*` и спрашивает про каждую найденную. Импортированные папки регистрируются там, где они находятся; ccam никогда не перемещает и не копирует данные входа.

## Общий доступ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Каждый элемент принимает `on` (группа `default`), `off` или имя группы. Подробнее в разделе [общий доступ](/docs/sharing).

## MCP-серверы

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Использование и переключение

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Обслуживание

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Файлы

| Путь | Содержимое |
| --- | --- |
| `~/.config/ccam/config.toml` | Аккаунты и то, что они делят |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` каждого аккаунта |
| `~/.local/share/ccam/groups/<group>/` | Общие настройки, MCP-серверы, расширения и история |
| `~/.local/state/ccam/` | Кэш использования и время последнего использования |
| `~/.local/bin/claude1`, `ccs` | Ссылки на ccam |

Переменные `XDG_CONFIG_HOME`, `XDG_DATA_HOME` и `XDG_STATE_HOME` позволяют перенести эти папки.
