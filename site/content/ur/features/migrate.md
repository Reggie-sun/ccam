---
title: "موجودہ سیٹ اپ کاپی کرنا"
description: "ccam migrate سے موجودہ ~/.claude کی سیٹنگز، CLAUDE.md، ایکسٹینشنز، MCP سرورز اور MCP سائن اِن، ccam کے ایک گروپ میں کاپی کریں۔"
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

اگر آپ ccam سے پہلے Claude Code استعمال کرتے تھے تو آپ کا سیٹ اپ `~/.claude` اور `~/.claude.json` میں ہے۔ `ccam migrate` اسے ایک شیئر گروپ میں کاپی کرتا ہے، تاکہ جو بھی اکاؤنٹ گروپ کے ذریعے وہ آئٹمز شیئر کرے وہ آپ کے موجودہ سیٹ اپ کے ساتھ شروع ہو۔ سورس فولڈر کو صرف پڑھا جاتا ہے، اور سادہ `claude` پہلے کی طرح کام کرتا رہتا ہے۔ اسے اپنے پہلے `ccam add` کے بعد ایک بار چلائیں، اور پھر جب بھی کوئی اکاؤنٹ MCP سرورز شیئر کرنا شروع کرے۔

## استعمال

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

ٹرمینل پر، آئٹم فلیگز اور `--yes` کے بغیر، ccam پہلے ایک چیک لسٹ دکھاتا ہے:

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

Space (یا `x`) کرسر کے نیچے والا آئٹم بدلتا ہے، تیر والی کیز یا `j`/`k` حرکت دیتی ہیں، Enter منصوبہ دکھاتا ہے، اور Esc، `q` یا Ctrl+C، `Cancelled.` پرنٹ کرتے ہیں۔ جن آئٹمز کے لیے سورس میں کچھ نہیں ہوتا ان پر `not found` لکھا ہوتا ہے۔ اگر آپ سب کچھ ان چیک کر دیں تو ccam، `Nothing selected.` پرنٹ کر کے رک جاتا ہے۔

پھر ccam منصوبہ پرنٹ کرتا ہے اور کاپی کرنے سے پہلے پوچھتا ہے:

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

