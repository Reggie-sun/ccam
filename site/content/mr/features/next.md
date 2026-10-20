---
title: पुढच्या खात्यावर जाणे
description: ccam next शेवटच्या खात्यानंतरच्या खात्याने Claude Code उघडते; --auto सर्वाधिक जागेचे खाते निवडते, --continue चालू ठेवते, --fallback -p jobs पुन्हा चालवते.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` तुम्ही नाव न देता दुसऱ्या खात्याने Claude Code उघडते. एखादे खाते वापर मर्यादेला पोहोचल्यावर ते वापरा: साधे `ccam next` तुम्ही शेवटी वापरलेल्या खात्यानंतरचे खाते घेते, `--auto` 5 तासांच्या window मध्ये सर्वात जास्त जागा असलेले खाते घेते, आणि `--continue` शेवटचे संभाषण नवीन खात्यात पुढे चालू ठेवते. `claude -p` चालवणाऱ्या scripts साठी `ccam --fallback` हा बदल आपोआप करते ([मर्यादेवर पुन्हा चालवणे](#run-again-at-the-limit) पाहा).

## वापर

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam हा बदल standard error वर `→ <from> → <to>` असा छापते, मग `ccam <name>` प्रमाणेच नवीन खात्याने Claude Code सुरू करते. `<from>` म्हणजे तुम्ही शेवटी वापरलेले खाते. `--` नंतरचे आर्ग्युमेंट्स न बदलता `claude` ला जातात.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--auto` | वापरावरून लॉग इन केलेले खाते निवडा (खाली पहा), 100% वरची खाती वगळून | बंद |
| `--max-usage <1-100>` | `--auto` सह: या वापराखालील खात्यांना प्राधान्य द्या | `auto_max_usage`, नाहीतर 100 |
| `--strategy headroom\|expiring` | `--auto` सह: सर्वात कमी 5 तासांचा वापर, किंवा सर्वात लवकर होणारा 5 तासांचा reset | `auto_strategy`, नाहीतर `headroom` |
| `--wait[=<duration>]` | `--auto` सह: कोणत्याही खात्यात जागा नसेल तर पहिले खाते reset होण्याची वाट पाहा, जास्तीत जास्त इतका वेळ | बंद; एकटे दिल्यास `5h` |
| `-c`, `--continue` | फक्त तुम्ही शेवटी वापरलेल्या खात्याशी संभाषण इतिहास शेअर करणारी खाती विचारात घ्या, आणि `claude --continue` सुरू करा | बंद |
| `-- <claude args>` | ccam च्या स्वतःच्या आर्ग्युमेंट्सनंतर `claude` ला दिले जाणारे आर्ग्युमेंट्स | काहीही नाही |
| `--dry-run` | पाळला जातो: Claude Code सुरू न करता किंवा सुरू केल्याची नोंद न करता `→ <from> → <to>` आणि चालवली जाणारी कमांड लाइन छापतो | बंद |

`--auto` आणि `--continue` एकत्र चालतात. `--continue` सह, `--continue` तुमच्या स्वतःच्या आर्ग्युमेंट्सच्या आधी ठेवला जातो.

## हे कसे काम करते

ccam मार्फत प्रत्येक वेळी सुरू करताना वेळ `~/.local/state/ccam/last-used/<name>` मध्ये लिहिली जाते. `ccam next` या फाइल्स वाचते, कोणती खाती लॉग इन आहेत ते तपासते, आणि निवड करते.

1. `cli.candidates` `config.toml` वरून प्रत्येक खात्यासाठी एक `next.Candidate` तयार करते: खाते, ते शेवटी सुरू केल्याची वेळ (`launcher.LastUsed`) आणि त्याला लॉगिन आहे का (`creds.Has`: वाचता येणारी `.credentials.json`, किंवा macOS वर Keychain item).
2. `--auto` शिवाय, `next.Rotate` खाती slot क्रमांकानुसार लावते आणि सर्वात अलीकडे सुरू केलेले खाते शोधते. ज्या खात्यावरून तुम्ही दुसरीकडे जात आहात ते हेच. अजून एकही खाते सुरू केले नसेल, तर ते सर्वात लहान slot असते.
3. तिथून ते slot क्रमाने पुढे जाते, शेवटच्या slot नंतर पुन्हा पहिल्यावर येते, आणि लॉग इन असलेले पहिले खाते परत देते.
4. `launcher.LaunchWith` निवडलेले खाते अगदी `ccam <name>` प्रमाणे तयार करून सुरू करते: ते खात्याच्या लिंक्स तपासते, MCP शेअरिंग लागू करते, session नोंदवते आणि `syscall.Exec` मार्फत ccam च्या जागी `claude` ठेवते.

शोध शेवटी वापरलेल्या खात्यापासून सुरू होतो आणि त्यानंतरच्या पहिल्या लॉग इन असलेल्या खात्यावर थांबतो:

<!-- motion:next-rotate -->

