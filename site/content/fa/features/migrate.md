---
title: کپی کردن یک پیکربندی موجود
description: تنظیمات، CLAUDE.md، افزونه‌ها، سرورهای MCP و ورودهای MCP یک ~/.claude موجود را با ccam migrate در یک گروه ccam کپی کنید.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 7cee52e75749
---

اگر پیش از ccam از Claude Code استفاده کرده‌اید، پیکربندی شما در `~/.claude` و `~/.claude.json` است. `ccam migrate` آن را در یک گروه اشتراک کپی می‌کند تا هر حسابی که این موارد را از طریق آن گروه به اشتراک می‌گذارد با همان پیکربندی‌ای که از قبل دارید شروع کند. پوشه مبدأ فقط خوانده می‌شود و `claude` معمولی مثل قبل کار می‌کند. آن را یک بار پس از نخستین `ccam add` اجرا کنید و هر بار که حسابی اشتراک سرورهای MCP را آغاز می‌کند دوباره اجرا کنید.

## نحوه استفاده

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

در ترمینال، بدون گزینه‌های مورد و بدون `--yes`، ccam ابتدا یک فهرست انتخاب نشان می‌دهد:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (یا `x`) مورد زیر مکان‌نما را روشن و خاموش می‌کند، کلیدهای جهت‌نما یا `j`/`k` جابه‌جا می‌کنند، Enter برنامه را نشان می‌دهد و Esc، `q` یا Ctrl+C پیام `Cancelled.` را چاپ می‌کند. مواردی که مبدأ چیزی برایشان ندارد با `not found` مشخص می‌شوند. اگر علامت همه موارد را بردارید، ccam پیام `Nothing selected.` را چاپ می‌کند و متوقف می‌شود.

سپس ccam برنامه را چاپ می‌کند و پیش از کپی کردن از شما می‌پرسد:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

