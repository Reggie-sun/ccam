---
title: शेअर केलेले MCP सर्व्हर
description: ccam mcp add, add-json, list, remove आणि import ने प्रत्येक गटासाठी MCP सर्व्हरची एक यादी ठेवा, जी --mcp-config ने दिली जाते किंवा प्रत्येक खात्यात sync होते.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
---

Claude Code user MCP सर्व्हर `.claude.json` मध्ये ठेवते, ज्या फाइलमध्ये खात्याचे लॉगिनही असते, त्यामुळे ती फाइल खात्यांमध्ये शेअर करता येत नाही. ccam शेअर केलेले सर्व्हर `mcp.json` या गट फाइलमध्ये ठेवते आणि सुरू करताना ती गटातील प्रत्येक खात्याला देते. तुमच्या सर्व खात्यांसाठी एकदाच सर्व्हर जोडण्यासाठी `ccam mcp` वापरा, आणि एका session मध्ये जोडलेला सर्व्हर इतर खात्यांकडे हलवण्यासाठी `ccam mcp import` वापरा.

## वापर

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

बदल प्रत्येक खाते पुढच्या वेळी सुरू होईल तेव्हापासून लागू होतो; चालू sessions ज्या सर्व्हरसह सुरू झाले ते ठेवतात:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

गटांच्या खाली, `ccam mcp list` इतर खात्यांपर्यंत न पोहोचणारा प्रत्येक सर्व्हर सांगते:

- एखाद्या खात्याने स्वतःसाठी जोडलेल्या सर्व्हरसाठी `only in account work: pg`;
- एकाच फोल्डरसाठी जोडलेल्या सर्व्हरसाठी (साधे `claude mcp add`, खाली पाहा) `only in account side, for /home/me/code/app only: notes`;
- खात्याच्या स्वतःच्या सर्व्हरला गटाच्या सर्व्हरचे नाव असेल तेव्हा `account work has its own …, used in its sessions in place of the group's`;
- या फोल्डरचे `.mcp.json`, ज्या खात्यांच्या मंजुरीची ते वाट पाहत आहे त्यांच्यासह (`waiting for approval in`);
- `in no account, only in ~/.claude.json (plain claude's)` असलेले सर्व्हर, ते कॉपी करणाऱ्या [`ccam migrate --mcp`](/features/migrate) कमांडसह.

गटाकडे आधीच असलेले नाव जोडल्यास त्याची व्याख्या बदलली जाते आणि `added` ऐवजी `updated` छापले जाते. सर्व्हर एखाद्या खात्यापर्यंत तेव्हाच पोहोचतात जेव्हा ते खाते त्या गटात असते आणि तो गट MCP शेअर करतो (`ccam group set <group> --add mcp`); गट तसे करत नसेल, तर `add` तसे सांगते.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` आणि `remove` साठी लक्ष्य गट | ज्या एकमेव गटामार्फत एखादे खाते MCP शेअर करते तो, नाहीतर `default` |
| `-- <command> [args…]` (`add`) | ही कमांड stdio सर्व्हर म्हणून चालवा | |
| `--url <url>` (`add`) | रिमोट सर्व्हरचा URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` सर्व्हरचा transport | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, अनेकदा देता येतो | काहीही नाही |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio सर्व्हरसाठी environment variable, अनेकदा देता येतो | काहीही नाही |
| `--force` (`import`) | त्याच नावाचा वेगळा गट सर्व्हर बदला; जुन्या `mcp.json` चा बॅकअप घेतला जातो | बंद |
| `--keep-local` (`import`) | खात्याच्या `.claude.json` मधील नोंदी ठेवा | बंद |
| `-h`, `--help` | प्रत्येक कमांडची मदत | |
| `--dry-run` (global) | फक्त `mcp list` ला समर्थित; `add`, `add-json`, `remove` आणि `import` तो नाकारतात आणि काहीही बदलत नाहीत | बंद |

