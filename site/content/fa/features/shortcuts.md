---
title: فرمان‌های میان‌بر
description: با ccam shortcut به یک حساب فرمانی مخصوص خودش مانند ccwork بدهید و پیوندهای ~/.local/bin را با ccam shims install دوباره بسازید.
updated: 2026-10-23
order: 3
translation: machine
sourceHash: 38e2d1aba397
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
ccam shims install            # recreate missing or outdated links of all accounts
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
- نه فایلی در `~/.local/bin` که از قبل پیوند ccam نباشد (یا میان‌بر خود همین حساب که ccam ساخته است)؛
- نه نام برنامه اجرایی دیگری در `PATH` شما.

`config.toml` هنگام بارگذاری با همین قواعد بررسی می‌شود: میان‌بری که دستی ویرایش شده و با الگو جور نیست، `claude` یا `ccam` است، یا به دو حساب داده شده است، باعث می‌شود هر فرمانی با پیامی که نام حساب را می‌برد ناموفق شود، تا وقتی رفع شود.

دادن یک فرمان همراه با `--remove` با پیام `give a command or --remove, not both` متوقف می‌شود. تنظیم نامی که حساب از قبل دارد، `work: ccwork (unchanged)` را چاپ می‌کند و اگر پیوند آن وجود نداشته باشد، دوباره آن را می‌سازد.

## چگونه کار می‌کند

تنظیم یک نام، `checkShortcut` را اجرا می‌کند که `shim.Check` را با پوشه bin و مسیر واقعی فایل اجرایی ccam فرا می‌خواند:

<!-- motion:shortcuts-check -->

نام به ترتیب از بررسی‌ها عبور می‌کند. کنار هر بررسی نامی آمده که آن بررسی رد می‌کند: `Work` با الگو جور نیست، `run` یکی از فرمان‌های ccam است، `ccside` به حساب دیگری تعلق دارد، `ccfile` فایلی موجود در `~/.local/bin` است و `git` فایل `/usr/bin/git` را پنهان می‌کند. فقط نامی که از همه بررسی‌ها عبور کند ذخیره و پیوند داده می‌شود.

1. `shim.Check` نام را با `^[a-z0-9][a-z0-9._-]{0,63}$`، سپس با `config.Reserved` و بعد با `cfg.FindShortcut` مقایسه می‌کند. `os.Lstat` را روی `~/.local/bin/<name>` فرا می‌خواند و برای هر پوشه دیگر در `PATH`، `os.Stat` را فرا می‌خواند و فایل اجرایی‌ای را که پیوندی به ccam نباشد رد می‌کند.
2. ccam فایل `config.toml` را با `fsutil.AtomicWrite` ذخیره می‌کند (فایل موقت، `fsync`، تغییر نام).
3. `shim.Remove` پیوند میان‌بر قدیمی را حذف می‌کند، اما فقط اگر پیوندی به ccam باشد.
4. `shim.Install` همه میان‌برهای پیکربندی‌شده را به ترتیب نام مرور می‌کند و پیوند هر کدام را که وجود ندارد با `os.Symlink` به فایل اجرایی ccam می‌سازد، که با `os.Executable` و `filepath.EvalSymlinks` پیدا می‌شود. فایلی که سر راه یک میان‌بر باشد مانع بقیه نمی‌شود؛ خطا نام هر کدام را می‌برد.

ccam هر پیوندی را که می‌سازد، همراه با فایل اجرایی‌ای که پیوند برای اجرای آن ساخته شده، در `~/.local/state/ccam/shortcuts.json` ثبت می‌کند. یک پیوند وقتی متعلق به خود ccam شمرده می‌شود که `os.Readlink` مسیر ccam را برگرداند، یا پیوند و مسیر ccam به یک فایل واحد برسند. پیوند ثبت‌شده‌ای که فایل اجرایی‌اش از آن پس جابه‌جا یا جایگزین شده است (ccamی که به پوشه دیگری منتقل شده، یک نصب دوباره) کهنه است و `shim.Install` آن را به این ccam هدایت می‌کند. پیوند ثبت‌شده‌ای به ccam دیگری که هنوز نصب است همان‌طور که هست می‌ماند. هر حذف و هر بررسی بازنویسی از این آزمون‌ها می‌گذرد، بنابراین ccam به فایل‌هایی که خودش نساخته دست نمی‌زند، مانند کپی خود شما از ccam یا پیوندی که خودتان ساخته‌اید.

وقتی یک میان‌بر را اجرا می‌کنید، نامی که تایپ کرده‌اید تنها چیزی است که ccam لازم دارد:

<!-- motion:shortcuts-link -->

پوسته `ccwork` را در `PATH` پیدا می‌کند، هسته پیوند را دنبال می‌کند و فایل اجرایی ccam را آغاز می‌کند، و `argv[0]` همچنان `ccwork` است. در `app.Main`، هر نامی به جز `ccam` فایل `config.toml` را یک بار بارگذاری می‌کند و `shim.Lookup` را فرا می‌خواند، که حسابی را پیدا می‌کند که `shortcut` آن برابر با `filepath.Base(argv[0])` است. سپس `launcher.LaunchWith` با پیکربندی‌ای که از قبل بارگذاری شده همان گام‌های [`ccam <name>`](/features/launch) را برمی‌دارد و `argv[1:]` را بدون تغییر به `claude` می‌دهد.

