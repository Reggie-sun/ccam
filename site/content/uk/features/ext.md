---
title: Перегляд розширень
description: ccam ext показує навички, команди, плагіни, сервери MCP і хуки, які завантажує кожен обліковий запис, звідки береться кожне з них і що зламано.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` відповідає на запитання «чому ця навичка (або сервер MCP) працює в одному обліковому записі й не працює в іншому?». Для кожного облікового запису він перелічує розширення, які завантажує Claude Code, і звідки береться кожне з них: із групи спільного використання, із власної теки облікового запису або з проєкту, у якому ви перебуваєте. Під таблицею він перелічує те, що заважає розширенню завантажуватися так, як ви очікуєте. Він лише читає.

## Використання

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Клітинка називає джерело цього елемента в цьому обліковому записі: групу спільного використання (`default`), `account`, коли елемент лежить у власній теці облікового запису, або `-`, коли в облікового запису його немає. `!` позначає елемент із проблемою, яку пояснено під таблицею.

## Що він читає

| Вид | Де | Показане джерело |
| --- | --- | --- |
| Навички | кожна тека в `skills/` з `SKILL.md` | група, на яку посилається `skills/`, або `account` |
| Команди, агенти | файли `.md` у `commands/` і `agents/`, підтеки як `folder/name` | група або `account` |
| Плагіни | `plugins/installed_plugins.json` | група або `account` |
| Сервери MCP | `mcp.json` групи (коли MCP спільний) і `.claude.json` облікового запису | група або `account` |
| Хуки | `hooks` у `settings.json`, з підрахунком за подіями | група або `account` |
| Цей проєкт | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` і `.mcp.json` у поточній теці | виводиться один раз, для всіх облікових записів |

## Про які проблеми він повідомляє

- **Зламані посилання**: посилання в `skills/`, `commands/` або `agents/`, ціль якого зникла, або спільна тека (`skills`, `commands`, `agents`, `plugins`), що посилається на відсутнє місце (виправте це командою `ccam doctor --fix`).
- **Тека навички без `SKILL.md`**: Claude Code її пропускає.
- **Сервер MCP, по-різному визначений у двох місцях**: у `mcp.json` групи та в `.claude.json` облікового запису або в одному з них і в `.mcp.json` проєкту. Однакові налаштування в обох місцях (їх записує `mcp_strategy = "sync"`) проблемою не вважаються. Перенесіть власний сервер облікового запису до групи командою `ccam mcp import <account> <server>`.
- **`installed_plugins.json`, який не є коректним JSON.**

## Безпека та граничні випадки

- `ccam ext` нічого не змінює й не надсилає мережевих запитів. Він читає назви серверів і кількість хуків, але ніколи не читає значень `env` чи заголовків у визначеннях MCP.
- Він показує те, що лежить на диску. Плагін, який установлено, але вимкнено в налаштуваннях, або сервер, заблокований керованими налаштуваннями, однаково відображається.
- Запустіть його в теці проєкту, щоб охопити теку `.claude` і `.mcp.json` цього проєкту.

## Пов'язані сторінки

- [Спільне використання між обліковими записами](/features/share)
- [Спільні сервери MCP](/features/mcp)
- [Перевірка та виправлення](/features/doctor)
