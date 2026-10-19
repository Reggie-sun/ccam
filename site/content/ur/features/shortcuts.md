---
title: "شارٹ کٹ کمانڈز"
description: "ccam shortcut سے کسی اکاؤنٹ کو ccwork جیسی اپنی کمانڈ دیں، اور ccam shims install سے ~/.local/bin میں لنکس دوبارہ بنائیں۔"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

شارٹ کٹ ایک کمانڈ ہے جس کا نام آپ خود رکھتے ہیں اور جو ایک اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے: `ccwork --resume` وہی کرتا ہے جو `ccam work --resume` کرتا ہے۔ یہ `~/.local/bin` میں ccam کا ایک لنک ہے، اس لیے یہ bash، zsh، اسکرپٹس اور ہر اس چیز میں کام کرتا ہے جو `PATH` سے پروگرام چلاتی ہے۔ شارٹ کٹس اختیاری ہیں؛ جو اکاؤنٹس آپ اکثر شروع کرتے ہیں ان کے لیے ایک بنا لیں۔

## استعمال

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam ہر تبدیلی کی تصدیق کرتا ہے:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

ایک اکاؤنٹ کا زیادہ سے زیادہ ایک شارٹ کٹ ہوتا ہے، جو `~/.config/ccam/config.toml` میں `shortcut` کے طور پر محفوظ ہوتا ہے۔ `ccam add` کا وزرڈ `Create a shortcut command for <name>? [y/N]` اور پھر `Shortcut name [cc<name>]:` پوچھتا ہے۔ وزرڈ کے باہر، `ccam add` صرف `--shortcut` کے ساتھ شارٹ کٹ بناتا ہے۔ `ccam list` ہر اکاؤنٹ کا شارٹ کٹ دکھاتا ہے، یا `-`۔

