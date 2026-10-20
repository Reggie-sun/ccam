---
title: استفاده از Remote Control با ccam
description: با ccam <name> rc، Remote Control در Claude Code را با هر حسابی اجرا کنید؛ گزینه‌های مؤثر، دلیل پیشنهاد ccam برای کپی سرورهای MCP مشترک و رفع مشکلات.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc` قابلیت Remote Control در Claude Code را با حساب `work` اجرا می‌کند تا بتوانید نشست‌های این رایانه را از claude.ai/code یا اپلیکیشن موبایل Claude هدایت کنید. این فرمان مانند `claude rc` کار می‌کند، با یک تفاوت: اگر حساب سرورهای MCP را به اشتراک بگذارد، ccam یک بار می‌پرسد که آیا آن‌ها را در `.claude.json` حساب کپی کند یا نه، چون نشست‌های Remote Control سرورها را از هیچ جای دیگری نمی‌خوانند.

## اجرای آن

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

یک [فرمان میان‌بر](/features/shortcuts) هم به همین شکل کار می‌کند: `ccwork rc` همان `ccam work rc` است. فرمان را در پوشه‌ای که می‌خواهید در آن کار کنید اجرا کنید، سپس claude.ai/code یا اپلیکیشن Claude را باز کنید، در حالی که با همان حساب Claude که `work` استفاده می‌کند وارد شده‌اید.

## گزینه‌ها

اینجا دو نوع گزینه به هم می‌رسند: گزینه‌های خود ccam که پیش از نام حساب می‌آیند یا در تنظیمات ccam قرار دارند، و گزینه‌های Remote Control که پس از `rc` می‌آیند و بدون تغییر به Claude Code می‌رسند.

### گزینه‌ها و تنظیمات ccam

| گزینه یا تنظیم | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--dry-run` (پیش از حساب) | محیط و فرمان را چاپ می‌کند، و اگر اعلان MCP وجود داشته باشد آن را هم. نه می‌پرسد و نه کپی می‌کند | خاموش |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | زمانی پرسیده می‌شود که حساب سرورهای MCP را به اشتراک بگذارد و `.claude.json` آن فهرست فعلی گروه را نداشته باشد. `y` آن‌ها را کپی می‌کند، `N` یا Enter بدون آن‌ها شروع می‌کند | `N` |
| `mcp_strategy = "sync"` در `~/.config/ccam/config.toml` | سرورهای گروه را پیش از هر اجرا، چه Remote Control باشد چه نباشد، در `.claude.json` کپی می‌کند، بنابراین هیچ‌وقت چیزی پرسیده نمی‌شود | `"flag"` |
| `ccam share work --mcp off` | اشتراک سرورهای MCP با حساب را متوقف می‌کند؛ در این صورت Remote Control فقط با سرورهای خود حساب شروع می‌شود و چیزی پرسیده نمی‌شود | روشن |

### گزینه‌های Remote Control

هر چه پس از `rc` بیاید متعلق به Claude Code است. این‌ها گزینه‌های `claude remote-control` در Claude Code 2.1.304 هستند؛ `ccam work rc --help` فهرست مربوط به نسخه شما را چاپ می‌کند.

| گزینه | کاری که انجام می‌دهد |
| --- | --- |
| `--name <name>` | نام نشست که در claude.ai/code نشان داده می‌شود |
| `--remote-control-session-name-prefix <prefix>` | پیشوند نام‌های ساخته‌شده برای نشست‌ها (پیش‌فرض: نام میزبان) |
| `-c`، `--continue` | دوباره به نشستی وصل می‌شود که Remote Control آخرین بار برای این پوشه ثبت کرده است، تا حدود ۴ ساعت |
| `--session-id <id>` | دوباره به یک نشست مشخص وصل می‌شود |
| `--permission-mode <mode>` | حالت مجوز برای نشست‌های جدید: `acceptEdits`، `auto`، `bypassPermissions`، `default`، `dontAsk` یا `plan` |
| `--spawn <mode>` | `same-dir` (پیش‌فرض)، `worktree` (یک git worktree برای هر نشست) یا `session` (یک نشست که با پایان آن خارج می‌شود) |
| `--capacity <N>` | بیشترین تعداد نشست هم‌زمان در حالت‌های `same-dir` و `worktree` (پیش‌فرض ۳۲) |
| `--[no-]create-session-in-dir` | هنگام شروع، نشست اول را در پوشه فعلی می‌سازد (پیش‌فرض روشن) |
| `--[no-]chrome` | Claude in Chrome برای نشست‌های جدید (پیش‌فرض خاموش) |
| `-d`، `--debug`، `--debug-file <path>`، `-v`، `--verbose` | خروجی اشکال‌زدایی و خروجی مفصل |

