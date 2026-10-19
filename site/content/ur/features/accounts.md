---
title: "اکاؤنٹس کا انتظام"
description: "ccam add، ccam login، ccam list، ccam rename، ccam remove اور ccam import سے Claude Code اکاؤنٹس شامل کریں، لاگ اِن کریں، ان کی فہرست دیکھیں، نام بدلیں، ہٹائیں اور امپورٹ کریں۔"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 97cef50f48f4
---

ccam کا ہر اکاؤنٹ Claude Code کا ایک الگ کنفیگ فولڈر ہے جس کا اپنا لاگ اِن ہوتا ہے۔ یہ کمانڈز اکاؤنٹس بناتی ہیں، ان میں سائن اِن کرتی ہیں، ان کی فہرست دکھاتی ہیں، ان کے نام بدلتی ہیں اور انہیں ہٹاتی ہیں، اور ان فولڈرز کو رجسٹر کرتی ہیں جو آپ کے پاس پہلے سے ہیں۔

## استعمال

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
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

جہاں بھی کوئی کمانڈ `<name|number>` لیتی ہے، آپ اکاؤنٹ کا نام یا اس کا سلاٹ نمبر دے سکتے ہیں، یعنی `ccam list` کا `#` کالم۔

ٹرمینل پر بغیر آپشنز کے `ccam add` ایک وزرڈ کھولتا ہے۔ یہ نام پوچھتا ہے، پھر شیئر کیے جانے والے آئٹمز کی ایک چیک لسٹ دکھاتا ہے (گفتگو کی ہسٹری کے سوا سب کچھ چیک ہوتا ہے)، پھر `Create a shortcut command for work? [y/N]` پوچھتا ہے اور `y` پر `Shortcut name [ccwork]:`۔ `--share`، `--shortcut` یا `--no-login` میں سے کوئی بھی، یا ٹرمینل کے بجائے pipe، وزرڈ کو چھوڑ دیتا ہے؛ تب نام دینا لازمی ہے اور کوئی شارٹ کٹ نہیں بنتا جب تک `--shortcut` کسی کا نام نہ دے۔

