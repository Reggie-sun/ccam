---
title: Vikomo vya matumizi
description: ccam usage huonyesha vikomo vya saa 5 na vya wiki vya kila akaunti kwa pamoja, na ccam statusline huweka jina la akaunti na matumizi yake katika Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` huonyesha kiasi cha kikomo cha saa 5, kikomo cha wiki na kikomo cha wiki cha kila model (kama Fable) ambacho kila akaunti imetumia, kama upau mmoja kwa kila kikomo pamoja na muda hadi kianze upya. Itumie kabla ya kuanza kazi ndefu, au Claude Code inapokuambia kikomo kiko karibu, ili kuona ni akaunti ipi bado ina nafasi. `ccam statusline` ni amri ya mstari wa hali (status line) kwa Claude Code: huonyesha jina la akaunti na matumizi ndani ya Claude Code na huhifadhi namba kwa ajili ya `ccam usage`.

## Matumizi

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Kila akaunti hupata mstari wenye namba yake, jina lake na umri wa namba zake; `(cached)` inamaanisha zinatoka kwenye faili iliyohifadhiwa, si kwenye endpoint. Chini yake, mstari mmoja kwa kila kikomo huonyesha upau wa visanduku 20, asilimia na muda hadi kipindi kianze upya: `5-hour`, `Weekly` (model zote), kisha kipindi cha wiki cha kila model yenyewe, kama `Fable`. Kipindi ambacho akaunti haina huonyesha `-`. Asilimia hugeuka manjano kuanzia 80% na nyekundu kwa 100%. Kwenye terminal kila akaunti huonekana mara moja kama `fetching…`, na mistari yake huchorwa inapojibu. Mstari wa mwisho hutaja akaunti ambayo [`ccam next --auto`](/features/next) ingechagua kwa kanuni zile zile (`auto_max_usage`, `auto_strategy`, vipindi vinavyoendeshwa, akaunti za shirika la akaunti iliyotumika mwisho pekee isipokuwa `cross_org` imewekwa, kila kipindi kilicho 100% kikirukwa, kwa hiyo akaunti ambayo wiki yake imejaa haitajwi kamwe), na huonekana tu ukiwa na zaidi ya akaunti moja. Akaunti ambayo haijaingia huonyesha `not logged in → ccam login <name>`, na akaunti ya API key huonyesha `no limits (API)`.

Baada ya `ccam statusline install`, Claude Code huonyesha mstari kama `[work] · Opus · 5h 24% · 7d 41%`.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--json` | Huchapisha ripoti kama JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | imezimwa |
| `--watch[=N]` | Huchora upya kila sekunde `N` hadi ubonyeze Ctrl+C. `--watch` peke yake, au thamani ya 0 au chini, inamaanisha 60 | imezimwa |
| `-v`, `--verbose` | Huongeza mstari wa `source:` chini ya kila akaunti: namba zake zilikotoka (endpoint ya matumizi, Claude Code, mstari wa hali au faili iliyohifadhiwa) | imezimwa |
| `--no-cache` | Hupuuza cache ya sekunde 60 na huuliza endpoint tena | imezimwa |
| `--source live\|claude\|cache` | `live` huuliza endpoint kwa tokeni ya akaunti; `claude` huiachia Claude Code itafute (hapa chini); `cache` huonyesha namba zilizohifadhiwa tu | `usage_source`, vinginevyo `live` |
| `--refresh` | Hufanya upya tokeni za ufikiaji zilizoisha muda, kwa akaunti ambazo haziendeshwi tu | imezimwa |
| `--dry-run` | Inakubaliwa: `ccam usage` haibadilishi chochote isipokuwa cache yake yenyewe. Hukataliwa pamoja na `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | imezimwa |

Thamani nyingine yoyote ya `--source` hushindwa na `--source must be live, claude or cache`.

### Matumizi kutoka Claude Code yenyewe

Kwa `usage_source = "claude"` katika `config.toml` (au `--source claude`), ccam haisomi tokeni ya kuingia hata kidogo. Kwa kila akaunti huanzisha `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` katika folda ya akaunti na hutuma ombi la `get_usage` la Claude Code kupitia standard input (Claude Code 2.1.270 au baadaye); Claude Code hutafuta matumizi kwa taarifa zake yenyewe za kuingia na hujibu kwa vipindi vya saa 5, vya wiki na vya model. `--safe-mode` huzuia hooks, seva za MCP na plugins za akaunti zisianze, na hakuna kipindi kinachohifadhiwa. Kila akaunti huhitaji Claude Code kuanzishwa, takriban sekunde moja, kwa hiyo chanzo hiki ni cha polepole kuliko `live`; akaunti huulizwa kwa wakati mmoja, kila moja ndani ya sekunde 20, na cache ya sekunde 60 hutumika kama kawaida. Claude Code huliweka ombi hili alama ya majaribio, kwa hiyo likiacha kujibu, ccam huonyesha namba zilizohifadhiwa na `ccam doctor` huliripoti. `--dry-run` haianzishi Claude Code (huandika kwenye folda ya akaunti) na huonyesha namba zilizohifadhiwa. `usage_source` hupokea `live` au `claude`; `cache` ni chaguo la `--source` pekee. `ccam usage` haipokei hoja.

`ccam statusline` husoma JSON ya mstari wa hali ya Claude Code kutoka standard input. Chaguo lake moja ni `--wrap <command>`: amri ya mstari wa hali iliyopo ambayo mstari wake wa kwanza wa matokeo huwekewa jina la akaunti mbele. `ccam statusline install` haipokei chaguo. Zote mbili hukataa `--dry-run` kwa `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Jinsi inavyofanya kazi

