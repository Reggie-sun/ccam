---
title: فحص الإضافات
description: يعرض ccam ext المهارات والأوامر والمكوّنات الإضافية وخوادم MCP والخطافات التي يحمّلها كل حساب، ومصدر كل منها، وما هو معطوب.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

يجيب `ccam ext` عن السؤال "لماذا تعمل هذه المهارة (أو خادم MCP هذا) في حساب ولا تعمل في الآخر؟". فهو يسرد لكل حساب الإضافات التي يحمّلها Claude Code ومصدر كل منها: مجموعة مشاركة، أو مجلد الحساب نفسه، أو المشروع الذي أنت فيه. وتحت الجدول يسرد ما يمنع إضافة من التحميل كما تتوقع. وهو يقرأ فقط.

## الاستخدام

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

تذكر كل خلية مصدر ذلك العنصر في ذلك الحساب: مجموعة المشاركة (`default`)، أو `account` عندما يكون في مجلد الحساب نفسه، أو `-` عندما لا يملكه الحساب. وتشير `!` إلى عنصر فيه مشكلة، مشروحة تحت الجدول.

## ما يقرؤه

| النوع | المكان | المصدر المعروض |
| --- | --- | --- |
| المهارات | كل مجلد في `skills/` فيه `SKILL.md` | المجموعة التي يرتبط بها `skills/`، أو `account` |
| الأوامر والوكلاء | ملفات `.md` تحت `commands/` و`agents/`، والمجلدات الفرعية بالشكل `folder/name` | المجموعة، أو `account` |
| المكوّنات الإضافية | `plugins/installed_plugins.json` | المجموعة، أو `account` |
| خوادم MCP | ملف `mcp.json` الخاص بالمجموعة (عندما تكون MCP مشتركة) و`.claude.json` الخاص بالحساب | المجموعة، أو `account` |
| الخطافات | `hooks` في `settings.json`، معدودة لكل حدث | المجموعة، أو `account` |
| هذا المشروع | `.claude/skills` و`.claude/commands` و`.claude/agents` و`.claude/settings.json` و`.claude/settings.local.json` و`.mcp.json` في المجلد الحالي | تُسرد مرة واحدة، لكل الحسابات |

## المشكلات التي يبلّغ عنها

- **الروابط المعطوبة**: رابط في `skills/` أو `commands/` أو `agents/` اختفى هدفه، أو مجلد مشترك (`skills` أو `commands` أو `agents` أو `plugins`) يرتبط بمكان مفقود (أصلحه باستخدام `ccam doctor --fix`).
- **مجلد مهارة دون `SKILL.md`**: يتخطاه Claude Code.
- **خادم MCP معرَّف بشكل مختلف في مكانين**: ملف `mcp.json` الخاص بالمجموعة و`.claude.json` الخاص بالحساب، أو أحدهما و`.mcp.json` الخاص بالمشروع. أما الإعدادات نفسها في المكانين (وهو ما يكتبه `mcp_strategy = "sync"`) فليست مشكلة. انقل خادم الحساب الخاص إلى المجموعة باستخدام `ccam mcp import <account> <server>`.
- **ملف `installed_plugins.json` ليس JSON صالحًا.**

## الأمان والحالات الخاصة

- لا يغيّر `ccam ext` شيئًا ولا يرسل أي طلبات شبكة. فهو يقرأ أسماء الخوادم وأعداد الخطافات، ولا يقرأ أبدًا قيم `env` أو الترويسات في تعريفات MCP.
- يعرض ما هو موجود على القرص. فالمكوّن الإضافي المثبّت لكنه معطّل في الإعدادات، أو الخادم الذي تحظره الإعدادات المُدارة، يظل يظهر.
- شغّله في مجلد المشروع ليشمل مجلد `.claude` و`.mcp.json` الخاصين بذلك المشروع.

## صفحات ذات صلة

- [المشاركة بين الحسابات](/features/share)
- [خوادم MCP المشتركة](/features/mcp)
- [الفحص والإصلاح](/features/doctor)
