---
title: ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು
description: ccam shortcut ಮೂಲಕ ಖಾತೆಗೆ ccwork ನಂತಹ ತನ್ನದೇ ಕಮಾಂಡ್ ನೀಡಿ, ಮತ್ತು ccam shims install ಮೂಲಕ ~/.local/bin ನಲ್ಲಿರುವ ಲಿಂಕ್‌ಗಳನ್ನು ಮತ್ತೆ ರಚಿಸಿ.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

ಶಾರ್ಟ್‌ಕಟ್ ಎಂದರೆ ಒಂದು ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುವ, ನೀವೇ ಹೆಸರಿಡುವ ಕಮಾಂಡ್: `ccam work --resume` ಮಾಡುವುದನ್ನೇ `ccwork --resume` ಮಾಡುತ್ತದೆ. ಇದು `~/.local/bin` ನಲ್ಲಿರುವ ccam ಗೆ ಒಂದು ಲಿಂಕ್, ಹೀಗಾಗಿ bash, zsh, ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳು ಮತ್ತು `PATH` ನಿಂದ ಪ್ರೋಗ್ರಾಂಗಳನ್ನು ರನ್ ಮಾಡುವ ಬೇರೆ ಯಾವುದರಲ್ಲೂ ಕೆಲಸ ಮಾಡುತ್ತದೆ. ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು ಐಚ್ಛಿಕ; ನೀವು ಆಗಾಗ ಆರಂಭಿಸುವ ಖಾತೆಗಳಿಗೆ ಒಂದನ್ನು ರಚಿಸಿ.

## ಬಳಕೆ

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam ಪ್ರತಿ ಬದಲಾವಣೆಯನ್ನು ದೃಢೀಕರಿಸುತ್ತದೆ:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

ಒಂದು ಖಾತೆಗೆ ಗರಿಷ್ಠ ಒಂದು ಶಾರ್ಟ್‌ಕಟ್ ಇರುತ್ತದೆ, ಅದು `~/.config/ccam/config.toml` ನಲ್ಲಿ `shortcut` ಆಗಿ ಸಂಗ್ರಹವಾಗುತ್ತದೆ. `ccam add` ವಿಜಾರ್ಡ್ `Create a shortcut command for <name>? [y/N]` ಎಂದು ಕೇಳುತ್ತದೆ, ನಂತರ `Shortcut name [cc<name>]:` ಎಂದು ಕೇಳುತ್ತದೆ. ವಿಜಾರ್ಡ್‌ನ ಹೊರಗೆ, `ccam add` `--shortcut` ಇದ್ದಾಗ ಮಾತ್ರ ಒಂದನ್ನು ರಚಿಸುತ್ತದೆ. `ccam list` ಪ್ರತಿ ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು, ಇಲ್ಲದಿದ್ದರೆ `-` ಅನ್ನು ತೋರಿಸುತ್ತದೆ.

ಶಾರ್ಟ್‌ಕಟ್ ಪ್ರತಿ ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು Claude Code ಗೆ ರವಾನಿಸುತ್ತದೆ, `ccam <name>` ಹೆಸರಿನ ನಂತರ ಮಾಡುವಂತೆಯೇ. ಇದರಲ್ಲಿ `--dry-run` ಮತ್ತು `--help` ಕೂಡ ಸೇರಿವೆ: ಆರಂಭವನ್ನು ಮುಂಚಿತವಾಗಿ ನೋಡಲು `ccam --dry-run work` ಬಳಸಿ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `--remove` | ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಮತ್ತು ಅದರ ಲಿಂಕ್ ಅನ್ನು ಅಳಿಸುತ್ತದೆ | ಆಫ್ |
| `-h`, `--help` | ಕಮಾಂಡ್‌ಗೆ ಸಹಾಯ | |
| `--dry-run` | `ccam shortcut` ಇದನ್ನು ಬೆಂಬಲಿಸುತ್ತದೆ: ಹೆಸರನ್ನು ಪರಿಶೀಲಿಸಿ, ಉಳಿಸದೆ ಅಥವಾ ಲಿಂಕ್ ಮಾಡದೆ `Would set the shortcut of work to ccwork.` ಅಥವಾ `Would remove ccwork, the shortcut of work.` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ | ಆಫ್ |

`ccam shims install` ಯಾವುದೇ ಆರ್ಗ್ಯುಮೆಂಟ್ ತೆಗೆದುಕೊಳ್ಳುವುದಿಲ್ಲ, ಮತ್ತು `--dry-run` ಅನ್ನು `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` ಸಂದೇಶದೊಂದಿಗೆ ನಿರಾಕರಿಸುತ್ತದೆ.

