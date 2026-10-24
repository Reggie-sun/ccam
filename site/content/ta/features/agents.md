---
title: எல்லாக் கணக்குகளின் பின்னணி அமர்வுகள்
description: ccam agents எல்லாக் கணக்குகளின் Claude Code பின்னணி agents-ஐயும் திறந்த அமர்வுகளையும் ஒரே அட்டவணையில் காட்டுகிறது; ccam attach ஒன்றை அதன் id மூலம் திறக்கிறது.
updated: 2026-10-23
order: 16
translation: machine
sourceHash: "092368f63363"
---

Claude Code பின்னணி அமர்வுகளை (`claude --bg`, agent view) கணக்கு வாரியாக வைத்திருக்கிறது: `claude agents` அது இயங்கும் கணக்கின் அமர்வுகளை மட்டுமே காட்டும். `ccam agents` எல்லாக் கணக்குகளிடமும் ஒரே நேரத்தில் கேட்டு அவற்றை ஒன்றாகப் பட்டியலிடுகிறது; `ccam attach <id|name>` ஒரு பின்னணி அமர்வை, அதைத் தொடங்கிய கணக்காகவே திறக்கிறது.

## பயன்பாடு

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` என்பது பின்னணி அமர்வின் நிலையும் (`working`, `blocked`, `done`, `failed`, `stopped`) அதன் process என்ன செய்கிறது என்பதும் (`busy`, `waiting`, `idle`); திறந்திருக்கும் interactive அமர்வு `open` எனக் காட்டும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `--all` | முடிந்த, தோல்வியடைந்த அல்லது நிறுத்தப்பட்ட பின்னணி அமர்வுகளையும் பட்டியலிடும் | off |
| `--json` | அமர்வுகளை JSON ஆக அச்சிடும், ஒவ்வொன்றும் அது சேர்ந்த `account` உடன் | off |

`ccam attach <id|name> [claude args]` id-க்குப் பிறகு வரும் அனைத்தையும், flags உட்பட, மாற்றமின்றி `claude attach`-க்கு அனுப்புகிறது: `ccam attach 49ec --dry-run` என்பது `--dry-run`-ஐ Claude Code-இடம் கொடுக்கிறது. அதற்குப் பதிலாகக் கட்டளையை அச்சிட, id-க்கு முன் `--dry-run` வையுங்கள் (`ccam --dry-run attach 49ec` அல்லது `ccam attach --dry-run 49ec`).

## இது எப்படி வேலை செய்கிறது

1. ஒவ்வொரு கணக்குக்கும், ccam அந்தக் கணக்கின் `CLAUDE_CONFIG_DIR` உடன் `claude agents --json`-ஐ (கொடுத்திருந்தால் `--all` உடன்) இயக்குகிறது; எல்லாம் ஒரே நேரத்தில், ஒவ்வொன்றும் 15 வினாடிகளுக்குள். இது கணக்கின் சொந்த அமர்வுப் பட்டியலைப் படிக்கிறது, Claude Code-இன் பின்னணி service-ஐத் தொடங்குவதில்லை.
2. ஒவ்வொரு உள்ளீடும் கணக்கின் பெயரைப் பெறுகிறது, பட்டியல் தொடக்க நேரத்தின்படி, புதியது முதலில் என வரிசைப்படுத்தப்படுகிறது. பதிலளிக்காத கணக்கு standard error-இல் பெயரிடப்பட்டு விலக்கப்படுகிறது, பிறகு ccam code 1 உடன் வெளியேறுகிறது, அதனால் ஒரு script தோல்வியை வெற்றுப் பட்டியல் என்று தவறாக எடுத்துக்கொள்ளாது.
3. Claude Code பின்னணி அமர்வை வைத்திருக்கும் இடமான `jobs/<id>`-ஐ `ccam attach` ஒவ்வொரு கணக்கின் கோப்புறையிலும் தேடுகிறது, பிறகு `ccam <name>` Claude Code-ஐத் தொடங்கும் அதே வழியில், அந்தக் கணக்காக `claude attach <id>`-ஐ இயக்குகிறது. எந்த id-இன் தொடக்கமாகவும் இல்லாத ஒரு சொல், பெரிய, சிறிய எழுத்து வேறுபாட்டைப் பொருட்படுத்தாமல், அமர்வுகளின் பெயர்களில் தேடப்படுகிறது; id போலத் தோன்றும் சொல் ஒருபோதும் பெயராக எடுத்துக்கொள்ளப்படுவதில்லை.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- பின்னணி அமர்வு ids ஒரு கணக்குக்கு மட்டுமே உரியவை: வேறொரு கணக்காக இயக்கப்படும் `claude attach` அவற்றைக் கண்டுபிடிக்க முடியாது; `ccam attach` தீர்க்கும் சிக்கல் இதுதான்.
- இரண்டு கணக்குகள் ஒரே மாதிரித் தொடங்கும் ids உள்ள, அல்லது பெயர்களில் அதே சொல் உள்ள அமர்வுகளை வைத்திருக்கலாம்; அப்போது `ccam attach` அவற்றைப் பெயரிட்டு, id அல்லது பெயரின் இன்னும் சில எழுத்துகளைக் கேட்கும்.
- ccam தன் பகிரப்பட்ட MCP சர்வர்களை ஒரே argument-இல் `--mcp-config=<file>` மூலம் Claude Code-க்குக் கொடுக்கிறது, அதை `claude agents`-உம் `claude attach`-உம் ஏற்கின்றன; பின்னணி அமர்வுகள் மீண்டும் தொடங்கும்போதும் அதே கோப்பைப் பயன்படுத்துகின்றன.
- API key கணக்கின் பின்னணி service அதைத் தொடங்கிய shell-இலிருந்து தன் key-ஐப் பெறுகிறது. அத்தகைய கணக்கின் முதல் பின்னணி அமர்வை ccam வழியாகத் தொடங்குங்கள் (`ccam gw --bg …`), அப்போதுதான் service-இடம் key இருக்கும்.

## தொடர்புடையவை

- [Claude Code-ஐத் தொடங்குதல்](/features/launch)
- [scripts-க்கு ஒரு கணக்கைத் தேர்ந்தெடுத்தல்](/features/pick)
- [டாஷ்போர்டு](/features/dashboard)
