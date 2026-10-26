---
title: షేర్ చేసిన MCP సర్వర్లు
description: ccam mcp add, add-json, list, remove, import తో ప్రతి గ్రూప్‌కు ఒక MCP సర్వర్ల జాబితా, అది --mcp-config తో పంపబడుతుంది లేదా ప్రతి ఖాతాలో సింక్ అవుతుంది.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
---

Claude Code యూజర్ MCP సర్వర్లను `.claude.json` లో ఉంచుతుంది, ఖాతా లాగిన్‌ను కలిగి ఉండే అదే ఫైల్, కాబట్టి ఆ ఫైల్‌ను ఖాతాల మధ్య షేర్ చేయలేము. ccam షేర్ చేసిన సర్వర్లను `mcp.json` అనే గ్రూప్ ఫైల్‌లో ఉంచి, ప్రారంభంలో గ్రూప్‌లోని ప్రతి ఖాతాకు దాన్ని అందిస్తుంది. మీ అన్ని ఖాతాల కోసం ఒక సర్వర్‌ను ఒక్కసారి జోడించడానికి `ccam mcp` ను, ఒక సెషన్ లోపల మీరు జోడించిన సర్వర్‌ను ఇతర ఖాతాలకు తరలించడానికి `ccam mcp import` ను ఉపయోగించండి.

## వినియోగం

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

మార్పు ప్రతి ఖాతా తదుపరి ప్రారంభం నుంచి వర్తిస్తుంది; రన్ అవుతున్న సెషన్‌లు తాము ప్రారంభమైన సర్వర్లనే ఉంచుకుంటాయి:

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

గ్రూప్‌ల కింద, ఇతర ఖాతాలను చేరని ప్రతి సర్వర్‌ను `ccam mcp list` పేర్కొంటుంది:

- ఒక ఖాతా తన కోసం జోడించిన సర్వర్‌కు `only in account work: pg`;
- ఒకే ఫోల్డర్ కోసం జోడించిన సర్వర్‌కు `only in account side, for /home/me/code/app only: notes` (సాధారణ `claude mcp add`, కింద చూడండి);
- ఖాతా సొంత సర్వర్‌కు గ్రూప్ సర్వర్ పేరే ఉన్నప్పుడు `account work has its own …, used in its sessions in place of the group's`;
- ఈ ఫోల్డర్ `.mcp.json`, అది ఆమోదం కోసం వేచి ఉన్న ఖాతాలతో (`waiting for approval in`);
- `in no account, only in ~/.claude.json (plain claude's)` అయిన సర్వర్లు, వాటిని కాపీ చేసే [`ccam migrate --mcp`](/features/migrate) కమాండ్‌తో.

గ్రూప్‌లో ఇప్పటికే ఉన్న పేరును జోడిస్తే దాని నిర్వచనం భర్తీ అవుతుంది, `added` బదులు `updated` అని ప్రింట్ అవుతుంది. ఒక ఖాతా ఆ గ్రూప్‌లో ఉండి, ఆ గ్రూప్ MCP ను షేర్ చేస్తేనే (`ccam group set <group> --add mcp`) సర్వర్లు దానికి అందుతాయి; గ్రూప్ MCP ను షేర్ చేయకపోతే, `add` ఆ విషయం చెబుతుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json`, `remove` కోసం లక్ష్య గ్రూప్ | ఏదో ఒక ఖాతా MCP ను షేర్ చేసే ఏకైక గ్రూప్, లేకపోతే `default` |
| `-- <command> [args…]` (`add`) | ఈ కమాండ్‌ను stdio సర్వర్‌గా రన్ చేస్తుంది | |
| `--url <url>` (`add`) | రిమోట్ సర్వర్ URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` సర్వర్ యొక్క ట్రాన్స్‌పోర్ట్ | `http` |
| `--header <'Name: value'>` (`add`) | HTTP హెడర్, మళ్లీ మళ్లీ ఇవ్వవచ్చు | ఏమీ లేదు |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio సర్వర్ కోసం ఎన్విరాన్‌మెంట్ వేరియబుల్, మళ్లీ మళ్లీ ఇవ్వవచ్చు | ఏమీ లేదు |
| `--force` (`import`) | వేరుగా ఉన్న అదే పేరు గల గ్రూప్ సర్వర్‌ను భర్తీ చేస్తుంది; పాత `mcp.json` బ్యాకప్ అవుతుంది | ఆఫ్ |
| `--keep-local` (`import`) | ఖాతా `.claude.json` లోని ఎంట్రీలను ఉంచుతుంది | ఆఫ్ |
| `-h`, `--help` | ప్రతి కమాండ్‌కు సహాయం | |
| `--dry-run` (గ్లోబల్) | `mcp list` మాత్రమే దీనికి మద్దతు ఇస్తుంది; `add`, `add-json`, `remove`, `import` దాన్ని తిరస్కరించి ఏదీ మార్చవు | ఆఫ్ |