گروپ صرف اپنے اکاؤنٹس تک پہنچتا ہے، اور صرف ان آئٹمز کے لیے جو وہ شیئر کرتا ہے، اس لیے خلاصہ گروپ کے اکاؤنٹس کی فہرست دیتا ہے۔ جب گروپ کاپی کیے گئے آئٹمز میں سے کوئی شیئر نہ کرتا ہو، یا کچھ اکاؤنٹس اس میں نہ ہوں، تو ccam وہ `ccam group set <group> --add <items>` اور `ccam group join <group> <accounts>` کمانڈز پرنٹ کرتا ہے جو اسے بدلتی ہیں۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `[folder]` | وہ Claude Code فولڈر جس سے کاپی کرنا ہے | `~/.claude` |
| `--settings` | `settings.json` کاپی کرتا ہے | بند |
| `--memory` | گلوبل `CLAUDE.md` کاپی کرتا ہے | بند |
| `--extensions` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/`، `rules/` کاپی کرتا ہے | بند |
| `--mcp` | یوزر MCP سرورز (`.claude.json` میں `mcpServers`) کاپی کرتا ہے | بند |
| `--history` | گفتگو کی ہسٹری (`projects/`، `history.jsonl`) کاپی کرتا ہے | بند |
| `-g`، `--group string` | وہ شیئر گروپ جس میں کاپی کرنا ہے | `default` |
| `-y`، `--yes` | پوچھے بغیر کاپی کرتا ہے؛ ٹکراؤ پر گروپ کی کاپی رکھی جاتی ہے | بند |
| `--dry-run` (گلوبل) | منصوبہ پرنٹ کرتا ہے اور کسی سوال یا کاپی سے پہلے رک جاتا ہے | بند |

کوئی آئٹم فلیگ نہ ہو تو ٹرمینل پر چیک لسٹ فیصلہ کرتی ہے۔ ٹرمینل کے باہر، یا `--yes` کے ساتھ، ccam گفتگو کی ہسٹری کے سوا سب کچھ کاپی کرتا ہے۔ ہسٹری ڈیفالٹ طور پر چھوڑ دی جاتی ہے کیونکہ `projects/` بڑا ہو سکتا ہے۔

`--group on` کا مطلب `default` ہے۔ `--group off` اور `--group local`، `--group takes a group name, not "off"` کے ساتھ مسترد ہوتے ہیں۔ `ccam migrate`، `--dry-run` کو لاگو کرتا ہے، اور ڈرائی رن میں بھی چیک لسٹ ظاہر ہوتی ہے تاکہ آپ کسی انتخاب کا پیشگی جائزہ لے سکیں۔

## یہ کیسے کام کرتا ہے

`internal/cli/migrate.go` آئٹمز چنتا ہے، پھر `migrate.Build` سورس اور گروپ کو پڑھ کر ایک منصوبہ لوٹاتا ہے۔ آپ کی تصدیق کے بعد `Plan.Apply` تک کچھ نہیں لکھا جاتا۔

1. **سورس تلاش کرنا۔** `checkSource` ایسا راستہ مسترد کرتا ہے جو موجود نہ ہو، فولڈر نہ ہو، یا جس میں Claude Code کی کوئی معلوم اینٹری نہ ہو۔ `migrate.ClaudeJSON` وہ `.claude.json` تلاش کرتا ہے جو فولڈر سے تعلق رکھتی ہے: اس کے اندر، یا `~/.claude` کے لیے جس کے اندر کوئی نہ ہو، اس کے ساتھ والی `~/.claude.json`۔ جب `CLAUDE_CONFIG_DIR` مقرر نہ ہو تو Claude Code اسے وہیں رکھتا ہے۔
2. **آئٹمز کو راستوں سے جوڑنا۔** ہر راستہ `share.Layout` سے آتا ہے، وہی اجازت شدہ فہرست جو گروپس استعمال کرتے ہیں۔ MCP سرورز `.claude.json` میں سے `mcp.LocalServers` سے پڑھے جاتے ہیں، جو پوری فائل پارس کرنے کے بجائے `profile.TopLevelValue` اسکینر استعمال کرتا ہے۔
3. **کاپی کرنا، کبھی منتقل نہیں۔** فائلیں `fsutil.AtomicWrite` یا `fsutil.CopyFile` سے `~/.local/share/ccam/groups/<group>/` میں جاتی ہیں؛ MCP سرورز گروپ کی `mcp.json` میں جاتے ہیں۔

ہر آئٹم ایک سورس راستے کو ایک گروپ راستے سے جوڑتا ہے:

<!-- motion:migrate-items -->

ہسٹری کی دو قطاریں ڈیش والی ہیں کیونکہ وہ صرف `--history` کے ساتھ کاپی ہوتی ہیں۔ `.credentials.json` اور `.claude.json`، `share.Never` میں درج ہیں اور کبھی گروپ تک نہیں پہنچتیں: اگر ان میں سے کوئی کبھی ترتیب میں نظر آئے تو `Build` اندرونی غلطی لوٹاتا ہے، اور فولڈر کی تلاش انہیں نام سے چھوڑ دیتی ہے۔

`buildEntry` ہر راستے کو ایک حالت دیتا ہے، اور منصوبے کی ہر لائن یہی بیان کرتی ہے:

<!-- motion:migrate-states -->

- **Missing**: سورس میں نہیں، یا جہاں فولڈر متوقع ہو وہاں فائل۔ ہر آئٹم کے لیے ایک بار `not found, skipped` کے طور پر دکھایا جاتا ہے۔
- **Linked**: سورس راستہ پہلے سے گروپ کی کاپی تک پہنچتا ہے، مثلاً جب سورس ccam کا کوئی پروفائل فولڈر ہو۔
- **Same**: گروپ میں پہلے سے بالکل وہی مواد ہے۔
- **New**: گروپ میں یہ نہیں ہے، اس لیے یہ پورا کاپی ہوتا ہے۔
- **Merge**: دونوں طرف ایک فولڈر۔ `copyMissing` صرف وہ فائلیں شامل کرتا ہے جو گروپ میں نہیں ہیں؛ مختلف فائلیں ویسی ہی رہتی ہیں جیسی گروپ میں ہیں اور منصوبہ `group copies kept` کہتا ہے۔ `plugins/` کے لیے، رجسٹریز `installed_plugins.json` اور `known_marketplaces.json` ملا دی جاتی ہیں، اور پلگ اِنز کے مکمل راستے گروپ کے فولڈر کی طرف دوبارہ لکھے جاتے ہیں۔
- **Conflict**: دونوں طرف ایک فائل جس کا مواد مختلف ہو۔ آپ ہر فائل کے لیے چنتے ہیں: `keep group`، `use mine`، اور `.json` فائلوں کے لیے `merge JSON` (گہرا انضمام جس میں ٹکراتی قدروں پر گروپ کو ترجیح ملتی ہے) یا `history.jsonl` کے لیے `merge by line`۔ جس MCP سرور کی تعریف مختلف ہو وہ بھی ٹکراؤ ہے، `keep group` یا `use mine` کے ساتھ۔ پہلا دکھایا گیا انتخاب Enter دبانے پر ڈیفالٹ ہے؛ `--yes` ہمیشہ گروپ کی کاپی رکھتا ہے۔

چیک لسٹ خود `tui.MigratePicker` ہے، جو اسی `checklist` جزو پر بنی ہے جو `ccam add` استعمال کرتا ہے؛ `migrate.Found` غائب آئٹمز کو lstat کالز اور MCP سرورز کی فہرست سے نشان زد کرتا ہے۔

### MCP سائن اِن

جو MCP سرورز OAuth سے سائن اِن کرتے ہیں وہ اپنے ٹوکنز سرورز کی فہرست میں نہیں بلکہ لاگ اِن فائل میں رکھتے ہیں، یعنی `.credentials.json` میں `mcpOAuth` کے تحت۔ Claude Code ہر اینٹری کی key، `<name>|` اور سرور کے `type`، `url` اور `headers` کے `sha256` hash کے پہلے 16 حروف سے بناتا ہے۔ key فولڈر پر منحصر نہیں، اس لیے کاپی کی گئی اینٹری دوسرے اکاؤنٹ میں تب تک کام کرتی ہے جب تک سرور کی تعریف بالکل وہی ہو۔ اسی لیے `internal/migrate/logins.go` میں `buildLogins` صرف ان سرورز کی اینٹریز کاپی کرتا ہے جو آخرکار گروپ میں سورس والی تعریف کے ساتھ ہوں: نئے، پہلے سے یکساں، یا ایسا ٹکراؤ جس میں آپ نے `use mine` چنا ہو۔ پلگ اِن سرورز کے سائن اِن (جن کی keys `plugin:` سے شروع ہوتی ہیں) تب ساتھ جاتے ہیں جب آپ ایکسٹینشنز کاپی کریں اور سورس میں `plugins/` ہو۔

<!-- motion:migrate-signins -->

اینٹریز ہر وصول کرنے والے اکاؤنٹ کی اپنی `.credentials.json` میں جاتی ہیں، جو `fsutil.AtomicWrite` سے موڈ `0600` پر لکھی جاتی ہے۔ یہ کبھی گروپ میں نہیں جاتیں۔ کسی اکاؤنٹ کو یہ صرف تب ملتی ہیں جب وہ ہدف گروپ کے ذریعے `mcp` (یا پلگ اِن سرورز کے لیے `extensions`) شیئر کرتا ہو، اور صرف وہ keys شامل ہوتی ہیں جو اس کے پاس نہیں: جو اینٹریز اس کے پاس پہلے سے ہیں، شاید زیادہ نئی، وہ برقرار رہتی ہیں۔ منصوبے کی `sign-ins` لائن سرورز، وصول کرنے والے اکاؤنٹس اور چھوڑے گئے اکاؤنٹس دکھاتی ہے۔ ccam کسی اکاؤنٹ کو چھوڑ دیتا ہے، اور وجہ کے ساتھ اس کا نام بتاتا ہے، جب وہ:

- `ccam import` سے اپنی جگہ پر رجسٹر کیا گیا ہو، کیونکہ سادہ `claude` اس کی لاگ اِن فائل استعمال کر رہا ہو سکتا ہے؛
- چل رہا ہو، کیونکہ Claude Code اس کی لاگ اِن فائل دوبارہ لکھتا ہے (لکھنے سے فوراً پہلے دوبارہ جانچا جاتا ہے، اور `started meanwhile` کے طور پر بتایا جاتا ہے)؛
- ابھی لاگ اِن نہ ہوا ہو؛
- اپنا لاگ اِن macOS Keychain میں رکھتا ہو، جس میں ccam نہیں لکھتا۔

چھوڑے گئے اکاؤنٹ میں سرورز میں `/mcp` سے سائن اِن کریں۔ جو اکاؤنٹ بعد میں MCP سرورز شیئر کرنا شروع کرے اسے سائن اِن تب ملتے ہیں جب آپ `ccam migrate` دوبارہ چلائیں؛ خلاصہ یہ ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` کے ساتھ بتاتا ہے۔

