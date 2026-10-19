---
title: ಖಾತೆಗಳ ನಿರ್ವಹಣೆ
description: ccam add, ccam login, ccam list, ccam rename, ccam remove ಮತ್ತು ccam import ಮೂಲಕ Claude Code ಖಾತೆಗಳನ್ನು ಸೇರಿಸಿ, ಲಾಗಿನ್ ಮಾಡಿ, ಪಟ್ಟಿ ಮಾಡಿ, ಮರುಹೆಸರಿಸಿ, ತೆಗೆದುಹಾಕಿ ಮತ್ತು ಇಂಪೋರ್ಟ್ ಮಾಡಿ.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

ಪ್ರತಿ ccam ಖಾತೆಯೂ ತನ್ನದೇ ಲಾಗಿನ್ ಇರುವ ಪ್ರತ್ಯೇಕ Claude Code ಕಾನ್ಫಿಗ್ ಫೋಲ್ಡರ್. ಈ ಕಮಾಂಡ್‌ಗಳು ಖಾತೆಗಳನ್ನು ರಚಿಸುತ್ತವೆ, ಸೈನ್ ಇನ್ ಮಾಡುತ್ತವೆ, ಪಟ್ಟಿ ಮಾಡುತ್ತವೆ, ಮರುಹೆಸರಿಸುತ್ತವೆ ಮತ್ತು ತೆಗೆದುಹಾಕುತ್ತವೆ, ಹಾಗೂ ನಿಮ್ಮ ಬಳಿ ಈಗಾಗಲೇ ಇರುವ ಫೋಲ್ಡರ್‌ಗಳನ್ನು ನೋಂದಾಯಿಸುತ್ತವೆ.

## ಬಳಕೆ

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

ಕಮಾಂಡ್ `<name|number>` ತೆಗೆದುಕೊಳ್ಳುವ ಎಲ್ಲೆಡೆ, ನೀವು ಖಾತೆಯ ಹೆಸರನ್ನು ಅಥವಾ ಅದರ ಸ್ಲಾಟ್ ಸಂಖ್ಯೆಯನ್ನು, ಅಂದರೆ `ccam list` ನ `#` ಕಾಲಮ್ ಅನ್ನು, ನೀಡಬಹುದು.

ಆಯ್ಕೆಗಳಿಲ್ಲದೆ ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ `ccam add` ರನ್ ಮಾಡಿದರೆ ವಿಜಾರ್ಡ್ ತೆರೆಯುತ್ತದೆ. ಅದು ಹೆಸರನ್ನು ಕೇಳುತ್ತದೆ, ನಂತರ ಹಂಚಿಕೊಳ್ಳುವ ಐಟಂಗಳ ಚೆಕ್‌ಲಿಸ್ಟ್ ತೋರಿಸುತ್ತದೆ (ಸಂಭಾಷಣೆಯ ಇತಿಹಾಸವನ್ನು ಹೊರತುಪಡಿಸಿ ಎಲ್ಲವನ್ನೂ ಗುರುತಿಸಲಾಗಿರುತ್ತದೆ), ನಂತರ `Create a shortcut command for work? [y/N]` ಎಂದು ಕೇಳುತ್ತದೆ ಮತ್ತು `y` ಎಂದರೆ `Shortcut name [ccwork]:` ಎಂದು ಕೇಳುತ್ತದೆ. `--share`, `--shortcut` ಅಥವಾ `--no-login` ಇವುಗಳಲ್ಲಿ ಯಾವುದಾದರೂ ಒಂದು ಇದ್ದರೆ, ಅಥವಾ ಟರ್ಮಿನಲ್ ಬದಲು ಪೈಪ್ ಇದ್ದರೆ, ವಿಜಾರ್ಡ್ ಬರುವುದಿಲ್ಲ; ಆಗ ಹೆಸರು ಕಡ್ಡಾಯ, ಮತ್ತು `--shortcut` ಒಂದು ಹೆಸರನ್ನು ನೀಡದಿದ್ದರೆ ಯಾವುದೇ ಶಾರ್ಟ್‌ಕಟ್ ರಚನೆಯಾಗುವುದಿಲ್ಲ.

