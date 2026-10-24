---
title: "شیئر کیے گئے MCP سرورز"
description: "ccam mcp add، add-json، list، remove اور import سے ہر گروپ کے لیے MCP سرورز کی ایک فہرست رکھیں، جو --mcp-config سے دی جاتی ہے یا ہر اکاؤنٹ میں sync کی جاتی ہے۔"
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code یوزر MCP سرورز کو `.claude.json` میں رکھتا ہے، وہی فائل جس میں اکاؤنٹ کا لاگ اِن ہوتا ہے، اس لیے یہ فائل اکاؤنٹس کے درمیان شیئر نہیں ہو سکتی۔ ccam شیئر کیے گئے سرورز ایک گروپ فائل `mcp.json` میں رکھتا ہے، اور لانچ کے وقت اسے گروپ کے ہر اکاؤنٹ کو دیتا ہے۔ کسی سرور کو اپنے تمام اکاؤنٹس کے لیے ایک بار شامل کرنے کے لیے `ccam mcp` استعمال کریں، اور کسی ایک سیشن کے اندر شامل کیے گئے سرور کو دوسروں تک لے جانے کے لیے `ccam mcp import`۔

## استعمال

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

تبدیلی ہر اکاؤنٹ کے اگلے لانچ سے لاگو ہوتی ہے؛ چلتے ہوئے سیشنز وہی سرورز رکھتے ہیں جن کے ساتھ وہ شروع ہوئے تھے:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

گروپس کے نیچے `ccam mcp list` ہر اس سرور کا نام بتاتا ہے جو دوسرے اکاؤنٹس تک نہیں پہنچتا: کسی اکاؤنٹ کے اپنے لیے شامل کیے گئے سرور کے لیے `only in account work: pg`، ایک فولڈر کے لیے شامل کیے گئے سرور کے لیے `only in account side, for /home/me/code/app only: notes` (سادہ `claude mcp add`، نیچے دیکھیں)، جب کسی اکاؤنٹ کے اپنے سرور کا نام گروپ کے کسی سرور جیسا ہو تو `account work has its own …, used in its sessions in place of the group's`، اس فولڈر کی `.mcp.json` ان اکاؤنٹس کے ساتھ جن کی منظوری کا وہ `waiting for approval in` ہے، اور وہ سرورز جو `in no account, only in ~/.claude.json (plain claude's)` ہیں، اس [`ccam migrate --mcp`](/features/migrate) کمانڈ کے ساتھ جو انہیں کاپی کرتی ہے۔

