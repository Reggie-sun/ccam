---
title: ಕಮಾಂಡ್ ಉಲ್ಲೇಖ
metaTitle: ccam ಕಮಾಂಡ್ ಉಲ್ಲೇಖ
description: ccam add ಮತ್ತು ccam <name> ನಿಂದ ccam usage, ccam next ಮತ್ತು ccam doctor ವರೆಗೆ, ccam ನ ಎಲ್ಲಾ ಕಮಾಂಡ್‌ಗಳು ಮತ್ತು ಆಯ್ಕೆಗಳು.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

ಕಮಾಂಡ್ ಏನು ಮಾಡುತ್ತದೆ ಎಂಬುದನ್ನು ನಿಜವಾಗಿ ಮಾಡದೆಯೇ ನೋಡಲು, ಅದರ ಮುಂದೆ `--dry-run` ಹಾಕಿ. ಉದಾಹರಣೆಗೆ `ccam --dry-run work`. ಇದು Claude Code ಆರಂಭಿಸುವಾಗ (`ccam <name>`, `ccam run`, `ccam next`) ಮತ್ತು `share`, `migrate`, `shortcut`, `update` ಹಾಗೂ `doctor --fix` ಜೊತೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ. ಏನನ್ನಾದರೂ ಬದಲಾಯಿಸುವ ಇತರ ಕಮಾಂಡ್‌ಗಳು ಚಾಲನೆಯಾಗುವ ಬದಲು ಇದನ್ನು ನಿರಾಕರಿಸುತ್ತವೆ.

<!-- commands -->

## Claude Code ಆರಂಭಿಸುವುದು

`ccam <name|number> [arguments]` ಆ ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಹೆಸರಿನ ನಂತರದ ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳನ್ನು ccam ಓದುವುದಿಲ್ಲ; ಅದು ಅವುಗಳನ್ನು ಅದೇ ಕ್ರಮದಲ್ಲಿ `claude` ಗೆ ರವಾನಿಸುತ್ತದೆ. ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್, ಎಕ್ಸಿಟ್ ಕೋಡ್ ಮತ್ತು ಸಿಗ್ನಲ್‌ಗಳು `claude` ನಲ್ಲಿ ಇರುವಂತೆಯೇ ವರ್ತಿಸುತ್ತವೆ. ಹೆಸರಿನ ನಂತರ `--` ಐಚ್ಛಿಕ.

ಆರ್ಗ್ಯುಮೆಂಟ್‌ಗಳಿದ್ದು ಖಾತೆಯ ಹೆಸರು ಇಲ್ಲದ `ccam`, ಉದಾಹರಣೆಗೆ `ccam -p "summarise"`, ಈಗಿನ ಫೋಲ್ಡರ್‌ಗೆ ಹೊಂದಿಸಿದ ಖಾತೆಯನ್ನು, ಇಲ್ಲದಿದ್ದರೆ ಡೀಫಾಲ್ಟ್ ಖಾತೆಯನ್ನು ಬಳಸುತ್ತದೆ. ಅವುಗಳನ್ನು `ccam folder <name> [folder]` ಮತ್ತು `ccam use <name>` ಮೂಲಕ ಹೊಂದಿಸಿ.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

