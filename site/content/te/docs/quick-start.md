---
title: క్విక్ స్టార్ట్
metaTitle: "క్విక్ స్టార్ట్: Claude Code ఖాతాలను జోడించి మార్చండి"
description: ccam add తో రెండు Claude Code ఖాతాలను జోడించి, వాటికి ccwork వంటి షార్ట్‌కట్‌లు ఇచ్చి, వాటి మధ్య మారి, అవి ఏమి షేర్ చేయాలో ఎంచుకోండి.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

`ccam add` తో ప్రతి ఖాతాను జోడించండి, తర్వాత `ccam <name>` తో లేదా `ccwork` వంటి మీరే పేరు పెట్టిన షార్ట్‌కట్‌తో Claude Code ను ప్రారంభించండి. చాలా మందికి కావలసింది ఇంతే.

## 1. మీ మొదటి ఖాతాను జోడించండి

```sh
ccam add
```

విజార్డ్ `personal` లేదా `work` వంటి ఒక పేరును, ఈ ఖాతా మీ ఇతర ఖాతాలతో ఏ అంశాలను షేర్ చేయాలో, ఒక షార్ట్‌కట్ కమాండ్ సృష్టించాలా వద్దా (దశ 4 చూడండి) అని అడుగుతుంది. తర్వాత అది `claude auth login` ను రన్ చేస్తుంది, ఇది మీరు సైన్ ఇన్ చేయడానికి బ్రౌజర్‌ను తెరుస్తుంది. రిమోట్ మెషీన్‌లో, బ్రౌజర్ లోడ్ కాని `http://localhost:…/callback` పేజీ వద్ద ఆగిపోవచ్చు; ఆ అడ్రస్‌ను టెర్మినల్‌లో పేస్ట్ చేస్తే ccam సైన్-ఇన్‌ను పూర్తి చేస్తుంది. పూర్తయినప్పుడు మీకు ఇది కనిపిస్తుంది:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. మీ ప్రస్తుత సెట్టింగ్‌లను తీసుకురండి

మీరు ఇంతకు ముందు Claude Code ఉపయోగించి ఉంటే, మీ సెట్టింగ్‌లు, `CLAUDE.md`, ఎక్స్‌టెన్షన్లు, MCP సర్వర్లు `~/.claude` లో ఉంటాయి. ప్రతి ఖాతాకు అవి అందేలా వాటిని షేర్డ్ సెట్టింగ్‌లలోకి కాపీ చేయండి:

```sh
ccam migrate
```

