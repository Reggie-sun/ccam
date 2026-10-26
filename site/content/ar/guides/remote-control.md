---
title: استخدام Remote Control مع ccam
description: شغّل Remote Control في Claude Code بأي حساب باستخدام ccam <name> rc، وتعرّف على الخيارات المتاحة، ولماذا يعرض ccam نسخ خوادم MCP المشتركة، وحلّ المشكلات.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 4d46f49c5dea
---

يشغّل `ccam work rc` ميزة Remote Control في Claude Code بالحساب `work`، فتتحكم في الجلسات على هذا الحاسوب من claude.ai/code أو من تطبيق Claude على الهاتف. وهو يعمل مثل `claude rc` مع فرق واحد: إذا كان الحساب يشارك خوادم MCP، يسألك ccam مرة واحدة هل ينسخها إلى ملف `.claude.json` الخاص بالحساب، لأن هذا الملف هو المكان الوحيد الذي تقرأ منه جلسات Remote Control الخوادم.

## التشغيل

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

يعمل [الأمر المختصر](/features/shortcuts) بالطريقة نفسها: `ccwork rc` هو `ccam work rc`. شغّل الأمر في المجلد الذي تريد العمل فيه، ثم افتح claude.ai/code أو تطبيق Claude بعد تسجيل الدخول بحساب Claude نفسه الذي يستخدمه `work`.

## الخيارات

يلتقي هنا نوعان من الخيارات: خيارات ccam نفسه، التي توضع قبل اسم الحساب أو تكون في إعدادات ccam، وخيارات Remote Control، التي توضع بعد `rc` وتصل إلى Claude Code دون تغيير.

### خيارات ccam وإعداداته

