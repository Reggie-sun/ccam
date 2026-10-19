---
title: اجرای Claude Code
description: Claude Code را با هر حسابی از طریق ccam <name|number>، ccam run یا فهرست انتخاب ccam اجرا کنید؛ آرگومان‌ها بدون تغییر به claude می‌رسند و --dry-run فرمان را نشان می‌دهد.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

`ccam work` اجرای Claude Code را با حساب `work` آغاز می‌کند و هر چه پس از نام بیاید همان‌طور که هست به `claude` می‌رسد. هر جا که `claude` را تایپ می‌کردید از آن استفاده کنید: در ترمینال، در یک اسکریپت یا با ورودی‌ای که از طریق pipe داده می‌شود.

## نحوه استفاده

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

نام یک حساب یا شماره آن در `ccam list` را بدهید. ccam آرگومان‌های پس از نام را نمی‌خواند، بنابراین `ccam work --help` راهنمای Claude Code را نشان می‌دهد. گذاشتن `--` پس از نام اختیاری است. ورودی استاندارد، کد خروج و سیگنال‌ها مانند `claude` رفتار می‌کنند.

وقتی آرگومان اول با `-` شروع شود و یکی از گزینه‌های خود ccam نباشد، ccam حساب پیش‌فرض را با همه آرگومان‌ها اجرا می‌کند. حساب پیش‌فرض همان حسابی است که با `ccam use` تعیین شده، یا اگر حسابی تعیین نکرده باشید، حسابی که کمترین شماره را دارد.

### فهرست انتخاب

`ccam` بدون آرگومان، وقتی ورودی و خروجی استاندارد ترمینال باشند فهرست انتخاب را باز می‌کند و در غیر این صورت `ccam list` را چاپ می‌کند. این فهرست شماره، ایمیل، مصرف ۵ ساعته، آخرین استفاده و موارد مشترک هر حساب را نشان می‌دهد؛ مصرف هر حساب به محض پاسخ دادن آن پر می‌شود.

| کلید | کار آن |
| --- | --- |
| `Enter` | اجرای Claude Code با حساب انتخاب‌شده |
| `1` تا `9` | اجرای حسابی که آن شماره را دارد |
| `Up`/`Down`، `k`/`j` | جابه‌جا کردن انتخاب |
| `a` | بستن فهرست انتخاب و اجرای دستیار `ccam add` |
| `s` | ویرایش موارد مشترک حساب انتخاب‌شده (`Space` یا `x` روشن و خاموش می‌کند، `Enter` اعمال می‌کند، `Esc` برمی‌گردد) |
| `d` | پیش‌فرض کردن حساب انتخاب‌شده |
| `r` | دریافت دوباره مصرف، بدون حافظه نهان کوتاه‌مدت |
| `q`، `Esc`، `Ctrl+C` | خروج |

فهرست انتخاب پیش از انجام کار بسته می‌شود، بنابراین ترمینال پیش از شروع Claude Code به حالت عادی برمی‌گردد. پس از `d` یا تغییر اشتراک‌گذاری، فهرست دوباره باز می‌شود. تغییر اشتراک‌گذاری، مانند `ccam share`، ابتدا برنامه را نشان می‌دهد و از شما می‌پرسد.

## گزینه‌ها

| گزینه | کار آن | پیش‌فرض |
| --- | --- | --- |
| `--dry-run` | محیط و فرمانی را که ccam اجرا می‌کرد چاپ می‌کند و خارج می‌شود | خاموش |
| `-v`، `--version` | `ccam <version>` را چاپ می‌کند (فقط به‌عنوان آرگومان اول) | |
| `-h`، `--help` | راهنمای ccam، به شکل `ccam --help` یا `ccam run --help` | |

قواعد آرگومان‌ها:

- ccam گزینه‌های خودش را فقط پیش از نام حساب می‌خواند. پس از نام، همه آرگومان‌ها متعلق به Claude Code هستند، از جمله `--dry-run` و `--help`.
- `--dry-run` به شکل‌های `ccam --dry-run work`، `ccam run --dry-run work` و `ccam --dry-run` (حساب پیش‌فرض) کار می‌کند. این گزینه پذیرفته می‌شود، نه رد.
- `ccam run <name|number> [claude args]` شکل بلند `ccam <name>` است. `ccam run` بدون نام از حساب پیش‌فرض استفاده می‌کند.
- نام حساب‌ها نمی‌تواند نام فرمان‌های ccam مانند `add`، `list`، `run` یا `usage` باشد، بنابراین `ccam <word>` هرگز مبهم نیست.

