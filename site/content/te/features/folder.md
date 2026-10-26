---
title: ఫోల్డర్ వారీగా ఖాతాలు
description: ఖాతా పేరు లేకుండా Claude Code ను ప్రారంభిస్తే, ccam folder ఒక ఫోల్డర్‌ను, దాని కింద ఉన్నవాటిని ఒకే ఖాతా వాడేలా చేస్తుంది, ఉదా. మీ కంపెనీ కోడ్ కోసం work.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: "25f82da64580"
---

`ccam folder work ~/code/company` తర్వాత, `~/code/company` లో గానీ దాని సబ్‌ఫోల్డర్లలో గానీ ఖాతా పేరు లేని ఏ లాంచ్ అయినా Claude Code ను `work` గా ప్రారంభిస్తుంది. దాని బయట, డిఫాల్ట్ ఖాతా మునుపటిలాగే ప్రారంభమవుతుంది. `ccam personal` తో లేదా `ccpersonal` వంటి షార్ట్‌కట్‌తో ఖాతా పేరు చెబితే, ఎప్పుడూ ఆ ఖాతానే ప్రారంభమవుతుంది.

## వినియోగం

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

ఈ లాంచ్‌లు ఫోల్డర్ ఖాతాను ఉపయోగిస్తాయి:

- `ccam -p "…"`, మొదటి ఆర్గ్యుమెంట్ Claude Code కోసం ఆప్షన్ అయిన ఇతర ఏ లాంచ్ అయినా
- పేరు లేని `ccam run`, `ccam --dry-run`
- `ccam --fallback -p "…"` దాని మొదటి ప్రయత్నం కోసం
- డాష్‌బోర్డ్, ఒంటరి `ccam`, ఇది ఫోల్డర్ ఖాతాపై కర్సర్‌తో ప్రారంభమవుతుంది, మీరు సింబాలిక్ లింక్ ద్వారా ఫోల్డర్‌కు చేరుకున్నప్పటికీ

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `<account>` | ఉపయోగించాల్సిన ఖాతా పేరు లేదా నంబర్ | |
| `[folder]` | ఫోల్డర్, relative లేదా absolute | ప్రస్తుత ఫోల్డర్ |
| `--remove` | సరిగ్గా ఆ ఫోల్డర్ ఎంట్రీని తీసివేస్తుంది | ఆఫ్ |
| `--dry-run` | మద్దతు ఉంది: మార్పును సేవ్ చేయకుండా ప్రింట్ చేస్తుంది | ఆఫ్ |

## ఇది ఎలా పనిచేస్తుంది

ప్రతి ఫోల్డర్ `~/.config/ccam/config.toml` లో ఒక `[[folder]]` ఎంట్రీ, ప్రతి లాంచ్ ఎలాగూ చదివే ఫైల్ ఇదే:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` ఫోల్డర్‌ను, సింబాలిక్ లింక్‌లు పరిష్కరించిన absolute పాత్‌గా మార్చి సేవ్ చేస్తుంది. ఒక ఫోల్డర్‌ను మళ్లీ సెట్ చేస్తే దాని ఖాతా మారుతుంది.
2. పేరు లేని లాంచ్ `launcher.HereProfile` ను పిలుస్తుంది. ఎంట్రీలు లేకపోతే అది వెంటనే డిఫాల్ట్ ఖాతాను తీసుకుంటుంది.
3. లేకపోతే అది working directory ని ప్రతి `path` తో strings గా పోలుస్తుంది: ఫోల్డర్ స్వయంగా లేదా దాని కింద ఉన్నది, అవే అక్షరాలతో మొదలయ్యే ఫోల్డర్ మాత్రం ఎప్పుడూ కాదు (`/code/company` అనేది `/code/companyx` ను కవర్ చేయదు). అత్యంత లోతైన ఎంట్రీ గెలుస్తుంది, కాబట్టి `~/code/company` `work` గా ఉన్నప్పుడు `~/code` `personal` గా ఉండవచ్చు.
4. working directory ఇచ్చిన విధంగా దేనితోనూ సరిపోలనప్పుడు మాత్రమే, ccam దాని సింబాలిక్ లింక్‌లను పరిష్కరించి మళ్లీ పోలుస్తుంది.
5. Windows, macOS లలో, వాటి ఫైల్ సిస్టమ్‌లు సాధారణంగా చేసినట్లే, పోలిక పెద్ద, చిన్న అక్షరాల తేడాను పట్టించుకోదు.

`config.toml` తప్ప మరే ఫైల్ చదవబడదు, కాబట్టి ఒక ఫోల్డర్ లోపల లాంచ్‌కు కొలవదగిన అదనపు సమయం పట్టదు (ఎంట్రీలతో లేదా లేకుండా 4.0 ms, `hyperfine`, 300 రన్‌లు).

## భద్రత, ప్రత్యేక సందర్భాలు

- `ccam rename` ఫోల్డర్ ఎంట్రీలను కొత్త పేరుకు తీసుకెళ్తుంది, `ccam remove` తొలగించిన ఖాతా ఎంట్రీలను తొలగిస్తుంది.
- ఖాతా ఇక లేని ఎంట్రీ ఆ ఫోల్డర్‌లో పేరు లేని లాంచ్‌ను ఆపుతుంది: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` అలాంటి ఎంట్రీలను, ఫోల్డర్ లేకుండా పోయిన ఎంట్రీలను నివేదిస్తుంది.
- పై ఫోల్డర్ నుంచి ఖాతాను పొందే సబ్‌ఫోల్డర్‌లో మీరు `ccam folder --remove` రన్ చేస్తే, ccam ఆ పై ఫోల్డర్ పేరును చెబుతుంది, దేన్నీ తీసివేయదు.
- `ccam next` రొటేషన్ లేదా వినియోగం ప్రకారం ఎంచుకుంటుంది, ఫోల్డర్‌లను చూడదు.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [షార్ట్‌కట్ కమాండ్‌లు](/features/shortcuts)
- [ఖాతాలను నిర్వహించడం](/features/accounts)
- [కమాండ్ రిఫరెన్స్](/docs/commands)
