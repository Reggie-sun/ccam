---
title: खात्यांमध्ये शेअरिंग
description: ccam share सेटिंग्ज, CLAUDE.md, एक्स्टेंशन्स आणि इतिहास योजना, बॅकअप आणि conflict पर्यायांसह गट फोल्डरशी कसे लिंक करते, आणि off कसे परत कॉपी करते.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 16de29bbdd57
---

प्रत्येक खात्याचा स्वतःचा फोल्डर असतो, पण त्यातील बहुतेक गोष्टी वेगळ्या असण्याची गरज नसते. `ccam share` एका वेळी एका खात्यासाठी आणि एका आयटमसाठी शेअरिंग चालू किंवा बंद करते: खात्याची फाइल किंवा फोल्डर एका गट फोल्डरमध्ये हलवला जातो आणि त्याच्या जागी symbolic link ठेवली जाते, त्यामुळे गटातील प्रत्येक खाते तीच प्रत वाचते आणि लिहिते. खाते काय शेअर करते ते बदलण्यासाठी, किंवा काही खाती त्यांच्या स्वतःच्या गटात ठेवण्यासाठी `ccam add` नंतर ते वापरा.

## वापर

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

आयटम पर्यायांशिवाय, `ccam share` सध्याची स्थिती छापते:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

पर्यायांसह, ccam योजना छापते, गटाच्या प्रतीपेक्षा वेगळ्या असलेल्या प्रत्येक फाइलबद्दल विचारते, आणि काहीही बदलण्यापूर्वी `Proceed? [y/N]` विचारते. नवीन खाते संभाषण इतिहास सोडून सर्व काही शेअर करते; [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing) प्रत्येक आयटमचे वापरकर्त्याच्या बाजूने वर्णन करते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP सर्व्हर | बदल नाही |
| `--settings <on\|off\|group>` | `settings.json` | बदल नाही |
| `--memory <on\|off\|group>` | global `CLAUDE.md` | बदल नाही |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | बदल नाही |
| `--history <on\|off\|group>` | `projects/` आणि `history.jsonl` | बदल नाही |
| `--fresh` | शेअरिंग बंद करताना, गटाचा मजकूर कॉपी करण्याऐवजी रिकाम्याने सुरुवात करा | बंद |
| `--prefer <group\|mine\|merge>` | conflicts कसे ठरवायचे | विचारणे; `--yes` सह प्रत्येक आयटमचा डीफॉल्ट |
| `-y`, `--yes` | न विचारता चालवा; conflicts `--prefer` पाळतात | बंद |
| `-h`, `--help` | `share` ची मदत | |
| `--dry-run` (global) | पाळला जातो: conflicts बद्दल विचारतो, योजना छापतो, काहीही बदलत नाही | बंद |

`on` मूल्य म्हणजे `default` गट, `off` म्हणजे local (शेअर न केलेले), आणि इतर कोणतेही मूल्य म्हणजे गटाचे नाव. गटांची नावे `~/.local/share/ccam/groups/` खाली फोल्डरची नावे बनतात, त्यामुळे `ccam share side --history ../x` `invalid group name: "../x"` देऊन अयशस्वी होते. खाते नावाने किंवा क्रमांकाने दिले जाते, आणि फक्त तुम्ही नाव दिलेले आयटम बदलतात.

## हे कसे काम करते

संबंधित फाइल्स:

- `~/.config/ccam/config.toml`: प्रत्येक खात्याचे `[profile.share]` table, प्रत्येक आयटमसाठी एक मूल्य.
- `~/.local/share/ccam/profiles/<name>/`: खात्याचा `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<group>/`: शेअर केलेल्या प्रती.

1. `internal/cli/share.go` config वाचते, प्रत्येक `--<item>` मूल्य `config.Share.Set` ने खात्याच्या सेटिंग्जच्या प्रतीवर लागू करते, आणि खाते चालू असेल तर नकार देते.
2. `share.Build` सध्याच्या सेटिंग्जची अपेक्षित सेटिंग्जशी तुलना करते आणि टप्प्यांची यादी करते. ते फक्त डिस्क वाचते: `share.Layout` मधील प्रत्येक पाथसाठी `os.Lstat` आणि `os.Readlink`. योजना म्हणजे operations ची यादी: गटातील रिकाम्या जागी हलवणे (adopt), गटाकडे नसलेल्या फाइल्स हलवणे (merge dir), सोडवलेला conflict गटात लिहिणे, बॅकअप, unlink, link, परत कॉपी करणे, plugin पाथ पुन्हा लिहिणे आणि plugin registry merge करणे.
3. conflicts prompt ने किंवा `--prefer` ने ठरवले जातात.
4. `Plan.Describe` योजना छापते, आणि ccam विचारते.
5. काहीही हलण्यापूर्वी अपेक्षित सेटिंग्ज `config.toml` मध्ये लिहिल्या जातात. `Plan.Apply` अर्ध्यावर थांबले, तर `ccam doctor --fix` आणि पुन्हा चालवणे सेव्ह केलेल्या अपेक्षित स्थितीकडे आपली योजना तयार करतात आणि बदल पूर्ण करतात.
6. `Plan.Apply` टप्पे क्रमाने चालवते. प्रत्येक टप्पा डेटा गटात किंवा बॅकअप फोल्डरमध्ये हलवतो; काहीही हटवले जात नाही.

