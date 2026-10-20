---
title: مرجع الأوامر
metaTitle: مرجع أوامر ccam
description: كل أوامر ccam وخياراتها، من ccam add وccam <name> إلى ccam usage وccam next وccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 5d3275878535
---

ضع `--dry-run` قبل أمر لترى ما سيفعله دون أن يفعله، مثلًا `ccam --dry-run work`. يعمل ذلك عند تشغيل Claude Code (`ccam <name>` و`ccam run` و`ccam next`) ومع `share` و`migrate` و`shortcut` و`update` و`doctor --fix`. أما الأوامر الأخرى التي تغيّر شيئًا فترفضه بدلًا من أن تعمل.

<!-- commands -->

## تشغيل Claude Code

يشغّل `ccam <name|number> [arguments]` برنامج Claude Code بذلك الحساب. لا يقرأ ccam الوسائط التي تأتي بعد الاسم؛ بل يمرّرها إلى `claude` بالترتيب نفسه، ويتصرف الإدخال القياسي ورمز الخروج والإشارات تمامًا كما مع `claude` نفسه. وضع `--` بعد الاسم اختياري.

أما `ccam` مع وسائط ولكن دون اسم حساب، مثل `ccam -p "summarise"`، فيستخدم الحساب المعيّن للمجلد الحالي، وإلا فالحساب الافتراضي. عيّنهما باستخدام `ccam folder <name> [folder]` و`ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

الخيار `--fallback` قبل اسم الحساب مخصص لتشغيلات `-p`/`--print`. عندما يبلغ الحساب حد استخدامه، يشغّل ccam الأمر نفسه مجددًا بالحساب المسجّل الدخول الذي لديه أكبر هامش، ولا يطبع إلا مخرج التشغيل الذي اكتمل. راجع [الحسابات حسب المجلد](/features/folder) و[إعادة التشغيل عند بلوغ الحد](/features/next#run-again-at-the-limit).

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

ترفض الأوامر `ccam rename` و`ccam remove` و`ccam share` تغيير حساب ما دام Claude Code يعمل بذلك الحساب. أنهِ الجلسة أولًا.

يطبع `ccam list --json` كل حساب مع رقم خانته واختصاره وبريده الإلكتروني ومجلده. ويربط الكائن `share` كل عنصر (`mcp` و`settings` و`memory` و`extensions` و`history`) باسم مجموعة أو بالقيمة `local`.

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

ينسخ `ccam migrate` الملفات `settings.json` و`CLAUDE.md` والإضافات وخوادم MCP الخاصة بالمستخدم من مجلد Claude Code إلى مجموعة، فتحصل عليها كل الحسابات التي تشارك هذه العناصر. يُقرأ المجلد فقط. تنسخ الخيارات `--settings` و`--memory` و`--extensions` و`--mcp` و`--history` العناصر التي تحددها فقط. ودونها تتيح لك قائمة اختيار في الطرفية أن تختار، مع تحديد كل شيء ما عدا سجل المحادثات؛ وفي غير ذلك يُنسخ كل شيء ما عدا سجل المحادثات. يختار `-g` أو `--group` المجموعة، وينسخ `-y` أو `--yes` دون سؤال ويحتفظ بنسخة المجموعة عندما يختلف ملف. كما ينسخ تسجيلات الدخول المحفوظة للخوادم إلى ملف تسجيل الدخول الخاص بكل حساب. راجع [نسخ إعداد موجود](/docs/sharing#copy-an-existing-setup).

## خوادم MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

ينقل `ccam mcp import <account>` الخوادم التي لا يملكها إلا ذلك الحساب إلى مجموعته. ويُتخطّى الخادم الذي تستخدم المجموعة اسمه بالفعل لتعريف مختلف؛ أضف `--force` لاستبدال نسخة المجموعة، ويحتفظ ccam بنسخة احتياطية منها أولًا.

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
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

يتحقق `ccam update` من مجموع SHA-256 الخاص بالإصدار قبل أن يستبدل ccam، ويثبّت `--version v0.2.0` إصدارًا محددًا. أما ccam المثبّت من حزمة `.deb` فيُحدَّث باستخدام apt. راجع [التحديث](/docs/install#update).

## الملفات

| المسار | المحتوى |
| --- | --- |
| `~/.config/ccam/config.toml` | الحسابات وما تشاركه |
| `~/.local/share/ccam/profiles/<name>/` | مجلد `CLAUDE_CONFIG_DIR` لكل حساب |
| `~/.local/share/ccam/groups/<group>/` | الإعدادات وخوادم MCP والإضافات والسجل المشتركة |
| `~/.local/state/ccam/` | ذاكرة التخزين المؤقت للاستخدام وأوقات آخر استخدام |
| `~/.local/bin/<shortcut>` | أوامر مختصرة مثل `ccwork`، وهي روابط إلى ccam |

تغيّر المتغيرات `XDG_CONFIG_HOME` و`XDG_DATA_HOME` و`XDG_STATE_HOME` مواقع هذه المجلدات. وعلى Windows، يكون `~` هو `%USERPROFILE%` وتكون الاختصارات `<shortcut>.exe`.
