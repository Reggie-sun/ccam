---
title: ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం
description: ccam migrate తో ఇప్పటికే ఉన్న ~/.claude లోని సెట్టింగ్‌లు, CLAUDE.md, ఎక్స్‌టెన్షన్లు, MCP సర్వర్లను ఒక ccam గ్రూప్‌లోకి కాపీ చేయండి.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

ccam కంటే ముందు మీరు Claude Code ఉపయోగించి ఉంటే, మీ సెటప్ `~/.claude`, `~/.claude.json` లో ఉంటుంది. `ccam migrate` దాన్ని ఒక షేర్ గ్రూప్‌లోకి కాపీ చేస్తుంది, కాబట్టి ఆ గ్రూప్ ద్వారా ఆ అంశాలను షేర్ చేసే ప్రతి ఖాతా మీకు ఇప్పటికే ఉన్న సెటప్‌తో మొదలవుతుంది. మూల ఫోల్డర్ చదవబడుతుంది మాత్రమే, సాధారణ `claude` మునుపటిలాగే పనిచేస్తుంది. మీ మొదటి `ccam add` తర్వాత దీన్ని ఒకసారి రన్ చేయండి; `~/.claude.json` లో MCP సర్వర్లు ఉంటే, మొదటి `ccam add` తన సూచనలో వాటి పేర్లు చెబుతుంది.

## వినియోగం

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

టెర్మినల్‌లో, అంశాల ఫ్లాగ్‌లు లేకుండా, `--yes` లేకుండా, ccam ముందుగా ఒక చెక్‌లిస్ట్‌ను చూపిస్తుంది:

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

Space (లేదా `x`) కర్సర్ కింద ఉన్న అంశాన్ని మార్చుతుంది, బాణం కీలు లేదా `j`/`k` కదిలిస్తాయి, Enter ప్లాన్‌ను చూపిస్తుంది, Esc, `q` లేదా Ctrl+C `Cancelled.` అని ప్రింట్ చేస్తాయి. మూల ఫోల్డర్‌లో ఏమీ లేని అంశాలకు `not found` అని గుర్తు ఉంటుంది. మీరు అన్నింటి టిక్ తీసేస్తే, ccam `Nothing selected.` అని ప్రింట్ చేసి ఆగుతుంది.

ఆపై ccam ప్రణాళికను ప్రింట్ చేసి, కాపీ చేసే ముందు అడుగుతుంది:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

