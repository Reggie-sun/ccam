---
title: کپی کردن یک پیکربندی موجود
description: تنظیمات، CLAUDE.md، افزونه‌ها و سرورهای MCP یک ~/.claude موجود را با ccam migrate در یک گروه ccam کپی کنید.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: 7bca9efb6eb0
---

اگر پیش از ccam از Claude Code استفاده کرده‌اید، پیکربندی شما در `~/.claude` و `~/.claude.json` است. `ccam migrate` آن را در یک گروه اشتراک کپی می‌کند تا هر حسابی که این موارد را از طریق آن گروه به اشتراک می‌گذارد با همان پیکربندی‌ای که از قبل دارید شروع کند. پوشه مبدأ فقط خوانده می‌شود و `claude` معمولی مثل قبل کار می‌کند. آن را یک بار پس از نخستین `ccam add` اجرا کنید؛ وقتی `~/.claude.json` سرورهای MCP داشته باشد، نخستین `ccam add` در نکته خود نامشان را می‌برد.

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
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
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

1. **یافتن مبدأ.** `checkSource` مسیری را که وجود ندارد، پوشه نیست یا هیچ مدخل شناخته‌شده‌ای از Claude Code ندارد رد می‌کند. برای سرورهای MCP به‌تنهایی، `~/.claude` خالی یا ناموجود اشکالی ندارد، به شرط آنکه `~/.claude.json` سرور داشته باشد. `migrate.ClaudeJSON` فایل `.claude.json` متعلق به پوشه را پیدا می‌کند: درون آن، یا برای `~/.claude` که چنین فایلی درونش نیست، `~/.claude.json` در کنار آن. این همان جایی است که Claude Code وقتی `CLAUDE_CONFIG_DIR` تنظیم نشده آن را نگه می‌دارد.
2. **نگاشت موارد به مسیرها.** همه مسیرها از `share.Layout` می‌آیند، همان فهرست مجازی که گروه‌ها استفاده می‌کنند. سرورهای MCP با `mcp.LocalServers` از `.claude.json` خوانده می‌شوند که به‌جای تجزیه کل فایل از پویشگر `profile.TopLevelValue` استفاده می‌کند.
3. **کپی، هرگز جابه‌جایی.** فایل‌ها با `fsutil.AtomicWrite` یا `fsutil.CopyFile` به `~/.local/share/ccam/groups/<group>/` می‌روند، که یک فایل موقت کنار مقصد می‌نویسد و نامش را تغییر می‌دهد، پس کپی‌ای که متوقف شود هرگز فایل نیمه‌کاره‌ای بر جا نمی‌گذارد. سرورهای MCP زیر قفل `mcp.json` گروه به آن می‌روند: سروری که ccam دیگری در این مدت اضافه کند می‌ماند، و کپی‌هایی که برای [Remote Control](/guides/remote-control) ساخته شده‌اند پس از آن به‌روز می‌شوند.

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

سرورهای MCP که با OAuth وارد می‌شوند توکن‌های خود را در فایل ورود، یعنی `.credentials.json` زیر `mcpOAuth`، نگه می‌دارند، نه در فهرست سرورها. ccam آن‌ها را کپی نمی‌کند. یک ورود کپی‌شده همان توکن تمدید نسخه اصلی را دارد و سرور آن توکن را در نخستین تمدید عوض می‌کند: آنگاه همه دارندگان دیگر، از جمله `~/.claude` مبدأ، از آن سرور خارج می‌شوند. بنابراین خط `sign-ins` در برنامه، برای هر سرور کپی‌شده‌ای که مبدأ به آن وارد شده است، حساب‌هایی از گروه را نام می‌برد که خودشان با `/mcp` وارد می‌شوند، و خلاصه پایانی آن را تکرار می‌کند. مدخلی که توکن ندارد ورود به حساب نمی‌آید و فهرست نمی‌شود، و همچنین حسابی که از قبل ورود خودش را دارد، یا سروری که گروه تعریف خودش را برایش نگه می‌دارد. هیچ فایل ورودی نوشته نمی‌شود.

## ایمنی و موارد خاص

- مبدأ فقط خوانده می‌شود. هیچ فایلی در آن جابه‌جا، پیوند یا تغییر داده نمی‌شود و ورود مبدأ هرگز نوشته نمی‌شود، از جمله در macOS که ممکن است از Keychain خوانده شود.
- پیش از آنکه یک فایل گروه جایگزین شود، با `use mine`، یک ادغام، یا سرورهای MCP جدید در یک `mcp.json` موجود، نسخه گروه در `<group>/.ccam-backup/<time>/` ذخیره می‌شود. نام پوشه در هر اجرا یکتاست، پس دو اجرا در یک ثانیه نسخه‌های پشتیبان یکدیگر را بازنویسی نمی‌کنند. خلاصه پایانی مسیر را چاپ می‌کند.
- با `--yes`، هر تعارض نسخه گروه را نگه می‌دارد و ccam می‌گوید چگونه به‌جای آن انتخاب کنید: همان فرمان را بدون `--yes` اجرا کنید.
- اگر چیزی تغییر نکند، ccam پیام `Nothing to copy.` را چاپ می‌کند؛ اگر همه تعارض‌ها نسخه گروه را نگه دارند، `Nothing to copy: the group keeps its own copies.` را.
- اگر کپی در میانه کار ناموفق شود، ccam با پیام `The source folder was not changed. Group files replaced so far are saved in` و پس از آن پوشه پشتیبان متوقف می‌شود. اجرای دوباره فایل‌هایی را که کامل نرسیده‌اند کپی می‌کند.
- وقتی چند گروه MCP را به اشتراک می‌گذارند، فرمان‌هایی که `ccam doctor`، `ccam ext` و `ccam mcp list` پیشنهاد می‌دهند به شکل `ccam migrate --mcp --group <a|b>` هستند.
- نبودن مبدأ با `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` پایان می‌یابد.

## مطالب مرتبط

- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [بررسی و تعمیر](/features/doctor)
- [چند حساب Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
