---
title: Kukagua na kurekebisha
description: Kagua viungo vya kushiriki, ruhusa za faili, kuingia na toleo la Claude Code kwa ccam doctor, na urekebishe kilicho salama kurekebisha kwa ccam doctor --fix.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` hukagua usakinishaji wako wa ccam: binary ya Claude Code, folda ya kila akaunti, taarifa za kuingia na viungo vya kushiriki, seva za MCP zinazoshirikiwa, endpoint ya matumizi na amri zako za njia fupi. `ccam doctor --fix` hurekebisha kile inachoweza bila kupoteza data, na kwanza huhamisha faili yoyote inayoibadilisha kwenda kwenye folda ya nakala rudufu. Iendeshe `ccam <name>` inaposimama kwa sababu ya tatizo la kiungo cha kushiriki, baada ya kurejesha faili kwa mkono, au baada ya kusasisha Claude Code.

## Matumizi

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

Mstari wa kwanza huweka alama kwa kila akaunti kwa pamoja, ili uone matatizo yalipo kabla ya kusoma kitu kingine chochote. Chini yake, kila akaunti hupata kisanduku, na matokeo yasiyo ya akaunti yoyote (Claude Code, vigezo vya mazingira, plugins, seva za MCP, njia fupi, kanuni za folda, vikundi na mipangilio, `auto_switch`) huenda kwenye kisanduku cha `General`. Kichwa cha kisanduku huonyesha alama mbaya zaidi ya akaunti na idadi ya matatizo yake. Kila safu ndani yake ni eneo moja lililokaguliwa, kila mara `Folder`, `Login` na `Share links` kwa mpangilio huo, kisha mengine yoyote, yakiwa na alama: `✔` sawa, `ℹ` taarifa, `⚠` onyo, `✘` hitilafu. Folda isiyo na tatizo lolote huonyesha `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Kwa `--fix`, kila matokeo yaliyorekebishwa huishia na `→ fixed`, na hayahesabiwi tena kama tatizo:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Kwenye terminal visanduku hulingana na upana wake na ujumbe mrefu hufungwa kwenye mstari unaofuata, kwa rangi za terminal yenyewe (`NO_COLOR` huzizima). Kisanduku kamwe hakiwi chembamba kuliko safu wima 29, kwa hiyo hakuna neno linalokatwa kwenye terminal nyembamba. Yakitumwa kwenye pipe au faili, matokeo hayana rangi na hakuna ujumbe unaofungwa wala kukatwa, kwa hiyo `ccam doctor | grep settings.json` hupata mstari mzima; kisha kila kisanduku huwa na upana wa mstari mrefu zaidi. Kichupo cha Health cha dashibodi huonyesha visanduku vile vile.

## Chaguo

