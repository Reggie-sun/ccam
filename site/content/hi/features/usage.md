---
title: उपयोग सीमाएँ
description: ccam usage हर अकाउंट की 5 घंटे और साप्ताहिक सीमाएँ एक साथ दिखाता है, और ccam statusline अकाउंट का नाम और उपयोग Claude Code में दिखाता है।
updated: 2026-10-20
order: 7
translation: machine
sourceHash: f08a2fdc6f0b
---

`ccam usage` दिखाता है कि हर अकाउंट ने 5 घंटे और साप्ताहिक सीमा का कितना हिस्सा इस्तेमाल किया है, और हर विंडो के रीसेट होने तक का समय भी। कोई लंबा काम शुरू करने से पहले, या जब Claude Code बताए कि कोई सीमा पास है, तब इसे इस्तेमाल करें, ताकि पता चले कि किस अकाउंट में गुंजाइश बची है। `ccam statusline`, Claude Code के लिए स्टेटस लाइन कमांड है: यह Claude Code के अंदर अकाउंट का नाम और उपयोग दिखाता है और `ccam usage` के लिए संख्याएँ सेव करता है।

## उपयोग

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

टर्मिनल पर table तुरंत हर पंक्ति में `fetching…` के साथ दिखती है, और हर पंक्ति उसके अकाउंट के जवाब देते ही फिर से बनती है। `As of` बताता है कि संख्याएँ कितनी पुरानी हैं; `(cached)` का मतलब है कि वे endpoint से नहीं, सेव की गई फ़ाइल से आई हैं। आख़िरी पंक्ति 100% से कम में सबसे कम 5 घंटे के उपयोग वाले अकाउंट का नाम बताती है, और केवल तभी दिखती है जब आपके पास एक से ज़्यादा अकाउंट हों। जिस अकाउंट में लॉगिन नहीं है, वह `not logged in → ccam login <name>` दिखाता है, और API key वाला अकाउंट `no limits (API)` दिखाता है।

