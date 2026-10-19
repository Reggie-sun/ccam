---
title: એકાઉન્ટનું સંચાલન
description: ccam add, ccam login, ccam list, ccam rename, ccam remove અને ccam import વડે Claude Code એકાઉન્ટ ઉમેરો, લૉગિન કરો, યાદી જુઓ, નામ બદલો, દૂર કરો અને ઇમ્પોર્ટ કરો.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

દરેક ccam એકાઉન્ટ પોતાના લૉગિન સાથેનું અલગ Claude Code config ફોલ્ડર છે. આ કમાન્ડ એકાઉન્ટ બનાવે છે, તેમાં સાઇન ઇન કરે છે, તેમની યાદી બતાવે છે, નામ બદલે છે અને તેમને દૂર કરે છે, તથા તમારી પાસે પહેલેથી હોય તેવાં ફોલ્ડર રજિસ્ટર કરે છે.

## ઉપયોગ

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

જ્યાં પણ કોઈ કમાન્ડ `<name|number>` લે છે, ત્યાં તમે એકાઉન્ટનું નામ કે તેનો સ્લૉટ નંબર, એટલે કે `ccam list` ની `#` કૉલમ, આપી શકો છો.

ટર્મિનલ પર કોઈ ઓપ્શન વગરનું `ccam add` એક વિઝાર્ડ ખોલે છે. તે નામ પૂછે છે, પછી શેર કરેલી આઇટમનું ચેકલિસ્ટ બતાવે છે (વાતચીતના ઇતિહાસ સિવાય બધું ચેક કરેલું હોય છે), પછી `Create a shortcut command for work? [y/N]` પૂછે છે અને `y` આપો તો `Shortcut name [ccwork]:` પૂછે છે. `--share`, `--shortcut` કે `--no-login` માંથી કોઈ પણ ઓપ્શન, અથવા ટર્મિનલને બદલે pipe, વિઝાર્ડ છોડી દે છે; ત્યારે નામ આપવું જરૂરી છે, અને `--shortcut` કોઈ નામ ન આપે ત્યાં સુધી કોઈ શૉર્ટકટ બનતો નથી.

`ccam list` સ્લૉટ, નામ, શૉર્ટકટ, ઇમેઇલ, શેર કરેલી આઇટમ અને છેલ્લો ઉપયોગ બતાવે છે, અને ડિફૉલ્ટ એકાઉન્ટની આગળ `*` મૂકે છે:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` સાથે, દરેક એકાઉન્ટ `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (એકાઉન્ટ ક્યારેય શરૂ ન થયું હોય ત્યારે છોડી દેવાય છે) અને `share` ધરાવતો ઑબ્જેક્ટ છે. `share` દરેક આઇટમ (`mcp`, `settings`, `memory`, `extensions`, `history`) ને તેના ગ્રુપના નામ સાથે અથવા `"local"` સાથે જોડે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `ccam add --no-login` | લૉગિન કર્યા વગર એકાઉન્ટ બનાવે છે; પછીથી `ccam login <name>` વડે લૉગિન કરો | off |
| `ccam add --share <item=group>` | કોઈ આઇટમ ગ્રુપ સાથે શેર કરે છે, અથવા તેને ખાનગી રાખવા `item=local`. ઘણી વાર આપી શકાય | `mcp`, `settings`, `memory`, `extensions` `default` માં; `history` local |
| `ccam add --shortcut <command>` | `~/.local/bin` માં આ શૉર્ટકટ કમાન્ડ પણ બનાવે છે | કોઈ નહીં |
| `ccam list --json` | એકાઉન્ટને JSON તરીકે છાપે છે | off |
| `ccam remove -y`, `--yes` | પૂછ્યા વગર દૂર કરે છે | off |
| `ccam import --name <name>` | તમે આપેલા ફોલ્ડર માટે એકાઉન્ટનું નામ | ફોલ્ડરના નામ પરથી |
| `ccam import -y`, `--yes` | મળેલું દરેક ફોલ્ડર પૂછ્યા વગર ઇમ્પોર્ટ કરે છે | off |
| `--dry-run` | `ccam list` અને `ccam completion` તેને માને છે, કારણ કે તેઓ કંઈ બદલતા નથી. `add`, `login`, `logout`, `use`, `rename`, `remove` અને `import` તેનો `` --dry-run is not supported by `ccam add`; nothing was changed `` (કમાન્ડના નામ સાથે) સંદેશ સાથે ઇનકાર કરે છે | off |

`ccam login`, `ccam logout`, `ccam use` અને `ccam rename` ના પોતાના કોઈ ઓપ્શન નથી. `ccam completion` `bash` કે `zsh` લે છે; બીજું કંઈ પણ `only bash and zsh are supported` સાથે નિષ્ફળ જાય છે.

