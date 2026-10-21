---
title: "کمانڈز کا حوالہ"
metaTitle: "ccam کمانڈز کا حوالہ"
description: "ccam کی ہر کمانڈ اور آپشن، ccam add اور ccam <name> سے لے کر ccam usage، ccam next اور ccam doctor تک۔"
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

کسی کمانڈ کو اصل میں چلائے بغیر یہ دیکھنے کے لیے کہ وہ کیا کرے گی، اس سے پہلے `--dry-run` لگائیں، مثلاً `ccam --dry-run work`۔ یہ Claude Code شروع کرتے وقت (`ccam <name>`، `ccam run`، `ccam next`، `ccam attach`) اور `group`، `migrate`، `shortcut`، `update` اور `doctor --fix` کے ساتھ کام کرتا ہے۔ کچھ تبدیل کرنے والی دوسری کمانڈز چلنے کے بجائے اسے مسترد کر دیتی ہیں۔

<!-- commands -->

## Claude Code شروع کرنا

`ccam <name|number> [arguments]` اس اکاؤنٹ کے طور پر Claude Code شروع کرتا ہے۔ ccam نام کے بعد آنے والے آرگیومنٹس نہیں پڑھتا؛ وہ انہیں اسی ترتیب میں `claude` کو دے دیتا ہے، اور اسٹینڈرڈ اِن پٹ، ایگزٹ کوڈ اور سگنلز بالکل ویسے ہی کام کرتے ہیں جیسے خود `claude` کے ساتھ۔ نام کے بعد `--` لگانا اختیاری ہے۔

اکاؤنٹ کے نام کے بغیر لیکن آرگیومنٹس کے ساتھ `ccam`، جیسے `ccam -p "summarise"`، موجودہ فولڈر کے لیے مقرر اکاؤنٹ، ورنہ ڈیفالٹ اکاؤنٹ استعمال کرتا ہے۔ انہیں `ccam folder <name> [folder]` اور `ccam use <name>` سے مقرر کریں۔

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

