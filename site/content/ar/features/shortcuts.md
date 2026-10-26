---
title: الأوامر المختصرة
description: امنح حسابًا أمرًا خاصًا به مثل ccwork باستخدام ccam shortcut، وأعد إنشاء الروابط في ~/.local/bin باستخدام ccam shims install.
updated: 2026-10-26
order: 3
translation: machine
sourceHash: 066e7e622517
---

الأمر المختصر أمر تختار اسمه بنفسك، ويشغّل Claude Code بحساب واحد: يفعل `ccwork --resume` ما يفعله `ccam work --resume`. وهو رابط إلى ccam في `~/.local/bin`، لذلك يعمل في bash وzsh والسكربتات وأي شيء آخر يشغّل البرامج من `PATH`. الأوامر المختصرة اختيارية؛ أنشئ واحدًا للحسابات التي تشغّلها كثيرًا.

## الاستخدام

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

يؤكد ccam كل تغيير:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

للحساب أمر مختصر واحد على الأكثر، ويُخزَّن باسم `shortcut` في `~/.config/ccam/config.toml`. يسأل معالج `ccam add` السؤال `Create a shortcut command for <name>? [y/N]` ثم `Shortcut name [cc<name>]:`. وخارج المعالج، لا ينشئ `ccam add` أمرًا مختصرًا إلا مع `--shortcut`. ويعرض `ccam list` الأمر المختصر لكل حساب، أو `-`.

يمرّر الأمر المختصر كل وسيط إلى Claude Code، بالطريقة نفسها التي يفعلها `ccam <name>` بعد الاسم. ويشمل ذلك `--dry-run` و`--help`: لمعاينة تشغيل، استخدم `ccam --dry-run work`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--remove` | يحذف الأمر المختصر للحساب ورابطه | معطّل |
| `-h`، `--help` | مساعدة الأمر | |
| `--dry-run` | يدعمه `ccam shortcut`: يتحقق من الاسم ويطبع `Would set the shortcut of work to ccwork.` أو `Would remove ccwork, the shortcut of work.` دون حفظ أو إنشاء رابط | معطّل |

لا يأخذ `ccam shims install` أي وسائط، ويرفض `--dry-run` مع الرسالة `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

قواعد الاسم، وتُفحص بهذا الترتيب:

- من 1 إلى 64 حرفًا من `a-z` و`0-9` و`.` و`-` و`_`، ويبدأ بحرف أو رقم؛
- ليس `ccam` أو `claude` أو أمرًا من أوامر ccam مثل `run` أو `usage`؛
- ليس أمرًا مختصرًا لحساب آخر؛
- ليس ملفًا موجودًا في `~/.local/bin`، إلا إذا كان رابطًا أنشأه ccam (مثل الأمر المختصر الخاص بهذا الحساب)؛
- ليس اسم برنامج تنفيذي آخر في `PATH`.

يُفحص `config.toml` وفق القواعد نفسها عند تحميله. فإذا عدّلت أمرًا مختصرًا يدويًا بحيث يخالف النمط، أو يكون `claude` أو `ccam`، أو يستخدمه حسابان، يفشل كل أمر مع رسالة تذكر الحساب حتى تُصلحه.

تحديد أمر مع `--remove` معًا يتوقف مع `give a command or --remove, not both`. وتعيين الاسم الذي يملكه الحساب بالفعل يطبع `work: ccwork (unchanged)` ويعيد إنشاء رابطه إذا كان مفقودًا.

## كيف يعمل

تعيين اسم يشغّل `checkShortcut`، الذي يستدعي `shim.Check` مع مجلد bin والمسار الحقيقي لملف ccam التنفيذي:

<!-- motion:shortcuts-check -->

ينتقل الاسم عبر الفحوص بالترتيب. وبجانب كل فحص اسم يرفضه: `Work` لا يطابق النمط، و`run` أمر من أوامر ccam، و`ccside` يخص حسابًا آخر، و`ccfile` ملف موجود في `~/.local/bin`، و`git` سيحجب `/usr/bin/git`. ولا يُحفظ ويُربط إلا الاسم الذي يجتاز كل الفحوص.