اجرای آزمایشی یک خط چاپ می‌کند که می‌توانید آن را در shell بچسبانید:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## چگونه کار می‌کند

هر اجرا، از خط فرمان تا `exec`، همین گام‌ها را طی می‌کند:

<!-- motion:launch -->

هر جعبه هنگام اجرای گام خود روشن می‌شود و خط‌های زیر آن، محیط و آرگومان‌هایی هستند که Claude Code با آن‌ها شروع می‌شود.

1. `cmd/ccam/main.go` تابع `app.Main` را فراخوانی می‌کند. اگر ccam با نام دیگری اجرا شده باشد، یک [میان‌بر](/features/shortcuts) است و مستقیم به اجراکننده می‌رود.
2. `app.parse` پیش از دخالت cobra به `os.Args` نگاه می‌کند. فقط `ccam` به‌تنهایی (فهرست انتخاب)، نام یک فرمان ccam (`config.Reserved`)، `-h`/`--help` و `-v`/`--version` به `cli.Execute` می‌روند. هر چیز دیگری یک اجراست.
3. `launcher.Launch` فایل `~/.config/ccam/config.toml` را می‌خواند، تنها فایلی که تجزیه می‌کند.
4. `launcher.Prepare` حساب را با نام یا شماره پیدا می‌کند، یک `--` آغازین را حذف می‌کند و فایل اجرایی `claude` را با `ResolveClaude` پیدا می‌کند. نتیجه جست‌وجو در PATH همراه با زمان تغییر فایل در `~/.local/state/ccam/claude-path` ذخیره می‌شود، بنابراین ccam در هر اجرا `PATH` را جست‌وجو نمی‌کند. اگر فایل تغییر کرده یا وجود نداشته باشد، جست‌وجوی تازه‌ای انجام می‌شود. مسیر `claude_bin` در `config.toml` همان‌طور که هست استفاده می‌شود. ccam فایل `claude` را که خود ccam باشد (با مقایسه inode) رد می‌کند، چون به یک حلقه بی‌پایان می‌انجامد.
5. `checkLinks` پوشه حساب `~/.local/share/ccam/profiles/<name>/` را با تنظیمات اشتراک‌گذاری آن با استفاده از `share.Check` مقایسه می‌کند، که فقط `lstat` و `readlink` را فراخوانی می‌کند، یعنی یک یا دو فراخوانی سیستمی برای هر مورد. پیوندی که وجود نداشته باشد در همان لحظه ساخته می‌شود. هر چیزی که ممکن است به از دست رفتن داده بینجامد اجرا را متوقف می‌کند.
6. `applyMCP` سرورهای MCP مشترک را اضافه می‌کند (پایین‌تر را ببینید).
7. `session.Register` نشست را ثبت می‌کند. ccam یک نشانگر به نام PID خود، که زمان شروع فرایند را در خود دارد، در `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` می‌نویسد (یا در `~/.local/state/ccam/locks/` اگر `XDG_RUNTIME_DIR` تعریف نشده باشد). برای حسابی که تاریخچه گفتگو را به اشتراک می‌گذارد، `historyCheck` ابتدا در `history-<group>-<hash>`، که hash آن از پوشه فعلی گرفته می‌شود، به دنبال نشستی در حال اجرا می‌گردد و آنجا هم ثبت می‌کند. `exec` همان PID را نگه می‌دارد، بنابراین نشانگر سپس به خود فرایند Claude Code اشاره می‌کند. از نشانگر استفاده می‌شود نه قفل نگه‌داشته‌شده، چون فرایندهایی که Claude Code اجرا می‌کند قفل را به ارث می‌برند و پس از پایان نشست آن را نگه می‌دارند.
8. `Plan.Exec` عنوان ترمینال را `claude · <name>` می‌گذارد (وقتی خروجی خطای استاندارد ترمینال باشد)، زمان اجرا را بدون `fsync` در `~/.local/state/ccam/last-used/<name>` می‌نویسد و `syscall.Exec` را فراخوانی می‌کند.