ಖಾತೆಯ ಹೆಸರಿನ ಮುಂದೆ ಇರುವ `--fallback`, `-p`/`--print` ರನ್‌ಗಳಿಗಾಗಿ. ಖಾತೆ ತನ್ನ ಬಳಕೆಯ ಮಿತಿ ತಲುಪಿದಾಗ, ccam ಅದೇ ಕಮಾಂಡ್ ಅನ್ನು ಹೆಚ್ಚು ಅವಕಾಶ ಇರುವ ಲಾಗಿನ್ ಆದ ಖಾತೆಯಾಗಿ ಮತ್ತೆ ರನ್ ಮಾಡುತ್ತದೆ ಮತ್ತು ಪೂರ್ಣಗೊಂಡ ರನ್‌ನ ಔಟ್‌ಪುಟ್ ಅನ್ನು ಮಾತ್ರ ಪ್ರಿಂಟ್ ಮಾಡುತ್ತದೆ. [ಫೋಲ್ಡರ್‌ವಾರು ಖಾತೆಗಳು](/features/folder) ಮತ್ತು [ಮಿತಿಯಲ್ಲಿ ಮತ್ತೆ ರನ್ ಮಾಡುವುದು](/features/next#run-again-at-the-limit) ನೋಡಿ.

`ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಕೂಡ ತನ್ನ ಖಾತೆಯನ್ನು ಇದೇ ರೀತಿ ಆರಂಭಿಸುತ್ತದೆ: `ccam work --resume` ಮಾಡುವುದನ್ನೇ `ccwork --resume` ಮಾಡುತ್ತದೆ. [ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು](/docs/commands#shortcuts) ನೋಡಿ.

## ಖಾತೆಗಳನ್ನು ಸೇರಿಸುವುದು ಮತ್ತು ತೆಗೆದುಹಾಕುವುದು

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

ಫೋಲ್ಡರ್ ನೀಡದೆ `ccam import` ರನ್ ಮಾಡಿದರೆ, ಅದು `~/.claude` ಮತ್ತು `~/.claude-*` ಗಾಗಿ ಹುಡುಕುತ್ತದೆ ಮತ್ತು ಪ್ರತಿಯೊಂದರ ಬಗ್ಗೆ ಕೇಳುತ್ತದೆ. ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಫೋಲ್ಡರ್‌ಗಳು ಇದ್ದಲ್ಲೇ ನೋಂದಣಿಯಾಗುತ್ತವೆ; ccam ಎಂದಿಗೂ ಲಾಗಿನ್ ಅನ್ನು ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ ಅಥವಾ ನಕಲಿಸುವುದಿಲ್ಲ.

`ccam add --api` Claude ಲಾಗಿನ್ ಬದಲು API ಕೀ ಬಳಸುವ ಖಾತೆಯನ್ನು ಸೇರಿಸುತ್ತದೆ; ccam ಪ್ರತಿ ಆರಂಭದಲ್ಲೂ ಕೀಯನ್ನು ಓದುತ್ತದೆ ಮತ್ತು ಅದನ್ನು ಎಂದಿಗೂ ಸಂಗ್ರಹಿಸುವುದಿಲ್ಲ. [API ಕೀ ಖಾತೆಗಳು](/features/api) ನೋಡಿ.

ಒಂದು ಖಾತೆಯಾಗಿ Claude Code ಚಾಲನೆಯಲ್ಲಿರುವಾಗ, `ccam rename`, `ccam remove` ಮತ್ತು `ccam share` ಆ ಖಾತೆಯನ್ನು ಬದಲಾಯಿಸಲು ನಿರಾಕರಿಸುತ್ತವೆ. ಮೊದಲು ಸೆಷನ್‌ನಿಂದ ನಿರ್ಗಮಿಸಿ.

`ccam list --json` ಪ್ರತಿ ಖಾತೆಯನ್ನು ಅದರ ಸ್ಲಾಟ್, ಶಾರ್ಟ್‌ಕಟ್, ಇಮೇಲ್ ಮತ್ತು ಫೋಲ್ಡರ್ ಜೊತೆ ಮುದ್ರಿಸುತ್ತದೆ. ಅದರ `share` ಆಬ್ಜೆಕ್ಟ್ ಪ್ರತಿ ಐಟಂ (`mcp`, `settings`, `memory`, `extensions`, `history`) ಅನ್ನು ಗುಂಪಿನ ಹೆಸರಿಗೆ ಅಥವಾ `local` ಗೆ ಹೊಂದಿಸುತ್ತದೆ.

## ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

ಶಾರ್ಟ್‌ಕಟ್ ಎಂದರೆ ಒಂದು ಖಾತೆಯಾಗಿ Claude Code ಅನ್ನು ಆರಂಭಿಸುವ, ನೀವೇ ಹೆಸರಿಡುವ ಕಮಾಂಡ್. `ccam work` ನಂತೆಯೇ, ಇದು ಪ್ರತಿ ಆರ್ಗ್ಯುಮೆಂಟ್ ಅನ್ನು ಬದಲಾವಣೆ ಇಲ್ಲದೆ Claude Code ಗೆ ರವಾನಿಸುತ್ತದೆ. ಹೀಗಾಗಿ `ccwork --resume` ಮತ್ತು `ccwork -p "hi"` ನೀವು ನಿರೀಕ್ಷಿಸಿದಂತೆಯೇ ಕೆಲಸ ಮಾಡುತ್ತವೆ. ಇದು `~/.local/bin` ನಲ್ಲಿರುವ ccam ಗೆ ಒಂದು ಲಿಂಕ್.

ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು ಐಚ್ಛಿಕ. `ccam add` ವಿಜಾರ್ಡ್ ಒಂದನ್ನು ರಚಿಸಬೇಕೇ ಎಂದು ಕೇಳುತ್ತದೆ, `ccam add <name> --shortcut <command>` ಕೇಳದೆಯೇ ಒಂದನ್ನು ರಚಿಸುತ್ತದೆ, ಮತ್ತು `ccam list` ಪ್ರತಿ ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು, ಇಲ್ಲದಿದ್ದರೆ `-` ಅನ್ನು ತೋರಿಸುತ್ತದೆ.

ಹೆಸರುಗಳಲ್ಲಿ a-z, 0-9, `.`, `-` ಮತ್ತು `_` ಅಕ್ಷರಗಳನ್ನು ಬಳಸಬಹುದು. ಹೆಸರು ccam ಕಮಾಂಡ್ ಆಗಿರಬಾರದು, ಇನ್ನೊಂದು ಖಾತೆಯ ಶಾರ್ಟ್‌ಕಟ್ ಆಗಿರಬಾರದು, ಅಥವಾ ನಿಮ್ಮ `PATH` ನಲ್ಲಿರುವ ಇನ್ನೊಂದು ಪ್ರೋಗ್ರಾಂನ ಹೆಸರಾಗಿರಬಾರದು.

`ccam rename` ಶಾರ್ಟ್‌ಕಟ್ ಹೆಸರನ್ನು ಹಾಗೆಯೇ ಇಡುತ್ತದೆ, ಮತ್ತು `ccam remove` ಖಾತೆಯ ಜೊತೆಗೆ ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನೂ ಅಳಿಸುತ್ತದೆ. `ccam shims install` ಎಲ್ಲಾ ಖಾತೆಗಳ ಲಿಂಕ್‌ಗಳನ್ನು ಮತ್ತೆ ರಚಿಸುತ್ತದೆ, ಉದಾಹರಣೆಗೆ ನೀವು ccam ಅನ್ನು ಬೇರೆಡೆಗೆ ಸ್ಥಳಾಂತರಿಸಿದ ನಂತರ.

## ಹಂಚಿಕೆ

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ಪ್ರತಿ ಐಟಂ `on` (`default` ಗುಂಪು), `off`, ಅಥವಾ ಒಂದು ಗುಂಪಿನ ಹೆಸರನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ. [ಹಂಚಿಕೆ](/docs/sharing) ನೋಡಿ.

## ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` ಒಂದು Claude Code ಫೋಲ್ಡರ್‌ನಿಂದ `settings.json`, `CLAUDE.md`, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು ಬಳಕೆದಾರ MCP ಸರ್ವರ್‌ಗಳನ್ನು ಒಂದು ಗುಂಪಿಗೆ ನಕಲಿಸುತ್ತದೆ. ಹೀಗಾಗಿ ಆ ಐಟಂಗಳನ್ನು ಹಂಚಿಕೊಳ್ಳುವ ಪ್ರತಿ ಖಾತೆಗೂ ಅವು ಸಿಗುತ್ತವೆ. ಫೋಲ್ಡರ್ ಅನ್ನು ಓದಲಾಗುತ್ತದೆ ಅಷ್ಟೇ. `--settings`, `--memory`, `--extensions`, `--mcp` ಮತ್ತು `--history` ಆಯ್ಕೆಗಳು ನೀವು ಹೆಸರಿಸಿದ ಐಟಂಗಳನ್ನು ಮಾತ್ರ ನಕಲಿಸುತ್ತವೆ. ಅವು ಇಲ್ಲದಿದ್ದರೆ, ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಒಂದು ಚೆಕ್‌ಲಿಸ್ಟ್ ನಿಮಗೆ ಆಯ್ಕೆಮಾಡಲು ಬಿಡುತ್ತದೆ, ಅದರಲ್ಲಿ ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವನ್ನೂ ಗುರುತಿಸಲಾಗಿರುತ್ತದೆ; ಬೇರೆಡೆ ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವೂ ನಕಲಾಗುತ್ತದೆ. `-g` ಅಥವಾ `--group` ಗುಂಪನ್ನು ಆಯ್ಕೆಮಾಡುತ್ತದೆ, ಮತ್ತು `-y` ಅಥವಾ `--yes` ಕೇಳದೆ ನಕಲಿಸುತ್ತದೆ ಹಾಗೂ ಫೈಲ್ ಬೇರೆಯಾಗಿದ್ದಾಗ ಗುಂಪಿನ ನಕಲನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ. ಇದು ಸರ್ವರ್‌ಗಳ ಉಳಿಸಿದ ಸೈನ್-ಇನ್‌ಗಳನ್ನು ಪ್ರತಿ ಖಾತೆಯ ಸ್ವಂತ ಲಾಗಿನ್ ಫೈಲ್‌ಗೂ ನಕಲಿಸುತ್ತದೆ. [ಈಗಿರುವ ಸೆಟಪ್ ಅನ್ನು ನಕಲಿಸುವುದು](/docs/sharing#copy-an-existing-setup) ನೋಡಿ.

## MCP ಸರ್ವರ್‌ಗಳು

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` ಆ ಖಾತೆಯಲ್ಲಿ ಮಾತ್ರ ಇರುವ ಸರ್ವರ್‌ಗಳನ್ನು ಅದರ ಗುಂಪಿಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ. ಗುಂಪು ಈಗಾಗಲೇ ಬೇರೆ ವ್ಯಾಖ್ಯಾನಕ್ಕಾಗಿ ಬಳಸುತ್ತಿರುವ ಹೆಸರಿನ ಸರ್ವರ್ ಅನ್ನು ಬಿಟ್ಟುಬಿಡಲಾಗುತ್ತದೆ; ಗುಂಪಿನ ನಕಲನ್ನು ಬದಲಾಯಿಸಲು `--force` ಸೇರಿಸಿ, ccam ಮೊದಲು ಅದರ ಬ್ಯಾಕಪ್ ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ.

## ಬಳಕೆ ಮತ್ತು ಬದಲಾವಣೆ

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` ಮತ್ತು `ccam pick` ಒಂದೇ ರೀತಿಯಲ್ಲಿ ಆಯ್ಕೆಮಾಡುತ್ತವೆ: ಮೊದಲು `auto_max_usage` (`config.toml`, ಡೀಫಾಲ್ಟ್ 100) ಗಿಂತ ಕೆಳಗಿರುವ ಖಾತೆಗಳು, `auto_strategy` (`headroom` ಅಥವಾ `expiring`) ಪ್ರಕಾರ, `--model` ನ ವಾರದ ಅವಧಿಯನ್ನೂ ಎಣಿಸಿ; 100% ತಲುಪಿದ ಖಾತೆಯನ್ನು ಎಂದಿಗೂ ಆಯ್ಕೆಮಾಡುವುದಿಲ್ಲ. ಯಾವುದನ್ನೂ ಆಯ್ಕೆಮಾಡಲಾಗದಿದ್ದಾಗ ಅವು ಕೋಡ್ 75 ನೊಂದಿಗೆ ನಿರ್ಗಮಿಸುತ್ತವೆ, ಅಥವಾ `--wait` ಇದ್ದರೆ ಮೊದಲ ಖಾತೆ ರೀಸೆಟ್ ಆಗುವವರೆಗೆ ಕಾಯುತ್ತವೆ. [ಮುಂದಿನ ಖಾತೆಗೆ ಬದಲಾಯಿಸುವುದು](/features/next) ಮತ್ತು [ಸ್ಕ್ರಿಪ್ಟ್‌ಗಳಿಗಾಗಿ ಖಾತೆ ಆಯ್ಕೆ](/features/pick) ನೋಡಿ.

## ನಿರ್ವಹಣೆ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

ccam ಅನ್ನು ಬದಲಾಯಿಸುವ ಮೊದಲು `ccam update` ಬಿಡುಗಡೆಯ SHA-256 ಚೆಕ್‌ಸಮ್ ಅನ್ನು ಪರಿಶೀಲಿಸುತ್ತದೆ, ಮತ್ತು `--version v0.2.0` ಒಂದು ನಿರ್ದಿಷ್ಟ ಬಿಡುಗಡೆಯನ್ನು ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡುತ್ತದೆ. `.deb` ಪ್ಯಾಕೇಜ್‌ನಿಂದ ಇನ್‌ಸ್ಟಾಲ್ ಮಾಡಿದ ccam ಅನ್ನು ಬದಲಿಗೆ apt ಮೂಲಕ ಅಪ್‌ಡೇಟ್ ಮಾಡಲಾಗುತ್ತದೆ. [ಅಪ್‌ಡೇಟ್](/docs/install#update) ನೋಡಿ.

## ಫೈಲ್‌ಗಳು

| ಪಾತ್ | ಒಳಗೊಂಡಿರುವುದು |
| --- | --- |
| `~/.config/ccam/config.toml` | ಖಾತೆಗಳು ಮತ್ತು ಅವು ಏನನ್ನು ಹಂಚಿಕೊಳ್ಳುತ್ತವೆ |
| `~/.local/share/ccam/profiles/<name>/` | ಪ್ರತಿ ಖಾತೆಯ `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | ಹಂಚಿಕೊಂಡ ಸೆಟ್ಟಿಂಗ್‌ಗಳು, MCP ಸರ್ವರ್‌ಗಳು, ವಿಸ್ತರಣೆಗಳು ಮತ್ತು ಇತಿಹಾಸ |
| `~/.local/state/ccam/` | ಬಳಕೆಯ ಕ್ಯಾಶ್ ಮತ್ತು ಕೊನೆಯ ಬಾರಿ ಬಳಸಿದ ಸಮಯಗಳು |
| `~/.local/bin/<shortcut>` | `ccwork` ನಂತಹ ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು, ccam ಗೆ ಲಿಂಕ್‌ಗಳು |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ಮತ್ತು `XDG_STATE_HOME` ವೇರಿಯೇಬಲ್‌ಗಳು ಈ ಫೋಲ್ಡರ್‌ಗಳ ಸ್ಥಳವನ್ನು ಬದಲಾಯಿಸುತ್ತವೆ. Windows ನಲ್ಲಿ, `~` ಎಂದರೆ `%USERPROFILE%`, ಮತ್ತು ಶಾರ್ಟ್‌ಕಟ್‌ಗಳು `<shortcut>.exe`.
