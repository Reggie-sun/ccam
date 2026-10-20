---
title: المشاركة بين الحسابات
description: كيف يربط ccam share الإعدادات وCLAUDE.md والإضافات والسجل بمجلد مجموعة، مع خطة ونسخ احتياطية وخيارات للتعارضات، وكيف يعيد off النسخ إلى الحساب.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 16de29bbdd57
---

لكل حساب مجلده الخاص، لكن معظم ما فيه لا يحتاج إلى أن يكون منفصلًا. يفعّل `ccam share` المشاركة أو يوقفها لحساب واحد وعنصر واحد في كل مرة: يُنقل ملف الحساب أو مجلده إلى مجلد مجموعة ويُستبدل برابط رمزي، فتقرأ كل حسابات المجموعة النسخة نفسها وتكتب فيها. استخدمه بعد `ccam add` لتغيير ما يشاركه حساب ما، أو لوضع بضعة حسابات في مجموعة خاصة بها.

## الاستخدام

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

دون خيارات العناصر، يطبع `ccam share` الحالة الحالية:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

مع الخيارات، يطبع ccam خطة، ويسألك عن كل ملف يختلف عن نسخة المجموعة، ويسأل `Proceed? [y/N]` قبل أن يغيّر أي شيء. يشارك الحساب الجديد كل شيء ما عدا سجل المحادثات؛ وتصف صفحة [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing) كل عنصر من جهة المستخدم.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | خوادم MCP | دون تغيير |
| `--settings <on\|off\|group>` | `settings.json` | دون تغيير |
| `--memory <on\|off\|group>` | ملف `CLAUDE.md` العام | دون تغيير |
| `--extensions <on\|off\|group>` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/`، `rules/` | دون تغيير |
| `--history <on\|off\|group>` | `projects/` و`history.jsonl` | دون تغيير |
| `--fresh` | عند إيقاف المشاركة، يبدأ فارغًا بدلًا من نسخ محتوى المجموعة | معطّل |
| `--prefer <group\|mine\|merge>` | كيفية حسم التعارضات | السؤال؛ ومع `--yes` الخيار الافتراضي لكل عنصر |
| `-y`، `--yes` | يعمل دون سؤال؛ وتتبع التعارضات `--prefer` | معطّل |
| `-h`، `--help` | مساعدة `share` | |
| `--dry-run` (عام) | مدعوم: يسأل عن التعارضات، ويطبع الخطة، ولا يغيّر شيئًا | معطّل |

القيمة `on` تعني المجموعة `default`، و`off` تعني محليًا (غير مشترك)، وأي قيمة أخرى اسم مجموعة. تصبح أسماء المجموعات أسماء مجلدات تحت `~/.local/share/ccam/groups/`، لذلك يفشل `ccam share side --history ../x` مع `invalid group name: "../x"`. يُحدَّد الحساب بالاسم أو بالرقم، ولا يتغيّر إلا ما تذكره من عناصر.

## كيف يعمل

الملفات المعنية:

- `~/.config/ccam/config.toml`: الجدول `[profile.share]` لكل حساب، بقيمة واحدة لكل عنصر.
- `~/.local/share/ccam/profiles/<name>/`: مجلد `CLAUDE_CONFIG_DIR` الخاص بالحساب.
- `~/.local/share/ccam/groups/<group>/`: النسخ المشتركة.

1. يقرأ `internal/cli/share.go` الإعداد، ويطبّق كل قيمة `--<item>` على نسخة من إعدادات الحساب باستخدام `config.Share.Set`، ويرفض إذا كان الحساب يعمل.
2. يقارن `share.Build` الإعدادات الحالية بالمستهدفة ويسرد الخطوات. وهو يقرأ القرص فقط: `os.Lstat` و`os.Readlink` لكل مسار في `share.Layout`. والخطة قائمة عمليات: النقل إلى خانة فارغة في المجموعة (تبنٍّ)، ونقل الملفات التي لا تملكها المجموعة (دمج مجلد)، وكتابة تعارض محسوم في المجموعة، والنسخ الاحتياطي، وإزالة الرابط، والربط، والنسخ إلى الحساب، وإعادة كتابة مسارات الإضافات، ودمج سجل الإضافات.
3. تُحسم التعارضات، بالسؤال أو بـ `--prefer`.
4. يطبع `Plan.Describe` الخطة، ويسألك ccam.
5. تُكتب الإعدادات المستهدفة في `config.toml` قبل أن يُنقل أي شيء. وإذا توقف `Plan.Apply` في منتصف الطريق، يبني `ccam doctor --fix` وإعادة التشغيل خطتهما نحو الهدف المحفوظ ويُكملان التغيير.
6. ينفّذ `Plan.Apply` الخطوات بالترتيب. كل خطوة تنقل البيانات إلى المجموعة أو إلى مجلد النسخ الاحتياطية؛ لا يُحذف شيء.

`share.Layout` قائمة سماح بالمسارات التي يغطيها كل عنصر. وأي شيء ليس فيها لا يمكن أن يصل إلى مجموعة أبدًا. يذكر `share.Never` الملفين `.credentials.json` و`.claude.json`، ويتوقف `Build` مع خطأ داخلي إذا ظهر أي منهما يومًا في البنية. وليس للعنصر `mcp` مسارات: فقائمة الخوادم موجودة داخل `.claude.json`، لذلك تعمل [مشاركة MCP](/features/mcp) عبر ملف `mcp.json` الخاص بالمجموعة بدلًا من ذلك، ولا يغيّر `--mcp` إلا الإعداد.

هذا ما يحدث لمجلد واحد، `projects/`، عندما يفعّل `side` مشاركة السجل وفي المجموعة محتوى بالفعل:

<!-- motion:share-on -->

إذا لم يكن لدى المجموعة شيء في ذلك المسار بعد، فالحساب هو أول من يشاركه، ويُنقل مجلده إلى هناك كما هو، دون نسخة احتياطية. وإلا ينقل `mergeDir` كل ملف لا تملكه المجموعة إلى المجموعة، وتبقى الملفات الموجودة في الجهتين في مجلد الحساب، فيُحتفظ بنسخة المجموعة. ثم يذهب مجلد الحساب إلى `.ccam-backup/<time>/` داخل مجلد الحساب، ويضع `os.Symlink` رابطًا مكانه. تُسمّى ملفات المحادثات باسم معرّف الجلسة، لذلك تُدمج السجلات القادمة من عدة حسابات دون تعارضات.

يستخدم النقل `rename(2)`، ويلجأ إلى النسخ ثم الحذف عبر أنظمة الملفات المختلفة. والرابط النسبي داخل الشجرة المنقولة الذي كان يشير إلى خارجها يُعاد كتابته رابطًا مطلقًا (`fsutil.PortableLink`)، فيظل يعمل من المكان الجديد. تخزّن سجلات `plugins/` (`installed_plugins.json` و`known_marketplaces.json`) مسارات مطلقة تمر عبر مجلد الحساب؛ فيوجّهها ccam إلى مجلد المجموعة، وعندما يكون لدى المجموعة سجل بالفعل، يدمج مدخلات الحساب فيه.

إيقاف عنصر يعكس الرابط:

<!-- motion:share-off -->

لا يزيل ccam الرابط إلا إذا كان يشير إلى داخل مجلد المجموعات، ثم ينسخ المحتوى الحالي للمجموعة إلى الحساب باستخدام `CopyTree`. ومع `--fresh` يُتخطّى النسخ ويبدأ الحساب دون شيء في ذلك المسار. وتحتفظ المجموعة بمحتواها للحسابات الأخرى. ونقل عنصر من مجموعة إلى أخرى يزيل رابطه ويربطه بالمجموعة الجديدة، تاركًا المجموعة القديمة كما هي؛ وتذكر الخطة ذلك في سطر `Note:`.

الملف الواحد الموجود في الجهتين بمحتوى مختلف تعارض:

<!-- motion:share-conflict -->

يسرد `share.NewConflict` الخيارات التي تناسب الملف. يعرض `history.jsonl` الخيارات `merge by line` و`keep group` و`use mine`، مع الدمج أولًا. وتعرض ملفات `.json` الأخرى `keep group` و`use mine` و`merge JSON`. وأي ملف آخر يعرض `keep group` و`use mine`. والخيار الأول هو الافتراضي:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

يدمج `merge JSON` (`MergeJSONBytes`) الكائنات بشكل متكرر، ويضم المصفوفات دون تكرار، ويحتفظ بقيمة المجموعة عندما تعيّن الجهتان المفتاح نفسه بقيمتين مختلفتين، فلا ترى الحسابات الأخرى أي تغيير في السلوك. ويحتفظ `merge by line` بأسطر المجموعة، ثم يضيف أسطرك التي لا تملكها. وكل خيار ما عدا `keep group` يحفظ أولًا نسخة المجموعة في `.ccam-backup/<time>/group/` داخل مجلد الحساب، ثم يكتب النتيجة باستخدام `fsutil.AtomicWrite`. ويذهب ملفك أنت دائمًا إلى النسخة الاحتياطية.

يحسم `--prefer group` أو `mine` أو `merge` كل تعارض دون سؤال؛ وتعني `merge` الخيار `merge JSON` أو `merge by line`، أيهما يعرضه الملف، والملف الذي لا يعرض أيًا منهما يحتفظ بالخيار الافتراضي. ويأخذ `--yes` وحده الخيار الافتراضي لكل ملف.

عند كل تشغيل، يقارن `share.Check` كل مسار في البنية بالإعداد باستخدام `lstat` و`readlink` فقط، مما يبقيه خفيفًا بما يكفي لمسار التشغيل. ويُنشأ الرابط المفقود في الحال، لأنه لا يوجد ما يُفقد. أما الملف الحقيقي حيث يجب أن يكون رابط، أو الرابط إلى المكان الخطأ، أو الرابط المتبقي، فيوقف التشغيل:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

يشغّل `ccam doctor --fix` الدالة `share.Build` مع الإعدادات المحفوظة في الجهتين، فيستنتج كل خطوة مما هو موجود على القرص، ويحتفظ بمحتوى المجموعة عند التعارضات، ويأخذ نسخة احتياطية من أي شيء يستبدله.

### على Windows

لا ينشئ Windows الروابط الرمزية إلا في وضع المطور (Developer Mode، في الإعدادات > النظام > للمطورين) أو للمسؤول. يحاول `fsutil.Link` إنشاء رابط رمزي أولًا. ودون هذا الحق، يصبح المجلد المشترك (`agents` و`skills` و`projects` وغيرها) نقطة وصل مجلدات (directory junction)، يستطيع أي مستخدم إنشاءها، وتعمل بالطريقة نفسها مع Claude Code. أما الملف المشترك (`settings.json` و`CLAUDE.md` و`history.jsonl`) فلا يمكن أن يكون إلا رابطًا رمزيًا، لذا يتوقف ccam برسالة `Windows links files only in Developer Mode` ويقترح تفعيل هذا الوضع أو إبقاء ذلك العنصر محليًا. ويقرأ `share.Check` نقاط الوصل بوصفها روابط (`fsutil.LinkTarget`)، وإزالة مجلد حساب تزيل نقطة الوصل نفسها، ولا تزيل أبدًا محتوى المجموعة الذي تشير إليه.

## الأمان والحالات الخاصة

- **الحسابات التي تعمل.** يرفض ccam تغيير المشاركة ما دام Claude Code يعمل بالحساب (`account work is running; quit it first, then change its sharing`)، لأن الجلسة تظل تكتب في الملفات التي يجري نقلها، والأسطر التي تضيفها في أثناء ذلك ستضيع. ويظل `--dry-run` يطبع الخطة.
- **النسخ الاحتياطية.** تحصل كل خطة على مجلد `.ccam-backup/<time>/` خاص بها داخل مجلد الحساب، مع لاحقة مرقّمة عندما تبدأ خطتان في الثانية نفسها. ولا يحذفها ccam أبدًا.
- **الإخفاقات الجزئية.** إذا فشلت خطوة، يطبع ccam الخطأ ومكان الملفات الأصلية، ويطلب منك الإكمال باستخدام `ccam doctor --fix`. والهدف محفوظ بالفعل، فيتابع الإصلاح نحوه.
- **روابطك الخاصة.** إذا كان مسار ما رابطًا أنشأته بنفسك بالفعل، مثلًا إلى مستودع dotfiles، وكانت المجموعة فارغة، ينسخ ccam هدف الرابط إلى المجموعة وينقل الرابط نفسه إلى النسخة الاحتياطية. ولا يُمسّ الهدف.
- **ملفات تسجيل الدخول.** لا يُربط `.credentials.json` و`.claude.json` أبدًا بمجموعة، ولا يُنقلان إليها ولا يُنسخان.
- **حسابان في مجلد واحد.** الحسابات التي تتشارك السجل تكتب في الملفات نفسها. وعندما يبدأ حساب ثانٍ من مجموعة السجل نفسها في مجلد يعمل فيه حساب بالفعل، يحذّرك ccam ويسألك قبل أن يتابع.

## صفحات ذات صلة

- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [خوادم MCP المشتركة](/features/mcp)
- [نسخ إعداد موجود](/features/migrate)
- [كل الأوامر](/docs/commands)
