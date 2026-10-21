---
title: ccam आणि इतर Claude Code account switchers
description: Claude Code खाती बदलण्यासाठी ccam ची claude-swap, CCS, clausona, clauth, CAAM, cux आणि इतर टूल्सशी तुलना, आणि कोणती पद्धत कोणत्या वापरासाठी योग्य आहे.
updated: 2026-10-21
order: 2
translation: machine
sourceHash: "c74be64d2978"
---

Claude Code खाती बदलणारी टूल्स दोनपैकी एक पद्धत वापरतात. बहुतेक टूल्स एकाच `~/.claude` फोल्डरमध्ये लॉगिन बदलतात, त्यामुळे एका वेळी एकच खाते सक्रिय असते. इतर टूल्स, ज्यात ccam सुद्धा आहे, `CLAUDE_CONFIG_DIR` ने प्रत्येक खात्याला स्वतःचा फोल्डर देतात, त्यामुळे प्रत्येक खाते लॉगिन राहते आणि अनेक खाती एकाच वेळी चालू शकतात. ccam यात शेअर केलेल्या सेटिंग्ज आणि MCP सर्व्हर जोडते, आणि संभाषण दुसऱ्या खात्याच्या फोल्डरमध्ये पुढे चालू ठेवते.

## एकच लॉगिन बदलायचे की वेगळे फोल्डर?

**लॉगिन बदलताना** सेटिंग्ज, MCP सर्व्हर आणि इतिहासाचा एकच संच असतो, कारण फोल्डर एकच असतो. खाते बदलणे जलद होते, आणि चालू सत्र नवीन लॉगिन घेऊ शकते. पण एका वेळी एकच खाते सक्रिय असते, सर्व टर्मिनल्स आणि editors साठी एकाच वेळी.

लॉगिन बदलण्यासाठी टूल इतर प्रत्येक खात्याच्या लॉगिनची एक प्रत ठेवते आणि गरजेनुसार ती योग्य जागी लिहिते. Claude Code लॉगिनचे नूतनीकरण करताना त्यातील refresh token बदलते, त्यामुळे प्रत्येक जतन केलेली प्रत चालू सत्रांशी जुळवून ठेवावी लागते. खाते बदलल्यानंतर लॉगआउट होणारी किंवा पुन्हा `/login` मागणारी खाती ही या प्रकल्पांत सर्वाधिक वेळा दुरुस्त केल्या जाणाऱ्या समस्यांपैकी आहेत.

**वेगळे फोल्डर** ही [Claude Code च्या दस्तऐवजात सांगितलेली पद्धत आहे](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). प्रत्येक खाते लॉगिन राहते, आणि तुम्ही work आणि personal सत्रे शेजारी शेजारी चालवू शकता. प्रत्येक लॉगिन त्याच्या स्वतःच्या फोल्डरमध्ये राहते, जिथे Claude Code नेहमीप्रमाणे त्याचे नूतनीकरण करते. याची किंमत म्हणजे, एखादे टूल शेअर करत नसेल तर प्रत्येक फोल्डर स्वतःच्या सेटिंग्ज, MCP सर्व्हर आणि इतिहासासह सुरू होतो.

## टूल्सची तुलना

2026-10-20 रोजी प्रत्येक प्रकल्पाच्या स्वतःच्या दस्तऐवजातून आणि source code मधून घेतले (clausona: v0.12.0-beta). "दस्तऐवजात नाही" म्हणजे प्रकल्प त्याबद्दल काही सांगत नाही, वैशिष्ट्य नाही असा अर्थ नाही.

