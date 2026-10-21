---
title: वापर मर्यादा
description: ccam usage प्रत्येक खात्याच्या 5 तासांच्या आणि साप्ताहिक मर्यादा एकाच वेळी दाखवते, आणि ccam statusline खात्याचे नाव आणि वापर Claude Code मध्ये दाखवते.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

प्रत्येक खात्याने 5 तासांची मर्यादा, साप्ताहिक मर्यादा आणि प्रत्येक मॉडेलची स्वतःची साप्ताहिक मर्यादा (उदा. Fable) यांपैकी किती वापरले ते `ccam usage` दाखवते, प्रत्येक मर्यादेसाठी एक bar आणि ती रीसेट होईपर्यंतचा वेळ यासह. कोणत्या खात्यात जागा उरली आहे ते पाहण्यासाठी, एखादे मोठे काम सुरू करण्यापूर्वी किंवा Claude Code मर्यादा जवळ आल्याचे सांगेल तेव्हा ते वापरा. `ccam statusline` ही Claude Code साठी स्टेटस लाइन कमांड आहे: ती Claude Code मध्ये खात्याचे नाव आणि वापर दाखवते आणि `ccam usage` साठी आकडे साठवते.

## वापर

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

प्रत्येक खात्याला एक ओळ मिळते, ज्यात त्याचा क्रमांक, त्याचे नाव आणि त्याचे आकडे किती जुने आहेत ते असते; `(cached)` म्हणजे ते endpoint कडून नाही, तर साठवलेल्या फाइलमधून आले आहेत. त्याखाली, प्रत्येक मर्यादेसाठी एक ओळ 20 cells चा bar, टक्केवारी आणि window रीसेट होईपर्यंतचा वेळ दाखवते: `5-hour`, `Weekly` (सर्व मॉडेल्स), मग प्रत्येक मॉडेलची स्वतःची साप्ताहिक window, उदा. `Fable`. खात्याकडे नसलेली window `-` दाखवते. टक्केवारी 80% पासून पिवळी आणि 100% वर लाल होते. टर्मिनलवर प्रत्येक खाते लगेच `fetching…` म्हणून दिसते, आणि ते उत्तर देईल तशा त्याच्या ओळी काढल्या जातात. शेवटची ओळ [`ccam next --auto`](/features/next) त्याच नियमांनी जे खाते निवडेल त्याचे नाव देते (`auto_max_usage`, `auto_strategy`, चालू सत्रे, `cross_org` सेट नसेल तर फक्त शेवटी वापरलेल्या खात्याच्या organization ची खाती, आणि 100% वरची प्रत्येक window वगळली जाते, त्यामुळे ज्या खात्याचा आठवडा भरला आहे त्याचे नाव कधीच येत नाही), आणि तुमच्याकडे एकापेक्षा जास्त खाती असतील तेव्हाच दिसते. लॉग इन नसलेले खाते `not logged in → ccam login <name>` दाखवते, आणि API key खाते `no limits (API)` दाखवते.

