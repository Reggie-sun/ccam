---
title: Копирование существующей настройки
description: Копируйте настройки, CLAUDE.md, расширения и MCP-серверы из существующей ~/.claude в группу ccam командой ccam migrate.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: "f835fbe84ade"
---

Если вы пользовались Claude Code до ccam, ваша настройка находится в `~/.claude` и `~/.claude.json`. `ccam migrate` копирует её в группу общего доступа, поэтому каждый аккаунт, который делит эти элементы через группу, начинает с уже имеющейся у вас настройки. Исходная папка только читается, а обычный `claude` работает как прежде. Выполните команду один раз после первого `ccam add`; если в `~/.claude.json` есть MCP-серверы, первый `ccam add` называет их в своей подсказке.

## Использование

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

В терминале, без флагов элементов и без `--yes`, ccam сначала показывает список с флажками:

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Пробел (или `x`) переключает элемент под курсором, клавиши со стрелками или `j`/`k` перемещают курсор, Enter показывает план, а Esc, `q` или Ctrl+C выводят `Cancelled.`. Элементы, для которых в источнике ничего нет, помечены как `not found`. Если снять все флажки, ccam выводит `Nothing selected.` и завершается.

Затем ccam выводит план и спрашивает подтверждение перед копированием:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Группа влияет только на свои аккаунты и только на те элементы, которые она делит, поэтому итог перечисляет аккаунты группы. Если группа не делит какой-то из скопированных элементов или некоторые аккаунты в неё не входят, ccam выводит команды `ccam group set <group> --add <items>` и `ccam group join <group> <accounts>`, которые это исправляют.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `[folder]` | Папка Claude Code, из которой копировать | `~/.claude` |
| `--settings` | Копировать `settings.json` | выкл. |
| `--memory` | Копировать глобальный `CLAUDE.md` | выкл. |
| `--extensions` | Копировать `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | выкл. |
| `--mcp` | Копировать пользовательские MCP-серверы (`mcpServers` в `.claude.json`) | выкл. |
| `--history` | Копировать историю разговоров (`projects/`, `history.jsonl`) | выкл. |
| `-g`, `--group <group>` | Группа общего доступа, в которую копировать | `default` |
| `-y`, `--yes` | Копировать без вопросов; при конфликте сохраняется копия группы | выкл. |
| `--dry-run` (глобальный) | Вывести план и остановиться до любых вопросов и копирования | выкл. |

В терминале, если вы не указали флаги элементов, элементы выбираются в списке с флажками. Вне терминала или с `--yes` ccam копирует всё, кроме истории разговоров. История по умолчанию не включается, потому что `projects/` может быть большой.

`--group on` означает `default`. `--group off` и `--group local` отклоняются с сообщением `--group takes a group name, not "off"`. `ccam migrate` поддерживает `--dry-run`, и список с флажками появляется и в пробном запуске, чтобы можно было заранее посмотреть результат выбора.

## Как это работает

ccam строит полный план, прежде чем что-либо копировать. Он выбирает элементы (`internal/cli/migrate.go`), затем читает источник и группу и строит план (`migrate.Build`). Ничего не записывается до `Plan.Apply`, то есть до вашего подтверждения.

1. **Поиск источника.** ccam отклоняет путь, который отсутствует, не является папкой или не содержит ни одной известной записи Claude Code (`checkSource`). Для одних только MCP-серверов пустая или отсутствующая `~/.claude` допустима, если в `~/.claude.json` есть серверы. ccam находит `.claude.json`, относящийся к папке (`migrate.ClaudeJSON`): внутри неё или, для `~/.claude` без собственного файла внутри, `~/.claude.json` рядом с ней. Именно там его хранит Claude Code, когда `CLAUDE_CONFIG_DIR` не задана.
2. **Сопоставление элементов с путями.** Каждый путь берётся из `share.Layout`, того же списка разрешённых путей, что используют группы. MCP-серверы читаются из `.claude.json` через `mcp.LocalServers`, который использует сканер `profile.TopLevelValue`, а не разбирает весь файл.
3. **Копирование, а не перемещение.** Файлы попадают в `~/.local/share/ccam/groups/<group>/` через `fsutil.AtomicWrite` или `fsutil.CopyFile`, который записывает временный файл рядом с целевым и переименовывает его, поэтому прерванное копирование никогда не оставляет недописанный файл. MCP-серверы попадают в `mcp.json` группы под его блокировкой: сервер, который за это время добавил другой процесс ccam, сохраняется, а копии, сделанные для [Remote Control](/guides/remote-control), затем обновляются.

Каждый элемент сопоставляет один исходный путь одному пути в группе:

<!-- motion:migrate-items -->

Две строки истории показаны пунктиром, потому что они копируются только с `--history`. `.credentials.json` и `.claude.json` перечислены в `share.Never` и никогда не попадают в группу: ccam останавливается с внутренней ошибкой, если один из них вдруг окажется в структуре (`Build`), а обход папки пропускает их по имени.

ccam назначает каждому пути состояние, которое и описывает каждая строка плана (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: нет в источнике или на месте ожидаемой папки находится файл. Указывается один раз для элемента как `not found, skipped`.
- **Linked**: исходный путь уже указывает на копию группы, например когда источником является папка профиля ccam.
- **Same**: в группе уже есть идентичное содержимое.
- **New**: в группе этого нет, поэтому копируется целиком.
- **Merge**: папка есть с обеих сторон. ccam добавляет только файлы, которых нет в группе (`copyMissing`); отличающиеся файлы остаются в версии группы, и план сообщает `group copies kept`. Для `plugins/` реестры `installed_plugins.json` и `known_marketplaces.json` объединяются, а абсолютные пути плагинов переписываются так, чтобы указывать на папку группы.
- **Conflict**: файл есть с обеих сторон, но с разным содержимым. Вы выбираете для каждого файла: `keep group`, `use mine`, а также `merge JSON` для файлов `.json` (глубокое объединение, при котором в спорных значениях побеждает группа) или `merge by line` для `history.jsonl`. MCP-сервер с отличающимся определением тоже считается конфликтом, с вариантами `keep group` или `use mine`. Первый вариант в списке выбирается по умолчанию при нажатии Enter; `--yes` всегда сохраняет копию группы.

Сам список с флажками реализован в `tui.MigratePicker` на основе компонента `checklist`, который использует `ccam add`; ccam помечает отсутствующие элементы с помощью вызовов lstat и списка MCP-серверов (`migrate.Found`).

### Входы MCP

MCP-серверы, которые входят через OAuth, хранят свои токены в файле входа, `.credentials.json` в разделе `mcpOAuth`, а не в списке серверов. ccam их не копирует. Скопированный вход содержит тот же токен обновления, что и оригинал, а сервер меняет этот токен при первом обновлении: тогда все остальные владельцы, включая исходную `~/.claude`, выходят из этого сервера. Поэтому для каждого копируемого сервера, в который вошёл источник, строка плана `sign-ins` перечисляет аккаунты группы, которым нужно самим войти через `/mcp`, а итог повторяет это. Запись без токена не считается входом и не указывается, как и аккаунт, у которого уже есть собственный вход, или сервер, определение которого группа сохраняет. Ни один файл входа не записывается.

## Безопасность и особые случаи

- Источник только читается. Ни один файл в нём не перемещается, не связывается и не изменяется, а вход источника никогда не записывается, в том числе в macOS, где он может читаться из Связки ключей.
- Перед заменой файла группы, через `use mine`, объединение или новые MCP-серверы в существующем `mcp.json`, копия группы сохраняется в `<group>/.ccam-backup/<time>/`. Имя папки уникально для каждого запуска, поэтому два запуска в одну и ту же секунду не перезаписывают резервные копии друг друга. Итог выводит путь.
- С `--yes` при каждом конфликте сохраняется копия группы, и ccam подсказывает, как выбрать иначе: выполнить ту же команду без `--yes`.
- Если ничего не изменилось бы, ccam выводит `Nothing to copy.`; если при каждом конфликте сохраняется копия группы, выводится `Nothing to copy: the group keeps its own copies.`
- Если копирование прервалось на середине, ccam останавливается с сообщением `The source folder was not changed. Group files replaced so far are saved in` и путём к папке резервной копии. Повторный запуск копирует файлы, которые не были скопированы полностью.
- Когда MCP делят несколько групп, команды, которые предлагают `ccam doctor`, `ccam ext` и `ccam mcp list`, выглядят как `ccam migrate --mcp --group <a|b>`.
- Отсутствующий источник завершается сообщением `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## См. также

- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Проверка и исправление](/features/doctor)
- [Несколько аккаунтов Claude Code на одном компьютере](/guides/multiple-claude-code-accounts)
- [Что такое CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