1. يطابق ccam الاسم مع `^[a-z0-9][a-z0-9._-]{0,63}$`، ثم مع `config.Reserved`، ثم مع `cfg.FindShortcut` (`shim.Check`). ويستدعي `os.Lstat` على `~/.local/bin/<name>`، ولكل مجلد آخر في `PATH` يستدعي `os.Stat` ويرفض الملف التنفيذي الذي ليس رابطًا إلى ccam.
2. يحفظ ccam الملف `config.toml` باستخدام `fsutil.AtomicWrite` (ملف مؤقت، ثم `fsync`، ثم إعادة تسمية).
3. يحذف ccam رابط الأمر المختصر القديم، ولكن فقط إذا كان رابطًا إلى ccam (`shim.Remove`).
4. يمرّ ccam على كل أمر مختصر مُعدّ بترتيب الأسماء وينشئ رابط كل أمر مفقود (`shim.Install`)، باستخدام `os.Symlink` إلى ملف ccam التنفيذي، الذي يُحدَّد مساره باستخدام `os.Executable` و`filepath.EvalSymlinks`. والملف الذي يعترض أمرًا مختصرًا واحدًا لا يوقف الأوامر الأخرى؛ ويذكر الخطأ كلًّا منها.

يسجّل ccam كل رابط ينشئه، والملف التنفيذي الذي يشغّله كل رابط، في `~/.local/state/ccam/shortcuts.json`. ويُعدّ الرابط خاصًا بـ ccam عندما يُرجع `os.Readlink` مسار ccam، أو عندما يشير الرابط ومسار ccam إلى الملف نفسه. والرابط المسجَّل الذي نُقل ملفه التنفيذي أو استُبدل منذ ذلك الحين (نُقل ملف ccam التنفيذي إلى مجلد آخر، أو أُعيد التثبيت) قديم، ويوجّهه ccam إلى ملف ccam التنفيذي هذا (`shim.Install`). أما الرابط المسجَّل إلى نسخة أخرى مثبّتة من ccam فيُترك كما هو. وكل عملية حذف وكل فحص استبدال يمرّ عبر هذه الاختبارات، فلا يلمس ccam ملفات لم ينشئها، مثل نسختك الخاصة من ccam أو رابط أنشأته بنفسك.

عندما تشغّل أمرًا مختصرًا، فالاسم الذي كتبته هو كل ما يحتاجه ccam:

<!-- motion:shortcuts-link -->

تجد الصدفة `ccwork` في `PATH`، ويتبع النواة الرابط ويشغّل ملف ccam التنفيذي، ويظل `argv[0]` هو `ccwork`. في `app.Main`، أي اسم غير `ccam` يحمّل `config.toml` مرة واحدة ويستدعي `shim.Lookup`، الذي يجد الحساب الذي يساوي `shortcut` فيه `filepath.Base(argv[0])`. ثم يتبع ccam الخطوات نفسها التي يتبعها [`ccam <name>`](/features/launch) مع الإعداد المحمَّل مسبقًا، ويمرّر `argv[1:]` إلى `claude` دون تغيير (`launcher.LaunchWith`).

لماذا رابط رمزي وليس اسمًا مستعارًا في الصدفة أو سكربتًا مغلِّفًا: الرابط الرمزي يعمل في كل صدفة وفي السكربتات، ويتجنب تشغيل صدفة إضافية كان سيكلّف نحو 1 إلى 2 ms، ويظل التشغيل يقرأ ملفًا واحدًا. يقيس `scripts/perf-check.sh` الأمر المختصر باستخدام `hyperfine` مقابل الحد نفسه البالغ 10 ms (p95) المطبق على `ccam <name>`.

### على Windows

الأمر المختصر هو `ccwork.exe` في `%USERPROFILE%\.local\bin`، ويحذف ccam اللاحقة `.exe` (دون تمييز حالة الأحرف) قبل البحث (`shim.Name`)، لذا تشغّل `ccwork` و`ccwork.exe` و`CCWORK.EXE` كلها الحساب `work`. ينشئ ccam أولًا رابطًا رمزيًا، يسمح به Windows في وضع المطور أو للمسؤول. ودون هذا الحق ينشئ رابطًا صلبًا (hard link): اسمًا ثانيًا للملف `ccam.exe` نفسه، لا يحتاج إلى أي صلاحية. ويُعدّ الرابط الصلب تابعًا لـ ccam عندما يكون هو الملف `ccam.exe` نفسه.

