---
title: ફોલ્ડર પ્રમાણે એકાઉન્ટ
description: ccam folder એકાઉન્ટના નામ વગર Claude Code શરૂ કરો ત્યારે ફોલ્ડર અને તેની નીચેની દરેક વસ્તુ માટે એક એકાઉન્ટ વાપરે છે, જેમ કે કંપનીના કોડ માટે work.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` એકાઉન્ટના નામ વગરના દરેક લૉન્ચને `~/code/company` ની અંદર, અને તેની નીચેના કોઈ પણ ફોલ્ડરમાં, Claude Code ને `work` તરીકે શરૂ કરતું કરે છે. તેની બહાર, ડિફૉલ્ટ એકાઉન્ટ પહેલાંની જેમ શરૂ થાય છે. એકાઉન્ટનું નામ આપવાથી, `ccam personal` વડે કે `ccpersonal` જેવા શૉર્ટકટ વડે, હંમેશાં એ જ એકાઉન્ટ શરૂ થાય છે.

## ઉપયોગ

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

આ લૉન્ચ ફોલ્ડરનું એકાઉન્ટ વાપરે છે:

- `ccam -p "…"` અને બીજું કોઈ પણ લૉન્ચ જેનો પહેલો આર્ગ્યુમેન્ટ Claude Code નો ઓપ્શન હોય
- નામ વગરનું `ccam run`, અને `ccam --dry-run`
- `ccam --fallback -p "…"` તેના પહેલા પ્રયાસ માટે
- picker, એકલું `ccam`, જે ફોલ્ડરના એકાઉન્ટ પર cursor સાથે શરૂ થાય છે

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `<account>` | વાપરવાના એકાઉન્ટનું નામ કે નંબર | |
| `[folder]` | ફોલ્ડર, relative કે absolute | હાલનું ફોલ્ડર |
| `--remove` | બરાબર એ જ ફોલ્ડરની entry દૂર કરે છે | off |
| `--dry-run` | માન્ય: ફેરફાર સેવ કર્યા વગર છાપે છે | off |

## તે કેવી રીતે કામ કરે છે

દરેક ફોલ્ડર `~/.config/ccam/config.toml` માં એક `[[folder]]` entry છે, જે ફાઇલ દરેક લૉન્ચ આમ પણ વાંચે છે:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` ફોલ્ડરને symbolic links ઉકેલીને absolute પાથમાં ફેરવે છે અને સેવ કરે છે. કોઈ ફોલ્ડર ફરી સેટ કરવાથી તેનું એકાઉન્ટ બદલાય છે.
2. નામ વગરનું લૉન્ચ `launcher.HereProfile` ને બોલાવે છે. કોઈ entry ન હોય તો તે તરત ડિફૉલ્ટ એકાઉન્ટ લે છે.
3. નહીં તો તે working directory ને દરેક `path` સાથે strings તરીકે સરખાવે છે: ફોલ્ડર પોતે કે તેની નીચેનું ફોલ્ડર, ફક્ત એ જ અક્ષરોથી શરૂ થતું ફોલ્ડર ક્યારેય નહીં (`/code/company` એ `/code/companyx` ને આવરી લેતું નથી). સૌથી ઊંડી entry જીતે છે, તેથી `~/code` `personal` હોઈ શકે જ્યારે `~/code/company` `work` હોય.
4. working directory જેમ છે તેમ કોઈ સાથે મેળ ન ખાય ત્યારે જ ccam તેની symbolic links ઉકેલે છે અને ફરી એક વાર સરખાવે છે.
5. Windows પર સરખામણી, Windows ની જેમ જ, નાના-મોટા અક્ષરોને અવગણે છે.

`config.toml` સિવાય કોઈ ફાઇલ વંચાતી નથી, તેથી ફોલ્ડરની અંદરના લૉન્ચમાં માપી શકાય એવો કોઈ સમય લાગતો નથી (entries સાથે કે વગર 4.0 ms, `hyperfine`, 300 runs).

## સલામતી અને ખાસ કિસ્સા

- `ccam rename` ફોલ્ડર entries ને નવા નામ પર લઈ જાય છે, અને `ccam remove` દૂર કરેલા એકાઉન્ટની entries કાઢી નાખે છે.
- જે entry નું એકાઉન્ટ હવે અસ્તિત્વમાં ન હોય તે એ ફોલ્ડરમાં નામ વગરના લૉન્ચને અટકાવે છે: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` આવી entries, અને જેમનું ફોલ્ડર ન રહ્યું હોય તેવી entries, બતાવે છે.
- ઉપરના ફોલ્ડરમાંથી એકાઉન્ટ મેળવતા ફોલ્ડરની અંદર `ccam folder --remove` કંઈ દૂર કરવાને બદલે એ ફોલ્ડરનું નામ બતાવે છે.
- `ccam next` rotation કે વપરાશ પ્રમાણે પસંદ કરે છે અને ફોલ્ડર જોતું નથી.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [શૉર્ટકટ કમાન્ડ](/features/shortcuts)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [કમાન્ડ સંદર્ભ](/docs/commands)
