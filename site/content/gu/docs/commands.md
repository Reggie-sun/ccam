---
title: કમાન્ડ સંદર્ભ
metaTitle: ccam કમાન્ડ સંદર્ભ
description: ccam ના બધા કમાન્ડ અને ઓપ્શન, ccam add અને ccam <name> થી લઈને ccam usage, ccam next અને ccam doctor સુધી.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: 338ac11a0b5a
---

કોઈ કમાન્ડ ખરેખર ચલાવ્યા વગર તે શું કરશે તે જોવા માટે તેની આગળ `--dry-run` મૂકો, ઉદાહરણ તરીકે `ccam --dry-run work`. તે Claude Code શરૂ કરતી વખતે (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) અને `group`, `migrate`, `shortcut`, `update` તથા `doctor --fix` સાથે કામ કરે છે. કંઈક બદલતા બીજા કમાન્ડ ચાલવાને બદલે તેનો ઇનકાર કરે છે.

<!-- commands -->

## Claude Code શરૂ કરવું

`ccam <name|number> [arguments]` તે એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે. ccam નામ પછીના આર્ગ્યુમેન્ટ્સ વાંચતું નથી; તે તેમને એ જ ક્રમમાં `claude` ને આપે છે, અને standard input, exit code અને સિગ્નલ બરાબર `claude` જેવું જ વર્તન કરે છે. નામ પછી `--` મૂકવું વૈકલ્પિક છે.

આર્ગ્યુમેન્ટ્સ સાથે પણ એકાઉન્ટના નામ વગરનું `ccam`, જેમ કે `ccam -p "summarise"`, હાલના ફોલ્ડર માટે સેટ કરેલું એકાઉન્ટ, નહીં તો ડિફૉલ્ટ એકાઉન્ટ વાપરે છે. તેમને `ccam folder <name> [folder]` અને `ccam use <name>` વડે સેટ કરો.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

