---
title: Копирование существующей настройки
description: Копируйте настройки, CLAUDE.md, расширения, MCP-серверы и входы MCP из существующей ~/.claude в группу общего доступа ccam командой ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

Если вы пользовались Claude Code до ccam, ваша настройка находится в `~/.claude` и `~/.claude.json`. `ccam migrate` копирует её в группу общего доступа, поэтому каждый аккаунт, который делит эти элементы через группу, начинает с уже имеющейся у вас настройки. Исходная папка только читается, а обычный `claude` работает как прежде. Выполните команду один раз после первого `ccam add` и снова каждый раз, когда аккаунт начинает делить MCP-серверы.

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
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Пробел (или `x`) переключает элемент под курсором, клавиши со стрелками или `j`/`k` перемещают курсор, Enter показывает план, а Esc, `q` или Ctrl+C выводят `Cancelled.`. Элементы, для которых в источнике ничего нет, помечены как `not found`. Если снять все флажки, ccam выводит `Nothing selected.` и завершается.

Затем ccam выводит план и спрашивает подтверждение перед копированием:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Группа доходит только до тех аккаунтов, которые делят через неё элемент, поэтому итог перечисляет их для каждого элемента. Для аккаунтов, которые этого не делают, ccam выводит команды `ccam share <name> --<item> on`, включающие общий доступ.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `[folder]` | Папка Claude Code, из которой копировать | `~/.claude` |
| `--settings` | Копировать `settings.json` | выкл. |
| `--memory` | Копировать глобальный `CLAUDE.md` | выкл. |
| `--extensions` | Копировать `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | выкл. |
| `--mcp` | Копировать пользовательские MCP-серверы (`mcpServers` в `.claude.json`) | выкл. |
| `--history` | Копировать историю разговоров (`projects/`, `history.jsonl`) | выкл. |
| `-g`, `--group string` | Группа общего доступа, в которую копировать | `default` |
| `-y`, `--yes` | Копировать без вопросов; при конфликте сохраняется копия группы | выкл. |
| `--dry-run` (глобальный) | Вывести план и остановиться до любых вопросов и копирования | выкл. |

Без флагов элементов в терминале выбор делается в списке с флажками. Вне терминала или с `--yes` ccam копирует всё, кроме истории разговоров. История по умолчанию не включается, потому что `projects/` может быть большой.

`--group on` означает `default`. `--group off` и `--group local` отклоняются с сообщением `--group takes a group name, not "off"`. `ccam migrate` учитывает `--dry-run`, и список с флажками появляется и в пробном запуске, чтобы можно было заранее посмотреть результат выбора.

## Как это работает

`internal/cli/migrate.go` выбирает элементы, затем `migrate.Build` читает источник и группу и возвращает план. Ничего не записывается до `Plan.Apply`, то есть до вашего подтверждения.

1. **Поиск источника.** `checkSource` отклоняет путь, который отсутствует, не является папкой или не содержит ни одной известной записи Claude Code. `migrate.ClaudeJSON` находит `.claude.json`, относящийся к папке: внутри неё или, для `~/.claude` без собственного файла внутри, `~/.claude.json` рядом с ней. Именно там его хранит Claude Code, когда `CLAUDE_CONFIG_DIR` не задана.
2. **Сопоставление элементов с путями.** Каждый путь берётся из `share.Layout`, того же списка разрешённых путей, что использует `ccam share`. MCP-серверы читаются из `.claude.json` через `mcp.LocalServers`, который использует сканер `profile.TopLevelValue`, а не разбирает весь файл.
3. **Копирование, а не перемещение.** Файлы попадают в `~/.local/share/ccam/groups/<group>/` через `fsutil.AtomicWrite` или `fsutil.CopyFile`; MCP-серверы попадают в `mcp.json` группы.

Каждый элемент сопоставляет один исходный путь одному пути в группе:

<!-- motion:migrate-items -->

Две строки истории показаны пунктиром, потому что они копируются только с `--history`. `.credentials.json` и `.claude.json` перечислены в `share.Never` и никогда не попадают в группу: `Build` возвращает внутреннюю ошибку, если один из них вдруг окажется в структуре, а обход папки пропускает их по имени.

`buildEntry` назначает каждому пути состояние, которое и описывает каждая строка плана:

<!-- motion:migrate-states -->

- **Missing**: нет в источнике или на месте ожидаемой папки находится файл. Указывается один раз для элемента как `not found, skipped`.
- **Linked**: исходный путь уже указывает на копию группы, например когда источником является папка профиля ccam.
- **Same**: в группе уже есть идентичное содержимое.
- **New**: в группе этого нет, поэтому копируется целиком.
- **Merge**: папка есть с обеих сторон. `copyMissing` добавляет только файлы, которых нет в группе; отличающиеся файлы остаются в версии группы, и план сообщает `group copies kept`. Для `plugins/` реестры `installed_plugins.json` и `known_marketplaces.json` объединяются, а абсолютные пути плагинов переписываются так, чтобы указывать на папку группы.
- **Conflict**: файл есть с обеих сторон, но с разным содержимым. Вы выбираете для каждого файла: `keep group`, `use mine`, а также `merge JSON` для файлов `.json` (глубокое объединение, при котором в спорных значениях побеждает группа) или `merge by line` для `history.jsonl`. MCP-сервер с отличающимся определением тоже считается конфликтом, с вариантами `keep group` или `use mine`. Первый вариант в списке выбирается по умолчанию при нажатии Enter; `--yes` всегда сохраняет копию группы.

Сам список с флажками реализован в `tui.MigratePicker` на основе компонента `checklist`, который использует `ccam add`; `migrate.Found` помечает отсутствующие элементы с помощью вызовов lstat и списка MCP-серверов.

### Входы MCP

MCP-серверы, которые входят через OAuth, хранят свои токены в файле входа, `.credentials.json` в разделе `mcpOAuth`, а не в списке серверов. Claude Code использует в качестве ключа каждой записи `<name>|` плюс первые 16 символов хэша `sha256` от полей сервера `type`, `url` и `headers`. Ключ не зависит от папки, поэтому скопированная запись работает в другом аккаунте, если определение сервера идентично. Поэтому `buildLogins` в `internal/migrate/logins.go` копирует только записи серверов, которые в итоге оказываются в группе в определении из источника: новых, уже совпадающих или конфликтных, для которых вы выбрали `use mine`. Входы серверов плагинов (ключи, начинающиеся с `plugin:`) копируются, когда вы копируете расширения и в источнике есть `plugins/`.

<!-- motion:migrate-signins -->

Записи попадают в собственный `.credentials.json` каждого получающего аккаунта, записываемый через `fsutil.AtomicWrite` с правами `0600`. В группу они никогда не попадают. Аккаунт получает их, только если делит `mcp` (или, для серверов плагинов, `extensions`) через целевую группу, и добавляются только ключи, которых у него нет: уже имеющиеся записи, возможно более новые, сохраняются. Строка плана `sign-ins` перечисляет серверы, получающие аккаунты и пропущенные аккаунты. ccam пропускает аккаунт и называет его с причиной, если аккаунт:

- зарегистрирован на месте через `ccam import`, поскольку обычный `claude` может использовать его файл входа;
- запущен, поскольку Claude Code перезаписывает его файл входа (проверяется ещё раз непосредственно перед записью, сообщается как `started meanwhile`);
- ещё не вошёл;
- хранит вход в Связке ключей macOS, в которую ccam не записывает.

В пропущенном аккаунте войдите в серверы через `/mcp`. Аккаунт, который начнёт делить MCP-серверы позже, получит входы, когда вы снова выполните `ccam migrate`; итог сообщает об этом строкой ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Безопасность и особые случаи

- Источник только читается. Ни один файл в нём не перемещается, не связывается и не изменяется, а вход источника никогда не записывается, в том числе в macOS, где он может читаться из Связки ключей.
- Перед заменой файла группы, через `use mine`, объединение или новые MCP-серверы в существующем `mcp.json`, копия группы сохраняется в `<group>/.ccam-backup/<time>/`. Имя папки уникально для каждого запуска, поэтому два запуска в одну и ту же секунду не перезаписывают резервные копии друг друга. Итог выводит путь.
- С `--yes` при каждом конфликте сохраняется копия группы, и ccam подсказывает, как выбрать иначе: выполнить ту же команду без `--yes`.
- Если ничего не изменилось бы, ccam выводит `Nothing to copy.`; если при каждом конфликте сохраняется копия группы, выводится `Nothing to copy: the group keeps its own copies.`
- Если копирование прервалось на середине, ccam останавливается с сообщением `The source folder was not changed. Group files replaced so far are saved in` и путём к папке резервной копии.
- Отсутствующий источник завершается сообщением `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## См. также

- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Проверка и исправление](/features/doctor)
- [Несколько аккаунтов Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
