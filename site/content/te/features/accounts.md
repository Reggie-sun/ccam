---
title: ఖాతాలను నిర్వహించడం
description: ccam add, ccam login, ccam list, ccam rename, ccam remove, ccam import తో Claude Code ఖాతాలను జోడించండి, లాగిన్ చేయండి, జాబితా చూడండి, పేరు మార్చండి, తీసివేయండి, ఇంపోర్ట్ చేయండి.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

ప్రతి ccam ఖాతా తన సొంత లాగిన్‌తో ఉన్న వేరే Claude Code కాన్ఫిగ్ ఫోల్డర్. ఈ కమాండ్‌లు ఖాతాలను సృష్టిస్తాయి, సైన్ ఇన్ చేస్తాయి, జాబితాగా చూపిస్తాయి, పేరు మారుస్తాయి, తీసివేస్తాయి, మీకు ఇప్పటికే ఉన్న ఫోల్డర్‌లను నమోదు చేస్తాయి.

## వినియోగం

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

ఏ కమాండ్ అయినా `<name|number>` తీసుకునే చోట, మీరు ఖాతా పేరును లేదా దాని స్లాట్ నంబర్‌ను, అంటే `ccam list` లోని `#` కాలమ్‌ను ఇవ్వవచ్చు.

ఆప్షన్లు లేకుండా టెర్మినల్‌లో `ccam add` ఒక విజార్డ్‌ను తెరుస్తుంది. అది పేరు అడుగుతుంది, ఆపై షేర్ చేసే అంశాల చెక్‌లిస్ట్‌ను చూపిస్తుంది (సంభాషణ చరిత్ర తప్ప అన్నీ టిక్ చేసి ఉంటాయి), ఆపై `Create a shortcut command for work? [y/N]` అని అడుగుతుంది, `y` అంటే `Shortcut name [ccwork]:` అని అడుగుతుంది. `--share`, `--shortcut`, `--no-login` లో ఏదైనా ఉంటే, లేదా టెర్మినల్ బదులు పైప్ ఉంటే, విజార్డ్ దాటవేయబడుతుంది; అప్పుడు పేరు తప్పనిసరి, `--shortcut` ఒక పేరు ఇస్తే తప్ప షార్ట్‌కట్ సృష్టించబడదు.

`ccam list` స్లాట్, పేరు, షార్ట్‌కట్, ఇమెయిల్, షేర్ చేసిన అంశాలు, చివరి వినియోగాన్ని చూపిస్తుంది, డిఫాల్ట్ ఖాతా ముందు `*` ఉంటుంది:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` తో, ప్రతి ఖాతా `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (ఖాతా ఎప్పుడూ ప్రారంభించబడకపోతే వదిలేయబడుతుంది), `share` ఉన్న ఒక ఆబ్జెక్ట్. `share` ప్రతి అంశాన్ని (`mcp`, `settings`, `memory`, `extensions`, `history`) దాని గ్రూప్ పేరుకు లేదా `"local"` కు మ్యాప్ చేస్తుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `ccam add --no-login` | లాగిన్ చేయకుండా ఖాతాను సృష్టిస్తుంది; తర్వాత `ccam login <name>` తో లాగిన్ చేయండి | ఆఫ్ |
| `ccam add --share <item=group>` | ఒక అంశాన్ని ఒక గ్రూప్‌తో షేర్ చేస్తుంది, లేదా దాన్ని ప్రైవేట్‌గా ఉంచడానికి `item=local`. మళ్లీ మళ్లీ ఇవ్వవచ్చు | `default` లో `mcp`, `settings`, `memory`, `extensions`; `history` లోకల్ |
| `ccam add --shortcut <command>` | `~/.local/bin` లో ఈ షార్ట్‌కట్ కమాండ్‌ను కూడా సృష్టిస్తుంది | ఏమీ లేదు |
| `ccam list --json` | ఖాతాలను JSON గా ప్రింట్ చేస్తుంది | ఆఫ్ |
| `ccam remove -y`, `--yes` | అడగకుండా తీసివేస్తుంది | ఆఫ్ |
| `ccam import --name <name>` | మీరు ఇచ్చే ఫోల్డర్ కోసం ఖాతా పేరు | ఫోల్డర్ పేరు నుంచి తీసుకోబడుతుంది |
| `ccam import -y`, `--yes` | కనిపించిన ప్రతి ఫోల్డర్‌ను అడగకుండా ఇంపోర్ట్ చేస్తుంది | ఆఫ్ |
| `--dry-run` | ఏమీ మార్చని `ccam list`, `ccam completion` దీన్ని పాటిస్తాయి. `add`, `login`, `logout`, `use`, `rename`, `remove`, `import` దీన్ని `` --dry-run is not supported by `ccam add`; nothing was changed `` (కమాండ్ పేరుతో) అని తిరస్కరిస్తాయి | ఆఫ్ |

