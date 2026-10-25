---
title: الفحص والإصلاح
description: افحص روابط المشاركة وأذونات الملفات وتسجيلات الدخول وإصدار Claude Code باستخدام ccam doctor، وأصلح ما يمكن إصلاحه بأمان باستخدام ccam doctor --fix.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

يفحص `ccam doctor` تثبيت ccam لديك: ملف Claude Code التنفيذي، ومجلد كل حساب وتسجيل دخوله وروابط المشاركة فيه، وخوادم MCP المشتركة، ونقطة نهاية الاستخدام، وأوامرك المختصرة. ويصلح `ccam doctor --fix` ما يمكنه إصلاحه دون فقدان بيانات، وينقل أولًا أي ملف يستبدله إلى مجلد نسخ احتياطية. شغّله عندما يتوقف `ccam <name>` بسبب مشكلة في رابط مشاركة، أو بعد أن تستعيد ملفات يدويًا، أو بعد تحديث Claude Code.

## الاستخدام

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

يضع السطر الأول علامة لكل حساب دفعة واحدة، فترى أين المشكلات قبل أن تقرأ أي شيء آخر. وتحته يحصل كل حساب على صندوق، وتوضع النتائج التي لا تخص أي حساب (Claude Code، ومتغيرات البيئة، والإضافات، وخوادم MCP، والأوامر المختصرة، وقواعد المجلدات، والمجموعات، والإعدادات، و`auto_switch`) في صندوق `General`. ويعرض عنوان الصندوق أسوأ علامة للحساب وعدد مشكلاته. وكل صف داخله مجال واحد جرى فحصه، دائمًا `Folder` ثم `Login` ثم `Share links` بهذا الترتيب، ثم أي مجالات أخرى، مع علامة: `✔` سليم، و`ℹ` معلومة، و`⚠` تحذير، و`✘` خطأ. ويعرض المجلد الذي لا مشكلة فيه `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

مع `--fix`، تنتهي كل نتيجة أُصلحت بـ `→ fixed`، ولا تُحسب بعد ذلك مشكلةً:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

في الطرفية تتسع الصناديق لعرضها وتلتف الرسائل الطويلة، بألوان الطرفية نفسها (ويوقفها `NO_COLOR`). ولا يكون الصندوق أبدًا أضيق من 29 عمودًا، فلا تُقطع أي كلمة في طرفية ضيقة. وعند إرسال المخرجات إلى أنبوب أو ملف، لا تكون فيها ألوان ولا تُلَف أي رسالة أو تُقطع، فيجد `ccam doctor | grep settings.json` السطر كاملًا؛ ويصبح عرض كل صندوق حينها بعرض أطول سطر. ويعرض تبويب Health في لوحة التحكم الصناديق نفسها.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--fix` | يصلح ما يمكن إصلاحه بأمان (مع نسخ احتياطية للملفات المتغيّرة) | معطّل |
| `--dry-run` (عام) | يجري كل الفحوص ولا يصلح شيئًا، حتى مع `--fix` | معطّل |

لا يأخذ `ccam doctor` أي وسائط. وهو يدعم `--dry-run`: يطبع `ccam --dry-run doctor --fix` نتائج `ccam doctor`، ويضيف ما سيضعه `--fix` جانبًا، ولا يغيّر شيئًا:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

رمز الخروج هو `1` عندما تبقى نتيجة `✘` واحدة على الأقل دون إصلاح، و`0` في غير ذلك. ولا تُفشل التحذيرات والمعلومات التشغيل أبدًا. وعندما يفشل التشغيل دون `--fix`، ينهي ccam المخرجات بسطر فارغ ثم `` `ccam doctor --fix` repairs what it can. ``. والإعداد الذي لا يمكن تحميله ينتهي أيضًا بـ `1` وخطأ يبدأ بـ `ccam:`.

## كيف يعمل

يحمّل `internal/cli/doctor.go` الملف `~/.config/ccam/config.toml` ويستدعي `doctor.Run` في `internal/doctor/doctor.go`، الذي يجري الفحوص بترتيب ثابت ويُرجع النتائج، كلٌّ منها مع الحساب الذي يخصه والمجال الذي فحصه. ويجمعها `tui.HealthReport` في الصناديق. ويُمرَّر `--fix` على أنه `Options.Fix`، ويُعيَّن إلى false في التشغيل التجريبي.

