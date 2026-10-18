---
title: Marejeo ya amri
description: Kila amri na chaguo la ccam, kuanzia ccam add na ccs <name> hadi ccam usage, ccam next na ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` na `ccs` ni programu moja; tumia unayopendelea. Weka `--dry-run` kabla ya amri yoyote ili kuona itakachofanya bila kukifanya, kwa mfano `ccs --dry-run work`.

<!-- commands -->

## Kuanzisha Claude Code

`ccs <name|number> [arguments]` huanzisha Claude Code kama akaunti hiyo. ccam haisomi hoja zinazokuja baada ya jina; huzipitisha kwa `claude` kwa mpangilio ule ule, na standard input, exit code na signals hufanya kazi sawasawa kama kwa `claude` yenyewe. `--` baada ya jina ni hiari.

`ccs` ikiwa na hoja lakini bila jina la akaunti, kama `ccs -p "summarise"`, hutumia akaunti chaguo-msingi. Iweke kwa `ccam use <name>`.

Njia fupi `claude1`, `claude2` na kadhalika ni viungo vya ccam katika `~/.local/bin`. `ccam shims install --names` pia huunda `claude-<name>` kwa kila akaunti.

## Kuongeza na kuondoa akaunti

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` bila folda hutafuta `~/.claude` na `~/.claude-*` na huuliza kuhusu kila moja. Folda zilizoingizwa husajiliwa pale zilipo; ccam kamwe haihamishi wala kunakili taarifa za kuingia.

## Kushiriki

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Kila kipengele hupokea `on` (kikundi cha `default`), `off`, au jina la kikundi. Tazama [kushiriki](/docs/sharing).

## Seva za MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Matumizi na kubadili

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Matengenezo

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Faili

| Njia | Yaliyomo |
| --- | --- |
| `~/.config/ccam/config.toml` | Akaunti na zinachoshiriki |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ya kila akaunti |
| `~/.local/share/ccam/groups/<group>/` | Mipangilio, seva za MCP, viendelezi na historia vinavyoshirikiwa |
| `~/.local/state/ccam/` | Cache ya matumizi na nyakati za matumizi ya mwisho |
| `~/.local/bin/claude1`, `ccs` | Viungo vya ccam |

Vigezo `XDG_CONFIG_HOME`, `XDG_DATA_HOME` na `XDG_STATE_HOME` huhamisha folda hizi.