ccam తాను ఏమి కాపీ చేస్తుందో చూపించి, ముందుగా అడుగుతుంది. ఇది `~/.claude` ను చదువుతుంది మాత్రమే, ఏ లాగిన్‌నూ కాపీ చేయదు, మీరు `--history` జోడిస్తే తప్ప సంభాషణ చరిత్రను వదిలేస్తుంది. వివరాల కోసం [ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం](/docs/sharing#copy-an-existing-setup) చూడండి.

## 3. రెండో ఖాతాను జోడించండి

```sh
ccam add work
```

బ్రౌజర్‌లో మరో ఖాతాతో సైన్ ఇన్ చేయండి.

## 4. ఖాతాలు మార్చండి

`ccam work` Claude Code ను work ఖాతాగా ప్రారంభిస్తుంది. ఇంకా చిన్నది కావాలంటే, ఖాతాకు `ccwork` వంటి దాని సొంత షార్ట్‌కట్ కమాండ్ ఇవ్వండి. షేర్ చేసే అంశాల తర్వాత విజార్డ్ ఇలా అడుగుతుంది:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` నొక్కి, ఒక పేరు టైప్ చేయండి, లేదా సూచించిన పేరునే ఉపయోగించడానికి Enter నొక్కండి:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

లాగిన్ తర్వాత ccam దాన్ని నిర్ధారిస్తుంది:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

మీరు ఆ ప్రశ్నను దాటవేసి ఉంటే, ఎప్పుడైనా షార్ట్‌కట్‌ను జోడించవచ్చు:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

ఇప్పుడు వీటిలో ఏదైనా work ఖాతాగా Claude Code ను ప్రారంభిస్తుంది:

```sh
ccwork
ccam work
ccam 2
```

`ccam` ను ఒక్కటే రన్ చేస్తే [డాష్‌బోర్డ్](/features/dashboard) తెరుచుకుంటుంది: యారో కీలతో లేదా ఖాతా సంఖ్యతో ఒక ఖాతాకు వెళ్లి, Enter నొక్కండి.

షార్ట్‌కట్ లేదా ఖాతా పేరు తర్వాత ఉన్నదంతా మార్పు లేకుండా Claude Code కు వెళ్తుంది, కాబట్టి మీకు తెలిసిన ఆప్షన్లనే ఉపయోగిస్తూ ఉండవచ్చు:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. ఏమి షేర్ చేయాలో ఎంచుకోండి

డిఫాల్ట్‌గా, MCP సర్వర్లు, `settings.json`, గ్లోబల్ `CLAUDE.md`, ఏజెంట్లు, కమాండ్‌లు, స్కిల్స్, ప్లగిన్లు `default` గ్రూప్ ద్వారా ఖాతాల మధ్య షేర్ అవుతాయి. సంభాషణ చరిత్ర షేర్ కాదు. ఖాతాల మధ్య సంభాషణలను కొనసాగించాలనుకుంటే, దాన్ని గ్రూప్‌కు జోడించండి:

```sh
ccam group set default --add history
```

ఏదైనా మార్చే ముందు ccam ప్రతి ఖాతా కోసం తాను ఏమి తరలిస్తుందో చూపిస్తుంది. వివరాల కోసం [గ్రూప్‌లతో షేరింగ్](/features/share) చూడండి.

## 6. ఒక ఖాతా పరిమితిని చేరినప్పుడు

```sh
ccam usage
ccam next --continue
```

`ccam usage` అన్ని ఖాతాల 5 గంటల, వారపు వినియోగాన్ని చూపిస్తుంది. `ccam next --continue` అదే సంభాషణను తర్వాతి ఖాతాలో తెరుస్తుంది. స్క్రిప్ట్‌లలో `ccam --fallback work -p "…"` తనంతట తానే మారుస్తుంది. [వినియోగ పరిమితులు](/docs/usage-limits) చూడండి.

**Team, Enterprise సీట్లు.** ఆటోమేటిక్ ఎంపిక (`ccam next`, `--fallback`, `ccam pick`, డాష్‌బోర్డ్‌లోని `n`, మీరు ఆన్ చేసినప్పుడు [`auto_switch`](/features/next#auto-switch)) ఒకే Claude Team లేదా Enterprise సంస్థలోని సీట్ల మధ్య కూడా మారుస్తుంది. Anthropic సహాయ కేంద్రం "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team ప్లాన్‌లోని ప్రతి సీటు పేరుతో నిర్దేశించిన ఒక వ్యక్తికి కేటాయించబడుతుంది, ఆ వ్యక్తి మాత్రమే దాన్ని ఉపయోగించవచ్చు) అని, "Users can't route their work through another user's account to use that user's remaining usage" (యూజర్లు మరో యూజర్‌కు మిగిలి ఉన్న వినియోగాన్ని వాడుకోవడానికి తమ పనిని ఆ యూజర్ ఖాతా ద్వారా మళ్లించకూడదు) అని చెబుతుంది ([Team ప్లాన్‌లు](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise ప్లాన్‌లకు](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) కూడా ఇదే వర్తిస్తుంది). కాబట్టి సీట్ల మధ్య ఆటోమేటిక్ ఎంపికను ఉపయోగించడం Anthropic నిబంధనల ఉల్లంఘన కావచ్చు. అక్కడ దాన్ని ఉపయోగించాలా వద్దా అనేది మీ నిర్ణయం: దానిపై ఆధారపడే ముందు మీ సంస్థను, Anthropic ను సంప్రదించి నిర్ధారించుకోండి. వ్యక్తిగత ప్లాన్‌ల (Pro, Max) గురించి [తరచుగా అడిగే ప్రశ్నలలోని నిబంధనల గమనిక](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) చెబుతుంది.

## ఇప్పటికే ~/.claude ఉపయోగిస్తున్నారా?

మళ్లీ సైన్ ఇన్ చేయడానికి బదులు దాన్ని నమోదు చేయండి:

```sh
ccam import
```

ఫోల్డర్ ఉన్న చోటే ఉంటుంది, సాధారణ `claude` మునుపటిలాగే పనిచేస్తుంది.
