---
title: ccam इन्स्टॉल करा
description: Linux किंवा macOS वर ccam एका कमांडने, .deb पॅकेजमधून किंवा सोर्सवरून इन्स्टॉल करा आणि ccam update ने ते अद्ययावत ठेवा. आधी Claude Code इन्स्टॉल असणे आवश्यक आहे.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: bde486b82b2e
---

Linux किंवा macOS वर टर्मिनलमध्ये हे चालवा:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

स्क्रिप्ट तुमच्या सिस्टमसाठी रिलीज डाउनलोड करते, त्याचा SHA-256 checksum तपासते आणि `ccam` ला `~/.local/bin` मध्ये ठेवते; तुम्ही नंतर निवडलेल्या शॉर्टकट कमांड्स, जसे `ccwork`, ccam याच फोल्डरमध्ये तयार करते. हा फोल्डर तुमच्या `PATH` मध्ये नसेल, तर `~/.bashrc` किंवा `~/.zshrc` मध्ये जोडायची ओळ स्क्रिप्ट प्रिंट करते.

ccam तुमच्यासाठी Claude Code सुरू करते, म्हणून आधी [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) इन्स्टॉल करा आणि `claude --version` चालते याची खात्री करा.

## समर्थित सिस्टिम्स

- x86-64 आणि ARM64 वरील Linux: Ubuntu, Debian, Fedora, Arch, Alpine आणि इतर distributions. बायनरी static आहे आणि तिला कोणत्याही dependencies लागत नाहीत.
- Apple silicon आणि Intel वरील macOS.
- Windows ला थेट समर्थन नाही. WSL मध्ये ccam इतर कोणत्याही Linux सिस्टिमप्रमाणे काम करते.

## इन्स्टॉल करण्याचे इतर मार्ग

**Debian आणि Ubuntu पॅकेज.** [रिलीज पान](https://github.com/jhyunwoo/ccam/releases) वरून `ccam_<version>_linux_amd64.deb` डाउनलोड करा आणि ते इन्स्टॉल करा:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

हे पॅकेज `/usr/bin/ccam` इन्स्टॉल करते.

**सोर्सवरून.** Go 1.26 किंवा नवीन आवृत्तीसह:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**विशिष्ट आवृत्ती.** `CCAM_VERSION` सेट करा:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## शेल completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion `ccam` च्या सर्व कमांड आणि पर्यायांसाठी, खात्यांच्या नावांसह, काम करते.

## इन्स्टॉलेशन तपासा

```sh
ccam version
ccam doctor
```

`ccam doctor` तपासते की Claude Code सापडते आणि प्रत्येक खात्याचा फोल्डर, लिंक आणि लॉगिन व्यवस्थित आहे. पुढे [क्विक स्टार्ट](/docs/quick-start) वाचा.

## अपडेट

```sh
ccam update --check
ccam update
```

`ccam update --check` फक्त नवीन रिलीज उपलब्ध आहे का ते सांगते. `ccam update` तुमच्या सिस्टमसाठी रिलीज डाउनलोड करते, रिलीजमधील `checksums.txt` शी त्याचा SHA-256 checksum जुळवून पाहते आणि मग एकाच टप्प्यात `ccam` फाइल बदलते. Checksum जुळला नाही, तर काहीही बदलले जात नाही. `ccwork` सारखे शॉर्टकट ccam च्या लिंक्स असल्यामुळे ते लगेच नवीन आवृत्ती चालवतात.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam स्वतःहून कधीही अपडेट तपासत नाही; तुम्ही ही कमांड चालवता तेव्हाच ते नवीन रिलीज शोधते. जुन्या रिलीजसह कोणताही विशिष्ट रिलीज इन्स्टॉल करण्यासाठी `--version` द्या. `ccam --dry-run update` काहीही न बदलता काय डाउनलोड होईल आणि काय बदलले जाईल ते दाखवते.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

इन्स्टॉल स्क्रिप्ट पुन्हा चालवल्यानेही `~/.local/bin` मधील प्रत अपडेट होते. तुम्ही ccam `go install` ने बिल्ड केले असेल, तर तीच कमांड पुन्हा चालवा.

**Debian आणि Ubuntu पॅकेज.** apt ने इन्स्टॉल केलेले `/usr/bin/ccam` `ccam update` बदलत नाही; त्याऐवजी ते पुढील पायऱ्या प्रिंट करते. [रिलीज पान](https://github.com/jhyunwoo/ccam/releases) वरून नवीन `ccam_<version>_linux_amd64.deb` डाउनलोड करा आणि जुन्या पॅकेजवरच ते इन्स्टॉल करा:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## अनइन्स्टॉल

ccam इन्स्टॉल असतानाच आधी खात्यांचे शॉर्टकट काढा. `ccam list` ते Shortcut स्तंभात दाखवते; प्रत्येक `ccam shortcut <account> --remove` ने काढा, किंवा त्या लिंक्स `~/.local/bin` मधून स्वतः हटवा. नंतर कमांड काढा, आणि खाती आता लागत नसतील तर डेटा फोल्डरही:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` हटवल्याने तिथे साठवलेले लॉगिन आणि संभाषण इतिहासही हटतात. तुम्ही [इम्पोर्ट केलेला](/docs/commands) `~/.claude` फोल्डर आहे तिथेच राहतो.
