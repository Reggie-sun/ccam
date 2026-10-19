---
title: கட்டளைக் குறிப்பு
description: ccam add, ccam <name> முதல் ccam usage, ccam next, ccam doctor வரை, ccam-இன் ஒவ்வொரு கட்டளையும் option-உம்.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

எந்தக் கட்டளையும் உண்மையில் எதையும் செய்யாமல் என்ன செய்யும் என்று பார்க்க, அதற்கு முன் `--dry-run` வையுங்கள், எடுத்துக்காட்டாக `ccam --dry-run work`.

<!-- commands -->

## Claude Code-ஐத் தொடங்குதல்

`ccam <name|number> [arguments]` அந்தக் கணக்காக Claude Code-ஐத் தொடங்குகிறது. பெயருக்குப் பின் உள்ள arguments-ஐ ccam படிப்பதில்லை; அவற்றை அதே வரிசையில் `claude`-க்கு அனுப்புகிறது. standard input, exit code, signals ஆகியவை `claude`-ஐ நேரடியாக இயக்குவது போலவே செயல்படும். பெயருக்குப் பின் `--` விருப்பத்தேர்வு மட்டுமே.

கணக்குப் பெயர் இல்லாமல் arguments உடன் இயக்கும் `ccam`, எடுத்துக்காட்டாக `ccam -p "summarise"`, இயல்புநிலைக் கணக்கைப் பயன்படுத்தும். அதை `ccam use <name>` மூலம் அமைக்கலாம்.

`ccwork` போன்ற ஒரு குறுக்குவழியும் அதன் கணக்கை இதே முறையில் தொடங்கும்: `ccam work --resume` செய்வதையே `ccwork --resume` செய்யும். [குறுக்குவழிகள்](/docs/commands#shortcuts) பகுதியைப் பாருங்கள்.

## கணக்குகளைச் சேர்த்தல் மற்றும் நீக்குதல்

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

கோப்புறை இல்லாமல் இயக்கும் `ccam import`, `~/.claude` மற்றும் `~/.claude-*`-ஐத் தேடி, ஒவ்வொன்றைப் பற்றியும் கேட்கும். இறக்குமதி செய்யப்பட்ட கோப்புறைகள் அவை இருக்கும் இடத்திலேயே பதிவு செய்யப்படும்; ccam ஒருபோதும் உள்நுழைவை நகர்த்துவதோ நகலெடுப்பதோ இல்லை.

## குறுக்குவழிகள்

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

குறுக்குவழி என்பது ஒரு கணக்காக Claude Code-ஐத் தொடங்க, நீங்களே பெயரிடும் கட்டளை. `ccam work` போலவே, அது ஒவ்வொரு argument-ஐயும் மாற்றமின்றி Claude Code-க்கு அனுப்பும், அதனால் `ccwork --resume`, `ccwork -p "hi"` ஆகியவை நீங்கள் எதிர்பார்ப்பது போலவே வேலை செய்யும். அது `~/.local/bin`-இல் உள்ள ccam-க்கான ஓர் இணைப்பு.

குறுக்குவழிகள் விருப்பத்தேர்வு மட்டுமே. `ccam add` wizard ஒன்றை உருவாக்க வேண்டுமா என்று கேட்கும், `ccam add <name> --shortcut <command>` கேட்காமலே ஒன்றை உருவாக்கும், `ccam list` ஒவ்வொரு கணக்கின் குறுக்குவழியையும், இல்லையென்றால் `-`-ஐயும் காட்டும்.

பெயர்களில் a-z, 0-9, `.`, `-`, `_` ஆகிய எழுத்துகளைப் பயன்படுத்தலாம். ஒரு பெயர் ccam கட்டளையாகவோ, மற்றொரு கணக்கின் குறுக்குவழியாகவோ, உங்கள் `PATH`-இல் உள்ள வேறொரு நிரலின் பெயராகவோ இருக்கக் கூடாது.

`ccam rename` குறுக்குவழிப் பெயரை மாற்றாமல் வைத்திருக்கும், `ccam remove` கணக்குடன் சேர்த்துக் குறுக்குவழியையும் நீக்கும். `ccam shims install` எல்லாக் கணக்குகளின் இணைப்புகளையும் மீண்டும் உருவாக்கும், எடுத்துக்காட்டாக நீங்கள் ccam-ஐ வேறிடத்துக்கு நகர்த்திய பிறகு.

## பகிர்வு

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ஒவ்வொரு உருப்படியும் `on` (`default` குழு), `off`, அல்லது ஒரு குழுவின் பெயரை ஏற்கும். [பகிர்வு](/docs/sharing) பக்கத்தைப் பாருங்கள்.

## ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`, ஒரு Claude Code கோப்புறையிலிருந்து `settings.json`, `CLAUDE.md`, நீட்டிப்புகள், பயனர் MCP சர்வர்கள் ஆகியவற்றை ஒரு குழுவுக்குள் நகலெடுக்கிறது, அதனால் அந்த உருப்படிகளைப் பகிரும் ஒவ்வொரு கணக்குக்கும் அவை கிடைக்கும். அந்தக் கோப்புறை படிக்கப்பட மட்டுமே செய்யும். `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` ஆகிய options நீங்கள் குறிப்பிடும் உருப்படிகளை மட்டுமே நகலெடுக்கும்; அவை இல்லாவிட்டால், உரையாடல் வரலாற்றைத் தவிர அனைத்தும் நகலெடுக்கப்படும். `-g` அல்லது `--group` குழுவைத் தேர்ந்தெடுக்கும், `-y` அல்லது `--yes` கேட்காமல் நகலெடுத்து, ஒரு கோப்பு வேறுபட்டால் குழுவின் நகலை வைத்திருக்கும். [ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்](/docs/sharing#copy-an-existing-setup) பகுதியைப் பாருங்கள்.

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
ccam shims install   # recreate the shortcut links
```

## கோப்புகள்

| பாதை | உள்ளடக்கம் |
| --- | --- |
| `~/.config/ccam/config.toml` | கணக்குகளும் அவை பகிர்வதும் |
| `~/.local/share/ccam/profiles/<name>/` | ஒவ்வொரு கணக்கின் `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | பகிரப்பட்ட அமைப்புகள், MCP சர்வர்கள், நீட்டிப்புகள், வரலாறு |
| `~/.local/state/ccam/` | பயன்பாட்டு cache மற்றும் கடைசியாகப் பயன்படுத்திய நேரங்கள் |
| `~/.local/bin/<shortcut>` | `ccwork` போன்ற குறுக்குவழிக் கட்டளைகள், ccam-க்கான இணைப்புகள் |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` மாறிகள் இந்தக் கோப்புறைகளின் இடத்தை மாற்றுகின்றன.
