---
title: "کمانڈز کا حوالہ"
description: "ccam کی ہر کمانڈ اور آپشن، ccam add اور ccam <name> سے لے کر ccam usage، ccam next اور ccam doctor تک۔"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a5a33b1a2f4
---

کوئی کمانڈ اصل میں چلائے بغیر یہ دیکھنے کے لیے کہ وہ کیا کرے گی، اس سے پہلے `--dry-run` لگائیں، مثلاً `ccam --dry-run work`۔

<!-- commands -->

## Claude Code شروع کرنا

`ccam <name|number> [arguments]` اس اکاؤنٹ کے طور پر Claude Code شروع کرتا ہے۔ ccam نام کے بعد آنے والے آرگیومنٹس نہیں پڑھتا؛ وہ انہیں اسی ترتیب میں `claude` کو دے دیتا ہے، اور اسٹینڈرڈ اِن پٹ، ایگزٹ کوڈ اور سگنلز بالکل ویسے ہی کام کرتے ہیں جیسے خود `claude` کے ساتھ۔ نام کے بعد `--` لگانا اختیاری ہے۔

اکاؤنٹ کے نام کے بغیر لیکن آرگیومنٹس کے ساتھ `ccam`، جیسے `ccam -p "summarise"`، ڈیفالٹ اکاؤنٹ استعمال کرتا ہے۔ ڈیفالٹ اکاؤنٹ `ccam use <name>` سے مقرر کریں۔

`ccwork` جیسا شارٹ کٹ اپنے اکاؤنٹ کو اسی طرح شروع کرتا ہے: `ccwork --resume` وہی کرتا ہے جو `ccam work --resume` کرتا ہے۔ [شارٹ کٹس](/docs/commands#shortcuts) دیکھیں۔

## اکاؤنٹس شامل کرنا اور ہٹانا

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

فولڈر بتائے بغیر `ccam import`، `~/.claude` اور `~/.claude-*` تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے۔ امپورٹ کیے گئے فولڈرز اپنی جگہ پر ہی رجسٹر ہوتے ہیں؛ ccam کبھی کوئی لاگ اِن منتقل یا کاپی نہیں کرتا۔

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
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ہر آئٹم کے لیے `on` (یعنی `default` گروپ)، `off`، یا کسی گروپ کا نام دیا جا سکتا ہے۔ [شیئرنگ](/docs/sharing) دیکھیں۔

## موجودہ سیٹ اپ کاپی کرنا

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` کسی Claude Code فولڈر سے `settings.json`، `CLAUDE.md`، ایکسٹینشنز اور یوزر MCP سرورز کو ایک گروپ میں کاپی کرتا ہے، تاکہ وہ آئٹمز شیئر کرنے والے ہر اکاؤنٹ کو یہ مل جائیں۔ فولڈر کو صرف پڑھا جاتا ہے۔ آپشنز `--settings`، `--memory`، `--extensions`، `--mcp` اور `--history` صرف وہی آئٹمز کاپی کرتے ہیں جن کا آپ نام لیں؛ ان کے بغیر گفتگو کی ہسٹری کے سوا سب کچھ کاپی ہوتا ہے۔ `-g` یا `--group` گروپ چنتا ہے، اور `-y` یا `--yes` بغیر پوچھے کاپی کرتا ہے اور فائل مختلف ہونے پر گروپ والی کاپی رکھتا ہے۔ [موجودہ سیٹ اپ کاپی کرنا](/docs/sharing#copy-an-existing-setup) دیکھیں۔

## MCP سرورز

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## استعمال اور اکاؤنٹ تبدیل کرنا

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## دیکھ بھال

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
```

`ccam update`، ccam کو بدلنے سے پہلے ریلیز کا SHA-256 چیک سم جانچتا ہے، اور `--version v0.2.0` کوئی مخصوص ریلیز انسٹال کرتا ہے۔ `.deb` پیکیج سے انسٹال کیا گیا ccam اس کے بجائے apt سے اپ ڈیٹ ہوتا ہے۔ [اپ ڈیٹ کرنا](/docs/install#update) دیکھیں۔

## فائلیں

| راستہ | مواد |
| --- | --- |
| `~/.config/ccam/config.toml` | اکاؤنٹس اور یہ کہ وہ کیا شیئر کرتے ہیں |
| `~/.local/share/ccam/profiles/<name>/` | ہر اکاؤنٹ کا `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | شیئر کی گئی سیٹنگز، MCP سرورز، ایکسٹینشنز اور ہسٹری |
| `~/.local/state/ccam/` | استعمال کا کیش اور آخری استعمال کے اوقات |
| `~/.local/bin/<shortcut>` | `ccwork` جیسی شارٹ کٹ کمانڈز، جو ccam کے لنکس ہیں |

`XDG_CONFIG_HOME`، `XDG_DATA_HOME` اور `XDG_STATE_HOME` ویری ایبلز ان فولڈرز کی جگہ بدل دیتے ہیں۔
