---
title: ccam अपडेट करना
description: ccam update, SHA-256 जाँच के बाद ccam को GitHub की सबसे नई रिलीज़ से बदलता है; --check केवल बताता है और --version कोई दिया गया रिलीज़ टैग इंस्टॉल करता है।
updated: 2026-10-21
order: 10
translation: machine
sourceHash: a92dbf37399f
---

`ccam update`, GitHub से ccam की एक रिलीज़ डाउनलोड करता है, उसे रिलीज़ की `checksums.txt` से जाँचता है और चल रही बाइनरी को एक rename में बदल देता है। इसे तब इस्तेमाल करें जब आपने ccam को install स्क्रिप्ट से या किसी archive से इंस्टॉल किया हो। [डैशबोर्ड](/features/dashboard) भी हर बार खुलने पर GitHub से नवीनतम रिलीज़ पूछता है और उसे `U` से पेश करता है, जो पूछने के बाद यही कदम चलाता है। कोई और कमांड GitHub से संपर्क नहीं करती, इसलिए वे offline भी वैसे ही काम करती हैं।

## उपयोग

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

जब कोई नई रिलीज़ नहीं होती, तो ccam `ccam 0.4.0 is up to date.` प्रिंट करता है और कुछ नहीं बदलता। `--version` के साथ वही वर्ज़न देने पर जो आप पहले से चला रहे हैं, यह `ccam 0.4.0 is already installed.` प्रिंट करता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--check` | केवल बताएँ कि कोई नई रिलीज़ है या नहीं | बंद |
| `--version <tag>` | यह रिलीज़ टैग इंस्टॉल करें, जैसे `v0.3.0`, मौजूदा से नया या पुराना | सबसे नई |
| `--dry-run` | माना जाता है: सबसे नया टैग खोजता है (जब तक `--version` कोई टैग न दे), बाइनरी पर वही जाँचें चलाता है, और डाउनलोड किए बिना डाउनलोड, checksum और बदलने वाले पाथ प्रिंट करता है | बंद |

`--check` और `--version` एक साथ इस्तेमाल नहीं हो सकते। टैग `vMAJOR.MINOR.PATCH` जैसा दिखना चाहिए, वैकल्पिक pre-release या build suffix के साथ; शुरुआत का `v` छोड़ा जा सकता है। कुछ और देने पर विफल होता है, उदाहरण के लिए `--version takes a release tag such as v0.3.0, not "1.2"`।

Environment variable `CCAM_RELEASES_URL`, mirror या tests के लिए `https://github.com/jhyunwoo/ccam/releases` की जगह लेता है। इसकी बनावट वैसी ही होनी चाहिए: `<url>/latest`, `<url>/tag/<tag>` पर redirect करता है, और assets `<url>/download/<tag>/<asset>` पर होते हैं।

## यह कैसे काम करता है

कोड `internal/update` में है, जिसे केवल `ccam update` import करता है। शुरू करने का रास्ता इसे कभी नहीं बुलाता।

1. `cli.updateCmd`, `os.Executable` से चल रही बाइनरी खोजता है और `filepath.EvalSymlinks` से symlinks हल करता है। वही असली पाथ बदला जाता है।
2. `--version` के बिना, `Updater.Latest` redirect का पालन किए बिना `<releases>/latest` पर request भेजता है, और `Location` header (`…/tag/v0.5.0`) से टैग पढ़ता है। यह GitHub API की जगह releases पेज इस्तेमाल करता है, इसलिए API की एक घंटे में 60 requests की सीमा लागू नहीं होती। इस खोज के लिए 10 सेकंड हैं।
3. `parseVersion` टैग को एक सख़्त pattern से जाँचता है। टैग डाउनलोड URLs में और प्रिंट की गई shell कमांड में जाता है, इसलिए `/`, `%`, spaces या shell अक्षरों वाली कोई चीज़ आगे नहीं जाती। `compareVersions` इसकी तुलना चल रहे वर्ज़न से करता है; pre-release अपनी रिलीज़ से पहले गिना जाता है।
4. `managed()` और `writable()` तय करते हैं कि ccam बाइनरी को बदल भी सकता है या नहीं (नीचे देखें)।
5. `Updater.Fetch`, `ccam_<os>_<arch>.tar.gz` (Windows पर `ccam_windows_<arch>.zip`) और `checksums.txt` एक साथ डाउनलोड करता है, दोनों के लिए 2 मिनट के साथ। हर डाउनलोड की सीमा 256 MB है।
6. कुछ और होने से पहले archive का SHA-256, `checksums.txt` में उसकी पंक्ति से मेल खाना चाहिए। पंक्ति न होने या checksum अलग होने पर अपडेट `nothing was changed` के साथ रुक जाता है।
7. `extract()` tar पढ़ता है और `ccam` नाम की पहली सामान्य फ़ाइल लेता है (`extractZip` zip से `ccam.exe` लेता है)। Directories, links और दूसरी entries अनदेखी की जाती हैं, और बाइनरी की सीमा भी 256 MB है।
8. `Updater.Install` नई बाइनरी `fsutil.AtomicWrite` से लिखता है: बाइनरी के बगल में एक अस्थायी फ़ाइल `.ccam.tmp*`, पुरानी फ़ाइल के permission bits के साथ, जो डिस्क पर flush होकर पुरानी फ़ाइल के ऊपर rename होती है।

