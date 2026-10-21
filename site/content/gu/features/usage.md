---
title: વપરાશ મર્યાદા
description: ccam usage દરેક એકાઉન્ટની 5 કલાકની અને સાપ્તાહિક મર્યાદા એકસાથે બતાવે છે, અને ccam statusline એકાઉન્ટનું નામ અને વપરાશ Claude Code માં મૂકે છે.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` બતાવે છે કે દરેક એકાઉન્ટે 5 કલાકની મર્યાદા, સાપ્તાહિક મર્યાદા અને દરેક મૉડલની પોતાની સાપ્તાહિક મર્યાદા (જેમ કે Fable) માંથી કેટલું વાપર્યું છે, દરેક મર્યાદા માટે એક બાર અને તે રીસેટ થવામાં બાકી સમય સાથે. કોઈ લાંબું કામ શરૂ કરતાં પહેલાં, અથવા Claude Code મર્યાદા નજીક હોવાનું કહે ત્યારે, કયા એકાઉન્ટ પાસે જગ્યા બાકી છે તે જોવા તેને વાપરો. `ccam statusline` એ Claude Code માટેનો સ્ટેટસ લાઇન કમાન્ડ છે: તે Claude Code ની અંદર એકાઉન્ટનું નામ અને વપરાશ બતાવે છે અને `ccam usage` માટે આંકડા સેવ કરે છે.

## ઉપયોગ

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

દરેક એકાઉન્ટને એક લાઇન મળે છે જેમાં તેનો નંબર, તેનું નામ અને તેના આંકડા કેટલા જૂના છે તે હોય છે; `(cached)` નો અર્થ છે કે તે endpoint માંથી નહીં, સેવ કરેલી ફાઇલમાંથી આવે છે. તેની નીચે, દરેક મર્યાદા માટે એક લાઇન 20 ખાનાંનો બાર, ટકાવારી અને વિન્ડો રીસેટ થવામાં બાકી સમય બતાવે છે: `5-hour`, `Weekly` (બધાં મૉડલ), પછી દરેક મૉડલની પોતાની સાપ્તાહિક વિન્ડો, જેમ કે `Fable`. જે વિન્ડો એકાઉન્ટ પાસે ન હોય તે `-` બતાવે છે. ટકાવારી 80% થી પીળી અને 100% પર લાલ થાય છે. ટર્મિનલ પર દરેક એકાઉન્ટ તરત `fetching…` તરીકે દેખાય છે, અને તે જવાબ આપે તેમ તેની લાઇનો દોરાય છે. છેલ્લી લાઇન એ એકાઉન્ટનું નામ આપે છે જે [`ccam next --auto`](/features/next) એ જ નિયમો પ્રમાણે પસંદ કરે (`auto_max_usage`, `auto_strategy`, ચાલી રહેલાં સેશન, `cross_org` સેટ ન હોય તો ફક્ત છેલ્લે વપરાયેલા એકાઉન્ટની સંસ્થાનાં એકાઉન્ટ, 100% પરની દરેક વિન્ડો છોડી દેવાય છે, તેથી જેનું અઠવાડિયું ભરાઈ ગયું હોય તે એકાઉન્ટનું નામ ક્યારેય આવતું નથી), અને તમારી પાસે એક કરતાં વધુ એકાઉન્ટ હોય ત્યારે જ દેખાય છે. લૉગિન ન થયેલું એકાઉન્ટ `not logged in → ccam login <name>` બતાવે છે, અને API key વાળું એકાઉન્ટ `no limits (API)` બતાવે છે.

`ccam statusline install` પછી, Claude Code `[work] · Opus · 5h 24% · 7d 41%` જેવી લાઇન બતાવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--json` | રિપોર્ટ JSON તરીકે છાપે છે: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | off |
| `--watch[=N]` | તમે Ctrl+C દબાવો ત્યાં સુધી દર `N` સેકન્ડે ફરીથી દોરે છે. એકલું `--watch`, અથવા 0 કે તેથી ઓછું મૂલ્ય, એટલે 60 | off |
| `-v`, `--verbose` | દરેક એકાઉન્ટ નીચે `source:` લાઇન ઉમેરે છે: તેના આંકડા ક્યાંથી આવ્યા (usage endpoint, Claude Code, સ્ટેટસ લાઇન કે સેવ કરેલી ફાઇલ) | off |
| `--no-cache` | 60 સેકન્ડનો કૅશ અવગણીને endpoint ને ફરીથી પૂછે છે | off |
| `--source live\|claude\|cache` | `live` એકાઉન્ટના token વડે endpoint ને પૂછે છે; `claude` Claude Code પાસે તે જોવડાવે છે (નીચે જુઓ); `cache` ફક્ત સેવ કરેલા આંકડા બતાવે છે | `usage_source`, નહીં તો `live` |
| `--refresh` | મુદત પૂરી થયેલા access token રિન્યૂ કરે છે, ફક્ત ચાલતાં ન હોય તેવાં એકાઉન્ટ માટે | off |
| `--dry-run` | સ્વીકારાય છે: `ccam usage` પોતાના કૅશ સિવાય કંઈ બદલતું નથી. `--refresh` સાથે ઇનકાર થાય છે: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

