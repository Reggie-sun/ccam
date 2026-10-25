---
title: Быстрый старт
metaTitle: "Быстрый старт: добавить и сменить аккаунт Claude Code"
description: Добавьте два аккаунта Claude Code командой ccam add, задайте им короткие команды вроде ccwork, переключайтесь между ними и выберите, что у них общее.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

Добавьте каждый аккаунт командой `ccam add`, затем запускайте Claude Code командой `ccam <name>` или короткой командой, имя которой вы выбираете сами, например `ccwork`. Большинству этого достаточно.

## 1. Добавьте первый аккаунт

```sh
ccam add
```

Мастер спрашивает имя, например `personal` или `work`, какие элементы аккаунт должен делить с другими вашими аккаунтами и нужно ли создать короткую команду (см. шаг 4). Затем он выполняет `claude auth login`, который открывает браузер, чтобы вы могли войти. На удалённой машине браузер может остановиться на странице `http://localhost:…/callback`, которая не загружается; вставьте этот адрес в терминал, и ccam завершит вход. По завершении вы увидите:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Перенесите существующие настройки

Если вы уже пользовались Claude Code, ваши настройки, `CLAUDE.md`, расширения и MCP-серверы находятся в `~/.claude`. Скопируйте их в общие настройки, чтобы они были у каждого аккаунта:

```sh
ccam migrate
```

ccam показывает, что он скопирует, и сначала спрашивает подтверждение. Он только читает `~/.claude`, никогда не копирует данные входа и не переносит историю разговоров, если не добавить `--history`. Подробности в разделе [копирование существующей настройки](/docs/sharing#copy-an-existing-setup).

## 3. Добавьте второй аккаунт

```sh
ccam add work
```

Войдите в браузере с другим аккаунтом.

## 4. Переключение аккаунтов

`ccam work` запускает Claude Code от имени рабочего аккаунта. Если хочется короче, дайте аккаунту собственную короткую команду, например `ccwork`. После выбора общих элементов мастер спрашивает:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Нажмите `y`, затем введите имя или нажмите Enter, чтобы принять предложенное:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

После входа ccam подтверждает результат:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Если вы пропустили этот вопрос, короткую команду можно добавить в любой момент:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Теперь любая из этих команд запускает Claude Code от имени рабочего аккаунта:

```sh
ccwork
ccam work
ccam 2
```

Если выполнить `ccam` без аргументов, откроется [панель](/features/dashboard): перейдите к аккаунту клавишами со стрелками или его номером и нажмите Enter.

Всё, что указано после короткой команды или имени аккаунта, передаётся в Claude Code без изменений, поэтому вы продолжаете пользоваться привычными параметрами:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Выберите, что будет общим

По умолчанию MCP-серверы, `settings.json`, глобальный `CLAUDE.md`, а также агенты, команды, навыки и плагины общие для всех аккаунтов через группу `default`. История разговоров не общая. Добавьте её в группу, если хотите продолжать разговоры в другом аккаунте:

```sh
ccam group set default --add history
```

Прежде чем что-либо изменить, ccam показывает, что он переместит для каждого аккаунта. Подробности в разделе [общий доступ через группы](/features/share).

## 6. Когда аккаунт достиг лимита

```sh
ccam usage
ccam next --continue
```

`ccam usage` показывает использование за 5 часов и за неделю для каждого аккаунта. `ccam next --continue` открывает тот же разговор в следующем аккаунте. В скриптах `ccam --fallback work -p "…"` переключается сам. См. раздел [лимиты использования](/docs/usage-limits).

**Места в тарифах Team и Enterprise.** При автоматическом выборе (`ccam next`, `--fallback`, `ccam pick`, клавиша `n` на панели, а также [`auto_switch`](/features/next#auto-switch), если вы его включите) ccam переключается и между местами одной и той же организации Claude Team или Enterprise. В справочном центре Anthropic ([тарифы Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); аналогично для [тарифов Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) сказано: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (каждое место в тарифе Team закреплено за одним конкретным человеком, и пользоваться им может только этот человек) и «Users can't route their work through another user's account to use that user's remaining usage» (пользователи не могут направлять свою работу через аккаунт другого пользователя, чтобы расходовать оставшийся лимит этого пользователя). Поэтому использование автоматического выбора между местами может нарушать условия Anthropic. Пользоваться ли им в этом случае, решать вам: прежде чем на него полагаться, уточните это у своей организации и у Anthropic. О личных тарифах (Pro, Max) см. [примечание об условиях в FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Уже используете ~/.claude?

Зарегистрируйте папку вместо повторного входа:

```sh
ccam import
```

Папка остаётся на месте, а обычный `claude` работает как прежде.
