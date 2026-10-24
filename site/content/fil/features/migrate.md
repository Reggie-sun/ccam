---
title: Pagkopya ng dati nang setup
description: Kopyahin ang settings, CLAUDE.md, mga extension at MCP server ng dati nang ~/.claude papunta sa isang ccam group gamit ang ccam migrate.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: 7bca9efb6eb0
---

Kung gumamit ka ng Claude Code bago ang ccam, nasa `~/.claude` at `~/.claude.json` ang iyong setup. Kinokopya ito ng `ccam migrate` sa isang share group, kaya nagsisimula ang bawat account na nagshe-share ng mga item na iyon sa pamamagitan ng group na dala ang setup na mayroon ka na. Binabasa lang ang source folder, at gumagana pa rin ang karaniwang `claude` gaya ng dati. Patakbuhin ito nang isang beses pagkatapos ng una mong `ccam add`; kapag may mga MCP server ang `~/.claude.json`, pinapangalanan ang mga ito ng unang `ccam add` sa tip nito.

## Paggamit

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Sa terminal, kapag walang item flag at walang `--yes`, magpapakita muna ang ccam ng checklist:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Tino-toggle ng Space (o `x`) ang item na nasa ilalim ng cursor, naglilipat ang mga arrow key o `j`/`k`, ipinapakita ng Enter ang plano, at nagpi-print ng `Cancelled.` ang Esc, `q` o Ctrl+C. Minamarkahan ng `not found` ang mga item na walang laman sa source. Kapag in-uncheck mo ang lahat, ipi-print ng ccam ang `Nothing selected.` at hihinto.

