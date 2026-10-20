---
title: مرجع فرمان‌ها
metaTitle: مرجع فرمان‌های ccam
description: همه فرمان‌ها و گزینه‌های ccam، از ccam add و ccam <name> تا ccam usage، ccam next و ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 5d3275878535
---

برای اینکه ببینید یک فرمان چه کاری انجام می‌دهد بدون آنکه واقعاً انجامش دهد، `--dry-run` را پیش از آن بگذارید، برای مثال `ccam --dry-run work`. این گزینه هنگام اجرای Claude Code (`ccam <name>`، `ccam run`، `ccam next`) و با `share`، `migrate`، `shortcut`، `update` و `doctor --fix` کار می‌کند. فرمان‌های دیگری که چیزی را تغییر می‌دهند، به‌جای اجرا شدن آن را رد می‌کنند.

<!-- commands -->

## اجرای Claude Code

`ccam <name|number> [arguments]` اجرای Claude Code با آن حساب را آغاز می‌کند. ccam آرگومان‌های پس از نام را نمی‌خواند؛ آن‌ها را به همان ترتیب به `claude` می‌دهد، و ورودی استاندارد، کد خروج و سیگنال‌ها دقیقاً مانند خود `claude` رفتار می‌کنند. گذاشتن `--` پس از نام اختیاری است.

