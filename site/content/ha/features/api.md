---
title: Asusan API key
description: ccam add --api yana ƙara asusun API key don API na Anthropic ko gateway; ana karanta key a lokacin farawa daga variable ko umarni, ba a taɓa ajiye shi.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: 771d749f1f90
---

Asusun API key yana fara Claude Code da API key maimakon shigar biyan kuɗin Claude. Key ɗin zai iya kasancewa na API na Anthropic ko na kowane endpoint da ke magana da tsarin Messages iri ɗaya, kamar gateway na LiteLLM ko OpenRouter. ccam ba ya ajiye key ɗin. Yana karanta shi a duk lokacin da asusun ya fara, daga environment variable ko daga umarni kamar password manager, kuma yana miƙa shi ga wannan process na Claude Code kawai.

## Amfani

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

Asusun API key yana shiga rukunoni kamar kowane asusu (ta tsohuwa rukunin `default`, wanda ke raba saituna, `CLAUDE.md`, extensions da sabobin MCP), kuma yana aiki da gajerun umarni, `ccam folder` da `ccam ext`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--api` | Ƙara asusun API key maimakon shiga | a kashe |
| `--key-from env:NAME` | Karanta key ɗin daga environment variable `NAME` | |
| `--key-from command:COMMAND` | Gudanar da `COMMAND` (da `sh -c`, ko a Windows kamar yadda aka rubuta shi a cikin `cmd /d /s /c "…"`, ana riƙe alamomin zance) kuma a yi amfani da layin farko da ya buga | |
| `--base-url <url>` | Endpoint ɗin (`ANTHROPIC_BASE_URL`) | API na Anthropic |
| `--auth api-key\|bearer` | Aika key ɗin a matsayin `ANTHROPIC_API_KEY` (`x-api-key`) ko a matsayin `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, wanda gateways suka fi so) | `api-key` |
| `--model <model>` | Model ɗin (`ANTHROPIC_MODEL`) | Tsohon model na Claude Code |

Ana ajiye saitunan a `config.toml`:

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

## Yadda yake aiki

A kowane ƙaddamar da asusun, bayan dubawa na yau da kullum kuma kafin Claude Code ya fara:

1. ccam yana karanta key ɗin: `env:` daga environment, `command:` ta gudanar da umarnin tare da iyakar lokaci ta daƙiƙa 30. Umarni zai iya neman kalmar sirri a terminal; output na kuskurensa ma yana zuwa terminal. Key mara komai yana tsayar da ƙaddamarwar da saƙon da ke faɗin inda ake sa ran key ɗin.
2. Yana cire kowane variable da aka gada wanda ke zaɓar yadda Claude Code ke shiga ko wane mai bayarwa yake amfani da shi:
   - key da token: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - endpoint: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - maɓallan zaɓin mai bayarwa: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   Key da aka saita a shell ɗinka don wani abu dabam ba zai iya kaiwa wannan asusun ba, kuma shigar biyan kuɗi ba za ta iya samun fifiko a kan key ɗinsa ba.
3. Yana saita `ANTHROPIC_BASE_URL` da `ANTHROPIC_MODEL` idan an ba su, da key ɗin a matsayin `ANTHROPIC_API_KEY` ko `ANTHROPIC_AUTH_TOKEN`.

Asusan da ke da shigar biyan kuɗi suna riƙe environment ɗin da suka gada, kamar da. A zaman mu'amala Claude Code zai iya tambaya sau ɗaya ko zai yi amfani da key ɗin da ya samu a environment.

## Tsaro da yanayi na musamman

- Ba a taɓa rubuta key ɗin a `config.toml`, cache ko log ba, kuma `ccam --dry-run gw` ba ya karanta shi: dry run ba ya gudanar da tushen `command:`. Maimakon haka dry run yana nuna endpoint ɗin, model ɗin da inda key ɗin yake fitowa:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Variable da aka gada wanda ƙaddamarwar ke cirewa yana bayyana a matsayin `env -u <name>` a farko. Da tushen `env:` wanda variable ɗinsa babu komai a ciki, dry run yana gazawa kamar yadda ƙaddamarwar za ta gaza.
- `ccam add --api` yana duba tushen `env:` nan take kuma yana gargaɗi idan variable ɗin babu komai a ciki; ba a gudanar da tushen `command:` sai lokacin ƙaddamarwa.
- `ccam login` da `ccam logout` suna ƙin asusun API key, domin ba shi da shiga. `ccam doctor` yana nuna endpoint ɗinsa da tushen key ɗinsa, kuma yana gargaɗi idan variable na `env:` babu komai a ciki.
- `ccam usage` yana nuna `no limits (API)` gare shi. `ccam next --auto`, `--fallback` da `ccam pick` suna zaɓa a cikin asusan biyan kuɗi kawai, domin API key ba shi da iyakar awa 5 ko ta mako da za a kwatanta.
- `settings.json` da aka raba mai sashen `env` ko `apiKeyHelper` yana aiki ga kowane asusu a rukunin, har da asusan API key. Ka cire saitunan mai bayarwa daga saitunan da aka raba, ko ka saka asusun API key a rukunin da ba ya raba `settings`.

## Masu alaƙa

- [Sarrafa asusai](/features/accounts)
- [Fara Claude Code](/features/launch)
- [Bayanin umarnai](/docs/commands)
