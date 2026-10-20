---
title: कमांड संदर्भ
metaTitle: "ccam कमांड संदर्भ"
description: ccam add आणि ccam <name> पासून ccam usage, ccam next आणि ccam doctor पर्यंत ccam ची प्रत्येक कमांड आणि पर्याय.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

एखादी कमांड प्रत्यक्ष न चालवता ती काय करेल हे पाहण्यासाठी तिच्या आधी `--dry-run` लावा, उदाहरणार्थ `ccam --dry-run work`. हे Claude Code सुरू करताना (`ccam <name>`, `ccam run`, `ccam next`) आणि `share`, `migrate`, `shortcut`, `update` व `doctor --fix` सोबत चालते. काहीतरी बदलणाऱ्या इतर कमांड्स चालण्याऐवजी हा पर्याय नाकारतात.

<!-- commands -->

## Claude Code सुरू करणे

`ccam <name|number> [arguments]` त्या खात्याने Claude Code सुरू करते. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही; ते त्याच क्रमाने `claude` ला पाठवते, आणि standard input, exit code आणि signals अगदी `claude` प्रमाणेच वागतात. नावानंतर `--` लावणे ऐच्छिक आहे.

खात्याच्या नावाशिवाय पण आर्ग्युमेंट्ससह `ccam`, जसे `ccam -p "summarise"`, सध्याच्या फोल्डरसाठी सेट केलेले खाते वापरते, नाहीतर डीफॉल्ट खाते. ती `ccam folder <name> [folder]` आणि `ccam use <name>` ने सेट करा.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

