---
title: Mga API key account
description: Nagdadagdag ang ccam add --api ng API key account para sa API ng Anthropic o isang gateway; binabasa ang key sa pagsisimula at hindi iniimbak.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Sinisimulan ng isang API key account ang Claude Code gamit ang API key sa halip na Claude subscription login: ang API ng Anthropic, o isang endpoint na gumagamit ng parehong Messages format, gaya ng LiteLLM o OpenRouter gateway. Hindi iniimbak ng ccam ang key. Binabasa nito ang key sa bawat pagsisimula ng account, mula sa isang environment variable o sa isang command gaya ng password manager, at ibinibigay lang ito sa prosesong iyon ng Claude Code.

## Paggamit

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

Nagshe-share ang isang API key account ng settings, `CLAUDE.md`, mga extension at mga MCP server gaya ng ibang account, at gumagana ito sa mga shortcut, sa `ccam folder` at sa `ccam ext`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--api` | Magdagdag ng API key account sa halip na mag-log in | off |
| `--key-from env:NAME` | Basahin ang key mula sa environment variable na `NAME` | |
| `--key-from command:COMMAND` | Patakbuhin ang `COMMAND` (gamit ang `sh -c`, o `cmd /C` sa Windows) at gamitin ang unang linyang ipi-print nito | |
| `--base-url <url>` | Ang endpoint (`ANTHROPIC_BASE_URL`) | Ang API ng Anthropic |
| `--auth api-key\|bearer` | Ipadala ang key bilang `ANTHROPIC_API_KEY` (`x-api-key`) o bilang `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, na karaniwang hinihingi ng mga gateway) | `api-key` |
| `--model <model>` | Ang model (`ANTHROPIC_MODEL`) | Default ng Claude Code |

Iniimbak ang mga setting sa `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Paano ito gumagana

Sa bawat pagsisimula ng account, pagkatapos ng karaniwang mga pagsusuri at bago magsimula ang Claude Code:

1. Binabasa ng ccam ang key: ang `env:` mula sa environment, ang `command:` sa pagpapatakbo ng command na may 30-segundong limit. Puwedeng humingi ng password sa terminal ang command; napupunta rin sa terminal ang error output nito. Pinapahinto ng walang laman na key ang pagsisimula na may mensaheng nagsasabi kung saan inaasahan ang key.
2. Inaalis nito ang bawat minanang variable na nagpapasya kung paano magsa-sign in ang Claude Code o kung aling provider ang gagamitin nito: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` at `CLAUDE_CODE_USE_FOUNDRY`. Hindi makakarating sa account na ito ang key na itinakda sa shell mo para sa ibang bagay, at hindi matatalo ng subscription login ang key nito.
3. Itinatakda nito ang `ANTHROPIC_BASE_URL` at `ANTHROPIC_MODEL` kapag ibinigay, at ang key bilang `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`.

Pinananatili ng mga account na may subscription login ang environment na minana nila, gaya ng dati. Sa isang interactive session, puwedeng magtanong nang isang beses ang Claude Code kung gagamitin ang key na nakita nito sa environment.

## Kaligtasan at mga espesyal na kaso

- Hindi kailanman isinusulat ang key sa `config.toml`, sa cache o sa log, at hindi ito binabasa ng `ccam --dry-run gw`: hindi pinapatakbo ng dry run ang `command:` na source.
- Agad na sinusuri ng `ccam add --api` ang `env:` na source at nagbababala kapag walang laman ang variable; unang pinapatakbo ang `command:` na source sa pagsisimula.
- Tinatanggihan ng `ccam login` at `ccam logout` ang API key account, dahil wala itong login. Ipinapakita ng `ccam doctor` ang endpoint at key source nito, at nagbababala kapag walang laman ang `env:` na variable.
- Ipinapakita ng `ccam usage` ang `API` para rito. Pumipili lang ang `ccam next --auto`, `--fallback` at `ccam pick` sa mga subscription account, dahil walang 5-oras o lingguhang limit ang API key na maikukumpara.
- Nalalapat sa bawat account sa group, kasama ang mga API key account, ang shared na `settings.json` na may `env` block o `apiKeyHelper`. Huwag ilagay ang mga setting ng provider sa shared na settings.

## Kaugnay

- [Pamamahala ng mga account](/features/accounts)
- [Pagsisimula ng Claude Code](/features/launch)
- [Sanggunian ng mga command](/docs/commands)
