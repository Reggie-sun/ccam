---
title: डैशबोर्ड
description: टर्मिनल पर ccam अकाउंट्स, ग्रुप्स, फ़ोल्डर नियमों, एक्सटेंशन और जाँचों का एक full-screen डैशबोर्ड खोलता है; हर बदलाव पहले अपनी योजना दिखाता है।
updated: 2026-10-24
order: 15
translation: machine
sourceHash: "e48ca2c5a643"
---

बिना आर्ग्युमेंट के `ccam`, standard input और output के टर्मिनल होने पर एक full-screen डैशबोर्ड खोलता है; वरना यह `ccam list` प्रिंट करता है। डैशबोर्ड हर अकाउंट को उसके उपयोग और ग्रुप्स के साथ दिखाता है, और वही करता है जो कमांड करती हैं: शुरू करना, जोड़ना, लॉगिन करना, ग्रुप बदलना, फ़ोल्डर नियम सेट करना, जाँचना और मरम्मत करना। शुरू करने पर यह बंद हो जाता है; अकाउंट जोड़ने या लॉगिन करने के लिए यह टर्मिनल पर जाता है और फिर लौट आता है।

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

टैब के नामों के नीचे वाली पंक्ति एक वाक्य में बताती है कि मौजूदा टैब किसलिए है। ज़्यादा चौड़े टर्मिनल पर Accounts टैब हर अकाउंट का ईमेल भी दिखाता है।

हर बार खुलने पर डैशबोर्ड background में GitHub से पूछता है कि क्या ccam की कोई नई रिलीज़ है; पहला frame कभी इसका इंतज़ार नहीं करता, और offline होने पर कुछ नहीं दिखता। अगर कोई हो, तो शीर्षक के नीचे एक पंक्ति उसकी सिफ़ारिश करती है:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` पहले पूछता है, फिर उसी तरह अपडेट करता है जैसे [`ccam update`](/features/update) करता है: रिलीज़ को उसके SHA-256 checksums से जाँचा जाता है और वह एक ही कदम में ccam की जगह ले लेती है। जिस डैशबोर्ड में आप हैं वह अभी भी पुराना संस्करण है, इसलिए बाहर निकलें और `ccam` फिर से शुरू करें। जब अपडेट यहाँ नहीं हो सकता, जैसे `.deb` पैकेज से इंस्टॉल किए गए ccam के लिए, तो डैशबोर्ड वही सलाह दिखाता है जो `ccam update` देता है और प्रस्ताव बनाए रखता है।

पहला frame केवल `config.toml` से बनता है, 100 ms से काफ़ी कम में; हर अकाउंट के जवाब देते ही उपयोग भरता जाता है। हर सीमा की अपनी पट्टी होती है: 5 घंटे की विंडो, सभी मॉडल की साप्ताहिक विंडो, और Fable की अपनी साप्ताहिक विंडो; `-` का मतलब है कि अकाउंट के पास ऐसी कोई विंडो नहीं है। प्रतिशत 80% से पीला और 100% पर लाल हो जाता है, और विफल खोज के बाद cache से आए मान हल्के रंग में दिखते हैं। जिस विंडो को Claude Code ने भरा हुआ बताया, वह अपने reset तक 100% दिखाती है। एक्सपायर हुआ access token हल्के रंग में `⚠ token expired` दिखाता है (Claude Code अगली शुरुआत पर उसे renew करता है), जबकि `⚠ login expired` का मतलब है कि अकाउंट को `ccam login` चाहिए।

## टैब

| टैब | क्या दिखाता है | Keys |
| --- | --- | --- |
| Accounts | नंबर, नाम, ईमेल, 5 घंटे, साप्ताहिक और Fable की पट्टियाँ, आख़िरी उपयोग, ग्रुप; `✔` डिफ़ॉल्ट, `▸` इस फ़ोल्डर का अकाउंट, `•` चल रहा (ccam के बिना शुरू हुआ Claude Code भी: सादा `claude`, कोई IDE, background sessions और उनका daemon) | नीचे देखें |
| Groups | हर ग्रुप, वह क्या साझा करता है और उसके अकाउंट, और वे अकाउंट जो किसी ग्रुप में नहीं हैं | `enter` बदलें, `n` नया, `e` नाम बदलें, `x` हटाएँ |
| Folders | [`ccam folder`](/features/folder) के नियम और वह नियम जो यहाँ लागू होता है, तब भी जब आप उस फ़ोल्डर तक symbolic link से पहुँचे हों (नए नियम links को resolve करके सहेजे जाते हैं) | `n` यहाँ कोई अकाउंट इस्तेमाल करें, `enter` बदलें, `x` हटाएँ |
| Extensions | जो [`ccam ext`](/features/ext) दिखाता है, हर एक्सटेंशन एक बार, उन ग्रुप्स के साथ जिनसे वह आता है, `off:` वे अकाउंट जहाँ वह plugin चालू नहीं है, किसी अकाउंट के अपने सर्वरों में वे सर्वर जो उसके पास केवल इस फ़ोल्डर के लिए हैं, `This folder's project`, और वे समस्याएँ जो किसी अकाउंट की नहीं हैं | `r` फिर से लोड करें |
| Health | जो [`ccam doctor`](/features/doctor) बताता है, हर अकाउंट के लिए एक box के रूप में (नीचे देखें) | `f` चुने गए box की मरम्मत करें, `F` सबकी मरम्मत करें (`doctor --fix`), `↑`/`↓` box के बीच जाएँ |