یک گروه فقط به حساب‌های خودش می‌رسد، و فقط برای مواردی که به اشتراک می‌گذارد، بنابراین خلاصه پایانی حساب‌های گروه را فهرست می‌کند. وقتی گروه یکی از موارد کپی‌شده را به اشتراک نمی‌گذارد، یا برخی حساب‌ها در آن نیستند، ccam فرمان‌های `ccam group set <group> --add <items>` و `ccam group join <group> <accounts>` را که این را تغییر می‌دهند چاپ می‌کند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `[folder]` | پوشه Claude Code که از آن کپی می‌شود | `~/.claude` |
| `--settings` | `settings.json` را کپی می‌کند | خاموش |
| `--memory` | فایل سراسری `CLAUDE.md` را کپی می‌کند | خاموش |
| `--extensions` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/` و `rules/` را کپی می‌کند | خاموش |
| `--mcp` | سرورهای MCP کاربر (`mcpServers` در `.claude.json`) را کپی می‌کند | خاموش |
| `--history` | تاریخچه گفتگو (`projects/`، `history.jsonl`) را کپی می‌کند | خاموش |
| `-g`، `--group string` | گروه اشتراکی که در آن کپی می‌شود | `default` |
| `-y`، `--yes` | بدون پرسش کپی می‌کند؛ در تعارض، نسخه گروه نگه داشته می‌شود | خاموش |
| `--dry-run` (سراسری) | برنامه را چاپ می‌کند و پیش از هر پرسش یا کپی متوقف می‌شود | خاموش |

بدون هیچ گزینه موردی، در ترمینال فهرست انتخاب تصمیم می‌گیرد. بیرون از ترمینال، یا با `--yes`، ccam همه چیز به جز تاریخچه گفتگو را کپی می‌کند. تاریخچه به‌طور پیش‌فرض کنار گذاشته می‌شود، چون `projects/` ممکن است بزرگ باشد.

`--group on` یعنی `default`. `--group off` و `--group local` با پیام `--group takes a group name, not "off"` رد می‌شوند. `ccam migrate` از `--dry-run` پیروی می‌کند و فهرست انتخاب در اجرای آزمایشی هم ظاهر می‌شود تا بتوانید پیش‌نمایش یک انتخاب را ببینید.

## سازوکار

`internal/cli/migrate.go` موارد را انتخاب می‌کند، سپس `migrate.Build` مبدأ و گروه را می‌خواند و یک برنامه برمی‌گرداند. تا پیش از `Plan.Apply`، یعنی پس از تأیید شما، چیزی نوشته نمی‌شود.

1. **یافتن مبدأ.** `checkSource` مسیری را که وجود ندارد، پوشه نیست یا هیچ مدخل شناخته‌شده‌ای از Claude Code ندارد رد می‌کند. `migrate.ClaudeJSON` فایل `.claude.json` متعلق به پوشه را پیدا می‌کند: درون آن، یا برای `~/.claude` که چنین فایلی درونش نیست، `~/.claude.json` در کنار آن. این همان جایی است که Claude Code وقتی `CLAUDE_CONFIG_DIR` تنظیم نشده آن را نگه می‌دارد.
2. **نگاشت موارد به مسیرها.** همه مسیرها از `share.Layout` می‌آیند، همان فهرست مجازی که گروه‌ها استفاده می‌کنند. سرورهای MCP با `mcp.LocalServers` از `.claude.json` خوانده می‌شوند که به‌جای تجزیه کل فایل از پویشگر `profile.TopLevelValue` استفاده می‌کند.
3. **کپی، هرگز جابه‌جایی.** فایل‌ها با `fsutil.AtomicWrite` یا `fsutil.CopyFile` به `~/.local/share/ccam/groups/<group>/` می‌روند؛ سرورهای MCP به `mcp.json` گروه می‌روند.

هر مورد یک مسیر مبدأ را به یک مسیر گروه نگاشت می‌کند:

<!-- motion:migrate-items -->

دو ردیف تاریخچه خط‌چین هستند، چون فقط با `--history` کپی می‌شوند. `.credentials.json` و `.claude.json` در `share.Never` فهرست شده‌اند و هرگز به گروه نمی‌رسند: اگر یکی از آن‌ها روزی در چیدمان ظاهر شود، `Build` یک خطای داخلی برمی‌گرداند و پیمایش پوشه آن‌ها را بر اساس نام رد می‌کند.

`buildEntry` به هر مسیر یک وضعیت می‌دهد و هر خط برنامه همین وضعیت را توصیف می‌کند:

<!-- motion:migrate-states -->

- **Missing**: در مبدأ نیست، یا جایی که پوشه انتظار می‌رود یک فایل است. برای هر مورد یک بار با `not found, skipped` فهرست می‌شود.
- **Linked**: مسیر مبدأ از قبل به نسخه گروه می‌رسد، برای مثال وقتی مبدأ پوشه یک پروفایل ccam است.
- **Same**: گروه از قبل محتوای یکسانی دارد.
- **New**: گروه آن را ندارد، پس به‌طور کامل کپی می‌شود.
- **Merge**: یک پوشه در هر دو طرف. `copyMissing` فقط فایل‌هایی را که گروه ندارد اضافه می‌کند؛ فایل‌های متفاوت همان‌طور که در گروه هستند می‌مانند و برنامه `group copies kept` را نشان می‌دهد. برای `plugins/`، فهرست‌های ثبت `installed_plugins.json` و `known_marketplaces.json` ادغام می‌شوند و مسیرهای مطلق افزونه‌ها طوری بازنویسی می‌شوند که به پوشه گروه اشاره کنند.
- **Conflict**: یک فایل در هر دو طرف با محتوای متفاوت. برای هر فایل انتخاب می‌کنید: `keep group`، `use mine`، به‌علاوه `merge JSON` برای فایل‌های `.json` (ادغامی عمیق که در آن مقدارهای متعارض به نفع گروه حل می‌شوند) یا `merge by line` برای `history.jsonl`. سرور MCP که تعریفش متفاوت باشد هم یک تعارض است، با `keep group` یا `use mine`. نخستین گزینه فهرست‌شده پیش‌فرض است وقتی Enter را می‌زنید؛ `--yes` همیشه نسخه گروه را نگه می‌دارد.

خود فهرست انتخاب `tui.MigratePicker` است که بر پایه مؤلفه `checklist` ساخته شده که `ccam add` هم از آن استفاده می‌کند؛ `migrate.Found` موارد ناموجود را با فراخوانی‌های lstat و فهرست سرورهای MCP مشخص می‌کند.

### ورودهای MCP

سرورهای MCP که با OAuth وارد می‌شوند توکن‌های خود را در فایل ورود، یعنی `.credentials.json` زیر `mcpOAuth`، نگه می‌دارند، نه در فهرست سرورها. Claude Code کلید هر مدخل را `<name>|` به‌علاوه ۱۶ نویسه نخست یک hash از نوع `sha256` روی `type`، `url` و `headers` سرور قرار می‌دهد. کلید به پوشه وابسته نیست، پس مدخل کپی‌شده در حساب دیگر هم کار می‌کند، به شرط آنکه تعریف سرور یکسان باشد. به همین دلیل `buildLogins` در `internal/migrate/logins.go` فقط مدخل‌های سرورهایی را کپی می‌کند که گروه در نهایت تعریف مبدأ آن‌ها را دارد: جدید، از قبل یکسان، یا تعارضی که در آن `use mine` را انتخاب کرده‌اید. ورودهای سرورهای افزونه (کلیدهایی که با `plugin:` شروع می‌شوند) وقتی افزونه‌ها را کپی کنید و مبدأ `plugins/` داشته باشد همراه آن‌ها می‌روند.

<!-- motion:migrate-signins -->

مدخل‌ها در `.credentials.json` خود هر حساب دریافت‌کننده نوشته می‌شوند، با `fsutil.AtomicWrite` و حالت `0600`. هرگز به گروه نمی‌روند. یک حساب فقط وقتی آن‌ها را دریافت می‌کند که `mcp` (یا برای سرورهای افزونه، `extensions`) را از طریق گروه مقصد به اشتراک بگذارد، و فقط کلیدهایی که ندارد اضافه می‌شوند: مدخل‌هایی که از قبل دارد، که شاید جدیدتر باشند، نگه داشته می‌شوند. خط `sign-ins` در برنامه، سرورها، حساب‌های دریافت‌کننده و حساب‌های ردشده را فهرست می‌کند. ccam یک حساب را رد می‌کند و نام آن را همراه با دلیل می‌گوید، وقتی آن حساب:

- با `ccam import` در همان محل ثبت شده باشد، چون ممکن است `claude` معمولی از فایل ورود آن استفاده کند؛
- در حال اجرا باشد، چون Claude Code فایل ورود خود را بازنویسی می‌کند (درست پیش از نوشتن دوباره بررسی می‌شود و با `started meanwhile` گزارش می‌شود)؛
- هنوز وارد نشده باشد؛
- ورود خود را در Keychain سیستم macOS نگه دارد، که ccam در آن نمی‌نویسد.

در حساب ردشده با `/mcp` وارد سرورها شوید. حسابی که بعداً اشتراک سرورهای MCP را آغاز کند، وقتی دوباره `ccam migrate` را اجرا کنید ورودها را دریافت می‌کند؛ خلاصه پایانی این را با ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` اعلام می‌کند.

