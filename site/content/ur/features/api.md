---
title: "API key اکاؤنٹس"
description: "ccam add --api، Anthropic کے API یا کسی گیٹ وے کے لیے API key اکاؤنٹ شامل کرتا ہے؛ key لانچ پر کسی ویری ایبل یا کمانڈ سے پڑھی جاتی ہے، کبھی محفوظ نہیں ہوتی۔"
updated: 2026-10-26
order: 13
translation: machine
sourceHash: 771d749f1f90
---

API key اکاؤنٹ، Claude سبسکرپشن لاگ اِن کے بجائے API key سے Claude Code شروع کرتا ہے۔ key، Anthropic کے API کے لیے ہو سکتی ہے یا کسی بھی ایسے اینڈ پوائنٹ کے لیے جو وہی Messages فارمیٹ سمجھتا ہو، جیسے LiteLLM یا OpenRouter گیٹ وے۔ ccam، key محفوظ نہیں کرتا۔ یہ ہر بار اکاؤنٹ شروع ہونے پر اسے کسی انوائرنمنٹ ویری ایبل سے، یا پاس ورڈ مینیجر جیسی کسی کمانڈ سے پڑھتا ہے، اور صرف اسی Claude Code پروسیس کو دیتا ہے۔

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
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

API key اکاؤنٹ کسی بھی دوسرے اکاؤنٹ کی طرح گروپس میں شامل ہوتا ہے (ڈیفالٹ طور پر `default` گروپ، جو سیٹنگز، `CLAUDE.md`، ایکسٹینشنز اور MCP سرورز شیئر کرتا ہے)، اور شارٹ کٹس، `ccam folder` اور `ccam ext` کے ساتھ کام کرتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--api` | لاگ اِن کے بجائے API key اکاؤنٹ شامل کرتا ہے | بند |
| `--key-from env:NAME` | key انوائرنمنٹ ویری ایبل `NAME` سے پڑھتا ہے | |
| `--key-from command:COMMAND` | `COMMAND` چلاتا ہے (`sh -c` سے، یا Windows پر `cmd /d /s /c "…"` میں بالکل ویسے ہی جیسے لکھی گئی ہو، کوٹیشن مارکس برقرار رکھتے ہوئے) اور اس کی پرنٹ کی ہوئی پہلی لائن استعمال کرتا ہے | |
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

1. ccam، key پڑھتا ہے: `env:` انوائرنمنٹ سے، `command:` کمانڈ کو 30 سیکنڈ کے ٹائم آؤٹ کے ساتھ چلا کر۔ کمانڈ ٹرمینل پر پاس ورڈ مانگ سکتی ہے؛ اس کا ایرر آؤٹ پٹ بھی ٹرمینل پر جاتا ہے۔ خالی key لانچ کو ایک پیغام کے ساتھ روک دیتی ہے جو بتاتا ہے کہ key کہاں متوقع تھی۔
2. یہ ہر وہ وراثتی ویری ایبل ہٹا دیتا ہے جو طے کرتا ہے کہ Claude Code کیسے سائن اِن کرے یا کون سا پرووائیڈر استعمال کرے:
   - keys اور ٹوکنز: `ANTHROPIC_API_KEY`، `ANTHROPIC_AUTH_TOKEN`، `CLAUDE_CODE_OAUTH_TOKEN`، `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`، `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`؛
   - اینڈ پوائنٹ: `ANTHROPIC_BASE_URL`، `ANTHROPIC_UNIX_SOCKET`، `ANTHROPIC_MODEL`؛
   - پرووائیڈر سوئچز: `CLAUDE_CODE_USE_BEDROCK`، `CLAUDE_CODE_USE_VERTEX`، `CLAUDE_CODE_USE_FOUNDRY`، `CLAUDE_CODE_USE_ANTHROPIC_AWS`، `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`، `CLAUDE_CODE_USE_MANTLE`، `CLAUDE_CODE_USE_GATEWAY`۔

   آپ کے شیل میں کسی اور کام کے لیے سیٹ کی گئی key اس اکاؤنٹ تک نہیں پہنچ سکتی، اور سبسکرپشن لاگ اِن اس کی key پر ترجیح نہیں لے سکتا۔
3. یہ دیے جانے پر `ANTHROPIC_BASE_URL` اور `ANTHROPIC_MODEL` سیٹ کرتا ہے، اور key کو `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN` کے طور پر۔

سبسکرپشن لاگ اِن والے اکاؤنٹس پہلے کی طرح وراثت میں ملا انوائرنمنٹ رکھتے ہیں۔ انٹرایکٹو سیشن میں Claude Code ایک بار پوچھ سکتا ہے کہ کیا انوائرنمنٹ میں ملنے والی key استعمال کی جائے۔

## حفاظت اور خاص صورتیں

- key کبھی `config.toml`، کسی کیش یا لاگ میں نہیں لکھی جاتی، اور `ccam --dry-run gw` اسے نہیں پڑھتا: ڈرائی رن `command:` ماخذ نہیں چلاتا۔ اس کے بجائے ڈرائی رن اینڈ پوائنٹ، ماڈل اور یہ دکھاتا ہے کہ key کہاں سے آتی ہے:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  لانچ جو وراثتی ویری ایبل ہٹاتا ہے وہ شروع میں `env -u <name>` کے طور پر دکھایا جاتا ہے۔ اگر `env:` ماخذ کا ویری ایبل خالی ہو تو ڈرائی رن اسی طرح ناکام ہوتا ہے جیسے لانچ ہوتا۔
- `ccam add --api`، `env:` ماخذ کو فوراً جانچتا ہے اور ویری ایبل خالی ہو تو انتباہ دیتا ہے؛ `command:` ماخذ لانچ تک نہیں چلتا۔
- `ccam login` اور `ccam logout`، API key اکاؤنٹ کو رد کر دیتے ہیں، کیونکہ اس کا کوئی لاگ اِن نہیں۔ `ccam doctor` اس کا اینڈ پوائنٹ اور key کا ماخذ دکھاتا ہے، اور جب کوئی `env:` ویری ایبل خالی ہو تو انتباہ دیتا ہے۔
- `ccam usage` اس کے لیے `no limits (API)` دکھاتا ہے۔ `ccam next --auto`، `--fallback` اور `ccam pick` صرف سبسکرپشن اکاؤنٹس میں سے چنتے ہیں، کیونکہ API key کی کوئی 5 گھنٹے یا ہفتہ وار حد نہیں جس کا موازنہ کیا جا سکے۔
- `env` بلاک یا `apiKeyHelper` والی شیئر شدہ `settings.json` گروپ کے ہر اکاؤنٹ پر لاگو ہوتی ہے، API key اکاؤنٹس سمیت۔ پرووائیڈر کی سیٹنگز کو شیئر شدہ سیٹنگز سے باہر رکھیں، یا API key اکاؤنٹ کو ایسے گروپ میں رکھیں جس میں `settings` نہ ہو۔

## متعلقہ

- [اکاؤنٹس کا انتظام](/features/accounts)
- [Claude Code شروع کرنا](/features/launch)
- [کمانڈز کا حوالہ](/docs/commands)
