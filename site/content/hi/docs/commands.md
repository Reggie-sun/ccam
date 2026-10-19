---
title: कमांड संदर्भ
description: ccam की हर कमांड और विकल्प, ccam add और ccam <name> से लेकर ccam usage, ccam next और ccam doctor तक।
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

कोई कमांड असल में चलाए बिना यह देखने के लिए कि वह क्या करेगी, उसके आगे `--dry-run` लगाएँ, उदाहरण के लिए `ccam --dry-run work`।

<!-- commands -->

## Claude Code शुरू करना

`ccam <name|number> [arguments]` Claude Code को उस अकाउंट से शुरू करता है। ccam नाम के बाद वाले आर्ग्युमेंट नहीं पढ़ता; वह उन्हें उसी क्रम में `claude` को दे देता है, और standard input, exit code और signals ठीक वैसे ही काम करते हैं जैसे ख़ुद `claude` के साथ। नाम के बाद `--` लगाना वैकल्पिक है।

आर्ग्युमेंट के साथ लेकिन बिना अकाउंट नाम के `ccam`, जैसे `ccam -p "summarise"`, डिफ़ॉल्ट अकाउंट इस्तेमाल करता है। इसे `ccam use <name>` से सेट करें।

`ccwork` जैसा शॉर्टकट अपने अकाउंट को इसी तरह शुरू करता है: `ccwork --resume` वही करता है जो `ccam work --resume` करता है। [शॉर्टकट](/docs/commands#shortcuts) देखें।

## अकाउंट जोड़ना और हटाना

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

फ़ोल्डर के बिना `ccam import`, `~/.claude` और `~/.claude-*` खोजता है और हर एक के बारे में पूछता है। Import किए गए फ़ोल्डर जहाँ हैं वहीं रजिस्टर होते हैं; ccam कभी लॉगिन को move या copy नहीं करता।

## शॉर्टकट

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

शॉर्टकट एक कमांड है जिसका नाम आप खुद रखते हैं और जो Claude Code को एक अकाउंट से शुरू करती है। `ccam work` की तरह यह हर आर्ग्युमेंट बिना बदलाव के Claude Code को देती है, इसलिए `ccwork --resume` और `ccwork -p "hi"` वैसे ही काम करते हैं जैसी आप उम्मीद करते हैं। यह `~/.local/bin` में ccam का एक लिंक है।

शॉर्टकट वैकल्पिक हैं। `ccam add` का विज़र्ड पूछता है कि शॉर्टकट बनाना है या नहीं, `ccam add <name> --shortcut <command>` बिना पूछे एक शॉर्टकट बनाता है, और `ccam list` हर अकाउंट का शॉर्टकट दिखाता है, या शॉर्टकट न होने पर `-`।

नामों में a-z, 0-9, `.`, `-` और `_` अक्षर इस्तेमाल होते हैं। नाम कोई ccam कमांड, किसी दूसरे अकाउंट का शॉर्टकट, या आपके `PATH` पर मौजूद किसी दूसरे प्रोग्राम का नाम नहीं हो सकता।

`ccam rename` शॉर्टकट का नाम वही रखता है, और `ccam remove` अकाउंट के साथ उसका शॉर्टकट भी हटा देता है। `ccam shims install` सभी अकाउंट्स के लिंक फिर से बनाता है, जैसे ccam को दूसरी जगह ले जाने के बाद।

## साझा करना

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

हर आइटम `on` (`default` ग्रुप), `off` या किसी ग्रुप का नाम लेता है। [साझा करना](/docs/sharing) देखें।

## मौजूदा सेटअप कॉपी करना

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` किसी Claude Code फ़ोल्डर से `settings.json`, `CLAUDE.md`, एक्सटेंशन और user MCP सर्वर एक ग्रुप में कॉपी करता है, ताकि उन आइटम को साझा करने वाले हर अकाउंट को वे मिल जाएँ। फ़ोल्डर को केवल पढ़ा जाता है। `--settings`, `--memory`, `--extensions`, `--mcp` और `--history` विकल्प केवल उन्हीं आइटम को कॉपी करते हैं जिनका आप नाम देते हैं; इनके बिना बातचीत के इतिहास को छोड़कर सब कुछ कॉपी होता है। `-g` या `--group` ग्रुप चुनता है, और `-y` या `--yes` बिना पूछे कॉपी करता है और कोई फ़ाइल अलग होने पर ग्रुप की कॉपी रखता है। [मौजूदा सेटअप कॉपी करना](/docs/sharing#copy-an-existing-setup) देखें।

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
ccam shims install   # recreate the shortcut links
```

## फ़ाइलें

| पाथ | सामग्री |
| --- | --- |
| `~/.config/ccam/config.toml` | अकाउंट्स और वे क्या साझा करते हैं |
| `~/.local/share/ccam/profiles/<name>/` | हर अकाउंट का `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | साझा सेटिंग्स, MCP सर्वर, एक्सटेंशन और इतिहास |
| `~/.local/state/ccam/` | उपयोग का cache और आख़िरी इस्तेमाल का समय |
| `~/.local/bin/<shortcut>` | `ccwork` जैसी शॉर्टकट कमांड, ccam के लिंक |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` और `XDG_STATE_HOME` variables इन फ़ोल्डरों की जगह बदल देते हैं।
