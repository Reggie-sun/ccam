---
title: مرجع فرمان‌ها
metaTitle: مرجع فرمان‌های ccam
description: همه فرمان‌ها و گزینه‌های ccam، از ccam add و ccam <name> تا ccam usage، ccam next و ccam doctor.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: 338ac11a0b5a
---

برای اینکه ببینید یک فرمان چه کاری انجام می‌دهد بدون آنکه واقعاً انجامش دهد، `--dry-run` را پیش از آن بگذارید، برای مثال `ccam --dry-run work`. این گزینه هنگام اجرای Claude Code (`ccam <name>`، `ccam run`، `ccam next`، `ccam attach`) و با `group`، `migrate`، `shortcut`، `update` و `doctor --fix` کار می‌کند. فرمان‌های دیگری که چیزی را تغییر می‌دهند، به‌جای اجرا شدن آن را رد می‌کنند.

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
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` بدون نام پوشه، به دنبال `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام از شما می‌پرسد. پوشه‌های واردشده در همان جایی که هستند ثبت می‌شوند؛ ccam هرگز یک ورود را جابه‌جا یا کپی نمی‌کند.

`ccam add --api` حسابی اضافه می‌کند که به جای ورود به Claude از کلید API استفاده می‌کند؛ ccam کلید را در هر اجرا می‌خواند و هرگز ذخیره‌اش نمی‌کند. [حساب‌های کلید API](/features/api) را ببینید.

تا وقتی Claude Code با یک حساب در حال اجراست، `ccam rename`، `ccam remove` و فرمان‌های `ccam group` آن حساب را تغییر نمی‌دهند. ابتدا نشست را ببندید.

`ccam list --json` هر حساب را همراه با شماره جایگاه، میان‌بر، ایمیل و پوشه آن چاپ می‌کند. آرایه `groups` آن گروه‌های حساب را فهرست می‌کند و شیء `share` آن هر مورد (`mcp`، `settings`، `memory`، `extensions`، `history`) را به گروهی که از آن می‌آید، یا به `local`، نگاشت می‌کند.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

یک گروه مواردی را بین حساب‌هایش به اشتراک می‌گذارد؛ یک حساب می‌تواند در چند گروه باشد، اما هر مورد فقط از یک گروه می‌آید. [اشتراک‌گذاری با گروه‌ها](/features/share) را ببینید.

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
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`، `--fallback` و `ccam pick` به یک شکل انتخاب می‌کنند: ابتدا حساب‌های زیر `auto_max_usage` (`config.toml`، پیش‌فرض ۱۰۰)، بر اساس `auto_strategy` (`headroom`، `expiring` یا `balance`)، با در نظر گرفتن پنجره هفتگی `--model`، هرگز حسابی که به ۱۰۰٪ رسیده است، و در سازمان حسابی که از آن شروع می‌کنند (`--cross-org` یا `cross_org = true` برای جست‌وجوی فراتر از آن). وقتی هیچ حسابی را نتوان انتخاب کرد، با کد ۷۵ خارج می‌شوند، یا با `--wait` منتظر می‌مانند تا نخستین حساب بازنشانی شود. `ccam next --continue` یک گفتگو را حتی بین حساب‌هایی که تاریخچه مشترک ندارند هم منتقل می‌کند: در این حالت Claude Code آن را در حساب جدید کپی می‌کند. برخلاف `ccam <name>`، `ccam next` گزینه‌های خودش را می‌خواند، پس گزینه‌های Claude Code پس از `--` می‌آیند؛ بدون آن ccam با `unknown flag` متوقف می‌شود. درون Claude Code، فرمان `/ccam:switch [account]` پلاگین یک گفتگو را به حساب دیگری می‌سپارد و `ccam hook` همان چیزی است که هوک‌های آن اجرا می‌کنند. [جابه‌جایی به حساب بعدی](/features/next)، [انتخاب حساب برای اسکریپت‌ها](/features/pick) و [نشست‌های پس‌زمینه همه حساب‌ها](/features/agents) را ببینید.

## نگهداری

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
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
| `~/.config/ccam/config.toml` | حساب‌ها، گروه‌ها و تنظیمات (`theme`، `usage_source`، `relay`، `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` هر حساب |
| `~/.local/share/ccam/groups/<group>/` | تنظیمات، سرورهای MCP، افزونه‌ها و تاریخچه مشترک |
| `~/.local/state/ccam/` | حافظه نهان مصرف، زمان‌های آخرین استفاده و یادداشت‌هایی که `/ccam:switch` و هوک سقف مصرف به جا می‌گذارند (`handoff/`) |
| `~/.local/bin/<shortcut>` | فرمان‌های میان‌بر مانند `ccwork`، که پیوندهایی به ccam هستند |

متغیرهای `XDG_CONFIG_HOME`، `XDG_DATA_HOME` و `XDG_STATE_HOME` محل این پوشه‌ها را تغییر می‌دهند. در Windows، `~` همان `%USERPROFILE%` است و میان‌برها `<shortcut>.exe` هستند.
