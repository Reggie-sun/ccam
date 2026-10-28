---
title: Встановлення ccam
metaTitle: Встановлення ccam на Linux, macOS і Windows
description: Встановіть ccam на Linux, macOS або Windows однією командою, з пакета .deb або через go install та оновлюйте командою ccam update. Спершу потрібен Claude Code.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam запускає Claude Code за вас, тому спочатку встановіть [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) і переконайтеся, що `claude --version` працює.

Виконайте це в терміналі на Linux або macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт завантажує випуск для вашої системи, перевіряє його контрольну суму SHA-256 і розміщує `ccam` у `~/.local/bin`; там само ccam згодом створює вибрані вами короткі команди, наприклад `ccwork`. Якщо цієї теки немає у вашому `PATH`, скрипт виводить рядок, який потрібно додати до `~/.bashrc` або `~/.zshrc`.

У Windows виконайте це в PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Він завантажує `ccam_windows_<arch>.zip`, перевіряє його контрольну суму SHA-256 і кладе `ccam.exe` у `%USERPROFILE%\.local\bin`, теку, куди інсталятор Claude Code кладе `claude.exe`. Якщо цієї теки ще немає у вашому користувацькому `PATH`, скрипт додає її, залишаючи записи на кшталт `%USERPROFILE%\bin` без змін. Після цього відкрийте новий термінал.

## Підтримувані системи

