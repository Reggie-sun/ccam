---
title: Mga usage limit
description: Ipinapakita ng ccam usage ang 5-oras at lingguhang limit ng bawat account nang sabay; inilalagay ng ccam statusline ang account at usage nito sa Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

Ipinapakita ng `ccam usage` kung gaano kalaki na ang nagamit ng bawat account sa 5-oras na limit, sa lingguhang limit at sa sariling lingguhang limit ng bawat model (gaya ng Fable), bilang isang bar bawat limit kasama ang oras bago ito mag-reset. Gamitin ito bago ka magsimula ng mahabang gawain, o kapag sinabi ng Claude Code na malapit na ang isang limit, para makita kung aling account ang may natitira pa. Ang `ccam statusline` ay ang status line command para sa Claude Code: ipinapakita nito ang pangalan ng account at usage sa loob ng Claude Code at sine-save ang mga numero para sa `ccam usage`.

## Paggamit

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

Bawat account ay may linyang may numero nito, pangalan nito at kung gaano na katagal ang mga numero nito; ang ibig sabihin ng `(cached)` ay galing ang mga ito sa naka-save na file, hindi sa endpoint. Sa ilalim nito, ipinapakita ng isang linya bawat limit ang 20-cell na bar, ang porsiyento at ang oras bago mag-reset ang window: `5-hour`, `Weekly` (lahat ng model), at pagkatapos ang sariling lingguhang window ng bawat model, gaya ng `Fable`. Nagpapakita ng `-` ang window na wala sa account. Nagiging dilaw ang porsiyento mula 80% at pula sa 100%. Sa terminal, lumalabas agad ang bawat account bilang `fetching…`, at iginuguhit ang mga linya nito habang sumasagot ito. Pinapangalanan ng huling linya ang account na pipiliin ng [`ccam next --auto`](/features/next) ayon sa parehong mga patakaran (`auto_max_usage`, `auto_strategy`, ang mga tumatakbong session, mga account lang ng organisasyon ng huling ginamit na account maliban kung nakatakda ang `cross_org`, nilalaktawan ang bawat window na nasa 100%, kaya hindi kailanman pinapangalanan ang account na puno na ang linggo), at lumalabas lang ito kapag mahigit sa isa ang account mo. Nagpapakita ng `not logged in → ccam login <name>` ang account na hindi naka-log in, at ng `no limits (API)` ang account na API key ang gamit.

Pagkatapos ng `ccam statusline install`, magpapakita ang Claude Code ng linya gaya ng `[work] · Opus · 5h 24% · 7d 41%`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--json` | I-print ang mga ulat bilang JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | off |
| `--watch[=N]` | Iguhit ulit tuwing `N` segundo hanggang pindutin mo ang Ctrl+C. Ang `--watch` lang, o value na 0 o mas mababa, ay nangangahulugang 60 | off |
| `-v`, `--verbose` | Magdagdag ng linyang `source:` sa ilalim ng bawat account: kung saan galing ang mga numero nito (ang usage endpoint, ang Claude Code, ang status line o isang naka-save na file) | off |
| `--no-cache` | Huwag pansinin ang 60 segundong cache at magtanong ulit sa endpoint | off |
| `--source live\|claude\|cache` | Nagtatanong sa endpoint ang `live` gamit ang token ng account; ipinapahanap ng `claude` ito sa Claude Code (tingnan sa ibaba); ipinapakita lang ng `cache` ang mga naka-save na numero | `usage_source`, kung wala ay `live` |
| `--refresh` | I-renew ang mga expired na access token, para lang sa mga account na hindi tumatakbo | off |
| `--dry-run` | Tinatanggap: walang binabago ang `ccam usage` maliban sa sarili nitong cache. Tinatanggihan kapag kasama ang `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | off |

Pumapalya ang anumang ibang value ng `--source` na may `--source must be live, claude or cache`.

### Usage mula mismo sa Claude Code

