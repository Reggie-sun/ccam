---
title: "جانچ اور مرمت"
description: "ccam doctor سے شیئر لنکس، فائلوں کی اجازتیں، لاگ اِنز اور Claude Code کا ورژن جانچیں، اور ccam doctor --fix سے وہ ٹھیک کریں جسے ٹھیک کرنا محفوظ ہو۔"
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` آپ کی ccam انسٹالیشن جانچتا ہے: Claude Code بائنری، ہر اکاؤنٹ کا فولڈر، لاگ اِن اور شیئر لنکس، شیئر کیے گئے MCP سرورز، usage اینڈ پوائنٹ اور آپ کی شارٹ کٹ کمانڈز۔ `ccam doctor --fix` وہ ٹھیک کرتا ہے جو ڈیٹا کھوئے بغیر ہو سکے، اور جو فائل بھی وہ بدلتا ہے اسے پہلے ایک بیک اپ فولڈر میں منتقل کرتا ہے۔ اسے تب چلائیں جب `ccam <name>` کسی شیئر لنک کے مسئلے پر رک جائے، جب آپ نے فائلیں ہاتھ سے بحال کی ہوں، یا Claude Code کی اپ ڈیٹ کے بعد۔

## استعمال

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

پہلی لائن ایک ساتھ ہر اکاؤنٹ پر نشان لگاتی ہے، تاکہ کچھ اور پڑھنے سے پہلے ہی آپ دیکھ لیں کہ مسائل کہاں ہیں۔ اس کے نیچے ہر اکاؤنٹ کو ایک باکس ملتا ہے، اور جو نتائج کسی اکاؤنٹ سے متعلق نہیں (Claude Code، ماحولیاتی متغیرات، پلگ اِنز، MCP سرورز، شارٹ کٹس، فولڈر کے اصول، گروپس اور سیٹنگز، `auto_switch`) وہ `General` باکس میں جاتے ہیں۔ باکس کا عنوان اکاؤنٹ کا سب سے برا نشان اور اس کے مسائل کی تعداد دکھاتا ہے۔ اندر کی ہر قطار ایک جانچا گیا شعبہ ہے، ہمیشہ پہلے `Folder`، `Login` اور `Share links` اسی ترتیب میں، پھر باقی، ایک نشان کے ساتھ: `✔` ٹھیک، `ℹ` معلومات، `⚠` انتباہ، `✘` غلطی۔ جس فولڈر میں کوئی خرابی نہ ہو وہ `✔ Folder  OK` دکھاتا ہے۔

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix` کے ساتھ، جو بھی نتیجہ ٹھیک کیا گیا ہو اس کے آخر میں `→ fixed` آتا ہے، اور وہ مزید مسئلہ نہیں گنا جاتا:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

