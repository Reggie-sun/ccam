---
title: extensions की जाँच
description: ccam ext दिखाता है कि हर अकाउंट कौन-से skills, commands, plugins, MCP सर्वर और hooks लोड करता है, हर एक कहाँ से आता है, और क्या टूटा है।
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
---

`ccam ext` इस सवाल का जवाब देता है: "यह skill (या MCP सर्वर) एक अकाउंट में काम करता है और दूसरे में क्यों नहीं?"। यह हर अकाउंट के लिए वे extensions दिखाता है जिन्हें Claude Code लोड करता है, और यह कि हर एक कहाँ से आता है: कोई share group, अकाउंट का अपना फ़ोल्डर, या वह प्रोजेक्ट जिसमें आप हैं। table के नीचे यह वे बातें दिखाता है जो किसी extension को आपकी उम्मीद के मुताबिक लोड होने से रोकती हैं। यह केवल पढ़ता है।

## उपयोग

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

हर cell उस अकाउंट में उस आइटम का source बताता है: share group (`default`), `account` जब वह अकाउंट के अपने फ़ोल्डर में हो, `local` ऐसे MCP सर्वर के लिए जो अकाउंट के पास केवल इस फ़ोल्डर के लिए है, या `-` जब अकाउंट के पास वह न हो। जो plugin इंस्टॉल है लेकिन चालू नहीं है, उसके source के बाद `(off)` दिखता है। `!` किसी समस्या वाले आइटम को चिह्नित करता है, जिसकी व्याख्या table के नीचे होती है; जो समस्याएँ किसी अकाउंट की नहीं हैं, जैसे ऐसे सर्वर जो केवल सादे `claude` के पास हैं, बिना नाम के आख़िर में आती हैं। `ccam ext --json` वही रिपोर्ट प्रिंट करता है (`accounts`, प्रोजेक्ट, `problems`), उन MCP सर्वरों पर `remote` के साथ जिन तक Claude Code नेटवर्क से पहुँचता है; कोई अकाउंट न होने पर यह `accounts: []` के साथ रिपोर्ट प्रिंट करता है।

## यह क्या पढ़ता है

| प्रकार | कहाँ | दिखाया गया source |
| --- | --- | --- |
| Skills | `skills/` में हर वह फ़ोल्डर जिसमें `SKILL.md` हो | वह group जिसकी ओर `skills/` link करता है, या `account` |
| Commands, agents | `commands/` और `agents/` के अंदर `.md` फ़ाइलें, सब-फ़ोल्डर `folder/name` के रूप में | group, या `account` |
| Plugins | `plugins/installed_plugins.json`, अकाउंट और प्रोजेक्ट की सेटिंग्स में `enabledPlugins` से चालू या बंद; केवल एक प्रोजेक्ट के लिए इंस्टॉल किया गया plugin केवल उसी प्रोजेक्ट में दिखता है | group, या `account`; चालू न होने पर `(off)` |
| MCP सर्वर | group का `mcp.json` (जब MCP साझा हो), अकाउंट का `.claude.json` (उसके अपने सर्वर, और वे जो उसके पास केवल इस फ़ोल्डर के लिए हैं), और उसके plugins के सर्वर, `plugin:<plugin>:<server>` के रूप में | group, `account` या `local` |
| Hooks | `settings.json` में `hooks`, हर event के हिसाब से गिने गए | group, या `account` |
| यह प्रोजेक्ट | मौजूदा फ़ोल्डर में `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` और `.mcp.json` | एक बार दिखाया जाता है; जिस अकाउंट ने फ़ोल्डर पर भरोसा नहीं किया है, उसे उसके `.mcp.json` सर्वरों के लिए एक समस्या मिलती है |

## यह कौन-सी समस्याएँ बताता है