`work` सर्वात अलीकडे सुरू केले होते, म्हणून शोध त्याच्यानंतर सुरू होतो. `side` ला लॉगिन नाही आणि ते वगळले जाते, आणि `lab` निवडले जाते. `lab` सुद्धा लॉग आउट असते, तर शोध फिरून `personal` वर आला असता (आकृतीतील `% len`). rotation कधीच तुम्ही जिथून सुरू केले ते खाते परत देत नाही: दुसरे लॉग इन असलेले खाते नसल्यास ते त्रुटीसह थांबते.

### --auto

`--auto` slot क्रमाऐवजी वापरावरून ठरवते. `cli.pickAuto` लॉग इन असलेली खाती ठेवते, `ccam usage` वापरते त्याच `Fetcher.FetchAll` ने त्या सर्वांचा वापर एकाच वेळी विचारते, आणि `usage.Policy.Choose()` ला निवडू देते:

<!-- motion:next-auto -->

खाते त्याच्या 5 तासांच्या आणि साप्ताहिक windows वरून ठरवले जाते, आणि Claude Code च्या arguments मध्ये `--model` असेल तर त्या model च्या साप्ताहिक window वरूनही (नावात `opus`, `sonnet`, `haiku` किंवा `fable`). `Policy.Choose()` खात्यांची तीन गटांत विभागणी करते:

1. **जागा आहे**: प्रत्येक window `auto_max_usage` (`config.toml`, डीफॉल्ट 100) खाली. त्यांच्यातून strategy निवडते: `headroom` (डीफॉल्ट) सर्वात कमी 5 तासांचा वापर घेते, `expiring` सर्वात लवकर reset होणारी 5 तासांची window घेते, म्हणजे तिचा उरलेला quota वाया जाण्याआधी वापरला जातो.
2. **मर्यादेवर**: 100% खाली, पण `auto_max_usage` वर किंवा त्याहून जास्त. कोणत्याही खात्यात जागा नसेल तेव्हाच वापरली जातात, सर्वात कमी 5 तासांचा वापर आधी.
3. **भरलेली**: एखादी window 100% वर, किंवा अज्ञात 5 तासांचा वापर. कधीच निवडली जात नाहीत.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

काहीच निवडता येत नसेल तर ccam code 75 (`EX_TEMPFAIL`) ने बाहेर पडते, जेणेकरून script "नंतर पुन्हा प्रयत्न करा" आणि अपयश यांतील फरक ओळखू शकेल. `--wait` सह, ते त्याऐवजी पहिल्या भरलेल्या खात्याची अडवणारी window reset होईपर्यंत थांबते (reset time नंतर 30 सेकंदांनी, जेव्हा endpoint नवीन window दाखवतो), `ccam: no account has room; waiting until 15:10` छापते, आणि पुन्हा निवडते; कोणतेही खाते कधी reset होणार ते सांगत नसेल तर ते 10 मिनिटांनी पुन्हा पाहते. वाट पाहणे संपण्याआधी काहीच मोकळे झाले नाही तर ते 75 ने बाहेर पडते. हा शोध `ccam usage` चे नियम पाळतो: 60 सेकंदांपेक्षा नवीन आकडे पुन्हा वापरले जातात, प्रत्येक खात्याला 3 सेकंद असतात, आणि उत्तर न देणारे किंवा लॉगिन कालबाह्य झालेले खाते त्याच्या साठवलेल्या आकड्यांवरून ठरवले जाते. reset time उलटून गेलेली साठवलेली window 0% मोजली जाते, त्यामुळे ज्या खात्याची window दरम्यान रीसेट झाली आहे ते डावलले जात नाही. `--auto` कधीच tokens नूतनीकरण करत नाही. साध्या rotation च्या विपरीत, `--auto` तुम्ही शेवटी वापरलेले खाते निवडू शकते, जर त्यात अजूनही सर्वात जास्त जागा असेल.

### --continue

Claude Code चा `--continue` सध्याच्या फोल्डरमधील सर्वात अलीकडचे संभाषण उघडतो, जे नवीन खात्याला ते संभाषण दिसत असेल तरच चालते. म्हणून `--continue` सह ccam फक्त तुम्ही शेवटी वापरलेल्या खात्याच्या त्याच `history` गटातील खाती विचारात घेते, मग तुमच्या आर्ग्युमेंट्सच्या आधी `--continue` ठेवते. ज्या खात्यांचा इतिहास `local` आहे ती कोणाशीही शेअर करत नाहीत, त्यामुळे अशा खात्यावरून `--continue` नेहमी अयशस्वी होतो. आधी `ccam share <name> --history on` ने शेअरिंग चालू करा ([शेअरिंग](/docs/sharing) पाहा).

## मर्यादेवर पुन्हा चालवणे

