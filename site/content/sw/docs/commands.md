---
title: Marejeo ya amri
description: Kila amri na chaguo la ccam, kuanzia ccam add na ccam <name> hadi ccam usage, ccam next na ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

Weka `--dry-run` kabla ya amri yoyote ili kuona itakachofanya bila kukifanya, kwa mfano `ccam --dry-run work`.

<!-- commands -->

## Kuanzisha Claude Code

`ccam <name|number> [arguments]` huanzisha Claude Code kama akaunti hiyo. ccam haisomi hoja zinazokuja baada ya jina; huzipitisha kwa `claude` kwa mpangilio ule ule, na standard input, exit code na signals hufanya kazi sawasawa kama kwa `claude` yenyewe. `--` baada ya jina ni hiari.

`ccam` ikiwa na hoja lakini bila jina la akaunti, kama `ccam -p "summarise"`, hutumia akaunti chaguo-msingi. Iweke kwa `ccam use <name>`.

Njia fupi kama `ccwork` huanzisha akaunti yake kwa njia ile ile: `ccwork --resume` hufanya kile ambacho `ccam work --resume` hufanya. Tazama [njia fupi](/docs/commands#shortcuts).

## Kuongeza na kuondoa akaunti

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` bila folda hutafuta `~/.claude` na `~/.claude-*` na huuliza kuhusu kila moja. Folda zilizoingizwa husajiliwa pale zilipo; ccam kamwe haihamishi wala kunakili taarifa za kuingia.

## Njia fupi

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Njia fupi ni amri unayoipa jina mwenyewe inayoanzisha Claude Code kama akaunti moja. Kama `ccam work`, hupitisha kila hoja kwa Claude Code bila kuibadilisha, kwa hiyo `ccwork --resume` na `ccwork -p "hi"` hufanya kazi kama unavyotarajia. Ni kiungo cha ccam katika `~/.local/bin`.

Njia fupi ni za hiari. Mchawi wa `ccam add` huuliza kama uunde moja, `ccam add <name> --shortcut <command>` huunda moja bila kuuliza, na `ccam list` huonyesha njia fupi ya kila akaunti, au `-` ikiwa haina.

Majina hutumia herufi a-z, 0-9, `.`, `-` na `_`. Jina haliwezi kuwa amri ya ccam, njia fupi ya akaunti nyingine, au jina la programu nyingine iliyo kwenye `PATH` yako.

`ccam rename` huhifadhi jina la njia fupi, na `ccam remove` hufuta njia fupi pamoja na akaunti. `ccam shims install` huunda upya viungo vya akaunti zote, kwa mfano baada ya kuhamisha ccam.

## Kushiriki

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Kila kipengele hupokea `on` (kikundi cha `default`), `off`, au jina la kikundi. Tazama [kushiriki](/docs/sharing).

## Kunakili usanidi uliopo

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` hunakili `settings.json`, `CLAUDE.md`, viendelezi na seva za MCP za mtumiaji kutoka folda ya Claude Code kwenda kwenye kikundi, ili kila akaunti inayoshiriki vipengele hivyo ivipate. Folda husomwa tu. Chaguo `--settings`, `--memory`, `--extensions`, `--mcp` na `--history` hunakili tu vipengele unavyovitaja; bila hivyo, kila kitu isipokuwa historia ya mazungumzo hunakiliwa. `-g` au `--group` huchagua kikundi, na `-y` au `--yes` hunakili bila kuuliza na huhifadhi nakala ya kikundi faili inapotofautiana. Tazama [kunakili usanidi uliopo](/docs/sharing#copy-an-existing-setup).

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
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Matengenezo

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
```

## Faili

| Njia | Yaliyomo |
| --- | --- |
| `~/.config/ccam/config.toml` | Akaunti na zinachoshiriki |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ya kila akaunti |
| `~/.local/share/ccam/groups/<group>/` | Mipangilio, seva za MCP, viendelezi na historia vinavyoshirikiwa |
| `~/.local/state/ccam/` | Cache ya matumizi na nyakati za matumizi ya mwisho |
| `~/.local/bin/<shortcut>` | Amri za njia fupi kama `ccwork`, ambazo ni viungo vya ccam |

Vigezo `XDG_CONFIG_HOME`, `XDG_DATA_HOME` na `XDG_STATE_HOME` huhamisha folda hizi.
