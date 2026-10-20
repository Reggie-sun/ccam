---
title: "جانچ اور مرمت"
description: "ccam doctor سے شیئر لنکس، فائلوں کی اجازتیں، لاگ اِنز اور Claude Code کا ورژن جانچیں، اور ccam doctor --fix سے وہ ٹھیک کریں جسے ٹھیک کرنا محفوظ ہو۔"
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` آپ کی ccam انسٹالیشن جانچتا ہے: Claude Code بائنری، ہر اکاؤنٹ کا فولڈر، لاگ اِن اور شیئر لنکس، شیئر کیے گئے MCP سرورز، usage اینڈ پوائنٹ اور آپ کی شارٹ کٹ کمانڈز۔ `ccam doctor --fix` وہ ٹھیک کرتا ہے جو ڈیٹا کھوئے بغیر ہو سکے، اور جو فائل بھی وہ بدلتا ہے اسے پہلے ایک بیک اپ فولڈر میں منتقل کرتا ہے۔ اسے تب چلائیں جب `ccam <name>` کسی شیئر لنک کے مسئلے پر رک جائے، جب آپ نے فائلیں ہاتھ سے بحال کی ہوں، یا Claude Code کی اپ ڈیٹ کے بعد۔

## استعمال

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

ہر نتیجہ ایک لائن ہے جس پر ایک نشان ہوتا ہے: `✓` ٹھیک، `·` معلومات، `!` انتباہ، `✗` غلطی۔ اکاؤنٹ کے نتائج اکاؤنٹ کے نام سے شروع ہوتے ہیں۔

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

`--fix` کے ساتھ، جو بھی نتیجہ ٹھیک کیا گیا ہو اس کے آخر میں `→ fixed` آتا ہے:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--fix` | وہ ٹھیک کرتا ہے جسے ٹھیک کرنا محفوظ ہو (بدلی گئی فائلوں کا بیک اپ لیا جاتا ہے) | بند |
| `--dry-run` (گلوبل) | ہر جانچ چلاتا ہے لیکن کچھ ٹھیک نہیں کرتا، `--fix` کے ساتھ بھی | بند |

`ccam doctor` کوئی آرگیومنٹ نہیں لیتا۔ یہ `--dry-run` کو لاگو کرتا ہے: `ccam --dry-run doctor --fix` وہی نتائج پرنٹ کرتا ہے جو `ccam doctor` کرتا ہے اور کچھ تبدیل نہیں کرتا۔

جب کم از کم ایک `✗` نتیجہ ٹھیک کیے بغیر رہ جائے تو ایگزٹ کوڈ `1` ہوتا ہے، ورنہ `0`۔ انتباہات اور معلومات کبھی رن کو ناکام نہیں کرتیں۔ جب رن `--fix` کے بغیر ناکام ہو تو ccam ایک خالی لائن اور `` `ccam doctor --fix` repairs what it can. `` پر ختم ہوتا ہے۔ جو کنفیگ لوڈ نہ ہو سکے وہ بھی `1` اور ایک `ccam:` غلطی پر ختم ہوتا ہے۔

## یہ کیسے کام کرتا ہے

`internal/cli/doctor.go`، `~/.config/ccam/config.toml` لوڈ کرتا ہے اور `internal/doctor/doctor.go` میں `doctor.Run` کو کال کرتا ہے، جو جانچیں ایک مقررہ ترتیب میں چلاتا ہے اور نتائج لوٹاتا ہے۔ `--fix` آگے `Options.Fix` کے طور پر جاتا ہے، جو ڈرائی رن میں false ہوتا ہے۔

<!-- motion:doctor-checks -->

درمیانی حصہ ہر اکاؤنٹ کے لیے ایک بار چلتا ہے۔ جانچیں، اور ہر ایک کے لیے `--fix` کیا کرتا ہے:

