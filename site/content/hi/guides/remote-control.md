---
title: ccam के साथ Remote Control का उपयोग
description: ccam <name> rc से किसी भी अकाउंट से Claude Code Remote Control शुरू करें; लागू विकल्प, साझा MCP सर्वर कॉपी करने का प्रस्ताव और समस्याओं के हल देखें।
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 4d46f49c5dea
---

`ccam work rc` Claude Code का Remote Control `work` अकाउंट के रूप में शुरू करता है, ताकि आप इस कंप्यूटर पर चल रहे sessions को claude.ai/code या Claude मोबाइल ऐप से नियंत्रित कर सकें। यह `claude rc` की तरह ही काम करता है, बस एक अंतर के साथ: अगर अकाउंट MCP सर्वर साझा करता है, तो ccam एक बार पूछता है कि क्या उन्हें अकाउंट की `.claude.json` में कॉपी किया जाए, क्योंकि वही फ़ाइल अकेली जगह है जहाँ से Remote Control के sessions सर्वर पढ़ते हैं।

## इसे शुरू करें

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[शॉर्टकट](/features/shortcuts) भी इसी तरह काम करता है: `ccwork rc` का मतलब `ccam work rc` है। कमांड उस फ़ोल्डर में चलाएँ जिसमें आप काम करना चाहते हैं, फिर claude.ai/code या Claude ऐप खोलें, उसी Claude अकाउंट से साइन इन करके जो `work` का है।

## विकल्प

यहाँ दो तरह के विकल्प मिलते हैं: ccam के अपने, जो अकाउंट के नाम से पहले आते हैं या ccam की सेटिंग्स में रहते हैं, और Remote Control के, जो `rc` के बाद आते हैं और बिना बदले Claude Code तक पहुँचते हैं।

### ccam के विकल्प और सेटिंग्स

| विकल्प या सेटिंग | यह क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--dry-run` (अकाउंट से पहले) | environment और कमांड प्रिंट करता है, और MCP सूचना भी, अगर कोई हो। न पूछता है, न कॉपी करता है | बंद |
| `~/.config/ccam/config.toml` में सबसे ऊपर `mcp_strategy = "sync"` | हर बार शुरू करने से पहले, चाहे Remote Control हो या नहीं, ग्रुप के सर्वर `.claude.json` में कॉपी करता है, इसलिए कभी कुछ नहीं पूछा जाता | `"flag"` |
| `ccam group leave <group> work` | अकाउंट को उस ग्रुप से निकालता है जो उसके साथ MCP साझा करता है। ग्रुप के सर्वर उसकी `.claude.json` में उसके अपने सर्वरों के रूप में कॉपी होते हैं, इसलिए Remote Control के पास वे होते हैं और कुछ नहीं पूछा जाता; `--fresh` के साथ यह केवल अकाउंट के अपने सर्वरों के साथ शुरू होता है | `default` में |

जब अकाउंट MCP सर्वर साझा करता है और उसकी `.claude.json` में ग्रुप की मौजूदा सूची नहीं होती, या कोई कॉपी किया गया सर्वर उससे ग़ायब हो, तो ccam पूछता है `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`। `y` उन्हें कॉपी करता है; `N` या Enter उनके बिना शुरू करता है।

### Remote Control के विकल्प

`rc` के बाद की हर चीज़ Claude Code की है। ये Claude Code 2.1.304 में `claude remote-control` के विकल्प हैं; `ccam work rc --help` आपके वर्शन की सूची प्रिंट करता है।

| विकल्प | यह क्या करता है |
| --- | --- |
| `--name <name>` | claude.ai/code में दिखने वाला session का नाम |
| `--remote-control-session-name-prefix <prefix>` | अपने-आप बनने वाले session नामों का prefix (डिफ़ॉल्ट: host का नाम) |
| `-c`, `--continue` | उस session से फिर जुड़ता है जिसे Remote Control ने इस फ़ोल्डर के लिए आख़िरी बार दर्ज किया था, लगभग 4 घंटे के भीतर |
| `--session-id <id>` | किसी ख़ास session से फिर जुड़ता है |
| `--permission-mode <mode>` | नए sessions के लिए permission mode: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` या `plan` |
| `--spawn <mode>` | `same-dir` (डिफ़ॉल्ट), `worktree` (हर session के लिए एक git worktree) या `session` (एक session, जो उसके ख़त्म होने पर बाहर निकल जाता है) |
| `--capacity <N>` | `same-dir` और `worktree` mode में एक साथ चलने वाले sessions की अधिकतम संख्या (डिफ़ॉल्ट 32) |
| `--[no-]create-session-in-dir` | शुरू होते समय मौजूदा फ़ोल्डर में पहला session बनाता है (डिफ़ॉल्ट चालू) |
| `--[no-]chrome` | नए sessions के लिए Claude in Chrome (डिफ़ॉल्ट बंद) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Debug और विस्तृत आउटपुट |

