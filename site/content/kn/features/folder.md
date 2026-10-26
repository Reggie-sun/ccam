---
title: ಫೋಲ್ಡರ್‌ವಾರು ಖಾತೆಗಳು
description: ಖಾತೆಯ ಹೆಸರಿಲ್ಲದೆ Claude Code ಆರಂಭಿಸಿದಾಗ, ccam folder ಒಂದು ಫೋಲ್ಡರ್ ಮತ್ತು ಅದರ ಕೆಳಗಿನ ಎಲ್ಲವೂ ಒಂದೇ ಖಾತೆಯನ್ನು ಬಳಸುವಂತೆ ಮಾಡುತ್ತದೆ, ಉದಾಹರಣೆಗೆ ಕಂಪನಿಯ ಕೋಡ್‌ಗೆ work.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: "25f82da64580"
---

`ccam folder work ~/code/company` ನಂತರ, `~/code/company` ಅಥವಾ ಅದರ ಉಪಫೋಲ್ಡರ್‌ಗಳಲ್ಲಿ ಖಾತೆಯ ಹೆಸರಿಲ್ಲದ ಯಾವುದೇ ಆರಂಭವು Claude Code ಅನ್ನು `work` ಆಗಿ ಆರಂಭಿಸುತ್ತದೆ. ಅದರ ಹೊರಗೆ, ಡೀಫಾಲ್ಟ್ ಖಾತೆ ಮೊದಲಿನಂತೆಯೇ ಆರಂಭವಾಗುತ್ತದೆ. `ccam personal` ಅಥವಾ `ccpersonal` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಮೂಲಕ ಖಾತೆಯನ್ನು ಹೆಸರಿಸಿದರೆ, ಯಾವಾಗಲೂ ಆ ಖಾತೆಯೇ ಆರಂಭವಾಗುತ್ತದೆ.

## ಬಳಕೆ

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

ಈ ಆರಂಭಗಳು ಫೋಲ್ಡರ್‌ನ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತವೆ:

- `ccam -p "…"` ಮತ್ತು ಮೊದಲ ಆರ್ಗ್ಯುಮೆಂಟ್ Claude Code ಗಾಗಿ ಒಂದು ಆಯ್ಕೆಯಾಗಿರುವ ಇತರ ಯಾವುದೇ ಆರಂಭ
- ಹೆಸರಿಲ್ಲದ `ccam run`, ಮತ್ತು `ccam --dry-run`
- `ccam --fallback -p "…"`, ಅದರ ಮೊದಲ ಪ್ರಯತ್ನಕ್ಕೆ
- ಡ್ಯಾಶ್‌ಬೋರ್ಡ್, ಅಂದರೆ ಕೇವಲ `ccam`, ಅದು ಫೋಲ್ಡರ್‌ನ ಖಾತೆಯ ಮೇಲೆ ಕರ್ಸರ್ ಇಟ್ಟು ಆರಂಭವಾಗುತ್ತದೆ, ನೀವು ಸಿಂಬಾಲಿಕ್ ಲಿಂಕ್ ಮೂಲಕ ಫೋಲ್ಡರ್ ತಲುಪಿದ್ದರೂ ಸಹ

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `<account>` | ಬಳಸಬೇಕಾದ ಖಾತೆಯ ಹೆಸರು ಅಥವಾ ಸಂಖ್ಯೆ | |
| `[folder]` | ಫೋಲ್ಡರ್, ಸಾಪೇಕ್ಷ ಅಥವಾ ಸಂಪೂರ್ಣ ಪಾತ್ | ಈಗಿನ ಫೋಲ್ಡರ್ |
| `--remove` | ನಿಖರವಾಗಿ ಆ ಫೋಲ್ಡರ್‌ನ ನಮೂದನ್ನು ತೆಗೆದುಹಾಕುತ್ತದೆ | ಆಫ್ |
| `--dry-run` | ಬೆಂಬಲಿತ: ಬದಲಾವಣೆಯನ್ನು ಉಳಿಸದೆ ಮುದ್ರಿಸುತ್ತದೆ | ಆಫ್ |

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

