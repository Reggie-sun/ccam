---
title: तपासणी आणि दुरुस्ती
description: ccam doctor ने शेअर लिंक्स, फाइल परवानग्या, लॉगिन आणि Claude Code ची आवृत्ती तपासा, आणि ccam doctor --fix ने सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` तुमचे ccam इन्स्टॉलेशन तपासते: Claude Code binary, प्रत्येक खात्याचा फोल्डर, लॉगिन आणि शेअर लिंक्स, शेअर केलेले MCP सर्व्हर, usage endpoint आणि तुमच्या शॉर्टकट कमांड्स. `ccam doctor --fix` डेटा न गमावता करता येईल ते दुरुस्त करते, आणि ज्या फाइलच्या जागी काही ठेवते ती आधी बॅकअप फोल्डरमध्ये हलवते. `ccam <name>` शेअर लिंकच्या समस्येमुळे थांबल्यावर, फाइल्स हाताने परत आणल्यानंतर किंवा Claude Code अपडेट झाल्यानंतर ते चालवा.

## वापर

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

पहिली ओळ सर्व खात्यांना एकाच वेळी चिन्ह देते, त्यामुळे बाकी काही वाचण्यापूर्वीच समस्या कुठे आहेत ते तुम्हाला दिसते. त्याखाली प्रत्येक खात्याला एक box मिळतो, आणि कोणत्याही खात्याशी संबंधित नसलेले निष्कर्ष (Claude Code, environment variables, plugins, MCP सर्व्हर, शॉर्टकट्स, फोल्डर नियम, गट आणि सेटिंग्ज, `auto_switch`) `General` box मध्ये जातात. box च्या शीर्षकात खात्याचे सर्वात वाईट चिन्ह आणि त्यात किती समस्या आहेत ते दिसते. आतली प्रत्येक ओळ तपासलेले एक क्षेत्र आहे, नेहमी `Folder`, `Login` आणि `Share links` याच क्रमाने, मग इतर कोणतीही, चिन्हासह: `✔` ठीक, `ℹ` माहिती, `⚠` इशारा, `✘` त्रुटी. काहीही चूक नसलेला फोल्डर `✔ Folder  OK` दाखवतो.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix` सह, दुरुस्त झालेला प्रत्येक निष्कर्ष `→ fixed` ने संपतो, आणि तो आता समस्या म्हणून मोजला जात नाही:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

