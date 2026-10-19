---
title: "شیئر کیے گئے MCP سرورز"
description: "ccam mcp add، add-json، list، remove اور import سے ہر گروپ کے لیے MCP سرورز کی ایک فہرست رکھیں، جو --mcp-config سے دی جاتی ہے یا ہر اکاؤنٹ میں sync کی جاتی ہے۔"
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

ایسا نام شامل کرنے پر جو گروپ میں پہلے سے ہو، اس کی تعریف بدل جاتی ہے اور `added` کے بجائے `updated` پرنٹ ہوتا ہے۔ سرورز کسی اکاؤنٹ تک صرف تب پہنچتے ہیں جب وہ اس گروپ کے ساتھ MCP شیئر کرتا ہو (`ccam share <name> --mcp on`)؛ جب ابھی کوئی اکاؤنٹ MCP شیئر نہ کرتا ہو تو `add` یہ بتا دیتا ہے۔

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: شیئر کیے گئے سرورز، `{"mcpServers": {…}}`، وہ فارمیٹ جو `--mcp-config` پڑھتا ہے۔ یہ `fsutil.AtomicWrite` اور موڈ `0600` کے ساتھ لکھی جاتی ہے، کیونکہ سرور اینٹریز میں اکثر API keys ہوتی ہیں۔
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code کی اپنی فائل جس میں اکاؤنٹ کے سرورز ہوتے ہیں۔ ccam اس کی صرف `mcpServers` قدر `profile.TopLevelValue` اسکینر سے پڑھتا ہے اور لانچ کے راستے پر کبھی پوری فائل پارس نہیں کرتا۔
- `~/.local/state/ccam/mcp/<name>/`: `seen` فائل، `sync-stamp` اور `.claude.json` کے بیک اپس۔

دو حکمتِ عملیاں ہیں، جو `~/.config/ccam/config.toml` میں `mcp_strategy` سے مقرر ہوتی ہیں۔ ڈیفالٹ فلیگ والی حکمتِ عملی ہے:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` اکاؤنٹ کا MCP گروپ تلاش کرتا ہے۔ `mcp.FlagArgs` ایک `stat` سے جانچتا ہے کہ گروپ کی `mcp.json` موجود ہے، اسے پارس کیے بغیر۔
2. `Plan.InsertArgs`، `--mcp-config=<path>` کو `claude` کے فوراً بعد، آپ کے اپنے آرگیومنٹس سے پہلے رکھتا ہے، اس لیے آپ کا اپنا `--mcp-config` یا `--strict-mcp-config` پھر بھی اس کے اوپر لاگو ہوتا ہے۔
3. `mcp.NewLocal` جانچتا ہے کہ اکاؤنٹ کی `.claude.json` پچھلے لانچ کے بعد بدلی ہے یا نہیں، اور صرف تب اس کے سرورز پارس کرتا ہے۔
4. `syscall.Exec`، ccam کی جگہ Claude Code لے آتا ہے۔

راستہ ایک آرگیومنٹ کے طور پر دیا جاتا ہے، `--mcp-config=<path>`، دو کے طور پر نہیں۔ Claude Code کا `--mcp-config` کئی قدریں لیتا ہے، اس لیے `claude --mcp-config a.json "prompt"` پرامپٹ کو دوسری کنفیگ فائل سمجھ کر پڑھے گا اور `MCP config file not found` کے ساتھ ناکام ہو گا۔ ایک آرگیومنٹ والی شکل وہیں رک جاتی ہے، اور بعد میں آنے والا پرامپٹ یا `mcp list` جیسی ذیلی کمانڈ کام کرتی ہے۔ Claude Code، `--mcp-config` کے سرورز اکاؤنٹ کے اپنے سرورز کے علاوہ لوڈ کرتا ہے؛ نام کے ٹکراؤ پر `.claude.json` کی اینٹری کو ترجیح ملتی ہے۔ فلیگ والی حکمتِ عملی کبھی `.claude.json` نہیں لکھتی۔

کسی سیشن کے اندر `claude mcp add` سے شامل کیا گیا سرور صرف اسی اکاؤنٹ کی `.claude.json` میں جاتا ہے۔ ccam اسے اس وقت نہیں دیکھ سکتا، کیونکہ Claude Code اس کی جگہ لے چکا ہوتا ہے، اس لیے یہ اگلے لانچ پر جانچتا ہے۔ `seen` فائل میں `.claude.json` کا ترمیم کا وقت اور پہلے سے بتائے گئے سرورز کے نام ہوتے ہیں۔ جب وقت مختلف ہو تو ccam ناموں کو گروپ سے موازنہ کرتا ہے اور ہر نئے سرور کے لیے ایک بار، اِن پٹ کا انتظار کیے بغیر، ایک اشارہ پرنٹ کرتا ہے:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` کے ساتھ، سرورز `claude mcp list` اور `/mcp` میں عام یوزر سرورز کی طرح نظر آتے ہیں، اور ccam لانچ سے پہلے انہیں ہر اکاؤنٹ کی `.claude.json` میں لکھتا ہے:

<!-- motion:mcp-sync -->

