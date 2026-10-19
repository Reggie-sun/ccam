---
title: "CLAUDE_CONFIG_DIR کی وضاحت"
description: "Claude Code میں CLAUDE_CONFIG_DIR کیا کرتا ہے، فولڈر میں کیا رکھا جاتا ہے، macOS لاگ اِن کہاں رکھتا ہے، اور کون سی فائلیں اکاؤنٹس کے درمیان شیئر ہو سکتی ہیں۔"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR`، Claude Code کو بتاتا ہے کہ `~/.claude` کی جگہ کون سا فولڈر استعمال کرنا ہے۔ اکاؤنٹ سے متعلق ہر چیز وہیں رہتی ہے: لاگ اِن، سیٹنگز، MCP سرورز، پلگ اِنز اور گفتگو کی ہسٹری۔ اسے مختلف فولڈرز کی طرف سیٹ کرنے سے آپ کو الگ الگ اکاؤنٹس مل جاتے ہیں۔

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## کیا CLAUDE_CONFIG_DIR سرکاری طور پر سپورٹڈ ہے؟

ہاں۔ Claude Code اسے اپنے [انوائرنمنٹ ویری ایبلز](https://code.claude.com/docs/en/env-vars) میں کنفیگریشن فولڈر بدلنے کے طریقے کے طور پر درج کرتا ہے، "کئی اکاؤنٹس ساتھ ساتھ چلانے کے لیے مفید"، اور اس کی [تصدیق (authentication) کی گائیڈ](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) کئی اکاؤنٹس سے لاگ اِن کرنے کے لیے اسی کو استعمال کرتی ہے۔ آپ اسے اپنے شیل میں یا یوزر سیٹنگز میں سیٹ کر سکتے ہیں، لیکن کسی پروجیکٹ کی سیٹنگز میں نہیں۔

## فولڈر کے اندر کیا ہوتا ہے؟

Claude Code 2.1.302 کے ساتھ جانچا گیا:

| راستہ | مواد | شیئر ہو سکتا ہے |
| --- | --- | --- |
| `.credentials.json` | لاگ اِن ٹوکنز (Linux، Windows)، MCP OAuth ٹوکنز | نہیں |
| `.claude.json` | اکاؤنٹ کی تفصیلات، یوزر MCP سرورز، ہر پروجیکٹ کی حالت | نہیں (MCP سرورز الگ سے) |
| `settings.json` | اجازتیں، ہُکس، ماڈل، اسٹیٹس لائن، پلگ اِنز | ہاں |
| `CLAUDE.md` | گلوبل ہدایات | ہاں |
| `agents/`، `commands/`، `skills/`، `plugins/` | ایکسٹینشنز | ہاں |
| `projects/` | `.jsonl` فائلوں کی صورت میں گفتگو، ہر پروجیکٹ کی میموری | ہاں |
| `history.jsonl` | پرامپٹس کی ہسٹری | ہاں |
| `sessions/`، `shell-snapshots/`، `cache/` اور دیگر | رن ٹائم حالت | نہیں |

## لوگ کن باتوں میں اکثر الجھ جاتے ہیں؟

**`CLAUDE_CONFIG_DIR` کے بغیر، `.claude.json` فائل `~/.claude` کے اندر نہیں ہوتی۔** اس صورت میں Claude Code اسے فولڈر کے برابر میں `~/.claude.json` پر رکھتا ہے۔ اگر آپ `CLAUDE_CONFIG_DIR=~/.claude` استعمال کرنا شروع کریں تو Claude Code اس کے بجائے `~/.claude/.claude.json` تلاش کرتا ہے اور آپ کا اکاؤنٹ نہیں ڈھونڈ پاتا۔ `ccam import` ان دونوں کو لنک کر کے یہ مسئلہ حل کرتا ہے۔

**macOS پر لاگ اِن فولڈر میں نہیں بلکہ Keychain میں ہوتا ہے۔** Claude Code اسے `Claude Code-credentials` نام کے تحت محفوظ کرتا ہے، اور جب `CLAUDE_CONFIG_DIR` سیٹ ہو تو اس نام کے بعد فولڈر کے راستے کا ایک مختصر ہیش لگتا ہے۔ اس لیے الگ الگ فولڈرز کے لیے Keychain میں الگ الگ اندراجات بنتے ہیں۔ اگر Keychain لکھنے کی اجازت نہ دے تو Claude Code متبادل کے طور پر فولڈر میں `.credentials.json` استعمال کرتا ہے، جیسا کہ [اس کی دستاویزات](https://code.claude.com/docs/en/authentication#credential-management) میں بیان کیا گیا ہے۔

**ہر لاگ اِن فولڈر میں نہیں ہوتا۔** API key کے بغیر Claude Console کا لاگ اِن فولڈر سے باہر محفوظ ہوتا ہے، اس لیے الگ فولڈرز ایسے دو لاگ اِنز کو الگ نہیں رکھ پاتے۔ claude.ai کے لاگ اِنز اور API keys ہر فولڈر میں الگ رکھے جاتے ہیں۔

## Windows پر فولڈر کہاں ہوتا ہے؟

ڈیفالٹ طور پر `%USERPROFILE%\.claude`، اور لاگ اِن اس کے اندر `.credentials.json` میں ہوتا ہے۔ ہر اضافی اکاؤنٹ کے لیے `CLAUDE_CONFIG_DIR` کو کسی دوسرے فولڈر پر سیٹ کریں، بالکل اسی طرح جیسے Linux اور macOS پر۔

## کیا فولڈرز فائلیں شیئر کر سکتے ہیں؟

سمبولک لنکس کام کرتے ہیں: Claude Code لنک کیے گئے `settings.json` یا `projects/` فولڈر کے ذریعے پڑھتا اور لکھتا ہے اور لنک کو برقرار رکھتا ہے۔ MCP سرورز اس سے مستثنیٰ ہیں، کیونکہ وہ لاگ اِن کے ساتھ `.claude.json` میں رہتے ہیں؛ انہیں اس کے بجائے `--mcp-config` سے دیں۔

ccam بالکل انہی اصولوں پر بنا ہے۔ یہ فولڈرز کو کیسے لنک کرتا ہے، یہ جاننے کے لیے [شیئرنگ](/docs/sharing) دیکھیں، اور مجموعی تصویر کے لیے [کئی اکاؤنٹس کی گائیڈ](/guides/multiple-claude-code-accounts) پڑھیں۔
