---
title: نصب ccam
description: نصب ccam روی Linux یا macOS با یک فرمان، از بسته .deb یا از کد منبع. ابتدا باید Claude Code نصب شده باشد.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

این فرمان را در یک ترمینال روی Linux یا macOS اجرا کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکریپت نسخه منتشرشده مناسب سیستم شما را دانلود می‌کند، checksum از نوع SHA-256 آن را بررسی می‌کند و دو فرمان را در `~/.local/bin` قرار می‌دهد: `ccam` و نام مستعار کوتاه آن `ccs`. اگر این پوشه در `PATH` شما نباشد، اسکریپت خطی را که باید به `~/.bashrc` یا `~/.zshrc` اضافه کنید چاپ می‌کند.

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

این بسته `/usr/bin/ccam` و پیوند `/usr/bin/ccs` را نصب می‌کند.

**از کد منبع.** با Go 1.26 یا جدیدتر:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

تکمیل خودکار هم `ccam` و هم `ccs` را پوشش می‌دهد، از جمله نام حساب‌ها.

## بررسی نصب

```sh
ccam version
ccam doctor
```

`ccam doctor` بررسی می‌کند که Claude Code پیدا شود و هر پوشه حساب، پیوند و ورود درست باشد. سپس به [شروع سریع](/docs/quick-start) بروید.

## حذف نصب

دو فرمان و میان‌برهای حساب‌ها را حذف کنید، سپس اگر دیگر به حساب‌ها نیاز ندارید، پوشه‌های داده را هم حذف کنید:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

پاک کردن `~/.local/share/ccam` اطلاعات ورود و تاریخچه گفتگوی ذخیره‌شده در آن را هم پاک می‌کند. پوشه `~/.claude` که آن را [وارد کرده‌اید](/docs/commands) در جای خود باقی می‌ماند.
