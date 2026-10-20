---
title: अकाउंट प्रबंधित करना
description: ccam add, ccam login, ccam list, ccam rename, ccam remove और ccam import से Claude Code अकाउंट जोड़ें, लॉगिन करें, सूची देखें, नाम बदलें, हटाएँ व import करें।
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

हर ccam अकाउंट एक अलग Claude Code config फ़ोल्डर है, जिसका अपना लॉगिन होता है। ये कमांड अकाउंट बनाती हैं, उनमें साइन इन करती हैं, उनकी सूची दिखाती हैं, उनका नाम बदलती हैं और उन्हें हटाती हैं, और आपके पास पहले से मौजूद फ़ोल्डर रजिस्टर करती हैं।

## उपयोग

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

जहाँ भी कोई कमांड `<name|number>` लेती है, वहाँ आप अकाउंट का नाम या उसका स्लॉट नंबर दे सकते हैं, यानी `ccam list` का `#` कॉलम।

टर्मिनल पर बिना विकल्पों के `ccam add` एक विज़र्ड खोलता है। यह नाम और अकाउंट के साइन इन करने का तरीक़ा पूछता है (अभी Claude लॉगिन, API key, या बाद में Claude लॉगिन), फिर शामिल होने के लिए ग्रुप दिखाता है जिनमें `default` चुना हुआ होता है (जब तक कोई ग्रुप न हो, तब वे आइटम जिन्हें ग्रुप `default` साझा करेगा), फिर `Create a shortcut command for work? [y/N]` पूछता है, और अंत में पुष्टि के लिए एक सारांश दिखाता है। `Esc` एक कदम पीछे ले जाता है। `--group`, `--no-group`, `--shortcut` या `--no-login` में से कोई भी विकल्प, या टर्मिनल की जगह pipe, विज़र्ड को छोड़ देता है; तब नाम देना ज़रूरी है और कोई शॉर्टकट तब तक नहीं बनता जब तक `--shortcut` उसका नाम न दे।

`--group` के बिना नया अकाउंट ग्रुप `default` में शामिल होता है। पहला अकाउंट इसे बनाता है, जो MCP सर्वर, सेटिंग्स, `CLAUDE.md` और एक्सटेंशन साझा करता है; [ग्रुप के साथ साझा करना](/features/share) देखें।

`ccam list` स्लॉट, नाम, शॉर्टकट, ईमेल, ग्रुप और आख़िरी उपयोग दिखाता है, और डिफ़ॉल्ट अकाउंट के आगे `*` लगाता है:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