એકાઉન્ટનાં નામ `a-z`, `0-9`, `-` અને `_` ના 1 થી 32 અક્ષરનાં હોય છે, અને અક્ષર કે અંકથી શરૂ થાય છે. નામ ફક્ત અંકોનું ન હોઈ શકે, કારણ કે `ccam 2` નો અર્થ સ્લૉટ 2 છે. નામ કોઈ ccam કમાન્ડ પણ ન હોઈ શકે (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), તેથી `ccam <name>` ક્યારેય અસ્પષ્ટ નથી. આ યાદી `internal/config` માં `config.Reserved` છે.

## તે કેવી રીતે કામ કરે છે

`ccam add work --shortcut ccwork` આ પગલાં `addProfile` (`internal/cli/profiles.go`) માં ચલાવે છે:

<!-- motion:accounts-add -->

1. કંઈ પણ બને તે પહેલાં `config.ValidateName` નામ તપાસે છે, અને `shim.Check` શૉર્ટકટ તપાસે છે.
2. `profile.Add` `~/.local/share/ccam/profiles/work/` ને mode `0700` સાથે બનાવે છે અને એકાઉન્ટને સૌથી નાનો ખાલી સ્લૉટ નંબર (`Config.NextSlot`) આપે છે, તેથી દૂર કરેલા એકાઉન્ટનો નંબર ફરી વપરાય છે. પહેલું એકાઉન્ટ ડિફૉલ્ટ બને છે.
3. `applyShare` નવા, ખાલી ફોલ્ડરને `share.Build` અને પછી `Apply` વડે તેના શેર ગ્રુપ સાથે લિંક કરે છે, ઉદાહરણ તરીકે `settings.json` ને `groups/default/settings.json` સાથે. જો આ નિષ્ફળ જાય તો નવું ફોલ્ડર દૂર થાય છે, કારણ કે બાકી રહેલું ફોલ્ડર દરેક ફરી પ્રયાસને `folder already exists` સાથે અટકાવી દે.
4. એકાઉન્ટ `fsutil.AtomicWrite` દ્વારા `~/.config/ccam/config.toml` માં સેવ થાય છે: એક કામચલાઉ ફાઇલ, `fsync`, પછી `rename`, તેથી crash થાય તો પણ config ક્યારેય અડધું રહેતું નથી.
5. `runLogin` `CLAUDE_CONFIG_DIR` ને એકાઉન્ટ ફોલ્ડર પર સેટ કરીને `claude auth login` શરૂ કરે છે. જો લૉગિન નિષ્ફળ જાય, તો એકાઉન્ટ રહે છે, અને ccam ભૂલ છાપે છે અને પછીથી `ccam login work` ચલાવવા કહે છે.
6. `profile.ReadEmail` `profiles/work/.claude.json` માંથી `oauthAccount.emailAddress` કાઢે છે. એ ફાઇલ કેટલાક મેગાબાઇટની હોઈ શકે છે, તેથી `TopLevelValue` આખી ફાઇલ parse કરવાને બદલે એ એક key શોધે છે.
7. `shim.Install` `~/.local/bin/ccwork` ને ccam binary સાથે લિંક કરે છે, જે શરૂ થાય ત્યારે `config.toml` માં `argv[0]` શોધે છે.

લૉગિનમાં મિનિટો લાગી શકે છે, અને એ દરમિયાન બીજું ccam config સેવ કરી શકે છે. તેથી રાહ જોતા કમાન્ડ `env.update` દ્વારા સેવ કરે છે, જે `config.toml` ફરી લોડ કરે છે, એક ફેરફાર (અહીં, ઇમેઇલ) લાગુ કરે છે અને સેવ કરે છે. આ જ કારણે `ccam remove` પોતાના પ્રશ્ન પછી ફરી લોડ કરે છે.

`ccam login` હાલના એકાઉન્ટ માટે પગલાં 5 અને 6 ચલાવે છે. `ccam logout` એકાઉન્ટમાં `claude auth logout` ચલાવે છે અને સેવ કરેલો ઇમેઇલ સાફ કરે છે. `ccam use` ફક્ત `config.toml` માં `default` ફરીથી લખે છે.

### રિમોટ મશીન પર સાઇન ઇન કરવું

SSH પર, `claude auth login` નો બ્રાઉઝરવાળો ભાગ તૂટી શકે છે. Claude Code સર્વર પર `http://localhost:<port>/callback` ના callback સર્વર પર સાઇન-ઇનની રાહ જુએ છે, પણ સાઇન-ઇન પૂરું કરતું બ્રાઉઝર તમારા કમ્પ્યુટર પર ચાલે છે, જ્યાં એ સરનામું ક્યાંય જતું નથી અને બ્રાઉઝર `ERR_CONNECTION_REFUSED` જેવી ભૂલ બતાવે છે. સરનામામાં તો પણ એક વખતનો કોડ હોય છે, તેથી ccam તે પહોંચાડી શકે છે:

<!-- motion:accounts-relay -->

`runLogin` તમારું ટર્મિનલ Claude Code ને સોંપતું નથી. તે Claude Code ના standard input ને એક pipe સાથે જોડે છે અને `relayLogin` તમે લખો છો તે દરેક લાઇન તેમાં કૉપિ કરે છે, તેથી Claude Code ના `Paste code here` પ્રશ્ન પર પેસ્ટ કરેલો કોડ પણ પહોંચે છે. જ્યારે કોઈ લાઇન callback સરનામું હોય, ત્યારે `callbackURL` તેને ફક્ત ત્યારે જ સ્વીકારે છે જ્યારે host `localhost`, `127.0.0.1` કે `[::1]` હોય, port આપેલો હોય, path `/callback` હોય અને `code` તથા `state` બંને હાજર હોય; `http://` છોડી શકાય છે. પછી `deliverCallback` સર્વર પર એ સરનામાને 10 સેકન્ડની મર્યાદા સાથે અને redirect અનુસર્યા વગર request મોકલે છે, જેમ બ્રાઉઝરે કર્યું હોત. request રાહ જોતા `claude auth login` ના callback સર્વર સુધી પહોંચે છે (ડાયાગ્રામમાં તૂટક લાઇન). Claude Code કોડ મેળવે છે અને લૉગિન સેવ કરે છે, Linux પર એકાઉન્ટ ફોલ્ડરની `.credentials.json` માં. પેસ્ટ કરેલી લાઇન ccam ને ક્યારેય બીજા host નો સંપર્ક કરાવી શકતી નથી.

લૉગિન શરૂ થાય તે પહેલાં ccam આ છાપે છે:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

સફળ થાય ત્યારે તે `→ Passed the sign-in to Claude Code on this machine.` છાપે છે.

લૉગિન ટર્મિનલ પર input ની રાહ જુએ ત્યારે, ccam ટર્મિનલનો echo અને line mode બંધ કરે છે અને તમે ટાઇપ કરો તે દરેક અક્ષર માટે એક `*` બતાવે છે, જેથી paste કરેલો code કે callback સરનામું સ્ક્રીન પર ન દેખાય. `readLine` line editing જાતે કરે છે: `Enter` લાઇન મોકલે છે, `Backspace` એક અક્ષર ભૂંસે છે, અને ખાલી લાઇન પર `Ctrl+D` input પૂરું કરે છે. લૉગિન પૂરું થાય ત્યારે ટર્મિનલ પહેલાં જેવું થઈ જાય છે, અને `Ctrl+C`, `SIGTERM` કે `SIGHUP` પર પણ, જે હજી પણ ccam ને બંધ કરે છે. standard input અને standard error પર ટર્મિનલ ન હોય ત્યારે, લાઇનો જેમ આવે તેમ કૉપિ થાય છે.

### નામ બદલવું, દૂર કરવું અને ઇમ્પોર્ટ કરવું

