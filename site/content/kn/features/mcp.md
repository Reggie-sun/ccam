---
title: ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು
description: ccam mcp add, add-json, list, remove ಮತ್ತು import ಮೂಲಕ ಪ್ರತಿ ಗುಂಪಿಗೆ MCP ಸರ್ವರ್‌ಗಳ ಒಂದೇ ಪಟ್ಟಿಯನ್ನು ಇಡಿ; ಅದು --mcp-config ಮೂಲಕ ರವಾನೆಯಾಗುತ್ತದೆ ಅಥವಾ ಪ್ರತಿ ಖಾತೆಗೆ ಸಿಂಕ್ ಆಗುತ್ತದೆ.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code ಬಳಕೆದಾರ MCP ಸರ್ವರ್‌ಗಳನ್ನು `.claude.json` ನಲ್ಲಿ ಇಡುತ್ತದೆ. ಖಾತೆಯ ಲಾಗಿನ್ ಕೂಡ ಇದೇ ಫೈಲ್‌ನಲ್ಲಿರುತ್ತದೆ, ಹೀಗಾಗಿ ಈ ಫೈಲ್ ಅನ್ನು ಖಾತೆಗಳ ನಡುವೆ ಹಂಚಿಕೊಳ್ಳಲು ಸಾಧ್ಯವಿಲ್ಲ. ccam ಹಂಚಿಕೊಂಡ ಸರ್ವರ್‌ಗಳನ್ನು `mcp.json` ಎಂಬ ಗುಂಪಿನ ಫೈಲ್‌ನಲ್ಲಿ ಇಡುತ್ತದೆ ಮತ್ತು ಆರಂಭದ ಸಮಯದಲ್ಲಿ ಅದನ್ನು ಗುಂಪಿನ ಪ್ರತಿ ಖಾತೆಗೆ ನೀಡುತ್ತದೆ. ನಿಮ್ಮ ಎಲ್ಲಾ ಖಾತೆಗಳಿಗೆ ಒಮ್ಮೆಯೇ ಸರ್ವರ್ ಸೇರಿಸಲು `ccam mcp` ಬಳಸಿ, ಮತ್ತು ಒಂದು ಸೆಷನ್‌ನೊಳಗೆ ನೀವು ಸೇರಿಸಿದ ಸರ್ವರ್ ಅನ್ನು ಇತರ ಖಾತೆಗಳಿಗೆ ಸ್ಥಳಾಂತರಿಸಲು `ccam mcp import` ಬಳಸಿ.

## ಬಳಕೆ

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

ಬದಲಾವಣೆ ಪ್ರತಿ ಖಾತೆಯ ಮುಂದಿನ ಆರಂಭದಿಂದ ಅನ್ವಯವಾಗುತ್ತದೆ; ಚಾಲನೆಯಲ್ಲಿರುವ ಸೆಷನ್‌ಗಳು ತಾವು ಆರಂಭವಾದಾಗಿನ ಸರ್ವರ್‌ಗಳನ್ನೇ ಉಳಿಸಿಕೊಳ್ಳುತ್ತವೆ:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