`ccam list` ಸ್ಲಾಟ್, ಹೆಸರು, ಶಾರ್ಟ್‌ಕಟ್, ಇಮೇಲ್, ಹಂಚಿಕೊಂಡ ಐಟಂಗಳು ಮತ್ತು ಕೊನೆಯ ಬಳಕೆಯನ್ನು ತೋರಿಸುತ್ತದೆ, ಡೀಫಾಲ್ಟ್ ಖಾತೆಯ ಮುಂದೆ `*` ಇರುತ್ತದೆ:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` ಜೊತೆ, ಪ್ರತಿ ಖಾತೆಯೂ `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (ಖಾತೆಯನ್ನು ಎಂದೂ ಆರಂಭಿಸದಿದ್ದರೆ ಇದು ಇರುವುದಿಲ್ಲ) ಮತ್ತು `share` ಇರುವ ಒಂದು ಆಬ್ಜೆಕ್ಟ್ ಆಗಿರುತ್ತದೆ. `share` ಪ್ರತಿ ಐಟಂ (`mcp`, `settings`, `memory`, `extensions`, `history`) ಅನ್ನು ಅದರ ಗುಂಪಿನ ಹೆಸರಿಗೆ ಅಥವಾ `"local"` ಗೆ ಹೊಂದಿಸುತ್ತದೆ.

## ಆಯ್ಕೆಗಳು

| ಆಯ್ಕೆ | ಏನು ಮಾಡುತ್ತದೆ | ಡೀಫಾಲ್ಟ್ |
| --- | --- | --- |
| `ccam add --no-login` | ಲಾಗಿನ್ ಮಾಡದೆ ಖಾತೆಯನ್ನು ರಚಿಸುತ್ತದೆ; ನಂತರ `ccam login <name>` ಮೂಲಕ ಲಾಗಿನ್ ಮಾಡಿ | ಆಫ್ |
| `ccam add --share <item=group>` | ಒಂದು ಐಟಂ ಅನ್ನು ಗುಂಪಿನೊಂದಿಗೆ ಹಂಚಿಕೊಳ್ಳುತ್ತದೆ, ಅಥವಾ ಖಾಸಗಿಯಾಗಿಡಲು `item=local`. ಹಲವು ಬಾರಿ ನೀಡಬಹುದು | `default` ನಲ್ಲಿ `mcp`, `settings`, `memory`, `extensions`; `history` ಸ್ಥಳೀಯ |
| `ccam add --shortcut <command>` | `~/.local/bin` ನಲ್ಲಿ ಈ ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್ ಅನ್ನೂ ರಚಿಸುತ್ತದೆ | ಇಲ್ಲ |
| `ccam list --json` | ಖಾತೆಗಳನ್ನು JSON ಆಗಿ ಮುದ್ರಿಸುತ್ತದೆ | ಆಫ್ |
| `ccam remove -y`, `--yes` | ಕೇಳದೆ ತೆಗೆದುಹಾಕುತ್ತದೆ | ಆಫ್ |
| `ccam import --name <name>` | ನೀವು ನೀಡುವ ಫೋಲ್ಡರ್‌ಗೆ ಖಾತೆಯ ಹೆಸರು | ಫೋಲ್ಡರ್ ಹೆಸರಿನಿಂದ ಪಡೆಯಲಾಗುತ್ತದೆ |
| `ccam import -y`, `--yes` | ಸಿಕ್ಕ ಪ್ರತಿ ಫೋಲ್ಡರ್ ಅನ್ನು ಕೇಳದೆ ಇಂಪೋರ್ಟ್ ಮಾಡುತ್ತದೆ | ಆಫ್ |
| `--dry-run` | ಏನನ್ನೂ ಬದಲಾಯಿಸದ `ccam list` ಮತ್ತು `ccam completion` ಇದನ್ನು ಪಾಲಿಸುತ್ತವೆ. `add`, `login`, `logout`, `use`, `rename`, `remove` ಮತ್ತು `import` ಇದನ್ನು `` --dry-run is not supported by `ccam add`; nothing was changed `` (ಕಮಾಂಡ್ ಹೆಸರಿನೊಂದಿಗೆ) ಎಂದು ನಿರಾಕರಿಸುತ್ತವೆ | ಆಫ್ |

`ccam login`, `ccam logout`, `ccam use` ಮತ್ತು `ccam rename` ಗೆ ತಮ್ಮದೇ ಆಯ್ಕೆಗಳಿಲ್ಲ. `ccam completion` `bash` ಅಥವಾ `zsh` ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ; ಬೇರೆ ಯಾವುದೇ ಮೌಲ್ಯ `only bash and zsh are supported` ಎಂದು ವಿಫಲವಾಗುತ್ತದೆ.