`add` కు `--` తర్వాత ఒక కమాండ్ గానీ `--url` గానీ అవసరం; లేకపోతే అది `a command to run (after --) or --url is required` తో విఫలమవుతుంది. `--header` `--url` సర్వర్లకు, `--env` కమాండ్ సర్వర్లకు వర్తిస్తాయి. `add-json` `claude mcp add-json` తీసుకునే అదే JSON ను తీసుకుంటుంది, నిర్వచనం `command` లేదా `url` ఉన్న ఆబ్జెక్ట్ అయి ఉండాలి. `import` ఎప్పుడూ ఖాతా MCP ను షేర్ చేసే గ్రూప్‌నే ఉపయోగిస్తుంది, `list` ప్రతి గ్రూప్‌ను చూపిస్తుంది, కాబట్టి `-g` వాటిని మార్చదు.

## ఇది ఎలా పనిచేస్తుంది

సంబంధిత ఫైళ్లు:

- `~/.local/share/ccam/groups/<group>/mcp.json`: షేర్ చేసిన సర్వర్లు, `{"mcpServers": {…}}`, `--mcp-config` చదివే ఫార్మాట్. సర్వర్ ఎంట్రీలలో తరచుగా API కీలు ఉంటాయి కాబట్టి, ఇది `fsutil.AtomicWrite` తో, `0600` మోడ్‌తో రాయబడుతుంది, అది కూడా `<group>/.mcp.json.lock` లాక్ కింద మాత్రమే, కాబట్టి ఒకే సమయంలో రెండు `ccam mcp add` కమాండ్‌లు రన్ అయినా రెండూ తమ సర్వర్లను ఉంచుకుంటాయి.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ఖాతా సర్వర్లు ఉన్న Claude Code సొంత ఫైల్. ccam `profile.TopLevelValue` స్కానర్‌తో దాని `mcpServers` విలువను మాత్రమే చదువుతుంది, ప్రారంభ మార్గంలో మొత్తం ఫైల్‌ను ఎప్పుడూ పార్స్ చేయదు.
- `~/.local/state/ccam/mcp/<name>/`: `seen` ఫైల్, `sync-stamp` (ccam రాసిన ప్రతి కాపీ హాష్‌తో, కాబట్టి ఎవరూ మార్చని కాపీని మీరు ఎడిట్ చేసిన దాని నుంచి వేరు చేయగలదు), `sync-checked` (sync విధానం చివరిగా చూసింది), `kept` (MCP షేరింగ్ నుంచి బయటకు వెళ్లినప్పుడు ఖాతా ఉంచుకున్న సర్వర్లు), `.claude.json` బ్యాకప్‌లు.

`~/.config/ccam/config.toml` లోని `mcp_strategy` తో సెట్ చేసే రెండు విధానాలు ఉన్నాయి. డిఫాల్ట్ ఫ్లాగ్ విధానం:

<!-- motion:mcp-flag -->

