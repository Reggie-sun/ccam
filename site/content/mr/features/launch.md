---
title: Claude Code सुरू करणे
description: ccam <name|number>, ccam run किंवा ccam picker ने कोणत्याही खात्याने Claude Code सुरू करा; आर्ग्युमेंट्स न बदलता claude ला जातात, आणि --dry-run कमांड दाखवते.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work` `work` खात्याने Claude Code सुरू करते, आणि नावानंतरचे सर्व काही जसेच्या तसे `claude` ला जाते. जिथे तुम्ही `claude` टाइप केले असते तिथे ते वापरा: टर्मिनलमध्ये, script मध्ये, किंवा pipe ने input देऊन.

## वापर

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

खात्याचे नाव किंवा `ccam list` मधील त्याचा क्रमांक द्या. नावानंतरचे आर्ग्युमेंट्स ccam वाचत नाही, त्यामुळे `ccam work --help` Claude Code ची मदत दाखवते. नावानंतर `--` लावणे ऐच्छिक आहे. standard input, exit code आणि signals `claude` प्रमाणेच वागतात.

पहिला आर्ग्युमेंट `-` ने सुरू होत असेल आणि तो ccam च्या स्वतःच्या flags पैकी नसेल, तर ccam सर्व आर्ग्युमेंट्ससह नावाशिवाय एक खाते सुरू करते. ते म्हणजे [`ccam folder`](/features/folder) ने सध्याच्या फोल्डरसाठी सेट केलेले खाते, नाहीतर डीफॉल्ट: `ccam use` ने सेट केलेले खाते, किंवा तुम्ही ते सेट केले नसेल तर सर्वात लहान क्रमांकाचे खाते.

### Picker

standard input आणि output टर्मिनल असताना आर्ग्युमेंट्सशिवाय `ccam` picker उघडते, आणि इतर वेळी `ccam list` छापते. ते प्रत्येक खात्याचा क्रमांक, ईमेल, 5 तासांचा वापर, शेवटचा वापर आणि शेअर केलेले आयटम दाखवते; प्रत्येक खाते उत्तर देईल तसा वापर भरला जातो.

| की | काय करते |
| --- | --- |
| `Enter` | निवडलेल्या खात्याने Claude Code सुरू करा |
| `1` to `9` | त्या क्रमांकाचे खाते सुरू करा |
| `Up`/`Down`, `k`/`j` | निवड हलवा |
| `a` | picker बंद करून `ccam add` विझार्ड चालवा |
| `s` | निवडलेल्या खात्याचे शेअर केलेले आयटम बदला (`Space` किंवा `x` टॉगल करते, `Enter` लागू करते, `Esc` मागे जाते) |
| `d` | निवडलेले खाते डीफॉल्ट करा |
| `r` | छोट्या cache शिवाय वापर पुन्हा आणा |
| `q`, `Esc`, `Ctrl+C` | बाहेर पडा |

picker काही करण्यापूर्वी बंद होतो, त्यामुळे Claude Code सुरू होण्यापूर्वी टर्मिनल पूर्ववत होतो. `d` किंवा शेअरिंगमधील बदलानंतर तो पुन्हा उघडतो. शेअरिंगमधील बदल `ccam share` प्रमाणेच आधी योजना दाखवतो आणि विचारतो.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `--dry-run` | ccam चालवेल ते environment आणि कमांड छापा, मग बाहेर पडा | बंद |
| `--fallback` | `-p`/`--print` runs साठी: खाते आपल्या वापर मर्यादेला पोहोचले की तीच कमांड सर्वात जास्त जागा असलेल्या लॉगिन केलेल्या खात्याने पुन्हा चालवा ([तपशील](/features/next#run-again-at-the-limit)) | बंद |
| `-v`, `--version` | `ccam <version>` छापा (फक्त पहिला आर्ग्युमेंट म्हणून) | |
| `-h`, `--help` | ccam ची मदत, `ccam --help` किंवा `ccam run --help` म्हणून | |

आर्ग्युमेंट्सचे नियम:

- ccam आपले स्वतःचे flags फक्त खात्याच्या नावाआधी वाचते. नावानंतरचा प्रत्येक आर्ग्युमेंट Claude Code चा असतो, `--dry-run` आणि `--help` सुद्धा.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work` आणि `ccam --dry-run` (नावाशिवायचे खाते) असे चालते. तो पाळला जातो, नाकारला जात नाही. `--dry-run` आणि `--fallback` कोणत्याही क्रमाने येऊ शकतात.
- `ccam run <name|number> [claude args]` हे `ccam <name>` चे लांब रूप आहे. नावाशिवाय `ccam run` फोल्डरचे खाते किंवा डीफॉल्ट खाते वापरते.
- खात्यांची नावे `add`, `list`, `run` किंवा `usage` सारखी ccam कमांड्सची नावे असू शकत नाहीत, त्यामुळे `ccam <word>` कधीच संदिग्ध होत नाही.

