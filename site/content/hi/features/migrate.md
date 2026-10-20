---
title: मौजूदा सेटअप कॉपी करना
description: ccam migrate से किसी मौजूदा ~/.claude की सेटिंग्स, CLAUDE.md, एक्सटेंशन, MCP सर्वर और MCP साइन-इन एक ccam ग्रुप में कॉपी करें।
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 7cee52e75749
---

अगर आपने ccam से पहले Claude Code इस्तेमाल किया है, तो आपका सेटअप `~/.claude` और `~/.claude.json` में है। `ccam migrate` उसे एक साझा ग्रुप में कॉपी करता है, ताकि ग्रुप के ज़रिए उन आइटम को साझा करने वाला हर अकाउंट आपके पहले से मौजूद सेटअप के साथ शुरू हो। Source फ़ोल्डर को केवल पढ़ा जाता है, और सादा `claude` पहले की तरह काम करता रहता है। इसे अपने पहले `ccam add` के बाद एक बार चलाएँ, और फिर जब भी कोई अकाउंट MCP सर्वर साझा करना शुरू करे।

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
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (या `x`) cursor के नीचे वाले आइटम को टॉगल करता है, arrow keys या `j`/`k` ऊपर-नीचे ले जाते हैं, Enter योजना दिखाता है, और Esc, `q` या Ctrl+C `Cancelled.` प्रिंट करते हैं। जिन आइटम के लिए source में कुछ नहीं है, उन पर `not found` लिखा होता है। अगर आप सब कुछ अचयनित कर दें, तो ccam `Nothing selected.` प्रिंट करके रुक जाता है।

फिर ccam योजना प्रिंट करता है और कॉपी करने से पहले पूछता है:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