ويشغّل الرابط الصلب دائمًا الملف الذي أُنشئ منه، لذا عندما يضع `ccam update` ملفًا جديدًا مكان `ccam.exe`، فإنه يجدّد كل أمر مختصر قائم على رابط صلب باستخدام `shim.Relink` ويطبع `✓ shortcut ccwork renewed`. وتشغيل `install.ps1` مرة أخرى يفعل الشيء نفسه عبر `ccam shims install` (`OK 2 shortcuts now run the new ccam`)، وكذلك يفعل `ccam shims install` و`ccam doctor --fix` لأي أمر مختصر هو نسخة من إصدار أقدم من ccam. ويبحث فحص الاسم عن `<name>` مع كل لاحقة في `PATHEXT` (`.exe` و`.cmd` وغيرهما) في مجلدات `PATH` الأخرى.

اكتب الأمر المختصر في صدفة (PowerShell أو cmd.exe أو Git Bash). فالرابط الرمزي الذي يُفتح من مستكشف الملفات (Explorer) أو مربع الحوار Run أو `Start-Process` يبدأ بوصفه `ccam` العادي، فيفتح لوحة التحكم بدلًا من الحساب.

## الأمان والحالات الخاصة

- **رابط لا يملكه أي حساب.** إذا كان في `~/.local/bin` رابط إلى ccam باسم لا يستخدمه أي حساب، مثلًا بعد أن عدّلت `shortcut` في `config.toml` يدويًا، فلا يلجأ ccam إلى الحساب الافتراضي. بل يتوقف مع `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. عيّن الاسم مرة أخرى باستخدام `ccam shortcut`، أو احذف الرابط.
- **الملفات الموجودة.** لا يستبدل ccam أبدًا ملفًا لم ينشئه. ويبلّغ `ccam shims install` بـ `<path> already exists and is not a ccam link, so it was left alone` ويظل ينشئ الروابط الأخرى. ثم يطبع `ccam add` و`ccam import` السطر `shortcut: ccwork (not created; see above)`.
- **نقل ملف ccam التنفيذي.** إذا نقلت ملف ccam التنفيذي، تشير الروابط التي أنشأها ccam إلى مسار لم يعد موجودًا. ويوجّهها `ccam shims install` (أو `ccam doctor --fix`) إلى ملف ccam التنفيذي الذي تشغّله به.
- **إعادة التسمية والإزالة.** يحتفظ `ccam rename` باسم الأمر المختصر، ويظل الرابط يعمل، لأنه يشير إلى ccam ويبحث ccam عن الاسم في `config.toml`. ويحذف `ccam remove` رابط الحساب. ويحذف `ccam shortcut <account> <new>` الرابط القديم قبل أن ينشئ الجديد.
- **PATH.** عندما ينشئ ccam رابطًا ولا يكون `~/.local/bin` في `PATH` لديك، يطبع `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` متبوعًا بسطر `export PATH=...` الذي تضيفه. وعلى Windows يطبع بدلًا من ذلك سطر PowerShell الذي يضيف المجلد إلى `PATH` الخاص بمستخدمك.
- **الفحص لاحقًا.** يبلّغ `ccam doctor` بـ `missing shortcuts in /home/me/.local/bin: ccwork`، وبالأوامر المختصرة التي تشغّل إصدارًا أقدم من ccam، ويحذّر عندما لا يكون المجلد في `PATH`. وينشئ `ccam doctor --fix` الروابط المفقودة ويجدّد القديمة، ولا يقول `→ fixed` إلا لما أنشأه. والأمر المختصر الذي يشغّل ccam آخر لا يزال مثبّتًا يعمل، لذا يسرده `ccam doctor` على أنه معلومة فقط؛ ويوجّهه `ccam doctor --fix` إلى ملف ccam التنفيذي الذي تشغّله به، ويعرضه `ccam add` و`ccam import` بالشكل ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [إدارة الحسابات](/features/accounts)
- [مرجع الأوامر](/docs/commands#shortcuts)
- [تثبيت ccam](/docs/install)