Claude Code با `CLAUDE_CONFIG_DIR` برابر با پوشه حساب و `CCAM_PROFILE` برابر با نام حساب اجرا می‌شود. ccam هر مقداری از این متغیرها را که به ارث برده باشد حذف می‌کند. برای پوشه `~/.claude` که با `ccam import` روی macOS ثبت شده، `CLAUDE_SECURESTORAGE_CONFIG_DIR` را هم با مقدار خالی تنظیم می‌کند تا آن حساب و `claude` معمولی از یک ورود Keychain استفاده کنند.

### چرا exec

`syscall.Exec` فرایند ccam را با `claude` جایگزین می‌کند، بنابراین ccam در حافظه باقی نمی‌ماند، کلیدها، سیگنال‌ها و کد خروج متعلق به خود Claude Code هستند و آرگومان‌ها به شکل آرایه می‌مانند، هرگز رشته shell نمی‌شوند. کارهای پس از نشست، مانند تشخیص سرورهایی که با `claude mcp add` اضافه شده‌اند، در اجرای بعدی انجام می‌شود: ccam زمان تغییر `.claude.json` را با `~/.local/state/ccam/mcp/<name>/seen` مقایسه می‌کند و فقط وقتی فایل تغییر کرده باشد آن را تجزیه می‌کند.

### مسیر یک خط فرمان

سه خط فرمان، سه مسیر در `internal/app`:

<!-- motion:launch-route -->

میان‌بر بر اساس `argv[0]` و `ccam <name>` بر اساس آرگومان اول تشخیص داده می‌شود و هر دو به `syscall.Exec` می‌رسند. فقط یک فرمان ccam مانند `ccam usage` به cobra می‌رسد.

### سرورهای MCP مشترک

