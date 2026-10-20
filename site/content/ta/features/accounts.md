---
title: கணக்குகளை நிர்வகித்தல்
description: ccam add, ccam login, ccam list, ccam rename, ccam remove, ccam import மூலம் Claude Code கணக்குகளைச் சேர்க்க, உள்நுழைய, பட்டியலிட, மறுபெயரிட, நீக்க.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

ஒவ்வொரு ccam கணக்கும் தனக்கென உள்நுழைவு கொண்ட ஒரு தனி Claude Code config கோப்புறை. இந்தக் கட்டளைகள் கணக்குகளை உருவாக்குகின்றன, உள்நுழைகின்றன, பட்டியலிடுகின்றன, பெயர் மாற்றுகின்றன, நீக்குகின்றன, உங்களிடம் ஏற்கெனவே உள்ள கோப்புறைகளைப் பதிவு செய்கின்றன.

## பயன்பாடு

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

ஒரு கட்டளை `<name|number>`-ஐ ஏற்கும் இடங்களில், கணக்கின் பெயரையோ அதன் slot எண்ணையோ, அதாவது `ccam list`-இன் `#` நெடுவரிசையை, கொடுக்கலாம்.

options இல்லாமல் டெர்மினலில் இயக்கும் `ccam add` ஒரு wizard-ஐத் திறக்கும். அது பெயரையும் கணக்கு எப்படி உள்நுழையும் என்பதையும் (இப்போதே Claude உள்நுழைவு, API key, அல்லது பின்னர் Claude உள்நுழைவு) கேட்டு, பிறகு சேர வேண்டிய குழுக்களை `default` தேர்வுசெய்யப்பட்ட நிலையில் காட்டி (குழுக்கள் இல்லாதபோது, `default` குழு பகிரப்போகும் உருப்படிகளை), பிறகு `Create a shortcut command for work? [y/N]` என்று கேட்டு, உறுதிப்படுத்த ஒரு சுருக்கத்துடன் முடியும். `Esc` ஒரு படி பின்செல்லும். `--group`, `--no-group`, `--shortcut`, `--no-login` ஆகியவற்றில் ஏதேனும் ஒன்று, அல்லது டெர்மினலுக்குப் பதிலாக ஒரு pipe, wizard-ஐத் தவிர்க்கும்; அப்போது பெயர் கட்டாயம், `--shortcut` ஒன்றைப் பெயரிடாவிட்டால் குறுக்குவழி உருவாக்கப்படாது.

`--group` இல்லாவிட்டால், புதிய கணக்கு `default` குழுவில் சேரும். முதல் கணக்கு அதை உருவாக்கும், அது MCP சர்வர்கள், அமைப்புகள், `CLAUDE.md`, நீட்டிப்புகள் ஆகியவற்றைப் பகிரும்; [குழுக்களுடன் பகிர்தல்](/features/share) பாருங்கள்.

`ccam list` slot, பெயர், குறுக்குவழி, மின்னஞ்சல், குழுக்கள், கடைசிப் பயன்பாடு ஆகியவற்றைக் காட்டும், இயல்புநிலைக் கணக்குக்கு முன் `*` இருக்கும்:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