`ccam login`, `ccam logout`, `ccam use`, `ccam rename` లకు సొంత ఆప్షన్లు లేవు. `ccam completion` `bash` లేదా `zsh` తీసుకుంటుంది; మరేదైనా `only bash and zsh are supported` తో విఫలమవుతుంది.

ఖాతా పేర్లు `a-z`, `0-9`, `-`, `_` అక్షరాలతో 1 నుంచి 32 అక్షరాల వరకు ఉంటాయి, ఒక అక్షరం లేదా అంకెతో మొదలవుతాయి. పేరు మొత్తం అంకెలు కాకూడదు, ఎందుకంటే `ccam 2` అంటే స్లాట్ 2. పేరు ccam కమాండ్ కూడా కాకూడదు (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), కాబట్టి `ccam <name>` ఎప్పుడూ అస్పష్టంగా ఉండదు. ఈ జాబితా `internal/config` లోని `config.Reserved`.

## ఇది ఎలా పనిచేస్తుంది

`ccam add work --shortcut ccwork` `addProfile` (`internal/cli/profiles.go`) లో ఈ దశలను రన్ చేస్తుంది:

<!-- motion:accounts-add -->

1. ఏదైనా సృష్టించే ముందు `config.ValidateName` పేరును, `shim.Check` షార్ట్‌కట్‌ను తనిఖీ చేస్తాయి.
2. `profile.Add` `0700` మోడ్‌తో `~/.local/share/ccam/profiles/work/` ను సృష్టించి, ఖాతాకు అతి తక్కువ ఖాళీ స్లాట్ నంబర్‌ను (`Config.NextSlot`) ఇస్తుంది, కాబట్టి తీసివేసిన ఖాతా నంబర్ మళ్లీ ఉపయోగించబడుతుంది. మొదటి ఖాతా డిఫాల్ట్ అవుతుంది.
3. `applyShare` కొత్త, ఖాళీ ఫోల్డర్‌ను `share.Build` తో, ఆపై `Apply` తో దాని షేర్ గ్రూప్‌లకు లింక్ చేస్తుంది, ఉదాహరణకు `settings.json` ను `groups/default/settings.json` కు. ఇది విఫలమైతే, కొత్త ఫోల్డర్ తీసివేయబడుతుంది, ఎందుకంటే మిగిలిపోయిన ఫోల్డర్ ప్రతి మళ్లీ ప్రయత్నాన్ని `folder already exists` తో అడ్డుకుంటుంది.
4. ఖాతా `fsutil.AtomicWrite` ద్వారా `~/.config/ccam/config.toml` లో సేవ్ అవుతుంది: ఒక తాత్కాలిక ఫైల్, `fsync`, ఆపై `rename`, కాబట్టి క్రాష్ అయినా సగం కాన్ఫిగ్ ఎప్పుడూ మిగలదు.
5. `runLogin` `CLAUDE_CONFIG_DIR` ను ఖాతా ఫోల్డర్‌కు సెట్ చేసి `claude auth login` ను ప్రారంభిస్తుంది. లాగిన్ విఫలమైతే, ఖాతా అలాగే ఉంటుంది, ccam ఎర్రర్‌ను ప్రింట్ చేసి తర్వాత `ccam login work` రన్ చేయమని చెబుతుంది.
6. `profile.ReadEmail` `profiles/work/.claude.json` నుంచి `oauthAccount.emailAddress` ను తీసుకుంటుంది. ఆ ఫైల్ అనేక మెగాబైట్లు ఉండవచ్చు, కాబట్టి `TopLevelValue` మొత్తం ఫైల్‌ను పార్స్ చేయకుండా ఆ ఒక్క కీ కోసం స్కాన్ చేస్తుంది.
7. `shim.Install` `~/.local/bin/ccwork` ను ccam బైనరీకి లింక్ చేస్తుంది, అది ప్రారంభమైనప్పుడు `config.toml` లో `argv[0]` కోసం చూస్తుంది.

