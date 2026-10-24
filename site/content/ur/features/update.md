---
title: "ccam اپ ڈیٹ کرنا"
description: "ccam update، SHA-256 چیک کے بعد ccam کو GitHub کی تازہ ترین ریلیز سے بدل دیتا ہے؛ --check صرف بتاتا ہے اور --version کسی مخصوص ریلیز ٹیگ کو انسٹال کرتا ہے۔"
updated: 2026-10-23
order: 10
translation: machine
sourceHash: "000d37de5f29"
---

`ccam update` ccam کی ایک ریلیز GitHub سے ڈاؤن لوڈ کرتا ہے، اسے ریلیز کی `checksums.txt` سے جانچتا ہے اور چلنے والی بائنری کو ایک ہی rename میں بدل دیتا ہے۔ اسے تب استعمال کریں جب آپ نے ccam انسٹال اسکرپٹ سے یا کسی آرکائیو سے انسٹال کیا ہو۔ [ڈیش بورڈ](/features/dashboard) بھی ہر بار کھلنے پر GitHub سے تازہ ترین ریلیز پوچھتا ہے اور اسے `U` کے ذریعے پیش کرتا ہے، جو پوچھنے کے بعد یہی مراحل چلاتا ہے۔ کوئی اور کمانڈ GitHub سے رابطہ نہیں کرتی، اس لیے وہ آف لائن بھی ویسے ہی کام کرتی ہیں۔

## استعمال

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

جب کوئی نئی ریلیز موجود نہ ہو تو ccam `ccam 0.4.0 is up to date.` پرنٹ کرتا ہے اور کچھ تبدیل نہیں کرتا۔ `--version` کے ساتھ وہی ورژن دینے پر جو آپ پہلے سے چلا رہے ہیں، یہ `ccam 0.4.0 is already installed.` پرنٹ کرتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--check` | صرف یہ بتاتا ہے کہ کوئی نئی ریلیز موجود ہے یا نہیں | بند |
| `--version <tag>` | یہ ریلیز ٹیگ انسٹال کرتا ہے، جیسے `v0.3.0`، جو موجودہ سے نیا یا پرانا ہو سکتا ہے | تازہ ترین |
| `--dry-run` | لاگو ہوتا ہے: تازہ ترین ٹیگ تلاش کرتا ہے (جب تک `--version` کوئی ٹیگ نہ بتائے)، بائنری پر وہی جانچ چلاتا ہے، اور ڈاؤن لوڈ کیے بغیر ڈاؤن لوڈ، چیک سم اور بدلنے کے راستے پرنٹ کرتا ہے | بند |

`--check` اور `--version` ایک ساتھ استعمال نہیں ہو سکتے۔ ٹیگ `vMAJOR.MINOR.PATCH` جیسا ہونا چاہیے، جس کے بعد اختیاری طور پر pre-release یا build کا لاحقہ ہو سکتا ہے؛ شروع کا `v` چھوڑا جا سکتا ہے۔ اس کے علاوہ کچھ بھی ناکام ہوتا ہے، مثلاً `--version takes a release tag such as v0.3.0, not "1.2"`۔

انوائرنمنٹ ویری ایبل `CCAM_RELEASES_URL`، کسی مرر یا ٹیسٹس کے لیے، `https://github.com/jhyunwoo/ccam/releases` کی جگہ لے لیتا ہے۔ اس کی ترتیب بھی ویسی ہی ہونی چاہیے: `<url>/latest` کو `<url>/tag/<tag>` پر ری ڈائریکٹ ہونا چاہیے، اور فائلیں `<url>/download/<tag>/<asset>` پر ہونی چاہئیں۔

## یہ کیسے کام کرتا ہے

کوڈ `internal/update` میں ہے، جسے صرف `ccam update` امپورٹ کرتا ہے۔ Claude Code شروع کرنے والا راستہ اسے کبھی کال نہیں کرتا۔

