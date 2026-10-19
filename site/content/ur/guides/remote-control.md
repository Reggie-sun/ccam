---
title: "ccam کے ساتھ Remote Control استعمال کرنا"
description: "ccam <name> rc سے Claude Code کا Remote Control کسی بھی اکاؤنٹ کے طور پر شروع کریں، دیکھیں کون سے آپشنز لاگو ہوتے ہیں، ccam شیئر کیے گئے MCP سرورز کاپی کرنے کا کیوں پوچھتا ہے، اور عام مسائل کیسے حل کریں۔"
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

`ccam work rc`، Claude Code کا Remote Control اکاؤنٹ `work` کے طور پر شروع کرتا ہے، تاکہ آپ اس کمپیوٹر پر چلنے والے سیشنز کو claude.ai/code یا Claude کی موبائل ایپ سے چلا سکیں۔ یہ `claude rc` کی طرح کام کرتا ہے، بس ایک فرق کے ساتھ: اگر اکاؤنٹ MCP سرورز شیئر کرتا ہے تو ccam ایک بار پوچھتا ہے کہ کیا انہیں اکاؤنٹ کی `.claude.json` میں کاپی کیا جائے، کیونکہ Remote Control کے سیشنز سرورز کہیں اور سے نہیں پڑھتے۔

## اسے شروع کریں

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[شارٹ کٹ](/features/shortcuts) بھی اسی طرح کام کرتا ہے: `ccwork rc` دراصل `ccam work rc` ہے۔ کمانڈ اس فولڈر میں چلائیں جس میں آپ کام کرنا چاہتے ہیں، پھر claude.ai/code یا Claude ایپ کھولیں، اسی Claude اکاؤنٹ سے سائن اِن کر کے جو `work` استعمال کرتا ہے۔

## آپشنز

یہاں دو طرح کے آپشنز ملتے ہیں: ccam کے اپنے، جو اکاؤنٹ کے نام سے پہلے آتے ہیں یا ccam کی سیٹنگز میں ہوتے ہیں، اور Remote Control کے، جو `rc` کے بعد آتے ہیں اور بغیر تبدیلی کے Claude Code تک پہنچتے ہیں۔

### ccam کے آپشنز اور سیٹنگز

