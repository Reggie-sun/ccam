---
title: "فولڈر کے حساب سے اکاؤنٹس"
description: "ccam folder کسی فولڈر اور اس کے نیچے ہر چیز کو ایک اکاؤنٹ استعمال کرواتا ہے جب آپ اکاؤنٹ کے نام کے بغیر Claude Code شروع کریں، جیسے کمپنی کے کوڈ کے لیے work۔"
updated: 2026-10-26
order: 11
translation: machine
sourceHash: 25f82da64580
---

`ccam folder work ~/code/company` کے بعد، `~/code/company` یا اس کے ذیلی فولڈرز میں اکاؤنٹ کے نام کے بغیر کوئی بھی لانچ Claude Code کو `work` کے طور پر شروع کرتا ہے۔ اس سے باہر، ڈیفالٹ اکاؤنٹ پہلے کی طرح شروع ہوتا ہے۔ کسی اکاؤنٹ کا نام لینا، `ccam personal` سے یا `ccpersonal` جیسے شارٹ کٹ سے، ہمیشہ وہی اکاؤنٹ شروع کرتا ہے۔

## استعمال

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

یہ لانچز فولڈر کا اکاؤنٹ استعمال کرتے ہیں:

- `ccam -p "…"` اور کوئی بھی دوسرا لانچ جس کا پہلا آرگیومنٹ Claude Code کا آپشن ہو
- نام کے بغیر `ccam run`، اور `ccam --dry-run`
- `ccam --fallback -p "…"` اپنی پہلی کوشش کے لیے
- ڈیش بورڈ، یعنی اکیلا `ccam`، جو کرسر کو فولڈر کے اکاؤنٹ پر رکھ کر شروع ہوتا ہے، تب بھی جب آپ اس فولڈر تک کسی symbolic link کے ذریعے پہنچے ہوں

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `<account>` | استعمال ہونے والے اکاؤنٹ کا نام یا نمبر | |
| `[folder]` | فولڈر، نسبتی یا مطلق | موجودہ فولڈر |
| `--remove` | ٹھیک اسی فولڈر کا اندراج ہٹاتا ہے | بند |
| `--dry-run` | قبول کیا جاتا ہے: تبدیلی محفوظ کیے بغیر پرنٹ کرتا ہے | بند |

## یہ کیسے کام کرتا ہے

ہر فولڈر `~/.config/ccam/config.toml` میں ایک `[[folder]]` اندراج ہے، وہی فائل جسے ہر لانچ ویسے بھی پڑھتا ہے:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` فولڈر کو symbolic links حل کر کے ایک مطلق پاتھ میں بدلتا ہے اور محفوظ کرتا ہے۔ کسی فولڈر کو دوبارہ سیٹ کرنے سے اس کا اکاؤنٹ بدل جاتا ہے۔
2. نام کے بغیر لانچ `launcher.HereProfile` کو کال کرتا ہے۔ کوئی اندراج نہ ہو تو یہ فوراً ڈیفالٹ اکاؤنٹ لے لیتا ہے۔
3. ورنہ یہ ورکنگ ڈائریکٹری کو ہر `path` سے اسٹرنگ کے طور پر ملاتا ہے: خود وہ فولڈر یا اس کے نیچے کوئی فولڈر، کبھی ایسا فولڈر نہیں جو صرف انہی حروف سے شروع ہوتا ہو (`/code/company` میں `/code/companyx` شامل نہیں)۔ سب سے گہرا اندراج جیتتا ہے، اس لیے `~/code` کے لیے `personal` ہو سکتا ہے جبکہ `~/code/company` کے لیے `work`۔
4. صرف اگر ورکنگ ڈائریکٹری دیے گئے روپ میں کسی سے نہ ملے، تب ccam اس کے symbolic links حل کرتا ہے اور دوبارہ ملاتا ہے۔
5. Windows اور macOS پر موازنہ بڑے اور چھوٹے حروف کو نظر انداز کرتا ہے، جیسا کہ ان کے فائل سسٹمز عموماً کرتے ہیں۔

`config.toml` کے سوا کوئی فائل نہیں پڑھی جاتی، اس لیے کسی فولڈر کے اندر لانچ میں کوئی قابلِ پیمائش وقت کا اضافہ نہیں ہوتا (اندراجات کے ساتھ یا ان کے بغیر 4.0 ms، `hyperfine`، 300 رنز)۔

## حفاظت اور خاص صورتیں

- `ccam rename` فولڈر کے اندراجات نئے نام پر منتقل کرتا ہے، اور `ccam remove` ہٹائے گئے اکاؤنٹ کے اندراجات حذف کر دیتا ہے۔
- جس اندراج کا اکاؤنٹ اب موجود نہ ہو وہ اس فولڈر میں نام کے بغیر لانچ کو روک دیتا ہے: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`۔ `ccam doctor` ایسے اندراجات کی اطلاع دیتا ہے، اور ان اندراجات کی بھی جن کا فولڈر ختم ہو چکا ہو۔
- اگر آپ کسی ایسے ذیلی فولڈر میں `ccam folder --remove` چلائیں جو اپنا اکاؤنٹ کسی پیرنٹ فولڈر سے لیتا ہو، تو ccam پیرنٹ فولڈر کا نام بتاتا ہے اور کچھ نہیں ہٹاتا۔
- `ccam next` گردش یا استعمال کے حساب سے چنتا ہے اور فولڈرز کو نہیں دیکھتا۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [شارٹ کٹ کمانڈز](/features/shortcuts)
- [اکاؤنٹس کا انتظام](/features/accounts)
- [کمانڈز کا حوالہ](/docs/commands)
