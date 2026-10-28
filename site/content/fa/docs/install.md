---
title: نصب ccam
metaTitle: نصب ccam روی Linux، macOS و Windows
description: ccam را روی Linux، macOS یا Windows با یک فرمان، بسته .deb یا go install نصب کنید و با ccam update به‌روزش کنید. ابتدا باید Claude Code نصب شده باشد.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: 193df30c9d9a
---

ccam خودش Claude Code را برای شما اجرا می‌کند، پس ابتدا [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) را نصب کنید و مطمئن شوید که `claude --version` کار می‌کند.

این فرمان را در یک ترمینال روی Linux یا macOS اجرا کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکریپت نسخه مناسب سیستم شما را دانلود می‌کند، checksum SHA-256 آن را بررسی می‌کند و `ccam` را در `~/.local/bin` قرار می‌دهد؛ ccam فرمان‌های میان‌بری را هم که بعداً انتخاب می‌کنید، مانند `ccwork`، همان‌جا می‌سازد. اگر آن پوشه در `PATH` شما نباشد، اسکریپت خطی را که باید به `~/.bashrc` یا `~/.zshrc` اضافه کنید چاپ می‌کند.

در Windows، این را در PowerShell اجرا کنید:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

این فرمان `ccam_windows_<arch>.zip` را دانلود می‌کند، چک‌سام SHA-256 آن را بررسی می‌کند و `ccam.exe` را در `%USERPROFILE%\.local\bin` می‌گذارد، همان پوشه‌ای که نصب‌کننده Claude Code فایل `claude.exe` را در آن قرار می‌دهد. اگر آن پوشه هنوز در `PATH` کاربری شما نباشد، اسکریپت آن را اضافه می‌کند و ورودی‌هایی مانند `%USERPROFILE%\bin` را همان‌طور که هستند نگه می‌دارد. پس از آن یک ترمینال جدید باز کنید.

## سیستم‌های پشتیبانی‌شده

