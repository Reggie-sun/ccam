---
title: ccam इंस्टॉल करें
metaTitle: "Linux, macOS और Windows पर ccam इंस्टॉल करें"
description: "ccam को Linux, macOS या Windows पर एक कमांड, .deb पैकेज या go install से इंस्टॉल करें, और ccam update से अपडेट करें। पहले Claude Code इंस्टॉल होना चाहिए।"
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam आपके लिए Claude Code शुरू करता है, इसलिए पहले [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) इंस्टॉल करें और पक्का करें कि `claude --version` काम करता है।

Linux या macOS पर terminal में यह चलाएँ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

यह स्क्रिप्ट आपके सिस्टम के लिए रिलीज़ डाउनलोड करती है, उसका SHA-256 checksum जाँचती है और `ccam` को `~/.local/bin` में रखती है; ccam बाद में आपके चुने हुए शॉर्टकट कमांड, जैसे `ccwork`, भी इसी फ़ोल्डर में बनाता है। अगर यह फ़ोल्डर आपके `PATH` में नहीं है, तो स्क्रिप्ट वह लाइन प्रिंट करती है जिसे `~/.bashrc` या `~/.zshrc` में जोड़ना है।

Windows पर इसे PowerShell में चलाएँ:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

यह `ccam_windows_<arch>.zip` डाउनलोड करता है, उसका SHA-256 checksum जाँचता है और `ccam.exe` को `%USERPROFILE%\.local\bin` में रखता है, यानी उसी फ़ोल्डर में जहाँ Claude Code का installer `claude.exe` रखता है। अगर वह फ़ोल्डर आपके user `PATH` में अभी नहीं है, तो स्क्रिप्ट उसे जोड़ देती है और `%USERPROFILE%\bin` जैसी entries को जस का तस रखती है। इसके बाद नया terminal खोलें।

## समर्थित सिस्टम

