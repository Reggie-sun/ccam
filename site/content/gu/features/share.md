---
title: ગ્રુપ વડે શેરિંગ
description: ccam નું ગ્રુપ પસંદ કરેલી આઇટમ એક ફોલ્ડર દ્વારા તેનાં એકાઉન્ટ વચ્ચે શેર કરે છે; ccam group યોજના, બૅકઅપ અને ટકરાવની પસંદગી સાથે ડેટા ખસેડે છે.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: e4f05b0fd572
---

દરેક એકાઉન્ટ પાસે પોતાનું ફોલ્ડર છે, પણ તેમાંનું મોટા ભાગનું અલગ રાખવાની જરૂર નથી. **ગ્રુપ** નક્કી કરે છે કે શું શેર થાય અને કોના દ્વારા: ગ્રુપનું દરેક એકાઉન્ટ એ આઇટમને એક જ ફોલ્ડર, `~/.local/share/ccam/groups/<group>/`, સાથે લિંક કરે છે, તેથી બધાં એક જ કૉપિ વાંચે અને લખે છે. લૉગિન ટોકન ક્યારેય શેર થતા નથી.

એક એકાઉન્ટ ઘણાં ગ્રુપમાં હોઈ શકે છે, પણ દરેક આઇટમ ફક્ત એક જ ગ્રુપમાંથી આવે છે. સામાન્ય સેટઅપમાં સેટિંગ્સ અને એક્સટેન્શન્સ માટે એક ગ્રુપ હોય છે, અને બે એકાઉન્ટ વચ્ચે વાતચીતના ઇતિહાસ માટે બીજું:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

તમે ઉમેરો તે પહેલું એકાઉન્ટ `default` બનાવે છે, અને તમે બીજું કંઈ ન કહો તો પછીનાં એકાઉન્ટ તેમાં જોડાય છે.

## ઉપયોગ

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

ડેટા ખસેડતો દરેક કમાન્ડ યોજના છાપે છે, ગ્રુપની કૉપિથી જુદી હોય તેવી દરેક ફાઇલ વિશે પૂછે છે, અને કંઈ બદલતાં પહેલાં `Proceed? [y/N]` પૂછે છે:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

આ જ ફેરફારો [ડેશબોર્ડ](/features/dashboard) ના **Groups** ટૅબ પર છે, અને કોઈ એકાઉન્ટ પર `s` તેનાં ગ્રુપ પસંદ કરે છે.

## આઇટમ

