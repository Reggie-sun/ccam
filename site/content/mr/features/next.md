---
title: पुढच्या खात्यावर जाणे
description: ccam next दुसऱ्या खात्याने Claude Code उघडते; --auto सर्वाधिक जागेचे खाते निवडते, --continue संभाषण पुढे नेते, --fallback -p jobs पुन्हा चालवते.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` तुम्ही नाव न देता दुसऱ्या खात्याने Claude Code उघडते. एखादे खाते वापर मर्यादेला पोहोचल्यावर ते वापरा: साधे `ccam next` तुम्ही शेवटी वापरलेल्या खात्यानंतरचे खाते घेते, `--auto` सर्वात जास्त जागा असलेले खाते घेते, आणि `--continue` या फोल्डरमधील शेवटचे संभाषण नवीन खात्यात पुढे नेते, दोन्ही खाती इतिहास शेअर करत असोत वा नसोत. Claude Code मध्ये, मर्यादेमुळे सत्र थांबले की [ccam plugin](#inside-claude-code) तुम्हाला सांगते आणि `/ccam:switch` संभाषण पुढे सोपवते. `claude -p` चालवणाऱ्या scripts साठी `ccam --fallback` हा बदल आपोआप करते ([मर्यादेवर पुन्हा चालवणे](#run-again-at-the-limit) पाहा).

## वापर

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam हा बदल standard error वर `→ <from> → <to>` असा छापते, `--auto` ने निवड केली असल्यास तिच्या कारणासह, मग `ccam <name>` प्रमाणेच नवीन खात्याने Claude Code सुरू करते. `<from>` म्हणजे तुम्ही शेवटी वापरलेले खाते, किंवा `--continue` सह ते संभाषण ज्या खात्याचे आहे ते खाते. `--` नंतरचे आर्ग्युमेंट्स न बदलता `claude` ला जातात.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--auto` | वापरावरून लॉग इन केलेले खाते निवडा (खाली पहा), 100% वरची खाती वगळून | बंद |
| `--max-usage <1-100>` | `--auto` सह: या वापराखालील खात्यांना प्राधान्य द्या | `auto_max_usage`, नाहीतर 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` सह: सर्वात कमी 5 तासांचा वापर, सर्वात लवकर होणारा 5 तासांचा reset, किंवा सर्वात आधी न वापरता वाया जाणारा साप्ताहिक quota | `auto_strategy`, नाहीतर `headroom` |
| `--cross-org` | `--auto` सह: दुसऱ्या organization ची खातीसुद्धा निवडा (खाली पहा) | `cross_org`, नाहीतर बंद |
| `--wait[=<duration>]` | `--auto` सह: कोणत्याही खात्यात जागा नसेल तर पहिले खाते reset होण्याची वाट पाहा, जास्तीत जास्त इतका वेळ | बंद; एकटे दिल्यास `5h` |
| `-c`, `--continue` | या फोल्डरमधील शेवटचे संभाषण नवीन खात्यात पुढे न्या (खाली पहा) | बंद |
| `-- <claude args>` | ccam च्या स्वतःच्या आर्ग्युमेंट्सनंतर `claude` ला दिले जाणारे आर्ग्युमेंट्स | काहीही नाही |
| `--dry-run` | पाळला जातो: Claude Code सुरू न करता किंवा सुरू केल्याची नोंद न करता `→ <from> → <to>` आणि चालवली जाणारी कमांड लाइन छापतो | बंद |

`--auto` आणि `--continue` एकत्र चालतात. `--continue` सह, `--resume` चे आर्ग्युमेंट्स तुमच्या स्वतःच्या आर्ग्युमेंट्सच्या आधी ठेवले जातात.

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

1. **जागा आहे**: प्रत्येक window `auto_max_usage` (`config.toml`, डीफॉल्ट 100) खाली. त्यांच्यातून strategy निवडते: `headroom` (डीफॉल्ट) सर्वात कमी 5 तासांचा वापर घेते, `expiring` सर्वात लवकर reset होणारी 5 तासांची window घेते, म्हणजे तिचा उरलेला quota वाया जाण्याआधी वापरला जातो; `balance` आठवडा reset होईपर्यंत दररोज सर्वात जास्त साप्ताहिक quota उरलेले खाते घेते, आणि चालू असलेले प्रत्येक ccam सत्र त्याच्या विरुद्ध मोजते, त्यामुळे साप्ताहिक quota न वापरता राहत नाही आणि समांतर सत्रे विखुरली जातात.
2. **मर्यादेवर**: 100% खाली, पण `auto_max_usage` वर किंवा त्याहून जास्त. कोणत्याही खात्यात जागा नसेल तेव्हाच वापरली जातात, सर्वात कमी 5 तासांचा वापर आधी.
3. **भरलेली**: एखादी window 100% वर, किंवा अज्ञात 5 तासांचा वापर. कधीच निवडली जात नाहीत.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizations नुसार मर्यादा.** आपोआप होणाऱ्या निवडी तुम्ही सोडत असलेल्या खात्यासारख्याच खात्यांमध्ये राहतात: तीच team किंवा enterprise organization, किंवा personal plans (Pro, Max) एकत्र, त्यामुळे कंपनीच्या खात्यावरील काम खासगी खात्यावर पुढे चालू होत नाही, किंवा उलटही नाही. ccam ला प्रत्येक खात्याची organization त्याच्या लॉगिनवरून कळते (`ccam login` ती नोंदवते; जुनी लॉगिन `.claude.json` मधून वाचली जातात); ज्या खात्याची organization ccam सांगू शकत नाही ते कुठेही बसते. एकाच Claude खात्यात लॉग इन असलेले दोन फोल्डर त्याच्या मर्यादा शेअर करतात, त्यामुळे एक दुसऱ्याची जागा कधीच घेत नाही. `--cross-org`, किंवा `config.toml` मधील `cross_org = true`, organization चा नियम उठवतो.

काहीच निवडता येत नसेल तर ccam code 75 (`EX_TEMPFAIL`) ने बाहेर पडते, जेणेकरून script "नंतर पुन्हा प्रयत्न करा" आणि अपयश यांतील फरक ओळखू शकेल. `--wait` सह, ते त्याऐवजी पहिल्या भरलेल्या खात्याची अडवणारी window reset होईपर्यंत थांबते (reset time नंतर 30 सेकंदांनी, जेव्हा endpoint नवीन window दाखवतो), `ccam: no account has room; waiting until 15:10` छापते, आणि पुन्हा निवडते; कोणतेही खाते कधी reset होणार ते सांगत नसेल तर ते 10 मिनिटांनी पुन्हा पाहते. वाट पाहणे संपण्याआधी काहीच मोकळे झाले नाही तर ते 75 ने बाहेर पडते. हा शोध `ccam usage` चे नियम पाळतो: 60 सेकंदांपेक्षा नवीन आकडे पुन्हा वापरले जातात, प्रत्येक खात्याला 3 सेकंद असतात, आणि उत्तर न देणारे किंवा लॉगिन कालबाह्य झालेले खाते त्याच्या साठवलेल्या आकड्यांवरून ठरवले जाते. reset time उलटून गेलेली साठवलेली window 0% मोजली जाते, त्यामुळे ज्या खात्याची window दरम्यान रीसेट झाली आहे ते डावलले जात नाही. `--auto` कधीच tokens नूतनीकरण करत नाही. साध्या rotation च्या विपरीत, `--auto` तुम्ही शेवटी वापरलेले खाते निवडू शकते, जर त्यात अजूनही सर्वात जास्त जागा असेल.

### --continue

1. ccam प्रत्येक खात्यात या फोल्डरचे शेवटचे संभाषण शोधते: `projects/<folder>/` मधील सर्वात नवीन `<session>.jsonl` (फोल्डरचे नाव म्हणजे कामाचा फोल्डर, ज्यातील अक्षर किंवा अंक नसलेले प्रत्येक चिन्ह `-` मध्ये बदललेले असते, जसे Claude Code नाव देते). सत्राने नोंद ठेवली असेल (वापर मर्यादा किंवा `/ccam:switch`, खाली पहा) आणि त्यानंतर काहीही नवीन लिहिले गेले नसेल, तर त्याऐवजी नेमके तेच सत्र वापरले जाते.
2. ते संभाषण ज्या खात्याचे आहे त्याच खात्यावरून तुम्ही बदलत असता, आणि नवीन खाते वरीलप्रमाणे, त्याच organization मधून निवडले जाते.
3. नवीन खाते तो इतिहास शेअर करत असेल (तोच गट फोल्डर), तर ccam `claude --resume <session id>` सुरू करते, म्हणजे नेमके तेच संभाषण (`--continue` फोल्डरमधील जे सर्वात नवीन असेल ते उघडेल).
4. नाहीतर ते `claude --resume <the conversation's file> --fork-session` सुरू करते: Claude Code संपूर्ण संभाषण नवीन session id खाली नवीन खात्याच्या स्वतःच्या इतिहासात कॉपी करते, आणि मूळ संभाषण होते तसेच राहते. त्यानंतर नवीन खात्याचा `--continue` ते शोधतो.

