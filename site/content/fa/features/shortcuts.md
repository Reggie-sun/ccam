---
title: فرمان‌های میان‌بر
description: با ccam shortcut به یک حساب فرمانی مخصوص خودش مانند ccwork بدهید و پیوندهای ~/.local/bin را با ccam shims install دوباره بسازید.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

میان‌بر فرمانی است که نامش را خودتان انتخاب می‌کنید و Claude Code را با یک حساب اجرا می‌کند: `ccwork --resume` همان کاری را می‌کند که `ccam work --resume` انجام می‌دهد. میان‌بر یک پیوند به ccam در `~/.local/bin` است، بنابراین در bash، zsh، اسکریپت‌ها و هر چیز دیگری که برنامه‌ها را از `PATH` اجرا می‌کند کار می‌کند. میان‌برها اختیاری‌اند؛ برای حساب‌هایی که زیاد اجرا می‌کنید یکی بسازید.

## کاربرد

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

ccam هر تغییر را تأیید می‌کند:

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

هر حساب حداکثر یک میان‌بر دارد که به‌صورت `shortcut` در `~/.config/ccam/config.toml` ذخیره می‌شود. دستیار `ccam add` ابتدا `Create a shortcut command for <name>? [y/N]` و سپس `Shortcut name [cc<name>]:` را می‌پرسد. بیرون از دستیار، `ccam add` فقط با `--shortcut` میان‌بر می‌سازد. `ccam list` میان‌بر هر حساب را نشان می‌دهد، یا `-` را.

