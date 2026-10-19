---
title: Sharing sa pagitan ng mga account
description: Paano nili-link ng ccam share ang settings, CLAUDE.md, mga extension at kasaysayan sa isang group folder, na may plano, backup at mga pagpipilian sa conflict, at paano kumokopya pabalik ang off.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

May sariling folder ang bawat account, pero hindi kailangang hiwalay ang karamihan ng laman nito. Ino-on o ino-off ng `ccam share` ang sharing para sa isang account at isang item sa bawat pagkakataon: inililipat ang file o folder ng account sa isang group folder at pinapalitan ng symbolic link, kaya binabasa at sinusulatan ng bawat account sa group ang iisang kopya. Gamitin ito pagkatapos ng `ccam add` para baguhin kung ano ang shine-share ng isang account, o para ilagay ang ilang account sa sarili nilang group.

## Paggamit

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Kapag walang item option, ipi-print ng `ccam share` ang kasalukuyang estado:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Kapag may mga option, nagpi-print ang ccam ng plano, nagtatanong tungkol sa bawat file na naiiba sa kopya ng group, at nagtatanong ng `Proceed? [y/N]` bago ito magbago ng anuman. Shine-share ng bagong account ang lahat maliban sa kasaysayan ng usapan; inilalarawan ng [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing) ang bawat item mula sa panig ng user.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Mga MCP server | hindi binabago |
| `--settings <on\|off\|group>` | `settings.json` | hindi binabago |
| `--memory <on\|off\|group>` | ang global na `CLAUDE.md` | hindi binabago |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | hindi binabago |
| `--history <on\|off\|group>` | `projects/` at `history.jsonl` | hindi binabago |
| `--fresh` | kapag ino-off ang sharing, magsimulang walang laman sa halip na kopyahin ang nilalaman ng group | off |
| `--prefer <group\|mine\|merge>` | paano pagpapasyahan ang mga conflict | magtanong; sa `--yes`, ang default ng bawat item |
| `-y`, `--yes` | tumakbo nang hindi nagtatanong; sinusunod ng mga conflict ang `--prefer` | off |
| `-h`, `--help` | help para sa `share` | |
| `--dry-run` (global) | sinusunod: nagtatanong tungkol sa mga conflict, nagpi-print ng plano, walang binabago | off |

Ang ibig sabihin ng value na `on` ay ang group na `default`, ang `off` ay local (hindi shine-share), at pangalan ng group ang anumang ibang value. Nagiging pangalan ng folder sa ilalim ng `~/.local/share/ccam/groups/` ang mga pangalan ng group, kaya pumapalya ang `ccam share side --history ../x` na may `invalid group name: "../x"`. Ibinibigay ang account ayon sa pangalan o numero, at ang mga item lang na pinangalanan mo ang nagbabago.

## Paano ito gumagana

Ang mga file na kasangkot:

- `~/.config/ccam/config.toml`: ang table na `[profile.share]` ng bawat account, isang value bawat item.
- `~/.local/share/ccam/profiles/<name>/`: ang `CLAUDE_CONFIG_DIR` ng account.
- `~/.local/share/ccam/groups/<group>/`: ang mga shared na kopya.

