---
title: जाँच और मरम्मत
description: ccam doctor से साझा लिंक, फ़ाइल अनुमतियाँ, लॉगिन और Claude Code का वर्ज़न जाँचें, और ccam doctor --fix से वह ठीक करें जिसे ठीक करना सुरक्षित है।
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` आपके ccam इंस्टॉलेशन की जाँच करता है: Claude Code बाइनरी, हर अकाउंट का फ़ोल्डर, लॉगिन और साझा लिंक, साझा MCP सर्वर, usage endpoint और आपकी शॉर्टकट कमांड। `ccam doctor --fix` वह ठीक करता है जो डेटा खोए बिना किया जा सकता है, और जिस फ़ाइल को बदलता है उसे पहले एक बैकअप फ़ोल्डर में move करता है। इसे तब चलाएँ जब `ccam <name>` किसी साझा लिंक की समस्या के कारण रुक जाए, जब आप फ़ाइलें हाथ से restore करें, या Claude Code अपडेट होने के बाद।

## उपयोग

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

हर नतीजा एक निशान वाली एक पंक्ति है: `✓` ठीक, `·` जानकारी, `!` चेतावनी, `✗` त्रुटि। अकाउंट से जुड़े नतीजे अकाउंट के नाम से शुरू होते हैं।

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

`--fix` के साथ, ठीक किया गया हर नतीजा `→ fixed` पर ख़त्म होता है:

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

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--fix` | वह ठीक करें जिसे ठीक करना सुरक्षित है (बदली गई फ़ाइलों का बैकअप लिया जाता है) | बंद |
| `--dry-run` (ग्लोबल) | हर जाँच चलाएँ लेकिन कुछ भी ठीक न करें, `--fix` के साथ भी | बंद |

`ccam doctor` कोई आर्ग्युमेंट नहीं लेता। यह `--dry-run` को मानता है: `ccam --dry-run doctor --fix` वही नतीजे प्रिंट करता है जो `ccam doctor` करता है, और कुछ नहीं बदलता।

जब कम से कम एक `✗` नतीजा बिना ठीक हुए रह जाता है, तो exit code `1` होता है, वरना `0`। चेतावनियाँ और जानकारी कभी रन को विफल नहीं करतीं। जब `--fix` के बिना रन विफल होता है, तो ccam एक ख़ाली पंक्ति और `` `ccam doctor --fix` repairs what it can. `` के साथ ख़त्म होता है। जो config लोड नहीं हो सकता, वह भी `1` और एक `ccam:` त्रुटि के साथ ख़त्म होता है।

## यह कैसे काम करता है

`internal/cli/doctor.go`, `~/.config/ccam/config.toml` लोड करता है और `internal/doctor/doctor.go` में `doctor.Run` को बुलाता है, जो जाँचों को एक तय क्रम में चलाता है और नतीजे लौटाता है। `--fix` को `Options.Fix` के रूप में आगे दिया जाता है, जो dry run में false होता है।

<!-- motion:doctor-checks -->

बीच वाला हिस्सा हर अकाउंट के लिए एक बार चलता है। जाँचें, और हर एक के लिए `--fix` क्या करता है:

1. **Claude Code वर्ज़न** (`checkClaude`)। `claude` बाइनरी को उसी तरह खोजता है जैसे launcher खोजता है, और `claude --version` चलाता है। अगर वर्ज़न `2.1.302` से शुरू नहीं होता, जो वह वर्ज़न है जिसके फ़ोल्डर लेआउट के साथ ccam की जाँच हुई है (`doctor.LayoutVersion`), तो यह `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` जोड़ता है। बाइनरी न मिलना या `--version` का विफल होना त्रुटि है। कोई fix नहीं।
2. **Environment में credentials** (`checkAuthEnv`)। जब `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` या `ANTHROPIC_AUTH_TOKEN` सेट हो, तो चेतावनी देता है: उस shell से शुरू होने वाला हर अकाउंट इसे विरासत में पाता है, और Claude Code इसे अकाउंट के अपने लॉगिन से ज़्यादा प्राथमिकता देता है। कोई fix नहीं; इसे ख़ुद unset करें।
3. **फ़ोल्डर** (`checkProfile`)। अकाउंट फ़ोल्डर का न होना त्रुटि है और उस अकाउंट की जाँचें वहीं ख़त्म हो जाती हैं। अगर अनुमतियों में group या other का कोई bit सेट हो, तो चेतावनी है; `--fix`, `chmod 0700` चलाता है। Windows में permission bits नहीं होते, इसलिए वहाँ यह जाँच और नीचे वाली `0600` जाँच छोड़ दी जाती हैं; आपकी user profile में मौजूद फ़ोल्डर उसकी access list से सुरक्षित रहता है।
4. **लॉगिन फ़ाइल** (`checkCredentials`)। `.credentials.json` को `lstat` से जाँचा जाता है। Symlink त्रुटि है, क्योंकि हर अकाउंट का अपना token होना चाहिए, और इसे ठीक नहीं किया जाता। `0600` के अलावा कोई मोड चेतावनी है; `--fix`, `chmod 0600` चलाता है। फिर लॉगिन पढ़ा जाता है (macOS पर पहले कीचेन से)। लॉगिन न होना एक चेतावनी है जो `ccam login <name>` की ओर इशारा करती है। एक्सपायर हुआ refresh token त्रुटि `login expired` है; एक्सपायर हुआ access token केवल जानकारी है, क्योंकि Claude Code अगली बार शुरू होने पर उसे renew कर देता है। इनमें से कुछ भी ठीक नहीं किया जाता।
5. **साझा लिंक** (`share.Check`)। `share.Layout` के हर पाथ की तुलना अकाउंट की साझा सेटिंग्स से करता है, केवल `lstat` और `readlink` का इस्तेमाल करके। यह चार तरह की समस्याएँ बताता है: लिंक का न होना, कहीं और इशारा करने वाला लिंक, लिंक की जगह असली फ़ाइल, और ऐसे आइटम के लिए बचा हुआ लिंक जो अब local है। लिंक का न होना चेतावनी है, बाकी त्रुटियाँ हैं।
6. **अज्ञात entries** (`checkUnknownEntries`)। अकाउंट फ़ोल्डर के top-level नामों में से वे दिखाता है जिन्हें ccam नहीं जानता (`share.Known`), इस संकेत के रूप में कि Claude Code का लेआउट बदल गया है। केवल जानकारी।
7. **Plugin पाथ** (`checkPlugins`)। Plugin registries (`installed_plugins.json`, `known_marketplaces.json`) absolute पाथ सेव करती हैं। हर ग्रुप के `plugins/` में, जो पाथ किसी एक अकाउंट के फ़ोल्डर से होकर जाते हैं, वे उस अकाउंट के हटने या नाम बदलने पर टूट जाते हैं। `--fix` उन्हें `profile.RewritePluginPaths` से ग्रुप फ़ोल्डर पर फिर से लिखता है, जो registry फ़ाइल की एक atomic write है।
8. **MCP सर्वर** (`checkMCP`)। हर ग्रुप की `mcp.json` एक बार लोड करता है। जो फ़ाइल पढ़ी न जा सके, वह त्रुटि है; जो सर्वर `mcp.Validate` में विफल हो, वह चेतावनी है। कोई fix नहीं।
9. **Usage response का format** (`checkUsageFormat`)। Cache के बिना और सामान्य 3 सेकंड की सीमा के साथ, usage endpoint से एक के बाद एक अकाउंट के लिए पूछता है, जब तक कोई जवाब न दे। अगर response अब parse नहीं होता, तो यह चेतावनी देता है कि `ccam usage` स्टेटस लाइन की cache पर निर्भर करेगा। नेटवर्क इस्तेमाल करने वाली यही एकमात्र जाँच है।
10. **शॉर्टकट** (`checkShims`)। केवल तब जब किसी अकाउंट का शॉर्टकट हो। `~/.local/bin` में लिंक न होना चेतावनी है; `--fix`, `shim.Install` को बुलाता है, जो केवल गायब लिंक बनाता है और दूसरी फ़ाइलों को कभी नहीं छूता। जो bin फ़ोल्डर `PATH` पर नहीं है, वह चेतावनी है और उसे ठीक नहीं किया जाता।
11. **फ़ोल्डर** (`checkFolders`)। केवल तब जब [`ccam folder`](/features/folder) की entries मौजूद हों। जिस entry का अकाउंट मौजूद नहीं है, वह त्रुटि है, और जिस entry का फ़ोल्डर अब नहीं है, वह चेतावनी है। दोनों में से कोई ठीक नहीं किया जाता, क्योंकि किसी entry को हटाने से बदल जाएगा कि वहाँ कौन-सा अकाउंट शुरू होता है।

