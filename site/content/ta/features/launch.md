---
title: Claude Code-ஐத் தொடங்குதல்
description: ccam <name|number>, ccam run அல்லது ccam தேர்வி மூலம் எந்தக் கணக்காகவும் Claude Code-ஐத் தொடங்குங்கள்; arguments மாறாமல் claude-க்குச் செல்லும்.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work` Claude Code-ஐ `work` கணக்காகத் தொடங்குகிறது, பெயருக்குப் பின் உள்ள அனைத்தும் அப்படியே `claude`-க்குச் செல்லும். நீங்கள் `claude` என்று தட்டச்சு செய்யும் எந்த இடத்திலும் இதைப் பயன்படுத்துங்கள்: டெர்மினலில், script-இல், அல்லது pipe மூலம் உள்ளீடு கொடுக்கும்போது.

## பயன்பாடு

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

ஒரு கணக்கின் பெயரையோ `ccam list`-இல் உள்ள அதன் எண்ணையோ கொடுங்கள். பெயருக்குப் பின் உள்ள arguments-ஐ ccam படிப்பதில்லை, அதனால் `ccam work --help` Claude Code-இன் உதவியைக் காட்டும். பெயருக்குப் பின் `--` விருப்பத்தேர்வு மட்டுமே. standard input, exit code, signals ஆகியவை `claude`-இல் உள்ளது போலவே செயல்படும்.

முதல் argument `-` என்று தொடங்கி, ccam-இன் சொந்த flags-இல் ஒன்றாக இல்லாதபோது, ccam எல்லா arguments-உடனும் பெயர் இல்லாமல் ஒரு கணக்கைத் தொடங்கும். அது [`ccam folder`](/features/folder) மூலம் தற்போதைய கோப்புறைக்கு அமைக்கப்பட்ட கணக்கு, இல்லையெனில் இயல்புநிலை: `ccam use` மூலம் அமைத்த கணக்கு, அல்லது நீங்கள் ஒன்றை அமைக்கவில்லை என்றால் மிகக் குறைந்த எண் உள்ள கணக்கு.

### தேர்வி

standard input-உம் output-உம் டெர்மினலாக இருக்கும்போது, arguments இல்லாத `ccam` தேர்வியைத் திறக்கும், இல்லையென்றால் `ccam list`-ஐ அச்சிடும். அது ஒவ்வொரு கணக்கின் எண், மின்னஞ்சல், 5 மணி நேரப் பயன்பாடு, கடைசிப் பயன்பாடு, பகிரப்பட்ட உருப்படிகள் ஆகியவற்றைக் காட்டும்; ஒவ்வொரு கணக்கும் பதில் அளிக்கும்போது பயன்பாடு நிரப்பப்படும்.

| விசை | என்ன செய்யும் |
| --- | --- |
| `Enter` | தேர்ந்தெடுத்த கணக்காக Claude Code-ஐத் தொடங்கும் |
| `1` முதல் `9` வரை | அந்த எண் உள்ள கணக்கைத் தொடங்கும் |
| `Up`/`Down`, `k`/`j` | தேர்வை நகர்த்தும் |
| `a` | தேர்வியை மூடி `ccam add` wizard-ஐ இயக்கும் |
| `s` | தேர்ந்தெடுத்த கணக்கின் பகிரப்பட்ட உருப்படிகளைத் திருத்தும் (`Space` அல்லது `x` மாற்றும், `Enter` பொருத்தும், `Esc` பின்செல்லும்) |
| `d` | தேர்ந்தெடுத்த கணக்கை இயல்புநிலை ஆக்கும் |
| `r` | குறுகிய cache இல்லாமல், பயன்பாட்டை மீண்டும் பெறும் |
| `q`, `Esc`, `Ctrl+C` | வெளியேறும் |

தேர்வி செயல்படும் முன் மூடிக்கொள்கிறது, அதனால் Claude Code தொடங்கும் முன் டெர்மினல் மீட்டமைக்கப்படும். `d` அல்லது பகிர்வு மாற்றத்துக்குப் பிறகு அது மீண்டும் திறக்கும். பகிர்வு மாற்றம், `ccam share` செய்வது போலவே, முதலில் திட்டத்தைக் காட்டிக் கேட்கும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--dry-run` | ccam இயக்கும் சூழலையும் கட்டளையையும் அச்சிட்டு வெளியேறும் | off |
| `--fallback` | `-p`/`--print` இயக்கங்களுக்கு: கணக்கு தன் பயன்பாட்டு வரம்பை எட்டும்போது, அதே கட்டளையை அதிக இடம் உள்ள உள்நுழைந்த கணக்காக மீண்டும் இயக்கும் ([விவரங்கள்](/features/next#run-again-at-the-limit)) | off |
| `-v`, `--version` | `ccam <version>`-ஐ அச்சிடும் (முதல் argument ஆக மட்டும்) | |
| `-h`, `--help` | `ccam --help` அல்லது `ccam run --help` ஆக, ccam-இன் உதவி | |

Argument விதிகள்:

- ccam தன் சொந்த flags-ஐக் கணக்குப் பெயருக்கு முன் மட்டுமே படிக்கிறது. பெயருக்குப் பின், `--dry-run` மற்றும் `--help` உட்பட, ஒவ்வொரு argument-உம் Claude Code-க்கு உரியது.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run` (பெயர் இல்லாத கணக்கு) ஆக வேலை செய்யும். அது மதிக்கப்படும், மறுக்கப்படாது. `--dry-run`, `--fallback` எந்த வரிசையிலும் வரலாம்.
- `ccam run <name|number> [claude args]` என்பது `ccam <name>`-இன் நீண்ட வடிவம். பெயர் இல்லாத `ccam run` கோப்புறையின் கணக்கை அல்லது இயல்புநிலைக் கணக்கைப் பயன்படுத்தும்.
- கணக்குப் பெயர்கள் `add`, `list`, `run`, `usage` போன்ற ccam கட்டளைப் பெயர்களாக இருக்க முடியாது, அதனால் `ccam <word>` ஒருபோதும் குழப்பமாக இருக்காது.

dry run, shell-இல், அல்லது Windows-இல் PowerShell-இல் ஒட்டக்கூடிய ஒரு வரியை அச்சிடும்:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## இது எப்படி வேலை செய்கிறது

ஒவ்வொரு தொடக்கமும், கட்டளை வரியிலிருந்து `exec` வரை, அதே படிகளை எடுக்கிறது:

<!-- motion:launch -->

ஒவ்வொரு படியும் இயங்கும்போது அதன் பெட்டி ஒளிரும், கீழே உள்ள வரிகள் Claude Code தொடங்கும் சூழலும் arguments-உம் ஆகும்.

1. `cmd/ccam/main.go` `app.Main`-ஐ அழைக்கிறது. ccam வேறொரு பெயரில் தொடங்கப்பட்டிருந்தால், அது ஒரு [குறுக்குவழி](/features/shortcuts), நேரடியாக launcher-க்குச் செல்லும்.
2. cobra ஈடுபடும் முன் `app.parse` `os.Args`-ஐப் பார்க்கிறது. தனியாக `ccam` (தேர்வி), ஒரு ccam கட்டளைப் பெயர் (`config.Reserved`), `-h`/`--help`, `-v`/`--version` ஆகியவை மட்டுமே `cli.Execute`-க்குச் செல்லும். மற்ற அனைத்தும் தொடக்கமே.
3. `launcher.Launch` `~/.config/ccam/config.toml`-ஐப் படிக்கிறது; அது parse செய்யும் ஒரே கோப்பு இதுதான்.
4. `launcher.Prepare` கணக்கைப் பெயராலோ எண்ணாலோ கண்டுபிடித்து (பெயர் இல்லையென்றால் `HereProfile` கோப்புறையின் கணக்கை எடுக்கிறது, working directory-ஐ `config.toml`-இன் `[[folder]]` பதிவுகளுடன் strings ஆக ஒப்பிட்டு), முன்னால் உள்ள ஒரு `--`-ஐ நீக்கி, `claude` executable-ஐ `ResolveClaude` மூலம் கண்டறிகிறது. PATH தேடல் கோப்பின் மாற்ற நேரத்துடன் சேர்த்து `~/.local/state/ccam/claude-path`-இல் cache செய்யப்படுகிறது, அதனால் ஒவ்வொரு தொடக்கத்திலும் ccam `PATH`-ஐத் தேடுவதில்லை. மாறிய அல்லது இல்லாத கோப்பு புதிய தேடலைத் தூண்டும். `config.toml`-இல் உள்ள `claude_bin` பாதை அப்படியே பயன்படுத்தப்படும். ccam-ஆகவே இருக்கும் `claude`-ஐ, inode மூலம் ஒப்பிட்டு, ccam மறுக்கிறது, ஏனெனில் அது முடிவில்லாச் சுழற்சியை ஏற்படுத்தும்.
5. `checkLinks` கணக்குக் கோப்புறை `~/.local/share/ccam/profiles/<name>/`-ஐ அதன் பகிர்வு அமைப்புகளுடன் `share.Check` மூலம் ஒப்பிடுகிறது; அது `lstat` மற்றும் `readlink`-ஐ மட்டுமே அழைக்கிறது, ஒவ்வோர் உருப்படிக்கும் ஒன்று அல்லது இரண்டு system calls. இல்லாத இணைப்பு அந்த இடத்திலேயே உருவாக்கப்படும். தரவை இழக்கக்கூடிய எதுவும் தொடக்கத்தை நிறுத்தும்.
6. `applyMCP` பகிரப்பட்ட MCP சர்வர்களைச் சேர்க்கிறது (கீழே பாருங்கள்).
7. `session.Register` அமர்வைப் பதிவு செய்கிறது. ccam தன் PID-இன் பெயரில், process தொடக்க நேரத்தைக் கொண்ட ஒரு குறிப்பானை `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`-இல் (அல்லது `XDG_RUNTIME_DIR` இல்லாவிட்டால் `~/.local/state/ccam/locks/`-இல்) எழுதுகிறது. உரையாடல் வரலாற்றைப் பகிரும் கணக்குக்கு, `historyCheck` முதலில் `history-<group>-<hash>`-இல் இயங்கும் அமர்வைத் தேடுகிறது, இங்கே hash தற்போதைய கோப்புறையினுடையது, அங்கேயும் பதிவு செய்கிறது. `exec` PID-ஐ வைத்திருக்கிறது, அதனால் குறிப்பான் பின்னர் Claude Code process-ஐயே பெயரிடுகிறது. பிடித்து வைத்த lock அல்ல, குறிப்பான், ஏனெனில் Claude Code தொடங்கும் processes ஒரு lock-ஐப் பெற்று, அமர்வு முடிந்த பிறகும் அதை வைத்திருக்கும்.
8. `Plan.Exec` டெர்மினல் தலைப்பை `claude · <name>` என அமைக்கிறது (standard error டெர்மினலாக இருக்கும்போது), தொடக்க நேரத்தை `fsync` இல்லாமல் `~/.local/state/ccam/last-used/<name>`-இல் எழுதுகிறது, `syscall.Exec`-ஐ அழைக்கிறது.

Claude Code, கணக்குக் கோப்புறைக்கு அமைக்கப்பட்ட `CLAUDE_CONFIG_DIR` உடனும் கணக்குப் பெயருக்கு அமைக்கப்பட்ட `CCAM_PROFILE` உடனும் இயங்குகிறது. இவற்றுக்கு மரபுரிமையாகப் பெற்ற மதிப்புகள் எதையும் ccam நீக்குகிறது. macOS-இல் `ccam import` மூலம் பதிவுசெய்யப்பட்ட `~/.claude` கோப்புறைக்கு, அது `CLAUDE_SECURESTORAGE_CONFIG_DIR`-ஐ வெற்று மதிப்புக்கும் அமைக்கிறது, அதனால் அந்தக் கணக்கும் வெறும் `claude`-உம் ஒரே Keychain உள்நுழைவைப் பயன்படுத்தும்.

### ஏன் exec

`syscall.Exec` ccam process-ஐ `claude`-ஆல் மாற்றுகிறது, அதனால் ccam நினைவகத்தில் தங்குவதில்லை, விசை அழுத்தங்கள், signals, exit code ஆகியவை Claude Code-இன் சொந்தமானவை, arguments ஒருபோதும் shell string ஆகாமல் array ஆகவே இருக்கும்.

Windows-இல் `exec` இல்லை. அங்கே `Plan.Exec` அதே சூழல், arguments, console-உடன் `claude.exe`-ஐ child process ஆகத் தொடங்கி, அதற்காகக் காத்திருந்து, அதன் exit code-உடன் வெளியேறுகிறது. Ctrl+C console-இலிருந்து நேரடியாக Claude Code-ஐ அடைகிறது, அதனால் ccam அதைப் புறக்கணிக்கிறது. அமர்வு marker ccam-இன் PID-ஐக் கொண்டிருக்கிறது, அது சரியாக Claude Code இருக்கும் வரை இருக்கும், Claude Code வெளியேறும்போது marker நீக்கப்படுகிறது. காத்திருக்கும்போது ccam வேறு எதுவும் செய்வதில்லை. அமர்வுக்குப் பிந்தைய வேலை, எடுத்துக்காட்டாக `claude mcp add` மூலம் சேர்க்கப்பட்ட சர்வர்களைக் கவனிப்பது, அடுத்த தொடக்கத்தில் நடக்கும்: ccam `.claude.json`-இன் மாற்ற நேரத்தை `~/.local/state/ccam/mcp/<name>/seen`-உடன் ஒப்பிட்டு, அது மாறியிருந்தால் மட்டுமே கோப்பை parse செய்கிறது.

### கட்டளை வரி எப்படி வழிநடத்தப்படுகிறது

மூன்று கட்டளை வரிகள், `internal/app` வழியாக மூன்று வழிகள்:

<!-- motion:launch-route -->

குறுக்குவழி `argv[0]`-இலும், `ccam <name>` முதல் argument-இலும் பொருத்தப்படுகின்றன, இரண்டும் `syscall.Exec`-இல் முடிகின்றன. `ccam usage` போன்ற ccam கட்டளை மட்டுமே cobra-வை அடைகிறது.

### பகிரப்பட்ட MCP சர்வர்கள்

இயல்புநிலை `flag` உத்தியுடன், ccam குழுவின் MCP சர்வர்களை உங்களுடையவற்றுக்கு முன் வைக்கப்படும் ஒரு argument ஆக Claude Code-க்கு அனுப்புகிறது:

<!-- motion:launch-mcp -->

குழுவின் `mcp.json` இருந்தால் மட்டுமே `mcp.FlagArgs` ஒரு argument-ஐத் திருப்பித் தருகிறது, `Plan.InsertArgs` அதை `claude`-க்கு உடனே பின் வைக்கிறது. அப்போதும் உங்கள் சொந்த `--mcp-config` அல்லது `--strict-mcp-config` அதற்கு மேல் பொருந்தும். அது `--mcp-config=<path>` என்ற ஒரே argument, ஏனெனில் Claude Code-இன் `--mcp-config` பல மதிப்புகளை ஏற்கும்: இரண்டு arguments ஆக எழுதினால், அதைத் தொடரும் prompt-ஐயும் config கோப்பாகப் படித்துவிடும்.

`config.toml`-இல் `mcp_strategy = "sync"` இருந்தால், ccam அதற்குப் பதிலாகச் சர்வர்களைக் கணக்கின் `.claude.json`-இல் எழுதும், கணக்கு ஏற்கெனவே இயங்கினால் அதைத் தவிர்க்கும்.

### Remote Control

`ccam work rc` (அல்லது `ccam work remote-control`) Claude Code-இன் Remote Control-ஐ `work` கணக்காகத் தொடங்குகிறது. `claude remote-control` தன் verb-க்கு முன் `--mcp-config`-ஐ மறுக்கிறது, அதற்குப் பின் எதையும் ஏற்பதில்லை; அதன் அமர்வுகள் MCP சர்வர்களைக் கணக்கின் `.claude.json`-இலிருந்து மட்டுமே படிக்கின்றன. அதனால் இங்கே flag உத்தி வேறு விதமாகச் செயல்படுகிறது:

<!-- motion:launch-rc -->

1. `remoteControl` கணக்குக்குப் பின் `-` எனத் தொடங்காத முதல் argument-ஐ Claude Code-இன் verb ஆக எடுத்துக்கொள்கிறது. அது `remote-control` அல்லது `rc` ஆக இருந்தால், `mcp.FlagArgs` எதையும் சேர்க்காது.
2. `mcp.Synced` குழுவின் `mcp.json` மாற்ற நேரத்தைக் கணக்கின் `sync-stamp`-உடன் ஒப்பிடுகிறது. `.claude.json`-இல் ஏற்கெனவே குழுவின் தற்போதைய சர்வர்கள் இருந்தால், அல்லது குழுவிடம் சர்வர்களே இல்லையென்றால், Remote Control கேள்வி இல்லாமல் தொடங்கும்.
3. இல்லையெனில், அமர்வுகளுக்குக் கிடைக்காமல் போகும் சர்வர்களை ccam பெயரிட்டுக் கேட்கும்:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` என்றால், `mcp.Sync` அவற்றை `sync` உத்தி செய்வது போலவே நகலெடுக்கிறது: Claude Code-இன் lock-இன் கீழ், காப்புப்பிரதி எடுத்த பிறகு, கணக்கின் சொந்தச் சர்வர்களை வைத்துக்கொண்டு. `N` என்றால், அல்லது டெர்மினல் இல்லாவிட்டால், Remote Control அவை இல்லாமல் தொடங்கும், ccam ஒவ்வொரு தொடக்கத்துக்கு முன்பும் அவற்றை நகலெடுக்கும் `mcp_strategy = "sync"`-ஐச் சுட்டிக்காட்டும்.

சர்வர்கள் நகலெடுக்கப்படும் வரை ஒவ்வொரு Remote Control தொடக்கத்திலும் ccam கேட்கும்; நகலெடுத்த பிறகு, குழுவின் `mcp.json` மீண்டும் மாறும்போது மட்டுமே கேட்கும். பின்வரும் சாதாரணத் தொடக்கங்களில் `.claude.json`-இல் உள்ள நகல் `--mcp-config`-ஐ விட முன்னுரிமை பெறும், அதனால் `ccam mcp add`, `remove`, `import` ஆகியவை மற்ற கணக்குகள் வைத்திருக்கும் நகல்களைப் புதுப்பிக்கின்றன ([பகிரப்பட்ட MCP சர்வர்கள்](/features/mcp)).

### 10 ms வரம்பு

Claude Code தொடங்கும் முன் ccam அதிகபட்சம் 10 ms (p95) மட்டுமே சேர்க்கிறது. தொடக்கப் பாதை ஒருபோதும் cobra கட்டளை மரத்தை உருவாக்குவதோ, Bubble Tea-வைத் தொடங்குவதோ, network இணைப்பைத் திறப்பதோ இல்லை, package `init()` functions மொத்தம் 5 ms-க்கு மேல் எடுத்தால் `TestPackageInitIsCheap` தோல்வியடையும். `scripts/perf-check.sh` `ccam --dry-run work`, `ccam work`, ஒரு குறுக்குவழி ஆகியவற்றை `hyperfine` மூலம் அளந்து, 10 ms-க்கு மேல் இருந்தால் தோல்வியடைகிறது, `internal/launcher`-இல் உள்ள `BenchmarkLaunchPath` படிகளை ஒரு loop-இல் இயக்குகிறது. ஒவ்வொரு படியின் நேரத்தையும் பார்க்க, `CCAM_TRACE=1` அமையுங்கள்:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- தெரியாத பெயர் அல்லது எண் `no account "<name>"` என்று நின்று `ccam list`-ஐச் சுட்டிக்காட்டும். கணக்குகளே இல்லையென்றால், `ccam add <name>` மூலம் ஒன்றைச் சேர்க்க ccam கேட்கும்.
- பகிரப்பட்ட பாதை சாதாரணக் கோப்பாக மாறியிருந்தால் அல்லது தவறான இடத்தைச் சுட்டினால், ccam `stopped because of share link problems:` என்று நின்று, அவற்றைப் பட்டியலிட்டு, `ccam doctor --fix`-ஐ இயக்கச் சொல்லும். தொடக்கப் பாதையில் அது தரவை நகர்த்துவதில்லை.
- இந்தக் கணக்குடன் உரையாடல் வரலாற்றைப் பகிரும் (அதே குழு) மற்றொரு Claude Code அதே கோப்புறையில் இயங்கினால், வரலாறு கலக்கக்கூடும் என்று ccam எச்சரித்து, `Run anyway? [y/N]` என்று கேட்கும். டெர்மினல் இல்லாவிட்டால் எச்சரிக்கையைக் காட்டித் தொடரும்.
- `--dry-run` MCP sync, புதிய சர்வர் அறிவிப்பு, வரலாற்றுச் சரிபார்ப்பு, அமர்வுக் குறிப்பான், தலைப்பு, கடைசிப் பயன்பாட்டு நேரம் ஆகியவற்றைத் தவிர்க்கும், exec செய்யாது. அது இன்னும் பகிர்வு இணைப்புகளைச் சரிபார்க்கும், சாதாரணத் தொடக்கம் செய்வது போலவே இல்லாத இணைப்பை உருவாக்கும்.
- `ccam work rc`, `work` ஏற்கெனவே இயங்கிக்கொண்டிருக்கும்போது, `.claude.json`-ஐத் தொடாமல் விட்டு, அதைச் சொல்லி, குழுவின் சர்வர்கள் இல்லாமல் Remote Control-ஐத் தொடங்கும்; அவற்றை நகலெடுக்க அமர்வு முடிந்த பிறகு மீண்டும் இயக்குங்கள். `--dry-run` உடன் அது செய்தியை அச்சிடும், ஆனால் கேட்கவோ நகலெடுக்கவோ செய்யாது.
- `mcp_strategy = "sync"` உடன், தோல்வியடைந்த sync `ccam: MCP sync failed (launching anyway): …` என்று அச்சிடும், Claude Code ஏற்கெனவே இருந்த சர்வர்களுடன் தொடங்கும்.

## தொடர்புடையவை

- [குறுக்குவழிக் கட்டளைகள்](/features/shortcuts)
- [ccam உடன் Remote Control-ஐப் பயன்படுத்துதல்](/guides/remote-control)
- [கணக்குகளை நிர்வகித்தல்](/features/accounts)
- [கட்டளைக் குறிப்பு](/docs/commands)
- [MCP சர்வர்கள், அமைப்புகள், வரலாற்றைப் பகிர்தல்](/docs/sharing)
- [CLAUDE_CONFIG_DIR எப்படி வேலை செய்கிறது](/guides/claude-config-dir)
