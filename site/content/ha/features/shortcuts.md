---
title: Gajerun umarni
description: Ba asusu umarninsa kamar ccwork da ccam shortcut, kuma sake ƙirƙirar mahaɗai a ~/.local/bin da ccam shims install.
updated: 2026-10-23
order: 3
translation: machine
sourceHash: 38e2d1aba397
---

Gajeren umarni umarni ne da kai da kanka ka sa masa suna, wanda ke fara Claude Code a matsayin asusu ɗaya: `ccwork --resume` yana yin abin da `ccam work --resume` ke yi. Mahaɗi ne zuwa ccam a cikin `~/.local/bin`, don haka yana aiki a bash, zsh, rubutun (scripts) da duk wani abu da ke gudanar da shirye-shirye daga `PATH`. Gajerun umarni zaɓi ne; ƙirƙiri ɗaya ga asusan da kake farawa sau da yawa.

## Amfani

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccam yana tabbatar da kowane canji:

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

Asusu yana da gajeren umarni ɗaya a mafi yawa, wanda aka ajiye a matsayin `shortcut` a `~/.config/ccam/config.toml`. Mataimakin `ccam add` yana tambaya `Create a shortcut command for <name>? [y/N]` sannan `Shortcut name [cc<name>]:`. A wajen mataimakin, `ccam add` yana ƙirƙirar ɗaya ne kawai da `--shortcut`. `ccam list` yana nuna gajeren umarnin kowane asusu, ko `-`.

Gajeren umarni yana miƙa kowane argument ga Claude Code, kamar yadda `ccam <name>` ke yi bayan sunan. Wannan ya haɗa da `--dry-run` da `--help`: don duba farawa kafin lokaci, yi amfani da `ccam --dry-run work`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--remove` | Goge gajeren umarnin asusun da mahaɗinsa | a kashe |
| `-h`, `--help` | Taimako ga umarnin | |
| `--dry-run` | `ccam shortcut` yana karɓar sa: yana duba sunan kuma yana buga `Would set the shortcut of work to ccwork.` ko `Would remove ccwork, the shortcut of work.` ba tare da ajiyewa ko haɗawa ba | a kashe |

`ccam shims install` ba ya karɓar arguments kuma yana ƙin `--dry-run` da `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Ƙa'idodin suna, ana duba su a wannan tsari:

- haruffa 1 zuwa 64 na `a-z`, `0-9`, `.`, `-` da `_`, suna farawa da harafi ko lamba;
- ba `ccam`, `claude` ko umarnin ccam kamar `run` ko `usage` ba;
- ba gajeren umarnin wani asusu ba;
- ba fayil a `~/.local/bin` wanda ba mahaɗin ccam ba ne tuni (ko gajeren umarnin wannan asusun kansa, wanda ccam ya yi) ba;
- ba sunan wani shiri a `PATH` naka ba.

Ana duba `config.toml` da ƙa'idodi iri ɗaya lokacin da yake lodawa: gajeren umarnin da aka gyara da hannu wanda ya karya tsarin, ko shi ne `claude` ko `ccam`, ko an ba asusu biyu shi, yana sa kowane umarni ya gaza da saƙon da ke ambaton asusun, har sai an gyara shi.

Ba da umarni da `--remove` tare yana tsayawa da `give a command or --remove, not both`. Saita sunan da asusu ya riga yana da shi yana buga `work: ccwork (unchanged)` kuma yana sake ƙirƙirar mahaɗinsa idan ya ɓace.

## Yadda yake aiki

Saita suna yana gudanar da `checkShortcut`, wanda ke kiran `shim.Check` tare da foldar bin da ainihin hanyar binary na ccam:

<!-- motion:shortcuts-check -->

Sunan yana sauka ta dubawar a tsari. Kusa da kowane dubawa akwai sunan da yake ƙi: `Work` ya gaza tsarin, `run` umarnin ccam ne, `ccside` na wani asusu ne, `ccfile` fayil ne da ke akwai a `~/.local/bin`, kuma `git` zai ɓoye `/usr/bin/git`. Sunan da ya wuce kowane dubawa ne kawai ake ajiyewa da haɗawa.

