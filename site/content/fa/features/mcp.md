---
title: سرورهای MCP مشترک
description: برای هر گروه یک فهرست سرور MCP با ccam mcp add، add-json، list، remove و import نگه دارید که با --mcp-config داده یا در هر حساب همگام‌سازی می‌شود.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: 6d3b3c38c6cb
---

Claude Code سرورهای MCP کاربر را در `.claude.json` نگه می‌دارد، همان فایلی که ورود حساب را در خود دارد، بنابراین این فایل را نمی‌توان بین حساب‌ها به اشتراک گذاشت. ccam سرورهای مشترک را در یک فایل گروه، `mcp.json`، نگه می‌دارد و هنگام اجرا آن را به همه حساب‌های گروه می‌دهد. با `ccam mcp` یک سرور را یک بار برای همه حساب‌هایتان اضافه کنید و با `ccam mcp import` سروری را که داخل یک نشست اضافه کرده‌اید به حساب‌های دیگر منتقل کنید.

## نحوه استفاده

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

هر تغییر از اجرای بعدی هر حساب اعمال می‌شود؛ نشست‌های در حال اجرا سرورهایی را که با آن‌ها شروع شده‌اند نگه می‌دارند:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

افزودن نامی که گروه از قبل دارد، تعریف آن را جایگزین می‌کند و به‌جای `added`، `updated` چاپ می‌کند. سرورها فقط به حسابی می‌رسند که در آن گروه باشد و گروه MCP را به اشتراک بگذارد (`ccam group set <group> --add mcp`)؛ وقتی گروه این کار را نمی‌کند، `add` این را اعلام می‌کند.

## گزینه‌ها

