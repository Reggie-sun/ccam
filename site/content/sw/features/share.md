---
title: Kushiriki kwa vikundi
description: Kikundi cha ccam hushiriki vipengele ulivyochagua kati ya akaunti zake kupitia folda moja; ccam group huhamisha data kwa mpango, nakala rudufu na chaguo.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Kila akaunti ina folda yake, lakini sehemu kubwa ya yaliyomo ndani yake haihitaji kuwa tofauti. **Kikundi** husema kinachoshirikiwa na nani: kila akaunti katika kikundi huunganisha vipengele hivyo kwenye folda moja, `~/.local/share/ccam/groups/<group>/`, ili zote zisome na ziandike nakala ile ile. Tokeni za kuingia kamwe hazishirikiwi.

Akaunti inaweza kuwa katika vikundi kadhaa, lakini kila kipengele hutoka kwenye kikundi kimoja tu. Mpangilio wa kawaida ni kikundi kimoja cha mipangilio na viendelezi, na cha pili cha historia ya mazungumzo kati ya akaunti mbili:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Akaunti ya kwanza unayoongeza huunda `default`, na akaunti zinazofuata hujiunga nacho usiposema vinginevyo.

## Matumizi

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Kila amri inayohamisha data huchapisha mpango, huuliza kuhusu kila faili inayotofautiana na nakala ya kikundi, na huuliza `Proceed? [y/N]` kabla ya kubadilisha chochote:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Mabadiliko hayo hayo yako kwenye kichupo cha **Groups** cha [dashibodi](/features/dashboard), na `s` kwenye akaunti huchagua vikundi vyake.

## Vipengele

