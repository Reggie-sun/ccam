---
title: खाती व्यवस्थापित करणे
description: ccam add, ccam login, ccam list, ccam rename, ccam remove आणि ccam import ने Claude Code खाती जोडा, लॉग इन करा, यादी पाहा, नाव बदला, काढा आणि इम्पोर्ट करा.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

प्रत्येक ccam खाते हा स्वतःचे लॉगिन असलेला वेगळा Claude Code config फोल्डर आहे. या कमांड्स खाती तयार करतात, त्यांत साइन इन करतात, त्यांची यादी दाखवतात, नाव बदलतात आणि ती काढतात, तसेच तुमच्याकडे आधीच असलेले फोल्डर नोंदवतात.

## वापर

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
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

जिथे एखादी कमांड `<name|number>` घेते, तिथे तुम्ही खात्याचे नाव किंवा त्याचा slot क्रमांक देऊ शकता, म्हणजे `ccam list` मधील `#` कॉलम.

कोणत्याही पर्यायाशिवाय टर्मिनलवर `ccam add` चालवल्यास विझार्ड उघडतो. तो नाव विचारतो, मग शेअर केलेल्या आयटम्सची चेकलिस्ट दाखवतो (संभाषण इतिहास सोडून सर्व काही निवडलेले असते), मग `Create a shortcut command for work? [y/N]` विचारतो आणि `y` दिल्यास `Shortcut name [ccwork]:` विचारतो. `--share`, `--shortcut` किंवा `--no-login` यांपैकी कोणताही पर्याय, किंवा टर्मिनलऐवजी pipe, विझार्ड वगळतो; तेव्हा नाव देणे आवश्यक असते आणि `--shortcut` ने नाव दिल्याशिवाय शॉर्टकट तयार होत नाही.

`ccam list` slot, नाव, शॉर्टकट, ईमेल, शेअर केलेले आयटम आणि शेवटचा वापर दाखवते, आणि डीफॉल्ट खात्याच्या आधी `*` असतो:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` सह प्रत्येक खाते एक object असते, ज्यात `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (खाते कधीच सुरू केले नसेल तर वगळले जाते) आणि `share` असतात. `share` प्रत्येक आयटम (`mcp`, `settings`, `memory`, `extensions`, `history`) त्याच्या गटाच्या नावाशी किंवा `"local"` शी जोडतो.

## पर्याय

| पर्याय | काय करतो | डीफॉल्ट |
| --- | --- | --- |
| `ccam add --no-login` | लॉग इन न करता खाते तयार करा; नंतर `ccam login <name>` ने लॉग इन करा | बंद |
| `ccam add --share <item=group>` | एखादा आयटम गटासोबत शेअर करा, किंवा तो खाजगी ठेवण्यासाठी `item=local`. अनेकदा देता येतो | `mcp`, `settings`, `memory`, `extensions` हे `default` मध्ये; `history` local |
| `ccam add --shortcut <command>` | `~/.local/bin` मध्ये ही शॉर्टकट कमांडही तयार करा | काहीही नाही |
| `ccam list --json` | खाती JSON म्हणून छापा | बंद |
| `ccam remove -y`, `--yes` | न विचारता काढा | बंद |
| `ccam import --name <name>` | तुम्ही दिलेल्या फोल्डरसाठी खात्याचे नाव | फोल्डरच्या नावावरून |
| `ccam import -y`, `--yes` | सापडलेला प्रत्येक फोल्डर न विचारता इम्पोर्ट करा | बंद |
| `--dry-run` | काहीही न बदलणाऱ्या `ccam list` आणि `ccam completion` कडून पाळला जातो. `add`, `login`, `logout`, `use`, `rename`, `remove` आणि `import` तो `` --dry-run is not supported by `ccam add`; nothing was changed `` (कमांडच्या नावासह) देऊन नाकारतात | बंद |

`ccam login`, `ccam logout`, `ccam use` आणि `ccam rename` यांना स्वतःचे पर्याय नाहीत. `ccam completion` `bash` किंवा `zsh` घेते; इतर काहीही दिल्यास ते `only bash and zsh are supported` देऊन अयशस्वी होते.

