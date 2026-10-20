---
title: Remote Control с ccam
description: "Запуск Remote Control от имени любого аккаунта через ccam <name> rc: действующие параметры, зачем ccam предлагает копировать общие MCP-серверы, решение проблем."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc` запускает Remote Control из Claude Code от имени аккаунта `work`, и вы можете управлять сессиями на этом компьютере из claude.ai/code или мобильного приложения Claude. Команда работает как `claude rc` с одним отличием: если аккаунт делит MCP-серверы, ccam один раз спрашивает, скопировать ли их в `.claude.json` аккаунта, потому что сессии Remote Control больше нигде серверы не ищут.

## Запуск

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Короткая команда](/features/shortcuts) работает так же: `ccwork rc` равносильна `ccam work rc`. Выполните команду в папке, где хотите работать, затем откройте claude.ai/code или приложение Claude, войдя в тот же аккаунт Claude, что и `work`.

## Параметры

Здесь встречаются параметры двух видов: собственные параметры ccam, которые ставятся перед именем аккаунта или задаются в настройках ccam, и параметры Remote Control, которые идут после `rc` и передаются в Claude Code без изменений.

### Параметры и настройки ccam

| Параметр или настройка | Действие | По умолчанию |
| --- | --- | --- |
| `--dry-run` (перед аккаунтом) | Вывести окружение и команду, а также уведомление о MCP, если оно есть. Ничего не спрашивает и не копирует | выкл. |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Задаётся, когда аккаунт делит MCP-серверы, а в его `.claude.json` нет текущего списка группы. `y` копирует их, `N` или Enter запускает без них | `N` |
| `mcp_strategy = "sync"` в `~/.config/ccam/config.toml` | Копировать серверы группы в `.claude.json` перед каждым запуском, с Remote Control или без, чтобы вопрос никогда не задавался | `"flag"` |
| `ccam share work --mcp off` | Перестать делить MCP-серверы с аккаунтом; тогда Remote Control запускается только с собственными серверами аккаунта и ничего не спрашивает | вкл. |

### Параметры Remote Control

Всё, что идёт после `rc`, принадлежит Claude Code. Ниже перечислены параметры `claude remote-control` в Claude Code 2.1.304; `ccam work rc --help` выводит список для вашей версии.

| Параметр | Действие |
| --- | --- |
| `--name <name>` | Имя сессии, которое показывается в claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Префикс для автоматически создаваемых имён сессий (по умолчанию имя хоста) |
| `-c`, `--continue` | Снова подключиться к сессии, которую Remote Control последней записал для этой папки, в пределах примерно 4 часов |
| `--session-id <id>` | Снова подключиться к определённой сессии |
| `--permission-mode <mode>` | Режим разрешений для новых сессий: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` или `plan` |
| `--spawn <mode>` | `same-dir` (по умолчанию), `worktree` (отдельный git worktree для каждой сессии) или `session` (одна сессия, завершается вместе с ней) |
| `--capacity <N>` | Наибольшее число одновременных сессий в режимах `same-dir` и `worktree` (по умолчанию 32) |
| `--[no-]create-session-in-dir` | Создать первую сессию в текущей папке при запуске (по умолчанию вкл.) |
| `--[no-]chrome` | Claude in Chrome для новых сессий (по умолчанию выкл.) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Отладочный и подробный вывод |

Ставьте `rc` сразу после имени аккаунта. ccam считает командой Claude Code первый аргумент, который не начинается с `-`, поэтому в `ccam work rc --name laptop` он видит `rc`, а в `ccam work --model opus rc` видит `opus`.

## Что делает ccam

Запуск проходит обычный путь ([Запуск Claude Code](/features/launch)) с одним изменением на шаге MCP:

1. ccam находит аккаунт, проверяет его ссылки общего доступа и находит `claude`, как при любом запуске.
2. Он видит, что команда — `rc` или `remote-control`, и не добавляет `--mcp-config`.
3. Если аккаунт делит MCP-серверы, ccam сравнивает `mcp.json` группы с `sync-stamp` в `~/.local/state/ccam/mcp/work/`. Если `.claude.json` уже содержит текущие серверы или у группы их нет, он сразу переходит к шагу 6.
4. В противном случае он называет серверы, которых не будет в сессиях:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. При ответе `y` он сохраняет резервную копию `.claude.json`, берёт собственную блокировку Claude Code, добавляет серверы группы в `mcpServers` и сохраняет остальные серверы аккаунта. При ответе `N`, без терминала или пока аккаунт запущен, он ничего не меняет и сообщает об этом:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Он регистрирует сессию, указывает в `CLAUDE_CONFIG_DIR` папку аккаунта и заменяет себя на `claude rc` через `exec`.

После этого Remote Control продолжает работать как сервер. Каждая сессия, которую он для вас запускает, является дочерним процессом этого сервера и поэтому работает как `work`: с входом, настройками и историей разговоров `work`, а расход учитывается в лимите использования `work`. Пока сервер работает, ccam считает аккаунт запущенным.

Пробный запуск показывает разницу. Обычный запуск передаёт серверы группы аргументом, а Remote Control — нет:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Почему так устроено

**Remote Control не принимает `--mcp-config`.** При обычных запусках ccam передаёт серверы группы в Claude Code как `--mcp-config=<path>` и никогда не записывает файлы Claude Code. `claude remote-control` отклоняет этот параметр перед именем своей команды и не принимает никаких после него, поэтому до ccam 0.6.0 `ccam work rc` завершался ошибкой для каждого аккаунта, который делил MCP-серверы. Теперь ccam этот параметр не добавляет.

