---
title: "ccam انسٹال کریں"
metaTitle: "ccam کو Linux، macOS اور Windows پر انسٹال کریں"
description: "ccam کو Linux، macOS یا Windows پر ایک کمانڈ، .deb پیکیج یا سورس سے انسٹال کریں اور ccam update سے اپ ڈیٹ کریں۔ اس سے پہلے Claude Code انسٹال ہونا ضروری ہے۔"
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "a7f9dac24f79"
---

Linux یا macOS پر ٹرمینل میں یہ کمانڈ چلائیں:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکرپٹ آپ کے سسٹم کے لیے ریلیز ڈاؤن لوڈ کرتا ہے، اس کا SHA-256 چیک سم جانچتا ہے اور `ccam` کو `~/.local/bin` میں رکھتا ہے، جہاں ccam وہ شارٹ کٹ کمانڈز بھی بناتا ہے جو آپ بعد میں چنیں، جیسے `ccwork`۔ اگر وہ فولڈر آپ کے `PATH` میں نہیں ہے تو اسکرپٹ وہ لائن دکھاتا ہے جو `~/.bashrc` یا `~/.zshrc` میں شامل کرنی ہے۔

Windows پر یہ PowerShell میں چلائیں:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

یہ `ccam_windows_<arch>.zip` ڈاؤن لوڈ کرتا ہے، اس کا SHA-256 چیک سم جانچتا ہے، `ccam.exe` کو `%USERPROFILE%\.local\bin` میں رکھتا ہے، وہی فولڈر جس میں Claude Code کا انسٹالر `claude.exe` رکھتا ہے، اور اگر وہ فولڈر ابھی آپ کے یوزر `PATH` میں نہ ہو تو اسے شامل کر دیتا ہے، اور `%USERPROFILE%\bin` جیسے اندراجات کو جوں کا توں رہنے دیتا ہے۔ اس کے بعد نیا ٹرمینل کھولیں۔ اسے دوبارہ چلانے سے ccam بدل جاتا ہے، تب بھی جب کوئی Claude Code سیشن ابھی پرانی بائنری سے چل رہا ہو (اس بائنری کا نام بدل کر `ccam.exe.old` کر دیا جاتا ہے، یا `ccam.exe.old-<number>` جب اس سے پرانی کوئی بائنری ابھی استعمال میں ہو)، اور جو شارٹ کٹس hard links ہیں انہیں نیا کر دیتا ہے: `OK 2 shortcuts now run the new ccam`۔

ccam آپ کے لیے Claude Code شروع کرتا ہے، اس لیے پہلے [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) انسٹال کریں اور یقینی بنائیں کہ `claude --version` کام کرتا ہے۔

## معاون سسٹمز

