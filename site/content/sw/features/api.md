---
title: Akaunti za API key
description: ccam add --api huongeza akaunti ya API key kwa API ya Anthropic au gateway; key husomwa wakati wa kuanza kutoka kwa kigezo au amri, haihifadhiwi kamwe.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Akaunti ya API key huanzisha Claude Code kwa API key badala ya kuingia kwa usajili wa Claude: API ya Anthropic, au endpoint inayozungumza muundo ule ule wa Messages, kama gateway ya LiteLLM au OpenRouter. ccam haihifadhi key. Huisoma kila akaunti inapoanza, kutoka kwa kigezo cha mazingira au kutoka kwa amri kama kidhibiti cha nywila, na huikabidhi kwa mchakato huo wa Claude Code pekee.

## Matumizi

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

Akaunti ya API key hushiriki mipangilio, `CLAUDE.md`, viendelezi na seva za MCP kama akaunti nyingine yoyote, na hufanya kazi na njia fupi, `ccam folder` na `ccam ext`.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--api` | Huongeza akaunti ya API key badala ya kuingia | imezimwa |
| `--key-from env:NAME` | Husoma key kutoka kwa kigezo cha mazingira `NAME` | |
| `--key-from command:COMMAND` | Huendesha `COMMAND` (kwa `sh -c`, au `cmd /C` kwenye Windows) na hutumia mstari wa kwanza inaouchapisha | |
| `--base-url <url>` | Endpoint (`ANTHROPIC_BASE_URL`) | API ya Anthropic |
| `--auth api-key\|bearer` | Hutuma key kama `ANTHROPIC_API_KEY` (`x-api-key`) au kama `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, ambayo gateway kwa kawaida huhitaji) | `api-key` |
| `--model <model>` | Modeli (`ANTHROPIC_MODEL`) | Chaguo-msingi la Claude Code |

Mipangilio huhifadhiwa katika `config.toml`:

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

## Jinsi inavyofanya kazi

Kila akaunti inapoanzishwa, baada ya ukaguzi wa kawaida na kabla Claude Code haijaanza:

1. ccam husoma key: `env:` kutoka kwa mazingira, `command:` kwa kuendesha amri yenye kikomo cha sekunde 30. Amri inaweza kuomba nywila kwenye terminal; output yake ya makosa pia huenda kwenye terminal. Key tupu husimamisha uanzishaji kwa ujumbe unaosema key ilitarajiwa wapi.
2. Huondoa kila kigezo kilichorithiwa kinachochagua jinsi Claude Code inavyoingia au mtoa huduma inayemtumia: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` na `CLAUDE_CODE_USE_FOUNDRY`. Key iliyowekwa kwenye shell yako kwa jambo jingine haiwezi kufikia akaunti hii, na kuingia kwa usajili hakuwezi kushinda key yake.
3. Huweka `ANTHROPIC_BASE_URL` na `ANTHROPIC_MODEL` vinapotolewa, na key kama `ANTHROPIC_API_KEY` au `ANTHROPIC_AUTH_TOKEN`.

Akaunti zenye kuingia kwa usajili huhifadhi mazingira zinayorithi, kama hapo awali. Katika kipindi shirikishi Claude Code inaweza kuuliza mara moja kama itumie key inayoipata katika mazingira.

## Usalama na hali maalum

- Key kamwe haiandikwi kwenye `config.toml`, akiba au kumbukumbu (log), na `ccam --dry-run gw` haiisomi: chanzo cha `command:` hakiendeshwi na dry run.
- `ccam add --api` hukagua chanzo cha `env:` papo hapo na huonya ikiwa kigezo ni tupu; chanzo cha `command:` huendeshwa kwa mara ya kwanza wakati wa kuanzisha.
- `ccam login` na `ccam logout` hukataa akaunti ya API key, kwa kuwa haina kuingia. `ccam doctor` huonyesha endpoint yake na chanzo cha key, na huonya kigezo cha `env:` kikiwa tupu.
- `ccam usage` huonyesha `API` kwa ajili yake. `ccam next --auto`, `--fallback` na `ccam pick` huchagua kati ya akaunti za usajili pekee, kwa sababu API key haina kikomo cha saa 5 wala cha wiki cha kulinganisha.
- `settings.json` inayoshirikiwa yenye sehemu ya `env` au `apiKeyHelper` hutumika kwa kila akaunti katika kundi, pamoja na akaunti za API key. Weka mipangilio ya mtoa huduma nje ya mipangilio inayoshirikiwa.

## Zinazohusiana

- [Kusimamia akaunti](/features/accounts)
- [Kuanzisha Claude Code](/features/launch)
- [Marejeo ya amri](/docs/commands)
