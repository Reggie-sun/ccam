---
title: ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ
description: Linux ಅಥವಾ macOS ನಲ್ಲಿ ಒಂದೇ ಕಮಾಂಡ್‌ನಿಂದ, .deb ಪ್ಯಾಕೇಜ್‌ನಿಂದ ಅಥವಾ ಸೋರ್ಸ್‌ನಿಂದ ccam ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ. ಮೊದಲು Claude Code ಇನ್‌ಸ್ಟಾಲ್ ಆಗಿರಬೇಕು.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Linux ಅಥವಾ macOS ನ ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಇದನ್ನು ರನ್ ಮಾಡಿ:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ಈ ಸ್ಕ್ರಿಪ್ಟ್ ನಿಮ್ಮ ಸಿಸ್ಟಮ್‌ಗೆ ಹೊಂದುವ ಬಿಡುಗಡೆಯನ್ನು ಡೌನ್‌ಲೋಡ್ ಮಾಡುತ್ತದೆ, ಅದರ SHA-256 ಚೆಕ್‌ಸಮ್ ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು `ccam` ಅನ್ನು `~/.local/bin` ನಲ್ಲಿ ಇರಿಸುತ್ತದೆ; ನಂತರ `ccam add` ಅದೇ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ `claude1`, `claude2`… ಶಾರ್ಟ್‌ಕಟ್‌ಗಳನ್ನು ರಚಿಸುತ್ತದೆ. ಆ ಫೋಲ್ಡರ್ ನಿಮ್ಮ `PATH` ನಲ್ಲಿ ಇಲ್ಲದಿದ್ದರೆ, `~/.bashrc` ಅಥವಾ `~/.zshrc` ಗೆ ಸೇರಿಸಬೇಕಾದ ಸಾಲನ್ನು ಸ್ಕ್ರಿಪ್ಟ್ ಮುದ್ರಿಸುತ್ತದೆ.

ccam ನಿಮಗಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಆದ್ದರಿಂದ ಮೊದಲು [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿ ಮತ್ತು `claude --version` ಕೆಲಸ ಮಾಡುತ್ತದೆಯೇ ಎಂದು ಖಚಿತಪಡಿಸಿಕೊಳ್ಳಿ.

## ಬೆಂಬಲಿತ ಸಿಸ್ಟಮ್‌ಗಳು

- x86-64 ಮತ್ತು ARM64 ನಲ್ಲಿ Linux: Ubuntu, Debian, Fedora, Arch, Alpine ಮತ್ತು ಇತರ ಡಿಸ್ಟ್ರಿಬ್ಯೂಷನ್‌ಗಳು. ಬೈನರಿ ಸ್ಟ್ಯಾಟಿಕ್ ಆಗಿದೆ ಮತ್ತು ಅದಕ್ಕೆ ಯಾವುದೇ ಡಿಪೆಂಡೆನ್ಸಿ ಇಲ್ಲ.
- Apple silicon ಮತ್ತು Intel ನಲ್ಲಿ macOS.
- Windows ಗೆ ನೇರ ಬೆಂಬಲ ಇಲ್ಲ. WSL ಒಳಗೆ ccam ಯಾವುದೇ Linux ಸಿಸ್ಟಮ್‌ನಲ್ಲಿ ಕೆಲಸ ಮಾಡುವಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತದೆ.

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

**ನಿರ್ದಿಷ್ಟ ಆವೃತ್ತಿ.** `CCAM_VERSION` ಹೊಂದಿಸಿ:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## ಶೆಲ್ completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion ಖಾತೆಗಳ ಹೆಸರುಗಳನ್ನೂ ಒಳಗೊಂಡಂತೆ `ccam` ನ ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು ಮತ್ತು ಆಯ್ಕೆಗಳಿಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ.

## ಇನ್‌ಸ್ಟಾಲೇಶನ್ ಪರಿಶೀಲಿಸಿ

```sh
ccam version
ccam doctor
```

`ccam doctor` Claude Code ಸಿಗುತ್ತದೆಯೇ ಮತ್ತು ಪ್ರತಿ ಖಾತೆಯ ಫೋಲ್ಡರ್, ಲಿಂಕ್ ಮತ್ತು ಲಾಗಿನ್ ಸರಿಯಾಗಿದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುತ್ತದೆ. ನಂತರ [ತ್ವರಿತ ಆರಂಭ](/docs/quick-start) ಗೆ ಮುಂದುವರಿಯಿರಿ.

## ಅನ್‌ಇನ್‌ಸ್ಟಾಲ್

ಕಮಾಂಡ್ ಮತ್ತು ಖಾತೆಗಳ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳನ್ನು ತೆಗೆದುಹಾಕಿ. ನಿಮಗೆ ಖಾತೆಗಳು ಇನ್ನು ಬೇಡವಾದರೆ, ನಂತರ ಡೇಟಾ ಫೋಲ್ಡರ್‌ಗಳನ್ನೂ ತೆಗೆದುಹಾಕಿ:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` ಅನ್ನು ಅಳಿಸಿದರೆ ಅಲ್ಲಿ ಸಂಗ್ರಹವಾಗಿರುವ ಲಾಗಿನ್‌ಗಳು ಮತ್ತು ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವೂ ಅಳಿಸಿಹೋಗುತ್ತವೆ. ನೀವು [ಇಂಪೋರ್ಟ್ ಮಾಡಿದ](/docs/commands) `~/.claude` ಫೋಲ್ಡರ್ ಇದ್ದಲ್ಲೇ ಉಳಿಯುತ್ತದೆ.