Namba zinatoka sehemu mbili: endpoint ya matumizi ya Anthropic isiyo na nyaraka, `GET https://api.anthropic.com/api/oauth/usage`, inayoitwa kwa taarifa za kuingia za kila akaunti yenyewe, na thamani ambazo Claude Code hukabidhi amri yake ya mstari wa hali. Yote haya yako katika `internal/usage`.

1. `cli.usageCmd` hupakia `config.toml` na huunda `usage.Account` moja kwa kila akaunti: folda yake `~/.local/share/ccam/profiles/<name>/`, folda yake ya vipindi (kuona kama inaendeshwa) na kama iliingizwa mahali ilipo.
2. `Fetcher.FetchAll` huanzisha goroutine moja kwa kila akaunti. Zote hushiriki `http.Client` moja, kwa hiyo miunganisho kwenye endpoint hutumiwa tena. Kuunganisha kunaweza kuchukua sekunde 2, na kila akaunti hupata sekunde 3 kwa jumla. Client haifuati uelekezaji, kwa kuwa uelekezaji ungetuma tokeni mahali pengine, na husoma kwa kiwango cha juu MB 1 ya jibu.
3. Kila akaunti inapomaliza, safu yake huchorwa upya. Akaunti inayoshindwa au isiyojibu kwa wakati hurudi kwenye namba zilizohifadhiwa, kwa hiyo akaunti moja ya polepole kamwe haizuii nyingine.

Kila akaunti huendeshwa katika goroutine yake dhidi ya endpoint; ile inayopitwa na muda huonyesha namba zake zilizohifadhiwa badala yake.

<!-- motion:usage-fetch -->

`personal` na `work` hujibu kwa wakati na huonyesha namba za moja kwa moja. `side` huchukua zaidi ya sekunde 3, kwa hiyo `fallback()` husoma faili zake mbili za cache na hutumia ile iliyoandikwa mwisho. `usage.Table()` huchora kizuizi kimoja kwa kila akaunti na `usage.Best()` huchagua akaunti inayotajwa kwenye mstari wa mwisho.

Kwa akaunti moja, `Fetcher.Fetch` hufanya ukaguzi huu kwa mpangilio:

<!-- motion:usage-cache -->

