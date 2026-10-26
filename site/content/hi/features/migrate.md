---
title: मौजूदा सेटअप कॉपी करना
description: ccam migrate से किसी मौजूदा ~/.claude की सेटिंग्स, CLAUDE.md, एक्सटेंशन और MCP सर्वर एक ccam ग्रुप में कॉपी करें।
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

अगर आपने ccam से पहले Claude Code इस्तेमाल किया है, तो आपका सेटअप `~/.claude` और `~/.claude.json` में है। `ccam migrate` उसे एक साझा ग्रुप में कॉपी करता है, ताकि ग्रुप के ज़रिए उन आइटम को साझा करने वाला हर अकाउंट आपके पहले से मौजूद सेटअप के साथ शुरू हो। Source फ़ोल्डर को केवल पढ़ा जाता है, और सादा `claude` पहले की तरह काम करता रहता है। इसे अपने पहले `ccam add` के बाद एक बार चलाएँ; जब `~/.claude.json` में MCP सर्वर हों, तो पहला `ccam add` अपने सुझाव में उनके नाम बताता है।

## उपयोग

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

टर्मिनल पर, आइटम flags और `--yes` के बिना, ccam पहले एक चेकलिस्ट दिखाता है:

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space (या `x`) cursor के नीचे वाले आइटम को टॉगल करता है, arrow keys या `j`/`k` ऊपर-नीचे ले जाते हैं, Enter योजना दिखाता है, और Esc, `q` या Ctrl+C `Cancelled.` प्रिंट करते हैं। जिन आइटम के लिए source में कुछ नहीं है, उन पर `not found` लिखा होता है। अगर आप सब कुछ अचयनित कर दें, तो ccam `Nothing selected.` प्रिंट करके रुक जाता है।

फिर ccam योजना प्रिंट करता है और कॉपी करने से पहले पूछता है:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

