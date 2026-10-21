---
title: حدود الاستخدام
description: يعرض ccam usage حد الخمس ساعات والحد الأسبوعي وحد كل نموذج (Fable) لكل الحسابات دفعة واحدة على شكل أشرطة، ويضع ccam statusline اسم الحساب واستخدامه داخل Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "0e9587b79ce8"
---

يعرض `ccam usage` مقدار ما استخدمه كل حساب من حد الخمس ساعات والحد الأسبوعي والحد الأسبوعي الخاص بكل نموذج (مثل Fable)، على شكل شريط لكل حد مع الوقت المتبقي حتى إعادة ضبطه. استخدمه قبل أن تبدأ مهمة طويلة، أو عندما يخبرك Claude Code أن حدًا ما يقترب، لترى أي حساب لا تزال لديه سعة. و`ccam statusline` هو أمر سطر الحالة في Claude Code: يعرض اسم الحساب واستخدامه داخل Claude Code ويحفظ الأرقام لـ `ccam usage`.

## الاستخدام

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

يحصل كل حساب على سطر فيه رقمه واسمه وعمر أرقامه؛ وتعني `(cached)` أنها تأتي من ملف محفوظ، لا من نقطة النهاية. وتحته يعرض سطر لكل حد شريطًا من 20 خانة والنسبة المئوية والوقت المتبقي حتى إعادة ضبط النافذة: `5-hour`، ثم `Weekly` (كل النماذج)، ثم النافذة الأسبوعية الخاصة بكل نموذج، مثل `Fable`. وتظهر النافذة التي لا يملكها الحساب بالرمز `-`. وتتحول النسبة إلى الأصفر بدءًا من 80% وإلى الأحمر عند 100%. في الطرفية تظهر كل الحسابات فورًا بالحالة `fetching…`، وتُرسم أسطر كل حساب عندما يصل رده. يذكر السطر الأخير الحساب الذي سيختاره [`ccam next --auto`](/features/next) وفق القواعد نفسها (`auto_max_usage` و`auto_strategy` والجلسات العاملة، وحسابات مؤسسة آخر حساب مستخدَم فقط ما لم يُعيَّن `cross_org`، مع تخطي كل نافذة بلغت 100%، فلا يُذكر أبدًا حساب امتلأ أسبوعه)، ولا يظهر إلا إذا كان لديك أكثر من حساب. ويعرض الحساب غير المسجَّل دخوله `not logged in → ccam login <name>`، ويعرض حساب مفتاح API الرسالة `no limits (API)`.

بعد `ccam statusline install`، يعرض Claude Code سطرًا مثل `[work] · Opus · 5h 24% · 7d 41%`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--json` | يطبع التقارير بصيغة JSON: `name` و`slot` و`status` و`limits` و`source` و`at` و`error` | معطّل |
| `--watch[=N]` | يعيد الرسم كل `N` ثانية حتى تضغط Ctrl+C. و`--watch` وحده، أو قيمة 0 أو أقل، يعني 60 | معطّل |
| `-v`، `--verbose` | يضيف سطر `source:` تحت كل حساب: مصدر أرقامه (نقطة نهاية الاستخدام، أو Claude Code، أو سطر الحالة، أو ملف محفوظ) | معطّل |
| `--no-cache` | يتجاهل ذاكرة التخزين المؤقت ذات الستين ثانية ويطلب من نقطة النهاية مرة أخرى | معطّل |
| `--source live\|claude\|cache` | يطلب `live` من نقطة النهاية باستخدام رمز الحساب؛ ويجعل `claude` برنامج Claude Code يستعلم عنه (انظر أدناه)؛ ويعرض `cache` الأرقام المحفوظة فقط | `usage_source`، وإلا `live` |
| `--refresh` | يجدّد رموز الوصول المنتهية، للحسابات التي لا تعمل فقط | معطّل |
| `--dry-run` | مقبول: لا يغيّر `ccam usage` شيئًا سوى ذاكرة التخزين المؤقت الخاصة به. ويُرفض مع `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | معطّل |

أي قيمة أخرى لـ `--source` تفشل مع `--source must be live, claude or cache`.

### الاستخدام من Claude Code نفسه