| الخيار أو الإعداد | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--dry-run` (قبل الحساب) | يطبع البيئة والأمر، وإشعار MCP إن وُجد. لا يسأل ولا ينسخ | معطّل |
| `mcp_strategy = "sync"` في أعلى `~/.config/ccam/config.toml` | ينسخ خوادم المجموعة إلى `.claude.json` قبل كل تشغيل، سواء أكان Remote Control أم لا، فلا يُطرح أي سؤال أبدًا | `"flag"` |
| `ccam group leave <group> work` | يُخرج الحساب من المجموعة التي تشاركه MCP. وتُنسخ خوادم المجموعة إلى ملف `.claude.json` الخاص به بوصفها خوادمه الخاصة، فتحصل عليها Remote Control ولا يُطرح أي سؤال؛ ومع `--fresh` تبدأ Remote Control بخوادم الحساب نفسه فقط | في `default` |

عندما يشارك الحساب خوادم MCP ولا يحتوي ملف `.claude.json` الخاص به على القائمة الحالية للمجموعة، أو عندما يختفي منه خادم منسوخ، يسأل ccam `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`. `y` ينسخها؛ و`N` أو Enter يبدأ بدونها.

### خيارات Remote Control

كل ما يأتي بعد `rc` يخص Claude Code. هذه هي خيارات `claude remote-control` في Claude Code 2.1.304؛ ويطبع `ccam work rc --help` القائمة الخاصة بإصدارك.

| الخيار | ما يفعله |
| --- | --- |
| `--name <name>` | اسم الجلسة كما يظهر في claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | بادئة لأسماء الجلسات المولَّدة (الافتراضي: اسم المضيف) |
| `-c`، `--continue` | يعيد الاتصال بالجلسة التي سجّلتها Remote Control أخيرًا لهذا المجلد، خلال نحو 4 ساعات |
| `--session-id <id>` | يعيد الاتصال بجلسة محددة |
| `--permission-mode <mode>` | وضع الأذونات للجلسات الجديدة: `acceptEdits` أو `auto` أو `bypassPermissions` أو `default` أو `dontAsk` أو `plan` |
| `--spawn <mode>` | `same-dir` (الافتراضي)، أو `worktree` (شجرة عمل git لكل جلسة)، أو `session` (جلسة واحدة، ويخرج عند انتهائها) |
| `--capacity <N>` | أقصى عدد من الجلسات في وقت واحد في الوضعين `same-dir` و`worktree` (الافتراضي 32) |
| `--[no-]create-session-in-dir` | ينشئ جلسة أولى في المجلد الحالي عند البدء (مفعّل افتراضيًا) |
| `--[no-]chrome` | Claude in Chrome للجلسات الجديدة (معطّل افتراضيًا) |
| `-d`، `--debug`، `--debug-file <path>`، `-v`، `--verbose` | مخرجات التصحيح والمخرجات المفصّلة |

يقرأ ccam أمر Claude Code كما يقرؤه Claude Code، متخطيًا الخيارات مع قيمها، لذلك يشغّل كلٌّ من `ccam work rc --name laptop` و`ccam work --model opus rc` ميزة Remote Control. والخيار الذي لا يعرفه ccam يُعدّ خيارًا بلا قيمة؛ ولمزيد من الأمان، ضع `rc` مباشرة بعد اسم الحساب.

## ما يفعله ccam

يسلك التشغيل المسار المعتاد ([تشغيل Claude Code](/features/launch)) مع تغيير واحد في خطوة MCP:

1. يجد ccam الحساب، ويفحص روابط المشاركة الخاصة به، ويحدد `claude`، كما في أي تشغيل.
2. يرى أن الأمر هو `rc` أو `remote-control` فلا يضيف `--mcp-config`.
3. إذا كان الحساب يشارك خوادم MCP، يقارن ملف `mcp.json` الخاص بالمجموعة مع `sync-stamp` في `~/.local/state/ccam/mcp/work/`، ويتحقق من أن كل خادم من المجموعة لا يزال في `.claude.json`. وعندما يحتوي `.claude.json` بالفعل على الخوادم الحالية، أو لا تكون لدى المجموعة أي خوادم، ينتقل مباشرة إلى الخطوة 6.
4. وإلا يسمّي الخوادم التي ستفتقدها الجلسات:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. عند `y` يحفظ نسخة احتياطية من `.claude.json`، ويأخذ قفل Claude Code نفسه، ويضيف خوادم المجموعة إلى `mcpServers` مع الإبقاء على خوادم الحساب الأخرى؛ وإذا كان لدى الحساب خادمه الخاص بالاسم نفسه الذي يحمله خادم في المجموعة، يحتفظ ccam بخادم الحساب ويخبرك بذلك. وعند `N`، أو دون طرفية، أو بينما الحساب يعمل، لا يغيّر شيئًا ويخبرك بذلك:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. يسجّل الجلسة، ويعيّن `CLAUDE_CONFIG_DIR` إلى مجلد الحساب، ويستبدل نفسه بـ `claude rc` عبر `exec`.

بعد ذلك تظل Remote Control تعمل بوصفها خادمًا. وكل جلسة تبدؤها لك هي عملية فرعية من تلك العملية، فتعمل بالحساب `work`: بتسجيل دخول `work` وإعداداته وسجل محادثاته، وتُحتسب من حد استخدام `work`. ويعدّ ccam الحساب قيد التشغيل ما دام خادم Remote Control شغّله ccam يعمل. أما الخادم الذي بدأ باستخدام `claude remote-control` العادي فلا يُحتسب إلا ما دامت لديه جلسات؛ ولا يستطيع ccam رؤية خادم لا جلسة له بعد.

يُظهر التشغيل التجريبي الفرق. التشغيل العادي يحمل خوادم المجموعة بوصفها وسيطًا، أما Remote Control فلا:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## لماذا يعمل بهذه الطريقة

**لا تأخذ Remote Control الخيار `--mcp-config`.** في عمليات التشغيل العادية يسلّم ccam خوادم المجموعة إلى Claude Code بالشكل `--mcp-config=<path>` ولا يكتب أبدًا في ملفات Claude Code. يرفض `claude remote-control` هذا الخيار قبل اسم أمره ولا يقبل أيًّا منه بعده، لذلك كان `ccam work rc` يفشل قبل ccam 0.6.0 مع كل حساب يشارك خوادم MCP. والآن يحذف ccam هذا الخيار.

**لا تقرأ جلساتها الخوادم إلا من `.claude.json`.** الجلسات التي تبدؤها Remote Control لا ترث الخيارات من سطر الأوامر. والمكان الوحيد الذي تجد فيه خوادم MCP الخاصة بالمستخدم هو ملف `.claude.json` الخاص بالحساب. ويحتوي هذا الملف أيضًا على تفاصيل تسجيل دخول الحساب وحالة كل مجلد، لذلك لا يستطيع ccam ربطه بنسخة مشتركة كما يفعل مع `settings.json`. ونسخ الخوادم إليه هو الطريقة الوحيدة لإعطاء الجلسات القائمة المشتركة.

**يسأل ccam قبل الكتابة.** مع استراتيجية الخيار الافتراضية، يترك ccam ملف `.claude.json` دون تغيير عند التشغيل. وتشغيل Remote Control هو التشغيل الوحيد الذي يضطر فيه إلى الكتابة هناك، لذلك يسأل أولًا، ويحتفظ بنسخة احتياطية، ولا يمس شيئًا سوى `mcpServers`. (وخارج التشغيل، تكتب هناك أيضًا مغادرةُ مشاركة MCP أو الانضمامُ إليها و`ccam mcp import`، كما تذكر خططها.)

**لا يكتب بينما الحساب يعمل.** يعيد Claude Code كتابة `.claude.json` في أثناء الجلسة، والتغيير الذي يأتي من الخارج سيضيع أو يُكتب فوقه. لهذا يتخطى ccam النسخ عندما تعمل أي جلسة للحساب، بما في ذلك خادم Remote Control آخر وClaude Code الذي بدأ دون ccam.

**يسأل مرة واحدة لكل تغيير.** بعد النسخ، يحتفظ `sync-stamp` بوقت تعديل ملف `mcp.json` الخاص بالمجموعة. ولا يسأل ccam مرة أخرى إلا عندما يتغير ذلك الملف، مثلًا بعد `ccam mcp add`، أو عندما يختفي خادم منسوخ من `.claude.json`. وإذا أجبت `N`، فلا توجد نسخة بعد، لذلك يسأل عند كل تشغيل لـ Remote Control.

**يُبقي `ccam mcp` النسخ محدّثة.** يتقدم المدخل الموجود في `.claude.json` على `--mcp-config` عند تعارض الأسماء، فالنسخة القديمة ستخفي تعريفًا أحدث في عمليات التشغيل العادية أيضًا. لذلك تحدّث الأوامر `ccam mcp add` و`add-json` و`remove` و`import`، وكذلك `ccam migrate`، النسخ في كل حساب لا يعمل، وتخبرك بذلك عندما تتغيّر نسخة:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**دون طرفية يكون الجواب لا.** كثيرًا ما تعمل Remote Control من سكربت أو خدمة أو `nohup`، حيث لا يوجد من يجيب. فيبدأ ccam حينها دون الخوادم بدلًا من الانتظار.

## المشكلات وحلولها

### خوادم MCP المشتركة مفقودة في جلسات Remote Control

كانت الإجابة عن السؤال `N`، أو لم تكن هناك طرفية، أو كان الحساب يعمل. أوقف Remote Control، وأغلق جلسات الحساب الأخرى، وشغّل `ccam work rc` مرة أخرى في طرفية؛ وأجب `y`. ولكي تتوقف عن تلقّي السؤال، عيّن `mcp_strategy = "sync"` في `~/.config/ccam/config.toml`. والخوادم الموجودة في جلسة لا تتغير ما دامت تعمل، لذلك ابدأ جلسة جديدة بعد النسخ.

### "Account work is running, so its .claude.json is left alone"

ما زالت جلسة لـ `work` مفتوحة: `ccam work` عادي، أو أمر مختصر، أو خادم Remote Control آخر. اخرج منها وشغّل `ccam work rc` مرة أخرى. ولا تُحتسب الجلسات التي يشغّلها حساب آخر.

### يسأل ccam في كل مرة

لم تنسخ الخوادم بعد. أجب `y` مرة واحدة، أو عيّن `mcp_strategy = "sync"`. وإذا كنت لا تريد الخوادم في Remote Control أبدًا، فأوقف مشاركة MCP مع الحساب دون الاحتفاظ بها: `ccam group leave default work --fresh`، أو ضع MCP في مجموعة خاصة به لا ينتمي إليها الحساب. أما `ccam group leave` العادي فينسخ خوادم المجموعة إلى ملف `.claude.json` الخاص بالحساب، فتحصل عليها Remote Control.

### يفشل `ccam work rc` فورًا مع خطأ يتعلق بـ `--mcp-config`

إصدار ccam لديك أقدم من 0.6.0 ويمرّر `--mcp-config` إلى Remote Control. شغّل `ccam update`. أما مع إصدار حديث من ccam، فالخطأ نفسه يعني أن خيارًا غير معروف له قيمة جاء قبل `rc`، فظنّ ccam خطأً أن تلك القيمة هي الأمر؛ ضع `rc` مباشرة بعد اسم الحساب.

### التغيير الذي أُجري باستخدام `ccam mcp` لا يصل إلى Remote Control

لا يحدّث `ccam mcp add` و`remove` النسخ إلا في الحسابات التي لا تعمل، ويخبرانك بذلك: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` أوقف Remote Control وشغّلها مرة أخرى؛ فيرى ccam أن `mcp.json` تغيّر ويسأل مجددًا.

