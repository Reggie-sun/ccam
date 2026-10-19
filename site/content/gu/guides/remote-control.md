---
title: ccam સાથે Remote Control વાપરવું
description: ccam <name> rc વડે કોઈ પણ એકાઉન્ટ તરીકે Claude Code Remote Control શરૂ કરો, કયા ઓપ્શન લાગુ પડે છે તે જુઓ, ccam શેર કરેલા MCP સર્વર કૉપિ કરવાનું કેમ પૂછે છે તે સમજો, અને સામાન્ય સમસ્યાઓ ઉકેલો.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

`ccam work rc` Claude Code નું Remote Control `work` એકાઉન્ટ તરીકે શરૂ કરે છે, જેથી તમે claude.ai/code કે Claude મોબાઇલ ઍપ પરથી આ કમ્પ્યૂટર પરનાં સેશન ચલાવી શકો. તે `claude rc` જેવું જ કામ કરે છે, ફક્ત એક ફરક સાથે: જો એકાઉન્ટ MCP સર્વર શેર કરતું હોય, તો ccam એક વાર પૂછે છે કે તેમને એકાઉન્ટની `.claude.json` માં કૉપિ કરવા કે નહીં, કારણ કે Remote Control સેશન સર્વર બીજે ક્યાંયથી વાંચતાં નથી.

## શરૂ કરો

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[શૉર્ટકટ](/features/shortcuts) પણ એ જ રીતે કામ કરે છે: `ccwork rc` એટલે `ccam work rc`. જે ફોલ્ડરમાં કામ કરવું હોય ત્યાં કમાન્ડ ચલાવો, પછી `work` જેવા જ Claude એકાઉન્ટથી સાઇન ઇન કરેલું claude.ai/code કે Claude ઍપ ખોલો.

## ઓપ્શન

અહીં બે પ્રકારના ઓપ્શન ભેગા થાય છે: ccam ના પોતાના, જે એકાઉન્ટના નામ પહેલાં આવે છે કે ccam નાં સેટિંગ્સમાં રહે છે, અને Remote Control ના, જે `rc` પછી આવે છે અને કોઈ ફેરફાર વગર Claude Code સુધી પહોંચે છે.

### ccam ના ઓપ્શન અને સેટિંગ્સ

| ઓપ્શન કે સેટિંગ | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--dry-run` (એકાઉન્ટ પહેલાં) | environment અને કમાન્ડ છાપે છે, અને MCP સૂચના હોય તો તે પણ. ન તો પૂછે છે, ન કૉપિ કરે છે | off |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | જ્યારે એકાઉન્ટ MCP સર્વર શેર કરતું હોય અને તેની `.claude.json` માં ગ્રુપની હાલની યાદી ન હોય ત્યારે પુછાય છે. `y` તેમને કૉપિ કરે છે, `N` કે Enter તેમના વગર શરૂ કરે છે | `N` |
| `~/.config/ccam/config.toml` માં `mcp_strategy = "sync"` | દરેક શરૂઆત પહેલાં, Remote Control હોય કે ન હોય, ગ્રુપના સર્વર `.claude.json` માં કૉપિ કરે છે, તેથી ક્યારેય કશું પુછાતું નથી | `"flag"` |
| `ccam share work --mcp off` | એકાઉન્ટ સાથે MCP સર્વર શેર કરવાનું બંધ કરે છે; પછી Remote Control ફક્ત એકાઉન્ટના પોતાના સર્વર સાથે શરૂ થાય છે અને કશું પુછાતું નથી | on |

### Remote Control ના ઓપ્શન

`rc` પછીનું બધું Claude Code નું છે. આ Claude Code 2.1.304 માં `claude remote-control` ના ઓપ્શન છે; તમારા વર્ઝનની યાદી `ccam work rc --help` છાપે છે.

| ઓપ્શન | શું કરે છે |
| --- | --- |
| `--name <name>` | claude.ai/code માં દેખાતું સેશનનું નામ |
| `--remote-control-session-name-prefix <prefix>` | આપમેળે બનતાં સેશન નામો માટેનો prefix (ડિફૉલ્ટ: host નું નામ) |
| `-c`, `--continue` | આ ફોલ્ડર માટે Remote Control એ છેલ્લે નોંધેલા સેશન સાથે, લગભગ 4 કલાકની અંદર, ફરી જોડાય છે |
| `--session-id <id>` | ચોક્કસ સેશન સાથે ફરી જોડાય છે |
| `--permission-mode <mode>` | નવાં સેશન માટે permission mode: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` કે `plan` |
| `--spawn <mode>` | `same-dir` (ડિફૉલ્ટ), `worktree` (દરેક સેશન માટે એક git worktree) કે `session` (એક સેશન, તે પૂરું થાય ત્યારે બહાર નીકળે છે) |
| `--capacity <N>` | `same-dir` અને `worktree` mode માં એક સાથે વધુમાં વધુ સેશન (ડિફૉલ્ટ 32) |
| `--[no-]create-session-in-dir` | શરૂઆતમાં હાલના ફોલ્ડરમાં પહેલું સેશન બનાવે છે (ડિફૉલ્ટ on) |
| `--[no-]chrome` | નવાં સેશન માટે Claude in Chrome (ડિફૉલ્ટ off) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | debug અને વિગતવાર output |

