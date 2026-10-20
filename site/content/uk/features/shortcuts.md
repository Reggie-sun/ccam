---
title: Короткі команди
description: Дайте обліковому запису власну команду, наприклад ccwork, через ccam shortcut і повторно створіть посилання в ~/.local/bin через ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Коротка команда — це команда з назвою, яку ви вибираєте самі, і вона запускає Claude Code від імені одного облікового запису: `ccwork --resume` робить те саме, що й `ccam work --resume`. Це посилання на ccam у `~/.local/bin`, тож вона працює в bash, zsh, скриптах і будь-чому іншому, що запускає програми з `PATH`. Короткі команди необов'язкові; створіть їх для облікових записів, які ви запускаєте часто.

## Використання

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam підтверджує кожну зміну:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Обліковий запис має не більше однієї короткої команди, яка зберігається як `shortcut` у `~/.config/ccam/config.toml`. Майстер `ccam add` запитує `Create a shortcut command for <name>? [y/N]`, а потім `Shortcut name [cc<name>]:`. Поза майстром `ccam add` створює її лише з `--shortcut`. `ccam list` показує коротку команду кожного облікового запису або `-`.

Коротка команда передає всі аргументи в Claude Code так само, як `ccam <name>` після назви. Це стосується й `--dry-run` та `--help`: щоб попередньо переглянути запуск, використовуйте `ccam --dry-run work`.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `--remove` | Видалити коротку команду облікового запису та її посилання | вимкнено |
| `-h`, `--help` | Довідка для команди | |
| `--dry-run` | `ccam shortcut` враховує його: перевіряє назву й виводить `Would set the shortcut of work to ccwork.` або `Would remove ccwork, the shortcut of work.`, нічого не зберігаючи й не створюючи посилань | вимкнено |

