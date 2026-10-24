---
title: ஏற்கெனவே உள்ள அமைப்பை நகலெடுத்தல்
description: ccam migrate மூலம் ~/.claude-இன் அமைப்புகள், CLAUDE.md, நீட்டிப்புகள், MCP சர்வர்களை ஒரு ccam குழுவுக்குள் நகலெடுங்கள்.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

ccam-க்கு முன்பே நீங்கள் Claude Code பயன்படுத்தியிருந்தால், உங்கள் அமைப்பு `~/.claude` மற்றும் `~/.claude.json`-இல் இருக்கும். `ccam migrate` அதை ஒரு பகிர்வுக் குழுவுக்குள் நகலெடுக்கிறது, அதனால் அந்த உருப்படிகளைக் குழு வழியாகப் பகிரும் ஒவ்வொரு கணக்கும் உங்களிடம் ஏற்கெனவே உள்ள அமைப்புடன் தொடங்கும். மூலக் கோப்புறை படிக்கப்பட மட்டுமே செய்யும், வெறும் `claude` முன்பு போலவே வேலை செய்யும். உங்கள் முதல் `ccam add`-க்குப் பிறகு ஒருமுறை இதை இயக்குங்கள்; `~/.claude.json`-இல் MCP சர்வர்கள் இருந்தால், முதல் `ccam add` தன் குறிப்பில் அவற்றைப் பெயரிடும்.

## பயன்பாடு

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

டெர்மினலில், உருப்படி flags இல்லாமலும் `--yes` இல்லாமலும், ccam முதலில் ஒரு தேர்வுப் பட்டியலைக் காட்டும்:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (அல்லது `x`) cursor-இன் கீழ் உள்ள உருப்படியை மாற்றும், அம்புக்குறி விசைகள் அல்லது `j`/`k` நகர்த்தும், Enter திட்டத்தைக் காட்டும், Esc, `q` அல்லது Ctrl+C `Cancelled.` என்று அச்சிடும். மூலத்தில் எதுவும் இல்லாத உருப்படிகள் `not found` எனக் குறிக்கப்படும். எல்லாவற்றையும் தேர்வுநீக்கினால், ccam `Nothing selected.` என்று அச்சிட்டு நின்றுவிடும்.

பிறகு ccam திட்டத்தை அச்சிட்டு, நகலெடுக்கும் முன் கேட்கும்:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