ccam Claude Code की कमांड वैसे ही पढ़ता है जैसे Claude Code पढ़ता है, विकल्पों को उनके मानों के साथ छोड़ते हुए, इसलिए `ccam work rc --name laptop` और `ccam work --model opus rc` दोनों Remote Control शुरू करते हैं। जिस विकल्प को ccam नहीं जानता, उसे बिना मान वाला विकल्प माना जाता है; सुरक्षित रहने के लिए `rc` को अकाउंट के नाम के ठीक बाद रखें।

## ccam क्या करता है

शुरुआत सामान्य रास्ते से होती है ([Claude Code शुरू करना](/features/launch)), बस MCP वाले चरण में एक बदलाव के साथ:

1. ccam अकाउंट ढूँढता है, उसके share links जाँचता है और `claude` का पता लगाता है, जैसे किसी भी शुरुआत में।
2. वह देखता है कि कमांड `rc` या `remote-control` है और `--mcp-config` नहीं जोड़ता।
3. अगर अकाउंट MCP सर्वर साझा करता है, तो वह ग्रुप की `mcp.json` की तुलना `~/.local/state/ccam/mcp/work/` में मौजूद `sync-stamp` से करता है, और जाँचता है कि हर ग्रुप सर्वर अभी भी `.claude.json` में है। जब `.claude.json` में पहले से मौजूदा सर्वर हों, या ग्रुप में कोई सर्वर न हो, तो वह सीधे चरण 6 पर जाता है।
4. वरना वह उन सर्वरों के नाम बताता है जो sessions को नहीं मिलेंगे:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. `y` पर वह `.claude.json` का बैकअप सहेजता है, Claude Code का अपना lock लेता है, ग्रुप के सर्वर `mcpServers` में जोड़ता है और अकाउंट के बाकी सर्वर बनाए रखता है; अगर अकाउंट का अपना कोई सर्वर किसी ग्रुप सर्वर के नाम वाला है, तो ccam अकाउंट का सर्वर रखता है और यह बताता है। `N` पर, टर्मिनल के बिना, या जब अकाउंट चल रहा हो, वह कुछ नहीं बदलता और यही बताता है:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. वह session दर्ज करता है, `CLAUDE_CONFIG_DIR` को अकाउंट के फ़ोल्डर पर सेट करता है और `exec` के ज़रिए ख़ुद को `claude rc` से बदल देता है।

इसके बाद Remote Control एक सर्वर के रूप में चलता रहता है। आपके लिए शुरू किया गया हर session उसी process का child होता है, इसलिए वह `work` के रूप में चलता है: `work` के लॉगिन, सेटिंग्स और बातचीत के इतिहास के साथ, और `work` की उपयोग सीमा में गिना जाता है। जब तक ccam का शुरू किया Remote Control सर्वर चल रहा है, ccam अकाउंट को चालू मानता है। सादे `claude remote-control` से शुरू हुआ सर्वर केवल तब तक गिना जाता है जब तक उसके sessions हों; जिसका अभी कोई session नहीं है, उसे ccam नहीं देख सकता।

Dry run अंतर दिखाता है। सामान्य शुरुआत ग्रुप के सर्वरों को एक आर्ग्युमेंट के रूप में साथ ले जाती है, Remote Control नहीं:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## यह इस तरह क्यों काम करता है

**Remote Control `--mcp-config` नहीं लेता।** सामान्य शुरुआत में ccam ग्रुप के सर्वर Claude Code को `--mcp-config=<path>` के रूप में देता है और Claude Code की फ़ाइलों में कभी नहीं लिखता। `claude remote-control` यह विकल्प अपने कमांड नाम से पहले अस्वीकार करता है और उसके बाद कोई नहीं लेता, इसलिए ccam 0.6.0 से पहले `ccam work rc` हर उस अकाउंट के लिए विफल होता था जो MCP सर्वर साझा करता था। अब ccam यह विकल्प छोड़ देता है।

