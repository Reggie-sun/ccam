---
title: Спільні сервери MCP
description: "Один список серверів MCP на групу: ccam mcp add, add-json, list, remove та import; він передається через --mcp-config або синхронізується в кожен запис."
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code зберігає користувацькі сервери MCP у `.claude.json`, тому самому файлі, що містить вхід облікового запису, тож цей файл не можна спільно використовувати між обліковими записами. ccam зберігає спільні сервери у файлі групи, `mcp.json`, і передає його кожному обліковому запису групи під час запуску. Використовуйте `ccam mcp`, щоб один раз додати сервер для всіх своїх облікових записів, і `ccam mcp import`, щоб перенести до інших сервер, який ви додали в межах одного сеансу.

## Використання

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
ccam mcp import work gh --force                                   # replace the group's different gh
```

Зміна застосовується з наступного запуску кожного облікового запису; запущені сеанси зберігають сервери, з якими вони стартували:

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

Під групами `ccam mcp list` називає кожен сервер, який не доходить до інших облікових записів: `only in account work: pg` для сервера, який обліковий запис додав для себе, `only in account side, for /home/me/code/app only: notes` для сервера, доданого для однієї теки (звичайний `claude mcp add`, див. нижче), `account work has its own …, used in its sessions in place of the group's`, коли власний сервер облікового запису має назву сервера групи, `.mcp.json` цієї теки з обліковими записами, схвалення яких він чекає (`waiting for approval in`), і сервери, які є `in no account, only in ~/.claude.json (plain claude's)`, разом із командою [`ccam migrate --mcp`](/features/migrate), яка їх копіює.

Додавання назви, яка вже є в групі, замінює її визначення й виводить `updated` замість `added`. Сервери потрапляють до облікового запису, лише якщо він входить до цієї групи, а група робить MCP спільним (`ccam group set <group> --add mcp`); коли група цього не робить, `add` про це повідомляє.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `-g`, `--group <group>` | цільова група для `add`, `add-json` і `remove` | єдина група, через яку якийсь обліковий запис спільно використовує MCP, інакше `default` |
| `-- <command> [args…]` (`add`) | запускати цю команду як сервер stdio | |
| `--url <url>` (`add`) | URL віддаленого сервера (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | транспорт сервера з `--url` | `http` |
| `--header <'Name: value'>` (`add`) | заголовок HTTP, можна повторювати | немає |
| `-e`, `--env <KEY=VALUE>` (`add`) | змінна середовища для сервера stdio, можна повторювати | немає |
| `--force` (`import`) | замінити відмінний сервер групи з тією самою назвою; для старого `mcp.json` створюється резервна копія | вимкнено |
| `--keep-local` (`import`) | залишити записи в `.claude.json` облікового запису | вимкнено |
| `-h`, `--help` | довідка для кожної команди | |
| `--dry-run` (глобальний) | враховується лише `mcp list`; `add`, `add-json`, `remove` і `import` відхиляють його й нічого не змінюють | вимкнено |

`add` потребує або команди після `--`, або `--url`; інакше він завершується помилкою `a command to run (after --) or --url is required`. `--header` стосується серверів із `--url`, а `--env` — серверів із командою. `add-json` приймає той самий JSON, що й `claude mcp add-json`, і визначення має бути об'єктом із `command` або `url`. `import` завжди використовує групу, через яку обліковий запис спільно використовує MCP, а `list` показує всі групи, тож `-g` на них не впливає.

## Як це працює

Задіяні файли:

- `~/.local/share/ccam/groups/<group>/mcp.json`: спільні сервери, `{"mcpServers": {…}}`, формат, який читає `--mcp-config`. Він записується через `fsutil.AtomicWrite` з режимом `0600`, бо записи серверів часто містять ключі API, і лише під блокуванням `<group>/.mcp.json.lock`, тож два одночасні `ccam mcp add` обидва зберігають свій сервер.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: власний файл Claude Code із серверами облікового запису. ccam читає лише його значення `mcpServers` за допомогою сканера `profile.TopLevelValue` і ніколи не розбирає весь файл на шляху запуску.
- `~/.local/state/ccam/mcp/<name>/`: файл `seen`, `sync-stamp` (з хешем кожної копії, яку записав ccam, щоб він міг відрізнити копію, якої ніхто не змінював, від відредагованої вами), `sync-checked` (те, що стратегія синхронізації бачила востаннє), `kept` (сервери, які обліковий запис зберіг, коли вийшов зі спільного доступу до MCP) і резервні копії `.claude.json`.

Є дві стратегії, які задає `mcp_strategy` у `~/.config/ccam/config.toml`. За замовчуванням діє стратегія прапорця:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` знаходить групу MCP облікового запису. `mcp.FlagArgs` одним викликом `stat` перевіряє, що `mcp.json` групи існує, не розбираючи його.
2. `Plan.InsertArgs` ставить `--mcp-config=<path>` одразу після `claude`, перед вашими аргументами, тож ваші власні `--mcp-config` чи `--strict-mcp-config` усе одно застосовуються поверх нього.
3. `mcp.NewLocal` перевіряє, чи змінився `.claude.json` облікового запису після останнього запуску, і лише тоді розбирає його сервери.
4. `syscall.Exec` замінює ccam на Claude Code.

