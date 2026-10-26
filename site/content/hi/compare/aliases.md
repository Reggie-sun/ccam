---
title: ccam और shell aliases की तुलना
description: कई Claude Code अकाउंट चलाने के लिए CLAUDE_CONFIG_DIR वाले shell aliases और ccam की हर सुविधा के हिसाब से तुलना।
updated: 2026-10-26
order: 1
translation: machine
sourceHash: c30d84bceba0
---

हर अकाउंट के लिए एक shell alias, कई Claude Code अकाउंट चलाने का सबसे आसान तरीका है। ccam अंदर से यही करता है और साथ में साझा करना, स्विच करना और उपयोग की निगरानी जोड़ता है।

## alias वाला सेटअप

हर अतिरिक्त अकाउंट के लिए `~/.bashrc` या `~/.zshrc` में एक लाइन जोड़ें, नया टर्मिनल खोलें और लॉगिन करने के लिए alias एक बार चलाएँ:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` आपका पहला अकाउंट बना रहता है और `claude-work` दूसरा अकाउंट शुरू करता है। यही वह तरीका है जिसका [Claude Code के दस्तावेज़ में वर्णन है](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)।

## ccam वाला सेटअप

ccam इंस्टॉल करें और अकाउंट जोड़ें, यहाँ उसकी अपनी शॉर्टकट कमांड के साथ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` और `ccwork` अकाउंट शुरू करते हैं, और दोनों scripts में भी काम करते हैं।

## आमने-सामने

| | Shell alias | ccam |
| --- | --- | --- |
| अलग-अलग लॉगिन | हाँ | हाँ |
| हर अकाउंट का सेटअप | `~/.bashrc` संपादित करें, लॉगिन करें | `ccam add` |
| Scripts और दूसरे प्रोग्राम में काम करता है | नहीं, aliases केवल shell में होते हैं | हाँ, `ccam work` और `ccwork` जैसे शॉर्टकट असली कमांड हैं |
| साझा सेटिंग्स, `CLAUDE.md`, एक्सटेंशन | हाथ से copy या link करें | हाँ, हर आइटम अलग से चुनें, बैकअप के साथ |
| साझा MCP सर्वर | हाथ से, हर अकाउंट में | एक बार `ccam mcp add` |
| दूसरे अकाउंट में बातचीत जारी रखना | नहीं | `ccam next --continue` |
| हर अकाउंट का उपयोग | नहीं | `ccam usage` |
| सबसे ज़्यादा गुंजाइश वाला अकाउंट चुनना | नहीं | `ccam next --auto` |
| सीमा पर `-p` job को दूसरे अकाउंट के रूप में फिर से चलाना | नहीं | `ccam --fallback` |
| हर प्रोजेक्ट फ़ोल्डर के लिए एक अकाउंट | नहीं | `ccam folder` |
| स्टेटस लाइन में अकाउंट का नाम | नहीं | `ccam statusline install` |
| टूटे लिंक ठीक करना | हाथ से | `ccam doctor --fix` |
| अतिरिक्त सॉफ़्टवेयर | कोई नहीं | एक static बाइनरी |

अगर आपके पास दो अकाउंट हैं और आप उनके बीच कुछ भी साझा नहीं करना चाहते, तो aliases चुनें। अगर आप चाहते हैं कि सभी अकाउंट एक ही सेटअप जैसे लगें, तो ccam चुनें। आप aliases से शुरू करके बाद में ccam पर जा सकते हैं: `ccam import` मौजूदा `~/.claude-*` फ़ोल्डरों को उसी जगह रजिस्टर करता है जहाँ वे हैं।

[एक कंप्यूटर पर कई अकाउंट कैसे इस्तेमाल करें](/guides/multiple-claude-code-accounts) दोनों तरीकों को क़दम-दर-क़दम समझाती है, और [ccam और दूसरे account switchers](/compare/account-switchers) ccam की तुलना उन टूल्स से करता है जो एक ही लॉगिन को बदलते हैं।