`←`/`→` या `Tab` टैब बदलते हैं, `?` मौजूदा टैब की keys दिखाता है, `r` छोटी cache के बिना उपयोग फिर से पता करता है, और `q` या `Esc` बाहर निकलते हैं।

## Accounts टैब की keys

| Key | क्या करती है |
| --- | --- |
| `Enter` | चुने गए अकाउंट से Claude Code शुरू करें |
| `1`–`9` | उस नंबर वाले अकाउंट पर जाएँ; 12 के लिए जल्दी से `1` `2` टाइप करें। केवल `Enter` शुरू करता है |
| `↑`/`↓`, `k`/`j` | ऊपर-नीचे जाएँ |
| `n` | वह अकाउंट शुरू करें जिसे [`ccam next --auto`](/features/next) चुनता |
| `o` | विकल्पों के साथ शुरू करें: `--continue`, `--resume`, दूसरे आर्ग्युमेंट, या dry run |
| `u` | उपयोग का विवरण: 5 घंटे, साप्ताहिक और मॉडल की विंडो, उनके reset के साथ, और संख्याएँ कहाँ से आती हैं |
| `a` | अकाउंट जोड़ें (नीचे वाला विज़र्ड) |
| `l` | फिर से लॉगिन करें |
| `s` | अकाउंट के ग्रुप चुनें |
| `d` | इसे डिफ़ॉल्ट अकाउंट बनाएँ |
| `f` | इसे मौजूदा फ़ोल्डर में इस्तेमाल करें |
| `e` | इसका नाम बदलें |
| `c` | इसकी शॉर्टकट कमांड सेट करें या हटाएँ |
| `x` | इसे हटाएँ (पहले पूछता है, No से शुरू करते हुए) |

