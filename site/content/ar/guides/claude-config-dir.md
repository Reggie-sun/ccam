---
title: شرح CLAUDE_CONFIG_DIR
description: ما وظيفة CLAUDE_CONFIG_DIR في Claude Code، وما المخزَّن في المجلد، وأين يحفظ macOS تسجيل الدخول، وأي الملفات يمكن مشاركتها بين الحسابات.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

يخبر `CLAUDE_CONFIG_DIR` برنامج Claude Code بالمجلد الذي يستخدمه بدلًا من `~/.claude`. كل ما يخص الحساب موجود هناك: تسجيل الدخول والإعدادات وخوادم MCP والإضافات وسجل المحادثات. وتوجيهه إلى مجلدات مختلفة يمنحك حسابات منفصلة.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## ما الموجود داخل المجلد

تم التحقق باستخدام Claude Code 2.1.302:

| المسار | المحتوى | قابل للمشاركة |
| --- | --- | --- |
| `.credentials.json` | رموز تسجيل الدخول (Linux)، ورموز OAuth الخاصة بـ MCP | لا |
| `.claude.json` | تفاصيل الحساب، وخوادم MCP الخاصة بالمستخدم، وحالة كل مشروع | لا (خوادم MCP بشكل منفصل) |
| `settings.json` | الأذونات والخطافات (hooks) والنموذج وسطر الحالة والإضافات | نعم |
| `CLAUDE.md` | التعليمات العامة | نعم |
| `agents/`، `commands/`، `skills/`، `plugins/` | الإضافات | نعم |
| `projects/` | المحادثات كملفات `.jsonl`، وذاكرة كل مشروع | نعم |
| `history.jsonl` | سجل الأوامر المكتوبة (prompts) | نعم |
| `sessions/`، `shell-snapshots/`، `cache/` وغيرها | حالة التشغيل | لا |

## تفصيلتان توقعان الناس في الخطأ

**دون `CLAUDE_CONFIG_DIR`، لا يكون `.claude.json` داخل `~/.claude`.** يحتفظ Claude Code به حينها في `~/.claude.json`، بجوار المجلد. إذا بدأت باستخدام `CLAUDE_CONFIG_DIR=~/.claude`، يبحث Claude Code عن `~/.claude/.claude.json` بدلًا من ذلك ولا يجد حسابك. يعالج `ccam import` ذلك بربط الملفين.

**على macOS، يكون تسجيل الدخول في سلسلة المفاتيح وليس في المجلد.** يخزّنه Claude Code تحت الاسم `Claude Code-credentials`، متبوعًا بتجزئة (hash) قصيرة لمسار المجلد عندما يكون `CLAUDE_CONFIG_DIR` معيّنًا. لذلك تحصل المجلدات المنفصلة على إدخالات منفصلة في سلسلة المفاتيح. وإذا رفضت سلسلة المفاتيح الكتابة، يعود Claude Code إلى استخدام `.credentials.json` في المجلد.

## المشاركة بين المجلدات

الروابط الرمزية تعمل: يقرأ Claude Code ويكتب عبر ملف `settings.json` أو مجلد `projects/` مرتبط، ويحافظ على الرابط. خوادم MCP هي الاستثناء، لأنها تتشارك `.claude.json` مع تسجيل الدخول؛ مرّرها باستخدام `--mcp-config` بدلًا من ذلك.

يُبنى ccam على هذه القواعد تحديدًا. راجع [المشاركة](/docs/sharing) لمعرفة كيف يربط المجلدات، و[دليل الحسابات المتعددة](/guides/multiple-claude-code-accounts) للصورة الأشمل.
