---
title: Встановлення ccam
metaTitle: Встановлення ccam на Linux, macOS і Windows
description: Встановіть ccam на Linux, macOS або Windows однією командою, з пакета .deb або з коду та оновлюйте командою ccam update. Спершу потрібен Claude Code.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "07fab3165030"
---

Виконайте це в терміналі на Linux або macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт завантажує випуск для вашої системи, перевіряє його контрольну суму SHA-256 і розміщує `ccam` у `~/.local/bin`; там само ccam згодом створює вибрані вами короткі команди, наприклад `ccwork`. Якщо цієї теки немає у вашому `PATH`, скрипт виводить рядок, який потрібно додати до `~/.bashrc` або `~/.zshrc`.

У Windows виконайте це в PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Він завантажує `ccam_windows_<arch>.zip`, перевіряє його контрольну суму SHA-256, кладе `ccam.exe` у `%USERPROFILE%\.local\bin`, теку, куди інсталятор Claude Code кладе `claude.exe`, і додає цю теку до вашого користувацького `PATH`, якщо її там ще немає. Після цього відкрийте новий термінал.

ccam запускає Claude Code за вас, тому спочатку встановіть [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) і переконайтеся, що `claude --version` працює.

## Підтримувані системи

- Linux на x86-64 і ARM64: Ubuntu, Debian, Fedora, Arch, Alpine та інші дистрибутиви. Двійковий файл статичний і не має залежностей.
- macOS на Apple silicon та Intel.
- Windows 10 і 11 на x86-64 та ARM64, у PowerShell, Windows Terminal або Git Bash. Спільні теки працюють для будь-якого користувача. Спільним файлам (`settings.json`, `CLAUDE.md`, `history.jsonl`) потрібен режим розробника (Параметри > Система > Для розробників), бо Windows дозволяє створювати посилання на файли лише в цьому режимі або адміністраторові; див. [спільне використання у Windows](/features/share#on-windows).
- Усередині WSL ccam працює так само, як у будь-якій системі Linux.

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

Він додає навичку, команди `/ccam:usage`, `/ccam:ext` і `/ccam:switch`, а також два хуки: коли сеанс зупиняється на ліміті використання, ви отримуєте сповіщення в терміналі, а `/ccam:switch [account]` передає розмову іншому обліковому запису ([подробиці](/features/next#inside-claude-code)). Хуки запускають `ccam`, тому він має бути в `PATH`, який бачить Claude Code. Навичка ніколи не виконує команд, які виконують вхід або переміщують файли; натомість вона їх пропонує.

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

ccam ніколи не перевіряє оновлення самостійно: новий випуск він шукає лише тоді, коли ви запускаєте цю команду. Щоб встановити певний випуск, зокрема старіший, передайте `--version`. `ccam --dry-run update` показує, що буде завантажено й замінено, нічого не змінюючи.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Повторний запуск скрипта встановлення також оновлює копію в `~/.local/bin`. Якщо ви зібрали ccam за допомогою `go install`, виконайте ту саму команду ще раз.

**Windows.** `ccam update` завантажує `.zip` для Windows. У Windows не можна замінити запущену програму, тому старий виконуваний файл перейменовується на `ccam.exe.old`, і наступне оновлення його видаляє.

**Пакет для Debian і Ubuntu.** `ccam update` не замінює `/usr/bin/ccam`, встановлений через apt, а натомість виводить ці кроки. Завантажте новий `ccam_<version>_linux_amd64.deb` зі [сторінки випусків](https://github.com/jhyunwoo/ccam/releases) і встановіть його поверх старого:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Видалення

Спершу видаліть короткі команди облікових записів, поки ccam ще встановлено. `ccam list` показує їх у стовпці Shortcut; видаліть кожну командою `ccam shortcut <account> --remove` або самостійно видаліть ці посилання з `~/.local/bin`. Потім видаліть програму, а також теки з даними, якщо облікові записи вам більше не потрібні:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

У Windows, у PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Видалення `~/.local/share/ccam` також видаляє дані входу та історію розмов, що там зберігаються. Тека `~/.claude`, яку ви [імпортували](/docs/commands), залишається там, де вона є.
