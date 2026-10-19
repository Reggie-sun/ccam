---
title: ccam తో Remote Control ఉపయోగించడం
description: ccam <name> rc తో ఏ ఖాతాగానైనా Claude Code Remote Control ను ప్రారంభించండి, ఏ ఆప్షన్లు వర్తిస్తాయో చూడండి, షేర్ చేసిన MCP సర్వర్లను కాపీ చేయాలా అని ccam ఎందుకు అడుగుతుందో తెలుసుకోండి, సాధారణ సమస్యలను పరిష్కరించండి.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

`ccam work rc` Claude Code యొక్క Remote Control ను `work` ఖాతాగా ప్రారంభిస్తుంది, కాబట్టి ఈ కంప్యూటర్‌లోని సెషన్‌లను claude.ai/code లేదా Claude మొబైల్ యాప్ నుంచి నడపవచ్చు. ఇది `claude rc` లాగే పనిచేస్తుంది, ఒకే ఒక్క తేడాతో: ఖాతా MCP సర్వర్లను షేర్ చేస్తే, వాటిని ఖాతా `.claude.json` లోకి కాపీ చేయాలా అని ccam ఒక్కసారి అడుగుతుంది, ఎందుకంటే Remote Control సెషన్‌లు సర్వర్లను మరెక్కడి నుంచీ చదవవు.

## ప్రారంభించడం

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[షార్ట్‌కట్](/features/shortcuts) కూడా ఇలాగే పనిచేస్తుంది: `ccwork rc` అంటే `ccam work rc`. మీరు పని చేయాలనుకునే ఫోల్డర్‌లో కమాండ్‌ను రన్ చేయండి, ఆ తర్వాత `work` ఉపయోగించే అదే Claude ఖాతాతో సైన్ ఇన్ అయిన claude.ai/code లేదా Claude యాప్‌ను తెరవండి.

## ఆప్షన్లు

ఇక్కడ రెండు రకాల ఆప్షన్లు కలుస్తాయి: ccam సొంత ఆప్షన్లు, ఇవి ఖాతా పేరుకు ముందు వస్తాయి లేదా ccam సెట్టింగ్‌లలో ఉంటాయి; Remote Control ఆప్షన్లు, ఇవి `rc` తర్వాత వచ్చి మార్పు లేకుండా Claude Code కు చేరతాయి.

### ccam ఆప్షన్లు, సెట్టింగ్‌లు