खात्यांची नावे `a-z`, `0-9`, `-` आणि `_` यांपासून बनलेली 1 ते 32 अक्षरांची असतात, आणि अक्षराने किंवा अंकाने सुरू होतात. नाव फक्त अंकांचे असू शकत नाही, कारण `ccam 2` म्हणजे slot 2. नाव ccam ची कमांडही असू शकत नाही (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), त्यामुळे `ccam <name>` कधीच संदिग्ध होत नाही. ही यादी `internal/config` मधील `config.Reserved` आहे.

## हे कसे काम करते

`ccam add work --shortcut ccwork` हे टप्पे `addProfile` (`internal/cli/profiles.go`) मध्ये चालवते:

<!-- motion:accounts-add -->

1. काहीही तयार होण्यापूर्वी `config.ValidateName` नाव तपासते आणि `shim.Check` शॉर्टकट तपासते.
2. `profile.Add` `0700` mode सह `~/.local/share/ccam/profiles/work/` तयार करते आणि खात्याला सर्वात लहान मोकळा slot क्रमांक (`Config.NextSlot`) देते, त्यामुळे काढलेल्या खात्याचा क्रमांक पुन्हा वापरला जातो. पहिले खाते डीफॉल्ट होते.
3. `applyShare` नवीन, रिकामा फोल्डर `share.Build` आणि मग `Apply` ने त्याच्या शेअर गटांशी लिंक करते, उदाहरणार्थ `settings.json` ला `groups/default/settings.json` शी. हे अयशस्वी झाले तर नवीन फोल्डर काढला जातो, कारण उरलेला फोल्डर प्रत्येक पुढच्या प्रयत्नाला `folder already exists` ने अडवेल.
4. खाते `fsutil.AtomicWrite` मार्फत `~/.config/ccam/config.toml` मध्ये सेव्ह होते: आधी तात्पुरती फाइल, मग `fsync`, मग `rename`, त्यामुळे crash झाला तरी अर्धवट config कधीच उरत नाही.
5. `runLogin` खात्याच्या फोल्डरवर सेट केलेल्या `CLAUDE_CONFIG_DIR` सह `claude auth login` सुरू करते. लॉगिन अयशस्वी झाले तरी खाते राहते, आणि ccam त्रुटी छापून नंतर `ccam login work` चालवायला सांगते.
6. `profile.ReadEmail` `profiles/work/.claude.json` मधून `oauthAccount.emailAddress` काढते. ती फाइल अनेक मेगाबाइट्सची असू शकते, म्हणून `TopLevelValue` संपूर्ण फाइल parse करण्याऐवजी फक्त ती एक key शोधते.
7. `shim.Install` `~/.local/bin/ccwork` ला ccam binary शी लिंक करते, जी सुरू होताना `config.toml` मध्ये `argv[0]` शोधते.

लॉगिनला काही मिनिटे लागू शकतात, आणि त्या दरम्यान दुसरे ccam config सेव्ह करू शकते. म्हणून थांबणाऱ्या कमांड्स `env.update` मार्फत सेव्ह करतात, जे `config.toml` पुन्हा लोड करते, तो एक बदल (इथे ईमेल) लागू करते आणि सेव्ह करते. याच कारणासाठी `ccam remove` प्रश्न विचारल्यानंतर पुन्हा लोड करते.

`ccam login` आधीपासून असलेल्या खात्यासाठी टप्पे 5 आणि 6 चालवते. `ccam logout` खात्यात `claude auth logout` चालवते आणि सेव्ह केलेला ईमेल काढून टाकते. `ccam use` फक्त `config.toml` मधील `default` पुन्हा लिहिते.

### रिमोट मशीनवर साइन इन करणे

