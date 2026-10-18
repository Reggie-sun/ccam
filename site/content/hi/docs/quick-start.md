---
title: क्विक स्टार्ट
description: ccam add से दो Claude Code अकाउंट जोड़ें, claude1, claude2 या ccs से उनके बीच स्विच करें, और चुनें कि वे क्या साझा करें।
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

हर अकाउंट को `ccam add` से जोड़ें, फिर Claude Code को `claude1`, `claude2` या `ccs <name>` से शुरू करें। ज़्यादातर लोगों को बस इतना ही चाहिए।

## 1. अपना पहला अकाउंट जोड़ें

```sh
ccam add
```

विज़र्ड एक नाम पूछता है, जैसे `personal` या `work`, और यह भी कि अकाउंट आपके दूसरे अकाउंट्स के साथ कौन-से आइटम साझा करे। फिर यह `claude auth login` चलाता है, जो साइन इन करने के लिए ब्राउज़र खोलता है। पूरा होने पर आपको यह दिखता है:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

फ़िलहाल ccam अपने संदेश कोरियाई में प्रिंट करता है; कमांड और विकल्प हर भाषा में एक जैसे हैं।

## 2. दूसरा अकाउंट जोड़ें

```sh
ccam add work
```

ब्राउज़र में दूसरे अकाउंट से साइन इन करें। यह `claude2` बन जाता है।

## 3. स्विच करें

इनमें से कोई भी कमांड Claude Code को work अकाउंट से शुरू करती है:

```sh
claude2
ccs work
ccs 2
```

अकाउंट के नाम के बाद का सब कुछ बिना बदलाव के Claude Code को जाता है, इसलिए आप वही विकल्प इस्तेमाल करते रहते हैं जो आप जानते हैं:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

सिर्फ़ `ccs` चलाने पर अकाउंट चुनने का मेन्यू खुलता है: arrow keys और Enter इस्तेमाल करें, या अकाउंट का नंबर दबाएँ।

## 4. चुनें कि क्या साझा करना है

डिफ़ॉल्ट रूप से MCP सर्वर, `settings.json`, ग्लोबल `CLAUDE.md`, और agents, commands, skills और plugins अकाउंट्स के बीच साझा होते हैं। बातचीत का इतिहास साझा नहीं होता। अगर आप अकाउंट्स के बीच बातचीत जारी रखना चाहते हैं, तो दोनों अकाउंट्स के लिए इसे चालू करें:

```sh
ccam share personal --history on
ccam share work --history on
```

कुछ भी बदलने से पहले ccam दिखाता है कि वह क्या move करेगा। विस्तार के लिए [साझा करना](/docs/sharing) देखें।

## 5. जब कोई अकाउंट अपनी सीमा तक पहुँच जाए

```sh
ccs usage
ccs next --continue
```

`ccs usage` हर अकाउंट का 5 घंटे और साप्ताहिक उपयोग दिखाता है। `ccs next --continue` वही बातचीत अगले अकाउंट में खोलता है। [उपयोग सीमाएँ](/docs/usage-limits) देखें।

## पहले से ~/.claude इस्तेमाल कर रहे हैं?

फिर से साइन इन करने की बजाय इसे रजिस्टर करें:

```sh
ccam import
```

फ़ोल्डर अपनी जगह पर रहता है, और सादा `claude` पहले की तरह काम करता रहता है।
