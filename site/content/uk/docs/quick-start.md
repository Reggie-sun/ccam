---
title: Швидкий старт
metaTitle: Додати й перемкнути обліковий запис Claude Code
description: Додайте два облікові записи Claude Code командою ccam add, дайте їм короткі команди на зразок ccwork, перемикайтеся між ними й виберіть, що в них спільне.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "e98d0637356a"
---

Додайте кожен обліковий запис командою `ccam add`, а потім запускайте Claude Code командою `ccam <name>` або короткою командою з назвою, яку ви вибираєте самі, наприклад `ccwork`. Більшості цього достатньо.

## 1. Додайте перший обліковий запис

```sh
ccam add
```

Майстер запитує назву, наприклад `personal` або `work`, які елементи обліковий запис має спільно використовувати з іншими вашими обліковими записами і чи створити коротку команду (див. крок 4). Потім він виконує `claude auth login`, що відкриває браузер, щоб ви могли увійти. На віддаленій машині браузер може зупинитися на сторінці `http://localhost:…/callback`, яка не завантажується; вставте цю адресу в термінал, і ccam завершить вхід. Після завершення ви побачите:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Перенесіть наявні налаштування

Якщо ви вже користувалися Claude Code, ваші налаштування, `CLAUDE.md`, розширення та сервери MCP містяться в `~/.claude`. Скопіюйте їх до спільних налаштувань, щоб вони були в кожному обліковому записі:

```sh
ccam migrate
```

ccam показує, що він скопіює, і спершу запитує дозволу. Він лише читає `~/.claude`, ніколи не копіює дані входу й не переносить історію розмов, якщо ви не додасте `--history`. Подробиці див. у розділі [копіювання наявного налаштування](/docs/sharing#copy-an-existing-setup).

## 3. Додайте другий обліковий запис

```sh
ccam add work
```

Увійдіть у браузері з іншим обліковим записом.

## 4. Перемикання облікових записів

`ccam work` запускає Claude Code від імені робочого облікового запису. Якщо хочеться коротше, дайте обліковому запису власну коротку команду, наприклад `ccwork`. Після вибору спільних елементів майстер запитує:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Натисніть `y`, а потім введіть назву або натисніть Enter, щоб прийняти запропоновану:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Після входу ccam це підтверджує:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Якщо ви пропустили це запитання, коротку команду можна додати будь-коли:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Тепер будь-яка з цих команд запускає Claude Code від імені робочого облікового запису:

```sh
ccwork
ccam work
ccam 2
```

Якщо запустити лише `ccam`, відкриється [панель керування](/features/dashboard): перейдіть до облікового запису клавішами зі стрілками або його номером, а потім натисніть Enter.

Усе, що стоїть після короткої команди або назви облікового запису, передається в Claude Code без змін, тож ви й далі користуєтеся звичними параметрами:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Виберіть, що спільно використовувати

За замовчуванням між обліковими записами спільно використовуються сервери MCP, `settings.json`, глобальний `CLAUDE.md`, а також агенти, команди, навички та плагіни через групу `default`. Історія розмов не є спільною. Додайте її до групи, якщо хочете продовжувати розмови між обліковими записами:

```sh
ccam group set default --add history
```

ccam показує, що він перемістить для кожного облікового запису, перш ніж щось змінити. Подробиці див. у розділі [спільне використання через групи](/features/share).

## 6. Коли обліковий запис досягає ліміту

```sh
ccam usage
ccam next --continue
```

`ccam usage` показує 5-годинне та тижневе використання всіх облікових записів. `ccam next --continue` відкриває ту саму розмову в наступному обліковому записі. У скриптах `ccam --fallback work -p "…"` перемикається сам. Див. [ліміти використання](/docs/usage-limits).

**Місця в планах Team і Enterprise.** Автоматичний вибір (`ccam next`, `--fallback`, `ccam pick`, `n` на панелі керування) також перемикає між місцями однієї організації Claude Team або Enterprise. У довідковому центрі Anthropic сказано: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (кожне місце в плані Team закріплене за однією конкретною особою, і користуватися ним може лише ця особа) та «Users can't route their work through another user's account to use that user's remaining usage» (користувачі не можуть спрямовувати свою роботу через обліковий запис іншого користувача, щоб використати залишок ліміту цього користувача) ([плани Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); так само для [планів Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Тому використання автоматичного вибору між місцями може порушувати умови Anthropic. Перш ніж на нього покладатися, уточніть це у своїй організації та в Anthropic. Про особисті плани (Pro, Max) ідеться в [примітці про умови в FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Уже користуєтеся ~/.claude?

Зареєструйте її замість повторного входу:

```sh
ccam import
```

Тека залишається там, де вона є, а звичайний `claude` працює як і раніше.
