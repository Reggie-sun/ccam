---
title: कमांड संदर्भ
description: ccam add आणि ccam <name> पासून ccam usage, ccam next आणि ccam doctor पर्यंत ccam ची प्रत्येक कमांड आणि पर्याय.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

कोणतीही कमांड प्रत्यक्ष न चालवता ती काय करेल हे पाहण्यासाठी तिच्या आधी `--dry-run` लावा, उदाहरणार्थ `ccam --dry-run work`.

<!-- commands -->

## Claude Code सुरू करणे

`ccam <name|number> [arguments]` त्या खात्याने Claude Code सुरू करते. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही; ते त्याच क्रमाने `claude` ला पाठवते, आणि standard input, exit code आणि signals अगदी `claude` प्रमाणेच वागतात. नावानंतर `--` लावणे ऐच्छिक आहे.

खात्याच्या नावाशिवाय पण आर्ग्युमेंट्ससह `ccam`, जसे `ccam -p "summarise"`, डीफॉल्ट खाते वापरते. ते `ccam use <name>` ने सेट करा.

`claude1`, `claude2` वगैरे शॉर्टकट हे `~/.local/bin` मधील ccam च्या लिंक्स आहेत. `ccam shims install --names` प्रत्येक खात्यासाठी `claude-<name>` देखील तयार करते.

## खाती जोडणे आणि काढणे

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

फोल्डर न देता `ccam import` चालवल्यास ते `~/.claude` आणि `~/.claude-*` शोधते आणि प्रत्येकाबद्दल विचारते. इम्पोर्ट केलेले फोल्डर आहेत त्याच जागी नोंदवले जातात; ccam कधीही लॉगिन हलवत किंवा कॉपी करत नाही.

## शेअरिंग

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

प्रत्येक आयटमसाठी `on` (`default` गट), `off`, किंवा गटाचे नाव देता येते. [शेअरिंग](/docs/sharing) पाहा.

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
ccam shims install
```

## फाइल्स

| पाथ | मजकूर |
| --- | --- |
| `~/.config/ccam/config.toml` | खाती आणि ती काय शेअर करतात |
| `~/.local/share/ccam/profiles/<name>/` | प्रत्येक खात्याचा `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | शेअर केलेल्या सेटिंग्ज, MCP सर्व्हर, एक्स्टेंशन्स आणि इतिहास |
| `~/.local/state/ccam/` | वापराचा cache आणि शेवटच्या वापराच्या वेळा |
| `~/.local/bin/claude1`, `claude2`… | ccam च्या लिंक्स |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` आणि `XDG_STATE_HOME` हे व्हेरिएबल्स या फोल्डर्सची जागा बदलतात.
