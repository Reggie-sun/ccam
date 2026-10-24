---
title: હાલનું સેટઅપ કૉપિ કરવું
description: ccam migrate વડે હાલના ~/.claude નાં સેટિંગ્સ, CLAUDE.md, એક્સટેન્શન્સ અને MCP સર્વર ccam ના ગ્રુપમાં કૉપિ કરો.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: 7bca9efb6eb0
---

જો તમે ccam પહેલાં Claude Code વાપર્યું હોય, તો તમારું સેટઅપ `~/.claude` અને `~/.claude.json` માં છે. `ccam migrate` તેને એક શેર ગ્રુપમાં કૉપિ કરે છે, જેથી ગ્રુપ દ્વારા તે આઇટમ શેર કરતું દરેક એકાઉન્ટ તમારી પાસે પહેલેથી હોય તે સેટઅપ સાથે શરૂ થાય. સ્રોત ફોલ્ડર ફક્ત વંચાય છે, અને સાદું `claude` પહેલાંની જેમ કામ કરતું રહે છે. તમારા પહેલા `ccam add` પછી તેને એક વાર ચલાવો; જ્યારે `~/.claude.json` માં MCP સર્વર હોય, ત્યારે પહેલું `ccam add` તેની ટિપમાં તેમનાં નામ આપે છે.

## ઉપયોગ

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

ટર્મિનલ પર, આઇટમ flag વગર અને `--yes` વગર, ccam પહેલાં ચેકલિસ્ટ બતાવે છે:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (કે `x`) cursor હેઠળની આઇટમ ટૉગલ કરે છે, arrow કી કે `j`/`k` ખસેડે છે, Enter યોજના બતાવે છે, અને Esc, `q` કે Ctrl+C `Cancelled.` છાપે છે. જે આઇટમ માટે સ્રોતમાં કંઈ નથી તેના પર `not found` લખેલું હોય છે. જો તમે બધું અનચેક કરો, તો ccam `Nothing selected.` છાપીને અટકે છે.

પછી ccam યોજના છાપે છે અને કૉપિ કરતાં પહેલાં પૂછે છે:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