- x86-64 اور ARM64 پر Linux: Ubuntu، Debian، Fedora، Arch، Alpine اور دیگر ڈسٹری بیوشنز۔ بائنری اسٹیٹک ہے اور اسے کسی اور چیز کی ضرورت نہیں۔
- Apple silicon اور Intel پر macOS۔
- x86-64 اور ARM64 پر Windows 10 اور 11، PowerShell، Windows Terminal یا Git Bash میں۔ شیئر کیے گئے فولڈرز ہر یوزر کے لیے کام کرتے ہیں۔ شیئر کی گئی فائلوں (`settings.json`، `CLAUDE.md`، `history.jsonl`) کے لیے Developer Mode درکار ہے (Settings > System > For developers)، کیونکہ Windows فائلوں کے لنکس بنانے کی اجازت صرف اسی موڈ میں، یا ایڈمنسٹریٹر کو، دیتا ہے؛ دیکھیں [Windows پر شیئرنگ](/features/share#on-windows)۔ اس کے بغیر پہلا `ccam add` گروپ `default` صرف MCP سرورز اور ایکسٹینشنز کے ساتھ بناتا ہے، اور بتاتا ہے کہ Developer Mode آن ہونے کے بعد باقی چیزیں کیسے شامل کریں۔
- WSL کے اندر ccam کسی بھی Linux سسٹم کی طرح کام کرتا ہے۔

## انسٹال کرنے کے دوسرے طریقے

**Debian اور Ubuntu پیکیج۔** [ریلیزز کے صفحے](https://github.com/jhyunwoo/ccam/releases) سے `ccam_<version>_linux_amd64.deb` ڈاؤن لوڈ کریں اور اسے انسٹال کریں:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

پیکیج `/usr/bin/ccam` انسٹال کرتا ہے۔

**سورس سے۔** Go 1.26 یا اس سے نئے ورژن کے ساتھ:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

کسی ریلیز ٹیگ (`@latest` یا `@vX.Y.Z`) سے بنی بائنری وہی ورژن بتاتی ہے، اس لیے اس کے لیے `ccam update` کام کرتا ہے؛ `@main` سے، یا کسی ٹیگ سے آگے بڑھ چکے لوکل checkout سے بنی بائنری `dev` بتاتی ہے۔

**کوئی مخصوص ورژن۔** `CCAM_VERSION` سیٹ کریں:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## شیل کمپلیشن

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

کمپلیشن `ccam` کی تمام کمانڈز اور آپشنز کے لیے کام کرتی ہے، اکاؤنٹس کے ناموں سمیت۔

## Claude Code پلگ اِن

ccam پلگ اِن خود Claude Code کو ccam استعمال کرنا سکھاتا ہے: یہ دیکھنا کہ کس اکاؤنٹ میں گنجائش ہے، `ccam pick` سے مختلف اکاؤنٹس پر متوازی `claude -p` ورکرز شروع کرنا، اور `ccam ext` سے یہ بتانا کہ کوئی اسکل یا MCP سرور کیوں غائب ہے۔ اسے Claude Code کے اندر انسٹال کریں:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

یہ ایک اسکل، کمانڈز `/ccam:usage`، `/ccam:ext` اور `/ccam:switch`، اور دو ہُکس شامل کرتا ہے: جب کوئی سیشن استعمال کی حد پر رک جائے تو آپ کو ٹرمینل میں اطلاع ملتی ہے (اور Linux اور macOS پر، [`auto_switch`](/features/next#auto-switch) آن ہو تو، ccam سیشن کو کسی ایسے دوسرے اکاؤنٹ میں منتقل کر دیتا ہے جس میں گنجائش ہو)، اور `/ccam:switch [account]` گفتگو کسی دوسرے اکاؤنٹ کو سونپ دیتا ہے ([تفصیلات](/features/next#inside-claude-code))۔ ہُکس `ccam` چلاتے ہیں، اس لیے اسے اس `PATH` پر ہونا چاہیے جو Claude Code دیکھتا ہے۔ اسکل کبھی ایسی کمانڈز نہیں چلاتی جو لاگ اِن کریں یا فائلیں منتقل کریں؛ اس کے بجائے انہیں تجویز کرتی ہے۔

## انسٹالیشن چیک کریں

```sh
ccam version
ccam doctor
```

`ccam doctor` چیک کرتا ہے کہ Claude Code مل رہا ہے اور ہر اکاؤنٹ فولڈر، لنک اور لاگ اِن درست حالت میں ہے۔ اس کے بعد [فوری آغاز](/docs/quick-start) پر جائیں۔

## اپ ڈیٹ کرنا

```sh
ccam update --check
ccam update
```

`ccam update --check` صرف یہ بتاتا ہے کہ کوئی نئی ریلیز موجود ہے یا نہیں۔ `ccam update` آپ کے سسٹم کے لیے ریلیز ڈاؤن لوڈ کرتا ہے، اس کا SHA-256 چیک سم ریلیز کی `checksums.txt` فائل سے ملاتا ہے اور پھر `ccam` فائل کو ایک ہی قدم میں بدل دیتا ہے۔ اگر چیک سم میل نہ کھائے تو کچھ بھی تبدیل نہیں ہوتا۔ `ccwork` جیسے شارٹ کٹس ccam کے لنکس ہیں، اس لیے وہ فوراً نیا ورژن چلاتے ہیں۔

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ڈیش بورڈ (ٹرمینل پر `ccam`) بھی ہر بار کھلنے پر نئی ریلیز ڈھونڈتا ہے اور اگر کوئی ہو تو اسے پیش کرتا ہے: `U` دبائیں اور تصدیق کریں۔ کوئی اور کمانڈ نہیں ڈھونڈتی، اور ccam پوچھے بغیر کبھی اپ ڈیٹ نہیں ہوتا۔ کوئی مخصوص ریلیز، پرانی ریلیز سمیت، انسٹال کرنے کے لیے `--version` دیں۔ `ccam --dry-run update` کچھ بھی تبدیل کیے بغیر دکھاتا ہے کہ کیا ڈاؤن لوڈ ہوگا اور کیا بدلا جائے گا۔

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

انسٹال اسکرپٹ دوبارہ چلانے سے بھی `~/.local/bin` والی کاپی اپ ڈیٹ ہو جاتی ہے۔ اگر آپ نے ccam کو کسی ریلیز ٹیگ سے `go install` کے ذریعے بنایا تھا تو `ccam update` بھی کام کرتا ہے، یا وہی کمانڈ دوبارہ چلائیں۔

**Windows۔** `ccam update` Windows کے لیے `.zip` ڈاؤن لوڈ کرتا ہے۔ Windows پر چلتے ہوئے پروگرام کو بدلا نہیں جا سکتا، اس لیے پرانی بائنری کا نام بدل کر `ccam.exe.old` کر دیا جاتا ہے، جسے اگلی اپ ڈیٹ حذف کر دیتی ہے؛ جب کوئی سیشن ابھی کسی پچھلی `ccam.exe.old` سے چل رہا ہو تو نام `ccam.exe.old-<time>` ہوتا ہے۔ جو شارٹ کٹس hard links ہیں وہ اس کے فوراً بعد نئے کر دیے جاتے ہیں۔

**Debian اور Ubuntu پیکیج۔** apt کے ذریعے انسٹال کی گئی `/usr/bin/ccam` کو `ccam update` نہیں بدلتا بلکہ اس کے بجائے یہ مراحل دکھاتا ہے۔ [ریلیزز کے صفحے](https://github.com/jhyunwoo/ccam/releases) سے نئی `ccam_<version>_linux_amd64.deb` ڈاؤن لوڈ کریں اور اسے پرانی کے اوپر انسٹال کریں:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## اَن انسٹال کرنا

پہلے، جب تک ccam انسٹال ہے، اکاؤنٹس کے شارٹ کٹس ہٹائیں۔ `ccam list` انہیں Shortcut کالم میں دکھاتا ہے؛ ہر ایک کو `ccam shortcut <account> --remove` سے ہٹائیں، یا وہ لنکس خود `~/.local/bin` سے حذف کریں۔ پھر کمانڈ ہٹائیں، اور اگر اکاؤنٹس کی ضرورت نہ رہی ہو تو ڈیٹا فولڈرز بھی:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows پر، PowerShell میں:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` حذف کرنے سے وہاں محفوظ لاگ اِن اور گفتگو کی ہسٹری بھی حذف ہو جاتی ہے۔ جو `~/.claude` فولڈر آپ نے [امپورٹ کیا تھا](/docs/commands) وہ اپنی جگہ پر رہتا ہے، لیکن جب تک وہ کسی گروپ میں ہے، جو فائلیں وہ شیئر کرتا ہے وہ `~/.local/share/ccam/groups/` میں رہتی ہیں۔ پہلے اس کے لیے `ccam remove <name>` چلائیں (یا `ccam group leave`): یہ اپنے گروپس سے نکل جاتا ہے اور اس کے پاس ہر اس چیز کی کاپی رہتی ہے جو اس نے شیئر کی تھی۔
