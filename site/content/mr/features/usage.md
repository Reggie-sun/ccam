---
title: वापर मर्यादा
description: ccam usage प्रत्येक खात्याच्या 5 तासांच्या आणि साप्ताहिक मर्यादा एकाच वेळी दाखवते, आणि ccam statusline खात्याचे नाव आणि वापर Claude Code मध्ये दाखवते.
updated: 2026-10-24
order: 7
translation: machine
sourceHash: "a14b1da71a1a"
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
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

प्रत्येक खात्याला एक ओळ मिळते, ज्यात त्याचा क्रमांक, त्याचे नाव आणि त्याचे आकडे किती जुने आहेत ते असते; `(cached)` म्हणजे ते endpoint कडून नाही, तर साठवलेल्या फाइलमधून आले आहेत. त्याखाली, प्रत्येक मर्यादेसाठी एक ओळ 20 cells चा bar, टक्केवारी आणि window रीसेट होईपर्यंतचा वेळ दाखवते: `5-hour`, `Weekly` (सर्व मॉडेल्स), मग प्रत्येक मॉडेलची स्वतःची साप्ताहिक window, उदा. `Fable`. खात्याकडे नसलेली window `-` दाखवते. टक्केवारी 80% पासून पिवळी आणि 100% वर लाल होते. वापर मर्यादेवर request नाकारताना Claude Code ने भरलेली म्हणून नोंदवलेली window, साठवलेले आकडे काहीही सांगोत, तिच्या reset पर्यंत 100% दाखवते ([मर्यादेची नोंद](#the-limit-record) पाहा). टर्मिनलवर प्रत्येक खाते लगेच `fetching…` म्हणून दिसते, आणि ते उत्तर देईल तशा त्याच्या ओळी काढल्या जातात. शेवटची ओळ [`ccam next --auto`](/features/next) त्याच नियमांनी जे खाते निवडेल त्याचे नाव देते (`auto_max_usage`, `auto_strategy`, चालू सत्रे, `cross_org` सेट नसेल तर फक्त शेवटी वापरलेल्या खात्याच्या organization ची खाती, आणि 100% वरची प्रत्येक window वगळली जाते, त्यामुळे ज्या खात्याचा आठवडा भरला आहे त्याचे नाव कधीच येत नाही), आणि तुमच्याकडे एकापेक्षा जास्त खाती असतील तेव्हाच दिसते. डीफॉल्ट `headroom` strategy सह ती `→ Most headroom: <name>` अशी असते; `expiring` किंवा `balance` सह ती कारण देते, जसे `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team आणि Enterprise सीट्स.** आपोआप होणारी निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n`) एकाच Claude Team किंवा Enterprise organization मधील एका सीटवरून दुसऱ्या सीटवरही जाते. Anthropic च्या मदत केंद्रात म्हटले आहे की "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लॅनमधील प्रत्येक सीट नावानिशी एका व्यक्तीला नेमून दिलेली असते, आणि फक्त तीच व्यक्ती ती वापरू शकते) आणि "Users can't route their work through another user's account to use that user's remaining usage" (दुसऱ्या वापरकर्त्याच्या उरलेल्या वापराचा लाभ घेण्यासाठी वापरकर्ते आपले काम त्याच्या खात्यामार्फत चालवू शकत नाहीत) ([Team प्लॅन](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लॅन](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) साठीही असेच). त्यामुळे सीट्सदरम्यान आपोआप होणारी निवड वापरल्यास Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा. वैयक्तिक प्लॅनबद्दल (Pro, Max) [FAQ मधील अटींविषयीची टीप](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) सांगते.

लॉग इन नसलेले खाते `not logged in → ccam login <name>` दाखवते, आणि API key खाते `no limits (API)` दाखवते. कालबाह्य access token ``access token expired; `ccam <name>` renews it when it starts`` दाखवते, साठवलेले आकडे असल्यास त्यांच्या वर: पुढच्या सुरुवातीला Claude Code त्याचे नूतनीकरण करते. refresh token कालबाह्य झाल्यामुळे, किंवा सर्व्हरने नकार दिल्यानंतर Claude Code ने काढून टाकल्यामुळे, संपूर्णपणे कालबाह्य झालेले लॉगिन `` login expired → `ccam login <name>` `` दाखवते.

`ccam statusline install` नंतर Claude Code `[work] · Opus · 5h 24% · 7d 41%` सारखी ओळ दाखवते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--json` | reports JSON म्हणून छापा: `name`, `slot`, `status`, `limits` (खात्यावर usage credits चालू असल्यास `credits` सह), `source`, `at`, `error`, आणि मर्यादेची नोंद windows 100% करत असेल तेव्हा `limited` आणि `limited_windows`. एकही खाते नसल्यास ते `[]` छापते | बंद |
| `--watch[=N]` | तुम्ही Ctrl+C दाबेपर्यंत दर `N` सेकंदांनी पुन्हा काढा. फक्त `--watch`, किंवा 0 किंवा त्यापेक्षा कमी मूल्य, म्हणजे 60. मूल्य `=` नंतर येते: `--watch 10` एका सूचनेसह नाकारले जाते | बंद |
| `-v`, `--verbose` | प्रत्येक खात्याखाली `source:` ओळ जोडा: त्याचे आकडे कुठून आले (usage endpoint, Claude Code, स्टेटस लाइन, साठवलेली फाइल, किंवा Claude Code ने नोंदवलेली वापर मर्यादा). मर्यादेची नोंद आकड्यांवर मात करत असेल, तर त्यानंतर `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` सारखी ओळ येते | बंद |
| `--no-cache` | 60 सेकंदांच्या cache कडे दुर्लक्ष करून endpoint ला पुन्हा विचारा | बंद |
| `--source live\|claude\|cache` | `live` खात्याच्या token ने endpoint ला विचारते; `claude` Claude Code कडून वापर पाहून घेते (खाली पाहा); `cache` फक्त साठवलेले आकडे दाखवते | `usage_source`, नाहीतर `live` |
| `--refresh` | कालबाह्य access tokens नूतनीकरण करा, फक्त चालू नसलेल्या खात्यांसाठी. एका नूतनीकरणाला 30 सेकंदांपर्यंत लागू शकतात | बंद |
| `--dry-run` | स्वीकारला जातो: `ccam usage` स्वतःच्या cache शिवाय काहीही बदलत नाही. `--refresh` सोबत नाकारला जातो: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | बंद |

`--source` चे इतर कोणतेही मूल्य `--source must be live, claude or cache` देऊन अयशस्वी होते.

### Claude Code कडूनच वापर

`config.toml` च्या सुरुवातीला (पहिल्या `[[group]]`, `[[profile]]` किंवा `[[folder]]` च्या वर) `usage_source = "claude"` असेल, किंवा `--source claude` दिले असेल, तर ccam लॉगिन token अजिबात वाचत नाही. प्रत्येक खात्यासाठी ते खात्याच्या फोल्डरमध्ये `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` सुरू करते आणि standard input वर Claude Code ची `get_usage` विनंती पाठवते (Claude Code 2.1.270 किंवा नंतरचे); Claude Code स्वतःच्या लॉगिनने वापर पाहते आणि 5 तासांच्या, साप्ताहिक आणि model windows सह उत्तर देते. `--safe-mode` खात्याचे hooks, MCP सर्व्हर आणि plugins सुरू होऊ देत नाही, आणि कोणतेही सत्र साठवले जात नाही. प्रत्येक खात्यासाठी Claude Code एकदा सुरू करावे लागते, सुमारे एक सेकंद, त्यामुळे हा source `live` पेक्षा हळू आहे; खात्यांना एकाच वेळी विचारले जाते, प्रत्येकाला 20 सेकंदांच्या आत, आणि 60 सेकंदांचा cache नेहमीप्रमाणे लागू होतो. ज्या लॉगिनचे access token कालबाह्य झाले आहे, किंवा 5 मिनिटांत कालबाह्य होणार आहे, त्याला 50 सेकंद मिळतात, कारण Claude Code सुरू होताना ते नूतनीकरण करते आणि अर्ध्यावर तोडलेल्या नूतनीकरणामुळे लॉगिन गमावले जाऊ शकते. सोडून दिलेला शोध Claude Code ला SIGTERM ने संपायला सांगतो, त्यामुळे ते सुरू केलेले नूतनीकरण पूर्ण करते (Windows वर असा signal नसतो; तिथे ते थांबवले जाते). संपूर्णपणे कालबाह्य झालेले लॉगिन Claude Code सुरू न करता `` login expired → `ccam login <name>` `` दाखवते. Claude Code ही विनंती experimental म्हणून चिन्हांकित करते, त्यामुळे तिने उत्तर देणे थांबवले तर ccam साठवलेले आकडे दाखवते आणि `ccam doctor` ते नोंदवते. Claude Code ने आपला स्वतःचा शोध अयशस्वी झाल्याचे उत्तर दिले, तर ccam मर्यादा नसलेला plan न दाखवता, `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` सह साठवलेले आकडे दाखवते. `--dry-run` Claude Code सुरू करत नाही (ते खात्याच्या फोल्डरमध्ये लिहिते) आणि साठवलेले आकडे दाखवते. `usage_source` `live` किंवा `claude` घेतो; `cache` फक्त `--source` चा पर्याय आहे. `ccam usage` कोणतेही आर्ग्युमेंट घेत नाही.

`ccam statusline` standard input वरून Claude Code चा स्टेटस लाइन JSON वाचते. त्याचा एकच पर्याय `--wrap <command>` आहे: आधीपासून असलेली स्टेटस लाइन कमांड, जिच्या आउटपुटच्या पहिल्या ओळीसमोर खात्याचे नाव लावले जाते. `ccam statusline install` कोणतेही पर्याय घेत नाही. दोन्ही `--dry-run` `` --dry-run is not supported by `ccam statusline`; nothing was changed `` देऊन नाकारतात.

## हे कसे काम करते

आकडे दोन ठिकाणांहून येतात: प्रत्येक खात्याच्या स्वतःच्या लॉगिनने बोलावलेला Anthropic चा दस्तऐवजीकृत नसलेला usage endpoint, `GET https://api.anthropic.com/api/oauth/usage`, आणि Claude Code आपल्या स्टेटस लाइन कमांडला देत असलेली मूल्ये. हे सर्व `internal/usage` मध्ये आहे.

1. `cli.usageCmd` `config.toml` लोड करते आणि प्रत्येक खात्यासाठी एक `usage.Account` तयार करते: त्याचा फोल्डर `~/.local/share/ccam/profiles/<name>/`, त्याचे session markers (खात्याच्या फोल्डरमधील Claude Code च्या स्वतःच्या नोंदींसह, ते खाते चालू आहे का ते सांगतात) आणि ते जागेवरच इम्पोर्ट केले आहे का.
2. `Fetcher.FetchAll` प्रत्येक खात्यासाठी एक goroutine सुरू करते. ते सर्व एकच `http.Client` वापरतात, त्यामुळे endpoint ची connections पुन्हा वापरली जातात. connect होण्यास 2 s लागू शकतात, आणि प्रत्येक खात्याला एकूण 3 s मिळतात. client redirects पाळत नाही, कारण redirect token दुसरीकडे पुन्हा पाठवेल, आणि उत्तराचा जास्तीत जास्त 1 MB वाचतो.
3. प्रत्येक खाते पूर्ण होईल तशी त्याची ओळ पुन्हा काढली जाते. अयशस्वी होणारे किंवा वेळेत उत्तर न देणारे खाते साठवलेल्या आकड्यांवर येते, त्यामुळे एक हळू खाते इतरांना कधीच अडवत नाही.

प्रत्येक खाते endpoint विरुद्ध आपल्या स्वतःच्या goroutine मध्ये चालते; वेळ संपलेले खाते त्याऐवजी आपले साठवलेले आकडे दाखवते.

<!-- motion:usage-fetch -->

`personal` आणि `work` वेळेत उत्तर देतात आणि live आकडे दाखवतात. `side` ला 3 s पेक्षा जास्त वेळ लागतो, म्हणून `fallback()` त्याऐवजी त्याचे साठवलेले आकडे दाखवते. `usage.Table()` प्रत्येक खात्यासाठी एक block काढते आणि `usage.Best()` शेवटच्या ओळीत नाव दिलेले खाते निवडते.

एका खात्यासाठी `Fetcher.Fetch` या तपासण्या क्रमाने करते:

<!-- motion:usage-cache -->

1. `--source cache` सह ते थेट `fallback()` कडे जाते.
2. `--no-cache` सेट नसेल तर, `~/.local/state/ccam/usage/<name>.json` मधील 60 सेकंदांपेक्षा (`CacheTTL`) कमी वेळापूर्वी मोजलेला live परिणाम वापरला जातो, आत्ताच्या वेळेनुसार तपासून: त्यानंतर रीसेट झालेली window 0% दाखवते, आणि घड्याळापेक्षा पुढची तारीख असलेली फाइल ताजी मानली जात नाही. `ccam usage` लागोपाठ दोनदा चालवल्यास प्रत्येक खात्यासाठी एकच request लागते.
3. `readCredentials()` `.credentials.json` (किंवा macOS Keychain) न बदलता वाचते, आणि फक्त `claudeAiOauth` key decode होते. subscription लॉगिन नसल्यास `not logged in`, किंवा `.claude.json` मध्ये API key असल्यास `no limits (API)`. refresh token कालबाह्य झालेले, किंवा Claude Code ने काढून टाकलेले लॉगिन `` login expired → `ccam login <name>` `` असते आणि त्याचा शोध घेतला जात नाही.
4. कालबाह्य access token वापरले जात नाही. `--refresh` शिवाय, खाते ``access token expired; `ccam <name>` renews it when it starts`` सह आपले साठवलेले आकडे दाखवते.
5. नाहीतर ccam endpoint बोलावते. 401 किंवा 403 कालबाह्य लॉगिन मानले जाते (`login expired: authentication failed (401)`), इतर कोणताही status किंवा वाचता न येणारा body त्रुटी मानली जाते, आणि दोन्ही वेळा साठवलेले आकडे वापरले जातात. Ctrl+C ने संपलेला शोध `cancelled` सांगतो. 429 म्हणजे endpoint ला खूप वेळा विचारले जात आहे: ccam पुन्हा विचारण्याच्या वेळेसह `<name>.pause` लिहिते (`Retry-After`, किंवा 5 मिनिटे, जास्तीत जास्त एक तास), `rate limited by the usage endpoint; asking again after 15:10` सह साठवलेले आकडे दाखवते, आणि तोपर्यंत त्या खात्यासाठी endpoint बोलावत नाही, `--no-cache` असले तरी. लगेच पुन्हा प्रयत्न केल्याने मर्यादा फक्त चालूच राहिली असती.
6. चांगले उत्तर मोजल्याच्या वेळेसह `<name>.json` मध्ये लिहिले जाते.

`fallback()` सर्वात नवीन साठवलेला शोध घेते (`<name>.json`, किंवा Claude Code च्या शोधांमधील `<name>.claude.json`) आणि, `<name>.statusline.json` त्यानंतर लिहिली गेली असेल तर, त्यावर स्टेटस लाइनच्या 5 तासांच्या आणि साप्ताहिक windows ठेवते; स्टेटस लाइनकडे नसलेल्या model windows शोधाने साठवल्या तशाच राहतात. साठवलेले आकडे दाखवण्यापूर्वी ते प्रत्येक window ची reset time तपासते: आकडे साठवल्यानंतर रीसेट झालेली window 0% दाखवते, कारण जुनी टक्केवारी आता लागू होत नाही. म्हणूनच वरती `personal` reset time शिवाय `0%` दाखवते.

endpoint चे दोन response formats आहेत. नवीन format मध्ये `session` (5 तास), `weekly_all` आणि `weekly_scoped` (मॉडेलनुसार) नोंदी असलेला `limits` array असतो. जुन्यामध्ये `five_hour`, `seven_day`, `seven_day_opus` आणि `seven_day_sonnet` objects असतात. `ParseLive` दोन्ही वाचते: `limits` array जिंकतो, आणि तो जे व्यापत नाही ते जुने objects भरतात. कोणत्याही मर्यादा नसलेला body बदललेला format म्हणून नोंदवला जातो, ज्याबद्दल `ccam doctor` तुम्हाला सांगते.

### मर्यादेची नोंद

Claude Code एखादा turn वापर मर्यादेवर थांबवते, तेव्हा [ccam plugin](/docs/install#claude-code-plugin) चा hook संभाषणातून नाकारलेली request वाचतो: कोणती window भरली होती (5 तास, साप्ताहिक, Opus किंवा Sonnet) आणि ती कधी रीसेट होते. hook ही माहिती `~/.local/state/ccam/usage/<name>.limit.json` मध्ये लिहितो. `Fetch` ही नोंद प्रत्येक source मधून त्याला माहीत असलेल्या इतर सर्व गोष्टींवर ठेवते: तिच्यातील प्रत्येक window तिच्या reset पर्यंत 100% मोजली जाते, साठवलेले आकडे कितीही जुने किंवा कमी असले तरी, त्यामुळे `ccam next --auto`, `ccam pick` आणि `--fallback` तोपर्यंत ते खाते पुन्हा निवडत नाहीत. दुसरे काहीच माहीत नसेल, तर खाते फक्त हीच नोंद दाखवते (`-v` सह `source: the usage limit Claude Code reported`). कोणीही समोर बसलेले नसलेली सत्रे (`-p`, background सत्रे) कोणतीही नोंद लिहीत नाहीत, तसेच organization ने थांबवल्यासही (usage संपलेला, $0 वर सेट केलेला group) नोंद लिहिली जात नाही, कारण तिचे admins ते कधीही उठवू शकतात. दुसऱ्या Claude खात्यात `ccam login` केल्याने नोंद काढून टाकली जाते; फाइल हटवल्यानेही ती पुसली जाते.

### Tokens नूतनीकरण करणे

Refresh tokens बदलत राहतात: एक नूतनीकरण केल्यावर तुम्हाला नवीन refresh token मिळते आणि जुने अवैध होते. Claude Code त्याच खात्याने चालू असताना ccam ने लॉगिन नूतनीकरण केले, तर ते session साइन आउट होईल. म्हणून `ccam usage` tokens वाचते आणि कधीच लिहीत नाही, तुम्ही `--refresh` दिल्याशिवाय. तेव्हा `Fetcher.refresh`:

1. `ccam import` ने जागेवरच नोंदवलेली खाती नाकारते, कारण साधे `claude` ccam ला न कळता ते लॉगिन वापरत असू शकते;
2. macOS Keychain मध्ये साठवलेली लॉगिन्स नाकारते;
3. एक lock फाइल घेते, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, त्यामुळे दोन ccam processes एकच token कधीच दोनदा नूतनीकरण करत नाहीत; तो धरलेला असताना खाते `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)` दाखवते. `ccam rename`, `remove` आणि `logout` सुद्धा तो घेतात, आणि नूतनीकरण चालू असताना नकार देतात;
4. lock खाली, खात्याचे session चालू असेल तर नकार देते, ते कसेही सुरू केलेले असो (ccam मार्फत, साधे `claude`, IDE, किंवा background सत्र आणि त्यामागचा daemon, जो स्वतः लॉगिन नूतनीकरण करतो), आणि दुसऱ्या process ने token आधीच नूतनीकरण केले असेल तर request वगळते;
5. Claude Code खात्याच्या फोल्डरमध्ये स्वतःचे नूतनीकरणाचे locks (`.oauth_refresh.lock`, `<folder>.lock`) धरून असताना काहीही पाठवत नाही: तेव्हा Claude Code नेमके तेच लॉगिन नूतनीकरण करत असते;
6. refresh token token endpoint ला post करते आणि उत्तराची 30 सेकंदांपर्यंत वाट पाहते, जितका वेळ Claude Code स्वतःच्या नूतनीकरणांना देते. एकदा पाठवलेली request कधीही तोडली जात नाही, Ctrl+C ने सुद्धा नाही, कारण सर्व्हरने refresh token आधीच बदललेला असू शकतो आणि आता फक्त त्याच्या उत्तरातच लॉगिन असते: Ctrl+C किंवा hangup इतर शोध संपवतात, आणि उत्तर सेव्ह झाल्यावर ccam बाहेर पडते (Ctrl+C नंतर 130);
7. Claude Code च्या write lock (`.storage-write.lock`) खाली सेव्ह करते आणि त्या वेळी असलेल्या `.credentials.json` मधील फक्त token fields बदलते: दरम्यान Claude Code ने लिहिलेले (MCP साइन-इन) तसेच राहते, दरम्यान Claude Code ने नूतनीकरण केलेले लॉगिन ठेवले जाते, आणि लॉगआउट कधीही उलटवले जात नाही. उत्तरात नवीन refresh token चे आयुष्य दिलेले असेल, तर ते Claude Code प्रमाणेच `refreshTokenExpiresAt` म्हणून साठवले जाते.

पूर्ण न झालेले refresh, खात्याच्या ओळीवर किंवा त्याच्या साठवलेल्या आकड्यांच्या वर, त्याचे कारण सांगते: उदाहरणार्थ `token not renewed: a session is running`, किंवा request पाठवली गेली पण तिचे उत्तर वापरता किंवा सेव्ह करता आले नाही तेव्हा `login may be lost: …`, आणि त्यानंतर Claude Code ने लॉग इन करायला सांगितले तर `ccam login <name>` चालवण्याच्या सल्ल्यासह.

### स्टेटस लाइन collector

Claude Code स्टेटस लाइन कमांड चालवते आणि तिला pipe ने एक JSON document देते. ccam त्यातून मॉडेलचे नाव आणि `rate_limits` वाचते:

<!-- motion:usage-statusline -->

1. `runStatusline` `CCAM_PROFILE` वरून खाते शोधते, जो ccam Claude Code सुरू करताना सेट करते, नाहीतर `CLAUDE_CONFIG_DIR` खात्यांच्या फोल्डर्सशी जुळवून.
2. `ParseStatusLine` `rate_limits.five_hour` आणि `rate_limits.seven_day` वाचते. API key खात्यांसाठी, किंवा session च्या पहिल्या response आधी, JSON मध्ये rate limits नसतात; तेव्हा काहीही साठवले जात नाही.
3. `SaveStatusLine` `~/.local/state/ccam/usage/<name>.statusline.json` लिहिते. endpoint उत्तर देत नाही तेव्हा `ccam usage` ती वापरते.
4. ccam `[name] · model · 5h N% · 7d N%` छापते. `--wrap` सह, ते तुमची कमांड त्याच input सह `sh -c` मार्फत चालवते आणि Claude Code स्टेटस लाइन दाखवते तसा तिचा output छापते: रिकामी नसलेली प्रत्येक ओळ, पहिल्या ओळीसमोर `[name]` लावून. अयशस्वी होणारी कमांड फक्त खात्याचे नाव दाखवते. Windows वर ते कमांड Claude Code प्रमाणेच चालवते: Git Bash इन्स्टॉल असल्यास त्याने, नाहीतर PowerShell ने.
5. 5 तासांचा किंवा साप्ताहिक वापर 90% किंवा त्याहून जास्त असेल आणि दुसरे खाते असेल, तर ओळीच्या शेवटी `→ ccam next --auto --continue` येते, उदाहरणार्थ `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

wrap केलेल्या कमांडला स्वतःची वेळेची मर्यादा नसते: Claude Code स्टेटस लाइन कमांडला 10 मिनिटे देते आणि नवीन update सुरू झाल्यावर चालू run थांबवते, त्यामुळे हळू कमांड ccam शिवाय जशी दिसली असती तशीच दिसते. Claude Code ccam संपवते तेव्हा कमांडही संपते (Linux आणि macOS वर तिने सुरू केलेल्या प्रत्येक process सह).

`ccam statusline install` प्रत्येक खात्याच्या `settings.json` मध्ये `statusLine` लिहिते आणि त्या object च्या इतर keys, जसे `padding`, `refreshInterval` आणि `hideVimModeIndicator`, ठेवते. खात्याकडे आधीच स्टेटस लाइन कमांड असेल, तर ccam ती `--wrap '<old command>'` म्हणून ठेवते; Windows वर `--wrap-base64 <base64>` म्हणून, ज्याला कोणत्याही shell मध्ये quoting लागत नाही. आधीच एखादे ccam चालवणारी स्टेटस लाइन, दुसऱ्या पाथवरून असो किंवा जुन्या ccam ने लिहिल्याप्रमाणे quote केलेली असो, wrap करण्याऐवजी बदलली जाते, ती जे wrap करते ते ठेवून. Windows वर ccam चा पाथ forward slashes आणि home फोल्डरसाठी `~` वापरून लिहिला जातो (`~/.local/bin/ccam.exe statusline`), जो Git Bash आणि PowerShell दोन्ही समजतात. quotes लागणारा पाथ Claude Code वापरत असलेल्या shell साठी पूर्ण लिहिला जातो: Git Bash साठी `'…'`, PowerShell साठी `& '…'`. गटामार्फत शेअर केलेली `settings.json` तिच्या लिंकमार्फत एकदाच लिहिली जाते, त्यामुळे गटातील प्रत्येक खात्याला ती मिळते. ते पुन्हा चालवल्यास आधीच ccam वापरणाऱ्या फाइल्ससाठी `(already set)` छापले जाते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- तुम्ही `--refresh` दिल्याशिवाय `ccam usage` कधीच token फाइल बदलत नाही. access token कालबाह्य झाल्यावर, `ccam <name>` ने खाते एकदा चालवा, म्हणजे Claude Code ते नूतनीकरण करते. संपूर्णपणे कालबाह्य झालेल्या लॉगिनला `ccam login <name>` लागते.
- `--refresh` चालू खाती (Claude Code ccam शिवाय सुरू केले असले तरीही), इम्पोर्ट केलेली खाती किंवा Keychain लॉगिन्सना कधीच हात लावत नाही. नाकारलेले refresh साठवलेले आकडे आणि कारण दाखवते (`token not renewed: …`).
- usage endpoint सार्वजनिकरित्या दस्तऐवजीकृत नाही आणि बदलू शकतो. तो बदलला तरी तुम्हाला साठवलेले आकडे दिसतात, आणि `ccam doctor` format मधील बदल नोंदवते.
- साठवलेले आकडे आणि उत्तर दोन्ही नसल्यास, खात्याची ओळ `lookup failed:` आणि कारण दाखवते, उदाहरणार्थ `--source cache` सह `no saved usage`.
- cache फाइल्समध्ये फक्त टक्केवारी आणि reset times असतात, कधीच tokens नसतात.

## संबंधित

- [वापर मर्यादा आणि खाते बदलणे](/docs/usage-limits)
- [पुढच्या खात्यावर जाणे](/features/next)
- [वापर मर्यादा workflow](/guides/usage-limit-workflow)
- [कमांड्स](/docs/commands)
