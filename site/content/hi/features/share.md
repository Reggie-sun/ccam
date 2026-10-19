---
title: अकाउंट्स के बीच साझा करना
description: ccam share कैसे सेटिंग्स, CLAUDE.md, एक्सटेंशन और इतिहास को योजना, बैकअप और टकराव के विकल्पों के साथ एक ग्रुप फ़ोल्डर से लिंक करता है, और off कैसे वापस कॉपी करता है।
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

हर अकाउंट का अपना फ़ोल्डर होता है, लेकिन उसकी ज़्यादातर सामग्री को अलग रखने की ज़रूरत नहीं होती। `ccam share` एक बार में एक अकाउंट और एक आइटम के लिए साझा करना चालू या बंद करता है: अकाउंट की फ़ाइल या फ़ोल्डर एक ग्रुप फ़ोल्डर में move होता है और उसकी जगह एक symbolic link लग जाता है, इसलिए ग्रुप का हर अकाउंट एक ही कॉपी पढ़ता और लिखता है। `ccam add` के बाद यह बदलने के लिए कि कोई अकाउंट क्या साझा करता है, या कुछ अकाउंट्स को उनके अपने ग्रुप में रखने के लिए इसे इस्तेमाल करें।

## उपयोग

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

आइटम विकल्पों के बिना, `ccam share` मौजूदा स्थिति प्रिंट करता है:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

विकल्पों के साथ, ccam एक योजना प्रिंट करता है, ग्रुप की कॉपी से अलग हर फ़ाइल के बारे में पूछता है, और कुछ भी बदलने से पहले `Proceed? [y/N]` पूछता है। नया अकाउंट बातचीत के इतिहास को छोड़कर सब कुछ साझा करता है; [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing) हर आइटम को उपयोगकर्ता की नज़र से बताता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP सर्वर | बिना बदलाव |
| `--settings <on\|off\|group>` | `settings.json` | बिना बदलाव |
| `--memory <on\|off\|group>` | ग्लोबल `CLAUDE.md` | बिना बदलाव |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | बिना बदलाव |
| `--history <on\|off\|group>` | `projects/` और `history.jsonl` | बिना बदलाव |
| `--fresh` | साझा करना बंद करते समय, ग्रुप की सामग्री कॉपी करने के बजाय ख़ाली शुरुआत करें | बंद |
| `--prefer <group\|mine\|merge>` | टकराव कैसे तय हों | पूछें; `--yes` के साथ, हर आइटम का डिफ़ॉल्ट |
| `-y`, `--yes` | बिना पूछे चलाएँ; टकराव `--prefer` के अनुसार तय होते हैं | बंद |
| `-h`, `--help` | `share` की help | |
| `--dry-run` (ग्लोबल) | माना जाता है: टकराव के बारे में पूछता है, योजना प्रिंट करता है, कुछ नहीं बदलता | बंद |

मान `on` का मतलब `default` ग्रुप है, `off` का मतलब local (साझा नहीं), और कोई भी दूसरा मान ग्रुप का नाम है। ग्रुप के नाम `~/.local/share/ccam/groups/` के नीचे फ़ोल्डर के नाम बनते हैं, इसलिए `ccam share side --history ../x`, `invalid group name: "../x"` के साथ विफल होता है। अकाउंट नाम या नंबर से दिया जाता है, और केवल वही आइटम बदलते हैं जिनका आप नाम देते हैं।

## यह कैसे काम करता है

शामिल फ़ाइलें:

- `~/.config/ccam/config.toml`: हर अकाउंट की `[profile.share]` table, हर आइटम के लिए एक मान।
- `~/.local/share/ccam/profiles/<name>/`: अकाउंट का `CLAUDE_CONFIG_DIR`।
- `~/.local/share/ccam/groups/<group>/`: साझा कॉपियाँ।

