---
title: கட்டளைக் குறிப்பு
metaTitle: ccam கட்டளைக் குறிப்பு
description: ccam add, ccam <name> முதல் ccam usage, ccam next, ccam doctor வரை, ccam-இன் ஒவ்வொரு கட்டளையும் option-உம்.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: "116c5b34ac9a"
---

ஒரு கட்டளை உண்மையில் எதையும் செய்யாமல் என்ன செய்யும் என்று பார்க்க, அதற்கு முன் `--dry-run` வையுங்கள், எடுத்துக்காட்டாக `ccam --dry-run work`. Claude Code-ஐத் தொடங்கும்போது (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) மற்றும் `group`, `migrate`, `shortcut`, `update`, `doctor --fix` ஆகியவற்றுடன் இது வேலை செய்யும். எதையாவது மாற்றும் மற்ற கட்டளைகள் இயங்குவதற்குப் பதிலாக இதை மறுக்கும்.

<!-- commands -->

## Claude Code-ஐத் தொடங்குதல்

`ccam <name|number> [arguments]` அந்தக் கணக்காக Claude Code-ஐத் தொடங்குகிறது. பெயருக்குப் பின் உள்ள arguments-ஐ ccam படிப்பதில்லை; அவற்றை அதே வரிசையில் `claude`-க்கு அனுப்புகிறது. standard input, exit code, signals ஆகியவை `claude`-ஐ நேரடியாக இயக்குவது போலவே செயல்படும். பெயருக்குப் பின் `--` விருப்பத்தேர்வு மட்டுமே.

கணக்குப் பெயர் இல்லாமல் arguments உடன் இயக்கும் `ccam`, எடுத்துக்காட்டாக `ccam -p "summarise"`, தற்போதைய கோப்புறைக்கு அமைக்கப்பட்ட கணக்கை, இல்லையெனில் இயல்புநிலைக் கணக்கைப் பயன்படுத்தும். அவற்றை `ccam folder <name> [folder]`, `ccam use <name>` மூலம் அமைக்கலாம்.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

