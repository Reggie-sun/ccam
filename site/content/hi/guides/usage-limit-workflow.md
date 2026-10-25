---
title: Claude Code की usage limit पूरी होने पर क्या करें
description: रीसेट का इंतज़ार करें या वही बातचीत तुरंत दूसरे Claude Code अकाउंट पर जारी रखें। 5 घंटे और साप्ताहिक सीमा कैसे काम करती है, क़दम-दर-क़दम।
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "18ef61a68c55"
---

जब किसी अकाउंट की उपयोग सीमा (जिसे अक्सर rate limit भी कहते हैं) पूरी हो जाती है, तो आपके पास दो रास्ते हैं: सीमा रीसेट होने तक इंतज़ार करें, या वही बातचीत अभी अपने किसी दूसरे अकाउंट पर जारी रखें। ccam के साथ दूसरा रास्ता बस एक कमांड है: `ccam next --auto --continue`।

## Claude Code की कौन-सी सीमाएँ हैं?

Pro और Max प्लान पर Claude Code की [पाँच घंटे की session सीमा और एक साप्ताहिक सीमा](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) होती है। दोनों वेब और ऐप्स वाले Claude के साथ साझा होती हैं, और उपयोग दोनों में एक साथ गिना जाता है। जब कोई एक सीमा ख़त्म हो जाती है, तो Claude Code "You've hit your session limit · resets 3:45pm" जैसा [संदेश दिखाता है](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)। मॉडल बदलने से मदद नहीं मिलती, क्योंकि सीमाएँ सभी मॉडलों पर लागू होती हैं।

## विकल्प 1: रीसेट का इंतज़ार करें

अगर आपके पास एक ही अकाउंट है, तो इंतज़ार करें। Claude Code के हाल के संस्करण [खुले session में इंतज़ार करते हैं](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) और सीमा रीसेट होते ही काम अपने-आप जारी रखते हैं। Claude Code के अंदर `/usage` दिखाता है कि आपने अपने प्लान की सीमाओं का कितना हिस्सा इस्तेमाल किया है।

## विकल्प 2: अभी दूसरे अकाउंट पर जारी रखें

अगर आपके पास दूसरा अकाउंट है, तो आप उसी बातचीत के साथ तुरंत आगे बढ़ सकते हैं। इसे एक बार सेट करें:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

अब काम करते समय स्टेटस लाइन अकाउंट और उसका 5 घंटे का उपयोग दिखाती है, इसलिए सीमा अचानक नहीं आती।

जब सीमा पूरी हो जाए, तो Claude Code से बाहर निकलें (`/exit`), फिर:

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

सबसे ज़्यादा गुंजाइश वाले अकाउंट में बातचीत खोलें:

```sh
ccam next --auto --continue
```

Claude Code `work` के रूप में शुरू होता है, मौजूदा फ़ोल्डर की सबसे हाल की बातचीत के साथ, जिसमें अब तक कही गई हर बात शामिल है। इतिहास साझा करने से (ऊपर `ccam group set default --add history`) यह एक ही बातचीत रहती है; इसके बिना Claude Code बातचीत को `work` में कॉपी कर देता है, और ccam बताता है कि वह पहला जवाब कितने tokens फिर से पढ़ता है।

