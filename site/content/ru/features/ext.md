---
title: Просмотр расширений
description: ccam ext показывает навыки, команды, плагины, MCP-серверы и хуки, которые загружает каждый аккаунт, откуда берётся каждый из них и что сломано.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` отвечает на вопрос «почему этот навык (или MCP-сервер) работает в одном аккаунте и не работает в другом?». Для каждого аккаунта он перечисляет расширения, которые загружает Claude Code, и откуда берётся каждое из них: из группы общего доступа, из собственной папки аккаунта или из проекта, в котором вы находитесь. Под таблицей он перечисляет то, что мешает расширению загружаться так, как вы ожидаете. Он только читает.

## Использование

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

Ячейка называет источник этого элемента в этом аккаунте: группу общего доступа (`default`), `account`, если элемент лежит в собственной папке аккаунта, или `-`, если у аккаунта его нет. `!` отмечает элемент с проблемой, которая объясняется под таблицей.

## Что он читает

| Вид | Где | Показанный источник |
| --- | --- | --- |
| Навыки | каждая папка в `skills/` с `SKILL.md` | группа, на которую ссылается `skills/`, или `account` |
| Команды, агенты | файлы `.md` в `commands/` и `agents/`, подпапки как `folder/name` | группа или `account` |
| Плагины | `plugins/installed_plugins.json` | группа или `account` |
| MCP-серверы | `mcp.json` группы (если MCP общий) и `.claude.json` аккаунта | группа или `account` |
| Хуки | `hooks` в `settings.json`, с подсчётом по событиям | группа или `account` |
| Этот проект | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` и `.mcp.json` в текущей папке | выводится один раз, для всех аккаунтов |

## Какие проблемы он сообщает

- **Сломанные ссылки**: ссылка в `skills/`, `commands/` или `agents/`, цель которой пропала, или общая папка (`skills`, `commands`, `agents`, `plugins`), которая ссылается на несуществующее место (исправьте это командой `ccam doctor --fix`).
- **Папка навыка без `SKILL.md`**: Claude Code её пропускает.
- **MCP-сервер, по-разному определённый в двух местах**: в `mcp.json` группы и в `.claude.json` аккаунта или в одном из них и в `.mcp.json` проекта. Одинаковые настройки в обоих местах (их записывает `mcp_strategy = "sync"`) проблемой не считаются. Перенесите собственный сервер аккаунта в группу командой `ccam mcp import <account> <server>`.
- **`installed_plugins.json`, который не является корректным JSON.**

## Безопасность и особые случаи

- `ccam ext` ничего не меняет и не выполняет сетевых запросов. Он читает имена серверов и количество хуков, но никогда не читает значения `env` или заголовков в определениях MCP.
- Он показывает то, что лежит на диске. Плагин, который установлен, но выключен в настройках, или сервер, заблокированный управляемыми настройками, всё равно отображается.
- Запустите его в папке проекта, чтобы включить папку `.claude` и `.mcp.json` этого проекта.

## См. также

- [Общий доступ между аккаунтами](/features/share)
- [Общие MCP-серверы](/features/mcp)
- [Проверка и исправление](/features/doctor)
