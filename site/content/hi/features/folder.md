---
title: फ़ोल्डर के हिसाब से अकाउंट
description: ccam folder से कोई फ़ोल्डर और उसके नीचे का सब कुछ एक अकाउंट इस्तेमाल करता है जब आप बिना अकाउंट नाम के Claude Code शुरू करें, जैसे कंपनी के code के लिए work।
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 46d1410a3d81
---

`ccam folder work ~/code/company` से `~/code/company` के अंदर, और उसके नीचे के किसी भी फ़ोल्डर में, बिना अकाउंट नाम वाला हर launch Claude Code को `work` के रूप में शुरू करता है। उसके बाहर, डिफ़ॉल्ट अकाउंट पहले की तरह शुरू होता है। अकाउंट का नाम देने पर, `ccam personal` से या `ccpersonal` जैसे शॉर्टकट से, हमेशा वही अकाउंट शुरू होता है।

## उपयोग

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

ये launches फ़ोल्डर का अकाउंट इस्तेमाल करते हैं:

- `ccam -p "…"` और कोई भी दूसरा launch जिसका पहला आर्ग्युमेंट Claude Code का विकल्प हो
- बिना नाम के `ccam run`, और `ccam --dry-run`
- `ccam --fallback -p "…"` अपने पहले प्रयास के लिए
- डैशबोर्ड, यानी अकेला `ccam`, जो फ़ोल्डर के अकाउंट पर cursor के साथ खुलता है

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `<account>` | इस्तेमाल किए जाने वाले अकाउंट का नाम या नंबर | |
| `[folder]` | फ़ोल्डर, relative या absolute | मौजूदा फ़ोल्डर |
| `--remove` | ठीक उसी फ़ोल्डर की entry हटाएँ | बंद |
| `--dry-run` | माना जाता है: बदलाव को सेव किए बिना प्रिंट करें | बंद |

## यह कैसे काम करता है

हर फ़ोल्डर `~/.config/ccam/config.toml` में एक `[[folder]]` entry है, वही फ़ाइल जिसे हर launch वैसे भी पढ़ता है:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` फ़ोल्डर को symbolic links सुलझाकर एक absolute पाथ में बदलता है और उसे सेव करता है। किसी फ़ोल्डर को फिर से सेट करने पर उसका अकाउंट बदल जाता है।
2. बिना नाम वाला launch `launcher.HereProfile` को बुलाता है। कोई entry न होने पर यह तुरंत डिफ़ॉल्ट अकाउंट लेता है।
3. वरना यह working directory की तुलना हर `path` से strings के रूप में करता है: ख़ुद वह फ़ोल्डर या उसके नीचे का कोई फ़ोल्डर, कभी ऐसा फ़ोल्डर नहीं जो सिर्फ़ उन्हीं अक्षरों से शुरू होता हो (`/code/company`, `/code/companyx` को कवर नहीं करता)। सबसे गहरी entry जीतती है, इसलिए `~/code` `personal` हो सकता है जबकि `~/code/company` `work` है।
4. केवल तभी जब working directory जैसी है वैसी किसी से मेल नहीं खाती, ccam उसके symbolic links सुलझाता है और एक बार फिर तुलना करता है।
5. Windows पर तुलना बड़े और छोटे अक्षरों को अनदेखा करती है, जैसे Windows करता है।

`config.toml` के अलावा कोई फ़ाइल नहीं पढ़ी जाती, इसलिए किसी फ़ोल्डर के अंदर launch में मापने लायक कोई समय नहीं लगता (entries के साथ या बिना 4.0 ms, `hyperfine`, 300 runs)।

## सुरक्षा और विशेष स्थितियाँ

- `ccam rename` फ़ोल्डर entries को नए नाम पर ले जाता है, और `ccam remove` हटाए गए अकाउंट की entries मिटा देता है।
- जिस entry का अकाउंट अब मौजूद नहीं है, वह उस फ़ोल्डर में बिना नाम वाले launch को रोक देती है: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`। `ccam doctor` ऐसी entries की, और उन entries की जिनका फ़ोल्डर अब नहीं है, रिपोर्ट करता है।
- जो फ़ोल्डर अपना अकाउंट ऊपर के किसी फ़ोल्डर से पाता है, उसके अंदर `ccam folder --remove` कुछ हटाने के बजाय उस ऊपर वाले फ़ोल्डर का नाम बताता है।
- `ccam next` रोटेशन या उपयोग के हिसाब से चुनता है और फ़ोल्डरों को नहीं देखता।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [शॉर्टकट कमांड](/features/shortcuts)
- [अकाउंट प्रबंधित करना](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
