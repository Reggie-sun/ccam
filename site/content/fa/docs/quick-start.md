---
title: شروع سریع
metaTitle: "شروع سریع: افزودن و تعویض حساب‌های Claude Code"
description: دو حساب Claude Code را با ccam add اضافه کنید، به آن‌ها میان‌برهایی مانند ccwork بدهید، بین آن‌ها جابه‌جا شوید و تعیین کنید چه چیزهایی مشترک باشد.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

هر حساب را با `ccam add` اضافه کنید، سپس Claude Code را با `ccam <name>` یا با میان‌بری که نامش را خودتان انتخاب می‌کنید، مانند `ccwork`، اجرا کنید. بیشتر کاربران به چیزی بیش از این نیاز ندارند.

## ۱. نخستین حساب را اضافه کنید

```sh
ccam add
```

دستیار راه‌اندازی یک نام مانند `personal` یا `work` از شما می‌پرسد، اینکه حساب کدام موارد را با حساب‌های دیگر شما به اشتراک بگذارد و اینکه یک فرمان میان‌بر ساخته شود یا نه (مرحله ۴ را ببینید). سپس `claude auth login` را اجرا می‌کند که مرورگر را باز می‌کند تا وارد حساب شوید. روی یک ماشین راه دور، ممکن است مرورگر در صفحه `http://localhost:…/callback` بماند که بارگذاری نمی‌شود؛ آن نشانی را در ترمینال بچسبانید تا ccam ورود را کامل کند. در پایان این را می‌بینید:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

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

اجرای `ccam` به‌تنهایی [داشبورد](/features/dashboard) را باز می‌کند: با کلیدهای جهت‌نما یا شماره حساب به آن بروید، سپس Enter را بزنید.

هر چه پس از میان‌بر یا نام حساب بیاید بدون تغییر به Claude Code می‌رسد، بنابراین می‌توانید همان گزینه‌هایی را که می‌شناسید به کار ببرید:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## ۵. تعیین کنید چه چیزهایی مشترک باشد

به‌طور پیش‌فرض، سرورهای MCP، `settings.json`، فایل سراسری `CLAUDE.md` و همچنین agents، commands، skills و plugins از طریق گروه `default` بین حساب‌ها مشترک‌اند. تاریخچه گفتگو مشترک نیست. اگر می‌خواهید گفتگوها را بین حساب‌ها ادامه دهید، آن را به گروه اضافه کنید:

```sh
ccam group set default --add history
```

ccam پیش از هر تغییری نشان می‌دهد برای هر حساب چه چیزهایی را جابه‌جا خواهد کرد. برای جزئیات، [اشتراک‌گذاری با گروه‌ها](/features/share) را ببینید.

## ۶. وقتی یک حساب به سقف خود می‌رسد

```sh
ccam usage
ccam next --continue
```

`ccam usage` مصرف ۵ ساعته و هفتگی همه حساب‌ها را نشان می‌دهد. `ccam next --continue` همان گفتگو را در حساب بعدی باز می‌کند. در اسکریپت‌ها، `ccam --fallback work -p "…"` خودش حساب را عوض می‌کند. [سقف مصرف](/docs/usage-limits) را ببینید.

**صندلی‌های Team و Enterprise.** انتخاب خودکار (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد، و [`auto_switch`](/features/next#auto-switch) وقتی آن را روشن کنید) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود. مرکز راهنمای Anthropic می‌گوید: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (هر صندلی در طرح Team به نام یک فرد مشخص اختصاص داده می‌شود و فقط همان شخص اجازه دارد از آن استفاده کند) و «Users can't route their work through another user's account to use that user's remaining usage» (کاربران نمی‌توانند کار خود را از طریق حساب کاربر دیگری انجام دهند تا از سهمیه مصرف باقی‌مانده آن کاربر استفاده کنند) ([طرح‌های Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ و همین‌طور برای [طرح‌های Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). بنابراین استفاده از انتخاب خودکار میان صندلی‌ها ممکن است شرایط Anthropic را نقض کند. اینکه در آنجا از آن استفاده کنید یا نه، تصمیم خودتان است: پیش از آنکه به آن تکیه کنید، از سازمان خود و Anthropic بپرسید. برای طرح‌های شخصی (Pro، Max)، [توضیح شرایط در پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را ببینید.

## از ~/.claude استفاده می‌کنید؟

به‌جای ورود دوباره، آن را ثبت کنید:

```sh
ccam import
```

پوشه در همان جایی که هست می‌ماند و `claude` معمولی مثل قبل کار می‌کند.
