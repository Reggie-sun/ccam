---
title: Marejeo ya amri
metaTitle: Marejeo ya amri za ccam
description: Kila amri na chaguo la ccam, kuanzia ccam add na ccam <name> hadi ccam usage, ccam next na ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 5d3275878535
---

Weka `--dry-run` kabla ya amri ili kuona itakachofanya bila kukifanya, kwa mfano `ccam --dry-run work`. Inafanya kazi unapoanzisha Claude Code (`ccam <name>`, `ccam run`, `ccam next`) na pamoja na `share`, `migrate`, `shortcut`, `update` na `doctor --fix`. Amri nyingine zinazobadilisha kitu huikataa badala ya kufanya kazi.

<!-- commands -->

## Kuanzisha Claude Code

`ccam <name|number> [arguments]` huanzisha Claude Code kama akaunti hiyo. ccam haisomi hoja zinazokuja baada ya jina; huzipitisha kwa `claude` kwa mpangilio ule ule, na standard input, exit code na signals hufanya kazi sawasawa kama kwa `claude` yenyewe. `--` baada ya jina ni hiari.

`ccam` ikiwa na hoja lakini bila jina la akaunti, kama `ccam -p "summarise"`, hutumia akaunti iliyowekwa kwa folda ya sasa, la sivyo akaunti chaguo-msingi. Ziweke kwa `ccam folder <name> [folder]` na `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` kabla ya jina la akaunti ni kwa uendeshaji wa `-p`/`--print`. Akaunti inapofikia kikomo chake cha matumizi, ccam huendesha amri ile ile tena kama akaunti iliyoingia yenye nafasi kubwa zaidi na huchapisha matokeo ya uendeshaji uliokamilika pekee. Tazama [Akaunti kwa folda](/features/folder) na [Endesha tena kwenye kikomo](/features/next#run-again-at-the-limit).

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

`ccam rename`, `ccam remove` na `ccam share` hukataa kubadilisha akaunti wakati Claude Code inaendeshwa kama akaunti hiyo. Funga kipindi (session) kwanza.

`ccam list --json` huchapisha kila akaunti pamoja na namba yake (slot), njia fupi, barua pepe na folda. Kitu chake `share` huhusisha kila kipengele (`mcp`, `settings`, `memory`, `extensions`, `history`) na jina la kikundi au `local`.

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

`ccam migrate` hunakili `settings.json`, `CLAUDE.md`, viendelezi na seva za MCP za mtumiaji kutoka folda ya Claude Code kwenda kwenye kikundi, ili kila akaunti inayoshiriki vipengele hivyo ivipate. Folda husomwa tu. Chaguo `--settings`, `--memory`, `--extensions`, `--mcp` na `--history` hunakili tu vipengele unavyovitaja. Bila hivyo, kwenye terminal orodha ya kuteua hukuruhusu kuchagua, huku kila kitu isipokuwa historia ya mazungumzo kikiwa kimeteuliwa; mahali pengine, kila kitu isipokuwa historia ya mazungumzo hunakiliwa. `-g` au `--group` huchagua kikundi, na `-y` au `--yes` hunakili bila kuuliza na huhifadhi nakala ya kikundi faili inapotofautiana. Pia hunakili taarifa za kuingia zilizohifadhiwa za seva hizo kwenye faili ya kuingia ya kila akaunti yenyewe. Tazama [kunakili usanidi uliopo](/docs/sharing#copy-an-existing-setup).

## Seva za MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` huhamisha seva ambazo akaunti hiyo pekee inazo kwenda kwenye kikundi chake. Seva ambayo jina lake tayari linatumiwa na kikundi kwa ufafanuzi tofauti huachwa; ongeza `--force` ili kubadilisha nakala ya kikundi, ambayo ccam huitengenezea nakala rudufu kwanza.

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
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` hukagua checksum ya SHA-256 ya toleo kabla ya kubadilisha ccam, na `--version v0.2.0` husakinisha toleo mahususi. ccam iliyosakinishwa kutoka kifurushi cha `.deb` husasishwa kwa apt badala yake. Tazama [kusasisha ccam](/docs/install#update).

## Faili

| Njia | Yaliyomo |
| --- | --- |
| `~/.config/ccam/config.toml` | Akaunti na zinachoshiriki |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ya kila akaunti |
| `~/.local/share/ccam/groups/<group>/` | Mipangilio, seva za MCP, viendelezi na historia vinavyoshirikiwa |
| `~/.local/state/ccam/` | Cache ya matumizi na nyakati za matumizi ya mwisho |
| `~/.local/bin/<shortcut>` | Amri za njia fupi kama `ccwork`, ambazo ni viungo vya ccam |

Vigezo `XDG_CONFIG_HOME`, `XDG_DATA_HOME` na `XDG_STATE_HOME` huhamisha folda hizi. Kwenye Windows, `~` ni `%USERPROFILE%` na njia fupi ni `<shortcut>.exe`.
