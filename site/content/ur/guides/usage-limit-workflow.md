---
title: "جب Claude Code استعمال کی حد تک پہنچ جائے تو کیا کریں"
description: "ری سیٹ کا انتظار کریں یا وہی گفتگو فوراً Claude Code کے دوسرے اکاؤنٹ پر جاری رکھیں۔ 5 گھنٹے اور ہفتہ وار حدیں کیسے کام کرتی ہیں، مرحلہ وار۔"
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

جب کوئی اکاؤنٹ استعمال کی حد تک پہنچ جائے (جسے اکثر rate limit کہا جاتا ہے) تو آپ کے پاس دو راستے ہیں: حد ری سیٹ ہونے تک انتظار کریں، یا وہی گفتگو ابھی اپنے کسی دوسرے اکاؤنٹ پر جاری رکھیں۔ ccam کے ساتھ دوسرے راستے کے لیے ایک کمانڈ کافی ہے، `ccam next --auto --continue`۔

## Claude Code کی کون سی حدیں ہیں؟

Pro اور Max پلانز پر Claude Code کی [5 گھنٹے کی سیشن حد اور ایک ہفتہ وار حد](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) ہوتی ہے۔ دونوں حدیں ویب اور ایپس والے Claude کے ساتھ مشترک ہیں، اور استعمال دونوں میں ایک ساتھ شمار ہوتا ہے۔ جب ایک ختم ہو جائے تو Claude Code [ایک پیغام دکھاتا ہے](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)، جیسے "You've hit your session limit · resets 3:45pm"۔ ماڈل بدلنے سے فائدہ نہیں ہوتا، کیونکہ یہ حدیں تمام ماڈلز پر لاگو ہوتی ہیں۔

## پہلا راستہ: ری سیٹ کا انتظار کریں

اگر آپ کے پاس ایک ہی اکاؤنٹ ہے تو انتظار کریں۔ Claude Code کے حالیہ ورژن [کھلے سیشن میں انتظار کرتے ہیں](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) اور حد ری سیٹ ہوتے ہی خود کام جاری رکھتے ہیں۔ Claude Code کے اندر `/usage` دکھاتا ہے کہ آپ نے اپنے پلان کی حدوں میں سے کتنا استعمال کر لیا ہے۔

## دوسرا راستہ: ابھی دوسرے اکاؤنٹ پر جاری رکھیں

اگر آپ کے پاس دوسرا اکاؤنٹ ہے تو آپ وہی گفتگو فوراً جاری رکھ سکتے ہیں۔ اسے ایک بار سیٹ اپ کریں:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

اب کام کے دوران اسٹیٹس لائن اکاؤنٹ اور اس کا 5 گھنٹے کا استعمال دکھاتی ہے، اس لیے حد اچانک سامنے نہیں آتی۔

جب حد پوری ہو جائے تو Claude Code سے باہر نکلیں (`/exit`)، پھر:

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

سب سے زیادہ گنجائش والے اکاؤنٹ میں گفتگو کھولیں:

```sh
ccam next --auto --continue
```

Claude Code، `work` کے طور پر موجودہ فولڈر کی سب سے حالیہ گفتگو کے ساتھ شروع ہوتا ہے، جس میں اب تک کی ہر بات شامل ہوتی ہے۔ ہسٹری شیئر کرنا (اوپر `ccam group set default --add history`) اسے ایک ہی گفتگو رکھتا ہے؛ اس کے بغیر Claude Code گفتگو کو `work` میں کاپی کر دیتا ہے، اور ccam بتاتا ہے کہ وہ پہلا جواب کتنے ٹوکنز دوبارہ پڑھے گا۔

