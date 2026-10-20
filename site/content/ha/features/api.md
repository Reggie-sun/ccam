---
title: Asusan API key
description: ccam add --api yana ƙara asusun API key don API na Anthropic ko gateway; ana karanta key a lokacin farawa daga variable ko umarni, ba a taɓa ajiye shi.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Asusun API key yana fara Claude Code da API key maimakon shigar biyan kuɗin Claude: API na Anthropic, ko endpoint da ke magana da tsarin Messages iri ɗaya, kamar gateway na LiteLLM ko OpenRouter. ccam ba ya ajiye key ɗin. Yana karanta shi a duk lokacin da asusun ya fara, daga environment variable ko daga umarni kamar password manager, kuma yana miƙa shi ga wannan process na Claude Code kawai.

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
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

Asusun API key yana raba saituna, `CLAUDE.md`, extensions da sabobin MCP kamar kowane asusu, kuma yana aiki da gajerun umarni, `ccam folder` da `ccam ext`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--api` | Ƙara asusun API key maimakon shiga | a kashe |
| `--key-from env:NAME` | Karanta key ɗin daga environment variable `NAME` | |
| `--key-from command:COMMAND` | Gudanar da `COMMAND` (da `sh -c`, ko `cmd /C` a Windows) kuma a yi amfani da layin farko da ya buga | |
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

1. ccam yana karanta key ɗin: `env:` daga environment, `command:` ta gudanar da umarnin tare da iyakar daƙiƙa 30. Umarni zai iya neman kalmar sirri a terminal; output na kuskurensa ma yana zuwa terminal. Key mara komai yana tsayar da ƙaddamarwar da saƙon da ke faɗin inda ake sa ran key ɗin.
2. Yana cire kowane variable da aka gada wanda ke zaɓar yadda Claude Code ke shiga ko wane mai bayarwa yake amfani da shi: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` da `CLAUDE_CODE_USE_FOUNDRY`. Key da aka saita a shell ɗinka don wani abu dabam ba zai iya kaiwa wannan asusun ba, kuma shigar biyan kuɗi ba za ta iya rinjayar key ɗinsa ba.
3. Yana saita `ANTHROPIC_BASE_URL` da `ANTHROPIC_MODEL` idan an ba su, da key ɗin a matsayin `ANTHROPIC_API_KEY` ko `ANTHROPIC_AUTH_TOKEN`.

Asusan da ke da shigar biyan kuɗi suna riƙe environment ɗin da suka gada, kamar da. A zaman mu'amala Claude Code zai iya tambaya sau ɗaya ko zai yi amfani da key ɗin da ya samu a environment.

## Tsaro da yanayi na musamman

- Ba a taɓa rubuta key ɗin a `config.toml`, cache ko log ba, kuma `ccam --dry-run gw` ba ya karanta shi: dry run ba ya gudanar da tushen `command:`.
- `ccam add --api` yana duba tushen `env:` nan take kuma yana gargaɗi idan variable ɗin babu komai a ciki; ana fara gudanar da tushen `command:` ne a lokacin ƙaddamarwa.
- `ccam login` da `ccam logout` suna ƙin asusun API key, domin ba shi da shiga. `ccam doctor` yana nuna endpoint ɗinsa da tushen key ɗinsa, kuma yana gargaɗi idan variable na `env:` babu komai a ciki.
- `ccam usage` yana nuna `API` gare shi. `ccam next --auto`, `--fallback` da `ccam pick` suna zaɓa a cikin asusan biyan kuɗi kawai, domin API key ba shi da iyakar awa 5 ko ta mako da za a kwatanta.
- `settings.json` da aka raba mai sashen `env` ko `apiKeyHelper` yana aiki ga kowane asusu a rukunin, har da asusan API key. Ka cire saitunan mai bayarwa daga saitunan da aka raba.

## Masu alaƙa

- [Sarrafa asusai](/features/accounts)
- [Fara Claude Code](/features/launch)
- [Bayanin umarnai](/docs/commands)