ಗುಂಪಿನಲ್ಲಿ ಈಗಾಗಲೇ ಇರುವ ಹೆಸರನ್ನು ಸೇರಿಸಿದರೆ ಅದರ ವ್ಯಾಖ್ಯಾನ ಬದಲಾಗುತ್ತದೆ, ಮತ್ತು `added` ಬದಲು `updated` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ. ಒಂದು ಖಾತೆ ಆ ಗುಂಪಿನೊಂದಿಗೆ MCP ಹಂಚಿಕೊಂಡಿದ್ದರೆ ಮಾತ್ರ (`ccam share <name> --mcp on`) ಸರ್ವರ್‌ಗಳು ಅದಕ್ಕೆ ತಲುಪುತ್ತವೆ; ಯಾವ ಖಾತೆಯೂ ಇನ್ನೂ MCP ಹಂಚಿಕೊಳ್ಳದಿದ್ದರೆ `add` ಅದನ್ನು ತಿಳಿಸುತ್ತದೆ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` ಮತ್ತು `remove` ಗಾಗಿ ಗುರಿ ಗುಂಪು | ಯಾವುದಾದರೂ ಖಾತೆ MCP ಹಂಚಿಕೊಳ್ಳುವ ಏಕೈಕ ಗುಂಪು, ಇಲ್ಲದಿದ್ದರೆ `default` |
| `-- <command> [args…]` (`add`) | ಈ ಕಮಾಂಡ್ ಅನ್ನು stdio ಸರ್ವರ್ ಆಗಿ ರನ್ ಮಾಡುತ್ತದೆ | |
| `--url <url>` (`add`) | ರಿಮೋಟ್ ಸರ್ವರ್‌ನ URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` ಸರ್ವರ್‌ನ ಟ್ರಾನ್ಸ್‌ಪೋರ್ಟ್ | `http` |
| `--header <'Name: value'>` (`add`) | HTTP ಹೆಡರ್, ಹಲವು ಬಾರಿ ನೀಡಬಹುದು | ಇಲ್ಲ |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio ಸರ್ವರ್‌ಗೆ ಎನ್ವಿರಾನ್‌ಮೆಂಟ್ ವೇರಿಯೇಬಲ್, ಹಲವು ಬಾರಿ ನೀಡಬಹುದು | ಇಲ್ಲ |
| `--force` (`import`) | ಬೇರೆಯಾಗಿರುವ ಅದೇ ಹೆಸರಿನ ಗುಂಪಿನ ಸರ್ವರ್ ಅನ್ನು ಬದಲಾಯಿಸುತ್ತದೆ; ಹಳೆಯ `mcp.json` ನ ಬ್ಯಾಕಪ್ ತೆಗೆಯಲಾಗುತ್ತದೆ | ಆಫ್ |
| `--keep-local` (`import`) | ಖಾತೆಯ `.claude.json` ನಲ್ಲಿರುವ ನಮೂದುಗಳನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ | ಆಫ್ |
| `-h`, `--help` | ಪ್ರತಿ ಕಮಾಂಡ್‌ನ ಸಹಾಯ | |
| `--dry-run` (ಜಾಗತಿಕ) | `mcp list` ಮಾತ್ರ ಇದನ್ನು ಪಾಲಿಸುತ್ತದೆ; `add`, `add-json`, `remove` ಮತ್ತು `import` ಇದನ್ನು ನಿರಾಕರಿಸುತ್ತವೆ ಮತ್ತು ಏನನ್ನೂ ಬದಲಾಯಿಸುವುದಿಲ್ಲ | ಆಫ್ |

`add` ಗೆ `--` ನಂತರ ಒಂದು ಕಮಾಂಡ್ ಅಥವಾ `--url` ಬೇಕು; ಇಲ್ಲದಿದ್ದರೆ ಅದು `a command to run (after --) or --url is required` ಎಂದು ವಿಫಲವಾಗುತ್ತದೆ. `--header` `--url` ಸರ್ವರ್‌ಗಳಿಗೆ ಮತ್ತು `--env` ಕಮಾಂಡ್ ಸರ್ವರ್‌ಗಳಿಗೆ ಅನ್ವಯವಾಗುತ್ತವೆ. `add-json` `claude mcp add-json` ನಂತೆಯೇ JSON ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ, ಮತ್ತು ವ್ಯಾಖ್ಯಾನವು `command` ಅಥವಾ `url` ಇರುವ ಆಬ್ಜೆಕ್ಟ್ ಆಗಿರಬೇಕು. `import` ಯಾವಾಗಲೂ ಖಾತೆ MCP ಹಂಚಿಕೊಳ್ಳುವ ಗುಂಪನ್ನೇ ಬಳಸುತ್ತದೆ, ಮತ್ತು `list` ಎಲ್ಲಾ ಗುಂಪುಗಳನ್ನು ತೋರಿಸುತ್ತದೆ; ಹೀಗಾಗಿ `-g` ಅವುಗಳನ್ನು ಬದಲಾಯಿಸುವುದಿಲ್ಲ.

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

ಒಳಗೊಂಡಿರುವ ಫೈಲ್‌ಗಳು:

- `~/.local/share/ccam/groups/<group>/mcp.json`: ಹಂಚಿಕೊಂಡ ಸರ್ವರ್‌ಗಳು, `{"mcpServers": {…}}`, ಅಂದರೆ `--mcp-config` ಓದುವ ಸ್ವರೂಪ. ಸರ್ವರ್ ನಮೂದುಗಳಲ್ಲಿ ಹಲವೊಮ್ಮೆ API ಕೀಗಳಿರುವುದರಿಂದ, ಇದನ್ನು `fsutil.AtomicWrite` ಮೂಲಕ `0600` ಮೋಡ್‌ನಲ್ಲಿ ಬರೆಯಲಾಗುತ್ತದೆ.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ಖಾತೆಯ ಸರ್ವರ್‌ಗಳಿರುವ Claude Code ನ ಸ್ವಂತ ಫೈಲ್. ccam `profile.TopLevelValue` ಸ್ಕ್ಯಾನರ್ ಮೂಲಕ ಇದರ `mcpServers` ಮೌಲ್ಯವನ್ನು ಮಾತ್ರ ಓದುತ್ತದೆ, ಮತ್ತು ಆರಂಭದ ಹಾದಿಯಲ್ಲಿ ಎಂದಿಗೂ ಇಡೀ ಫೈಲ್ ಅನ್ನು ಪಾರ್ಸ್ ಮಾಡುವುದಿಲ್ಲ.
- `~/.local/state/ccam/mcp/<name>/`: `seen` ಫೈಲ್, `sync-stamp` ಮತ್ತು `.claude.json` ಬ್ಯಾಕಪ್‌ಗಳು.

