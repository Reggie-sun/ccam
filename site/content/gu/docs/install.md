---
title: ccam ઇન્સ્ટૉલ કરો
description: Linux કે macOS પર ccam એક કમાન્ડથી, .deb પેકેજથી કે સોર્સમાંથી ઇન્સ્ટૉલ કરો. પહેલાં Claude Code ઇન્સ્ટૉલ થયેલું હોવું જોઈએ.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux કે macOS પર ટર્મિનલમાં આ ચલાવો:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

સ્ક્રિપ્ટ તમારી સિસ્ટમ માટેનું રિલીઝ ડાઉનલોડ કરે છે, તેનો SHA-256 ચેકસમ તપાસે છે અને `~/.local/bin` માં બે કમાન્ડ મૂકે છે: `ccam` અને તેનું ટૂંકું ઉપનામ `ccs`. જો એ ફોલ્ડર તમારા `PATH` માં ન હોય, તો સ્ક્રિપ્ટ `~/.bashrc` કે `~/.zshrc` માં ઉમેરવાની લાઇન છાપે છે.

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

પેકેજ `/usr/bin/ccam` અને `/usr/bin/ccs` લિંક ઇન્સ્ટૉલ કરે છે.

**સોર્સમાંથી.** Go 1.26 કે તેથી નવા વર્ઝન સાથે:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

Completion `ccam` અને `ccs` બંનેને આવરી લે છે, એકાઉન્ટનાં નામ સહિત.

## ઇન્સ્ટૉલેશન તપાસો

```sh
ccam version
ccam doctor
```

`ccam doctor` તપાસે છે કે Claude Code મળે છે અને દરેક એકાઉન્ટ ફોલ્ડર, લિંક અને લૉગિન બરાબર છે. પછી [ક્વિક સ્ટાર્ટ](/docs/quick-start) સાથે આગળ વધો.

## અનઇન્સ્ટૉલ કરો

બંને કમાન્ડ અને એકાઉન્ટ શૉર્ટકટ દૂર કરો, પછી જો તમને એકાઉન્ટની હવે જરૂર ન હોય તો ડેટા ફોલ્ડર પણ દૂર કરો:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` કાઢી નાખવાથી તેમાં રાખેલાં લૉગિન અને વાતચીતનો ઇતિહાસ પણ કાઢી નખાય છે. તમે [ઇમ્પોર્ટ કરેલું](/docs/commands) `~/.claude` ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે.