1. Binabasa ng `internal/cli/share.go` ang config, inilalapat ang bawat value ng `--<item>` sa kopya ng settings ng account gamit ang `config.Share.Set`, at tumatanggi kung tumatakbo ang account.
2. Ikinukumpara ng `share.Build` ang kasalukuyang settings sa target at inililista ang mga hakbang. Binabasa lang nito ang disk: `os.Lstat` at `os.Readlink` para sa bawat path sa `share.Layout`. Listahan ng mga operasyon ang plano: ilipat sa bakanteng puwesto sa group (adopt), ilipat ang mga file na wala sa group (merge dir), isulat sa group ang naresolbang conflict, mag-back up, mag-unlink, mag-link, kumopya pabalik, isulat ulit ang mga plugin path at pagsamahin ang plugin registry.
3. Pinagpapasyahan ang mga conflict, sa pamamagitan ng tanong o ng `--prefer`.
4. Ipi-print ng `Plan.Describe` ang plano, at magtatanong ang ccam.
5. Isinusulat ang target na settings sa `config.toml` bago may ilipat. Kapag huminto ang `Plan.Apply` sa kalagitnaan, binubuo ng `ccam doctor --fix` at ng muling pagpapatakbo ang kanilang plano patungo sa naka-save na target at tinatapos ang pagbabago.
6. Pinapatakbo ng `Plan.Apply` ang mga hakbang nang sunod-sunod. Inililipat ng bawat hakbang ang data sa group o sa backup folder; walang binubura.

Ang `share.Layout` ay listahan ng mga pinapayagang path na sakop ng bawat item. Hindi kailanman makakaabot sa group ang anumang wala rito. Pinapangalanan ng `share.Never` ang `.credentials.json` at `.claude.json`, at humihinto ang `Build` na may internal error kung lumabas man ang alinman sa layout. Walang path ang item na `mcp`: nasa loob ng `.claude.json` ang listahan ng server, kaya gumagana ang [MCP sharing](/features/mcp) sa pamamagitan ng `mcp.json` ng group sa halip, at binabago lang ng `--mcp` ang setting.

Ito ang nangyayari sa isang folder, ang `projects/`, kapag ino-on ng `side` ang kasaysayan at may laman na ang group:

<!-- motion:share-on -->

Kung wala pang laman ang group sa path na iyon, ang account ang unang nagshe-share nito at inililipat doon ang folder nito nang walang pagbabago, nang walang backup. Kung hindi, inililipat ng `mergeDir` sa group ang bawat file na wala sa group, at nananatili sa folder ng account ang mga file na mayroon sa magkabilang panig, kaya pinapanatili ang bersyon ng group. Pagkatapos ay napupunta ang folder ng account sa `.ccam-backup/<time>/` sa loob ng folder ng account, at naglalagay ang `os.Symlink` ng link sa puwesto nito. Ipinapangalan ang mga file ng usapan ayon sa kanilang session ID, kaya nagsasama ang mga kasaysayan mula sa ilang account nang walang banggaan.