లాగిన్‌కు నిమిషాలు పట్టవచ్చు, ఈలోగా మరో ccam కాన్ఫిగ్‌ను సేవ్ చేయవచ్చు. కాబట్టి వేచి ఉండే కమాండ్‌లు `env.update` ద్వారా సేవ్ చేస్తాయి, అది `config.toml` ను మళ్లీ లోడ్ చేసి, ఆ ఒక్క మార్పును (ఇక్కడ, ఇమెయిల్) వర్తింపజేసి సేవ్ చేస్తుంది. అదే కారణంతో `ccam remove` తన ప్రశ్న తర్వాత మళ్లీ లోడ్ చేస్తుంది.

`ccam login` ఇప్పటికే ఉన్న ఖాతా కోసం దశ 5, 6 ను రన్ చేస్తుంది. `ccam logout` ఖాతాలో `claude auth logout` ను రన్ చేసి సేవ్ చేసిన ఇమెయిల్‌ను తొలగిస్తుంది. `ccam use` `config.toml` లోని `default` ను మాత్రమే మళ్లీ రాస్తుంది.

### రిమోట్ మెషీన్‌లో సైన్ ఇన్ చేయడం

SSH ద్వారా, `claude auth login` యొక్క బ్రౌజర్ భాగం పాడవ్వవచ్చు. Claude Code సర్వర్‌లో `http://localhost:<port>/callback` వద్ద ఉన్న కాల్‌బ్యాక్ సర్వర్‌లో సైన్-ఇన్ కోసం వేచి ఉంటుంది, కానీ సైన్-ఇన్‌ను పూర్తి చేసే బ్రౌజర్ మీ కంప్యూటర్‌లో రన్ అవుతుంది, అక్కడ ఆ అడ్రస్ ఎక్కడికీ దారితీయదు, బ్రౌజర్ `ERR_CONNECTION_REFUSED` వంటి ఎర్రర్‌ను చూపిస్తుంది. అయినా ఆ అడ్రస్‌లో వన్-టైమ్ కోడ్ ఉంటుంది, కాబట్టి ccam దాన్ని అందించగలదు:

<!-- motion:accounts-relay -->

`runLogin` మీ టెర్మినల్‌ను Claude Code కు అప్పగించదు. అది Claude Code స్టాండర్డ్ ఇన్‌పుట్‌ను ఒక పైప్‌కు కలుపుతుంది, మీరు టైప్ చేసే ప్రతి లైన్‌ను `relayLogin` అందులోకి కాపీ చేస్తుంది, కాబట్టి Claude Code యొక్క `Paste code here` ప్రాంప్ట్ వద్ద పేస్ట్ చేసిన కోడ్ కూడా చేరుతుంది. ఒక లైన్ కాల్‌బ్యాక్ అడ్రస్ అయినప్పుడు, హోస్ట్ `localhost`, `127.0.0.1` లేదా `[::1]` అయి, పోర్ట్ ఇవ్వబడి, పాత్ `/callback` అయి, `code`, `state` రెండూ ఉంటేనే `callbackURL` దాన్ని అంగీకరిస్తుంది; `http://` వదిలేయవచ్చు. ఆపై `deliverCallback` బ్రౌజర్ చేసి ఉండే విధంగానే, 10 సెకన్ల పరిమితితో, రీడైరెక్ట్‌లను అనుసరించకుండా, సర్వర్‌లో ఆ అడ్రస్‌ను అభ్యర్థిస్తుంది. ఆ అభ్యర్థన వేచి ఉన్న `claude auth login` యొక్క కాల్‌బ్యాక్ సర్వర్‌ను చేరుతుంది (డయాగ్రామ్‌లోని గీతల లైన్). Claude Code కోడ్‌ను అందుకుని లాగిన్‌ను సేవ్ చేస్తుంది, Linux లో ఖాతా ఫోల్డర్‌లోని `.credentials.json` లో. పేస్ట్ చేసిన లైన్ ccam ను మరో హోస్ట్‌ను సంప్రదించేలా ఎప్పుడూ చేయలేదు.

