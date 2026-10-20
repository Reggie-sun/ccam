---
title: आधीचा सेटअप कॉपी करणे
description: ccam migrate ने आधीपासून असलेल्या ~/.claude मधील सेटिंग्ज, CLAUDE.md, एक्स्टेंशन्स, MCP सर्व्हर आणि MCP साइन-इन्स ccam गटात कॉपी करा.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

तुम्ही ccam आधी Claude Code वापरले असेल, तर तुमचा सेटअप `~/.claude` आणि `~/.claude.json` मध्ये आहे. `ccam migrate` तो एका शेअर गटात कॉपी करते, त्यामुळे गटामार्फत ते आयटम शेअर करणारे प्रत्येक खाते तुमच्याकडे आधीपासून असलेल्या सेटअपने सुरू होते. स्रोत फोल्डर फक्त वाचला जातो, आणि साधे `claude` आधीसारखेच चालत राहते. तुमच्या पहिल्या `ccam add` नंतर ते एकदा चालवा, आणि एखादे खाते MCP सर्व्हर शेअर करू लागेल तेव्हा पुन्हा चालवा.

## वापर

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

टर्मिनलवर, आयटम flags शिवाय आणि `--yes` शिवाय, ccam आधी एक चेकलिस्ट दाखवते:

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

Space (किंवा `x`) cursor खालील आयटम टॉगल करते, arrow keys किंवा `j`/`k` हलवतात, Enter योजना दाखवते, आणि Esc, `q` किंवा Ctrl+C `Cancelled.` छापतात. ज्या आयटम्ससाठी स्रोतात काहीच नाही त्यांच्यावर `not found` चिन्ह असते. तुम्ही सर्व काही अनचेक केले, तर ccam `Nothing selected.` छापून थांबते.

मग ccam योजना छापते आणि कॉपी करण्यापूर्वी विचारते:

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