SSH वरून `claude auth login` चा browser वाला भाग बिघडू शकतो. Claude Code सर्व्हरवरील `http://localhost:<port>/callback` या callback सर्व्हरवर साइन-इनची वाट पाहते, पण साइन-इन पूर्ण करणारा browser तुमच्या कॉम्प्युटरवर चालतो, जिथे तो पत्ता कुठेच जात नाही आणि browser `ERR_CONNECTION_REFUSED` सारखी त्रुटी दाखवतो. तरीही त्या पत्त्यात एकदाच वापरता येणारा code असतो, त्यामुळे ccam तो पोहोचवू शकते:

<!-- motion:accounts-relay -->

`runLogin` तुमचा टर्मिनल Claude Code च्या हाती देत नाही. ते Claude Code चे standard input एका pipe ला जोडते आणि `relayLogin` तुम्ही टाइप केलेली प्रत्येक ओळ त्यात कॉपी करते, त्यामुळे Claude Code च्या `Paste code here` prompt वर paste केलेला code तरीही पोहोचतो. एखादी ओळ callback पत्ता असेल, तर `callbackURL` ती फक्त तेव्हाच स्वीकारते जेव्हा host `localhost`, `127.0.0.1` किंवा `[::1]` असेल, port दिलेला असेल, path `/callback` असेल आणि `code` व `state` दोन्ही असतील; `http://` वगळता येते. मग `deliverCallback` तो पत्ता सर्व्हरवर 10 सेकंदांच्या मर्यादेसह आणि redirects न पाळता request करते, जसे browser ने केले असते. ही request वाट पाहणाऱ्या `claude auth login` च्या callback सर्व्हरपर्यंत पोहोचते (आकृतीतील तुटक रेषा). Claude Code ला code मिळतो आणि ते लॉगिन सेव्ह करते, Linux वर खात्याच्या फोल्डरमधील `.credentials.json` मध्ये. paste केलेली ओळ ccam ला कधीही दुसऱ्या host शी संपर्क करायला लावू शकत नाही.

लॉगिन सुरू होण्यापूर्वी ccam हे छापते:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

यशस्वी झाल्यावर ते `→ Passed the sign-in to Claude Code on this machine.` छापते.

लॉगिन टर्मिनलवर input ची वाट पाहत असताना, ccam टर्मिनलचा echo आणि line mode बंद करते आणि तुम्ही टाइप केलेल्या प्रत्येक अक्षरासाठी एक `*` दाखवते, जेणेकरून paste केलेला code किंवा callback पत्ता स्क्रीनवर दिसत नाही. `readLine` line editing स्वतः करते: `Enter` ओळ पाठवते, `Backspace` एक अक्षर पुसते, आणि रिकाम्या ओळीवर `Ctrl+D` input संपवते. लॉगिन संपल्यावर टर्मिनल पूर्वस्थितीत येतो, आणि `Ctrl+C`, `SIGTERM` किंवा `SIGHUP` वरही येतो, जे तरीही ccam बंद करतात. standard input आणि standard error वर टर्मिनल नसल्यास, ओळी जशा येतात तशा कॉपी होतात.

### नाव बदलणे, काढणे आणि इम्पोर्ट करणे

