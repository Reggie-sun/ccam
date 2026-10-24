---
title: દરેક એકાઉન્ટનાં બૅકગ્રાઉન્ડ સેશન
description: ccam agents તમારાં બધાં એકાઉન્ટના Claude Code બૅકગ્રાઉન્ડ agents અને ખુલ્લાં સેશન એક જ table માં બતાવે છે, અને ccam attach તેના id વડે એક ખોલે છે.
updated: 2026-10-23
order: 16
translation: machine
sourceHash: 092368f63363
---

Claude Code બૅકગ્રાઉન્ડ સેશન (`claude --bg`, agent view) એકાઉન્ટ પ્રમાણે રાખે છે: `claude agents` ફક્ત એ જે એકાઉન્ટ તરીકે ચાલે છે તેનાં સેશન બતાવે છે. `ccam agents` બધાં એકાઉન્ટને એકસાથે પૂછે છે અને તેમને સાથે બતાવે છે, અને `ccam attach <id|name>` બૅકગ્રાઉન્ડ સેશનને એ જે એકાઉન્ટે શરૂ કર્યું હોય તે એકાઉન્ટ તરીકે ખોલે છે.

## ઉપયોગ

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` બૅકગ્રાઉન્ડ સેશનની સ્થિતિ (`working`, `blocked`, `done`, `failed`, `stopped`) અને તેની process શું કરી રહી છે (`busy`, `waiting`, `idle`) તે બતાવે છે; ખુલ્લું interactive સેશન `open` બતાવે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--all` | પૂરાં થયેલાં, નિષ્ફળ ગયેલાં કે બંધ થયેલાં બૅકગ્રાઉન્ડ સેશન પણ બતાવે છે | off |
| `--json` | સેશનને JSON તરીકે છાપે છે, દરેક તે જે `account` નું છે તેની સાથે | off |

`ccam attach <id|name> [claude args]` id પછીનું બધું, flags સહિત, કોઈ ફેરફાર વગર `claude attach` ને આપે છે: `ccam attach 49ec --dry-run` `--dry-run` ને Claude Code ને સોંપે છે. તેના બદલે કમાન્ડ છાપવા માટે, `--dry-run` ને id ની પહેલાં મૂકો (`ccam --dry-run attach 49ec` કે `ccam attach --dry-run 49ec`).

## તે કેવી રીતે કામ કરે છે

1. દરેક એકાઉન્ટ માટે, ccam એકાઉન્ટના `CLAUDE_CONFIG_DIR` સાથે `claude agents --json` (આપ્યું હોય ત્યારે `--all` સાથે) ચલાવે છે, બધાં એકસાથે, દરેક 15 સેકન્ડની અંદર. આ એકાઉન્ટની પોતાની સેશન યાદી વાંચે છે અને Claude Code ની બૅકગ્રાઉન્ડ સેવા શરૂ કરતું નથી.
2. દરેક એન્ટ્રીને એકાઉન્ટનું નામ મળે છે, અને યાદી શરૂઆતના સમય પ્રમાણે, સૌથી નવું પહેલાં, ગોઠવાય છે. જવાબ ન આપતા એકાઉન્ટનું નામ standard error પર આવે છે અને તે બહાર રહે છે, અને પછી ccam code 1 સાથે બહાર નીકળે છે, જેથી કોઈ script નિષ્ફળતાને ખાલી યાદી ન સમજી બેસે.
3. `ccam attach` દરેક એકાઉન્ટના ફોલ્ડરમાં `jobs/<id>` શોધે છે, જ્યાં Claude Code બૅકગ્રાઉન્ડ સેશન રાખે છે, અને `ccam <name>` જે રીતે Claude Code શરૂ કરે છે તે જ રીતે એ એકાઉન્ટ તરીકે `claude attach <id>` ચલાવે છે. જે શબ્દ કોઈ id ની શરૂઆત ન હોય તેને સેશનનાં નામોમાં શોધવામાં આવે છે, અક્ષરો મોટા હોય કે નાના; id જેવો દેખાતો શબ્દ ક્યારેય નામ તરીકે લેવાતો નથી.

## સલામતી અને ખાસ કિસ્સા

- બૅકગ્રાઉન્ડ સેશનના id એક જ એકાઉન્ટના હોય છે: બીજા એકાઉન્ટ તરીકે ચલાવેલું `claude attach` તેમને શોધી શકતું નથી, અને `ccam attach` આ જ ઉકેલે છે.
- બે એકાઉન્ટ પાસે એવાં સેશન હોઈ શકે જેમના id એક જ રીતે શરૂ થાય, કે જેમનાં નામમાં એ જ શબ્દ હોય; ત્યારે `ccam attach` તેમનાં નામ આપે છે અને id નો કે નામનો વધુ ભાગ માગે છે.
- ccam તેના શેર કરેલા MCP સર્વર Claude Code ને એક જ આર્ગ્યુમેન્ટમાં `--mcp-config=<file>` વડે આપે છે, જે `claude agents` અને `claude attach` સ્વીકારે છે; બૅકગ્રાઉન્ડ સેશન ફરી શરૂ થાય ત્યારે પણ એ જ ફાઇલ વાપરતાં રહે છે.
- API key એકાઉન્ટની બૅકગ્રાઉન્ડ સેવા તેની key તેને શરૂ કરનાર શેલ પાસેથી લે છે. એવા એકાઉન્ટનું પહેલું બૅકગ્રાઉન્ડ સેશન ccam દ્વારા શરૂ કરો (`ccam gw --bg …`), જેથી સેવા પાસે key હોય.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [scripts માટે એકાઉન્ટ પસંદ કરવું](/features/pick)
- [ડેશબોર્ડ](/features/dashboard)
