---
title: કમાન્ડ સંદર્ભ
description: ccam ના બધા કમાન્ડ અને ઓપ્શન, ccam add અને ccam <name> થી લઈને ccam usage, ccam next અને ccam doctor સુધી.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

કોઈ પણ કમાન્ડ ખરેખર ચલાવ્યા વગર તે શું કરશે તે જોવા માટે તેની આગળ `--dry-run` મૂકો, ઉદાહરણ તરીકે `ccam --dry-run work`.

<!-- commands -->

## Claude Code શરૂ કરવું

`ccam <name|number> [arguments]` તે એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી; તે તેમને એ જ ક્રમમાં `claude` ને આપે છે, અને standard input, exit code અને સિગ્નલ બરાબર `claude` જેવું જ વર્તન કરે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે.

આર્ગ્યુમેન્ટ્સ સાથે પણ એકાઉન્ટના નામ વગરનું `ccam`, જેમ કે `ccam -p "summarise"`, ડિફૉલ્ટ એકાઉન્ટ વાપરે છે. તેને `ccam use <name>` વડે સેટ કરો.

`claude1`, `claude2` વગેરે શૉર્ટકટ `~/.local/bin` માં ccam ની લિંક છે. `ccam shims install --names` દરેક એકાઉન્ટ માટે `claude-<name>` પણ બનાવે છે.

## એકાઉન્ટ ઉમેરવાં અને દૂર કરવાં

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ફોલ્ડર વગરનું `ccam import` `~/.claude` અને `~/.claude-*` શોધે છે અને દરેક વિશે પૂછે છે. ઇમ્પોર્ટ કરેલાં ફોલ્ડર જ્યાં છે ત્યાં જ રજિસ્ટર થાય છે; ccam ક્યારેય લૉગિન ખસેડતું કે કૉપિ કરતું નથી.

## શેરિંગ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

દરેક આઇટમ `on` (`default` ગ્રુપ), `off` કે કોઈ ગ્રુપનું નામ લે છે. [શેરિંગ](/docs/sharing) જુઓ.

## MCP સર્વર

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## વપરાશ અને સ્વિચિંગ

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## જાળવણી

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ફાઇલો

| પાથ | સામગ્રી |
| --- | --- |
| `~/.config/ccam/config.toml` | એકાઉન્ટ અને તેઓ શું શેર કરે છે |
| `~/.local/share/ccam/profiles/<name>/` | દરેક એકાઉન્ટનું `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | શેર કરેલાં સેટિંગ્સ, MCP સર્વર, એક્સટેન્શન્સ અને ઇતિહાસ |
| `~/.local/state/ccam/` | વપરાશ કૅશ અને છેલ્લે વાપર્યાના સમય |
| `~/.local/bin/claude1`, `claude2`… | ccam ની લિંક |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` અને `XDG_STATE_HOME` વેરિએબલ આ ફોલ્ડરનું સ્થાન બદલે છે.