### साझा लिंक की मरम्मत

शुरू करने वाला रास्ता गायब लिंक को तुरंत बना देता है, लेकिन जिस भी चीज़ से डेटा खो सकता है, उस पर रुक जाता है और यहाँ की ओर इशारा करता है:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

किसी भी साझा लिंक समस्या के लिए `--fix`, अकाउंट की साझा सेटिंग्स को मौजूदा और लक्ष्य दोनों स्थिति मानकर `share.Build` को बुलाता है, फिर `Plan.Apply` को। `Build` हर चरण को डिस्क पर मौजूद चीज़ों से फिर से तय करता है, वही योजना जो `ccam share` इस्तेमाल करता है, इसलिए मरम्मत उन्हीं नियमों का पालन करती है जो साझा करना चालू करने पर लागू होते हैं।

<!-- motion:doctor-repair -->

यहाँ `settings.json` एक असली फ़ाइल है। क्योंकि ग्रुप के पास पहले से एक अलग `settings.json` है, टकराव में डिफ़ॉल्ट विकल्प `KeepGroup` चुना जाता है, इसलिए ग्रुप की कॉपी बनी रहती है। अकाउंट की फ़ाइल अकाउंट फ़ोल्डर के अंदर `.ccam-backup/<time>/settings.json` में move होती है, और `os.Symlink` पाथ को ग्रुप से लिंक करता है। दूसरी स्थितियाँ:

- अगर ग्रुप के पास अभी कोई कॉपी नहीं है, तो अकाउंट की फ़ाइल ग्रुप में move होती है और उसकी पहली सामग्री बनती है।
- फ़ोल्डर के लिए, जो फ़ाइलें ग्रुप में नहीं हैं वे ग्रुप में move होती हैं; बाकी बैकअप में जाती हैं।
- `history.jsonl` के लिए डिफ़ॉल्ट `merge by line` है: ग्रुप की पंक्तियाँ, फिर आपकी नई पंक्तियाँ। ग्रुप की पिछली कॉपी `.ccam-backup/<time>/group/` में सेव होती है।
- किसी दूसरे ग्रुप का लिंक बदल दिया जाता है। ccam के बाहर किसी फ़ोल्डर का लिंक बैकअप में जाता है, और अगर ग्रुप ख़ाली है तो उसकी सामग्री ग्रुप की शुरुआती सामग्री बनती है।
- ऐसे आइटम के लिए बचा हुआ लिंक जो अब local है, हटा दिया जाता है, और अकाउंट को ग्रुप की सामग्री की एक कॉपी मिलती है।

लिंक हटाने के अलावा, हर चरण डेटा को ग्रुप या बैकअप फ़ोल्डर में move करता है। कुछ भी मिटाया नहीं जाता।

## सुरक्षा और विशेष स्थितियाँ

- सादा `ccam doctor` कभी कुछ नहीं बदलता। स्क्रिप्ट में पक्का करने के लिए `--fix` के साथ `--dry-run` इस्तेमाल करें।
- `doctor --fix` यह नहीं जाँचता कि कोई अकाउंट चल रहा है या नहीं। किसी अकाउंट के साझा लिंक ठीक करने से पहले उसका Claude Code बंद करें, क्योंकि session उन फ़ाइलों में लिखता रहता है।
- अगर साझा लिंक की मरम्मत विफल हो, तो ccam उस अकाउंट के लिए `could not repair share links: <error>` बताता है और उसकी अज्ञात entries वाली जाँच छोड़ देता है। दूसरे अकाउंट्स की जाँच फिर भी होती है।
- अगर `ccam share` बीच में विफल हो, तो यह पहले लक्ष्य सेटिंग्स सेव करता है और बताता है कि `ccam doctor --fix` से काम पूरा करें।
- Symlink वाली `.credentials.json`, एक्सपायर हुआ लॉगिन और environment के credentials बताए जाते हैं, लेकिन कभी बदले नहीं जाते।
- Plugin registry को फिर से लिखना atomic रूप से होता है, लेकिन उसकी कॉपी `.ccam-backup/` में नहीं जाती।

## संबंधित

- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [कमांड](/docs/commands)
