---
title: ccam इन्स्टॉल करा
metaTitle: "Linux, macOS आणि Windows वर ccam इन्स्टॉल करा"
description: "Linux, macOS किंवा Windows वर ccam एका कमांडने, .deb पॅकेज किंवा सोर्सवरून इन्स्टॉल करा, आणि ccam update ने अपडेट करा. आधी Claude Code इन्स्टॉल असावे."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "07fab3165030"
---

Linux किंवा macOS वर टर्मिनलमध्ये हे चालवा:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

स्क्रिप्ट तुमच्या सिस्टमसाठी रिलीज डाउनलोड करते, त्याचा SHA-256 checksum तपासते आणि `ccam` ला `~/.local/bin` मध्ये ठेवते; तुम्ही नंतर निवडलेल्या शॉर्टकट कमांड्स, जसे `ccwork`, ccam याच फोल्डरमध्ये तयार करते. हा फोल्डर तुमच्या `PATH` मध्ये नसेल, तर `~/.bashrc` किंवा `~/.zshrc` मध्ये जोडायची ओळ स्क्रिप्ट प्रिंट करते.

Windows वर हे PowerShell मध्ये चालवा:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

ते `ccam_windows_<arch>.zip` डाउनलोड करते, त्याचा SHA-256 checksum तपासते, `ccam.exe` ला `%USERPROFILE%\.local\bin` मध्ये ठेवते, म्हणजे त्याच फोल्डरमध्ये जिथे Claude Code चा installer `claude.exe` ठेवतो, आणि तो फोल्डर तुमच्या user `PATH` मध्ये नसेल तर जोडते. त्यानंतर नवीन टर्मिनल उघडा.

ccam तुमच्यासाठी Claude Code सुरू करते, म्हणून आधी [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) इन्स्टॉल करा आणि `claude --version` चालते याची खात्री करा.

## समर्थित सिस्टिम्स

- x86-64 आणि ARM64 वरील Linux: Ubuntu, Debian, Fedora, Arch, Alpine आणि इतर distributions. बायनरी static आहे आणि तिला कोणत्याही dependencies लागत नाहीत.
- Apple silicon आणि Intel वरील macOS.
- x86-64 आणि ARM64 वरील Windows 10 आणि 11, PowerShell, Windows Terminal किंवा Git Bash मध्ये. शेअर केलेले फोल्डर प्रत्येक user साठी चालतात. शेअर केलेल्या फाइल्सना (`settings.json`, `CLAUDE.md`, `history.jsonl`) Developer Mode (Settings > System > For developers) लागतो, कारण Windows फक्त या mode मध्ये, किंवा administrator ला, फाइल्सच्या लिंक्स तयार करू देते; [Windows वर शेअरिंग](/features/share#on-windows) पाहा.
- WSL मध्ये ccam इतर कोणत्याही Linux सिस्टिमप्रमाणे काम करते.

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

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## शेल completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion `ccam` च्या सर्व कमांड आणि पर्यायांसाठी, खात्यांच्या नावांसह, काम करते.

## Claude Code प्लगइन

ccam प्लगइन Claude Code ला स्वतःच ccam वापरायला शिकवते: कोणत्या खात्यात जागा आहे ते तपासणे, `ccam pick` ने वेगवेगळ्या खात्यांवर समांतर `claude -p` workers सुरू करणे, आणि एखादे skill किंवा MCP सर्व्हर का सापडत नाही हे `ccam ext` ने समजावणे. ते Claude Code मध्येच इन्स्टॉल करा:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

ते एक skill, `/ccam:usage`, `/ccam:ext` आणि `/ccam:switch` या कमांड, आणि दोन hooks जोडते: सत्र वापर मर्यादेवर थांबले की तुम्हाला terminal मध्ये सूचना मिळते, आणि `/ccam:switch [account]` संभाषण दुसऱ्या खात्याकडे सोपवते ([तपशील](/features/next#inside-claude-code)). hooks `ccam` चालवतात, त्यामुळे ते Claude Code ला दिसणाऱ्या `PATH` वर असले पाहिजे. हे skill लॉगिन करणाऱ्या किंवा फाइल्स हलवणाऱ्या कमांड कधीच चालवत नाही; त्याऐवजी ते त्या सुचवते.

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

**Windows.** `ccam update` Windows साठीचा `.zip` डाउनलोड करते. Windows वर चालू असलेला प्रोग्राम बदलता येत नाही, त्यामुळे जुनी बायनरी `ccam.exe.old` असे नाव बदलली जाते, जी पुढचा अपडेट हटवतो.

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

Windows वर, PowerShell मध्ये:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` हटवल्याने तिथे साठवलेले लॉगिन आणि संभाषण इतिहासही हटतात. तुम्ही [इम्पोर्ट केलेला](/docs/commands) `~/.claude` फोल्डर आहे तिथेच राहतो.
