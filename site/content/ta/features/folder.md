---
title: கோப்புறை வாரியான கணக்குகள்
description: கணக்குப் பெயரின்றி Claude Code-ஐத் தொடங்கினால், ccam folder-ஆல் கோப்புறையும் உட்கோப்புறைகளும் ஒரே கணக்கைப் பயன்படுத்தும், எ.கா. நிறுவனக் குறியீட்டுக்கு work.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` கணக்குப் பெயர் இல்லாத ஒவ்வொரு தொடக்கத்தையும், `~/code/company`-க்குள்ளும் அதன் கீழ் உள்ள எந்தக் கோப்புறையிலும், Claude Code-ஐ `work` ஆகத் தொடங்க வைக்கிறது. அதற்கு வெளியே, இயல்புநிலைக் கணக்கு முன்பு போலவே தொடங்கும். `ccam personal` மூலமோ `ccpersonal` போன்ற குறுக்குவழி மூலமோ ஒரு கணக்கைப் பெயரிட்டால், எப்போதும் அந்தக் கணக்கே தொடங்கும்.

## பயன்பாடு

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

இந்தத் தொடக்கங்கள் கோப்புறையின் கணக்கைப் பயன்படுத்துகின்றன:

- `ccam -p "…"` மற்றும் முதல் argument Claude Code-க்கான option ஆக உள்ள வேறு எந்தத் தொடக்கமும்
- பெயர் இல்லாத `ccam run`, மற்றும் `ccam --dry-run`
- `ccam --fallback -p "…"` அதன் முதல் முயற்சிக்கு
- தேர்வி, தனியாக `ccam`, இது கோப்புறையின் கணக்கின் மீது cursor உடன் தொடங்குகிறது

## Options

| Option | என்ன செய்யும் | இயல்புநிலை |
| --- | --- | --- |
| `<account>` | பயன்படுத்த வேண்டிய கணக்கின் பெயர் அல்லது எண் | |
| `[folder]` | கோப்புறை, relative அல்லது absolute | தற்போதைய கோப்புறை |
| `--remove` | சரியாக அந்தக் கோப்புறைக்கான பதிவை நீக்கும் | off |
| `--dry-run` | மதிக்கப்படும்: மாற்றத்தைச் சேமிக்காமல் அச்சிடும் | off |

## இது எப்படி வேலை செய்கிறது

ஒவ்வொரு கோப்புறையும் `~/.config/ccam/config.toml`-இல் ஒரு `[[folder]]` பதிவு; ஒவ்வொரு தொடக்கமும் எப்படியும் படிக்கும் கோப்பு இது:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` கோப்புறையை, symbolic links தீர்க்கப்பட்ட absolute பாதையாக மாற்றிச் சேமிக்கிறது. ஒரு கோப்புறையை மீண்டும் அமைத்தால் அதன் கணக்கு மாற்றப்படும்.
2. பெயர் இல்லாத தொடக்கம் `launcher.HereProfile`-ஐ அழைக்கிறது. பதிவுகள் இல்லையென்றால் அது உடனே இயல்புநிலைக் கணக்கை எடுக்கிறது.
3. இல்லையெனில் அது working directory-ஐ ஒவ்வொரு `path`-உடனும் strings ஆக ஒப்பிடுகிறது: கோப்புறையே அல்லது அதன் கீழ் உள்ள ஒன்று, அதே எழுத்துகளில் தொடங்கும் கோப்புறை மட்டும் ஒருபோதும் அல்ல (`/code/company` என்பது `/code/companyx`-ஐ உள்ளடக்காது). மிக ஆழமான பதிவு வெல்லும், அதனால் `~/code/company` `work` ஆக இருக்கும்போது `~/code` `personal` ஆக இருக்கலாம்.
4. working directory அப்படியே எதனுடனும் பொருந்தாதபோது மட்டுமே, ccam அதன் symbolic links-ஐத் தீர்த்து மீண்டும் ஒருமுறை ஒப்பிடுகிறது.
5. Windows-இல், Windows செய்வது போலவே, ஒப்பீடு பெரிய, சிறிய எழுத்து வேறுபாட்டைப் புறக்கணிக்கிறது.

`config.toml` தவிர வேறு எந்தக் கோப்பும் படிக்கப்படுவதில்லை, அதனால் ஒரு கோப்புறைக்குள் தொடங்குவதற்கு அளவிடக்கூடிய நேரம் எதுவும் ஆகாது (பதிவுகளுடனோ இல்லாமலோ 4.0 ms, `hyperfine`, 300 இயக்கங்கள்).

## பாதுகாப்பும் விளிம்பு நிலைகளும்

- `ccam rename` கோப்புறைப் பதிவுகளைப் புதிய பெயருக்குக் கொண்டு செல்கிறது, `ccam remove` நீக்கப்பட்ட கணக்கின் பதிவுகளை நீக்குகிறது.
- கணக்கு இனி இல்லாத ஒரு பதிவு, அந்தக் கோப்புறையில் பெயர் இல்லாத தொடக்கத்தை நிறுத்துகிறது: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` அத்தகைய பதிவுகளையும், கோப்புறை இல்லாமல் போன பதிவுகளையும் தெரிவிக்கிறது.
- மேலே உள்ள கோப்புறையிலிருந்து கணக்கைப் பெறும் கோப்புறைக்குள் `ccam folder --remove` எதையும் நீக்குவதற்குப் பதிலாக அந்தக் கோப்புறையின் பெயரைச் சொல்கிறது.
- `ccam next` சுழற்சி அல்லது பயன்பாட்டின்படி தேர்ந்தெடுக்கிறது, கோப்புறைகளைப் பார்ப்பதில்லை.

## தொடர்புடையவை

- [Claude Code-ஐத் தொடங்குதல்](/features/launch)
- [குறுக்குவழிக் கட்டளைகள்](/features/shortcuts)
- [கணக்குகளை நிர்வகித்தல்](/features/accounts)
- [கட்டளைக் குறிப்பு](/docs/commands)
