---
title: Kukagua na kurekebisha
description: Kagua viungo vya kushiriki, ruhusa za faili, kuingia na toleo la Claude Code kwa ccam doctor, na urekebishe kilicho salama kurekebisha kwa ccam doctor --fix.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` hukagua usakinishaji wako wa ccam: binary ya Claude Code, folda ya kila akaunti, taarifa za kuingia na viungo vya kushiriki, seva za MCP zinazoshirikiwa, endpoint ya matumizi na amri zako za njia fupi. `ccam doctor --fix` hurekebisha kile inachoweza bila kupoteza data, na kwanza huhamisha faili yoyote inayoibadilisha kwenda kwenye folda ya nakala rudufu. Iendeshe `ccam <name>` inaposimama kwa sababu ya tatizo la kiungo cha kushiriki, baada ya kurejesha faili kwa mkono, au baada ya kusasisha Claude Code.

## Matumizi

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Kila matokeo ni mstari mmoja wenye alama: `✓` sawa, `·` taarifa, `!` onyo, `✗` hitilafu. Matokeo ya akaunti huanza na jina la akaunti.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Kwa `--fix`, kila matokeo yaliyorekebishwa huishia na `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Chaguo

| Chaguo | Kinachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--fix` | Rekebisha kilicho salama kurekebisha (faili zilizobadilishwa huhifadhiwa nakala rudufu) | imezimwa |
| `--dry-run` (global) | Endesha kila ukaguzi lakini usirekebishe chochote, hata kwa `--fix` | imezimwa |

`ccam doctor` haipokei hoja. Hukubali `--dry-run`: `ccam --dry-run doctor --fix` huchapisha matokeo yale yale kama `ccam doctor` na haibadilishi kitu.

Exit code ni `1` wakati angalau tokeo moja la `✗` limebaki bila kurekebishwa, na `0` vinginevyo. Maonyo na taarifa kamwe hayasababishi kushindwa. Uendeshaji unaposhindwa bila `--fix`, ccam humalizia kwa mstari mtupu na `` `ccam doctor --fix` repairs what it can. ``. Usanidi usioweza kupakiwa pia humalizika kwa `1` na hitilafu ya `ccam:`.

## Jinsi inavyofanya kazi

`internal/cli/doctor.go` hupakia `~/.config/ccam/config.toml` na huita `doctor.Run` katika `internal/doctor/doctor.go`, ambayo huendesha ukaguzi kwa mpangilio maalum na hurudisha matokeo. `--fix` hupitishwa kama `Options.Fix`, ambayo huwekwa false katika dry run.

<!-- motion:doctor-checks -->

Sehemu ya kati huendeshwa mara moja kwa kila akaunti. Ukaguzi, na kile `--fix` hufanya kwa kila mmoja:

1. **Toleo la Claude Code** (`checkClaude`). Hutafuta binary ya `claude` kama launcher inavyofanya na huendesha `claude --version`. Ikiwa toleo halianzi na `2.1.302`, toleo ambalo mpangilio wa folda wa ccam ulithibitishwa dhidi yake (`doctor.LayoutVersion`), huongeza `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Binary inayokosekana au `--version` inayoshindwa ni hitilafu. Hakuna urekebishaji.
2. **Taarifa za kuingia kwenye mazingira** (`checkAuthEnv`). Huonya wakati `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` au `ANTHROPIC_AUTH_TOKEN` imewekwa: kila akaunti inayoanzishwa kutoka shell hiyo huirithi, na Claude Code huipendelea kuliko taarifa za kuingia za akaunti yenyewe. Hakuna urekebishaji; iondoe mwenyewe.
3. **Folda** (`checkProfile`). Folda ya akaunti inayokosekana ni hitilafu na humaliza ukaguzi wa akaunti hiyo. Ruhusa zenye biti zozote za group au other ni onyo; `--fix` huendesha `chmod 0700`.
4. **Faili ya kuingia** (`checkCredentials`). `.credentials.json` hukaguliwa kwa `lstat`. Symlink ni hitilafu, kwa sababu kila akaunti lazima iwe na tokeni yake yenyewe, na hairekebishwi. Mode isiyo `0600` ni onyo; `--fix` huendesha `chmod 0600`. Kisha taarifa za kuingia husomwa (kwenye macOS kutoka Keychain kwanza). Kukosa taarifa za kuingia ni onyo linaloelekeza kwa `ccam login <name>`. Refresh token iliyoisha muda ni hitilafu `login expired`; access token iliyoisha muda ni taarifa tu, kwa kuwa Claude Code huifanya upya wakati wa kuanza unaofuata. Hakuna kati ya haya kinachorekebishwa.
5. **Viungo vya kushiriki** (`share.Check`). Hulinganisha kila njia katika `share.Layout` na mipangilio ya kushiriki ya akaunti, ikitumia `lstat` na `readlink` tu. Huripoti aina nne: kiungo kinachokosekana, kiungo kinachoelekea mahali pengine, faili halisi mahali ambapo kiungo kinapaswa kuwa, na kiungo kilichobaki kwa kipengele ambacho sasa ni cha ndani. Kiungo kinachokosekana ni onyo, vingine ni hitilafu.
6. **Vipengee visivyojulikana** (`checkUnknownEntries`). Huorodhesha majina ya ngazi ya juu katika folda ya akaunti ambayo ccam haiyajui (`share.Known`), kama ishara kwamba mpangilio wa Claude Code umebadilika. Taarifa tu.
7. **Njia za plugins** (`checkPlugins`). Rejista za plugins (`installed_plugins.json`, `known_marketplaces.json`) huhifadhi njia kamili. Katika `plugins/` ya kila kikundi, njia zinazopita kwenye folda ya akaunti moja huvunjika akaunti hiyo inapoondolewa au kubadilishwa jina. `--fix` huziandika upya kuelekea folda ya kikundi kwa `profile.RewritePluginPaths`, uandikaji wa atomiki wa faili ya rejista.
8. **Seva za MCP** (`checkMCP`). Hupakia `mcp.json` ya kila kikundi mara moja. Faili isiyoweza kusomwa ni hitilafu; seva inayoshindwa `mcp.Validate` ni onyo. Hakuna urekebishaji.
9. **Muundo wa jibu la matumizi** (`checkUsageFormat`). Huuliza endpoint ya matumizi, bila cache na kwa kikomo cha kawaida cha sekunde 3, kwa akaunti moja baada ya nyingine hadi moja ijibu. Ikiwa jibu haliwezi tena kuchanganuliwa, huonya kwamba `ccam usage` hurudi kwenye cache ya mstari wa hali. Huu ndio ukaguzi pekee unaotumia mtandao.
10. **Njia fupi** (`checkShims`). Tu wakati akaunti ina njia fupi. Viungo vinavyokosekana katika `~/.local/bin` ni onyo; `--fix` huita `shim.Install`, ambayo huunda tu viungo vinavyokosekana na kamwe haigusi faili nyingine. Folda ya bin isiyo kwenye `PATH` ni onyo na hairekebishwi.

