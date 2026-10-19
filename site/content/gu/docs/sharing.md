---
title: MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો
description: ccam દરેક લૉગિન અલગ રાખીને Claude Code એકાઉન્ટ વચ્ચે MCP સર્વર, settings.json, CLAUDE.md, agents, skills અને વાતચીતનો ઇતિહાસ કેવી રીતે શેર કરે છે.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

દરેક એકાઉન્ટ પોતાનું લૉગિન અલગ રાખે છે. બાકીનું બધું આઇટમ પ્રમાણે શેર કરી શકાય છે: ccam એકાઉન્ટની ફાઇલ કે ફોલ્ડરને શેર કરેલા ગ્રુપ ફોલ્ડર સાથે લિંક કરે છે, તેથી એક એકાઉન્ટમાં કરેલો ફેરફાર એ ગ્રુપનાં બધાં એકાઉન્ટમાં દેખાય છે.

## શું શેર કરી શકાય

| આઇટમ | ઓપ્શન | શું શેર થાય છે |
| --- | --- | --- |
| MCP સર્વર | `--mcp` | ગ્રુપની `mcp.json`, જે શરૂઆતમાં Claude Code ને આપવામાં આવે છે |
| સેટિંગ્સ | `--settings` | `settings.json`: પરવાનગીઓ, hooks, મૉડલ, સ્ટેટસ લાઇન |
| ગ્લોબલ મેમરી | `--memory` | `CLAUDE.md` |
| એક્સટેન્શન્સ | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| વાતચીતનો ઇતિહાસ | `--history` | `projects/` અને `history.jsonl` |

નવાં એકાઉન્ટ વાતચીતના ઇતિહાસ સિવાય બધું શેર કરે છે.

લૉગિન ફાઇલો `.credentials.json` અને `.claude.json` ક્યારેય શેર થતી નથી. ccam ફક્ત નિશ્ચિત યાદીના પાથ જ લિંક કરે છે, અને આ બે ફાઇલો તે યાદીમાં નથી.

## ગ્રુપ

`on` નો અર્થ `default` ગ્રુપ છે. તમે કોઈ નામ પણ આપી શકો છો, જેથી અમુક એકાઉન્ટ જ કોઈ આઇટમ શેર કરે:

```sh
ccam share work --history team
ccam share side --history team
```

હવે `work` અને `side` એકબીજા સાથે વાતચીત શેર કરે છે, પણ તમારાં બીજાં એકાઉન્ટ સાથે નહીં.

## શેરિંગ ચાલુ કરવું

જ્યારે કોઈ એકાઉન્ટ કોઈ આઇટમ શેર કરવાનું શરૂ કરે છે, ત્યારે ccam:

1. જે દરેક ફાઇલ તે ખસેડશે, કૉપિ કરશે કે લિંક કરશે તેની યોજના છાપે છે, અને શરૂ કરતાં પહેલાં પૂછે છે;
2. જો ગ્રુપ હજી ખાલી હોય તો એકાઉન્ટની ફાઇલો ગ્રુપ ફોલ્ડરમાં ખસેડે છે, નહીં તો ગ્રુપમાં હજી ન હોય તેવી ફાઇલો ઉમેરે છે;
3. જ્યારે એક જ ફાઇલ જુદી હોય, ત્યારે પૂછે છે કે ગ્રુપનું વર્ઝન રાખવું, તમારું વાપરવું કે JSON મર્જ કરવું;
4. તમારી મૂળ ફાઇલ એકાઉન્ટ ફોલ્ડરની અંદર `.ccam-backup/<time>/` માં રાખે છે;
5. મૂળ ફાઇલની જગ્યાએ ગ્રુપની લિંક મૂકે છે.

વાતચીતની ફાઇલોનાં નામ તેમના session ID પરથી હોય છે, તેથી ઘણાં એકાઉન્ટનો ઇતિહાસ કોઈ ટકરાવ વગર ભેગો થાય છે.

## શેરિંગ બંધ કરવું

```sh
ccam share work --history off
```

લિંક દૂર થાય છે અને એકાઉન્ટને ગ્રુપની હાલની સામગ્રીની પોતાની નકલ મળે છે. તેના બદલે ખાલીથી શરૂ કરવા માટે `--fresh` ઉમેરો. બીજાં એકાઉન્ટ માટે ગ્રુપ પોતાની સામગ્રી રાખે છે.

## હાલનું સેટઅપ કૉપિ કરો

જો તમે ccam પહેલાં Claude Code વાપર્યું હોય, તો તમારાં સેટિંગ્સ `~/.claude` માં છે. `ccam migrate` તેમને એક ગ્રુપમાં કૉપિ કરે છે, જેથી તે આઇટમ શેર કરતા દરેક એકાઉન્ટને તે મળે:

```sh
ccam migrate
```

તે શું કૉપિ કરે છે:

- `settings.json` અને ગ્લોબલ `CLAUDE.md`;
- એક્સટેન્શન્સ `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` અને `rules/`;
- તમારા user MCP સર્વર, એટલે કે `~/.claude.json` ની `mcpServers` એન્ટ્રી, ગ્રુપની `mcp.json` માં;
- વાતચીતનો ઇતિહાસ (`projects/` અને `history.jsonl`), પણ ફક્ત તમે `--history` ઉમેરો ત્યારે જ.