`rc` ને એકાઉન્ટના નામ પછી તરત મૂકો. ccam `-` થી શરૂ ન થતા પહેલા આર્ગ્યુમેન્ટને Claude Code નો કમાન્ડ ગણે છે, તેથી `ccam work rc --name laptop` માં તે `rc` જુએ છે, પણ `ccam work --model opus rc` માં તે `opus` જુએ છે.

## ccam શું કરે છે

શરૂઆત સામાન્ય રસ્તે જ થાય છે ([Claude Code શરૂ કરવું](/features/launch)), ફક્ત MCP પગલામાં એક ફેરફાર સાથે:

1. ccam એકાઉન્ટ શોધે છે, તેની શેર લિંક તપાસે છે અને `claude` શોધે છે, બીજી કોઈ પણ શરૂઆતની જેમ.
2. તે જુએ છે કે કમાન્ડ `rc` કે `remote-control` છે અને `--mcp-config` છોડી દે છે.
3. જો એકાઉન્ટ MCP સર્વર શેર કરતું હોય, તો તે ગ્રુપની `mcp.json` ની `~/.local/state/ccam/mcp/work/` માંના `sync-stamp` સાથે સરખામણી કરે છે. જો `.claude.json` માં પહેલેથી હાલના સર્વર હોય, અથવા ગ્રુપમાં કોઈ સર્વર જ ન હોય, તો તે સીધું પગલું 6 પર જાય છે.
4. નહીંતર સેશનને જે સર્વર નહીં મળે તેમનાં નામ તે જણાવે છે:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y` પર તે `.claude.json` નો બૅકઅપ સેવ કરે છે, Claude Code નું પોતાનું lock લે છે, ગ્રુપના સર્વર `mcpServers` માં ઉમેરે છે અને એકાઉન્ટના બીજા સર્વર જાળવી રાખે છે. `N` પર, ટર્મિનલ વગર, કે એકાઉન્ટ ચાલતું હોય ત્યારે, તે કશું બદલતું નથી અને એ જણાવે છે:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. તે સેશન નોંધે છે, `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર પર સેટ કરે છે અને `exec` વડે પોતાને `claude rc` થી બદલે છે.

પછી Remote Control સર્વર તરીકે ચાલતું રહે છે. તે તમારા માટે શરૂ કરે તે દરેક સેશન એ process નું child છે, તેથી તે `work` તરીકે ચાલે છે: `work` ના લૉગિન, સેટિંગ્સ અને વાતચીતના ઇતિહાસ સાથે, અને તે `work` ની વપરાશ મર્યાદામાં ગણાય છે. સર્વર ચાલે ત્યાં સુધી ccam એકાઉન્ટને ચાલતું ગણે છે.

dry run આ ફરક બતાવે છે. સામાન્ય શરૂઆત ગ્રુપના સર્વર આર્ગ્યુમેન્ટ તરીકે સાથે લઈ જાય છે, Remote Control નહીં:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## તે આ રીતે કેમ કામ કરે છે

**Remote Control કોઈ `--mcp-config` લેતું નથી.** સામાન્ય શરૂઆતમાં ccam ગ્રુપના સર્વર Claude Code ને `--mcp-config=<path>` તરીકે આપે છે અને Claude Code ની ફાઇલોમાં ક્યારેય લખતું નથી. `claude remote-control` તેના કમાન્ડ નામ પહેલાં એ ઓપ્શન નકારે છે અને પછી એક પણ સ્વીકારતું નથી, તેથી ccam 0.6.0 પહેલાં MCP સર્વર શેર કરતા દરેક એકાઉન્ટ માટે `ccam work rc` નિષ્ફળ જતું. હવે ccam એ ઓપ્શન છોડી દે છે.

