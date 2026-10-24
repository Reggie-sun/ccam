---
title: "فوری آغاز"
metaTitle: "فوری آغاز: Claude Code اکاؤنٹس شامل کریں اور بدلیں"
description: ccam add سے Claude Code کے دو اکاؤنٹس شامل کریں، انہیں ccwork جیسے شارٹ کٹ دیں، ان کے درمیان تبدیل کریں اور طے کریں کہ وہ کیا شیئر کریں۔
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "e98d0637356a"
---

ہر اکاؤنٹ `ccam add` سے شامل کریں، پھر Claude Code کو `ccam <name>` سے یا ایسے شارٹ کٹ سے شروع کریں جس کا نام آپ خود رکھیں، جیسے `ccwork`۔ زیادہ تر لوگوں کو بس اتنا ہی چاہیے۔

## 1. اپنا پہلا اکاؤنٹ شامل کریں

```sh
ccam add
```

وزرڈ ایک نام پوچھتا ہے، جیسے `personal` یا `work`، یہ کہ اکاؤنٹ آپ کے دوسرے اکاؤنٹس کے ساتھ کون سے آئٹمز شیئر کرے، اور یہ کہ شارٹ کٹ کمانڈ بنائی جائے یا نہیں (مرحلہ 4 دیکھیں)۔ پھر یہ `claude auth login` چلاتا ہے، جو براؤزر کھولتا ہے تاکہ آپ سائن اِن کر سکیں۔ ریموٹ مشین پر براؤزر ایک ایسے `http://localhost:…/callback` صفحے پر رک سکتا ہے جو لوڈ نہیں ہوتا؛ وہ ایڈریس ٹرمینل میں پیسٹ کریں اور ccam سائن اِن مکمل کر دیتا ہے۔ کام مکمل ہونے پر آپ کو یہ نظر آتا ہے:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. اپنی موجودہ سیٹنگز ساتھ لائیں

اگر آپ پہلے Claude Code استعمال کر چکے ہیں تو آپ کی سیٹنگز، `CLAUDE.md`، ایکسٹینشنز اور MCP سرورز `~/.claude` میں ہیں۔ انہیں شیئر شدہ سیٹنگز میں کاپی کریں تاکہ ہر اکاؤنٹ کو یہ مل جائیں:

```sh
ccam migrate
```

ccam دکھاتا ہے کہ وہ کیا کاپی کرے گا اور پہلے پوچھتا ہے۔ یہ `~/.claude` کو صرف پڑھتا ہے، کبھی لاگ اِن کاپی نہیں کرتا، اور گفتگو کی ہسٹری تب تک شامل نہیں کرتا جب تک آپ `--history` نہ لگائیں۔ تفصیلات کے لیے [موجودہ سیٹ اپ کاپی کرنا](/docs/sharing#copy-an-existing-setup) دیکھیں۔

## 3. دوسرا اکاؤنٹ شامل کریں

```sh
ccam add work
```

براؤزر میں دوسرے اکاؤنٹ سے سائن اِن کریں۔

## 4. اکاؤنٹ تبدیل کریں

`ccam work` ورک اکاؤنٹ کے طور پر Claude Code شروع کرتا ہے۔ اس سے چھوٹی کمانڈ کے لیے اکاؤنٹ کو اس کی اپنی شارٹ کٹ کمانڈ دیں، جیسے `ccwork`۔ شیئر کیے جانے والے آئٹمز کے بعد وزرڈ پوچھتا ہے:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` دبائیں، پھر کوئی نام لکھیں، یا تجویز کردہ نام استعمال کرنے کے لیے Enter دبائیں:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

لاگ اِن کے بعد ccam اس کی تصدیق کرتا ہے:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

اگر آپ نے یہ سوال چھوڑ دیا تھا تو کسی بھی وقت شارٹ کٹ شامل کریں:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

اب ان میں سے کوئی بھی کمانڈ ورک اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے:

```sh
ccwork
ccam work
ccam 2
```

صرف `ccam` چلانے سے [ڈیش بورڈ](/features/dashboard) کھلتا ہے: ایرو کیز یا اکاؤنٹ کے نمبر سے کسی اکاؤنٹ پر جائیں، پھر Enter دبائیں۔

شارٹ کٹ یا اکاؤنٹ کے نام کے بعد کی ہر چیز بغیر تبدیلی کے Claude Code کو جاتی ہے، اس لیے آپ وہی آپشنز استعمال کرتے رہتے ہیں جو آپ جانتے ہیں:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. طے کریں کہ کیا شیئر کرنا ہے

ڈیفالٹ طور پر MCP سرورز، `settings.json`، گلوبل `CLAUDE.md`، اور ایجنٹس، کمانڈز، اسکلز اور پلگ اِنز گروپ `default` کے ذریعے اکاؤنٹس کے درمیان شیئر ہوتے ہیں۔ گفتگو کی ہسٹری شیئر نہیں ہوتی۔ اگر آپ گفتگو کو ایک اکاؤنٹ سے دوسرے میں جاری رکھنا چاہتے ہیں تو اسے گروپ میں شامل کریں:

```sh
ccam group set default --add history
```

ccam کچھ بھی تبدیل کرنے سے پہلے دکھاتا ہے کہ وہ ہر اکاؤنٹ کے لیے کیا منتقل کرے گا۔ تفصیلات کے لیے [گروپس کے ساتھ شیئرنگ](/features/share) دیکھیں۔

## 6. جب کوئی اکاؤنٹ اپنی حد تک پہنچ جائے

```sh
ccam usage
ccam next --continue
```

`ccam usage` ہر اکاؤنٹ کا 5 گھنٹے اور ہفتہ وار استعمال دکھاتا ہے۔ `ccam next --continue` وہی گفتگو اگلے اکاؤنٹ میں کھولتا ہے۔ اسکرپٹس میں `ccam --fallback work -p "…"` خود اکاؤنٹ بدل دیتا ہے۔ [استعمال کی حدیں](/docs/usage-limits) دیکھیں۔

**Team اور Enterprise سیٹیں۔** خودکار انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے۔ Anthropic کا ہیلپ سینٹر کہتا ہے کہ "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (کسی Team پلان کی ہر سیٹ ایک نامزد فرد کو تفویض کی جاتی ہے، اور صرف اسی شخص کو اسے استعمال کرنے کی اجازت ہے) اور "Users can't route their work through another user's account to use that user's remaining usage" (صارفین کسی دوسرے صارف کا باقی ماندہ استعمال کام میں لانے کے لیے اپنا کام اس کے اکاؤنٹ کے ذریعے نہیں چلا سکتے) ([Team پلانز](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ اسی طرح [Enterprise پلانز](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) کے لیے بھی)۔ اس لیے سیٹوں کے درمیان خودکار انتخاب استعمال کرنے سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ اس پر انحصار کرنے سے پہلے اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔ ذاتی پلانز (Pro، Max) کی بات [عام سوالات میں شرائط والے نوٹ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) میں کی گئی ہے۔

## کیا آپ پہلے سے ~/.claude استعمال کر رہے ہیں؟

دوبارہ سائن اِن کرنے کے بجائے اسے رجسٹر کریں:

```sh
ccam import
```

فولڈر اپنی جگہ پر رہتا ہے، اور سادہ `claude` پہلے کی طرح کام کرتا رہتا ہے۔
