---
title: கட்டளைக் குறிப்பு
description: ccam add, ccam <name> முதல் ccam usage, ccam next, ccam doctor வரை, ccam-இன் ஒவ்வொரு கட்டளையும் option-உம்.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

எந்தக் கட்டளையும் உண்மையில் எதையும் செய்யாமல் என்ன செய்யும் என்று பார்க்க, அதற்கு முன் `--dry-run` வையுங்கள், எடுத்துக்காட்டாக `ccam --dry-run work`.

<!-- commands -->

## Claude Code-ஐத் தொடங்குதல்

`ccam <name|number> [arguments]` அந்தக் கணக்காக Claude Code-ஐத் தொடங்குகிறது. பெயருக்குப் பின் உள்ள arguments-ஐ ccam படிப்பதில்லை; அவற்றை அதே வரிசையில் `claude`-க்கு அனுப்புகிறது. standard input, exit code, signals ஆகியவை `claude`-ஐ நேரடியாக இயக்குவது போலவே செயல்படும். பெயருக்குப் பின் `--` விருப்பத்தேர்வு மட்டுமே.

கணக்குப் பெயர் இல்லாமல் arguments உடன் இயக்கும் `ccam`, எடுத்துக்காட்டாக `ccam -p "summarise"`, இயல்புநிலைக் கணக்கைப் பயன்படுத்தும். அதை `ccam use <name>` மூலம் அமைக்கலாம்.

`claude1`, `claude2` போன்ற குறுக்குவழிகள், `~/.local/bin`-இல் உள்ள ccam-க்கான இணைப்புகள். `ccam shims install --names` ஒவ்வொரு கணக்குக்கும் `claude-<name>`-ஐயும் உருவாக்குகிறது.

## கணக்குகளைச் சேர்த்தல் மற்றும் நீக்குதல்

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

கோப்புறை இல்லாமல் இயக்கும் `ccam import`, `~/.claude` மற்றும் `~/.claude-*`-ஐத் தேடி, ஒவ்வொன்றைப் பற்றியும் கேட்கும். இறக்குமதி செய்யப்பட்ட கோப்புறைகள் அவை இருக்கும் இடத்திலேயே பதிவு செய்யப்படும்; ccam ஒருபோதும் உள்நுழைவை நகர்த்துவதோ நகலெடுப்பதோ இல்லை.

## பகிர்வு

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ஒவ்வொரு உருப்படியும் `on` (`default` குழு), `off`, அல்லது ஒரு குழுவின் பெயரை ஏற்கும். [பகிர்வு](/docs/sharing) பக்கத்தைப் பாருங்கள்.

## MCP சர்வர்கள்

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## பயன்பாடு மற்றும் கணக்கு மாற்றம்

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## பராமரிப்பு

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## கோப்புகள்

| பாதை | உள்ளடக்கம் |
| --- | --- |
| `~/.config/ccam/config.toml` | கணக்குகளும் அவை பகிர்வதும் |
| `~/.local/share/ccam/profiles/<name>/` | ஒவ்வொரு கணக்கின் `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | பகிரப்பட்ட அமைப்புகள், MCP சர்வர்கள், நீட்டிப்புகள், வரலாறு |
| `~/.local/state/ccam/` | பயன்பாட்டு cache மற்றும் கடைசியாகப் பயன்படுத்திய நேரங்கள் |
| `~/.local/bin/claude1`, `claude2`… | ccam-க்கான இணைப்புகள் |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` மாறிகள் இந்தக் கோப்புறைகளின் இடத்தை மாற்றுகின்றன.
