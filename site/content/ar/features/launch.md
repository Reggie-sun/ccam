---
title: تشغيل Claude Code
description: شغّل Claude Code بأي حساب باستخدام ccam <name|number> أو ccam run أو أداة الاختيار في ccam؛ تُمرَّر الوسائط إلى claude دون تغيير، ويعرض --dry-run الأمر.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

يشغّل `ccam work` برنامج Claude Code بالحساب `work`، وكل ما يأتي بعد الاسم يُمرَّر إلى `claude` كما هو. استخدمه في أي مكان كنت ستكتب فيه `claude`: في الطرفية، أو في سكربت، أو مع إدخال يصل عبر أنبوب.

## الاستخدام

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

حدّد اسم الحساب أو رقمه من `ccam list`. لا يقرأ ccam الوسائط التي تأتي بعد الاسم، لذلك يعرض `ccam work --help` مساعدة Claude Code. ووضع `--` بعد الاسم اختياري. ويتصرف الإدخال القياسي ورمز الخروج والإشارات كما مع `claude`.

عندما يبدأ الوسيط الأول بـ `-` ولا يكون من خيارات ccam نفسه، يشغّل ccam الحساب الافتراضي مع كل الوسائط. والحساب الافتراضي هو الذي عيّنته باستخدام `ccam use`، أو صاحب أدنى رقم إن لم تعيّن واحدًا.

### أداة الاختيار

يفتح `ccam` دون وسائط أداة الاختيار عندما يكون الإدخال والإخراج القياسيان طرفية، ويطبع `ccam list` في غير ذلك. وتعرض لكل حساب رقمه وبريده الإلكتروني واستخدامه في نافذة الخمس ساعات وآخر استخدام له والعناصر المشتركة؛ ويظهر الاستخدام تباعًا كلما وصل رد حساب.

| المفتاح | ما يفعله |
| --- | --- |
| `Enter` | يشغّل Claude Code بالحساب المحدد |
| من `1` إلى `9` | يشغّل الحساب صاحب هذا الرقم |
| `Up`/`Down`، `k`/`j` | ينقل التحديد |
| `a` | يغلق أداة الاختيار ويشغّل معالج `ccam add` |
| `s` | يعدّل العناصر المشتركة للحساب المحدد (`Space` أو `x` للتبديل، و`Enter` للتطبيق، و`Esc` للرجوع) |
| `d` | يجعل الحساب المحدد هو الافتراضي |
| `r` | يجلب الاستخدام مرة أخرى، دون ذاكرة التخزين المؤقت القصيرة |
| `q`، `Esc`، `Ctrl+C` | خروج |

