---
title: अगले अकाउंट पर स्विच करना
description: ccam next किसी दूसरे अकाउंट से Claude Code खोलता है; --auto सबसे ज़्यादा गुंजाइश वाला चुनता है, --continue बातचीत आगे ले जाता है, --fallback -p jobs दोहराता है।
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` आपके नाम बताए बिना Claude Code को किसी दूसरे अकाउंट से खोलता है। इसे तब इस्तेमाल करें जब कोई अकाउंट अपनी उपयोग सीमा तक पहुँच जाए: सादा `ccam next` पिछली बार इस्तेमाल किए गए अकाउंट के बाद वाला अकाउंट लेता है, `--auto` वह लेता है जिसमें सबसे ज़्यादा गुंजाइश हो, और `--continue` इस फ़ोल्डर की आख़िरी बातचीत को नए अकाउंट में ले जाता है, चाहे दोनों इतिहास साझा करते हों या नहीं। Claude Code के अंदर, [ccam plugin](#inside-claude-code) बताता है कि कब कोई सीमा session को रोकती है, और `/ccam:switch` बातचीत आगे सौंप देता है। `claude -p` चलाने वाली scripts के लिए, `ccam --fallback` यह स्विच अपने-आप करता है ([सीमा पर फिर से चलाना](#run-again-at-the-limit) देखें)।

## उपयोग

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam स्विच को standard error पर `→ <from> → <to>` के रूप में प्रिंट करता है, `--auto` वाले चुनाव के कारण के साथ, फिर नए अकाउंट से Claude Code उसी तरह शुरू करता है जैसे `ccam <name>` करता है। `<from>` वह अकाउंट है जिसे आपने पिछली बार इस्तेमाल किया, या `--continue` के साथ वह अकाउंट जिसकी बातचीत है। `--` के बाद वाले आर्ग्युमेंट बिना बदलाव के `claude` को जाते हैं।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--auto` | उपयोग के हिसाब से लॉगिन किया हुआ अकाउंट चुनें (नीचे देखें), 100% वाले अकाउंट छोड़कर | बंद |
| `--max-usage <1-100>` | `--auto` के साथ: इस उपयोग से नीचे वाले अकाउंट्स को प्राथमिकता दें | `auto_max_usage`, वरना 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` के साथ: सबसे कम 5 घंटे का उपयोग, सबसे जल्दी होने वाला 5 घंटे का reset, या वह साप्ताहिक कोटा जो सबसे पहले बिना इस्तेमाल के रह जाएगा | `auto_strategy`, वरना `headroom` |
| `--cross-org` | `--auto` के साथ: किसी दूसरे संगठन के अकाउंट भी चुनें (नीचे देखें) | `cross_org`, वरना बंद |
| `--wait[=<duration>]` | `--auto` के साथ: जब किसी अकाउंट में गुंजाइश न हो, तो पहले अकाउंट के reset होने का इंतज़ार करें, ज़्यादा से ज़्यादा इतनी देर | बंद; अकेले देने पर `5h` |
| `-c`, `--continue` | इस फ़ोल्डर की आख़िरी बातचीत को नए अकाउंट के रूप में आगे ले जाएँ (नीचे देखें) | बंद |
| `-- <claude args>` | ccam के अपने आर्ग्युमेंट के बाद `claude` को दिए जाने वाले आर्ग्युमेंट | कोई नहीं |
| `--dry-run` | माना जाता है: Claude Code शुरू किए बिना या शुरुआत दर्ज किए बिना `→ <from> → <to>` और चलाई जाने वाली कमांड लाइन प्रिंट करता है | बंद |

`--auto` और `--continue` एक साथ काम करते हैं। `--continue` के साथ, `--resume` वाले आर्ग्युमेंट आपके अपने आर्ग्युमेंट से पहले रखे जाते हैं।

Claude Code के अपने विकल्प, जैसे `--dangerously-skip-permissions` या `--model`, `--` के बाद आते हैं। `ccam <name>` के उलट, जो नाम के बाद की हर चीज़ आगे भेज देता है, `ccam next` अपने आर्ग्युमेंट ख़ुद पढ़ता है, इसलिए `--` के बिना यह `ccam: unknown flag: --dangerously-skip-permissions` के साथ रुक जाता है।

## यह कैसे काम करता है

ccam से होने वाली हर शुरुआत समय को `~/.local/state/ccam/last-used/<name>` में लिखती है। `ccam next` इन फ़ाइलों को पढ़ता है, जाँचता है कि कौन से अकाउंट लॉगिन किए हुए हैं, और चुनता है।

1. `cli.candidates`, `config.toml` से हर अकाउंट के लिए एक `next.Candidate` बनाता है: अकाउंट, उसके आख़िरी बार शुरू होने का समय (`launcher.LastUsed`) और यह कि उसका लॉगिन है या नहीं (`creds.Has`: पढ़ने लायक `.credentials.json`, या macOS पर कीचेन item)।
2. `--auto` के बिना, `next.Rotate` अकाउंट्स को स्लॉट नंबर से क्रम में रखता है और सबसे हाल में शुरू हुआ अकाउंट खोजता है। यही वह अकाउंट है जिससे आप स्विच कर रहे हैं। अगर अभी तक कोई अकाउंट शुरू नहीं हुआ, तो यह सबसे छोटा स्लॉट होता है।
3. वहाँ से यह स्लॉट क्रम में आगे बढ़ता है, आख़िरी स्लॉट से वापस पहले पर घूमते हुए, और पहला लॉगिन किया हुआ अकाउंट लौटाता है।
4. `launcher.LaunchWith` चुने गए अकाउंट को ठीक वैसे तैयार करके शुरू करता है जैसे `ccam <name>` करता: यह अकाउंट के लिंक जाँचता है, MCP साझा करना लागू करता है, session रजिस्टर करता है और `syscall.Exec` के ज़रिए ccam की जगह `claude` रखता है।

खोज पिछली बार इस्तेमाल किए गए अकाउंट से शुरू होती है और उसके बाद पहले लॉगिन किए हुए अकाउंट पर रुकती है:

<!-- motion:next-rotate -->

`work` सबसे हाल में शुरू हुआ था, इसलिए खोज उसके बाद से शुरू होती है। `side` का लॉगिन नहीं है और उसे छोड़ दिया जाता है, और `lab` चुना जाता है। अगर `lab` भी लॉग आउट होता, तो खोज घूमकर `personal` पर आती (आरेख में `% len`)। रोटेशन कभी वह अकाउंट नहीं लौटाता जिससे आपने शुरू किया: कोई दूसरा लॉगिन किया हुआ अकाउंट न होने पर यह त्रुटि के साथ रुकता है।

### --auto

`--auto` स्लॉट क्रम की जगह उपयोग से तय करता है। `cli.pickAuto` लॉगिन किए हुए अकाउंट रखता है, उसी `Fetcher.FetchAll` से जिसे `ccam usage` इस्तेमाल करता है, एक साथ उन सबका उपयोग पूछता है, और `usage.Policy.Choose()` को चुनने देता है:

<!-- motion:next-auto -->

हर अकाउंट को उसकी 5 घंटे और साप्ताहिक विंडो से आँका जाता है, और जब Claude Code के आर्ग्युमेंट में `--model` हो (`--` के बाद), तो उस मॉडल की साप्ताहिक विंडो से भी (नाम में `opus`, `sonnet`, `haiku` या `fable`)। `Policy.Choose()` अकाउंट्स को तीन हिस्सों में बाँटता है:

1. **गुंजाइश**: हर विंडो `auto_max_usage` (`config.toml`, डिफ़ॉल्ट 100) से नीचे। इनमें से strategy चुनती है: `headroom` (डिफ़ॉल्ट) सबसे कम 5 घंटे का उपयोग लेती है; `expiring` वह 5 घंटे की विंडो लेती है जो सबसे जल्दी reset होती है, ताकि उसका बचा हुआ quota ख़त्म होने से पहले इस्तेमाल हो जाए; `balance` वह अकाउंट लेती है जिसके हफ़्ते के reset होने तक हर दिन के हिसाब से सबसे ज़्यादा साप्ताहिक quota बचा हो, और उस अकाउंट पर चल रहे हर ccam session को उसके ख़िलाफ़ गिनती है, ताकि साप्ताहिक quota बिना इस्तेमाल के न रह जाए और समानांतर sessions बँट जाएँ।
2. **सीमा से ऊपर**: 100% से नीचे, लेकिन `auto_max_usage` पर या उससे ऊपर। केवल तब इस्तेमाल होते हैं जब किसी अकाउंट में गुंजाइश न हो, सबसे कम 5 घंटे के उपयोग वाला पहले।
3. **भरे हुए**: कोई विंडो 100% पर, या अज्ञात 5 घंटे का उपयोग। कभी नहीं चुने जाते।

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**संगठन।** अपने-आप होने वाले चुनाव उसी तरह के अकाउंट्स के साथ रहते हैं जैसा वह अकाउंट जिसे आप छोड़ रहे हैं: वही team या enterprise संगठन, या personal प्लान (Pro, Max) आपस में, ताकि किसी कंपनी अकाउंट का काम किसी निजी अकाउंट पर, या इसका उल्टा, जारी न हो। ccam हर अकाउंट का संगठन उसके लॉगिन से जानता है (`ccam login` इसे दर्ज करता है; पुराने लॉगिन `.claude.json` से पढ़े जाते हैं); जिस अकाउंट का संगठन वह नहीं बता पाता, वह कहीं भी फ़िट होता है। एक ही Claude अकाउंट में लॉगिन दो फ़ोल्डर उसकी सीमाएँ साझा करते हैं, इसलिए उनमें से एक कभी दूसरे की जगह नहीं लेता। `--cross-org`, या `config.toml` में `cross_org = true`, संगठन वाला नियम हटा देता है।

जब कुछ भी नहीं चुना जा सकता, तो ccam code 75 (`EX_TEMPFAIL`) के साथ बाहर निकलता है, ताकि script "बाद में फिर कोशिश करें" को विफलता से अलग पहचान सके। `--wait` के साथ, यह इसके बजाय तब तक रुकता है जब तक पहले भरे हुए अकाउंट की रोकने वाली विंडो reset न हो जाए (reset समय के 30 सेकंड बाद, जब endpoint नई विंडो दिखाता है), `ccam: no account has room; waiting until 15:10` प्रिंट करता है, और फिर से चुनता है; जब कोई अकाउंट नहीं बताता कि वह कब reset होगा, तो यह 10 मिनट बाद फिर देखता है। अगर इंतज़ार ख़त्म होने से पहले कुछ ख़ाली नहीं होता, तो यह 75 के साथ बाहर निकलता है। यह खोज `ccam usage` के नियमों का पालन करती है: 60 सेकंड से कम पुरानी संख्याएँ दोबारा इस्तेमाल होती हैं, हर अकाउंट के पास 3 सेकंड होते हैं, और जो अकाउंट जवाब नहीं देता, या जिसका लॉगिन एक्सपायर हो गया है, उसे उसकी सेव की गई संख्याओं से आँका जाता है। सेव की गई जिस विंडो का रीसेट समय बीत चुका है, उसे 0% माना जाता है, ताकि जिस अकाउंट की विंडो तब से रीसेट हो चुकी है, वह छूट न जाए। `--auto` कभी tokens renew नहीं करता। सादे रोटेशन के उलट, `--auto` पिछली बार इस्तेमाल किया गया अकाउंट भी चुन सकता है, अगर उसमें अब भी सबसे ज़्यादा गुंजाइश हो।

### --continue

1. ccam हर अकाउंट में इस फ़ोल्डर की आख़िरी बातचीत खोजता है: `projects/<folder>/` में सबसे नई `<session>.jsonl` (फ़ोल्डर का नाम काम करने वाला फ़ोल्डर होता है, जिसमें अक्षर या अंक के अलावा हर character `-` में बदल जाता है, जैसे Claude Code इसका नाम रखता है)। अगर session ने कोई नोट छोड़ा है (उपयोग सीमा या `/ccam:switch`, नीचे देखें) और तब से कुछ नया नहीं लिखा गया, तो इसके बजाय ठीक वही session इस्तेमाल होता है।
2. जिस अकाउंट की वह बातचीत है, उसी से आप स्विच करते हैं, और नया अकाउंट ऊपर बताए तरीक़े से, उसी संगठन के भीतर, चुना जाता है।
3. अगर नया अकाउंट वह इतिहास साझा करता है (वही ग्रुप फ़ोल्डर), तो ccam `claude --resume <session id>` शुरू करता है, यानी ठीक वही बातचीत (`--continue` फ़ोल्डर की सबसे नई बातचीत खोलता, चाहे वह कोई भी हो)।
4. नहीं तो यह `claude --resume <the conversation's file> --fork-session` शुरू करता है: Claude Code पूरी बातचीत को एक नए session id के तहत नए अकाउंट के अपने इतिहास में कॉपी करता है, और मूल बातचीत जैसी थी वैसी रहती है। उसके बाद से नए अकाउंट का `--continue` उसे ढूँढ लेता है।

