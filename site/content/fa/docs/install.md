---
title: نصب ccam
metaTitle: نصب ccam روی Linux، macOS و Windows
description: ccam را روی Linux، macOS یا Windows با یک فرمان، بسته .deb یا از کد منبع نصب کنید و با ccam update به‌روزش کنید. ابتدا باید Claude Code نصب شده باشد.
updated: 2026-10-21
order: 1
translation: machine
sourceHash: 1e3a53623de0
---

این فرمان را در یک ترمینال روی Linux یا macOS اجرا کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکریپت نسخه مناسب سیستم شما را دانلود می‌کند، checksum SHA-256 آن را بررسی می‌کند و `ccam` را در `~/.local/bin` قرار می‌دهد؛ ccam فرمان‌های میان‌بری را هم که بعداً انتخاب می‌کنید، مانند `ccwork`، همان‌جا می‌سازد. اگر آن پوشه در `PATH` شما نباشد، اسکریپت خطی را که باید به `~/.bashrc` یا `~/.zshrc` اضافه کنید چاپ می‌کند.

در Windows، این را در PowerShell اجرا کنید:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

این فرمان `ccam_windows_<arch>.zip` را دانلود می‌کند، چک‌سام SHA-256 آن را بررسی می‌کند، `ccam.exe` را در `%USERPROFILE%\.local\bin` می‌گذارد، همان پوشه‌ای که نصب‌کننده Claude Code فایل `claude.exe` را در آن قرار می‌دهد، و اگر آن پوشه هنوز در `PATH` کاربری شما نباشد، آن را اضافه می‌کند. پس از آن یک ترمینال جدید باز کنید.

ccam خودش Claude Code را برای شما اجرا می‌کند، پس ابتدا [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) را نصب کنید و مطمئن شوید که `claude --version` کار می‌کند.

## سیستم‌های پشتیبانی‌شده