गट फक्त त्याच्या स्वतःच्या खात्यांपर्यंत पोहोचतो, आणि तेही फक्त तो शेअर करत असलेल्या आयटम्ससाठी, म्हणून सारांश गटाच्या खात्यांची यादी देतो. कॉपी केलेल्यांपैकी एखादा आयटम गट शेअर करत नसेल, किंवा काही खाती त्यात नसतील, तर ccam ते बदलणाऱ्या `ccam group set <group> --add <items>` आणि `ccam group join <group> <accounts>` कमांड्स छापते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `[folder]` | ज्यातून कॉपी करायचे तो Claude Code फोल्डर | `~/.claude` |
| `--settings` | `settings.json` कॉपी करा | बंद |
| `--memory` | global `CLAUDE.md` कॉपी करा | बंद |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` कॉपी करा | बंद |
| `--mcp` | user MCP सर्व्हर कॉपी करा (`.claude.json` मधील `mcpServers`) | बंद |
| `--history` | संभाषण इतिहास कॉपी करा (`projects/`, `history.jsonl`) | बंद |
| `-g`, `--group string` | ज्यात कॉपी करायचे तो शेअर गट | `default` |
| `-y`, `--yes` | न विचारता कॉपी करा; conflict असल्यास गटाची प्रत ठेवली जाते | बंद |
| `--dry-run` (global) | योजना छापा आणि कोणताही प्रश्न किंवा कॉपीपूर्वी थांबा | बंद |

कोणताही आयटम flag नसल्यास, टर्मिनलवर चेकलिस्ट ठरवते. टर्मिनलबाहेर, किंवा `--yes` सह, ccam संभाषण इतिहास सोडून सर्व काही कॉपी करते. इतिहास डीफॉल्टनुसार वगळला जातो कारण `projects/` मोठा असू शकतो.

`--group on` म्हणजे `default`. `--group off` आणि `--group local` `--group takes a group name, not "off"` देऊन नाकारले जातात. `ccam migrate` `--dry-run` पाळते, आणि dry run मध्येही चेकलिस्ट दिसते, त्यामुळे तुम्ही निवडीचे पूर्वावलोकन करू शकता.

## हे कसे काम करते

`internal/cli/migrate.go` आयटम निवडते, मग `migrate.Build` स्रोत आणि गट वाचते आणि योजना परत देते. तुम्ही खात्री दिल्यानंतर `Plan.Apply` होईपर्यंत काहीही लिहिले जात नाही.

1. **स्रोत शोधणे.** नसलेला, फोल्डर नसलेला किंवा Claude Code ची एकही ओळखीची नोंद नसलेला पाथ `checkSource` नाकारते. `migrate.ClaudeJSON` फोल्डरची `.claude.json` शोधते: त्याच्या आत, किंवा आत नसलेल्या `~/.claude` साठी त्याच्या शेजारी असलेली `~/.claude.json`. `CLAUDE_CONFIG_DIR` सेट नसताना Claude Code ती तिथेच ठेवते.
2. **आयटम्स पाथशी जोडणे.** प्रत्येक पाथ `share.Layout` मधून येतो, गट वापरतात तीच अनुमत यादी. MCP सर्व्हर `.claude.json` मधून `mcp.LocalServers` ने वाचले जातात, जे संपूर्ण फाइल parse करण्याऐवजी `profile.TopLevelValue` scanner वापरते.
3. **कॉपी, कधीच हलवणे नाही.** फाइल्स `fsutil.AtomicWrite` किंवा `fsutil.CopyFile` ने `~/.local/share/ccam/groups/<group>/` मध्ये जातात; MCP सर्व्हर गटाच्या `mcp.json` मध्ये जातात.

प्रत्येक आयटम एक स्रोत पाथ एका गट पाथशी जोडतो:

<!-- motion:migrate-items -->

इतिहासाच्या दोन ओळी तुटक आहेत कारण त्या फक्त `--history` सह कॉपी होतात. `.credentials.json` आणि `.claude.json` `share.Never` मध्ये नोंदवलेल्या आहेत आणि कधीच गटापर्यंत पोहोचत नाहीत: त्यांपैकी एखादी layout मध्ये दिसलीच तर `Build` internal error परत देते, आणि फोल्डर तपासताना त्या नावाने वगळल्या जातात.

`buildEntry` प्रत्येक पाथला एक स्थिती देते, आणि योजनेची प्रत्येक ओळ तेच वर्णन करते:

<!-- motion:migrate-states -->

- **Missing (नसलेले)**: स्रोतात नाही, किंवा जिथे फोल्डर अपेक्षित आहे तिथे फाइल. प्रत्येक आयटमसाठी एकदा `not found, skipped` म्हणून दाखवले जाते.
- **Linked (लिंक केलेले)**: स्रोत पाथ आधीच गटाच्या प्रतीकडे जातो, उदाहरणार्थ स्रोत ccam profile फोल्डर असतो तेव्हा.
- **Same (सारखे)**: गटाकडे आधीच तोच मजकूर आहे.
- **New (नवीन)**: गटाकडे ते नाही, म्हणून ते पूर्ण कॉपी होते.
- **Merge (एकत्रीकरण)**: दोन्ही बाजूंना फोल्डर. `copyMissing` फक्त गटाकडे नसलेल्या फाइल्स जोडते; वेगळ्या असलेल्या फाइल्स गटाकडे आहेत तशाच राहतात आणि योजना `group copies kept` सांगते. `plugins/` साठी, `installed_plugins.json` आणि `known_marketplaces.json` या registries merge होतात, आणि absolute plugin पाथ गटाच्या फोल्डरकडे निर्देश करतील असे पुन्हा लिहिले जातात.
- **Conflict (संघर्ष)**: दोन्ही बाजूंना वेगळा मजकूर असलेली फाइल. तुम्ही प्रत्येक फाइलसाठी निवडता: `keep group`, `use mine`, आणि `.json` फाइल्ससाठी `merge JSON` (deep merge, ज्यात जुळणाऱ्या मूल्यांवर गट जिंकतो) किंवा `history.jsonl` साठी `merge by line`. ज्या MCP सर्व्हरची व्याख्या वेगळी आहे तोही conflict आहे, `keep group` किंवा `use mine` सह. Enter दाबल्यावर यादीतील पहिला पर्याय डीफॉल्ट असतो; `--yes` नेहमी गटाची प्रत ठेवतो.

चेकलिस्ट स्वतः `tui.MigratePicker` आहे, जी `ccam add` वापरत असलेल्या `checklist` component वर बांधली आहे; `migrate.Found` lstat calls आणि MCP सर्व्हर यादीने नसलेल्या आयटम्सवर चिन्ह लावते.

### MCP साइन-इन्स

OAuth ने साइन इन करणारे MCP सर्व्हर आपले tokens सर्व्हर यादीत नाही, तर लॉगिन फाइलमध्ये, `.credentials.json` मध्ये `mcpOAuth` खाली ठेवतात. Claude Code प्रत्येक नोंदीची key `<name>|` आणि सर्व्हरच्या `type`, `url` आणि `headers` च्या `sha256` hash मधील पहिली 16 अक्षरे अशी ठेवते. key फोल्डरवर अवलंबून नसते, त्यामुळे सर्व्हरची व्याख्या अगदी तीच असेपर्यंत कॉपी केलेली नोंद दुसऱ्या खात्यात चालते. म्हणून `internal/migrate/logins.go` मधील `buildLogins` फक्त अशाच सर्व्हरच्या नोंदी कॉपी करते ज्या गटाकडे शेवटी स्रोताच्या व्याख्येसह असतात: नवीन, आधीच सारखे, किंवा जिथे तुम्ही `use mine` निवडले असा conflict. plugin सर्व्हरची साइन-इन्स (`plugin:` ने सुरू होणाऱ्या keys) तुम्ही एक्स्टेंशन्स कॉपी करता आणि स्रोतात `plugins/` असतो तेव्हा सोबत जातात.

<!-- motion:migrate-signins -->

नोंदी मिळणाऱ्या प्रत्येक खात्याच्या स्वतःच्या `.credentials.json` मध्ये जातात, `fsutil.AtomicWrite` ने `0600` mode सह लिहिलेल्या. त्या कधीच गटात जात नाहीत. खात्याला त्या फक्त तेव्हाच मिळतात जेव्हा ते लक्ष्य गटामार्फत `mcp` (किंवा plugin सर्व्हरसाठी `extensions`) शेअर करते, आणि फक्त त्याच्याकडे नसलेल्या keys जोडल्या जातात: त्याच्याकडे आधीच असलेल्या, कदाचित नवीन, नोंदी ठेवल्या जातात. योजनेतील `sign-ins` ओळ सर्व्हर, मिळणारी खाती आणि वगळलेली खाती दाखवते. खाते खालीलपैकी असेल तेव्हा ccam ते वगळते आणि कारणासह त्याचे नाव सांगते:

- `ccam import` ने जागेवरच नोंदवलेले, कारण साधे `claude` त्याची लॉगिन फाइल वापरत असू शकते;
- चालू, कारण Claude Code त्याची लॉगिन फाइल पुन्हा लिहिते (लिहिण्याच्या अगदी आधी पुन्हा तपासले जाते, `started meanwhile` म्हणून नोंदवले जाते);
- अजून लॉग इन न केलेले;
- आपले लॉगिन macOS Keychain मध्ये ठेवणारे, ज्यात ccam लिहीत नाही.

वगळलेल्या खात्यात `/mcp` ने सर्व्हरमध्ये साइन इन करा. जे खाते नंतर MCP सर्व्हर शेअर करू लागते, त्याला तुम्ही `ccam migrate` पुन्हा चालवल्यावर साइन-इन्स मिळतात; सारांश ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` ने तसे सांगतो.