## حفاظت اور خاص صورتیں

- سورس کو صرف پڑھا جاتا ہے۔ اس کی کوئی فائل منتقل، لنک یا تبدیل نہیں ہوتی، اور سورس کا لاگ اِن کبھی نہیں لکھا جاتا، macOS پر بھی جہاں اسے Keychain سے پڑھا جا سکتا ہے۔
- گروپ کی کوئی فائل بدلنے سے پہلے، چاہے `use mine` سے، انضمام سے، یا موجودہ `mcp.json` میں نئے MCP سرورز سے، گروپ کی کاپی `<group>/.ccam-backup/<time>/` میں محفوظ ہوتی ہے۔ فولڈر کا نام ہر رن کے لیے منفرد ہوتا ہے، اس لیے ایک ہی سیکنڈ میں دو رن ایک دوسرے کے بیک اپ اوور رائٹ نہیں کرتے۔ خلاصہ راستہ پرنٹ کرتا ہے۔
- `--yes` کے ساتھ ہر ٹکراؤ میں گروپ کی کاپی رکھی جاتی ہے، اور ccam بتاتا ہے کہ اس کے بجائے خود کیسے چنیں: وہی کمانڈ `--yes` کے بغیر چلائیں۔
- اگر کچھ نہ بدلے تو ccam، `Nothing to copy.` پرنٹ کرتا ہے؛ اگر ہر ٹکراؤ میں گروپ کی کاپی رکھی جائے تو `Nothing to copy: the group keeps its own copies.`
- اگر کاپی آدھے راستے میں ناکام ہو تو ccam، `The source folder was not changed. Group files replaced so far are saved in` اور اس کے بعد بیک اپ فولڈر کے ساتھ رک جاتا ہے۔
- سورس نہ ہو تو عمل `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` پر ختم ہوتا ہے۔

## متعلقہ

- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [جانچ اور مرمت](/features/doctor)
- [Claude Code کے ایک سے زیادہ اکاؤنٹس](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
