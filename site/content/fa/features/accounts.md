---
title: مدیریت حساب‌ها
description: افزودن، ورود، فهرست، تغییر نام، حذف و وارد کردن حساب‌های Claude Code با ccam add، ccam login، ccam list، ccam rename، ccam remove و ccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

هر حساب ccam یک پوشه پیکربندی جداگانه Claude Code با ورود مخصوص خودش است. این فرمان‌ها حساب‌ها را می‌سازند، وارد می‌کنند، فهرست می‌کنند، تغییر نام می‌دهند و حذف می‌کنند، و پوشه‌هایی را که از قبل دارید ثبت می‌کنند.

## نحوه استفاده

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

هر جا فرمانی `<name|number>` بپذیرد، می‌توانید نام حساب یا شماره جایگاه آن، یعنی ستون `#` در `ccam list`، را بدهید.

`ccam add` در ترمینال و بدون گزینه یک دستیار باز می‌کند. ابتدا نام و روش ورود حساب را می‌پرسد (ورود به Claude همین حالا، کلید API، یا ورود به Claude بعداً)، سپس گروه‌هایی را که می‌توان به آن‌ها پیوست نشان می‌دهد و `default` را علامت می‌زند (تا وقتی گروهی وجود ندارد، مواردی را که گروه `default` به اشتراک خواهد گذاشت)، بعد می‌پرسد `Create a shortcut command for work? [y/N]` و در پایان خلاصه‌ای برای تأیید نشان می‌دهد. `Esc` یک مرحله به عقب برمی‌گردد. هر یک از `--group`، `--no-group`، `--shortcut` یا `--no-login`، یا یک pipe به‌جای ترمینال، دستیار را کنار می‌گذارد؛ در این حالت نام الزامی است و میان‌بری ساخته نمی‌شود مگر اینکه `--shortcut` نامی برای آن بدهد.

بدون `--group`، حساب جدید به گروه `default` می‌پیوندد. نخستین حساب این گروه را می‌سازد که سرورهای MCP، تنظیمات، `CLAUDE.md` و افزونه‌ها را به اشتراک می‌گذارد؛ [اشتراک‌گذاری با گروه‌ها](/features/share) را ببینید.

`ccam list` جایگاه، نام، میان‌بر، ایمیل، گروه‌ها و آخرین استفاده را نشان می‌دهد و پیش از حساب پیش‌فرض `*` می‌گذارد:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

با `--json`، هر حساب یک شیء با `slot`، `shortcut`، `name`، `email`، `default`، `dir`، `last_used` (که اگر حساب هرگز اجرا نشده باشد حذف می‌شود)، `groups` و `share` است. `share` هر مورد (`mcp`، `settings`، `memory`، `extensions`، `history`) را به گروهی که از آن می‌آید، یا به `"local"`، نگاشت می‌کند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `ccam add --no-login` | حساب را بدون ورود می‌سازد؛ بعداً با `ccam login <name>` وارد شوید | خاموش |
| `ccam add --group <group>` | به‌جای `default` به این گروه می‌پیوندد. قابل تکرار؛ اگر `default` وجود نداشته باشد ساخته می‌شود | `default` |
| `ccam add --no-group` | به هیچ گروهی نمی‌پیوندد: حساب چیزی را به اشتراک نمی‌گذارد | خاموش |
| `ccam add --shortcut <command>` | این فرمان میان‌بر را هم در `~/.local/bin` می‌سازد | هیچ |
| `ccam add --api --key-from <source>` | حسابی که به جای ورود از کلید API استفاده می‌کند؛ [حساب‌های کلید API](/features/api) را ببینید | خاموش |
| `ccam list --json` | حساب‌ها را به‌صورت JSON چاپ می‌کند | خاموش |
| `ccam remove -y`، `--yes` | بدون پرسش حذف می‌کند | خاموش |
| `ccam import --name <name>` | نام حساب برای پوشه‌ای که می‌دهید | برگرفته از نام پوشه |
| `ccam import -y`، `--yes` | همه پوشه‌های پیداشده را بدون پرسش وارد می‌کند | خاموش |
| `--dry-run` | `ccam list` و `ccam completion` که چیزی را تغییر نمی‌دهند آن را می‌پذیرند. `add`، `login`، `logout`، `use`، `rename`، `remove` و `import` آن را با پیام `` --dry-run is not supported by `ccam add`; nothing was changed `` (با نام همان فرمان) رد می‌کنند | خاموش |

`ccam login`، `ccam logout`، `ccam use` و `ccam rename` گزینه‌ای از خود ندارند. `ccam completion` مقدار `bash`، `zsh` یا `powershell` را می‌پذیرد؛ هر چیز دیگری با `only bash, zsh and powershell are supported` ناموفق می‌شود.