లాగిన్ మొదలయ్యే ముందు ccam ఇది ప్రింట్ చేస్తుంది:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

విజయవంతమైతే అది `→ Passed the sign-in to Claude Code on this machine.` అని ప్రింట్ చేస్తుంది.

లాగిన్ టెర్మినల్‌లో ఇన్‌పుట్ కోసం వేచి ఉన్నప్పుడు, ccam టెర్మినల్ echo ను, లైన్ మోడ్‌ను ఆపివేసి, మీరు టైప్ చేసే ప్రతి అక్షరానికి ఒక `*` చూపిస్తుంది, కాబట్టి పేస్ట్ చేసిన కోడ్ లేదా కాల్‌బ్యాక్ అడ్రస్ స్క్రీన్‌పై కనిపించదు. `readLine` లైన్ ఎడిటింగ్‌ను తానే చేస్తుంది: `Enter` లైన్‌ను పంపుతుంది, `Backspace` ఒక అక్షరాన్ని చెరిపేస్తుంది, ఖాళీ లైన్‌పై `Ctrl+D` ఇన్‌పుట్‌ను ముగిస్తుంది. లాగిన్ ముగిసినప్పుడు, అలాగే `Ctrl+C`, `SIGTERM` లేదా `SIGHUP` వచ్చినప్పుడు కూడా టెర్మినల్ పునరుద్ధరించబడుతుంది; ఇవి ఇప్పటికీ ccam ను ముగిస్తాయి. స్టాండర్డ్ ఇన్‌పుట్, స్టాండర్డ్ ఎర్రర్‌పై టెర్మినల్ లేకపోతే, లైన్‌లు వచ్చినట్లే కాపీ అవుతాయి.

### పేరు మార్చడం, తీసివేయడం, ఇంపోర్ట్ చేయడం

`ccam rename`, `ccam remove` ఖాతా ఫోల్డర్‌పై, ccam ఇతర చోట్ల ఖాతా పేరుతో ఉంచే వాటిపై పనిచేస్తాయి. `nameState` ఆ పాత్‌లను జాబితా చేస్తుంది: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, వినియోగ కాష్‌లు `~/.local/state/ccam/usage/<name>.json`, `<name>.statusline.json`. `ccam remove work` వాటిని ఈ క్రమంలో నిర్వహిస్తుంది:

<!-- motion:accounts-remove -->

ఖాతా రన్ అవుతుంటే `ccam remove` ముందుగా తిరస్కరిస్తుంది, ఆపై అడుగుతుంది. ఖాతా ఎక్స్‌టెన్షన్లను షేర్ చేస్తే, `RewritePluginPaths` గ్రూప్ యొక్క `installed_plugins.json`, `known_marketplaces.json` ను మళ్లీ రాస్తుంది, ఎందుకంటే Claude Code ఈ ఖాతా ఫోల్డర్ ద్వారా వెళ్లే అవకాశం ఉన్న పూర్తి ప్లగిన్ పాత్‌లను నిల్వ చేస్తుంది. `profile.Trash` ఫోల్డర్‌ను `~/.local/share/ccam/trash/<name>-<time>/` కు తరలిస్తుంది, ఖాతా `config.toml` నుంచి తొలగించబడుతుంది, అదే పేరుతో తర్వాత వచ్చే ఖాతాకు అవి వారసత్వంగా రాకుండా `nameState` పాత్‌లు తొలగించబడతాయి, షార్ట్‌కట్ లింక్ తీసివేయబడుతుంది. తిరిగి పొందడానికి, ఫోల్డర్‌ను ట్రాష్ నుంచి వెనక్కి తరలించండి.

`ccam rename` ఫోల్డర్ పేరును మారుస్తుంది, దాని `plugins/` లోని ప్లగిన్ పాత్‌లను పాత ఫోల్డర్ నుంచి కొత్తదానికి మళ్లీ రాస్తుంది, డిఫాల్ట్‌ను కూడా తీసుకువెళ్తుంది, `nameState` పాత్‌లను కొత్త పేరుకు తరలిస్తుంది. షార్ట్‌కట్ కమాండ్ తన పేరును ఉంచుకుంటుంది. ప్లగిన్ పాత్‌లను మళ్లీ రాయడం విఫలమైతే, కాన్ఫిగ్ సేవ్ కాలేదు కాబట్టి ఫోల్డర్ పేరు వెనక్కి మార్చబడుతుంది.

