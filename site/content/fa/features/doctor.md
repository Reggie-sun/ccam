---
title: بررسی و تعمیر
description: پیوندهای اشتراک، مجوزهای فایل، ورودها و نسخه Claude Code را با ccam doctor بررسی کنید و آنچه را تعمیرش بی‌خطر است با ccam doctor --fix تعمیر کنید.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` نصب ccam شما را بررسی می‌کند: فایل اجرایی Claude Code، پوشه، ورود و پیوندهای اشتراک هر حساب، سرورهای MCP مشترک، نقطه پایانی مصرف و فرمان‌های میان‌بر شما. `ccam doctor --fix` آنچه را بدون از دست رفتن داده ممکن است تعمیر می‌کند و هر فایلی را که جایگزین می‌کند ابتدا به یک پوشه پشتیبان منتقل می‌کند. وقتی `ccam <name>` به دلیل مشکل پیوند اشتراک متوقف می‌شود، پس از بازگرداندن دستی فایل‌ها یا پس از به‌روزرسانی Claude Code آن را اجرا کنید.

## نحوه استفاده

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

هر یافته یک خط با یک علامت است: `✓` درست، `·` اطلاعات، `!` هشدار، `✗` خطا. یافته‌های مربوط به حساب‌ها با نام حساب شروع می‌شوند.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

با `--fix`، هر یافته‌ای که تعمیر شده با `→ fixed` تمام می‌شود:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--fix` | آنچه را تعمیرش بی‌خطر است تعمیر می‌کند (از فایل‌های تغییریافته پشتیبان گرفته می‌شود) | خاموش |
| `--dry-run` (سراسری) | همه بررسی‌ها را اجرا می‌کند اما حتی با `--fix` چیزی را تعمیر نمی‌کند | خاموش |

`ccam doctor` آرگومانی نمی‌پذیرد. از `--dry-run` پیروی می‌کند: `ccam --dry-run doctor --fix` همان یافته‌های `ccam doctor` را چاپ می‌کند و چیزی را تغییر نمی‌دهد.

کد خروج وقتی دست‌کم یک یافته `✗` تعمیرنشده بماند `1` و در غیر این صورت `0` است. هشدارها و اطلاعات هرگز اجرا را ناموفق نمی‌کنند. وقتی اجرا بدون `--fix` ناموفق باشد، ccam با یک خط خالی و `` `ccam doctor --fix` repairs what it can. `` تمام می‌شود. پیکربندی‌ای که بارگذاری نشود هم با `1` و یک خطای `ccam:` تمام می‌شود.

## نحوه کار

`internal/cli/doctor.go` فایل `~/.config/ccam/config.toml` را بارگذاری می‌کند و `doctor.Run` را در `internal/doctor/doctor.go` فرا می‌خواند که بررسی‌ها را به ترتیبی ثابت اجرا می‌کند و یافته‌ها را برمی‌گرداند. `--fix` به‌صورت `Options.Fix` منتقل می‌شود که در اجرای آزمایشی false است.

<!-- motion:doctor-checks -->

بخش میانی برای هر حساب یک بار اجرا می‌شود. بررسی‌ها و کاری که `--fix` برای هر یک انجام می‌دهد:

1. **نسخه Claude Code** (`checkClaude`). فایل اجرایی `claude` را همان‌طور که اجراکننده پیدا می‌کند پیدا می‌کند و `claude --version` را اجرا می‌کند. اگر نسخه با `2.1.302`، یعنی نسخه‌ای که ساختار پوشه ccam با آن تأیید شده (`doctor.LayoutVersion`)، شروع نشود، `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` را اضافه می‌کند. نبودن فایل اجرایی یا ناموفق بودن `--version` خطاست. تعمیری ندارد.
2. **اطلاعات ورود در محیط** (`checkAuthEnv`). وقتی `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` تنظیم شده باشد هشدار می‌دهد: هر حسابی که از آن shell اجرا شود آن را به ارث می‌برد و Claude Code آن را به ورود خود حساب ترجیح می‌دهد. تعمیری ندارد؛ خودتان آن را حذف کنید.
3. **پوشه** (`checkProfile`). نبودن پوشه حساب خطاست و بررسی‌های آن حساب را پایان می‌دهد. مجوزهایی که هر بیتی برای گروه یا دیگران داشته باشند هشدارند؛ `--fix` فرمان `chmod 0700` را اجرا می‌کند. Windows بیت مجوز ندارد، پس این بررسی و بررسی `0600` در ادامه آنجا انجام نمی‌شوند؛ پوشه‌ای که در پروفایل کاربری شماست با فهرست دسترسی خودش محافظت می‌شود.
4. **فایل ورود** (`checkCredentials`). `.credentials.json` با `lstat` بررسی می‌شود. پیوند نمادین خطاست، چون هر حساب باید توکن خودش را داشته باشد، و تعمیر نمی‌شود. مجوزی غیر از `0600` هشدار است؛ `--fix` فرمان `chmod 0600` را اجرا می‌کند. سپس ورود خوانده می‌شود (روی macOS ابتدا از Keychain). نبودن ورود هشداری است که به `ccam login <name>` ارجاع می‌دهد. توکن تمدید منقضی‌شده خطای `login expired` است؛ توکن دسترسی منقضی‌شده فقط اطلاعات است، چون Claude Code آن را در اجرای بعدی تمدید می‌کند. هیچ‌کدام از این‌ها تعمیر نمی‌شوند.
5. **پیوندهای اشتراک** (`share.Check`). هر مسیر در `share.Layout` را با تنظیمات اشتراک حساب مقایسه می‌کند و فقط از `lstat` و `readlink` استفاده می‌کند. چهار نوع مشکل را گزارش می‌دهد: پیوند گم‌شده، پیوندی که به جای دیگری اشاره می‌کند، فایل واقعی در جایی که باید پیوند باشد، و پیوند باقی‌مانده برای موردی که اکنون محلی است. پیوند گم‌شده هشدار است و بقیه خطا هستند.
6. **ورودی‌های ناشناخته** (`checkUnknownEntries`). نام‌های سطح بالای پوشه حساب را که ccam نمی‌شناسد (`share.Known`) فهرست می‌کند، به‌عنوان نشانه‌ای از اینکه ساختار Claude Code تغییر کرده است. فقط اطلاعات است.
7. **مسیرهای افزونه** (`checkPlugins`). فهرست‌های افزونه (`installed_plugins.json`، `known_marketplaces.json`) مسیرهای مطلق ذخیره می‌کنند. در `plugins/` هر گروه، مسیرهایی که از پوشه یک حساب می‌گذرند با حذف یا تغییر نام آن حساب خراب می‌شوند. `--fix` آن‌ها را با `profile.RewritePluginPaths`، که فایل فهرست را به‌صورت اتمی می‌نویسد، به پوشه گروه بازنویسی می‌کند.
8. **سرورهای MCP** (`checkMCP`). `mcp.json` هر گروه را یک بار بارگذاری می‌کند. فایلی که خوانده نشود خطاست؛ سروری که از `mcp.Validate` رد نشود هشدار است. تعمیری ندارد.
9. **قالب پاسخ مصرف** (`checkUsageFormat`). بدون حافظه نهان و با محدودیت معمول ۳ ثانیه، از نقطه پایانی مصرف برای حساب‌ها یکی پس از دیگری پرس‌وجو می‌کند تا یکی پاسخ دهد. اگر پاسخ دیگر قابل تجزیه نباشد، هشدار می‌دهد که `ccam usage` به حافظه نهان خط وضعیت برمی‌گردد. این تنها بررسی‌ای است که از شبکه استفاده می‌کند.
10. **میان‌برها** (`checkShims`). فقط وقتی حسابی میان‌بر داشته باشد. پیوندهای گم‌شده در `~/.local/bin` هشدارند؛ `--fix` تابع `shim.Install` را فرا می‌خواند که فقط پیوندهای گم‌شده را می‌سازد و هرگز به فایل‌های دیگر دست نمی‌زند. پوشه bin که در `PATH` نباشد هشدار است و تعمیر نمی‌شود.
11. **پوشه‌ها** (`checkFolders`). فقط وقتی ورودی‌های [`ccam folder`](/features/folder) وجود داشته باشند. ورودی‌ای که حسابش وجود ندارد خطاست و ورودی‌ای که پوشه‌اش از بین رفته هشدار است. هیچ‌کدام تعمیر نمی‌شوند، چون حذف یک ورودی حسابی را که آنجا اجرا می‌شود تغییر می‌دهد.