### Kurekebisha viungo vya kushiriki

Njia ya kuanzisha huunda kiungo kinachokosekana papo hapo, lakini husimama kwa jambo lolote linaloweza kupoteza data na huelekeza hapa:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Kwa tatizo lolote la kiungo cha kushiriki, `--fix` huita `share.Build` ikiwa na mipangilio ya kushiriki ya akaunti kama hali ya sasa na pia hali lengwa, kisha `Plan.Apply`. `Build` hutengeneza upya kila hatua kutoka kilichopo kwenye diski, mpango ule ule ambao `ccam share` hutumia, kwa hiyo urekebishaji hufuata sheria zile zile kama kuwasha kushiriki.

<!-- motion:doctor-repair -->

Hapa `settings.json` ni faili halisi. Kwa kuwa kikundi tayari kina `settings.json` tofauti, mgongano huchukua chaguo lake la msingi, `KeepGroup`, kwa hiyo nakala ya kikundi hubaki. Faili ya akaunti huhamia `.ccam-backup/<time>/settings.json` ndani ya folda ya akaunti, na `os.Symlink` huunganisha njia hiyo na kikundi. Hali nyingine:

- Ikiwa kikundi bado hakina nakala, faili ya akaunti huhamia kwenye kikundi na kuwa maudhui yake ya kwanza.
- Kwa folda, faili ambazo kikundi hakina huhamia kwenye kikundi; zilizobaki huenda kwenye nakala rudufu.
- Kwa `history.jsonl`, chaguo-msingi ni `merge by line`: mistari ya kikundi, kisha mistari yako mipya. Nakala ya awali ya kikundi huhifadhiwa katika `.ccam-backup/<time>/group/`.
- Kiungo kinachoelekea kikundi kingine hubadilishwa. Kiungo kinachoelekea folda nje ya ccam huenda kwenye nakala rudufu, na maudhui yake huanzisha kikundi ikiwa kikundi kiko tupu.
- Kiungo kilichobaki kwa kipengele ambacho sasa ni cha ndani huondolewa, na akaunti hupata nakala ya maudhui ya kikundi.

Mbali na kuondoa viungo, kila hatua huhamisha data kwenda kwenye kikundi au folda ya nakala rudufu. Hakuna kinachofutwa.

## Usalama na hali maalum

- `ccam doctor` ya kawaida kamwe haibadilishi chochote. Tumia `--dry-run` pamoja na `--fix` ili kuwa na uhakika katika scripts.
- `doctor --fix` haikagui kama akaunti inaendeshwa. Funga Claude Code ya akaunti kabla ya kurekebisha viungo vyake vya kushiriki, kwa sababu kipindi kinaendelea kuandika kwenye faili hizo.
- Urekebishaji wa kushiriki ukishindwa, ccam huripoti `could not repair share links: <error>` kwa akaunti hiyo na huruka ukaguzi wake wa vipengee visivyojulikana. Akaunti nyingine bado hukaguliwa.
- Ikiwa `ccam share` inashindwa katikati, kwanza huhifadhi mipangilio lengwa na hukuambia umalizie kwa `ccam doctor --fix`.
- `.credentials.json` iliyo symlink, taarifa za kuingia zilizoisha muda na taarifa za kuingia za mazingira huripotiwa lakini kamwe hazibadilishwi.
- Uandikaji upya wa rejista za plugins hufanywa kwa atomiki lakini haunakiliwi kwenda `.ccam-backup/`.

## Zinazohusiana

- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Kunakili usanidi uliopo](/features/migrate)
- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Amri](/docs/commands)