| آپشن یا سیٹنگ | یہ کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--dry-run` (اکاؤنٹ سے پہلے) | انوائرنمنٹ اور کمانڈ پرنٹ کرتا ہے، اور اگر MCP کا نوٹس ہو تو وہ بھی۔ نہ پوچھتا ہے نہ کاپی کرتا ہے | بند |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | تب پوچھا جاتا ہے جب اکاؤنٹ MCP سرورز شیئر کرتا ہو اور اس کی `.claude.json` میں گروپ کی موجودہ فہرست نہ ہو۔ `y` انہیں کاپی کرتا ہے، `N` یا Enter ان کے بغیر شروع کرتا ہے | `N` |
| `~/.config/ccam/config.toml` میں `mcp_strategy = "sync"` | ہر لانچ سے پہلے، چاہے Remote Control ہو یا نہیں، گروپ کے سرورز `.claude.json` میں کاپی کرتا ہے، اس لیے کبھی کچھ نہیں پوچھا جاتا | `"flag"` |
| `ccam share work --mcp off` | اکاؤنٹ کے ساتھ MCP سرورز کی شیئرنگ بند کرتا ہے؛ پھر Remote Control صرف اکاؤنٹ کے اپنے سرورز کے ساتھ شروع ہوتا ہے اور کچھ نہیں پوچھا جاتا | آن |

### Remote Control کے آپشنز

`rc` کے بعد ہر چیز Claude Code کی ہے۔ یہ Claude Code 2.1.304 میں `claude remote-control` کے آپشنز ہیں؛ `ccam work rc --help` آپ کے ورژن کی فہرست پرنٹ کرتا ہے۔

| آپشن | یہ کیا کرتا ہے |
| --- | --- |
| `--name <name>` | سیشن کا نام جو claude.ai/code میں دکھایا جاتا ہے |
| `--remote-control-session-name-prefix <prefix>` | خودکار بنائے گئے سیشن ناموں کا سابقہ (ڈیفالٹ: ہوسٹ کا نام) |
| `-c`، `--continue` | اس سیشن سے دوبارہ جڑتا ہے جو Remote Control نے اس فولڈر کے لیے آخری بار ریکارڈ کیا، تقریباً 4 گھنٹوں کے اندر |
| `--session-id <id>` | کسی مخصوص سیشن سے دوبارہ جڑتا ہے |
| `--permission-mode <mode>` | نئے سیشنز کے لیے اجازت کا موڈ: `acceptEdits`، `auto`، `bypassPermissions`، `default`، `dontAsk` یا `plan` |
| `--spawn <mode>` | `same-dir` (ڈیفالٹ)، `worktree` (ہر سیشن کے لیے ایک git worktree) یا `session` (ایک سیشن، جو ختم ہونے پر باہر نکل جاتا ہے) |
| `--capacity <N>` | `same-dir` اور `worktree` موڈ میں ایک وقت میں زیادہ سے زیادہ سیشنز (ڈیفالٹ 32) |
| `--[no-]create-session-in-dir` | شروع میں موجودہ فولڈر میں پہلا سیشن بناتا ہے (ڈیفالٹ آن) |
| `--[no-]chrome` | نئے سیشنز کے لیے Claude in Chrome (ڈیفالٹ بند) |
| `-d`، `--debug`، `--debug-file <path>`، `-v`، `--verbose` | ڈیبگ اور تفصیلی آؤٹ پٹ |

`rc` کو اکاؤنٹ کے نام کے فوراً بعد رکھیں۔ ccam پہلے ایسے آرگیومنٹ کو جو `-` سے شروع نہ ہو، Claude Code کی کمانڈ مانتا ہے، اس لیے `ccam work rc --name laptop` میں اسے `rc` نظر آتا ہے، لیکن `ccam work --model opus rc` میں اسے `opus` نظر آتا ہے۔

## ccam کیا کرتا ہے

لانچ معمول کا راستہ اختیار کرتا ہے ([Claude Code شروع کرنا](/features/launch))، بس MCP والے مرحلے میں ایک تبدیلی کے ساتھ:

1. ccam اکاؤنٹ ڈھونڈتا ہے، اس کے شیئر لنکس جانچتا ہے اور `claude` کا راستہ طے کرتا ہے، جیسے کسی بھی لانچ میں۔
2. وہ دیکھتا ہے کہ کمانڈ `rc` یا `remote-control` ہے اور `--mcp-config` شامل نہیں کرتا۔
3. اگر اکاؤنٹ MCP سرورز شیئر کرتا ہے تو وہ گروپ کی `mcp.json` کا موازنہ `~/.local/state/ccam/mcp/work/` میں موجود `sync-stamp` سے کرتا ہے۔ جب `.claude.json` میں پہلے سے موجودہ سرورز ہوں، یا گروپ کے پاس کوئی سرور نہ ہو، تو وہ سیدھا مرحلہ 6 پر چلا جاتا ہے۔
4. ورنہ وہ ان سرورز کے نام بتاتا ہے جو سیشنز کو نہیں ملیں گے:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y` پر وہ `.claude.json` کا بیک اپ محفوظ کرتا ہے، Claude Code کا اپنا لاک لیتا ہے، گروپ کے سرورز `mcpServers` میں شامل کرتا ہے اور اکاؤنٹ کے دوسرے سرورز برقرار رکھتا ہے۔ `N` پر، ٹرمینل کے بغیر، یا جب اکاؤنٹ چل رہا ہو، وہ کچھ تبدیل نہیں کرتا اور یہی بتاتا ہے:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. وہ سیشن ریکارڈ کرتا ہے، `CLAUDE_CONFIG_DIR` کو اکاؤنٹ کے فولڈر پر سیٹ کرتا ہے اور `exec` کے ذریعے خود کو `claude rc` سے بدل دیتا ہے۔

اس کے بعد Remote Control ایک سرور کے طور پر چلتا رہتا ہے۔ آپ کے لیے شروع کیا گیا ہر سیشن اسی پروسیس کا چائلڈ ہوتا ہے، اس لیے وہ `work` کے طور پر چلتا ہے: `work` کے لاگ اِن، سیٹنگز اور گفتگو کی ہسٹری کے ساتھ، اور `work` کی استعمال کی حد میں شمار ہوتا ہے۔ جب تک سرور چل رہا ہے، ccam اکاؤنٹ کو چلتا ہوا مانتا ہے۔

ڈرائی رن فرق دکھاتا ہے۔ عام لانچ گروپ کے سرورز کو ایک آرگیومنٹ کے طور پر ساتھ لے جاتا ہے، Remote Control نہیں:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## یہ اس طرح کیوں کام کرتا ہے

**Remote Control، `--mcp-config` نہیں لیتا۔** عام لانچز میں ccam گروپ کے سرورز Claude Code کو `--mcp-config=<path>` کے طور پر دیتا ہے اور Claude Code کی فائلوں میں کبھی نہیں لکھتا۔ `claude remote-control` یہ آپشن اپنی کمانڈ کے نام سے پہلے مسترد کرتا ہے اور اس کے بعد کوئی نہیں لیتا، اس لیے ccam 0.6.0 سے پہلے `ccam work rc` ہر اس اکاؤنٹ کے لیے ناکام ہو جاتا تھا جو MCP سرورز شیئر کرتا تھا۔ اب ccam یہ آپشن چھوڑ دیتا ہے۔

