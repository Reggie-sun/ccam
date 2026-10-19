---
title: ಕಮಾಂಡ್ ಉಲ್ಲೇಖ
description: ccam add ಮತ್ತು ccam <name> ನಿಂದ ccam usage, ccam next ಮತ್ತು ccam doctor ವರೆಗೆ, ccam ನ ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು ಮತ್ತು ಆಯ್ಕೆಗಳು.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

ಯಾವುದೇ ಕಮಾಂಡ್ ಏನು ಮಾಡುತ್ತದೆ ಎಂಬುದನ್ನು ನಿಜವಾಗಿ ಮಾಡದೆಯೇ ನೋಡಲು, ಅದರ ಮುಂದೆ `--dry-run` ಹಾಕಿ. ಉದಾಹರಣೆಗೆ `ccam --dry-run work`.

<!-- commands -->

## Claude Code ಆರಂಭಿಸುವುದು

`ccam <name|number> [arguments]` ಆ ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಹೆಸರಿನ ನಂತರದ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು ccam ಓದುವುದಿಲ್ಲ; ಅದು ಅವುಗಳನ್ನು ಅದೇ ಕ್ರಮದಲ್ಲಿ `claude` ಗೆ ರವಾನಿಸುತ್ತದೆ. ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್, ಎಕ್ಸಿಟ್ ಕೋಡ್ ಮತ್ತು ಸಿಗ್ನಲ್‌ಗಳು `claude` ನಲ್ಲಿ ಇರುವಂತೆಯೇ ವರ್ತಿಸುತ್ತವೆ. ಹೆಸರಿನ ನಂತರ `--` ಐಚ್ಛಿಕ.

ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿದ್ದು ಖಾತೆಯ ಹೆಸರು ಇಲ್ಲದ `ccam`, ಉದಾಹರಣೆಗೆ `ccam -p "summarise"`, ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತದೆ. ಅದನ್ನು `ccam use <name>` ಮೂಲಕ ಹೊಂದಿಸಿ.

`claude1`, `claude2` ಮುಂತಾದ ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು `~/.local/bin` ನಲ್ಲಿರುವ ccam ಗೆ ಲಿಂಕ್‌ಗಳು. `ccam shims install --names` ಪ್ರತಿ ಖಾತೆಗೂ `claude-<name>` ಅನ್ನು ಸಹ ರಚಿಸುತ್ತದೆ.

## ಖಾತೆಗಳನ್ನು ಸೇರಿಸುವುದು ಮತ್ತು ತೆಗೆದುಹಾಕುವುದು

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ಫೋಲ್ಡರ್ ನೀಡದೆ `ccam import` ರನ್ ಮಾಡಿದರೆ, ಅದು `~/.claude` ಮತ್ತು `~/.claude-*` ಗಾಗಿ ಹುಡುಕುತ್ತದೆ ಮತ್ತು ಪ್ರತಿಯೊಂದರ ಬಗ್ಗೆ ಕೇಳುತ್ತದೆ. ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಫೋಲ್ಡರ್‌ಗಳು ಇದ್ದಲ್ಲೇ ನೋಂದಣಿಯಾಗುತ್ತವೆ; ccam ಎಂದಿಗೂ ಲಾಗಿನ್ ಅನ್ನು ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ ಅಥವಾ ನಕಲಿಸುವುದಿಲ್ಲ.

## ಹಂಚಿಕೆ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ಪ್ರತಿ ಐಟಂ `on` (`default` ಗುಂಪು), `off`, ಅಥವಾ ಒಂದು ಗುಂಪಿನ ಹೆಸರನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ. [ಹಂಚಿಕೆ](/docs/sharing) ನೋಡಿ.

## MCP ಸರ್ವರ್‌ಗಳು

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## ಬಳಕೆ ಮತ್ತು ಬದಲಾವಣೆ

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## ನಿರ್ವಹಣೆ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ಫೈಲ್‌ಗಳು

| ಪಾತ್ | ಒಳಗೊಂಡಿರುವುದು |
| --- | --- |
| `~/.config/ccam/config.toml` | ಖಾತೆಗಳು ಮತ್ತು ಅವು ಏನನ್ನು ಹಂಚಿಕೊಳ್ಳುತ್ತವೆ |
| `~/.local/share/ccam/profiles/<name>/` | ಪ್ರತಿ ಖಾತೆಯ `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | ಹಂಚಿಕೊಂಡ ಸೆಟ್ಟಿಂಗ್‌ಗಳು, MCP ಸರ್ವರ್‌ಗಳು, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು ಇತಿಹಾಸ |
| `~/.local/state/ccam/` | ಬಳಕೆಯ ಕ್ಯಾಶ್ ಮತ್ತು ಕೊನೆಯ ಬಾರಿ ಬಳಸಿದ ಸಮಯಗಳು |
| `~/.local/bin/claude1`, `claude2`… | ccam ಗೆ ಲಿಂಕ್‌ಗಳು |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ಮತ್ತು `XDG_STATE_HOME` ವೇರಿಯೇಬಲ್‌ಗಳು ಈ ಫೋಲ್ಡರ್‌ಗಳ ಸ್ಥಳವನ್ನು ಬದಲಾಯಿಸುತ್ತವೆ.
