---
title: "کمانڈز کا حوالہ"
metaTitle: "ccam کمانڈز کا حوالہ"
description: "ccam کی ہر کمانڈ اور آپشن، ccam add اور ccam <name> سے لے کر ccam usage، ccam next اور ccam doctor تک۔"
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "41fb781fd146"
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

**Team اور Enterprise سیٹیں۔** خودکار انتخاب (`ccam next`، `--fallback`، `ccam pick`، ڈیش بورڈ کی `n` کلید، اور جب آپ اسے آن کریں تو [`auto_switch`](/features/next#auto-switch)) ایک ہی Claude Team یا Enterprise تنظیم کی سیٹوں کے درمیان بھی اکاؤنٹ بدلتا ہے۔ Anthropic کا ہیلپ سینٹر کہتا ہے کہ "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (کسی Team پلان کی ہر سیٹ ایک نامزد فرد کو تفویض کی جاتی ہے، اور صرف اسی شخص کو اسے استعمال کرنے کی اجازت ہے) اور "Users can't route their work through another user's account to use that user's remaining usage" (صارفین کسی دوسرے صارف کا باقی ماندہ استعمال کام میں لانے کے لیے اپنا کام اس کے اکاؤنٹ کے ذریعے نہیں چلا سکتے) ([Team پلانز](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ اسی طرح [Enterprise پلانز](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) کے لیے بھی)۔ اس لیے سیٹوں کے درمیان خودکار انتخاب استعمال کرنے سے Anthropic کی شرائط کی خلاف ورزی ہو سکتی ہے۔ وہاں اسے استعمال کرنا ہے یا نہیں، یہ آپ کا فیصلہ ہے: اس پر انحصار کرنے سے پہلے اپنی تنظیم اور Anthropic سے تصدیق کر لیں۔ ذاتی پلانز (Pro، Max) کی بات [عام سوالات میں شرائط والے نوٹ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) میں کی گئی ہے۔

`ccwork` جیسا شارٹ کٹ اپنے اکاؤنٹ کو اسی طرح شروع کرتا ہے: `ccwork --resume` وہی کرتا ہے جو `ccam work --resume` کرتا ہے۔ [شارٹ کٹس](/docs/commands#shortcuts) دیکھیں۔

## اکاؤنٹس شامل کرنا اور ہٹانا

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

فولڈر بتائے بغیر `ccam import`، `~/.claude` اور `~/.claude-*` تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے؛ جو فولڈر ناکام ہو اس کی اطلاع `Skipped <path>: <reason>` کے طور پر دی جاتی ہے اور باقی رجسٹر ہو جاتے ہیں، اور جس فولڈر کا نام پہلے سے کسی اکاؤنٹ کا ہو اسے `--name` چاہیے، جو فولڈر کے ساتھ دیا جاتا ہے۔ امپورٹ کیے گئے فولڈرز اپنی جگہ پر ہی رجسٹر ہوتے ہیں؛ ccam کبھی کوئی لاگ اِن منتقل یا کاپی نہیں کرتا۔

`ccam login --expired` ان اکاؤنٹس میں بھی لاگ اِن کرتا ہے جن کا لاگ اِن Claude Code نے مٹا دیا ہو۔ `ccam remove` کسی منیجڈ اکاؤنٹ کا فولڈر ٹریش میں جانے سے پہلے اسے لاگ آؤٹ کرتا ہے (تب نہیں جب کسی دوسرے اکاؤنٹ کے پاس وہی لاگ اِن ہو)، اور کسی امپورٹ کیے گئے اکاؤنٹ کا `ccam remove` اسے اس کے گروپس سے نکال دیتا ہے، اور اس کے پاس اس چیز کی ایک کاپی رہتی ہے جو وہ گروپس شیئر کرتے ہیں۔ macOS پر `ccam rename` فولڈر کا نام برقرار رکھتا ہے، کیونکہ Keychain لاگ اِن کو فولڈر کے راستے کے تحت رکھتی ہے۔

`ccam add --api` ایسا اکاؤنٹ شامل کرتا ہے جو Claude لاگ اِن کے بجائے API key استعمال کرتا ہے؛ ccam ہر لانچ پر key پڑھتا ہے اور اسے کبھی محفوظ نہیں کرتا۔ دیکھیں [API key اکاؤنٹس](/features/api)۔

جب Claude Code کسی اکاؤنٹ کے طور پر چل رہا ہو، چاہے وہ کسی بھی طرح شروع ہوا ہو (بیک گراؤنڈ سیشنز بھی)، تو `ccam rename`، `ccam remove` اور `ccam group` کمانڈز اس اکاؤنٹ میں تبدیلی سے انکار کر دیتی ہیں۔ پہلے وہ سیشن بند کریں۔ جب گروپس کی کوئی تبدیلی ادھوری ہو تب بھی یہ انکار کرتی ہیں: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``۔

`ccam list --json` ہر اکاؤنٹ کو اس کے سلاٹ، شارٹ کٹ، ای میل اور فولڈر کے ساتھ پرنٹ کرتا ہے۔ اس کی `groups` array اکاؤنٹ کے گروپس کی فہرست دیتی ہے، اور اس کا `share` آبجیکٹ ہر آئٹم (`mcp`، `settings`، `memory`، `extensions`، `history`) کو اس گروپ سے جوڑتا ہے جس سے وہ آتا ہے، یا `local` سے۔

## شارٹ کٹس

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

شارٹ کٹ ایک کمانڈ ہے جس کا نام آپ خود رکھتے ہیں اور جو ایک اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے۔ `ccam work` کی طرح یہ ہر آرگیومنٹ بغیر تبدیلی کے Claude Code کو دیتی ہے، اس لیے `ccwork --resume` اور `ccwork -p "hi"` ویسے ہی کام کرتے ہیں جیسے آپ توقع کرتے ہیں۔ یہ `~/.local/bin` میں ccam کا ایک لنک ہے۔

شارٹ کٹس اختیاری ہیں۔ `ccam add` کا وزرڈ پوچھتا ہے کہ شارٹ کٹ بنایا جائے یا نہیں، `ccam add <name> --shortcut <command>` بغیر پوچھے ایک بنا دیتا ہے، اور `ccam list` ہر اکاؤنٹ کا شارٹ کٹ دکھاتا ہے، یا شارٹ کٹ نہ ہو تو `-`۔

ناموں میں a-z، 0-9، `.`، `-` اور `_` حروف استعمال ہوتے ہیں۔ نام کوئی ccam کمانڈ، کسی دوسرے اکاؤنٹ کا شارٹ کٹ، یا آپ کے `PATH` پر موجود کسی دوسرے پروگرام کا نام نہیں ہو سکتا، اور نہ ہی `~/.local/bin` میں کوئی ایسی فائل جو ccam نے نہ بنائی ہو، جیسے ccam کی آپ کی اپنی کاپی یا لنک۔ `config.toml` میں ہاتھ سے بدلا گیا کوئی شارٹ کٹ جو ان اصولوں کو توڑے، `claude` یا `ccam` ہو، یا دو اکاؤنٹس کو دیا گیا ہو، `config.toml` کو لوڈ ہونے سے روک دیتا ہے، اور اکاؤنٹ کا نام بتاتا ہے۔

`ccam rename` شارٹ کٹ کا نام برقرار رکھتا ہے، اور `ccam remove` اکاؤنٹ کے ساتھ شارٹ کٹ بھی حذف کر دیتا ہے۔ `ccam shims install` تمام اکاؤنٹس کے لنکس دوبارہ بناتا ہے، مثلاً جب آپ نے ccam کو کسی اور جگہ منتقل کیا ہو؛ کسی ایک شارٹ کٹ کی جگہ پر پڑی فائل باقیوں کو نہیں روکتی۔

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

`ccam migrate` کسی Claude Code فولڈر سے `settings.json`، `CLAUDE.md`، ایکسٹینشنز اور یوزر MCP سرورز کو ایک گروپ میں کاپی کرتا ہے، تاکہ وہ آئٹمز شیئر کرنے والے ہر اکاؤنٹ کو یہ مل جائیں۔ فولڈر کو صرف پڑھا جاتا ہے۔ آپشنز `--settings`، `--memory`، `--extensions`، `--mcp` اور `--history` صرف وہی آئٹمز کاپی کرتے ہیں جن کا آپ نام لیں۔ ان کے بغیر، ٹرمینل پر ایک چیک لسٹ سے آپ خود چن سکتے ہیں، جس میں گفتگو کی ہسٹری کے سوا سب کچھ چیک ہوتا ہے؛ باقی صورتوں میں گفتگو کی ہسٹری کے سوا سب کچھ کاپی ہوتا ہے۔ `-g` یا `--group` گروپ چنتا ہے، اور `-y` یا `--yes` بغیر پوچھے کاپی کرتا ہے اور فائل مختلف ہونے پر گروپ والی کاپی رکھتا ہے۔ MCP سائن اِن کاپی نہیں ہوتے، کیونکہ کاپی کیا گیا سائن اِن دوسرے رکھنے والوں کو لاگ آؤٹ کر دیتا ہے؛ پلان ان اکاؤنٹس کے نام بتاتا ہے جو `/mcp` سے سائن اِن کریں۔ `--mcp` تب بھی کام کرتا ہے جب `~/.claude` خالی ہو یا موجود نہ ہو لیکن `~/.claude.json` میں سرورز ہوں۔ [موجودہ سیٹ اپ کاپی کرنا](/docs/sharing#copy-an-existing-setup) دیکھیں۔

## MCP سرورز

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` وہ سرور جو صرف اسی اکاؤنٹ میں ہیں اس کے گروپ میں منتقل کرتا ہے، وہ بھی جو اس کے پاس صرف ایک فولڈر کے لیے ہیں (سادہ `claude mcp add`)، جو اس کے بعد ہر فولڈر میں لوڈ ہوتے ہیں۔ جس سرور کا نام گروپ پہلے سے کسی مختلف تعریف کے لیے استعمال کر رہا ہو اسے چھوڑ دیا جاتا ہے؛ گروپ کی کاپی بدلنے کے لیے `--force` شامل کریں، ccam پہلے اس کا بیک اپ لیتا ہے۔

## استعمال اور اکاؤنٹ تبدیل کرنا

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
ccam ccusage claude daily       # Claude Code's tokens and costs, every account (ccusage)
ccam ccusage -a work blocks --active   # one account's current 5-hour block
```

`next --auto`، `--fallback` اور `ccam pick` ایک ہی طرح چنتے ہیں: پہلے `auto_max_usage` (`config.toml`، ڈیفالٹ 100) سے کم والے اکاؤنٹس، `auto_strategy` (`headroom`، `expiring` یا `balance`) کے مطابق، `--model` کی ہفتہ وار ونڈو کو شمار کرتے ہوئے، کبھی 100% والا اکاؤنٹ نہیں (جس ونڈو کو Claude Code نے بھری ہوئی بتایا ہو وہ اپنے ری سیٹ تک 100% شمار ہوتی ہے)، کبھی API key لاگ اِن نہیں، اور اسی تنظیم کے اندر جس کے اکاؤنٹ سے وہ شروع کرتے ہیں (آگے دیکھنے کے لیے `--cross-org` یا `cross_org = true`)۔ `--wait` اپنی قدر `=` کے بعد لیتا ہے (`--wait=2h`)؛ اسپیس کے ساتھ لکھا جائے تو ایک اشارے کے ساتھ رد کر دیا جاتا ہے۔ جب کوئی اکاؤنٹ نہ چنا جا سکے تو یہ کوڈ 75 کے ساتھ باہر نکلتے ہیں، یا `--wait` کے ساتھ پہلے اکاؤنٹ کے ری سیٹ ہونے کا انتظار کرتے ہیں۔ `ccam next --continue` گفتگو کو ان اکاؤنٹس کے درمیان بھی لے جاتا ہے جو ہسٹری شیئر نہیں کرتے: تب Claude Code اسے نئے اکاؤنٹ میں کاپی کر دیتا ہے۔ `--brief` کے ساتھ نیا اکاؤنٹ اس کے بجائے ایک نئی گفتگو شروع کرتا ہے، اس خلاصے سے جو ccam کسی ماڈل کے بغیر پچھلی گفتگو کا لکھتا ہے ([`--brief`](/features/next#brief))۔ `ccam <name>` کے برعکس، `ccam next` اپنے آپشنز خود پڑھتا ہے، اس لیے Claude Code کے آپشنز `--` کے بعد آتے ہیں؛ اس کے بغیر ccam `unknown flag` کے ساتھ رک جاتا ہے۔ Claude Code کے اندر، پلگ اِن کا `/ccam:switch [account]` گفتگو آگے سونپتا ہے، اور `ccam hook` وہ ہے جو اس کے ہُکس چلاتے ہیں۔ `config.toml` میں `auto_switch = true` کے ساتھ (Linux اور macOS؛ ڈیفالٹ طور پر بند) حد والا ہُک ccam کے شروع کیے ہوئے انٹرایکٹو سیشن کو ختم بھی کرتا ہے اور اسے گنجائش والے کسی دوسرے اکاؤنٹ کے طور پر آگے لے جاتا ہے، صرف اسی جیسے اکاؤنٹ میں: اسی تنظیم کی کوئی دوسری سیٹ، یا کوئی دوسرا ذاتی پلان ([`auto_switch`](/features/next#auto-switch))۔ دیکھیں [اگلے اکاؤنٹ پر جانا](/features/next)، [اسکرپٹس کے لیے اکاؤنٹ چننا](/features/pick) اور [ہر اکاؤنٹ کے پس منظر کے سیشنز](/features/agents)۔

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
| `~/.config/ccam/config.toml` | اکاؤنٹس، گروپس اور سیٹنگز (`theme`، `usage_source`، `relay`، `auto_switch`، `cross_org` …)۔ سیٹنگز سب سے اوپر، پہلے `[[group]]`، `[[profile]]` یا `[[folder]]` سے اوپر آتی ہیں: کسی ٹیبل کے نیچے والی سیٹنگ اسی ٹیبل کی ہوتی ہے، اس لیے ccam ہر لانچ پر اور `ccam doctor` میں اس کا نام بتاتا ہے، اور کنفیگ محفوظ کرنے والی کمانڈز اس کے منتقل ہونے تک انکار کرتی ہیں۔ جب ccam فائل محفوظ کرتا ہے تو تبصرے برقرار نہیں رہتے |
| `~/.local/share/ccam/profiles/<name>/` | ہر اکاؤنٹ کا `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | شیئر کی گئی سیٹنگز، MCP سرورز، ایکسٹینشنز اور ہسٹری |
| `~/.local/state/ccam/` | استعمال کا کیش اور حد کے ریکارڈ (`usage/<name>.json`، `usage/<name>.limit.json`)، آخری استعمال کے اوقات، سیشن مارکرز اور لاکس (`locks/`)، وہ شارٹ کٹس جو ccam نے بنائے (`shortcuts.json`)، گروپس کی کوئی تبدیلی جو مکمل نہیں ہوئی (`group-change.json`)، وہ نوٹس جو `/ccam:switch` اور حد والا ہُک چھوڑتے ہیں (`handoff/<account>/<session>.json`)، `ccam pick` یا `auto_switch` نے ہر اکاؤنٹ کو آخری بار کب چنا (`picks/<name>`؛ کوئی انتخاب 2 منٹ تک شمار ہوتا ہے)، اور `--brief` کے خلاصے، جو 7 دن رکھے جاتے ہیں (`brief/<session>.md`) |
| `~/.local/bin/<shortcut>` | `ccwork` جیسی شارٹ کٹ کمانڈز، جو ccam کے لنکس ہیں |

`XDG_CONFIG_HOME`، `XDG_DATA_HOME` اور `XDG_STATE_HOME` ویری ایبلز ان فولڈرز کی جگہ بدل دیتے ہیں۔ Windows پر `~` سے مراد `%USERPROFILE%` ہے اور شارٹ کٹس `<shortcut>.exe` ہوتے ہیں۔
