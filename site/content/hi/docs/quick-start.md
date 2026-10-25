---
title: क्विक स्टार्ट
metaTitle: "क्विक स्टार्ट: Claude Code अकाउंट जोड़ें, स्विच करें"
description: "ccam add से दो Claude Code अकाउंट जोड़ें, उन्हें ccwork जैसे शॉर्टकट दें, उनके बीच स्विच करें, और चुनें कि वे क्या साझा करें।"
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

हर अकाउंट को `ccam add` से जोड़ें, फिर Claude Code को `ccam <name>` से या अपने रखे नाम वाले शॉर्टकट, जैसे `ccwork`, से शुरू करें। ज़्यादातर लोगों को बस इतना ही चाहिए।

## 1. अपना पहला अकाउंट जोड़ें

```sh
ccam add
```

विज़र्ड एक नाम पूछता है, जैसे `personal` या `work`, यह भी कि अकाउंट आपके दूसरे अकाउंट्स के साथ कौन-से आइटम साझा करे, और क्या शॉर्टकट कमांड बनाना है (चरण 4 देखें)। फिर यह `claude auth login` चलाता है, जो साइन इन करने के लिए ब्राउज़र खोलता है। रिमोट मशीन पर ब्राउज़र आख़िर में ऐसे `http://localhost:…/callback` पेज पर पहुँच सकता है जो लोड नहीं होता; वह पता टर्मिनल में पेस्ट करें और ccam साइन-इन पूरा कर देगा। पूरा होने पर आपको यह दिखता है:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. अपनी मौजूदा सेटिंग्स साथ लाएँ

अगर आपने पहले Claude Code इस्तेमाल किया है, तो आपकी सेटिंग्स, `CLAUDE.md`, एक्सटेंशन और MCP सर्वर `~/.claude` में हैं। उन्हें साझा सेटिंग्स में कॉपी करें, ताकि हर अकाउंट को वे मिल जाएँ:

```sh
ccam migrate
```

ccam पहले दिखाता है कि वह क्या कॉपी करेगा, और फिर पूछता है। यह `~/.claude` को केवल पढ़ता है, कभी कोई लॉगिन कॉपी नहीं करता, और जब तक आप `--history` न जोड़ें, बातचीत का इतिहास छोड़ देता है। विस्तार के लिए [मौजूदा सेटअप कॉपी करना](/docs/sharing#copy-an-existing-setup) देखें।

## 3. दूसरा अकाउंट जोड़ें

```sh
ccam add work
```

ब्राउज़र में दूसरे अकाउंट से साइन इन करें।

## 4. अकाउंट स्विच करें

`ccam work` Claude Code को work अकाउंट से शुरू करता है। इससे छोटा कुछ चाहिए, तो अकाउंट को उसकी अपनी शॉर्टकट कमांड दें, जैसे `ccwork`। साझा आइटम के बाद विज़र्ड पूछता है:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` दबाएँ, फिर कोई नाम टाइप करें, या सुझाया गया नाम इस्तेमाल करने के लिए Enter दबाएँ:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

लॉगिन के बाद ccam इसकी पुष्टि करता है:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

अगर आपने यह सवाल छोड़ दिया था, तो कभी भी शॉर्टकट जोड़ सकते हैं:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

अब इनमें से कोई भी कमांड Claude Code को work अकाउंट से शुरू करती है:

```sh
ccwork
ccam work
ccam 2
```

सिर्फ़ `ccam` चलाने पर [डैशबोर्ड](/features/dashboard) खुलता है: arrow keys या अकाउंट के नंबर से किसी अकाउंट पर जाएँ, फिर Enter दबाएँ।

शॉर्टकट या अकाउंट के नाम के बाद का सब कुछ बिना बदलाव के Claude Code को जाता है, इसलिए आप वही विकल्प इस्तेमाल करते रहते हैं जो आप जानते हैं:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. चुनें कि क्या साझा करना है

डिफ़ॉल्ट रूप से MCP सर्वर, `settings.json`, ग्लोबल `CLAUDE.md`, और agents, commands, skills और plugins ग्रुप `default` के ज़रिए अकाउंट्स के बीच साझा होते हैं। बातचीत का इतिहास साझा नहीं होता। अगर आप अकाउंट्स के बीच बातचीत जारी रखना चाहते हैं, तो इसे ग्रुप में जोड़ें:

```sh
ccam group set default --add history
```

कुछ भी बदलने से पहले ccam हर अकाउंट के लिए दिखाता है कि वह क्या move करेगा। विस्तार के लिए [ग्रुप के साथ साझा करना](/features/share) देखें।

## 6. जब कोई अकाउंट अपनी सीमा तक पहुँच जाए

```sh
ccam usage
ccam next --continue
```

`ccam usage` हर अकाउंट का 5 घंटे और साप्ताहिक उपयोग दिखाता है। `ccam next --continue` वही बातचीत अगले अकाउंट में खोलता है। Scripts में, `ccam --fallback work -p "…"` अपने-आप स्विच करता है। [उपयोग सीमाएँ](/docs/usage-limits) देखें।

**Team और Enterprise सीटें।** अपने-आप होने वाला चुनाव (`ccam next`, `--fallback`, `ccam pick`, डैशबोर्ड का `n`, और चालू करने पर [`auto_switch`](/features/next#auto-switch)) एक ही Claude Team या Enterprise संगठन की सीटों के बीच भी स्विच करता है। Anthropic का सहायता केंद्र कहता है कि "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लान की हर सीट एक नामित व्यक्ति को सौंपी जाती है, और केवल वही व्यक्ति उसका इस्तेमाल कर सकता है) और "Users can't route their work through another user's account to use that user's remaining usage" (उपयोगकर्ता किसी दूसरे उपयोगकर्ता के बचे हुए उपयोग का इस्तेमाल करने के लिए अपना काम उसके अकाउंट के ज़रिए नहीं चला सकते) ([Team प्लान](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लान](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) के लिए भी इसी तरह)। इसलिए सीटों के बीच अपने-आप होने वाला चुनाव इस्तेमाल करना Anthropic की शर्तों का उल्लंघन हो सकता है। वहाँ इसका इस्तेमाल करना है या नहीं, यह आपका फ़ैसला है: इस पर भरोसा करने से पहले अपने संगठन और Anthropic से पता कर लें। Personal प्लान (Pro, Max) के बारे में [FAQ का शर्तों वाला नोट](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) बताता है।

## पहले से ~/.claude इस्तेमाल कर रहे हैं?

फिर से साइन इन करने की बजाय इसे रजिस्टर करें:

```sh
ccam import
```

फ़ोल्डर अपनी जगह पर रहता है, और सादा `claude` पहले की तरह काम करता रहता है।
