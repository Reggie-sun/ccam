---
title: ccam ಜೊತೆಗೆ Remote Control ಬಳಸುವುದು
description: ccam <name> rc ಮೂಲಕ ಯಾವುದೇ ಖಾತೆಯಾಗಿ Claude Code Remote Control ಆರಂಭಿಸಿ; ಅನ್ವಯಿಸುವ ಆಯ್ಕೆಗಳು, ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳ ನಕಲನ್ನು ccam ಏಕೆ ಸೂಚಿಸುತ್ತದೆ, ಸಮಸ್ಯೆ ಪರಿಹಾರ.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc` Claude Code ನ Remote Control ಅನ್ನು `work` ಖಾತೆಯಾಗಿ ಆರಂಭಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಈ ಕಂಪ್ಯೂಟರ್‌ನಲ್ಲಿನ ಸೆಷನ್‌ಗಳನ್ನು claude.ai/code ಅಥವಾ Claude ಮೊಬೈಲ್ ಆ್ಯಪ್‌ನಿಂದ ನಡೆಸಬಹುದು. ಇದು `claude rc` ನಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತದೆ, ಒಂದೇ ವ್ಯತ್ಯಾಸದೊಂದಿಗೆ: ಖಾತೆ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಹಂಚಿಕೊಂಡಿದ್ದರೆ, ಅವುಗಳನ್ನು ಖಾತೆಯ `.claude.json` ಗೆ ನಕಲಿಸಬೇಕೇ ಎಂದು ccam ಒಮ್ಮೆ ಕೇಳುತ್ತದೆ, ಏಕೆಂದರೆ Remote Control ಸೆಷನ್‌ಗಳು ಬೇರೆ ಯಾವ ಸ್ಥಳದಿಂದಲೂ ಸರ್ವರ್‌ಗಳನ್ನು ಓದುವುದಿಲ್ಲ.

## ಆರಂಭಿಸಿ

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[ಶಾರ್ಟ್‌ಕಟ್](/features/shortcuts) ಕೂಡ ಇದೇ ರೀತಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ: `ccwork rc` ಎಂದರೆ `ccam work rc`. ನೀವು ಕೆಲಸ ಮಾಡಬೇಕಾದ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಕಮಾಂಡ್ ರನ್ ಮಾಡಿ, ನಂತರ `work` ಬಳಸುವ ಅದೇ Claude ಖಾತೆಯಿಂದ ಸೈನ್ ಇನ್ ಆಗಿರುವ claude.ai/code ಅಥವಾ Claude ಆ್ಯಪ್ ಅನ್ನು ತೆರೆಯಿರಿ.

## ಆಯ್ಕೆಗಳು

ಇಲ್ಲಿ ಎರಡು ಬಗೆಯ ಆಯ್ಕೆಗಳು ಸೇರುತ್ತವೆ: ccam ನ ಸ್ವಂತ ಆಯ್ಕೆಗಳು, ಇವು ಖಾತೆಯ ಹೆಸರಿನ ಮೊದಲು ಬರುತ್ತವೆ ಅಥವಾ ccam ನ ಸೆಟ್ಟಿಂಗ್‌ಗಳಲ್ಲಿ ಇರುತ್ತವೆ; ಮತ್ತು Remote Control ನ ಆಯ್ಕೆಗಳು, ಇವು `rc` ನಂತರ ಬರುತ್ತವೆ ಮತ್ತು ಬದಲಾವಣೆ ಇಲ್ಲದೆ Claude Code ಗೆ ತಲುಪುತ್ತವೆ.

### ccam ಆಯ್ಕೆಗಳು ಮತ್ತು ಸೆಟ್ಟಿಂಗ್‌ಗಳು

| ಆಯ್ಕೆ ಅಥವಾ ಸೆಟ್ಟಿಂಗ್ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `--dry-run` (ಖಾತೆಯ ಮೊದಲು) | ಪರಿಸರ ಮತ್ತು ಕಮಾಂಡ್ ಅನ್ನು, ಮತ್ತು MCP ಸೂಚನೆ ಇದ್ದರೆ ಅದನ್ನೂ ಮುದ್ರಿಸುತ್ತದೆ. ಕೇಳುವುದೂ ಇಲ್ಲ, ನಕಲಿಸುವುದೂ ಇಲ್ಲ | ಆಫ್ |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | ಖಾತೆ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಹಂಚಿಕೊಂಡಿದ್ದು, ಅದರ `.claude.json` ನಲ್ಲಿ ಗುಂಪಿನ ಪ್ರಸ್ತುತ ಪಟ್ಟಿ ಇಲ್ಲದಿದ್ದಾಗ ಕೇಳಲಾಗುತ್ತದೆ. `y` ಅವುಗಳನ್ನು ನಕಲಿಸುತ್ತದೆ, `N` ಅಥವಾ Enter ಅವುಗಳಿಲ್ಲದೆ ಆರಂಭಿಸುತ್ತದೆ | `N` |
| `~/.config/ccam/config.toml` ನಲ್ಲಿ `mcp_strategy = "sync"` | Remote Control ಆಗಿರಲಿ ಇಲ್ಲದಿರಲಿ, ಪ್ರತಿ ಆರಂಭದ ಮೊದಲು ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳನ್ನು `.claude.json` ಗೆ ನಕಲಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಏನನ್ನೂ ಎಂದಿಗೂ ಕೇಳಲಾಗುವುದಿಲ್ಲ | `"flag"` |
| `ccam share work --mcp off` | ಖಾತೆಯೊಂದಿಗೆ MCP ಸರ್ವರ್‌ಗಳ ಹಂಚಿಕೆಯನ್ನು ನಿಲ್ಲಿಸುತ್ತದೆ; ಆಗ Remote Control ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್‌ಗಳೊಂದಿಗೆ ಮಾತ್ರ ಆರಂಭವಾಗುತ್ತದೆ ಮತ್ತು ಏನನ್ನೂ ಕೇಳಲಾಗುವುದಿಲ್ಲ | ಆನ್ |

### Remote Control ಆಯ್ಕೆಗಳು

`rc` ನಂತರದ ಎಲ್ಲವೂ Claude Code ಗೆ ಸೇರಿದ್ದು. ಇವು Claude Code 2.1.304 ನಲ್ಲಿನ `claude remote-control` ನ ಆಯ್ಕೆಗಳು; ನಿಮ್ಮ ಆವೃತ್ತಿಯ ಪಟ್ಟಿಯನ್ನು `ccam work rc --help` ಮುದ್ರಿಸುತ್ತದೆ.

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ |
| --- | --- |
| `--name <name>` | claude.ai/code ನಲ್ಲಿ ತೋರುವ ಸೆಷನ್‌ನ ಹೆಸರು |
| `--remote-control-session-name-prefix <prefix>` | ಸ್ವಯಂ ರಚಿತ ಸೆಷನ್ ಹೆಸರುಗಳ ಪೂರ್ವಪ್ರತ್ಯಯ (ಡೀಫಾಲ್ಟ್: ಹೋಸ್ಟ್ ಹೆಸರು) |
| `-c`, `--continue` | ಈ ಫೋಲ್ಡರ್‌ಗಾಗಿ Remote Control ಕೊನೆಯದಾಗಿ ದಾಖಲಿಸಿದ ಸೆಷನ್‌ಗೆ, ಸುಮಾರು 4 ಗಂಟೆಗಳೊಳಗೆ, ಮತ್ತೆ ಸಂಪರ್ಕಿಸುತ್ತದೆ |
| `--session-id <id>` | ನಿರ್ದಿಷ್ಟ ಸೆಷನ್‌ಗೆ ಮತ್ತೆ ಸಂಪರ್ಕಿಸುತ್ತದೆ |
| `--permission-mode <mode>` | ಹೊಸ ಸೆಷನ್‌ಗಳ ಅನುಮತಿ ಮೋಡ್: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` ಅಥವಾ `plan` |
| `--spawn <mode>` | `same-dir` (ಡೀಫಾಲ್ಟ್), `worktree` (ಪ್ರತಿ ಸೆಷನ್‌ಗೆ ಒಂದು git worktree) ಅಥವಾ `session` (ಒಂದೇ ಸೆಷನ್, ಅದು ಮುಗಿದಾಗ ನಿರ್ಗಮಿಸುತ್ತದೆ) |
| `--capacity <N>` | `same-dir` ಮತ್ತು `worktree` ಮೋಡ್‌ನಲ್ಲಿ ಒಂದೇ ಬಾರಿಗೆ ಗರಿಷ್ಠ ಸೆಷನ್‌ಗಳು (ಡೀಫಾಲ್ಟ್ 32) |
| `--[no-]create-session-in-dir` | ಆರಂಭದಲ್ಲೇ ಪ್ರಸ್ತುತ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಮೊದಲ ಸೆಷನ್ ರಚಿಸುತ್ತದೆ (ಡೀಫಾಲ್ಟ್ ಆನ್) |
| `--[no-]chrome` | ಹೊಸ ಸೆಷನ್‌ಗಳಿಗೆ Claude in Chrome (ಡೀಫಾಲ್ಟ್ ಆಫ್) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | ಡೀಬಗ್ ಮತ್ತು ವಿವರವಾದ ಔಟ್‌ಪುಟ್ |