`--json` உடன், ஒவ்வொரு கணக்கும் `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (கணக்கு ஒருபோதும் தொடங்கப்படவில்லை என்றால் விடப்படும்), `groups`, `share` ஆகியவற்றைக் கொண்ட ஒரு object. `share` ஒவ்வோர் உருப்படியையும் (`mcp`, `settings`, `memory`, `extensions`, `history`) அது வரும் குழுவுடன், அல்லது `"local"` உடன் இணைக்கும்.

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `ccam add --no-login` | உள்நுழையாமல் கணக்கை உருவாக்கும்; பின்னர் `ccam login <name>` மூலம் உள்நுழையுங்கள் | off |
| `ccam add --group <group>` | `default`-க்குப் பதிலாக இந்தக் குழுவில் சேரும். பலமுறை கொடுக்கலாம்; `default` இல்லையென்றால் உருவாக்கப்படும் | `default` |
| `ccam add --no-group` | எந்தக் குழுவிலும் சேராது: கணக்கு எதையும் பகிராது | off |
| `ccam add --shortcut <command>` | இந்தக் குறுக்குவழிக் கட்டளையையும் `~/.local/bin`-இல் உருவாக்கும் | none |
| `ccam add --api --key-from <source>` | உள்நுழைவுக்குப் பதிலாக API key பயன்படுத்தும் கணக்கு; [API key கணக்குகள்](/features/api) பாருங்கள் | off |
| `ccam list --json` | கணக்குகளை JSON ஆக அச்சிடும் | off |
| `ccam remove -y`, `--yes` | கேட்காமல் நீக்கும் | off |
| `ccam import --name <name>` | நீங்கள் கொடுக்கும் கோப்புறைக்கான கணக்குப் பெயர் | கோப்புறைப் பெயரிலிருந்து பெறப்படும் |
| `ccam import -y`, `--yes` | கண்டுபிடித்த ஒவ்வொரு கோப்புறையையும் கேட்காமல் இறக்குமதி செய்யும் | off |
| `--dry-run` | எதையும் மாற்றாத `ccam list` மற்றும் `ccam completion` இதை மதிக்கும். `add`, `login`, `logout`, `use`, `rename`, `remove`, `import` இதை `` --dry-run is not supported by `ccam add`; nothing was changed `` (கட்டளையைப் பெயரிட்டு) என்று மறுக்கும் | off |

`ccam login`, `ccam logout`, `ccam use`, `ccam rename` ஆகியவற்றுக்குச் சொந்த options இல்லை. `ccam completion` `bash`, `zsh` அல்லது `powershell`-ஐ ஏற்கும்; வேறு எதுவும் `only bash, zsh and powershell are supported` என்று தோல்வியடையும்.

கணக்குப் பெயர்கள் `a-z`, `0-9`, `-`, `_` ஆகியவற்றால் ஆன 1 முதல் 32 எழுத்துகள், ஓர் எழுத்து அல்லது இலக்கத்தில் தொடங்க வேண்டும். `ccam 2` என்றால் slot 2 என்பதால், பெயர் முழுவதும் இலக்கங்களாக இருக்க முடியாது. பெயர் ஒரு ccam கட்டளையாகவும் இருக்க முடியாது (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), அதனால் `ccam <name>` ஒருபோதும் குழப்பமாக இருக்காது. இந்தப் பட்டியல் `internal/config`-இல் உள்ள `config.Reserved`.

## இது எப்படி வேலை செய்கிறது

`ccam add work --shortcut ccwork` `addProfile`-இல் (`internal/cli/profiles.go`) இந்தப் படிகளை இயக்குகிறது:

<!-- motion:accounts-add -->

1. எதுவும் உருவாக்கப்படும் முன், `config.ValidateName` பெயரையும், `shim.Check` குறுக்குவழியையும் சரிபார்க்கின்றன.
2. `profile.Add` `0700` mode-உடன் `~/.local/share/ccam/profiles/work/`-ஐ உருவாக்கி, கணக்குக்கு மிகக் குறைந்த காலி slot எண்ணைக் (`Config.NextSlot`) கொடுக்கிறது, அதனால் நீக்கப்பட்ட கணக்கின் எண் மீண்டும் பயன்படுத்தப்படும். முதல் கணக்கு இயல்புநிலை ஆகிறது.
3. கணக்கு தன் குழுக்களில் சேர்கிறது (`--group` அல்லது `--no-group` வேறு எதுவும் சொல்லாவிட்டால் `default`; முதல் கணக்கு `default`-ஐ உருவாக்குகிறது), `applyShare` புதிய, காலியான கோப்புறையை `share.Build` மற்றும் பிறகு `Apply` மூலம் அவற்றுடன் இணைக்கிறது, எடுத்துக்காட்டாக `settings.json`-ஐ `groups/default/settings.json`-உடன். இது தோல்வியடைந்தால், புதிய கோப்புறை நீக்கப்படும், ஏனெனில் மீதமிருக்கும் கோப்புறை ஒவ்வொரு மறுமுயற்சியையும் `folder already exists` என்று தடுக்கும்.
4. கணக்கு `fsutil.AtomicWrite` மூலம் `~/.config/ccam/config.toml`-இல் சேமிக்கப்படுகிறது: ஒரு தற்காலிகக் கோப்பு, `fsync`, பிறகு `rename`, அதனால் crash ஒருபோதும் பாதி config-ஐ விட்டுச்செல்லாது.
5. `runLogin` கணக்குக் கோப்புறைக்கு அமைக்கப்பட்ட `CLAUDE_CONFIG_DIR` உடனும், உள்நுழைவுப் பக்கம் நிரப்பப்பட்டிருக்கக் கணக்கின் கடைசி மின்னஞ்சலுடன் `--email` உடனும், `claude auth login`-ஐத் தொடங்குகிறது. உள்நுழைவு தோல்வியடைந்தால், கணக்கு அப்படியே இருக்கும், ccam பிழையை அச்சிட்டு, பின்னர் `ccam login work`-ஐ இயக்கச் சொல்லும்.
6. `profile.ReadIdentity` `profiles/work/.claude.json`-இலிருந்து `oauthAccount`-ஐ எடுக்கிறது: மின்னஞ்சல், கணக்கு மற்றும் நிறுவன ids, நிறுவனத்தின் பெயர், plan (`claude_max`, `claude_team` …). அந்தக் கோப்பு பல megabytes ஆக இருக்கலாம், அதனால் `TopLevelValue` முழுக் கோப்பையும் parse செய்வதற்குப் பதிலாக அந்த ஒரு key-ஐ மட்டும் தேடுகிறது. ccam அவற்றை `config.toml`-இல் `[profile.login]`-இல் வைத்திருக்கிறது.
7. `shim.Install` `~/.local/bin/ccwork`-ஐ ccam binary-உடன் இணைக்கிறது; அது தொடங்கும்போது `argv[0]`-ஐ `config.toml`-இல் தேடுகிறது.

உள்நுழைவுக்குச் சில நிமிடங்கள் ஆகலாம், இடையில் மற்றொரு ccam config-ஐச் சேமிக்கக்கூடும். அதனால் காத்திருக்கும் கட்டளைகள் `env.update` மூலம் சேமிக்கின்றன; அது `config.toml`-ஐ மீண்டும் ஏற்றி, அந்த ஒரு மாற்றத்தை (இங்கே, மின்னஞ்சல்) பொருத்தி, சேமிக்கிறது. இதே காரணத்துக்காக `ccam remove` தன் கேள்விக்குப் பிறகு மீண்டும் ஏற்றுகிறது.

`ccam login` ஏற்கெனவே உள்ள கணக்குக்குப் படி 5 மற்றும் 6-ஐ இயக்குகிறது. Claude Code மின்னஞ்சலை உள்நுழைவுப் பக்கத்துக்கான ஒரு குறிப்பாக மட்டுமே எடுத்துக்கொள்கிறது, அதனால் ccam கணக்கு முன்பு இருந்த Claude கணக்கைத் தவிர வேறொரு Claude கணக்கில் உள்நுழைவு முடியக்கூடும்; ccam புதிய உள்நுழைவை `[profile.login]`-உடன் ஒப்பிட்டு அதைச் சொல்கிறது (`ccam: this login is other@company.com, but work was dev@company.com before.`). அது புதிய உள்நுழைவை வைத்துக்கொண்டு, அதை எப்படித் திரும்பப் பெறுவது என்று சொல்கிறது: `ccam logout work`, பிறகு `ccam login work`, சரியான மின்னஞ்சலுடன் உள்நுழையுங்கள். வேறொரு ccam கணக்கு அதே Claude கணக்கில் உள்நுழைந்திருந்தால், அந்த இரண்டும் தங்கள் பயன்பாட்டு வரம்புகளைப் பகிர்ந்துகொள்கின்றன என்று ccam எச்சரிக்கிறது. `ccam login --expired` ஒருபோதும் உள்நுழையாத, அல்லது உள்நுழைவு (அதன் refresh token) காலாவதியான ஒவ்வொரு கணக்கிலும், ஒன்றன்பின் ஒன்றாக உள்நுழைகிறது.

`ccam logout` கணக்கில் `claude auth logout`-ஐ இயக்கி, சேமித்த மின்னஞ்சலை அழிக்கிறது; உள்நுழைவுப் பதிவு அப்படியே இருக்கும், அதனால் அடுத்த உள்நுழைவு அதனுடன் சரிபார்க்கப்படும். `ccam use` `config.toml`-இல் `default`-ஐ மட்டுமே மாற்றி எழுதுகிறது.

கணக்குகள் ஒன்றுக்கு மேற்பட்ட நிறுவனங்களைச் சேர்ந்தவையாக இருக்கும்போது `ccam list` ஒரு `Organization` நெடுவரிசையைச் சேர்க்கிறது (தனிப்பட்ட plans `personal` என்று காட்டும்), `ccam list --json`-இல் ஒவ்வொரு கணக்குக்கும் `plan`, `org` இருக்கும்.

### தொலைநிலைக் கணினியில் உள்நுழைதல்

SSH வழியாக, `claude auth login`-இன் browser பகுதி உடைந்துவிடக்கூடும். Claude Code சர்வரில் `http://localhost:<port>/callback`-இல் உள்ள ஒரு callback சர்வரில் உள்நுழைவுக்காகக் காத்திருக்கிறது, ஆனால் உள்நுழைவை முடிக்கும் browser உங்கள் கணினியில் இயங்குகிறது, அங்கே அந்த முகவரி எங்கும் செல்லாது, browser `ERR_CONNECTION_REFUSED` போன்ற பிழையைக் காட்டும். முகவரியில் இன்னும் ஒருமுறைக் குறியீடு உள்ளது, அதனால் ccam அதைச் சேர்க்க முடியும்:

<!-- motion:accounts-relay -->

`runLogin` உங்கள் டெர்மினலை Claude Code-இடம் ஒப்படைப்பதில்லை. அது Claude Code-இன் standard input-ஐ ஒரு pipe-உடன் இணைக்கிறது, நீங்கள் தட்டச்சு செய்யும் ஒவ்வொரு வரியையும் `relayLogin` அதற்குள் நகலெடுக்கிறது, அதனால் Claude Code-இன் `Paste code here` கேள்வியில் ஒட்டப்படும் குறியீடு இன்னும் வந்துசேரும். ஒரு வரி callback முகவரியாக இருந்தால், host `localhost`, `127.0.0.1` அல்லது `[::1]` ஆக இருந்து, port கொடுக்கப்பட்டு, பாதை `/callback` ஆக இருந்து, `code` மற்றும் `state` இரண்டும் இருந்தால் மட்டுமே `callbackURL` அதை ஏற்கிறது; `http://`-ஐ விட்டுவிடலாம். பிறகு `deliverCallback` browser செய்திருப்பது போலவே, 10 வினாடி வரம்புடனும் redirects-ஐப் பின்தொடராமலும், சர்வரில் அந்த முகவரியைக் கோருகிறது. கோரிக்கை காத்திருக்கும் `claude auth login`-இன் callback சர்வரை அடைகிறது (வரைபடத்தில் கோடிட்ட வரி). Claude Code குறியீட்டைப் பெற்று உள்நுழைவைச் சேமிக்கிறது, Linux-இல் கணக்குக் கோப்புறையில் உள்ள `.credentials.json`-இல். ஒட்டப்பட்ட வரி ccam-ஐ ஒருபோதும் வேறொரு host-ஐத் தொடர்புகொள்ள வைக்க முடியாது.