ಖಾತೆಯ ಹೆಸರುಗಳು `a-z`, `0-9`, `-` ಮತ್ತು `_` ಅಕ್ಷರಗಳಿಂದ ಕೂಡಿದ 1 ರಿಂದ 32 ಅಕ್ಷರಗಳಾಗಿರುತ್ತವೆ ಮತ್ತು ಅಕ್ಷರ ಅಥವಾ ಅಂಕಿಯಿಂದ ಆರಂಭವಾಗುತ್ತವೆ. ಹೆಸರು ಸಂಪೂರ್ಣವಾಗಿ ಅಂಕಿಗಳಾಗಿರಬಾರದು, ಏಕೆಂದರೆ `ccam 2` ಎಂದರೆ ಸ್ಲಾಟ್ 2. ಹೆಸರು ccam ಕಮಾಂಡ್ ಕೂಡ ಆಗಿರಬಾರದು (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), ಹೀಗಾಗಿ `ccam <name>` ಎಂದಿಗೂ ಅಸ್ಪಷ್ಟವಾಗುವುದಿಲ್ಲ. ಈ ಪಟ್ಟಿ `internal/config` ನಲ್ಲಿರುವ `config.Reserved`.

## ಇದು ಹೇಗೆ ಕೆಲಸ ಮಾಡುತ್ತದೆ

`ccam add work --shortcut ccwork` `addProfile` (`internal/cli/profiles.go`) ನಲ್ಲಿ ಈ ಹಂತಗಳನ್ನು ನಡೆಸುತ್ತದೆ:

<!-- motion:accounts-add -->

1. ಏನನ್ನಾದರೂ ರಚಿಸುವ ಮೊದಲು `config.ValidateName` ಹೆಸರನ್ನು ಮತ್ತು `shim.Check` ಶಾರ್ಟ್‌ಕಟ್ ಅನ್ನು ಪರಿಶೀಲಿಸುತ್ತವೆ.
2. `profile.Add` `0700` ಮೋಡ್‌ನೊಂದಿಗೆ `~/.local/share/ccam/profiles/work/` ಅನ್ನು ರಚಿಸುತ್ತದೆ ಮತ್ತು ಖಾತೆಗೆ ಅತ್ಯಂತ ಕಡಿಮೆ ಖಾಲಿ ಸ್ಲಾಟ್ ಸಂಖ್ಯೆಯನ್ನು (`Config.NextSlot`) ನೀಡುತ್ತದೆ, ಹೀಗಾಗಿ ತೆಗೆದುಹಾಕಿದ ಖಾತೆಯ ಸಂಖ್ಯೆ ಮತ್ತೆ ಬಳಕೆಯಾಗುತ್ತದೆ. ಮೊದಲ ಖಾತೆ ಡೀಫಾಲ್ಟ್ ಆಗುತ್ತದೆ.
3. `applyShare` ಹೊಸ, ಖಾಲಿ ಫೋಲ್ಡರ್ ಅನ್ನು `share.Build` ಮತ್ತು ನಂತರ `Apply` ಮೂಲಕ ಅದರ ಹಂಚಿಕೆ ಗುಂಪುಗಳಿಗೆ ಲಿಂಕ್ ಮಾಡುತ್ತದೆ, ಉದಾಹರಣೆಗೆ `settings.json` ಅನ್ನು `groups/default/settings.json` ಗೆ. ಇದು ವಿಫಲವಾದರೆ ಹೊಸ ಫೋಲ್ಡರ್ ತೆಗೆದುಹಾಕಲಾಗುತ್ತದೆ, ಏಕೆಂದರೆ ಉಳಿದ ಫೋಲ್ಡರ್ ಪ್ರತಿ ಮರುಪ್ರಯತ್ನವನ್ನೂ `folder already exists` ಎಂದು ತಡೆಯುತ್ತದೆ.
4. ಖಾತೆಯನ್ನು `fsutil.AtomicWrite` ಮೂಲಕ `~/.config/ccam/config.toml` ಗೆ ಉಳಿಸಲಾಗುತ್ತದೆ: ಒಂದು ತಾತ್ಕಾಲಿಕ ಫೈಲ್, `fsync`, ನಂತರ `rename`. ಹೀಗಾಗಿ ಕ್ರ್ಯಾಶ್ ಆದರೂ ಅರ್ಧ ಕಾನ್ಫಿಗ್ ಉಳಿಯುವುದಿಲ್ಲ.
5. `runLogin` `CLAUDE_CONFIG_DIR` ಅನ್ನು ಖಾತೆಯ ಫೋಲ್ಡರ್‌ಗೆ ಹೊಂದಿಸಿ `claude auth login` ಅನ್ನು ಆರಂಭಿಸುತ್ತದೆ. ಲಾಗಿನ್ ವಿಫಲವಾದರೆ ಖಾತೆ ಉಳಿಯುತ್ತದೆ, ಮತ್ತು ccam ದೋಷವನ್ನು ಮುದ್ರಿಸಿ ನಂತರ `ccam login work` ರನ್ ಮಾಡಲು ಹೇಳುತ್ತದೆ.
6. `profile.ReadEmail` `profiles/work/.claude.json` ನಿಂದ `oauthAccount.emailAddress` ಅನ್ನು ತೆಗೆದುಕೊಳ್ಳುತ್ತದೆ. ಆ ಫೈಲ್ ಹಲವು ಮೆಗಾಬೈಟ್ ಇರಬಹುದು, ಹೀಗಾಗಿ `TopLevelValue` ಇಡೀ ಫೈಲ್ ಅನ್ನು ಪಾರ್ಸ್ ಮಾಡುವ ಬದಲು ಆ ಒಂದು ಕೀಗಾಗಿ ಮಾತ್ರ ಹುಡುಕುತ್ತದೆ.
7. `shim.Install` `~/.local/bin/ccwork` ಅನ್ನು ccam ಬೈನರಿಗೆ ಲಿಂಕ್ ಮಾಡುತ್ತದೆ; ಅದು ಆರಂಭವಾದಾಗ `config.toml` ನಲ್ಲಿ `argv[0]` ಅನ್ನು ಹುಡುಕುತ್ತದೆ.

