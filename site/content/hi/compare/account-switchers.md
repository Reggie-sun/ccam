---
title: ccam और दूसरे Claude Code account switchers
description: Claude Code अकाउंट स्विच करने के लिए ccam की तुलना claude-swap, CCS, clauth, CAAM, cux और दूसरे टूल्स से, और कौन-सा तरीका किस काम के लिए ठीक है।
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Claude Code अकाउंट स्विच करने वाले टूल दो में से एक तरीका अपनाते हैं। ज़्यादातर टूल एक ही `~/.claude` फ़ोल्डर के अंदर लॉगिन बदलते हैं, इसलिए एक समय में एक ही अकाउंट सक्रिय रहता है। दूसरे टूल, जिनमें ccam भी है, `CLAUDE_CONFIG_DIR` से हर अकाउंट को उसका अपना फ़ोल्डर देते हैं, इसलिए हर अकाउंट लॉग इन रहता है और कई अकाउंट एक साथ चल सकते हैं। ccam इसमें साझा सेटिंग्स और MCP सर्वर जोड़ता है, और किसी बातचीत को दूसरे अकाउंट के फ़ोल्डर में जारी रखता है।

## एक लॉगिन बदलना या अलग-अलग फ़ोल्डर?

**लॉगिन बदलने** में सेटिंग्स, MCP सर्वर और इतिहास का एक ही सेट रहता है, क्योंकि फ़ोल्डर एक ही है। स्विच करना तेज़ है, और चल रहा session नया लॉगिन ले सकता है। लेकिन एक समय में एक ही अकाउंट सक्रिय रहता है, सभी terminals और editors के लिए एक साथ।

**अलग-अलग फ़ोल्डर** वह तरीका है जो [Claude Code के दस्तावेज़ में बताया गया है](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)। हर अकाउंट लॉग इन रहता है, और आप work और personal session साथ-साथ चला सकते हैं। इसकी क़ीमत यह है कि हर फ़ोल्डर अपनी अलग सेटिंग्स, MCP सर्वर और इतिहास के साथ शुरू होता है, जब तक कोई टूल उन्हें साझा न करे।

## टूल्स की तुलना

2026-10-19 को हर प्रोजेक्ट के अपने दस्तावेज़ से लिया गया। "दस्तावेज़ में नहीं" का मतलब है कि प्रोजेक्ट इसके बारे में कुछ नहीं कहता, यह नहीं कि सुविधा मौजूद नहीं है।

| टूल | तरीका | एक साथ कई अकाउंट | 5 घंटे और साप्ताहिक उपयोग | सीमा पूरी होने पर | प्लेटफ़ॉर्म |
| --- | --- | --- | --- | --- | --- |
| ccam | अलग-अलग फ़ोल्डर | हाँ | हाँ, `ccam usage` | आप `ccam next --auto --continue` चलाते हैं, और बातचीत सबसे ज़्यादा गुंजाइश वाले अकाउंट में जारी रहती है | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | लॉगिन बदलता है; `cswap run` से अलग-अलग फ़ोल्डर | हाँ, `cswap run` से | हाँ | `cswap auto` एक सीमा-बिंदु पर लॉगिन बदलता है; चल रहा session जारी रहता है | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | अलग-अलग फ़ोल्डर | हाँ | आंशिक रूप से, macOS menu bar ऐप में | Account pools एक local proxy के ज़रिए requests अगले अकाउंट को भेजते हैं | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | लॉगिन बदलता है; `clauth start` से अलग-अलग फ़ोल्डर | हाँ, `clauth start` से | हाँ | fallback chain में अगले अकाउंट पर स्विच करता है | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | लॉगिन बदलता है, या isolated profiles | हाँ, isolated profiles से | हाँ | `caam run` non-interactive कमांड को दूसरे अकाउंट पर फिर से चलाता है | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | लॉगिन बदलता है और `claude` को wrap करता है | दस्तावेज़ में नहीं | हाँ | लॉगिन बदलता है और वही बातचीत फिर से शुरू करता है | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | लॉगिन बदलता है | नहीं | हाँ | दस्तावेज़ में नहीं | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | लॉगिन बदलता है | नहीं | हाँ | 100% पर वैकल्पिक स्विच | macOS |

ccswap, claude-swap का एक fork है जो Codex लॉगिन भी संभालता है। इनमें से कई टूल Codex या Gemini जैसे दूसरे coding agents भी संभालते हैं।

## ccam किस पर ध्यान देता है

- हर अकाउंट लॉग इन रहता है। स्विच करने से कोई अकाउंट लॉग आउट नहीं होता, और `ccam work` और `ccam personal` दो terminals में एक साथ चल सकते हैं।
- MCP सर्वर, `settings.json`, `CLAUDE.md`, agents, skills और plugins फ़ोल्डरों के बीच links से साझा होते हैं, हर आइटम और हर अकाउंट के हिसाब से। लॉगिन कभी साझा नहीं होते।
- साझा इतिहास के साथ, एक कमांड, `ccam next --auto --continue`, वही बातचीत सबसे ज़्यादा गुंजाइश वाले अकाउंट के फ़ोल्डर में खोलती है।
- ccam Claude Code को शुरू करता है और ख़ुद को उससे बदल देता है, इसलिए background में कुछ नहीं चलता। यह आपके उपयोग पर नज़र नहीं रखता और अपने-आप स्विच नहीं करता; कब स्विच करना है, यह आप तय करते हैं।

## जहाँ दूसरे टूल ज़्यादा करते हैं

- claude-swap, CCS, clauth, CAAM और cux, WSL के बिना Windows पर चलते हैं। ccam नहीं चलता।
- claude-swap, clauth और cux, किसी अकाउंट के उपयोग की तय सीमा पार होने पर अपने-आप स्विच करते हैं। ccam स्विच करना आप पर छोड़ता है।
- CCSwitcher, Claude Switcher और claude-swap उपयोग को macOS menu bar में दिखाते हैं।
- CCS, CAAM, ccswap और Claude Switcher, Codex और दूसरे coding agents भी संभालते हैं।

## कौन-सा ठीक रहेगा?

अगर आप एक समय में एक ही अकाउंट इस्तेमाल करते हैं, चाहते हैं कि स्विच अपने-आप हो, या WSL के बिना Windows पर हैं, तो लॉगिन बदलने वाला टूल चुनें। अगर आप अकाउंट्स को साथ-साथ चलाना, एक ही साझा सेटअप रखना और अगले अकाउंट में वही बातचीत जारी रखना चाहते हैं, तो ccam चुनें। अगर आपको कुछ भी साझा नहीं करना है, तो [साधारण shell aliases](/compare/aliases) काफ़ी हैं।

कई अकाउंट इस्तेमाल करने की अनुमति है या नहीं, यह हर अकाउंट की शर्तों पर निर्भर करता है। ऐसा करने से पहले उन्हें जाँच लें।
