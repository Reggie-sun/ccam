---
title: ತ್ವರಿತ ಆರಂಭ
metaTitle: "ತ್ವರಿತ ಆರಂಭ: Claude Code ಖಾತೆ ಸೇರಿಸಿ, ಬದಲಾಯಿಸಿ"
description: ccam add ಮೂಲಕ ಎರಡು Claude Code ಖಾತೆಗಳನ್ನು ಸೇರಿಸಿ, ccwork ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್‌ ನೀಡಿ, ಅವುಗಳ ನಡುವೆ ಬದಲಾಯಿಸಿ ಮತ್ತು ಏನನ್ನು ಹಂಚಬೇಕು ಎಂದು ಆಯ್ಕೆಮಾಡಿ.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

ಪ್ರತಿ ಖಾತೆಯನ್ನು `ccam add` ಮೂಲಕ ಸೇರಿಸಿ, ನಂತರ `ccam <name>` ಮೂಲಕ ಅಥವಾ `ccwork` ನಂತಹ ನೀವೇ ಹೆಸರಿಡುವ ಶಾರ್ಟ್‌ಕಟ್ ಮೂಲಕ Claude Code ಆರಂಭಿಸಿ. ಹೆಚ್ಚಿನವರಿಗೆ ಬೇಕಾಗಿರುವುದು ಇಷ್ಟೇ.

## 1. ನಿಮ್ಮ ಮೊದಲ ಖಾತೆಯನ್ನು ಸೇರಿಸಿ

```sh
ccam add
```

ವಿಜಾರ್ಡ್ `personal` ಅಥವಾ `work` ನಂತಹ ಒಂದು ಹೆಸರನ್ನು, ಈ ಖಾತೆ ನಿಮ್ಮ ಇತರ ಖಾತೆಗಳೊಂದಿಗೆ ಯಾವ ಐಟಂಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳಬೇಕು ಎಂಬುದನ್ನು, ಮತ್ತು ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್ ರಚಿಸಬೇಕೇ ಎಂಬುದನ್ನು (ಹಂತ 4 ನೋಡಿ) ಕೇಳುತ್ತದೆ. ನಂತರ ಅದು `claude auth login` ರನ್ ಮಾಡುತ್ತದೆ. ಇದು ಬ್ರೌಸರ್ ತೆರೆಯುತ್ತದೆ, ಅಲ್ಲಿ ನೀವು ಸೈನ್ ಇನ್ ಮಾಡಬಹುದು. ರಿಮೋಟ್ ಯಂತ್ರದಲ್ಲಿ, ಬ್ರೌಸರ್ ಕೊನೆಗೆ ಲೋಡ್ ಆಗದ `http://localhost:…/callback` ಪುಟವನ್ನು ತಲುಪಬಹುದು; ಆ ವಿಳಾಸವನ್ನು ಟರ್ಮಿನಲ್‌ಗೆ ಅಂಟಿಸಿ, ccam ಸೈನ್-ಇನ್ ಅನ್ನು ಪೂರ್ಣಗೊಳಿಸುತ್ತದೆ. ಮುಗಿದಾಗ ನಿಮಗೆ ಇದು ಕಾಣುತ್ತದೆ:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. ನಿಮ್ಮ ಈಗಿರುವ ಸೆಟ್ಟಿಂಗ್‌ಗಳನ್ನು ತನ್ನಿ

ನೀವು ಹಿಂದೆ Claude Code ಬಳಸಿದ್ದರೆ, ನಿಮ್ಮ ಸೆಟ್ಟಿಂಗ್‌ಗಳು, `CLAUDE.md`, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು MCP ಸರ್ವರ್‌ಗಳು `~/.claude` ನಲ್ಲಿವೆ. ಪ್ರತಿ ಖಾತೆಗೂ ಅವು ಸಿಗುವಂತೆ ಅವುಗಳನ್ನು ಹಂಚಿಕೆಯ ಸೆಟ್ಟಿಂಗ್‌ಗಳಿಗೆ ನಕಲಿಸಿ:

```sh
ccam migrate
```

