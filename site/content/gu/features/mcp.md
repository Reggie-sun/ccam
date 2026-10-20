---
title: શેર કરેલા MCP સર્વર
description: ccam mcp add, add-json, list, remove અને import વડે દરેક ગ્રુપ માટે MCP સર્વરની એક યાદી રાખો, જે --mcp-config વડે આપવામાં આવે છે કે દરેક એકાઉન્ટમાં sync થાય છે.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: 6d3b3c38c6cb
---

Claude Code user MCP સર્વર `.claude.json` માં રાખે છે, એ જ ફાઇલ જેમાં એકાઉન્ટનું લૉગિન હોય છે, તેથી આ ફાઇલ એકાઉન્ટ વચ્ચે શેર થઈ શકતી નથી. ccam શેર કરેલા સર્વર એક ગ્રુપ ફાઇલ, `mcp.json`, માં રાખે છે, અને શરૂઆત વખતે તે ગ્રુપના દરેક એકાઉન્ટને આપે છે. તમારાં બધાં એકાઉન્ટ માટે એક જ વાર સર્વર ઉમેરવા `ccam mcp` વાપરો, અને કોઈ એક સેશનમાં ઉમેરેલો સર્વર બીજાં એકાઉન્ટમાં ખસેડવા `ccam mcp import` વાપરો.

## ઉપયોગ

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

ફેરફાર દરેક એકાઉન્ટની આગલી શરૂઆતથી લાગુ પડે છે; ચાલી રહેલાં સેશન જે સર્વર સાથે શરૂ થયાં હતાં તે જ રાખે છે:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

ગ્રુપ પાસે પહેલેથી હોય તેવું નામ ઉમેરવાથી તેની વ્યાખ્યા બદલાય છે અને `added` ને બદલે `updated` છપાય છે. સર્વર કોઈ એકાઉન્ટ સુધી ત્યારે જ પહોંચે છે જ્યારે તે એ ગ્રુપમાં હોય અને ગ્રુપ MCP શેર કરતું હોય (`ccam group set <group> --add mcp`); ગ્રુપ MCP શેર ન કરતું હોય, ત્યારે `add` એ જણાવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` અને `remove` માટેનું લક્ષ્ય ગ્રુપ | જેના દ્વારા કોઈ એકાઉન્ટ MCP શેર કરે છે તે એકમાત્ર ગ્રુપ, નહીં તો `default` |
| `-- <command> [args…]` (`add`) | આ કમાન્ડને stdio સર્વર તરીકે ચલાવે છે | |
| `--url <url>` (`add`) | રિમોટ સર્વરનું URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` સર્વરનું transport | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, ઘણી વાર આપી શકાય | કોઈ નહીં |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio સર્વર માટે environment variable, ઘણી વાર આપી શકાય | કોઈ નહીં |
| `--force` (`import`) | એ જ નામના પણ જુદા ગ્રુપ સર્વરને બદલે છે; જૂની `mcp.json` નો બૅકઅપ લેવાય છે | off |
| `--keep-local` (`import`) | એકાઉન્ટની `.claude.json` માંની એન્ટ્રીઓ રાખે છે | off |
| `-h`, `--help` | દરેક કમાન્ડની મદદ | |
| `--dry-run` (global) | ફક્ત `mcp list` તેને માને છે; `add`, `add-json`, `remove` અને `import` તેનો ઇનકાર કરે છે અને કંઈ બદલતા નથી | off |

`add` ને `--` પછી કમાન્ડ કે `--url` માંથી એકની જરૂર છે; નહીં તો તે `a command to run (after --) or --url is required` સાથે નિષ્ફળ જાય છે. `--header` `--url` સર્વરને લાગુ પડે છે અને `--env` કમાન્ડ સર્વરને. `add-json` `claude mcp add-json` જેવું જ JSON લે છે, અને વ્યાખ્યા `command` કે `url` ધરાવતો ઑબ્જેક્ટ હોવી જોઈએ. `import` હંમેશાં એ ગ્રુપ વાપરે છે જેના દ્વારા એકાઉન્ટ MCP શેર કરે છે, અને `list` દરેક ગ્રુપ બતાવે છે, તેથી `-g` તેમને બદલતું નથી.

