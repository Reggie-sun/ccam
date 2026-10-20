---
title: Копіювання наявного налаштування
description: Скопіюйте налаштування, CLAUDE.md, розширення, сервери MCP і входи MCP з наявної ~/.claude до групи ccam за допомогою ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

Якщо ви користувалися Claude Code до ccam, ваше налаштування міститься в `~/.claude` і `~/.claude.json`. `ccam migrate` копіює його до спільної групи, тож кожен обліковий запис, який спільно використовує ці елементи через групу, починає з налаштування, яке у вас уже є. Вихідна тека лише читається, а звичайний `claude` і далі працює як раніше. Виконайте команду один раз після першого `ccam add` і знову щоразу, коли якийсь обліковий запис починає спільно використовувати сервери MCP.

## Використання

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

У терміналі, без прапорців елементів і без `--yes`, ccam спершу показує контрольний список:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (або `x`) перемикає елемент під курсором, клавіші зі стрілками або `j`/`k` переміщують курсор, Enter показує план, а Esc, `q` або Ctrl+C виводить `Cancelled.`. Елементи, для яких у джерелі нічого немає, позначено `not found`. Якщо зняти всі позначки, ccam виводить `Nothing selected.` і зупиняється.

Потім ccam виводить план і запитує підтвердження перед копіюванням:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Група доходить лише до власних облікових записів і лише для елементів, які вона робить спільними, тому підсумок перелічує облікові записи групи. Коли група не робить спільним один зі скопійованих елементів або деякі облікові записи до неї не входять, ccam виводить команди `ccam group set <group> --add <items>` і `ccam group join <group> <accounts>`, які це змінюють.

## Параметри

