---
title: Аккаунты по папкам
description: ccam folder заставляет папку и всё, что в ней, использовать один аккаунт, когда вы запускаете Claude Code без имени аккаунта, например work для кода компании.
updated: 2026-10-23
order: 11
translation: machine
sourceHash: "48de98d42b99"
---

`ccam folder work ~/code/company` делает так, что каждый запуск без имени аккаунта внутри `~/code/company` и в любой вложенной папке запускает Claude Code от имени `work`. За её пределами, как и раньше, запускается аккаунт по умолчанию. Если указать аккаунт, через `ccam personal` или короткую команду, например `ccpersonal`, всегда запускается именно он.

## Использование

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Аккаунт папки используют такие запуски:

- `ccam -p "…"` и любой другой запуск, первый аргумент которого является параметром Claude Code
- `ccam run` без имени и `ccam --dry-run`
- `ccam --fallback -p "…"` для первой попытки
- панель, `ccam` без аргументов, которая открывается с курсором на аккаунте папки, в том числе когда вы попали в папку через символическую ссылку

## Параметры

| Параметр | Действие | По умолчанию |
| --- | --- | --- |
| `<account>` | Имя или номер аккаунта, который нужно использовать | |
| `[folder]` | Папка, относительный или абсолютный путь | текущая папка |
| `--remove` | Удалить запись именно для этой папки | выкл. |
| `--dry-run` | Учитывается: вывести изменение, не сохраняя его | выкл. |

## Как это работает

Каждая папка — это одна запись `[[folder]]` в `~/.config/ccam/config.toml`, файле, который каждый запуск читает в любом случае:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` превращает папку в абсолютный путь с раскрытыми символическими ссылками и сохраняет его. Если задать папку ещё раз, её аккаунт заменяется.
2. Запуск без имени вызывает `launcher.HereProfile`. Если записей нет, он сразу берёт аккаунт по умолчанию.
3. Иначе он сравнивает рабочий каталог с каждым `path` как строки: подходит сама папка или вложенная в неё, но никогда не папка, имя которой лишь начинается с тех же букв (`/code/company` не охватывает `/code/companyx`). Побеждает самая глубокая запись, поэтому `~/code` может быть `personal`, а `~/code/company` — `work`.
4. Только если рабочий каталог в исходном виде ни с чем не совпал, ccam раскрывает его символические ссылки и сравнивает ещё раз.
5. В Windows и macOS сравнение не учитывает регистр букв, как обычно и их файловые системы.

Кроме `config.toml`, никакие файлы не читаются, поэтому запуск внутри папки не занимает измеримого времени (4,0 мс с записями и без них, `hyperfine`, 300 запусков).

## Безопасность и особые случаи

- `ccam rename` переносит записи папок на новое имя, а `ccam remove` удаляет записи удалённого аккаунта.
- Запись, аккаунт которой больше не существует, останавливает запуск без имени в этой папке: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` сообщает о таких записях, а также о записях, папки которых больше нет.
- `ccam folder --remove` внутри папки, которая получает аккаунт от папки выше, называет эту папку, а не удаляет что-либо.
- `ccam next` выбирает по ротации или по использованию и не смотрит на папки.

## См. также

- [Запуск Claude Code](/features/launch)
- [Короткие команды](/features/shortcuts)
- [Управление аккаунтами](/features/accounts)
- [Справочник команд](/docs/commands)
