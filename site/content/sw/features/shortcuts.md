---
title: Amri za njia fupi
description: Ipe akaunti amri yake yenyewe kama ccwork kwa ccam shortcut, na uunde upya viungo katika ~/.local/bin kwa ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Njia fupi ni amri unayoipa jina mwenyewe inayoanzisha Claude Code kama akaunti moja: `ccwork --resume` hufanya kile ambacho `ccam work --resume` hufanya. Ni kiungo cha ccam katika `~/.local/bin`, kwa hiyo inafanya kazi katika bash, zsh, scripts na kitu kingine chochote kinachoendesha programu kutoka `PATH`. Njia fupi ni za hiari; unda moja kwa akaunti unazoanzisha mara nyingi.

## Matumizi

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

ccam huthibitisha kila badiliko:

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

Akaunti ina njia fupi moja tu kwa kiwango cha juu, inayohifadhiwa kama `shortcut` katika `~/.config/ccam/config.toml`. Mchawi wa `ccam add` huuliza `Create a shortcut command for <name>? [y/N]` kisha `Shortcut name [cc<name>]:`. Nje ya mchawi, `ccam add` huunda moja tu kwa `--shortcut`. `ccam list` huonyesha njia fupi ya kila akaunti, au `-`.

Njia fupi hupitisha kila hoja kwa Claude Code, kwa njia ile ile ambayo `ccam <name>` hufanya baada ya jina. Hiyo inajumuisha `--dry-run` na `--help`: ili kuona uanzishaji kabla, tumia `ccam --dry-run work`.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--remove` | Hufuta njia fupi ya akaunti na kiungo chake | imezimwa |
| `-h`, `--help` | Msaada wa amri | |
| `--dry-run` | `ccam shortcut` huiheshimu: hukagua jina na huchapisha `Would set the shortcut of work to ccwork.` au `Would remove ccwork, the shortcut of work.` bila kuhifadhi wala kuunganisha | imezimwa |

`ccam shims install` haipokei hoja na hukataa `--dry-run` kwa `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Kanuni za majina, hukaguliwa kwa mpangilio huu:

- herufi 1 hadi 64 za `a-z`, `0-9`, `.`, `-` na `_`, zikianza kwa herufi au tarakimu;
- si `ccam`, `claude` wala amri ya ccam kama `run` au `usage`;
- si njia fupi ya akaunti nyingine;
- si faili katika `~/.local/bin` ambayo tayari si kiungo cha ccam;
- si jina la programu nyingine inayoendeshwa iliyo kwenye `PATH` yako.

Kutoa amri pamoja na `--remove` husimama na `give a command or --remove, not both`. Kuweka jina ambalo akaunti tayari inalo huchapisha `work: ccwork (unchanged)` na huunda upya kiungo chake ikiwa hakipo.

## Jinsi inavyofanya kazi

Kuweka jina huendesha `checkShortcut`, ambayo huita `shim.Check` pamoja na folda ya bin na njia halisi ya binary ya ccam:

<!-- motion:shortcuts-check -->

Jina hupita ukaguzi mmoja baada ya mwingine kwa mpangilio. Kando ya kila ukaguzi kuna jina ambalo unalikataa: `Work` halilingani na mchoro wa majina, `run` ni amri ya ccam, `ccside` ni ya akaunti nyingine, `ccfile` ni faili iliyopo katika `~/.local/bin`, na `git` ingeficha `/usr/bin/git`. Jina linalopita kila ukaguzi pekee ndilo huhifadhiwa na kuunganishwa.

1. `shim.Check` hulinganisha jina na `^[a-z0-9][a-z0-9._-]{0,63}$`, kisha `config.Reserved`, kisha `cfg.FindShortcut`. Huita `os.Lstat` kwa `~/.local/bin/<name>`, na kwa kila folda nyingine kwenye `PATH` huita `os.Stat` na hukataa faili inayoendeshwa ambayo si kiungo cha ccam.
2. ccam huhifadhi `config.toml` kwa `fsutil.AtomicWrite` (faili ya muda, `fsync`, kubadilisha jina).
3. `shim.Remove` hufuta kiungo cha njia fupi ya zamani, lakini tu ikiwa ni kiungo cha ccam.
4. `shim.Install` huunda kiungo cha kila njia fupi iliyosanidiwa ambayo haipo, kwa `os.Symlink` kuelekea binary ya ccam, iliyotafutwa kwa `os.Executable` na `filepath.EvalSymlinks`. Kamwe haiondoi chochote.

