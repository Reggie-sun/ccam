---
title: ಪ್ರತಿ ಖಾತೆಯ ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳು
description: ccam agents ಎಲ್ಲಾ ಖಾತೆಗಳ Claude Code ಹಿನ್ನೆಲೆ ಏಜೆಂಟ್‌ಗಳು ಮತ್ತು ತೆರೆದ ಸೆಷನ್‌ಗಳನ್ನು ಒಂದೇ ಟೇಬಲ್‌ನಲ್ಲಿ ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ; ccam attach ಒಂದನ್ನು ID ಮೂಲಕ ತೆರೆಯುತ್ತದೆ.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: "06281fb13119"
---

Claude Code ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳನ್ನು (`claude --bg`, agent view) ಖಾತೆವಾರು ಇಡುತ್ತದೆ: `claude agents` ತಾನು ಯಾವ ಖಾತೆಯಾಗಿ ರನ್ ಆಗುತ್ತದೆಯೋ ಆ ಖಾತೆಯ ಸೆಷನ್‌ಗಳನ್ನು ಮಾತ್ರ ತೋರಿಸುತ್ತದೆ. `ccam agents` ಎಲ್ಲಾ ಖಾತೆಗಳನ್ನು ಒಮ್ಮೆಗೇ ಕೇಳಿ ಅವನ್ನು ಒಟ್ಟಿಗೆ ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ, ಮತ್ತು `ccam attach <id|name>` ಒಂದು ಹಿನ್ನೆಲೆ ಸೆಷನ್ ಅನ್ನು ಅದನ್ನು ಆರಂಭಿಸಿದ ಖಾತೆಯಾಗಿ ತೆರೆಯುತ್ತದೆ.

