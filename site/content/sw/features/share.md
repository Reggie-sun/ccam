---
title: Kushiriki kati ya akaunti
description: Jinsi ccam share huunganisha mipangilio, CLAUDE.md, viendelezi na historia kwa folda ya kikundi kwa mpango, nakala rudufu na chaguo za mgongano; off hurudisha.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Kila akaunti ina folda yake, lakini sehemu kubwa ya yaliyomo ndani yake haihitaji kuwa tofauti. `ccam share` huwasha au huzima kushiriki kwa akaunti moja na kipengele kimoja kwa wakati: faili au folda ya akaunti huhamishwa kwenda folda ya kikundi na hubadilishwa kwa kiungo (symbolic link), ili kila akaunti katika kikundi isome na iandike nakala ile ile. Itumie baada ya `ccam add` kubadilisha kile akaunti inachoshiriki, au kuweka akaunti chache katika kikundi chao wenyewe.

## Matumizi

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Bila chaguo za vipengele, `ccam share` huchapisha hali ya sasa:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Ukiwa na chaguo, ccam huchapisha mpango, huuliza kuhusu kila faili inayotofautiana na nakala ya kikundi, na huuliza `Proceed? [y/N]` kabla ya kubadilisha chochote. Akaunti mpya hushiriki kila kitu isipokuwa historia ya mazungumzo; [Shiriki seva za MCP, mipangilio na historia](/docs/sharing) inaeleza kila kipengele kwa mtazamo wa mtumiaji.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Seva za MCP | haibadiliki |
| `--settings <on\|off\|group>` | `settings.json` | haibadiliki |
| `--memory <on\|off\|group>` | `CLAUDE.md` ya jumla | haibadiliki |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | haibadiliki |
| `--history <on\|off\|group>` | `projects/` na `history.jsonl` | haibadiliki |
| `--fresh` | unapozima kushiriki, anza ukiwa tupu badala ya kunakili maudhui ya kikundi | imezimwa |
| `--prefer <group\|mine\|merge>` | jinsi ya kuamua migongano | huuliza; kwa `--yes`, chaguo-msingi la kila kipengele |
| `-y`, `--yes` | huendesha bila kuuliza; migongano hufuata `--prefer` | imezimwa |
| `-h`, `--help` | msaada wa `share` | |
| `--dry-run` (jumla) | inaheshimiwa: huuliza kuhusu migongano, huchapisha mpango, haibadilishi chochote | imezimwa |

Thamani `on` inamaanisha kikundi `default`, `off` inamaanisha ya ndani (haishirikiwi), na thamani nyingine yoyote ni jina la kikundi. Majina ya vikundi huwa majina ya folda chini ya `~/.local/share/ccam/groups/`, kwa hiyo `ccam share side --history ../x` hushindwa na `invalid group name: "../x"`. Akaunti hutajwa kwa jina au kwa namba, na vipengele unavyovitaja pekee ndivyo hubadilika.

## Jinsi inavyofanya kazi

Faili zinazohusika:

- `~/.config/ccam/config.toml`: jedwali la `[profile.share]` la kila akaunti, thamani moja kwa kila kipengele.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` ya akaunti.
- `~/.local/share/ccam/groups/<group>/`: nakala zinazoshirikiwa.

1. `internal/cli/share.go` husoma usanidi, hutumia kila thamani ya `--<item>` kwenye nakala ya mipangilio ya akaunti kwa `config.Share.Set`, na hukataa ikiwa akaunti inaendeshwa.
2. `share.Build` hulinganisha mipangilio ya sasa na lengo na huorodhesha hatua. Husoma diski tu: `os.Lstat` na `os.Readlink` kwa kila njia katika `share.Layout`. Mpango ni orodha ya operesheni: kuhamisha kwenda nafasi tupu ya kikundi (adopt), kuhamisha faili ambazo kikundi hakina (merge dir), kuandika mgongano uliotatuliwa kwenye kikundi, kuhifadhi nakala rudufu, kuondoa kiungo, kuunganisha, kunakili kurudi, kuandika upya njia za plugins na kuunganisha rejista ya plugins.
3. Migongano huamuliwa, kwa swali au kwa `--prefer`.
4. `Plan.Describe` huchapisha mpango, na ccam huuliza.
5. Mipangilio lengwa huandikwa kwenye `config.toml` kabla ya kitu chochote kuhamishwa. Ikiwa `Plan.Apply` itasimama katikati, `ccam doctor --fix` na kuendesha tena huunda mpango wao kuelekea lengo lililohifadhiwa na hukamilisha badiliko.
6. `Plan.Apply` huendesha hatua kwa mpangilio. Kila hatua huhamisha data kwenda kwenye kikundi au kwenye folda ya nakala rudufu; hakuna kinachofutwa.

`share.Layout` ni orodha ya ruhusa ya njia ambazo kila kipengele hujumuisha. Kitu chochote kisichokuwamo kamwe hakiwezi kufika kwenye kikundi. `share.Never` hutaja `.credentials.json` na `.claude.json`, na `Build` husimama na kosa la ndani ikiwa yoyote kati ya hizo itaonekana kwenye mpangilio. Kipengele cha `mcp` hakina njia: orodha ya seva iko ndani ya `.claude.json`, kwa hiyo [kushiriki MCP](/features/mcp) hufanya kazi kupitia `mcp.json` ya kikundi badala yake, na `--mcp` hubadilisha mpangilio tu.

Hiki ndicho kinachotokea kwa folda moja, `projects/`, wakati `side` inawasha historia na kikundi tayari kina maudhui:

<!-- motion:share-on -->

Ikiwa kikundi bado hakina chochote kwenye njia hiyo, akaunti ndiyo ya kwanza kuishiriki na folda yake huhamishwa huko kama ilivyo, bila nakala rudufu. Vinginevyo `mergeDir` huhamisha kila faili ambayo kikundi hakina kwenda kwenye kikundi, na faili zilizopo pande zote mbili hubaki katika folda ya akaunti, kwa hiyo toleo la kikundi huhifadhiwa. Kisha folda ya akaunti huenda `.ccam-backup/<time>/` ndani ya folda ya akaunti, na `os.Symlink` huweka kiungo mahali pake. Faili za mazungumzo hupewa majina kwa kufuata session ID yake, kwa hiyo historia kutoka akaunti kadhaa huungana bila migongano.

Kuhamisha hutumia `rename(2)` na hurudi kwenye kunakili na kufuta kati ya mifumo ya faili tofauti. Kiungo cha jamaa ndani ya mti uliohamishwa kilichoelekea nje yake huandikwa upya kama kiungo kamili (`fsutil.PortableLink`), ili bado kielekee mahali sahihi kutoka mahali papya. Rejista za `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) huhifadhi njia kamili kupitia folda ya akaunti; ccam huzielekeza kwenye folda ya kikundi na, kikundi kikiwa tayari na rejista, huunganisha vipengele vya akaunti ndani yake.

Kuzima kipengele hugeuza kiungo:

<!-- motion:share-off -->

ccam huondoa kiungo tu ikiwa kinaelekea ndani ya folda ya vikundi, kisha hunakili maudhui ya sasa ya kikundi kwenda kwenye akaunti kwa `CopyTree`. Kwa `--fresh` kunakili kunarukwa na akaunti huanza bila chochote kwenye njia hiyo. Kikundi huendelea kuhifadhi maudhui yake kwa ajili ya akaunti nyingine. Kuhamisha kipengele kutoka kikundi kimoja kwenda kingine huondoa kiungo na hukiunganisha na kikundi kipya, kikiacha kikundi cha zamani kama kilivyo; mpango husema hivyo katika mstari wa `Note:`.

Faili moja iliyopo pande zote mbili yenye maudhui tofauti ni mgongano:

<!-- motion:share-conflict -->

