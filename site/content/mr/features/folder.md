---
title: फोल्डरनुसार खाती
description: ccam folder मुळे खात्याचे नाव न देता Claude Code सुरू केल्यावर फोल्डर आणि त्याखालील सर्व काही एकच खाते वापरते, जसे तुमच्या कंपनीच्या code साठी work.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: "46d1410a3d81"
---

`ccam folder work ~/code/company` मुळे `~/code/company` मध्ये, आणि त्याखालील कोणत्याही फोल्डरमध्ये, खात्याच्या नावाशिवाय होणारा प्रत्येक launch Claude Code `work` म्हणून सुरू करतो. त्या फोल्डरबाहेर, डीफॉल्ट खाते पूर्वीसारखेच सुरू होते. `ccam personal` ने किंवा `ccpersonal` सारख्या शॉर्टकटने खात्याचे नाव दिल्यास नेहमी तेच खाते सुरू होते.

## वापर

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

हे launches फोल्डरचे खाते वापरतात:

- `ccam -p "…"` आणि ज्याचा पहिला आर्ग्युमेंट Claude Code साठीचा पर्याय आहे असा इतर कोणताही launch
- नावाशिवाय `ccam run`, आणि `ccam --dry-run`
- `ccam --fallback -p "…"` त्याच्या पहिल्या प्रयत्नासाठी
- डॅशबोर्ड, म्हणजे एकटे `ccam`, जो फोल्डरच्या खात्यावर cursor ठेवून सुरू होतो

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `<account>` | वापरायच्या खात्याचे नाव किंवा क्रमांक | |
| `[folder]` | फोल्डर, relative किंवा absolute | सध्याचा फोल्डर |
| `--remove` | नेमक्या त्याच फोल्डरची नोंद काढा | बंद |
| `--dry-run` | पाळला जातो: बदल सेव्ह न करता छापा | बंद |

## हे कसे काम करते

प्रत्येक फोल्डर `~/.config/ccam/config.toml` मधील एक `[[folder]]` नोंद आहे, जी फाइल प्रत्येक launch तसाही वाचतो:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` फोल्डरला symbolic links सोडवलेल्या absolute पाथमध्ये बदलते आणि तो सेव्ह करते. एखादा फोल्डर पुन्हा सेट केल्यास त्याचे खाते बदलते.
2. नावाशिवायचा launch `launcher.HereProfile` बोलावतो. एकही नोंद नसल्यास ते लगेच डीफॉल्ट खाते घेते.
3. नाहीतर ते working directory ची प्रत्येक `path` शी strings म्हणून तुलना करते: तो फोल्डर स्वतः किंवा त्याखालचा फोल्डर, पण फक्त त्याच अक्षरांनी सुरू होणारा फोल्डर कधीच नाही (`/code/company` मध्ये `/code/companyx` येत नाही). सर्वात खोलची नोंद जिंकते, त्यामुळे `~/code` `personal` असू शकते तर `~/code/company` `work`.
4. working directory जशी आहे तशी कशाशीच जुळत नसेल, तरच ccam तिच्या symbolic links सोडवते आणि पुन्हा एकदा तुलना करते.
5. Windows वर तुलना Windows प्रमाणेच लहान आणि मोठ्या अक्षरांकडे दुर्लक्ष करते.

`config.toml` शिवाय दुसरी कोणतीही फाइल वाचली जात नाही, त्यामुळे फोल्डरमधील launch ला मोजता येईल इतका वेळ लागत नाही (नोंदींसह किंवा नोंदींशिवाय 4.0 ms, `hyperfine`, 300 runs).

## सुरक्षितता आणि अपवादात्मक स्थिती

- `ccam rename` फोल्डर नोंदी नवीन नावावर नेते, आणि `ccam remove` काढलेल्या खात्याच्या नोंदी हटवते.
- ज्या नोंदीचे खाते आता अस्तित्वात नाही, ती त्या फोल्डरमधील नावाशिवायचा launch थांबवते: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` अशा नोंदी, आणि ज्यांचा फोल्डर नाहीसा झाला आहे त्या नोंदी कळवते.
- वरच्या फोल्डरमधून खाते मिळणाऱ्या फोल्डरमध्ये `ccam folder --remove` काहीही काढण्याऐवजी त्या वरच्या फोल्डरचे नाव सांगते.
- `ccam next` rotation किंवा वापरावरून निवडते आणि फोल्डर्सकडे पाहत नाही.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [शॉर्टकट कमांड्स](/features/shortcuts)
- [खाती व्यवस्थापित करणे](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
