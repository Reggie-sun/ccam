---
title: ccam-ஐப் புதுப்பித்தல்
description: ccam update SHA-256 சரிபார்ப்புக்குப் பின் ccam-ஐச் சமீபத்திய GitHub வெளியீட்டால் மாற்றுகிறது; --check தெரிவிக்கும், --version tag-ஐ நிறுவும்.
updated: 2026-10-20
order: 10
translation: machine
sourceHash: d91499b95d3d
---

`ccam update` GitHub-இலிருந்து ஒரு ccam வெளியீட்டைப் பதிவிறக்கி, அதை வெளியீட்டின் `checksums.txt`-உடன் சரிபார்த்து, இயங்கும் binary-ஐ ஒரே rename மூலம் மாற்றுகிறது. install script மூலமாகவோ ஓர் archive-இலிருந்தோ ccam-ஐ நிறுவியிருந்தால் இதைப் பயன்படுத்துங்கள். ccam தானாகப் புதுப்பிப்புகளைச் சரிபார்ப்பதில்லை: இந்தக் கட்டளை மட்டுமே GitHub-ஐத் தொடர்புகொள்கிறது, அதனால் மற்ற எல்லாக் கட்டளைகளும் இணையம் இல்லாமலும் அதே போல வேலை செய்யும்.

## பயன்பாடு

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

புதிய வெளியீடு எதுவும் இல்லையென்றால், ccam `ccam 0.4.0 is up to date.` என்று அச்சிட்டு, எதையும் மாற்றாது. நீங்கள் ஏற்கெனவே இயக்கும் பதிப்பை `--version` உடன் கொடுத்தால், அது `ccam 0.4.0 is already installed.` என்று அச்சிடும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--check` | புதிய வெளியீடு உள்ளதா என்று மட்டும் தெரிவிக்கும் | off |
| `--version <tag>` | `v0.3.0` போன்ற இந்த வெளியீட்டு tag-ஐ நிறுவும், அது தற்போதையதை விடப் புதியதாகவோ பழையதாகவோ இருக்கலாம் | latest |
| `--dry-run` | மதிக்கப்படும்: சமீபத்திய tag-ஐத் தேடி (`--version` ஒன்றைக் குறிப்பிடாவிட்டால்), binary-இல் அதே சரிபார்ப்புகளை இயக்கி, பதிவிறக்காமலே download, checksum, replace பாதைகளை அச்சிடும் | off |

`--check`, `--version` இரண்டையும் சேர்த்துப் பயன்படுத்த முடியாது. tag `vMAJOR.MINOR.PATCH` வடிவத்தில் இருக்க வேண்டும், விருப்பத்தேர்வாக ஒரு pre-release அல்லது build பின்னொட்டுடன்; முன்னால் உள்ள `v`-ஐ விட்டுவிடலாம். வேறு எதுவும் தோல்வியடையும், எடுத்துக்காட்டாக `--version takes a release tag such as v0.3.0, not "1.2"`.

`CCAM_RELEASES_URL` சூழல் மாறி, ஒரு mirror-க்காகவோ சோதனைகளுக்காகவோ, `https://github.com/jhyunwoo/ccam/releases`-ஐ மாற்றுகிறது. அது அதே அமைப்பில் இருக்க வேண்டும்: `<url>/latest` ஆனது `<url>/tag/<tag>`-க்குத் திருப்பிவிட வேண்டும், assets `<url>/download/<tag>/<asset>`-இல் இருக்க வேண்டும்.

## இது எப்படி வேலை செய்கிறது

குறியீடு `internal/update`-இல் உள்ளது, அதை `ccam update` மட்டுமே import செய்கிறது. தொடக்கப் பாதை அதை ஒருபோதும் அழைப்பதில்லை.

