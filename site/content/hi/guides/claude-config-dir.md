---
title: CLAUDE_CONFIG_DIR क्या है
description: Claude Code में CLAUDE_CONFIG_DIR क्या करता है, फ़ोल्डर में क्या रखा जाता है, macOS लॉगिन कहाँ रखता है, और अकाउंट्स के बीच कौन-सी फ़ाइलें साझा हो सकती हैं।
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` Claude Code को बताता है कि `~/.claude` की बजाय कौन-सा फ़ोल्डर इस्तेमाल करना है। अकाउंट से जुड़ी हर चीज़ वहीं रहती है: लॉगिन, सेटिंग्स, MCP सर्वर, plugins और बातचीत का इतिहास। इसे अलग-अलग फ़ोल्डरों पर सेट करने से आपको अलग-अलग अकाउंट मिलते हैं।

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## फ़ोल्डर के अंदर क्या है

Claude Code 2.1.302 के साथ जाँचा गया:

| पाथ | सामग्री | साझा हो सकता है |
| --- | --- | --- |
| `.credentials.json` | लॉगिन tokens (Linux), MCP OAuth tokens | नहीं |
| `.claude.json` | अकाउंट की जानकारी, user MCP सर्वर, हर प्रोजेक्ट की स्थिति | नहीं (MCP सर्वर अलग से) |
| `settings.json` | अनुमतियाँ, hooks, मॉडल, स्टेटस लाइन, plugins | हाँ |
| `CLAUDE.md` | ग्लोबल निर्देश | हाँ |
| `agents/`, `commands/`, `skills/`, `plugins/` | एक्सटेंशन | हाँ |
| `projects/` | `.jsonl` फ़ाइलों के रूप में बातचीत, हर प्रोजेक्ट की मेमोरी | हाँ |
| `history.jsonl` | Prompt इतिहास | हाँ |
| `sessions/`, `shell-snapshots/`, `cache/` और दूसरे | रनटाइम स्थिति | नहीं |

## दो बातें जिनमें लोग अक्सर फँसते हैं

**`CLAUDE_CONFIG_DIR` के बिना, `.claude.json` फ़ाइल `~/.claude` के अंदर नहीं होती।** तब Claude Code उसे फ़ोल्डर के बगल में `~/.claude.json` पर रखता है। अगर आप `CLAUDE_CONFIG_DIR=~/.claude` इस्तेमाल करना शुरू करते हैं, तो Claude Code इसकी बजाय `~/.claude/.claude.json` खोजता है और आपका अकाउंट नहीं पाता। `ccam import` इन दोनों को लिंक करके इसे संभालता है।

**macOS पर लॉगिन फ़ोल्डर में नहीं, कीचेन में होता है।** Claude Code उसे `Claude Code-credentials` नाम से रखता है, और जब `CLAUDE_CONFIG_DIR` सेट हो, तो इस नाम के बाद फ़ोल्डर पाथ का एक छोटा hash जुड़ता है। इसलिए अलग-अलग फ़ोल्डरों को कीचेन में अलग-अलग entries मिलती हैं। अगर कीचेन लिखने से मना कर दे, तो Claude Code फ़ोल्डर में `.credentials.json` का इस्तेमाल करता है।

## फ़ोल्डरों के बीच साझा करना

Symbolic links काम करते हैं: Claude Code लिंक की गई `settings.json` या `projects/` फ़ोल्डर के ज़रिए पढ़ता और लिखता है और लिंक को बनाए रखता है। MCP सर्वर इसका अपवाद हैं, क्योंकि वे लॉगिन के साथ `.claude.json` साझा करते हैं; उन्हें इसकी बजाय `--mcp-config` से दें।

ccam ठीक इन्हीं नियमों पर बना है। ccam फ़ोल्डरों को कैसे लिंक करता है, यह जानने के लिए [साझा करना](/docs/sharing) देखें, और पूरी तस्वीर के लिए [कई अकाउंट्स की गाइड](/guides/multiple-claude-code-accounts) देखें।