Шлях передається одним аргументом, `--mcp-config=<path>`, а не двома. `--mcp-config` у Claude Code приймає кілька значень, тож `claude --mcp-config a.json "prompt"` прочитав би запит як другий файл конфігурації й завершився б помилкою `MCP config file not found`. Форма з одним аргументом на цьому зупиняється, і запит чи підкоманда, що йде далі, працює. `ccam work mcp list`, `get` і `remove` не отримують прапорця, бо Claude Code там його відхиляє; натомість ccam називає сервери групи в стандартному потоці помилок ([докладніше](/features/launch#shared-mcp-servers)). Claude Code завантажує сервери з `--mcp-config` додатково до власних серверів облікового запису; у разі збігу назв перемагає запис у `.claude.json`. Стратегія прапорця записує `.claude.json` лише для Remote Control, коли обліковий запис виходить зі спільного доступу до MCP або знову приєднується до нього ([нижче](#leaving-and-joining)), і для `ccam mcp import`.

`claude remote-control` (`rc`) не приймає `--mcp-config`, а його сеанси читають сервери лише з `.claude.json`. Для `ccam work rc` ccam не додає прапорець і, якщо є термінал, запитує, чи скопіювати сервери групи в `.claude.json` із резервною копією, як це робить стратегія `sync`. Він запитує, доки копії немає, а потім знову лише після зміни `mcp.json` групи або коли скопійований сервер зникає ([Remote Control](/features/launch#remote-control)). Оскільки така копія перемагає `--mcp-config` у разі збігу назв, `ccam mcp add`, `add-json`, `remove` і `import` оновлюють копії, що є в інших облікових записах групи:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Рядок `✓ Updated` з'являється, лише коли копія змінилася. Копію, створену старішим ccam, оновлюють так само, як і копію після `ccam migrate`.

Сервер, доданий у сеансі через `claude mcp add`, потрапляє лише до `.claude.json` цього облікового запису. Звичайний `claude mcp add` використовує локальну область: сервер завантажується для цього одного облікового запису в цій одній теці (`-s user` робить його сервером облікового запису в усіх теках). ccam не може побачити його в цей момент, бо Claude Code його замінив, тому перевіряє під час наступного запуску. Файл `seen` містить час зміни `.claude.json` і назви серверів, про які вже повідомлено. Коли час відрізняється, ccam порівнює назви з групою й один раз виводить підказку для кожного нового сервера, не чекаючи введення, за обох стратегій:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Підказка також називає власний сервер облікового запису, який має назву сервера групи, бо в його сеансах перемагає сервер облікового запису.

З `mcp_strategy = "sync"` сервери з'являються в `claude mcp list` і `/mcp` як звичайні користувацькі сервери, а ccam перед запуском записує їх у `.claude.json` кожного облікового запису:

<!-- motion:mcp-sync -->

1. `mcp.Sync` порівнює час зміни `mcp.json` групи з `sync-stamp`, а `.claude.json` облікового запису — з тим, що бачив востаннє. Якщо ніщо не змінилося, більше нічого не відбувається. Інакше він також перевіряє, що кожен сервер групи досі є в `.claude.json`, і повертає копію, яка зникла (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Він бере власне блокування Claude Code, теку `.claude.json.lock` поруч із файлом, створену через `mkdir`. Він чекає до 2 секунд і вважає блокування, старіше за 10 секунд, застарілим, як це робить бібліотека блокувань Claude Code.
3. Він зберігає поточний `.claude.json` до `backups/` у теці стану й залишає 5 найновіших копій.
4. Він задає кожен сервер групи в `mcpServers`, видаляє сервери, які записав раніше і яких у групі вже немає, і атомарно записує файл з його початковим режимом, лише коли щось змінилося. Він змінює або видаляє лише ті записані ним копії, яких відтоді ніхто не змінював. Власний сервер облікового запису з назвою сервера групи лишається за обліковим записом, і ccam один раз про це повідомляє: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. Він записує нову позначку часу з хешем кожної копії, і Claude Code запускається без `--mcp-config`.

Синхронізація пропускається з повідомленням `ccam: account work is already running, skipping MCP sync`, коли сеанс облікового запису запущено, хоч би як було запущено Claude Code, бо Claude Code переписує `.claude.json`, поки працює. Невдала синхронізація виводить `ccam: MCP sync failed (launching anyway): …` і все одно виконує запуск.

### Вихід і приєднання

Обліковий запис, який припиняє спільно використовувати MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`), зберігає сервери групи: план копіює їх у його `.claude.json` із резервною копією, а його власні сервери з тими самими назвами лишаються. З `--fresh` він не отримує жодного, а план перелічує сервери, які перестануть завантажуватися. Коли він знову приєднується до групи, що має ці сервери, збережені копії, яких ніхто не змінював, видаляються, тож пізніші зміни групи до нього доходять.

`ccam mcp import` переміщує власні сервери облікового запису до його групи:

<!-- motion:mcp-import -->

ccam читає сервери облікового запису, його власні й ті, що він має лише для однієї теки, і порівнює кожен із сервером групи з тією самою назвою за допомогою `mcp.SameDef`, яка ігнорує порядок ключів і пробіли. Сервер теки стає сервером групи для всіх тек (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); одну назву з двома різними визначеннями для тек буде відхилено. Сервери плагінів не пропонуються, бо їх приносить плагін. Сервер, якого в групі немає або який має там те саме визначення, копіюється. Відмінний сервер пропускається, бо визначення групи використовують усі інші облікові записи:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

З `--force` сервер групи замінюється, а попередній `mcp.json` спершу зберігається до `.ccam-backup/<time>/mcp.json` усередині теки групи. Потім `mcp.Remove` видаляє скопійовані записи з `.claude.json` облікового запису, його власні й записи теки, одним записом, під тим самим блокуванням і після створення резервної копії в теці стану. `mcp.json` групи лишається заблокованим від його завантаження до зникнення цих записів, тож ніхто не видалить сервер із групи тим часом. Це потрібно, бо в разі збігу назв запис облікового запису переміг би запис групи. `--keep-local` їх залишає.

## Безпека та граничні випадки

- **Входи.** Сервери, що виконують вхід через OAuth, зберігають свої токени в `.credentials.json` кожного облікового запису (`mcpOAuth`), а не в списку серверів, тож спільне використання сервера не робить спільним його вхід. Увійдіть один раз у кожному обліковому записі через `/mcp`. ccam ніколи не копіює входи, навіть [`ccam migrate`](/features/migrate): скопійований токен оновлення виводить із системи всіх інших власників. [`ccam ext`](/features/ext) і `ccam doctor` називають сервери, у які обліковий запис не увійшов, тоді як інший обліковий запис увійшов, а `ccam next --continue` повідомляє, яких бракує новому обліковому запису. Група ніколи не отримує файлу входу.
- **Запущені облікові записи.** `import` копіює до групи, але залишає власні записи запущеного облікового запису й повідомляє про це; виконайте його знову після завершення сеансу. Синхронізація пропускає запущені облікові записи. Обидві враховують і Claude Code, запущений без ccam.
- **Копії для Remote Control.** Копію запущеного облікового запису не змінюють: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` З `mcp_strategy = "sync"` нічого оновлювати не треба, бо кожен запуск виконує синхронізацію.
- **Пошкоджені файли.** Якщо `.claude.json` не є дійсним JSON, синхронізація не чіпає його й повідомляє про помилку. Якщо взяти блокування за 2 секунди не вдається, ccam повідомляє `could not lock … (Claude Code is writing the file)`.
- **Облікові записи, які не використовують MCP спільно.** `import` відмовляє й радить спершу додати обліковий запис до групи, яка робить MCP спільним (`ccam group join <group> side`).
- **Назви груп.** `-g` приймає лише назву групи; вона стає текою в `groups/`, тож назви на кшталт `../x` відхиляються.

## Пов'язані сторінки

- [Спільне використання між обліковими записами](/features/share)
- [Remote Control з ccam](/guides/remote-control)
- [Спільні сервери MCP, налаштування та історія](/docs/sharing)
- [Копіювання наявного налаштування](/features/migrate)
- [Усі команди](/docs/commands)