1. `cli.updateCmd`، `os.Executable` سے چلنے والی بائنری تلاش کرتا ہے اور `filepath.EvalSymlinks` سے symlinks حل کرتا ہے۔ یہی اصل راستہ بدلا جاتا ہے۔
2. `--version` کے بغیر، `Updater.Latest` ری ڈائریکٹ پر عمل کیے بغیر `<releases>/latest` کی درخواست کرتا ہے، اور ٹیگ `Location` ہیڈر سے پڑھتا ہے (`…/tag/v0.5.0`)۔ اس میں GitHub API کے بجائے ریلیزز کا صفحہ استعمال ہوتا ہے، اس لیے API کی 60 درخواستیں فی گھنٹہ کی حد لاگو نہیں ہوتی۔ اس تلاش کے لیے 10 سیکنڈ ہوتے ہیں۔
3. `parseVersion` ٹیگ کو ایک سخت پیٹرن سے جانچتا ہے۔ ٹیگ ڈاؤن لوڈ URLs اور پرنٹ کی گئی شیل کمانڈ میں جاتا ہے، اس لیے `/`، `%`، اسپیس یا شیل کے حروف والی کوئی چیز آگے نہیں جاتی۔ `compareVersions` اسے چلنے والے ورژن سے موازنہ کرتا ہے؛ pre-release اپنی ریلیز سے پہلے شمار ہوتی ہے۔
4. `managed()` اور `writable()` طے کرتے ہیں کہ ccam بائنری بدل بھی سکتا ہے یا نہیں (نیچے دیکھیں)۔
5. `Updater.Fetch`، `ccam_<os>_<arch>.tar.gz` (Windows پر `ccam_windows_<arch>.zip`) اور `checksums.txt` ایک ساتھ ڈاؤن لوڈ کرتا ہے، دونوں کے لیے 2 منٹ۔ ہر ڈاؤن لوڈ کی حد 256 MB ہے۔
6. کچھ اور ہونے سے پہلے آرکائیو کا SHA-256، `checksums.txt` میں اس کی لائن سے ملنا چاہیے۔ لائن نہ ہو یا چیک سم مختلف ہو تو اپ ڈیٹ `nothing was changed` کے ساتھ رک جاتا ہے۔
7. `extract()` ٹار پڑھتا ہے اور `ccam` نام کی پہلی عام فائل لیتا ہے (`extractZip`، zip سے `ccam.exe` لیتا ہے)۔ فولڈرز، لنکس اور دوسری اینٹریز نظر انداز ہوتی ہیں، اور بائنری کی حد بھی 256 MB ہے۔
8. `Updater.Install` نئی بائنری `fsutil.AtomicWrite` سے لکھتا ہے: بائنری کے ساتھ ایک عارضی فائل `.ccam.tmp*`، پرانی فائل کی اجازتوں کے ساتھ، جو ڈسک پر محفوظ ہو کر پرانی فائل کی جگہ rename ہو جاتی ہے۔

ٹیگ تلاش کرنے سے لے کر rename تک پورا عمل:

<!-- motion:update-pipeline -->

آرکائیو اور `checksums.txt` ساتھ ساتھ ڈاؤن لوڈ ہوتے ہیں، اور چیک سم ملنے تک کچھ نہیں لکھا جاتا۔ rename ایٹامک ہے، اس لیے کریش یا بھری ہوئی ڈسک کے بعد یا تو پرانی بائنری رہتی ہے یا نئی، کبھی آدھی نہیں۔ `~/.local/bin/ccwork` جیسی شارٹ کٹ کمانڈز ccam بائنری کے symlinks ہیں، اس لیے وہ بغیر کسی تبدیلی کے نیا ورژن چلاتی ہیں۔

Windows پر چلتے ہوئے پروگرام کی فائل نہ بدلی جا سکتی ہے نہ حذف کی جا سکتی ہے، صرف اس کا نام بدلا جا سکتا ہے۔ `Updater.Install` نئی بائنری `ccam.exe.new` میں لکھتا ہے، چلتی ہوئی `ccam.exe` کا نام بدل کر `ccam.exe.old` کرتا ہے، اور نئی فائل کا نام `ccam.exe` کر دیتا ہے۔ اگر آخری قدم ناکام ہو تو پرانی بائنری کا نام واپس بدل دیا جاتا ہے۔ اگلی اپ ڈیٹ وہ `ccam.exe.old` فائلیں حذف کر دیتی ہے جنہیں وہ حذف کر سکے؛ جب کوئی سیشن ابھی کسی پچھلی `ccam.exe.old` سے چل رہا ہو تو چلتی ہوئی بائنری اس کے بجائے `ccam.exe.old-<time>` بن جاتی ہے۔ جو شارٹ کٹس symbolic links کے بجائے hard links ہیں ([شارٹ کٹ کمانڈز](/features/shortcuts)) وہ فوراً بعد نئے کیے جاتے ہیں، اس لیے وہ بھی نیا ورژن چلاتے ہیں؛ اس میں ccam کے بنائے ہوئے کسی بھی پرانے ccam کے hard links شامل ہیں، نہ کہ صرف ابھی بدلی گئی بائنری کے۔

