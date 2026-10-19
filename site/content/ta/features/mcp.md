---
title: பகிரப்பட்ட MCP சர்வர்கள்
description: ccam mcp add, add-json, list, remove, import மூலம் ஒவ்வொரு குழுவுக்கும் MCP சர்வர்களின் ஒரே பட்டியலை வைத்திருங்கள்; அது --mcp-config மூலம் அனுப்பப்படும் அல்லது ஒவ்வொரு கணக்குக்குள்ளும் sync செய்யப்படும்.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code பயனர் MCP சர்வர்களை `.claude.json`-இல் வைக்கிறது, கணக்கின் உள்நுழைவையும் வைத்திருப்பது அதே கோப்புதான், அதனால் அந்தக் கோப்பைக் கணக்குகளுக்கு இடையே பகிர முடியாது. ccam பகிரப்பட்ட சர்வர்களை `mcp.json` என்ற குழுக் கோப்பில் வைத்து, தொடக்கத்தில் அதைக் குழுவின் ஒவ்வொரு கணக்குக்கும் கொடுக்கிறது. உங்கள் எல்லாக் கணக்குகளுக்கும் ஒரே முறை ஒரு சர்வரைச் சேர்க்க `ccam mcp`-ஐயும், ஓர் அமர்வுக்குள் நீங்கள் சேர்த்த சர்வரை மற்றவற்றுக்கு நகர்த்த `ccam mcp import`-ஐயும் பயன்படுத்துங்கள்.

## பயன்பாடு

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

ஒவ்வொரு கணக்கின் அடுத்த தொடக்கத்திலிருந்து மாற்றம் பொருந்தும்; இயங்கும் அமர்வுகள் தாங்கள் தொடங்கிய சர்வர்களையே வைத்திருக்கும்:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

