---
title: حساب‌های کلید API
description: ccam add --api یک حساب کلید API برای API شرکت Anthropic یا یک درگاه اضافه می‌کند؛ کلید هنگام اجرا از یک متغیر یا فرمان خوانده و هرگز ذخیره نمی‌شود.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: "a903c92f2c50"
---

حساب کلید API، Claude Code را به جای ورود با اشتراک Claude با یک کلید API آغاز می‌کند: API شرکت Anthropic، یا نقطه پایانی‌ای که همان قالب Messages را می‌فهمد، مانند یک درگاه LiteLLM یا OpenRouter. ccam کلید را ذخیره نمی‌کند. هر بار که حساب آغاز می‌شود آن را از یک متغیر محیطی یا از فرمانی مانند یک مدیر گذرواژه می‌خواند و فقط به همان فرایند Claude Code می‌دهد.

## نحوه استفاده

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

حساب کلید API مانند هر حساب دیگری به گروه‌ها می‌پیوندد (به‌طور پیش‌فرض `default`: تنظیمات، `CLAUDE.md`، افزونه‌ها و سرورهای MCP) و با میان‌برها، `ccam folder` و `ccam ext` کار می‌کند.

## گزینه‌ها

| گزینه | کاری که انجام می‌دهد | پیش‌فرض |
| --- | --- | --- |
| `--api` | به جای ورود، یک حساب کلید API اضافه می‌کند | خاموش |
| `--key-from env:NAME` | کلید را از متغیر محیطی `NAME` می‌خواند | |
| `--key-from command:COMMAND` | `COMMAND` را اجرا می‌کند (با `sh -c`، یا در Windows با `cmd /C`) و نخستین خطی را که چاپ می‌کند به کار می‌برد | |
| `--base-url <url>` | نقطه پایانی (`ANTHROPIC_BASE_URL`) | API شرکت Anthropic |
| `--auth api-key\|bearer` | کلید را به‌عنوان `ANTHROPIC_API_KEY` (`x-api-key`) یا `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`، که درگاه‌ها معمولاً می‌خواهند) می‌فرستد | `api-key` |
| `--model <model>` | مدل (`ANTHROPIC_MODEL`) | پیش‌فرض Claude Code |

تنظیمات در `config.toml` ذخیره می‌شوند:

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

## سازوکار

در هر اجرای حساب، پس از بررسی‌های معمول و پیش از آغاز Claude Code:

1. ccam کلید را می‌خواند: `env:` از محیط، و `command:` با اجرای فرمان با محدودیت ۳۰ ثانیه. فرمان می‌تواند در ترمینال گذرواژه بخواهد؛ خروجی خطای آن هم به ترمینال می‌رود. کلید خالی اجرا را با پیامی متوقف می‌کند که می‌گوید کلید کجا انتظار می‌رفت.
2. هر متغیر به‌ارث‌رسیده‌ای را که تعیین می‌کند Claude Code چگونه وارد شود یا از کدام ارائه‌دهنده استفاده کند حذف می‌کند: کلیدها و توکن‌ها (`ANTHROPIC_API_KEY`، `ANTHROPIC_AUTH_TOKEN`، `CLAUDE_CODE_OAUTH_TOKEN`، `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`، `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`)، نقطه پایانی (`ANTHROPIC_BASE_URL`، `ANTHROPIC_UNIX_SOCKET`، `ANTHROPIC_MODEL`) و سوئیچ‌های انتخاب ارائه‌دهنده (`CLAUDE_CODE_USE_BEDROCK`، `CLAUDE_CODE_USE_VERTEX`، `CLAUDE_CODE_USE_FOUNDRY`، `CLAUDE_CODE_USE_ANTHROPIC_AWS`، `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`، `CLAUDE_CODE_USE_MANTLE`، `CLAUDE_CODE_USE_GATEWAY`). کلیدی که در شل شما برای کار دیگری تنظیم شده نمی‌تواند به این حساب برسد، و ورود با اشتراک نمی‌تواند بر کلید آن غلبه کند.
3. `ANTHROPIC_BASE_URL` و `ANTHROPIC_MODEL` را در صورت داده شدن تنظیم می‌کند، و کلید را به‌عنوان `ANTHROPIC_API_KEY` یا `ANTHROPIC_AUTH_TOKEN`.

حساب‌هایی که ورود با اشتراک دارند، مانند قبل، محیطی را که به ارث می‌برند نگه می‌دارند. در یک نشست تعاملی، Claude Code ممکن است یک بار بپرسد که آیا از کلیدی که در محیط پیدا می‌کند استفاده کند یا نه.

## ایمنی و موارد خاص

- کلید هرگز در `config.toml`، حافظه نهان یا گزارش نوشته نمی‌شود، و `ccam --dry-run gw` آن را نمی‌خواند: اجرای آزمایشی منبع `command:` را اجرا نمی‌کند.
- `ccam add --api` منبع `env:` را بی‌درنگ بررسی می‌کند و اگر متغیر خالی باشد هشدار می‌دهد؛ منبع `command:` نخستین بار هنگام اجرا اجرا می‌شود.
- `ccam login` و `ccam logout` حساب کلید API را رد می‌کنند، چون ورودی ندارد. `ccam doctor` نقطه پایانی و منبع کلید آن را نشان می‌دهد و وقتی یک متغیر `env:` خالی باشد هشدار می‌دهد.
- `ccam usage` برای آن `API` نشان می‌دهد. `ccam next --auto`، `--fallback` و `ccam pick` فقط از میان حساب‌های اشتراکی انتخاب می‌کنند، چون کلید API سقف ۵ ساعته یا هفتگی‌ای برای مقایسه ندارد.
- یک `settings.json` مشترک که بلوک `env` یا `apiKeyHelper` دارد، روی همه حساب‌های گروه اعمال می‌شود، از جمله حساب‌های کلید API. تنظیمات ارائه‌دهنده را بیرون از تنظیمات مشترک نگه دارید، یا حساب کلید API را در گروهی بدون `settings` بگذارید.

## مطالب مرتبط

- [مدیریت حساب‌ها](/features/accounts)
- [اجرای Claude Code](/features/launch)
- [مرجع فرمان‌ها](/docs/commands)