| टूल | पद्धत | एकाच वेळी अनेक खाती | 5 तासांचा आणि साप्ताहिक वापर | मर्यादा गाठल्यावर | प्लॅटफॉर्म |
| --- | --- | --- | --- | --- | --- |
| ccam | वेगळे फोल्डर | हो | हो, `ccam usage` | सत्रात एक सूचना; `ccam next --auto --continue` किंवा `/ccam:switch` संभाषण सर्वात जास्त जागा असलेल्या खात्यात नेते, इतिहास शेअर केलेला असो वा नसो; `ccam --fallback` `-p` jobs तिथे आपोआप पुन्हा चालवते, आणि `ccam pick` workers ना खात्यांमध्ये वाटून देते | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | लॉगिन बदलते; `cswap run` ने वेगळे फोल्डर | हो, `cswap run` ने | हो | `cswap auto` ठरावीक पातळीवर लॉगिन बदलते; चालू सत्र पुढे चालू राहते | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | वेगळे फोल्डर | हो | अंशतः | खात्यांचे pools एका local proxy मधून विनंत्या पुढच्या खात्याकडे पाठवतात | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | वेगळे फोल्डर | हो | हो | Routes प्रत्येक run आधी quota असलेले खाते निवडतात, reset ची वाट पाहू शकतात, आणि workers ना वेगवेगळी खाती देतात | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | लॉगिन बदलते; `clauth start` ने वेगळे फोल्डर | हो, `clauth start` ने | हो | fallback साखळीनुसार पुढच्या खात्यावर जाते | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | लॉगिन बदलते, किंवा isolated profiles | हो, isolated profiles ने | हो | `caam run` non-interactive कमांड दुसऱ्या खात्यावर पुन्हा चालवते | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | लॉगिन बदलते आणि `claude` ला wrap करते | नाही, एकच सक्रिय खाते | हो | लॉगिन बदलते आणि तेच संभाषण पुन्हा सुरू करते | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | लॉगिन बदलते | नाही | हो | ठरावीक पातळीवर ऐच्छिक बदल | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | लॉगिन बदलते | नाही | हो | 100% वर ऐच्छिक बदल | macOS |

ccswap हा claude-swap चा fork आहे, जो Codex लॉगिनही सांभाळतो. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) आणि [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) सुद्धा प्रत्येक खात्याला स्वतःचा फोल्डर देतात. यांपैकी अनेक टूल्स Codex किंवा Gemini सारखे इतर coding agents सुद्धा सांभाळतात.

## ही टूल्स कशी बनवली आहेत

| टूल | कोणत्या भाषेत लिहिले | Claude Code चालू असताना | लॉगिन एका जागेवरून दुसरीकडे कॉपी करते | खाती काय शेअर करू शकतात |
| --- | --- | --- | --- | --- |
| ccam | Go, एक static binary | काहीही नाही: ccam स्वतःची जागा `claude` ला देते (Windows वर, जिथे हे शक्य नाही, ते त्याची वाट पाहते; एक ऐच्छिक relay त्याच्या शेजारी राहतो) | नाही | Groups: प्रत्येक group तो शेअर करत असलेले आयटम आणि त्याची खाती ठरवतो; एक खाते अनेक groups मध्ये असू शकते |
| claude-swap | Python 3.12+, uv किंवा pipx ने इन्स्टॉल | `cswap run` सह काहीही नाही; `cswap auto` आणि menu bar ॲप चालू राहतात | हो, लॉगिन बदलताना | `cswap run` सह: सेटिंग्ज, `CLAUDE.md` आणि extensions एकत्र, इतिहास ऐच्छिक; प्रत्येक launch ला MCP सर्व्हर कॉपी होतात |
| CCS | Node.js 18+ किंवा Bun | CCS parent process म्हणून राहते; pools एक local proxy जोडतात | त्याच्या खाते फोल्डर्ससाठी नाही; pools स्वतःचे लॉगिन ठेवतात | प्रत्येक खात्यासाठी शेअर केलेले किंवा local, इतिहास context groups मध्ये; MCP सर्व्हर कॉपी होतात |
| clausona | TypeScript, Node.js 20+ | साधे `claude` एका shell function मधून चालते; `clausona run` Node ला parent process म्हणून ठेवते; प्रत्येक exit नंतर एक background Node process | नाही, पण idle खात्यांचा वापर वाचण्यासाठी ते स्वतःच त्यांचे tokens नूतनीकरण करते | प्रत्येक खात्याच्या स्वतःच्या फाइल्स सोडून सर्व काही, एका primary मधून; इतिहास ऐच्छिक; `claude mcp add` ने जोडलेले सर्व्हर शेअर होत नाहीत |
| clauth | Rust | clauth parent process म्हणून राहते | हो, लॉगिन बदलताना | संपूर्ण `~/.claude`, किंवा काहीही नाही |
| CAAM | Go | mode वर अवलंबून; एक ऐच्छिक daemon | हो, त्याच्या vault mode मध्ये | mode नुसार ठरलेले; प्रत्येक launch ला MCP सर्व्हर कॉपी होतात |
| cux | Go | `claude` भोवती एक wrapper राहतो | हो | सर्व काही: एकच फोल्डर |
| CCSwitcher, Claude Switcher | macOS ॲप्स (Swift, Python) | menu bar ॲप | हो | सर्व काही: एकच फोल्डर |