`rc` را درست پس از نام حساب بگذارید. ccam نخستین آرگومانی را که با `-` شروع نشود فرمان Claude Code در نظر می‌گیرد، بنابراین در `ccam work rc --name laptop` آن را `rc` می‌بیند، اما در `ccam work --model opus rc` آن را `opus` می‌بیند.

## کاری که ccam انجام می‌دهد

اجرا همان مسیر همیشگی را طی می‌کند ([اجرای Claude Code](/features/launch)) با یک تغییر در گام MCP:

1. ccam حساب را پیدا می‌کند، پیوندهای اشتراک آن را بررسی می‌کند و `claude` را پیدا می‌کند، مانند هر اجرای دیگر.
2. می‌بیند که فرمان `rc` یا `remote-control` است و `--mcp-config` را کنار می‌گذارد.
3. اگر حساب سرورهای MCP را به اشتراک بگذارد، `mcp.json` گروه را با `sync-stamp` در `~/.local/state/ccam/mcp/work/` مقایسه می‌کند. وقتی `.claude.json` از قبل سرورهای فعلی را داشته باشد، یا گروه هیچ سروری نداشته باشد، مستقیم به گام ۶ می‌رود.
4. در غیر این صورت سرورهایی را که نشست‌ها از دست می‌دهند نام می‌برد:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. با `y` از `.claude.json` نسخه پشتیبان می‌گیرد، قفل خود Claude Code را می‌گیرد، سرورهای گروه را به `mcpServers` اضافه می‌کند و سرورهای دیگر حساب را نگه می‌دارد. با `N`، بدون ترمینال، یا وقتی حساب در حال اجراست، چیزی را تغییر نمی‌دهد و این را اعلام می‌کند:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. نشست را ثبت می‌کند، `CLAUDE_CONFIG_DIR` را روی پوشه حساب تنظیم می‌کند و خودش را از طریق `exec` با `claude rc` جایگزین می‌کند.

پس از آن Remote Control به‌عنوان یک سرور در حال اجرا می‌ماند. هر نشستی که برای شما شروع می‌کند فرزند همان فرایند است، پس با حساب `work` اجرا می‌شود: با ورود، تنظیمات و تاریخچه گفتگوی `work`، و از سقف استفاده `work` کم می‌شود. تا وقتی سرور در حال اجراست، ccam حساب را در حال اجرا می‌داند.

اجرای آزمایشی تفاوت را نشان می‌دهد. یک اجرای عادی سرورهای گروه را به‌صورت آرگومان می‌برد، اما Remote Control نه:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## چرا این‌طور کار می‌کند

**Remote Control گزینه `--mcp-config` را نمی‌پذیرد.** در اجراهای عادی، ccam سرورهای گروه را به‌صورت `--mcp-config=<path>` به Claude Code می‌دهد و هرگز در فایل‌های Claude Code نمی‌نویسد. `claude remote-control` این گزینه را پیش از نام فرمانش رد می‌کند و پس از آن هم هیچ‌کدام را نمی‌پذیرد، بنابراین پیش از ccam 0.6.0، `ccam work rc` برای هر حسابی که سرورهای MCP را به اشتراک می‌گذاشت شکست می‌خورد. اکنون ccam این گزینه را کنار می‌گذارد.