تُغلق أداة الاختيار قبل أن تنفّذ أي إجراء، فتُستعاد الطرفية قبل أن يبدأ Claude Code. وبعد `d` أو تغيير في المشاركة تُفتح مرة أخرى. ويعرض تغيير المشاركة الخطة ويسألك أولًا، كما يفعل `ccam share`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--dry-run` | يطبع البيئة والأمر اللذين كان ccam سيشغّلهما، ثم يخرج | معطّل |
| `-v`، `--version` | يطبع `ccam <version>` (فقط بوصفه الوسيط الأول) | |
| `-h`، `--help` | مساعدة ccam، بالشكل `ccam --help` أو `ccam run --help` | |

قواعد الوسائط:

- لا يقرأ ccam خياراته الخاصة إلا قبل اسم الحساب. وبعد الاسم، كل وسيط يخص Claude Code، بما في ذلك `--dry-run` و`--help`.
- يعمل `--dry-run` بالأشكال `ccam --dry-run work` و`ccam run --dry-run work` و`ccam --dry-run` (الحساب الافتراضي). وهو مدعوم، لا مرفوض.
- `ccam run <name|number> [claude args]` هو الشكل الطويل لـ `ccam <name>`. ويستخدم `ccam run` دون اسم الحساب الافتراضي.
- لا يجوز أن تكون أسماء الحسابات أسماء أوامر ccam مثل `add` أو `list` أو `run` أو `usage`، فلا يكون `ccam <word>` ملتبسًا أبدًا.

يطبع التشغيل التجريبي سطرًا واحدًا يمكنك لصقه في الصدفة:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## كيف يعمل

يمرّ كل تشغيل بالخطوات نفسها، من سطر الأوامر إلى `exec`:

<!-- motion:launch -->

يضيء كل مربع عندما تعمل خطوته، والأسطر في الأسفل هي البيئة والوسائط التي يبدأ بها Claude Code.

1. يستدعي `cmd/ccam/main.go` الدالة `app.Main`. وإذا شُغّل ccam باسم آخر، فهو [أمر مختصر](/features/shortcuts) ويذهب مباشرة إلى المشغّل.
2. ينظر `app.parse` في `os.Args` قبل أن يتدخل cobra. لا يذهب إلى `cli.Execute` إلا `ccam` وحده (أداة الاختيار)، واسم أمر من أوامر ccam (`config.Reserved`)، و`-h`/`--help` و`-v`/`--version`. وكل ما عدا ذلك تشغيل.
3. يقرأ `launcher.Launch` الملف `~/.config/ccam/config.toml`، وهو الملف الوحيد الذي يحلّله.
4. يجد `launcher.Prepare` الحساب بالاسم أو الرقم، ويحذف `--` واحدًا في البداية، ويحدد الملف التنفيذي `claude` باستخدام `ResolveClaude`. تُخزَّن نتيجة البحث في PATH مؤقتًا في `~/.local/state/ccam/claude-path` مع وقت تعديل الملف، فلا يبحث ccam في `PATH` عند كل تشغيل. والملف المتغيّر أو المفقود يؤدي إلى بحث جديد. ويُستخدم مسار `claude_bin` في `config.toml` كما هو. ويرفض ccam أي `claude` يكون هو ccam نفسه، بالمقارنة حسب inode، لأن ذلك سيؤدي إلى حلقة.
5. يقارن `checkLinks` مجلد الحساب `~/.local/share/ccam/profiles/<name>/` بإعدادات المشاركة الخاصة به باستخدام `share.Check`، الذي لا يستدعي إلا `lstat` و`readlink`، أي استدعاء نظام واحد أو اثنين لكل عنصر. ويُنشأ الرابط المفقود في الحال. وأي شيء قد يؤدي إلى فقدان بيانات يوقف التشغيل.
6. يضيف `applyMCP` خوادم MCP المشتركة (انظر أدناه).
7. يسجّل `session.Register` الجلسة. يكتب ccam علامة باسم معرّف العملية (PID) الخاص به، تحتوي على وقت بدء العملية، في `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (أو `~/.local/state/ccam/locks/` دون `XDG_RUNTIME_DIR`). وللحساب الذي يشارك سجل المحادثات، يبحث `historyCheck` أولًا عن جلسة تعمل في `history-<group>-<hash>`، حيث التجزئة هي تجزئة المجلد الحالي، ويسجّل هناك أيضًا. يحتفظ `exec` بمعرّف العملية، فتشير العلامة بعد ذلك إلى عملية Claude Code نفسها. وهي علامة لا قفل محتجز، لأن العمليات التي يشغّلها Claude Code سترث القفل وتحتفظ به بعد انتهاء الجلسة.
8. يعيّن `Plan.Exec` عنوان الطرفية إلى `claude · <name>` (عندما يكون الخطأ القياسي طرفية)، ويكتب وقت التشغيل في `~/.local/state/ccam/last-used/<name>` دون `fsync`، ويستدعي `syscall.Exec`.

يعمل Claude Code مع تعيين `CLAUDE_CONFIG_DIR` إلى مجلد الحساب و`CCAM_PROFILE` إلى اسم الحساب. ويزيل ccam أي قيم لهذين المتغيرين ورثها. وبالنسبة إلى مجلد `~/.claude` المسجل باستخدام `ccam import` على macOS، يعيّن أيضًا `CLAUDE_SECURESTORAGE_CONFIG_DIR` إلى قيمة فارغة، كي يستخدم الحساب و`claude` العادي تسجيل الدخول نفسه في سلسلة المفاتيح.

### لماذا exec

يستبدل `syscall.Exec` عملية ccam بـ `claude`، فلا يبقى ccam في الذاكرة، وتكون ضغطات المفاتيح والإشارات ورمز الخروج خاصة بـ Claude Code نفسه، وتبقى الوسائط مصفوفة، لا نصًا للصدفة أبدًا. والعمل الذي يأتي بعد الجلسة، مثل ملاحظة الخوادم المضافة باستخدام `claude mcp add`، يحدث عند التشغيل التالي: يقارن ccam وقت تعديل `.claude.json` مع `~/.local/state/ccam/mcp/<name>/seen` ولا يحلّل الملف إلا إذا تغيّر.

