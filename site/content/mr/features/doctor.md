---
title: तपासणी आणि दुरुस्ती
description: ccam doctor ने शेअर लिंक्स, फाइल परवानग्या, लॉगिन आणि Claude Code ची आवृत्ती तपासा, आणि ccam doctor --fix ने सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा.
updated: 2026-10-21
order: 9
translation: machine
sourceHash: "9be63350cb69"
---

`ccam doctor` तुमचे ccam इन्स्टॉलेशन तपासते: Claude Code binary, प्रत्येक खात्याचा फोल्डर, लॉगिन आणि शेअर लिंक्स, शेअर केलेले MCP सर्व्हर, usage endpoint आणि तुमच्या शॉर्टकट कमांड्स. `ccam doctor --fix` डेटा न गमावता करता येईल ते दुरुस्त करते, आणि ज्या फाइलच्या जागी काही ठेवते ती आधी बॅकअप फोल्डरमध्ये हलवते. `ccam <name>` शेअर लिंकच्या समस्येमुळे थांबल्यावर, फाइल्स हाताने परत आणल्यानंतर किंवा Claude Code अपडेट झाल्यानंतर ते चालवा.

## वापर

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

पहिली ओळ सर्व खात्यांना एकाच वेळी चिन्ह देते, त्यामुळे बाकी काही वाचण्यापूर्वीच समस्या कुठे आहेत ते तुम्हाला दिसते. त्याखाली प्रत्येक खात्याला एक box मिळतो, आणि कोणत्याही खात्याशी संबंधित नसलेले निष्कर्ष (Claude Code, environment variables, plugins, MCP सर्व्हर, शॉर्टकट्स, फोल्डर नियम, गट आणि सेटिंग्ज) `General` box मध्ये जातात. box च्या शीर्षकात खात्याचे सर्वात वाईट चिन्ह आणि त्यात किती समस्या आहेत ते दिसते. आतली प्रत्येक ओळ तपासलेले एक क्षेत्र आहे, नेहमी `Folder`, `Login` आणि `Share links` याच क्रमाने, मग इतर कोणतीही, चिन्हासह: `✔` ठीक, `ℹ` माहिती, `⚠` इशारा, `✘` त्रुटी. काहीही चूक नसलेला फोल्डर `✔ Folder  OK` दाखवतो.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix` सह, दुरुस्त झालेला प्रत्येक निष्कर्ष `→ fixed` ने संपतो, आणि तो आता समस्या म्हणून मोजला जात नाही:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

टर्मिनलवर boxes त्याच्या रुंदीत बसतात आणि लांब संदेश पुढच्या ओळीत जातात, टर्मिनलच्या स्वतःच्या रंगांत (`NO_COLOR` ते बंद करते). pipe किंवा फाइलकडे पाठवल्यास आउटपुटमध्ये रंग नसतात आणि कोणताही संदेश तोडला किंवा कापला जात नाही, त्यामुळे `ccam doctor | grep settings.json` संपूर्ण ओळ शोधते; तेव्हा प्रत्येक box सर्वात लांब ओळीइतका रुंद असतो. डॅशबोर्डचा Health टॅब हेच boxes दाखवतो.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--fix` | सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा (बदललेल्या फाइल्सचा बॅकअप घेतला जातो) | बंद |
| `--dry-run` (global) | प्रत्येक तपासणी चालवा पण काहीही दुरुस्त करू नका, `--fix` असला तरी | बंद |

`ccam doctor` कोणतेही आर्ग्युमेंट घेत नाही. ते `--dry-run` पाळते: `ccam --dry-run doctor --fix` `ccam doctor` सारखेच निष्कर्ष छापते आणि काहीही बदलत नाही.

किमान एक `✘` निष्कर्ष दुरुस्त न होता उरला तर exit code `1` असतो, नाहीतर `0`. इशारे आणि माहितीमुळे run कधीच अयशस्वी होत नाही. `--fix` शिवाय run अयशस्वी झाल्यास ccam एक रिकामी ओळ आणि `` `ccam doctor --fix` repairs what it can. `` ने संपते. लोड न होऊ शकणारा config देखील `1` आणि `ccam:` त्रुटीने संपतो.

