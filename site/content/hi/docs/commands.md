---
title: कमांड संदर्भ
description: ccam की हर कमांड और विकल्प, ccam add और ccam <name> से लेकर ccam usage, ccam next और ccam doctor तक।
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

कोई कमांड असल में चलाए बिना यह देखने के लिए कि वह क्या करेगी, उसके आगे `--dry-run` लगाएँ, उदाहरण के लिए `ccam --dry-run work`।

<!-- commands -->

## Claude Code शुरू करना

`ccam <name|number> [arguments]` Claude Code को उस अकाउंट से शुरू करता है। ccam नाम के बाद वाले आर्ग्युमेंट नहीं पढ़ता; वह उन्हें उसी क्रम में `claude` को दे देता है, और standard input, exit code और signals ठीक वैसे ही काम करते हैं जैसे ख़ुद `claude` के साथ। नाम के बाद `--` लगाना वैकल्पिक है।

आर्ग्युमेंट के साथ लेकिन बिना अकाउंट नाम के `ccam`, जैसे `ccam -p "summarise"`, डिफ़ॉल्ट अकाउंट इस्तेमाल करता है। इसे `ccam use <name>` से सेट करें।

शॉर्टकट `claude1`, `claude2` आदि `~/.local/bin` में ccam के लिंक हैं। `ccam shims install --names` हर अकाउंट के लिए `claude-<name>` भी बनाता है।

## अकाउंट जोड़ना और हटाना

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

फ़ोल्डर के बिना `ccam import`, `~/.claude` और `~/.claude-*` खोजता है और हर एक के बारे में पूछता है। Import किए गए फ़ोल्डर जहाँ हैं वहीं रजिस्टर होते हैं; ccam कभी लॉगिन को move या copy नहीं करता।

## साझा करना

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

हर आइटम `on` (`default` ग्रुप), `off` या किसी ग्रुप का नाम लेता है। [साझा करना](/docs/sharing) देखें।

## MCP सर्वर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## उपयोग और स्विच करना

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## रखरखाव

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## फ़ाइलें

| पाथ | सामग्री |
| --- | --- |
| `~/.config/ccam/config.toml` | अकाउंट्स और वे क्या साझा करते हैं |
| `~/.local/share/ccam/profiles/<name>/` | हर अकाउंट का `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | साझा सेटिंग्स, MCP सर्वर, एक्सटेंशन और इतिहास |
| `~/.local/state/ccam/` | उपयोग का cache और आख़िरी इस्तेमाल का समय |
| `~/.local/bin/claude1`, `claude2`… | ccam के लिंक |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` और `XDG_STATE_HOME` variables इन फ़ोल्डरों की जगह बदल देते हैं।