கணக்குப் பெயருக்கு முன் வரும் `--fallback` என்பது `-p`/`--print` இயக்கங்களுக்கானது. கணக்கு தன் பயன்பாட்டு வரம்பை எட்டும்போது, ccam அதே கட்டளையை அதிக இடம் உள்ள உள்நுழைந்த கணக்காக மீண்டும் இயக்கி, முடிந்த இயக்கத்தின் output-ஐ மட்டும் அச்சிடுகிறது. பார்க்க: [கோப்புறை வாரியான கணக்குகள்](/features/folder), [வரம்பில் மீண்டும் இயக்குதல்](/features/next#run-again-at-the-limit).

`ccwork` போன்ற ஒரு குறுக்குவழியும் அதன் கணக்கை இதே முறையில் தொடங்கும்: `ccam work --resume` செய்வதையே `ccwork --resume` செய்யும். [குறுக்குவழிகள்](/docs/commands#shortcuts) பகுதியைப் பாருங்கள்.

## கணக்குகளைச் சேர்த்தல் மற்றும் நீக்குதல்

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

கோப்புறை இல்லாமல் இயக்கும் `ccam import`, `~/.claude` மற்றும் `~/.claude-*`-ஐத் தேடி, ஒவ்வொன்றைப் பற்றியும் கேட்கும். இறக்குமதி செய்யப்பட்ட கோப்புறைகள் அவை இருக்கும் இடத்திலேயே பதிவு செய்யப்படும்; ccam ஒருபோதும் உள்நுழைவை நகர்த்துவதோ நகலெடுப்பதோ இல்லை.

`ccam add --api` Claude உள்நுழைவுக்குப் பதிலாக API key பயன்படுத்தும் கணக்கைச் சேர்க்கும்; ccam ஒவ்வொரு தொடக்கத்திலும் key-ஐப் படிக்கும், அதை ஒருபோதும் சேமிக்காது. [API key கணக்குகள்](/features/api) பாருங்கள்.

ஒரு கணக்காக Claude Code இயங்கிக்கொண்டிருக்கும்போது, `ccam rename`, `ccam remove`, `ccam group` கட்டளைகள் ஆகியவை அந்தக் கணக்கை மாற்ற மறுக்கும். முதலில் அமர்வை முடியுங்கள்.

`ccam list --json` ஒவ்வொரு கணக்கையும் அதன் slot எண், குறுக்குவழி, மின்னஞ்சல், கோப்புறை ஆகியவற்றுடன் அச்சிடும். அதிலுள்ள `groups` array கணக்கின் குழுக்களைப் பட்டியலிடும், `share` object ஒவ்வொரு உருப்படியையும் (`mcp`, `settings`, `memory`, `extensions`, `history`) அது வரும் குழுவுடன், அல்லது `local` உடன் இணைக்கும்.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

ஒரு குழு தன் கணக்குகளுக்கு இடையே உருப்படிகளைப் பகிர்கிறது; ஒரு கணக்கு பல குழுக்களில் இருக்கலாம், ஆனால் ஒவ்வோர் உருப்படியும் ஒரே ஒரு குழுவிலிருந்து மட்டுமே வரும். [குழுக்களுடன் பகிர்தல்](/features/share) பாருங்கள்.

## ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`, ஒரு Claude Code கோப்புறையிலிருந்து `settings.json`, `CLAUDE.md`, நீட்டிப்புகள், பயனர் MCP சர்வர்கள் ஆகியவற்றை ஒரு குழுவுக்குள் நகலெடுக்கிறது, அதனால் அந்த உருப்படிகளைப் பகிரும் ஒவ்வொரு கணக்குக்கும் அவை கிடைக்கும். அந்தக் கோப்புறை படிக்கப்பட மட்டுமே செய்யும். `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` ஆகிய options நீங்கள் குறிப்பிடும் உருப்படிகளை மட்டுமே நகலெடுக்கும். அவை இல்லாவிட்டால், டெர்மினலில் ஒரு தேர்வுப் பட்டியல் (checklist) நீங்களே தேர்ந்தெடுக்க உதவும், அதில் உரையாடல் வரலாற்றைத் தவிர அனைத்தும் தேர்வுசெய்யப்பட்டிருக்கும்; மற்ற இடங்களில் உரையாடல் வரலாற்றைத் தவிர அனைத்தும் நகலெடுக்கப்படும். `-g` அல்லது `--group` குழுவைத் தேர்ந்தெடுக்கும், `-y` அல்லது `--yes` கேட்காமல் நகலெடுத்து, ஒரு கோப்பு வேறுபட்டால் குழுவின் நகலை வைத்திருக்கும். சர்வர்களின் சேமிக்கப்பட்ட உள்நுழைவுகளையும் அது ஒவ்வொரு கணக்கின் சொந்த உள்நுழைவுக் கோப்புக்குள் நகலெடுக்கிறது. [ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்](/docs/sharing#copy-an-existing-setup) பகுதியைப் பாருங்கள்.

## MCP சர்வர்கள்

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` அந்தக் கணக்கிடம் மட்டும் உள்ள சர்வர்களை அதன் குழுவுக்கு நகர்த்தும். குழு ஏற்கெனவே வேறொரு வரையறைக்குப் பயன்படுத்தும் பெயருடைய சர்வர் தவிர்க்கப்படும்; குழுவின் நகலை மாற்ற `--force` சேருங்கள், ccam முதலில் அதை backup எடுக்கும்.

## பயன்பாடு மற்றும் கணக்கு மாற்றம்

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback`, `ccam pick` மூன்றும் ஒரே முறையில் தேர்ந்தெடுக்கும்: முதலில் `auto_max_usage`-க்குக் (`config.toml`, இயல்புநிலை 100) கீழ் உள்ளவை, `auto_strategy` (`headroom`, `expiring` அல்லது `balance`) படி, `--model`-இன் வாராந்திரக் கால அளவையும் கணக்கிட்டு, 100%-இல் உள்ள கணக்கை ஒருபோதும் அல்ல, அவை தொடங்கும் கணக்கின் நிறுவனத்துக்குள்ளேயே (அதற்கு அப்பாலும் தேட `--cross-org` அல்லது `cross_org = true`). எதையும் தேர்ந்தெடுக்க முடியாதபோது அவை code 75 உடன் வெளியேறும், அல்லது `--wait` உடன் முதல் கணக்கு மீட்டமைக்கப்படும் வரை காத்திருக்கும். வரலாற்றைப் பகிராத கணக்குகளுக்கு இடையிலும் `ccam next --continue` உரையாடலைக் கொண்டுசெல்கிறது: அப்போது Claude Code அதைப் புதிய கணக்குக்குள் நகலெடுக்கிறது. Claude Code-க்குள், plugin-இன் `/ccam:switch [account]` உரையாடலை ஒப்படைக்கிறது, அதன் hooks இயக்குவது `ccam hook`. [அடுத்த கணக்குக்கு மாறுதல்](/features/next), [scripts-க்கு ஒரு கணக்கைத் தேர்ந்தெடுத்தல்](/features/pick), [எல்லாக் கணக்குகளின் பின்னணி அமர்வுகள்](/features/agents) பாருங்கள்.

## பராமரிப்பு

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` ccam-ஐ மாற்றுவதற்கு முன் வெளியீட்டின் SHA-256 checksum-ஐச் சரிபார்க்கிறது, `--version v0.2.0` ஒரு குறிப்பிட்ட வெளியீட்டை நிறுவுகிறது. `.deb` தொகுப்பிலிருந்து நிறுவிய ccam-ஐ apt மூலம் புதுப்பியுங்கள். [புதுப்பித்தல்](/docs/install#update) பகுதியைப் பாருங்கள்.

## கோப்புகள்

| பாதை | உள்ளடக்கம் |
| --- | --- |
| `~/.config/ccam/config.toml` | கணக்குகள், குழுக்கள், அமைப்புகள் (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | ஒவ்வொரு கணக்கின் `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | பகிரப்பட்ட அமைப்புகள், MCP சர்வர்கள், நீட்டிப்புகள், வரலாறு |
| `~/.local/state/ccam/` | பயன்பாட்டு cache, கடைசியாகப் பயன்படுத்திய நேரங்கள், `/ccam:switch`-உம் வரம்பு hook-உம் விட்டுச்செல்லும் குறிப்புகள் (`handoff/`) |
| `~/.local/bin/<shortcut>` | `ccwork` போன்ற குறுக்குவழிக் கட்டளைகள், ccam-க்கான இணைப்புகள் |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` மாறிகள் இந்தக் கோப்புறைகளின் இடத்தை மாற்றுகின்றன. Windows-இல் `~` என்பது `%USERPROFILE%`, குறுக்குவழிகள் `<shortcut>.exe` ஆக இருக்கும்.
