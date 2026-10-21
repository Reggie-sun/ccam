---
title: Справочник команд
metaTitle: Справочник команд ccam
description: Все команды и параметры ccam, от ccam add и ccam <name> до ccam usage, ccam next и ccam doctor.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

Поставьте `--dry-run` перед командой, чтобы увидеть, что она сделает, ничего не выполняя, например `ccam --dry-run work`. Это работает при запуске Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) и с `group`, `migrate`, `shortcut`, `update` и `doctor --fix`. Другие команды, которые что-то меняют, отклоняют этот флаг и не выполняются.

<!-- commands -->

## Запуск Claude Code

`ccam <name|number> [arguments]` запускает Claude Code от имени этого аккаунта. ccam не читает аргументы после имени; он передаёт их в `claude` в том же порядке, а стандартный ввод, код возврата и сигналы работают точно так же, как с самим `claude`. Указывать `--` после имени необязательно.

`ccam` с аргументами, но без имени аккаунта, например `ccam -p "summarise"`, использует аккаунт, заданный для текущей папки, а если его нет — аккаунт по умолчанию. Задать их можно командами `ccam folder <name> [folder]` и `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` перед именем аккаунта предназначен для запусков с `-p`/`--print`. Когда аккаунт достигает лимита использования, ccam снова выполняет ту же команду от имени аккаунта с входом и наибольшим запасом и выводит только результат завершившегося запуска. См. [Аккаунты по папкам](/features/folder) и [Повторный запуск при лимите](/features/next#run-again-at-the-limit).

Короткая команда, например `ccwork`, запускает свой аккаунт так же: `ccwork --resume` делает то же, что и `ccam work --resume`. См. раздел [короткие команды](/docs/commands#shortcuts).

## Добавление и удаление аккаунтов

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` без указания папки ищет `~/.claude` и `~/.claude-*` и спрашивает про каждую найденную. Импортированные папки регистрируются там, где они находятся; ccam никогда не перемещает и не копирует данные входа.

`ccam add --api` добавляет аккаунт, который использует ключ API вместо входа Claude; ccam читает ключ при каждом запуске и никогда его не сохраняет. См. [Аккаунты с ключом API](/features/api).

`ccam rename`, `ccam remove` и команды `ccam group` не изменяют аккаунт, пока Claude Code работает от его имени. Сначала завершите сессию.

`ccam list --json` выводит все аккаунты с номером (слотом), короткой командой, email и папкой каждого. Массив `groups` перечисляет группы аккаунта, а объект `share` сопоставляет каждому элементу (`mcp`, `settings`, `memory`, `extensions`, `history`) группу, из которой он берётся, или `local`.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Группа делит элементы между своими аккаунтами; аккаунт может состоять в нескольких группах, но каждый элемент берётся только из одной группы. Подробнее в разделе [общий доступ через группы](/features/share).

## Копирование существующей настройки

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` копирует `settings.json`, `CLAUDE.md`, расширения и пользовательские MCP-серверы из папки Claude Code в группу, чтобы их получил каждый аккаунт, который делит эти элементы. Папка только читается. Параметры `--settings`, `--memory`, `--extensions`, `--mcp` и `--history` копируют только указанные элементы. Без них в терминале появляется список с флажками, где отмечено всё, кроме истории разговоров; в остальных случаях копируется всё, кроме истории разговоров. `-g` или `--group` выбирает группу, а `-y` или `--yes` копирует без вопросов и оставляет копию группы, если файл отличается. Кроме того, он копирует сохранённые данные входа серверов в собственный файл входа каждого аккаунта. См. раздел [копирование существующей настройки](/docs/sharing#copy-an-existing-setup).

## MCP-серверы

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` переносит в группу аккаунта серверы, которые есть только у этого аккаунта. Сервер, чьё имя в группе уже занято другим определением, пропускается; добавьте `--force`, чтобы заменить копию в группе, — ccam сначала сохранит её резервную копию.

## Использование и переключение

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` и `ccam pick` выбирают одинаково: сначала аккаунты ниже `auto_max_usage` (`config.toml`, по умолчанию 100), по стратегии `auto_strategy` (`headroom`, `expiring` или `balance`), с учётом недельного окна `--model`, никогда не аккаунт на 100% и в пределах организации аккаунта, с которого они начинают (`--cross-org` или `cross_org = true`, чтобы искать дальше). Когда выбрать некого, они завершаются с кодом 75 или, с `--wait`, ждут сброса первого аккаунта. `ccam next --continue` переносит разговор даже между аккаунтами без общей истории: тогда Claude Code копирует его в новый аккаунт. В отличие от `ccam <name>`, `ccam next` сам разбирает свои опции, поэтому опции Claude Code указываются после `--`; без него ccam останавливается с `unknown flag`. Внутри Claude Code команда плагина `/ccam:switch [account]` передаёт разговор дальше, а `ccam hook` — то, что запускают его хуки. См. [Переключение на следующий аккаунт](/features/next), [Выбор аккаунта для скриптов](/features/pick) и [Фоновые сессии всех аккаунтов](/features/agents).

## Обслуживание

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

`ccam update` проверяет контрольную сумму SHA-256 выпуска, прежде чем заменить ccam, а `--version v0.2.0` устанавливает определённый выпуск. ccam, установленный из пакета `.deb`, обновляется через apt. Подробности в разделе [обновление](/docs/install#update).

## Файлы

| Путь | Содержимое |
| --- | --- |
| `~/.config/ccam/config.toml` | Аккаунты, группы и настройки (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` каждого аккаунта |
| `~/.local/share/ccam/groups/<group>/` | Общие настройки, MCP-серверы, расширения и история |
| `~/.local/state/ccam/` | Кэш использования, время последнего использования и заметки, которые оставляют `/ccam:switch` и хук лимита (`handoff/`) |
| `~/.local/bin/<shortcut>` | Короткие команды, например `ccwork`; ссылки на ccam |

Переменные `XDG_CONFIG_HOME`, `XDG_DATA_HOME` и `XDG_STATE_HOME` позволяют перенести эти папки. В Windows `~` — это `%USERPROFILE%`, а короткие команды — это `<shortcut>.exe`.
