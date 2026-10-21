---
title: Pagsusuri at pagkukumpuni
description: Suriin ang mga share link, permission ng file, login at bersyon ng Claude Code gamit ang ccam doctor, at ayusin ang ligtas ayusin gamit ang ccam doctor --fix.
updated: 2026-10-21
order: 9
translation: machine
sourceHash: "9be63350cb69"
---

Sinusuri ng `ccam doctor` ang iyong ccam installation: ang Claude Code binary, ang folder, login at mga share link ng bawat account, ang mga shared na MCP server, ang usage endpoint at ang iyong mga shortcut command. Kinukumpuni ng `ccam doctor --fix` ang kaya nitong ayusin nang walang nawawalang data, at inililipat muna sa backup folder ang anumang file na papalitan nito. Patakbuhin ito kapag huminto ang `ccam <name>` dahil sa problema sa share link, pagkatapos mong mano-manong ibalik ang mga file, o pagkatapos ng update ng Claude Code.

## Paggamit

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Minamarkahan ng unang linya ang bawat account nang sabay, kaya makikita mo kung nasaan ang mga problema bago ka magbasa ng iba pa. Sa ilalim nito, may sariling kahon ang bawat account, at napupunta sa kahong `General` ang mga natuklasang hindi kabilang sa anumang account (Claude Code, mga environment variable, mga plugin, mga MCP server, mga shortcut, mga patakaran ng folder, mga group at mga setting). Ipinapakita ng pamagat ng kahon ang pinakamalalang marka ng account at kung ilan ang problema nito. Ang bawat hilera sa loob ay isang bahaging sinuri, laging `Folder`, `Login` at `Share links` sa ganoong pagkakasunod-sunod, saka ang iba pa, na may marka: `✔` OK, `ℹ` impormasyon, `⚠` babala, `✘` error. Nagpapakita ng `✔ Folder  OK` ang folder na walang problema.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Sa `--fix`, nagtatapos sa `→ fixed` ang bawat natuklasang nakumpuni, at hindi na ito binibilang bilang problema:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

