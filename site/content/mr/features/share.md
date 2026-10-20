---
title: गटांसह शेअरिंग
description: ccam गट निवडलेले आयटम एका फोल्डरमार्फत त्याच्या खात्यांमध्ये शेअर करतो; ccam group योजना, बॅकअप आणि conflict पर्यायांसह डेटा हलवते.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

प्रत्येक खात्याचा स्वतःचा फोल्डर असतो, पण त्यातील बहुतेक गोष्टी वेगळ्या असण्याची गरज नसते. **गट** काय शेअर होते आणि कोण शेअर करते ते ठरवतो: गटातील प्रत्येक खाते ते आयटम एका फोल्डरशी, `~/.local/share/ccam/groups/<group>/`, लिंक करते, त्यामुळे ती सर्व तीच प्रत वाचतात आणि लिहितात. लॉगिन tokens कधीच शेअर होत नाहीत.

एक खाते अनेक गटांमध्ये असू शकते, पण प्रत्येक आयटम फक्त एकाच गटातून येतो. नेहमीचा सेटअप म्हणजे सेटिंग्ज आणि एक्स्टेंशन्ससाठी एक गट, आणि दोन खात्यांमधील संभाषण इतिहासासाठी दुसरा:

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

तुम्ही जोडलेले पहिले खाते `default` तयार करते, आणि तुम्ही वेगळे सांगितले नाही तर नंतरची खाती त्यात सामील होतात.

## वापर

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

डेटा हलवणारी प्रत्येक कमांड योजना छापते, गटाच्या प्रतीपेक्षा वेगळ्या असलेल्या प्रत्येक फाइलबद्दल विचारते, आणि काहीही बदलण्यापूर्वी `Proceed? [y/N]` विचारते:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

हेच बदल [डॅशबोर्ड](/features/dashboard) च्या **Groups** tab वर आहेत, आणि खात्यावर `s` दाबल्यास त्याचे गट निवडता येतात.

## आयटम