| ఆప్షన్ లేదా సెట్టింగ్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--dry-run` (ఖాతాకు ముందు) | ఎన్విరాన్‌మెంట్, కమాండ్, MCP నోటీసు ఉంటే అదీ ప్రింట్ చేస్తుంది. అడగదు, కాపీ చేయదు | ఆఫ్ |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | ఖాతా MCP సర్వర్లను షేర్ చేస్తూ, దాని `.claude.json` లో గ్రూప్ ప్రస్తుత జాబితా లేనప్పుడు అడుగుతుంది. `y` వాటిని కాపీ చేస్తుంది, `N` లేదా Enter అవి లేకుండా ప్రారంభిస్తుంది | `N` |
| `~/.config/ccam/config.toml` లో `mcp_strategy = "sync"` | Remote Control అయినా కాకపోయినా, ప్రతి ప్రారంభానికి ముందు గ్రూప్ సర్వర్లను `.claude.json` లోకి కాపీ చేస్తుంది, కాబట్టి ఎప్పుడూ ఏమీ అడగదు | `"flag"` |
| `ccam share work --mcp off` | ఖాతాతో MCP సర్వర్లను షేర్ చేయడం ఆపుతుంది; అప్పుడు Remote Control ఖాతా సొంత సర్వర్లతో మాత్రమే ప్రారంభమవుతుంది, ఏమీ అడగదు | ఆన్ |

### Remote Control ఆప్షన్లు

`rc` తర్వాత ఉన్నదంతా Claude Code కు చెందినది. ఇవి Claude Code 2.1.304 లోని `claude remote-control` ఆప్షన్లు; మీ వెర్షన్ జాబితాను `ccam work rc --help` ప్రింట్ చేస్తుంది.

| ఆప్షన్ | ఏమి చేస్తుంది |
| --- | --- |
| `--name <name>` | claude.ai/code లో కనిపించే సెషన్ పేరు |
| `--remote-control-session-name-prefix <prefix>` | ఆటోమేటిక్‌గా సృష్టించే సెషన్ పేర్లకు ప్రిఫిక్స్ (డిఫాల్ట్: హోస్ట్ పేరు) |
| `-c`, `--continue` | ఈ ఫోల్డర్ కోసం Remote Control చివరిగా నమోదు చేసిన సెషన్‌కు, సుమారు 4 గంటలలోపు, మళ్లీ కనెక్ట్ అవుతుంది |
| `--session-id <id>` | నిర్దిష్ట సెషన్‌కు మళ్లీ కనెక్ట్ అవుతుంది |
| `--permission-mode <mode>` | కొత్త సెషన్‌లకు అనుమతి మోడ్: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` లేదా `plan` |
| `--spawn <mode>` | `same-dir` (డిఫాల్ట్), `worktree` (ప్రతి సెషన్‌కు ఒక git worktree) లేదా `session` (ఒకే సెషన్, అది ముగిసినప్పుడు నిష్క్రమిస్తుంది) |
| `--capacity <N>` | `same-dir`, `worktree` మోడ్‌లలో ఒకేసారి గరిష్ఠ సెషన్‌లు (డిఫాల్ట్ 32) |
| `--[no-]create-session-in-dir` | ప్రారంభంలోనే ప్రస్తుత ఫోల్డర్‌లో మొదటి సెషన్‌ను సృష్టిస్తుంది (డిఫాల్ట్ ఆన్) |
| `--[no-]chrome` | కొత్త సెషన్‌లకు Claude in Chrome (డిఫాల్ట్ ఆఫ్) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | డీబగ్, వివరణాత్మక అవుట్‌పుట్ |

`rc` ను ఖాతా పేరు తర్వాత వెంటనే ఉంచండి. `-` తో ప్రారంభం కాని మొదటి ఆర్గ్యుమెంట్‌ను ccam Claude Code కమాండ్‌గా తీసుకుంటుంది, కాబట్టి `ccam work rc --name laptop` లో అది `rc` ను చూస్తుంది, కానీ `ccam work --model opus rc` లో `opus` ను చూస్తుంది.

## ccam ఏమి చేస్తుంది

ప్రారంభం సాధారణ మార్గంలోనే జరుగుతుంది ([Claude Code ను ప్రారంభించడం](/features/launch)), MCP దశలో మాత్రం ఒక మార్పుతో:

