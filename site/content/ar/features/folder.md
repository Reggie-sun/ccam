---
title: الحسابات حسب المجلد
description: يجعل ccam folder مجلدًا وكل ما تحته يستخدم حسابًا واحدًا عندما تشغّل Claude Code دون اسم حساب، مثل work لشيفرة شركتك.
updated: 2026-10-23
order: 11
translation: machine
sourceHash: 48de98d42b99
---

يجعل `ccam folder work ~/code/company` كل تشغيل دون اسم حساب داخل `~/code/company`، وفي أي مجلد تحته، يبدأ Claude Code بالحساب `work`. وخارجه يبدأ الحساب الافتراضي كما من قبل. أما تسمية حساب، باستخدام `ccam personal` أو أمر مختصر مثل `ccpersonal`، فتشغّل ذلك الحساب دائمًا.

## الاستخدام

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

تستخدم عمليات التشغيل هذه حساب المجلد:

- `ccam -p "…"` وأي تشغيل آخر يكون وسيطه الأول خيارًا لـ Claude Code
- `ccam run` دون اسم، و`ccam --dry-run`
- `ccam --fallback -p "…"` في محاولته الأولى
- لوحة التحكم، أي `ccam` وحده، التي تبدأ ومؤشرها على حساب المجلد، حتى عندما تصل إلى المجلد عبر رابط رمزي

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `<account>` | اسم الحساب المراد استخدامه أو رقمه | |
| `[folder]` | المجلد، نسبيًا أو مطلقًا | المجلد الحالي |
| `--remove` | يزيل الإدخال الخاص بذلك المجلد بالضبط | معطّل |
| `--dry-run` | مدعوم: يطبع التغيير دون حفظه | معطّل |

## كيف يعمل

كل مجلد هو إدخال `[[folder]]` واحد في `~/.config/ccam/config.toml`، الملف الذي يقرؤه كل تشغيل في كل الأحوال:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. يحوّل `ccam folder` المجلد إلى مسار مطلق بعد حل الروابط الرمزية ويحفظه. وتعيين مجلد مرة أخرى يستبدل حسابه.
2. يستدعي التشغيل دون اسم `launcher.HereProfile`. وإذا لم توجد إدخالات، يأخذ الحساب الافتراضي فورًا.
3. وإلا فإنه يقارن مجلد العمل مع كل `path` كسلاسل نصية: المجلد نفسه أو مجلد تحته، ولا يكون أبدًا مجلدًا يبدأ بالأحرف نفسها فقط (لا يشمل `/code/company` المجلد `/code/companyx`). ويفوز الإدخال الأعمق، فيمكن أن يكون `~/code` للحساب `personal` بينما يكون `~/code/company` للحساب `work`.
4. فقط إذا لم يطابق مجلد العمل كما هو أي إدخال، يحل ccam روابطه الرمزية ويقارن مرة أخرى.
5. على Windows وmacOS تتجاهل المقارنة حالة الأحرف الكبيرة والصغيرة، كما تفعل أنظمة ملفاتهما عادةً.

لا يُقرأ أي ملف سوى `config.toml`، لذا لا يكلّف التشغيل داخل مجلد وقتًا يمكن قياسه (4.0 ms مع الإدخالات أو دونها، `hyperfine`، 300 تشغيل).

## الأمان والحالات الخاصة

- ينقل `ccam rename` إدخالات المجلدات إلى الاسم الجديد، ويحذف `ccam remove` إدخالات الحساب المحذوف.
- الإدخال الذي لم يعد حسابه موجودًا يوقف التشغيل دون اسم في ذلك المجلد: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. ويبلّغ `ccam doctor` عن مثل هذه الإدخالات، وعن الإدخالات التي اختفت مجلداتها.
- عند تشغيل `ccam folder --remove` داخل مجلد يأخذ حسابه من مجلد أعلى، يذكر ccam ذلك المجلد الأعلى بدلًا من أن يزيل شيئًا.
- يختار `ccam next` بالتدوير أو بحسب الاستخدام ولا ينظر إلى المجلدات.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [الأوامر المختصرة](/features/shortcuts)
- [إدارة الحسابات](/features/accounts)
- [مرجع الأوامر](/docs/commands)