`add` ला `--` नंतर कमांड किंवा `--url` यांपैकी एक आवश्यक आहे; नाहीतर ते `a command to run (after --) or --url is required` देऊन अयशस्वी होते. `--header` `--url` सर्व्हरना लागू होतो आणि `--env` कमांड सर्व्हरना. `add-json` `claude mcp add-json` सारखाच JSON घेते, आणि व्याख्या `command` किंवा `url` असलेला object असणे आवश्यक आहे. `import` नेहमी तो गट वापरते ज्यामार्फत खाते MCP शेअर करते, आणि `list` प्रत्येक गट दाखवते, त्यामुळे `-g` ने त्यांच्यात फरक पडत नाही.

## हे कसे काम करते

संबंधित फाइल्स:

- `~/.local/share/ccam/groups/<group>/mcp.json`: शेअर केलेले सर्व्हर, `{"mcpServers": {…}}`, `--mcp-config` वाचतो तो format. ती `fsutil.AtomicWrite` ने आणि `0600` mode सह लिहिली जाते, कारण सर्व्हरच्या नोंदींमध्ये अनेकदा API keys असतात, आणि फक्त `<group>/.mcp.json.lock` या lock खालीच, त्यामुळे एकाच वेळी चाललेल्या दोन `ccam mcp add` कमांड्स आपापले सर्व्हर ठेवतात.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: खात्याचे सर्व्हर असलेली Claude Code ची स्वतःची फाइल. ccam `profile.TopLevelValue` scanner ने फक्त तिचे `mcpServers` मूल्य वाचते आणि सुरू करण्याच्या मार्गावर संपूर्ण फाइल कधीच parse करत नाही.
- `~/.local/state/ccam/mcp/<name>/`: `seen` फाइल, `sync-stamp` (ccam ने लिहिलेल्या प्रत्येक प्रतीच्या hash सह, जेणेकरून कोणीही न बदललेली प्रत आणि तुम्ही संपादित केलेली प्रत यांतील फरक त्याला कळतो), `sync-checked` (sync strategy ने शेवटी काय पाहिले), `kept` (MCP शेअरिंग सोडताना खात्याने ठेवलेले सर्व्हर) आणि `.claude.json` चे बॅकअप.

दोन strategies आहेत, ज्या `~/.config/ccam/config.toml` मधील `mcp_strategy` ने सेट होतात. डीफॉल्ट flag strategy आहे:

<!-- motion:mcp-flag -->

1. ccam खात्याचा MCP गट शोधते (`launcher.applyMCP`) आणि एका `stat` ने गटाची `mcp.json` आहे का ते तपासते, ती parse न करता (`mcp.FlagArgs`).
2. ccam `--mcp-config=<path>` `claude` च्या लगेच नंतर, तुमच्या स्वतःच्या आर्ग्युमेंट्सच्या आधी ठेवते, त्यामुळे तुमचे स्वतःचे `--mcp-config` किंवा `--strict-mcp-config` त्यावर लागू होत राहतात (`Plan.InsertArgs`).
3. ccam मागच्या वेळी सुरू केल्यापासून खात्याची `.claude.json` बदलली का ते तपासते, आणि तेव्हाच तिचे सर्व्हर parse करते (`mcp.NewLocal`).
4. ccam स्वतःच्या जागी Claude Code ठेवते (`syscall.Exec`).

