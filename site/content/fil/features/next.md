---
title: Paglipat sa susunod na account
description: Sa ccam next, bubukas ang Claude Code bilang ibang account; --auto para sa pinakamaluwag, --continue para dalhin ang usapan, --fallback para ulitin ang -p job.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

Binubuksan ng `ccam next` ang Claude Code bilang ibang account nang hindi mo ito pinapangalanan. Gamitin ito kapag naabot ng isang account ang usage limit nito: kinukuha ng simpleng `ccam next` ang account na kasunod ng huli mong ginamit, kinukuha ng `--auto` ang may pinakamaraming natitira, at dinadala ng `--continue` ang huling usapan ng folder na ito sa bagong account, nagshe-share man o hindi ng kasaysayan ang dalawa. Sa loob ng Claude Code, sinasabi sa iyo ng [ccam plugin](#inside-claude-code) kapag pinahinto ng limit ang isang session, at ipinapasa ng `/ccam:switch` ang usapan. Para sa mga script na nagpapatakbo ng `claude -p`, kusang ginagawa ng `ccam --fallback` ang paglipat (tingnan ang [Patakbuhin ulit sa limit](#run-again-at-the-limit)).

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
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

Ipi-print ng ccam ang paglipat bilang `→ <from> → <to>` sa standard error, kasama ang dahilan ng pagpili ng `--auto`, saka sisimulan ang Claude Code bilang ang bagong account sa parehong paraan ng `ccam <name>`. Ang `<from>` ay ang account na huli mong ginamit, o sa `--continue`, ang account na kinabibilangan ng usapan. Ipinapasa sa `claude` nang walang pagbabago ang mga argument na kasunod ng `--`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--auto` | Piliin ang naka-log in na account ayon sa usage (tingnan sa ibaba), at laktawan ang mga account na nasa 100% | off |
| `--max-usage <1-100>` | Kasama ang `--auto`: unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring\|balance` | Kasama ang `--auto`: pinakamababang 5-oras na usage, ang pinakamalapit na 5-oras na reset, o ang lingguhang quota na unang masasayang kapag hindi nagamit | `auto_strategy`, kung wala ay `headroom` |
| `--cross-org` | Kasama ang `--auto`: pumili rin ng mga account ng ibang organisasyon (tingnan sa ibaba) | `cross_org`, kung wala ay off |
| `--wait[=<duration>]` | Kasama ang `--auto`: kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito | off; `5h` kapag ibinigay nang mag-isa |
| `-c`, `--continue` | Dalhin ang huling usapan ng folder na ito sa bagong account (tingnan sa ibaba) | off |
| `-- <claude args>` | Mga argument na ipinapasa sa `claude` pagkatapos ng sariling argument ng ccam | wala |
| `--dry-run` | Sinusunod: ipi-print ang `→ <from> → <to>` at ang command line na patatakbuhin nito, nang hindi sinisimulan ang Claude Code o itinatala ang pagsisimula | off |

Gumagana nang magkasama ang `--auto` at `--continue`. Sa `--continue`, inilalagay ang mga argument ng `--resume` bago ang sarili mong mga argument.

## Paano ito gumagana

Isinusulat ng bawat pagsisimula sa pamamagitan ng ccam ang oras sa `~/.local/state/ccam/last-used/<name>`. Binabasa ng `ccam next` ang mga file na ito, tinitingnan kung aling mga account ang naka-log in, at pumipili.

1. Bumubuo ang `cli.candidates` ng isang `next.Candidate` bawat account mula sa `config.toml`: ang account, ang huling oras ng pagsisimula nito (`launcher.LastUsed`) at kung may login ito (`creds.Has`: isang nababasang `.credentials.json`, o sa macOS isang Keychain item).
2. Kapag walang `--auto`, inaayos ng `next.Rotate` ang mga account ayon sa slot number at hinahanap ang pinakahuling sinimulan. Iyon ang account na nililisan mo. Kung wala pang account na nasimulan, ito ang pinakamababang slot.
3. Mula roon, sumusulong ito ayon sa pagkakasunod-sunod ng slot, bumabalik sa una mula sa huling slot, at ibinabalik ang unang account na naka-log in.
4. Inihahanda at sinisimulan ng `launcher.LaunchWith` ang napiling account nang eksaktong gaya ng gagawin ng `ccam <name>`: sinusuri nito ang mga link ng account, inilalapat ang MCP sharing, nirerehistro ang session at pinapalitan ang ccam ng `claude` sa pamamagitan ng `syscall.Exec`.

Nagsisimula ang pagdaan sa account na huling ginamit at humihinto sa unang naka-log in na account kasunod nito:

<!-- motion:next-rotate -->

Ang `work` ang pinakahuling sinimulan, kaya nagsisimula ang pagdaan pagkatapos nito. Walang login ang `side` at nilalaktawan ito, at pinipili ang `lab`. Kung naka-log out din sana ang `lab`, babalik ang pagdaan sa `personal` (`% len` sa diagram). Hindi kailanman ibinabalik ng rotation ang account na pinagmulan mo: kapag walang ibang naka-log in na account, humihinto ito na may error.

### --auto

Nagpapasya ang `--auto` ayon sa usage sa halip na sa pagkakasunod-sunod ng slot. Pinapanatili ng `cli.pickAuto` ang mga naka-log in na account, sabay-sabay na hinihingi ang usage nilang lahat gamit ang parehong `Fetcher.FetchAll` na ginagamit ng `ccam usage`, at hinahayaang pumili ang `usage.Policy.Choose()`:

<!-- motion:next-auto -->

Hinuhusgahan ang isang account ayon sa 5-oras at lingguhang window nito, at, kapag kasama ang `--model` sa mga argument para sa Claude Code, ayon din sa lingguhang window ng model na iyon (`opus`, `sonnet`, `haiku` o `fable` sa pangalan). Hinahati ng `Policy.Choose()` ang mga account sa tatlo:

1. **May natitira**: bawat window ay mas mababa sa `auto_max_usage` (`config.toml`, default 100). Ang strategy ang pumipili sa kanila: kinukuha ng `headroom` (default) ang pinakamababang 5-oras na usage; kinukuha ng `expiring` ang 5-oras na window na pinakamalapit mag-reset, para magamit ang natitirang quota nito bago ito mawala; kinukuha ng `balance` ang account na may pinakamaraming natitirang lingguhang quota bawat araw hanggang mag-reset ang linggo nito, na ibinibilang laban dito ang bawat tumatakbong ccam session, para hindi masayang ang lingguhang quota at kumalat ang mga parallel session.
2. **Lampas**: mas mababa sa 100% pero nasa o lampas sa `auto_max_usage`. Ginagamit lang kapag walang account na may natitira, inuuna ang pinakamababang 5-oras na usage.
3. **Puno**: may window na nasa 100%, o hindi alam ang 5-oras na usage. Hindi kailanman pinipili.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Mga organisasyon.** Nananatili ang mga awtomatikong pagpili sa mga account na kauri ng nililisan mo: ang parehong team o enterprise na organisasyon, o ang mga personal na plan (Pro, Max) nang magkakasama, kaya hindi naitutuloy sa pribadong account ang trabaho mula sa account ng kompanya, o ang kabaligtaran. Alam ng ccam ang organisasyon ng bawat account mula sa login nito (itinatala ito ng `ccam login`; binabasa mula sa `.claude.json` ang mas lumang mga login); puwede kahit saan ang account na hindi nito matukoy. Iisa ang limit ng dalawang folder na naka-log in sa parehong Claude account, kaya hindi kailanman pumapalit ang isa sa isa pa. Inaalis ng `--cross-org`, o ng `cross_org = true` sa `config.toml`, ang patakaran ng organisasyon.

Kapag walang mapipili, lumalabas ang ccam na may code 75 (`EX_TEMPFAIL`), para matukoy ng script ang "subukan mamaya" mula sa isang pagkabigo. Sa `--wait`, natutulog ito sa halip hanggang mag-reset ang humaharang na window ng unang punong account (30 segundo pagkatapos ng oras ng reset, kapag ipinapakita na ng endpoint ang bagong window), nagpi-print ng `ccam: no account has room; waiting until 15:10`, at pumipili ulit; kapag walang account na nagsasabi kung kailan ito magre-reset, titingin ulit ito pagkalipas ng 10 minuto. Kung walang lumuwag bago matapos ang paghihintay, lumalabas ito na may 75. Sinusunod ng paghahanap ang mga patakaran ng `ccam usage`: muling ginagamit ang mga numerong wala pang 60 segundo, may 3 segundo ang bawat account, at ang account na hindi sumasagot, o nag-expire ang login, ay hinuhusgahan ayon sa mga naka-save na numero nito. Itinuturing na 0% ang naka-save na window na lumipas na ang oras ng reset, kaya hindi nalalampasan ang account na nag-reset na ang window mula noon. Hindi kailanman nire-renew ng `--auto` ang mga token. Hindi tulad ng simpleng rotation, puwedeng piliin ng `--auto` ang account na huli mong ginamit kung ito pa rin ang may pinakamaraming natitira.

### --continue

1. Hinahanap ng ccam ang huling usapan ng folder na ito sa bawat account: ang pinakabagong `<session>.jsonl` sa `projects/<folder>/` (ang pangalan ng folder ay ang working folder na ginawang `-` ang bawat character na hindi titik o numero, gaya ng pagpapangalan ng Claude Code). Kung nag-iwan ng tala ang session (isang usage limit o `/ccam:switch`, tingnan sa ibaba) at walang mas bagong naisulat mula noon, ang eksaktong session na iyon ang ginagamit.
2. Ang account na kinabibilangan nito ang account na nililisan mo, at pinipili ang bagong account gaya ng nasa itaas, sa loob ng parehong organisasyon.
3. Kung shine-share ng bagong account ang kasaysayang iyon (ang parehong group folder), sinisimulan ng ccam ang `claude --resume <session id>`, ang mismong usapang iyon (bubuksan ng `--continue` kung alinman ang pinakabago sa folder).
4. Kung hindi, sinisimulan nito ang `claude --resume <the conversation's file> --fork-session`: kinokopya ng Claude Code ang buong usapan, sa ilalim ng bagong session id, sa sariling kasaysayan ng bagong account, at nananatili ang orihinal gaya ng dati. Mula noon, nahahanap na ito ng `--continue` ng bagong account.

Sinasabi ng ccam kung gaano kalaki ang usapan: binabasa ulit ng bagong account ang lahat ng ito sa unang sagot nito. Hindi kailanman shine-share ang mga prompt cache sa pagitan ng mga organisasyon, kaya sa `--cross-org`, buong laki ang gastos ng unang sagot na iyon. Kapag shine-share ang `history` sa isang group, naiiwasan ang pagkopya, nananatiling iisa ang usapan, at shine-share din ang mga checkpoint na pinagbabalikan ng code ng `/rewind` (`file-history/`).

## Sa loob ng Claude Code

Nagdadagdag ang [ccam plugin](/docs/install#claude-code-plugin) ng dalawang hook (Claude Code 2.1.78 o mas bago):

- Kapag huminto ang isang turn sa usage limit (`You've hit your session limit`, ang lingguhan o isang model limit), itinatala nito ang usapan at nagpapakita ng desktop notification sa mga terminal na sumusuporta nito (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Hindi usage limit ang maikling throttle ng server o ang rate limit ng isang API key, at hinahayaan lang ang mga iyon. Naghihintay rin sa session ang Claude Code 2.1.234 at mas bago, at kusang nagpapatuloy pagkatapos ng 5-oras na reset; binibigyan ka ng notification ng pagpipiliang lumipat na lang.
- Itinatala ng `/ccam:switch [account]` ang usapan at tinatapos ang turn bago gumawa ng anumang request: i-type ang `/exit`, saka ipagpapatuloy ito ng `ccam next --continue`, bilang account na pinangalanan mo o bilang kasunod na account.

Nasa `~/.local/state/ccam/handoff/<account>.json` ang tala nang 12 oras at isang beses lang ginagamit, kapag magsisimula na ang susunod na session; kapag pumalya ang pagsisimula, naiiwan ito para sa susunod na subok. Nakatali ito sa folder kung saan nagsimula ang session, kahit pagkatapos ng `cd` sa shell ng Claude Code. Kung napalitan na ng pangalan o natanggal na ang account na binabanggit nito, sasabihin iyon ng ccam at pipili ito gaya ng dati. Kapag may `relay = true` sa `config.toml`, nananatili ang ccam sa tabi ng mga interactive session sa halip na palitan ang sarili nito ng Claude Code, at pagkatapos ng `/ccam:switch`, kusa nitong ipinagpapatuloy ang usapan ng session na iyon (hindi nito pinapansin ang talang iniwan ng ibang session). Isang naghihintay na ccam process bawat session ang kapalit nito, kaya naka-off ito bilang default. Nagsisimula ang ipinasang session sa `ccam next --continue`, nang wala ang mga argument na ginamit sa pagsisimula ng una, gaya ng `--model`.

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
2. Isang beses binabasa ang standard input at ibinibigay sa bawat pagtatangka.
3. Sinisimulan ng ccam ang `claude` bilang child process sa halip na `exec`, pinipigil ang standard output nito, at ipinapasa ang standard error habang dumarating ito.
4. Kapag lumabas ang `claude` nang may 0, o pumalya dahil sa ibang dahilan, ipinapasa nang walang pagbabago ang pinigil na output at ang exit code.
5. Itinuturing na usage limit ang pagkabigo kapag sinasabi iyon ng mensahe ng Claude Code (`You've hit your limit`, `usage limit reached`), o kapag ipinapakita ng usage lookup ng account na nasa 100% ang 5-oras o lingguhang window nito.
6. Saka itinatapon ang output ng pagtatangkang iyon at pinapatakbo ang parehong command bilang ibang naka-log in na account: ang may pinakamababang 5-oras na usage, saka ang mga account na hindi matingnan ang usage, ayon sa pagkakasunod-sunod ng numero. Nilalaktawan ang mga account na nasa 100%, at sinusubukan ang bawat account nang hindi hihigit sa isang beses.
7. Kapag kasama sa mga argument ang `-c`, `--continue`, `-r` o `--resume`, ang mga account lang sa parehong history group ang sinusubukan, para mahanap ang usapan.
8. Ang mga account lang ng parehong organisasyon ang sinusubukan, gaya sa `--auto` (inaalis ito ng `cross_org = true`).

Pinipili ang susunod na account ayon sa parehong mga patakaran ng `--auto` (`auto_max_usage`, `auto_strategy`, ang window ng `--model`). Kapag wala nang natitirang account, ipi-print ng ccam ang output ng huling pagtatangka, sasabihing walang naka-log in na account na may natitira pa, at lalabas na may code 75. Sa `--wait` bago ang pangalan ng account (`ccam --fallback --wait=2h work -p …`), hinihintay nito sa halip na mag-reset ang unang account, saka puwedeng subukan ulit ang bawat account; lagi itong natutulog nang kahit isang beses muna, kaya hindi agad pinapatakbo ulit ang account na kaaabot lang sa limit nito. Ipi-print ng `ccam --dry-run --fallback work -p …` ang command ng unang pagtatangka.

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
- Pinipigil ng `--fallback` ang standard output hanggang matapos ang isang pagtatangka, kaya dumarating nang buo sa dulo ang streaming output gaya ng `--output-format stream-json`. Nananatiling binago ang mga file na binago ng unang pagtatangka bago nito naabot ang limit; sinisimulan ulit ng susunod na pagtatangka ang prompt.

## Kaugnay

- [Mga usage limit](/features/usage)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Workflow kapag naabot ang usage limit](/guides/usage-limit-workflow)
