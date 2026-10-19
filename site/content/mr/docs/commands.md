---
title: कमांड संदर्भ
description: ccam add आणि ccam <name> पासून ccam usage, ccam next आणि ccam doctor पर्यंत ccam ची प्रत्येक कमांड आणि पर्याय.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

कोणतीही कमांड प्रत्यक्ष न चालवता ती काय करेल हे पाहण्यासाठी तिच्या आधी `--dry-run` लावा, उदाहरणार्थ `ccam --dry-run work`.

<!-- commands -->

## Claude Code सुरू करणे

`ccam <name|number> [arguments]` त्या खात्याने Claude Code सुरू करते. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही; ते त्याच क्रमाने `claude` ला पाठवते, आणि standard input, exit code आणि signals अगदी `claude` प्रमाणेच वागतात. नावानंतर `--` लावणे ऐच्छिक आहे.

खात्याच्या नावाशिवाय पण आर्ग्युमेंट्ससह `ccam`, जसे `ccam -p "summarise"`, डीफॉल्ट खाते वापरते. ते `ccam use <name>` ने सेट करा.

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
```

फोल्डर न देता `ccam import` चालवल्यास ते `~/.claude` आणि `~/.claude-*` शोधते आणि प्रत्येकाबद्दल विचारते. इम्पोर्ट केलेले फोल्डर आहेत त्याच जागी नोंदवले जातात; ccam कधीही लॉगिन हलवत किंवा कॉपी करत नाही.

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

`ccam migrate` एका Claude Code फोल्डरमधून `settings.json`, `CLAUDE.md`, एक्स्टेंशन्स आणि user MCP सर्व्हर एका गटात कॉपी करते, म्हणजे ते आयटम शेअर करणाऱ्या प्रत्येक खात्याला ते मिळतात. फोल्डर फक्त वाचला जातो. `--settings`, `--memory`, `--extensions`, `--mcp` आणि `--history` हे पर्याय फक्त तुम्ही नाव दिलेले आयटम कॉपी करतात; त्यांच्याशिवाय संभाषण इतिहास सोडून सर्व काही कॉपी होते. `-g` किंवा `--group` गट निवडतो, आणि `-y` किंवा `--yes` न विचारता कॉपी करतो आणि एखादी फाइल वेगळी असल्यास गटाची प्रत ठेवतो. [आधीचा सेटअप कॉपी करणे](/docs/sharing#copy-an-existing-setup) पाहा.

## MCP सर्व्हर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## वापर आणि खाते बदलणे

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## देखभाल

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
```

## फाइल्स

| पाथ | मजकूर |
| --- | --- |
| `~/.config/ccam/config.toml` | खाती आणि ती काय शेअर करतात |
| `~/.local/share/ccam/profiles/<name>/` | प्रत्येक खात्याचा `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | शेअर केलेल्या सेटिंग्ज, MCP सर्व्हर, एक्स्टेंशन्स आणि इतिहास |
| `~/.local/state/ccam/` | वापराचा cache आणि शेवटच्या वापराच्या वेळा |
| `~/.local/bin/<shortcut>` | `ccwork` सारख्या शॉर्टकट कमांड्स, ccam च्या लिंक्स |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` आणि `XDG_STATE_HOME` हे व्हेरिएबल्स या फोल्डर्सची जागा बदलतात.
