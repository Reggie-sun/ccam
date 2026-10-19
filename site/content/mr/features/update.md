---
title: ccam अपडेट करणे
description: ccam update SHA-256 तपासणीनंतर ccam च्या जागी GitHub वरील सर्वात नवीन रिलीज ठेवते; --check फक्त कळवते आणि --version दिलेला रिलीज tag इन्स्टॉल करते.
updated: 2026-10-19
order: 10
translation: machine
sourceHash: 83fc20191328
---

`ccam update` GitHub वरून ccam रिलीज डाउनलोड करते, तो रिलीजच्या `checksums.txt` शी तपासते आणि चालू binary एका rename ने बदलते. तुम्ही ccam install script ने किंवा archive मधून इन्स्टॉल केले असेल तेव्हा ते वापरा. ccam कधीच स्वतःहून अपडेट्स तपासत नाही: फक्त हीच कमांड GitHub शी संपर्क करते, त्यामुळे इतर प्रत्येक कमांड offline सुद्धा तशीच चालते.

## वापर

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

नवीन रिलीज नसेल तेव्हा ccam `ccam 0.4.0 is up to date.` छापते आणि काहीही बदलत नाही. `--version` सह तुम्ही आधीच चालवत असलेली आवृत्ती दिल्यास ते `ccam 0.4.0 is already installed.` छापते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--check` | फक्त नवीन रिलीज आहे का ते कळवा | बंद |
| `--version <tag>` | `v0.3.0` सारखा हा रिलीज tag इन्स्टॉल करा, सध्याच्यापेक्षा नवीन किंवा जुना | latest |
| `--dry-run` | पाळला जातो: सर्वात नवीन tag शोधतो (`--version` ने एखादा दिला नसेल तर), binary वर त्याच तपासण्या चालवतो, आणि डाउनलोड न करता download, checksum आणि replace पाथ छापतो | बंद |

`--check` आणि `--version` एकत्र वापरता येत नाहीत. tag `vMAJOR.MINOR.PATCH` सारखा असणे आवश्यक आहे, ऐच्छिक pre-release किंवा build suffix सह; सुरुवातीचा `v` वगळता येतो. इतर काहीही अयशस्वी होते, उदाहरणार्थ `--version takes a release tag such as v0.3.0, not "1.2"`.

`CCAM_RELEASES_URL` हा environment variable mirror साठी किंवा tests साठी `https://github.com/jhyunwoo/ccam/releases` ची जागा घेतो. त्याची रचना तशीच असणे आवश्यक आहे: `<url>/latest` `<url>/tag/<tag>` कडे redirect करतो, आणि assets `<url>/download/<tag>/<asset>` वर असतात.

## हे कसे काम करते

कोड `internal/update` मध्ये आहे, जो फक्त `ccam update` import करते. सुरू करण्याचा मार्ग ते कधीच बोलावत नाही.

1. `cli.updateCmd` `os.Executable` ने चालू binary शोधते आणि `filepath.EvalSymlinks` ने symlinks सोडवते. हाच खरा पाथ बदलला जातो.
2. `--version` शिवाय, `Updater.Latest` redirect न पाळता `<releases>/latest` ची request करते, आणि `Location` header (`…/tag/v0.5.0`) मधून tag वाचते. हे GitHub API ऐवजी releases page वापरते, त्यामुळे API ची तासाला 60 requests ची मर्यादा लागू होत नाही. या शोधाला 10 सेकंद असतात.
3. `parseVersion` tag एका कडक pattern शी तपासते. tag download URLs मध्ये आणि छापलेल्या shell कमांडमध्ये जातो, त्यामुळे `/`, `%`, spaces किंवा shell अक्षरे असलेले काहीही पुढे जात नाही. `compareVersions` त्याची चालू आवृत्तीशी तुलना करते; pre-release त्याच्या रिलीजच्या आधी येतो.
4. `managed()` आणि `writable()` ठरवतात की ccam binary बदलू शकते की नाही (खाली पाहा).
5. `Updater.Fetch` `ccam_<os>_<arch>.tar.gz` आणि `checksums.txt` एकाच वेळी डाउनलोड करते, दोन्हीसाठी 2 मिनिटे. प्रत्येक डाउनलोड 256 MB पर्यंत मर्यादित असतो.
6. इतर काहीही होण्यापूर्वी archive चा SHA-256 `checksums.txt` मधील त्याच्या ओळीशी जुळणे आवश्यक आहे. ओळ नसल्यास किंवा checksum वेगळा असल्यास अपडेट `nothing was changed` सह थांबतो.
7. `extract()` tar वाचते आणि `ccam` नावाची पहिली regular फाइल घेते. फोल्डर, लिंक्स आणि इतर नोंदी दुर्लक्षित केल्या जातात, आणि binary सुद्धा 256 MB पर्यंत मर्यादित असते.
8. `Updater.Install` नवीन binary `fsutil.AtomicWrite` ने लिहिते: binary शेजारी `.ccam.tmp*` ही तात्पुरती फाइल, जुन्या फाइलच्या permission bits सह, डिस्कवर flush करून जुन्या फाइलवर rename केलेली.