## ccam कशात चांगले आहे

### लॉगिन Claude Code ने ठेवले तिथेच राहतात

ccam कधीही Claude लॉगिन एका जागेवरून दुसरीकडे कॉपी करत नाही. प्रत्येक खात्याचे लॉगिन त्याच्या स्वतःच्या फोल्डरमध्ये राहते, किंवा macOS वर त्या फोल्डरच्या स्वतःच्या Keychain item मध्ये, आणि त्याचे नूतनीकरण Claude Code च करते. विसंगत होऊ शकेल अशी कोणतीही जतन केलेली प्रत नसते. आधीचा `~/.claude`, जो `ccam import` ने नोंदवला आहे, साधे `claude` वापरते तोच Keychain item ठेवतो, त्यामुळे दोन्ही प्रकारे तो लॉगिन राहतो.

ccam कालबाह्य token चे नूतनीकरण स्वतः फक्त तुम्ही `ccam usage --refresh` चालवल्यावर करते. ते आधी lock घेते आणि ते खाते चालू असेल तर नकार देते, त्यामुळे चालू सत्राखालून ते कधीही token बदलत नाही.

### तुमच्या आणि Claude Code च्या मध्ये काहीही चालत नाही

ccam एक फाइल वाचते, आपले links तपासते आणि स्वतःची जागा `claude` ला देते (`exec`). तुम्ही काम करत असताना कोणतीही ccam process, wrapper, proxy किंवा daemon नसते: standard input, signals आणि exit code Claude Code चे असतात, आणि त्याच्या विनंत्या थेट Anthropic कडे जातात. Windows process बदलू शकत नाही, त्यामुळे तिथे ccam Claude Code त्याच console सह सुरू करते आणि त्याच्या exit code ची वाट पाहते, आणि दुसरे काहीही करत नाही.

अशा प्रकारे खाते सुरू केल्याने सुमारे 4 ms वाढतात, 95 व्या percentile ला 6 ms पेक्षा कमी, आणि peak memory 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) hyperfine सह, 300 runs). त्याच laptop वर, एखादा प्रोग्राम काही करण्याआधीच फक्त Python interpreter सुरू व्हायला 13 ms आणि Node.js ला 31 ms लागतात.

### काय शेअर होते ते groups ठरवतात

MCP सर्व्हर, `settings.json`, `CLAUDE.md`, extensions (agents, commands, skills, plugins) आणि संभाषणाचा इतिहास हे पाच वेगळे आयटम आहेत. group तो शेअर करत असलेले आयटम आणि त्यातील खाती ठरवतो, त्यामुळे दोन कामाची खाती `work` group मध्ये इतिहास शेअर करू शकतात, तर वैयक्तिक खाते स्वतःचा इतिहास वेगळा ठेवते, आणि तिन्ही खाती `default` मधून तेच MCP सर्व्हर शेअर करतात. `ccam group set work --add memory` group मधील प्रत्येक खाते एकाच योजनेत हलवते; [डॅशबोर्ड](/features/dashboard) हेच checkboxes ने करतो, आणि एखाद्या खात्याला दुसऱ्या group मधून आधीच मिळणारा आयटम राखाडी करून दाखवतो.

