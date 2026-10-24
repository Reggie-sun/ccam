---
title: API కీ ఖాతాలు
description: ccam add --api, Anthropic API లేదా gateway కోసం API కీ ఖాతాను జోడిస్తుంది; కీ ప్రారంభంలో variable లేదా కమాండ్ నుంచి చదవబడుతుంది, ఎప్పుడూ సేవ్ కాదు.
updated: 2026-10-23
order: 13
translation: machine
sourceHash: "232a71525e42"
---

API కీ ఖాతా Claude సబ్‌స్క్రిప్షన్ లాగిన్ బదులు API కీతో Claude Code ను ప్రారంభిస్తుంది: Anthropic యొక్క API, లేదా అదే Messages ఫార్మాట్ మాట్లాడే ఎండ్‌పాయింట్, ఉదాహరణకు LiteLLM లేదా OpenRouter gateway. ccam కీని సేవ్ చేయదు. ఖాతా ప్రారంభమైన ప్రతిసారీ, అది ఒక environment variable నుంచి లేదా password manager వంటి కమాండ్ నుంచి కీని చదివి, ఆ Claude Code ప్రాసెస్‌కు మాత్రమే ఇస్తుంది.

## వినియోగం

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

API కీ ఖాతా ఇతర ఖాతాల లాగే గ్రూప్‌లలో చేరుతుంది (డిఫాల్ట్‌గా `default`: సెట్టింగ్‌లు, `CLAUDE.md`, extensions, MCP సర్వర్లు), షార్ట్‌కట్‌లు, `ccam folder`, `ccam ext` తో పనిచేస్తుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--api` | లాగిన్ చేయడం బదులు API కీ ఖాతాను జోడిస్తుంది | ఆఫ్ |
| `--key-from env:NAME` | environment variable `NAME` నుంచి కీని చదువుతుంది | |
| `--key-from command:COMMAND` | `COMMAND` ను రన్ చేసి (`sh -c` తో, లేదా Windows లో టైప్ చేసినట్లే `cmd /d /s /c "…"` లో, కోట్‌లను అలాగే ఉంచి) అది ప్రింట్ చేసే మొదటి లైన్‌ను ఉపయోగిస్తుంది | |
| `--base-url <url>` | ఎండ్‌పాయింట్ (`ANTHROPIC_BASE_URL`) | Anthropic యొక్క API |
| `--auth api-key\|bearer` | కీని `ANTHROPIC_API_KEY` (`x-api-key`) గా లేదా `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, gateways సాధారణంగా కోరేది) గా పంపుతుంది | `api-key` |
| `--model <model>` | మోడల్ (`ANTHROPIC_MODEL`) | Claude Code డిఫాల్ట్ |

సెట్టింగ్‌లు `config.toml` లో సేవ్ అవుతాయి:

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

## ఇది ఎలా పనిచేస్తుంది

ఖాతా ప్రతి ప్రారంభంలో, సాధారణ తనిఖీల తర్వాత, Claude Code ప్రారంభమయ్యే ముందు:

1. ccam కీని చదువుతుంది: `env:` environment నుంచి, `command:` కమాండ్‌ను 30 సెకన్ల పరిమితితో రన్ చేసి. కమాండ్ టెర్మినల్‌లో password అడగవచ్చు; దాని ఎర్రర్ అవుట్‌పుట్ కూడా టెర్మినల్‌కే వెళ్తుంది. ఖాళీ కీ ప్రారంభాన్ని ఆపుతుంది, కీ ఎక్కడ ఆశించబడిందో సందేశం చెబుతుంది.
2. Claude Code ఎలా సైన్ ఇన్ అవుతుందో లేదా ఏ provider ను ఉపయోగిస్తుందో నిర్ణయించే, వారసత్వంగా వచ్చిన ప్రతి variable ను అది తీసివేస్తుంది: కీలు, టోకెన్‌లు (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), ఎండ్‌పాయింట్ (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`), provider స్విచ్‌లు (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). మీ shell లో వేరే దాని కోసం సెట్ చేసిన కీ ఈ ఖాతాను చేరలేదు, సబ్‌స్క్రిప్షన్ లాగిన్ దాని కీని అధిగమించలేదు.
3. ఇచ్చినప్పుడు `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL` ను సెట్ చేస్తుంది, కీని `ANTHROPIC_API_KEY` లేదా `ANTHROPIC_AUTH_TOKEN` గా సెట్ చేస్తుంది.

సబ్‌స్క్రిప్షన్ లాగిన్ ఉన్న ఖాతాలు మునుపటిలాగే వారసత్వంగా వచ్చిన environment ను ఉంచుకుంటాయి. ఇంటరాక్టివ్ సెషన్‌లో, environment లో కనిపించిన కీని ఉపయోగించాలా అని Claude Code ఒక్కసారి అడగవచ్చు.

## భద్రత, ప్రత్యేక సందర్భాలు

- కీ ఎప్పుడూ `config.toml`, క్యాష్ లేదా లాగ్‌లో రాయబడదు, `ccam --dry-run gw` దాన్ని చదవదు: డ్రై రన్ `command:` సోర్స్‌ను రన్ చేయదు. బదులుగా డ్రై రన్ ఎండ్‌పాయింట్‌ను, మోడల్‌ను, కీ ఎక్కడి నుంచి వస్తుందో చూపిస్తుంది:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  ప్రారంభం తీసివేసే వారసత్వ variable మొదట్లో `env -u <name>` గా కనిపిస్తుంది. variable ఖాళీగా ఉన్న `env:` సోర్స్‌తో, ప్రారంభం విఫలమయ్యే విధంగానే డ్రై రన్ కూడా విఫలమవుతుంది.
- `ccam add --api` `env:` సోర్స్‌ను వెంటనే తనిఖీ చేసి, variable ఖాళీగా ఉంటే హెచ్చరిస్తుంది; `command:` సోర్స్ మొదటిసారి ప్రారంభంలో రన్ అవుతుంది.
- API కీ ఖాతాకు లాగిన్ లేనందున `ccam login`, `ccam logout` దాన్ని తిరస్కరిస్తాయి. `ccam doctor` దాని ఎండ్‌పాయింట్‌ను, కీ సోర్స్‌ను చూపిస్తుంది, `env:` variable ఖాళీగా ఉన్నప్పుడు హెచ్చరిస్తుంది.
- `ccam usage` దానికి `no limits (API)` అని చూపిస్తుంది. API కీకి పోల్చడానికి 5 గంటల లేదా వారపు పరిమితి లేనందున, `ccam next --auto`, `--fallback`, `ccam pick` సబ్‌స్క్రిప్షన్ ఖాతాల నుంచి మాత్రమే ఎంచుకుంటాయి.
- `env` బ్లాక్ లేదా `apiKeyHelper` ఉన్న షేర్ చేసిన `settings.json`, API కీ ఖాతాలతో సహా గ్రూప్‌లోని ప్రతి ఖాతాకు వర్తిస్తుంది. provider సెట్టింగ్‌లను షేర్ చేసిన సెట్టింగ్‌లకు దూరంగా ఉంచండి, లేదా API కీ ఖాతాను `settings` లేని గ్రూప్‌లో ఉంచండి.

## సంబంధిత పేజీలు

- [ఖాతాలను నిర్వహించడం](/features/accounts)
- [Claude Code ను ప్రారంభించడం](/features/launch)
- [కమాండ్ రిఫరెన్స్](/docs/commands)