dry run एक ओळ छापतो, जी तुम्ही shell मध्ये, किंवा Windows वर PowerShell मध्ये, paste करू शकता:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## हे कसे काम करते

प्रत्येक वेळी सुरू करताना, कमांड लाइनपासून `exec` पर्यंत तेच टप्पे होतात:

<!-- motion:launch -->

प्रत्येक टप्पा चालताना त्याचा box उजळतो, आणि खालच्या ओळी म्हणजे Claude Code ज्या environment आणि आर्ग्युमेंट्ससह सुरू होते ते.

1. `cmd/ccam/main.go` `app.Main` बोलावते. ccam दुसऱ्या नावाने सुरू झाले असेल, तर तो [शॉर्टकट](/features/shortcuts) आहे आणि थेट launcher कडे जातो.
2. cobra सहभागी होण्यापूर्वी `app.parse` `os.Args` पाहते. फक्त एकटे `ccam` (picker), ccam कमांडचे नाव (`config.Reserved`), `-h`/`--help` आणि `-v`/`--version` `cli.Execute` कडे जातात. बाकी सर्व काही म्हणजे Claude Code सुरू करणे.
3. `launcher.Launch` `~/.config/ccam/config.toml` वाचते, ही एकमेव फाइल जी ते parse करते.
4. `launcher.Prepare` नाव किंवा क्रमांकावरून खाते शोधते (नाव नसल्यास `HereProfile` फोल्डरचे खाते घेते, working directory ची `config.toml` मधील `[[folder]]` नोंदींशी strings म्हणून तुलना करून), सुरुवातीचा एक `--` काढते, आणि `ResolveClaude` ने `claude` executable शोधते. PATH मधील शोध फाइलच्या modification time सह `~/.local/state/ccam/claude-path` मध्ये cache होतो, त्यामुळे ccam प्रत्येक वेळी `PATH` शोधत नाही. फाइल बदलली किंवा नसली तर नवीन शोध होतो. `config.toml` मधील `claude_bin` पाथ जसाच्या तसा वापरला जातो. जो `claude` स्वतः ccam आहे, तो ccam inode ने तुलना करून नाकारते, कारण त्याने loop होईल.
5. `checkLinks` खात्याचा फोल्डर `~/.local/share/ccam/profiles/<name>/` त्याच्या शेअरिंग सेटिंग्जशी `share.Check` ने तुलना करते, जे फक्त `lstat` आणि `readlink` बोलावते, प्रत्येक आयटमसाठी एक किंवा दोन system calls. नसलेली लिंक तिथल्या तिथे तयार होते. डेटा गमावू शकेल अशी कोणतीही गोष्ट सुरू करणे थांबवते.
6. `applyMCP` शेअर केलेले MCP सर्व्हर जोडते (खाली पाहा).
7. `session.Register` session नोंदवते. ccam आपल्या PID च्या नावाचा, process सुरू झाल्याची वेळ असलेला marker `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` मध्ये लिहिते (`XDG_RUNTIME_DIR` नसल्यास `~/.local/state/ccam/locks/` मध्ये). संभाषण इतिहास शेअर करणाऱ्या खात्यासाठी, `historyCheck` आधी `history-<group>-<hash>` मध्ये चालू session शोधते, जिथे hash सध्याच्या फोल्डरचा असतो, आणि तिथेही नोंदवते. `exec` PID तोच ठेवते, त्यामुळे marker नंतर Claude Code process चेच नाव सांगतो. धरलेले lock नाही तर marker, कारण Claude Code सुरू करत असलेल्या processes ना lock वारसाने मिळेल आणि session संपल्यानंतरही ते राहील.
8. `Plan.Exec` टर्मिनलचे title `claude · <name>` करते (standard error टर्मिनल असताना), सुरू केल्याची वेळ `fsync` शिवाय `~/.local/state/ccam/last-used/<name>` मध्ये लिहिते, आणि `syscall.Exec` बोलावते.

