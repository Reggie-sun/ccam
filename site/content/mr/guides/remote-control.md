---
title: ccam सोबत Remote Control वापरणे
description: ccam <name> rc ने कोणत्याही खात्याने Claude Code Remote Control सुरू करा; लागू पर्याय, शेअर केलेले MCP सर्व्हर कॉपी करण्याची सूचना आणि अडचणींवरील उपाय पाहा.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "4d46f49c5dea"
---

`ccam work rc` Claude Code चे Remote Control `work` खात्याने सुरू करते, त्यामुळे या कॉम्प्युटरवरचे sessions तुम्ही claude.ai/code किंवा Claude मोबाइल ॲपवरून नियंत्रित करू शकता. ते `claude rc` सारखेच काम करते, फक्त एका फरकाने: खाते MCP सर्व्हर शेअर करत असेल, तर ccam एकदा विचारते की ते खात्याच्या `.claude.json` मध्ये कॉपी करायचे का, कारण Remote Control sessions फक्त याच फाइलमधून सर्व्हर वाचतात.

## सुरू करा

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[शॉर्टकट](/features/shortcuts) सुद्धा असेच काम करतो: `ccwork rc` म्हणजे `ccam work rc`. ज्या फोल्डरमध्ये काम करायचे आहे तिथे कमांड चालवा, मग `work` वापरत असलेल्या त्याच Claude खात्याने साइन इन केलेले claude.ai/code किंवा Claude ॲप उघडा.

## पर्याय

इथे दोन प्रकारचे पर्याय एकत्र येतात: ccam चे स्वतःचे, जे खात्याच्या नावाआधी येतात किंवा ccam च्या सेटिंग्जमध्ये असतात, आणि Remote Control चे, जे `rc` नंतर येतात आणि न बदलता Claude Code पर्यंत पोहोचतात.

### ccam चे पर्याय आणि सेटिंग्ज

| पर्याय किंवा सेटिंग | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--dry-run` (खात्याच्या आधी) | environment आणि कमांड छापतो, आणि MCP सूचना असल्यास तीही. ना विचारतो, ना कॉपी करतो | बंद |
| `~/.config/ccam/config.toml` च्या सुरुवातीला `mcp_strategy = "sync"` | प्रत्येक वेळी सुरू करण्यापूर्वी, Remote Control असो वा नसो, गटाचे सर्व्हर `.claude.json` मध्ये कॉपी करतो, त्यामुळे कधीच काहीही विचारले जात नाही | `"flag"` |
| `ccam group leave <group> work` | खात्यासोबत MCP शेअर करणाऱ्या गटातून ते खाते बाहेर काढतो. गटाचे सर्व्हर त्याच्या `.claude.json` मध्ये त्याचे स्वतःचे म्हणून कॉपी होतात, त्यामुळे Remote Control कडे ते असतात आणि काहीही विचारले जात नाही; `--fresh` सह ते फक्त खात्याच्या स्वतःच्या सर्व्हरसह सुरू होते | `default` मध्ये |

खाते MCP सर्व्हर शेअर करत असेल आणि त्याच्या `.claude.json` मध्ये गटाची सध्याची यादी नसेल, किंवा कॉपी केलेला सर्व्हर त्यातून नाहीसा झाला असेल, तेव्हा ccam विचारते `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`. `y` ते कॉपी करतो; `N` किंवा Enter त्यांच्याशिवाय सुरू करतो.

### Remote Control चे पर्याय

`rc` नंतरचे सर्व काही Claude Code चे आहे. हे Claude Code 2.1.304 मधील `claude remote-control` चे पर्याय आहेत; तुमच्या आवृत्तीची यादी `ccam work rc --help` छापते.

| पर्याय | काय करतो |
| --- | --- |
| `--name <name>` | claude.ai/code मध्ये दिसणारे session चे नाव |
| `--remote-control-session-name-prefix <prefix>` | आपोआप तयार होणाऱ्या session नावांसाठी prefix (डीफॉल्ट: host चे नाव) |
| `-c`, `--continue` | या फोल्डरसाठी Remote Control ने शेवटी नोंदवलेल्या session ला, साधारण 4 तासांच्या आत, पुन्हा जोडतो |
| `--session-id <id>` | ठराविक session ला पुन्हा जोडतो |
| `--permission-mode <mode>` | नवीन sessions साठी permission mode: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` किंवा `plan` |
| `--spawn <mode>` | `same-dir` (डीफॉल्ट), `worktree` (प्रत्येक session साठी एक git worktree) किंवा `session` (एकच session, ते संपल्यावर बाहेर पडतो) |
| `--capacity <N>` | `same-dir` आणि `worktree` mode मध्ये एकाच वेळी चालणाऱ्या sessions ची कमाल संख्या (डीफॉल्ट 32) |
| `--[no-]create-session-in-dir` | सुरुवातीलाच सध्याच्या फोल्डरमध्ये पहिले session तयार करतो (डीफॉल्ट चालू) |
| `--[no-]chrome` | नवीन sessions साठी Claude in Chrome (डीफॉल्ट बंद) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | debug आणि सविस्तर output |