**اس کے سیشنز سرورز صرف `.claude.json` سے پڑھتے ہیں۔** Remote Control جو سیشنز شروع کرتا ہے وہ کمانڈ لائن سے آپشنز وراثت میں نہیں لیتے۔ یوزر MCP سرورز انہیں صرف ایک جگہ ملتے ہیں: اکاؤنٹ کی `.claude.json`۔ اسی فائل میں اکاؤنٹ کی لاگ اِن تفصیلات اور ہر فولڈر کی حالت بھی ہوتی ہے، اس لیے ccam اسے `settings.json` کی طرح کسی شیئر کی گئی کاپی سے لنک نہیں کر سکتا۔ سیشنز تک شیئر کی گئی فہرست پہنچانے کا واحد طریقہ سرورز کو اس میں کاپی کرنا ہے۔

**ccam لکھنے سے پہلے پوچھتا ہے۔** ڈیفالٹ فلیگ والی حکمتِ عملی میں ccam وعدہ کرتا ہے کہ `.claude.json` کو نہیں چھیڑے گا۔ Remote Control واحد صورت ہے جس میں اسے وہاں لکھنا پڑتا ہے، اس لیے وہ پہلے پوچھتا ہے، بیک اپ رکھتا ہے، اور `mcpServers` کے سوا کچھ نہیں چھیڑتا۔

**اکاؤنٹ چلنے کے دوران یہ نہیں لکھتا۔** Claude Code سیشن کے دوران `.claude.json` دوبارہ لکھتا ہے، اور باہر سے کی گئی تبدیلی ضائع ہو جائے گی یا اس پر لکھ دیا جائے گا۔ اسی لیے جب اکاؤنٹ کا کوئی بھی سیشن چل رہا ہو، بشمول کوئی دوسرا Remote Control سرور، تو ccam کاپی چھوڑ دیتا ہے۔

**یہ ہر تبدیلی پر ایک بار پوچھتا ہے۔** کاپی کے بعد `sync-stamp` گروپ کی `mcp.json` کا ترمیم کا وقت رکھتا ہے۔ ccam دوبارہ صرف تب پوچھتا ہے جب وہ فائل بدلے، مثلاً `ccam mcp add` کے بعد۔ اگر آپ `N` جواب دیں تو ابھی کوئی کاپی نہیں ہوتی، اس لیے یہ ہر Remote Control لانچ پر پوچھتا ہے۔

**`ccam mcp` کاپیاں تازہ رکھتا ہے۔** نام ٹکرانے پر `.claude.json` کی اینٹری کو `--mcp-config` پر ترجیح ملتی ہے، اس لیے پرانی کاپی عام لانچز میں بھی نئی تعریف کو چھپا دے گی۔ اسی وجہ سے `ccam mcp add`، `add-json`، `remove` اور `import` ہر اس اکاؤنٹ کی کاپیاں اپ ڈیٹ کرتے ہیں جو نہیں چل رہا:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**ٹرمینل کے بغیر جواب "نہیں" ہوتا ہے۔** Remote Control اکثر کسی اسکرپٹ، سروس یا `nohup` سے چلتا ہے، جہاں جواب دینے والا کوئی نہیں ہوتا۔ تب ccam انتظار کرنے کے بجائے سرورز کے بغیر شروع کرتا ہے۔

## مسائل اور ان کے حل

### Remote Control کے سیشنز میں شیئر کیے گئے MCP سرورز موجود نہیں

سوال کا جواب `N` دیا گیا، ٹرمینل نہیں تھا، یا اکاؤنٹ چل رہا تھا۔ Remote Control بند کریں، اکاؤنٹ کے دوسرے سیشنز بند کریں، اور ٹرمینل میں `ccam work rc` دوبارہ چلائیں؛ جواب `y` دیں۔ کبھی نہ پوچھے جانے کے لیے `~/.config/ccam/config.toml` میں `mcp_strategy = "sync"` سیٹ کریں۔ کسی سیشن میں پہلے سے موجود سرورز اس کے چلتے ہوئے نہیں بدلتے، اس لیے کاپی کے بعد نیا سیشن شروع کریں۔

### "Account work is running, so its .claude.json is left alone"

`work` کا کوئی سیشن ابھی کھلا ہے: عام `ccam work`، کوئی شارٹ کٹ، یا کوئی دوسرا Remote Control سرور۔ اس سے باہر نکلیں اور `ccam work rc` دوبارہ چلائیں۔ کسی دوسرے اکاؤنٹ کے چلائے ہوئے سیشنز شمار نہیں ہوتے۔

### ccam ہر بار پوچھتا ہے

