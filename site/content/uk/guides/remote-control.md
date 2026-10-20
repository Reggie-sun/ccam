---
title: Remote Control з ccam
description: "Remote Control з будь-якого облікового запису через ccam <name> rc: чинні параметри, чому ccam пропонує копіювати спільні сервери MCP і як усунути проблеми."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "5f0b90ea53d6"
---

`ccam work rc` запускає Remote Control з Claude Code від імені облікового запису `work`, тож ви можете керувати сеансами на цьому комп'ютері з claude.ai/code або мобільного застосунку Claude. Команда працює як `claude rc` з однією відмінністю: якщо обліковий запис спільно використовує сервери MCP, ccam один раз запитує, чи скопіювати їх у `.claude.json` облікового запису, бо сеанси Remote Control більше ніде не шукають сервери.

## Запуск

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Коротка команда](/features/shortcuts) працює так само: `ccwork rc` — це `ccam work rc`. Виконайте команду в теці, де хочете працювати, потім відкрийте claude.ai/code або застосунок Claude, увійшовши в той самий обліковий запис Claude, що й `work`.

## Параметри

Тут зустрічаються параметри двох видів: власні параметри ccam, які ставлять перед назвою облікового запису або задають у налаштуваннях ccam, і параметри Remote Control, які йдуть після `rc` і передаються в Claude Code без змін.

### Параметри й налаштування ccam

| Параметр або налаштування | Що робить | За замовчуванням |
| --- | --- | --- |
| `--dry-run` (перед обліковим записом) | Вивести середовище й команду, а також повідомлення про MCP, якщо воно є. Нічого не запитує й не копіює | вимкнено |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Запитується, коли обліковий запис спільно використовує сервери MCP, а його `.claude.json` не містить поточного списку групи. `y` копіює їх, `N` або Enter запускає без них | `N` |
| `mcp_strategy = "sync"` у `~/.config/ccam/config.toml` | Копіювати сервери групи в `.claude.json` перед кожним запуском, з Remote Control чи без, щоб запитання ніколи не з'являлося | `"flag"` |
| `ccam group leave <group> work` | Вивести обліковий запис із групи, яка робить для нього MCP спільним; тоді Remote Control запускається лише з власними серверами облікового запису й нічого не запитує | у `default` |

### Параметри Remote Control

Усе, що йде після `rc`, належить Claude Code. Нижче наведено параметри `claude remote-control` у Claude Code 2.1.304; `ccam work rc --help` виводить список для вашої версії.