`ccam rename` आणि `ccam remove` खात्याच्या फोल्डरवर आणि ccam खात्याच्या नावाखाली इतरत्र ठेवत असलेल्या गोष्टींवर काम करतात. `nameState` ते पाथ नोंदवते: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` आणि वापराचे caches `~/.local/state/ccam/usage/<name>.json` व `<name>.statusline.json`. `ccam remove work` ते या क्रमाने हाताळते:

<!-- motion:accounts-remove -->

खाते चालू असेल तर `ccam remove` आधी नकार देते, मग विचारते. खाते एक्स्टेंशन्स शेअर करत असेल, तर `RewritePluginPaths` गटाच्या `installed_plugins.json` आणि `known_marketplaces.json` पुन्हा लिहिते, कारण Claude Code plugins चे absolute पाथ साठवते, जे या खात्याच्या फोल्डरमधून जाऊ शकतात. `profile.Trash` फोल्डर `~/.local/share/ccam/trash/<name>-<time>/` मध्ये हलवते, खाते `config.toml` मधून निघते, `nameState` पाथ हटवले जातात जेणेकरून नंतरच्या त्याच नावाच्या खात्याला ते वारसाने मिळू नयेत, आणि शॉर्टकट लिंक काढली जाते. परत मिळवण्यासाठी फोल्डर trash मधून परत हलवा.

`ccam rename` फोल्डरचे नाव बदलते, त्याच्या `plugins/` मधील plugin पाथ जुन्या फोल्डरवरून नवीन फोल्डरवर पुन्हा लिहिते, डीफॉल्ट सोबत नेते, आणि `nameState` पाथ नवीन नावावर हलवते. शॉर्टकट कमांडचे नाव तेच राहते. plugin पाथ पुन्हा लिहिणे अयशस्वी झाले तर फोल्डरचे नाव परत पूर्वीसारखे केले जाते, कारण config सेव्ह झालेला नसतो.

`ccam import` फोल्डर आहे त्याच जागी, `config.toml` मध्ये `dir` म्हणून नोंदवते. ते त्याला कधीही हलवत किंवा कॉपी करत नाही: कॉपी केलेल्या refresh token मुळे प्रती एकमेकींना लॉग आउट करतील. फोल्डर न दिल्यास ते न नोंदवलेले `~/.claude` आणि `~/.claude-*` फोल्डर शोधते आणि प्रत्येकाबद्दल विचारते; नावे फोल्डरवरून ठरतात (`~/.claude` म्हणजे `main`, `~/.claude-work` म्हणजे `work`, `~/.claude-2` म्हणजे `acct2`), आणि ज्या फोल्डरचे नाव ठरवता येत नाही तो वापरायच्या `--name` कमांडसह दाखवला जातो. `CLAUDE_CONFIG_DIR` नसताना Claude Code `.claude.json` तुमच्या home फोल्डरमध्ये `~/.claude` शेजारी ठेवते, म्हणून तुम्ही `~/.claude` इम्पोर्ट करता आणि त्यात स्वतःची `.claude.json` नसते, तेव्हा ccam `~/.claude/.claude.json` ते `~/.claude.json` अशी लिंक तयार करते. मग साधे `claude` आणि `ccam main` एकच खाते पाहतात. इम्पोर्ट केलेली खाती तुम्ही `ccam share` चालवेपर्यंत काहीही शेअर करत नाहीत, आणि शॉर्टकट तयार होत नाही.

`ccam completion bash|zsh` एक completion script छापते, जी खात्यांची नावेही पूर्ण करते.

## सुरक्षितता आणि अपवादात्मक स्थिती

- खाते चालू असताना `ccam rename` आणि `ccam remove` (`refuseRunning`) `work is running; quit its Claude Code session first` देऊन नकार देतात, कारण session फोल्डरमध्ये लिहीत राहते.
- इम्पोर्ट केलेल्या खात्यावर `ccam remove` फक्त त्याची नोंदणी रद्द करते: `✓ main unregistered (folder ... left in place)`. तुमचा `~/.claude` trash मध्ये हलवला जात नाही.
- ज्या नावाचा फोल्डर `profiles/` मध्ये आधीच आहे ते नाव `ccam add` नाकारते, आणि `ccam import` कडे निर्देश करते.
- आधीच नोंदवलेला फोल्डर, फोल्डर नसलेला पाथ आणि आधीच घेतलेले नाव `ccam import` नाकारते.
- `ccam logout` आणि `ccam login` खाते चालू आहे का ते तपासत नाहीत.
- relay ला Claude Code पर्यंत पोहोचता आले नाही, उदाहरणार्थ लॉगिन संपल्यानंतर, तर ते कारणासह `ccam: could not pass the sign-in to Claude Code:` छापते. HTTP त्रुटीचा अर्थ सहसा code कालबाह्य झाला असा असतो; लॉगिन पुन्हा सुरू करा.

## संबंधित

- [Claude Code सुरू करणे](/features/launch)
- [शॉर्टकट कमांड्स](/features/shortcuts)
- [MCP सर्व्हर, सेटिंग्ज आणि इतिहास शेअर करा](/docs/sharing)
- [अनेक Claude Code खाती](/guides/multiple-claude-code-accounts)