એકાઉન્ટના નામ પહેલાં `--fallback` એ `-p`/`--print` runs માટે છે. એકાઉન્ટ તેની વપરાશ મર્યાદાએ પહોંચે ત્યારે ccam એ જ command સૌથી વધુ જગ્યાવાળા લૉગિન થયેલા એકાઉન્ટ તરીકે ફરી ચલાવે છે અને ફક્ત પૂરા થયેલા run નું output છાપે છે. જુઓ [ફોલ્ડર પ્રમાણે એકાઉન્ટ](/features/folder) અને [મર્યાદાએ ફરી ચલાવવું](/features/next#run-again-at-the-limit).

`ccwork` જેવો શૉર્ટકટ તેના એકાઉન્ટને આ જ રીતે શરૂ કરે છે: `ccwork --resume` એ જ કરે છે જે `ccam work --resume` કરે છે. [શૉર્ટકટ](/docs/commands#shortcuts) જુઓ.

## એકાઉન્ટ ઉમેરવાં અને દૂર કરવાં

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

ફોલ્ડર વગરનું `ccam import` `~/.claude` અને `~/.claude-*` શોધે છે અને દરેક વિશે પૂછે છે. ઇમ્પોર્ટ કરેલાં ફોલ્ડર જ્યાં છે ત્યાં જ રજિસ્ટર થાય છે; ccam ક્યારેય લૉગિન ખસેડતું કે કૉપિ કરતું નથી.

`ccam add --api` એવું એકાઉન્ટ ઉમેરે છે જે Claude લૉગિનને બદલે API key વાપરે છે; ccam દરેક લૉન્ચ વખતે key વાંચે છે અને તેને ક્યારેય સેવ કરતું નથી. [API key એકાઉન્ટ](/features/api) જુઓ.

Claude Code કોઈ એકાઉન્ટ તરીકે ચાલી રહ્યું હોય ત્યારે `ccam rename`, `ccam remove` અને `ccam group` કમાન્ડ તે એકાઉન્ટમાં ફેરફાર કરવાનો ઇનકાર કરે છે. પહેલાં સેશન બંધ કરો.

`ccam list --json` દરેક એકાઉન્ટને તેના સ્લૉટ, શૉર્ટકટ, ઇમેઇલ અને ફોલ્ડર સાથે છાપે છે. તેનો `groups` array એકાઉન્ટનાં ગ્રુપની યાદી આપે છે, અને તેનો `share` ઑબ્જેક્ટ દરેક આઇટમ (`mcp`, `settings`, `memory`, `extensions`, `history`) માટે તે જે ગ્રુપમાંથી આવે છે તે, અથવા `local` આપે છે.

## શૉર્ટકટ

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

શૉર્ટકટ એ તમે જાતે નામ આપેલો કમાન્ડ છે, જે Claude Code ને એક એકાઉન્ટ તરીકે શરૂ કરે છે. `ccam work` ની જેમ તે દરેક આર્ગ્યુમેન્ટ કોઈ ફેરફાર વગર Claude Code ને આપે છે, તેથી `ccwork --resume` અને `ccwork -p "hi"` અપેક્ષા મુજબ કામ કરે છે. તે `~/.local/bin` માં ccam ની લિંક છે.

શૉર્ટકટ વૈકલ્પિક છે. `ccam add` વિઝાર્ડ શૉર્ટકટ બનાવવો કે નહીં તે પૂછે છે, `ccam add <name> --shortcut <command>` પૂછ્યા વગર એક બનાવે છે, અને `ccam list` દરેક એકાઉન્ટનો શૉર્ટકટ બતાવે છે, અથવા શૉર્ટકટ ન હોય ત્યારે `-`.

નામમાં a-z, 0-9, `.`, `-` અને `_` અક્ષરો વાપરી શકાય છે. નામ કોઈ ccam કમાન્ડ, બીજા એકાઉન્ટનો શૉર્ટકટ, કે તમારા `PATH` પરના બીજા પ્રોગ્રામનું નામ ન હોઈ શકે.

`ccam rename` શૉર્ટકટનું નામ એ જ રાખે છે, અને `ccam remove` એકાઉન્ટની સાથે તેનો શૉર્ટકટ પણ કાઢી નાખે છે. `ccam shims install` બધાં એકાઉન્ટની લિંક ફરીથી બનાવે છે, ઉદાહરણ તરીકે ccam ને બીજે ખસેડ્યા પછી.

## શેરિંગ

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ગ્રુપ તેનાં એકાઉન્ટ વચ્ચે આઇટમ શેર કરે છે; એક એકાઉન્ટ ઘણાં ગ્રુપમાં હોઈ શકે છે, પણ દરેક આઇટમ ફક્ત એક જ ગ્રુપમાંથી આવે છે. [ગ્રુપ વડે શેરિંગ](/features/share) જુઓ.

## હાલનું સેટઅપ કૉપિ કરવું

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` કોઈ Claude Code ફોલ્ડરમાંથી `settings.json`, `CLAUDE.md`, એક્સટેન્શન્સ અને user MCP સર્વર એક ગ્રુપમાં કૉપિ કરે છે, જેથી તે આઇટમ શેર કરતા દરેક એકાઉન્ટને તે મળે. ફોલ્ડર ફક્ત વંચાય છે. `--settings`, `--memory`, `--extensions`, `--mcp` અને `--history` ઓપ્શન્સ ફક્ત તમે નામ આપેલી આઇટમ કૉપિ કરે છે. તેમના વગર, ટર્મિનલ પર એક ચેકલિસ્ટ તમને પસંદ કરવા દે છે, જેમાં વાતચીતના ઇતિહાસ સિવાય બધું ચેક કરેલું હોય છે; બીજે બધે વાતચીતના ઇતિહાસ સિવાય બધું કૉપિ થાય છે. `-g` કે `--group` ગ્રુપ પસંદ કરે છે, અને `-y` કે `--yes` પૂછ્યા વગર કૉપિ કરે છે અને કોઈ ફાઇલ જુદી હોય ત્યારે ગ્રુપની કૉપિ રાખે છે. તે સર્વરનાં સેવ કરેલાં સાઇન-ઇન પણ દરેક એકાઉન્ટની પોતાની લૉગિન ફાઇલમાં કૉપિ કરે છે. [હાલનું સેટઅપ કૉપિ કરવું](/docs/sharing#copy-an-existing-setup) જુઓ.

## MCP સર્વર

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` ફક્ત તે એકાઉન્ટ પાસે જ હોય તેવા સર્વર તેના ગ્રુપમાં ખસેડે છે. જે સર્વરનું નામ ગ્રુપ પહેલેથી બીજી વ્યાખ્યા માટે વાપરે છે તેને છોડી દેવામાં આવે છે; ગ્રુપની કૉપિ બદલવા માટે `--force` ઉમેરો, ccam પહેલાં તેનો બૅકઅપ લે છે.

## વપરાશ અને સ્વિચિંગ

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` અને `ccam pick` એક જ રીતે પસંદ કરે છે: પહેલાં `auto_max_usage` (`config.toml`, ડિફૉલ્ટ 100) થી નીચેનાં, `auto_strategy` (`headroom`, `expiring` કે `balance`) પ્રમાણે, `--model` ની સાપ્તાહિક વિન્ડો ગણીને, 100% પરનું એકાઉન્ટ ક્યારેય નહીં, અને જે એકાઉન્ટથી તેઓ શરૂ કરે તેની સંસ્થાની અંદર (વધુ આગળ જોવા માટે `--cross-org` કે `cross_org = true`). કોઈ પસંદ ન થઈ શકે ત્યારે તેઓ code 75 સાથે બહાર નીકળે છે, અથવા `--wait` સાથે પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે. `ccam next --continue` ઇતિહાસ શેર ન કરતાં એકાઉન્ટ વચ્ચે પણ વાતચીત આગળ લઈ જાય છે: ત્યારે Claude Code તેની નકલ નવા એકાઉન્ટમાં કરે છે. `ccam <name>` થી અલગ, `ccam next` પોતાના ઓપ્શન જાતે વાંચે છે, તેથી Claude Code ના ઓપ્શન `--` પછી આવે છે; તેના વગર ccam `unknown flag` સાથે અટકી જાય છે. Claude Code ની અંદર, plugin નું `/ccam:switch [account]` વાતચીત આગળ સોંપે છે, અને તેના hooks `ccam hook` ચલાવે છે. જુઓ [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next), [scripts માટે એકાઉન્ટ પસંદ કરવું](/features/pick) અને [દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન](/features/agents).

## જાળવણી

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

ccam ને બદલતાં પહેલાં `ccam update` રિલીઝનો SHA-256 checksum તપાસે છે, અને `--version v0.2.0` કોઈ ચોક્કસ રિલીઝ ઇન્સ્ટૉલ કરે છે. `.deb` પેકેજથી ઇન્સ્ટૉલ કરેલું ccam તેના બદલે apt વડે અપડેટ થાય છે. [અપડેટ](/docs/install#update) જુઓ.

## ફાઇલો

| પાથ | સામગ્રી |
| --- | --- |
| `~/.config/ccam/config.toml` | એકાઉન્ટ, ગ્રુપ અને સેટિંગ્સ (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | દરેક એકાઉન્ટનું `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | શેર કરેલાં સેટિંગ્સ, MCP સર્વર, એક્સટેન્શન્સ અને ઇતિહાસ |
| `~/.local/state/ccam/` | વપરાશ કૅશ, છેલ્લે વાપર્યાના સમય, અને `/ccam:switch` તથા મર્યાદા hook છોડી જાય તે નોંધો (`handoff/`) |
| `~/.local/bin/<shortcut>` | `ccwork` જેવા શૉર્ટકટ કમાન્ડ, ccam ની લિંક |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` અને `XDG_STATE_HOME` વેરિએબલ આ ફોલ્ડરનું સ્થાન બદલે છે. Windows પર `~` એટલે `%USERPROFILE%`, અને શૉર્ટકટ `<shortcut>.exe` હોય છે.
