---
title: extensions की जाँच
description: ccam ext दिखाता है कि हर अकाउंट कौन-से skills, commands, plugins, MCP सर्वर और hooks लोड करता है, हर एक कहाँ से आता है, और क्या टूटा है।
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

हर cell उस अकाउंट में उस आइटम का source बताता है: share group (`default`), `account` जब वह अकाउंट के अपने फ़ोल्डर में हो, या `-` जब अकाउंट के पास वह न हो। `!` किसी समस्या वाले आइटम को चिह्नित करता है, जिसकी व्याख्या table के नीचे होती है।

## यह क्या पढ़ता है

| प्रकार | कहाँ | दिखाया गया source |
| --- | --- | --- |
| Skills | `skills/` में हर वह फ़ोल्डर जिसमें `SKILL.md` हो | वह group जिसकी ओर `skills/` link करता है, या `account` |
| Commands, agents | `commands/` और `agents/` के अंदर `.md` फ़ाइलें, सब-फ़ोल्डर `folder/name` के रूप में | group, या `account` |
| Plugins | `plugins/installed_plugins.json` | group, या `account` |
| MCP सर्वर | group का `mcp.json` (जब MCP साझा हो) और अकाउंट का `.claude.json` | group, या `account` |
| Hooks | `settings.json` में `hooks`, हर event के हिसाब से गिने गए | group, या `account` |
| यह प्रोजेक्ट | मौजूदा फ़ोल्डर में `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` और `.mcp.json` | एक बार दिखाया जाता है, हर अकाउंट के लिए |

## यह कौन-सी समस्याएँ बताता है

- **टूटे links**: `skills/`, `commands/` या `agents/` में कोई link जिसका target अब नहीं है, या कोई साझा फ़ोल्डर (`skills`, `commands`, `agents`, `plugins`) जो किसी ग़ायब जगह की ओर link करता है (इसे `ccam doctor --fix` से ठीक करें)।
- **`SKILL.md` के बिना skill फ़ोल्डर**: Claude Code उसे छोड़ देता है।
- **दो जगह अलग-अलग परिभाषित MCP सर्वर**: group का `mcp.json` और अकाउंट का `.claude.json`, या इनमें से एक और प्रोजेक्ट का `.mcp.json`। दोनों जगह एक जैसी सेटिंग्स (जो `mcp_strategy = "sync"` लिखता है) कोई समस्या नहीं है। किसी अकाउंट के अपने सर्वर को `ccam mcp import <account> <server>` से group में ले जाएँ।
- **`installed_plugins.json` जो मान्य JSON नहीं है।**

## सुरक्षा और विशेष स्थितियाँ

- `ccam ext` कुछ नहीं बदलता और कोई नेटवर्क request नहीं भेजता। यह सर्वर के नाम और hooks की गिनती पढ़ता है, MCP definitions में `env` या headers के मान कभी नहीं।
- यह वही दिखाता है जो disk पर है। कोई plugin जो इंस्टॉल है लेकिन सेटिंग्स में बंद है, या managed settings से रोका गया सर्वर, फिर भी दिखता है।
- उस प्रोजेक्ट का `.claude` फ़ोल्डर और `.mcp.json` शामिल करने के लिए इसे प्रोजेक्ट फ़ोल्डर में चलाएँ।

## संबंधित

- [अकाउंट्स के बीच साझा करना](/features/share)
- [साझा MCP सर्वर](/features/mcp)
- [जाँच और मरम्मत](/features/doctor)
