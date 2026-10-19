---
title: مشاركة خوادم MCP والإعدادات والسجل
description: كيف يشارك ccam خوادم MCP وsettings.json وCLAUDE.md والوكلاء والمهارات وسجل المحادثات بين حسابات Claude Code مع إبقاء كل تسجيل دخول منفصلًا.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

يحتفظ كل حساب بتسجيل دخوله الخاص. أما كل ما عدا ذلك فيمكن مشاركته عنصرًا عنصرًا: يربط ccam ملف الحساب أو مجلده بمجلد مجموعة مشترك، فيظهر التغيير الذي يُجرى في حساب واحد في كل حسابات تلك المجموعة.

## ما الذي يمكن مشاركته

| العنصر | الخيار | ما يُشارَك |
| --- | --- | --- |
| خوادم MCP | `--mcp` | ملف `mcp.json` الخاص بالمجموعة، ويُمرَّر إلى Claude Code عند التشغيل |
| الإعدادات | `--settings` | `settings.json`: الأذونات والخطافات (hooks) والنموذج وسطر الحالة |
| الذاكرة العامة | `--memory` | `CLAUDE.md` |
| الإضافات | `--extensions` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/`، `rules/` |
| سجل المحادثات | `--history` | `projects/` و`history.jsonl` |

تشارك الحسابات الجديدة كل شيء ما عدا سجل المحادثات.

ملفا تسجيل الدخول `.credentials.json` و`.claude.json` لا يُشارَكان أبدًا. لا يربط ccam إلا مسارات من قائمة ثابتة، وهذان الملفان ليسا فيها.

## المجموعات

القيمة `on` تعني المجموعة `default`. ويمكنك أيضًا تحديد اسم، بحيث لا يشارك عنصرًا ما إلا بعض الحسابات:

```sh
ccam share work --history team
ccam share side --history team
```

أصبح `work` و`side` الآن يتشاركان المحادثات فيما بينهما، ولكن ليس مع حساباتك الأخرى.

## تفعيل المشاركة

عندما يبدأ حساب بمشاركة عنصر ما، فإن ccam:

1. يطبع خطة بكل ملف سينقله أو ينسخه أو يربطه، ويسألك قبل أن يبدأ؛
2. ينقل ملفات الحساب إلى مجلد المجموعة إذا كانت المجموعة لا تزال فارغة، أو يضيف الملفات التي لا تملكها المجموعة بعد؛
3. عندما يختلف الملف نفسه، يسألك هل تريد الاحتفاظ بنسخة المجموعة، أم استخدام نسختك، أم دمج JSON؛
4. يحتفظ بملفك الأصلي في `.ccam-backup/<time>/` داخل مجلد الحساب؛
5. يستبدل الملف الأصلي برابط إلى المجموعة.

تُسمّى ملفات المحادثات باسم معرّف الجلسة، لذلك تُدمج السجلات القادمة من عدة حسابات دون تعارضات.

## إيقاف المشاركة

```sh
ccam share work --history off
```

يُزال الرابط ويحصل الحساب على نسخته الخاصة من المحتوى الحالي للمجموعة. أضف `--fresh` لتبدأ فارغًا بدلًا من ذلك. وتحتفظ المجموعة بمحتواها للحسابات الأخرى.

## نسخ إعداد موجود

إذا كنت قد استخدمت Claude Code قبل ccam، فإعداداتك موجودة في `~/.claude`. ينسخها `ccam migrate` إلى مجموعة، فتحصل عليها كل الحسابات التي تشارك هذه العناصر:

```sh
ccam migrate
```

ما الذي ينسخه:

- `settings.json` وملف `CLAUDE.md` العام؛
- الإضافات `agents/` و`commands/` و`skills/` و`plugins/` و`output-styles/` و`rules/`؛
- خوادم MCP الخاصة بالمستخدم، أي المدخل `mcpServers` في `~/.claude.json`، إلى ملف `mcp.json` الخاص بالمجموعة؛
- سجل المحادثات (`projects/` و`history.jsonl`)، ولكن فقط عندما تضيف `--history`.

يُقرأ المجلد المصدر فقط: لا يُنقل أي شيء فيه ولا يُربط ولا يُغيَّر، ويستمر `claude` العادي في العمل كما كان. ولا تُنسخ ملفات تسجيل الدخول `.credentials.json` و`.claude.json` أبدًا.

يطبع ccam الخطة ويسألك قبل أن ينسخ أي شيء:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

في النهاية، يسرد ccam الحسابات التي تشارك كل عنصر عبر المجموعة. وإذا كانت هناك حسابات لا تشاركه، يطبع أوامر `ccam share` التي تفعّل المشاركة لها.

عندما تملك المجموعة نسخة مختلفة من ملف ما، يسألك ccam كما يفعل عند تفعيل المشاركة: الاحتفاظ بنسخة المجموعة أم استخدام نسختك. ويمكن أيضًا دمج `settings.json` وملفات JSON الأخرى، ودمج `history.jsonl` سطرًا بسطر. تُحفظ النسخ السابقة للمجموعة في `.ccam-backup/<time>/` داخل مجلد المجموعة.

| الخيار | التأثير |
| --- | --- |
| `--settings`، `--memory`، `--extensions`، `--mcp`، `--history` | ينسخ العناصر التي تحددها فقط؛ ودونها ينسخ كل شيء ما عدا سجل المحادثات |
| `-g`، `--group <group>` | ينسخ إلى هذه المجموعة بدلًا من `default` |
| `-y`، `--yes` | ينسخ دون سؤال؛ وعندما يختلف ملف، يحتفظ بنسخة المجموعة |

ضع `--dry-run` قبل الأمر، كما في `ccam --dry-run migrate`، لترى الخطة دون نسخ أي شيء. وللنسخ من مجلد آخر، حدّده: `ccam migrate ~/.claude-work --group work`.

## خوادم MCP

يخزّن Claude Code خوادم MCP في `.claude.json` مع تسجيل الدخول، لذلك لا يمكن مشاركة هذا الملف. يحتفظ ccam بالخوادم المشتركة في ملف `mcp.json` الخاص بالمجموعة، ويشغّل Claude Code مع `--mcp-config=<group>/mcp.json`. وتظل خيارات `--mcp-config` الخاصة بك سارية فوق ذلك.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

الخادم الذي تضيفه داخل جلسة باستخدام `claude mcp add` يُحفظ لذلك الحساب فقط. في المرة التالية التي تشغّل فيها ذلك الحساب، يذكر ccam اسم الخادم والأمر الذي ينقله إلى المجموعة:

```sh
ccam mcp import work gh
```

خوادم MCP التي تسجّل الدخول عبر OAuth تخزّن رموزها مع تسجيل دخول الحساب، لذا امنحها التفويض مرة واحدة في كل حساب.

إذا كنت تفضّل أن تظهر الخوادم في `claude mcp list` و`/mcp` كخوادم مستخدم عادية، فعيّن `mcp_strategy = "sync"` في `~/.config/ccam/config.toml`. عندها يكتبها ccam في ملف `.claude.json` الخاص بكل حساب قبل التشغيل، مستخدمًا قفل الملف الخاص بـ Claude Code نفسه ومع نسخة احتياطية.

## حسابان في المجلد نفسه

الحسابات التي تتشارك سجل المحادثات تكتب في الملفات نفسها. إذا شغّلت نسخة ثانية من Claude Code في مجلد يعمل فيه بالفعل حساب آخر من المجموعة نفسها، يحذّرك ccam ويسألك هل تريد المتابعة.

## عندما يتعطل شيء ما

```sh
ccam doctor
ccam doctor --fix
```

إذا تحوّل مسار مشترك إلى ملف عادي، مثلًا بعد أن استعدت نسخة احتياطية يدويًا، يتوقف ccam قبل تشغيل Claude Code ويوجّهك إلى `ccam doctor --fix`، الذي يدمج الملف في المجموعة ويحتفظ بنسخة احتياطية.