1. ఏ ప్రారంభంలోనైనా లాగే, ccam ఖాతాను కనుగొని, దాని షేర్ లింక్‌లను తనిఖీ చేసి, `claude` ను గుర్తిస్తుంది.
2. కమాండ్ `rc` లేదా `remote-control` అని గమనించి `--mcp-config` ను వదిలేస్తుంది.
3. ఖాతా MCP సర్వర్లను షేర్ చేస్తే, గ్రూప్ `mcp.json` ను `~/.local/state/ccam/mcp/work/` లోని `sync-stamp` తో పోలుస్తుంది. `.claude.json` లో ఇప్పటికే ప్రస్తుత సర్వర్లు ఉంటే, లేదా గ్రూప్‌కు సర్వర్లే లేకపోతే, నేరుగా దశ 6 కు వెళ్తుంది.
4. లేకపోతే సెషన్‌లకు అందని సర్వర్ల పేర్లు చెబుతుంది:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y` అయితే, `.claude.json` బ్యాకప్‌ను సేవ్ చేసి, Claude Code సొంత లాక్‌ను తీసుకుని, గ్రూప్ సర్వర్లను `mcpServers` కు జోడించి, ఖాతా ఇతర సర్వర్లను అలాగే ఉంచుతుంది. `N` అయితే, టెర్మినల్ లేకపోతే, లేదా ఖాతా రన్ అవుతుంటే, ఏదీ మార్చకుండా ఆ విషయం చెబుతుంది:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. సెషన్‌ను నమోదు చేసి, `CLAUDE_CONFIG_DIR` ను ఖాతా ఫోల్డర్‌కు సెట్ చేసి, `exec` ద్వారా తనను తాను `claude rc` తో భర్తీ చేసుకుంటుంది.

ఆ తర్వాత Remote Control సర్వర్‌గా రన్ అవుతూనే ఉంటుంది. అది మీ కోసం ప్రారంభించే ప్రతి సెషన్ ఆ ప్రాసెస్ యొక్క చైల్డ్, కాబట్టి అది `work` గా రన్ అవుతుంది: `work` లాగిన్, సెట్టింగ్‌లు, సంభాషణ చరిత్రతో, `work` వినియోగ పరిమితిలో లెక్కించబడుతుంది. సర్వర్ రన్ అవుతున్నంత కాలం ccam ఆ ఖాతాను రన్ అవుతున్నట్లుగా పరిగణిస్తుంది.

డ్రై రన్ ఈ తేడాను చూపిస్తుంది. సాధారణ ప్రారంభం గ్రూప్ సర్వర్లను ఆర్గ్యుమెంట్‌గా తీసుకువెళ్తుంది, Remote Control తీసుకువెళ్లదు:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## ఇది ఇలా ఎందుకు పనిచేస్తుంది

**Remote Control `--mcp-config` ను తీసుకోదు.** సాధారణ ప్రారంభాల్లో ccam గ్రూప్ సర్వర్లను Claude Code కు `--mcp-config=<path>` గా అందిస్తుంది, Claude Code ఫైళ్లలో ఎప్పుడూ రాయదు. `claude remote-control` తన కమాండ్ పేరుకు ముందు ఆ ఆప్షన్‌ను తిరస్కరిస్తుంది, తర్వాత ఏదీ అంగీకరించదు; కాబట్టి ccam 0.6.0 కి ముందు, MCP సర్వర్లను షేర్ చేసిన ప్రతి ఖాతాకు `ccam work rc` విఫలమయ్యేది. ఇప్పుడు ccam ఆ ఆప్షన్‌ను వదిలేస్తుంది.

**దాని సెషన్‌లు సర్వర్లను `.claude.json` నుంచి మాత్రమే చదువుతాయి.** Remote Control ప్రారంభించే సెషన్‌లకు కమాండ్ లైన్ ఆప్షన్లు వారసత్వంగా రావు. యూజర్ MCP సర్వర్లు వాటికి దొరికే ఏకైక చోటు ఖాతా `.claude.json`. ఆ ఫైల్‌లో ఖాతా లాగిన్ వివరాలు, ఫోల్డర్ వారీ స్థితి కూడా ఉంటాయి, కాబట్టి `settings.json` కు చేసినట్లు ccam దాన్ని షేర్ చేసిన కాపీకి లింక్ చేయలేదు. సెషన్‌లకు షేర్ చేసిన జాబితాను ఇవ్వడానికి సర్వర్లను అందులోకి కాపీ చేయడం ఒక్కటే మార్గం.

**రాసే ముందు ccam అడుగుతుంది.** డిఫాల్ట్ ఫ్లాగ్ విధానంలో, `.claude.json` ను ముట్టుకోనని ccam హామీ ఇస్తుంది. అక్కడ తప్పనిసరిగా రాయాల్సిన ఏకైక సందర్భం Remote Control, కాబట్టి అది ముందుగా అడుగుతుంది, బ్యాకప్ ఉంచుతుంది, `mcpServers` తప్ప మరేదీ ముట్టుకోదు.

**ఖాతా రన్ అవుతున్నప్పుడు అది రాయదు.** సెషన్ జరుగుతున్నప్పుడు Claude Code `.claude.json` ను మళ్లీ రాస్తుంది, బయటి నుంచి చేసిన మార్పు పోతుంది లేదా దానిపై మళ్లీ రాయబడుతుంది. అందుకే మరో Remote Control సర్వర్‌తో సహా, ఖాతా యొక్క ఏ సెషన్ రన్ అవుతున్నా ccam కాపీని దాటవేస్తుంది.

**ప్రతి మార్పుకు అది ఒక్కసారే అడుగుతుంది.** కాపీ తర్వాత `sync-stamp` గ్రూప్ `mcp.json` మార్పు సమయాన్ని ఉంచుకుంటుంది. ఆ ఫైల్ మారినప్పుడు మాత్రమే, ఉదాహరణకు `ccam mcp add` తర్వాత, ccam మళ్లీ అడుగుతుంది. మీరు `N` అని సమాధానం ఇస్తే ఇంకా కాపీ ఉండదు, కాబట్టి ప్రతి Remote Control ప్రారంభంలో అది అడుగుతుంది.

**`ccam mcp` కాపీలను తాజాగా ఉంచుతుంది.** పేరు ఘర్షణ అయితే `.claude.json` లోని ఎంట్రీ `--mcp-config` పై గెలుస్తుంది, కాబట్టి పాత కాపీ సాధారణ ప్రారంభాల్లో కూడా కొత్త నిర్వచనాన్ని దాచేస్తుంది. అందుకే `ccam mcp add`, `add-json`, `remove`, `import` రన్ అవ్వని ప్రతి ఖాతా కాపీని అప్‌డేట్ చేస్తాయి:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**టెర్మినల్ లేకపోతే సమాధానం "కాదు".** Remote Control తరచుగా స్క్రిప్ట్, సర్వీస్ లేదా `nohup` నుంచి రన్ అవుతుంది, అక్కడ సమాధానం ఇవ్వడానికి ఎవరూ ఉండరు. అప్పుడు ccam వేచి ఉండకుండా సర్వర్లు లేకుండా ప్రారంభిస్తుంది.

## సమస్యలు, పరిష్కారాలు

### Remote Control సెషన్‌లలో షేర్ చేసిన MCP సర్వర్లు లేవు

ప్రశ్నకు `N` అని సమాధానం ఇచ్చారు, టెర్మినల్ లేదు, లేదా ఖాతా రన్ అవుతోంది. Remote Control ను ఆపి, ఖాతా ఇతర సెషన్‌లను మూసి, టెర్మినల్‌లో `ccam work rc` ను మళ్లీ రన్ చేయండి; `y` అని సమాధానం ఇవ్వండి. ఎప్పుడూ అడగకుండా ఉండాలంటే, `~/.config/ccam/config.toml` లో `mcp_strategy = "sync"` సెట్ చేయండి. సెషన్‌లో ఇప్పటికే ఉన్న సర్వర్లు అది రన్ అవుతున్నప్పుడు మారవు, కాబట్టి కాపీ తర్వాత కొత్త సెషన్ ప్రారంభించండి.

### "Account work is running, so its .claude.json is left alone"

`work` యొక్క ఒక సెషన్ ఇంకా తెరిచి ఉంది: సాధారణ `ccam work`, ఒక షార్ట్‌కట్, లేదా మరో Remote Control సర్వర్. దాని నుంచి నిష్క్రమించి `ccam work rc` ను మళ్లీ రన్ చేయండి. వేరే ఖాతా నడుపుతున్న సెషన్‌లు లెక్కలోకి రావు.

### ccam ప్రతిసారీ అడుగుతుంది

మీరు ఇంకా సర్వర్లను కాపీ చేయలేదు. ఒక్కసారి `y` అని సమాధానం ఇవ్వండి, లేదా `mcp_strategy = "sync"` సెట్ చేయండి. Remote Control లో సర్వర్లు ఎప్పటికీ వద్దనుకుంటే, ఖాతాతో MCP షేరింగ్‌ను ఆపండి: `ccam share work --mcp off`.

### `ccam work rc` వెంటనే `--mcp-config` గురించిన ఎర్రర్‌తో విఫలమవుతుంది

మీ ccam 0.6.0 కంటే పాతది, అది Remote Control కు `--mcp-config` ను పంపుతుంది. `ccam update` రన్ చేయండి. ప్రస్తుత ccam తో కూడా ఇదే ఎర్రర్ వస్తే, విలువ తీసుకునే ఏదో ఆప్షన్ `rc` కంటే ముందు వచ్చిందని అర్థం, `ccam work --model opus rc` లో లాగా; అందుకే ccam `opus` ను కమాండ్‌గా తీసుకుంది. `rc` ను ఖాతా పేరు తర్వాత వెంటనే ఉంచండి.

### `ccam mcp` తో చేసిన మార్పు Remote Control కు చేరడం లేదు

`ccam mcp add`, `remove` రన్ అవ్వని ఖాతాల్లో మాత్రమే కాపీలను అప్‌డేట్ చేస్తాయి, ఆ విషయం చెబుతాయి కూడా: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control ను ఆపి మళ్లీ ప్రారంభించండి; ccam మారిన `mcp.json` ను గమనించి మళ్లీ అడుగుతుంది.

### కాపీ నా సొంత సర్వర్‌ను భర్తీ చేసింది

గ్రూప్ సర్వర్ పేరుతోనే ఖాతాకు సొంత సర్వర్ ఉంటే, కాపీ దాన్ని భర్తీ చేసింది. మునుపటి `.claude.json` `~/.local/state/ccam/mcp/work/backups/` లో ఉంది, అక్కడ సరికొత్త ఐదు బ్యాకప్‌లు ఉంటాయి. మీ వెర్షన్‌ను అన్ని ఖాతాలకు ఉంచాలంటే, తర్వాతి Remote Control ప్రారంభానికి ముందు `ccam mcp import work <server> --force` తో దాన్ని గ్రూప్‌లోకి తరలించండి.

### కాపీ అయిన సర్వర్లను తొలగించాలి

`ccam mcp remove <server>` సర్వర్‌ను గ్రూప్ నుంచి, ప్రతి కాపీ నుంచి తొలగిస్తుంది. కాపీకి ముందు స్థితికి తిరిగి వెళ్లాలంటే, ఖాతా యొక్క ఏ సెషన్ రన్ అవ్వనప్పుడు `~/.local/state/ccam/mcp/work/backups/` నుంచి సరికొత్త బ్యాకప్‌ను పునరుద్ధరించండి.

### Remote Control ఫోల్డర్‌ను నమ్మాలా అని అడుగుతుంది, లేదా టెర్మినల్ లేకుండా నిష్క్రమిస్తుంది

ఫోల్డర్ నమ్మకం ప్రతి ఖాతాకు విడిగా `.claude.json` లో నిల్వ అవుతుంది, కాబట్టి `personal` గా మీరు నమ్మిన ఫోల్డర్ `work` కు కొత్తది. ఆ ఫోల్డర్‌లో ఒక్కసారి `ccam work` రన్ చేసి అంగీకరించండి, ఆ తర్వాత అక్కడ Remote Control ప్రారంభించండి.

### "You must be logged in with a Claude account that has a subscription"

Remote Control కు Claude సబ్‌స్క్రిప్షన్ ఉన్న లాగిన్ అవసరం. `ccam list` లో ఖాతా ఈమెయిల్‌ను తనిఖీ చేయండి, అవసరమైతే `ccam login work` తో మళ్లీ సైన్ ఇన్ చేయండి.

### సెషన్ ఫోన్‌లో కనిపించడం లేదు

Claude యాప్ లేదా claude.ai/code వేరే Claude ఖాతాతో సైన్ ఇన్ అయి ఉంది. `work` కోసం `ccam list` చూపించే ఈమెయిల్‌తో దాన్ని పోల్చండి.

### ఈ ఫోల్డర్‌లో మరో Claude Code సంభాషణ చరిత్రను షేర్ చేస్తోందని ccam హెచ్చరిస్తుంది

అదే చరిత్ర గ్రూప్‌లోని మరో ఖాతా అదే ఫోల్డర్‌లో రన్ అవుతోంది. రెండు ఖాతాల నుంచి ఒకే సంభాషణను కొనసాగించే సెషన్‌లు చరిత్రను కలగలపవచ్చు. అవి వేర్వేరు సంభాషణలపై పనిచేస్తాయి అనుకుంటేనే `y` అని సమాధానం ఇవ్వండి, లేదా Remote Control ను మరో ఫోల్డర్‌లో ప్రారంభించండి.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [షేర్ చేసిన MCP సర్వర్లు](/features/mcp)
- [ఖాతాల మధ్య షేరింగ్](/features/share)
- [ఖాతాలను నిర్వహించడం](/features/accounts)
