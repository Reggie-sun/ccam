---
title: شرح CLAUDE_CONFIG_DIR
description: ما وظيفة CLAUDE_CONFIG_DIR في Claude Code، وما المخزَّن في المجلد، وأين يحفظ macOS تسجيل الدخول، وأي الملفات يمكن مشاركتها بين الحسابات.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: 2cf08478085f
---

يخبر `CLAUDE_CONFIG_DIR` برنامج Claude Code بالمجلد الذي يستخدمه بدلًا من `~/.claude`. كل ما يخص الحساب موجود هناك: تسجيل الدخول والإعدادات وخوادم MCP والإضافات وسجل المحادثات. وتوجيهه إلى مجلدات مختلفة يمنحك حسابات منفصلة.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## هل CLAUDE_CONFIG_DIR مدعوم رسميًا؟

نعم. يذكره Claude Code ضمن [متغيرات البيئة](https://code.claude.com/docs/en/env-vars) الخاصة به بوصفه طريقة لتغيير مجلد الإعدادات، "مفيدة لتشغيل عدة حسابات جنبًا إلى جنب"، ويستخدمه [دليل المصادقة](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) لتسجيل الدخول بعدة حسابات. يمكنك تعيينه في الصدفة أو في إعدادات المستخدم، لكن ليس في إعدادات مشروع.

## ما الموجود داخل المجلد؟

تم التحقق باستخدام Claude Code 2.1.309:

| المسار | المحتوى | قابل للمشاركة |
| --- | --- | --- |
| `.credentials.json` | رموز تسجيل الدخول (Linux وWindows)، ورموز OAuth الخاصة بـ MCP | لا |
| `.claude.json` | تفاصيل الحساب، وخوادم MCP الخاصة بالمستخدم، وحالة كل مشروع | لا (خوادم MCP بشكل منفصل) |
| `settings.json` | الأذونات والخطافات (hooks) والنموذج وشريط الحالة والإضافات | نعم |
| `CLAUDE.md` | التعليمات العامة | نعم |
| `agents/`، `commands/`، `skills/`، `plugins/` | الإضافات | نعم |
| `projects/` | المحادثات كملفات `.jsonl`، وذاكرة كل مشروع | نعم |
| `history.jsonl` | سجل الأوامر المكتوبة (prompts) | نعم |
| `sessions/`، `shell-snapshots/`، `cache/` وغيرها | حالة التشغيل | لا |

## ما الذي يوقع الناس في الخطأ؟

**دون `CLAUDE_CONFIG_DIR`، لا يكون `.claude.json` داخل `~/.claude`.** يحتفظ Claude Code به حينها في `~/.claude.json`، بجوار المجلد. إذا بدأت باستخدام `CLAUDE_CONFIG_DIR=~/.claude`، يبحث Claude Code عن `~/.claude/.claude.json` بدلًا من ذلك ولا يجد حسابك. يعالج `ccam import` ذلك بربط الملفين. وإذا كان في `~/.claude` بالفعل ملف `.claude.json` خاص به، من تشغيل سابق باستخدام `CLAUDE_CONFIG_DIR=~/.claude`، يترك ccam الملفين ويخبرك بذلك، لأنك وحدك تعرف أيهما الصحيح.

**على macOS، يكون تسجيل الدخول في سلسلة المفاتيح وليس في المجلد.** يخزّنه Claude Code تحت الاسم `Claude Code-credentials`، متبوعًا بـ `-` وأول 8 أرقام ست عشرية من تجزئة SHA-256 لمسار المجلد عندما يكون `CLAUDE_CONFIG_DIR` معيّنًا، ولحساب المستخدم الذي يحدده `$USER`. لذلك تحصل المجلدات المنفصلة على إدخالات منفصلة في سلسلة المفاتيح، والمجلد الذي يُنقل أو تُعاد تسميته يفقد تسجيل دخوله؛ ولهذا يحتفظ `ccam rename` باسم المجلد على macOS. وإذا رفضت سلسلة المفاتيح الكتابة، يعود Claude Code إلى استخدام `.credentials.json` في المجلد، كما [يصف توثيقه](https://code.claude.com/docs/en/authentication#credential-management).

**ليس كل تسجيل دخول محفوظًا في المجلد.** تسجيل الدخول إلى Claude Console دون مفتاح API يُخزَّن خارجه، لذلك لا تفصل المجلدات المنفصلة بين تسجيلَي دخول من هذا النوع. أما تسجيلات الدخول إلى claude.ai ومفاتيح API فتُحفظ لكل مجلد على حدة.

## أين يوجد المجلد على Windows؟

افتراضيًا في `%USERPROFILE%\.claude`، مع تسجيل الدخول في `.credentials.json` داخله، أو في Windows Credential Manager عندما يُضبط Claude Code لاستخدامه، تحت اسم مشتق من مسار المجلد كما على macOS. عيّن `CLAUDE_CONFIG_DIR` إلى مجلد آخر لكل حساب إضافي، بالطريقة نفسها كما على Linux وmacOS.

## هل يمكن للمجلدات مشاركة الملفات؟

الروابط الرمزية تعمل: يقرأ Claude Code ويكتب عبر ملف `settings.json` أو مجلد `projects/` مرتبط، ويحافظ على الرابط. خوادم MCP هي الاستثناء، لأنها تتشارك `.claude.json` مع تسجيل الدخول؛ مرّرها باستخدام `--mcp-config` بدلًا من ذلك.

يُبنى ccam على هذه القواعد تحديدًا. راجع [المشاركة](/docs/sharing) لمعرفة كيف يربط المجلدات، و[دليل الحسابات المتعددة](/guides/multiple-claude-code-accounts) للصورة الأشمل.