ಎರಡು ವಿಧಾನಗಳಿವೆ, ಅವುಗಳನ್ನು `~/.config/ccam/config.toml` ನಲ್ಲಿನ `mcp_strategy` ಹೊಂದಿಸುತ್ತದೆ. ಡೀಫಾಲ್ಟ್ ಫ್ಲ್ಯಾಗ್ ವಿಧಾನ:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` ಖಾತೆಯ MCP ಗುಂಪನ್ನು ಕಂಡುಹಿಡಿಯುತ್ತದೆ. `mcp.FlagArgs` ಗುಂಪಿನ `mcp.json` ಇದೆಯೇ ಎಂದು ಪಾರ್ಸ್ ಮಾಡದೆ ಒಂದೇ `stat` ಮೂಲಕ ಪರಿಶೀಲಿಸುತ್ತದೆ.
2. `Plan.InsertArgs` `--mcp-config=<path>` ಅನ್ನು `claude` ನ ತಕ್ಷಣ ನಂತರ, ನಿಮ್ಮ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿಗಿಂತ ಮೊದಲು ಇಡುತ್ತದೆ. ಹೀಗಾಗಿ ನಿಮ್ಮದೇ `--mcp-config` ಅಥವಾ `--strict-mcp-config` ಅದರ ಮೇಲೆ ಅನ್ವಯವಾಗುತ್ತವೆ.
3. `mcp.NewLocal` ಕೊನೆಯ ಆರಂಭದ ನಂತರ ಖಾತೆಯ `.claude.json` ಬದಲಾಗಿದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುತ್ತದೆ, ಬದಲಾಗಿದ್ದರೆ ಮಾತ್ರ ಅದರ ಸರ್ವರ್‌ಗಳನ್ನು ಪಾರ್ಸ್ ಮಾಡುತ್ತದೆ.
4. `syscall.Exec` ccam ನ ಜಾಗದಲ್ಲಿ Claude Code ಅನ್ನು ತರುತ್ತದೆ.

ಪಾತ್ ಅನ್ನು ಎರಡು ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಾಗಿ ಅಲ್ಲ, `--mcp-config=<path>` ಎಂಬ ಒಂದೇ ಆರ್ಗ್ಯುಮೆಂಟ್ ಆಗಿ ರವಾನಿಸಲಾಗುತ್ತದೆ. Claude Code ನ `--mcp-config` ಹಲವು ಮೌಲ್ಯಗಳನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ, ಹೀಗಾಗಿ `claude --mcp-config a.json "prompt"` ಪ್ರಾಂಪ್ಟ್ ಅನ್ನು ಎರಡನೇ ಕಾನ್ಫಿಗ್ ಫೈಲ್ ಎಂದು ಓದಿ `MCP config file not found` ಎಂದು ವಿಫಲವಾಗುತ್ತದೆ. ಒಂದೇ ಆರ್ಗ್ಯುಮೆಂಟ್‌ನ ರೂಪ ಅಲ್ಲಿಗೇ ನಿಲ್ಲುತ್ತದೆ, ಮತ್ತು ನಂತರ ಬರುವ ಪ್ರಾಂಪ್ಟ್ ಅಥವಾ `mcp list` ನಂತಹ ಉಪ-ಕಮಾಂಡ್ ಸರಿಯಾಗಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ. Claude Code `--mcp-config` ಸರ್ವರ್‌ಗಳನ್ನು ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್‌ಗಳ ಜೊತೆಗೆ ಲೋಡ್ ಮಾಡುತ್ತದೆ; ಹೆಸರು ಒಂದೇ ಆದಾಗ `.claude.json` ನಲ್ಲಿನ ನಮೂದು ಗೆಲ್ಲುತ್ತದೆ. ಫ್ಲ್ಯಾಗ್ ವಿಧಾನ Remote Control ಗಾಗಿ ಮಾತ್ರ `.claude.json` ಗೆ ಬರೆಯುತ್ತದೆ.

`claude remote-control` (`rc`) ಯಾವುದೇ `--mcp-config` ತೆಗೆದುಕೊಳ್ಳುವುದಿಲ್ಲ, ಮತ್ತು ಅದರ ಸೆಷನ್‌ಗಳು ಸರ್ವರ್‌ಗಳನ್ನು `.claude.json` ನಿಂದ ಮಾತ್ರ ಓದುತ್ತವೆ. `ccam work rc` ಗಾಗಿ, ccam ಫ್ಲ್ಯಾಗ್ ಅನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ ಮತ್ತು ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ, ಗುಂಪಿನ ಸರ್ವರ್‌ಗಳನ್ನು ಬ್ಯಾಕಪ್‌ನೊಂದಿಗೆ `.claude.json` ಗೆ ನಕಲಿಸಬೇಕೇ ಎಂದು ಕೇಳುತ್ತದೆ, `sync` ವಿಧಾನ ಮಾಡುವಂತೆಯೇ. ನಕಲು ಆಗುವವರೆಗೆ ಅದು ಕೇಳುತ್ತಲೇ ಇರುತ್ತದೆ, ನಂತರ ಗುಂಪಿನ `mcp.json` ಬದಲಾದ ಮೇಲೆ ಮಾತ್ರ ಮತ್ತೆ ಕೇಳುತ್ತದೆ ([Remote Control](/features/launch#remote-control)). ಹೆಸರು ಒಂದೇ ಆದಾಗ ಅಂತಹ ನಕಲು `--mcp-config` ಮೇಲೆ ಗೆಲ್ಲುವುದರಿಂದ, `ccam mcp add`, `add-json`, `remove` ಮತ್ತು `import` ಗುಂಪಿನ ಇತರ ಖಾತೆಗಳಲ್ಲಿರುವ ನಕಲುಗಳನ್ನು ನವೀಕರಿಸುತ್ತವೆ:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

ಸೆಷನ್‌ನೊಳಗೆ `claude mcp add` ಮೂಲಕ ಸೇರಿಸಿದ ಸರ್ವರ್ ಆ ಖಾತೆಯ `.claude.json` ನಲ್ಲಿ ಮಾತ್ರ ಸೇರುತ್ತದೆ. Claude Code ccam ನ ಜಾಗವನ್ನು ಪಡೆದಿರುವುದರಿಂದ ಆ ಸಮಯದಲ್ಲಿ ccam ಅದನ್ನು ನೋಡಲಾರದು, ಹೀಗಾಗಿ ಮುಂದಿನ ಆರಂಭದಲ್ಲಿ ಪರಿಶೀಲಿಸುತ್ತದೆ. `seen` ಫೈಲ್ `.claude.json` ನ ಮಾರ್ಪಾಡಿನ ಸಮಯ ಮತ್ತು ಈಗಾಗಲೇ ತಿಳಿಸಿದ ಸರ್ವರ್ ಹೆಸರುಗಳನ್ನು ಇಡುತ್ತದೆ. ಸಮಯ ಬೇರೆಯಾಗಿದ್ದರೆ, ccam ಹೆಸರುಗಳನ್ನು ಗುಂಪಿನೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ ಮತ್ತು ಪ್ರತಿ ಹೊಸ ಸರ್ವರ್‌ಗೆ ಒಮ್ಮೆ, ಇನ್‌ಪುಟ್‌ಗಾಗಿ ಕಾಯದೆ, ಒಂದು ಸಲಹೆಯನ್ನು ಮುದ್ರಿಸುತ್ತದೆ:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` ಇದ್ದಾಗ, ಸರ್ವರ್‌ಗಳು `claude mcp list` ಮತ್ತು `/mcp` ನಲ್ಲಿ ಸಾಮಾನ್ಯ ಬಳಕೆದಾರ ಸರ್ವರ್‌ಗಳಾಗಿ ಕಾಣಿಸುತ್ತವೆ, ಮತ್ತು ccam ಆರಂಭಕ್ಕೆ ಮೊದಲು ಅವುಗಳನ್ನು ಪ್ರತಿ ಖಾತೆಯ `.claude.json` ಗೆ ಬರೆಯುತ್ತದೆ:

<!-- motion:mcp-sync -->

1. `mcp.Sync` ಗುಂಪಿನ `mcp.json` ನ ಮಾರ್ಪಾಡಿನ ಸಮಯವನ್ನು `sync-stamp` ಜೊತೆ ಹೋಲಿಸುತ್ತದೆ. ಅದು ಬದಲಾಗಿಲ್ಲದಿದ್ದರೆ ಖಾತೆಯ `.claude.json` ಅನ್ನು ತೆರೆಯುವುದೇ ಇಲ್ಲ.
2. ಅದು Claude Code ನ ಸ್ವಂತ ಲಾಕ್ ಅನ್ನು ಪಡೆಯುತ್ತದೆ: ಫೈಲ್‌ನ ಪಕ್ಕದಲ್ಲಿ `mkdir` ಮೂಲಕ ರಚಿಸಿದ `.claude.json.lock` ಫೋಲ್ಡರ್. ಇದು ಗರಿಷ್ಠ 2 ಸೆಕೆಂಡ್ ಕಾಯುತ್ತದೆ, ಮತ್ತು Claude Code ನ ಲಾಕ್ ಲೈಬ್ರರಿಯಂತೆಯೇ 10 ಸೆಕೆಂಡ್‌ಗಿಂತ ಹಳೆಯ ಲಾಕ್ ಅನ್ನು ಅವಧಿ ಮೀರಿದ್ದು ಎಂದು ಪರಿಗಣಿಸುತ್ತದೆ.
3. ಈಗಿನ `.claude.json` ಅನ್ನು ಸ್ಟೇಟ್ ಫೋಲ್ಡರ್‌ನ `backups/` ನಲ್ಲಿ ಉಳಿಸುತ್ತದೆ ಮತ್ತು ಇತ್ತೀಚಿನ 5 ಅನ್ನು ಇಟ್ಟುಕೊಳ್ಳುತ್ತದೆ.
4. ಪ್ರತಿ ಗುಂಪಿನ ಸರ್ವರ್ ಅನ್ನು `mcpServers` ನಲ್ಲಿ ಹೊಂದಿಸುತ್ತದೆ, ತಾನು ಹಿಂದೆ ಬರೆದಿದ್ದ ಆದರೆ ಈಗ ಗುಂಪಿನಲ್ಲಿ ಇಲ್ಲದ ಸರ್ವರ್‌ಗಳನ್ನು ತೆಗೆದುಹಾಕುತ್ತದೆ, ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್‌ಗಳನ್ನು ಉಳಿಸುತ್ತದೆ, ಮತ್ತು ಫೈಲ್ ಅನ್ನು ಅದರ ಮೂಲ ಮೋಡ್‌ನೊಂದಿಗೆ ಅಟಾಮಿಕ್ ಆಗಿ ಬರೆಯುತ್ತದೆ.
5. ಹೊಸ ಸ್ಟ್ಯಾಂಪ್ ಮತ್ತು ತಾನು ನಿರ್ವಹಿಸುವ ಹೆಸರುಗಳನ್ನು ದಾಖಲಿಸುತ್ತದೆ, ಮತ್ತು Claude Code `--mcp-config` ಇಲ್ಲದೆ ಆರಂಭವಾಗುತ್ತದೆ.

