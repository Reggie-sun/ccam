---
title: I-share ang mga MCP server, settings at kasaysayan
description: Paano shine-share ng ccam ang MCP server, settings.json, CLAUDE.md, agent, skill at kasaysayan ng usapan ng mga Claude Code account; hiwalay ang bawat login.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 0df8a7d51c87
---

May sariling login ang bawat account. Lahat ng iba pa ay puwedeng i-share nang isa-isang item: nili-link ng ccam ang file o folder ng account sa isang shared na group folder, kaya ang pagbabagong ginawa sa isang account ay lumalabas sa lahat ng account ng group na iyon.

## Ano ang puwedeng i-share

| Item | Option | Ano ang shine-share |
| --- | --- | --- |
| Mga MCP server | `--mcp` | Ang `mcp.json` ng group, na ipinapasa sa Claude Code sa pagsisimula |
| Settings | `--settings` | `settings.json`: mga permission, hook, model, status line |
| Global na memory | `--memory` | `CLAUDE.md` |
| Mga extension | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Kasaysayan ng usapan | `--history` | `projects/` at `history.jsonl` |

Shine-share ng mga bagong account ang lahat maliban sa kasaysayan ng usapan.

Hindi kailanman shine-share ang mga login file na `.credentials.json` at `.claude.json`. Nagli-link lang ang ccam ng mga path mula sa isang nakapirming listahan, at wala sa listahang iyon ang dalawang ito.

## Mga group

Ang ibig sabihin ng `on` ay ang `default` group. Puwede ka ring magbigay ng pangalan, para iilang account lang ang mag-share ng isang item:

```sh
ccam share work --history team
ccam share side --history team
```

Nagshe-share na ngayon ng mga usapan ang `work` at `side` sa isa't isa, pero hindi sa iba mo pang account.

## Pag-on ng sharing

Kapag nagsimulang mag-share ng item ang isang account, ang ccam ay:

1. nagpi-print ng plano ng bawat file na ililipat, kokopyahin o ili-link nito, at nagtatanong bago magsimula;
2. inililipat ang mga file ng account sa group folder kung wala pang laman ang group, o idinadagdag ang mga file na wala pa sa group;
3. kapag magkaiba ang parehong file, nagtatanong kung pananatilihin ang bersyon ng group, gagamitin ang sa iyo, o pagsasamahin ang JSON;
4. iniimbak ang orihinal mo sa `.ccam-backup/<time>/` sa loob ng account folder;
5. pinapalitan ang orihinal ng link sa group.

Ipinapangalan ang mga file ng usapan ayon sa kanilang session ID, kaya nagsasama ang mga kasaysayan mula sa iba't ibang account nang walang conflict.

Hindi binabago ng ccam ang sharing habang tumatakbo ang Claude Code bilang ang account, dahil patuloy na nagsusulat ang session sa mga file na inililipat. Isara muna ito.

## Pag-off ng sharing

```sh
ccam share work --history off
```

Inaalis ang link at nakakakuha ang account ng sarili nitong kopya ng kasalukuyang nilalaman ng group. Idagdag ang `--fresh` para magsimulang walang laman sa halip. Pinapanatili ng group ang nilalaman nito para sa ibang account.

## Kopyahin ang dati nang setup

Kung gumamit ka ng Claude Code bago ang ccam, nasa `~/.claude` ang iyong settings. Kinokopya ng `ccam migrate` ang mga ito sa isang group, para makuha ng bawat account na nagshe-share ng mga item na iyon:

```sh
ccam migrate
```

Sa terminal, magpapakita muna ang ccam ng checklist ng mga item. Naka-check ang lahat maliban sa kasaysayan ng usapan; pindutin ang Space para baguhin ang isang item at Enter para makita ang plano. Minamarkahan ng `not found` ang mga item na walang laman sa folder.

Ang kinokopya nito:

- ang `settings.json` at ang global na `CLAUDE.md`;
- ang mga extension na `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` at `rules/`;
- ang iyong mga user MCP server, ang entry na `mcpServers` ng `~/.claude.json`, papunta sa `mcp.json` ng group;
- ang mga naka-save na sign-in (OAuth) ng mga MCP server na kinokopya nito, papunta sa sariling login file ng bawat account;
- ang kasaysayan ng usapan (`projects/` at `history.jsonl`), pero kapag idinagdag mo lang ang `--history`.

Binabasa lang ang source folder: walang inililipat, nili-link o binabago rito, at gumagana pa rin ang karaniwang `claude` gaya ng dati. Hindi kailanman kinokopya sa group ang mga login file na `.credentials.json` at `.claude.json`.

