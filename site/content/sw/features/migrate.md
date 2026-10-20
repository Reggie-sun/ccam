---
title: Kunakili usanidi uliopo
description: Nakili mipangilio, CLAUDE.md, viendelezi, seva za MCP na taarifa za kuingia za MCP za ~/.claude iliyopo hadi kikundi cha ccam kwa ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

Ikiwa ulitumia Claude Code kabla ya ccam, usanidi wako uko katika `~/.claude` na `~/.claude.json`. `ccam migrate` huunakili kwenye kikundi cha kushiriki, ili kila akaunti inayoshiriki vipengele hivyo kupitia kikundi ianze na usanidi ulio nao tayari. Folda chanzo husomwa tu, na `claude` ya kawaida inaendelea kufanya kazi kama awali. Iendeshe mara moja baada ya `ccam add` yako ya kwanza, na tena kila akaunti inapoanza kushiriki seva za MCP.

## Matumizi

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Kwenye terminal, bila chaguo za vipengele na bila `--yes`, ccam huonyesha kwanza orodha ya kuteua:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (au `x`) hubadilisha kipengele kilicho chini ya kishale, vitufe vya mishale au `j`/`k` husogeza, Enter huonyesha mpango, na Esc, `q` au Ctrl+C huchapisha `Cancelled.`. Vipengele ambavyo chanzo hakina chochote kwavyo huwekwa alama `not found`. Ukiondoa uteuzi wa kila kitu, ccam huchapisha `Nothing selected.` na husimama.