ccam बताता है कि बातचीत कितनी बड़ी है: नया अकाउंट अपने पहले जवाब के साथ उसे पूरा फिर से पढ़ता है। Prompt caches कभी संगठनों के बीच साझा नहीं होते, इसलिए `--cross-org` के साथ उस पहले जवाब की लागत पूरे आकार जितनी होती है। किसी ग्रुप में `history` साझा करने से कॉपी की ज़रूरत नहीं पड़ती, एक ही बातचीत रहती है, और वे checkpoints भी साझा होते हैं जिनसे `/rewind` कोड वापस लाता है (`file-history/`)।

## Claude Code के अंदर

[ccam plugin](/docs/install#claude-code-plugin) दो hooks जोड़ता है (Claude Code 2.1.78 या उसके बाद):

- जब कोई turn उपयोग सीमा पर रुकता है (`You've hit your session limit`, साप्ताहिक या किसी मॉडल की सीमा), तो यह बातचीत को दर्ज करता है और उन टर्मिनलों में desktop notification दिखाता है जो इसे समर्थन करते हैं (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` सर्वर का थोड़ी देर का throttle या किसी API key की rate limit उपयोग सीमा नहीं है और उसे छोड़ दिया जाता है। Claude Code 2.1.234 और उसके बाद के वर्ज़न session में इंतज़ार भी करते हैं और 5 घंटे के reset के बाद अपने-आप जारी रखते हैं; notification आपको इसके बजाय स्विच करने का विकल्प देता है।
- `/ccam:switch [account]` बातचीत को दर्ज करता है और कोई request भेजे जाने से पहले turn ख़त्म कर देता है: `/exit` टाइप करें, फिर `ccam next --continue` उसे आपके बताए अकाउंट के रूप में, या अगले अकाउंट के रूप में, आगे ले जाता है।

नोट `~/.local/state/ccam/handoff/<account>.json` में 12 घंटे तक रहता है और एक बार इस्तेमाल होता है, जब अगला session शुरू होने वाला हो; जो शुरुआत विफल हो जाए, वह नोट को अगली कोशिश के लिए छोड़ देती है। नोट उसी फ़ोल्डर का होता है जिसमें session शुरू हुआ था, Claude Code के shell में `cd` करने के बाद भी। अगर नोट में जिस अकाउंट का नाम है, उसका नाम तब से बदल दिया गया हो या उसे हटा दिया गया हो, तो ccam यह बताता है और हमेशा की तरह चुनता है। `config.toml` में `relay = true` होने पर ccam ख़ुद को Claude Code से बदलने के बजाय interactive sessions के साथ-साथ बना रहता है, और `/ccam:switch` के बाद उस session की बातचीत अपने-आप आगे ले जाता है (किसी दूसरे session के छोड़े नोट पर यह कुछ नहीं करता)। इसके लिए हर session पर एक इंतज़ार करता ccam process लगता है, इसीलिए यह डिफ़ॉल्ट रूप से बंद है। इस तरह आगे ले जाया गया session `ccam next --continue` से शुरू होता है, उन आर्ग्युमेंट के बिना जिनसे पहला session शुरू हुआ था, जैसे `--model`।

## सीमा पर फिर से चलाना

`ccam next` उस interactive session के लिए है जो रुक गया है। एक non-interactive run, यानी किसी script या scheduled job में `claude -p`, आपका इंतज़ार नहीं कर सकता, इसलिए अकाउंट नाम से पहले `--fallback` लगाने से ccam उसके लिए स्विच करता है:

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

1. `--fallback` केवल `-p` या `--print` के साथ काम करता है; वरना ccam रुक जाता है और `ccam next --auto --continue` की ओर इशारा करता है।
2. Standard input एक बार पढ़ा जाता है और हर प्रयास को दिया जाता है।
3. ccam `exec` के बजाय `claude` को एक child process के रूप में शुरू करता है, उसका standard output रोककर रखता है, और standard error को आते ही आगे भेजता है।
4. जब `claude` 0 के साथ बाहर निकलता है, या किसी दूसरे कारण से विफल होता है, तो रोका गया output और exit code बिना बदलाव आगे दिए जाते हैं।
5. विफलता को उपयोग सीमा तब माना जाता है जब Claude Code का संदेश यही कहे (`You've hit your limit`, `usage limit reached`), या जब अकाउंट की उपयोग खोज उसकी 5 घंटे या साप्ताहिक विंडो को 100% पर दिखाए।
6. तब उस प्रयास का output छोड़ दिया जाता है और वही कमांड किसी दूसरे लॉग इन अकाउंट के रूप में चलती है: सबसे कम 5 घंटे के उपयोग वाला, फिर वे अकाउंट जिनका उपयोग पता नहीं किया जा सकता, नंबर के क्रम में। 100% वाले अकाउंट छोड़ दिए जाते हैं, और हर अकाउंट को ज़्यादा से ज़्यादा एक बार आज़माया जाता है।
7. आर्ग्युमेंट में `-c`, `--continue`, `-r` या `--resume` होने पर, केवल उसी इतिहास ग्रुप के अकाउंट आज़माए जाते हैं, ताकि बातचीत मिल सके।
8. `--auto` की तरह, केवल उसी संगठन के अकाउंट आज़माए जाते हैं (`cross_org = true` यह पाबंदी हटाता है)।

अगला अकाउंट `--auto` वाले नियमों से ही चुना जाता है (`auto_max_usage`, `auto_strategy`, `--model` की विंडो)। जब कोई अकाउंट नहीं बचता, तो ccam आख़िरी प्रयास का output छापता है, बताता है कि किसी लॉग इन अकाउंट में गुंजाइश नहीं बची, और code 75 के साथ बाहर निकलता है। अकाउंट नाम से पहले `--wait` के साथ (`ccam --fallback --wait=2h work -p …`), यह इसके बजाय पहले अकाउंट के reset होने का इंतज़ार करता है, फिर हर अकाउंट को दोबारा आज़मा सकता है; यह हमेशा पहले कम से कम एक बार रुकता है, ताकि अभी-अभी सीमा पर पहुँचा अकाउंट तुरंत दोबारा न चले। `ccam --dry-run --fallback work -p …` पहले प्रयास की कमांड छापता है।

Interactive sessions अपने-आप स्विच नहीं होते। इसके बजाय, `ccam statusline install` के साथ, 5 घंटे या साप्ताहिक उपयोग 90% तक पहुँचते ही स्टेटस लाइन `→ ccam next --auto --continue` जोड़ देती है ([उपयोग सीमाएँ](/features/usage))।

## सुरक्षा और विशेष स्थितियाँ

- बिना लॉगिन वाले अकाउंट कभी नहीं चुने जाते। कोई दूसरा लॉगिन किया हुआ अकाउंट न होने पर, सादा `ccam next` `no logged-in account other than <name>` के साथ रुकता है।
- `--continue` के साथ और फ़ोल्डर में कोई बातचीत न हो: `no account has a conversation in <folder> to continue`।
- जब संगठन वाला नियम कोई अकाउंट नहीं छोड़ता: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``।
- `--auto` के साथ और कोई लॉगिन किया हुआ अकाउंट न हो: `no logged-in account to pick`। जब हर अकाउंट भरा हो या उसका उपयोग अज्ञात हो: `` no account with known usage and headroom (see `ccam usage`) ``। दोनों exit code 75 के साथ ख़त्म होते हैं।
- `--max-usage`, `--strategy`, `--wait` और `--cross-org` केवल `--auto` के साथ काम करते हैं, और 1-100 से बाहर का मान या कोई दूसरी strategy अस्वीकार की जाती है।
- Fork बातचीत को कॉपी करता है, अकाउंट की दूसरी state को नहीं: मूल अकाउंट के `file-history/` के checkpoints तब तक साथ नहीं आते जब तक इतिहास साझा न हो।
- `--auto`, `--dry-run` के साथ भी नेटवर्क requests भेजता है। यह केवल पढ़ता है, और token फ़ाइलों में कभी नहीं लिखता।
- चुना गया अकाउंट उसी रास्ते से शुरू होता है जिससे `ccam <name>`, इसलिए वही जाँचें लागू होती हैं, उस चेतावनी सहित जब वही इतिहास साझा करने वाला कोई दूसरा Claude Code इस फ़ोल्डर में पहले से चल रहा हो।
- `last-used` केवल तभी लिखा जाता है जब Claude Code वास्तव में शुरू होता है। `--dry-run` इसे नहीं छूता, इसलिए dry run रोटेशन को आगे नहीं बढ़ाता।
- `--fallback` किसी प्रयास के ख़त्म होने तक standard output रोककर रखता है, इसलिए `--output-format stream-json` जैसा streaming output आख़िर में एक साथ आता है। सीमा तक पहुँचने से पहले पहले प्रयास ने जो फ़ाइलें बदलीं, वे बदली ही रहती हैं; अगला प्रयास prompt को फिर से शुरू करता है।

## संबंधित

- [उपयोग सीमाएँ](/features/usage)
- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [उपयोग सीमा के साथ काम करने का तरीक़ा](/guides/usage-limit-workflow)