ఒక గ్రూప్ తన సొంత ఖాతాలకు మాత్రమే, తాను షేర్ చేసే అంశాలకు మాత్రమే అందుతుంది, కాబట్టి సారాంశం గ్రూప్ ఖాతాలను జాబితా చేస్తుంది. కాపీ చేసిన అంశాల్లో ఒకదాన్ని గ్రూప్ షేర్ చేయకపోతే, లేదా కొన్ని ఖాతాలు అందులో లేకపోతే, దాన్ని మార్చే `ccam group set <group> --add <items>`, `ccam group join <group> <accounts>` కమాండ్‌లను ccam ప్రింట్ చేస్తుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `[folder]` | కాపీ చేయాల్సిన Claude Code ఫోల్డర్ | `~/.claude` |
| `--settings` | `settings.json` ను కాపీ చేస్తుంది | ఆఫ్ |
| `--memory` | గ్లోబల్ `CLAUDE.md` ను కాపీ చేస్తుంది | ఆఫ్ |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` ను కాపీ చేస్తుంది | ఆఫ్ |
| `--mcp` | యూజర్ MCP సర్వర్లను (`.claude.json` లోని `mcpServers`) కాపీ చేస్తుంది | ఆఫ్ |
| `--history` | సంభాషణ చరిత్రను (`projects/`, `history.jsonl`) కాపీ చేస్తుంది | ఆఫ్ |
| `-g`, `--group string` | కాపీ చేయాల్సిన షేర్ గ్రూప్ | `default` |
| `-y`, `--yes` | అడగకుండా కాపీ చేస్తుంది; ఘర్షణ అయితే గ్రూప్ కాపీ ఉంచబడుతుంది | ఆఫ్ |
| `--dry-run` (గ్లోబల్) | ప్రణాళికను ప్రింట్ చేసి, ఏ ప్రశ్న లేదా కాపీకి ముందే ఆగుతుంది | ఆఫ్ |

అంశాల ఫ్లాగ్ లేకపోతే, టెర్మినల్‌లో చెక్‌లిస్ట్ నిర్ణయిస్తుంది. టెర్మినల్ బయట, లేదా `--yes` తో, ccam సంభాషణ చరిత్ర తప్ప అన్నీ కాపీ చేస్తుంది. `projects/` పెద్దది కావచ్చు కాబట్టి చరిత్ర డిఫాల్ట్‌గా వదిలేయబడుతుంది.

`--group on` అంటే `default`. `--group off`, `--group local` `--group takes a group name, not "off"` తో తిరస్కరించబడతాయి. `ccam migrate` `--dry-run` ను పాటిస్తుంది, మీరు ఒక ఎంపికను ముందుగా చూడటానికి డ్రై రన్‌లో కూడా చెక్‌లిస్ట్ కనిపిస్తుంది.

## ఇది ఎలా పనిచేస్తుంది

`internal/cli/migrate.go` అంశాలను ఎంచుకుంటుంది, ఆపై `migrate.Build` మూల ఫోల్డర్‌ను, గ్రూప్‌ను చదివి ఒక ప్రణాళికను తిరిగి ఇస్తుంది. మీరు నిర్ధారించిన తర్వాత, `Plan.Apply` వరకు ఏదీ రాయబడదు.

1. **మూలాన్ని కనుగొనడం.** లేని, ఫోల్డర్ కాని, లేదా తెలిసిన Claude Code ఎంట్రీ ఏదీ లేని పాత్‌ను `checkSource` తిరస్కరిస్తుంది. MCP సర్వర్ల కోసం మాత్రమే అయితే, `~/.claude.json` లో సర్వర్లు ఉన్నంత వరకు ఖాళీగా ఉన్న లేదా అసలు లేని `~/.claude` సరిపోతుంది. `migrate.ClaudeJSON` ఫోల్డర్‌కు చెందిన `.claude.json` ను కనుగొంటుంది: దాని లోపల, లేదా లోపల లేని `~/.claude` కోసం, దాని పక్కన ఉన్న `~/.claude.json`. `CLAUDE_CONFIG_DIR` సెట్ కానప్పుడు Claude Code దాన్ని అక్కడే ఉంచుతుంది.
2. **అంశాలను పాత్‌లకు మ్యాప్ చేయడం.** ప్రతి పాత్ `share.Layout` నుంచి వస్తుంది, గ్రూప్‌లు ఉపయోగించే అదే అనుమతి జాబితా. MCP సర్వర్లు `mcp.LocalServers` తో `.claude.json` నుంచి చదవబడతాయి, అది మొత్తం ఫైల్‌ను పార్స్ చేయకుండా `profile.TopLevelValue` స్కానర్‌ను ఉపయోగిస్తుంది.
3. **కాపీ చేయడం, ఎప్పుడూ తరలించకపోవడం.** ఫైళ్లు `fsutil.AtomicWrite` లేదా `fsutil.CopyFile` తో `~/.local/share/ccam/groups/<group>/` కు వెళ్తాయి, అది లక్ష్యం పక్కన ఒక తాత్కాలిక ఫైల్‌ను రాసి దాని పేరు మారుస్తుంది, కాబట్టి మధ్యలో ఆగిన కాపీ ఎప్పుడూ సగం కత్తిరించిన ఫైల్‌ను వదలదు. MCP సర్వర్లు గ్రూప్ `mcp.json` లోకి దాని లాక్ కింద వెళ్తాయి: ఈలోగా మరో ccam జోడించిన సర్వర్ అలాగే ఉంటుంది, [Remote Control](/guides/remote-control) కోసం చేసిన కాపీలు ఆ తర్వాత అప్‌డేట్ అవుతాయి.

ప్రతి అంశం ఒక మూల పాత్‌ను ఒక గ్రూప్ పాత్‌కు మ్యాప్ చేస్తుంది:

<!-- motion:migrate-items -->

రెండు చరిత్ర వరుసలు `--history` తో మాత్రమే కాపీ అవుతాయి కాబట్టి అవి గీతలతో ఉన్నాయి. `.credentials.json`, `.claude.json` `share.Never` లో జాబితా చేయబడ్డాయి, ఎప్పుడూ గ్రూప్‌ను చేరవు: వాటిలో ఒకటి ఎప్పుడైనా లేఅవుట్‌లో కనిపిస్తే `Build` అంతర్గత ఎర్రర్‌ను ఇస్తుంది, ఫోల్డర్ వాక్ వాటిని పేరు ద్వారా దాటవేస్తుంది.

`buildEntry` ప్రతి పాత్‌కు ఒక స్థితిని ఇస్తుంది, ప్రతి ప్లాన్ లైన్ వివరించేది అదే:

<!-- motion:migrate-states -->

- **Missing**: మూలంలో లేదు, లేదా ఫోల్డర్ ఆశించిన చోట ఫైల్ ఉంది. ప్రతి అంశానికి ఒకసారి `not found, skipped` గా జాబితా చేయబడుతుంది.
- **Linked**: మూల పాత్ ఇప్పటికే గ్రూప్ కాపీకి చేరుతుంది, ఉదాహరణకు మూలం ఒక ccam ప్రొఫైల్ ఫోల్డర్ అయినప్పుడు.
- **Same**: గ్రూప్‌లో ఇప్పటికే అదే కంటెంట్ ఉంది.
- **New**: గ్రూప్‌లో అది లేదు, కాబట్టి మొత్తంగా కాపీ అవుతుంది.
- **Merge**: రెండు వైపులా ఉన్న ఫోల్డర్. `copyMissing` గ్రూప్‌లో లేని ఫైళ్లను మాత్రమే జోడిస్తుంది; వేరుగా ఉన్న ఫైళ్లు గ్రూప్‌లో ఉన్నట్లే ఉంటాయి, ప్లాన్ `group copies kept` అని చెబుతుంది. `plugins/` కోసం, రిజిస్ట్రీలు `installed_plugins.json`, `known_marketplaces.json` మెర్జ్ అవుతాయి, పూర్తి ప్లగిన్ పాత్‌లు గ్రూప్ ఫోల్డర్‌కు చూపేలా మళ్లీ రాయబడతాయి.
- **Conflict**: రెండు వైపులా వేరే కంటెంట్‌తో ఉన్న ఫైల్. మీరు ప్రతి ఫైల్‌కు ఎంచుకుంటారు: `keep group`, `use mine`, అలాగే `.json` ఫైళ్లకు `merge JSON` (ఘర్షణ పడే విలువలలో గ్రూప్ గెలిచే డీప్ మెర్జ్) లేదా `history.jsonl` కు `merge by line`. నిర్వచనం వేరుగా ఉన్న MCP సర్వర్ కూడా `keep group` లేదా `use mine` తో ఘర్షణే. మీరు Enter నొక్కినప్పుడు జాబితాలోని మొదటి ఎంపిక డిఫాల్ట్; `--yes` ఎప్పుడూ గ్రూప్ కాపీని ఉంచుతుంది.

చెక్‌లిస్ట్ స్వయంగా `ccam add` ఉపయోగించే `checklist` కాంపోనెంట్‌పై నిర్మించిన `tui.MigratePicker`; `migrate.Found` lstat కాల్స్‌తో, MCP సర్వర్ జాబితాతో లేని అంశాలను గుర్తిస్తుంది.

### MCP సైన్-ఇన్‌లు

OAuth తో సైన్ ఇన్ అయ్యే MCP సర్వర్లు తమ టోకెన్‌లను సర్వర్ జాబితాలో కాకుండా లాగిన్ ఫైల్‌లో, `mcpOAuth` కింద `.credentials.json` లో ఉంచుతాయి. ccam వాటిని కాపీ చేయదు. కాపీ చేసిన సైన్-ఇన్‌లో అసలు దానిలోని అదే refresh token ఉంటుంది, మొదటి రెన్యూవల్‌లో సర్వర్ ఆ టోకెన్‌ను మారుస్తుంది: అప్పుడు మిగతా ప్రతి హోల్డర్, మూల `~/.claude` తో సహా, ఆ సర్వర్ నుంచి లాగ్ అవుట్ అవుతుంది. కాబట్టి ప్లాన్ లైన్ `sign-ins`, మూలం సైన్ ఇన్ అయిన ప్రతి కాపీ చేసిన సర్వర్ కోసం, `/mcp` తో తామే సైన్ ఇన్ చేయాల్సిన గ్రూప్ ఖాతాలను పేర్కొంటుంది, సారాంశం దాన్ని మళ్లీ చెబుతుంది. టోకెన్ లేని ఎంట్రీ సైన్-ఇన్ కాదు, అది జాబితాలో ఉండదు; ఇప్పటికే సొంత సైన్-ఇన్ ఉన్న ఖాతా, గ్రూప్ తన నిర్వచనాన్ని ఉంచుకునే సర్వర్ కూడా జాబితాలో ఉండవు. ఏ లాగిన్ ఫైల్ రాయబడదు.

## భద్రత, ప్రత్యేక సందర్భాలు

- మూలం చదవబడుతుంది మాత్రమే. అందులోని ఏ ఫైల్ తరలించబడదు, లింక్ చేయబడదు, మార్చబడదు, మూల లాగిన్‌లో ఎప్పుడూ రాయబడదు, Keychain నుంచి చదవబడే అవకాశం ఉన్న macOS లో కూడా.
- `use mine`, మెర్జ్, లేదా ఇప్పటికే ఉన్న `mcp.json` లో కొత్త MCP సర్వర్ల ద్వారా గ్రూప్ ఫైల్ భర్తీ అయ్యే ముందు, గ్రూప్ కాపీ `<group>/.ccam-backup/<time>/` లో సేవ్ అవుతుంది. ఫోల్డర్ పేరు ప్రతి రన్‌కు ప్రత్యేకం, కాబట్టి ఒకే సెకనులో జరిగే రెండు రన్‌లు ఒకదాని బ్యాకప్‌లను మరొకటి ఓవర్‌రైట్ చేయవు. సారాంశం పాత్‌ను ప్రింట్ చేస్తుంది.
- `--yes` తో, ప్రతి ఘర్షణ గ్రూప్ కాపీని ఉంచుతుంది, బదులుగా ఎలా ఎంచుకోవాలో ccam చెబుతుంది: అదే కమాండ్‌ను `--yes` లేకుండా రన్ చేయండి.
- ఏదీ మారకపోతే, ccam `Nothing to copy.` అని ప్రింట్ చేస్తుంది; ప్రతి ఘర్షణ గ్రూప్ కాపీని ఉంచితే, `Nothing to copy: the group keeps its own copies.`
- కాపీ మధ్యలో విఫలమైతే, ccam `The source folder was not changed. Group files replaced so far are saved in` తో, దాని తర్వాత బ్యాకప్ ఫోల్డర్‌తో ఆగుతుంది. మళ్లీ రన్ చేస్తే పూర్తిగా చేరని ఫైళ్లు కాపీ అవుతాయి.
- అనేక గ్రూప్‌లు MCP ని షేర్ చేసినప్పుడు, `ccam doctor`, `ccam ext`, `ccam mcp list` సూచించే కమాండ్‌లు `ccam migrate --mcp --group <a|b>` గా ఉంటాయి.
- మూలం లేకపోతే `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` తో ముగుస్తుంది.

## సంబంధిత పేజీలు

- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [తనిఖీ చేయడం, సరిచేయడం](/features/doctor)
- [అనేక Claude Code ఖాతాలు](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
