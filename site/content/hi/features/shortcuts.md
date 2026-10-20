---
title: शॉर्टकट कमांड
description: ccam shortcut से किसी अकाउंट को ccwork जैसी उसकी अपनी कमांड दें, और ccam shims install से ~/.local/bin में लिंक फिर से बनाएँ।
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

शॉर्टकट एक कमांड है जिसका नाम आप ख़ुद रखते हैं और जो Claude Code को एक अकाउंट से शुरू करती है: `ccwork --resume` वही करता है जो `ccam work --resume` करता है। यह `~/.local/bin` में ccam का एक लिंक है, इसलिए यह bash, zsh, स्क्रिप्ट और `PATH` से प्रोग्राम चलाने वाली हर चीज़ में काम करता है। शॉर्टकट वैकल्पिक हैं; उन अकाउंट्स के लिए बनाएँ जिन्हें आप अक्सर शुरू करते हैं।

## उपयोग

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam हर बदलाव की पुष्टि करता है:

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

एक अकाउंट का अधिकतम एक शॉर्टकट होता है, जो `~/.config/ccam/config.toml` में `shortcut` के रूप में सेव होता है। `ccam add` का विज़र्ड `Create a shortcut command for <name>? [y/N]` और फिर `Shortcut name [cc<name>]:` पूछता है। विज़र्ड के बाहर, `ccam add` केवल `--shortcut` के साथ शॉर्टकट बनाता है। `ccam list` हर अकाउंट का शॉर्टकट, या `-` दिखाता है।