ಲಾಗಿನ್‌ಗೆ ಕೆಲವು ನಿಮಿಷಗಳು ಬೇಕಾಗಬಹುದು, ಮತ್ತು ಅಷ್ಟರಲ್ಲಿ ಇನ್ನೊಂದು ccam ಕಾನ್ಫಿಗ್ ಅನ್ನು ಉಳಿಸಬಹುದು. ಆದ್ದರಿಂದ ಕಾಯುವ ಕಮಾಂಡ್‌ಗಳು `env.update` ಮೂಲಕ ಉಳಿಸುತ್ತವೆ: ಅದು `config.toml` ಅನ್ನು ಮತ್ತೆ ಓದುತ್ತದೆ, ಆ ಒಂದು ಬದಲಾವಣೆಯನ್ನು (ಇಲ್ಲಿ ಇಮೇಲ್) ಅನ್ವಯಿಸುತ್ತದೆ ಮತ್ತು ಉಳಿಸುತ್ತದೆ. ಇದೇ ಕಾರಣಕ್ಕೆ `ccam remove` ತನ್ನ ಪ್ರಶ್ನೆಯ ನಂತರ ಮತ್ತೆ ಓದುತ್ತದೆ.

`ccam login` ಈಗಿರುವ ಖಾತೆಗೆ ಹಂತ 5 ಮತ್ತು 6 ಅನ್ನು ನಡೆಸುತ್ತದೆ. `ccam logout` ಖಾತೆಯಲ್ಲಿ `claude auth logout` ರನ್ ಮಾಡುತ್ತದೆ ಮತ್ತು ಉಳಿಸಿದ ಇಮೇಲ್ ಅನ್ನು ಅಳಿಸುತ್ತದೆ. `ccam use` `config.toml` ನಲ್ಲಿ `default` ಅನ್ನು ಮಾತ್ರ ಬದಲಾಯಿಸುತ್ತದೆ.

### ರಿಮೋಟ್ ಯಂತ್ರದಲ್ಲಿ ಸೈನ್ ಇನ್ ಮಾಡುವುದು

SSH ಮೂಲಕ, `claude auth login` ನ ಬ್ರೌಸರ್ ಭಾಗ ಮುರಿಯಬಹುದು. Claude Code ಸರ್ವರ್‌ನಲ್ಲಿ `http://localhost:<port>/callback` ನಲ್ಲಿರುವ ಕಾಲ್‌ಬ್ಯಾಕ್ ಸರ್ವರ್‌ನಲ್ಲಿ ಸೈನ್-ಇನ್‌ಗಾಗಿ ಕಾಯುತ್ತದೆ, ಆದರೆ ಸೈನ್-ಇನ್ ಮುಗಿಸುವ ಬ್ರೌಸರ್ ನಿಮ್ಮ ಕಂಪ್ಯೂಟರ್‌ನಲ್ಲಿ ಚಾಲನೆಯಾಗುತ್ತದೆ. ಅಲ್ಲಿ ಆ ವಿಳಾಸ ಎಲ್ಲಿಗೂ ತಲುಪುವುದಿಲ್ಲ, ಮತ್ತು ಬ್ರೌಸರ್ `ERR_CONNECTION_REFUSED` ನಂತಹ ದೋಷ ತೋರಿಸುತ್ತದೆ. ಆದರೂ ಆ ವಿಳಾಸದಲ್ಲಿ ಒಂದು ಬಾರಿಯ ಕೋಡ್ ಇರುತ್ತದೆ, ಹೀಗಾಗಿ ccam ಅದನ್ನು ತಲುಪಿಸಬಲ್ಲದು:

<!-- motion:accounts-relay -->

`runLogin` ನಿಮ್ಮ ಟರ್ಮಿನಲ್ ಅನ್ನು Claude Code ಗೆ ಒಪ್ಪಿಸುವುದಿಲ್ಲ. ಅದು Claude Code ನ ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್ ಅನ್ನು ಒಂದು ಪೈಪ್‌ಗೆ ಜೋಡಿಸುತ್ತದೆ, ಮತ್ತು `relayLogin` ನೀವು ಟೈಪ್ ಮಾಡುವ ಪ್ರತಿ ಸಾಲನ್ನು ಅದಕ್ಕೆ ನಕಲಿಸುತ್ತದೆ. ಹೀಗಾಗಿ Claude Code ನ `Paste code here` ಪ್ರಾಂಪ್ಟ್‌ನಲ್ಲಿ ಪೇಸ್ಟ್ ಮಾಡಿದ ಕೋಡ್ ಕೂಡ ತಲುಪುತ್ತದೆ. ಒಂದು ಸಾಲು ಕಾಲ್‌ಬ್ಯಾಕ್ ವಿಳಾಸವಾಗಿದ್ದರೆ, ಹೋಸ್ಟ್ `localhost`, `127.0.0.1` ಅಥವಾ `[::1]` ಆಗಿದ್ದರೆ, ಪೋರ್ಟ್ ನೀಡಿದ್ದರೆ, ಪಾತ್ `/callback` ಆಗಿದ್ದರೆ ಮತ್ತು `code` ಹಾಗೂ `state` ಎರಡೂ ಇದ್ದರೆ ಮಾತ್ರ `callbackURL` ಅದನ್ನು ಸ್ವೀಕರಿಸುತ್ತದೆ; `http://` ಬಿಟ್ಟುಬಿಡಬಹುದು. ನಂತರ `deliverCallback` ಬ್ರೌಸರ್ ಮಾಡುತ್ತಿದ್ದಂತೆಯೇ, 10 ಸೆಕೆಂಡ್ ಮಿತಿಯೊಂದಿಗೆ ಮತ್ತು ರೀಡೈರೆಕ್ಟ್‌ಗಳನ್ನು ಅನುಸರಿಸದೆ, ಸರ್ವರ್‌ನಲ್ಲಿ ಆ ವಿಳಾಸವನ್ನು ವಿನಂತಿಸುತ್ತದೆ. ವಿನಂತಿ ಕಾಯುತ್ತಿರುವ `claude auth login` ನ ಕಾಲ್‌ಬ್ಯಾಕ್ ಸರ್ವರ್ ಅನ್ನು ತಲುಪುತ್ತದೆ (ಚಿತ್ರದಲ್ಲಿನ ಡ್ಯಾಶ್ ಗೆರೆ). Claude Code ಕೋಡ್ ಪಡೆದು ಲಾಗಿನ್ ಅನ್ನು ಉಳಿಸುತ್ತದೆ, Linux ನಲ್ಲಿ ಖಾತೆಯ ಫೋಲ್ಡರ್‌ನಲ್ಲಿರುವ `.credentials.json` ನಲ್ಲಿ. ಪೇಸ್ಟ್ ಮಾಡಿದ ಸಾಲು ccam ಅನ್ನು ಎಂದಿಗೂ ಬೇರೆ ಹೋಸ್ಟ್ ಅನ್ನು ಸಂಪರ್ಕಿಸುವಂತೆ ಮಾಡಲಾರದು.

ಲಾಗಿನ್ ಆರಂಭವಾಗುವ ಮೊದಲು ccam ಇದನ್ನು ಮುದ್ರಿಸುತ್ತದೆ:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

ಯಶಸ್ವಿಯಾದಾಗ ಅದು `→ Passed the sign-in to Claude Code on this machine.` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ.

