---
title: CLAUDE_CONFIG_DIR क्या है
description: Claude Code में CLAUDE_CONFIG_DIR क्या करता है, फ़ोल्डर में क्या रखा जाता है, macOS लॉगिन कहाँ रखता है, और अकाउंट्स के बीच कौन-सी फ़ाइलें साझा हो सकती हैं।
updated: 2026-10-23
order: 2
translation: machine
sourceHash: 2cf08478085f
---

`CLAUDE_CONFIG_DIR` Claude Code को बताता है कि `~/.claude` की बजाय कौन-सा फ़ोल्डर इस्तेमाल करना है। अकाउंट से जुड़ी हर चीज़ वहीं रहती है: लॉगिन, सेटिंग्स, MCP सर्वर, plugins और बातचीत का इतिहास। इसे अलग-अलग फ़ोल्डरों पर सेट करने से आपको अलग-अलग अकाउंट मिलते हैं।

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## क्या CLAUDE_CONFIG_DIR आधिकारिक रूप से समर्थित है?

हाँ। Claude Code इसे अपने [environment variables](https://code.claude.com/docs/en/env-vars) में configuration फ़ोल्डर बदलने के तरीक़े के रूप में बताता है, जो "कई अकाउंट साथ-साथ चलाने के लिए उपयोगी" है, और उसकी [authentication गाइड](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) कई अकाउंट्स से लॉगिन करने के लिए इसी का इस्तेमाल करती है। आप इसे अपने shell में या user settings में सेट कर सकते हैं, लेकिन किसी प्रोजेक्ट की settings में नहीं।

## फ़ोल्डर के अंदर क्या है?

Claude Code 2.1.309 के साथ जाँचा गया:

| पाथ | सामग्री | साझा हो सकता है |
| --- | --- | --- |
| `.credentials.json` | लॉगिन tokens (Linux, Windows), MCP OAuth tokens | नहीं |
| `.claude.json` | अकाउंट की जानकारी, user MCP सर्वर, हर प्रोजेक्ट की स्थिति | नहीं (MCP सर्वर अलग से) |
| `settings.json` | अनुमतियाँ, hooks, मॉडल, स्टेटस लाइन, plugins | हाँ |
| `CLAUDE.md` | ग्लोबल निर्देश | हाँ |
| `agents/`, `commands/`, `skills/`, `plugins/` | एक्सटेंशन | हाँ |
| `projects/` | `.jsonl` फ़ाइलों के रूप में बातचीत, हर प्रोजेक्ट की मेमोरी | हाँ |
| `history.jsonl` | Prompt इतिहास | हाँ |
| `sessions/`, `shell-snapshots/`, `cache/` और दूसरे | रनटाइम स्थिति | नहीं |

## लोग अक्सर कहाँ फँसते हैं?

**`CLAUDE_CONFIG_DIR` के बिना, `.claude.json` फ़ाइल `~/.claude` के अंदर नहीं होती।** तब Claude Code उसे फ़ोल्डर के बगल में `~/.claude.json` पर रखता है। अगर आप `CLAUDE_CONFIG_DIR=~/.claude` इस्तेमाल करना शुरू करते हैं, तो Claude Code इसकी बजाय `~/.claude/.claude.json` खोजता है और आपका अकाउंट नहीं पाता। `ccam import` इन दोनों को लिंक करके इसे संभालता है। अगर `~/.claude` के पास पहले से अपनी `.claude.json` है, जो `CLAUDE_CONFIG_DIR=~/.claude` के साथ पहले चलाने से बनी थी, तो ccam दोनों फ़ाइलें छोड़ देता है और यह बताता है, क्योंकि केवल आप जानते हैं कि कौन-सी सही है।

**macOS पर लॉगिन फ़ोल्डर में नहीं, कीचेन में होता है।** Claude Code उसे `Claude Code-credentials` नाम से रखता है, और जब `CLAUDE_CONFIG_DIR` सेट हो, तो इस नाम के बाद `-` और फ़ोल्डर पाथ के SHA-256 के पहले 8 hex अंक जुड़ते हैं, उस account के लिए जिसका नाम `$USER` बताता है। इसलिए अलग-अलग फ़ोल्डरों को कीचेन में अलग-अलग entries मिलती हैं, और जो फ़ोल्डर move किया जाए या जिसका नाम बदला जाए, उसका लॉगिन खो जाता है; इसीलिए macOS पर `ccam rename` फ़ोल्डर का नाम बनाए रखता है। अगर कीचेन लिखने से मना कर दे, तो Claude Code फ़ोल्डर में `.credentials.json` का इस्तेमाल करता है, जैसा [उसके दस्तावेज़](https://code.claude.com/docs/en/authentication#credential-management) में बताया गया है।

**हर sign-in फ़ोल्डर में नहीं रहता।** API key के बिना किया गया Claude Console sign-in फ़ोल्डर के बाहर रखा जाता है, इसलिए अलग-अलग फ़ोल्डर ऐसे दो sign-in को अलग नहीं रख पाते। claude.ai लॉगिन और API keys हर फ़ोल्डर में अलग रखे जाते हैं।

## Windows पर फ़ोल्डर कहाँ है?

डिफ़ॉल्ट रूप से `%USERPROFILE%\.claude`, और लॉगिन उसके अंदर `.credentials.json` में, या जब Claude Code को Windows Credential Manager इस्तेमाल करने के लिए सेट किया गया हो, तो उसमें, macOS की तरह फ़ोल्डर पाथ से बने नाम के तहत। हर अतिरिक्त अकाउंट के लिए `CLAUDE_CONFIG_DIR` को किसी दूसरे फ़ोल्डर पर सेट करें, ठीक वैसे ही जैसे Linux और macOS पर।

## क्या फ़ोल्डर फ़ाइलें साझा कर सकते हैं?

Symbolic links काम करते हैं: Claude Code लिंक की गई `settings.json` या `projects/` फ़ोल्डर के ज़रिए पढ़ता और लिखता है और लिंक को बनाए रखता है। MCP सर्वर इसका अपवाद हैं, क्योंकि वे लॉगिन के साथ `.claude.json` साझा करते हैं; उन्हें इसकी बजाय `--mcp-config` से दें।

ccam ठीक इन्हीं नियमों पर बना है। ccam फ़ोल्डरों को कैसे लिंक करता है, यह जानने के लिए [साझा करना](/docs/sharing) देखें, और पूरी तस्वीर के लिए [कई अकाउंट्स की गाइड](/guides/multiple-claude-code-accounts) देखें।