ಪ್ರತಿ ಫೋಲ್ಡರ್ `~/.config/ccam/config.toml` ನಲ್ಲಿ ಒಂದು `[[folder]]` ನಮೂದು; ಪ್ರತಿ ಆರಂಭವೂ ಹೇಗಿದ್ದರೂ ಓದುವ ಫೈಲ್ ಇದೇ:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` ಫೋಲ್ಡರ್ ಅನ್ನು ಸಿಂಬಾಲಿಕ್ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಹರಿಸಿದ ಸಂಪೂರ್ಣ ಪಾತ್ ಆಗಿ ಮಾಡಿ ಉಳಿಸುತ್ತದೆ. ಫೋಲ್ಡರ್ ಅನ್ನು ಮತ್ತೆ ಹೊಂದಿಸಿದರೆ ಅದರ ಖಾತೆ ಬದಲಾಗುತ್ತದೆ.
2. ಹೆಸರಿಲ್ಲದ ಆರಂಭ `launcher.HereProfile` ಅನ್ನು ಕರೆಯುತ್ತದೆ. ಯಾವ ನಮೂದೂ ಇಲ್ಲದಿದ್ದರೆ ಅದು ತಕ್ಷಣ ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ.
3. ಇಲ್ಲದಿದ್ದರೆ ಅದು ಕೆಲಸದ ಡೈರೆಕ್ಟರಿಯನ್ನು ಪ್ರತಿ `path` ಜೊತೆ ಸ್ಟ್ರಿಂಗ್‌ಗಳಾಗಿ ಹೋಲಿಸುತ್ತದೆ: ಫೋಲ್ಡರ್ ಸ್ವತಃ ಅಥವಾ ಅದರ ಕೆಳಗಿನದು, ಅದೇ ಅಕ್ಷರಗಳಿಂದ ಆರಂಭವಾಗುವುದಷ್ಟೇ ಆದ ಫೋಲ್ಡರ್ ಎಂದಿಗೂ ಅಲ್ಲ (`/code/company` `/code/companyx` ಅನ್ನು ಒಳಗೊಳ್ಳುವುದಿಲ್ಲ). ಅತ್ಯಂತ ಆಳದ ನಮೂದು ಗೆಲ್ಲುತ್ತದೆ, ಹೀಗಾಗಿ `~/code/company` `work` ಆಗಿರುವಾಗ `~/code` `personal` ಆಗಿರಬಹುದು.
4. ಕೆಲಸದ ಡೈರೆಕ್ಟರಿ ನೀಡಿದಂತೆಯೇ ಯಾವುದಕ್ಕೂ ಹೊಂದದಿದ್ದರೆ ಮಾತ್ರ, ccam ಅದರ ಸಿಂಬಾಲಿಕ್ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಹರಿಸಿ ಮತ್ತೆ ಹೋಲಿಸುತ್ತದೆ.
5. Windows ಮತ್ತು macOS ನಲ್ಲಿ, ಅವುಗಳ ಫೈಲ್ ಸಿಸ್ಟಮ್‌ಗಳು ಸಾಮಾನ್ಯವಾಗಿ ಮಾಡುವಂತೆಯೇ, ಹೋಲಿಕೆ ದೊಡ್ಡ ಮತ್ತು ಸಣ್ಣ ಅಕ್ಷರಗಳನ್ನು ಲೆಕ್ಕಿಸುವುದಿಲ್ಲ.

`config.toml` ಹೊರತುಪಡಿಸಿ ಬೇರೆ ಯಾವ ಫೈಲ್ ಅನ್ನೂ ಓದಲಾಗುವುದಿಲ್ಲ, ಹೀಗಾಗಿ ಫೋಲ್ಡರ್‌ನೊಳಗಿನ ಆರಂಭಕ್ಕೆ ಅಳೆಯಬಹುದಾದ ಯಾವುದೇ ಸಮಯ ಸೇರುವುದಿಲ್ಲ (ನಮೂದುಗಳಿದ್ದರೂ ಇಲ್ಲದಿದ್ದರೂ 4.0 ms, `hyperfine`, 300 ರನ್‌ಗಳು).

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- `ccam rename` ಫೋಲ್ಡರ್ ನಮೂದುಗಳನ್ನು ಹೊಸ ಹೆಸರಿಗೆ ಒಯ್ಯುತ್ತದೆ, ಮತ್ತು `ccam remove` ತೆಗೆದುಹಾಕಿದ ಖಾತೆಯ ನಮೂದುಗಳನ್ನು ಅಳಿಸುತ್ತದೆ.
- ಖಾತೆ ಇನ್ನು ಅಸ್ತಿತ್ವದಲ್ಲಿಲ್ಲದ ನಮೂದು ಆ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಹೆಸರಿಲ್ಲದ ಆರಂಭವನ್ನು ನಿಲ್ಲಿಸುತ್ತದೆ: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` ಅಂತಹ ನಮೂದುಗಳನ್ನು, ಮತ್ತು ಫೋಲ್ಡರ್ ಇಲ್ಲವಾದ ನಮೂದುಗಳನ್ನು ವರದಿ ಮಾಡುತ್ತದೆ.
- ಮೇಲಿನ ಫೋಲ್ಡರ್‌ನಿಂದ ಖಾತೆಯನ್ನು ಪಡೆಯುವ ಉಪಫೋಲ್ಡರ್‌ನಲ್ಲಿ ನೀವು `ccam folder --remove` ರನ್ ಮಾಡಿದರೆ, ccam ಆ ಮೇಲಿನ ಫೋಲ್ಡರ್ ಅನ್ನು ಹೆಸರಿಸುತ್ತದೆ ಮತ್ತು ಏನನ್ನೂ ತೆಗೆದುಹಾಕುವುದಿಲ್ಲ.
- `ccam next` ರೊಟೇಶನ್ ಅಥವಾ ಬಳಕೆಯ ಆಧಾರದಲ್ಲಿ ಆಯ್ಕೆಮಾಡುತ್ತದೆ ಮತ್ತು ಫೋಲ್ಡರ್‌ಗಳನ್ನು ನೋಡುವುದಿಲ್ಲ.

## ಸಂಬಂಧಿತ

- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು](/features/shortcuts)
- [ಖಾತೆಗಳ ನಿರ್ವಹಣೆ](/features/accounts)
- [ಕಮಾಂಡ್ ಉಲ್ಲೇಖ](/docs/commands)