உள்நுழைவு தொடங்கும் முன் ccam இதை அச்சிடும்:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

வெற்றிபெற்றால் அது `→ Passed the sign-in to Claude Code on this machine.` என்று அச்சிடும்.

உள்நுழைவு டெர்மினலில் உள்ளீட்டுக்குக் காத்திருக்கும்போது, ccam டெர்மினலின் echo மற்றும் line mode-ஐ அணைத்து, நீங்கள் தட்டச்சு செய்யும் ஒவ்வொரு எழுத்துக்கும் ஒரு `*`-ஐக் காட்டுகிறது, அதனால் ஒட்டப்படும் குறியீடோ callback முகவரியோ திரையில் தெரியாது. `readLine` வரித் திருத்தத்தைத் தானே செய்கிறது: `Enter` வரியை அனுப்புகிறது, `Backspace` ஓர் எழுத்தை அழிக்கிறது, வெற்று வரியில் `Ctrl+D` உள்ளீட்டை முடிக்கிறது. உள்நுழைவு முடியும்போதும், `Ctrl+C`, `SIGTERM` அல்லது `SIGHUP` வரும்போதும் டெர்மினல் மீட்டமைக்கப்படும்; இவை இன்னும் ccam-ஐ முடிக்கும். standard input மற்றும் standard error-இல் டெர்மினல் இல்லாவிட்டால், வரிகள் வந்தபடியே நகலெடுக்கப்படும்.

