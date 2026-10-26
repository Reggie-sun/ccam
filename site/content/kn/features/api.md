---
title: API ಕೀ ಖಾತೆಗಳು
description: ccam add --api, Anthropic API ಅಥವಾ ಗೇಟ್‌ವೇಗಾಗಿ API ಕೀ ಖಾತೆ ಸೇರಿಸುತ್ತದೆ; ಕೀಯನ್ನು ಆರಂಭದಲ್ಲಿ ವೇರಿಯಬಲ್ ಅಥವಾ ಕಮಾಂಡ್‌ನಿಂದ ಓದುತ್ತದೆ, ಸಂಗ್ರಹಿಸುವುದಿಲ್ಲ.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: "771d749f1f90"
---

API ಕೀ ಖಾತೆ Claude ಚಂದಾದಾರಿಕೆ ಲಾಗಿನ್ ಬದಲು API ಕೀ ಮೂಲಕ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಕೀ Anthropic ನ API ಗಾಗಿ, ಅಥವಾ LiteLLM ಅಥವಾ OpenRouter ಗೇಟ್‌ವೇಯಂತಹ ಅದೇ Messages ಫಾರ್ಮ್ಯಾಟ್ ಬಳಸುವ ಯಾವುದೇ ಎಂಡ್‌ಪಾಯಿಂಟ್‌ಗಾಗಿ ಇರಬಹುದು. ccam ಕೀಯನ್ನು ಸಂಗ್ರಹಿಸುವುದಿಲ್ಲ. ಖಾತೆ ಆರಂಭವಾದ ಪ್ರತಿ ಬಾರಿಯೂ ಅದನ್ನು ಎನ್ವಿರಾನ್‌ಮೆಂಟ್ ವೇರಿಯಬಲ್‌ನಿಂದ ಅಥವಾ ಪಾಸ್‌ವರ್ಡ್ ಮ್ಯಾನೇಜರ್‌ನಂತಹ ಕಮಾಂಡ್‌ನಿಂದ ಓದುತ್ತದೆ, ಮತ್ತು ಆ Claude Code ಪ್ರೊಸೆಸ್‌ಗೆ ಮಾತ್ರ ನೀಡುತ್ತದೆ.

## ಬಳಕೆ

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