`ccam rename` અને `ccam remove` એકાઉન્ટના ફોલ્ડર પર અને ccam બીજે એકાઉન્ટના નામે જે રાખે છે તેના પર કામ કરે છે. `nameState` એ પાથની યાદી આપે છે: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` અને વપરાશ કૅશ `~/.local/state/ccam/usage/<name>.json` તથા `<name>.statusline.json`. `ccam remove work` તેમને આ ક્રમમાં સંભાળે છે:

<!-- motion:accounts-remove -->

`ccam remove` પહેલાં એકાઉન્ટ ચાલતું હોય તો ઇનકાર કરે છે, પછી પૂછે છે. જો એકાઉન્ટ એક્સટેન્શન્સ શેર કરતું હોય, તો `RewritePluginPaths` ગ્રુપની `installed_plugins.json` અને `known_marketplaces.json` ફરીથી લખે છે, કારણ કે Claude Code plugin ના absolute પાથ સાચવે છે જે આ એકાઉન્ટના ફોલ્ડરમાંથી પસાર થતા હોઈ શકે. `profile.Trash` ફોલ્ડરને `~/.local/share/ccam/trash/<name>-<time>/` માં ખસેડે છે, એકાઉન્ટ `config.toml` માંથી નીકળી જાય છે, `nameState` પાથ કાઢી નાખવામાં આવે છે જેથી પછીથી એ જ નામનું એકાઉન્ટ તેમને વારસામાં ન મેળવે, અને શૉર્ટકટની લિંક દૂર થાય છે. પાછું મેળવવા માટે, ફોલ્ડરને trash માંથી પાછું ખસેડો.

`ccam rename` ફોલ્ડરનું નામ બદલે છે, તેના `plugins/` માંના plugin પાથ જૂના ફોલ્ડરથી નવા ફોલ્ડર પર ફરીથી લખે છે, ડિફૉલ્ટને સાથે લઈ જાય છે, અને `nameState` પાથ નવા નામ પર ખસેડે છે. શૉર્ટકટ કમાન્ડ પોતાનું નામ રાખે છે. જો plugin પાથ ફરીથી લખવામાં નિષ્ફળ જાય, તો ફોલ્ડરનું નામ પાછું બદલાય છે, કારણ કે config સેવ થતું નથી.

`ccam import` ફોલ્ડરને જ્યાં છે ત્યાં જ, `config.toml` માં `dir` તરીકે રજિસ્ટર કરે છે. તે તેને ક્યારેય ખસેડતું કે કૉપિ કરતું નથી: કૉપિ થયેલો refresh token નકલોને એકબીજાને લૉગઆઉટ કરાવી દે. ફોલ્ડર આર્ગ્યુમેન્ટ વગર તે રજિસ્ટર ન થયેલાં `~/.claude` અને `~/.claude-*` ફોલ્ડર શોધે છે અને દરેક વિશે પૂછે છે; નામ ફોલ્ડર પરથી આવે છે (`~/.claude` એ `main`, `~/.claude-work` એ `work`, `~/.claude-2` એ `acct2`), અને જે ફોલ્ડરનું નામ કાઢી ન શકાય તેને વાપરવાના `--name` કમાન્ડ સાથે બતાવવામાં આવે છે. `CLAUDE_CONFIG_DIR` વગર, Claude Code `.claude.json` ને તમારા home ફોલ્ડરમાં `~/.claude` ની બાજુમાં રાખે છે, તેથી જ્યારે તમે `~/.claude` ઇમ્પોર્ટ કરો અને તેની પોતાની `.claude.json` ન હોય, ત્યારે ccam `~/.claude/.claude.json` થી `~/.claude.json` ની લિંક બનાવે છે. પછી સાદું `claude` અને `ccam main` એક જ એકાઉન્ટ જુએ છે. તમે `ccam share` ન ચલાવો ત્યાં સુધી ઇમ્પોર્ટ કરેલાં એકાઉન્ટ કંઈ શેર કરતાં નથી, અને કોઈ શૉર્ટકટ બનતો નથી.

`ccam completion bash|zsh` એક completion સ્ક્રિપ્ટ છાપે છે, જે એકાઉન્ટનાં નામ પણ પૂરાં કરે છે.

## સલામતી અને ખાસ કિસ્સા

- એકાઉન્ટ ચાલતું હોય ત્યારે `ccam rename` અને `ccam remove` (`refuseRunning`) `work is running; quit its Claude Code session first` સાથે ઇનકાર કરે છે, કારણ કે સેશન ફોલ્ડરમાં લખવાનું ચાલુ રાખે છે.
- ઇમ્પોર્ટ કરેલા એકાઉન્ટ પર `ccam remove` ફક્ત તેનું રજિસ્ટ્રેશન રદ કરે છે: `✓ main unregistered (folder ... left in place)`. તમારું `~/.claude` trash માં ખસેડાતું નથી.
- જેનું ફોલ્ડર `profiles/` માં પહેલેથી હોય તેવા નામનો `ccam add` ઇનકાર કરે છે, અને `ccam import` તરફ નિર્દેશ કરે છે.
- `ccam import` પહેલેથી રજિસ્ટર થયેલા ફોલ્ડરનો, ફોલ્ડર ન હોય તેવા પાથનો અને વપરાયેલા નામનો ઇનકાર કરે છે.
- `ccam logout` અને `ccam login` એકાઉન્ટ ચાલી રહ્યું છે કે નહીં તે તપાસતા નથી.
- જો relay Claude Code સુધી ન પહોંચી શકે, ઉદાહરણ તરીકે લૉગિન પૂરું થઈ ગયા પછી, તો તે કારણ સાથે `ccam: could not pass the sign-in to Claude Code:` છાપે છે. HTTP ભૂલનો અર્થ સામાન્ય રીતે એ છે કે કોડની મુદત પૂરી થઈ ગઈ; લૉગિન ફરીથી શરૂ કરો.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [શૉર્ટકટ કમાન્ડ](/features/shortcuts)
- [MCP સર્વર, સેટિંગ્સ અને ઇતિહાસ શેર કરો](/docs/sharing)
- [ઘણાં Claude Code એકાઉન્ટ](/guides/multiple-claude-code-accounts)
