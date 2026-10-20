---
title: Pagsusuri at pagkukumpuni
description: Suriin ang mga share link, permission ng file, login at bersyon ng Claude Code gamit ang ccam doctor, at ayusin ang ligtas ayusin gamit ang ccam doctor --fix.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

Sinusuri ng `ccam doctor` ang iyong ccam installation: ang Claude Code binary, ang folder, login at mga share link ng bawat account, ang mga shared na MCP server, ang usage endpoint at ang iyong mga shortcut command. Kinukumpuni ng `ccam doctor --fix` ang kaya nitong ayusin nang walang nawawalang data, at inililipat muna sa backup folder ang anumang file na papalitan nito. Patakbuhin ito kapag huminto ang `ccam <name>` dahil sa problema sa share link, pagkatapos mong mano-manong ibalik ang mga file, o pagkatapos ng update ng Claude Code.

## Paggamit

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Isang linya na may marka ang bawat natuklasan: `✓` OK, `·` impormasyon, `!` babala, `✗` error. Nagsisimula sa pangalan ng account ang mga natuklasan tungkol sa account.

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

Sa `--fix`, nagtatapos sa `→ fixed` ang bawat natuklasang nakumpuni:

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

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--fix` | Kumpunihin ang ligtas kumpunihin (bina-back up ang mga binagong file) | off |
| `--dry-run` (global) | Patakbuhin ang bawat pagsusuri pero walang kukumpunihin, kahit may `--fix` | off |

Walang tinatanggap na argument ang `ccam doctor`. Sinusunod nito ang `--dry-run`: nagpi-print ang `ccam --dry-run doctor --fix` ng parehong mga natuklasan gaya ng `ccam doctor` at walang binabago.

`1` ang exit code kapag may kahit isang natuklasang `✗` na hindi naayos, at `0` kung wala. Hindi kailanman nagpapapalya ng run ang mga babala at impormasyon. Kapag pumalya ang run nang walang `--fix`, nagtatapos ang ccam sa isang blangkong linya at `` `ccam doctor --fix` repairs what it can. ``. Nagtatapos din sa `1` at isang `ccam:` error ang config na hindi ma-load.

## Paano ito gumagana

Nilo-load ng `internal/cli/doctor.go` ang `~/.config/ccam/config.toml` at tinatawag ang `doctor.Run` sa `internal/doctor/doctor.go`, na nagpapatakbo ng mga pagsusuri sa nakapirming pagkakasunod-sunod at nagbabalik ng mga natuklasan. Ipinapasa ang `--fix` bilang `Options.Fix`, na false sa dry run.

<!-- motion:doctor-checks -->

Isang beses bawat account tumatakbo ang gitnang bahagi. Ang mga pagsusuri, at ang ginagawa ng `--fix` sa bawat isa:

1. **Bersyon ng Claude Code** (`checkClaude`). Hinahanap ang `claude` binary gaya ng ginagawa ng launcher at pinapatakbo ang `claude --version`. Kung hindi nagsisimula sa `2.1.302` ang bersyon, ang bersyong pinagbatayan ng pag-verify sa folder layout ng ccam (`doctor.LayoutVersion`), idinadagdag nito ang `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Error ang nawawalang binary o pumapalyang `--version`. Walang fix.
2. **Mga credential sa environment** (`checkAuthEnv`). Nagbababala kapag nakatakda ang `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: minamana ito ng bawat account na sinimulan mula sa shell na iyon, at mas inuuna ito ng Claude Code kaysa sa sariling login ng account. Walang fix; ikaw mismo ang mag-unset nito.
3. **Folder** (`checkProfile`). Error ang nawawalang folder ng account at tinatapos nito ang mga pagsusuri ng account na iyon. Babala ang mga permission na may nakatakdang group o other bit; pinapatakbo ng `--fix` ang `chmod 0700`. Walang permission bit ang Windows, kaya nilalaktawan doon ang pagsusuring ito at ang pagsusuri ng `0600` sa ibaba; pinoprotektahan ng access list nito ang folder sa user profile mo.
4. **Login file** (`checkCredentials`). Sinusuri ang `.credentials.json` gamit ang `lstat`. Error ang symlink, dahil dapat may sariling token ang bawat account, at hindi ito kinukumpuni. Babala ang mode na hindi `0600`; pinapatakbo ng `--fix` ang `chmod 0600`. Pagkatapos ay binabasa ang login (sa macOS, mula muna sa Keychain). Babala ang walang login at itinuturo nito ang `ccam login <name>`. Ang expired na refresh token ay ang error na `login expired`; impormasyon lang ang expired na access token, dahil nire-refresh ito ng Claude Code sa susunod na pagsisimula. Wala sa mga ito ang kinukumpuni.
5. **Mga share link** (`share.Check`). Ikinukumpara ang bawat path sa `share.Layout` sa share settings ng account, gamit lang ang `lstat` at `readlink`. Apat na uri ang iniuulat nito: nawawalang link, link na tumuturo sa ibang lugar, totoong file kung saan dapat may link, at link na naiwan para sa item na local na ngayon. Babala ang nawawalang link, error ang iba.
6. **Mga hindi kilalang entry** (`checkUnknownEntries`). Inililista ang mga pangalan sa itaas na antas ng folder ng account na hindi kilala ng ccam (`share.Known`), bilang pahiwatig na nagbago ang layout ng Claude Code. Impormasyon lang.
7. **Mga plugin path** (`checkPlugins`). Nag-iimbak ng mga absolute na path ang mga plugin registry (`installed_plugins.json`, `known_marketplaces.json`). Sa `plugins/` ng bawat group, nasisira ang mga path na dumadaan sa folder ng isang account kapag inalis o pinalitan ang pangalan ng account na iyon. Isinusulat ulit ng `--fix` ang mga ito papunta sa group folder gamit ang `profile.RewritePluginPaths`, isang atomic na pagsulat ng registry file.
8. **Mga MCP server** (`checkMCP`). Nilo-load nang isang beses ang `mcp.json` ng bawat group. Error ang file na hindi mabasa; babala ang server na pumapalya sa `mcp.Validate`. Walang fix.
9. **Format ng usage response** (`checkUsageFormat`). Nagtatanong sa usage endpoint, nang walang cache at may karaniwang 3 segundong limit, para sa isang account kasunod ng isa pa hanggang may sumagot. Kapag hindi na ma-parse ang response, nagbababala ito na babalik ang `ccam usage` sa status line cache. Ito lang ang pagsusuring gumagamit ng network.
10. **Mga shortcut** (`checkShims`). Kapag may shortcut lang ang isang account. Babala ang mga nawawalang link sa `~/.local/bin`; tinatawag ng `--fix` ang `shim.Install`, na gumagawa lang ng mga nawawalang link at hindi kailanman ginagalaw ang ibang file. Babala ang bin folder na wala sa `PATH` at hindi ito kinukumpuni.
11. **Mga folder** (`checkFolders`). Kapag may mga entry lang ng [`ccam folder`](/features/folder). Error ang entry na hindi umiiral ang account, at babala ang entry na wala na ang folder. Hindi kinukumpuni ang alinman, dahil mababago ng pag-alis ng entry kung aling account ang nagsisimula roon.

### Pagkukumpuni ng mga share link

Ginagawa agad ng launch path ang nawawalang link, pero humihinto ito sa anumang puwedeng magpawala ng data at itinuturo ka rito:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Para sa anumang problema sa share link, tinatawag ng `--fix` ang `share.Build` na ang share settings ng account ang kasalukuyan at target na estado, saka ang `Plan.Apply`. Hinuhango ulit ng `Build` ang bawat hakbang mula sa nasa disk, ang parehong plano na ginagamit ng `ccam share`, kaya sinusunod ng pagkukumpuni ang parehong mga patakaran ng pag-on ng sharing.

<!-- motion:doctor-repair -->

Dito, totoong file ang `settings.json`. Dahil may ibang `settings.json` na ang group, kinukuha ng conflict ang default na pagpipilian nito, ang `KeepGroup`, kaya nananatili ang kopya ng group. Lumilipat ang file ng account sa `.ccam-backup/<time>/settings.json` sa loob ng folder ng account, at nili-link ng `os.Symlink` ang path sa group. Iba pang kaso:

- Kung wala pang kopya ang group, lumilipat ang file ng account sa group at nagiging unang nilalaman nito.
- Para sa folder, lumilipat sa group ang mga file na wala sa group; napupunta sa backup ang iba.
- Para sa `history.jsonl`, `merge by line` ang default: ang mga linya ng group, saka ang mga bago mong linya. Sine-save ang dating kopya ng group sa `.ccam-backup/<time>/group/`.
- Pinapalitan ang link sa ibang group. Napupunta sa backup ang link sa folder na nasa labas ng ccam, at ang nilalaman nito ang nagiging simula ng group kung walang laman ang group.
- Inaalis ang link na naiwan para sa item na local na ngayon, at nakakakuha ang account ng kopya ng nilalaman ng group.

Bukod sa pag-alis ng mga link, inililipat ng bawat hakbang ang data sa group o sa backup folder. Walang binubura.

## Kaligtasan at mga espesyal na kaso

- Walang binabago kailanman ang simpleng `ccam doctor`. Gamitin ang `--dry-run` kasama ng `--fix` para makasiguro sa mga script.
- Hindi sinusuri ng `doctor --fix` kung tumatakbo ang isang account. Isara ang Claude Code para sa isang account bago kumpunihin ang mga share link nito, dahil patuloy na nagsusulat ang session sa mga file na iyon.
- Kapag pumalya ang pagkukumpuni ng share, iniuulat ng ccam ang `could not repair share links: <error>` para sa account na iyon at nilalaktawan ang pagsusuri nito ng mga hindi kilalang entry. Sinusuri pa rin ang ibang account.
- Kapag pumalya ang `ccam share` sa kalagitnaan, sine-save muna nito ang target na settings at sasabihing tapusin mo gamit ang `ccam doctor --fix`.
- Iniuulat pero hindi kailanman binabago ang naka-symlink na `.credentials.json`, ang expired na login at ang mga credential sa environment.
- Atomic na isinusulat ang pagsulat ulit ng plugin registry pero hindi kinokopya sa `.ccam-backup/`.

## Kaugnay

- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Pagkopya ng dati nang setup](/features/migrate)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [Mga command](/docs/commands)
