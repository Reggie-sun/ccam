---
title: ccam इन्स्टॉल करा
description: Linux किंवा macOS वर ccam एका कमांडने, .deb पॅकेजमधून किंवा सोर्सवरून इन्स्टॉल करा. आधी Claude Code इन्स्टॉल असणे आवश्यक आहे.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux किंवा macOS वर टर्मिनलमध्ये हे चालवा:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ही स्क्रिप्ट तुमच्या सिस्टिमसाठीचे रिलीज डाउनलोड करते, त्याचा SHA-256 checksum तपासते आणि `~/.local/bin` मध्ये दोन कमांड ठेवते: `ccam` आणि त्याचे छोटे नाव `ccs`. तो फोल्डर तुमच्या `PATH` वर नसेल, तर `~/.bashrc` किंवा `~/.zshrc` मध्ये जोडायची ओळ स्क्रिप्ट प्रिंट करते.

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

हे पॅकेज `/usr/bin/ccam` आणि `/usr/bin/ccs` लिंक इन्स्टॉल करते.

**सोर्सवरून.** Go 1.26 किंवा नवीन आवृत्तीसह:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

Completion `ccam` आणि `ccs` दोन्हींसाठी, खात्यांच्या नावांसह, काम करते.

## इन्स्टॉलेशन तपासा

```sh
ccam version
ccam doctor
```

`ccam doctor` तपासते की Claude Code सापडते आणि प्रत्येक खात्याचा फोल्डर, लिंक आणि लॉगिन व्यवस्थित आहे. पुढे [क्विक स्टार्ट](/docs/quick-start) वाचा.

## अनइन्स्टॉल

दोन कमांड आणि खात्यांचे शॉर्टकट काढून टाका, आणि खाती आता लागत नसतील तर डेटा फोल्डरही काढा:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` हटवल्याने तिथे साठवलेले लॉगिन आणि संभाषण इतिहासही हटतात. तुम्ही [इम्पोर्ट केलेला](/docs/commands) `~/.claude` फोल्डर आहे तिथेच राहतो.
