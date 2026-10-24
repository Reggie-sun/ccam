---
title: انتخاب حساب برای اسکریپت‌ها
description: ccam pick حسابی را که بیشترین ظرفیت را دارد برای اسکریپت چاپ می‌کند، به کارگرهای هم‌زمان حساب‌های متفاوت می‌دهد و وقتی ظرفیتی نماند با کد 75 خارج می‌شود.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: a8baeb64de99
---

`ccam pick` نام حساب واردشده‌ای را که بیشترین ظرفیت را دارد چاپ می‌کند تا اسکریپت Claude Code را با آن آغاز کند. کارگرهایی که هم‌زمان آغاز می‌شوند، تا وقتی حساب‌های کافی با ظرفیت باشد، هر کدام حساب متفاوتی می‌گیرند؛ پس دسته‌ای از کارهای `claude -p` به جای خالی کردن یک حساب، میان حساب‌های شما پخش می‌شود.

**صندلی‌های Team و Enterprise.** انتخاب خودکار (`ccam next`، `--fallback`، `ccam pick`، کلید `n` در داشبورد) میان صندلی‌های همان سازمان Claude Team یا Enterprise هم جابه‌جا می‌شود. مرکز راهنمای Anthropic می‌گوید: «Each seat on a Team plan is assigned to one named individual, and only that person may use it» (هر صندلی در طرح Team به نام یک فرد مشخص اختصاص داده می‌شود و فقط همان شخص اجازه دارد از آن استفاده کند) و «Users can't route their work through another user's account to use that user's remaining usage» (کاربران نمی‌توانند کار خود را از طریق حساب کاربر دیگری انجام دهند تا از سهمیه مصرف باقی‌مانده آن کاربر استفاده کنند) ([طرح‌های Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ و همین‌طور برای [طرح‌های Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). بنابراین استفاده از انتخاب خودکار میان صندلی‌ها ممکن است شرایط Anthropic را نقض کند. پیش از آنکه به آن تکیه کنید، از سازمان خود و Anthropic بپرسید. برای طرح‌های شخصی (Pro، Max)، [توضیح شرایط در پرسش‌های متداول](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) را ببینید.

## نحوه استفاده

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

کدهای خروج: `0` با یک نام در خروجی استاندارد، `75` (`EX_TEMPFAIL`) وقتی هیچ حسابی ظرفیت ندارد یا هیچ حسابی وارد نشده است، و `1` برای خطاهای دیگر مانند یک گزینه نادرست.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--json` | حساب، مصرف ۵ ساعته و هفتگی آن، منبع عددها و دلیل انتخاب آن را چاپ می‌کند | خاموش |
| `--model <model>` | پنجره هفتگی آن مدل را هم در نظر می‌گیرد، مانند `opus` | هیچ |
| `--max-usage <1-100>` | حساب‌هایی را که مصرفشان زیر این مقدار است ترجیح می‌دهد | `auto_max_usage`، وگرنه ۱۰۰ |
| `--strategy headroom\|expiring\|balance` | کمترین مصرف ۵ ساعته، زودترین بازنشانی ۵ ساعته، یا سهمیه هفتگی‌ای که زودتر از همه بی‌استفاده می‌ماند | `auto_strategy`، وگرنه `headroom` |
| `--cross-org` | حساب‌های سازمان‌های دیگر را هم انتخاب می‌کند | `cross_org`، وگرنه خاموش |
| `--wait[=<duration>]` | وقتی هیچ حسابی ظرفیت ندارد، منتظر می‌ماند تا نخستین حساب بازنشانی شود، حداکثر به این مدت. مقدار پس از `=` می‌آید: `--wait 30m` با راهنمایی `--wait takes its value after '=': --wait=30m` رد می‌شود | خاموش؛ اگر تنها داده شود `5h` |
| `--dry-run` | انتخاب می‌کند بدون آنکه انتخاب را ثبت کند | خاموش |

## سازوکار

1. ccam حساب‌های واردشده را در سازمان حسابی که این پوشه اجرا می‌کرد فهرست می‌کند (حساب‌های کلید API و ورودهای Console سقفی ندارند و کنار گذاشته می‌شوند؛ [سازمان‌ها](/features/next#auto) را ببینید)؛ پوشه‌ای که به همان حساب Claude آن حساب وارد شده باشد هم کنار گذاشته می‌شود.
2. یک قفل می‌گیرد (`~/.local/state/ccam/picks.lock`)، تا انتخاب‌هایی که هم‌زمان اجرا می‌شوند یکی پس از دیگری انجام شوند.
3. مصرف همه حساب‌ها را هم‌زمان می‌پرسد، به همان شکلی که `ccam usage` می‌پرسد، با حافظه نهان ۶۰ ثانیه‌ای. پنجره‌ای که Claude Code پر گزارش کرده تا بازنشانی‌اش ۱۰۰٪ حساب می‌شود ([ثبت سقف مصرف](/features/usage#the-limit-record))، پس حسابی که همین حالا به سقفش رسیده بیرون می‌ماند، حتی وقتی عددهای ذخیره‌شده‌اش ظرفیت نشان می‌دهند.
4. ابتدا حساب‌هایی که در ۲ دقیقه گذشته انتخاب شده‌اند (`~/.local/state/ccam/picks/<name>`؛ انتخابی که تاریخش پس از زمان ساعت باشد حساب نمی‌شود) کنار گذاشته می‌شوند و [قواعد انتخاب](/features/next#auto) `ccam next --auto` از میان بقیه انتخاب می‌کنند: ابتدا حساب‌های زیر `auto_max_usage`، بر اساس راهبرد، و هرگز حسابی که به ۱۰۰٪ رسیده است.
5. اگر همه حساب‌هایی که ظرفیت دارند به‌تازگی انتخاب شده باشند، قواعد دوباره از میان همه آن‌ها انتخاب می‌کنند، پس کارگر چهارم روی سه حساب، یکی از آن‌ها را با دیگری شریک می‌شود.
6. انتخاب ثبت و قفل آزاد می‌شود و سپس نام چاپ می‌شود.

پنجره ۲ دقیقه‌ای اهمیت دارد، چون مدتی طول می‌کشد تا مصرف یک کار تازه در عددها دیده شود: بدون آن، پنج کارگری که با هم آغاز شده‌اند همگی یک حساب را حسابی با بیشترین ظرفیت می‌دیدند.

## ایمنی و موارد خاص

- `ccam pick` فقط مصرف را می‌خواند و زمان انتخاب را می‌نویسد؛ هرگز توکن‌ها را تمدید نمی‌کند و بیش از آنچه `ccam usage` می‌خواند، آن‌ها را نمی‌خواند.
- انتخاب، حساب را رزرو نمی‌کند. دو اسکریپت ممکن است پس از ۲ دقیقه باز هم روی یک حساب آغاز شوند؛ برای بازیابی وقتی آن حساب به سقفش می‌رسد از `--fallback` استفاده کنید.
- انتظار با `--wait` هیچ قفلی نگه نمی‌دارد، پس انتخاب‌های دیگر در این مدت ادامه می‌یابند.
- `ccam pick` به [`ccam folder`](/features/folder) توجهی نمی‌کند: فقط بر اساس مصرف انتخاب می‌کند.

## مطالب مرتبط

- [جابه‌جایی به حساب بعدی](/features/next)
- [اجرای دوباره هنگام رسیدن به سقف](/features/next#run-again-at-the-limit)
- [سقف مصرف](/features/usage)
- [حساب‌های کلید API](/features/api)
