---
title: "گروپس کے ساتھ شیئرنگ"
description: "ccam کا گروپ چنے ہوئے آئٹمز اپنے اکاؤنٹس کے درمیان ایک فولڈر سے شیئر کرتا ہے؛ ccam group ڈیٹا کو منصوبے، بیک اپ اور ٹکراؤ کے انتخاب کے ساتھ منتقل کرتا ہے۔"
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "e4f05b0fd572"
---

ہر اکاؤنٹ کا اپنا فولڈر ہوتا ہے، لیکن اس میں موجود زیادہ تر چیزوں کو الگ رکھنے کی ضرورت نہیں۔ ایک **گروپ** بتاتا ہے کہ کیا شیئر ہوتا ہے اور کن کے درمیان: گروپ کا ہر اکاؤنٹ ان آئٹمز کو ایک فولڈر، `~/.local/share/ccam/groups/<group>/`، سے لنک کرتا ہے، اس لیے وہ سب ایک ہی کاپی پڑھتے اور لکھتے ہیں۔ لاگ اِن ٹوکنز کبھی شیئر نہیں ہوتے۔

ایک اکاؤنٹ کئی گروپس میں ہو سکتا ہے، لیکن ہر آئٹم صرف ایک گروپ سے آتا ہے۔ ایک عام سیٹ اپ میں سیٹنگز اور ایکسٹینشنز کے لیے ایک گروپ ہوتا ہے، اور دو اکاؤنٹس کے درمیان گفتگو کی ہسٹری کے لیے دوسرا:

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

آپ کا شامل کیا گیا پہلا اکاؤنٹ `default` بناتا ہے، اور بعد کے اکاؤنٹس اس میں شامل ہو جاتے ہیں جب تک آپ کچھ اور نہ کہیں۔

## استعمال

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
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

ڈیٹا منتقل کرنے والی ہر کمانڈ ایک منصوبہ پرنٹ کرتی ہے، ہر اس فائل کے بارے میں پوچھتی ہے جو گروپ کی کاپی سے مختلف ہو، اور کچھ بھی تبدیل کرنے سے پہلے `Proceed? [y/N]` پوچھتی ہے:

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

یہی تبدیلیاں [ڈیش بورڈ](/features/dashboard) کے **Groups** ٹیب پر بھی ہیں، اور کسی اکاؤنٹ پر `s` اس کے گروپس چنتا ہے۔

## آئٹمز

