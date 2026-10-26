---
title: Общие MCP-серверы
description: "Один список MCP-серверов на группу: ccam mcp add, add-json, list, remove и import; он передаётся через --mcp-config или синхронизируется в каждый аккаунт."
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
---

Claude Code хранит пользовательские MCP-серверы в `.claude.json`, том же файле, где лежит вход аккаунта, поэтому этот файл нельзя сделать общим для аккаунтов. ccam хранит общие серверы в файле группы `mcp.json` и передаёт его каждому аккаунту группы при запуске. Используйте `ccam mcp`, чтобы добавить сервер один раз для всех аккаунтов, и `ccam mcp import`, чтобы перенести в остальные аккаунты сервер, добавленный внутри одной сессии.

## Использование

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

Изменение вступает в силу со следующего запуска каждого аккаунта; запущенные сессии сохраняют серверы, с которыми они стартовали:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Под группами `ccam mcp list` называет каждый сервер, который не доходит до остальных аккаунтов:

- `only in account work: pg` для сервера, который аккаунт добавил для себя;
- `only in account side, for /home/me/code/app only: notes` для сервера, добавленного для одной папки (обычный `claude mcp add`, см. ниже);
- `account work has its own …, used in its sessions in place of the group's`, когда собственный сервер аккаунта носит имя сервера группы;
- `.mcp.json` этой папки с аккаунтами, одобрения которых он ждёт (`waiting for approval in`);
- серверы, которые находятся `in no account, only in ~/.claude.json (plain claude's)`, с командой [`ccam migrate --mcp`](/features/migrate), которая их копирует.

