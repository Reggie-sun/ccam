---
title: Claude Code ను ప్రారంభించడం
description: ccam <name|number>, ccam run లేదా పికర్‌తో ఏ ఖాతాగానైనా Claude Code ప్రారంభించండి; ఆర్గ్యుమెంట్‌లు మారకుండా claude కు వెళ్తాయి, --dry-run కమాండ్ చూపుతుంది.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work` Claude Code ను `work` ఖాతాగా ప్రారంభిస్తుంది, పేరు తర్వాత ఉన్నదంతా ఉన్నది ఉన్నట్లుగా `claude` కు వెళ్తుంది. మీరు `claude` అని టైప్ చేసే ప్రతి చోటా దీన్ని ఉపయోగించండి: టెర్మినల్‌లో, స్క్రిప్ట్‌లో, లేదా పైప్ ద్వారా ఇన్‌పుట్ ఇచ్చి.

## వినియోగం

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

ఖాతా పేరును లేదా `ccam list` లోని దాని నంబర్‌ను ఇవ్వండి. పేరు తర్వాతి ఆర్గ్యుమెంట్‌లను ccam చదవదు, కాబట్టి `ccam work --help` Claude Code సహాయాన్ని చూపిస్తుంది. పేరు తర్వాత `--` ఐచ్ఛికం. స్టాండర్డ్ ఇన్‌పుట్, ఎగ్జిట్ కోడ్, సిగ్నల్స్ `claude` తో ఉన్నట్లే పనిచేస్తాయి.

మొదటి ఆర్గ్యుమెంట్ `-` తో మొదలై, ccam సొంత ఫ్లాగ్‌లలో ఒకటి కాకపోతే, ccam అన్ని ఆర్గ్యుమెంట్‌లతో పేరు లేకుండా ఒక ఖాతాను ప్రారంభిస్తుంది. అది [`ccam folder`](/features/folder) తో ప్రస్తుత ఫోల్డర్‌కు సెట్ చేసిన ఖాతా, లేకపోతే డిఫాల్ట్: `ccam use` తో సెట్ చేసిన ఖాతా, మీరు ఏదీ సెట్ చేయకపోతే అతి తక్కువ నంబర్ ఉన్న ఖాతా.

### పికర్

స్టాండర్డ్ ఇన్‌పుట్, అవుట్‌పుట్ టెర్మినల్ అయినప్పుడు ఆర్గ్యుమెంట్‌లు లేని `ccam` పికర్‌ను తెరుస్తుంది, లేకపోతే `ccam list` ను ప్రింట్ చేస్తుంది. ఇది ప్రతి ఖాతా నంబర్, ఇమెయిల్, 5 గంటల వినియోగం, చివరి వినియోగం, షేర్ చేసిన అంశాలను చూపిస్తుంది; ప్రతి ఖాతా సమాధానం ఇచ్చినప్పుడు వినియోగం నిండుతుంది.

| కీ | ఏమి చేస్తుంది |
| --- | --- |
| `Enter` | ఎంచుకున్న ఖాతాగా Claude Code ను ప్రారంభిస్తుంది |
| `1` నుంచి `9` | ఆ నంబర్ ఉన్న ఖాతాను ప్రారంభిస్తుంది |
| `Up`/`Down`, `k`/`j` | ఎంపికను కదిలిస్తుంది |
| `a` | పికర్‌ను మూసి `ccam add` విజార్డ్‌ను రన్ చేస్తుంది |
| `s` | ఎంచుకున్న ఖాతా షేర్ చేసే అంశాలను మార్చుతుంది (`Space` లేదా `x` మార్చుతుంది, `Enter` వర్తింపజేస్తుంది, `Esc` వెనక్కి వెళ్తుంది) |
| `d` | ఎంచుకున్న ఖాతాను డిఫాల్ట్‌గా చేస్తుంది |
| `r` | చిన్న కాష్ లేకుండా వినియోగాన్ని మళ్లీ తెస్తుంది |
| `q`, `Esc`, `Ctrl+C` | నిష్క్రమిస్తుంది |

పికర్ పని చేసే ముందు మూసుకుంటుంది, కాబట్టి Claude Code ప్రారంభమయ్యే ముందు టెర్మినల్ పునరుద్ధరించబడుతుంది. `d` లేదా షేరింగ్ మార్పు తర్వాత అది మళ్లీ తెరుచుకుంటుంది. షేరింగ్ మార్పు `ccam share` లాగే ముందుగా ప్రణాళికను చూపించి అడుగుతుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--dry-run` | ccam రన్ చేసే ఎన్విరాన్‌మెంట్, కమాండ్‌ను ప్రింట్ చేసి నిష్క్రమిస్తుంది | ఆఫ్ |
| `--fallback` | `-p`/`--print` రన్‌ల కోసం: ఖాతా తన వినియోగ పరిమితిని చేరినప్పుడు, అదే కమాండ్‌ను ఎక్కువ అవకాశం ఉన్న లాగిన్ అయిన ఖాతాగా మళ్లీ రన్ చేస్తుంది ([వివరాలు](/features/next#run-again-at-the-limit)) | ఆఫ్ |
| `-v`, `--version` | `ccam <version>` ను ప్రింట్ చేస్తుంది (మొదటి ఆర్గ్యుమెంట్‌గా మాత్రమే) | |
| `-h`, `--help` | `ccam --help` లేదా `ccam run --help` గా, ccam సహాయం | |

ఆర్గ్యుమెంట్ నియమాలు:

- ccam తన సొంత ఫ్లాగ్‌లను ఖాతా పేరుకు ముందు మాత్రమే చదువుతుంది. పేరు తర్వాత, `--dry-run`, `--help` తో సహా ప్రతి ఆర్గ్యుమెంట్ Claude Code కు చెందుతుంది.
- `--dry-run` `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run` (పేరు లేని ఖాతా) గా పనిచేస్తుంది. ఇది పాటించబడుతుంది, తిరస్కరించబడదు. `--dry-run`, `--fallback` ఏ క్రమంలోనైనా రావచ్చు.
- `ccam run <name|number> [claude args]` అనేది `ccam <name>` యొక్క పొడవైన రూపం. పేరు లేని `ccam run` ఫోల్డర్ ఖాతాను లేదా డిఫాల్ట్ ఖాతాను ఉపయోగిస్తుంది.
- ఖాతా పేర్లు `add`, `list`, `run` లేదా `usage` వంటి ccam కమాండ్ పేర్లు కాకూడదు, కాబట్టి `ccam <word>` ఎప్పుడూ అస్పష్టంగా ఉండదు.

డ్రై రన్ మీరు షెల్‌లో, లేదా Windows లో PowerShell లో పేస్ట్ చేయగల ఒక లైన్‌ను ప్రింట్ చేస్తుంది:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## ఇది ఎలా పనిచేస్తుంది

ప్రతి ప్రారంభం కమాండ్ లైన్ నుంచి `exec` వరకు అవే దశలను తీసుకుంటుంది:

<!-- motion:launch -->

ప్రతి బాక్స్ దాని దశ రన్ అయినప్పుడు వెలుగుతుంది, కింద ఉన్న లైన్లు Claude Code ప్రారంభమయ్యే ఎన్విరాన్‌మెంట్, ఆర్గ్యుమెంట్‌లు.

1. `cmd/ccam/main.go` `app.Main` ను పిలుస్తుంది. ccam మరో పేరుతో ప్రారంభించబడితే, అది ఒక [షార్ట్‌కట్](/features/shortcuts), నేరుగా లాంచర్‌కు వెళ్తుంది.
2. cobra జోక్యం చేసుకునే ముందే `app.parse` `os.Args` ను చూస్తుంది. ఒంటరి `ccam` (పికర్), ఒక ccam కమాండ్ పేరు (`config.Reserved`), `-h`/`--help`, `-v`/`--version` మాత్రమే `cli.Execute` కు వెళ్తాయి. మిగతావన్నీ ప్రారంభాలు.
3. `launcher.Launch` `~/.config/ccam/config.toml` ను చదువుతుంది, అది పార్స్ చేసే ఏకైక ఫైల్ ఇదే.
4. `launcher.Prepare` పేరు లేదా నంబర్ ద్వారా ఖాతాను కనుగొని (పేరు లేకపోతే `HereProfile` ఫోల్డర్ ఖాతాను తీసుకుంటుంది, working directory ని `config.toml` లోని `[[folder]]` ఎంట్రీలతో strings గా పోల్చి), ముందున్న ఒక `--` ను తొలగించి, `ResolveClaude` తో `claude` ఎగ్జిక్యూటబుల్‌ను కనుగొంటుంది. PATH లుకప్ ఫైల్ మార్పు సమయంతో పాటు `~/.local/state/ccam/claude-path` లో కాష్ అవుతుంది, కాబట్టి ccam ప్రతి ప్రారంభంలో `PATH` ను వెతకదు. ఫైల్ మారినా లేకపోయినా కొత్త లుకప్ జరుగుతుంది. `config.toml` లోని `claude_bin` పాత్ ఉన్నది ఉన్నట్లుగా ఉపయోగించబడుతుంది. ccam తానే అయిన `claude` ను, inode ద్వారా పోల్చి, ccam తిరస్కరిస్తుంది, ఎందుకంటే అది లూప్ అవుతుంది.
5. `checkLinks` ఖాతా ఫోల్డర్ `~/.local/share/ccam/profiles/<name>/` ను దాని షేరింగ్ సెట్టింగ్‌లతో `share.Check` ద్వారా పోలుస్తుంది, అది `lstat`, `readlink` ను మాత్రమే పిలుస్తుంది, ఒక్కో అంశానికి ఒకటి లేదా రెండు సిస్టమ్ కాల్స్. లేని లింక్ అక్కడికక్కడే సృష్టించబడుతుంది. డేటా కోల్పోయే అవకాశం ఉన్న ఏదైనా ప్రారంభాన్ని ఆపుతుంది.
6. `applyMCP` షేర్ చేసిన MCP సర్వర్లను జోడిస్తుంది (కింద చూడండి).
7. `session.Register` సెషన్‌ను నమోదు చేస్తుంది. ccam తన PID పేరుతో, ప్రాసెస్ ప్రారంభ సమయాన్ని కలిగిన ఒక మార్కర్‌ను `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` లో (లేదా `XDG_RUNTIME_DIR` లేకపోతే `~/.local/state/ccam/locks/` లో) రాస్తుంది. సంభాషణ చరిత్రను షేర్ చేసే ఖాతా కోసం, `historyCheck` ముందుగా `history-<group>-<hash>` లో రన్ అవుతున్న సెషన్ కోసం చూస్తుంది, ఇక్కడ హాష్ ప్రస్తుత ఫోల్డర్‌ది, అక్కడ కూడా నమోదు చేస్తుంది. `exec` PID ను ఉంచుతుంది, కాబట్టి ఆ తర్వాత మార్కర్ Claude Code ప్రాసెస్‌నే సూచిస్తుంది. పట్టుకున్న లాక్ కాకుండా మార్కర్ ఎందుకంటే, Claude Code ప్రారంభించే ప్రాసెస్‌లు లాక్‌ను వారసత్వంగా పొంది సెషన్ ముగిసిన తర్వాత కూడా దాన్ని ఉంచుకుంటాయి.
8. `Plan.Exec` టెర్మినల్ టైటిల్‌ను `claude · <name>` గా సెట్ చేస్తుంది (స్టాండర్డ్ ఎర్రర్ టెర్మినల్ అయినప్పుడు), `fsync` లేకుండా ప్రారంభ సమయాన్ని `~/.local/state/ccam/last-used/<name>` లో రాస్తుంది, `syscall.Exec` ను పిలుస్తుంది.

Claude Code `CLAUDE_CONFIG_DIR` ఖాతా ఫోల్డర్‌కు, `CCAM_PROFILE` ఖాతా పేరుకు సెట్ చేయబడి రన్ అవుతుంది. ccam వారసత్వంగా వచ్చిన వీటి విలువలను తొలగిస్తుంది. macOS లో `ccam import` తో నమోదు చేసిన `~/.claude` ఫోల్డర్ కోసం, ఇది `CLAUDE_SECURESTORAGE_CONFIG_DIR` ను కూడా ఖాళీ విలువకు సెట్ చేస్తుంది, తద్వారా ఆ ఖాతా, సాధారణ `claude` ఒకే Keychain లాగిన్‌ను ఉపయోగిస్తాయి.

### exec ఎందుకు

`syscall.Exec` ccam ప్రాసెస్‌ను `claude` తో భర్తీ చేస్తుంది, కాబట్టి ccam మెమరీలో ఉండదు, కీస్ట్రోక్‌లు, సిగ్నల్స్, ఎగ్జిట్ కోడ్ Claude Code సొంతమే, ఆర్గ్యుమెంట్‌లు ఎప్పుడూ షెల్ స్ట్రింగ్ కాకుండా అర్రేగానే ఉంటాయి.

Windows లో `exec` లేదు. అక్కడ `Plan.Exec` అదే ఎన్విరాన్‌మెంట్, ఆర్గ్యుమెంట్‌లు, console తో `claude.exe` ను చైల్డ్ ప్రాసెస్‌గా ప్రారంభించి, దాని కోసం వేచి ఉండి, దాని ఎగ్జిట్ కోడ్‌తో నిష్క్రమిస్తుంది. Ctrl+C console నుంచి నేరుగా Claude Code కు చేరుతుంది, కాబట్టి ccam దాన్ని పట్టించుకోదు. సెషన్ మార్కర్‌లో ccam PID ఉంటుంది, అది సరిగ్గా Claude Code ఉన్నంత కాలం ఉంటుంది, Claude Code నిష్క్రమించినప్పుడు మార్కర్ తొలగించబడుతుంది. వేచి ఉన్నప్పుడు ccam ఇంకేమీ చేయదు. `claude mcp add` తో జోడించిన సర్వర్లను గమనించడం వంటి సెషన్ తర్వాతి పని, తదుపరి ప్రారంభంలో జరుగుతుంది: ccam `.claude.json` మార్పు సమయాన్ని `~/.local/state/ccam/mcp/<name>/seen` తో పోల్చి, ఫైల్ మారినప్పుడు మాత్రమే దాన్ని పార్స్ చేస్తుంది.

### కమాండ్ లైన్ ఎలా దారి మళ్లించబడుతుంది

మూడు కమాండ్ లైన్లు, `internal/app` ద్వారా మూడు మార్గాలు:

<!-- motion:launch-route -->

షార్ట్‌కట్ `argv[0]` ద్వారా, `ccam <name>` మొదటి ఆర్గ్యుమెంట్ ద్వారా గుర్తించబడతాయి, రెండూ `syscall.Exec` తో ముగుస్తాయి. `ccam usage` వంటి ccam కమాండ్ మాత్రమే cobra ను చేరుతుంది.

### షేర్ చేసిన MCP సర్వర్లు

డిఫాల్ట్ `flag` విధానంతో, ccam గ్రూప్ MCP సర్వర్లను మీ ఆర్గ్యుమెంట్‌ల ముందు, ఒక ఆర్గ్యుమెంట్‌గా Claude Code కు పంపుతుంది:

<!-- motion:launch-mcp -->

గ్రూప్ `mcp.json` ఉంటేనే `mcp.FlagArgs` ఒక ఆర్గ్యుమెంట్‌ను ఇస్తుంది, `Plan.InsertArgs` దాన్ని `claude` తర్వాత వెంటనే పెడుతుంది. ఆపై మీ సొంత `--mcp-config` లేదా `--strict-mcp-config` దానిపై అదనంగా వర్తిస్తుంది. ఇది ఒకే ఆర్గ్యుమెంట్, `--mcp-config=<path>`, ఎందుకంటే Claude Code యొక్క `--mcp-config` అనేక విలువలను తీసుకుంటుంది: రెండు ఆర్గ్యుమెంట్‌లుగా రాస్తే, తర్వాత వచ్చే ప్రాంప్ట్‌ను కూడా కాన్ఫిగ్ ఫైల్‌గా చదువుతుంది.

`config.toml` లో `mcp_strategy = "sync"` తో, ccam బదులుగా సర్వర్లను ఖాతా `.claude.json` లో రాస్తుంది, ఖాతా ఇప్పటికే రన్ అవుతుంటే దాన్ని దాటవేస్తుంది.

### Remote Control

`ccam work rc` (లేదా `ccam work remote-control`) Claude Code యొక్క Remote Control ను `work` గా ప్రారంభిస్తుంది. `claude remote-control` తన వెర్బ్‌కు ముందు `--mcp-config` ను తిరస్కరిస్తుంది, దాని తర్వాత ఏదీ తీసుకోదు, దాని సెషన్‌లు MCP సర్వర్లను ఖాతా `.claude.json` నుంచి మాత్రమే చదువుతాయి. కాబట్టి ఇక్కడ ఫ్లాగ్ విధానం వేరుగా పనిచేస్తుంది:

<!-- motion:launch-rc -->

1. `remoteControl` ఖాతా తర్వాత `-` తో ప్రారంభం కాని మొదటి ఆర్గ్యుమెంట్‌ను Claude Code వెర్బ్‌గా తీసుకుంటుంది. అది `remote-control` లేదా `rc` అయితే, `mcp.FlagArgs` ఏమీ జోడించదు.
2. `mcp.Synced` గ్రూప్ `mcp.json` మార్పు సమయాన్ని ఖాతా `sync-stamp` తో పోలుస్తుంది. `.claude.json` లో ఇప్పటికే గ్రూప్ ప్రస్తుత సర్వర్లు ఉంటే, లేదా గ్రూప్‌కు సర్వర్లే లేకపోతే, Remote Control ప్రశ్న లేకుండా ప్రారంభమవుతుంది.
3. లేకపోతే ccam సెషన్‌లకు అందని సర్వర్ల పేర్లు చెప్పి అడుగుతుంది:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` అయితే, `mcp.Sync` వాటిని `sync` విధానం లాగే కాపీ చేస్తుంది: Claude Code లాక్ కింద, బ్యాకప్ తర్వాత, ఖాతా సొంత సర్వర్లను ఉంచుతూ. `N` అయితే, లేదా టెర్మినల్ లేకపోతే, Remote Control అవి లేకుండా ప్రారంభమవుతుంది, ccam ప్రతి ప్రారంభానికి ముందు వాటిని కాపీ చేసే `mcp_strategy = "sync"` ను సూచిస్తుంది.

సర్వర్లు కాపీ అయ్యే వరకు ప్రతి Remote Control ప్రారంభంలో ccam అడుగుతుంది; కాపీ తర్వాత, గ్రూప్ `mcp.json` మళ్లీ మారినప్పుడు మాత్రమే అడుగుతుంది. తర్వాతి సాధారణ ప్రారంభాల్లో `.claude.json` లోని కాపీ `--mcp-config` పై గెలుస్తుంది, కాబట్టి `ccam mcp add`, `remove`, `import` ఇతర ఖాతాలు ఉంచుకున్న కాపీలను అప్‌డేట్ చేస్తాయి ([షేర్ చేసిన MCP సర్వర్లు](/features/mcp)).

### 10 ms బడ్జెట్

Claude Code ప్రారంభమయ్యే ముందు ccam గరిష్ఠంగా 10 ms (p95) జోడిస్తుంది. ప్రారంభ మార్గం ఎప్పుడూ cobra కమాండ్ ట్రీని నిర్మించదు, Bubble Tea ను ప్రారంభించదు, నెట్‌వర్క్ కనెక్షన్‌ను తెరవదు, ప్యాకేజీ `init()` ఫంక్షన్‌లు మొత్తం 5 ms కంటే ఎక్కువ తీసుకుంటే `TestPackageInitIsCheap` విఫలమవుతుంది. `scripts/perf-check.sh` `hyperfine` తో `ccam --dry-run work`, `ccam work`, ఒక షార్ట్‌కట్‌ను కొలిచి 10 ms దాటితే విఫలమవుతుంది, `internal/launcher` లోని `BenchmarkLaunchPath` దశలను ఒక లూప్‌లో రన్ చేస్తుంది. ఒక్కో దశ సమయాన్ని చూడటానికి, `CCAM_TRACE=1` ను సెట్ చేయండి:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## భద్రత, ప్రత్యేక సందర్భాలు

- తెలియని పేరు లేదా నంబర్ `no account "<name>"` తో ఆగి `ccam list` ను సూచిస్తుంది. అసలు ఖాతాలే లేకపోతే, ccam `ccam add <name>` తో ఒకటి జోడించమని అడుగుతుంది.
- ఒక షేర్డ్ పాత్ సాధారణ ఫైల్‌గా మారితే లేదా తప్పు చోటికి చూపిస్తే, ccam `stopped because of share link problems:` తో ఆగి, వాటిని జాబితా చేసి, `ccam doctor --fix` రన్ చేయమని అడుగుతుంది. ప్రారంభ మార్గంలో అది డేటాను తరలించదు.
- ఈ ఖాతాతో సంభాషణ చరిత్రను షేర్ చేసే (అదే గ్రూప్) మరో Claude Code అదే ఫోల్డర్‌లో రన్ అవుతుంటే, చరిత్ర కలిసిపోవచ్చని ccam హెచ్చరించి `Run anyway? [y/N]` అని అడుగుతుంది. టెర్మినల్ లేకపోతే అది హెచ్చరికను చూపించి కొనసాగుతుంది.
- `--dry-run` MCP సింక్, కొత్త సర్వర్ నోటీసు, చరిత్ర తనిఖీ, సెషన్ మార్కర్, టైటిల్, చివరి వినియోగ సమయాన్ని దాటవేస్తుంది, exec చేయదు. అయినా అది షేర్ లింక్‌లను తనిఖీ చేస్తుంది, సాధారణ ప్రారంభం లాగే లేని లింక్‌ను సృష్టిస్తుంది.
- `ccam work rc`, `work` ఇప్పటికే రన్ అవుతున్నప్పుడు, `.claude.json` ను ముట్టుకోకుండా, ఆ విషయం చెప్పి, గ్రూప్ సర్వర్లు లేకుండా Remote Control ను ప్రారంభిస్తుంది; వాటిని కాపీ చేయడానికి సెషన్ ముగిసిన తర్వాత మళ్లీ రన్ చేయండి. `--dry-run` తో అది సందేశాన్ని ప్రింట్ చేస్తుంది కానీ అడగదు, కాపీ చేయదు.
- `mcp_strategy = "sync"` తో, సింక్ విఫలమైతే `ccam: MCP sync failed (launching anyway): …` అని ప్రింట్ అవుతుంది, Claude Code తనకు ఇప్పటికే ఉన్న సర్వర్లతో ప్రారంభమవుతుంది.

## సంబంధిత పేజీలు

- [షార్ట్‌కట్ కమాండ్‌లు](/features/shortcuts)
- [ccam తో Remote Control ఉపయోగించడం](/guides/remote-control)
- [ఖాతాలను నిర్వహించడం](/features/accounts)
- [కమాండ్ రిఫరెన్స్](/docs/commands)
- [MCP సర్వర్లు, సెట్టింగ్‌లు, చరిత్రను షేర్ చేయండి](/docs/sharing)
- [CLAUDE_CONFIG_DIR ఎలా పనిచేస్తుంది](/guides/claude-config-dir)
