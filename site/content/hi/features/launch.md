---
title: Claude Code शुरू करना
description: ccam <name|number>, ccam run या ccam picker से Claude Code को किसी भी अकाउंट से शुरू करें; आर्ग्युमेंट बिना बदलाव के claude को जाते हैं, और --dry-run कमांड दिखाता है।
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

`ccam work` Claude Code को `work` अकाउंट से शुरू करता है, और नाम के बाद का सब कुछ जैसा है वैसा ही `claude` को जाता है। इसे वहाँ इस्तेमाल करें जहाँ आप `claude` टाइप करते: टर्मिनल में, किसी स्क्रिप्ट में, या pipe से आए input के साथ।

## उपयोग

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

अकाउंट का नाम या `ccam list` से उसका नंबर दें। ccam नाम के बाद वाले आर्ग्युमेंट नहीं पढ़ता, इसलिए `ccam work --help` Claude Code की help दिखाता है। नाम के बाद `--` लगाना वैकल्पिक है। Standard input, exit code और signals वैसे ही काम करते हैं जैसे `claude` के साथ।

जब पहला आर्ग्युमेंट `-` से शुरू होता है और ccam के अपने flags में से नहीं है, तो ccam सभी आर्ग्युमेंट के साथ डिफ़ॉल्ट अकाउंट शुरू करता है। डिफ़ॉल्ट वह अकाउंट है जिसे `ccam use` से सेट किया गया, या अगर आपने कोई सेट नहीं किया तो सबसे छोटे नंबर वाला अकाउंट।

### Picker

बिना आर्ग्युमेंट के `ccam`, standard input और output के टर्मिनल होने पर picker खोलता है, वरना `ccam list` प्रिंट करता है। यह हर अकाउंट का नंबर, ईमेल, 5 घंटे का उपयोग, आख़िरी उपयोग और साझा आइटम दिखाता है; हर अकाउंट के जवाब देते ही उपयोग भरता जाता है।

| Key | क्या करती है |
| --- | --- |
| `Enter` | चुने गए अकाउंट से Claude Code शुरू करें |
| `1` से `9` | उस नंबर वाला अकाउंट शुरू करें |
| `Up`/`Down`, `k`/`j` | चयन ऊपर-नीचे करें |
| `a` | Picker बंद करें और `ccam add` विज़र्ड चलाएँ |
| `s` | चुने गए अकाउंट के साझा आइटम बदलें (`Space` या `x` टॉगल करता है, `Enter` लागू करता है, `Esc` वापस जाता है) |
| `d` | चुने गए अकाउंट को डिफ़ॉल्ट बनाएँ |
| `r` | छोटी cache के बिना उपयोग फिर से लाएँ |
| `q`, `Esc`, `Ctrl+C` | बाहर निकलें |