`ccam shims install` не приймає аргументів і відхиляє `--dry-run` з повідомленням `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Правила для назв, які перевіряються в такому порядку:

- від 1 до 64 символів із `a-z`, `0-9`, `.`, `-` і `_`, що починаються з літери або цифри;
- не `ccam`, `claude` і не команда ccam, як-от `run` чи `usage`;
- не коротка команда іншого облікового запису;
- не файл у `~/.local/bin`, який ще не є посиланням на ccam;
- не назва іншого виконуваного файлу у вашому `PATH`.

Якщо вказати команду разом із `--remove`, виконання зупиниться з повідомленням `give a command or --remove, not both`. Якщо задати назву, яку обліковий запис уже має, буде виведено `work: ccwork (unchanged)`, а посилання буде створено заново, якщо його немає.

## Як це працює

Задання назви запускає `checkShortcut`, яка викликає `shim.Check` з текою bin і справжнім шляхом до виконуваного файлу ccam:

<!-- motion:shortcuts-check -->

Назва проходить перевірки по черзі. Біля кожної перевірки показано назву, яку вона відхиляє: `Work` не відповідає шаблону, `run` є командою ccam, `ccside` належить іншому обліковому запису, `ccfile` є наявним файлом у `~/.local/bin`, а `git` приховав би `/usr/bin/git`. Лише назва, що проходить усі перевірки, зберігається й отримує посилання.

1. `shim.Check` зіставляє назву з `^[a-z0-9][a-z0-9._-]{0,63}$`, потім з `config.Reserved`, потім з `cfg.FindShortcut`. Він викликає `os.Lstat` для `~/.local/bin/<name>`, а для кожної іншої теки в `PATH` викликає `os.Stat` і відхиляє виконуваний файл, який не є посиланням на ccam.
2. ccam зберігає `config.toml` через `fsutil.AtomicWrite` (тимчасовий файл, `fsync`, перейменування).
3. `shim.Remove` видаляє посилання старої короткої команди, але лише якщо воно є посиланням на ccam.
4. `shim.Install` створює посилання для кожної налаштованої короткої команди, якого бракує, через `os.Symlink` на виконуваний файл ccam, шлях до якого визначається через `os.Executable` і `filepath.EvalSymlinks`. Він ніколи нічого не видаляє.

Посилання вважається власним посиланням ccam, коли `os.Readlink` повертає шлях до ccam або коли посилання й шлях до ccam ведуть до того самого файлу. Кожне видалення й кожна перевірка перезапису проходять цей тест, тож ccam не чіпає файлів, яких він не створював.

Коли ви запускаєте коротку команду, ccam потрібна лише назва, яку ви ввели:

<!-- motion:shortcuts-link -->

Оболонка знаходить `ccwork` у `PATH`, ядро переходить за посиланням і запускає виконуваний файл ccam, а `argv[0]` усе ще дорівнює `ccwork`. В `app.Main` будь-яка назва, відмінна від `ccam`, один раз завантажує `config.toml` і викликає `shim.Lookup`, який знаходить обліковий запис, чий `shortcut` дорівнює `filepath.Base(argv[0])`. Потім `launcher.LaunchWith` проходить ті самі кроки, що й [`ccam <name>`](/features/launch), з уже завантаженою конфігурацією й передає `argv[1:]` у `claude` без змін.

Чому символьне посилання, а не псевдонім оболонки чи скрипт-обгортка: символьне посилання працює в кожній оболонці та в скриптах, не додає запуску ще однієї оболонки (приблизно 1–2 мс), а запуск і далі читає один файл. `scripts/perf-check.sh` вимірює коротку команду за допомогою `hyperfine` з тим самим обмеженням у 10 мс (p95), що й для `ccam <name>`.

### У Windows

Коротка команда — це `ccwork.exe` у `%USERPROFILE%\.local\bin`, а `shim.Name` відкидає `.exe` (у будь-якому регістрі) перед пошуком, тож `ccwork`, `ccwork.exe` і `CCWORK.EXE` запускають `work`. Спершу ccam створює символьне посилання, що Windows дозволяє в режимі розробника або адміністраторові. Без цього права він створює жорстке посилання: друге ім'я того самого файлу `ccam.exe`, для якого привілеї не потрібні. Жорстке посилання вважається власним посиланням ccam, коли це той самий файл, що й `ccam.exe`. Оскільки `ccam update` ставить новий файл на місце `ccam.exe`, після цього він оновлює кожну коротку команду з жорстким посиланням через `shim.Relink` і виводить `✓ shortcut ccwork renewed`. Перевірка назви шукає `<name>` з кожним розширенням із `PATHEXT` (`.exe`, `.cmd` тощо) в інших теках `PATH`.

## Безпека та граничні випадки

- **Посилання, яке не належить жодному обліковому запису.** Якщо посилання на ccam у `~/.local/bin` має назву, яку не використовує жоден обліковий запис, наприклад після того, як ви вручну змінили `shortcut` у `config.toml`, ccam не переходить до облікового запису за замовчуванням. Він зупиняється з повідомленням `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Задайте назву знову через `ccam shortcut` або видаліть посилання.
- **Наявні файли.** ccam ніколи не замінює файл, якого він не створював. `ccam shims install` зупиняється з повідомленням `<path> already exists and is not a ccam link, so it was left alone`.
- **Переміщений виконуваний файл ccam.** Старі посилання тоді вказують на шлях, якого більше немає, і ccam не вважає їх своїми. Видаліть їх, а потім виконайте `ccam shims install`.
- **Перейменування та видалення.** `ccam rename` зберігає назву короткої команди, і посилання далі працює, бо воно вказує на ccam, а ccam шукає назву в `config.toml`. `ccam remove` видаляє посилання облікового запису. `ccam shortcut <account> <new>` видаляє старе посилання, перш ніж створити нове.
- **PATH.** Коли ccam створює посилання, а `~/.local/bin` немає у вашому `PATH`, він виводить `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`, а за ним рядок `export PATH=...`, який треба додати. У Windows натомість він виводить рядок PowerShell, який додає теку до вашого користувацького `PATH`.
- **Пізніша перевірка.** `ccam doctor` повідомляє `missing shortcuts in /home/me/.local/bin: ccwork` і попереджає, коли теки немає в `PATH`. `ccam doctor --fix` створює відсутні посилання.

## Пов'язані сторінки

- [Запуск Claude Code](/features/launch)
- [Керування обліковими записами](/features/accounts)
- [Довідник команд](/docs/commands#shortcuts)
- [Встановлення ccam](/docs/install)
