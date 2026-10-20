---
title: البدء السريع
metaTitle: "البدء السريع: أضف حسابات Claude Code وبدّل بينها"
description: أضف حسابَي Claude Code باستخدام ccam add، وامنحهما أوامر مختصرة مثل ccwork، وبدّل بينهما، واختر ما يشاركانه.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d5b0e295f3f0
---

أضف كل حساب باستخدام `ccam add`، ثم شغّل Claude Code باستخدام `ccam <name>` أو بأمر مختصر تختار اسمه بنفسك، مثل `ccwork`. هذا كل ما يحتاجه معظم الناس.

## 1. أضف حسابك الأول

```sh
ccam add
```

يسألك المعالج عن اسم مثل `personal` أو `work`، وعن العناصر التي يجب أن يشاركها الحساب مع حساباتك الأخرى، وعمّا إذا كنت تريد إنشاء أمر مختصر (راجع الخطوة 4). ثم يشغّل `claude auth login`، الذي يفتح المتصفح لتسجّل الدخول. على جهاز بعيد، قد ينتهي المتصفح عند صفحة `http://localhost:…/callback` لا تُحمَّل؛ الصق هذا العنوان في الطرفية وسيُكمل ccam تسجيل الدخول. عند الانتهاء ترى:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. انقل إعداداتك الحالية

إذا كنت قد استخدمت Claude Code من قبل، فإعداداتك و`CLAUDE.md` والإضافات وخوادم MCP موجودة في `~/.claude`. انسخها إلى الإعدادات المشتركة لتحصل عليها كل الحسابات:

```sh
ccam migrate
```

يعرض ccam ما سينسخه ويسألك أولًا. وهو يقرأ `~/.claude` فقط، ولا ينسخ أي تسجيل دخول أبدًا، ولا ينسخ سجل المحادثات إلا إذا أضفت `--history`. راجع [نسخ إعداد موجود](/docs/sharing#copy-an-existing-setup) للتفاصيل.

## 3. أضف الحساب الثاني

```sh
ccam add work
```

سجّل الدخول بالحساب الآخر في المتصفح.

## 4. بدّل بين الحسابات

يشغّل `ccam work` برنامج Claude Code بحساب العمل. ولأمر أقصر، امنح الحساب أمرًا مختصرًا خاصًا به، مثل `ccwork`. بعد سؤالك عن العناصر المشتركة، يسأل المعالج:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

اضغط `y`، ثم اكتب اسمًا، أو اضغط Enter لاستخدام الاسم المقترح:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

بعد تسجيل الدخول، يؤكد ccam ذلك:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

إذا تخطيت السؤال، يمكنك إضافة أمر مختصر في أي وقت:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

الآن أيٌّ من هذه الأوامر يشغّل Claude Code بحساب العمل:

```sh
ccwork
ccam work
ccam 2
```

تشغيل `ccam` وحده يفتح [لوحة التحكم](/features/dashboard): انتقل إلى حساب بمفاتيح الأسهم أو برقمه، ثم اضغط Enter.

كل ما يأتي بعد الأمر المختصر أو اسم الحساب يُمرَّر إلى Claude Code دون تغيير، فتواصل استخدام الخيارات التي تعرفها:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. اختر ما تريد مشاركته

افتراضيًا، تُشارَك بين الحسابات عبر المجموعة `default` خوادم MCP و`settings.json` وملف `CLAUDE.md` العام، والوكلاء والأوامر والمهارات والإضافات. أما سجل المحادثات فلا يُشارَك. أضفه إلى المجموعة إذا أردت متابعة المحادثات بين الحسابات:

```sh
ccam group set default --add history
```

يعرض ccam ما سينقله لكل حساب قبل أن يغيّر أي شيء. راجع [المشاركة عبر المجموعات](/features/share) للتفاصيل.

## 6. عندما يبلغ حساب حدّه

```sh
ccam usage
ccam next --continue
```

يعرض `ccam usage` استخدام كل حساب خلال 5 ساعات وخلال الأسبوع. ويفتح `ccam next --continue` المحادثة نفسها في الحساب التالي. وفي السكربتات، يبدّل `ccam --fallback work -p "…"` الحساب بنفسه. راجع [حدود الاستخدام](/docs/usage-limits).

## هل تستخدم ~/.claude بالفعل؟

سجّله بدلًا من تسجيل الدخول من جديد:

```sh
ccam import
```

يبقى المجلد في مكانه، ويستمر `claude` العادي في العمل كما كان.