குழுவிடம் ஏற்கெனவே உள்ள பெயரைச் சேர்த்தால் அதன் வரையறை மாற்றப்பட்டு, `added`-க்குப் பதிலாக `updated` என்று அச்சிடப்படும். ஒரு கணக்கு அந்தக் குழுவுடன் MCP-ஐப் பகிர்ந்தால் மட்டுமே (`ccam share <name> --mcp on`) சர்வர்கள் அதை அடையும்; இன்னும் எந்தக் கணக்கும் MCP-ஐப் பகிரவில்லை என்றால், `add` அதைச் சொல்லும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json`, `remove` ஆகியவற்றுக்கான இலக்குக் குழு | ஏதாவது ஒரு கணக்கு MCP-ஐப் பகிரும் ஒரே குழு, இல்லையென்றால் `default` |
| `-- <command> [args…]` (`add`) | இந்தக் கட்டளையை stdio சர்வராக இயக்கும் | |
| `--url <url>` (`add`) | தொலைநிலைச் சர்வர் URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | ஒரு `--url` சர்வரின் transport | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, பலமுறை கொடுக்கலாம் | none |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio சர்வருக்கான சூழல் மாறி, பலமுறை கொடுக்கலாம் | none |
| `--force` (`import`) | அதே பெயருடைய, வேறுபடும் குழுச் சர்வரை மாற்றும்; பழைய `mcp.json` காப்புப்பிரதி எடுக்கப்படும் | off |
| `--keep-local` (`import`) | கணக்கின் `.claude.json`-இல் உள்ள உள்ளீடுகளை வைத்திருக்கும் | off |
| `-h`, `--help` | ஒவ்வொரு கட்டளைக்குமான உதவி | |
| `--dry-run` (global) | `mcp list` மட்டுமே இதை மதிக்கும்; `add`, `add-json`, `remove`, `import` இதை மறுத்து, எதையும் மாற்றாது | off |

`add`-க்கு `--`-க்குப் பின் ஒரு கட்டளையோ `--url`-ஓ தேவை; இல்லையென்றால் அது `a command to run (after --) or --url is required` என்று தோல்வியடையும். `--header` ஆனது `--url` சர்வர்களுக்கும், `--env` கட்டளைச் சர்வர்களுக்கும் பொருந்தும். `add-json` `claude mcp add-json` ஏற்கும் அதே JSON-ஐ ஏற்கும், வரையறை ஒரு `command` அல்லது `url` உள்ள object ஆக இருக்க வேண்டும். `import` எப்போதும் கணக்கு MCP-ஐப் பகிரும் குழுவையே பயன்படுத்தும், `list` ஒவ்வொரு குழுவையும் காட்டும், அதனால் `-g` அவற்றை மாற்றாது.

## இது எப்படி வேலை செய்கிறது

தொடர்புடைய கோப்புகள்:

- `~/.local/share/ccam/groups/<group>/mcp.json`: பகிரப்பட்ட சர்வர்கள், `{"mcpServers": {…}}`, `--mcp-config` படிக்கும் வடிவம். சர்வர் உள்ளீடுகளில் பெரும்பாலும் API keys இருப்பதால், அது `fsutil.AtomicWrite` மூலம் `0600` mode-உடன் எழுதப்படுகிறது.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: கணக்கின் சர்வர்களுடன் உள்ள Claude Code-இன் சொந்தக் கோப்பு. ccam அதன் `mcpServers` மதிப்பை மட்டும் `profile.TopLevelValue` scanner மூலம் படிக்கிறது, தொடக்கப் பாதையில் முழுக் கோப்பையும் ஒருபோதும் parse செய்வதில்லை.
- `~/.local/state/ccam/mcp/<name>/`: `seen` கோப்பு, `sync-stamp`, `.claude.json` காப்புப்பிரதிகள்.

`~/.config/ccam/config.toml`-இல் உள்ள `mcp_strategy` அமைக்கும் இரண்டு உத்திகள் உள்ளன. இயல்புநிலை flag உத்தி:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` கணக்கின் MCP குழுவைக் கண்டுபிடிக்கிறது. `mcp.FlagArgs` குழுவின் `mcp.json` உள்ளதா என்பதை parse செய்யாமல், ஒரே ஒரு `stat` மூலம் சரிபார்க்கிறது.
2. `Plan.InsertArgs` `--mcp-config=<path>`-ஐ `claude`-க்கு உடனே பின், உங்கள் சொந்த arguments-க்கு முன் வைக்கிறது, அதனால் உங்கள் சொந்த `--mcp-config` அல்லது `--strict-mcp-config` அதற்கு மேல் தொடர்ந்து பொருந்தும்.
3. `mcp.NewLocal` கடைசித் தொடக்கத்துக்குப் பிறகு கணக்கின் `.claude.json` மாறியதா என்று சரிபார்த்து, மாறியிருந்தால் மட்டுமே அதன் சர்வர்களை parse செய்கிறது.
4. `syscall.Exec` ccam-ஐ Claude Code-ஆல் மாற்றுகிறது.

பாதை இரண்டாக அல்ல, `--mcp-config=<path>` என்ற ஒரே argument ஆக அனுப்பப்படுகிறது. Claude Code-இன் `--mcp-config` பல மதிப்புகளை ஏற்கும், அதனால் `claude --mcp-config a.json "prompt"` prompt-ஐ இரண்டாவது config கோப்பாகப் படித்து `MCP config file not found` என்று தோல்வியடையும். ஒரே argument வடிவம் அங்கேயே நின்றுவிடும், அதைத் தொடரும் prompt-உம் `mcp list` போன்ற subcommand-உம் வேலை செய்யும். Claude Code `--mcp-config` சர்வர்களைக் கணக்கின் சொந்தச் சர்வர்களுடன் கூடுதலாக ஏற்றுகிறது; பெயர் மோதலில் `.claude.json`-இல் உள்ள உள்ளீடு முன்னுரிமை பெறும். flag உத்தி Remote Control-க்கு மட்டுமே `.claude.json`-இல் எழுதுகிறது.