## ಬಳಕೆ

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` ಎಂದರೆ ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ನ ಸ್ಥಿತಿ (`working`, `blocked`, `done`, `failed`, `stopped`) ಮತ್ತು ಅದರ ಪ್ರೊಸೆಸ್ ಏನು ಮಾಡುತ್ತಿದೆ (`busy`, `waiting`, `idle`); ತೆರೆದಿರುವ ಇಂಟರ್‌ಆಕ್ಟಿವ್ ಸೆಷನ್ `open` ಎಂದು ತೋರಿಸುತ್ತದೆ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `--all` | ಮುಗಿದ, ವಿಫಲವಾದ ಅಥವಾ ನಿಲ್ಲಿಸಿದ ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳನ್ನೂ ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ | ಆಫ್ |
| `--json` | ಸೆಷನ್‌ಗಳನ್ನು JSON ಆಗಿ ಮುದ್ರಿಸುತ್ತದೆ, ಪ್ರತಿಯೊಂದೂ ಅದು ಸೇರಿರುವ `account` ಜೊತೆ | ಆಫ್ |

`ccam attach <id|name> [claude args]` ID ನಂತರದ ಎಲ್ಲವನ್ನೂ, ಫ್ಲ್ಯಾಗ್‌ಗಳನ್ನೂ ಸೇರಿಸಿ, ಬದಲಾವಣೆ ಇಲ್ಲದೆ `claude attach` ಗೆ ರವಾನಿಸುತ್ತದೆ: `ccam attach 49ec --dry-run` `--dry-run` ಅನ್ನು Claude Code ಗೆ ನೀಡುತ್ತದೆ. ಬದಲಿಗೆ ಕಮಾಂಡ್ ಅನ್ನು ಮುದ್ರಿಸಲು, `--dry-run` ಅನ್ನು ID ಗಿಂತ ಮುಂದೆ ಹಾಕಿ (`ccam --dry-run attach 49ec` ಅಥವಾ `ccam attach --dry-run 49ec`).

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

1. ಪ್ರತಿ ಖಾತೆಗೂ, ccam ಖಾತೆಯ `CLAUDE_CONFIG_DIR` ಜೊತೆ `claude agents --json` ಅನ್ನು (ನೀಡಿದ್ದರೆ `--all` ಜೊತೆ) ರನ್ ಮಾಡುತ್ತದೆ, ಸಮಾನಾಂತರವಾಗಿ, ಪ್ರತಿಯೊಂದಕ್ಕೂ 15 ಸೆಕೆಂಡ್ ಟೈಮ್‌ಔಟ್‌ನೊಂದಿಗೆ. ಇದು ಖಾತೆಯ ಸ್ವಂತ ಸೆಷನ್ ಪಟ್ಟಿಯನ್ನು ಓದುತ್ತದೆ ಮತ್ತು Claude Code ನ ಹಿನ್ನೆಲೆ ಸೇವೆಯನ್ನು ಆರಂಭಿಸುವುದಿಲ್ಲ.
2. ಪ್ರತಿ ನಮೂದಿಗೂ ಖಾತೆಯ ಹೆಸರು ಸೇರುತ್ತದೆ, ಮತ್ತು ಪಟ್ಟಿಯನ್ನು ಆರಂಭದ ಸಮಯದ ಪ್ರಕಾರ, ಹೊಸದು ಮೊದಲು, ವಿಂಗಡಿಸಲಾಗುತ್ತದೆ. ಉತ್ತರಿಸದ ಖಾತೆಯನ್ನು ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಎರರ್‌ನಲ್ಲಿ ಹೆಸರಿಸಿ ಬಿಟ್ಟುಬಿಡಲಾಗುತ್ತದೆ, ಮತ್ತು ನಂತರ ccam ಕೋಡ್ 1 ರೊಂದಿಗೆ ನಿರ್ಗಮಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಸ್ಕ್ರಿಪ್ಟ್ ವೈಫಲ್ಯವನ್ನು ಖಾಲಿ ಪಟ್ಟಿ ಎಂದು ತಪ್ಪಾಗಿ ತಿಳಿಯುವುದಿಲ್ಲ.
3. `ccam attach` ಪ್ರತಿ ಖಾತೆಯ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ `jobs/<id>` ಗಾಗಿ ಹುಡುಕುತ್ತದೆ, ಅಲ್ಲಿ Claude Code ಹಿನ್ನೆಲೆ ಸೆಷನ್ ಅನ್ನು ಇಡುತ್ತದೆ, ಮತ್ತು `ccam <name>` Claude Code ಅನ್ನು ಆರಂಭಿಸುವಂತೆಯೇ, ಆ ಖಾತೆಯಾಗಿ `claude attach <id>` ಅನ್ನು ರನ್ ಮಾಡುತ್ತದೆ. ID ಪೂರ್ವಪ್ರತ್ಯಯವಲ್ಲದ ಪದವನ್ನು ಸೆಷನ್ ಹೆಸರುಗಳೊಂದಿಗೆ, ದೊಡ್ಡ ಅಥವಾ ಸಣ್ಣ ಅಕ್ಷರಗಳನ್ನು ಲೆಕ್ಕಿಸದೆ, ಹೊಂದಿಸಲಾಗುತ್ತದೆ; ID ಯಂತೆ ಕಾಣುವ ಪದವನ್ನು ಎಂದಿಗೂ ಹೆಸರು ಎಂದು ಪರಿಗಣಿಸಲಾಗುವುದಿಲ್ಲ.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- ಹಿನ್ನೆಲೆ ಸೆಷನ್ ID ಗಳು ಒಂದೇ ಖಾತೆಗೆ ಸೇರಿರುತ್ತವೆ: ಇನ್ನೊಂದು ಖಾತೆಯಾಗಿ ರನ್ ಮಾಡಿದ `claude attach` ಅವನ್ನು ಕಂಡುಹಿಡಿಯಲಾರದು, ಮತ್ತು `ccam attach` ಬಗೆಹರಿಸುವುದು ಇದನ್ನೇ.
- ಎರಡು ಖಾತೆಗಳ ಸೆಷನ್‌ಗಳು ಒಂದೇ ID ಪೂರ್ವಪ್ರತ್ಯಯ (prefix) ಅಥವಾ ಹೆಸರಿನ ಒಂದೇ ಪದವನ್ನು ಹೊಂದಿದ್ದರೆ, `ccam attach` ಅವುಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡಿ ಇನ್ನಷ್ಟು ಉದ್ದದ ID ಅಥವಾ ಹೆಸರನ್ನು ಕೇಳುತ್ತದೆ.
- ccam ತನ್ನ ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳನ್ನು Claude Code ಗೆ ಒಂದೇ ಆರ್ಗ್ಯುಮೆಂಟ್‌ನಲ್ಲಿ `--mcp-config=<file>` ಮೂಲಕ ನೀಡುತ್ತದೆ, ಅದನ್ನು `claude agents` ಮತ್ತು `claude attach` ಸ್ವೀಕರಿಸುತ್ತವೆ; ಹಿನ್ನೆಲೆ ಸೆಷನ್‌ಗಳು ಮರುಆರಂಭವಾದಾಗ ಅದೇ ಫೈಲ್ ಅನ್ನು ಬಳಸುತ್ತಲೇ ಇರುತ್ತವೆ.
- API ಕೀ ಖಾತೆಯ ಹಿನ್ನೆಲೆ ಸೇವೆ ತನ್ನ ಕೀಯನ್ನು ಅದನ್ನು ಆರಂಭಿಸಿದ ಶೆಲ್‌ನಿಂದ ಪಡೆಯುತ್ತದೆ. ಅಂತಹ ಖಾತೆಯ ಮೊದಲ ಹಿನ್ನೆಲೆ ಸೆಷನ್ ಅನ್ನು ccam ಮೂಲಕ ಆರಂಭಿಸಿ (`ccam gw --bg …`), ಹೀಗಾಗಿ ಸೇವೆಗೆ ಕೀ ಸಿಗುತ್ತದೆ.

## ಸಂಬಂಧಿತ

- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳಿಗಾಗಿ ಖಾತೆ ಆಯ್ಕೆ](/features/pick)
- [ಡ್ಯಾಶ್‌ಬೋರ್ಡ್](/features/dashboard)
