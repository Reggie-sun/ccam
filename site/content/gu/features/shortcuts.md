---
title: શૉર્ટકટ કમાન્ડ
description: ccam shortcut વડે કોઈ એકાઉન્ટને ccwork જેવો તેનો પોતાનો કમાન્ડ આપો, અને ccam shims install વડે ~/.local/bin માંની લિંક ફરીથી બનાવો.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

શૉર્ટકટ એ તમે જાતે નામ આપેલો કમાન્ડ છે, જે Claude Code ને એક એકાઉન્ટ તરીકે શરૂ કરે છે: `ccwork --resume` એ જ કરે છે જે `ccam work --resume` કરે છે. તે `~/.local/bin` માં ccam ની લિંક છે, તેથી તે bash, zsh, સ્ક્રિપ્ટ અને `PATH` પરથી પ્રોગ્રામ ચલાવતી બીજી કોઈ પણ જગ્યાએ કામ કરે છે. શૉર્ટકટ વૈકલ્પિક છે; તમે વારંવાર શરૂ કરતા હો તે એકાઉન્ટ માટે એક બનાવો.

## ઉપયોગ

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam દરેક ફેરફારની પુષ્ટિ કરે છે:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

એક એકાઉન્ટ પાસે વધુમાં વધુ એક શૉર્ટકટ હોય છે, જે `~/.config/ccam/config.toml` માં `shortcut` તરીકે સેવ થાય છે. `ccam add` વિઝાર્ડ `Create a shortcut command for <name>? [y/N]` અને પછી `Shortcut name [cc<name>]:` પૂછે છે. વિઝાર્ડની બહાર, `ccam add` ફક્ત `--shortcut` સાથે જ શૉર્ટકટ બનાવે છે. `ccam list` દરેક એકાઉન્ટનો શૉર્ટકટ બતાવે છે, અથવા `-`.