`ccam statusline install` के बाद, Claude Code `[work] · Opus · 5h 24% · 7d 41%` जैसी पंक्ति दिखाता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--json` | रिपोर्ट को JSON के रूप में प्रिंट करें: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | बंद |
| `--watch[=N]` | Ctrl+C दबाने तक हर `N` सेकंड में फिर से बनाएँ। अकेले `--watch`, या 0 या उससे कम मान का मतलब 60 है | बंद |
| `-v`, `--verbose` | मॉडल के हिसाब से साप्ताहिक सीमाएँ और एक `source:` पंक्ति (`live`, `cache` या `statusline`) दिखाएँ | बंद |
| `--no-cache` | 60 सेकंड की cache अनदेखी करें और endpoint से फिर पूछें | बंद |
| `--source live\|cache` | `cache` कभी endpoint से नहीं पूछता और केवल सेव की गई संख्याएँ दिखाता है; `live` डिफ़ॉल्ट व्यवहार है | live |
| `--refresh` | एक्सपायर हुए access tokens renew करें, केवल उन अकाउंट्स के लिए जो चल नहीं रहे | बंद |
| `--dry-run` | स्वीकार किया जाता है: `ccam usage` अपनी cache के अलावा कुछ नहीं बदलता। `--refresh` के साथ अस्वीकार किया जाता है: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | बंद |

`--source` का कोई और मान `--source must be live or cache` के साथ विफल होता है। `ccam usage` कोई आर्ग्युमेंट नहीं लेता।

`ccam statusline`, standard input से Claude Code का स्टेटस लाइन JSON पढ़ता है। इसका एकमात्र विकल्प `--wrap <command>` है: एक मौजूदा स्टेटस लाइन कमांड, जिसके आउटपुट की पहली पंक्ति के आगे अकाउंट का नाम लगता है। `ccam statusline install` कोई विकल्प नहीं लेता। दोनों `--dry-run` को `` --dry-run is not supported by `ccam statusline`; nothing was changed `` के साथ अस्वीकार करते हैं।

## यह कैसे काम करता है

संख्याएँ दो जगहों से आती हैं: Anthropic का undocumented usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, जिसे हर अकाउंट के अपने लॉगिन से बुलाया जाता है, और वे मान जो Claude Code अपनी स्टेटस लाइन कमांड को देता है। यह सब `internal/usage` में है।

1. `cli.usageCmd`, `config.toml` लोड करता है और हर अकाउंट के लिए एक `usage.Account` बनाता है: उसका फ़ोल्डर `~/.local/share/ccam/profiles/<name>/`, उसका session फ़ोल्डर (यह देखने के लिए कि वह चल रहा है या नहीं) और यह कि वह अपनी जगह पर import किया गया था या नहीं।
2. `Fetcher.FetchAll` हर अकाउंट के लिए एक goroutine शुरू करता है। ये सभी एक ही `http.Client` साझा करते हैं, इसलिए endpoint से कनेक्शन दोबारा इस्तेमाल होते हैं। कनेक्ट होने में 2 s लग सकते हैं, और हर अकाउंट को कुल 3 s मिलते हैं। Client redirects का पालन नहीं करता, क्योंकि redirect token को कहीं और भेज देता, और जवाब का अधिकतम 1 MB पढ़ता है।
3. हर अकाउंट के पूरा होते ही उसकी पंक्ति फिर से बनती है। जो अकाउंट विफल होता है या समय पर जवाब नहीं देता, वह सेव की गई संख्याओं पर लौट आता है, इसलिए एक धीमा अकाउंट कभी दूसरों को नहीं रोकता।

हर अकाउंट endpoint के सामने अपने goroutine में चलता है; जिसका समय ख़त्म हो जाता है, वह इसके बजाय अपनी सेव की गई संख्याएँ दिखाता है।

<!-- motion:usage-fetch -->

`personal` और `work` समय पर जवाब देते हैं और live संख्याएँ दिखाते हैं। `side` को 3 s से ज़्यादा लगते हैं, इसलिए `fallback()` उसकी दोनों cache फ़ाइलें पढ़ता है और जो आख़िर में लिखी गई, उसे इस्तेमाल करता है। `usage.Table()` पंक्तियाँ बनाता है और `usage.Best()` आख़िरी पंक्ति में बताया जाने वाला अकाउंट चुनता है।

एक अकाउंट के लिए `Fetcher.Fetch` ये जाँचें इस क्रम में करता है:

<!-- motion:usage-cache -->

1. `--source cache` के साथ यह सीधे `fallback()` पर जाता है।
2. जब तक `--no-cache` सेट न हो, `~/.local/state/ccam/usage/<name>.json` में 60 सेकंड (`CacheTTL`) से कम पुराना live नतीजा जैसा है वैसा लौटाया जाता है। लगातार दो बार `ccam usage` चलाने पर हर अकाउंट के लिए एक ही request लगती है।
3. `readCredentials()`, `.credentials.json` (या macOS कीचेन) को बदले बिना पढ़ता है, और केवल `claudeAiOauth` key decode होती है। Subscription लॉगिन न होने का मतलब `not logged in` है, या `no limits (API)` जब `.claude.json` में API key हो।
4. एक्सपायर हुआ access token इस्तेमाल नहीं होता। `--refresh` के बिना, अकाउंट अपनी सेव की गई संख्याएँ expired के निशान के साथ दिखाता है।
5. वरना ccam endpoint को बुलाता है। 401 या 403 को एक्सपायर हुआ लॉगिन माना जाता है, कोई दूसरा status या न पढ़ा जा सकने वाला body त्रुटि माना जाता है, और दोनों सेव की गई संख्याओं पर लौटते हैं।
6. अच्छा जवाब मापे जाने के समय के साथ `<name>.json` में लिखा जाता है।

`fallback()`, `<name>.json` और `<name>.statusline.json` दोनों पढ़ता है और नई वाली रखता है। सेव की गई संख्याएँ दिखाने से पहले यह हर विंडो का रीसेट समय जाँचता है: जो विंडो संख्याएँ सेव होने के बाद रीसेट हो चुकी है, वह 0% दिखाती है, क्योंकि पुराना प्रतिशत अब सही नहीं है। इसीलिए ऊपर `personal` बिना रीसेट समय के `0%` दिखाता है।

Endpoint के दो response formats हैं। नए में `session` (5 घंटे), `weekly_all` और `weekly_scoped` (मॉडल के हिसाब से) entries वाला `limits` array होता है। पुराने में `five_hour`, `seven_day`, `seven_day_opus` और `seven_day_sonnet` ऑब्जेक्ट होते हैं। `ParseLive` दोनों पढ़ता है: `limits` array को प्राथमिकता मिलती है, और पुराने ऑब्जेक्ट वह भरते हैं जो वह कवर नहीं करता। बिना किसी सीमा वाला body बदले हुए format के रूप में बताया जाता है, जिसके बारे में `ccam doctor` आपको बताता है।

### Tokens renew करना

Refresh tokens बदलते रहते हैं: एक को renew करने पर नया refresh token मिलता है और पुराना अमान्य हो जाता है। अगर ccam उसी अकाउंट से Claude Code चलते समय लॉगिन renew करता, तो वह session साइन आउट हो जाता। इसलिए `ccam usage` tokens पढ़ता है और कभी नहीं लिखता, जब तक आप `--refresh` न दें। तब `Fetcher.refresh`:

1. `ccam import` से अपनी जगह पर रजिस्टर किए गए अकाउंट अस्वीकार करता है, क्योंकि सादा `claude` ccam की जानकारी के बिना वह लॉगिन इस्तेमाल कर सकता है;
2. macOS कीचेन में रखे लॉगिन अस्वीकार करता है;
3. एक lock फ़ाइल लेता है, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (वह variable सेट न हो तो `~/.local/state/ccam/locks/` के नीचे), ताकि दो ccam processes कभी एक ही token को दो बार renew न करें;
4. lock के तहत, अकाउंट का कोई session चल रहा हो तो मना करता है, और अगर किसी दूसरे process ने token पहले ही renew कर दिया हो तो request छोड़ देता है;
5. refresh token को token endpoint पर post करता है और नए tokens को `fsutil.AtomicWrite` से वापस `.credentials.json` में लिखता है, फ़ाइल की बाकी हर key रखते हुए।

### स्टेटस लाइन collector

Claude Code स्टेटस लाइन कमांड चलाता है और उसे pipe से एक JSON document देता है। ccam उससे मॉडल का नाम और `rate_limits` पढ़ता है:

<!-- motion:usage-statusline -->

1. `runStatusline` अकाउंट को `CCAM_PROFILE` से खोजता है, जिसे ccam Claude Code शुरू करते समय सेट करता है, या वरना `CLAUDE_CONFIG_DIR` का अकाउंट फ़ोल्डरों से मिलान करके।
2. `ParseStatusLine`, `rate_limits.five_hour` और `rate_limits.seven_day` पढ़ता है। API key वाले अकाउंट्स के लिए, या session के पहले response से पहले, JSON में rate limits नहीं होतीं; तब कुछ सेव नहीं होता।
3. `SaveStatusLine`, `~/.local/state/ccam/usage/<name>.statusline.json` लिखता है। जब endpoint जवाब नहीं देता, तो `ccam usage` इसे इस्तेमाल करता है।
4. ccam `[name] · model · 5h N% · 7d N%` प्रिंट करता है। `--wrap` के साथ, यह आपकी कमांड को उसी input के साथ `sh -c` से चलाता है और इसके बजाय `[name]` और उसके बाद कमांड की पहली पंक्ति प्रिंट करता है। Windows पर यह कमांड को उसी तरह चलाता है जैसे Claude Code: Git Bash इंस्टॉल हो तो उससे, वरना PowerShell से।
5. जब 5 घंटे या साप्ताहिक उपयोग 90% या उससे ज़्यादा हो और कोई दूसरा अकाउंट मौजूद हो, तो पंक्ति `→ ccam next --auto --continue` पर ख़त्म होती है, उदाहरण के लिए `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`।

Wrap की गई कमांड को 2 सेकंड मिलते हैं। उसके बाद उसे बंद कर दिया जाता है और पंक्ति में केवल अकाउंट का नाम दिखता है, ताकि कोई धीमी कमांड Claude Code की स्टेटस लाइन न रोक सके।

`ccam statusline install` हर अकाउंट की `settings.json` में `statusLine` लिखता है। जब किसी अकाउंट के पास पहले से स्टेटस लाइन कमांड हो, तो ccam उसे `--wrap '<old command>'` के रूप में रखता है। Windows पर ccam का पाथ forward slashes और home फ़ोल्डर के लिए `~` के साथ लिखा जाता है (`~/.local/bin/ccam.exe statusline`), जिसे Git Bash और PowerShell दोनों समझते हैं। ग्रुप के ज़रिए साझा की गई `settings.json` उसके लिंक से एक बार लिखी जाती है, इसलिए ग्रुप के हर अकाउंट को यह मिल जाता है। इसे फिर से चलाने पर, जो फ़ाइलें पहले से ccam इस्तेमाल करती हैं उनके लिए `(already set)` प्रिंट होता है।

## सुरक्षा और विशेष स्थितियाँ

- जब तक आप `--refresh` न दें, `ccam usage` कभी token फ़ाइल नहीं बदलता। जब कोई लॉगिन एक्सपायर हो जाए, तो अकाउंट को एक बार `ccam <name>` से चलाएँ, Claude Code उसे renew कर देगा।
- `--refresh` कभी चल रहे अकाउंट्स, import किए गए अकाउंट्स या कीचेन लॉगिन को नहीं छूता। अस्वीकार किया गया refresh सेव की गई संख्याएँ और कारण दिखाता है।
- Usage endpoint सार्वजनिक रूप से documented नहीं है और बदल सकता है। ऐसा होने पर भी आपको सेव की गई संख्याएँ दिखती हैं, और `ccam doctor` format में बदलाव बताता है।
- सेव की गई संख्याएँ और जवाब दोनों न होने पर, पंक्ति `lookup failed:` और कारण दिखाती है, उदाहरण के लिए `--source cache` के साथ `no saved usage`।
- Cache फ़ाइलों में केवल प्रतिशत और रीसेट समय होते हैं, कभी tokens नहीं।

## संबंधित

- [उपयोग सीमाएँ और अकाउंट स्विच करना](/docs/usage-limits)
- [अगले अकाउंट पर स्विच करना](/features/next)
- [उपयोग सीमा के साथ काम करने का तरीक़ा](/guides/usage-limit-workflow)
- [कमांड](/docs/commands)