| Chaguo | Kinachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--fix` | Rekebisha kilicho salama kurekebisha (faili zilizobadilishwa huhifadhiwa nakala rudufu) | imezimwa |
| `--dry-run` (global) | Endesha kila ukaguzi lakini usirekebishe chochote, hata kwa `--fix` | imezimwa |

`ccam doctor` haipokei hoja. Hukubali `--dry-run`: `ccam --dry-run doctor --fix` huchapisha matokeo ya `ccam doctor`, huongeza kile ambacho `--fix` ingeweka kando, na haibadilishi kitu:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Exit code ni `1` wakati angalau tokeo moja la `✘` limebaki bila kurekebishwa, na `0` vinginevyo. Maonyo na taarifa kamwe hayasababishi kushindwa. Uendeshaji unaposhindwa bila `--fix`, ccam humalizia kwa mstari mtupu na `` `ccam doctor --fix` repairs what it can. ``. Usanidi usioweza kupakiwa pia humalizika kwa `1` na hitilafu ya `ccam:`.

## Jinsi inavyofanya kazi

`internal/cli/doctor.go` hupakia `~/.config/ccam/config.toml` na huita `doctor.Run` katika `internal/doctor/doctor.go`, ambayo huendesha ukaguzi kwa mpangilio maalum na hurudisha matokeo, kila moja likiwa na akaunti linalohusika nayo na eneo lililokagua. `tui.HealthReport` huyapanga kwenye visanduku. `--fix` hupitishwa kama `Options.Fix`, ambayo huwekwa false katika dry run.

<!-- motion:doctor-checks -->

Sehemu ya kati huendeshwa mara moja kwa kila akaunti. Ukaguzi, na kile `--fix` hufanya kwa kila mmoja:

1. **Toleo la Claude Code** (`checkClaude`). Hutafuta binary ya `claude` kama launcher inavyofanya na huendesha `claude --version`. Binary inayokosekana, `claude_bin` ambayo si faili inayotekelezeka, au `--version` inayoshindwa ni hitilafu. Hakuna urekebishaji.
2. **Taarifa za kuingia kwenye mazingira** (`checkAuthEnv`). Huonya wakati `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` au `ANTHROPIC_AUTH_TOKEN` imewekwa: kila akaunti inayoanzishwa kutoka shell hiyo huirithi, na Claude Code huipendelea kuliko taarifa za kuingia za akaunti yenyewe. Hakuna urekebishaji; iondoe mwenyewe.
3. **Folda** (`checkProfile`). Folda ya akaunti inayokosekana ni hitilafu na humaliza ukaguzi wa akaunti hiyo. Ruhusa zenye biti zozote za group au other ni onyo; `--fix` huendesha `chmod 0700`. `~/.claude` iliyoingizwa yenye `.claude.json` yake yenyewe kando ya `~/.claude.json` ni taarifa: hapo `ccam main` na `claude` ya kawaida huhifadhi taarifa za kuingia na seva za MCP tofauti. Windows haina biti za ruhusa, kwa hiyo ukaguzi huu na ukaguzi wa `0600` ulio hapa chini hurukwa huko; folda iliyo katika profaili yako ya mtumiaji inalindwa na orodha yake ya ufikiaji.
4. **Faili ya kuingia** (`checkCredentials`). `.credentials.json` hukaguliwa kwa `lstat`. Symlink ni hitilafu, kwa sababu kila akaunti lazima iwe na tokeni yake yenyewe, na hairekebishwi. Mode isiyo `0600` ni onyo; `--fix` huendesha `chmod 0600`. Kisha taarifa za kuingia husomwa (kwenye macOS kutoka Keychain kwanza). Kukosa taarifa za kuingia, au faili ya kuingia isiyo nazo, ni onyo linaloelekeza kwa `ccam login <name>`. Refresh token iliyoisha muda, au taarifa za kuingia ambazo Claude Code ilizifuta baada ya seva kuzikataa, ni hitilafu `` login expired → `ccam login <name>` ``; access token iliyoisha muda ni taarifa tu, kwa kuwa Claude Code huifanya upya wakati wa kuanza unaofuata. Taarifa za kuingia zilizohifadhiwa katika Credential Manager ya Windows ni taarifa, kwa kuwa ccam haiwezi kuzisoma. Taarifa halali za kuingia huonyesha barua pepe ambayo folda imeingia nayo sasa na mpango (`login valid (dev@company.com, Claude Team · Acme)`). Hakuna kati ya haya kinachorekebishwa; `ccam login --expired` huingiza kila akaunti inayohitaji.
5. **Viungo vya kushiriki** (`share.Check`). Hulinganisha kila njia katika `share.Layout` na vikundi vya akaunti, ikitumia `lstat` na `readlink` tu. Huripoti aina tano: kiungo kinachokosekana, kiungo kinachoelekea mahali pengine, faili halisi mahali ambapo kiungo kinapaswa kuwa, kiungo kilichobaki kwa kipengele ambacho sasa ni cha ndani, na kiungo sahihi ambacho folda yake ya kikundi inakosekana (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; uanzishaji haukagui hili). Kiungo kinachokosekana na folda ya kikundi inayokosekana ni maonyo, vingine ni hitilafu. Unakili wa kurudisha maudhui kutoka kwenye kikundi uliosimama nusu njia pia huripotiwa: `--fix` huweka mahali pake nakala kamili iliyoandaliwa, lakini ikiwa njia hiyo imeundwa upya tangu hapo (na uanzishaji), hutaja folda yenye nakala kamili ili uishughulikie mwenyewe. Njia zilizoongezwa kwenye kipengele katika toleo la baadaye la ccam (`file-history/`, `tasks/` na `paste-cache/` kwa historia, `keybindings.json` kwa mipangilio) ni onyo tu wakati akaunti bado ina nakala yake yenyewe: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Uanzishaji kamwe hausimami kwa sababu yake.
6. **Vipengee visivyojulikana** (`checkUnknownEntries`). Huorodhesha majina ya ngazi ya juu katika folda ya akaunti ambayo ccam haiyajui (`share.Known`), kama ishara kwamba mpangilio wa Claude Code umebadilika. Taarifa tu.
7. **Njia za plugins** (`checkPlugins`). Rejista za plugins (`installed_plugins.json`, `known_marketplaces.json`) huhifadhi njia kamili. Katika `plugins/` ya kila kikundi, njia zinazopita kwenye folda ya akaunti moja huvunjika akaunti hiyo inapoondolewa au kubadilishwa jina. `--fix` huziandika upya kuelekea folda ya kikundi kwa `profile.RewritePluginPaths`, uandikaji wa atomiki wa faili ya rejista.
8. **Seva za MCP** (`checkMCP`). Hupakia `mcp.json` ya kila kikundi mara moja. Faili isiyoweza kusomwa ni hitilafu; seva inayoshindwa `mcp.Validate` ni onyo. Seva zilizo katika `~/.claude.json` ya `claude` ya kawaida ambazo hakuna akaunti inayozipakia ni onyo linalotaja amri ya `ccam migrate --mcp` inayozinakili (pamoja na `--group <a|b>` wakati vikundi kadhaa vinashiriki MCP). Seva ambazo akaunti moja tu ya kikundi inazo ni taarifa pamoja na amri ya `ccam mcp import` inayozishiriki, na akaunti ambayo haijaingia kwenye seva za MCP ambazo akaunti nyingine zimeingia hupata onyo la kuingia kwa `/mcp` (ccam kamwe hainakili kuingia). Hakuna urekebishaji.
9. **Muundo wa jibu la matumizi** (`checkUsageFormat`). Huuliza endpoint ya matumizi, bila cache, kwa akaunti moja baada ya nyingine hadi moja ijibu, zote ndani ya kikomo kimoja cha sekunde 4. Kwa `usage_source = "claude"` badala yake huiuliza Claude Code kwa akaunti moja. Ikiwa jibu haliwezi tena kuchanganuliwa, huonya kwamba `ccam usage` hurudi kwenye cache ya mstari wa hali; hakuna akaunti ikijibu, husema hivyo kama taarifa (`could not check the usage response format: …`). Huu ndio ukaguzi pekee unaotumia mtandao.
10. **Njia fupi** (`checkShims`). Tu wakati akaunti ina njia fupi. Viungo vinavyokosekana katika `~/.local/bin` ni onyo, vivyo hivyo njia fupi ambazo ccam iliunda ambazo bado zinaendesha ccam ya zamani, iliyobadilishwa au kuhamishwa (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` huita `shim.Install`, ambayo huunda na kusasisha tu viungo ambavyo ccam iliunda na kurekodi katika `~/.local/state/ccam/shortcuts.json`, na kamwe haigusi faili nyingine. Onyo huishia na `→ fixed` tu wakati kila kimoja cha viungo hivyo kimeundwa; kimoja kisipoweza kuundwa, sababu huongezwa. Njia fupi inayoendesha ccam nyingine ambayo bado imesakinishwa hufanya kazi, na ni taarifa; `--fix` huielekeza kwa ccam hii. Folda ya bin isiyo kwenye `PATH` ni onyo na hairekebishwi.
11. **Folda** (`checkFolders`). Tu wakati maingizo ya [`ccam folder`](/features/folder) yapo. Ingizo ambalo akaunti yake haipo ni hitilafu, na ingizo ambalo folda yake imetoweka ni onyo. Hakuna kati yao kinachorekebishwa, kwa sababu kuondoa ingizo kungebadilisha akaunti inayoanza hapo.
12. **Vikundi na mipangilio** (`checkGroups`). Mwanachama wa kikundi ambaye si akaunti, akaunti ambayo jina lake tangu hapo limekuwa amri ya ccam (ianzishe kwa `ccam run <name>` au ibadilishe jina), `theme` au `usage_source` isiyojulikana, na mistari ya `config.toml` ambayo ccam haiisomi ni maonyo. Mistari hiyo ya mwisho kwa kawaida ni mpangilio ulioandikwa chini ya jedwali la `[[group]]`, `[[profile]]` au `[[folder]]`, ambapo huwa sehemu ya jedwali hilo (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); amri zinazobadilisha usanidi hukataa hadi hilo lirekebishwe. Hakuna urekebishaji.
13. **Kuingia** (`checkIdentities`). Husoma kila folda imeingia kama nani kutoka `.claude.json` yake. Folda ambayo sasa imeingia kwenye akaunti nyingine kuliko ile ambayo `ccam login` ilirekodi (kwa mfano baada ya `/login` ndani ya kipindi) ni onyo; akaunti isiyo na rekodi ya kuingia, kama iliyoingizwa au iliyoundwa na ccam ya zamani, hulinganishwa kwa barua pepe yake iliyorekodiwa. Folda iliyoingia kwenye akaunti ile ile ya Claude kama nyingine pia ni onyo: hizo mbili hushiriki vikomo vyao vya matumizi. Folda inayoshikilia tokeni ile ile ya kuingia kama nyingine, yaani folda iliyonakiliwa (ikiwemo `~/.claude` ya `claude` ya kawaida), ni hitilafu, kwa kuwa ile itakayoifanya upya kwanza huitoa nyingine nje ya akaunti; endesha `ccam login <name>` ili kuipa tokeni yake yenyewe. Hakuna urekebishaji.
14. **`auto_switch`** (`checkAutoSwitch`). Ni pale tu `auto_switch = true`. Taarifa katika kisanduku cha `General` husema `auto_switch` hufanya nini na huunganisha na [maelezo yake](/features/next#auto-switch). Kwa kila mpango wenye viti kati ya akaunti (Claude Team, Claude Enterprise), mstari mmoja huzitaja, husema kwamba `auto_switch` huhamisha kipindi kati ya viti vya shirika moja, hunukuu kanuni ya viti ya Anthropic pamoja na kiungo cha kituo cha usaidizi, na hukuachia uamuzi. Claude Code ya zamani kuliko 2.1.214, ambayo haiambii hooks mchakato wake (`CLAUDE_PID`), ni onyo huko, kwa kuwa hakuna kipindi kinachoweza kuhamishwa. Akaunti isiyo na plugin ya ccam, au ambayo plugin imezimwa, hupata onyo katika kisanduku chake yenyewe pamoja na amri inayolirekebisha ndani ya Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (au `/plugin enable ccam@ccam`). Akaunti za API key haziguswi. Kwenye Windows mstari pekee ni `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` Hakuna urekebishaji.

Kabla ya akaunti, **badiliko la vikundi** lililosimama nusu njia ni hitilafu katika kisanduku cha `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` hulikamilisha kwa mpangilio uliopangwa, pamoja na majibu yaliyotolewa kwa migongano yake; hadi wakati huo hakuna badiliko jingine la vikundi, kubadilisha jina wala kuondoa kunakoanza. Akaunti zilizomo huundiwa tu viungo vyake vinavyokosekana hadi likamilike. Lisipoweza kukamilishwa, hitilafu husema kwamba kuondoa faili hiyo kunamaanisha kuachana na sehemu iliyobaki ya badiliko. `ccam --dry-run doctor --fix` huonyesha mpangilio (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) na kile ambacho sehemu ya kila akaunti ingeweka kando. Badiliko ambalo ccam nyingine inalitekeleza sasa hivi ni taarifa, na wakati linaendelea `--fix` hairekebishi viungo vya kushiriki (`its share links are not repaired while another ccam changes the groups; …`). `--fix` hushikilia `~/.local/state/ccam/group-change.lock` wakati inarekebisha, na husoma `config.toml` tena baada ya kukamilisha badiliko.

### Kurekebisha viungo vya kushiriki

Njia ya kuanzisha huunda kiungo kinachokosekana papo hapo, lakini husimama kwa jambo lolote linaloweza kupoteza data na huelekeza hapa:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Kwa tatizo lolote la kiungo cha kushiriki, `--fix` huita `share.Build` ikiwa na mipangilio ya kushiriki ya akaunti kama hali ya sasa na pia hali lengwa, kisha `Plan.Apply`. `Build` hutengeneza upya kila hatua kutoka kilichopo kwenye diski, mpango ule ule ambao badiliko la `ccam group` hutumia, kwa hiyo urekebishaji hufuata sheria zile zile kama kujiunga na kikundi.

<!-- motion:doctor-repair -->

Hapa `settings.json` ni faili halisi. Kwa kuwa kikundi tayari kina `settings.json` tofauti, mgongano huchukua chaguo lake la msingi, `KeepGroup`, kwa hiyo nakala ya kikundi hubaki. Faili ya akaunti huhamia `.ccam-backup/<time>/settings.json` ndani ya folda ya akaunti, na `os.Symlink` huunganisha njia hiyo na kikundi. `--fix` husema hivyo: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Hali nyingine:

- Ikiwa kikundi bado hakina nakala, faili ya akaunti huhamia kwenye kikundi na kuwa maudhui yake ya kwanza.
- Kwa folda, faili ambazo kikundi hakina huhamia kwenye kikundi. Faili ambazo kikundi kinazo zikiwa sawa baiti kwa baiti hutupwa, na zile zinazotofautiana huenda kwenye nakala rudufu, zikitajwa katika ujumbe (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Folda ya kikundi inayokosekana huundwa kwanza.
- Kwa `history.jsonl`, chaguo-msingi ni `merge by line`: mistari ya kikundi, kisha mistari yako mipya. Nakala ya awali ya kikundi huhifadhiwa katika `.ccam-backup/<time>/group/`.
- Kiungo kinachoelekea kikundi kingine hubadilishwa. Kiungo kinachoelekea folda nje ya ccam huenda kwenye nakala rudufu, na maudhui yake huanzisha kikundi ikiwa kikundi kiko tupu.
- Kiungo kilichobaki kwa kipengele ambacho sasa ni cha ndani huondolewa, na akaunti hupata nakala ya maudhui ya kikundi.

Mbali na kuondoa viungo na kutupa faili ambazo kikundi tayari kinazo zikiwa sawa baiti kwa baiti, kila hatua huhamisha data kwenda kwenye kikundi au folda ya nakala rudufu. Hakuna kinachopotea.

## Usalama na hali maalum

- `ccam doctor` ya kawaida kamwe haibadilishi chochote. Tumia `--dry-run` pamoja na `--fix` ili kuwa na uhakika katika scripts.
- Kwa akaunti inayoendeshwa, bila kujali Claude Code ilianzishwa vipi, `doctor --fix` huunda tu viungo vinavyokosekana na folda za vikundi zinazokosekana, jambo ambalo halipotezi chochote. Kitu kingine kikibaki, huongeza `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, kwa sababu kipindi kinaendelea kuandika kwenye faili ambazo urekebishaji ungezihamisha.
- Urekebishaji wa kushiriki ukishindwa, ccam huripoti `could not repair share links: <error>` kwa akaunti hiyo na huruka ukaguzi wake wa vipengee visivyojulikana. Akaunti nyingine bado hukaguliwa.
- Ikiwa badiliko la `ccam group` linashindwa katikati, vikundi vipya tayari vimehifadhiwa na badiliko hurekodiwa katika `~/.local/state/ccam/group-change.json`; ccam hukuambia umalizie kwa `ccam doctor --fix`, ambayo hutekeleza sehemu iliyobaki, pamoja na hatua za MCP.
- `.credentials.json` iliyo symlink, taarifa za kuingia zilizoisha muda na taarifa za kuingia za mazingira huripotiwa lakini kamwe hazibadilishwi.
- Uandikaji upya wa rejista za plugins hufanywa kwa atomiki lakini haunakiliwi kwenda `.ccam-backup/`.

## Zinazohusiana

- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Kunakili usanidi uliopo](/features/migrate)
- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Amri](/docs/commands)
