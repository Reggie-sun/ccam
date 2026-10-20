---
title: शेल aliases विरुद्ध ccam
description: अनेक Claude Code खाती चालवण्यासाठी CLAUDE_CONFIG_DIR सह शेल aliases आणि ccam यांची शेजारी शेजारी तुलना.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: ea24bb091caf
---

प्रत्येक खात्यासाठी एक शेल alias हा अनेक Claude Code खाती चालवण्याचा सर्वात सोपा मार्ग आहे. ccam आतून तेच करते आणि त्यावर शेअरिंग, खाते बदलणे आणि वापराचा मागोवा जोडते.

## alias वापरून सेटअप

प्रत्येक जादा खात्यासाठी `~/.bashrc` किंवा `~/.zshrc` मध्ये एक ओळ जोडा, नवीन टर्मिनल उघडा आणि लॉगिन करण्यासाठी alias एकदा चालवा:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` तुमचे पहिले खाते राहते आणि `claude-work` दुसरे खाते सुरू करते. हीच पद्धत [Claude Code च्या दस्तऐवजात सांगितली आहे](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## ccam वापरून सेटअप

ccam इन्स्टॉल करा आणि खाते जोडा, येथे त्याच्या स्वतःच्या शॉर्टकट कमांडसह:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` आणि `ccwork` खाते सुरू करतात, आणि दोन्ही scripts मध्येही काम करतात.

## शेजारी शेजारी

| | शेल alias | ccam |
| --- | --- | --- |
| वेगळे लॉगिन | हो | हो |
| प्रत्येक खात्यासाठी सेटअप | `~/.bashrc` संपादित करा, लॉगिन करा | `ccam add` |
| स्क्रिप्ट्स आणि इतर प्रोग्राम्समध्ये चालते | नाही, aliases फक्त शेलपुरते | हो, `ccam work` आणि `ccwork` सारखे शॉर्टकट ह्या खऱ्या कमांड्स आहेत |
| शेअर केलेल्या सेटिंग्ज, `CLAUDE.md`, एक्स्टेंशन्स | हाताने कॉपी किंवा लिंक | आयटमनुसार, बॅकअपसह |
| शेअर केलेले MCP सर्व्हर | हाताने, प्रत्येक खात्यासाठी | एकदा `ccam mcp add` |
| दुसऱ्या खात्यात संभाषण पुढे चालू ठेवणे | नाही | `ccam next --continue` |
| प्रत्येक खात्याचा वापर | नाही | `ccam usage` |
| सर्वात जास्त जागा असलेले खाते निवडणे | नाही | `ccam next --auto` |
| स्टेटस लाइनमध्ये खात्याचे नाव | नाही | `ccam statusline install` |
| तुटलेल्या लिंक्स दुरुस्त करणे | हाताने | `ccam doctor --fix` |
| अतिरिक्त सॉफ्टवेअर | काहीही नाही | एक static बायनरी |

तुमच्याकडे दोन खाती असतील आणि त्यांच्यात काहीही शेअर करायचे नसेल तर aliases निवडा. खाती एकाच सेटअपसारखी वाटावीत असे हवे असेल तर ccam निवडा. तुम्ही aliases ने सुरुवात करून नंतर ccam कडे जाऊ शकता: `ccam import` आधीपासून असलेले `~/.claude-*` फोल्डर आहेत त्याच जागी नोंदवते.

[एका कॉम्प्युटरवर अनेक खाती कशी वापरायची](/guides/multiple-claude-code-accounts) हे दोन्ही पद्धती टप्प्याटप्प्याने समजावून सांगते, आणि [ccam आणि इतर account switchers](/compare/account-switchers) एकच लॉगिन बदलणाऱ्या टूल्सशी ccam ची तुलना करते.