1. `internal/cli/share.go` config पढ़ता है, `config.Share.Set` से हर `--<item>` मान को अकाउंट की सेटिंग्स की एक कॉपी पर लागू करता है, और अकाउंट चल रहा हो तो मना कर देता है।
2. `share.Build` मौजूदा सेटिंग्स की तुलना लक्ष्य से करता है और चरणों की सूची बनाता है। यह डिस्क को केवल पढ़ता है: `share.Layout` के हर पाथ के लिए `os.Lstat` और `os.Readlink`। योजना operations की एक सूची है: ग्रुप की ख़ाली जगह में move करना (adopt), वे फ़ाइलें move करना जो ग्रुप में नहीं हैं (merge dir), हल किए गए टकराव को ग्रुप में लिखना, बैकअप लेना, unlink करना, link करना, वापस कॉपी करना, plugin पाथ फिर से लिखना और plugin registry merge करना।
3. टकराव prompt से या `--prefer` से तय होते हैं।
4. `Plan.Describe` योजना प्रिंट करता है, और ccam पूछता है।
5. कुछ भी move होने से पहले लक्ष्य सेटिंग्स `config.toml` में लिखी जाती हैं। अगर `Plan.Apply` बीच में रुक जाए, तो `ccam doctor --fix` और दोबारा चलाना अपनी योजना सेव किए गए लक्ष्य की ओर बनाते हैं और बदलाव पूरा करते हैं।
6. `Plan.Apply` चरणों को क्रम से चलाता है। हर चरण डेटा को ग्रुप में या बैकअप फ़ोल्डर में move करता है; कुछ भी मिटाया नहीं जाता।

`share.Layout` उन पाथ की अनुमत सूची है जिन्हें हर आइटम कवर करता है। जो इस सूची में नहीं है, वह कभी किसी ग्रुप तक नहीं पहुँच सकता। `share.Never` में `.credentials.json` और `.claude.json` के नाम हैं, और अगर इनमें से कोई कभी लेआउट में दिखे तो `Build` एक आंतरिक त्रुटि के साथ रुक जाता है। `mcp` आइटम का कोई पाथ नहीं है: सर्वर सूची `.claude.json` के अंदर रहती है, इसलिए [MCP साझा करना](/features/mcp) इसके बजाय ग्रुप की `mcp.json` से काम करता है, और `--mcp` केवल सेटिंग बदलता है।

जब `side` इतिहास चालू करता है और ग्रुप में पहले से सामग्री है, तो एक फ़ोल्डर `projects/` के साथ यह होता है:

<!-- motion:share-on -->

अगर ग्रुप के पास उस पाथ पर अभी कुछ नहीं है, तो यह अकाउंट उसे साझा करने वाला पहला है और उसका फ़ोल्डर बिना बैकअप के, जैसा है वैसा वहाँ move हो जाता है। वरना `mergeDir` हर वह फ़ाइल ग्रुप में move करता है जो ग्रुप में नहीं है, और दोनों तरफ़ मौजूद फ़ाइलें अकाउंट के फ़ोल्डर में रहती हैं, इसलिए ग्रुप का वर्ज़न रखा जाता है। फिर अकाउंट का फ़ोल्डर अकाउंट फ़ोल्डर के अंदर `.ccam-backup/<time>/` में जाता है, और `os.Symlink` उसकी जगह एक लिंक लगाता है। बातचीत की फ़ाइलों के नाम उनके session ID पर होते हैं, इसलिए कई अकाउंट्स के इतिहास बिना टकराव के एक साथ जुड़ जाते हैं।

Move करने के लिए `rename(2)` इस्तेमाल होता है, और अलग file systems के बीच यह कॉपी करके मिटाने पर लौट आता है। Move किए गए ट्री के अंदर का ऐसा relative link जो उसके बाहर इशारा करता था, absolute link के रूप में फिर से लिखा जाता है (`fsutil.PortableLink`), ताकि वह नई जगह से भी सही जगह पहुँचे। `plugins/` की registries (`installed_plugins.json`, `known_marketplaces.json`) अकाउंट फ़ोल्डर से होकर जाने वाले absolute पाथ सेव करती हैं; ccam उन्हें ग्रुप के फ़ोल्डर की ओर मोड़ता है, और जब ग्रुप के पास पहले से registry हो, तो अकाउंट की entries उसमें merge करता है।

किसी आइटम को बंद करना लिंक को उलट देता है:

<!-- motion:share-off -->

ccam लिंक तभी हटाता है जब वह groups फ़ोल्डर के अंदर इशारा करता हो, फिर `CopyTree` से ग्रुप की मौजूदा सामग्री अकाउंट में कॉपी करता है। `--fresh` के साथ कॉपी छोड़ दी जाती है और अकाउंट उस पाथ पर कुछ नहीं के साथ शुरू होता है। ग्रुप अपनी सामग्री दूसरे अकाउंट्स के लिए रखे रहता है। किसी आइटम को एक ग्रुप से दूसरे में ले जाना उसे unlink करके नए ग्रुप से link करता है, और पुराने ग्रुप को जैसा है वैसा छोड़ देता है; योजना यह एक `Note:` पंक्ति में बताती है।

दोनों तरफ़ अलग सामग्री के साथ मौजूद एक अकेली फ़ाइल टकराव है:

<!-- motion:share-conflict -->

