---
title: शॉर्टकट कमांड्स
description: ccam shortcut ने खात्याला ccwork सारखी स्वतःची कमांड द्या, आणि ccam shims install ने ~/.local/bin मधील लिंक्स पुन्हा तयार करा.
updated: 2026-10-26
order: 3
translation: machine
sourceHash: "066e7e622517"
---

शॉर्टकट ही तुम्ही स्वतः नाव दिलेली कमांड आहे, जी Claude Code एका खात्याने सुरू करते: `ccwork --resume` तेच करते जे `ccam work --resume` करते. ती `~/.local/bin` मधील ccam ची लिंक आहे, त्यामुळे ती bash, zsh, scripts आणि `PATH` वरून प्रोग्राम चालवणाऱ्या इतर सर्व गोष्टींमध्ये चालते. शॉर्टकट ऐच्छिक आहेत; तुम्ही वारंवार सुरू करता त्या खात्यांसाठी एक तयार करा.

## वापर

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

ccam प्रत्येक बदलाची खात्री देते:

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

एका खात्याला जास्तीत जास्त एक शॉर्टकट असतो, जो `~/.config/ccam/config.toml` मध्ये `shortcut` म्हणून साठवला जातो. `ccam add` विझार्ड `Create a shortcut command for <name>? [y/N]` आणि मग `Shortcut name [cc<name>]:` विचारतो. विझार्डबाहेर, `ccam add` फक्त `--shortcut` सह शॉर्टकट तयार करते. `ccam list` प्रत्येक खात्याचा शॉर्टकट दाखवते, किंवा `-`.

