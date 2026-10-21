---
title: Iyakokin amfani
description: ccam usage yana nuna iyakokin awa 5, na mako da na kowane model (Fable) na kowane asusu a lokaci ɗaya a matsayin sanduna, kuma ccam statusline yana saka sunan asusun da amfaninsa a Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "0e9587b79ce8"
---

`ccam usage` yana nuna nawa kowane asusu ya yi amfani da shi daga iyakar awa 5, iyakar mako da iyakar mako ta kowane model (kamar Fable), a matsayin sanda ɗaya ga kowace iyaka tare da lokacin da ya rage kafin ta sake saitawa. Yi amfani da shi kafin ka fara dogon aiki, ko idan Claude Code ya gaya maka iyaka ta kusa, don ganin asusun da ke da sauran sarari. `ccam statusline` shi ne umarnin layin matsayi (status line) na Claude Code: yana nuna sunan asusun da amfaninsa a cikin Claude Code kuma yana ajiye lambobin don `ccam usage`.

## Amfani

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

Kowane asusu yana samun layi mai lambarsa, sunansa da tsawon lokacin da lambobinsa suka yi; `(cached)` yana nufin sun fito daga fayil da aka ajiye, ba daga endpoint ba. A ƙarƙashinsa, layi ɗaya ga kowace iyaka yana nuna sanda mai ƙofofi 20, kashi da lokacin da ya rage kafin taga ta sake saitawa: `5-hour`, `Weekly` (dukkan model), sannan taga ta mako ta kowane model, kamar `Fable`. Tagar da asusun ba shi da ita tana nuna `-`. Kashi yana zama rawaya daga 80% kuma ja a 100%. A terminal kowane asusu yana bayyana nan take a matsayin `fetching…`, kuma ana zana layukansa yayin da ya amsa. Layin ƙarshe yana ambaton asusun da [`ccam next --auto`](/features/next) zai zaɓa bisa ƙa'idodi iri ɗaya (`auto_max_usage`, `auto_strategy`, zaman da ke gudana, asusan ƙungiyar asusun da aka yi amfani da shi na ƙarshe kawai sai dai idan an saita `cross_org`, ana tsallake duk tagar da ta kai 100%, don haka asusun da makonsa ya cika ba a taɓa ambaton sa ba), kuma yana bayyana ne kawai idan kana da asusu fiye da ɗaya. Asusun da ba a shiga da shi ba yana nuna `not logged in → ccam login <name>`, kuma asusun API key yana nuna `no limits (API)`.

Bayan `ccam statusline install`, Claude Code yana nuna layi kamar `[work] · Opus · 5h 24% · 7d 41%`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--json` | Buga rahotannin a matsayin JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | a kashe |
| `--watch[=N]` | Sake zanawa kowane daƙiƙa `N` har sai ka danna Ctrl+C. `--watch` shi kaɗai, ko ƙima ta 0 ko ƙasa da haka, yana nufin 60 | a kashe |
| `-v`, `--verbose` | Ƙara layin `source:` ƙarƙashin kowane asusu: inda lambobinsa suka fito (endpoint na amfani, Claude Code, layin matsayi ko fayil da aka ajiye) | a kashe |
| `--no-cache` | Yi watsi da cache na daƙiƙa 60 kuma sake tambayar endpoint | a kashe |
| `--source live\|claude\|cache` | `live` yana tambayar endpoint da token na asusun; `claude` yana sa Claude Code ya nemo shi (a ƙasa); `cache` yana nuna lambobin da aka ajiye kawai | `usage_source`, in ba haka ba `live` |
| `--refresh` | Sabunta access tokens da suka ƙare, ga asusan da ba sa aiki kawai | a kashe |
| `--dry-run` | Ana karɓar sa: `ccam usage` ba ya canza komai sai cache nasa. Ana ƙin sa tare da `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | a kashe |

Duk wata ƙimar `--source` dabam tana gazawa da `--source must be live, claude or cache`.

### Amfani daga Claude Code kansa