Picker कुछ करने से पहले बंद होता है, इसलिए Claude Code शुरू होने से पहले टर्मिनल अपनी सामान्य स्थिति में आ जाता है। `d` या साझा करने में बदलाव के बाद यह फिर से खुलता है। साझा करने में बदलाव पहले योजना दिखाता है और पूछता है, जैसे `ccam share` करता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--dry-run` | वह environment और कमांड प्रिंट करें जो ccam चलाता, फिर बाहर निकलें | बंद |
| `-v`, `--version` | `ccam <version>` प्रिंट करें (केवल पहले आर्ग्युमेंट के रूप में) | |
| `-h`, `--help` | ccam की help, `ccam --help` या `ccam run --help` के रूप में | |

आर्ग्युमेंट के नियम:

- ccam अपने flags केवल अकाउंट के नाम से पहले पढ़ता है। नाम के बाद हर आर्ग्युमेंट Claude Code का है, `--dry-run` और `--help` भी।
- `--dry-run`, `ccam --dry-run work`, `ccam run --dry-run work` और `ccam --dry-run` (डिफ़ॉल्ट अकाउंट) के रूप में काम करता है। इसे माना जाता है, अस्वीकार नहीं किया जाता।
- `ccam run <name|number> [claude args]`, `ccam <name>` का लंबा रूप है। बिना नाम के `ccam run` डिफ़ॉल्ट अकाउंट इस्तेमाल करता है।
- अकाउंट के नाम `add`, `list`, `run` या `usage` जैसे ccam कमांड के नाम नहीं हो सकते, इसलिए `ccam <word>` कभी अस्पष्ट नहीं होता।

Dry run एक पंक्ति प्रिंट करता है जिसे आप shell में paste कर सकते हैं:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## यह कैसे काम करता है

हर बार शुरू करने में, कमांड लाइन से `exec` तक, वही चरण होते हैं:

<!-- motion:launch -->

हर बॉक्स अपना चरण चलते समय जलता है, और नीचे की पंक्तियाँ वह environment और आर्ग्युमेंट हैं जिनके साथ Claude Code शुरू होता है।

1. `cmd/ccam/main.go`, `app.Main` को बुलाता है। अगर ccam किसी दूसरे नाम से शुरू हुआ, तो वह एक [शॉर्टकट](/features/shortcuts) है और सीधे launcher पर जाता है।
2. `app.parse`, cobra के शामिल होने से पहले `os.Args` देखता है। केवल अकेला `ccam` (picker), कोई ccam कमांड नाम (`config.Reserved`), `-h`/`--help` और `-v`/`--version` ही `cli.Execute` पर जाते हैं। बाकी सब कुछ शुरू करना है।
3. `launcher.Launch`, `~/.config/ccam/config.toml` पढ़ता है, जो इकलौती फ़ाइल है जिसे यह parse करता है।
4. `launcher.Prepare` नाम या नंबर से अकाउंट खोजता है, शुरुआत का एक `--` हटाता है, और `ResolveClaude` से `claude` executable खोजता है। PATH की खोज फ़ाइल के modification time के साथ `~/.local/state/ccam/claude-path` में cache होती है, इसलिए ccam हर बार शुरू होने पर `PATH` नहीं खोजता। बदली हुई या गायब फ़ाइल नई खोज करवाती है। `config.toml` में दिया `claude_bin` पाथ जैसा है वैसा इस्तेमाल होता है। ccam ऐसे `claude` को अस्वीकार करता है जो ख़ुद ccam हो, inode से तुलना करके, क्योंकि उससे loop बन जाता।
5. `checkLinks`, अकाउंट फ़ोल्डर `~/.local/share/ccam/profiles/<name>/` की तुलना उसकी साझा सेटिंग्स से `share.Check` के ज़रिए करता है, जो केवल `lstat` और `readlink` बुलाता है, यानी हर आइटम के लिए एक या दो system call। गायब लिंक तुरंत बना दिया जाता है। जिस भी चीज़ से डेटा खो सकता है, वह शुरू करना रोक देती है।
6. `applyMCP` साझा MCP सर्वर जोड़ता है (नीचे देखें)।
7. `session.Register` session दर्ज करता है। ccam अपने PID के नाम वाला एक marker, जिसमें process के शुरू होने का समय होता है, `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` में लिखता है (या `XDG_RUNTIME_DIR` न होने पर `~/.local/state/ccam/locks/` में)। बातचीत का इतिहास साझा करने वाले अकाउंट के लिए, `historyCheck` पहले `history-<group>-<hash>` में कोई चल रहा session खोजता है, जहाँ hash मौजूदा फ़ोल्डर का है, और वहाँ भी रजिस्टर करता है। `exec` PID बनाए रखता है, इसलिए marker फिर ख़ुद Claude Code process का नाम बताता है। यह marker है, पकड़ा हुआ lock नहीं, क्योंकि Claude Code जो processes शुरू करता है वे lock विरासत में पाकर session ख़त्म होने के बाद भी उसे रखे रहते।
8. `Plan.Exec` टर्मिनल का शीर्षक `claude · <name>` सेट करता है (जब standard error टर्मिनल हो), शुरू करने का समय `fsync` के बिना `~/.local/state/ccam/last-used/<name>` में लिखता है, और `syscall.Exec` बुलाता है।

Claude Code, अकाउंट फ़ोल्डर पर सेट `CLAUDE_CONFIG_DIR` और अकाउंट के नाम पर सेट `CCAM_PROFILE` के साथ चलता है। ccam इनके विरासत में मिले मान हटा देता है। macOS पर `ccam import` से रजिस्टर किए गए `~/.claude` फ़ोल्डर के लिए, यह `CLAUDE_SECURESTORAGE_CONFIG_DIR` को ख़ाली मान पर भी सेट करता है, ताकि अकाउंट और सादा `claude` एक ही कीचेन लॉगिन इस्तेमाल करें।

### exec क्यों

`syscall.Exec`, ccam process की जगह `claude` को रख देता है, इसलिए ccam memory में नहीं रहता, keystrokes, signals और exit code Claude Code के अपने होते हैं, और आर्ग्युमेंट एक array बने रहते हैं, कभी shell string नहीं बनते। Session के बाद का काम, जैसे `claude mcp add` से जोड़े गए सर्वर पहचानना, अगली बार शुरू होने पर होता है: ccam `.claude.json` के modification time की तुलना `~/.local/state/ccam/mcp/<name>/seen` से करता है और फ़ाइल को केवल तभी parse करता है जब वह बदली हो।

### कमांड लाइन का रास्ता कैसे तय होता है

तीन कमांड लाइनें, `internal/app` से होकर तीन रास्ते:

<!-- motion:launch-route -->

शॉर्टकट का मिलान `argv[0]` पर होता है, `ccam <name>` का पहले आर्ग्युमेंट पर, और दोनों `syscall.Exec` पर ख़त्म होते हैं। केवल `ccam usage` जैसी ccam कमांड ही cobra तक पहुँचती है।

### साझा MCP सर्वर

डिफ़ॉल्ट `flag` strategy के साथ, ccam ग्रुप के MCP सर्वर Claude Code को एक आर्ग्युमेंट के रूप में देता है, आपके आर्ग्युमेंट से पहले:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` केवल तभी आर्ग्युमेंट लौटाता है जब ग्रुप की `mcp.json` मौजूद हो, और `Plan.InsertArgs` उसे `claude` के ठीक बाद रखता है। तब आपका अपना `--mcp-config` या `--strict-mcp-config` उसके ऊपर फिर भी लागू होता है। यह एक आर्ग्युमेंट है, `--mcp-config=<path>`, क्योंकि Claude Code का `--mcp-config` कई मान लेता है: दो आर्ग्युमेंट के रूप में लिखने पर वह बाद में आने वाले prompt को भी config फ़ाइल समझ लेता।