ಖಾತೆಯ ಸೆಷನ್ ಚಾಲನೆಯಲ್ಲಿರುವಾಗ ಸಿಂಕ್ ಅನ್ನು `ccam: account work is already running, skipping MCP sync` ಸಂದೇಶದೊಂದಿಗೆ ಬಿಟ್ಟುಬಿಡಲಾಗುತ್ತದೆ, ಏಕೆಂದರೆ Claude Code ಚಾಲನೆಯಲ್ಲಿರುವಾಗ `.claude.json` ಅನ್ನು ಮತ್ತೆ ಬರೆಯುತ್ತದೆ. ವಿಫಲವಾದ ಸಿಂಕ್ `ccam: MCP sync failed (launching anyway): …` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ ಮತ್ತು ಹೇಗಿದ್ದರೂ ಆರಂಭಿಸುತ್ತದೆ. ಸೆಷನ್‌ನೊಳಗೆ ಸೇರಿಸಿದ ಸರ್ವರ್‌ಗಳ ಸಲಹೆ ಫ್ಲ್ಯಾಗ್ ವಿಧಾನಕ್ಕೆ ಮಾತ್ರ ಸೇರಿದ್ದು.

`ccam mcp import` ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್‌ಗಳನ್ನು ಅದರ ಗುಂಪಿಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ:

<!-- motion:mcp-import -->

