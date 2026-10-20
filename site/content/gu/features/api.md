---
title: API key એકાઉન્ટ
description: ccam add --api Anthropic ના API કે gateway માટે API key એકાઉન્ટ ઉમેરે છે; key લૉન્ચ વખતે variable કે command માંથી વંચાય છે, ક્યારેય સેવ થતી નથી.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API key એકાઉન્ટ Claude subscription લૉગિનને બદલે API key વડે Claude Code શરૂ કરે છે: Anthropic નો API, અથવા એ જ Messages format બોલતું endpoint, જેમ કે LiteLLM કે OpenRouter gateway. ccam key સેવ કરતું નથી. એકાઉન્ટ શરૂ થાય ત્યારે દર વખતે તે environment variable માંથી કે password manager જેવા command માંથી તેને વાંચે છે, અને ફક્ત એ Claude Code process ને જ આપે છે.

## ઉપયોગ

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

API key એકાઉન્ટ બીજાં કોઈ પણ એકાઉન્ટની જેમ સેટિંગ્સ, `CLAUDE.md`, extensions અને MCP સર્વર શેર કરે છે, અને શૉર્ટકટ, `ccam folder` અને `ccam ext` સાથે કામ કરે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--api` | લૉગિન કરવાને બદલે API key એકાઉન્ટ ઉમેરે છે | off |
| `--key-from env:NAME` | environment variable `NAME` માંથી key વાંચે છે | |
| `--key-from command:COMMAND` | `COMMAND` ચલાવે છે (`sh -c` વડે, કે Windows પર `cmd /C` વડે) અને તે છાપે તે પહેલી લાઇન વાપરે છે | |
| `--base-url <url>` | endpoint (`ANTHROPIC_BASE_URL`) | Anthropic નો API |
| `--auth api-key\|bearer` | key ને `ANTHROPIC_API_KEY` (`x-api-key`) તરીકે કે `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, જે gateways સામાન્ય રીતે માગે છે) તરીકે મોકલે છે | `api-key` |
| `--model <model>` | મૉડલ (`ANTHROPIC_MODEL`) | Claude Code નું ડિફૉલ્ટ |

સેટિંગ્સ `config.toml` માં સેવ થાય છે:

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

## તે કેવી રીતે કામ કરે છે

એકાઉન્ટના દરેક લૉન્ચ વખતે, સામાન્ય તપાસો પછી અને Claude Code શરૂ થાય તે પહેલાં:

1. ccam key વાંચે છે: `env:` environment માંથી, `command:` command ને 30 સેકન્ડની મર્યાદા સાથે ચલાવીને. command ટર્મિનલ પર password માગી શકે છે; તેનું error output પણ ટર્મિનલ પર જાય છે. ખાલી key લૉન્ચ અટકાવે છે, અને સંદેશ કહે છે કે key ક્યાં અપેક્ષિત હતી.
2. તે વારસામાં મળેલા એવા દરેક variable ને દૂર કરે છે જે નક્કી કરે છે કે Claude Code કેવી રીતે સાઇન ઇન કરે કે કયો provider વાપરે: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` અને `CLAUDE_CODE_USE_FOUNDRY`. તમારા shell માં બીજા કામ માટે સેટ કરેલી key આ એકાઉન્ટ સુધી પહોંચી શકતી નથી, અને subscription લૉગિન તેની key પર જીતી શકતું નથી.
3. આપેલાં હોય ત્યારે તે `ANTHROPIC_BASE_URL` અને `ANTHROPIC_MODEL` સેટ કરે છે, અને key ને `ANTHROPIC_API_KEY` કે `ANTHROPIC_AUTH_TOKEN` તરીકે.

subscription લૉગિનવાળાં એકાઉન્ટ પહેલાંની જેમ વારસામાં મળેલું environment રાખે છે. interactive સેશનમાં Claude Code એક વાર પૂછી શકે છે કે environment માં મળેલી key વાપરવી કે નહીં.

## સલામતી અને ખાસ કિસ્સા

- key ક્યારેય `config.toml`, cache કે log માં લખાતી નથી, અને `ccam --dry-run gw` તેને વાંચતું નથી: dry run `command:` source ચલાવતું નથી.
- `ccam add --api` `env:` source તરત તપાસે છે અને variable ખાલી હોય તો ચેતવણી આપે છે; `command:` source પહેલી વાર લૉન્ચ વખતે ચાલે છે.
- `ccam login` અને `ccam logout` API key એકાઉન્ટનો ઇનકાર કરે છે, કારણ કે તેનું કોઈ લૉગિન નથી. `ccam doctor` તેનું endpoint અને key source બતાવે છે, અને `env:` variable ખાલી હોય ત્યારે ચેતવણી આપે છે.
- `ccam usage` તેના માટે `API` બતાવે છે. `ccam next --auto`, `--fallback` અને `ccam pick` ફક્ત subscription એકાઉન્ટમાંથી પસંદ કરે છે, કારણ કે API key પાસે સરખાવવા માટે કોઈ 5 કલાકની કે સાપ્તાહિક મર્યાદા નથી.
- `env` block કે `apiKeyHelper` વાળી શેર કરેલી `settings.json` group ના દરેક એકાઉન્ટને લાગુ પડે છે, API key એકાઉન્ટ સહિત. provider સેટિંગ્સ શેર કરેલાં સેટિંગ્સથી બહાર રાખો.

## સંબંધિત

- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [Claude Code શરૂ કરવું](/features/launch)
- [કમાન્ડ સંદર્ભ](/docs/commands)