**તેનાં સેશન સર્વર ફક્ત `.claude.json` માંથી વાંચે છે.** Remote Control જે સેશન શરૂ કરે તે કમાન્ડ લાઇનના ઓપ્શન વારસામાં મેળવતાં નથી. user MCP સર્વર તેમને એક જ જગ્યાએ મળે છે: એકાઉન્ટની `.claude.json`. એ ફાઇલમાં એકાઉન્ટની લૉગિન વિગતો અને ફોલ્ડર દીઠ state પણ હોય છે, તેથી ccam તેને `settings.json` ની જેમ શેર કરેલી નકલ સાથે લિંક કરી શકતું નથી. સેશનને શેર કરેલી યાદી આપવાનો એકમાત્ર રસ્તો સર્વરને તેમાં કૉપિ કરવાનો છે.

**ccam લખતાં પહેલાં પૂછે છે.** ડિફૉલ્ટ flag વ્યૂહરચનામાં ccam `.claude.json` ને ન અડવાનું વચન આપે છે. Remote Control એ એકમાત્ર કિસ્સો છે જ્યાં તેને ત્યાં લખવું પડે છે, તેથી તે પહેલાં પૂછે છે, બૅકઅપ રાખે છે, અને `mcpServers` સિવાય કશાને અડતું નથી.

**એકાઉન્ટ ચાલતું હોય ત્યારે તે લખતું નથી.** Claude Code સેશન દરમિયાન `.claude.json` ફરીથી લખે છે, અને બહારથી કરેલો ફેરફાર ખોવાઈ જાય કે તેના પર બીજું લખાઈ જાય. એટલે જ એકાઉન્ટનું કોઈ પણ સેશન ચાલતું હોય ત્યારે, બીજા Remote Control સર્વર સહિત, ccam કૉપિ છોડી દે છે.

**દરેક ફેરફાર દીઠ તે એક જ વાર પૂછે છે.** કૉપિ પછી `sync-stamp` ગ્રુપની `mcp.json` નો modification time રાખે છે. ccam ફરી ત્યારે જ પૂછે છે જ્યારે એ ફાઇલ બદલાય, ઉદાહરણ તરીકે `ccam mcp add` પછી. જો તમે `N` જવાબ આપો, તો હજી કોઈ કૉપિ નથી, તેથી તે દરેક Remote Control શરૂઆતે પૂછે છે.

**`ccam mcp` કૉપિને અદ્યતન રાખે છે.** નામ ટકરાય ત્યારે `.claude.json` ની એન્ટ્રી `--mcp-config` સામે જીતે છે, તેથી જૂની કૉપિ સામાન્ય શરૂઆતમાં પણ નવી વ્યાખ્યાને ઢાંકી દે. એટલે `ccam mcp add`, `add-json`, `remove` અને `import` ચાલતાં ન હોય તેવાં દરેક એકાઉન્ટની કૉપિ અપડેટ કરે છે:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**ટર્મિનલ વગર જવાબ "ના" ગણાય છે.** Remote Control ઘણી વાર સ્ક્રિપ્ટ, સર્વિસ કે `nohup` માંથી ચાલે છે, જ્યાં કોઈ જવાબ આપી શકતું નથી. ત્યારે ccam રાહ જોવાને બદલે સર્વર વગર શરૂ કરે છે.

## સમસ્યાઓ અને ઉકેલ

### Remote Control સેશનમાં શેર કરેલા MCP સર્વર નથી

પ્રશ્નનો જવાબ `N` આપવામાં આવ્યો હતો, ટર્મિનલ નહોતું, કે એકાઉન્ટ ચાલતું હતું. Remote Control બંધ કરો, એકાઉન્ટનાં બીજાં સેશન બંધ કરો, અને ટર્મિનલમાં `ccam work rc` ફરી ચલાવો; `y` જવાબ આપો. ક્યારેય ન પુછાય તે માટે `~/.config/ccam/config.toml` માં `mcp_strategy = "sync"` સેટ કરો. સેશનમાં પહેલેથી રહેલા સર્વર તે ચાલે ત્યાં સુધી બદલાતા નથી, તેથી કૉપિ પછી નવું સેશન શરૂ કરો.

### "Account work is running, so its .claude.json is left alone"

`work` નું એક સેશન હજી ખુલ્લું છે: સામાન્ય `ccam work`, શૉર્ટકટ, કે બીજું Remote Control સર્વર. તેમાંથી બહાર નીકળો અને `ccam work rc` ફરી ચલાવો. બીજું એકાઉન્ટ ચલાવતું હોય તેવાં સેશન ગણાતાં નથી.

### ccam દર વખતે પૂછે છે

તમે હજી સર્વર કૉપિ કર્યા નથી. એક વાર `y` જવાબ આપો, અથવા `mcp_strategy = "sync"` સેટ કરો. જો Remote Control માં તમારે સર્વર ક્યારેય ન જોઈતા હોય, તો એકાઉન્ટ સાથે MCP શેર કરવાનું બંધ કરો: `ccam share work --mcp off`.

