---
title: اختيار حساب للسكربتات
description: يطبع ccam pick للسكربت الحساب الذي لديه أكبر هامش، ويعطي العمّال الذين يبدؤون معًا حسابات مختلفة، ويخرج بالرمز 75 عندما لا يتبقى هامش.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: a8baeb64de99
---

يطبع `ccam pick` اسم الحساب المسجَّل دخوله الذي لديه أكبر هامش، ليشغّل به السكربت برنامج Claude Code. ويحصل كل عامل من العمّال الذين يبدؤون في الوقت نفسه على حساب مختلف ما دامت هناك حسابات كافية لديها هامش، فتتوزع دفعة من مهام `claude -p` على حساباتك بدلًا من أن تستنزف حسابًا واحدًا.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise. يقول مركز المساعدة لدى Anthropic: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (يُخصَّص كل مقعد في خطة Team لشخص واحد محدد بالاسم، ولا يجوز أن يستخدمه إلا ذلك الشخص) و"Users can't route their work through another user's account to use that user's remaining usage" (لا يجوز للمستخدمين تمرير عملهم عبر حساب مستخدم آخر للاستفادة من الاستخدام المتبقي لذلك المستخدم) ([خطط Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ وبالمثل في [خطط Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). لذا فإن استخدام الاختيار التلقائي بين المقاعد قد يخالف شروط Anthropic. تحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه. أما الخطط الشخصية (Pro وMax) فتتناولها [ملاحظة الشروط في الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## الاستخدام

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

رموز الخروج: `0` مع اسم على المخرج القياسي، و`75` (`EX_TEMPFAIL`) عندما لا يكون لأي حساب هامش أو لا يكون أي حساب مسجَّل الدخول، و`1` للأخطاء الأخرى مثل خيار غير صحيح.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--json` | يطبع الحساب واستخدامه في نافذتي الخمس ساعات والأسبوع ومصدر الأرقام وسبب اختياره | معطّل |
| `--model <model>` | يحتسب أيضًا النافذة الأسبوعية لذلك النموذج، مثل `opus` | لا شيء |
| `--max-usage <1-100>` | يفضّل الحسابات التي يقل استخدامها عن هذه القيمة | `auto_max_usage`، وإلا 100 |
| `--strategy headroom\|expiring\|balance` | أدنى استخدام في نافذة الخمس ساعات، أو أقرب إعادة ضبط لنافذة الخمس ساعات، أو الحصة الأسبوعية التي ستضيع أولًا دون استخدام | `auto_strategy`، وإلا `headroom` |
| `--cross-org` | يختار أيضًا حسابات المؤسسات الأخرى | `cross_org`، وإلا معطّل |
| `--wait[=<duration>]` | عندما لا يكون لأي حساب هامش، ينتظر حتى يُعاد ضبط أولها، بحد أقصى هذه المدة. وتأتي القيمة بعد `=`: فيُرفض `--wait 30m` مع التلميح `--wait takes its value after '=': --wait=30m` | معطّل؛ و`5h` عند تحديده وحده |
| `--dry-run` | يختار دون تسجيل الاختيار | معطّل |

## كيف يعمل

1. يسرد ccam الحسابات المسجَّل دخولها (حسابات مفاتيح API وتسجيلات دخول Console ليس لها حدود فتُستبعد) في مؤسسة الحساب الذي سيشغّله هذا المجلد (راجع [المؤسسات](/features/next#auto))؛ ويُستبعد أيضًا المجلد المسجَّل دخوله بحساب Claude نفسه الذي يستخدمه ذلك الحساب.
2. يأخذ قفلًا (`~/.local/state/ccam/picks.lock`)، فتجري عمليات الاختيار المتزامنة واحدة بعد أخرى.
3. يستعلم عن استخدام كل الحسابات في آن واحد، بالطريقة نفسها التي يتبعها `ccam usage`، مع الذاكرة المؤقتة لمدة 60 ثانية. والنافذة التي أبلغ Claude Code بامتلائها تُحتسب 100% حتى إعادة ضبطها ([سجل الحد](/features/usage#the-limit-record))، فيبقى الحساب الذي بلغ حده للتو مستبعدًا حتى عندما تُظهر أرقامه المحفوظة هامشًا.
4. تُستبعد أولًا الحسابات التي اختيرت خلال آخر دقيقتين (`~/.local/state/ccam/picks/<name>`؛ ولا يُحتسب الاختيار المؤرخ بعد الوقت الحالي للساعة)، وتختار [قواعد الاختيار](/features/next#auto) الخاصة بـ `ccam next --auto` من بين البقية: ما يقل عن `auto_max_usage` أولًا، وفق الاستراتيجية، ولا تختار أبدًا حسابًا بلغ 100%.
5. إذا كانت كل الحسابات التي لديها هامش قد اختيرت مؤخرًا، تختار القواعد من بينها كلها من جديد، فيتشارك العامل الرابع على ثلاثة حسابات حسابًا مع غيره.
6. يُسجَّل الاختيار ويُحرَّر القفل، ثم يُطبع الاسم.

تهم نافذة الدقيقتين لأن استخدام المهمة الجديدة يحتاج إلى بعض الوقت حتى يظهر في الأرقام: ودونها، سيرى خمسة عمّال بدؤوا معًا الحساب نفسه على أنه صاحب أكبر هامش.

## الأمان والحالات الخاصة

- لا يفعل `ccam pick` إلا قراءة الاستخدام وكتابة وقت الاختيار؛ ولا يجدّد الرموز أبدًا ولا يقرؤها أكثر مما يفعل `ccam usage`.
- الاختيار لا يحجز الحساب. فقد يبدأ سكربتان على الحساب نفسه بعد دقيقتين؛ استخدم `--fallback` للتعافي عندما يبلغ ذلك الحساب حده.
- لا يحتفظ الانتظار مع `--wait` بأي قفل، فتستمر عمليات الاختيار الأخرى في أثنائه.
- يتجاهل `ccam pick` الأمر [`ccam folder`](/features/folder): فهو يختار حسب الاستخدام فقط.

## صفحات ذات صلة

- [التبديل إلى الحساب التالي](/features/next)
- [إعادة التشغيل عند بلوغ الحد](/features/next#run-again-at-the-limit)
- [حدود الاستخدام](/features/usage)
- [حسابات مفاتيح API](/features/api)
