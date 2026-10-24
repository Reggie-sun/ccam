---
title: ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು
description: ಈಗಿರುವ ~/.claude ನ ಸೆಟ್ಟಿಂಗ್‌ಗಳು, CLAUDE.md, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು MCP ಸರ್ವರ್‌ಗಳನ್ನು ccam migrate ಮೂಲಕ ccam ಗುಂಪಿಗೆ ನಕಲಿಸಿ.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

ccam ಗಿಂತ ಮೊದಲು ನೀವು Claude Code ಬಳಸಿದ್ದರೆ, ನಿಮ್ಮ ಸೆಟಪ್ `~/.claude` ಮತ್ತು `~/.claude.json` ನಲ್ಲಿದೆ. `ccam migrate` ಅದನ್ನು ಒಂದು ಹಂಚಿಕೆಯ ಗುಂಪಿಗೆ ನಕಲಿಸುತ್ತದೆ. ಹೀಗಾಗಿ ಆ ಗುಂಪಿನ ಮೂಲಕ ಆ ಐಟಂಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳುವ ಪ್ರತಿ ಖಾತೆಯೂ ನಿಮ್ಮ ಈಗಿನ ಸೆಟಪ್‌ನೊಂದಿಗೇ ಆರಂಭವಾಗುತ್ತದೆ. ಮೂಲ ಫೋಲ್ಡರ್ ಅನ್ನು ಓದಲಾಗುತ್ತದೆ ಅಷ್ಟೇ, ಮತ್ತು ಸಾಮಾನ್ಯ `claude` ಮೊದಲಿನಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತದೆ. ನಿಮ್ಮ ಮೊದಲ `ccam add` ನಂತರ ಇದನ್ನು ಒಮ್ಮೆ ರನ್ ಮಾಡಿ; `~/.claude.json` ನಲ್ಲಿ MCP ಸರ್ವರ್‌ಗಳಿದ್ದರೆ, ಮೊದಲ `ccam add` ತನ್ನ ಸಲಹೆಯಲ್ಲಿ ಅವುಗಳನ್ನು ಹೆಸರಿಸುತ್ತದೆ.

## ಬಳಕೆ

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ, ಐಟಂ ಆಯ್ಕೆಗಳಿಲ್ಲದೆ ಮತ್ತು `--yes` ಇಲ್ಲದೆ, ccam ಮೊದಲು ಒಂದು ಚೆಕ್‌ಲಿಸ್ಟ್ ತೋರಿಸುತ್ತದೆ:

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

Space (ಅಥವಾ `x`) ಕರ್ಸರ್ ಕೆಳಗಿನ ಐಟಂ ಅನ್ನು ಗುರುತಿಸುತ್ತದೆ ಅಥವಾ ಗುರುತು ತೆಗೆಯುತ್ತದೆ, ಬಾಣದ ಕೀಗಳು ಅಥವಾ `j`/`k` ಸರಿಸುತ್ತವೆ, Enter ಯೋಜನೆಯನ್ನು ತೋರಿಸುತ್ತದೆ, ಮತ್ತು Esc, `q` ಅಥವಾ Ctrl+C `Cancelled.` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ. ಮೂಲದಲ್ಲಿ ಏನೂ ಇಲ್ಲದ ಐಟಂಗಳನ್ನು `not found` ಎಂದು ಗುರುತಿಸಲಾಗುತ್ತದೆ. ಎಲ್ಲದರ ಗುರುತು ತೆಗೆದರೆ, ccam `Nothing selected.` ಎಂದು ಮುದ್ರಿಸಿ ನಿಲ್ಲುತ್ತದೆ.

ನಂತರ ccam ಯೋಜನೆಯನ್ನು ಮುದ್ರಿಸುತ್ತದೆ ಮತ್ತು ನಕಲಿಸುವ ಮೊದಲು ಕೇಳುತ್ತದೆ:

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

