---
title: به‌روزرسانی ccam
description: ccam update پس از بررسی SHA-256، ccam را با آخرین نسخه GitHub جایگزین می‌کند؛ --check فقط گزارش می‌دهد و --version یک برچسب نسخه مشخص را نصب می‌کند.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: 598931e4801f
---

`ccam update` یک نسخه ccam را از GitHub دانلود می‌کند، آن را با `checksums.txt` همان نسخه بررسی می‌کند و فایل اجرایی در حال اجرا را با یک تغییر نام جایگزین می‌کند. اگر ccam را با اسکریپت نصب یا از یک فایل فشرده نصب کرده‌اید، از این فرمان استفاده کنید. [داشبورد](/features/dashboard) هم هر بار که باز می‌شود آخرین نسخه را از GitHub می‌پرسد و آن را با `U` پیشنهاد می‌کند، که پس از پرسیدن همین مراحل را اجرا می‌کند. فرمان‌های دیگر هرگز با GitHub ارتباط نمی‌گیرند، بنابراین آن‌ها به‌صورت آفلاین هم همان‌طور کار می‌کنند.

## استفاده

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
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

وقتی نسخه جدیدتری وجود نداشته باشد، ccam پیام `ccam 0.4.0 is up to date.` را چاپ می‌کند و چیزی را تغییر نمی‌دهد. اگر `--version` همان نسخه‌ای را نام ببرد که اکنون دارید، `ccam 0.4.0 is already installed.` را چاپ می‌کند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--check` | فقط گزارش می‌دهد که نسخه جدیدتری وجود دارد یا نه | خاموش |
| `--version <tag>` | این برچسب نسخه را نصب می‌کند، مانند `v0.3.0`، چه جدیدتر و چه قدیمی‌تر از نسخه فعلی | آخرین نسخه |
| `--dry-run` | پشتیبانی می‌شود: آخرین برچسب را پیدا می‌کند (مگر اینکه `--version` برچسبی را نام ببرد)، همان بررسی‌ها را روی فایل اجرایی انجام می‌دهد و مسیرهای دانلود، checksum و جایگزینی را بدون دانلود چاپ می‌کند | خاموش |

`--check` و `--version` را نمی‌توان با هم استفاده کرد. برچسب باید به شکل `vMAJOR.MINOR.PATCH` باشد و می‌تواند یک پسوند پیش‌انتشار یا build هم داشته باشد؛ می‌توان `v` ابتدای آن را حذف کرد. هر چیز دیگری خطا می‌دهد، برای مثال `--version takes a release tag such as v0.3.0, not "1.2"`.

متغیر محیطی `CCAM_RELEASES_URL` برای یک mirror یا برای آزمون‌ها جایگزین `https://github.com/jhyunwoo/ccam/releases` می‌شود. ساختار آن باید یکسان باشد: `<url>/latest` به `<url>/tag/<tag>` هدایت می‌شود و فایل‌ها در `<url>/download/<tag>/<asset>` قرار دارند.

## چگونه کار می‌کند

ccam پیش از آنکه به فایل اجرایی نصب‌شده دست بزند، دانلود را بررسی می‌کند و فایل اجرایی را با یک تغییر نام جایگزین می‌کند. کد در `internal/update` است که فقط `ccam update` آن را import می‌کند. مسیر اجرا هرگز آن را فراخوانی نمی‌کند.

1. ccam فایل اجرایی در حال اجرا را با `os.Executable` پیدا می‌کند و پیوندهای نمادین را با `filepath.EvalSymlinks` دنبال می‌کند (`cli.updateCmd`). همین مسیر واقعی است که جایگزین می‌شود.
2. بدون `--version`، ccam نشانی `<releases>/latest` را بدون دنبال کردن هدایت درخواست می‌کند (`Updater.Latest`) و برچسب را از سرآیند `Location` (`…/tag/v0.5.0`) می‌خواند. این روش به‌جای API GitHub از صفحه نسخه‌ها استفاده می‌کند، بنابراین محدودیت ۶۰ درخواست در ساعت API شامل آن نمی‌شود. مهلت این جست‌وجو پس از ۱۰ ثانیه تمام می‌شود.
3. ccam برچسب را با یک الگوی سخت‌گیرانه بررسی می‌کند (`parseVersion`). برچسب در نشانی‌های دانلود و در یک فرمان shell چاپ‌شده قرار می‌گیرد، پس هیچ چیزی که `/`، `%`، فاصله یا نویسه‌های shell داشته باشد پذیرفته نمی‌شود. سپس ccam آن را با نسخه در حال اجرا مقایسه می‌کند (`compareVersions`)؛ یک پیش‌انتشار پیش از نسخه نهایی خود قرار می‌گیرد.
4. ccam تعیین می‌کند که آیا اصلاً اجازه دارد فایل اجرایی را جایگزین کند یا نه (`managed()`، `writable()`؛ پایین‌تر را ببینید).
5. ccam فایل‌های `ccam_<os>_<arch>.tar.gz` (در Windows `ccam_windows_<arch>.zip`) و `checksums.txt` را هم‌زمان دانلود می‌کند و برای هر دو روی هم ۲ دقیقه فرصت دارد (`Updater.Fetch`). حجم هر دانلود حداکثر ۲۵۶ MB است.
6. پیش از هر کار دیگری، SHA-256 فایل فشرده باید با خط مربوط به آن در `checksums.txt` یکی باشد. نبودن خط یا متفاوت بودن checksum به‌روزرسانی را با `nothing was changed` متوقف می‌کند.
7. ccam فایل tar را می‌خواند و نخستین فایل عادی با نام `ccam` را برمی‌دارد (`extract()`؛ `extractZip` فایل `ccam.exe` را از zip برمی‌دارد). پوشه‌ها، پیوندها و دیگر ورودی‌ها نادیده گرفته می‌شوند و حجم فایل اجرایی هم حداکثر ۲۵۶ MB است.
8. ccam فایل اجرایی جدید را با `fsutil.AtomicWrite` می‌نویسد (`Updater.Install`): یک فایل موقت `.ccam.tmp*` کنار فایل اجرایی، با بیت‌های مجوز فایل قبلی، که روی دیسک ذخیره و سپس با تغییر نام جایگزین فایل قبلی می‌شود.