## हे कसे काम करते

`internal/cli/doctor.go` `~/.config/ccam/config.toml` लोड करते आणि `internal/doctor/doctor.go` मधील `doctor.Run` बोलावते, जे तपासण्या एका निश्चित क्रमाने चालवते आणि निष्कर्ष परत देते, प्रत्येकासोबत तो ज्या खात्याचा आहे ते खाते आणि त्याने तपासलेले क्षेत्र. `tui.HealthReport` त्यांना boxes मध्ये गटबद्ध करते. `--fix` `Options.Fix` म्हणून पुढे दिला जातो, जो dry run मध्ये false असतो.

<!-- motion:doctor-checks -->

मधला भाग प्रत्येक खात्यासाठी एकदा चालतो. तपासण्या, आणि प्रत्येकीसाठी `--fix` काय करते:

1. **Claude Code आवृत्ती** (`checkClaude`). launcher प्रमाणेच `claude` binary शोधते आणि `claude --version` चालवते. आवृत्ती `2.1.302` ने सुरू होत नसेल, म्हणजे ज्या आवृत्तीवर ccam चा फोल्डर layout पडताळला गेला ती (`doctor.LayoutVersion`), तर ते माहिती `ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` जोडते. binary नसणे किंवा `--version` अयशस्वी होणे ही त्रुटी आहे. दुरुस्ती नाही.
2. **environment मधील credentials** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` किंवा `ANTHROPIC_AUTH_TOKEN` सेट असेल तर इशारा देते: त्या shell मधून सुरू होणाऱ्या प्रत्येक खात्याला ते वारसाने मिळते, आणि Claude Code खात्याच्या स्वतःच्या लॉगिनपेक्षा त्याला प्राधान्य देते. दुरुस्ती नाही; ते स्वतः unset करा.
3. **फोल्डर** (`checkProfile`). खात्याचा फोल्डर नसणे ही त्रुटी आहे आणि त्यानंतर त्या खात्याच्या तपासण्या थांबतात. group किंवा other bits सेट असलेल्या परवानग्या म्हणजे इशारा; `--fix` `chmod 0700` चालवते. Windows वर परवानगी bits नसतात, त्यामुळे तिथे ही तपासणी आणि खालची `0600` तपासणी वगळल्या जातात; तुमच्या user profile मधील फोल्डर त्याच्या access list ने सुरक्षित असतो.
4. **लॉगिन फाइल** (`checkCredentials`). `.credentials.json` `lstat` ने तपासली जाते. symlink ही त्रुटी आहे, कारण प्रत्येक खात्याचे स्वतःचे token असणे आवश्यक आहे, आणि ती दुरुस्त होत नाही. `0600` शिवाय इतर mode म्हणजे इशारा; `--fix` `chmod 0600` चालवते. मग लॉगिन वाचले जाते (macOS वर आधी Keychain मधून). लॉगिन नसणे हा इशारा आहे, जो `ccam login <name>` कडे निर्देश करतो. कालबाह्य refresh token ही `login expired` त्रुटी आहे; कालबाह्य access token ही फक्त माहिती आहे, कारण Claude Code पुढच्या वेळी सुरू होताना ते नूतनीकरण करते. वैध लॉगिन email आणि plan दाखवते (`login valid (dev@company.com, Claude Team · Acme)`). यांपैकी काहीही दुरुस्त होत नाही; `ccam login --expired` गरज असलेल्या प्रत्येक खात्यात लॉगिन करते.
5. **शेअर लिंक्स** (`share.Check`). `share.Layout` मधील प्रत्येक पाथची खात्याच्या गटांशी तुलना करते, फक्त `lstat` आणि `readlink` वापरून. ते चार प्रकार नोंदवते: नसलेली लिंक, दुसरीकडे निर्देश करणारी लिंक, लिंकच्या जागी असलेली खरी फाइल, आणि आता local असलेल्या आयटमसाठी उरलेली लिंक. नसलेली लिंक म्हणजे इशारा, बाकीच्या त्रुटी. ccam च्या नंतरच्या आवृत्तीत एखाद्या आयटममध्ये जोडलेले पाथ (history साठी `file-history/`, `tasks/` आणि `paste-cache/`, settings साठी `keybindings.json`) खात्याकडे अजून स्वतःची प्रत असेपर्यंत फक्त इशारा असतात: `is not shared with the group yet (ccam doctor --fix moves it in)`. त्यांच्यामुळे launch कधीच थांबत नाही.
6. **अज्ञात नोंदी** (`checkUnknownEntries`). खात्याच्या फोल्डरमधील ccam ला माहीत नसलेली (`share.Known`) top-level नावे दाखवते, Claude Code चा layout बदलला असल्याची सूचना म्हणून. फक्त माहिती.
7. **Plugin पाथ** (`checkPlugins`). Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) absolute पाथ साठवतात. प्रत्येक गटाच्या `plugins/` मध्ये, एखाद्या खात्याच्या फोल्डरमधून जाणारे पाथ ते खाते काढल्यावर किंवा त्याचे नाव बदलल्यावर तुटतात. `--fix` ते `profile.RewritePluginPaths` ने गट फोल्डरकडे पुन्हा लिहिते, registry फाइलचे atomic write करून.
8. **MCP सर्व्हर** (`checkMCP`). प्रत्येक गटाची `mcp.json` एकदा लोड करते. वाचता न येणारी फाइल ही त्रुटी आहे; `mcp.Validate` मध्ये अयशस्वी होणारा सर्व्हर हा इशारा आहे. दुरुस्ती नाही.
9. **Usage response format** (`checkUsageFormat`). cache शिवाय आणि नेहमीच्या 3 सेकंदांच्या मर्यादेसह, एकामागून एक खात्यासाठी usage endpoint ला विचारते, जोपर्यंत एखादे उत्तर देत नाही. response आता parse होत नसेल, तर ते इशारा देते की `ccam usage` स्टेटस लाइन cache वर अवलंबून राहील. नेटवर्क वापरणारी ही एकमेव तपासणी आहे.
10. **शॉर्टकट** (`checkShims`). फक्त एखाद्या खात्याला शॉर्टकट असेल तेव्हा. `~/.local/bin` मध्ये नसलेल्या लिंक्स म्हणजे इशारा; `--fix` `shim.Install` बोलावते, जे फक्त नसलेल्या लिंक्स तयार करते आणि इतर फाइल्सना कधीच हात लावत नाही. `PATH` वर नसलेला bin फोल्डर हा इशारा आहे आणि दुरुस्त होत नाही.
11. **फोल्डर नोंदी** (`checkFolders`). फक्त [`ccam folder`](/features/folder) च्या नोंदी असतील तेव्हा. ज्या नोंदीचे खाते अस्तित्वात नाही ती त्रुटी आहे, आणि ज्या नोंदीचा फोल्डर नाहीसा झाला आहे तो इशारा आहे. दोन्हीपैकी काहीही दुरुस्त होत नाही, कारण नोंद काढल्याने तिथे कोणते खाते सुरू होते ते बदलेल.
12. **गट आणि सेटिंग्ज** (`checkGroups`). खाते नसलेला गट सदस्य, ज्याचे नाव नंतर ccam कमांड बनले आहे असे खाते (ते `ccam run <name>` ने सुरू करा किंवा त्याचे नाव बदला), आणि `config.toml` मधील अज्ञात `theme` किंवा `usage_source` हे इशारे आहेत. दुरुस्ती नाही.
13. **लॉगिन** (`checkIdentities`). प्रत्येक फोल्डर कोणत्या खात्याने लॉगिन आहे ते त्याच्या `.claude.json` मधून वाचते. `ccam login` ने नोंदवलेल्या खात्याऐवजी आता दुसऱ्या खात्यात लॉगिन असलेला फोल्डर (उदाहरणार्थ session मध्ये `/login` केल्यानंतर) हा इशारा आहे, आणि दुसऱ्या फोल्डरच्याच Claude खात्यात लॉगिन असलेला फोल्डरही: ती दोन्ही आपल्या वापर मर्यादा शेअर करतात. दुरुस्ती नाही.

### शेअर लिंक्स दुरुस्त करणे

सुरू करण्याचा मार्ग नसलेली लिंक तिथल्या तिथे तयार करतो, पण डेटा गमावू शकेल अशा कोणत्याही गोष्टीसाठी थांबतो आणि इथे निर्देश करतो:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

कोणत्याही शेअर लिंक समस्येसाठी `--fix` खात्याच्या शेअर सेटिंग्ज सध्याची आणि अपेक्षित दोन्ही स्थिती म्हणून देऊन `share.Build` बोलावते, मग `Plan.Apply`. `Build` प्रत्येक टप्पा डिस्कवर जे आहे त्यावरून पुन्हा ठरवते, `ccam group` बदल वापरतो तीच योजना, त्यामुळे दुरुस्ती गटात सामील होण्याचेच नियम पाळते.

<!-- motion:doctor-repair -->

इथे `settings.json` ही खरी फाइल आहे. गटाकडे आधीच वेगळी `settings.json` असल्यामुळे conflict त्याचा डीफॉल्ट पर्याय `KeepGroup` घेतो, त्यामुळे गटाची प्रत राहते. खात्याची फाइल खात्याच्या फोल्डरमधील `.ccam-backup/<time>/settings.json` मध्ये हलते, आणि `os.Symlink` तो पाथ गटाशी लिंक करते. इतर स्थिती:

- गटाकडे अजून प्रत नसेल, तर खात्याची फाइल गटात हलते आणि त्याचा पहिला मजकूर बनते.
- फोल्डरसाठी, गटाकडे नसलेल्या फाइल्स गटात हलतात; बाकी बॅकअपमध्ये जाते.
- `history.jsonl` साठी डीफॉल्ट `merge by line` आहे: गटाच्या ओळी, मग तुमच्या नवीन ओळी. गटाची आधीची प्रत `.ccam-backup/<time>/group/` मध्ये जतन होते.
- दुसऱ्या गटाकडे जाणारी लिंक बदलली जाते. ccam बाहेरील फोल्डरकडे जाणारी लिंक बॅकअपमध्ये जाते, आणि गट रिकामा असेल तर तिचा मजकूर गटाचा पहिला मजकूर बनतो.
- आता local असलेल्या आयटमसाठी उरलेली लिंक काढली जाते, आणि खात्याला गटाच्या मजकुराची प्रत मिळते.

लिंक्स काढणे सोडल्यास, प्रत्येक टप्पा डेटा गटात किंवा बॅकअप फोल्डरमध्ये हलवतो. काहीही हटवले जात नाही.

## सुरक्षितता आणि अपवादात्मक स्थिती

- साधे `ccam doctor` कधीच काहीही बदलत नाही. scripts मध्ये खात्री करण्यासाठी `--fix` सोबत `--dry-run` वापरा.
- चालू असलेल्या खात्यासाठी `doctor --fix` फक्त नसलेल्या लिंक्स बनवते आणि `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` जोडते, कारण दुरुस्ती ज्या फाइल्स हलवेल त्यांत session लिहीत राहते.
- शेअर दुरुस्ती अयशस्वी झाली, तर ccam त्या खात्यासाठी `could not repair share links: <error>` नोंदवते आणि त्याची अज्ञात नोंदींची तपासणी वगळते. इतर खाती तरीही तपासली जातात.
- `ccam group` बदल अर्ध्यावर अयशस्वी झाला, तर नवीन गट आधीच सेव्ह झालेले असतात, आणि ccam `ccam doctor --fix` ने पूर्ण करायला सांगते.
- symlink केलेली `.credentials.json`, कालबाह्य लॉगिन आणि environment मधील credentials नोंदवले जातात पण कधीच बदलले जात नाहीत.
- Plugin registry चे पुनर्लेखन atomic पद्धतीने लिहिले जाते, पण `.ccam-backup/` मध्ये कॉपी होत नाही.

## संबंधित

- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [आधीचा सेटअप कॉपी करणे](/features/migrate)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [कमांड्स](/docs/commands)