اکاؤنٹ کے نام سے پہلے `--fallback` اُن رنز کے لیے ہے جو `-p`/`--print` استعمال کرتے ہیں۔ جب اکاؤنٹ اپنے استعمال کی حد کو پہنچ جائے تو ccam وہی کمانڈ سب سے زیادہ گنجائش والے لاگ اِن اکاؤنٹ سے دوبارہ چلاتا ہے اور صرف اس رن کا آؤٹ پٹ پرنٹ کرتا ہے جو مکمل ہوا۔ دیکھیں [فولڈر کے حساب سے اکاؤنٹس](/features/folder) اور [حد پر دوبارہ چلانا](/features/next#run-again-at-the-limit)۔

`ccwork` جیسا شارٹ کٹ اپنے اکاؤنٹ کو اسی طرح شروع کرتا ہے: `ccwork --resume` وہی کرتا ہے جو `ccam work --resume` کرتا ہے۔ [شارٹ کٹس](/docs/commands#shortcuts) دیکھیں۔

## اکاؤنٹس شامل کرنا اور ہٹانا

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

فولڈر بتائے بغیر `ccam import`، `~/.claude` اور `~/.claude-*` تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے۔ امپورٹ کیے گئے فولڈرز اپنی جگہ پر ہی رجسٹر ہوتے ہیں؛ ccam کبھی کوئی لاگ اِن منتقل یا کاپی نہیں کرتا۔

`ccam add --api` ایسا اکاؤنٹ شامل کرتا ہے جو Claude لاگ اِن کے بجائے API key استعمال کرتا ہے؛ ccam ہر لانچ پر key پڑھتا ہے اور اسے کبھی محفوظ نہیں کرتا۔ دیکھیں [API key اکاؤنٹس](/features/api)۔

جب Claude Code کسی اکاؤنٹ کے طور پر چل رہا ہو تو `ccam rename`، `ccam remove` اور `ccam group` کمانڈز اس اکاؤنٹ میں تبدیلی سے انکار کر دیتی ہیں۔ پہلے وہ سیشن بند کریں۔

`ccam list --json` ہر اکاؤنٹ کو اس کے سلاٹ، شارٹ کٹ، ای میل اور فولڈر کے ساتھ پرنٹ کرتا ہے۔ اس کی `groups` array اکاؤنٹ کے گروپس کی فہرست دیتی ہے، اور اس کا `share` آبجیکٹ ہر آئٹم (`mcp`، `settings`، `memory`، `extensions`، `history`) کو اس گروپ سے جوڑتا ہے جس سے وہ آتا ہے، یا `local` سے۔

## شارٹ کٹس

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

شارٹ کٹ ایک کمانڈ ہے جس کا نام آپ خود رکھتے ہیں اور جو ایک اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے۔ `ccam work` کی طرح یہ ہر آرگیومنٹ بغیر تبدیلی کے Claude Code کو دیتی ہے، اس لیے `ccwork --resume` اور `ccwork -p "hi"` ویسے ہی کام کرتے ہیں جیسے آپ توقع کرتے ہیں۔ یہ `~/.local/bin` میں ccam کا ایک لنک ہے۔

شارٹ کٹس اختیاری ہیں۔ `ccam add` کا وزرڈ پوچھتا ہے کہ شارٹ کٹ بنایا جائے یا نہیں، `ccam add <name> --shortcut <command>` بغیر پوچھے ایک بنا دیتا ہے، اور `ccam list` ہر اکاؤنٹ کا شارٹ کٹ دکھاتا ہے، یا شارٹ کٹ نہ ہو تو `-`۔

ناموں میں a-z، 0-9، `.`، `-` اور `_` حروف استعمال ہوتے ہیں۔ نام کوئی ccam کمانڈ، کسی دوسرے اکاؤنٹ کا شارٹ کٹ، یا آپ کے `PATH` پر موجود کسی دوسرے پروگرام کا نام نہیں ہو سکتا۔

`ccam rename` شارٹ کٹ کا نام برقرار رکھتا ہے، اور `ccam remove` اکاؤنٹ کے ساتھ شارٹ کٹ بھی حذف کر دیتا ہے۔ `ccam shims install` تمام اکاؤنٹس کے لنکس دوبارہ بناتا ہے، مثلاً جب آپ نے ccam کو کسی اور جگہ منتقل کیا ہو۔

## شیئرنگ

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ایک گروپ اپنے اکاؤنٹس کے درمیان آئٹمز شیئر کرتا ہے؛ ایک اکاؤنٹ کئی گروپس میں ہو سکتا ہے، لیکن ہر آئٹم صرف ایک گروپ سے آتا ہے۔ [گروپس کے ساتھ شیئرنگ](/features/share) دیکھیں۔

## موجودہ سیٹ اپ کاپی کرنا

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` کسی Claude Code فولڈر سے `settings.json`، `CLAUDE.md`، ایکسٹینشنز اور یوزر MCP سرورز کو ایک گروپ میں کاپی کرتا ہے، تاکہ وہ آئٹمز شیئر کرنے والے ہر اکاؤنٹ کو یہ مل جائیں۔ فولڈر کو صرف پڑھا جاتا ہے۔ آپشنز `--settings`، `--memory`، `--extensions`، `--mcp` اور `--history` صرف وہی آئٹمز کاپی کرتے ہیں جن کا آپ نام لیں۔ ان کے بغیر، ٹرمینل پر ایک چیک لسٹ سے آپ خود چن سکتے ہیں، جس میں گفتگو کی ہسٹری کے سوا سب کچھ چیک ہوتا ہے؛ باقی صورتوں میں گفتگو کی ہسٹری کے سوا سب کچھ کاپی ہوتا ہے۔ `-g` یا `--group` گروپ چنتا ہے، اور `-y` یا `--yes` بغیر پوچھے کاپی کرتا ہے اور فائل مختلف ہونے پر گروپ والی کاپی رکھتا ہے۔ یہ سرورز کے محفوظ کیے گئے سائن اِن بھی ہر اکاؤنٹ کی اپنی لاگ اِن فائل میں کاپی کرتا ہے۔ [موجودہ سیٹ اپ کاپی کرنا](/docs/sharing#copy-an-existing-setup) دیکھیں۔

## MCP سرورز

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` وہ سرور جو صرف اسی اکاؤنٹ میں ہیں اس کے گروپ میں منتقل کرتا ہے۔ جس سرور کا نام گروپ پہلے سے کسی مختلف تعریف کے لیے استعمال کر رہا ہو اسے چھوڑ دیا جاتا ہے؛ گروپ کی کاپی بدلنے کے لیے `--force` شامل کریں، ccam پہلے اس کا بیک اپ لیتا ہے۔

## استعمال اور اکاؤنٹ تبدیل کرنا

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

`next --auto`، `--fallback` اور `ccam pick` ایک ہی طرح چنتے ہیں: پہلے `auto_max_usage` (`config.toml`، ڈیفالٹ 100) سے کم والے اکاؤنٹس، `auto_strategy` (`headroom`، `expiring` یا `balance`) کے مطابق، `--model` کی ہفتہ وار ونڈو کو شمار کرتے ہوئے، کبھی 100% والا اکاؤنٹ نہیں، اور اسی تنظیم کے اندر جس کے اکاؤنٹ سے وہ شروع کرتے ہیں (آگے دیکھنے کے لیے `--cross-org` یا `cross_org = true`)۔ جب کوئی اکاؤنٹ نہ چنا جا سکے تو یہ کوڈ 75 کے ساتھ باہر نکلتے ہیں، یا `--wait` کے ساتھ پہلے اکاؤنٹ کے ری سیٹ ہونے کا انتظار کرتے ہیں۔ `ccam next --continue` گفتگو کو ان اکاؤنٹس کے درمیان بھی لے جاتا ہے جو ہسٹری شیئر نہیں کرتے: تب Claude Code اسے نئے اکاؤنٹ میں کاپی کر دیتا ہے۔ `ccam <name>` کے برعکس، `ccam next` اپنے آپشنز خود پڑھتا ہے، اس لیے Claude Code کے آپشنز `--` کے بعد آتے ہیں؛ اس کے بغیر ccam `unknown flag` کے ساتھ رک جاتا ہے۔ Claude Code کے اندر، پلگ اِن کا `/ccam:switch [account]` گفتگو آگے سونپتا ہے، اور `ccam hook` وہ ہے جو اس کے ہُکس چلاتے ہیں۔ دیکھیں [اگلے اکاؤنٹ پر جانا](/features/next)، [اسکرپٹس کے لیے اکاؤنٹ چننا](/features/pick) اور [ہر اکاؤنٹ کے پس منظر کے سیشنز](/features/agents)۔

## دیکھ بھال

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

`ccam update`، ccam کو بدلنے سے پہلے ریلیز کا SHA-256 چیک سم جانچتا ہے، اور `--version v0.2.0` کوئی مخصوص ریلیز انسٹال کرتا ہے۔ `.deb` پیکیج سے انسٹال کیا گیا ccam اس کے بجائے apt سے اپ ڈیٹ ہوتا ہے۔ [اپ ڈیٹ کرنا](/docs/install#update) دیکھیں۔

## فائلیں

| راستہ | مواد |
| --- | --- |
| `~/.config/ccam/config.toml` | اکاؤنٹس، گروپس اور سیٹنگز (`theme`، `usage_source`، `relay`، `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | ہر اکاؤنٹ کا `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | شیئر کی گئی سیٹنگز، MCP سرورز، ایکسٹینشنز اور ہسٹری |
| `~/.local/state/ccam/` | استعمال کا کیش، آخری استعمال کے اوقات اور وہ نوٹس جو `/ccam:switch` اور حد والا ہُک چھوڑتے ہیں (`handoff/`) |
| `~/.local/bin/<shortcut>` | `ccwork` جیسی شارٹ کٹ کمانڈز، جو ccam کے لنکس ہیں |

`XDG_CONFIG_HOME`، `XDG_DATA_HOME` اور `XDG_STATE_HOME` ویری ایبلز ان فولڈرز کی جگہ بدل دیتے ہیں۔ Windows پر `~` سے مراد `%USERPROFILE%` ہے اور شارٹ کٹس `<shortcut>.exe` ہوتے ہیں۔