<!-- motion:doctor-checks -->

تعمل الكتلة الوسطى مرة واحدة لكل حساب. الفحوص، وما يفعله `--fix` لكل منها:

1. **إصدار Claude Code** (`checkClaude`). يحدد الملف التنفيذي `claude` كما يفعل المشغّل، ويشغّل `claude --version`. وغياب الملف التنفيذي، أو `claude_bin` ليس ملفًا تنفيذيًا، أو فشل `--version` خطأ. لا إصلاح.
2. **بيانات الاعتماد في البيئة** (`checkAuthEnv`). يحذّر عندما يكون `CLAUDE_CODE_OAUTH_TOKEN` أو `ANTHROPIC_API_KEY` أو `ANTHROPIC_AUTH_TOKEN` معيّنًا: فكل حساب يُشغَّل من تلك الصدفة يرثه، ويفضّله Claude Code على تسجيل دخول الحساب نفسه. لا إصلاح؛ ألغِ تعيينه بنفسك.
3. **المجلد** (`checkProfile`). غياب مجلد الحساب خطأ وينهي فحوص ذلك الحساب. والأذونات التي فيها أي بت للمجموعة أو للآخرين تحذير؛ ويشغّل `--fix` الأمر `chmod 0700`. ومجلد `~/.claude` المستورد الذي فيه ملف `.claude.json` خاص به بجانب `~/.claude.json` معلومة: فعندها يحتفظ `ccam main` و`claude` العادي بتفاصيل تسجيل دخول وخوادم MCP منفصلة. لا توجد بتات أذونات في Windows، لذا يُتخطى هذا الفحص وفحص `0600` أدناه هناك؛ فالمجلد الموجود في ملف تعريف المستخدم محمي بقائمة التحكم في الوصول الخاصة به.
4. **ملف تسجيل الدخول** (`checkCredentials`). يُفحص `.credentials.json` باستخدام `lstat`. الرابط الرمزي خطأ، لأن كل حساب يجب أن يملك رمزه الخاص، ولا يُصلح. والوضع المختلف عن `0600` تحذير؛ ويشغّل `--fix` الأمر `chmod 0600`. ثم يُقرأ تسجيل الدخول (على macOS من سلسلة المفاتيح أولًا). وعدم وجود تسجيل دخول، أو ملف تسجيل دخول لا يحتوي على واحد، تحذير يوجّهك إلى `ccam login <name>`. ورمز التحديث المنتهي، أو تسجيل الدخول الذي مسحه Claude Code بعد أن رفضه الخادم، هو الخطأ `` login expired → `ccam login <name>` ``؛ أما رمز الوصول المنتهي فمعلومة فقط، لأن Claude Code يجدّده عند التشغيل التالي. وتسجيل الدخول المحفوظ في Windows Credential Manager معلومة، لأن ccam لا يستطيع قراءته. ويعرض تسجيل الدخول الصالح البريد الإلكتروني الذي سُجّل دخول المجلد به الآن والخطة (`login valid (dev@company.com, Claude Team · Acme)`). ولا يُصلح أي من هذه؛ ويسجّل `ccam login --expired` دخول كل حساب يحتاج إلى ذلك.
5. **روابط المشاركة** (`share.Check`). يقارن كل مسار في `share.Layout` بمجموعات الحساب، مستخدمًا `lstat` و`readlink` فقط. ويبلّغ عن خمسة أنواع: رابط مفقود، ورابط يشير إلى مكان آخر، وملف حقيقي حيث يجب أن يكون رابط، ورابط متبقٍّ لعنصر أصبح محليًا الآن، ورابط صحيح مجلد مجموعته مفقود (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`؛ ولا يفحص التشغيل هذا). الرابط المفقود ومجلد المجموعة المفقود تحذيران، والبقية أخطاء. ويُبلَّغ أيضًا عن النسخ من المجموعة الذي توقف في منتصفه: يضع `--fix` نسخة مرحلية كاملة في مكانها، لكن عندما يكون المسار قد أُنشئ من جديد منذ ذلك الحين (عبر تشغيل)، يذكر المجلد الذي فيه النسخة الكاملة لتتولى أمره بنفسك. أما المسارات التي أُضيفت إلى عنصر في إصدار لاحق من ccam (`file-history/` و`tasks/` و`paste-cache/` للسجل، و`keybindings.json` للإعدادات) فهي تحذير فقط ما دام الحساب يملك نسخته الخاصة منها: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. ولا يتوقف التشغيل بسببها أبدًا.
6. **المدخلات غير المعروفة** (`checkUnknownEntries`). يسرد الأسماء في المستوى الأعلى من مجلد الحساب التي لا يعرفها ccam (`share.Known`)، كإشارة إلى أن بنية Claude Code تغيّرت. معلومة فقط.
7. **مسارات الإضافات** (`checkPlugins`). تخزّن سجلات الإضافات (`installed_plugins.json` و`known_marketplaces.json`) مسارات مطلقة. وفي `plugins/` الخاص بكل مجموعة، تتعطل المسارات التي تمر عبر مجلد حساب واحد عندما يُزال ذلك الحساب أو يُعاد تسميته. يعيد `--fix` كتابتها إلى مجلد المجموعة باستخدام `profile.RewritePluginPaths`، عبر كتابة ذرّية لملف السجل.
8. **خوادم MCP** (`checkMCP`). يحمّل ملف `mcp.json` لكل مجموعة مرة واحدة. الملف الذي لا يمكن قراءته خطأ؛ والخادم الذي يفشل في `mcp.Validate` تحذير. والخوادم الموجودة في `~/.claude.json` الخاص بـ `claude` العادي التي لا يحمّلها أي حساب تحذير يذكر أمر `ccam migrate --mcp` الذي ينسخها (مع `--group <a|b>` عندما تشارك عدة مجموعات MCP). والخوادم التي لا يملكها إلا حساب واحد في مجموعة معلومة مع أمر `ccam mcp import` الذي يشاركها، والحساب الذي لم يسجّل الدخول إلى خوادم MCP سجّلت حسابات أخرى الدخول إليها يحصل على تحذير بأن يسجّل الدخول باستخدام `/mcp` (ولا ينسخ ccam تسجيلات الدخول أبدًا). لا إصلاح.
9. **تنسيق استجابة الاستخدام** (`checkUsageFormat`). يطلب من نقطة نهاية الاستخدام، دون ذاكرة التخزين المؤقت، بيانات حساب تلو الآخر حتى يستجيب أحدها، وكل ذلك ضمن حد واحد قدره 4 ثوانٍ. ومع `usage_source = "claude"` يسأل Claude Code عن حساب واحد بدلًا من ذلك. وإذا لم يعد بالإمكان تحليل الاستجابة، يحذّر من أن `ccam usage` سيلجأ إلى ذاكرة التخزين المؤقت لسطر الحالة؛ وعندما لا يستجيب أي حساب، يذكر ذلك على أنه معلومة (`could not check the usage response format: …`). وهذا هو الفحص الوحيد الذي يستخدم الشبكة.
10. **الأوامر المختصرة** (`checkShims`). فقط عندما يكون لحساب ما أمر مختصر. الروابط المفقودة في `~/.local/bin` تحذير، وكذلك الأوامر المختصرة التي أنشأها ccam ولا تزال تشغّل ccam أقدم، استُبدل أو نُقل (`shortcuts in <bin> that run another ccam (one replaced or moved): …`)؛ ويستدعي `--fix` الدالة `shim.Install`، التي لا تنشئ ولا تجدّد إلا الروابط التي أنشأها ccam وسجّلها في `~/.local/state/ccam/shortcuts.json`، ولا تلمس الملفات الأخرى أبدًا. ولا ينتهي التحذير بـ `→ fixed` إلا عندما يُنشأ كل رابط من تلك الروابط؛ وإذا تعذّر إنشاء أحدها، يُضاف السبب. والأمر المختصر الذي يشغّل ccam آخر لا يزال مثبّتًا يعمل، وهو معلومة؛ ويوجّهه `--fix` إلى ccam هذا. ومجلد bin غير الموجود في `PATH` تحذير ولا يُصلح.
11. **المجلدات** (`checkFolders`). فقط عند وجود إدخالات [`ccam folder`](/features/folder). الإدخال الذي لا يوجد حسابه خطأ، والإدخال الذي اختفى مجلده تحذير. ولا يُصلح أي منهما، لأن إزالة إدخال ستغيّر الحساب الذي يبدأ هناك.
12. **المجموعات والإعدادات** (`checkGroups`). العضو في مجموعة الذي ليس حسابًا، والحساب الذي أصبح اسمه منذ ذلك الحين أمرًا من أوامر ccam (شغّله باستخدام `ccam run <name>` أو أعد تسميته)، والقيمة غير المعروفة لـ `theme` أو `usage_source`، وأسطر `config.toml` التي لا يقرؤها ccam، كلها تحذيرات. والأخيرة تكون عادةً إعدادًا مكتوبًا تحت جدول `[[group]]` أو `[[profile]]` أو `[[folder]]`، حيث ينتمي إلى ذلك الجدول (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`)؛ وترفض الأوامر التي تغيّر الإعداد العمل حتى يُصلح. لا إصلاح.
13. **تسجيلات الدخول** (`checkIdentities`). يقرأ من `.claude.json` في كل مجلد الحسابَ الذي سُجّل الدخول به. والمجلد المسجّل دخوله الآن بحساب غير الذي سجّله `ccam login` (مثلًا بعد `/login` داخل جلسة) تحذير؛ والحساب الذي لا سجل تسجيل دخول له، مثل الحساب المستورد أو الذي أنشأه ccam أقدم، يُقارَن ببريده الإلكتروني المسجَّل. والمجلد المسجّل دخوله بحساب Claude نفسه الذي يستخدمه مجلد آخر تحذير أيضًا: فالاثنان يتشاركان حدود الاستخدام. والمجلد الذي يحمل رمز تسجيل الدخول نفسه الذي يحمله مجلد آخر، أي مجلد منسوخ (بما في ذلك `~/.claude` الخاص بـ `claude` العادي)، خطأ، لأن أيهما يجدّده أولًا يُخرج الآخر؛ شغّل `ccam login <name>` لتمنحه تسجيل دخول خاصًا به. لا إصلاح.
14. **`auto_switch`** (`checkAutoSwitch`). فقط عند `auto_switch = true`. تذكر معلومة في صندوق `General` ما يفعله `auto_switch` وتربط إلى [وصفه](/features/next#auto-switch). ولكل خطة ذات مقاعد بين الحسابات (Claude Team وClaude Enterprise) سطر يسمّيها، ويذكر أن `auto_switch` ينقل الجلسة بين مقاعد المؤسسة الواحدة، ويقتبس قاعدة المقاعد لدى Anthropic مع رابط مركز المساعدة، ويترك القرار لك. وClaude Code أقدم من 2.1.214، الذي لا يُخبر الخطافات بعمليته (`CLAUDE_PID`)، تحذير هناك، إذ لا يمكن نقل أي جلسة. والحساب الذي لا يملك إضافة ccam، أو الذي أُوقفت فيه الإضافة، يحصل على تحذير في صندوقه مع الأمر الذي يصلح ذلك داخل Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (أو `/plugin enable ccam@ccam`). وتُستثنى حسابات مفاتيح API. وعلى Windows يكون السطر الوحيد `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` لا إصلاح.

قبل الحسابات، يكون **تغيير المجموعات** الذي توقف في منتصفه خطأً في صندوق `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. ويُكمله `--fix` بالترتيب الذي خُطط له، مع الإجابات التي أُعطيت لتعارضاته؛ وحتى ذلك الحين لا يبدأ أي تغيير آخر في المجموعات أو إعادة تسمية أو إزالة. والحسابات المشمولة به لا تُنشأ لها إلا روابطها المفقودة حتى يكتمل. وإذا تعذّر إكماله، يقول الخطأ إن إزالة ذلك الملف تعني التخلي عن بقية التغيير. ويعرض `ccam --dry-run doctor --fix` الترتيب (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) وما سيضعه جزء كل حساب جانبًا. والتغيير الذي ينفّذه ccam آخر في هذه اللحظة معلومة، وما دام يعمل لا يصلح `--fix` أي روابط مشاركة (`its share links are not repaired while another ccam changes the groups; …`). ويمسك `--fix` بالقفل `~/.local/state/ccam/group-change.lock` أثناء الإصلاح، ويعيد قراءة `config.toml` بعد إكمال تغيير.

### إصلاح روابط المشاركة

ينشئ مسار التشغيل الرابط المفقود في الحال، لكنه يتوقف عند أي شيء قد يؤدي إلى فقدان بيانات ويوجّهك إلى هنا:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

لأي مشكلة في روابط المشاركة، يستدعي `--fix` الدالة `share.Build` مع إعدادات المشاركة للحساب بوصفها الحالة الحالية والحالة المستهدفة معًا، ثم `Plan.Apply`. يعيد `Build` استنتاج كل خطوة مما هو موجود على القرص، وهي الخطة نفسها التي يستخدمها أي تغيير عبر `ccam group`، لذلك يتبع الإصلاح القواعد نفسها التي يتبعها الانضمام إلى مجموعة.

<!-- motion:doctor-repair -->

هنا `settings.json` ملف حقيقي. ولأن لدى المجموعة بالفعل ملف `settings.json` مختلفًا، يأخذ التعارض خياره الافتراضي `KeepGroup`، فتبقى نسخة المجموعة. وينتقل ملف الحساب إلى `.ccam-backup/<time>/settings.json` داخل مجلد الحساب، ويربط `os.Symlink` المسار بالمجموعة. ويخبرك `--fix` بذلك: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. حالات أخرى:

- إذا لم تكن لدى المجموعة نسخة بعد، ينتقل ملف الحساب إلى المجموعة ويصبح أول محتواها.
- في حالة المجلد، تنتقل الملفات التي لا تملكها المجموعة إلى المجموعة. أما الملفات التي تملكها المجموعة مطابقة بايتًا ببايت فتُسقط، والملفات المختلفة تذهب إلى النسخة الاحتياطية، وتُذكر في ملاحظة (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- يُنشأ مجلد المجموعة المفقود أولًا.
- بالنسبة إلى `history.jsonl`، الخيار الافتراضي هو `merge by line`: أسطر المجموعة، ثم أسطرك الجديدة. وتُحفظ النسخة السابقة للمجموعة في `.ccam-backup/<time>/group/`.
- يُستبدل الرابط إلى مجموعة أخرى. والرابط إلى مجلد خارج ccam يذهب إلى النسخة الاحتياطية، ويصبح محتواه بداية المجموعة إذا كانت المجموعة فارغة.
- يُزال الرابط المتبقي لعنصر أصبح محليًا الآن، ويحصل الحساب على نسخة من محتوى المجموعة.

باستثناء إزالة الروابط وإسقاط الملفات التي تملكها المجموعة مطابقة بايتًا ببايت، كل خطوة تنقل البيانات إلى المجموعة أو إلى مجلد النسخ الاحتياطية. لا يضيع شيء.

## الأمان والحالات الخاصة

- لا يغيّر `ccam doctor` وحده أي شيء أبدًا. استخدم `--dry-run` مع `--fix` لتكون متأكدًا في السكربتات.
- لحساب قيد التشغيل، أيًّا كانت طريقة تشغيل Claude Code، لا ينشئ `doctor --fix` إلا الروابط المفقودة ومجلدات المجموعات المفقودة، وهذا لا يُضيّع شيئًا. وعندما يتبقى شيء آخر، يضيف `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``، لأن الجلسة تظل تكتب في الملفات التي سينقلها الإصلاح.
- إذا فشل إصلاح المشاركة، يبلّغ ccam بـ `could not repair share links: <error>` لذلك الحساب ويتخطى فحص المدخلات غير المعروفة فيه. وتظل الحسابات الأخرى تُفحص.
- إذا فشل تغيير عبر `ccam group` في منتصف الطريق، تكون المجموعات الجديدة قد حُفظت بالفعل ويُسجَّل التغيير في `~/.local/state/ccam/group-change.json`؛ ويطلب منك ccam الإكمال باستخدام `ccam doctor --fix`، الذي ينفّذ بقيته، بما في ذلك خطوات MCP.
- الملف `.credentials.json` المربوط رمزيًا وتسجيل الدخول المنتهي وبيانات الاعتماد في البيئة يُبلَّغ عنها ولا تُغيَّر أبدًا.
- إعادة كتابة سجلات الإضافات تتم ذرّيًا لكنها لا تُنسخ إلى `.ccam-backup/`.

## صفحات ذات صلة

- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [نسخ إعداد موجود](/features/migrate)
- [حدود الاستخدام والتبديل بين الحسابات](/docs/usage-limits)
- [الأوامر](/docs/commands)