### كيف يُوجَّه سطر الأوامر

ثلاثة أسطر أوامر، وثلاثة مسارات عبر `internal/app`:

<!-- motion:launch-route -->

يُطابَق الأمر المختصر على `argv[0]`، و`ccam <name>` على الوسيط الأول، وكلاهما ينتهي في `syscall.Exec`. ولا يصل إلى cobra إلا أمر من أوامر ccam مثل `ccam usage`.

### خوادم MCP المشتركة

مع الاستراتيجية الافتراضية `flag`، يمرّر ccam خوادم MCP الخاصة بالمجموعة إلى Claude Code بوصفها وسيطًا، يوضع قبل وسائطك:

<!-- motion:launch-mcp -->

لا يُرجع `mcp.FlagArgs` وسيطًا إلا إذا كان ملف `mcp.json` الخاص بالمجموعة موجودًا، ويضعه `Plan.InsertArgs` مباشرة بعد `claude`. ثم يظل `--mcp-config` أو `--strict-mcp-config` الخاص بك ساريًا فوق ذلك. وهو وسيط واحد، `--mcp-config=<path>`، لأن `--mcp-config` في Claude Code يأخذ عدة قيم: فلو كُتب وسيطين، لقرأ أيضًا الطلب الذي يليه على أنه ملف إعداد.

مع `mcp_strategy = "sync"` في `config.toml`، يكتب ccam الخوادم في ملف `.claude.json` الخاص بالحساب بدلًا من ذلك، ويتخطى ذلك إذا كان الحساب يعمل بالفعل.

### ميزانية 10 ms

يضيف ccam ما لا يزيد على 10 ms (p95) قبل أن يبدأ Claude Code. لا يبني مسار التشغيل أبدًا شجرة أوامر cobra، ولا يشغّل Bubble Tea، ولا يفتح اتصال شبكة، ويفشل `TestPackageInitIsCheap` إذا استغرقت دوال `init()` في الحزم أكثر من 5 ms إجمالًا. يقيس `scripts/perf-check.sh` الأوامر `ccam --dry-run work` و`ccam work` وأمرًا مختصرًا باستخدام `hyperfine` ويفشل فوق 10 ms، ويشغّل `BenchmarkLaunchPath` في `internal/launcher` الخطوات في حلقة. ولرؤية الوقت لكل خطوة، عيّن `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## الأمان والحالات الخاصة

- الاسم أو الرقم غير المعروف يتوقف مع `no account "<name>"` ويوجّهك إلى `ccam list`. وإن لم تكن هناك حسابات على الإطلاق، يطلب منك ccam إضافة واحد باستخدام `ccam add <name>`.
- إذا أصبح مسار مشترك ملفًا عاديًا أو صار يشير إلى المكان الخطأ، يتوقف ccam مع `stopped because of share link problems:`، ويسرد المشكلات، ويطلب منك تشغيل `ccam doctor --fix`. ولا ينقل بيانات في مسار التشغيل.
- إذا كانت نسخة أخرى من Claude Code تشارك سجل المحادثات مع هذا الحساب (المجموعة نفسها) تعمل في المجلد نفسه، يحذّرك ccam من أن السجل قد يختلط ويسأل `Run anyway? [y/N]`. ودون طرفية يعرض التحذير ويتابع.
- يتخطى `--dry-run` مزامنة MCP، وإشعار الخادم الجديد، وفحص السجل، وعلامة الجلسة، والعنوان، ووقت آخر استخدام، ولا يستدعي exec. لكنه يظل يفحص روابط المشاركة، وينشئ الرابط المفقود كما يفعل التشغيل العادي.
- مع `mcp_strategy = "sync"`، تطبع المزامنة الفاشلة `ccam: MCP sync failed (launching anyway): …` ويبدأ Claude Code بالخوادم التي كانت لديه بالفعل.

## صفحات ذات صلة

- [الأوامر المختصرة](/features/shortcuts)
- [إدارة الحسابات](/features/accounts)
- [مرجع الأوامر](/docs/commands)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [كيف يعمل CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
