---
title: Marejeo ya amri
metaTitle: Marejeo ya amri za ccam
description: Kila amri na chaguo la ccam, kuanzia ccam add na ccam <name> hadi ccam usage, ccam next na ccam doctor.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

Weka `--dry-run` kabla ya amri ili kuona itakachofanya bila kukifanya, kwa mfano `ccam --dry-run work`. Inafanya kazi unapoanzisha Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) na pamoja na `group`, `migrate`, `shortcut`, `update` na `doctor --fix`. Amri nyingine zinazobadilisha kitu huikataa badala ya kufanya kazi.

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
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` bila folda hutafuta `~/.claude` na `~/.claude-*` na huuliza kuhusu kila moja. Folda zilizoingizwa husajiliwa pale zilipo; ccam kamwe haihamishi wala kunakili taarifa za kuingia.

`ccam add --api` huongeza akaunti inayotumia API key badala ya kuingia kwa Claude; ccam husoma key kila inapoanzishwa na kamwe haiihifadhi. Tazama [Akaunti za API key](/features/api).

`ccam rename`, `ccam remove` na amri za `ccam group` hukataa kubadilisha akaunti wakati Claude Code inaendeshwa kama akaunti hiyo. Funga kipindi (session) kwanza.

`ccam list --json` huchapisha kila akaunti pamoja na namba yake (slot), njia fupi, barua pepe na folda. Safu yake `groups` huorodhesha vikundi vya akaunti, na kitu chake `share` huhusisha kila kipengele (`mcp`, `settings`, `memory`, `extensions`, `history`) na kikundi kinachotoka, au `local`.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Kikundi hushiriki vipengele kati ya akaunti zake; akaunti inaweza kuwa katika vikundi kadhaa, lakini kila kipengele hutoka kwenye kikundi kimoja tu. Tazama [kushiriki kwa vikundi](/features/share).

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
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` na `ccam pick` huchagua kwa njia ile ile: kwanza zilizo chini ya `auto_max_usage` (`config.toml`, chaguo-msingi 100), kwa `auto_strategy` (`headroom`, `expiring` au `balance`), zikihesabu kipindi cha wiki cha `--model`, kamwe si akaunti iliyo kwenye 100%, na ndani ya shirika la akaunti wanayoanzia (`--cross-org` au `cross_org = true` ili kutafuta zaidi). Isipowezekana kuchagua yoyote, hutoka na msimbo 75, au kwa `--wait` husubiri kipindi cha akaunti ya kwanza kianze upya. `ccam next --continue` huhamisha mazungumzo hata kati ya akaunti zisizoshiriki historia: hapo Claude Code huyanakili kwenye akaunti mpya. Tofauti na `ccam <name>`, `ccam next` husoma chaguo zake yenyewe, kwa hiyo chaguo za Claude Code huwekwa baada ya `--`; bila hiyo ccam husimama na `unknown flag`. Ndani ya Claude Code, `/ccam:switch [account]` ya plugin hukabidhi mazungumzo, na `ccam hook` ndicho ambacho hooks zake huendesha. Tazama [Kubadili kwenda akaunti inayofuata](/features/next), [Kuchagua akaunti kwa ajili ya script](/features/pick) na [Vipindi vya chinichini vya kila akaunti](/features/agents).

## Matengenezo

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

`ccam update` hukagua checksum ya SHA-256 ya toleo kabla ya kubadilisha ccam, na `--version v0.2.0` husakinisha toleo mahususi. ccam iliyosakinishwa kutoka kifurushi cha `.deb` husasishwa kwa apt badala yake. Tazama [kusasisha ccam](/docs/install#update).

## Faili

| Njia | Yaliyomo |
| --- | --- |
| `~/.config/ccam/config.toml` | Akaunti, vikundi na mipangilio (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` ya kila akaunti |
| `~/.local/share/ccam/groups/<group>/` | Mipangilio, seva za MCP, viendelezi na historia vinavyoshirikiwa |
| `~/.local/state/ccam/` | Cache ya matumizi, nyakati za matumizi ya mwisho na maelezo ambayo `/ccam:switch` na hook ya kikomo huacha (`handoff/`) |
| `~/.local/bin/<shortcut>` | Amri za njia fupi kama `ccwork`, ambazo ni viungo vya ccam |

Vigezo `XDG_CONFIG_HOME`, `XDG_DATA_HOME` na `XDG_STATE_HOME` huhamisha folda hizi. Kwenye Windows, `~` ni `%USERPROFILE%` na njia fupi ni `<shortcut>.exe`.
