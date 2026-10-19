---
title: షార్ట్‌కట్ కమాండ్‌లు
description: ccam shortcut తో ఒక ఖాతాకు ccwork వంటి సొంత కమాండ్‌ను ఇవ్వండి, ccam shims install తో ~/.local/bin లోని లింక్‌లను మళ్లీ సృష్టించండి.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

షార్ట్‌కట్ అంటే ఒక ఖాతాగా Claude Code ను ప్రారంభించే, మీరే పేరు పెట్టిన కమాండ్: `ccam work --resume` చేసేదే `ccwork --resume` చేస్తుంది. ఇది `~/.local/bin` లో ccam కు ఉన్న లింక్, కాబట్టి bash, zsh, స్క్రిప్ట్‌లు, `PATH` నుంచి ప్రోగ్రామ్‌లను రన్ చేసే మరేదైనా దీనితో పనిచేస్తుంది. షార్ట్‌కట్‌లు ఐచ్ఛికం; మీరు తరచుగా ప్రారంభించే ఖాతాల కోసం ఒకటి సృష్టించండి.

## వినియోగం

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam ప్రతి మార్పును నిర్ధారిస్తుంది:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

ఒక ఖాతాకు గరిష్ఠంగా ఒక షార్ట్‌కట్ ఉంటుంది, అది `~/.config/ccam/config.toml` లో `shortcut` గా నిల్వ అవుతుంది. `ccam add` విజార్డ్ `Create a shortcut command for <name>? [y/N]` అని, ఆపై `Shortcut name [cc<name>]:` అని అడుగుతుంది. విజార్డ్ బయట, `ccam add` `--shortcut` తో మాత్రమే ఒకటి సృష్టిస్తుంది. `ccam list` ప్రతి ఖాతా షార్ట్‌కట్‌ను చూపిస్తుంది, లేకపోతే `-` చూపిస్తుంది.

`ccam <name>` పేరు తర్వాత చేసే విధంగానే, షార్ట్‌కట్ ప్రతి ఆర్గ్యుమెంట్‌ను Claude Code కు పంపుతుంది. ఇందులో `--dry-run`, `--help` కూడా ఉన్నాయి: ఒక ప్రారంభాన్ని ముందుగా చూడటానికి, `ccam --dry-run work` ను ఉపయోగించండి.

## ఆప్షన్లు

| ఆప్షన్ | ఏమి చేస్తుంది | డిఫాల్ట్ |
| --- | --- | --- |
| `--remove` | ఖాతా షార్ట్‌కట్‌ను, దాని లింక్‌ను తొలగిస్తుంది | ఆఫ్ |
| `-h`, `--help` | కమాండ్ కోసం సహాయం | |
| `--dry-run` | `ccam shortcut` దీన్ని పాటిస్తుంది: పేరును తనిఖీ చేసి, సేవ్ చేయకుండా, లింక్ చేయకుండా `Would set the shortcut of work to ccwork.` లేదా `Would remove ccwork, the shortcut of work.` అని ప్రింట్ చేస్తుంది | ఆఫ్ |

`ccam shims install` ఆర్గ్యుమెంట్‌లు తీసుకోదు, `--dry-run` ను `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` తో తిరస్కరిస్తుంది.

పేరు నియమాలు, ఈ క్రమంలో తనిఖీ చేయబడతాయి:

- `a-z`, `0-9`, `.`, `-`, `_` అక్షరాలతో 1 నుంచి 64 అక్షరాలు, ఒక అక్షరం లేదా అంకెతో మొదలవ్వాలి;
- `ccam`, `claude` లేదా `run`, `usage` వంటి ccam కమాండ్ కాకూడదు;
- మరో ఖాతా షార్ట్‌కట్ కాకూడదు;
- ఇప్పటికే ccam కు లింక్ కాని, `~/.local/bin` లోని ఫైల్ కాకూడదు;
- మీ `PATH` లోని మరో ఎగ్జిక్యూటబుల్ పేరు కాకూడదు.

