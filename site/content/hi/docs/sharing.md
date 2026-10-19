---
title: MCP सर्वर, सेटिंग्स और इतिहास साझा करें
description: ccam कैसे Claude Code अकाउंट्स के बीच MCP सर्वर, settings.json, CLAUDE.md, agents, skills और बातचीत का इतिहास साझा करता है, और हर लॉगिन अलग रखता है।
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

हर अकाउंट का अपना लॉगिन रहता है। बाक़ी सब कुछ आइटम के हिसाब से साझा किया जा सकता है: ccam अकाउंट की फ़ाइल या फ़ोल्डर को एक साझा ग्रुप फ़ोल्डर से लिंक करता है, इसलिए एक अकाउंट में किया गया बदलाव उस ग्रुप के सभी अकाउंट्स में दिखता है।

## क्या साझा किया जा सकता है

| आइटम | विकल्प | क्या साझा होता है |
| --- | --- | --- |
| MCP सर्वर | `--mcp` | ग्रुप की `mcp.json`, जो शुरू होते समय Claude Code को दी जाती है |
| सेटिंग्स | `--settings` | `settings.json`: अनुमतियाँ, hooks, मॉडल, स्टेटस लाइन |
| ग्लोबल मेमोरी | `--memory` | `CLAUDE.md` |
| एक्सटेंशन | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| बातचीत का इतिहास | `--history` | `projects/` और `history.jsonl` |

नए अकाउंट बातचीत के इतिहास को छोड़कर सब कुछ साझा करते हैं।

लॉगिन फ़ाइलें `.credentials.json` और `.claude.json` कभी साझा नहीं होतीं। ccam केवल एक तय सूची के पाथ लिंक करता है, और ये दोनों उस सूची में नहीं हैं।

## ग्रुप

`on` का मतलब `default` ग्रुप है। आप कोई नाम भी दे सकते हैं, ताकि कोई आइटम केवल कुछ अकाउंट्स के बीच साझा हो:

```sh
ccam share work --history team
ccam share side --history team
```

अब `work` और `side` आपस में बातचीत साझा करते हैं, लेकिन आपके दूसरे अकाउंट्स के साथ नहीं।

## साझा करना चालू करना

जब कोई अकाउंट किसी आइटम को साझा करना शुरू करता है, तो ccam:

1. हर उस फ़ाइल की योजना प्रिंट करता है जिसे वह move, copy या link करेगा, और शुरू करने से पहले पूछता है;
2. अगर ग्रुप अभी ख़ाली है तो अकाउंट की फ़ाइलें ग्रुप फ़ोल्डर में move करता है, वरना वे फ़ाइलें जोड़ता है जो ग्रुप में अभी नहीं हैं;
3. जब एक ही फ़ाइल अलग हो, तो पूछता है कि ग्रुप का वर्ज़न रखें, आपका इस्तेमाल करें, या JSON को merge करें;
4. आपकी मूल फ़ाइल को अकाउंट फ़ोल्डर के अंदर `.ccam-backup/<time>/` में रखता है;
5. मूल फ़ाइल की जगह ग्रुप का लिंक लगा देता है।

बातचीत की फ़ाइलों के नाम उनके session ID पर होते हैं, इसलिए कई अकाउंट्स के इतिहास बिना टकराव के एक साथ जुड़ जाते हैं।

## साझा करना बंद करना

```sh
ccam share work --history off
```

लिंक हट जाता है और अकाउंट को ग्रुप की मौजूदा सामग्री की अपनी कॉपी मिल जाती है। इसकी बजाय ख़ाली शुरुआत करने के लिए `--fresh` जोड़ें। ग्रुप अपनी सामग्री दूसरे अकाउंट्स के लिए रखे रहता है।

## मौजूदा सेटअप कॉपी करें

अगर आपने ccam से पहले Claude Code इस्तेमाल किया है, तो आपकी सेटिंग्स `~/.claude` में हैं। `ccam migrate` उन्हें एक ग्रुप में कॉपी करता है, ताकि उन आइटम को साझा करने वाले हर अकाउंट को वे मिल जाएँ:

```sh
ccam migrate
```

यह क्या कॉपी करता है:

- `settings.json` और ग्लोबल `CLAUDE.md`;
- एक्सटेंशन `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` और `rules/`;
- आपके user MCP सर्वर, यानी `~/.claude.json` की `mcpServers` एंट्री, ग्रुप की `mcp.json` में;
- बातचीत का इतिहास (`projects/` और `history.jsonl`), लेकिन केवल तब जब आप `--history` जोड़ें।