چرا پیوند نمادین و نه alias پوسته یا اسکریپت پوشاننده: پیوند نمادین در هر پوسته و در اسکریپت‌ها کار می‌کند، هیچ اجرای اضافه پوسته‌ای (حدود ۱ تا ۲ میلی‌ثانیه) اضافه نمی‌کند و اجرا همچنان فقط یک فایل را می‌خواند. `scripts/perf-check.sh` یک میان‌بر را با `hyperfine` در برابر همان حد ۱۰ میلی‌ثانیه (p95) که برای `ccam <name>` است اندازه می‌گیرد.

### در Windows

یک میان‌بر فایل `ccwork.exe` در `%USERPROFILE%\.local\bin` است و `shim.Name` پیش از جست‌وجو پسوند `.exe` را (با هر حالت حروف) حذف می‌کند، پس `ccwork`، `ccwork.exe` و `CCWORK.EXE` همه `work` را اجرا می‌کنند. ccam ابتدا یک پیوند نمادین می‌سازد که Windows آن را در Developer Mode یا برای مدیر سیستم مجاز می‌داند. بدون این حق، یک hard link می‌سازد: نام دومی برای همان فایل `ccam.exe` که به هیچ دسترسی ویژه‌ای نیاز ندارد. یک hard link وقتی متعلق به ccam شمرده می‌شود که همان فایل `ccam.exe` باشد. یک hard link همچنان فایلی را اجرا می‌کند که از آن ساخته شده است، پس وقتی `ccam update` فایل جدیدی را جای `ccam.exe` می‌گذارد، هر میان‌بری را که hard link است با `shim.Relink` تازه می‌کند و `✓ shortcut ccwork renewed` را چاپ می‌کند. اجرای دوباره `install.ps1` همین کار را از طریق `ccam shims install` انجام می‌دهد (`OK 2 shortcuts now run the new ccam`)، و `ccam shims install` و `ccam doctor --fix` هم برای هر میان‌بری که کپی یک ccam قدیمی‌تر است همین کار را می‌کنند. بررسی نام، `<name>` را با هر پسوند موجود در `PATHEXT` (`.exe`، `.cmd` و غیره) در دیگر پوشه‌های `PATH` جست‌وجو می‌کند.

میان‌بر را در یک پوسته تایپ کنید (PowerShell، cmd.exe، Git Bash). پیوند نمادینی که از Explorer، پنجره Run یا `Start-Process` باز شود به‌صورت `ccam` ساده اجرا می‌شود که به جای حساب، داشبورد را باز می‌کند.

## ایمنی و موارد خاص

- **پیوندی که هیچ حسابی ادعایش را ندارد.** اگر پیوندی به ccam در `~/.local/bin` نامی داشته باشد که هیچ حسابی از آن استفاده نمی‌کند، برای مثال پس از آنکه `shortcut` را در `config.toml` دستی ویرایش کرده‌اید، ccam به حساب پیش‌فرض برنمی‌گردد. با پیام `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` متوقف می‌شود. نام را دوباره با `ccam shortcut` تنظیم کنید، یا پیوند را حذف کنید.
- **فایل‌های موجود.** ccam هرگز فایلی را که خودش نساخته جایگزین نمی‌کند. `ccam shims install` پیام `<path> already exists and is not a ccam link, so it was left alone` را گزارش می‌کند و با این حال پیوندهای دیگر را می‌سازد. در این حالت `ccam add` و `ccam import` پیام `shortcut: ccwork (not created; see above)` را چاپ می‌کنند.
- **فایل اجرایی ccam جابه‌جاشده.** در این حالت پیوندهایی که ccam ساخته به مسیری اشاره می‌کنند که دیگر وجود ندارد. `ccam shims install` (یا `ccam doctor --fix`) آن‌ها را به ccamی هدایت می‌کند که با آن اجرایش می‌کنید.
- **تغییر نام و حذف.** `ccam rename` نام میان‌بر را نگه می‌دارد و پیوند همچنان کار می‌کند، چون به ccam اشاره می‌کند و ccam نام را در `config.toml` جست‌وجو می‌کند. `ccam remove` پیوند حساب را حذف می‌کند. `ccam shortcut <account> <new>` پیش از ساختن پیوند جدید، پیوند قدیمی را حذف می‌کند.
- **PATH.** وقتی ccam یک پیوند می‌سازد و `~/.local/bin` در `PATH` شما نیست، `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` را چاپ می‌کند و پس از آن خط `export PATH=...` را که باید اضافه کنید. در Windows به جای آن خط PowerShell را چاپ می‌کند که پوشه را به `PATH` کاربری شما اضافه می‌کند.
- **بررسی در زمان بعد.** `ccam doctor` پیام `missing shortcuts in /home/me/.local/bin: ccwork` و میان‌برهایی را که یک ccam قدیمی‌تر را اجرا می‌کنند گزارش می‌کند و وقتی پوشه در `PATH` نباشد هشدار می‌دهد. `ccam doctor --fix` پیوندهای ناموجود را می‌سازد و پیوندهای کهنه را تازه می‌کند، و فقط برای آن‌هایی که ساخته است `→ fixed` می‌گوید. میان‌بری که ccam دیگری را اجرا می‌کند که هنوز نصب است کار می‌کند و فقط اطلاعات است؛ `ccam doctor --fix` آن را به ccamی هدایت می‌کند که با آن اجرایش می‌کنید، و `ccam add` و `ccam import` آن را به‌صورت ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)`` نشان می‌دهند.

## مرتبط

- [اجرای Claude Code](/features/launch)
- [مدیریت حساب‌ها](/features/accounts)
- [مرجع فرمان‌ها](/docs/commands#shortcuts)
- [نصب ccam](/docs/install)