पाथ दोन नाही तर `--mcp-config=<path>` हा एक आर्ग्युमेंट म्हणून दिला जातो. Claude Code चा `--mcp-config` अनेक मूल्ये घेतो, त्यामुळे `claude --mcp-config a.json "prompt"` prompt ला दुसरी config फाइल म्हणून वाचेल आणि `MCP config file not found` देऊन अयशस्वी होईल. एका आर्ग्युमेंटचे रूप हे टाळते, त्यामुळे त्यानंतर येणारा prompt किंवा subcommand चालतो. `ccam work mcp list`, `get` आणि `remove` यांना flag मिळत नाही, कारण Claude Code तिथे तो नाकारते; त्याऐवजी ccam standard error वर गटाच्या सर्व्हरची नावे सांगते ([तपशील](/features/launch#shared-mcp-servers)). Claude Code `--mcp-config` चे सर्व्हर खात्याच्या स्वतःच्या सर्व्हरसोबत लोड करते; नाव जुळल्यास `.claude.json` मधील नोंद जिंकते. flag strategy `.claude.json` फक्त Remote Control साठी, एखादे खाते MCP शेअरिंग सोडते किंवा त्यात पुन्हा सामील होते तेव्हा ([खाली](#leaving-and-joining)), आणि `ccam mcp import` साठी लिहिते.

`claude remote-control` (`rc`) कोणताही `--mcp-config` घेत नाही, आणि त्याचे sessions सर्व्हर फक्त `.claude.json` मधून वाचतात. `ccam work rc` साठी, ccam flag वगळते आणि टर्मिनलवर विचारते की गटाचे सर्व्हर बॅकअपसह `.claude.json` मध्ये कॉपी करायचे का, जसे `sync` strategy करते. कॉपी तयार होईपर्यंत ते विचारते, नंतर फक्त गटाची `mcp.json` बदलल्यानंतर किंवा कॉपी केलेला सर्व्हर नाहीसा झाल्यानंतरच ([Remote Control](/features/launch#remote-control)). नाव जुळल्यास अशी कॉपी `--mcp-config` वर जिंकत असल्याने, `ccam mcp add`, `add-json`, `remove` आणि `import` गटातील इतर खात्यांकडील कॉपी अपडेट करतात:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated` ओळ फक्त एखादी प्रत बदलली तेव्हाच दिसते. ccam च्या जुन्या आवृत्तीने तयार केलेली प्रतही अशाच प्रकारे अद्ययावत केली जाते, आणि `ccam migrate` नंतरची प्रतही.

session मध्ये `claude mcp add` ने जोडलेला सर्व्हर फक्त त्या खात्याच्या `.claude.json` मध्ये जातो. साधे `claude mcp add` local scope वापरते: सर्व्हर फक्त त्या एका खात्यासाठी, त्या एका फोल्डरमध्ये लोड होतो (`-s user` तो प्रत्येक फोल्डरमध्ये त्या खात्यासाठी उपलब्ध करतो). हे घडते तेव्हा ccam ते पाहू शकत नाही, कारण Claude Code ने आधीच ccam process ची जागा घेतलेली असते, म्हणून ते पुढच्या वेळी सुरू करताना तपासते. `seen` फाइलमध्ये `.claude.json` ची modification time आणि आधी नोंदवलेली सर्व्हरची नावे असतात. वेळ वेगळी असल्यास, ccam नावांची गटाशी तुलना करते आणि प्रत्येक नवीन सर्व्हरसाठी एकदा सूचना छापते, input ची वाट न पाहता, दोन्ही strategies सह:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

गटाच्या सर्व्हरचे नाव असलेल्या खात्याच्या स्वतःच्या सर्व्हरचे नावही ही सूचना सांगते, कारण त्याच्या सत्रांमध्ये खात्याचा सर्व्हर जिंकतो.

`mcp_strategy = "sync"` सह, सर्व्हर `claude mcp list` आणि `/mcp` मध्ये साध्या user सर्व्हरप्रमाणे दिसतात, आणि ccam सुरू करण्यापूर्वी ते प्रत्येक खात्याच्या `.claude.json` मध्ये लिहिते:

<!-- motion:mcp-sync -->

1. ccam गटाच्या `mcp.json` च्या modification time ची `sync-stamp` शी, आणि खात्याच्या `.claude.json` ची त्याने शेवटी पाहिलेल्या स्थितीशी तुलना करते (`mcp.Sync`). दोन्हीपैकी काहीही बदलले नसेल, तर पुढे काहीही होत नाही. नाहीतर ते गटाचा प्रत्येक सर्व्हर अजून `.claude.json` मध्ये आहे का तेही तपासते, आणि नाहीशी झालेली प्रत परत ठेवते (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. ते Claude Code चे स्वतःचे lock घेते, म्हणजे फाइलशेजारी `mkdir` ने तयार केलेला `.claude.json.lock` फोल्डर. ते 2 सेकंदांपर्यंत थांबते आणि 10 सेकंदांपेक्षा जुने lock शिळे मानते, जसे Claude Code ची lock library करते.
3. ते सध्याची `.claude.json` state फोल्डरमधील `backups/` मध्ये सेव्ह करते आणि सर्वात नवीन 5 ठेवते.
4. ते `mcpServers` मध्ये प्रत्येक गट सर्व्हर सेट करते, आधी स्वतः लिहिलेले पण आता गटाकडे नसलेले सर्व्हर काढते, आणि काहीतरी बदलले असेल तरच फाइल तिच्या मूळ mode सह atomic पद्धतीने लिहिते. ते फक्त स्वतः लिहिलेल्या, आणि त्यानंतर कोणीही न बदललेल्या, प्रती बदलते किंवा काढते. खात्याकडे गटातील सर्व्हरच्याच नावाचा स्वतःचा सर्व्हर असेल, तर ccam खात्याचा सर्व्हर ठेवते आणि एकदा तसे सांगते: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. ते प्रत्येक प्रतीच्या hash सह नवीन stamp नोंदवते, आणि Claude Code `--mcp-config` शिवाय सुरू होते.

खात्याचे session चालू असताना, Claude Code कसेही सुरू केलेले असो, sync `ccam: account work is already running, skipping MCP sync` सह वगळला जातो, कारण चालू असताना Claude Code `.claude.json` पुन्हा लिहिते. sync अयशस्वी झाल्यास `ccam: MCP sync failed (launching anyway): …` छापले जाते आणि तरीही सुरू केले जाते.

### सोडणे आणि पुन्हा सामील होणे

MCP शेअर करणे थांबवणारे खाते (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) गटाचे सर्व्हर ठेवते: योजना ते बॅकअपसह त्याच्या `.claude.json` मध्ये कॉपी करते, आणि त्याच नावांचे त्याचे स्वतःचे सर्व्हर तसेच राहतात. `--fresh` सह त्याला एकही मिळत नाही, आणि योजना लोड होणे थांबणाऱ्या सर्व्हरची यादी देते. ते सर्व्हर असलेल्या गटात ते पुन्हा सामील होते, तेव्हा कोणीही न बदललेल्या ठेवलेल्या प्रती काढल्या जातात, त्यामुळे गटाचे नंतरचे बदल त्याच्यापर्यंत पोहोचतात.

`ccam mcp import` खात्याचे स्वतःचे सर्व्हर त्याच्या गटात हलवते:

<!-- motion:mcp-import -->

ccam खात्याचे सर्व्हर, त्याचे स्वतःचे आणि फक्त एकाच फोल्डरसाठी असलेले, वाचते आणि प्रत्येकाची गटातील त्याच नावाच्या सर्व्हरशी `mcp.SameDef` ने तुलना करते, जे keys चा क्रम आणि spacing दुर्लक्षित करते. फोल्डरचा सर्व्हर प्रत्येक फोल्डरसाठीचा गट सर्व्हर बनतो (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); दोन वेगवेगळ्या फोल्डर व्याख्या असलेले एकच नाव नाकारले जाते. Plugin सर्व्हर सुचवले जात नाहीत, कारण plugin ते स्वतः आणते. गटाकडे नसलेला किंवा त्याच व्याख्येसह असलेला सर्व्हर कॉपी होतो. वेगळा असलेला वगळला जातो, कारण इतर प्रत्येक खाते गटाचीच व्याख्या वापरते:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` सह गटाचा सर्व्हर बदलला जातो, आणि आधीची `mcp.json` आधी गट फोल्डरमधील `.ccam-backup/<time>/mcp.json` मध्ये सेव्ह होते. मग ccam कॉपी झालेल्या नोंदी खात्याच्या `.claude.json` मधून काढते, त्याच्या स्वतःच्या आणि फोल्डरच्या नोंदी एकाच write मध्ये, त्याच lock खाली आणि state फोल्डरमध्ये बॅकअप घेतल्यानंतर (`mcp.Remove`). ccam गटाची `mcp.json` वाचल्यापासून त्या नोंदी काढल्या जाईपर्यंत ती locked ठेवते, त्यामुळे दरम्यान कोणीही गटातून तो सर्व्हर काढत नाही. हे आवश्यक आहे कारण नाव जुळल्यास खात्याची नोंद गटाच्या नोंदीवर जिंकेल. `--keep-local` सह, त्या नोंदी खात्याच्या `.claude.json` मध्ये तशाच राहतात.

## सुरक्षितता आणि अपवादात्मक स्थिती

- **साइन-इन्स.** OAuth ने साइन इन करणारे सर्व्हर आपले tokens सर्व्हर यादीत नाही, तर प्रत्येक खात्याच्या `.credentials.json` (`mcpOAuth`) मध्ये ठेवतात, त्यामुळे सर्व्हर शेअर केल्याने त्याचे साइन-इन शेअर होत नाही. प्रत्येक खात्यात `/mcp` ने एकदा साइन इन करा. ccam कधीही साइन-इन्स कॉपी करत नाही, आणि [`ccam migrate`](/features/migrate) सुद्धा करत नाही: कॉपी केलेला refresh token ते असलेल्या इतर प्रत्येकाला लॉग आउट करतो. एखादे खाते ज्या सर्व्हरमध्ये साइन इन नाही पण दुसरे खाते आहे, त्या सर्व्हरची नावे [`ccam ext`](/features/ext) आणि `ccam doctor` सांगतात, आणि नवीन खात्याकडे कोणते नाहीत ते `ccam next --continue` सांगते. गटाला कधीच लॉगिन फाइल मिळत नाही.
- **चालू खाती.** `import` गटात कॉपी करते पण चालू खात्याच्या स्वतःच्या नोंदी तशाच ठेवते आणि तसे सांगते; session संपल्यानंतर ते पुन्हा चालवा. sync चालू खाती वगळतो. दोन्ही ccam शिवाय सुरू केलेली Claude Code सत्रेही ओळखतात.
- **Remote Control साठीच्या कॉपी.** चालू खात्याच्या कॉपीला हात लावला जात नाही: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` सह काहीही अपडेट करण्याची गरज नाही, कारण प्रत्येक वेळी सुरू करताना sync होते.
- **बिघडलेल्या फाइल्स.** `.claude.json` वैध JSON नसेल, तर sync तिला हात लावत नाही आणि त्रुटी नोंदवतो. 2 सेकंदांत lock घेता आले नाही, तर ccam `could not lock … (Claude Code is writing the file)` नोंदवते.
- **MCP शेअर न करणारी खाती.** `import` नकार देते आणि आधी ते खाते MCP शेअर करणाऱ्या गटात टाकायला सांगते (`ccam group join <group> side`).
- **गटांची नावे.** `-g` फक्त गटाचे नाव घेतो; ते `groups/` खाली फोल्डर बनते, त्यामुळे `../x` सारखी नावे नाकारली जातात.

## संबंधित

- [गटांसह शेअरिंग](/features/share)
- [ccam सोबत Remote Control वापरणे](/guides/remote-control)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [आधीचा सेटअप कॉपी करणे](/features/migrate)
- [कमांड संदर्भ](/docs/commands)
