---
title: شروع سریع
description: دو حساب Claude Code را با ccam add اضافه کنید، به آن‌ها میان‌برهایی مانند ccwork بدهید، بین آن‌ها جابه‌جا شوید و تعیین کنید چه چیزهایی مشترک باشد.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: b73914ff371f
---

هر حساب را با `ccam add` اضافه کنید، سپس Claude Code را با `ccam <name>` یا با میان‌بری که نامش را خودتان انتخاب می‌کنید، مانند `ccwork`، اجرا کنید. بیشتر کاربران به چیزی بیش از این نیاز ندارند.

## ۱. نخستین حساب را اضافه کنید

```sh
ccam add
```

دستیار راه‌اندازی یک نام مانند `personal` یا `work` از شما می‌پرسد، اینکه حساب کدام موارد را با حساب‌های دیگر شما به اشتراک بگذارد و اینکه یک فرمان میان‌بر ساخته شود یا نه (مرحله ۴ را ببینید). سپس `claude auth login` را اجرا می‌کند که مرورگر را باز می‌کند تا وارد حساب شوید. در پایان این را می‌بینید:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

ccam پیام‌های خود را به زبان انگلیسی چاپ می‌کند؛ فرمان‌ها و گزینه‌ها در همه زبان‌ها یکسان‌اند.

## ۲. تنظیمات فعلی خود را منتقل کنید

اگر پیش‌تر از Claude Code استفاده کرده‌اید، تنظیمات، `CLAUDE.md`، افزونه‌ها و سرورهای MCP شما در `~/.claude` هستند. آن‌ها را در تنظیمات مشترک کپی کنید تا همه حساب‌ها از آن‌ها بهره ببرند:

```sh
ccam migrate
```

ccam نشان می‌دهد چه چیزی را کپی خواهد کرد و پیش از آن از شما می‌پرسد. فقط `~/.claude` را می‌خواند، هرگز اطلاعات ورود را کپی نمی‌کند و تاریخچه گفتگو را تنها وقتی کپی می‌کند که `--history` را اضافه کنید. برای جزئیات، [کپی کردن یک پیکربندی موجود](/docs/sharing#copy-an-existing-setup) را ببینید.

## ۳. حساب دوم را اضافه کنید

```sh
ccam add work
```

در مرورگر با حساب دیگر وارد شوید.

## ۴. بین حساب‌ها جابه‌جا شوید

`ccam work` برنامه Claude Code را با حساب کاری اجرا می‌کند. برای فرمانی کوتاه‌تر، به حساب یک فرمان میان‌بر مخصوص خودش بدهید، مانند `ccwork`. دستیار راه‌اندازی پس از پرسیدن موارد مشترک این را می‌پرسد:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` را بزنید، سپس یک نام بنویسید یا Enter را بزنید تا نام پیشنهادی استفاده شود:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

پس از ورود، ccam آن را تأیید می‌کند:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

اگر از این پرسش گذشتید، هر زمان که بخواهید می‌توانید میان‌بری اضافه کنید:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

اکنون هر یک از این فرمان‌ها Claude Code را با حساب کاری اجرا می‌کند:

```sh
ccwork
ccam work
ccam 2
```

اجرای `ccam` به‌تنهایی یک انتخابگر باز می‌کند: از کلیدهای جهت‌نما و Enter استفاده کنید یا شماره حساب را بزنید.

هر چه پس از میان‌بر یا نام حساب بیاید بدون تغییر به Claude Code می‌رسد، بنابراین می‌توانید همان گزینه‌هایی را که می‌شناسید به کار ببرید:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## ۵. تعیین کنید چه چیزهایی مشترک باشد

به‌طور پیش‌فرض، سرورهای MCP، `settings.json`، فایل سراسری `CLAUDE.md` و همچنین agents، commands، skills و plugins بین حساب‌ها مشترک‌اند. تاریخچه گفتگو مشترک نیست. اگر می‌خواهید گفتگوها را بین حساب‌ها ادامه دهید، آن را برای هر دو حساب روشن کنید:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam پیش از هر تغییری نشان می‌دهد چه چیزهایی را جابه‌جا خواهد کرد. برای جزئیات، [اشتراک‌گذاری](/docs/sharing) را ببینید.

## ۶. وقتی یک حساب به سقف خود می‌رسد

```sh
ccam usage
ccam next --continue
```

`ccam usage` مصرف ۵ ساعته و هفتگی همه حساب‌ها را نشان می‌دهد. `ccam next --continue` همان گفتگو را در حساب بعدی باز می‌کند. [سقف مصرف](/docs/usage-limits) را ببینید.

## از ~/.claude استفاده می‌کنید؟

به‌جای ورود دوباره، آن را ثبت کنید:

```sh
ccam import
```

پوشه در همان جایی که هست می‌ماند و `claude` معمولی مثل قبل کار می‌کند.
