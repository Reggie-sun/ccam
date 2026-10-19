---
title: એકાઉન્ટ વચ્ચે શેરિંગ
description: ccam share કેવી રીતે સેટિંગ્સ, CLAUDE.md, એક્સટેન્શન્સ અને ઇતિહાસને યોજના, બૅકઅપ અને ટકરાવની પસંદગીઓ સાથે ગ્રુપ ફોલ્ડરમાં લિંક કરે છે, અને off કેવી રીતે પાછું કૉપિ કરે છે.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

દરેક એકાઉન્ટ પાસે પોતાનું ફોલ્ડર છે, પણ તેમાંનું મોટા ભાગનું અલગ રાખવાની જરૂર નથી. `ccam share` એક સમયે એક એકાઉન્ટ અને એક આઇટમ માટે શેરિંગ ચાલુ કે બંધ કરે છે: એકાઉન્ટની ફાઇલ કે ફોલ્ડર ગ્રુપ ફોલ્ડરમાં ખસેડાય છે અને તેની જગ્યાએ symbolic link મુકાય છે, તેથી ગ્રુપનું દરેક એકાઉન્ટ એક જ કૉપિ વાંચે અને લખે છે. કોઈ એકાઉન્ટ શું શેર કરે તે બદલવા, કે થોડાં એકાઉન્ટને તેમના પોતાના ગ્રુપમાં મૂકવા, `ccam add` પછી તેને વાપરો.

## ઉપયોગ

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

આઇટમ ઓપ્શન વગર, `ccam share` હાલની સ્થિતિ છાપે છે:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

ઓપ્શન સાથે, ccam યોજના છાપે છે, ગ્રુપની કૉપિથી જુદી હોય તેવી દરેક ફાઇલ વિશે પૂછે છે, અને કંઈ બદલતાં પહેલાં `Proceed? [y/N]` પૂછે છે. નવું એકાઉન્ટ વાતચીતના ઇતિહાસ સિવાય બધું શેર કરે છે; [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing) દરેક આઇટમને વપરાશકર્તાની દૃષ્ટિએ વર્ણવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP સર્વર | બદલાતું નથી |
| `--settings <on\|off\|group>` | `settings.json` | બદલાતું નથી |
| `--memory <on\|off\|group>` | ગ્લોબલ `CLAUDE.md` | બદલાતું નથી |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | બદલાતું નથી |
| `--history <on\|off\|group>` | `projects/` અને `history.jsonl` | બદલાતું નથી |
| `--fresh` | શેરિંગ બંધ કરતી વખતે, ગ્રુપની સામગ્રી કૉપિ કરવાને બદલે ખાલીથી શરૂ કરે છે | off |
| `--prefer <group\|mine\|merge>` | ટકરાવ કેવી રીતે નક્કી કરવા | પૂછે છે; `--yes` સાથે, દરેક આઇટમનો ડિફૉલ્ટ |
| `-y`, `--yes` | પૂછ્યા વગર ચલાવે છે; ટકરાવ `--prefer` પ્રમાણે નક્કી થાય છે | off |
| `-h`, `--help` | `share` ની મદદ | |
| `--dry-run` (global) | માનવામાં આવે છે: ટકરાવ વિશે પૂછે છે, યોજના છાપે છે, કંઈ બદલતું નથી | off |

મૂલ્ય `on` નો અર્થ ગ્રુપ `default` છે, `off` નો અર્થ local (શેર નહીં) છે, અને બીજું કોઈ પણ મૂલ્ય ગ્રુપનું નામ છે. ગ્રુપનાં નામ `~/.local/share/ccam/groups/` હેઠળ ફોલ્ડરનાં નામ બને છે, તેથી `ccam share side --history ../x` `invalid group name: "../x"` સાથે નિષ્ફળ જાય છે. એકાઉન્ટ નામથી કે નંબરથી આપવામાં આવે છે, અને ફક્ત તમે નામ આપેલી આઇટમ જ બદલાય છે.

## તે કેવી રીતે કામ કરે છે

સંકળાયેલી ફાઇલો:

- `~/.config/ccam/config.toml`: દરેક એકાઉન્ટનું `[profile.share]` table, આઇટમ દીઠ એક મૂલ્ય.
- `~/.local/share/ccam/profiles/<name>/`: એકાઉન્ટનું `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<group>/`: શેર કરેલી કૉપિ.

1. `internal/cli/share.go` config વાંચે છે, `config.Share.Set` વડે દરેક `--<item>` મૂલ્ય એકાઉન્ટનાં સેટિંગ્સની કૉપિ પર લાગુ કરે છે, અને એકાઉન્ટ ચાલતું હોય તો ઇનકાર કરે છે.
2. `share.Build` હાલનાં સેટિંગ્સની લક્ષ્ય સાથે સરખામણી કરે છે અને પગલાંની યાદી બનાવે છે. તે ફક્ત ડિસ્ક વાંચે છે: `share.Layout` ના દરેક પાથ માટે `os.Lstat` અને `os.Readlink`. યોજના operations ની યાદી છે: ગ્રુપની ખાલી જગ્યામાં ખસેડવું (adopt), ગ્રુપમાં ન હોય તેવી ફાઇલો ખસેડવી (merge dir), ઉકેલાયેલો ટકરાવ ગ્રુપમાં લખવો, બૅકઅપ લેવો, unlink, link, પાછું કૉપિ કરવું, plugin પાથ ફરીથી લખવા અને plugin registry મર્જ કરવી.
3. ટકરાવ પ્રશ્ન દ્વારા કે `--prefer` દ્વારા નક્કી થાય છે.
4. `Plan.Describe` યોજના છાપે છે, અને ccam પૂછે છે.
5. કંઈ પણ ખસે તે પહેલાં લક્ષ્ય સેટિંગ્સ `config.toml` માં લખાય છે. જો `Plan.Apply` અધવચ્ચે અટકે, તો `ccam doctor --fix` અને ફરીથી ચલાવવું સેવ કરેલા લક્ષ્ય તરફ પોતાની યોજના બનાવે છે અને ફેરફાર પૂરો કરે છે.
6. `Plan.Apply` પગલાં ક્રમમાં ચલાવે છે. દરેક પગલું ડેટાને ગ્રુપમાં કે બૅકઅપ ફોલ્ડરમાં ખસેડે છે; કંઈ કાઢી નાખવામાં આવતું નથી.

`share.Layout` દરેક આઇટમમાં આવતા પાથની મંજૂર યાદી છે. તેમાં ન હોય તેવું કંઈ પણ ક્યારેય ગ્રુપ સુધી પહોંચી શકતું નથી. `share.Never` `.credentials.json` અને `.claude.json` નાં નામ આપે છે, અને જો તેમાંથી કોઈ ક્યારેય layout માં દેખાય તો `Build` આંતરિક ભૂલ સાથે અટકે છે. `mcp` આઇટમ પાસે કોઈ પાથ નથી: સર્વરની યાદી `.claude.json` ની અંદર રહે છે, તેથી [MCP શેરિંગ](/features/mcp) તેના બદલે ગ્રુપની `mcp.json` દ્વારા કામ કરે છે, અને `--mcp` ફક્ત સેટિંગ બદલે છે.

જ્યારે `side` ઇતિહાસ ચાલુ કરે અને ગ્રુપ પાસે પહેલેથી સામગ્રી હોય, ત્યારે એક ફોલ્ડર, `projects/`, સાથે આ થાય છે:

<!-- motion:share-on -->

જો ગ્રુપ પાસે એ પાથ પર હજી કંઈ ન હોય, તો એકાઉન્ટ તેને શેર કરનાર પહેલું છે અને તેનું ફોલ્ડર જેમ છે તેમ, બૅકઅપ વગર, ત્યાં ખસેડાય છે. નહીં તો `mergeDir` ગ્રુપમાં ન હોય તેવી દરેક ફાઇલ ગ્રુપમાં ખસેડે છે, અને બંને બાજુ હોય તેવી ફાઇલો એકાઉન્ટના ફોલ્ડરમાં રહે છે, તેથી ગ્રુપનું વર્ઝન રહે છે. પછી એકાઉન્ટનું ફોલ્ડર એકાઉન્ટ ફોલ્ડરની અંદર `.ccam-backup/<time>/` માં જાય છે, અને `os.Symlink` તેની જગ્યાએ લિંક મૂકે છે. વાતચીતની ફાઇલોનાં નામ તેમના session ID પરથી હોય છે, તેથી ઘણાં એકાઉન્ટનો ઇતિહાસ કોઈ ટકરાવ વગર ભેગો થાય છે.

