---
title: Mga background session ng bawat account
description: Inililista ng ccam agents ang mga background agent at bukas na session ng Claude Code ng lahat ng account mo; binubuksan ng ccam attach ang isa ayon sa id.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Nag-iimbak ang Claude Code ng mga background session (`claude --bg`, agent view) bawat account: ang mga session lang ng account na pinagtatakbuhan nito ang ipinapakita ng `claude agents`. Sabay-sabay na tinatanong ng `ccam agents` ang bawat account at inililista ang mga ito nang magkakasama, at binubuksan ng `ccam attach <id>` ang isang background session bilang alinmang account na nagsimula nito.

## Paggamit

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

Ang `State` ay ang estado ng background session (`working`, `blocked`, `done`, `failed`, `stopped`) at kung ano ang ginagawa ng process nito (`busy`, `waiting`, `idle`); `open` ang ipinapakita ng bukas na interactive session.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--all` | Ilista rin ang mga background session na tapos na, pumalya o huminto | off |
| `--json` | I-print ang mga session bilang JSON, bawat isa kasama ang `account` na kinabibilangan nito | off |

Ipinapasa ng `ccam attach <id> [claude args]` sa `claude attach` ang mga argument na kasunod ng id. Sa `--dry-run`, ipi-print nito ang command.

## Paano ito gumagana

1. Para sa bawat account, pinapatakbo ng ccam ang `claude agents --json` (kasama ang `--all` kapag ibinigay) na may `CLAUDE_CONFIG_DIR` ng account, nang sabay-sabay, bawat isa sa loob ng 15 segundo. Binabasa nito ang sariling listahan ng session ng account at hindi nito sinisimulan ang background service ng Claude Code.
2. Nilalagyan ang bawat entry ng pangalan ng account, at inaayos ang listahan ayon sa oras ng pagsisimula, nauuna ang pinakabago. Pinapangalanan sa standard error at hindi isinasama ang account na hindi sumasagot, at saka lalabas ang ccam na may code 1, para hindi mapagkamalan ng script na walang laman na listahan ang isang pagkabigo.
3. Hinahanap ng `ccam attach` ang `jobs/<id>` sa folder ng bawat account, kung saan iniimbak ng Claude Code ang isang background session, at pinapatakbo ang `claude attach <id>` bilang account na iyon, sa parehong paraan ng pagsisimula ng `ccam <name>` sa Claude Code.

## Kaligtasan at mga espesyal na kaso

- Pag-aari ng iisang account ang mga id ng background session: hindi mahahanap ng `claude attach` na pinatakbo bilang ibang account ang mga ito, at iyon ang inaayos ng `ccam attach`.
- Puwedeng may mga session ang dalawang account na pareho ang simula ng id; sa ganitong kaso, pinapangalanan ng `ccam attach` ang dalawa at humihingi ng mas mahabang bahagi ng id.
- Ipinapasa ng ccam ang mga shared na MCP server nito sa Claude Code gamit ang `--mcp-config=<file>` sa iisang argument, na tinatanggap ng `claude agents` at `claude attach`; patuloy na ginagamit ng mga background session ang file na iyon kapag nagre-restart sila.
- Kinukuha ng background service ng isang API key account ang key nito mula sa shell na nagsimula rito. Simulan sa pamamagitan ng ccam (`ccam gw --bg …`) ang unang background session ng ganitong account, para makuha ng service ang key.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Pagpili ng account para sa mga script](/features/pick)
- [Ang dashboard](/features/dashboard)