1. `cli.updateCmd` இயங்கும் binary-ஐ `os.Executable` மூலம் கண்டுபிடித்து, symlinks-ஐ `filepath.EvalSymlinks` மூலம் தீர்க்கிறது. அந்த உண்மையான பாதையே மாற்றப்படுகிறது.
2. `--version` இல்லாதபோது, `Updater.Latest` redirect-ஐப் பின்தொடராமல் `<releases>/latest`-ஐக் கோருகிறது, `Location` header-இலிருந்து (`…/tag/v0.5.0`) tag-ஐப் படிக்கிறது. இது GitHub API-க்குப் பதிலாக வெளியீடுகள் பக்கத்தைப் பயன்படுத்துகிறது, அதனால் மணிக்கு 60 கோரிக்கைகள் என்ற API வரம்பு பொருந்தாது. இந்தத் தேடலுக்கு 10 வினாடிகள் உள்ளன.
3. `parseVersion` tag-ஐ ஒரு கண்டிப்பான pattern-உடன் சரிபார்க்கிறது. tag பதிவிறக்க URLs-இலும் அச்சிடப்படும் shell கட்டளையிலும் இடம்பெறுகிறது, அதனால் `/`, `%`, இடைவெளிகள், shell எழுத்துகள் உள்ள எதுவும் உள்ளே வராது. `compareVersions` அதை இயங்கும் பதிப்புடன் ஒப்பிடுகிறது; ஒரு pre-release அதன் வெளியீட்டுக்கு முன்னால் வரிசைப்படுத்தப்படும்.
4. ccam binary-ஐ மாற்றலாமா என்பதை `managed()` மற்றும் `writable()` முடிவு செய்கின்றன (கீழே பாருங்கள்).
5. `Updater.Fetch` `ccam_<os>_<arch>.tar.gz` (Windows-இல் `ccam_windows_<arch>.zip`) மற்றும் `checksums.txt`-ஐ ஒரே நேரத்தில், இரண்டுக்கும் சேர்த்து 2 நிமிடங்களுக்குள், பதிவிறக்குகிறது. ஒவ்வொரு பதிவிறக்கமும் 256 MB-க்குள் வரம்பிடப்படுகிறது.
6. வேறு எதுவும் நடப்பதற்கு முன், archive-இன் SHA-256 `checksums.txt`-இல் உள்ள அதன் வரியுடன் பொருந்த வேண்டும். வரி இல்லையென்றாலோ checksum வேறுபட்டாலோ, புதுப்பிப்பு `nothing was changed` என்று நின்றுவிடும்.
7. `extract()` tar-ஐப் படித்து, `ccam` என்ற பெயருடைய முதல் சாதாரணக் கோப்பை எடுக்கிறது (`extractZip` zip-இலிருந்து `ccam.exe`-ஐ எடுக்கிறது). கோப்புறைகள், இணைப்புகள், பிற உள்ளீடுகள் புறக்கணிக்கப்படும், binary-உம் 256 MB-க்குள் வரம்பிடப்படும்.
8. `Updater.Install` புதிய binary-ஐ `fsutil.AtomicWrite` மூலம் எழுதுகிறது: binary-க்கு அருகில் ஒரு தற்காலிகக் கோப்பு `.ccam.tmp*`, பழைய கோப்பின் அனுமதி bits-உடன், வட்டில் flush செய்யப்பட்டு, பழையதன் மேல் rename செய்யப்படுகிறது.

tag தேடலிலிருந்து rename வரை, முழு pipeline:

<!-- motion:update-pipeline -->

archive-உம் `checksums.txt`-உம் அருகருகே பதிவிறக்கப்படுகின்றன, checksum பொருந்தும் வரை எதுவும் எழுதப்படுவதில்லை. rename atomic ஆனது, அதனால் crash அல்லது நிரம்பிய வட்டு பழைய binary-ஐயோ புதியதையோ விட்டுச்செல்லும், ஒருபோதும் பாதியை அல்ல. `~/.local/bin/ccwork` போன்ற குறுக்குவழிக் கட்டளைகள் ccam binary-க்கான symlinks, அதனால் எந்த மாற்றமும் இல்லாமல் அவை புதிய பதிப்பை இயக்கும்.