Kapag may `usage_source = "claude"` sa `config.toml` (o `--source claude`), hindi binabasa ng ccam ang login token kahit kailan. Para sa bawat account, sinisimulan nito ang `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` sa folder ng account at ipinapadala ang `get_usage` request ng Claude Code sa standard input (Claude Code 2.1.270 o mas bago); hinahanap ng Claude Code ang usage gamit ang sarili nitong login at sumasagot ito ng 5-oras, lingguhan at model na mga window. Pinipigilan ng `--safe-mode` na magsimula ang mga hook, MCP server at plugin ng account, at walang session na sine-save. Isang pagsisimula ng Claude Code, mga isang segundo, ang kailangan ng bawat account, kaya mas mabagal ang source na ito kaysa sa `live`; sabay-sabay na tinatanong ang mga account, bawat isa sa loob ng 20 segundo, at gumagana ang 60-segundong cache gaya ng dati. Itinuturing ng Claude Code na experimental ang request, kaya kapag huminto itong sumagot, ipinapakita ng ccam ang mga naka-save na numero at iniuulat ito ng `ccam doctor`. Hindi sinisimulan ng `--dry-run` ang Claude Code (nagsusulat ito sa folder ng account) at ipinapakita nito ang mga naka-save na numero. Tumatanggap ang `usage_source` ng `live` o `claude`; pagpipilian lang ng `--source` ang `cache`. Walang tinatanggap na argument ang `ccam usage`.

