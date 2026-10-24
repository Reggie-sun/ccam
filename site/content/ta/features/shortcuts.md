---
title: குறுக்குவழிக் கட்டளைகள்
description: ccam shortcut மூலம் ஒரு கணக்குக்கு ccwork போன்ற தனிக் கட்டளையைக் கொடுங்கள், ccam shims install மூலம் ~/.local/bin-இல் உள்ள இணைப்புகளை மீண்டும் உருவாக்குங்கள்.
updated: 2026-10-23
order: 3
translation: machine
sourceHash: "38e2d1aba397"
---

குறுக்குவழி என்பது ஒரு கணக்காக Claude Code-ஐத் தொடங்க நீங்களே பெயரிடும் கட்டளை: `ccam work --resume` செய்வதையே `ccwork --resume` செய்யும். அது `~/.local/bin`-இல் உள்ள ccam-க்கான ஓர் இணைப்பு, அதனால் bash, zsh, scripts, `PATH`-இலிருந்து நிரல்களை இயக்கும் வேறு எதிலும் அது வேலை செய்யும். குறுக்குவழிகள் விருப்பத்தேர்வு மட்டுமே; நீங்கள் அடிக்கடி தொடங்கும் கணக்குகளுக்கு ஒன்றை உருவாக்குங்கள்.

## பயன்பாடு

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

ccam ஒவ்வொரு மாற்றத்தையும் உறுதிப்படுத்தும்:

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

ஒரு கணக்குக்கு அதிகபட்சம் ஒரு குறுக்குவழி மட்டுமே இருக்கும்; அது `~/.config/ccam/config.toml`-இல் `shortcut` ஆகச் சேமிக்கப்படும். `ccam add` wizard முதலில் `Create a shortcut command for <name>? [y/N]` என்றும், பிறகு `Shortcut name [cc<name>]:` என்றும் கேட்கும். wizard-க்கு வெளியே, `--shortcut` கொடுத்தால் மட்டுமே `ccam add` ஒன்றை உருவாக்கும். `ccam list` ஒவ்வொரு கணக்கின் குறுக்குவழியையும், இல்லையென்றால் `-`-ஐயும் காட்டும்.