1. **Claude Code کا ورژن** (`checkClaude`)۔ `claude` بائنری کو اسی طرح تلاش کرتا ہے جیسے لانچر کرتا ہے اور `claude --version` چلاتا ہے۔ اگر ورژن `2.1.302` سے شروع نہ ہو، جو وہ ورژن ہے جس کے مقابلے میں ccam کے فولڈر کی ترتیب جانچی گئی تھی (`doctor.LayoutVersion`)، تو یہ `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` شامل کرتا ہے۔ بائنری کا نہ ہونا یا `--version` کا ناکام ہونا غلطی ہے۔ کوئی مرمت نہیں۔
2. **انوائرنمنٹ میں اسناد** (`checkAuthEnv`)۔ جب `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` مقرر ہو تو انتباہ دیتا ہے: اس شیل سے لانچ ہونے والا ہر اکاؤنٹ اسے وراثت میں پاتا ہے، اور Claude Code اسے اکاؤنٹ کے اپنے لاگ اِن پر ترجیح دیتا ہے۔ کوئی مرمت نہیں؛ اسے خود ہٹائیں۔
3. **فولڈر** (`checkProfile`)۔ اکاؤنٹ فولڈر کا نہ ہونا غلطی ہے اور اس اکاؤنٹ کی جانچ ختم کر دیتا ہے۔ اگر اجازتوں میں group یا other کا کوئی بٹ مقرر ہو تو انتباہ ہے؛ `--fix`، `chmod 0700` چلاتا ہے۔ Windows میں اجازت کے بٹس نہیں ہوتے، اس لیے وہاں یہ جانچ اور نیچے والی `0600` جانچ چھوڑ دی جاتی ہے؛ آپ کے یوزر پروفائل میں موجود فولڈر اپنی ایکسیس لسٹ سے محفوظ رہتا ہے۔
4. **لاگ اِن فائل** (`checkCredentials`)۔ `.credentials.json` کو `lstat` سے جانچا جاتا ہے۔ symlink غلطی ہے، کیونکہ ہر اکاؤنٹ کا اپنا ٹوکن ہونا چاہیے، اور اسے ٹھیک نہیں کیا جاتا۔ `0600` کے سوا کوئی موڈ انتباہ ہے؛ `--fix`، `chmod 0600` چلاتا ہے۔ پھر لاگ اِن پڑھا جاتا ہے (macOS پر پہلے Keychain سے)۔ لاگ اِن نہ ہونا انتباہ ہے جو `ccam login <name>` کی طرف رہنمائی کرتا ہے۔ میعاد ختم ہونے والا ریفریش ٹوکن `login expired` غلطی ہے؛ میعاد ختم ہونے والا ایکسیس ٹوکن صرف معلومات ہے، کیونکہ Claude Code اگلے لانچ پر اس کی تجدید کر دیتا ہے۔ ان میں سے کوئی بھی ٹھیک نہیں کیا جاتا۔
5. **شیئر لنکس** (`share.Check`)۔ `share.Layout` کے ہر راستے کو اکاؤنٹ کی شیئر سیٹنگز سے موازنہ کرتا ہے، صرف `lstat` اور `readlink` استعمال کرتے ہوئے۔ یہ چار قسمیں بتاتا ہے: غائب لنک، کسی اور جگہ اشارہ کرنے والا لنک، جہاں لنک ہونا چاہیے وہاں اصل فائل، اور کسی ایسے آئٹم کا بچا ہوا لنک جو اب لوکل ہے۔ غائب لنک انتباہ ہے، باقی غلطیاں ہیں۔
6. **نامعلوم اینٹریز** (`checkUnknownEntries`)۔ اکاؤنٹ فولڈر میں اوپری سطح کے وہ نام دکھاتا ہے جو ccam نہیں جانتا (`share.Known`)، اس اشارے کے طور پر کہ Claude Code کی ترتیب بدل گئی ہے۔ صرف معلومات۔
7. **پلگ اِن کے راستے** (`checkPlugins`)۔ پلگ اِن رجسٹریز (`installed_plugins.json`، `known_marketplaces.json`) مکمل راستے محفوظ کرتی ہیں۔ ہر گروپ کے `plugins/` میں، جو راستے کسی ایک اکاؤنٹ کے فولڈر سے گزرتے ہیں وہ اس اکاؤنٹ کے ہٹائے جانے یا نام بدلنے پر ٹوٹ جاتے ہیں۔ `--fix` انہیں `profile.RewritePluginPaths` سے گروپ فولڈر پر دوبارہ لکھتا ہے، جو رجسٹری فائل کی ایٹامک تحریر ہے۔
8. **MCP سرورز** (`checkMCP`)۔ ہر گروپ کی `mcp.json` ایک بار لوڈ کرتا ہے۔ جو فائل پڑھی نہ جا سکے وہ غلطی ہے؛ جو سرور `mcp.Validate` میں ناکام ہو وہ انتباہ ہے۔ کوئی مرمت نہیں۔
9. **usage جواب کا فارمیٹ** (`checkUsageFormat`)۔ کیش کے بغیر اور معمول کی 3 سیکنڈ کی حد کے ساتھ، usage اینڈ پوائنٹ سے ایک کے بعد ایک اکاؤنٹ کے لیے پوچھتا ہے یہاں تک کہ کوئی جواب دے۔ اگر جواب اب پارس نہ ہو تو یہ انتباہ دیتا ہے کہ `ccam usage` اسٹیٹس لائن کیش پر چلا جاتا ہے۔ یہ واحد جانچ ہے جو نیٹ ورک استعمال کرتی ہے۔
10. **شارٹ کٹس** (`checkShims`)۔ صرف تب جب کسی اکاؤنٹ کا شارٹ کٹ ہو۔ `~/.local/bin` میں غائب لنکس انتباہ ہیں؛ `--fix`، `shim.Install` کو کال کرتا ہے، جو صرف غائب لنکس بناتا ہے اور دوسری فائلوں کو کبھی نہیں چھیڑتا۔ جو bin فولڈر `PATH` پر نہ ہو وہ انتباہ ہے اور ٹھیک نہیں کیا جاتا۔
11. **فولڈرز** (`checkFolders`)۔ صرف تب جب [`ccam folder`](/features/folder) اندراجات موجود ہوں۔ جس اندراج کا اکاؤنٹ موجود نہ ہو وہ غلطی ہے، اور جس اندراج کا فولڈر ختم ہو چکا ہو وہ انتباہ ہے۔ دونوں میں سے کوئی ٹھیک نہیں کیا جاتا، کیونکہ اندراج ہٹانے سے بدل جائے گا کہ وہاں کون سا اکاؤنٹ شروع ہوتا ہے۔