با راهبرد پیش‌فرض `flag`، ccam سرورهای MCP گروه را به‌صورت یک آرگومان، پیش از آرگومان‌های شما، به Claude Code می‌دهد:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` فقط وقتی فایل `mcp.json` گروه وجود داشته باشد یک آرگومان برمی‌گرداند و `Plan.InsertArgs` آن را درست پس از `claude` قرار می‌دهد. گزینه‌های `--mcp-config` یا `--strict-mcp-config` خود شما همچنان روی آن اعمال می‌شوند. این یک آرگومان واحد است، `--mcp-config=<path>`، چون `--mcp-config` در Claude Code چند مقدار می‌پذیرد: اگر به شکل دو آرگومان نوشته شود، پرامپت بعدی را هم به‌عنوان فایل پیکربندی می‌خواند.

با `mcp_strategy = "sync"` در `config.toml`، ccam به‌جای آن سرورها را در `.claude.json` حساب می‌نویسد و اگر حساب از قبل در حال اجرا باشد این کار را انجام نمی‌دهد.

### Remote Control

`ccam work rc` (یا `ccam work remote-control`) قابلیت Remote Control در Claude Code را با حساب `work` اجرا می‌کند. `claude remote-control` گزینه `--mcp-config` را پیش از فعل خود رد می‌کند و پس از آن هم هیچ‌کدام را نمی‌پذیرد، و نشست‌های آن سرورهای MCP را فقط از `.claude.json` حساب می‌خوانند. بنابراین راهبرد flag اینجا طور دیگری کار می‌کند:

<!-- motion:launch-rc -->

1. `remoteControl` نخستین آرگومان پس از حساب را که با `-` شروع نشود فعل Claude Code در نظر می‌گیرد. وقتی این فعل `remote-control` یا `rc` باشد، `mcp.FlagArgs` چیزی اضافه نمی‌کند.
2. `mcp.Synced` زمان تغییر `mcp.json` گروه را با `sync-stamp` حساب مقایسه می‌کند. اگر `.claude.json` از قبل سرورهای فعلی گروه را داشته باشد، یا گروه هیچ سروری نداشته باشد، Remote Control بدون پرسش شروع می‌شود.
3. در غیر این صورت ccam سرورهایی را که نشست‌ها از دست می‌دهند نام می‌برد و می‌پرسد:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. با `y`، `mcp.Sync` آن‌ها را همان‌طور کپی می‌کند که راهبرد `sync` کپی می‌کند: زیر قفل Claude Code، پس از گرفتن نسخه پشتیبان و با نگه داشتن سرورهای خود حساب. با `N`، یا بدون ترمینال، Remote Control بدون آن‌ها شروع می‌شود و ccam به `mcp_strategy = "sync"` اشاره می‌کند که آن‌ها را پیش از هر اجرا کپی می‌کند.

ccam در هر اجرای Remote Control می‌پرسد تا سرورها کپی شوند، و پس از کپی فقط وقتی دوباره می‌پرسد که `mcp.json` گروه باز هم تغییر کند. در اجراهای عادی بعدی، کپی موجود در `.claude.json` بر `--mcp-config` برتری دارد، بنابراین `ccam mcp add`، `remove` و `import` کپی‌هایی را که حساب‌های دیگر نگه می‌دارند به‌روز می‌کنند ([سرورهای MCP مشترک](/features/mcp)).

### بودجه ۱۰ میلی‌ثانیه

ccam پیش از شروع Claude Code حداکثر ۱۰ میلی‌ثانیه (p95) اضافه می‌کند. مسیر اجرا هرگز درخت فرمان‌های cobra را نمی‌سازد، Bubble Tea را راه نمی‌اندازد و اتصال شبکه باز نمی‌کند، و اگر مجموع زمان توابع `init()` بسته‌ها بیش از ۵ میلی‌ثانیه شود، `TestPackageInitIsCheap` شکست می‌خورد. `scripts/perf-check.sh` زمان `ccam --dry-run work`، `ccam work` و یک میان‌بر را با `hyperfine` اندازه می‌گیرد و بالای ۱۰ میلی‌ثانیه شکست می‌خورد، و `BenchmarkLaunchPath` در `internal/launcher` گام‌ها را در یک حلقه اجرا می‌کند. برای دیدن زمان هر گام، `CCAM_TRACE=1` را تنظیم کنید:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## ایمنی و حالت‌های خاص

- نام یا شماره ناشناخته با `no account "<name>"` متوقف می‌شود و شما را به `ccam list` ارجاع می‌دهد. اگر هیچ حسابی وجود نداشته باشد، ccam از شما می‌خواهد با `ccam add <name>` یکی اضافه کنید.
- اگر یک مسیر مشترک به فایل عادی تبدیل شده یا به جای اشتباهی اشاره کند، ccam با `stopped because of share link problems:` متوقف می‌شود، مشکلات را فهرست می‌کند و از شما می‌خواهد `ccam doctor --fix` را اجرا کنید. ccam در مسیر اجرا داده‌ای را جابه‌جا نمی‌کند.
- اگر Claude Code دیگری که تاریخچه گفتگو را با این حساب به اشتراک می‌گذارد (همان گروه) در همان پوشه در حال اجرا باشد، ccam هشدار می‌دهد که تاریخچه ممکن است در هم آمیخته شود و می‌پرسد `Run anyway? [y/N]`. بدون ترمینال، هشدار را نشان می‌دهد و ادامه می‌دهد.
- `--dry-run` از همگام‌سازی MCP، اعلان سرور جدید، بررسی تاریخچه، نشانگر نشست، عنوان و زمان آخرین استفاده می‌گذرد و exec نمی‌کند. با این حال پیوندهای اشتراک را بررسی می‌کند و پیوندی را که وجود نداشته باشد، مانند یک اجرای عادی، می‌سازد.
- اجرای `ccam work rc` وقتی `work` از قبل در حال اجراست به `.claude.json` دست نمی‌زند، این را اعلام می‌کند و Remote Control را بدون سرورهای گروه شروع می‌کند؛ پس از پایان نشست دوباره آن را اجرا کنید تا سرورها کپی شوند. با `--dry-run` پیام را چاپ می‌کند اما نه می‌پرسد و نه کپی می‌کند.
- با `mcp_strategy = "sync"`، اگر همگام‌سازی ناموفق باشد `ccam: MCP sync failed (launching anyway): …` چاپ می‌شود و Claude Code با سرورهایی که از قبل داشت شروع می‌شود.

## مطالب مرتبط

- [فرمان‌های میان‌بر](/features/shortcuts)
- [مدیریت حساب‌ها](/features/accounts)
- [مرجع فرمان‌ها](/docs/commands)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [CLAUDE_CONFIG_DIR چگونه کار می‌کند](/guides/claude-config-dir)
