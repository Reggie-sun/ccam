---
title: المشاركة عبر المجموعات
description: تشارك مجموعة ccam عناصر مختارة بين حساباتها عبر مجلد واحد؛ وينقل ccam group البيانات وفق خطة، مع نسخ احتياطية وخيارات للتعارضات.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 3f94dd13dbf1
---

لكل حساب مجلده الخاص، لكن معظم ما فيه لا يحتاج إلى أن يكون منفصلًا. تحدد **المجموعة** ما يُشارَك ومن يشاركه: يربط كل حساب في المجموعة تلك العناصر بمجلد واحد، `~/.local/share/ccam/groups/<group>/`، فتقرأ الحسابات كلها النسخة نفسها وتكتب فيها. ولا تُشارَك رموز تسجيل الدخول أبدًا.

يمكن أن يكون الحساب في عدة مجموعات، لكن كل عنصر يأتي من مجموعة واحدة فقط. ومن الإعدادات الشائعة مجموعة للإعدادات والإضافات، ومجموعة ثانية لسجل المحادثات بين حسابين:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

ينشئ أول حساب تضيفه المجموعة `default`، وتنضم إليها الحسابات اللاحقة ما لم تحدد غير ذلك.

## الاستخدام

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

كل أمر ينقل البيانات يطبع خطة، ويسألك عن كل ملف يختلف عن نسخة المجموعة، ويسأل `Proceed? [y/N]` قبل أن يغيّر أي شيء:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

يمكنك إجراء التغييرات نفسها في تبويب **Groups** في [لوحة التحكم](/features/dashboard)؛ واضغط `s` على أي حساب لتختار مجموعاته.

## العناصر