Pagkatapos ay ipi-print ng ccam ang plano at magtatanong bago ito kumopya:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Umaabot lang ang isang group sa sarili nitong mga account, at para lang sa mga item na shine-share nito, kaya inililista ng buod ang mga account ng group. Kapag hindi shine-share ng group ang isa sa mga kinopyang item, o kapag wala rito ang ilang account, ipi-print ng ccam ang mga command na `ccam group set <group> --add <items>` at `ccam group join <group> <accounts>` na nagbabago niyon.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `[folder]` | Ang Claude Code folder na pagkokopyahan | `~/.claude` |
| `--settings` | Kopyahin ang `settings.json` | off |
| `--memory` | Kopyahin ang global na `CLAUDE.md` | off |
| `--extensions` | Kopyahin ang `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | off |
| `--mcp` | Kopyahin ang mga user MCP server (`mcpServers` sa `.claude.json`) | off |
| `--history` | Kopyahin ang kasaysayan ng usapan (`projects/`, `history.jsonl`) | off |
| `-g`, `--group string` | Ang share group na kokopyahan | `default` |
| `-y`, `--yes` | Kumopya nang hindi nagtatanong; kapag may conflict, pinapanatili ang kopya ng group | off |
| `--dry-run` (global) | I-print ang plano at huminto bago ang anumang tanong o pagkopya | off |

Kapag walang item flag, ang checklist ang nagpapasya sa terminal. Sa labas ng terminal, o kapag may `--yes`, kinokopya ng ccam ang lahat maliban sa kasaysayan ng usapan. Hindi kasama ang kasaysayan bilang default dahil puwedeng malaki ang `projects/`.

Ang ibig sabihin ng `--group on` ay `default`. Tinatanggihan ang `--group off` at `--group local` na may `--group takes a group name, not "off"`. Sinusunod ng `ccam migrate` ang `--dry-run`, at lumalabas pa rin ang checklist sa dry run para ma-preview mo ang isang pagpili.

## Paano ito gumagana

Pinipili ng `internal/cli/migrate.go` ang mga item, saka binabasa ng `migrate.Build` ang source at ang group at nagbabalik ng plano. Walang isinusulat hanggang `Plan.Apply`, pagkatapos mong kumpirmahin.

1. **Hanapin ang source.** Tinatanggihan ng `checkSource` ang path na wala, hindi folder, o walang kilalang entry ng Claude Code. Para sa mga MCP server lang, ayos lang ang walang laman o nawawalang `~/.claude` basta may mga server ang `~/.claude.json`. Hinahanap ng `migrate.ClaudeJSON` ang `.claude.json` na kabilang sa folder: sa loob nito, o para sa `~/.claude` na walang ganito sa loob, ang `~/.claude.json` na katabi nito. Doon ito iniimbak ng Claude Code kapag hindi nakatakda ang `CLAUDE_CONFIG_DIR`.
2. **Itugma ang mga item sa mga path.** Galing ang bawat path sa `share.Layout`, ang parehong listahan ng pinapayagan na ginagamit ng mga group. Binabasa ang mga MCP server mula sa `.claude.json` gamit ang `mcp.LocalServers`, na gumagamit ng scanner na `profile.TopLevelValue` sa halip na i-parse ang buong file.
3. **Kopyahin, huwag kailanman ilipat.** Napupunta ang mga file sa `~/.local/share/ccam/groups/<group>/` gamit ang `fsutil.AtomicWrite` o `fsutil.CopyFile`, na sumusulat ng pansamantalang file katabi ng target at saka pinapalitan ang pangalan nito, kaya hindi kailanman nag-iiwan ng putol na file ang pagkopyang huminto. Napupunta ang mga MCP server sa `mcp.json` ng group sa ilalim ng lock nito: nananatili ang server na idinagdag ng ibang ccam habang nangyayari ito, at ina-update pagkatapos ang mga kopyang ginawa para sa [Remote Control](/guides/remote-control).

Itinutugma ng bawat item ang isang source path sa isang group path:

<!-- motion:migrate-items -->

Putol-putol ang dalawang hilera ng kasaysayan dahil kinokopya lang ang mga ito kapag may `--history`. Nakalista ang `.credentials.json` at `.claude.json` sa `share.Never` at hindi kailanman umaabot sa group: nagbabalik ang `Build` ng internal error kung lumabas man ang isa sa layout, at nilalaktawan ang mga ito ayon sa pangalan habang dinadaanan ang folder.

Binibigyan ng `buildEntry` ng estado ang bawat path, at iyon ang inilalarawan ng bawat linya ng plano:

<!-- motion:migrate-states -->

- **Missing**: wala sa source, o file kung saan inaasahan ang folder. Inililista nang isang beses bawat item bilang `not found, skipped`.
- **Linked**: tumuturo na sa kopya ng group ang source path, halimbawa kapag folder ng ccam profile ang source.
- **Same**: may kaparehong nilalaman na ang group.
- **New**: wala ito sa group, kaya kinokopya nang buo.
- **Merge**: folder sa magkabilang panig. Idinadagdag lang ng `copyMissing` ang mga file na wala sa group; nananatili ayon sa group ang mga file na magkaiba at sinasabi ng plano ang `group copies kept`. Para sa `plugins/`, pinagsasama ang mga registry na `installed_plugins.json` at `known_marketplaces.json`, at isinusulat ulit ang mga absolute na plugin path para tumuro sa folder ng group.
- **Conflict**: file sa magkabilang panig na magkaiba ang nilalaman. Ikaw ang pipili bawat file: `keep group`, `use mine`, at `merge JSON` para sa mga `.json` file (isang deep merge kung saan nananaig ang group sa nagbabanggaang value) o `merge by line` para sa `history.jsonl`. Conflict din ang MCP server na naiiba ang definition, na may `keep group` o `use mine`. Ang unang nakalistang pagpipilian ang default kapag pinindot mo ang Enter; laging pinapanatili ng `--yes` ang kopya ng group.

Ang checklist mismo ay ang `tui.MigratePicker` na binuo sa component na `checklist` na ginagamit ng `ccam add`; minamarkahan ng `migrate.Found` ang mga nawawalang item gamit ang mga lstat call at ang listahan ng MCP server.

### Mga MCP sign-in

Iniimbak ng mga MCP server na nagsa-sign in gamit ang OAuth ang kanilang mga token sa login file, ang `.credentials.json` sa ilalim ng `mcpOAuth`, hindi sa listahan ng server. Hindi kinokopya ng ccam ang mga ito. Hawak ng kinopyang sign-in ang parehong refresh token ng orihinal, at pinapalitan ng server ang token na iyon sa unang pag-renew: nala-log out na sa server na iyon ang bawat ibang may hawak, kasama ang source na `~/.claude`. Kaya pinapangalanan ng linya ng plano na `sign-ins`, para sa bawat kinopyang server na naka-sign in ang source, ang mga account ng group na magsa-sign in nang mag-isa gamit ang `/mcp`, at inuulit ito ng buod. Hindi sign-in ang entry na walang token at hindi ito inililista, gayundin ang account na may sarili nang sign-in, o ang server na ang definition ng group ang pinapanatili. Walang login file na isinusulat.

## Kaligtasan at mga espesyal na kaso

- Binabasa lang ang source. Walang file doon na inililipat, nili-link o binabago, at hindi kailanman sinusulatan ang login ng source, kasama sa macOS kung saan baka binabasa ito mula sa Keychain.
- Bago palitan ang isang group file, dahil sa `use mine`, sa isang merge, o sa mga bagong MCP server sa umiiral na `mcp.json`, sine-save ang kopya ng group sa `<group>/.ccam-backup/<time>/`. Natatangi bawat run ang pangalan ng folder, kaya hindi nagpapatungan ang mga backup ng dalawang run sa parehong segundo. Ipi-print ng buod ang path.
- Sa `--yes`, pinapanatili ng bawat conflict ang kopya ng group, at sinasabi ng ccam kung paano pumili sa halip: patakbuhin ang parehong command nang walang `--yes`.
- Kapag walang magbabago, ipi-print ng ccam ang `Nothing to copy.`; kapag pinapanatili ng bawat conflict ang kopya ng group, `Nothing to copy: the group keeps its own copies.`
- Kapag pumalya ang pagkopya sa kalagitnaan, humihinto ang ccam na may `The source folder was not changed. Group files replaced so far are saved in` kasunod ang backup folder. Kinokopya ng muling pagpapatakbo ang mga file na hindi dumating nang buo.
- Kapag ilang group ang nagshe-share ng MCP, `ccam migrate --mcp --group <a|b>` ang nakasulat sa mga command na iminumungkahi ng `ccam doctor`, `ccam ext` at `ccam mcp list`.
- Nagtatapos ang nawawalang source na may `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Kaugnay

- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Pagsusuri at pagkukumpuni](/features/doctor)
- [Maraming Claude Code account](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
