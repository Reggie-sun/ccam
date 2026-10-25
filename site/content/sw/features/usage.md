---
title: Vikomo vya matumizi
description: ccam usage huonyesha vikomo vya saa 5 na vya wiki vya kila akaunti kwa pamoja, na ccam statusline huweka jina la akaunti na matumizi yake katika Claude Code.
updated: 2026-10-25
order: 7
translation: machine
sourceHash: "91ad20327289"
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
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
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
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Kila akaunti hupata mstari wenye namba yake, jina lake na umri wa namba zake; `(cached)` inamaanisha zinatoka kwenye faili iliyohifadhiwa, si kwenye endpoint. Chini yake, mstari mmoja kwa kila kikomo huonyesha upau wa visanduku 20, asilimia na muda hadi kipindi kianze upya: `5-hour`, `Weekly` (model zote), kisha kipindi cha wiki cha kila model yenyewe, kama `Fable`. Kipindi ambacho akaunti haina huonyesha `-`. Asilimia hugeuka manjano kuanzia 80% na nyekundu kwa 100%. Kipindi ambacho Claude Code iliripoti kuwa kimejaa, ilipokataa ombi kwenye kikomo cha matumizi, huonyesha 100% hadi kianze upya, bila kujali namba zilizohifadhiwa zinasema nini (tazama [rekodi ya kikomo](#the-limit-record)). Kwenye terminal kila akaunti huonekana mara moja kama `fetching…`, na mistari yake huchorwa inapojibu. Mstari wa mwisho hutaja akaunti ambayo [`ccam next --auto`](/features/next) ingechagua kwa kanuni zile zile (`auto_max_usage`, `auto_strategy`, vipindi vinavyoendeshwa, akaunti za shirika la akaunti iliyotumika mwisho pekee isipokuwa `cross_org` imewekwa, kila kipindi kilicho 100% kikirukwa, kwa hiyo akaunti ambayo wiki yake imejaa haitajwi kamwe), na huonekana tu ukiwa na zaidi ya akaunti moja. Kwa mkakati chaguo-msingi wa `headroom` husomeka `→ Most headroom: <name>`; kwa `expiring` au `balance` hutoa sababu, kama `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Viti (seats) vya Team na Enterprise.** Chaguo za kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi, na [`auto_switch`](/features/next#auto-switch) unapokiwasha) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kwa hiyo, kutumia chaguo za kiotomatiki kati ya viti kunaweza kukiuka masharti ya Anthropic. Kama utazitumia huko ni uamuzi wako: wasiliana na shirika lako na Anthropic kabla ya kuzitegemea. Mipango ya binafsi (Pro, Max) inaelezwa katika [dokezo la masharti kwenye Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Akaunti ambayo haijaingia huonyesha `not logged in → ccam login <name>`, na akaunti ya API key huonyesha `no limits (API)`. Tokeni ya ufikiaji iliyoisha muda huonyesha ``access token expired; `ccam <name>` renews it when it starts``, juu ya namba zilizohifadhiwa zikiwepo: Claude Code huifanya upya katika uanzishaji unaofuata. Taarifa za kuingia zilizoisha muda kabisa, kwa sababu refresh token iliisha muda au Claude Code ilizifuta baada ya seva kuzikataa, huonyesha `` login expired → `ccam login <name>` ``.

Baada ya `ccam statusline install`, Claude Code huonyesha mstari kama `[work] · Opus · 5h 24% · 7d 41%`.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--json` | Huchapisha ripoti kama JSON: `name`, `slot`, `status`, `limits` (pamoja na `credits` akaunti ikiwa na credits za matumizi zimewashwa), `source`, `at`, `error`, na `limited` na `limited_windows` wakati rekodi ya kikomo inaweka vipindi kuwa 100%. Bila akaunti yoyote huchapisha `[]` | imezimwa |
| `--watch[=N]` | Huchora upya kila sekunde `N` hadi ubonyeze Ctrl+C. `--watch` peke yake, au thamani ya 0 au chini, inamaanisha 60. Thamani huwekwa baada ya `=`: `--watch 10` hukataliwa pamoja na kidokezo | imezimwa |
| `-v`, `--verbose` | Huongeza mstari wa `source:` chini ya kila akaunti: namba zake zilikotoka (endpoint ya matumizi, Claude Code, mstari wa hali, faili iliyohifadhiwa, au kikomo cha matumizi ambacho Claude Code iliripoti). Rekodi ya kikomo inapokuwa juu ya namba, hufuatiwa na mstari kama `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | imezimwa |
| `--no-cache` | Hupuuza cache ya sekunde 60 na huuliza endpoint tena | imezimwa |
| `--source live\|claude\|cache` | `live` huuliza endpoint kwa tokeni ya akaunti; `claude` huiachia Claude Code itafute (hapa chini); `cache` huonyesha namba zilizohifadhiwa tu | `usage_source`, vinginevyo `live` |
| `--refresh` | Hufanya upya tokeni za ufikiaji zilizoisha muda, kwa akaunti ambazo haziendeshwi tu. Kufanya upya kunaweza kuchukua hadi sekunde 30 | imezimwa |
| `--dry-run` | Inakubaliwa: `ccam usage` haibadilishi chochote isipokuwa cache yake yenyewe. Hukataliwa pamoja na `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | imezimwa |

Thamani nyingine yoyote ya `--source` hushindwa na `--source must be live, claude or cache`.

### Matumizi kutoka Claude Code yenyewe

Kwa `usage_source = "claude"` juu kabisa ya `config.toml` (kabla ya `[[group]]`, `[[profile]]` au `[[folder]]` ya kwanza), au `--source claude`, ccam haisomi tokeni ya kuingia hata kidogo. Kwa kila akaunti huanzisha `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` katika folda ya akaunti na hutuma ombi la `get_usage` la Claude Code kupitia standard input (Claude Code 2.1.270 au baadaye); Claude Code hutafuta matumizi kwa taarifa zake yenyewe za kuingia na hujibu kwa vipindi vya saa 5, vya wiki na vya model. `--safe-mode` huzuia hooks, seva za MCP na plugins za akaunti zisianze, na hakuna kipindi kinachohifadhiwa. Kila akaunti huhitaji Claude Code kuanzishwa, takriban sekunde moja, kwa hiyo chanzo hiki ni cha polepole kuliko `live`; akaunti huulizwa kwa wakati mmoja, kila moja ndani ya sekunde 20, na cache ya sekunde 60 hutumika kama kawaida. Kuingia ambako tokeni yake ya ufikiaji imeisha muda, au itaisha ndani ya dakika 5, hupata sekunde 50, kwa kuwa Claude Code huifanya upya inapoanza na kufanya upya kunakokatizwa kunaweza kupoteza kuingia. Utafutaji unaoachwa huiomba Claude Code imalize kwa SIGTERM, ili ikamilishe ufanyaji upya iliouanza (Windows haina ishara kama hiyo; huko husimamishwa). Taarifa za kuingia zilizoisha muda kabisa huonyesha `` login expired → `ccam login <name>` `` bila kuanzisha Claude Code. Claude Code huliweka ombi hili alama ya majaribio, kwa hiyo likiacha kujibu, ccam huonyesha namba zilizohifadhiwa na `ccam doctor` huliripoti. Claude Code inapojibu kwamba utafutaji wake yenyewe umeshindwa, ccam huonyesha namba zilizohifadhiwa pamoja na `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, si mpango usio na vikomo. `--dry-run` haianzishi Claude Code (huandika kwenye folda ya akaunti) na huonyesha namba zilizohifadhiwa. `usage_source` hupokea `live` au `claude`; `cache` ni chaguo la `--source` pekee. `ccam usage` haipokei hoja.

`ccam statusline` husoma JSON ya mstari wa hali ya Claude Code kutoka standard input. Chaguo lake moja ni `--wrap <command>`: amri ya mstari wa hali iliyopo ambayo mstari wake wa kwanza wa matokeo huwekewa jina la akaunti mbele. `ccam statusline install` haipokei chaguo. Zote mbili hukataa `--dry-run` kwa `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Jinsi inavyofanya kazi

Namba zinatoka sehemu mbili: endpoint ya matumizi ya Anthropic isiyo na nyaraka, `GET https://api.anthropic.com/api/oauth/usage`, inayoitwa kwa taarifa za kuingia za kila akaunti yenyewe, na thamani ambazo Claude Code hukabidhi amri yake ya mstari wa hali. Yote haya yako katika `internal/usage`.

1. `cli.usageCmd` hupakia `config.toml` na huunda `usage.Account` moja kwa kila akaunti: folda yake `~/.local/share/ccam/profiles/<name>/`, alama zake za vipindi (pamoja na rekodi za Claude Code yenyewe katika folda ya akaunti, huonyesha kama inaendeshwa) na kama iliingizwa mahali ilipo.
2. `Fetcher.FetchAll` huanzisha goroutine moja kwa kila akaunti. Zote hushiriki `http.Client` moja, kwa hiyo miunganisho kwenye endpoint hutumiwa tena. Kuunganisha kunaweza kuchukua sekunde 2, na kila akaunti hupata sekunde 3 kwa jumla. Client haifuati uelekezaji, kwa kuwa uelekezaji ungetuma tokeni mahali pengine, na husoma kwa kiwango cha juu MB 1 ya jibu.
3. Kila akaunti inapomaliza, safu yake huchorwa upya. Akaunti inayoshindwa au isiyojibu kwa wakati hurudi kwenye namba zilizohifadhiwa, kwa hiyo akaunti moja ya polepole kamwe haizuii nyingine.

Kila akaunti huendeshwa katika goroutine yake dhidi ya endpoint; ile inayopitwa na muda huonyesha namba zake zilizohifadhiwa badala yake.

<!-- motion:usage-fetch -->

`personal` na `work` hujibu kwa wakati na huonyesha namba za moja kwa moja. `side` huchukua zaidi ya sekunde 3, kwa hiyo `fallback()` huonyesha namba zake zilizohifadhiwa badala yake. `usage.Table()` huchora kizuizi kimoja kwa kila akaunti na `usage.Best()` huchagua akaunti inayotajwa kwenye mstari wa mwisho.

Kwa akaunti moja, `Fetcher.Fetch` hufanya ukaguzi huu kwa mpangilio:

<!-- motion:usage-cache -->

1. Kwa `--source cache`, huenda moja kwa moja kwa `fallback()`.
2. Isipokuwa `--no-cache` imewekwa, matokeo ya moja kwa moja katika `~/.local/state/ccam/usage/<name>.json` yaliyopimwa chini ya sekunde 60 zilizopita (`CacheTTL`) hutumiwa, yakipimwa kwa hali ya sasa: kipindi kilichoanza upya tangu hapo huonyesha 0%, na faili yenye tarehe iliyo mbele ya saa ya kompyuta haihesabiwi kuwa mpya. Kuendesha `ccam usage` mara mbili mfululizo hugharimu ombi moja kwa kila akaunti.
3. `readCredentials()` husoma `.credentials.json` (au Keychain ya macOS) bila kuibadilisha, na ufunguo wa `claudeAiOauth` pekee ndio hufasiriwa. Kukosekana kwa kuingia kwa usajili kunamaanisha `not logged in`, au `no limits (API)` wakati `.claude.json` ina API key. Taarifa za kuingia ambazo refresh token yake imeisha muda, au ambazo Claude Code imezifuta, ni `` login expired → `ccam login <name>` `` na hazitafutwi.
4. Tokeni ya ufikiaji iliyoisha muda haitumiwi. Bila `--refresh`, akaunti huonyesha namba zake zilizohifadhiwa pamoja na ``access token expired; `ccam <name>` renews it when it starts``.
5. Vinginevyo ccam huita endpoint. 401 au 403 huhesabiwa kama kuingia kulikoisha muda (`login expired: authentication failed (401)`), hali nyingine yoyote au mwili wa jibu usioweza kusomwa huhesabiwa kama kosa, na zote mbili hurudi kwenye namba zilizohifadhiwa. Utafutaji uliokatishwa kwa Ctrl+C husema `cancelled`. 429 inamaanisha endpoint inaulizwa mara nyingi mno: ccam huandika `<name>.pause` pamoja na muda wa kuuliza tena (`Retry-After`, au dakika 5, si zaidi ya saa moja), huonyesha namba zilizohifadhiwa pamoja na `rate limited by the usage endpoint; asking again after 15:10`, na haiiti endpoint kwa akaunti hiyo hadi wakati huo, hata kwa `--no-cache`. Kujaribu tena mara moja kungeendeleza kikomo tu.
6. Jibu zuri huandikwa kwenye `<name>.json` pamoja na muda lilipopimwa.

`fallback()` huchukua utafutaji uliohifadhiwa mpya zaidi (`<name>.json`, au `<name>.claude.json` kutoka kwa utafutaji wa Claude Code) na, `<name>.statusline.json` ikiwa iliandikwa baada yake, huweka vipindi vya saa 5 na vya wiki vya mstari wa hali juu yake; vipindi vya model, ambavyo mstari wa hali hauna, hubaki kama utafutaji ulivyovihifadhi. Kabla ya kuonyesha namba zilizohifadhiwa, hukagua muda wa kuanza upya wa kila kipindi: kipindi kilichoanza upya tangu namba zihifadhiwe huonyesha 0%, kwa sababu asilimia ya zamani haitumiki tena. Ndiyo sababu `personal` hapo juu huonyesha `0%` bila muda wa kuanza upya.

Endpoint ina miundo miwili ya majibu. Mpya zaidi una safu ya `limits` yenye vipengele vya `session` (saa 5), `weekly_all` na `weekly_scoped` (kwa kila model). Wa zamani una vitu vya `five_hour`, `seven_day`, `seven_day_opus` na `seven_day_sonnet`. `ParseLive` husoma yote miwili: safu ya `limits` hushinda, na vitu vya zamani hujaza kile ambacho haikihusishi. Mwili wa jibu usio na vikomo kabisa huripotiwa kama muundo uliobadilika, jambo ambalo `ccam doctor` hukujulisha.

### Rekodi ya kikomo

Claude Code inaposimamisha zamu kwenye kikomo cha matumizi, hook ya [plugin ya ccam](/docs/install#claude-code-plugin) husoma ombi lililokataliwa kutoka kwenye mazungumzo: ni kipindi kipi kilichojaa (saa 5, wiki, Opus au Sonnet) na lini kinaanza upya. Huandika hilo kwenye `~/.local/state/ccam/usage/<name>.limit.json`. `Fetch` huweka rekodi hiyo juu ya kila kitu kingine inachojua, kutoka kila chanzo: kila kipindi kilichomo huhesabiwa kama 100% hadi kianze upya, bila kujali namba zilizohifadhiwa ni za zamani au za chini kiasi gani, kwa hiyo `ccam next --auto`, `ccam pick`, `--fallback` na [`auto_switch`](/features/next#auto-switch) hazichagui akaunti hiyo tena kabla ya hapo. Kama hakuna kingine kinachojulikana, akaunti huonyesha rekodi hiyo pekee (`source: the usage limit Claude Code reported` kwa `-v`). Vipindi ambavyo hakuna mtu anayevitumia (`-p`, vipindi vya chinichini) haviandiki rekodi, wala usimamishaji wa shirika (matumizi yameisha, kikundi kilichowekwa $0), ambao wasimamizi wake wanaweza kuuondoa wakati wowote. `ccam login` kwenye akaunti nyingine ya Claude huondoa rekodi; kufuta faili pia huifuta.

### Kufanya upya tokeni

Refresh tokens huzunguka: kufanya upya moja hukupa refresh token mpya na hubatilisha ile ya zamani. Ikiwa ccam ingefanya upya kuingia wakati Claude Code inaendeshwa kama akaunti ile ile, kipindi hicho kingetolewa nje ya akaunti. Kwa hiyo `ccam usage` husoma tokeni na kamwe haiziandiki, isipokuwa utoe `--refresh`. Kisha `Fetcher.refresh`:

1. hukataa akaunti zilizosajiliwa mahali zilipo kwa `ccam import`, kwa kuwa `claude` ya kawaida inaweza kutumia kuingia huko bila ccam kujua;
2. hukataa taarifa za kuingia zilizohifadhiwa katika Keychain ya macOS;
3. huchukua faili ya kufuli, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, ili michakato miwili ya ccam kamwe isifanye upya tokeni ile ile mara mbili; wakati inashikiliwa, akaunti huonyesha `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` na `logout` nazo huichukua, na hukataa wakati kufanya upya kunaendelea;
4. chini ya kufuli, hukataa ikiwa kipindi cha akaunti kinaendeshwa, bila kujali kilianzishwa vipi (kupitia ccam, `claude` ya kawaida, IDE, au kipindi cha chinichini na daemon iliyo nyuma yake, ambayo hufanya upya kuingia yenyewe), na huruka ombi ikiwa mchakato mwingine tayari umefanya upya tokeni;
5. haitumi chochote wakati Claude Code inashikilia kufuli zake yenyewe za kufanya upya katika folda ya akaunti (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code inafanya upya kuingia huko wakati huo huo;
6. hutuma refresh token kwenye endpoint ya tokeni na husubiri jibu hadi sekunde 30, muda ambao Claude Code huupa ufanyaji upya wake yenyewe. Ombi likishatumwa kamwe halikatizwi, hata kwa Ctrl+C, kwa kuwa huenda seva tayari imebadilisha refresh token na jibu lake pekee ndilo bado linashikilia kuingia: Ctrl+C au hangup humaliza utafutaji mwingine, na ccam hutoka (130 baada ya Ctrl+C) jibu likishahifadhiwa;
7. huhifadhi chini ya kufuli ya uandishi ya Claude Code (`.storage-write.lock`) na hubadilisha tu sehemu za tokeni za `.credentials.json` kama ilivyo wakati huo: kile ambacho Claude Code iliandika katikati (kuingia kwa MCP) hubaki, kuingia ambako Claude Code ilifanya upya katikati huhifadhiwa, na kutoka kwenye akaunti kamwe hakutenguliwi. Jibu linapotoa muda wa kuishi wa refresh token mpya, huhifadhiwa kama `refreshTokenExpiresAt`, kama Claude Code inavyofanya.

Ufanyaji upya usiofanikiwa husema sababu, kwenye mstari wa akaunti au juu ya namba zake zilizohifadhiwa: `token not renewed: a session is running`, kwa mfano, au `login may be lost: …` wakati ombi lilitumwa lakini jibu lake halikuweza kutumiwa au kuhifadhiwa, pamoja na ushauri wa kuendesha `ccam login <name>` ikiwa baadaye Claude Code itakuomba uingie.

### Mkusanyaji wa mstari wa hali

Claude Code huendesha amri ya mstari wa hali na huipitishia hati ya JSON. ccam husoma jina la model na `rate_limits` kutoka humo:

<!-- motion:usage-statusline -->

1. `runStatusline` hupata akaunti kutoka `CCAM_PROFILE`, ambayo ccam huiweka inapoanzisha Claude Code, au vinginevyo kwa kulinganisha `CLAUDE_CONFIG_DIR` na folda za akaunti.
2. `ParseStatusLine` husoma `rate_limits.five_hour` na `rate_limits.seven_day`. JSON haina vikomo vya maombi kwa akaunti za API key, au kabla ya jibu la kwanza la kipindi; basi hakuna kinachohifadhiwa.
3. `SaveStatusLine` huandika `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` huitumia wakati endpoint haijibu.
4. ccam huchapisha `[name] · model · 5h N% · 7d N%`. Kwa `--wrap`, huendesha amri yako kupitia `sh -c` kwa ingizo lile lile na huchapisha matokeo yake kama Claude Code inavyoonyesha mstari wa hali: kila mstari usio mtupu, `[name]` ikiwa mbele ya wa kwanza. Amri inayoshindwa huonyesha jina la akaunti peke yake. Kwenye Windows huendesha amri kama Claude Code inavyofanya: kwa Git Bash ikiwa imesakinishwa, la sivyo kwa PowerShell.
5. Matumizi ya saa 5 au ya wiki yanapokuwa 90% au zaidi na kuna akaunti nyingine, mstari huishia na `→ ccam next --auto --continue`, kwa mfano `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Amri iliyozungukwa haina kikomo chake yenyewe cha muda: Claude Code huipa amri ya mstari wa hali dakika 10 na husimamisha uendeshaji sasisho jipya linapoanza, kwa hiyo amri ya polepole huonekana kama ingeonekana bila ccam. Claude Code inapomaliza ccam, amri nayo humalizika (kwenye Linux na macOS pamoja na kila mchakato iliouanzisha).

`ccam statusline install` huandika `statusLine` katika `settings.json` ya kila akaunti na huhifadhi funguo nyingine za kitu hicho, kama `padding`, `refreshInterval` na `hideVimModeIndicator`. Akaunti ikiwa tayari ina amri ya mstari wa hali, ccam huihifadhi kama `--wrap '<old command>'`; kwenye Windows kama `--wrap-base64 <base64>`, ambayo haihitaji alama za nukuu katika shell yoyote kati ya hizo mbili. Mstari wa hali ambao tayari unaendesha ccam, kutoka njia nyingine au ukiwa na nukuu kwa namna ccam ya zamani ilivyouandika, hubadilishwa badala ya kuzungukwa, ukihifadhi kile unachokizunguka. Kwenye Windows njia ya ccam huandikwa kwa mikwaju ya mbele na `~` kwa folda ya nyumbani (`~/.local/bin/ccam.exe statusline`), ambayo Git Bash na PowerShell zote huielewa. Njia inayohitaji nukuu huandikwa kamili kwa shell ambayo Claude Code hutumia: `'…'` kwa Git Bash, `& '…'` kwa PowerShell. `settings.json` inayoshirikiwa kupitia kikundi huandikwa mara moja, kupitia kiungo chake, kwa hiyo kila akaunti katika kikundi huipata. Kuiendesha tena huchapisha `(already set)` kwa faili ambazo tayari zinatumia ccam.

## Usalama na hali maalum

- `ccam usage` kamwe haibadilishi faili ya tokeni isipokuwa utoe `--refresh`. Tokeni ya ufikiaji ikiisha muda, anzisha akaunti mara moja kwa `ccam <name>` na Claude Code itaifanya upya. Taarifa za kuingia zilizoisha muda kabisa huhitaji `ccam login <name>`.
- `--refresh` kamwe haigusi akaunti zinazoendeshwa (hata Claude Code ilipoanzishwa bila ccam), akaunti zilizoingizwa wala taarifa za kuingia za Keychain. Ufanyaji upya uliokataliwa huonyesha namba zilizohifadhiwa na sababu (`token not renewed: …`).
- Endpoint ya matumizi haijaandikiwa nyaraka hadharani na inaweza kubadilika. Ikibadilika, bado unaona namba zilizohifadhiwa, na `ccam doctor` huripoti badiliko la muundo.
- Bila namba zilizohifadhiwa na bila jibu, mstari wa akaunti huonyesha `lookup failed:` na sababu, kwa mfano `no saved usage` kwa `--source cache`.
- Faili za cache huhifadhi asilimia na nyakati za kuanza upya tu, kamwe si tokeni.

## Zinazohusiana

- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Kubadili kwenda akaunti inayofuata](/features/next)
- [Hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow)
- [Amri](/docs/commands)
