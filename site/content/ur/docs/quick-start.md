---
title: "فوری آغاز"
metaTitle: "فوری آغاز: Claude Code اکاؤنٹس شامل کریں اور بدلیں"
description: ccam add سے Claude Code کے دو اکاؤنٹس شامل کریں، انہیں ccwork جیسے شارٹ کٹ دیں، ان کے درمیان تبدیل کریں اور طے کریں کہ وہ کیا شیئر کریں۔
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

ہر اکاؤنٹ `ccam add` سے شامل کریں، پھر Claude Code کو `ccam <name>` سے یا ایسے شارٹ کٹ سے شروع کریں جس کا نام آپ خود رکھیں، جیسے `ccwork`۔ زیادہ تر لوگوں کو بس اتنا ہی چاہیے۔

کیا آپ پہلے سے Claude Code استعمال کرتے ہیں؟ شروع کرنے سے پہلے اس صفحے کے آخر میں [کیا آپ پہلے سے ~/.claude استعمال کر رہے ہیں؟](#already-using-claude) پڑھیں۔

## 1. اپنا پہلا اکاؤنٹ شامل کریں

```sh
ccam add
```

وزرڈ ایک نام پوچھتا ہے، جیسے `personal` یا `work`، یہ کہ اکاؤنٹ آپ کے دوسرے اکاؤنٹس کے ساتھ کون سے آئٹمز شیئر کرے، اور یہ کہ شارٹ کٹ کمانڈ بنائی جائے یا نہیں (مرحلہ 3 دیکھیں)۔ پھر یہ `claude auth login` چلاتا ہے، جو براؤزر کھولتا ہے تاکہ آپ سائن اِن کر سکیں۔ ریموٹ مشین پر آپ کا براؤزر ایک ایسے `http://localhost:…/callback` صفحے پر رک سکتا ہے جو لوڈ نہیں ہوتا؛ وہ ایڈریس ٹرمینل میں پیسٹ کریں اور ccam سائن اِن مکمل کر دیتا ہے۔ سائن اِن مکمل ہونے پر آپ کو یہ نظر آتا ہے:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. اپنی موجودہ سیٹنگز ساتھ لائیں

اگر آپ پہلے Claude Code استعمال کر چکے ہیں تو آپ کی سیٹنگز، `CLAUDE.md` اور ایکسٹینشنز `~/.claude` میں ہیں، اور آپ کے MCP سرورز `~/.claude.json` میں ہیں۔ انہیں شیئر شدہ `default` گروپ میں کاپی کریں تاکہ ہر اکاؤنٹ کو یہ مل جائیں:

```sh
ccam migrate
```

ccam دکھاتا ہے کہ وہ کیا کاپی کرے گا اور پہلے پوچھتا ہے۔ یہ `~/.claude` کو صرف پڑھتا ہے، کبھی لاگ اِن کاپی نہیں کرتا، اور گفتگو کی ہسٹری تب تک شامل نہیں کرتا جب تک آپ `--history` نہ لگائیں۔ تفصیلات کے لیے [موجودہ سیٹ اپ کاپی کرنا](/docs/sharing#copy-an-existing-setup) دیکھیں۔

## 3. دوسرا اکاؤنٹ شامل کریں

```sh
ccam add work
```

وزرڈ نام نہیں پوچھتا، کیونکہ وہ آپ نے کمانڈ لائن پر دے دیا ہے، اور شامل ہونے کے لیے گروپس دکھاتا ہے جن میں `default` چیک ہوتا ہے۔ اس بار اکاؤنٹ کو اس کی اپنی شارٹ کٹ کمانڈ دیں، جیسے `ccwork`، تاکہ آپ اسے ایک لفظ سے شروع کر سکیں۔ گروپس کے بعد وزرڈ پوچھتا ہے:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` دبائیں، پھر کوئی نام لکھیں، یا تجویز کردہ نام استعمال کرنے کے لیے Enter دبائیں:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

پھر براؤزر میں دوسرے اکاؤنٹ سے سائن اِن کریں۔ ccam اکاؤنٹ اور اس کے شارٹ کٹ کی تصدیق کرتا ہے:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. اکاؤنٹ تبدیل کریں

`ccam work` ورک اکاؤنٹ کے طور پر Claude Code شروع کرتا ہے۔ اگر آپ نے شارٹ کٹ والا سوال چھوڑ دیا تھا تو کسی بھی وقت شارٹ کٹ شامل کریں:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

اب ان میں سے ہر کمانڈ ورک اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے (`2` اس کا `ccam list` میں نمبر ہے):

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

ڈیفالٹ طور پر MCP سرورز، `settings.json`، گلوبل `CLAUDE.md` اور ایکسٹینشنز (ایجنٹس، کمانڈز، اسکلز اور پلگ اِنز) `default` گروپ کے ذریعے اکاؤنٹس کے درمیان شیئر ہوتے ہیں۔ گفتگو کی ہسٹری شیئر نہیں ہوتی۔ اگر آپ گفتگو کو ایک اکاؤنٹ سے دوسرے میں جاری رکھنا چاہتے ہیں تو اسے گروپ میں شامل کریں:

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

**Team اور Enterprise سیٹیں۔** خودکار اکاؤنٹ انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید اور، جب آپ اسے آن کریں تو، [`auto_switch`](/features/next#auto-switch)) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے، جس سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ اس پر انحصار کرنے سے پہلے [عام سوالات](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) پڑھیں اور اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔

## کیا آپ پہلے سے ~/.claude استعمال کر رہے ہیں؟

دوبارہ سائن اِن کرنے کے بجائے اسے رجسٹر کریں:

```sh
ccam import
```

فولڈر اپنی جگہ پر رہتا ہے، اور سادہ `claude` پہلے کی طرح کام کرتا رہتا ہے۔
