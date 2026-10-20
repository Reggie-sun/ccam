---
title: ccam-ஐ நிறுவுதல்
metaTitle: Linux, macOS, Windows-இல் ccam-ஐ நிறுவுதல்
description: Linux, macOS அல்லது Windows-இல் ccam-ஐ ஒரே கட்டளை, .deb தொகுப்பு அல்லது மூலக் குறியீட்டால் நிறுவி, ccam update-ஆல் புதுப்பியுங்கள். முதலில் Claude Code தேவை.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07eadc0769c8
---

Linux அல்லது macOS-இல் ஒரு டெர்மினலில் இதை இயக்குங்கள்:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

இந்த script உங்கள் கணினி அமைப்புக்கான வெளியீட்டைப் பதிவிறக்கி, அதன் SHA-256 checksum-ஐச் சரிபார்த்து, `ccam`-ஐ `~/.local/bin`-இல் வைக்கிறது; பின்னர் நீங்கள் தேர்ந்தெடுக்கும் `ccwork` போன்ற குறுக்குவழிக் கட்டளைகளையும் ccam அதே கோப்புறையில் உருவாக்கும். அந்தக் கோப்புறை உங்கள் `PATH`-இல் இல்லையென்றால், `~/.bashrc` அல்லது `~/.zshrc`-இல் சேர்க்க வேண்டிய வரியை script அச்சிடும்.

Windows-இல், PowerShell-இல் இதை இயக்குங்கள்:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

இது `ccam_windows_<arch>.zip`-ஐப் பதிவிறக்கி, அதன் SHA-256 checksum-ஐச் சரிபார்த்து, Claude Code-இன் நிறுவி `claude.exe`-ஐ வைக்கும் அதே கோப்புறையான `%USERPROFILE%\.local\bin`-இல் `ccam.exe`-ஐ வைக்கிறது, அந்தக் கோப்புறை உங்கள் பயனர் `PATH`-இல் இன்னும் இல்லையென்றால் அதைச் சேர்க்கிறது. பிறகு புதிய terminal-ஐத் திறங்கள்.

ccam உங்களுக்காக Claude Code-ஐத் தொடங்குகிறது, அதனால் முதலில் [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)-ஐ நிறுவி, `claude --version` வேலை செய்கிறதா என்று உறுதிசெய்யுங்கள்.

## ஆதரிக்கப்படும் கணினி அமைப்புகள்