ಹೆಸರಿನ ನಿಯಮಗಳು, ಈ ಕ್ರಮದಲ್ಲಿ ಪರಿಶೀಲಿಸಲಾಗುತ್ತದೆ:

- `a-z`, `0-9`, `.`, `-` ಮತ್ತು `_` ಅಕ್ಷರಗಳ 1 ರಿಂದ 64 ಅಕ್ಷರಗಳು, ಅಕ್ಷರ ಅಥವಾ ಅಂಕಿಯಿಂದ ಆರಂಭವಾಗಬೇಕು;
- `ccam`, `claude` ಅಥವಾ `run`, `usage` ನಂತಹ ccam ಕಮಾಂಡ್ ಆಗಿರಬಾರದು;
- ಇನ್ನೊಂದು ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಆಗಿರಬಾರದು;
- ಈಗಾಗಲೇ ccam ಗೆ ಲಿಂಕ್ ಅಲ್ಲದ, `~/.local/bin` ನಲ್ಲಿರುವ ಫೈಲ್ ಆಗಿರಬಾರದು;
- ನಿಮ್ಮ `PATH` ನಲ್ಲಿರುವ ಇನ್ನೊಂದು ಎಕ್ಸಿಕ್ಯೂಟಬಲ್‌ನ ಹೆಸರಾಗಿರಬಾರದು.

ಕಮಾಂಡ್ ಮತ್ತು `--remove` ಎರಡನ್ನೂ ಒಟ್ಟಿಗೆ ನೀಡಿದರೆ `give a command or --remove, not both` ದೋಷದೊಂದಿಗೆ ನಿಲ್ಲುತ್ತದೆ. ಖಾತೆಗೆ ಈಗಾಗಲೇ ಇರುವ ಹೆಸರನ್ನೇ ಹೊಂದಿಸಿದರೆ `work: ccwork (unchanged)` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ, ಮತ್ತು ಅದರ ಲಿಂಕ್ ಕಾಣೆಯಾಗಿದ್ದರೆ ಮತ್ತೆ ರಚಿಸುತ್ತದೆ.

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

ಹೆಸರನ್ನು ಹೊಂದಿಸಿದಾಗ `checkShortcut` ರನ್ ಆಗುತ್ತದೆ. ಅದು bin ಫೋಲ್ಡರ್ ಮತ್ತು ccam ಬೈನರಿಯ ನಿಜವಾದ ಪಾತ್‌ನೊಂದಿಗೆ `shim.Check` ಅನ್ನು ಕರೆಯುತ್ತದೆ:

<!-- motion:shortcuts-check -->

ಹೆಸರು ಪರಿಶೀಲನೆಗಳ ಮೂಲಕ ಕ್ರಮವಾಗಿ ಕೆಳಗೆ ಸಾಗುತ್ತದೆ. ಪ್ರತಿ ಪರಿಶೀಲನೆಯ ಪಕ್ಕದಲ್ಲಿ ಅದು ನಿರಾಕರಿಸುವ ಒಂದು ಹೆಸರಿದೆ: `Work` ಮಾದರಿಗೆ ಹೊಂದುವುದಿಲ್ಲ, `run` ಒಂದು ccam ಕಮಾಂಡ್, `ccside` ಇನ್ನೊಂದು ಖಾತೆಗೆ ಸೇರಿದೆ, `ccfile` `~/.local/bin` ನಲ್ಲಿ ಈಗಾಗಲೇ ಇರುವ ಫೈಲ್, ಮತ್ತು `git` `/usr/bin/git` ಅನ್ನು ಮರೆಮಾಡುತ್ತದೆ. ಎಲ್ಲಾ ಪರಿಶೀಲನೆಗಳನ್ನು ದಾಟಿದ ಹೆಸರನ್ನು ಮಾತ್ರ ಉಳಿಸಿ ಲಿಂಕ್ ಮಾಡಲಾಗುತ್ತದೆ.