کل این فرایند، از یافتن برچسب تا تغییر نام:

<!-- motion:update-pipeline -->

فایل فشرده و `checksums.txt` به‌صورت موازی دانلود می‌شوند و تا وقتی checksum مطابقت نداشته باشد چیزی نوشته نمی‌شود. تغییر نام اتمی است، بنابراین یک خرابی یا پر شدن دیسک یا فایل اجرایی قبلی را باقی می‌گذارد یا فایل جدید را، هرگز نیمی از یکی را. فرمان‌های میان‌بر مانند `~/.local/bin/ccwork` پیوندهای نمادین به فایل اجرایی ccam هستند، بنابراین بدون هیچ تغییری نسخه جدید را اجرا می‌کنند.

در Windows فایل برنامه‌ای که در حال اجراست قابل جایگزینی یا حذف نیست، فقط می‌توان نامش را تغییر داد. ccam فایل اجرایی جدید را در `ccam.exe.new` می‌نویسد، نام `ccam.exe` در حال اجرا را به `ccam.exe.old` تغییر می‌دهد و نام فایل جدید را به `ccam.exe` تغییر می‌دهد (`Updater.Install`). اگر گام آخر ناموفق شود، نام فایل اجرایی قدیمی به حالت قبل برگردانده می‌شود. به‌روزرسانی بعدی فایل‌های `ccam.exe.old` را که بتواند پاک می‌کند؛ وقتی یک نشست هنوز از یک `ccam.exe.old` قدیمی‌تر اجرا می‌شود، فایل اجرایی در حال اجرا به جای آن `ccam.exe.old-<time>` می‌شود. میان‌برهایی که به جای پیوند نمادین hard link هستند ([فرمان‌های میان‌بر](/features/shortcuts)) بلافاصله پس از آن تازه می‌شوند، پس آن‌ها هم نسخه جدید را اجرا می‌کنند؛ این شامل hard linkهای هر نسخه قدیمی‌تری از ccam که ccam ساخته است می‌شود، نه فقط فایل اجرایی‌ای که همین حالا جایگزین شد.

### نصب‌های مدیریت‌شده با بسته

نسخه‌ای از ccam که از `.deb` نصب شده متعلق به apt است و جایگزین کردن آن بدون اطلاع apt، پایگاه داده بسته‌ها را با فایل ناهماهنگ می‌کند. ccam پیش از دانلود بررسی می‌کند که فایل اجرایی متعلق به کیست:

<!-- motion:update-managed -->

1. فایل اجرایی در `/usr/bin`، `/bin`، `/usr/sbin` یا `/sbin` مدیریت‌شده با بسته به حساب می‌آید.
2. در جاهای دیگر، اگر `dpkg` نصب باشد و `dpkg -S <path>` موفق شود، آن هم مدیریت‌شده با بسته به حساب می‌آید.
3. در غیر این صورت ccam یک فایل موقت `.ccam-update-*` کنار فایل اجرایی می‌سازد و حذف می‌کند تا مطمئن شود می‌تواند در آنجا بنویسد.

برای فایل اجرایی مدیریت‌شده با بسته، ccam متوقف می‌شود و گام‌های apt را چاپ می‌کند:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

در سیستم‌هایی غیر از Linux، فقط خط اسکریپت نصب چاپ می‌شود: در Windows `irm https://ccam.dev/install.ps1 | iex`. وقتی پوشه قابل نوشتن نباشد، ccam پیشنهاد می‌کند اسکریپت نصب را دوباره اجرا کنید که ccam را در `~/.local/bin` نصب می‌کند.

## ایمنی و حالت‌های خاص

- یک build توسعه نسخه خود را `dev` گزارش می‌کند و نسخه‌ای برای مقایسه ندارد، بنابراین `ccam update` و `ccam update --check` آن را رد می‌کنند. `ccam update --version vX.Y.Z` یک نسخه را روی آن نصب می‌کند. برای buildهایی که با `go install` ساخته شده‌اند، [نصب ccam](/docs/install#other-ways-to-install) را ببینید.
- دانلود ناموفق، نبودن فایل مخصوص پلتفرم شما یا عدم تطابق checksum، فایل اجرایی نصب‌شده را دست‌نخورده باقی می‌گذارد.
- نسخه‌ای که وجود ندارد با `there is no release <tag>` خطا می‌دهد. وقتی GitHub با HTTP 403 یا 429 درخواست را رد کند، ccam می‌گوید احتمالاً محدودیت نرخ درخواست است و از شما می‌خواهد بعداً دوباره امتحان کنید.
- `--check` فقط می‌خواند. بررسی نمی‌کند که فایل اجرایی مدیریت‌شده با بسته یا قابل نوشتن است.
- استفاده از `--version` با یک برچسب قدیمی‌تر، راه بازگشت به نسخه قبلی است.

## مرتبط

- [نصب ccam](/docs/install)
- [مرجع فرمان‌ها](/docs/commands)
- [پرسش‌های متداول](/docs/faq)
