---
title: ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ
metaTitle: Linux, macOS ಮತ್ತು Windows ನಲ್ಲಿ ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ
description: Linux, macOS ಅಥವಾ Windows ನಲ್ಲಿ ಒಂದೇ ಕಮಾಂಡ್, .deb ಪ್ಯಾಕೇಜ್ ಅಥವಾ go install ಮೂಲಕ ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ, ccam update ಮೂಲಕ ಅಪ್‌ಡೇಟ್ ಮಾಡಿ. ಮೊದಲು Claude Code ಬೇಕು.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "2ada7b2650d6"
---

ccam ನಿಮಗಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಆದ್ದರಿಂದ ಮೊದಲು [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ ಮತ್ತು `claude --version` ಕೆಲಸ ಮಾಡುತ್ತದೆಯೇ ಎಂದು ಖಚಿತಪಡಿಸಿಕೊಳ್ಳಿ.

Linux ಅಥವಾ macOS ನ ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಇದನ್ನು ರನ್ ಮಾಡಿ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ಈ ಸ್ಕ್ರಿಪ್ಟ್ ನಿಮ್ಮ ಸಿಸ್ಟಮ್‌ಗೆ ಹೊಂದುವ ಬಿಡುಗಡೆಯನ್ನು ಡೌನ್‌ಲೋಡ್ ಮಾಡುತ್ತದೆ, ಅದರ SHA-256 ಚೆಕ್‌ಸಮ್ ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು `ccam` ಅನ್ನು `~/.local/bin` ನಲ್ಲಿ ಇರಿಸುತ್ತದೆ; ನೀವು ನಂತರ ಆಯ್ಕೆಮಾಡುವ `ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳನ್ನೂ ccam ಅದೇ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ರಚಿಸುತ್ತದೆ. ಆ ಫೋಲ್ಡರ್ ನಿಮ್ಮ `PATH` ನಲ್ಲಿ ಇಲ್ಲದಿದ್ದರೆ, `~/.bashrc` ಅಥವಾ `~/.zshrc` ಗೆ ಸೇರಿಸಬೇಕಾದ ಸಾಲನ್ನು ಸ್ಕ್ರಿಪ್ಟ್ ಮುದ್ರಿಸುತ್ತದೆ.

Windows ನಲ್ಲಿ, ಇದನ್ನು PowerShell ನಲ್ಲಿ ರನ್ ಮಾಡಿ:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

ಇದು `ccam_windows_<arch>.zip` ಅನ್ನು ಡೌನ್‌ಲೋಡ್ ಮಾಡುತ್ತದೆ, ಅದರ SHA-256 ಚೆಕ್‌ಸಮ್ ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು `ccam.exe` ಅನ್ನು `%USERPROFILE%\.local\bin` ನಲ್ಲಿ ಇರಿಸುತ್ತದೆ (Claude Code ನ ಇನ್‌ಸ್ಟಾಲರ್ `claude.exe` ಅನ್ನು ಇರಿಸುವ ಫೋಲ್ಡರ್ ಇದೇ). ಆ ಫೋಲ್ಡರ್ ನಿಮ್ಮ ಬಳಕೆದಾರ `PATH` ನಲ್ಲಿ ಇನ್ನೂ ಇಲ್ಲದಿದ್ದರೆ, ಸ್ಕ್ರಿಪ್ಟ್ ಅದನ್ನು ಸೇರಿಸುತ್ತದೆ ಮತ್ತು `%USERPROFILE%\bin` ನಂತಹ ನಮೂದುಗಳನ್ನು ಇದ್ದಂತೆಯೇ ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ. ನಂತರ ಹೊಸ ಟರ್ಮಿನಲ್ ತೆರೆಯಿರಿ.

## ಬೆಂಬಲಿತ ಸಿಸ್ಟಮ್‌ಗಳು

- x86-64 ಮತ್ತು ARM64 ನಲ್ಲಿ Linux: Ubuntu, Debian, Fedora, Arch, Alpine ಮತ್ತು ಇತರ ಡಿಸ್ಟ್ರಿಬ್ಯೂಷನ್‌ಗಳು. ಬೈನರಿ ಸ್ಟ್ಯಾಟಿಕ್ ಆಗಿದೆ ಮತ್ತು ಅದಕ್ಕೆ ಯಾವುದೇ ಡಿಪೆಂಡೆನ್ಸಿ ಇಲ್ಲ.
- Apple silicon ಮತ್ತು Intel ನಲ್ಲಿ macOS.
- x86-64 ಮತ್ತು ARM64 ನಲ್ಲಿ Windows 10 ಮತ್ತು 11, PowerShell, Windows Terminal ಅಥವಾ Git Bash ನಲ್ಲಿ:
  - ಹಂಚಿಕೊಂಡ ಫೋಲ್ಡರ್‌ಗಳು Developer Mode ಇಲ್ಲದೆಯೂ ಕೆಲಸ ಮಾಡುತ್ತವೆ.
  - ಹಂಚಿಕೊಂಡ ಫೈಲ್‌ಗಳಿಗೆ (`settings.json`, `CLAUDE.md`, `history.jsonl`) Developer Mode (Settings > System > For developers) ಬೇಕು, ಏಕೆಂದರೆ Windows ಆ ಮೋಡ್‌ಗೆ ಅಥವಾ ನಿರ್ವಾಹಕರಿಗೆ ಮಾತ್ರ ಫೈಲ್‌ಗಳಿಗೆ ಲಿಂಕ್ ರಚಿಸಲು ಬಿಡುತ್ತದೆ; [Windows ನಲ್ಲಿ ಹಂಚಿಕೆ](/features/share#on-windows) ನೋಡಿ.
  - Developer Mode ಇಲ್ಲದಿದ್ದರೆ, ಮೊದಲ `ccam add` MCP ಸರ್ವರ್‌ಗಳು ಮತ್ತು ವಿಸ್ತರಣೆಗಳನ್ನು ಮಾತ್ರ ಹೊಂದಿರುವ `default` ಗುಂಪನ್ನು ರಚಿಸುತ್ತದೆ, ಮತ್ತು Developer Mode ಆನ್ ಆದ ನಂತರ ಉಳಿದವನ್ನು ಹೇಗೆ ಸೇರಿಸುವುದು ಎಂದು ತಿಳಿಸುತ್ತದೆ.
- WSL ಒಳಗೆ ccam ಬೇರೆ ಯಾವುದೇ Linux ಸಿಸ್ಟಮ್‌ನಲ್ಲಿ ಕೆಲಸ ಮಾಡುವಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತದೆ.

## ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡುವ ಇತರ ವಿಧಾನಗಳು

**Debian ಮತ್ತು Ubuntu ಪ್ಯಾಕೇಜ್.** [ಬಿಡುಗಡೆಗಳ ಪುಟದಿಂದ](https://github.com/jhyunwoo/ccam/releases) `ccam_<version>_linux_amd64.deb` ಡೌನ್‌ಲೋಡ್ ಮಾಡಿ ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

ಪ್ಯಾಕೇಜ್ `/usr/bin/ccam` ಅನ್ನು ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡುತ್ತದೆ.

**ಸೋರ್ಸ್‌ನಿಂದ.** Go 1.26 ಅಥವಾ ಹೊಸ ಆವೃತ್ತಿಯೊಂದಿಗೆ:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

ನೀವು ಬಿಡುಗಡೆಯ ಟ್ಯಾಗ್ (`@latest` ಅಥವಾ `@vX.Y.Z`) ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದರೆ, ccam ಆ ಆವೃತ್ತಿಯನ್ನು ತಿಳಿಸುತ್ತದೆ, ಮತ್ತು `ccam update` ಹಾಗೂ ಡ್ಯಾಶ್‌ಬೋರ್ಡ್‌ನ ಅಪ್‌ಡೇಟ್ ಆಫರ್ ಕೆಲಸ ಮಾಡುತ್ತವೆ. `@main` ನ ಬಿಲ್ಡ್, ಅಥವಾ ಟ್ಯಾಗ್‌ನ ನಂತರ commit‌ಗಳಿರುವ ಸ್ಥಳೀಯ checkout ನ ಬಿಲ್ಡ್ `dev` ಎಂದು ತಿಳಿಸುತ್ತದೆ, ಮತ್ತು ನೀವು `--version` ನೀಡದ ಹೊರತು `ccam update` ಅದನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ.

**ನಿರ್ದಿಷ್ಟ ಆವೃತ್ತಿ.** `CCAM_VERSION` ಹೊಂದಿಸಿ:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## ಶೆಲ್ completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion ಖಾತೆಗಳ ಹೆಸರುಗಳನ್ನೂ ಒಳಗೊಂಡಂತೆ `ccam` ನ ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು ಮತ್ತು ಆಯ್ಕೆಗಳಿಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ.

## Claude Code ಪ್ಲಗಿನ್

ccam ಪ್ಲಗಿನ್ Claude Code ಗೇ ccam ಬಳಸುವುದನ್ನು ಕಲಿಸುತ್ತದೆ: ಯಾವ ಖಾತೆಯಲ್ಲಿ ಅವಕಾಶ ಇದೆ ಎಂದು ಪರಿಶೀಲಿಸಲು, `ccam pick` ಮೂಲಕ ಬೇರೆ ಬೇರೆ ಖಾತೆಗಳಲ್ಲಿ ಸಮಾನಾಂತರ `claude -p` ವರ್ಕರ್‌ಗಳನ್ನು ಆರಂಭಿಸಲು, ಮತ್ತು ಕಾಣೆಯಾದ ಸ್ಕಿಲ್ ಅಥವಾ MCP ಸರ್ವರ್ ಅನ್ನು `ccam ext` ಮೂಲಕ ವಿವರಿಸಲು. ಇದನ್ನು Claude Code ಒಳಗೆ ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

ಇದು ಸೇರಿಸುವುದು:

- ಲಾಗಿನ್ ಮಾಡುವ ಅಥವಾ ಫೈಲ್‌ಗಳನ್ನು ಸ್ಥಳಾಂತರಿಸುವ ಕಮಾಂಡ್‌ಗಳನ್ನು ಸೂಚಿಸುವ, ಆದರೆ ಎಂದಿಗೂ ರನ್ ಮಾಡದ ಒಂದು ಸ್ಕಿಲ್;
- `/ccam:usage`, `/ccam:ext` ಮತ್ತು `/ccam:switch` ಕಮಾಂಡ್‌ಗಳು;
- ಎರಡು ಹುಕ್‌ಗಳು ([ವಿವರಗಳು](/features/next#inside-claude-code)). ಸೆಷನ್ ಬಳಕೆಯ ಮಿತಿಯಲ್ಲಿ ನಿಂತಾಗ ಒಂದು ಹುಕ್ ಟರ್ಮಿನಲ್ ಅಧಿಸೂಚನೆಯನ್ನು ತೋರಿಸುತ್ತದೆ; Linux ಹಾಗೂ macOS ನಲ್ಲಿ, [`auto_switch`](/features/next#auto-switch) ಆನ್ ಆಗಿದ್ದರೆ, ccam ನಂತರ ಸೆಷನ್ ಅನ್ನು ಅವಕಾಶ ಇರುವ ಇನ್ನೊಂದು ಖಾತೆಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ. ಇನ್ನೊಂದು ಹುಕ್ `/ccam:switch [account]` ಸಂಭಾಷಣೆಯನ್ನು ಇನ್ನೊಂದು ಖಾತೆಗೆ ಒಪ್ಪಿಸಲು ಬಿಡುತ್ತದೆ.

ಹುಕ್‌ಗಳು `ccam` ಅನ್ನು ರನ್ ಮಾಡುತ್ತವೆ, ಆದ್ದರಿಂದ ಅದು Claude Code ನೋಡುವ `PATH` ನಲ್ಲಿ ಇರಬೇಕು.

## ಇನ್‌ಸ್ಟಾಲೇಶನ್ ಪರಿಶೀಲಿಸಿ

```sh
ccam version
ccam doctor
```

`ccam doctor` Claude Code ಸಿಗುತ್ತದೆಯೇ ಮತ್ತು ಪ್ರತಿ ಖಾತೆಯ ಫೋಲ್ಡರ್, ಲಿಂಕ್ ಮತ್ತು ಲಾಗಿನ್ ಸರಿಯಾಗಿದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುತ್ತದೆ. ನಂತರ [ತ್ವರಿತ ಆರಂಭ](/docs/quick-start) ಗೆ ಮುಂದುವರಿಯಿರಿ.

## ಅಪ್‌ಡೇಟ್

```sh
ccam update --check
ccam update
```

`ccam update --check` ಹೊಸ ಬಿಡುಗಡೆ ಇದೆಯೇ ಎಂಬುದನ್ನು ಮಾತ್ರ ತಿಳಿಸುತ್ತದೆ. `ccam update` ನಿಮ್ಮ ಸಿಸ್ಟಮ್‌ಗೆ ಹೊಂದುವ ಬಿಡುಗಡೆಯನ್ನು ಡೌನ್‌ಲೋಡ್ ಮಾಡುತ್ತದೆ, ಅದರ SHA-256 ಚೆಕ್‌ಸಮ್ ಅನ್ನು ಬಿಡುಗಡೆಯ `checksums.txt` ಜೊತೆ ಹೋಲಿಸಿ ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು ನಂತರ `ccam` ಫೈಲ್ ಅನ್ನು ಒಂದೇ ಹಂತದಲ್ಲಿ ಬದಲಾಯಿಸುತ್ತದೆ. ಚೆಕ್‌ಸಮ್ ಹೊಂದಿಕೆಯಾಗದಿದ್ದರೆ ಏನನ್ನೂ ಬದಲಾಯಿಸುವುದಿಲ್ಲ. `ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು ccam ಗೆ ಲಿಂಕ್‌ಗಳಾಗಿರುವುದರಿಂದ, ಅವು ತಕ್ಷಣವೇ ಹೊಸ ಆವೃತ್ತಿಯನ್ನು ರನ್ ಮಾಡುತ್ತವೆ.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ಡ್ಯಾಶ್‌ಬೋರ್ಡ್ (ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿಲ್ಲದ `ccam`) ಕೂಡ ತೆರೆದ ಪ್ರತಿ ಬಾರಿ ಹೊಸ ಬಿಡುಗಡೆ ಇದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು ಒಂದು ಇದ್ದರೆ ಅದನ್ನು ನೀಡುತ್ತದೆ: `U` ಒತ್ತಿ ಮತ್ತು ದೃಢೀಕರಿಸಿ. ಬೇರೆ ಯಾವ ಕಮಾಂಡ್ ಕೂಡ ಅಪ್‌ಡೇಟ್‌ಗಳಿಗಾಗಿ ಪರಿಶೀಲಿಸುವುದಿಲ್ಲ, ಮತ್ತು ccam ಕೇಳದೆ ಎಂದಿಗೂ ಅಪ್‌ಡೇಟ್ ಆಗುವುದಿಲ್ಲ. ಹಳೆಯ ಬಿಡುಗಡೆಯೂ ಸೇರಿದಂತೆ ಒಂದು ನಿರ್ದಿಷ್ಟ ಬಿಡುಗಡೆಯನ್ನು ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಲು `--version` ನೀಡಿ. `ccam --dry-run update` ಏನನ್ನೂ ಬದಲಾಯಿಸದೆ, ಏನು ಡೌನ್‌ಲೋಡ್ ಆಗುತ್ತದೆ ಮತ್ತು ಏನು ಬದಲಾಗುತ್ತದೆ ಎಂಬುದನ್ನು ತೋರಿಸುತ್ತದೆ.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ಇನ್‌ಸ್ಟಾಲ್ ಸ್ಕ್ರಿಪ್ಟ್ ಅನ್ನು ಮತ್ತೆ ರನ್ ಮಾಡಿದರೂ `~/.local/bin` ನಲ್ಲಿರುವ ಪ್ರತಿ ಅಪ್‌ಡೇಟ್ ಆಗುತ್ತದೆ. ನೀವು ccam ಅನ್ನು ಬಿಡುಗಡೆಯ ಟ್ಯಾಗ್‌ನಿಂದ `go install` ಮೂಲಕ ಬಿಲ್ಡ್ ಮಾಡಿದ್ದರೆ, `ccam update` ಅಥವಾ ಅದೇ `go install` ಕಮಾಂಡ್ ಅನ್ನು ಮತ್ತೆ ರನ್ ಮಾಡಬಹುದು.

**Windows.** `ccam update` Windows ಗಾಗಿ `.zip` ಅನ್ನು ಡೌನ್‌ಲೋಡ್ ಮಾಡುತ್ತದೆ. Windows ನಲ್ಲಿ ರನ್ ಆಗುತ್ತಿರುವ ಪ್ರೋಗ್ರಾಂ ಅನ್ನು ಬದಲಾಯಿಸಲು ಸಾಧ್ಯವಿಲ್ಲ, ಆದ್ದರಿಂದ ಹಳೆಯ ಬೈನರಿಯನ್ನು `ccam.exe.old` ಎಂದು ಮರುಹೆಸರಿಸಲಾಗುತ್ತದೆ; ಮುಂದಿನ ಅಪ್‌ಡೇಟ್ ಅದನ್ನು ಅಳಿಸುತ್ತದೆ. ಹಿಂದಿನ `ccam.exe.old` ನಿಂದ ಸೆಷನ್ ಇನ್ನೂ ರನ್ ಆಗುತ್ತಿದ್ದರೆ, ಹೆಸರು `ccam.exe.old-<time>` ಆಗಿರುತ್ತದೆ. ಹಾರ್ಡ್ ಲಿಂಕ್‌ಗಳಾಗಿರುವ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳನ್ನು ತಕ್ಷಣವೇ ನವೀಕರಿಸಲಾಗುತ್ತದೆ. `install.ps1` ಅನ್ನು ಮತ್ತೆ ರನ್ ಮಾಡುವುದೂ ಇದೇ ರೀತಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ, ಹಳೆಯ ಬೈನರಿಯಿಂದ Claude Code ಸೆಷನ್ ಇನ್ನೂ ರನ್ ಆಗುತ್ತಿರುವಾಗಲೂ, ಮತ್ತು ಹಾರ್ಡ್ ಲಿಂಕ್ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳನ್ನು ನವೀಕರಿಸಿದಾಗ `OK 2 shortcuts now run the new ccam` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ.

**Debian ಮತ್ತು Ubuntu ಪ್ಯಾಕೇಜ್.** apt ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದ `/usr/bin/ccam` ಅನ್ನು `ccam update` ಬದಲಾಯಿಸುವುದಿಲ್ಲ; ಬದಲಿಗೆ ಈ ಹಂತಗಳನ್ನು ಮುದ್ರಿಸುತ್ತದೆ. [ಬಿಡುಗಡೆಗಳ ಪುಟದಿಂದ](https://github.com/jhyunwoo/ccam/releases) ಹೊಸ `ccam_<version>_linux_amd64.deb` ಡೌನ್‌ಲೋಡ್ ಮಾಡಿ, ಹಳೆಯ ಪ್ಯಾಕೇಜ್‌ನ ಮೇಲೆಯೇ ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## ಅನ್‌ಇನ್‌ಸ್ಟಾಲ್

`~/.local/share/ccam` ಅನ್ನು ಅಳಿಸಿದರೆ ಅಲ್ಲಿ ಸಂಗ್ರಹವಾಗಿರುವ ಲಾಗಿನ್‌ಗಳು ಮತ್ತು ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವೂ ಅಳಿಸಿಹೋಗುತ್ತವೆ. ನೀವು [ಇಂಪೋರ್ಟ್ ಮಾಡಿದ](/docs/commands) `~/.claude` ಫೋಲ್ಡರ್ ಇದ್ದಲ್ಲೇ ಉಳಿಯುತ್ತದೆ, ಆದರೆ ಅದು ಒಂದು ಗುಂಪಿನಲ್ಲಿರುವವರೆಗೆ, ಅದು ಹಂಚಿಕೊಳ್ಳುವ ಫೈಲ್‌ಗಳು `~/.local/share/ccam/groups/` ನಲ್ಲಿ ಇರುತ್ತವೆ. ಏನನ್ನಾದರೂ ಅಳಿಸುವ ಮೊದಲು ಅದಕ್ಕಾಗಿ `ccam remove <name>` ರನ್ ಮಾಡಿ (ಅಥವಾ `ccam group leave`): ಇದು ಅದನ್ನು ಅದರ ಗುಂಪುಗಳಿಂದ ಹೊರತೆಗೆದು, ಅದು ಹಂಚಿಕೊಂಡ ಎಲ್ಲದರ ಒಂದು ಪ್ರತಿಯನ್ನು ಮತ್ತೆ ಫೋಲ್ಡರ್‌ಗೆ ಇರಿಸುತ್ತದೆ.

ccam ಇನ್ನೂ ಇನ್‌ಸ್ಟಾಲ್ ಆಗಿರುವಾಗಲೇ ಮೊದಲು ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳನ್ನು ತೆಗೆದುಹಾಕಿ. `ccam list` ಅವುಗಳನ್ನು Shortcut ಕಾಲಂನಲ್ಲಿ ತೋರಿಸುತ್ತದೆ; ಪ್ರತಿಯೊಂದನ್ನೂ `ccam shortcut <account> --remove` ಮೂಲಕ ತೆಗೆದುಹಾಕಿ, ಅಥವಾ ಆ ಲಿಂಕ್‌ಗಳನ್ನು `~/.local/bin` ನಿಂದ ನೀವೇ ಅಳಿಸಿ. ನಂತರ ಕಮಾಂಡ್ ಅನ್ನು, ಮತ್ತು ಖಾತೆಗಳು ಇನ್ನು ಬೇಡವಾದರೆ ಡೇಟಾ ಫೋಲ್ಡರ್‌ಗಳನ್ನು ತೆಗೆದುಹಾಕಿ:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows ನಲ್ಲಿ, PowerShell ನಲ್ಲಿ:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```