API ಕೀ ಖಾತೆ ಇತರ ಯಾವುದೇ ಖಾತೆಯಂತೆ ಗುಂಪುಗಳನ್ನು ಸೇರುತ್ತದೆ (ಡೀಫಾಲ್ಟ್ ಆಗಿ `default` ಗುಂಪು, ಅದು ಸೆಟ್ಟಿಂಗ್‌ಗಳು, `CLAUDE.md`, ಎಕ್ಸ್‌ಟೆನ್ಶನ್‌ಗಳು ಮತ್ತು MCP ಸರ್ವರ್‌ಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳುತ್ತದೆ), ಮತ್ತು ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು, `ccam folder` ಮತ್ತು `ccam ext` ಜೊತೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `--api` | ಲಾಗಿನ್ ಮಾಡುವ ಬದಲು API ಕೀ ಖಾತೆಯನ್ನು ಸೇರಿಸುತ್ತದೆ | ಆಫ್ |
| `--key-from env:NAME` | ಎನ್ವಿರಾನ್‌ಮೆಂಟ್ ವೇರಿಯಬಲ್ `NAME` ನಿಂದ ಕೀಯನ್ನು ಓದುತ್ತದೆ | |
| `--key-from command:COMMAND` | `COMMAND` ಅನ್ನು ರನ್ ಮಾಡುತ್ತದೆ (`sh -c` ಮೂಲಕ, ಅಥವಾ Windows ನಲ್ಲಿ ಟೈಪ್ ಮಾಡಿದಂತೆಯೇ `cmd /d /s /c "…"` ಒಳಗೆ, ಉದ್ಧರಣ ಚಿಹ್ನೆಗಳನ್ನು ಉಳಿಸಿಕೊಂಡು) ಮತ್ತು ಅದು ಮುದ್ರಿಸುವ ಮೊದಲ ಸಾಲನ್ನು ಬಳಸುತ್ತದೆ | |
| `--base-url <url>` | ಎಂಡ್‌ಪಾಯಿಂಟ್ (`ANTHROPIC_BASE_URL`) | Anthropic ನ API |
| `--auth api-key\|bearer` | ಕೀಯನ್ನು `ANTHROPIC_API_KEY` (`x-api-key`) ಆಗಿ ಅಥವಾ `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, ಗೇಟ್‌ವೇಗಳು ಸಾಮಾನ್ಯವಾಗಿ ಬಯಸುವುದು) ಆಗಿ ಕಳುಹಿಸುತ್ತದೆ | `api-key` |
| `--model <model>` | ಮಾಡೆಲ್ (`ANTHROPIC_MODEL`) | Claude Code ನ ಡೀಫಾಲ್ಟ್ |

ಸೆಟ್ಟಿಂಗ್‌ಗಳನ್ನು `config.toml` ನಲ್ಲಿ ಸಂಗ್ರಹಿಸಲಾಗುತ್ತದೆ:

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

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

ಖಾತೆಯ ಪ್ರತಿ ಆರಂಭದಲ್ಲೂ, ಸಾಮಾನ್ಯ ಪರಿಶೀಲನೆಗಳ ನಂತರ ಮತ್ತು Claude Code ಆರಂಭವಾಗುವ ಮೊದಲು:

1. ccam ಕೀಯನ್ನು ಓದುತ್ತದೆ: `env:` ಎನ್ವಿರಾನ್‌ಮೆಂಟ್‌ನಿಂದ, `command:` ಕಮಾಂಡ್ ಅನ್ನು 30 ಸೆಕೆಂಡ್ ಟೈಮ್‌ಔಟ್‌ನೊಂದಿಗೆ ರನ್ ಮಾಡಿ. ಕಮಾಂಡ್ ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಪಾಸ್‌ವರ್ಡ್ ಕೇಳಬಹುದು; ಅದರ ಎರರ್ ಔಟ್‌ಪುಟ್ ಕೂಡ ಟರ್ಮಿನಲ್‌ಗೆ ಹೋಗುತ್ತದೆ. ಖಾಲಿ ಕೀ, ಕೀ ಎಲ್ಲಿ ನಿರೀಕ್ಷಿಸಲಾಗಿತ್ತು ಎಂದು ಹೇಳುವ ಸಂದೇಶದೊಂದಿಗೆ ಆರಂಭವನ್ನು ನಿಲ್ಲಿಸುತ್ತದೆ.
2. Claude Code ಹೇಗೆ ಸೈನ್ ಇನ್ ಆಗುತ್ತದೆ ಅಥವಾ ಯಾವ ಪ್ರೊವೈಡರ್ ಬಳಸುತ್ತದೆ ಎಂದು ನಿರ್ಧರಿಸುವ ಪ್ರತಿ ಆನುವಂಶಿಕ ವೇರಿಯಬಲ್ ಅನ್ನು ಅದು ತೆಗೆದುಹಾಕುತ್ತದೆ:
   - ಕೀಗಳು ಮತ್ತು ಟೋಕನ್‌ಗಳು: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - ಎಂಡ್‌ಪಾಯಿಂಟ್: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - ಪ್ರೊವೈಡರ್ ಸ್ವಿಚ್‌ಗಳು: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   ಬೇರೆ ಉದ್ದೇಶಕ್ಕಾಗಿ ನಿಮ್ಮ ಶೆಲ್‌ನಲ್ಲಿ ಹೊಂದಿಸಿದ ಕೀ ಈ ಖಾತೆಯನ್ನು ತಲುಪಲಾರದು, ಮತ್ತು ಚಂದಾದಾರಿಕೆ ಲಾಗಿನ್ ಅದರ ಕೀಗಿಂತ ಆದ್ಯತೆ ಪಡೆಯಲಾರದು.
3. ನೀಡಿದ್ದರೆ `ANTHROPIC_BASE_URL` ಮತ್ತು `ANTHROPIC_MODEL` ಅನ್ನು, ಮತ್ತು ಕೀಯನ್ನು `ANTHROPIC_API_KEY` ಅಥವಾ `ANTHROPIC_AUTH_TOKEN` ಆಗಿ ಅದು ಹೊಂದಿಸುತ್ತದೆ.

ಚಂದಾದಾರಿಕೆ ಲಾಗಿನ್ ಇರುವ ಖಾತೆಗಳು ಮೊದಲಿನಂತೆಯೇ ತಾವು ಪಡೆದ ಎನ್ವಿರಾನ್‌ಮೆಂಟ್ ಅನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತವೆ. ಇಂಟರಾಕ್ಟಿವ್ ಸೆಷನ್‌ನಲ್ಲಿ, ಎನ್ವಿರಾನ್‌ಮೆಂಟ್‌ನಲ್ಲಿ ಕಂಡ ಕೀಯನ್ನು ಬಳಸಬೇಕೇ ಎಂದು Claude Code ಒಮ್ಮೆ ಕೇಳಬಹುದು.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- ಕೀಯನ್ನು ಎಂದಿಗೂ `config.toml`, ಕ್ಯಾಶ್ ಅಥವಾ ಲಾಗ್‌ಗೆ ಬರೆಯಲಾಗುವುದಿಲ್ಲ, ಮತ್ತು `ccam --dry-run gw` ಅದನ್ನು ಓದುವುದಿಲ್ಲ: ಡ್ರೈ ರನ್ `command:` ಮೂಲವನ್ನು ರನ್ ಮಾಡುವುದಿಲ್ಲ. ಬದಲಿಗೆ ಡ್ರೈ ರನ್ ಎಂಡ್‌ಪಾಯಿಂಟ್, ಮಾಡೆಲ್ ಮತ್ತು ಕೀ ಎಲ್ಲಿಂದ ಬರುತ್ತದೆ ಎಂಬುದನ್ನು ತೋರಿಸುತ್ತದೆ:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  ಆರಂಭವು ತೆಗೆದುಹಾಕುವ ಆನುವಂಶಿಕ ವೇರಿಯಬಲ್ ಕಮಾಂಡ್‌ನ ಆರಂಭದಲ್ಲಿ `env -u <name>` ಎಂದು ಕಾಣಿಸುತ್ತದೆ. ವೇರಿಯಬಲ್ ಖಾಲಿಯಾಗಿರುವ `env:` ಮೂಲದೊಂದಿಗೆ, ಆರಂಭ ವಿಫಲವಾಗುವಂತೆಯೇ ಡ್ರೈ ರನ್ ಕೂಡ ವಿಫಲವಾಗುತ್ತದೆ.
- `ccam add --api` `env:` ಮೂಲವನ್ನು ತಕ್ಷಣ ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು ವೇರಿಯಬಲ್ ಖಾಲಿಯಾಗಿದ್ದರೆ ಎಚ್ಚರಿಸುತ್ತದೆ; `command:` ಮೂಲವನ್ನು ಆರಂಭದವರೆಗೆ ರನ್ ಮಾಡಲಾಗುವುದಿಲ್ಲ.
- `ccam login` ಮತ್ತು `ccam logout` API ಕೀ ಖಾತೆಯನ್ನು ನಿರಾಕರಿಸುತ್ತವೆ, ಏಕೆಂದರೆ ಅದಕ್ಕೆ ಲಾಗಿನ್ ಇಲ್ಲ. `ccam doctor` ಅದರ ಎಂಡ್‌ಪಾಯಿಂಟ್ ಮತ್ತು ಕೀ ಮೂಲವನ್ನು ತೋರಿಸುತ್ತದೆ, ಮತ್ತು `env:` ವೇರಿಯಬಲ್ ಖಾಲಿಯಾಗಿದ್ದರೆ ಎಚ್ಚರಿಸುತ್ತದೆ.
- `ccam usage` ಅದಕ್ಕೆ `no limits (API)` ತೋರಿಸುತ್ತದೆ. API ಕೀಗೆ ಹೋಲಿಸಲು 5 ಗಂಟೆಯ ಅಥವಾ ವಾರದ ಮಿತಿ ಇಲ್ಲದಿರುವುದರಿಂದ, `ccam next --auto`, `--fallback` ಮತ್ತು `ccam pick` ಚಂದಾದಾರಿಕೆ ಖಾತೆಗಳಲ್ಲಿ ಮಾತ್ರ ಆಯ್ಕೆಮಾಡುತ್ತವೆ.
- `env` ಬ್ಲಾಕ್ ಅಥವಾ `apiKeyHelper` ಇರುವ ಹಂಚಿಕೊಂಡ `settings.json` ಗುಂಪಿನ ಪ್ರತಿ ಖಾತೆಗೂ, API ಕೀ ಖಾತೆಗಳನ್ನೂ ಒಳಗೊಂಡಂತೆ, ಅನ್ವಯಿಸುತ್ತದೆ. ಪ್ರೊವೈಡರ್ ಸೆಟ್ಟಿಂಗ್‌ಗಳನ್ನು ಹಂಚಿಕೊಂಡ ಸೆಟ್ಟಿಂಗ್‌ಗಳಿಂದ ಹೊರಗಿಡಿ, ಅಥವಾ API ಕೀ ಖಾತೆಯನ್ನು `settings` ಇಲ್ಲದ ಗುಂಪಿಗೆ ಸೇರಿಸಿ.

## ಸಂಬಂಧಿತ

- [ಖಾತೆಗಳ ನಿರ್ವಹಣೆ](/features/accounts)
- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಕಮಾಂಡ್ ಉಲ್ಲೇಖ](/docs/commands)
