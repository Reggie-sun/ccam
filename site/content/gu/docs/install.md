---
title: ccam ઇન્સ્ટૉલ કરો
metaTitle: Linux, macOS અને Windows પર ccam ઇન્સ્ટૉલ કરો
description: Linux, macOS કે Windows પર ccam એક કમાન્ડથી, .deb પેકેજથી કે સોર્સમાંથી ઇન્સ્ટૉલ કરો, અને ccam update વડે અપડેટ કરો. પહેલાં Claude Code ઇન્સ્ટૉલ હોવું જોઈએ.
updated: 2026-10-21
order: 1
translation: machine
sourceHash: 1e3a53623de0
---

Linux કે macOS પર ટર્મિનલમાં આ ચલાવો:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

સ્ક્રિપ્ટ તમારી સિસ્ટમ માટેની રિલીઝ ડાઉનલોડ કરે છે, તેનો SHA-256 checksum તપાસે છે અને `ccam` ને `~/.local/bin` માં મૂકે છે; તમે પછીથી પસંદ કરો તે શૉર્ટકટ કમાન્ડ, જેમ કે `ccwork`, ccam આ જ ફોલ્ડરમાં બનાવે છે. જો આ ફોલ્ડર તમારા `PATH` માં ન હોય, તો સ્ક્રિપ્ટ `~/.bashrc` કે `~/.zshrc` માં ઉમેરવાની લાઇન છાપે છે.

Windows પર, PowerShell માં આ ચલાવો:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

તે `ccam_windows_<arch>.zip` ડાઉનલોડ કરે છે, તેનો SHA-256 checksum તપાસે છે, `ccam.exe` ને `%USERPROFILE%\.local\bin` માં મૂકે છે, એ જ ફોલ્ડર જ્યાં Claude Code નું ઇન્સ્ટૉલર `claude.exe` મૂકે છે, અને જો તે ફોલ્ડર તમારા યુઝર `PATH` માં ન હોય તો તેને ઉમેરે છે. પછી નવું ટર્મિનલ ખોલો.

