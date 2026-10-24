---
title: Установка ccam
metaTitle: Установка ccam на Linux, macOS и Windows
description: Установите ccam на Linux, macOS или Windows одной командой, из пакета .deb или из исходного кода и обновляйте через ccam update. Сначала нужен Claude Code.
updated: 2026-10-23
order: 1
translation: machine
sourceHash: "af21dcc24169"
---

Выполните это в терминале на Linux или macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт скачивает выпуск для вашей системы, проверяет его контрольную сумму SHA-256 и помещает `ccam` в `~/.local/bin`; там же ccam позже создаёт выбранные вами короткие команды, например `ccwork`. Если этой папки нет в вашем `PATH`, скрипт выводит строку, которую нужно добавить в `~/.bashrc` или `~/.zshrc`.

В Windows выполните это в PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Он скачивает `ccam_windows_<arch>.zip`, проверяет его контрольную сумму SHA-256, помещает `ccam.exe` в `%USERPROFILE%\.local\bin`, папку, куда установщик Claude Code помещает `claude.exe`, и добавляет эту папку в ваш пользовательский `PATH`, если её там ещё нет, оставляя записи вроде `%USERPROFILE%\bin` как есть. После этого откройте новый терминал. Повторный запуск заменяет ccam, даже пока сессия Claude Code ещё работает из старого исполняемого файла (он переименовывается в `ccam.exe.old` или в `ccam.exe.old-<number>`, если более старый ещё используется), и обновляет короткие команды, которые являются жёсткими ссылками: `OK 2 shortcuts now run the new ccam`.

ccam запускает Claude Code за вас, поэтому сначала установите [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) и убедитесь, что `claude --version` работает.

## Поддерживаемые системы

- Linux на x86-64 и ARM64: Ubuntu, Debian, Fedora, Arch, Alpine и другие дистрибутивы. Исполняемый файл статический и не имеет зависимостей.
- macOS на Apple silicon и Intel.
- Windows 10 и 11 на x86-64 и ARM64, в PowerShell, Windows Terminal или Git Bash. Общие папки работают у любого пользователя. Общим файлам (`settings.json`, `CLAUDE.md`, `history.jsonl`) нужен режим разработчика (Параметры > Система > Для разработчиков), потому что Windows разрешает создавать ссылки на файлы только в этом режиме или администратору; см. [общий доступ в Windows](/features/share#on-windows). Без него первый `ccam add` создаёт группу `default` только с MCP-серверами и расширениями и сообщает, как добавить остальное после включения режима разработчика.
- Внутри WSL ccam работает так же, как в любой системе Linux.

## Другие способы установки

**Пакет для Debian и Ubuntu.** Скачайте `ccam_<version>_linux_amd64.deb` со [страницы выпусков](https://github.com/jhyunwoo/ccam/releases) и установите его:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Пакет устанавливает `/usr/bin/ccam`.

**Из исходного кода.** С Go 1.26 или новее:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

Сборка тега выпуска (`@latest` или `@vX.Y.Z`) сообщает эту версию, поэтому для неё работает `ccam update`; сборка `@main` или локальной копии новее тега сообщает `dev`.

**Определённая версия.** Задайте `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Автодополнение в оболочке

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Автодополнение работает для всех команд и параметров `ccam`, включая имена аккаунтов.

## Плагин Claude Code

Плагин ccam учит сам Claude Code пользоваться ccam: проверять, у какого аккаунта есть запас, запускать параллельные рабочие процессы `claude -p` в разных аккаунтах с помощью `ccam pick` и объяснять отсутствие навыка или MCP-сервера с помощью `ccam ext`. Установите его внутри Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Он добавляет навык, команды `/ccam:usage`, `/ccam:ext` и `/ccam:switch` и два хука: когда сессия останавливается на лимите использования, вы получаете уведомление в терминале, а `/ccam:switch [account]` передаёт разговор другому аккаунту ([подробнее](/features/next#inside-claude-code)). Хуки запускают `ccam`, поэтому он должен быть в `PATH`, который видит Claude Code. Навык никогда не выполняет команды, которые выполняют вход или перемещают файлы; вместо этого он их предлагает.

## Проверка установки

```sh
ccam version
ccam doctor
```

`ccam doctor` проверяет, что Claude Code найден и что все папки аккаунтов, ссылки и входы в порядке. Дальше переходите к [быстрому старту](/docs/quick-start).

## Обновление

```sh
ccam update --check
ccam update
```

`ccam update --check` только сообщает, есть ли более новый выпуск. `ccam update` скачивает выпуск для вашей системы, сверяет его контрольную сумму SHA-256 с файлом выпуска `checksums.txt` и затем за один шаг заменяет файл `ccam`. Если контрольная сумма не совпадает, ничего не меняется. Короткие команды, например `ccwork`, являются ссылками на ccam, поэтому сразу запускают новую версию.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

Панель (`ccam` в терминале) тоже ищет новый выпуск при каждом открытии и, если он есть, предлагает его: нажмите `U` и подтвердите. Больше ни одна команда его не ищет, и ccam никогда не обновляется, не спросив вас. Чтобы установить определённый выпуск, в том числе более старый, укажите `--version`. `ccam --dry-run update` показывает, что будет скачано и заменено, ничего не меняя.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Повторный запуск скрипта установки тоже обновляет копию в `~/.local/bin`. Если вы собрали ccam с помощью `go install` из тега выпуска, работает и `ccam update`, или выполните ту же команду ещё раз.

**Windows.** `ccam update` скачивает `.zip` для Windows. В Windows нельзя заменить работающую программу, поэтому старый исполняемый файл переименовывается в `ccam.exe.old`, и следующее обновление его удаляет; пока сессия ещё работает из более раннего `ccam.exe.old`, имя будет `ccam.exe.old-<time>`. Короткие команды, которые являются жёсткими ссылками, обновляются сразу после этого.

**Пакет для Debian и Ubuntu.** `ccam update` не заменяет `/usr/bin/ccam`, установленный через apt, а вместо этого выводит следующие шаги. Скачайте новый `ccam_<version>_linux_amd64.deb` со [страницы выпусков](https://github.com/jhyunwoo/ccam/releases) и установите его поверх старого:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Удаление

Сначала удалите короткие команды аккаунтов, пока ccam ещё установлен. `ccam list` показывает их в столбце Shortcut; удалите каждую командой `ccam shortcut <account> --remove` или сами удалите эти ссылки из `~/.local/bin`. Затем удалите программу, а также папки с данными, если аккаунты вам больше не нужны:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

В Windows, в PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Удаление `~/.local/share/ccam` удаляет также данные входа и историю разговоров, которые там хранятся. Папка `~/.claude`, которую вы [импортировали](/docs/commands), остаётся на месте, но пока она входит в группу, общие файлы лежат в `~/.local/share/ccam/groups/`. Сначала выполните для неё `ccam remove <name>` (или `ccam group leave`): она выйдет из своих групп с копией всего, что делала общим.