## તે કેવી રીતે કામ કરે છે

સંકળાયેલી ફાઇલો:

- `~/.local/share/ccam/groups/<group>/mcp.json`: શેર કરેલા સર્વર, `{"mcpServers": {…}}`, એ ફૉર્મેટ જે `--mcp-config` વાંચે છે. તે `fsutil.AtomicWrite` અને mode `0600` સાથે લખાય છે, કારણ કે સર્વરની એન્ટ્રીઓમાં ઘણી વાર API key હોય છે.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: એકાઉન્ટના સર્વર ધરાવતી Claude Code ની પોતાની ફાઇલ. ccam `profile.TopLevelValue` scanner વડે ફક્ત તેનું `mcpServers` મૂલ્ય વાંચે છે અને શરૂઆતના પાથ પર ક્યારેય આખી ફાઇલ parse કરતું નથી.
- `~/.local/state/ccam/mcp/<name>/`: `seen` ફાઇલ, `sync-stamp` અને `.claude.json` ના બૅકઅપ.

બે વ્યૂહરચના છે, જે `~/.config/ccam/config.toml` માં `mcp_strategy` વડે સેટ થાય છે. ડિફૉલ્ટ flag વ્યૂહરચના છે:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` એકાઉન્ટનું MCP ગ્રુપ શોધે છે. `mcp.FlagArgs` એક `stat` વડે તપાસે છે કે ગ્રુપની `mcp.json` છે, તેને parse કર્યા વગર.
2. `Plan.InsertArgs` `--mcp-config=<path>` ને `claude` પછી તરત, તમારા પોતાના આર્ગ્યુમેન્ટ્સ પહેલાં મૂકે છે, તેથી તમારું પોતાનું `--mcp-config` કે `--strict-mcp-config` તો પણ તેના પર લાગુ પડે છે.
3. `mcp.NewLocal` તપાસે છે કે છેલ્લી શરૂઆત પછી એકાઉન્ટની `.claude.json` બદલાઈ છે કે નહીં, અને ત્યારે જ તેના સર્વર parse કરે છે.
4. `syscall.Exec` ccam ને Claude Code થી બદલે છે.

પાથ બે નહીં, એક આર્ગ્યુમેન્ટ, `--mcp-config=<path>`, તરીકે આપવામાં આવે છે. Claude Code નું `--mcp-config` ઘણાં મૂલ્યો લે છે, તેથી `claude --mcp-config a.json "prompt"` prompt ને બીજી config ફાઇલ તરીકે વાંચે અને `MCP config file not found` સાથે નિષ્ફળ જાય. એક આર્ગ્યુમેન્ટવાળું સ્વરૂપ ત્યાં જ અટકે છે, અને પછી આવતો prompt કે `mcp list` જેવો subcommand કામ કરે છે. Claude Code `--mcp-config` ના સર્વર એકાઉન્ટના પોતાના સર્વર ઉપરાંત લોડ કરે છે; નામ ટકરાય ત્યારે `.claude.json` ની એન્ટ્રી જીતે છે. flag વ્યૂહરચના ફક્ત Remote Control માટે `.claude.json` માં લખે છે.

`claude remote-control` (`rc`) કોઈ `--mcp-config` લેતું નથી, અને તેનાં સેશન સર્વર ફક્ત `.claude.json` માંથી વાંચે છે. `ccam work rc` માટે, ccam flag છોડી દે છે અને ટર્મિનલ પર પૂછે છે કે ગ્રુપના સર્વર બૅકઅપ સાથે `.claude.json` માં કૉપિ કરવા કે નહીં, જેમ `sync` વ્યૂહરચના કરે છે. કૉપિ ન બને ત્યાં સુધી તે પૂછે છે, પછી ફક્ત ગ્રુપની `mcp.json` બદલાય ત્યાર પછી ([Remote Control](/features/launch#remote-control)). નામ ટકરાય ત્યારે આવી કૉપિ `--mcp-config` સામે જીતતી હોવાથી, `ccam mcp add`, `add-json`, `remove` અને `import` ગ્રુપનાં બીજાં એકાઉન્ટ પાસે રહેલી કૉપિ અપડેટ કરે છે:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

સેશનની અંદર `claude mcp add` વડે ઉમેરેલો સર્વર ફક્ત એ એકાઉન્ટની `.claude.json` માં જાય છે. ccam તે સમયે તેને જોઈ શકતું નથી, કારણ કે Claude Code એ તેની જગ્યા લઈ લીધી છે, તેથી તે આગલી શરૂઆતે તપાસે છે. `seen` ફાઇલમાં `.claude.json` નો modification time અને પહેલેથી જણાવેલાં સર્વરનાં નામ હોય છે. સમય જુદો હોય ત્યારે, ccam નામોની ગ્રુપ સાથે સરખામણી કરે છે અને દરેક નવા સર્વર માટે, input ની રાહ જોયા વગર, એક વાર સૂચના છાપે છે:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` સાથે, સર્વર `claude mcp list` અને `/mcp` માં સામાન્ય user સર્વર તરીકે દેખાય છે, અને ccam શરૂઆત પહેલાં તેમને દરેક એકાઉન્ટની `.claude.json` માં લખે છે:

<!-- motion:mcp-sync -->

1. `mcp.Sync` ગ્રુપની `mcp.json` ના modification time ની `sync-stamp` સાથે સરખામણી કરે છે. જો તે બદલાયો ન હોય, તો એકાઉન્ટની `.claude.json` ખોલવામાં આવતી નથી.
2. તે Claude Code નું પોતાનું lock લે છે, ફાઇલની બાજુનું `.claude.json.lock` ફોલ્ડર, જે `mkdir` વડે બને છે. તે 2 સેકન્ડ સુધી રાહ જુએ છે અને Claude Code ની lock library ની જેમ 10 સેકન્ડથી જૂના lock ને stale ગણે છે.
3. તે હાલની `.claude.json` ને state ફોલ્ડરના `backups/` માં સેવ કરે છે અને સૌથી નવા 5 રાખે છે.
4. તે `mcpServers` માં દરેક ગ્રુપ સર્વર સેટ કરે છે, પહેલાં પોતે લખેલા પણ હવે ગ્રુપમાં ન હોય તેવા સર્વર દૂર કરે છે, એકાઉન્ટના પોતાના સર્વર રાખે છે, અને ફાઇલને તેના મૂળ mode સાથે atomic રીતે લખે છે.
5. તે નવો stamp અને પોતે સંભાળે છે તે નામો નોંધે છે, અને Claude Code `--mcp-config` વગર શરૂ થાય છે.

એકાઉન્ટનું કોઈ સેશન ચાલતું હોય ત્યારે sync `ccam: account work is already running, skipping MCP sync` સાથે છોડી દેવાય છે, કારણ કે Claude Code ચાલતી વખતે `.claude.json` ફરીથી લખે છે. નિષ્ફળ sync `ccam: MCP sync failed (launching anyway): …` છાપે છે અને તો પણ શરૂ કરે છે. સેશનની અંદર ઉમેરેલા સર્વરની સૂચના ફક્ત flag વ્યૂહરચનામાં જ છે.

`ccam mcp import` એકાઉન્ટના પોતાના સર્વર તેના ગ્રુપમાં ખસેડે છે:

<!-- motion:mcp-import -->