टैग की खोज से rename तक पूरी pipeline:

<!-- motion:update-pipeline -->

Archive और `checksums.txt` साथ-साथ डाउनलोड होते हैं, और checksum मेल खाने तक कुछ नहीं लिखा जाता। Rename atomic है, इसलिए crash या भरी हुई डिस्क के बाद या तो पुरानी बाइनरी बचती है या नई, कभी आधी नहीं। `~/.local/bin/ccwork` जैसी शॉर्टकट कमांड ccam बाइनरी के symlinks हैं, इसलिए वे बिना किसी बदलाव के नया वर्ज़न चलाती हैं।

Windows पर चल रहे प्रोग्राम की फ़ाइल को बदला या मिटाया नहीं जा सकता, केवल उसका नाम बदला जा सकता है। `Updater.Install` नई बाइनरी `ccam.exe.new` में लिखता है, चल रहे `ccam.exe` का नाम बदलकर `ccam.exe.old` करता है, और नई फ़ाइल का नाम बदलकर `ccam.exe` करता है। अगर आख़िरी चरण विफल होता है, तो पुरानी बाइनरी का नाम वापस बदल दिया जाता है। अगला अपडेट `ccam.exe.old` मिटा देता है। जो शॉर्टकट symbolic links के बजाय hard links हैं ([शॉर्टकट कमांड](/features/shortcuts)), वे तुरंत बाद नए कर दिए जाते हैं, ताकि वे भी नया वर्ज़न चलाएँ।

### Package manager से हुए इंस्टॉलेशन

`.deb` से इंस्टॉल किया गया ccam apt का है, और apt की जानकारी के बिना उसे बदलने से package database ग़लत हो जाता। डाउनलोड करने से पहले ccam जाँचता है कि बाइनरी किसकी है:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` या `/sbin` में मौजूद बाइनरी package-managed मानी जाती है।
2. दूसरी जगहों पर, अगर `dpkg` इंस्टॉल है और `dpkg -S <path>` सफल होता है, तो वह भी package-managed मानी जाती है।
3. वरना ccam बाइनरी के बगल में एक अस्थायी `.ccam-update-*` फ़ाइल बनाकर हटाता है, ताकि पक्का हो कि वह वहाँ लिख सकता है।

Package-managed बाइनरी के लिए ccam रुक जाता है और apt के चरण प्रिंट करता है:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux के अलावा दूसरे सिस्टम पर केवल install स्क्रिप्ट वाली पंक्ति प्रिंट होती है: Windows पर `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`। जब फ़ोल्डर में लिखा नहीं जा सकता, तो ccam install स्क्रिप्ट फिर से चलाने का सुझाव देता है, जो ccam को `~/.local/bin` में इंस्टॉल करती है।

## सुरक्षा और विशेष स्थितियाँ

- Development build अपना वर्ज़न `dev` बताता है और उसके पास तुलना के लिए कोई रिलीज़ नहीं होती, इसलिए `ccam update` और `ccam update --check` इसे अस्वीकार करते हैं। `ccam update --version vX.Y.Z` उसके ऊपर एक रिलीज़ इंस्टॉल करता है।
- विफल डाउनलोड, आपके प्लेटफ़ॉर्म का asset न होना, या checksum का मेल न खाना इंस्टॉल की गई बाइनरी को वैसा ही छोड़ देता है।
- जो रिलीज़ मौजूद नहीं है, वह `there is no release <tag>` के साथ विफल होती है। जब GitHub HTTP 403 या 429 के साथ मना करता है, तो ccam बताता है कि यह शायद rate limit है और बाद में फिर कोशिश करने को कहता है।
- `--check` केवल पढ़ता है। यह नहीं जाँचता कि बाइनरी package-managed है या उसमें लिखा जा सकता है।
- पुराने टैग के साथ `--version` ही पिछले वर्ज़न पर लौटने का तरीक़ा है।

## संबंधित

- [इंस्टॉल](/docs/install)
- [कमांड](/docs/commands)
- [अक्सर पूछे जाने वाले सवाल](/docs/faq)
