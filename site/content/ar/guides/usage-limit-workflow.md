---
title: ماذا تفعل عندما يبلغ Claude Code حد الاستخدام
description: انتظر إعادة التعيين أو تابع المحادثة نفسها فورًا على حساب Claude Code آخر. كيف يعمل حد الخمس ساعات والحد الأسبوعي، خطوة بخطوة.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

عندما يبلغ حساب حد الاستخدام (ويُسمّى غالبًا rate limit)، أمامك خياران: أن تنتظر حتى يُعاد تعيين الحد، أو أن تتابع المحادثة نفسها الآن على حساب آخر تملكه. مع ccam يتطلب الخيار الثاني أمرًا واحدًا، `ccam next --auto --continue`.

## ما الحدود التي يفرضها Claude Code؟

في خطتَي Pro وMax، لدى Claude Code [حد جلسة مدته خمس ساعات وحد أسبوعي](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). يُشارَك الحدّان مع Claude على الويب وفي التطبيقات، ويُحتسب الاستخدام من الاثنين في الوقت نفسه. عندما ينفد أحدهما، [يعرض Claude Code رسالة](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) مثل "You've hit your session limit · resets 3:45pm". لا يفيد تغيير النموذج، لأن الحدود تشمل كل النماذج.

## الخيار 1: انتظر إعادة التعيين

إذا كان لديك حساب واحد، فانتظر. الإصدارات الحديثة من Claude Code [تنتظر في الجلسة المفتوحة](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) وتتابع المهمة من تلقاء نفسها بعد إعادة تعيين الحد. ويعرض `/usage` داخل Claude Code مقدار ما استخدمته من حدود خطتك.

## الخيار 2: تابع على حساب آخر الآن

إذا كان لديك حساب ثانٍ، يمكنك المتابعة فورًا بالمحادثة نفسها. أعِدّ ذلك مرة واحدة:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

يعرض شريط الحالة الآن الحساب واستخدامه خلال 5 ساعات أثناء عملك، فلا يفاجئك الحد.

عند بلوغ الحد، اخرج من Claude Code (`/exit`)، ثم:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

افتح المحادثة في الحساب الذي لديه أكبر هامش متبقٍّ:

```sh
ccam next --auto --continue
```

يبدأ Claude Code بالحساب `work` مع أحدث محادثة في المجلد الحالي، بما فيها كل ما قيل حتى الآن. ومشاركة السجل (`ccam group set default --add history` أعلاه) تُبقيها محادثة واحدة؛ ودونها ينسخ Claude Code المحادثة إلى `work`، ويخبرك ccam بعدد الرموز (tokens) التي سيقرؤها ذلك الرد الأول من جديد.

ومع [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) في Claude Code، تعرض الجلسة التي تتوقف عند حدها إشعارًا وتسجّل النافذة الممتلئة، فلا يختار `ccam next --auto` ذلك الحساب مجددًا قبل إعادة تعيينها، ويدوّن `/ccam:switch` (أو `/ccam:switch work`) المحادثة لـ `ccam next --continue` دون أن يترك شيئًا للتخمين.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise. يقول مركز المساعدة لدى Anthropic: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (يُخصَّص كل مقعد في خطة Team لشخص واحد محدد بالاسم، ولا يجوز أن يستخدمه إلا ذلك الشخص) و"Users can't route their work through another user's account to use that user's remaining usage" (لا يجوز للمستخدمين تمرير عملهم عبر حساب مستخدم آخر للاستفادة من الاستخدام المتبقي لذلك المستخدم) ([خطط Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ وبالمثل في [خطط Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). لذا فإن استخدام الاختيار التلقائي بين المقاعد قد يخالف شروط Anthropic. تحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه. أما الخطط الشخصية (Pro وMax) فتتناولها [ملاحظة الشروط في الأسئلة الشائعة](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## هل يبدّل ccam الحسابات تلقائيًا؟

لا. لا يعمل ccam أثناء عملك: فهو يشغّل Claude Code ويستبدل نفسه به. أنت من يقرر متى تبدّل، ويختار `ccam next --auto` الحساب الأقل استخدامًا خلال 5 ساعات ويتخطى الحسابات التي بلغت 100% أو غير المسجّلة الدخول أو التابعة لمؤسسة أخرى. ومع `relay = true` في `config.toml`، ينتظر ccam بجانب الجلسة بدلًا من ذلك، فيتابع `/ccam:switch` بنفسه، مع الإبقاء على خيارات الجلسة ووضع الأذونات ومستوى الجهد. بعض الأدوات الأخرى تبدّل تسجيل الدخول في جلسة قيد التشغيل بدلًا من ذلك؛ ويشرح [ccam وأدوات تبديل الحسابات الأخرى](/compare/account-switchers) الفرق.

## أمور يجب مراعاتها

- شغّل `ccam next --continue` في المجلد نفسه الذي عملت فيه؛ فالمحادثات تُخزَّن لكل مجلد مشروع.
- تجنّب تشغيل حسابين على المحادثة نفسها في الوقت نفسه. ينبّهك ccam إذا حاولت ذلك.
- دون سجل مشترك، تُنسخ المحادثة إلى الحساب الآخر وتمضي كل نسخة في طريقها؛ وتبقى نقاط الحفظ التي يستعيد منها `/rewind` الشيفرة مع الحساب الأول.
- قد تستهلك الرسالة الأولى على الحساب الجديد من حدّه أكثر من المعتاد، لأن ذلك الحساب لم يخزّن المحادثة مؤقتًا بعد.
- يعتمد السماح بالتبديل بين الحسابات بهذه الطريقة على شروط كل حساب.

راجع [حدود الاستخدام](/docs/usage-limits) لمعرفة كل خيارات `ccam usage` و`ccam next`.