Windows-இல் இயங்கும் நிரலின் கோப்பை மாற்றவோ நீக்கவோ முடியாது, பெயர் மாற்ற மட்டுமே முடியும். `Updater.Install` புதிய binary-ஐ `ccam.exe.new`-இல் எழுதி, இயங்கும் `ccam.exe`-ஐ `ccam.exe.old` எனப் பெயர் மாற்றி, புதிய கோப்பை `ccam.exe` எனப் பெயர் மாற்றுகிறது. கடைசிப் படி தோல்வியடைந்தால், பழைய binary மீண்டும் பழைய பெயருக்கு மாற்றப்படும். அடுத்த புதுப்பிப்பு `ccam.exe.old`-ஐ நீக்குகிறது. symbolic links-க்குப் பதிலாக hard links ஆக உள்ள குறுக்குவழிகள் ([குறுக்குவழிக் கட்டளைகள்](/features/shortcuts)) உடனே புதுப்பிக்கப்படுகின்றன, அதனால் அவையும் புதிய பதிப்பை இயக்கும்.

### தொகுப்பு மேலாளர் மூலம் நிறுவியவை

`.deb`-இலிருந்து நிறுவிய ccam apt-க்குச் சொந்தமானது, apt-க்குத் தெரியாமல் அதை மாற்றினால் தொகுப்புத் தரவுத்தளம் தவறாகிவிடும். பதிவிறக்கும் முன், binary யாருக்குச் சொந்தமானது என்று ccam சரிபார்க்கிறது:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` அல்லது `/sbin`-இல் உள்ள binary தொகுப்பு மேலாளரால் நிர்வகிக்கப்படுவதாகக் கருதப்படும்.
2. வேறிடத்தில், `dpkg` நிறுவப்பட்டிருந்து `dpkg -S <path>` வெற்றிபெற்றால், அதுவும் தொகுப்பு மேலாளரால் நிர்வகிக்கப்படுவதாகக் கருதப்படும்.
3. இல்லையென்றால், அங்கே எழுத முடியும் என்பதை உறுதிசெய்ய, ccam binary-க்கு அருகில் ஒரு தற்காலிக `.ccam-update-*` கோப்பை உருவாக்கி நீக்குகிறது.

தொகுப்பு மேலாளரால் நிர்வகிக்கப்படும் binary-க்கு, ccam நின்று apt படிகளை அச்சிடும்:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux அல்லாத கணினிகளில், install script வரி மட்டுமே அச்சிடப்படும்: Windows-இல் `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`. கோப்புறையில் எழுத முடியாதபோது, install script-ஐ மீண்டும் இயக்க ccam பரிந்துரைக்கும்; அது ccam-ஐ `~/.local/bin`-இல் நிறுவும்.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- development build தன் பதிப்பை `dev` எனத் தெரிவிக்கும், ஒப்பிட அதற்கு வெளியீடு இல்லை, அதனால் `ccam update`, `ccam update --check` இரண்டும் அதை மறுக்கும். `ccam update --version vX.Y.Z` அதன் மேல் ஒரு வெளியீட்டை நிறுவும்.
- தோல்வியடைந்த பதிவிறக்கம், உங்கள் தளத்துக்கான asset இல்லாமை, அல்லது checksum பொருந்தாமை ஆகியவை நிறுவப்பட்ட binary-ஐ அப்படியே விட்டுவிடும்.
- இல்லாத வெளியீடு `there is no release <tag>` என்று தோல்வியடையும். GitHub HTTP 403 அல்லது 429 மூலம் மறுத்தால், அது பெரும்பாலும் rate limit என்று ccam சொல்லி, பின்னர் மீண்டும் முயலச் சொல்லும்.
- `--check` படிக்க மட்டுமே செய்யும். binary தொகுப்பு மேலாளரால் நிர்வகிக்கப்படுகிறதா, எழுதக்கூடியதா என்று அது சரிபார்ப்பதில்லை.
- பழைய tag-உடன் `--version` பயன்படுத்துவதே முந்தைய பதிப்புக்குத் திரும்பும் வழி.

## தொடர்புடையவை

- [நிறுவுதல்](/docs/install)
- [கட்டளைகள்](/docs/commands)
- [அடிக்கடி கேட்கப்படும் கேள்விகள்](/docs/faq)
