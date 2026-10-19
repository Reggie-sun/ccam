---
title: ccam इंस्टॉल करें
description: ccam को Linux या macOS पर एक कमांड से, .deb पैकेज से या सोर्स से इंस्टॉल करें। पहले Claude Code इंस्टॉल होना चाहिए।
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Linux या macOS पर terminal में यह चलाएँ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

यह स्क्रिप्ट आपके सिस्टम के लिए रिलीज़ डाउनलोड करती है, उसका SHA-256 checksum जाँचती है और `ccam` को `~/.local/bin` में रखती है; बाद में `ccam add` इसी फ़ोल्डर में `claude1`, `claude2`… शॉर्टकट बनाता है। अगर यह फ़ोल्डर आपके `PATH` में नहीं है, तो स्क्रिप्ट वह लाइन प्रिंट करती है जिसे `~/.bashrc` या `~/.zshrc` में जोड़ना है।

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

## अनइंस्टॉल करें

कमांड और अकाउंट शॉर्टकट हटाएँ, फिर अगर अब अकाउंट्स की ज़रूरत नहीं है तो डेटा फ़ोल्डर भी हटाएँ:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` हटाने से उसमें रखे लॉगिन और बातचीत का इतिहास भी हट जाते हैं। जो `~/.claude` फ़ोल्डर आपने [import किया था](/docs/commands), वह अपनी जगह पर बना रहता है।