### تعمیر پیوندهای اشتراک

مسیر اجرا پیوند گم‌شده را همان‌جا می‌سازد، اما برای هر چیزی که ممکن است داده را از بین ببرد متوقف می‌شود و شما را به اینجا ارجاع می‌دهد:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

برای هر مشکل پیوند اشتراک، `--fix` تابع `share.Build` را با تنظیمات اشتراک حساب هم به‌عنوان وضعیت فعلی و هم وضعیت هدف فرا می‌خواند و سپس `Plan.Apply` را. `Build` همه مراحل را از آنچه روی دیسک هست دوباره به دست می‌آورد، همان برنامه‌ای که `ccam share` استفاده می‌کند، بنابراین تعمیر از همان قواعد روشن کردن اشتراک‌گذاری پیروی می‌کند.

<!-- motion:doctor-repair -->

در اینجا `settings.json` یک فایل واقعی است. چون گروه از قبل `settings.json` متفاوتی دارد، تعارض با انتخاب پیش‌فرض خود، `KeepGroup`، حل می‌شود و نسخه گروه باقی می‌ماند. فایل حساب به `.ccam-backup/<time>/settings.json` داخل پوشه حساب منتقل می‌شود و `os.Symlink` مسیر را به گروه پیوند می‌دهد. موارد دیگر:

- اگر گروه هنوز نسخه‌ای نداشته باشد، فایل حساب به گروه منتقل می‌شود و نخستین محتوای آن می‌شود.
- برای یک پوشه، فایل‌هایی که گروه ندارد به گروه منتقل می‌شوند و بقیه به پشتیبان می‌روند.
- برای `history.jsonl`، پیش‌فرض `merge by line` است: خطوط گروه و سپس خطوط جدید شما. نسخه قبلی گروه در `.ccam-backup/<time>/group/` ذخیره می‌شود.
- پیوند به گروه دیگر جایگزین می‌شود. پیوند به پوشه‌ای بیرون از ccam به پشتیبان می‌رود و اگر گروه خالی باشد، محتوای آن نخستین محتوای گروه می‌شود.
- پیوندی که برای موردی که اکنون محلی است باقی مانده حذف می‌شود و حساب یک نسخه از محتوای گروه را می‌گیرد.

به جز حذف پیوندها، هر مرحله داده را به گروه یا پوشه پشتیبان منتقل می‌کند. هیچ چیزی پاک نمی‌شود.

## ایمنی و موارد خاص

- `ccam doctor` بدون گزینه هرگز چیزی را تغییر نمی‌دهد. برای اطمینان در اسکریپت‌ها، `--dry-run` را همراه `--fix` به کار ببرید.
- `doctor --fix` بررسی نمی‌کند که حسابی در حال اجرا هست یا نه. پیش از تعمیر پیوندهای اشتراک یک حساب، Claude Code آن را ببندید، چون نشست همچنان در آن فایل‌ها می‌نویسد.
- اگر تعمیر اشتراک ناموفق باشد، ccam برای آن حساب `could not repair share links: <error>` را گزارش می‌دهد و بررسی ورودی‌های ناشناخته آن را رد می‌کند. حساب‌های دیگر همچنان بررسی می‌شوند.
- اگر `ccam share` در میانه کار ناموفق شود، ابتدا تنظیمات هدف را ذخیره می‌کند و به شما می‌گوید کار را با `ccam doctor --fix` تمام کنید.
- `.credentials.json` پیوندشده، ورود منقضی‌شده و اطلاعات ورود محیطی گزارش می‌شوند اما هرگز تغییر داده نمی‌شوند.
- بازنویسی فهرست‌های افزونه به‌صورت اتمی نوشته می‌شود اما در `.ccam-backup/` کپی نمی‌شود.

## مطالب مرتبط

- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [کپی کردن یک پیکربندی موجود](/features/migrate)
- [سقف مصرف و جابه‌جایی بین حساب‌ها](/docs/usage-limits)
- [فرمان‌ها](/docs/commands)