ग्रुप केवल अपने अकाउंट्स पर असर डालता है, और केवल उन आइटम पर जिन्हें वह साझा करता है, इसलिए सारांश ग्रुप के अकाउंट्स की सूची देता है। जब ग्रुप कॉपी किए गए आइटम में से कोई साझा नहीं करता, या कुछ अकाउंट उसमें नहीं हैं, तो ccam इसे बदलने वाली `ccam group set <group> --add <items>` और `ccam group join <group> <accounts>` कमांड प्रिंट करता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `[folder]` | वह Claude Code फ़ोल्डर जिससे कॉपी करना है | `~/.claude` |
| `--settings` | `settings.json` कॉपी करें | बंद |
| `--memory` | ग्लोबल `CLAUDE.md` कॉपी करें | बंद |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` कॉपी करें | बंद |
| `--mcp` | user MCP सर्वर कॉपी करें (`.claude.json` में `mcpServers`) | बंद |
| `--history` | बातचीत का इतिहास कॉपी करें (`projects/`, `history.jsonl`) | बंद |
| `-g`, `--group <group>` | वह साझा ग्रुप जिसमें कॉपी करना है | `default` |
| `-y`, `--yes` | बिना पूछे कॉपी करें; टकराव होने पर ग्रुप की कॉपी रखी जाती है | बंद |
| `--dry-run` (ग्लोबल) | योजना प्रिंट करें और किसी सवाल या कॉपी से पहले रुक जाएँ | बंद |

टर्मिनल पर, अगर आप कोई आइटम flag नहीं देते, तो आप चेकलिस्ट में आइटम चुनते हैं। टर्मिनल के बाहर, या `--yes` के साथ, ccam बातचीत के इतिहास को छोड़कर सब कुछ कॉपी करता है। इतिहास डिफ़ॉल्ट रूप से छोड़ा जाता है, क्योंकि `projects/` बड़ा हो सकता है।

`--group on` का मतलब `default` है। `--group off` और `--group local` को `--group takes a group name, not "off"` के साथ अस्वीकार किया जाता है। `ccam migrate`, `--dry-run` का समर्थन करता है, और dry run में भी चेकलिस्ट दिखती है, ताकि आप किसी चयन का पूर्वावलोकन कर सकें।

## यह कैसे काम करता है

ccam कुछ भी कॉपी करने से पहले एक पूरी योजना बनाता है। यह आइटम चुनता है (`internal/cli/migrate.go`), फिर source और ग्रुप को पढ़कर योजना बनाता है (`migrate.Build`)। आपकी पुष्टि के बाद `Plan.Apply` तक कुछ भी नहीं लिखा जाता।

1. **Source खोजना।** ccam ऐसे पाथ को अस्वीकार करता है जो मौजूद नहीं है, फ़ोल्डर नहीं है, या जिसमें Claude Code की कोई ज्ञात entry नहीं है (`checkSource`)। केवल MCP सर्वरों के लिए, ख़ाली या ग़ायब `~/.claude` ठीक है, बशर्ते `~/.claude.json` में सर्वर हों। ccam फ़ोल्डर से जुड़ी `.claude.json` खोजता है (`migrate.ClaudeJSON`): उसके अंदर, या ऐसे `~/.claude` के लिए जिसके अंदर यह नहीं है, उसके बगल में `~/.claude.json`। `CLAUDE_CONFIG_DIR` सेट न होने पर Claude Code इसे वहीं रखता है।
2. **आइटम को पाथ से जोड़ना।** हर पाथ `share.Layout` से आता है, वही अनुमत सूची जिसे ग्रुप इस्तेमाल करते हैं। MCP सर्वर `.claude.json` से `mcp.LocalServers` के ज़रिए पढ़े जाते हैं, जो पूरी फ़ाइल parse करने के बजाय `profile.TopLevelValue` scanner इस्तेमाल करता है।
3. **कॉपी, कभी move नहीं।** फ़ाइलें `~/.local/share/ccam/groups/<group>/` में `fsutil.AtomicWrite` या `fsutil.CopyFile` से जाती हैं, जो target के बगल में एक अस्थायी फ़ाइल लिखता है और उसका नाम बदलता है, ताकि बाधित कॉपी कभी अधूरी फ़ाइल न छोड़े। MCP सर्वर ग्रुप की `mcp.json` में उसके lock के तहत जाते हैं: इस बीच कोई दूसरा ccam process जो सर्वर जोड़ता है वह रखा जाता है, और [Remote Control](/guides/remote-control) के लिए बनी कॉपियाँ बाद में अपडेट की जाती हैं।

हर आइटम एक source पाथ को एक ग्रुप पाथ से जोड़ता है:

<!-- motion:migrate-items -->

इतिहास की दोनों पंक्तियाँ डैश वाली हैं, क्योंकि वे केवल `--history` के साथ कॉपी होती हैं। `.credentials.json` और `.claude.json`, `share.Never` में दर्ज हैं और कभी ग्रुप तक नहीं पहुँचतीं: अगर इनमें से कोई कभी लेआउट में दिखे, तो ccam एक आंतरिक त्रुटि के साथ रुक जाता है (`Build`), और फ़ोल्डर की खोज इन्हें नाम से छोड़ देती है।

ccam हर पाथ को एक स्थिति देता है, और योजना की हर पंक्ति यही बताती है (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: source में नहीं है, या जहाँ फ़ोल्डर अपेक्षित है वहाँ फ़ाइल है। हर आइटम के लिए एक बार `not found, skipped` के रूप में दिखाया जाता है।
- **Linked**: source पाथ पहले से ग्रुप की कॉपी पर पहुँचता है, उदाहरण के लिए जब source कोई ccam profile फ़ोल्डर हो।
- **Same**: ग्रुप में पहले से एक जैसी सामग्री है।
- **New**: ग्रुप में यह नहीं है, इसलिए यह पूरा कॉपी होता है।
- **Merge**: दोनों तरफ़ एक फ़ोल्डर। ccam केवल वे फ़ाइलें जोड़ता है जो ग्रुप में नहीं हैं (`copyMissing`); अलग फ़ाइलें वैसी ही रहती हैं जैसी ग्रुप में हैं, और योजना `group copies kept` कहती है। `plugins/` के लिए, registries `installed_plugins.json` और `known_marketplaces.json` merge होती हैं, और plugins के absolute पाथ ग्रुप के फ़ोल्डर की ओर इशारा करने के लिए फिर से लिखे जाते हैं।
- **Conflict**: दोनों तरफ़ अलग सामग्री वाली एक फ़ाइल। आप हर फ़ाइल के लिए चुनते हैं: `keep group`, `use mine`, साथ में `.json` फ़ाइलों के लिए `merge JSON` (एक deep merge जिसमें टकराने वाले मानों पर ग्रुप जीतता है) या `history.jsonl` के लिए `merge by line`। जिस MCP सर्वर की definition अलग है, वह भी टकराव है, `keep group` या `use mine` के साथ। सूची का पहला विकल्प Enter दबाने पर डिफ़ॉल्ट है; `--yes` हमेशा ग्रुप की कॉपी रखता है।

चेकलिस्ट ख़ुद `tui.MigratePicker` है, जो उसी `checklist` component पर बनी है जिसे `ccam add` इस्तेमाल करता है; ccam lstat calls और MCP सर्वर सूची से गायब आइटम चिह्नित करता है (`migrate.Found`)।

### MCP साइन-इन

OAuth से साइन इन करने वाले MCP सर्वर अपने tokens लॉगिन फ़ाइल `.credentials.json` में `mcpOAuth` के नीचे रखते हैं, सर्वर सूची में नहीं। ccam उन्हें कॉपी नहीं करता। कॉपी किए गए साइन-इन में मूल वाला ही refresh token होता है, और सर्वर पहले renewal पर उस token को बदल देता है: तब हर दूसरा धारक, source `~/.claude` भी, उस सर्वर से लॉग आउट हो जाता है। इसलिए हर उस कॉपी किए गए सर्वर के लिए जिसमें source साइन इन है, योजना की `sign-ins` पंक्ति ग्रुप के उन अकाउंट्स की सूची देती है जिन्हें `/mcp` से ख़ुद साइन इन करना होगा, और सारांश इसे दोहराता है। बिना token वाली entry साइन-इन नहीं है और सूची में नहीं आती, न ही ऐसा अकाउंट जिसके पास पहले से अपना साइन-इन है, या ऐसा सर्वर जिसकी definition ग्रुप अपनी रखता है। कोई लॉगिन फ़ाइल नहीं लिखी जाती।

## सुरक्षा और विशेष स्थितियाँ

- Source को केवल पढ़ा जाता है। इसकी कोई फ़ाइल move, link या बदली नहीं जाती, और source का लॉगिन कभी नहीं लिखा जाता, macOS पर भी, जहाँ इसे कीचेन से पढ़ा जा सकता है।
- `use mine`, merge, या किसी मौजूदा `mcp.json` में नए MCP सर्वरों से ग्रुप फ़ाइल बदले जाने से पहले, ग्रुप की कॉपी `<group>/.ccam-backup/<time>/` में सेव होती है। फ़ोल्डर का नाम हर रन के लिए अलग होता है, इसलिए एक ही सेकंड में दो रन एक-दूसरे के बैकअप नहीं मिटाते। सारांश पाथ प्रिंट करता है।
- `--yes` के साथ, हर टकराव में ग्रुप की कॉपी रखी जाती है, और ccam बताता है कि इसकी जगह कैसे चुनें: वही कमांड `--yes` के बिना चलाएँ।
- अगर कुछ नहीं बदलेगा, तो ccam `Nothing to copy.` प्रिंट करता है; अगर हर टकराव में ग्रुप की कॉपी रखी जाती है, तो `Nothing to copy: the group keeps its own copies.`
- अगर कॉपी बीच में विफल हो, तो ccam `The source folder was not changed. Group files replaced so far are saved in` और उसके बाद बैकअप फ़ोल्डर के साथ रुकता है। इसे दोबारा चलाने पर वे फ़ाइलें कॉपी होती हैं जो पूरी तरह कॉपी नहीं हुई थीं।
- जब कई ग्रुप MCP साझा करते हैं, तो `ccam doctor`, `ccam ext` और `ccam mcp list` जो कमांड सुझाते हैं वे `ccam migrate --mcp --group <a|b>` होती हैं।
- Source न मिलने पर यह `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` के साथ ख़त्म होता है।

## संबंधित

- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [जाँच और मरम्मत](/features/doctor)
- [एक कंप्यूटर पर कई Claude Code अकाउंट इस्तेमाल करें](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR क्या है](/guides/claude-config-dir)
