---
title: اشتراک‌گذاری با گروه‌ها
description: یک گروه ccam موارد انتخاب‌شده را از طریق یک پوشه بین حساب‌هایش مشترک می‌کند؛ ccam group داده‌ها را با برنامه، نسخه پشتیبان و انتخاب در تعارض‌ها جابه‌جا می‌کند.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 3f94dd13dbf1
---

هر حساب پوشه خود را دارد، اما بیشتر چیزهایی که در آن است لازم نیست جدا باشد. یک **گروه** تعیین می‌کند چه چیزی و توسط چه کسانی مشترک شود: هر حساب عضو گروه آن موارد را به یک پوشه، `~/.local/share/ccam/groups/<group>/`، پیوند می‌دهد، بنابراین همه آن‌ها یک نسخه واحد را می‌خوانند و در آن می‌نویسند. توکن‌های ورود هرگز مشترک نمی‌شوند.

یک حساب می‌تواند در چند گروه باشد، اما هر مورد فقط از یک گروه می‌آید. یک پیکربندی معمول، یک گروه برای تنظیمات و افزونه‌ها و گروه دومی برای تاریخچه گفتگو بین دو حساب است:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

نخستین حسابی که اضافه می‌کنید `default` را می‌سازد و حساب‌های بعدی به آن می‌پیوندند، مگر اینکه چیز دیگری بگویید.

## کاربرد

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

هر فرمانی که داده را جابه‌جا می‌کند یک برنامه چاپ می‌کند، درباره هر فایلی که با نسخه گروه متفاوت است می‌پرسد و پیش از تغییر هر چیزی `Proceed? [y/N]` را می‌پرسد:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

همین تغییرها را می‌توانید در زبانه **Groups** در [داشبورد](/features/dashboard) هم انجام دهید؛ برای انتخاب گروه‌های یک حساب، روی آن `s` را بزنید.

## موارد