खात्याच्या नावाआधी `--fallback` हे `-p`/`--print` runs साठी आहे. खाते आपल्या वापर मर्यादेला पोहोचले की ccam तीच कमांड सर्वात जास्त जागा असलेल्या लॉगिन केलेल्या खात्याने पुन्हा चालवते आणि फक्त पूर्ण झालेल्या run चा output छापते. [फोल्डरनुसार खाती](/features/folder) आणि [मर्यादेवर पुन्हा चालवणे](/features/next#run-again-at-the-limit) पाहा.

`ccwork` सारखा शॉर्टकट त्याचे खाते याच प्रकारे सुरू करतो: `ccwork --resume` तेच करते जे `ccam work --resume` करते. [शॉर्टकट](/docs/commands#shortcuts) पाहा.

## खाती जोडणे आणि काढणे

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

फोल्डर न देता `ccam import` चालवल्यास ते `~/.claude` आणि `~/.claude-*` शोधते आणि प्रत्येकाबद्दल विचारते. इम्पोर्ट केलेले फोल्डर आहेत त्याच जागी नोंदवले जातात; ccam कधीही लॉगिन हलवत किंवा कॉपी करत नाही.

`ccam add --api` असे खाते जोडते जे Claude लॉगिनऐवजी API key वापरते; ccam प्रत्येक launch ला key वाचते आणि ती कधीही साठवत नाही. [API key खाती](/features/api) पहा.

एखाद्या खात्याने Claude Code चालू असताना `ccam rename`, `ccam remove` आणि `ccam share` त्या खात्यात बदल करण्यास नकार देतात. आधी session बंद करा.

`ccam list --json` प्रत्येक खाते त्याचा slot क्रमांक, शॉर्टकट, ईमेल आणि फोल्डरसह छापते. त्यातील `share` object प्रत्येक घटक (`mcp`, `settings`, `memory`, `extensions`, `history`) एखाद्या गटाच्या नावाशी किंवा `local` शी जोडतो.

## शॉर्टकट

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

शॉर्टकट ही तुम्ही स्वतः नाव दिलेली कमांड आहे, जी Claude Code एका खात्याने सुरू करते. `ccam work` प्रमाणेच ती प्रत्येक आर्ग्युमेंट न बदलता Claude Code ला देते, त्यामुळे `ccwork --resume` आणि `ccwork -p "hi"` अपेक्षेप्रमाणे चालतात. ती `~/.local/bin` मधील ccam ची लिंक आहे.

शॉर्टकट ऐच्छिक आहेत. `ccam add` चा विझार्ड शॉर्टकट तयार करायचा का ते विचारतो, `ccam add <name> --shortcut <command>` न विचारता एक तयार करते, आणि `ccam list` प्रत्येक खात्याचा शॉर्टकट दाखवते, किंवा शॉर्टकट नसेल तर `-`.

नावांमध्ये a-z, 0-9, `.`, `-` आणि `_` ही अक्षरे चालतात. नाव ccam ची कमांड, दुसऱ्या खात्याचा शॉर्टकट, किंवा तुमच्या `PATH` वरील दुसऱ्या प्रोग्रामचे नाव असू शकत नाही.

`ccam rename` शॉर्टकटचे नाव तेच ठेवते, आणि `ccam remove` खात्यासोबत त्याचा शॉर्टकटही काढून टाकते. `ccam shims install` सर्व खात्यांच्या लिंक्स पुन्हा तयार करते, उदाहरणार्थ ccam दुसरीकडे हलवल्यानंतर.

## शेअरिंग

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

प्रत्येक आयटमसाठी `on` (`default` गट), `off`, किंवा गटाचे नाव देता येते. [शेअरिंग](/docs/sharing) पाहा.

## आधीचा सेटअप कॉपी करणे

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` एका Claude Code फोल्डरमधून `settings.json`, `CLAUDE.md`, एक्स्टेंशन्स आणि user MCP सर्व्हर एका गटात कॉपी करते, म्हणजे ते आयटम शेअर करणाऱ्या प्रत्येक खात्याला ते मिळतात. फोल्डर फक्त वाचला जातो. `--settings`, `--memory`, `--extensions`, `--mcp` आणि `--history` हे पर्याय फक्त तुम्ही नाव दिलेले आयटम कॉपी करतात. त्यांच्याशिवाय, टर्मिनलवर एक चेकलिस्ट तुम्हाला निवडू देते, ज्यात संभाषण इतिहास सोडून सर्व काही निवडलेले असते; इतरत्र संभाषण इतिहास सोडून सर्व काही कॉपी होते. `-g` किंवा `--group` गट निवडतो, आणि `-y` किंवा `--yes` न विचारता कॉपी करतो आणि एखादी फाइल वेगळी असल्यास गटाची प्रत ठेवतो. ते सर्व्हरची सेव्ह केलेली साइन-इन्स प्रत्येक खात्याच्या स्वतःच्या लॉगिन फाइलमध्येही कॉपी करते. [आधीचा सेटअप कॉपी करणे](/docs/sharing#copy-an-existing-setup) पाहा.

## MCP सर्व्हर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` फक्त त्याच खात्याकडे असलेले सर्व्हर त्याच्या गटात हलवते. ज्या सर्व्हरचे नाव गट आधीच वेगळ्या व्याख्येसाठी वापरत असेल तो वगळला जातो; गटाची प्रत बदलण्यासाठी `--force` जोडा, ccam आधी तिचा बॅकअप घेते.

## वापर आणि खाते बदलणे

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` आणि `ccam pick` एकाच पद्धतीने निवडतात: आधी `auto_max_usage` (`config.toml`, डीफॉल्ट 100) खालची खाती, `auto_strategy` (`headroom` किंवा `expiring`) नुसार, `--model` ची साप्ताहिक window मोजून, आणि 100% वरचे खाते कधीच नाही. काहीच निवडता येत नसेल तर ते code 75 ने बाहेर पडतात, किंवा `--wait` सह पहिले खाते reset होण्याची वाट पाहतात. [पुढच्या खात्यावर जाणे](/features/next) आणि [scripts साठी खाते निवडणे](/features/pick) पहा.

## देखभाल

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` ccam बदलण्यापूर्वी रिलीजचा SHA-256 checksum तपासते, आणि `--version v0.2.0` विशिष्ट रिलीज इन्स्टॉल करते. `.deb` पॅकेजमधून इन्स्टॉल केलेले ccam त्याऐवजी apt ने अपडेट करा. तपशीलासाठी [अपडेट](/docs/install#update) पाहा.

## फाइल्स

| पाथ | मजकूर |
| --- | --- |
| `~/.config/ccam/config.toml` | खाती आणि ती काय शेअर करतात |
| `~/.local/share/ccam/profiles/<name>/` | प्रत्येक खात्याचा `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | शेअर केलेल्या सेटिंग्ज, MCP सर्व्हर, एक्स्टेंशन्स आणि इतिहास |
| `~/.local/state/ccam/` | वापराचा cache आणि शेवटच्या वापराच्या वेळा |
| `~/.local/bin/<shortcut>` | `ccwork` सारख्या शॉर्टकट कमांड्स, ccam च्या लिंक्स |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` आणि `XDG_STATE_HOME` हे व्हेरिएबल्स या फोल्डर्सची जागा बदलतात. Windows वर `~` म्हणजे `%USERPROFILE%`, आणि शॉर्टकट `<shortcut>.exe` असतात.