`ccam` با آرگومان ولی بدون نام حساب، مانند `ccam -p "summarise"`، از حسابی که برای پوشه فعلی تنظیم شده و در غیر این صورت از حساب پیش‌فرض استفاده می‌کند. آن‌ها را با `ccam folder <name> [folder]` و `ccam use <name>` تعیین کنید.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` پیش از نام حساب برای اجراهای `-p`/`--print` است. وقتی حساب به سقف مصرفش برسد، ccam همان فرمان را دوباره با حساب واردشده‌ای که بیشترین ظرفیت را دارد اجرا می‌کند و فقط خروجی اجرایی را که به پایان رسیده چاپ می‌کند. [حساب‌ها بر اساس پوشه](/features/folder) و [اجرای دوباره هنگام رسیدن به سقف](/features/next#run-again-at-the-limit) را ببینید.

میان‌بری مانند `ccwork` حساب خود را به همین شکل اجرا می‌کند: `ccwork --resume` همان کاری را می‌کند که `ccam work --resume` انجام می‌دهد. [میان‌برها](/docs/commands#shortcuts) را ببینید.

## افزودن و حذف حساب‌ها

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` بدون نام پوشه، به دنبال `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام از شما می‌پرسد. پوشه‌های واردشده در همان جایی که هستند ثبت می‌شوند؛ ccam هرگز یک ورود را جابه‌جا یا کپی نمی‌کند.

تا وقتی Claude Code با یک حساب در حال اجراست، `ccam rename`، `ccam remove` و `ccam share` آن حساب را تغییر نمی‌دهند. ابتدا نشست را ببندید.

`ccam list --json` هر حساب را همراه با شماره جایگاه، میان‌بر، ایمیل و پوشه آن چاپ می‌کند. شیء `share` آن هر مورد (`mcp`، `settings`، `memory`، `extensions`، `history`) را به نام یک گروه یا `local` نگاشت می‌کند.

## میان‌برها

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

میان‌بر فرمانی است که نامش را خودتان انتخاب می‌کنید و Claude Code را با یک حساب اجرا می‌کند. مانند `ccam work`، همه آرگومان‌ها را بدون تغییر به Claude Code می‌دهد، بنابراین `ccwork --resume` و `ccwork -p "hi"` همان‌طور که انتظار دارید کار می‌کنند. میان‌بر یک پیوند به ccam در `~/.local/bin` است.

میان‌برها اختیاری‌اند. دستیار `ccam add` می‌پرسد که میان‌بری ساخته شود یا نه، `ccam add <name> --shortcut <command>` بدون پرسش یکی می‌سازد و `ccam list` میان‌بر هر حساب را نشان می‌دهد، یا `-` اگر میان‌بری نداشته باشد.

نام‌ها از نویسه‌های a-z، 0-9، `.`، `-` و `_` تشکیل می‌شوند. نام نمی‌تواند یکی از فرمان‌های ccam، میان‌بر حساب دیگری یا نام برنامه دیگری در `PATH` شما باشد.

`ccam rename` نام میان‌بر را نگه می‌دارد و `ccam remove` میان‌بر را همراه با حساب حذف می‌کند. `ccam shims install` پیوندهای همه حساب‌ها را دوباره می‌سازد، برای نمونه پس از اینکه ccam را جابه‌جا کرده‌اید.

## اشتراک‌گذاری

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

هر مورد یکی از این مقدارها را می‌پذیرد: `on` (گروه `default`)، `off` یا نام یک گروه. [اشتراک‌گذاری](/docs/sharing) را ببینید.

## کپی کردن یک پیکربندی موجود

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` فایل‌های `settings.json` و `CLAUDE.md`، افزونه‌ها و سرورهای MCP کاربر را از یک پوشه Claude Code در یک گروه کپی می‌کند تا همه حساب‌هایی که این موارد را به اشتراک می‌گذارند آن‌ها را داشته باشند. پوشه فقط خوانده می‌شود. گزینه‌های `--settings`، `--memory`، `--extensions`، `--mcp` و `--history` فقط مواردی را که نام می‌برید کپی می‌کنند. بدون آن‌ها، در ترمینال یک فهرست انتخاب به شما امکان انتخاب می‌دهد که در آن همه چیز به جز تاریخچه گفتگو علامت خورده است؛ در جاهای دیگر همه چیز به جز تاریخچه گفتگو کپی می‌شود. `-g` یا `--group` گروه را انتخاب می‌کند و `-y` یا `--yes` بدون پرسش کپی می‌کند و وقتی فایلی متفاوت باشد، نسخه گروه را نگه می‌دارد. همچنین ورودهای ذخیره‌شده سرورها را در فایل ورود خود هر حساب کپی می‌کند. [کپی کردن یک پیکربندی موجود](/docs/sharing#copy-an-existing-setup) را ببینید.

## سرورهای MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` سرورهایی را که فقط آن حساب دارد به گروه آن منتقل می‌کند. سروری که گروه از نامش برای تعریف دیگری استفاده می‌کند رد می‌شود؛ برای جایگزین کردن نسخه گروه `--force` را اضافه کنید، ccam پیش از آن از نسخه گروه پشتیبان می‌گیرد.

## مصرف و جابه‌جایی

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## نگهداری

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` پیش از جایگزین کردن ccam، checksum SHA-256 نسخه را بررسی می‌کند و `--version v0.2.0` یک نسخه مشخص را نصب می‌کند. ccam نصب‌شده از بسته `.deb` را باید با apt به‌روز کنید. [به‌روزرسانی](/docs/install#update) را ببینید.

## فایل‌ها

| مسیر | محتوا |
| --- | --- |
| `~/.config/ccam/config.toml` | حساب‌ها و آنچه به اشتراک می‌گذارند |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` هر حساب |
| `~/.local/share/ccam/groups/<group>/` | تنظیمات، سرورهای MCP، افزونه‌ها و تاریخچه مشترک |
| `~/.local/state/ccam/` | حافظه نهان مصرف و زمان‌های آخرین استفاده |
| `~/.local/bin/<shortcut>` | فرمان‌های میان‌بر مانند `ccwork`، که پیوندهایی به ccam هستند |

متغیرهای `XDG_CONFIG_HOME`، `XDG_DATA_HOME` و `XDG_STATE_HOME` محل این پوشه‌ها را تغییر می‌دهند. در Windows، `~` همان `%USERPROFILE%` است و میان‌برها `<shortcut>.exe` هستند.