## ایمنی و موارد خاص

- مبدأ فقط خوانده می‌شود. هیچ فایلی در آن جابه‌جا، پیوند یا تغییر داده نمی‌شود و ورود مبدأ هرگز نوشته نمی‌شود، از جمله در macOS که ممکن است از Keychain خوانده شود.
- پیش از آنکه یک فایل گروه جایگزین شود، با `use mine`، یک ادغام، یا سرورهای MCP جدید در یک `mcp.json` موجود، نسخه گروه در `<group>/.ccam-backup/<time>/` ذخیره می‌شود. نام پوشه در هر اجرا یکتاست، پس دو اجرا در یک ثانیه نسخه‌های پشتیبان یکدیگر را بازنویسی نمی‌کنند. خلاصه پایانی مسیر را چاپ می‌کند.
- با `--yes`، هر تعارض نسخه گروه را نگه می‌دارد و ccam می‌گوید چگونه به‌جای آن انتخاب کنید: همان فرمان را بدون `--yes` اجرا کنید.
- اگر چیزی تغییر نکند، ccam پیام `Nothing to copy.` را چاپ می‌کند؛ اگر همه تعارض‌ها نسخه گروه را نگه دارند، `Nothing to copy: the group keeps its own copies.` را.
- اگر کپی در میانه کار ناموفق شود، ccam با پیام `The source folder was not changed. Group files replaced so far are saved in` و پس از آن پوشه پشتیبان متوقف می‌شود.
- نبودن مبدأ با `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` پایان می‌یابد.

## مطالب مرتبط

- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [بررسی و تعمیر](/features/doctor)
- [چند حساب Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