نام حساب‌ها ۱ تا ۳۲ نویسه از `a-z`، `0-9`، `-` و `_` است و با یک حرف یا رقم شروع می‌شود. نام نمی‌تواند فقط از رقم تشکیل شود، چون `ccam 2` یعنی جایگاه ۲. نام همچنین نمی‌تواند یکی از فرمان‌های ccam باشد (`add`، `list`، `ls`، `run`، `use`، `next`، `usage`، `u`، `share`، `mcp`، `login`، `logout`، `rename`، `remove`، `rm`، `import`، `migrate`، `doctor`، `shims`، `shortcut`، `completion`، `help`، `version`، `statusline`، `update`، `folder`، `pick`، `ext`، `group`، `groups`)، تا `ccam <name>` هرگز مبهم نباشد. این فهرست `config.Reserved` در `internal/config` است.

## نحوه کار

`ccam add work --shortcut ccwork` این مراحل را در `addProfile` (`internal/cli/profiles.go`) اجرا می‌کند:

<!-- motion:accounts-add -->

1. `config.ValidateName` نام را و `shim.Check` میان‌بر را پیش از ساختن هر چیزی بررسی می‌کنند.
2. `profile.Add` پوشه `~/.local/share/ccam/profiles/work/` را با مجوز `0700` می‌سازد و کوچک‌ترین شماره جایگاه آزاد (`Config.NextSlot`) را به حساب می‌دهد، بنابراین شماره یک حساب حذف‌شده دوباره استفاده می‌شود. نخستین حساب، حساب پیش‌فرض می‌شود.
3. حساب به گروه‌هایش می‌پیوندد (`default`، مگر اینکه `--group` یا `--no-group` چیز دیگری بگوید؛ نخستین حساب `default` را می‌سازد) و `applyShare` پوشه جدید و خالی را با `share.Build` و سپس `Apply` به آن‌ها پیوند می‌دهد، برای مثال `settings.json` را به `groups/default/settings.json`. اگر این کار ناموفق باشد، پوشه جدید حذف می‌شود، چون پوشه باقی‌مانده هر تلاش دوباره را با `folder already exists` متوقف می‌کرد.
4. حساب از طریق `fsutil.AtomicWrite` در `~/.config/ccam/config.toml` ذخیره می‌شود: یک فایل موقت، `fsync` و سپس `rename`، بنابراین یک خرابی هرگز پیکربندی نیمه‌کاره بر جا نمی‌گذارد.
5. `runLogin` فرمان `claude auth login` را با `CLAUDE_CONFIG_DIR` تنظیم‌شده روی پوشه حساب، و با `--email` همراه با آخرین ایمیل حساب، اجرا می‌کند تا صفحه ورود از پیش پر شده باشد. اگر ورود ناموفق باشد، حساب باقی می‌ماند و ccam خطا را چاپ می‌کند و می‌گوید بعداً `ccam login work` را اجرا کنید.
6. `profile.ReadIdentity` مقدار `oauthAccount` را از `profiles/work/.claude.json` می‌خواند: ایمیل، شناسه‌های حساب و سازمان، نام سازمان و طرح اشتراک (`claude_max`، `claude_team` …). این فایل ممکن است چند مگابایت باشد، بنابراین `TopLevelValue` به‌جای تجزیه کل فایل، فقط به دنبال همان یک کلید می‌گردد. ccam آن‌ها را در `[profile.login]` در `config.toml` نگه می‌دارد.
7. `shim.Install` پیوند `~/.local/bin/ccwork` را به فایل اجرایی ccam می‌سازد که هنگام شروع، `argv[0]` را در `config.toml` جست‌وجو می‌کند.

ورود ممکن است چند دقیقه طول بکشد و در این مدت ccam دیگری ممکن است پیکربندی را ذخیره کند. به همین دلیل فرمان‌هایی که منتظر می‌مانند از طریق `env.update` ذخیره می‌کنند که `config.toml` را دوباره می‌خواند، همان یک تغییر (در اینجا ایمیل) را اعمال می‌کند و ذخیره می‌کند. `ccam remove` هم به همین دلیل پس از پرسش خود دوباره می‌خواند.

`ccam login` مراحل ۵ و ۶ را برای یک حساب موجود اجرا می‌کند. Claude Code ایمیل را فقط به‌عنوان راهنمایی برای صفحه ورود می‌گیرد، پس ممکن است ورود به حساب Claude دیگری برسد جز آنچه حساب ccam پیش‌تر بود؛ ccam ورود جدید را با `[profile.login]` مقایسه می‌کند و این را می‌گوید (`ccam: this login is other@company.com, but work was dev@company.com before.`). ورود جدید را نگه می‌دارد و می‌گوید چگونه آن را برگردانید: `ccam logout work`، سپس `ccam login work`، و با ایمیل درست وارد شوید. وقتی حساب ccam دیگری به همان حساب Claude وارد شده باشد، ccam هشدار می‌دهد که آن دو سقف‌های مصرف مشترکی دارند. `ccam login --expired` یکی پس از دیگری به همه حساب‌هایی وارد می‌شود که هرگز وارد نشده‌اند یا ورودشان (توکن تمدیدشان) منقضی شده است.

