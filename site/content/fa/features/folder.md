---
title: حساب‌ها بر اساس پوشه
description: ccam folder باعث می‌شود یک پوشه و هر چه زیر آن است، وقتی Claude Code را بدون نام حساب اجرا می‌کنید، از یک حساب استفاده کند، مثلاً work برای کد شرکت‌تان.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 46d1410a3d81
---

`ccam folder work ~/code/company` باعث می‌شود هر اجرای بدون نام حساب در `~/code/company` و در هر پوشه‌ای زیر آن، Claude Code را با `work` آغاز کند. بیرون از آن، حساب پیش‌فرض مثل قبل اجرا می‌شود. نام بردن از یک حساب، با `ccam personal` یا میان‌بری مانند `ccpersonal`، همیشه همان حساب را اجرا می‌کند.

## نحوه استفاده

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

این اجراها از حساب پوشه استفاده می‌کنند:

- `ccam -p "…"` و هر اجرای دیگری که آرگومان اولش گزینه‌ای برای Claude Code باشد
- `ccam run` بدون نام، و `ccam --dry-run`
- `ccam --fallback -p "…"` برای نخستین تلاشش
- داشبورد، یعنی `ccam` به‌تنهایی، که با مکان‌نما روی حساب پوشه آغاز می‌شود

## گزینه‌ها

| گزینه | کار آن | پیش‌فرض |
| --- | --- | --- |
| `<account>` | نام یا شماره حسابی که باید استفاده شود | |
| `[folder]` | پوشه، نسبی یا مطلق | پوشه فعلی |
| `--remove` | ورودی دقیقاً همان پوشه را حذف می‌کند | خاموش |
| `--dry-run` | پشتیبانی می‌شود: تغییر را بدون ذخیره کردن چاپ می‌کند | خاموش |

## سازوکار

هر پوشه یک ورودی `[[folder]]` در `~/.config/ccam/config.toml` است، همان فایلی که هر اجرا در هر حال می‌خواند:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` پوشه را به یک مسیر مطلق با پیوندهای نمادین حل‌شده تبدیل و ذخیره می‌کند. تنظیم دوباره یک پوشه، حساب آن را جایگزین می‌کند.
2. اجرای بدون نام، `launcher.HereProfile` را فرا می‌خواند. اگر هیچ ورودی‌ای نباشد، بلافاصله حساب پیش‌فرض را برمی‌دارد.
3. در غیر این صورت پوشه کاری را به‌صورت رشته با هر `path` مقایسه می‌کند: خود پوشه یا پوشه‌ای زیر آن، هرگز پوشه‌ای که فقط با همان حروف شروع شود (`/code/company` شامل `/code/companyx` نمی‌شود). عمیق‌ترین ورودی برنده است، پس `~/code` می‌تواند `personal` باشد در حالی که `~/code/company`، `work` است.
4. فقط اگر پوشه کاری همان‌طور که هست با هیچ ورودی‌ای مطابقت نکند، ccam پیوندهای نمادین آن را حل می‌کند و یک بار دیگر مقایسه می‌کند.
5. در Windows مقایسه، مانند خود Windows، به بزرگی و کوچکی حروف حساس نیست.

هیچ فایلی جز `config.toml` خوانده نمی‌شود، پس اجرا درون یک پوشه هیچ زمان قابل اندازه‌گیری‌ای نمی‌گیرد (4.0 ms با ورودی یا بدون آن، `hyperfine`، ۳۰۰ اجرا).

## ایمنی و موارد خاص

- `ccam rename` ورودی‌های پوشه را به نام جدید منتقل می‌کند و `ccam remove` ورودی‌های حساب حذف‌شده را پاک می‌کند.
- ورودی‌ای که حسابش دیگر وجود ندارد، اجرای بدون نام در آن پوشه را متوقف می‌کند: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` چنین ورودی‌هایی را گزارش می‌کند، و نیز ورودی‌هایی را که پوشه‌شان از بین رفته است.
- `ccam folder --remove` درون پوشه‌ای که حسابش را از پوشه‌ای بالاتر می‌گیرد، به جای حذف چیزی، نام آن پوشه بالاتر را می‌گوید.
- `ccam next` بر اساس چرخش یا مصرف انتخاب می‌کند و به پوشه‌ها نگاه نمی‌کند.

## مطالب مرتبط

- [اجرای Claude Code](/features/launch)
- [فرمان‌های میان‌بر](/features/shortcuts)
- [مدیریت حساب‌ها](/features/accounts)
- [مرجع فرمان‌ها](/docs/commands)