### پیکیج سے انسٹال شدہ ccam

`.deb` سے انسٹال کیا گیا ccam، apt کی ملکیت ہے، اور اسے apt کے علم کے بغیر بدلنے سے پیکیج ڈیٹا بیس غلط ہو جائے گا۔ ڈاؤن لوڈ سے پہلے ccam دیکھتا ہے کہ بائنری کس کی ملکیت ہے:

<!-- motion:update-managed -->

1. `/usr/bin`، `/bin`، `/usr/sbin` یا `/sbin` میں موجود بائنری پیکیج سے منظم شمار ہوتی ہے۔
2. کسی اور جگہ، اگر `dpkg` انسٹال ہو اور `dpkg -S <path>` کامیاب ہو، تو وہ بھی پیکیج سے منظم شمار ہوتی ہے۔
3. ورنہ ccam بائنری کے ساتھ ایک عارضی `.ccam-update-*` فائل بنا کر ہٹاتا ہے تاکہ یقین ہو جائے کہ وہ وہاں لکھ سکتا ہے۔

پیکیج سے منظم بائنری کے لیے ccam رک جاتا ہے اور apt کے مراحل پرنٹ کرتا ہے:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux کے علاوہ دوسرے سسٹمز پر صرف انسٹال اسکرپٹ والی لائن پرنٹ ہوتی ہے: Windows پر `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`۔ جب فولڈر میں لکھا نہ جا سکے تو ccam انسٹال اسکرپٹ دوبارہ چلانے کا مشورہ دیتا ہے، جو ccam کو `~/.local/bin` میں انسٹال کرتا ہے۔

## حفاظت اور خاص صورتیں

- ڈیولپمنٹ بلڈ اپنا ورژن `dev` بتاتا ہے اور اس کے موازنے کے لیے کوئی ریلیز نہیں ہوتی، اس لیے `ccam update` اور `ccam update --check` اسے مسترد کر دیتے ہیں۔ `ccam update --version vX.Y.Z` اس کے اوپر ایک ریلیز انسٹال کر دیتا ہے۔ `go install github.com/jhyunwoo/ccam/cmd/ccam@vX.Y.Z` (یا `@latest`) سے بنا بلڈ وہی ریلیز بتاتا ہے، اس لیے `ccam update` اور ڈیش بورڈ کی پیشکش اس کے لیے کام کرتی ہیں؛ `@main` یا کسی ٹیگ سے آگے کا لوکل بلڈ `dev` ہی رہتا ہے۔
- ناکام ڈاؤن لوڈ، آپ کے پلیٹ فارم کی فائل کا نہ ہونا، یا چیک سم کا نہ ملنا انسٹال شدہ بائنری کو ویسا ہی چھوڑ دیتا ہے۔
- جو ریلیز موجود نہ ہو وہ `there is no release <tag>` کے ساتھ ناکام ہوتی ہے۔ جب GitHub HTTP 403 یا 429 سے انکار کرے تو ccam بتاتا ہے کہ یہ غالباً ریٹ لمٹ ہے اور بعد میں دوبارہ کوشش کرنے کو کہتا ہے۔
- `--check` صرف پڑھتا ہے۔ یہ نہیں جانچتا کہ بائنری پیکیج سے منظم ہے یا اس میں لکھا جا سکتا ہے۔
- پرانے ٹیگ کے ساتھ `--version` ہی پچھلے ورژن پر واپس جانے کا طریقہ ہے۔

## متعلقہ

- [انسٹال کرنا](/docs/install)
- [کمانڈز](/docs/commands)
- [عام سوالات](/docs/faq)
