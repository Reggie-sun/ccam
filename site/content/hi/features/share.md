---
title: ग्रुप के साथ साझा करना
description: ccam का ग्रुप चुने हुए आइटम अपने अकाउंट्स के बीच एक फ़ोल्डर से साझा करता है; ccam group योजना, बैकअप और टकराव विकल्पों के साथ डेटा move करता है।
updated: 2026-10-23
order: 4
translation: machine
sourceHash: e4f05b0fd572
---

हर अकाउंट का अपना फ़ोल्डर होता है, लेकिन उसकी ज़्यादातर सामग्री को अलग रखने की ज़रूरत नहीं होती। एक **ग्रुप** बताता है कि क्या साझा होगा और किनके बीच: ग्रुप का हर अकाउंट उन आइटम को एक ही फ़ोल्डर, `~/.local/share/ccam/groups/<group>/`, से लिंक करता है, इसलिए वे सब एक ही कॉपी पढ़ते और लिखते हैं। लॉगिन tokens कभी साझा नहीं होते।

एक अकाउंट कई ग्रुप में हो सकता है, लेकिन हर आइटम केवल एक ही ग्रुप से आता है। एक आम सेटअप में सेटिंग्स और एक्सटेंशन के लिए एक ग्रुप होता है, और दो अकाउंट्स के बीच बातचीत के इतिहास के लिए दूसरा:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

आपका जोड़ा पहला अकाउंट `default` बनाता है, और बाद के अकाउंट उसमें शामिल होते हैं, जब तक आप कुछ और न कहें।

## उपयोग

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

डेटा move करने वाली हर कमांड एक योजना प्रिंट करती है, ग्रुप की कॉपी से अलग हर फ़ाइल के बारे में पूछती है, और कुछ भी बदलने से पहले `Proceed? [y/N]` पूछती है:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

यही बदलाव [डैशबोर्ड](/features/dashboard) के **Groups** टैब पर भी हैं, और किसी अकाउंट पर `s` उसके ग्रुप चुनता है।

## आइटम