| आयटम | काय व्यापतो |
| --- | --- |
| `mcp` | MCP सर्व्हर, गटाच्या `mcp.json` मार्फत ([शेअर केलेले MCP सर्व्हर](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | global `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` आणि `history.jsonl`, जेणेकरून दुसरे खाते `--continue` करू शकेल |

खात्याचा कोणताही गट शेअर करत नसलेला आयटम खात्याच्या स्वतःच्या फोल्डरमध्ये राहतो.

## कमांड्स आणि पर्याय

| कमांड | काय करते |
| --- | --- |
| `ccam group` किंवा `ccam group list` | गटांची यादी दाखवा |
| `ccam group show <group>` | गट काय शेअर करतो, त्याची खाती, त्याचा फोल्डर आणि MCP सर्व्हर |
| `ccam group add <group> --items <items> [--members <accounts>]` | गट तयार करा; त्याच्या खात्यांचे आयटम त्यात हलतात |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | गट काय शेअर करतो ते बदला, त्यातील प्रत्येक खात्यासाठी |
| `ccam group join <group> <account>...` | खाती गटात टाका |
| `ccam group leave <group> <account>...` | खाती बाहेर काढा; प्रत्येक खाते एक प्रत ठेवते |
| `ccam group rename <group> <new-name>` | फक्त config मध्ये नाव बदला |
| `ccam group remove <group> [--leave]` | गट काढून टाका |

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: ज्या खात्याची प्रत आधी गटात जाते ते खाते | डीफॉल्ट खाते, मग क्रमांकानुसार |
| `--replace` | `join`: तेच आयटम शेअर करणाऱ्या गटांमधून आधी बाहेर पडा | बंद |
| `--fresh` | एखादा आयटम शेअर करणे थांबवणारी खाती प्रतीशिवाय, रिकाम्याने सुरू होतात | बंद |
| `--leave` | `remove`: आधी गटाची खाती बाहेर काढा | गटात खाती असल्यास आवश्यक |
| `--prefer <group\|mine\|merge>` | conflicts कसे ठरवायचे | विचारणे; `--yes` सह प्रत्येक फाइलचा डीफॉल्ट |
| `-y`, `--yes` | न विचारता चालवा; conflicts `--prefer` पाळतात | बंद |
| `--dry-run` (global) | योजना छापा, काहीही बदलू नका | बंद |

खाती नावाने किंवा क्रमांकाने दिली जातात; आयटम स्वल्पविरामाने वेगळ्या केलेल्या यादीत. `ccam add` `--group <group>` (अनेकदा देता येतो) किंवा `--no-group` घेते.

## हे कसे काम करते

संबंधित फाइल्स:

- `~/.config/ccam/config.toml`: `[[group]]` नोंदी: `name`, `items`, `members` आणि, फोल्डरचे नाव वेगळे असल्यास, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: खात्याचा `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<folder>/`: शेअर केलेल्या प्रती.

`config.toml` वाचताना प्रत्येक खात्याचे गट प्रत्येक आयटमसाठी एका फोल्डरमध्ये ठरवले जातात. पाच खात्यांसाठी याला सुमारे 25 µs लागतात, त्यामुळे सुरू करताना अजूनही फक्त `config.toml` वाचली जाते. ज्या config मध्ये एखाद्या खात्याला एकच आयटम दोन गटांतून मिळेल ती त्रुटी आहे, आणि ती दुरुस्त होईपर्यंत काहीही सुरू होत नाही.

बदल config च्या प्रतीवर काम करतो:

1. कमांड प्रत संपादित करते (`AddGroup`, `SetGroupItems`, `JoinGroup` …), जी नियम तपासते.
2. `internal/group` जुन्या आणि नवीनची तुलना करते, आणि ज्या खात्यांचे आयटम बदलतात त्या प्रत्येकासाठी `share.Build` ने योजना तयार करते. योजना एकामागून एक तयार होतात, आणि प्रत्येक योजनेला आधीच्या योजना गट फोल्डर्समध्ये काय ठेवतील ते दिसते, त्यामुळे फक्त पहिले खाते आपली प्रत आत हलवते. पहिली प्रत गट फोल्डरकडे आधीच असल्यास तिथून येते, मग `--from`, मग डीफॉल्ट खाते, मग बाकीची खाती क्रमांकानुसार.
3. त्यांपैकी एखादे खाते चालू असल्यास ccam नकार देते, आणि Developer Mode शिवाय Windows वर एखाद्या फाइलला लिंक लागणार असल्यासही.
4. conflicts prompt ने किंवा `--prefer` ने ठरवले जातात, योजना छापली जाते, आणि ccam विचारते.
5. काहीही हलण्यापूर्वी नवीन config सेव्ह होते. एखादा टप्पा अर्ध्यावर थांबला, तर `ccam doctor --fix` आणि पुन्हा चालवणे सेव्ह केलेल्या गटांकडे आपली योजना तयार करतात आणि बदल पूर्ण करतात.
6. योजना क्रमाने चालतात. प्रत्येक टप्पा डेटा गटात किंवा बॅकअप फोल्डरमध्ये हलवतो; काहीही हटवले जात नाही.

`share.Layout` ही प्रत्येक आयटम कोणते पाथ व्यापतो त्याची अनुमत यादी आहे. तिच्यात नसलेले काहीही कधीच गटापर्यंत पोहोचू शकत नाही. `share.Never` `.credentials.json` आणि `.claude.json` यांची नावे देते, आणि यांपैकी कोणतीही layout मध्ये दिसलीच तर `Build` internal error सह थांबते. `mcp` आयटमला पाथ नाहीत: सर्व्हर यादी `.claude.json` च्या आत असते, म्हणून [MCP शेअरिंग](/features/mcp) गटाच्या `mcp.json` मार्फत चालते.

`side` इतिहास शेअर करू लागते आणि गटाकडे आधीच मजकूर असतो तेव्हा `projects/` या एका फोल्डरचे असे होते:

<!-- motion:share-on -->

गटाकडे त्या पाथवर अजून काहीच नसेल, तर ते शेअर करणारे हे पहिले खाते असते आणि त्याचा फोल्डर जसाच्या तसा तिथे हलवला जातो. नाहीतर गटाकडे नसलेली प्रत्येक फाइल गटात हलवली जाते, आणि दोन्ही बाजूंना असलेल्या फाइल्स खात्याच्या फोल्डरमध्ये राहतात, त्यामुळे गटाची आवृत्ती ठेवली जाते. मग खात्याचा फोल्डर खात्याच्या फोल्डरमधील `.ccam-backup/<time>/` मध्ये जातो, आणि त्याच्या जागी लिंक येते. संभाषण फाइल्सना त्यांच्या session ID वरून नावे दिली जातात, त्यामुळे अनेक खात्यांचे इतिहास कोणत्याही टक्करीशिवाय एकत्र होतात. `plugins/` registries खात्याच्या फोल्डरमधून जाणारे absolute पाथ साठवतात; ccam त्यांना गटाच्या फोल्डरकडे निर्देशित करते आणि गटाकडे आधीच असलेल्या registry मध्ये खात्याच्या नोंदी merge करते.

गटातून बाहेर पडल्याने लिंक उलट होते:

<!-- motion:share-off -->

लिंक groups फोल्डरमध्ये निर्देश करत असेल तरच ccam ती काढते, मग गटाचा सध्याचा मजकूर खात्यात कॉपी करते. `--fresh` सह कॉपी वगळली जाते. गट इतर खात्यांसाठी आपला मजकूर ठेवतो. ज्या गटाचा फोल्डर अजून रिकामा आहे त्या गटात एखादा आयटम नेणारे खाते त्याला दिसत होते त्याची प्रत सोबत आणते, त्यामुळे त्याच्या फाइल्स राहतात; MCP सर्व्हरही याच प्रकारे सोबत येतात. MCP शेअर करणे थांबवणारे खाते ccam ने त्याच्या `.claude.json` मध्ये कॉपी केलेले सर्व्हर गमावते (बॅकअप ठेवला जातो), कारण ते गटाकडून आले होते.

दोन्ही बाजूंना वेगळ्या मजकुरासह असलेली एकच फाइल म्हणजे conflict:

<!-- motion:share-conflict -->

`history.jsonl` `merge by line`, `keep group` आणि `use mine` देते, merge पहिला. इतर `.json` फाइल्स `keep group`, `use mine` आणि `merge JSON` देतात. बाकी सर्व `keep group` आणि `use mine` देतात. `merge JSON` objects recursively merge करते, arrays duplicates शिवाय जोडते, आणि दोन्ही बाजूंनी तीच key वेगळी सेट केली असेल तर गटाचे मूल्य ठेवते, त्यामुळे इतर खात्यांना कोणताही बदल दिसत नाही. `keep group` सोडून प्रत्येक पर्याय आधी गटाची प्रत खात्याच्या `.ccam-backup/<time>/group/` मध्ये सेव्ह करतो. तुमची स्वतःची फाइल नेहमी बॅकअपमध्ये जाते.

प्रत्येक वेळी सुरू करताना, `share.Check` layout मधील प्रत्येक पाथची खात्याच्या गटांशी फक्त `lstat` आणि `readlink` वापरून तुलना करते. नसलेली लिंक तिथल्या तिथे तयार होते. लिंकच्या जागी खरी फाइल, चुकीच्या जागी जाणारी लिंक, किंवा शिळी लिंक सुरू करणे थांबवते आणि `ccam doctor --fix` कडे निर्देश करते.

### नाव बदलणे, काढणे आणि जुना format

- `ccam group rename` फक्त config बदलते. फोल्डरचे नाव तेच राहते आणि गट ते `dir` म्हणून नोंदवतो, त्यामुळे काहीही हलत नाही आणि काहीही हरवू शकत नाही.
- `ccam group remove` खाती असलेला गट नाकारते, जोपर्यंत तुम्ही `--leave` जोडत नाही. मग प्रत्येक खाते एक प्रत ठेवते, आणि दुसरा गट तो फोल्डर वापरत नसेल तर तो `~/.local/share/ccam/trash/` मध्ये जातो.
- गट येण्यापूर्वी `config.toml` प्रत्येक खात्यासाठी `[profile.share]` table साठवत असे. ccam तो format काहीही न हलवता गट म्हणून वाचते. पहिल्या सेव्हच्या वेळी ते नवीन format लिहिते आणि जुनी फाइल `config.toml.bak` म्हणून ठेवते. तीच खाती शेअर करत असलेले आयटम फोल्डरच्या नावाचा एक गट बनतात; बाकीच्यांना `default-history` सारखी नावे मिळतात आणि ते तोच फोल्डर वापरत राहतात.

### Windows वर

Windows symbolic links फक्त Developer Mode (Settings > System > For developers) मध्ये किंवा administrator साठी तयार करते. तो अधिकार नसल्यास, शेअर केलेला फोल्डर directory junction बनतो, जो Claude Code साठी तसाच चालतो. शेअर केलेली फाइल (`settings.json`, `CLAUDE.md`, `history.jsonl`) फक्त symbolic link च असू शकते, त्यामुळे तिची गरज असलेला गट बदल काहीही हलण्यापूर्वीच नाकारला जातो.

## सुरक्षितता आणि अपवादात्मक स्थिती

- **चालू खाती.** बदल ज्या खात्यांना हलवतो त्यांपैकी एखाद्या खात्याने Claude Code चालू असताना तो बदल नाकारला जातो, कारण session त्या फाइल्समध्ये लिहीत राहते. बाहेर पडणे आणि सामील होणे फक्त नाव दिलेल्या खात्याशी संबंधित असते. `--dry-run` तरीही योजना छापतो.
- **बॅकअप.** प्रत्येक योजनेला खात्याच्या फोल्डरमध्ये स्वतःचा `.ccam-backup/<time>/` फोल्डर मिळतो. ccam ते कधीच हटवत नाही.
- **तुमच्या स्वतःच्या लिंक्स.** एखादा पाथ आधीच तुम्ही केलेली लिंक असेल, उदाहरणार्थ dotfiles repository कडे, आणि गट रिकामा असेल, तर ccam लिंकचे target गटात कॉपी करते आणि लिंक स्वतः बॅकअपमध्ये हलवते. target ला हात लावला जात नाही.
- **इम्पोर्ट केलेले फोल्डर.** `ccam import` ने आहे त्या जागी नोंदवलेले खाते (जसे `~/.claude`) त्याच जागी बदलते, त्यामुळे साध्या `claude` लाही बदल दिसतो; योजना तसे सांगते.
- **एकाच फोल्डरमध्ये दोन खाती.** इतिहास शेअर करणारी खाती त्याच फाइल्समध्ये लिहितात. ज्या फोल्डरमध्ये एक खाते आधीच चालू आहे तिथे दुसरे सुरू झाले, तर ccam इशारा देते आणि विचारते.
- **हाताने केलेले बदल.** खाते नसलेला सदस्य दुर्लक्षित केला जातो आणि `ccam doctor` तो नोंदवते. `local`, `on`, `off` आणि `none` ही गटांची नावे असू शकत नाहीत.

## संबंधित

- [डॅशबोर्ड](/features/dashboard)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [शेअर केलेले MCP सर्व्हर](/features/mcp)
- [आधीचा सेटअप कॉपी करा](/features/migrate)
- [सर्व कमांड्स](/docs/commands)