संभाषण किती मोठे आहे ते ccam सांगते: नवीन खाते आपल्या पहिल्या उत्तराबरोबर ते सगळे पुन्हा वाचते. Prompt caches कधीही organizations मध्ये शेअर होत नाहीत, त्यामुळे `--cross-org` सह त्या पहिल्या उत्तराला पूर्ण आकाराचा खर्च येतो. गटात `history` शेअर केल्याने ही प्रत टळते, एकच संभाषण राहते, आणि `/rewind` ज्या checkpoints मधून code परत आणते ते (`file-history/`) सुद्धा शेअर होतात.

## Claude Code मध्ये

[ccam plugin](/docs/install#claude-code-plugin) दोन hooks जोडते (Claude Code 2.1.78 किंवा नंतरचे):

- एखादा turn वापर मर्यादेवर थांबला (`You've hit your session limit`, साप्ताहिक किंवा model ची मर्यादा), की ते संभाषणाची नोंद ठेवते आणि सूचना देऊ शकणाऱ्या terminals मध्ये (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) desktop सूचना दाखवते: `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` सर्व्हरचे थोड्या वेळाचे throttling किंवा API key ची rate limit ही वापर मर्यादा नाही, आणि तिला हात लावला जात नाही. Claude Code 2.1.234 आणि नंतरचे सत्रात थांबून राहते आणि 5 तासांच्या reset नंतर आपोआप पुढे चालू ठेवते; सूचना तुम्हाला त्याऐवजी खाते बदलण्याचा पर्याय देते.
- `/ccam:switch [account]` संभाषणाची नोंद ठेवते आणि कोणतीही request होण्यापूर्वी turn संपवते: `/exit` टाइप करा, मग `ccam next --continue` ते तुम्ही नाव दिलेल्या खात्याने, किंवा पुढच्या खात्याने, पुढे चालू ठेवते.

ही नोंद `~/.local/state/ccam/handoff/<account>.json` मध्ये 12 तास राहते आणि पुढचे सत्र सुरू होणार असताना एकदाच वापरली जाते; सुरू करणे अयशस्वी झाले तर ती पुढच्या प्रयत्नासाठी तशीच राहते. ती सत्र ज्या फोल्डरमध्ये सुरू झाले त्या फोल्डरची असते, Claude Code च्या shell मध्ये `cd` केल्यानंतरही. तिच्यात नाव असलेल्या खात्याचे नाव नंतर बदलले असेल किंवा ते खाते काढून टाकले असेल, तर ccam तसे सांगते आणि नेहमीप्रमाणे निवडते. `config.toml` मध्ये `relay = true` असेल, तर ccam स्वतःची जागा Claude Code ला देण्याऐवजी interactive सत्रांच्या शेजारी राहते, आणि `/ccam:switch` नंतर त्या सत्राचे संभाषण आपोआप पुढे नेते (दुसऱ्या सत्राने ठेवलेल्या नोंदीवर कृती केली जात नाही). त्यासाठी प्रत्येक सत्रामागे वाट पाहणारी एक ccam process लागते, म्हणूनच हे डीफॉल्टनुसार बंद असते. पुढे नेलेले सत्र `ccam next --continue` ने सुरू होते, पहिले सत्र ज्या आर्ग्युमेंट्सनी सुरू झाले होते (जसे `--model`) त्यांच्याशिवाय.

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
8. `--auto` प्रमाणेच फक्त त्याच organization ची खाती वापरून पाहिली जातात (`cross_org = true` हे उठवते).

पुढचे खाते `--auto` च्याच नियमांनी निवडले जाते (`auto_max_usage`, `auto_strategy`, `--model` ची window). एकही खाते उरले नाही तर ccam शेवटच्या प्रयत्नाचा output छापते, कोणत्याही लॉग इन असलेल्या खात्यात जागा उरलेली नाही असे सांगते, आणि code 75 ने बाहेर पडते. खात्याच्या नावाआधी `--wait` दिल्यास (`ccam --fallback --wait=2h work -p …`), ते त्याऐवजी पहिले खाते reset होण्याची वाट पाहते, मग प्रत्येक खाते पुन्हा वापरून पाहू शकते; ते नेहमी आधी किमान एकदा थांबते, जेणेकरून नुकतीच मर्यादा गाठलेले खाते लगेच पुन्हा चालवले जात नाही. `ccam --dry-run --fallback work -p …` पहिल्या प्रयत्नाची कमांड छापते.

interactive sessions आपोआप बदलले जात नाहीत. त्याऐवजी, `ccam statusline install` सह, 5 तासांचा किंवा साप्ताहिक वापर 90% ला पोहोचला की स्टेटस लाइन `→ ccam next --auto --continue` जोडते ([वापर मर्यादा](/features/usage)).

## सुरक्षितता आणि अपवादात्मक स्थिती

- लॉगिन नसलेली खाती कधीच निवडली जात नाहीत. दुसरे लॉग इन असलेले खाते नसल्यास, साधे `ccam next` `no logged-in account other than <name>` ने थांबते.
- `--continue` सह आणि फोल्डरमध्ये एकही संभाषण नसल्यास: `no account has a conversation in <folder> to continue`.
- organization च्या नियमामुळे एकही खाते उरले नाही तर: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` सह आणि एकही लॉग इन असलेले खाते नसल्यास: `no logged-in account to pick`. प्रत्येक खाते भरलेले असेल किंवा त्याचा वापर अज्ञात असेल तेव्हा: `` no account with known usage and headroom (see `ccam usage`) ``. दोन्ही exit code 75 ने संपतात.
- `--max-usage`, `--strategy`, `--wait` आणि `--cross-org` फक्त `--auto` सह काम करतात, आणि 1-100 बाहेरचे मूल्य किंवा दुसरी strategy नाकारली जाते.
- fork संभाषण कॉपी करतो, खात्याची इतर स्थिती नाही: इतिहास शेअर केलेला नसेल तर मूळ खात्याच्या `file-history/` मधील checkpoints सोबत येत नाहीत.
- `--dry-run` सह सुद्धा `--auto` नेटवर्क requests करते. ते फक्त वाचते, आणि token फाइल्समध्ये कधीच लिहीत नाही.
- निवडलेले खाते `ccam <name>` च्याच मार्गाने सुरू होते, त्यामुळे त्याच तपासण्या लागू होतात, यात तोच इतिहास शेअर करणारे दुसरे Claude Code या फोल्डरमध्ये आधीच चालू असल्याचा इशाराही येतो.
- `last-used` फक्त Claude Code प्रत्यक्ष सुरू झाल्यावर लिहिले जाते. `--dry-run` त्याला हात लावत नाही, त्यामुळे dry run rotation पुढे सरकवत नाही.
- `--fallback` प्रयत्न संपेपर्यंत standard output राखून ठेवते, त्यामुळे `--output-format stream-json` सारखा streaming output शेवटी एकाच वेळी येतो. पहिल्या प्रयत्नाने मर्यादा गाठण्यापूर्वी बदललेल्या फाइल्स बदललेल्याच राहतात; पुढचा प्रयत्न prompt पुन्हा सुरुवातीपासून सुरू करतो.

## संबंधित

- [वापर मर्यादा](/features/usage)
- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [वापर मर्यादा workflow](/guides/usage-limit-workflow)
