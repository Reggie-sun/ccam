---
title: نسخ إعداد موجود
description: انسخ الإعدادات وCLAUDE.md والإضافات وخوادم MCP من ~/.claude موجود إلى مجموعة في ccam باستخدام ccam migrate.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

إذا كنت قد استخدمت Claude Code قبل ccam، فإعدادك موجود في `~/.claude` و`~/.claude.json`. ينسخه `ccam migrate` إلى مجموعة مشاركة، فيبدأ كل حساب يشارك هذه العناصر عبر المجموعة بالإعداد الذي لديك بالفعل. يُقرأ المجلد المصدر فقط، ويستمر `claude` العادي في العمل كما كان. شغّله مرة واحدة بعد أول `ccam add`؛ وعندما يحتوي `~/.claude.json` على خوادم MCP، يذكرها أول `ccam add` في نصيحته.

## الاستخدام

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

في الطرفية، ودون خيارات العناصر ودون `--yes`، يعرض ccam أولًا قائمة اختيار:

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

يبدّل مفتاح المسافة (أو `x`) العنصر الذي تحت المؤشر، وتنقل مفاتيح الأسهم أو `j`/`k` التحديد، ويعرض Enter الخطة، ويطبع Esc أو `q` أو Ctrl+C الرسالة `Cancelled.`. والعناصر التي لا يحتوي المصدر على شيء لها تُعلَّم بـ `not found`. وإذا ألغيت تحديد كل شيء، يطبع ccam الرسالة `Nothing selected.` ويتوقف.

ثم يطبع ccam الخطة ويسألك قبل أن ينسخ:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

لا تؤثر المجموعة إلا في حساباتها، وفي العناصر التي تشاركها فقط، لذلك يسرد الملخص حسابات المجموعة. وعندما لا تشارك المجموعة أحد العناصر المنسوخة، أو لا تكون بعض الحسابات فيها، يطبع ccam الأمرين `ccam group set <group> --add <items>` و`ccam group join <group> <accounts>` اللذين يغيّران ذلك.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `[folder]` | مجلد Claude Code الذي يُنسخ منه | `~/.claude` |
| `--settings` | ينسخ `settings.json` | معطّل |
| `--memory` | ينسخ ملف `CLAUDE.md` العام | معطّل |
| `--extensions` | ينسخ `agents/` و`commands/` و`skills/` و`plugins/` و`output-styles/` و`rules/` | معطّل |
| `--mcp` | ينسخ خوادم MCP الخاصة بالمستخدم (`mcpServers` في `.claude.json`) | معطّل |
| `--history` | ينسخ سجل المحادثات (`projects/` و`history.jsonl`) | معطّل |
| `-g`، `--group <group>` | مجموعة المشاركة التي يُنسخ إليها | `default` |
| `-y`، `--yes` | ينسخ دون سؤال؛ وعند التعارض يُحتفظ بنسخة المجموعة | معطّل |
| `--dry-run` (عام) | يطبع الخطة ويتوقف قبل أي سؤال أو نسخ | معطّل |

في الطرفية، إذا لم تحدد أي خيارات للعناصر، تختار العناصر في قائمة الاختيار. وخارج الطرفية، أو مع `--yes`، ينسخ ccam كل شيء ما عدا سجل المحادثات. ويُستثنى السجل افتراضيًا لأن `projects/` قد يكون كبيرًا.

القيمة `--group on` تعني `default`. وتُرفض `--group off` و`--group local` مع `--group takes a group name, not "off"`. يدعم `ccam migrate` الخيار `--dry-run`، وتظل قائمة الاختيار تظهر في التشغيل التجريبي لتتمكن من معاينة اختيار ما.

## كيف يعمل

يبني ccam خطة كاملة قبل أن ينسخ أي شيء. فيختار العناصر (`internal/cli/migrate.go`)، ثم يقرأ المصدر والمجموعة ويبني الخطة (`migrate.Build`). ولا يُكتب شيء حتى `Plan.Apply`، بعد أن تؤكد.

1. **إيجاد المصدر.** يرفض ccam المسار المفقود، أو الذي ليس مجلدًا، أو الذي لا يحتوي على أي مدخل معروف من Claude Code (`checkSource`). ولخوادم MCP وحدها، لا بأس بأن يكون `~/.claude` فارغًا أو مفقودًا ما دام `~/.claude.json` يحتوي على خوادم. ويجد ccam ملف `.claude.json` الذي يخص المجلد (`migrate.ClaudeJSON`): داخله، أو بالنسبة إلى `~/.claude` الذي ليس فيه واحد، `~/.claude.json` بجانبه. فهناك يحتفظ به Claude Code عندما لا يكون `CLAUDE_CONFIG_DIR` معيّنًا.
2. **ربط العناصر بالمسارات.** يأتي كل مسار من `share.Layout`، وهي قائمة السماح نفسها التي تستخدمها المجموعات. وتُقرأ خوادم MCP من `.claude.json` باستخدام `mcp.LocalServers`، الذي يستخدم الماسح `profile.TopLevelValue` بدلًا من تحليل الملف كله.
3. **النسخ لا النقل.** تذهب الملفات إلى `~/.local/share/ccam/groups/<group>/` باستخدام `fsutil.AtomicWrite` أو `fsutil.CopyFile`، الذي يكتب ملفًا مؤقتًا بجانب الهدف ثم يعيد تسميته، فلا يترك النسخ المقطوع ملفًا جزئيًا أبدًا. وتذهب خوادم MCP إلى ملف `mcp.json` الخاص بالمجموعة تحت قفله: فالخادم الذي تضيفه عملية ccam أخرى في الأثناء يُحتفظ به، وتُحدَّث النسخ المنشأة لـ [Remote Control](/guides/remote-control) بعد ذلك.