1. Kwa `--source cache`, huenda moja kwa moja kwa `fallback()`.
2. Isipokuwa `--no-cache` imewekwa, matokeo ya moja kwa moja katika `~/.local/state/ccam/usage/<name>.json` yenye umri wa chini ya sekunde 60 (`CacheTTL`) hurudishwa kama yalivyo. Kuendesha `ccam usage` mara mbili mfululizo hugharimu ombi moja kwa kila akaunti.
3. `readCredentials()` husoma `.credentials.json` (au Keychain ya macOS) bila kuibadilisha, na ufunguo wa `claudeAiOauth` pekee ndio hufasiriwa. Kukosekana kwa kuingia kwa usajili kunamaanisha `not logged in`, au `no limits (API)` wakati `.claude.json` ina API key.
4. Tokeni ya ufikiaji iliyoisha muda haitumiwi. Bila `--refresh`, akaunti huonyesha namba zake zilizohifadhiwa zikiwa na alama ya kuisha muda.
5. Vinginevyo ccam huita endpoint. 401 au 403 huhesabiwa kama kuingia kulikoisha muda, hali nyingine yoyote au mwili wa jibu usioweza kusomwa huhesabiwa kama kosa, na zote mbili hurudi kwenye namba zilizohifadhiwa. 429 inamaanisha endpoint inaulizwa mara nyingi mno: ccam huandika `<name>.pause` pamoja na muda wa kuuliza tena (`Retry-After`, au dakika 5, si zaidi ya saa moja), huonyesha namba zilizohifadhiwa pamoja na `rate limited by the usage endpoint; asking again after 15:10`, na haiiti endpoint kwa akaunti hiyo hadi wakati huo, hata kwa `--no-cache`. Kujaribu tena mara moja kungeendeleza kikomo tu.
6. Jibu zuri huandikwa kwenye `<name>.json` pamoja na muda lilipopimwa.

`fallback()` husoma `<name>.json` na `<name>.statusline.json` zote mbili na hubakiza ile mpya zaidi. Kabla ya kuonyesha namba zilizohifadhiwa, hukagua muda wa kuanza upya wa kila kipindi: kipindi kilichoanza upya tangu namba zihifadhiwe huonyesha 0%, kwa sababu asilimia ya zamani haitumiki tena. Ndiyo sababu `personal` hapo juu huonyesha `0%` bila muda wa kuanza upya.

Endpoint ina miundo miwili ya majibu. Mpya zaidi una safu ya `limits` yenye vipengele vya `session` (saa 5), `weekly_all` na `weekly_scoped` (kwa kila model). Wa zamani una vitu vya `five_hour`, `seven_day`, `seven_day_opus` na `seven_day_sonnet`. `ParseLive` husoma yote miwili: safu ya `limits` hushinda, na vitu vya zamani hujaza kile ambacho haikihusishi. Mwili wa jibu usio na vikomo kabisa huripotiwa kama muundo uliobadilika, jambo ambalo `ccam doctor` hukujulisha.

### Kufanya upya tokeni

Refresh tokens huzunguka: kufanya upya moja hukupa refresh token mpya na hubatilisha ile ya zamani. Ikiwa ccam ingefanya upya kuingia wakati Claude Code inaendeshwa kama akaunti ile ile, kipindi hicho kingetolewa nje ya akaunti. Kwa hiyo `ccam usage` husoma tokeni na kamwe haiziandiki, isipokuwa utoe `--refresh`. Kisha `Fetcher.refresh`:

1. hukataa akaunti zilizosajiliwa mahali zilipo kwa `ccam import`, kwa kuwa `claude` ya kawaida inaweza kutumia kuingia huko bila ccam kujua;
2. hukataa taarifa za kuingia zilizohifadhiwa katika Keychain ya macOS;
3. huchukua faili ya kufuli, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (chini ya `~/.local/state/ccam/locks/` kigezo hicho kisipowekwa), ili michakato miwili ya ccam kamwe isifanye upya tokeni ile ile mara mbili;
4. chini ya kufuli, hukataa ikiwa kipindi cha akaunti kinaendeshwa, na huruka ombi ikiwa mchakato mwingine tayari umefanya upya tokeni;
5. hutuma refresh token kwenye endpoint ya tokeni na huandika tokeni mpya kurudi kwenye `.credentials.json` kwa `fsutil.AtomicWrite`, ikihifadhi kila ufunguo mwingine katika faili.

