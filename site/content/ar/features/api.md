---
title: حسابات مفاتيح API
description: يضيف ccam add --api حساب مفتاح API لواجهة Anthropic أو لبوابة؛ ويُقرأ المفتاح عند التشغيل من متغير أو أمر، ولا يُخزَّن أبدًا.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

يشغّل حساب مفتاح API برنامج Claude Code بمفتاح API بدلًا من تسجيل الدخول باشتراك Claude: واجهة Anthropic البرمجية، أو نقطة نهاية تتحدث صيغة Messages نفسها، مثل بوابة LiteLLM أو OpenRouter. لا يخزّن ccam المفتاح. بل يقرؤه في كل مرة يبدأ فيها الحساب، من متغير بيئة أو من أمر مثل مدير كلمات المرور، ولا يعطيه إلا لعملية Claude Code تلك.

## الاستخدام

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

يشارك حساب مفتاح API الإعدادات و`CLAUDE.md` والإضافات وخوادم MCP مثل أي حساب آخر، ويعمل مع الأوامر المختصرة و`ccam folder` و`ccam ext`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--api` | يضيف حساب مفتاح API بدلًا من تسجيل الدخول | معطّل |
| `--key-from env:NAME` | يقرأ المفتاح من متغير البيئة `NAME` | |
| `--key-from command:COMMAND` | يشغّل `COMMAND` (عبر `sh -c`، أو `cmd /C` على Windows) ويستخدم أول سطر يطبعه | |
| `--base-url <url>` | نقطة النهاية (`ANTHROPIC_BASE_URL`) | واجهة Anthropic البرمجية |
| `--auth api-key\|bearer` | يرسل المفتاح بوصفه `ANTHROPIC_API_KEY` (`x-api-key`) أو `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`، وهو ما تريده البوابات عادةً) | `api-key` |
| `--model <model>` | النموذج (`ANTHROPIC_MODEL`) | النموذج الافتراضي في Claude Code |

تُحفظ الإعدادات في `config.toml`:

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

## كيف يعمل

عند كل تشغيل للحساب، بعد الفحوص المعتادة وقبل أن يبدأ Claude Code:

1. يقرأ ccam المفتاح: `env:` من البيئة، و`command:` بتشغيل الأمر مع حد قدره 30 ثانية. ويمكن للأمر أن يطلب كلمة مرور في الطرفية؛ ويذهب مخرج أخطائه إلى الطرفية أيضًا. والمفتاح الفارغ يوقف التشغيل مع رسالة تذكر أين كان المفتاح متوقعًا.
2. يزيل كل متغير موروث يحدد طريقة تسجيل دخول Claude Code أو المزوّد الذي يستخدمه: `ANTHROPIC_API_KEY` و`ANTHROPIC_AUTH_TOKEN` و`CLAUDE_CODE_OAUTH_TOKEN` و`ANTHROPIC_BASE_URL` و`ANTHROPIC_MODEL` و`CLAUDE_CODE_USE_BEDROCK` و`CLAUDE_CODE_USE_VERTEX` و`CLAUDE_CODE_USE_FOUNDRY`. فلا يصل إلى هذا الحساب مفتاح مضبوط في صدفتك لغرض آخر، ولا يتغلب تسجيل دخول باشتراك على مفتاحه.
3. يضبط `ANTHROPIC_BASE_URL` و`ANTHROPIC_MODEL` إذا حُدّدا، والمفتاح بوصفه `ANTHROPIC_API_KEY` أو `ANTHROPIC_AUTH_TOKEN`.

تحتفظ الحسابات التي لديها تسجيل دخول باشتراك بالبيئة التي ترثها، كما من قبل. وفي الجلسة التفاعلية قد يسأل Claude Code مرة واحدة هل يستخدم المفتاح الذي يجده في البيئة.

## الأمان والحالات الخاصة

- لا يُكتب المفتاح أبدًا في `config.toml` ولا في ذاكرة مؤقتة ولا في سجل، ولا يقرؤه `ccam --dry-run gw`: فالتشغيل التجريبي لا يشغّل مصدر `command:`.
- يتحقق `ccam add --api` من مصدر `env:` فورًا ويحذّر إذا كان المتغير فارغًا؛ أما مصدر `command:` فيُشغَّل أول مرة عند التشغيل.
- يرفض `ccam login` و`ccam logout` حساب مفتاح API، لأنه لا تسجيل دخول له. ويعرض `ccam doctor` نقطة نهايته ومصدر مفتاحه، ويحذّر عندما يكون متغير `env:` فارغًا.
- يعرض `ccam usage` القيمة `API` له. ويختار `ccam next --auto` و`--fallback` و`ccam pick` من بين حسابات الاشتراك فقط، لأن مفتاح API ليس له حد للخمس ساعات أو للأسبوع يمكن المقارنة به.
- ملف `settings.json` المشترك الذي فيه كتلة `env` أو `apiKeyHelper` ينطبق على كل حساب في المجموعة، بما في ذلك حسابات مفاتيح API. أبقِ إعدادات المزوّد خارج الإعدادات المشتركة.

## صفحات ذات صلة

- [إدارة الحسابات](/features/accounts)
- [تشغيل Claude Code](/features/launch)
- [مرجع الأوامر](/docs/commands)
