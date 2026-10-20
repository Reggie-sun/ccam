---
title: Vipindi vya chinichini vya kila akaunti
description: ccam agents huorodhesha agents wa chinichini na vipindi wazi vya Claude Code vya akaunti zote katika jedwali moja, na ccam attach hufungua kimoja kwa id yake.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code huhifadhi vipindi vya chinichini (`claude --bg`, mwonekano wa agents) kwa kila akaunti: `claude agents` huonyesha tu vipindi vya akaunti ambayo inaendeshwa kama hiyo. `ccam agents` huuliza kila akaunti kwa wakati mmoja na huviorodhesha pamoja, na `ccam attach <id>` hufungua kipindi cha chinichini kama akaunti yoyote iliyokianzisha.

## Matumizi

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

`State` ni hali ya kipindi cha chinichini (`working`, `blocked`, `done`, `failed`, `stopped`) na kile ambacho mchakato wake unafanya (`busy`, `waiting`, `idle`); kipindi cha mwingiliano kilicho wazi huonyesha `open`.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--all` | Orodhesha pia vipindi vya chinichini vilivyokamilika, vilivyoshindwa au vilivyosimamishwa | imezimwa |
| `--json` | Chapisha vipindi kama JSON, kila kimoja pamoja na `account` kinachomilikiwa nayo | imezimwa |

`ccam attach <id> [claude args]` hupitisha hoja zilizo baada ya id kwa `claude attach`. Kwa `--dry-run` huchapisha amri.

## Jinsi inavyofanya kazi

1. Kwa kila akaunti, ccam huendesha `claude agents --json` (pamoja na `--all` kikitolewa) kwa `CLAUDE_CONFIG_DIR` ya akaunti hiyo, zote kwa wakati mmoja, kila moja ndani ya sekunde 15. Hili husoma orodha ya vipindi ya akaunti yenyewe na halianzishi huduma ya chinichini ya Claude Code.
2. Kila ingizo hupata jina la akaunti, na orodha hupangwa kwa muda wa kuanza, vipya kwanza. Akaunti isiyojibu hutajwa kwenye standard error na huachwa nje, kisha ccam hutoka na msimbo 1, ili script isichukulie kushindwa kuwa orodha tupu.
3. `ccam attach` hutafuta `jobs/<id>` katika folda ya kila akaunti, mahali ambapo Claude Code huhifadhi kipindi cha chinichini, na huendesha `claude attach <id>` kama akaunti hiyo, kwa njia ile ile ambayo `ccam <name>` huanzisha Claude Code.

## Usalama na hali maalum

- Id za vipindi vya chinichini ni za akaunti moja: `claude attach` ikiendeshwa kama akaunti nyingine haiwezi kuzipata, na hilo ndilo `ccam attach` hutatua.
- Akaunti mbili zinaweza kuwa na vipindi ambavyo id zake zinaanza kwa namna ile ile; hapo `ccam attach` huzitaja zote mbili na huomba sehemu zaidi ya id.
- ccam hupitisha seva zake za MCP zinazoshirikiwa kwa Claude Code kwa `--mcp-config=<file>` katika hoja moja, ambayo `claude agents` na `claude attach` hukubali; vipindi vya chinichini huendelea kutumia faili hiyo vinapoanza upya.
- Huduma ya chinichini ya akaunti ya API key huchukua ufunguo wake kutoka kwenye shell iliyoianzisha. Anzisha kipindi cha kwanza cha chinichini cha akaunti kama hiyo kupitia ccam (`ccam gw --bg …`), ili huduma iwe na ufunguo.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Kuchagua akaunti kwa ajili ya script](/features/pick)
- [Dashibodi](/features/dashboard)
