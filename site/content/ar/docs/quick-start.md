---
title: البدء السريع
description: أضف حسابَي Claude Code باستخدام ccam add، وبدّل بينهما باستخدام claude1 أو claude2 أو ccs، واختر ما يشاركانه.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

أضف كل حساب باستخدام `ccam add`، ثم شغّل Claude Code باستخدام `claude1` أو `claude2` أو `ccs <name>`. هذا كل ما يحتاجه معظم الناس.

## 1. أضف حسابك الأول

```sh
ccam add
```

يسألك المعالج عن اسم مثل `personal` أو `work`، وعن العناصر التي يجب أن يشاركها الحساب مع حساباتك الأخرى. ثم يشغّل `claude auth login`، الذي يفتح المتصفح لتسجّل الدخول. عند الانتهاء ترى:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

يطبع ccam رسائله الخاصة باللغة الكورية حاليًا؛ أما الأوامر والخيارات فهي نفسها في كل اللغات.

## 2. أضف الحساب الثاني

```sh
ccam add work
```

سجّل الدخول بالحساب الآخر في المتصفح. يصبح هذا الحساب `claude2`.

## 3. بدّل بين الحسابات

أيٌّ من هذه الأوامر يشغّل Claude Code بحساب العمل:

```sh
claude2
ccs work
ccs 2
```

كل ما يأتي بعد اسم الحساب يُمرَّر إلى Claude Code دون تغيير، فتواصل استخدام الخيارات التي تعرفها:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

تشغيل `ccs` وحده يفتح أداة اختيار: استخدم مفاتيح الأسهم وEnter، أو اضغط رقم الحساب.

## 4. اختر ما تريد مشاركته

افتراضيًا، تُشارَك بين الحسابات خوادم MCP و`settings.json` وملف `CLAUDE.md` العام، والوكلاء والأوامر والمهارات والإضافات. أما سجل المحادثات فلا يُشارَك. فعّله لكلا الحسابين إذا أردت متابعة المحادثات بينهما:

```sh
ccam share personal --history on
ccam share work --history on
```

يعرض ccam ما سينقله قبل أن يغيّر أي شيء. راجع [المشاركة](/docs/sharing) للتفاصيل.

## 5. عندما يبلغ حساب حدّه

```sh
ccs usage
ccs next --continue
```

يعرض `ccs usage` استخدام كل حساب خلال 5 ساعات وخلال الأسبوع. ويفتح `ccs next --continue` المحادثة نفسها في الحساب التالي. راجع [حدود الاستخدام](/docs/usage-limits).

## هل تستخدم ~/.claude بالفعل؟

سجّله بدلًا من تسجيل الدخول من جديد:

```sh
ccam import
```

يبقى المجلد في مكانه، ويستمر `claude` العادي في العمل كما كان.