Da `usage_source = "claude"` a `config.toml` (ko `--source claude`), ccam ba ya karanta token na shiga ko kaɗan. Ga kowane asusu yana fara `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` a foldar asusun kuma yana aika buƙatar `get_usage` ta Claude Code a standard input (Claude Code 2.1.270 ko sabo); Claude Code yana nemo amfanin da nasa shigar kuma yana amsawa da tagogin awa 5, na mako da na model. `--safe-mode` yana hana hooks, sabobin MCP da plugins na asusun farawa, kuma ba a ajiye zama. Kowane asusu yana ɗaukar farawar Claude Code ɗaya, kusan daƙiƙa ɗaya, don haka wannan tushe ya fi `live` jinkiri; ana tambayar asusan a lokaci ɗaya, kowanne a cikin daƙiƙa 20, kuma cache na daƙiƙa 60 yana aiki kamar yadda aka saba. Claude Code yana yi wa buƙatar alama a matsayin ta gwaji, don haka idan ta daina amsawa, ccam yana nuna lambobin da aka ajiye kuma `ccam doctor` yana bayar da rahoto. `--dry-run` ba ya fara Claude Code (yana rubutu a foldar asusun) kuma yana nuna lambobin da aka ajiye. `usage_source` yana karɓar `live` ko `claude`; `cache` zaɓi ne na `--source` kawai. `ccam usage` ba ya karɓar arguments.

`ccam statusline` yana karanta JSON na layin matsayi na Claude Code daga standard input. Zaɓinsa ɗaya shi ne `--wrap <command>`: umarnin layin matsayi da ke akwai wanda layin fitarwa na farko ke samun sunan asusun a gabansa. `ccam statusline install` ba ya karɓar zaɓuɓɓuka. Duka biyun suna ƙin `--dry-run` da `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Yadda yake aiki

Lambobin suna fitowa daga wurare biyu: endpoint na amfani na Anthropic da ba a rubuta bayaninsa ba, `GET https://api.anthropic.com/api/oauth/usage`, wanda ake kira da shigar kowane asusu, da ƙimomin da Claude Code ke miƙa wa umarnin layin matsayinsa. Duk wannan yana cikin `internal/usage`.

1. `cli.usageCmd` yana loda `config.toml` kuma yana gina `usage.Account` ɗaya ga kowane asusu: foldarsa `~/.local/share/ccam/profiles/<name>/`, foldar zamansa (don ganin ko yana aiki) da ko an shigo da shi a wurinsa.
2. `Fetcher.FetchAll` yana fara goroutine ɗaya ga kowane asusu. Dukkansu suna raba `http.Client` ɗaya, don haka ana sake amfani da haɗi zuwa endpoint ɗin. Haɗawa na iya ɗaukar 2 s, kuma kowane asusu yana samun 3 s gaba ɗaya. Client ɗin ba ya bin turawa (redirects), tunda turawa za ta sake aika token zuwa wani wuri, kuma yana karanta 1 MB na amsar a mafi yawa.
3. Yayin da kowane asusu ya kammala, ana sake zana layinsa. Asusun da ya gaza ko bai amsa a kan lokaci ba yana komawa ga lambobin da aka ajiye, don haka asusu ɗaya mai jinkiri ba ya taɓa riƙe sauran.

Kowane asusu yana gudana a goroutine nasa a kan endpoint ɗin; wanda lokacinsa ya ƙare yana nuna lambobinsa da aka ajiye maimakon haka.

<!-- motion:usage-fetch -->

`personal` da `work` suna amsawa a kan lokaci kuma suna nuna lambobi kai tsaye. `side` yana ɗaukar fiye da 3 s, don haka `fallback()` yana karanta fayilolin cache nasa biyu kuma yana amfani da wanda aka rubuta na ƙarshe. `usage.Table()` yana zana toshe ɗaya ga kowane asusu kuma `usage.Best()` yana zaɓar asusun da aka ambata a layin ƙarshe.

Ga asusu ɗaya, `Fetcher.Fetch` yana yin waɗannan dubawa a tsari:

<!-- motion:usage-cache -->