- MCP सर्व्हर launch च्या वेळी `--mcp-config` ने दिले जातात, त्यामुळे default नुसार ccam कधीही `.claude.json` संपादित करत नाही; ही फाइल Claude Code चालू असताना पुन्हा लिहिते. `ccam mcp add` संपूर्ण group साठी एकदाच सर्व्हर जोडते.
- बाकीचे आयटम group फोल्डरकडे जाणारे links आहेत, त्यामुळे एका खात्यात केलेला बदल लगेच इतर खात्यांत दिसतो. launch च्या वेळी काहीही कॉपी होत नाही.
- लॉगिन, `.claude.json` आणि चालू स्थिती कधीही शेअर होत नाहीत.

### प्रत्येक बदल आधी दाखवला जातो आणि त्याचा backup घेतला जातो

`ccam group` आणि `ccam migrate` काहीही हलवण्यापूर्वी योजना दाखवतात आणि विचारतात. एखादी फाइल खाते आणि त्याचा group यांच्यात वेगळी असेल, तर तुम्ही निवडता: group ची ठेवायची, तुमची वापरायची, की JSON merge करायचे. मूळ फाइल्स `.ccam-backup/<time>/` मध्ये जातात, काढलेल्या खात्याचा फोल्डर trash फोल्डरमध्ये जातो, आणि फाइल्स एकाच rename मध्ये बदलल्या जातात, त्यामुळे अर्धवट थांबलेले लेखन जुनी फाइल तशीच ठेवते. `--dry-run` दाखवते की launch, sharing, migration, shortcuts, updates आणि दुरुस्त्या काय करतील. `ccam doctor --fix` तुटलेले links आणि permissions दुरुस्त करते.

### पुढच्या खात्यात तेच संभाषण

`ccam next --auto --continue` सर्व खात्यांना एकाच वेळी त्यांचा वापर विचारते, सर्वात जास्त जागा असलेले खाते उघडते आणि संभाषण पुढे चालू ठेवते. इतिहास शेअर केलेला असेल तर ते तेच संभाषण असते; नसेल तर Claude Code ते नवीन खात्यात कॉपी करते (`--resume <file> --fork-session`), आणि पहिले उत्तर किती tokens पुन्हा वाचेल ते ccam सांगते. मर्यादेमुळे सत्र थांबले की ccam plugin चा hook एक सूचना दाखवतो, आणि `/ccam:switch` संभाषण पुढे सोपवते. कधी बदलायचे ते तुम्ही ठरवता; background मध्ये तुमच्या वापरावर काहीही लक्ष ठेवत नाही, आणि आपोआप होणाऱ्या निवडी एकाच organization मध्ये राहतात, त्यामुळे कंपनीचे संभाषण चुकूनही खासगी खात्यावर जात नाही. एखाद्या खात्याने 90% गाठले की status line ही कमांड सुचवते.

Scripts तुमची वाट पाहू शकत नाहीत, त्यामुळे `ccam --fallback work -p "…"` त्यांच्यासाठी खाते बदलते: खात्याने मर्यादा गाठली की तीच कमांड सर्वात जास्त जागा असलेल्या खात्याने पुन्हा चालते, आणि फक्त पूर्ण झालेल्या run चा output छापला जातो. वापर मर्यादेव्यतिरिक्त इतर कोणतेही अपयश जसेच्या तसे पुढे दिले जाते. समांतर workers साठी `ccam pick` प्रत्येकाला जागा असलेले वेगळे खाते देते, आणि कोणत्याही खात्यात जागा उरली नाही तर दोन्ही code 75 ने बाहेर पडतात.

### प्रत्येक फोल्डरमध्ये योग्य खाते