ایسا نام شامل کرنے پر جو گروپ میں پہلے سے ہو، اس کی تعریف بدل جاتی ہے اور `added` کے بجائے `updated` پرنٹ ہوتا ہے۔ سرورز کسی اکاؤنٹ تک صرف تب پہنچتے ہیں جب وہ اس گروپ میں ہو اور گروپ MCP شیئر کرتا ہو (`ccam group set <group> --add mcp`)؛ جب گروپ ایسا نہ کرتا ہو تو `add` یہ بتا دیتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `-g`، `--group <group>` | `add`، `add-json` اور `remove` کے لیے ہدف گروپ | وہ واحد گروپ جس کے ذریعے کوئی اکاؤنٹ MCP شیئر کرتا ہے، ورنہ `default` |
| `-- <command> [args…]` (`add`) | یہ کمانڈ stdio سرور کے طور پر چلاتا ہے | |
| `--url <url>` (`add`) | ریموٹ سرور کا URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` سرور کا ٹرانسپورٹ | `http` |
| `--header <'Name: value'>` (`add`) | HTTP ہیڈر، کئی بار دیا جا سکتا ہے | کوئی نہیں |
| `-e`، `--env <KEY=VALUE>` (`add`) | stdio سرور کے لیے انوائرنمنٹ ویری ایبل، کئی بار دیا جا سکتا ہے | کوئی نہیں |
| `--force` (`import`) | اسی نام کا گروپ سرور جو مختلف ہو بدل دیتا ہے؛ پرانی `mcp.json` کا بیک اپ لیا جاتا ہے | بند |
| `--keep-local` (`import`) | اکاؤنٹ کی `.claude.json` میں اینٹریز برقرار رکھتا ہے | بند |
| `-h`، `--help` | ہر کمانڈ کی مدد | |
| `--dry-run` (گلوبل) | صرف `mcp list` اسے لاگو کرتا ہے؛ `add`، `add-json`، `remove` اور `import` اسے مسترد کرتے ہیں اور کچھ تبدیل نہیں کرتے | بند |

`add` کو یا تو `--` کے بعد ایک کمانڈ چاہیے یا `--url`؛ ورنہ یہ `a command to run (after --) or --url is required` کے ساتھ ناکام ہوتا ہے۔ `--header`، `--url` سرورز پر لاگو ہوتا ہے اور `--env` کمانڈ والے سرورز پر۔ `add-json` وہی JSON لیتا ہے جو `claude mcp add-json` لیتا ہے، اور تعریف ایک آبجیکٹ ہونی چاہیے جس میں `command` یا `url` ہو۔ `import` ہمیشہ وہ گروپ استعمال کرتا ہے جس کے ذریعے اکاؤنٹ MCP شیئر کرتا ہے، اور `list` ہر گروپ دکھاتا ہے، اس لیے `-g` ان پر اثر نہیں ڈالتا۔

## یہ کیسے کام کرتا ہے

متعلقہ فائلیں:

- `~/.local/share/ccam/groups/<group>/mcp.json`: شیئر کیے گئے سرورز، `{"mcpServers": {…}}`، وہ فارمیٹ جو `--mcp-config` پڑھتا ہے۔ یہ `fsutil.AtomicWrite` اور موڈ `0600` کے ساتھ لکھی جاتی ہے، کیونکہ سرور اینٹریز میں اکثر API keys ہوتی ہیں، اور صرف `<group>/.mcp.json.lock` لاک کے تحت، تاکہ بیک وقت چلنے والے دو `ccam mcp add` دونوں اپنا سرور رکھیں۔
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code کی اپنی فائل جس میں اکاؤنٹ کے سرورز ہوتے ہیں۔ ccam اس کی صرف `mcpServers` قدر `profile.TopLevelValue` اسکینر سے پڑھتا ہے اور لانچ کے راستے پر کبھی پوری فائل پارس نہیں کرتا۔
- `~/.local/state/ccam/mcp/<name>/`: `seen` فائل، `sync-stamp` (ccam کی لکھی ہوئی ہر کاپی کے hash کے ساتھ، تاکہ یہ ایسی کاپی جسے کسی نے نہیں بدلا اس کاپی سے الگ پہچان سکے جس میں آپ نے ترمیم کی)، `sync-checked` (جو sync حکمتِ عملی نے آخری بار دیکھا)، `kept` (وہ سرورز جو اکاؤنٹ نے MCP شیئرنگ چھوڑتے وقت رکھے) اور `.claude.json` کے بیک اپس۔

دو حکمتِ عملیاں ہیں، جو `~/.config/ccam/config.toml` میں `mcp_strategy` سے مقرر ہوتی ہیں۔ ڈیفالٹ فلیگ والی حکمتِ عملی ہے:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` اکاؤنٹ کا MCP گروپ تلاش کرتا ہے۔ `mcp.FlagArgs` ایک `stat` سے جانچتا ہے کہ گروپ کی `mcp.json` موجود ہے، اسے پارس کیے بغیر۔
2. `Plan.InsertArgs`، `--mcp-config=<path>` کو `claude` کے فوراً بعد، آپ کے اپنے آرگیومنٹس سے پہلے رکھتا ہے، اس لیے آپ کا اپنا `--mcp-config` یا `--strict-mcp-config` پھر بھی اس کے اوپر لاگو ہوتا ہے۔
3. `mcp.NewLocal` جانچتا ہے کہ اکاؤنٹ کی `.claude.json` پچھلے لانچ کے بعد بدلی ہے یا نہیں، اور صرف تب اس کے سرورز پارس کرتا ہے۔
4. `syscall.Exec`، ccam کی جگہ Claude Code لے آتا ہے۔