ccam ತಾನು ಏನನ್ನು ನಕಲಿಸುತ್ತದೆ ಎಂದು ತೋರಿಸಿ ಮೊದಲು ಕೇಳುತ್ತದೆ. ಇದು `~/.claude` ಅನ್ನು ಓದುತ್ತದೆ ಅಷ್ಟೇ, ಯಾವುದೇ ಲಾಗಿನ್ ಅನ್ನು ನಕಲಿಸುವುದಿಲ್ಲ, ಮತ್ತು ನೀವು `--history` ಸೇರಿಸದ ಹೊರತು ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ. ವಿವರಗಳಿಗೆ [ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು](/docs/sharing#copy-an-existing-setup) ನೋಡಿ.

## 3. ಎರಡನೇ ಖಾತೆಯನ್ನು ಸೇರಿಸಿ

```sh
ccam add work
```

ಬ್ರೌಸರ್‌ನಲ್ಲಿ ಇನ್ನೊಂದು ಖಾತೆಯಿಂದ ಸೈನ್ ಇನ್ ಮಾಡಿ.

## 4. ಖಾತೆಗಳನ್ನು ಬದಲಾಯಿಸಿ

`ccam work` Claude Code ಅನ್ನು work ಖಾತೆಯಾಗಿ ಆರಂಭಿಸುತ್ತದೆ. ಇನ್ನೂ ಚಿಕ್ಕದು ಬೇಕಿದ್ದರೆ, ಖಾತೆಗೆ `ccwork` ನಂತಹ ಅದರದೇ ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್ ನೀಡಿ. ಹಂಚಿಕೊಳ್ಳುವ ಐಟಂಗಳ ನಂತರ ವಿಜಾರ್ಡ್ ಹೀಗೆ ಕೇಳುತ್ತದೆ:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` ಒತ್ತಿ, ನಂತರ ಒಂದು ಹೆಸರನ್ನು ಟೈಪ್ ಮಾಡಿ, ಅಥವಾ ಸೂಚಿಸಿದ ಹೆಸರನ್ನು ಬಳಸಲು Enter ಒತ್ತಿ:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

ಲಾಗಿನ್ ನಂತರ ccam ಅದನ್ನು ದೃಢಪಡಿಸುತ್ತದೆ:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

ನೀವು ಆ ಪ್ರಶ್ನೆಯನ್ನು ಬಿಟ್ಟಿದ್ದರೆ, ಯಾವಾಗ ಬೇಕಾದರೂ ಶಾರ್ಟ್‌ಕಟ್ ಸೇರಿಸಬಹುದು:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

ಈಗ ಇವುಗಳಲ್ಲಿ ಯಾವುದಾದರೂ ಒಂದು Claude Code ಅನ್ನು work ಖಾತೆಯಾಗಿ ಆರಂಭಿಸುತ್ತದೆ:

```sh
ccwork
ccam work
ccam 2
```

`ccam` ಅನ್ನು ಮಾತ್ರ ರನ್ ಮಾಡಿದರೆ [ಡ್ಯಾಶ್‌ಬೋರ್ಡ್](/features/dashboard) ತೆರೆಯುತ್ತದೆ: ಬಾಣದ ಕೀಲಿಗಳು ಅಥವಾ ಖಾತೆಯ ಸಂಖ್ಯೆಯಿಂದ ಖಾತೆಗೆ ಹೋಗಿ, ನಂತರ Enter ಒತ್ತಿ.

ಶಾರ್ಟ್‌ಕಟ್ ಅಥವಾ ಖಾತೆಯ ಹೆಸರಿನ ನಂತರದ ಎಲ್ಲವೂ ಯಾವುದೇ ಬದಲಾವಣೆ ಇಲ್ಲದೆ Claude Code ಗೆ ಹೋಗುತ್ತದೆ. ಹೀಗಾಗಿ ನಿಮಗೆ ಗೊತ್ತಿರುವ ಆಯ್ಕೆಗಳನ್ನೇ ಬಳಸುತ್ತಿರಬಹುದು:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. ಏನನ್ನು ಹಂಚಿಕೊಳ್ಳಬೇಕು ಎಂದು ಆಯ್ಕೆಮಾಡಿ

ಡೀಫಾಲ್ಟ್ ಆಗಿ MCP ಸರ್ವರ್‌ಗಳು, `settings.json`, ಜಾಗತಿಕ `CLAUDE.md`, ಹಾಗೂ agents, commands, skills ಮತ್ತು plugins `default` ಗುಂಪಿನ ಮೂಲಕ ಖಾತೆಗಳ ನಡುವೆ ಹಂಚಿಕೆಯಾಗುತ್ತವೆ. ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸ ಹಂಚಿಕೆಯಾಗುವುದಿಲ್ಲ. ಖಾತೆಗಳ ನಡುವೆ ಸಂಭಾಷಣೆಗಳನ್ನು ಮುಂದುವರಿಸಬೇಕಿದ್ದರೆ, ಅದನ್ನು ಗುಂಪಿಗೆ ಸೇರಿಸಿ:

```sh
ccam group set default --add history
```

ಯಾವುದನ್ನಾದರೂ ಬದಲಾಯಿಸುವ ಮೊದಲು, ccam ಪ್ರತಿ ಖಾತೆಗೂ ತಾನು ಏನನ್ನು ಸ್ಥಳಾಂತರಿಸಲಿದೆ ಎಂದು ತೋರಿಸುತ್ತದೆ. ವಿವರಗಳಿಗೆ [ಗುಂಪುಗಳೊಂದಿಗೆ ಹಂಚಿಕೆ](/features/share) ನೋಡಿ.

## 6. ಖಾತೆ ತನ್ನ ಮಿತಿ ತಲುಪಿದಾಗ

```sh
ccam usage
ccam next --continue
```

`ccam usage` ಎಲ್ಲಾ ಖಾತೆಗಳ 5 ಗಂಟೆಯ ಮತ್ತು ವಾರದ ಬಳಕೆಯನ್ನು ತೋರಿಸುತ್ತದೆ. `ccam next --continue` ಅದೇ ಸಂಭಾಷಣೆಯನ್ನು ಮುಂದಿನ ಖಾತೆಯಲ್ಲಿ ತೆರೆಯುತ್ತದೆ. ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳಲ್ಲಿ, `ccam --fallback work -p "…"` ತಾನಾಗಿಯೇ ಬದಲಾಯಿಸುತ್ತದೆ. [ಬಳಕೆಯ ಮಿತಿಗಳು](/docs/usage-limits) ನೋಡಿ.

**Team ಮತ್ತು Enterprise ಸೀಟುಗಳು.** ಸ್ವಯಂಚಾಲಿತ ಆಯ್ಕೆ (`ccam next`, `--fallback`, `ccam pick`, ಡ್ಯಾಶ್‌ಬೋರ್ಡ್‌ನ `n`, ಮತ್ತು ನೀವು ಆನ್ ಮಾಡಿದಾಗ [`auto_switch`](/features/next#auto-switch)) ಒಂದೇ Claude Team ಅಥವಾ Enterprise ಸಂಸ್ಥೆಯ ಸೀಟುಗಳ ನಡುವೆಯೂ ಬದಲಾಯಿಸುತ್ತದೆ. Anthropic ನ ಸಹಾಯ ಕೇಂದ್ರವು "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team ಪ್ಲಾನ್‌ನ ಪ್ರತಿ ಸೀಟನ್ನು ಹೆಸರಿಸಲಾದ ಒಬ್ಬ ವ್ಯಕ್ತಿಗೆ ನಿಯೋಜಿಸಲಾಗುತ್ತದೆ, ಮತ್ತು ಆ ವ್ಯಕ್ತಿ ಮಾತ್ರ ಅದನ್ನು ಬಳಸಬಹುದು) ಮತ್ತು "Users can't route their work through another user's account to use that user's remaining usage" (ಬಳಕೆದಾರರು ಇನ್ನೊಬ್ಬ ಬಳಕೆದಾರರ ಉಳಿದಿರುವ ಬಳಕೆಯನ್ನು ಉಪಯೋಗಿಸಿಕೊಳ್ಳಲು ತಮ್ಮ ಕೆಲಸವನ್ನು ಅವರ ಖಾತೆಯ ಮೂಲಕ ನಡೆಸುವಂತಿಲ್ಲ) ಎಂದು ಹೇಳುತ್ತದೆ ([Team ಪ್ಲಾನ್‌ಗಳು](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise ಪ್ಲಾನ್‌ಗಳಿಗೂ](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ಇದೇ ಅನ್ವಯಿಸುತ್ತದೆ). ಆದ್ದರಿಂದ ಸೀಟುಗಳ ನಡುವೆ ಸ್ವಯಂಚಾಲಿತ ಆಯ್ಕೆಯನ್ನು ಬಳಸುವುದು Anthropic ನ ನಿಯಮಗಳನ್ನು ಉಲ್ಲಂಘಿಸಬಹುದು. ಅಲ್ಲಿ ಇದನ್ನು ಬಳಸಬೇಕೆ ಎಂಬುದು ನಿಮ್ಮ ನಿರ್ಧಾರ: ಇದನ್ನು ಅವಲಂಬಿಸುವ ಮೊದಲು ನಿಮ್ಮ ಸಂಸ್ಥೆ ಮತ್ತು Anthropic ಬಳಿ ವಿಚಾರಿಸಿ ಖಚಿತಪಡಿಸಿಕೊಳ್ಳಿ. ವೈಯಕ್ತಿಕ ಪ್ಲಾನ್‌ಗಳ (Pro, Max) ಬಗ್ಗೆ [ಪದೇ ಪದೇ ಕೇಳುವ ಪ್ರಶ್ನೆಗಳಲ್ಲಿನ ನಿಯಮಗಳ ಟಿಪ್ಪಣಿ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ತಿಳಿಸುತ್ತದೆ.

## ಈಗಾಗಲೇ ~/.claude ಬಳಸುತ್ತಿದ್ದೀರಾ?

ಮತ್ತೆ ಸೈನ್ ಇನ್ ಮಾಡುವ ಬದಲು ಅದನ್ನು ನೋಂದಾಯಿಸಿ:

```sh
ccam import
```

ಫೋಲ್ಡರ್ ಇದ್ದಲ್ಲೇ ಉಳಿಯುತ್ತದೆ, ಮತ್ತು ಸಾಮಾನ್ಯ `claude` ಮೊದಲಿನಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತದೆ.