### பெயர் மாற்றம், நீக்கம், இறக்குமதி

`ccam rename` மற்றும் `ccam remove` கணக்கின் கோப்புறையிலும், கணக்கின் பெயரில் ccam வேறிடங்களில் வைத்திருப்பவற்றிலும் செயல்படுகின்றன. `nameState` அந்தப் பாதைகளைப் பட்டியலிடுகிறது: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, usage caches `~/.local/state/ccam/usage/<name>.json` மற்றும் `<name>.statusline.json`. `ccam remove work` அவற்றை இந்த வரிசையில் கையாள்கிறது:

<!-- motion:accounts-remove -->

`ccam remove` முதலில் கணக்கு இயங்கிக்கொண்டிருந்தால் மறுத்து, பிறகு கேட்கும். கணக்கு நீட்டிப்புகளைப் பகிர்ந்தால், `RewritePluginPaths` குழுவின் `installed_plugins.json` மற்றும் `known_marketplaces.json`-ஐ மாற்றி எழுதுகிறது, ஏனெனில் Claude Code இந்தக் கணக்கின் கோப்புறை வழியாகச் செல்லக்கூடிய absolute plugin பாதைகளைச் சேமிக்கிறது. `profile.Trash` கோப்புறையை `~/.local/share/ccam/trash/<name>-<time>/`-க்கு நகர்த்துகிறது, கணக்கு `config.toml`-இலிருந்து நீங்குகிறது, அதே பெயருடைய பிந்தைய கணக்கு அவற்றைப் பெறாமல் இருக்க `nameState` பாதைகள் நீக்கப்படுகின்றன, குறுக்குவழி இணைப்பும் நீக்கப்படுகிறது. மீட்டெடுக்க, கோப்புறையைக் குப்பையிலிருந்து திரும்ப நகர்த்துங்கள்.

`ccam rename` கோப்புறையின் பெயரை மாற்றி, அதன் `plugins/`-இல் உள்ள plugin பாதைகளைப் பழைய கோப்புறையிலிருந்து புதியதுக்கு மாற்றி எழுதி, இயல்புநிலையையும் உடன் கொண்டு சென்று, `nameState` பாதைகளைப் புதிய பெயருக்கு நகர்த்துகிறது. குறுக்குவழிக் கட்டளை தன் பெயரை வைத்திருக்கும். plugin மாற்றி எழுதுதல் தோல்வியடைந்தால், config சேமிக்கப்படாததால், கோப்புறையின் பெயர் மீண்டும் பழையதாக மாற்றப்படும்.

