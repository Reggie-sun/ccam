---
title: Встановлення ccam
description: Встановіть ccam на Linux або macOS однією командою, з пакета .deb або з вихідного коду. Спочатку має бути встановлено Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Виконайте це в терміналі на Linux або macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Скрипт завантажує випуск для вашої системи, перевіряє його контрольну суму SHA-256 і розміщує дві команди в `~/.local/bin`: `ccam` і його короткий псевдонім `ccs`. Якщо цієї теки немає у вашому `PATH`, скрипт виводить рядок, який потрібно додати до `~/.bashrc` або `~/.zshrc`.

ccam запускає Claude Code за вас, тому спочатку встановіть [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) і переконайтеся, що `claude --version` працює.

## Підтримувані системи

- Linux на x86-64 і ARM64: Ubuntu, Debian, Fedora, Arch, Alpine та інші дистрибутиви. Двійковий файл статичний і не має залежностей.
- macOS на Apple silicon та Intel.
- Windows безпосередньо не підтримується. Усередині WSL ccam працює так само, як у будь-якій системі Linux.

## Інші способи встановлення

**Пакет для Debian і Ubuntu.** Завантажте `ccam_<version>_linux_amd64.deb` зі [сторінки випусків](https://github.com/jhyunwoo/ccam/releases) і встановіть його:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Пакет встановлює `/usr/bin/ccam` і посилання `/usr/bin/ccs`.

**З вихідного коду.** З Go 1.26 або новішою версією:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**Певна версія.** Задайте `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Автодоповнення в оболонці

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Автодоповнення працює і для `ccam`, і для `ccs`, включно з назвами облікових записів.

## Перевірка встановлення

```sh
ccam version
ccam doctor
```

`ccam doctor` перевіряє, що Claude Code знайдено і що кожна тека облікового запису, посилання та вхід у порядку. Далі перейдіть до [швидкого старту](/docs/quick-start).

## Видалення

Видаліть дві команди та короткі команди облікових записів, а потім теки з даними, якщо облікові записи вам більше не потрібні:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Видалення `~/.local/share/ccam` також видаляє дані входу та історію розмов, що там зберігаються. Тека `~/.claude`, яку ви [імпортували](/docs/commands), залишається там, де вона є.
