---
title: API key अकाउंट
description: ccam add --api, Anthropic के API या किसी gateway के लिए API key अकाउंट जोड़ता है; key launch पर किसी variable या कमांड से पढ़ी जाती है, कभी सेव नहीं होती।
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API key अकाउंट Claude subscription लॉगिन के बजाय API key से Claude Code शुरू करता है: Anthropic का API, या कोई ऐसा endpoint जो वही Messages format समझता है, जैसे LiteLLM या OpenRouter gateway। ccam key सेव नहीं करता। अकाउंट जब भी शुरू होता है, यह key किसी environment variable से या password manager जैसी किसी कमांड से पढ़ता है, और उसे केवल उसी Claude Code process को देता है।

## उपयोग

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

API key अकाउंट किसी भी दूसरे अकाउंट की तरह सेटिंग्स, `CLAUDE.md`, extensions और MCP सर्वर साझा करता है, और शॉर्टकट, `ccam folder` और `ccam ext` के साथ काम करता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--api` | लॉगिन करने के बजाय API key अकाउंट जोड़ें | बंद |
| `--key-from env:NAME` | key को environment variable `NAME` से पढ़ें | |
| `--key-from command:COMMAND` | `COMMAND` चलाएँ (`sh -c` से, या Windows पर `cmd /C` से) और उसकी प्रिंट की गई पहली पंक्ति इस्तेमाल करें | |
| `--base-url <url>` | endpoint (`ANTHROPIC_BASE_URL`) | Anthropic का API |
| `--auth api-key\|bearer` | key को `ANTHROPIC_API_KEY` (`x-api-key`) के रूप में भेजें या `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, जो gateways आमतौर पर चाहते हैं) के रूप में | `api-key` |
| `--model <model>` | मॉडल (`ANTHROPIC_MODEL`) | Claude Code का डिफ़ॉल्ट |

ये सेटिंग्स `config.toml` में सेव होती हैं:

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

## यह कैसे काम करता है

अकाउंट के हर launch पर, सामान्य जाँचों के बाद और Claude Code शुरू होने से पहले:

1. ccam key पढ़ता है: `env:` को environment से, `command:` को 30 सेकंड की सीमा के साथ कमांड चलाकर। कोई कमांड terminal पर password माँग सकती है; उसका error output भी terminal पर जाता है। ख़ाली key launch को रोक देती है, एक संदेश के साथ जो बताता है कि key कहाँ अपेक्षित थी।
2. यह हर विरासत में मिला variable हटाता है जो तय करता है कि Claude Code कैसे sign in करता है या कौन-सा provider इस्तेमाल करता है: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` और `CLAUDE_CODE_USE_FOUNDRY`। आपके shell में किसी और काम के लिए सेट की गई key इस अकाउंट तक नहीं पहुँच सकती, और कोई subscription लॉगिन इसकी key पर हावी नहीं हो सकता।
3. दिए जाने पर यह `ANTHROPIC_BASE_URL` और `ANTHROPIC_MODEL` सेट करता है, और key को `ANTHROPIC_API_KEY` या `ANTHROPIC_AUTH_TOKEN` के रूप में।

subscription लॉगिन वाले अकाउंट पहले की तरह विरासत में मिला environment ही रखते हैं। interactive session में Claude Code एक बार पूछ सकता है कि environment में मिली key इस्तेमाल करनी है या नहीं।

## सुरक्षा और विशेष स्थितियाँ

- key कभी `config.toml`, किसी cache या log में नहीं लिखी जाती, और `ccam --dry-run gw` उसे नहीं पढ़ता: dry run में `command:` source नहीं चलाया जाता।
- `ccam add --api` किसी `env:` source को तुरंत जाँचता है और variable ख़ाली होने पर चेतावनी देता है; `command:` source पहली बार launch पर चलता है।
- `ccam login` और `ccam logout` API key अकाउंट को मना कर देते हैं, क्योंकि उसका कोई लॉगिन नहीं होता। `ccam doctor` उसका endpoint और key source दिखाता है, और `env:` variable ख़ाली होने पर चेतावनी देता है।
- `ccam usage` उसके लिए `API` दिखाता है। `ccam next --auto`, `--fallback` और `ccam pick` केवल subscription अकाउंट्स में से चुनते हैं, क्योंकि API key की तुलना करने लायक कोई 5 घंटे या साप्ताहिक सीमा नहीं होती।
- `env` block या `apiKeyHelper` वाला साझा `settings.json` ग्रुप के हर अकाउंट पर लागू होता है, API key अकाउंट्स पर भी। provider की सेटिंग्स को साझा सेटिंग्स से बाहर रखें।

## संबंधित

- [अकाउंट प्रबंधित करना](/features/accounts)
- [Claude Code शुरू करना](/features/launch)
- [कमांड संदर्भ](/docs/commands)
