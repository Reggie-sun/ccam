---
title: આગલા એકાઉન્ટ પર સ્વિચ કરવું
description: ccam next બીજા એકાઉન્ટ તરીકે Claude Code ખોલે છે; --auto સૌથી વધુ જગ્યાવાળું લે છે, --continue વાતચીત આગળ લઈ જાય છે, --fallback -p jobs ફરી ચલાવે છે.
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` તમે નામ આપ્યા વગર Claude Code ને બીજા એકાઉન્ટ તરીકે ખોલે છે. જ્યારે કોઈ એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે તેને વાપરો: સાદું `ccam next` છેલ્લે વાપરેલા એકાઉન્ટ પછીનું એકાઉન્ટ લે છે, `--auto` સૌથી વધુ જગ્યાવાળું એકાઉન્ટ લે છે, અને `--continue` આ ફોલ્ડરની છેલ્લી વાતચીત નવા એકાઉન્ટમાં લઈ જાય છે, બંને ઇતિહાસ શેર કરતાં હોય કે નહીં. Claude Code ની અંદર, [ccam plugin](#inside-claude-code) કોઈ મર્યાદા સેશનને અટકાવે ત્યારે તમને જણાવે છે અને `/ccam:switch` વાતચીત આગળ સોંપે છે. `claude -p` ચલાવતી scripts માટે, `ccam --fallback` સ્વિચ જાતે કરે છે (જુઓ [મર્યાદાએ ફરી ચલાવવું](#run-again-at-the-limit)).

## ઉપયોગ

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam સ્વિચને standard error પર `→ <from> → <to>` તરીકે છાપે છે, `--auto` ની પસંદગી હોય તો તેના કારણ સાથે, પછી `ccam <name>` ની જેમ જ Claude Code ને નવા એકાઉન્ટ તરીકે શરૂ કરે છે. `<from>` એ તમે છેલ્લે વાપરેલું એકાઉન્ટ છે, અથવા `--continue` સાથે, જે એકાઉન્ટની વાતચીત છે તે. `--` પછીના આર્ગ્યુમેન્ટ્સ કોઈ ફેરફાર વગર `claude` ને જાય છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--auto` | વપરાશ પ્રમાણે લૉગિન થયેલું એકાઉન્ટ પસંદ કરે છે (નીચે જુઓ), 100% પરનાં એકાઉન્ટ છોડીને | off |
| `--max-usage <1-100>` | `--auto` સાથે: આ વપરાશથી નીચેનાં એકાઉન્ટને પ્રાથમિકતા આપે છે | `auto_max_usage`, નહીં તો 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` સાથે: સૌથી ઓછો 5 કલાકનો વપરાશ, સૌથી વહેલો 5 કલાકનો રીસેટ, અથવા જે સાપ્તાહિક quota સૌથી પહેલાં વણવપરાયેલો રહી જાય | `auto_strategy`, નહીં તો `headroom` |
| `--cross-org` | `--auto` સાથે: બીજી સંસ્થાનાં એકાઉન્ટ પણ પસંદ કરે છે (નીચે જુઓ) | `cross_org`, નહીં તો off |
| `--wait[=<duration>]` | `--auto` સાથે: કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે, પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, વધુમાં વધુ આટલો સમય | off; એકલું આપો તો `5h` |
| `-c`, `--continue` | આ ફોલ્ડરની છેલ્લી વાતચીત નવા એકાઉન્ટ તરીકે આગળ લઈ જાય છે (નીચે જુઓ) | off |
| `-- <claude args>` | ccam ના પોતાના આર્ગ્યુમેન્ટ્સ પછી `claude` ને આપવામાં આવતા આર્ગ્યુમેન્ટ્સ | કોઈ નહીં |
| `--dry-run` | માનવામાં આવે છે: Claude Code શરૂ કર્યા વગર કે શરૂઆત નોંધ્યા વગર, `→ <from> → <to>` અને જે કમાન્ડ લાઇન ચલાવશે તે છાપે છે | off |

`--auto` અને `--continue` સાથે કામ કરે છે. `--continue` સાથે, `--resume` આર્ગ્યુમેન્ટ્સ તમારા પોતાના આર્ગ્યુમેન્ટ્સ પહેલાં મુકાય છે.

Claude Code ના પોતાના ઓપ્શન, જેમ કે `--dangerously-skip-permissions` કે `--model`, `--` પછી આવે છે. `ccam <name>` નામ પછીનું બધું આગળ મોકલે છે, પણ `ccam next` પોતાના આર્ગ્યુમેન્ટ્સ જાતે વાંચે છે, તેથી `--` વગર તે `ccam: unknown flag: --dangerously-skip-permissions` સાથે અટકી જાય છે.

## તે કેવી રીતે કામ કરે છે

ccam દ્વારા થતી દરેક શરૂઆત સમયને `~/.local/state/ccam/last-used/<name>` માં લખે છે. `ccam next` આ ફાઇલો વાંચે છે, કયાં એકાઉન્ટ લૉગિન થયેલાં છે તે તપાસે છે, અને પસંદ કરે છે.

1. `cli.candidates` `config.toml` પરથી દરેક એકાઉન્ટ માટે એક `next.Candidate` બનાવે છે: એકાઉન્ટ, તેની છેલ્લી શરૂઆતનો સમય (`launcher.LastUsed`) અને તેની પાસે લૉગિન છે કે નહીં (`creds.Has`: વાંચી શકાય તેવી `.credentials.json`, અથવા macOS પર Keychain આઇટમ).
2. `--auto` વગર, `next.Rotate` એકાઉન્ટને સ્લૉટ નંબર પ્રમાણે ગોઠવે છે અને સૌથી તાજેતરમાં શરૂ થયેલું એકાઉન્ટ શોધે છે. એ જ એકાઉન્ટ છે જેમાંથી તમે સ્વિચ કરી રહ્યા છો. જો હજી કોઈ એકાઉન્ટ શરૂ ન થયું હોય, તો તે સૌથી નાનો સ્લૉટ છે.
3. ત્યાંથી તે સ્લૉટના ક્રમમાં આગળ વધે છે, છેલ્લા સ્લૉટ પછી ફરી પહેલા પર આવે છે, અને લૉગિન થયેલું પહેલું એકાઉન્ટ પાછું આપે છે.
4. `launcher.LaunchWith` પસંદ કરેલું એકાઉન્ટ બરાબર `ccam <name>` ની જેમ તૈયાર કરે છે અને શરૂ કરે છે: તે એકાઉન્ટની લિંક તપાસે છે, MCP શેરિંગ લાગુ કરે છે, સેશન નોંધે છે અને `syscall.Exec` દ્વારા ccam ને `claude` થી બદલે છે.

આ શોધ છેલ્લે વાપરેલા એકાઉન્ટથી શરૂ થાય છે અને તેના પછીના લૉગિન થયેલા પહેલા એકાઉન્ટ પર અટકે છે:

<!-- motion:next-rotate -->

`work` સૌથી તાજેતરમાં શરૂ થયું હતું, તેથી શોધ તેના પછીથી શરૂ થાય છે. `side` પાસે લૉગિન નથી અને તે છોડી દેવાય છે, અને `lab` પસંદ થાય છે. જો `lab` પણ લૉગઆઉટ હોત, તો શોધ ફરીને `personal` પર આવત (ડાયાગ્રામમાં `% len`). rotation ક્યારેય તમે જ્યાંથી શરૂ કર્યું તે એકાઉન્ટ પાછું આપતું નથી: લૉગિન થયેલું બીજું કોઈ એકાઉન્ટ ન હોય, ત્યારે તે ભૂલ સાથે અટકે છે.

### --auto

`--auto` સ્લૉટના ક્રમને બદલે વપરાશ પરથી નક્કી કરે છે. `cli.pickAuto` લૉગિન થયેલાં એકાઉન્ટ રાખે છે, `ccam usage` જે વાપરે છે તે જ `Fetcher.FetchAll` વડે તેમનો બધો વપરાશ એકસાથે પૂછે છે, અને `usage.Policy.Choose()` ને પસંદ કરવા દે છે:

<!-- motion:next-auto -->

એકાઉન્ટને તેની 5 કલાકની અને સાપ્તાહિક વિન્ડો પરથી મૂલવવામાં આવે છે, અને જ્યારે Claude Code માટેના આર્ગ્યુમેન્ટ્સમાં `--model` હોય (`--` પછી) ત્યારે એ મૉડલની સાપ્તાહિક વિન્ડો પરથી પણ (નામમાં `opus`, `sonnet`, `haiku` કે `fable`). `Policy.Choose()` એકાઉન્ટને ત્રણ ભાગમાં વહેંચે છે:

1. **જગ્યા**: દરેક વિન્ડો `auto_max_usage` (`config.toml`, ડિફૉલ્ટ 100) થી નીચે. તેમાંથી strategy પસંદ કરે છે: `headroom` (ડિફૉલ્ટ) સૌથી ઓછો 5 કલાકનો વપરાશ લે છે; `expiring` સૌથી વહેલી રીસેટ થતી 5 કલાકની વિન્ડો લે છે, જેથી તેનો બાકીનો quota ગુમાવાય તે પહેલાં વપરાઈ જાય; `balance` એવું એકાઉન્ટ લે છે જેની પાસે તેનું અઠવાડિયું રીસેટ થાય ત્યાં સુધી દરરોજ માટે સૌથી વધુ સાપ્તાહિક quota બાકી હોય, ચાલુ દરેક ccam સેશનને તેની સામે ગણીને, જેથી સાપ્તાહિક quota વણવપરાયેલો ન રહે અને સમાંતર સેશન વહેંચાઈ જાય.
2. **ઉપર**: 100% થી નીચે પણ `auto_max_usage` પર કે તેથી વધુ. કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે જ વપરાય છે, સૌથી ઓછો 5 કલાકનો વપરાશ પહેલાં.
3. **ભરેલું**: 100% પરની કોઈ વિન્ડો, અથવા અજાણ્યો 5 કલાકનો વપરાશ. ક્યારેય પસંદ થતું નથી.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**સંસ્થાઓ.** આપમેળે થતી પસંદગી તમે છોડો છો તે એકાઉન્ટ જેવાં એકાઉન્ટમાં જ રહે છે: એ જ team કે enterprise સંસ્થા, અથવા અંગત પ્લાન (Pro, Max) એકસાથે, જેથી કંપનીના એકાઉન્ટનું કામ ખાનગી એકાઉન્ટ પર ચાલુ ન થાય, કે ઊલટું. ccam દરેક એકાઉન્ટની સંસ્થા તેના લૉગિન પરથી જાણે છે (`ccam login` તેની નોંધ કરે છે; જૂનાં લૉગિન `.claude.json` માંથી વંચાય છે); જે એકાઉન્ટ વિશે તે કહી ન શકે તે ગમે ત્યાં બંધબેસે છે. એ જ Claude એકાઉન્ટમાં લૉગિન થયેલાં બે ફોલ્ડર તેની મર્યાદા શેર કરે છે, તેથી એક ક્યારેય બીજાનું કામ સંભાળતું નથી. `--cross-org`, અથવા `config.toml` માં `cross_org = true`, સંસ્થાનો નિયમ હટાવે છે.

કંઈ પસંદ ન થઈ શકે ત્યારે ccam code 75 (`EX_TEMPFAIL`) સાથે બહાર નીકળે છે, જેથી script "પછી પ્રયાસ કરો" ને નિષ્ફળતાથી અલગ ઓળખી શકે. `--wait` સાથે, તેના બદલે તે પહેલા ભરેલા એકાઉન્ટની અટકાવતી વિન્ડો રીસેટ થાય ત્યાં સુધી ઊંઘે છે (રીસેટ સમયના 30 સેકન્ડ પછી, જ્યારે endpoint નવી વિન્ડો બતાવે), `ccam: no account has room; waiting until 15:10` છાપે છે, અને ફરી પસંદ કરે છે; જ્યારે કોઈ એકાઉન્ટ ક્યારે રીસેટ થશે તે ન કહે, ત્યારે તે 10 મિનિટ પછી ફરી જુએ છે. રાહ પૂરી થાય તે પહેલાં કંઈ ખાલી ન થાય, તો તે 75 સાથે બહાર નીકળે છે. lookup `ccam usage` ના નિયમો અનુસરે છે: 60 સેકન્ડથી ઓછા જૂના આંકડા ફરી વપરાય છે, દરેક એકાઉન્ટ પાસે 3 સેકન્ડ હોય છે, અને જે એકાઉન્ટ જવાબ ન આપે કે જેના લૉગિનની મુદત પૂરી થઈ ગઈ હોય, તેને તેના સેવ કરેલા આંકડા પરથી મૂલવવામાં આવે છે. જે સેવ કરેલી વિન્ડોનો રીસેટ સમય વીતી ગયો હોય તે 0% ગણાય છે, તેથી ત્યાર પછી જેની વિન્ડો રીસેટ થઈ હોય તે એકાઉન્ટ અવગણાતું નથી. `--auto` ક્યારેય ટોકન રિન્યૂ કરતું નથી. સાદા rotation થી વિપરીત, `--auto` છેલ્લે વાપરેલું એકાઉન્ટ પણ પસંદ કરી શકે છે, જો તેની પાસે હજી સૌથી વધુ જગ્યા હોય.

### --continue

1. ccam દરેક એકાઉન્ટમાં આ ફોલ્ડરની છેલ્લી વાતચીત શોધે છે: `projects/<folder>/` માંની સૌથી નવી `<session>.jsonl` (ફોલ્ડરનું નામ એ કામનું ફોલ્ડર છે જેમાં અક્ષર કે અંક સિવાયનો દરેક અક્ષર `-` માં ફેરવાય છે, જેમ Claude Code તેને નામ આપે છે). જો સેશને કોઈ નોંધ છોડી હોય (વપરાશ મર્યાદા કે `/ccam:switch`, નીચે જુઓ) અને ત્યાર પછી કંઈ નવું લખાયું ન હોય, તો તેના બદલે બરાબર એ જ સેશન વપરાય છે.
2. તે જે એકાઉન્ટની છે તે જ એ એકાઉન્ટ છે જેમાંથી તમે સ્વિચ કરો છો, અને નવું એકાઉન્ટ ઉપર પ્રમાણે, એ જ સંસ્થાની અંદર પસંદ થાય છે.
3. જો નવું એકાઉન્ટ એ ઇતિહાસ શેર કરતું હોય (એ જ ગ્રુપ ફોલ્ડર), તો ccam `claude --resume <session id>` શરૂ કરે છે, એટલે કે બરાબર એ જ વાતચીત (`--continue` વાપરવાથી ફોલ્ડરમાંની સૌથી નવી વાતચીત ખૂલે).
4. નહીં તો તે `claude --resume <the conversation's file> --fork-session` શરૂ કરે છે: Claude Code આખી વાતચીતની નકલ નવા session id હેઠળ નવા એકાઉન્ટના પોતાના ઇતિહાસમાં કરે છે, અને મૂળ જેમ હતી તેમ રહે છે. ત્યાર પછી નવા એકાઉન્ટનું `--continue` તેને શોધી લે છે.