1. `shim.Check` yana kwatanta sunan da `^[a-z0-9][a-z0-9._-]{0,63}$`, sannan `config.Reserved`, sannan `cfg.FindShortcut`. Yana kiran `os.Lstat` a kan `~/.local/bin/<name>`, kuma ga kowace foldar dabam a `PATH` yana kiran `os.Stat` kuma yana ƙin fayil mai gudana (executable) wanda ba mahaɗi zuwa ccam ba ne.
2. ccam yana ajiye `config.toml` da `fsutil.AtomicWrite` (fayil na wucin gadi, `fsync`, rename).
3. `shim.Remove` yana goge mahaɗin tsohon gajeren umarnin, amma sai idan mahaɗi ne zuwa ccam.
4. `shim.Install` yana bi ta kowane gajeren umarni da aka saita bisa tsarin suna kuma yana ƙirƙirar mahaɗin kowanne da ya ɓace, da `os.Symlink` zuwa binary na ccam, wanda aka nemo da `os.Executable` da `filepath.EvalSymlinks`. Fayil da ke tare hanyar gajeren umarni ɗaya ba ya hana sauran; kuskuren yana ambaton kowanne.

ccam yana rubuta kowane mahaɗi da ya yi, da binary da ya sa shi ya gudanar, a `~/.local/state/ccam/shortcuts.json`. Ana ɗaukar mahaɗi a matsayin na ccam idan `os.Readlink` ya mayar da hanyar ccam, ko idan mahaɗin da hanyar ccam suna kai wa fayil ɗaya. Mahaɗin da aka rubuta wanda aka matsar ko aka maye gurbin binary ɗinsa tun daga lokacin (ccam da aka matsar zuwa wata folda, sake girkawa) ya tsufa, kuma `shim.Install` yana nuna shi zuwa wannan ccam. Mahaɗin da aka rubuta zuwa wani ccam da har yanzu yake a girke ana barin sa yadda yake. Kowace cirewa da kowane dubawar maye gurbi suna bin waɗannan gwaje-gwaje, don haka ccam ba ya taɓa fayilolin da bai ƙirƙira ba, kamar kwafinka na ccam ko mahaɗin da ka yi da kanka.

Idan ka gudanar da gajeren umarni, sunan da ka rubuta shi ne kawai abin da ccam ke buƙata:

<!-- motion:shortcuts-link -->

Shell yana samun `ccwork` a `PATH`, kernel yana bin mahaɗin kuma yana fara binary na ccam, kuma `argv[0]` har yanzu `ccwork` ne. A `app.Main`, duk sunan da ba `ccam` ba yana loda `config.toml` sau ɗaya kuma yana kiran `shim.Lookup`, wanda ke nemo asusun da `shortcut` ɗinsa ya yi daidai da `filepath.Base(argv[0])`. Sannan `launcher.LaunchWith` yana bin matakai iri ɗaya da [`ccam <name>`](/features/launch) da saitin da aka riga aka loda, kuma yana miƙa `argv[1:]` ga `claude` ba tare da taɓawa ba.

Me ya sa symlink ba alias na shell ko rubutun naɗewa (wrapper script) ba: symlink yana aiki a kowane shell da a rubutun, ba ya ƙara farawar shell (kusan 1 zuwa 2 ms), kuma farawar har yanzu tana karanta fayil ɗaya. `scripts/perf-check.sh` yana auna gajeren umarni da `hyperfine` a kan iyakar 10 ms (p95) ɗaya da `ccam <name>`.

### A Windows

