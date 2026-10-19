---
title: Управление аккаунтами
description: Добавляйте, входите, просматривайте, переименовывайте, удаляйте и импортируйте аккаунты Claude Code командами ccam add, ccam login, ccam list, ccam rename, ccam remove и ccam import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 97cef50f48f4
---

Каждый аккаунт ccam представляет собой отдельную папку конфигурации Claude Code с собственным входом. Эти команды создают аккаунты, выполняют вход, выводят список, переименовывают и удаляют аккаунты, а также регистрируют уже существующие папки.

## Использование

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Везде, где команда принимает `<name|number>`, можно указать имя аккаунта или номер его слота, то есть столбец `#` в выводе `ccam list`.

`ccam add` в терминале без параметров открывает мастер. Он спрашивает имя, затем показывает список общих элементов с флажками (отмечено всё, кроме истории разговоров), затем задаёт вопрос `Create a shortcut command for work? [y/N]` и после ответа `y` спрашивает `Shortcut name [ccwork]:`. Любой из параметров `--share`, `--shortcut` или `--no-login`, а также конвейер вместо терминала отключают мастер; тогда имя обязательно, а короткая команда создаётся, только если её имя задано через `--shortcut`.

`ccam list` показывает слот, имя, короткую команду, email, общие элементы и время последнего использования, а перед аккаунтом по умолчанию ставит `*`:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

С `--json` каждый аккаунт выводится как объект с полями `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (отсутствует, если аккаунт ни разу не запускался) и `share`. `share` сопоставляет каждому элементу (`mcp`, `settings`, `memory`, `extensions`, `history`) имя группы или `"local"`.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `ccam add --no-login` | Создать аккаунт без входа; войти позже можно командой `ccam login <name>` | выкл. |
| `ccam add --share <item=group>` | Сделать элемент общим с группой или указать `item=local`, чтобы оставить его личным. Можно повторять | `mcp`, `settings`, `memory`, `extensions` в `default`; `history` локально |
| `ccam add --shortcut <command>` | Дополнительно создать эту короткую команду в `~/.local/bin` | нет |
| `ccam list --json` | Вывести аккаунты в формате JSON | выкл. |
| `ccam remove -y`, `--yes` | Удалить без вопросов | выкл. |
| `ccam import --name <name>` | Имя аккаунта для указанной папки | выводится из имени папки |
| `ccam import -y`, `--yes` | Импортировать все найденные папки без вопросов | выкл. |
| `--dry-run` | Учитывается командами `ccam list` и `ccam completion`, которые ничего не меняют. Команды `add`, `login`, `logout`, `use`, `rename`, `remove` и `import` отклоняют его с сообщением `` --dry-run is not supported by `ccam add`; nothing was changed `` (с именем команды) | выкл. |

У `ccam login`, `ccam logout`, `ccam use` и `ccam rename` нет собственных параметров. `ccam completion` принимает `bash` или `zsh`; любое другое значение завершается ошибкой `only bash and zsh are supported`.

Имя аккаунта содержит от 1 до 32 символов из `a-z`, `0-9`, `-` и `_` и начинается с буквы или цифры. Имя не может состоять только из цифр, потому что `ccam 2` означает слот 2. Имя также не может совпадать с командой ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), поэтому `ccam <name>` никогда не бывает неоднозначным. Этот список задан в `config.Reserved` в `internal/config`.

## Как это работает

`ccam add work --shortcut ccwork` выполняет следующие шаги в `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` проверяет имя, а `shim.Check` проверяет короткую команду, прежде чем что-либо будет создано.
2. `profile.Add` создаёт `~/.local/share/ccam/profiles/work/` с правами `0700` и назначает аккаунту наименьший свободный номер слота (`Config.NextSlot`), поэтому номер удалённого аккаунта используется повторно. Первый аккаунт становится аккаунтом по умолчанию.
3. `applyShare` связывает новую пустую папку с её группами общего доступа через `share.Build` и затем `Apply`, например `settings.json` с `groups/default/settings.json`. Если это не удалось, новая папка удаляется, потому что оставшаяся папка блокировала бы каждую повторную попытку ошибкой `folder already exists`.
4. Аккаунт сохраняется в `~/.config/ccam/config.toml` через `fsutil.AtomicWrite`: временный файл, `fsync`, затем `rename`, поэтому сбой никогда не оставляет конфигурацию записанной наполовину.
5. `runLogin` запускает `claude auth login` с `CLAUDE_CONFIG_DIR`, указывающим на папку аккаунта. Если вход не удался, аккаунт остаётся, а ccam выводит ошибку и предлагает позже выполнить `ccam login work`.
6. `profile.ReadEmail` берёт `oauthAccount.emailAddress` из `profiles/work/.claude.json`. Этот файл может занимать несколько мегабайт, поэтому `TopLevelValue` ищет один ключ, а не разбирает весь файл.
7. `shim.Install` связывает `~/.local/bin/ccwork` с исполняемым файлом ccam, который при запуске ищет `argv[0]` в `config.toml`.

Вход может занять несколько минут, и за это время другой ccam может сохранить конфигурацию. Поэтому команды, которые ждут, сохраняют через `env.update`: он заново читает `config.toml`, применяет одно изменение (здесь email) и сохраняет. По той же причине `ccam remove` заново читает конфигурацию после своего вопроса.

`ccam login` выполняет шаги 5 и 6 для существующего аккаунта. `ccam logout` запускает `claude auth logout` в аккаунте и удаляет сохранённый email. `ccam use` только перезаписывает `default` в `config.toml`.

### Вход на удалённой машине

Через SSH браузерная часть `claude auth login` может не сработать. Claude Code ждёт завершения входа на сервере обратного вызова по адресу `http://localhost:<port>/callback` на сервере, но браузер, в котором завершается вход, работает на вашем компьютере, где этот адрес никуда не ведёт, и браузер показывает ошибку, например `ERR_CONNECTION_REFUSED`. Адрес при этом всё равно содержит одноразовый код, поэтому ccam может его доставить:

<!-- motion:accounts-relay -->

`runLogin` не отдаёт ваш терминал Claude Code. Он подключает стандартный ввод Claude Code к каналу, а `relayLogin` копирует туда каждую введённую вами строку, поэтому код, вставленный в ответ на приглашение Claude Code `Paste code here`, всё равно доходит. Если строка является адресом обратного вызова, `callbackURL` принимает её, только если хост равен `localhost`, `127.0.0.1` или `[::1]`, порт указан, путь равен `/callback` и присутствуют и `code`, и `state`; `http://` можно опустить. Затем `deliverCallback` запрашивает этот адрес на сервере с ограничением в 10 секунд и без перехода по перенаправлениям, как сделал бы браузер. Запрос попадает на сервер обратного вызова ожидающего `claude auth login` (пунктирная линия на схеме). Claude Code получает код и сохраняет вход, в Linux в `.credentials.json` в папке аккаунта. Вставленная строка никогда не может заставить ccam обратиться к другому хосту.

Перед началом входа ccam выводит:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

При успехе он выводит `→ Passed the sign-in to Claude Code on this machine.`

### Переименование, удаление и импорт

`ccam rename` и `ccam remove` работают с папкой аккаунта и с тем, что ccam хранит под именем аккаунта в других местах. `nameState` перечисляет эти пути: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` и кэши использования `~/.local/state/ccam/usage/<name>.json` и `<name>.statusline.json`. `ccam remove work` обрабатывает их в таком порядке:

<!-- motion:accounts-remove -->

Сначала `ccam remove` отказывается работать, если аккаунт запущен, затем задаёт вопрос. Если аккаунт делит расширения, `RewritePluginPaths` переписывает файлы группы `installed_plugins.json` и `known_marketplaces.json`, потому что Claude Code хранит абсолютные пути к плагинам, которые могут проходить через папку этого аккаунта. `profile.Trash` перемещает папку в `~/.local/share/ccam/trash/<name>-<time>/`, аккаунт удаляется из `config.toml`, пути `nameState` удаляются, чтобы будущий аккаунт с тем же именем их не унаследовал, а ссылка короткой команды удаляется. Чтобы восстановить аккаунт, верните папку из корзины.

`ccam rename` переименовывает папку, переписывает пути плагинов в её `plugins/` со старой папки на новую, сохраняет статус аккаунта по умолчанию и переносит пути `nameState` на новое имя. Короткая команда сохраняет своё имя. Если переписать пути плагинов не удалось, папке возвращается прежнее имя, так как конфигурация не сохраняется.

`ccam import` регистрирует папку там, где она находится, как `dir` в `config.toml`. Он никогда не перемещает и не копирует её: скопированный токен обновления заставил бы копии выводить друг друга из аккаунта. Без аргумента с папкой он ищет незарегистрированные папки `~/.claude` и `~/.claude-*` и спрашивает про каждую; имена берутся из папки (`~/.claude` становится `main`, `~/.claude-work` становится `work`, `~/.claude-2` становится `acct2`), а папка, для которой имя вывести нельзя, выводится вместе с командой `--name`, которую нужно использовать. Без `CLAUDE_CONFIG_DIR` Claude Code хранит `.claude.json` в домашней папке рядом с `~/.claude`, поэтому, когда вы импортируете `~/.claude` и в ней нет собственного `.claude.json`, ccam создаёт ссылку `~/.claude/.claude.json` на `~/.claude.json`. После этого обычный `claude` и `ccam main` видят один и тот же аккаунт. Импортированные аккаунты ничего не делят, пока вы не выполните `ccam share`, и короткая команда для них не создаётся.

`ccam completion bash|zsh` выводит скрипт автодополнения, который также дополняет имена аккаунтов.

## Безопасность и особые случаи

- `ccam rename` и `ccam remove` отказываются работать, пока аккаунт запущен (`refuseRunning`), с сообщением `work is running; quit its Claude Code session first`, потому что сессия продолжает записывать в папку.
- `ccam remove` для импортированного аккаунта только отменяет его регистрацию: `✓ main unregistered (folder ... left in place)`. Ваша `~/.claude` не перемещается в корзину.
- `ccam add` отклоняет имя, папка для которого уже существует в `profiles/`, и предлагает `ccam import`.
- `ccam import` отклоняет уже зарегистрированную папку, путь, который не является папкой, и занятое имя.
- `ccam logout` и `ccam login` не проверяют, запущен ли аккаунт.
- Если передать вход в Claude Code не удаётся, например после завершения входа, выводится `ccam: could not pass the sign-in to Claude Code:` с причиной. Ошибка HTTP обычно означает, что срок действия кода истёк; начните вход заново.

## См. также

- [Запуск Claude Code](/features/launch)
- [Короткие команды](/features/shortcuts)
- [Общие MCP-серверы, настройки и история](/docs/sharing)
- [Несколько аккаунтов Claude Code](/guides/multiple-claude-code-accounts)