ಲಾಗಿನ್ ಟರ್ಮಿನಲ್‌ನಲ್ಲಿ ಇನ್‌ಪುಟ್‌ಗಾಗಿ ಕಾಯುತ್ತಿರುವಾಗ, ccam ಟರ್ಮಿನಲ್‌ನ echo ಮತ್ತು line mode ಅನ್ನು ಆಫ್ ಮಾಡುತ್ತದೆ ಮತ್ತು ನೀವು ಟೈಪ್ ಮಾಡುವ ಪ್ರತಿ ಅಕ್ಷರಕ್ಕೂ ಒಂದು `*` ತೋರಿಸುತ್ತದೆ, ಹೀಗಾಗಿ ಪೇಸ್ಟ್ ಮಾಡಿದ ಕೋಡ್ ಅಥವಾ callback ವಿಳಾಸ ಪರದೆಯ ಮೇಲೆ ಕಾಣುವುದಿಲ್ಲ. `readLine` ಸಾಲಿನ ಸಂಪಾದನೆಯನ್ನು ತಾನೇ ಮಾಡುತ್ತದೆ: `Enter` ಸಾಲನ್ನು ಕಳುಹಿಸುತ್ತದೆ, `Backspace` ಒಂದು ಅಕ್ಷರವನ್ನು ಅಳಿಸುತ್ತದೆ, ಮತ್ತು ಖಾಲಿ ಸಾಲಿನಲ್ಲಿ `Ctrl+D` ಇನ್‌ಪುಟ್ ಅನ್ನು ಮುಗಿಸುತ್ತದೆ. ಲಾಗಿನ್ ಮುಗಿದಾಗ ಟರ್ಮಿನಲ್ ಹಿಂದಿನ ಸ್ಥಿತಿಗೆ ಮರಳುತ್ತದೆ, ಮತ್ತು `Ctrl+C`, `SIGTERM` ಅಥವಾ `SIGHUP` ಬಂದಾಗಲೂ ಮರಳುತ್ತದೆ; ಇವು ಈಗಲೂ ccam ಅನ್ನು ಕೊನೆಗೊಳಿಸುತ್ತವೆ. ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಇನ್‌ಪುಟ್ ಮತ್ತು ಸ್ಟ್ಯಾಂಡರ್ಡ್ ಎರರ್‌ನಲ್ಲಿ ಟರ್ಮಿನಲ್ ಇಲ್ಲದಿದ್ದರೆ, ಸಾಲುಗಳು ಬಂದಂತೆಯೇ ನಕಲಾಗುತ್ತವೆ.

### ಮರುಹೆಸರಿಸುವುದು, ತೆಗೆದುಹಾಕುವುದು ಮತ್ತು ಇಂಪೋರ್ಟ್

