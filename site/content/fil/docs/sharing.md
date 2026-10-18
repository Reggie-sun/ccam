---
title: I-share ang mga MCP server, settings at kasaysayan
description: Paano shine-share ng ccam ang MCP server, settings.json, CLAUDE.md, agent, skill at kasaysayan ng usapan ng mga Claude Code account; hiwalay ang bawat login.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
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

## Pag-off ng sharing

```sh
ccam share work --history off
```

Inaalis ang link at nakakakuha ang account ng sarili nitong kopya ng kasalukuyang nilalaman ng group. Idagdag ang `--fresh` para magsimulang walang laman sa halip. Pinapanatili ng group ang nilalaman nito para sa ibang account.

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

Iniimbak ng mga MCP server na nagsa-sign in gamit ang OAuth ang kanilang mga token kasama ng login ng account, kaya i-authorize sila nang isang beses sa bawat account.

Kung gusto mong lumabas ang mga server sa `claude mcp list` at `/mcp` bilang karaniwang user server, itakda ang `mcp_strategy = "sync"` sa `~/.config/ccam/config.toml`. Isusulat na ng ccam ang mga ito sa `.claude.json` ng bawat account bago magsimula, sa ilalim ng sariling file lock ng Claude Code at may backup.

## Dalawang account sa iisang folder

Sumusulat sa parehong mga file ang mga account na nagshe-share ng kasaysayan ng usapan. Kung magsisimula ka ng pangalawang Claude Code sa isang folder kung saan tumatakbo na ang isa pang account ng parehong group, babalaan ka ng ccam at tatanungin kung magpapatuloy.

## Kapag may nasira

```sh
ccam doctor
ccam doctor --fix
```

Kung naging karaniwang file ang isang shared na path, halimbawa pagkatapos mong i-restore nang mano-mano ang isang backup, hihinto ang ccam bago simulan ang Claude Code at ituturo ang `ccam doctor --fix`, na pinagsasama ang file sa group at nag-iimbak ng backup.
