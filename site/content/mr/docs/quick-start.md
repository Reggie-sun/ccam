---
title: क्विक स्टार्ट
description: ccam add ने दोन Claude Code खाती जोडा, claude1, claude2 किंवा ccam <name> ने त्यांच्यात बदल करा, आणि ती काय शेअर करतील ते निवडा.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

प्रत्येक खाते `ccam add` ने जोडा, नंतर `claude1`, `claude2` किंवा `ccam <name>` ने Claude Code सुरू करा. बहुतेक लोकांना एवढेच लागते.

## 1. तुमचे पहिले खाते जोडा

```sh
ccam add
```

विझार्ड `personal` किंवा `work` सारखे नाव विचारतो, आणि हे खाते तुमच्या इतर खात्यांशी कोणते आयटम शेअर करेल ते विचारतो. नंतर तो `claude auth login` चालवतो, जे साइन इन करण्यासाठी ब्राउझर उघडते. पूर्ण झाल्यावर तुम्हाला हे दिसते:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam स्वतःचे संदेश इंग्रजीमध्ये प्रिंट करते; कमांड आणि पर्याय प्रत्येक भाषेत सारखेच आहेत.

## 2. दुसरे खाते जोडा

```sh
ccam add work
```

ब्राउझरमध्ये दुसऱ्या खात्याने साइन इन करा. ते `claude2` बनते.

## 3. खाते बदला

यांपैकी कोणतीही कमांड work खात्याने Claude Code सुरू करते:

```sh
claude2
ccam work
ccam 2
```

खात्याच्या नावानंतरचे सर्व काही जसेच्या तसे Claude Code कडे जाते, त्यामुळे तुम्ही ओळखीचे पर्याय वापरत राहू शकता:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

फक्त `ccam` चालवल्यास निवडण्याची यादी उघडते: बाण की आणि Enter वापरा, किंवा खात्याचा क्रमांक दाबा.

## 4. काय शेअर करायचे ते निवडा

डीफॉल्टनुसार MCP सर्व्हर, `settings.json`, global `CLAUDE.md`, आणि agents, commands, skills आणि plugins खात्यांमध्ये शेअर होतात. संभाषण इतिहास शेअर होत नाही. खात्यांमध्ये संभाषणे पुढे चालू ठेवायची असतील तर दोन्ही खात्यांसाठी तो चालू करा:

```sh
ccam share personal --history on
ccam share work --history on
```

काहीही बदलण्यापूर्वी ccam काय हलवणार आहे ते दाखवते. तपशीलासाठी [शेअरिंग](/docs/sharing) पाहा.

## 5. खाते मर्यादेला पोहोचल्यावर

```sh
ccam usage
ccam next --continue
```

`ccam usage` सर्व खात्यांचा 5 तासांचा आणि साप्ताहिक वापर दाखवते. `ccam next --continue` तेच संभाषण पुढच्या खात्यात उघडते. [वापर मर्यादा](/docs/usage-limits) पाहा.

## आधीच ~/.claude वापरत आहात?

पुन्हा साइन इन करण्याऐवजी तो नोंदवा:

```sh
ccam import
```

फोल्डर आहे तिथेच राहतो, आणि साधे `claude` पूर्वीसारखेच काम करत राहते.