Claude Code खात्याच्या फोल्डरवर सेट केलेल्या `CLAUDE_CONFIG_DIR` सह आणि खात्याच्या नावावर सेट केलेल्या `CCAM_PROFILE` सह चालते. यांची वारसाने मिळालेली कोणतीही मूल्ये ccam काढून टाकते. macOS वर `ccam import` ने नोंदवलेल्या `~/.claude` फोल्डरसाठी ते `CLAUDE_SECURESTORAGE_CONFIG_DIR` रिकाम्या मूल्यावर सेट करते, जेणेकरून ते खाते आणि साधे `claude` एकच Keychain लॉगिन वापरतील.

### exec का

`syscall.Exec` ccam process च्या जागी `claude` ठेवते, त्यामुळे ccam memory मध्ये राहत नाही, keystrokes, signals आणि exit code Claude Code चेच असतात, आणि आर्ग्युमेंट्स array राहतात, कधीच shell string बनत नाहीत.

Windows मध्ये `exec` नाही. तिथे `Plan.Exec` `claude.exe` ला त्याच environment, आर्ग्युमेंट्स आणि console सह child process म्हणून सुरू करते, त्याची वाट पाहते आणि त्याच्याच exit code ने बाहेर पडते. Ctrl+C console मधून थेट Claude Code पर्यंत पोहोचतो, त्यामुळे ccam त्याकडे दुर्लक्ष करते. session marker मध्ये ccam चा PID असतो, जो Claude Code इतकाच काळ जिवंत राहतो, आणि Claude Code बाहेर पडल्यावर marker काढला जातो. वाट पाहत असताना ccam दुसरे काहीही करत नाही. session नंतरचे काम, जसे `claude mcp add` ने जोडलेले सर्व्हर लक्षात घेणे, पुढच्या वेळी सुरू करताना होते: ccam `.claude.json` च्या modification time ची `~/.local/state/ccam/mcp/<name>/seen` शी तुलना करते आणि फाइल बदलली असेल तरच ती parse करते.

### कमांड लाइन कुठे जाते

तीन कमांड लाइन्स, `internal/app` मधून तीन मार्ग:

<!-- motion:launch-route -->

शॉर्टकट `argv[0]` वरून ओळखला जातो, `ccam <name>` पहिल्या आर्ग्युमेंटवरून, आणि दोन्ही `syscall.Exec` मध्ये संपतात. फक्त `ccam usage` सारखी ccam कमांड cobra पर्यंत पोहोचते.

### शेअर केलेले MCP सर्व्हर

डीफॉल्ट `flag` strategy सह, ccam गटाचे MCP सर्व्हर Claude Code ला एक आर्ग्युमेंट म्हणून, तुमच्या आर्ग्युमेंट्सच्या आधी देते:

<!-- motion:launch-mcp -->

गटाची `mcp.json` असेल तरच `mcp.FlagArgs` आर्ग्युमेंट परत देते, आणि `Plan.InsertArgs` तो `claude` च्या लगेच नंतर ठेवते. तुमचे स्वतःचे `--mcp-config` किंवा `--strict-mcp-config` त्यावर लागू होत राहतात. तो `--mcp-config=<path>` हा एकच आर्ग्युमेंट आहे, कारण Claude Code चा `--mcp-config` अनेक मूल्ये घेतो: दोन आर्ग्युमेंट्स म्हणून लिहिल्यास तो पुढचा prompt सुद्धा config फाइल म्हणून वाचेल.

`config.toml` मध्ये `mcp_strategy = "sync"` असल्यास, ccam त्याऐवजी सर्व्हर खात्याच्या `.claude.json` मध्ये लिहिते, आणि खाते आधीच चालू असेल तर ते वगळते.

### Remote Control

`ccam work rc` (किंवा `ccam work remote-control`) Claude Code चे Remote Control `work` म्हणून सुरू करते. `claude remote-control` त्याच्या verb च्या आधी `--mcp-config` नाकारतो आणि नंतर एकही घेत नाही, आणि त्याचे sessions MCP सर्व्हर फक्त खात्याच्या `.claude.json` मधून वाचतात. म्हणून इथे flag strategy वेगळ्या प्रकारे काम करते:

<!-- motion:launch-rc -->