### Mkusanyaji wa mstari wa hali

Claude Code huendesha amri ya mstari wa hali na huipitishia hati ya JSON. ccam husoma jina la model na `rate_limits` kutoka humo:

<!-- motion:usage-statusline -->

1. `runStatusline` hupata akaunti kutoka `CCAM_PROFILE`, ambayo ccam huiweka inapoanzisha Claude Code, au vinginevyo kwa kulinganisha `CLAUDE_CONFIG_DIR` na folda za akaunti.
2. `ParseStatusLine` husoma `rate_limits.five_hour` na `rate_limits.seven_day`. JSON haina vikomo vya maombi kwa akaunti za API key, au kabla ya jibu la kwanza la kipindi; basi hakuna kinachohifadhiwa.
3. `SaveStatusLine` huandika `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` huitumia wakati endpoint haijibu.
4. ccam huchapisha `[name] · model · 5h N% · 7d N%`. Kwa `--wrap`, huendesha amri yako kupitia `sh -c` kwa ingizo lile lile na huchapisha `[name]` ikifuatiwa na mstari wa kwanza wa amri badala yake. Kwenye Windows huendesha amri kama Claude Code inavyofanya: kwa Git Bash ikiwa imesakinishwa, la sivyo kwa PowerShell.
5. Matumizi ya saa 5 au ya wiki yanapokuwa 90% au zaidi na kuna akaunti nyingine, mstari huishia na `→ ccam next --auto --continue`, kwa mfano `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Amri iliyozungukwa hupata sekunde 2. Baada ya hapo husimamishwa kwa nguvu na mstari huonyesha jina la akaunti peke yake, kwa hiyo amri ya polepole haiwezi kuzuia mstari wa hali wa Claude Code.

`ccam statusline install` huandika `statusLine` katika `settings.json` ya kila akaunti. Akaunti ikiwa tayari ina amri ya mstari wa hali, ccam huihifadhi kama `--wrap '<old command>'`. Kwenye Windows njia ya ccam huandikwa kwa mikwaju ya mbele na `~` kwa folda ya nyumbani (`~/.local/bin/ccam.exe statusline`), ambayo Git Bash na PowerShell zote huielewa. `settings.json` inayoshirikiwa kupitia kikundi huandikwa mara moja, kupitia kiungo chake, kwa hiyo kila akaunti katika kikundi huipata. Kuiendesha tena huchapisha `(already set)` kwa faili ambazo tayari zinatumia ccam.

## Usalama na hali maalum

- `ccam usage` kamwe haibadilishi faili ya tokeni isipokuwa utoe `--refresh`. Kuingia kukiisha muda, anzisha akaunti mara moja kwa `ccam <name>` na Claude Code itakufanya upya.
- `--refresh` kamwe haigusi akaunti zinazoendeshwa, akaunti zilizoingizwa wala taarifa za kuingia za Keychain. Ufanyaji upya uliokataliwa huonyesha namba zilizohifadhiwa na sababu.
- Endpoint ya matumizi haijaandikiwa nyaraka hadharani na inaweza kubadilika. Ikibadilika, bado unaona namba zilizohifadhiwa, na `ccam doctor` huripoti badiliko la muundo.
- Bila namba zilizohifadhiwa na bila jibu, mstari wa akaunti huonyesha `lookup failed:` na sababu, kwa mfano `no saved usage` kwa `--source cache`.
- Faili za cache huhifadhi asilimia na nyakati za kuanza upya tu, kamwe si tokeni.

## Zinazohusiana

- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Kubadili kwenda akaunti inayofuata](/features/next)
- [Hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow)
- [Amri](/docs/commands)