**Его сессии читают серверы только из `.claude.json`.** Сессии, которые запускает Remote Control, не наследуют параметры командной строки. Единственное место, где они находят пользовательские MCP-серверы, — `.claude.json` аккаунта. В этом файле также хранятся данные входа аккаунта и состояние по папкам, поэтому ccam не может связать его с общей копией, как делает это с `settings.json`. Скопировать в него серверы — единственный способ дать сессиям общий список.

**ccam спрашивает перед записью.** Со стратегией по умолчанию (флаг) ccam обещает не трогать `.claude.json`. Remote Control — единственный случай, когда ему приходится туда писать, поэтому он сначала спрашивает, сохраняет резервную копию и не трогает ничего, кроме `mcpServers`.

**Он не пишет, пока аккаунт запущен.** Claude Code перезаписывает `.claude.json` во время сессии, и изменение, сделанное снаружи, потерялось бы или было бы перезаписано. Поэтому ccam пропускает копирование, когда запущена любая сессия аккаунта, в том числе другой сервер Remote Control.

**Он спрашивает один раз на каждое изменение.** После копирования `sync-stamp` хранит время изменения `mcp.json` группы. ccam спрашивает снова, только когда этот файл меняется, например после `ccam mcp add`. Если вы ответили `N`, копии ещё нет, поэтому вопрос задаётся при каждом запуске Remote Control.

**`ccam mcp` поддерживает копии в актуальном состоянии.** При совпадении имён запись в `.claude.json` побеждает `--mcp-config`, поэтому старая копия скрыла бы более новое определение и при обычных запусках. Поэтому `ccam mcp add`, `add-json`, `remove` и `import` обновляют копии во всех аккаунтах, которые не запущены:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Без терминала ответ — «нет».** Remote Control часто запускают из скрипта, службы или через `nohup`, где ответить некому. Тогда ccam запускается без серверов, а не ждёт ответа.

## Проблемы и решения

### В сессиях Remote Control нет общих MCP-серверов

На вопрос ответили `N`, терминала не было или аккаунт был запущен. Остановите Remote Control, закройте другие сессии аккаунта и снова выполните `ccam work rc` в терминале; ответьте `y`. Чтобы вопрос никогда не задавался, задайте `mcp_strategy = "sync"` в `~/.config/ccam/config.toml`. Серверы уже запущенной сессии не меняются во время её работы, поэтому после копирования начните новую сессию.

### "Account work is running, so its .claude.json is left alone"

Сессия `work` всё ещё открыта: обычный `ccam work`, короткая команда или другой сервер Remote Control. Завершите её и снова выполните `ccam work rc`. Сессии, запущенные другим аккаунтом, не учитываются.

### ccam спрашивает каждый раз

Вы ещё не скопировали серверы. Один раз ответьте `y` или задайте `mcp_strategy = "sync"`. Если серверы в Remote Control вам вообще не нужны, перестаньте делить MCP с этим аккаунтом: `ccam share work --mcp off`.

### `ccam work rc` сразу завершается ошибкой про `--mcp-config`

Ваша версия ccam старше 0.6.0 и передаёт `--mcp-config` в Remote Control. Выполните `ccam update`. Та же ошибка в актуальной версии ccam означает, что перед `rc` стоял параметр со значением, как в `ccam work --model opus rc`, и ccam принял `opus` за команду; ставьте `rc` сразу после имени аккаунта.

### Изменение, сделанное через `ccam mcp`, не доходит до Remote Control

`ccam mcp add` и `remove` обновляют копии только в аккаунтах, которые не запущены, и сообщают об этом: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Остановите Remote Control и запустите его снова; ccam увидит изменившийся `mcp.json` и спросит ещё раз.

### Копия заменила мой собственный сервер

Если у аккаунта был собственный сервер с тем же именем, что и сервер группы, копия его заменила. Прежний `.claude.json` лежит в `~/.local/state/ccam/mcp/work/backups/`, где хранятся пять самых новых копий. Чтобы ваша версия была у всех аккаунтов, перенесите её в группу командой `ccam mcp import work <server> --force` до следующего запуска Remote Control.

### Я хочу удалить скопированные серверы

`ccam mcp remove <server>` удаляет сервер из группы и из всех копий. Чтобы вернуть состояние до копирования, восстановите самую новую резервную копию из `~/.local/state/ccam/mcp/work/backups/`, когда ни одна сессия аккаунта не запущена.

### Remote Control спрашивает, доверять ли папке, или завершается без терминала

Доверие к папкам хранится для каждого аккаунта отдельно в `.claude.json`, поэтому папка, которой вы доверяли как `personal`, для `work` новая. Один раз выполните `ccam work` в этой папке и подтвердите доверие, затем запускайте там Remote Control.

### "You must be logged in with a Claude account that has a subscription"

Для Remote Control нужен вход с подпиской Claude. Проверьте email аккаунта в `ccam list` и при необходимости войдите заново командой `ccam login work`.

### Сессия не появляется на телефоне

В приложении Claude или на claude.ai/code выполнен вход в другой аккаунт Claude. Сравните его с email, который `ccam list` показывает для `work`.

### ccam предупреждает, что другой Claude Code делит историю разговоров в этой папке

В той же папке запущен другой аккаунт из той же группы истории. Сессии, продолжающие один и тот же разговор из обоих аккаунтов, могут перемешать историю. Отвечайте `y`, только если они будут работать с разными разговорами, или запустите Remote Control в другой папке.

## См. также

- [Запуск Claude Code](/features/launch)
- [Общие MCP-серверы](/features/mcp)
- [Общий доступ между аккаунтами](/features/share)
- [Управление аккаунтами](/features/accounts)
