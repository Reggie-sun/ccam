---
title: शेअर केलेले MCP सर्व्हर
description: ccam mcp add, add-json, list, remove आणि import ने प्रत्येक गटासाठी MCP सर्व्हरची एक यादी ठेवा, जी --mcp-config ने दिली जाते किंवा प्रत्येक खात्यात sync होते.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
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
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
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
```

गटाकडे आधीच असलेले नाव जोडल्यास त्याची व्याख्या बदलली जाते आणि `added` ऐवजी `updated` छापले जाते. सर्व्हर एखाद्या खात्यापर्यंत तेव्हाच पोहोचतात जेव्हा ते खाते त्या गटासोबत MCP शेअर करते (`ccam share <name> --mcp on`); अजून कोणतेही खाते MCP शेअर करत नसेल, तर `add` तसे सांगते.

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
| `--dry-run` (global) | फक्त `mcp list` कडून पाळला जातो; `add`, `add-json`, `remove` आणि `import` तो नाकारतात आणि काहीही बदलत नाहीत | बंद |

`add` ला `--` नंतर कमांड किंवा `--url` यांपैकी एक आवश्यक आहे; नाहीतर ते `a command to run (after --) or --url is required` देऊन अयशस्वी होते. `--header` `--url` सर्व्हरना लागू होतो आणि `--env` कमांड सर्व्हरना. `add-json` `claude mcp add-json` सारखाच JSON घेते, आणि व्याख्या `command` किंवा `url` असलेला object असणे आवश्यक आहे. `import` नेहमी तो गट वापरते ज्यामार्फत खाते MCP शेअर करते, आणि `list` प्रत्येक गट दाखवते, त्यामुळे `-g` ने त्यांच्यात फरक पडत नाही.

## हे कसे काम करते

संबंधित फाइल्स:

- `~/.local/share/ccam/groups/<group>/mcp.json`: शेअर केलेले सर्व्हर, `{"mcpServers": {…}}`, `--mcp-config` वाचतो तो format. ती `fsutil.AtomicWrite` ने आणि `0600` mode सह लिहिली जाते, कारण सर्व्हरच्या नोंदींमध्ये अनेकदा API keys असतात.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: खात्याचे सर्व्हर असलेली Claude Code ची स्वतःची फाइल. ccam `profile.TopLevelValue` scanner ने फक्त तिचे `mcpServers` मूल्य वाचते आणि सुरू करण्याच्या मार्गावर संपूर्ण फाइल कधीच parse करत नाही.
- `~/.local/state/ccam/mcp/<name>/`: `seen` फाइल, `sync-stamp` आणि `.claude.json` चे बॅकअप.

दोन strategies आहेत, ज्या `~/.config/ccam/config.toml` मधील `mcp_strategy` ने सेट होतात. डीफॉल्ट flag strategy आहे:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` खात्याचा MCP गट शोधते. `mcp.FlagArgs` एका `stat` ने गटाची `mcp.json` आहे का ते तपासते, ती parse न करता.
2. `Plan.InsertArgs` `--mcp-config=<path>` `claude` च्या लगेच नंतर, तुमच्या स्वतःच्या आर्ग्युमेंट्सच्या आधी ठेवते, त्यामुळे तुमचे स्वतःचे `--mcp-config` किंवा `--strict-mcp-config` त्यावर लागू होत राहतात.
3. `mcp.NewLocal` मागच्या वेळी सुरू केल्यापासून खात्याची `.claude.json` बदलली का ते तपासते, आणि तेव्हाच तिचे सर्व्हर parse करते.
4. `syscall.Exec` ccam च्या जागी Claude Code ठेवते.

पाथ दोन नाही तर `--mcp-config=<path>` हा एक आर्ग्युमेंट म्हणून दिला जातो. Claude Code चा `--mcp-config` अनेक मूल्ये घेतो, त्यामुळे `claude --mcp-config a.json "prompt"` prompt ला दुसरी config फाइल म्हणून वाचेल आणि `MCP config file not found` देऊन अयशस्वी होईल. एका आर्ग्युमेंटचे रूप तिथेच थांबते, आणि त्यानंतर येणारा prompt किंवा `mcp list` सारखी subcommand चालते. Claude Code `--mcp-config` चे सर्व्हर खात्याच्या स्वतःच्या सर्व्हरसोबत लोड करते; नाव जुळल्यास `.claude.json` मधील नोंद जिंकते. flag strategy फक्त Remote Control साठी `.claude.json` लिहिते.