Claude Code में [ccam plugin](/docs/install#claude-code-plugin) होने पर, सीमा पर रुकने वाला session एक नोटिफ़िकेशन दिखाता है और भरी हुई विंडो दर्ज करता है, ताकि `ccam next --auto` उस अकाउंट को उसके reset से पहले फिर से न चुने, और `/ccam:switch` (या `/ccam:switch work`) बातचीत को `ccam next --continue` के लिए दर्ज कर देता है, ताकि कुछ भी अनुमान पर न छोड़ा जाए।

**Team और Enterprise सीटें।** अपने-आप होने वाला चुनाव (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n`, और चालू करने पर [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है। Anthropic का सहायता केंद्र कहता है कि "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लान की हर सीट एक नामित व्यक्ति को सौंपी जाती है, और केवल वही व्यक्ति उसका इस्तेमाल कर सकता है) और "Users can't route their work through another user's account to use that user's remaining usage" (उपयोगकर्ता किसी दूसरे उपयोगकर्ता के बचे हुए उपयोग का इस्तेमाल करने के लिए अपना काम उसके अकाउंट के ज़रिए नहीं चला सकते) ([Team प्लान](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लान](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) के लिए भी इसी तरह)। इसलिए सीटों के बीच अपने-आप होने वाला चुनाव इस्तेमाल करना Anthropic की शर्तों का उल्लंघन हो सकता है। वहाँ इसका इस्तेमाल करना है या नहीं, यह आपका फ़ैसला है: इस पर भरोसा करने से पहले अपने संगठन और Anthropic से पता कर लें। Personal प्लान (Pro, Max) के बारे में [FAQ का शर्तों वाला नोट](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) बताता है।

## क्या ccam अपने-आप अकाउंट बदलता है?

केवल तभी जब आप इसे चालू करें। डिफ़ॉल्ट रूप से आपके काम करते समय ccam चल नहीं रहा होता: यह Claude Code को शुरू करता है और ख़ुद को उससे बदल देता है, इसलिए स्विच कब करना है, यह आप तय करते हैं, और `ccam next --auto` सबसे कम 5 घंटे के उपयोग वाला अकाउंट चुनता है, उन अकाउंट्स को छोड़ देता है जो 100% पर हैं या लॉग इन नहीं हैं, और एक ही Team या Enterprise संगठन के भीतर, या personal प्लान के बीच, रहता है। `config.toml` में `relay = true` होने पर ccam इसके बजाय session के साथ-साथ इंतज़ार करता है, इसलिए `/ccam:switch` अपने-आप आगे बढ़ जाता है, session के विकल्पों, permission mode और effort level को बनाए रखते हुए। `auto_switch = true` होने पर (Linux और macOS; प्रायोगिक), ccam स्विच ख़ुद भी करता है: जब Claude Code session को उसकी सीमा पर रोकता है और उसी तरह के किसी दूसरे अकाउंट में गुंजाइश होती है, तो ccam session बंद करता है और बातचीत वहाँ जारी रखता है ([शर्तें और क़ीमत](/features/next#auto-switch))। कुछ दूसरे टूल इसकी बजाय चल रहे session का लॉगिन बदलते हैं; [ccam और दूसरे account switchers](/compare/account-switchers) इसका अंतर समझाता है।

## ध्यान रखने वाली बातें

- `ccam next --continue` उसी फ़ोल्डर में चलाएँ जिसमें आप काम कर रहे थे; बातचीत हर प्रोजेक्ट फ़ोल्डर के हिसाब से रखी जाती है।
- एक ही बातचीत पर एक ही समय में दो अकाउंट चलाने से बचें। अगर आप ऐसा करने की कोशिश करते हैं, तो ccam चेतावनी देता है।
- साझा इतिहास के बिना, बातचीत दूसरे अकाउंट में कॉपी होती है और दोनों कॉपियाँ अलग-अलग आगे बढ़ती हैं; जिन checkpoints से `/rewind` कोड वापस लाता है, वे पहले अकाउंट के पास ही रहते हैं।
- नए अकाउंट पर पहला संदेश उसकी सीमा का सामान्य से ज़्यादा हिस्सा ले सकता है, क्योंकि उस अकाउंट ने अभी तक बातचीत को cache नहीं किया है। `ccam next --continue --brief` इसके बजाय उसे एक छोटे सारांश से शुरू करता है।
- आप इस तरह अकाउंट्स के बीच स्विच कर सकते हैं या नहीं, यह हर अकाउंट की शर्तों पर निर्भर करता है।

`ccam usage` और `ccam next` के हर विकल्प के लिए [उपयोग सीमाएँ](/docs/usage-limits) देखें।
