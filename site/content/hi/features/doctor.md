---
title: जाँच और मरम्मत
description: ccam doctor से साझा लिंक, फ़ाइल अनुमतियाँ, लॉगिन और Claude Code का वर्ज़न जाँचें, और ccam doctor --fix से वह ठीक करें जिसे ठीक करना सुरक्षित है।
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` आपके ccam इंस्टॉलेशन की जाँच करता है: Claude Code बाइनरी, हर अकाउंट का फ़ोल्डर, लॉगिन और साझा लिंक, साझा MCP सर्वर, usage endpoint और आपकी शॉर्टकट कमांड। `ccam doctor --fix` वह ठीक करता है जो डेटा खोए बिना किया जा सकता है, और जिस फ़ाइल को बदलता है उसे पहले एक बैकअप फ़ोल्डर में move करता है। इसे तब चलाएँ जब `ccam <name>` किसी साझा लिंक की समस्या के कारण रुक जाए, जब आप फ़ाइलें हाथ से restore करें, या Claude Code अपडेट होने के बाद।

## उपयोग

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

पहली पंक्ति हर अकाउंट को एक साथ निशान देती है, ताकि कुछ और पढ़ने से पहले ही आप देख लें कि समस्याएँ कहाँ हैं। उसके नीचे हर अकाउंट को एक box मिलता है, और जो नतीजे किसी अकाउंट से जुड़े नहीं हैं (Claude Code, environment variables, प्लगइन, MCP सर्वर, शॉर्टकट, फ़ोल्डर नियम, ग्रुप और सेटिंग्स, `auto_switch`) वे `General` box में जाते हैं। box का शीर्षक अकाउंट का सबसे बुरा निशान और उसकी समस्याओं की संख्या दिखाता है। अंदर की हर पंक्ति एक जाँचा गया क्षेत्र है, हमेशा पहले `Folder`, `Login` और `Share links` इसी क्रम में, फिर बाक़ी, एक निशान के साथ: `✔` ठीक, `ℹ` जानकारी, `⚠` चेतावनी, `✘` त्रुटि। जिस फ़ोल्डर में कोई गड़बड़ी नहीं है, वह `✔ Folder  OK` दिखाता है।

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

`--fix` के साथ, ठीक किया गया हर नतीजा `→ fixed` पर ख़त्म होता है, और उसे अब समस्या नहीं गिना जाता:

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

टर्मिनल पर box उसकी चौड़ाई में समा जाते हैं और लंबे संदेश अगली पंक्ति में चले जाते हैं, टर्मिनल के अपने रंगों में (`NO_COLOR` उन्हें बंद करता है)। कोई box कभी 29 कॉलम से संकरा नहीं होता, ताकि संकरे टर्मिनल पर कोई शब्द न कटे। pipe या फ़ाइल में भेजने पर आउटपुट में कोई रंग नहीं होता और कोई संदेश तोड़ा या काटा नहीं जाता, इसलिए `ccam doctor | grep settings.json` पूरी पंक्ति ढूँढ लेता है; तब हर box सबसे लंबी पंक्ति जितना चौड़ा होता है। डैशबोर्ड का Health टैब यही box दिखाता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--fix` | वह ठीक करें जिसे ठीक करना सुरक्षित है (बदली गई फ़ाइलों का बैकअप लिया जाता है) | बंद |
| `--dry-run` (ग्लोबल) | हर जाँच चलाएँ लेकिन कुछ भी ठीक न करें, `--fix` के साथ भी | बंद |

`ccam doctor` कोई आर्ग्युमेंट नहीं लेता। यह `--dry-run` को मानता है: `ccam --dry-run doctor --fix` `ccam doctor` के नतीजे प्रिंट करता है, उनमें यह जोड़ता है कि `--fix` क्या अलग रखता, और कुछ नहीं बदलता:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