| العنصر | ما يغطيه |
| --- | --- |
| `mcp` | خوادم MCP، عبر ملف `mcp.json` الخاص بالمجموعة ([خوادم MCP المشتركة](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | ملف `CLAUDE.md` العام |
| `extensions` | `agents/` و`commands/` و`skills/` و`plugins/` و`output-styles/` و`rules/` |
| `history` | `projects/` و`history.jsonl`، فيستطيع حساب آخر استخدام `--continue` |

إذا لم تشارك أي مجموعة من مجموعات الحساب عنصرًا ما، يبقى ذلك العنصر في مجلد الحساب نفسه.

## الأوامر والخيارات

| الأمر | ما يفعله |
| --- | --- |
| `ccam group` أو `ccam group list` | يسرد المجموعات |
| `ccam group show <group>` | ما تشاركه المجموعة، وحساباتها، ومجلدها، وخوادم MCP فيها |
| `ccam group add <group> --items <items> [--members <accounts>]` | ينشئ مجموعة؛ وتنتقل إليها عناصر حساباتها |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | يغيّر ما تشاركه المجموعة، لكل حساب فيها |
| `ccam group join <group> <account>...` | يضع حسابات في المجموعة |
| `ccam group leave <group> <account>...` | يُخرج حسابات منها؛ ويحتفظ كل منها بنسخة، بما في ذلك خوادم MCP |
| `ccam group rename <group> <new-name>` | يعيد التسمية في الإعداد فقط |
| `ccam group remove <group> [--leave]` | يزيل مجموعة |

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--from <account>` | مع `add` و`set`: الحساب الذي تدخل نسخته إلى المجموعة أولًا | الحساب الافتراضي، ثم حسب الرقم |
| `--replace` | مع `join`: يغادر أولًا المجموعات التي تشارك العناصر نفسها | معطّل |
| `--fresh` | الحسابات التي تتوقف عن مشاركة عنصر تبدأ دونه بدلًا من البدء بنسخة منه (ويظل المجلد المستورد يحتفظ بنسخة) | معطّل |
| `--leave` | مع `remove`: يُخرج حسابات المجموعة أولًا | مطلوب عندما تكون فيها حسابات |
| `--prefer <group\|mine\|merge>` | كيفية حسم التعارضات | السؤال؛ ومع `--yes` الخيار الافتراضي لكل ملف |
| `-y`، `--yes` | يعمل دون سؤال؛ وتتبع التعارضات `--prefer` | معطّل |
| `--dry-run` (عام) | يطبع الخطة، ولا يغيّر شيئًا | معطّل |

تُحدَّد الحسابات بالاسم أو بالرقم، والعناصر قائمةً مفصولة بفواصل. ويأخذ `ccam add` الخيار `--group <group>` (ويمكن تكراره) أو `--no-group`.

## كيف يعمل

الملفات المعنية:

- `~/.config/ccam/config.toml`: إدخالات `[[group]]`: `name` و`items` و`members`، و`dir` عندما يختلف اسم المجلد.
- `~/.local/share/ccam/profiles/<name>/`: مجلد `CLAUDE_CONFIG_DIR` الخاص بالحساب.
- `~/.local/share/ccam/groups/<folder>/`: النسخ المشتركة.

عند قراءة `config.toml`، تُحَل مجموعات كل حساب إلى مجلد واحد لكل عنصر. ويستغرق ذلك نحو 25 µs لخمسة حسابات، فلا يزال التشغيل يقرأ `config.toml` وحده. والإعداد الذي يحصل فيه حساب على عنصر من مجموعتين خطأ، ولا يبدأ شيء حتى يُصلح.

يعمل التغيير على نسخة من الإعداد:

1. يعدّل الأمر النسخة (`AddGroup` و`SetGroupItems` و`JoinGroup` …)، وهذه الدوال تتحقق من القواعد.
2. يقارن ccam القديم بالجديد، ولكل حساب تتغيّر عناصره يبني خطة (`internal/group` و`share.Build`). وتُبنى الخطط واحدة بعد أخرى، وترى كل منها ما ستضعه الخطط السابقة في مجلدات المجموعات، فلا ينقل نسخته إليها إلا الحساب الأول. وتأتي النسخة الأولى للمجموعة، بهذا الترتيب، من: مجلد المجموعة (إن كانت فيه نسخة بالفعل)، ثم الحساب المحدد في `--from`، ثم الحساب الافتراضي، ثم الحسابات الأخرى حسب الرقم.
3. يرفض ccam إذا كان أحد تلك الحسابات يعمل، أيًّا كانت طريقة تشغيل Claude Code (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``)، وعلى Windows دون وضع المطور إذا كان ملف ما سيحتاج إلى رابط.
4. يحسم ccam التعارضات (بالسؤال، أو باستخدام `--prefer`)، ويطبع الخطة ويطلب التأكيد.
5. قبيل الحفظ مباشرة، يتحقق ccam مرة أخرى. فإذا بدأ حساب بينما كان ccam ينتظر إجاباتك، يتوقف التغيير؛ ويتوقف أيضًا إذا غيّرت عملية ccam أخرى المجموعات نفسها في الأثناء (`the groups changed while this plan waited; nothing was moved, make the change again`). ويُحفظ الإعداد الجديد قبل أن يُنقل أي شيء، ويُسجَّل التغيير في `~/.local/state/ccam/group-change.json`. وإذا توقفت خطوة في منتصف الطريق، يبقى التسجيل: يُكمل `ccam doctor --fix` التغيير بترتيبه، مع الإجابات التي أعطيتها، بما في ذلك خطوات MCP، وحتى ذلك الحين لا يبدأ أي تغيير آخر في المجموعات (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. تُنفَّذ خطوات MCP أولًا. وعندما تُخرج الحسابات الإضافات من مجموعة، تُوجَّه قوائم المكوّنات الإضافية للمجموعة إلى مجلد المجموعة حيثما تمر مساراتها عبر مجلد حساب (`group plugin lists:` في الخطة). ثم تُنفَّذ الخطط بالترتيب. كل خطوة تنقل البيانات إلى المجموعة أو إلى مجلد النسخ الاحتياطية؛ ولا تُحذف بدلًا من نسخها احتياطيًا إلا الملفات المطابقة لنسخة المجموعة.

يحتفظ ccam بقائمة سماح بالمسارات التي يغطيها كل عنصر (`share.Layout`). وأي شيء ليس فيها لا يمكن أن يصل إلى مجموعة أبدًا. ويسمّي ccam أيضًا الملفين `.credentials.json` و`.claude.json` بوصفهما لا يُشارَكان أبدًا (`share.Never`)، ويتوقف بخطأ داخلي إذا ظهر أي منهما يومًا في البنية (`Build`). وليس للعنصر `mcp` مسارات: فقائمة الخوادم موجودة داخل `.claude.json`، لذلك تعمل [مشاركة MCP](/features/mcp) عبر ملف `mcp.json` الخاص بالمجموعة.

هذا ما يحدث لمجلد واحد، `projects/`، عندما يبدأ `side` مشاركة السجل وفي المجموعة محتوى بالفعل:

<!-- motion:share-on -->

إذا لم يكن لدى المجموعة شيء في ذلك المسار بعد، فالحساب هو أول من يشاركه، ويُنقل مجلده إلى هناك كما هو. وإلا يُنقل كل ملف لا تملكه المجموعة إلى المجموعة، وتبقى الملفات الموجودة في الجهتين في مجلد الحساب، فيُحتفظ بنسخة المجموعة. وتُحذف الملفات المطابقة لملفات المجموعة، وتذكر الخطة الملفات المختلفة (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

ثم يذهب مجلد الحساب إلى `.ccam-backup/<time>/` داخل مجلد الحساب، ويحل رابط محله. وإذا كان مدخل المجموعة نفسه رابطًا إلى المدخل الخاص بالحساب المنضم، يذهب محتوى الحساب إلى المجموعة مكانه. تُسمّى ملفات المحادثات باسم معرّف الجلسة، لذلك تُدمج السجلات القادمة من عدة حسابات دون تعارضات. وتخزّن سجلات `plugins/` مسارات مطلقة تمر عبر مجلد الحساب؛ فيوجّهها ccam إلى مجلد المجموعة، ويدمج مدخلات الحساب في السجل الذي تملكه المجموعة بالفعل.

مغادرة المجموعة تعكس الرابط:

<!-- motion:share-off -->

ينسخ ccam أولًا المحتوى الحالي للمجموعة إلى `.ccam-backup/incomplete-<name>` في مجلد الحساب، ولا يستبدل الرابط إلا بعد اكتمال النسخ. وإذا توقف النسخ، يبقى الرابط في مكانه وينسخ `ccam doctor --fix` من جديد (وتُظهر خطة هذا الإكمال `restore <entry>: … (copied back from the group before a stop)`). وتحتفظ المجموعة بمحتواها للحسابات الأخرى.

- لا يزيل ccam الرابط إلا إذا كان لا يزال يشير إلى داخل مجلد المجموعات كما كان عند إنشاء الخطة. والمسار الذي تغيّر منذ ذلك الحين، مثل ملف حقيقي حيث كان الرابط، يُترك كما هو ويتوقف التغيير.
- مع `--fresh`، يُتخطّى النسخ.
- الحساب الذي ينقل عنصرًا إلى مجموعة لا يزال مجلدها فارغًا يجلب معه نسخة مما كان يراه، فيحتفظ بملفاته. وتُجهَّز تلك النسخة في `<group>/.ccam-incomplete-<name>` حتى تكتمل.
- الحساب الذي يتوقف عن مشاركة MCP يحتفظ بخوادم المجموعة: تذهب إلى ملف `.claude.json` الخاص به (مع الاحتفاظ بنسخة احتياطية)، إلا حيث يكون لدى الحساب بالفعل خادم بالاسم نفسه، فيبقى كما هو. ومع `--fresh` تسرد الخطة بدلًا من ذلك الخوادم التي سيتوقف تحميلها. والانضمام مرة أخرى يزيل النسخ التي لم يغيّرها أحد.
- مغادرة `extensions` توجّه أيضًا قوائم المكوّنات الإضافية للمجموعة إلى مجلد المجموعة حيث كانت مساراتها تمر عبر مجلد الحساب المغادر؛ وبالنسبة إلى الحسابات الأخرى، تقود المسارات المعاد كتابتها إلى المكوّنات الإضافية نفسها.

الملف الواحد الموجود في الجهتين بمحتوى مختلف تعارض:

<!-- motion:share-conflict -->

تختار أنت للملفات المفردة مثل `settings.json`؛ أما داخل المجلدات فتبقى نسخة المجموعة، وتذكر الخطة الملفات المختلفة. وتعتمد الخيارات على الملف؛ والخيار الأول هو الافتراضي:

| الملف | الخيارات |
| --- | --- |
| `history.jsonl` | `merge by line`، `keep group`، `use mine` |
| ملفات `.json` الأخرى | `keep group`، `use mine`، `merge JSON` |
| أي ملف آخر | `keep group`، `use mine` |

ويدمج `merge JSON` الكائنات بشكل متكرر، ويضم المصفوفات دون تكرار، ويحتفظ بقيمة المجموعة عندما تعيّن الجهتان المفتاح نفسه بقيمتين مختلفتين، فلا ترى الحسابات الأخرى أي تغيير. وكل خيار ما عدا `keep group` يحفظ أولًا نسخة المجموعة في المجلد `.ccam-backup/<time>/group/` الخاص بالحساب. ويذهب ملفك أنت دائمًا إلى النسخة الاحتياطية.

عند كل تشغيل، يقارن ccam كل مسار في البنية بمجموعات الحساب باستخدام `lstat` و`readlink` فقط (`share.Check`). ويُنشأ الرابط المفقود في الحال. أما الملف الحقيقي حيث يجب أن يكون رابط، أو الرابط إلى المكان الخطأ، أو الرابط المتبقي، فيوقف التشغيل ويوجّهك إلى `ccam doctor --fix`.

### إعادة التسمية والإزالة والتنسيق القديم

- لا يغيّر `ccam group rename` إلا الإعداد. ويحتفظ المجلد باسمه وتسجّله المجموعة في `dir`، فلا يُنقل شيء ولا يمكن أن يضيع شيء.
- يرفض `ccam group remove` إزالة مجموعة فيها حسابات ما لم تضف `--leave`. وعندها يحتفظ كل حساب بنسخة، ويذهب المجلد إلى `~/.local/share/ccam/trash/` ما لم تستخدمه مجموعة أخرى.
- قبل المجموعات، كان `config.toml` يخزّن جدول `[profile.share]` لكل حساب. ويقرأ ccam ذلك التنسيق بوصفه مجموعات دون نقل أي شيء. وعند أول حفظ يكتب التنسيق الجديد ويحتفظ بالملف القديم باسم `config.toml.bak`. والعناصر التي تشاركها الحسابات نفسها تصبح مجموعة واحدة تحمل اسم المجلد؛ وتحصل البقية على أسماء مثل `default-history` وتواصل استخدام المجلد نفسه.

### على Windows

لا ينشئ Windows الروابط الرمزية إلا في وضع المطور (Developer Mode، في الإعدادات > النظام > للمطورين) أو للمسؤول. ودون هذا الحق، يصبح المجلد المشترك نقطة وصل مجلدات (directory junction)، وهي تعمل بالطريقة نفسها مع Claude Code. أما الملف المشترك (`settings.json` و`CLAUDE.md` و`history.jsonl`) فلا يمكن أن يكون إلا رابطًا رمزيًا، لذا يُرفض تغيير المجموعة الذي يحتاج إلى رابط كهذا قبل أن يُنقل أي شيء. وللسبب نفسه ينشئ أول `ccam add` هناك المجموعة `default` بالعنصرين `mcp` و`extensions` فقط ويخبرك بذلك، ويُرفض `ccam add --group` إلى مجموعة تشارك الإعدادات أو الذاكرة أو السجل.

## الأمان والحالات الخاصة

- **الحسابات التي تعمل.** يُرفض التغيير ما دام Claude Code يعمل بأحد الحسابات التي ينقلها، لأن الجلسة تظل تكتب في تلك الملفات، سواء شغّلها ccam أم لا (`claude` العادي، أو بيئة تطوير، أو الجلسات الخلفية). ويُكرَّر الفحص قبيل الحفظ مباشرة. ولا تخص المغادرة والانضمام إلا الحساب المذكور؛ ومغادرة `extensions` تعيد أيضًا كتابة المسارات في قوائم المكوّنات الإضافية للمجموعة، دون تغيير معناها بالنسبة إلى الآخرين. ويظل `--dry-run` يطبع الخطة.
- **النسخ الاحتياطية.** تحصل كل خطة على مجلد `.ccam-backup/<time>/` خاص بها داخل مجلد الحساب. ولا يحذفها ccam أبدًا.
- **روابطك الخاصة.** إذا كان مسار ما رابطًا أنشأته بنفسك بالفعل، مثلًا إلى مستودع dotfiles، وكانت المجموعة فارغة، ينسخ ccam هدف الرابط إلى المجموعة وينقل الرابط نفسه إلى النسخة الاحتياطية. ولا يُمسّ الهدف.
- **المجلدات المستوردة.** الحساب المسجَّل في مكانه باستخدام `ccam import` (مثل `~/.claude`) يتغيّر في مكانه، فيرى `claude` العادي التغيير أيضًا؛ وتذكر الخطة ذلك. وعندما يغادر مجموعة يحتفظ دائمًا بنسخة، حتى مع `--fresh`، ويُخرجه `ccam remove` من مجموعاته بهذه الطريقة قبل أن ينسى الحساب. وملفاته المشتركة موجودة في مجلد المجموعة ما دام في مجموعة، لذا أزل الحساب (أو غادر المجموعات) قبل أن تحذف `~/.local/share/ccam`.
- **حسابان في مجلد واحد.** الحسابات التي تتشارك السجل تكتب في الملفات نفسها. وعندما يبدأ حساب ثانٍ في مجلد يعمل فيه حساب بالفعل، يحذّرك ccam ويسألك.
- **التعديلات اليدوية.** إذا ذكرت `members` اسمًا ليس حسابًا، يتجاهله ccam ويبلّغ عنه `ccam doctor`. ولا يمكن أن تكون الأسماء `local` و`on` و`off` و`none` أسماء مجموعات.

## صفحات ذات صلة

- [لوحة التحكم](/features/dashboard)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [خوادم MCP المشتركة](/features/mcp)
- [نسخ إعداد موجود](/features/migrate)
- [مرجع الأوامر](/docs/commands)