શૉર્ટકટ દરેક આર્ગ્યુમેન્ટ Claude Code ને આપે છે, બરાબર જેમ `ccam <name>` નામ પછી આપે છે. તેમાં `--dry-run` અને `--help` પણ સામેલ છે: શરૂઆતનું પૂર્વાવલોકન કરવા માટે `ccam --dry-run work` વાપરો.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--remove` | એકાઉન્ટનો શૉર્ટકટ અને તેની લિંક કાઢી નાખે છે | off |
| `-h`, `--help` | કમાન્ડની મદદ | |
| `--dry-run` | `ccam shortcut` તેને માને છે: તે નામ તપાસે છે અને સેવ કે લિંક કર્યા વગર `Would set the shortcut of work to ccwork.` કે `Would remove ccwork, the shortcut of work.` છાપે છે | off |

`ccam shims install` કોઈ આર્ગ્યુમેન્ટ લેતું નથી અને `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` સાથે `--dry-run` નો ઇનકાર કરે છે.

નામના નિયમો, આ ક્રમમાં તપાસાય છે:

- `a-z`, `0-9`, `.`, `-` અને `_` ના 1 થી 64 અક્ષર, અક્ષર કે અંકથી શરૂ થતા;
- `ccam`, `claude` કે `run` અથવા `usage` જેવો ccam કમાન્ડ નહીં;
- બીજા એકાઉન્ટનો શૉર્ટકટ નહીં;
- `~/.local/bin` માં એવી ફાઇલ નહીં જે પહેલેથી ccam ની લિંક ન હોય;
- તમારા `PATH` પરના બીજા executable નું નામ નહીં.

કમાન્ડ અને `--remove` બંને સાથે આપવાથી `give a command or --remove, not both` સાથે અટકે છે. એકાઉન્ટ પાસે પહેલેથી હોય તે જ નામ સેટ કરવાથી `work: ccwork (unchanged)` છપાય છે, અને લિંક ખૂટતી હોય તો ફરીથી બને છે.

## તે કેવી રીતે કામ કરે છે

નામ સેટ કરવાથી `checkShortcut` ચાલે છે, જે bin ફોલ્ડર અને ccam binary ના સાચા પાથ સાથે `shim.Check` ને બોલાવે છે:

<!-- motion:shortcuts-check -->

નામ તપાસોમાંથી ક્રમમાં નીચે જાય છે. દરેક તપાસની બાજુમાં એક નામ છે જેનો તે ઇનકાર કરે છે: `Work` pattern માં નિષ્ફળ જાય છે, `run` ccam કમાન્ડ છે, `ccside` બીજા એકાઉન્ટનો છે, `ccfile` `~/.local/bin` માં હાલની ફાઇલ છે, અને `git` `/usr/bin/git` ને છુપાવી દે. ફક્ત દરેક તપાસ પાર કરતું નામ જ સેવ અને લિંક થાય છે.

1. `shim.Check` નામને `^[a-z0-9][a-z0-9._-]{0,63}$` સાથે, પછી `config.Reserved` સાથે, પછી `cfg.FindShortcut` સાથે મેળવે છે. તે `~/.local/bin/<name>` પર `os.Lstat` ને બોલાવે છે, અને `PATH` પરના દરેક બીજા ફોલ્ડર માટે `os.Stat` ને બોલાવે છે અને ccam ની લિંક ન હોય તેવી executable ફાઇલનો ઇનકાર કરે છે.
2. ccam `config.toml` ને `fsutil.AtomicWrite` (કામચલાઉ ફાઇલ, `fsync`, rename) વડે સેવ કરે છે.
3. `shim.Remove` જૂના શૉર્ટકટની લિંક કાઢી નાખે છે, પણ ફક્ત ત્યારે જ જ્યારે તે ccam ની લિંક હોય.
4. `shim.Install` ગોઠવેલા દરેક શૉર્ટકટની ખૂટતી લિંક `os.Symlink` વડે ccam binary પર બનાવે છે, જે `os.Executable` અને `filepath.EvalSymlinks` વડે શોધાય છે. તે ક્યારેય કંઈ દૂર કરતું નથી.

જ્યારે `os.Readlink` ccam નો પાથ આપે, અથવા લિંક અને ccam નો પાથ એક જ ફાઇલ પર પહોંચે, ત્યારે લિંક ccam ની પોતાની ગણાય છે. દરેક remove અને દરેક overwrite ની તપાસ આ કસોટીમાંથી પસાર થાય છે, તેથી ccam પોતે ન બનાવી હોય તેવી ફાઇલોને અડતું નથી.

જ્યારે તમે શૉર્ટકટ ચલાવો, ત્યારે તમે લખેલું નામ જ ccam ને જરૂરી એકમાત્ર વસ્તુ છે:

<!-- motion:shortcuts-link -->

shell `PATH` પર `ccwork` શોધે છે, kernel લિંક અનુસરીને ccam binary શરૂ કરે છે, અને `argv[0]` હજી પણ `ccwork` હોય છે. `app.Main` માં, `ccam` સિવાયનું કોઈ પણ નામ `config.toml` એક વાર લોડ કરે છે અને `shim.Lookup` ને બોલાવે છે, જે એ એકાઉન્ટ શોધે છે જેનો `shortcut` `filepath.Base(argv[0])` જેટલો હોય. પછી `launcher.LaunchWith` પહેલેથી લોડ થયેલા config સાથે [`ccam <name>`](/features/launch) જેવાં જ પગલાં લે છે, અને `argv[1:]` ને કોઈ ફેરફાર વગર `claude` ને આપે છે.

shell alias કે wrapper સ્ક્રિપ્ટ નહીં પણ symlink શા માટે: symlink દરેક shell માં અને સ્ક્રિપ્ટમાં કામ કરે છે, તે shell ની વધારાની શરૂઆત (આશરે 1 થી 2 ms) ઉમેરતી નથી, અને શરૂઆત હજી એક જ ફાઇલ વાંચે છે. `scripts/perf-check.sh` `hyperfine` વડે શૉર્ટકટને `ccam <name>` જેવી જ 10 ms (p95) મર્યાદા સામે માપે છે.

## સલામતી અને ખાસ કિસ્સા

- **કોઈ એકાઉન્ટની ન હોય તેવી લિંક.** જો `~/.local/bin` માંની ccam ની કોઈ લિંકનું નામ કોઈ એકાઉન્ટ ન વાપરતું હોય, ઉદાહરણ તરીકે તમે `config.toml` માં `shortcut` જાતે બદલ્યા પછી, તો ccam ડિફૉલ્ટ એકાઉન્ટ પર જતું નથી. તે `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` સાથે અટકે છે. `ccam shortcut` વડે નામ ફરી સેટ કરો, અથવા લિંક કાઢી નાખો.
- **હાલની ફાઇલો.** ccam પોતે ન બનાવી હોય તેવી ફાઇલ ક્યારેય બદલતું નથી. `ccam shims install` `<path> already exists and is not a ccam link, so it was left alone` સાથે અટકે છે.
- **ખસેડેલી ccam binary.** ત્યારે જૂની લિંક એવા પાથ તરફ નિર્દેશ કરે છે જે હવે નથી, અને ccam તેમને પોતાની ગણતું નથી. તેમને કાઢી નાખો, પછી `ccam shims install` ચલાવો.
- **નામ બદલવું અને દૂર કરવું.** `ccam rename` શૉર્ટકટનું નામ રાખે છે, અને લિંક કામ કરતી રહે છે, કારણ કે તે ccam તરફ નિર્દેશ કરે છે અને ccam `config.toml` માં નામ શોધે છે. `ccam remove` એકાઉન્ટની લિંક કાઢી નાખે છે. `ccam shortcut <account> <new>` નવી લિંક બનાવતાં પહેલાં જૂની લિંક કાઢી નાખે છે.
- **PATH.** જ્યારે ccam લિંક બનાવે અને `~/.local/bin` તમારા `PATH` પર ન હોય, ત્યારે તે `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` અને પછી ઉમેરવાની `export PATH=...` લાઇન છાપે છે.
- **પછીથી તપાસવું.** `ccam doctor` `missing shortcuts in /home/me/.local/bin: ccwork` બતાવે છે અને ફોલ્ડર `PATH` પર ન હોય ત્યારે ચેતવણી આપે છે. `ccam doctor --fix` ખૂટતી લિંક બનાવે છે.

## સંબંધિત

- [Claude Code શરૂ કરવું](/features/launch)
- [એકાઉન્ટનું સંચાલન](/features/accounts)
- [કમાન્ડ સંદર્ભ](/docs/commands#shortcuts)
- [ccam ઇન્સ્ટૉલ કરો](/docs/install)