مع `usage_source = "claude"` في `config.toml` (أو `--source claude`)، لا يقرأ ccam رمز تسجيل الدخول إطلاقًا. فلكل حساب يشغّل `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` في مجلد الحساب، ويرسل طلب `get_usage` الخاص بـ Claude Code عبر الإدخال القياسي (Claude Code 2.1.270 أو أحدث)؛ فيستعلم Claude Code عن الاستخدام بتسجيل دخوله هو، ويجيب بنوافذ الخمس ساعات والأسبوع والنماذج. ويمنع `--safe-mode` خطافات الحساب وخوادم MCP والمكوّنات الإضافية من البدء، ولا تُحفظ أي جلسة. ويتطلب كل حساب تشغيلًا لـ Claude Code يستغرق نحو ثانية، لذلك هذا المصدر أبطأ من `live`؛ وتُسأل الحسابات كلها في آن واحد، كلٌّ منها خلال 20 ثانية، وتنطبق الذاكرة المؤقتة لمدة 60 ثانية كالمعتاد. ويصنّف Claude Code هذا الطلب على أنه تجريبي، لذلك إذا توقف عن الاستجابة، يعرض ccam الأرقام المحفوظة ويبلّغ عنه `ccam doctor`. ولا يشغّل `--dry-run` برنامج Claude Code (لأنه يكتب في مجلد الحساب)، بل يعرض الأرقام المحفوظة. ويقبل `usage_source` القيمة `live` أو `claude`؛ أما `cache` فهو خيار لـ `--source` فقط. ولا يأخذ `ccam usage` أي وسائط.

يقرأ `ccam statusline` بيانات JSON الخاصة بسطر الحالة في Claude Code من الإدخال القياسي. وخياره الوحيد هو `--wrap <command>`: أمر سطر حالة موجود يُضاف اسم الحساب أمام أول سطر من مخرجاته. ولا يأخذ `ccam statusline install` أي خيارات. ويرفض كلاهما `--dry-run` مع `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## كيف يعمل

تأتي الأرقام من مكانين: نقطة نهاية الاستخدام غير الموثقة لدى Anthropic، أي `GET https://api.anthropic.com/api/oauth/usage`، التي تُستدعى بتسجيل الدخول الخاص بكل حساب، والقيم التي يمرّرها Claude Code إلى أمر سطر الحالة. وكل ذلك موجود في `internal/usage`.

1. يحمّل `cli.usageCmd` الملف `config.toml` وينشئ `usage.Account` واحدًا لكل حساب: مجلده `~/.local/share/ccam/profiles/<name>/`، ومجلد جلساته (لمعرفة هل يعمل)، وهل استُورد في مكانه.
2. يبدأ `Fetcher.FetchAll` روتينًا متزامنًا (goroutine) واحدًا لكل حساب. وكلها تتشارك `http.Client` واحدًا، فيُعاد استخدام الاتصالات بنقطة النهاية. قد يستغرق الاتصال 2 s، ولكل حساب 3 s إجمالًا. ولا يتبع العميل إعادة التوجيه، لأن إعادة التوجيه سترسل الرمز إلى مكان آخر، ويقرأ 1 MB على الأكثر من الرد.
3. عندما ينتهي كل حساب، يُعاد رسم صفه. والحساب الذي يفشل أو لا يستجيب في الوقت المحدد يلجأ إلى الأرقام المحفوظة، فلا يؤخّر حساب بطيء الحسابات الأخرى أبدًا.

يعمل كل حساب في روتينه الخاص مقابل نقطة النهاية؛ والحساب الذي تنتهي مهلته يعرض أرقامه المحفوظة بدلًا من ذلك.

<!-- motion:usage-fetch -->

يستجيب `personal` و`work` في الوقت المحدد ويعرضان أرقامًا حية. أما `side` فيستغرق أكثر من 3 s، لذلك يقرأ `fallback()` ملفي ذاكرة التخزين المؤقت الخاصين به ويستخدم الأحدث كتابةً منهما. يرسم `usage.Table()` كتلة لكل حساب، ويختار `usage.Best()` الحساب المذكور في السطر الأخير.

لحساب واحد، يجري `Fetcher.Fetch` هذه الفحوص بالترتيب:

<!-- motion:usage-cache -->