ccam Claude Code प्रमाणेच Claude Code ची कमांड वाचते, पर्याय त्यांच्या मूल्यांसह वगळून, त्यामुळे `ccam work rc --name laptop` आणि `ccam work --model opus rc` दोन्ही Remote Control सुरू करतात. ccam ला माहीत नसलेला पर्याय मूल्य नसलेला पर्याय म्हणून मोजला जातो; खात्रीसाठी, `rc` खात्याच्या नावानंतर लगेच ठेवा.

## ccam काय करते

सुरुवात नेहमीच्याच मार्गाने होते ([Claude Code सुरू करणे](/features/launch)), फक्त MCP टप्प्यात एक बदल असतो:

1. इतर कोणत्याही सुरुवातीप्रमाणे, ccam खाते शोधते, त्याच्या share links तपासते आणि `claude` शोधते.
2. कमांड `rc` किंवा `remote-control` आहे हे ओळखून ते `--mcp-config` वगळते.
3. खाते MCP सर्व्हर शेअर करत असेल, तर ते गटाच्या `mcp.json` ची `~/.local/state/ccam/mcp/work/` मधील `sync-stamp` शी तुलना करते, आणि गटाचा प्रत्येक सर्व्हर अजून `.claude.json` मध्ये आहे का ते तपासते. `.claude.json` मध्ये आधीच सध्याचे सर्व्हर असतील, किंवा गटात एकही सर्व्हर नसेल, तर ते थेट टप्पा 6 वर जाते.
4. नाहीतर sessions ना जे सर्व्हर मिळणार नाहीत त्यांची नावे ते सांगते:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. `y` वर ते `.claude.json` चा बॅकअप जतन करते, Claude Code चे स्वतःचे lock घेते, गटाचे सर्व्हर `mcpServers` मध्ये जोडते आणि खात्याचे इतर सर्व्हर तसेच ठेवते; खात्याकडे गटातील सर्व्हरच्याच नावाचा स्वतःचा सर्व्हर असेल, तर ccam खात्याचा सर्व्हर ठेवते आणि तसे सांगते. `N` वर, टर्मिनल नसल्यास, किंवा खाते चालू असताना, ते काहीही बदलत नाही आणि तसे सांगते:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. ते session नोंदवते, `CLAUDE_CONFIG_DIR` खात्याच्या फोल्डरवर सेट करते आणि `exec` द्वारे स्वतःला `claude rc` ने बदलते.

त्यानंतर Remote Control सर्व्हर म्हणून चालू राहते. ते तुमच्यासाठी सुरू करत असलेले प्रत्येक session त्या process चे child असते, त्यामुळे ते `work` म्हणून चालते: `work` चे लॉगिन, सेटिंग्ज आणि संभाषण इतिहास वापरून, आणि ते `work` च्या वापर मर्यादेत मोजले जाते. ccam ने सुरू केलेला Remote Control सर्व्हर चालू असेपर्यंत ccam खाते चालू असल्याचे मानते. साध्या `claude remote-control` ने सुरू केलेला सर्व्हर त्याला sessions असतील तेव्हाच मोजला जातो; अजून एकही session नसलेला असा सर्व्हर ccam ला दिसत नाही.

dry run हा फरक दाखवतो. सामान्य सुरुवात गटाचे सर्व्हर आर्ग्युमेंट म्हणून सोबत नेते, Remote Control नेत नाही:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## हे असे का काम करते

**Remote Control कोणताही `--mcp-config` घेत नाही.** सामान्य सुरुवातीला ccam गटाचे सर्व्हर Claude Code ला `--mcp-config=<path>` म्हणून देते आणि Claude Code च्या फाइल्समध्ये कधीच लिहीत नाही. `claude remote-control` त्याच्या कमांड नावाआधी तो पर्याय नाकारतो आणि नंतर एकही स्वीकारत नाही, त्यामुळे ccam 0.6.0 पूर्वी MCP सर्व्हर शेअर करणाऱ्या प्रत्येक खात्यासाठी `ccam work rc` अयशस्वी होत असे. आता ccam तो पर्याय वगळते.