1. `shim.Check` ಹೆಸರನ್ನು `^[a-z0-9][a-z0-9._-]{0,63}$` ನೊಂದಿಗೆ, ನಂತರ `config.Reserved` ನೊಂದಿಗೆ, ನಂತರ `cfg.FindShortcut` ನೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ. ಅದು `~/.local/bin/<name>` ಮೇಲೆ `os.Lstat` ಕರೆಯುತ್ತದೆ, ಮತ್ತು `PATH` ನಲ್ಲಿರುವ ಪ್ರತಿ ಇತರ ಫೋಲ್ಡರ್‌ಗೆ `os.Stat` ಕರೆದು, ccam ಗೆ ಲಿಂಕ್ ಅಲ್ಲದ ಎಕ್ಸಿಕ್ಯೂಟಬಲ್ ಫೈಲ್ ಅನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ.
2. ccam `config.toml` ಅನ್ನು `fsutil.AtomicWrite` ಮೂಲಕ ಉಳಿಸುತ್ತದೆ (ತಾತ್ಕಾಲಿಕ ಫೈಲ್, `fsync`, ಮರುಹೆಸರಿಸುವಿಕೆ).
3. `shim.Remove` ಹಳೆಯ ಶಾರ್ಟ್‌ಕಟ್‌ನ ಲಿಂಕ್ ಅನ್ನು ಅಳಿಸುತ್ತದೆ, ಆದರೆ ಅದು ccam ಗೆ ಲಿಂಕ್ ಆಗಿದ್ದರೆ ಮಾತ್ರ.
4. `shim.Install` ಕಾನ್ಫಿಗರ್ ಮಾಡಿದ, ಕಾಣೆಯಾಗಿರುವ ಪ್ರತಿ ಶಾರ್ಟ್‌ಕಟ್‌ನ ಲಿಂಕ್ ಅನ್ನು `os.Symlink` ಮೂಲಕ ccam ಬೈನರಿಗೆ ರಚಿಸುತ್ತದೆ; ಬೈನರಿಯ ಪಾತ್ ಅನ್ನು `os.Executable` ಮತ್ತು `filepath.EvalSymlinks` ಮೂಲಕ ಕಂಡುಹಿಡಿಯಲಾಗುತ್ತದೆ. ಇದು ಎಂದಿಗೂ ಏನನ್ನೂ ತೆಗೆದುಹಾಕುವುದಿಲ್ಲ.

`os.Readlink` ccam ಪಾತ್ ಅನ್ನು ಹಿಂತಿರುಗಿಸಿದಾಗ, ಅಥವಾ ಲಿಂಕ್ ಮತ್ತು ccam ಪಾತ್ ಒಂದೇ ಫೈಲ್‌ಗೆ ತಲುಪಿದಾಗ, ಆ ಲಿಂಕ್ ಅನ್ನು ccam ನದೇ ಎಂದು ಪರಿಗಣಿಸಲಾಗುತ್ತದೆ. ಪ್ರತಿ ಅಳಿಸುವಿಕೆ ಮತ್ತು ಪ್ರತಿ ಮೇಲ್ಬರಹ ಪರಿಶೀಲನೆ ಈ ಪರೀಕ್ಷೆಯ ಮೂಲಕ ಹೋಗುತ್ತದೆ, ಹೀಗಾಗಿ ccam ತಾನು ರಚಿಸದ ಫೈಲ್‌ಗಳನ್ನು ಮುಟ್ಟುವುದಿಲ್ಲ.

ನೀವು ಶಾರ್ಟ್‌ಕಟ್ ರನ್ ಮಾಡಿದಾಗ, ನೀವು ಟೈಪ್ ಮಾಡಿದ ಹೆಸರು ಮಾತ್ರ ccam ಗೆ ಬೇಕಾಗುತ್ತದೆ:

<!-- motion:shortcuts-link -->

ಶೆಲ್ `PATH` ನಲ್ಲಿ `ccwork` ಅನ್ನು ಕಂಡುಹಿಡಿಯುತ್ತದೆ, ಕರ್ನಲ್ ಲಿಂಕ್ ಅನ್ನು ಅನುಸರಿಸಿ ccam ಬೈನರಿಯನ್ನು ಆರಂಭಿಸುತ್ತದೆ, ಮತ್ತು `argv[0]` ಇನ್ನೂ `ccwork` ಆಗಿರುತ್ತದೆ. `app.Main` ನಲ್ಲಿ, `ccam` ಅಲ್ಲದ ಯಾವುದೇ ಹೆಸರು `config.toml` ಅನ್ನು ಒಮ್ಮೆ ಲೋಡ್ ಮಾಡಿ `shim.Lookup` ಅನ್ನು ಕರೆಯುತ್ತದೆ, ಅದು `shortcut` ಮೌಲ್ಯ `filepath.Base(argv[0])` ಗೆ ಸಮನಾಗಿರುವ ಖಾತೆಯನ್ನು ಕಂಡುಹಿಡಿಯುತ್ತದೆ. ನಂತರ `launcher.LaunchWith` ಈಗಾಗಲೇ ಲೋಡ್ ಆದ ಕಾನ್ಫಿಗ್‌ನೊಂದಿಗೆ [`ccam <name>`](/features/launch) ನ ಅದೇ ಹಂತಗಳನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ, ಮತ್ತು `argv[1:]` ಅನ್ನು ಬದಲಾಯಿಸದೆ `claude` ಗೆ ರವಾನಿಸುತ್ತದೆ.