`claude remote-control` (`rc`) `--mcp-config`-ஐ ஏற்பதில்லை, அதன் அமர்வுகள் சர்வர்களை `.claude.json`-இலிருந்து மட்டுமே படிக்கின்றன. `ccam work rc`-க்கு, ccam அந்த flag-ஐ விட்டுவிட்டு, டெர்மினலில் இருந்தால், குழுவின் சர்வர்களை `.claude.json`-க்குள் நகலெடுக்கவா என்று கேட்கும்; `sync` உத்தி செய்வது போலவே காப்புப்பிரதியுடன். நகல் உருவாகும் வரை அது கேட்கும், பிறகு குழுவின் `mcp.json` மாறிய பிறகு மட்டுமே மீண்டும் கேட்கும் ([Remote Control](/features/launch#remote-control)). பெயர் மோதலில் அத்தகைய நகல் `--mcp-config`-ஐ விட முன்னுரிமை பெறுவதால், `ccam mcp add`, `add-json`, `remove`, `import` ஆகியவை குழுவின் மற்ற கணக்குகள் வைத்திருக்கும் நகல்களைப் புதுப்பிக்கின்றன:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

ஓர் அமர்வுக்குள் `claude mcp add` மூலம் சேர்க்கப்படும் சர்வர் அந்தக் கணக்கின் `.claude.json`-இல் மட்டுமே சேரும். Claude Code ccam-ஐ மாற்றிவிட்டதால் அந்த நேரத்தில் ccam அதைப் பார்க்க முடியாது, அதனால் அடுத்த தொடக்கத்தில் சரிபார்க்கிறது. `seen` கோப்பு `.claude.json`-இன் மாற்ற நேரத்தையும் ஏற்கெனவே தெரிவிக்கப்பட்ட சர்வர் பெயர்களையும் வைத்திருக்கிறது. நேரம் வேறுபட்டால், ccam பெயர்களைக் குழுவுடன் ஒப்பிட்டு, ஒவ்வொரு புதிய சர்வருக்கும் ஒருமுறை, உள்ளீட்டுக்குக் காத்திருக்காமல், ஒரு குறிப்பை அச்சிடும்:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` உடன், சர்வர்கள் `claude mcp list` மற்றும் `/mcp`-இல் சாதாரண user சர்வர்களாகத் தெரியும், ccam தொடக்கத்துக்கு முன் அவற்றை ஒவ்வொரு கணக்கின் `.claude.json`-இலும் எழுதும்:

<!-- motion:mcp-sync -->

1. `mcp.Sync` குழுவின் `mcp.json` மாற்ற நேரத்தை `sync-stamp`-உடன் ஒப்பிடுகிறது. அது மாறவில்லை என்றால், கணக்கின் `.claude.json` திறக்கப்படுவதில்லை.
2. அது Claude Code-இன் சொந்த lock-ஐ எடுக்கிறது: கோப்புக்கு அருகில் `mkdir` மூலம் உருவாக்கப்படும் `.claude.json.lock` கோப்புறை. அது 2 வினாடிகள் வரை காத்திருக்கிறது, Claude Code-இன் lock library செய்வது போலவே 10 வினாடிகளுக்கு மேல் பழைய lock-ஐக் காலாவதியானதாகக் கருதுகிறது.
3. அது தற்போதைய `.claude.json`-ஐ state கோப்புறையில் உள்ள `backups/`-இல் சேமித்து, புதிய 5-ஐ வைத்திருக்கிறது.
4. அது ஒவ்வொரு குழுச் சர்வரையும் `mcpServers`-இல் அமைக்கிறது, முன்பு தான் எழுதிய, இப்போது குழுவில் இல்லாத சர்வர்களை நீக்குகிறது, கணக்கின் சொந்தச் சர்வர்களை வைத்திருக்கிறது, கோப்பை அதன் அசல் mode-உடன் atomic ஆக எழுதுகிறது.
5. அது புதிய stamp-ஐயும் தான் நிர்வகிக்கும் பெயர்களையும் பதிவு செய்கிறது, Claude Code `--mcp-config` இல்லாமல் தொடங்குகிறது.

Claude Code இயங்கும்போது `.claude.json`-ஐ மாற்றி எழுதுவதால், அந்தக் கணக்கின் அமர்வு இயங்கிக்கொண்டிருந்தால் `ccam: account work is already running, skipping MCP sync` என்று sync தவிர்க்கப்படும். தோல்வியடைந்த sync `ccam: MCP sync failed (launching anyway): …` என்று அச்சிட்டு, எப்படியும் தொடங்கும். அமர்வுக்குள் சேர்க்கப்பட்ட சர்வர்களுக்கான குறிப்பு flag உத்திக்கு மட்டுமே உரியது.

`ccam mcp import` ஒரு கணக்கின் சொந்தச் சர்வர்களை அதன் குழுவுக்குள் நகர்த்துகிறது:

<!-- motion:mcp-import -->

ccam கணக்கின் சர்வர்களை `mcp.LocalServers` மூலம் படித்து, ஒவ்வொன்றையும் அதே பெயருடைய குழுச் சர்வருடன் `mcp.SameDef` மூலம் ஒப்பிடுகிறது; அது key வரிசையையும் இடைவெளிகளையும் புறக்கணிக்கிறது. குழுவிடம் இல்லாத, அல்லது அதே வரையறையுடன் உள்ள சர்வர் நகலெடுக்கப்படும். வேறுபட்ட ஒன்று தவிர்க்கப்படும், ஏனெனில் மற்ற ஒவ்வொரு கணக்கும் பயன்படுத்துவது குழுவின் வரையறையைத்தான்:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` உடன் குழுவின் சர்வர் மாற்றப்படும், முந்தைய `mcp.json` முதலில் குழுக் கோப்புறைக்குள் `.ccam-backup/<time>/mcp.json`-இல் சேமிக்கப்படும். பிறகு `mcp.RemoveLocal` நகலெடுத்த உள்ளீடுகளைக் கணக்கின் `.claude.json`-இலிருந்து, அதே lock-இன் கீழ், state கோப்புறையில் காப்புப்பிரதி எடுத்த பிறகு, நீக்குகிறது. பெயர் மோதலில் கணக்கின் உள்ளீடு குழுவினுடையதை மீறிவிடும் என்பதால் இது தேவை. `--keep-local` அவற்றை விட்டுவிடும்.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- **உள்நுழைவுகள்.** OAuth மூலம் உள்நுழையும் சர்வர்கள் தங்கள் tokens-ஐச் சர்வர் பட்டியலில் அல்ல, ஒவ்வொரு கணக்கின் `.credentials.json`-இல் (`mcpOAuth`) வைத்திருக்கின்றன, அதனால் ஒரு சர்வரைப் பகிர்வது அதன் உள்நுழைவைப் பகிர்வதில்லை. ஒவ்வொரு கணக்கிலும் `/mcp` மூலம் ஒருமுறை உள்நுழையுங்கள், அல்லது [`ccam migrate`](/features/migrate) மூலம் `~/.claude`-இலிருந்து உள்நுழைவுகளை நகலெடுங்கள். குழு ஒருபோதும் உள்நுழைவுக் கோப்பைப் பெறாது.
- **இயங்கும் கணக்குகள்.** `import` குழுவுக்கு நகலெடுக்கும், ஆனால் இயங்கும் கணக்கின் சொந்த உள்ளீடுகளை விட்டுவிட்டு அதைச் சொல்லும்; அமர்வு முடிந்ததும் மீண்டும் இயக்குங்கள். Sync இயங்கும் கணக்குகளைத் தவிர்க்கும்.
- **Remote Control-க்கான நகல்கள்.** இயங்கும் கணக்கின் நகல் தொடப்படாது: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` உடன் எதையும் புதுப்பிக்கத் தேவையில்லை, ஏனெனில் ஒவ்வொரு தொடக்கமும் sync செய்கிறது.
- **உடைந்த கோப்புகள்.** `.claude.json` சரியான JSON இல்லையென்றால், sync அதைத் தொடாமல் விட்டு பிழையைத் தெரிவிக்கும். 2 வினாடிகளுக்குள் lock-ஐ எடுக்க முடியாவிட்டால், ccam `could not lock … (Claude Code is writing the file)` என்று தெரிவிக்கும்.
- **MCP-ஐப் பகிராத கணக்குகள்.** `import` மறுத்து, முதலில் `ccam share side --mcp on`-ஐ இயக்கச் சொல்லும்.
- **குழுப் பெயர்கள்.** `-g` ஒரு குழுப் பெயரை மட்டுமே ஏற்கும்; அது `groups/`-இன் கீழ் ஒரு கோப்புறையாகிறது, அதனால் `../x` போன்ற பெயர்கள் மறுக்கப்படும்.

## தொடர்புடையவை

- [கணக்குகளுக்கு இடையே பகிர்தல்](/features/share)
- [MCP சர்வர்கள், அமைப்புகள், வரலாற்றைப் பகிர்தல்](/docs/sharing)
- [ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்](/features/migrate)
- [எல்லாக் கட்டளைகளும்](/docs/commands)
