---
title: Керування обліковими записами
description: Додавайте, входьте, переглядайте, перейменовуйте, видаляйте й імпортуйте облікові записи Claude Code командами ccam add, ccam login, ccam list, ccam rename, ccam remove і ccam import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 97cef50f48f4
---

Кожен обліковий запис ccam є окремою текою конфігурації Claude Code із власним входом. Ці команди створюють облікові записи, виконують вхід, виводять їх список, перейменовують і видаляють їх, а також реєструють теки, які у вас уже є.

## Використання

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

Там, де команда приймає `<name|number>`, можна вказати назву облікового запису або номер його слота, тобто стовпець `#` у `ccam list`.

`ccam add` у терміналі без параметрів відкриває майстер. Він запитує назву, потім показує контрольний список спільних елементів (позначено все, крім історії розмов), далі запитує `Create a shortcut command for work? [y/N]` і, якщо відповісти `y`, `Shortcut name [ccwork]:`. Будь-який із параметрів `--share`, `--shortcut` чи `--no-login` або конвеєр замість термінала пропускає майстер; тоді назва обов'язкова, а коротка команда не створюється, якщо `--shortcut` її не вказує.

`ccam list` показує слот, назву, коротку команду, електронну пошту, спільні елементи й останнє використання, а перед обліковим записом за замовчуванням стоїть `*`:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

З `--json` кожен обліковий запис є об'єктом із полями `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (відсутнє, якщо обліковий запис ще ніколи не запускали) і `share`. `share` зіставляє кожен елемент (`mcp`, `settings`, `memory`, `extensions`, `history`) з назвою його групи або з `"local"`.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `ccam add --no-login` | Створити обліковий запис без входу; увійти пізніше через `ccam login <name>` | вимкнено |
| `ccam add --share <item=group>` | Спільно використовувати елемент із групою або `item=local`, щоб залишити його приватним. Можна повторювати | `mcp`, `settings`, `memory`, `extensions` у `default`; `history` локально |
| `ccam add --shortcut <command>` | Також створити цю коротку команду в `~/.local/bin` | немає |
| `ccam list --json` | Вивести облікові записи у форматі JSON | вимкнено |
| `ccam remove -y`, `--yes` | Видалити без запитання | вимкнено |
| `ccam import --name <name>` | Назва облікового запису для вказаної теки | виводиться з назви теки |
| `ccam import -y`, `--yes` | Імпортувати кожну знайдену теку без запитання | вимкнено |
| `--dry-run` | Враховується командами `ccam list` і `ccam completion`, які нічого не змінюють. Відхиляється командами `add`, `login`, `logout`, `use`, `rename`, `remove` і `import` з повідомленням `` --dry-run is not supported by `ccam add`; nothing was changed `` (з назвою команди) | вимкнено |

`ccam login`, `ccam logout`, `ccam use` і `ccam rename` не мають власних параметрів. `ccam completion` приймає `bash` або `zsh`; будь-що інше завершується помилкою `only bash and zsh are supported`.

Назви облікових записів мають від 1 до 32 символів із `a-z`, `0-9`, `-` і `_` та починаються з літери або цифри. Назва не може складатися лише з цифр, бо `ccam 2` означає слот 2. Назва також не може збігатися з командою ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), тож `ccam <name>` ніколи не буває неоднозначним. Цей список міститься в `config.Reserved` у `internal/config`.

## Як це працює

`ccam add work --shortcut ccwork` виконує такі кроки в `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` перевіряє назву, а `shim.Check` перевіряє коротку команду, перш ніж щось буде створено.
2. `profile.Add` створює `~/.local/share/ccam/profiles/work/` з режимом `0700` і призначає обліковому запису найменший вільний номер слота (`Config.NextSlot`), тож номер видаленого облікового запису використовується повторно. Перший обліковий запис стає обліковим записом за замовчуванням.
3. `applyShare` пов'язує нову порожню теку з її спільними групами через `share.Build`, а потім `Apply`, наприклад `settings.json` з `groups/default/settings.json`. Якщо це не вдається, нова тека видаляється, бо залишена тека блокувала б кожну повторну спробу помилкою `folder already exists`.
4. Обліковий запис зберігається в `~/.config/ccam/config.toml` через `fsutil.AtomicWrite`: тимчасовий файл, `fsync`, потім `rename`, тож збій ніколи не залишає половину конфігурації.
5. `runLogin` запускає `claude auth login` зі змінною `CLAUDE_CONFIG_DIR`, що вказує на теку облікового запису. Якщо вхід не вдається, обліковий запис залишається, а ccam виводить помилку й радить пізніше виконати `ccam login work`.
6. `profile.ReadEmail` бере `oauthAccount.emailAddress` з `profiles/work/.claude.json`. Цей файл може мати кілька мегабайтів, тому `TopLevelValue` шукає один ключ замість розбору всього файлу.
7. `shim.Install` пов'язує `~/.local/bin/ccwork` з виконуваним файлом ccam, який під час запуску шукає `argv[0]` у `config.toml`.

Вхід може тривати кілька хвилин, і за цей час інший ccam може зберегти конфігурацію. Тому команди, які чекають, зберігають через `env.update`: вона повторно завантажує `config.toml`, застосовує одну зміну (тут електронну пошту) і зберігає. З тієї самої причини `ccam remove` повторно завантажує конфігурацію після свого запитання.

`ccam login` виконує кроки 5 і 6 для наявного облікового запису. `ccam logout` запускає `claude auth logout` в обліковому записі й очищає збережену електронну пошту. `ccam use` лише переписує `default` у `config.toml`.