| آئٹم | کیا شامل ہے |
| --- | --- |
| `mcp` | MCP سرورز، گروپ کی `mcp.json` کے ذریعے ([شیئر کیے گئے MCP سرورز](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | گلوبل `CLAUDE.md` |
| `extensions` | `agents/`، `commands/`، `skills/`، `plugins/`، `output-styles/`، `rules/` |
| `history` | `projects/` اور `history.jsonl`، تاکہ کوئی دوسرا اکاؤنٹ `--continue` کر سکے |

جو آئٹم اکاؤنٹ کا کوئی گروپ شیئر نہیں کرتا وہ اکاؤنٹ کے اپنے فولڈر میں رہتا ہے۔

## کمانڈز اور آپشنز

| کمانڈ | کیا کرتی ہے |
| --- | --- |
| `ccam group` یا `ccam group list` | گروپس کی فہرست دکھاتی ہے |
| `ccam group show <group>` | گروپ کیا شیئر کرتا ہے، اس کے اکاؤنٹس، اس کا فولڈر اور MCP سرورز |
| `ccam group add <group> --items <items> [--members <accounts>]` | گروپ بناتی ہے؛ اس کے اکاؤنٹس کے آئٹمز اس میں منتقل ہو جاتے ہیں |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | گروپ میں موجود ہر اکاؤنٹ کے لیے بدلتی ہے کہ گروپ کیا شیئر کرتا ہے |
| `ccam group join <group> <account>...` | اکاؤنٹس کو گروپ میں ڈالتی ہے |
| `ccam group leave <group> <account>...` | اکاؤنٹس کو باہر نکالتی ہے؛ ہر ایک ایک کاپی رکھتا ہے، MCP سرورز سمیت |
| `ccam group rename <group> <new-name>` | صرف کنفیگ میں نام بدلتی ہے |
| `ccam group remove <group> [--leave]` | گروپ ہٹاتی ہے |

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--from <account>` | `add`، `set`: وہ اکاؤنٹ جس کی کاپی سب سے پہلے گروپ میں جاتی ہے | ڈیفالٹ اکاؤنٹ، پھر نمبر کے حساب سے |
| `--replace` | `join`: پہلے ان گروپس کو چھوڑتا ہے جو وہی آئٹمز شیئر کرتے ہیں | بند |
| `--fresh` | جو اکاؤنٹس کسی آئٹم کی شیئرنگ بند کریں وہ کاپی کے بجائے اس کے بغیر شروع ہوتے ہیں (امپورٹ کیا گیا فولڈر پھر بھی ایک کاپی رکھتا ہے) | بند |
| `--leave` | `remove`: پہلے گروپ کے اکاؤنٹس کو باہر نکالتا ہے | جب گروپ میں اکاؤنٹس ہوں تو ضروری |
| `--prefer <group\|mine\|merge>` | ٹکراؤ کا فیصلہ کیسے ہو | پوچھا جاتا ہے؛ `--yes` کے ساتھ ہر فائل کا ڈیفالٹ |
| `-y`، `--yes` | پوچھے بغیر چلاتا ہے؛ ٹکراؤ `--prefer` کے مطابق طے ہوتے ہیں | بند |
| `--dry-run` (گلوبل) | منصوبہ پرنٹ کرتا ہے، کچھ تبدیل نہیں کرتا | بند |

اکاؤنٹس نام یا نمبر سے دیے جاتے ہیں؛ آئٹمز کاما سے الگ کی گئی فہرست کے طور پر۔ `ccam add`، `--group <group>` (کئی بار دیا جا سکتا ہے) یا `--no-group` لیتا ہے۔

## یہ کیسے کام کرتا ہے

متعلقہ فائلیں:

- `~/.config/ccam/config.toml`: `[[group]]` اندراجات: `name`، `items`، `members` اور، جب فولڈر کا نام مختلف ہو، `dir`۔
- `~/.local/share/ccam/profiles/<name>/`: اکاؤنٹ کا `CLAUDE_CONFIG_DIR`۔
- `~/.local/share/ccam/groups/<folder>/`: شیئر کی گئی کاپیاں۔

جب `config.toml` پڑھا جاتا ہے تو ہر اکاؤنٹ کے گروپس ہر آئٹم کے لیے ایک فولڈر میں حل کیے جاتے ہیں۔ پانچ اکاؤنٹس کے لیے اس میں تقریباً 25 µs لگتے ہیں، اس لیے لانچ اب بھی صرف `config.toml` پڑھتا ہے۔ ایسی کنفیگ جس میں کسی اکاؤنٹ کو ایک آئٹم دو گروپس سے ملے، غلطی ہے، اور جب تک اسے ٹھیک نہ کیا جائے کچھ شروع نہیں ہوتا۔

تبدیلی کنفیگ کی ایک کاپی پر کام کرتی ہے:

1. کمانڈ کاپی میں ترمیم کرتی ہے (`AddGroup`، `SetGroupItems`، `JoinGroup` …)، جو اصولوں کی جانچ کرتی ہے۔
2. `internal/group` پرانی اور نئی کنفیگ کا موازنہ کرتا ہے اور ہر اس اکاؤنٹ کے لیے جس کے آئٹمز بدلیں، `share.Build` سے ایک منصوبہ بناتا ہے۔ منصوبے ایک کے بعد ایک بنتے ہیں، اور ہر منصوبہ دیکھتا ہے کہ پچھلے منصوبے گروپ فولڈرز میں کیا ڈالیں گے، اس لیے صرف پہلا اکاؤنٹ اپنی کاپی اندر منتقل کرتا ہے۔ پہلی کاپی گروپ فولڈر سے آتی ہے اگر اس میں پہلے سے ہو، پھر `--from` سے، پھر ڈیفالٹ اکاؤنٹ سے، پھر باقی اکاؤنٹس سے نمبر کے حساب سے۔
3. اگر ان میں سے کوئی اکاؤنٹ چل رہا ہو، چاہے Claude Code کسی بھی طرح شروع ہوا ہو، تو ccam انکار کرتا ہے (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``)، اور Developer Mode کے بغیر Windows پر بھی، اگر کسی فائل کو لنک کی ضرورت ہو۔
4. ٹکراؤ کا فیصلہ ہوتا ہے، سوال کے ذریعے یا `--prefer` سے، منصوبہ پرنٹ ہوتا ہے، اور ccam پوچھتا ہے۔
5. محفوظ کرنے سے عین پہلے ccam دوبارہ جانچتا ہے: جو اکاؤنٹ سوالات کے انتظار کے دوران شروع ہوا ہو وہ تبدیلی روک دیتا ہے، اور انہی گروپس میں کسی دوسرے ccam کی تبدیلی بھی (`the groups changed while this plan waited; nothing was moved, make the change again`)۔ کچھ بھی منتقل ہونے سے پہلے نئی کنفیگ محفوظ ہو جاتی ہے، اور تبدیلی `~/.local/state/ccam/group-change.json` میں درج ہوتی ہے۔ اگر کوئی مرحلہ آدھے راستے میں رک جائے تو یہ اندراج باقی رہتا ہے: `ccam doctor --fix` تبدیلی کو اس کی ترتیب میں، آپ کے دیے گئے جوابات کے ساتھ، MCP کے مراحل سمیت، مکمل کرتا ہے، اور تب تک گروپس کی کوئی اور تبدیلی شروع نہیں ہوتی (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``)۔
6. MCP کے مراحل پہلے چلتے ہیں۔ جب اکاؤنٹس کسی گروپ سے ایکسٹینشنز نکالتے ہیں تو گروپ کی پلگ اِن فہرستیں، جہاں بھی ان کے راستے کسی اکاؤنٹ کے فولڈر سے گزرتے ہوں، گروپ فولڈر کی طرف موڑ دی جاتی ہیں (منصوبے میں `group plugin lists:`)۔ پھر منصوبے ترتیب سے چلتے ہیں۔ ہر مرحلہ ڈیٹا کو گروپ یا بیک اپ فولڈر میں منتقل کرتا ہے؛ صرف وہ فائلیں جو گروپ کے پاس پہلے سے بائٹ بہ بائٹ موجود ہوں بیک اپ کے بجائے ہٹا دی جاتی ہیں۔

`share.Layout` ان راستوں کی اجازت شدہ فہرست ہے جو ہر آئٹم میں آتے ہیں۔ جو چیز اس میں نہیں وہ کبھی گروپ تک نہیں پہنچ سکتی۔ `share.Never` میں `.credentials.json` اور `.claude.json` کے نام ہیں، اور اگر ان میں سے کوئی کبھی layout میں نظر آئے تو `Build` اندرونی غلطی کے ساتھ رک جاتا ہے۔ `mcp` آئٹم کا کوئی راستہ نہیں: سرورز کی فہرست `.claude.json` کے اندر ہوتی ہے، اس لیے [MCP شیئرنگ](/features/mcp) گروپ کی `mcp.json` کے ذریعے کام کرتی ہے۔

جب `side` ہسٹری شیئر کرنا شروع کرے اور گروپ میں پہلے سے مواد ہو، تو ایک فولڈر `projects/` کے ساتھ یہ ہوتا ہے:

<!-- motion:share-on -->

اگر گروپ میں اس راستے پر ابھی کچھ نہ ہو تو یہ اکاؤنٹ اسے شیئر کرنے والا پہلا ہے اور اس کا فولڈر جوں کا توں وہاں منتقل ہو جاتا ہے۔ ورنہ ہر وہ فائل جو گروپ میں نہیں گروپ میں منتقل ہوتی ہے، اور جو فائلیں دونوں طرف موجود ہوں وہ اکاؤنٹ کے فولڈر میں رہتی ہیں، اس لیے گروپ کا ورژن برقرار رہتا ہے۔ جو فائلیں گروپ والی فائلوں جیسی ہوں وہ ہٹا دی جاتی ہیں، اور منصوبہ ان فائلوں کے نام بتاتا ہے جو مختلف ہوں (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`)۔ پھر اکاؤنٹ کا فولڈر اکاؤنٹ فولڈر کے اندر `.ccam-backup/<time>/` میں جاتا ہے، اور ایک لنک اس کی جگہ لے لیتا ہے۔ اگر گروپ کی اینٹری خود شامل ہونے والے اکاؤنٹ کی اپنی اینٹری کا لنک ہو، تو اکاؤنٹ کا مواد اس کی جگہ گروپ میں چلا جاتا ہے۔ گفتگو کی فائلوں کے نام ان کے سیشن ID پر رکھے جاتے ہیں، اس لیے کئی اکاؤنٹس کی ہسٹریاں بغیر ٹکراؤ کے یکجا ہو جاتی ہیں۔ `plugins/` کی رجسٹریز اکاؤنٹ فولڈر سے گزرنے والے مکمل راستے محفوظ کرتی ہیں؛ ccam انہیں گروپ کے فولڈر کی طرف موڑ دیتا ہے اور اکاؤنٹ کی اینٹریز اس رجسٹری میں ملا دیتا ہے جو گروپ میں پہلے سے ہو۔

چھوڑنا لنک کو الٹ دیتا ہے:

<!-- motion:share-off -->

ccam پہلے گروپ کا موجودہ مواد اکاؤنٹ فولڈر میں `.ccam-backup/incomplete-<name>` میں کاپی کرتا ہے، اور لنک کو صرف تب بدلتا ہے جب کاپی مکمل ہو جائے، اس لیے رک جانے والی کاپی لنک کو اپنی جگہ چھوڑ دیتی ہے اور `ccam doctor --fix` دوبارہ کاپی کرتا ہے (ایسی تکمیل کا منصوبہ `restore <entry>: … (copied back from the group before a stop)` دکھاتا ہے)۔ یہ لنک صرف تب ہٹاتا ہے جب وہ اب بھی groups فولڈر کے اندر اسی طرح اشارہ کرتا ہو جیسے منصوبہ بنتے وقت کرتا تھا؛ جو راستہ اس کے بعد بدل گیا ہو، جیسے لنک کی جگہ کوئی اصل فائل، اسے جوں کا توں چھوڑ دیا جاتا ہے اور تبدیلی رک جاتی ہے۔ `--fresh` کے ساتھ کاپی چھوڑ دی جاتی ہے۔ گروپ اپنا مواد دوسرے اکاؤنٹس کے لیے برقرار رکھتا ہے۔ جو اکاؤنٹ کسی آئٹم کو ایسے گروپ میں لے جائے جس کا فولڈر ابھی خالی ہو، وہ جو کچھ دیکھ رہا تھا اس کی ایک کاپی ساتھ لاتا ہے، اس لیے اس کی فائلیں برقرار رہتی ہیں؛ یہ کاپی مکمل ہونے تک `<group>/.ccam-incomplete-<name>` میں تیار کی جاتی ہے۔ جو اکاؤنٹ MCP کی شیئرنگ بند کرے وہ گروپ کے سرورز رکھتا ہے: وہ اس کے اپنے سرورز کے ساتھ اس کی اپنی `.claude.json` میں جاتے ہیں (ایک بیک اپ رکھا جاتا ہے)، اور اسی نام کے اس کے اپنے سرورز برقرار رہتے ہیں۔ `--fresh` کے ساتھ منصوبہ اس کے بجائے ان سرورز کی فہرست دیتا ہے جو لوڈ ہونا بند ہو جائیں گے۔ دوبارہ شامل ہونا وہ کاپیاں ہٹا دیتا ہے جنہیں کسی نے نہیں بدلا۔ `extensions` چھوڑنے سے گروپ کی پلگ اِن فہرستیں بھی، جہاں ان کے راستے چھوڑنے والے اکاؤنٹ کے فولڈر سے گزرتے تھے، گروپ فولڈر کی طرف موڑ دی جاتی ہیں؛ دوسرے اکاؤنٹس کے لیے ان کا مطلب وہی رہتا ہے۔

ایک فائل جو دونوں طرف مختلف مواد کے ساتھ موجود ہو، ٹکراؤ ہے:

<!-- motion:share-conflict -->

`settings.json` جیسی اکیلی فائلوں کے لیے آپ چنتے ہیں؛ فولڈرز کے اندر گروپ کی کاپی رہتی ہے، اور منصوبہ ان فائلوں کے نام بتاتا ہے جو مختلف ہوں۔ `history.jsonl` میں `merge by line`، `keep group` اور `use mine` ملتے ہیں، جن میں انضمام پہلے ہے۔ دوسری `.json` فائلوں میں `keep group`، `use mine` اور `merge JSON` ملتے ہیں۔ باقی ہر چیز میں `keep group` اور `use mine`۔ `merge JSON` آبجیکٹس کو recursive طریقے سے ملاتا ہے، arrays کو بغیر تکرار کے جوڑتا ہے، اور جب دونوں طرف ایک ہی key کی قدر مختلف ہو تو گروپ کی قدر رکھتا ہے، تاکہ دوسرے اکاؤنٹس کو کوئی تبدیلی نظر نہ آئے۔ `keep group` کے سوا ہر انتخاب پہلے گروپ کی کاپی اکاؤنٹ کے `.ccam-backup/<time>/group/` میں محفوظ کرتا ہے۔ آپ کی اپنی فائل ہمیشہ بیک اپ میں جاتی ہے۔

ہر لانچ پر `share.Check` layout کے ہر راستے کو اکاؤنٹ کے گروپس سے صرف `lstat` اور `readlink` کے ذریعے موازنہ کرتا ہے۔ غائب لنک فوراً بنا دیا جاتا ہے۔ جہاں لنک ہونا چاہیے وہاں اصل فائل، غلط جگہ کا لنک، یا باسی لنک لانچ روک دیتا ہے اور `ccam doctor --fix` کی طرف اشارہ کرتا ہے۔

### نام بدلنا، ہٹانا اور پرانا فارمیٹ

- `ccam group rename` صرف کنفیگ بدلتا ہے۔ فولڈر اپنا نام برقرار رکھتا ہے اور گروپ اسے `dir` کے طور پر درج کرتا ہے، اس لیے کچھ منتقل نہیں ہوتا اور کچھ کھو نہیں سکتا۔
- `ccam group remove` اکاؤنٹس والے گروپ کو ہٹانے سے انکار کرتا ہے جب تک آپ `--leave` شامل نہ کریں۔ تب ہر اکاؤنٹ ایک کاپی رکھتا ہے، اور فولڈر `~/.local/share/ccam/trash/` میں چلا جاتا ہے، جب تک کوئی دوسرا گروپ اسے استعمال نہ کر رہا ہو۔
- گروپس سے پہلے، `config.toml` ہر اکاؤنٹ کے لیے ایک `[profile.share]` جدول محفوظ کرتا تھا۔ ccam اس فارمیٹ کو کچھ منتقل کیے بغیر گروپس کے طور پر پڑھتا ہے۔ پہلی بار محفوظ کرنے پر یہ نیا فارمیٹ لکھتا ہے اور پرانی فائل `config.toml.bak` کے طور پر رکھتا ہے۔ جو آئٹمز ایک ہی اکاؤنٹس شیئر کرتے ہوں وہ فولڈر کے نام والا ایک گروپ بن جاتے ہیں؛ باقی کو `default-history` جیسے نام ملتے ہیں اور وہ وہی فولڈر استعمال کرتے رہتے ہیں۔

### Windows پر

Windows symbolic links صرف Developer Mode (Settings > System > For developers) میں یا ایڈمنسٹریٹر کے لیے بناتا ہے۔ یہ حق نہ ہو تو شیئر کیا گیا فولڈر ایک directory junction بن جاتا ہے، جو Claude Code کے لیے ویسے ہی کام کرتا ہے۔ شیئر کی گئی فائل (`settings.json`، `CLAUDE.md`، `history.jsonl`) صرف symbolic link ہی ہو سکتی ہے، اس لیے جس گروپ تبدیلی کو اس کی ضرورت ہو وہ کچھ بھی منتقل ہونے سے پہلے مسترد کر دی جاتی ہے۔ اسی وجہ سے وہاں پہلا `ccam add` گروپ `default` کو صرف `mcp` اور `extensions` کے ساتھ بناتا ہے اور یہ بتا دیتا ہے، اور کسی ایسے گروپ میں `ccam add --group` مسترد کر دیا جاتا ہے جو سیٹنگز، میموری یا ہسٹری شیئر کرتا ہو۔

## حفاظت اور خاص صورتیں

- **چلتے ہوئے اکاؤنٹس۔** جب Claude Code ان اکاؤنٹس میں سے کسی کے طور پر چل رہا ہو جنہیں تبدیلی منتقل کرتی ہے، تو تبدیلی مسترد کر دی جاتی ہے، کیونکہ سیشن ان فائلوں میں لکھتا رہتا ہے، چاہے اسے ccam نے شروع کیا ہو یا نہیں (سادہ `claude`، کوئی IDE، بیک گراؤنڈ سیشنز)۔ یہ جانچ محفوظ کرنے سے عین پہلے دہرائی جاتی ہے۔ چھوڑنا اور شامل ہونا صرف نام لیے گئے اکاؤنٹ سے متعلق ہیں؛ `extensions` چھوڑنا گروپ کی پلگ اِن فہرستوں میں راستے بھی دوبارہ لکھتا ہے، دوسروں کے لیے ان کا مطلب بدلے بغیر۔ `--dry-run` پھر بھی منصوبہ پرنٹ کرتا ہے۔
- **بیک اپس۔** ہر منصوبے کو اکاؤنٹ فولڈر میں اپنا `.ccam-backup/<time>/` فولڈر ملتا ہے۔ ccam انہیں کبھی حذف نہیں کرتا۔
- **آپ کے اپنے لنکس۔** اگر کوئی راستہ پہلے سے آپ کا بنایا ہوا لنک ہو، مثلاً کسی dotfiles ریپوزیٹری کی طرف، اور گروپ خالی ہو، تو ccam لنک کا ہدف گروپ میں کاپی کرتا ہے اور خود لنک کو بیک اپ میں منتقل کرتا ہے۔ ہدف کو نہیں چھیڑا جاتا۔
- **امپورٹ کیے گئے فولڈرز۔** جو اکاؤنٹ `ccam import` سے اپنی جگہ پر رجسٹر کیا گیا ہو (جیسے `~/.claude`) وہ اپنی جگہ پر ہی بدلتا ہے، اس لیے سادہ `claude` کو بھی تبدیلی نظر آتی ہے؛ منصوبہ یہ بتاتا ہے۔ جب یہ کسی گروپ کو چھوڑتا ہے تو ہمیشہ ایک کاپی رکھتا ہے، `--fresh` کے ساتھ بھی، اور `ccam remove` اکاؤنٹ کو بھولنے سے پہلے اسے اسی طرح اس کے گروپس سے نکالتا ہے۔ جب تک یہ کسی گروپ میں ہے، اس کی شیئر کی گئی فائلیں گروپ فولڈر میں رہتی ہیں، اس لیے `~/.local/share/ccam` حذف کرنے سے پہلے اکاؤنٹ ہٹا دیں (یا گروپس چھوڑ دیں)۔
- **ایک فولڈر میں دو اکاؤنٹس۔** جو اکاؤنٹس ہسٹری شیئر کرتے ہیں وہ ایک ہی فائلوں میں لکھتے ہیں۔ جب دوسرا اکاؤنٹ کسی ایسے فولڈر میں شروع ہو جہاں ایک پہلے سے چل رہا ہو، تو ccam انتباہ دیتا ہے اور پوچھتا ہے۔
- **ہاتھ سے ترامیم۔** جو رکن اکاؤنٹ نہ ہو اسے نظر انداز کیا جاتا ہے اور `ccam doctor` اس کی اطلاع دیتا ہے۔ `local`، `on`، `off` اور `none` گروپس کے نام نہیں ہو سکتے۔

## متعلقہ

- [ڈیش بورڈ](/features/dashboard)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [شیئر کیے گئے MCP سرورز](/features/mcp)
- [موجودہ سیٹ اپ کاپی کریں](/features/migrate)
- [تمام کمانڈز](/docs/commands)