`ccam import` ఒక ఫోల్డర్‌ను ఉన్న చోటే, `config.toml` లో `dir` గా నమోదు చేస్తుంది. అది దాన్ని ఎప్పుడూ తరలించదు, కాపీ చేయదు: కాపీ చేసిన రిఫ్రెష్ టోకెన్ కాపీలు ఒకదాన్ని మరొకటి లాగ్ అవుట్ చేసేలా చేస్తుంది. ఫోల్డర్ ఆర్గ్యుమెంట్ లేకుండా, అది నమోదు కాని `~/.claude`, `~/.claude-*` ఫోల్డర్‌ల కోసం వెతికి ప్రతిదాని గురించి అడుగుతుంది; పేర్లు ఫోల్డర్ నుంచి వస్తాయి (`~/.claude` అంటే `main`, `~/.claude-work` అంటే `work`, `~/.claude-2` అంటే `acct2`), పేరు తీసుకోలేని ఫోల్డర్ ఉపయోగించాల్సిన `--name` కమాండ్‌తో జాబితాలో చూపబడుతుంది. `CLAUDE_CONFIG_DIR` లేకుండా, Claude Code `.claude.json` ను మీ హోమ్ ఫోల్డర్‌లో `~/.claude` పక్కన ఉంచుతుంది, కాబట్టి మీరు `~/.claude` ను ఇంపోర్ట్ చేసినప్పుడు దానికి సొంత `.claude.json` లేకపోతే, ccam `~/.claude/.claude.json` నుంచి `~/.claude.json` కు లింక్‌ను సృష్టిస్తుంది. అప్పుడు సాధారణ `claude`, `ccam main` ఒకే ఖాతాను చూస్తాయి. మీరు `ccam share` రన్ చేసే వరకు ఇంపోర్ట్ చేసిన ఖాతాలు ఏదీ షేర్ చేయవు, షార్ట్‌కట్ సృష్టించబడదు.

`ccam completion bash|zsh` ఖాతా పేర్లను కూడా పూర్తి చేసే కంప్లీషన్ స్క్రిప్ట్‌ను ప్రింట్ చేస్తుంది.

## భద్రత, ప్రత్యేక సందర్భాలు

- సెషన్ ఫోల్డర్‌లో రాస్తూనే ఉంటుంది కాబట్టి, ఖాతా రన్ అవుతున్నప్పుడు `ccam rename`, `ccam remove` `work is running; quit its Claude Code session first` తో తిరస్కరిస్తాయి (`refuseRunning`).
- ఇంపోర్ట్ చేసిన ఖాతాపై `ccam remove` దాని నమోదును మాత్రమే తొలగిస్తుంది: `✓ main unregistered (folder ... left in place)`. మీ `~/.claude` ట్రాష్‌కు తరలించబడదు.
- `profiles/` లో ఫోల్డర్ ఇప్పటికే ఉన్న పేరును `ccam add` తిరస్కరించి, `ccam import` ను సూచిస్తుంది.
- ఇప్పటికే నమోదైన ఫోల్డర్‌ను, ఫోల్డర్ కాని పాత్‌ను, తీసుకున్న పేరును `ccam import` తిరస్కరిస్తుంది.
- ఖాతా రన్ అవుతుందో లేదో `ccam logout`, `ccam login` తనిఖీ చేయవు.
- రిలే Claude Code ను చేరలేకపోతే, ఉదాహరణకు లాగిన్ ముగిసిన తర్వాత, అది కారణంతో `ccam: could not pass the sign-in to Claude Code:` అని ప్రింట్ చేస్తుంది. HTTP ఎర్రర్ సాధారణంగా కోడ్ గడువు ముగిసిందని అర్థం; లాగిన్‌ను మళ్లీ ప్రారంభించండి.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [షార్ట్‌కట్ కమాండ్‌లు](/features/shortcuts)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [అనేక Claude Code ఖాతాలు](/guides/multiple-claude-code-accounts)