1. `mcp.Sync` گروپ کی `mcp.json` کے ترمیم کے وقت کو `sync-stamp` سے موازنہ کرتا ہے۔ اگر یہ نہ بدلا ہو تو اکاؤنٹ کی `.claude.json` کھولی نہیں جاتی۔
2. یہ Claude Code کا اپنا لاک لیتا ہے، یعنی فائل کے ساتھ `.claude.json.lock` فولڈر، جو `mkdir` سے بنایا جاتا ہے۔ یہ 2 سیکنڈ تک انتظار کرتا ہے اور 10 سیکنڈ سے پرانے لاک کو باسی سمجھتا ہے، جیسے Claude Code کی لاک لائبریری کرتی ہے۔
3. یہ موجودہ `.claude.json` کو اسٹیٹ فولڈر میں `backups/` میں محفوظ کرتا ہے اور سب سے نئی 5 رکھتا ہے۔
4. یہ ہر گروپ سرور کو `mcpServers` میں مقرر کرتا ہے، وہ سرورز ہٹاتا ہے جو اس نے پہلے لکھے تھے اور اب گروپ میں نہیں ہیں، اکاؤنٹ کے اپنے سرورز برقرار رکھتا ہے، اور فائل کو اس کے اصل موڈ کے ساتھ ایٹامک طور پر لکھتا ہے۔
5. یہ نیا stamp اور اپنے منظم کردہ نام ریکارڈ کرتا ہے، اور Claude Code، `--mcp-config` کے بغیر شروع ہوتا ہے۔

جب اکاؤنٹ کا کوئی سیشن چل رہا ہو تو sync، `ccam: account work is already running, skipping MCP sync` کے ساتھ چھوڑ دیا جاتا ہے، کیونکہ Claude Code چلتے وقت `.claude.json` دوبارہ لکھتا ہے۔ ناکام sync، `ccam: MCP sync failed (launching anyway): …` پرنٹ کرتا ہے اور پھر بھی لانچ کرتا ہے۔ سیشن کے اندر شامل کیے گئے سرورز کا اشارہ صرف فلیگ والی حکمتِ عملی کا حصہ ہے۔

`ccam mcp import` کسی اکاؤنٹ کے اپنے سرورز اس کے گروپ میں منتقل کرتا ہے:

<!-- motion:mcp-import -->

ccam اکاؤنٹ کے سرورز `mcp.LocalServers` سے پڑھتا ہے اور ہر ایک کو گروپ کے اسی نام کے سرور سے `mcp.SameDef` کے ذریعے موازنہ کرتا ہے، جو keys کی ترتیب اور اسپیسنگ کو نظر انداز کرتا ہے۔ جو سرور گروپ میں نہ ہو یا اسی تعریف کے ساتھ ہو وہ کاپی ہو جاتا ہے۔ مختلف سرور چھوڑ دیا جاتا ہے، کیونکہ گروپ کی تعریف ہی وہ ہے جو باقی ہر اکاؤنٹ استعمال کرتا ہے:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` کے ساتھ گروپ کا سرور بدل دیا جاتا ہے، اور پچھلی `mcp.json` پہلے گروپ فولڈر کے اندر `.ccam-backup/<time>/mcp.json` میں محفوظ کی جاتی ہے۔ پھر `mcp.RemoveLocal` کاپی کی گئی اینٹریز اکاؤنٹ کی `.claude.json` سے ہٹاتا ہے، اسی لاک کے تحت اور اسٹیٹ فولڈر میں بیک اپ کے بعد۔ یہ ضروری ہے کیونکہ نام کے ٹکراؤ پر اکاؤنٹ کی اینٹری کو گروپ کی اینٹری پر ترجیح ملے گی۔ `--keep-local` انہیں رہنے دیتا ہے۔

## حفاظت اور خاص صورتیں

- **سائن اِن۔** جو سرورز OAuth سے سائن اِن کرتے ہیں وہ اپنے ٹوکنز سرورز کی فہرست میں نہیں بلکہ ہر اکاؤنٹ کی `.credentials.json` (`mcpOAuth`) میں رکھتے ہیں، اس لیے سرور شیئر کرنے سے اس کا سائن اِن شیئر نہیں ہوتا۔ ہر اکاؤنٹ میں ایک بار `/mcp` سے سائن اِن کریں، یا [`ccam migrate`](/features/migrate) سے `~/.claude` کے سائن اِن کاپی کریں۔ گروپ کو کبھی لاگ اِن فائل نہیں ملتی۔
- **چلتے ہوئے اکاؤنٹس۔** `import` گروپ میں کاپی کرتا ہے لیکن چلتے ہوئے اکاؤنٹ کی اپنی اینٹریز رہنے دیتا ہے اور یہ بتا دیتا ہے؛ سیشن ختم ہونے کے بعد اسے دوبارہ چلائیں۔ sync چلتے ہوئے اکاؤنٹس کو چھوڑ دیتا ہے۔
- **خراب فائلیں۔** اگر `.claude.json` درست JSON نہ ہو تو sync اسے نہیں چھیڑتا اور غلطی بتاتا ہے۔ اگر لاک 2 سیکنڈ میں نہ مل سکے تو ccam، `could not lock … (Claude Code is writing the file)` بتاتا ہے۔
- **وہ اکاؤنٹس جو MCP شیئر نہیں کرتے۔** `import` انکار کرتا ہے اور کہتا ہے کہ پہلے `ccam share side --mcp on` چلائیں۔
- **گروپس کے نام۔** `-g` صرف گروپ کا نام لیتا ہے؛ یہ `groups/` کے نیچے ایک فولڈر بن جاتا ہے، اس لیے `../x` جیسے نام مسترد کیے جاتے ہیں۔

## متعلقہ

- [اکاؤنٹس کے درمیان شیئرنگ](/features/share)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [موجودہ سیٹ اپ کاپی کریں](/features/migrate)
- [تمام کمانڈز](/docs/commands)