**त्याचे sessions सर्व्हर फक्त `.claude.json` मधून वाचतात.** Remote Control सुरू करत असलेल्या sessions ना कमांड लाइनचे पर्याय वारशाने मिळत नाहीत. user MCP सर्व्हर त्यांना एकाच ठिकाणी मिळतात: खात्याच्या `.claude.json` मध्ये. त्या फाइलमध्ये खात्याचे लॉगिन तपशील आणि प्रत्येक फोल्डरची स्थितीही असते, त्यामुळे `settings.json` प्रमाणे ccam ती शेअर केलेल्या प्रतीशी link करू शकत नाही. sessions ना शेअर केलेली यादी देण्याचा एकमेव मार्ग म्हणजे सर्व्हर त्या फाइलमध्ये कॉपी करणे.

**लिहिण्यापूर्वी ccam विचारते.** डीफॉल्ट flag strategy मध्ये ccam सुरुवातीला `.claude.json` ला हात लावत नाही. Remote Control ही एकमेव सुरुवात आहे जिथे त्याला तिथे लिहावेच लागते, म्हणून ते आधी विचारते, बॅकअप ठेवते, आणि `mcpServers` शिवाय कशालाही हात लावत नाही. (सुरुवातीबाहेर, MCP शेअरिंग सोडणे किंवा त्यात पुन्हा सामील होणे आणि `ccam mcp import` सुद्धा तिथे लिहितात, त्यांच्या योजना सांगतात तसे.)

**खाते चालू असताना ते लिहीत नाही.** session चालू असताना Claude Code `.claude.json` पुन्हा लिहिते, आणि बाहेरून केलेला बदल हरवेल किंवा त्यावर पुन्हा लिहिले जाईल. म्हणूनच खात्याचे कोणतेही session चालू असताना, दुसऱ्या Remote Control सर्व्हरसह आणि ccam शिवाय सुरू केलेल्या Claude Code सह, ccam कॉपी वगळते.

**प्रत्येक बदलासाठी ते एकदाच विचारते.** कॉपीनंतर `sync-stamp` मध्ये गटाच्या `mcp.json` ची modification time असते. ती फाइल बदलली, उदाहरणार्थ `ccam mcp add` नंतर, किंवा कॉपी केलेला सर्व्हर `.claude.json` मधून नाहीसा झाला, तरच ccam पुन्हा विचारते. तुम्ही `N` उत्तर दिल्यास अजून कॉपी झालेली नसते, त्यामुळे ते प्रत्येक वेळी Remote Control सुरू करताना विचारते.

**`ccam mcp` कॉपी अद्ययावत ठेवते.** नाव जुळल्यास `.claude.json` मधील नोंद `--mcp-config` वर जिंकते, त्यामुळे जुनी कॉपी सामान्य सुरुवातीलाही नवीन व्याख्या झाकून टाकेल. म्हणून `ccam mcp add`, `add-json`, `remove` आणि `import`, आणि `ccam migrate`, चालू नसलेल्या प्रत्येक खात्याची कॉपी अपडेट करतात, आणि एखादी कॉपी बदलली तर तसे सांगतात:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**टर्मिनल नसल्यास उत्तर "नाही" असते.** Remote Control अनेकदा script, service किंवा `nohup` मधून चालते, जिथे उत्तर द्यायला कोणी नसते. तेव्हा ccam वाट न पाहता सर्व्हरशिवाय सुरू करते.

## अडचणी आणि उपाय

### Remote Control sessions मध्ये शेअर केलेले MCP सर्व्हर दिसत नाहीत

प्रश्नाला `N` उत्तर दिले होते, टर्मिनल नव्हते, किंवा खाते चालू होते. Remote Control थांबवा, खात्याचे इतर sessions बंद करा, आणि टर्मिनलमध्ये `ccam work rc` पुन्हा चालवा; `y` उत्तर द्या. पुन्हा विचारले जाऊ नये म्हणून `~/.config/ccam/config.toml` मध्ये `mcp_strategy = "sync"` सेट करा. session मध्ये आधीपासून असलेले सर्व्हर ते चालू असताना बदलत नाहीत, त्यामुळे कॉपीनंतर नवीन session सुरू करा.

### "Account work is running, so its .claude.json is left alone"

`work` चे एक session अजून उघडे आहे: सामान्य `ccam work`, शॉर्टकट, किंवा दुसरा Remote Control सर्व्हर. त्यातून बाहेर पडा आणि `ccam work rc` पुन्हा चालवा. दुसरे खाते चालवत असलेले sessions मोजले जात नाहीत.

### ccam प्रत्येक वेळी विचारते