`ccam statusline install` नंतर Claude Code `[work] · Opus · 5h 24% · 7d 41%` सारखी ओळ दाखवते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--json` | reports JSON म्हणून छापा: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | बंद |
| `--watch[=N]` | तुम्ही Ctrl+C दाबेपर्यंत दर `N` सेकंदांनी पुन्हा काढा. फक्त `--watch`, किंवा 0 किंवा त्यापेक्षा कमी मूल्य, म्हणजे 60 | बंद |
| `-v`, `--verbose` | प्रत्येक खात्याखाली `source:` ओळ जोडा: त्याचे आकडे कुठून आले (usage endpoint, Claude Code, स्टेटस लाइन किंवा साठवलेली फाइल) | बंद |
| `--no-cache` | 60 सेकंदांच्या cache कडे दुर्लक्ष करून endpoint ला पुन्हा विचारा | बंद |
| `--source live\|claude\|cache` | `live` खात्याच्या token ने endpoint ला विचारते; `claude` Claude Code कडून वापर पाहून घेते (खाली पाहा); `cache` फक्त साठवलेले आकडे दाखवते | `usage_source`, नाहीतर `live` |
| `--refresh` | कालबाह्य access tokens नूतनीकरण करा, फक्त चालू नसलेल्या खात्यांसाठी | बंद |
| `--dry-run` | स्वीकारला जातो: `ccam usage` स्वतःच्या cache शिवाय काहीही बदलत नाही. `--refresh` सोबत नाकारला जातो: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | बंद |

`--source` चे इतर कोणतेही मूल्य `--source must be live, claude or cache` देऊन अयशस्वी होते.

### Claude Code कडूनच वापर

`config.toml` मध्ये `usage_source = "claude"` असेल (किंवा `--source claude`), तर ccam लॉगिन token अजिबात वाचत नाही. प्रत्येक खात्यासाठी ते खात्याच्या फोल्डरमध्ये `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` सुरू करते आणि standard input वर Claude Code ची `get_usage` विनंती पाठवते (Claude Code 2.1.270 किंवा नंतरचे); Claude Code स्वतःच्या लॉगिनने वापर पाहते आणि 5 तासांच्या, साप्ताहिक आणि model windows सह उत्तर देते. `--safe-mode` खात्याचे hooks, MCP सर्व्हर आणि plugins सुरू होऊ देत नाही, आणि कोणतेही सत्र साठवले जात नाही. प्रत्येक खात्यासाठी Claude Code एकदा सुरू करावे लागते, सुमारे एक सेकंद, त्यामुळे हा source `live` पेक्षा हळू आहे; खात्यांना एकाच वेळी विचारले जाते, प्रत्येकाला 20 सेकंदांच्या आत, आणि 60 सेकंदांचा cache नेहमीप्रमाणे लागू होतो. Claude Code ही विनंती experimental म्हणून चिन्हांकित करते, त्यामुळे तिने उत्तर देणे थांबवले तर ccam साठवलेले आकडे दाखवते आणि `ccam doctor` ते नोंदवते. `--dry-run` Claude Code सुरू करत नाही (ते खात्याच्या फोल्डरमध्ये लिहिते) आणि साठवलेले आकडे दाखवते. `usage_source` `live` किंवा `claude` घेतो; `cache` फक्त `--source` चा पर्याय आहे. `ccam usage` कोणतेही आर्ग्युमेंट घेत नाही.

`ccam statusline` standard input वरून Claude Code चा स्टेटस लाइन JSON वाचते. त्याचा एकच पर्याय `--wrap <command>` आहे: आधीपासून असलेली स्टेटस लाइन कमांड, जिच्या आउटपुटच्या पहिल्या ओळीसमोर खात्याचे नाव लावले जाते. `ccam statusline install` कोणतेही पर्याय घेत नाही. दोन्ही `--dry-run` `` --dry-run is not supported by `ccam statusline`; nothing was changed `` देऊन नाकारतात.

## हे कसे काम करते

आकडे दोन ठिकाणांहून येतात: प्रत्येक खात्याच्या स्वतःच्या लॉगिनने बोलावलेला Anthropic चा दस्तऐवजीकृत नसलेला usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, आणि Claude Code आपल्या स्टेटस लाइन कमांडला देत असलेली मूल्ये. हे सर्व `internal/usage` मध्ये आहे.

1. `cli.usageCmd` `config.toml` लोड करते आणि प्रत्येक खात्यासाठी एक `usage.Account` तयार करते: त्याचा फोल्डर `~/.local/share/ccam/profiles/<name>/`, त्याचा session फोल्डर (ते चालू आहे का ते पाहण्यासाठी) आणि ते जागेवरच इम्पोर्ट केले आहे का.
2. `Fetcher.FetchAll` प्रत्येक खात्यासाठी एक goroutine सुरू करते. ते सर्व एकच `http.Client` वापरतात, त्यामुळे endpoint ची connections पुन्हा वापरली जातात. connect होण्यास 2 s लागू शकतात, आणि प्रत्येक खात्याला एकूण 3 s मिळतात. client redirects पाळत नाही, कारण redirect token दुसरीकडे पुन्हा पाठवेल, आणि उत्तराचा जास्तीत जास्त 1 MB वाचतो.
3. प्रत्येक खाते पूर्ण होईल तशी त्याची ओळ पुन्हा काढली जाते. अयशस्वी होणारे किंवा वेळेत उत्तर न देणारे खाते साठवलेल्या आकड्यांवर येते, त्यामुळे एक हळू खाते इतरांना कधीच अडवत नाही.

प्रत्येक खाते endpoint विरुद्ध आपल्या स्वतःच्या goroutine मध्ये चालते; वेळ संपलेले खाते त्याऐवजी आपले साठवलेले आकडे दाखवते.

<!-- motion:usage-fetch -->

`personal` आणि `work` वेळेत उत्तर देतात आणि live आकडे दाखवतात. `side` ला 3 s पेक्षा जास्त वेळ लागतो, म्हणून `fallback()` त्याच्या दोन cache फाइल्स वाचते आणि जी शेवटी लिहिली गेली ती वापरते. `usage.Table()` प्रत्येक खात्यासाठी एक block काढते आणि `usage.Best()` शेवटच्या ओळीत नाव दिलेले खाते निवडते.

एका खात्यासाठी `Fetcher.Fetch` या तपासण्या क्रमाने करते:

<!-- motion:usage-cache -->

1. `--source cache` सह ते थेट `fallback()` कडे जाते.
2. `--no-cache` सेट नसेल तर, `~/.local/state/ccam/usage/<name>.json` मधील 60 सेकंदांपेक्षा (`CacheTTL`) नवीन live परिणाम जसाच्या तसा परत दिला जातो. `ccam usage` लागोपाठ दोनदा चालवल्यास प्रत्येक खात्यासाठी एकच request लागते.
3. `readCredentials()` `.credentials.json` (किंवा macOS Keychain) न बदलता वाचते, आणि फक्त `claudeAiOauth` key decode होते. subscription लॉगिन नसल्यास `not logged in`, किंवा `.claude.json` मध्ये API key असल्यास `no limits (API)`.
4. कालबाह्य access token वापरले जात नाही. `--refresh` शिवाय, खाते कालबाह्य म्हणून चिन्हांकित केलेले साठवलेले आकडे दाखवते.
5. नाहीतर ccam endpoint बोलावते. 401 किंवा 403 कालबाह्य लॉगिन मानले जाते, इतर कोणताही status किंवा वाचता न येणारा body त्रुटी मानली जाते, आणि दोन्ही वेळा साठवलेले आकडे वापरले जातात. 429 म्हणजे endpoint ला खूप वेळा विचारले जात आहे: ccam पुन्हा विचारण्याच्या वेळेसह `<name>.pause` लिहिते (`Retry-After`, किंवा 5 मिनिटे, जास्तीत जास्त एक तास), `rate limited by the usage endpoint; asking again after 15:10` सह साठवलेले आकडे दाखवते, आणि तोपर्यंत त्या खात्यासाठी endpoint बोलावत नाही, `--no-cache` असले तरी. लगेच पुन्हा प्रयत्न केल्याने मर्यादा फक्त चालूच राहिली असती.
6. चांगले उत्तर मोजल्याच्या वेळेसह `<name>.json` मध्ये लिहिले जाते.

`fallback()` `<name>.json` आणि `<name>.statusline.json` दोन्ही वाचते आणि नवीन असलेले ठेवते. साठवलेले आकडे दाखवण्यापूर्वी ते प्रत्येक window ची reset time तपासते: आकडे साठवल्यानंतर रीसेट झालेली window 0% दाखवते, कारण जुनी टक्केवारी आता लागू होत नाही. म्हणूनच वरती `personal` reset time शिवाय `0%` दाखवते.

endpoint चे दोन response formats आहेत. नवीन format मध्ये `session` (5 तास), `weekly_all` आणि `weekly_scoped` (मॉडेलनुसार) नोंदी असलेला `limits` array असतो. जुन्यामध्ये `five_hour`, `seven_day`, `seven_day_opus` आणि `seven_day_sonnet` objects असतात. `ParseLive` दोन्ही वाचते: `limits` array जिंकतो, आणि तो जे व्यापत नाही ते जुने objects भरतात. कोणत्याही मर्यादा नसलेला body बदललेला format म्हणून नोंदवला जातो, ज्याबद्दल `ccam doctor` तुम्हाला सांगते.

### Tokens नूतनीकरण करणे

Refresh tokens बदलत राहतात: एक नूतनीकरण केल्यावर तुम्हाला नवीन refresh token मिळते आणि जुने अवैध होते. Claude Code त्याच खात्याने चालू असताना ccam ने लॉगिन नूतनीकरण केले, तर ते session साइन आउट होईल. म्हणून `ccam usage` tokens वाचते आणि कधीच लिहीत नाही, तुम्ही `--refresh` दिल्याशिवाय. तेव्हा `Fetcher.refresh`:

1. `ccam import` ने जागेवरच नोंदवलेली खाती नाकारते, कारण साधे `claude` ccam ला न कळता ते लॉगिन वापरत असू शकते;
2. macOS Keychain मध्ये साठवलेली लॉगिन्स नाकारते;
3. एक lock फाइल घेते, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (तो variable सेट नसताना `~/.local/state/ccam/locks/` खाली), त्यामुळे दोन ccam processes एकच token कधीच दोनदा नूतनीकरण करत नाहीत;
4. lock खाली, खात्याचे session चालू असेल तर नकार देते, आणि दुसऱ्या process ने token आधीच नूतनीकरण केले असेल तर request वगळते;
5. refresh token token endpoint ला post करते आणि नवीन tokens `fsutil.AtomicWrite` ने `.credentials.json` मध्ये परत लिहिते, फाइलमधील इतर प्रत्येक key ठेवून.

### स्टेटस लाइन collector

Claude Code स्टेटस लाइन कमांड चालवते आणि तिला pipe ने एक JSON document देते. ccam त्यातून मॉडेलचे नाव आणि `rate_limits` वाचते:

<!-- motion:usage-statusline -->

1. `runStatusline` `CCAM_PROFILE` वरून खाते शोधते, जो ccam Claude Code सुरू करताना सेट करते, नाहीतर `CLAUDE_CONFIG_DIR` खात्यांच्या फोल्डर्सशी जुळवून.
2. `ParseStatusLine` `rate_limits.five_hour` आणि `rate_limits.seven_day` वाचते. API key खात्यांसाठी, किंवा session च्या पहिल्या response आधी, JSON मध्ये rate limits नसतात; तेव्हा काहीही साठवले जात नाही.
3. `SaveStatusLine` `~/.local/state/ccam/usage/<name>.statusline.json` लिहिते. endpoint उत्तर देत नाही तेव्हा `ccam usage` ती वापरते.
4. ccam `[name] · model · 5h N% · 7d N%` छापते. `--wrap` सह, ते तुमची कमांड त्याच input सह `sh -c` मार्फत चालवते आणि त्याऐवजी `[name]` आणि त्यानंतर कमांडची पहिली ओळ छापते. Windows वर ते कमांड Claude Code प्रमाणेच चालवते: Git Bash इन्स्टॉल असल्यास त्याने, नाहीतर PowerShell ने.
5. 5 तासांचा किंवा साप्ताहिक वापर 90% किंवा त्याहून जास्त असेल आणि दुसरे खाते असेल, तर ओळीच्या शेवटी `→ ccam next --auto --continue` येते, उदाहरणार्थ `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