`rc` ಅನ್ನು ಖಾತೆಯ ಹೆಸರಿನ ನಂತರ ತಕ್ಷಣ ಇಡಿ. `-` ನಿಂದ ಆರಂಭವಾಗದ ಮೊದಲ ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು ccam Claude Code ನ ಕಮಾಂಡ್ ಎಂದು ಪರಿಗಣಿಸುತ್ತದೆ, ಹೀಗಾಗಿ `ccam work rc --name laptop` ನಲ್ಲಿ ಅದು `rc` ಅನ್ನು ನೋಡುತ್ತದೆ, ಆದರೆ `ccam work --model opus rc` ನಲ್ಲಿ `opus` ಅನ್ನು ನೋಡುತ್ತದೆ.

## ccam ಏನು ಮಾಡುತ್ತದೆ

ಆರಂಭವು ಸಾಮಾನ್ಯ ಮಾರ್ಗದಲ್ಲೇ ನಡೆಯುತ್ತದೆ ([Claude Code ಆರಂಭಿಸುವುದು](/features/launch)), MCP ಹಂತದಲ್ಲಿ ಮಾತ್ರ ಒಂದು ಬದಲಾವಣೆಯೊಂದಿಗೆ:

1. ಯಾವುದೇ ಆರಂಭದಲ್ಲಿ ಮಾಡುವಂತೆ, ccam ಖಾತೆಯನ್ನು ಹುಡುಕುತ್ತದೆ, ಅದರ ಹಂಚಿಕೆ ಲಿಂಕ್‌ಗಳನ್ನು ಪರಿಶೀಲಿಸುತ್ತದೆ ಮತ್ತು `claude` ಅನ್ನು ಪತ್ತೆಹಚ್ಚುತ್ತದೆ.
2. ಕಮಾಂಡ್ `rc` ಅಥವಾ `remote-control` ಎಂದು ಗುರುತಿಸಿ `--mcp-config` ಅನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.
3. ಖಾತೆ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಹಂಚಿಕೊಂಡಿದ್ದರೆ, ಗುಂಪಿನ `mcp.json` ಅನ್ನು `~/.local/state/ccam/mcp/work/` ನಲ್ಲಿರುವ `sync-stamp` ನೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ. `.claude.json` ನಲ್ಲಿ ಈಗಾಗಲೇ ಪ್ರಸ್ತುತ ಸರ್ವರ್‌ಗಳಿದ್ದರೆ, ಅಥವಾ ಗುಂಪಿನಲ್ಲಿ ಯಾವುದೇ ಸರ್ವರ್ ಇಲ್ಲದಿದ್ದರೆ, ನೇರವಾಗಿ ಹಂತ 6 ಕ್ಕೆ ಹೋಗುತ್ತದೆ.
4. ಇಲ್ಲದಿದ್ದರೆ ಸೆಷನ್‌ಗಳಿಗೆ ಸಿಗದ ಸರ್ವರ್‌ಗಳ ಹೆಸರುಗಳನ್ನು ತಿಳಿಸುತ್ತದೆ:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y` ನೀಡಿದರೆ, ಅದು `.claude.json` ನ ಬ್ಯಾಕಪ್ ಉಳಿಸುತ್ತದೆ, Claude Code ನ ಸ್ವಂತ ಲಾಕ್ ಪಡೆಯುತ್ತದೆ, ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳನ್ನು `mcpServers` ಗೆ ಸೇರಿಸುತ್ತದೆ ಮತ್ತು ಖಾತೆಯ ಇತರ ಸರ್ವರ್‌ಗಳನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ. `N` ನೀಡಿದರೆ, ಟರ್ಮಿನಲ್ ಇಲ್ಲದಿದ್ದರೆ, ಅಥವಾ ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿದ್ದರೆ, ಅದು ಏನನ್ನೂ ಬದಲಾಯಿಸದೆ ಅದನ್ನು ತಿಳಿಸುತ್ತದೆ:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. ಅದು ಸೆಷನ್ ಅನ್ನು ದಾಖಲಿಸುತ್ತದೆ, `CLAUDE_CONFIG_DIR` ಅನ್ನು ಖಾತೆಯ ಫೋಲ್ಡರ್‌ಗೆ ಹೊಂದಿಸುತ್ತದೆ ಮತ್ತು `exec` ಮೂಲಕ ತನ್ನನ್ನು `claude rc` ನಿಂದ ಬದಲಾಯಿಸಿಕೊಳ್ಳುತ್ತದೆ.

ನಂತರ Remote Control ಸರ್ವರ್ ಆಗಿ ಚಾಲನೆಯಲ್ಲಿರುತ್ತದೆ. ಅದು ನಿಮಗಾಗಿ ಆರಂಭಿಸುವ ಪ್ರತಿ ಸೆಷನ್ ಆ ಪ್ರೋಸೆಸ್‌ನ ಚೈಲ್ಡ್ ಆಗಿರುತ್ತದೆ, ಹೀಗಾಗಿ ಅದು `work` ಆಗಿ ರನ್ ಆಗುತ್ತದೆ: `work` ನ ಲಾಗಿನ್, ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮತ್ತು ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸದೊಂದಿಗೆ, ಮತ್ತು `work` ನ ಬಳಕೆಯ ಮಿತಿಗೆ ಲೆಕ್ಕವಾಗುತ್ತದೆ. ಸರ್ವರ್ ಚಾಲನೆಯಲ್ಲಿರುವವರೆಗೆ ccam ಖಾತೆಯನ್ನು ಚಾಲನೆಯಲ್ಲಿದೆ ಎಂದು ಪರಿಗಣಿಸುತ್ತದೆ.

ಡ್ರೈ ರನ್ ಈ ವ್ಯತ್ಯಾಸವನ್ನು ತೋರಿಸುತ್ತದೆ. ಸಾಮಾನ್ಯ ಆರಂಭ ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳನ್ನು ಆರ್ಗ್ಯುಮೆಂಟ್ ಆಗಿ ಒಯ್ಯುತ್ತದೆ, Remote Control ಒಯ್ಯುವುದಿಲ್ಲ:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## ಇದು ಏಕೆ ಹೀಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

**Remote Control ಯಾವುದೇ `--mcp-config` ತೆಗೆದುಕೊಳ್ಳುವುದಿಲ್ಲ.** ಸಾಮಾನ್ಯ ಆರಂಭಗಳಲ್ಲಿ ccam ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳನ್ನು Claude Code ಗೆ `--mcp-config=<path>` ಆಗಿ ನೀಡುತ್ತದೆ ಮತ್ತು Claude Code ನ ಫೈಲ್‌ಗಳಿಗೆ ಎಂದಿಗೂ ಬರೆಯುವುದಿಲ್ಲ. `claude remote-control` ತನ್ನ ಕಮಾಂಡ್ ಹೆಸರಿನ ಮೊದಲು ಆ ಆಯ್ಕೆಯನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ ಮತ್ತು ನಂತರ ಯಾವುದನ್ನೂ ಸ್ವೀಕರಿಸುವುದಿಲ್ಲ, ಹೀಗಾಗಿ ccam 0.6.0 ಕ್ಕಿಂತ ಮೊದಲು, MCP ಸರ್ವರ್‌ಗಳನ್ನು ಹಂಚಿಕೊಂಡ ಪ್ರತಿ ಖಾತೆಗೂ `ccam work rc` ವಿಫಲವಾಗುತ್ತಿತ್ತು. ಈಗ ccam ಆ ಆಯ್ಕೆಯನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.

**ಅದರ ಸೆಷನ್‌ಗಳು ಸರ್ವರ್‌ಗಳನ್ನು `.claude.json` ನಿಂದ ಮಾತ್ರ ಓದುತ್ತವೆ.** Remote Control ಆರಂಭಿಸುವ ಸೆಷನ್‌ಗಳು ಕಮಾಂಡ್ ಲೈನ್‌ನ ಆಯ್ಕೆಗಳನ್ನು ಪಡೆದುಕೊಳ್ಳುವುದಿಲ್ಲ. ಬಳಕೆದಾರ MCP ಸರ್ವರ್‌ಗಳು ಅವುಗಳಿಗೆ ಸಿಗುವ ಏಕೈಕ ಸ್ಥಳ ಖಾತೆಯ `.claude.json`. ಆ ಫೈಲ್‌ನಲ್ಲಿ ಖಾತೆಯ ಲಾಗಿನ್ ವಿವರಗಳು ಮತ್ತು ಫೋಲ್ಡರ್‌ವಾರು ಸ್ಥಿತಿಯೂ ಇರುತ್ತದೆ, ಹೀಗಾಗಿ `settings.json` ಗೆ ಮಾಡುವಂತೆ ccam ಅದನ್ನು ಹಂಚಿಕೊಂಡ ಪ್ರತಿಗೆ ಲಿಂಕ್ ಮಾಡಲು ಸಾಧ್ಯವಿಲ್ಲ. ಸರ್ವರ್‌ಗಳನ್ನು ಅದರೊಳಗೆ ನಕಲಿಸುವುದೇ ಸೆಷನ್‌ಗಳಿಗೆ ಹಂಚಿಕೊಂಡ ಪಟ್ಟಿಯನ್ನು ನೀಡುವ ಏಕೈಕ ಮಾರ್ಗ.

**ಬರೆಯುವ ಮೊದಲು ccam ಕೇಳುತ್ತದೆ.** ಡೀಫಾಲ್ಟ್ ಫ್ಲ್ಯಾಗ್ ವಿಧಾನದಲ್ಲಿ, `.claude.json` ಅನ್ನು ಮುಟ್ಟುವುದಿಲ್ಲ ಎಂದು ccam ಭರವಸೆ ನೀಡುತ್ತದೆ. ಅಲ್ಲಿ ಬರೆಯಲೇಬೇಕಾದ ಏಕೈಕ ಸಂದರ್ಭ Remote Control, ಹೀಗಾಗಿ ಅದು ಮೊದಲು ಕೇಳುತ್ತದೆ, ಬ್ಯಾಕಪ್ ಇಡುತ್ತದೆ, ಮತ್ತು `mcpServers` ಹೊರತು ಬೇರೇನನ್ನೂ ಮುಟ್ಟುವುದಿಲ್ಲ.

**ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿರುವಾಗ ಅದು ಬರೆಯುವುದಿಲ್ಲ.** ಸೆಷನ್ ನಡೆಯುವಾಗ Claude Code `.claude.json` ಅನ್ನು ಮತ್ತೆ ಬರೆಯುತ್ತದೆ, ಮತ್ತು ಹೊರಗಿನಿಂದ ಮಾಡಿದ ಬದಲಾವಣೆ ಕಳೆದುಹೋಗುತ್ತದೆ ಅಥವಾ ಅದರ ಮೇಲೆ ಮತ್ತೆ ಬರೆಯಲಾಗುತ್ತದೆ. ಅದಕ್ಕಾಗಿಯೇ ಇನ್ನೊಂದು Remote Control ಸರ್ವರ್ ಸೇರಿದಂತೆ, ಖಾತೆಯ ಯಾವುದೇ ಸೆಷನ್ ಚಾಲನೆಯಲ್ಲಿದ್ದಾಗ ccam ನಕಲನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.

**ಪ್ರತಿ ಬದಲಾವಣೆಗೆ ಅದು ಒಮ್ಮೆ ಮಾತ್ರ ಕೇಳುತ್ತದೆ.** ನಕಲಿನ ನಂತರ `sync-stamp` ಗುಂಪಿನ `mcp.json` ನ ಮಾರ್ಪಾಡಿನ ಸಮಯವನ್ನು ಹೊಂದಿರುತ್ತದೆ. ಆ ಫೈಲ್ ಬದಲಾದಾಗ ಮಾತ್ರ, ಉದಾಹರಣೆಗೆ `ccam mcp add` ನಂತರ, ccam ಮತ್ತೆ ಕೇಳುತ್ತದೆ. ನೀವು `N` ಎಂದು ಉತ್ತರಿಸಿದರೆ ಇನ್ನೂ ನಕಲು ಇರುವುದಿಲ್ಲ, ಹೀಗಾಗಿ ಪ್ರತಿ Remote Control ಆರಂಭದಲ್ಲೂ ಅದು ಕೇಳುತ್ತದೆ.

**`ccam mcp` ನಕಲುಗಳನ್ನು ನವೀಕೃತವಾಗಿ ಇಡುತ್ತದೆ.** ಹೆಸರು ಒಂದೇ ಆದಾಗ `.claude.json` ನಲ್ಲಿನ ನಮೂದು `--mcp-config` ಮೇಲೆ ಗೆಲ್ಲುತ್ತದೆ, ಹೀಗಾಗಿ ಹಳೆಯ ನಕಲು ಸಾಮಾನ್ಯ ಆರಂಭಗಳಲ್ಲೂ ಹೊಸ ವ್ಯಾಖ್ಯಾನವನ್ನು ಮರೆಮಾಡುತ್ತದೆ. ಆದ್ದರಿಂದ `ccam mcp add`, `add-json`, `remove` ಮತ್ತು `import` ಚಾಲನೆಯಲ್ಲಿಲ್ಲದ ಪ್ರತಿ ಖಾತೆಯ ನಕಲನ್ನು ನವೀಕರಿಸುತ್ತವೆ:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**ಟರ್ಮಿನಲ್ ಇಲ್ಲದಿದ್ದರೆ ಉತ್ತರ "ಇಲ್ಲ".** Remote Control ಅನೇಕ ಬಾರಿ ಸ್ಕ್ರಿಪ್ಟ್, ಸರ್ವಿಸ್ ಅಥವಾ `nohup` ನಿಂದ ಚಾಲನೆಯಾಗುತ್ತದೆ, ಅಲ್ಲಿ ಉತ್ತರಿಸಲು ಯಾರೂ ಇರುವುದಿಲ್ಲ. ಆಗ ccam ಕಾಯುವ ಬದಲು ಸರ್ವರ್‌ಗಳಿಲ್ಲದೆ ಆರಂಭಿಸುತ್ತದೆ.

## ಸಮಸ್ಯೆಗಳು ಮತ್ತು ಪರಿಹಾರಗಳು

### Remote Control ಸೆಷನ್‌ಗಳಲ್ಲಿ ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು ಕಾಣೆಯಾಗಿವೆ

ಪ್ರಶ್ನೆಗೆ `N` ಎಂದು ಉತ್ತರಿಸಲಾಗಿತ್ತು, ಟರ್ಮಿನಲ್ ಇರಲಿಲ್ಲ, ಅಥವಾ ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿತ್ತು. Remote Control ಅನ್ನು ನಿಲ್ಲಿಸಿ, ಖಾತೆಯ ಇತರ ಸೆಷನ್‌ಗಳನ್ನು ಮುಚ್ಚಿ, ಮತ್ತು ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ `ccam work rc` ಅನ್ನು ಮತ್ತೆ ರನ್ ಮಾಡಿ; `y` ಎಂದು ಉತ್ತರಿಸಿ. ಎಂದಿಗೂ ಕೇಳಬಾರದೆಂದರೆ, `~/.config/ccam/config.toml` ನಲ್ಲಿ `mcp_strategy = "sync"` ಹೊಂದಿಸಿ. ಸೆಷನ್‌ನಲ್ಲಿ ಈಗಾಗಲೇ ಇರುವ ಸರ್ವರ್‌ಗಳು ಅದು ನಡೆಯುತ್ತಿರುವಾಗ ಬದಲಾಗುವುದಿಲ್ಲ, ಹೀಗಾಗಿ ನಕಲಿನ ನಂತರ ಹೊಸ ಸೆಷನ್ ಆರಂಭಿಸಿ.

### "Account work is running, so its .claude.json is left alone"

`work` ನ ಒಂದು ಸೆಷನ್ ಇನ್ನೂ ತೆರೆದಿದೆ: ಸಾಮಾನ್ಯ `ccam work`, ಶಾರ್ಟ್‌ಕಟ್, ಅಥವಾ ಇನ್ನೊಂದು Remote Control ಸರ್ವರ್. ಅದರಿಂದ ನಿರ್ಗಮಿಸಿ `ccam work rc` ಅನ್ನು ಮತ್ತೆ ರನ್ ಮಾಡಿ. ಬೇರೆ ಖಾತೆ ನಡೆಸುತ್ತಿರುವ ಸೆಷನ್‌ಗಳು ಲೆಕ್ಕಕ್ಕೆ ಬರುವುದಿಲ್ಲ.

### ccam ಪ್ರತಿ ಬಾರಿಯೂ ಕೇಳುತ್ತದೆ

ನೀವು ಇನ್ನೂ ಸರ್ವರ್‌ಗಳನ್ನು ನಕಲಿಸಿಲ್ಲ. ಒಮ್ಮೆ `y` ಎಂದು ಉತ್ತರಿಸಿ, ಅಥವಾ `mcp_strategy = "sync"` ಹೊಂದಿಸಿ. Remote Control ನಲ್ಲಿ ಸರ್ವರ್‌ಗಳು ಎಂದಿಗೂ ಬೇಡವಾದರೆ, ಖಾತೆಯೊಂದಿಗೆ MCP ಹಂಚಿಕೆಯನ್ನು ನಿಲ್ಲಿಸಿ: `ccam share work --mcp off`.

### `ccam work rc` ತಕ್ಷಣವೇ `--mcp-config` ಕುರಿತ ದೋಷದೊಂದಿಗೆ ವಿಫಲವಾಗುತ್ತದೆ

ನಿಮ್ಮ ccam 0.6.0 ಗಿಂತ ಹಳೆಯದು ಮತ್ತು Remote Control ಗೆ `--mcp-config` ಅನ್ನು ರವಾನಿಸುತ್ತದೆ. `ccam update` ರನ್ ಮಾಡಿ. ಪ್ರಸ್ತುತ ccam ನಲ್ಲೂ ಇದೇ ದೋಷ ಬಂದರೆ, ಮೌಲ್ಯವಿರುವ ಯಾವುದೋ ಆಯ್ಕೆ `rc` ಗಿಂತ ಮೊದಲು ಬಂದಿದೆ ಎಂದರ್ಥ, ಉದಾಹರಣೆಗೆ `ccam work --model opus rc` ನಲ್ಲಿ, ಹೀಗಾಗಿ ccam `opus` ಅನ್ನು ಕಮಾಂಡ್ ಎಂದು ತೆಗೆದುಕೊಂಡಿದೆ; `rc` ಅನ್ನು ಖಾತೆಯ ಹೆಸರಿನ ನಂತರ ತಕ್ಷಣ ಇಡಿ.

### `ccam mcp` ಮೂಲಕ ಮಾಡಿದ ಬದಲಾವಣೆ Remote Control ಗೆ ತಲುಪುವುದಿಲ್ಲ

`ccam mcp add` ಮತ್ತು `remove` ಚಾಲನೆಯಲ್ಲಿಲ್ಲದ ಖಾತೆಗಳಲ್ಲಿ ಮಾತ್ರ ನಕಲುಗಳನ್ನು ನವೀಕರಿಸುತ್ತವೆ, ಮತ್ತು ಅದನ್ನು ತಿಳಿಸುತ್ತವೆ: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control ಅನ್ನು ನಿಲ್ಲಿಸಿ ಮತ್ತೆ ಆರಂಭಿಸಿ; ccam ಬದಲಾದ `mcp.json` ಅನ್ನು ಗಮನಿಸಿ ಮತ್ತೆ ಕೇಳುತ್ತದೆ.

### ನಕಲು ನನ್ನದೇ ಆದ ಒಂದು ಸರ್ವರ್ ಅನ್ನು ಬದಲಾಯಿಸಿತು

ಗುಂಪಿನ ಸರ್ವರ್‌ನ ಹೆಸರಿನದೇ ಸ್ವಂತ ಸರ್ವರ್ ಖಾತೆಯಲ್ಲಿದ್ದರೆ, ನಕಲು ಅದನ್ನು ಬದಲಾಯಿಸಿದೆ. ಹಿಂದಿನ `.claude.json` `~/.local/state/ccam/mcp/work/backups/` ನಲ್ಲಿದೆ, ಅಲ್ಲಿ ಇತ್ತೀಚಿನ ಐದು ಬ್ಯಾಕಪ್‌ಗಳು ಉಳಿಯುತ್ತವೆ. ನಿಮ್ಮ ಆವೃತ್ತಿಯನ್ನು ಎಲ್ಲಾ ಖಾತೆಗಳಿಗೂ ಉಳಿಸಿಕೊಳ್ಳಲು, ಮುಂದಿನ Remote Control ಆರಂಭದ ಮೊದಲು `ccam mcp import work <server> --force` ಮೂಲಕ ಅದನ್ನು ಗುಂಪಿಗೆ ಸ್ಥಳಾಂತರಿಸಿ.

### ನಕಲಿಸಿದ ಸರ್ವರ್‌ಗಳನ್ನು ತೆಗೆದುಹಾಕಬೇಕು

`ccam mcp remove <server>` ಸರ್ವರ್ ಅನ್ನು ಗುಂಪಿನಿಂದ ಮತ್ತು ಪ್ರತಿ ನಕಲಿನಿಂದ ತೆಗೆದುಹಾಕುತ್ತದೆ. ನಕಲಿನ ಮೊದಲಿನ ಸ್ಥಿತಿಗೆ ಮರಳಲು, ಖಾತೆಯ ಯಾವುದೇ ಸೆಷನ್ ಚಾಲನೆಯಲ್ಲಿಲ್ಲದಾಗ `~/.local/state/ccam/mcp/work/backups/` ನಿಂದ ಇತ್ತೀಚಿನ ಬ್ಯಾಕಪ್ ಅನ್ನು ಮರುಸ್ಥಾಪಿಸಿ.

### Remote Control ಫೋಲ್ಡರ್ ಅನ್ನು ನಂಬಬೇಕೇ ಎಂದು ಕೇಳುತ್ತದೆ, ಅಥವಾ ಟರ್ಮಿನಲ್ ಇಲ್ಲದೆ ನಿರ್ಗಮಿಸುತ್ತದೆ

ಫೋಲ್ಡರ್ ಮೇಲಿನ ನಂಬಿಕೆ ಪ್ರತಿ ಖಾತೆಗೂ ಪ್ರತ್ಯೇಕವಾಗಿ `.claude.json` ನಲ್ಲಿ ಸಂಗ್ರಹವಾಗುತ್ತದೆ, ಹೀಗಾಗಿ `personal` ಆಗಿ ನೀವು ನಂಬಿದ ಫೋಲ್ಡರ್ `work` ಗೆ ಹೊಸದು. ಆ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಒಮ್ಮೆ `ccam work` ರನ್ ಮಾಡಿ ಒಪ್ಪಿಕೊಳ್ಳಿ, ನಂತರ ಅಲ್ಲಿ Remote Control ಆರಂಭಿಸಿ.

### "You must be logged in with a Claude account that has a subscription"

Remote Control ಗೆ Claude ಚಂದಾದಾರಿಕೆಯ ಲಾಗಿನ್ ಬೇಕು. `ccam list` ನಲ್ಲಿ ಖಾತೆಯ ಇಮೇಲ್ ಪರಿಶೀಲಿಸಿ, ಮತ್ತು ಅಗತ್ಯವಿದ್ದರೆ `ccam login work` ಮೂಲಕ ಮತ್ತೆ ಸೈನ್ ಇನ್ ಮಾಡಿ.

### ಸೆಷನ್ ಫೋನ್‌ನಲ್ಲಿ ಕಾಣಿಸುವುದಿಲ್ಲ

Claude ಆ್ಯಪ್ ಅಥವಾ claude.ai/code ಬೇರೆ Claude ಖಾತೆಯಿಂದ ಸೈನ್ ಇನ್ ಆಗಿದೆ. `work` ಗಾಗಿ `ccam list` ತೋರಿಸುವ ಇಮೇಲ್‌ನೊಂದಿಗೆ ಅದನ್ನು ಹೋಲಿಸಿ.

### ಈ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಇನ್ನೊಂದು Claude Code ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳುತ್ತಿದೆ ಎಂದು ccam ಎಚ್ಚರಿಸುತ್ತದೆ

ಅದೇ ಇತಿಹಾಸ ಗುಂಪಿನ ಇನ್ನೊಂದು ಖಾತೆ ಅದೇ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಚಾಲನೆಯಲ್ಲಿದೆ. ಎರಡೂ ಖಾತೆಗಳಿಂದ ಒಂದೇ ಸಂಭಾಷಣೆಯನ್ನು ಮುಂದುವರಿಸುವ ಸೆಷನ್‌ಗಳು ಇತಿಹಾಸವನ್ನು ಬೆರೆಸಬಹುದು. ಅವು ಬೇರೆ ಬೇರೆ ಸಂಭಾಷಣೆಗಳಲ್ಲಿ ಕೆಲಸ ಮಾಡುತ್ತವೆ ಎಂದಾದರೆ ಮಾತ್ರ `y` ಎಂದು ಉತ್ತರಿಸಿ, ಅಥವಾ Remote Control ಅನ್ನು ಬೇರೆ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಆರಂಭಿಸಿ.

## ಸಂಬಂಧಿತ

- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು](/features/mcp)
- [ಖಾತೆಗಳ ನಡುವೆ ಹಂಚಿಕೆ](/features/share)
- [ಖಾತೆಗಳ ನಿರ್ವಹಣೆ](/features/accounts)