| مورد | آنچه پوشش می‌دهد |
| --- | --- |
| `mcp` | سرورهای MCP، از طریق `mcp.json` گروه ([سرورهای MCP مشترک](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | فایل سراسری `CLAUDE.md` |
| `extensions` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/`، `rules/` |
| `history` | `projects/` و `history.jsonl`، تا حساب دیگری بتواند `--continue` کند |

اگر هیچ‌یک از گروه‌های یک حساب موردی را به اشتراک نگذارد، آن مورد در پوشه خود حساب می‌ماند.

## فرمان‌ها و گزینه‌ها

| فرمان | کاری که انجام می‌دهد |
| --- | --- |
| `ccam group` یا `ccam group list` | فهرست گروه‌ها |
| `ccam group show <group>` | آنچه گروه به اشتراک می‌گذارد، حساب‌هایش، پوشه‌اش و سرورهای MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | یک گروه می‌سازد؛ موارد حساب‌هایش به آن منتقل می‌شوند |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | آنچه گروه به اشتراک می‌گذارد را برای همه حساب‌های آن تغییر می‌دهد |
| `ccam group join <group> <account>...` | حساب‌ها را در گروه قرار می‌دهد |
| `ccam group leave <group> <account>...` | حساب‌ها را بیرون می‌برد؛ هر کدام یک نسخه نگه می‌دارد، از جمله سرورهای MCP |
| `ccam group rename <group> <new-name>` | فقط در پیکربندی نام را تغییر می‌دهد |
| `ccam group remove <group> [--leave]` | یک گروه را حذف می‌کند |

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--from <account>` | `add`، `set`: حسابی که نسخه‌اش ابتدا وارد گروه می‌شود | حساب پیش‌فرض، سپس به ترتیب شماره |
| `--replace` | `join`: ابتدا از گروه‌هایی که همان موارد را به اشتراک می‌گذارند بیرون می‌رود | خاموش |
| `--fresh` | حساب‌هایی که اشتراک یک مورد را متوقف می‌کنند، به‌جای یک نسخه، بدون آن شروع می‌کنند (یک پوشه واردشده باز هم یک نسخه نگه می‌دارد) | خاموش |
| `--leave` | `remove`: ابتدا حساب‌های گروه را بیرون می‌برد | وقتی گروه حساب دارد لازم است |
| `--prefer <group\|mine\|merge>` | نحوه تصمیم‌گیری در تعارض‌ها | پرسیدن؛ با `--yes`، پیش‌فرض هر فایل |
| `-y`، `--yes` | بدون پرسش اجرا می‌کند؛ تعارض‌ها از `--prefer` پیروی می‌کنند | خاموش |
| `--dry-run` (سراسری) | برنامه را چاپ می‌کند و چیزی را تغییر نمی‌دهد | خاموش |

حساب‌ها با نام یا شماره داده می‌شوند و موارد به‌صورت فهرستی جداشده با ویرگول. `ccam add` گزینه `--group <group>` (قابل تکرار) یا `--no-group` را می‌پذیرد.

## چگونه کار می‌کند

فایل‌های درگیر:

- `~/.config/ccam/config.toml`: ورودی‌های `[[group]]`: `name`، `items`، `members` و، وقتی نام پوشه متفاوت است، `dir`.
- `~/.local/share/ccam/profiles/<name>/`: پوشه `CLAUDE_CONFIG_DIR` حساب.
- `~/.local/share/ccam/groups/<folder>/`: نسخه‌های مشترک.

وقتی `config.toml` خوانده می‌شود، گروه‌های هر حساب به یک پوشه برای هر مورد تبدیل می‌شوند. این کار برای پنج حساب حدود 25 µs طول می‌کشد، بنابراین یک اجرا همچنان فقط `config.toml` را می‌خواند. پیکربندی‌ای که در آن یک حساب یک مورد را از دو گروه بگیرد خطاست و تا وقتی اصلاح نشود هیچ چیز اجرا نمی‌شود.

یک تغییر روی یک کپی از پیکربندی کار می‌کند:

1. فرمان کپی را ویرایش می‌کند (`AddGroup`، `SetGroupItems`، `JoinGroup` …) که قواعد را بررسی می‌کند.
2. ccam پیکربندی قدیم و جدید را مقایسه می‌کند و برای هر حسابی که مواردش تغییر می‌کند، یک برنامه می‌سازد (`internal/group`، `share.Build`). برنامه‌ها یکی پس از دیگری ساخته می‌شوند و هر کدام می‌بیند برنامه‌های قبلی چه چیزی را در پوشه‌های گروه خواهند گذاشت، بنابراین فقط نخستین حساب نسخه خود را وارد می‌کند. نخستین نسخه گروه به این ترتیب از این‌جاها می‌آید: پوشه گروه (اگر از قبل یکی داشته باشد)، حساب `--from`، حساب پیش‌فرض و سپس حساب‌های دیگر به ترتیب شماره.
3. اگر یکی از آن حساب‌ها در حال اجرا باشد، Claude Code به هر شکلی که شروع شده باشد، ccam کار را رد می‌کند (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``)، و همچنین در Windows بدون Developer Mode اگر فایلی به پیوند نیاز داشته باشد.
4. ccam تعارض‌ها را حل می‌کند (با پرسیدن از شما یا با `--prefer`)، برنامه را چاپ می‌کند و برای تأیید از شما می‌پرسد.
5. درست پیش از ذخیره، ccam دوباره بررسی می‌کند. اگر حسابی در مدتی که ccam منتظر پاسخ‌های شما بود شروع شده باشد، تغییر متوقف می‌شود؛ همچنین اگر فرایند ccam دیگری در این فاصله همان گروه‌ها را تغییر داده باشد، تغییر متوقف می‌شود (`the groups changed while this plan waited; nothing was moved, make the change again`). پیکربندی جدید پیش از جابه‌جایی هر چیزی ذخیره می‌شود و تغییر در `~/.local/state/ccam/group-change.json` یادداشت می‌شود. اگر گامی در میانه کار متوقف شود، یادداشت می‌ماند: `ccam doctor --fix` تغییر را به ترتیب خودش، با پاسخ‌هایی که داده‌اید و از جمله مرحله‌های MCP، تمام می‌کند و تا آن زمان هیچ تغییر دیگری در گروه‌ها شروع نمی‌شود (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. مرحله‌های MCP اول اجرا می‌شوند. وقتی حساب‌ها افزونه‌ها را از یک گروه بیرون می‌برند، فهرست‌های پلاگین گروه هر جا که مسیرهایشان از پوشه یک حساب می‌گذرد به پوشه گروه هدایت می‌شوند (`group plugin lists:` در برنامه). سپس برنامه‌ها به ترتیب اجرا می‌شوند. هر گام داده را به گروه یا به پوشه پشتیبان منتقل می‌کند؛ فقط فایل‌هایی که با نسخه گروه یکسان‌اند، به‌جای رفتن به پشتیبان، حذف می‌شوند.

ccam فهرست مجازی از مسیرهایی که هر مورد پوشش می‌دهد نگه می‌دارد (`share.Layout`). هر چیزی که در آن نباشد هرگز به یک گروه نمی‌رسد. ccam همچنین `.credentials.json` و `.claude.json` را به‌عنوان مواردی که هرگز به اشتراک گذاشته نمی‌شوند نام می‌برد (`share.Never`) و اگر هر کدام از آن‌ها روزی در layout ظاهر شود، با یک خطای داخلی متوقف می‌شود (`Build`). مورد `mcp` هیچ مسیری ندارد: فهرست سرورها درون `.claude.json` است، بنابراین [اشتراک‌گذاری MCP](/features/mcp) از طریق `mcp.json` گروه کار می‌کند.

این اتفاقی است که برای یک پوشه، `projects/`، می‌افتد وقتی `side` اشتراک تاریخچه را آغاز می‌کند و گروه از قبل محتوا دارد:

<!-- motion:share-on -->

اگر گروه هنوز چیزی در آن مسیر نداشته باشد، حساب اولین حسابی است که آن را به اشتراک می‌گذارد و پوشه‌اش همان‌طور که هست به آنجا منتقل می‌شود. در غیر این صورت هر فایلی که گروه ندارد به گروه منتقل می‌شود و فایل‌هایی که در هر دو طرف وجود دارند در پوشه حساب می‌مانند، بنابراین نسخه گروه حفظ می‌شود. فایل‌هایی که با نسخه گروه یکسان‌اند حذف می‌شوند و برنامه آن‌هایی را که فرق دارند نام می‌برد (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

سپس پوشه حساب به `.ccam-backup/<time>/` درون پوشه حساب می‌رود و یک پیوند جای آن را می‌گیرد. اگر ورودی گروه خودش پیوندی به ورودی خودِ حسابی باشد که می‌پیوندد، محتوای حساب به جای آن به گروه می‌رود. نام فایل‌های گفتگو از شناسه نشست آن‌ها گرفته می‌شود، بنابراین تاریخچه چند حساب بدون تداخل با هم ترکیب می‌شود. رجیستری‌های `plugins/` مسیرهای مطلق را از طریق پوشه حساب ذخیره می‌کنند؛ ccam آن‌ها را به پوشه گروه هدایت می‌کند و ورودی‌های حساب را در رجیستری‌ای که گروه از قبل دارد ادغام می‌کند.

بیرون رفتن از گروه، پیوند را برعکس می‌کند:

<!-- motion:share-off -->

ccam ابتدا محتوای فعلی گروه را در `.ccam-backup/incomplete-<name>` درون پوشه حساب کپی می‌کند و فقط وقتی کپی کامل شد پیوند را جایگزین می‌کند. اگر کپی متوقف شود، پیوند سر جایش می‌ماند و `ccam doctor --fix` دوباره کپی می‌کند (برنامه چنین تکمیلی `restore <entry>: … (copied back from the group before a stop)` را نشان می‌دهد). گروه محتوای خود را برای حساب‌های دیگر نگه می‌دارد.

- ccam پیوند را فقط وقتی حذف می‌کند که هنوز، مانند زمان ساختن برنامه، به پوشه groups اشاره کند. مسیری که از آن پس تغییر کرده باشد، مانند یک فایل واقعی در جای پیوند، همان‌طور که هست می‌ماند و تغییر متوقف می‌شود.
- با `--fresh` کپی انجام نمی‌شود.
- حسابی که یک مورد را به گروهی منتقل می‌کند که پوشه‌اش هنوز خالی است، یک نسخه از آنچه می‌دید با خود می‌آورد، بنابراین فایل‌هایش را نگه می‌دارد. آن نسخه تا کامل شدن در `<group>/.ccam-incomplete-<name>` آماده می‌شود.
- حسابی که اشتراک MCP را متوقف می‌کند سرورهای گروه را نگه می‌دارد: آن‌ها به `.claude.json` خود آن حساب می‌روند (یک نسخه پشتیبان نگه داشته می‌شود)، مگر آنجا که حساب از قبل سروری با همان نام داشته باشد، که همان‌طور که هست می‌ماند. با `--fresh` به جای آن برنامه سرورهایی را که دیگر بارگذاری نمی‌شوند فهرست می‌کند. پیوستن دوباره کپی‌هایی را که کسی تغییرشان نداده حذف می‌کند.
- بیرون رفتن از `extensions` همچنین فهرست‌های پلاگین گروه را، هر جا که مسیرهایشان از پوشه حسابِ بیرون‌رونده می‌گذشت، به پوشه گروه هدایت می‌کند؛ برای حساب‌های دیگر، مسیرهای بازنویسی‌شده به همان پلاگین‌ها می‌رسند.

یک فایل واحد که در هر دو طرف با محتوای متفاوت وجود دارد، یک تعارض است:

<!-- motion:share-conflict -->

برای فایل‌های تکی مانند `settings.json` شما انتخاب می‌کنید؛ درون پوشه‌ها نسخه گروه می‌ماند و برنامه فایل‌هایی را که فرق دارند نام می‌برد. گزینه‌ها به فایل بستگی دارند؛ گزینه اول پیش‌فرض است:

| فایل | گزینه‌ها |
| --- | --- |
| `history.jsonl` | `merge by line`، `keep group`، `use mine` |
| دیگر فایل‌های `.json` | `keep group`، `use mine`، `merge JSON` |
| هر فایل دیگر | `keep group`، `use mine` |

`merge JSON` اشیا را به‌صورت بازگشتی ادغام می‌کند، آرایه‌ها را بدون تکرار به هم می‌پیوندد و وقتی دو طرف یک کلید را متفاوت تنظیم کرده باشند، مقدار گروه را نگه می‌دارد، بنابراین حساب‌های دیگر هیچ تغییری نمی‌بینند. هر انتخابی به جز `keep group` ابتدا نسخه گروه را در `.ccam-backup/<time>/group/` حساب ذخیره می‌کند. فایل خود شما همیشه به پوشه پشتیبان می‌رود.

در هر اجرا، ccam هر مسیر در layout را فقط با `lstat` و `readlink` با گروه‌های حساب مقایسه می‌کند (`share.Check`). پیوندی که وجود ندارد همان لحظه ساخته می‌شود. یک فایل واقعی در جایی که باید پیوند باشد، پیوندی به جای اشتباه، یا یک پیوند کهنه، اجرا را متوقف می‌کند و به `ccam doctor --fix` ارجاع می‌دهد.

### تغییر نام، حذف و قالب قدیمی

- `ccam group rename` فقط پیکربندی را تغییر می‌دهد. پوشه نام خود را نگه می‌دارد و گروه آن را به‌صورت `dir` ثبت می‌کند، بنابراین هیچ چیز جابه‌جا نمی‌شود و هیچ چیز نمی‌تواند از دست برود.
- `ccam group remove` گروهی را که حساب دارد رد می‌کند، مگر اینکه `--leave` را اضافه کنید. در این صورت هر حساب یک نسخه نگه می‌دارد و پوشه به `~/.local/share/ccam/trash/` می‌رود، مگر اینکه گروه دیگری از آن استفاده کند.
- پیش از گروه‌ها، `config.toml` برای هر حساب یک جدول `[profile.share]` ذخیره می‌کرد. ccam آن قالب را بدون جابه‌جا کردن چیزی به‌صورت گروه می‌خواند. در نخستین ذخیره، قالب جدید را می‌نویسد و فایل قدیمی را به‌صورت `config.toml.bak` نگه می‌دارد. مواردی که حساب‌های یکسانی به اشتراک می‌گذارند به یک گروه با نام پوشه تبدیل می‌شوند؛ بقیه نام‌هایی مانند `default-history` می‌گیرند و همچنان از همان پوشه استفاده می‌کنند.

### در Windows

Windows پیوندهای نمادین را فقط در Developer Mode (Settings > System > For developers) یا برای مدیر سیستم می‌سازد. بدون این حق، یک پوشه مشترک به یک directory junction تبدیل می‌شود که برای Claude Code همان‌طور کار می‌کند. یک فایل مشترک (`settings.json`، `CLAUDE.md`، `history.jsonl`) فقط می‌تواند پیوند نمادین باشد، بنابراین تغییر گروهی که به آن نیاز دارد، پیش از جابه‌جایی هر چیزی رد می‌شود. به همین دلیل، نخستین `ccam add` در آنجا گروه `default` را فقط با `mcp` و `extensions` می‌سازد و این را می‌گوید، و `ccam add --group` به گروهی که تنظیمات، حافظه یا تاریخچه را به اشتراک می‌گذارد رد می‌شود.

## ایمنی و موارد خاص

- **حساب‌های در حال اجرا.** تا وقتی Claude Code با یکی از حساب‌هایی که تغییر جابه‌جا می‌کند در حال اجراست، تغییر رد می‌شود، چون نشست همچنان در آن فایل‌ها می‌نویسد، چه ccam آن را شروع کرده باشد چه نه (`claude` معمولی، یک IDE، نشست‌های پس‌زمینه). این بررسی درست پیش از ذخیره تکرار می‌شود. پیوستن و بیرون رفتن فقط روی حسابی که نام برده شده اثر دارد؛ بیرون رفتن از `extensions` مسیرهای فهرست‌های پلاگین گروه را هم بازنویسی می‌کند، بی‌آنکه معنای آن‌ها را برای دیگران تغییر دهد. `--dry-run` همچنان برنامه را چاپ می‌کند.
- **نسخه‌های پشتیبان.** هر برنامه پوشه `.ccam-backup/<time>/` خود را در پوشه حساب می‌گیرد. ccam هرگز آن‌ها را حذف نمی‌کند.
- **پیوندهای خود شما.** اگر یک مسیر از قبل پیوندی باشد که خودتان ساخته‌اید، برای مثال به یک مخزن dotfiles، و گروه خالی باشد، ccam مقصد پیوند را در گروه کپی می‌کند و خود پیوند را به پوشه پشتیبان منتقل می‌کند. به مقصد دست زده نمی‌شود.
- **پوشه‌های واردشده.** حسابی که با `ccam import` در همان جای خود ثبت شده است (مانند `~/.claude`) در همان جا تغییر می‌کند، بنابراین `claude` معمولی هم تغییر را می‌بیند؛ برنامه این را می‌گوید. وقتی از یک گروه بیرون می‌رود همیشه یک نسخه نگه می‌دارد، حتی با `--fresh`، و `ccam remove` پیش از آنکه حساب را فراموش کند آن را به همین شکل از گروه‌هایش بیرون می‌برد. تا وقتی در یک گروه است، فایل‌های مشترکش در پوشه گروه قرار دارند، پس پیش از پاک کردن `~/.local/share/ccam` حساب را حذف کنید (یا از گروه‌ها بیرون بروید).
- **دو حساب در یک پوشه.** حساب‌هایی که تاریخچه را به اشتراک می‌گذارند در فایل‌های یکسانی می‌نویسند. وقتی حساب دومی در پوشه‌ای شروع به کار کند که یکی از قبل در آن در حال اجراست، ccam هشدار می‌دهد و می‌پرسد.
- **ویرایش دستی.** اگر `members` نامی را فهرست کند که حساب نیست، ccam آن را نادیده می‌گیرد و `ccam doctor` آن را گزارش می‌کند. نام‌های `local`، `on`، `off` و `none` نمی‌توانند نام گروه باشند.

## مرتبط

- [داشبورد](/features/dashboard)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [سرورهای MCP مشترک](/features/mcp)
- [کپی کردن یک پیکربندی موجود](/features/migrate)
- [مرجع فرمان‌ها](/docs/commands)