`ccam import` ஒரு கோப்புறையை அது இருக்கும் இடத்திலேயே, `config.toml`-இல் `dir` ஆகப் பதிவு செய்கிறது. அது அதை ஒருபோதும் நகர்த்துவதோ நகலெடுப்பதோ இல்லை: நகலெடுக்கப்பட்ட refresh token நகல்களை ஒன்றையொன்று வெளியேற்ற வைக்கும். கோப்புறை argument இல்லாமல், அது பதிவுசெய்யப்படாத `~/.claude` மற்றும் `~/.claude-*` கோப்புறைகளைத் தேடி, ஒவ்வொன்றைப் பற்றியும் கேட்கும்; பெயர்கள் கோப்புறையிலிருந்து வருகின்றன (`~/.claude` என்பது `main`, `~/.claude-work` என்பது `work`, `~/.claude-2` என்பது `acct2`), பெயரைப் பெற முடியாத கோப்புறை பயன்படுத்த வேண்டிய `--name` கட்டளையுடன் பட்டியலிடப்படும். `CLAUDE_CONFIG_DIR` இல்லாமல், Claude Code `.claude.json`-ஐ உங்கள் home கோப்புறையில் `~/.claude`-க்கு அருகில் வைக்கிறது, அதனால் நீங்கள் `~/.claude`-ஐ இறக்குமதி செய்து அதற்குச் சொந்த `.claude.json` இல்லையென்றால், ccam `~/.claude/.claude.json`-இலிருந்து `~/.claude.json`-க்கு இணைப்பை உருவாக்கும். அப்போது வெறும் `claude`-உம் `ccam main`-உம் ஒரே கணக்கைப் பார்க்கும். நீங்கள் ஒரு குழுவில் சேர்க்கும் வரை (`ccam group join <group> <name>`) இறக்குமதி செய்யப்பட்ட கணக்குகள் எந்தக் குழுவிலும் இருக்காது, குறுக்குவழியும் உருவாக்கப்படாது.

`ccam completion bash|zsh` கணக்குப் பெயர்களையும் நிறைவுசெய்யும் ஒரு completion script-ஐ அச்சிடும்.

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- கணக்கு இயங்கும்போது `ccam rename` மற்றும் `ccam remove` மறுக்கும் (`refuseRunning`), `work is running; quit its Claude Code session first` என்று, ஏனெனில் அமர்வு கோப்புறையில் தொடர்ந்து எழுதுகிறது.
- இறக்குமதி செய்யப்பட்ட கணக்கில் `ccam remove` அதன் பதிவை மட்டுமே நீக்கும்: `✓ main unregistered (folder ... left in place)`. உங்கள் `~/.claude` குப்பைக்கு நகர்த்தப்படாது.
- `profiles/`-இல் ஏற்கெனவே கோப்புறை உள்ள பெயரை `ccam add` மறுத்து, `ccam import`-ஐச் சுட்டிக்காட்டும்.
- `ccam import` ஏற்கெனவே பதிவுசெய்யப்பட்ட கோப்புறை, கோப்புறை அல்லாத பாதை, ஏற்கெனவே எடுக்கப்பட்ட பெயர் ஆகியவற்றை மறுக்கும்.
- `ccam logout` மற்றும் `ccam login` கணக்கு இயங்குகிறதா என்று சரிபார்ப்பதில்லை.
- relay Claude Code-ஐ அடைய முடியாவிட்டால், எடுத்துக்காட்டாக உள்நுழைவு முடிந்த பிறகு, அது காரணத்துடன் `ccam: could not pass the sign-in to Claude Code:` என்று அச்சிடும். HTTP பிழை பொதுவாகக் குறியீடு காலாவதியானது என்று பொருள்; உள்நுழைவை மீண்டும் தொடங்குங்கள்.

## தொடர்புடையவை

- [Claude Code-ஐத் தொடங்குதல்](/features/launch)
- [குறுக்குவழிக் கட்டளைகள்](/features/shortcuts)
- [MCP சர்வர்கள், அமைப்புகள், வரலாற்றைப் பகிர்தல்](/docs/sharing)
- [பல Claude Code கணக்குகள்](/guides/multiple-claude-code-accounts)