| Kipengele | Kinachojumuisha |
| --- | --- |
| `mcp` | seva za MCP, kupitia `mcp.json` ya kikundi ([Seva za MCP zinazoshirikiwa](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` ya jumla |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` na `history.jsonl`, ili akaunti nyingine iweze `--continue` |

Kipengele ambacho hakuna kikundi cha akaunti kinachokishiriki hubaki katika folda ya akaunti yenyewe.

## Amri na chaguo

| Amri | Inachofanya |
| --- | --- |
| `ccam group` au `ccam group list` | huorodhesha vikundi |
| `ccam group show <group>` | kile kikundi kinachoshiriki, akaunti zake, folda yake na seva za MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | huunda kikundi; vipengele vya akaunti zake huhamia ndani yake |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | hubadilisha kile kikundi kinachoshiriki, kwa kila akaunti iliyomo |
| `ccam group join <group> <account>...` | huweka akaunti katika kikundi |
| `ccam group leave <group> <account>...` | hutoa akaunti; kila moja hubaki na nakala |
| `ccam group rename <group> <new-name>` | hubadilisha jina katika usanidi pekee |
| `ccam group remove <group> [--leave]` | huondoa kikundi |

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: akaunti ambayo nakala yake huingia kwenye kikundi kwanza | akaunti chaguo-msingi, kisha kwa namba |
| `--replace` | `join`: ondoka kwanza kwenye vikundi vinavyoshiriki vipengele hivyo hivyo | imezimwa |
| `--fresh` | akaunti zinazoacha kushiriki kipengele huanza bila kipengele hicho badala ya kuwa na nakala | imezimwa |
| `--leave` | `remove`: toa akaunti za kikundi kwanza | inahitajika kikundi kikiwa na akaunti |
| `--prefer <group\|mine\|merge>` | jinsi ya kuamua migongano | huuliza; kwa `--yes`, chaguo-msingi la kila faili |
| `-y`, `--yes` | huendesha bila kuuliza; migongano hufuata `--prefer` | imezimwa |
| `--dry-run` (jumla) | huchapisha mpango, haibadilishi chochote | imezimwa |

Akaunti hutajwa kwa jina au kwa namba; vipengele kama orodha iliyotenganishwa kwa koma. `ccam add` hupokea `--group <group>` (inaweza kurudiwa) au `--no-group`.

## Jinsi inavyofanya kazi

Faili zinazohusika:

- `~/.config/ccam/config.toml`: maingizo ya `[[group]]`: `name`, `items`, `members` na, jina la folda likiwa tofauti, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` ya akaunti.
- `~/.local/share/ccam/groups/<folder>/`: nakala zinazoshirikiwa.

`config.toml` inaposomwa, vikundi vya kila akaunti hutatuliwa kuwa folda moja kwa kila kipengele. Hili huchukua takriban 25 µs kwa akaunti tano, kwa hiyo uanzishaji bado husoma `config.toml` pekee. Usanidi ambao akaunti ingepata kipengele kimoja kutoka vikundi viwili ni kosa, na hakuna kinachoanza hadi lirekebishwe.

Badiliko hufanya kazi kwenye nakala ya usanidi:

1. Amri huhariri nakala hiyo (`AddGroup`, `SetGroupItems`, `JoinGroup` …), ambayo hukagua kanuni.
2. `internal/group` hulinganisha ya zamani na mpya na, kwa kila akaunti ambayo vipengele vyake vinabadilika, huunda mpango kwa `share.Build`. Mipango huundwa mmoja baada ya mwingine, na kila mmoja huona kile ambacho mipango ya awali itaweka kwenye folda za vikundi, kwa hiyo akaunti ya kwanza pekee ndiyo huhamisha nakala yake ndani. Nakala ya kwanza hutoka kwenye folda ya kikundi ikiwa tayari ina moja, kisha `--from`, kisha akaunti chaguo-msingi, kisha nyingine kwa namba.
3. ccam hukataa ikiwa mojawapo ya akaunti hizo inaendeshwa, na kwenye Windows bila Developer Mode ikiwa faili ingehitaji kiungo.
4. Migongano huamuliwa, kwa swali au kwa `--prefer`, mpango huchapishwa, na ccam huuliza.
5. Usanidi mpya huhifadhiwa kabla ya kitu chochote kuhamishwa. Hatua ikisimama katikati, `ccam doctor --fix` na kuendesha tena huunda mpango wao kuelekea vikundi vilivyohifadhiwa na hukamilisha badiliko.
6. Mipango huendeshwa kwa mpangilio. Kila hatua huhamisha data kwenda kwenye kikundi au kwenye folda ya nakala rudufu; hakuna kinachofutwa.

`share.Layout` ni orodha ya ruhusa ya njia ambazo kila kipengele hujumuisha. Kitu chochote kisichokuwamo kamwe hakiwezi kufika kwenye kikundi. `share.Never` hutaja `.credentials.json` na `.claude.json`, na `Build` husimama na kosa la ndani ikiwa yoyote kati ya hizo itaonekana kwenye mpangilio. Kipengele cha `mcp` hakina njia: orodha ya seva iko ndani ya `.claude.json`, kwa hiyo [kushiriki MCP](/features/mcp) hufanya kazi kupitia `mcp.json` ya kikundi.

Hiki ndicho kinachotokea kwa folda moja, `projects/`, wakati `side` inaanza kushiriki historia na kikundi tayari kina maudhui:

<!-- motion:share-on -->

Ikiwa kikundi bado hakina chochote kwenye njia hiyo, akaunti ndiyo ya kwanza kuishiriki na folda yake huhamishwa huko kama ilivyo. Vinginevyo kila faili ambayo kikundi hakina huhamishwa kwenda kwenye kikundi, na faili zilizopo pande zote mbili hubaki katika folda ya akaunti, kwa hiyo toleo la kikundi huhifadhiwa. Kisha folda ya akaunti huenda `.ccam-backup/<time>/` ndani ya folda ya akaunti, na kiungo huchukua nafasi yake. Faili za mazungumzo hupewa majina kwa kufuata session ID yake, kwa hiyo historia kutoka akaunti kadhaa huungana bila migongano. Rejista za `plugins/` huhifadhi njia kamili kupitia folda ya akaunti; ccam huzielekeza kwenye folda ya kikundi na huunganisha vipengele vya akaunti kwenye rejista ambayo kikundi tayari kinayo.

Kuondoka hugeuza kiungo:

<!-- motion:share-off -->

ccam huondoa kiungo tu ikiwa kinaelekea ndani ya folda ya vikundi, kisha hunakili maudhui ya sasa ya kikundi kwenda kwenye akaunti. Kwa `--fresh` kunakili kunarukwa. Kikundi huendelea kuhifadhi maudhui yake kwa ajili ya akaunti nyingine. Akaunti inayohamisha kipengele kwenda kikundi ambacho folda yake bado iko tupu huleta nakala ya kile ilichokiona, kwa hiyo huhifadhi faili zake; seva za MCP huja vivyo hivyo. Akaunti inayoacha kushiriki MCP hupoteza seva ambazo ccam ilikuwa imezinakili kwenye `.claude.json` yake (nakala rudufu huhifadhiwa), kwa kuwa zilitoka kwenye kikundi.

Faili moja iliyopo pande zote mbili yenye maudhui tofauti ni mgongano:

<!-- motion:share-conflict -->

`history.jsonl` hutoa `merge by line`, `keep group` na `use mine`, uunganishaji ukiwa wa kwanza. Faili nyingine za `.json` hutoa `keep group`, `use mine` na `merge JSON`. Kitu kingine chochote hutoa `keep group` na `use mine`. `merge JSON` huunganisha vitu kwa kujirudia, huunganisha safu bila marudio, na huhifadhi thamani ya kikundi wakati pande zote mbili zinaweka ufunguo ule ule kwa njia tofauti, kwa hiyo akaunti nyingine hazioni badiliko. Kila chaguo isipokuwa `keep group` kwanza huhifadhi nakala ya kikundi katika `.ccam-backup/<time>/group/` ya akaunti. Faili yako mwenyewe daima huenda kwenye nakala rudufu.

Katika kila uanzishaji, `share.Check` hulinganisha kila njia katika mpangilio na vikundi vya akaunti kwa kutumia `lstat` na `readlink` tu. Kiungo kinachokosekana huundwa papo hapo. Faili halisi mahali ambapo kiungo kinapaswa kuwa, kiungo kinachoelekea mahali pasipo sahihi, au kiungo kilichopitwa na wakati husimamisha uanzishaji na huelekeza kwenye `ccam doctor --fix`.

### Kubadilisha jina, kuondoa na muundo wa zamani

- `ccam group rename` hubadilisha usanidi pekee. Folda huhifadhi jina lake na kikundi hulirekodi kama `dir`, kwa hiyo hakuna kinachohama na hakuna kinachoweza kupotea.
- `ccam group remove` hukataa kikundi chenye akaunti usipoongeza `--leave`. Kisha kila akaunti hubaki na nakala, na folda huenda `~/.local/share/ccam/trash/` isipokuwa kikundi kingine kinaitumia.
- Kabla ya vikundi, `config.toml` ilihifadhi jedwali la `[profile.share]` kwa kila akaunti. ccam husoma muundo huo kama vikundi bila kuhamisha chochote. Katika uhifadhi wa kwanza huandika muundo mpya na huhifadhi faili ya zamani kama `config.toml.bak`. Vipengele vinavyoshirikiwa na akaunti zile zile huwa kikundi kimoja chenye jina la folda; vingine hupata majina kama `default-history` na huendelea kutumia folda ile ile.

### Kwenye Windows

Windows huunda viungo vya ishara (symbolic links) katika Developer Mode pekee (Settings > System > For developers) au kwa msimamizi. Bila haki hiyo, folda inayoshirikiwa huwa directory junction, ambayo hufanya kazi vivyo hivyo kwa Claude Code. Faili inayoshirikiwa (`settings.json`, `CLAUDE.md`, `history.jsonl`) inaweza kuwa kiungo cha ishara tu, kwa hiyo badiliko la kikundi linalohitaji kiungo kama hicho hukataliwa kabla ya kitu chochote kuhamishwa.

## Usalama na hali maalum

- **Akaunti zinazoendeshwa.** Badiliko hukataliwa wakati Claude Code inaendeshwa kama mojawapo ya akaunti linazozihamisha, kwa sababu kipindi kinaendelea kuandika kwenye faili hizo. Kuondoka na kujiunga huhusu akaunti iliyotajwa pekee. `--dry-run` bado huchapisha mpango.
- **Nakala rudufu.** Kila mpango hupata folda yake ya `.ccam-backup/<time>/` katika folda ya akaunti. ccam kamwe haizifuti.
- **Viungo vyako mwenyewe.** Ikiwa njia tayari ni kiungo ulichokiunda, kwa mfano kuelekea hazina ya dotfiles, na kikundi kiko tupu, ccam hunakili lengo la kiungo kwenda kwenye kikundi na huhamisha kiungo chenyewe kwenye nakala rudufu. Lengo haliguswi.
- **Folda zilizoingizwa.** Akaunti iliyosajiliwa pale ilipo kwa `ccam import` (kama `~/.claude`) hubadilika pale ilipo, kwa hiyo `claude` ya kawaida huona badiliko pia; mpango husema hivyo.
- **Akaunti mbili katika folda moja.** Akaunti zinazoshiriki historia huandika kwenye faili zile zile. Akaunti ya pili inapoanza katika folda ambamo moja tayari inaendeshwa, ccam huonya na huuliza.
- **Mabadiliko ya mkono.** Mwanachama ambaye si akaunti hupuuzwa na huripotiwa na `ccam doctor`. Majina `local`, `on`, `off` na `none` hayawezi kuwa majina ya vikundi.

## Zinazohusiana

- [Dashibodi](/features/dashboard)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Seva za MCP zinazoshirikiwa](/features/mcp)
- [Nakili usanidi uliopo](/features/migrate)
- [Amri zote](/docs/commands)
