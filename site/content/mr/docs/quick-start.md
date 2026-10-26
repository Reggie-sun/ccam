---
title: क्विक स्टार्ट
metaTitle: "क्विक स्टार्ट: Claude Code खाती जोडा आणि बदला"
description: "ccam add ने दोन Claude Code खाती जोडा, त्यांना ccwork सारखे शॉर्टकट द्या, खाती बदला, आणि ती काय शेअर करतील ते निवडा."
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

प्रत्येक खाते `ccam add` ने जोडा, नंतर `ccam <name>` ने किंवा तुम्ही स्वतः नाव दिलेल्या शॉर्टकटने, जसे `ccwork`, Claude Code सुरू करा. बहुतेक लोकांना एवढेच लागते.

आधीच Claude Code वापरता? सुरू करण्यापूर्वी या पानाच्या शेवटी असलेले [आधीच ~/.claude वापरत आहात?](#already-using-claude) वाचा.

## 1. तुमचे पहिले खाते जोडा

```sh
ccam add
```

विझार्ड `personal` किंवा `work` सारखे नाव विचारतो, खाते तुमच्या इतर खात्यांशी कोणते आयटम शेअर करेल ते विचारतो, आणि शॉर्टकट कमांड तयार करायची का हेही विचारतो (पायरी 3 पाहा). नंतर तो `claude auth login` चालवतो, जे साइन इन करण्यासाठी ब्राउझर उघडते. रिमोट मशीनवर तुमचा ब्राउझर शेवटी लोड न होणाऱ्या `http://localhost:…/callback` पेजवर पोहोचू शकतो; तो पत्ता टर्मिनलमध्ये पेस्ट करा आणि ccam साइन इन पूर्ण करते. साइन इन पूर्ण झाल्यावर तुम्हाला हे दिसते:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. तुमची आधीची सेटिंग्ज आणा

तुम्ही आधी Claude Code वापरले असेल, तर तुमची सेटिंग्ज, `CLAUDE.md` आणि एक्स्टेंशन्स `~/.claude` मध्ये आहेत, आणि तुमचे MCP सर्व्हर `~/.claude.json` मध्ये आहेत. प्रत्येक खात्याला ती मिळावीत म्हणून ती शेअर केलेल्या `default` गटात कॉपी करा:

```sh
ccam migrate
```

ccam आधी काय कॉपी करणार ते दाखवते आणि मग विचारते. ते `~/.claude` फक्त वाचते, कधीच लॉगिन कॉपी करत नाही, आणि तुम्ही `--history` जोडल्याशिवाय संभाषण इतिहास वगळते. तपशीलासाठी [आधीचा सेटअप कॉपी करणे](/docs/sharing#copy-an-existing-setup) पाहा.

## 3. दुसरे खाते जोडा

```sh
ccam add work
```

नाव तुम्ही कमांड लाइनवर दिल्यामुळे विझार्ड ते विचारत नाही, आणि `default` निवडलेला ठेवून सामील व्हायचे गट दाखवतो. या वेळी खात्याला त्याची स्वतःची शॉर्टकट कमांड द्या, जसे `ccwork`, म्हणजे तुम्ही ते एका शब्दाने सुरू करू शकाल. गटांनंतर विझार्ड विचारतो:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` दाबा, नंतर नाव टाइप करा, किंवा सुचवलेले नाव वापरण्यासाठी Enter दाबा:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

नंतर ब्राउझरमध्ये दुसऱ्या खात्याने साइन इन करा. ccam खाते आणि त्याचा शॉर्टकट यांची खात्री देते:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. खाते बदला

`ccam work` Claude Code work खात्याने सुरू करते. तुम्ही शॉर्टकटचा प्रश्न वगळला असेल, तर कधीही शॉर्टकट जोडू शकता:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

आता यांपैकी प्रत्येक कमांड Claude Code work खात्याने सुरू करते (`ccam list` मध्ये `2` हा त्याचा क्रमांक आहे):

```sh
ccwork
ccam work
ccam 2
```

फक्त `ccam` चालवल्यास [डॅशबोर्ड](/features/dashboard) उघडतो: बाण की किंवा खात्याच्या क्रमांकाने खात्यावर जा, मग Enter दाबा.

शॉर्टकट किंवा खात्याच्या नावानंतरचे सर्व काही न बदलता Claude Code ला जाते, त्यामुळे तुम्ही ओळखीचे पर्याय वापरत राहता:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. काय शेअर करायचे ते निवडा

डीफॉल्टनुसार MCP सर्व्हर, `settings.json`, global `CLAUDE.md` आणि एक्स्टेंशन्स (agents, commands, skills आणि plugins) `default` गटामार्फत खात्यांमध्ये शेअर होतात. संभाषण इतिहास शेअर होत नाही. खात्यांमध्ये संभाषणे पुढे चालू ठेवायची असतील तर तो गटात जोडा:

```sh
ccam group set default --add history
```

काहीही बदलण्यापूर्वी ccam प्रत्येक खात्यासाठी काय हलवणार आहे ते दाखवते. तपशीलासाठी [गटांसह शेअरिंग](/features/share) पाहा.

## 6. खाते मर्यादेला पोहोचल्यावर

```sh
ccam usage
ccam next --continue
```

`ccam usage` सर्व खात्यांचा 5 तासांचा आणि साप्ताहिक वापर दाखवते. `ccam next --continue` तेच संभाषण पुढच्या खात्यात उघडते. Scripts मध्ये `ccam --fallback work -p "…"` आपोआप खाते बदलते. [वापर मर्यादा](/docs/usage-limits) पाहा.

**Team आणि Enterprise सीट्स.** आपोआप होणारी खाते निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n` आणि, तुम्ही चालू केल्यावर, [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील सीट्सदरम्यानही खाते बदलते, आणि यामुळे Anthropic च्या अटींचे उल्लंघन होऊ शकते. त्यावर अवलंबून राहण्यापूर्वी [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) वाचा आणि तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा.

## आधीच ~/.claude वापरत आहात?

पुन्हा साइन इन करण्याऐवजी तो नोंदवा:

```sh
ccam import
```

फोल्डर आहे तिथेच राहतो, आणि साधे `claude` पूर्वीसारखेच काम करत राहते.