`ccam rename` ಮತ್ತು `ccam remove` ಖಾತೆಯ ಫೋಲ್ಡರ್ ಮೇಲೆ ಮತ್ತು ccam ಖಾತೆಯ ಹೆಸರಿನಲ್ಲಿ ಬೇರೆಡೆ ಇಟ್ಟಿರುವುದರ ಮೇಲೆ ಕೆಲಸ ಮಾಡುತ್ತವೆ. `nameState` ಆ ಪಾತ್‌ಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` ಮತ್ತು ಬಳಕೆಯ ಕ್ಯಾಶ್‌ಗಳಾದ `~/.local/state/ccam/usage/<name>.json` ಹಾಗೂ `<name>.statusline.json`. `ccam remove work` ಅವುಗಳನ್ನು ಈ ಕ್ರಮದಲ್ಲಿ ನಿಭಾಯಿಸುತ್ತದೆ:

<!-- motion:accounts-remove -->

ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿದ್ದರೆ `ccam remove` ಮೊದಲು ನಿರಾಕರಿಸುತ್ತದೆ, ಇಲ್ಲದಿದ್ದರೆ ಕೇಳುತ್ತದೆ. ಖಾತೆ ವಿಸ್ತರಣೆಗಳನ್ನು ಹಂಚಿಕೊಂಡಿದ್ದರೆ, `RewritePluginPaths` ಗುಂಪಿನ `installed_plugins.json` ಮತ್ತು `known_marketplaces.json` ಅನ್ನು ಮರುಬರೆಯುತ್ತದೆ, ಏಕೆಂದರೆ Claude Code ಈ ಖಾತೆಯ ಫೋಲ್ಡರ್ ಮೂಲಕ ಹಾದುಹೋಗಬಹುದಾದ ಸಂಪೂರ್ಣ ಪ್ಲಗಿನ್ ಪಾತ್‌ಗಳನ್ನು ಸಂಗ್ರಹಿಸುತ್ತದೆ. `profile.Trash` ಫೋಲ್ಡರ್ ಅನ್ನು `~/.local/share/ccam/trash/<name>-<time>/` ಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ, ಖಾತೆ `config.toml` ನಿಂದ ಹೊರಹೋಗುತ್ತದೆ, ನಂತರ ಅದೇ ಹೆಸರಿನ ಖಾತೆ ಅವುಗಳನ್ನು ಪಡೆಯದಂತೆ `nameState` ಪಾತ್‌ಗಳನ್ನು ಅಳಿಸಲಾಗುತ್ತದೆ, ಮತ್ತು ಶಾರ್ಟ್‌ಕಟ್ ಲಿಂಕ್ ತೆಗೆದುಹಾಕಲಾಗುತ್ತದೆ. ಮರಳಿ ಪಡೆಯಲು, ಫೋಲ್ಡರ್ ಅನ್ನು ಟ್ರ್ಯಾಶ್‌ನಿಂದ ಹಿಂದಕ್ಕೆ ಸ್ಥಳಾಂತರಿಸಿ.

`ccam rename` ಫೋಲ್ಡರ್ ಅನ್ನು ಮರುಹೆಸರಿಸುತ್ತದೆ, ಅದರ `plugins/` ನಲ್ಲಿನ ಪ್ಲಗಿನ್ ಪಾತ್‌ಗಳನ್ನು ಹಳೆಯ ಫೋಲ್ಡರ್‌ನಿಂದ ಹೊಸದಕ್ಕೆ ಮರುಬರೆಯುತ್ತದೆ, ಡೀಫಾಲ್ಟ್ ಅನ್ನು ಜೊತೆಗೆ ಕೊಂಡೊಯ್ಯುತ್ತದೆ, ಮತ್ತು `nameState` ಪಾತ್‌ಗಳನ್ನು ಹೊಸ ಹೆಸರಿಗೆ ಸ್ಥಳಾಂತರಿಸುತ್ತದೆ. ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್ ತನ್ನ ಹೆಸರನ್ನು ಉಳಿಸಿಕೊಳ್ಳುತ್ತದೆ. ಪ್ಲಗಿನ್ ಮರುಬರಹ ವಿಫಲವಾದರೆ, ಕಾನ್ಫಿಗ್ ಉಳಿಸಲಾಗದ ಕಾರಣ ಫೋಲ್ಡರ್‌ಗೆ ಹಳೆಯ ಹೆಸರನ್ನೇ ಮರಳಿ ನೀಡಲಾಗುತ್ತದೆ.

`ccam import` ಒಂದು ಫೋಲ್ಡರ್ ಅನ್ನು ಅದು ಇರುವಲ್ಲೇ, `config.toml` ನಲ್ಲಿ `dir` ಆಗಿ ನೋಂದಾಯಿಸುತ್ತದೆ. ಅದು ಅದನ್ನು ಎಂದಿಗೂ ಸ್ಥಳಾಂತರಿಸುವುದಿಲ್ಲ ಅಥವಾ ನಕಲಿಸುವುದಿಲ್ಲ: ನಕಲಿಸಿದ ರಿಫ್ರೆಶ್ ಟೋಕನ್ ಪ್ರತಿಗಳು ಒಂದನ್ನೊಂದು ಲಾಗ್ ಔಟ್ ಮಾಡುವಂತೆ ಮಾಡುತ್ತದೆ. ಫೋಲ್ಡರ್ ಆರ್ಗ್ಯುಮೆಂಟ್ ಇಲ್ಲದಿದ್ದರೆ, ಅದು ನೋಂದಾಯಿಸದ `~/.claude` ಮತ್ತು `~/.claude-*` ಫೋಲ್ಡರ್‌ಗಳಿಗಾಗಿ ಹುಡುಕಿ ಪ್ರತಿಯೊಂದರ ಬಗ್ಗೆ ಕೇಳುತ್ತದೆ; ಹೆಸರುಗಳು ಫೋಲ್ಡರ್‌ನಿಂದ ಬರುತ್ತವೆ (`~/.claude` ಎಂದರೆ `main`, `~/.claude-work` ಎಂದರೆ `work`, `~/.claude-2` ಎಂದರೆ `acct2`), ಮತ್ತು ಹೆಸರನ್ನು ಪಡೆಯಲಾಗದ ಫೋಲ್ಡರ್ ಅನ್ನು ಬಳಸಬೇಕಾದ `--name` ಕಮಾಂಡ್ ಜೊತೆ ಪಟ್ಟಿ ಮಾಡಲಾಗುತ್ತದೆ. `CLAUDE_CONFIG_DIR` ಇಲ್ಲದಿದ್ದಾಗ Claude Code `.claude.json` ಅನ್ನು ನಿಮ್ಮ ಹೋಮ್ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ `~/.claude` ಪಕ್ಕದಲ್ಲಿ ಇಡುತ್ತದೆ. ಹೀಗಾಗಿ ನೀವು `~/.claude` ಅನ್ನು ಇಂಪೋರ್ಟ್ ಮಾಡಿದಾಗ ಅದಕ್ಕೆ ತನ್ನದೇ `.claude.json` ಇಲ್ಲದಿದ್ದರೆ, ccam `~/.claude/.claude.json` ನಿಂದ `~/.claude.json` ಗೆ ಲಿಂಕ್ ರಚಿಸುತ್ತದೆ. ಆಗ ಸಾಮಾನ್ಯ `claude` ಮತ್ತು `ccam main` ಒಂದೇ ಖಾತೆಯನ್ನು ನೋಡುತ್ತವೆ. ನೀವು `ccam share` ರನ್ ಮಾಡುವವರೆಗೆ ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಖಾತೆಗಳು ಏನನ್ನೂ ಹಂಚಿಕೊಳ್ಳುವುದಿಲ್ಲ, ಮತ್ತು ಯಾವುದೇ ಶಾರ್ಟ್‌ಕಟ್ ರಚನೆಯಾಗುವುದಿಲ್ಲ.

`ccam completion bash|zsh` ಖಾತೆಯ ಹೆಸರುಗಳನ್ನೂ ಪೂರ್ಣಗೊಳಿಸುವ ಕಂಪ್ಲೀಷನ್ ಸ್ಕ್ರಿಪ್ಟ್ ಅನ್ನು ಮುದ್ರಿಸುತ್ತದೆ.

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿರುವಾಗ `ccam rename` ಮತ್ತು `ccam remove` `work is running; quit its Claude Code session first` ಎಂದು ನಿರಾಕರಿಸುತ್ತವೆ (`refuseRunning`), ಏಕೆಂದರೆ ಸೆಷನ್ ಫೋಲ್ಡರ್‌ಗೆ ಬರೆಯುತ್ತಲೇ ಇರುತ್ತದೆ.
- ಇಂಪೋರ್ಟ್ ಮಾಡಿದ ಖಾತೆಯ ಮೇಲೆ `ccam remove` ಅದರ ನೋಂದಣಿಯನ್ನು ಮಾತ್ರ ರದ್ದುಗೊಳಿಸುತ್ತದೆ: `✓ main unregistered (folder ... left in place)`. ನಿಮ್ಮ `~/.claude` ಟ್ರ್ಯಾಶ್‌ಗೆ ಹೋಗುವುದಿಲ್ಲ.
- `profiles/` ನಲ್ಲಿ ಈಗಾಗಲೇ ಫೋಲ್ಡರ್ ಇರುವ ಹೆಸರನ್ನು `ccam add` ನಿರಾಕರಿಸುತ್ತದೆ ಮತ್ತು `ccam import` ಅನ್ನು ಸೂಚಿಸುತ್ತದೆ.
- `ccam import` ಈಗಾಗಲೇ ನೋಂದಾಯಿಸಿದ ಫೋಲ್ಡರ್, ಫೋಲ್ಡರ್ ಅಲ್ಲದ ಪಾತ್, ಮತ್ತು ಈಗಾಗಲೇ ಬಳಕೆಯಲ್ಲಿರುವ ಹೆಸರನ್ನು ನಿರಾಕರಿಸುತ್ತದೆ.
- `ccam logout` ಮತ್ತು `ccam login` ಖಾತೆ ಚಾಲನೆಯಲ್ಲಿದೆಯೇ ಎಂದು ಪರಿಶೀಲಿಸುವುದಿಲ್ಲ.
- ರಿಲೇ Claude Code ಅನ್ನು ತಲುಪಲಾಗದಿದ್ದರೆ, ಉದಾಹರಣೆಗೆ ಲಾಗಿನ್ ಮುಗಿದ ನಂತರ, ಅದು ಕಾರಣದೊಂದಿಗೆ `ccam: could not pass the sign-in to Claude Code:` ಎಂದು ಮುದ್ರಿಸುತ್ತದೆ. HTTP ದೋಷ ಸಾಮಾನ್ಯವಾಗಿ ಕೋಡ್‌ನ ಅವಧಿ ಮುಗಿದಿದೆ ಎಂದರ್ಥ; ಲಾಗಿನ್ ಅನ್ನು ಮತ್ತೆ ಆರಂಭಿಸಿ.

## ಸಂಬಂಧಿತ

- [Claude Code ಆರಂಭಿಸುವುದು](/features/launch)
- [ಶಾರ್ಟ್‌ಕಟ್ ಕಮಾಂಡ್‌ಗಳು](/features/shortcuts)
- [MCP ಸರ್ವರ್‌ಗಳು, ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಮತ್ತು ಇತಿಹಾಸವನ್ನು ಹಂಚಿಕೊಳ್ಳಿ](/docs/sharing)
- [ಹಲವು Claude Code ಖಾತೆಗಳು](/guides/multiple-claude-code-accounts)