`ccam list` سلاٹ، نام، شارٹ کٹ، ای میل، شیئر کیے گئے آئٹمز اور آخری استعمال دکھاتا ہے، اور ڈیفالٹ اکاؤنٹ سے پہلے `*` ہوتا ہے:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` کے ساتھ ہر اکاؤنٹ ایک آبجیکٹ ہے جس میں `slot`، `shortcut`، `name`، `email`، `default`، `dir`، `last_used` (جب اکاؤنٹ کبھی شروع نہ ہوا ہو تو شامل نہیں) اور `share` ہوتے ہیں۔ `share` ہر آئٹم (`mcp`، `settings`، `memory`، `extensions`، `history`) کو اس کے گروپ کے نام یا `"local"` سے جوڑتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `ccam add --no-login` | لاگ اِن کیے بغیر اکاؤنٹ بناتا ہے؛ بعد میں `ccam login <name>` سے لاگ اِن کریں | بند |
| `ccam add --share <item=group>` | کسی آئٹم کو ایک گروپ کے ساتھ شیئر کرتا ہے، یا اسے نجی رکھنے کے لیے `item=local`۔ کئی بار دیا جا سکتا ہے | `mcp`، `settings`، `memory`، `extensions`، `default` میں؛ `history` لوکل |
| `ccam add --shortcut <command>` | `~/.local/bin` میں یہ شارٹ کٹ کمانڈ بھی بناتا ہے | کوئی نہیں |
| `ccam list --json` | اکاؤنٹس کو JSON میں پرنٹ کرتا ہے | بند |
| `ccam remove -y`، `--yes` | پوچھے بغیر ہٹاتا ہے | بند |
| `ccam import --name <name>` | آپ کے دیے گئے فولڈر کے لیے اکاؤنٹ کا نام | فولڈر کے نام سے اخذ |
| `ccam import -y`، `--yes` | ملنے والا ہر فولڈر پوچھے بغیر امپورٹ کرتا ہے | بند |
| `--dry-run` | `ccam list` اور `ccam completion` اسے لاگو کرتے ہیں، کیونکہ وہ کچھ تبدیل نہیں کرتے۔ `add`، `login`، `logout`، `use`، `rename`، `remove` اور `import` اسے `` --dry-run is not supported by `ccam add`; nothing was changed `` کے ساتھ مسترد کرتے ہیں (کمانڈ کے نام کے ساتھ) | بند |

`ccam login`، `ccam logout`، `ccam use` اور `ccam rename` کے اپنے کوئی آپشنز نہیں۔ `ccam completion`، `bash` یا `zsh` لیتا ہے؛ اس کے علاوہ کچھ بھی `only bash and zsh are supported` کے ساتھ ناکام ہوتا ہے۔

اکاؤنٹ کے نام `a-z`، `0-9`، `-` اور `_` کے 1 سے 32 حروف ہوتے ہیں، جو کسی حرف یا ہندسے سے شروع ہوں۔ نام صرف ہندسوں پر مشتمل نہیں ہو سکتا، کیونکہ `ccam 2` کا مطلب سلاٹ 2 ہے۔ نام کوئی ccam کمانڈ بھی نہیں ہو سکتا (`add`، `list`، `ls`، `run`، `use`، `next`، `usage`، `u`، `share`، `mcp`، `login`، `logout`، `rename`، `remove`، `rm`، `import`، `migrate`، `doctor`، `shims`، `shortcut`، `completion`، `help`، `version`، `statusline`، `update`)، اس لیے `ccam <name>` کبھی مبہم نہیں ہوتا۔ یہ فہرست `internal/config` میں `config.Reserved` ہے۔

## یہ کیسے کام کرتا ہے

`ccam add work --shortcut ccwork`، `addProfile` (`internal/cli/profiles.go`) میں یہ مراحل چلاتا ہے:

<!-- motion:accounts-add -->

1. کچھ بھی بنانے سے پہلے `config.ValidateName` نام جانچتا ہے، اور `shim.Check` شارٹ کٹ جانچتا ہے۔
2. `profile.Add`، `~/.local/share/ccam/profiles/work/` کو موڈ `0700` کے ساتھ بناتا ہے اور اکاؤنٹ کو سب سے کم خالی سلاٹ نمبر (`Config.NextSlot`) دیتا ہے، اس لیے ہٹائے گئے اکاؤنٹ کا نمبر دوبارہ استعمال ہوتا ہے۔ پہلا اکاؤنٹ ڈیفالٹ بن جاتا ہے۔
3. `applyShare` نئے، خالی فولڈر کو `share.Build` اور پھر `Apply` سے اس کے شیئر گروپس سے لنک کرتا ہے، مثلاً `settings.json` کو `groups/default/settings.json` سے۔ اگر یہ ناکام ہو تو نیا فولڈر ہٹا دیا جاتا ہے، کیونکہ بچا ہوا فولڈر ہر دوبارہ کوشش کو `folder already exists` کے ساتھ روک دے گا۔
4. اکاؤنٹ `fsutil.AtomicWrite` کے ذریعے `~/.config/ccam/config.toml` میں محفوظ ہوتا ہے: ایک عارضی فائل، `fsync`، پھر `rename`، اس لیے کریش کبھی آدھا کنفیگ نہیں چھوڑتا۔
5. `runLogin`، `CLAUDE_CONFIG_DIR` کو اکاؤنٹ فولڈر پر مقرر کر کے `claude auth login` شروع کرتا ہے۔ اگر لاگ اِن ناکام ہو تو اکاؤنٹ برقرار رہتا ہے، اور ccam غلطی پرنٹ کر کے بتاتا ہے کہ بعد میں `ccam login work` چلائیں۔
6. `profile.ReadEmail`، `profiles/work/.claude.json` سے `oauthAccount.emailAddress` نکالتا ہے۔ یہ فائل کئی میگا بائٹ کی ہو سکتی ہے، اس لیے `TopLevelValue` پوری فائل پارس کرنے کے بجائے صرف اسی ایک key کو تلاش کرتا ہے۔
7. `shim.Install`، `~/.local/bin/ccwork` کو ccam بائنری سے لنک کرتا ہے، جو شروع ہوتے وقت `argv[0]` کو `config.toml` میں تلاش کرتی ہے۔

لاگ اِن میں کئی منٹ لگ سکتے ہیں، اور اس دوران کوئی دوسرا ccam کنفیگ محفوظ کر سکتا ہے۔ اس لیے انتظار کرنے والی کمانڈز `env.update` کے ذریعے محفوظ کرتی ہیں، جو `config.toml` دوبارہ لوڈ کرتا ہے، صرف ایک تبدیلی (یہاں، ای میل) لاگو کرتا ہے اور محفوظ کرتا ہے۔ اسی وجہ سے `ccam remove` اپنے سوال کے بعد دوبارہ لوڈ کرتا ہے۔

`ccam login` کسی موجودہ اکاؤنٹ کے لیے مرحلہ 5 اور 6 چلاتا ہے۔ `ccam logout` اکاؤنٹ میں `claude auth logout` چلاتا ہے اور محفوظ ای میل صاف کر دیتا ہے۔ `ccam use` صرف `config.toml` میں `default` دوبارہ لکھتا ہے۔

### ریموٹ مشین پر سائن اِن

SSH پر `claude auth login` کا براؤزر والا حصہ ٹوٹ سکتا ہے۔ Claude Code سرور پر `http://localhost:<port>/callback` کے ایک callback سرور پر سائن اِن کا انتظار کرتا ہے، لیکن جو براؤزر سائن اِن مکمل کرتا ہے وہ آپ کے کمپیوٹر پر چلتا ہے، جہاں یہ پتا کہیں نہیں جاتا اور براؤزر `ERR_CONNECTION_REFUSED` جیسی غلطی دکھاتا ہے۔ اس پتے میں پھر بھی ایک بار استعمال ہونے والا کوڈ ہوتا ہے، اس لیے ccam اسے پہنچا سکتا ہے:

<!-- motion:accounts-relay -->

`runLogin` آپ کا ٹرمینل Claude Code کے حوالے نہیں کرتا۔ یہ Claude Code کے اسٹینڈرڈ اِن پٹ کو ایک pipe سے جوڑتا ہے اور `relayLogin` آپ کی ٹائپ کی ہوئی ہر لائن اس میں کاپی کرتا ہے، اس لیے Claude Code کے `Paste code here` پرامپٹ پر چسپاں کیا گیا کوڈ پھر بھی پہنچ جاتا ہے۔ جب کوئی لائن callback پتا ہو تو `callbackURL` اسے صرف تب قبول کرتا ہے جب ہوسٹ `localhost`، `127.0.0.1` یا `[::1]` ہو، پورٹ دیا گیا ہو، راستہ `/callback` ہو اور `code` اور `state` دونوں موجود ہوں؛ `http://` چھوڑا جا سکتا ہے۔ پھر `deliverCallback` سرور پر اس پتے کی درخواست 10 سیکنڈ کی حد کے ساتھ اور ری ڈائریکٹس پر عمل کیے بغیر کرتا ہے، جیسے براؤزر کرتا۔ درخواست انتظار کرنے والے `claude auth login` کے callback سرور تک پہنچتی ہے (ڈایاگرام میں ڈیش والی لائن)۔ Claude Code کوڈ وصول کر کے لاگ اِن محفوظ کرتا ہے، Linux پر اکاؤنٹ فولڈر کی `.credentials.json` میں۔ چسپاں کی گئی کوئی لائن کبھی ccam کو کسی دوسرے ہوسٹ سے رابطہ نہیں کروا سکتی۔

لاگ اِن شروع ہونے سے پہلے ccam یہ پرنٹ کرتا ہے:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

کامیابی پر یہ `→ Passed the sign-in to Claude Code on this machine.` پرنٹ کرتا ہے۔

### نام بدلنا، ہٹانا اور امپورٹ کرنا

`ccam rename` اور `ccam remove` اکاؤنٹ کے فولڈر پر اور ان چیزوں پر عمل کرتے ہیں جو ccam اکاؤنٹ کے نام سے دوسری جگہوں پر رکھتا ہے۔ `nameState` ان راستوں کی فہرست ہے: `~/.local/state/ccam/last-used/<name>`، `~/.local/state/ccam/mcp/<name>` اور استعمال کے کیش `~/.local/state/ccam/usage/<name>.json` اور `<name>.statusline.json`۔ `ccam remove work` انہیں اس ترتیب سے سنبھالتا ہے:

<!-- motion:accounts-remove -->

`ccam remove` پہلے انکار کرتا ہے اگر اکاؤنٹ چل رہا ہو، پھر پوچھتا ہے۔ اگر اکاؤنٹ ایکسٹینشنز شیئر کرتا ہو تو `RewritePluginPaths` گروپ کی `installed_plugins.json` اور `known_marketplaces.json` دوبارہ لکھتا ہے، کیونکہ Claude Code پلگ اِنز کے مکمل راستے محفوظ کرتا ہے جو اس اکاؤنٹ کے فولڈر سے گزر سکتے ہیں۔ `profile.Trash` فولڈر کو `~/.local/share/ccam/trash/<name>-<time>/` میں منتقل کرتا ہے، اکاؤنٹ `config.toml` سے نکل جاتا ہے، `nameState` کے راستے حذف ہو جاتے ہیں تاکہ بعد میں اسی نام کا اکاؤنٹ انہیں وراثت میں نہ پائے، اور شارٹ کٹ کا لنک ہٹا دیا جاتا ہے۔ واپس لانے کے لیے فولڈر کو trash سے واپس منتقل کریں۔