1. `remoteControl` खात्यानंतरचा `-` ने सुरू न होणारा पहिला आर्ग्युमेंट Claude Code चा verb मानते. तो `remote-control` किंवा `rc` असल्यास, `mcp.FlagArgs` काहीही जोडत नाही.
2. `mcp.Synced` गटाच्या `mcp.json` ची modification time खात्याच्या `sync-stamp` शी तुलना करते. `.claude.json` मध्ये आधीच गटाचे सध्याचे सर्व्हर असतील, किंवा गटात एकही सर्व्हर नसेल, तर Remote Control प्रश्न न विचारता सुरू होते.
3. नाहीतर sessions ना जे सर्व्हर मिळणार नाहीत त्यांची नावे ccam सांगते आणि विचारते:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` वर, `mcp.Sync` ते `sync` strategy प्रमाणेच कॉपी करते: Claude Code च्या lock खाली, बॅकअप घेतल्यानंतर, खात्याचे स्वतःचे सर्व्हर तसेच ठेवून. `N` वर, किंवा टर्मिनल नसल्यास, Remote Control त्यांच्याशिवाय सुरू होते, आणि ccam `mcp_strategy = "sync"` सुचवते, जे प्रत्येक वेळी सुरू करण्यापूर्वी ते कॉपी करते.

सर्व्हर कॉपी होईपर्यंत ccam प्रत्येक वेळी Remote Control सुरू करताना विचारते, आणि कॉपीनंतर फक्त गटाची `mcp.json` पुन्हा बदलल्यावरच. नंतर सामान्यपणे सुरू करताना `.claude.json` मधील कॉपी `--mcp-config` वर जिंकते, म्हणून `ccam mcp add`, `remove` आणि `import` इतर खात्यांकडील कॉपी अपडेट करतात ([शेअर केलेले MCP सर्व्हर](/features/mcp)).

### 10 ms ची मर्यादा

Claude Code सुरू होण्यापूर्वी ccam जास्तीत जास्त 10 ms (p95) जोडते. सुरू करण्याचा मार्ग कधीच cobra कमांड tree तयार करत नाही, Bubble Tea सुरू करत नाही किंवा नेटवर्क कनेक्शन उघडत नाही, आणि package च्या `init()` functions ना एकूण 5 ms पेक्षा जास्त लागल्यास `TestPackageInitIsCheap` अयशस्वी होते. `scripts/perf-check.sh` `hyperfine` ने `ccam --dry-run work`, `ccam work` आणि एक शॉर्टकट मोजते आणि 10 ms पेक्षा जास्त झाल्यास अयशस्वी होते, आणि `internal/launcher` मधील `BenchmarkLaunchPath` हे टप्पे loop मध्ये चालवते. प्रत्येक टप्प्याचा वेळ पाहण्यासाठी `CCAM_TRACE=1` सेट करा:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## सुरक्षितता आणि अपवादात्मक स्थिती

- अज्ञात नाव किंवा क्रमांक `no account "<name>"` ने थांबतो आणि `ccam list` कडे निर्देश करतो. एकही खाते नसल्यास, ccam `ccam add <name>` ने एक जोडायला सांगते.
- एखादा शेअर केलेला पाथ साधी फाइल बनला असेल किंवा चुकीच्या जागी निर्देश करत असेल, तर ccam `stopped because of share link problems:` ने थांबते, त्यांची यादी देते, आणि `ccam doctor --fix` चालवायला सांगते. सुरू करण्याच्या मार्गावर ते डेटा हलवत नाही.
- या खात्याशी संभाषण इतिहास शेअर करणारे (त्याच गटातील) दुसरे Claude Code त्याच फोल्डरमध्ये चालू असेल, तर ccam इतिहास मिसळू शकतो असा इशारा देते आणि `Run anyway? [y/N]` विचारते. टर्मिनल नसताना ते इशारा दाखवते आणि पुढे जाते.
- `--dry-run` MCP sync, नवीन सर्व्हरची सूचना, इतिहासाची तपासणी, session marker, title आणि शेवटच्या वापराची वेळ वगळते, आणि exec करत नाही. तरीही ते शेअर लिंक्स तपासते, आणि सामान्यपणे सुरू करताना होईल तशी नसलेली लिंक तयार करते.
- `ccam work rc` हे `work` आधीच चालू असताना `.claude.json` ला हात लावत नाही, तसे सांगते, आणि गटाच्या सर्व्हरशिवाय Remote Control सुरू करते; ते कॉपी करण्यासाठी session संपल्यानंतर ते पुन्हा चालवा. `--dry-run` सह ते संदेश छापते पण विचारत नाही आणि कॉपीही करत नाही.
- `mcp_strategy = "sync"` सह sync अयशस्वी झाल्यास `ccam: MCP sync failed (launching anyway): …` छापले जाते आणि Claude Code त्याच्याकडे आधीपासून असलेल्या सर्व्हरसह सुरू होते.

## संबंधित

- [शॉर्टकट कमांड्स](/features/shortcuts)
- [ccam सोबत Remote Control वापरणे](/guides/remote-control)
- [खाती व्यवस्थापित करणे](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [CLAUDE_CONFIG_DIR कसे काम करते](/guides/claude-config-dir)