- x86-64 மற்றும் ARM64-இல் Linux: Ubuntu, Debian, Fedora, Arch, Alpine மற்றும் பிற distributions. binary static ஆனது, அதற்கு எந்தச் சார்புகளும் இல்லை.
- Apple silicon மற்றும் Intel-இல் macOS.
- x86-64 மற்றும் ARM64-இல் Windows 10 மற்றும் 11, PowerShell, Windows Terminal அல்லது Git Bash-இல். பகிரப்பட்ட கோப்புறைகள் எல்லாப் பயனர்களுக்கும் வேலை செய்யும். பகிரப்பட்ட கோப்புகளுக்கு (`settings.json`, `CLAUDE.md`, `history.jsonl`) Developer Mode (Settings > System > For developers) தேவை, ஏனெனில் Windows அந்த mode-இல் அல்லது administrator-க்கு மட்டுமே கோப்புகளுக்கு இணைப்பு உருவாக்க அனுமதிக்கிறது; பார்க்க [Windows-இல் பகிர்தல்](/features/share#on-windows).
- WSL-க்குள், ccam மற்ற எந்த Linux அமைப்பிலும் போலவே வேலை செய்யும்.

## நிறுவுவதற்கான பிற வழிகள்

**Debian மற்றும் Ubuntu தொகுப்பு.** [வெளியீடுகள் பக்கத்திலிருந்து](https://github.com/jhyunwoo/ccam/releases) `ccam_<version>_linux_amd64.deb`-ஐப் பதிவிறக்கி நிறுவுங்கள்:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

இந்தத் தொகுப்பு `/usr/bin/ccam`-ஐ நிறுவுகிறது.

**மூலக் குறியீட்டிலிருந்து.** Go 1.26 அல்லது புதிய பதிப்புடன்:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**ஒரு குறிப்பிட்ட பதிப்பு.** `CCAM_VERSION`-ஐ அமையுங்கள்:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion, கணக்குப் பெயர்கள் உட்பட `ccam`-இன் எல்லாக் கட்டளைகளுக்கும் விருப்பங்களுக்கும் வேலை செய்யும்.

## நிறுவலைச் சரிபார்த்தல்

```sh
ccam version
ccam doctor
```

Claude Code கண்டறியப்படுகிறதா, ஒவ்வொரு கணக்குக் கோப்புறையும் இணைப்பும் உள்நுழைவும் சரியாக உள்ளதா என்று `ccam doctor` சரிபார்க்கிறது. [விரைவுத் தொடக்கத்துடன்](/docs/quick-start) தொடருங்கள்.

## புதுப்பித்தல்

```sh
ccam update --check
ccam update
```

`ccam update --check` புதிய வெளியீடு உள்ளதா என்று மட்டுமே சொல்லும். `ccam update` உங்கள் கணினி அமைப்புக்கான வெளியீட்டைப் பதிவிறக்கி, அந்த வெளியீட்டின் `checksums.txt`-உடன் ஒப்பிட்டு அதன் SHA-256 checksum-ஐச் சரிபார்த்து, பிறகு `ccam` கோப்பை ஒரே படியில் மாற்றுகிறது. checksum பொருந்தவில்லை என்றால், எதுவும் மாற்றப்படாது. `ccwork` போன்ற குறுக்குவழிகள் ccam-க்கான இணைப்புகள் என்பதால், அவை உடனே புதிய பதிப்பை இயக்கும்.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam தானாகப் புதுப்பிப்புகளைச் சரிபார்ப்பதில்லை; நீங்கள் இந்தக் கட்டளையை இயக்கும்போது மட்டுமே புதிய வெளியீட்டைத் தேடும். பழைய வெளியீடு உட்பட, ஒரு குறிப்பிட்ட வெளியீட்டை நிறுவ `--version`-ஐக் கொடுங்கள். `ccam --dry-run update` எதையும் மாற்றாமல், எது பதிவிறக்கப்படும், எது மாற்றப்படும் என்பதைக் காட்டும்.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

நிறுவல் script-ஐ மீண்டும் இயக்கினாலும் `~/.local/bin`-இல் உள்ள நகல் புதுப்பிக்கப்படும். ccam-ஐ `go install` மூலம் build செய்திருந்தால், அதே கட்டளையை மீண்டும் இயக்குங்கள்.

**Windows.** `ccam update` Windows-க்கான `.zip`-ஐப் பதிவிறக்குகிறது. Windows-இல் இயங்கும் நிரலை மாற்ற முடியாது, அதனால் பழைய binary `ccam.exe.old` எனப் பெயர் மாற்றப்படுகிறது, அடுத்த புதுப்பிப்பு அதை நீக்குகிறது.

**Debian மற்றும் Ubuntu தொகுப்பு.** apt நிறுவிய `/usr/bin/ccam`-ஐ `ccam update` மாற்றாது, அதற்குப் பதிலாக இந்தப் படிகளை அச்சிடும். [வெளியீடுகள் பக்கத்திலிருந்து](https://github.com/jhyunwoo/ccam/releases) புதிய `ccam_<version>_linux_amd64.deb`-ஐப் பதிவிறக்கி, பழையதன் மேலேயே நிறுவுங்கள்:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## நீக்குதல்

ccam நிறுவப்பட்டிருக்கும்போதே முதலில் கணக்குக் குறுக்குவழிகளை நீக்குங்கள். `ccam list` அவற்றை Shortcut நெடுவரிசையில் காட்டும்; ஒவ்வொன்றையும் `ccam shortcut <account> --remove` மூலம் நீக்குங்கள், அல்லது அந்த இணைப்புகளை `~/.local/bin`-இலிருந்து நீங்களே நீக்குங்கள். பிறகு கட்டளையை நீக்குங்கள், கணக்குகள் இனி தேவையில்லை என்றால் தரவுக் கோப்புறைகளையும் நீக்குங்கள்:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows-இல், PowerShell-இல்:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam`-ஐ நீக்கினால், அங்கே சேமிக்கப்பட்ட உள்நுழைவுகளும் உரையாடல் வரலாறும் நீக்கப்படும். நீங்கள் [இறக்குமதி செய்த](/docs/commands) `~/.claude` கோப்புறை அது இருக்கும் இடத்திலேயே இருக்கும்.