`ccam next` हे थांबलेल्या interactive session साठी आहे. non-interactive run, म्हणजे script किंवा scheduled job मधील `claude -p`, तुमची वाट पाहू शकत नाही, त्यामुळे खात्याच्या नावाआधी `--fallback` लावल्यास ccam त्याच्यासाठी खाते बदलते:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` फक्त `-p` किंवा `--print` सोबतच चालते; नाहीतर ccam थांबते आणि `ccam next --auto --continue` कडे निर्देश करते.
2. standard input एकदाच वाचला जातो आणि प्रत्येक प्रयत्नाला दिला जातो.
3. ccam `exec` ऐवजी `claude` ला child process म्हणून सुरू करते, त्याचा standard output राखून ठेवते, आणि standard error येईल तसा पुढे देते.
4. `claude` 0 ने बाहेर पडल्यास, किंवा दुसऱ्या कारणाने अयशस्वी झाल्यास, राखून ठेवलेला output आणि exit code न बदलता पुढे दिले जातात.
5. Claude Code चा संदेश तसे सांगत असेल (`You've hit your limit`, `usage limit reached`), किंवा खात्याच्या वापराच्या शोधात त्याची 5 तासांची किंवा साप्ताहिक window 100% वर दिसत असेल, तर ते अपयश वापर मर्यादा म्हणून मोजले जाते.
6. मग त्या प्रयत्नाचा output टाकून दिला जातो आणि तीच कमांड दुसऱ्या लॉग इन असलेल्या खात्याने चालते: सर्वात कमी 5 तासांचा वापर असलेले खाते, मग ज्यांचा वापर शोधता येत नाही ती खाती, क्रमांकाच्या क्रमाने. 100% वरची खाती वगळली जातात, आणि प्रत्येक खाते जास्तीत जास्त एकदाच वापरून पाहिले जाते.
7. आर्ग्युमेंट्समध्ये `-c`, `--continue`, `-r` किंवा `--resume` असल्यास, फक्त त्याच इतिहास गटातील खाती वापरून पाहिली जातात, जेणेकरून संभाषण सापडेल.

पुढचे खाते `--auto` च्याच नियमांनी निवडले जाते (`auto_max_usage`, `auto_strategy`, `--model` ची window). एकही खाते उरले नाही तर ccam शेवटच्या प्रयत्नाचा output छापते, कोणत्याही लॉग इन असलेल्या खात्यात जागा उरलेली नाही असे सांगते, आणि code 75 ने बाहेर पडते. खात्याच्या नावाआधी `--wait` दिल्यास (`ccam --fallback --wait=2h work -p …`), ते त्याऐवजी पहिले खाते reset होण्याची वाट पाहते, मग प्रत्येक खाते पुन्हा वापरून पाहू शकते; ते नेहमी आधी किमान एकदा थांबते, जेणेकरून नुकतीच मर्यादा गाठलेले खाते लगेच पुन्हा चालवले जात नाही. `ccam --dry-run --fallback work -p …` पहिल्या प्रयत्नाची कमांड छापते.

interactive sessions आपोआप बदलले जात नाहीत. त्याऐवजी, `ccam statusline install` सह, 5 तासांचा किंवा साप्ताहिक वापर 90% ला पोहोचला की स्टेटस लाइन `→ ccam next --auto --continue` जोडते ([वापर मर्यादा](/features/usage)).

## सुरक्षितता आणि अपवादात्मक स्थिती

- लॉगिन नसलेली खाती कधीच निवडली जात नाहीत. दुसरे लॉग इन असलेले खाते नसल्यास, साधे `ccam next` `no logged-in account other than <name>` ने थांबते.
- `--continue` सह आणि इतिहास शेअर करणारे एकही खाते नसल्यास: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- `--auto` सह आणि एकही लॉग इन असलेले खाते नसल्यास: `no logged-in account to pick`. प्रत्येक खाते भरलेले असेल किंवा त्याचा वापर अज्ञात असेल तेव्हा: `` no account with known usage and headroom (see `ccam usage`) ``. दोन्ही exit code 75 ने संपतात.
- `--max-usage`, `--strategy` आणि `--wait` फक्त `--auto` सह काम करतात, आणि 1-100 बाहेरचे मूल्य किंवा दुसरी strategy नाकारली जाते.
- `--dry-run` सह सुद्धा `--auto` नेटवर्क requests करते. ते फक्त वाचते, आणि token फाइल्समध्ये कधीच लिहीत नाही.
- निवडलेले खाते `ccam <name>` च्याच मार्गाने सुरू होते, त्यामुळे त्याच तपासण्या लागू होतात, यात तोच इतिहास शेअर करणारे दुसरे Claude Code या फोल्डरमध्ये आधीच चालू असल्याचा इशाराही येतो.
- `last-used` फक्त Claude Code प्रत्यक्ष सुरू झाल्यावर लिहिले जाते. `--dry-run` त्याला हात लावत नाही, त्यामुळे dry run rotation पुढे सरकवत नाही.
- `--fallback` प्रयत्न संपेपर्यंत standard output राखून ठेवते, त्यामुळे `--output-format stream-json` सारखा streaming output शेवटी एकाच वेळी येतो. पहिल्या प्रयत्नाने मर्यादा गाठण्यापूर्वी बदललेल्या फाइल्स बदललेल्याच राहतात; पुढचा प्रयत्न prompt पुन्हा सुरुवातीपासून सुरू करतो.

## संबंधित

- [वापर मर्यादा](/features/usage)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [वापर मर्यादा workflow](/guides/usage-limit-workflow)