तुम्ही अजून सर्व्हर कॉपी केलेले नाहीत. एकदा `y` उत्तर द्या, किंवा `mcp_strategy = "sync"` सेट करा. Remote Control मध्ये सर्व्हर कधीच नको असतील, तर ते न ठेवता खात्यासोबत MCP शेअर करणे थांबवा: `ccam group leave default work --fresh`, किंवा MCP अशा स्वतंत्र गटात ठेवा ज्यात ते खाते नाही. साधे `ccam group leave` गटाचे सर्व्हर खात्याच्या `.claude.json` मध्ये कॉपी करते, त्यामुळे Remote Control कडे ते असतात.

### `ccam work rc` लगेचच `--mcp-config` बद्दलच्या त्रुटीसह अयशस्वी होते

तुमचे ccam 0.6.0 पेक्षा जुने आहे आणि Remote Control ला `--mcp-config` देते. `ccam update` चालवा. सध्याच्या ccam सोबत हीच त्रुटी आली, तर याचा अर्थ अनोळखी, मूल्य घेणारा एखादा पर्याय `rc` च्या आधी आला होता, आणि ccam ने चुकून ते मूल्य कमांड समजले; `rc` खात्याच्या नावानंतर लगेच ठेवा.

### `ccam mcp` ने केलेला बदल Remote Control पर्यंत पोहोचत नाही

`ccam mcp add` आणि `remove` फक्त चालू नसलेल्या खात्यांमधील कॉपी अपडेट करतात, आणि तसे सांगतात: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control थांबवून पुन्हा सुरू करा; ccam ला बदललेली `mcp.json` दिसते आणि ते पुन्हा विचारते.

### गटाच्या सर्व्हरऐवजी माझा स्वतःचा सर्व्हर वापरला जातो

खात्याकडे गटातील सर्व्हरच्याच नावाचा स्वतःचा सर्व्हर असेल, तर कॉपी तो जागेवरच ठेवते आणि तसे सांगते: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`. गटाचा सर्व्हर वापरायचा असल्यास खात्याचा सर्व्हर काढा (`claude mcp remove -s user <name>`); तुमची आवृत्ती प्रत्येक खात्याला द्यायची असल्यास `ccam mcp import work <server> --force` ने ती गटात हलवा. `.claude.json` चे बॅकअप `~/.local/state/ccam/mcp/work/backups/` मध्ये आहेत, जिथे सर्वात नवीन पाच ठेवले जातात.

### मला कॉपी केलेले सर्व्हर काढून टाकायचे आहेत

`ccam mcp remove <server>` सर्व्हर गटातून आणि प्रत्येक कॉपीतून काढून टाकते. कॉपीपूर्वीच्या स्थितीत परत जाण्यासाठी, खात्याचे कोणतेही session चालू नसताना `~/.local/state/ccam/mcp/work/backups/` मधील सर्वात नवीन बॅकअप पुनर्संचयित करा.

### Remote Control फोल्डरवर विश्वास ठेवायचा का असे विचारते, किंवा टर्मिनल नसताना बाहेर पडते

फोल्डरवरचा विश्वास प्रत्येक खात्यासाठी स्वतंत्रपणे `.claude.json` मध्ये साठवला जातो, त्यामुळे `personal` म्हणून तुम्ही विश्वास ठेवलेला फोल्डर `work` साठी नवीन असतो. त्या फोल्डरमध्ये एकदा `ccam work` चालवून स्वीकारा, मग तिथे Remote Control सुरू करा.

### "You must be logged in with a Claude account that has a subscription"

Remote Control ला Claude subscription असलेले लॉगिन लागते. `ccam list` मध्ये खात्याचा ईमेल तपासा, आणि गरज असल्यास `ccam login work` ने पुन्हा साइन इन करा.

### session फोनवर दिसत नाही

Claude ॲप किंवा claude.ai/code वेगळ्या Claude खात्याने साइन इन केलेले आहे. `work` साठी `ccam list` दाखवत असलेल्या ईमेलशी त्याची तुलना करा.

### या फोल्डरमध्ये दुसरे Claude Code session संभाषण इतिहास शेअर करत असल्याची ccam चेतावणी देते

त्याच इतिहास गटातील दुसरे खाते त्याच फोल्डरमध्ये चालू आहे. दोन्ही खात्यांमधून एकच संभाषण पुढे चालवणारे sessions इतिहासात गोंधळ करू शकतात. ते वेगवेगळ्या संभाषणांवर काम करणार असतील तरच `y` उत्तर द्या, किंवा Remote Control दुसऱ्या फोल्डरमध्ये सुरू करा.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [शेअर केलेले MCP सर्व्हर](/features/mcp)
- [गटांसह शेअरिंग](/features/share)
- [खाती व्यवस्थापित करणे](/features/accounts)
