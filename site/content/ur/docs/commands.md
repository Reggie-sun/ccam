---
title: "کمانڈز کا حوالہ"
description: "ccam کی ہر کمانڈ اور آپشن، ccam add اور ccs <name> سے لے کر ccam usage، ccam next اور ccam doctor تک۔"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` اور `ccs` ایک ہی پروگرام ہیں؛ جو آپ کو پسند ہو وہ استعمال کریں۔ کوئی کمانڈ اصل میں چلائے بغیر یہ دیکھنے کے لیے کہ وہ کیا کرے گی، اس سے پہلے `--dry-run` لگائیں، مثلاً `ccs --dry-run work`۔

<!-- commands -->

## Claude Code شروع کرنا

`ccs <name|number> [arguments]` اس اکاؤنٹ کے طور پر Claude Code شروع کرتا ہے۔ ccam نام کے بعد آنے والے آرگیومنٹس نہیں پڑھتا؛ وہ انہیں اسی ترتیب میں `claude` کو دے دیتا ہے، اور اسٹینڈرڈ اِن پٹ، ایگزٹ کوڈ اور سگنلز بالکل ویسے ہی کام کرتے ہیں جیسے خود `claude` کے ساتھ۔ نام کے بعد `--` لگانا اختیاری ہے۔

اکاؤنٹ کے نام کے بغیر لیکن آرگیومنٹس کے ساتھ `ccs`، جیسے `ccs -p "summarise"`، ڈیفالٹ اکاؤنٹ استعمال کرتا ہے۔ ڈیفالٹ اکاؤنٹ `ccam use <name>` سے مقرر کریں۔

شارٹ کٹس `claude1`، `claude2` وغیرہ `~/.local/bin` میں ccam کے لنکس ہیں۔ `ccam shims install --names` ہر اکاؤنٹ کے لیے `claude-<name>` بھی بناتا ہے۔

## اکاؤنٹس شامل کرنا اور ہٹانا

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

فولڈر بتائے بغیر `ccam import`، `~/.claude` اور `~/.claude-*` تلاش کرتا ہے اور ہر ایک کے بارے میں پوچھتا ہے۔ امپورٹ کیے گئے فولڈرز اپنی جگہ پر ہی رجسٹر ہوتے ہیں؛ ccam کبھی کوئی لاگ اِن منتقل یا کاپی نہیں کرتا۔

## شیئرنگ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ہر آئٹم کے لیے `on` (یعنی `default` گروپ)، `off`، یا کسی گروپ کا نام دیا جا سکتا ہے۔ [شیئرنگ](/docs/sharing) دیکھیں۔

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
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## دیکھ بھال

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## فائلیں

| راستہ | مواد |
| --- | --- |
| `~/.config/ccam/config.toml` | اکاؤنٹس اور یہ کہ وہ کیا شیئر کرتے ہیں |
| `~/.local/share/ccam/profiles/<name>/` | ہر اکاؤنٹ کا `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | شیئر کی گئی سیٹنگز، MCP سرورز، ایکسٹینشنز اور ہسٹری |
| `~/.local/state/ccam/` | استعمال کا کیش اور آخری استعمال کے اوقات |
| `~/.local/bin/claude1`، `ccs` | ccam کے لنکس |

`XDG_CONFIG_HOME`، `XDG_DATA_HOME` اور `XDG_STATE_HOME` ویری ایبلز ان فولڈرز کی جگہ بدل دیتے ہیں۔