`share.Layout` ही प्रत्येक आयटम कोणते पाथ व्यापतो त्याची अनुमत यादी आहे. तिच्यात नसलेले काहीही कधीच गटापर्यंत पोहोचू शकत नाही. `share.Never` `.credentials.json` आणि `.claude.json` यांची नावे देते, आणि यांपैकी कोणतीही layout मध्ये दिसलीच तर `Build` internal error सह थांबते. `mcp` आयटमला पाथ नाहीत: सर्व्हर यादी `.claude.json` च्या आत असते, म्हणून [MCP शेअरिंग](/features/mcp) त्याऐवजी गटाच्या `mcp.json` मार्फत चालते, आणि `--mcp` फक्त सेटिंग बदलतो.

`side` इतिहास चालू करते आणि गटाकडे आधीच मजकूर असतो तेव्हा `projects/` या एका फोल्डरचे असे होते:

<!-- motion:share-on -->

गटाकडे त्या पाथवर अजून काहीच नसेल, तर ते शेअर करणारे हे पहिले खाते असते आणि त्याचा फोल्डर जसाच्या तसा, बॅकअपशिवाय तिथे हलवला जातो. नाहीतर `mergeDir` गटाकडे नसलेली प्रत्येक फाइल गटात हलवते, आणि दोन्ही बाजूंना असलेल्या फाइल्स खात्याच्या फोल्डरमध्ये राहतात, त्यामुळे गटाची आवृत्ती ठेवली जाते. मग खात्याचा फोल्डर खात्याच्या फोल्डरमधील `.ccam-backup/<time>/` मध्ये जातो, आणि `os.Symlink` त्याच्या जागी लिंक ठेवते. संभाषण फाइल्सना त्यांच्या session ID वरून नावे दिली जातात, त्यामुळे अनेक खात्यांचे इतिहास कोणत्याही conflict शिवाय एकत्र होतात.

हलवण्यासाठी `rename(2)` वापरले जाते, आणि file systems वेगळ्या असल्यास कॉपी करून हटवणे वापरले जाते. हलवलेल्या tree मधील, त्याच्या बाहेर निर्देश करणारी relative link absolute link म्हणून पुन्हा लिहिली जाते (`fsutil.PortableLink`), त्यामुळे ती नवीन जागेवरूनही बरोबर जागी जाते. `plugins/` registries (`installed_plugins.json`, `known_marketplaces.json`) खात्याच्या फोल्डरमधून जाणारे absolute पाथ साठवतात; ccam त्यांना गटाच्या फोल्डरकडे निर्देशित करते आणि, गटाकडे आधीच registry असल्यास, खात्याच्या नोंदी त्यात merge करते.

आयटम बंद केल्याने लिंक उलट होते:

<!-- motion:share-off -->

लिंक groups फोल्डरमध्ये निर्देश करत असेल तरच ccam ती काढते, मग `CopyTree` ने गटाचा सध्याचा मजकूर खात्यात कॉपी करते. `--fresh` सह कॉपी वगळली जाते आणि खाते त्या पाथवर काहीही नसताना सुरू होते. गट इतर खात्यांसाठी आपला मजकूर ठेवतो. एखादा आयटम एका गटातून दुसऱ्या गटात नेल्यास तो unlink होऊन नवीन गटाशी लिंक होतो, आणि जुना गट आहे तसाच राहतो; योजना हे `Note:` ओळीत सांगते.

दोन्ही बाजूंना वेगळ्या मजकुरासह असलेली एकच फाइल म्हणजे conflict:

<!-- motion:share-conflict -->

`share.NewConflict` फाइलला लागू होणारे पर्याय देते. `history.jsonl` `merge by line`, `keep group` आणि `use mine` देते, merge पहिला. इतर `.json` फाइल्स `keep group`, `use mine` आणि `merge JSON` देतात. बाकी सर्व `keep group` आणि `use mine` देतात. पहिला पर्याय डीफॉल्ट असतो:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) objects recursively merge करते, arrays duplicates शिवाय जोडते, आणि दोन्ही बाजूंनी तीच key वेगळी सेट केली असेल तर गटाचे मूल्य ठेवते, त्यामुळे इतर खात्यांना वागण्यात कोणताही बदल दिसत नाही. `merge by line` गटाच्या ओळी ठेवते, मग गटाकडे नसलेल्या तुमच्या ओळी जोडते. `keep group` सोडून प्रत्येक पर्याय आधी गटाची प्रत खात्याच्या फोल्डरमधील `.ccam-backup/<time>/group/` मध्ये सेव्ह करतो, मग परिणाम `fsutil.AtomicWrite` ने लिहितो. तुमची स्वतःची फाइल नेहमी बॅकअपमध्ये जाते.

