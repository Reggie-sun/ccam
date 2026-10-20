---
title: Короткие команды
description: Дайте аккаунту собственную команду, например ccwork, через ccam shortcut и заново создавайте ссылки в ~/.local/bin командой ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Короткая команда представляет собой команду с выбранным вами именем, которая запускает Claude Code от имени одного аккаунта: `ccwork --resume` делает то же, что и `ccam work --resume`. Это ссылка на ccam в `~/.local/bin`, поэтому она работает в bash, zsh, скриптах и везде, где программы запускаются из `PATH`. Короткие команды необязательны; создайте их для аккаунтов, которые вы запускаете часто.

## Использование

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

ccam подтверждает каждое изменение:

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

У аккаунта может быть не больше одной короткой команды, она хранится как `shortcut` в `~/.config/ccam/config.toml`. Мастер `ccam add` спрашивает `Create a shortcut command for <name>? [y/N]`, а затем `Shortcut name [cc<name>]:`. Вне мастера `ccam add` создаёт её только с `--shortcut`. `ccam list` показывает короткую команду каждого аккаунта или `-`.

Короткая команда передаёт в Claude Code каждый аргумент так же, как `ccam <name>` после имени. Это относится и к `--dry-run` и `--help`: чтобы заранее посмотреть запуск, используйте `ccam --dry-run work`.

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `--remove` | Удалить короткую команду аккаунта и её ссылку | выкл. |
| `-h`, `--help` | Справка по команде | |
| `--dry-run` | `ccam shortcut` учитывает его: проверяет имя и выводит `Would set the shortcut of work to ccwork.` или `Would remove ccwork, the shortcut of work.`, ничего не сохраняя и не создавая ссылок | выкл. |