ஒரு குழு தன் சொந்தக் கணக்குகளை மட்டுமே, தான் பகிரும் உருப்படிகளுக்கு மட்டுமே அடையும், அதனால் சுருக்கம் குழுவின் கணக்குகளைப் பட்டியலிடுகிறது. நகலெடுத்த உருப்படிகளில் ஒன்றைக் குழு பகிரவில்லை என்றால், அல்லது சில கணக்குகள் அதில் இல்லையென்றால், அதை மாற்றும் `ccam group set <group> --add <items>`, `ccam group join <group> <accounts>` கட்டளைகளை ccam அச்சிடும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `[folder]` | நகலெடுக்க வேண்டிய Claude Code கோப்புறை | `~/.claude` |
| `--settings` | `settings.json`-ஐ நகலெடுக்கும் | off |
| `--memory` | உலகளாவிய `CLAUDE.md`-ஐ நகலெடுக்கும் | off |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` ஆகியவற்றை நகலெடுக்கும் | off |
| `--mcp` | பயனர் MCP சர்வர்களை (`.claude.json`-இல் உள்ள `mcpServers`) நகலெடுக்கும் | off |
| `--history` | உரையாடல் வரலாற்றை (`projects/`, `history.jsonl`) நகலெடுக்கும் | off |
| `-g`, `--group string` | நகலெடுக்க வேண்டிய பகிர்வுக் குழு | `default` |
| `-y`, `--yes` | கேட்காமல் நகலெடுக்கும்; முரண்பாட்டில் குழுவின் நகல் வைக்கப்படும் | off |
| `--dry-run` (global) | திட்டத்தை அச்சிட்டு, எந்தக் கேள்விக்கும் நகலுக்கும் முன் நின்றுவிடும் | off |

உருப்படி flag இல்லாதபோது, டெர்மினலில் தேர்வுப் பட்டியல் முடிவு செய்யும். டெர்மினலுக்கு வெளியே, அல்லது `--yes` உடன், ccam உரையாடல் வரலாற்றைத் தவிர அனைத்தையும் நகலெடுக்கும். `projects/` பெரியதாக இருக்கக்கூடும் என்பதால் வரலாறு இயல்பாக விடப்படுகிறது.

`--group on` என்றால் `default`. `--group off` மற்றும் `--group local` `--group takes a group name, not "off"` என்று மறுக்கப்படும். `ccam migrate` `--dry-run`-ஐ மதிக்கும், ஒரு தேர்வை முன்னோட்டமாகப் பார்க்க dry run-இலும் தேர்வுப் பட்டியல் தோன்றும்.

## இது எப்படி வேலை செய்கிறது

`internal/cli/migrate.go` உருப்படிகளைத் தேர்ந்தெடுக்கிறது, பிறகு `migrate.Build` மூலத்தையும் குழுவையும் படித்து ஒரு திட்டத்தைத் திருப்பித் தருகிறது. நீங்கள் உறுதிப்படுத்திய பிறகு `Plan.Apply` வரும் வரை எதுவும் எழுதப்படுவதில்லை.

1. **மூலத்தைக் கண்டறிதல்.** இல்லாத, கோப்புறை அல்லாத, அல்லது தெரிந்த Claude Code உள்ளீடு எதுவும் இல்லாத பாதையை `checkSource` மறுக்கிறது. MCP சர்வர்களுக்கு மட்டும் என்றால், `~/.claude.json`-இல் சர்வர்கள் இருக்கும் வரை, காலியான அல்லது இல்லாத `~/.claude` பரவாயில்லை. `migrate.ClaudeJSON` கோப்புறைக்குரிய `.claude.json`-ஐக் கண்டுபிடிக்கிறது: அதற்குள், அல்லது உள்ளே ஒன்று இல்லாத `~/.claude`-க்கு, அதன் அருகில் உள்ள `~/.claude.json`. `CLAUDE_CONFIG_DIR` அமைக்கப்படாதபோது Claude Code அதை வைத்திருக்கும் இடம் அதுவே.
2. **உருப்படிகளைப் பாதைகளுடன் பொருத்துதல்.** ஒவ்வொரு பாதையும் குழுக்கள் பயன்படுத்தும் அதே அனுமதிப் பட்டியலான `share.Layout`-இலிருந்து வருகிறது. MCP சர்வர்கள் `.claude.json`-இலிருந்து `mcp.LocalServers` மூலம் படிக்கப்படுகின்றன; அது முழுக் கோப்பையும் parse செய்வதற்குப் பதிலாக `profile.TopLevelValue` scanner-ஐப் பயன்படுத்துகிறது.
3. **நகலெடுக்கும், ஒருபோதும் நகர்த்தாது.** கோப்புகள் `fsutil.AtomicWrite` அல்லது `fsutil.CopyFile` மூலம் `~/.local/share/ccam/groups/<group>/`-க்குச் செல்கின்றன; அது இலக்குக்கு அருகில் ஒரு தற்காலிகக் கோப்பை எழுதி அதன் பெயரை மாற்றுகிறது, அதனால் நின்றுபோன நகல் ஒருபோதும் பாதி வெட்டப்பட்ட கோப்பை விட்டுச்செல்லாது. MCP சர்வர்கள் குழுவின் `mcp.json`-க்குள் அதன் lock-இன் கீழ் செல்கின்றன: இடையில் மற்றொரு ccam சேர்க்கும் சர்வர் அப்படியே இருக்கும், [Remote Control](/guides/remote-control)-க்காக உருவாக்கப்பட்ட நகல்கள் பிறகு புதுப்பிக்கப்படுகின்றன.

ஒவ்வோர் உருப்படியும் ஒரு மூலப் பாதையை ஒரு குழுப் பாதையுடன் பொருத்துகிறது:

<!-- motion:migrate-items -->

இரண்டு வரலாற்று வரிசைகளும் `--history` உடன் மட்டுமே நகலெடுக்கப்படுவதால் கோடிட்டுக் காட்டப்படுகின்றன. `.credentials.json` மற்றும் `.claude.json` `share.Never`-இல் பட்டியலிடப்பட்டுள்ளன, அவை ஒருபோதும் குழுவை அடைவதில்லை: அவற்றில் ஒன்று layout-இல் தோன்றினால் `Build` ஓர் உள் பிழையைத் திருப்பித் தரும், கோப்புறைத் தேடல் அவற்றைப் பெயரால் தவிர்க்கும்.

`buildEntry` ஒவ்வொரு பாதைக்கும் ஒரு நிலையைக் கொடுக்கிறது; ஒவ்வொரு திட்ட வரியும் விவரிப்பது அதைத்தான்:

<!-- motion:migrate-states -->

- **Missing**: மூலத்தில் இல்லை, அல்லது கோப்புறை எதிர்பார்க்கப்படும் இடத்தில் கோப்பு உள்ளது. ஒவ்வோர் உருப்படிக்கும் ஒருமுறை `not found, skipped` எனப் பட்டியலிடப்படும்.
- **Linked**: மூலப் பாதை ஏற்கெனவே குழுவின் நகலுக்குத் தீர்க்கப்படுகிறது, எடுத்துக்காட்டாக மூலம் ஒரு ccam profile கோப்புறையாக இருக்கும்போது.
- **Same**: குழுவிடம் ஏற்கெனவே ஒரே மாதிரியான உள்ளடக்கம் உள்ளது.
- **New**: குழுவிடம் அது இல்லை, அதனால் முழுமையாக நகலெடுக்கப்படும்.
- **Merge**: இரு பக்கங்களிலும் உள்ள கோப்புறை. `copyMissing` குழுவிடம் இல்லாத கோப்புகளை மட்டுமே சேர்க்கிறது; வேறுபடும் கோப்புகள் குழுவிடம் உள்ளபடியே இருக்கும், திட்டம் `group copies kept` என்று சொல்லும். `plugins/`-க்கு, `installed_plugins.json` மற்றும் `known_marketplaces.json` registries ஒன்றிணைக்கப்படும், absolute plugin பாதைகள் குழுவின் கோப்புறையைச் சுட்டும்படி மாற்றி எழுதப்படும்.
- **Conflict**: இரு பக்கங்களிலும் வேறுபட்ட உள்ளடக்கத்துடன் உள்ள கோப்பு. கோப்பு வாரியாக நீங்கள் தேர்ந்தெடுப்பீர்கள்: `keep group`, `use mine`, அத்துடன் `.json` கோப்புகளுக்கு `merge JSON` (மோதும் மதிப்புகளில் குழு முன்னுரிமை பெறும் ஆழமான ஒன்றிணைப்பு) அல்லது `history.jsonl`-க்கு `merge by line`. வரையறை வேறுபடும் MCP சர்வரும் முரண்பாடுதான், `keep group` அல்லது `use mine` உடன். Enter அழுத்தும்போது பட்டியலின் முதல் தேர்வே இயல்புநிலை; `--yes` எப்போதும் குழுவின் நகலை வைத்திருக்கும்.

தேர்வுப் பட்டியல் என்பது `ccam add` பயன்படுத்தும் `checklist` component-இன் மேல் கட்டப்பட்ட `tui.MigratePicker`; `migrate.Found` lstat அழைப்புகள் மற்றும் MCP சர்வர் பட்டியல் மூலம் இல்லாத உருப்படிகளைக் குறிக்கிறது.

### MCP உள்நுழைவுகள்

OAuth மூலம் உள்நுழையும் MCP சர்வர்கள் தங்கள் tokens-ஐச் சர்வர் பட்டியலில் அல்ல, உள்நுழைவுக் கோப்பான `.credentials.json`-இல் `mcpOAuth`-இன் கீழ் வைத்திருக்கின்றன. ccam அவற்றை நகலெடுப்பதில்லை. நகலெடுத்த உள்நுழைவு மூலத்தின் அதே refresh token-ஐக் கொண்டிருக்கும், முதல் புதுப்பிப்பில் சர்வர் அந்த token-ஐ மாற்றிவிடும்: அப்போது அதை வைத்திருக்கும் மற்ற அனைவரும், மூல `~/.claude` உட்பட, அந்தச் சர்வரிலிருந்து வெளியேற்றப்படுவார்கள். அதனால் `sign-ins` திட்ட வரி, மூலம் உள்நுழைந்துள்ள ஒவ்வொரு நகலெடுக்கப்பட்ட சர்வருக்கும், `/mcp` மூலம் தாங்களே உள்நுழைய வேண்டிய குழுவின் கணக்குகளைப் பெயரிடுகிறது, சுருக்கம் அதை மீண்டும் சொல்கிறது. token இல்லாத உள்ளீடு உள்நுழைவு அல்ல, அது பட்டியலிடப்படாது; ஏற்கெனவே சொந்த உள்நுழைவு உள்ள கணக்கும், குழு தன் வரையறையை வைத்துக்கொள்ளும் சர்வரும் பட்டியலிடப்படாது. எந்த உள்நுழைவுக் கோப்பும் எழுதப்படுவதில்லை.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- மூலம் படிக்கப்பட மட்டுமே செய்யும். அதிலுள்ள எந்தக் கோப்பும் நகர்த்தப்படவோ, இணைக்கப்படவோ, மாற்றப்படவோ மாட்டாது, மூல உள்நுழைவு ஒருபோதும் எழுதப்படாது, அது Keychain-இலிருந்து படிக்கப்படக்கூடிய macOS-இலும் கூட.
- `use mine`, ஓர் ஒன்றிணைப்பு, அல்லது ஏற்கெனவே உள்ள `mcp.json`-இல் புதிய MCP சர்வர்கள் மூலம் ஒரு குழுக் கோப்பு மாற்றப்படும் முன், குழுவின் நகல் `<group>/.ccam-backup/<time>/`-இல் சேமிக்கப்படும். கோப்புறைப் பெயர் ஒவ்வோர் இயக்கத்துக்கும் தனித்துவமானது, அதனால் ஒரே வினாடியில் நடக்கும் இரண்டு இயக்கங்கள் ஒன்றின் காப்புப்பிரதியை மற்றொன்று மேலெழுதாது. சுருக்கம் பாதையை அச்சிடும்.
- `--yes` உடன், ஒவ்வொரு முரண்பாடும் குழுவின் நகலை வைத்திருக்கும், அதற்குப் பதிலாக எப்படித் தேர்ந்தெடுப்பது என்று ccam சொல்லும்: அதே கட்டளையை `--yes` இல்லாமல் இயக்குங்கள்.
- எதுவும் மாறாது என்றால், ccam `Nothing to copy.` என்று அச்சிடும்; ஒவ்வொரு முரண்பாடும் குழுவின் நகலை வைத்திருந்தால், `Nothing to copy: the group keeps its own copies.`
- ஒரு நகல் பாதியில் தோல்வியடைந்தால், ccam `The source folder was not changed. Group files replaced so far are saved in` என்றும், அதைத் தொடர்ந்து காப்புப்பிரதிக் கோப்புறையையும் சொல்லி நின்றுவிடும். மீண்டும் இயக்கினால், முழுமையாக வந்துசேராத கோப்புகள் நகலெடுக்கப்படும்.
- பல குழுக்கள் MCP-ஐப் பகிரும்போது, `ccam doctor`, `ccam ext`, `ccam mcp list` பரிந்துரைக்கும் கட்டளைகள் `ccam migrate --mcp --group <a|b>` என இருக்கும்.
- மூலம் இல்லையென்றால் `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` என்று முடியும்.

## தொடர்புடையவை

- [MCP சர்வர்கள், அமைப்புகள், வரலாற்றைப் பகிர்தல்](/docs/sharing)
- [சரிபார்த்தலும் பழுதுநீக்கமும்](/features/doctor)
- [பல Claude Code கணக்குகள்](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