**نشست‌های آن سرورها را فقط از `.claude.json` می‌خوانند.** نشست‌هایی که Remote Control شروع می‌کند گزینه‌ها را از خط فرمان به ارث نمی‌برند. تنها جایی که سرورهای MCP کاربر را پیدا می‌کنند `.claude.json` حساب است. این فایل جزئیات ورود حساب و وضعیت هر پوشه را هم نگه می‌دارد، بنابراین ccam نمی‌تواند آن را مانند `settings.json` به یک نسخه مشترک پیوند دهد. کپی کردن سرورها در آن تنها راه رساندن فهرست مشترک به نشست‌هاست.

**ccam پیش از نوشتن می‌پرسد.** با راهبرد پیش‌فرض flag، ccam قول می‌دهد به `.claude.json` دست نزند. Remote Control تنها موردی است که ناچار است آنجا بنویسد، پس اول می‌پرسد، نسخه پشتیبان نگه می‌دارد و به چیزی جز `mcpServers` دست نمی‌زند.

**وقتی حساب در حال اجراست نمی‌نویسد.** Claude Code در طول نشست `.claude.json` را بازنویسی می‌کند، و تغییری که از بیرون اعمال شود از بین می‌رود یا رونویسی می‌شود. به همین دلیل وقتی هر نشستی از حساب در حال اجرا باشد، از جمله یک سرور Remote Control دیگر، ccam کپی را انجام نمی‌دهد.

**برای هر تغییر یک بار می‌پرسد.** پس از کپی، `sync-stamp` زمان تغییر `mcp.json` گروه را نگه می‌دارد. ccam فقط وقتی دوباره می‌پرسد که آن فایل تغییر کند، برای مثال پس از `ccam mcp add`. اگر `N` پاسخ دهید، هنوز کپی‌ای وجود ندارد، پس در هر اجرای Remote Control می‌پرسد.

**`ccam mcp` کپی‌ها را به‌روز نگه می‌دارد.** در صورت تداخل نام، ورودی `.claude.json` بر `--mcp-config` برتری دارد، بنابراین یک کپی قدیمی در اجراهای عادی هم تعریف جدیدتر را پنهان می‌کند. از این رو `ccam mcp add`، `add-json`، `remove` و `import` کپی‌های همه حساب‌هایی را که در حال اجرا نیستند به‌روز می‌کنند:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**بدون ترمینال، پاسخ منفی است.** Remote Control اغلب از یک اسکریپت، یک سرویس یا `nohup` اجرا می‌شود، جایی که کسی نیست تا پاسخ دهد. در این صورت ccam به‌جای منتظر ماندن، بدون سرورها شروع می‌کند.

## مشکلات و راه‌حل‌ها

### سرورهای MCP مشترک در نشست‌های Remote Control نیستند

به پرسش `N` پاسخ داده شده، ترمینالی در کار نبوده، یا حساب در حال اجرا بوده است. Remote Control را متوقف کنید، نشست‌های دیگر حساب را ببندید و `ccam work rc` را دوباره در یک ترمینال اجرا کنید؛ پاسخ `y` بدهید. برای اینکه هرگز پرسیده نشود، `mcp_strategy = "sync"` را در `~/.config/ccam/config.toml` تنظیم کنید. سرورهای یک نشست تا وقتی در حال اجراست تغییر نمی‌کنند، پس پس از کپی یک نشست جدید شروع کنید.

### "Account work is running, so its .claude.json is left alone"

یک نشست از `work` هنوز باز است: یک `ccam work` عادی، یک میان‌بر یا یک سرور Remote Control دیگر. از آن خارج شوید و `ccam work rc` را دوباره اجرا کنید. نشست‌هایی که حساب دیگری اجرا می‌کند حساب نمی‌شوند.

### ccam هر بار می‌پرسد