સ્રોત ફોલ્ડર ફક્ત વંચાય છે: તેમાંનું કંઈ પણ ખસેડાતું, લિંક થતું કે બદલાતું નથી, અને સાદું `claude` પહેલાંની જેમ જ કામ કરતું રહે છે. લૉગિન ફાઇલો `.credentials.json` અને `.claude.json` ક્યારેય કૉપિ થતી નથી.

કંઈ પણ કૉપિ કરતાં પહેલાં ccam યોજના છાપે છે અને પૂછે છે:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

અંતે ccam ગ્રુપ દ્વારા દરેક આઇટમ શેર કરતાં એકાઉન્ટની યાદી બતાવે છે. જો કેટલાંક એકાઉન્ટ એમ ન કરતાં હોય, તો તેમના માટે શેરિંગ ચાલુ કરતા `ccam share` કમાન્ડ છાપે છે.

જ્યારે ગ્રુપમાં કોઈ ફાઇલની જુદી કૉપિ પહેલેથી હોય, ત્યારે ccam શેરિંગ ચાલુ કરતી વખતે પૂછે છે તેમ જ પૂછે છે: ગ્રુપની કૉપિ રાખવી કે તમારી વાપરવી. `settings.json` અને બીજી JSON ફાઇલો મર્જ પણ કરી શકાય છે, અને `history.jsonl` લાઇન પ્રમાણે. ગ્રુપની અગાઉની કૉપિઓ ગ્રુપ ફોલ્ડરની અંદર `.ccam-backup/<time>/` માં સાચવાય છે.

| ઓપ્શન | અસર |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | ફક્ત તમે નામ આપેલી આઇટમ કૉપિ કરો; તેમના વગર વાતચીતના ઇતિહાસ સિવાય બધું |
| `-g`, `--group <group>` | `default` ને બદલે આ ગ્રુપમાં કૉપિ કરો |
| `-y`, `--yes` | પૂછ્યા વગર કૉપિ કરો; કોઈ ફાઇલ જુદી હોય ત્યારે ગ્રુપની કૉપિ રાખો |

કંઈ પણ કૉપિ કર્યા વગર યોજના જોવા માટે કમાન્ડની આગળ `--dry-run` મૂકો, જેમ કે `ccam --dry-run migrate`. બીજા ફોલ્ડરમાંથી કૉપિ કરવા તેનું નામ આપો: `ccam migrate ~/.claude-work --group work`.

## MCP સર્વર

Claude Code MCP સર્વરને લૉગિન સાથે `.claude.json` માં રાખે છે, તેથી એ ફાઇલ શેર કરી શકાતી નથી. ccam શેર કરેલા સર્વર ગ્રુપની `mcp.json` માં રાખે છે અને Claude Code ને `--mcp-config=<group>/mcp.json` સાથે શરૂ કરે છે. તમારા પોતાના `--mcp-config` ઓપ્શન તેની ઉપર લાગુ પડતા રહે છે.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

સેશનની અંદર `claude mcp add` વડે ઉમેરેલો સર્વર ફક્ત તે એકાઉન્ટ માટે સેવ થાય છે. તમે તે એકાઉન્ટ આગલી વખતે શરૂ કરો ત્યારે ccam સર્વરનું નામ અને તેને ગ્રુપમાં ખસેડવાનો કમાન્ડ બતાવે છે:

```sh
ccam mcp import work gh
```

OAuth થી સાઇન ઇન કરતા MCP સર્વર તેમના ટોકન એકાઉન્ટના લૉગિન સાથે રાખે છે, તેથી દરેક એકાઉન્ટમાં તેમને એક વાર અધિકૃત કરો.

જો તમે ઇચ્છતા હો કે સર્વર `claude mcp list` અને `/mcp` માં સામાન્ય user સર્વર તરીકે દેખાય, તો `~/.config/ccam/config.toml` માં `mcp_strategy = "sync"` સેટ કરો. પછી ccam શરૂ કરતાં પહેલાં તેમને દરેક એકાઉન્ટની `.claude.json` માં લખે છે, Claude Code ના પોતાના file lock હેઠળ અને બૅકઅપ સાથે.

## એક જ ફોલ્ડરમાં બે એકાઉન્ટ

વાતચીતનો ઇતિહાસ શેર કરતાં એકાઉન્ટ એક જ ફાઇલોમાં લખે છે. જો તમે એવા ફોલ્ડરમાં બીજું Claude Code શરૂ કરો જ્યાં એ જ ગ્રુપનું બીજું એકાઉન્ટ પહેલેથી ચાલી રહ્યું હોય, તો ccam તમને ચેતવે છે અને આગળ વધવું કે નહીં તે પૂછે છે.

## જ્યારે કંઈક તૂટે

```sh
ccam doctor
ccam doctor --fix
```

જો કોઈ શેર કરેલો પાથ સામાન્ય ફાઇલ બની ગયો હોય, ઉદાહરણ તરીકે તમે હાથથી બૅકઅપ પાછો મૂક્યો હોય પછી, તો ccam Claude Code શરૂ કરતાં પહેલાં અટકી જાય છે અને `ccam doctor --fix` તરફ નિર્દેશ કરે છે, જે ફાઇલને ગ્રુપમાં મર્જ કરે છે અને બૅકઅપ રાખે છે.