| आइटम | क्या कवर करता है |
| --- | --- |
| `mcp` | MCP सर्वर, ग्रुप की `mcp.json` के ज़रिए ([साझा MCP सर्वर](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | ग्लोबल `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` और `history.jsonl`, ताकि कोई दूसरा अकाउंट `--continue` कर सके |

जिस आइटम को अकाउंट का कोई ग्रुप साझा नहीं करता, वह अकाउंट के अपने फ़ोल्डर में रहता है।

## कमांड और विकल्प

| कमांड | क्या करती है |
| --- | --- |
| `ccam group` या `ccam group list` | ग्रुप्स की सूची दिखाएँ |
| `ccam group show <group>` | ग्रुप क्या साझा करता है, उसके अकाउंट, उसका फ़ोल्डर और MCP सर्वर |
| `ccam group add <group> --items <items> [--members <accounts>]` | ग्रुप बनाएँ; उसके अकाउंट्स के आइटम उसमें move होते हैं |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | ग्रुप क्या साझा करता है, यह उसके हर अकाउंट के लिए बदलें |
| `ccam group join <group> <account>...` | अकाउंट्स को ग्रुप में रखें |
| `ccam group leave <group> <account>...` | अकाउंट्स को बाहर निकालें; हर एक के पास एक कॉपी रहती है, MCP सर्वर समेत |
| `ccam group rename <group> <new-name>` | केवल config में नाम बदलें |
| `ccam group remove <group> [--leave]` | ग्रुप हटाएँ |

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: वह अकाउंट जिसकी कॉपी सबसे पहले ग्रुप में जाती है | डिफ़ॉल्ट अकाउंट, फिर नंबर के क्रम में |
| `--replace` | `join`: पहले उन ग्रुप्स को छोड़ें जो वही आइटम साझा करते हैं | बंद |
| `--fresh` | जो अकाउंट किसी आइटम को साझा करना बंद करते हैं, वे कॉपी के बजाय उसके बिना शुरू होते हैं (import किया गया फ़ोल्डर फिर भी एक रखता है) | बंद |
| `--leave` | `remove`: पहले ग्रुप के अकाउंट्स को बाहर निकालें | जब ग्रुप में अकाउंट हों तो ज़रूरी |
| `--prefer <group\|mine\|merge>` | टकराव कैसे तय हों | पूछें; `--yes` के साथ, हर फ़ाइल का डिफ़ॉल्ट |
| `-y`, `--yes` | बिना पूछे चलाएँ; टकराव `--prefer` के अनुसार तय होते हैं | बंद |
| `--dry-run` (ग्लोबल) | योजना प्रिंट करें, कुछ न बदलें | बंद |

अकाउंट नाम या नंबर से दिए जाते हैं; आइटम कॉमा से अलग की गई सूची के रूप में। `ccam add`, `--group <group>` (कई बार दिया जा सकता है) या `--no-group` लेता है।

## यह कैसे काम करता है

शामिल फ़ाइलें:

- `~/.config/ccam/config.toml`: `[[group]]` entries: `name`, `items`, `members` और, जब फ़ोल्डर का नाम अलग हो, `dir`।
- `~/.local/share/ccam/profiles/<name>/`: अकाउंट का `CLAUDE_CONFIG_DIR`।
- `~/.local/share/ccam/groups/<folder>/`: साझा कॉपियाँ।

जब `config.toml` पढ़ा जाता है, तो हर अकाउंट के ग्रुप्स को हर आइटम के लिए एक फ़ोल्डर में बदल दिया जाता है। पाँच अकाउंट्स के लिए इसमें लगभग 25 µs लगते हैं, इसलिए launch अब भी केवल `config.toml` पढ़ता है। ऐसा config जिसमें किसी अकाउंट को एक आइटम दो ग्रुप्स से मिलता, त्रुटि है, और उसके ठीक होने तक कुछ भी शुरू नहीं होता।

बदलाव config की एक कॉपी पर काम करता है:

1. कमांड कॉपी को बदलती है (`AddGroup`, `SetGroupItems`, `JoinGroup` …), जो नियम जाँचता है।
2. `internal/group` पुराने और नए की तुलना करता है और हर उस अकाउंट के लिए जिसके आइटम बदलते हैं, `share.Build` से एक योजना बनाता है। योजनाएँ एक के बाद एक बनती हैं, और हर योजना देखती है कि पहले वाली योजनाएँ ग्रुप फ़ोल्डरों में क्या रखेंगी, इसलिए केवल पहला अकाउंट अपनी कॉपी अंदर move करता है। पहली कॉपी ग्रुप फ़ोल्डर से आती है, अगर उसमें पहले से कोई हो, फिर `--from` से, फिर डिफ़ॉल्ट अकाउंट से, फिर बाक़ी से नंबर के क्रम में।
3. अगर इनमें से कोई अकाउंट चल रहा हो, चाहे Claude Code किसी भी तरह शुरू हुआ हो (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), तो ccam मना कर देता है, और Developer Mode के बिना Windows पर भी, अगर किसी फ़ाइल को लिंक की ज़रूरत हो।
4. टकराव prompt से या `--prefer` से तय होते हैं, योजना प्रिंट होती है, और ccam पूछता है।
5. सेव करने से ठीक पहले ccam फिर से जाँचता है: जो अकाउंट सवालों के इंतज़ार के दौरान शुरू हुआ, वह बदलाव रोक देता है, और उन्हीं ग्रुप्स का किसी दूसरे ccam का बदलाव भी (`the groups changed while this plan waited; nothing was moved, make the change again`)। कुछ भी move होने से पहले नया config सेव होता है, और बदलाव `~/.local/state/ccam/group-change.json` में दर्ज होता है। अगर कोई चरण बीच में रुक जाए, तो नोट बना रहता है: `ccam doctor --fix` बदलाव को उसके क्रम में, आपके दिए जवाबों के साथ, MCP चरणों समेत पूरा करता है, और तब तक ग्रुप का कोई और बदलाव शुरू नहीं होता (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``)।
6. MCP चरण पहले चलते हैं। जब अकाउंट किसी ग्रुप से एक्सटेंशन निकालते हैं, तो ग्रुप की plugin सूचियाँ जहाँ भी उनके पाथ किसी अकाउंट के फ़ोल्डर से होकर जाते हैं, ग्रुप फ़ोल्डर की ओर मोड़ दी जाती हैं (योजना में `group plugin lists:`)। फिर योजनाएँ क्रम से चलती हैं। हर चरण डेटा को ग्रुप में या बैकअप फ़ोल्डर में move करता है; केवल वे फ़ाइलें बैकअप के बजाय हटाई जाती हैं जो ग्रुप में byte-दर-byte पहले से मौजूद हैं।

`share.Layout` उन पाथ की अनुमत सूची है जिन्हें हर आइटम कवर करता है। जो इस सूची में नहीं है, वह कभी किसी ग्रुप तक नहीं पहुँच सकता। `share.Never` में `.credentials.json` और `.claude.json` के नाम हैं, और अगर इनमें से कोई कभी लेआउट में दिखे तो `Build` एक आंतरिक त्रुटि के साथ रुक जाता है। `mcp` आइटम का कोई पाथ नहीं है: सर्वर सूची `.claude.json` के अंदर रहती है, इसलिए [MCP साझा करना](/features/mcp) ग्रुप की `mcp.json` से काम करता है।

जब `side` इतिहास साझा करना शुरू करता है और ग्रुप में पहले से सामग्री है, तो एक फ़ोल्डर `projects/` के साथ यह होता है:

<!-- motion:share-on -->

अगर ग्रुप के पास उस पाथ पर अभी कुछ नहीं है, तो यह अकाउंट उसे साझा करने वाला पहला है और उसका फ़ोल्डर जैसा है वैसा वहाँ move हो जाता है। वरना हर वह फ़ाइल जो ग्रुप में नहीं है, ग्रुप में move होती है, और दोनों तरफ़ मौजूद फ़ाइलें अकाउंट के फ़ोल्डर में रहती हैं, इसलिए ग्रुप का वर्ज़न रखा जाता है। ग्रुप की फ़ाइलों से हूबहू मिलने वाली फ़ाइलें हटा दी जाती हैं, और योजना उनके नाम बताती है जो अलग हैं (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`)। फिर अकाउंट का फ़ोल्डर अकाउंट फ़ोल्डर के अंदर `.ccam-backup/<time>/` में जाता है, और उसकी जगह एक लिंक लग जाता है। अगर ग्रुप की entry ख़ुद शामिल होने वाले अकाउंट की अपनी entry का लिंक है, तो अकाउंट की सामग्री उसकी जगह ग्रुप में जाती है। बातचीत की फ़ाइलों के नाम उनके session ID पर होते हैं, इसलिए कई अकाउंट्स के इतिहास बिना टकराव के एक साथ जुड़ जाते हैं। `plugins/` की registries अकाउंट फ़ोल्डर से होकर जाने वाले absolute पाथ सेव करती हैं; ccam उन्हें ग्रुप के फ़ोल्डर की ओर मोड़ता है और अकाउंट की entries को उस registry में merge करता है जो ग्रुप के पास पहले से है।

ग्रुप छोड़ना लिंक को उलट देता है:

<!-- motion:share-off -->

ccam पहले ग्रुप की मौजूदा सामग्री अकाउंट फ़ोल्डर में `.ccam-backup/incomplete-<name>` में कॉपी करता है, और लिंक को केवल तब बदलता है जब कॉपी पूरी हो जाए, ताकि बीच में रुकी कॉपी लिंक को जगह पर छोड़ दे और `ccam doctor --fix` फिर से कॉपी करे (ऐसे पूरा करने की योजना `restore <entry>: … (copied back from the group before a stop)` दिखाती है)। यह लिंक तभी हटाता है जब वह अभी भी groups फ़ोल्डर के अंदर उसी तरह इशारा करता हो जैसे योजना बनते समय करता था; जो पाथ तब से बदल गया, जैसे लिंक की जगह असली फ़ाइल, उसे जस का तस छोड़ दिया जाता है और बदलाव रुक जाता है। `--fresh` के साथ कॉपी छोड़ दी जाती है। ग्रुप अपनी सामग्री दूसरे अकाउंट्स के लिए रखे रहता है। जो अकाउंट किसी आइटम को ऐसे ग्रुप में ले जाता है जिसका फ़ोल्डर अभी ख़ाली है, वह जो देख रहा था उसकी एक कॉपी साथ लाता है, इसलिए उसकी फ़ाइलें बनी रहती हैं; वह कॉपी पूरी होने तक `<group>/.ccam-incomplete-<name>` में तैयार रखी जाती है। जो अकाउंट MCP साझा करना बंद करता है, वह ग्रुप के सर्वर रखता है: वे उसकी अपनी `.claude.json` में जाते हैं (एक बैकअप रखा जाता है), उन्हीं नामों वाले उसके अपने सर्वरों के बगल में, जो बने रहते हैं। `--fresh` के साथ योजना इसके बजाय उन सर्वरों की सूची देती है जो लोड होना बंद कर देते हैं। फिर से जुड़ने पर वे कॉपियाँ हट जाती हैं जिन्हें किसी ने नहीं बदला। `extensions` छोड़ने पर ग्रुप की plugin सूचियाँ भी, जहाँ उनके पाथ छोड़ने वाले अकाउंट के फ़ोल्डर से होकर जाते थे, ग्रुप फ़ोल्डर की ओर मोड़ दी जाती हैं; दूसरे अकाउंट्स के लिए उनका मतलब वही रहता है।

दोनों तरफ़ अलग सामग्री के साथ मौजूद एक अकेली फ़ाइल टकराव है:

<!-- motion:share-conflict -->

आप `settings.json` जैसी अकेली फ़ाइलों के लिए चुनते हैं; फ़ोल्डरों के अंदर ग्रुप की कॉपी बनी रहती है, और योजना उन फ़ाइलों के नाम बताती है जो अलग हैं। `history.jsonl` पर `merge by line`, `keep group` और `use mine` मिलते हैं, merge पहले। दूसरी `.json` फ़ाइलों पर `keep group`, `use mine` और `merge JSON` मिलते हैं। बाकी सब पर `keep group` और `use mine`। `merge JSON` ऑब्जेक्ट को recursively merge करता है, arrays को बिना दोहराव के जोड़ता है, और जब दोनों तरफ़ एक ही key अलग मान पर सेट हो तो ग्रुप का मान रखता है, ताकि दूसरे अकाउंट्स को कोई बदलाव न दिखे। `keep group` को छोड़कर हर विकल्प पहले ग्रुप की कॉपी अकाउंट के `.ccam-backup/<time>/group/` में सेव करता है। आपकी अपनी फ़ाइल हमेशा बैकअप में जाती है।

हर बार शुरू होने पर, `share.Check` लेआउट के हर पाथ की तुलना अकाउंट के ग्रुप्स से केवल `lstat` और `readlink` से करता है। गायब लिंक तुरंत बना दिया जाता है। लिंक की जगह असली फ़ाइल, ग़लत जगह का लिंक, या पुराना लिंक शुरू करना रोक देता है और `ccam doctor --fix` की ओर इशारा करता है।

### नाम बदलना, हटाना और पुराना फ़ॉर्मैट

- `ccam group rename` केवल config बदलता है। फ़ोल्डर का नाम वही रहता है और ग्रुप उसे `dir` के रूप में दर्ज करता है, इसलिए कुछ भी move नहीं होता और कुछ भी खो नहीं सकता।
- `ccam group remove` अकाउंट वाले ग्रुप को तब तक हटाने से मना करता है जब तक आप `--leave` न जोड़ें। तब हर अकाउंट के पास एक कॉपी रहती है, और फ़ोल्डर `~/.local/share/ccam/trash/` में जाता है, जब तक कोई दूसरा ग्रुप उसे इस्तेमाल न करता हो।
- ग्रुप्स से पहले, `config.toml` हर अकाउंट के लिए एक `[profile.share]` table सेव करता था। ccam उस फ़ॉर्मैट को कुछ भी move किए बिना ग्रुप्स के रूप में पढ़ता है। पहली बार सेव करने पर यह नया फ़ॉर्मैट लिखता है और पुरानी फ़ाइल को `config.toml.bak` के रूप में रखता है। एक ही अकाउंट्स द्वारा साझा किए गए आइटम एक ग्रुप बन जाते हैं जिसका नाम फ़ोल्डर पर होता है; बाक़ी को `default-history` जैसे नाम मिलते हैं और वे वही फ़ोल्डर इस्तेमाल करते रहते हैं।

### Windows पर

Windows symbolic links केवल Developer Mode (Settings > System > For developers) में या किसी administrator के लिए बनाता है। यह अधिकार न होने पर, साझा फ़ोल्डर एक directory junction बन जाता है, जो Claude Code के लिए वैसे ही काम करता है। साझा फ़ाइल (`settings.json`, `CLAUDE.md`, `history.jsonl`) केवल symbolic link ही हो सकती है, इसलिए जिस ग्रुप बदलाव को इसकी ज़रूरत हो, उसे कुछ भी move होने से पहले मना कर दिया जाता है। इसी कारण वहाँ पहला `ccam add` केवल `mcp` और `extensions` के साथ ग्रुप `default` बनाता है और यह बताता है, और सेटिंग्स, memory या इतिहास साझा करने वाले ग्रुप में `ccam add --group` अस्वीकार होता है।

## सुरक्षा और विशेष स्थितियाँ

- **चल रहे अकाउंट।** जब Claude Code उन अकाउंट्स में से किसी के रूप में चल रहा हो जिन्हें बदलाव move करता है, तो बदलाव मना कर दिया जाता है, क्योंकि session उन फ़ाइलों में लिखता रहता है, चाहे उसे ccam ने शुरू किया हो या नहीं (सादा `claude`, कोई IDE, background sessions)। सेव करने से ठीक पहले जाँच दोहराई जाती है। छोड़ना और शामिल होना केवल नाम दिए गए अकाउंट से जुड़े हैं; `extensions` छोड़ना ग्रुप की plugin सूचियों में पाथ भी फिर से लिखता है, दूसरों के लिए उनका मतलब बदले बिना। `--dry-run` फिर भी योजना प्रिंट करता है।
- **बैकअप।** हर योजना को अकाउंट फ़ोल्डर में अपना `.ccam-backup/<time>/` फ़ोल्डर मिलता है। ccam इन्हें कभी नहीं मिटाता।
- **आपके अपने लिंक।** अगर कोई पाथ पहले से आपका बनाया लिंक है, उदाहरण के लिए किसी dotfiles repository का, और ग्रुप ख़ाली है, तो ccam लिंक के target को ग्रुप में कॉपी करता है और लिंक को ख़ुद बैकअप में move करता है। Target को छुआ नहीं जाता।
- **Import किए गए फ़ोल्डर।** `ccam import` से अपनी जगह पर रजिस्टर किया गया अकाउंट (जैसे `~/.claude`) अपनी जगह पर ही बदलता है, इसलिए सादा `claude` भी बदलाव देखता है; योजना यह बताती है। जब वह कोई ग्रुप छोड़ता है तो हमेशा एक कॉपी रखता है, `--fresh` के साथ भी, और `ccam remove` अकाउंट को भूलने से पहले उसे इसी तरह उसके ग्रुप्स से निकालता है। ग्रुप में रहते समय उसकी साझा फ़ाइलें ग्रुप फ़ोल्डर में रहती हैं, इसलिए `~/.local/share/ccam` मिटाने से पहले अकाउंट हटाएँ (या ग्रुप छोड़ें)।
- **एक फ़ोल्डर में दो अकाउंट।** इतिहास साझा करने वाले अकाउंट एक ही फ़ाइलों में लिखते हैं। जब दूसरा अकाउंट ऐसे फ़ोल्डर में शुरू होता है जहाँ एक पहले से चल रहा है, तो ccam चेतावनी देता है और पूछता है।
- **हाथ से किए बदलाव।** ऐसा सदस्य जो अकाउंट नहीं है, अनदेखा किया जाता है और `ccam doctor` उसकी सूचना देता है। `local`, `on`, `off` और `none` ग्रुप के नाम नहीं हो सकते।

## संबंधित

- [डैशबोर्ड](/features/dashboard)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [साझा MCP सर्वर](/features/mcp)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [सभी कमांड](/docs/commands)
