---
title: Облікові записи за теками
description: ccam folder змушує теку й усе в ній використовувати один обліковий запис, коли ви запускаєте Claude Code без його назви, наприклад work для коду компанії.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` робить так, що кожен запуск без назви облікового запису всередині `~/code/company` і в будь-якій вкладеній теці запускає Claude Code від імені `work`. Поза нею, як і раніше, запускається обліковий запис за замовчуванням. Якщо назвати обліковий запис, через `ccam personal` або коротку команду, як-от `ccpersonal`, завжди запускається саме він.

## Використання

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

Обліковий запис теки використовують такі запуски:

- `ccam -p "…"` та будь-який інший запуск, перший аргумент якого є параметром Claude Code
- `ccam run` без назви та `ccam --dry-run`
- `ccam --fallback -p "…"` для першої спроби
- вікно вибору, `ccam` без аргументів, яке відкривається з курсором на обліковому записі теки

## Параметри

| Параметр | Що робить | За замовчуванням |
| --- | --- | --- |
| `<account>` | Назва або номер облікового запису, який слід використовувати | |
| `[folder]` | Тека, відносний або абсолютний шлях | поточна тека |
| `--remove` | Видалити запис саме для цієї теки | вимкнено |
| `--dry-run` | Враховується: вивести зміну, не зберігаючи її | вимкнено |

## Як це працює

Кожна тека — це один запис `[[folder]]` у `~/.config/ccam/config.toml`, файлі, який кожен запуск читає в будь-якому разі:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` перетворює теку на абсолютний шлях із розкритими символьними посиланнями й зберігає його. Якщо задати теку ще раз, її обліковий запис замінюється.
2. Запуск без назви викликає `launcher.HereProfile`. Якщо записів немає, він одразу бере обліковий запис за замовчуванням.
3. Інакше він порівнює робочий каталог із кожним `path` як рядки: підходить сама тека або вкладена в неї, але ніколи не тека, назва якої лише починається з тих самих літер (`/code/company` не охоплює `/code/companyx`). Перемагає найглибший запис, тож `~/code` може бути `personal`, а `~/code/company` — `work`.
4. Лише якщо робочий каталог у початковому вигляді ні з чим не збігся, ccam розкриває його символьні посилання й порівнює ще раз.
5. У Windows порівняння не враховує регістр літер, як і сама Windows.

Крім `config.toml`, жодні файли не читаються, тож запуск усередині теки не забирає вимірного часу (4,0 мс із записами й без них, `hyperfine`, 300 запусків).

## Безпека та граничні випадки

- `ccam rename` переносить записи тек на нову назву, а `ccam remove` видаляє записи видаленого облікового запису.
- Запис, обліковий запис якого більше не існує, зупиняє запуск без назви в цій теці: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` повідомляє про такі записи, а також про записи, теки яких більше немає.
- `ccam folder --remove` усередині теки, яка отримує обліковий запис від теки вище, називає ту теку, а не видаляє щось.
- `ccam next` вибирає за ротацією або за використанням і не зважає на теки.

## Пов'язані сторінки

- [Запуск Claude Code](/features/launch)
- [Короткі команди](/features/shortcuts)
- [Керування обліковими записами](/features/accounts)
- [Довідник команд](/docs/commands)