కమాండ్‌ను, `--remove` ను కలిపి ఇస్తే `give a command or --remove, not both` తో ఆగుతుంది. ఖాతాకు ఇప్పటికే ఉన్న పేరునే సెట్ చేస్తే `work: ccwork (unchanged)` అని ప్రింట్ చేస్తుంది, లింక్ లేకపోతే దాన్ని మళ్లీ సృష్టిస్తుంది.

## ఇది ఎలా పనిచేస్తుంది

పేరును సెట్ చేయడం `checkShortcut` ను రన్ చేస్తుంది, అది bin ఫోల్డర్‌తో, ccam బైనరీ అసలు పాత్‌తో `shim.Check` ను పిలుస్తుంది:

<!-- motion:shortcuts-check -->

పేరు తనిఖీల గుండా క్రమంగా కిందికి వెళ్తుంది. ప్రతి తనిఖీ పక్కన అది తిరస్కరించే ఒక పేరు ఉంటుంది: `Work` ప్యాటర్న్‌లో విఫలమవుతుంది, `run` ఒక ccam కమాండ్, `ccside` మరో ఖాతాకు చెందినది, `ccfile` `~/.local/bin` లో ఇప్పటికే ఉన్న ఫైల్, `git` `/usr/bin/git` ను దాచేస్తుంది. ప్రతి తనిఖీలో నెగ్గిన పేరు మాత్రమే సేవ్ చేయబడి లింక్ అవుతుంది.

1. `shim.Check` పేరును `^[a-z0-9][a-z0-9._-]{0,63}$` తో, ఆపై `config.Reserved` తో, ఆపై `cfg.FindShortcut` తో సరిపోలుస్తుంది. అది `~/.local/bin/<name>` పై `os.Lstat` ను పిలుస్తుంది, `PATH` లోని ప్రతి ఇతర ఫోల్డర్‌కు `os.Stat` ను పిలిచి, ccam కు లింక్ కాని ఎగ్జిక్యూటబుల్ ఫైల్‌ను తిరస్కరిస్తుంది.
2. ccam `fsutil.AtomicWrite` (తాత్కాలిక ఫైల్, `fsync`, rename) తో `config.toml` ను సేవ్ చేస్తుంది.
3. `shim.Remove` పాత షార్ట్‌కట్ లింక్‌ను తొలగిస్తుంది, కానీ అది ccam కు లింక్ అయితేనే.
4. `shim.Install` కాన్ఫిగర్ చేసిన ప్రతి షార్ట్‌కట్‌కు లేని లింక్‌ను, `os.Executable`, `filepath.EvalSymlinks` తో కనుగొన్న ccam బైనరీకి `os.Symlink` తో సృష్టిస్తుంది. అది ఎప్పుడూ దేనినీ తీసివేయదు.

`os.Readlink` ccam పాత్‌ను ఇచ్చినప్పుడు, లేదా లింక్, ccam పాత్ ఒకే ఫైల్‌కు చేరినప్పుడు, ఆ లింక్ ccam సొంతమైనదిగా లెక్కించబడుతుంది. ప్రతి తొలగింపు, ప్రతి ఓవర్‌రైట్ తనిఖీ ఈ పరీక్ష ద్వారా వెళ్తుంది, కాబట్టి ccam తాను సృష్టించని ఫైళ్లను ముట్టుకోదు.

మీరు షార్ట్‌కట్‌ను రన్ చేసినప్పుడు, ccam కు కావలసింది మీరు టైప్ చేసిన పేరు ఒక్కటే:

<!-- motion:shortcuts-link -->