- Linux на x86-64 і ARM64: Ubuntu, Debian, Fedora, Arch, Alpine та інші дистрибутиви. Двійковий файл статичний і не має залежностей.
- macOS на Apple silicon та Intel.
- Windows 10 і 11 на x86-64 та ARM64, у PowerShell, Windows Terminal або Git Bash:
  - Спільні теки працюють без режиму розробника.
  - Спільним файлам (`settings.json`, `CLAUDE.md`, `history.jsonl`) потрібен режим розробника (Параметри > Система > Для розробників), бо Windows дозволяє створювати посилання на файли лише в цьому режимі або адміністраторові; див. [спільне використання у Windows](/features/share#on-windows).
  - Без режиму розробника перший `ccam add` створює групу `default` лише із серверами MCP і розширеннями та пояснює, як додати решту, коли режим розробника буде ввімкнено.
- Усередині WSL ccam працює так само, як у будь-якій іншій системі Linux.

## Інші способи встановлення

**Пакет для Debian і Ubuntu.** Завантажте `ccam_<version>_linux_amd64.deb` зі [сторінки випусків](https://github.com/jhyunwoo/ccam/releases) і встановіть його:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Пакет встановлює `/usr/bin/ccam`.

**З вихідного коду.** З Go 1.26 або новішою версією:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

Якщо встановити тег випуску (`@latest` або `@vX.Y.Z`), ccam повідомляє цю версію, і працюють `ccam update` та пропозиція оновлення на панелі керування. Збірка з `@main` або з локальної копії з комітами після тегу повідомляє `dev`, і `ccam update` відмовляється її оновлювати, якщо не передати `--version`.

**Певна версія.** Задайте `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Автодоповнення в оболонці

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Автодоповнення працює для всіх команд і параметрів `ccam`, включно з назвами облікових записів.

## Плагін Claude Code

Плагін ccam навчає сам Claude Code користуватися ccam: перевіряти, у якого облікового запису є запас, запускати паралельні робочі процеси `claude -p` у різних облікових записах за допомогою `ccam pick` і пояснювати відсутність навички чи сервера MCP за допомогою `ccam ext`. Установіть його всередині Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Він додає:

- навичку, яка пропонує команди, що виконують вхід або переміщують файли, але ніколи їх не виконує;
- команди `/ccam:usage`, `/ccam:ext` і `/ccam:switch`;
- два хуки ([подробиці](/features/next#inside-claude-code)). Один показує сповіщення в терміналі, коли сеанс зупиняється на ліміті використання; на Linux і macOS з увімкненим [`auto_switch`](/features/next#auto-switch) ccam після цього переносить сеанс до іншого облікового запису, у якому ще є запас. Інший дає змогу `/ccam:switch [account]` передати розмову іншому обліковому запису.

Хуки запускають `ccam`, тому він має бути в `PATH`, який бачить Claude Code.

## Перевірка встановлення

```sh
ccam version
ccam doctor
```

`ccam doctor` перевіряє, що Claude Code знайдено і що кожна тека облікового запису, посилання та вхід у порядку. Далі перейдіть до [швидкого старту](/docs/quick-start).

## Оновлення

```sh
ccam update --check
ccam update
```

`ccam update --check` лише повідомляє, чи є новіший випуск. `ccam update` завантажує випуск для вашої системи, звіряє його контрольну суму SHA-256 з файлом `checksums.txt` цього випуску, а потім замінює файл `ccam` за один крок. Якщо контрольна сума не збігається, нічого не змінюється. Короткі команди, наприклад `ccwork`, є посиланнями на ccam, тому вони одразу запускають нову версію.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

Панель керування (`ccam` без аргументів, у терміналі) теж перевіряє, чи є новий випуск, щоразу, коли відкривається, і, якщо він є, пропонує його: натисніть `U` і підтвердьте. Жодна інша команда не перевіряє оновлення, і ccam ніколи не оновлюється, не запитавши вас. Щоб встановити певний випуск, зокрема старіший, передайте `--version`. `ccam --dry-run update` показує, що буде завантажено й замінено, нічого не змінюючи.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Повторний запуск скрипта встановлення також оновлює копію в `~/.local/bin`. Якщо ви зібрали ccam за допомогою `go install` з тегу випуску, можна виконати `ccam update` або ту саму команду `go install` ще раз.

**Windows.** `ccam update` завантажує `.zip` для Windows. У Windows не можна замінити запущену програму, тому старий виконуваний файл перейменовується на `ccam.exe.old`, і наступне оновлення його видаляє; поки сеанс ще працює з попереднього `ccam.exe.old`, назва буде `ccam.exe.old-<time>`. Короткі команди, які є жорсткими посиланнями, оновлюються одразу після цього. Повторний запуск `install.ps1` працює так само, навіть поки сеанс Claude Code ще працює зі старого виконуваного файлу, і виводить `OK 2 shortcuts now run the new ccam`, коли оновлює короткі команди, що є жорсткими посиланнями.

**Пакет для Debian і Ubuntu.** `ccam update` не замінює `/usr/bin/ccam`, встановлений через apt, а натомість виводить ці кроки. Завантажте новий `ccam_<version>_linux_amd64.deb` зі [сторінки випусків](https://github.com/jhyunwoo/ccam/releases) і встановіть його поверх старого:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Видалення

`ccam uninstall` видаляє ccam і все, що він створив. Спершу він показує список того, що буде видалено, і питає:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Кожен обліковий запис видаляється так само, як це робить `ccam remove`. З облікового запису в теках ccam спершу виконується вихід, а потім його видаляють разом із текою, історією розмов і налаштуваннями.
- Тека `~/.claude`, яку ви [імпортували](/docs/commands), залишається там, де вона є, для звичайного `claude`. Якщо вона входить до групи, то отримує назад копію всього, що мала спільним. Рядок стану ccam у її `settings.json` поступається місцем рядку стану, який обгортав ccam, якщо такий був. Якщо в ній є плагін ccam, ccam виводить команду `/plugin uninstall`, яку треба виконати в Claude Code.
- Короткі команди, `~/.config/ccam`, `~/.local/share/ccam` (разом із кошиком ccam) і `~/.local/state/ccam` видаляються, а потім і сама команда `ccam`.

ccam відмовляється працювати, поки якийсь обліковий запис запущено; спершу завершіть його сеанси Claude Code. `/usr/bin/ccam`, встановлений через apt, залишається; видаліть його потім командою `sudo apt remove ccam`. У Windows запущена програма не може видалити сама себе, тож `ccam uninstall` наприкінці виводить команду `Remove-Item` для `ccam.exe`.