1. مع `--source cache`، ينتقل مباشرة إلى `fallback()`.
2. ما لم يُعيَّن `--no-cache`، تُرجَع النتيجة الحية الموجودة في `~/.local/state/ccam/usage/<name>.json` كما هي إذا كان عمرها أقل من 60 ثانية (`CacheTTL`). فتشغيل `ccam usage` مرتين متتاليتين يكلّف طلبًا واحدًا لكل حساب.
3. يقرأ `readCredentials()` الملف `.credentials.json` (أو سلسلة المفاتيح على macOS) دون تغييره، ولا يُفك إلا المفتاح `claudeAiOauth`. وعدم وجود تسجيل دخول باشتراك يعني `not logged in`، أو `no limits (API)` عندما يحتوي `.claude.json` على مفتاح API.
4. لا يُستخدم رمز الوصول المنتهي. ودون `--refresh`، يعرض الحساب أرقامه المحفوظة معلَّمة بأنها منتهية.
5. وإلا يستدعي ccam نقطة النهاية. يُعدّ الرد 401 أو 403 تسجيل دخول منتهيًا، وأي حالة أخرى أو محتوى لا يستطيع قراءته يُعدّ خطأ، وفي الحالتين يلجأ إلى الأرقام المحفوظة. ويعني الرد 429 أن نقطة النهاية تُسأل أكثر مما ينبغي: يكتب ccam الملف `<name>.pause` مع وقت السؤال التالي (`Retry-After`، أو 5 دقائق، وساعة على الأكثر)، ويعرض الأرقام المحفوظة مع `rate limited by the usage endpoint; asking again after 15:10`، ولا يستدعي نقطة النهاية لذلك الحساب حتى ذلك الحين، حتى مع `--no-cache`. فإعادة المحاولة فورًا لن تؤدي إلا إلى إبقاء الحد قائمًا.
6. يُكتب الرد السليم في `<name>.json` مع وقت قياسه.

يقرأ `fallback()` الملفين `<name>.json` و`<name>.statusline.json` ويحتفظ بالأحدث. وقبل أن يعرض الأرقام المحفوظة، يتحقق من وقت إعادة ضبط كل نافذة: النافذة التي أُعيد ضبطها منذ حفظ الأرقام تظهر بقيمة 0%، لأن النسبة القديمة لم تعد صحيحة. ولهذا يعرض `personal` أعلاه `0%` دون وقت لإعادة الضبط.

لنقطة النهاية تنسيقان للاستجابة. الأحدث فيه مصفوفة `limits` تحتوي على المدخلات `session` (الخمس ساعات) و`weekly_all` و`weekly_scoped` (لكل نموذج). والأقدم فيه الكائنات `five_hour` و`seven_day` و`seven_day_opus` و`seven_day_sonnet`. يقرأ `ParseLive` كليهما: تتقدم مصفوفة `limits`، وتكمل الكائنات القديمة ما لا تغطيه. والمحتوى الذي لا يحتوي على أي حدود يُبلَّغ عنه على أنه تنسيق متغيّر، ويخبرك `ccam doctor` بذلك.

### تجديد الرموز

رموز التحديث تتبدّل: تجديد واحد منها يمنحك رمز تحديث جديدًا ويبطل القديم. لو جدّد ccam تسجيل الدخول بينما يعمل Claude Code بالحساب نفسه، لسُجّل خروج تلك الجلسة. لذلك يقرأ `ccam usage` الرموز ولا يكتبها أبدًا، ما لم تمرّر `--refresh`. عندها فإن `Fetcher.refresh`:

1. يرفض الحسابات المسجلة في مكانها باستخدام `ccam import`، لأن `claude` العادي قد يستخدم ذلك التسجيل دون علم ccam؛
2. يرفض تسجيلات الدخول المخزنة في سلسلة المفاتيح على macOS؛
3. يأخذ ملف قفل، `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (أو تحت `~/.local/state/ccam/locks/` عندما لا يكون ذلك المتغير معيّنًا)، فلا تجدّد عمليتان من ccam الرمز نفسه مرتين أبدًا؛
4. تحت القفل، يرفض إذا كانت جلسة للحساب تعمل، ويتخطى الطلب إذا كانت عملية أخرى قد جدّدت الرمز بالفعل؛
5. يرسل رمز التحديث إلى نقطة نهاية الرموز، ويكتب الرموز الجديدة في `.credentials.json` باستخدام `fsutil.AtomicWrite`، مع الإبقاء على كل المفاتيح الأخرى في الملف.

### جامع بيانات سطر الحالة

يشغّل Claude Code أمر سطر الحالة ويمرّر إليه مستند JSON. يقرأ ccam منه اسم النموذج و`rate_limits`:

<!-- motion:usage-statusline -->

1. يجد `runStatusline` الحساب من `CCAM_PROFILE`، الذي يعيّنه ccam عندما يشغّل Claude Code، وإلا فبمطابقة `CLAUDE_CONFIG_DIR` مع مجلدات الحسابات.
2. يقرأ `ParseStatusLine` القيمتين `rate_limits.five_hour` و`rate_limits.seven_day`. لا يحتوي JSON على حدود للطلبات في حسابات مفاتيح API، أو قبل أول رد في الجلسة؛ وعندها لا يُحفظ شيء.
3. يكتب `SaveStatusLine` الملف `~/.local/state/ccam/usage/<name>.statusline.json`. ويستخدمه `ccam usage` عندما لا تستجيب نقطة النهاية.
4. يطبع ccam السطر `[name] · model · 5h N% · 7d N%`. ومع `--wrap`، يشغّل أمرك عبر `sh -c` بالإدخال نفسه ويطبع بدلًا من ذلك `[name]` متبوعًا بأول سطر من مخرجات الأمر. وعلى Windows يشغّل الأمر كما يفعل Claude Code: باستخدام Git Bash إذا كان مثبّتًا، وإلا باستخدام PowerShell.
5. عندما يكون الاستخدام في نافذة الخمس ساعات أو النافذة الأسبوعية 90% أو أكثر ويوجد حساب آخر، ينتهي السطر بـ `→ ccam next --auto --continue`، مثلًا `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

للأمر المغلَّف ثانيتان. بعدها يُنهى، ويعرض السطر اسم الحساب وحده، فلا يستطيع أمر بطيء أن يعطّل سطر الحالة في Claude Code.

يكتب `ccam statusline install` المدخل `statusLine` في ملف `settings.json` لكل حساب. وعندما يكون لدى الحساب أمر لسطر الحالة بالفعل، يحتفظ به ccam بالشكل `--wrap '<old command>'`. وعلى Windows يُكتب مسار ccam بشرطات مائلة أمامية و`~` لمجلد المنزل (`~/.local/bin/ccam.exe statusline`)، وهو ما يفهمه Git Bash وPowerShell كلاهما. وملف `settings.json` المشترك عبر مجموعة يُكتب مرة واحدة، عبر رابطه، فتحصل عليه كل حسابات المجموعة. وتشغيله مرة أخرى يطبع `(already set)` للملفات التي تستخدم ccam بالفعل.

## الأمان والحالات الخاصة

- لا يغيّر `ccam usage` أبدًا ملف رموز ما لم تمرّر `--refresh`. وعندما ينتهي تسجيل دخول، شغّل الحساب مرة واحدة باستخدام `ccam <name>` وسيجدّده Claude Code.
- لا يلمس `--refresh` أبدًا الحسابات التي تعمل أو الحسابات المستوردة أو تسجيلات الدخول في سلسلة المفاتيح. والتجديد المرفوض يعرض الأرقام المحفوظة والسبب.
- نقطة نهاية الاستخدام غير موثقة علنًا وقد تتغير. وعندما يحدث ذلك، تظل ترى الأرقام المحفوظة، ويبلّغ `ccam doctor` عن تغيّر التنسيق.
- عندما لا توجد أرقام محفوظة ولا رد، يعرض سطر الحساب `lookup failed:` والسبب، مثلًا `no saved usage` مع `--source cache`.
- لا تحتوي ملفات ذاكرة التخزين المؤقت إلا على النسب وأوقات إعادة الضبط، ولا تحتوي على رموز أبدًا.

## صفحات ذات صلة

- [حدود الاستخدام والتبديل بين الحسابات](/docs/usage-limits)
- [التبديل إلى الحساب التالي](/features/next)
- [سير العمل عند بلوغ حد الاستخدام](/guides/usage-limit-workflow)
- [الأوامر](/docs/commands)