`--source` નું બીજું કોઈ પણ મૂલ્ય `--source must be live, claude or cache` સાથે નિષ્ફળ જાય છે.

### Claude Code પાસેથી જ વપરાશ

`config.toml` માં `usage_source = "claude"` (અથવા `--source claude`) સાથે, ccam લૉગિન token બિલકુલ વાંચતું નથી. દરેક એકાઉન્ટ માટે તે એકાઉન્ટના ફોલ્ડરમાં `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` શરૂ કરે છે અને standard input પર Claude Code ની `get_usage` request મોકલે છે (Claude Code 2.1.270 કે પછીનું); Claude Code પોતાના લૉગિન વડે વપરાશ જુએ છે અને 5 કલાકની, સાપ્તાહિક અને મૉડલ વિન્ડો સાથે જવાબ આપે છે. `--safe-mode` એકાઉન્ટના hooks, MCP સર્વર અને plugins ને શરૂ થતાં રોકે છે, અને કોઈ સેશન સેવ થતું નથી. દરેક એકાઉન્ટ માટે એક Claude Code શરૂઆત, લગભગ એક સેકન્ડ, લાગે છે, તેથી આ source `live` કરતાં ધીમો છે; એકાઉન્ટને એકસાથે પૂછવામાં આવે છે, દરેકને 20 સેકન્ડની અંદર, અને 60 સેકન્ડનો cache હંમેશની જેમ લાગુ પડે છે. Claude Code આ request ને પ્રાયોગિક ગણાવે છે, તેથી જો તે જવાબ આપવાનું બંધ કરે, તો ccam સેવ કરેલા આંકડા બતાવે છે અને `ccam doctor` તેની જાણ કરે છે. `--dry-run` Claude Code શરૂ કરતું નથી (તે એકાઉન્ટના ફોલ્ડરમાં લખે છે) અને સેવ કરેલા આંકડા બતાવે છે. `usage_source` માં `live` અથવા `claude` આપી શકાય છે; `cache` ફક્ત `--source` નો વિકલ્પ છે. `ccam usage` કોઈ આર્ગ્યુમેન્ટ લેતું નથી.

`ccam statusline` standard input માંથી Claude Code નું સ્ટેટસ લાઇન JSON વાંચે છે. તેનો એકમાત્ર ઓપ્શન `--wrap <command>` છે: હાલનો સ્ટેટસ લાઇન કમાન્ડ, જેની આઉટપુટની પહેલી લાઇનની આગળ એકાઉન્ટનું નામ મુકાય છે. `ccam statusline install` કોઈ ઓપ્શન લેતું નથી. બંને `` --dry-run is not supported by `ccam statusline`; nothing was changed `` સાથે `--dry-run` નો ઇનકાર કરે છે.

## તે કેવી રીતે કામ કરે છે

આંકડા બે જગ્યાએથી આવે છે: Anthropic નું દસ્તાવેજીકૃત ન થયેલું usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, જેને દરેક એકાઉન્ટના પોતાના લૉગિન સાથે બોલાવવામાં આવે છે, અને Claude Code પોતાના સ્ટેટસ લાઇન કમાન્ડને જે મૂલ્યો આપે છે તે. આ બધું `internal/usage` માં છે.