Claude Code میں [ccam پلگ اِن](/docs/install#claude-code-plugin) کے ساتھ، جو سیشن اپنی حد پر رک جائے وہ ایک اطلاع دکھاتا ہے اور بھری ہوئی ونڈو ریکارڈ کر لیتا ہے، اس لیے `ccam next --auto` اس اکاؤنٹ کو اس کے ری سیٹ سے پہلے دوبارہ نہیں چنتا، اور `/ccam:switch` (یا `/ccam:switch work`) گفتگو کو `ccam next --continue` کے لیے نوٹ کر لیتا ہے تاکہ اندازہ لگانے کو کچھ نہ رہے۔

**Team اور Enterprise سیٹیں۔** خودکار انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے۔ Anthropic کا ہیلپ سینٹر کہتا ہے کہ "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (کسی Team پلان کی ہر سیٹ ایک نامزد فرد کو تفویض کی جاتی ہے، اور صرف اسی شخص کو اسے استعمال کرنے کی اجازت ہے) اور "Users can't route their work through another user's account to use that user's remaining usage" (صارفین کسی دوسرے صارف کا باقی ماندہ استعمال کام میں لانے کے لیے اپنا کام اس کے اکاؤنٹ کے ذریعے نہیں چلا سکتے) ([Team پلانز](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ اسی طرح [Enterprise پلانز](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) کے لیے بھی)۔ اس لیے سیٹوں کے درمیان خودکار انتخاب استعمال کرنے سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ اس پر انحصار کرنے سے پہلے اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔ ذاتی پلانز (Pro، Max) کی بات [عام سوالات میں شرائط والے نوٹ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) میں کی گئی ہے۔

## کیا ccam خود بخود اکاؤنٹ تبدیل کرتا ہے؟

نہیں۔ آپ کے کام کے دوران ccam چل نہیں رہا ہوتا: یہ Claude Code کو شروع کرتا ہے اور خود کو اس سے بدل دیتا ہے۔ اکاؤنٹ کب بدلنا ہے، یہ آپ طے کرتے ہیں، اور `ccam next --auto` وہ اکاؤنٹ چنتا ہے جس کا 5 گھنٹے کا استعمال سب سے کم ہو، اور ان اکاؤنٹس کو چھوڑ دیتا ہے جو 100% پر ہوں، لاگ اِن نہ ہوں، یا کسی دوسری تنظیم کے ہوں۔ `config.toml` میں `relay = true` کے ساتھ ccam اس کے بجائے سیشن کے ساتھ انتظار کرتا ہے، اس لیے `/ccam:switch` خود آگے بڑھ جاتا ہے، سیشن کے آپشنز، اجازتوں کا موڈ اور effort لیول برقرار رکھتے ہوئے۔ کچھ دوسرے ٹولز اس کی بجائے چلتے ہوئے سیشن کا لاگ اِن بدل دیتے ہیں؛ [ccam اور دوسرے اکاؤنٹ سوئچرز](/compare/account-switchers) میں یہ فرق بیان کیا گیا ہے۔

## یاد رکھنے کی باتیں

- `ccam next --continue` اسی فولڈر میں چلائیں جس میں آپ کام کر رہے تھے؛ گفتگو ہر پروجیکٹ فولڈر کے حساب سے الگ محفوظ ہوتی ہے۔
- ایک ہی گفتگو پر بیک وقت دو اکاؤنٹس چلانے سے گریز کریں۔ اگر آپ ایسا کرنے کی کوشش کریں تو ccam آپ کو خبردار کرتا ہے۔
- شیئر شدہ ہسٹری کے بغیر گفتگو دوسرے اکاؤنٹ میں کاپی ہو جاتی ہے اور دونوں کاپیاں اپنے اپنے راستے چلتی ہیں؛ وہ checkpoints جن سے `/rewind` کوڈ بحال کرتا ہے، پہلے اکاؤنٹ کے پاس رہتے ہیں۔
- نئے اکاؤنٹ پر پہلا پیغام معمول سے زیادہ حد استعمال کر سکتا ہے، کیونکہ اس اکاؤنٹ نے ابھی گفتگو کو کیش نہیں کیا ہوتا۔
- آپ اس طرح اکاؤنٹس کے درمیان تبدیل کر سکتے ہیں یا نہیں، یہ ہر اکاؤنٹ کی شرائط پر منحصر ہے۔

`ccam usage` اور `ccam next` کے تمام آپشنز کے لیے [استعمال کی حدیں](/docs/usage-limits) دیکھیں۔