آپ نے ابھی تک سرورز کاپی نہیں کیے۔ ایک بار `y` جواب دیں، یا `mcp_strategy = "sync"` سیٹ کریں۔ اگر آپ Remote Control میں سرورز کبھی نہیں چاہتے تو اکاؤنٹ کے ساتھ MCP کی شیئرنگ بند کریں: `ccam share work --mcp off`۔

### `ccam work rc` فوراً `--mcp-config` سے متعلق ایرر کے ساتھ ناکام ہو جاتا ہے

آپ کا ccam، 0.6.0 سے پرانا ہے اور Remote Control کو `--mcp-config` دیتا ہے۔ `ccam update` چلائیں۔ موجودہ ccam کے ساتھ یہی ایرر آنے کا مطلب ہے کہ کوئی قدر والا آپشن `rc` سے پہلے آیا، جیسے `ccam work --model opus rc` میں، اس لیے ccam نے `opus` کو کمانڈ سمجھ لیا؛ `rc` کو اکاؤنٹ کے نام کے فوراً بعد رکھیں۔

### `ccam mcp` سے کی گئی تبدیلی Remote Control تک نہیں پہنچتی

`ccam mcp add` اور `remove` صرف ان اکاؤنٹس میں کاپیاں اپ ڈیٹ کرتے ہیں جو نہیں چل رہے، اور یہ بتاتے بھی ہیں: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control بند کر کے دوبارہ شروع کریں؛ ccam بدلی ہوئی `mcp.json` دیکھ کر دوبارہ پوچھے گا۔

### کاپی نے میرا اپنا ایک سرور بدل دیا

اگر اکاؤنٹ کا اپنا کوئی سرور گروپ کے کسی سرور کے ہم نام تھا تو کاپی نے اسے بدل دیا۔ پچھلی `.claude.json`، `~/.local/state/ccam/mcp/work/backups/` میں ہے، جو تازہ ترین پانچ بیک اپ رکھتا ہے۔ اپنا ورژن ہر اکاؤنٹ کے لیے رکھنے کے لیے، اگلے Remote Control لانچ سے پہلے اسے `ccam mcp import work <server> --force` سے گروپ میں منتقل کریں۔

### میں کاپی کیے گئے سرورز ہٹانا چاہتا ہوں

`ccam mcp remove <server>` سرور کو گروپ سے اور ہر کاپی سے ہٹا دیتا ہے۔ کاپی سے پہلے والی حالت پر واپس جانے کے لیے، جب اکاؤنٹ کا کوئی سیشن نہ چل رہا ہو، `~/.local/state/ccam/mcp/work/backups/` سے تازہ ترین بیک اپ بحال کریں۔

### Remote Control پوچھتا ہے کہ کیا فولڈر پر بھروسہ ہے، یا ٹرمینل کے بغیر بند ہو جاتا ہے

فولڈر پر بھروسہ ہر اکاؤنٹ کے لیے `.claude.json` میں الگ محفوظ ہوتا ہے، اس لیے جس فولڈر پر آپ نے `personal` کے طور پر بھروسہ کیا وہ `work` کے لیے نیا ہے۔ اس فولڈر میں ایک بار `ccam work` چلائیں اور قبول کریں، پھر وہاں Remote Control شروع کریں۔

### "You must be logged in with a Claude account that has a subscription"

Remote Control کو Claude سبسکرپشن والے لاگ اِن کی ضرورت ہے۔ `ccam list` میں اکاؤنٹ کی ای میل چیک کریں، اور ضرورت ہو تو `ccam login work` سے دوبارہ سائن اِن کریں۔

### سیشن فون پر نظر نہیں آتا

Claude ایپ یا claude.ai/code کسی دوسرے Claude اکاؤنٹ سے سائن اِن ہے۔ اس کا موازنہ اس ای میل سے کریں جو `ccam list`، `work` کے لیے دکھاتا ہے۔

### ccam خبردار کرتا ہے کہ اس فولڈر میں کوئی دوسرا Claude Code گفتگو کی ہسٹری شیئر کر رہا ہے

اسی ہسٹری گروپ کا کوئی دوسرا اکاؤنٹ اسی فولڈر میں چل رہا ہے۔ اگر دونوں اکاؤنٹس کے سیشنز ایک ہی گفتگو جاری رکھیں تو ہسٹری گڈمڈ ہو سکتی ہے۔ `y` صرف تب جواب دیں جب وہ الگ الگ گفتگو پر کام کریں گے، ورنہ Remote Control کسی دوسرے فولڈر میں شروع کریں۔

## متعلقہ

- [Claude Code شروع کرنا](/features/launch)
- [شیئر کیے گئے MCP سرورز](/features/mcp)
- [اکاؤنٹس کے درمیان شیئرنگ](/features/share)
- [اکاؤنٹس کا انتظام](/features/accounts)