ccam જણાવે છે કે વાતચીત કેટલી મોટી છે: નવું એકાઉન્ટ તેના પહેલા જવાબ સાથે તે આખી ફરી વાંચે છે. prompt caches ક્યારેય સંસ્થાઓ વચ્ચે શેર થતા નથી, તેથી `--cross-org` સાથે એ પહેલા જવાબનો ખર્ચ આખા કદ જેટલો થાય છે. ગ્રુપમાં `history` શેર કરવાથી નકલ ટળે છે, એક જ વાતચીત રહે છે, અને `/rewind` જેમાંથી કોડ પાછો લાવે છે તે checkpoints (`file-history/`) પણ શેર થાય છે.

## Claude Code ની અંદર

[ccam plugin](/docs/install#claude-code-plugin) બે hooks ઉમેરે છે (Claude Code 2.1.78 કે પછીનું):

- જ્યારે કોઈ turn વપરાશ મર્યાદાએ અટકે (`You've hit your session limit`, સાપ્તાહિક કે મૉડલની મર્યાદા), ત્યારે તે વાતચીતની નોંધ કરે છે અને તેને ટેકો આપતાં ટર્મિનલ (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) માં desktop સૂચના બતાવે છે: `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` સર્વરનું ટૂંકું throttling કે API key ની rate limit વપરાશ મર્યાદા નથી, અને તેને એમ જ છોડી દેવાય છે. Claude Code 2.1.234 અને પછીનાં વર્ઝન સેશનમાં રાહ પણ જુએ છે અને 5 કલાકના રીસેટ પછી જાતે ચાલુ રાખે છે; સૂચના તમને તેના બદલે સ્વિચ કરવાનો વિકલ્પ આપે છે.
- `/ccam:switch [account]` વાતચીતની નોંધ કરે છે અને કોઈ request થાય તે પહેલાં turn પૂરો કરે છે: `/exit` લખો, પછી `ccam next --continue` તેને તમે નામ આપેલા એકાઉન્ટ તરીકે, કે આગલા એકાઉન્ટ તરીકે આગળ લઈ જાય છે.

નોંધ `~/.local/state/ccam/handoff/<account>.json` માં 12 કલાક રહે છે અને આગલું સેશન શરૂ થવાનું હોય ત્યારે એક જ વાર વપરાય છે; નિષ્ફળ જતી શરૂઆત તેને આગલા પ્રયાસ માટે રહેવા દે છે. તે જે ફોલ્ડરમાં સેશન શરૂ થયું હતું તે ફોલ્ડરની જ રહે છે, Claude Code ના shell માં `cd` કર્યા પછી પણ. જો તેમાં લખેલા એકાઉન્ટનું ત્યાર પછી નામ બદલાયું હોય કે તે કાઢી નખાયું હોય, તો ccam તે જણાવે છે અને હંમેશની જેમ પસંદ કરે છે. `config.toml` માં `relay = true` સાથે, ccam પોતાની જગ્યા Claude Code ને આપવાને બદલે interactive સેશનની બાજુમાં રહે છે, અને `/ccam:switch` પછી એ સેશનની વાતચીતને જાતે આગળ લઈ જાય છે (બીજા સેશને છોડેલી નોંધ પર તે કંઈ કરતું નથી). તેનો ખર્ચ દરેક સેશન દીઠ એક રાહ જોતી ccam process છે, તેથી તે ડિફૉલ્ટ રીતે બંધ છે. આ રીતે આગળ ચાલતું સેશન `ccam next --continue` થી શરૂ થાય છે, પહેલું સેશન જે આર્ગ્યુમેન્ટ્સ સાથે શરૂ થયું હતું (જેમ કે `--model`) તેમના વગર.

## મર્યાદાએ ફરી ચલાવવું

`ccam next` અટકી ગયેલા interactive સેશન માટે છે. non-interactive run, એટલે કે script કે શેડ્યૂલ કરેલા job માંનું `claude -p`, તમારી રાહ જોઈ શકતું નથી, તેથી એકાઉન્ટના નામ પહેલાં `--fallback` મૂકવાથી ccam તેના માટે સ્વિચ કરે છે:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` ફક્ત `-p` કે `--print` સાથે જ કામ કરે છે; નહીં તો ccam અટકી જાય છે અને `ccam next --auto --continue` તરફ ધ્યાન દોરે છે.
2. standard input એક વાર વંચાય છે અને દરેક પ્રયાસને અપાય છે.
3. ccam `exec` ને બદલે `claude` ને child process તરીકે શરૂ કરે છે, તેનું standard output રોકી રાખે છે, અને standard error આવે તેમ આગળ આપે છે.
4. `claude` 0 સાથે બહાર નીકળે, કે બીજા કોઈ કારણે નિષ્ફળ જાય, ત્યારે રોકેલું output અને exit code બદલ્યા વગર આગળ અપાય છે.
5. નિષ્ફળતાને વપરાશ મર્યાદા ત્યારે ગણાય છે જ્યારે Claude Code નો સંદેશ એમ કહે (`You've hit your limit`, `usage limit reached`), અથવા જ્યારે એકાઉન્ટનો વપરાશ lookup તેની 5 કલાકની કે સાપ્તાહિક વિન્ડો 100% પર બતાવે.
6. પછી એ પ્રયાસનું output છોડી દેવાય છે અને એ જ command બીજા લૉગિન થયેલા એકાઉન્ટ તરીકે ચાલે છે: સૌથી ઓછા 5 કલાકના વપરાશવાળું, પછી જેમનો વપરાશ જાણી ન શકાય તે એકાઉન્ટ, નંબરના ક્રમમાં. 100% પરનાં એકાઉન્ટ છોડી દેવાય છે, અને દરેક એકાઉન્ટ વધુમાં વધુ એક વાર અજમાવાય છે.
7. આર્ગ્યુમેન્ટ્સમાં `-c`, `--continue`, `-r` કે `--resume` હોય ત્યારે ફક્ત એ જ history ગ્રુપનાં એકાઉન્ટ અજમાવાય છે, જેથી વાતચીત મળી શકે.
8. `--auto` ની જેમ, ફક્ત એ જ સંસ્થાનાં એકાઉન્ટ અજમાવાય છે (`cross_org = true` આ હટાવે છે).

આગલું એકાઉન્ટ `--auto` ના જ નિયમોથી પસંદ થાય છે (`auto_max_usage`, `auto_strategy`, `--model` ની વિન્ડો). કોઈ એકાઉન્ટ ન બચે ત્યારે ccam છેલ્લા પ્રયાસનું output છાપે છે, કહે છે કે કોઈ લૉગિન થયેલા એકાઉન્ટમાં જગ્યા બાકી નથી, અને code 75 સાથે બહાર નીકળે છે. એકાઉન્ટના નામ પહેલાં `--wait` સાથે (`ccam --fallback --wait=2h work -p …`), તે તેના બદલે પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, પછી બધાં એકાઉન્ટ ફરી અજમાવી શકે છે; તે હંમેશાં પહેલાં ઓછામાં ઓછું એક વાર ઊંઘે છે, જેથી હમણાં જ મર્યાદાએ પહોંચેલું એકાઉન્ટ તરત ફરી ન ચાલે. `ccam --dry-run --fallback work -p …` પહેલા પ્રયાસનો command છાપે છે.

interactive સેશન જાતે સ્વિચ થતાં નથી. તેના બદલે, `ccam statusline install` સાથે, 5 કલાકનો કે સાપ્તાહિક વપરાશ 90% પર પહોંચે એટલે સ્ટેટસ લાઇન `→ ccam next --auto --continue` ઉમેરે છે ([વપરાશ મર્યાદા](/features/usage)).

## સલામતી અને ખાસ કિસ્સા

- લૉગિન વગરનાં એકાઉન્ટ ક્યારેય પસંદ થતાં નથી. લૉગિન થયેલું બીજું કોઈ એકાઉન્ટ ન હોય, ત્યારે સાદું `ccam next` `no logged-in account other than <name>` સાથે અટકે છે.
- `--continue` સાથે અને ફોલ્ડરમાં કોઈ વાતચીત ન હોય ત્યારે: `no account has a conversation in <folder> to continue`.
- જ્યારે સંસ્થાનો નિયમ કોઈ એકાઉન્ટ બાકી ન રાખે: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` સાથે અને લૉગિન થયેલું કોઈ એકાઉન્ટ ન હોય ત્યારે: `no logged-in account to pick`. જ્યારે દરેક એકાઉન્ટ ભરાઈ ગયું હોય કે તેનો વપરાશ અજાણ્યો હોય: `` no account with known usage and headroom (see `ccam usage`) ``. બંને exit code 75 સાથે પૂરા થાય છે.
- `--max-usage`, `--strategy`, `--wait` અને `--cross-org` ફક્ત `--auto` સાથે કામ કરે છે, અને 1-100 ની બહારનું મૂલ્ય કે બીજી strategy નકારવામાં આવે છે.
- fork વાતચીતની નકલ કરે છે, એકાઉન્ટની બાકીની સ્થિતિની નહીં: ઇતિહાસ શેર ન થતો હોય તો મૂળ એકાઉન્ટના `file-history/` ના checkpoints સાથે આવતા નથી.
- `--auto` `--dry-run` સાથે પણ નેટવર્ક request કરે છે. તે ફક્ત વાંચે છે, અને ક્યારેય ટોકન ફાઇલો લખતું નથી.
- પસંદ કરેલું એકાઉન્ટ `ccam <name>` ના જ પાથથી શરૂ થાય છે, તેથી એ જ તપાસો લાગુ પડે છે, જેમાં એ જ ઇતિહાસ શેર કરતું બીજું Claude Code આ ફોલ્ડરમાં પહેલેથી ચાલતું હોય ત્યારની ચેતવણી પણ સામેલ છે.
- `last-used` ત્યારે જ લખાય છે જ્યારે Claude Code ખરેખર શરૂ થાય. `--dry-run` તેને અડતું નથી, તેથી dry run rotation ને આગળ ખસેડતું નથી.
- `--fallback` પ્રયાસ પૂરો થાય ત્યાં સુધી standard output રોકી રાખે છે, તેથી `--output-format stream-json` જેવું streaming output અંતે એકસાથે આવે છે. મર્યાદાએ પહોંચતાં પહેલાં પહેલા પ્રયાસે બદલેલી ફાઇલો બદલાયેલી જ રહે છે; આગલો પ્રયાસ prompt ફરીથી શરૂ કરે છે.

## સંબંધિત

- [વપરાશ મર્યાદા](/features/usage)
- [વપરાશ મર્યાદા અને એકાઉન્ટ સ્વિચ કરવું](/docs/usage-limits)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [વપરાશ મર્યાદા સાથે કામ કરવાની રીત](/guides/usage-limit-workflow)
