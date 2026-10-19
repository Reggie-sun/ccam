---
title: Справочник команд
description: Все команды и параметры ccam, от ccam add и ccam <name> до ccam usage, ccam next и ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

Поставьте `--dry-run` перед любой командой, чтобы увидеть, что она сделает, ничего не выполняя, например `ccam --dry-run work`.

<!-- commands -->

## Запуск Claude Code

`ccam <name|number> [arguments]` запускает Claude Code от имени этого аккаунта. ccam не читает аргументы после имени; он передаёт их в `claude` в том же порядке, а стандартный ввод, код возврата и сигналы работают точно так же, как с самим `claude`. Указывать `--` после имени необязательно.

`ccam` с аргументами, но без имени аккаунта, например `ccam -p "summarise"`, использует аккаунт по умолчанию. Задать его можно командой `ccam use <name>`.

Короткая команда, например `ccwork`, запускает свой аккаунт так же: `ccwork --resume` делает то же, что и `ccam work --resume`. См. раздел [короткие команды](/docs/commands#shortcuts).

## Добавление и удаление аккаунтов

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` без указания папки ищет `~/.claude` и `~/.claude-*` и спрашивает про каждую найденную. Импортированные папки регистрируются там, где они находятся; ccam никогда не перемещает и не копирует данные входа.

## Короткие команды

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Короткая команда представляет собой команду с выбранным вами именем, которая запускает Claude Code от имени одного аккаунта. Как и `ccam work`, она передаёт все аргументы в Claude Code без изменений, поэтому `ccwork --resume` и `ccwork -p "hi"` работают так, как вы ожидаете. Это ссылка на ccam в `~/.local/bin`.

Короткие команды необязательны. Мастер `ccam add` спрашивает, создать ли её, `ccam add <name> --shortcut <command>` создаёт её без вопросов, а `ccam list` показывает короткую команду каждого аккаунта или `-`, если её нет.

В именах допустимы символы a-z, 0-9, `.`, `-` и `_`. Имя не может совпадать с командой ccam, с короткой командой другого аккаунта или с именем другой программы в вашем `PATH`.

`ccam rename` сохраняет имя короткой команды, а `ccam remove` удаляет её вместе с аккаунтом. `ccam shims install` заново создаёт ссылки всех аккаунтов, например после того как вы переместили ccam.

## Общий доступ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Каждый элемент принимает `on` (группа `default`), `off` или имя группы. Подробнее в разделе [общий доступ](/docs/sharing).

## Копирование существующей настройки

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` копирует `settings.json`, `CLAUDE.md`, расширения и пользовательские MCP-серверы из папки Claude Code в группу, чтобы их получил каждый аккаунт, который делит эти элементы. Папка только читается. Параметры `--settings`, `--memory`, `--extensions`, `--mcp` и `--history` копируют только указанные элементы; без них копируется всё, кроме истории разговоров. `-g` или `--group` выбирает группу, а `-y` или `--yes` копирует без вопросов и оставляет копию группы, если файл отличается. См. раздел [копирование существующей настройки](/docs/sharing#copy-an-existing-setup).

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
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Обслуживание

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
```

## Файлы

| Путь | Содержимое |
| --- | --- |
| `~/.config/ccam/config.toml` | Аккаунты и то, что они делят |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` каждого аккаунта |
| `~/.local/share/ccam/groups/<group>/` | Общие настройки, MCP-серверы, расширения и история |
| `~/.local/state/ccam/` | Кэш использования и время последнего использования |
| `~/.local/bin/<shortcut>` | Короткие команды, например `ccwork`; ссылки на ccam |

Переменные `XDG_CONFIG_HOME`, `XDG_DATA_HOME` и `XDG_STATE_HOME` позволяют перенести эти папки.
