---
title: તપાસ અને સમારકામ
description: ccam doctor વડે શેર લિંક, ફાઇલ પરવાનગીઓ, લૉગિન અને Claude Code નું વર્ઝન તપાસો, અને ccam doctor --fix વડે સુરક્ષિત રીતે સુધારી શકાય તે સુધારો.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: 0f5ae4ab9a54
---

`ccam doctor` તમારા ccam ઇન્સ્ટૉલેશનને તપાસે છે: Claude Code binary, દરેક એકાઉન્ટનું ફોલ્ડર, લૉગિન અને શેર લિંક, શેર કરેલા MCP સર્વર, usage endpoint અને તમારા શૉર્ટકટ કમાન્ડ. `ccam doctor --fix` ડેટા ગુમાવ્યા વગર જે સુધારી શકાય તે સુધારે છે, અને જે ફાઇલ બદલે તેને પહેલાં બૅકઅપ ફોલ્ડરમાં ખસેડે છે. જ્યારે `ccam <name>` શેર લિંકની સમસ્યાને કારણે અટકી જાય, તમે ફાઇલો જાતે પાછી મૂકો, અથવા Claude Code અપડેટ થાય, ત્યારે તેને ચલાવો.

## ઉપયોગ

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

પહેલી લાઇન બધાં એકાઉન્ટને એકસાથે ચિહ્નિત કરે છે, જેથી બીજું કંઈ વાંચતાં પહેલાં તમે જોઈ શકો કે સમસ્યાઓ ક્યાં છે. તેની નીચે, દરેક એકાઉન્ટને એક બૉક્સ મળે છે, અને કોઈ એકાઉન્ટનાં ન હોય તેવાં તારણ (Claude Code, environment variables, plugin, MCP સર્વર, શૉર્ટકટ, ફોલ્ડર નિયમો, ગ્રુપ અને સેટિંગ્સ, `auto_switch`) `General` બૉક્સમાં જાય છે. બૉક્સનું શીર્ષક એકાઉન્ટનું સૌથી ખરાબ ચિહ્ન અને તેમાં કેટલી સમસ્યાઓ છે તે બતાવે છે. અંદરની દરેક હરોળ તપાસાયેલો એક વિસ્તાર છે, હંમેશાં એ ક્રમમાં `Folder`, `Login` અને `Share links`, પછી બાકીના, એક ચિહ્ન સાથે: `✔` બરાબર, `ℹ` માહિતી, `⚠` ચેતવણી, `✘` ભૂલ. જેમાં કંઈ ખોટું ન હોય તેવું ફોલ્ડર `✔ Folder  OK` બતાવે છે.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix` સાથે, સુધારેલું દરેક તારણ `→ fixed` થી પૂરું થાય છે, અને તે હવે સમસ્યા તરીકે ગણાતું નથી:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

ટર્મિનલ પર બૉક્સ તેની પહોળાઈમાં બંધબેસે છે અને લાંબા સંદેશા આગલી લાઇનમાં જાય છે, ટર્મિનલના પોતાના રંગોમાં (`NO_COLOR` તેમને બંધ કરે છે). બૉક્સ ક્યારેય 29 કૉલમથી સાંકડું હોતું નથી, તેથી સાંકડા ટર્મિનલ પર કોઈ શબ્દ કપાતો નથી. pipe કે ફાઇલમાં મોકલાય ત્યારે, આઉટપુટમાં રંગ હોતા નથી અને કોઈ સંદેશો તૂટતો કે કપાતો નથી, તેથી `ccam doctor | grep settings.json` આખી લાઇન શોધે છે; ત્યારે દરેક બૉક્સ સૌથી લાંબી લાઇન જેટલું પહોળું હોય છે. ડેશબોર્ડની Health ટૅબ એ જ બૉક્સ બતાવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--fix` | સુરક્ષિત રીતે સુધારી શકાય તે સુધારે છે (બદલાયેલી ફાઇલોનો બૅકઅપ લેવાય છે) | off |
| `--dry-run` (global) | દરેક તપાસ ચલાવે છે પણ `--fix` સાથે પણ કંઈ સુધારતું નથી | off |