| આઇટમ | તેમાં શું આવે છે |
| --- | --- |
| `mcp` | MCP સર્વર, ગ્રુપની `mcp.json` દ્વારા ([શેર કરેલા MCP સર્વર](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | ગ્લોબલ `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` અને `history.jsonl`, જેથી બીજું એકાઉન્ટ `--continue` કરી શકે |

જે આઇટમ એકાઉન્ટનું કોઈ ગ્રુપ શેર ન કરતું હોય તે એકાઉન્ટના પોતાના ફોલ્ડરમાં રહે છે.

## કમાન્ડ અને ઓપ્શન

| કમાન્ડ | શું કરે છે |
| --- | --- |
| `ccam group` અથવા `ccam group list` | ગ્રુપની યાદી બતાવે છે |
| `ccam group show <group>` | ગ્રુપ શું શેર કરે છે, તેનાં એકાઉન્ટ, તેનું ફોલ્ડર અને MCP સર્વર |
| `ccam group add <group> --items <items> [--members <accounts>]` | ગ્રુપ બનાવે છે; તેનાં એકાઉન્ટની આઇટમ તેમાં ખસેડાય છે |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | ગ્રુપ શું શેર કરે છે તે બદલે છે, તેમાંનાં દરેક એકાઉન્ટ માટે |
| `ccam group join <group> <account>...` | એકાઉન્ટને ગ્રુપમાં મૂકે છે |
| `ccam group leave <group> <account>...` | એકાઉન્ટને બહાર કાઢે છે; દરેક એક કૉપિ રાખે છે, MCP સર્વર સહિત |
| `ccam group rename <group> <new-name>` | ફક્ત config માં નામ બદલે છે |
| `ccam group remove <group> [--leave]` | ગ્રુપ દૂર કરે છે |

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: જે એકાઉન્ટની કૉપિ સૌથી પહેલાં ગ્રુપમાં જાય છે | ડિફૉલ્ટ એકાઉન્ટ, પછી નંબર પ્રમાણે |
| `--replace` | `join`: પહેલાં એ જ આઇટમ શેર કરતાં ગ્રુપ છોડે છે | off |
| `--fresh` | કોઈ આઇટમ શેર કરવાનું બંધ કરતાં એકાઉન્ટ કૉપિને બદલે તેના વગર શરૂ થાય છે (ઇમ્પોર્ટ કરેલું ફોલ્ડર તો પણ કૉપિ રાખે છે) | off |
| `--leave` | `remove`: પહેલાં ગ્રુપનાં એકાઉન્ટને બહાર કાઢે છે | ગ્રુપમાં એકાઉન્ટ હોય ત્યારે જરૂરી |
| `--prefer <group\|mine\|merge>` | ટકરાવ કેવી રીતે નક્કી કરવા | પૂછે છે; `--yes` સાથે, દરેક ફાઇલનો ડિફૉલ્ટ |
| `-y`, `--yes` | પૂછ્યા વગર ચલાવે છે; ટકરાવ `--prefer` પ્રમાણે નક્કી થાય છે | off |
| `--dry-run` (global) | યોજના છાપે છે, કંઈ બદલતું નથી | off |

એકાઉન્ટ નામથી કે નંબરથી આપવામાં આવે છે; આઇટમ અલ્પવિરામથી અલગ કરેલી યાદી તરીકે. `ccam add` `--group <group>` (ઘણી વાર આપી શકાય) કે `--no-group` લે છે.

## તે કેવી રીતે કામ કરે છે

સંકળાયેલી ફાઇલો:

- `~/.config/ccam/config.toml`: `[[group]]` એન્ટ્રીઓ: `name`, `items`, `members` અને, ફોલ્ડરનું નામ જુદું હોય ત્યારે, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: એકાઉન્ટનું `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<folder>/`: શેર કરેલી કૉપિ.

`config.toml` વંચાય ત્યારે દરેક એકાઉન્ટનાં ગ્રુપ આઇટમ દીઠ એક ફોલ્ડરમાં ઉકેલાય છે. પાંચ એકાઉન્ટ માટે આમાં લગભગ 25 µs લાગે છે, તેથી શરૂઆત હજી પણ ફક્ત `config.toml` જ વાંચે છે. જે config માં કોઈ એકાઉન્ટને એક આઇટમ બે ગ્રુપમાંથી મળે તે ભૂલ છે, અને તે સુધારાય નહીં ત્યાં સુધી કંઈ શરૂ થતું નથી.

ફેરફાર config ની કૉપિ પર કામ કરે છે:

1. કમાન્ડ કૉપિમાં ફેરફાર કરે છે (`AddGroup`, `SetGroupItems`, `JoinGroup` …), જે નિયમો તપાસે છે.
2. `internal/group` જૂના અને નવાની સરખામણી કરે છે અને, જે દરેક એકાઉન્ટની આઇટમ બદલાય છે તેના માટે, `share.Build` વડે યોજના બનાવે છે. યોજનાઓ એક પછી એક બને છે, અને દરેક જુએ છે કે અગાઉની યોજનાઓ ગ્રુપ ફોલ્ડરમાં શું મૂકશે, તેથી ફક્ત પહેલું એકાઉન્ટ જ પોતાની કૉપિ અંદર ખસેડે છે. પહેલી કૉપિ ગ્રુપ ફોલ્ડર પાસે પહેલેથી હોય તો ત્યાંથી આવે છે, પછી `--from`, પછી ડિફૉલ્ટ એકાઉન્ટ, પછી બાકીનાં નંબર પ્રમાણે.
3. જો એમાંનું કોઈ એકાઉન્ટ ચાલી રહ્યું હોય, Claude Code ભલે ગમે તે રીતે શરૂ થયું હોય, તો ccam ઇનકાર કરે છે (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), અને Developer Mode વગરના Windows પર જો કોઈ ફાઇલને લિંકની જરૂર પડે તો પણ.
4. ટકરાવ પ્રશ્ન દ્વારા કે `--prefer` દ્વારા નક્કી થાય છે, યોજના છપાય છે, અને ccam પૂછે છે.
5. સેવ કરતાં ઠીક પહેલાં, ccam ફરી તપાસે છે: પ્રશ્નો રાહ જોતા હતા ત્યારે શરૂ થયેલું એકાઉન્ટ ફેરફાર અટકાવે છે, અને એ જ ગ્રુપમાં બીજા ccam નો ફેરફાર પણ (`the groups changed while this plan waited; nothing was moved, make the change again`). કંઈ પણ ખસે તે પહેલાં નવો config સેવ થાય છે, અને ફેરફાર `~/.local/state/ccam/group-change.json` માં નોંધાય છે. જો કોઈ પગલું અધવચ્ચે અટકે, તો નોંધ રહે છે: `ccam doctor --fix` ફેરફારને તેના ક્રમમાં, તમે આપેલા જવાબો સાથે, MCP પગલાં સહિત પૂરો કરે છે, અને ત્યાં સુધી ગ્રુપનો બીજો કોઈ ફેરફાર શરૂ થતો નથી (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``).
6. MCP પગલાં પહેલાં ચાલે છે. જ્યારે એકાઉન્ટ કોઈ ગ્રુપમાંથી extensions બહાર કાઢે, ત્યારે ગ્રુપની plugin યાદીઓના જે પાથ કોઈ એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા હોય, તેમને ગ્રુપ ફોલ્ડર તરફ વાળવામાં આવે છે (યોજનામાં `group plugin lists:`). પછી યોજનાઓ ક્રમમાં ચાલે છે. દરેક પગલું ડેટાને ગ્રુપમાં કે બૅકઅપ ફોલ્ડરમાં ખસેડે છે; ફક્ત ગ્રુપ પાસે બાઇટે બાઇટ પહેલેથી હોય તેવી ફાઇલો બૅકઅપને બદલે કાઢી નખાય છે.

`share.Layout` દરેક આઇટમમાં આવતા પાથની મંજૂર યાદી છે. તેમાં ન હોય તેવું કંઈ પણ ક્યારેય ગ્રુપ સુધી પહોંચી શકતું નથી. `share.Never` `.credentials.json` અને `.claude.json` નાં નામ આપે છે, અને જો તેમાંથી કોઈ ક્યારેય layout માં દેખાય તો `Build` આંતરિક ભૂલ સાથે અટકે છે. `mcp` આઇટમ પાસે કોઈ પાથ નથી: સર્વરની યાદી `.claude.json` ની અંદર રહે છે, તેથી [MCP શેરિંગ](/features/mcp) ગ્રુપની `mcp.json` દ્વારા કામ કરે છે.

જ્યારે `side` ઇતિહાસ શેર કરવાનું શરૂ કરે અને ગ્રુપ પાસે પહેલેથી સામગ્રી હોય, ત્યારે એક ફોલ્ડર, `projects/`, સાથે આ થાય છે:

<!-- motion:share-on -->

જો ગ્રુપ પાસે એ પાથ પર હજી કંઈ ન હોય, તો એકાઉન્ટ તેને શેર કરનાર પહેલું છે અને તેનું ફોલ્ડર જેમ છે તેમ ત્યાં ખસેડાય છે. નહીં તો ગ્રુપમાં ન હોય તેવી દરેક ફાઇલ ગ્રુપમાં ખસેડાય છે, અને બંને બાજુ હોય તેવી ફાઇલો એકાઉન્ટના ફોલ્ડરમાં રહે છે, તેથી ગ્રુપનું વર્ઝન રહે છે. ગ્રુપની ફાઇલો જેવી જ ફાઇલો કાઢી નખાય છે, અને યોજના જુદી પડતી ફાઇલોનાં નામ જણાવે છે (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`). પછી એકાઉન્ટનું ફોલ્ડર એકાઉન્ટ ફોલ્ડરની અંદર `.ccam-backup/<time>/` માં જાય છે, અને તેની જગ્યાએ લિંક મુકાય છે. જો ગ્રુપની એન્ટ્રી પોતે જ જોડાતા એકાઉન્ટની પોતાની એન્ટ્રીની લિંક હોય, તો એકાઉન્ટની સામગ્રી તેની જગ્યાએ ગ્રુપમાં જાય છે. વાતચીતની ફાઇલોનાં નામ તેમના session ID પરથી હોય છે, તેથી ઘણાં એકાઉન્ટનો ઇતિહાસ કોઈ ટકરાવ વગર ભેગો થાય છે. `plugins/` ની registry એકાઉન્ટ ફોલ્ડરમાંથી પસાર થતા absolute પાથ સાચવે છે; ccam તેમને ગ્રુપના ફોલ્ડર તરફ વાળે છે અને ગ્રુપ પાસે પહેલેથી હોય તે registry માં એકાઉન્ટની એન્ટ્રીઓ મર્જ કરે છે.

ગ્રુપ છોડવાથી લિંક ઊલટી થાય છે:

<!-- motion:share-off -->

ccam પહેલાં ગ્રુપની હાલની સામગ્રી એકાઉન્ટ ફોલ્ડરમાંના `.ccam-backup/incomplete-<name>` માં કૉપિ કરે છે, અને કૉપિ પૂરી થાય પછી જ લિંક બદલે છે, તેથી અટકી ગયેલી કૉપિ લિંકને તેની જગ્યાએ રહેવા દે છે અને `ccam doctor --fix` ફરી કૉપિ કરે છે (આવી રીતે પૂરું કરવાની યોજના `restore <entry>: … (copied back from the group before a stop)` બતાવે છે). તે લિંક ફક્ત ત્યારે જ દૂર કરે છે જ્યારે તે, યોજના બની ત્યારે હતી તેમ, હજી groups ફોલ્ડરની અંદર નિર્દેશ કરતી હોય; ત્યાર પછી બદલાયેલો પાથ, જેમ કે જ્યાં લિંક હતી ત્યાં સાચી ફાઇલ, જેમ છે તેમ રહે છે અને ફેરફાર અટકે છે. `--fresh` સાથે કૉપિ છોડી દેવાય છે. બીજાં એકાઉન્ટ માટે ગ્રુપ પોતાની સામગ્રી રાખે છે. જે એકાઉન્ટ કોઈ આઇટમને એવા ગ્રુપમાં ખસેડે જેનું ફોલ્ડર હજી ખાલી હોય, તે પોતે જે જોયું હતું તેની કૉપિ સાથે લાવે છે, તેથી તેની ફાઇલો રહે છે; એ કૉપિ પૂરી ન થાય ત્યાં સુધી `<group>/.ccam-incomplete-<name>` માં તૈયાર રખાય છે. MCP શેર કરવાનું બંધ કરતું એકાઉન્ટ ગ્રુપના સર્વર રાખે છે: તે તેની પોતાની `.claude.json` માં, તેના પોતાના સર્વરની બાજુમાં જાય છે (બૅકઅપ રખાય છે); એ જ નામના તેના પોતાના સર્વર જેમ છે તેમ રહે છે. `--fresh` સાથે, તેના બદલે યોજના લોડ થવાનું બંધ થતા સર્વરની યાદી આપે છે. ફરી જોડાવાથી કોઈએ ન બદલેલી કૉપિ દૂર થાય છે. `extensions` છોડવાથી, ગ્રુપની plugin યાદીઓના જે પાથ છોડતા એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા હતા, તેમને પણ ગ્રુપ ફોલ્ડર તરફ વાળવામાં આવે છે; બીજાં એકાઉન્ટ માટે તેમનો અર્થ એ જ રહે છે.

બંને બાજુ જુદી સામગ્રી સાથે હોય તેવી એક ફાઇલ ટકરાવ છે:

<!-- motion:share-conflict -->

`settings.json` જેવી એકલી ફાઇલો માટે તમે પસંદ કરો છો; ફોલ્ડરની અંદર ગ્રુપની કૉપિ રહે છે, અને યોજના જુદી પડતી ફાઇલોનાં નામ જણાવે છે. `history.jsonl` `merge by line`, `keep group` અને `use mine` આપે છે, જેમાં મર્જ પહેલું છે. બીજી `.json` ફાઇલો `keep group`, `use mine` અને `merge JSON` આપે છે. બાકીનું બધું `keep group` અને `use mine` આપે છે. `merge JSON` objects ને recursive રીતે મર્જ કરે છે, arrays ને duplicate વગર જોડે છે, અને બંને બાજુ એક જ key જુદી રીતે સેટ હોય ત્યારે ગ્રુપનું મૂલ્ય રાખે છે, તેથી બીજાં એકાઉન્ટને કોઈ ફેરફાર દેખાતો નથી. `keep group` સિવાયની દરેક પસંદગી પહેલાં ગ્રુપની કૉપિ એકાઉન્ટના `.ccam-backup/<time>/group/` માં સેવ કરે છે. તમારી પોતાની ફાઇલ હંમેશાં બૅકઅપમાં જાય છે.

દરેક શરૂઆતે, `share.Check` ફક્ત `lstat` અને `readlink` વાપરીને layout ના દરેક પાથની એકાઉન્ટનાં ગ્રુપ સાથે સરખામણી કરે છે. ખૂટતી લિંક ત્યાં જ બને છે. જ્યાં લિંક હોવી જોઈએ ત્યાં સાચી ફાઇલ, ખોટી જગ્યાની લિંક, કે stale લિંક શરૂઆત અટકાવે છે અને `ccam doctor --fix` તરફ નિર્દેશ કરે છે.

### નામ બદલવું, દૂર કરવું અને જૂનું ફૉર્મેટ

- `ccam group rename` ફક્ત config બદલે છે. ફોલ્ડર પોતાનું નામ રાખે છે અને ગ્રુપ તેને `dir` તરીકે નોંધે છે, તેથી કંઈ ખસતું નથી અને કંઈ ગુમાઈ શકતું નથી.
- `ccam group remove` એકાઉન્ટવાળા ગ્રુપનો ઇનકાર કરે છે, સિવાય કે તમે `--leave` ઉમેરો. ત્યારે દરેક એકાઉન્ટ એક કૉપિ રાખે છે, અને બીજું કોઈ ગ્રુપ ફોલ્ડર ન વાપરતું હોય તો તે `~/.local/share/ccam/trash/` માં જાય છે.
- ગ્રુપ પહેલાં, `config.toml` દરેક એકાઉન્ટ માટે `[profile.share]` table રાખતું હતું. ccam એ ફૉર્મેટને કંઈ ખસેડ્યા વગર ગ્રુપ તરીકે વાંચે છે. પહેલી વાર સેવ કરતી વખતે તે નવું ફૉર્મેટ લખે છે અને જૂની ફાઇલ `config.toml.bak` તરીકે રાખે છે. એ જ એકાઉન્ટ દ્વારા શેર થતી આઇટમ ફોલ્ડરના નામનું એક ગ્રુપ બને છે; બાકીનાને `default-history` જેવાં નામ મળે છે અને તે એ જ ફોલ્ડર વાપરતી રહે છે.

### Windows પર

Windows symbolic links ફક્ત Developer Mode (Settings > System > For developers) માં કે ઍડમિનિસ્ટ્રેટર માટે જ બનાવે છે. એ અધિકાર વગર, શેર કરેલું ફોલ્ડર directory junction બને છે, જે Claude Code માટે એ જ રીતે કામ કરે છે. શેર કરેલી ફાઇલ (`settings.json`, `CLAUDE.md`, `history.jsonl`) ફક્ત symbolic link જ હોઈ શકે, તેથી જેને તેની જરૂર પડે તેવો ગ્રુપ ફેરફાર કંઈ ખસે તે પહેલાં નકારાય છે. આ જ કારણે ત્યાં પહેલું `ccam add` ફક્ત `mcp` અને `extensions` સાથે `default` ગ્રુપ બનાવે છે અને તેમ જણાવે છે, અને સેટિંગ્સ, મેમરી કે ઇતિહાસ શેર કરતા ગ્રુપમાં `ccam add --group` નો ઇનકાર થાય છે.

## સલામતી અને ખાસ કિસ્સા

- **ચાલતાં એકાઉન્ટ.** ફેરફાર જે એકાઉન્ટને ખસેડે છે તેમાંના કોઈ એક તરીકે Claude Code ચાલતું હોય ત્યારે ફેરફાર નકારાય છે, કારણ કે સેશન એ ફાઇલોમાં લખવાનું ચાલુ રાખે છે, ભલે ccam એ તેને શરૂ કર્યું હોય કે નહીં (સાદું `claude`, IDE, બૅકગ્રાઉન્ડ સેશન). સેવ કરતાં ઠીક પહેલાં તપાસ ફરી થાય છે. ગ્રુપ છોડવું અને તેમાં જોડાવું ફક્ત નામ આપેલા એકાઉન્ટને લાગુ પડે છે; `extensions` છોડવાથી ગ્રુપની plugin યાદીઓમાંના પાથ પણ ફરી લખાય છે, બીજાં માટે તેમનો અર્થ બદલ્યા વગર. `--dry-run` તો પણ યોજના છાપે છે.
- **બૅકઅપ.** દરેક યોજનાને એકાઉન્ટ ફોલ્ડરમાં પોતાનું `.ccam-backup/<time>/` ફોલ્ડર મળે છે. ccam તેમને ક્યારેય કાઢી નાખતું નથી.
- **તમારી પોતાની લિંક.** જો કોઈ પાથ પહેલેથી તમે બનાવેલી લિંક હોય, ઉદાહરણ તરીકે dotfiles repository ની, અને ગ્રુપ ખાલી હોય, તો ccam લિંકનું લક્ષ્ય ગ્રુપમાં કૉપિ કરે છે અને લિંકને પોતે બૅકઅપમાં ખસેડે છે. લક્ષ્યને અડવામાં આવતું નથી.
- **ઇમ્પોર્ટ કરેલાં ફોલ્ડર.** `ccam import` વડે જ્યાં છે ત્યાં જ નોંધાયેલું એકાઉન્ટ (જેમ કે `~/.claude`) જ્યાં છે ત્યાં જ બદલાય છે, તેથી સાદું `claude` પણ ફેરફાર જુએ છે; યોજના એ જણાવે છે. જ્યારે તે કોઈ ગ્રુપ છોડે ત્યારે હંમેશાં કૉપિ રાખે છે, `--fresh` સાથે પણ, અને એકાઉન્ટને ભૂલી જતાં પહેલાં `ccam remove` તેને એ જ રીતે તેનાં ગ્રુપમાંથી બહાર કાઢે છે. તે કોઈ ગ્રુપમાં હોય ત્યાં સુધી તેની શેર કરેલી ફાઇલો ગ્રુપ ફોલ્ડરમાં રહે છે, તેથી `~/.local/share/ccam` કાઢી નાખતાં પહેલાં એકાઉન્ટ દૂર કરો (કે ગ્રુપ છોડો).
- **એક ફોલ્ડરમાં બે એકાઉન્ટ.** ઇતિહાસ શેર કરતાં એકાઉન્ટ એક જ ફાઇલોમાં લખે છે. જ્યારે બીજું એકાઉન્ટ એવા ફોલ્ડરમાં શરૂ થાય જ્યાં એક પહેલેથી ચાલે છે, ત્યારે ccam ચેતવણી આપે છે અને પૂછે છે.
- **હાથથી કરેલા ફેરફાર.** એકાઉન્ટ ન હોય તેવો સભ્ય અવગણાય છે અને `ccam doctor` તેની જાણ કરે છે. `local`, `on`, `off` અને `none` નામ ગ્રુપનાં નામ ન હોઈ શકે.

## સંબંધિત

- [ડેશબોર્ડ](/features/dashboard)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [શેર કરેલા MCP સર્વર](/features/mcp)
- [હાલનું સેટઅપ કૉપિ કરો](/features/migrate)
- [બધા કમાન્ડ](/docs/commands)