ccam તમારા માટે Claude Code શરૂ કરે છે, તેથી પહેલાં [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ઇન્સ્ટૉલ કરો અને ખાતરી કરો કે `claude --version` કામ કરે છે.

## સમર્થિત સિસ્ટમ

- x86-64 અને ARM64 પર Linux: Ubuntu, Debian, Fedora, Arch, Alpine અને બીજાં ડિસ્ટ્રિબ્યુશન. બાઇનરી સ્ટેટિક છે અને તેને કોઈ ડિપેન્ડન્સીની જરૂર નથી.
- Apple silicon અને Intel પર macOS.
- x86-64 અને ARM64 પર Windows 10 અને 11, PowerShell, Windows Terminal કે Git Bash માં. શેર કરેલાં ફોલ્ડર દરેક યુઝર માટે કામ કરે છે. શેર કરેલી ફાઇલો (`settings.json`, `CLAUDE.md`, `history.jsonl`) માટે Developer Mode (Settings > System > For developers) જોઈએ, કારણ કે Windows ફક્ત એ મોડમાં, કે ઍડમિનિસ્ટ્રેટરને જ, ફાઇલની લિંક બનાવવા દે છે; જુઓ [Windows પર શેરિંગ](/features/share#on-windows).
- WSL ની અંદર ccam કોઈ પણ Linux સિસ્ટમની જેમ કામ કરે છે.

## ઇન્સ્ટૉલ કરવાની બીજી રીતો

**Debian અને Ubuntu પેકેજ.** [રિલીઝ પેજ](https://github.com/jhyunwoo/ccam/releases) પરથી `ccam_<version>_linux_amd64.deb` ડાઉનલોડ કરો અને તેને ઇન્સ્ટૉલ કરો:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

પેકેજ `/usr/bin/ccam` ઇન્સ્ટૉલ કરે છે.

**સોર્સમાંથી.** Go 1.26 કે તેથી નવા વર્ઝન સાથે:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**કોઈ ચોક્કસ વર્ઝન.** `CCAM_VERSION` સેટ કરો:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## શેલ completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion `ccam` ના બધા કમાન્ડ અને વિકલ્પોને આવરી લે છે, એકાઉન્ટનાં નામ સહિત.

## Claude Code plugin

ccam plugin Claude Code ને પોતાને ccam વાપરતાં શીખવે છે: કયા એકાઉન્ટમાં જગ્યા છે તે તપાસવું, `ccam pick` વડે જુદાં એકાઉન્ટ પર સમાંતર `claude -p` workers શરૂ કરવા, અને `ccam ext` વડે ખૂટતી skill કે MCP સર્વર સમજાવવું. તેને Claude Code ની અંદર ઇન્સ્ટૉલ કરો:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

તે એક skill, `/ccam:usage`, `/ccam:ext` અને `/ccam:switch` commands, અને બે hooks ઉમેરે છે: કોઈ સેશન વપરાશ મર્યાદાએ અટકે ત્યારે તમને ટર્મિનલમાં સૂચના મળે છે, અને `/ccam:switch [account]` વાતચીત બીજા એકાઉન્ટને સોંપે છે ([વિગતો](/features/next#inside-claude-code)). hooks `ccam` ચલાવે છે, તેથી Claude Code જે `PATH` જુએ છે તેમાં તે હોવું જોઈએ. skill ક્યારેય લૉગિન કરતા કે ફાઇલો ખસેડતા commands ચલાવતી નથી; તેના બદલે તે તેમને સૂચવે છે.

## ઇન્સ્ટૉલેશન તપાસો

```sh
ccam version
ccam doctor
```

`ccam doctor` તપાસે છે કે Claude Code મળે છે અને દરેક એકાઉન્ટ ફોલ્ડર, લિંક અને લૉગિન બરાબર છે. પછી [ક્વિક સ્ટાર્ટ](/docs/quick-start) સાથે આગળ વધો.

## અપડેટ

```sh
ccam update --check
ccam update
```

`ccam update --check` ફક્ત એટલું જણાવે છે કે નવી રિલીઝ છે કે નહીં. `ccam update` તમારી સિસ્ટમ માટેની રિલીઝ ડાઉનલોડ કરે છે, રિલીઝની `checksums.txt` સાથે મેળવીને તેનો SHA-256 checksum તપાસે છે અને પછી `ccam` ફાઇલને એક જ પગલામાં બદલે છે. જો checksum મેળ ન ખાય, તો કંઈ બદલાતું નથી. `ccwork` જેવા શૉર્ટકટ ccam ની લિંક છે, તેથી તે તરત જ નવું વર્ઝન ચલાવે છે.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ડેશબોર્ડ (ટર્મિનલ પર `ccam`) પણ દરેક વખતે ખૂલે ત્યારે નવી રિલીઝ શોધે છે અને જો હોય તો તે ઑફર કરે છે: `U` દબાવો અને પુષ્ટિ કરો. બીજો કોઈ કમાન્ડ શોધતો નથી, અને ccam પૂછ્યા વગર ક્યારેય અપડેટ થતું નથી. કોઈ ચોક્કસ રિલીઝ, જૂની હોય તો પણ, ઇન્સ્ટૉલ કરવા માટે `--version` આપો. `ccam --dry-run update` કંઈ બદલ્યા વગર બતાવે છે કે શું ડાઉનલોડ થશે અને શું બદલાશે.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ઇન્સ્ટૉલ સ્ક્રિપ્ટ ફરીથી ચલાવવાથી પણ `~/.local/bin` માંની કૉપી અપડેટ થાય છે. જો તમે ccam ને `go install` વડે બનાવ્યું હોય, તો એ જ કમાન્ડ ફરીથી ચલાવો.

**Windows.** `ccam update` Windows માટેની `.zip` ડાઉનલોડ કરે છે. Windows પર ચાલતો પ્રોગ્રામ બદલી શકાતો નથી, તેથી જૂની બાઇનરીનું નામ બદલીને `ccam.exe.old` કરાય છે, જેને આગલું અપડેટ કાઢી નાખે છે.

**Debian અને Ubuntu પેકેજ.** apt એ ઇન્સ્ટૉલ કરેલા `/usr/bin/ccam` ને `ccam update` બદલતું નથી, અને તેના બદલે આ પગલાં છાપે છે. [રિલીઝ પેજ](https://github.com/jhyunwoo/ccam/releases) પરથી નવું `ccam_<version>_linux_amd64.deb` ડાઉનલોડ કરો અને તેને જૂના પેકેજ પર ઇન્સ્ટૉલ કરો:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## અનઇન્સ્ટૉલ કરો

ccam ઇન્સ્ટોલ હોય ત્યારે જ પહેલાં એકાઉન્ટના શૉર્ટકટ દૂર કરો. `ccam list` તેમને Shortcut કૉલમમાં બતાવે છે; દરેકને `ccam shortcut <account> --remove` વડે દૂર કરો, અથવા એ લિંક `~/.local/bin` માંથી જાતે કાઢી નાખો. પછી કમાન્ડ દૂર કરો, અને જો હવે એકાઉન્ટની જરૂર ન હોય તો ડેટા ફોલ્ડર પણ:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows પર, PowerShell માં:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` કાઢી નાખવાથી તેમાં રાખેલાં લૉગિન અને વાતચીતનો ઇતિહાસ પણ કાઢી નખાય છે. તમે [ઇમ્પોર્ટ કરેલું](/docs/commands) `~/.claude` ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે.