`--json` के साथ हर अकाउंट एक ऑब्जेक्ट होता है, जिसमें `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (जो अकाउंट कभी शुरू नहीं हुआ, उसमें यह नहीं होता), `groups` और `share` होते हैं। `share` हर आइटम (`mcp`, `settings`, `memory`, `extensions`, `history`) को उस ग्रुप से जोड़ता है जिससे वह आता है, या `"local"` से।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `ccam add --no-login` | लॉगिन किए बिना अकाउंट बनाएँ; बाद में `ccam login <name>` से लॉगिन करें | बंद |
| `ccam add --group <group>` | `default` के बजाय इस ग्रुप में शामिल हों। कई बार दिया जा सकता है; `default` न हो तो बनाया जाता है | `default` |
| `ccam add --no-group` | किसी ग्रुप में शामिल न हों: अकाउंट कुछ साझा नहीं करता | बंद |
| `ccam add --shortcut <command>` | `~/.local/bin` में यह शॉर्टकट कमांड भी बनाएँ | कोई नहीं |
| `ccam add --api --key-from <source>` | ऐसा अकाउंट जो लॉगिन के बजाय API key इस्तेमाल करता है; [API key अकाउंट](/features/api) देखें | बंद |
| `ccam list --json` | अकाउंट्स को JSON के रूप में प्रिंट करें | बंद |
| `ccam remove -y`, `--yes` | बिना पूछे हटाएँ | बंद |
| `ccam import --name <name>` | आपके दिए फ़ोल्डर के लिए अकाउंट का नाम | फ़ोल्डर के नाम से बना |
| `ccam import -y`, `--yes` | मिले हर फ़ोल्डर को बिना पूछे import करें | बंद |
| `--dry-run` | `ccam list` और `ccam completion` इसे मानते हैं, क्योंकि वे कुछ नहीं बदलते। `add`, `login`, `logout`, `use`, `rename`, `remove` और `import` इसे `` --dry-run is not supported by `ccam add`; nothing was changed `` (कमांड के नाम के साथ) कहकर अस्वीकार करते हैं | बंद |

`ccam login`, `ccam logout`, `ccam use` और `ccam rename` के अपने कोई विकल्प नहीं हैं। `ccam completion` `bash`, `zsh` या `powershell` लेता है; कुछ और देने पर यह `only bash, zsh and powershell are supported` के साथ विफल होता है।

अकाउंट के नाम 1 से 32 अक्षरों के होते हैं, जिनमें `a-z`, `0-9`, `-` और `_` आते हैं, और वे किसी अक्षर या अंक से शुरू होते हैं। नाम केवल अंकों का नहीं हो सकता, क्योंकि `ccam 2` का मतलब स्लॉट 2 है। नाम कोई ccam कमांड भी नहीं हो सकता (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), ताकि `ccam <name>` कभी अस्पष्ट न हो। यह सूची `internal/config` में `config.Reserved` है।

## यह कैसे काम करता है

`ccam add work --shortcut ccwork` ये चरण `addProfile` (`internal/cli/profiles.go`) में चलाता है:

<!-- motion:accounts-add -->

1. कुछ भी बनाने से पहले `config.ValidateName` नाम जाँचता है, और `shim.Check` शॉर्टकट जाँचता है।
2. `profile.Add` मोड `0700` के साथ `~/.local/share/ccam/profiles/work/` बनाता है और अकाउंट को सबसे छोटा ख़ाली स्लॉट नंबर देता है (`Config.NextSlot`), इसलिए हटाए गए अकाउंट का नंबर दोबारा इस्तेमाल होता है। पहला अकाउंट डिफ़ॉल्ट बनता है।
3. अकाउंट अपने ग्रुप्स में शामिल होता है (`default`, जब तक `--group` या `--no-group` कुछ और न कहे; पहला अकाउंट `default` बनाता है), और `applyShare` नए, ख़ाली फ़ोल्डर को `share.Build` और फिर `Apply` से उनसे लिंक करता है, उदाहरण के लिए `settings.json` को `groups/default/settings.json` से। अगर यह विफल हो, तो नया फ़ोल्डर हटा दिया जाता है, क्योंकि बचा हुआ फ़ोल्डर हर अगली कोशिश को `folder already exists` के साथ रोक देता।
4. अकाउंट `fsutil.AtomicWrite` के ज़रिए `~/.config/ccam/config.toml` में सेव होता है: एक अस्थायी फ़ाइल, `fsync`, फिर `rename`, ताकि crash होने पर भी आधा config कभी न बचे।
5. `runLogin` अकाउंट फ़ोल्डर पर सेट `CLAUDE_CONFIG_DIR` के साथ, और अकाउंट के पिछले ईमेल वाले `--email` के साथ `claude auth login` शुरू करता है, ताकि लॉगिन पेज पहले से भरा हो। अगर लॉगिन विफल हो, तो अकाउंट बना रहता है, और ccam त्रुटि प्रिंट करके बताता है कि बाद में `ccam login work` चलाएँ।
6. `profile.ReadIdentity`, `profiles/work/.claude.json` से `oauthAccount` निकालता है: ईमेल, अकाउंट और संगठन के id, संगठन का नाम और प्लान (`claude_max`, `claude_team` …)। वह फ़ाइल कई मेगाबाइट की हो सकती है, इसलिए `TopLevelValue` पूरी फ़ाइल parse करने के बजाय केवल उस एक key को खोजता है। ccam इन्हें `config.toml` में `[profile.login]` में रखता है।
7. `shim.Install`, `~/.local/bin/ccwork` को ccam बाइनरी से लिंक करता है, जो शुरू होते समय `config.toml` में `argv[0]` खोजती है।

लॉगिन में कई मिनट लग सकते हैं, और इस बीच कोई दूसरा ccam config सेव कर सकता है। इसलिए इंतज़ार करने वाली कमांड `env.update` के ज़रिए सेव करती हैं, जो `config.toml` दोबारा लोड करता है, केवल वह एक बदलाव (यहाँ, ईमेल) लागू करता है और सेव करता है। इसी कारण `ccam remove` भी अपने सवाल के बाद दोबारा लोड करता है।

`ccam login` किसी मौजूदा अकाउंट के लिए चरण 5 और 6 चलाता है। Claude Code ईमेल को लॉगिन पेज के लिए केवल एक संकेत की तरह लेता है, इसलिए लॉगिन उस Claude अकाउंट के बजाय किसी दूसरे Claude अकाउंट में हो सकता है जो ccam अकाउंट पहले था; ccam नए लॉगिन की तुलना `[profile.login]` से करता है और यह बताता है (`ccam: this login is other@company.com, but work was dev@company.com before.`)। यह नया लॉगिन रखता है और बताता है कि इसे कैसे पूर्ववत करें: `ccam logout work`, फिर `ccam login work`, और सही ईमेल से लॉगिन करें। जब कोई दूसरा ccam अकाउंट उसी Claude अकाउंट में लॉगिन हो, तो ccam चेतावनी देता है कि दोनों अपनी उपयोग सीमाएँ साझा करते हैं। `ccam login --expired` एक के बाद एक हर उस अकाउंट में लॉगिन करता है जिसमें कभी लॉगिन नहीं हुआ या जिसका लॉगिन (उसका refresh token) समाप्त हो गया है।

`ccam logout` अकाउंट में `claude auth logout` चलाता है और सेव किया गया ईमेल हटा देता है; लॉगिन का रिकॉर्ड बना रहता है, इसलिए अगला लॉगिन उससे जाँचा जाता है। `ccam use` केवल `config.toml` में `default` को फिर से लिखता है।

जब अकाउंट एक से ज़्यादा संगठनों के हों, तो `ccam list` एक `Organization` कॉलम जोड़ता है (personal प्लान `personal` दिखाते हैं), और `ccam list --json` में हर अकाउंट के लिए `plan` और `org` होते हैं।

### रिमोट मशीन पर साइन इन करना

SSH पर `claude auth login` का ब्राउज़र वाला हिस्सा टूट सकता है। Claude Code सर्वर पर `http://localhost:<port>/callback` वाले callback सर्वर पर साइन-इन का इंतज़ार करता है, लेकिन साइन-इन पूरा करने वाला ब्राउज़र आपके कंप्यूटर पर चलता है, जहाँ वह पता कहीं नहीं ले जाता और ब्राउज़र `ERR_CONNECTION_REFUSED` जैसी त्रुटि दिखाता है। पते में एक बार इस्तेमाल होने वाला कोड फिर भी मौजूद रहता है, इसलिए ccam उसे पहुँचा सकता है:

<!-- motion:accounts-relay -->

`runLogin` आपका टर्मिनल Claude Code को नहीं सौंपता। यह Claude Code के standard input को एक pipe से जोड़ता है और `relayLogin` आपकी टाइप की गई हर पंक्ति उसमें कॉपी करता है, इसलिए Claude Code के `Paste code here` prompt पर paste किया गया कोड फिर भी पहुँचता है। जब कोई पंक्ति callback पता हो, तो `callbackURL` उसे केवल तभी स्वीकार करता है जब host `localhost`, `127.0.0.1` या `[::1]` हो, port दिया गया हो, path `/callback` हो और `code` तथा `state` दोनों मौजूद हों; `http://` छोड़ा जा सकता है। फिर `deliverCallback` सर्वर पर 10 सेकंड की सीमा के साथ और redirects का पालन किए बिना उस पते पर request भेजता है, ठीक वैसे जैसे ब्राउज़र भेजता। Request इंतज़ार कर रहे `claude auth login` के callback सर्वर तक पहुँचती है (आरेख में डैश वाली रेखा)। Claude Code कोड पाता है और लॉगिन सेव करता है, Linux पर अकाउंट फ़ोल्डर की `.credentials.json` में। Paste की गई कोई पंक्ति ccam से कभी किसी दूसरे host से संपर्क नहीं करवा सकती।

लॉगिन शुरू होने से पहले ccam यह प्रिंट करता है:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

सफल होने पर यह `→ Passed the sign-in to Claude Code on this machine.` प्रिंट करता है।

जब लॉगिन टर्मिनल पर input का इंतज़ार करता है, तब ccam टर्मिनल का echo और line mode बंद कर देता है और आपके टाइप किए हर अक्षर के लिए एक `*` दिखाता है, ताकि paste किया गया code या callback पता स्क्रीन पर न दिखे। `readLine` line editing ख़ुद करता है: `Enter` पंक्ति भेजता है, `Backspace` एक अक्षर मिटाता है, और ख़ाली पंक्ति पर `Ctrl+D` input ख़त्म करता है। लॉगिन ख़त्म होने पर टर्मिनल पहले जैसी स्थिति में आ जाता है, और `Ctrl+C`, `SIGTERM` या `SIGHUP` पर भी, जो फिर भी ccam को बंद करते हैं। Standard input और standard error पर टर्मिनल न होने पर, पंक्तियाँ जैसे आती हैं वैसे ही कॉपी होती हैं।

### नाम बदलना, हटाना और import करना