टर्मिनलवर boxes त्याच्या रुंदीत बसतात आणि लांब संदेश पुढच्या ओळीत जातात, टर्मिनलच्या स्वतःच्या रंगांत (`NO_COLOR` ते बंद करते). box कधीही 29 columns पेक्षा अरुंद नसतो, त्यामुळे अरुंद टर्मिनलवर कोणताही शब्द कापला जात नाही. pipe किंवा फाइलकडे पाठवल्यास आउटपुटमध्ये रंग नसतात आणि कोणताही संदेश तोडला किंवा कापला जात नाही, त्यामुळे `ccam doctor | grep settings.json` संपूर्ण ओळ शोधते; तेव्हा प्रत्येक box सर्वात लांब ओळीइतका रुंद असतो. डॅशबोर्डचा Health टॅब हेच boxes दाखवतो.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--fix` | सुरक्षितपणे दुरुस्त करता येईल ते दुरुस्त करा (बदललेल्या फाइल्सचा बॅकअप घेतला जातो) | बंद |
| `--dry-run` (global) | प्रत्येक तपासणी चालवा पण काहीही दुरुस्त करू नका, `--fix` असला तरी | बंद |

`ccam doctor` कोणतेही आर्ग्युमेंट घेत नाही. ते `--dry-run` पाळते: `ccam --dry-run doctor --fix` `ccam doctor` चे निष्कर्ष छापते, `--fix` काय बाजूला ठेवेल ते जोडते, आणि काहीही बदलत नाही:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

किमान एक `✘` निष्कर्ष दुरुस्त न होता उरला तर exit code `1` असतो, नाहीतर `0`. इशारे आणि माहितीमुळे run कधीच अयशस्वी होत नाही. `--fix` शिवाय run अयशस्वी झाल्यास ccam एक रिकामी ओळ आणि `` `ccam doctor --fix` repairs what it can. `` ने संपते. लोड न होऊ शकणारा config देखील `1` आणि `ccam:` त्रुटीने संपतो.

## हे कसे काम करते

`internal/cli/doctor.go` `~/.config/ccam/config.toml` लोड करते आणि `internal/doctor/doctor.go` मधील `doctor.Run` बोलावते, जे तपासण्या एका निश्चित क्रमाने चालवते आणि निष्कर्ष परत देते, प्रत्येकासोबत तो ज्या खात्याचा आहे ते खाते आणि त्याने तपासलेले क्षेत्र. `tui.HealthReport` त्यांना boxes मध्ये गटबद्ध करते. `--fix` `Options.Fix` म्हणून पुढे दिला जातो, जो dry run मध्ये false असतो.

<!-- motion:doctor-checks -->

मधला भाग प्रत्येक खात्यासाठी एकदा चालतो. तपासण्या, आणि प्रत्येकीसाठी `--fix` काय करते:

1. **Claude Code आवृत्ती** (`checkClaude`). launcher प्रमाणेच `claude` binary शोधते आणि `claude --version` चालवते. binary नसणे, executable फाइल नसलेला `claude_bin`, किंवा `--version` अयशस्वी होणे ही त्रुटी आहे. दुरुस्ती नाही.
2. **environment मधील credentials** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` किंवा `ANTHROPIC_AUTH_TOKEN` सेट असेल तर इशारा देते: त्या shell मधून सुरू होणाऱ्या प्रत्येक खात्याला ते वारसाने मिळते, आणि Claude Code खात्याच्या स्वतःच्या लॉगिनपेक्षा त्याला प्राधान्य देते. दुरुस्ती नाही; ते स्वतः unset करा.
3. **फोल्डर** (`checkProfile`). खात्याचा फोल्डर नसणे ही त्रुटी आहे आणि त्यानंतर त्या खात्याच्या तपासण्या थांबतात. group किंवा other bits सेट असलेल्या परवानग्या म्हणजे इशारा; `--fix` `chmod 0700` चालवते. `~/.claude.json` शेजारी स्वतःची `.claude.json` असलेला इम्पोर्ट केलेला `~/.claude` ही माहिती आहे: तेव्हा `ccam main` आणि साधे `claude` वेगवेगळे लॉगिन तपशील आणि MCP सर्व्हर ठेवतात. Windows वर परवानगी bits नसतात, त्यामुळे तिथे ही तपासणी आणि खालची `0600` तपासणी वगळल्या जातात; तुमच्या user profile मधील फोल्डर त्याच्या access list ने सुरक्षित असतो.
4. **लॉगिन फाइल** (`checkCredentials`). `.credentials.json` `lstat` ने तपासली जाते. symlink ही त्रुटी आहे, कारण प्रत्येक खात्याचे स्वतःचे token असणे आवश्यक आहे, आणि ती दुरुस्त होत नाही. `0600` शिवाय इतर mode म्हणजे इशारा; `--fix` `chmod 0600` चालवते. मग लॉगिन वाचले जाते (macOS वर आधी Keychain मधून). लॉगिन नसणे, किंवा लॉगिन नसलेली लॉगिन फाइल, हा इशारा आहे, जो `ccam login <name>` कडे निर्देश करतो. कालबाह्य refresh token, किंवा सर्व्हरने नकार दिल्यानंतर Claude Code ने काढून टाकलेले लॉगिन, ही `` login expired → `ccam login <name>` `` त्रुटी आहे; कालबाह्य access token ही फक्त माहिती आहे, कारण Claude Code पुढच्या वेळी सुरू होताना ते नूतनीकरण करते. Windows Credential Manager मध्ये ठेवलेले लॉगिन ही माहिती आहे, कारण ccam ते वाचू शकत नाही. वैध लॉगिन, फोल्डर सध्या ज्या email ने लॉगिन आहे तो email आणि plan दाखवते (`login valid (dev@company.com, Claude Team · Acme)`). यांपैकी काहीही दुरुस्त होत नाही; `ccam login --expired` गरज असलेल्या प्रत्येक खात्यात लॉगिन करते.
5. **शेअर लिंक्स** (`share.Check`). `share.Layout` मधील प्रत्येक पाथची खात्याच्या गटांशी तुलना करते, फक्त `lstat` आणि `readlink` वापरून. ते पाच प्रकार नोंदवते: नसलेली लिंक, दुसरीकडे निर्देश करणारी लिंक, लिंकच्या जागी असलेली खरी फाइल, आता local असलेल्या आयटमसाठी उरलेली लिंक, आणि जिचा गट फोल्डरच नाही अशी योग्य लिंक (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; launch हे तपासत नाही). नसलेली लिंक आणि नसलेला गट फोल्डर हे इशारे आहेत, बाकीच्या त्रुटी. गटातून परत होणारी पण अर्ध्यावर थांबलेली कॉपीही नोंदवली जाते: `--fix` तयार करून ठेवलेली पूर्ण प्रत जागेवर ठेवते, पण त्यानंतर तो पाथ (launch मुळे) नव्याने तयार झाला असेल, तर तुम्ही स्वतः निस्तरावे म्हणून ते पूर्ण प्रत असलेल्या फोल्डरचे नाव सांगते. ccam च्या नंतरच्या आवृत्तीत एखाद्या आयटममध्ये जोडलेले पाथ (history साठी `file-history/`, `tasks/` आणि `paste-cache/`, settings साठी `keybindings.json`) खात्याकडे अजून स्वतःची प्रत असेपर्यंत फक्त इशारा असतात: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. त्यांच्यामुळे launch कधीच थांबत नाही.
6. **अज्ञात नोंदी** (`checkUnknownEntries`). खात्याच्या फोल्डरमधील ccam ला माहीत नसलेली (`share.Known`) top-level नावे दाखवते, Claude Code चा layout बदलला असल्याची सूचना म्हणून. फक्त माहिती.
7. **Plugin पाथ** (`checkPlugins`). Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) absolute पाथ साठवतात. प्रत्येक गटाच्या `plugins/` मध्ये, एखाद्या खात्याच्या फोल्डरमधून जाणारे पाथ ते खाते काढल्यावर किंवा त्याचे नाव बदलल्यावर तुटतात. `--fix` ते `profile.RewritePluginPaths` ने गट फोल्डरकडे पुन्हा लिहिते, registry फाइलचे atomic write करून.
8. **MCP सर्व्हर** (`checkMCP`). प्रत्येक गटाची `mcp.json` एकदा लोड करते. वाचता न येणारी फाइल ही त्रुटी आहे; `mcp.Validate` मध्ये अयशस्वी होणारा सर्व्हर हा इशारा आहे. साध्या `claude` च्या `~/.claude.json` मधील, कोणतेही खाते लोड करत नाही असे सर्व्हर हा इशारा आहे, जो ते कॉपी करणारी `ccam migrate --mcp` कमांड सांगतो (अनेक गट MCP शेअर करत असतील तेव्हा `--group <a|b>` सह). गटातील फक्त एकाच खात्याकडे असलेले सर्व्हर ही माहिती आहे, ते शेअर करणाऱ्या `ccam mcp import` कमांडसह, आणि इतर खाती ज्या MCP सर्व्हरमध्ये साइन इन आहेत त्यांत साइन इन नसलेल्या खात्याला `/mcp` ने साइन इन करण्याचा इशारा मिळतो (ccam कधीही साइन-इन्स कॉपी करत नाही). दुरुस्ती नाही.
9. **Usage response format** (`checkUsageFormat`). cache शिवाय, एकामागून एक खात्यासाठी usage endpoint ला विचारते, जोपर्यंत एखादे उत्तर देत नाही, आणि हे सर्व 4 सेकंदांच्या एकाच मर्यादेत. `usage_source = "claude"` सह ते त्याऐवजी एका खात्यासाठी Claude Code ला विचारते. response आता parse होत नसेल, तर ते इशारा देते की `ccam usage` स्टेटस लाइन cache वर अवलंबून राहील; कोणत्याही खात्याने उत्तर दिले नाही, तर ते माहिती म्हणून तसे सांगते (`could not check the usage response format: …`). नेटवर्क वापरणारी ही एकमेव तपासणी आहे.
10. **शॉर्टकट** (`checkShims`). फक्त एखाद्या खात्याला शॉर्टकट असेल तेव्हा. `~/.local/bin` मध्ये नसलेल्या लिंक्स म्हणजे इशारा, आणि ccam ने बनवलेले पण अजूनही जुने ccam, म्हणजे बदललेले किंवा हलवलेले ccam, चालवणारे शॉर्टकटही (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` `shim.Install` बोलावते, जे फक्त ccam ने बनवलेल्या आणि `~/.local/state/ccam/shortcuts.json` मध्ये नोंदवलेल्या लिंक्स तयार करते आणि नव्या करते, आणि इतर फाइल्सना कधीच हात लावत नाही. त्यातील प्रत्येक लिंक तयार झाली तरच इशारा `→ fixed` ने संपतो; एखादी तयार करता आली नाही, तर कारण जोडले जाते. अजूनही इन्स्टॉल असलेले दुसरे ccam चालवणारा शॉर्टकट चालतो, आणि ती माहिती असते; `--fix` तो या ccam कडे वळवते. `PATH` वर नसलेला bin फोल्डर हा इशारा आहे आणि दुरुस्त होत नाही.
11. **फोल्डर नोंदी** (`checkFolders`). फक्त [`ccam folder`](/features/folder) च्या नोंदी असतील तेव्हा. ज्या नोंदीचे खाते अस्तित्वात नाही ती त्रुटी आहे, आणि ज्या नोंदीचा फोल्डर नाहीसा झाला आहे तो इशारा आहे. दोन्हीपैकी काहीही दुरुस्त होत नाही, कारण नोंद काढल्याने तिथे कोणते खाते सुरू होते ते बदलेल.
12. **गट आणि सेटिंग्ज** (`checkGroups`). खाते नसलेला गट सदस्य, ज्याचे नाव नंतर ccam कमांड बनले आहे असे खाते (ते `ccam run <name>` ने सुरू करा किंवा त्याचे नाव बदला), अज्ञात `theme` किंवा `usage_source`, आणि ccam वाचत नाही अशा `config.toml` मधील ओळी हे इशारे आहेत. शेवटच्या प्रकारात सहसा `[[group]]`, `[[profile]]` किंवा `[[folder]]` table च्या खाली लिहिलेली सेटिंग असते, जिथे ती त्या table ची ठरते (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); ते दुरुस्त होईपर्यंत config बदलणाऱ्या कमांड्स नकार देतात. दुरुस्ती नाही.
13. **लॉगिन** (`checkIdentities`). प्रत्येक फोल्डर कोणत्या खात्याने लॉगिन आहे ते त्याच्या `.claude.json` मधून वाचते. `ccam login` ने नोंदवलेल्या खात्याऐवजी आता दुसऱ्या खात्यात लॉगिन असलेला फोल्डर (उदाहरणार्थ session मध्ये `/login` केल्यानंतर) हा इशारा आहे; लॉगिनची नोंद नसलेल्या खात्याची, जसे इम्पोर्ट केलेले किंवा जुन्या ccam ने बनवलेले खाते, तुलना त्याच्या नोंदवलेल्या email वरून होते. दुसऱ्या फोल्डरच्याच Claude खात्यात लॉगिन असलेला फोल्डरही इशारा आहे: ती दोन्ही आपल्या वापर मर्यादा शेअर करतात. दुसऱ्या फोल्डरसारखाच लॉगिन token असलेला फोल्डर, म्हणजे कॉपी केलेला फोल्डर (साध्या `claude` चा `~/.claude` सुद्धा), ही त्रुटी आहे, कारण जो आधी त्याचे नूतनीकरण करेल तो दुसऱ्याला लॉग आउट करतो; त्याला स्वतःचे लॉगिन देण्यासाठी `ccam login <name>` चालवा. दुरुस्ती नाही.
14. **`auto_switch`** (`checkAutoSwitch`). फक्त `auto_switch = true` असेल तेव्हा. `General` box मधील माहिती `auto_switch` काय करते ते सांगते आणि [त्याच्या वर्णनाची](/features/next#auto-switch) लिंक देते. खात्यांमध्ये सीट्स असलेल्या प्रत्येक plan साठी (Claude Team, Claude Enterprise) एक ओळ त्या खात्यांची नावे देते, `auto_switch` एकाच organization च्या सीट्सदरम्यान सत्र हलवते हे सांगते, मदत केंद्राच्या लिंकसह Anthropic चा सीटविषयीचा नियम उद्धृत करते, आणि निर्णय तुमच्यावर सोडते. 2.1.214 पेक्षा जुने Claude Code, जे hooks ना आपली process (`CLAUDE_PID`) सांगत नाही, तिथे इशारा आहे, कारण कोणतेही सत्र हलवता येत नाही. ccam plugin नसलेल्या, किंवा plugin बंद असलेल्या खात्याला त्याच्या स्वतःच्या box मध्ये इशारा मिळतो, Claude Code मध्ये ते दुरुस्त करणाऱ्या कमांडसह: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (किंवा `/plugin enable ccam@ccam`). API key खाती वगळली जातात. Windows वर एकच ओळ असते: `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` दुरुस्ती नाही.

खात्यांच्या आधी, अर्ध्यावर थांबलेला **गटांचा बदल** ही `General` box मधील त्रुटी आहे: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` तो ठरलेल्या क्रमाने, त्याच्या conflicts ना दिलेल्या उत्तरांसह पूर्ण करते; तोपर्यंत गटांचा दुसरा कोणताही बदल, नाव बदलणे किंवा काढणे सुरू होत नाही. तो पूर्ण होईपर्यंत त्यातील खात्यांच्या फक्त नसलेल्या लिंक्स तयार केल्या जातात. तो पूर्ण करता येत नसेल, तर त्रुटी सांगते की ती फाइल काढल्यास बदलाचा उरलेला भाग सोडून दिला जातो. `ccam --dry-run doctor --fix` क्रम (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) आणि प्रत्येक खात्याचा भाग काय बाजूला ठेवेल ते दाखवते. दुसरे ccam आत्ता करत असलेला बदल ही माहिती आहे, आणि तो चालू असताना `--fix` कोणत्याही शेअर लिंक्स दुरुस्त करत नाही (`its share links are not repaired while another ccam changes the groups; …`). दुरुस्ती करताना `--fix` `~/.local/state/ccam/group-change.lock` धरून ठेवते, आणि बदल पूर्ण केल्यानंतर `config.toml` पुन्हा वाचते.

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

इथे `settings.json` ही खरी फाइल आहे. गटाकडे आधीच वेगळी `settings.json` असल्यामुळे conflict त्याचा डीफॉल्ट पर्याय `KeepGroup` घेतो, त्यामुळे गटाची प्रत राहते. खात्याची फाइल खात्याच्या फोल्डरमधील `.ccam-backup/<time>/settings.json` मध्ये हलते, आणि `os.Symlink` तो पाथ गटाशी लिंक करते. `--fix` तसे सांगते: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. इतर स्थिती:

- गटाकडे अजून प्रत नसेल, तर खात्याची फाइल गटात हलते आणि त्याचा पहिला मजकूर बनते.
- फोल्डरसाठी, गटाकडे नसलेल्या फाइल्स गटात हलतात. गटाकडे byte-न-byte तशाच असलेल्या फाइल्स टाकून दिल्या जातात, आणि वेगळ्या असलेल्या फाइल्स बॅकअपमध्ये जातात, ज्यांची नावे एका टिपेत दिली जातात (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- नसलेला गट फोल्डर आधी तयार केला जातो.
- `history.jsonl` साठी डीफॉल्ट `merge by line` आहे: गटाच्या ओळी, मग तुमच्या नवीन ओळी. गटाची आधीची प्रत `.ccam-backup/<time>/group/` मध्ये जतन होते.
- दुसऱ्या गटाकडे जाणारी लिंक बदलली जाते. ccam बाहेरील फोल्डरकडे जाणारी लिंक बॅकअपमध्ये जाते, आणि गट रिकामा असेल तर तिचा मजकूर गटाचा पहिला मजकूर बनतो.
- आता local असलेल्या आयटमसाठी उरलेली लिंक काढली जाते, आणि खात्याला गटाच्या मजकुराची प्रत मिळते.

लिंक्स काढणे आणि गटाकडे आधीच byte-न-byte तशाच असलेल्या फाइल्स टाकून देणे सोडल्यास, प्रत्येक टप्पा डेटा गटात किंवा बॅकअप फोल्डरमध्ये हलवतो. काहीही गमावले जात नाही.

## सुरक्षितता आणि अपवादात्मक स्थिती

- साधे `ccam doctor` कधीच काहीही बदलत नाही. scripts मध्ये खात्री करण्यासाठी `--fix` सोबत `--dry-run` वापरा.
- चालू असलेल्या खात्यासाठी, Claude Code कसेही सुरू केलेले असो, `doctor --fix` फक्त नसलेल्या लिंक्स आणि नसलेले गट फोल्डर बनवते, ज्यात काहीही गमावले जात नाही. आणखी काही उरले असेल, तर ते `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` जोडते, कारण दुरुस्ती ज्या फाइल्स हलवेल त्यांत session लिहीत राहते.
- शेअर दुरुस्ती अयशस्वी झाली, तर ccam त्या खात्यासाठी `could not repair share links: <error>` नोंदवते आणि त्याची अज्ञात नोंदींची तपासणी वगळते. इतर खाती तरीही तपासली जातात.
- `ccam group` बदल अर्ध्यावर अयशस्वी झाला, तर नवीन गट आधीच सेव्ह झालेले असतात आणि बदलाची नोंद `~/.local/state/ccam/group-change.json` मध्ये होते; ccam `ccam doctor --fix` ने पूर्ण करायला सांगते, जे MCP च्या पायऱ्यांसह त्याचा उरलेला भाग पूर्ण करते.
- symlink केलेली `.credentials.json`, कालबाह्य लॉगिन आणि environment मधील credentials नोंदवले जातात पण कधीच बदलले जात नाहीत.
- Plugin registry चे पुनर्लेखन atomic पद्धतीने लिहिले जाते, पण `.ccam-backup/` मध्ये कॉपी होत नाही.

## संबंधित

- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [आधीचा सेटअप कॉपी करणे](/features/migrate)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [कमांड्स](/docs/commands)