**Team और Enterprise सीटें।** अपने-आप होने वाला चुनाव (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n`) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है। Anthropic का सहायता केंद्र कहता है कि "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लान की हर सीट एक नामित व्यक्ति को सौंपी जाती है, और केवल वही व्यक्ति उसका इस्तेमाल कर सकता है) और "Users can't route their work through another user's account to use that user's remaining usage" (उपयोगकर्ता किसी दूसरे उपयोगकर्ता के बचे हुए उपयोग का इस्तेमाल करने के लिए अपना काम उसके अकाउंट के ज़रिए नहीं चला सकते) ([Team प्लान](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लान](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) के लिए भी इसी तरह)। इसलिए सीटों के बीच अपने-आप होने वाला चुनाव इस्तेमाल करना Anthropic की शर्तों का उल्लंघन हो सकता है। इस पर भरोसा करने से पहले अपने संगठन और Anthropic से पता कर लें। Personal प्लान (Pro, Max) के बारे में [FAQ का शर्तों वाला नोट](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) बताता है।

## Health टैब

Health टैब खुलते ही [`ccam doctor`](/features/doctor) की जाँचें चलाता है और वही दृश्य बनाता है: पहली पंक्ति में हर अकाउंट का एक निशान, फिर हर अकाउंट का एक box। box की हर पंक्ति एक जाँचा गया क्षेत्र है, पहले `Folder`, `Login` और `Share links`, जिन पर `✔` ठीक, `ℹ` जानकारी, `⚠` चेतावनी या `✘` त्रुटि का निशान होता है, ताकि आप तुरंत देख लें कि किस अकाउंट में समस्या है और कहाँ। जो नतीजे किसी अकाउंट से जुड़े नहीं हैं, वे `General` box में जाते हैं।

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` एक box से दूसरे पर ले जाते हैं, और चुने गए box का border accent रंग का होता है। `f` सिर्फ़ चुने गए box की मरम्मत करता है, उसी तरह जैसे `ccam doctor --fix` करता: किसी अकाउंट के box पर वह सिर्फ़ उसी अकाउंट की मरम्मत करता है, और `General` box पर उन नतीजों की जो किसी अकाउंट से जुड़े नहीं हैं (प्लगइन पाथ और शॉर्टकट); बाक़ी हर अकाउंट सिर्फ़ फिर से जाँचा जाता है। यह पहले पूछता है (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), और जिस box में मरम्मत के लिए कुछ नहीं है वह बस `nothing to repair in <name>` कहता है। `F` एक साथ सबकी मरम्मत करता है, `ccam doctor --fix` की तरह, और `r` फिर से जाँचता है।

## ग्रुप बदलना

Groups टैब पर `Enter` editor खोलता है: बाईं ओर आइटम, दाईं ओर अकाउंट, बदलने के लिए `Tab`, टॉगल करने के लिए `Space`। जो आइटम किसी अकाउंट को पहले से किसी दूसरे ग्रुप से मिलता है, वह कारण के साथ धुंधला दिखता है, और वह अकाउंट भी जिसे कोई आइटम दो बार मिलता; एक आइटम केवल एक ही ग्रुप से आ सकता है। `Enter` योजना दिखाता है:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

ग्रुप की कॉपी से अलग फ़ाइल के बारे में पहले पूछा जाता है (ग्रुप वाली रखें, मेरी इस्तेमाल करें, या merge करें)। फिर बदलाव background में उसी कोड से चलता है जो [`ccam group`](/features/share) का है: config पहले सेव होता है, हर अकाउंट का डेटा क्रम से move होता है, और मूल फ़ाइलें `.ccam-backup/` में जाती हैं; जो अकाउंट MCP sharing छोड़ता है, वह ग्रुप के सर्वर अपनी `.claude.json` में रखता है। जो बदलाव किसी चल रहे अकाउंट का डेटा move करता, उसे अकाउंट के नाम के साथ मना कर दिया जाता है, तब भी जब अकाउंट योजना खुली रहते शुरू हुआ हो, और जिस योजना को किसी दूसरे ccam के ग्रुप बदलाव ने पुराना कर दिया, वह `the groups changed while this plan waited; nothing was moved, make the change again` के साथ मना होती है। जब कोई दूसरा ccam ग्रुप बदल रहा हो, तो नया बदलाव, नाम बदलना या हटाना `another ccam is changing the groups right now` कहता है; बीच में रुके बदलाव के बाद, वे तब तक इंतज़ार करते हैं जब तक Health टैब पर `F` (या `ccam doctor --fix`) उसे पूरा न कर दे।

## अकाउंट जोड़ना

`a` वह विज़र्ड खोलता है जिसे `ccam add` भी टर्मिनल पर इस्तेमाल करता है:

1. अकाउंट का नाम।
2. यह कैसे साइन इन करता है: अभी Claude subscription (ब्राउज़र में), [API key](/features/api), या बाद में Claude subscription।
3. API key के लिए: ccam हर launch पर key कहाँ से पढ़ता है (`env:NAME` या `command:COMMAND`), endpoint, key कैसे भेजी जाती है और मॉडल।
4. शामिल होने के लिए ग्रुप, जिनमें `default` चुना हुआ होता है; जब तक कोई ग्रुप न हो, तब वह जो ग्रुप `default` साझा करता है।
5. एक शॉर्टकट कमांड, अगर आप चाहें।
6. पुष्टि के लिए एक सारांश।

`Esc` एक कदम पीछे जाता है, और पहले कदम पर विज़र्ड बंद कर देता है।

## अकाउंट हटाना

`x` कुछ भी होने से पहले पूछता है, No से शुरू करते हुए, और बताता है कि हटाना क्या करता है: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` API key अकाउंट के लिए revoke करने को कोई लॉगिन नहीं होता, और जब कोई दूसरा अकाउंट वही लॉगिन रखता हो, तो वह लॉग आउट नहीं किया जाता। import किया गया अकाउंट अपना फ़ोल्डर रखता है: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## रंग

रंग Claude Code के टर्मिनल UI की भूमिकाओं के हिसाब से होते हैं, इसलिए ccam उसके बगल में अपनापन लिए दिखता है; ग्रुप के नामों को उनके नाम से एक रंग मिलता है। ccam टर्मिनल से उसका background पूछता है और light या dark theme चुनता है। इसे ख़ुद `config.toml` में सबसे ऊपर, पहले `[[group]]`, `[[profile]]` या `[[folder]]` से ऊपर, सेट करें:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` themes लाल-हरे रंग की पहचान में कमी (color blindness) के लिए हैं, और `-ansi` themes टर्मिनल के अपने 16 रंग इस्तेमाल करती हैं। `NO_COLOR` रंग बंद कर देता है।

## सुरक्षा और विशेष स्थितियाँ

- डैशबोर्ड वही functions बुलाता है जो कमांड बुलाती हैं, इसलिए हर नियम लागू रहता है: डेटा move होने से पहले योजना, चल रहे अकाउंट्स के लिए मनाही, बैकअप, और लॉगिन tokens जो कभी किसी अकाउंट से बाहर नहीं जाते।
- डैशबोर्ड में dry run नहीं है: बिना नाम के `ccam --dry-run` प्रिंट करता है कि इस फ़ोल्डर का अकाउंट कैसे शुरू होगा, और किसी अकाउंट पर `o` सिर्फ़ उसी एक शुरुआत का dry run देता है।
- Claude Code शुरू होने या लॉगिन खुलने से पहले टर्मिनल अपनी सामान्य स्थिति में आ जाता है।
- किसी dialog के लिए बहुत छोटी स्क्रीन पहले टैब की पंक्तियाँ हटाती है; पंक्तियाँ wrap होने के बजाय `…` से काटी जाती हैं। छोटी ऊँचाई वाली स्क्रीन पर इसके बाद ख़ाली पंक्तियाँ और टैब का विवरण हटते हैं, ताकि keys वाली पंक्ति बनी रहे। 40 कॉलम या 8 पंक्तियों से कम पर डैशबोर्ड केवल `Too small: ccam needs 40x8 (now 36x7)` दिखाता है।
- जब कोई बदलाव, हटाना, अपडेट या मरम्मत लागू हो रही हो, तब `Ctrl+C` उसका इंतज़ार करता है: स्टेटस लाइन `ccam quits when it is done.` कहती है।
- `U` से ccam अपडेट होने के बाद, डैशबोर्ड के अकाउंट जोड़ने या लॉगिन से लौटने पर प्रस्ताव फिर से नहीं आता।

## संबंधित

- [ग्रुप के साथ साझा करना](/features/share)
- [Claude Code शुरू करना](/features/launch)
- [अकाउंट प्रबंधित करना](/features/accounts)
- [उपयोग सीमाएँ](/features/usage)