`ccam <name>` பெயருக்குப் பின் உள்ளவற்றை அனுப்புவது போலவே, குறுக்குவழி ஒவ்வொரு argument-ஐயும் Claude Code-க்கு அனுப்பும். இதில் `--dry-run` மற்றும் `--help`-உம் அடங்கும்: ஒரு தொடக்கத்தை முன்னோட்டமாகப் பார்க்க `ccam --dry-run work` பயன்படுத்துங்கள்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--remove` | கணக்கின் குறுக்குவழியையும் அதன் இணைப்பையும் நீக்கும் | off |
| `-h`, `--help` | கட்டளைக்கான உதவி | |
| `--dry-run` | `ccam shortcut` இதை மதிக்கும்: பெயரைச் சரிபார்த்து, சேமிக்காமலும் இணைக்காமலும் `Would set the shortcut of work to ccwork.` அல்லது `Would remove ccwork, the shortcut of work.` என்று அச்சிடும் | off |

`ccam shims install` எந்த argument-ஐயும் ஏற்காது, `--dry-run`-ஐ `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` என்று மறுக்கும்.

பெயர் விதிகள், இந்த வரிசையில் சரிபார்க்கப்படும்:

- `a-z`, `0-9`, `.`, `-`, `_` ஆகியவற்றால் ஆன 1 முதல் 64 எழுத்துகள், ஓர் எழுத்து அல்லது இலக்கத்தில் தொடங்க வேண்டும்;
- `ccam`, `claude`, அல்லது `run`, `usage` போன்ற ccam கட்டளையாக இருக்கக் கூடாது;
- மற்றொரு கணக்கின் குறுக்குவழியாக இருக்கக் கூடாது;
- ஏற்கெனவே ccam-இன் இணைப்பாக (அல்லது ccam உருவாக்கிய, இந்தக் கணக்கின் சொந்தக் குறுக்குவழியாக) இல்லாத, `~/.local/bin`-இல் உள்ள கோப்பாக இருக்கக் கூடாது;
- உங்கள் `PATH`-இல் உள்ள வேறொரு executable-இன் பெயராக இருக்கக் கூடாது.

`config.toml` ஏற்றப்படும்போது இதே விதிகளுடன் சரிபார்க்கப்படுகிறது: pattern-ஐ மீறும், `claude` அல்லது `ccam` ஆக உள்ள, அல்லது இரண்டு கணக்குகளுக்குக் கொடுக்கப்பட்ட, கையால் திருத்திய குறுக்குவழி, அது சரிசெய்யப்படும் வரை, கணக்கைப் பெயரிடும் செய்தியுடன் ஒவ்வொரு கட்டளையையும் தோல்வியடையச் செய்யும்.

ஒரு கட்டளையையும் `--remove`-ஐயும் சேர்த்துக் கொடுத்தால் `give a command or --remove, not both` என்று நின்றுவிடும். ஒரு கணக்கிடம் ஏற்கெனவே உள்ள பெயரையே அமைத்தால் `work: ccwork (unchanged)` என்று அச்சிட்டு, இணைப்பு இல்லையென்றால் அதை மீண்டும் உருவாக்கும்.

## இது எப்படி வேலை செய்கிறது

பெயரை அமைக்கும்போது `checkShortcut` இயங்குகிறது; அது bin கோப்புறையையும் ccam binary-இன் உண்மையான பாதையையும் கொடுத்து `shim.Check`-ஐ அழைக்கிறது:

<!-- motion:shortcuts-check -->

பெயர் சரிபார்ப்புகளை வரிசையாகக் கடந்து செல்கிறது. ஒவ்வொரு சரிபார்ப்பின் அருகிலும் அது மறுக்கும் ஒரு பெயர் உள்ளது: `Work` pattern-இல் தோல்வியடைகிறது, `run` ஒரு ccam கட்டளை, `ccside` மற்றொரு கணக்குக்குச் சொந்தமானது, `ccfile` என்பது `~/.local/bin`-இல் ஏற்கெனவே உள்ள கோப்பு, `git` ஆனது `/usr/bin/git`-ஐ மறைத்துவிடும். எல்லாச் சரிபார்ப்புகளையும் கடக்கும் பெயர் மட்டுமே சேமிக்கப்பட்டு இணைக்கப்படும்.

1. `shim.Check` பெயரை முதலில் `^[a-z0-9][a-z0-9._-]{0,63}$`-உடனும், பிறகு `config.Reserved`, `cfg.FindShortcut` ஆகியவற்றுடனும் ஒப்பிடுகிறது. அது `~/.local/bin/<name>`-இல் `os.Lstat`-ஐ அழைக்கிறது, `PATH`-இல் உள்ள மற்ற ஒவ்வொரு கோப்புறைக்கும் `os.Stat`-ஐ அழைத்து, ccam-க்கான இணைப்பாக இல்லாத executable கோப்பை மறுக்கிறது.
2. ccam `fsutil.AtomicWrite` மூலம் (தற்காலிகக் கோப்பு, `fsync`, rename) `config.toml`-ஐச் சேமிக்கிறது.
3. `shim.Remove` பழைய குறுக்குவழியின் இணைப்பை நீக்குகிறது, ஆனால் அது ccam-க்கான இணைப்பாக இருந்தால் மட்டுமே.
4. `shim.Install` அமைக்கப்பட்ட ஒவ்வொரு குறுக்குவழியையும் பெயர் வரிசையில் கடந்து, இல்லாத ஒவ்வொன்றின் இணைப்பையும், `os.Executable` மற்றும் `filepath.EvalSymlinks` மூலம் கண்டறிந்த ccam binary-க்கு `os.Symlink` கொண்டு உருவாக்குகிறது. ஒரு குறுக்குவழியின் இடத்தில் உள்ள கோப்பு மற்றவற்றைத் தடுக்காது; பிழை ஒவ்வொன்றையும் பெயரிடுகிறது.

ccam தான் உருவாக்கும் ஒவ்வொரு இணைப்பையும், அது இயக்க வைத்த binary-ஐயும் `~/.local/state/ccam/shortcuts.json`-இல் பதிவு செய்கிறது. `os.Readlink` ccam பாதையைத் திருப்பித் தந்தால், அல்லது இணைப்பும் ccam பாதையும் ஒரே கோப்பைச் சுட்டினால், அந்த இணைப்பு ccam-இன் சொந்த இணைப்பாகக் கருதப்படும். பதிவு செய்யப்பட்ட இணைப்பின் binary அதன் பிறகு நகர்த்தப்பட்டிருந்தால் அல்லது மாற்றப்பட்டிருந்தால் (வேறொரு கோப்புறைக்கு நகர்த்தப்பட்ட ccam, மறு நிறுவல்), அது காலாவதியானது, `shim.Install` அதை இந்த ccam-க்குச் சுட்ட வைக்கிறது. இன்னும் நிறுவப்பட்டுள்ள வேறொரு ccam-க்கான பதிவு செய்யப்பட்ட இணைப்பு அப்படியே விடப்படும். ஒவ்வொரு நீக்கமும் ஒவ்வொரு மேலெழுதல் சரிபார்ப்பும் இந்தச் சோதனைகள் வழியாகச் செல்கின்றன, அதனால் ccam-இன் உங்கள் சொந்த நகல் அல்லது நீங்களே உருவாக்கிய இணைப்பு போன்ற, தான் உருவாக்காத கோப்புகளை ccam தொடுவதில்லை.

ஒரு குறுக்குவழியை இயக்கும்போது, நீங்கள் தட்டச்சு செய்த பெயர் மட்டுமே ccam-க்குத் தேவை:

<!-- motion:shortcuts-link -->

shell `PATH`-இல் `ccwork`-ஐக் கண்டுபிடிக்கிறது, kernel இணைப்பைப் பின்தொடர்ந்து ccam binary-ஐத் தொடங்குகிறது, `argv[0]` இன்னும் `ccwork` ஆகவே இருக்கிறது. `app.Main`-இல், `ccam` அல்லாத எந்தப் பெயரும் `config.toml`-ஐ ஒருமுறை ஏற்றி `shim.Lookup`-ஐ அழைக்கிறது; அது `shortcut` மதிப்பு `filepath.Base(argv[0])`-க்குச் சமமான கணக்கைக் கண்டுபிடிக்கிறது. பிறகு `launcher.LaunchWith`, ஏற்கெனவே ஏற்றப்பட்ட config-உடன் [`ccam <name>`](/features/launch) எடுக்கும் அதே படிகளை எடுத்து, `argv[1:]`-ஐ மாற்றாமல் `claude`-க்கு அனுப்புகிறது.

shell alias அல்லது wrapper script-க்குப் பதிலாக ஏன் symlink: symlink ஒவ்வொரு shell-இலும் scripts-இலும் வேலை செய்யும், கூடுதல் shell தொடக்கத்தைச் (சுமார் 1 முதல் 2 ms) சேர்க்காது, தொடக்கம் இன்னும் ஒரே ஒரு கோப்பைப் படிக்கும். `scripts/perf-check.sh` ஒரு குறுக்குவழியை `hyperfine` மூலம், `ccam <name>`-க்கான அதே 10 ms (p95) வரம்புக்கு எதிராக அளவிடுகிறது.

### Windows-இல்

குறுக்குவழி என்பது `%USERPROFILE%\.local\bin`-இல் உள்ள `ccwork.exe`, `shim.Name` தேடலுக்கு முன் `.exe`-ஐ (எந்த எழுத்து வடிவிலும்) நீக்குகிறது, அதனால் `ccwork`, `ccwork.exe`, `CCWORK.EXE` அனைத்தும் `work`-ஐத் தொடங்கும். ccam முதலில் symbolic link-ஐ உருவாக்குகிறது, அதை Windows, Developer Mode-இல் அல்லது administrator-க்கு மட்டுமே அனுமதிக்கிறது. அந்த உரிமை இல்லாமல் அது hard link-ஐ உருவாக்குகிறது: அதே `ccam.exe` கோப்புக்கான இரண்டாவது பெயர், அதற்கு எந்தச் சிறப்புரிமையும் தேவையில்லை. ஒரு hard link `ccam.exe`-இன் அதே கோப்பாக இருக்கும்போது அது ccam-இன் சொந்தமாகக் கருதப்படுகிறது. hard link தான் உருவாக்கப்பட்ட கோப்பையே தொடர்ந்து இயக்குகிறது, அதனால் `ccam update` `ccam.exe`-இன் இடத்தில் புதிய கோப்பை வைக்கும்போது, அது ஒவ்வொரு hard-link குறுக்குவழியையும் `shim.Relink` மூலம் புதுப்பித்து `✓ shortcut ccwork renewed` என்று அச்சிடுகிறது. `install.ps1`-ஐ மீண்டும் இயக்குவதும் `ccam shims install` வழியாக இதையே செய்கிறது (`OK 2 shortcuts now run the new ccam`); பழைய ccam-இன் நகலாக உள்ள எந்தக் குறுக்குவழிக்கும் `ccam shims install`-உம் `ccam doctor --fix`-உம் அப்படியே செய்கின்றன. பெயர்ச் சரிபார்ப்பு மற்ற `PATH` கோப்புறைகளில் `PATHEXT`-இல் உள்ள ஒவ்வொரு extension-உடனும் (`.exe`, `.cmd` போன்றவை) `<name>`-ஐத் தேடுகிறது.

குறுக்குவழியை ஒரு shell-இல் (PowerShell, cmd.exe, Git Bash) தட்டச்சு செய்யுங்கள். Explorer, Run dialog அல்லது `Start-Process`-இலிருந்து திறக்கப்படும் symbolic link வெறும் `ccam` ஆகத் தொடங்கும், அது கணக்குக்குப் பதிலாக டாஷ்போர்டைத் திறக்கும்.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- **எந்தக் கணக்கும் உரிமை கோராத இணைப்பு.** `~/.local/bin`-இல் உள்ள ccam-க்கான ஓர் இணைப்பின் பெயரை எந்தக் கணக்கும் பயன்படுத்தவில்லை என்றால், எடுத்துக்காட்டாக `config.toml`-இல் `shortcut`-ஐ நீங்கள் கையால் திருத்திய பிறகு, ccam இயல்புநிலைக் கணக்குக்கு மாறாது. அது `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` என்று நின்றுவிடும். `ccam shortcut` மூலம் பெயரை மீண்டும் அமையுங்கள், அல்லது இணைப்பை நீக்குங்கள்.
- **ஏற்கெனவே உள்ள கோப்புகள்.** தான் உருவாக்காத கோப்பை ccam ஒருபோதும் மாற்றுவதில்லை. `ccam shims install` `<path> already exists and is not a ccam link, so it was left alone` என்று தெரிவித்து, மற்ற இணைப்புகளை அப்போதும் உருவாக்கும். அப்போது `ccam add`-உம் `ccam import`-உம் `shortcut: ccwork (not created; see above)` என்று அச்சிடும்.
- **நகர்த்தப்பட்ட ccam binary.** அப்போது ccam உருவாக்கிய இணைப்புகள் இனி இல்லாத பாதையைச் சுட்டும். `ccam shims install` (அல்லது `ccam doctor --fix`) நீங்கள் அதை இயக்கும் ccam-க்கு அவற்றைச் சுட்ட வைக்கிறது.
- **பெயர் மாற்றமும் நீக்கமும்.** `ccam rename` குறுக்குவழிப் பெயரை வைத்திருக்கும், இணைப்பும் தொடர்ந்து வேலை செய்யும், ஏனெனில் அது ccam-ஐச் சுட்டுகிறது, ccam பெயரை `config.toml`-இல் தேடுகிறது. `ccam remove` கணக்கின் இணைப்பை நீக்கும். `ccam shortcut <account> <new>` புதியதை உருவாக்கும் முன் பழைய இணைப்பை நீக்கும்.
- **PATH.** ccam ஓர் இணைப்பை உருவாக்கும்போது `~/.local/bin` உங்கள் `PATH`-இல் இல்லையென்றால், அது `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` என்றும், அதைத் தொடர்ந்து சேர்க்க வேண்டிய `export PATH=...` வரியையும் அச்சிடும். Windows-இல் அதற்குப் பதிலாக, கோப்புறையை உங்கள் பயனர் `PATH`-இல் சேர்க்கும் PowerShell வரியை அச்சிடும்.
- **பின்னர் சரிபார்த்தல்.** `ccam doctor` `missing shortcuts in /home/me/.local/bin: ccwork` என்றும், பழைய ccam-ஐ இயக்கும் குறுக்குவழிகளையும் தெரிவிக்கும், கோப்புறை `PATH`-இல் இல்லையென்றால் எச்சரிக்கும். `ccam doctor --fix` இல்லாத இணைப்புகளை உருவாக்கி, காலாவதியானவற்றைப் புதுப்பிக்கும், தான் உருவாக்கியவற்றுக்கு மட்டுமே `→ fixed` என்று சொல்லும். இன்னும் நிறுவப்பட்டுள்ள வேறொரு ccam-ஐ இயக்கும் குறுக்குவழி வேலை செய்யும், அது தகவல் மட்டுமே; `ccam doctor --fix` நீங்கள் அதை இயக்கும் ccam-க்கு அதைச் சுட்ட வைக்கிறது, `ccam add`-உம் `ccam import`-உம் அதை ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)`` எனக் காட்டும்.

## தொடர்புடையவை

- [Claude Code-ஐத் தொடங்குதல்](/features/launch)
- [கணக்குகளை நிர்வகித்தல்](/features/accounts)
- [கட்டளைக் குறிப்பு](/docs/commands#shortcuts)
- [ccam-ஐ நிறுவுதல்](/docs/install)