ٹرمینل پر باکس اس کی چوڑائی میں سما جاتے ہیں اور لمبے پیغامات اگلی لائن میں چلے جاتے ہیں، ٹرمینل کے اپنے رنگوں میں (`NO_COLOR` انہیں بند کر دیتا ہے)۔ کوئی باکس کبھی 29 کالمز سے کم چوڑا نہیں ہوتا، اس لیے تنگ ٹرمینل پر کوئی لفظ نہیں کٹتا۔ پائپ یا فائل میں بھیجنے پر آؤٹ پٹ میں کوئی رنگ نہیں ہوتا اور کوئی پیغام توڑا یا کاٹا نہیں جاتا، اس لیے `ccam doctor | grep settings.json` پوری لائن ڈھونڈ لیتا ہے؛ تب ہر باکس سب سے لمبی لائن جتنا چوڑا ہوتا ہے۔ ڈیش بورڈ کا Health ٹیب یہی باکس دکھاتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--fix` | وہ ٹھیک کرتا ہے جسے ٹھیک کرنا محفوظ ہو (بدلی گئی فائلوں کا بیک اپ لیا جاتا ہے) | بند |
| `--dry-run` (گلوبل) | ہر جانچ چلاتا ہے لیکن کچھ ٹھیک نہیں کرتا، `--fix` کے ساتھ بھی | بند |

`ccam doctor` کوئی آرگیومنٹ نہیں لیتا۔ یہ `--dry-run` کو لاگو کرتا ہے: `ccam --dry-run doctor --fix`، `ccam doctor` کے نتائج پرنٹ کرتا ہے، یہ بھی بتاتا ہے کہ `--fix` کیا الگ رکھ دے گا، اور کچھ تبدیل نہیں کرتا:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

جب کم از کم ایک `✘` نتیجہ ٹھیک کیے بغیر رہ جائے تو ایگزٹ کوڈ `1` ہوتا ہے، ورنہ `0`۔ انتباہات اور معلومات کبھی رن کو ناکام نہیں کرتیں۔ جب رن `--fix` کے بغیر ناکام ہو تو ccam ایک خالی لائن اور `` `ccam doctor --fix` repairs what it can. `` پر ختم ہوتا ہے۔ جو کنفیگ لوڈ نہ ہو سکے وہ بھی `1` اور ایک `ccam:` غلطی پر ختم ہوتا ہے۔

## یہ کیسے کام کرتا ہے

`internal/cli/doctor.go`، `~/.config/ccam/config.toml` لوڈ کرتا ہے اور `internal/doctor/doctor.go` میں `doctor.Run` کو کال کرتا ہے، جو جانچیں ایک مقررہ ترتیب میں چلاتا ہے اور نتائج لوٹاتا ہے، ہر ایک اس اکاؤنٹ کے ساتھ جس سے وہ متعلق ہے اور اس شعبے کے ساتھ جو اس نے جانچا۔ `tui.HealthReport` انہیں باکسز میں گروپ کرتا ہے۔ `--fix` آگے `Options.Fix` کے طور پر جاتا ہے، جو ڈرائی رن میں false ہوتا ہے۔

<!-- motion:doctor-checks -->

درمیانی حصہ ہر اکاؤنٹ کے لیے ایک بار چلتا ہے۔ جانچیں، اور ہر ایک کے لیے `--fix` کیا کرتا ہے:

1. **Claude Code کا ورژن** (`checkClaude`)۔ `claude` بائنری کو اسی طرح تلاش کرتا ہے جیسے لانچر کرتا ہے اور `claude --version` چلاتا ہے۔ بائنری کا نہ ہونا، ایسا `claude_bin` جو قابلِ عمل فائل نہ ہو، یا `--version` کا ناکام ہونا غلطی ہے۔ کوئی مرمت نہیں۔
2. **انوائرنمنٹ میں اسناد** (`checkAuthEnv`)۔ جب `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` مقرر ہو تو انتباہ دیتا ہے: اس شیل سے لانچ ہونے والا ہر اکاؤنٹ اسے وراثت میں پاتا ہے، اور Claude Code اسے اکاؤنٹ کے اپنے لاگ اِن پر ترجیح دیتا ہے۔ کوئی مرمت نہیں؛ اسے خود ہٹائیں۔
3. **فولڈر** (`checkProfile`)۔ اکاؤنٹ فولڈر کا نہ ہونا غلطی ہے اور اس اکاؤنٹ کی جانچ ختم کر دیتا ہے۔ اگر اجازتوں میں group یا other کا کوئی بٹ مقرر ہو تو انتباہ ہے؛ `--fix`، `chmod 0700` چلاتا ہے۔ امپورٹ کیا گیا `~/.claude` جس میں `~/.claude.json` کے علاوہ اپنی ایک `.claude.json` بھی ہو، معلومات ہے: تب `ccam main` اور سادہ `claude` الگ الگ لاگ اِن تفصیلات اور MCP سرورز رکھتے ہیں۔ Windows میں اجازت کے بٹس نہیں ہوتے، اس لیے وہاں یہ جانچ اور نیچے والی `0600` جانچ چھوڑ دی جاتی ہے؛ آپ کے یوزر پروفائل میں موجود فولڈر اپنی ایکسیس لسٹ سے محفوظ رہتا ہے۔
4. **لاگ اِن فائل** (`checkCredentials`)۔ `.credentials.json` کو `lstat` سے جانچا جاتا ہے۔ symlink غلطی ہے، کیونکہ ہر اکاؤنٹ کا اپنا ٹوکن ہونا چاہیے، اور اسے ٹھیک نہیں کیا جاتا۔ `0600` کے سوا کوئی موڈ انتباہ ہے؛ `--fix`، `chmod 0600` چلاتا ہے۔ پھر لاگ اِن پڑھا جاتا ہے (macOS پر پہلے Keychain سے)۔ لاگ اِن نہ ہونا، یا ایسی لاگ اِن فائل جس میں لاگ اِن نہ ہو، انتباہ ہے جو `ccam login <name>` کی طرف رہنمائی کرتا ہے۔ میعاد ختم ہونے والا ریفریش ٹوکن، یا ایسا لاگ اِن جسے سرور کے انکار کے بعد Claude Code نے مٹا دیا ہو، `` login expired → `ccam login <name>` `` غلطی ہے؛ میعاد ختم ہونے والا ایکسیس ٹوکن صرف معلومات ہے، کیونکہ Claude Code اگلے لانچ پر اس کی تجدید کر دیتا ہے۔ Windows Credential Manager میں رکھا گیا لاگ اِن معلومات ہے، کیونکہ ccam اسے پڑھ نہیں سکتا۔ درست لاگ اِن وہ ای میل دکھاتا ہے جس سے فولڈر اس وقت لاگ اِن ہے، اور پلان (`login valid (dev@company.com, Claude Team · Acme)`)۔ ان میں سے کوئی بھی ٹھیک نہیں کیا جاتا؛ `ccam login --expired` ہر اس اکاؤنٹ میں لاگ اِن کرتا ہے جسے اس کی ضرورت ہو۔
5. **شیئر لنکس** (`share.Check`)۔ `share.Layout` کے ہر راستے کو اکاؤنٹ کے گروپس سے موازنہ کرتا ہے، صرف `lstat` اور `readlink` استعمال کرتے ہوئے۔ یہ پانچ قسمیں بتاتا ہے: غائب لنک، کسی اور جگہ اشارہ کرنے والا لنک، جہاں لنک ہونا چاہیے وہاں اصل فائل، کسی ایسے آئٹم کا بچا ہوا لنک جو اب لوکل ہے، اور درست لنک جس کا گروپ فولڈر غائب ہو (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`؛ لانچ اسے نہیں جانچتا)۔ غائب لنک اور غائب گروپ فولڈر انتباہ ہیں، باقی غلطیاں ہیں۔ گروپ سے واپس کاپی کا کوئی عمل جو ادھورا رک گیا ہو، وہ بھی بتایا جاتا ہے: `--fix` تیار رکھی گئی مکمل کاپی کو اس کی جگہ لگا دیتا ہے، لیکن اگر اس کے بعد وہ راستہ (کسی لانچ نے) نئے سرے سے بنا دیا ہو تو یہ مکمل کاپی والے فولڈر کا نام بتاتا ہے تاکہ آپ خود اسے سلجھا لیں۔ ccam کے کسی بعد کے ورژن میں کسی آئٹم میں شامل کیے گئے راستے (ہسٹری کے لیے `file-history/`، `tasks/` اور `paste-cache/`، سیٹنگز کے لیے `keybindings.json`) صرف انتباہ ہیں جب تک اکاؤنٹ کے پاس ابھی اپنی کاپی ہو: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`۔ لانچ ان کی وجہ سے کبھی نہیں رکتا۔
6. **نامعلوم اینٹریز** (`checkUnknownEntries`)۔ اکاؤنٹ فولڈر میں اوپری سطح کے وہ نام دکھاتا ہے جو ccam نہیں جانتا (`share.Known`)، اس اشارے کے طور پر کہ Claude Code کی ترتیب بدل گئی ہے۔ صرف معلومات۔
7. **پلگ اِن کے راستے** (`checkPlugins`)۔ پلگ اِن رجسٹریز (`installed_plugins.json`، `known_marketplaces.json`) مکمل راستے محفوظ کرتی ہیں۔ ہر گروپ کے `plugins/` میں، جو راستے کسی ایک اکاؤنٹ کے فولڈر سے گزرتے ہیں وہ اس اکاؤنٹ کے ہٹائے جانے یا نام بدلنے پر ٹوٹ جاتے ہیں۔ `--fix` انہیں `profile.RewritePluginPaths` سے گروپ فولڈر پر دوبارہ لکھتا ہے، جو رجسٹری فائل کی ایٹامک تحریر ہے۔
8. **MCP سرورز** (`checkMCP`)۔ ہر گروپ کی `mcp.json` ایک بار لوڈ کرتا ہے۔ جو فائل پڑھی نہ جا سکے وہ غلطی ہے؛ جو سرور `mcp.Validate` میں ناکام ہو وہ انتباہ ہے۔ سادہ `claude` کی `~/.claude.json` میں ایسے سرورز جنہیں کوئی اکاؤنٹ لوڈ نہیں کرتا، انتباہ ہیں جو وہ `ccam migrate --mcp` کمانڈ بتاتا ہے جو انہیں کاپی کرتی ہے (جب کئی گروپس MCP شیئر کرتے ہوں تو `--group <a|b>` کے ساتھ)۔ جو سرورز کسی گروپ کے صرف ایک اکاؤنٹ کے پاس ہوں وہ معلومات ہیں، اس `ccam mcp import` کمانڈ کے ساتھ جو انہیں شیئر کرتی ہے، اور جو اکاؤنٹ ان MCP سرورز میں سائن اِن نہ ہو جن میں دوسرے اکاؤنٹس سائن اِن ہیں، اسے `/mcp` سے سائن اِن کرنے کا انتباہ ملتا ہے (ccam کبھی سائن اِن کاپی نہیں کرتا)۔ کوئی مرمت نہیں۔
9. **usage جواب کا فارمیٹ** (`checkUsageFormat`)۔ کیش کے بغیر، usage اینڈ پوائنٹ سے ایک کے بعد ایک اکاؤنٹ کے لیے پوچھتا ہے یہاں تک کہ کوئی جواب دے، یہ سب 4 سیکنڈ کی ایک ہی حد کے اندر۔ `usage_source = "claude"` کے ساتھ یہ اس کے بجائے ایک اکاؤنٹ کے لیے Claude Code سے پوچھتا ہے۔ اگر جواب اب پارس نہ ہو تو یہ انتباہ دیتا ہے کہ `ccam usage` اسٹیٹس لائن کیش پر چلا جاتا ہے؛ جب کسی اکاؤنٹ نے جواب نہ دیا ہو تو یہ بات معلومات کے طور پر بتاتا ہے (`could not check the usage response format: …`)۔ یہ واحد جانچ ہے جو نیٹ ورک استعمال کرتی ہے۔
10. **شارٹ کٹس** (`checkShims`)۔ صرف تب جب کسی اکاؤنٹ کا شارٹ کٹ ہو۔ `~/.local/bin` میں غائب لنکس انتباہ ہیں، اور ccam کے بنائے ہوئے وہ شارٹ کٹس بھی جو ابھی تک کوئی پرانا ccam چلاتے ہیں، جو بدل دیا گیا یا منتقل کر دیا گیا ہو (`shortcuts in <bin> that run another ccam (one replaced or moved): …`)؛ `--fix`، `shim.Install` کو کال کرتا ہے، جو صرف وہی لنکس بناتا اور نئے کرتا ہے جو ccam نے بنائے اور `~/.local/state/ccam/shortcuts.json` میں درج کیے، اور دوسری فائلوں کو کبھی نہیں چھیڑتا۔ انتباہ `→ fixed` پر صرف تب ختم ہوتا ہے جب ان میں سے ہر لنک بن گیا ہو؛ اگر کوئی نہ بن سکا ہو تو وجہ شامل کر دی جاتی ہے۔ جو شارٹ کٹ کوئی ایسا دوسرا ccam چلاتا ہو جو اب بھی انسٹال ہے، وہ کام کرتا ہے اور معلومات ہے؛ `--fix` اسے اس ccam کی طرف کر دیتا ہے۔ جو bin فولڈر `PATH` پر نہ ہو وہ انتباہ ہے اور ٹھیک نہیں کیا جاتا۔
11. **فولڈرز** (`checkFolders`)۔ صرف تب جب [`ccam folder`](/features/folder) اندراجات موجود ہوں۔ جس اندراج کا اکاؤنٹ موجود نہ ہو وہ غلطی ہے، اور جس اندراج کا فولڈر ختم ہو چکا ہو وہ انتباہ ہے۔ دونوں میں سے کوئی ٹھیک نہیں کیا جاتا، کیونکہ اندراج ہٹانے سے بدل جائے گا کہ وہاں کون سا اکاؤنٹ شروع ہوتا ہے۔
12. **گروپس اور سیٹنگز** (`checkGroups`)۔ گروپ کا ایسا رکن جو کوئی اکاؤنٹ نہ ہو، ایسا اکاؤنٹ جس کا نام بعد میں ccam کمانڈ بن گیا ہو (اسے `ccam run <name>` سے شروع کریں یا نام بدلیں)، نامعلوم `theme` یا `usage_source`، اور `config.toml` کی وہ لائنیں جو ccam نہیں پڑھتا، انتباہ ہیں۔ آخری والی عموماً کسی `[[group]]`، `[[profile]]` یا `[[folder]]` ٹیبل کے نیچے لکھی گئی کوئی سیٹنگ ہوتی ہے، جہاں وہ اسی ٹیبل کا حصہ بن جاتی ہے (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`)؛ کنفیگ بدلنے والی کمانڈز اس کے ٹھیک ہونے تک انکار کرتی ہیں۔ کوئی مرمت نہیں۔
13. **لاگ اِنز** (`checkIdentities`)۔ ہر فولڈر کی `.claude.json` سے پڑھتا ہے کہ وہ کس اکاؤنٹ میں لاگ اِن ہے۔ جو فولڈر اب اس اکاؤنٹ کے بجائے کسی دوسرے اکاؤنٹ میں لاگ اِن ہو جو `ccam login` نے درج کیا تھا (مثلاً سیشن کے اندر `/login` کے بعد) وہ انتباہ ہے؛ جس اکاؤنٹ کا لاگ اِن ریکارڈ نہ ہو، جیسے امپورٹ کیا گیا یا کسی پرانے ccam کا بنایا ہوا، اس کا موازنہ اس کی درج شدہ ای میل سے کیا جاتا ہے۔ جو فولڈر کسی دوسرے فولڈر والے Claude اکاؤنٹ میں ہی لاگ اِن ہو وہ بھی انتباہ ہے: دونوں اپنی استعمال کی حدیں شیئر کرتے ہیں۔ جس فولڈر میں وہی لاگ اِن ٹوکن ہو جو کسی دوسرے میں ہے، یعنی کاپی کیا گیا فولڈر (سادہ `claude` کا `~/.claude` بھی)، وہ غلطی ہے، کیونکہ جو بھی پہلے اس کی تجدید کرے وہ دوسرے کو لاگ آؤٹ کر دیتا ہے؛ اسے اپنا لاگ اِن دینے کے لیے `ccam login <name>` چلائیں۔ کوئی مرمت نہیں۔
14. **`auto_switch`** (`checkAutoSwitch`)۔ صرف تب جب `auto_switch = true` ہو۔ `General` باکس میں معلومات بتاتی ہے کہ `auto_switch` کیا کرتا ہے اور [اس کی تفصیل](/features/next#auto-switch) کا لنک دیتی ہے۔ اکاؤنٹس میں سیٹوں والے ہر پلان (Claude Team، Claude Enterprise) کے لیے ایک لائن ان کے نام بتاتی ہے، کہتی ہے کہ `auto_switch` سیشن کو ایک ہی تنظیم کی سیٹوں کے درمیان منتقل کرتا ہے، ہیلپ سینٹر کے لنک کے ساتھ Anthropic کا سیٹ والا اصول نقل کرتی ہے، اور فیصلہ آپ پر چھوڑتی ہے۔ 2.1.214 سے پرانا Claude Code، جو ہُکس کو اپنا پروسیس (`CLAUDE_PID`) نہیں بتاتا، وہاں انتباہ ہے، کیونکہ کوئی سیشن منتقل نہیں ہو سکتا۔ جس اکاؤنٹ میں ccam پلگ اِن نہ ہو، یا پلگ اِن بند ہو، اسے اپنے باکس میں انتباہ ملتا ہے، اس کمانڈ کے ساتھ جو Claude Code میں اسے ٹھیک کرتی ہے: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (یا `/plugin enable ccam@ccam`)۔ API key اکاؤنٹس شامل نہیں کیے جاتے۔ Windows پر واحد لائن یہ ہے: `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` کوئی مرمت نہیں۔

اکاؤنٹس سے پہلے، **گروپس کی تبدیلی** جو ادھوری رک گئی ہو `General` باکس میں غلطی ہے: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``۔ `--fix` اسے اسی ترتیب میں مکمل کرتا ہے جس میں اس کا منصوبہ بنا تھا، اس کے ٹکراؤ کے لیے دیے گئے جوابات کے ساتھ؛ تب تک گروپس کی کوئی اور تبدیلی، نام بدلنا یا ہٹانا شروع نہیں ہوتا۔ اس میں شامل اکاؤنٹس کے لیے اس کے مکمل ہونے تک صرف غائب لنکس بنائے جاتے ہیں۔ اگر اسے مکمل نہ کیا جا سکے تو غلطی بتاتی ہے کہ وہ فائل ہٹانے سے تبدیلی کا باقی حصہ چھوڑ دیا جاتا ہے۔ `ccam --dry-run doctor --fix` ترتیب (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) اور یہ دکھاتا ہے کہ ہر اکاؤنٹ کا حصہ کیا الگ رکھ دے گا۔ جو تبدیلی اس وقت کوئی دوسرا ccam کر رہا ہو وہ معلومات ہے، اور جب تک وہ چل رہی ہو `--fix` کوئی شیئر لنک ٹھیک نہیں کرتا (`its share links are not repaired while another ccam changes the groups; …`)۔ مرمت کے دوران `--fix`، `~/.local/state/ccam/group-change.lock` تھامے رکھتا ہے، اور کوئی تبدیلی مکمل کرنے کے بعد `config.toml` دوبارہ پڑھتا ہے۔

### شیئر لنکس کی مرمت

لانچ کا راستہ غائب لنک فوراً بنا دیتا ہے، لیکن جس چیز سے ڈیٹا کھو سکتا ہو اس پر رک جاتا ہے اور یہاں کی طرف رہنمائی کرتا ہے:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

شیئر لنک کے کسی بھی مسئلے کے لیے `--fix`، اکاؤنٹ کی شیئر سیٹنگز کو موجودہ اور مطلوبہ دونوں حالتوں کے طور پر دے کر `share.Build` کو کال کرتا ہے، پھر `Plan.Apply` کو۔ `Build` ہر قدم ڈسک پر موجود چیزوں سے دوبارہ اخذ کرتا ہے، وہی منصوبہ جو `ccam group` کی کوئی تبدیلی استعمال کرتی ہے، اس لیے مرمت انہی اصولوں پر چلتی ہے جن پر کسی گروپ میں شامل ہونا۔

<!-- motion:doctor-repair -->

یہاں `settings.json` ایک اصل فائل ہے۔ چونکہ گروپ کے پاس پہلے سے ایک مختلف `settings.json` ہے، ٹکراؤ اپنا ڈیفالٹ انتخاب `KeepGroup` لیتا ہے، اس لیے گروپ کی کاپی برقرار رہتی ہے۔ اکاؤنٹ کی فائل اکاؤنٹ فولڈر کے اندر `.ccam-backup/<time>/settings.json` میں منتقل ہو جاتی ہے، اور `os.Symlink` اس راستے کو گروپ سے لنک کر دیتا ہے۔ `--fix` یہ بتا دیتا ہے: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`۔ دوسری صورتیں:

- اگر گروپ کے پاس ابھی کوئی کاپی نہ ہو تو اکاؤنٹ کی فائل گروپ میں منتقل ہو کر اس کا پہلا مواد بن جاتی ہے۔
- فولڈر کے لیے، جو فائلیں گروپ میں نہیں ہیں وہ گروپ میں منتقل ہو جاتی ہیں۔ جو فائلیں گروپ کے پاس بائٹ بہ بائٹ ویسی ہی موجود ہوں وہ ہٹا دی جاتی ہیں، اور جو مختلف ہوں وہ بیک اپ میں جاتی ہیں، جن کے نام ایک نوٹ میں بتائے جاتے ہیں (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`)۔
- غائب گروپ فولڈر پہلے بنایا جاتا ہے۔
- `history.jsonl` کے لیے ڈیفالٹ `merge by line` ہے: پہلے گروپ کی لائنیں، پھر آپ کی نئی لائنیں۔ گروپ کی پچھلی کاپی `.ccam-backup/<time>/group/` میں محفوظ ہوتی ہے۔
- کسی دوسرے گروپ کا لنک بدل دیا جاتا ہے۔ ccam سے باہر کسی فولڈر کا لنک بیک اپ میں جاتا ہے، اور اگر گروپ خالی ہو تو اس کا مواد گروپ کا ابتدائی مواد بن جاتا ہے۔
- کسی ایسے آئٹم کا بچا ہوا لنک جو اب لوکل ہے ہٹا دیا جاتا ہے، اور اکاؤنٹ کو گروپ کے مواد کی ایک کاپی مل جاتی ہے۔