ccam `mcp.LocalServers` વડે એકાઉન્ટના સર્વર વાંચે છે અને `mcp.SameDef` વડે દરેકની એ જ નામના ગ્રુપ સર્વર સાથે સરખામણી કરે છે, જે key ના ક્રમ અને spacing ને અવગણે છે. ગ્રુપમાં ન હોય કે એ જ વ્યાખ્યા સાથે હોય તેવો સર્વર કૉપિ થાય છે. જુદો સર્વર છોડી દેવાય છે, કારણ કે બીજું દરેક એકાઉન્ટ ગ્રુપની વ્યાખ્યા વાપરે છે:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` સાથે ગ્રુપનો સર્વર બદલાય છે, અને અગાઉની `mcp.json` પહેલાં ગ્રુપ ફોલ્ડરની અંદર `.ccam-backup/<time>/mcp.json` માં સેવ થાય છે. પછી `mcp.RemoveLocal` એ જ lock હેઠળ અને state ફોલ્ડરમાં બૅકઅપ લીધા પછી, કૉપિ થયેલી એન્ટ્રીઓ એકાઉન્ટની `.claude.json` માંથી દૂર કરે છે. આ જરૂરી છે, કારણ કે નામ ટકરાય ત્યારે એકાઉન્ટની એન્ટ્રી ગ્રુપની એન્ટ્રી પર જીતે. `--keep-local` તેમને રહેવા દે છે.

## સલામતી અને ખાસ કિસ્સા

- **સાઇન-ઇન.** OAuth થી સાઇન ઇન કરતા સર્વર તેમના ટોકન સર્વર યાદીમાં નહીં, પણ દરેક એકાઉન્ટની `.credentials.json` (`mcpOAuth`) માં રાખે છે, તેથી સર્વર શેર કરવાથી તેનું સાઇન-ઇન શેર થતું નથી. દરેક એકાઉન્ટમાં `/mcp` વડે એક વાર સાઇન ઇન કરો, અથવા [`ccam migrate`](/features/migrate) વડે `~/.claude` માંથી સાઇન-ઇન કૉપિ કરો. ગ્રુપને ક્યારેય લૉગિન ફાઇલ મળતી નથી.
- **ચાલતાં એકાઉન્ટ.** `import` ગ્રુપમાં કૉપિ કરે છે પણ ચાલતા એકાઉન્ટની પોતાની એન્ટ્રીઓ રહેવા દે છે અને તે જણાવે છે; સેશન પૂરું થયા પછી તેને ફરી ચલાવો. sync ચાલતાં એકાઉન્ટ છોડી દે છે.
- **Remote Control માટેની કૉપિ.** ચાલતા એકાઉન્ટની કૉપિને અડવામાં આવતું નથી: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` સાથે કશું અપડેટ કરવાની જરૂર નથી, કારણ કે દરેક શરૂઆત sync કરે છે.
- **તૂટેલી ફાઇલો.** જો `.claude.json` માન્ય JSON ન હોય, તો sync તેને અડ્યા વગર રહેવા દે છે અને ભૂલ બતાવે છે. જો 2 સેકન્ડમાં lock ન લઈ શકાય, તો ccam `could not lock … (Claude Code is writing the file)` બતાવે છે.
- **MCP શેર ન કરતાં એકાઉન્ટ.** `import` ઇનકાર કરે છે અને પહેલાં એકાઉન્ટને MCP શેર કરતા ગ્રુપમાં મૂકવા કહે છે (`ccam group join <group> side`).
- **ગ્રુપનાં નામ.** `-g` ફક્ત ગ્રુપનું નામ લે છે; તે `groups/` હેઠળ ફોલ્ડર બને છે, તેથી `../x` જેવાં નામનો ઇનકાર થાય છે.

## સંબંધિત

- [એકાઉન્ટ વચ્ચે શેરિંગ](/features/share)
- [ccam સાથે Remote Control વાપરવું](/guides/remote-control)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [હાલનું સેટઅપ કૉપિ કરો](/features/migrate)
- [બધા કમાન્ડ](/docs/commands)