- x86-64 और ARM64 पर Linux: Ubuntu, Debian, Fedora, Arch, Alpine और दूसरे डिस्ट्रीब्यूशन। बाइनरी static है और उसकी कोई dependency नहीं है।
- Apple silicon और Intel पर macOS।
- x86-64 और ARM64 पर Windows 10 और 11, PowerShell, Windows Terminal या Git Bash में:
  - साझा फ़ोल्डर Developer Mode के बिना भी काम करते हैं।
  - साझा फ़ाइलों (`settings.json`, `CLAUDE.md`, `history.jsonl`) के लिए Developer Mode (Settings > System > For developers) चाहिए, क्योंकि Windows सिर्फ़ इसी mode में, या किसी administrator को, फ़ाइलों के लिंक बनाने देता है; [Windows पर साझा करना](/features/share#on-windows) देखें।
  - Developer Mode के बिना, पहला `ccam add` केवल MCP सर्वर और एक्सटेंशन वाला `default` ग्रुप बनाता है, और बताता है कि Developer Mode चालू होने के बाद बाक़ी चीज़ें कैसे जोड़ें।
- WSL के अंदर ccam वैसे ही काम करता है जैसे किसी भी दूसरे Linux सिस्टम पर।

## इंस्टॉल करने के दूसरे तरीके

**Debian और Ubuntu पैकेज।** [रिलीज़ पेज](https://github.com/jhyunwoo/ccam/releases) से `ccam_<version>_linux_amd64.deb` डाउनलोड करें और इसे इंस्टॉल करें:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

यह पैकेज `/usr/bin/ccam` इंस्टॉल करता है।

**सोर्स से।** Go 1.26 या नए वर्ज़न के साथ:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

अगर आप कोई रिलीज़ tag (`@latest` या `@vX.Y.Z`) इंस्टॉल करते हैं, तो ccam वही वर्ज़न बताता है, और `ccam update` और डैशबोर्ड का अपडेट ऑफ़र काम करते हैं। `@main` का, या किसी tag के बाद के commits वाले local checkout का build `dev` बताता है, और `ccam update` उसे तब तक मना करता है जब तक आप `--version` न दें।

**कोई ख़ास वर्ज़न।** `CCAM_VERSION` सेट करें:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion `ccam` की सभी कमांड और विकल्पों के लिए काम करता है, अकाउंट के नामों सहित।

## Claude Code प्लगइन

ccam प्लगइन ख़ुद Claude Code को ccam इस्तेमाल करना सिखाता है: यह जाँचना कि किस अकाउंट में गुंजाइश है, `ccam pick` से अलग-अलग अकाउंट्स पर समानांतर `claude -p` workers शुरू करना, और `ccam ext` से यह समझाना कि कोई skill या MCP सर्वर क्यों नहीं मिल रहा। इसे Claude Code के अंदर इंस्टॉल करें:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

यह जोड़ता है:

- एक skill, जो लॉगिन करने वाली या फ़ाइलें move करने वाली कमांड सुझाता है, लेकिन उन्हें कभी चलाता नहीं;
- कमांड `/ccam:usage`, `/ccam:ext` और `/ccam:switch`;
- दो hook ([विवरण](/features/next#inside-claude-code))। एक hook, जब कोई session उपयोग सीमा पर रुकता है, तो टर्मिनल नोटिफ़िकेशन दिखाता है; Linux और macOS पर, [`auto_switch`](/features/next#auto-switch) चालू होने पर, ccam फिर session को गुंजाइश वाले किसी दूसरे अकाउंट पर ले जाता है। दूसरा hook `/ccam:switch [account]` को बातचीत दूसरे अकाउंट को सौंपने देता है।

hook `ccam` चलाते हैं, इसलिए यह उस `PATH` पर होना चाहिए जो Claude Code देखता है।

## इंस्टॉलेशन जाँचें

```sh
ccam version
ccam doctor
```

`ccam doctor` जाँचता है कि Claude Code मिल रहा है और हर अकाउंट फ़ोल्डर, लिंक और लॉगिन ठीक है। इसके बाद [क्विक स्टार्ट](/docs/quick-start) पर जाएँ।

## अपडेट

```sh
ccam update --check
ccam update
```

`ccam update --check` केवल यह बताता है कि कोई नई रिलीज़ उपलब्ध है या नहीं। `ccam update` आपके सिस्टम के लिए रिलीज़ डाउनलोड करता है, रिलीज़ की `checksums.txt` से उसका SHA-256 checksum मिलाता है और फिर `ccam` फ़ाइल को एक ही चरण में बदल देता है। अगर checksum मेल नहीं खाता, तो कुछ भी नहीं बदला जाता। `ccwork` जैसे शॉर्टकट ccam के लिंक हैं, इसलिए वे तुरंत नया वर्ज़न चलाते हैं।

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

डैशबोर्ड (टर्मिनल में बिना arguments के `ccam`) भी हर बार खुलने पर नई रिलीज़ की जाँच करता है और अगर कोई हो तो उसे पेश करता है: `U` दबाएँ और पुष्टि करें। कोई और कमांड अपडेट की जाँच नहीं करती, और ccam बिना पूछे कभी अपडेट नहीं होता। कोई खास रिलीज़ इंस्टॉल करने के लिए, चाहे वह पुरानी ही क्यों न हो, `--version` दें। `ccam --dry-run update` बिना कुछ बदले दिखाता है कि क्या डाउनलोड होगा और क्या बदला जाएगा।

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

इंस्टॉल स्क्रिप्ट दोबारा चलाने से भी `~/.local/bin` वाली कॉपी अपडेट हो जाती है। अगर आपने ccam को किसी रिलीज़ tag से `go install` से बनाया था, तो आप `ccam update` चला सकते हैं, या वही `go install` कमांड दोबारा चला सकते हैं।

**Windows.** `ccam update` Windows के लिए `.zip` डाउनलोड करता है। Windows पर चल रहे प्रोग्राम को बदला नहीं जा सकता, इसलिए पुरानी बाइनरी का नाम बदलकर `ccam.exe.old` कर दिया जाता है, जिसे अगला अपडेट मिटा देता है; जब तक कोई session किसी पुराने `ccam.exe.old` से अभी चल रहा हो, नाम `ccam.exe.old-<time>` होता है। जो शॉर्टकट hard links हैं, वे तुरंत बाद नए कर दिए जाते हैं। `install.ps1` दोबारा चलाना भी इसी तरह काम करता है, तब भी जब कोई Claude Code session पुरानी बाइनरी से अभी चल रहा हो, और hard-link वाले शॉर्टकट नए करने पर यह `OK 2 shortcuts now run the new ccam` प्रिंट करता है।

**Debian और Ubuntu पैकेज।** `ccam update` apt से इंस्टॉल किए गए `/usr/bin/ccam` को नहीं बदलता, बल्कि उसकी जगह ये चरण प्रिंट करता है। [रिलीज़ पेज](https://github.com/jhyunwoo/ccam/releases) से नया `ccam_<version>_linux_amd64.deb` डाउनलोड करें और उसे पुराने पैकेज के ऊपर इंस्टॉल करें:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## अनइंस्टॉल करें

`ccam uninstall` ccam को और उसकी बनाई हर चीज़ को हटा देता है। यह पहले दिखाता है कि क्या हटेगा, और पूछता है:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- हर अकाउंट वैसे ही हटता है जैसे `ccam remove` हटाता है। ccam के फ़ोल्डरों में मौजूद अकाउंट पहले लॉग आउट होता है और फिर अपने फ़ोल्डर के साथ मिटा दिया जाता है, उसकी बातचीत का इतिहास और सेटिंग्स भी।
- जो `~/.claude` फ़ोल्डर आपने [import किया था](/docs/commands), वह सीधे `claude` के लिए अपनी जगह पर बना रहता है। अगर वह किसी ग्रुप में है, तो उसने जो कुछ साझा किया था उसकी एक कॉपी उसे वापस मिल जाती है। उसकी `settings.json` में ccam की स्टेटस लाइन की जगह वह स्टेटस लाइन लौट आती है जिसे ccam ने wrap किया था, अगर कोई थी। अगर उसमें ccam plugin है, तो ccam वह `/plugin uninstall` कमांड प्रिंट करता है जिसे Claude Code में चलाना है।
- शॉर्टकट कमांड, `~/.config/ccam`, `~/.local/share/ccam` (ccam के trash समेत) और `~/.local/state/ccam` मिटा दिए जाते हैं, फिर खुद `ccam` कमांड।

कोई अकाउंट चल रहा हो तो ccam मना कर देता है; पहले उसके Claude Code session बंद करें। apt से इंस्टॉल हुआ `/usr/bin/ccam` बना रहता है; उसे बाद में `sudo apt remove ccam` से हटाएँ। Windows पर चलता हुआ प्रोग्राम खुद को मिटा नहीं सकता, इसलिए `ccam uninstall` अंत में `ccam.exe` के लिए `Remove-Item` कमांड प्रिंट करता है।
