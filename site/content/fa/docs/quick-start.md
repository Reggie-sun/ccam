---
title: شروع سریع
metaTitle: "شروع سریع: افزودن و تعویض حساب‌های Claude Code"
description: دو حساب Claude Code را با ccam add اضافه کنید، به آن‌ها میان‌برهایی مانند ccwork بدهید، بین آن‌ها جابه‌جا شوید و تعیین کنید چه چیزهایی مشترک باشد.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

هر حساب را با `ccam add` اضافه کنید، سپس Claude Code را با `ccam <name>` یا با میان‌بری که نامش را خودتان انتخاب می‌کنید، مانند `ccwork`، اجرا کنید. بیشتر کاربران به چیزی بیش از این نیاز ندارند.

از قبل از Claude Code استفاده می‌کنید؟ پیش از شروع، بخش [از ~/.claude استفاده می‌کنید؟](#already-using-claude) را در انتهای همین صفحه بخوانید.

## ۱. نخستین حساب را اضافه کنید

```sh
ccam add
```

دستیار راه‌اندازی یک نام مانند `personal` یا `work` از شما می‌پرسد، اینکه حساب کدام موارد را با حساب‌های دیگر شما به اشتراک بگذارد و اینکه یک فرمان میان‌بر ساخته شود یا نه (مرحله ۳ را ببینید). سپس `claude auth login` را اجرا می‌کند که مرورگر را باز می‌کند تا وارد حساب شوید. روی یک ماشین راه دور، ممکن است مرورگر شما در صفحه `http://localhost:…/callback` بماند که بارگذاری نمی‌شود؛ آن نشانی را در ترمینال بچسبانید تا ccam ورود را کامل کند. وقتی ورود کامل شد، این را می‌بینید:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## ۲. تنظیمات فعلی خود را منتقل کنید

اگر پیش‌تر از Claude Code استفاده کرده‌اید، تنظیمات، `CLAUDE.md` و افزونه‌های شما در `~/.claude` و سرورهای MCP شما در `~/.claude.json` هستند. آن‌ها را در گروه مشترک `default` کپی کنید تا همه حساب‌ها از آن‌ها بهره ببرند:

```sh
ccam migrate
```

ccam نشان می‌دهد چه چیزی را کپی خواهد کرد و پیش از آن از شما می‌پرسد. فقط `~/.claude` را می‌خواند، هرگز اطلاعات ورود را کپی نمی‌کند و تاریخچه گفتگو را تنها وقتی کپی می‌کند که `--history` را اضافه کنید. برای جزئیات، [کپی کردن یک پیکربندی موجود](/docs/sharing#copy-an-existing-setup) را ببینید.

## ۳. حساب دوم را اضافه کنید

```sh
ccam add work
```

دستیار راه‌اندازی از نام می‌گذرد، چون آن را در خط فرمان داده‌اید، و گروه‌هایی را که می‌توان به آن‌ها پیوست نشان می‌دهد، با `default` علامت‌خورده. این بار به حساب یک فرمان میان‌بر مخصوص خودش بدهید، مانند `ccwork`، تا بتوانید آن را با یک کلمه اجرا کنید. دستیار راه‌اندازی پس از گروه‌ها این را می‌پرسد:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` را بزنید، سپس یک نام بنویسید یا Enter را بزنید تا نام پیشنهادی استفاده شود:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

سپس در مرورگر با حساب دیگر وارد شوید. ccam حساب و میان‌بر آن را تأیید می‌کند:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## ۴. بین حساب‌ها جابه‌جا شوید

`ccam work` برنامه Claude Code را با حساب کاری اجرا می‌کند. اگر از پرسش میان‌بر گذشتید، هر زمان که بخواهید می‌توانید میان‌بری اضافه کنید:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

اکنون هر یک از این فرمان‌ها Claude Code را با حساب کاری اجرا می‌کند (`2` شماره آن در `ccam list` است):

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

به‌طور پیش‌فرض، سرورهای MCP، `settings.json`، فایل سراسری `CLAUDE.md` و افزونه‌ها (agents، commands، skills و plugins) از طریق گروه `default` بین حساب‌ها مشترک‌اند. تاریخچه گفتگو مشترک نیست. اگر می‌خواهید گفتگوها را بین حساب‌ها ادامه دهید، آن را به گروه اضافه کنید:

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

**صندلی‌های Team و Enterprise.** انتخاب خودکار حساب (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد و، وقتی آن را روشن کنید، [`auto_switch`](/features/next#auto-switch)) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود، که ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، [پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را بخوانید و از سازمان خود و Anthropic بپرسید.

## از ~/.claude استفاده می‌کنید؟

به‌جای ورود دوباره، آن را ثبت کنید:

```sh
ccam import
```

پوشه در همان جایی که هست می‌ماند و `claude` معمولی مثل قبل کار می‌کند.
