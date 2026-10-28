---
title: Установка ccam
metaTitle: Установка ccam на Linux, macOS и Windows
description: Установите ccam на Linux, macOS или Windows одной командой, из пакета .deb или через go install и обновляйте через ccam update. Сначала нужен Claude Code.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam запускает Claude Code за вас, поэтому сначала установите [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) и убедитесь, что `claude --version` работает.

Выполните это в терминале на Linux или macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт скачивает выпуск для вашей системы, проверяет его контрольную сумму SHA-256 и помещает `ccam` в `~/.local/bin`; там же ccam позже создаёт выбранные вами короткие команды, например `ccwork`. Если этой папки нет в вашем `PATH`, скрипт выводит строку, которую нужно добавить в `~/.bashrc` или `~/.zshrc`.

В Windows выполните это в PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Он скачивает `ccam_windows_<arch>.zip`, проверяет его контрольную сумму SHA-256 и помещает `ccam.exe` в `%USERPROFILE%\.local\bin`, папку, куда установщик Claude Code помещает `claude.exe`. Если этой папки ещё нет в вашем пользовательском `PATH`, скрипт добавляет её, оставляя записи вроде `%USERPROFILE%\bin` как есть. После этого откройте новый терминал.

## Поддерживаемые системы

- Linux на x86-64 и ARM64: Ubuntu, Debian, Fedora, Arch, Alpine и другие дистрибутивы. Исполняемый файл статический и не имеет зависимостей.
- macOS на Apple silicon и Intel.
- Windows 10 и 11 на x86-64 и ARM64, в PowerShell, Windows Terminal или Git Bash:
  - Общие папки работают без режима разработчика.
  - Общим файлам (`settings.json`, `CLAUDE.md`, `history.jsonl`) нужен режим разработчика (Параметры > Система > Для разработчиков), потому что Windows разрешает создавать ссылки на файлы только в этом режиме или администратору; см. [общий доступ в Windows](/features/share#on-windows).
  - Без режима разработчика первый `ccam add` создаёт группу `default` только с MCP-серверами и расширениями и сообщает, как добавить остальное после включения режима разработчика.
- Внутри WSL ccam работает так же, как в любой другой системе Linux.

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

Если вы устанавливаете тег выпуска (`@latest` или `@vX.Y.Z`), ccam сообщает эту версию, и работают `ccam update` и предложение обновления на панели. Сборка `@main` или локальной копии с коммитами после тега сообщает `dev`, и `ccam update` отказывается её обновлять, если не указать `--version`.

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

Он добавляет:

- навык, который предлагает команды для входа или перемещения файлов, но никогда их не выполняет;
- команды `/ccam:usage`, `/ccam:ext` и `/ccam:switch`;
- два хука ([подробнее](/features/next#inside-claude-code)). Один показывает уведомление в терминале, когда сессия останавливается на лимите использования; в Linux и macOS при включённом [`auto_switch`](/features/next#auto-switch) ccam затем переносит сессию в другой аккаунт, где есть запас. Другой позволяет `/ccam:switch [account]` передать разговор другому аккаунту.

Хуки запускают `ccam`, поэтому он должен быть в `PATH`, который видит Claude Code.

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

Панель (`ccam` без аргументов, в терминале) тоже проверяет наличие нового выпуска при каждом открытии и, если он есть, предлагает его: нажмите `U` и подтвердите. Больше ни одна команда не проверяет обновления, и ccam никогда не обновляется, не спросив вас. Чтобы установить определённый выпуск, в том числе более старый, укажите `--version`. `ccam --dry-run update` показывает, что будет скачано и заменено, ничего не меняя.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Повторный запуск скрипта установки тоже обновляет копию в `~/.local/bin`. Если вы собрали ccam с помощью `go install` из тега выпуска, можно выполнить `ccam update` или ту же команду `go install` ещё раз.

**Windows.** `ccam update` скачивает `.zip` для Windows. В Windows нельзя заменить работающую программу, поэтому старый исполняемый файл переименовывается в `ccam.exe.old`, и следующее обновление его удаляет; пока сессия ещё работает из более раннего `ccam.exe.old`, имя будет `ccam.exe.old-<time>`. Короткие команды, которые являются жёсткими ссылками, обновляются сразу после этого. Повторный запуск `install.ps1` работает так же, даже пока сессия Claude Code ещё работает из старого исполняемого файла, и выводит `OK 2 shortcuts now run the new ccam`, когда обновляет короткие команды — жёсткие ссылки.

**Пакет для Debian и Ubuntu.** `ccam update` не заменяет `/usr/bin/ccam`, установленный через apt, а вместо этого выводит следующие шаги. Скачайте новый `ccam_<version>_linux_amd64.deb` со [страницы выпусков](https://github.com/jhyunwoo/ccam/releases) и установите его поверх старого:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Удаление

`ccam uninstall` удаляет ccam и всё, что он создал. Сначала он показывает, что будет удалено, и спрашивает подтверждение:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Каждый аккаунт удаляется так же, как это делает `ccam remove`. Из аккаунта в папках ccam сначала выполняется выход, затем он удаляется вместе со своей папкой, а значит, и с историей разговоров и настройками.
- Папка `~/.claude`, которую вы [импортировали](/docs/commands), остаётся на месте для обычного `claude`. Если она входит в группу, она получает обратно копию всего, что делала общим. Строка состояния ccam в её `settings.json` уступает место строке состояния, которую ccam оборачивал, если такая была. Если в ней есть плагин ccam, ccam выводит команду `/plugin uninstall`, которую нужно выполнить в Claude Code.
- Короткие команды, `~/.config/ccam`, `~/.local/share/ccam` (вместе с корзиной ccam) и `~/.local/state/ccam` удаляются, а затем и сама команда `ccam`.

ccam отказывается работать, пока запущен какой-либо аккаунт; сначала завершите его сессии Claude Code. `/usr/bin/ccam`, установленный через apt, остаётся; удалите его потом командой `sudo apt remove ccam`. В Windows запущенная программа не может удалить сама себя, поэтому в конце `ccam uninstall` выводит команду `Remove-Item` для `ccam.exe`.
