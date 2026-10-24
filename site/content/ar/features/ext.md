---
title: فحص الإضافات
description: يعرض ccam ext المهارات والأوامر والمكوّنات الإضافية وخوادم MCP والخطافات التي يحمّلها كل حساب، ومصدر كل منها، وما هو معطوب.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

تذكر كل خلية مصدر ذلك العنصر في ذلك الحساب: مجموعة المشاركة (`default`)، أو `account` عندما يكون في مجلد الحساب نفسه، أو `local` لخادم MCP يملكه الحساب لهذا المجلد فقط، أو `-` عندما لا يملكه الحساب. والمكوّن الإضافي المثبّت لكنه غير مشغَّل يُظهر `(off)` بعد مصدره. وتشير `!` إلى عنصر فيه مشكلة، مشروحة تحت الجدول؛ والمشكلات التي لا تخص أي حساب، مثل الخوادم التي لا يملكها إلا `claude` العادي، تأتي في النهاية دون اسم. ويطبع `ccam ext --json` التقرير نفسه (`accounts`، والمشروع، و`problems`)، مع `remote` على خوادم MCP التي يصل إليها Claude Code عبر الشبكة؛ ودون حسابات يطبع التقرير مع `accounts: []`.

## ما يقرؤه

| النوع | المكان | المصدر المعروض |
| --- | --- | --- |
| المهارات | كل مجلد في `skills/` فيه `SKILL.md` | المجموعة التي يرتبط بها `skills/`، أو `account` |
| الأوامر والوكلاء | ملفات `.md` تحت `commands/` و`agents/`، والمجلدات الفرعية بالشكل `folder/name` | المجموعة، أو `account` |
| المكوّنات الإضافية | `plugins/installed_plugins.json`، ويشغّلها `enabledPlugins` أو يوقفها في إعدادات الحساب والمشروع؛ والمكوّن الإضافي المثبّت لمشروع واحد فقط يظهر في ذلك المشروع | المجموعة، أو `account`؛ و`(off)` عندما لا يكون مشغَّلًا |
| خوادم MCP | ملف `mcp.json` الخاص بالمجموعة (عندما تكون MCP مشتركة)، و`.claude.json` الخاص بالحساب (خوادمه الخاصة، والخوادم التي يملكها لهذا المجلد فقط)، وخوادم مكوّناته الإضافية، بالشكل `plugin:<plugin>:<server>` | المجموعة، أو `account` أو `local` |
| الخطافات | `hooks` في `settings.json`، معدودة لكل حدث | المجموعة، أو `account` |
| هذا المشروع | `.claude/skills` و`.claude/commands` و`.claude/agents` و`.claude/settings.json` و`.claude/settings.local.json` و`.mcp.json` في المجلد الحالي | تُسرد مرة واحدة؛ والحساب الذي لم يثق بالمجلد يحصل على مشكلة لخوادم `.mcp.json` |

## المشكلات التي يبلّغ عنها

- **الروابط المعطوبة**: رابط في `skills/` أو `commands/` أو `agents/` اختفى هدفه، أو مجلد مشترك (`skills` أو `commands` أو `agents` أو `plugins`) يرتبط بمكان مفقود (أصلحه باستخدام `ccam doctor --fix`، الذي ينشئ مجلد المجموعة المفقود من جديد).
- **مجلد مهارة دون `SKILL.md`**: يتخطاه Claude Code.
- **خادم MCP معرَّف بشكل مختلف في مكانين**: ملف `mcp.json` الخاص بالمجموعة و`.claude.json` الخاص بالحساب، أو أحدهما و`.mcp.json` الخاص بالمشروع. أما الإعدادات نفسها في المكانين (وهو ما يكتبه `mcp_strategy = "sync"`) فليست مشكلة. انقل خادم الحساب الخاص إلى المجموعة باستخدام `ccam mcp import <account> <server>`.
- **ملف `installed_plugins.json` ليس JSON صالحًا**، ومكوّن إضافي لا يذكره `enabledPlugins`، ومكوّن إضافي مجلده مفقود.
- **خوادم لا يحمّلها إلا حساب واحد**: خادم موجود فقط في `.claude.json` الخاص بهذا الحساب، ولا تحمّله الحسابات الأخرى في مجموعته، وخادم في النطاق المحلي للحساب لهذا المجلد (وهو ما يفعله `claude mcp add` ما لم يُعطَ `-s user`). ويأتي كل منها مع أمر `ccam mcp import` الذي يشاركه. والخوادم الموجودة في `~/.claude.json` الخاص بـ `claude` العادي التي لا يحمّلها أي حساب تأتي مع أمر `ccam migrate --mcp`.
- **خادم مشروع ينتظر الموافقة**: لم يثق الحساب بالمجلد بعد، فيسأل Claude Code مرة واحدة، وحتى ذلك الحين تنتظر خوادم `.mcp.json` الخاصة بالمشروع. وتذكر المشكلة أين توافق عليها: `"enabledMcpjsonServers"` في `settings.json` الخاص بحسابات مجموعة الإعدادات، أو بهذا الحساب. والثقة الممنوحة لشجرة العمل الرئيسية لمستودع git تُحتسب في أشجار العمل الأخرى.
- **تسجيلات الدخول والموصلات**: خادم MCP يحتاج إلى تسجيل دخول في هذا الحساب بينما سجّل حساب آخر الدخول إليه (تسجيل الدخول يسري لحساب واحد فقط؛ سجّل الدخول باستخدام `/mcp`)، وموصلات claude.ai التي استخدمها حساب آخر ولم يستخدمها هذا الحساب.

## الأمان والحالات الخاصة

- لا يغيّر `ccam ext` شيئًا ولا يرسل أي طلبات شبكة. فهو يقرأ أسماء الخوادم وأعداد الخطافات، ولا يقرأ أبدًا قيم `env` أو الترويسات في تعريفات MCP. ولفحص تسجيلات الدخول يقرأ أسماء تسجيلات دخول MCP في تسجيل دخول كل حساب (على macOS من سلسلة المفاتيح)، لا الرموز أبدًا، وفقط للحسابات التي تحمّل خادمًا يحمّله حساب آخر أيضًا.
- يعرض ما هو موجود على القرص. فالخادم الذي تحظره الإعدادات المُدارة يظل يظهر؛ والمكوّن الإضافي المعطّل في الإعدادات يُظهر `(off)`.
- شغّله في مجلد المشروع ليشمل مجلد `.claude` و`.mcp.json` الخاصين بذلك المشروع.

## صفحات ذات صلة

- [المشاركة بين الحسابات](/features/share)
- [خوادم MCP المشتركة](/features/mcp)
- [الفحص والإصلاح](/features/doctor)
