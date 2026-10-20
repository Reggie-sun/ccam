---
title: તપાસ અને સમારકામ
description: ccam doctor વડે શેર લિંક, ફાઇલ પરવાનગીઓ, લૉગિન અને Claude Code નું વર્ઝન તપાસો, અને ccam doctor --fix વડે સુરક્ષિત રીતે સુધારી શકાય તે સુધારો.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` તમારા ccam ઇન્સ્ટૉલેશનને તપાસે છે: Claude Code binary, દરેક એકાઉન્ટનું ફોલ્ડર, લૉગિન અને શેર લિંક, શેર કરેલા MCP સર્વર, usage endpoint અને તમારા શૉર્ટકટ કમાન્ડ. `ccam doctor --fix` ડેટા ગુમાવ્યા વગર જે સુધારી શકાય તે સુધારે છે, અને જે ફાઇલ બદલે તેને પહેલાં બૅકઅપ ફોલ્ડરમાં ખસેડે છે. જ્યારે `ccam <name>` શેર લિંકની સમસ્યાને કારણે અટકી જાય, તમે ફાઇલો જાતે પાછી મૂકો, અથવા Claude Code અપડેટ થાય, ત્યારે તેને ચલાવો.

## ઉપયોગ

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

દરેક તારણ એક ચિહ્ન સાથેની એક લાઇન છે: `✓` બરાબર, `·` માહિતી, `!` ચેતવણી, `✗` ભૂલ. એકાઉન્ટનાં તારણ એકાઉન્ટના નામથી શરૂ થાય છે.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

`--fix` સાથે, સુધારેલું દરેક તારણ `→ fixed` થી પૂરું થાય છે:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--fix` | સુરક્ષિત રીતે સુધારી શકાય તે સુધારે છે (બદલાયેલી ફાઇલોનો બૅકઅપ લેવાય છે) | off |
| `--dry-run` (global) | દરેક તપાસ ચલાવે છે પણ `--fix` સાથે પણ કંઈ સુધારતું નથી | off |

`ccam doctor` કોઈ આર્ગ્યુમેન્ટ લેતું નથી. તે `--dry-run` ને માને છે: `ccam --dry-run doctor --fix` `ccam doctor` જેવાં જ તારણ છાપે છે અને કંઈ બદલતું નથી.

ઓછામાં ઓછું એક `✗` તારણ સુધાર્યા વગર બાકી રહે ત્યારે exit code `1` હોય છે, નહીં તો `0`. ચેતવણીઓ અને માહિતી ક્યારેય રનને નિષ્ફળ બનાવતી નથી. જ્યારે `--fix` વગર રન નિષ્ફળ જાય, ત્યારે ccam એક ખાલી લાઇન અને `` `ccam doctor --fix` repairs what it can. `` સાથે પૂરું થાય છે. લોડ ન થઈ શકે તેવું config પણ `1` અને `ccam:` ભૂલ સાથે પૂરું થાય છે.

## તે કેવી રીતે કામ કરે છે

`internal/cli/doctor.go` `~/.config/ccam/config.toml` લોડ કરે છે અને `internal/doctor/doctor.go` માં `doctor.Run` ને બોલાવે છે, જે તપાસો નિશ્ચિત ક્રમમાં ચલાવે છે અને તારણ પાછાં આપે છે. `--fix` `Options.Fix` તરીકે આગળ જાય છે, જે dry run માં false હોય છે.

<!-- motion:doctor-checks -->

વચ્ચેનો ભાગ દરેક એકાઉન્ટ માટે એક વાર ચાલે છે. તપાસો, અને દરેક માટે `--fix` શું કરે છે:

1. **Claude Code વર્ઝન** (`checkClaude`). launcher ની જેમ જ `claude` binary શોધે છે અને `claude --version` ચલાવે છે. જો વર્ઝન `2.1.302` થી શરૂ ન થાય, એટલે કે ccam ના ફોલ્ડર લેઆઉટની જે વર્ઝન સામે ચકાસણી થઈ છે (`doctor.LayoutVersion`), તો તે `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` ઉમેરે છે. binary ન મળે કે `--version` નિષ્ફળ જાય તો તે ભૂલ છે. કોઈ સુધારો નહીં.
2. **environment માં credentials** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` કે `ANTHROPIC_AUTH_TOKEN` સેટ હોય ત્યારે ચેતવણી આપે છે: એ shell માંથી શરૂ થતું દરેક એકાઉન્ટ તેને વારસામાં મેળવે છે, અને Claude Code એકાઉન્ટના પોતાના લૉગિન કરતાં તેને પ્રાધાન્ય આપે છે. કોઈ સુધારો નહીં; તેને જાતે unset કરો.
3. **ફોલ્ડર** (`checkProfile`). એકાઉન્ટ ફોલ્ડર ન હોય તો તે ભૂલ છે અને એ એકાઉન્ટની તપાસો ત્યાં પૂરી થાય છે. group કે other bit સેટ હોય તેવી પરવાનગીઓ ચેતવણી છે; `--fix` `chmod 0700` ચલાવે છે. Windows માં permission bits નથી, તેથી ત્યાં આ તપાસ અને નીચેની `0600` તપાસ છોડી દેવાય છે; તમારી યુઝર પ્રોફાઇલમાંનું ફોલ્ડર તેની access list વડે સુરક્ષિત છે.
4. **લૉગિન ફાઇલ** (`checkCredentials`). `.credentials.json` ને `lstat` વડે તપાસવામાં આવે છે. symlink ભૂલ છે, કારણ કે દરેક એકાઉન્ટ પાસે પોતાનો ટોકન હોવો જોઈએ, અને તે સુધારાતી નથી. `0600` સિવાયનો mode ચેતવણી છે; `--fix` `chmod 0600` ચલાવે છે. પછી લૉગિન વંચાય છે (macOS પર પહેલાં Keychain માંથી). લૉગિન ન હોય તો તે `ccam login <name>` તરફ નિર્દેશ કરતી ચેતવણી છે. મુદત પૂરી થયેલો refresh token `login expired` ભૂલ છે; મુદત પૂરી થયેલો access token ફક્ત માહિતી છે, કારણ કે Claude Code આગલી વાર શરૂ થાય ત્યારે તેને રિન્યૂ કરે છે. આમાંનું કંઈ સુધારાતું નથી.
5. **શેર લિંક** (`share.Check`). ફક્ત `lstat` અને `readlink` વાપરીને `share.Layout` ના દરેક પાથની એકાઉન્ટનાં શેર સેટિંગ્સ સાથે સરખામણી કરે છે. તે ચાર પ્રકાર બતાવે છે: ખૂટતી લિંક, બીજે ક્યાંક નિર્દેશ કરતી લિંક, જ્યાં લિંક હોવી જોઈએ ત્યાં સાચી ફાઇલ, અને હવે local થયેલી આઇટમ માટે બાકી રહેલી લિંક. ખૂટતી લિંક ચેતવણી છે, બાકીની ભૂલો છે.
6. **અજાણી એન્ટ્રીઓ** (`checkUnknownEntries`). એકાઉન્ટ ફોલ્ડરમાંનાં ટોચના સ્તરનાં જે નામ ccam જાણતું નથી (`share.Known`) તેમની યાદી આપે છે, Claude Code નું લેઆઉટ બદલાયું હોવાના સંકેત તરીકે. ફક્ત માહિતી.
7. **plugin પાથ** (`checkPlugins`). plugin registry (`installed_plugins.json`, `known_marketplaces.json`) absolute પાથ સાચવે છે. દરેક ગ્રુપના `plugins/` માં, કોઈ એક એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા પાથ એ એકાઉન્ટ દૂર થાય કે તેનું નામ બદલાય ત્યારે તૂટી જાય છે. `--fix` તેમને `profile.RewritePluginPaths` વડે ગ્રુપ ફોલ્ડર પર ફરીથી લખે છે, જે registry ફાઇલનું atomic write છે.
8. **MCP સર્વર** (`checkMCP`). દરેક ગ્રુપની `mcp.json` એક વાર લોડ કરે છે. ન વંચાય તેવી ફાઇલ ભૂલ છે; `mcp.Validate` માં નિષ્ફળ જતો સર્વર ચેતવણી છે. કોઈ સુધારો નહીં.
9. **વપરાશના જવાબનું ફૉર્મેટ** (`checkUsageFormat`). કૅશ વગર અને સામાન્ય 3 સેકન્ડની મર્યાદા સાથે, કોઈ એક જવાબ આપે ત્યાં સુધી એક પછી એક એકાઉન્ટ માટે usage endpoint ને પૂછે છે. જો જવાબ હવે parse ન થાય, તો તે ચેતવણી આપે છે કે `ccam usage` સ્ટેટસ લાઇન કૅશ પર આધાર રાખશે. નેટવર્ક વાપરતી આ એકમાત્ર તપાસ છે.
10. **શૉર્ટકટ** (`checkShims`). ફક્ત ત્યારે જ જ્યારે કોઈ એકાઉન્ટ પાસે શૉર્ટકટ હોય. `~/.local/bin` માં ખૂટતી લિંક ચેતવણી છે; `--fix` `shim.Install` ને બોલાવે છે, જે ફક્ત ખૂટતી લિંક બનાવે છે અને બીજી ફાઇલોને ક્યારેય અડતું નથી. `PATH` પર ન હોય તેવું bin ફોલ્ડર ચેતવણી છે અને સુધારાતું નથી.
11. **ફોલ્ડર** (`checkFolders`). ફક્ત ત્યારે જ્યારે [`ccam folder`](/features/folder) entries હોય. જે entry નું એકાઉન્ટ અસ્તિત્વમાં ન હોય તે ભૂલ છે, અને જે entry નું ફોલ્ડર ન રહ્યું હોય તે ચેતવણી છે. બેમાંથી એકેય સુધારાતું નથી, કારણ કે entry કાઢી નાખવાથી ત્યાં કયું એકાઉન્ટ શરૂ થાય તે બદલાઈ જાય.

### શેર લિંકનું સમારકામ

શરૂ કરવાનો પાથ ખૂટતી લિંક તરત બનાવે છે, પણ ડેટા ગુમાવી શકે તેવી કોઈ પણ બાબત માટે અટકી જાય છે અને અહીં નિર્દેશ કરે છે:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

શેર લિંકની કોઈ પણ સમસ્યા માટે, `--fix` એકાઉન્ટનાં શેર સેટિંગ્સને હાલની અને લક્ષ્ય બંને સ્થિતિ તરીકે આપીને `share.Build` ને બોલાવે છે, પછી `Plan.Apply`. `Build` ડિસ્ક પર જે છે તેના પરથી દરેક પગલું ફરી નક્કી કરે છે, એ જ યોજના જે `ccam share` વાપરે છે, તેથી સમારકામ શેરિંગ ચાલુ કરવાના જ નિયમો અનુસરે છે.

<!-- motion:doctor-repair -->

અહીં `settings.json` સાચી ફાઇલ છે. ગ્રુપ પાસે પહેલેથી જુદી `settings.json` હોવાથી, ટકરાવમાં તેની ડિફૉલ્ટ પસંદગી `KeepGroup` લેવાય છે, તેથી ગ્રુપની કૉપિ રહે છે. એકાઉન્ટની ફાઇલ એકાઉન્ટ ફોલ્ડરની અંદર `.ccam-backup/<time>/settings.json` માં ખસે છે, અને `os.Symlink` પાથને ગ્રુપ સાથે લિંક કરે છે. બીજા કિસ્સા:

- જો ગ્રુપ પાસે હજી કૉપિ ન હોય, તો એકાઉન્ટની ફાઇલ ગ્રુપમાં ખસે છે અને તેની પહેલી સામગ્રી બને છે.
- ફોલ્ડર માટે, ગ્રુપમાં ન હોય તેવી ફાઇલો ગ્રુપમાં ખસે છે; બાકીનું બૅકઅપમાં જાય છે.
- `history.jsonl` માટે, ડિફૉલ્ટ `merge by line` છે: ગ્રુપની લાઇનો, પછી તમારી નવી લાઇનો. ગ્રુપની અગાઉની કૉપિ `.ccam-backup/<time>/group/` માં સેવ થાય છે.
- બીજા ગ્રુપની લિંક બદલાઈ જાય છે. ccam ની બહારના ફોલ્ડરની લિંક બૅકઅપમાં જાય છે, અને ગ્રુપ ખાલી હોય તો તેની સામગ્રી ગ્રુપની શરૂઆતની સામગ્રી બને છે.
- હવે local થયેલી આઇટમ માટે બાકી રહેલી લિંક દૂર થાય છે, અને એકાઉન્ટને ગ્રુપની સામગ્રીની કૉપિ મળે છે.

લિંક દૂર કરવા સિવાય, દરેક પગલું ડેટાને ગ્રુપમાં કે બૅકઅપ ફોલ્ડરમાં ખસેડે છે. કંઈ પણ કાઢી નાખવામાં આવતું નથી.

## સલામતી અને ખાસ કિસ્સા

- સાદું `ccam doctor` ક્યારેય કંઈ બદલતું નથી. સ્ક્રિપ્ટમાં ખાતરી માટે `--fix` સાથે `--dry-run` વાપરો.
- `doctor --fix` એકાઉન્ટ ચાલી રહ્યું છે કે નહીં તે તપાસતું નથી. કોઈ એકાઉન્ટની શેર લિંક સુધારતાં પહેલાં તેનું Claude Code બંધ કરો, કારણ કે સેશન એ ફાઇલોમાં લખવાનું ચાલુ રાખે છે.
- જો શેરનું સમારકામ નિષ્ફળ જાય, તો ccam એ એકાઉન્ટ માટે `could not repair share links: <error>` બતાવે છે અને તેની અજાણી એન્ટ્રીઓની તપાસ છોડી દે છે. બીજાં એકાઉન્ટ તો પણ તપાસાય છે.
- જો `ccam share` અધવચ્ચે નિષ્ફળ જાય, તો તે પહેલાં લક્ષ્ય સેટિંગ્સ સેવ કરે છે અને `ccam doctor --fix` વડે પૂરું કરવા કહે છે.
- symlink થયેલી `.credentials.json`, મુદત પૂરી થયેલું લૉગિન અને environment માંનાં credentials બતાવવામાં આવે છે પણ ક્યારેય બદલાતાં નથી.
- plugin registry નું ફરીથી લખાણ atomic રીતે થાય છે પણ `.ccam-backup/` માં કૉપિ થતું નથી.

## સંબંધિત

- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [હાલનું સેટઅપ કૉપિ કરવું](/features/migrate)
- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [કમાન્ડ](/docs/commands)