- Linux روی x86-64 و ARM64: Ubuntu، Debian، Fedora، Arch، Alpine و توزیع‌های دیگر. فایل اجرایی ایستا است و به هیچ وابستگی‌ای نیاز ندارد.
- macOS روی Apple silicon و Intel.
- Windows 10 و 11 روی x86-64 و ARM64، در PowerShell، Windows Terminal یا Git Bash. پوشه‌های مشترک برای همه کاربران کار می‌کنند. فایل‌های مشترک (`settings.json`، `CLAUDE.md`، `history.jsonl`) به Developer Mode نیاز دارند (Settings > System > For developers)، چون Windows فقط در این حالت، یا به مدیر سیستم، اجازه ساختن پیوند به فایل‌ها را می‌دهد؛ [اشتراک در Windows](/features/share#on-windows) را ببینید.
- داخل WSL، ccam مانند هر سیستم Linux دیگری کار می‌کند.

## روش‌های دیگر نصب

**بسته Debian و Ubuntu.** فایل `ccam_<version>_linux_amd64.deb` را از [صفحه نسخه‌ها](https://github.com/jhyunwoo/ccam/releases) دانلود و نصب کنید:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

این بسته `/usr/bin/ccam` را نصب می‌کند.

**از کد منبع.** با Go 1.26 یا جدیدتر:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**یک نسخه مشخص.** متغیر `CCAM_VERSION` را تنظیم کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## تکمیل خودکار در شل

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

تکمیل خودکار همه فرمان‌ها و گزینه‌های `ccam` را پوشش می‌دهد، از جمله نام حساب‌ها.

## پلاگین Claude Code

پلاگین ccam به خود Claude Code یاد می‌دهد که از ccam استفاده کند: بررسی کند کدام حساب ظرفیت دارد، با `ccam pick` کارگرهای موازی `claude -p` را روی حساب‌های مختلف آغاز کند و با `ccam ext` توضیح دهد چرا یک مهارت یا سرور MCP در دسترس نیست. آن را درون Claude Code نصب کنید:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

این پلاگین یک مهارت، فرمان‌های `/ccam:usage`، `/ccam:ext` و `/ccam:switch` و دو هوک اضافه می‌کند: وقتی یک نشست در سقف مصرف متوقف شود یک اعلان در ترمینال دریافت می‌کنید و `/ccam:switch [account]` گفتگو را به حساب دیگری می‌سپارد ([جزئیات](/features/next#inside-claude-code)). هوک‌ها `ccam` را اجرا می‌کنند، پس باید در `PATH`ی باشد که Claude Code می‌بیند. این مهارت هرگز فرمان‌هایی را که وارد حساب می‌شوند یا فایل‌ها را جابه‌جا می‌کنند اجرا نمی‌کند؛ به جای آن، آن‌ها را پیشنهاد می‌دهد.

## بررسی نصب

```sh
ccam version
ccam doctor
```

`ccam doctor` بررسی می‌کند که Claude Code پیدا شود و هر پوشه حساب، پیوند و ورود درست باشد. سپس به [شروع سریع](/docs/quick-start) بروید.

## به‌روزرسانی

```sh
ccam update --check
ccam update
```

`ccam update --check` فقط به شما می‌گوید که آیا نسخه جدیدتری وجود دارد یا نه. `ccam update` نسخه مناسب سیستم شما را دانلود می‌کند، checksum SHA-256 آن را با فایل `checksums.txt` همان نسخه بررسی می‌کند و سپس فایل `ccam` را در یک مرحله جایگزین می‌کند. اگر checksum مطابقت نداشته باشد، هیچ چیزی تغییر نمی‌کند. میان‌برهایی مانند `ccwork` پیوندهایی به ccam هستند، پس بلافاصله نسخه جدید را اجرا می‌کنند.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

داشبورد (`ccam` در ترمینال) هم هر بار که باز می‌شود دنبال نسخه جدید می‌گردد و اگر باشد، آن را پیشنهاد می‌کند: `U` را بزنید و تأیید کنید. هیچ فرمان دیگری دنبال آن نمی‌گردد و ccam هرگز بدون پرسیدن به‌روز نمی‌شود. برای نصب یک نسخه مشخص، حتی نسخه‌ای قدیمی‌تر، گزینه `--version` را بدهید. `ccam --dry-run update` بدون اینکه چیزی را تغییر دهد نشان می‌دهد چه چیزی دانلود و جایگزین می‌شود.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

اجرای دوباره اسکریپت نصب هم نسخه موجود در `~/.local/bin` را به‌روز می‌کند. اگر ccam را با `go install` ساخته‌اید، همان فرمان را دوباره اجرا کنید.

**Windows.** `ccam update` فایل `.zip` مخصوص Windows را دانلود می‌کند. در Windows نمی‌توان برنامه‌ای را که در حال اجراست جایگزین کرد، پس نام فایل اجرایی قدیمی به `ccam.exe.old` تغییر می‌کند و به‌روزرسانی بعدی آن را پاک می‌کند.

**بسته Debian و Ubuntu.** `ccam update` فایل `/usr/bin/ccam` را که apt نصب کرده است جایگزین نمی‌کند و به‌جای آن این مراحل را چاپ می‌کند. فایل جدید `ccam_<version>_linux_amd64.deb` را از [صفحه نسخه‌ها](https://github.com/jhyunwoo/ccam/releases) دانلود کنید و روی نسخه قبلی نصب کنید:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## حذف نصب

ابتدا، تا وقتی ccam هنوز نصب است، میان‌برهای حساب‌ها را حذف کنید. `ccam list` آن‌ها را در ستون Shortcut نشان می‌دهد؛ هر کدام را با `ccam shortcut <account> --remove` حذف کنید، یا خودتان آن پیوندها را از `~/.local/bin` پاک کنید. سپس فرمان را حذف کنید و اگر دیگر به حساب‌ها نیاز ندارید، پوشه‌های داده را هم حذف کنید:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

در Windows، در PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

پاک کردن `~/.local/share/ccam` اطلاعات ورود و تاریخچه گفتگوی ذخیره‌شده در آن را هم پاک می‌کند. پوشه `~/.claude` که آن را [وارد کرده‌اید](/docs/commands) در جای خود باقی می‌ماند.
