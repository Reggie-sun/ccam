---
title: الفحص والإصلاح
description: افحص روابط المشاركة وأذونات الملفات وتسجيلات الدخول وإصدار Claude Code باستخدام ccam doctor، وأصلح ما يمكن إصلاحه بأمان باستخدام ccam doctor --fix.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

يفحص `ccam doctor` تثبيت ccam لديك: ملف Claude Code التنفيذي، ومجلد كل حساب وتسجيل دخوله وروابط المشاركة فيه، وخوادم MCP المشتركة، ونقطة نهاية الاستخدام، وأوامرك المختصرة. ويصلح `ccam doctor --fix` ما يمكنه إصلاحه دون فقدان بيانات، وينقل أولًا أي ملف يستبدله إلى مجلد نسخ احتياطية. شغّله عندما يتوقف `ccam <name>` بسبب مشكلة في رابط مشاركة، أو بعد أن تستعيد ملفات يدويًا، أو بعد تحديث Claude Code.

## الاستخدام

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

كل نتيجة سطر واحد مع علامة: `✓` سليم، و`·` معلومة، و`!` تحذير، و`✗` خطأ. وتبدأ نتائج الحسابات باسم الحساب.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

مع `--fix`، تنتهي كل نتيجة أُصلحت بـ `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--fix` | يصلح ما يمكن إصلاحه بأمان (مع نسخ احتياطية للملفات المتغيّرة) | معطّل |
| `--dry-run` (عام) | يجري كل الفحوص ولا يصلح شيئًا، حتى مع `--fix` | معطّل |

لا يأخذ `ccam doctor` أي وسائط. وهو يدعم `--dry-run`: يطبع `ccam --dry-run doctor --fix` النتائج نفسها التي يطبعها `ccam doctor` ولا يغيّر شيئًا.

رمز الخروج هو `1` عندما تبقى نتيجة `✗` واحدة على الأقل دون إصلاح، و`0` في غير ذلك. ولا تُفشل التحذيرات والمعلومات التشغيل أبدًا. وعندما يفشل التشغيل دون `--fix`، ينهي ccam المخرجات بسطر فارغ ثم `` `ccam doctor --fix` repairs what it can. ``. والإعداد الذي لا يمكن تحميله ينتهي أيضًا بـ `1` وخطأ يبدأ بـ `ccam:`.

## كيف يعمل

يحمّل `internal/cli/doctor.go` الملف `~/.config/ccam/config.toml` ويستدعي `doctor.Run` في `internal/doctor/doctor.go`، الذي يجري الفحوص بترتيب ثابت ويُرجع النتائج. ويُمرَّر `--fix` على أنه `Options.Fix`، ويُعيَّن إلى false في التشغيل التجريبي.

<!-- motion:doctor-checks -->

تعمل الكتلة الوسطى مرة واحدة لكل حساب. الفحوص، وما يفعله `--fix` لكل منها:

1. **إصدار Claude Code** (`checkClaude`). يحدد الملف التنفيذي `claude` كما يفعل المشغّل، ويشغّل `claude --version`. إذا لم يبدأ الإصدار بـ `2.1.302`، وهو الإصدار الذي تحقق ccam من بنية مجلده (`doctor.LayoutVersion`)، يضيف `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. وغياب الملف التنفيذي أو فشل `--version` خطأ. لا إصلاح.
2. **بيانات الاعتماد في البيئة** (`checkAuthEnv`). يحذّر عندما يكون `CLAUDE_CODE_OAUTH_TOKEN` أو `ANTHROPIC_API_KEY` أو `ANTHROPIC_AUTH_TOKEN` معيّنًا: فكل حساب يُشغَّل من تلك الصدفة يرثه، ويفضّله Claude Code على تسجيل دخول الحساب نفسه. لا إصلاح؛ ألغِ تعيينه بنفسك.
3. **المجلد** (`checkProfile`). غياب مجلد الحساب خطأ وينهي فحوص ذلك الحساب. والأذونات التي فيها أي بت للمجموعة أو للآخرين تحذير؛ ويشغّل `--fix` الأمر `chmod 0700`. لا توجد بتات أذونات في Windows، لذا يُتخطى هذا الفحص وفحص `0600` أدناه هناك؛ فالمجلد الموجود في ملف تعريف المستخدم محمي بقائمة التحكم في الوصول الخاصة به.
4. **ملف تسجيل الدخول** (`checkCredentials`). يُفحص `.credentials.json` باستخدام `lstat`. الرابط الرمزي خطأ، لأن كل حساب يجب أن يملك رمزه الخاص، ولا يُصلح. والوضع المختلف عن `0600` تحذير؛ ويشغّل `--fix` الأمر `chmod 0600`. ثم يُقرأ تسجيل الدخول (على macOS من سلسلة المفاتيح أولًا). وعدم وجود تسجيل دخول تحذير يوجّهك إلى `ccam login <name>`. ورمز التحديث المنتهي هو الخطأ `login expired`؛ أما رمز الوصول المنتهي فمعلومة فقط، لأن Claude Code يجدّده عند التشغيل التالي. ولا يُصلح أي من هذه.
5. **روابط المشاركة** (`share.Check`). يقارن كل مسار في `share.Layout` بإعدادات المشاركة للحساب، مستخدمًا `lstat` و`readlink` فقط. ويبلّغ عن أربعة أنواع: رابط مفقود، ورابط يشير إلى مكان آخر، وملف حقيقي حيث يجب أن يكون رابط، ورابط متبقٍّ لعنصر أصبح محليًا الآن. الرابط المفقود تحذير، والبقية أخطاء.
6. **المدخلات غير المعروفة** (`checkUnknownEntries`). يسرد الأسماء في المستوى الأعلى من مجلد الحساب التي لا يعرفها ccam (`share.Known`)، كإشارة إلى أن بنية Claude Code تغيّرت. معلومة فقط.
7. **مسارات الإضافات** (`checkPlugins`). تخزّن سجلات الإضافات (`installed_plugins.json` و`known_marketplaces.json`) مسارات مطلقة. وفي `plugins/` الخاص بكل مجموعة، تتعطل المسارات التي تمر عبر مجلد حساب واحد عندما يُزال ذلك الحساب أو يُعاد تسميته. يعيد `--fix` كتابتها إلى مجلد المجموعة باستخدام `profile.RewritePluginPaths`، عبر كتابة ذرّية لملف السجل.
8. **خوادم MCP** (`checkMCP`). يحمّل ملف `mcp.json` لكل مجموعة مرة واحدة. الملف الذي لا يمكن قراءته خطأ؛ والخادم الذي يفشل في `mcp.Validate` تحذير. لا إصلاح.
9. **تنسيق استجابة الاستخدام** (`checkUsageFormat`). يطلب من نقطة نهاية الاستخدام، دون ذاكرة التخزين المؤقت ومع الحد المعتاد البالغ 3 ثوانٍ، بيانات حساب تلو الآخر حتى يستجيب أحدها. وإذا لم يعد بالإمكان تحليل الاستجابة، يحذّر من أن `ccam usage` سيلجأ إلى ذاكرة التخزين المؤقت لسطر الحالة. وهذا هو الفحص الوحيد الذي يستخدم الشبكة.
10. **الأوامر المختصرة** (`checkShims`). فقط عندما يكون لحساب ما أمر مختصر. الروابط المفقودة في `~/.local/bin` تحذير؛ ويستدعي `--fix` الدالة `shim.Install`، التي لا تنشئ إلا الروابط المفقودة ولا تلمس الملفات الأخرى أبدًا. ومجلد bin غير الموجود في `PATH` تحذير ولا يُصلح.
11. **المجلدات** (`checkFolders`). فقط عند وجود إدخالات [`ccam folder`](/features/folder). الإدخال الذي لا يوجد حسابه خطأ، والإدخال الذي اختفى مجلده تحذير. ولا يُصلح أي منهما، لأن إزالة إدخال ستغيّر الحساب الذي يبدأ هناك.

### إصلاح روابط المشاركة

ينشئ مسار التشغيل الرابط المفقود في الحال، لكنه يتوقف عند أي شيء قد يؤدي إلى فقدان بيانات ويوجّهك إلى هنا:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

لأي مشكلة في روابط المشاركة، يستدعي `--fix` الدالة `share.Build` مع إعدادات المشاركة للحساب بوصفها الحالة الحالية والحالة المستهدفة معًا، ثم `Plan.Apply`. يعيد `Build` استنتاج كل خطوة مما هو موجود على القرص، وهي الخطة نفسها التي يستخدمها `ccam share`، لذلك يتبع الإصلاح القواعد نفسها التي يتبعها تفعيل المشاركة.

<!-- motion:doctor-repair -->

هنا `settings.json` ملف حقيقي. ولأن لدى المجموعة بالفعل ملف `settings.json` مختلفًا، يأخذ التعارض خياره الافتراضي `KeepGroup`، فتبقى نسخة المجموعة. وينتقل ملف الحساب إلى `.ccam-backup/<time>/settings.json` داخل مجلد الحساب، ويربط `os.Symlink` المسار بالمجموعة. حالات أخرى:

- إذا لم تكن لدى المجموعة نسخة بعد، ينتقل ملف الحساب إلى المجموعة ويصبح أول محتواها.
- في حالة المجلد، تنتقل الملفات التي لا تملكها المجموعة إلى المجموعة؛ ويذهب الباقي إلى النسخة الاحتياطية.
- بالنسبة إلى `history.jsonl`، الخيار الافتراضي هو `merge by line`: أسطر المجموعة، ثم أسطرك الجديدة. وتُحفظ النسخة السابقة للمجموعة في `.ccam-backup/<time>/group/`.
- يُستبدل الرابط إلى مجموعة أخرى. والرابط إلى مجلد خارج ccam يذهب إلى النسخة الاحتياطية، ويصبح محتواه بداية المجموعة إذا كانت المجموعة فارغة.
- يُزال الرابط المتبقي لعنصر أصبح محليًا الآن، ويحصل الحساب على نسخة من محتوى المجموعة.

باستثناء إزالة الروابط، كل خطوة تنقل البيانات إلى المجموعة أو إلى مجلد النسخ الاحتياطية. لا يُحذف شيء.

## الأمان والحالات الخاصة

- لا يغيّر `ccam doctor` وحده أي شيء أبدًا. استخدم `--dry-run` مع `--fix` لتكون متأكدًا في السكربتات.
- لا يتحقق `doctor --fix` هل الحساب يعمل. أنهِ Claude Code لحساب ما قبل إصلاح روابط المشاركة فيه، لأن الجلسة تظل تكتب في تلك الملفات.
- إذا فشل إصلاح المشاركة، يبلّغ ccam بـ `could not repair share links: <error>` لذلك الحساب ويتخطى فحص المدخلات غير المعروفة فيه. وتظل الحسابات الأخرى تُفحص.
- إذا فشل `ccam share` في منتصف الطريق، يحفظ الإعدادات المستهدفة أولًا ويطلب منك الإكمال باستخدام `ccam doctor --fix`.
- الملف `.credentials.json` المربوط رمزيًا وتسجيل الدخول المنتهي وبيانات الاعتماد في البيئة يُبلَّغ عنها ولا تُغيَّر أبدًا.
- إعادة كتابة سجلات الإضافات تتم ذرّيًا لكنها لا تُنسخ إلى `.ccam-backup/`.

## صفحات ذات صلة

- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [نسخ إعداد موجود](/features/migrate)
- [حدود الاستخدام والتبديل بين الحسابات](/docs/usage-limits)
- [الأوامر](/docs/commands)