`ccam logout` فرمان `claude auth logout` را در حساب اجرا می‌کند و ایمیل ذخیره‌شده را پاک می‌کند؛ سابقه ورود باقی می‌ماند، پس ورود بعدی با آن مقایسه می‌شود. `ccam use` فقط `default` را در `config.toml` بازنویسی می‌کند.

وقتی حساب‌ها به بیش از یک سازمان تعلق داشته باشند، `ccam list` یک ستون `Organization` اضافه می‌کند (طرح‌های شخصی `personal` نشان می‌دهند) و `ccam list --json` برای هر حساب `plan` و `org` دارد.

### ورود روی یک ماشین راه دور

از طریق SSH، بخش مرورگر `claude auth login` ممکن است کار نکند. Claude Code روی سرور منتظر ورود در یک سرور callback در `http://localhost:<port>/callback` می‌ماند، اما مرورگری که ورود را تمام می‌کند روی رایانه شما اجرا می‌شود، جایی که آن نشانی به هیچ جا نمی‌رسد و مرورگر خطایی مانند `ERR_CONNECTION_REFUSED` نشان می‌دهد. نشانی همچنان کد یک‌بارمصرف را در خود دارد، پس ccam می‌تواند آن را برساند:

<!-- motion:accounts-relay -->

`runLogin` ترمینال شما را به Claude Code نمی‌سپارد. ورودی استاندارد Claude Code را به یک pipe وصل می‌کند و `relayLogin` هر خطی را که تایپ می‌کنید در آن کپی می‌کند، بنابراین کدی که در اعلان `Paste code here` در Claude Code بچسبانید همچنان می‌رسد. وقتی یک خط نشانی callback باشد، `callbackURL` فقط در صورتی آن را می‌پذیرد که میزبان `localhost`، `127.0.0.1` یا `[::1]` باشد، درگاه داده شده باشد، مسیر `/callback` باشد و هر دو `code` و `state` وجود داشته باشند؛ `http://` را می‌توان ننوشت. سپس `deliverCallback` آن نشانی را روی سرور با محدودیت ۱۰ ثانیه و بدون دنبال کردن تغییر مسیرها درخواست می‌کند، همان‌طور که مرورگر می‌کرد. درخواست به سرور callback فرمان `claude auth login` منتظر می‌رسد (خط‌چین در نمودار). Claude Code کد را دریافت می‌کند و ورود را ذخیره می‌کند؛ روی Linux در `.credentials.json` داخل پوشه حساب. یک خط چسبانده‌شده هرگز نمی‌تواند ccam را به تماس با میزبان دیگری وادار کند.

ccam پیش از شروع ورود این را چاپ می‌کند:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

در صورت موفقیت `→ Passed the sign-in to Claude Code on this machine.` را چاپ می‌کند.

وقتی ورود در یک ترمینال منتظر ورودی است، ccam پژواک (echo) و حالت خطی ترمینال را خاموش می‌کند و برای هر نویسه‌ای که تایپ می‌کنید یک `*` نشان می‌دهد، بنابراین کد یا نشانی callback چسبانده‌شده روی صفحه دیده نمی‌شود. `readLine` ویرایش خط را خودش انجام می‌دهد: `Enter` خط را می‌فرستد، `Backspace` یک نویسه را پاک می‌کند و `Ctrl+D` روی خط خالی ورودی را تمام می‌کند. ترمینال وقتی ورود تمام شود به حالت قبل برمی‌گردد، و همچنین با `Ctrl+C`، `SIGTERM` یا `SIGHUP` که همچنان ccam را پایان می‌دهند. اگر ورودی استاندارد و خطای استاندارد به ترمینال وصل نباشند، خط‌ها همان‌طور که می‌رسند کپی می‌شوند.

### تغییر نام، حذف و وارد کردن

`ccam rename` و `ccam remove` روی پوشه حساب و روی آنچه ccam در جاهای دیگر به نام آن حساب نگه می‌دارد عمل می‌کنند. `nameState` این مسیرها را فهرست می‌کند: `~/.local/state/ccam/last-used/<name>`، `~/.local/state/ccam/mcp/<name>` و حافظه‌های نهان مصرف `~/.local/state/ccam/usage/<name>.json` و `<name>.statusline.json`. `ccam remove work` آن‌ها را به این ترتیب انجام می‌دهد:

<!-- motion:accounts-remove -->

`ccam remove` ابتدا اگر حساب در حال اجرا باشد آن را رد می‌کند و سپس می‌پرسد. اگر حساب افزونه‌ها را به اشتراک بگذارد، `RewritePluginPaths` فایل‌های `installed_plugins.json` و `known_marketplaces.json` گروه را بازنویسی می‌کند، چون Claude Code مسیرهای مطلق افزونه‌ها را ذخیره می‌کند که ممکن است از پوشه این حساب بگذرند. `profile.Trash` پوشه را به `~/.local/share/ccam/trash/<name>-<time>/` منتقل می‌کند، حساب از `config.toml` حذف می‌شود، مسیرهای `nameState` پاک می‌شوند تا حساب بعدی با همین نام آن‌ها را به ارث نبرد، و پیوند میان‌بر حذف می‌شود. برای بازیابی، پوشه را از سطل زباله برگردانید.

`ccam rename` نام پوشه را تغییر می‌دهد، مسیرهای افزونه در `plugins/` آن را از پوشه قدیمی به پوشه جدید بازنویسی می‌کند، وضعیت پیش‌فرض را حفظ می‌کند و مسیرهای `nameState` را به نام جدید منتقل می‌کند. فرمان میان‌بر نام خود را نگه می‌دارد. اگر بازنویسی افزونه‌ها ناموفق باشد، نام پوشه به حالت قبل برمی‌گردد، چون پیکربندی ذخیره نمی‌شود.

`ccam import` یک پوشه را در همان جایی که هست، به‌صورت `dir` در `config.toml` ثبت می‌کند. هرگز آن را جابه‌جا یا کپی نمی‌کند: یک توکن تمدید کپی‌شده باعث می‌شد نسخه‌ها یکدیگر را از حساب خارج کنند. بدون آرگومان پوشه، به دنبال پوشه‌های ثبت‌نشده `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام می‌پرسد؛ نام‌ها از پوشه گرفته می‌شوند (`~/.claude` می‌شود `main`، `~/.claude-work` می‌شود `work`، `~/.claude-2` می‌شود `acct2`) و پوشه‌ای که نامی از آن به دست نیاید همراه با فرمان `--name` لازم فهرست می‌شود. بدون `CLAUDE_CONFIG_DIR`، Claude Code فایل `.claude.json` را در پوشه خانه شما کنار `~/.claude` نگه می‌دارد، بنابراین وقتی `~/.claude` را وارد می‌کنید و `.claude.json` خودش را ندارد، ccam پیوند `~/.claude/.claude.json` را به `~/.claude.json` می‌سازد. در این صورت `claude` معمولی و `ccam main` یک حساب را می‌بینند. حساب‌های واردشده تا وقتی آن‌ها را به گروهی اضافه نکنید (`ccam group join <group> <name>`) در هیچ گروهی نیستند و میان‌بری هم برایشان ساخته نمی‌شود.

`ccam completion bash|zsh` یک اسکریپت تکمیل خودکار چاپ می‌کند که نام حساب‌ها را هم تکمیل می‌کند.

## ایمنی و موارد خاص

- `ccam rename` و `ccam remove` تا وقتی حساب در حال اجراست (`refuseRunning`) با `work is running; quit its Claude Code session first` رد می‌کنند، چون نشست همچنان در پوشه می‌نویسد.
- `ccam remove` روی یک حساب واردشده فقط ثبت آن را لغو می‌کند: `✓ main unregistered (folder ... left in place)`. `~/.claude` شما به سطل زباله منتقل نمی‌شود.
- `ccam add` نامی را که پوشه‌اش از قبل در `profiles/` وجود دارد رد می‌کند و شما را به `ccam import` ارجاع می‌دهد.
- `ccam import` پوشه‌ای را که از قبل ثبت شده، مسیری را که پوشه نیست و نامی را که گرفته شده رد می‌کند.
- `ccam logout` و `ccam login` بررسی نمی‌کنند که حساب در حال اجرا هست یا نه.
- اگر واسطه نتواند به Claude Code برسد، برای مثال پس از پایان ورود، `ccam: could not pass the sign-in to Claude Code:` را همراه با دلیل چاپ می‌کند. خطای HTTP معمولاً یعنی کد منقضی شده است؛ ورود را دوباره شروع کنید.

## مطالب مرتبط

- [اجرای Claude Code](/features/launch)
- [فرمان‌های میان‌بر](/features/shortcuts)
- [اشتراک‌گذاری سرورهای MCP، تنظیمات و تاریخچه](/docs/sharing)
- [چند حساب Claude Code](/guides/multiple-claude-code-accounts)
