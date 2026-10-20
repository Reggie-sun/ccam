---
title: "ایکسٹینشنز کی جانچ"
description: "ccam ext وہ اسکلز، کمانڈز، پلگ اِنز، MCP سرورز اور ہُکس دکھاتا ہے جو ہر اکاؤنٹ لوڈ کرتا ہے، ہر ایک کہاں سے آتا ہے، اور کیا خراب ہے۔"
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` اس سوال کا جواب دیتا ہے کہ "یہ اسکل (یا MCP سرور) ایک اکاؤنٹ میں کیوں کام کرتی ہے اور دوسرے میں کیوں نہیں؟"۔ یہ ہر اکاؤنٹ کے لیے وہ ایکسٹینشنز دکھاتا ہے جو Claude Code لوڈ کرتا ہے، اور یہ کہ ہر ایک کہاں سے آتی ہے: کوئی شیئر گروپ، اکاؤنٹ کا اپنا فولڈر، یا وہ پروجیکٹ جس میں آپ ہیں۔ ٹیبل کے نیچے یہ وہ چیزیں دکھاتا ہے جو کسی ایکسٹینشن کو آپ کی توقع کے مطابق لوڈ ہونے سے روکتی ہیں۔ یہ صرف پڑھتا ہے۔

## استعمال

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

ہر خانہ اس اکاؤنٹ میں اس آئٹم کا ماخذ بتاتا ہے: شیئر گروپ (`default`)، `account` جب وہ اکاؤنٹ کے اپنے فولڈر میں ہو، یا `-` جب اکاؤنٹ کے پاس وہ نہ ہو۔ `!` کسی مسئلے والے آئٹم کو نشان زد کرتا ہے، جس کی وضاحت ٹیبل کے نیچے ہوتی ہے۔

## یہ کیا پڑھتا ہے

| قسم | کہاں | دکھایا گیا ماخذ |
| --- | --- | --- |
| اسکلز | `skills/` میں ہر وہ فولڈر جس میں `SKILL.md` ہو | وہ گروپ جس سے `skills/` لنک ہے، یا `account` |
| کمانڈز، ایجنٹس | `commands/` اور `agents/` کے نیچے `.md` فائلیں، ذیلی فولڈرز `folder/name` کی صورت میں | گروپ، یا `account` |
| پلگ اِنز | `plugins/installed_plugins.json` | گروپ، یا `account` |
| MCP سرورز | گروپ کی `mcp.json` (جب MCP شیئر ہو) اور اکاؤنٹ کی `.claude.json` | گروپ، یا `account` |
| ہُکس | `settings.json` میں `hooks`، ہر ایونٹ کے حساب سے گنے گئے | گروپ، یا `account` |
| یہ پروجیکٹ | موجودہ فولڈر میں `.claude/skills`، `.claude/commands`، `.claude/agents`، `.claude/settings.json`، `.claude/settings.local.json` اور `.mcp.json` | ایک بار، ہر اکاؤنٹ کے لیے دکھایا جاتا ہے |

## یہ کون سے مسائل بتاتا ہے

- **ٹوٹے ہوئے لنکس**: `skills/`، `commands/` یا `agents/` میں کوئی لنک جس کا ہدف موجود نہ رہا ہو، یا کوئی شیئر شدہ فولڈر (`skills`، `commands`، `agents`، `plugins`) جو کسی غائب جگہ سے لنک ہو (اسے `ccam doctor --fix` سے ٹھیک کریں)۔
- **`SKILL.md` کے بغیر اسکل فولڈر**: Claude Code اسے چھوڑ دیتا ہے۔
- **دو جگہوں پر مختلف طرح سے متعین MCP سرور**: گروپ کی `mcp.json` اور اکاؤنٹ کی `.claude.json`، یا ان میں سے ایک اور پروجیکٹ کی `.mcp.json`۔ دونوں جگہوں پر ایک جیسی سیٹنگز (جو `mcp_strategy = "sync"` لکھتا ہے) مسئلہ نہیں۔ کسی اکاؤنٹ کا اپنا سرور `ccam mcp import <account> <server>` سے گروپ میں منتقل کریں۔
- **`installed_plugins.json` جو درست JSON نہ ہو۔**

## حفاظت اور خاص صورتیں

- `ccam ext` کچھ نہیں بدلتا اور کوئی نیٹ ورک درخواست نہیں بھیجتا۔ یہ سرورز کے نام اور ہُکس کی تعداد پڑھتا ہے، MCP تعریفوں میں `env` یا ہیڈرز کی قدریں کبھی نہیں۔
- یہ وہی دکھاتا ہے جو ڈسک پر ہے۔ جو پلگ اِن انسٹال ہو لیکن سیٹنگز میں بند ہو، یا جو سرور منیجڈ سیٹنگز نے روکا ہو، وہ پھر بھی دکھائی دیتا ہے۔
- اسے پروجیکٹ کے فولڈر میں چلائیں تاکہ اس پروجیکٹ کا `.claude` فولڈر اور `.mcp.json` بھی شامل ہوں۔

## متعلقہ

- [اکاؤنٹس کے درمیان شیئرنگ](/features/share)
- [شیئر کیے گئے MCP سرورز](/features/mcp)
- [جانچ اور مرمت](/features/doctor)