wrap केलेल्या कमांडला 2 सेकंद मिळतात. त्यानंतर ती थांबवली जाते आणि ओळीत फक्त खात्याचे नाव दिसते, त्यामुळे हळू कमांड Claude Code ची स्टेटस लाइन अडवू शकत नाही.

`ccam statusline install` प्रत्येक खात्याच्या `settings.json` मध्ये `statusLine` लिहिते. खात्याकडे आधीच स्टेटस लाइन कमांड असेल, तर ccam ती `--wrap '<old command>'` म्हणून ठेवते. Windows वर ccam चा पाथ forward slashes आणि home फोल्डरसाठी `~` वापरून लिहिला जातो (`~/.local/bin/ccam.exe statusline`), जो Git Bash आणि PowerShell दोन्ही समजतात. गटामार्फत शेअर केलेली `settings.json` तिच्या लिंकमार्फत एकदाच लिहिली जाते, त्यामुळे गटातील प्रत्येक खात्याला ती मिळते. ते पुन्हा चालवल्यास आधीच ccam वापरणाऱ्या फाइल्ससाठी `(already set)` छापले जाते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- तुम्ही `--refresh` दिल्याशिवाय `ccam usage` कधीच token फाइल बदलत नाही. लॉगिन कालबाह्य झाल्यावर, `ccam <name>` ने खाते एकदा चालवा, म्हणजे Claude Code ते नूतनीकरण करते.
- `--refresh` चालू खाती, इम्पोर्ट केलेली खाती किंवा Keychain लॉगिन्सना कधीच हात लावत नाही. नाकारलेले refresh साठवलेले आकडे आणि कारण दाखवते.
- usage endpoint सार्वजनिकरित्या दस्तऐवजीकृत नाही आणि बदलू शकतो. तो बदलला तरी तुम्हाला साठवलेले आकडे दिसतात, आणि `ccam doctor` format मधील बदल नोंदवते.
- साठवलेले आकडे आणि उत्तर दोन्ही नसल्यास, खात्याची ओळ `lookup failed:` आणि कारण दाखवते, उदाहरणार्थ `--source cache` सह `no saved usage`.
- cache फाइल्समध्ये फक्त टक्केवारी आणि reset times असतात, कधीच tokens नसतात.

## संबंधित

- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [पुढच्या खात्यावर जाणे](/features/next)
- [वापर मर्यादा workflow](/guides/usage-limit-workflow)
- [कमांड्स](/docs/commands)