ಒಂದು ಗುಂಪು ತನ್ನದೇ ಖಾತೆಗಳನ್ನು ಮಾತ್ರ, ಅದೂ ತಾನು ಹಂಚಿಕೊಳ್ಳುವ ಐಟಂಗಳಿಗೆ ಮಾತ್ರ ತಲುಪುತ್ತದೆ, ಹೀಗಾಗಿ ಸಾರಾಂಶವು ಗುಂಪಿನ ಖಾತೆಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ. ನಕಲಿಸಿದ ಐಟಂಗಳಲ್ಲಿ ಒಂದನ್ನು ಗುಂಪು ಹಂಚಿಕೊಳ್ಳದಿದ್ದರೆ, ಅಥವಾ ಕೆಲವು ಖಾತೆಗಳು ಅದರಲ್ಲಿ ಇಲ್ಲದಿದ್ದರೆ, ccam ಅದನ್ನು ಬದಲಾಯಿಸುವ `ccam group set <group> --add <items>` ಮತ್ತು `ccam group join <group> <accounts>` ಕಮಾಂಡ್‌ಗಳನ್ನು ಮುದ್ರಿಸುತ್ತದೆ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `[folder]` | ನಕಲಿಸಬೇಕಾದ Claude Code ಫೋಲ್ಡರ್ | `~/.claude` |
| `--settings` | `settings.json` ಅನ್ನು ನಕಲಿಸುತ್ತದೆ | ಆಫ್ |
| `--memory` | ಜಾಗತಿಕ `CLAUDE.md` ಅನ್ನು ನಕಲಿಸುತ್ತದೆ | ಆಫ್ |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` ಅನ್ನು ನಕಲಿಸುತ್ತದೆ | ಆಫ್ |
| `--mcp` | ಬಳಕೆದಾರ MCP ಸರ್ವರ್‌ಗಳನ್ನು (`.claude.json` ನಲ್ಲಿನ `mcpServers`) ನಕಲಿಸುತ್ತದೆ | ಆಫ್ |
| `--history` | ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು (`projects/`, `history.jsonl`) ನಕಲಿಸುತ್ತದೆ | ಆಫ್ |
| `-g`, `--group string` | ನಕಲಿಸಬೇಕಾದ ಹಂಚಿಕೆಯ ಗುಂಪು | `default` |
| `-y`, `--yes` | ಕೇಳದೆ ನಕಲಿಸುತ್ತದೆ; ಸಂಘರ್ಷವಿದ್ದಾಗ ಗುಂಪಿನ ನಕಲನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ | ಆಫ್ |
| `--dry-run` (ಜಾಗತಿಕ) | ಯೋಜನೆಯನ್ನು ಮುದ್ರಿಸಿ, ಯಾವುದೇ ಪ್ರಶ್ನೆ ಅಥವಾ ನಕಲಿಗೆ ಮೊದಲು ನಿಲ್ಲುತ್ತದೆ | ಆಫ್ |

ಯಾವುದೇ ಐಟಂ ಆಯ್ಕೆ ಇಲ್ಲದಿದ್ದರೆ, ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಚೆಕ್‌ಲಿಸ್ಟ್ ನಿರ್ಧರಿಸುತ್ತದೆ. ಟರ್ಮಿನಲ್‌ನ ಹೊರಗೆ, ಅಥವಾ `--yes` ಇದ್ದಾಗ, ccam ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವನ್ನೂ ನಕಲಿಸುತ್ತದೆ. `projects/` ತುಂಬಾ ದೊಡ್ಡದಾಗಿರಬಹುದಾದ್ದರಿಂದ ಇತಿಹಾಸವನ್ನು ಡೀಫಾಲ್ಟ್ ಆಗಿ ಬಿಡಲಾಗುತ್ತದೆ.

`--group on` ಎಂದರೆ `default`. `--group off` ಮತ್ತು `--group local` ಅನ್ನು `--group takes a group name, not "off"` ಎಂದು ನಿರಾಕರಿಸಲಾಗುತ್ತದೆ. `ccam migrate` `--dry-run` ಅನ್ನು ಪಾಲಿಸುತ್ತದೆ, ಮತ್ತು ಡ್ರೈ ರನ್‌ನಲ್ಲೂ ಚೆಕ್‌ಲಿಸ್ಟ್ ಕಾಣಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ನೀವು ಆಯ್ಕೆಯನ್ನು ಮುಂಚಿತವಾಗಿ ನೋಡಬಹುದು.

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

`internal/cli/migrate.go` ಐಟಂಗಳನ್ನು ಆಯ್ಕೆಮಾಡುತ್ತದೆ, ನಂತರ `migrate.Build` ಮೂಲ ಮತ್ತು ಗುಂಪನ್ನು ಓದಿ ಒಂದು ಯೋಜನೆಯನ್ನು ಹಿಂತಿರುಗಿಸುತ್ತದೆ. ನೀವು ದೃಢೀಕರಿಸಿದ ನಂತರ `Plan.Apply` ಆಗುವವರೆಗೆ ಏನನ್ನೂ ಬರೆಯಲಾಗುವುದಿಲ್ಲ.

1. **ಮೂಲವನ್ನು ಹುಡುಕುವುದು.** ಇಲ್ಲದ, ಫೋಲ್ಡರ್ ಅಲ್ಲದ, ಅಥವಾ ಯಾವುದೇ ತಿಳಿದಿರುವ Claude Code ನಮೂದು ಇಲ್ಲದ ಪಾತ್ ಅನ್ನು `checkSource` ನಿರಾಕರಿಸುತ್ತದೆ. MCP ಸರ್ವರ್‌ಗಳನ್ನು ಮಾತ್ರ ನಕಲಿಸುವಾಗ, `~/.claude.json` ನಲ್ಲಿ ಸರ್ವರ್‌ಗಳಿರುವವರೆಗೆ ಖಾಲಿ ಅಥವಾ ಇಲ್ಲದ `~/.claude` ಪರವಾಗಿಲ್ಲ. `migrate.ClaudeJSON` ಫೋಲ್ಡರ್‌ಗೆ ಸೇರಿದ `.claude.json` ಅನ್ನು ಹುಡುಕುತ್ತದೆ: ಅದರೊಳಗೆ, ಅಥವಾ ಒಳಗೆ ಅದಿಲ್ಲದ `~/.claude` ಗಾಗಿ, ಅದರ ಪಕ್ಕದ `~/.claude.json`. `CLAUDE_CONFIG_DIR` ಹೊಂದಿಸದಿದ್ದಾಗ Claude Code ಅದನ್ನು ಇಡುವುದು ಅಲ್ಲೇ.
2. **ಐಟಂಗಳನ್ನು ಪಾತ್‌ಗಳಿಗೆ ಹೊಂದಿಸುವುದು.** ಪ್ರತಿ ಪಾತ್ `share.Layout` ನಿಂದ ಬರುತ್ತದೆ, ಇದು ಗುಂಪುಗಳು ಬಳಸುವ ಅದೇ ಅನುಮತಿ ಪಟ್ಟಿ. MCP ಸರ್ವರ್‌ಗಳನ್ನು `.claude.json` ನಿಂದ `mcp.LocalServers` ಮೂಲಕ ಓದಲಾಗುತ್ತದೆ, ಇದು ಇಡೀ ಫೈಲ್ ಅನ್ನು ಪಾರ್ಸ್ ಮಾಡುವ ಬದಲು `profile.TopLevelValue` ಸ್ಕ್ಯಾನರ್ ಬಳಸುತ್ತದೆ.
3. **ನಕಲಿಸುವುದು, ಎಂದಿಗೂ ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ.** ಫೈಲ್‌ಗಳು `fsutil.AtomicWrite` ಅಥವಾ `fsutil.CopyFile` ಮೂಲಕ `~/.local/share/ccam/groups/<group>/` ಗೆ ಹೋಗುತ್ತವೆ; ಅದು ಗುರಿಯ ಪಕ್ಕದಲ್ಲಿ ತಾತ್ಕಾಲಿಕ ಫೈಲ್ ಬರೆದು ಅದನ್ನು ಮರುಹೆಸರಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಅರ್ಧದಲ್ಲಿ ನಿಂತ ನಕಲು ಎಂದಿಗೂ ಕತ್ತರಿಸಿದ ಫೈಲ್ ಅನ್ನು ಬಿಡುವುದಿಲ್ಲ. MCP ಸರ್ವರ್‌ಗಳು ಗುಂಪಿನ `mcp.json` ಗೆ ಅದರ ಲಾಕ್ ಅಡಿಯಲ್ಲಿ ಹೋಗುತ್ತವೆ: ಈ ನಡುವೆ ಇನ್ನೊಂದು ccam ಸೇರಿಸುವ ಸರ್ವರ್ ಉಳಿಯುತ್ತದೆ, ಮತ್ತು [Remote Control](/guides/remote-control) ಗಾಗಿ ಮಾಡಿದ ನಕಲುಗಳನ್ನು ನಂತರ ನವೀಕರಿಸಲಾಗುತ್ತದೆ.

ಪ್ರತಿ ಐಟಂ ಒಂದು ಮೂಲ ಪಾತ್ ಅನ್ನು ಒಂದು ಗುಂಪಿನ ಪಾತ್‌ಗೆ ಹೊಂದಿಸುತ್ತದೆ:

<!-- motion:migrate-items -->

ಇತಿಹಾಸದ ಎರಡು ಸಾಲುಗಳು ಡ್ಯಾಶ್ ರೇಖೆಯಲ್ಲಿವೆ, ಏಕೆಂದರೆ ಅವುಗಳನ್ನು `--history` ಇದ್ದಾಗ ಮಾತ್ರ ನಕಲಿಸಲಾಗುತ್ತದೆ. `.credentials.json` ಮತ್ತು `.claude.json` `share.Never` ನಲ್ಲಿ ಪಟ್ಟಿಯಾಗಿವೆ ಮತ್ತು ಎಂದಿಗೂ ಗುಂಪನ್ನು ತಲುಪುವುದಿಲ್ಲ: ಅವುಗಳಲ್ಲಿ ಒಂದು ಲೇಔಟ್‌ನಲ್ಲಿ ಕಾಣಿಸಿಕೊಂಡರೆ `Build` ಆಂತರಿಕ ದೋಷವನ್ನು ಹಿಂತಿರುಗಿಸುತ್ತದೆ, ಮತ್ತು ಫೋಲ್ಡರ್ ಪರಿಶೀಲನೆ ಅವುಗಳನ್ನು ಹೆಸರಿನಿಂದಲೇ ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.

`buildEntry` ಪ್ರತಿ ಪಾತ್‌ಗೆ ಒಂದು ಸ್ಥಿತಿಯನ್ನು ನೀಡುತ್ತದೆ; ಯೋಜನೆಯ ಪ್ರತಿ ಸಾಲು ವಿವರಿಸುವುದು ಅದನ್ನೇ:

<!-- motion:migrate-states -->

- **Missing**: ಮೂಲದಲ್ಲಿ ಇಲ್ಲ, ಅಥವಾ ಫೋಲ್ಡರ್ ಇರಬೇಕಾದಲ್ಲಿ ಫೈಲ್ ಇದೆ. ಪ್ರತಿ ಐಟಂಗೆ ಒಮ್ಮೆ `not found, skipped` ಎಂದು ಪಟ್ಟಿಯಾಗುತ್ತದೆ.
- **Linked**: ಮೂಲ ಪಾತ್ ಈಗಾಗಲೇ ಗುಂಪಿನ ನಕಲಿಗೆ ತಲುಪುತ್ತದೆ, ಉದಾಹರಣೆಗೆ ಮೂಲವು ccam ಪ್ರೊಫೈಲ್ ಫೋಲ್ಡರ್ ಆಗಿದ್ದಾಗ.
- **Same**: ಗುಂಪಿನಲ್ಲಿ ಈಗಾಗಲೇ ಒಂದೇ ರೀತಿಯ ವಿಷಯ ಇದೆ.
- **New**: ಗುಂಪಿನಲ್ಲಿ ಅದಿಲ್ಲ, ಹೀಗಾಗಿ ಅದನ್ನು ಪೂರ್ತಿಯಾಗಿ ನಕಲಿಸಲಾಗುತ್ತದೆ.
- **Merge**: ಎರಡೂ ಕಡೆ ಇರುವ ಫೋಲ್ಡರ್. `copyMissing` ಗುಂಪಿನಲ್ಲಿ ಇಲ್ಲದ ಫೈಲ್‌ಗಳನ್ನು ಮಾತ್ರ ಸೇರಿಸುತ್ತದೆ; ಬೇರೆಯಾಗಿರುವ ಫೈಲ್‌ಗಳು ಗುಂಪಿನಲ್ಲಿರುವಂತೆಯೇ ಉಳಿಯುತ್ತವೆ ಮತ್ತು ಯೋಜನೆ `group copies kept` ಎಂದು ಹೇಳುತ್ತದೆ. `plugins/` ಗಾಗಿ, `installed_plugins.json` ಮತ್ತು `known_marketplaces.json` ರಿಜಿಸ್ಟ್ರಿಗಳನ್ನು ವಿಲೀನಗೊಳಿಸಲಾಗುತ್ತದೆ, ಮತ್ತು ಪ್ಲಗಿನ್‌ಗಳ ಸಂಪೂರ್ಣ ಪಾತ್‌ಗಳನ್ನು ಗುಂಪಿನ ಫೋಲ್ಡರ್‌ಗೆ ತೋರುವಂತೆ ಮತ್ತೆ ಬರೆಯಲಾಗುತ್ತದೆ.
- **Conflict**: ಎರಡೂ ಕಡೆ ಇರುವ, ಬೇರೆ ಬೇರೆ ವಿಷಯದ ಫೈಲ್. ನೀವು ಫೈಲ್‌ವಾರು ಆಯ್ಕೆಮಾಡುತ್ತೀರಿ: `keep group`, `use mine`, ಜೊತೆಗೆ `.json` ಫೈಲ್‌ಗಳಿಗೆ `merge JSON` (ಮೌಲ್ಯಗಳು ಸಂಘರ್ಷಿಸಿದಾಗ ಗುಂಪು ಗೆಲ್ಲುವ ಆಳವಾದ ವಿಲೀನ) ಅಥವಾ `history.jsonl` ಗೆ `merge by line`. ವ್ಯಾಖ್ಯಾನ ಬೇರೆಯಾಗಿರುವ MCP ಸರ್ವರ್ ಕೂಡ ಸಂಘರ್ಷವೇ, `keep group` ಅಥವಾ `use mine` ಆಯ್ಕೆಗಳೊಂದಿಗೆ. ನೀವು Enter ಒತ್ತಿದಾಗ ಪಟ್ಟಿಯ ಮೊದಲ ಆಯ್ಕೆ ಡೀಫಾಲ್ಟ್; `--yes` ಯಾವಾಗಲೂ ಗುಂಪಿನ ನಕಲನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ.

ಚೆಕ್‌ಲಿಸ್ಟ್ ಸ್ವತಃ `tui.MigratePicker`, ಇದು `ccam add` ಬಳಸುವ `checklist` ಕಾಂಪೊನೆಂಟ್ ಮೇಲೆ ನಿರ್ಮಿತವಾಗಿದೆ; `migrate.Found` lstat ಕರೆಗಳು ಮತ್ತು MCP ಸರ್ವರ್ ಪಟ್ಟಿಯ ಮೂಲಕ ಇಲ್ಲದ ಐಟಂಗಳನ್ನು ಗುರುತಿಸುತ್ತದೆ.

### MCP ಸೈನ್-ಇನ್‌ಗಳು

OAuth ಮೂಲಕ ಸೈನ್ ಇನ್ ಆಗುವ MCP ಸರ್ವರ್‌ಗಳು ತಮ್ಮ ಟೋಕನ್‌ಗಳನ್ನು ಸರ್ವರ್ ಪಟ್ಟಿಯಲ್ಲಿ ಅಲ್ಲ, ಲಾಗಿನ್ ಫೈಲ್ `.credentials.json` ನಲ್ಲಿ `mcpOAuth` ಅಡಿಯಲ್ಲಿ ಇಡುತ್ತವೆ. ccam ಅವುಗಳನ್ನು ನಕಲಿಸುವುದಿಲ್ಲ. ನಕಲಿಸಿದ ಸೈನ್-ಇನ್ ಮೂಲದಲ್ಲಿರುವ ಅದೇ ರಿಫ್ರೆಶ್ ಟೋಕನ್ ಅನ್ನು ಹೊಂದಿರುತ್ತದೆ, ಮತ್ತು ಮೊದಲ ನವೀಕರಣದಲ್ಲೇ ಸರ್ವರ್ ಆ ಟೋಕನ್ ಅನ್ನು ಬದಲಾಯಿಸುತ್ತದೆ: ಆಗ ಅದನ್ನು ಹೊಂದಿರುವ ಉಳಿದೆಲ್ಲರೂ, ಮೂಲ `~/.claude` ಸೇರಿದಂತೆ, ಆ ಸರ್ವರ್‌ನಿಂದ ಲಾಗ್ ಔಟ್ ಆಗುತ್ತಾರೆ. ಆದ್ದರಿಂದ ಯೋಜನೆಯ `sign-ins` ಸಾಲು, ಮೂಲ ಸೈನ್ ಇನ್ ಆಗಿರುವ ಪ್ರತಿ ನಕಲಿಸಿದ ಸರ್ವರ್‌ಗೆ, `/mcp` ಮೂಲಕ ತಾವೇ ಸೈನ್ ಇನ್ ಆಗುವ ಗುಂಪಿನ ಖಾತೆಗಳನ್ನು ಹೆಸರಿಸುತ್ತದೆ, ಮತ್ತು ಸಾರಾಂಶ ಅದನ್ನು ಪುನರಾವರ್ತಿಸುತ್ತದೆ. ಟೋಕನ್ ಇಲ್ಲದ ನಮೂದು ಸೈನ್-ಇನ್ ಅಲ್ಲ ಮತ್ತು ಅದನ್ನು ಪಟ್ಟಿ ಮಾಡಲಾಗುವುದಿಲ್ಲ; ಈಗಾಗಲೇ ತನ್ನದೇ ಸೈನ್-ಇನ್ ಹೊಂದಿರುವ ಖಾತೆಯನ್ನೂ, ಗುಂಪು ತನ್ನ ವ್ಯಾಖ್ಯಾನವನ್ನು ಉಳಿಸಿಕೊಳ್ಳುವ ಸರ್ವರ್ ಅನ್ನೂ ಪಟ್ಟಿ ಮಾಡಲಾಗುವುದಿಲ್ಲ. ಯಾವುದೇ ಲಾಗಿನ್ ಫೈಲ್ ಅನ್ನು ಬರೆಯಲಾಗುವುದಿಲ್ಲ.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- ಮೂಲವನ್ನು ಓದಲಾಗುತ್ತದೆ ಅಷ್ಟೇ. ಅದರಲ್ಲಿನ ಯಾವ ಫೈಲ್ ಅನ್ನೂ ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ, ಲಿಂಕ್ ಮಾಡುವುದಿಲ್ಲ ಅಥವಾ ಬದಲಾಯಿಸುವುದಿಲ್ಲ, ಮತ್ತು ಮೂಲದ ಲಾಗಿನ್‌ಗೆ ಎಂದಿಗೂ ಬರೆಯುವುದಿಲ್ಲ; Keychain ನಿಂದ ಓದಬಹುದಾದ macOS ನಲ್ಲೂ ಇದು ಅನ್ವಯಿಸುತ್ತದೆ.
- `use mine`, ವಿಲೀನ, ಅಥವಾ ಈಗಿರುವ `mcp.json` ನಲ್ಲಿ ಹೊಸ MCP ಸರ್ವರ್‌ಗಳಿಂದ ಗುಂಪಿನ ಫೈಲ್ ಬದಲಾಗುವ ಮೊದಲು, ಗುಂಪಿನ ನಕಲನ್ನು `<group>/.ccam-backup/<time>/` ನಲ್ಲಿ ಉಳಿಸಲಾಗುತ್ತದೆ. ಫೋಲ್ಡರ್ ಹೆಸರು ಪ್ರತಿ ರನ್‌ಗೆ ಅನನ್ಯ, ಹೀಗಾಗಿ ಒಂದೇ ಸೆಕೆಂಡಿನ ಎರಡು ರನ್‌ಗಳು ಪರಸ್ಪರರ ಬ್ಯಾಕಪ್‌ಗಳನ್ನು ತಿದ್ದಿ ಬರೆಯುವುದಿಲ್ಲ. ಸಾರಾಂಶವು ಪಾತ್ ಅನ್ನು ಮುದ್ರಿಸುತ್ತದೆ.
- `--yes` ಇದ್ದಾಗ, ಪ್ರತಿ ಸಂಘರ್ಷದಲ್ಲೂ ಗುಂಪಿನ ನಕಲು ಉಳಿಯುತ್ತದೆ, ಮತ್ತು ಬದಲಿಗೆ ಹೇಗೆ ಆಯ್ಕೆಮಾಡಬೇಕು ಎಂದು ccam ಹೇಳುತ್ತದೆ: ಅದೇ ಕಮಾಂಡ್ ಅನ್ನು `--yes` ಇಲ್ಲದೆ ರನ್ ಮಾಡಿ.
- ಏನೂ ಬದಲಾಗದಿದ್ದರೆ, ccam `Nothing to copy.` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ; ಪ್ರತಿ ಸಂಘರ್ಷದಲ್ಲೂ ಗುಂಪಿನ ನಕಲು ಉಳಿದರೆ, `Nothing to copy: the group keeps its own copies.`
- ನಕಲು ಅರ್ಧದಲ್ಲಿ ವಿಫಲವಾದರೆ, ccam `The source folder was not changed. Group files replaced so far are saved in` ಮತ್ತು ಅದರ ನಂತರ ಬ್ಯಾಕಪ್ ಫೋಲ್ಡರ್‌ನೊಂದಿಗೆ ನಿಲ್ಲುತ್ತದೆ. ಮತ್ತೆ ರನ್ ಮಾಡಿದರೆ ಪೂರ್ತಿಯಾಗಿ ತಲುಪದ ಫೈಲ್‌ಗಳು ನಕಲಾಗುತ್ತವೆ.
- ಹಲವು ಗುಂಪುಗಳು MCP ಹಂಚಿಕೊಂಡಾಗ, `ccam doctor`, `ccam ext` ಮತ್ತು `ccam mcp list` ಸೂಚಿಸುವ ಕಮಾಂಡ್‌ಗಳು `ccam migrate --mcp --group <a|b>` ಎಂದಿರುತ್ತವೆ.
- ಮೂಲ ಇಲ್ಲದಿದ್ದರೆ `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` ಎಂದು ಮುಗಿಯುತ್ತದೆ.

## ಸಂಬಂಧಿತ

- [MCP ಸರ್ವರ್‌ಗಳು, ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮತ್ತು ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳಿ](/docs/sharing)
- [ಪರಿಶೀಲನೆ ಮತ್ತು ದುರಸ್ತಿ](/features/doctor)
- [ಹಲವು Claude Code ಖಾತೆಗಳು](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