Gumagamit ang paglilipat ng `rename(2)` at bumabalik sa copy at delete kapag magkaibang file system. Ang relative na link sa loob ng inilipat na tree na tumuturo sa labas nito ay isinusulat ulit bilang absolute na link (`fsutil.PortableLink`), para gumana pa rin ito mula sa bagong lugar. Nag-iimbak ang mga registry ng `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) ng mga absolute na path na dumadaan sa folder ng account; itinuturo ng ccam ang mga ito sa folder ng group at, kapag may registry na ang group, pinagsasama roon ang mga entry ng account.

Binabaligtad ng pag-off ng isang item ang link:

<!-- motion:share-off -->

Inaalis lang ng ccam ang link kung tumuturo ito sa groups folder, saka kinokopya ang kasalukuyang nilalaman ng group sa account gamit ang `CopyTree`. Sa `--fresh`, nilalaktawan ang pagkopya at nagsisimula ang account na walang laman sa path na iyon. Pinapanatili ng group ang nilalaman nito para sa ibang account. Kapag inilipat ang isang item mula sa isang group papunta sa iba, ina-unlink ito at nili-link sa bagong group, at hindi ginagalaw ang lumang group; sinasabi iyon ng plano sa linyang `Note:`.

Conflict ang iisang file na mayroon sa magkabilang panig na magkaiba ang nilalaman:

<!-- motion:share-conflict -->

Inililista ng `share.NewConflict` ang mga pagpipiliang angkop sa file. Nag-aalok ang `history.jsonl` ng `merge by line`, `keep group` at `use mine`, na nauuna ang pagsasama. Nag-aalok ang ibang `.json` file ng `keep group`, `use mine` at `merge JSON`. Nag-aalok ang anumang iba ng `keep group` at `use mine`. Ang unang pagpipilian ang default:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

Pinagsasama ng `merge JSON` (`MergeJSONBytes`) ang mga object nang recursive, pinagdudugtong ang mga array nang walang duplicate, at pinapanatili ang value ng group kapag magkaiba ang itinakda ng magkabilang panig sa parehong key, kaya walang nakikitang pagbabago sa kilos ang ibang account. Pinapanatili ng `merge by line` ang mga linya ng group, saka idinadagdag ang mga linya mo na wala roon. Bawat pagpipilian maliban sa `keep group` ay sine-save muna ang kopya ng group sa `.ccam-backup/<time>/group/` sa folder ng account, saka isinusulat ang resulta gamit ang `fsutil.AtomicWrite`. Laging napupunta sa backup ang sarili mong file.

Pinagpapasyahan ng `--prefer group`, `mine` o `merge` ang bawat conflict nang hindi nagtatanong; ang ibig sabihin ng `merge` ay `merge JSON` o `merge by line`, alinman ang inaalok ng file, at pinapanatili ang default ng file na hindi nag-aalok ng alinman. Kinukuha ng `--yes` lang ang default ng bawat file.

Sa bawat pagsisimula, ikinukumpara ng `share.Check` ang bawat path sa layout sa setting gamit lang ang `lstat` at `readlink`, kaya sapat itong magaan para sa launch path. Ginagawa agad ang nawawalang link, dahil walang mawawala. Pinapahinto ang pagsisimula ng totoong file kung saan dapat may link, ng link sa maling lugar, o ng lumang link:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Pinapatakbo ng `ccam doctor --fix` ang `share.Build` na ang naka-save na settings ang nasa magkabilang panig, na hinahango ang bawat hakbang mula sa nasa disk, pinapanatili ang nilalaman ng group sa mga conflict at bina-back up ang anumang pinapalitan nito.

## Kaligtasan at mga espesyal na kaso

- **Mga tumatakbong account.** Tumatanggi ang ccam na baguhin ang sharing habang tumatakbo ang Claude Code bilang ang account (`account work is running; quit it first, then change its sharing`), dahil patuloy na nagsusulat ang session sa mga file na inililipat at mawawala ang mga linyang idinadagdag nito habang nangyayari iyon. Ipi-print pa rin ng `--dry-run` ang plano.
- **Mga backup.** May sariling `.ccam-backup/<time>/` folder ang bawat plano sa folder ng account, na may numerong suffix kapag nagsimula ang dalawang plano sa parehong segundo. Hindi kailanman binubura ng ccam ang mga ito.
- **Bahagyang pagpalya.** Kapag pumalya ang isang hakbang, ipi-print ng ccam ang error at kung nasaan ang mga orihinal, at sasabihing tapusin mo gamit ang `ccam doctor --fix`. Naka-save na ang target, kaya nagpapatuloy ang pagkukumpuni patungo roon.
- **Sarili mong mga link.** Kung link na ginawa mo na ang isang path, halimbawa sa isang dotfiles repository, at walang laman ang group, kinokopya ng ccam ang target ng link sa group at inililipat ang link mismo sa backup. Hindi ginagalaw ang target.
- **Mga login file.** Hindi kailanman nili-link, inililipat o kinokopya sa group ang `.credentials.json` at `.claude.json`.
- **Dalawang account sa isang folder.** Sinusulatan ng mga account na nagshe-share ng kasaysayan ang parehong mga file. Kapag nagsimula ang pangalawang account ng parehong history group sa folder kung saan may tumatakbo na, nagbababala ang ccam at nagtatanong bago magpatuloy.

## Kaugnay

- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Mga shared na MCP server](/features/mcp)
- [Kopyahin ang dati nang setup](/features/migrate)
- [Lahat ng command](/docs/commands)