| گزینه | کار آن | پیش‌فرض |
| --- | --- | --- |
| `-g`، `--group <group>` | گروه مقصد برای `add`، `add-json` و `remove` | تنها گروهی که حسابی MCP را از طریق آن به اشتراک می‌گذارد، وگرنه `default` |
| `-- <command> [args…]` (`add`) | این فرمان را به‌عنوان سرور stdio اجرا می‌کند | |
| `--url <url>` (`add`) | نشانی سرور راه دور (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | نوع انتقال سروری که با `--url` تعریف شده | `http` |
| `--header <'Name: value'>` (`add`) | سرآیند HTTP، قابل تکرار | هیچ |
| `-e`، `--env <KEY=VALUE>` (`add`) | متغیر محیطی برای سرور stdio، قابل تکرار | هیچ |
| `--force` (`import`) | سرور هم‌نام و متفاوت گروه را جایگزین می‌کند؛ از `mcp.json` قدیمی پشتیبان گرفته می‌شود | خاموش |
| `--keep-local` (`import`) | ورودی‌ها را در `.claude.json` حساب نگه می‌دارد | خاموش |
| `-h`، `--help` | راهنمای هر فرمان | |
| `--dry-run` (سراسری) | فقط `mcp list` آن را می‌پذیرد؛ `add`، `add-json`، `remove` و `import` آن را رد می‌کنند و چیزی را تغییر نمی‌دهند | خاموش |

`add` یا به یک فرمان پس از `--` نیاز دارد یا به `--url`؛ در غیر این صورت با `a command to run (after --) or --url is required` شکست می‌خورد. `--header` برای سرورهای `--url` و `--env` برای سرورهای فرمانی است. `add-json` همان JSON فرمان `claude mcp add-json` را می‌پذیرد و تعریف باید شیئی با `command` یا `url` باشد. `import` همیشه از گروهی استفاده می‌کند که حساب MCP را از طریق آن به اشتراک می‌گذارد و `list` همه گروه‌ها را نشان می‌دهد، بنابراین `-g` روی آن‌ها اثری ندارد.

## چگونه کار می‌کند

فایل‌های درگیر:

- `~/.local/share/ccam/groups/<group>/mcp.json`: سرورهای مشترک، `{"mcpServers": {…}}`، همان قالبی که `--mcp-config` می‌خواند. این فایل با `fsutil.AtomicWrite` و حالت `0600` نوشته می‌شود، چون ورودی‌های سرور اغلب کلیدهای API دارند.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: فایل خود Claude Code با سرورهای حساب. ccam فقط مقدار `mcpServers` آن را با پویشگر `profile.TopLevelValue` می‌خواند و هرگز کل فایل را در مسیر اجرا تجزیه نمی‌کند.
- `~/.local/state/ccam/mcp/<name>/`: فایل `seen`، فایل `sync-stamp` و نسخه‌های پشتیبان `.claude.json`.

دو راهبرد وجود دارد که با `mcp_strategy` در `~/.config/ccam/config.toml` تعیین می‌شوند. پیش‌فرض، راهبرد flag است:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` گروه MCP حساب را پیدا می‌کند. `mcp.FlagArgs` با یک `stat` بررسی می‌کند که `mcp.json` گروه وجود دارد، بدون اینکه آن را تجزیه کند.
2. `Plan.InsertArgs` گزینه `--mcp-config=<path>` را درست پس از `claude` و پیش از آرگومان‌های خود شما قرار می‌دهد، بنابراین `--mcp-config` یا `--strict-mcp-config` خود شما همچنان روی آن اعمال می‌شوند.
3. `mcp.NewLocal` بررسی می‌کند که آیا `.claude.json` حساب از اجرای قبلی تغییر کرده است و فقط در آن صورت سرورهای آن را تجزیه می‌کند.
4. `syscall.Exec` ccam را با Claude Code جایگزین می‌کند.

مسیر به‌صورت یک آرگومان، `--mcp-config=<path>`، داده می‌شود نه دو آرگومان. `--mcp-config` در Claude Code چند مقدار می‌پذیرد، بنابراین `claude --mcp-config a.json "prompt"` پرامپت را به‌عنوان فایل پیکربندی دوم می‌خواند و با `MCP config file not found` شکست می‌خورد. شکل تک‌آرگومانی همان‌جا تمام می‌شود و پرامپت یا زیرفرمانی مانند `mcp list` که پس از آن بیاید کار می‌کند. Claude Code سرورهای `--mcp-config` را علاوه بر سرورهای خود حساب بارگذاری می‌کند؛ در صورت تداخل نام، ورودی `.claude.json` برنده است. راهبرد flag فقط برای Remote Control در `.claude.json` می‌نویسد.

`claude remote-control` (`rc`) گزینه `--mcp-config` را نمی‌پذیرد و نشست‌های آن سرورها را فقط از `.claude.json` می‌خوانند. برای `ccam work rc`، ccam این گزینه را کنار می‌گذارد و در ترمینال می‌پرسد که آیا سرورهای گروه را، با نسخه پشتیبان، در `.claude.json` کپی کند یا نه، همان‌طور که راهبرد `sync` انجام می‌دهد. تا وقتی کپی‌ای وجود نداشته باشد می‌پرسد و پس از آن فقط بعد از تغییر `mcp.json` گروه دوباره می‌پرسد ([Remote Control](/features/launch#remote-control)). چون چنین کپی‌ای در صورت تداخل نام بر `--mcp-config` برتری دارد، `ccam mcp add`، `add-json`، `remove` و `import` کپی‌هایی را که حساب‌های دیگر گروه نگه می‌دارند به‌روز می‌کنند:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

سروری که داخل یک نشست با `claude mcp add` اضافه شود فقط در `.claude.json` همان حساب قرار می‌گیرد. ccam در آن لحظه نمی‌تواند آن را ببیند، چون Claude Code جای آن را گرفته است، بنابراین در اجرای بعدی بررسی می‌کند. فایل `seen` زمان تغییر `.claude.json` و نام سرورهایی را که قبلاً گزارش شده‌اند نگه می‌دارد. وقتی زمان متفاوت باشد، ccam نام‌ها را با گروه مقایسه می‌کند و برای هر سرور جدید یک بار یک راهنمایی چاپ می‌کند، بدون اینکه منتظر ورودی بماند:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

با `mcp_strategy = "sync"`، سرورها در `claude mcp list` و `/mcp` به شکل سرورهای عادی کاربر دیده می‌شوند و ccam پیش از اجرا آن‌ها را در `.claude.json` هر حساب می‌نویسد:

<!-- motion:mcp-sync -->

1. `mcp.Sync` زمان تغییر `mcp.json` گروه را با `sync-stamp` مقایسه می‌کند. اگر تغییری نکرده باشد، `.claude.json` حساب باز نمی‌شود.
2. قفل خود Claude Code را می‌گیرد، یعنی پوشه `.claude.json.lock` کنار فایل که با `mkdir` ساخته می‌شود. تا ۲ ثانیه صبر می‌کند و قفلی را که بیش از ۱۰ ثانیه از عمرش گذشته باشد کهنه در نظر می‌گیرد، همان کاری که کتابخانه قفل Claude Code انجام می‌دهد.
3. `.claude.json` فعلی را در `backups/` در پوشه وضعیت ذخیره می‌کند و ۵ نسخه جدیدتر را نگه می‌دارد.
4. همه سرورهای گروه را در `mcpServers` قرار می‌دهد، سرورهایی را که قبلاً نوشته و گروه دیگر ندارد حذف می‌کند، سرورهای خود حساب را نگه می‌دارد و فایل را به‌صورت اتمی با حالت اصلی‌اش می‌نویسد.
5. مهر زمانی جدید و نام‌هایی را که مدیریت می‌کند ثبت می‌کند و Claude Code بدون `--mcp-config` شروع می‌شود.

وقتی نشستی از حساب در حال اجرا باشد، همگام‌سازی با پیام `ccam: account work is already running, skipping MCP sync` انجام نمی‌شود، چون Claude Code در حین اجرا `.claude.json` را بازنویسی می‌کند. همگام‌سازی ناموفق `ccam: MCP sync failed (launching anyway): …` را چاپ می‌کند و با این حال اجرا را ادامه می‌دهد. راهنمایی درباره سرورهایی که داخل نشست اضافه شده‌اند فقط مربوط به راهبرد flag است.

`ccam mcp import` سرورهای خود یک حساب را به گروه آن منتقل می‌کند:

<!-- motion:mcp-import -->

ccam سرورهای حساب را با `mcp.LocalServers` می‌خواند و هر کدام را با سرور هم‌نام گروه با استفاده از `mcp.SameDef` مقایسه می‌کند، که ترتیب کلیدها و فاصله‌ها را نادیده می‌گیرد. سروری که گروه ندارد یا با همان تعریف دارد کپی می‌شود. سرور متفاوت رد می‌شود، چون تعریف گروه همان چیزی است که همه حساب‌های دیگر از آن استفاده می‌کنند:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

با `--force` سرور گروه جایگزین می‌شود و `mcp.json` قبلی ابتدا در `.ccam-backup/<time>/mcp.json` داخل پوشه گروه ذخیره می‌شود. سپس `mcp.RemoveLocal` ورودی‌های کپی‌شده را از `.claude.json` حساب حذف می‌کند، با همان قفل و پس از گرفتن پشتیبان در پوشه وضعیت. این کار لازم است، چون در صورت تداخل نام، ورودی حساب بر ورودی گروه برنده می‌شود. `--keep-local` آن‌ها را باقی می‌گذارد.

## ایمنی و حالت‌های خاص

- **ورودها.** سرورهایی که با OAuth وارد می‌شوند توکن‌های خود را در `.credentials.json` هر حساب (`mcpOAuth`) نگه می‌دارند، نه در فهرست سرورها، بنابراین اشتراک یک سرور ورود آن را به اشتراک نمی‌گذارد. در هر حساب یک بار با `/mcp` وارد شوید، یا ورودها را با [`ccam migrate`](/features/migrate) از `~/.claude` کپی کنید. گروه هرگز فایل ورود دریافت نمی‌کند.
- **حساب‌های در حال اجرا.** `import` در گروه کپی می‌کند، اما ورودی‌های خود حسابی را که در حال اجراست باقی می‌گذارد و این را اعلام می‌کند؛ پس از پایان نشست آن را دوباره اجرا کنید. همگام‌سازی از حساب‌های در حال اجرا می‌گذرد.
- **کپی‌های مربوط به Remote Control.** به کپی حسابی که در حال اجراست دست زده نمی‌شود: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` با `mcp_strategy = "sync"` چیزی نیاز به به‌روزرسانی ندارد، چون هر اجرا همگام‌سازی می‌کند.
- **فایل‌های خراب.** اگر `.claude.json` JSON معتبر نباشد، همگام‌سازی به آن دست نمی‌زند و خطا را گزارش می‌کند. اگر قفل در ۲ ثانیه گرفته نشود، ccam پیام `could not lock … (Claude Code is writing the file)` را گزارش می‌کند.
- **حساب‌هایی که MCP را به اشتراک نمی‌گذارند.** `import` رد می‌کند و از شما می‌خواهد ابتدا حساب را در گروهی بگذارید که MCP را به اشتراک می‌گذارد (`ccam group join <group> side`).
- **نام گروه‌ها.** `-g` فقط نام گروه می‌پذیرد؛ این نام به پوشه‌ای زیر `groups/` تبدیل می‌شود، بنابراین نام‌هایی مانند `../x` رد می‌شوند.

## مطالب مرتبط

- [اشتراک‌گذاری بین حساب‌ها](/features/share)
- [استفاده از Remote Control با ccam](/guides/remote-control)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [کپی کردن یک پیکربندی موجود](/features/migrate)
- [همه فرمان‌ها](/docs/commands)
