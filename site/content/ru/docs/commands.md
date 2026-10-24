---
title: Справочник команд
metaTitle: Справочник команд ccam
description: Все команды и параметры ccam, от ccam add и ccam <name> до ccam usage, ccam next и ccam doctor.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "36af7e209cd9"
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

**Места в тарифах Team и Enterprise.** При автоматическом выборе (`ccam next`, `--fallback`, `ccam pick`, клавиша `n` на панели) ccam переключается и между местами одной и той же организации Claude Team или Enterprise. В справочном центре Anthropic ([тарифы Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); аналогично для [тарифов Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) сказано: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (каждое место в тарифе Team закреплено за одним конкретным человеком, и пользоваться им может только этот человек) и «Users can't route their work through another user's account to use that user's remaining usage» (пользователи не могут направлять свою работу через аккаунт другого пользователя, чтобы расходовать оставшийся лимит этого пользователя). Поэтому использование автоматического выбора между местами может нарушать условия Anthropic. Прежде чем на него полагаться, уточните это у своей организации и у Anthropic. О личных тарифах (Pro, Max) см. [примечание об условиях в FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Короткая команда, например `ccwork`, запускает свой аккаунт так же: `ccwork --resume` делает то же, что и `ccam work --resume`. См. раздел [короткие команды](/docs/commands#shortcuts).

## Добавление и удаление аккаунтов

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` без указания папки ищет `~/.claude` и `~/.claude-*` и спрашивает про каждую найденную; о папке, которую не удалось зарегистрировать, сообщается как `Skipped <path>: <reason>`, а остальные регистрируются, и для папки, чьё имя уже занято аккаунтом, нужен `--name`, который указывается вместе с папкой. Импортированные папки регистрируются там, где они находятся; ccam никогда не перемещает и не копирует данные входа.

`ccam login --expired` также выполняет вход в аккаунты, чей вход очистил Claude Code. `ccam remove` выполняет выход из управляемого аккаунта, прежде чем его папка попадёт в корзину (кроме случая, когда тот же вход есть у другого аккаунта), а `ccam remove` импортированного аккаунта выводит его из групп с копией того, что они делают общим. В macOS `ccam rename` сохраняет имя папки, поскольку Связка ключей хранит вход под путём папки.

`ccam add --api` добавляет аккаунт, который использует ключ API вместо входа Claude; ccam читает ключ при каждом запуске и никогда его не сохраняет. См. [Аккаунты с ключом API](/features/api).

`ccam rename`, `ccam remove` и команды `ccam group` не изменяют аккаунт, пока Claude Code работает от его имени, как бы он ни был запущен (в том числе фоновые сессии). Сначала завершите сессию. Они также отказываются работать, пока изменение групп не завершено: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` выводит все аккаунты с номером (слотом), короткой командой, email и папкой каждого. Массив `groups` перечисляет группы аккаунта, а объект `share` сопоставляет каждому элементу (`mcp`, `settings`, `memory`, `extensions`, `history`) группу, из которой он берётся, или `local`.

## Короткие команды

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Короткая команда представляет собой команду с выбранным вами именем, которая запускает Claude Code от имени одного аккаунта. Как и `ccam work`, она передаёт все аргументы в Claude Code без изменений, поэтому `ccwork --resume` и `ccwork -p "hi"` работают так, как вы ожидаете. Это ссылка на ccam в `~/.local/bin`.

Короткие команды необязательны. Мастер `ccam add` спрашивает, создать ли её, `ccam add <name> --shortcut <command>` создаёт её без вопросов, а `ccam list` показывает короткую команду каждого аккаунта или `-`, если её нет.

В именах допустимы символы a-z, 0-9, `.`, `-` и `_`. Имя не может совпадать с командой ccam, с короткой командой другого аккаунта или с именем другой программы в вашем `PATH`, а также с файлом в `~/.local/bin`, который создал не ccam, например с вашей собственной копией ccam или ссылкой на него. Если короткая команда, отредактированная вручную в `config.toml`, нарушает эти правила, равна `claude` или `ccam` или назначена двум аккаунтам, `config.toml` не загружается, и ccam называет аккаунт.

`ccam rename` сохраняет имя короткой команды, а `ccam remove` удаляет её вместе с аккаунтом. `ccam shims install` заново создаёт ссылки всех аккаунтов, например после того как вы переместили ccam; файл, мешающий одной короткой команде, не останавливает остальные.

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

`ccam migrate` копирует `settings.json`, `CLAUDE.md`, расширения и пользовательские MCP-серверы из папки Claude Code в группу, чтобы их получил каждый аккаунт, который делит эти элементы. Папка только читается. Параметры `--settings`, `--memory`, `--extensions`, `--mcp` и `--history` копируют только указанные элементы. Без них в терминале появляется список с флажками, где отмечено всё, кроме истории разговоров; в остальных случаях копируется всё, кроме истории разговоров. `-g` или `--group` выбирает группу, а `-y` или `--yes` копирует без вопросов и оставляет копию группы, если файл отличается. Данные входа MCP не копируются, поскольку скопированный вход выводит из системы других его владельцев; план называет аккаунты, которые входят через `/mcp`. `--mcp` работает и тогда, когда `~/.claude` пуст или отсутствует, но в `~/.claude.json` есть серверы. См. раздел [копирование существующей настройки](/docs/sharing#copy-an-existing-setup).

## MCP-серверы

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` переносит в группу аккаунта серверы, которые есть только у этого аккаунта, в том числе серверы только для одной папки (обычный `claude mcp add`), которые после этого загружаются во всех папках. Сервер, чьё имя в группе уже занято другим определением, пропускается; добавьте `--force`, чтобы заменить копию в группе, — ccam сначала сохранит её резервную копию.

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
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` и `ccam pick` выбирают одинаково: сначала аккаунты ниже `auto_max_usage` (`config.toml`, по умолчанию 100), по стратегии `auto_strategy` (`headroom`, `expiring` или `balance`), с учётом недельного окна `--model`, никогда не аккаунт на 100% (окно, которое Claude Code сообщил заполненным, считается заполненным на 100% до своего сброса), никогда не аккаунт со входом по ключу API и в пределах организации аккаунта, с которого они начинают (`--cross-org` или `cross_org = true`, чтобы искать дальше). `--wait` принимает значение после `=` (`--wait=2h`); если записать его через пробел, он отклоняется с подсказкой. Когда выбрать некого, они завершаются с кодом 75 или, с `--wait`, ждут сброса первого аккаунта. `ccam next --continue` переносит разговор даже между аккаунтами без общей истории: тогда Claude Code копирует его в новый аккаунт. В отличие от `ccam <name>`, `ccam next` сам разбирает свои опции, поэтому опции Claude Code указываются после `--`; без него ccam останавливается с `unknown flag`. Внутри Claude Code команда плагина `/ccam:switch [account]` передаёт разговор дальше, а `ccam hook` — то, что запускают его хуки. См. [Переключение на следующий аккаунт](/features/next), [Выбор аккаунта для скриптов](/features/pick) и [Фоновые сессии всех аккаунтов](/features/agents).

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
| `~/.config/ccam/config.toml` | Аккаунты, группы и настройки (`theme`, `usage_source`, `relay`, `cross_org` …). Настройки пишутся в начале, над первым `[[group]]`, `[[profile]]` или `[[folder]]`: настройка под таблицей относится к этой таблице, поэтому ccam называет её при каждом запуске и в `ccam doctor`, а команды, которые сохраняют конфигурацию, отказываются работать, пока её не перенесут. Комментарии не сохраняются, когда ccam записывает файл |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` каждого аккаунта |
| `~/.local/share/ccam/groups/<group>/` | Общие настройки, MCP-серверы, расширения и история |
| `~/.local/state/ccam/` | Кэш использования и записи о лимите (`usage/<name>.json`, `usage/<name>.limit.json`), время последнего использования, метки сессий и блокировки (`locks/`), короткие команды, созданные ccam (`shortcuts.json`), незавершённое изменение групп (`group-change.json`) и заметки, которые оставляют `/ccam:switch` и хук лимита (`handoff/<account>/<session>.json`) |
| `~/.local/bin/<shortcut>` | Короткие команды, например `ccwork`; ссылки на ccam |

Переменные `XDG_CONFIG_HOME`, `XDG_DATA_HOME` и `XDG_STATE_HOME` позволяют перенести эти папки. В Windows `~` — это `%USERPROFILE%`, а короткие команды — это `<shortcut>.exe`.