`--prefer group`, `mine` किंवा `merge` न विचारता प्रत्येक conflict ठरवतो; `merge` म्हणजे फाइल जे देते ते, `merge JSON` किंवा `merge by line`, आणि दोन्हीपैकी काहीच न देणारी फाइल डीफॉल्ट ठेवते. फक्त `--yes` प्रत्येक फाइलचा डीफॉल्ट घेतो.

प्रत्येक वेळी सुरू करताना, `share.Check` layout मधील प्रत्येक पाथची सेटिंगशी फक्त `lstat` आणि `readlink` वापरून तुलना करते, त्यामुळे ते सुरू करण्याच्या मार्गासाठी पुरेसे हलके राहते. नसलेली लिंक तिथल्या तिथे तयार होते, कारण त्यात गमावण्यासारखे काही नसते. लिंकच्या जागी खरी फाइल, चुकीच्या जागी जाणारी लिंक, किंवा शिळी लिंक सुरू करणे थांबवते:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` दोन्ही बाजूंना सेव्ह केलेल्या सेटिंग्ज देऊन `share.Build` चालवते, जे प्रत्येक टप्पा डिस्कवर जे आहे त्यावरून ठरवते, conflicts मध्ये गटाचा मजकूर ठेवते आणि ज्याच्या जागी काही ठेवते त्याचा बॅकअप घेते.

### Windows वर

Windows symbolic links फक्त Developer Mode (Settings > System > For developers) मध्ये किंवा administrator साठी तयार करते. `fsutil.Link` आधी symbolic link वापरून पाहते. तो अधिकार नसल्यास, शेअर केलेला फोल्डर (`agents`, `skills`, `projects` आणि इतर) directory junction बनतो, जे कोणताही user तयार करू शकतो, आणि Claude Code साठी तो तसाच चालतो. शेअर केलेली फाइल (`settings.json`, `CLAUDE.md`, `history.jsonl`) फक्त symbolic link च असू शकते, त्यामुळे ccam `Windows links files only in Developer Mode` ने थांबते आणि तो mode चालू करण्याचे किंवा तो आयटम local ठेवण्याचे सुचवते. `share.Check` junctions ना links म्हणून वाचते (`fsutil.LinkTarget`), आणि खात्याचा फोल्डर काढताना junction स्वतःच काढले जाते, ते ज्या गटाच्या मजकुराकडे निर्देश करते तो कधीच नाही.

## सुरक्षितता आणि अपवादात्मक स्थिती

- **चालू खाती.** खात्याने Claude Code चालू असताना ccam शेअरिंग बदलण्यास नकार देते (`account work is running; quit it first, then change its sharing`), कारण session हलवल्या जाणाऱ्या फाइल्समध्ये लिहीत राहते आणि दरम्यान जोडलेल्या ओळी हरवतील. `--dry-run` तरीही योजना छापतो.
- **बॅकअप.** प्रत्येक योजनेला खात्याच्या फोल्डरमध्ये स्वतःचा `.ccam-backup/<time>/` फोल्डर मिळतो, आणि दोन योजना एकाच सेकंदात सुरू झाल्यास क्रमांक असलेला suffix. ccam ते कधीच हटवत नाही.
- **अर्धवट अपयश.** एखादा टप्पा अयशस्वी झाला, तर ccam त्रुटी आणि मूळ फाइल्स कुठे आहेत ते छापते, आणि `ccam doctor --fix` ने पूर्ण करायला सांगते. अपेक्षित स्थिती आधीच सेव्ह असते, त्यामुळे दुरुस्ती तिच्याकडेच पुढे जाते.
- **तुमच्या स्वतःच्या लिंक्स.** एखादा पाथ आधीच तुम्ही केलेली लिंक असेल, उदाहरणार्थ dotfiles repository कडे, आणि गट रिकामा असेल, तर ccam लिंकचे target गटात कॉपी करते आणि लिंक स्वतः बॅकअपमध्ये हलवते. target ला हात लावला जात नाही.
- **लॉगिन फाइल्स.** `.credentials.json` आणि `.claude.json` कधीच गटाशी लिंक, गटात हलवल्या किंवा कॉपी केल्या जात नाहीत.
- **एकाच फोल्डरमध्ये दोन खाती.** इतिहास शेअर करणारी खाती त्याच फाइल्समध्ये लिहितात. ज्या फोल्डरमध्ये त्याच इतिहास गटातील एक खाते आधीच चालू आहे तिथे दुसरे सुरू झाले, तर ccam इशारा देते आणि पुढे जाण्यापूर्वी विचारते.

## संबंधित

- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [शेअर केलेले MCP सर्व्हर](/features/mcp)
- [आधीचा सेटअप कॉपी करा](/features/migrate)
- [सर्व कमांड्स](/docs/commands)