`ccam shims install` не принимает аргументов и отклоняет `--dry-run` с сообщением `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Правила для имени, проверяемые в таком порядке:

- от 1 до 64 символов из `a-z`, `0-9`, `.`, `-` и `_`, начинается с буквы или цифры;
- не `ccam`, не `claude` и не команда ccam, например `run` или `usage`;
- не короткая команда другого аккаунта;
- не файл в `~/.local/bin`, который ещё не является ссылкой на ccam;
- не имя другой исполняемой программы в вашем `PATH`.

Если указать и команду, и `--remove`, выполнение останавливается с сообщением `give a command or --remove, not both`. Если задать имя, которое у аккаунта уже есть, выводится `work: ccwork (unchanged)`, а отсутствующая ссылка создаётся заново.

## Как это работает

Задание имени запускает `checkShortcut`, который вызывает `shim.Check` с папкой bin и реальным путём к исполняемому файлу ccam:

<!-- motion:shortcuts-check -->

Имя проходит проверки по порядку. Рядом с каждой проверкой показано имя, которое она отклоняет: `Work` не соответствует шаблону, `run` является командой ccam, `ccside` принадлежит другому аккаунту, `ccfile` уже существует как файл в `~/.local/bin`, а `git` скрыл бы `/usr/bin/git`. Сохраняется и связывается только имя, прошедшее все проверки.

1. `shim.Check` сопоставляет имя с `^[a-z0-9][a-z0-9._-]{0,63}$`, затем с `config.Reserved`, затем с `cfg.FindShortcut`. Он вызывает `os.Lstat` для `~/.local/bin/<name>`, а для каждой другой папки в `PATH` вызывает `os.Stat` и отклоняет исполняемый файл, который не является ссылкой на ccam.
2. ccam сохраняет `config.toml` через `fsutil.AtomicWrite` (временный файл, `fsync`, переименование).
3. `shim.Remove` удаляет ссылку старой короткой команды, но только если это ссылка на ccam.
4. `shim.Install` создаёт ссылку для каждой настроенной короткой команды, у которой её нет, через `os.Symlink` на исполняемый файл ccam, путь к которому определяется через `os.Executable` и `filepath.EvalSymlinks`. Он никогда ничего не удаляет.

Ссылка считается собственной ссылкой ccam, если `os.Readlink` возвращает путь к ccam или если ссылка и путь к ccam указывают на один и тот же файл. Каждое удаление и каждая проверка перед перезаписью проходят через эту проверку, поэтому ccam не трогает файлы, которые создал не он.

Когда вы запускаете короткую команду, ccam нужно только набранное вами имя:

<!-- motion:shortcuts-link -->

Оболочка находит `ccwork` в `PATH`, ядро следует по ссылке и запускает исполняемый файл ccam, а `argv[0]` по-прежнему равен `ccwork`. В `app.Main` любое имя, кроме `ccam`, один раз загружает `config.toml` и вызывает `shim.Lookup`, который находит аккаунт, у которого `shortcut` равен `filepath.Base(argv[0])`. Затем `launcher.LaunchWith` выполняет те же шаги, что и [`ccam <name>`](/features/launch), с уже загруженной конфигурацией и передаёт `argv[1:]` в `claude` без изменений.

Почему символическая ссылка, а не псевдоним оболочки или скрипт-обёртка: символическая ссылка работает в любой оболочке и в скриптах, не добавляет лишнего запуска оболочки (около 1–2 мс), а запуск по-прежнему читает один файл. `scripts/perf-check.sh` измеряет короткую команду с помощью `hyperfine` с тем же пределом 10 мс (p95), что и для `ccam <name>`.

### В Windows

Короткая команда — это `ccwork.exe` в `%USERPROFILE%\.local\bin`, а `shim.Name` отбрасывает `.exe` (в любом регистре) перед поиском, поэтому `ccwork`, `ccwork.exe` и `CCWORK.EXE` запускают `work`. Сначала ccam создаёт символическую ссылку, что Windows разрешает в режиме разработчика или администратору. Без этого права он создаёт жёсткую ссылку: второе имя того же файла `ccam.exe`, для которого привилегии не нужны. Жёсткая ссылка считается ссылкой ccam, если это тот же файл, что и `ccam.exe`. Поскольку `ccam update` ставит новый файл на место `ccam.exe`, после этого он обновляет каждую короткую команду с жёсткой ссылкой через `shim.Relink` и выводит `✓ shortcut ccwork renewed`. Проверка имени ищет `<name>` с каждым расширением из `PATHEXT` (`.exe`, `.cmd` и так далее) в остальных папках `PATH`.

## Безопасность и особые случаи

- **Ссылка, которая не принадлежит ни одному аккаунту.** Если у ссылки на ccam в `~/.local/bin` имя, которое не использует ни один аккаунт, например после ручного изменения `shortcut` в `config.toml`, ccam не переходит на аккаунт по умолчанию. Он останавливается с сообщением `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Задайте имя заново через `ccam shortcut` или удалите ссылку.
- **Существующие файлы.** ccam никогда не заменяет файл, который создал не он. `ccam shims install` останавливается с сообщением `<path> already exists and is not a ccam link, so it was left alone`.
- **Перемещённый исполняемый файл ccam.** Старые ссылки тогда указывают на путь, которого больше нет, и ccam не считает их своими. Удалите их, затем выполните `ccam shims install`.
- **Переименование и удаление.** `ccam rename` сохраняет имя короткой команды, и ссылка продолжает работать, потому что она указывает на ccam, а ccam ищет имя в `config.toml`. `ccam remove` удаляет ссылку аккаунта. `ccam shortcut <account> <new>` удаляет старую ссылку, прежде чем создать новую.
- **PATH.** Если ccam создаёт ссылку, а `~/.local/bin` нет в вашем `PATH`, он выводит `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` и следом строку `export PATH=...`, которую нужно добавить. В Windows вместо этого он выводит строку PowerShell, которая добавляет папку в ваш пользовательский `PATH`.
- **Последующая проверка.** `ccam doctor` сообщает `missing shortcuts in /home/me/.local/bin: ccwork` и предупреждает, если папки нет в `PATH`. `ccam doctor --fix` создаёт отсутствующие ссылки.

## См. также

- [Запуск Claude Code](/features/launch)
- [Управление аккаунтами](/features/accounts)
- [Справочник команд](/docs/commands#shortcuts)
- [Установка ccam](/docs/install)