1. `cli.usageCmd` `config.toml` લોડ કરે છે અને દરેક એકાઉન્ટ માટે એક `usage.Account` બનાવે છે: તેનું ફોલ્ડર `~/.local/share/ccam/profiles/<name>/`, તેનું સેશન ફોલ્ડર (તે ચાલી રહ્યું છે કે નહીં તે જોવા) અને તે જ્યાં છે ત્યાં જ ઇમ્પોર્ટ થયું હતું કે નહીં.
2. `Fetcher.FetchAll` દરેક એકાઉન્ટ માટે એક goroutine શરૂ કરે છે. તે બધાં એક જ `http.Client` વાપરે છે, તેથી endpoint સાથેનાં કનેક્શન ફરી વપરાય છે. કનેક્ટ થવામાં 2 s લાગી શકે છે, અને દરેક એકાઉન્ટને કુલ 3 s મળે છે. client redirect અનુસરતો નથી, કારણ કે redirect ટોકનને બીજે ક્યાંક ફરીથી મોકલી દે, અને જવાબનો વધુમાં વધુ 1 MB વાંચે છે.
3. દરેક એકાઉન્ટ પૂરું થાય તેમ તેની હરોળ ફરીથી દોરાય છે. જે એકાઉન્ટ નિષ્ફળ જાય કે સમયસર જવાબ ન આપે તે સેવ કરેલા આંકડા પર આધાર રાખે છે, તેથી એક ધીમું એકાઉન્ટ ક્યારેય બીજાંને રોકતું નથી.

દરેક એકાઉન્ટ endpoint સામે પોતાની goroutine માં ચાલે છે; જેનો સમય પૂરો થાય તે તેના બદલે પોતાના સેવ કરેલા આંકડા બતાવે છે.

<!-- motion:usage-fetch -->

`personal` અને `work` સમયસર જવાબ આપે છે અને live આંકડા બતાવે છે. `side` ને 3 s કરતાં વધુ લાગે છે, તેથી `fallback()` તેની બે કૅશ ફાઇલો વાંચે છે અને જે છેલ્લે લખાઈ હોય તે વાપરે છે. `usage.Table()` દરેક એકાઉન્ટ માટે એક બ્લૉક દોરે છે અને `usage.Best()` છેલ્લી લાઇનમાં જેનું નામ આવે છે તે એકાઉન્ટ પસંદ કરે છે.

એક એકાઉન્ટ માટે, `Fetcher.Fetch` આ તપાસો ક્રમમાં કરે છે:

<!-- motion:usage-cache -->

1. `--source cache` સાથે, તે સીધું `fallback()` પર જાય છે.
2. `--no-cache` સેટ ન હોય ત્યાં સુધી, `~/.local/state/ccam/usage/<name>.json` માંનું 60 સેકન્ડ (`CacheTTL`) કરતાં ઓછું જૂનું live પરિણામ જેમ છે તેમ પાછું અપાય છે. `ccam usage` સળંગ બે વાર ચલાવવાથી એકાઉન્ટ દીઠ એક જ request થાય છે.
3. `readCredentials()` `.credentials.json` (કે macOS Keychain) ને બદલ્યા વગર વાંચે છે, અને ફક્ત `claudeAiOauth` key જ decode થાય છે. subscription લૉગિન ન હોય તો `not logged in`, અથવા `.claude.json` માં API key હોય ત્યારે `no limits (API)`.
4. મુદત પૂરી થયેલો access token વપરાતો નથી. `--refresh` વગર, એકાઉન્ટ પોતાના સેવ કરેલા આંકડા મુદત પૂરી થયેલા તરીકે ચિહ્નિત કરીને બતાવે છે.
5. નહીં તો ccam endpoint ને બોલાવે છે. 401 કે 403 મુદત પૂરી થયેલું લૉગિન ગણાય છે, બીજું કોઈ પણ status કે ન વાંચી શકાય તેવો body ભૂલ ગણાય છે, અને બંને સેવ કરેલા આંકડા પર આધાર રાખે છે. 429 નો અર્થ છે કે endpoint ને બહુ વાર પૂછવામાં આવી રહ્યું છે: ccam ફરી પૂછવાના સમય (`Retry-After`, અથવા 5 મિનિટ, વધુમાં વધુ એક કલાક) સાથે `<name>.pause` લખે છે, `rate limited by the usage endpoint; asking again after 15:10` સાથે સેવ કરેલા આંકડા બતાવે છે, અને ત્યાં સુધી, `--no-cache` સાથે પણ, તે એકાઉન્ટ માટે endpoint ને બોલાવતું નથી. તરત ફરી પ્રયાસ કરવાથી મર્યાદા ફક્ત લંબાય.
6. સારો જવાબ માપ્યાના સમય સાથે `<name>.json` માં લખાય છે.

