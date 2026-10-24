---
title: وقتی Claude Code به سقف مصرف می‌رسد چه کنیم
description: منتظر بازنشانی بمانید یا همان گفتگو را فوراً در حساب Claude Code دیگری ادامه دهید. سقف ۵ ساعته و هفتگی چگونه کار می‌کند، گام به گام.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

وقتی یک حساب به سقف مصرف خود می‌رسد (که اغلب rate limit نامیده می‌شود)، دو راه دارید: صبر کنید تا سقف بازنشانی شود، یا همین حالا همان گفتگو را در حساب دیگری که متعلق به خودتان است ادامه دهید. با ccam راه دوم فقط یک فرمان است: `ccam next --auto --continue`.

## Claude Code چه سقف‌هایی دارد؟

در طرح‌های Pro و Max، Claude Code یک [سقف نشست ۵ ساعته و یک سقف هفتگی](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) دارد. هر دو با Claude در وب و در اپلیکیشن‌ها مشترک‌اند و مصرف هم‌زمان از هر دو کم می‌شود. وقتی یکی تمام شود، Claude Code [پیامی نشان می‌دهد](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) مانند "You've hit your session limit · resets 3:45pm". عوض کردن مدل کمکی نمی‌کند، چون سقف‌ها همه مدل‌ها را در بر می‌گیرند.

## راه ۱: منتظر بازنشانی بمانید

اگر یک حساب دارید، صبر کنید. نسخه‌های جدید Claude Code [در همان نشست باز منتظر می‌مانند](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) و پس از بازنشانی سقف، کار را خودشان ادامه می‌دهند. `/usage` داخل Claude Code نشان می‌دهد چه مقدار از سقف‌های طرح خود را مصرف کرده‌اید.

## راه ۲: همین حالا در حساب دیگری ادامه دهید

اگر حساب دومی دارید، می‌توانید بی‌درنگ همان گفتگو را ادامه دهید. یک بار آن را راه‌اندازی کنید:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

اکنون خط وضعیت هنگام کار، حساب و مصرف ۵ ساعته آن را نشان می‌دهد، بنابراین رسیدن به سقف غافلگیرتان نمی‌کند.

وقتی به سقف رسیدید، از Claude Code خارج شوید (`/exit`)، سپس:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

گفتگو را در حسابی که بیشترین ظرفیت را دارد باز کنید:

```sh
ccam next --auto --continue
```

Claude Code با حساب `work` و با آخرین گفتگوی پوشه فعلی اجرا می‌شود، همراه با همه آنچه تا این لحظه گفته شده است. اشتراک تاریخچه (`ccam group set default --add history` در بالا) آن را یک گفتگو نگه می‌دارد؛ بدون آن، Claude Code گفتگو را در `work` کپی می‌کند و ccam به شما می‌گوید پاسخ اول چند توکن را دوباره می‌خواند.

با [پلاگین ccam](/docs/install#claude-code-plugin) در Claude Code، نشستی که در سقف خود متوقف شود یک اعلان نشان می‌دهد و پنجره پر را ثبت می‌کند، پس `ccam next --auto` آن حساب را پیش از بازنشانی‌اش دوباره انتخاب نمی‌کند، و `/ccam:switch` (یا `/ccam:switch work`) گفتگو را برای `ccam next --continue` یادداشت می‌کند، بی‌آنکه چیزی برای حدس زدن باقی بگذارد.

**صندلی‌های Team و Enterprise.** انتخاب خودکار (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود. مرکز راهنمای Anthropic می‌گوید: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (هر صندلی در طرح Team به نام یک فرد مشخص اختصاص داده می‌شود و فقط همان شخص اجازه دارد از آن استفاده کند) و «Users can't route their work through another user's account to use that user's remaining usage» (کاربران نمی‌توانند کار خود را از طریق حساب کاربر دیگری انجام دهند تا از سهمیه مصرف باقی‌مانده آن کاربر استفاده کنند) ([طرح‌های Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ و همین‌طور برای [طرح‌های Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). بنابراین استفاده از انتخاب خودکار میان صندلی‌ها ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، از سازمان خود و Anthropic بپرسید. برای طرح‌های شخصی (Pro، Max)، [توضیح شرایط در پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را ببینید.

## آیا ccam حساب‌ها را خودکار عوض می‌کند؟

خیر. ccam هنگام کار شما در حال اجرا نیست: Claude Code را اجرا می‌کند و خودش را با آن جایگزین می‌کند. شما تصمیم می‌گیرید کی جابه‌جا شوید، و `ccam next --auto` حسابی را انتخاب می‌کند که کمترین مصرف ۵ ساعته را دارد و از حساب‌هایی که به ۱۰۰٪ رسیده‌اند، وارد نشده‌اند یا متعلق به سازمان دیگری هستند می‌گذرد. با `relay = true` در `config.toml`، ccam به‌جای آن در کنار نشست منتظر می‌ماند، پس `/ccam:switch` خودش کار را ادامه می‌دهد و گزینه‌ها، حالت مجوز و سطح تلاش نشست را نگه می‌دارد. برخی ابزارهای دیگر به‌جای آن ورود یک نشست در حال اجرا را عوض می‌کنند؛ [ccam و دیگر ابزارهای تعویض حساب](/compare/account-switchers) تفاوت را توضیح می‌دهد.

## نکته‌هایی که باید در نظر داشته باشید

- `ccam next --continue` را در همان پوشه‌ای اجرا کنید که در آن کار می‌کردید؛ گفتگوها برای هر پوشه پروژه جداگانه ذخیره می‌شوند.
- از اجرای هم‌زمان دو حساب روی یک گفتگو پرهیز کنید. اگر چنین کاری کنید، ccam به شما هشدار می‌دهد.
- بدون تاریخچه مشترک، گفتگو در حساب دیگر کپی می‌شود و دو نسخه راه خود را می‌روند؛ نقطه‌های بازگشتی که `/rewind` کد را از آن‌ها بازمی‌گرداند در حساب اول می‌مانند.
- نخستین پیام در حساب جدید ممکن است بیش از معمول از سقف آن مصرف کند، چون آن حساب هنوز گفتگو را کش نکرده است.
- اینکه اجازه دارید این‌گونه بین حساب‌ها جابه‌جا شوید یا نه، به شرایط هر حساب بستگی دارد.

برای همه گزینه‌های `ccam usage` و `ccam next`، [سقف مصرف](/docs/usage-limits) را ببینید.
