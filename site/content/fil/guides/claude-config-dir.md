---
title: Ipinaliwanag ang CLAUDE_CONFIG_DIR
description: Ano ang ginagawa ng CLAUDE_CONFIG_DIR sa Claude Code, ano ang naka-store sa folder, saan itinatago ng macOS ang login, at aling file ang puwedeng i-share.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: dc2fe42e11e1
---

Sinasabi ng `CLAUDE_CONFIG_DIR` sa Claude Code kung aling folder ang gagamitin sa halip na `~/.claude`. Naroon ang lahat ng partikular sa account: ang login, settings, mga MCP server, plugin at kasaysayan ng usapan. Kapag itinuro mo ito sa magkakaibang folder, magkakaroon ka ng magkakahiwalay na account.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Opisyal bang suportado ang CLAUDE_CONFIG_DIR?

Oo. Nakalista ito sa mga [environment variable](https://code.claude.com/docs/en/env-vars) ng Claude Code bilang paraan para palitan ang configuration folder, "kapaki-pakinabang sa pagpapatakbo ng maraming account nang magkatabi", at ginagamit ito ng [gabay nito sa authentication](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) para mag-log in gamit ang maraming account. Puwede mo itong itakda sa shell mo o sa user settings, pero hindi sa settings ng isang project.

## Ano ang nasa loob ng folder?

Sinuri gamit ang Claude Code 2.1.309:

| Path | Nilalaman | Puwedeng i-share |
| --- | --- | --- |
| `.credentials.json` | Mga login token (Linux, Windows), mga MCP OAuth token | Hindi |
| `.claude.json` | Detalye ng account, mga user MCP server, state ng bawat project | Hindi (hiwalay ang mga MCP server) |
| `settings.json` | Mga permission, hook, model, status line, plugin | Oo |
| `CLAUDE.md` | Mga global na instruction | Oo |
| `agents/`, `commands/`, `skills/`, `plugins/` | Mga extension | Oo |
| `projects/` | Mga usapan bilang mga `.jsonl` file, memory ng bawat project | Oo |
| `history.jsonl` | Kasaysayan ng mga prompt | Oo |
| `sessions/`, `shell-snapshots/`, `cache/` at iba pa | Runtime state | Hindi |

## Ano ang madalas ikalito?

**Kapag walang `CLAUDE_CONFIG_DIR`, wala sa loob ng `~/.claude` ang `.claude.json`.** Sa ganoong kaso, iniimbak ito ng Claude Code sa `~/.claude.json`, katabi ng folder. Kung sisimulan mong gamitin ang `CLAUDE_CONFIG_DIR=~/.claude`, hahanapin ng Claude Code ang `~/.claude/.claude.json` sa halip at hindi nito mahahanap ang account mo. Inaayos ito ng `ccam import` sa pamamagitan ng pag-link sa dalawa. Kung may sarili nang `.claude.json` ang `~/.claude`, mula sa naunang pagpapatakbo gamit ang `CLAUDE_CONFIG_DIR=~/.claude`, hinahayaan ng ccam ang parehong file at sinasabi ito, dahil ikaw lang ang nakakaalam kung alin ang tama.

**Sa macOS, nasa Keychain ang login, wala sa folder.** Iniimbak ito ng Claude Code sa ilalim ng pangalang `Claude Code-credentials`. Kapag nakatakda ang `CLAUDE_CONFIG_DIR`, nagtatapos din ang pangalan sa `-` at sa unang 8 hex digit ng SHA-256 ng path ng folder. `$USER` ang account name ng Keychain item. Kaya magkakahiwalay na Keychain entry ang nakukuha ng magkakahiwalay na folder, at nawawala ang login ng folder na inilipat o pinalitan ng pangalan; iyan ang dahilan kung bakit pinapanatili ng `ccam rename` ang pangalan ng folder sa macOS. Kung tinanggihan ng Keychain ang pagsulat, babalik ang Claude Code sa `.credentials.json` sa folder, gaya ng inilalarawan sa [dokumentasyon nito](https://code.claude.com/docs/en/authentication#credential-management).

**Hindi lahat ng sign-in ay nasa folder.** Ang Claude Console sign-in na walang API key ay naka-store sa labas nito, kaya hindi kayang paghiwalayin ng magkakahiwalay na folder ang dalawang ganoong sign-in. Ang mga claude.ai login at API key ay naka-store bawat folder.

## Nasaan ang folder sa Windows?

Bilang default, ito ay `%USERPROFILE%\.claude`. Nasa `.credentials.json` sa loob nito ang login, o nasa Windows Credential Manager kung nakatakdang gamitin ito ng Claude Code, sa ilalim ng pangalang hinango mula sa path ng folder gaya ng sa macOS. Itakda ang `CLAUDE_CONFIG_DIR` sa ibang folder para sa bawat karagdagang account, gaya ng sa Linux at macOS.

## Puwede bang mag-share ng file ang mga folder?

Gumagana ang mga symbolic link: nagbabasa at nagsusulat ang Claude Code sa pamamagitan ng naka-link na `settings.json` o `projects/` folder at pinapanatili ang link. Eksepsiyon ang mga MCP server, dahil kasama nila ang login sa `.claude.json`; ipasa sila gamit ang `--mcp-config` sa halip.

Nakabatay ang ccam sa mismong mga patakarang ito. Tingnan ang [sharing](/docs/sharing) para sa kung paano nito nili-link ang mga folder, at ang [gabay sa maraming account](/guides/multiple-claude-code-accounts) para sa mas malawak na larawan.