### يُستخدم خادمي الخاص بدلًا من خادم المجموعة

إذا كان للحساب خادم خاص يحمل اسم خادم من المجموعة نفسه، فإن النسخ يتركه في مكانه ويخبرك بذلك: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`. ولاستخدام خادم المجموعة، أزل خادم الحساب (`claude mcp remove -s user <name>`)؛ ولإعطاء نسختك لكل الحسابات، انقلها إلى المجموعة باستخدام `ccam mcp import work <server> --force`. وتوجد النسخ الاحتياطية من `.claude.json` في `~/.local/state/ccam/mcp/work/backups/`، الذي يحتفظ بأحدث خمس نسخ.

### أريد إزالة الخوادم المنسوخة

يزيل `ccam mcp remove <server>` الخادم من المجموعة ومن كل نسخة. وللعودة إلى الحالة التي سبقت النسخ، استعد أحدث نسخة احتياطية من `~/.local/state/ccam/mcp/work/backups/` بينما لا تعمل أي جلسة للحساب.

### تسأل Remote Control هل تثق بالمجلد، أو تخرج دون طرفية

تُخزَّن الثقة بالمجلدات لكل حساب في `.claude.json`، فالمجلد الذي وثقت به بالحساب `personal` جديد على `work`. شغّل `ccam work` مرة واحدة في ذلك المجلد ووافق، ثم ابدأ Remote Control هناك.

### "You must be logged in with a Claude account that has a subscription"

تحتاج Remote Control إلى تسجيل دخول باشتراك Claude. تحقق من البريد الإلكتروني للحساب في `ccam list`، وسجّل الدخول مرة أخرى باستخدام `ccam login work` إن لزم الأمر.

### لا تظهر الجلسة على الهاتف

تطبيق Claude أو claude.ai/code مسجَّل الدخول بحساب Claude مختلف. قارنه بالبريد الإلكتروني الذي يعرضه `ccam list` للحساب `work`.

### يحذّر ccam من أن جلسة Claude Code أخرى تشارك سجل المحادثات في هذا المجلد

يعمل حساب آخر من مجموعة السجل نفسها في المجلد نفسه. والجلسات التي تكمل المحادثة نفسها من الحسابين قد تخلط السجل. أجب `y` فقط إذا كانا سيعملان على محادثات مختلفة، أو ابدأ Remote Control في مجلد آخر.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [خوادم MCP المشتركة](/features/mcp)
- [المشاركة عبر المجموعات](/features/share)
- [إدارة الحسابات](/features/accounts)