لنکس ہٹانے اور ان فائلوں کو ہٹانے کے سوا جو گروپ کے پاس پہلے سے بائٹ بہ بائٹ موجود ہیں، ہر قدم ڈیٹا کو گروپ یا بیک اپ فولڈر میں منتقل کرتا ہے۔ کچھ بھی ضائع نہیں ہوتا۔

## حفاظت اور خاص صورتیں

- سادہ `ccam doctor` کبھی کچھ تبدیل نہیں کرتا۔ اسکرپٹس میں یقینی بنانے کے لیے `--fix` کے ساتھ `--dry-run` استعمال کریں۔
- جو اکاؤنٹ چل رہا ہو، چاہے Claude Code کسی بھی طرح شروع ہوا ہو، اس کے لیے `doctor --fix` صرف غائب لنکس اور غائب گروپ فولڈرز بناتا ہے، جس سے کچھ ضائع نہیں ہوتا۔ جب کچھ اور باقی رہ جائے تو یہ `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` شامل کرتا ہے، کیونکہ سیشن ان فائلوں میں لکھتا رہتا ہے جنہیں مرمت منتقل کرتی۔
- اگر شیئر کی مرمت ناکام ہو تو ccam اس اکاؤنٹ کے لیے `could not repair share links: <error>` بتاتا ہے اور اس کی نامعلوم اینٹریز کی جانچ چھوڑ دیتا ہے۔ دوسرے اکاؤنٹس پھر بھی جانچے جاتے ہیں۔
- اگر `ccam group` کی کوئی تبدیلی آدھے راستے میں ناکام ہو تو نئے گروپس پہلے ہی محفوظ ہو چکے ہوتے ہیں اور تبدیلی `~/.local/state/ccam/group-change.json` میں درج ہوتی ہے؛ ccam بتاتا ہے کہ `ccam doctor --fix` سے مکمل کریں، جو اس کا باقی حصہ، MCP کے مراحل سمیت، پورا کرتا ہے۔
- symlink والی `.credentials.json`، میعاد ختم ہونے والا لاگ اِن اور انوائرنمنٹ کی اسناد بتائی جاتی ہیں لیکن کبھی تبدیل نہیں کی جاتیں۔
- پلگ اِن رجسٹری کی دوبارہ تحریریں ایٹامک طور پر لکھی جاتی ہیں لیکن `.ccam-backup/` میں کاپی نہیں ہوتیں۔

## متعلقہ

- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [موجودہ سیٹ اپ کاپی کرنا](/features/migrate)
- [استعمال کی حدیں اور اکاؤنٹ تبدیل کرنا](/docs/usage-limits)
- [کمانڈز](/docs/commands)