**इसके sessions सर्वर केवल `.claude.json` से पढ़ते हैं।** Remote Control जो sessions शुरू करता है, वे कमांड लाइन से विकल्प नहीं लेते। उन्हें user MCP सर्वर केवल एक जगह मिलते हैं: अकाउंट की `.claude.json`। उसी फ़ाइल में अकाउंट की लॉगिन जानकारी और हर फ़ोल्डर की स्थिति भी रहती है, इसलिए ccam उसे `settings.json` की तरह किसी साझा कॉपी से link नहीं कर सकता। sessions तक साझा सूची पहुँचाने का एकमात्र तरीका सर्वरों को उसमें कॉपी करना है।

**ccam लिखने से पहले पूछता है।** डिफ़ॉल्ट flag strategy में ccam शुरुआत के समय `.claude.json` को नहीं छूता। Remote Control अकेली शुरुआत है जिसमें उसे वहाँ लिखना पड़ता है, इसलिए वह पहले पूछता है, बैकअप रखता है, और `mcpServers` के अलावा कुछ नहीं छूता। (शुरुआत के बाहर, MCP sharing छोड़ना या उसमें फिर से जुड़ना और `ccam mcp import` भी वहाँ लिखते हैं, जैसा उनकी योजनाएँ बताती हैं।)

**अकाउंट चलते समय यह नहीं लिखता।** Claude Code session के दौरान `.claude.json` दोबारा लिखता है, और बाहर से किया गया बदलाव खो जाएगा या उस पर लिख दिया जाएगा। इसीलिए जब अकाउंट का कोई भी session चल रहा हो, किसी दूसरे Remote Control सर्वर और ccam के बिना शुरू हुए Claude Code समेत, तो ccam कॉपी छोड़ देता है।

**हर बदलाव पर यह एक बार पूछता है।** कॉपी के बाद `sync-stamp` ग्रुप की `mcp.json` का modification time रखता है। ccam दोबारा केवल तब पूछता है जब वह फ़ाइल बदलती है, जैसे `ccam mcp add` के बाद, या जब कोई कॉपी किया गया सर्वर `.claude.json` से ग़ायब हो। अगर आप `N` जवाब देते हैं, तो अभी कोई कॉपी नहीं होती, इसलिए यह हर बार Remote Control शुरू करने पर पूछता है।

**`ccam mcp` कॉपियों को अद्यतन रखता है।** नाम टकराने पर `.claude.json` की entry `--mcp-config` पर जीतती है, इसलिए पुरानी कॉपी सामान्य शुरुआत में भी नई परिभाषा को छिपा देगी। इसलिए `ccam mcp add`, `add-json`, `remove` और `import`, और `ccam migrate`, हर उस अकाउंट की कॉपियाँ अपडेट करते हैं जो नहीं चल रहा, और कोई कॉपी बदलने पर यह बताते हैं:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**टर्मिनल के बिना जवाब "नहीं" है।** Remote Control अक्सर किसी script, service या `nohup` से चलता है, जहाँ जवाब देने वाला कोई नहीं होता। तब ccam इंतज़ार करने के बजाय सर्वरों के बिना शुरू करता है।

## समस्याएँ और समाधान

### Remote Control sessions में साझा MCP सर्वर नहीं हैं

सवाल का जवाब `N` दिया गया, टर्मिनल नहीं था, या अकाउंट चल रहा था। Remote Control रोकें, अकाउंट के दूसरे sessions बंद करें, और टर्मिनल में `ccam work rc` फिर से चलाएँ; जवाब `y` दें। पूछा जाना बंद करने के लिए `~/.config/ccam/config.toml` में `mcp_strategy = "sync"` सेट करें। किसी session में पहले से मौजूद सर्वर उसके चलते रहने तक नहीं बदलते, इसलिए कॉपी के बाद नया session शुरू करें।

### "Account work is running, so its .claude.json is left alone"

`work` का कोई session अभी खुला है: सामान्य `ccam work`, कोई शॉर्टकट, या कोई दूसरा Remote Control सर्वर। उससे बाहर निकलें और `ccam work rc` फिर से चलाएँ। किसी दूसरे अकाउंट के चलाए sessions नहीं गिने जाते।

### ccam हर बार पूछता है