ccam `mcp.LocalServers` ಮೂಲಕ ಖಾತೆಯ ಸರ್ವರ್‌ಗಳನ್ನು ಓದುತ್ತದೆ, ಮತ್ತು ಕೀಗಳ ಕ್ರಮ ಹಾಗೂ ಸ್ಪೇಸಿಂಗ್ ಅನ್ನು ನಿರ್ಲಕ್ಷಿಸುವ `mcp.SameDef` ಬಳಸಿ ಪ್ರತಿಯೊಂದನ್ನು ಗುಂಪಿನ ಅದೇ ಹೆಸರಿನ ಸರ್ವರ್‌ನೊಂದಿಗೆ ಹೋಲಿಸುತ್ತದೆ. ಗುಂಪಿನಲ್ಲಿ ಇಲ್ಲದ ಅಥವಾ ಅದೇ ವ್ಯಾಖ್ಯಾನದೊಂದಿಗೆ ಇರುವ ಸರ್ವರ್ ಅನ್ನು ನಕಲಿಸಲಾಗುತ್ತದೆ. ಬೇರೆಯಾಗಿರುವುದನ್ನು ಬಿಟ್ಟುಬಿಡಲಾಗುತ್ತದೆ, ಏಕೆಂದರೆ ಉಳಿದ ಪ್ರತಿ ಖಾತೆಯೂ ಬಳಸುವುದು ಗುಂಪಿನ ವ್ಯಾಖ್ಯಾನವನ್ನೇ:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` ಇದ್ದಾಗ ಗುಂಪಿನ ಸರ್ವರ್ ಬದಲಾಗುತ್ತದೆ, ಮತ್ತು ಹಿಂದಿನ `mcp.json` ಅನ್ನು ಮೊದಲು ಗುಂಪಿನ ಫೋಲ್ಡರ್‌ನೊಳಗಿನ `.ccam-backup/<time>/mcp.json` ನಲ್ಲಿ ಉಳಿಸಲಾಗುತ್ತದೆ. ನಂತರ `mcp.RemoveLocal` ಅದೇ ಲಾಕ್ ಅಡಿಯಲ್ಲಿ, ಸ್ಟೇಟ್ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಬ್ಯಾಕಪ್ ತೆಗೆದ ನಂತರ, ನಕಲಿಸಿದ ನಮೂದುಗಳನ್ನು ಖಾತೆಯ `.claude.json` ನಿಂದ ತೆಗೆದುಹಾಕುತ್ತದೆ. ಹೆಸರು ಒಂದೇ ಆದಾಗ ಖಾತೆಯ ನಮೂದು ಗುಂಪಿನದನ್ನು ಮೀರಿ ಗೆಲ್ಲುವುದರಿಂದ ಇದು ಅಗತ್ಯ. `--keep-local` ಅವುಗಳನ್ನು ಹಾಗೆಯೇ ಬಿಡುತ್ತದೆ.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- **ಸೈನ್-ಇನ್‌ಗಳು.** OAuth ಮೂಲಕ ಸೈನ್ ಇನ್ ಆಗುವ ಸರ್ವರ್‌ಗಳು ತಮ್ಮ ಟೋಕನ್‌ಗಳನ್ನು ಸರ್ವರ್ ಪಟ್ಟಿಯಲ್ಲಿ ಅಲ್ಲ, ಪ್ರತಿ ಖಾತೆಯ `.credentials.json` (`mcpOAuth`) ನಲ್ಲಿ ಇಡುತ್ತವೆ. ಹೀಗಾಗಿ ಸರ್ವರ್ ಹಂಚಿಕೊಂಡರೆ ಅದರ ಸೈನ್-ಇನ್ ಹಂಚಿಕೆಯಾಗುವುದಿಲ್ಲ. ಪ್ರತಿ ಖಾತೆಯಲ್ಲಿ `/mcp` ಮೂಲಕ ಒಮ್ಮೆ ಸೈನ್ ಇನ್ ಆಗಿ, ಅಥವಾ [`ccam migrate`](/features/migrate) ಮೂಲಕ `~/.claude` ನಿಂದ ಸೈನ್-ಇನ್‌ಗಳನ್ನು ನಕಲಿಸಿ. ಗುಂಪಿಗೆ ಎಂದಿಗೂ ಲಾಗಿನ್ ಫೈಲ್ ಸಿಗುವುದಿಲ್ಲ.
- **ಚಾಲನೆಯಲ್ಲಿರುವ ಖಾತೆಗಳು.** `import` ಗುಂಪಿಗೆ ನಕಲಿಸುತ್ತದೆ, ಆದರೆ ಚಾಲನೆಯಲ್ಲಿರುವ ಖಾತೆಯ ಸ್ವಂತ ನಮೂದುಗಳನ್ನು ಹಾಗೆಯೇ ಬಿಟ್ಟು ಅದನ್ನು ತಿಳಿಸುತ್ತದೆ; ಸೆಷನ್ ಮುಗಿದ ನಂತರ ಮತ್ತೆ ರನ್ ಮಾಡಿ. ಸಿಂಕ್ ಚಾಲನೆಯಲ್ಲಿರುವ ಖಾತೆಗಳನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.
- **Remote Control ಗಾಗಿ ನಕಲುಗಳು.** ಚಾಲನೆಯಲ್ಲಿರುವ ಖಾತೆಯ ನಕಲನ್ನು ಮುಟ್ಟಲಾಗುವುದಿಲ್ಲ: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` ಇದ್ದಾಗ ಏನನ್ನೂ ನವೀಕರಿಸುವ ಅಗತ್ಯವಿಲ್ಲ, ಏಕೆಂದರೆ ಪ್ರತಿ ಆರಂಭವೂ ಸಿಂಕ್ ಮಾಡುತ್ತದೆ.
- **ಹಾಳಾದ ಫೈಲ್‌ಗಳು.** `.claude.json` ಮಾನ್ಯ JSON ಅಲ್ಲದಿದ್ದರೆ, ಸಿಂಕ್ ಅದನ್ನು ಮುಟ್ಟದೆ ದೋಷವನ್ನು ತಿಳಿಸುತ್ತದೆ. 2 ಸೆಕೆಂಡ್‌ಗಳಲ್ಲಿ ಲಾಕ್ ಸಿಗದಿದ್ದರೆ, ccam `could not lock … (Claude Code is writing the file)` ಎಂದು ತಿಳಿಸುತ್ತದೆ.
- **MCP ಹಂಚಿಕೊಳ್ಳದ ಖಾತೆಗಳು.** `import` ನಿರಾಕರಿಸುತ್ತದೆ ಮತ್ತು ಮೊದಲು `ccam share side --mcp on` ರನ್ ಮಾಡಲು ಹೇಳುತ್ತದೆ.
- **ಗುಂಪಿನ ಹೆಸರುಗಳು.** `-g` ಗುಂಪಿನ ಹೆಸರನ್ನು ಮಾತ್ರ ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ; ಅದು `groups/` ಅಡಿಯಲ್ಲಿ ಒಂದು ಫೋಲ್ಡರ್ ಆಗುತ್ತದೆ, ಹೀಗಾಗಿ `../x` ನಂತಹ ಹೆಸರುಗಳನ್ನು ನಿರಾಕರಿಸಲಾಗುತ್ತದೆ.

## ಸಂಬಂಧಿತ

- [ಖಾತೆಗಳ ನಡುವೆ ಹಂಚಿಕೆ](/features/share)
- [MCP ಸರ್ವರ್‌ಗಳು, ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮತ್ತು ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳಿ](/docs/sharing)
- [ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು](/features/migrate)
- [ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು](/docs/commands)