ग्रुप केवल अपने अकाउंट्स तक पहुँचता है, और केवल उन आइटम के लिए जिन्हें वह साझा करता है, इसलिए सारांश ग्रुप के अकाउंट्स की सूची देता है। जब ग्रुप कॉपी किए गए आइटम में से कोई साझा नहीं करता, या कुछ अकाउंट उसमें नहीं हैं, तो ccam इसे बदलने वाली `ccam group set <group> --add <items>` और `ccam group join <group> <accounts>` कमांड प्रिंट करता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `[folder]` | वह Claude Code फ़ोल्डर जिससे कॉपी करना है | `~/.claude` |
| `--settings` | `settings.json` कॉपी करें | बंद |
| `--memory` | ग्लोबल `CLAUDE.md` कॉपी करें | बंद |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` कॉपी करें | बंद |
| `--mcp` | user MCP सर्वर कॉपी करें (`.claude.json` में `mcpServers`) | बंद |
| `--history` | बातचीत का इतिहास कॉपी करें (`projects/`, `history.jsonl`) | बंद |
| `-g`, `--group string` | वह साझा ग्रुप जिसमें कॉपी करना है | `default` |
| `-y`, `--yes` | बिना पूछे कॉपी करें; टकराव होने पर ग्रुप की कॉपी रखी जाती है | बंद |
| `--dry-run` (ग्लोबल) | योजना प्रिंट करें और किसी सवाल या कॉपी से पहले रुक जाएँ | बंद |

कोई आइटम flag न होने पर, टर्मिनल पर चेकलिस्ट तय करती है। टर्मिनल के बाहर, या `--yes` के साथ, ccam बातचीत के इतिहास को छोड़कर सब कुछ कॉपी करता है। इतिहास डिफ़ॉल्ट रूप से छोड़ा जाता है, क्योंकि `projects/` बड़ा हो सकता है।

`--group on` का मतलब `default` है। `--group off` और `--group local` को `--group takes a group name, not "off"` के साथ अस्वीकार किया जाता है। `ccam migrate`, `--dry-run` को मानता है, और dry run में भी चेकलिस्ट दिखती है, ताकि आप किसी चयन का पूर्वावलोकन कर सकें।

## यह कैसे काम करता है

`internal/cli/migrate.go` आइटम चुनता है, फिर `migrate.Build` source और ग्रुप को पढ़कर एक योजना लौटाता है। आपकी पुष्टि के बाद `Plan.Apply` तक कुछ भी नहीं लिखा जाता।

1. **Source खोजना।** `checkSource` ऐसे पाथ को अस्वीकार करता है जो मौजूद नहीं है, फ़ोल्डर नहीं है, या जिसमें Claude Code की कोई ज्ञात entry नहीं है। `migrate.ClaudeJSON` फ़ोल्डर से जुड़ी `.claude.json` खोजता है: उसके अंदर, या ऐसे `~/.claude` के लिए जिसके अंदर यह नहीं है, उसके बगल में `~/.claude.json`। `CLAUDE_CONFIG_DIR` सेट न होने पर Claude Code इसे वहीं रखता है।
2. **आइटम को पाथ से जोड़ना।** हर पाथ `share.Layout` से आता है, वही अनुमत सूची जिसे ग्रुप इस्तेमाल करते हैं। MCP सर्वर `.claude.json` से `mcp.LocalServers` के ज़रिए पढ़े जाते हैं, जो पूरी फ़ाइल parse करने के बजाय `profile.TopLevelValue` scanner इस्तेमाल करता है।
3. **कॉपी, कभी move नहीं।** फ़ाइलें `fsutil.AtomicWrite` या `fsutil.CopyFile` से `~/.local/share/ccam/groups/<group>/` में जाती हैं; MCP सर्वर ग्रुप की `mcp.json` में जाते हैं।

हर आइटम एक source पाथ को एक ग्रुप पाथ से जोड़ता है:

<!-- motion:migrate-items -->

इतिहास की दोनों पंक्तियाँ डैश वाली हैं, क्योंकि वे केवल `--history` के साथ कॉपी होती हैं। `.credentials.json` और `.claude.json`, `share.Never` में दर्ज हैं और कभी ग्रुप तक नहीं पहुँचतीं: अगर इनमें से कोई कभी लेआउट में दिखे, तो `Build` एक आंतरिक त्रुटि लौटाता है, और फ़ोल्डर की खोज इन्हें नाम से छोड़ देती है।

`buildEntry` हर पाथ को एक स्थिति देता है, और योजना की हर पंक्ति यही बताती है:

<!-- motion:migrate-states -->

- **Missing**: source में नहीं है, या जहाँ फ़ोल्डर अपेक्षित है वहाँ फ़ाइल है। हर आइटम के लिए एक बार `not found, skipped` के रूप में दिखाया जाता है।
- **Linked**: source पाथ पहले से ग्रुप की कॉपी पर पहुँचता है, उदाहरण के लिए जब source कोई ccam profile फ़ोल्डर हो।
- **Same**: ग्रुप में पहले से एक जैसी सामग्री है।
- **New**: ग्रुप में यह नहीं है, इसलिए यह पूरा कॉपी होता है।
- **Merge**: दोनों तरफ़ एक फ़ोल्डर। `copyMissing` केवल वे फ़ाइलें जोड़ता है जो ग्रुप में नहीं हैं; अलग फ़ाइलें वैसी ही रहती हैं जैसी ग्रुप में हैं, और योजना `group copies kept` कहती है। `plugins/` के लिए, registries `installed_plugins.json` और `known_marketplaces.json` merge होती हैं, और plugins के absolute पाथ ग्रुप के फ़ोल्डर की ओर इशारा करने के लिए फिर से लिखे जाते हैं।
- **Conflict**: दोनों तरफ़ अलग सामग्री वाली एक फ़ाइल। आप हर फ़ाइल के लिए चुनते हैं: `keep group`, `use mine`, साथ में `.json` फ़ाइलों के लिए `merge JSON` (एक deep merge जिसमें टकराने वाले मानों पर ग्रुप जीतता है) या `history.jsonl` के लिए `merge by line`। जिस MCP सर्वर की definition अलग है, वह भी टकराव है, `keep group` या `use mine` के साथ। सूची का पहला विकल्प Enter दबाने पर डिफ़ॉल्ट है; `--yes` हमेशा ग्रुप की कॉपी रखता है।

चेकलिस्ट ख़ुद `tui.MigratePicker` है, जो उसी `checklist` component पर बनी है जिसे `ccam add` इस्तेमाल करता है; `migrate.Found` lstat calls और MCP सर्वर सूची से गायब आइटम चिह्नित करता है।

### MCP साइन-इन

OAuth से साइन इन करने वाले MCP सर्वर अपने tokens लॉगिन फ़ाइल `.credentials.json` में `mcpOAuth` के नीचे रखते हैं, सर्वर सूची में नहीं। Claude Code हर entry की key `<name>|` और सर्वर के `type`, `url` और `headers` के `sha256` hash के पहले 16 अक्षरों से बनाता है। Key फ़ोल्डर पर निर्भर नहीं करती, इसलिए कॉपी की गई entry दूसरे अकाउंट में तब तक काम करती है जब तक सर्वर की definition एक जैसी हो। इसलिए `internal/migrate/logins.go` में `buildLogins` केवल उन सर्वरों की entries कॉपी करता है जो ग्रुप में आख़िर में source की definition के साथ रहते हैं: नए, पहले से एक जैसे, या ऐसा टकराव जिसमें आपने `use mine` चुना। Plugin सर्वरों के साइन-इन (`plugin:` से शुरू होने वाली keys) तब साथ जाते हैं जब आप एक्सटेंशन कॉपी करते हैं और source में `plugins/` होता है।

<!-- motion:migrate-signins -->

Entries पाने वाले हर अकाउंट की अपनी `.credentials.json` में जाती हैं, जो `fsutil.AtomicWrite` से मोड `0600` पर लिखी जाती है। वे कभी ग्रुप में नहीं जातीं। कोई अकाउंट उन्हें केवल तभी पाता है जब वह लक्ष्य ग्रुप के ज़रिए `mcp` (या plugin सर्वरों के लिए `extensions`) साझा करता हो, और केवल वे keys जोड़ी जाती हैं जो उसके पास नहीं हैं: जो entries उसके पास पहले से हैं, शायद नई, वे रखी जाती हैं। योजना की `sign-ins` पंक्ति सर्वर, पाने वाले अकाउंट और छोड़े गए अकाउंट दिखाती है। ccam किसी अकाउंट को छोड़ देता है, और कारण के साथ उसका नाम बताता है, जब वह:

- `ccam import` से अपनी जगह पर रजिस्टर हो, क्योंकि सादा `claude` उसकी लॉगिन फ़ाइल इस्तेमाल कर रहा हो सकता है;
- चल रहा हो, क्योंकि Claude Code उसकी लॉगिन फ़ाइल फिर से लिखता है (लिखने से ठीक पहले फिर जाँचा जाता है, और `started meanwhile` के रूप में बताया जाता है);
- अभी लॉगिन न हुआ हो;
- अपना लॉगिन macOS कीचेन में रखता हो, जिसमें ccam नहीं लिखता।

छोड़े गए अकाउंट में `/mcp` से सर्वरों में साइन इन करें। जो अकाउंट बाद में MCP सर्वर साझा करना शुरू करता है, उसे साइन-इन तब मिलते हैं जब आप `ccam migrate` फिर से चलाते हैं; सारांश यह ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` के साथ बताता है।

