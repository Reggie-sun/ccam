---
title: "ہر اکاؤنٹ کے پس منظر کے سیشنز"
description: "ccam agents آپ کے تمام اکاؤنٹس کے Claude Code بیک گراؤنڈ ایجنٹس اور کھلے سیشنز ایک ٹیبل میں دکھاتا ہے، اور ccam attach کسی ایک کو اس کی ID سے کھولتا ہے۔"
updated: 2026-10-26
order: 16
translation: machine
sourceHash: 06281fb13119
---

Claude Code پس منظر کے سیشنز (`claude --bg`، ایجنٹ ویو) ہر اکاؤنٹ کے لیے الگ رکھتا ہے: `claude agents` صرف اس اکاؤنٹ کے سیشنز دکھاتا ہے جس کے طور پر وہ چل رہا ہو۔ `ccam agents` ہر اکاؤنٹ سے ایک ساتھ پوچھتا ہے اور سب کی فہرست اکٹھی دکھاتا ہے، اور `ccam attach <id|name>` پس منظر کا کوئی سیشن اسی اکاؤنٹ کے طور پر کھولتا ہے جس نے اسے شروع کیا تھا۔

## استعمال

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

`State` پس منظر کے سیشن کی حالت ہے (`working`، `blocked`، `done`، `failed`، `stopped`) اور یہ کہ اس کا پروسیس کیا کر رہا ہے (`busy`، `waiting`، `idle`)؛ کھلا انٹرایکٹو سیشن `open` دکھاتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--all` | وہ پس منظر کے سیشنز بھی دکھاتا ہے جو مکمل، ناکام یا بند ہو چکے ہوں | بند |
| `--json` | سیشنز JSON کے طور پر پرنٹ کرتا ہے، ہر ایک کے ساتھ وہ `account` جس کا وہ ہے | بند |

`ccam attach <id|name> [claude args]` ID کے بعد والی ہر چیز، فلیگز سمیت، بغیر تبدیلی کے `claude attach` کو دیتا ہے: `ccam attach 49ec --dry-run`، `--dry-run` کو Claude Code کے حوالے کرتا ہے۔ اس کے بجائے کمانڈ پرنٹ کرنے کے لیے `--dry-run` کو ID سے پہلے لگائیں (`ccam --dry-run attach 49ec` یا `ccam attach --dry-run 49ec`)۔

## یہ کیسے کام کرتا ہے

1. ہر اکاؤنٹ کے لیے ccam اکاؤنٹ کے `CLAUDE_CONFIG_DIR` کے ساتھ `claude agents --json` چلاتا ہے (دیا گیا ہو تو `--all` کے ساتھ)، متوازی طور پر، ہر ایک 15 سیکنڈ کے ٹائم آؤٹ کے ساتھ۔ یہ اکاؤنٹ کی اپنی سیشن فہرست پڑھتا ہے اور Claude Code کی بیک گراؤنڈ سروس شروع نہیں کرتا۔
2. ہر اندراج کو اکاؤنٹ کا نام ملتا ہے، اور فہرست شروع ہونے کے وقت کے لحاظ سے ترتیب دی جاتی ہے، نیا پہلے۔ جو اکاؤنٹ جواب نہ دے اس کا نام اسٹینڈرڈ ایرر پر بتایا جاتا ہے اور اسے چھوڑ دیا جاتا ہے، پھر ccam کوڈ 1 کے ساتھ باہر نکلتا ہے، تاکہ کوئی اسکرپٹ ناکامی کو خالی فہرست نہ سمجھ لے۔
3. `ccam attach` ہر اکاؤنٹ کے فولڈر میں `jobs/<id>` تلاش کرتا ہے، جہاں Claude Code پس منظر کا سیشن رکھتا ہے، اور اس اکاؤنٹ کے طور پر `claude attach <id>` چلاتا ہے، اسی طرح جیسے `ccam <name>` Claude Code شروع کرتا ہے۔ جو لفظ ID کا ابتدائی حصہ نہ ہو اسے سیشنز کے ناموں سے ملایا جاتا ہے، حروف کے چھوٹے بڑے ہونے سے قطع نظر؛ جو لفظ ID جیسا لگے اسے کبھی نام نہیں سمجھا جاتا۔

## حفاظت اور خاص صورتیں

- پس منظر کے سیشنز کی IDs ایک اکاؤنٹ کی ہوتی ہیں: کسی دوسرے اکاؤنٹ کے طور پر چلایا گیا `claude attach` انہیں نہیں ڈھونڈ سکتا، اور `ccam attach` یہی مسئلہ حل کرتا ہے۔
- اگر دو اکاؤنٹس کے سیشنز کی IDs کا ابتدائی حصہ یا نام کا کوئی لفظ ایک جیسا ہو تو `ccam attach` ان کی فہرست دکھاتا ہے اور لمبی ID یا نام مانگتا ہے۔
- ccam اپنے شیئر کیے گئے MCP سرورز Claude Code کو ایک ہی آرگیومنٹ میں `--mcp-config=<file>` کے ساتھ دیتا ہے، جسے `claude agents` اور `claude attach` قبول کرتے ہیں؛ پس منظر کے سیشنز دوبارہ شروع ہونے پر بھی وہی فائل استعمال کرتے رہتے ہیں۔
- API key اکاؤنٹ کی بیک گراؤنڈ سروس اپنی کلید اس شیل سے لیتی ہے جس نے اسے شروع کیا۔ ایسے اکاؤنٹ کا پہلا پس منظر کا سیشن ccam کے ذریعے شروع کریں (`ccam gw --bg …`)، تاکہ سروس کے پاس کلید ہو۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [اسکرپٹس کے لیے اکاؤنٹ چننا](/features/pick)
- [ڈیش بورڈ](/features/dashboard)