जब कम से कम एक `✘` नतीजा बिना ठीक हुए रह जाता है, तो exit code `1` होता है, वरना `0`। चेतावनियाँ और जानकारी कभी रन को विफल नहीं करतीं। जब `--fix` के बिना रन विफल होता है, तो ccam एक ख़ाली पंक्ति और `` `ccam doctor --fix` repairs what it can. `` के साथ ख़त्म होता है। जो config लोड नहीं हो सकता, वह भी `1` और एक `ccam:` त्रुटि के साथ ख़त्म होता है।

## यह कैसे काम करता है

`internal/cli/doctor.go`, `~/.config/ccam/config.toml` लोड करता है और `internal/doctor/doctor.go` में `doctor.Run` को बुलाता है, जो जाँचों को एक तय क्रम में चलाता है और नतीजे लौटाता है, हर एक उस अकाउंट के साथ जिससे वह जुड़ा है और उस क्षेत्र के साथ जिसे उसने जाँचा। `tui.HealthReport` उन्हें box में बाँटता है। `--fix` को `Options.Fix` के रूप में आगे दिया जाता है, जो dry run में false होता है।

<!-- motion:doctor-checks -->

बीच वाला हिस्सा हर अकाउंट के लिए एक बार चलता है। जाँचें, और हर एक के लिए `--fix` क्या करता है:

1. **Claude Code वर्ज़न** (`checkClaude`)। `claude` बाइनरी को उसी तरह खोजता है जैसे launcher खोजता है, और `claude --version` चलाता है। बाइनरी न मिलना, ऐसा `claude_bin` जो executable फ़ाइल नहीं है, या `--version` का विफल होना त्रुटि है। कोई fix नहीं।
2. **Environment में credentials** (`checkAuthEnv`)। जब `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` या `ANTHROPIC_AUTH_TOKEN` सेट हो, तो चेतावनी देता है: उस shell से शुरू होने वाला हर अकाउंट इसे विरासत में पाता है, और Claude Code इसे अकाउंट के अपने लॉगिन से ज़्यादा प्राथमिकता देता है। कोई fix नहीं; इसे ख़ुद unset करें।
3. **फ़ोल्डर** (`checkProfile`)। अकाउंट फ़ोल्डर का न होना त्रुटि है और उस अकाउंट की जाँचें वहीं ख़त्म हो जाती हैं। अगर अनुमतियों में group या other का कोई bit सेट हो, तो चेतावनी है; `--fix`, `chmod 0700` चलाता है। import किया गया `~/.claude` जिसकी `~/.claude.json` के बगल में अपनी `.claude.json` हो, जानकारी है: तब `ccam main` और सादा `claude` अलग लॉगिन विवरण और MCP सर्वर रखते हैं। Windows में permission bits नहीं होते, इसलिए वहाँ यह जाँच और नीचे वाली `0600` जाँच छोड़ दी जाती हैं; आपकी user profile में मौजूद फ़ोल्डर उसकी access list से सुरक्षित रहता है।
4. **लॉगिन फ़ाइल** (`checkCredentials`)। `.credentials.json` को `lstat` से जाँचा जाता है। Symlink त्रुटि है, क्योंकि हर अकाउंट का अपना token होना चाहिए, और इसे ठीक नहीं किया जाता। `0600` के अलावा कोई मोड चेतावनी है; `--fix`, `chmod 0600` चलाता है। फिर लॉगिन पढ़ा जाता है (macOS पर पहले कीचेन से)। लॉगिन न होना, या बिना लॉगिन वाली लॉगिन फ़ाइल, एक चेतावनी है जो `ccam login <name>` की ओर इशारा करती है। एक्सपायर हुआ refresh token, या ऐसा लॉगिन जिसे सर्वर के मना करने के बाद Claude Code ने हटा दिया, त्रुटि `` login expired → `ccam login <name>` `` है; एक्सपायर हुआ access token केवल जानकारी है, क्योंकि Claude Code अगली बार शुरू होने पर उसे renew कर देता है। Windows Credential Manager में रखा लॉगिन जानकारी है, क्योंकि ccam उसे पढ़ नहीं सकता। मान्य लॉगिन वह ईमेल दिखाता है जिससे फ़ोल्डर अभी लॉगिन है, और प्लान (`login valid (dev@company.com, Claude Team · Acme)`)। इनमें से कुछ भी ठीक नहीं किया जाता; `ccam login --expired` हर उस अकाउंट में लॉगिन करता है जिसे इसकी ज़रूरत है।
5. **साझा लिंक** (`share.Check`)। `share.Layout` के हर पाथ की तुलना अकाउंट के ग्रुप्स से करता है, केवल `lstat` और `readlink` का इस्तेमाल करके। यह पाँच तरह की समस्याएँ बताता है: लिंक का न होना, कहीं और इशारा करने वाला लिंक, लिंक की जगह असली फ़ाइल, ऐसे आइटम के लिए बचा हुआ लिंक जो अब local है, और सही लिंक जिसका ग्रुप फ़ोल्डर ग़ायब है (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; launch इसकी जाँच नहीं करता)। लिंक का न होना और ग्रुप फ़ोल्डर का ग़ायब होना चेतावनियाँ हैं, बाकी त्रुटियाँ हैं। ग्रुप से वापस कॉपी जो बीच में रुक गई, वह भी बताई जाती है: `--fix` तैयार रखी पूरी कॉपी को उसकी जगह रखता है, लेकिन जब पाथ तब से (किसी launch से) नया बन गया हो, तो यह पूरी कॉपी वाले फ़ोल्डर का नाम बताता है ताकि आप ख़ुद सुलझाएँ। ccam के किसी बाद के वर्ज़न में किसी आइटम में जोड़े गए पाथ (history के लिए `file-history/`, `tasks/` और `paste-cache/`, settings के लिए `keybindings.json`) तब तक केवल चेतावनी हैं जब तक अकाउंट के पास अभी भी अपनी कॉपी है: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`। इनके कारण launch कभी नहीं रुकता।
6. **अज्ञात entries** (`checkUnknownEntries`)। अकाउंट फ़ोल्डर के top-level नामों में से वे दिखाता है जिन्हें ccam नहीं जानता (`share.Known`), इस संकेत के रूप में कि Claude Code का लेआउट बदल गया है। केवल जानकारी।
7. **Plugin पाथ** (`checkPlugins`)। Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) absolute पाथ सेव करती हैं। हर ग्रुप के `plugins/` में, जो पाथ किसी एक अकाउंट के फ़ोल्डर से होकर जाते हैं, वे उस अकाउंट के हटने या नाम बदलने पर टूट जाते हैं। `--fix` उन्हें `profile.RewritePluginPaths` से ग्रुप फ़ोल्डर पर फिर से लिखता है, जो registry फ़ाइल की एक atomic write है।
8. **MCP सर्वर** (`checkMCP`)। हर ग्रुप की `mcp.json` एक बार लोड करता है। जो फ़ाइल पढ़ी न जा सके, वह त्रुटि है; जो सर्वर `mcp.Validate` में विफल हो, वह चेतावनी है। सादे `claude` के `~/.claude.json` में ऐसे सर्वर जिन्हें कोई अकाउंट लोड नहीं करता, एक चेतावनी हैं जो उन्हें कॉपी करने वाली `ccam migrate --mcp` कमांड बताती है (जब कई ग्रुप MCP साझा करते हों तो `--group <a|b>` के साथ)। ऐसे सर्वर जो किसी ग्रुप के केवल एक अकाउंट के पास हैं, उन्हें साझा करने वाली `ccam mcp import` कमांड के साथ जानकारी हैं, और जो अकाउंट उन MCP सर्वरों में साइन इन नहीं है जिनमें दूसरे अकाउंट हैं, उसे `/mcp` से साइन इन करने की चेतावनी मिलती है (ccam कभी साइन-इन कॉपी नहीं करता)। कोई fix नहीं।
9. **Usage response का format** (`checkUsageFormat`)। Cache के बिना, usage endpoint से एक के बाद एक अकाउंट के लिए पूछता है, जब तक कोई जवाब न दे, सब कुछ 4 सेकंड की एक ही सीमा में। `usage_source = "claude"` के साथ यह इसके बजाय एक अकाउंट के लिए Claude Code से पूछता है। अगर response अब parse नहीं होता, तो यह चेतावनी देता है कि `ccam usage` स्टेटस लाइन की cache पर निर्भर करेगा; जब किसी अकाउंट ने जवाब नहीं दिया, तो यह जानकारी के रूप में बताता है (`could not check the usage response format: …`)। नेटवर्क इस्तेमाल करने वाली यही एकमात्र जाँच है।
10. **शॉर्टकट** (`checkShims`)। केवल तब जब किसी अकाउंट का शॉर्टकट हो। `~/.local/bin` में लिंक न होना चेतावनी है, और ccam के बनाए ऐसे शॉर्टकट भी जो अभी भी किसी पुराने ccam को चलाते हैं, जिसे बदला या move किया गया (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix`, `shim.Install` को बुलाता है, जो केवल वही लिंक बनाता और नए करता है जिन्हें ccam ने बनाया और `~/.local/state/ccam/shortcuts.json` में दर्ज किया, और दूसरी फ़ाइलों को कभी नहीं छूता। चेतावनी `→ fixed` पर तभी ख़त्म होती है जब उनमें से हर लिंक बन गया हो; अगर कोई नहीं बन सका, तो कारण जोड़ा जाता है। ऐसा शॉर्टकट जो किसी दूसरे, अभी भी इंस्टॉल ccam को चलाता है, काम करता है, और जानकारी है; `--fix` उसे इस ccam की ओर कर देता है। जो bin फ़ोल्डर `PATH` पर नहीं है, वह चेतावनी है और उसे ठीक नहीं किया जाता।
11. **फ़ोल्डर** (`checkFolders`)। केवल तब जब [`ccam folder`](/features/folder) की entries मौजूद हों। जिस entry का अकाउंट मौजूद नहीं है, वह त्रुटि है, और जिस entry का फ़ोल्डर अब नहीं है, वह चेतावनी है। दोनों में से कोई ठीक नहीं किया जाता, क्योंकि किसी entry को हटाने से बदल जाएगा कि वहाँ कौन-सा अकाउंट शुरू होता है।
12. **ग्रुप और सेटिंग्स** (`checkGroups`)। ग्रुप का ऐसा सदस्य जो अकाउंट नहीं है, ऐसा अकाउंट जिसका नाम बाद में ccam कमांड बन गया है (इसे `ccam run <name>` से शुरू करें या इसका नाम बदलें), अज्ञात `theme` या `usage_source`, और `config.toml` की वे पंक्तियाँ जिन्हें ccam नहीं पढ़ता, चेतावनियाँ हैं। आख़िरी वाली आमतौर पर किसी `[[group]]`, `[[profile]]` या `[[folder]]` table के नीचे लिखी सेटिंग होती है, जहाँ वह उस table की होती है (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); config बदलने वाली कमांड तब तक मना करती हैं जब तक इसे ठीक न किया जाए। कोई fix नहीं।
13. **लॉगिन** (`checkIdentities`)। हर फ़ोल्डर की `.claude.json` से पढ़ता है कि वह किस रूप में लॉगिन है। जो फ़ोल्डर अब उस अकाउंट के बजाय किसी दूसरे अकाउंट में लॉगिन है जो `ccam login` ने दर्ज किया था (उदाहरण के लिए session के अंदर `/login` के बाद), वह चेतावनी है; लॉगिन रिकॉर्ड के बिना वाले अकाउंट की, जैसे import किया गया या किसी पुराने ccam से बना, तुलना उसके दर्ज ईमेल से होती है। जो फ़ोल्डर किसी दूसरे फ़ोल्डर वाले उसी Claude अकाउंट में लॉगिन है, वह भी चेतावनी है: दोनों अपनी उपयोग सीमाएँ साझा करते हैं। जो फ़ोल्डर किसी दूसरे वाला ही लॉगिन token रखता है, यानी कॉपी किया गया फ़ोल्डर (सादे `claude` का `~/.claude` भी), वह त्रुटि है, क्योंकि जो पहले उसे renew करता है वह दूसरे को लॉग आउट कर देता है; उसे अपना लॉगिन देने के लिए `ccam login <name>` चलाएँ। कोई fix नहीं।
14. **`auto_switch`** (`checkAutoSwitch`)। केवल तब जब `auto_switch = true` हो। `General` box में जानकारी बताती है कि `auto_switch` क्या करता है और [उसके विवरण](/features/next#auto-switch) का लिंक देती है। अकाउंट्स में सीटों वाले हर प्लान (Claude Team, Claude Enterprise) के लिए एक पंक्ति उनके नाम बताती है, कहती है कि `auto_switch` किसी session को एक संगठन की सीटों के बीच ले जाता है, सहायता केंद्र के लिंक के साथ Anthropic का सीट वाला नियम उद्धृत करती है, और चुनाव आप पर छोड़ती है। 2.1.214 से पुराना Claude Code, जो hooks को अपना process (`CLAUDE_PID`) नहीं बताता, वहाँ चेतावनी है, क्योंकि कोई session ले जाया नहीं जा सकता। ccam plugin के बिना, या plugin बंद वाले, अकाउंट को उसके अपने box में चेतावनी मिलती है, साथ में वह कमांड जो इसे Claude Code में ठीक करती है: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (या `/plugin enable ccam@ccam`)। API key वाले अकाउंट छोड़ दिए जाते हैं। Windows पर अकेली पंक्ति यह है: `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` कोई fix नहीं।

अकाउंट्स से पहले, बीच में रुका **ग्रुप का बदलाव** `General` box में एक त्रुटि है: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``। `--fix` उसे उसी क्रम में पूरा करता है जिसमें उसकी योजना बनी थी, उसके टकरावों के दिए गए जवाबों के साथ; तब तक ग्रुप का कोई और बदलाव, नाम बदलना या हटाना शुरू नहीं होता। उसमें शामिल अकाउंट्स के, उसके पूरा होने तक, केवल ग़ायब लिंक बनाए जाते हैं। अगर वह पूरा नहीं हो सकता, तो त्रुटि बताती है कि उस फ़ाइल को हटाने से बदलाव का बाक़ी हिस्सा छोड़ दिया जाता है। `ccam --dry-run doctor --fix` क्रम (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) और यह दिखाता है कि हर अकाउंट का हिस्सा क्या अलग रखता। जो बदलाव कोई दूसरा ccam अभी कर रहा है, वह जानकारी है, और उसके चलते `--fix` कोई साझा लिंक ठीक नहीं करता (`its share links are not repaired while another ccam changes the groups; …`)। मरम्मत करते समय `--fix` `~/.local/state/ccam/group-change.lock` रखता है, और बदलाव पूरा करने के बाद `config.toml` फिर से पढ़ता है।

### साझा लिंक की मरम्मत

शुरू करने वाला रास्ता गायब लिंक को तुरंत बना देता है, लेकिन जिस भी चीज़ से डेटा खो सकता है, उस पर रुक जाता है और यहाँ की ओर इशारा करता है:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

किसी भी साझा लिंक समस्या के लिए `--fix`, अकाउंट की साझा सेटिंग्स को मौजूदा और लक्ष्य दोनों स्थिति मानकर `share.Build` को बुलाता है, फिर `Plan.Apply` को। `Build` हर चरण को डिस्क पर मौजूद चीज़ों से फिर से तय करता है, वही योजना जो `ccam group` का कोई बदलाव इस्तेमाल करता है, इसलिए मरम्मत उन्हीं नियमों का पालन करती है जो किसी ग्रुप में शामिल होने पर लागू होते हैं।

<!-- motion:doctor-repair -->

यहाँ `settings.json` एक असली फ़ाइल है। क्योंकि ग्रुप के पास पहले से एक अलग `settings.json` है, टकराव में डिफ़ॉल्ट विकल्प `KeepGroup` चुना जाता है, इसलिए ग्रुप की कॉपी बनी रहती है। अकाउंट की फ़ाइल अकाउंट फ़ोल्डर के अंदर `.ccam-backup/<time>/settings.json` में move होती है, और `os.Symlink` पाथ को ग्रुप से लिंक करता है। `--fix` यह बताता है: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`। दूसरी स्थितियाँ:

- अगर ग्रुप के पास अभी कोई कॉपी नहीं है, तो अकाउंट की फ़ाइल ग्रुप में move होती है और उसकी पहली सामग्री बनती है।
- फ़ोल्डर के लिए, जो फ़ाइलें ग्रुप में नहीं हैं वे ग्रुप में move होती हैं। जो फ़ाइलें ग्रुप में byte-दर-byte वैसी ही हैं, वे हटा दी जाती हैं, और जो अलग हैं वे बैकअप में जाती हैं, जिनके नाम एक नोट में होते हैं (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`)।
- ग़ायब ग्रुप फ़ोल्डर पहले बनाया जाता है।
- `history.jsonl` के लिए डिफ़ॉल्ट `merge by line` है: ग्रुप की पंक्तियाँ, फिर आपकी नई पंक्तियाँ। ग्रुप की पिछली कॉपी `.ccam-backup/<time>/group/` में सेव होती है।
- किसी दूसरे ग्रुप का लिंक बदल दिया जाता है। ccam के बाहर किसी फ़ोल्डर का लिंक बैकअप में जाता है, और अगर ग्रुप ख़ाली है तो उसकी सामग्री ग्रुप की शुरुआती सामग्री बनती है।
- ऐसे आइटम के लिए बचा हुआ लिंक जो अब local है, हटा दिया जाता है, और अकाउंट को ग्रुप की सामग्री की एक कॉपी मिलती है।