షెల్ `PATH` లో `ccwork` ను కనుగొంటుంది, కెర్నల్ లింక్‌ను అనుసరించి ccam బైనరీని ప్రారంభిస్తుంది, `argv[0]` ఇంకా `ccwork` గానే ఉంటుంది. `app.Main` లో, `ccam` కాని ఏ పేరైనా `config.toml` ను ఒకసారి లోడ్ చేసి `shim.Lookup` ను పిలుస్తుంది, అది `shortcut` `filepath.Base(argv[0])` కు సమానమైన ఖాతాను కనుగొంటుంది. ఆపై `launcher.LaunchWith` ఇప్పటికే లోడ్ అయిన కాన్ఫిగ్‌తో [`ccam <name>`](/features/launch) తీసుకునే అవే దశలను తీసుకుని, `argv[1:]` ను మార్పు లేకుండా `claude` కు పంపుతుంది.

షెల్ అలియాస్ లేదా వ్రాపర్ స్క్రిప్ట్ కాకుండా సింలింక్ ఎందుకు: సింలింక్ ప్రతి షెల్‌లో, స్క్రిప్ట్‌లలో పనిచేస్తుంది, అదనపు షెల్ ప్రారంభాన్ని (సుమారు 1 నుంచి 2 ms) జోడించదు, ప్రారంభం ఇంకా ఒకే ఫైల్‌ను చదువుతుంది. `scripts/perf-check.sh` `ccam <name>` కు ఉన్న అదే 10 ms (p95) పరిమితితో `hyperfine` ద్వారా షార్ట్‌కట్‌ను కొలుస్తుంది.

## భద్రత, ప్రత్యేక సందర్భాలు

- **ఏ ఖాతాకూ చెందని లింక్.** `~/.local/bin` లో ccam కు ఉన్న లింక్ పేరును ఏ ఖాతా ఉపయోగించకపోతే, ఉదాహరణకు మీరు `config.toml` లో `shortcut` ను చేతితో మార్చిన తర్వాత, ccam డిఫాల్ట్ ఖాతాకు మారదు. అది `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` తో ఆగుతుంది. `ccam shortcut` తో పేరును మళ్లీ సెట్ చేయండి, లేదా లింక్‌ను తొలగించండి.
- **ఇప్పటికే ఉన్న ఫైళ్లు.** ccam తాను సృష్టించని ఫైల్‌ను ఎప్పుడూ భర్తీ చేయదు. `ccam shims install` `<path> already exists and is not a ccam link, so it was left alone` తో ఆగుతుంది.
- **తరలించిన ccam బైనరీ.** అప్పుడు పాత లింక్‌లు ఇక లేని పాత్‌కు చూపిస్తాయి, ccam వాటిని తన సొంతవిగా లెక్కించదు. వాటిని తొలగించి, ఆపై `ccam shims install` ను రన్ చేయండి.
- **పేరు మార్చడం, తీసివేయడం.** `ccam rename` షార్ట్‌కట్ పేరును ఉంచుతుంది, లింక్ పనిచేస్తూనే ఉంటుంది, ఎందుకంటే అది ccam కు చూపిస్తుంది, ccam `config.toml` లో పేరును చూస్తుంది. `ccam remove` ఖాతా లింక్‌ను తొలగిస్తుంది. `ccam shortcut <account> <new>` కొత్తది సృష్టించే ముందు పాత లింక్‌ను తొలగిస్తుంది.
- **PATH.** ccam ఒక లింక్‌ను సృష్టించినప్పుడు `~/.local/bin` మీ `PATH` లో లేకపోతే, అది `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` అని, దాని తర్వాత జోడించాల్సిన `export PATH=...` లైన్‌ను ప్రింట్ చేస్తుంది.
- **తర్వాత తనిఖీ చేయడం.** `ccam doctor` `missing shortcuts in /home/me/.local/bin: ccwork` అని నివేదిస్తుంది, ఫోల్డర్ `PATH` లో లేనప్పుడు హెచ్చరిస్తుంది. `ccam doctor --fix` లేని లింక్‌లను సృష్టిస్తుంది.

## సంబంధిత పేజీలు

- [Claude Code ను ప్రారంభించడం](/features/launch)
- [ఖాతాలను నిర్వహించడం](/features/accounts)
- [కమాండ్ రిఫరెన్స్](/docs/commands#shortcuts)
- [ccam ను ఇన్‌స్టాల్ చేయండి](/docs/install)