`config.toml` में `mcp_strategy = "sync"` होने पर, ccam इसके बजाय सर्वर अकाउंट की `.claude.json` में लिखता है, और अगर अकाउंट पहले से चल रहा हो तो यह छोड़ देता है।

### Remote Control

`ccam work rc` (या `ccam work remote-control`) Claude Code का Remote Control `work` के रूप में शुरू करता है। `claude remote-control` अपनी verb से पहले `--mcp-config` अस्वीकार करता है और उसके बाद कोई नहीं लेता, और उसके sessions MCP सर्वर केवल अकाउंट की `.claude.json` से पढ़ते हैं। इसलिए यहाँ flag strategy अलग तरह से काम करती है:

<!-- motion:launch-rc -->

1. `remoteControl` अकाउंट के बाद वाले पहले ऐसे आर्ग्युमेंट को, जो `-` से शुरू नहीं होता, Claude Code की verb मानता है। जब वह `remote-control` या `rc` हो, तो `mcp.FlagArgs` कुछ नहीं जोड़ता।
2. `mcp.Synced` ग्रुप की `mcp.json` के modification time की तुलना अकाउंट के `sync-stamp` से करता है। अगर `.claude.json` में पहले से ग्रुप के मौजूदा सर्वर हों, या ग्रुप में कोई सर्वर न हो, तो Remote Control बिना कुछ पूछे शुरू होता है।
3. वरना ccam उन सर्वरों के नाम बताता है जो sessions को नहीं मिलेंगे, और पूछता है:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` पर, `mcp.Sync` उन्हें वैसे ही कॉपी करता है जैसे `sync` strategy करती है: Claude Code के lock के तहत, बैकअप के बाद, अकाउंट के अपने सर्वर बनाए रखते हुए। `N` पर, या टर्मिनल के बिना, Remote Control उनके बिना शुरू होता है, और ccam `mcp_strategy = "sync"` की ओर इशारा करता है, जो हर बार शुरू करने से पहले उन्हें कॉपी करता है।

सर्वर कॉपी होने तक ccam हर बार Remote Control शुरू करने पर पूछता है, और कॉपी के बाद केवल तब जब ग्रुप की `mcp.json` फिर से बदले। बाद में सामान्य रूप से शुरू करने पर `.claude.json` की कॉपी `--mcp-config` पर जीतती है, इसलिए `ccam mcp add`, `remove` और `import` दूसरे अकाउंट्स के पास रखी कॉपियाँ अपडेट करते हैं ([साझा MCP सर्वर](/features/mcp))।

### 10 ms का बजट

Claude Code शुरू होने से पहले ccam अधिकतम 10 ms (p95) जोड़ता है। शुरू करने का रास्ता कभी cobra का कमांड ट्री नहीं बनाता, Bubble Tea शुरू नहीं करता और कोई नेटवर्क कनेक्शन नहीं खोलता, और अगर package के `init()` functions कुल मिलाकर 5 ms से ज़्यादा लें तो `TestPackageInitIsCheap` विफल होता है। `scripts/perf-check.sh`, `hyperfine` से `ccam --dry-run work`, `ccam work` और एक शॉर्टकट को मापता है और 10 ms से ऊपर विफल होता है, और `internal/launcher` में `BenchmarkLaunchPath` इन चरणों को loop में चलाता है। हर चरण का समय देखने के लिए `CCAM_TRACE=1` सेट करें:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## सुरक्षा और विशेष स्थितियाँ

- अज्ञात नाम या नंबर `no account "<name>"` के साथ रुकता है और `ccam list` की ओर इशारा करता है। अगर कोई अकाउंट ही न हो, तो ccam `ccam add <name>` से एक जोड़ने को कहता है।
- अगर कोई साझा पाथ सामान्य फ़ाइल बन गया हो या ग़लत जगह इशारा करता हो, तो ccam `stopped because of share link problems:` के साथ रुकता है, उनकी सूची दिखाता है, और `ccam doctor --fix` चलाने को कहता है। शुरू करने के रास्ते पर यह डेटा move नहीं करता।
- अगर इस अकाउंट के साथ बातचीत का इतिहास साझा करने वाला (उसी ग्रुप का) कोई दूसरा Claude Code उसी फ़ोल्डर में चल रहा हो, तो ccam चेतावनी देता है कि इतिहास मिल सकता है और `Run anyway? [y/N]` पूछता है। टर्मिनल के बिना यह चेतावनी दिखाकर आगे बढ़ता है।
- `--dry-run`, MCP sync, नए सर्वर की सूचना, इतिहास की जाँच, session marker, शीर्षक और आख़िरी उपयोग का समय छोड़ देता है, और exec नहीं करता। यह फिर भी साझा लिंक जाँचता है, और सामान्य शुरुआत की तरह गायब लिंक बनाता है।
- `ccam work rc`, जब `work` पहले से चल रहा हो, `.claude.json` को नहीं छूता, यह बताता है, और Remote Control को ग्रुप के सर्वरों के बिना शुरू करता है; उन्हें कॉपी करने के लिए session ख़त्म होने के बाद इसे फिर से चलाएँ। `--dry-run` के साथ यह संदेश प्रिंट करता है, लेकिन न पूछता है न कॉपी करता है।
- `mcp_strategy = "sync"` के साथ, विफल sync `ccam: MCP sync failed (launching anyway): …` प्रिंट करता है और Claude Code उन्हीं सर्वरों के साथ शुरू होता है जो उसके पास पहले से थे।

## संबंधित

- [शॉर्टकट कमांड](/features/shortcuts)
- [अकाउंट प्रबंधित करना](/features/accounts)
- [कमांड संदर्भ](/docs/commands)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [CLAUDE_CONFIG_DIR कैसे काम करता है](/guides/claude-config-dir)