ಶೆಲ್ alias ಅಥವಾ ವ್ರ್ಯಾಪರ್ ಸ್ಕ್ರಿಪ್ಟ್ ಬದಲು symlink ಏಕೆ: symlink ಪ್ರತಿ ಶೆಲ್‌ನಲ್ಲೂ ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳಲ್ಲೂ ಕೆಲಸ ಮಾಡುತ್ತದೆ, ಹೆಚ್ಚುವರಿ ಶೆಲ್ ಆರಂಭವನ್ನು (ಸುಮಾರು 1 ರಿಂದ 2 ms) ಸೇರಿಸುವುದಿಲ್ಲ, ಮತ್ತು ಆರಂಭವು ಈಗಲೂ ಒಂದೇ ಫೈಲ್ ಅನ್ನು ಓದುತ್ತದೆ. `scripts/perf-check.sh` ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು `hyperfine` ಮೂಲಕ, `ccam <name>` ನ ಅದೇ 10 ms (p95) ಮಿತಿಗೆ ಎದುರಾಗಿ ಅಳೆಯುತ್ತದೆ.

### Windows ನಲ್ಲಿ

ಶಾರ್ಟ್‌ಕಟ್ ಎಂದರೆ `%USERPROFILE%\.local\bin` ನಲ್ಲಿರುವ `ccwork.exe`, ಮತ್ತು `shim.Name` ಹುಡುಕುವ ಮೊದಲು `.exe` ಅನ್ನು (ಯಾವುದೇ ಕೇಸ್‌ನಲ್ಲಿ) ತೆಗೆದುಹಾಕುತ್ತದೆ, ಹೀಗಾಗಿ `ccwork`, `ccwork.exe` ಮತ್ತು `CCWORK.EXE` ಎಲ್ಲವೂ `work` ಅನ್ನು ಆರಂಭಿಸುತ್ತವೆ. ccam ಮೊದಲು symlink ರಚಿಸುತ್ತದೆ, ಇದನ್ನು Windows Developer Mode ನಲ್ಲಿ ಅಥವಾ ನಿರ್ವಾಹಕರಿಗೆ ಅನುಮತಿಸುತ್ತದೆ. ಆ ಹಕ್ಕು ಇಲ್ಲದಿದ್ದರೆ ಅದು hard link ರಚಿಸುತ್ತದೆ: ಅದೇ `ccam.exe` ಫೈಲ್‌ಗೆ ಎರಡನೇ ಹೆಸರು, ಇದಕ್ಕೆ ಯಾವ ವಿಶೇಷ ಹಕ್ಕೂ ಬೇಕಿಲ್ಲ. hard link `ccam.exe` ನ ಅದೇ ಫೈಲ್ ಆಗಿದ್ದಾಗ ಅದನ್ನು ccam ನದೇ ಎಂದು ಎಣಿಸಲಾಗುತ್ತದೆ. `ccam update` `ccam.exe` ನ ಜಾಗದಲ್ಲಿ ಹೊಸ ಫೈಲ್ ಇಡುವುದರಿಂದ, ಅದು ನಂತರ `shim.Relink` ಮೂಲಕ hard link ಆಗಿರುವ ಪ್ರತಿ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು ನವೀಕರಿಸಿ `✓ shortcut ccwork renewed` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ. ಹೆಸರಿನ ಪರಿಶೀಲನೆ ಇತರ `PATH` ಫೋಲ್ಡರ್‌ಗಳಲ್ಲಿ `PATHEXT` ನ ಪ್ರತಿ ವಿಸ್ತರಣೆಯೊಂದಿಗೆ (`.exe`, `.cmd` ಇತ್ಯಾದಿ) `<name>` ಗಾಗಿ ಹುಡುಕುತ್ತದೆ.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- **ಯಾವ ಖಾತೆಗೂ ಸೇರದ ಲಿಂಕ್.** `~/.local/bin` ನಲ್ಲಿರುವ ccam ಗೆ ಇರುವ ಲಿಂಕ್‌ನ ಹೆಸರನ್ನು ಯಾವ ಖಾತೆಯೂ ಬಳಸದಿದ್ದರೆ, ಉದಾಹರಣೆಗೆ ನೀವು `config.toml` ನಲ್ಲಿ `shortcut` ಅನ್ನು ಕೈಯಿಂದ ಬದಲಾಯಿಸಿದ ನಂತರ, ccam ಡೀಫಾಲ್ಟ್ ಖಾತೆಗೆ ಹಿಂತಿರುಗುವುದಿಲ್ಲ. ಅದು `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` ದೋಷದೊಂದಿಗೆ ನಿಲ್ಲುತ್ತದೆ. `ccam shortcut` ಮೂಲಕ ಹೆಸರನ್ನು ಮತ್ತೆ ಹೊಂದಿಸಿ, ಅಥವಾ ಲಿಂಕ್ ಅನ್ನು ಅಳಿಸಿ.
- **ಈಗಿರುವ ಫೈಲ್‌ಗಳು.** ccam ತಾನು ರಚಿಸದ ಫೈಲ್ ಅನ್ನು ಎಂದಿಗೂ ಬದಲಾಯಿಸುವುದಿಲ್ಲ. `ccam shims install` `<path> already exists and is not a ccam link, so it was left alone` ಸಂದೇಶದೊಂದಿಗೆ ನಿಲ್ಲುತ್ತದೆ.
- **ಸ್ಥಳಾಂತರಿಸಿದ ccam ಬೈನರಿ.** ಆಗ ಹಳೆಯ ಲಿಂಕ್‌ಗಳು ಇನ್ನು ಅಸ್ತಿತ್ವದಲ್ಲಿಲ್ಲದ ಪಾತ್‌ಗೆ ಸೂಚಿಸುತ್ತವೆ, ಮತ್ತು ccam ಅವುಗಳನ್ನು ತನ್ನವು ಎಂದು ಪರಿಗಣಿಸುವುದಿಲ್ಲ. ಅವುಗಳನ್ನು ಅಳಿಸಿ, ನಂತರ `ccam shims install` ರನ್ ಮಾಡಿ.
- **ಮರುಹೆಸರಿಸುವುದು ಮತ್ತು ತೆಗೆದುಹಾಕುವುದು.** `ccam rename` ಶಾರ್ಟ್‌ಕಟ್ ಹೆಸರನ್ನು ಹಾಗೆಯೇ ಇಡುತ್ತದೆ, ಮತ್ತು ಲಿಂಕ್ ಕೆಲಸ ಮಾಡುತ್ತಲೇ ಇರುತ್ತದೆ, ಏಕೆಂದರೆ ಅದು ccam ಗೆ ಸೂಚಿಸುತ್ತದೆ ಮತ್ತು ccam ಹೆಸರನ್ನು `config.toml` ನಲ್ಲಿ ಹುಡುಕುತ್ತದೆ. `ccam remove` ಖಾತೆಯ ಲಿಂಕ್ ಅನ್ನು ಅಳಿಸುತ್ತದೆ. `ccam shortcut <account> <new>` ಹೊಸ ಲಿಂಕ್ ರಚಿಸುವ ಮೊದಲು ಹಳೆಯದನ್ನು ಅಳಿಸುತ್ತದೆ.
- **PATH.** ccam ಲಿಂಕ್ ರಚಿಸಿದಾಗ `~/.local/bin` ನಿಮ್ಮ `PATH` ನಲ್ಲಿ ಇಲ್ಲದಿದ್ದರೆ, ಅದು `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` ಎಂದು ಮುದ್ರಿಸಿ, ನಂತರ ಸೇರಿಸಬೇಕಾದ `export PATH=...` ಸಾಲನ್ನು ತೋರಿಸುತ್ತದೆ. Windows ನಲ್ಲಿ ಬದಲಿಗೆ ಆ ಫೋಲ್ಡರ್ ಅನ್ನು ನಿಮ್ಮ ಬಳಕೆದಾರ `PATH` ಗೆ ಸೇರಿಸುವ PowerShell ಸಾಲನ್ನು ಮುದ್ರಿಸುತ್ತದೆ.
- **ನಂತರ ಪರಿಶೀಲಿಸುವುದು.** `ccam doctor` `missing shortcuts in /home/me/.local/bin: ccwork` ಎಂದು ವರದಿ ಮಾಡುತ್ತದೆ, ಮತ್ತು ಫೋಲ್ಡರ್ `PATH` ನಲ್ಲಿ ಇಲ್ಲದಿದ್ದಾಗ ಎಚ್ಚರಿಸುತ್ತದೆ. `ccam doctor --fix` ಕಾಣೆಯಾದ ಲಿಂಕ್‌ಗಳನ್ನು ರಚಿಸುತ್ತದೆ.

## ಸಂಬಂಧಿತ

- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಖಾತೆಗಳ ನಿರ್ವಹಣೆ](/features/accounts)
- [ಕಮಾಂಡ್ ಉಲ್ಲೇಖ](/docs/commands#shortcuts)
- [ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ](/docs/install)
