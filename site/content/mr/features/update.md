---
title: ccam अपडेट करणे
description: ccam update SHA-256 तपासणीनंतर ccam च्या जागी GitHub वरील सर्वात नवीन रिलीज ठेवते; --check फक्त कळवते आणि --version दिलेला रिलीज tag इन्स्टॉल करते.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: "598931e4801f"
---

`ccam update` GitHub वरून ccam रिलीज डाउनलोड करते, तो रिलीजच्या `checksums.txt` शी तपासते आणि चालू binary एका rename ने बदलते. तुम्ही ccam install script ने किंवा archive मधून इन्स्टॉल केले असेल तेव्हा ते वापरा. [डॅशबोर्ड](/features/dashboard) देखील प्रत्येक वेळी उघडताना GitHub ला नवीनतम रिलीज विचारतो आणि `U` ने तो देऊ करतो, जे विचारल्यानंतर हीच पावले चालवते. इतर कमांड्स GitHub शी कधीच संपर्क करत नाहीत, त्यामुळे त्या offline सुद्धा तशाच चालतात.

## वापर

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
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

नवीन रिलीज नसेल तेव्हा ccam `ccam 0.4.0 is up to date.` छापते आणि काहीही बदलत नाही. `--version` तुमच्याकडे आधीच असलेल्या आवृत्तीचे नाव देत असेल, तर ते `ccam 0.4.0 is already installed.` छापते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--check` | फक्त नवीन रिलीज आहे का ते कळवा | बंद |
| `--version <tag>` | `v0.3.0` सारखा हा रिलीज tag इन्स्टॉल करा, सध्याच्यापेक्षा नवीन किंवा जुना | latest |
| `--dry-run` | समर्थित: सर्वात नवीन tag शोधतो (`--version` ने एखादा दिला नसेल तर), binary वर त्याच तपासण्या चालवतो, आणि डाउनलोड न करता download, checksum आणि replace पाथ छापतो | बंद |

`--check` आणि `--version` एकत्र वापरता येत नाहीत. tag `vMAJOR.MINOR.PATCH` सारखा असणे आवश्यक आहे, ऐच्छिक pre-release किंवा build suffix सह; सुरुवातीचा `v` वगळता येतो. इतर काहीही अयशस्वी होते, उदाहरणार्थ `--version takes a release tag such as v0.3.0, not "1.2"`.

`CCAM_RELEASES_URL` हा environment variable mirror साठी किंवा tests साठी `https://github.com/jhyunwoo/ccam/releases` ची जागा घेतो. त्याची रचना तशीच असणे आवश्यक आहे: `<url>/latest` `<url>/tag/<tag>` कडे redirect करतो, आणि assets `<url>/download/<tag>/<asset>` वर असतात.

## हे कसे काम करते

इन्स्टॉल केलेल्या binary ला हात लावण्यापूर्वी ccam download तपासते, आणि binary एका rename ने बदलते. कोड `internal/update` मध्ये आहे, जो फक्त `ccam update` import करते. सुरू करण्याचा मार्ग ते कधीच बोलावत नाही.

1. ccam `os.Executable` ने चालू binary शोधते आणि `filepath.EvalSymlinks` ने symlinks सोडवते (`cli.updateCmd`). हाच खरा पाथ बदलला जातो.
2. `--version` शिवाय, ccam redirect न पाळता `<releases>/latest` ची request करते (`Updater.Latest`), आणि `Location` header (`…/tag/v0.5.0`) मधून tag वाचते. हे GitHub API ऐवजी releases page वापरते, त्यामुळे API ची तासाला 60 requests ची मर्यादा लागू होत नाही. 10 सेकंदांनंतर हा शोध timeout होतो.
3. ccam tag एका कडक pattern शी तपासते (`parseVersion`). tag download URLs मध्ये आणि छापलेल्या shell कमांडमध्ये जातो, त्यामुळे `/`, `%`, spaces किंवा shell अक्षरे असलेले काहीही पुढे जात नाही. मग ccam त्याची चालू आवृत्तीशी तुलना करते (`compareVersions`); pre-release त्याच्या रिलीजच्या आधी येतो.
4. ccam ठरवते की ते binary बदलू शकते की नाही (`managed()`, `writable()`; खाली पाहा).
5. ccam `ccam_<os>_<arch>.tar.gz` (Windows वर `ccam_windows_<arch>.zip`) आणि `checksums.txt` एकाच वेळी डाउनलोड करते, दोन्हीसाठी 2 मिनिटे (`Updater.Fetch`). प्रत्येक डाउनलोड 256 MB पर्यंत मर्यादित असतो.
6. इतर काहीही होण्यापूर्वी archive चा SHA-256 `checksums.txt` मधील त्याच्या ओळीशी जुळणे आवश्यक आहे. ओळ नसल्यास किंवा checksum वेगळा असल्यास अपडेट `nothing was changed` सह थांबतो.
7. ccam tar वाचते आणि `ccam` नावाची पहिली regular फाइल घेते (`extract()`; `extractZip` zip मधून `ccam.exe` घेते). फोल्डर, लिंक्स आणि इतर नोंदी दुर्लक्षित केल्या जातात, आणि binary सुद्धा 256 MB पर्यंत मर्यादित असते.
8. ccam नवीन binary `fsutil.AtomicWrite` ने लिहिते (`Updater.Install`): binary शेजारी `.ccam.tmp*` ही तात्पुरती फाइल, जुन्या फाइलच्या permission bits सह, डिस्कवर flush करून जुन्या फाइलवर rename केलेली.

