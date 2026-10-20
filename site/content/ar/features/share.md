---
title: المشاركة عبر المجموعات
description: تشارك مجموعة ccam عناصر مختارة بين حساباتها عبر مجلد واحد؛ وينقل ccam group البيانات وفق خطة، مع نسخ احتياطية وخيارات للتعارضات.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 0c0338362ba2
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
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
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
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

التغييرات نفسها متاحة في تبويب **Groups** في [لوحة التحكم](/features/dashboard)، ويختار `s` على أي حساب مجموعاتِه.

## العناصر

| العنصر | ما يغطيه |
| --- | --- |
| `mcp` | خوادم MCP، عبر ملف `mcp.json` الخاص بالمجموعة ([خوادم MCP المشتركة](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | ملف `CLAUDE.md` العام |
| `extensions` | `agents/` و`commands/` و`skills/` و`plugins/` و`output-styles/` و`rules/` |
| `history` | `projects/` و`history.jsonl`، فيستطيع حساب آخر استخدام `--continue` |

العنصر الذي لا تشاركه أي مجموعة من مجموعات الحساب يبقى في مجلد الحساب نفسه.

## الأوامر والخيارات

| الأمر | ما يفعله |
| --- | --- |
| `ccam group` أو `ccam group list` | يسرد المجموعات |
| `ccam group show <group>` | ما تشاركه المجموعة، وحساباتها، ومجلدها، وخوادم MCP فيها |
| `ccam group add <group> --items <items> [--members <accounts>]` | ينشئ مجموعة؛ وتنتقل إليها عناصر حساباتها |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | يغيّر ما تشاركه المجموعة، لكل حساب فيها |
| `ccam group join <group> <account>...` | يضع حسابات في المجموعة |
| `ccam group leave <group> <account>...` | يُخرج حسابات منها؛ ويحتفظ كل منها بنسخة |
| `ccam group rename <group> <new-name>` | يعيد التسمية في الإعداد فقط |
| `ccam group remove <group> [--leave]` | يزيل مجموعة |

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--from <account>` | مع `add` و`set`: الحساب الذي تدخل نسخته إلى المجموعة أولًا | الحساب الافتراضي، ثم حسب الرقم |
| `--replace` | مع `join`: يغادر أولًا المجموعات التي تشارك العناصر نفسها | معطّل |
| `--fresh` | الحسابات التي تتوقف عن مشاركة عنصر تبدأ دونه بدلًا من البدء بنسخة منه | معطّل |
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
2. يقارن `internal/group` القديم بالجديد، ولكل حساب تتغيّر عناصره يبني خطة باستخدام `share.Build`. وتُبنى الخطط واحدة بعد أخرى، وترى كل منها ما ستضعه الخطط السابقة في مجلدات المجموعات، فلا ينقل نسخته إليها إلا الحساب الأول. وتأتي النسخة الأولى من مجلد المجموعة إن كانت فيه نسخة بالفعل، ثم من `--from`، ثم من الحساب الافتراضي، ثم من البقية حسب الرقم.
3. يرفض ccam إذا كان أحد تلك الحسابات يعمل، وعلى Windows دون وضع المطور إذا كان ملف ما سيحتاج إلى رابط.
4. تُحسم التعارضات، بالسؤال أو بـ `--prefer`، وتُطبع الخطة، ثم يسألك ccam.
5. يُحفظ الإعداد الجديد قبل أن يُنقل أي شيء. وإذا توقفت خطوة في منتصف الطريق، يبني `ccam doctor --fix` وإعادة التشغيل خطتهما نحو المجموعات المحفوظة ويُكملان التغيير.
6. تُنفَّذ الخطط بالترتيب. كل خطوة تنقل البيانات إلى المجموعة أو إلى مجلد النسخ الاحتياطية؛ لا يُحذف شيء.

`share.Layout` قائمة سماح بالمسارات التي يغطيها كل عنصر. وأي شيء ليس فيها لا يمكن أن يصل إلى مجموعة أبدًا. يذكر `share.Never` الملفين `.credentials.json` و`.claude.json`، ويتوقف `Build` مع خطأ داخلي إذا ظهر أي منهما يومًا في البنية. وليس للعنصر `mcp` مسارات: فقائمة الخوادم موجودة داخل `.claude.json`، لذلك تعمل [مشاركة MCP](/features/mcp) عبر ملف `mcp.json` الخاص بالمجموعة.

هذا ما يحدث لمجلد واحد، `projects/`، عندما يبدأ `side` مشاركة السجل وفي المجموعة محتوى بالفعل:

<!-- motion:share-on -->

إذا لم يكن لدى المجموعة شيء في ذلك المسار بعد، فالحساب هو أول من يشاركه، ويُنقل مجلده إلى هناك كما هو. وإلا يُنقل كل ملف لا تملكه المجموعة إلى المجموعة، وتبقى الملفات الموجودة في الجهتين في مجلد الحساب، فيُحتفظ بنسخة المجموعة. ثم يذهب مجلد الحساب إلى `.ccam-backup/<time>/` داخل مجلد الحساب، ويحل رابط محله. تُسمّى ملفات المحادثات باسم معرّف الجلسة، لذلك تُدمج السجلات القادمة من عدة حسابات دون تعارضات. وتخزّن سجلات `plugins/` مسارات مطلقة تمر عبر مجلد الحساب؛ فيوجّهها ccam إلى مجلد المجموعة، ويدمج مدخلات الحساب في السجل الذي تملكه المجموعة بالفعل.

مغادرة المجموعة تعكس الرابط:

<!-- motion:share-off -->

لا يزيل ccam الرابط إلا إذا كان يشير إلى داخل مجلد المجموعات، ثم ينسخ المحتوى الحالي للمجموعة إلى الحساب. ومع `--fresh` يُتخطّى النسخ. وتحتفظ المجموعة بمحتواها للحسابات الأخرى. والحساب الذي ينقل عنصرًا إلى مجموعة لا يزال مجلدها فارغًا يجلب معه نسخة مما كان يراه، فيحتفظ بملفاته؛ وتنتقل خوادم MCP بالطريقة نفسها. أما الحساب الذي يتوقف عن مشاركة MCP فيفقد الخوادم التي نسخها ccam إلى ملف `.claude.json` الخاص به (مع الاحتفاظ بنسخة احتياطية)، لأنها جاءت من المجموعة.

الملف الواحد الموجود في الجهتين بمحتوى مختلف تعارض:

<!-- motion:share-conflict -->

يعرض `history.jsonl` الخيارات `merge by line` و`keep group` و`use mine`، مع الدمج أولًا. وتعرض ملفات `.json` الأخرى `keep group` و`use mine` و`merge JSON`. وأي ملف آخر يعرض `keep group` و`use mine`. ويدمج `merge JSON` الكائنات بشكل متكرر، ويضم المصفوفات دون تكرار، ويحتفظ بقيمة المجموعة عندما تعيّن الجهتان المفتاح نفسه بقيمتين مختلفتين، فلا ترى الحسابات الأخرى أي تغيير. وكل خيار ما عدا `keep group` يحفظ أولًا نسخة المجموعة في المجلد `.ccam-backup/<time>/group/` الخاص بالحساب. ويذهب ملفك أنت دائمًا إلى النسخة الاحتياطية.

عند كل تشغيل، يقارن `share.Check` كل مسار في البنية بمجموعات الحساب باستخدام `lstat` و`readlink` فقط. ويُنشأ الرابط المفقود في الحال. أما الملف الحقيقي حيث يجب أن يكون رابط، أو الرابط إلى المكان الخطأ، أو الرابط المتبقي، فيوقف التشغيل ويوجّهك إلى `ccam doctor --fix`.

### إعادة التسمية والإزالة والتنسيق القديم

- لا يغيّر `ccam group rename` إلا الإعداد. ويحتفظ المجلد باسمه وتسجّله المجموعة في `dir`، فلا يُنقل شيء ولا يمكن أن يضيع شيء.
- يرفض `ccam group remove` إزالة مجموعة فيها حسابات ما لم تضف `--leave`. وعندها يحتفظ كل حساب بنسخة، ويذهب المجلد إلى `~/.local/share/ccam/trash/` ما لم تستخدمه مجموعة أخرى.
- قبل المجموعات، كان `config.toml` يخزّن جدول `[profile.share]` لكل حساب. ويقرأ ccam ذلك التنسيق بوصفه مجموعات دون نقل أي شيء. وعند أول حفظ يكتب التنسيق الجديد ويحتفظ بالملف القديم باسم `config.toml.bak`. والعناصر التي تشاركها الحسابات نفسها تصبح مجموعة واحدة تحمل اسم المجلد؛ وتحصل البقية على أسماء مثل `default-history` وتواصل استخدام المجلد نفسه.

### على Windows

لا ينشئ Windows الروابط الرمزية إلا في وضع المطور (Developer Mode، في الإعدادات > النظام > للمطورين) أو للمسؤول. ودون هذا الحق، يصبح المجلد المشترك نقطة وصل مجلدات (directory junction)، وهي تعمل بالطريقة نفسها مع Claude Code. أما الملف المشترك (`settings.json` و`CLAUDE.md` و`history.jsonl`) فلا يمكن أن يكون إلا رابطًا رمزيًا، لذا يُرفض تغيير المجموعة الذي يحتاج إلى رابط كهذا قبل أن يُنقل أي شيء.

## الأمان والحالات الخاصة

- **الحسابات التي تعمل.** يُرفض التغيير ما دام Claude Code يعمل بأحد الحسابات التي ينقلها، لأن الجلسة تظل تكتب في تلك الملفات. ولا تخص المغادرة والانضمام إلا الحساب المذكور. ويظل `--dry-run` يطبع الخطة.
- **النسخ الاحتياطية.** تحصل كل خطة على مجلد `.ccam-backup/<time>/` خاص بها داخل مجلد الحساب. ولا يحذفها ccam أبدًا.
- **روابطك الخاصة.** إذا كان مسار ما رابطًا أنشأته بنفسك بالفعل، مثلًا إلى مستودع dotfiles، وكانت المجموعة فارغة، ينسخ ccam هدف الرابط إلى المجموعة وينقل الرابط نفسه إلى النسخة الاحتياطية. ولا يُمسّ الهدف.
- **المجلدات المستوردة.** الحساب المسجَّل في مكانه باستخدام `ccam import` (مثل `~/.claude`) يتغيّر في مكانه، فيرى `claude` العادي التغيير أيضًا؛ وتذكر الخطة ذلك.
- **حسابان في مجلد واحد.** الحسابات التي تتشارك السجل تكتب في الملفات نفسها. وعندما يبدأ حساب ثانٍ في مجلد يعمل فيه حساب بالفعل، يحذّرك ccam ويسألك.
- **التعديلات اليدوية.** العضو الذي ليس حسابًا يُتجاهَل ويبلّغ عنه `ccam doctor`. ولا يمكن أن تكون الأسماء `local` و`on` و`off` و`none` أسماء مجموعات.

## صفحات ذات صلة

- [لوحة التحكم](/features/dashboard)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [خوادم MCP المشتركة](/features/mcp)
- [نسخ إعداد موجود](/features/migrate)
- [كل الأوامر](/docs/commands)
