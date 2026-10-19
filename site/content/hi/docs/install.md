---
title: ccam इंस्टॉल करें
description: ccam को Linux या macOS पर एक कमांड से, .deb पैकेज से या सोर्स से इंस्टॉल करें, और ccam update से इसे अप-टू-डेट रखें। पहले Claude Code इंस्टॉल होना चाहिए।
updated: 2026-10-19
order: 1
translation: machine
sourceHash: bde486b82b2e
---

Linux या macOS पर terminal में यह चलाएँ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

यह स्क्रिप्ट आपके सिस्टम के लिए रिलीज़ डाउनलोड करती है, उसका SHA-256 checksum जाँचती है और `ccam` को `~/.local/bin` में रखती है; ccam बाद में आपके चुने हुए शॉर्टकट कमांड, जैसे `ccwork`, भी इसी फ़ोल्डर में बनाता है। अगर यह फ़ोल्डर आपके `PATH` में नहीं है, तो स्क्रिप्ट वह लाइन प्रिंट करती है जिसे `~/.bashrc` या `~/.zshrc` में जोड़ना है।

ccam आपके लिए Claude Code शुरू करता है, इसलिए पहले [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) इंस्टॉल करें और पक्का करें कि `claude --version` काम करता है।

## समर्थित सिस्टम

- x86-64 और ARM64 पर Linux: Ubuntu, Debian, Fedora, Arch, Alpine और दूसरे डिस्ट्रीब्यूशन। बाइनरी static है और उसकी कोई dependency नहीं है।
- Apple silicon और Intel पर macOS।
- Windows सीधे समर्थित नहीं है। WSL के अंदर ccam किसी भी Linux सिस्टम की तरह काम करता है।

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

**कोई ख़ास वर्ज़न।** `CCAM_VERSION` सेट करें:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion `ccam` की सभी कमांड और विकल्पों के लिए काम करता है, अकाउंट के नामों सहित।

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

ccam कभी अपने आप अपडेट की जाँच नहीं करता; वह नई रिलीज़ तभी खोजता है जब आप यह कमांड चलाते हैं। कोई खास रिलीज़ इंस्टॉल करने के लिए, चाहे वह पुरानी ही क्यों न हो, `--version` दें। `ccam --dry-run update` बिना कुछ बदले दिखाता है कि क्या डाउनलोड होगा और क्या बदला जाएगा।

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

इंस्टॉल स्क्रिप्ट दोबारा चलाने से भी `~/.local/bin` वाली कॉपी अपडेट हो जाती है। अगर आपने ccam को `go install` से बनाया था, तो वही कमांड दोबारा चलाएँ।

**Debian और Ubuntu पैकेज।** `ccam update` apt से इंस्टॉल किए गए `/usr/bin/ccam` को नहीं बदलता, बल्कि उसकी जगह ये चरण प्रिंट करता है। [रिलीज़ पेज](https://github.com/jhyunwoo/ccam/releases) से नया `ccam_<version>_linux_amd64.deb` डाउनलोड करें और उसे पुराने पैकेज के ऊपर इंस्टॉल करें:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## अनइंस्टॉल करें

पहले अकाउंट शॉर्टकट हटाएँ, जब तक ccam इंस्टॉल है। `ccam list` उन्हें Shortcut कॉलम में दिखाता है; हर एक को `ccam shortcut <account> --remove` से हटाएँ, या वे लिंक `~/.local/bin` से खुद मिटा दें। फिर कमांड हटाएँ, और अगर अब अकाउंट्स की ज़रूरत नहीं है तो डेटा फ़ोल्डर भी:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` हटाने से उसमें रखे लॉगिन और बातचीत का इतिहास भी हट जाते हैं। जो `~/.claude` फ़ोल्डर आपने [import किया था](/docs/commands), वह अपनी जगह पर बना रहता है।