`ccam folder work ~/code/company` मुळे त्या फोल्डरमध्ये खात्याच्या नावाशिवाय होणारा प्रत्येक launch `work` वापरतो, तर तुमचे बाकीचे प्रकल्प default खाते वापरत राहतात. ही सेटिंग प्रत्येक launch वाचतो त्याच `config.toml` मध्ये असते, त्यामुळे यात वेळ जात नाही.

### एक binary, तपासलेले downloads

ccam सुमारे 10 MB ची static binary आहे, ज्यासाठी कोणताही runtime इन्स्टॉल करावा लागत नाही, Linux, macOS आणि Windows साठी. install scripts (`install.sh`, `install.ps1`) आणि `ccam update` प्रत्येक download release च्या SHA-256 checksums शी जुळवून तपासतात, आणि ccam तुम्हाला विचारल्याशिवाय कधीही update होत नाही. त्यात telemetry नाही: ते फक्त तुमच्या खात्यांचा वापर जाणून घेण्यासाठी Anthropic कडे विनंत्या पाठवते, समांतरपणे, 3 सेकंदांच्या मर्यादेसह आणि 60 सेकंदांच्या cache सह, डॅशबोर्ड उघडताना नवीनतम release बद्दल GitHub ला विचारते, आणि तुम्ही मागितलेले downloads करते. `usage_source = "claude"` सह त्या विनंत्याही Claude Code मधूनच जातात, आणि ccam कधीही token हाताळत नाही.

## इतर टूल्स जिथे जास्त करतात

- claude-swap, clauth, cux, CAAM आणि CCSwitcher एखाद्या खात्याचा वापर ठरावीक पातळी ओलांडला की चालू interactive सत्र आपोआप दुसऱ्या खात्यावर नेऊ शकतात. ccam `-p` runs `--fallback` ने बदलते आणि interactive सत्रे तुमच्यावर सोडते, मर्यादेवर एक सूचना आणि `/ccam:switch` सह.
- Developer Mode शिवाय Windows वर ccam फोल्डर शेअर करते, पण `settings.json` सारख्या एकट्या फाइल्स शेअर करू शकत नाही; claude-swap आणि CCS तिथे त्या कॉपी करतात.
- CCSwitcher, Claude Switcher आणि claude-swap वापर macOS menu bar मध्ये दाखवतात.
- CCS, clausona, CAAM, clauth, ccswap आणि Claude Switcher, Codex आणि इतर coding agents सुद्धा सांभाळतात.
- clausona प्रत्येक खात्यासाठी extensions चालू आणि बंदही करते (ccam चे `ccam ext` ते फक्त दाखवते), local token खर्चाचा हिशेब ठेवते, आणि tmux व इतर terminals मध्ये workers चे संपूर्ण ताफे चालवणारे Claude Code plugins देते.
- CAAM, CCSwitcher आणि Claude Switcher, Homebrew ने इन्स्टॉल होतात, आणि clauth आणि CAAM त्यांच्या release checksums वर स्वाक्षरी करतात. ccam त्याच्या script ने, `.deb` ने किंवा `go install` ने इन्स्टॉल होते, आणि स्वाक्षरीशिवाय checksums तपासते.

## कोणते योग्य ठरेल?

तुम्ही एका वेळी एकच खाते वापरत असाल किंवा चालू सत्रे आपोआप बदलावीत असे हवे असेल, तर लॉगिन बदलणारे टूल निवडा. Linux, macOS किंवा Windows वर खाती शेजारी शेजारी, एकच शेअर केलेला सेटअप, पुढच्या खात्यात तेच संभाषण आणि Claude Code सोबत इतर काहीही चालू नसावे असे हवे असेल, तर ccam निवडा. काहीही शेअर करायचे नसेल, तर [साधे शेल aliases](/compare/aliases) पुरेसे आहेत.

अनेक खाती वापरण्याची परवानगी आहे की नाही, हे प्रत्येक खात्याच्या अटींवर अवलंबून आहे. वापरण्यापूर्वी त्या तपासा.
