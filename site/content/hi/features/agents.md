---
title: हर अकाउंट के background sessions
description: ccam agents आपके सभी अकाउंट्स के Claude Code background agents और खुले sessions एक ही table में दिखाता है, और ccam attach किसी एक को उसके id से खोलता है।
updated: 2026-10-23
order: 16
translation: machine
sourceHash: "092368f63363"
---

Claude Code background sessions (`claude --bg`, agent view) हर अकाउंट के लिए अलग रखता है: `claude agents` केवल उसी अकाउंट के sessions दिखाता है जिसके रूप में वह चलता है। `ccam agents` हर अकाउंट से एक साथ पूछता है और उन्हें एक साथ दिखाता है, और `ccam attach <id|name>` किसी background session को उसी अकाउंट के रूप में खोलता है जिसने उसे शुरू किया था।

## उपयोग

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` background session की स्थिति है (`working`, `blocked`, `done`, `failed`, `stopped`) और यह भी कि उसका process क्या कर रहा है (`busy`, `waiting`, `idle`); खुला interactive session `open` दिखाता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--all` | वे background sessions भी दिखाएँ जो पूरे हो चुके, विफल हुए या रोके गए हैं | बंद |
| `--json` | sessions को JSON के रूप में प्रिंट करें, हर एक के साथ वह `account` जिसका वह है | बंद |

`ccam attach <id|name> [claude args]` id के बाद का सब कुछ, flags समेत, बिना बदले `claude attach` को देता है: `ccam attach 49ec --dry-run` `--dry-run` को Claude Code को सौंप देता है। इसके बजाय कमांड प्रिंट करने के लिए `--dry-run` को id से पहले रखें (`ccam --dry-run attach 49ec` या `ccam attach --dry-run 49ec`)।

## यह कैसे काम करता है

1. हर अकाउंट के लिए ccam अकाउंट के `CLAUDE_CONFIG_DIR` के साथ `claude agents --json` चलाता है (दिए जाने पर `--all` के साथ), सभी एक साथ, हर एक 15 सेकंड के भीतर। यह अकाउंट की अपनी session सूची पढ़ता है और Claude Code की background service शुरू नहीं करता।
2. हर entry को अकाउंट का नाम मिलता है, और सूची शुरू होने के समय से क्रम में रखी जाती है, सबसे नई पहले। जो अकाउंट जवाब नहीं देता, उसका नाम standard error पर बताया जाता है और उसे छोड़ दिया जाता है, और फिर ccam code 1 के साथ बाहर निकलता है, ताकि कोई script विफलता को खाली सूची न समझ ले।
3. `ccam attach` हर अकाउंट के फ़ोल्डर में `jobs/<id>` खोजता है, जहाँ Claude Code background session रखता है, और उस अकाउंट के रूप में `claude attach <id>` चलाता है, उसी तरह जैसे `ccam <name>` Claude Code शुरू करता है। जो शब्द किसी id की शुरुआत नहीं है, उसे sessions के नामों में खोजा जाता है, बड़े-छोटे अक्षरों की परवाह किए बिना; जो शब्द id जैसा दिखता है, उसे कभी नाम नहीं माना जाता।

## सुरक्षा और विशेष स्थितियाँ

- Background session के id एक ही अकाउंट के होते हैं: किसी दूसरे अकाउंट के रूप में चलाया गया `claude attach` उन्हें नहीं ढूँढ सकता, और `ccam attach` यही सुलझाता है।
- दो अकाउंट्स में ऐसे sessions हो सकते हैं जिनके id एक जैसे शुरू होते हैं, या जिनके नामों में वही शब्द है; तब `ccam attach` उनके नाम बताता है और id या नाम का और हिस्सा माँगता है।
- ccam अपने साझा MCP सर्वर Claude Code को एक ही आर्ग्युमेंट में `--mcp-config=<file>` से देता है, जिसे `claude agents` और `claude attach` स्वीकार करते हैं; background sessions दोबारा शुरू होने पर भी वही फ़ाइल इस्तेमाल करते रहते हैं।
- API key अकाउंट की background service अपनी key उस shell से लेती है जिसने उसे शुरू किया। ऐसे अकाउंट का पहला background session ccam के ज़रिए शुरू करें (`ccam gw --bg …`), ताकि service के पास key हो।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [scripts के लिए अकाउंट चुनना](/features/pick)
- [डैशबोर्ड](/features/dashboard)
