---
title: Claude Code શરૂ કરવું
description: ccam <name|number>, ccam run કે ccam picker વડે કોઈ પણ એકાઉન્ટ તરીકે Claude Code શરૂ કરો; આર્ગ્યુમેન્ટ્સ બદલાયા વગર claude ને જાય છે, --dry-run કમાન્ડ બતાવે છે.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work` Claude Code ને `work` એકાઉન્ટ તરીકે શરૂ કરે છે, અને નામ પછીનું બધું જેમ છે તેમ `claude` ને જાય છે. જ્યાં તમે `claude` લખતા હો ત્યાં તેને વાપરો: ટર્મિનલમાં, સ્ક્રિપ્ટમાં, કે pipe થી input આપીને.

## ઉપયોગ

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

એકાઉન્ટનું નામ કે `ccam list` માંનો તેનો નંબર આપો. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી, તેથી `ccam work --help` Claude Code ની મદદ બતાવે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે. standard input, exit code અને સિગ્નલ `claude` જેવું જ વર્તન કરે છે.

જ્યારે પહેલો આર્ગ્યુમેન્ટ `-` થી શરૂ થાય અને ccam ના પોતાના flag માંનો ન હોય, ત્યારે ccam બધા આર્ગ્યુમેન્ટ્સ સાથે નામ વગરનું એકાઉન્ટ શરૂ કરે છે. તે [`ccam folder`](/features/folder) વડે હાલના ફોલ્ડર માટે સેટ કરેલું એકાઉન્ટ છે, નહીં તો ડિફૉલ્ટ: `ccam use` વડે સેટ કરેલું એકાઉન્ટ, અથવા તમે સેટ ન કર્યું હોય તો સૌથી નાના નંબરવાળું એકાઉન્ટ.

### picker

standard input અને output ટર્મિનલ હોય ત્યારે આર્ગ્યુમેન્ટ વગરનું `ccam` picker ખોલે છે, નહીં તો `ccam list` છાપે છે. તે દરેક એકાઉન્ટનો નંબર, ઇમેઇલ, 5 કલાકનો વપરાશ, છેલ્લો ઉપયોગ અને શેર કરેલી આઇટમ બતાવે છે; દરેક એકાઉન્ટ જવાબ આપે તેમ વપરાશ ભરાતો જાય છે.

| કી | શું કરે છે |
| --- | --- |
| `Enter` | પસંદ કરેલા એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે |
| `1` થી `9` | એ નંબરવાળું એકાઉન્ટ શરૂ કરે છે |
| `Up`/`Down`, `k`/`j` | પસંદગી ખસેડે છે |
| `a` | picker બંધ કરીને `ccam add` વિઝાર્ડ ચલાવે છે |
| `s` | પસંદ કરેલા એકાઉન્ટની શેર કરેલી આઇટમ બદલે છે (`Space` કે `x` ટૉગલ કરે છે, `Enter` લાગુ કરે છે, `Esc` પાછું જાય છે) |
| `d` | પસંદ કરેલા એકાઉન્ટને ડિફૉલ્ટ બનાવે છે |
| `r` | ટૂંકા કૅશ વગર વપરાશ ફરીથી મેળવે છે |
| `q`, `Esc`, `Ctrl+C` | બહાર નીકળે છે |

picker કંઈ કરતાં પહેલાં બંધ થાય છે, તેથી Claude Code શરૂ થાય તે પહેલાં ટર્મિનલ પાછું સામાન્ય થઈ જાય છે. `d` કે શેરિંગના ફેરફાર પછી તે ફરી ખૂલે છે. શેરિંગનો ફેરફાર `ccam share` ની જેમ પહેલાં યોજના બતાવે છે અને પૂછે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--dry-run` | ccam જે environment અને કમાન્ડ ચલાવશે તે છાપીને બહાર નીકળે છે | off |
| `--fallback` | `-p`/`--print` runs માટે: એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે એ જ command સૌથી વધુ જગ્યાવાળા લૉગિન થયેલા એકાઉન્ટ તરીકે ફરી ચલાવે છે ([વિગતો](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | `ccam <version>` છાપે છે (ફક્ત પહેલા આર્ગ્યુમેન્ટ તરીકે) | |
| `-h`, `--help` | ccam ની મદદ, `ccam --help` કે `ccam run --help` તરીકે | |

આર્ગ્યુમેન્ટના નિયમો:

- ccam પોતાના flag ફક્ત એકાઉન્ટના નામ પહેલાં જ વાંચે છે. નામ પછી, `--dry-run` અને `--help` સહિત દરેક આર્ગ્યુમેન્ટ Claude Code નો છે.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work` અને `ccam --dry-run` (નામ વગરનું એકાઉન્ટ) તરીકે કામ કરે છે. તેને માનવામાં આવે છે, ઇનકાર થતો નથી. `--dry-run` અને `--fallback` કોઈ પણ ક્રમમાં આવી શકે.
- `ccam run <name|number> [claude args]` એ `ccam <name>` નું લાંબું સ્વરૂપ છે. નામ વગરનું `ccam run` ફોલ્ડરનું એકાઉન્ટ કે ડિફૉલ્ટ એકાઉન્ટ વાપરે છે.
- એકાઉન્ટનાં નામ `add`, `list`, `run` કે `usage` જેવાં ccam કમાન્ડનાં નામ ન હોઈ શકે, તેથી `ccam <word>` ક્યારેય અસ્પષ્ટ નથી.

dry run એક લાઇન છાપે છે, જે તમે shell માં, કે Windows પર PowerShell માં, પેસ્ટ કરી શકો:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## તે કેવી રીતે કામ કરે છે

દરેક શરૂઆત કમાન્ડ લાઇનથી `exec` સુધી એ જ પગલાં લે છે:

<!-- motion:launch -->

દરેક બૉક્સ તેનું પગલું ચાલે ત્યારે પ્રકાશિત થાય છે, અને નીચેની લાઇનો એ environment અને આર્ગ્યુમેન્ટ્સ છે જેની સાથે Claude Code શરૂ થાય છે.

1. `cmd/ccam/main.go` `app.Main` ને બોલાવે છે. જો ccam બીજા નામે શરૂ થયું હોય, તો તે [શૉર્ટકટ](/features/shortcuts) છે અને સીધું launcher પર જાય છે.
2. cobra સામેલ થાય તે પહેલાં `app.parse` `os.Args` જુએ છે. ફક્ત એકલું `ccam` (picker), ccam કમાન્ડનું નામ (`config.Reserved`), `-h`/`--help` અને `-v`/`--version` `cli.Execute` પર જાય છે. બાકીનું બધું શરૂઆત છે.
3. `launcher.Launch` `~/.config/ccam/config.toml` વાંચે છે, જે એકમાત્ર ફાઇલ તે parse કરે છે.
4. `launcher.Prepare` નામ કે નંબરથી એકાઉન્ટ શોધે છે (નામ ન હોય તો `HereProfile` ફોલ્ડરનું એકાઉન્ટ લે છે, working directory ને `config.toml` ની `[[folder]]` entries સાથે strings તરીકે સરખાવીને), શરૂઆતનું એક `--` કાઢી નાખે છે, અને `ResolveClaude` વડે `claude` executable શોધે છે. PATH lookup ફાઇલના modification time સાથે `~/.local/state/ccam/claude-path` માં કૅશ થાય છે, તેથી ccam દરેક શરૂઆતે `PATH` શોધતું નથી. બદલાયેલી કે ગુમ થયેલી ફાઇલ નવું lookup કરાવે છે. `config.toml` માંનો `claude_bin` પાથ જેમ છે તેમ વપરાય છે. જે `claude` પોતે ccam હોય તેનો ccam inode વડે સરખામણી કરીને ઇનકાર કરે છે, કારણ કે તે loop બનાવે.
5. `checkLinks` એકાઉન્ટ ફોલ્ડર `~/.local/share/ccam/profiles/<name>/` ની તેનાં શેરિંગ સેટિંગ્સ સાથે `share.Check` વડે સરખામણી કરે છે, જે ફક્ત `lstat` અને `readlink` ને બોલાવે છે, આઇટમ દીઠ એક કે બે system call. ખૂટતી લિંક ત્યાં જ બને છે. ડેટા ગુમાવી શકે તેવું કંઈ પણ શરૂઆત અટકાવે છે.
6. `applyMCP` શેર કરેલા MCP સર્વર ઉમેરે છે (નીચે જુઓ).
7. `session.Register` સેશન નોંધે છે. ccam પોતાના PID ના નામવાળું, process ના શરૂ થવાના સમય સાથેનું marker `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` માં (અથવા `XDG_RUNTIME_DIR` વગર `~/.local/state/ccam/locks/` માં) લખે છે. વાતચીતનો ઇતિહાસ શેર કરતા એકાઉન્ટ માટે, `historyCheck` પહેલાં `history-<group>-<hash>` માં ચાલતું સેશન શોધે છે, જ્યાં hash હાલના ફોલ્ડરનો છે, અને ત્યાં પણ નોંધણી કરે છે. `exec` PID રાખે છે, તેથી પછી marker Claude Code process ને જ દર્શાવે છે. પકડી રાખેલું lock નહીં પણ marker, કારણ કે Claude Code જે process શરૂ કરે તે lock વારસામાં મેળવે અને સેશન પૂરું થયા પછી પણ રાખે.
8. `Plan.Exec` ટર્મિનલનું શીર્ષક `claude · <name>` કરે છે (standard error ટર્મિનલ હોય ત્યારે), શરૂઆતનો સમય `fsync` વગર `~/.local/state/ccam/last-used/<name>` માં લખે છે, અને `syscall.Exec` ને બોલાવે છે.

Claude Code `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર પર અને `CCAM_PROFILE` ને એકાઉન્ટના નામ પર સેટ કરીને ચાલે છે. ccam આના વારસામાં મળેલાં મૂલ્યો દૂર કરે છે. macOS પર `ccam import` વડે રજિસ્ટર કરેલા `~/.claude` ફોલ્ડર માટે, તે `CLAUDE_SECURESTORAGE_CONFIG_DIR` ને ખાલી મૂલ્ય પર પણ સેટ કરે છે, જેથી એકાઉન્ટ અને સાદું `claude` એક જ Keychain લૉગિન વાપરે.

### exec શા માટે

`syscall.Exec` ccam process ને `claude` થી બદલી નાખે છે, તેથી ccam memory માં રહેતું નથી, keystroke, સિગ્નલ અને exit code Claude Code ના પોતાના હોય છે, અને આર્ગ્યુમેન્ટ્સ array જ રહે છે, ક્યારેય shell string બનતા નથી.

Windows માં `exec` નથી. ત્યાં `Plan.Exec` એ જ environment, આર્ગ્યુમેન્ટ્સ અને console સાથે `claude.exe` ને child process તરીકે શરૂ કરે છે, તેની રાહ જુએ છે અને તેના exit code સાથે બહાર નીકળે છે. Ctrl+C console માંથી સીધું Claude Code સુધી પહોંચે છે, તેથી ccam તેને અવગણે છે. સેશન marker માં ccam નો PID હોય છે, જે બરાબર Claude Code જેટલો જ સમય જીવે છે, અને Claude Code બહાર નીકળે ત્યારે marker દૂર થાય છે. રાહ જોતી વખતે ccam બીજું કંઈ કરતું નથી. સેશન પછીનું કામ, જેમ કે `claude mcp add` વડે ઉમેરેલા સર્વર ધ્યાનમાં લેવા, આગલી શરૂઆતે થાય છે: ccam `.claude.json` ના modification time ની `~/.local/state/ccam/mcp/<name>/seen` સાથે સરખામણી કરે છે અને ફાઇલ બદલાઈ હોય ત્યારે જ તેને parse કરે છે.

### કમાન્ડ લાઇન કયા રસ્તે જાય છે

ત્રણ કમાન્ડ લાઇન, `internal/app` માંથી ત્રણ રસ્તા:

<!-- motion:launch-route -->

શૉર્ટકટ `argv[0]` પરથી ઓળખાય છે, `ccam <name>` પહેલા આર્ગ્યુમેન્ટ પરથી, અને બંને `syscall.Exec` માં પૂરા થાય છે. ફક્ત `ccam usage` જેવો ccam કમાન્ડ જ cobra સુધી પહોંચે છે.

### શેર કરેલા MCP સર્વર

ડિફૉલ્ટ `flag` વ્યૂહરચના સાથે, ccam ગ્રુપના MCP સર્વર Claude Code ને આર્ગ્યુમેન્ટ તરીકે, તમારા આર્ગ્યુમેન્ટ્સની પહેલાં, આપે છે:

<!-- motion:launch-mcp -->

ગ્રુપની `mcp.json` હોય તો જ `mcp.FlagArgs` આર્ગ્યુમેન્ટ આપે છે, અને `Plan.InsertArgs` તેને `claude` પછી તરત મૂકે છે. તમારું પોતાનું `--mcp-config` કે `--strict-mcp-config` પછી પણ તેના પર લાગુ પડે છે. તે એક જ આર્ગ્યુમેન્ટ, `--mcp-config=<path>`, છે, કારણ કે Claude Code નું `--mcp-config` ઘણાં મૂલ્યો લે છે: બે આર્ગ્યુમેન્ટ તરીકે લખેલું હોય, તો તે પછીના prompt ને પણ config ફાઇલ તરીકે વાંચે.

`config.toml` માં `mcp_strategy = "sync"` સાથે, ccam તેના બદલે સર્વર એકાઉન્ટની `.claude.json` માં લખે છે, અને એકાઉન્ટ પહેલેથી ચાલતું હોય તો એ છોડી દે છે.

### Remote Control

`ccam work rc` (અથવા `ccam work remote-control`) Claude Code નું Remote Control `work` તરીકે શરૂ કરે છે. `claude remote-control` તેના verb પહેલાં `--mcp-config` નકારે છે અને પછી એક પણ લેતું નથી, અને તેનાં સેશન MCP સર્વર ફક્ત એકાઉન્ટની `.claude.json` માંથી વાંચે છે. તેથી અહીં flag વ્યૂહરચના જુદી રીતે કામ કરે છે:

<!-- motion:launch-rc -->

1. `remoteControl` એકાઉન્ટ પછીના, `-` થી શરૂ ન થતા પહેલા આર્ગ્યુમેન્ટને Claude Code ના verb તરીકે લે છે. તે `remote-control` કે `rc` હોય ત્યારે, `mcp.FlagArgs` કશું ઉમેરતું નથી.
2. `mcp.Synced` ગ્રુપની `mcp.json` ના modification time ની એકાઉન્ટના `sync-stamp` સાથે સરખામણી કરે છે. જો `.claude.json` માં પહેલેથી ગ્રુપના હાલના સર્વર હોય, અથવા ગ્રુપમાં કોઈ સર્વર ન હોય, તો Remote Control કોઈ પ્રશ્ન વગર શરૂ થાય છે.
3. નહીંતર સેશનને જે સર્વર નહીં મળે તેમનાં નામ ccam જણાવે છે અને પૂછે છે:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` પર, `mcp.Sync` તેમને `sync` વ્યૂહરચનાની જેમ જ કૉપિ કરે છે: Claude Code ના lock હેઠળ, બૅકઅપ પછી, એકાઉન્ટના પોતાના સર્વર જાળવીને. `N` પર, અથવા ટર્મિનલ વગર, Remote Control તેમના વગર શરૂ થાય છે, અને ccam `mcp_strategy = "sync"` તરફ ધ્યાન દોરે છે, જે દરેક શરૂઆત પહેલાં તેમને કૉપિ કરે છે.

સર્વર કૉપિ ન થાય ત્યાં સુધી ccam દરેક Remote Control શરૂઆતે પૂછે છે, અને કૉપિ પછી ફક્ત ત્યારે જ, જ્યારે ગ્રુપની `mcp.json` ફરી બદલાય. પછીની સામાન્ય શરૂઆતોમાં `.claude.json` ની કૉપિ `--mcp-config` સામે જીતે છે, તેથી `ccam mcp add`, `remove` અને `import` બીજાં એકાઉન્ટ પાસે રહેલી કૉપિ અપડેટ કરે છે ([શેર કરેલા MCP સર્વર](/features/mcp)).

### 10 ms ની મર્યાદા

Claude Code શરૂ થાય તે પહેલાં ccam વધુમાં વધુ 10 ms (p95) ઉમેરે છે. શરૂઆતનો પાથ ક્યારેય cobra command tree બનાવતો નથી, Bubble Tea શરૂ કરતો નથી કે નેટવર્ક કનેક્શન ખોલતો નથી, અને package ના `init()` functions કુલ 5 ms થી વધુ લે તો `TestPackageInitIsCheap` નિષ્ફળ જાય છે. `scripts/perf-check.sh` `hyperfine` વડે `ccam --dry-run work`, `ccam work` અને એક શૉર્ટકટ માપે છે અને 10 ms થી ઉપર નિષ્ફળ જાય છે, અને `internal/launcher` માંનું `BenchmarkLaunchPath` પગલાં loop માં ચલાવે છે. દરેક પગલાનો સમય જોવા માટે `CCAM_TRACE=1` સેટ કરો:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## સલામતી અને ખાસ કિસ્સા

- અજાણ્યું નામ કે નંબર `no account "<name>"` સાથે અટકે છે અને `ccam list` તરફ નિર્દેશ કરે છે. કોઈ એકાઉન્ટ જ ન હોય, ત્યારે ccam તમને `ccam add <name>` વડે એક ઉમેરવા કહે છે.
- જો કોઈ શેર કરેલો પાથ સામાન્ય ફાઇલ બની ગયો હોય કે ખોટી જગ્યાએ નિર્દેશ કરતો હોય, તો ccam `stopped because of share link problems:` સાથે અટકે છે, તેમની યાદી આપે છે, અને `ccam doctor --fix` ચલાવવા કહે છે. શરૂઆતના પાથ પર તે ડેટા ખસેડતું નથી.
- જો આ એકાઉન્ટ સાથે વાતચીતનો ઇતિહાસ શેર કરતું (એ જ ગ્રુપનું) બીજું Claude Code એ જ ફોલ્ડરમાં ચાલતું હોય, તો ccam ચેતવણી આપે છે કે ઇતિહાસ ભળી શકે છે અને `Run anyway? [y/N]` પૂછે છે. ટર્મિનલ વગર તે ચેતવણી બતાવીને આગળ વધે છે.
- `--dry-run` MCP sync, નવા સર્વરની સૂચના, ઇતિહાસની તપાસ, સેશન marker, શીર્ષક અને છેલ્લા ઉપયોગનો સમય છોડી દે છે, અને exec કરતું નથી. તે તો પણ શેર લિંક તપાસે છે, અને સામાન્ય શરૂઆતની જેમ ખૂટતી લિંક બનાવે છે.
- `ccam work rc`, જ્યારે `work` પહેલેથી ચાલતું હોય, ત્યારે `.claude.json` ને અડતું નથી, એ જણાવે છે, અને ગ્રુપના સર્વર વગર Remote Control શરૂ કરે છે; તેમને કૉપિ કરવા માટે સેશન પૂરું થયા પછી તેને ફરી ચલાવો. `--dry-run` સાથે તે સંદેશ છાપે છે પણ ન તો પૂછે છે, ન કૉપિ કરે છે.
- `mcp_strategy = "sync"` સાથે, નિષ્ફળ sync `ccam: MCP sync failed (launching anyway): …` છાપે છે અને Claude Code પાસે પહેલેથી હતા તે સર્વર સાથે શરૂ થાય છે.

## સંબંધિત

- [શૉર્ટકટ કમાન્ડ](/features/shortcuts)
- [ccam સાથે Remote Control વાપરવું](/guides/remote-control)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [કમાન્ડ સંદર્ભ](/docs/commands)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [CLAUDE_CONFIG_DIR કેવી રીતે કામ કરે છે](/guides/claude-config-dir)