Binabasa ng `ccam statusline` ang status line JSON ng Claude Code mula sa standard input. Iisa ang option nito, ang `--wrap <command>`: isang umiiral na status line command na ang unang linya ng output ay lalagyan ng pangalan ng account sa unahan. Walang option ang `ccam statusline install`. Parehong tinatanggihan ng dalawa ang `--dry-run` na may `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Paano ito gumagana

Galing sa dalawang lugar ang mga numero: ang hindi dokumentadong usage endpoint ng Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, na tinatawag gamit ang sariling login ng bawat account, at ang mga value na ibinibigay ng Claude Code sa status line command nito. Nasa `internal/usage` ang lahat ng ito.

1. Nilo-load ng `cli.usageCmd` ang `config.toml` at bumubuo ng isang `usage.Account` bawat account: ang folder nito na `~/.local/share/ccam/profiles/<name>/`, ang session folder nito (para makita kung tumatakbo ito) at kung in-import ito sa kinalalagyan nito.
2. Nagsisimula ang `Fetcher.FetchAll` ng isang goroutine bawat account. Iisang `http.Client` ang gamit nilang lahat, kaya muling ginagamit ang mga koneksyon sa endpoint. Puwedeng umabot nang 2 s ang pagkonekta, at may kabuuang 3 s ang bawat account. Hindi sumusunod sa mga redirect ang client, dahil ipapadala ulit ng redirect ang token sa ibang lugar, at hanggang 1 MB lang ng sagot ang binabasa nito.
3. Habang natatapos ang bawat account, iginuguhit ulit ang hilera nito. Ang account na pumapalya o hindi sumasagot sa oras ay bumabalik sa mga naka-save na numero, kaya hindi kailanman napapabagal ng isang mabagal na account ang iba.

Tumatakbo ang bawat account sa sarili nitong goroutine laban sa endpoint; ang nag-time out ay nagpapakita ng mga naka-save na numero nito sa halip.

<!-- motion:usage-fetch -->

Sumasagot sa oras ang `personal` at `work` at nagpapakita ng mga live na numero. Lumalampas sa 3 s ang `side`, kaya binabasa ng `fallback()` ang dalawang cache file nito at ginagamit ang pinakahuling naisulat. Iginuguhit ng `usage.Table()` ang isang block bawat account at pinipili ng `usage.Best()` ang account na pinapangalanan sa huling linya.

Para sa isang account, ginagawa ng `Fetcher.Fetch` ang mga pagsusuring ito nang sunod-sunod:

<!-- motion:usage-cache -->

1. Sa `--source cache`, direkta itong pumupunta sa `fallback()`.
2. Maliban kung nakatakda ang `--no-cache`, ibinabalik nang walang pagbabago ang live na resulta sa `~/.local/state/ccam/usage/<name>.json` na wala pang 60 segundo (`CacheTTL`). Isang request bawat account lang ang gastos kapag dalawang beses na sunod na pinatakbo ang `ccam usage`.
3. Binabasa ng `readCredentials()` ang `.credentials.json` (o ang macOS Keychain) nang hindi ito binabago, at ang key na `claudeAiOauth` lang ang dine-decode. Kapag walang subscription login, `not logged in` ang resulta, o `no limits (API)` kapag may API key ang `.claude.json`.
4. Hindi ginagamit ang expired na access token. Kapag walang `--refresh`, ipinapakita ng account ang mga naka-save na numero nito na may markang expired.
5. Kung hindi, tinatawag ng ccam ang endpoint. Itinuturing na expired na login ang 401 o 403, itinuturing na error ang anumang ibang status o body na hindi nito mabasa, at parehong bumabalik ang mga ito sa mga naka-save na numero. Ang 429 ay nangangahulugang masyadong madalas tinatanong ang endpoint: isinusulat ng ccam ang `<name>.pause` kasama ang oras kung kailan magtatanong ulit (`Retry-After`, o 5 minuto, hindi hihigit sa isang oras), ipinapakita ang mga naka-save na numero na may `rate limited by the usage endpoint; asking again after 15:10`, at hindi tinatawag ang endpoint para sa account na iyon hanggang sa oras na iyon, kahit may `--no-cache`. Lalo lang mapapatagal ng agarang pagsubok ulit ang limit.
6. Isinusulat ang maayos na sagot sa `<name>.json` kasama ang oras ng pagsukat.

Binabasa ng `fallback()` ang parehong `<name>.json` at `<name>.statusline.json` at pinapanatili ang mas bago. Bago nito ipakita ang mga naka-save na numero, sinusuri nito ang oras ng reset ng bawat window: nagpapakita ng 0% ang window na nag-reset na mula nang i-save ang mga numero, dahil hindi na totoo ang lumang porsiyento. Kaya ang `personal` sa itaas ay nagpapakita ng `0%` na walang oras ng reset.

May dalawang format ng response ang endpoint. May array na `limits` ang mas bago na may mga entry na `session` (5-oras), `weekly_all` at `weekly_scoped` (bawat model). May mga object na `five_hour`, `seven_day`, `seven_day_opus` at `seven_day_sonnet` ang mas luma. Binabasa ng `ParseLive` ang dalawa: nananaig ang array na `limits`, at pinupunan ng mga lumang object ang hindi nito sakop. Iniuulat bilang nagbagong format ang body na walang kahit anong limit, at sasabihin iyon sa iyo ng `ccam doctor`.

### Pag-renew ng mga token

Nagro-rotate ang mga refresh token: kapag ni-renew ang isa, makakakuha ka ng bagong refresh token at magiging invalid ang luma. Kung ire-renew ng ccam ang login habang tumatakbo ang Claude Code bilang ang parehong account, masa-sign out ang session na iyon. Kaya binabasa lang ng `ccam usage` ang mga token at hindi kailanman sinusulatan, maliban kung ipapasa mo ang `--refresh`. Pagkatapos, ang `Fetcher.refresh` ay:

1. tumatanggi sa mga account na nirehistro sa kinalalagyan nila gamit ang `ccam import`, dahil baka ginagamit ng karaniwang `claude` ang login na iyon nang hindi alam ng ccam;
2. tumatanggi sa mga login na naka-store sa macOS Keychain;
3. kumukuha ng lock file, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (sa ilalim ng `~/.local/state/ccam/locks/` kapag hindi nakatakda ang variable na iyon), para hindi kailanman ma-renew nang dalawang beses ng dalawang ccam process ang parehong token;
4. habang hawak ang lock, tumatanggi kung may tumatakbong session ng account, at nilalaktawan ang request kung ni-renew na ng ibang process ang token;
5. ipinapadala ang refresh token sa token endpoint at isinusulat pabalik ang mga bagong token sa `.credentials.json` gamit ang `fsutil.AtomicWrite`, at pinapanatili ang bawat ibang key sa file.

### Ang status line collector

Pinapatakbo ng Claude Code ang status line command at ipina-pipe dito ang isang JSON document. Binabasa ng ccam mula rito ang pangalan ng model at ang `rate_limits`:

<!-- motion:usage-statusline -->

1. Hinahanap ng `runStatusline` ang account mula sa `CCAM_PROFILE`, na itinatakda ng ccam kapag sinisimulan nito ang Claude Code, o kung hindi ay sa pagtutugma ng `CLAUDE_CONFIG_DIR` sa mga folder ng account.
2. Binabasa ng `ParseStatusLine` ang `rate_limits.five_hour` at `rate_limits.seven_day`. Walang rate limit ang JSON para sa mga account na API key ang gamit, o bago ang unang response ng session; sa ganitong kaso, walang sine-save.
3. Isinusulat ng `SaveStatusLine` ang `~/.local/state/ccam/usage/<name>.statusline.json`. Ginagamit ito ng `ccam usage` kapag hindi sumasagot ang endpoint.
4. Ipi-print ng ccam ang `[name] · model · 5h N% · 7d N%`. Sa `--wrap`, pinapatakbo nito ang iyong command sa pamamagitan ng `sh -c` na may parehong input at ipi-print ang `[name]` kasunod ang unang linya ng command sa halip. Sa Windows, pinapatakbo nito ang command sa paraan ng Claude Code: gamit ang Git Bash kapag naka-install ito, o kung hindi, gamit ang PowerShell.
5. Kapag nasa 90% o higit pa ang 5-oras o lingguhang usage at may iba pang account, nagtatapos ang linya sa `→ ccam next --auto --continue`, halimbawa `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