ખસેડવા માટે `rename(2)` વપરાય છે, અને જુદી file system વચ્ચે તે કૉપિ અને delete પર આધાર રાખે છે. ખસેડેલા tree ની અંદરની જે relative link તેની બહાર નિર્દેશ કરતી હતી તે absolute link તરીકે ફરીથી લખાય છે (`fsutil.PortableLink`), જેથી નવી જગ્યાએથી પણ તે સાચી જગ્યાએ પહોંચે. `plugins/` ની registry (`installed_plugins.json`, `known_marketplaces.json`) એકાઉન્ટ ફોલ્ડરમાંથી પસાર થતા absolute પાથ સાચવે છે; ccam તેમને ગ્રુપના ફોલ્ડર તરફ વાળે છે અને, ગ્રુપ પાસે પહેલેથી registry હોય ત્યારે, એકાઉન્ટની એન્ટ્રીઓ તેમાં મર્જ કરે છે.

કોઈ આઇટમ બંધ કરવાથી લિંક ઊલટી થાય છે:

<!-- motion:share-off -->

ccam લિંક ફક્ત ત્યારે જ દૂર કરે છે જ્યારે તે groups ફોલ્ડરની અંદર નિર્દેશ કરતી હોય, પછી `CopyTree` વડે ગ્રુપની હાલની સામગ્રી એકાઉન્ટમાં કૉપિ કરે છે. `--fresh` સાથે કૉપિ છોડી દેવાય છે અને એકાઉન્ટ એ પાથ પર કંઈ વગર શરૂ થાય છે. બીજાં એકાઉન્ટ માટે ગ્રુપ પોતાની સામગ્રી રાખે છે. કોઈ આઇટમને એક ગ્રુપમાંથી બીજામાં ખસેડવાથી તે unlink થાય છે અને નવા ગ્રુપ સાથે લિંક થાય છે, જૂનું ગ્રુપ જેમ છે તેમ રહે છે; યોજના `Note:` લાઇનમાં એ જણાવે છે.

બંને બાજુ જુદી સામગ્રી સાથે હોય તેવી એક ફાઇલ ટકરાવ છે:

<!-- motion:share-conflict -->

`share.NewConflict` ફાઇલને બંધબેસતી પસંદગીઓની યાદી આપે છે. `history.jsonl` `merge by line`, `keep group` અને `use mine` આપે છે, જેમાં મર્જ પહેલું છે. બીજી `.json` ફાઇલો `keep group`, `use mine` અને `merge JSON` આપે છે. બાકીનું બધું `keep group` અને `use mine` આપે છે. પહેલી પસંદગી ડિફૉલ્ટ છે:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) objects ને recursive રીતે મર્જ કરે છે, arrays ને duplicate વગર જોડે છે, અને બંને બાજુ એક જ key જુદી રીતે સેટ હોય ત્યારે ગ્રુપનું મૂલ્ય રાખે છે, તેથી બીજાં એકાઉન્ટને વર્તનમાં કોઈ ફેરફાર દેખાતો નથી. `merge by line` ગ્રુપની લાઇનો રાખે છે, પછી તેમાં ન હોય તેવી તમારી લાઇનો ઉમેરે છે. `keep group` સિવાયની દરેક પસંદગી પહેલાં ગ્રુપની કૉપિ એકાઉન્ટ ફોલ્ડરમાં `.ccam-backup/<time>/group/` માં સેવ કરે છે, પછી પરિણામ `fsutil.AtomicWrite` વડે લખે છે. તમારી પોતાની ફાઇલ હંમેશાં બૅકઅપમાં જાય છે.

