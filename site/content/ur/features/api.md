---
title: "API key اکاؤنٹس"
description: "ccam add --api، Anthropic کے API یا کسی گیٹ وے کے لیے API key اکاؤنٹ شامل کرتا ہے؛ key لانچ پر کسی ویری ایبل یا کمانڈ سے پڑھی جاتی ہے، کبھی محفوظ نہیں ہوتی۔"
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API key اکاؤنٹ، Claude سبسکرپشن لاگ اِن کے بجائے API key سے Claude Code شروع کرتا ہے: Anthropic کا API، یا کوئی ایسا اینڈ پوائنٹ جو وہی Messages فارمیٹ سمجھتا ہو، جیسے LiteLLM یا OpenRouter گیٹ وے۔ ccam، key محفوظ نہیں کرتا۔ یہ ہر بار اکاؤنٹ شروع ہونے پر اسے کسی انوائرنمنٹ ویری ایبل سے، یا پاس ورڈ مینیجر جیسی کسی کمانڈ سے پڑھتا ہے، اور صرف اسی Claude Code پروسیس کو دیتا ہے۔

## استعمال

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

API key اکاؤنٹ کسی بھی دوسرے اکاؤنٹ کی طرح سیٹنگز، `CLAUDE.md`، ایکسٹینشنز اور MCP سرورز شیئر کرتا ہے، اور شارٹ کٹس، `ccam folder` اور `ccam ext` کے ساتھ کام کرتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--api` | لاگ اِن کے بجائے API key اکاؤنٹ شامل کرتا ہے | بند |
| `--key-from env:NAME` | key انوائرنمنٹ ویری ایبل `NAME` سے پڑھتا ہے | |
| `--key-from command:COMMAND` | `COMMAND` چلاتا ہے (`sh -c` سے، یا Windows پر `cmd /C` سے) اور اس کی پرنٹ کی ہوئی پہلی لائن استعمال کرتا ہے | |
| `--base-url <url>` | اینڈ پوائنٹ (`ANTHROPIC_BASE_URL`) | Anthropic کا API |
| `--auth api-key\|bearer` | key کو `ANTHROPIC_API_KEY` (`x-api-key`) یا `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`، جو گیٹ ویز عموماً چاہتے ہیں) کے طور پر بھیجتا ہے | `api-key` |
| `--model <model>` | ماڈل (`ANTHROPIC_MODEL`) | Claude Code کا ڈیفالٹ |

سیٹنگز `config.toml` میں محفوظ ہوتی ہیں:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## یہ کیسے کام کرتا ہے

اکاؤنٹ کے ہر لانچ پر، معمول کی جانچ کے بعد اور Claude Code شروع ہونے سے پہلے:

1. ccam، key پڑھتا ہے: `env:` انوائرنمنٹ سے، `command:` کمانڈ کو 30 سیکنڈ کی حد کے ساتھ چلا کر۔ کمانڈ ٹرمینل پر پاس ورڈ مانگ سکتی ہے؛ اس کا ایرر آؤٹ پٹ بھی ٹرمینل پر جاتا ہے۔ خالی key لانچ کو ایک پیغام کے ساتھ روک دیتی ہے جو بتاتا ہے کہ key کہاں متوقع تھی۔
2. یہ ہر وہ وراثتی ویری ایبل ہٹا دیتا ہے جو طے کرتا ہے کہ Claude Code کیسے سائن اِن کرے یا کون سا پرووائیڈر استعمال کرے: `ANTHROPIC_API_KEY`، `ANTHROPIC_AUTH_TOKEN`، `CLAUDE_CODE_OAUTH_TOKEN`، `ANTHROPIC_BASE_URL`، `ANTHROPIC_MODEL`، `CLAUDE_CODE_USE_BEDROCK`، `CLAUDE_CODE_USE_VERTEX` اور `CLAUDE_CODE_USE_FOUNDRY`۔ آپ کے شیل میں کسی اور کام کے لیے سیٹ کی گئی key اس اکاؤنٹ تک نہیں پہنچ سکتی، اور سبسکرپشن لاگ اِن اس کی key پر غالب نہیں آ سکتا۔
3. یہ دیے جانے پر `ANTHROPIC_BASE_URL` اور `ANTHROPIC_MODEL` سیٹ کرتا ہے، اور key کو `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` کے طور پر۔

سبسکرپشن لاگ اِن والے اکاؤنٹس پہلے کی طرح وراثت میں ملا انوائرنمنٹ رکھتے ہیں۔ انٹرایکٹو سیشن میں Claude Code ایک بار پوچھ سکتا ہے کہ کیا انوائرنمنٹ میں ملنے والی key استعمال کی جائے۔

## حفاظت اور خاص صورتیں

- key کبھی `config.toml`، کسی کیش یا لاگ میں نہیں لکھی جاتی، اور `ccam --dry-run gw` اسے نہیں پڑھتا: ڈرائی رن `command:` ماخذ نہیں چلاتا۔
- `ccam add --api`، `env:` ماخذ کو فوراً جانچتا ہے اور ویری ایبل خالی ہو تو انتباہ دیتا ہے؛ `command:` ماخذ پہلی بار لانچ پر چلتا ہے۔
- `ccam login` اور `ccam logout`، API key اکاؤنٹ کو رد کر دیتے ہیں، کیونکہ اس کا کوئی لاگ اِن نہیں۔ `ccam doctor` اس کا اینڈ پوائنٹ اور key کا ماخذ دکھاتا ہے، اور جب کوئی `env:` ویری ایبل خالی ہو تو انتباہ دیتا ہے۔
- `ccam usage` اس کے لیے `API` دکھاتا ہے۔ `ccam next --auto`، `--fallback` اور `ccam pick` صرف سبسکرپشن اکاؤنٹس میں سے چنتے ہیں، کیونکہ API key کی کوئی 5 گھنٹے یا ہفتہ وار حد نہیں جس کا موازنہ کیا جا سکے۔
- `env` بلاک یا `apiKeyHelper` والی شیئر شدہ `settings.json` گروپ کے ہر اکاؤنٹ پر لاگو ہوتی ہے، API key اکاؤنٹس سمیت۔ پرووائیڈر کی سیٹنگز کو شیئر شدہ سیٹنگز سے باہر رکھیں۔

## متعلقہ

- [اکاؤنٹس کا انتظام](/features/accounts)
- [Claude Code شروع کرنا](/features/launch)
- [کمانڈز کا حوالہ](/docs/commands)
