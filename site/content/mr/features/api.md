---
title: API key खाती
description: ccam add --api, Anthropic च्या API किंवा gateway साठी API key खाते जोडते; key launch वेळी variable किंवा कमांडमधून वाचली जाते, कधीही साठवली जात नाही.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: "771d749f1f90"
---

API key खाते Claude subscription लॉगिनऐवजी API key ने Claude Code सुरू करते. ही key Anthropic च्या API साठी असू शकते, किंवा तोच Messages format समजणाऱ्या कोणत्याही endpoint साठी, जसे LiteLLM किंवा OpenRouter gateway. ccam key साठवत नाही. खाते प्रत्येक वेळी सुरू होताना ते key एखाद्या environment variable मधून किंवा password manager सारख्या कमांडमधून वाचते, आणि फक्त त्या Claude Code process ला देते.

## वापर

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

API key खाते इतर कोणत्याही खात्याप्रमाणे गटांमध्ये सामील होते (डीफॉल्टनुसार `default` गट, जो सेटिंग्ज, `CLAUDE.md`, extensions आणि MCP सर्व्हर शेअर करतो), आणि शॉर्टकट, `ccam folder` आणि `ccam ext` सह काम करते.

## पर्याय

| पर्याय | काय करते | डीफॉल्ट |
| --- | --- | --- |
| `--api` | लॉगिन करण्याऐवजी API key खाते जोडा | बंद |
| `--key-from env:NAME` | key environment variable `NAME` मधून वाचा | |
| `--key-from command:COMMAND` | `COMMAND` चालवा (`sh -c` ने, किंवा Windows वर जशी टाइप केली तशी `cmd /d /s /c "…"` मध्ये, quotes तसेच ठेवून) आणि ती छापत असलेली पहिली ओळ वापरा | |
| `--base-url <url>` | endpoint (`ANTHROPIC_BASE_URL`) | Anthropic चा API |
| `--auth api-key\|bearer` | key `ANTHROPIC_API_KEY` (`x-api-key`) म्हणून किंवा `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, जे gateways सहसा मागतात) म्हणून पाठवा | `api-key` |
| `--model <model>` | model (`ANTHROPIC_MODEL`) | Claude Code चा डीफॉल्ट |

या सेटिंग्ज `config.toml` मध्ये साठवल्या जातात:

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

## हे कसे काम करते

खात्याच्या प्रत्येक launch वेळी, नेहमीच्या तपासण्यांनंतर आणि Claude Code सुरू होण्यापूर्वी:

1. ccam key वाचते: `env:` environment मधून, `command:` 30 सेकंदांच्या timeout सह कमांड चालवून. एखादी कमांड terminal वर password मागू शकते; तिचा error output सुद्धा terminal वर जातो. रिकामी key launch थांबवते, आणि key कुठे अपेक्षित होती ते सांगणारा संदेश दाखवते.
2. Claude Code कसे sign in करते किंवा कोणता provider वापरते हे ठरवणारा प्रत्येक वारशाने मिळालेला variable ते काढून टाकते:
   - keys आणि tokens: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - endpoint: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - provider switches: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   तुमच्या shell मध्ये इतर कामासाठी सेट केलेली key या खात्यापर्यंत पोहोचू शकत नाही, आणि subscription लॉगिनला या खात्याच्या key पेक्षा प्राधान्य मिळू शकत नाही.
3. दिले असल्यास ते `ANTHROPIC_BASE_URL` आणि `ANTHROPIC_MODEL` सेट करते, आणि key `ANTHROPIC_API_KEY` किंवा `ANTHROPIC_AUTH_TOKEN` म्हणून.

subscription लॉगिन असलेली खाती पूर्वीप्रमाणेच वारशाने मिळालेले environment ठेवतात. interactive session मध्ये Claude Code environment मध्ये सापडलेली key वापरायची का, हे एकदा विचारू शकते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- key कधीही `config.toml`, cache किंवा log मध्ये लिहिली जात नाही, आणि `ccam --dry-run gw` ती वाचत नाही: dry run मध्ये `command:` source चालवला जात नाही. त्याऐवजी dry run endpoint, model आणि key कुठून येते ते दाखवतो:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  launch काढून टाकणारा वारशाने मिळालेला variable सुरुवातीला `env -u <name>` म्हणून दिसतो. `env:` source चा variable रिकामा असेल, तर launch जसे अयशस्वी होईल तसाच dry run अयशस्वी होतो.
- `ccam add --api` `env:` source लगेच तपासते आणि variable रिकामा असल्यास इशारा देते; `command:` source launch होईपर्यंत चालवला जात नाही.
- `ccam login` आणि `ccam logout` API key खाते नाकारतात, कारण त्याला लॉगिन नसते. `ccam doctor` त्याचा endpoint आणि key source दाखवते, आणि `env:` variable रिकामा असल्यास इशारा देते.
- `ccam usage` त्यासाठी `no limits (API)` दाखवते. `ccam next --auto`, `--fallback` आणि `ccam pick` फक्त subscription खात्यांमधून निवडतात, कारण API key ला तुलना करण्यासारखी 5 तासांची किंवा साप्ताहिक मर्यादा नसते.
- `env` block किंवा `apiKeyHelper` असलेली शेअर केलेली `settings.json` गटातील प्रत्येक खात्याला लागू होते, API key खात्यांनाही. provider सेटिंग्ज शेअर केलेल्या सेटिंग्जबाहेर ठेवा, किंवा API key खाते `settings` शेअर न करणाऱ्या गटात ठेवा.

## संबंधित

- [खाती व्यवस्थापित करणे](/features/accounts)
- [Claude Code सुरू करणे](/features/launch)
- [कमांड संदर्भ](/docs/commands)