शॉर्टकट हर आर्ग्युमेंट Claude Code को देता है, ठीक वैसे जैसे `ccam <name>` नाम के बाद देता है। इसमें `--dry-run` और `--help` भी शामिल हैं: शुरुआत का पूर्वावलोकन करने के लिए `ccam --dry-run work` इस्तेमाल करें।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--remove` | अकाउंट का शॉर्टकट और उसका लिंक मिटाएँ | बंद |
| `-h`, `--help` | कमांड की help | |
| `--dry-run` | `ccam shortcut` इसे मानता है: यह नाम जाँचता है और सेव या लिंक किए बिना `Would set the shortcut of work to ccwork.` या `Would remove ccwork, the shortcut of work.` प्रिंट करता है | बंद |

`ccam shims install` कोई आर्ग्युमेंट नहीं लेता और `--dry-run` को `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` के साथ अस्वीकार करता है।

नाम के नियम, इसी क्रम में जाँचे जाते हैं:

- 1 से 64 अक्षर, जिनमें `a-z`, `0-9`, `.`, `-` और `_` आते हैं, और जो किसी अक्षर या अंक से शुरू हों;
- `ccam`, `claude` या `run` या `usage` जैसी कोई ccam कमांड नहीं;
- किसी दूसरे अकाउंट का शॉर्टकट नहीं;
- `~/.local/bin` की ऐसी फ़ाइल नहीं जो पहले से ccam का लिंक न हो;
- आपके `PATH` पर किसी दूसरे executable का नाम नहीं।

कमांड और `--remove` एक साथ देने पर `give a command or --remove, not both` के साथ रुकता है। अकाउंट का पहले से मौजूद नाम सेट करने पर `work: ccwork (unchanged)` प्रिंट होता है और लिंक गायब हो तो फिर से बनता है।

## यह कैसे काम करता है

नाम सेट करने पर `checkShortcut` चलता है, जो bin फ़ोल्डर और ccam बाइनरी के असली पाथ के साथ `shim.Check` को बुलाता है:

<!-- motion:shortcuts-check -->

नाम क्रम से जाँचों में नीचे जाता है। हर जाँच के बगल में एक नाम है जिसे वह अस्वीकार करती है: `Work` pattern में विफल होता है, `run` एक ccam कमांड है, `ccside` किसी दूसरे अकाउंट का है, `ccfile`, `~/.local/bin` में मौजूद एक फ़ाइल है, और `git`, `/usr/bin/git` को छिपा देता। केवल हर जाँच पास करने वाला नाम ही सेव और लिंक होता है।

1. `shim.Check` नाम का मिलान `^[a-z0-9][a-z0-9._-]{0,63}$` से, फिर `config.Reserved` से, फिर `cfg.FindShortcut` से करता है। यह `~/.local/bin/<name>` पर `os.Lstat` बुलाता है, और `PATH` के हर दूसरे फ़ोल्डर के लिए `os.Stat` बुलाकर ऐसी executable फ़ाइल को अस्वीकार करता है जो ccam का लिंक नहीं है।
2. ccam `fsutil.AtomicWrite` (अस्थायी फ़ाइल, `fsync`, rename) से `config.toml` सेव करता है।
3. `shim.Remove` पुराने शॉर्टकट का लिंक मिटाता है, लेकिन केवल तभी जब वह ccam का लिंक हो।
4. `shim.Install` हर configured शॉर्टकट का गायब लिंक बनाता है, ccam बाइनरी के लिए `os.Symlink` से, जिसे `os.Executable` और `filepath.EvalSymlinks` से खोजा जाता है। यह कभी कुछ नहीं हटाता।

कोई लिंक तब ccam का अपना माना जाता है जब `os.Readlink` ccam का पाथ लौटाए, या जब लिंक और ccam पाथ एक ही फ़ाइल पर पहुँचें। हर हटाना और हर overwrite जाँच इसी परीक्षण से गुज़रती है, इसलिए ccam उन फ़ाइलों को नहीं छूता जो उसने नहीं बनाईं।

जब आप शॉर्टकट चलाते हैं, तो ccam को केवल आपका टाइप किया नाम चाहिए:

<!-- motion:shortcuts-link -->

Shell `PATH` पर `ccwork` खोजता है, kernel लिंक का पालन करके ccam बाइनरी शुरू करता है, और `argv[0]` अब भी `ccwork` होता है। `app.Main` में `ccam` के अलावा कोई भी नाम `config.toml` एक बार लोड करता है और `shim.Lookup` को बुलाता है, जो वह अकाउंट खोजता है जिसका `shortcut`, `filepath.Base(argv[0])` के बराबर हो। फिर `launcher.LaunchWith` पहले से लोड किए गए config के साथ वही चरण अपनाता है जो [`ccam <name>`](/features/launch) अपनाता है, और `argv[1:]` को बिना छुए `claude` को देता है।

Shell alias या wrapper स्क्रिप्ट की जगह symlink क्यों: symlink हर shell में और स्क्रिप्ट में काम करता है, यह कोई अतिरिक्त shell शुरुआत नहीं जोड़ता (लगभग 1 से 2 ms), और शुरुआत अब भी एक ही फ़ाइल पढ़ती है। `scripts/perf-check.sh`, `hyperfine` से शॉर्टकट को उसी 10 ms (p95) सीमा के सामने मापता है जो `ccam <name>` पर लागू है।

### Windows पर

शॉर्टकट `%USERPROFILE%\.local\bin` में `ccwork.exe` होता है, और `shim.Name` खोज से पहले `.exe` (किसी भी case में) हटा देता है, इसलिए `ccwork`, `ccwork.exe` और `CCWORK.EXE` सभी `work` शुरू करते हैं। ccam पहले symbolic link बनाता है, जिसकी Windows Developer Mode में या किसी administrator को अनुमति देता है। यह अधिकार न होने पर यह hard link बनाता है: उसी `ccam.exe` फ़ाइल का दूसरा नाम, जिसके लिए किसी विशेषाधिकार की ज़रूरत नहीं होती। hard link को ccam का अपना तब माना जाता है जब वह `ccam.exe` वाली ही फ़ाइल हो। क्योंकि `ccam update`, `ccam.exe` की जगह एक नई फ़ाइल रखता है, इसलिए फिर यह `shim.Relink` से हर hard link वाले शॉर्टकट को नया करता है और `✓ shortcut ccwork renewed` प्रिंट करता है। नाम की जाँच दूसरे `PATH` फ़ोल्डरों में `PATHEXT` के हर extension (`.exe`, `.cmd` वगैरह) के साथ `<name>` खोजती है।

## सुरक्षा और विशेष स्थितियाँ

- **ऐसा लिंक जिसे कोई अकाउंट नहीं मानता।** अगर `~/.local/bin` में ccam के किसी लिंक का नाम कोई अकाउंट इस्तेमाल नहीं करता, उदाहरण के लिए `config.toml` में `shortcut` को हाथ से बदलने के बाद, तो ccam डिफ़ॉल्ट अकाउंट पर नहीं लौटता। यह `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` के साथ रुकता है। `ccam shortcut` से नाम फिर से सेट करें, या लिंक मिटा दें।
- **मौजूदा फ़ाइलें।** ccam कभी ऐसी फ़ाइल नहीं बदलता जो उसने नहीं बनाई। `ccam shims install`, `<path> already exists and is not a ccam link, so it was left alone` के साथ रुकता है।
- **दूसरी जगह ले जाई गई ccam बाइनरी।** तब पुराने लिंक ऐसे पाथ की ओर इशारा करते हैं जो अब मौजूद नहीं है, और ccam उन्हें अपना नहीं मानता। उन्हें मिटाएँ, फिर `ccam shims install` चलाएँ।
- **नाम बदलना और हटाना।** `ccam rename` शॉर्टकट का नाम वही रखता है, और लिंक काम करता रहता है, क्योंकि वह ccam की ओर इशारा करता है और ccam नाम को `config.toml` में खोजता है। `ccam remove` अकाउंट का लिंक मिटाता है। `ccam shortcut <account> <new>` नया लिंक बनाने से पहले पुराना लिंक मिटाता है।
- **PATH।** जब ccam लिंक बनाता है और `~/.local/bin` आपके `PATH` पर नहीं है, तो यह `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` और उसके बाद जोड़ने वाली `export PATH=...` पंक्ति प्रिंट करता है। Windows पर यह इसके बजाय वह PowerShell पंक्ति प्रिंट करता है जो फ़ोल्डर को आपके user `PATH` में जोड़ती है।
- **बाद में जाँचना।** `ccam doctor`, `missing shortcuts in /home/me/.local/bin: ccwork` बताता है और फ़ोल्डर `PATH` पर न हो तो चेतावनी देता है। `ccam doctor --fix` गायब लिंक बनाता है।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [अकाउंट प्रबंधित करना](/features/accounts)
- [कमांड संदर्भ](/docs/commands#shortcuts)
- [ccam इंस्टॉल करें](/docs/install)