`fallback()` `<name>.json` અને `<name>.statusline.json` બંને વાંચે છે અને વધુ નવું રાખે છે. સેવ કરેલા આંકડા બતાવતાં પહેલાં, તે દરેક વિન્ડોનો રીસેટ સમય તપાસે છે: આંકડા સેવ થયા પછી રીસેટ થયેલી વિન્ડો 0% બતાવે છે, કારણ કે જૂની ટકાવારી હવે લાગુ પડતી નથી. તેથી જ ઉપરનું `personal` રીસેટ સમય વગર `0%` બતાવે છે.

endpoint ના બે જવાબ ફૉર્મેટ છે. નવામાં `session` (5 કલાક), `weekly_all` અને `weekly_scoped` (મૉડલ પ્રમાણે) એન્ટ્રીઓ સાથેનો `limits` array છે. જૂનામાં `five_hour`, `seven_day`, `seven_day_opus` અને `seven_day_sonnet` objects છે. `ParseLive` બંને વાંચે છે: `limits` array જીતે છે, અને તે જે આવરી ન લે તે જૂના objects પૂરું કરે છે. કોઈ મર્યાદા વગરનો body બદલાયેલા ફૉર્મેટ તરીકે નોંધાય છે, જે વિશે `ccam doctor` તમને જણાવે છે.

### ટોકન રિન્યૂ કરવા

refresh token ફરતા રહે છે: એક રિન્યૂ કરવાથી તમને નવો refresh token મળે છે અને જૂનો અમાન્ય થાય છે. જો Claude Code એ જ એકાઉન્ટ તરીકે ચાલતું હોય ત્યારે ccam લૉગિન રિન્યૂ કરે, તો એ સેશન સાઇન આઉટ થઈ જાય. તેથી `ccam usage` ટોકન વાંચે છે અને ક્યારેય લખતું નથી, સિવાય કે તમે `--refresh` આપો. ત્યારે `Fetcher.refresh`:

1. `ccam import` વડે જ્યાં છે ત્યાં જ રજિસ્ટર થયેલાં એકાઉન્ટનો ઇનકાર કરે છે, કારણ કે સાદું `claude` ccam ની જાણ વગર એ લૉગિન વાપરતું હોઈ શકે;
2. macOS Keychain માં રાખેલાં લૉગિનનો ઇનકાર કરે છે;
3. lock ફાઇલ, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (એ variable સેટ ન હોય ત્યારે `~/.local/state/ccam/locks/` હેઠળ), લે છે, જેથી બે ccam process ક્યારેય એક જ ટોકન બે વાર રિન્યૂ ન કરે;
4. lock હેઠળ, એકાઉન્ટનું કોઈ સેશન ચાલતું હોય તો ઇનકાર કરે છે, અને બીજા process એ ટોકન પહેલેથી રિન્યૂ કર્યો હોય તો request છોડી દે છે;
5. refresh token ને token endpoint પર post કરે છે અને નવા ટોકન `fsutil.AtomicWrite` વડે `.credentials.json` માં પાછા લખે છે, ફાઇલની બીજી દરેક key રાખીને.

### સ્ટેટસ લાઇન collector

Claude Code સ્ટેટસ લાઇન કમાન્ડ ચલાવે છે અને તેને pipe દ્વારા JSON document આપે છે. ccam તેમાંથી મૉડલનું નામ અને `rate_limits` વાંચે છે:

<!-- motion:usage-statusline -->