Kisha ccam huchapisha mpango na huuliza kabla ya kunakili:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Kikundi hufikia akaunti zake tu, na kwa vipengele kinavyoshiriki tu, kwa hiyo muhtasari huorodhesha akaunti za kikundi. Kikundi kisiposhiriki mojawapo ya vipengele vilivyonakiliwa, au baadhi ya akaunti hazimo ndani yake, ccam huchapisha amri za `ccam group set <group> --add <items>` na `ccam group join <group> <accounts>` zinazobadilisha hilo.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `[folder]` | Folda ya Claude Code ya kunakili kutoka | `~/.claude` |
| `--settings` | Hunakili `settings.json` | imezimwa |
| `--memory` | Hunakili `CLAUDE.md` ya jumla | imezimwa |
| `--extensions` | Hunakili `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | imezimwa |
| `--mcp` | Hunakili seva za MCP za mtumiaji (`mcpServers` katika `.claude.json`) | imezimwa |
| `--history` | Hunakili historia ya mazungumzo (`projects/`, `history.jsonl`) | imezimwa |
| `-g`, `--group string` | Kikundi cha kushiriki cha kunakili ndani yake | `default` |
| `-y`, `--yes` | Hunakili bila kuuliza; kukiwa na mgongano nakala ya kikundi huhifadhiwa | imezimwa |
| `--dry-run` (jumla) | Huchapisha mpango na husimama kabla ya swali au unakili wowote | imezimwa |

Bila chaguo la kipengele, orodha ya kuteua huamua kwenye terminal. Nje ya terminal, au kwa `--yes`, ccam hunakili kila kitu isipokuwa historia ya mazungumzo. Historia huachwa kwa chaguo-msingi kwa sababu `projects/` inaweza kuwa kubwa.

`--group on` inamaanisha `default`. `--group off` na `--group local` hukataliwa kwa `--group takes a group name, not "off"`. `ccam migrate` huheshimu `--dry-run`, na orodha ya kuteua bado huonekana katika dry run ili uweze kuona uteuzi kabla.

## Jinsi inavyofanya kazi

`internal/cli/migrate.go` huchagua vipengele, kisha `migrate.Build` husoma chanzo na kikundi na hurudisha mpango. Hakuna kinachoandikwa hadi `Plan.Apply`, baada ya wewe kuthibitisha.

1. **Kupata chanzo.** `checkSource` hukataa njia inayokosekana, isiyo folda, au isiyo na kipengee chochote kinachojulikana cha Claude Code. `migrate.ClaudeJSON` hupata `.claude.json` inayohusiana na folda: ndani yake, au kwa `~/.claude` isiyo na moja ndani, `~/.claude.json` kando yake. Hapo ndipo Claude Code huihifadhi wakati `CLAUDE_CONFIG_DIR` haijawekwa.
2. **Kuhusisha vipengele na njia.** Kila njia hutoka `share.Layout`, orodha ile ile ya ruhusa ambayo vikundi hutumia. Seva za MCP husomwa kutoka `.claude.json` kwa `mcp.LocalServers`, ambayo hutumia kisomaji cha `profile.TopLevelValue` badala ya kufasiri faili nzima.
3. **Kunakili, kamwe si kuhamisha.** Faili huenda `~/.local/share/ccam/groups/<group>/` kwa `fsutil.AtomicWrite` au `fsutil.CopyFile`; seva za MCP huenda kwenye `mcp.json` ya kikundi.

Kila kipengele huhusisha njia moja ya chanzo na njia moja ya kikundi:

<!-- motion:migrate-items -->

Safu mbili za historia zina mistari ya vistari kwa sababu hunakiliwa tu kwa `--history`. `.credentials.json` na `.claude.json` zimeorodheshwa katika `share.Never` na kamwe hazifiki kwenye kikundi: `Build` hurudisha kosa la ndani ikiwa moja wapo itaonekana kwenye mpangilio, na upitiaji wa folda huziruka kwa jina.

`buildEntry` huipa kila njia hali, ambayo ndiyo kila mstari wa mpango huieleza:

<!-- motion:migrate-states -->

- **Missing**: haipo kwenye chanzo, au ni faili mahali ambapo folda inatarajiwa. Huorodheshwa mara moja kwa kila kipengele kama `not found, skipped`.
- **Linked**: njia ya chanzo tayari inaelekea nakala ya kikundi, kwa mfano chanzo kinapokuwa folda ya wasifu wa ccam.
- **Same**: kikundi tayari kina maudhui yanayofanana kabisa.
- **New**: kikundi hakina, kwa hiyo hunakiliwa yote.
- **Merge**: folda iliyo pande zote mbili. `copyMissing` huongeza tu faili ambazo kikundi hakina; faili zinazotofautiana hubaki kama kikundi kilivyo nazo na mpango husema `group copies kept`. Kwa `plugins/`, rejista `installed_plugins.json` na `known_marketplaces.json` huunganishwa, na njia kamili za plugins huandikwa upya ili kuelekea folda ya kikundi.
- **Conflict**: faili iliyo pande zote mbili yenye maudhui tofauti. Unachagua kwa kila faili: `keep group`, `use mine`, pamoja na `merge JSON` kwa faili za `.json` (uunganishaji wa kina ambapo kikundi hushinda kwenye thamani zinazogongana) au `merge by line` kwa `history.jsonl`. Seva ya MCP ambayo ufafanuzi wake unatofautiana pia ni mgongano, ikiwa na `keep group` au `use mine`. Chaguo la kwanza lililoorodheshwa ndilo chaguo-msingi unapobonyeza Enter; `--yes` daima huhifadhi nakala ya kikundi.

Orodha ya kuteua yenyewe ni `tui.MigratePicker` iliyojengwa juu ya kipengele cha `checklist` ambacho `ccam add` hutumia; `migrate.Found` huweka alama vipengele vinavyokosekana kwa miito ya lstat na orodha ya seva za MCP.

### Taarifa za kuingia za MCP

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao kwenye faili ya kuingia, `.credentials.json` chini ya `mcpOAuth`, si kwenye orodha ya seva. Claude Code huipa kila kipengele ufunguo wa `<name>|` pamoja na herufi 16 za kwanza za hash ya `sha256` ya `type`, `url` na `headers` za seva. Ufunguo hautegemei folda, kwa hiyo kipengele kilichonakiliwa hufanya kazi katika akaunti nyingine mradi ufafanuzi wa seva unafanana kabisa. Kwa hiyo `buildLogins` katika `internal/migrate/logins.go` hunakili tu vipengele vya seva ambazo kikundi kinaishia kuwa nazo kwa ufafanuzi wa chanzo: mpya, tayari zinafanana, au mgongano ambapo ulichagua `use mine`. Taarifa za kuingia za seva za plugins (funguo zinazoanza na `plugin:`) huenda pamoja unaponakili viendelezi na chanzo kina `plugins/`.

<!-- motion:migrate-signins -->

Vipengele huenda kwenye `.credentials.json` ya kila akaunti inayopokea yenyewe, vikiandikwa kwa `fsutil.AtomicWrite` kwa mode `0600`. Kamwe haviendi kwenye kikundi. Akaunti huvipokea tu ikiwa inashiriki `mcp` (au, kwa seva za plugins, `extensions`) kupitia kikundi lengwa, na funguo ambazo haina pekee ndizo huongezwa: vipengele ambavyo tayari inavyo, huenda vikiwa vipya zaidi, huhifadhiwa. Mstari wa mpango wa `sign-ins` huorodhesha seva, akaunti zinazopokea na zile zilizorukwa. ccam huruka akaunti, na huitaja pamoja na sababu, wakati:

- imesajiliwa mahali ilipo kwa `ccam import`, kwa kuwa `claude` ya kawaida inaweza kuwa inatumia faili yake ya kuingia;
- inaendeshwa, kwa kuwa Claude Code huandika upya faili yake ya kuingia (hukaguliwa tena kabla tu ya kuandika, huripotiwa kama `started meanwhile`);
- bado haijaingia;
- huhifadhi taarifa zake za kuingia katika Keychain ya macOS, ambayo ccam haiandiki.

Ingia kwenye seva katika akaunti iliyorukwa kwa `/mcp`. Akaunti inayoanza kushiriki seva za MCP baadaye hupata taarifa hizo za kuingia unapoendesha `ccam migrate` tena; muhtasari husema hivyo kwa ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Usalama na hali maalum

- Chanzo husomwa tu. Hakuna faili ndani yake inayohamishwa, kuunganishwa au kubadilishwa, na taarifa za kuingia za chanzo kamwe haziandikwi, pamoja na kwenye macOS ambapo zinaweza kusomwa kutoka Keychain.
- Kabla ya faili ya kikundi kubadilishwa, kwa `use mine`, uunganishaji, au seva mpya za MCP katika `mcp.json` iliyopo, nakala ya kikundi huhifadhiwa katika `<group>/.ccam-backup/<time>/`. Jina la folda ni la kipekee kwa kila uendeshaji, kwa hiyo uendeshaji miwili katika sekunde ile ile haiandiki juu ya nakala rudufu za mwingine. Muhtasari huchapisha njia.
- Kwa `--yes`, kila mgongano huhifadhi nakala ya kikundi, na ccam husema jinsi ya kuchagua badala yake: endesha amri ile ile bila `--yes`.
- Ikiwa hakuna kitakachobadilika, ccam huchapisha `Nothing to copy.`; ikiwa kila mgongano huhifadhi nakala ya kikundi, `Nothing to copy: the group keeps its own copies.`
- Unakili ukishindwa katikati, ccam husimama na `The source folder was not changed. Group files replaced so far are saved in` ikifuatiwa na folda ya nakala rudufu.
- Chanzo kinachokosekana huishia na `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Zinazohusiana

- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Kukagua na kurekebisha](/features/doctor)
- [Akaunti kadhaa za Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