Sa terminal, kasya ang mga kahon sa lapad nito at nababalot ang mahahabang mensahe, sa sariling mga kulay ng terminal (pinapatay ng `NO_COLOR` ang mga ito). Kapag ipinadala sa pipe o file, walang kulay ang output at walang mensaheng binabalot o pinuputol, kaya nahahanap ng `ccam doctor | grep settings.json` ang buong linya; kasinlapad noon ng pinakamahabang linya ang bawat kahon. Ipinapakita ng Health tab ng dashboard ang parehong mga kahon.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--fix` | Kumpunihin ang ligtas kumpunihin (bina-back up ang mga binagong file) | off |
| `--dry-run` (global) | Patakbuhin ang bawat pagsusuri pero walang kukumpunihin, kahit may `--fix` | off |

Walang tinatanggap na argument ang `ccam doctor`. Sinusunod nito ang `--dry-run`: nagpi-print ang `ccam --dry-run doctor --fix` ng parehong mga natuklasan gaya ng `ccam doctor` at walang binabago.

`1` ang exit code kapag may kahit isang natuklasang `✘` na hindi naayos, at `0` kung wala. Hindi kailanman nagpapapalya ng run ang mga babala at impormasyon. Kapag pumalya ang run nang walang `--fix`, nagtatapos ang ccam sa isang blangkong linya at `` `ccam doctor --fix` repairs what it can. ``. Nagtatapos din sa `1` at isang `ccam:` error ang config na hindi ma-load.

## Paano ito gumagana

Nilo-load ng `internal/cli/doctor.go` ang `~/.config/ccam/config.toml` at tinatawag ang `doctor.Run` sa `internal/doctor/doctor.go`, na nagpapatakbo ng mga pagsusuri sa nakapirming pagkakasunod-sunod at nagbabalik ng mga natuklasan, bawat isa kasama ang account na kinabibilangan nito at ang bahaging sinuri nito. Pinapangkat ng `tui.HealthReport` ang mga ito sa mga kahon. Ipinapasa ang `--fix` bilang `Options.Fix`, na false sa dry run.

<!-- motion:doctor-checks -->

Isang beses bawat account tumatakbo ang gitnang bahagi. Ang mga pagsusuri, at ang ginagawa ng `--fix` sa bawat isa:

1. **Bersyon ng Claude Code** (`checkClaude`). Hinahanap ang `claude` binary gaya ng ginagawa ng launcher at pinapatakbo ang `claude --version`. Kung hindi nagsisimula sa `2.1.302` ang bersyon, ang bersyong pinagbatayan ng pag-verify sa folder layout ng ccam (`doctor.LayoutVersion`), idinadagdag nito ang impormasyong `ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Error ang nawawalang binary o pumapalyang `--version`. Walang fix.
2. **Mga credential sa environment** (`checkAuthEnv`). Nagbababala kapag nakatakda ang `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: minamana ito ng bawat account na sinimulan mula sa shell na iyon, at mas inuuna ito ng Claude Code kaysa sa sariling login ng account. Walang fix; ikaw mismo ang mag-unset nito.
3. **Folder** (`checkProfile`). Error ang nawawalang folder ng account at tinatapos nito ang mga pagsusuri ng account na iyon. Babala ang mga permission na may nakatakdang group o other bit; pinapatakbo ng `--fix` ang `chmod 0700`. Walang permission bit ang Windows, kaya nilalaktawan doon ang pagsusuring ito at ang pagsusuri ng `0600` sa ibaba; pinoprotektahan ng access list nito ang folder sa user profile mo.
4. **Login file** (`checkCredentials`). Sinusuri ang `.credentials.json` gamit ang `lstat`. Error ang symlink, dahil dapat may sariling token ang bawat account, at hindi ito kinukumpuni. Babala ang mode na hindi `0600`; pinapatakbo ng `--fix` ang `chmod 0600`. Pagkatapos ay binabasa ang login (sa macOS, mula muna sa Keychain). Babala ang walang login at itinuturo nito ang `ccam login <name>`. Ang expired na refresh token ay ang error na `login expired`; impormasyon lang ang expired na access token, dahil nire-refresh ito ng Claude Code sa susunod na pagsisimula. Ipinapakita ng valid na login ang email at ang plan (`login valid (dev@company.com, Claude Team · Acme)`). Wala sa mga ito ang kinukumpuni; ina-log in ng `ccam login --expired` ang bawat account na nangangailangan nito.
5. **Mga share link** (`share.Check`). Ikinukumpara ang bawat path sa `share.Layout` sa mga group ng account, gamit lang ang `lstat` at `readlink`. Apat na uri ang iniuulat nito: nawawalang link, link na tumuturo sa ibang lugar, totoong file kung saan dapat may link, at link na naiwan para sa item na local na ngayon. Babala ang nawawalang link, error ang iba. Ang mga path na idinagdag sa isang item sa mas bagong bersyon ng ccam (`file-history/`, `tasks/` at `paste-cache/` para sa history, `keybindings.json` para sa settings) ay babala lang habang may sariling kopya pa ang account: `is not shared with the group yet (ccam doctor --fix moves it in)`. Hindi kailanman humihinto ang isang launch dahil sa mga ito.
6. **Mga hindi kilalang entry** (`checkUnknownEntries`). Inililista ang mga pangalan sa itaas na antas ng folder ng account na hindi kilala ng ccam (`share.Known`), bilang pahiwatig na nagbago ang layout ng Claude Code. Impormasyon lang.
7. **Mga plugin path** (`checkPlugins`). Nag-iimbak ng mga absolute na path ang mga plugin registry (`installed_plugins.json`, `known_marketplaces.json`). Sa `plugins/` ng bawat group, nasisira ang mga path na dumadaan sa folder ng isang account kapag inalis o pinalitan ang pangalan ng account na iyon. Isinusulat ulit ng `--fix` ang mga ito papunta sa group folder gamit ang `profile.RewritePluginPaths`, isang atomic na pagsulat ng registry file.
8. **Mga MCP server** (`checkMCP`). Nilo-load nang isang beses ang `mcp.json` ng bawat group. Error ang file na hindi mabasa; babala ang server na pumapalya sa `mcp.Validate`. Walang fix.
9. **Format ng usage response** (`checkUsageFormat`). Nagtatanong sa usage endpoint, nang walang cache at may karaniwang 3 segundong limit, para sa isang account kasunod ng isa pa hanggang may sumagot. Kapag hindi na ma-parse ang response, nagbababala ito na babalik ang `ccam usage` sa status line cache. Ito lang ang pagsusuring gumagamit ng network.
10. **Mga shortcut** (`checkShims`). Kapag may shortcut lang ang isang account. Babala ang mga nawawalang link sa `~/.local/bin`; tinatawag ng `--fix` ang `shim.Install`, na gumagawa lang ng mga nawawalang link at hindi kailanman ginagalaw ang ibang file. Babala ang bin folder na wala sa `PATH` at hindi ito kinukumpuni.
11. **Mga folder** (`checkFolders`). Kapag may mga entry lang ng [`ccam folder`](/features/folder). Error ang entry na hindi umiiral ang account, at babala ang entry na wala na ang folder. Hindi kinukumpuni ang alinman, dahil mababago ng pag-alis ng entry kung aling account ang nagsisimula roon.
12. **Mga group at setting** (`checkGroups`). Babala ang miyembro ng group na hindi account, ang account na naging ccam command na ang pangalan (simulan ito gamit ang `ccam run <name>` o palitan ang pangalan nito), at ang hindi kilalang `theme` o `usage_source` sa `config.toml`. Walang fix.
13. **Mga login** (`checkIdentities`). Binabasa mula sa `.claude.json` ng bawat folder kung kanino ito naka-log in. Babala ang folder na naka-log in na ngayon sa ibang account kaysa sa itinala ng `ccam login` (halimbawa pagkatapos ng `/login` sa loob ng isang session), at babala rin ang folder na naka-log in sa parehong Claude account ng isa pa: iisa ang usage limit ng dalawa. Walang fix.

### Pagkukumpuni ng mga share link

Ginagawa agad ng launch path ang nawawalang link, pero humihinto ito sa anumang puwedeng magpawala ng data at itinuturo ka rito:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Para sa anumang problema sa share link, tinatawag ng `--fix` ang `share.Build` na ang share settings ng account ang kasalukuyan at target na estado, saka ang `Plan.Apply`. Hinuhango ulit ng `Build` ang bawat hakbang mula sa nasa disk, ang parehong plano na ginagamit ng pagbabago sa `ccam group`, kaya sinusunod ng pagkukumpuni ang parehong mga patakaran ng pagsali sa isang group.

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
- Para sa account na tumatakbo, gumagawa lang ang `doctor --fix` ng mga nawawalang link at idinadagdag ang `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, dahil patuloy na nagsusulat ang session sa mga file na ililipat ng pagkukumpuni.
- Kapag pumalya ang pagkukumpuni ng share, iniuulat ng ccam ang `could not repair share links: <error>` para sa account na iyon at nilalaktawan ang pagsusuri nito ng mga hindi kilalang entry. Sinusuri pa rin ang ibang account.
- Kapag pumalya ang isang pagbabago sa `ccam group` sa kalagitnaan, naka-save na ang mga bagong group, at sasabihin ng ccam na tapusin mo gamit ang `ccam doctor --fix`.
- Iniuulat pero hindi kailanman binabago ang naka-symlink na `.credentials.json`, ang expired na login at ang mga credential sa environment.
- Atomic na isinusulat ang pagsulat ulit ng plugin registry pero hindi kinokopya sa `.ccam-backup/`.

## Kaugnay

- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Pagkopya ng dati nang setup](/features/migrate)
- [Mga usage limit at paglipat ng account](/docs/usage-limits)
- [Mga command](/docs/commands)