Добавление имени, которое уже есть в группе, заменяет его определение, и вместо `added` выводится `updated`. Серверы попадают в аккаунт, только если он состоит в этой группе и группа делит MCP (`ccam group set <group> --add mcp`); если группа этого не делает, `add` сообщает об этом.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `-g`, `--group <group>` | целевая группа для `add`, `add-json` и `remove` | единственная группа, через которую какой-либо аккаунт делит MCP, иначе `default` |
| `-- <command> [args…]` (`add`) | запускать эту команду как stdio-сервер | |
| `--url <url>` (`add`) | URL удалённого сервера (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | транспорт сервера с `--url` | `http` |
| `--header <'Name: value'>` (`add`) | HTTP-заголовок, можно повторять | нет |
| `-e`, `--env <KEY=VALUE>` (`add`) | переменная окружения для stdio-сервера, можно повторять | нет |
| `--force` (`import`) | заменить отличающийся сервер группы с тем же именем; для старого `mcp.json` создаётся резервная копия | выкл. |
| `--keep-local` (`import`) | оставить записи в `.claude.json` аккаунта | выкл. |
| `-h`, `--help` | справка по каждой команде | |
| `--dry-run` (глобальный) | поддерживается только командой `mcp list`; `add`, `add-json`, `remove` и `import` отклоняют его и ничего не меняют | выкл. |

`add` требует либо команду после `--`, либо `--url`; иначе он завершается ошибкой `a command to run (after --) or --url is required`. `--header` применяется к серверам с `--url`, а `--env` к серверам с командой. `add-json` принимает тот же JSON, что и `claude mcp add-json`, и определение должно быть объектом с `command` или `url`. `import` всегда использует группу, через которую аккаунт делит MCP, а `list` показывает все группы, поэтому `-g` на них не влияет.

## Как это работает

Участвующие файлы:

- `~/.local/share/ccam/groups/<group>/mcp.json`: общие серверы, `{"mcpServers": {…}}`, в формате, который читает `--mcp-config`. Он записывается через `fsutil.AtomicWrite` с правами `0600`, потому что записи серверов часто содержат ключи API, и только под блокировкой `<group>/.mcp.json.lock`, поэтому две одновременно запущенные команды `ccam mcp add` сохраняют свои серверы.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: собственный файл Claude Code с серверами аккаунта. ccam читает из него только значение `mcpServers` сканером `profile.TopLevelValue` и никогда не разбирает весь файл на пути запуска.
- `~/.local/state/ccam/mcp/<name>/`: файл `seen`, `sync-stamp` (с хэшем каждой копии, записанной ccam, чтобы отличать копию, которую никто не менял, от отредактированной вами), `sync-checked` (что стратегия синхронизации видела в последний раз), `kept` (серверы, которые аккаунт сохранил, когда вышел из общего доступа к MCP) и резервные копии `.claude.json`.

Есть две стратегии, которые задаются параметром `mcp_strategy` в `~/.config/ccam/config.toml`. По умолчанию используется стратегия с флагом:

<!-- motion:mcp-flag -->

1. ccam находит MCP-группу аккаунта (`launcher.applyMCP`) и одним вызовом `stat` проверяет, что `mcp.json` группы существует, не разбирая его (`mcp.FlagArgs`).
2. ccam ставит `--mcp-config=<path>` сразу после `claude`, перед вашими аргументами, поэтому ваши собственные `--mcp-config` или `--strict-mcp-config` по-прежнему применяются поверх (`Plan.InsertArgs`).
3. ccam проверяет, изменился ли `.claude.json` аккаунта с последнего запуска, и только в этом случае разбирает его серверы (`mcp.NewLocal`).
4. ccam заменяет себя на Claude Code (`syscall.Exec`).

Путь передаётся одним аргументом, `--mcp-config=<path>`, а не двумя. `--mcp-config` в Claude Code принимает несколько значений, поэтому `claude --mcp-config a.json "prompt"` прочитал бы запрос как второй файл конфигурации и завершился бы ошибкой `MCP config file not found`. Форма с одним аргументом этого избегает, поэтому следующий за ней запрос или подкоманда работают. `ccam work mcp list`, `get` и `remove` флаг не получают, поскольку там Claude Code его отклоняет; вместо этого ccam называет серверы группы в стандартном потоке ошибок ([подробнее](/features/launch#shared-mcp-servers)). Claude Code загружает серверы из `--mcp-config` в дополнение к собственным серверам аккаунта; при совпадении имён побеждает запись в `.claude.json`. Стратегия с флагом записывает `.claude.json` только для Remote Control, когда аккаунт выходит из общего доступа к MCP или снова в него вступает ([ниже](#leaving-and-joining)), и для `ccam mcp import`.

`claude remote-control` (`rc`) не принимает `--mcp-config`, а его сессии читают серверы только из `.claude.json`. Для `ccam work rc` ccam не добавляет флаг и, если есть терминал, спрашивает, скопировать ли серверы группы в `.claude.json` с резервной копией, как это делает стратегия `sync`. Он спрашивает, пока копии нет, а затем снова только после изменения `mcp.json` группы или пропажи скопированного сервера ([Remote Control](/features/launch#remote-control)). Поскольку такая копия побеждает `--mcp-config` при совпадении имён, `ccam mcp add`, `add-json`, `remove` и `import` обновляют копии, которые хранятся у других аккаунтов группы:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Строка `✓ Updated` появляется, только если копия изменилась. Копия, созданная более старой версией ccam, обновляется так же, как и копия после `ccam migrate`.

Сервер, добавленный внутри сессии через `claude mcp add`, попадает только в `.claude.json` этого аккаунта. Обычный `claude mcp add` использует локальную область: сервер загружается для этого одного аккаунта в этой одной папке (с `-s user` он становится доступен аккаунту во всех папках). В этот момент ccam его не видит, потому что Claude Code уже заменил процесс ccam, поэтому проверка выполняется при следующем запуске. Файл `seen` хранит время изменения `.claude.json` и имена серверов, о которых уже сообщалось. Если время отличается, ccam сравнивает имена с группой и один раз выводит подсказку для каждого нового сервера, не ожидая ввода, при обеих стратегиях:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Подсказка также называет собственный сервер аккаунта с именем сервера группы, поскольку в сессиях аккаунта побеждает его сервер.

С `mcp_strategy = "sync"` серверы отображаются в `claude mcp list` и `/mcp` как обычные пользовательские серверы, а ccam перед запуском записывает их в `.claude.json` каждого аккаунта:

<!-- motion:mcp-sync -->

1. ccam сравнивает время изменения `mcp.json` группы с `sync-stamp`, а `.claude.json` аккаунта — с тем, что видел в последний раз (`mcp.Sync`). Если ничего не изменилось, больше ничего не происходит. Иначе он также проверяет, что каждый сервер группы всё ещё есть в `.claude.json`, и возвращает пропавшую копию (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Он берёт собственную блокировку Claude Code, папку `.claude.json.lock` рядом с файлом, создаваемую через `mkdir`. Он ждёт до 2 секунд и считает блокировку старше 10 секунд устаревшей, как это делает библиотека блокировок Claude Code.
3. Он сохраняет текущий `.claude.json` в `backups/` в папке состояния и хранит 5 самых новых копий.
4. Он записывает каждый сервер группы в `mcpServers`, удаляет серверы, которые записал ранее и которых больше нет в группе, и атомарно записывает файл с исходными правами, только если что-то изменилось. Он меняет или удаляет только записанные им копии, которые с тех пор никто не менял. Если у аккаунта есть собственный сервер с тем же именем, что и сервер группы, ccam сохраняет сервер аккаунта и один раз сообщает об этом: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Он сохраняет новую отметку с хэшем каждой копии, и Claude Code запускается без `--mcp-config`.

Синхронизация пропускается с сообщением `ccam: account work is already running, skipping MCP sync`, когда запущена сессия аккаунта, как бы ни был запущен Claude Code, потому что Claude Code перезаписывает `.claude.json` во время работы. При неудачной синхронизации выводится `ccam: MCP sync failed (launching anyway): …`, и запуск всё равно выполняется.

### Выход и возвращение

Аккаунт, который перестаёт делить MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`), сохраняет серверы группы: план копирует их в его `.claude.json` с резервной копией, а его собственные серверы с теми же именами остаются. С `--fresh` он не получает ни одного, и план перечисляет серверы, которые перестанут загружаться. Когда он снова вступает в группу с этими серверами, сохранённые копии, которые никто не менял, удаляются, чтобы до него доходили последующие изменения группы.

`ccam mcp import` переносит собственные серверы аккаунта в его группу:

<!-- motion:mcp-import -->

ccam читает серверы аккаунта, собственные и те, что есть у него только для одной папки, и сравнивает каждый с сервером группы с тем же именем через `mcp.SameDef`, который не учитывает порядок ключей и пробелы. Сервер папки становится сервером группы для всех папок (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); одно имя с двумя разными определениями для папок отклоняется. Серверы плагинов не предлагаются, поскольку их приносит плагин. Сервер, которого нет в группе или который есть в ней с тем же определением, копируется. Отличающийся сервер пропускается, потому что определение группы используют все остальные аккаунты:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

С `--force` сервер группы заменяется, а прежний `mcp.json` сначала сохраняется в `.ccam-backup/<time>/mcp.json` внутри папки группы. Затем ccam удаляет скопированные записи из `.claude.json` аккаунта, собственные и для папки, одной записью, под той же блокировкой и после создания резервной копии в папке состояния (`mcp.Remove`). ccam держит `mcp.json` группы заблокированным с момента чтения файла, пока эти записи не удалены, чтобы никто не удалил сервер из группы в промежутке. Это нужно, потому что при совпадении имён запись аккаунта победила бы запись группы. С `--keep-local` записи остаются в `.claude.json` аккаунта.

## Безопасность и особые случаи

- **Входы.** Серверы, которые входят через OAuth, хранят свои токены в `.credentials.json` каждого аккаунта (`mcpOAuth`), а не в списке серверов, поэтому общий сервер не делает общим свой вход. Войдите один раз в каждом аккаунте через `/mcp`. ccam никогда не копирует входы, и [`ccam migrate`](/features/migrate) тоже: скопированный токен обновления выводит из системы всех остальных его владельцев. [`ccam ext`](/features/ext) и `ccam doctor` называют серверы, в которые аккаунт не вошёл, хотя другой аккаунт вошёл, а `ccam next --continue` сообщает, каких из них не хватает новому аккаунту. Группа никогда не получает файл входа.
- **Запущенные аккаунты.** `import` копирует в группу, но оставляет собственные записи запущенного аккаунта и сообщает об этом; выполните его снова после завершения сессии. Синхронизация пропускает запущенные аккаунты. Обе также обнаруживают сессии Claude Code, запущенные без ccam.
- **Копии для Remote Control.** Копия запущенного аккаунта не изменяется: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` С `mcp_strategy = "sync"` ничего обновлять не нужно, потому что каждый запуск выполняет синхронизацию.
- **Повреждённые файлы.** Если `.claude.json` не является корректным JSON, синхронизация оставляет его нетронутым и сообщает об ошибке. Если блокировку не удаётся взять за 2 секунды, ccam сообщает `could not lock … (Claude Code is writing the file)`.
- **Аккаунты, которые не делят MCP.** `import` отказывается работать и предлагает сначала добавить аккаунт в группу, которая делит MCP (`ccam group join <group> side`).
- **Имена групп.** `-g` принимает только имя группы; оно становится папкой внутри `groups/`, поэтому имена вроде `../x` отклоняются.

## См. также

- [Общий доступ через группы](/features/share)
- [Remote Control с ccam](/guides/remote-control)
- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Копирование существующей настройки](/features/migrate)
- [Справочник команд](/docs/commands)
