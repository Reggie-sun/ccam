---
title: क्विक स्टार्ट
description: "ccam add ने दोन Claude Code खाती जोडा, त्यांना ccwork सारखे शॉर्टकट द्या, खाती बदला, आणि ती काय शेअर करतील ते निवडा."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: b73914ff371f
---

प्रत्येक खाते `ccam add` ने जोडा, नंतर `ccam <name>` ने किंवा तुम्ही स्वतः नाव दिलेल्या शॉर्टकटने, जसे `ccwork`, Claude Code सुरू करा. बहुतेक लोकांना एवढेच लागते.

## 1. तुमचे पहिले खाते जोडा

```sh
ccam add
```

विझार्ड `personal` किंवा `work` सारखे नाव विचारतो, खाते तुमच्या इतर खात्यांशी कोणते आयटम शेअर करेल ते विचारतो, आणि शॉर्टकट कमांड तयार करायची का हेही विचारतो (पायरी 4 पाहा). नंतर तो `claude auth login` चालवतो, जे साइन इन करण्यासाठी ब्राउझर उघडते. पूर्ण झाल्यावर तुम्हाला हे दिसते:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

ccam स्वतःचे संदेश इंग्रजीमध्ये प्रिंट करते; कमांड आणि पर्याय प्रत्येक भाषेत सारखेच आहेत.

## 2. तुमची आधीची सेटिंग्ज आणा

तुम्ही आधी Claude Code वापरले असेल, तर तुमची सेटिंग्ज, `CLAUDE.md`, एक्स्टेंशन्स आणि MCP सर्व्हर `~/.claude` मध्ये आहेत. प्रत्येक खात्याला ती मिळावीत म्हणून ती शेअर केलेल्या सेटिंग्जमध्ये कॉपी करा:

```sh
ccam migrate
```

ccam आधी काय कॉपी करणार ते दाखवते आणि मग विचारते. ते `~/.claude` फक्त वाचते, कधीच लॉगिन कॉपी करत नाही, आणि तुम्ही `--history` जोडल्याशिवाय संभाषण इतिहास वगळते. तपशीलासाठी [आधीचा सेटअप कॉपी करणे](/docs/sharing#copy-an-existing-setup) पाहा.

## 3. दुसरे खाते जोडा

```sh
ccam add work
```

ब्राउझरमध्ये दुसऱ्या खात्याने साइन इन करा.

## 4. खाते बदला

`ccam work` Claude Code work खात्याने सुरू करते. याहून छोटे काही हवे असल्यास खात्याला त्याची स्वतःची शॉर्टकट कमांड द्या, जसे `ccwork`. शेअर करायच्या आयटमनंतर विझार्ड विचारतो:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` दाबा, नंतर नाव टाइप करा, किंवा सुचवलेले नाव वापरण्यासाठी Enter दाबा:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

लॉगिननंतर ccam याची खात्री देते:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

तुम्ही हा प्रश्न वगळला असेल, तर कधीही शॉर्टकट जोडू शकता:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

आता यांपैकी कोणतीही कमांड Claude Code work खात्याने सुरू करते:

```sh
ccwork
ccam work
ccam 2
```

फक्त `ccam` चालवल्यास निवडण्याची यादी उघडते: बाण की आणि Enter वापरा, किंवा खात्याचा क्रमांक दाबा.

शॉर्टकट किंवा खात्याच्या नावानंतरचे सर्व काही न बदलता Claude Code ला जाते, त्यामुळे तुम्ही ओळखीचे पर्याय वापरत राहता:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. काय शेअर करायचे ते निवडा

डीफॉल्टनुसार MCP सर्व्हर, `settings.json`, global `CLAUDE.md`, आणि agents, commands, skills आणि plugins खात्यांमध्ये शेअर होतात. संभाषण इतिहास शेअर होत नाही. खात्यांमध्ये संभाषणे पुढे चालू ठेवायची असतील तर दोन्ही खात्यांसाठी तो चालू करा:

```sh
ccam share personal --history on
ccam share work --history on
```

काहीही बदलण्यापूर्वी ccam काय हलवणार आहे ते दाखवते. तपशीलासाठी [शेअरिंग](/docs/sharing) पाहा.

## 6. खाते मर्यादेला पोहोचल्यावर

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
