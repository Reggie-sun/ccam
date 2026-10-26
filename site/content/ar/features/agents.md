---
title: الجلسات الخلفية لكل الحسابات
description: يعرض ccam agents وكلاء Claude Code الخلفيين والجلسات المفتوحة لكل حساباتك في جدول واحد، ويفتح ccam attach إحداها بمعرّفها.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: 06281fb13119
---

يحتفظ Claude Code بالجلسات الخلفية (`claude --bg`، وعرض الوكلاء) لكل حساب على حدة: فلا يعرض `claude agents` إلا جلسات الحساب الذي يعمل به. أما `ccam agents` فيسأل كل الحسابات في آن واحد ويسرد جلساتها معًا، ويفتح `ccam attach <id|name>` الجلسة الخلفية بالحساب الذي بدأها، أيًّا كان.

## الاستخدام

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

يبيّن `State` حالة الجلسة الخلفية (`working` أو `blocked` أو `done` أو `failed` أو `stopped`) وما تفعله عمليتها (`busy` أو `waiting` أو `idle`)؛ وتظهر الجلسة التفاعلية المفتوحة بالقيمة `open`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--all` | يسرد أيضًا الجلسات الخلفية التي انتهت أو فشلت أو أُوقفت | معطّل |
| `--json` | يطبع الجلسات بصيغة JSON، كلٌّ منها مع `account` الذي تنتمي إليه | معطّل |

يمرّر `ccam attach <id|name> [claude args]` كل ما يلي المعرّف إلى `claude attach` دون تغيير، بما في ذلك الخيارات: فالأمر `ccam attach 49ec --dry-run` يسلّم `--dry-run` إلى Claude Code. ولطباعة الأمر بدلًا من ذلك، ضع `--dry-run` قبل المعرّف (`ccam --dry-run attach 49ec` أو `ccam attach --dry-run 49ec`).

## كيف يعمل

1. لكل حساب، يشغّل ccam الأمر `claude agents --json` (مع `--all` إذا حُدّد) بقيمة `CLAUDE_CONFIG_DIR` الخاصة بالحساب، للحسابات كلها بالتوازي، مع مهلة 15 ثانية لكل منها. وهذا يقرأ قائمة جلسات الحساب نفسه ولا يشغّل الخدمة الخلفية لـ Claude Code.
2. يُضاف اسم الحساب إلى كل مدخل، وتُرتَّب القائمة حسب وقت البدء، الأحدث أولًا. والحساب الذي لا يجيب يُذكر اسمه في الخطأ القياسي ويُستبعد، ثم يخرج ccam بالرمز 1، حتى لا يظن السكربت أن الفشل قائمة فارغة.
3. يبحث `ccam attach` عن `jobs/<id>` في مجلد كل حساب، حيث يحتفظ Claude Code بالجلسة الخلفية، ويشغّل `claude attach <id>` بذلك الحساب، بالطريقة نفسها التي يبدأ بها `ccam <name>` برنامج Claude Code. والكلمة التي ليست بادئة معرّف تُطابَق مع أسماء الجلسات، دون تمييز بين الأحرف الكبيرة والصغيرة؛ والكلمة التي تشبه معرّفًا لا تُعامل أبدًا على أنها اسم.

## الأمان والحالات الخاصة

- معرّفات الجلسات الخلفية تخص حسابًا واحدًا: فإذا شُغّل `claude attach` بحساب آخر لم يجدها، وهذا ما يحلّه `ccam attach`.
- إذا اشتركت جلسات في حسابين في بادئة معرّف أو في كلمة من الاسم، يسردها `ccam attach` ويطلب معرّفًا أو اسمًا أطول.
- يمرّر ccam خوادم MCP المشتركة إلى Claude Code باستخدام `--mcp-config=<file>` في وسيط واحد، وهو ما يقبله `claude agents` و`claude attach`؛ وتواصل الجلسات الخلفية استخدام ذلك الملف عند إعادة تشغيلها.
- تأخذ الخدمة الخلفية لحساب مفتاح API مفتاحها من الصدفة التي بدأتها. لذلك ابدأ أول جلسة خلفية لمثل هذا الحساب عبر ccam (`ccam gw --bg …`)، لتحصل الخدمة على المفتاح.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [اختيار حساب للسكربتات](/features/pick)
- [لوحة التحكم](/features/dashboard)