`ccam rename` और `ccam remove` अकाउंट के फ़ोल्डर पर, और उस पर काम करते हैं जो ccam अकाउंट के नाम से दूसरी जगहों पर रखता है। `nameState` उन पाथ की सूची देता है: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` और उपयोग की cache फ़ाइलें `~/.local/state/ccam/usage/<name>.json` और `<name>.statusline.json`। `ccam remove work` इन्हें इस क्रम में संभालता है:

<!-- motion:accounts-remove -->

`ccam remove` पहले अकाउंट चल रहा हो तो मना कर देता है, फिर पूछता है। अगर अकाउंट एक्सटेंशन साझा करता है, तो `RewritePluginPaths` ग्रुप की `installed_plugins.json` और `known_marketplaces.json` को फिर से लिखता है, क्योंकि Claude Code plugins के absolute पाथ सेव करता है जो इस अकाउंट के फ़ोल्डर से होकर जा सकते हैं। `profile.Trash` फ़ोल्डर को `~/.local/share/ccam/trash/<name>-<time>/` में move करता है, अकाउंट `config.toml` से हट जाता है, `nameState` पाथ मिटा दिए जाते हैं ताकि बाद में उसी नाम का अकाउंट उन्हें विरासत में न पाए, और शॉर्टकट लिंक हटा दिया जाता है। वापस पाने के लिए फ़ोल्डर को trash से वापस move करें।

`ccam rename` फ़ोल्डर का नाम बदलता है, उसके `plugins/` में plugin पाथ को पुराने फ़ोल्डर से नए फ़ोल्डर पर फिर से लिखता है, डिफ़ॉल्ट को साथ ले जाता है, और `nameState` पाथ को नए नाम पर move करता है। शॉर्टकट कमांड का नाम वही रहता है। अगर plugin पाथ फिर से लिखना विफल हो, तो फ़ोल्डर का नाम वापस बदल दिया जाता है, क्योंकि config सेव नहीं होता।

`ccam import` किसी फ़ोल्डर को जहाँ वह है वहीं, `config.toml` में `dir` के रूप में रजिस्टर करता है। यह उसे कभी move या copy नहीं करता: कॉपी किया गया refresh token कॉपियों को एक-दूसरे से लॉग आउट करवा देता। फ़ोल्डर आर्ग्युमेंट के बिना यह बिना रजिस्टर वाले `~/.claude` और `~/.claude-*` फ़ोल्डर खोजता है और हर एक के बारे में पूछता है; नाम फ़ोल्डर से आते हैं (`~/.claude` का `main`, `~/.claude-work` का `work`, `~/.claude-2` का `acct2`), और जिस फ़ोल्डर से नाम नहीं बन सकता, उसे इस्तेमाल करने वाली `--name` कमांड के साथ दिखाया जाता है। `CLAUDE_CONFIG_DIR` के बिना Claude Code `.claude.json` को आपके home फ़ोल्डर में `~/.claude` के बगल में रखता है, इसलिए जब आप `~/.claude` import करते हैं और उसकी अपनी कोई `.claude.json` नहीं होती, तो ccam `~/.claude/.claude.json` से `~/.claude.json` का लिंक बनाता है। तब सादा `claude` और `ccam main` एक ही अकाउंट देखते हैं। Import किए गए अकाउंट तब तक किसी ग्रुप में नहीं होते जब तक आप उन्हें किसी ग्रुप में न जोड़ें (`ccam group join <group> <name>`), और कोई शॉर्टकट नहीं बनता।

`ccam completion bash|zsh` एक completion स्क्रिप्ट प्रिंट करता है जो अकाउंट के नाम भी पूरे करती है।

## सुरक्षा और विशेष स्थितियाँ

- अकाउंट चलते समय `ccam rename` और `ccam remove` मना कर देते हैं (`refuseRunning`), `work is running; quit its Claude Code session first` के साथ, क्योंकि session फ़ोल्डर में लिखता रहता है।
- Import किए गए अकाउंट पर `ccam remove` केवल उसका रजिस्ट्रेशन हटाता है: `✓ main unregistered (folder ... left in place)`। आपका `~/.claude` trash में नहीं जाता।
- `ccam add` ऐसा नाम अस्वीकार करता है जिसका फ़ोल्डर `profiles/` में पहले से मौजूद है, और `ccam import` की ओर इशारा करता है।
- `ccam import` पहले से रजिस्टर फ़ोल्डर, ऐसा पाथ जो फ़ोल्डर नहीं है, और पहले से लिया गया नाम अस्वीकार करता है।
- `ccam logout` और `ccam login` यह नहीं जाँचते कि अकाउंट चल रहा है या नहीं।
- अगर relay Claude Code तक नहीं पहुँच पाता, उदाहरण के लिए लॉगिन ख़त्म होने के बाद, तो यह कारण के साथ `ccam: could not pass the sign-in to Claude Code:` प्रिंट करता है। HTTP त्रुटि का आमतौर पर मतलब है कि कोड एक्सपायर हो गया; लॉगिन फिर से शुरू करें।

## संबंधित

- [Claude Code शुरू करना](/features/launch)
- [शॉर्टकट कमांड](/features/shortcuts)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [कई Claude Code अकाउंट](/guides/multiple-claude-code-accounts)
