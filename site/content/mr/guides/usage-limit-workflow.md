---
title: Claude Code वापर मर्यादेला पोहोचल्यावर काय करायचे
description: रीसेटची वाट पाहा किंवा तेच संभाषण लगेच दुसऱ्या Claude Code खात्यावर पुढे चालू ठेवा. 5 तासांची आणि साप्ताहिक मर्यादा कशी चालते, टप्प्याटप्प्याने.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

एखादे खाते वापर मर्यादेला (याला अनेकदा rate limit म्हणतात) पोहोचले की तुमच्याकडे दोन पर्याय असतात: मर्यादा रीसेट होईपर्यंत वाट पाहा, किंवा तेच संभाषण आत्ताच तुमच्या दुसऱ्या खात्यावर पुढे चालू ठेवा. ccam सोबत दुसरा पर्याय म्हणजे एकच कमांड: `ccam next --auto --continue`.

## Claude Code च्या कोणत्या मर्यादा आहेत?

Pro आणि Max प्लॅनवर Claude Code ला [पाच तासांची सत्र मर्यादा आणि एक साप्ताहिक मर्यादा](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) असते. दोन्ही वेब आणि ॲप्समधील Claude सोबत शेअर होतात, आणि वापर दोन्हींमधून एकाच वेळी मोजला जातो. एखादी मर्यादा संपली की Claude Code "You've hit your session limit · resets 3:45pm" असा [संदेश दाखवते](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit). मॉडेल बदलून उपयोग होत नाही, कारण मर्यादा सर्व मॉडेल्सना लागू होतात.

## पर्याय 1: रीसेटची वाट पाहा

तुमच्याकडे एकच खाते असेल तर वाट पाहा. Claude Code च्या अलीकडच्या आवृत्त्या [उघड्या सत्रातच वाट पाहतात](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) आणि मर्यादा रीसेट झाल्यावर काम आपोआप पुढे चालू ठेवतात. Claude Code मध्ये `/usage` तुमच्या प्लॅनच्या मर्यादांपैकी किती वापरले ते दाखवते.

## पर्याय 2: आत्ताच दुसऱ्या खात्यावर पुढे चालू ठेवा

तुमच्याकडे दुसरे खाते असेल, तर तुम्ही त्याच संभाषणासह लगेच पुढे जाऊ शकता. एकदा सेट करा:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

आता काम करताना स्टेटस लाइन खाते आणि त्याचा 5 तासांचा वापर दाखवते, त्यामुळे मर्यादा अचानक येत नाही.

मर्यादा गाठली की Claude Code मधून बाहेर पडा (`/exit`), मग:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

सर्वात जास्त जागा असलेल्या खात्यात संभाषण उघडा:

```sh
ccam next --auto --continue
```

Claude Code `work` म्हणून सुरू होते, सध्याच्या फोल्डरमधील सर्वात अलीकडील संभाषणासह, आतापर्यंत झालेल्या सगळ्या बोलण्यासह. इतिहास शेअर केल्याने (वरील `ccam group set default --add history`) ते एकच संभाषण राहते; त्याशिवाय Claude Code संभाषण `work` मध्ये कॉपी करते, आणि ते पहिले उत्तर किती tokens पुन्हा वाचेल ते ccam सांगते.

Claude Code मध्ये [ccam plugin](/docs/install#claude-code-plugin) असेल, तर मर्यादेवर थांबलेले सत्र एक सूचना दाखवते आणि भरलेली window नोंदवते, त्यामुळे `ccam next --auto` ते खाते त्याच्या reset पूर्वी पुन्हा निवडत नाही, आणि `/ccam:switch` (किंवा `/ccam:switch work`) `ccam next --continue` साठी संभाषणाची नोंद ठेवते, त्यामुळे काहीही अंदाजावर सोडले जात नाही.

**Team आणि Enterprise सीट्स.** आपोआप होणारी निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n`) एकाच Claude Team किंवा Enterprise organization मधील एका सीटवरून दुसऱ्या सीटवरही जाते. Anthropic च्या मदत केंद्रात म्हटले आहे की "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लॅनमधील प्रत्येक सीट नावानिशी एका व्यक्तीला नेमून दिलेली असते, आणि फक्त तीच व्यक्ती ती वापरू शकते) आणि "Users can't route their work through another user's account to use that user's remaining usage" (दुसऱ्या वापरकर्त्याच्या उरलेल्या वापराचा लाभ घेण्यासाठी वापरकर्ते आपले काम त्याच्या खात्यामार्फत चालवू शकत नाहीत) ([Team प्लॅन](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लॅन](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) साठीही असेच). त्यामुळे सीट्सदरम्यान आपोआप होणारी निवड वापरल्यास Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा. वैयक्तिक प्लॅनबद्दल (Pro, Max) [FAQ मधील अटींविषयीची टीप](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) सांगते.

## ccam आपोआप खाते बदलते का?

नाही. तुम्ही काम करत असताना ccam चालू नसते: ते Claude Code सुरू करते आणि स्वतःची जागा त्याला देते. कधी बदलायचे ते तुम्ही ठरवता, आणि `ccam next --auto` सर्वात कमी 5 तासांचा वापर असलेले खाते निवडते आणि 100% वर असलेली, लॉगिन नसलेली किंवा दुसऱ्या organization ची खाती वगळते. `config.toml` मध्ये `relay = true` असेल, तर ccam त्याऐवजी सत्राच्या शेजारी थांबून राहते, त्यामुळे `/ccam:switch` सत्राचे पर्याय, permission mode आणि effort level कायम ठेवून आपोआप पुढे चालू ठेवते. काही इतर टूल्स त्याऐवजी चालू सत्राचे लॉगिन बदलतात; [ccam आणि इतर account switchers](/compare/account-switchers) हा फरक समजावते.

## लक्षात ठेवण्यासारख्या गोष्टी

- `ccam next --continue` तुम्ही ज्या फोल्डरमध्ये काम करत होता त्याच फोल्डरमध्ये चालवा; संभाषणे प्रत्येक प्रोजेक्ट फोल्डरनुसार साठवली जातात.
- एकाच संभाषणावर एकाच वेळी दोन खाती चालवणे टाळा. तसा प्रयत्न केल्यास ccam सावध करते.
- इतिहास शेअर केलेला नसेल, तर संभाषण दुसऱ्या खात्यात कॉपी होते आणि दोन्ही प्रती आपापल्या वाटेने जातात; `/rewind` ज्या checkpoints मधून code परत आणते ते पहिल्या खात्याकडेच राहतात.
- नवीन खात्यावरील पहिला संदेश नेहमीपेक्षा जास्त मर्यादा वापरू शकतो, कारण त्या खात्याने अजून संभाषण cache केलेले नसते.
- अशा प्रकारे खाती बदलण्याची परवानगी आहे की नाही, हे प्रत्येक खात्याच्या अटींवर अवलंबून आहे.

`ccam usage` आणि `ccam next` च्या प्रत्येक पर्यायासाठी [वापर मर्यादा](/docs/usage-limits) पाहा.
