---
title: ccam-ஐ நிறுவுதல்
description: Linux அல்லது macOS-இல் ccam-ஐ ஒரே கட்டளையில், .deb தொகுப்பில் அல்லது மூலக் குறியீட்டிலிருந்து நிறுவுங்கள். முதலில் Claude Code நிறுவப்பட்டிருக்க வேண்டும்.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux அல்லது macOS-இல் ஒரு டெர்மினலில் இதை இயக்குங்கள்:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

இந்த script உங்கள் கணினி அமைப்புக்கான வெளியீட்டைப் பதிவிறக்கி, அதன் SHA-256 checksum-ஐச் சரிபார்த்து, `~/.local/bin`-இல் இரண்டு கட்டளைகளை வைக்கிறது: `ccam` மற்றும் அதன் சுருக்கப் பெயர் `ccs`. அந்தக் கோப்புறை உங்கள் `PATH`-இல் இல்லையென்றால், `~/.bashrc` அல்லது `~/.zshrc`-இல் சேர்க்க வேண்டிய வரியை script அச்சிடும்.

ccam உங்களுக்காக Claude Code-ஐத் தொடங்குகிறது, அதனால் முதலில் [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)-ஐ நிறுவி, `claude --version` வேலை செய்கிறதா என்று உறுதிசெய்யுங்கள்.

## ஆதரிக்கப்படும் கணினி அமைப்புகள்

- x86-64 மற்றும் ARM64-இல் Linux: Ubuntu, Debian, Fedora, Arch, Alpine மற்றும் பிற distributions. binary static ஆனது, அதற்கு எந்தச் சார்புகளும் இல்லை.
- Apple silicon மற்றும் Intel-இல் macOS.
- Windows நேரடியாக ஆதரிக்கப்படவில்லை. WSL-க்குள், ccam மற்ற எந்த Linux அமைப்பிலும் போலவே வேலை செய்யும்.

## நிறுவுவதற்கான பிற வழிகள்

**Debian மற்றும் Ubuntu தொகுப்பு.** [வெளியீடுகள் பக்கத்திலிருந்து](https://github.com/jhyunwoo/ccam/releases) `ccam_<version>_linux_amd64.deb`-ஐப் பதிவிறக்கி நிறுவுங்கள்:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

இந்தத் தொகுப்பு `/usr/bin/ccam`-ஐயும் `/usr/bin/ccs` இணைப்பையும் நிறுவுகிறது.

**மூலக் குறியீட்டிலிருந்து.** Go 1.26 அல்லது புதிய பதிப்புடன்:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**ஒரு குறிப்பிட்ட பதிப்பு.** `CCAM_VERSION`-ஐ அமையுங்கள்:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion, கணக்குப் பெயர்கள் உட்பட `ccam` மற்றும் `ccs` இரண்டுக்கும் வேலை செய்யும்.

## நிறுவலைச் சரிபார்த்தல்

```sh
ccam version
ccam doctor
```

Claude Code கண்டறியப்படுகிறதா, ஒவ்வொரு கணக்குக் கோப்புறையும் இணைப்பும் உள்நுழைவும் சரியாக உள்ளதா என்று `ccam doctor` சரிபார்க்கிறது. [விரைவுத் தொடக்கத்துடன்](/docs/quick-start) தொடருங்கள்.

## நீக்குதல்

இரண்டு கட்டளைகளையும் கணக்குக் குறுக்குவழிகளையும் நீக்குங்கள். கணக்குகள் இனி தேவையில்லை என்றால், பிறகு தரவுக் கோப்புறைகளையும் நீக்குங்கள்:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam`-ஐ நீக்கினால், அங்கே சேமிக்கப்பட்ட உள்நுழைவுகளும் உரையாடல் வரலாறும் நீக்கப்படும். நீங்கள் [இறக்குமதி செய்த](/docs/commands) `~/.claude` கோப்புறை அது இருக்கும் இடத்திலேயே இருக்கும்.
