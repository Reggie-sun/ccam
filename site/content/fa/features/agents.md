---
title: نشست‌های پس‌زمینه همه حساب‌ها
description: ccam agents عامل‌های پس‌زمینه و نشست‌های باز Claude Code همه حساب‌های شما را در یک جدول فهرست می‌کند و ccam attach یکی از آن‌ها را با شناسه‌اش باز می‌کند.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: "06281fb13119"
---

Claude Code نشست‌های پس‌زمینه (`claude --bg`، نمای عامل‌ها) را برای هر حساب جداگانه نگه می‌دارد: `claude agents` فقط نشست‌های حسابی را نشان می‌دهد که با آن اجرا شده است. `ccam agents` از همه حساب‌ها هم‌زمان می‌پرسد و آن‌ها را با هم فهرست می‌کند، و `ccam attach <id|name>` یک نشست پس‌زمینه را با هر حسابی که آن را آغاز کرده باز می‌کند.

## نحوه استفاده

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` وضعیت نشست پس‌زمینه (`working`، `blocked`، `done`، `failed`، `stopped`) و کاری است که فرایند آن انجام می‌دهد (`busy`، `waiting`، `idle`)؛ یک نشست تعاملی باز `open` را نشان می‌دهد.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--all` | نشست‌های پس‌زمینه‌ای را هم که تمام شده، ناموفق شده یا متوقف شده‌اند فهرست می‌کند | خاموش |
| `--json` | نشست‌ها را به‌صورت JSON چاپ می‌کند، هر کدام همراه با `account`ی که به آن تعلق دارد | خاموش |

`ccam attach <id|name> [claude args]` همه چیز پس از شناسه را بدون تغییر، از جمله پرچم‌ها، به `claude attach` می‌دهد: `ccam attach 49ec --dry-run` گزینه `--dry-run` را به Claude Code می‌دهد. برای چاپ فرمان، `--dry-run` را پیش از شناسه بگذارید (`ccam --dry-run attach 49ec` یا `ccam attach --dry-run 49ec`).

## سازوکار

1. ccam برای هر حساب، `claude agents --json` را (با `--all` اگر داده شده باشد) با `CLAUDE_CONFIG_DIR` همان حساب اجرا می‌کند، به‌صورت موازی و هر کدام با مهلت ۱۵ ثانیه. این کار فهرست نشست‌های خود حساب را می‌خواند و سرویس پس‌زمینه Claude Code را راه‌اندازی نمی‌کند.
2. هر ورودی نام حساب را می‌گیرد و فهرست بر اساس زمان شروع، از جدیدترین، مرتب می‌شود. حسابی که پاسخ ندهد در خطای استاندارد نام برده می‌شود و کنار گذاشته می‌شود، و سپس ccam با کد ۱ خارج می‌شود تا اسکریپت شکست را با فهرستی خالی اشتباه نگیرد.
3. `ccam attach` در پوشه هر حساب به دنبال `jobs/<id>` می‌گردد، جایی که Claude Code یک نشست پس‌زمینه را نگه می‌دارد، و `claude attach <id>` را با آن حساب اجرا می‌کند، به همان شکلی که `ccam <name>` Claude Code را آغاز می‌کند. واژه‌ای که پیشوند شناسه نباشد، بدون توجه به بزرگی و کوچکی حروف، با نام نشست‌ها تطبیق داده می‌شود؛ واژه‌ای که شبیه شناسه باشد هرگز نام در نظر گرفته نمی‌شود.

## ایمنی و موارد خاص

- شناسه‌های نشست پس‌زمینه به یک حساب تعلق دارند: `claude attach` که با حساب دیگری اجرا شود نمی‌تواند آن‌ها را پیدا کند، و `ccam attach` همین مشکل را حل می‌کند.
- اگر نشست‌هایی در دو حساب پیشوند شناسه یا واژه‌ای از نام مشترک داشته باشند، `ccam attach` آن‌ها را فهرست می‌کند و شناسه بلندتر یا نام را می‌خواهد.
- ccam سرورهای MCP مشترک خود را با `--mcp-config=<file>` در یک آرگومان به Claude Code می‌دهد که `claude agents` و `claude attach` آن را می‌پذیرند؛ نشست‌های پس‌زمینه هنگام راه‌اندازی دوباره همچنان از همان فایل استفاده می‌کنند.
- سرویس پس‌زمینه یک حساب کلید API کلید خود را از shellی می‌گیرد که آن را آغاز کرده است. نخستین نشست پس‌زمینه چنین حسابی را از طریق ccam آغاز کنید (`ccam gw --bg …`) تا سرویس کلید را داشته باشد.

## مطالب مرتبط

- [اجرای Claude Code](/features/launch)
- [انتخاب حساب برای اسکریپت‌ها](/features/pick)
- [داشبورد](/features/dashboard)