Gajeren umarni shi ne `ccwork.exe` a `%USERPROFILE%\.local\bin`, kuma `shim.Name` yana cire `.exe` (a kowane irin rubutun haruffa) kafin nema, don haka `ccwork`, `ccwork.exe` da `CCWORK.EXE` duka suna fara `work`. ccam yana fara ƙirƙirar symbolic link, wanda Windows ke bari a Developer Mode ko ga administrator. Ba tare da wannan izinin ba yana ƙirƙirar hard link: suna na biyu ga fayil ɗin `ccam.exe` ɗaya, wanda ba ya buƙatar wani izini na musamman. Hard link yana ƙirguwa a matsayin na ccam idan fayil ɗaya ne da `ccam.exe`. Hard link yana ci gaba da gudanar da fayil ɗin da aka yi shi daga gare shi, don haka idan `ccam update` ya sa sabon fayil a wurin `ccam.exe`, yana sabunta kowane gajeren umarni da ke hard link da `shim.Relink` kuma yana buga `✓ shortcut ccwork renewed`. Sake gudanar da `install.ps1` yana yin haka ta `ccam shims install` (`OK 2 shortcuts now run the new ccam`), haka ma `ccam shims install` da `ccam doctor --fix` ga duk gajeren umarnin da ke kwafin tsohon ccam. Duba sunan yana neman `<name>` da kowane extension a `PATHEXT` (`.exe`, `.cmd` da sauransu) a sauran foldojin `PATH`.

Rubuta gajeren umarni a shell (PowerShell, cmd.exe, Git Bash). Symbolic link da aka buɗe daga Explorer, akwatin Run ko `Start-Process` yana farawa a matsayin `ccam` kawai, wanda ke buɗe dashboard maimakon asusun.

## Tsaro da yanayi na musamman

- **Mahaɗin da babu asusun da ke da shi.** Idan mahaɗi zuwa ccam a `~/.local/bin` yana da sunan da babu asusun da ke amfani da shi, misali bayan ka gyara `shortcut` a `config.toml` da hannu, ccam ba ya komawa ga asusun tsoho. Yana tsayawa da `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Sake saita sunan da `ccam shortcut`, ko goge mahaɗin.
- **Fayilolin da ke akwai.** ccam ba ya taɓa maye gurbin fayil da bai ƙirƙira ba. `ccam shims install` yana bayar da rahoton `<path> already exists and is not a ccam link, so it was left alone` kuma har yanzu yana ƙirƙirar sauran mahaɗan. Sannan `ccam add` da `ccam import` suna buga `shortcut: ccwork (not created; see above)`.
- **Binary na ccam da aka matsar.** Mahaɗan da ccam ya yi sai su nuna hanyar da babu ita kuma. `ccam shims install` (ko `ccam doctor --fix`) yana nuna su zuwa ccam ɗin da ka gudanar da shi da shi.
- **Sake suna da cirewa.** `ccam rename` yana riƙe sunan gajeren umarnin, kuma mahaɗin yana ci gaba da aiki, domin yana nuna ccam kuma ccam yana neman sunan a `config.toml`. `ccam remove` yana goge mahaɗin asusun. `ccam shortcut <account> <new>` yana goge tsohon mahaɗin kafin ya ƙirƙiri sabon.
- **PATH.** Idan ccam ya ƙirƙiri mahaɗi kuma `~/.local/bin` ba ya cikin `PATH` naka, yana buga `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` sannan layin `export PATH=...` da za a ƙara. A Windows yana buga layin PowerShell da ke ƙara foldar a `PATH` na mai amfaninka maimakon haka.
- **Dubawa daga baya.** `ccam doctor` yana bayar da rahoton `missing shortcuts in /home/me/.local/bin: ccwork`, gajerun umarnin da ke gudanar da tsohon ccam, kuma yana gargaɗi idan foldar ba ta cikin `PATH`. `ccam doctor --fix` yana ƙirƙirar mahaɗan da suka ɓace kuma yana sabunta waɗanda suka tsufa, kuma yana cewa `→ fixed` ne kawai ga waɗanda ya yi. Gajeren umarnin da ke gudanar da wani ccam da har yanzu yake a girke yana aiki kuma bayani ne kawai; `ccam doctor --fix` yana nuna shi zuwa ccam ɗin da ka gudanar da shi da shi, kuma `ccam add` da `ccam import` suna nuna shi a matsayin ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Sarrafa asusai](/features/accounts)
- [Bayanin umarnai](/docs/commands#shortcuts)
- [Girka ccam](/docs/install)
