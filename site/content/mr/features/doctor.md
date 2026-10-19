---
title: तपासणी आणि दुरुस्ती
description: ccam doctor ने शेअर लिंक्स, फाइल परवानग्या, लॉगिन आणि Claude Code ची आवृत्ती तपासा, आणि ccam doctor --fix ने सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` तुमचे ccam इन्स्टॉलेशन तपासते: Claude Code binary, प्रत्येक खात्याचा फोल्डर, लॉगिन आणि शेअर लिंक्स, शेअर केलेले MCP सर्व्हर, usage endpoint आणि तुमच्या शॉर्टकट कमांड्स. `ccam doctor --fix` डेटा न गमावता करता येईल ते दुरुस्त करते, आणि ज्या फाइलच्या जागी काही ठेवते ती आधी बॅकअप फोल्डरमध्ये हलवते. `ccam <name>` शेअर लिंकच्या समस्येमुळे थांबल्यावर, फाइल्स हाताने परत आणल्यानंतर किंवा Claude Code अपडेट झाल्यानंतर ते चालवा.

## वापर

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

प्रत्येक निष्कर्ष एका ओळीत चिन्हासह असतो: `✓` ठीक, `·` माहिती, `!` इशारा, `✗` त्रुटी. खात्याचे निष्कर्ष खात्याच्या नावाने सुरू होतात.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

`--fix` सह, दुरुस्त झालेला प्रत्येक निष्कर्ष `→ fixed` ने संपतो:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--fix` | सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा (बदललेल्या फाइल्सचा बॅकअप घेतला जातो) | बंद |
| `--dry-run` (global) | प्रत्येक तपासणी चालवा पण काहीही दुरुस्त करू नका, `--fix` असला तरी | बंद |

`ccam doctor` कोणतेही आर्ग्युमेंट घेत नाही. ते `--dry-run` पाळते: `ccam --dry-run doctor --fix` `ccam doctor` सारखेच निष्कर्ष छापते आणि काहीही बदलत नाही.

किमान एक `✗` निष्कर्ष दुरुस्त न होता उरला तर exit code `1` असतो, नाहीतर `0`. इशारे आणि माहितीमुळे run कधीच अयशस्वी होत नाही. `--fix` शिवाय run अयशस्वी झाल्यास ccam एक रिकामी ओळ आणि `` `ccam doctor --fix` repairs what it can. `` ने संपते. लोड न होऊ शकणारा config देखील `1` आणि `ccam:` त्रुटीने संपतो.

## हे कसे काम करते

`internal/cli/doctor.go` `~/.config/ccam/config.toml` लोड करते आणि `internal/doctor/doctor.go` मधील `doctor.Run` बोलावते, जे तपासण्या एका निश्चित क्रमाने चालवते आणि निष्कर्ष परत देते. `--fix` `Options.Fix` म्हणून पुढे दिला जातो, जो dry run मध्ये false असतो.

<!-- motion:doctor-checks -->

मधला भाग प्रत्येक खात्यासाठी एकदा चालतो. तपासण्या, आणि प्रत्येकीसाठी `--fix` काय करते:

1. **Claude Code आवृत्ती** (`checkClaude`). launcher प्रमाणेच `claude` binary शोधते आणि `claude --version` चालवते. आवृत्ती `2.1.302` ने सुरू होत नसेल, म्हणजे ज्या आवृत्तीवर ccam चा फोल्डर layout पडताळला गेला ती (`doctor.LayoutVersion`), तर ते `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` जोडते. binary नसणे किंवा `--version` अयशस्वी होणे ही त्रुटी आहे. दुरुस्ती नाही.
2. **environment मधील credentials** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` किंवा `ANTHROPIC_AUTH_TOKEN` सेट असेल तर इशारा देते: त्या shell मधून सुरू होणाऱ्या प्रत्येक खात्याला ते वारसाने मिळते, आणि Claude Code खात्याच्या स्वतःच्या लॉगिनपेक्षा त्याला प्राधान्य देते. दुरुस्ती नाही; ते स्वतः unset करा.
3. **फोल्डर** (`checkProfile`). खात्याचा फोल्डर नसणे ही त्रुटी आहे आणि त्यानंतर त्या खात्याच्या तपासण्या थांबतात. group किंवा other bits सेट असलेल्या परवानग्या म्हणजे इशारा; `--fix` `chmod 0700` चालवते.
4. **लॉगिन फाइल** (`checkCredentials`). `.credentials.json` `lstat` ने तपासली जाते. symlink ही त्रुटी आहे, कारण प्रत्येक खात्याचे स्वतःचे token असणे आवश्यक आहे, आणि ती दुरुस्त होत नाही. `0600` शिवाय इतर mode म्हणजे इशारा; `--fix` `chmod 0600` चालवते. मग लॉगिन वाचले जाते (macOS वर आधी Keychain मधून). लॉगिन नसणे हा इशारा आहे, जो `ccam login <name>` कडे निर्देश करतो. कालबाह्य refresh token ही `login expired` त्रुटी आहे; कालबाह्य access token ही फक्त माहिती आहे, कारण Claude Code पुढच्या वेळी सुरू होताना ते नूतनीकरण करते. यांपैकी काहीही दुरुस्त होत नाही.
5. **शेअर लिंक्स** (`share.Check`). `share.Layout` मधील प्रत्येक पाथची खात्याच्या शेअर सेटिंग्जशी तुलना करते, फक्त `lstat` आणि `readlink` वापरून. ते चार प्रकार नोंदवते: नसलेली लिंक, दुसरीकडे निर्देश करणारी लिंक, लिंकच्या जागी असलेली खरी फाइल, आणि आता local असलेल्या आयटमसाठी उरलेली लिंक. नसलेली लिंक म्हणजे इशारा, बाकीच्या त्रुटी.
6. **अज्ञात नोंदी** (`checkUnknownEntries`). खात्याच्या फोल्डरमधील ccam ला माहीत नसलेली (`share.Known`) top-level नावे दाखवते, Claude Code चा layout बदलला असल्याची सूचना म्हणून. फक्त माहिती.
7. **Plugin पाथ** (`checkPlugins`). Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) absolute पाथ साठवतात. प्रत्येक गटाच्या `plugins/` मध्ये, एखाद्या खात्याच्या फोल्डरमधून जाणारे पाथ ते खाते काढल्यावर किंवा त्याचे नाव बदलल्यावर तुटतात. `--fix` ते `profile.RewritePluginPaths` ने गट फोल्डरकडे पुन्हा लिहिते, registry फाइलचे atomic write करून.
8. **MCP सर्व्हर** (`checkMCP`). प्रत्येक गटाची `mcp.json` एकदा लोड करते. वाचता न येणारी फाइल ही त्रुटी आहे; `mcp.Validate` मध्ये अयशस्वी होणारा सर्व्हर हा इशारा आहे. दुरुस्ती नाही.
9. **Usage response format** (`checkUsageFormat`). cache शिवाय आणि नेहमीच्या 3 सेकंदांच्या मर्यादेसह, एकामागून एक खात्यासाठी usage endpoint ला विचारते, जोपर्यंत एखादे उत्तर देत नाही. response आता parse होत नसेल, तर ते इशारा देते की `ccam usage` स्टेटस लाइन cache वर अवलंबून राहील. नेटवर्क वापरणारी ही एकमेव तपासणी आहे.
10. **शॉर्टकट** (`checkShims`). फक्त एखाद्या खात्याला शॉर्टकट असेल तेव्हा. `~/.local/bin` मध्ये नसलेल्या लिंक्स म्हणजे इशारा; `--fix` `shim.Install` बोलावते, जे फक्त नसलेल्या लिंक्स तयार करते आणि इतर फाइल्सना कधीच हात लावत नाही. `PATH` वर नसलेला bin फोल्डर हा इशारा आहे आणि दुरुस्त होत नाही.

### शेअर लिंक्स दुरुस्त करणे

सुरू करण्याचा मार्ग नसलेली लिंक तिथल्या तिथे तयार करतो, पण डेटा गमावू शकेल अशा कोणत्याही गोष्टीसाठी थांबतो आणि इथे निर्देश करतो:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

कोणत्याही शेअर लिंक समस्येसाठी `--fix` खात्याच्या शेअर सेटिंग्ज सध्याची आणि अपेक्षित दोन्ही स्थिती म्हणून देऊन `share.Build` बोलावते, मग `Plan.Apply`. `Build` प्रत्येक टप्पा डिस्कवर जे आहे त्यावरून पुन्हा ठरवते, `ccam share` वापरते तीच योजना, त्यामुळे दुरुस्ती शेअरिंग चालू करण्याचेच नियम पाळते.

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
- `doctor --fix` खाते चालू आहे का ते तपासत नाही. खात्याच्या शेअर लिंक्स दुरुस्त करण्यापूर्वी त्या खात्याचे Claude Code बंद करा, कारण session त्या फाइल्समध्ये लिहीत राहते.
- शेअर दुरुस्ती अयशस्वी झाली, तर ccam त्या खात्यासाठी `could not repair share links: <error>` नोंदवते आणि त्याची अज्ञात नोंदींची तपासणी वगळते. इतर खाती तरीही तपासली जातात.
- `ccam share` अर्ध्यावर अयशस्वी झाले, तर ते आधी अपेक्षित सेटिंग्ज सेव्ह करते आणि `ccam doctor --fix` ने पूर्ण करायला सांगते.
- symlink केलेली `.credentials.json`, कालबाह्य लॉगिन आणि environment मधील credentials नोंदवले जातात पण कधीच बदलले जात नाहीत.
- Plugin registry चे पुनर्लेखन atomic पद्धतीने लिहिले जाते, पण `.ccam-backup/` मध्ये कॉपी होत नाही.

## संबंधित

- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [आधीचा सेटअप कॉपी करणे](/features/migrate)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [कमांड्स](/docs/commands)