`ccam rename` فولڈر کا نام بدلتا ہے، اس کے `plugins/` میں پلگ اِن کے راستے پرانے فولڈر سے نئے فولڈر پر دوبارہ لکھتا ہے، ڈیفالٹ کو ساتھ لے جاتا ہے، اور `nameState` کے راستوں کو نئے نام پر منتقل کرتا ہے۔ شارٹ کٹ کمانڈ اپنا نام برقرار رکھتی ہے۔ اگر پلگ اِن کے راستے دوبارہ لکھنا ناکام ہو تو فولڈر کا نام واپس بدل دیا جاتا ہے، کیونکہ کنفیگ محفوظ نہیں ہوتا۔

`ccam import` کسی فولڈر کو اس کی جگہ پر ہی، `config.toml` میں `dir` کے طور پر رجسٹر کرتا ہے۔ یہ اسے کبھی منتقل یا کاپی نہیں کرتا: کاپی کیا گیا ریفریش ٹوکن کاپیوں کو ایک دوسرے کو لاگ آؤٹ کروا دے گا۔ فولڈر آرگیومنٹ کے بغیر یہ غیر رجسٹر شدہ `~/.claude` اور `~/.claude-*` فولڈرز تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے؛ نام فولڈر سے آتے ہیں (`~/.claude` سے `main`، `~/.claude-work` سے `work`، `~/.claude-2` سے `acct2`)، اور جس فولڈر کا نام اخذ نہ ہو سکے وہ استعمال کرنے کے لیے `--name` کمانڈ کے ساتھ دکھایا جاتا ہے۔ `CLAUDE_CONFIG_DIR` کے بغیر Claude Code، `.claude.json` کو آپ کے ہوم فولڈر میں `~/.claude` کے ساتھ رکھتا ہے، اس لیے جب آپ `~/.claude` امپورٹ کریں اور اس میں اپنی `.claude.json` نہ ہو تو ccam، `~/.claude/.claude.json` سے `~/.claude.json` کا لنک بناتا ہے۔ تب سادہ `claude` اور `ccam main` ایک ہی اکاؤنٹ دیکھتے ہیں۔ امپورٹ کیے گئے اکاؤنٹس تب تک کچھ شیئر نہیں کرتے جب تک آپ `ccam share` نہ چلائیں، اور کوئی شارٹ کٹ نہیں بنتا۔

`ccam completion bash|zsh` ایک completion اسکرپٹ پرنٹ کرتا ہے جو اکاؤنٹس کے نام بھی مکمل کرتا ہے۔

## حفاظت اور خاص صورتیں

- `ccam rename` اور `ccam remove` اکاؤنٹ چلنے کے دوران (`refuseRunning`) `work is running; quit its Claude Code session first` کے ساتھ انکار کرتے ہیں، کیونکہ سیشن فولڈر میں لکھتا رہتا ہے۔
- امپورٹ کیے گئے اکاؤنٹ پر `ccam remove` صرف اس کی رجسٹریشن ختم کرتا ہے: `✓ main unregistered (folder ... left in place)`۔ آپ کا `~/.claude`، trash میں منتقل نہیں ہوتا۔
- `ccam add` ایسا نام مسترد کرتا ہے جس کا فولڈر `profiles/` میں پہلے سے موجود ہو، اور `ccam import` کی طرف رہنمائی کرتا ہے۔
- `ccam import` پہلے سے رجسٹر شدہ فولڈر، ایسا راستہ جو فولڈر نہ ہو، اور پہلے سے لیا گیا نام مسترد کرتا ہے۔
- `ccam logout` اور `ccam login` یہ نہیں جانچتے کہ اکاؤنٹ چل رہا ہے یا نہیں۔
- اگر ریلے Claude Code تک نہ پہنچ سکے، مثلاً لاگ اِن ختم ہونے کے بعد، تو یہ وجہ کے ساتھ `ccam: could not pass the sign-in to Claude Code:` پرنٹ کرتا ہے۔ HTTP غلطی کا عام طور پر مطلب ہے کہ کوڈ کی میعاد ختم ہو گئی؛ لاگ اِن دوبارہ شروع کریں۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [شارٹ کٹ کمانڈز](/features/shortcuts)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [Claude Code کے ایک سے زیادہ اکاؤنٹس](/guides/multiple-claude-code-accounts)