### شیئر لنکس کی مرمت

لانچ کا راستہ غائب لنک فوراً بنا دیتا ہے، لیکن جس چیز سے ڈیٹا کھو سکتا ہو اس پر رک جاتا ہے اور یہاں کی طرف رہنمائی کرتا ہے:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

شیئر لنک کے کسی بھی مسئلے کے لیے `--fix`، اکاؤنٹ کی شیئر سیٹنگز کو موجودہ اور مطلوبہ دونوں حالتوں کے طور پر دے کر `share.Build` کو کال کرتا ہے، پھر `Plan.Apply` کو۔ `Build` ہر قدم ڈسک پر موجود چیزوں سے دوبارہ اخذ کرتا ہے، وہی منصوبہ جو `ccam share` استعمال کرتا ہے، اس لیے مرمت انہی اصولوں پر چلتی ہے جن پر شیئرنگ آن کرنا۔

<!-- motion:doctor-repair -->

یہاں `settings.json` ایک اصل فائل ہے۔ چونکہ گروپ کے پاس پہلے سے ایک مختلف `settings.json` ہے، ٹکراؤ اپنا ڈیفالٹ انتخاب `KeepGroup` لیتا ہے، اس لیے گروپ کی کاپی برقرار رہتی ہے۔ اکاؤنٹ کی فائل اکاؤنٹ فولڈر کے اندر `.ccam-backup/<time>/settings.json` میں منتقل ہو جاتی ہے، اور `os.Symlink` اس راستے کو گروپ سے لنک کر دیتا ہے۔ دوسری صورتیں:

- اگر گروپ کے پاس ابھی کوئی کاپی نہ ہو تو اکاؤنٹ کی فائل گروپ میں منتقل ہو کر اس کا پہلا مواد بن جاتی ہے۔
- فولڈر کے لیے، جو فائلیں گروپ میں نہیں ہیں وہ گروپ میں منتقل ہو جاتی ہیں؛ باقی بیک اپ میں جاتی ہیں۔
- `history.jsonl` کے لیے ڈیفالٹ `merge by line` ہے: پہلے گروپ کی لائنیں، پھر آپ کی نئی لائنیں۔ گروپ کی پچھلی کاپی `.ccam-backup/<time>/group/` میں محفوظ ہوتی ہے۔
- کسی دوسرے گروپ کا لنک بدل دیا جاتا ہے۔ ccam سے باہر کسی فولڈر کا لنک بیک اپ میں جاتا ہے، اور اگر گروپ خالی ہو تو اس کا مواد گروپ کا ابتدائی مواد بن جاتا ہے۔
- کسی ایسے آئٹم کا بچا ہوا لنک جو اب لوکل ہے ہٹا دیا جاتا ہے، اور اکاؤنٹ کو گروپ کے مواد کی ایک کاپی مل جاتی ہے۔

لنکس ہٹانے کے سوا، ہر قدم ڈیٹا کو گروپ یا بیک اپ فولڈر میں منتقل کرتا ہے۔ کچھ بھی حذف نہیں ہوتا۔

## حفاظت اور خاص صورتیں

- سادہ `ccam doctor` کبھی کچھ تبدیل نہیں کرتا۔ اسکرپٹس میں یقینی بنانے کے لیے `--fix` کے ساتھ `--dry-run` استعمال کریں۔
- `doctor --fix` یہ نہیں جانچتا کہ کوئی اکاؤنٹ چل رہا ہے یا نہیں۔ کسی اکاؤنٹ کے شیئر لنکس کی مرمت سے پہلے اس کا Claude Code بند کریں، کیونکہ سیشن ان فائلوں میں لکھتا رہتا ہے۔
- اگر شیئر کی مرمت ناکام ہو تو ccam اس اکاؤنٹ کے لیے `could not repair share links: <error>` بتاتا ہے اور اس کی نامعلوم اینٹریز کی جانچ چھوڑ دیتا ہے۔ دوسرے اکاؤنٹس پھر بھی جانچے جاتے ہیں۔
- اگر `ccam share` آدھے راستے میں ناکام ہو تو یہ پہلے مطلوبہ سیٹنگز محفوظ کرتا ہے اور بتاتا ہے کہ `ccam doctor --fix` سے مکمل کریں۔
- symlink والی `.credentials.json`، میعاد ختم ہونے والا لاگ اِن اور انوائرنمنٹ کی اسناد بتائی جاتی ہیں لیکن کبھی تبدیل نہیں کی جاتیں۔
- پلگ اِن رجسٹری کی دوبارہ تحریریں ایٹامک طور پر لکھی جاتی ہیں لیکن `.ccam-backup/` میں کاپی نہیں ہوتیں۔

## متعلقہ

- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [موجودہ سیٹ اپ کاپی کرنا](/features/migrate)
- [استعمال کی حدیں اور اکاؤنٹ تبدیل کرنا](/docs/usage-limits)
- [کمانڈز](/docs/commands)
