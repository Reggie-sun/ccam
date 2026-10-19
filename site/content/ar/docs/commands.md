---
title: مرجع الأوامر
description: كل أوامر ccam وخياراتها، من ccam add وccam <name> إلى ccam usage وccam next وccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

ضع `--dry-run` قبل أي أمر لترى ما سيفعله دون أن يفعله، مثلًا `ccam --dry-run work`.

<!-- commands -->

## تشغيل Claude Code

يشغّل `ccam <name|number> [arguments]` برنامج Claude Code بذلك الحساب. لا يقرأ ccam الوسائط التي تأتي بعد الاسم؛ بل يمرّرها إلى `claude` بالترتيب نفسه، ويتصرف الإدخال القياسي ورمز الخروج والإشارات تمامًا كما مع `claude` نفسه. وضع `--` بعد الاسم اختياري.

أما `ccam` مع وسائط ولكن دون اسم حساب، مثل `ccam -p "summarise"`، فيستخدم الحساب الافتراضي. عيّنه باستخدام `ccam use <name>`.

يشغّل الأمر المختصر مثل `ccwork` حسابه بالطريقة نفسها: يفعل `ccwork --resume` ما يفعله `ccam work --resume`. راجع [الأوامر المختصرة](/docs/commands#shortcuts).

## إضافة الحسابات وإزالتها

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

عند تشغيل `ccam import` دون تحديد مجلد، يبحث عن `~/.claude` و`~/.claude-*` ويسألك عن كل واحد منها. تُسجَّل المجلدات المستوردة في أماكنها؛ ولا ينقل ccam أي تسجيل دخول أو ينسخه أبدًا.

## الأوامر المختصرة

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

الأمر المختصر أمر تختار اسمه بنفسك، ويشغّل Claude Code بحساب واحد. ومثل `ccam work`، يمرّر كل وسيط إلى Claude Code دون تغيير، لذلك يعمل `ccwork --resume` و`ccwork -p "hi"` كما تتوقع. وهو رابط إلى ccam في `~/.local/bin`.

الأوامر المختصرة اختيارية. يسألك معالج `ccam add` هل تريد إنشاء واحد، وينشئ `ccam add <name> --shortcut <command>` واحدًا دون سؤال، ويعرض `ccam list` الأمر المختصر لكل حساب، أو `-` إذا لم يكن له أمر.

تستخدم الأسماء الأحرف a-z و0-9 و`.` و`-` و`_`. ولا يجوز أن يكون الاسم أمرًا من أوامر ccam، أو أمرًا مختصرًا لحساب آخر، أو اسم برنامج آخر موجود في `PATH`.

يحتفظ `ccam rename` باسم الأمر المختصر، ويحذف `ccam remove` الأمر المختصر مع الحساب. ويعيد `ccam shims install` إنشاء روابط كل الحسابات، مثلًا بعد نقل ccam.

## المشاركة

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

يقبل كل عنصر القيمة `on` (المجموعة `default`) أو `off` أو اسم مجموعة. راجع [المشاركة](/docs/sharing).

## نسخ إعداد موجود

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

ينسخ `ccam migrate` الملفات `settings.json` و`CLAUDE.md` والإضافات وخوادم MCP الخاصة بالمستخدم من مجلد Claude Code إلى مجموعة، فتحصل عليها كل الحسابات التي تشارك هذه العناصر. يُقرأ المجلد فقط. تنسخ الخيارات `--settings` و`--memory` و`--extensions` و`--mcp` و`--history` العناصر التي تحددها فقط؛ ودونها يُنسخ كل شيء ما عدا سجل المحادثات. يختار `-g` أو `--group` المجموعة، وينسخ `-y` أو `--yes` دون سؤال ويحتفظ بنسخة المجموعة عندما يختلف ملف. راجع [نسخ إعداد موجود](/docs/sharing#copy-an-existing-setup).

## خوادم MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## الاستخدام والتبديل

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## الصيانة

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
```

## الملفات

| المسار | المحتوى |
| --- | --- |
| `~/.config/ccam/config.toml` | الحسابات وما تشاركه |
| `~/.local/share/ccam/profiles/<name>/` | مجلد `CLAUDE_CONFIG_DIR` لكل حساب |
| `~/.local/share/ccam/groups/<group>/` | الإعدادات وخوادم MCP والإضافات والسجل المشتركة |
| `~/.local/state/ccam/` | ذاكرة التخزين المؤقت للاستخدام وأوقات آخر استخدام |
| `~/.local/bin/<shortcut>` | أوامر مختصرة مثل `ccwork`، وهي روابط إلى ccam |

تغيّر المتغيرات `XDG_CONFIG_HOME` و`XDG_DATA_HOME` و`XDG_STATE_HOME` مواقع هذه المجلدات.
