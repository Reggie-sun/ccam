---
title: نصب ccam
metaTitle: نصب ccam روی Linux و macOS
description: ccam را روی Linux یا macOS با یک فرمان، بسته .deb یا از کد منبع نصب کنید و با ccam update به‌روزش کنید. ابتدا باید Claude Code نصب شده باشد.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
---

این فرمان را در یک ترمینال روی Linux یا macOS اجرا کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکریپت نسخه مناسب سیستم شما را دانلود می‌کند، checksum SHA-256 آن را بررسی می‌کند و `ccam` را در `~/.local/bin` قرار می‌دهد؛ ccam فرمان‌های میان‌بری را هم که بعداً انتخاب می‌کنید، مانند `ccwork`، همان‌جا می‌سازد. اگر آن پوشه در `PATH` شما نباشد، اسکریپت خطی را که باید به `~/.bashrc` یا `~/.zshrc` اضافه کنید چاپ می‌کند.

ccam خودش Claude Code را برای شما اجرا می‌کند، پس ابتدا [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) را نصب کنید و مطمئن شوید که `claude --version` کار می‌کند.

## سیستم‌های پشتیبانی‌شده

- Linux روی x86-64 و ARM64: Ubuntu، Debian، Fedora، Arch، Alpine و توزیع‌های دیگر. فایل اجرایی ایستا است و به هیچ وابستگی‌ای نیاز ندارد.
- macOS روی Apple silicon و Intel.
- Windows مستقیماً پشتیبانی نمی‌شود. داخل WSL، ccam مانند هر سیستم Linux دیگری کار می‌کند.

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

## تکمیل خودکار در شل

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

تکمیل خودکار همه فرمان‌ها و گزینه‌های `ccam` را پوشش می‌دهد، از جمله نام حساب‌ها.

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

ccam هرگز خودش به دنبال به‌روزرسانی نمی‌گردد و فقط وقتی این فرمان را اجرا کنید سراغ نسخه جدید می‌رود. برای نصب یک نسخه مشخص، حتی نسخه‌ای قدیمی‌تر، گزینه `--version` را بدهید. `ccam --dry-run update` بدون اینکه چیزی را تغییر دهد نشان می‌دهد چه چیزی دانلود و جایگزین می‌شود.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

اجرای دوباره اسکریپت نصب هم نسخه موجود در `~/.local/bin` را به‌روز می‌کند. اگر ccam را با `go install` ساخته‌اید، همان فرمان را دوباره اجرا کنید.

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

پاک کردن `~/.local/share/ccam` اطلاعات ورود و تاریخچه گفتگوی ذخیره‌شده در آن را هم پاک می‌کند. پوشه `~/.claude` که آن را [وارد کرده‌اید](/docs/commands) در جای خود باقی می‌ماند.