`ccam doctor` કોઈ આર્ગ્યુમેન્ટ લેતું નથી. તે `--dry-run` ને માને છે: `ccam --dry-run doctor --fix` `ccam doctor` નાં તારણ છાપે છે, `--fix` જે બાજુ પર મૂકશે તે ઉમેરે છે, અને કંઈ બદલતું નથી:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

ઓછામાં ઓછું એક `✘` તારણ સુધાર્યા વગર બાકી રહે ત્યારે exit code `1` હોય છે, નહીં તો `0`. ચેતવણીઓ અને માહિતી ક્યારેય રનને નિષ્ફળ બનાવતી નથી. જ્યારે `--fix` વગર રન નિષ્ફળ જાય, ત્યારે ccam એક ખાલી લાઇન અને `` `ccam doctor --fix` repairs what it can. `` સાથે પૂરું થાય છે. લોડ ન થઈ શકે તેવું config પણ `1` અને `ccam:` ભૂલ સાથે પૂરું થાય છે.

## તે કેવી રીતે કામ કરે છે

`internal/cli/doctor.go` `~/.config/ccam/config.toml` લોડ કરે છે અને `internal/doctor/doctor.go` માં `doctor.Run` ને બોલાવે છે, જે તપાસો નિશ્ચિત ક્રમમાં ચલાવે છે અને તારણ પાછાં આપે છે, દરેક તે જે એકાઉન્ટનું છે અને જે વિસ્તાર તપાસ્યો તેની સાથે. `tui.HealthReport` તેમને બૉક્સમાં ગોઠવે છે. `--fix` `Options.Fix` તરીકે આગળ જાય છે, જે dry run માં false હોય છે.

<!-- motion:doctor-checks -->

વચ્ચેનો ભાગ દરેક એકાઉન્ટ માટે એક વાર ચાલે છે. તપાસો, અને દરેક માટે `--fix` શું કરે છે:

1. **Claude Code વર્ઝન** (`checkClaude`). launcher ની જેમ જ `claude` binary શોધે છે અને `claude --version` ચલાવે છે. binary ન મળે, `claude_bin` ચલાવી શકાય તેવી ફાઇલ ન હોય, કે `--version` નિષ્ફળ જાય તો તે ભૂલ છે. કોઈ સુધારો નહીં.
2. **environment માં credentials** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` કે `ANTHROPIC_AUTH_TOKEN` સેટ હોય ત્યારે ચેતવણી આપે છે: એ shell માંથી શરૂ થતું દરેક એકાઉન્ટ તેને વારસામાં મેળવે છે, અને Claude Code એકાઉન્ટના પોતાના લૉગિન કરતાં તેને પ્રાધાન્ય આપે છે. કોઈ સુધારો નહીં; તેને જાતે unset કરો.
3. **ફોલ્ડર** (`checkProfile`). એકાઉન્ટ ફોલ્ડર ન હોય તો તે ભૂલ છે અને એ એકાઉન્ટની તપાસો ત્યાં પૂરી થાય છે. group કે other bit સેટ હોય તેવી પરવાનગીઓ ચેતવણી છે; `--fix` `chmod 0700` ચલાવે છે. `~/.claude.json` ની બાજુમાં પોતાની `.claude.json` ધરાવતું ઇમ્પોર્ટ કરેલું `~/.claude` માહિતી છે: ત્યારે `ccam main` અને સાદું `claude` અલગ લૉગિન વિગતો અને MCP સર્વર રાખે છે. Windows માં permission bits નથી, તેથી ત્યાં આ તપાસ અને નીચેની `0600` તપાસ છોડી દેવાય છે; તમારી યુઝર પ્રોફાઇલમાંનું ફોલ્ડર તેની access list વડે સુરક્ષિત છે.
4. **લૉગિન ફાઇલ** (`checkCredentials`). `.credentials.json` ને `lstat` વડે તપાસવામાં આવે છે. symlink ભૂલ છે, કારણ કે દરેક એકાઉન્ટ પાસે પોતાનો ટોકન હોવો જોઈએ, અને તે સુધારાતી નથી. `0600` સિવાયનો mode ચેતવણી છે; `--fix` `chmod 0600` ચલાવે છે. પછી લૉગિન વંચાય છે (macOS પર પહેલાં Keychain માંથી). લૉગિન ન હોય, કે લૉગિન વગરની લૉગિન ફાઇલ હોય, તો તે `ccam login <name>` તરફ નિર્દેશ કરતી ચેતવણી છે. મુદત પૂરી થયેલો refresh token, કે સર્વરે ઇનકાર કર્યા પછી Claude Code એ સાફ કરેલું લૉગિન, `` login expired → `ccam login <name>` `` ભૂલ છે; મુદત પૂરી થયેલો access token ફક્ત માહિતી છે, કારણ કે Claude Code આગલી વાર શરૂ થાય ત્યારે તેને રિન્યૂ કરે છે. Windows Credential Manager માં રાખેલું લૉગિન માહિતી છે, કારણ કે ccam તેને વાંચી શકતું નથી. માન્ય લૉગિન, ફોલ્ડર અત્યારે જે ઇમેઇલથી લૉગિન છે તે અને પ્લાન બતાવે છે (`login valid (dev@company.com, Claude Team · Acme)`). આમાંનું કંઈ સુધારાતું નથી; `ccam login --expired` જરૂર હોય તે દરેક એકાઉન્ટમાં લૉગિન કરે છે.
5. **શેર લિંક** (`share.Check`). ફક્ત `lstat` અને `readlink` વાપરીને `share.Layout` ના દરેક પાથની એકાઉન્ટનાં ગ્રુપ સાથે સરખામણી કરે છે. તે પાંચ પ્રકાર બતાવે છે: ખૂટતી લિંક, બીજે ક્યાંક નિર્દેશ કરતી લિંક, જ્યાં લિંક હોવી જોઈએ ત્યાં સાચી ફાઇલ, હવે local થયેલી આઇટમ માટે બાકી રહેલી લિંક, અને જેનું ગ્રુપ ફોલ્ડર ખૂટતું હોય તેવી યોગ્ય લિંક (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; લૉન્ચ આ તપાસતું નથી). ખૂટતી લિંક અને ખૂટતું ગ્રુપ ફોલ્ડર ચેતવણીઓ છે, બાકીની ભૂલો છે. ગ્રુપમાંથી પાછી કૉપિ કરવાની જે પ્રક્રિયા અધવચ્ચે અટકી હોય તે પણ બતાવાય છે: `--fix` તૈયાર રાખેલી પૂરી કૉપિ તેની જગ્યાએ મૂકે છે, પણ જો ત્યાર પછી પાથ (લૉન્ચ દ્વારા) નવેસરથી બન્યો હોય, તો તે પૂરી કૉપિવાળા ફોલ્ડરનું નામ જણાવે છે, જેથી તમે જાતે ઉકેલી શકો. ccam ના પછીના વર્ઝનમાં કોઈ આઇટમમાં ઉમેરાયેલા પાથ (ઇતિહાસ માટે `file-history/`, `tasks/` અને `paste-cache/`, સેટિંગ્સ માટે `keybindings.json`) જ્યાં સુધી એકાઉન્ટ પાસે તેમની પોતાની કૉપિ હોય ત્યાં સુધી ફક્ત ચેતવણી છે: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. તેમના કારણે લૉન્ચ ક્યારેય અટકતું નથી.
6. **અજાણી એન્ટ્રીઓ** (`checkUnknownEntries`). એકાઉન્ટ ફોલ્ડરમાંનાં ટોચના સ્તરનાં જે નામ ccam જાણતું નથી (`share.Known`) તેમની યાદી આપે છે, Claude Code નું લેઆઉટ બદલાયું હોવાના સંકેત તરીકે. ફક્ત માહિતી.
7. **plugin પાથ** (`checkPlugins`). plugin registry (`installed_plugins.json`, `known_marketplaces.json`) absolute પાથ સાચવે છે. દરેક ગ્રુપના `plugins/` માં, કોઈ એક એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા પાથ એ એકાઉન્ટ દૂર થાય કે તેનું નામ બદલાય ત્યારે તૂટી જાય છે. `--fix` તેમને `profile.RewritePluginPaths` વડે ગ્રુપ ફોલ્ડર પર ફરીથી લખે છે, જે registry ફાઇલનું atomic write છે.
8. **MCP સર્વર** (`checkMCP`). દરેક ગ્રુપની `mcp.json` એક વાર લોડ કરે છે. ન વંચાય તેવી ફાઇલ ભૂલ છે; `mcp.Validate` માં નિષ્ફળ જતો સર્વર ચેતવણી છે. સાદા `claude` ની `~/.claude.json` માંના જે સર્વર કોઈ એકાઉન્ટ લોડ કરતું નથી તે ચેતવણી છે, જે તેમને કૉપિ કરતો `ccam migrate --mcp` કમાન્ડ જણાવે છે (ઘણાં ગ્રુપ MCP શેર કરતાં હોય ત્યારે `--group <a|b>` સાથે). ગ્રુપના ફક્ત એક એકાઉન્ટ પાસે હોય તેવા સર્વર માહિતી છે, તેમને શેર કરતા `ccam mcp import` કમાન્ડ સાથે, અને બીજાં એકાઉન્ટ જે MCP સર્વરમાં સાઇન ઇન છે તેમાં સાઇન ઇન ન હોય તેવા એકાઉન્ટને `/mcp` વડે સાઇન ઇન કરવાની ચેતવણી મળે છે (ccam ક્યારેય સાઇન-ઇન કૉપિ કરતું નથી). કોઈ સુધારો નહીં.
9. **વપરાશના જવાબનું ફૉર્મેટ** (`checkUsageFormat`). કૅશ વગર, કોઈ એક જવાબ આપે ત્યાં સુધી એક પછી એક એકાઉન્ટ માટે usage endpoint ને પૂછે છે, બધું એક જ 4 સેકન્ડની મર્યાદામાં. `usage_source = "claude"` સાથે તે તેના બદલે એક એકાઉન્ટ માટે Claude Code ને પૂછે છે. જો જવાબ હવે parse ન થાય, તો તે ચેતવણી આપે છે કે `ccam usage` સ્ટેટસ લાઇન કૅશ પર આધાર રાખશે; કોઈ એકાઉન્ટે જવાબ ન આપ્યો હોય ત્યારે, તે માહિતી તરીકે તેમ જણાવે છે (`could not check the usage response format: …`). નેટવર્ક વાપરતી આ એકમાત્ર તપાસ છે.
10. **શૉર્ટકટ** (`checkShims`). ફક્ત ત્યારે જ જ્યારે કોઈ એકાઉન્ટ પાસે શૉર્ટકટ હોય. `~/.local/bin` માં ખૂટતી લિંક ચેતવણી છે, અને ccam એ બનાવેલા જે શૉર્ટકટ હજી જૂનું ccam ચલાવે છે, જે બદલાઈ ગયું કે ખસેડાયું હોય, તે પણ (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` `shim.Install` ને બોલાવે છે, જે ફક્ત ccam એ બનાવેલી અને `~/.local/state/ccam/shortcuts.json` માં નોંધેલી લિંક બનાવે છે અને નવી કરે છે, અને બીજી ફાઇલોને ક્યારેય અડતું નથી. ચેતવણી `→ fixed` થી ત્યારે જ પૂરી થાય છે જ્યારે એ દરેક લિંક બની હોય; જો કોઈ ન બની શકી હોય, તો કારણ ઉમેરાય છે. હજી ઇન્સ્ટૉલ હોય તેવું બીજું ccam ચલાવતો શૉર્ટકટ કામ કરે છે, અને તે માહિતી છે; `--fix` તેને આ ccam તરફ વાળે છે. `PATH` પર ન હોય તેવું bin ફોલ્ડર ચેતવણી છે અને સુધારાતું નથી.
11. **ફોલ્ડર** (`checkFolders`). ફક્ત ત્યારે જ્યારે [`ccam folder`](/features/folder) entries હોય. જે entry નું એકાઉન્ટ અસ્તિત્વમાં ન હોય તે ભૂલ છે, અને જે entry નું ફોલ્ડર ન રહ્યું હોય તે ચેતવણી છે. બેમાંથી એકેય સુધારાતું નથી, કારણ કે entry કાઢી નાખવાથી ત્યાં કયું એકાઉન્ટ શરૂ થાય તે બદલાઈ જાય.
12. **ગ્રુપ અને સેટિંગ્સ** (`checkGroups`). જે ગ્રુપ સભ્ય એકાઉન્ટ ન હોય, જે એકાઉન્ટનું નામ ત્યાર પછી ccam કમાન્ડ બની ગયું હોય (તેને `ccam run <name>` વડે શરૂ કરો અથવા તેનું નામ બદલો), અજાણ્યું `theme` કે `usage_source`, અને `config.toml` ની જે લાઇનો ccam વાંચતું નથી તે ચેતવણીઓ છે. છેલ્લી વાળી સામાન્ય રીતે `[[group]]`, `[[profile]]` કે `[[folder]]` table ની નીચે લખેલું સેટિંગ હોય છે, જ્યાં તે એ table નું ગણાય છે (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); તે સુધારાય નહીં ત્યાં સુધી config બદલતા કમાન્ડ ઇનકાર કરે છે. કોઈ સુધારો નથી.
13. **લૉગિન** (`checkIdentities`). દરેક ફોલ્ડર કોના તરીકે લૉગિન છે તે તેની `.claude.json` માંથી વાંચે છે. `ccam login` એ નોંધ્યું હતું તેના કરતાં બીજા એકાઉન્ટમાં હવે લૉગિન થયેલું ફોલ્ડર (ઉદાહરણ તરીકે સેશનની અંદર `/login` પછી) ચેતવણી છે; લૉગિનની નોંધ વગરના એકાઉન્ટની, જેમ કે ઇમ્પોર્ટ કરેલું કે જૂના ccam એ બનાવેલું, સરખામણી તેના નોંધેલા ઇમેઇલથી થાય છે. બીજા ફોલ્ડર જેવા જ Claude એકાઉન્ટમાં લૉગિન થયેલું ફોલ્ડર પણ ચેતવણી છે: બંને તેમની વપરાશ મર્યાદા શેર કરે છે. બીજા ફોલ્ડર જેવો જ લૉગિન ટોકન ધરાવતું ફોલ્ડર, એટલે કે કૉપિ કરેલું ફોલ્ડર (સાદા `claude` નું `~/.claude` પણ), ભૂલ છે, કારણ કે જે પહેલું તેને રિન્યૂ કરે તે બીજાને લૉગઆઉટ કરી દે છે; તેને પોતાનું લૉગિન આપવા માટે `ccam login <name>` ચલાવો. કોઈ સુધારો નથી.
14. **`auto_switch`** (`checkAutoSwitch`). ફક્ત `auto_switch = true` હોય ત્યારે. `General` બૉક્સમાંની માહિતી કહે છે કે `auto_switch` શું કરે છે અને [તેના વર્ણન](/features/next#auto-switch) ની લિંક આપે છે. એકાઉન્ટમાં સીટવાળા દરેક પ્લાન (Claude Team, Claude Enterprise) માટે, એક લાઇન તેમનાં નામ આપે છે, કહે છે કે `auto_switch` એક સંસ્થાની સીટો વચ્ચે સેશન ખસેડે છે, સહાય કેન્દ્રની લિંક સાથે Anthropic નો સીટ અંગેનો નિયમ ટાંકે છે, અને પસંદગી તમારા પર છોડે છે. 2.1.214 કરતાં જૂનું Claude Code, જે hooks ને તેની process (`CLAUDE_PID`) જણાવતું નથી, ત્યાં ચેતવણી છે, કારણ કે કોઈ સેશન ખસેડી શકાતું નથી. ccam plugin વગરના, કે plugin બંધ કરેલા, એકાઉન્ટને તેના પોતાના બૉક્સમાં ચેતવણી મળે છે, Claude Code માં તે સુધારતા command સાથે: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (અથવા `/plugin enable ccam@ccam`). API key એકાઉન્ટ બાકાત રખાય છે. Windows પર એકમાત્ર લાઇન `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` છે. કોઈ સુધારો નથી.

એકાઉન્ટ પહેલાં, અધવચ્ચે અટકેલો **ગ્રુપનો ફેરફાર** `General` બૉક્સમાં ભૂલ છે: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` તેને જે ક્રમમાં યોજાયો હતો તે ક્રમમાં, તેના ટકરાવોને અપાયેલા જવાબો સાથે પૂરો કરે છે; ત્યાં સુધી ગ્રુપનો બીજો કોઈ ફેરફાર, નામ બદલવું કે દૂર કરવું શરૂ થતું નથી. તે પૂરો ન થાય ત્યાં સુધી તેમાંનાં એકાઉન્ટની ફક્ત ખૂટતી લિંક બને છે. જો તે પૂરો ન થઈ શકે, તો ભૂલ જણાવે છે કે એ ફાઇલ દૂર કરવાથી ફેરફારનો બાકીનો ભાગ જતો કરાય છે. `ccam --dry-run doctor --fix` ક્રમ (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) અને દરેક એકાઉન્ટનો ભાગ શું બાજુ પર મૂકશે તે બતાવે છે. બીજું ccam અત્યારે કરી રહ્યું હોય તેવો ફેરફાર માહિતી છે, અને તે ચાલે ત્યાં સુધી `--fix` કોઈ શેર લિંક સુધારતું નથી (`its share links are not repaired while another ccam changes the groups; …`). સમારકામ કરતી વખતે `--fix` `~/.local/state/ccam/group-change.lock` પકડી રાખે છે, અને ફેરફાર પૂરો કર્યા પછી `config.toml` ફરી વાંચે છે.

### શેર લિંકનું સમારકામ

શરૂ કરવાનો પાથ ખૂટતી લિંક તરત બનાવે છે, પણ ડેટા ગુમાવી શકે તેવી કોઈ પણ બાબત માટે અટકી જાય છે અને અહીં નિર્દેશ કરે છે:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

શેર લિંકની કોઈ પણ સમસ્યા માટે, `--fix` એકાઉન્ટનાં શેર સેટિંગ્સને હાલની અને લક્ષ્ય બંને સ્થિતિ તરીકે આપીને `share.Build` ને બોલાવે છે, પછી `Plan.Apply`. `Build` ડિસ્ક પર જે છે તેના પરથી દરેક પગલું ફરી નક્કી કરે છે, એ જ યોજના જે `ccam group` નો ફેરફાર વાપરે છે, તેથી સમારકામ ગ્રુપમાં જોડાવાના જ નિયમો અનુસરે છે.

<!-- motion:doctor-repair -->

અહીં `settings.json` સાચી ફાઇલ છે. ગ્રુપ પાસે પહેલેથી જુદી `settings.json` હોવાથી, ટકરાવમાં તેની ડિફૉલ્ટ પસંદગી `KeepGroup` લેવાય છે, તેથી ગ્રુપની કૉપિ રહે છે. એકાઉન્ટની ફાઇલ એકાઉન્ટ ફોલ્ડરની અંદર `.ccam-backup/<time>/settings.json` માં ખસે છે, અને `os.Symlink` પાથને ગ્રુપ સાથે લિંક કરે છે. `--fix` તેમ જણાવે છે: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. બીજા કિસ્સા:

- જો ગ્રુપ પાસે હજી કૉપિ ન હોય, તો એકાઉન્ટની ફાઇલ ગ્રુપમાં ખસે છે અને તેની પહેલી સામગ્રી બને છે.
- ફોલ્ડર માટે, ગ્રુપમાં ન હોય તેવી ફાઇલો ગ્રુપમાં ખસે છે. ગ્રુપ પાસે બાઇટે બાઇટ એ જ હોય તેવી ફાઇલો કાઢી નખાય છે, અને જુદી પડતી ફાઇલો બૅકઅપમાં જાય છે, જેમનાં નામ એક નોંધમાં આવે છે (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- ખૂટતું ગ્રુપ ફોલ્ડર પહેલાં બનાવાય છે.
- `history.jsonl` માટે, ડિફૉલ્ટ `merge by line` છે: ગ્રુપની લાઇનો, પછી તમારી નવી લાઇનો. ગ્રુપની અગાઉની કૉપિ `.ccam-backup/<time>/group/` માં સેવ થાય છે.
- બીજા ગ્રુપની લિંક બદલાઈ જાય છે. ccam ની બહારના ફોલ્ડરની લિંક બૅકઅપમાં જાય છે, અને ગ્રુપ ખાલી હોય તો તેની સામગ્રી ગ્રુપની શરૂઆતની સામગ્રી બને છે.
- હવે local થયેલી આઇટમ માટે બાકી રહેલી લિંક દૂર થાય છે, અને એકાઉન્ટને ગ્રુપની સામગ્રીની કૉપિ મળે છે.

લિંક દૂર કરવા અને ગ્રુપ પાસે બાઇટે બાઇટ પહેલેથી હોય તેવી ફાઇલો કાઢી નાખવા સિવાય, દરેક પગલું ડેટાને ગ્રુપમાં કે બૅકઅપ ફોલ્ડરમાં ખસેડે છે. કંઈ ગુમાવાતું નથી.

## સલામતી અને ખાસ કિસ્સા

- સાદું `ccam doctor` ક્યારેય કંઈ બદલતું નથી. સ્ક્રિપ્ટમાં ખાતરી માટે `--fix` સાથે `--dry-run` વાપરો.
- ચાલી રહેલા એકાઉન્ટ માટે, Claude Code ભલે ગમે તે રીતે શરૂ થયું હોય, `doctor --fix` ફક્ત ખૂટતી લિંક અને ખૂટતાં ગ્રુપ ફોલ્ડર બનાવે છે, જેમાં કંઈ ગુમાવાતું નથી. બીજું કંઈ બાકી રહે ત્યારે, તે `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` ઉમેરે છે, કારણ કે સેશન એ ફાઇલોમાં લખવાનું ચાલુ રાખે છે જે સમારકામ ખસેડશે.
- જો શેરનું સમારકામ નિષ્ફળ જાય, તો ccam એ એકાઉન્ટ માટે `could not repair share links: <error>` બતાવે છે અને તેની અજાણી એન્ટ્રીઓની તપાસ છોડી દે છે. બીજાં એકાઉન્ટ તો પણ તપાસાય છે.
- જો `ccam group` નો ફેરફાર અધવચ્ચે નિષ્ફળ જાય, તો નવાં ગ્રુપ પહેલેથી સેવ થયેલાં હોય છે અને ફેરફાર `~/.local/state/ccam/group-change.json` માં નોંધાય છે; ccam તમને `ccam doctor --fix` વડે પૂરું કરવા કહે છે, જે તેનો બાકીનો ભાગ, MCP પગલાં સહિત, પૂરો કરે છે.
- symlink થયેલી `.credentials.json`, મુદત પૂરી થયેલું લૉગિન અને environment માંનાં credentials બતાવવામાં આવે છે પણ ક્યારેય બદલાતાં નથી.
- plugin registry નું ફરીથી લખાણ atomic રીતે થાય છે પણ `.ccam-backup/` માં કૉપિ થતું નથી.

## સંબંધિત

- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [હાલનું સેટઅપ કૉપિ કરવું](/features/migrate)
- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [કમાન્ડ](/docs/commands)