- Linux روی x86-64 و ARM64: Ubuntu، Debian، Fedora، Arch، Alpine و توزیع‌های دیگر. فایل اجرایی ایستا است و به هیچ وابستگی‌ای نیاز ندارد.
- macOS روی Apple silicon و Intel.
- Windows 10 و 11 روی x86-64 و ARM64، در PowerShell، Windows Terminal یا Git Bash:
  - پوشه‌های مشترک بدون Developer Mode کار می‌کنند.
  - فایل‌های مشترک (`settings.json`، `CLAUDE.md`، `history.jsonl`) به Developer Mode نیاز دارند (Settings > System > For developers)، چون Windows فقط در این حالت، یا به مدیر سیستم، اجازه ساختن پیوند به فایل‌ها را می‌دهد؛ [اشتراک در Windows](/features/share#on-windows) را ببینید.
  - بدون Developer Mode، نخستین `ccam add` گروه `default` را فقط با سرورهای MCP و افزونه‌ها می‌سازد و می‌گوید پس از روشن شدن Developer Mode بقیه را چگونه اضافه کنید.
- داخل WSL، ccam همان‌طور کار می‌کند که روی هر سیستم Linux دیگری.

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

اگر یک برچسب نسخه (`@latest` یا `@vX.Y.Z`) را نصب کنید، ccam همان نسخه را گزارش می‌کند و `ccam update` و پیشنهاد به‌روزرسانی داشبورد کار می‌کنند. ساختی از `@main`، یا از یک checkout محلی که پس از یک برچسب commit دارد، `dev` را گزارش می‌کند و `ccam update` آن را رد می‌کند، مگر اینکه `--version` را بدهید.

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

این پلاگین این‌ها را اضافه می‌کند:

- یک مهارت که فرمان‌هایی را که وارد حساب می‌شوند یا فایل‌ها را جابه‌جا می‌کنند پیشنهاد می‌دهد اما هرگز اجرایشان نمی‌کند؛
- فرمان‌های `/ccam:usage`، `/ccam:ext` و `/ccam:switch`؛
- دو هوک ([جزئیات](/features/next#inside-claude-code)). یکی وقتی یک نشست در سقف مصرف متوقف شود یک اعلان در ترمینال نشان می‌دهد؛ در Linux و macOS، با روشن بودن [`auto_switch`](/features/next#auto-switch)، ccam سپس نشست را به حساب دیگری که ظرفیت دارد منتقل می‌کند. دیگری به `/ccam:switch [account]` امکان می‌دهد گفتگو را به حساب دیگری بسپارد.

هوک‌ها `ccam` را اجرا می‌کنند، پس باید در `PATH`ی باشد که Claude Code می‌بیند.

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

داشبورد (`ccam` بدون آرگومان، در ترمینال) هم هر بار که باز می‌شود نسخه جدید را بررسی می‌کند و اگر باشد، آن را پیشنهاد می‌کند: `U` را بزنید و تأیید کنید. هیچ فرمان دیگری به‌روزرسانی‌ها را بررسی نمی‌کند و ccam هرگز بدون پرسیدن به‌روز نمی‌شود. برای نصب یک نسخه مشخص، حتی نسخه‌ای قدیمی‌تر، گزینه `--version` را بدهید. `ccam --dry-run update` بدون اینکه چیزی را تغییر دهد نشان می‌دهد چه چیزی دانلود و جایگزین می‌شود.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

اجرای دوباره اسکریپت نصب هم نسخه موجود در `~/.local/bin` را به‌روز می‌کند. اگر ccam را با `go install` از یک برچسب نسخه ساخته‌اید، می‌توانید `ccam update` یا همان فرمان `go install` را دوباره اجرا کنید.

**Windows.** `ccam update` فایل `.zip` مخصوص Windows را دانلود می‌کند. در Windows نمی‌توان برنامه‌ای را که در حال اجراست جایگزین کرد، پس نام فایل اجرایی قدیمی به `ccam.exe.old` تغییر می‌کند و به‌روزرسانی بعدی آن را پاک می‌کند؛ وقتی یک نشست هنوز از یک `ccam.exe.old` قدیمی‌تر در حال اجراست، نام `ccam.exe.old-<time>` است. میان‌برهایی که hard link هستند بلافاصله پس از آن تازه می‌شوند. اجرای دوباره `install.ps1` هم به همین شکل کار می‌کند، حتی وقتی یک نشست Claude Code هنوز از فایل اجرایی قدیمی در حال اجراست، و وقتی میان‌برهای hard link را تازه می‌کند `OK 2 shortcuts now run the new ccam` را چاپ می‌کند.

**بسته Debian و Ubuntu.** `ccam update` فایل `/usr/bin/ccam` را که apt نصب کرده است جایگزین نمی‌کند و به‌جای آن این مراحل را چاپ می‌کند. فایل جدید `ccam_<version>_linux_amd64.deb` را از [صفحه نسخه‌ها](https://github.com/jhyunwoo/ccam/releases) دانلود کنید و روی نسخه قبلی نصب کنید:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## حذف نصب

`ccam uninstall` خود ccam و هر چیزی را که ساخته است حذف می‌کند. ابتدا فهرست چیزهایی را که حذف می‌شوند نشان می‌دهد و می‌پرسد:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- هر حساب همان‌طور حذف می‌شود که `ccam remove` حذف می‌کند. حسابی که در پوشه‌های ccam است ابتدا از سیستم خارج می‌شود و سپس همراه با پوشه‌اش، و همچنین تاریخچه گفتگو و تنظیماتش، پاک می‌شود.
- پوشه `~/.claude` که آن را [وارد کرده‌اید](/docs/commands) برای `claude` معمولی در جای خود باقی می‌ماند. اگر در یک گروه باشد، یک کپی از همه آنچه به اشتراک می‌گذاشت را پس می‌گیرد. خط وضعیت ccam در `settings.json` آن جای خود را به خط وضعیتی می‌دهد که ccam آن را در بر گرفته بود، اگر چنین خطی وجود داشت. اگر پلاگین ccam را داشته باشد، ccam فرمان `/plugin uninstall` را چاپ می‌کند تا آن را در Claude Code اجرا کنید.
- فرمان‌های میان‌بر، `~/.config/ccam`، `~/.local/share/ccam` (همراه با سطل زباله ccam) و `~/.local/state/ccam` پاک می‌شوند و سپس خود فرمان `ccam`.

تا وقتی حسابی در حال اجراست، ccam از این کار خودداری می‌کند؛ ابتدا نشست‌های Claude Code آن را ببندید. `/usr/bin/ccam` که apt نصب کرده باقی می‌ماند؛ بعداً آن را با `sudo apt remove ccam` حذف کنید. در Windows برنامه در حال اجرا نمی‌تواند خودش را پاک کند، پس `ccam uninstall` در پایان فرمان `Remove-Item` را برای `ccam.exe` چاپ می‌کند.