شارٹ کٹ ہر آرگیومنٹ Claude Code کو دیتا ہے، بالکل اسی طرح جیسے `ccam <name>` نام کے بعد والے آرگیومنٹس دیتا ہے۔ اس میں `--dry-run` اور `--help` بھی شامل ہیں: کسی لانچ کا پیشگی جائزہ لینے کے لیے `ccam --dry-run work` استعمال کریں۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--remove` | اکاؤنٹ کا شارٹ کٹ اور اس کا لنک حذف کرتا ہے | بند |
| `-h`، `--help` | کمانڈ کی مدد | |
| `--dry-run` | `ccam shortcut` اسے لاگو کرتا ہے: یہ نام جانچتا ہے اور محفوظ یا لنک کیے بغیر `Would set the shortcut of work to ccwork.` یا `Would remove ccwork, the shortcut of work.` پرنٹ کرتا ہے | بند |

`ccam shims install` کوئی آرگیومنٹ نہیں لیتا اور `--dry-run` کو `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` کے ساتھ مسترد کرتا ہے۔

نام کے اصول، جو اسی ترتیب سے جانچے جاتے ہیں:

- `a-z`، `0-9`، `.`، `-` اور `_` کے 1 سے 64 حروف، جو کسی حرف یا ہندسے سے شروع ہوں؛
- `ccam`، `claude` یا کوئی ccam کمانڈ جیسے `run` یا `usage` نہیں؛
- کسی دوسرے اکاؤنٹ کا شارٹ کٹ نہیں؛
- `~/.local/bin` میں کوئی ایسی فائل نہیں جو پہلے سے ccam کا لنک نہ ہو؛
- آپ کے `PATH` پر کسی دوسرے قابلِ عمل پروگرام کا نام نہیں۔

کمانڈ اور `--remove` ایک ساتھ دینے پر `give a command or --remove, not both` کے ساتھ عمل رک جاتا ہے۔ وہی نام دوبارہ دینے پر جو اکاؤنٹ کا پہلے سے ہے، `work: ccwork (unchanged)` پرنٹ ہوتا ہے اور لنک غائب ہو تو دوبارہ بن جاتا ہے۔

## یہ کیسے کام کرتا ہے

نام مقرر کرنے پر `checkShortcut` چلتا ہے، جو bin فولڈر اور ccam بائنری کے اصل راستے کے ساتھ `shim.Check` کو کال کرتا ہے:

<!-- motion:shortcuts-check -->

نام ترتیب سے ہر جانچ سے گزرتا ہے۔ ہر جانچ کے ساتھ ایک ایسا نام ہے جسے وہ مسترد کرتی ہے: `Work` پیٹرن پر پورا نہیں اترتا، `run` ایک ccam کمانڈ ہے، `ccside` کسی دوسرے اکاؤنٹ کا ہے، `ccfile`، `~/.local/bin` میں ایک موجودہ فائل ہے، اور `git`، `/usr/bin/git` کو چھپا دے گا۔ صرف وہ نام محفوظ اور لنک ہوتا ہے جو ہر جانچ سے گزر جائے۔

1. `shim.Check` نام کو `^[a-z0-9][a-z0-9._-]{0,63}$`، پھر `config.Reserved`، پھر `cfg.FindShortcut` سے ملاتا ہے۔ یہ `~/.local/bin/<name>` پر `os.Lstat` کال کرتا ہے، اور `PATH` کے ہر دوسرے فولڈر کے لیے `os.Stat` کال کر کے ایسی قابلِ عمل فائل کو مسترد کرتا ہے جو ccam کا لنک نہ ہو۔
2. ccam، `config.toml` کو `fsutil.AtomicWrite` سے محفوظ کرتا ہے (عارضی فائل، `fsync`، rename)۔
3. `shim.Remove` پرانے شارٹ کٹ کا لنک حذف کرتا ہے، لیکن صرف تب جب وہ ccam کا لنک ہو۔
4. `shim.Install` ہر مقرر شدہ شارٹ کٹ کا غائب لنک بناتا ہے، ccam بائنری کی طرف `os.Symlink` کے ساتھ، جس کا راستہ `os.Executable` اور `filepath.EvalSymlinks` سے حل کیا جاتا ہے۔ یہ کبھی کچھ نہیں ہٹاتا۔

لنک تب ccam کا اپنا شمار ہوتا ہے جب `os.Readlink`، ccam کا راستہ لوٹائے، یا جب لنک اور ccam کا راستہ ایک ہی فائل تک پہنچیں۔ ہر حذف اور ہر اوور رائٹ کی جانچ اسی ٹیسٹ سے گزرتی ہے، اس لیے ccam ان فائلوں کو نہیں چھیڑتا جو اس نے نہیں بنائیں۔

جب آپ شارٹ کٹ چلاتے ہیں تو ccam کو صرف وہ نام چاہیے جو آپ نے ٹائپ کیا:

<!-- motion:shortcuts-link -->

شیل `PATH` پر `ccwork` تلاش کرتا ہے، کرنل لنک کے پیچھے جا کر ccam بائنری شروع کرتا ہے، اور `argv[0]` اب بھی `ccwork` ہوتا ہے۔ `app.Main` میں، `ccam` کے سوا کوئی بھی نام `config.toml` ایک بار لوڈ کرتا ہے اور `shim.Lookup` کو کال کرتا ہے، جو وہ اکاؤنٹ تلاش کرتا ہے جس کا `shortcut`، `filepath.Base(argv[0])` کے برابر ہو۔ پھر `launcher.LaunchWith` پہلے سے لوڈ شدہ کنفیگ کے ساتھ وہی مراحل اٹھاتا ہے جو [`ccam <name>`](/features/launch)، اور `argv[1:]` کو بغیر چھیڑے `claude` کو دیتا ہے۔

شیل alias یا wrapper اسکرپٹ کے بجائے symlink کیوں: symlink ہر شیل اور اسکرپٹس میں کام کرتا ہے، کوئی اضافی شیل شروع نہیں کرتا (تقریباً 1 سے 2 ms)، اور لانچ اب بھی ایک ہی فائل پڑھتا ہے۔ `scripts/perf-check.sh`، `hyperfine` سے شارٹ کٹ کو اسی 10 ms (p95) کی حد کے مقابلے میں ناپتا ہے جو `ccam <name>` کی ہے۔

## حفاظت اور خاص صورتیں

- **ایسا لنک جو کسی اکاؤنٹ کا نہ ہو۔** اگر `~/.local/bin` میں ccam کے کسی لنک کا نام کوئی اکاؤنٹ استعمال نہ کرتا ہو، مثلاً جب آپ نے `config.toml` میں `shortcut` ہاتھ سے بدلا ہو، تو ccam ڈیفالٹ اکاؤنٹ پر نہیں جاتا۔ یہ `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` کے ساتھ رک جاتا ہے۔ `ccam shortcut` سے نام دوبارہ مقرر کریں، یا لنک حذف کر دیں۔
- **موجودہ فائلیں۔** ccam کبھی ایسی فائل نہیں بدلتا جو اس نے نہ بنائی ہو۔ `ccam shims install`، `<path> already exists and is not a ccam link, so it was left alone` کے ساتھ رک جاتا ہے۔
- **منتقل کی گئی ccam بائنری۔** تب پرانے لنکس ایسے راستے کی طرف اشارہ کرتے ہیں جو اب موجود نہیں، اور ccam انہیں اپنا شمار نہیں کرتا۔ انہیں حذف کریں، پھر `ccam shims install` چلائیں۔
- **نام بدلنا اور ہٹانا۔** `ccam rename` شارٹ کٹ کا نام برقرار رکھتا ہے، اور لنک کام کرتا رہتا ہے، کیونکہ وہ ccam کی طرف اشارہ کرتا ہے اور ccam نام `config.toml` میں تلاش کرتا ہے۔ `ccam remove` اکاؤنٹ کا لنک حذف کرتا ہے۔ `ccam shortcut <account> <new>` نیا لنک بنانے سے پہلے پرانا حذف کرتا ہے۔
- **PATH۔** جب ccam لنک بناتا ہے اور `~/.local/bin` آپ کے `PATH` پر نہ ہو، تو یہ `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` پرنٹ کرتا ہے اور اس کے بعد شامل کرنے کے لیے `export PATH=...` والی لائن۔
- **بعد میں جانچنا۔** `ccam doctor`، `missing shortcuts in /home/me/.local/bin: ccwork` کی اطلاع دیتا ہے اور فولڈر `PATH` پر نہ ہو تو انتباہ دیتا ہے۔ `ccam doctor --fix` غائب لنکس بنا دیتا ہے۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [اکاؤنٹس کا انتظام](/features/accounts)
- [کمانڈز کا حوالہ](/docs/commands#shortcuts)
- [ccam انسٹال کریں](/docs/install)