يربط كل عنصر مسارًا واحدًا في المصدر بمسار واحد في المجموعة:

<!-- motion:migrate-items -->

صفّا السجل متقطعان لأنهما لا يُنسخان إلا مع `--history`. ويُدرج `.credentials.json` و`.claude.json` في `share.Never` ولا يصلان إلى المجموعة أبدًا: يتوقف ccam بخطأ داخلي إذا ظهر أحدهما يومًا في البنية (`Build`)، ويتخطاهما المرور على المجلد بالاسم.

يمنح ccam كل مسار حالة، وهي ما يصفه كل سطر في الخطة (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: غير موجود في المصدر، أو ملف حيث يُتوقع مجلد. يُذكر مرة واحدة لكل عنصر بالشكل `not found, skipped`.
- **Linked**: مسار المصدر يشير بالفعل إلى نسخة المجموعة، مثلًا عندما يكون المصدر مجلد حساب في ccam.
- **Same**: لدى المجموعة بالفعل محتوى مطابق.
- **New**: لا تملكه المجموعة، فيُنسخ كاملًا.
- **Merge**: مجلد موجود في الجهتين. لا يضيف ccam إلا الملفات التي لا تملكها المجموعة (`copyMissing`)؛ والملفات المختلفة تبقى كما هي في المجموعة، وتقول الخطة `group copies kept`. وبالنسبة إلى `plugins/`، يُدمج السجلان `installed_plugins.json` و`known_marketplaces.json`، وتُعاد كتابة مسارات الإضافات المطلقة لتشير إلى مجلد المجموعة.
- **Conflict**: ملف موجود في الجهتين بمحتوى مختلف. تختار لكل ملف: `keep group` أو `use mine`، إضافة إلى `merge JSON` لملفات `.json` (دمج عميق تتقدم فيه المجموعة عند تعارض القيم) أو `merge by line` لـ `history.jsonl`. وخادم MCP الذي يختلف تعريفه تعارض أيضًا، مع `keep group` أو `use mine`. والخيار الأول المذكور هو الافتراضي عندما تضغط Enter؛ ويحتفظ `--yes` دائمًا بنسخة المجموعة.

قائمة الاختيار نفسها هي `tui.MigratePicker` المبنية على المكوّن `checklist` الذي يستخدمه `ccam add`؛ ويعلّم ccam العناصر المفقودة باستدعاءات lstat وقائمة خوادم MCP (`migrate.Found`).

### تسجيلات الدخول إلى MCP

خوادم MCP التي تسجّل الدخول عبر OAuth تحتفظ برموزها في ملف تسجيل الدخول، `.credentials.json` تحت `mcpOAuth`، لا في قائمة الخوادم. ولا ينسخها ccam. فتسجيل الدخول المنسوخ يحمل رمز التحديث نفسه الذي يحمله الأصل، ويبدّل الخادم ذلك الرمز عند أول تجديد: فيُسجَّل خروج كل حامل آخر من ذلك الخادم، بما في ذلك `~/.claude` المصدر. لذلك، لكل خادم منسوخ سجّل المصدر الدخول إليه، يسرد سطر `sign-ins` في الخطة حسابات المجموعة التي يجب أن تسجّل الدخول بنفسها باستخدام `/mcp`، ويكرر الملخص ذلك. والمدخل الذي لا رمز فيه ليس تسجيل دخول ولا يُسرد، ولا يُسرد الحساب الذي لديه تسجيل دخوله الخاص بالفعل، ولا الخادم الذي تحتفظ المجموعة بتعريفه. ولا يُكتب أي ملف تسجيل دخول.

## الأمان والحالات الخاصة

- يُقرأ المصدر فقط. لا يُنقل أي ملف فيه ولا يُربط ولا يُغيَّر، ولا يُكتب تسجيل الدخول في المصدر أبدًا، بما في ذلك على macOS حيث قد يُقرأ من سلسلة المفاتيح.
- قبل استبدال ملف في المجموعة، سواء بـ `use mine` أو بدمج أو بخوادم MCP جديدة في ملف `mcp.json` موجود، تُحفظ نسخة المجموعة في `<group>/.ccam-backup/<time>/`. واسم المجلد فريد لكل تشغيل، فلا يكتب تشغيلان في الثانية نفسها فوق النسخ الاحتياطية لبعضهما. ويطبع الملخص المسار.
- مع `--yes`، يحتفظ كل تعارض بنسخة المجموعة، ويخبرك ccam كيف تختار بدلًا من ذلك: شغّل الأمر نفسه دون `--yes`.
- إذا لم يكن هناك ما سيتغيّر، يطبع ccam الرسالة `Nothing to copy.`؛ وإذا احتفظ كل تعارض بنسخة المجموعة، يطبع `Nothing to copy: the group keeps its own copies.`
- إذا فشل النسخ في منتصف الطريق، يتوقف ccam مع `The source folder was not changed. Group files replaced so far are saved in` متبوعًا بمجلد النسخ الاحتياطية. وتشغيله مرة أخرى ينسخ الملفات التي لم تُنسخ كاملة.
- عندما تشارك عدة مجموعات MCP، تكون الأوامر التي يقترحها `ccam doctor` و`ccam ext` و`ccam mcp list` بالشكل `ccam migrate --mcp --group <a|b>`.
- ينتهي المصدر المفقود بالرسالة `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## صفحات ذات صلة

- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [الفحص والإصلاح](/features/doctor)
- [استخدام عدة حسابات Claude Code على حاسوب واحد](/guides/multiple-claude-code-accounts)
- [شرح CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
