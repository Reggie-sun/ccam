---
title: extensions-ஐ ஆய்வு செய்தல்
description: ccam ext ஒவ்வொரு கணக்கும் ஏற்றும் skills, commands, plugins, MCP சர்வர்கள், hooks, ஒவ்வொன்றும் எங்கிருந்து வருகிறது, எது உடைந்துள்ளது என்பதைக் காட்டுகிறது.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

"இந்த skill (அல்லது MCP சர்வர்) ஏன் ஒரு கணக்கில் வேலை செய்கிறது, மற்றொன்றில் இல்லை?" என்ற கேள்விக்கு `ccam ext` பதில் தருகிறது. ஒவ்வொரு கணக்குக்கும், Claude Code ஏற்றும் extensions-ஐயும் ஒவ்வொன்றும் எங்கிருந்து வருகிறது என்பதையும் அது பட்டியலிடுகிறது: ஒரு பகிர்வுக் குழு, கணக்கின் சொந்தக் கோப்புறை, அல்லது நீங்கள் இருக்கும் திட்டம். அட்டவணைக்குக் கீழே, ஒரு extension நீங்கள் எதிர்பார்ப்பது போல் ஏறுவதைத் தடுப்பவற்றை அது பட்டியலிடுகிறது. அது படிக்க மட்டுமே செய்கிறது.

## பயன்பாடு

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

ஒவ்வொரு கட்டமும் அந்தக் கணக்கில் அந்த உருப்படியின் மூலத்தைக் குறிக்கிறது: பகிர்வுக் குழு (`default`), அது கணக்கின் சொந்தக் கோப்புறையில் இருந்தால் `account`, அல்லது கணக்கிடம் அது இல்லையென்றால் `-`. `!` சிக்கல் உள்ள உருப்படியைக் குறிக்கிறது, அது அட்டவணைக்குக் கீழே விளக்கப்படுகிறது.

## அது எதைப் படிக்கிறது

| வகை | எங்கே | காட்டப்படும் மூலம் |
| --- | --- | --- |
| Skills | `skills/`-இல் `SKILL.md` உள்ள ஒவ்வொரு கோப்புறையும் | `skills/` இணைக்கும் குழு, அல்லது `account` |
| Commands, agents | `commands/`, `agents/` கீழ் உள்ள `.md` கோப்புகள், துணைக் கோப்புறைகள் `folder/name` ஆக | குழு, அல்லது `account` |
| Plugins | `plugins/installed_plugins.json` | குழு, அல்லது `account` |
| MCP சர்வர்கள் | குழுவின் `mcp.json` (MCP பகிரப்படும்போது), கணக்கின் `.claude.json` | குழு, அல்லது `account` |
| Hooks | `settings.json`-இல் உள்ள `hooks`, நிகழ்வு வாரியாக எண்ணப்படுகின்றன | குழு, அல்லது `account` |
| இந்தத் திட்டம் | தற்போதைய கோப்புறையில் உள்ள `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, `.mcp.json` | ஒருமுறை, எல்லாக் கணக்குகளுக்கும் |

## அது தெரிவிக்கும் சிக்கல்கள்

- **உடைந்த இணைப்புகள்**: `skills/`, `commands/` அல்லது `agents/`-இல் இலக்கு காணாமல் போன இணைப்பு, அல்லது இல்லாத இடத்துக்கு இணைக்கும் பகிரப்பட்ட கோப்புறை (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` மூலம் சரிசெய்யுங்கள்).
- **`SKILL.md` இல்லாத skill கோப்புறை**: Claude Code அதைத் தவிர்க்கிறது.
- **இரண்டு இடங்களில் வேறுவிதமாக வரையறுக்கப்பட்ட MCP சர்வர்**: குழுவின் `mcp.json`, கணக்கின் `.claude.json`, அல்லது அவற்றில் ஒன்றும் திட்டத்தின் `.mcp.json`-உம். இரண்டு இடங்களிலும் ஒரே அமைப்புகள் இருந்தால் (`mcp_strategy = "sync"` எழுதுவது அதுதான்) அது சிக்கல் அல்ல. கணக்கின் சொந்தச் சர்வரை `ccam mcp import <account> <server>` மூலம் குழுவுக்கு நகர்த்துங்கள்.
- **சரியான JSON இல்லாத `installed_plugins.json`.**

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- `ccam ext` எதையும் மாற்றுவதில்லை, network கோரிக்கைகள் எதையும் செய்வதில்லை. அது சர்வர் பெயர்களையும் hook எண்ணிக்கைகளையும் படிக்கிறது, MCP வரையறைகளில் உள்ள `env` அல்லது headers-இன் மதிப்புகளை ஒருபோதும் அல்ல.
- அது disk-இல் உள்ளதைக் காட்டுகிறது. நிறுவப்பட்டு ஆனால் அமைப்புகளில் அணைக்கப்பட்ட plugin, அல்லது managed settings-ஆல் தடுக்கப்பட்ட சர்வர், அப்போதும் தோன்றும்.
- அந்தத் திட்டத்தின் `.claude` கோப்புறையையும் `.mcp.json`-ஐயும் சேர்க்க, அதைத் திட்டக் கோப்புறையில் இயக்குங்கள்.

## தொடர்புடையவை

- [கணக்குகளுக்கு இடையே பகிர்தல்](/features/share)
- [பகிரப்பட்ட MCP சர்வர்கள்](/features/mcp)
- [சரிபார்த்தலும் பழுதுநீக்கமும்](/features/doctor)