راستہ ایک آرگیومنٹ کے طور پر دیا جاتا ہے، `--mcp-config=<path>`، دو کے طور پر نہیں۔ Claude Code کا `--mcp-config` کئی قدریں لیتا ہے، اس لیے `claude --mcp-config a.json "prompt"` پرامپٹ کو دوسری کنفیگ فائل سمجھ کر پڑھے گا اور `MCP config file not found` کے ساتھ ناکام ہو گا۔ ایک آرگیومنٹ والی شکل وہیں رک جاتی ہے، اور بعد میں آنے والا پرامپٹ یا ذیلی کمانڈ کام کرتی ہے۔ `ccam work mcp list`، `get` اور `remove` کو فلیگ نہیں ملتا، کیونکہ Claude Code وہاں اسے مسترد کرتا ہے؛ اس کے بجائے ccam اسٹینڈرڈ ایرر پر گروپ کے سرورز کے نام بتاتا ہے ([تفصیلات](/features/launch#shared-mcp-servers))۔ Claude Code، `--mcp-config` کے سرورز اکاؤنٹ کے اپنے سرورز کے علاوہ لوڈ کرتا ہے؛ نام کے ٹکراؤ پر `.claude.json` کی اینٹری کو ترجیح ملتی ہے۔ فلیگ والی حکمتِ عملی `.claude.json` صرف Remote Control کے لیے، جب کوئی اکاؤنٹ MCP شیئرنگ چھوڑے یا دوبارہ شامل ہو ([نیچے](#leaving-and-joining))، اور `ccam mcp import` کے لیے لکھتی ہے۔

`claude remote-control` (`rc`)، `--mcp-config` نہیں لیتا، اور اس کے سیشنز سرورز صرف `.claude.json` سے پڑھتے ہیں۔ `ccam work rc` کے لیے ccam یہ فلیگ چھوڑ دیتا ہے اور، ٹرمینل پر، پوچھتا ہے کہ کیا گروپ کے سرورز بیک اپ کے ساتھ `.claude.json` میں کاپی کیے جائیں، جیسے `sync` حکمتِ عملی کرتی ہے۔ یہ تب تک پوچھتا ہے جب تک کاپی موجود نہ ہو، پھر دوبارہ صرف تب جب گروپ کی `mcp.json` بدلے یا کاپی کیا گیا کوئی سرور غائب ہو جائے ([Remote Control](/features/launch#remote-control))۔ چونکہ نام کے ٹکراؤ پر ایسی کاپی کو `--mcp-config` پر ترجیح ملتی ہے، اس لیے `ccam mcp add`، `add-json`، `remove` اور `import` وہ کاپیاں اپ ڈیٹ کرتے ہیں جو گروپ کے دوسرے اکاؤنٹس کے پاس ہیں:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated` والی لائن صرف تب آتی ہے جب کوئی کاپی بدلی ہو۔ کسی پرانے ccam کی بنائی ہوئی کاپی بھی اسی طرح تازہ کی جاتی ہے، اور `ccam migrate` کے بعد والی کاپی بھی۔

کسی سیشن کے اندر `claude mcp add` سے شامل کیا گیا سرور صرف اسی اکاؤنٹ کی `.claude.json` میں جاتا ہے۔ سادہ `claude mcp add` local اسکوپ استعمال کرتا ہے: سرور صرف اسی ایک اکاؤنٹ کے لیے اسی ایک فولڈر میں لوڈ ہوتا ہے (`-s user` اسے ہر فولڈر میں اس اکاؤنٹ کا سرور بنا دیتا ہے)۔ ccam اسے اس وقت نہیں دیکھ سکتا، کیونکہ Claude Code اس کی جگہ لے چکا ہوتا ہے، اس لیے یہ اگلے لانچ پر جانچتا ہے۔ `seen` فائل میں `.claude.json` کا ترمیم کا وقت اور پہلے سے بتائے گئے سرورز کے نام ہوتے ہیں۔ جب وقت مختلف ہو تو ccam ناموں کو گروپ سے موازنہ کرتا ہے اور ہر نئے سرور کے لیے ایک بار، اِن پٹ کا انتظار کیے بغیر، ایک اشارہ پرنٹ کرتا ہے، دونوں حکمتِ عملیوں کے ساتھ:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

یہ اشارہ کسی اکاؤنٹ کے ایسے اپنے سرور کا نام بھی بتاتا ہے جس کا نام گروپ کے کسی سرور جیسا ہو، کیونکہ اس کے سیشنز میں اکاؤنٹ والے سرور کو ترجیح ملتی ہے۔

`mcp_strategy = "sync"` کے ساتھ، سرورز `claude mcp list` اور `/mcp` میں عام یوزر سرورز کی طرح نظر آتے ہیں، اور ccam لانچ سے پہلے انہیں ہر اکاؤنٹ کی `.claude.json` میں لکھتا ہے:

<!-- motion:mcp-sync -->

1. `mcp.Sync` گروپ کی `mcp.json` کے ترمیم کے وقت کو `sync-stamp` سے، اور اکاؤنٹ کی `.claude.json` کو اس سے موازنہ کرتا ہے جو اس نے آخری بار دیکھا تھا۔ اگر دونوں میں سے کوئی نہ بدلا ہو تو اور کچھ نہیں ہوتا۔ ورنہ یہ یہ بھی جانچتا ہے کہ گروپ کا ہر سرور اب بھی `.claude.json` میں ہے، اور غائب ہو جانے والی کاپی واپس رکھ دیتا ہے (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`)۔
2. یہ Claude Code کا اپنا لاک لیتا ہے، یعنی فائل کے ساتھ `.claude.json.lock` فولڈر، جو `mkdir` سے بنایا جاتا ہے۔ یہ 2 سیکنڈ تک انتظار کرتا ہے اور 10 سیکنڈ سے پرانے لاک کو باسی سمجھتا ہے، جیسے Claude Code کی لاک لائبریری کرتی ہے۔
3. یہ موجودہ `.claude.json` کو اسٹیٹ فولڈر میں `backups/` میں محفوظ کرتا ہے اور سب سے نئی 5 رکھتا ہے۔
4. یہ ہر گروپ سرور کو `mcpServers` میں مقرر کرتا ہے، وہ سرورز ہٹاتا ہے جو اس نے پہلے لکھے تھے اور اب گروپ میں نہیں ہیں، اور فائل کو اس کے اصل موڈ کے ساتھ ایٹامک طور پر لکھتا ہے، صرف تب جب کچھ بدلا ہو۔ یہ صرف وہی کاپیاں بدلتا یا ہٹاتا ہے جو اس نے لکھیں اور جنہیں بعد میں کسی نے نہیں بدلا۔ جو سرور اکاؤنٹ کا اپنا ہو اور گروپ کے کسی سرور کے نام سے ہو وہ اکاؤنٹ کا ہی رہتا ہے، اور ccam ایک بار یہ بتا دیتا ہے: `account work has its own MCP servers fs, so group default's of those names are not copied there`۔
5. یہ نیا stamp ہر کاپی کے hash کے ساتھ ریکارڈ کرتا ہے، اور Claude Code، `--mcp-config` کے بغیر شروع ہوتا ہے۔

جب اکاؤنٹ کا کوئی سیشن چل رہا ہو، چاہے Claude Code کسی بھی طرح شروع ہوا ہو، تو sync، `ccam: account work is already running, skipping MCP sync` کے ساتھ چھوڑ دیا جاتا ہے، کیونکہ Claude Code چلتے وقت `.claude.json` دوبارہ لکھتا ہے۔ ناکام sync، `ccam: MCP sync failed (launching anyway): …` پرنٹ کرتا ہے اور پھر بھی لانچ کرتا ہے۔

### چھوڑنا اور شامل ہونا

جو اکاؤنٹ MCP شیئر کرنا بند کرے (`ccam group leave`، `ccam group set --remove mcp`، `ccam group remove --leave`) وہ گروپ کے سرورز رکھتا ہے: پلان انہیں بیک اپ کے ساتھ اس کی `.claude.json` میں کاپی کرتا ہے، اور اسی نام کے اس کے اپنے سرورز برقرار رہتے ہیں۔ `--fresh` کے ساتھ اسے کوئی سرور نہیں ملتا، اور پلان ان سرورز کی فہرست دیتا ہے جو لوڈ ہونا بند ہو جائیں گے۔ جب وہ دوبارہ کسی ایسے گروپ میں شامل ہوتا ہے جس میں وہ سرورز ہوں، تو رکھی گئی وہ کاپیاں جنہیں کسی نے نہیں بدلا ہٹا دی جاتی ہیں، تاکہ گروپ کی بعد کی ترامیم اس تک پہنچیں۔

`ccam mcp import` کسی اکاؤنٹ کے اپنے سرورز اس کے گروپ میں منتقل کرتا ہے:

<!-- motion:mcp-import -->

ccam اکاؤنٹ کے سرورز پڑھتا ہے، اس کے اپنے اور وہ جو اس کے پاس صرف ایک فولڈر کے لیے ہیں، اور ہر ایک کو گروپ کے اسی نام کے سرور سے `mcp.SameDef` کے ذریعے موازنہ کرتا ہے، جو keys کی ترتیب اور اسپیسنگ کو نظر انداز کرتا ہے۔ فولڈر والا سرور ہر فولڈر کے لیے گروپ سرور بن جاتا ہے (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`)؛ ایک نام جس کی دو مختلف فولڈر تعریفیں ہوں مسترد کر دیا جاتا ہے۔ پلگ اِن سرورز پیش نہیں کیے جاتے، کیونکہ پلگ اِن انہیں خود لاتا ہے۔ جو سرور گروپ میں نہ ہو یا اسی تعریف کے ساتھ ہو وہ کاپی ہو جاتا ہے۔ مختلف سرور چھوڑ دیا جاتا ہے، کیونکہ گروپ کی تعریف ہی وہ ہے جو باقی ہر اکاؤنٹ استعمال کرتا ہے:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` کے ساتھ گروپ کا سرور بدل دیا جاتا ہے، اور پچھلی `mcp.json` پہلے گروپ فولڈر کے اندر `.ccam-backup/<time>/mcp.json` میں محفوظ کی جاتی ہے۔ پھر `mcp.Remove` کاپی کی گئی اینٹریز اکاؤنٹ کی `.claude.json` سے ہٹاتا ہے، اس کی اپنی اور فولڈر والی ایک ہی تحریر میں، اسی لاک کے تحت اور اسٹیٹ فولڈر میں بیک اپ کے بعد۔ گروپ کی `mcp.json` اسے لوڈ کرنے کے وقت سے لے کر ان اینٹریز کے ہٹنے تک لاک رہتی ہے، تاکہ اس دوران کوئی گروپ سے سرور نہ ہٹا دے۔ یہ ضروری ہے کیونکہ نام کے ٹکراؤ پر اکاؤنٹ کی اینٹری کو گروپ کی اینٹری پر ترجیح ملے گی۔ `--keep-local` انہیں رہنے دیتا ہے۔

## حفاظت اور خاص صورتیں

- **سائن اِن۔** جو سرورز OAuth سے سائن اِن کرتے ہیں وہ اپنے ٹوکنز سرورز کی فہرست میں نہیں بلکہ ہر اکاؤنٹ کی `.credentials.json` (`mcpOAuth`) میں رکھتے ہیں، اس لیے سرور شیئر کرنے سے اس کا سائن اِن شیئر نہیں ہوتا۔ ہر اکاؤنٹ میں ایک بار `/mcp` سے سائن اِن کریں۔ ccam کبھی سائن اِن کاپی نہیں کرتا، [`ccam migrate`](/features/migrate) بھی نہیں: کاپی کیا گیا ریفریش ٹوکن باقی ہر رکھنے والے کو لاگ آؤٹ کر دیتا ہے۔ [`ccam ext`](/features/ext) اور `ccam doctor` ان سرورز کے نام بتاتے ہیں جن میں کوئی اکاؤنٹ سائن اِن نہیں جبکہ کوئی دوسرا اکاؤنٹ سائن اِن ہے، اور `ccam next --continue` بتاتا ہے کہ نئے اکاؤنٹ میں ان میں سے کون سے نہیں ہیں۔ گروپ کو کبھی لاگ اِن فائل نہیں ملتی۔
- **چلتے ہوئے اکاؤنٹس۔** `import` گروپ میں کاپی کرتا ہے لیکن چلتے ہوئے اکاؤنٹ کی اپنی اینٹریز رہنے دیتا ہے اور یہ بتا دیتا ہے؛ سیشن ختم ہونے کے بعد اسے دوبارہ چلائیں۔ sync چلتے ہوئے اکاؤنٹس کو چھوڑ دیتا ہے۔ دونوں ccam کے بغیر شروع ہونے والے Claude Code کو بھی شمار کرتے ہیں۔
- **Remote Control کے لیے کاپیاں۔** چلتے ہوئے اکاؤنٹ کی کاپی کو نہیں چھیڑا جاتا: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` کے ساتھ کچھ اپ ڈیٹ کرنے کی ضرورت نہیں، کیونکہ ہر لانچ sync کرتا ہے۔
- **خراب فائلیں۔** اگر `.claude.json` درست JSON نہ ہو تو sync اسے نہیں چھیڑتا اور غلطی بتاتا ہے۔ اگر لاک 2 سیکنڈ میں نہ مل سکے تو ccam، `could not lock … (Claude Code is writing the file)` بتاتا ہے۔
- **وہ اکاؤنٹس جو MCP شیئر نہیں کرتے۔** `import` انکار کرتا ہے اور کہتا ہے کہ پہلے اکاؤنٹ کو کسی ایسے گروپ میں شامل کریں جو MCP شیئر کرتا ہو (`ccam group join <group> side`)۔
- **گروپس کے نام۔** `-g` صرف گروپ کا نام لیتا ہے؛ یہ `groups/` کے نیچے ایک فولڈر بن جاتا ہے، اس لیے `../x` جیسے نام مسترد کیے جاتے ہیں۔

## متعلقہ

- [اکاؤنٹس کے درمیان شیئرنگ](/features/share)
- [ccam کے ساتھ Remote Control استعمال کرنا](/guides/remote-control)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [موجودہ سیٹ اپ کاپی کریں](/features/migrate)
- [تمام کمانڈز](/docs/commands)
