---
title: ప్రతి ఖాతా బ్యాక్‌గ్రౌండ్ సెషన్‌లు
description: ccam agents మీ అన్ని ఖాతాల Claude Code బ్యాక్‌గ్రౌండ్ ఏజెంట్‌లను, తెరిచి ఉన్న సెషన్‌లను ఒకే టేబుల్‌లో చూపిస్తుంది, ccam attach ఒకదాన్ని దాని id తో తెరుస్తుంది.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code బ్యాక్‌గ్రౌండ్ సెషన్‌లను (`claude --bg`, agent view) ఖాతా వారీగా ఉంచుతుంది: `claude agents` తాను రన్ అవుతున్న ఖాతా సెషన్‌లను మాత్రమే చూపిస్తుంది. `ccam agents` ప్రతి ఖాతాను ఒకేసారి అడిగి వాటన్నింటినీ కలిపి జాబితా చేస్తుంది, `ccam attach <id>` ఒక బ్యాక్‌గ్రౌండ్ సెషన్‌ను దాన్ని ప్రారంభించిన ఖాతాగా తెరుస్తుంది.

## వినియోగం

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` అనేది బ్యాక్‌గ్రౌండ్ సెషన్ స్థితి (`working`, `blocked`, `done`, `failed`, `stopped`), దాని ప్రాసెస్ ఏమి చేస్తోందో (`busy`, `waiting`, `idle`); తెరిచి ఉన్న ఇంటరాక్టివ్ సెషన్ `open` అని చూపిస్తుంది.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--all` | పూర్తయిన, విఫలమైన లేదా ఆగిపోయిన బ్యాక్‌గ్రౌండ్ సెషన్‌లను కూడా జాబితా చేస్తుంది | ఆఫ్ |
| `--json` | సెషన్‌లను JSON గా ప్రింట్ చేస్తుంది, ప్రతిదానికి అది చెందిన `account` తో | ఆఫ్ |

`ccam attach <id> [claude args]` id తర్వాతి ఆర్గ్యుమెంట్‌లను `claude attach` కు పంపుతుంది. `--dry-run` తో అది కమాండ్‌ను ప్రింట్ చేస్తుంది.

## ఇది ఎలా పనిచేస్తుంది

1. ప్రతి ఖాతా కోసం, ccam ఖాతా `CLAUDE_CONFIG_DIR` తో `claude agents --json` ను (ఇస్తే `--all` తో) రన్ చేస్తుంది, అన్నీ ఒకేసారి, ప్రతిదీ 15 సెకన్లలోపు. ఇది ఖాతా సొంత సెషన్ జాబితాను చదువుతుంది, Claude Code బ్యాక్‌గ్రౌండ్ సర్వీస్‌ను ప్రారంభించదు.
2. ప్రతి ఎంట్రీకి ఖాతా పేరు జోడించబడుతుంది, జాబితా ప్రారంభ సమయం ప్రకారం, కొత్తవి ముందుగా క్రమబద్ధీకరించబడుతుంది. సమాధానం ఇవ్వని ఖాతా పేరు స్టాండర్డ్ ఎర్రర్‌పై చెప్పబడి, జాబితా నుంచి వదిలేయబడుతుంది, ఆ తర్వాత ccam కోడ్ 1 తో నిష్క్రమిస్తుంది, కాబట్టి స్క్రిప్ట్ వైఫల్యాన్ని ఖాళీ జాబితాగా పొరబడదు.
3. `ccam attach` ప్రతి ఖాతా ఫోల్డర్‌లో, Claude Code బ్యాక్‌గ్రౌండ్ సెషన్‌ను ఉంచే `jobs/<id>` కోసం చూసి, `ccam <name>` Claude Code ను ప్రారంభించే విధంగానే, ఆ ఖాతాగా `claude attach <id>` ను రన్ చేస్తుంది.

## భద్రత, ప్రత్యేక సందర్భాలు

- బ్యాక్‌గ్రౌండ్ సెషన్ ids ఒకే ఖాతాకు చెందుతాయి: మరో ఖాతాగా రన్ చేసిన `claude attach` వాటిని కనుగొనలేదు, ఆ సమస్యనే `ccam attach` పరిష్కరిస్తుంది.
- రెండు ఖాతాలలో ఒకేలా మొదలయ్యే ids ఉన్న సెషన్‌లు ఉండవచ్చు; అప్పుడు `ccam attach` రెండింటి పేర్లు చెప్పి, id లో ఇంకొంత భాగం అడుగుతుంది.
- ccam తన షేర్ చేసిన MCP సర్వర్లను ఒకే ఆర్గ్యుమెంట్‌లో `--mcp-config=<file>` తో Claude Code కు పంపుతుంది, దాన్ని `claude agents`, `claude attach` అంగీకరిస్తాయి; బ్యాక్‌గ్రౌండ్ సెషన్‌లు మళ్లీ ప్రారంభమైనప్పుడు ఆ ఫైల్‌నే ఉపయోగిస్తూ ఉంటాయి.
- API కీ ఖాతా బ్యాక్‌గ్రౌండ్ సర్వీస్ తనను ప్రారంభించిన షెల్ నుంచి కీని తీసుకుంటుంది. అలాంటి ఖాతా మొదటి బ్యాక్‌గ్రౌండ్ సెషన్‌ను ccam ద్వారా ప్రారంభించండి (`ccam gw --bg …`), తద్వారా సర్వీస్‌కు కీ అందుతుంది.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [స్క్రిప్ట్‌ల కోసం ఖాతాను ఎంచుకోవడం](/features/pick)
- [డాష్‌బోర్డ్](/features/dashboard)