`share.NewConflict` huorodhesha chaguo zinazofaa faili hiyo. `history.jsonl` hutoa `merge by line`, `keep group` na `use mine`, uunganishaji ukiwa wa kwanza. Faili nyingine za `.json` hutoa `keep group`, `use mine` na `merge JSON`. Kitu kingine chochote hutoa `keep group` na `use mine`. Chaguo la kwanza ndilo chaguo-msingi:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) huunganisha vitu kwa kujirudia, huunganisha safu bila marudio, na huhifadhi thamani ya kikundi wakati pande zote mbili zinaweka ufunguo ule ule kwa njia tofauti, kwa hiyo akaunti nyingine hazioni badiliko la tabia. `merge by line` huhifadhi mistari ya kikundi, kisha huongeza mistari yako ambayo haina. Kila chaguo isipokuwa `keep group` kwanza huhifadhi nakala ya kikundi katika `.ccam-backup/<time>/group/` ndani ya folda ya akaunti, kisha huandika matokeo kwa `fsutil.AtomicWrite`. Faili yako mwenyewe daima huenda kwenye nakala rudufu.

`--prefer group`, `mine` au `merge` huamua kila mgongano bila kuuliza; `merge` inamaanisha `merge JSON` au `merge by line`, yoyote ambayo faili inatoa, na faili isiyotoa yoyote kati ya hizo huhifadhi chaguo-msingi. `--yes` peke yake huchukua chaguo-msingi la kila faili.

Katika kila uanzishaji, `share.Check` hulinganisha kila njia katika mpangilio na mpangilio wa kushiriki kwa kutumia `lstat` na `readlink` tu, jambo linaloifanya kuwa nafuu vya kutosha kwa njia ya kuanzisha. Kiungo kinachokosekana huundwa papo hapo, kwa sababu hakuna cha kupoteza. Faili halisi mahali ambapo kiungo kinapaswa kuwa, kiungo kinachoelekea mahali pasipo sahihi, au kiungo kilichopitwa na wakati husimamisha uanzishaji:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` huendesha `share.Build` kwa mipangilio iliyohifadhiwa pande zote mbili, jambo linalotoa kila hatua kutokana na kilichopo kwenye diski, huhifadhi maudhui ya kikundi kwenye migongano na huhifadhi nakala rudufu ya chochote inachokibadilisha.

## Usalama na hali maalum

- **Akaunti zinazoendeshwa.** ccam hukataa kubadilisha kushiriki wakati Claude Code inaendeshwa kama akaunti hiyo (`account work is running; quit it first, then change its sharing`), kwa sababu kipindi kinaendelea kuandika kwenye faili zinazohamishwa na mistari inayoongeza wakati huo ingepotea. `--dry-run` bado huchapisha mpango.
- **Nakala rudufu.** Kila mpango hupata folda yake ya `.ccam-backup/<time>/` katika folda ya akaunti, yenye kiambishi cha namba mipango miwili inapoanza katika sekunde ile ile. ccam kamwe haizifuti.
- **Kushindwa kwa sehemu.** Hatua ikishindwa, ccam huchapisha kosa na mahali faili asili zilipo, na hukuambia ukamilishe kwa `ccam doctor --fix`. Lengo tayari limehifadhiwa, kwa hiyo urekebishaji huendelea kuelekea huko.
- **Viungo vyako mwenyewe.** Ikiwa njia tayari ni kiungo ulichokiunda, kwa mfano kuelekea hazina ya dotfiles, na kikundi kiko tupu, ccam hunakili lengo la kiungo kwenda kwenye kikundi na huhamisha kiungo chenyewe kwenye nakala rudufu. Lengo haliguswi.
- **Faili za kuingia.** `.credentials.json` na `.claude.json` kamwe haziunganishwi, kuhamishwa wala kunakiliwa kwenye kikundi.
- **Akaunti mbili katika folda moja.** Akaunti zinazoshiriki historia huandika kwenye faili zile zile. Akaunti ya pili ya kikundi kile kile cha historia inapoanza katika folda ambamo moja tayari inaendeshwa, ccam huonya na huuliza kabla ya kuendelea.

## Zinazohusiana

- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Seva za MCP zinazoshirikiwa](/features/mcp)
- [Nakili usanidi uliopo](/features/migrate)
- [Amri zote](/docs/commands)