1. Da `--source cache`, yana zuwa `fallback()` kai tsaye.
2. Sai dai idan an saita `--no-cache`, sakamako kai tsaye a `~/.local/state/ccam/usage/<name>.json` da bai kai daƙiƙa 60 ba (`CacheTTL`) ana mayar da shi kamar yadda yake. Gudanar da `ccam usage` sau biyu a jere yana kashe buƙata ɗaya ga kowane asusu.
3. `readCredentials()` yana karanta `.credentials.json` (ko Keychain na macOS) ba tare da canza shi ba, kuma maɓallin `claudeAiOauth` ne kawai ake karantawa. Rashin shigar biyan kuɗi (subscription) yana nufin `not logged in`, ko `no limits (API)` idan `.claude.json` yana riƙe API key.
4. Ba a amfani da access token da ya ƙare. Ba tare da `--refresh` ba, asusun yana nuna lambobinsa da aka ajiye tare da alamar ya ƙare.
5. In ba haka ba ccam yana kiran endpoint ɗin. 401 ko 403 ana ɗaukar su a matsayin shigar da ta ƙare, duk wani matsayi dabam ko jikin amsa da ba zai iya karantawa ba ana ɗaukar sa a matsayin kuskure, kuma duka biyun suna komawa ga lambobin da aka ajiye. 429 yana nufin ana tambayar endpoint ɗin da yawa: ccam yana rubuta `<name>.pause` da lokacin sake tambaya (`Retry-After`, ko minti 5, ba fiye da awa ɗaya ba), yana nuna lambobin da aka ajiye tare da `rate limited by the usage endpoint; asking again after 15:10`, kuma ba ya kiran endpoint ɗin ga wannan asusun har zuwa lokacin, ko da `--no-cache`. Sake gwadawa nan take zai ci gaba da riƙe iyakar ne kawai.
6. Ana rubuta amsa mai kyau a `<name>.json` tare da lokacin da aka auna ta.

`fallback()` yana karanta duka `<name>.json` da `<name>.statusline.json` kuma yana riƙe wanda ya fi sabo. Kafin ya nuna lambobin da aka ajiye, yana duba lokacin sake saitin kowace taga: tagar da ta sake saitawa tun lokacin da aka ajiye lambobin tana nuna 0%, domin tsohon kashi ba ya aiki kuma. Shi ya sa `personal` a sama yana nuna `0%` ba tare da lokacin sake saiti ba.

Endpoint ɗin yana da tsarin amsa biyu. Sabon yana da array na `limits` mai shigarwar `session` (awa 5), `weekly_all` da `weekly_scoped` (na kowane model). Tsohon yana da abubuwan `five_hour`, `seven_day`, `seven_day_opus` da `seven_day_sonnet`. `ParseLive` yana karanta duka biyun: array na `limits` ne ke nasara, kuma tsofaffin abubuwan suna cike abin da bai ƙunsa ba. Jikin amsa da ba shi da iyaka ko ɗaya ana bayar da rahotonsa a matsayin tsarin da ya canza, wanda `ccam doctor` ke gaya maka game da shi.

### Sabunta token

Refresh tokens suna juyawa: sabunta ɗaya yana ba ka sabon refresh token kuma yana soke tsohon. Idan ccam ya sabunta shigar yayin da Claude Code ke gudana a matsayin asusu ɗaya, za a fitar da wannan zaman. Don haka `ccam usage` yana karanta token kuma ba ya taɓa rubuta su, sai dai idan ka ba da `--refresh`. Sannan `Fetcher.refresh`:

1. yana ƙin asusan da aka yi wa rajista a wurinsu da `ccam import`, tunda `claude` na yau da kullum zai iya amfani da wannan shigar ba tare da ccam ya sani ba;
2. yana ƙin shigar da aka ajiye a Keychain na macOS;
3. yana ɗaukar fayil ɗin kulle, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (ƙarƙashin `~/.local/state/ccam/locks/` idan ba a saita wannan mai canjin ba), don haka processes biyu na ccam ba sa taɓa sabunta token ɗaya sau biyu;
4. ƙarƙashin kullen, yana ƙi idan zaman asusun yana gudana, kuma yana tsallake buƙatar idan wani process ya riga ya sabunta token ɗin;
5. yana aika refresh token zuwa endpoint na token kuma yana rubuta sababbin token a `.credentials.json` da `fsutil.AtomicWrite`, yana riƙe kowane sauran maɓalli a cikin fayil ɗin.