tag शोधण्यापासून rename पर्यंत संपूर्ण pipeline:

<!-- motion:update-pipeline -->

archive आणि `checksums.txt` शेजारी शेजारी डाउनलोड होतात, आणि checksum जुळेपर्यंत काहीही लिहिले जात नाही. rename atomic आहे, त्यामुळे crash किंवा भरलेली डिस्क एकतर जुनी binary किंवा नवीन binary ठेवते, कधीच अर्धवट नाही. `~/.local/bin/ccwork` सारख्या शॉर्टकट कमांड्स ccam binary च्या symlinks आहेत, त्यामुळे त्या कोणत्याही बदलाशिवाय नवीन आवृत्ती चालवतात.

### पॅकेजने इन्स्टॉल केलेले

`.deb` मधून इन्स्टॉल केलेले ccam apt चे असते, आणि apt च्या नकळत ते बदलल्यास package database चुकीचा राहील. डाउनलोड करण्यापूर्वी ccam binary कोणाची आहे ते तपासते:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` किंवा `/sbin` मधील binary पॅकेजने सांभाळलेली मानली जाते.
2. इतरत्र, `dpkg` इन्स्टॉल असेल आणि `dpkg -S <path>` यशस्वी झाले, तर तीही पॅकेजने सांभाळलेली मानली जाते.
3. नाहीतर ccam तिथे लिहिता येते याची खात्री करण्यासाठी binary शेजारी `.ccam-update-*` ही तात्पुरती फाइल तयार करून काढते.

पॅकेजने सांभाळलेल्या binary साठी ccam थांबते आणि apt चे टप्पे छापते:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux शिवाय इतर systems वर फक्त install script ची ओळ छापली जाते. फोल्डरमध्ये लिहिता येत नसेल, तर ccam install script पुन्हा चालवायला सुचवते, जे ccam `~/.local/bin` मध्ये इन्स्टॉल करते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- development build आपली आवृत्ती `dev` म्हणून सांगते आणि तिच्याशी तुलना करायला रिलीज नसतो, त्यामुळे `ccam update` आणि `ccam update --check` ती नाकारतात. `ccam update --version vX.Y.Z` तिच्यावर रिलीज इन्स्टॉल करते.
- अयशस्वी डाउनलोड, तुमच्या platform साठी asset नसणे, किंवा checksum न जुळणे यांमुळे इन्स्टॉल केलेली binary आहे तशीच राहते.
- अस्तित्वात नसलेला रिलीज `there is no release <tag>` देऊन अयशस्वी होतो. GitHub HTTP 403 किंवा 429 ने नकार देते तेव्हा, ccam सांगते की ही बहुधा rate limit आहे आणि नंतर पुन्हा प्रयत्न करायला सांगते.
- `--check` फक्त वाचते. binary पॅकेजने सांभाळलेली आहे का किंवा तिच्यात लिहिता येते का ते ते तपासत नाही.
- जुन्या tag सह `--version` वापरून तुम्ही आधीच्या आवृत्तीवर परत जाता.

## संबंधित

- [इन्स्टॉल](/docs/install)
- [कमांड्स](/docs/commands)
- [वारंवार विचारले जाणारे प्रश्न](/docs/faq)
