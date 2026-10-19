---
title: Установка ccam
description: Установите ccam на Linux или macOS одной командой, из пакета .deb или из исходного кода. Сначала должен быть установлен Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Выполните это в терминале на Linux или macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт скачивает выпуск для вашей системы, проверяет его контрольную сумму SHA-256 и помещает `ccam` в `~/.local/bin`; позже `ccam add` создаёт там же короткие команды `claude1`, `claude2` и так далее. Если этой папки нет в вашем `PATH`, скрипт выводит строку, которую нужно добавить в `~/.bashrc` или `~/.zshrc`.

ccam запускает Claude Code за вас, поэтому сначала установите [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) и убедитесь, что `claude --version` работает.

## Поддерживаемые системы

- Linux на x86-64 и ARM64: Ubuntu, Debian, Fedora, Arch, Alpine и другие дистрибутивы. Исполняемый файл статический и не имеет зависимостей.
- macOS на Apple silicon и Intel.
- Windows напрямую не поддерживается. Внутри WSL ccam работает так же, как в любой системе Linux.

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

**Определённая версия.** Задайте `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Автодополнение в оболочке

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Автодополнение работает для всех команд и параметров `ccam`, включая имена аккаунтов.

## Проверка установки

```sh
ccam version
ccam doctor
```

`ccam doctor` проверяет, что Claude Code найден и что все папки аккаунтов, ссылки и входы в порядке. Дальше переходите к [быстрому старту](/docs/quick-start).

## Удаление

Удалите программу ccam и короткие команды аккаунтов, а затем папки с данными, если аккаунты вам больше не нужны:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Удаление `~/.local/share/ccam` удаляет также данные входа и историю разговоров, которые там хранятся. Папка `~/.claude`, которую вы [импортировали](/docs/commands), остаётся на месте.