### Mai tattara layin matsayi

Claude Code yana gudanar da umarnin layin matsayi kuma yana tura masa takardar JSON ta bututu (pipe). ccam yana karanta sunan model da `rate_limits` daga cikinta:

<!-- motion:usage-statusline -->

1. `runStatusline` yana nemo asusun daga `CCAM_PROFILE`, wanda ccam ke saitawa lokacin da yake fara Claude Code, ko kuma ta kwatanta `CLAUDE_CONFIG_DIR` da foldojin asusai.
2. `ParseStatusLine` yana karanta `rate_limits.five_hour` da `rate_limits.seven_day`. JSON ɗin ba shi da iyakokin amfani ga asusan API key, ko kafin amsar farko ta zaman; to ba a ajiye komai.
3. `SaveStatusLine` yana rubuta `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` yana amfani da shi idan endpoint ɗin bai amsa ba.
4. ccam yana buga `[name] · model · 5h N% · 7d N%`. Da `--wrap`, yana gudanar da umarninka ta `sh -c` da shigarwa iri ɗaya kuma yana buga `[name]` sannan layin farko na umarnin maimakon haka. A Windows yana gudanar da umarnin yadda Claude Code ke yi: da Git Bash idan an girka shi, in ba haka ba da PowerShell.
5. Idan amfanin awa 5 ko na mako ya kai 90% ko fiye kuma akwai wani asusu, layin yana ƙarewa da `→ ccam next --auto --continue`, misali `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Umarnin da aka naɗe yana samun daƙiƙa 2. Bayan haka ana kashe shi kuma layin yana nuna sunan asusun kaɗai, don haka umarni mai jinkiri ba zai iya toshe layin matsayi na Claude Code ba.

`ccam statusline install` yana rubuta `statusLine` a cikin `settings.json` na kowane asusu. Idan asusu ya riga yana da umarnin layin matsayi, ccam yana riƙe shi a matsayin `--wrap '<old command>'`. A Windows ana rubuta hanyar ccam da slash na gaba da `~` don foldar gida (`~/.local/bin/ccam.exe statusline`), wanda Git Bash da PowerShell duka suke fahimta. `settings.json` da aka raba ta rukuni ana rubuta shi sau ɗaya, ta mahaɗinsa, don haka kowane asusu a rukunin yana samun sa. Sake gudanar da shi yana buga `(already set)` ga fayilolin da suka riga suna amfani da ccam.

## Tsaro da yanayi na musamman

- `ccam usage` ba ya taɓa canza fayil ɗin token sai dai idan ka ba da `--refresh`. Idan shiga ta ƙare, fara asusun sau ɗaya da `ccam <name>` kuma Claude Code zai sabunta ta.
- `--refresh` ba ya taɓa asusan da ke aiki, asusan da aka shigo da su ko shigar Keychain. Sabuntawar da aka ƙi tana nuna lambobin da aka ajiye da dalilin.
- Ba a rubuta bayanin endpoint na amfani a bainar jama'a ba, kuma yana iya canzawa. Idan ya canza, har yanzu kana ganin lambobin da aka ajiye, kuma `ccam doctor` yana bayar da rahoton canjin tsarin.
- Idan babu lambobin da aka ajiye kuma babu amsa, layin asusun yana nuna `lookup failed:` da dalilin, misali `no saved usage` da `--source cache`.
- Fayilolin cache suna riƙe kashi da lokutan sake saiti kawai, ba taɓa token ba.

## Masu alaƙa

- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Sauyawa zuwa asusu na gaba](/features/next)
- [Matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow)
- [Umarnai](/docs/commands)