`--prefer group`, `mine` કે `merge` પૂછ્યા વગર દરેક ટકરાવ નક્કી કરે છે; `merge` નો અર્થ `merge JSON` કે `merge by line` છે, ફાઇલ જે આપે તે, અને જે ફાઇલ બેમાંથી એક પણ ન આપે તે ડિફૉલ્ટ રાખે છે. એકલું `--yes` દરેક ફાઇલનો ડિફૉલ્ટ લે છે.

દરેક શરૂઆતે, `share.Check` ફક્ત `lstat` અને `readlink` વાપરીને layout ના દરેક પાથની સેટિંગ સાથે સરખામણી કરે છે, જે તેને શરૂઆતના પાથ માટે પૂરતું હળવું રાખે છે. ખૂટતી લિંક ત્યાં જ બને છે, કારણ કે તેમાં ગુમાવવા જેવું કંઈ નથી. જ્યાં લિંક હોવી જોઈએ ત્યાં સાચી ફાઇલ, ખોટી જગ્યાની લિંક, કે stale લિંક શરૂઆત અટકાવે છે:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` બંને બાજુ સેવ કરેલાં સેટિંગ્સ સાથે `share.Build` ચલાવે છે, જે ડિસ્ક પર જે છે તેના પરથી દરેક પગલું નક્કી કરે છે, ટકરાવમાં ગ્રુપની સામગ્રી રાખે છે અને જે બદલે તેનો બૅકઅપ લે છે.

## સલામતી અને ખાસ કિસ્સા

- **ચાલતાં એકાઉન્ટ.** Claude Code એકાઉન્ટ તરીકે ચાલતું હોય ત્યારે ccam શેરિંગ બદલવાનો ઇનકાર કરે છે (`account work is running; quit it first, then change its sharing`), કારણ કે સેશન ખસેડાતી ફાઇલોમાં લખવાનું ચાલુ રાખે છે અને એ દરમિયાન ઉમેરેલી લાઇનો ગુમાઈ જાય. `--dry-run` તો પણ યોજના છાપે છે.
- **બૅકઅપ.** દરેક યોજનાને એકાઉન્ટ ફોલ્ડરમાં પોતાનું `.ccam-backup/<time>/` ફોલ્ડર મળે છે, અને બે યોજના એક જ સેકન્ડમાં શરૂ થાય ત્યારે ક્રમાંકિત પ્રત્યય સાથે. ccam તેમને ક્યારેય કાઢી નાખતું નથી.
- **આંશિક નિષ્ફળતા.** જો કોઈ પગલું નિષ્ફળ જાય, તો ccam ભૂલ અને મૂળ ફાઇલો ક્યાં છે તે છાપે છે, અને `ccam doctor --fix` વડે પૂરું કરવા કહે છે. લક્ષ્ય પહેલેથી સેવ થયેલું છે, તેથી સમારકામ તેની તરફ આગળ વધે છે.
- **તમારી પોતાની લિંક.** જો કોઈ પાથ પહેલેથી તમે બનાવેલી લિંક હોય, ઉદાહરણ તરીકે dotfiles repository ની, અને ગ્રુપ ખાલી હોય, તો ccam લિંકનું લક્ષ્ય ગ્રુપમાં કૉપિ કરે છે અને લિંકને પોતે બૅકઅપમાં ખસેડે છે. લક્ષ્યને અડવામાં આવતું નથી.
- **લૉગિન ફાઇલો.** `.credentials.json` અને `.claude.json` ક્યારેય ગ્રુપમાં લિંક, ખસેડવામાં કે કૉપિ કરવામાં આવતી નથી.
- **એક ફોલ્ડરમાં બે એકાઉન્ટ.** ઇતિહાસ શેર કરતાં એકાઉન્ટ એક જ ફાઇલોમાં લખે છે. જ્યારે એ જ ઇતિહાસ ગ્રુપનું બીજું એકાઉન્ટ એવા ફોલ્ડરમાં શરૂ થાય જ્યાં એક પહેલેથી ચાલે છે, ત્યારે ccam ચેતવણી આપે છે અને આગળ વધતાં પહેલાં પૂછે છે.

## સંબંધિત

- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [શેર કરેલા MCP સર્વર](/features/mcp)
- [હાલનું સેટઅપ કૉપિ કરો](/features/migrate)
- [બધા કમાન્ડ](/docs/commands)