| Параметр | Що робить |
| --- | --- |
| `--name <name>` | Назва сеансу, яку показує claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Префікс для автоматично створених назв сеансів (за замовчуванням назва хоста) |
| `-c`, `--continue` | Знову під'єднатися до сеансу, який Remote Control востаннє записав для цієї теки, у межах приблизно 4 годин |
| `--session-id <id>` | Знову під'єднатися до певного сеансу |
| `--permission-mode <mode>` | Режим дозволів для нових сеансів: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` або `plan` |
| `--spawn <mode>` | `same-dir` (за замовчуванням), `worktree` (окремий git worktree для кожного сеансу) або `session` (один сеанс, завершується разом із ним) |
| `--capacity <N>` | Найбільша кількість одночасних сеансів у режимах `same-dir` і `worktree` (за замовчуванням 32) |
| `--[no-]create-session-in-dir` | Створити перший сеанс у поточній теці під час запуску (за замовчуванням увімкнено) |
| `--[no-]chrome` | Claude in Chrome для нових сеансів (за замовчуванням вимкнено) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Налагоджувальний і докладний вивід |

Ставте `rc` одразу після назви облікового запису. ccam вважає командою Claude Code перший аргумент, що не починається з `-`, тож у `ccam work rc --name laptop` він бачить `rc`, а в `ccam work --model opus rc` бачить `opus`.

## Що робить ccam

Запуск проходить звичайний шлях ([Запуск Claude Code](/features/launch)) з однією зміною на кроці MCP:

1. ccam знаходить обліковий запис, перевіряє його спільні посилання й знаходить `claude`, як під час будь-якого запуску.
2. Він бачить, що команда — `rc` або `remote-control`, і не додає `--mcp-config`.
3. Якщо обліковий запис спільно використовує сервери MCP, ccam порівнює `mcp.json` групи з `sync-stamp` у `~/.local/state/ccam/mcp/work/`. Якщо `.claude.json` уже містить поточні сервери або група їх не має, він одразу переходить до кроку 6.
4. Інакше він називає сервери, яких бракуватиме сеансам:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Після відповіді `y` він зберігає резервну копію `.claude.json`, бере власне блокування Claude Code, додає сервери групи до `mcpServers` і зберігає інші сервери облікового запису. Після відповіді `N`, без термінала або поки обліковий запис запущено, він нічого не змінює й повідомляє про це:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Він реєструє сеанс, указує в `CLAUDE_CONFIG_DIR` теку облікового запису й замінює себе на `claude rc` через `exec`.

Після цього Remote Control продовжує працювати як сервер. Кожен сеанс, який він для вас запускає, є дочірнім процесом цього сервера, тож працює як `work`: із входом, налаштуваннями й історією розмов `work`, а витрати зараховуються до ліміту використання `work`. Поки сервер працює, ccam вважає обліковий запис запущеним.

Пробний запуск показує різницю. Звичайний запуск передає сервери групи аргументом, а Remote Control — ні:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Чому це працює саме так

**Remote Control не приймає `--mcp-config`.** Під час звичайних запусків ccam передає сервери групи в Claude Code як `--mcp-config=<path>` і ніколи не записує файли Claude Code. `claude remote-control` відхиляє цей параметр перед назвою своєї команди й не приймає жодних після неї, тож до ccam 0.6.0 `ccam work rc` завершувався помилкою для кожного облікового запису, який спільно використовував сервери MCP. Тепер ccam цей параметр не додає.

**Його сеанси читають сервери лише з `.claude.json`.** Сеанси, які запускає Remote Control, не успадковують параметри командного рядка. Єдине місце, де вони знаходять користувацькі сервери MCP, — `.claude.json` облікового запису. Цей файл також містить дані входу облікового запису й стан для кожної теки, тож ccam не може пов'язати його зі спільною копією, як робить це з `settings.json`. Скопіювати в нього сервери — єдиний спосіб дати сеансам спільний список.

**ccam запитує перед записом.** Зі стратегією за замовчуванням (прапорець) ccam обіцяє не чіпати `.claude.json`. Remote Control — єдиний випадок, коли йому доводиться туди писати, тож він спершу запитує, зберігає резервну копію й не чіпає нічого, крім `mcpServers`.

**Він не пише, поки обліковий запис запущено.** Claude Code переписує `.claude.json` під час сеансу, і зміна, зроблена ззовні, загубилася б або була б перезаписана. Тому ccam пропускає копіювання, коли запущено будь-який сеанс облікового запису, зокрема інший сервер Remote Control.

**Він запитує один раз на кожну зміну.** Після копіювання `sync-stamp` зберігає час зміни `mcp.json` групи. ccam запитує знову лише тоді, коли цей файл змінюється, наприклад після `ccam mcp add`. Якщо ви відповіли `N`, копії ще немає, тож запитання з'являється під час кожного запуску Remote Control.

**`ccam mcp` підтримує копії актуальними.** У разі збігу назв запис у `.claude.json` перемагає `--mcp-config`, тож стара копія приховала б новіше визначення і під час звичайних запусків. Тому `ccam mcp add`, `add-json`, `remove` і `import` оновлюють копії в усіх облікових записах, які не запущено:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Без термінала відповідь — «ні».** Remote Control часто запускають зі скрипта, служби або через `nohup`, де ніхто не може відповісти. Тоді ccam запускається без серверів, а не чекає.

## Проблеми та розв'язання

### У сеансах Remote Control немає спільних серверів MCP

На запитання відповіли `N`, термінала не було або обліковий запис було запущено. Зупиніть Remote Control, закрийте інші сеанси облікового запису й знову виконайте `ccam work rc` у терміналі; відповідайте `y`. Щоб запитання ніколи не з'являлося, задайте `mcp_strategy = "sync"` у `~/.config/ccam/config.toml`. Сервери вже запущеного сеансу не змінюються під час його роботи, тож після копіювання почніть новий сеанс.

### "Account work is running, so its .claude.json is left alone"

Сеанс `work` досі відкритий: звичайний `ccam work`, коротка команда або інший сервер Remote Control. Завершіть його й знову виконайте `ccam work rc`. Сеанси, запущені іншим обліковим записом, не враховуються.

### ccam запитує щоразу

Ви ще не скопіювали сервери. Один раз відповідайте `y` або задайте `mcp_strategy = "sync"`. Якщо сервери в Remote Control вам узагалі не потрібні, припиніть спільне використання MCP з цим обліковим записом: виведіть його з тієї групи (`ccam group leave default work`) або тримайте MCP в окремій групі, до якої цей обліковий запис не входить.

### `ccam work rc` одразу завершується помилкою щодо `--mcp-config`

Ваш ccam старший за 0.6.0 і передає `--mcp-config` у Remote Control. Виконайте `ccam update`. Та сама помилка в актуальному ccam означає, що перед `rc` стояв параметр зі значенням, як у `ccam work --model opus rc`, і ccam сприйняв `opus` як команду; ставте `rc` одразу після назви облікового запису.

### Зміна, зроблена через `ccam mcp`, не доходить до Remote Control

`ccam mcp add` і `remove` оновлюють копії лише в облікових записах, які не запущено, і повідомляють про це: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Зупиніть Remote Control і запустіть його знову; ccam побачить змінений `mcp.json` і запитає ще раз.

### Копія замінила мій власний сервер

Якщо обліковий запис мав власний сервер з тією самою назвою, що й сервер групи, копія його замінила. Попередній `.claude.json` лежить у `~/.local/state/ccam/mcp/work/backups/`, де зберігаються п'ять найновіших копій. Щоб ваша версія була в усіх облікових записах, перенесіть її до групи командою `ccam mcp import work <server> --force` до наступного запуску Remote Control.

### Я хочу видалити скопійовані сервери

`ccam mcp remove <server>` видаляє сервер із групи й з усіх копій. Щоб повернути стан до копіювання, відновіть найновішу резервну копію з `~/.local/state/ccam/mcp/work/backups/`, коли жоден сеанс облікового запису не запущено.

### Remote Control запитує, чи довіряти теці, або завершується без термінала

Довіра до тек зберігається для кожного облікового запису окремо в `.claude.json`, тож тека, якій ви довіряли як `personal`, для `work` нова. Один раз виконайте `ccam work` у цій теці й підтвердьте довіру, потім запускайте там Remote Control.

### "You must be logged in with a Claude account that has a subscription"

Для Remote Control потрібен вхід із підпискою Claude. Перевірте email облікового запису в `ccam list` і за потреби увійдіть знову командою `ccam login work`.

### Сеанс не з'являється на телефоні

У застосунку Claude або на claude.ai/code виконано вхід в інший обліковий запис Claude. Порівняйте його з email, який `ccam list` показує для `work`.

### ccam попереджає, що інший Claude Code спільно використовує історію розмов у цій теці

У тій самій теці запущено інший обліковий запис із тієї самої групи історії. Сеанси, що продовжують ту саму розмову з обох облікових записів, можуть змішати історію. Відповідайте `y`, лише якщо вони працюватимуть із різними розмовами, або запустіть Remote Control в іншій теці.

## Пов'язані сторінки

- [Запуск Claude Code](/features/launch)
- [Спільні сервери MCP](/features/mcp)
- [Спільне використання між обліковими записами](/features/share)
- [Керування обліковими записами](/features/accounts)
