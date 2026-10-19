---
title: કમાન્ડ સંદર્ભ
description: ccam ના બધા કમાન્ડ અને ઓપ્શન, ccam add અને ccam <name> થી લઈને ccam usage, ccam next અને ccam doctor સુધી.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a5a33b1a2f4
---

કોઈ પણ કમાન્ડ ખરેખર ચલાવ્યા વગર તે શું કરશે તે જોવા માટે તેની આગળ `--dry-run` મૂકો, ઉદાહરણ તરીકે `ccam --dry-run work`.

<!-- commands -->

## Claude Code શરૂ કરવું

`ccam <name|number> [arguments]` તે એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી; તે તેમને એ જ ક્રમમાં `claude` ને આપે છે, અને standard input, exit code અને સિગ્નલ બરાબર `claude` જેવું જ વર્તન કરે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે.

આર્ગ્યુમેન્ટ્સ સાથે પણ એકાઉન્ટના નામ વગરનું `ccam`, જેમ કે `ccam -p "summarise"`, ડિફૉલ્ટ એકાઉન્ટ વાપરે છે. તેને `ccam use <name>` વડે સેટ કરો.

`ccwork` જેવો શૉર્ટકટ તેના એકાઉન્ટને આ જ રીતે શરૂ કરે છે: `ccwork --resume` એ જ કરે છે જે `ccam work --resume` કરે છે. [શૉર્ટકટ](/docs/commands#shortcuts) જુઓ.

## એકાઉન્ટ ઉમેરવાં અને દૂર કરવાં

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ફોલ્ડર વગરનું `ccam import` `~/.claude` અને `~/.claude-*` શોધે છે અને દરેક વિશે પૂછે છે. ઇમ્પોર્ટ કરેલાં ફોલ્ડર જ્યાં છે ત્યાં જ રજિસ્ટર થાય છે; ccam ક્યારેય લૉગિન ખસેડતું કે કૉપિ કરતું નથી.

## શૉર્ટકટ

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

શૉર્ટકટ એ તમે જાતે નામ આપેલો કમાન્ડ છે, જે Claude Code ને એક એકાઉન્ટ તરીકે શરૂ કરે છે. `ccam work` ની જેમ તે દરેક આર્ગ્યુમેન્ટ કોઈ ફેરફાર વગર Claude Code ને આપે છે, તેથી `ccwork --resume` અને `ccwork -p "hi"` અપેક્ષા મુજબ કામ કરે છે. તે `~/.local/bin` માં ccam ની લિંક છે.

શૉર્ટકટ વૈકલ્પિક છે. `ccam add` વિઝાર્ડ શૉર્ટકટ બનાવવો કે નહીં તે પૂછે છે, `ccam add <name> --shortcut <command>` પૂછ્યા વગર એક બનાવે છે, અને `ccam list` દરેક એકાઉન્ટનો શૉર્ટકટ બતાવે છે, અથવા શૉર્ટકટ ન હોય ત્યારે `-`.

નામમાં a-z, 0-9, `.`, `-` અને `_` અક્ષરો વાપરી શકાય છે. નામ કોઈ ccam કમાન્ડ, બીજા એકાઉન્ટનો શૉર્ટકટ, કે તમારા `PATH` પરના બીજા પ્રોગ્રામનું નામ ન હોઈ શકે.

`ccam rename` શૉર્ટકટનું નામ એ જ રાખે છે, અને `ccam remove` એકાઉન્ટની સાથે તેનો શૉર્ટકટ પણ કાઢી નાખે છે. `ccam shims install` બધાં એકાઉન્ટની લિંક ફરીથી બનાવે છે, ઉદાહરણ તરીકે ccam ને બીજે ખસેડ્યા પછી.

## શેરિંગ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

દરેક આઇટમ `on` (`default` ગ્રુપ), `off` કે કોઈ ગ્રુપનું નામ લે છે. [શેરિંગ](/docs/sharing) જુઓ.

## હાલનું સેટઅપ કૉપિ કરવું

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` કોઈ Claude Code ફોલ્ડરમાંથી `settings.json`, `CLAUDE.md`, એક્સટેન્શન્સ અને user MCP સર્વર એક ગ્રુપમાં કૉપિ કરે છે, જેથી તે આઇટમ શેર કરતા દરેક એકાઉન્ટને તે મળે. ફોલ્ડર ફક્ત વંચાય છે. `--settings`, `--memory`, `--extensions`, `--mcp` અને `--history` ઓપ્શન્સ ફક્ત તમે નામ આપેલી આઇટમ કૉપિ કરે છે; તેમના વગર વાતચીતના ઇતિહાસ સિવાય બધું કૉપિ થાય છે. `-g` કે `--group` ગ્રુપ પસંદ કરે છે, અને `-y` કે `--yes` પૂછ્યા વગર કૉપિ કરે છે અને કોઈ ફાઇલ જુદી હોય ત્યારે ગ્રુપની કૉપિ રાખે છે. [હાલનું સેટઅપ કૉપિ કરવું](/docs/sharing#copy-an-existing-setup) જુઓ.

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
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
```

ccam ને બદલતાં પહેલાં `ccam update` રિલીઝનો SHA-256 checksum તપાસે છે, અને `--version v0.2.0` કોઈ ચોક્કસ રિલીઝ ઇન્સ્ટૉલ કરે છે. `.deb` પેકેજથી ઇન્સ્ટૉલ કરેલું ccam તેના બદલે apt વડે અપડેટ થાય છે. [અપડેટ](/docs/install#update) જુઓ.

## ફાઇલો

| પાથ | સામગ્રી |
| --- | --- |
| `~/.config/ccam/config.toml` | એકાઉન્ટ અને તેઓ શું શેર કરે છે |
| `~/.local/share/ccam/profiles/<name>/` | દરેક એકાઉન્ટનું `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | શેર કરેલાં સેટિંગ્સ, MCP સર્વર, એક્સટેન્શન્સ અને ઇતિહાસ |
| `~/.local/state/ccam/` | વપરાશ કૅશ અને છેલ્લે વાપર્યાના સમય |
| `~/.local/bin/<shortcut>` | `ccwork` જેવા શૉર્ટકટ કમાન્ડ, ccam ની લિંક |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` અને `XDG_STATE_HOME` વેરિએબલ આ ફોલ્ડરનું સ્થાન બદલે છે.