## सुरक्षा और विशेष स्थितियाँ

- Source को केवल पढ़ा जाता है। इसकी कोई फ़ाइल move, link या बदली नहीं जाती, और source का लॉगिन कभी नहीं लिखा जाता, macOS पर भी, जहाँ इसे कीचेन से पढ़ा जा सकता है।
- `use mine`, merge, या किसी मौजूदा `mcp.json` में नए MCP सर्वरों से ग्रुप फ़ाइल बदले जाने से पहले, ग्रुप की कॉपी `<group>/.ccam-backup/<time>/` में सेव होती है। फ़ोल्डर का नाम हर रन के लिए अलग होता है, इसलिए एक ही सेकंड में दो रन एक-दूसरे के बैकअप नहीं मिटाते। सारांश पाथ प्रिंट करता है।
- `--yes` के साथ, हर टकराव में ग्रुप की कॉपी रखी जाती है, और ccam बताता है कि इसकी जगह कैसे चुनें: वही कमांड `--yes` के बिना चलाएँ।
- अगर कुछ नहीं बदलेगा, तो ccam `Nothing to copy.` प्रिंट करता है; अगर हर टकराव में ग्रुप की कॉपी रखी जाती है, तो `Nothing to copy: the group keeps its own copies.`
- अगर कॉपी बीच में विफल हो, तो ccam `The source folder was not changed. Group files replaced so far are saved in` और उसके बाद बैकअप फ़ोल्डर के साथ रुकता है।
- Source न मिलने पर यह `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` के साथ ख़त्म होता है।

## संबंधित

- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [जाँच और मरम्मत](/features/doctor)
- [कई Claude Code अकाउंट](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