ગ્રુપ ફક્ત પોતાનાં એકાઉન્ટ સુધી અને ફક્ત પોતે શેર કરે છે તે આઇટમ માટે જ પહોંચે છે, તેથી સારાંશ ગ્રુપનાં એકાઉન્ટની યાદી આપે છે. જ્યારે ગ્રુપ કૉપિ કરેલી આઇટમમાંથી કોઈ શેર ન કરતું હોય, અથવા કેટલાંક એકાઉન્ટ તેમાં ન હોય, ત્યારે ccam તે બદલતા `ccam group set <group> --add <items>` અને `ccam group join <group> <accounts>` કમાન્ડ છાપે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `[folder]` | જેમાંથી કૉપિ કરવાનું છે તે Claude Code ફોલ્ડર | `~/.claude` |
| `--settings` | `settings.json` કૉપિ કરે છે | off |
| `--memory` | ગ્લોબલ `CLAUDE.md` કૉપિ કરે છે | off |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` કૉપિ કરે છે | off |
| `--mcp` | user MCP સર્વર (`.claude.json` માંનું `mcpServers`) કૉપિ કરે છે | off |
| `--history` | વાતચીતનો ઇતિહાસ (`projects/`, `history.jsonl`) કૉપિ કરે છે | off |
| `-g`, `--group string` | જેમાં કૉપિ કરવાનું છે તે શેર ગ્રુપ | `default` |
| `-y`, `--yes` | પૂછ્યા વગર કૉપિ કરે છે; ટકરાવમાં ગ્રુપની કૉપિ રહે છે | off |
| `--dry-run` (global) | યોજના છાપે છે અને કોઈ પ્રશ્ન કે કૉપિ પહેલાં અટકે છે | off |

કોઈ આઇટમ flag ન હોય, ત્યારે ટર્મિનલ પર ચેકલિસ્ટ નક્કી કરે છે. ટર્મિનલની બહાર, કે `--yes` સાથે, ccam વાતચીતના ઇતિહાસ સિવાય બધું કૉપિ કરે છે. ઇતિહાસ ડિફૉલ્ટ રીતે છોડી દેવાય છે, કારણ કે `projects/` મોટું હોઈ શકે છે.

`--group on` નો અર્થ `default` છે. `--group off` અને `--group local` નો `--group takes a group name, not "off"` સાથે ઇનકાર થાય છે. `ccam migrate` `--dry-run` ને માને છે, અને dry run માં પણ ચેકલિસ્ટ દેખાય છે, જેથી તમે પસંદગીનું પૂર્વાવલોકન કરી શકો.

## તે કેવી રીતે કામ કરે છે

`internal/cli/migrate.go` આઇટમ પસંદ કરે છે, પછી `migrate.Build` સ્રોત અને ગ્રુપ વાંચે છે અને એક યોજના પાછી આપે છે. તમે ખાતરી કરો પછી `Plan.Apply` સુધી કંઈ લખાતું નથી.

1. **સ્રોત શોધો.** `checkSource` એવા પાથનો ઇનકાર કરે છે જે ન હોય, ફોલ્ડર ન હોય, કે જેમાં Claude Code ની કોઈ જાણીતી એન્ટ્રી ન હોય. ફક્ત MCP સર્વર માટે, `~/.claude.json` માં સર્વર હોય ત્યાં સુધી ખાલી કે ગેરહાજર `~/.claude` ચાલે છે. `migrate.ClaudeJSON` ફોલ્ડરની `.claude.json` શોધે છે: તેની અંદર, અથવા અંદર ન હોય તેવા `~/.claude` માટે, તેની બાજુની `~/.claude.json`. `CLAUDE_CONFIG_DIR` સેટ ન હોય ત્યારે Claude Code તેને ત્યાં જ રાખે છે.
2. **આઇટમને પાથ સાથે જોડો.** દરેક પાથ `share.Layout` માંથી આવે છે, એ જ મંજૂર યાદી જે ગ્રુપ વાપરે છે. MCP સર્વર `mcp.LocalServers` વડે `.claude.json` માંથી વંચાય છે, જે આખી ફાઇલ parse કરવાને બદલે `profile.TopLevelValue` scanner વાપરે છે.
3. **કૉપિ કરો, ક્યારેય ખસેડો નહીં.** ફાઇલો `fsutil.AtomicWrite` કે `fsutil.CopyFile` વડે `~/.local/share/ccam/groups/<group>/` માં જાય છે, જે લક્ષ્યની બાજુમાં એક કામચલાઉ ફાઇલ લખીને તેનું નામ બદલે છે, તેથી અટકી ગયેલી કૉપિ ક્યારેય અધૂરી ફાઇલ છોડતી નથી. MCP સર્વર ગ્રુપની `mcp.json` માં તેના lock હેઠળ જાય છે: એ દરમિયાન બીજું ccam ઉમેરે તે સર્વર રહે છે, અને [Remote Control](/guides/remote-control) માટે બનાવેલી કૉપિ પછીથી અપડેટ થાય છે.

દરેક આઇટમ સ્રોતના એક પાથને ગ્રુપના એક પાથ સાથે જોડે છે:

<!-- motion:migrate-items -->

ઇતિહાસની બે હરોળ તૂટક છે, કારણ કે તે ફક્ત `--history` સાથે જ કૉપિ થાય છે. `.credentials.json` અને `.claude.json` `share.Never` માં છે અને ક્યારેય ગ્રુપ સુધી પહોંચતી નથી: જો તેમાંથી કોઈ ક્યારેય layout માં દેખાય તો `Build` આંતરિક ભૂલ આપે છે, અને ફોલ્ડરમાં ફરતી વખતે તેમને નામથી છોડી દેવાય છે.

`buildEntry` દરેક પાથને એક સ્થિતિ આપે છે, જે યોજનાની દરેક લાઇન વર્ણવે છે:

<!-- motion:migrate-states -->

- **Missing**: સ્રોતમાં નથી, અથવા જ્યાં ફોલ્ડર અપેક્ષિત છે ત્યાં ફાઇલ છે. દરેક આઇટમ માટે એક વાર `not found, skipped` તરીકે બતાવાય છે.
- **Linked**: સ્રોતનો પાથ પહેલેથી ગ્રુપની કૉપિ તરફ જાય છે, ઉદાહરણ તરીકે જ્યારે સ્રોત ccam profile ફોલ્ડર હોય.
- **Same**: ગ્રુપ પાસે પહેલેથી એકસરખી સામગ્રી છે.
- **New**: ગ્રુપમાં તે નથી, તેથી આખું કૉપિ થાય છે.
- **Merge**: બંને બાજુ ફોલ્ડર. `copyMissing` ફક્ત ગ્રુપમાં ન હોય તેવી ફાઇલો ઉમેરે છે; જુદી ફાઇલો ગ્રુપમાં છે તેમ જ રહે છે અને યોજના `group copies kept` કહે છે. `plugins/` માટે, registry `installed_plugins.json` અને `known_marketplaces.json` મર્જ થાય છે, અને plugin ના absolute પાથ ગ્રુપના ફોલ્ડર તરફ નિર્દેશ કરે તે રીતે ફરીથી લખાય છે.
- **Conflict**: બંને બાજુ જુદી સામગ્રીવાળી ફાઇલ. તમે દરેક ફાઇલ માટે પસંદ કરો છો: `keep group`, `use mine`, ઉપરાંત `.json` ફાઇલો માટે `merge JSON` (deep merge જેમાં ટકરાતાં મૂલ્યોમાં ગ્રુપ જીતે છે) કે `history.jsonl` માટે `merge by line`. જેની વ્યાખ્યા જુદી હોય તેવો MCP સર્વર પણ `keep group` કે `use mine` સાથેનો ટકરાવ છે. Enter દબાવો ત્યારે યાદીની પહેલી પસંદગી ડિફૉલ્ટ છે; `--yes` હંમેશાં ગ્રુપની કૉપિ રાખે છે.

ચેકલિસ્ટ પોતે `tui.MigratePicker` છે, જે `ccam add` જે વાપરે છે તે `checklist` component પર બનેલું છે; `migrate.Found` lstat call અને MCP સર્વરની યાદી વડે ખૂટતી આઇટમ ચિહ્નિત કરે છે.

### MCP સાઇન-ઇન

OAuth થી સાઇન ઇન કરતા MCP સર્વર તેમના ટોકન સર્વર યાદીમાં નહીં, પણ લૉગિન ફાઇલમાં, `.credentials.json` માં `mcpOAuth` હેઠળ રાખે છે. ccam તેમને કૉપિ કરતું નથી. કૉપિ કરેલા સાઇન-ઇનમાં મૂળ જેવો જ refresh token હોય છે, અને પહેલી વાર રિન્યૂ થાય ત્યારે સર્વર એ token બદલી નાખે છે: ત્યારે બાકીનો દરેક ધારક, સ્રોત `~/.claude` પણ, તે સર્વરમાંથી લૉગઆઉટ થઈ જાય છે. તેથી યોજનાની `sign-ins` લાઇન, સ્રોત જેમાં સાઇન ઇન છે તેવા દરેક કૉપિ થયેલા સર્વર માટે, ગ્રુપનાં એ એકાઉન્ટનાં નામ આપે છે જે `/mcp` વડે જાતે સાઇન ઇન કરે છે, અને સારાંશ તે ફરી જણાવે છે. token વગરની એન્ટ્રી સાઇન-ઇન નથી અને યાદીમાં આવતી નથી, ન તો જેની પાસે પહેલેથી પોતાનું સાઇન-ઇન હોય તેવું એકાઉન્ટ, ન તો જેની વ્યાખ્યા ગ્રુપ રાખે છે તેવો સર્વર. કોઈ લૉગિન ફાઇલ લખાતી નથી.

## સલામતી અને ખાસ કિસ્સા

- સ્રોત ફક્ત વંચાય છે. તેમાંની કોઈ ફાઇલ ખસેડાતી, લિંક થતી કે બદલાતી નથી, અને સ્રોતનું લૉગિન ક્યારેય લખાતું નથી, macOS પર પણ, જ્યાં તે Keychain માંથી વંચાઈ શકે છે.
- `use mine`, મર્જ, કે હાલની `mcp.json` માં નવા MCP સર્વર દ્વારા ગ્રુપની ફાઇલ બદલાય તે પહેલાં, ગ્રુપની કૉપિ `<group>/.ccam-backup/<time>/` માં સેવ થાય છે. ફોલ્ડરનું નામ દરેક રન માટે અનોખું છે, તેથી એક જ સેકન્ડમાં થયેલા બે રન એકબીજાના બૅકઅપ પર લખતા નથી. સારાંશ પાથ છાપે છે.
- `--yes` સાથે, દરેક ટકરાવમાં ગ્રુપની કૉપિ રહે છે, અને ccam જણાવે છે કે તેના બદલે કેવી રીતે પસંદ કરવું: એ જ કમાન્ડ `--yes` વગર ચલાવો.
- જો કંઈ બદલાવાનું ન હોય, તો ccam `Nothing to copy.` છાપે છે; જો દરેક ટકરાવમાં ગ્રુપની કૉપિ રહે, તો `Nothing to copy: the group keeps its own copies.`
- જો કૉપિ અધવચ્ચે નિષ્ફળ જાય, તો ccam `The source folder was not changed. Group files replaced so far are saved in` અને પછી બૅકઅપ ફોલ્ડર સાથે અટકે છે. ફરી ચલાવવાથી જે ફાઇલો આખી ન પહોંચી હોય તે કૉપિ થાય છે.
- જ્યારે ઘણાં ગ્રુપ MCP શેર કરતાં હોય, ત્યારે `ccam doctor`, `ccam ext` અને `ccam mcp list` જે કમાન્ડ સૂચવે છે તે `ccam migrate --mcp --group <a|b>` હોય છે.
- સ્રોત ન હોય તો `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` સાથે પૂરું થાય છે.

## સંબંધિત

- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [તપાસ અને સમારકામ](/features/doctor)
- [ઘણાં Claude Code એકાઉન્ટ](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
