---
title: प्रत्येक खात्याची पार्श्वभूमीतील सत्रे
description: ccam agents तुमच्या सर्व खात्यांचे Claude Code background agents आणि उघडी सत्रे एकाच table मध्ये दाखवते, आणि ccam attach त्यांपैकी एक त्याच्या id ने उघडते.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code पार्श्वभूमीतील सत्रे (`claude --bg`, agent view) प्रत्येक खात्यासाठी वेगळी ठेवते: `claude agents` फक्त ते ज्या खात्याने चालते त्या खात्याची सत्रे दाखवते. `ccam agents` प्रत्येक खात्याला एकाच वेळी विचारते आणि ती सर्व एकत्र दाखवते, आणि `ccam attach <id>` पार्श्वभूमीतील सत्र ते ज्या खात्याने सुरू केले त्याच खात्याने उघडते.

## वापर

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` म्हणजे पार्श्वभूमीतील सत्राची स्थिती (`working`, `blocked`, `done`, `failed`, `stopped`) आणि त्याची process काय करत आहे (`busy`, `waiting`, `idle`); उघडे interactive सत्र `open` दाखवते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--all` | पूर्ण झालेली, अयशस्वी किंवा थांबवलेली पार्श्वभूमीतील सत्रेसुद्धा दाखवा | बंद |
| `--json` | सत्रे JSON म्हणून छापा, प्रत्येक ज्या `account` चे आहे त्यासह | बंद |

`ccam attach <id> [claude args]` id नंतरचे आर्ग्युमेंट्स `claude attach` ला देते. `--dry-run` सह ते कमांड छापते.

## हे कसे काम करते

1. प्रत्येक खात्यासाठी ccam खात्याच्या `CLAUDE_CONFIG_DIR` सह `claude agents --json` चालवते (दिले असल्यास `--all` सह), सर्व एकाच वेळी, प्रत्येक 15 सेकंदांच्या आत. हे खात्याची स्वतःची सत्रांची यादी वाचते आणि Claude Code ची background service सुरू करत नाही.
2. प्रत्येक नोंदीला खात्याचे नाव मिळते, आणि यादी सुरू झाल्याच्या वेळेनुसार, सर्वात नवीन आधी, लावली जाते. उत्तर न देणाऱ्या खात्याचे नाव standard error वर सांगितले जाते आणि ते वगळले जाते, आणि मग ccam code 1 ने बाहेर पडते, जेणेकरून script अपयशाला रिकामी यादी समजणार नाही.
3. `ccam attach` प्रत्येक खात्याच्या फोल्डरमध्ये `jobs/<id>` शोधते, जिथे Claude Code पार्श्वभूमीतील सत्र ठेवते, आणि `ccam <name>` Claude Code सुरू करते त्याच प्रकारे त्या खात्याने `claude attach <id>` चालवते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- पार्श्वभूमीतील सत्रांचे ids एकाच खात्याचे असतात: दुसऱ्या खात्याने चालवलेले `claude attach` ते शोधू शकत नाही, आणि नेमका हाच प्रश्न `ccam attach` सोडवते.
- दोन खात्यांकडे एकसारख्या सुरुवातीचे ids असलेली सत्रे असू शकतात; तेव्हा `ccam attach` दोघांची नावे सांगते आणि id चा आणखी भाग मागते.
- ccam आपले शेअर केलेले MCP सर्व्हर Claude Code ला एकाच आर्ग्युमेंटमध्ये `--mcp-config=<file>` ने देते, जे `claude agents` आणि `claude attach` स्वीकारतात; पार्श्वभूमीतील सत्रे पुन्हा सुरू झाल्यावरही तीच फाइल वापरत राहतात.
- API key खात्याची background service तिची key ज्या shell ने ती सुरू केली त्यातून घेते. अशा खात्याचे पहिले पार्श्वभूमीतील सत्र ccam मार्फत सुरू करा (`ccam gw --bg …`), जेणेकरून service ला key मिळेल.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [scripts साठी खाते निवडणे](/features/pick)
- [डॅशबोर्ड](/features/dashboard)
