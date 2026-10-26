---
title: API key கணக்குகள்
description: ccam add --api, Anthropic API அல்லது gateway-க்கு API key கணக்கைச் சேர்க்கும்; key தொடக்கத்தில் variable அல்லது கட்டளையிலிருந்து படிக்கப்படும், சேமிப்பதில்லை.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: 771d749f1f90
---

ஒரு API key கணக்கு, Claude subscription உள்நுழைவுக்குப் பதிலாக API key மூலம் Claude Code-ஐத் தொடங்குகிறது. அந்த key Anthropic-இன் API-க்கானதாகவோ, அதே Messages வடிவத்தைப் பேசும் எந்த endpoint-க்கானதாகவோ இருக்கலாம், உதாரணமாக LiteLLM அல்லது OpenRouter gateway. ccam key-ஐச் சேமிப்பதில்லை. கணக்கு தொடங்கும் ஒவ்வொரு முறையும், ஒரு environment variable-இலிருந்தோ password manager போன்ற கட்டளையிலிருந்தோ அதைப் படித்து, அந்த Claude Code process-க்கு மட்டுமே கொடுக்கிறது.

## பயன்பாடு

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

API key கணக்கு மற்ற எந்தக் கணக்கையும் போலவே குழுக்களில் சேர்கிறது (இயல்பாக `default` குழு, அது அமைப்புகள், `CLAUDE.md`, extensions, MCP சர்வர்கள் ஆகியவற்றைப் பகிர்கிறது), குறுக்குவழிகள், `ccam folder`, `ccam ext` ஆகியவற்றுடன் வேலை செய்கிறது.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--api` | உள்நுழைவதற்குப் பதிலாக API key கணக்கைச் சேர்க்கும் | off |
| `--key-from env:NAME` | environment variable `NAME`-இலிருந்து key-ஐப் படிக்கும் | |
| `--key-from command:COMMAND` | `COMMAND`-ஐ இயக்கி (`sh -c` மூலம், அல்லது Windows-இல் `cmd /d /s /c "…"`-இல் தட்டச்சு செய்தது போலவே, மேற்கோள் குறிகள் அப்படியே) அது அச்சிடும் முதல் வரியைப் பயன்படுத்தும் | |
| `--base-url <url>` | endpoint (`ANTHROPIC_BASE_URL`) | Anthropic-இன் API |
| `--auth api-key\|bearer` | key-ஐ `ANTHROPIC_API_KEY` (`x-api-key`) ஆக அல்லது `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, gateways பொதுவாக இதைக் கேட்கும்) ஆக அனுப்பும் | `api-key` |
| `--model <model>` | model (`ANTHROPIC_MODEL`) | Claude Code-இன் இயல்புநிலை |

அமைப்புகள் `config.toml`-இல் சேமிக்கப்படுகின்றன:

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

## இது எப்படி வேலை செய்கிறது

கணக்கின் ஒவ்வொரு தொடக்கத்திலும், வழக்கமான சரிபார்ப்புகளுக்குப் பிறகு, Claude Code தொடங்கும் முன்:

1. ccam key-ஐப் படிக்கிறது: `env:` environment-இலிருந்து, `command:` கட்டளையை 30 வினாடி timeout உடன் இயக்கி. கட்டளை terminal-இல் password கேட்கலாம்; அதன் error output-உம் terminal-க்குச் செல்லும். காலியான key தொடக்கத்தை நிறுத்தும், key எங்கே எதிர்பார்க்கப்பட்டது என்று செய்தி சொல்லும்.
2. Claude Code எப்படி உள்நுழைகிறது அல்லது எந்த provider-ஐப் பயன்படுத்துகிறது என்பதைத் தீர்மானிக்கும், மரபுரிமையாகப் பெற்ற ஒவ்வொரு variable-ஐயும் அது நீக்குகிறது:
   - keys மற்றும் tokens: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - endpoint: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - provider switches: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   உங்கள் shell-இல் வேறொன்றுக்காக அமைத்த key இந்தக் கணக்கை அடைய முடியாது, subscription உள்நுழைவு அதன் key-ஐ விட முன்னுரிமை பெற முடியாது.
3. கொடுக்கப்பட்டால் `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL` ஆகியவற்றை அமைக்கிறது, key-ஐ `ANTHROPIC_API_KEY` அல்லது `ANTHROPIC_AUTH_TOKEN` ஆக அமைக்கிறது.

subscription உள்நுழைவு உள்ள கணக்குகள் முன்பு போலவே மரபுரிமையாகப் பெற்ற environment-ஐ வைத்திருக்கின்றன. interactive அமர்வில், environment-இல் கண்டுபிடிக்கும் key-ஐப் பயன்படுத்தலாமா என்று Claude Code ஒருமுறை கேட்கலாம்.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- key ஒருபோதும் `config.toml`, cache அல்லது log-இல் எழுதப்படுவதில்லை, `ccam --dry-run gw` அதைப் படிப்பதில்லை: dry run ஒரு `command:` மூலத்தை இயக்காது. அதற்குப் பதிலாக dry run, endpoint-ஐயும் model-ஐயும் key எங்கிருந்து வருகிறது என்பதையும் காட்டுகிறது:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  தொடக்கம் நீக்கும் மரபுரிமையாகப் பெற்ற variable, வரியின் தொடக்கத்தில் `env -u <name>` ஆகக் காட்டப்படும். variable காலியாக உள்ள `env:` மூலத்துடன், தொடக்கம் தோல்வியடைவது போலவே dry run-உம் தோல்வியடையும்.
- `ccam add --api` ஒரு `env:` மூலத்தை உடனே சரிபார்த்து, variable காலியாக இருந்தால் எச்சரிக்கிறது; ஒரு `command:` மூலம் தொடக்கம் வரை இயங்குவதில்லை.
- API key கணக்குக்கு உள்நுழைவு இல்லாததால், `ccam login`, `ccam logout` அதை மறுக்கின்றன. `ccam doctor` அதன் endpoint-ஐயும் key மூலத்தையும் காட்டி, ஒரு `env:` variable காலியாக இருக்கும்போது எச்சரிக்கிறது.
- `ccam usage` அதற்கு `no limits (API)` எனக் காட்டுகிறது. API key-க்கு ஒப்பிட 5 மணி நேர அல்லது வாராந்திர வரம்பு இல்லாததால், `ccam next --auto`, `--fallback`, `ccam pick` ஆகியவை subscription கணக்குகளிலிருந்து மட்டுமே தேர்ந்தெடுக்கின்றன.
- `env` block அல்லது `apiKeyHelper` உள்ள பகிரப்பட்ட `settings.json`, API key கணக்குகள் உட்படக் குழுவின் ஒவ்வொரு கணக்குக்கும் பொருந்தும். provider அமைப்புகளைப் பகிரப்பட்ட அமைப்புகளுக்கு வெளியே வையுங்கள், அல்லது API key கணக்கை `settings` இல்லாத ஒரு குழுவில் வையுங்கள்.

## தொடர்புடையவை

- [கணக்குகளை நிர்வகித்தல்](/features/accounts)
- [Claude Code-ஐத் தொடங்குதல்](/features/launch)
- [கட்டளைக் குறிப்பு](/docs/commands)