هنوز سرورها را کپی نکرده‌اید. یک بار `y` پاسخ دهید، یا `mcp_strategy = "sync"` را تنظیم کنید. اگر هرگز سرورها را در Remote Control نمی‌خواهید، اشتراک MCP با حساب را متوقف کنید: `ccam share work --mcp off`.

### `ccam work rc` بلافاصله با خطایی درباره `--mcp-config` شکست می‌خورد

ccam شما قدیمی‌تر از 0.6.0 است و `--mcp-config` را به Remote Control می‌دهد. `ccam update` را اجرا کنید. همین خطا با یک ccam به‌روز یعنی گزینه‌ای مقدار‌دار پیش از `rc` آمده است، مانند `ccam work --model opus rc`، و ccam مقدار `opus` را فرمان گرفته است؛ `rc` را درست پس از نام حساب بگذارید.

### تغییری که با `ccam mcp` داده شده به Remote Control نمی‌رسد

`ccam mcp add` و `remove` کپی‌ها را فقط در حساب‌هایی به‌روز می‌کنند که در حال اجرا نیستند، و این را اعلام می‌کنند: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control را متوقف کنید و دوباره شروع کنید؛ ccam تغییر `mcp.json` را می‌بیند و دوباره می‌پرسد.

### کپی یکی از سرورهای خودم را جایگزین کرد

اگر حساب سرور خودش را با همان نام یکی از سرورهای گروه داشت، کپی آن را جایگزین کرده است. `.claude.json` قبلی در `~/.local/state/ccam/mcp/work/backups/` است که پنج نسخه آخر را نگه می‌دارد. برای اینکه نسخه خودتان را برای همه حساب‌ها نگه دارید، پیش از اجرای بعدی Remote Control آن را با `ccam mcp import work <server> --force` به گروه منتقل کنید.

### می‌خواهم سرورهای کپی‌شده را حذف کنم

`ccam mcp remove <server>` یک سرور را از گروه و از همه کپی‌ها حذف می‌کند. برای برگشتن به وضعیت پیش از کپی، جدیدترین نسخه پشتیبان را از `~/.local/state/ccam/mcp/work/backups/` بازگردانید، در حالی که هیچ نشستی از حساب در حال اجرا نیست.

### Remote Control می‌پرسد که آیا به پوشه اعتماد دارید، یا بدون ترمینال خارج می‌شود

اعتماد به پوشه برای هر حساب جداگانه در `.claude.json` ذخیره می‌شود، بنابراین پوشه‌ای که با `personal` به آن اعتماد کرده‌اید برای `work` تازه است. `ccam work` را یک بار در آن پوشه اجرا کنید و تأیید کنید، سپس Remote Control را آنجا شروع کنید.

### "You must be logged in with a Claude account that has a subscription"

Remote Control به ورود با اشتراک Claude نیاز دارد. ایمیل حساب را در `ccam list` بررسی کنید و در صورت نیاز با `ccam login work` دوباره وارد شوید.

### نشست روی گوشی نمایش داده نمی‌شود

اپلیکیشن Claude یا claude.ai/code با حساب Claude دیگری وارد شده است. آن را با ایمیلی که `ccam list` برای `work` نشان می‌دهد مقایسه کنید.

### ccam هشدار می‌دهد که Claude Code دیگری تاریخچه گفتگو را در این پوشه به اشتراک می‌گذارد

حساب دیگری از همان گروه تاریخچه در همان پوشه در حال اجراست. نشست‌هایی که از هر دو حساب همان گفتگو را ادامه دهند ممکن است تاریخچه را در هم بیامیزند. فقط وقتی `y` پاسخ دهید که روی گفتگوهای متفاوتی کار می‌کنند، یا Remote Control را در پوشه دیگری شروع کنید.

## مطالب مرتبط

- [اجرای Claude Code](/features/launch)
- [سرورهای MCP مشترک](/features/mcp)
- [اشتراک‌گذاری بین حساب‌ها](/features/share)
- [مدیریت حساب‌ها](/features/accounts)