`share.NewConflict` फ़ाइल के लिए उपयुक्त विकल्पों की सूची देता है। `history.jsonl` पर `merge by line`, `keep group` और `use mine` मिलते हैं, merge पहले। दूसरी `.json` फ़ाइलों पर `keep group`, `use mine` और `merge JSON` मिलते हैं। बाकी सब पर `keep group` और `use mine`। पहला विकल्प डिफ़ॉल्ट है:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) ऑब्जेक्ट को recursively merge करता है, arrays को बिना दोहराव के जोड़ता है, और जब दोनों तरफ़ एक ही key अलग मान पर सेट हो तो ग्रुप का मान रखता है, ताकि दूसरे अकाउंट्स को व्यवहार में कोई बदलाव न दिखे। `merge by line` ग्रुप की पंक्तियाँ रखता है, फिर आपकी वे पंक्तियाँ जोड़ता है जो उसमें नहीं हैं। `keep group` को छोड़कर हर विकल्प पहले ग्रुप की कॉपी अकाउंट फ़ोल्डर के `.ccam-backup/<time>/group/` में सेव करता है, फिर नतीजा `fsutil.AtomicWrite` से लिखता है। आपकी अपनी फ़ाइल हमेशा बैकअप में जाती है।

`--prefer group`, `mine` या `merge` हर टकराव बिना पूछे तय करता है; `merge` का मतलब `merge JSON` या `merge by line` है, जो भी फ़ाइल पर उपलब्ध हो, और जिस फ़ाइल पर दोनों में से कोई नहीं है वह डिफ़ॉल्ट रखती है। अकेला `--yes` हर फ़ाइल का डिफ़ॉल्ट लेता है।

हर बार शुरू होने पर, `share.Check` लेआउट के हर पाथ की तुलना सेटिंग से केवल `lstat` और `readlink` से करता है, जिससे यह शुरू करने के रास्ते के लिए काफ़ी हल्का रहता है। गायब लिंक तुरंत बना दिया जाता है, क्योंकि उसमें खोने को कुछ नहीं है। लिंक की जगह असली फ़ाइल, ग़लत जगह का लिंक, या पुराना लिंक शुरू करना रोक देता है:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` दोनों तरफ़ सेव की गई सेटिंग्स के साथ `share.Build` चलाता है, जो हर चरण डिस्क पर मौजूद चीज़ों से तय करता है, टकराव में ग्रुप की सामग्री रखता है और जो कुछ बदलता है उसका बैकअप लेता है।

## सुरक्षा और विशेष स्थितियाँ

- **चल रहे अकाउंट।** जब Claude Code उस अकाउंट के रूप में चल रहा हो, तो ccam साझा करना बदलने से मना करता है (`account work is running; quit it first, then change its sharing`), क्योंकि session move हो रही फ़ाइलों में लिखता रहता है और इस बीच जोड़ी गई पंक्तियाँ खो जातीं। `--dry-run` फिर भी योजना प्रिंट करता है।
- **बैकअप।** हर योजना को अकाउंट फ़ोल्डर में अपना `.ccam-backup/<time>/` फ़ोल्डर मिलता है, और जब दो योजनाएँ एक ही सेकंड में शुरू हों तो नंबर वाला suffix जुड़ता है। ccam इन्हें कभी नहीं मिटाता।
- **आंशिक विफलता।** अगर कोई चरण विफल हो, तो ccam त्रुटि और मूल फ़ाइलों की जगह प्रिंट करता है, और `ccam doctor --fix` से काम पूरा करने को कहता है। लक्ष्य पहले से सेव है, इसलिए मरम्मत उसी की ओर आगे बढ़ती है।
- **आपके अपने लिंक।** अगर कोई पाथ पहले से आपका बनाया लिंक है, उदाहरण के लिए किसी dotfiles repository का, और ग्रुप ख़ाली है, तो ccam लिंक के target को ग्रुप में कॉपी करता है और लिंक को ख़ुद बैकअप में move करता है। Target को छुआ नहीं जाता।
- **लॉगिन फ़ाइलें।** `.credentials.json` और `.claude.json` कभी किसी ग्रुप में link, move या copy नहीं होतीं।
- **एक फ़ोल्डर में दो अकाउंट।** इतिहास साझा करने वाले अकाउंट एक ही फ़ाइलों में लिखते हैं। जब उसी इतिहास ग्रुप का दूसरा अकाउंट ऐसे फ़ोल्डर में शुरू होता है जहाँ एक पहले से चल रहा है, तो ccam चेतावनी देता है और आगे बढ़ने से पहले पूछता है।

## संबंधित

- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [साझा MCP सर्वर](/features/mcp)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [सभी कमांड](/docs/commands)