आपने अभी तक सर्वर कॉपी नहीं किए हैं। एक बार `y` जवाब दें, या `mcp_strategy = "sync"` सेट करें। अगर आपको Remote Control में सर्वर कभी नहीं चाहिए, तो अकाउंट के साथ MCP साझा करना उन्हें रखे बिना बंद करें: `ccam group leave default work --fresh`, या MCP को उसके अपने ग्रुप में रखें जिसमें यह अकाउंट न हो। सादा `ccam group leave` ग्रुप के सर्वरों को अकाउंट की `.claude.json` में कॉपी करता है, इसलिए Remote Control के पास वे होते हैं।

### `ccam work rc` तुरंत `--mcp-config` से जुड़ी त्रुटि के साथ विफल होता है

आपका ccam 0.6.0 से पुराना है और Remote Control को `--mcp-config` देता है। `ccam update` चलाएँ। नए ccam के साथ, यही त्रुटि आने का मतलब है कि कोई अज्ञात विकल्प, मान के साथ, `rc` से पहले आया, और ccam ने ग़लती से उस मान को कमांड मान लिया; `rc` को अकाउंट के नाम के ठीक बाद रखें।

### `ccam mcp` से किया गया बदलाव Remote Control तक नहीं पहुँचता

`ccam mcp add` और `remove` केवल उन अकाउंट्स में कॉपियाँ अपडेट करते हैं जो नहीं चल रहे, और यह बताते भी हैं: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control रोकें और फिर से शुरू करें; ccam बदली हुई `mcp.json` देखकर फिर पूछेगा।

### ग्रुप के सर्वर के बजाय मेरा अपना सर्वर इस्तेमाल होता है

अगर अकाउंट का अपना कोई सर्वर ग्रुप के किसी सर्वर के नाम वाला है, तो कॉपी उसे उसकी जगह रहने देती है और यह बताती है: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`। ग्रुप वाला इस्तेमाल करने के लिए, अकाउंट वाला हटाएँ (`claude mcp remove -s user <name>`); अपना वर्शन हर अकाउंट को देने के लिए, उसे `ccam mcp import work <server> --force` से ग्रुप में ले जाएँ। `.claude.json` के बैकअप `~/.local/state/ccam/mcp/work/backups/` में हैं, जो सबसे नए पाँच बैकअप रखता है।

### मैं कॉपी किए गए सर्वर हटाना चाहता हूँ

`ccam mcp remove <server>` सर्वर को ग्रुप से और हर कॉपी से हटा देता है। कॉपी से पहले की स्थिति पर लौटने के लिए, जब अकाउंट का कोई session न चल रहा हो, `~/.local/state/ccam/mcp/work/backups/` से सबसे नया बैकअप वापस लाएँ।

### Remote Control पूछता है कि क्या फ़ोल्डर पर भरोसा है, या टर्मिनल के बिना बंद हो जाता है

फ़ोल्डर पर भरोसा हर अकाउंट के लिए `.claude.json` में अलग से रखा जाता है, इसलिए जिस फ़ोल्डर पर आपने `personal` के रूप में भरोसा किया, वह `work` के लिए नया है। उस फ़ोल्डर में एक बार `ccam work` चलाएँ और स्वीकार करें, फिर वहाँ Remote Control शुरू करें।

### "You must be logged in with a Claude account that has a subscription"

Remote Control को Claude subscription वाले लॉगिन की ज़रूरत है। `ccam list` में अकाउंट का ईमेल जाँचें, और ज़रूरत हो तो `ccam login work` से फिर से साइन इन करें।

### फ़ोन पर session नहीं दिखता

Claude ऐप या claude.ai/code किसी दूसरे Claude अकाउंट से साइन इन है। उसकी तुलना उस ईमेल से करें जो `ccam list` `work` के लिए दिखाता है।

### ccam चेतावनी देता है कि इस फ़ोल्डर में कोई दूसरा Claude Code session बातचीत का इतिहास साझा करता है

उसी इतिहास ग्रुप का कोई दूसरा अकाउंट उसी फ़ोल्डर में चल रहा है। दोनों अकाउंट्स से एक ही बातचीत आगे बढ़ाने वाले sessions इतिहास को गड़बड़ा सकते हैं। `y` तभी जवाब दें जब वे अलग-अलग बातचीत पर काम करेंगे, वरना Remote Control किसी दूसरे फ़ोल्डर में शुरू करें।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [साझा MCP सर्वर](/features/mcp)
- [ग्रुप के साथ साझा करना](/features/share)
- [अकाउंट प्रबंधित करना](/features/accounts)