लिंक हटाने और ग्रुप में byte-दर-byte पहले से मौजूद फ़ाइलें हटाने के अलावा, हर चरण डेटा को ग्रुप या बैकअप फ़ोल्डर में move करता है। कुछ भी खोता नहीं।

## सुरक्षा और विशेष स्थितियाँ

- सादा `ccam doctor` कभी कुछ नहीं बदलता। स्क्रिप्ट में पक्का करने के लिए `--fix` के साथ `--dry-run` इस्तेमाल करें।
- चल रहे अकाउंट के लिए, चाहे Claude Code किसी भी तरह शुरू हुआ हो, `doctor --fix` केवल ग़ायब लिंक और ग़ायब ग्रुप फ़ोल्डर बनाता है, जिससे कुछ नहीं खोता। जब कुछ और बाक़ी रहे, तो यह `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` जोड़ता है, क्योंकि session उन फ़ाइलों में लिखता रहता है जिन्हें मरम्मत हटाती।
- अगर साझा लिंक की मरम्मत विफल हो, तो ccam उस अकाउंट के लिए `could not repair share links: <error>` बताता है और उसकी अज्ञात entries वाली जाँच छोड़ देता है। दूसरे अकाउंट्स की जाँच फिर भी होती है।
- अगर `ccam group` का कोई बदलाव बीच में विफल हो, तो नए ग्रुप पहले ही सेव हो चुके होते हैं और बदलाव `~/.local/state/ccam/group-change.json` में दर्ज होता है; ccam बताता है कि `ccam doctor --fix` से काम पूरा करें, जो उसका बाक़ी हिस्सा, MCP चरणों समेत, पूरा करता है।
- Symlink वाली `.credentials.json`, एक्सपायर हुआ लॉगिन और environment के credentials बताए जाते हैं, लेकिन कभी बदले नहीं जाते।
- Plugin registry को फिर से लिखना atomic रूप से होता है, लेकिन उसकी कॉपी `.ccam-backup/` में नहीं जाती।

## संबंधित

- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [कमांड](/docs/commands)
