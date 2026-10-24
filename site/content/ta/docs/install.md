---
title: ccam-ஐ நிறுவுதல்
metaTitle: Linux, macOS, Windows-இல் ccam-ஐ நிறுவுதல்
description: Linux, macOS அல்லது Windows-இல் ccam-ஐ ஒரே கட்டளை, .deb தொகுப்பு அல்லது மூலக் குறியீட்டால் நிறுவி, ccam update-ஆல் புதுப்பியுங்கள். முதலில் Claude Code தேவை.
updated: 2026-10-23
order: 1
translation: machine
sourceHash: "af21dcc24169"
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

இது `ccam_windows_<arch>.zip`-ஐப் பதிவிறக்கி, அதன் SHA-256 checksum-ஐச் சரிபார்த்து, Claude Code-இன் நிறுவி `claude.exe`-ஐ வைக்கும் அதே கோப்புறையான `%USERPROFILE%\.local\bin`-இல் `ccam.exe`-ஐ வைக்கிறது, அந்தக் கோப்புறை உங்கள் பயனர் `PATH`-இல் இன்னும் இல்லையென்றால் அதைச் சேர்க்கிறது; `%USERPROFILE%\bin` போன்ற பதிவுகளை அப்படியே வைத்திருக்கிறது. பிறகு புதிய terminal-ஐத் திறங்கள். அதை மீண்டும் இயக்கினால் ccam மாற்றப்படும், பழைய binary-இலிருந்து ஒரு Claude Code அமர்வு இன்னும் இயங்கிக்கொண்டிருந்தாலும் கூட (அது `ccam.exe.old` எனப் பெயர் மாற்றப்படும், அதைவிடப் பழையது இன்னும் பயன்பாட்டில் இருந்தால் `ccam.exe.old-<number>` என), hard links ஆக உள்ள குறுக்குவழிகளும் புதுப்பிக்கப்படும்: `OK 2 shortcuts now run the new ccam`.

ccam உங்களுக்காக Claude Code-ஐத் தொடங்குகிறது, அதனால் முதலில் [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)-ஐ நிறுவி, `claude --version` வேலை செய்கிறதா என்று உறுதிசெய்யுங்கள்.

## ஆதரிக்கப்படும் கணினி அமைப்புகள்

