---
title: Paglipat sa susunod na account
description: Sa ccam next, bubukas ang Claude Code bilang ibang account; --auto para sa pinakamaluwag, --continue para dalhin ang usapan, --fallback para ulitin ang -p job.
updated: 2026-10-24
order: 8
translation: machine
sourceHash: "a38daed5dd97"
---

Binubuksan ng `ccam next` ang Claude Code bilang ibang account nang hindi mo ito pinapangalanan. Gamitin ito kapag naabot ng isang account ang usage limit nito: kinukuha ng simpleng `ccam next` ang account na kasunod ng huli mong ginamit, kinukuha ng `--auto` ang may pinakamaraming natitira, at dinadala ng `--continue` ang huling usapan ng folder na ito sa bagong account, nagshe-share man o hindi ng kasaysayan ang dalawa. Sa loob ng Claude Code, sinasabi sa iyo ng [ccam plugin](#inside-claude-code) kapag pinahinto ng limit ang isang session, at ipinapasa ng `/ccam:switch` ang usapan. Para sa mga script na nagpapatakbo ng `claude -p`, kusang ginagawa ng `ccam --fallback` ang paglipat (tingnan ang [Patakbuhin ulit sa limit](#run-again-at-the-limit)).

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon. Sinasabi ng help center ng Anthropic na "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (nakatalaga ang bawat seat ng Team plan sa iisang partikular na tao, at siya lang ang puwedeng gumamit nito) at "Users can't route their work through another user's account to use that user's remaining usage" (hindi puwedeng idaan ng mga user ang trabaho nila sa account ng ibang user para magamit ang natitirang usage ng user na iyon) ([mga Team plan](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ganoon din sa [mga Enterprise plan](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kaya maaaring lumabag sa terms ng Anthropic ang paggamit ng awtomatikong pagpili sa pagitan ng mga seat. Magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito. Saklaw ng [paalala tungkol sa terms sa FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ang mga personal na plan (Pro, Max).

## Paggamit

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

Ipi-print ng ccam ang paglipat bilang `→ <from> → <to>` sa standard error, kasama ang dahilan ng pagpili ng `--auto`, saka sisimulan ang Claude Code bilang ang bagong account sa parehong paraan ng `ccam <name>`. Ang `<from>` ay ang account na huli mong ginamit, o sa `--continue`, ang account na kinabibilangan ng usapan. Ipinapasa sa `claude` nang walang pagbabago ang mga argument na kasunod ng `--`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--auto` | Piliin ang naka-log in na account ayon sa usage (tingnan sa ibaba), at laktawan ang mga account na nasa 100% | off |
| `--max-usage <1-100>` | Kasama ang `--auto`: unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring\|balance` | Kasama ang `--auto`: pinakamababang 5-oras na usage, ang pinakamalapit na 5-oras na reset, o ang lingguhang quota na unang masasayang kapag hindi nagamit | `auto_strategy`, kung wala ay `headroom` |
| `--cross-org` | Kasama ang `--auto`: pumili rin ng mga account ng ibang organisasyon (tingnan sa ibaba) | `cross_org`, kung wala ay off |
| `--wait[=<duration>]` | Kasama ang `--auto`: kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito. Inilalagay ang value pagkatapos ng `=`: tinatanggihan ang `--wait 2h` kasama ang pahiwatig na `--wait takes its value after '=': --wait=2h` | off; `5h` kapag ibinigay nang mag-isa |
| `-c`, `--continue` | Dalhin ang huling usapan ng folder na ito sa bagong account (tingnan sa ibaba) | off |
| `-- <claude args>` | Mga argument na ipinapasa sa `claude` pagkatapos ng sariling argument ng ccam | wala |
| `--dry-run` | Sinusunod: ipi-print ang `→ <from> → <to>` at ang command line na patatakbuhin nito, nang hindi sinisimulan ang Claude Code o itinatala ang pagsisimula | off |

Gumagana nang magkasama ang `--auto` at `--continue`. Sa `--continue`, inilalagay ang mga argument ng `--resume` bago ang sarili mong mga argument.

Ang sariling mga option ng Claude Code, gaya ng `--dangerously-skip-permissions` o `--model`, ay inilalagay pagkatapos ng `--`. Hindi tulad ng `ccam <name>`, na ipinapasa ang lahat ng kasunod ng pangalan, binabasa mismo ng `ccam next` ang mga argument nito, kaya kapag walang `--` ay humihinto ito sa `ccam: unknown flag: --dangerously-skip-permissions`.

## Paano ito gumagana

Isinusulat ng bawat pagsisimula sa pamamagitan ng ccam ang oras sa `~/.local/state/ccam/last-used/<name>`. Binabasa ng `ccam next` ang mga file na ito, tinitingnan kung aling mga account ang naka-log in, at pumipili.

1. Bumubuo ang `cli.candidates` ng isang `next.Candidate` bawat account mula sa `config.toml`: ang account, ang huling oras ng pagsisimula nito (`launcher.LastUsed`) at kung may gumaganang login ito (`creds.Has`: isang Claude login sa `.credentials.json`, o sa macOS sa Keychain, na hindi pa buong nag-expire o binura ng Claude Code). Walang plan limit ang folder na gumagamit ng API key sa halip (isang Console login, o account na idinagdag gamit ang `ccam add --api`) at hindi ito kailanman pinipili.
2. Kapag walang `--auto`, inaayos ng `next.Rotate` ang mga account ayon sa slot number at hinahanap ang pinakahuling sinimulan. Iyon ang account na nililisan mo. Kung wala pang account na nasimulan, ito ang pinakamababang slot.
3. Mula roon, sumusulong ito ayon sa pagkakasunod-sunod ng slot, bumabalik sa una mula sa huling slot, at ibinabalik ang unang account na naka-log in.
4. Inihahanda at sinisimulan ng `launcher.LaunchWith` ang napiling account nang eksaktong gaya ng gagawin ng `ccam <name>`: sinusuri nito ang mga link ng account, inilalapat ang MCP sharing, nirerehistro ang session at pinapalitan ang ccam ng `claude` sa pamamagitan ng `syscall.Exec`.

Nagsisimula ang pagdaan sa account na huling ginamit at humihinto sa unang naka-log in na account kasunod nito:

<!-- motion:next-rotate -->

Ang `work` ang pinakahuling sinimulan, kaya nagsisimula ang pagdaan pagkatapos nito. Walang login ang `side` at nilalaktawan ito, at pinipili ang `lab`. Kung naka-log out din sana ang `lab`, babalik ang pagdaan sa `personal` (`% len` sa diagram). Hindi kailanman ibinabalik ng rotation ang account na pinagmulan mo: kapag walang ibang naka-log in na account, humihinto ito na may error.

### --auto

Nagpapasya ang `--auto` ayon sa usage sa halip na sa pagkakasunod-sunod ng slot. Pinapanatili ng `cli.pickAuto` ang mga naka-log in na account, sabay-sabay na hinihingi ang usage nilang lahat gamit ang parehong `Fetcher.FetchAll` na ginagamit ng `ccam usage`, at hinahayaang pumili ang `usage.Policy.Choose()`:

<!-- motion:next-auto -->

Hinuhusgahan ang isang account ayon sa 5-oras at lingguhang window nito, at, kapag kasama ang `--model` sa mga argument para sa Claude Code (pagkatapos ng `--`), ayon din sa lingguhang window ng model na iyon (`opus`, `sonnet`, `haiku` o `fable` sa pangalan). Hinahati ng `Policy.Choose()` ang mga account sa tatlo:

1. **May natitira**: bawat window ay mas mababa sa `auto_max_usage` (`config.toml`, default 100). Ang strategy ang pumipili sa kanila: kinukuha ng `headroom` (default) ang pinakamababang 5-oras na usage; kinukuha ng `expiring` ang 5-oras na window na pinakamalapit mag-reset, para magamit ang natitirang quota nito bago ito mawala; kinukuha ng `balance` ang account na may pinakamaraming natitirang lingguhang quota bawat araw hanggang mag-reset ang linggo nito, na ibinibilang laban dito ang bawat tumatakbong ccam session, para hindi masayang ang lingguhang quota at kumalat ang mga parallel session.
2. **Lampas**: mas mababa sa 100% pero nasa o lampas sa `auto_max_usage`. Ginagamit lang kapag walang account na may natitira, inuuna ang pinakamababang 5-oras na usage.
3. **Puno**: may window na nasa 100%, o hindi alam ang 5-oras na usage. Hindi kailanman pinipili.

Itinuturing na 100% ang isang window hanggang sa reset nito kapag tinanggihan doon ng Claude Code ang isang request: itinatala ito ng hook ng [ccam plugin](#inside-claude-code) (`~/.local/state/ccam/usage/<name>.limit.json`), anuman ang sabihin ng mga naka-save na numero. Kaya hindi pinipili ulit bago ang reset nito ang account na kaaabot lang sa limit nito, kahit luma na ang mga numero nito.

Inilalagay ang mga setting para sa lahat ng account sa itaas ng `config.toml`, bago ang unang `[[group]]`, `[[profile]]` o `[[folder]]`. Ang setting na nasa ibaba ng isa sa mga iyon ay pag-aari ng table na iyon: pinapangalanan ito ng ccam sa bawat pagsisimula, at tumatanggi ang mga command na nagse-save ng config hanggang mailipat ito.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Mga organisasyon.** Nananatili ang mga awtomatikong pagpili sa mga account na kauri ng nililisan mo: ang parehong team o enterprise na organisasyon, o ang mga personal na plan (Pro, Max) nang magkakasama, kaya hindi naitutuloy sa pribadong account ang trabaho mula sa account ng kompanya, o ang kabaligtaran. Alam ng ccam ang organisasyon ng bawat account mula sa login nito (itinatala ito ng `ccam login`; binabasa mula sa `.claude.json` ang mas lumang mga login); puwede kahit saan ang account na hindi nito matukoy. Iisa ang limit ng dalawang folder na naka-log in sa parehong Claude account, kaya hindi kailanman pumapalit ang isa sa isa pa. Inaalis ng `--cross-org`, o ng `cross_org = true` sa `config.toml`, ang patakaran ng organisasyon.

Kapag walang mapipili, lumalabas ang ccam na may code 75 (`EX_TEMPFAIL`), para matukoy ng script ang "subukan mamaya" mula sa isang pagkabigo. Sa `--wait`, natutulog ito sa halip hanggang mag-reset ang humaharang na window ng unang punong account (30 segundo pagkatapos ng oras ng reset, kapag ipinapakita na ng endpoint ang bagong window), nagpi-print ng `ccam: no account has room; waiting until 15:10`, at pumipili ulit; kapag walang account na nagsasabi kung kailan ito magre-reset, titingin ulit ito pagkalipas ng 10 minuto. Kung walang lumuwag bago matapos ang paghihintay, lumalabas ito na may 75. Sinusunod ng paghahanap ang mga patakaran ng `ccam usage`: muling ginagamit ang mga numerong wala pang 60 segundo, may 3 segundo ang bawat account, ang account na hindi sumasagot, o nag-expire ang access token, ay hinuhusgahan ayon sa mga naka-save na numero nito, at nakapatong sa mga iyon ang tala ng limit. Itinuturing na 0% ang naka-save na window na lumipas na ang oras ng reset, kaya hindi nalalampasan ang account na nag-reset na ang window mula noon. Hindi kailanman nire-renew ng `--auto` ang mga token. Hindi tulad ng simpleng rotation, puwedeng piliin ng `--auto` ang account na huli mong ginamit kung ito pa rin ang may pinakamaraming natitira.

### --continue

1. Hinahanap ng ccam ang huling usapan ng folder na ito sa bawat account: sa mga `<session>.jsonl` file sa `projects/<folder>/` (ang pangalan ng folder ay ang working folder, na naresolba ang mga link, na ginawang `-` ang bawat character na hindi titik o numero, gaya ng pagpapangalan ng Claude Code), ang may pinakabagong prompt o sagot. Nilalaktawan ang mga session na wala nito, gaya ng session na nagpatakbo lang ng `/mcp`, `/usage` o `/status`, o ang walang lamang file. Kung nag-iwan ng tala ang isang session (isang usage limit o `/ccam:switch`, tingnan sa ibaba) at walang usapang may prompt o sagot na lampas sa 2 minuto pagkatapos nito, ang eksaktong session na iyon ang ginagamit; hindi binibilang ang mga hilerang isinusulat ng Claude Code sa `/exit`.
2. Ang account na kinabibilangan nito ang account na nililisan mo, at pinipili ang bagong account gaya ng nasa itaas, sa loob ng parehong organisasyon. Kapag ilang account ang nagshe-share ng kasaysayang iyon, kabilang ito sa account na ang `.claude.json` ay nagsasabing ito ang huling nagpatakbo ng usapan (`lastSessionId`), kung hindi ay sa huling ginamit.
3. Kung shine-share ng bagong account ang kasaysayang iyon (ang parehong group folder), sinisimulan ng ccam ang `claude --resume <session id>`, ang mismong usapang iyon (bubuksan ng `--continue` kung alinman ang pinakabago sa folder).
4. Kung hindi, sinisimulan nito ang `claude --resume <the conversation's file> --fork-session`: kinokopya ng Claude Code ang buong usapan, sa ilalim ng bagong session id, sa sariling kasaysayan ng bagong account, at nananatili ang orihinal gaya ng dati. Mula noon, nahahanap na ito ng `--continue` ng bagong account.

Sinasabi ng ccam kung gaano kalaki ang usapan: binabasa ulit ng bagong account ang lahat ng ito sa unang sagot nito. Hindi kailanman shine-share ang mga prompt cache sa pagitan ng mga organisasyon, at sariling organisasyon ang bawat personal na plan, kaya kapag magkaiba ang kinabibilangan ng dalawang account, idinadagdag ng ccam ang `prompt caches are not shared between organizations, so the first reply reads it all again`. Kapag naka-sign in ang lumang account sa mga MCP server na hindi pa naka-sign in ang bago, pinapangalanan ng ccam ang mga ito: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. Kapag shine-share ang `history` sa isang group, naiiwasan ang pagkopya, nananatiling iisa ang usapan, at shine-share din ang mga checkpoint na pinagbabalikan ng code ng `/rewind` (`file-history/`).

## Sa loob ng Claude Code

Nagdadagdag ang [ccam plugin](/docs/install#claude-code-plugin) ng dalawang hook (Claude Code 2.1.78 o mas bago):

- Kapag huminto ang isang turn sa usage limit, itinatala nito ang usapan at nagpapakita ng desktop notification sa mga terminal na sumusuporta nito (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` Kilala nito ang bawat mensaheng itinuturing ng Claude Code na usage limit: ang 5-oras, lingguhan, karaniwan at model na limit (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` para sa allowance ng isang model), ang mga usage credit at spend limit, at ang organisasyong nagpahinto sa account. Kapag naisulat na ng Claude Code sa usapan ang tinanggihang request, itinatala rin ng hook ang punong window hanggang sa reset nito ([ang tala ng limit](/features/usage#the-limit-record)), maliban sa pagpapahinto ng isang organisasyon, na puwedeng alisin ng mga admin nito anumang oras. Hindi usage limit ang maikling throttle ng server o ang rate limit ng isang API key, at hinahayaan lang ang mga iyon, gayundin ang paghinto ng isang subagent at ang mga session na walang nakaupong tao (`-p`, mga background session). Naghihintay rin sa session ang Claude Code 2.1.234 at mas bago, at kusang nagpapatuloy pagkatapos ng 5-oras na reset; binibigyan ka ng notification ng pagpipiliang lumipat na lang.
- Itinatala ng `/ccam:switch [account]` ang usapan at tinatapos ang turn bago gumawa ng anumang request: i-type ang `/exit`, saka ipagpapatuloy ito ng `ccam next --continue`, bilang account na pinangalanan mo o bilang kasunod na account. Sa session na wala pang usapan, sinasabi nitong walang maipagpapatuloy.

Nasa `~/.local/state/ccam/handoff/<account>/<session>.json` ang tala, isang file bawat session, kaya hindi kailanman napapatungan ng dalawang session ng iisang account ang tala ng isa't isa. Itinatago ito nang 12 oras at isang beses lang ginagamit, kapag magsisimula na ang susunod na session; kapag pumalya ang pagsisimula, naiiwan ito para sa susunod na subok. Nakatali ito sa folder kung saan nagsimula ang session, kahit pagkatapos ng `cd` sa shell ng Claude Code. Ina-update ng `ccam rename` ang mga pangalan dito; kung natanggal na ang account na binabanggit nito, sasabihin iyon ng ccam at pipili ito gaya ng dati.

### Relay

Kapag may `relay = true` sa itaas ng `config.toml`, nananatili ang ccam sa tabi ng mga interactive session sa halip na palitan ang sarili nito ng Claude Code, at pagkatapos ng `/ccam:switch`, kusa nitong ipinagpapatuloy ang usapan ng session na iyon (hindi nito pinapansin ang talang iniwan ng ibang session): tinatapos ng hook ang Claude Code, at pinapatakbo ng ccam ang `ccam next --continue` gamit ang tala ng session na iyon. Kapag walang pinangalanang account, kinukuha nito ang may pinakamaraming natitira, na isinasaalang-alang ang model ng usapan. Sa Windows, hindi kayang tapusin ng hook ang Claude Code: i-type ang `/exit`, at magpapatuloy ang ccam. Isang naghihintay na ccam process bawat session ang kapalit nito, kaya naka-off ito bilang default.

Pinapanatili ng ipinasang session ang hindi naibabalik ng pag-resume ng usapan: mga option gaya ng `--add-dir`, `--mcp-config`, ang mga listahan ng tool, `--settings`, `--plugin-dir` at ang mga option ng system prompt, ang permission mode ng usapan at ang effort level ng session. Ang Claude Code mismo ang nagbabalik ng `--model`, `--name` at `--agent` mula sa usapan. Pinapangalanan sa `ccam: not carried over: …` ang anumang iba pa (`--worktree`, `--debug`, `--remote-control`, mga option na hindi kilala ng ccam). Hindi kailanman itinatanong ng paglipat ang tanong tungkol sa shared na kasaysayan, at itinatapon ang anumang na-type mo habang tumatakbo ito.

Nagpapatuloy lang ang relay kapag natapos ang Claude Code sa paraang hiniling ng `/ccam:switch`: hindi kapag nauna kang nag-type ng `/exit`, pumindot ng Ctrl+C habang natatapos ito, o nakatanggap ang ccam ng hangup o ibang signal, at hindi kapag wala na ang terminal. Sa ganoong kaso, naiiwan ang tala, at ipinagpapatuloy ng `ccam next --continue` ang usapan sa ibang pagkakataon. Hindi tinatapos ng hook ang session na sinimulan sa ilalim ng relay ng mas lumang ccam, o ang session na walang nakaupong tao (`-p`, mga background session); ang mensaheng nagsasabing i-type ang `/exit` ang nakukuha nito sa halip. Tumatakbo nang walang relay ang `claude remote-control`.

## Patakbuhin ulit sa limit

Para sa interactive session na huminto ang `ccam next`. Hindi ka kayang hintayin ng non-interactive na run, ang `claude -p` sa isang script o naka-iskedyul na job, kaya hinahayaan ng `--fallback` bago ang pangalan ng account na ang ccam ang lumipat para dito:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. Gumagana lang ang `--fallback` kasama ang `-p` o `--print`; kung hindi, humihinto ang ccam at itinuturo ang `ccam next --auto --continue`.
2. Isang beses binabasa ang piped na standard input, hanggang 10 MiB, at ibinibigay sa bawat pagtatangka. Gaya ng Claude Code, naghihintay ang ccam nang hindi hihigit sa 3 segundo para magsimula ito (hindi sa `--input-format stream-json`); walang input na naibibigay sa run ang pipe na nananatiling bukas at walang ipinapadala, at may linya sa standard error.
3. Sinisimulan ng ccam ang `claude` bilang child process sa halip na `exec`, pinipigil ang standard output nito, at ipinapasa ang standard error habang dumarating ito.
4. Kapag lumabas ang `claude` nang may 0, o pumalya dahil sa ibang dahilan, ipinapasa nang walang pagbabago ang pinigil na output at ang exit code. Ganoon din ang run na tinapos ng isang signal (exit code na lampas sa 128): hiniling itong huminto, at hindi na ito pinapatakbo ulit.
5. Itinuturing na usage limit ang pagkabigo kapag sinasabi iyon ng mensahe ng Claude Code, alinman sa mga mensahe ng usage limit na kilala ng [hook](#inside-claude-code), pati sa loob ng output na `--output-format json` o `stream-json`, o kapag ipinapakita ng usage ng account na nasa 100% ang 5-oras, lingguhan o `--model` na window nito (kasama ang tala ng limit).
6. Saka itinatapon ang output ng pagtatangkang iyon at pinapatakbo ang parehong command bilang ibang naka-log in na account: ang may pinakamababang 5-oras na usage, saka ang mga account na hindi matingnan ang usage, ayon sa pagkakasunod-sunod ng numero. Nilalaktawan ang mga account na nasa 100%, hindi kailanman sinusubukan ang mga account na may API key login, at sinusubukan ang bawat account nang hindi hihigit sa isang beses.
7. Kapag kasama sa mga argument ang `-c`, `--continue`, `-r` o `--resume`, ang mga account lang sa parehong history group ang sinusubukan, para mahanap ang usapan.
8. Ang mga account lang ng parehong organisasyon ang sinusubukan, gaya sa `--auto` (inaalis ito ng `cross_org = true`).

Pinipili ang susunod na account ayon sa parehong mga patakaran ng `--auto` (`auto_max_usage`, `auto_strategy`, ang window ng `--model`). Kapag wala nang natitirang account, ipi-print ng ccam ang output ng huling pagtatangka, sasabihing walang naka-log in na account na may natitira pa, at lalabas na may code 75. Sa `--wait` bago ang pangalan ng account (`ccam --fallback --wait=2h work -p …`), hinihintay nito sa halip na mag-reset ang unang account, saka puwedeng subukan ulit ang bawat account; lagi itong natutulog nang kahit isang beses muna, kaya hindi agad pinapatakbo ulit ang account na kaaabot lang sa limit nito. Saklaw ng tagal ang lahat ng paghihintay ng run nang magkakasama, hindi ang bawat paghihintay. Kapag isinulat nang may espasyo (`--wait 30m`), tinatanggihan ito kasama ang parehong pahiwatig ng `ccam next`. Gumagana ang `--` sa unahan ng mga argument ng Claude Code gaya ng sa anumang pagsisimula (`ccam --fallback work -- -p …`). Ipi-print ng `ccam --dry-run --fallback work -p …` ang command ng unang pagtatangka.

Hindi kusang inililipat ang mga interactive session. Sa halip, kapag may `ccam statusline install`, idinadagdag ng status line ang `→ ccam next --auto --continue` kapag umabot sa 90% ang 5-oras o lingguhang usage ([Mga usage limit](/features/usage)).

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman pinipili ang mga account na walang login. Kapag walang ibang naka-log in na account, humihinto ang simpleng `ccam next` na may `no logged-in account other than <name>`.
- Sa `--continue` at walang usapan sa folder: `no account has a conversation in <folder> to continue`.
- Kapag walang natirang account dahil sa patakaran ng organisasyon: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Sa `--auto` at walang naka-log in na account: `no logged-in account to pick`. Kapag puno ang bawat account o hindi alam ang usage nito: `` no account with known usage and headroom (see `ccam usage`) ``. Parehong nagtatapos ang dalawa sa exit code 75.
- Gumagana lang ang `--max-usage`, `--strategy`, `--wait` at `--cross-org` kasama ang `--auto`, at tinatanggihan ang value na labas sa 1-100 o ibang strategy.
- Kinokopya ng fork ang usapan, hindi ang iba pang state ng account: hindi kasamang dinadala ang mga checkpoint ng `file-history/` ng orihinal na account maliban kung shine-share ang kasaysayan.
- Gumagawa ng mga network request ang `--auto` kahit may `--dry-run`. Nagbabasa lang ito, at hindi kailanman sumusulat sa mga token file.
- Nagsisimula ang napiling account sa parehong daan ng `ccam <name>`, kaya nalalapat ang parehong mga pagsusuri, kasama ang babala kapag may tumatakbo nang ibang Claude Code na nagshe-share ng parehong kasaysayan sa folder na ito.
- Isinusulat lang ang `last-used` kapag talagang nagsimula ang Claude Code. Hindi ito ginagalaw ng `--dry-run`, kaya hindi isinusulong ng dry run ang rotation.
- Pinipigil ng `--fallback` ang standard output hanggang matapos ang isang pagtatangka, kaya dumarating nang buo sa dulo ang streaming output gaya ng `--output-format stream-json`. Hindi magagamit ang `--fallback` ng client na nagsi-stream ng `stream-json` na input at naghihintay ng mga sagot bago magpadala ng higit pa. Nananatiling binago ang mga file na binago ng unang pagtatangka bago nito naabot ang limit; sinisimulan ulit ng susunod na pagtatangka ang prompt.

## Kaugnay

- [Mga usage limit](/features/usage)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Workflow kapag naabot ang usage limit](/guides/usage-limit-workflow)