Kiungo huhesabiwa kuwa cha ccam yenyewe wakati `os.Readlink` inaporudisha njia ya ccam, au wakati kiungo na njia ya ccam zinaelekea faili ile ile. Kila uondoaji na kila ukaguzi wa kuandika juu hupitia jaribio hili, kwa hiyo ccam haigusi faili ambazo haikuziunda.

Unapoendesha njia fupi, jina ulilocharaza ndilo kitu pekee ambacho ccam inahitaji:

<!-- motion:shortcuts-link -->

Shell hupata `ccwork` kwenye `PATH`, kernel hufuata kiungo na huanzisha binary ya ccam, na `argv[0]` bado ni `ccwork`. Katika `app.Main`, jina lolote lisilo `ccam` hupakia `config.toml` mara moja na huita `shim.Lookup`, ambayo hupata akaunti ambayo `shortcut` yake ni sawa na `filepath.Base(argv[0])`. Kisha `launcher.LaunchWith` huchukua hatua zile zile kama [`ccam <name>`](/features/launch) ikiwa usanidi tayari umepakiwa, na hupitisha `argv[1:]` kwa `claude` bila kuigusa.

Kwa nini symlink na si alias ya shell au script ya kuzunguka: symlink hufanya kazi katika kila shell na katika scripts, haiongezi uanzishaji wa ziada wa shell (takriban ms 1 hadi 2), na uanzishaji bado husoma faili moja. `scripts/perf-check.sh` hupima njia fupi kwa `hyperfine` dhidi ya kikomo kile kile cha ms 10 (p95) kama `ccam <name>`.

### Kwenye Windows

Njia fupi ni `ccwork.exe` katika `%USERPROFILE%\.local\bin`, na `shim.Name` huondoa `.exe` (kwa herufi kubwa au ndogo) kabla ya kutafuta, kwa hiyo `ccwork`, `ccwork.exe` na `CCWORK.EXE` zote huanzisha `work`. ccam kwanza huunda kiungo cha ishara (symbolic link), ambacho Windows huruhusu katika Developer Mode au kwa msimamizi. Bila haki hiyo huunda hard link: jina la pili la faili ile ile ya `ccam.exe`, ambalo halihitaji ruhusa maalum. Hard link huhesabiwa kuwa ya ccam inapokuwa faili ile ile kama `ccam.exe`. Kwa sababu `ccam update` huweka faili jipya mahali pa `ccam.exe`, baadaye hufanya upya kila njia fupi ya hard link kwa `shim.Relink` na huchapisha `✓ shortcut ccwork renewed`. Ukaguzi wa jina hutafuta `<name>` pamoja na kila kiendelezi katika `PATHEXT` (`.exe`, `.cmd` na kadhalika) katika folda nyingine za `PATH`.

## Usalama na hali maalum

- **Kiungo kisicho cha akaunti yoyote.** Ikiwa kiungo cha ccam katika `~/.local/bin` kina jina ambalo hakuna akaunti inayolitumia, kwa mfano baada ya kuhariri `shortcut` katika `config.toml` kwa mkono, ccam hairudi kwenye akaunti chaguo-msingi. Husimama na `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Weka jina tena kwa `ccam shortcut`, au futa kiungo.
- **Faili zilizopo.** ccam kamwe haibadilishi faili ambayo haikuiunda. `ccam shims install` husimama na `<path> already exists and is not a ccam link, so it was left alone`.
- **Binary ya ccam iliyohamishwa.** Viungo vya zamani basi huelekea njia ambayo haipo tena, na ccam haivihesabu kuwa vyake. Vifute, kisha endesha `ccam shims install`.
- **Kubadilisha jina na kuondoa.** `ccam rename` huhifadhi jina la njia fupi, na kiungo huendelea kufanya kazi, kwa sababu kinaelekea ccam na ccam hutafuta jina katika `config.toml`. `ccam remove` hufuta kiungo cha akaunti. `ccam shortcut <account> <new>` hufuta kiungo cha zamani kabla ya kuunda kipya.
- **PATH.** ccam inapounda kiungo na `~/.local/bin` haiko kwenye `PATH` yako, huchapisha `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` ikifuatiwa na mstari wa `export PATH=...` wa kuongeza. Kwenye Windows huchapisha badala yake mstari wa PowerShell unaoongeza folda hiyo kwenye `PATH` ya mtumiaji wako.
- **Kukagua baadaye.** `ccam doctor` huripoti `missing shortcuts in /home/me/.local/bin: ccwork` na huonya wakati folda haiko kwenye `PATH`. `ccam doctor --fix` huunda viungo vinavyokosekana.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Kusimamia akaunti](/features/accounts)
- [Marejeo ya amri](/docs/commands#shortcuts)
- [Sakinisha ccam](/docs/install)