### `ccam work rc` તરત જ `--mcp-config` વિશેની ભૂલ સાથે નિષ્ફળ જાય છે

તમારું ccam 0.6.0 કરતાં જૂનું છે અને Remote Control ને `--mcp-config` આપે છે. `ccam update` ચલાવો. હાલના ccam સાથે આ જ ભૂલનો અર્થ એ છે કે મૂલ્ય લેતો કોઈ ઓપ્શન `rc` પહેલાં આવ્યો હતો, જેમ કે `ccam work --model opus rc` માં, તેથી ccam એ `opus` ને કમાન્ડ માની લીધો; `rc` ને એકાઉન્ટના નામ પછી તરત મૂકો.

### `ccam mcp` વડે કરેલો ફેરફાર Remote Control સુધી પહોંચતો નથી

`ccam mcp add` અને `remove` ફક્ત ચાલતાં ન હોય તેવાં એકાઉન્ટમાં જ કૉપિ અપડેટ કરે છે, અને એ જણાવે છે: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control બંધ કરીને ફરી શરૂ કરો; ccam બદલાયેલી `mcp.json` જુએ છે અને ફરી પૂછે છે.

### કૉપિએ મારો પોતાનો સર્વર બદલી નાખ્યો

જો એકાઉન્ટ પાસે ગ્રુપ સર્વર જેવા જ નામનો પોતાનો સર્વર હતો, તો કૉપિએ તેને બદલી નાખ્યો. અગાઉની `.claude.json` `~/.local/state/ccam/mcp/work/backups/` માં છે, જે સૌથી નવા પાંચ બૅકઅપ રાખે છે. તમારું વર્ઝન દરેક એકાઉન્ટ માટે રાખવા, આગલી Remote Control શરૂઆત પહેલાં `ccam mcp import work <server> --force` વડે તેને ગ્રુપમાં ખસેડો.

### મારે કૉપિ થયેલા સર્વર દૂર કરવા છે

`ccam mcp remove <server>` સર્વરને ગ્રુપમાંથી અને દરેક કૉપિમાંથી દૂર કરે છે. કૉપિ પહેલાંની સ્થિતિમાં પાછા જવા, એકાઉન્ટનું કોઈ સેશન ચાલતું ન હોય ત્યારે `~/.local/state/ccam/mcp/work/backups/` માંથી સૌથી નવો બૅકઅપ પાછો મૂકો.

### Remote Control ફોલ્ડર પર વિશ્વાસ કરવો કે નહીં તે પૂછે છે, અથવા ટર્મિનલ વગર બહાર નીકળી જાય છે

ફોલ્ડર પરનો વિશ્વાસ દરેક એકાઉન્ટ માટે અલગથી `.claude.json` માં સંગ્રહાય છે, તેથી `personal` તરીકે તમે જે ફોલ્ડર પર વિશ્વાસ કર્યો હતો તે `work` માટે નવું છે. એ ફોલ્ડરમાં એક વાર `ccam work` ચલાવીને સ્વીકારો, પછી ત્યાં Remote Control શરૂ કરો.

### "You must be logged in with a Claude account that has a subscription"

Remote Control ને Claude subscription વાળું લૉગિન જોઈએ. `ccam list` માં એકાઉન્ટનો ઇમેઇલ તપાસો, અને જરૂર હોય તો `ccam login work` વડે ફરી સાઇન ઇન કરો.

### સેશન ફોન પર દેખાતું નથી

Claude ઍપ કે claude.ai/code બીજા Claude એકાઉન્ટથી સાઇન ઇન થયેલું છે. `work` માટે `ccam list` જે ઇમેઇલ બતાવે છે તેની સાથે તેની સરખામણી કરો.

### ccam ચેતવણી આપે છે કે આ ફોલ્ડરમાં બીજું Claude Code વાતચીતનો ઇતિહાસ શેર કરે છે

એ જ ઇતિહાસ ગ્રુપનું બીજું એકાઉન્ટ એ જ ફોલ્ડરમાં ચાલી રહ્યું છે. બંને એકાઉન્ટમાંથી એક જ વાતચીત આગળ ચલાવતાં સેશન ઇતિહાસને ગૂંચવી શકે છે. બંને જુદી જુદી વાતચીત પર કામ કરવાનાં હોય તો જ `y` જવાબ આપો, અથવા Remote Control બીજા ફોલ્ડરમાં શરૂ કરો.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [શેર કરેલા MCP સર્વર](/features/mcp)
- [એકાઉન્ટ વચ્ચે શેરિંગ](/features/share)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
