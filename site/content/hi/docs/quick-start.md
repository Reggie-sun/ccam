---
title: क्विक स्टार्ट
description: "ccam add से दो Claude Code अकाउंट जोड़ें, उन्हें ccwork जैसे शॉर्टकट दें, उनके बीच स्विच करें, और चुनें कि वे क्या साझा करें।"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: eb078efa2ee2
---

हर अकाउंट को `ccam add` से जोड़ें, फिर Claude Code को `ccam <name>` से या अपने रखे नाम वाले शॉर्टकट, जैसे `ccwork`, से शुरू करें। ज़्यादातर लोगों को बस इतना ही चाहिए।

## 1. अपना पहला अकाउंट जोड़ें

```sh
ccam add
```

विज़र्ड एक नाम पूछता है, जैसे `personal` या `work`, यह भी कि अकाउंट आपके दूसरे अकाउंट्स के साथ कौन-से आइटम साझा करे, और क्या शॉर्टकट कमांड बनाना है (चरण 4 देखें)। फिर यह `claude auth login` चलाता है, जो साइन इन करने के लिए ब्राउज़र खोलता है। पूरा होने पर आपको यह दिखता है:

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

सिर्फ़ `ccam` चलाने पर अकाउंट चुनने का मेन्यू खुलता है: arrow keys और Enter इस्तेमाल करें, या अकाउंट का नंबर दबाएँ।

शॉर्टकट या अकाउंट के नाम के बाद का सब कुछ बिना बदलाव के Claude Code को जाता है, इसलिए आप वही विकल्प इस्तेमाल करते रहते हैं जो आप जानते हैं:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. चुनें कि क्या साझा करना है

डिफ़ॉल्ट रूप से MCP सर्वर, `settings.json`, ग्लोबल `CLAUDE.md`, और agents, commands, skills और plugins अकाउंट्स के बीच साझा होते हैं। बातचीत का इतिहास साझा नहीं होता। अगर आप अकाउंट्स के बीच बातचीत जारी रखना चाहते हैं, तो दोनों अकाउंट्स के लिए इसे चालू करें:

```sh
ccam share personal --history on
ccam share work --history on
```

कुछ भी बदलने से पहले ccam दिखाता है कि वह क्या move करेगा। विस्तार के लिए [साझा करना](/docs/sharing) देखें।

## 6. जब कोई अकाउंट अपनी सीमा तक पहुँच जाए

```sh
ccam usage
ccam next --continue
```

`ccam usage` हर अकाउंट का 5 घंटे और साप्ताहिक उपयोग दिखाता है। `ccam next --continue` वही बातचीत अगले अकाउंट में खोलता है। [उपयोग सीमाएँ](/docs/usage-limits) देखें।

## पहले से ~/.claude इस्तेमाल कर रहे हैं?

फिर से साइन इन करने की बजाय इसे रजिस्टर करें:

```sh
ccam import
```

फ़ोल्डर अपनी जगह पर रहता है, और सादा `claude` पहले की तरह काम करता रहता है।