May 2 segundo ang naka-wrap na command. Pagkatapos noon, pinapatay ito at pangalan lang ng account ang ipinapakita ng linya, kaya hindi mahaharangan ng mabagal na command ang status line ng Claude Code.

Isinusulat ng `ccam statusline install` ang `statusLine` sa `settings.json` ng bawat account. Kapag may status line command na ang isang account, pinapanatili ito ng ccam bilang `--wrap '<old command>'`. Sa Windows, isinusulat ang path papunta sa ccam gamit ang forward slash at `~` para sa home folder (`~/.local/bin/ccam.exe statusline`), na naiintindihan ng Git Bash at ng PowerShell. Isang beses lang isinusulat ang `settings.json` na shine-share sa pamamagitan ng group, sa pamamagitan ng link nito, kaya nakukuha ito ng bawat account sa group. Kapag pinatakbo ulit, ipi-print ang `(already set)` para sa mga file na gumagamit na ng ccam.

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman binabago ng `ccam usage` ang token file maliban kung ipapasa mo ang `--refresh`. Kapag nag-expire ang isang login, patakbuhin nang isang beses ang account gamit ang `ccam <name>` at ire-renew ito ng Claude Code.
- Hindi kailanman ginagalaw ng `--refresh` ang mga tumatakbong account, in-import na account o Keychain login. Ipinapakita ng tinanggihang refresh ang mga naka-save na numero at ang dahilan.
- Hindi pampublikong dokumentado ang usage endpoint at puwede itong magbago. Kapag nangyari iyon, makikita mo pa rin ang mga naka-save na numero, at iuulat ng `ccam doctor` ang pagbabago sa format.
- Kapag walang naka-save na numero at walang sagot, ipinapakita ng linya ng account ang `lookup failed:` at ang dahilan, halimbawa `no saved usage` kapag may `--source cache`.
- Mga porsiyento at oras ng reset lang ang laman ng mga cache file, hindi kailanman mga token.

## Kaugnay

- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [Paglipat sa susunod na account](/features/next)
- [Workflow kapag naabot ang usage limit](/guides/usage-limit-workflow)
- [Mga command](/docs/commands)