| Параметр | Дія | За замовчуванням |
| --- | --- | --- |
| `[folder]` | Тека Claude Code, з якої копіювати | `~/.claude` |
| `--settings` | Копіювати `settings.json` | вимкнено |
| `--memory` | Копіювати глобальний `CLAUDE.md` | вимкнено |
| `--extensions` | Копіювати `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | вимкнено |
| `--mcp` | Копіювати користувацькі сервери MCP (`mcpServers` у `.claude.json`) | вимкнено |
| `--history` | Копіювати історію розмов (`projects/`, `history.jsonl`) | вимкнено |
| `-g`, `--group string` | Спільна група, до якої копіювати | `default` |
| `-y`, `--yes` | Копіювати без запитань; у разі конфлікту зберігається копія групи | вимкнено |
| `--dry-run` (глобальний) | Вивести план і зупинитися перед будь-яким запитанням чи копіюванням | вимкнено |

Без прапорця елемента в терміналі вирішує контрольний список. Поза терміналом або з `--yes` ccam копіює все, крім історії розмов. Історію за замовчуванням пропущено, бо `projects/` може бути великою.

`--group on` означає `default`. `--group off` і `--group local` відхиляються з повідомленням `--group takes a group name, not "off"`. `ccam migrate` враховує `--dry-run`, і контрольний список з'являється й під час пробного запуску, тож можна попередньо переглянути вибір.

## Як це працює

`internal/cli/migrate.go` вибирає елементи, потім `migrate.Build` читає джерело та групу й повертає план. Нічого не записується до `Plan.Apply`, тобто доки ви не підтвердите.

1. **Пошук джерела.** `checkSource` відхиляє шлях, який відсутній, не є текою або не містить жодного відомого запису Claude Code. `migrate.ClaudeJSON` знаходить `.claude.json`, що належить теці: усередині неї або, для `~/.claude` без власного файлу всередині, `~/.claude.json` поруч із нею. Саме там Claude Code зберігає його, коли `CLAUDE_CONFIG_DIR` не задано.
2. **Зіставлення елементів зі шляхами.** Кожен шлях береться з `share.Layout`, того самого списку дозволених шляхів, який використовують групи. Сервери MCP зчитуються з `.claude.json` через `mcp.LocalServers`, який використовує сканер `profile.TopLevelValue` замість розбору всього файлу.
3. **Копіювати, ніколи не переміщувати.** Файли потрапляють до `~/.local/share/ccam/groups/<group>/` через `fsutil.AtomicWrite` або `fsutil.CopyFile`; сервери MCP потрапляють до `mcp.json` групи.

Кожен елемент зіставляє один вихідний шлях з одним шляхом групи:

<!-- motion:migrate-items -->

Два рядки історії позначено пунктиром, бо вони копіюються лише з `--history`. `.credentials.json` і `.claude.json` перелічено в `share.Never`, і вони ніколи не потрапляють до групи: `Build` повертає внутрішню помилку, якщо хтось із них колись з'явиться в структурі, а обхід теки пропускає їх за назвою.

`buildEntry` призначає кожному шляху стан, і саме його описує кожен рядок плану:

<!-- motion:migrate-states -->

- **Missing**: немає в джерелі або файл там, де очікується тека. Виводиться один раз для елемента як `not found, skipped`.
- **Linked**: вихідний шлях уже вказує на копію групи, наприклад коли джерелом є тека профілю ccam.
- **Same**: група вже має ідентичний вміст.
- **New**: у групі цього немає, тож його копіюють повністю.
- **Merge**: тека з обох боків. `copyMissing` додає лише файли, яких бракує групі; відмінні файли залишаються такими, як у групі, і план повідомляє `group copies kept`. Для `plugins/` реєстри `installed_plugins.json` і `known_marketplaces.json` об'єднуються, а абсолютні шляхи до плагінів переписуються так, щоб вказувати на теку групи.
- **Conflict**: файл з обох боків із різним вмістом. Ви вибираєте для кожного файлу: `keep group`, `use mine`, а також `merge JSON` для файлів `.json` (глибоке об'єднання, у якому при конфлікті значень перемагає група) або `merge by line` для `history.jsonl`. Сервер MCP з відмінним визначенням теж є конфліктом, з варіантами `keep group` або `use mine`. Перший наведений варіант діє за замовчуванням, коли ви натискаєте Enter; `--yes` завжди зберігає копію групи.

Сам контрольний список — це `tui.MigratePicker`, побудований на компоненті `checklist`, який використовує `ccam add`; `migrate.Found` позначає відсутні елементи за допомогою викликів lstat і списку серверів MCP.

### Входи MCP

Сервери MCP, що виконують вхід через OAuth, зберігають свої токени у файлі входу, `.credentials.json` у `mcpOAuth`, а не в списку серверів. Claude Code формує ключ кожного запису як `<name>|` плюс перші 16 символів хешу `sha256` від `type`, `url` і `headers` сервера. Ключ не залежить від теки, тож скопійований запис працює в іншому обліковому записі, якщо визначення сервера ідентичне. Тому `buildLogins` в `internal/migrate/logins.go` копіює лише записи серверів, які зрештою опиняються в групі з визначенням із джерела: нові, уже однакові або конфліктні, для яких ви вибрали `use mine`. Входи серверів плагінів (ключі, що починаються з `plugin:`) копіюються разом із розширеннями, якщо джерело має `plugins/`.

<!-- motion:migrate-signins -->

Записи потрапляють до власного `.credentials.json` кожного облікового запису-одержувача, записаного через `fsutil.AtomicWrite` з режимом `0600`. До групи вони ніколи не потрапляють. Обліковий запис отримує їх, лише якщо спільно використовує `mcp` (або, для серверів плагінів, `extensions`) через цільову групу, і додаються лише ключі, яких у нього немає: записи, які в нього вже є, можливо новіші, зберігаються. Рядок плану `sign-ins` перелічує сервери, облікові записи-одержувачі та пропущені. ccam пропускає обліковий запис і називає його з причиною, коли він:

- зареєстрований на місці через `ccam import`, бо звичайний `claude` може використовувати його файл входу;
- запущений, бо Claude Code переписує свій файл входу (перевіряється знову безпосередньо перед записом, повідомляється як `started meanwhile`);
- ще не має входу;
- зберігає вхід у в'язці ключів macOS, у яку ccam не записує.

Увійдіть у сервери в пропущеному обліковому записі через `/mcp`. Обліковий запис, який почне спільно використовувати сервери MCP пізніше, отримає входи, коли ви знову виконаєте `ccam migrate`; підсумок повідомляє про це рядком ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Безпека та граничні випадки

- Джерело лише читається. Жоден файл у ньому не переміщується, не пов'язується й не змінюється, і вхід джерела ніколи не записується, зокрема на macOS, де його можуть читати з в'язки ключів.
- Перш ніж файл групи буде замінено через `use mine`, об'єднання або нові сервери MCP у наявному `mcp.json`, копію групи зберігають у `<group>/.ccam-backup/<time>/`. Назва теки унікальна для кожного запуску, тож два запуски в ту саму секунду не перезаписують резервні копії один одного. Підсумок виводить шлях.
- З `--yes` кожен конфлікт зберігає копію групи, і ccam пояснює, як вибрати натомість: виконати ту саму команду без `--yes`.
- Якщо нічого не змінилося б, ccam виводить `Nothing to copy.`; якщо кожен конфлікт зберігає копію групи — `Nothing to copy: the group keeps its own copies.`
- Якщо копіювання зазнає невдачі посередині, ccam зупиняється з повідомленням `The source folder was not changed. Group files replaced so far are saved in`, після якого йде тека резервних копій.
- Відсутнє джерело завершується повідомленням `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Пов'язані сторінки

- [Спільні сервери MCP, налаштування та історія](/docs/sharing)
- [Перевірка та виправлення](/features/doctor)
- [Кілька облікових записів Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