Iniimbak ng mga MCP server na nagsa-sign in gamit ang OAuth ang kanilang mga token sa login, hindi sa listahan ng server. Para hindi ka na mag-sign in ulit sa mga ito, kinokopya ng ccam ang mga sign-in na iyon sa login file ng bawat account na nakakakuha ng mga server sa pamamagitan ng group, maliban kung may sarili na ang account. Nilalaktawan nito ang mga account na tumatakbo, hindi pa naka-log in, nirehistro sa kinalalagyan nito gamit ang `ccam import`, o nag-iimbak ng login sa macOS Keychain, at pinapangalanan ang mga ito; doon, mag-sign in sa mga server gamit ang `/mcp`. Ang account na magsisimulang mag-share ng mga MCP server sa ibang pagkakataon ay makakakuha ng mga sign-in kapag pinatakbo mo ulit ang `ccam migrate`.

Bago kumopya ng anuman, ipi-print ng ccam ang plano at magtatanong:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Sa dulo, inililista ng ccam ang mga account na nagshe-share ng bawat item sa pamamagitan ng group. Kung may mga account na hindi pa, ipi-print nito ang mga command na `ccam share` na mag-o-on ng sharing para sa kanila.

Kapag may iba nang kopya ng isang file sa group, magtatanong ang ccam, gaya ng pag-on ng sharing: pananatilihin ang kopya ng group o gagamitin ang sa iyo. Puwede ring pagsamahin ang `settings.json` at iba pang JSON file, at ang `history.jsonl` nang linya por linya. Sine-save ang mga dating kopya ng group sa `.ccam-backup/<time>/` sa loob ng group folder.

| Option | Epekto |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Kopyahin lang ang mga item na pinangalanan mo, nang walang checklist; kung wala ang mga ito at nasa labas ng terminal, lahat maliban sa kasaysayan ng usapan |
| `-g`, `--group <group>` | Kopyahin sa group na ito sa halip na sa `default` |
| `-y`, `--yes` | Kopyahin nang hindi nagtatanong; kapag magkaiba ang isang file, panatilihin ang kopya ng group |

Para makita ang plano nang walang kinokopya, ilagay ang `--dry-run` bago ang command, gaya ng `ccam --dry-run migrate`. Para kumopya mula sa ibang folder, pangalanan ito: `ccam migrate ~/.claude-work --group work`.

## Mga MCP server

Iniimbak ng Claude Code ang mga MCP server sa `.claude.json`, kasama ng login, kaya hindi puwedeng i-share ang file na iyon. Iniimbak ng ccam ang mga shared na server sa `mcp.json` ng group at sinisimulan ang Claude Code gamit ang `--mcp-config=<group>/mcp.json`. Umiiral pa rin sa ibabaw nito ang sarili mong mga `--mcp-config` option.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Ang server na idinagdag mo sa loob ng session gamit ang `claude mcp add` ay naka-save para sa account na iyon lang. Sa susunod na simulan mo ang account na iyon, babanggitin ng ccam ang server at ang command na maglilipat nito sa group:

```sh
ccam mcp import work gh
```

Kung may ibang server na ang group na may ganoong pangalan, nilalaktawan ito ng ccam at sinasabi ito. Pinapalitan ng `--force` ang server ng group at itinatago ang lumang `mcp.json` sa `.ccam-backup/` folder ng group.

Iniimbak ng mga MCP server na nagsa-sign in gamit ang OAuth ang kanilang mga token kasama ng login ng account, kaya i-authorize sila nang isang beses sa bawat account, o kopyahin ang mga sign-in mula sa `~/.claude` gamit ang `ccam migrate`.

Kung gusto mong lumabas ang mga server sa `claude mcp list` at `/mcp` bilang karaniwang user server, itakda ang `mcp_strategy = "sync"` sa `~/.config/ccam/config.toml`. Isusulat na ng ccam ang mga ito sa `.claude.json` ng bawat account bago magsimula, sa ilalim ng sariling file lock ng Claude Code at may backup.

## Dalawang account sa iisang folder

Sumusulat sa parehong mga file ang mga account na nagshe-share ng kasaysayan ng usapan. Kung magsisimula ka ng pangalawang Claude Code sa isang folder kung saan tumatakbo na ang isa pang account ng parehong group, babalaan ka ng ccam at tatanungin kung magpapatuloy.

## Kapag may nasira

```sh
ccam doctor
ccam doctor --fix
```

Kung naging karaniwang file ang isang shared na path, halimbawa pagkatapos mong i-restore nang mano-mano ang isang backup, hihinto ang ccam bago simulan ang Claude Code at ituturo ang `ccam doctor --fix`, na pinagsasama ang file sa group at nag-iimbak ng backup.