1. ccam ఖాతా MCP గ్రూప్‌ను కనుగొని (`launcher.applyMCP`), గ్రూప్ `mcp.json` ఉందో లేదో, దాన్ని పార్స్ చేయకుండా, ఒక `stat` తో తనిఖీ చేస్తుంది (`mcp.FlagArgs`).
2. ccam `--mcp-config=<path>` ను `claude` తర్వాత వెంటనే, మీ సొంత ఆర్గ్యుమెంట్‌ల ముందు పెడుతుంది, కాబట్టి మీ సొంత `--mcp-config` లేదా `--strict-mcp-config` దానిపై అదనంగా వర్తిస్తాయి (`Plan.InsertArgs`).
3. చివరి ప్రారంభం తర్వాత ఖాతా `.claude.json` మారిందో లేదో ccam తనిఖీ చేసి, అప్పుడు మాత్రమే దాని సర్వర్లను పార్స్ చేస్తుంది (`mcp.NewLocal`).
4. ccam తన స్థానంలో Claude Code ను పెడుతుంది (`syscall.Exec`).

పాత్ రెండు కాకుండా ఒకే ఆర్గ్యుమెంట్‌గా, `--mcp-config=<path>` గా పంపబడుతుంది. Claude Code యొక్క `--mcp-config` అనేక విలువలను తీసుకుంటుంది, కాబట్టి `claude --mcp-config a.json "prompt"` ప్రాంప్ట్‌ను రెండో కాన్ఫిగ్ ఫైల్‌గా చదివి `MCP config file not found` తో విఫలమవుతుంది. ఒకే ఆర్గ్యుమెంట్ రూపం దీన్ని నివారిస్తుంది, కాబట్టి దాని తర్వాత వచ్చే ప్రాంప్ట్ లేదా సబ్‌కమాండ్ పనిచేస్తుంది. `ccam work mcp list`, `get`, `remove` కు ఫ్లాగ్ రాదు, ఎందుకంటే Claude Code అక్కడ దాన్ని తిరస్కరిస్తుంది; బదులుగా ccam స్టాండర్డ్ ఎర్రర్‌పై గ్రూప్ సర్వర్ల పేర్లు చెబుతుంది ([వివరాలు](/features/launch#shared-mcp-servers)). Claude Code ఖాతా సొంత సర్వర్లతో పాటు `--mcp-config` సర్వర్లను లోడ్ చేస్తుంది; పేరు ఘర్షణ అయితే `.claude.json` లోని ఎంట్రీ గెలుస్తుంది. ఫ్లాగ్ విధానం Remote Control కోసం, ఒక ఖాతా MCP షేరింగ్ నుంచి బయటకు వెళ్లినప్పుడు లేదా మళ్లీ చేరినప్పుడు ([కింద](#leaving-and-joining)), `ccam mcp import` కోసం మాత్రమే `.claude.json` లో రాస్తుంది.

`claude remote-control` (`rc`) `--mcp-config` ను తీసుకోదు, దాని సెషన్‌లు సర్వర్లను `.claude.json` నుంచి మాత్రమే చదువుతాయి. `ccam work rc` కోసం, ccam ఆ ఫ్లాగ్‌ను వదిలేసి, టెర్మినల్‌లో అయితే, గ్రూప్ సర్వర్లను `.claude.json` లోకి కాపీ చేయాలా అని అడుగుతుంది, `sync` విధానం లాగే బ్యాకప్‌తో. కాపీ ఏర్పడే వరకు అది అడుగుతుంది, ఆ తర్వాత గ్రూప్ `mcp.json` మారిన తర్వాత లేదా కాపీ చేసిన సర్వర్ పోయిన తర్వాత మాత్రమే మళ్లీ అడుగుతుంది ([Remote Control](/features/launch#remote-control)). పేరు ఘర్షణ అయితే అలాంటి కాపీ `--mcp-config` పై గెలుస్తుంది కాబట్టి, `ccam mcp add`, `add-json`, `remove`, `import` గ్రూప్‌లోని ఇతర ఖాతాలు ఉంచుకున్న కాపీలను అప్‌డేట్ చేస్తాయి:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

కాపీ మారినప్పుడు మాత్రమే `✓ Updated` లైన్ కనిపిస్తుంది. ccam పాత వెర్షన్ సృష్టించిన కాపీ కూడా ఇదే విధంగా తాజాగా చేయబడుతుంది, `ccam migrate` తర్వాతి కాపీ కూడా.

సెషన్ లోపల `claude mcp add` తో జోడించిన సర్వర్ ఆ ఖాతా `.claude.json` లో మాత్రమే చేరుతుంది. సాధారణ `claude mcp add` local scope ను ఉపయోగిస్తుంది: సర్వర్ ఆ ఒక్క ఖాతాకు ఆ ఒక్క ఫోల్డర్‌లో లోడ్ అవుతుంది (`-s user` దాన్ని ప్రతి ఫోల్డర్‌లోనూ ఆ ఖాతాకు అందుబాటులో ఉంచుతుంది). అది జరిగినప్పుడు ccam దాన్ని చూడలేదు, ఎందుకంటే Claude Code అప్పటికే ccam ప్రాసెస్‌ను భర్తీ చేసింది, కాబట్టి అది తదుపరి ప్రారంభంలో తనిఖీ చేస్తుంది. `seen` ఫైల్‌లో `.claude.json` మార్పు సమయం, ఇప్పటికే నివేదించిన సర్వర్ పేర్లు ఉంటాయి. సమయం వేరుగా ఉన్నప్పుడు, ccam పేర్లను గ్రూప్‌తో పోల్చి, ఇన్‌పుట్ కోసం వేచి ఉండకుండా, రెండు విధానాలతోనూ, ప్రతి కొత్త సర్వర్‌కు ఒకసారి సూచనను ప్రింట్ చేస్తుంది:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

గ్రూప్ సర్వర్ పేరే ఉన్న ఖాతా సొంత సర్వర్‌ను కూడా సూచన పేర్కొంటుంది, ఎందుకంటే ఆ ఖాతా సెషన్‌లలో ఖాతాదే గెలుస్తుంది.

`mcp_strategy = "sync"` తో, సర్వర్లు `claude mcp list`, `/mcp` లో సాధారణ యూజర్ సర్వర్లుగా కనిపిస్తాయి, ccam ప్రారంభానికి ముందు వాటిని ప్రతి ఖాతా `.claude.json` లో రాస్తుంది:

<!-- motion:mcp-sync -->

1. ccam గ్రూప్ `mcp.json` మార్పు సమయాన్ని `sync-stamp` తో, ఖాతా `.claude.json` ను తాను చివరిగా చూసిన దానితో పోలుస్తుంది (`mcp.Sync`). ఏదీ మారకపోతే, ఇంకేమీ జరగదు. లేకపోతే ప్రతి గ్రూప్ సర్వర్ ఇంకా `.claude.json` లో ఉందో లేదో కూడా తనిఖీ చేసి, పోయిన కాపీని తిరిగి పెడుతుంది (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. అది Claude Code సొంత లాక్‌ను తీసుకుంటుంది, అంటే ఫైల్ పక్కన `mkdir` తో సృష్టించే `.claude.json.lock` ఫోల్డర్. Claude Code లాక్ లైబ్రరీ లాగే, అది గరిష్ఠంగా 2 సెకన్లు వేచి ఉంటుంది, 10 సెకన్ల కంటే పాత లాక్‌ను పాతబడినదిగా పరిగణిస్తుంది.
3. ప్రస్తుత `.claude.json` ను స్టేట్ ఫోల్డర్‌లోని `backups/` లో సేవ్ చేసి, కొత్తవైన 5 ను ఉంచుతుంది.
4. ప్రతి గ్రూప్ సర్వర్‌ను `mcpServers` లో సెట్ చేస్తుంది, తాను ముందు రాసిన కానీ ఇప్పుడు గ్రూప్‌లో లేని సర్వర్లను తీసివేస్తుంది, ఏదైనా మారినప్పుడు మాత్రమే ఫైల్‌ను దాని అసలు మోడ్‌తో అటామిక్‌గా రాస్తుంది. తాను రాసిన, ఆ తర్వాత ఎవరూ మార్చని కాపీలను మాత్రమే అది మారుస్తుంది లేదా తీసివేస్తుంది. ఒక గ్రూప్ సర్వర్ పేరుతోనే ఖాతాకు సొంత సర్వర్ ఉంటే, ccam ఖాతా సర్వర్‌నే ఉంచి, ఆ విషయం ఒక్కసారి చెబుతుంది: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. ప్రతి కాపీ హాష్‌తో కొత్త స్టాంప్‌ను నమోదు చేస్తుంది, Claude Code `--mcp-config` లేకుండా ప్రారంభమవుతుంది.

Claude Code రన్ అవుతున్నప్పుడు `.claude.json` ను మళ్లీ రాస్తుంది కాబట్టి, ఖాతా సెషన్ రన్ అవుతుంటే, Claude Code ఎలా ప్రారంభమైనా సరే, సింక్ `ccam: account work is already running, skipping MCP sync` తో దాటవేయబడుతుంది. సింక్ విఫలమైతే `ccam: MCP sync failed (launching anyway): …` అని ప్రింట్ చేసి, అయినా ప్రారంభిస్తుంది.

### బయటకు వెళ్లడం, మళ్లీ చేరడం

MCP ను షేర్ చేయడం ఆపే ఖాతా (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) గ్రూప్ సర్వర్లను ఉంచుకుంటుంది: ప్లాన్ వాటిని బ్యాకప్‌తో దాని `.claude.json` లోకి కాపీ చేస్తుంది, అదే పేర్లతో ఉన్న దాని సొంత సర్వర్లు అలాగే ఉంటాయి. `--fresh` తో దానికి ఏదీ రాదు, లోడ్ అవ్వడం ఆగిపోయే సర్వర్లను ప్లాన్ జాబితాగా చూపిస్తుంది. ఆ సర్వర్లు ఉన్న గ్రూప్‌లో అది మళ్లీ చేరినప్పుడు, ఎవరూ మార్చని ఉంచుకున్న కాపీలు తీసివేయబడతాయి, కాబట్టి గ్రూప్ తర్వాతి మార్పులు దానికి చేరతాయి.

`ccam mcp import` ఒక ఖాతా సొంత సర్వర్లను దాని గ్రూప్‌లోకి తరలిస్తుంది:

<!-- motion:mcp-import -->

ccam ఖాతా సర్వర్లను (దాని సొంతవి, ఒకే ఫోల్డర్ కోసం దానికి ఉన్నవి) చదివి, ప్రతిదాన్ని అదే పేరు గల గ్రూప్ సర్వర్‌తో `mcp.SameDef` ద్వారా పోలుస్తుంది, అది కీల క్రమాన్ని, ఖాళీలను పట్టించుకోదు. ఫోల్డర్ సర్వర్ ప్రతి ఫోల్డర్ కోసం గ్రూప్ సర్వర్ అవుతుంది (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); రెండు వేర్వేరు ఫోల్డర్ నిర్వచనాలు ఉన్న ఒకే పేరు తిరస్కరించబడుతుంది. Plugin సర్వర్లను plugin తానే తెస్తుంది కాబట్టి అవి అందించబడవు. గ్రూప్‌లో లేని లేదా అదే నిర్వచనంతో ఉన్న సర్వర్ కాపీ అవుతుంది. వేరుగా ఉన్నది దాటవేయబడుతుంది, ఎందుకంటే మిగతా ప్రతి ఖాతా ఉపయోగించేది గ్రూప్ నిర్వచనమే:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` తో గ్రూప్ సర్వర్ భర్తీ అవుతుంది, మునుపటి `mcp.json` ముందుగా గ్రూప్ ఫోల్డర్‌లోని `.ccam-backup/<time>/mcp.json` లో సేవ్ అవుతుంది. ఆపై ccam అదే లాక్ కింద, స్టేట్ ఫోల్డర్‌లో బ్యాకప్ తర్వాత, కాపీ చేసిన ఎంట్రీలను (దాని సొంతవి, ఫోల్డర్‌వి ఒకే రైట్‌లో) ఖాతా `.claude.json` నుంచి తీసివేస్తుంది (`mcp.Remove`). ఫైల్‌ను చదివినప్పటి నుంచి ఆ ఎంట్రీలు తీసివేయబడే వరకు ccam గ్రూప్ `mcp.json` ను లాక్‌లోనే ఉంచుతుంది, కాబట్టి ఈ మధ్యలో ఎవరూ గ్రూప్ నుంచి సర్వర్‌ను తీసివేయలేరు. పేరు ఘర్షణ అయితే ఖాతా ఎంట్రీ గ్రూప్ ఎంట్రీపై గెలుస్తుంది కాబట్టి ఇది అవసరం. `--keep-local` తో, ఎంట్రీలు ఖాతా `.claude.json` లోనే ఉంటాయి.

## భద్రత, ప్రత్యేక సందర్భాలు

- **సైన్-ఇన్‌లు.** OAuth తో సైన్ ఇన్ అయ్యే సర్వర్లు తమ టోకెన్‌లను సర్వర్ జాబితాలో కాకుండా ప్రతి ఖాతా `.credentials.json` (`mcpOAuth`) లో ఉంచుతాయి, కాబట్టి సర్వర్‌ను షేర్ చేస్తే దాని సైన్-ఇన్ షేర్ కాదు. `/mcp` తో ప్రతి ఖాతాలో ఒక్కసారి సైన్ ఇన్ చేయండి. ccam ఎప్పుడూ సైన్-ఇన్‌లను కాపీ చేయదు, [`ccam migrate`](/features/migrate) కూడా చేయదు: కాపీ చేసిన refresh token మిగతా ప్రతి హోల్డర్‌ను లాగ్ అవుట్ చేస్తుంది. మరో ఖాతా సైన్ ఇన్ అయి ఉండగా ఒక ఖాతా సైన్ ఇన్ కాని సర్వర్లను [`ccam ext`](/features/ext), `ccam doctor` పేర్కొంటాయి, కొత్త ఖాతాకు ఏవి లేవో `ccam next --continue` చెబుతుంది. గ్రూప్‌కు ఎప్పుడూ లాగిన్ ఫైల్ అందదు.
- **రన్ అవుతున్న ఖాతాలు.** `import` గ్రూప్‌కు కాపీ చేస్తుంది కానీ రన్ అవుతున్న ఖాతా సొంత ఎంట్రీలను వదిలేసి ఆ విషయం చెబుతుంది; సెషన్ ముగిసిన తర్వాత దాన్ని మళ్లీ రన్ చేయండి. సింక్ రన్ అవుతున్న ఖాతాలను దాటవేస్తుంది. రెండూ ccam లేకుండా ప్రారంభించిన Claude Code సెషన్‌లను కూడా గుర్తిస్తాయి.
- **Remote Control కోసం కాపీలు.** రన్ అవుతున్న ఖాతా కాపీని ccam ముట్టుకోదు: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` తో ఏదీ అప్‌డేట్ చేయాల్సిన అవసరం లేదు, ఎందుకంటే ప్రతి ప్రారంభం సింక్ చేస్తుంది.
- **పాడైన ఫైళ్లు.** `.claude.json` చెల్లుబాటయ్యే JSON కాకపోతే, సింక్ దాన్ని ముట్టుకోకుండా ఎర్రర్‌ను నివేదిస్తుంది. 2 సెకన్లలో లాక్ తీసుకోలేకపోతే, ccam `could not lock … (Claude Code is writing the file)` అని నివేదిస్తుంది.
- **MCP ను షేర్ చేయని ఖాతాలు.** `import` తిరస్కరించి, ముందుగా ఖాతాను MCP ను షేర్ చేసే గ్రూప్‌లో చేర్చమని చెబుతుంది (`ccam group join <group> side`).
- **గ్రూప్ పేర్లు.** `-g` గ్రూప్ పేరును మాత్రమే తీసుకుంటుంది; అది `groups/` కింద ఒక ఫోల్డర్ అవుతుంది, కాబట్టి `../x` వంటి పేర్లు తిరస్కరించబడతాయి.

## సంబంధిత పేజీలు

- [గ్రూప్‌లతో షేరింగ్](/features/share)
- [ccam తో Remote Control ఉపయోగించడం](/guides/remote-control)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం](/features/migrate)
- [కమాండ్ రిఫరెన్స్](/docs/commands)