1. `runStatusline` `CCAM_PROFILE` પરથી એકાઉન્ટ શોધે છે, જે ccam Claude Code શરૂ કરતી વખતે સેટ કરે છે, અથવા નહીં તો `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર સાથે મેળવીને.
2. `ParseStatusLine` `rate_limits.five_hour` અને `rate_limits.seven_day` વાંચે છે. API key વાળાં એકાઉન્ટ માટે, અથવા સેશનના પહેલા જવાબ પહેલાં, JSON માં rate limit હોતી નથી; ત્યારે કંઈ સેવ થતું નથી.
3. `SaveStatusLine` `~/.local/state/ccam/usage/<name>.statusline.json` લખે છે. endpoint જવાબ ન આપે ત્યારે `ccam usage` તેને વાપરે છે.
4. ccam `[name] · model · 5h N% · 7d N%` છાપે છે. `--wrap` સાથે, તે એ જ input સાથે તમારો કમાન્ડ `sh -c` દ્વારા ચલાવે છે અને તેના બદલે `[name]` અને પછી કમાન્ડની પહેલી લાઇન છાપે છે. Windows પર તે કમાન્ડ Claude Code ની જેમ ચલાવે છે: Git Bash ઇન્સ્ટૉલ હોય તો તેનાથી, નહીં તો PowerShell થી.
5. જ્યારે 5 કલાકનો કે સાપ્તાહિક વપરાશ 90% કે વધુ હોય અને બીજું એકાઉન્ટ હોય, ત્યારે લાઇનના અંતે `→ ccam next --auto --continue` આવે છે, ઉદાહરણ તરીકે `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

wrap કરેલા કમાન્ડને 2 સેકન્ડ મળે છે. ત્યાર પછી તેને બંધ કરવામાં આવે છે અને લાઇન ફક્ત એકાઉન્ટનું નામ બતાવે છે, તેથી ધીમો કમાન્ડ Claude Code ની સ્ટેટસ લાઇનને અટકાવી શકતો નથી.

`ccam statusline install` દરેક એકાઉન્ટની `settings.json` માં `statusLine` લખે છે. જ્યારે કોઈ એકાઉન્ટ પાસે પહેલેથી સ્ટેટસ લાઇન કમાન્ડ હોય, ત્યારે ccam તેને `--wrap '<old command>'` તરીકે રાખે છે. Windows પર ccam નો પાથ forward slashes અને હોમ ફોલ્ડર માટે `~` સાથે લખાય છે (`~/.local/bin/ccam.exe statusline`), જે Git Bash અને PowerShell બંને સમજે છે. ગ્રુપ દ્વારા શેર થયેલી `settings.json` તેની લિંક દ્વારા એક જ વાર લખાય છે, તેથી ગ્રુપના દરેક એકાઉન્ટને તે મળે છે. તેને ફરી ચલાવવાથી પહેલેથી ccam વાપરતી ફાઇલો માટે `(already set)` છપાય છે.

## સલામતી અને ખાસ કિસ્સા

- તમે `--refresh` ન આપો ત્યાં સુધી `ccam usage` ક્યારેય ટોકન ફાઇલ બદલતું નથી. જ્યારે લૉગિનની મુદત પૂરી થઈ જાય, ત્યારે એકાઉન્ટને `ccam <name>` વડે એક વાર ચલાવો, અને Claude Code તેને રિન્યૂ કરે છે.
- `--refresh` ક્યારેય ચાલતાં એકાઉન્ટ, ઇમ્પોર્ટ કરેલાં એકાઉન્ટ કે Keychain લૉગિનને અડતું નથી. ઇનકાર થયેલું refresh સેવ કરેલા આંકડા અને કારણ બતાવે છે.
- usage endpoint જાહેર રીતે દસ્તાવેજીકૃત નથી અને બદલાઈ શકે છે. જ્યારે તે બદલાય, ત્યારે પણ તમે સેવ કરેલા આંકડા જુઓ છો, અને `ccam doctor` ફૉર્મેટના ફેરફાર વિશે જણાવે છે.
- સેવ કરેલા આંકડા અને જવાબ બંને ન હોય ત્યારે, એકાઉન્ટ લાઇન `lookup failed:` અને કારણ બતાવે છે, ઉદાહરણ તરીકે `--source cache` સાથે `no saved usage`.
- કૅશ ફાઇલોમાં ફક્ત ટકાવારી અને રીસેટ સમય હોય છે, ક્યારેય ટોકન નહીં.

## સંબંધિત

- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next)
- [વપરાશ મર્યાદા સાથે કામ કરવાની રીત](/guides/usage-limit-workflow)
- [કમાન્ડ](/docs/commands)