tag शोधण्यापासून rename पर्यंत संपूर्ण pipeline:

<!-- motion:update-pipeline -->

archive आणि `checksums.txt` समांतरपणे डाउनलोड होतात, आणि checksum जुळेपर्यंत काहीही लिहिले जात नाही. rename atomic आहे, त्यामुळे crash किंवा भरलेली डिस्क एकतर जुनी binary किंवा नवीन binary ठेवते, कधीच अर्धवट नाही. `~/.local/bin/ccwork` सारख्या शॉर्टकट कमांड्स ccam binary च्या symlinks आहेत, त्यामुळे त्या कोणत्याही बदलाशिवाय नवीन आवृत्ती चालवतात.

Windows वर चालू प्रोग्रामची फाइल बदलता किंवा हटवता येत नाही, फक्त तिचे नाव बदलता येते. ccam नवीन binary `ccam.exe.new` मध्ये लिहिते, चालू `ccam.exe` चे नाव `ccam.exe.old` करते, आणि नवीन फाइलचे नाव `ccam.exe` करते (`Updater.Install`). शेवटचा टप्पा अयशस्वी झाला, तर जुन्या binary चे नाव परत पूर्वीसारखे केले जाते. पुढचा अपडेट हटवता येतील त्या `ccam.exe.old` फाइल्स हटवतो; आधीच्या `ccam.exe.old` वरून एखादे सत्र अजून चालू असेल, तर चालू binary त्याऐवजी `ccam.exe.old-<time>` बनते. symbolic links ऐवजी hard links असलेले शॉर्टकट ([शॉर्टकट कमांड्स](/features/shortcuts)) लगेच नंतर नवे केले जातात, त्यामुळे तेही नवीन आवृत्ती चालवतात; यात फक्त नुकत्याच बदललेल्या binary च्याच नव्हे, तर ccam ने तयार केलेल्या, ccam च्या कोणत्याही जुन्या आवृत्तीच्या hard links चा समावेश होतो.

### पॅकेजने इन्स्टॉल केलेले

`.deb` मधून इन्स्टॉल केलेली ccam ची प्रत apt ची असते, आणि apt च्या नकळत ती बदलल्यास package database फाइलशी जुळणार नाही. डाउनलोड करण्यापूर्वी ccam binary कोणाची आहे ते तपासते:

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
  curl -fsSL https://ccam.dev/install.sh | sh
```

Linux शिवाय इतर systems वर फक्त install script ची ओळ छापली जाते: Windows वर `irm https://ccam.dev/install.ps1 | iex`. फोल्डरमध्ये लिहिता येत नसेल, तर ccam install script पुन्हा चालवायला सुचवते, जे ccam `~/.local/bin` मध्ये इन्स्टॉल करते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- development build आपली आवृत्ती `dev` म्हणून सांगते आणि तिच्याशी तुलना करायला रिलीज नसतो, त्यामुळे `ccam update` आणि `ccam update --check` ती नाकारतात. `ccam update --version vX.Y.Z` तिच्यावर रिलीज इन्स्टॉल करते. `go install` ने बनवलेल्या builds साठी [ccam इन्स्टॉल करा](/docs/install#other-ways-to-install) पाहा.
- अयशस्वी डाउनलोड, तुमच्या platform साठी asset नसणे, किंवा checksum न जुळणे यांमुळे इन्स्टॉल केलेली binary आहे तशीच राहते.
- अस्तित्वात नसलेला रिलीज `there is no release <tag>` देऊन अयशस्वी होतो. GitHub HTTP 403 किंवा 429 ने नकार देते तेव्हा, ccam सांगते की ही बहुधा rate limit आहे आणि नंतर पुन्हा प्रयत्न करायला सांगते.
- `--check` फक्त वाचते. binary पॅकेजने सांभाळलेली आहे का किंवा तिच्यात लिहिता येते का ते ते तपासत नाही.
- जुन्या tag सह `--version` वापरून तुम्ही आधीच्या आवृत्तीवर परत जाता.

## संबंधित

- [ccam इन्स्टॉल करा](/docs/install)
- [कमांड संदर्भ](/docs/commands)
- [वारंवार विचारले जाणारे प्रश्न](/docs/faq)