`claude remote-control` (`rc`) कोणताही `--mcp-config` घेत नाही, आणि त्याचे sessions सर्व्हर फक्त `.claude.json` मधून वाचतात. `ccam work rc` साठी, ccam flag वगळते आणि टर्मिनलवर विचारते की गटाचे सर्व्हर बॅकअपसह `.claude.json` मध्ये कॉपी करायचे का, जसे `sync` strategy करते. कॉपी तयार होईपर्यंत ते विचारते, नंतर फक्त गटाची `mcp.json` बदलल्यानंतरच ([Remote Control](/features/launch#remote-control)). नाव जुळल्यास अशी कॉपी `--mcp-config` वर जिंकत असल्याने, `ccam mcp add`, `add-json`, `remove` आणि `import` गटातील इतर खात्यांकडील कॉपी अपडेट करतात:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

session मध्ये `claude mcp add` ने जोडलेला सर्व्हर फक्त त्या खात्याच्या `.claude.json` मध्ये जातो. त्या वेळी ccam तो पाहू शकत नाही, कारण Claude Code ने त्याची जागा घेतलेली असते, म्हणून ते पुढच्या वेळी सुरू करताना तपासते. `seen` फाइलमध्ये `.claude.json` ची modification time आणि आधी नोंदवलेली सर्व्हरची नावे असतात. वेळ वेगळी असल्यास, ccam नावांची गटाशी तुलना करते आणि प्रत्येक नवीन सर्व्हरसाठी एकदा सूचना छापते, input ची वाट न पाहता:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` सह, सर्व्हर `claude mcp list` आणि `/mcp` मध्ये साध्या user सर्व्हरप्रमाणे दिसतात, आणि ccam सुरू करण्यापूर्वी ते प्रत्येक खात्याच्या `.claude.json` मध्ये लिहिते:

<!-- motion:mcp-sync -->

1. `mcp.Sync` गटाच्या `mcp.json` च्या modification time ची `sync-stamp` शी तुलना करते. ती बदलली नसेल, तर खात्याची `.claude.json` उघडली जात नाही.
2. ते Claude Code चे स्वतःचे lock घेते, म्हणजे फाइलशेजारी `mkdir` ने तयार केलेला `.claude.json.lock` फोल्डर. ते 2 सेकंदांपर्यंत थांबते आणि 10 सेकंदांपेक्षा जुने lock शिळे मानते, जसे Claude Code ची lock library करते.
3. ते सध्याची `.claude.json` state फोल्डरमधील `backups/` मध्ये सेव्ह करते आणि सर्वात नवीन 5 ठेवते.
4. ते `mcpServers` मध्ये प्रत्येक गट सर्व्हर सेट करते, आधी स्वतः लिहिलेले पण आता गटाकडे नसलेले सर्व्हर काढते, खात्याचे स्वतःचे सर्व्हर ठेवते, आणि फाइल तिच्या मूळ mode सह atomic पद्धतीने लिहिते.
5. ते नवीन stamp आणि ते सांभाळत असलेली नावे नोंदवते, आणि Claude Code `--mcp-config` शिवाय सुरू होते.

खात्याचे session चालू असताना sync `ccam: account work is already running, skipping MCP sync` सह वगळला जातो, कारण चालू असताना Claude Code `.claude.json` पुन्हा लिहिते. sync अयशस्वी झाल्यास `ccam: MCP sync failed (launching anyway): …` छापले जाते आणि तरीही सुरू केले जाते. session मध्ये जोडलेल्या सर्व्हरबद्दलची सूचना फक्त flag strategy मध्ये असते.

`ccam mcp import` खात्याचे स्वतःचे सर्व्हर त्याच्या गटात हलवते:

<!-- motion:mcp-import -->

ccam `mcp.LocalServers` ने खात्याचे सर्व्हर वाचते आणि प्रत्येकाची गटातील त्याच नावाच्या सर्व्हरशी `mcp.SameDef` ने तुलना करते, जे keys चा क्रम आणि spacing दुर्लक्षित करते. गटाकडे नसलेला किंवा त्याच व्याख्येसह असलेला सर्व्हर कॉपी होतो. वेगळा असलेला वगळला जातो, कारण इतर प्रत्येक खाते गटाचीच व्याख्या वापरते:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` सह गटाचा सर्व्हर बदलला जातो, आणि आधीची `mcp.json` आधी गट फोल्डरमधील `.ccam-backup/<time>/mcp.json` मध्ये सेव्ह होते. मग `mcp.RemoveLocal` कॉपी झालेल्या नोंदी खात्याच्या `.claude.json` मधून काढते, त्याच lock खाली आणि state फोल्डरमध्ये बॅकअप घेतल्यानंतर. हे आवश्यक आहे कारण नाव जुळल्यास खात्याची नोंद गटाच्या नोंदीवर जिंकेल. `--keep-local` त्या तशाच ठेवते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- **साइन-इन्स.** OAuth ने साइन इन करणारे सर्व्हर आपले tokens सर्व्हर यादीत नाही, तर प्रत्येक खात्याच्या `.credentials.json` (`mcpOAuth`) मध्ये ठेवतात, त्यामुळे सर्व्हर शेअर केल्याने त्याचे साइन-इन शेअर होत नाही. प्रत्येक खात्यात `/mcp` ने एकदा साइन इन करा, किंवा [`ccam migrate`](/features/migrate) ने `~/.claude` मधून साइन-इन्स कॉपी करा. गटाला कधीच लॉगिन फाइल मिळत नाही.
- **चालू खाती.** `import` गटात कॉपी करते पण चालू खात्याच्या स्वतःच्या नोंदी तशाच ठेवते आणि तसे सांगते; session संपल्यानंतर ते पुन्हा चालवा. sync चालू खाती वगळतो.
- **Remote Control साठीच्या कॉपी.** चालू खात्याच्या कॉपीला हात लावला जात नाही: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` सह काहीही अपडेट करण्याची गरज नाही, कारण प्रत्येक वेळी सुरू करताना sync होते.
- **बिघडलेल्या फाइल्स.** `.claude.json` वैध JSON नसेल, तर sync तिला हात लावत नाही आणि त्रुटी नोंदवतो. 2 सेकंदांत lock घेता आले नाही, तर ccam `could not lock … (Claude Code is writing the file)` नोंदवते.
- **MCP शेअर न करणारी खाती.** `import` नकार देते आणि आधी `ccam share side --mcp on` चालवायला सांगते.
- **गटांची नावे.** `-g` फक्त गटाचे नाव घेतो; ते `groups/` खाली फोल्डर बनते, त्यामुळे `../x` सारखी नावे नाकारली जातात.

## संबंधित

- [खात्यांमध्ये शेअरिंग](/features/share)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [आधीचा सेटअप कॉपी करा](/features/migrate)
- [सर्व कमांड्स](/docs/commands)
