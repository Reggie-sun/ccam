---
title: ccam ઇન્સ્ટૉલ કરો
description: Linux કે macOS પર ccam એક કમાન્ડથી, .deb પેકેજથી કે સોર્સમાંથી ઇન્સ્ટૉલ કરો, અને ccam update વડે તેને અપ-ટુ-ડેટ રાખો.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 4d7291e0f86e
---

Linux કે macOS પર ટર્મિનલમાં આ ચલાવો:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

સ્ક્રિપ્ટ તમારી સિસ્ટમ માટેની રિલીઝ ડાઉનલોડ કરે છે, તેનો SHA-256 checksum તપાસે છે અને `ccam` ને `~/.local/bin` માં મૂકે છે; તમે પછીથી પસંદ કરો તે શૉર્ટકટ કમાન્ડ, જેમ કે `ccwork`, ccam આ જ ફોલ્ડરમાં બનાવે છે. જો આ ફોલ્ડર તમારા `PATH` માં ન હોય, તો સ્ક્રિપ્ટ `~/.bashrc` કે `~/.zshrc` માં ઉમેરવાની લાઇન છાપે છે.

ccam તમારા માટે Claude Code શરૂ કરે છે, તેથી પહેલાં [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ઇન્સ્ટૉલ કરો અને ખાતરી કરો કે `claude --version` કામ કરે છે.

## સમર્થિત સિસ્ટમ

- x86-64 અને ARM64 પર Linux: Ubuntu, Debian, Fedora, Arch, Alpine અને બીજાં ડિસ્ટ્રિબ્યુશન. બાઇનરી સ્ટેટિક છે અને તેને કોઈ ડિપેન્ડન્સીની જરૂર નથી.
- Apple silicon અને Intel પર macOS.
- Windows સીધું સમર્થિત નથી. WSL ની અંદર ccam કોઈ પણ Linux સિસ્ટમની જેમ કામ કરે છે.

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

## શેલ completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion `ccam` ના બધા કમાન્ડ અને વિકલ્પોને આવરી લે છે, એકાઉન્ટનાં નામ સહિત.

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

ccam ક્યારેય જાતે અપડેટ તપાસતું નથી; તમે આ કમાન્ડ ચલાવો ત્યારે જ તે નવી રિલીઝ શોધે છે. કોઈ ચોક્કસ રિલીઝ, જૂની હોય તો પણ, ઇન્સ્ટૉલ કરવા માટે `--version` આપો. `ccam --dry-run update` કંઈ બદલ્યા વગર બતાવે છે કે શું ડાઉનલોડ થશે અને શું બદલાશે.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ઇન્સ્ટૉલ સ્ક્રિપ્ટ ફરીથી ચલાવવાથી પણ `~/.local/bin` માંની કૉપી અપડેટ થાય છે. જો તમે ccam ને `go install` વડે બનાવ્યું હોય, તો એ જ કમાન્ડ ફરીથી ચલાવો.

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

`~/.local/share/ccam` કાઢી નાખવાથી તેમાં રાખેલાં લૉગિન અને વાતચીતનો ઇતિહાસ પણ કાઢી નખાય છે. તમે [ઇમ્પોર્ટ કરેલું](/docs/commands) `~/.claude` ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે.