शॉर्टकट प्रत्येक आर्ग्युमेंट Claude Code ला देतो, जसे `ccam <name>` नावानंतर करते. यात `--dry-run` आणि `--help` सुद्धा येतात: सुरू करण्याचे पूर्वावलोकन करण्यासाठी `ccam --dry-run work` वापरा.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--remove` | खात्याचा शॉर्टकट आणि त्याची लिंक हटवा | बंद |
| `-h`, `--help` | कमांडची मदत | |
| `--dry-run` | `ccam shortcut` त्याला समर्थन देते: ते नाव तपासते आणि सेव्ह किंवा लिंक न करता `Would set the shortcut of work to ccwork.` किंवा `Would remove ccwork, the shortcut of work.` छापते | बंद |

`ccam shims install` कोणतेही आर्ग्युमेंट घेत नाही आणि `--dry-run` `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` देऊन नाकारते.

नावांचे नियम, या क्रमाने तपासले जातात:

- `a-z`, `0-9`, `.`, `-` आणि `_` यांपासून बनलेली 1 ते 64 अक्षरे, अक्षराने किंवा अंकाने सुरू होणारी;
- `ccam`, `claude` किंवा `run` किंवा `usage` सारखी ccam कमांड नाही;
- दुसऱ्या खात्याचा शॉर्टकट नाही;
- `~/.local/bin` मधील आधीपासून असलेली फाइल नाही, ती ccam ने तयार केलेली लिंक (जसे या खात्याचा स्वतःचा शॉर्टकट) असल्याशिवाय;
- तुमच्या `PATH` वरील दुसऱ्या executable चे नाव नाही.

`config.toml` लोड होताना त्याच नियमांनुसार तपासले जाते. तुम्ही शॉर्टकट हाताने असा संपादित केला की तो pattern मोडतो, `claude` किंवा `ccam` आहे, किंवा दोन खाती तो वापरतात, तर तुम्ही तो दुरुस्त करेपर्यंत प्रत्येक कमांड त्या खात्याचे नाव सांगणाऱ्या संदेशासह अयशस्वी होते.

कमांड आणि `--remove` एकत्र दिल्यास `give a command or --remove, not both` ने थांबते. खात्याकडे आधीच असलेले नाव सेट केल्यास `work: ccwork (unchanged)` छापले जाते आणि लिंक नसल्यास ती पुन्हा तयार होते.

## हे कसे काम करते

नाव सेट केल्यावर `checkShortcut` चालते, जे bin फोल्डर आणि ccam binary चा खरा पाथ देऊन `shim.Check` बोलावते:

<!-- motion:shortcuts-check -->

नाव तपासण्यांमधून क्रमाने खाली जाते. प्रत्येक तपासणीशेजारी ती नाकारत असलेले एक नाव आहे: `Work` pattern मध्ये बसत नाही, `run` ही ccam कमांड आहे, `ccside` दुसऱ्या खात्याचे आहे, `ccfile` ही `~/.local/bin` मधील आधीपासून असलेली फाइल आहे, आणि `git` `/usr/bin/git` लपवेल. प्रत्येक तपासणी पार करणारे नावच सेव्ह आणि लिंक होते.

1. ccam नाव `^[a-z0-9][a-z0-9._-]{0,63}$` शी, मग `config.Reserved` शी, मग `cfg.FindShortcut` शी जुळवते (`shim.Check`). ते `~/.local/bin/<name>` वर `os.Lstat` बोलावते, आणि `PATH` वरील इतर प्रत्येक फोल्डरसाठी `os.Stat` बोलावून ccam ची लिंक नसलेली executable फाइल नाकारते.
2. ccam `config.toml` `fsutil.AtomicWrite` ने सेव्ह करते (तात्पुरती फाइल, `fsync`, rename).
3. ccam जुन्या शॉर्टकटची लिंक हटवते, पण ती ccam ची लिंक असेल तरच (`shim.Remove`).
4. ccam config मधील प्रत्येक शॉर्टकट नावाच्या क्रमाने पाहते आणि नसलेल्या प्रत्येकाची लिंक तयार करते (`shim.Install`), `os.Symlink` ने ccam binary कडे, जी `os.Executable` आणि `filepath.EvalSymlinks` ने शोधली जाते. एखाद्या शॉर्टकटच्या जागी असलेली फाइल इतर शॉर्टकट थांबवत नाही; त्रुटी प्रत्येकाचे नाव सांगते.

ccam तयार करत असलेली प्रत्येक लिंक, आणि प्रत्येक लिंक कोणती binary चालवते, `~/.local/state/ccam/shortcuts.json` मध्ये नोंदवते. `os.Readlink` ccam चा पाथ परत देते तेव्हा, किंवा लिंक आणि ccam चा पाथ एकाच फाइलकडे जातात तेव्हा, ती लिंक ccam ची स्वतःची मानली जाते. नोंदवलेल्या लिंकची binary नंतर हलवली किंवा बदलली गेली असेल (दुसऱ्या फोल्डरमध्ये हलवलेली ccam binary, पुन्हा केलेले इन्स्टॉलेशन), तर ती लिंक जुनी ठरते, आणि ccam तिला या ccam binary कडे वळवते (`shim.Install`). ccam च्या दुसऱ्या इन्स्टॉल असलेल्या प्रतीकडे जाणारी नोंदवलेली लिंक आहे तशीच ठेवली जाते. प्रत्येक remove आणि प्रत्येक overwrite तपासणी या चाचण्यांमधून जाते, त्यामुळे ccam ने तयार न केलेल्या फाइल्सना, जसे ccam ची तुमची स्वतःची प्रत किंवा तुम्ही स्वतः केलेली लिंक, ते हात लावत नाही.

तुम्ही शॉर्टकट चालवता तेव्हा, तुम्ही टाइप केलेले नाव हीच एक गोष्ट ccam ला लागते:

<!-- motion:shortcuts-link -->

shell ला `PATH` वर `ccwork` सापडतो, kernel लिंक पाळून ccam binary सुरू करतो, आणि `argv[0]` तरीही `ccwork` असतो. `app.Main` मध्ये, `ccam` शिवाय इतर कोणतेही नाव `config.toml` एकदा लोड करते आणि `shim.Lookup` बोलावते, जे ज्या खात्याचा `shortcut` `filepath.Base(argv[0])` शी जुळतो ते शोधते. मग ccam आधीच लोड झालेल्या config सह [`ccam <name>`](/features/launch) सारखेच टप्पे घेते, आणि `argv[1:]` न बदलता `claude` ला देते (`launcher.LaunchWith`).

shell alias किंवा wrapper script नाही तर symlink का: symlink प्रत्येक shell मध्ये आणि scripts मध्ये चालते, ती अतिरिक्त shell सुरू करणे टाळते, ज्याला सुमारे 1 ते 2 ms लागले असते, आणि सुरू करताना तरीही एकच फाइल वाचली जाते. `scripts/perf-check.sh` `hyperfine` ने शॉर्टकट मोजते, `ccam <name>` च्याच 10 ms (p95) मर्यादेवर.

### Windows वर

शॉर्टकट म्हणजे `%USERPROFILE%\.local\bin` मधील `ccwork.exe`, आणि ccam शोधण्यापूर्वी `.exe` (case कोणताही असो) काढून टाकते (`shim.Name`), त्यामुळे `ccwork`, `ccwork.exe` आणि `CCWORK.EXE` तिन्ही `work` सुरू करतात. ccam आधी symbolic link तयार करते, जी Windows Developer Mode मध्ये किंवा administrator साठी परवानगी देते. तो अधिकार नसल्यास ते hard link तयार करते: त्याच `ccam.exe` फाइलचे दुसरे नाव, ज्याला कोणताही विशेषाधिकार लागत नाही. hard link `ccam.exe` चीच फाइल असेल तर ती ccam ची स्वतःची मानली जाते.

hard link नेहमी ज्या फाइलपासून तयार झाली तीच फाइल चालवते, म्हणून `ccam update` `ccam.exe` च्या जागी नवीन फाइल ठेवते तेव्हा, ते `shim.Relink` ने hard link असलेला प्रत्येक शॉर्टकट नवा करते आणि `✓ shortcut ccwork renewed` छापते. `install.ps1` पुन्हा चालवल्याने `ccam shims install` मार्फत हेच होते (`OK 2 shortcuts now run the new ccam`), आणि ccam च्या जुन्या आवृत्तीची प्रत असलेल्या कोणत्याही शॉर्टकटसाठी `ccam shims install` आणि `ccam doctor --fix` सुद्धा हेच करतात. नावाची तपासणी इतर `PATH` फोल्डर्समध्ये `PATHEXT` मधील प्रत्येक extension (`.exe`, `.cmd` वगैरे) सह `<name>` शोधते.

शॉर्टकट shell मध्ये टाइप करा (PowerShell, cmd.exe, Git Bash). Explorer, Run dialog किंवा `Start-Process` मधून उघडलेली symbolic link साधे `ccam` म्हणून सुरू होते, जे खात्याऐवजी डॅशबोर्ड उघडते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- **कोणत्याही खात्याची नसलेली लिंक.** `~/.local/bin` मधील ccam च्या लिंकचे नाव कोणतेही खाते वापरत नसेल, उदाहरणार्थ तुम्ही `config.toml` मधील `shortcut` हाताने बदलल्यानंतर, तर ccam डीफॉल्ट खात्यावर जात नाही. ते `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` ने थांबते. `ccam shortcut` ने नाव पुन्हा सेट करा, किंवा लिंक हटवा.
- **आधीपासून असलेल्या फाइल्स.** ccam ने तयार न केलेली फाइल ते कधीच बदलत नाही. `ccam shims install` `<path> already exists and is not a ccam link, so it was left alone` कळवते आणि तरीही इतर लिंक्स तयार करते. तेव्हा `ccam add` आणि `ccam import` `shortcut: ccwork (not created; see above)` छापतात.
- **हलवलेली ccam binary.** तुम्ही ccam binary हलवल्यास, ccam ने तयार केलेल्या लिंक्स आता अस्तित्वात नसलेल्या पाथकडे निर्देश करतात. `ccam shims install` (किंवा `ccam doctor --fix`) त्यांना तुम्ही ज्या ccam binary ने ते चालवता तिच्याकडे वळवते.
- **नाव बदलणे आणि काढणे.** `ccam rename` शॉर्टकटचे नाव तेच ठेवते, आणि लिंक चालत राहते, कारण ती ccam कडे निर्देश करते आणि ccam `config.toml` मध्ये नाव शोधते. `ccam remove` खात्याची लिंक हटवते. `ccam shortcut <account> <new>` नवीन लिंक तयार करण्यापूर्वी जुनी हटवते.
- **PATH.** ccam लिंक तयार करते आणि `~/.local/bin` तुमच्या `PATH` वर नसतो, तेव्हा ते `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` आणि त्यानंतर जोडायची `export PATH=...` ओळ छापते. Windows वर त्याऐवजी ते फोल्डर तुमच्या user `PATH` मध्ये जोडणारी PowerShell ओळ छापते.
- **नंतर तपासणे.** `ccam doctor` `missing shortcuts in /home/me/.local/bin: ccwork` आणि ccam ची जुनी आवृत्ती चालवणारे शॉर्टकट नोंदवते, आणि फोल्डर `PATH` वर नसल्यास इशारा देते. `ccam doctor --fix` नसलेल्या लिंक्स तयार करते आणि जुन्या लिंक्स नव्या करते, आणि फक्त स्वतः तयार केलेल्यांसाठीच `→ fixed` सांगते. अजूनही इन्स्टॉल असलेले दुसरे ccam चालवणारा शॉर्टकट चालतो, म्हणून `ccam doctor` तो फक्त माहिती म्हणून दाखवते; `ccam doctor --fix` तो तुम्ही ज्या ccam binary ने ते चालवता तिच्याकडे वळवते, आणि `ccam add` व `ccam import` तो ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)`` असा दाखवतात.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [खाती व्यवस्थापित करणे](/features/accounts)
- [कमांड संदर्भ](/docs/commands#shortcuts)
- [ccam इन्स्टॉल करा](/docs/install)