میان‌بر همه آرگومان‌ها را به Claude Code می‌دهد، همان‌طور که `ccam <name>` آرگومان‌های پس از نام را می‌دهد. این شامل `--dry-run` و `--help` هم می‌شود: برای پیش‌نمایش یک اجرا از `ccam --dry-run work` استفاده کنید.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--remove` | میان‌بر حساب و پیوند آن را حذف می‌کند | خاموش |
| `-h`، `--help` | راهنمای فرمان | |
| `--dry-run` | `ccam shortcut` از آن پشتیبانی می‌کند: نام را بررسی می‌کند و بدون ذخیره یا پیوند دادن، `Would set the shortcut of work to ccwork.` یا `Would remove ccwork, the shortcut of work.` را چاپ می‌کند | خاموش |

`ccam shims install` هیچ آرگومانی نمی‌گیرد و `--dry-run` را با پیام `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` رد می‌کند.

قواعد نام‌گذاری، به همین ترتیب بررسی می‌شوند:

- ۱ تا ۶۴ نویسه از `a-z`، `0-9`، `.`، `-` و `_`، که با یک حرف یا رقم شروع شود؛
- نه `ccam`، `claude` یا یکی از فرمان‌های ccam مانند `run` یا `usage`؛
- نه میان‌بر حساب دیگری؛
- نه فایلی در `~/.local/bin` که از قبل پیوندی به ccam نباشد؛
- نه نام برنامه اجرایی دیگری در `PATH` شما.

دادن یک فرمان همراه با `--remove` با پیام `give a command or --remove, not both` متوقف می‌شود. تنظیم نامی که حساب از قبل دارد، `work: ccwork (unchanged)` را چاپ می‌کند و اگر پیوند آن وجود نداشته باشد، دوباره آن را می‌سازد.

## چگونه کار می‌کند

تنظیم یک نام، `checkShortcut` را اجرا می‌کند که `shim.Check` را با پوشه bin و مسیر واقعی فایل اجرایی ccam فرا می‌خواند:

<!-- motion:shortcuts-check -->

نام به ترتیب از بررسی‌ها عبور می‌کند. کنار هر بررسی نامی آمده که آن بررسی رد می‌کند: `Work` با الگو جور نیست، `run` یکی از فرمان‌های ccam است، `ccside` به حساب دیگری تعلق دارد، `ccfile` فایلی موجود در `~/.local/bin` است و `git` فایل `/usr/bin/git` را پنهان می‌کند. فقط نامی که از همه بررسی‌ها عبور کند ذخیره و پیوند داده می‌شود.

1. `shim.Check` نام را با `^[a-z0-9][a-z0-9._-]{0,63}$`، سپس با `config.Reserved` و بعد با `cfg.FindShortcut` مقایسه می‌کند. `os.Lstat` را روی `~/.local/bin/<name>` فرا می‌خواند و برای هر پوشه دیگر در `PATH`، `os.Stat` را فرا می‌خواند و فایل اجرایی‌ای را که پیوندی به ccam نباشد رد می‌کند.
2. ccam فایل `config.toml` را با `fsutil.AtomicWrite` ذخیره می‌کند (فایل موقت، `fsync`، تغییر نام).
3. `shim.Remove` پیوند میان‌بر قدیمی را حذف می‌کند، اما فقط اگر پیوندی به ccam باشد.
4. `shim.Install` پیوند هر میان‌بر پیکربندی‌شده‌ای را که وجود ندارد با `os.Symlink` به فایل اجرایی ccam می‌سازد، که با `os.Executable` و `filepath.EvalSymlinks` پیدا می‌شود. هرگز چیزی را حذف نمی‌کند.

یک پیوند وقتی متعلق به خود ccam شمرده می‌شود که `os.Readlink` مسیر ccam را برگرداند، یا پیوند و مسیر ccam به یک فایل واحد برسند. هر حذف و هر بررسی بازنویسی از این آزمون می‌گذرد، بنابراین ccam به فایل‌هایی که خودش نساخته دست نمی‌زند.

وقتی یک میان‌بر را اجرا می‌کنید، نامی که تایپ کرده‌اید تنها چیزی است که ccam لازم دارد:

<!-- motion:shortcuts-link -->

پوسته `ccwork` را در `PATH` پیدا می‌کند، هسته پیوند را دنبال می‌کند و فایل اجرایی ccam را آغاز می‌کند، و `argv[0]` همچنان `ccwork` است. در `app.Main`، هر نامی به جز `ccam` فایل `config.toml` را یک بار بارگذاری می‌کند و `shim.Lookup` را فرا می‌خواند، که حسابی را پیدا می‌کند که `shortcut` آن برابر با `filepath.Base(argv[0])` است. سپس `launcher.LaunchWith` با پیکربندی‌ای که از قبل بارگذاری شده همان گام‌های [`ccam <name>`](/features/launch) را برمی‌دارد و `argv[1:]` را بدون تغییر به `claude` می‌دهد.

چرا پیوند نمادین و نه alias پوسته یا اسکریپت پوشاننده: پیوند نمادین در هر پوسته و در اسکریپت‌ها کار می‌کند، هیچ اجرای اضافه پوسته‌ای (حدود ۱ تا ۲ میلی‌ثانیه) اضافه نمی‌کند و اجرا همچنان فقط یک فایل را می‌خواند. `scripts/perf-check.sh` یک میان‌بر را با `hyperfine` در برابر همان حد ۱۰ میلی‌ثانیه (p95) که برای `ccam <name>` است اندازه می‌گیرد.

## ایمنی و موارد خاص

- **پیوندی که هیچ حسابی ادعایش را ندارد.** اگر پیوندی به ccam در `~/.local/bin` نامی داشته باشد که هیچ حسابی از آن استفاده نمی‌کند، برای مثال پس از آنکه `shortcut` را در `config.toml` دستی ویرایش کرده‌اید، ccam به حساب پیش‌فرض برنمی‌گردد. با پیام `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` متوقف می‌شود. نام را دوباره با `ccam shortcut` تنظیم کنید، یا پیوند را حذف کنید.
- **فایل‌های موجود.** ccam هرگز فایلی را که خودش نساخته جایگزین نمی‌کند. `ccam shims install` با پیام `<path> already exists and is not a ccam link, so it was left alone` متوقف می‌شود.
- **فایل اجرایی ccam جابه‌جاشده.** در این حالت پیوندهای قدیمی به مسیری اشاره می‌کنند که دیگر وجود ندارد و ccam آن‌ها را متعلق به خود نمی‌شمارد. آن‌ها را حذف کنید، سپس `ccam shims install` را اجرا کنید.
- **تغییر نام و حذف.** `ccam rename` نام میان‌بر را نگه می‌دارد و پیوند همچنان کار می‌کند، چون به ccam اشاره می‌کند و ccam نام را در `config.toml` جست‌وجو می‌کند. `ccam remove` پیوند حساب را حذف می‌کند. `ccam shortcut <account> <new>` پیش از ساختن پیوند جدید، پیوند قدیمی را حذف می‌کند.
- **PATH.** وقتی ccam یک پیوند می‌سازد و `~/.local/bin` در `PATH` شما نیست، `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` را چاپ می‌کند و پس از آن خط `export PATH=...` را که باید اضافه کنید.
- **بررسی در زمان بعد.** `ccam doctor` پیام `missing shortcuts in /home/me/.local/bin: ccwork` را گزارش می‌کند و وقتی پوشه در `PATH` نباشد هشدار می‌دهد. `ccam doctor --fix` پیوندهای ناموجود را می‌سازد.

## مرتبط

- [اجرای Claude Code](/features/launch)
- [مدیریت حساب‌ها](/features/accounts)
- [مرجع فرمان‌ها](/docs/commands#shortcuts)
- [نصب ccam](/docs/install)