- **टूटे links**: `skills/`, `commands/` या `agents/` में कोई link जिसका target अब नहीं है, या कोई साझा फ़ोल्डर (`skills`, `commands`, `agents`, `plugins`) जो किसी ग़ायब जगह की ओर link करता है (इसे `ccam doctor --fix` से ठीक करें, जो ग़ायब group फ़ोल्डर फिर से बनाता है)।
- **`SKILL.md` के बिना skill फ़ोल्डर**: Claude Code उसे छोड़ देता है।
- **दो जगह अलग-अलग परिभाषित MCP सर्वर**: group का `mcp.json` और अकाउंट का `.claude.json`, या इनमें से एक और प्रोजेक्ट का `.mcp.json`। दोनों जगह एक जैसी सेटिंग्स (जो `mcp_strategy = "sync"` लिखता है) कोई समस्या नहीं है। किसी अकाउंट के अपने सर्वर को `ccam mcp import <account> <server>` से group में ले जाएँ।
- **`installed_plugins.json` जो मान्य JSON नहीं है**, ऐसा plugin जिसका नाम `enabledPlugins` में नहीं है, और ऐसा plugin जिसका फ़ोल्डर ग़ायब है।
- **ऐसे सर्वर जिन्हें केवल एक अकाउंट लोड करता है**: ऐसा सर्वर जो केवल इस अकाउंट के `.claude.json` में है, जिसे उसके group के दूसरे अकाउंट लोड नहीं करते, और ऐसा सर्वर जो इस फ़ोल्डर के लिए अकाउंट के local scope में है (जो `claude mcp add` करता है, जब तक `-s user` न दिया जाए)। हर एक के साथ वह `ccam mcp import` कमांड आती है जो उसे साझा करती है। सादे `claude` के `~/.claude.json` में ऐसे सर्वर जिन्हें कोई अकाउंट लोड नहीं करता, `ccam migrate --mcp` कमांड के साथ आते हैं।
- **मंज़ूरी का इंतज़ार कर रहा प्रोजेक्ट सर्वर**: अकाउंट ने अभी फ़ोल्डर पर भरोसा नहीं किया है, इसलिए Claude Code एक बार पूछता है, और तब तक प्रोजेक्ट के `.mcp.json` सर्वर इंतज़ार करते हैं। समस्या बताती है कि उन्हें कहाँ मंज़ूर करें: settings group के अकाउंट्स के, या इस अकाउंट के, `settings.json` में `"enabledMcpjsonServers"`। किसी git repository के main worktree को दिया गया भरोसा उसके दूसरे worktrees में भी गिना जाता है।
- **साइन-इन और connectors**: ऐसा MCP सर्वर जिसे इस अकाउंट में साइन-इन चाहिए जबकि कोई दूसरा अकाउंट उसमें साइन इन है (साइन-इन केवल एक अकाउंट के लिए होता है; `/mcp` से साइन इन करें), और claude.ai connectors जिन्हें किसी दूसरे अकाउंट ने इस्तेमाल किया है और इसने नहीं।

## सुरक्षा और विशेष स्थितियाँ

- `ccam ext` कुछ नहीं बदलता और कोई नेटवर्क request नहीं भेजता। यह सर्वर के नाम और hooks की गिनती पढ़ता है, MCP definitions में `env` या headers के मान कभी नहीं। साइन-इन की जाँच के लिए यह हर अकाउंट के लॉगिन में MCP साइन-इन के नाम पढ़ता है (macOS पर कीचेन से), tokens कभी नहीं, और केवल उन अकाउंट्स के लिए जो ऐसा सर्वर लोड करते हैं जिसे कोई दूसरा अकाउंट भी लोड करता है।
- यह वही दिखाता है जो disk पर है। managed settings से रोका गया सर्वर फिर भी दिखता है; सेटिंग्स में बंद किया गया plugin `(off)` दिखाता है।
- उस प्रोजेक्ट का `.claude` फ़ोल्डर और `.mcp.json` शामिल करने के लिए इसे प्रोजेक्ट फ़ोल्डर में चलाएँ।

## संबंधित

- [अकाउंट्स के बीच साझा करना](/features/share)
- [साझा MCP सर्वर](/features/mcp)
- [जाँच और मरम्मत](/features/doctor)