### Вхід на віддаленому комп'ютері

Через SSH браузерна частина `claude auth login` може не спрацювати. Claude Code чекає на вхід на сервері зворотного виклику за адресою `http://localhost:<port>/callback` на сервері, але браузер, який завершує вхід, працює на вашому комп'ютері, де ця адреса нікуди не веде, і браузер показує помилку на кшталт `ERR_CONNECTION_REFUSED`. Адреса все одно містить одноразовий код, тож ccam може його доставити:

<!-- motion:accounts-relay -->

`runLogin` не передає ваш термінал у Claude Code. Він під'єднує стандартне введення Claude Code до каналу, а `relayLogin` копіює туди кожен введений вами рядок, тож код, вставлений у запит Claude Code `Paste code here`, усе одно надходить. Коли рядок є адресою зворотного виклику, `callbackURL` приймає його, лише якщо хост — `localhost`, `127.0.0.1` або `[::1]`, вказано порт, шлях — `/callback` і присутні обидва параметри `code` і `state`; `http://` можна пропустити. Потім `deliverCallback` запитує цю адресу на сервері з обмеженням у 10 секунд і без переходу за переспрямуваннями, як це зробив би браузер. Запит надходить до сервера зворотного виклику `claude auth login`, що очікує (пунктирна лінія на діаграмі). Claude Code отримує код і зберігає вхід, на Linux — у `.credentials.json` у теці облікового запису. Вставлений рядок ніколи не може змусити ccam звернутися до іншого хоста.

Перед початком входу ccam виводить:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

У разі успіху він виводить `→ Passed the sign-in to Claude Code on this machine.`

### Перейменування, видалення та імпорт

`ccam rename` і `ccam remove` діють на теку облікового запису та на те, що ccam зберігає під назвою облікового запису в інших місцях. `nameState` перелічує ці шляхи: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` і кеші використання `~/.local/state/ccam/usage/<name>.json` та `<name>.statusline.json`. `ccam remove work` обробляє їх у такому порядку:

<!-- motion:accounts-remove -->

`ccam remove` спершу відмовляється, якщо обліковий запис запущено, а потім запитує підтвердження. Якщо обліковий запис спільно використовує розширення, `RewritePluginPaths` переписує `installed_plugins.json` і `known_marketplaces.json` групи, бо Claude Code зберігає абсолютні шляхи до плагінів, які можуть проходити через теку цього облікового запису. `profile.Trash` переміщує теку до `~/.local/share/ccam/trash/<name>-<time>/`, обліковий запис зникає з `config.toml`, шляхи `nameState` видаляються, щоб пізніший обліковий запис із тією самою назвою їх не успадкував, а посилання короткої команди видаляється. Щоб відновити обліковий запис, перемістіть теку назад із кошика.

`ccam rename` перейменовує теку, переписує шляхи до плагінів у її `plugins/` зі старої теки на нову, переносить позначку облікового запису за замовчуванням і переміщує шляхи `nameState` під нову назву. Коротка команда зберігає свою назву. Якщо переписати шляхи до плагінів не вдається, теку перейменовують назад, оскільки конфігурацію не збережено.

`ccam import` реєструє теку там, де вона є, як `dir` у `config.toml`. Він ніколи не переміщує й не копіює її: скопійований токен оновлення змушував би копії виводити одна одну із системи. Без аргументу-теки він шукає незареєстровані теки `~/.claude` і `~/.claude-*` та запитує про кожну; назви походять від теки (`~/.claude` — це `main`, `~/.claude-work` — `work`, `~/.claude-2` — `acct2`), а теку, назву якої вивести не вдається, показано разом із командою з `--name`, яку слід використати. Без `CLAUDE_CONFIG_DIR` Claude Code зберігає `.claude.json` у вашій домашній теці поруч із `~/.claude`, тож коли ви імпортуєте `~/.claude`, а власного `.claude.json` у ній немає, ccam створює посилання `~/.claude/.claude.json` на `~/.claude.json`. Тоді звичайний `claude` і `ccam main` бачать той самий обліковий запис. Імпортовані облікові записи нічого не використовують спільно, доки ви не виконаєте `ccam share`, і коротка команда не створюється.

`ccam completion bash|zsh` виводить скрипт автодоповнення, який також доповнює назви облікових записів.

## Безпека та граничні випадки

- `ccam rename` і `ccam remove` відмовляються працювати, поки обліковий запис запущено (`refuseRunning`), з повідомленням `work is running; quit its Claude Code session first`, бо сеанс і далі записує в теку.
- `ccam remove` для імпортованого облікового запису лише скасовує його реєстрацію: `✓ main unregistered (folder ... left in place)`. Ваша `~/.claude` не переміщується в кошик.
- `ccam add` відхиляє назву, тека якої вже існує в `profiles/`, і радить `ccam import`.
- `ccam import` відхиляє вже зареєстровану теку, шлях, який не є текою, і зайняту назву.
- `ccam logout` і `ccam login` не перевіряють, чи запущено обліковий запис.
- Якщо ретранслятор не може зв'язатися з Claude Code, наприклад після завершення входу, він виводить `ccam: could not pass the sign-in to Claude Code:` із причиною. Помилка HTTP зазвичай означає, що код прострочено; почніть вхід заново.

## Пов'язані сторінки

- [Запуск Claude Code](/features/launch)
- [Короткі команди](/features/shortcuts)
- [Спільні сервери MCP, налаштування та історія](/docs/sharing)
- [Кілька облікових записів Claude Code](/guides/multiple-claude-code-accounts)
