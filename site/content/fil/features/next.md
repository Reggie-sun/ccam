---
title: Paglipat sa susunod na account
description: Binubuksan ng ccam next ang Claude Code bilang account na kasunod ng huli mong ginamit; pinipili ng --auto ang pinakamababang 5-oras na usage at itinutuloy ng --continue ang usapan.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

Binubuksan ng `ccam next` ang Claude Code bilang ibang account nang hindi mo ito pinapangalanan. Gamitin ito kapag naabot ng isang account ang usage limit nito: kinukuha ng simpleng `ccam next` ang account na kasunod ng huli mong ginamit, kinukuha ng `--auto` ang may pinakamaraming natitira sa 5-oras na window nito, at itinutuloy ng `--continue` ang huling usapan sa bagong account.

## Paggamit

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

Ipi-print ng ccam ang paglipat bilang `→ <from> → <to>` sa standard error, saka sisimulan ang Claude Code bilang ang bagong account sa parehong paraan ng `ccam <name>`. Ang `<from>` ay ang account na huli mong ginamit. Ipinapasa sa `claude` nang walang pagbabago ang mga argument na kasunod ng `--`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--auto` | Piliin ang naka-log in na account na may pinakamababang 5-oras na usage, at laktawan ang mga account na nasa 100% | off |
| `-c`, `--continue` | Isaalang-alang lang ang mga account na nagshe-share ng kasaysayan ng usapan sa account na huli mong ginamit, at simulan ang `claude --continue` | off |
| `-- <claude args>` | Mga argument na ipinapasa sa `claude` pagkatapos ng sariling argument ng ccam | wala |
| `--dry-run` | Sinusunod: ipi-print ang `→ <from> → <to>` at ang command line na patatakbuhin nito, nang hindi sinisimulan ang Claude Code o itinatala ang pagsisimula | off |

Gumagana nang magkasama ang `--auto` at `--continue`. Sa `--continue`, inilalagay ang `--continue` bago ang sarili mong mga argument.

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

Nagpapasya ang `--auto` ayon sa usage sa halip na sa pagkakasunod-sunod ng slot. Pinapanatili ng `cli.pickAuto` ang mga naka-log in na account, sabay-sabay na hinihingi ang usage nilang lahat gamit ang parehong `Fetcher.FetchAll` na ginagamit ng `ccam usage`, at hinahayaang pumili ang `usage.Best()`:

<!-- motion:next-auto -->

Kinukuha ng `usage.Best()` ang pinakamababang 5-oras na porsiyento at nilalaktawan ang mga account na nasa 100% o hindi alam ang 5-oras na usage. Sinusunod ng paghahanap ang mga patakaran ng `ccam usage`: muling ginagamit ang mga numerong wala pang 60 segundo, may 3 segundo ang bawat account, at ang account na hindi sumasagot, o nag-expire ang login, ay hinuhusgahan ayon sa mga naka-save na numero nito. Itinuturing na 0% ang naka-save na window na lumipas na ang oras ng reset, kaya hindi nalalampasan ang account na nag-reset na ang window mula noon. Hindi kailanman nire-renew ng `--auto` ang mga token. Hindi tulad ng simpleng rotation, puwedeng piliin ng `--auto` ang account na huli mong ginamit kung ito pa rin ang may pinakamaraming natitira.

### --continue

Binubuksan ng `--continue` ng Claude Code ang pinakahuling usapan sa kasalukuyang folder, na gumagana lang kung nakikita ng bagong account ang usapang iyon. Kaya sa `--continue`, isinasaalang-alang lang ng ccam ang mga account na nasa parehong `history` group ng account na huli mong ginamit, saka inilalagay ang `--continue` sa unahan ng iyong mga argument. Walang pinagshe-share-an ang mga account na `local` ang kasaysayan, kaya laging pumapalya ang `--continue` mula sa ganitong account. I-on muna ang sharing gamit ang `ccam share <name> --history on` (tingnan ang [sharing](/docs/sharing)).

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman pinipili ang mga account na walang login. Kapag walang ibang naka-log in na account, humihinto ang simpleng `ccam next` na may `no logged-in account other than <name>`.
- Sa `--continue` at walang account na nagshe-share ng kasaysayan: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Sa `--auto` at walang naka-log in na account: `no logged-in account to pick`. Kapag puno ang bawat account o hindi alam ang usage nito: `` no account with known usage and headroom (see `ccam usage`) ``.
- Gumagawa ng mga network request ang `--auto` kahit may `--dry-run`. Nagbabasa lang ito, at hindi kailanman sumusulat sa mga token file.
- Nagsisimula ang napiling account sa parehong daan ng `ccam <name>`, kaya nalalapat ang parehong mga pagsusuri, kasama ang babala kapag may tumatakbo nang ibang Claude Code na nagshe-share ng parehong kasaysayan sa folder na ito.
- Isinusulat lang ang `last-used` kapag talagang nagsimula ang Claude Code. Hindi ito ginagalaw ng `--dry-run`, kaya hindi isinusulong ng dry run ang rotation.

## Kaugnay

- [Mga usage limit](/features/usage)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Workflow kapag naabot ang usage limit](/guides/usage-limit-workflow)