स्रोत फ़ोल्डर को केवल पढ़ा जाता है: उसमें कुछ भी move, link या बदला नहीं जाता, और सादा `claude` पहले की तरह काम करता रहता है। लॉगिन फ़ाइलें `.credentials.json` और `.claude.json` कभी कॉपी नहीं होतीं।

कुछ भी कॉपी करने से पहले ccam योजना प्रिंट करता है और पूछता है:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

आख़िर में ccam उन अकाउंट्स की सूची दिखाता है जो ग्रुप के ज़रिए हर आइटम साझा करते हैं। अगर कुछ अकाउंट ऐसा नहीं करते, तो यह वे `ccam share` कमांड प्रिंट करता है जो उनके लिए साझा करना चालू करती हैं।

जब ग्रुप में किसी फ़ाइल की अलग कॉपी पहले से हो, तो ccam वैसे ही पूछता है जैसे साझा करना चालू करते समय: ग्रुप की कॉपी रखें या अपनी इस्तेमाल करें। `settings.json` और दूसरी JSON फ़ाइलें merge भी की जा सकती हैं, और `history.jsonl` लाइन-दर-लाइन। ग्रुप की पिछली कॉपियाँ ग्रुप फ़ोल्डर के अंदर `.ccam-backup/<time>/` में सहेजी जाती हैं।

| विकल्प | असर |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | केवल वही आइटम कॉपी करें जिनका आप नाम देते हैं; इनके बिना बातचीत के इतिहास को छोड़कर सब कुछ |
| `-g`, `--group <group>` | `default` की बजाय इस ग्रुप में कॉपी करें |
| `-y`, `--yes` | बिना पूछे कॉपी करें; कोई फ़ाइल अलग हो तो ग्रुप की कॉपी रखें |

कुछ भी कॉपी किए बिना योजना देखने के लिए कमांड से पहले `--dry-run` लगाएँ, जैसे `ccam --dry-run migrate`। किसी दूसरे फ़ोल्डर से कॉपी करने के लिए उसका नाम दें: `ccam migrate ~/.claude-work --group work`।

## MCP सर्वर

Claude Code MCP सर्वर `.claude.json` में, लॉगिन के साथ ही रखता है, इसलिए यह फ़ाइल साझा नहीं की जा सकती। ccam साझा सर्वर ग्रुप की `mcp.json` में रखता है और Claude Code को `--mcp-config=<group>/mcp.json` के साथ शुरू करता है। आपके अपने `--mcp-config` विकल्प इसके ऊपर लागू होते रहते हैं।

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Session के अंदर `claude mcp add` से जोड़ा गया सर्वर केवल उसी अकाउंट के लिए सेव होता है। अगली बार जब आप वह अकाउंट शुरू करते हैं, तो ccam उस सर्वर का नाम बताता है और वह कमांड भी जो उसे ग्रुप में move करती है:

```sh
ccam mcp import work gh
```

जो MCP सर्वर OAuth से साइन इन करते हैं, वे अपने tokens अकाउंट के लॉगिन के साथ रखते हैं, इसलिए उन्हें हर अकाउंट में एक बार अधिकृत करें।

अगर आप चाहते हैं कि सर्वर `claude mcp list` और `/mcp` में सामान्य user सर्वर की तरह दिखें, तो `~/.config/ccam/config.toml` में `mcp_strategy = "sync"` सेट करें। तब ccam शुरू करने से पहले उन्हें हर अकाउंट की `.claude.json` में लिखता है, Claude Code के अपने file lock के तहत और बैकअप के साथ।

## एक ही फ़ोल्डर में दो अकाउंट

जो अकाउंट बातचीत का इतिहास साझा करते हैं, वे एक ही फ़ाइलों में लिखते हैं। अगर आप किसी ऐसे फ़ोल्डर में दूसरा Claude Code शुरू करते हैं जहाँ उसी ग्रुप का कोई दूसरा अकाउंट पहले से चल रहा है, तो ccam आपको चेतावनी देता है और पूछता है कि जारी रखना है या नहीं।

## जब कुछ टूट जाए

```sh
ccam doctor
ccam doctor --fix
```

अगर कोई साझा पाथ सामान्य फ़ाइल बन गया है, उदाहरण के लिए जब आपने हाथ से कोई बैकअप restore किया हो, तो ccam Claude Code शुरू करने से पहले रुक जाता है और `ccam doctor --fix` की ओर इशारा करता है, जो फ़ाइल को ग्रुप में merge करता है और बैकअप रखता है।