- x86-64 மற்றும் ARM64-இல் Linux: Ubuntu, Debian, Fedora, Arch, Alpine மற்றும் பிற distributions. binary static ஆனது, அதற்கு எந்தச் சார்புகளும் இல்லை.
- Apple silicon மற்றும் Intel-இல் macOS.
- x86-64 மற்றும் ARM64-இல் Windows 10 மற்றும் 11, PowerShell, Windows Terminal அல்லது Git Bash-இல். பகிரப்பட்ட கோப்புறைகள் எல்லாப் பயனர்களுக்கும் வேலை செய்யும். பகிரப்பட்ட கோப்புகளுக்கு (`settings.json`, `CLAUDE.md`, `history.jsonl`) Developer Mode (Settings > System > For developers) தேவை, ஏனெனில் Windows அந்த mode-இல் அல்லது administrator-க்கு மட்டுமே கோப்புகளுக்கு இணைப்பு உருவாக்க அனுமதிக்கிறது; பார்க்க [Windows-இல் பகிர்தல்](/features/share#on-windows). அது இல்லாமல், முதல் `ccam add` MCP சர்வர்களையும் நீட்டிப்புகளையும் மட்டும் கொண்ட `default` குழுவை உருவாக்கி, Developer Mode இயக்கப்பட்ட பிறகு மீதியை எப்படிச் சேர்ப்பது என்று சொல்லும்.
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

ஒரு release tag-இன் build (`@latest` அல்லது `@vX.Y.Z`) அந்தப் பதிப்பைத் தெரிவிக்கும், அதனால் அதற்கு `ccam update` வேலை செய்யும்; `@main`-இன் build அல்லது ஒரு tag-க்குப் பிந்தைய local checkout-இன் build `dev` எனத் தெரிவிக்கும்.

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

## Claude Code plugin

ccam plugin, Claude Code-க்கே ccam-ஐப் பயன்படுத்தக் கற்றுக்கொடுக்கிறது: எந்தக் கணக்கில் இடம் உள்ளது என்று சரிபார்க்க, `ccam pick` மூலம் வெவ்வேறு கணக்குகளில் இணையான `claude -p` workers-ஐத் தொடங்க, காணாமல் போன skill அல்லது MCP சர்வரை `ccam ext` மூலம் விளக்க. அதை Claude Code-க்குள் நிறுவுங்கள்:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

அது ஒரு skill-ஐயும், `/ccam:usage`, `/ccam:ext`, `/ccam:switch` ஆகிய commands-ஐயும், இரண்டு hooks-ஐயும் சேர்க்கிறது: ஓர் அமர்வு பயன்பாட்டு வரம்பில் நிற்கும்போது உங்களுக்கு ஒரு டெர்மினல் அறிவிப்பு கிடைக்கும், `/ccam:switch [account]` உரையாடலை வேறொரு கணக்கிடம் ஒப்படைக்கும் ([விவரங்கள்](/features/next#inside-claude-code)). hooks `ccam`-ஐ இயக்குகின்றன, அதனால் Claude Code பார்க்கும் `PATH`-இல் அது இருக்க வேண்டும். அந்த skill உள்நுழையும் அல்லது கோப்புகளை நகர்த்தும் கட்டளைகளை ஒருபோதும் இயக்காது; அதற்குப் பதிலாக அவற்றைப் பரிந்துரைக்கும்.

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

டாஷ்போர்டும் (டெர்மினலில் `ccam`) ஒவ்வொரு முறை திறக்கும்போதும் புதிய வெளியீட்டைத் தேடி, ஒன்று இருந்தால் அதை வழங்கும்: `U`-ஐ அழுத்தி உறுதிப்படுத்துங்கள். வேறு எந்தக் கட்டளையும் தேடுவதில்லை, ccam கேட்காமல் ஒருபோதும் புதுப்பித்துக்கொள்வதில்லை. பழைய வெளியீடு உட்பட, ஒரு குறிப்பிட்ட வெளியீட்டை நிறுவ `--version`-ஐக் கொடுங்கள். `ccam --dry-run update` எதையும் மாற்றாமல், எது பதிவிறக்கப்படும், எது மாற்றப்படும் என்பதைக் காட்டும்.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

நிறுவல் script-ஐ மீண்டும் இயக்கினாலும் `~/.local/bin`-இல் உள்ள நகல் புதுப்பிக்கப்படும். ccam-ஐ ஒரு release tag-இலிருந்து `go install` மூலம் build செய்திருந்தால், `ccam update`-உம் வேலை செய்யும், அல்லது அதே கட்டளையை மீண்டும் இயக்குங்கள்.

**Windows.** `ccam update` Windows-க்கான `.zip`-ஐப் பதிவிறக்குகிறது. Windows-இல் இயங்கும் நிரலை மாற்ற முடியாது, அதனால் பழைய binary `ccam.exe.old` எனப் பெயர் மாற்றப்படுகிறது, அடுத்த புதுப்பிப்பு அதை நீக்குகிறது; முந்தைய `ccam.exe.old`-இலிருந்து ஓர் அமர்வு இன்னும் இயங்கிக்கொண்டிருந்தால், பெயர் `ccam.exe.old-<time>` ஆக இருக்கும். hard links ஆக உள்ள குறுக்குவழிகள் உடனே புதுப்பிக்கப்படும்.

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

`~/.local/share/ccam`-ஐ நீக்கினால், அங்கே சேமிக்கப்பட்ட உள்நுழைவுகளும் உரையாடல் வரலாறும் நீக்கப்படும். நீங்கள் [இறக்குமதி செய்த](/docs/commands) `~/.claude` கோப்புறை அது இருக்கும் இடத்திலேயே இருக்கும், ஆனால் அது ஒரு குழுவில் இருக்கும்வரை, அது பகிரும் கோப்புகள் `~/.local/share/ccam/groups/`-இல் இருக்கும். முதலில் அதற்கு `ccam remove <name>` இயக்குங்கள் (அல்லது `ccam group leave`): அது தான் பகிர்ந்த எல்லாவற்றின் நகலுடன் தன் குழுக்களிலிருந்து விலகும்.