## सुरक्षितता आणि अपवादात्मक स्थिती

- स्रोत फक्त वाचला जातो. त्यातील कोणतीही फाइल हलवली, लिंक केली किंवा बदलली जात नाही, आणि स्रोताचे लॉगिन कधीच लिहिले जात नाही, macOS वरही, जिथे ते Keychain मधून वाचले जाऊ शकते.
- `use mine`, merge, किंवा आधीपासून असलेल्या `mcp.json` मधील नवीन MCP सर्व्हर यांमुळे गटाची फाइल बदलण्यापूर्वी, गटाची प्रत `<group>/.ccam-backup/<time>/` मध्ये सेव्ह होते. फोल्डरचे नाव प्रत्येक run साठी वेगळे असते, त्यामुळे एकाच सेकंदातील दोन runs एकमेकांचे बॅकअप overwrite करत नाहीत. सारांश पाथ छापतो.
- `--yes` सह प्रत्येक conflict गटाची प्रत ठेवतो, आणि ccam त्याऐवजी कसे निवडायचे ते सांगते: तीच कमांड `--yes` शिवाय चालवा.
- काहीही बदलणार नसेल, तर ccam `Nothing to copy.` छापते; प्रत्येक conflict गटाची प्रत ठेवत असेल, तर `Nothing to copy: the group keeps its own copies.`
- कॉपी अर्ध्यावर अयशस्वी झाली, तर ccam `The source folder was not changed. Group files replaced so far are saved in` आणि त्यानंतर बॅकअप फोल्डर छापून थांबते.
- स्रोत नसल्यास `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` ने संपते.

## संबंधित

- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [तपासणी आणि दुरुस्ती](/features/doctor)
- [अनेक Claude Code खाती](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
